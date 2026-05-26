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
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
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
// +kubebuilder:rbac:groups="",resources=services;secrets;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses;networkpolicies,verbs=get;list;watch;create;update;patch;delete

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

	// Reconcile the harbor-global-pull secret.
	if err := r.reconcileGlobalPullSecret(ctx, knowledgeBase); err != nil {
		log.Error(err, "Failed to reconcile harbor-global-pull secret")
		return ctrl.Result{}, err
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

	// Reconcile the Ingress.
	if err := r.reconcileIngress(ctx, knowledgeBase); err != nil {
		log.Error(err, "Failed to reconcile Ingress")
		return ctrl.Result{}, err
	}

	// Reconcile the NetworkPolicy.
	if err := r.reconcileNetworkPolicy(ctx, knowledgeBase); err != nil {
		log.Error(err, "Failed to reconcile NetworkPolicy")
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (r *KnowledgeBaseReconciler) reconcileGlobalPullSecret(ctx context.Context, kb *corev1alpha1.KnowledgeBase) error {
	log := log.FromContext(ctx)

	// Define the source and target names.
	sourceNamespace := "ai-models"
	secretName := "harbor-global-pull"

	// Check if the secret already exists in the target namespace.
	targetSecret := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{Name: secretName, Namespace: kb.Namespace}, targetSecret)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	// Fetch the source secret.
	sourceSecret := &corev1.Secret{}
	if err := r.Get(ctx, client.ObjectKey{Name: secretName, Namespace: sourceNamespace}, sourceSecret); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Source harbor-global-pull secret not found in ai-models namespace, skipping copy")
			return nil
		}
		return err
	}

	if err == nil {
		// Secret exists, check if it needs update.
		// We use reflect.DeepEqual or just check the data.
		sourceData := string(sourceSecret.Data[".dockerconfigjson"])
		targetData := string(targetSecret.Data[".dockerconfigjson"])

		if sourceData == targetData {
			// Already in sync.
			return nil
		}

		log.Info("Updating harbor-global-pull secret to match ai-models source", "namespace", kb.Namespace)
		targetSecret.Data = sourceSecret.Data
		targetSecret.Type = sourceSecret.Type
		return r.Update(ctx, targetSecret)
	}

	// Create the new secret in the target namespace.
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: kb.Namespace,
		},
		Type: sourceSecret.Type,
		Data: sourceSecret.Data,
	}

	log.Info("Copying harbor-global-pull secret from ai-models to target namespace", "namespace", kb.Namespace)
	if err := r.Create(ctx, newSecret); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}

func (r *KnowledgeBaseReconciler) reconcileDeployment(ctx context.Context, kb *corev1alpha1.KnowledgeBase) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kb.Name,
			Namespace: kb.Namespace,
		},
	}

	apiImage := getEnv("API_IMAGE", "harbor.ai-agents.private/teamknowl/api:latest")
	uiImage := getEnv("UI_IMAGE", "harbor.ai-agents.private/teamknowl/ui:latest")
	syncImage := getEnv("GIT_SYNC_IMAGE", "registry.k8s.io/git-sync/git-sync:v4.2.3")

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		labels := map[string]string{
			"app":       "teamknowl",
			"instance":  kb.Name,
		}
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		deployment.Spec.Template.ObjectMeta.Labels = labels

		// Define the containers.
		// Container 1: The TeamKnowl API.
		apiContainer := corev1.Container{
			Name:  "api",
			Image: apiImage,
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
			Image: syncImage,
			Args: []string{
				fmt.Sprintf("--repo=%s", kb.Spec.Repository.RepositoryURL),
				fmt.Sprintf("--ref=%s", kb.Spec.Repository.BranchName),
				"--root=/docs",
				"--link=repo",
				"--period=30s",
			},
			VolumeMounts: []corev1.VolumeMount{
				{Name: "docs", MountPath: "/docs"},
			},
		}

		// Handle Git-Sync Authentication.
		if kb.Spec.Repository.CredentialsSecretReference != "" {
			syncContainer.Env = append(syncContainer.Env,
				corev1.EnvVar{
					Name: "GIT_SYNC_USERNAME",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: kb.Spec.Repository.CredentialsSecretReference,
							},
							Key: "username",
						},
					},
				},
				corev1.EnvVar{
					Name: "GIT_SYNC_PASSWORD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: kb.Spec.Repository.CredentialsSecretReference,
							},
							Key: "password",
						},
					},
				},
			)
		}

		// Handle API/UI Storage configuration.
		// If using S3, we need to pass credentials to the API.
		if kb.Spec.Storage.Provider == "s3" && kb.Spec.Storage.S3Config != nil {
			apiContainer.Env = append(apiContainer.Env,
				corev1.EnvVar{Name: "STORAGE_PROVIDER", Value: "s3"},
				corev1.EnvVar{Name: "S3_ENDPOINT", Value: kb.Spec.Storage.S3Config.Endpoint},
				corev1.EnvVar{Name: "S3_BUCKET", Value: kb.Spec.Storage.S3Config.BucketName},
			)

			if kb.Spec.Storage.S3Config.CredentialsSecretName != "" {
				apiContainer.Env = append(apiContainer.Env,
					corev1.EnvVar{
						Name: "S3_ACCESS_KEY",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: kb.Spec.Storage.S3Config.CredentialsSecretName,
								},
								Key: "accessKey",
							},
						},
					},
					corev1.EnvVar{
						Name: "S3_SECRET_KEY",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: kb.Spec.Storage.S3Config.CredentialsSecretName,
								},
								Key: "secretKey",
							},
						},
					},
				)
			}
		}

		containers := []corev1.Container{apiContainer, syncContainer}

		// Container 3: The UI (if enabled).
		if kb.Spec.UserInterface.Enabled {
			uiContainer := corev1.Container{
				Name:  "ui",
				Image: uiImage,
				Ports: []corev1.ContainerPort{{ContainerPort: 3000}},
				Env: []corev1.EnvVar{
					{Name: "NEXT_PUBLIC_API_URL", Value: "http://localhost:8080"},
				},
			}
			containers = append(containers, uiContainer)
		}

		deployment.Spec.Template.Spec.Containers = containers
		deployment.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: "harbor-global-pull"},
		}
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
				Name:       "ui",
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(3000),
			},
			{
				Name:       "api",
				Protocol:   corev1.ProtocolTCP,
				Port:       8080,
				TargetPort: intstr.FromInt(8080),
			},
		}
		return controllerutil.SetControllerReference(kb, service, r.Scheme)
	})

	return err
}

