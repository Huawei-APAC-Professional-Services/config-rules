package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
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
	EvaluationTime int64                        `json:"evaluation_time"`
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

// https://rms.myhuaweicloud.com/v1/resource-manager/domains/{domain_id}/policy-states

type ConfigComplianceStatuesReportRequest struct {
	PolicyResource       ConfigPolicyResource `json:"policy_resource"`
	TriggerType          *string              `json:"trigger_type"`
	ComplianceState      string               `json:"compliance_state"`
	PolicyAssignmentId   *string              `json:"policy_assignment_id"`
	PolicyAssignmentName *string              `json:"policy_assignment_name"`
	FunctionURN          *string              `json:"function_urn"`
	EvaluationTime       string               `json:"evaluation_time"`
	EvalutationHash      *string              `json:"evaluation_hash"`
}

func (c *ConfigComplianceStatuesReportRequest) UpdatePolicyState(token string) error {
	endpoint := "https://rms.myhuaweicloud.com/v1/resource-manager/domains/https://rms.myhuaweicloud.com/v1/resource-manager/domains/" + *c.PolicyResource.DomainId + "/policy-states"
	slog.Info("marshar policy state body")
	body, err := json.Marshal(c)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("finish marshal policy state body")
	slog.Info(string(body))
	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("X-Security-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Info(err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Info("Failed to update policy state")
		defer resp.Body.Close()
		respData, _ := io.ReadAll(resp.Body)
		slog.Info(string(respData))
		return errors.New(resp.Status)
	}
	return nil
}
