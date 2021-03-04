/*
Copyright 2021 Genta Kamitani.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	svnv1alpha1 "github.com/genkami/svn-operator/api/v1alpha1"
)

const (
	VolumeNameRepos  = "repos"
	VolumePathRepos  = "/svn"
	VolumeNameConfig = "config"
	VolumePathConfig = "/etc/svn-config/"

	ContainerNameSVN = "svn"

	LabelAppKey          = "app"
	LabelAppValue        = "subversion"
	LabelInstanceNameKey = "svn.k8s.oyasumi.club/name"

	ConfigMapKeyAuthUserFile       = "AuthUserFile"
	ConfigMapKeyAuthzSVNAccessFile = "AuthzSVNAccessFile"
	ConfigMapKeyRepos              = "Repos"
)

// SVNServerReconciler reconciles a SVNServer object
type SVNServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	// DefaultSVNServerImage is a Docker image name to run SVN server.
	DefaultSVNServerImage string
}

// +kubebuilder:rbac:groups=svn.k8s.oyasumi.club,resources=svnservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=svn.k8s.oyasumi.club,resources=svnservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=svn.k8s.oyasumi.club,resources=svnservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SVNServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *SVNServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("svnserver", req.NamespacedName)

	svnServer := &svnv1alpha1.SVNServer{}
	err := r.Get(ctx, req.NamespacedName, svnServer)
	if err != nil {
		if errors.IsNotFound(err) {
			// The object cloud have been deleted asynchronously.
			log.Info("SVNServer not found; ignoring.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get SVNServer")
		return ctrl.Result{}, err
	}

	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: svnServer.Name, Namespace: svnServer.Namespace}, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			if err = r.createService(ctx, log, svnServer); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	ss := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: svnServer.Name, Namespace: svnServer.Namespace}, ss)
	if err != nil {
		if errors.IsNotFound(err) {
			if err = r.createStatefulSet(ctx, log, svnServer); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		log.Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}

	cm := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: svnServer.Name, Namespace: svnServer.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			if err = r.createConfigMap(ctx, log, svnServer); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		log.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	desiredSS := ss.DeepCopy()
	r.overrideWithPodTemplate(svnServer, desiredSS)
	if !reflect.DeepEqual(desiredSS, ss) {
		if err := r.Update(ctx, desiredSS); err != nil {
			log.Error(err, "Failed to update StatefulSet")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	desiredCM := r.configMapFor(svnServer)
	if !reflect.DeepEqual(desiredCM.Data, cm.Data) {
		if err := r.Update(ctx, desiredCM); err != nil {
			log.Error(err, "Failed to update ConfigMap")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// TODO: Update SVNServer.Status
	return ctrl.Result{}, nil
}

// Creates a StatefulSet and is corresponding Service
func (r *SVNServerReconciler) createStatefulSet(ctx context.Context, log logr.Logger, svn *svnv1alpha1.SVNServer) error {
	ss := r.statefulSetFor(svn)
	log = log.WithValues("StatefulSet.Namespace", ss.Namespace, "StatefulSet.Name", ss.Name)
	log.Info("Creating a new StatefulSet")
	if err := r.Create(ctx, ss); err != nil {
		log.Error(err, "Failed to create new StatefulSet")
		return err
	}
	return nil
}

func (r *SVNServerReconciler) createService(ctx context.Context, log logr.Logger, svn *svnv1alpha1.SVNServer) error {
	svc := r.serviceFor(svn)
	log = log.WithValues("Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
	log.Info("Creating a new Service")
	if err := r.Create(ctx, svc); err != nil {
		log.Error(err, "Failed to create new Service")
		return err
	}
	return nil
}

func (r *SVNServerReconciler) createConfigMap(ctx context.Context, log logr.Logger, svn *svnv1alpha1.SVNServer) error {
	svc := r.configMapFor(svn)
	log = log.WithValues("ConfigMap.Namespace", svc.Namespace, "ConfigMap.Name", svc.Name)
	log.Info("Creating a new ConfigMap")
	if err := r.Create(ctx, svc); err != nil {
		log.Error(err, "Failed to create new ConfigMap")
		return err
	}
	return nil
}

func (r *SVNServerReconciler) statefulSetFor(s *svnv1alpha1.SVNServer) *appsv1.StatefulSet {
	labels := r.labelsFor(s)
	replicas := int32(1)
	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: VolumeNameRepos,
				},
			}},
			ServiceName: s.Name,
		},
	}
	r.overrideWithPodTemplate(s, ss)
	ctrl.SetControllerReference(s, ss, r.Scheme)
	return ss
}

func (r *SVNServerReconciler) overrideWithPodTemplate(s *svnv1alpha1.SVNServer, ss *appsv1.StatefulSet) {
	var volumeClaimIndex int = -1
	for i := range ss.Spec.VolumeClaimTemplates {
		pvc := &ss.Spec.VolumeClaimTemplates[i]
		if pvc.Name == VolumeNameRepos {
			volumeClaimIndex = i
			break
		}
	}
	if volumeClaimIndex < 0 {
		pvc := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: VolumeNameRepos,
			},
		}
		ss.Spec.VolumeClaimTemplates = append(ss.Spec.VolumeClaimTemplates, pvc)
		volumeClaimIndex = len(ss.Spec.VolumeClaimTemplates) - 1
	}
	ss.Spec.VolumeClaimTemplates[volumeClaimIndex] = *s.Spec.VolumeClaimTemplate.DeepCopy()
	ss.Spec.VolumeClaimTemplates[volumeClaimIndex].Name = VolumeNameRepos

	var volume *corev1.Volume
	for i := range ss.Spec.Template.Spec.Volumes {
		v := &ss.Spec.Template.Spec.Volumes[i]
		if v.Name == VolumeNameConfig {
			volume = v
			break
		}
	}
	if volume == nil {
		v := &corev1.Volume{Name: VolumeNameConfig}
		ss.Spec.Template.Spec.Volumes = append(ss.Spec.Template.Spec.Volumes, *v)
		volume = &ss.Spec.Template.Spec.Volumes[len(ss.Spec.Template.Spec.Volumes)-1]
	}
	volume.VolumeSource = corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: s.Name,
			},
		},
	}

	var container *corev1.Container
	for i := range ss.Spec.Template.Spec.Containers {
		c := &ss.Spec.Template.Spec.Containers[i]
		if c.Name == ContainerNameSVN {
			container = c
			break
		}
	}
	if container == nil {
		ss.Spec.Template.Spec.Containers = append(ss.Spec.Template.Spec.Containers, r.svnContainerFor(s))
		container = &ss.Spec.Template.Spec.Containers[len(ss.Spec.Template.Spec.Containers)-1]
	}
	if s.Spec.PodTemplate.Image != "" {
		container.Image = s.Spec.PodTemplate.Image
	} else if container.Image == "" {
		container.Image = r.DefaultSVNServerImage
	}

	if len(s.Spec.PodTemplate.NodeSelector) > 0 {
		ss.Spec.Template.Spec.NodeSelector = map[string]string{}
		for k, v := range s.Spec.PodTemplate.NodeSelector {
			ss.Spec.Template.Spec.NodeSelector[k] = v
		}
	}
	if s.Spec.PodTemplate.ServiceAccountName != "" {
		ss.Spec.Template.Spec.ServiceAccountName = s.Spec.PodTemplate.ServiceAccountName
	}
	if len(s.Spec.PodTemplate.ImagePullSecrets) > 0 {
		ss.Spec.Template.Spec.ImagePullSecrets = make([]corev1.LocalObjectReference, len(s.Spec.PodTemplate.ImagePullSecrets))
		copy(ss.Spec.Template.Spec.ImagePullSecrets, s.Spec.PodTemplate.ImagePullSecrets)
	}
	if s.Spec.PodTemplate.Affinity != nil {
		affinity := *s.Spec.PodTemplate.Affinity
		ss.Spec.Template.Spec.Affinity = &affinity
	}
	if len(s.Spec.PodTemplate.Tolerations) > 0 {
		ss.Spec.Template.Spec.Tolerations = make([]corev1.Toleration, len(s.Spec.PodTemplate.Tolerations))
		copy(ss.Spec.Template.Spec.Tolerations, s.Spec.PodTemplate.Tolerations)
	}
}

func (r *SVNServerReconciler) svnContainerFor(s *svnv1alpha1.SVNServer) corev1.Container {
	return corev1.Container{
		Name:  ContainerNameSVN,
		Image: r.DefaultSVNServerImage,
		Ports: []corev1.ContainerPort{{
			ContainerPort: 80,
			Name:          "http",
		}},
		ReadinessProbe: &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/",
					Port: intstr.FromInt(80),
				},
			},
		},
		LivenessProbe: &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/",
					Port: intstr.FromInt(80),
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      VolumeNameRepos,
				MountPath: VolumePathRepos,
			},
			{
				Name:      VolumeNameConfig,
				MountPath: VolumePathConfig,
			},
		},
	}
}

func (r *SVNServerReconciler) serviceFor(s *svnv1alpha1.SVNServer) *corev1.Service {
	labels := r.labelsFor(s)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: 80,
			}},
			Selector:  labels,
			ClusterIP: "None",
		},
	}
	ctrl.SetControllerReference(s, svc, r.Scheme)
	return svc
}

// TODO: Use SVNRepository, SVNUser, and SVNGroup
func (r *SVNServerReconciler) configMapFor(s *svnv1alpha1.SVNServer) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Data: map[string]string{
			ConfigMapKeyAuthUserFile:       r.authUserFileFor(s),
			ConfigMapKeyAuthzSVNAccessFile: r.authzSVNAccessFileFor(s),
			ConfigMapKeyRepos:              r.reposConfigFor(s),
		},
	}
	ctrl.SetControllerReference(s, cm, r.Scheme)
	return cm
}

func (r *SVNServerReconciler) authUserFileFor(s *svnv1alpha1.SVNServer) string {
	// TODO
	// admin:hogefuga (for test)
	return `
admin:{SHA}LxoHQl0nHaVtqqtU9KO/J8O75JM=
`
}

func (r *SVNServerReconciler) authzSVNAccessFileFor(s *svnv1alpha1.SVNServer) string {
	// TODO
	return `
[groups]
all = admin

[example-rw-repo:/]
* =
@all = rw

[example-r-repo:/]
* =
@all = r

[example-private-repo:/]
* =
@all =
`
}

func (r *SVNServerReconciler) reposConfigFor(s *svnv1alpha1.SVNServer) string {
	// TODO
	return `
repositories:
- name: example-rw-repo
- name: example-r-repo
- name: example-private-repo
`
}

func (r *SVNServerReconciler) labelsFor(s *svnv1alpha1.SVNServer) map[string]string {
	return map[string]string{
		LabelAppKey:          LabelAppValue,
		LabelInstanceNameKey: s.Name,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SVNServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&svnv1alpha1.SVNServer{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
