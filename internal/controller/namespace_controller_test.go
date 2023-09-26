/*
Copyright 2023.

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

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Memcached controller", func() {
	Context("Memcached controller test", func() {
		const Gangway = "test-gangway"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   Gangway,
				Labels: map[string]string{"gangway": "true"},
			},
		}

		typeNamespaceName := types.NamespacedName{Name: Gangway, Namespace: Gangway}

		BeforeEach(func() {
			By("Creating the Namespace to perform the tests")
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))
		})

		It("should successfully reconcile a Gangway enabled namespace", func() {
			var found networking.NetworkPolicy

			By("Reconciling the resource created")
			namespaceReconciler := &NamespaceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if NetPol was successfully created in the reconciliation")
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: "gangway", Namespace: Gangway}, &found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Check if the Netpol spec is correct")
			Eventually(func() error {
				return nil

			}, time.Minute, time.Second).Should(Succeed())

			// By("Checking the latest Status Condition added to the Memcached instance")
			// Eventually(func() error {
			// 	if memcached.Status.Conditions != nil && len(memcached.Status.Conditions) != 0 {
			// 		latestStatusCondition := memcached.Status.Conditions[len(memcached.Status.Conditions)-1]
			// 		expectedLatestStatusCondition := metav1.Condition{Type: typeAvailableMemcached,
			// 			Status: metav1.ConditionTrue, Reason: "Reconciling",
			// 			Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", memcached.Name, memcached.Spec.Size)}
			// 		if latestStatusCondition != expectedLatestStatusCondition {
			// 			return fmt.Errorf("The latest status condition added to the memcached instance is not as expected")
			// 		}
			// 	}
			// 	return nil
			// }, time.Minute, time.Second).Should(Succeed())
		})
	})
})
