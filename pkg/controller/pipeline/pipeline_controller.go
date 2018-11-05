/*
Copyright 2018 Aljabr Inc.

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

package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	agentsv1alpha1 "github.com/AljabrIO/koalja-operator/pkg/apis/agents/v1alpha1"
	koaljav1alpha1 "github.com/AljabrIO/koalja-operator/pkg/apis/koalja/v1alpha1"
	"github.com/AljabrIO/koalja-operator/pkg/constants"
	"github.com/AljabrIO/koalja-operator/pkg/util"
	"github.com/rs/zerolog"
)

// Add creates a new Pipeline Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePipeline{
		log:           util.MustCreateLogger(),
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		eventRecorder: mgr.GetRecorder("controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
// +kubebuilder:rbac:groups=v1,resources=serviceaccounts,verbs=list;watch
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("pipeline-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Pipeline
	err = c.Watch(&source.Kind{Type: &koaljav1alpha1.Pipeline{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch Deployments since we're launching a various Agent Deployments
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &koaljav1alpha1.Pipeline{},
	})
	if err != nil {
		return err
	}

	// Watch StatefulSets since we're launching a various Agent StatefulSets
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &koaljav1alpha1.Pipeline{},
	})
	if err != nil {
		return err
	}

	// Watch Services since we're launching a various Agent Services
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &koaljav1alpha1.Pipeline{},
	})
	if err != nil {
		return err
	}

	// Watch ServiceAccounts since we're launching a various ServiceAccounts as identify for the pipeline
	// components.
	err = c.Watch(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &koaljav1alpha1.Pipeline{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcilePipeline{}

// ReconcilePipeline reconciles a Pipeline object
type ReconcilePipeline struct {
	log zerolog.Logger
	client.Client
	scheme        *runtime.Scheme
	eventRecorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a Pipeline object and makes changes based on the state read
// and what is in the Pipeline.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=v1,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=koalja.aljabr.io,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcilePipeline) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	log := r.log.With().
		Str("name", request.Name).
		Str("namespace", request.Namespace).
		Logger()
	log.Debug().Msg("reconcile")

	// Fetch the Pipeline instance
	instance := &koaljav1alpha1.Pipeline{}
	if err := r.Get(ctx, request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Validate the spec
	if err := instance.Spec.Validate(); err != nil {
		log.Warn().Err(err).Msg("Pipeline is invalid")
		r.eventRecorder.Eventf(instance, "Warning", "PipelineValidation", "Pipeline is not valid: %s", err)
		return reconcile.Result{}, nil
	} else {
		r.eventRecorder.Event(instance, "Normal", "PipelineValidation", "Pipeline is valid")
	}

	// Ensure all resources are created
	result := reconcile.Result{
		Requeue:      true,
		RequeueAfter: time.Minute,
	}

	// Ensure agents ServiceAccount is created
	if lresult, err := r.ensureAgentsServiceAccount(ctx, instance); err != nil {
		log.Error().Err(err).Msg("ensureAgentsServiceAccount failed")
		return lresult, err
	} else {
		result = MergeReconcileResult(result, lresult)
	}

	// Ensure executors ServiceAccount is created
	if lresult, err := r.ensureExecutorsServiceAccount(ctx, instance); err != nil {
		log.Error().Err(err).Msg("ensureExecutorsServiceAccount failed")
		return lresult, err
	} else {
		result = MergeReconcileResult(result, lresult)
	}

	// Ensure agents Role & RoleBinding are created
	if lresult, err := r.ensureAgentsRoleAndBinding(ctx, instance); err != nil {
		log.Error().Err(err).Msg("ensureAgentsRoleAndBinding failed")
		return lresult, err
	} else {
		result = MergeReconcileResult(result, lresult)
	}

	// Ensure executors Role & RoleBinding are created
	if lresult, err := r.ensureExecutorsRoleAndBinding(ctx, instance); err != nil {
		log.Error().Err(err).Msg("ensureExecutorsRoleAndBinding failed")
		return lresult, err
	} else {
		result = MergeReconcileResult(result, lresult)
	}

	// Ensure pipeline agent is created
	if lresult, err := r.ensurePipelineAgent(ctx, instance); err != nil {
		log.Error().Err(err).Msg("ensurePipelineAgent failed")
		return lresult, err
	} else {
		result = MergeReconcileResult(result, lresult)
	}

	// Ensure link agents are created
	for _, l := range instance.Spec.Links {
		if lresult, err := r.ensureLinkAgent(ctx, instance, l); err != nil {
			log.Error().Err(err).Msg("ensureLinkAgent failed")
			return lresult, err
		} else {
			result = MergeReconcileResult(result, lresult)
		}
	}

	// Ensure task agents are created
	for _, t := range instance.Spec.Tasks {
		if lresult, err := r.ensureTaskAgent(ctx, instance, t); err != nil {
			log.Error().Err(err).Msg("ensureTaskAgent failed")
			return lresult, err
		} else {
			result = MergeReconcileResult(result, lresult)
		}
	}

	return result, nil
}

// ensureAgentsServiceAccount ensures that a service account exists that is the identity
// for all pipeline, link & task agents.
// +kubebuilder:rbac:groups=v1,resources=serviceaccounts,verbs=get;create;update;patch;delete
func (r *ReconcilePipeline) ensureAgentsServiceAccount(ctx context.Context, instance *koaljav1alpha1.Pipeline) (reconcile.Result, error) {
	serviceAccountName := CreatePipelineAgentsServiceAccountName(instance.Name)
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: instance.Namespace,
		},
	}

	log := r.log.With().
		Str("name", serviceAccount.Name).
		Str("namespace", serviceAccount.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, serviceAccount, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agents ServiceAccount already exists
	if err := util.EnsureServiceAccount(ctx, log, r.Client, serviceAccount, "Pipeline Agents ServiceAccount"); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// ensureExecutorsServiceAccount ensures that a service account exists that is the identity
// for all task executors.
// +kubebuilder:rbac:groups=v1,resources=serviceaccounts,verbs=get;create;update;patch;delete
func (r *ReconcilePipeline) ensureExecutorsServiceAccount(ctx context.Context, instance *koaljav1alpha1.Pipeline) (reconcile.Result, error) {
	serviceAccountName := CreatePipelineExecutorsServiceAccountName(instance.Name)
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: instance.Namespace,
		},
	}

	log := r.log.With().
		Str("name", serviceAccount.Name).
		Str("namespace", serviceAccount.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, serviceAccount, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agents ServiceAccount already exists
	if err := util.EnsureServiceAccount(ctx, log, r.Client, serviceAccount, "Pipeline Executors ServiceAccount"); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// ensureAgentsRoleAndBinding ensures that a Role & RoleBinding exists for the service account that is the identity
// for all pipeline, link & task agents.
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;create;update;patch;delete
func (r *ReconcilePipeline) ensureAgentsRoleAndBinding(ctx context.Context, instance *koaljav1alpha1.Pipeline) (reconcile.Result, error) {
	// Role
	roleName := CreatePipelineAgentsRoleName(instance.Name)
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: instance.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"pods", "services"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
			},
		},
	}
	log := r.log.With().
		Str("name", role.Name).
		Str("namespace", role.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agents Role already exists
	if err := util.EnsureRole(ctx, log, r.Client, role, "Pipeline Agents Role"); err != nil {
		return reconcile.Result{}, err
	}

	// RoleBinding
	roleBindingName := CreatePipelineAgentsRoleBindingName(instance.Name)
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: instance.Namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      CreatePipelineAgentsServiceAccountName(instance.Name),
				Namespace: instance.Namespace,
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agents RoleBinding already exists
	if err := util.EnsureRoleBinding(ctx, log, r.Client, roleBinding, "Pipeline Agents RoleBinding"); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// ensureExecutorsRoleAndBinding ensures that a Role & RoleBinding exists for the service account that is the identity
// for all task executors.
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;create;update;patch;delete
func (r *ReconcilePipeline) ensureExecutorsRoleAndBinding(ctx context.Context, instance *koaljav1alpha1.Pipeline) (reconcile.Result, error) {
	// Role
	roleName := CreatePipelineExecutorsRoleName(instance.Name)
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: instance.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			// Allow reading pods
			rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			// Allow creating services
			rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
			},
		},
	}
	log := r.log.With().
		Str("name", role.Name).
		Str("namespace", role.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agents Role already exists
	if err := util.EnsureRole(ctx, log, r.Client, role, "Pipeline Executors Role"); err != nil {
		return reconcile.Result{}, err
	}

	// RoleBinding
	roleBindingName := CreatePipelineExecutorsRoleBindingName(instance.Name)
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: instance.Namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      CreatePipelineExecutorsServiceAccountName(instance.Name),
				Namespace: instance.Namespace,
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agents RoleBinding already exists
	if err := util.EnsureRoleBinding(ctx, log, r.Client, roleBinding, "Pipeline Executors RoleBinding"); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// ensurePipelineAgent ensures that a pipeline agent is launched for the given pipeline instance.
// +kubebuilder:rbac:groups=agents.aljabr.io,resources=pipelines,verbs=get;list;watch
// +kubebuilder:rbac:groups=agents.aljabr.io,resources=eventregistries,verbs=get;list;watch
func (r *ReconcilePipeline) ensurePipelineAgent(ctx context.Context, instance *koaljav1alpha1.Pipeline) (reconcile.Result, error) {
	// Search for pipeline agent resource
	var plAgentList agentsv1alpha1.PipelineList
	if err := r.List(ctx, &client.ListOptions{}, &plAgentList); err != nil {
		return reconcile.Result{}, err
	}
	if len(plAgentList.Items) == 0 {
		// No pipeline agent resource found
		r.log.Warn().Msg("No Pipeline Agents found")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 10,
		}, nil
	}
	agentCont := *plAgentList.Items[0].Spec.Container
	SetAgentContainerDefaults(&agentCont, true)
	SetContainerEnvVars(&agentCont, map[string]string{
		constants.EnvAPIPort:               strconv.Itoa(constants.AgentAPIPort),
		constants.EnvAPIHTTPPort:           strconv.Itoa(constants.AgentAPIHTTPPort),
		constants.EnvPipelineName:          instance.Name,
		constants.EnvAgentRegistryAddress:  CreateAgentRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvStatisticsSinkAddress: CreateAgentRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvEventRegistryAddress:  net.JoinHostPort("localhost", strconv.Itoa(constants.EventRegistryAPIPort)),
		constants.EnvDNSName:               CreatePipelineAgentDNSName(instance.Name, instance.Namespace),
	})

	// Search for event registry resource
	var evtRegistryList agentsv1alpha1.EventRegistryList
	if err := r.List(ctx, &client.ListOptions{}, &evtRegistryList); err != nil {
		return reconcile.Result{}, err
	}
	if len(evtRegistryList.Items) == 0 {
		// No event registry resource found
		r.log.Warn().Msg("No Event Registries found")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 10,
		}, nil
	}
	evtRegistryCont := *evtRegistryList.Items[0].Spec.Container
	SetEventRegistryContainerDefaults(&evtRegistryCont)
	SetContainerEnvVars(&evtRegistryCont, map[string]string{
		constants.EnvAPIPort:      strconv.Itoa(constants.EventRegistryAPIPort),
		constants.EnvPipelineName: instance.Name,
	})

	// Define the desired Deployment object for pipeline agent
	deplName := CreatePipelineAgentName(instance.Name)
	createDeplLabels := func() map[string]string {
		return map[string]string{
			"statefulset": deplName,
		}
	}
	deploy := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deplName,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: createDeplLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: createDeplLabels()},
				Spec: corev1.PodSpec{
					Containers:         []corev1.Container{agentCont, evtRegistryCont},
					ServiceAccountName: CreatePipelineAgentsServiceAccountName(instance.Name),
				},
			},
		},
	}
	log := r.log.With().
		Str("name", deploy.Name).
		Str("namespace", deploy.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pipeline agent StatefulSet already exists
	if err := util.EnsureStatefulSet(ctx, log, r.Client, deploy, "Pipeline Agent StatefulSet"); err != nil {
		return reconcile.Result{}, err
	}

	{
		// Define the desired Service object for pipeline agent
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deplName,
				Namespace: instance.Namespace,
				Labels:    createDeplLabels(),
			},
			Spec: corev1.ServiceSpec{
				Selector: createDeplLabels(),
				Type:     corev1.ServiceTypeClusterIP,
				Ports: []corev1.ServicePort{
					corev1.ServicePort{
						Name:       "grpc-agent-api",
						Port:       constants.AgentAPIPort,
						TargetPort: intstr.FromInt(constants.AgentAPIPort),
						Protocol:   corev1.ProtocolTCP,
					},
					corev1.ServicePort{
						Name:       "http-agent-api",
						Port:       constants.AgentAPIHTTPPort,
						TargetPort: intstr.FromInt(constants.AgentAPIHTTPPort),
						Protocol:   corev1.ProtocolTCP,
					},
					corev1.ServicePort{
						Name:       "grpc-event-registry-api",
						Port:       constants.EventRegistryAPIPort,
						TargetPort: intstr.FromInt(constants.EventRegistryAPIPort),
						Protocol:   corev1.ProtocolTCP,
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if the pipeline agent Service already exists
		if err := util.EnsureService(ctx, log, r.Client, service, "Pipeline Agent Service"); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// ensureLinkAgent ensures that a link agent is launched for the given link in given pipeline instance.
// +kubebuilder:rbac:groups=agents.aljabr.io,resources=links,verbs=get;list;watch
func (r *ReconcilePipeline) ensureLinkAgent(ctx context.Context, instance *koaljav1alpha1.Pipeline, link koaljav1alpha1.LinkSpec) (reconcile.Result, error) {
	// Search for link agent resource
	var linkAgentList agentsv1alpha1.LinkList
	if err := r.List(ctx, &client.ListOptions{}, &linkAgentList); err != nil {
		return reconcile.Result{}, err
	}
	if len(linkAgentList.Items) == 0 {
		// No link agent resource found
		r.log.Warn().Msg("No Link Agents found")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 10,
		}, nil
	}
	c := *linkAgentList.Items[0].Spec.Container
	SetAgentContainerDefaults(&c, false)
	SetContainerEnvVars(&c, map[string]string{
		constants.EnvAPIPort:               strconv.Itoa(constants.AgentAPIPort),
		constants.EnvPipelineName:          instance.Name,
		constants.EnvLinkName:              link.Name,
		constants.EnvAgentRegistryAddress:  CreateAgentRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvStatisticsSinkAddress: CreateAgentRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvEventRegistryAddress:  CreateEventRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvDNSName:               CreateLinkAgentDNSName(instance.Name, link.Name, instance.Namespace),
	})

	// Define the desired StatefulSet object for link agent
	deplName := CreateLinkAgentName(instance.Name, link.Name)
	createDeplLabels := func() map[string]string {
		return map[string]string{
			"statefulset": deplName,
			"link":        link.Name,
		}
	}
	deploy := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deplName,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: createDeplLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: createDeplLabels()},
				Spec: corev1.PodSpec{
					Containers:         []corev1.Container{c},
					ServiceAccountName: CreatePipelineAgentsServiceAccountName(instance.Name),
				},
			},
		},
	}
	log := r.log.With().
		Str("name", deploy.Name).
		Str("namespace", deploy.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the link agent StatefulSet already exists
	if err := util.EnsureStatefulSet(ctx, log, r.Client, deploy, "Link Agent StatefulSet"); err != nil {
		return reconcile.Result{}, err
	}

	{
		// Define the desired Service object for link agent
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deplName,
				Namespace: instance.Namespace,
				Labels:    createDeplLabels(),
			},
			Spec: corev1.ServiceSpec{
				Selector: createDeplLabels(),
				Type:     corev1.ServiceTypeClusterIP,
				Ports: []corev1.ServicePort{
					corev1.ServicePort{
						Name:       "grpc-api",
						Port:       constants.AgentAPIPort,
						TargetPort: intstr.FromInt(constants.AgentAPIPort),
						Protocol:   corev1.ProtocolTCP,
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if the link agent Service already exists
		if err := util.EnsureService(ctx, log, r.Client, service, "Link Agent Service"); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// ensureTaskAgent ensures that a task agent is launched for the given task in given pipeline instance.
// +kubebuilder:rbac:groups=agents.aljabr.io,resources=tasks,verbs=get;list;watch
// +kubebuilder:rbac:groups=agents.aljabr.io,resources=taskexecutors,verbs=get;list;watch
func (r *ReconcilePipeline) ensureTaskAgent(ctx context.Context, instance *koaljav1alpha1.Pipeline, task koaljav1alpha1.TaskSpec) (reconcile.Result, error) {
	// Search for FileSystem service
	var svcList corev1.ServiceList
	labelSel, err := labels.Parse(constants.LabelServiceType + "=" + constants.ServiceTypeFilesystem)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := r.List(ctx, &client.ListOptions{
		LabelSelector: labelSel,
	}, &svcList); err != nil {
		return reconcile.Result{}, err
	}
	if len(svcList.Items) == 0 {
		// No task agent resource found
		r.log.Warn().Msg("No FileSystem service found")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 10,
		}, nil
	}
	filesystemServiceAddress := CreateServiceAddress(svcList.Items[0])

	// Search for matching taskexecutor (if needed)
	var annTaskExecutorContainer string
	if task.Type != "" {
		var taskExecList agentsv1alpha1.TaskExecutorList
		if err := r.List(ctx, &client.ListOptions{Namespace: instance.Namespace}, &taskExecList); err != nil {
			return reconcile.Result{}, err
		}
		found := false
		for _, entry := range taskExecList.Items {
			if string(entry.Spec.Type) == string(task.Type) {
				// Found
				found = true
				encoded, err := json.Marshal(entry.Spec.Container)
				if err != nil {
					return reconcile.Result{}, err
				}
				annTaskExecutorContainer = string(encoded)
				break
			}
		}
		if !found {
			r.eventRecorder.Eventf(instance, "Warning", "PipelineValidation", "No TaskExecutor of type '%s' found for task '%s'", task.Type, task.Name)
			r.log.Warn().Msgf("No TaskExecutor of type '%s' found", task.Type)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second * 10,
			}, nil
		}
	}

	// Search for task agent resource
	var taskAgentList agentsv1alpha1.TaskList
	if err := r.List(ctx, &client.ListOptions{}, &taskAgentList); err != nil {
		return reconcile.Result{}, err
	}
	if len(taskAgentList.Items) == 0 {
		// No task agent resource found
		r.log.Warn().Msg("No Task Agents found")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 10,
		}, nil
	}
	c := *taskAgentList.Items[0].Spec.Container
	SetAgentContainerDefaults(&c, false)
	SetContainerEnvVars(&c, map[string]string{
		constants.EnvAPIPort:               strconv.Itoa(constants.AgentAPIPort),
		constants.EnvPipelineName:          instance.Name,
		constants.EnvTaskName:              task.Name,
		constants.EnvAgentRegistryAddress:  CreateAgentRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvStatisticsSinkAddress: CreateAgentRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvEventRegistryAddress:  CreateEventRegistryAddress(instance.Name, instance.Namespace),
		constants.EnvFileSystemAddress:     filesystemServiceAddress,
		constants.EnvDNSName:               CreateTaskAgentDNSName(instance.Name, task.Name, instance.Namespace),
		constants.EnvServiceAccountName:    CreatePipelineExecutorsServiceAccountName(instance.Name),
	})

	// Define the desired StatefulSet object for task agent
	deplName := CreateTaskAgentName(instance.Name, task.Name)
	createDeplLabels := func() map[string]string {
		return map[string]string{
			"statefulset": deplName,
			"task":        task.Name,
		}
	}
	// Create annotations to pass address of input links
	annotations := make(map[string]string)
	for _, tis := range task.Inputs {
		annKey := constants.CreateInputLinkAddressAnnotationName(tis.Name)
		ref := task.Name + "/" + tis.Name
		link, found := instance.Spec.LinkByDestinationRef(ref)
		if !found {
			r.log.Error().Str("ref", ref).Msg("No link found for DestinationRef")
			return reconcile.Result{}, fmt.Errorf("No link found with DestinationRef '%s'", ref)
		}
		annotations[annKey] = CreateLinkAgentEventSourceAddress(instance.Name, link.Name, instance.Namespace)
	}
	// Create annotations to pass addresses of output links
	for _, tos := range task.Outputs {
		annKey := constants.CreateOutputLinkAddressesAnnotationName(tos.Name)
		ref := task.Name + "/" + tos.Name
		links := instance.Spec.LinksBySourceRef(ref)
		if len(links) > 0 {
			addresses := make([]string, 0, len(links))
			for _, link := range links {
				addresses = append(addresses, CreateLinkAgentEventPublisherAddress(instance.Name, link.Name, instance.Namespace))
			}
			annotations[annKey] = strings.Join(addresses, ",")
		} else {
			// Task output is not connected. Connect it to the pipeline agent
			annotations[annKey] = CreatePipelineAgentEventPublisherAddress(instance.Name, instance.Namespace)
		}
	}
	// Create annotation containing task executor container (if any)
	if annTaskExecutorContainer != "" {
		annotations[constants.AnnTaskExecutorContainer] = annTaskExecutorContainer
	}
	deploy := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deplName,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: createDeplLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      createDeplLabels(),
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					Containers:         []corev1.Container{c},
					ServiceAccountName: CreatePipelineAgentsServiceAccountName(instance.Name),
				},
			},
		},
	}
	log := r.log.With().
		Str("name", deploy.Name).
		Str("namespace", deploy.Namespace).
		Logger()
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the task agent StatefulSet already exists
	if err := util.EnsureStatefulSet(ctx, log, r.Client, deploy, "Task Agent StatefulSet"); err != nil {
		return reconcile.Result{}, err
	}

	{
		// Define the desired Service object for task agent
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deplName,
				Namespace: instance.Namespace,
				Labels:    createDeplLabels(),
			},
			Spec: corev1.ServiceSpec{
				Selector: createDeplLabels(),
				Type:     corev1.ServiceTypeClusterIP,
				Ports: []corev1.ServicePort{
					corev1.ServicePort{
						Name:       "grpc-api",
						Port:       constants.AgentAPIPort,
						TargetPort: intstr.FromInt(constants.AgentAPIPort),
						Protocol:   corev1.ProtocolTCP,
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if the task agent Service already exists
		if err := util.EnsureService(ctx, log, r.Client, service, "Task Agent Service"); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
