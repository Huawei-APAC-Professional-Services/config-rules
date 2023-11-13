package service

import (
	"errors"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hwregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
)

type ConfigIAMClient struct {
	client *iam.IamClient
}

func NewIAMClient(ak, sk, domainId, region string) *ConfigIAMClient {
	var client = ConfigIAMClient{}
	auth := global.NewCredentialsBuilder().WithAk(ak).WithSk(sk).WithDomainId(domainId).Build()
	client.client = iam.NewIamClient(iam.IamClientBuilder().WithRegion(hwregion.ValueOf(region)).WithCredential(auth).Build())
	return &client
}

func (i *ConfigIAMClient) HasOnlyOneEnterpriseAdministrator() (*event.PeriodReportResource, error) {
	result := event.PeriodReportResource{ResourceProvider: "iam", ResourceType: "groups"}
	var adminGroupName = "admin"
	groupqueryResult, err := i.client.KeystoneListGroups(&model.KeystoneListGroupsRequest{Name: &adminGroupName})
	if err != nil {
		return nil, err
	}
	result.ResourceName = adminGroupName
	if len(*groupqueryResult.Groups) != 1 {
		return nil, errors.New("there is more than one admin group")
	}
	adminGroupId := (*groupqueryResult.Groups)[0].Id
	result.ResourceId = adminGroupId
	groupusersResult, err := i.client.KeystoneListUsersForGroupByAdmin(&model.KeystoneListUsersForGroupByAdminRequest{GroupId: adminGroupId})
	if len(*groupusersResult.Users) <= 1 {
		result.ComplianceStatus = event.CompliantResult
	} else {
		result.ComplianceStatus = event.NonCompliantResult
	}
	return &result, nil
}
