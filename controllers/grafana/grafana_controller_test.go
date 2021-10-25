/*
Copyright 2021.

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

package grafana

import (
	"context"
	"time"

	integreatlyorgv1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Grafana controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		GrafanaName        = "test-grafana"
		GrafanaServiceName = "grafana-service"
		GrafanaNamespace   = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating Grafana Status", func() {
		It("Should set Grafana Status.message to success", func() {
			By("setting up a new Grafana instance")
			ctx := context.Background()
			grafana := &integreatlyorgv1alpha1.Grafana{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "integreatly.org/v1alpha1",
					Kind:       "Grafana",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      GrafanaName,
					Namespace: GrafanaNamespace,
				},
				Spec: integreatlyorgv1alpha1.GrafanaSpec{
					DashboardLabelSelector: []*metav1.LabelSelector{
						{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values: []string{
										"grafana",
									},
								},
							},
						},
					},
					Service: &integreatlyorgv1alpha1.GrafanaService{
						Name: "grafana-service",
						Labels: map[string]string{
							"app":  "grafana",
							"type": "grafana-service",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, grafana)).Should(Succeed())

			grafanaLookupKey := types.NamespacedName{Name: GrafanaName, Namespace: GrafanaNamespace}
			//grafanaServiceLookupKey := types.NamespacedName{Name: GrafanaName, Namespace: GrafanaNamespace}
			createdGrafana := &integreatlyorgv1alpha1.Grafana{}

			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, grafanaLookupKey, createdGrafana)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(createdGrafana.ObjectMeta.Name).Should(Equal("test-grafana"))
			//By("By checking the service that got created")

			/*
				createdService := &apiv1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      GrafanaServiceName,
						Namespace: GrafanaNamespace,
					},
					Spec: apiv1.ServiceSpec{
						Ports: []apiv1.ServicePort{
							{
								Protocol:   "TCP",
								Port:       3000,
								TargetPort: intstr.FromString("grafana-http"),
							},
						},
					},
				}
			*/

			//createdService := &apiv1.Service{}
			//podList := &corev1.PodList{}
			createdService := &corev1.ServiceList{}

			Eventually(func() bool {
				//err := k8sClient.Get(ctx, grafanaServiceLookupKey, createdService)
				err := k8sClient.List(ctx, createdService, &client.ListOptions{})
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(createdService.Items).Should(Equal(GrafanaServiceName))

		})
	})
})
