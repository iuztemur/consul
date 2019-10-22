package structs

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/types"
)

// CheckDefinition is used to JSON decode the Check definitions
type CheckDefinition struct {
	ID        types.CheckID
	Name      string
	Notes     string
	ServiceID string
	Token     string
	Status    string

	// Copied fields from CheckType without the fields
	// already present in CheckDefinition:
	//
	//   ID (CheckID), Name, Status, Notes
	//
	ScriptArgs                     []string
	HTTP                           string
	Header                         map[string][]string
	Method                         string
	TCP                            string
	Interval                       time.Duration
	DockerContainerID              string
	Shell                          string
	GRPC                           string
	GRPCUseTLS                     bool
	TLSSkipVerify                  bool
	AliasNode                      string
	AliasService                   string
	Timeout                        time.Duration
	TTL                            time.Duration
	SuccessBeforePassing           int
	FailuresBeforeCritical         int
	DeregisterCriticalServiceAfter time.Duration
	OutputMaxSize                  int
}

func (t *CheckDefinition) UnmarshalJSON(data []byte) (err error) {
	type Alias CheckDefinition
	aux := &struct {
		// Parse special values
		Interval                       interface{}
		Timeout                        interface{}
		TTL                            interface{}
		DeregisterCriticalServiceAfter interface{}

		// Translate fields
		ArgsCamel                           []string    `json:"args"`
		ScriptArgsCamel                     []string    `json:"script_args"`
		DeregisterCriticalServiceAfterCamel interface{} `json:"deregister_critical_service_after"`
		DockerContainerIDCamel              string      `json:"docker_container_id"`
		TLSSkipVerifyCamel                  bool        `json:"tls_skip_verify"`
		ServiceIDCamel                      string      `json:"service_id"`

		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Translate Fields
	if aux.DeregisterCriticalServiceAfter == nil {
		aux.DeregisterCriticalServiceAfter = aux.DeregisterCriticalServiceAfterCamel
	}
	if t.ScriptArgs == nil {
		t.ScriptArgs = aux.ArgsCamel
	}
	if t.ScriptArgs == nil {
		t.ScriptArgs = aux.ScriptArgsCamel
	}
	if t.DockerContainerID == "" {
		t.DockerContainerID = aux.DockerContainerIDCamel
	}
	if aux.TLSSkipVerifyCamel {
		t.TLSSkipVerify = aux.TLSSkipVerifyCamel
	}
	if t.ServiceID == "" {
		t.ServiceID = aux.ServiceIDCamel
	}

	// Parse special values
	if aux.Interval != nil {
		switch v := aux.TTL.(type) {
		case string:
			if t.Interval, err = time.ParseDuration(v); err != nil {
				return err
			}
		case float64:
			t.Interval = time.Duration(v)
		}
	}
	if aux.Timeout != nil {
		switch v := aux.TTL.(type) {
		case string:
			if t.Timeout, err = time.ParseDuration(v); err != nil {
				return err
			}
		case float64:
			t.Timeout = time.Duration(v)
		}
	}
	if aux.TTL != nil {
		switch v := aux.TTL.(type) {
		case string:
			if t.TTL, err = time.ParseDuration(v); err != nil {
				return err
			}
		case float64:
			t.TTL = time.Duration(v)
		}
	}
	if aux.DeregisterCriticalServiceAfter != nil {
		switch v := aux.DeregisterCriticalServiceAfter.(type) {
		case string:
			if t.DeregisterCriticalServiceAfter, err = time.ParseDuration(v); err != nil {
				return err
			}
		case float64:
			t.DeregisterCriticalServiceAfter = time.Duration(v)
		}
	}

	return nil
}

func (c *CheckDefinition) HealthCheck(node string) *HealthCheck {
	health := &HealthCheck{
		Node:      node,
		CheckID:   c.ID,
		Name:      c.Name,
		Status:    api.HealthCritical,
		Notes:     c.Notes,
		ServiceID: c.ServiceID,
	}
	if c.Status != "" {
		health.Status = c.Status
	}
	if health.CheckID == "" && health.Name != "" {
		health.CheckID = types.CheckID(health.Name)
	}
	return health
}

func (c *CheckDefinition) CheckType() *CheckType {
	return &CheckType{
		CheckID: c.ID,
		Name:    c.Name,
		Status:  c.Status,
		Notes:   c.Notes,

		ScriptArgs:                     c.ScriptArgs,
		AliasNode:                      c.AliasNode,
		AliasService:                   c.AliasService,
		HTTP:                           c.HTTP,
		GRPC:                           c.GRPC,
		GRPCUseTLS:                     c.GRPCUseTLS,
		Header:                         c.Header,
		Method:                         c.Method,
		OutputMaxSize:                  c.OutputMaxSize,
		TCP:                            c.TCP,
		Interval:                       c.Interval,
		DockerContainerID:              c.DockerContainerID,
		Shell:                          c.Shell,
		TLSSkipVerify:                  c.TLSSkipVerify,
		Timeout:                        c.Timeout,
		TTL:                            c.TTL,
		SuccessBeforePassing:           c.SuccessBeforePassing,
		FailuresBeforeCritical:         c.FailuresBeforeCritical,
		DeregisterCriticalServiceAfter: c.DeregisterCriticalServiceAfter,
	}
}
