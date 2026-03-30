/*
File: teamknowl/operator/api/v1alpha1/knowledgebase_types.go
Purpose: Defines the KnowledgeBase Custom Resource Schema for TeamKnowl.
Product/business importance: Enables the management of multiple knowledge base instances within a Kubernetes cluster, allowing for separate DevOps, Support, and Engineering contexts.

Copyright (c) 2026 John K Johansen
License: MIT (see LICENSE)
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RepositoryConfig defines the Git source configuration for the KnowledgeBase.
type RepositoryConfig struct {
	// RepositoryURL is the canonical address of the Git repository containing the Markdown files.
	// +required
	RepositoryURL string `json:"repositoryUrl"`

	// CredentialsSecretReference points to a Kubernetes Secret containing Git authentication tokens or SSH keys.
	// +optional
	CredentialsSecretReference string `json:"credentialsSecretReference,omitempty"`

	// BranchName specifies which Git branch to track for documentation updates.
	// +optional
	// +kubebuilder:default=main
	BranchName string `json:"branchName,omitempty"`
}

// UserInterfaceConfig defines the appearance and accessibility of the Obsidian-like web UI.
type UserInterfaceConfig struct {
	// Enabled determines if the web interface should be deployed for this KnowledgeBase.
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// Theme specifies the visual style (e.g., "dark", "light") for the web frontend.
	// +optional
	// +kubebuilder:default=dark
	Theme string `json:"theme,omitempty"`
}

// APIConfig defines the headless and agentic integration settings.
type APIConfig struct {
	// Enabled determines if the REST API service should be deployed.
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// HeadlessMode disables the browser-facing UI and focuses solely on providing raw context for AI agents.
	// +optional
	// +kubebuilder:default=false
	HeadlessMode bool `json:"headlessMode"`
}

// KnowledgeBaseSpec defines the desired state of a TeamKnowl KnowledgeBase.
type KnowledgeBaseSpec struct {
	// Repository holds the Git-sync configuration for this instance.
	// +required
	Repository RepositoryConfig `json:"repository"`

	// SyncIntervalDuration specifies how often the operator should pull updates from the Git repository (e.g., "5m").
	// +optional
	// +kubebuilder:default="5m"
	SyncIntervalDuration string `json:"syncIntervalDuration,omitempty"`

	// UserInterface holds configuration for the web frontend.
	// +optional
	UserInterface UserInterfaceConfig `json:"userInterface,omitempty"`

	// API holds configuration for the agentic context engine.
	// +optional
	API APIConfig `json:"api,omitempty"`
}

// KnowledgeBaseStatus defines the observed state of a KnowledgeBase.
type KnowledgeBaseStatus struct {
	// LastSynchronizedCommit records the hash of the last successfully indexed Git commit.
	// This helps AI agents and humans verify the freshness of the documentation.
	// +optional
	LastSynchronizedCommit string `json:"lastSynchronizedCommit,omitempty"`

	// Conditions represent the current state of the KnowledgeBase resource (e.g., "Available", "Degraded").
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KnowledgeBase is the Schema for the TeamKnowl knowledgebases API.
type KnowledgeBase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KnowledgeBaseSpec `json:"spec"`
	// +optional
	Status *KnowledgeBaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KnowledgeBaseList contains a list of KnowledgeBase.
type KnowledgeBaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KnowledgeBase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KnowledgeBase{}, &KnowledgeBaseList{})
}
