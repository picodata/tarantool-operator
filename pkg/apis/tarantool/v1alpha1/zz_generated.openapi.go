// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.Cluster":                  schema_pkg_apis_tarantool_v1alpha1_Cluster(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ClusterSpec":              schema_pkg_apis_tarantool_v1alpha1_ClusterSpec(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ClusterStatus":            schema_pkg_apis_tarantool_v1alpha1_ClusterStatus(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplate":       schema_pkg_apis_tarantool_v1alpha1_ReplicasetTemplate(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplateSpec":   schema_pkg_apis_tarantool_v1alpha1_ReplicasetTemplateSpec(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplateStatus": schema_pkg_apis_tarantool_v1alpha1_ReplicasetTemplateStatus(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.Role":                     schema_pkg_apis_tarantool_v1alpha1_Role(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.RoleSpec":                 schema_pkg_apis_tarantool_v1alpha1_RoleSpec(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.RoleStatus":               schema_pkg_apis_tarantool_v1alpha1_RoleStatus(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolCluster":         schema_pkg_apis_tarantool_v1alpha1_TarantoolCluster(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolClusterSpec":     schema_pkg_apis_tarantool_v1alpha1_TarantoolClusterSpec(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolClusterStatus":   schema_pkg_apis_tarantool_v1alpha1_TarantoolClusterStatus(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.Topology":                 schema_pkg_apis_tarantool_v1alpha1_Topology(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TopologySpec":             schema_pkg_apis_tarantool_v1alpha1_TopologySpec(ref),
		"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TopologyStatus":           schema_pkg_apis_tarantool_v1alpha1_TopologyStatus(ref),
	}
}

func schema_pkg_apis_tarantool_v1alpha1_Cluster(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Cluster is the Schema for the clusters API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ClusterSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ClusterStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ClusterSpec", "gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ClusterStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_ClusterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ClusterSpec defines the desired state of Cluster",
				Properties: map[string]spec.Schema{
					"topologyService": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"topologyServiceType": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"roles": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.Role"),
									},
								},
							},
						},
					},
					"selector": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.Role", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_ClusterStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ClusterStatus defines the observed state of Cluster",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_ReplicasetTemplate(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ReplicasetTemplate is the Schema for the replicasettemplates API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplateSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplateStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplateSpec", "gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.ReplicasetTemplateStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_ReplicasetTemplateSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ReplicasetTemplateSpec defines the desired state of ReplicasetTemplate",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_ReplicasetTemplateStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ReplicasetTemplateStatus defines the observed state of ReplicasetTemplate",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_Role(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Role is the Schema for the roles API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.RoleSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.RoleStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.RoleSpec", "gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.RoleStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_RoleSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "RoleSpec defines the desired state of Role",
				Properties: map[string]spec.Schema{
					"replicas": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"serviceTemplate": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/api/core/v1.ServiceSpec"),
						},
					},
					"selector": {
						SchemaProps: spec.SchemaProps{
							Description: "StorageTemplate *StatefulSetTemplateSpec `json:\"storageTemplate,omitempty\"`",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
				},
				Required: []string{"serviceTemplate"},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.ServiceSpec", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_RoleStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "RoleStatus defines the observed state of Role",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_TarantoolCluster(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TarantoolCluster is the Schema for the tarantoolclusters API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolClusterSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolClusterStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolClusterSpec", "gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TarantoolClusterStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_TarantoolClusterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TarantoolClusterSpec defines the desired state of TarantoolCluster",
				Properties: map[string]spec.Schema{
					"selector": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html\n\tReplicasetTemplate appsv1.StatefulSetSpec `json:\"template,omitempty\"`",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_TarantoolClusterStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TarantoolClusterStatus defines the observed state of TarantoolCluster",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_Topology(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Topology is the Schema for the topologies API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TopologySpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TopologyStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TopologySpec", "gitlab.com/tarantool/sandbox/tarantool-operator/pkg/apis/tarantool/v1alpha1.TopologyStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_TopologySpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TopologySpec defines the desired state of Topology",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_tarantool_v1alpha1_TopologyStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TopologyStatus defines the observed state of Topology",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
