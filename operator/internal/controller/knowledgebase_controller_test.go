/*
File: teamknowl/operator/internal/controller/knowledgebase_controller_test.go
Purpose: Unit tests for the KnowledgeBase controller reconciliation loop.
Product/business importance: Ensures the stability and correctness of the KnowledgeBase lifecycle management.

Copyright (c) 2026 John K Johansen
License: MIT (see LICENSE)
*/

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1alpha1 "github.com/johnkjohansen/teamknowl/api/v1alpha1"
)

var _ = Describe("KnowledgeBase Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-kb"
		const resourceNamespace = "default"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}
		knowledgebase := &corev1alpha1.KnowledgeBase{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind KnowledgeBase")
			err := k8sClient.Get(ctx, typeNamespacedName, knowledgebase)
			if err != nil && errors.IsNotFound(err) {
				resource := &corev1alpha1.KnowledgeBase{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: resourceNamespace,
					},
					Spec: corev1alpha1.KnowledgeBaseSpec{
						Repository: corev1alpha1.RepositoryConfig{
							RepositoryURL: "https://github.com/johnkjohansen/teamknowl-test-docs.git",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			By("Cleanup the specific resource instance KnowledgeBase")
			resource := &corev1alpha1.KnowledgeBase{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should successfully reconcile the resource and initialize conditions", func() {
			By("Reconciling the created resource")
			controllerReconciler := &KnowledgeBaseReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying that the status conditions were initialized")
			updatedKB := &corev1alpha1.KnowledgeBase{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, updatedKB)).To(Succeed())
			Expect(updatedKB.Status.Conditions).To(HaveLen(1))
			Expect(updatedKB.Status.Conditions[0].Type).To(Equal("Progressing"))
			Expect(updatedKB.Status.Conditions[0].Status).To(Equal(metav1.ConditionTrue))
		})
	})
})
