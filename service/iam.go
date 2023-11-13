package service

import (
	"errors"
	"log/slog"

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

func (i *ConfigIAMClient) HasOnlyOneEnterpriseAdministrator() (bool, error) {
	var adminGroupName = "admin"
	groupqueryResult, err := i.client.KeystoneListGroups(&model.KeystoneListGroupsRequest{Name: &adminGroupName})
	if err != nil {
		return false, err
	}
	slog.Info("Check admin group id")
	if len(*groupqueryResult.Groups) != 1 {
		return false, errors.New("there is more than one admin group")
	}
	slog.Info("Get admin group id")
	adminGroupId := (*groupqueryResult.Groups)[0].Id
	groupusersResult, err := i.client.KeystoneListUsersForGroupByAdmin(&model.KeystoneListUsersForGroupByAdminRequest{GroupId: adminGroupId})
	if len(*groupusersResult.Users) <= 1 {
		return true, nil
	}
	return false, nil
}
