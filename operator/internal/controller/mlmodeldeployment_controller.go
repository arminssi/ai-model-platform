/*
Copyright 2025.

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
	"fmt"
	"os/exec"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	mlv1alpha1 "github.com/arminssi/ai-model-platform/operator/api/v1alpha1"
)

// MLModelDeploymentReconciler reconciles a MLModelDeployment object
type MLModelDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ml.ml.armin.dev,resources=mlmodeldeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ml.ml.armin.dev,resources=mlmodeldeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ml.ml.armin.dev,resources=mlmodeldeployments/finalizers,verbs=update

// Reconcile contains the business logic to deploy a Helm chart when the CR is created/updated
func (r *MLModelDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Fetch the MLModelDeployment instance
	var model mlv1alpha1.MLModelDeployment
	if err := r.Get(ctx, req.NamespacedName, &model); err != nil {
		logger.Error(err, "unable to fetch MLModelDeployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Generate Helm release name based on the resource name
	releaseName := model.Name
	image := model.Spec.Image
	replicas := fmt.Sprintf("%d", model.Spec.Replicas)

	logger.Info("Installing Helm release", "release", releaseName, "image", image, "replicas", replicas)

	// Run helm install/upgrade
	cmd := exec.Command(
		"helm", "upgrade", "--install", releaseName,
		"../helm/model-service", // update if path differs
		"-n", "ml-models",
		"--set", fmt.Sprintf("image.repository=%s", image),
		"--set", fmt.Sprintf("replicaCount=%s", replicas),
		"--create-namespace",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err, "Helm release failed", "output", string(output))
		return ctrl.Result{}, err
	}

	logger.Info("Helm release successful", "output", string(output))

	// Update CR status
	model.Status.Phase = "Deployed"
	if err := r.Status().Update(ctx, &model); err != nil {
		logger.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the controller with the manager
func (r *MLModelDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mlv1alpha1.MLModelDeployment{}).
		Owns(&appsv1.Deployment{}).
		Named("mlmodeldeployment").
		Complete(r)
}
