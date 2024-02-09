package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRoleDataSource_Configure(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.cloudfoundry_role.ds"
	t.Run("happy path - read role", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "ds",
						Id:            strtostrptr(testRoleGUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", testRoleGUID),
						resource.TestCheckResourceAttr(dataSourceName, "org", testOrgGUID),
						resource.TestCheckResourceAttr(dataSourceName, "user", testUser2ID),
					),
				},
			},
		})
	})
	t.Run("error path - role does not exist", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_role_invalid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "ds",
						Id:            strtostrptr(invalidOrgGUID),
					}),
					ExpectError: regexp.MustCompile(`API Error Fetching Role`),
				},
			},
		})
	})
}
