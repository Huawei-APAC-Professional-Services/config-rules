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
	DomainId       string              `json:"domain_id"`
	AssignmentId   string              `json:"policy_assignment_id"`
	AssignmentName string              `json:"policy_assignment_name"`
	FunctionURN    string              `json:"function_urn"`
	TriggerType    string              `json:"trigger_type"`
	EvaluationTime int64               `json:"evaluation_time"`
	EvaluationHash string              `json:"evaluation_hash"`
	RuleParameter  map[string]string   `json:"rule_parameter"`
	InvokingEvent  ConfigInvokingEvent `json:"invoking_event"`
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

type ConfigPolicyResource struct {
	DomainId         string `json:"domain_id"`
	RegionId         string `json:"region_id"`
	ResourceId       string `json:"resource_id"`
	ResourceName     string `json:"resource_name"`
	ResourceProvider string `json:"resource_provider"`
	ResourceType     string `json:"resource_type"`
}

type ConfigComplianceStatuesReportRequest struct {
	PolicyResource       ConfigPolicyResource `json:"policy_resource"`
	TriggerType          string               `json:"trigger_type"`
	ComplianceState      string               `json:"compliance_state"`
	PolicyAssignmentId   string               `json:"policy_assignment_id"`
	PolicyAssignmentName string               `json:"policy_assignment_name"`
	FunctionURN          string               `json:"function_urn"`
	EvaluationTime       int64                `json:"evaluation_time"`
	EvalutationHash      string               `json:"evaluation_hash"`
}

// Report Compliance Status to RMS
func (e *ConfigEvent) ReportComplianceStatus(status, token string) error {
	var policyResrouce ConfigPolicyResource
	var complianceRequestData ConfigComplianceStatuesReportRequest
	reportURL := ConfigEndpoint + "/v1/resource-manager/domains/" + e.DomainId + "/policy-states"
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
