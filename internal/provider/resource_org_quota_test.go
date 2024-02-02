package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceOrgQuota(t *testing.T) {
	t.Parallel()
	resourceName := "cloudfoundry_org_quota.rs"
	t.Run("happy path - create org quota", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_org_quota")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclOrgQuota(&OrgQuotaModelPtr{
						Name:                  strtostrptr("tf-unit-test"),
						AllowPaidServicePlans: booltoboolptr(true),
						HclType:               "resource",
						HclObjectName:         "rs",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr(resourceName, "updated_at", regexpValidRFC3999Format),
					),
				},
			},
		})
	})
	t.Run("happy path - import org quota", func(t *testing.T) {
		resourceName := "cloudfoundry_org_quota.rs_import"
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_org_quota_import")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclOrgQuota(&OrgQuotaModelPtr{
						HclType:               hclObjectResource,
						HclObjectName:         "rs_import",
						AllowPaidServicePlans: booltoboolptr(false),
						Name:                  strtostrptr("tf-unit-test-import"),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
					),
				},
				{
					ResourceName:      resourceName,
					ImportStateIdFunc: getIdForImport(resourceName),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}
