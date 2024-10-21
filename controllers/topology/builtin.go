package topology

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/machinebox/graphql"
	corev1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type EditTopologyResponse struct {
	EditTopology EditTopology `json:"edit_topology"`
}

type EditTopology struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	UUID string `json:"uuid"`
}

// ResponseError .
type ResponseError struct {
	Message string `json:"message"`
}

// JoinResponseData .
type JoinResponseData struct {
	JoinInstance bool `json:"joinInstanceResponse"`
}

// JoinResponse .
type JoinResponse struct {
	Errors []*ResponseError  `json:"errors,omitempty"`
	Data   *JoinResponseData `json:"data,omitempty"`
}

// ExpelResponseData .
type ExpelResponseData struct {
	ExpelInstance bool `json:"expel_instance"`
}

// ExpelResponse .
type ExpelResponse struct {
	Errors []*ResponseError   `json:"errors,omitempty"`
	Data   *ExpelResponseData `json:"data,omitempty"`
}

// BootstrapVshardData .
type BootstrapVshardData struct {
	BootstrapVshard bool `json:"bootstrapVshardResponse"`
}

// BootstrapVshardResponse .
type BootstrapVshardResponse struct {
	Data   *BootstrapVshardData `json:"data,omitempty"`
	Errors []*ResponseError     `json:"errors,omitempty"`
}

type FailoverParams struct {
	Mode string `json:"mode"`
}

type ClusterData struct {
	FailoverParams *FailoverParams `json:"failover_params"`
}
type FailoverResponseData struct {
	Cluster *ClusterData `json:"cluster"`
}

// BuiltInTopologyService .
type BuiltInTopologyService struct {
	serviceHost string
	clusterID   string
}

// EditReplicasetResponse .
type EditReplicasetResponse struct {
	Response bool `json:"editReplicasetResponse"`
}

// GetServerStatResponse .
type GetServerStatResponse struct {
	Data   *ServerStatData  `json:"data"`
	Errors []*ResponseError `json:"errors,omitempty"`
}

// ServerStatData .
type ServerStatData struct {
	Stats []*ServerStat `json:"serverStat"`
}

// ServerStat .
type ServerStat struct {
	Statistics Statistics `json:"statistics"`
	UUID       string     `json:"uuid"`
	URI        string     `json:"uri"`
}

// ReplicasetsQueryResponse struct for returning replicasets data
type ReplicasetsQueryResponse struct {
	Replicasets []*ReplicasetData `json:"replicasets"`
}

// ReplicasetData .
type ReplicasetData struct {
	UUID   string   `json:"uuid"`
	Roles  []string `json:"roles"`
	Weight *int     `json:"weight"`
}

// Statistics .
type Statistics struct {
	ItemsUsedRatio string `json:"items_used_ratio"`
	ArenaUsedRatio string `json:"arena_used_ratio"`
	QuotaSize      int    `json:"quotaSize"`
	ArenaUsed      int    `json:"arenaUsed"`
	QuotaUsedRatio string `json:"quota_used_ratio"`
	BucketsCount   int    `json:"bucketsCount"`
}

// join_servers array for editTopology
type EditTopologyServer struct {
	Uri string `json:"uri"`
}

// Replicasets
type Replicasets struct {
	Alias        string               `json:"alias"`
	Roles        []string             `json:"roles"`
	Weight       int                  `json:"weight"`
	All_rw       bool                 `json:"all_rw"`
	Vshard_group string               `json:"vshard_group"`
	Join_servers []EditTopologyServer `json:"join_servers"`
}

var log = logf.Log.WithName("topology")

var (
	errTopologyIsDown      = errors.New("topology service is down")
	errAlreadyJoined       = errors.New("already joined")
	errAlreadyBootstrapped = errors.New("already bootstrapped")
)