func (r *KnowledgeBaseReconciler) reconcileIngress(ctx context.Context, kb *corev1alpha1.KnowledgeBase) error {
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kb.Name,
			Namespace: kb.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		pathType := networkingv1.PathTypePrefix
		host := fmt.Sprintf("%s.ai-agents.private", kb.Name)

		ingress.Spec.Rules = []networkingv1.IngressRule{
			{
				Host: host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: kb.Name,
										Port: networkingv1.ServiceBackendPort{
											Name: "ui",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		// Optional: Add TLS if cert-manager is in use.
		// For now, we assume simple HTTP for the internal LAN.
		return controllerutil.SetControllerReference(kb, ingress, r.Scheme)
	})

	return err
}

func (r *KnowledgeBaseReconciler) reconcileNetworkPolicy(ctx context.Context, kb *corev1alpha1.KnowledgeBase) error {
	networkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kb.Name,
			Namespace: kb.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, networkPolicy, func() error {
		labels := map[string]string{
			"app":      "teamknowl",
			"instance": kb.Name,
		}
		networkPolicy.Spec.PodSelector = metav1.LabelSelector{
			MatchLabels: labels,
		}
		networkPolicy.Spec.PolicyTypes = []networkingv1.PolicyType{
			networkingv1.PolicyTypeIngress,
			networkingv1.PolicyTypeEgress,
		}

		// Allow ingress from the same namespace (UI to API) and from Ingress Controller.
		networkPolicy.Spec.Ingress = []networkingv1.NetworkPolicyIngressRule{
			{
				From: []networkingv1.NetworkPolicyPeer{
					{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
					},
					{
						NamespaceSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"kubernetes.io/metadata.name": "ingress-nginx",
							},
						},
					},
				},
				Ports: []networkingv1.NetworkPolicyPort{
					{
						Protocol: ptrProto(corev1.ProtocolTCP),
						Port:     ptrIntStr(8080),
					},
					{
						Protocol: ptrProto(corev1.ProtocolTCP),
						Port:     ptrIntStr(3000),
					},
				},
			},
		}

		// Allow egress to DNS and potentially S3/Git.
		networkPolicy.Spec.Egress = []networkingv1.NetworkPolicyEgressRule{
			{
				To: []networkingv1.NetworkPolicyPeer{
					{
						NamespaceSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"kubernetes.io/metadata.name": "kube-system",
							},
						},
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"k8s-app": "kube-dns",
							},
						},
					},
				},
				Ports: []networkingv1.NetworkPolicyPort{
					{
						Protocol: ptrProto(corev1.ProtocolUDP),
						Port:     ptrIntStr(53),
					},
				},
			},
			{
				To: []networkingv1.NetworkPolicyPeer{
					{
						IPBlock: &networkingv1.IPBlock{
							CIDR: "0.0.0.0/0",
						},
					},
				},
			},
		}

		return controllerutil.SetControllerReference(kb, networkPolicy, r.Scheme)
	})

	return err
}

func ptrProto(p corev1.Protocol) *corev1.Protocol {
	return &p
}

func ptrIntStr(p int) *intstr.IntOrString {
	i := intstr.FromInt(p)
	return &i
}

// SetupWithManager sets up the controller with the Manager.
func (r *KnowledgeBaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.KnowledgeBase{}).
		Named("knowledgebase").
		Complete(r)
}
