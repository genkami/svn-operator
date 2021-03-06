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
		interval = 250 * time.Millisecond
	)

	defaultSVNServer := func() *svnv1alpha1.SVNServer {
		return &svnv1alpha1.SVNServer{
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
	}

	Context("When SVNServer Status is created", func() {
		It("Creates corresponding resources automatically", func() {
			By("creating a new SVNServer")
			ctx := context.Background()
			svnServer := defaultSVNServer()
			Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
			defer func() {
				Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
			}()

			By("checking whether the StatefulSet is created")
			statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
			statefulSet := &appsv1.StatefulSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			defer func() {
				Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
			}()

			By("checking whether SVNServer.Status.Conditions is updated")
			svnServerLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
			createdSVNServer := &svnv1alpha1.SVNServer{}
			Eventually(func() (int, error) {
				err := k8sClient.Get(ctx, svnServerLookupKey, createdSVNServer)
				if err != nil {
					return -1, err
				}
				return len(createdSVNServer.Status.Conditions), nil
			}, timeout, interval).Should(BeNumerically(">=", 1))
			conds := createdSVNServer.Status.Conditions
			Expect(conds[len(conds)-1].Type).To(Equal(svnv1alpha1.ConditionTypeSynced))
		})
	})

	Describe(".Spec.PodTemplate.Image", func() {
		Context("when the field is not set", func() {
			It("uses the default value", func() {
				ctx := context.Background()
				svnServer := defaultSVNServer()
				Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
				defer func() {
					Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
				}()

				statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
				statefulSet := &appsv1.StatefulSet{}
				Eventually(func() (string, error) {
					err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
					if err != nil {
						return "", err
					}
					return statefulSet.Spec.Template.Spec.Containers[0].Image, nil
				}, timeout, interval).Should(Equal(defaultSVNServerImageForTest))
				defer func() {
					Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
				}()
			})
		})

		Context("when the field is set", func() {
			It("uses the given value", func() {
				ctx := context.Background()
				svnServer := defaultSVNServer()
				svnServer.Spec.PodTemplate.Image = "my-image:latest"
				Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
				defer func() {
					Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
				}()

				statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
				statefulSet := &appsv1.StatefulSet{}
				Eventually(func() (string, error) {
					err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
					if err != nil {
						return "", err
					}
					return statefulSet.Spec.Template.Spec.Containers[0].Image, nil
				}, timeout, interval).Should(Equal("my-image:latest"))
				defer func() {
					Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
				}()
			})
		})
	})

	Describe(".Spec.PodTemplate.NodeSelector", func() {
		Context("when the field is not set", func() {
			It("uses the default value", func() {
				ctx := context.Background()
				svnServer := defaultSVNServer()
				Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
				defer func() {
					Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
				}()

				statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
				statefulSet := &appsv1.StatefulSet{}
				Eventually(func() (map[string]string, error) {
					err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
					if err != nil {
						return nil, err
					}
					return statefulSet.Spec.Template.Spec.NodeSelector, nil
				}, timeout, interval).Should(BeZero())
				defer func() {
					Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
				}()

			})
		})

		Context("when the field is set", func() {
			It("uses the given value", func() {
				ctx := context.Background()
				svnServer := defaultSVNServer()
				svnServer.Spec.PodTemplate.NodeSelector = map[string]string{
					"some-label": "some-value",
				}
				Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
				defer func() {
					Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
				}()

				statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
				statefulSet := &appsv1.StatefulSet{}
				Eventually(func() (map[string]string, error) {
					err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
					if err != nil {
						return nil, err
					}
					return statefulSet.Spec.Template.Spec.NodeSelector, nil
				}, timeout, interval).Should(Equal(map[string]string{
					"some-label": "some-value",
				}))
				defer func() {
					Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
				}()
			})
		})
	})

	Describe(".Spec.PodTemplate.ServiceAccountName", func() {
		Context("when the field is not set", func() {
			It("uses the default value", func() {
				ctx := context.Background()
				svnServer := defaultSVNServer()
				Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
				defer func() {
					Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
				}()

				statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
				statefulSet := &appsv1.StatefulSet{}
				Eventually(func() (string, error) {
					err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
					if err != nil {
						return "", err
					}
					return statefulSet.Spec.Template.Spec.ServiceAccountName, nil
				}, timeout, interval).Should(BeZero())
				defer func() {
					Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
				}()
			})
		})

		Context("when the field is set", func() {
			It("uses the given value", func() {
				ctx := context.Background()
				svnServer := defaultSVNServer()
				svnServer.Spec.PodTemplate.ServiceAccountName = "my-account"
				Expect(k8sClient.Create(ctx, svnServer)).To(Succeed())
				defer func() {
					Expect(k8sClient.Delete(ctx, svnServer)).To(Succeed())
				}()

				statefulSetLookupKey := types.NamespacedName{Name: SVNServerName, Namespace: SVNServerNamespace}
				statefulSet := &appsv1.StatefulSet{}
				Eventually(func() (string, error) {
					err := k8sClient.Get(ctx, statefulSetLookupKey, statefulSet)
					if err != nil {
						return "", err
					}
					return statefulSet.Spec.Template.Spec.ServiceAccountName, nil
				}, timeout, interval).Should(Equal("my-account"))
				defer func() {
					Expect(k8sClient.Delete(ctx, statefulSet)).To(Succeed())
				}()
			})
		})
	})

})
