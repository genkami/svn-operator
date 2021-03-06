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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	svnv1alpha1 "github.com/genkami/svn-operator/api/v1alpha1"
)

var _ = Describe("SVNServer Controller", func() {
	const (
		SVNServerName      = "test-svnserver"
		SVNServerNamespace = "default"

		timeout  = 10 * time.Second
		duration = 10 * time.Second
		interval = 250 * time.Millisecond
	)

	Context("When updating SVNServer Status", func() {
		It("Updates its conditions", func() {
			By("creating a new SVNServer")
			ctx := context.Background()
			svnServer := &svnv1alpha1.SVNServer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "svn.k8s.oyasumi.club/v1alpha1",
					Kind:       "SVNServer",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      SVNServerName,
					Namespace: SVNServerNamespace,
				},
				Spec: svnv1alpha1.SVNServerSpec{
					VolumeClaimTemplate: corev1.PersistentVolumeClaim{
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{
								corev1.ReadWriteOnce,
							},
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceStorage: resource.MustParse("1G"),
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())

			statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
			statefulSet := &appsv1.StatefulSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})
