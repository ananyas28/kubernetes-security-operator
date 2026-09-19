/*
Copyright 2026.

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
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	securityv1 "github.com/ananyas28/kubernetes-security-operator/api/v1"
)

var (
	SecurityMonitorReconcileTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "securitymonitor_reconcile_total",
			Help: "Total number of SecurityMonitor reconciliations.",
		},
	)

	SecurityMonitorReconcileErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "securitymonitor_reconcile_errors_total",
			Help: "Total number of SecurityMonitor reconciliation errors.",
		},
	)
)

func init() {
	metrics.Registry.MustRegister(
		SecurityMonitorReconcileTotal,
		SecurityMonitorReconcileErrors,
	)

}

// SecurityMonitorReconciler reconciles a SecurityMonitor object
type SecurityMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=security.example.com,resources=securitymonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.example.com,resources=securitymonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=security.example.com,resources=securitymonitors/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SecurityMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.0/pkg/reconcile
func (r *SecurityMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	SecurityMonitorReconcileTotal.Inc()
	defer func() {
		if err != nil {
			SecurityMonitorReconcileErrors.Inc()
		}
	}()
	logger := logf.FromContext(ctx)

	var securityMonitor securityv1.SecurityMonitor
	if err := r.Get(ctx, req.NamespacedName, &securityMonitor); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("SecurityMonitor found",
		"name", securityMonitor.Name,
		"namespace", securityMonitor.Namespace,
	)

	configMapName := securityMonitor.Name + "-config"

	desiredData := map[string]string{
		"enabled":  strconv.FormatBool(securityMonitor.Spec.Enabled),
		"severity": securityMonitor.Spec.Severity,
		"alertOn":  strings.Join(securityMonitor.Spec.AlertOn, ","),
	}

	configMap := &corev1.ConfigMap{}

	err = r.Get(
		ctx,
		client.ObjectKey{
			Name:      configMapName,
			Namespace: securityMonitor.Namespace,
		},
		configMap,
	)

	if err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		configMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: securityMonitor.Namespace,
			},
			Data: desiredData,
		}

		if err := r.Create(ctx, configMap); err != nil {
			return ctrl.Result{}, err
		}

		logger.Info("SecurityMonitor ConfigMap created",
			"name", configMapName,
		)
	} else {
		if configMap.Data == nil {
			configMap.Data = map[string]string{}
		}

		configMap.Data = desiredData

		if err := r.Update(ctx, configMap); err != nil {
			return ctrl.Result{}, err
		}

		logger.Info("SecurityMonitor ConfigMap updated",
			"name", configMapName,
		)
	}

	securityMonitor.Status.Conditions = []metav1.Condition{
		{
			Type:               "Available",
			Status:             metav1.ConditionTrue,
			Reason:             "SecurityMonitorReady",
			Message:            "SecurityMonitor is being observed by the operator",
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: securityMonitor.Generation,
		},
	}

	if err := r.Status().Update(ctx, &securityMonitor); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecurityMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.SecurityMonitor{}).
		Owns(&corev1.ConfigMap{}).
		Named("securitymonitor").
		Complete(r)
}
