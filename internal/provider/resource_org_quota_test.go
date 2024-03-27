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
	t.Run("happy path - update org quota", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_org_quota_update")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclOrgQuota(&OrgQuotaModelPtr{
						Name:                  strtostrptr("tf-unit-test-update"),
						AllowPaidServicePlans: booltoboolptr(true),
						InstanceMemory:        inttointptr(2048),
						TotalMemory:           inttointptr(51200),
						TotalAppInstances:     inttointptr(100),
						TotalAppTasks:         inttointptr(100),
						TotalPrivateDomains:   inttointptr(100),
						TotalServices:         inttointptr(100),
						TotalAppLogRateLimit:  inttointptr(100),
						HclType:               "resource",
						HclObjectName:         "rs",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestCheckResourceAttr(resourceName, "allow_paid_service_plans", "true"),
						resource.TestCheckResourceAttr(resourceName, "instance_memory", "2048"),
						resource.TestCheckResourceAttr(resourceName, "total_memory", "51200"),
						resource.TestCheckResourceAttr(resourceName, "total_app_instances", "100"),
						resource.TestCheckResourceAttr(resourceName, "total_app_tasks", "100"),
						resource.TestCheckResourceAttr(resourceName, "total_private_domains", "100"),
						resource.TestCheckResourceAttr(resourceName, "total_services", "100"),
						resource.TestCheckResourceAttr(resourceName, "total_app_log_rate_limit", "100"),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr(resourceName, "updated_at", regexpValidRFC3999Format),
					),
				},
				{
					Config: hclProvider(nil) + hclOrgQuota(&OrgQuotaModelPtr{
						Name:                  strtostrptr("tf-unit-test-update"),
						AllowPaidServicePlans: booltoboolptr(false),
						InstanceMemory:        inttointptr(1024),
						TotalMemory:           inttointptr(51200),
						TotalAppInstances:     inttointptr(100),
						TotalAppTasks:         inttointptr(100),
						HclType:               "resource",
						HclObjectName:         "rs",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestCheckResourceAttr(resourceName, "allow_paid_service_plans", "false"),
						resource.TestCheckResourceAttr(resourceName, "instance_memory", "1024"),
						resource.TestCheckResourceAttr(resourceName, "total_memory", "51200"),
						resource.TestCheckResourceAttr(resourceName, "total_app_instances", "100"),
						resource.TestCheckResourceAttr(resourceName, "total_app_tasks", "100"),
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
						resource.TestCheckResourceAttr(resourceName, "allow_paid_service_plans", "false"),
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