var joinMutation = `mutation
	do_join_server(
		$uri: String!,
		$instance_uuid: String!,
		$replicaset_uuid: String!,
		$roles: [String!],
		$vshard_group: String!
	) {
	joinInstanceResponse: join_server(
		uri: $uri,
		instance_uuid: $instance_uuid,
		replicaset_uuid: $replicaset_uuid,
		roles: $roles,
		timeout: 10,
		vshard_group: $vshard_group
	)
}`

var setRsWeightMutation = `mutation editReplicaset($uuid: String!, $weight: Float) {
	editReplicasetResponse: edit_replicaset(uuid: $uuid, weight: $weight)
}`

var getRsWeightQuery = `query ($uuid: String!) {
	replicasets(uuid: $uuid) { weight }
}`

var setRsRolesMutation = `mutation editReplicaset($uuid: String!, $roles: [String!]) {
	editReplicasetResponse: edit_replicaset(uuid: $uuid, roles: $roles)
}`

var getRsRolesQuery = `query ($uuid: String!) {
	replicasets(uuid: $uuid) { roles }
}`

var getServerStatQuery = `query serverList {
	serverStat: servers {
		uuid
		uri
		statistics {
			quotaSize: quota_size
			arenaUsed: arena_used
			bucketsCount: vshard_buckets_count
			quota_used_ratio
			arena_used_ratio
			items_used_ratio
		}
	}
}`

var getFailoverStateQuery = `query {
	cluster {
		failover_params {
			mode
		}
	}
}`

var setFailoverStateQuery = `mutation setFailoverMode($mode: String) {
	cluster {
		failover_params(mode: $mode) {
		  mode
		}
	}
}`

var editTopologyQuery = `mutation editTopology($replicasets: [EditReplicasetInput!], $servers: [EditServerInput!]) {
    cluster {
      edit_topology(replicasets: $replicasets, servers: $servers) {
        servers {
          uuid
        }
      }
    }
}`

// An interface describing an object with accessor methods for labels and annotations
type ObjectWithMeta interface {
	GetLabels() map[string]string
	GetAnnotations() map[string]string
}

// GetRoles comment
func GetRoles(obj ObjectWithMeta) ([]string, error) {
	labels := obj.GetLabels()
	annotations := obj.GetAnnotations()

	rolesFromAnnotations, ok := annotations["tarantool.io/rolesToAssign"]
	if !ok {
		rolesFromLabels, ok := labels["tarantool.io/rolesToAssign"]
		if !ok {
			return nil, errors.New("role undefined")
		}

		roles := strings.Split(rolesFromLabels, ".")
		log.Info("roles", "roles", roles)

		return roles, nil
	}

	var singleRole string
	var roleArray []string

	err := json.Unmarshal([]byte(rolesFromAnnotations), &singleRole)
	if err == nil {
		log.Info("roles", "roles", singleRole)
		return []string{singleRole}, nil
	}

	err = json.Unmarshal([]byte(rolesFromAnnotations), &roleArray)
	if err == nil {
		log.Info("roles", "roles", roleArray)
		return roleArray, nil
	}

	return nil, errors.New("failed to parse roles from annotations")
}

