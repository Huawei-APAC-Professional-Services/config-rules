package service

import (
	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	config "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/config/v1"
	configModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/config/v1/model"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hwregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
)

type ConfigClient struct {
	config *config.ConfigClient
	iam    *iam.IamClient
}

func NewConfigClient(auth *global.Credentials, region string) *ConfigClient {
	var client = ConfigClient{}
	client.config = config.NewConfigClient(config.ConfigClientBuilder().WithCredential(auth).WithRegion(hwregion.ValueOf(region)).Build())
	client.iam = iam.NewIamClient(iam.IamClientBuilder().WithRegion(hwregion.ValueOf(region)).WithCredential(auth).Build())
	return &client
}

func (c *ConfigClient) UpdateComplianceStatus(event *event.ConfigEvent, policyResource *configModel.PolicyResource, isCompliance bool) error {
	req := configModel.UpdatePolicyStateRequest{}
	if isCompliance {
		req = configModel.UpdatePolicyStateRequest{
			Body: &configModel.PolicyStateRequestBody{
				PolicyResource:       policyResource,
				TriggerType:          configModel.GetPolicyStateRequestBodyTriggerTypeEnum().PERIOD,
				ComplianceState:      configModel.GetPolicyStateRequestBodyComplianceStateEnum().COMPLIANT,
				PolicyAssignmentId:   *event.AssignmentId,
				PolicyAssignmentName: event.AssignmentName,
				EvaluationTime:       *event.EvaluationTime,
				EvaluationHash:       *event.EvaluationHash,
			},
		}

	} else {
		req = configModel.UpdatePolicyStateRequest{
			Body: &configModel.PolicyStateRequestBody{
				PolicyResource:       &configModel.PolicyResource{},
				TriggerType:          configModel.GetPolicyStateRequestBodyTriggerTypeEnum().PERIOD,
				ComplianceState:      configModel.GetPolicyStateRequestBodyComplianceStateEnum().NON_COMPLIANT,
				PolicyAssignmentId:   *event.AssignmentId,
				PolicyAssignmentName: event.AssignmentName,
				EvaluationTime:       *event.EvaluationTime,
				EvaluationHash:       *event.EvaluationHash,
			},
		}
	}
	_, err := c.config.UpdatePolicyState(&req)
	if err != nil {
		return err
	}
	return nil
}
