/*
Copyright 2024.

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

	"github.com/zszabo-rh/issues-operator/gitclient"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	trainingv1alpha1 "github.com/zszabo-rh/issues-operator/api/v1alpha1"
)

// GithubIssueReconciler reconciles a GithubIssue object
type GithubIssueReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=training.redhat.com,resources=githubissues,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=training.redhat.com,resources=githubissues/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=training.redhat.com,resources=githubissues/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GithubIssue object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *GithubIssueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info(">>>>>>>>>>>> Starting reconcile (v6) <<<<<<<<<<<<<<<<<<<<<<<")

	githubissue := &trainingv1alpha1.GithubIssue{}
	err := r.Get(ctx, req.NamespacedName, githubissue)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "Issue not found!")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	clientissue := gitclient.GitIssue{}
	repo := githubissue.Spec.Repository
	clientissue.Title = githubissue.Spec.Title
	clientissue.Description = githubissue.Spec.Description

	issues, err := gitclient.GetIssues(repo)

	if err != nil {
		log.Error(err, "GetIssues("+gitclient.GetUrl(repo)+") failed")
		return ctrl.Result{}, err
	}

	log.Info("Github responded successfully, comparing issues")

	found := false
	for _, issue := range issues {
		if issue.Title == clientissue.Title {
			log.Info("Match! Updating description")
			found = true
			updatedissue, err := gitclient.UpdateIssue(repo, issue.Id, clientissue.Title, clientissue.Description)
			if err != nil {
				log.Error(err, "UpdateIssue("+repo+", "+fmt.Sprintf("%v", clientissue)+") failed")
				return ctrl.Result{}, err
			}

			log.Info("Updating spec")
			err = r.Update(ctx, githubissue)
			if err != nil {
				return ctrl.Result{}, err
			}

			githubissue.Status.State = updatedissue.Status
			githubissue.Status.LastUpdated = updatedissue.LastUpdated

			log.Info("Updating status: " + githubissue.Status.State + ", " + githubissue.Status.LastUpdated)
			err = r.Status().Update(ctx, githubissue)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	if !found {
		log.Info("No issues matched! Creating new github issue")
		newissue, err := gitclient.AddIssue(repo, clientissue.Title, clientissue.Description)
		if err != nil {
			log.Error(err, "AddIssue("+repo+", "+fmt.Sprintf("%v", clientissue)+") failed")
			return ctrl.Result{}, err
		}
		log.Info("Updating spec")

		err = r.Update(ctx, githubissue)
		if err != nil {
			return ctrl.Result{}, err
		}

		githubissue.Status.State = newissue.Status
		githubissue.Status.LastUpdated = newissue.LastUpdated

		log.Info("Updating status: " + githubissue.Status.State + ", " + githubissue.Status.LastUpdated)
		err = r.Status().Update(ctx, githubissue)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GithubIssueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trainingv1alpha1.GithubIssue{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