// Harness of edit topology query
func (s *BuiltInTopologyService) EditTopology(replicas map[string]*corev1.Pod, rsetName string) error {
	var replicasUri []EditTopologyServer
	var replicasetsValues Replicasets
	var replicasetsValuesArray []Replicasets
	var lastNode *corev1.Pod

	// Generating URI's array
	for key, pod := range replicas {
		fmt.Println("Recieved values to join: ", key)
		thisPodLabels := pod.GetLabels()
		clusterDomainName, ok := thisPodLabels["tarantool.io/cluster-domain-name"]
		if !ok {
			clusterDomainName = "cluster.local"
		}

		advURI := fmt.Sprintf("%s.%s.%s.svc.%s:3301",
			pod.GetObjectMeta().GetName(),      // Instance name
			s.clusterID,                        // Cartridge cluster name
			pod.GetObjectMeta().GetNamespace(), // Namespace
			clusterDomainName)                  // Cluster domain name
		replicasUri = append(replicasUri, EditTopologyServer{Uri: advURI})

		log.Info("Pod URI: " + advURI)

		lastNode = pod
	}

	// Constructing graphQL variables structure
	// Extracting all from the last node in set
	thisPodLabels := lastNode.GetLabels()

	roles, err := GetRoles(lastNode)
	if err != nil {
		return err
	}
	log.Info("roles", "roles", roles)

	vshardGroup := "default"
	useVshardGroups, ok := thisPodLabels["tarantool.io/useVshardGroups"]
	if !ok {
		return errors.New("failed to get label tarantool.io/useVshardGroups")
	}
	if useVshardGroups == "1" {
		vshardGroup, ok = thisPodLabels["tarantool.io/vshardGroupName"]
		if !ok {
			return errors.New("vshard_group undefined")
		}
	}

	// Generating replicaset name
	// TODO: So far weight can't be read from values.yaml
	replicasetsValues.Alias = rsetName + "_replicaset"
	replicasetsValues.Roles = roles
	replicasetsValues.Weight = 10
	replicasetsValues.Vshard_group = vshardGroup
	replicasetsValues.All_rw = false
	replicasetsValues.Join_servers = replicasUri

	// GraphQL request
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(editTopologyQuery)
	replicasetsValuesArray = append(replicasetsValuesArray, replicasetsValues)
	req.Var("replicasets", replicasetsValuesArray)
	resp := &EditTopologyResponse{}

	if err := client.Run(context.TODO(), req, resp); err != nil {
		return err
	}

	return nil
}

// Join comment
func (s *BuiltInTopologyService) Join(pod *corev1.Pod) error {

	thisPodLabels := pod.GetLabels()
	clusterDomainName, ok := thisPodLabels["tarantool.io/cluster-domain-name"]
	if !ok {
		clusterDomainName = "cluster.local"
	}

	advURI := fmt.Sprintf("%s.%s.%s.svc.%s:3301",
		pod.GetObjectMeta().GetName(),      // Instance name
		s.clusterID,                        // Cartridge cluster name
		pod.GetObjectMeta().GetNamespace(), // Namespace
		clusterDomainName)                  // Cluster domain name

	replicasetUUID, ok := thisPodLabels["tarantool.io/replicaset-uuid"]
	if !ok {
		return errors.New("replicaset uuid empty")
	}

	log.Info("payload", "advURI", advURI, "replicasetUUID", replicasetUUID)

	instanceUUID, ok := thisPodLabels["tarantool.io/instance-uuid"]
	if !ok {
		return errors.New("instance uuid empty")
	}

	roles, err := GetRoles(pod)
	if err != nil {
		return err
	}
	log.Info("roles", "roles", roles)

	vshardGroup := "default"
	useVshardGroups, ok := thisPodLabels["tarantool.io/useVshardGroups"]
	if !ok {
		return errors.New("failed to get label tarantool.io/useVshardGroups")
	}

	if useVshardGroups == "1" {
		vshardGroup, ok = thisPodLabels["tarantool.io/vshardGroupName"]
		if !ok {
			return errors.New("vshard_group undefined")
		}
	}

	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(joinMutation)

	req.Var("uri", advURI)
	req.Var("instance_uuid", instanceUUID)
	req.Var("replicaset_uuid", replicasetUUID)
	req.Var("roles", roles)
	req.Var("vshard_group", vshardGroup)

	resp := &JoinResponseData{}
	if err := client.Run(context.TODO(), req, resp); err != nil {
		if strings.Contains(err.Error(), "already joined") {
			return errAlreadyJoined
		}
		if strings.Contains(err.Error(), "This instance isn't bootstrapped yet") {
			return errTopologyIsDown
		}

		return err
	}

	if resp.JoinInstance {
		return nil
	}

	return errors.New("something really bad happened")
}

