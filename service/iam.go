package service

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	configEvent "github.com/Huawei-APAC-Professional-Services/config-rules/event"
	configModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/config/v1/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
)

func (c *ConfigClient) HasOnlyOneEnterpriseAdministrator(event *event.ConfigEvent) (*configModel.PolicyResource, bool, error) {
	provider := "iam"
	resourceType := "groups"
	comlianceResource := configModel.PolicyResource{}
	comlianceResource.DomainId = event.DomainId
	comlianceResource.RegionId = event.InvokingEvent.RegionId
	comlianceResource.ResourceProvider = &provider
	comlianceResource.ResourceType = &resourceType
	var adminGroupName = "admin"
	slog.Info("Starting List Groups")
	groupqueryResult, err := c.iam.KeystoneListGroups(&model.KeystoneListGroupsRequest{Name: &adminGroupName})
	if err != nil {
		slog.Error("errof", "msg", err.Error())
		return nil, false, err
	}
	slog.Info("Finished List Groups")
	comlianceResource.ResourceName = &adminGroupName
	if len(*groupqueryResult.Groups) != 1 {
		return nil, false, errors.New("there is more than one admin group")
	}
	adminGroupId := (*groupqueryResult.Groups)[0].Id
	slog.Info("Get Admin Group Id", "adminId", adminGroupId)
	comlianceResource.ResourceId = &adminGroupId
	groupusersResult, err := c.iam.KeystoneListUsersForGroupByAdmin(&model.KeystoneListUsersForGroupByAdminRequest{GroupId: adminGroupId})
	if err != nil {
		return nil, false, err
	}
	slog.Info("Finished Listing Users")
	if len(*groupusersResult.Users) <= 1 {
		return &comlianceResource, true, nil
	} else {
		return &comlianceResource, false, nil
	}
}

func (c *ConfigClient) EnsureHasOnlyOneEnterpriseAdministratorPeriodCheck(event *event.ConfigEvent, region string) (*configEvent.ConfigComplianceStatuesReportRequest, error) {
	provider := "iam"
	resourceType := "groups"
	var adminGroupName = "admin"
	result := &configEvent.ConfigComplianceStatuesReportRequest{
		PolicyResource: configEvent.ConfigPolicyResource{
			DomainId:         event.DomainId,
			RegionId:         &region,
			ResourceName:     &adminGroupName,
			ResourceProvider: &provider,
			ResourceType:     &resourceType,
		},
		TriggerType:          event.TriggerType,
		PolicyAssignmentId:   event.AssignmentId,
		PolicyAssignmentName: event.AssignmentName,
		FunctionURN:          event.FunctionURN,
		EvaluationTime:       strconv.FormatInt(event.EvaluationTime, 10),
		EvalutationHash:      event.EvaluationHash,
	}

	groupqueryResult, err := c.iam.KeystoneListGroups(&model.KeystoneListGroupsRequest{Name: &adminGroupName})
	if err != nil {
		slog.Error("errof", "msg", err.Error())
		return nil, err
	}
	slog.Info("Finished List Groups")
	if len(*groupqueryResult.Groups) != 1 {
		return nil, errors.New("there is more than one admin group")
	}
	adminGroupId := (*groupqueryResult.Groups)[0].Id
	slog.Info("Get Admin Group Id", "adminId", adminGroupId)
	result.PolicyResource.ResourceId = &adminGroupId
	groupusersResult, err := c.iam.KeystoneListUsersForGroupByAdmin(&model.KeystoneListUsersForGroupByAdminRequest{GroupId: adminGroupId})
	if err != nil {
		return nil, err
	}
	slog.Info("Finished Listing Users")
	if len(*groupusersResult.Users) <= 1 {
		result.ComplianceState = configEvent.CompliantResult
	} else {
		result.ComplianceState = configEvent.NonCompliantResult
	}
	return result, nil
}
