/*
File: teamknowl/operator/internal/controller/knowledgebase_controller.go
Purpose: Implements the reconciliation logic for the KnowledgeBase Custom Resource.
Product/business importance: Ensures that every KnowledgeBase requested by a user results in a functioning, synchronized documentation instance within the cluster.

Copyright (c) 2026 John K Johansen
License: MIT (see LICENSE)
*/

package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1alpha1 "github.com/johnkjohansen/teamknowl/api/v1alpha1"
)

// KnowledgeBaseReconciler reconciles a KnowledgeBase object.
type KnowledgeBaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.teamknowl.io,resources=knowledgebases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.teamknowl.io,resources=knowledgebases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.teamknowl.io,resources=knowledgebases/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile coordinates the cluster state with the desired KnowledgeBase specification.
func (r *KnowledgeBaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the KnowledgeBase instance from the cluster.
	// We use this to understand what the user wants to achieve.
	knowledgeBase := &corev1alpha1.KnowledgeBase{}
	if err := r.Get(ctx, req.NamespacedName, knowledgeBase); err != nil {
		if apierrors.IsNotFound(err) {
			// Resource not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to fetch KnowledgeBase resource")
		return ctrl.Result{}, err
	}

	// Initialize status if it's missing.
	if knowledgeBase.Status == nil {
		knowledgeBase.Status = &corev1alpha1.KnowledgeBaseStatus{}
	}

	// Initialize status conditions if they are missing.
	// This provides immediate feedback to the user that the operator has acknowledged the request.
	if len(knowledgeBase.Status.Conditions) == 0 {
		knowledgeBase.Status.Conditions = []metav1.Condition{
			{
				Type:               "Progressing",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "Initialization",
				Message:            "Operator is beginning to reconcile the KnowledgeBase",
			},
		}
		if err := r.Status().Update(ctx, knowledgeBase); err != nil {
			log.Error(err, "Failed to initialize KnowledgeBase status")
			return ctrl.Result{}, err
		}
	}

	// TODO: Implement deployment of the TeamKnowl API and Git-sync sidecar.
	// For now, we simply acknowledge the resource exists.
	log.Info("Successfully reconciled KnowledgeBase", "name", knowledgeBase.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KnowledgeBaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.KnowledgeBase{}).
		Named("knowledgebase").
		Complete(r)
}