// SetFailover enables cluster failover
func (s *BuiltInTopologyService) SetFailover(enabled bool) error {
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(setFailoverStateQuery)

	mode := "eventual"
	if !enabled {
		mode = "disabled"
	}
	req.Var("mode", mode)

	resp := &FailoverResponseData{}
	if err := client.Run(context.TODO(), req, resp); err != nil {
		log.Error(err, "failoverError")
		return errors.New("failed to enable cluster failover")
	}

	return nil
}

func (s *BuiltInTopologyService) GetFailover() (bool, error) {
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(getFailoverStateQuery)

	resp := &FailoverResponseData{}
	if err := client.Run(context.TODO(), req, resp); err != nil {
		log.Error(err, "failoverError")
		return false, errors.New("failed to get info about failover")
	}

	if resp == nil || resp.Cluster == nil || resp.Cluster.FailoverParams == nil || resp.Cluster.FailoverParams.Mode == "" {
		return false, errors.New("failed to get info about failover: broken gql response")
	}

	mode := resp.Cluster.FailoverParams.Mode

	return mode != "disabled", nil
}

// Expel removes an instance from the replicaset
func (s *BuiltInTopologyService) Expel(pod *corev1.Pod) error {
	req := fmt.Sprintf("mutation {expel_instance:expel_server(uuid:\\\"%s\\\")}", pod.GetAnnotations()["tarantool.io/instance_uuid"])
	j := fmt.Sprintf("{\"query\": \"%s\"}", req)
	rawResp, err := http.Post(s.serviceHost, "application/json", strings.NewReader(j))
	if err != nil {
		return err
	}
	defer rawResp.Body.Close()

	resp := &ExpelResponse{Errors: []*ResponseError{}, Data: &ExpelResponseData{}}
	if err := json.NewDecoder(rawResp.Body).Decode(resp); err != nil {
		return err
	}

	if !resp.Data.ExpelInstance && (resp.Errors == nil || len(resp.Errors) == 0) {
		return errors.New("something really bad happened")
	}

	return nil
}

// SetWeight sets weight of a replicaset
func (s *BuiltInTopologyService) SetWeight(replicasetUUID string, replicaWeight string) error {
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(setRsWeightMutation)

	reqLogger := log.WithValues("namespace", "topology.builtin")

	weightParam, err := strconv.ParseUint(replicaWeight, 10, 32)
	if err != nil {
		return err
	}

	reqLogger.Info("setting cluster weight", "uuid", replicasetUUID, "weight", replicaWeight)

	req.Var("uuid", replicasetUUID)
	req.Var("weight", weightParam)

	resp := &EditReplicasetResponse{}
	if err := client.Run(context.TODO(), req, resp); err != nil {
		return err
	}

	if resp.Response {
		return nil
	}

	return errors.New("something really bad happened")
}

// GetWeight gets weight of a replicaset
func (s *BuiltInTopologyService) GetWeight(replicasetUUID string) (int, error) {
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(getRsWeightQuery)

	reqLogger := log.WithValues("namespace", "topology.builtin")

	reqLogger.Info("getting cluster weight", "uuid", replicasetUUID)

	req.Var("uuid", replicasetUUID)

	resp := &ReplicasetsQueryResponse{}
	if err := client.Run(context.TODO(), req, resp); err != nil {
		return -1, err
	}

	if len(resp.Replicasets) == 0 {
		return -1, fmt.Errorf("replicaset with uuid: '%s' not found", replicasetUUID)
	}

	// Instance without role vshard-storage returns null as weight
	if resp.Replicasets[0].Weight == nil {
		return -1, nil
	}

	return *resp.Replicasets[0].Weight, nil
}

