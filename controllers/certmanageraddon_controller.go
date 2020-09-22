/*


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
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/status"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/danniel1205/cluster-addons-mgmt/api/v1alpha1"
)

var _ reconcile.Reconciler = &CertManagerAddonReconciler{}

// CertManagerAddonReconciler reconciles a CertManagerAddon object
type CertManagerAddonReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	declarative.Reconciler
}

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=certmanageraddons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=certmanageraddons/status,verbs=get;update;patch

func (r *CertManagerAddonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//loaders.FlagChannel = "https://raw.githubusercontent.com/danniel1205/cluster-addons-channels/master"
	addon.Init()

	labels := map[string]string{
		"k8s-app": "certmanageraddon",
	}

	watchLabels := declarative.SourceLabel(mgr.GetScheme())

	if err := r.Reconciler.Init(mgr, &api.CertManagerAddon{},
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(watchLabels),
		declarative.WithStatus(status.NewKstatusCheck(mgr.GetClient(), &r.Reconciler)),
		declarative.WithObjectTransform(addon.ApplyPatches),
		declarative.WithApplyPrune(),
	); err != nil {
		return err
	}

	c, err := controller.New("certmanageraddon-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to CertManagerAddon
	err = c.Watch(&source.Kind{Type: &api.CertManagerAddon{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to deployed objects
	_, err = declarative.WatchAll(mgr.GetConfig(), c, r, watchLabels)
	if err != nil {
		return err
	}

	return nil
}

// for WithApplyPrune
// +kubebuilder:rbac:groups=*,resources=*,verbs=list
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create;delete
