/*
BSD 2-Clause License

Copyright (c) 2019, Tarantool
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/google/uuid"
	tarantooliov1alpha1 "github.com/tarantool/tarantool-operator/api/v1alpha1"
	"github.com/tarantool/tarantool-operator/controllers/tarantool"
	"github.com/tarantool/tarantool-operator/controllers/topology"
)

var space = uuid.MustParse("73692FF6-EB42-46C2-92B6-65C45191368D")

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Checking for a leader in the cluster Endpoint annotation
func IsLeaderExists(ep *corev1.Endpoints) bool {
	leader, ok := ep.Annotations["tarantool.io/leader"]
	if !ok || leader == "" {
		return false
	}

	for _, addr := range ep.Subsets[0].Addresses {
		if leader == fmt.Sprintf("%s:%s", addr.IP, "8081") {
			return true
		}
	}
	return false
}

// HasInstanceUUID .
func HasInstanceUUID(o *corev1.Pod) bool {
	annotations := o.Labels
	if _, ok := annotations["tarantool.io/instance-uuid"]; ok {
		return true
	}

	return false
}

// SetInstanceUUID .
func SetInstanceUUID(o *corev1.Pod) *corev1.Pod {
	labels := o.Labels
	if len(o.GetName()) == 0 {
		return o
	}
	instanceUUID := uuid.NewSHA1(space, []byte(o.GetName()))
	labels["tarantool.io/instance-uuid"] = instanceUUID.String()

	o.SetLabels(labels)
	return o
}

//+kubebuilder:rbac:groups=tarantool.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tarantool.io,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tarantool.io,resources=clusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=endpoints,verbs=get;create;update;watch;list;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Reconciling Cluster")
	reqLogger.Info("Namespace:" + req.Namespace)

	// do nothing if no Cluster
	cluster := &tarantooliov1alpha1.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, cluster); err != nil {
		return ctrl.Result{RequeueAfter: time.Duration(1 * time.Second)}, err
	}
	clusterSelector, err := metav1.LabelSelectorAsSelector(cluster.Spec.Selector)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
	}

	roleList := &tarantooliov1alpha1.RoleList{}
	if err := r.List(context.TODO(), roleList, &client.ListOptions{LabelSelector: clusterSelector, Namespace: req.NamespacedName.Namespace}); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
		}

		return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
	}
	for _, role := range roleList.Items {
		if metav1.IsControlledBy(&role, cluster) {
			reqLogger.Info("Already owned", "Role.Name", role.Name)
			continue
		}
		annotations := role.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations["tarantool.io/cluster-id"] = cluster.GetName()
		role.SetAnnotations(annotations)
		if err := controllerutil.SetControllerReference(cluster, &role, r.Scheme); err != nil {
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
		}
		if err := r.Update(context.TODO(), &role); err != nil {
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
		}

		reqLogger.Info("Set role ownership", "Role.Name", role.GetName(), "Cluster.Name", cluster.GetName())
	}

	reqLogger.Info("Roles reconciled, moving to pod reconcile")

	// ensure cluster wide Service exists
	svc := &corev1.Service{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: cluster.GetNamespace(), Name: cluster.GetName()}, svc); err != nil {
		if errors.IsNotFound(err) {
			svc.Name = cluster.GetName()
			svc.Namespace = cluster.GetNamespace()
			svc.Spec = corev1.ServiceSpec{
				Selector:  cluster.Spec.Selector.MatchLabels,
				ClusterIP: "None",
				Ports: []corev1.ServicePort{
					{
						Name:     "app",
						Port:     3301,
						Protocol: "TCP",
					},
				},
			}

			if err := controllerutil.SetControllerReference(cluster, svc, r.Scheme); err != nil {
				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
			}

			if err := r.Create(context.TODO(), svc); err != nil {
				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
			}
		}
	}

	// ensure Cluster leader elected
	ep := &corev1.Endpoints{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: cluster.GetNamespace(), Name: cluster.GetName()}, ep); err != nil {
		return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
	}
	if len(ep.Subsets) == 0 || len(ep.Subsets[0].Addresses) == 0 {
		reqLogger.Info("No available Endpoint resource configured for Cluster, waiting")
		return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, nil
	}

	if !IsLeaderExists(ep) {
		leader := fmt.Sprintf("%s:%s", ep.Subsets[0].Addresses[0].IP, "8081")

		if ep.Annotations == nil {
			ep.Annotations = make(map[string]string)
		}

		ep.Annotations["tarantool.io/leader"] = leader
		if err := r.Update(context.TODO(), ep); err != nil {
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
		}
	}

	stsList := &appsv1.StatefulSetList{}
	if err := r.List(context.TODO(), stsList, &client.ListOptions{LabelSelector: clusterSelector, Namespace: req.NamespacedName.Namespace}); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, nil
		}

		return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
	}

	topologyClient := topology.NewBuiltInTopologyService(
		topology.WithTopologyEndpoint(fmt.Sprintf("http://%s/admin/api", ep.Annotations["tarantool.io/leader"])),
		topology.WithClusterID(cluster.GetName()),
	)

	reqLogger.Info("Wait for all instances to start")

	for _, sts := range stsList.Items {
		reqLogger.Info("Checking STS: " + sts.GetName())
		for i := 0; i < int(*sts.Spec.Replicas); i++ {
			pod := &corev1.Pod{}
			name := types.NamespacedName{
				Namespace: req.Namespace,
				Name:      fmt.Sprintf("%s-%d", sts.GetName(), i),
			}
			podName := sts.GetName() + "-" + fmt.Sprint(i)

			if err := r.Get(context.TODO(), name, pod); err != nil {
				if errors.IsNotFound(err) {
					return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
				}

				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
			}

			if pod.Status.Phase != corev1.PodRunning {
				reqLogger.Info("Pod is not yet running: " + podName)
				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, nil
			}

			isPodReady := false
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
					isPodReady = true
					break
				}
			}

			if !isPodReady {
				reqLogger.Info("Pod is not yet running: " + podName)
				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, nil
			}

		}
		if sts.Status.ReadyReplicas < *sts.Spec.Replicas {
			reqLogger.Info("StatefulSet hasn't started, requeue...")
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, nil
		}
	}

	reqLogger.Info("All instances have started")

	var replicasets []topology.Replicaset
	// Create replicasets, using EditTopology method
	for _, sts := range stsList.Items {
		isJoined := false
		replicas := make(map[string]*corev1.Pod)

		// Iterate through pods in replicaset, wait until
		// all pods are ready to be deployed
		for i := 0; i < int(*sts.Spec.Replicas); i++ {
			pod := &corev1.Pod{}
			name := types.NamespacedName{
				Namespace: req.Namespace,
				Name:      fmt.Sprintf("%s-%d", sts.GetName(), i),
			}
			podName := sts.GetName() + "-" + fmt.Sprint(i)
			reqLogger.Info("Start editing topology for replicaset: " + podName)

			if err := r.Get(context.TODO(), name, pod); err != nil {
				if errors.IsNotFound(err) {
					return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
				}

				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
			}

			if tarantool.IsJoined(pod) {
				reqLogger.Info("Already joined", "Pod.Name", pod.Name)
				isJoined = true
				break
			}
			_, exists := replicas[podName]
			if !exists {
				reqLogger.Info("Found new node in replicaset " + podName)
				replicas[podName] = pod
			} else {
				reqLogger.Info("Node has already been explored " + podName)
			}
		}
		if isJoined {
			continue
		}

		if len(replicas) != int(*sts.Spec.Replicas) {
			reqLogger.Info("Hasn't explored all replicas in replicaset - requeue")
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
		}

		reqLogger.Info("All replicas in replicaset explored, creating graphql struct")
		if curReplicaset, err := topologyClient.GetReplicasetValues(replicas, sts.GetName()); err != nil {
			reqLogger.Info("Failed to execute EditTopology")
			return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
		} else {
			replicasets = append(replicasets, curReplicaset)
		}

		reqLogger.Info("Successfully generated graphQL variables for replicaset: " + sts.GetName())
	}

	if len(replicasets) != 0 {
		reqLogger.Info("Sending clusterwide EditTopology request")

		if _, err := topologyClient.ExecEditTopology(replicasets); err != nil {
			reqLogger.Info("Failed to execute EditTopology request")
			return ctrl.Result{RequeueAfter: time.Duration(20 * time.Second)}, err
		}
		reqLogger.Info("Successfull EditTopology request, marking pods as joined")

		for _, sts := range stsList.Items {
			for i := 0; i < int(*sts.Spec.Replicas); i++ {
				pod := &corev1.Pod{}
				name := types.NamespacedName{
					Namespace: req.Namespace,
					Name:      fmt.Sprintf("%s-%d", sts.GetName(), i),
				}

				if err := r.Get(context.TODO(), name, pod); err != nil {
					return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
				}

				tarantool.MarkJoined(pod)
				if err := r.Update(context.TODO(), pod); err != nil {
					return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
				}
			}
		}
	}

	// Bootstrap vshard storages
	for _, sts := range stsList.Items {
		stsAnnotations := sts.GetAnnotations()
		if stsAnnotations["tarantool.io/isBootstrapped"] != "1" {
			reqLogger.Info("cluster is not bootstrapped, bootstrapping", "Statefulset.Name", sts.GetName())
			if err := topologyClient.BootstrapVshard(); err != nil {
				if topology.IsAlreadyBootstrapped(err) {
					stsAnnotations["tarantool.io/isBootstrapped"] = "1"
					sts.SetAnnotations(stsAnnotations)

					if err := r.Update(context.TODO(), &sts); err != nil {
						reqLogger.Error(err, "failed to set bootstrapped annotation")
					}

					reqLogger.Info("Added bootstrapped annotation", "StatefulSet.Name", sts.GetName())

					cluster.Status.State = "Ready"
					err = r.Status().Update(context.TODO(), cluster)
					if err != nil {
						return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
					}
					return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, nil
				}

				reqLogger.Error(err, "Bootstrap vshard error")
				return ctrl.Result{RequeueAfter: time.Duration(60 * time.Second)}, err
			}
		} else {
			reqLogger.Info("cluster is already bootstrapped, not retrying", "Statefulset.Name", sts.GetName())
		}

		if stsAnnotations["tarantool.io/failoverEnabled"] == "1" {
			reqLogger.Info("failover is enabled, not retrying")
		} else {
			enabled, err := topologyClient.GetFailover()
			if err != nil {
				reqLogger.Error(err, "failed to get failover status")
				continue
			}

			if !enabled {
				if err := topologyClient.SetFailover(true); err != nil {
					reqLogger.Error(err, "failed to enable cluster failover")
					continue
				}
				reqLogger.Info("enabled failover")
			}

			stsAnnotations["tarantool.io/failoverEnabled"] = "1"
			sts.SetAnnotations(stsAnnotations)
			if err := r.Update(context.TODO(), &sts); err != nil {
				reqLogger.Error(err, "failed to set failover enabled annotation")
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(10 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tarantooliov1alpha1.Cluster{}).
		Watches(&source.Kind{Type: &tarantooliov1alpha1.Cluster{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
			if a.GetLabels() == nil {
				return []ctrl.Request{}
			}
			return []ctrl.Request{
				{NamespacedName: types.NamespacedName{
					Namespace: a.GetNamespace(),
					Name:      a.GetLabels()["tarantool.io/cluster-id"],
				}},
			}
		})).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 1,
			RateLimiter:             workqueue.NewItemExponentialFailureRateLimiter(1*time.Second, 10*time.Second),
		}).
		Complete(r)
}
