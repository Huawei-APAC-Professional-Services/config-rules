package event

import (
	"time"
)

const CompliantResult string = "Compliant"
const NonCompliantResult string = "NonCompliant"
const ConfigEndpoint string = "https://rms.myhuaweicloud.com"

// Doc: https://support.huaweicloud.com/intl/en-us/usermanual-rms/rms_05_0506.html
type ConfigEvent struct {
	DomainId       *string                      `json:"domain_id"`
	AssignmentId   *string                      `json:"policy_assignment_id"`
	AssignmentName *string                      `json:"policy_assignment_name"`
	FunctionURN    *string                      `json:"function_urn"`
	TriggerType    *string                      `json:"trigger_type"`
	EvaluationTime *string                      `json:"evaluation_time"`
	EvaluationHash *string                      `json:"evaluation_hash"`
	RuleParameter  map[string]map[string]string `json:"rule_parameter"`
	InvokingEvent  ConfigInvokingEvent          `json:"invoking_event"`
}

type ConfigInvokingEvent struct {
	Id                *string           `json:"id"`
	Name              *string           `json:"name"`
	Provider          *string           `json:"provider"`
	Type              *string           `json:"type"`
	Tags              map[string]string `json:"tags"`
	CreateTime        time.Time         `json:"created"`
	UpdateTime        time.Time         `json:"updated"`
	Properties        interface{}       `json:"properties"`
	EP_Id             *string           `json:"ep_id"`
	ProjectId         *string           `json:"project_id"`
	RegionId          *string           `json:"region_id"`
	ProvisioningState *string           `json:"provisioning_state"`
}

type ConfigPolicyResource struct {
	DomainId         *string `json:"domain_id"`
	RegionId         *string `json:"region_id"`
	ResourceId       *string `json:"resource_id"`
	ResourceName     *string `json:"resource_name"`
	ResourceProvider *string `json:"resource_provider"`
	ResourceType     *string `json:"resource_type"`
}

type PeriodReportResource struct {
	ResourceId       string
	ResourceName     string
	ResourceProvider string
	ResourceType     string
	ComplianceStatus string
}

type ConfigComplianceStatuesReportRequest struct {
	PolicyResource       ConfigPolicyResource `json:"policy_resource"`
	TriggerType          *string              `json:"trigger_type"`
	ComplianceState      string               `json:"compliance_state"`
	PolicyAssignmentId   *string              `json:"policy_assignment_id"`
	PolicyAssignmentName *string              `json:"policy_assignment_name"`
	FunctionURN          *string              `json:"function_urn"`
	EvaluationTime       int64                `json:"evaluation_time"`
	EvalutationHash      *string              `json:"evaluation_hash"`
}
