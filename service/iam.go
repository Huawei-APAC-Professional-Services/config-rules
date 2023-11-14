package service

import (
	"errors"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
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
	groupqueryResult, err := c.iam.KeystoneListGroups(&model.KeystoneListGroupsRequest{Name: &adminGroupName})
	if err != nil {
		return nil, false, err
	}
	comlianceResource.ResourceName = &adminGroupName
	if len(*groupqueryResult.Groups) != 1 {
		return nil, false, errors.New("there is more than one admin group")
	}
	adminGroupId := (*groupqueryResult.Groups)[0].Id
	comlianceResource.ResourceId = &adminGroupId
	groupusersResult, err := c.iam.KeystoneListUsersForGroupByAdmin(&model.KeystoneListUsersForGroupByAdminRequest{GroupId: adminGroupId})
	if err != nil {
		return nil, false, err
	}
	if len(*groupusersResult.Users) <= 1 {
		return &comlianceResource, true, nil
	} else {
		return &comlianceResource, false, nil
	}
}
