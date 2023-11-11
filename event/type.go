package event

import "time"

const CompliantResult string = "Compliant"
const NonCompliantResult string = "NonCompliant"

// Doc: https://support.huaweicloud.com/intl/en-us/usermanual-rms/rms_05_0506.html
type ConfigEvent struct {
	DomainId       string                 `json:"domain_id"`
	AssignmentId   string                 `json:"policy_assignment_id"`
	AssignmentName string                 `json:"policy_assignment_name"`
	FunctionURN    string                 `json:"function_urn"`
	TriggerType    string                 `json:"trigger_type"`
	EvaluationTime int64                  `json:"evaluation_time"`
	EvaluationHash string                 `json:"evaluation_hash"`
	RuleParameter  map[string]interface{} `json:"rule_parameter"`
	InvokingEvent  ConfigInvokingEvent    `json:"invoking_event"`
}

type ConfigInvokingEvent struct {
	Id                string            `json:"id"`
	Name              string            `json:"name"`
	Provider          string            `json:"provider"`
	Type              string            `json:"type"`
	Tags              map[string]string `json:"tags"`
	CreateTime        time.Time         `json:"created"`
	UpdateTime        time.Time         `json:"updated"`
	Properties        interface{}       `json:"properties"`
	EP_Id             string            `json:"ep_id"`
	ProjectId         string            `json:"project_id"`
	RegionId          string            `json:"region_id"`
	ProvisioningState string            `json:"provisioning_state"`
}
