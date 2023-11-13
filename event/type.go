package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const CompliantResult string = "Compliant"
const NonCompliantResult string = "NonCompliant"
const ConfigEndpoint string = "https://rms.myhuaweicloud.com"

// Doc: https://support.huaweicloud.com/intl/en-us/usermanual-rms/rms_05_0506.html
type ConfigEvent struct {
	DomainId       *string                      `json:"domain_id,omitempty"`
	AssignmentId   *string                      `json:"policy_assignment_id,omitempty"`
	AssignmentName *string                      `json:"policy_assignment_name,omitempty"`
	FunctionURN    *string                      `json:"function_urn,omitempty"`
	TriggerType    *string                      `json:"trigger_type,omitempty"`
	EvaluationTime int64                        `json:"evaluation_time,omitempty"`
	EvaluationHash *string                      `json:"evaluation_hash,omitempty"`
	RuleParameter  map[string]map[string]string `json:"rule_parameter,omitempty"`
	InvokingEvent  ConfigInvokingEvent          `json:"invoking_event,omitempty"`
}

type ConfigInvokingEvent struct {
	Id                *string           `json:"id,omitempty"`
	Name              *string           `json:"name,omitempty"`
	Provider          *string           `json:"provider,omitempty"`
	Type              *string           `json:"type,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
	CreateTime        time.Time         `json:"created,omitempty"`
	UpdateTime        time.Time         `json:"updated,omitempty"`
	Properties        interface{}       `json:"properties,omitempty"`
	EP_Id             *string           `json:"ep_id,omitempty"`
	ProjectId         *string           `json:"project_id,omitempty"`
	RegionId          *string           `json:"region_id,omitempty"`
	ProvisioningState *string           `json:"provisioning_state,omitempty"`
}

type ConfigPolicyResource struct {
	DomainId         *string `json:"domain_id,omitempty"`
	RegionId         *string `json:"region_id,omitempty"`
	ResourceId       *string `json:"resource_id,omitempty"`
	ResourceName     *string `json:"resource_name,omitempty"`
	ResourceProvider *string `json:"resource_provider,omitempty"`
	ResourceType     *string `json:"resource_type,omitempty"`
}

type ConfigComplianceStatuesReportRequest struct {
	PolicyResource       ConfigPolicyResource `json:"policy_resource,omitempty"`
	TriggerType          *string              `json:"trigger_type,omitempty"`
	ComplianceState      string               `json:"compliance_state,omitempty"`
	PolicyAssignmentId   *string              `json:"policy_assignment_id,omitempty"`
	PolicyAssignmentName *string              `json:"policy_assignment_name,omitempty"`
	FunctionURN          *string              `json:"function_urn,omitempty"`
	EvaluationTime       int64                `json:"evaluation_time,omitempty"`
	EvalutationHash      *string              `json:"evaluation_hash,omitempty"`
}

// Report Compliance Status to RMS
func (e *ConfigEvent) ReportComplianceStatus(status, token string) error {
	var policyResrouce ConfigPolicyResource
	var complianceRequestData ConfigComplianceStatuesReportRequest
	reportURL := ConfigEndpoint + "/v1/resource-manager/domains/" + *e.DomainId + "/policy-states"
	policyResrouce.DomainId = e.DomainId
	policyResrouce.RegionId = e.InvokingEvent.RegionId
	policyResrouce.ResourceId = e.InvokingEvent.Id
	policyResrouce.ResourceName = e.InvokingEvent.Name
	policyResrouce.ResourceProvider = e.InvokingEvent.Provider
	policyResrouce.ResourceType = e.InvokingEvent.Type
	complianceRequestData.PolicyResource = policyResrouce
	complianceRequestData.TriggerType = e.TriggerType
	complianceRequestData.ComplianceState = status
	complianceRequestData.PolicyAssignmentId = e.AssignmentId
	complianceRequestData.PolicyAssignmentName = e.AssignmentName
	complianceRequestData.FunctionURN = e.FunctionURN
	complianceRequestData.EvaluationTime = e.EvaluationTime
	complianceRequestData.EvalutationHash = e.EvaluationHash
	reqData, err := json.Marshal(complianceRequestData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, reportURL, bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", token)
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if result.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New(result.Status)
}

func (e *ConfigEvent) PeriodReportComplianceStatus(status, token string) error {
	var policyResrouce ConfigPolicyResource
	var complianceRequestData ConfigComplianceStatuesReportRequest
	reportURL := ConfigEndpoint + "/v1/resource-manager/domains/" + *e.DomainId + "/policy-states"
	policyResrouce.DomainId = e.DomainId
	policyResrouce.RegionId = e.InvokingEvent.RegionId
	policyResrouce.ResourceName = e.InvokingEvent.Name
	complianceRequestData.PolicyResource = policyResrouce
	complianceRequestData.TriggerType = e.TriggerType
	complianceRequestData.ComplianceState = status
	complianceRequestData.PolicyAssignmentId = e.AssignmentId
	complianceRequestData.PolicyAssignmentName = e.AssignmentName
	complianceRequestData.FunctionURN = e.FunctionURN
	complianceRequestData.EvaluationTime = e.EvaluationTime
	complianceRequestData.EvalutationHash = e.EvaluationHash
	reqData, err := json.Marshal(complianceRequestData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, reportURL, bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", token)
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if result.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New(result.Status)
}