// SetReplicasetRoles set roles list of replicaset in the Tarantool service
func (s *BuiltInTopologyService) SetReplicasetRoles(replicasetUUID string, roles []string) error {
	reqLogger := log.WithValues("namespace", "topology.builtin")
	reqLogger.Info("setting replicaset roles", "uuid", replicasetUUID, "weight", roles)

	req := graphql.NewRequest(setRsRolesMutation)
	req.Var("uuid", replicasetUUID)
	req.Var("roles", roles)

	resp := &EditReplicasetResponse{}
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))

	if err := client.Run(context.TODO(), req, resp); err != nil {
		return err
	}
	return nil
}

// GetReplicasetRolesFromService get roles list of replicaset from the Tarantool service
func (s *BuiltInTopologyService) GetReplicasetRolesFromService(replicasetUUID string) ([]string, error) {
	reqLogger := log.WithValues("namespace", "topology.builtin")
	reqLogger.Info("getting replicaset roles", "uuid", replicasetUUID)

	req := graphql.NewRequest(getRsRolesQuery)
	req.Var("uuid", replicasetUUID)

	resp := &ReplicasetsQueryResponse{}
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	if err := client.Run(context.TODO(), req, resp); err != nil {
		return nil, err
	}

	if len(resp.Replicasets) == 0 {
		return nil, fmt.Errorf("replicaset with uuid: '%s' not found", replicasetUUID)
	}
	return resp.Replicasets[0].Roles, nil
}

// GetServerStat Fetch the replicaset as reported by cartridge
func (s *BuiltInTopologyService) GetServerStat() (ServerStatData, error) {
	client := graphql.NewClient(s.serviceHost, graphql.WithHTTPClient(&http.Client{Timeout: time.Duration(time.Second * 5)}))
	req := graphql.NewRequest(getServerStatQuery)

	reqLogger := log.WithValues("function", "GetServerStat")

	reqLogger.Info("fetching server stats")

	resp := ServerStatData{}
	if err := client.Run(context.TODO(), req, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// BootstrapVshard enable the vshard service on the cluster
func (s *BuiltInTopologyService) BootstrapVshard() error {
	reqLogger := log.WithValues("namespace", "topology.builtin")

	reqLogger.Info("Bootstrapping vshard")

	req := "mutation bootstrap {bootstrapVshardResponse: bootstrap_vshard}"
	j := fmt.Sprintf("{\"query\": \"%s\"}", req)
	rawResp, err := http.Post(s.serviceHost, "application/json", strings.NewReader(j))
	if err != nil {
		return err
	}

	defer rawResp.Body.Close()

	resp := &BootstrapVshardResponse{Data: &BootstrapVshardData{}}
	if err := json.NewDecoder(rawResp.Body).Decode(resp); err != nil {
		return err
	}

	if resp.Data.BootstrapVshard {
		return nil
	}
	if resp.Errors != nil && len(resp.Errors) > 0 {
		if strings.Contains(resp.Errors[0].Message, "already bootstrapped") {
			return errAlreadyBootstrapped
		}

		return errors.New(resp.Errors[0].Message)
	}

	return errors.New("unknown error")
}

// IsTopologyDown .
func IsTopologyDown(err error) bool {
	return err == errTopologyIsDown
}

// IsAlreadyJoined .
func IsAlreadyJoined(err error) bool {
	return err == errAlreadyJoined
}

// IsAlreadyBootstrapped .
func IsAlreadyBootstrapped(err error) bool {
	return err == errAlreadyBootstrapped
}

// Option .
type Option func(s *BuiltInTopologyService)

// WithTopologyEndpoint .
func WithTopologyEndpoint(url string) Option {
	return func(s *BuiltInTopologyService) {
		s.serviceHost = url
	}
}

// WithClusterID .
func WithClusterID(id string) Option {
	return func(s *BuiltInTopologyService) {
		s.clusterID = id
	}
}

// NewBuiltInTopologyService .
func NewBuiltInTopologyService(opts ...Option) *BuiltInTopologyService {
	s := &BuiltInTopologyService{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}
