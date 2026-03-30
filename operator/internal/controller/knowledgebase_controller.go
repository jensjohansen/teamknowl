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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	// Reconcile the Deployment.
	if err := r.reconcileDeployment(ctx, knowledgeBase); err != nil {
		log.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}

	// Reconcile the Service.
	if err := r.reconcileService(ctx, knowledgeBase); err != nil {
		log.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}

	// Update the status to reflect that we are ready.
	if knowledgeBase.Status.Conditions[0].Reason != "Deployed" {
		knowledgeBase.Status.Conditions[0].Reason = "Deployed"
		knowledgeBase.Status.Conditions[0].Message = "KnowledgeBase API and Sync services are running"
		if err := r.Status().Update(ctx, knowledgeBase); err != nil {
			log.Error(err, "Failed to update status to Deployed")
			return ctrl.Result{}, err
		}
	}

	log.Info("Successfully reconciled KnowledgeBase", "name", knowledgeBase.Name)
	return ctrl.Result{}, nil
}

func (r *KnowledgeBaseReconciler) reconcileDeployment(ctx context.Context, kb *corev1alpha1.KnowledgeBase) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kb.Name,
			Namespace: kb.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		labels := map[string]string{
			"app":       "teamknowl",
			"instance":  kb.Name,
			"component": "api",
		}
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		deployment.Spec.Template.ObjectMeta.Labels = labels

		// Define the containers.
		// Container 1: The TeamKnowl API.
		apiContainer := corev1.Container{
			Name:  "api",
			Image: "ghcr.io/johnkjohansen/teamknowl-api:latest",
			Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
			Env: []corev1.EnvVar{
				{Name: "DOCS_DIR", Value: "/docs"},
				{Name: "PORT", Value: "8080"},
			},
			VolumeMounts: []corev1.VolumeMount{
				{Name: "docs", MountPath: "/docs", ReadOnly: true},
			},
		}

		// Container 2: Git-Sync Sidecar.
		syncContainer := corev1.Container{
			Name:  "git-sync",
			Image: "registry.k8s.io/git-sync/git-sync:v4.2.3",
			Args: []string{
				fmt.Sprintf("--repo=%s", kb.Spec.Repository.RepositoryURL),
				fmt.Sprintf("--branch=%s", kb.Spec.Repository.BranchName),
				"--root=/docs",
				"--dest=repo",
				"--wait=30", // Sync every 30 seconds for dev/test
			},
			VolumeMounts: []corev1.VolumeMount{
				{Name: "docs", MountPath: "/docs"},
			},
		}

		deployment.Spec.Template.Spec.Containers = []corev1.Container{apiContainer, syncContainer}
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "docs",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}

		return controllerutil.SetControllerReference(kb, deployment, r.Scheme)
	})

	return err
}

func (r *KnowledgeBaseReconciler) reconcileService(ctx context.Context, kb *corev1alpha1.KnowledgeBase) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kb.Name,
			Namespace: kb.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		service.Spec.Selector = map[string]string{
			"app":      "teamknowl",
			"instance": kb.Name,
		}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(8080),
			},
		}
		return controllerutil.SetControllerReference(kb, service, r.Scheme)
	})

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *KnowledgeBaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.KnowledgeBase{}).
		Named("knowledgebase").
		Complete(r)
}
