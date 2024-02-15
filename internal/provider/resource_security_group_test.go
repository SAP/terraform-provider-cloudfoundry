package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSecurityGroupResource_Configure(t *testing.T) {
	t.Parallel()
	t.Run("happy path - create/read/update/delete security group", func(t *testing.T) {
		resourceName := "cloudfoundry_security_group.ds"
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_security_group_crud")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclSecurityGroup(&SecurityGroupModelPtr{
						HclType:                hclObjectResource,
						HclObjectName:          "ds",
						Name:                   strtostrptr("tf-unit-test"),
						GloballyEnabledRunning: booltoboolptr(false),
						GloballyEnabledStaging: booltoboolptr(false),
						Rules:                  strtostrptr(createRules),
						RunningSpaces:          &runningSpaces,
						StagingSpaces:          &stagingSpaces,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestCheckResourceAttr(resourceName, "name", "tf-unit-test"),
						resource.TestCheckResourceAttr(resourceName, "globally_enabled_running", "false"),
						resource.TestCheckResourceAttr(resourceName, "globally_enabled_staging", "false"),
						resource.TestCheckResourceAttr(resourceName, "rules.#", "3"),
						resource.TestCheckResourceAttr(resourceName, "running_spaces.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "staging_spaces.#", "2"),
					),
				},
				{
					Config: hclProvider(nil) + hclSecurityGroup(&SecurityGroupModelPtr{
						HclType:                hclObjectResource,
						HclObjectName:          "ds",
						Name:                   strtostrptr("tf-unit-test1"),
						GloballyEnabledRunning: booltoboolptr(true),
						GloballyEnabledStaging: booltoboolptr(true),
						RunningSpaces:          &stagingSpaces,
						StagingSpaces:          &runningSpaces,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tf-unit-test1"),
						resource.TestCheckResourceAttr(resourceName, "globally_enabled_running", "true"),
						resource.TestCheckResourceAttr(resourceName, "globally_enabled_staging", "true"),
						resource.TestCheckNoResourceAttr(resourceName, "rules"),
						resource.TestCheckResourceAttr(resourceName, "running_spaces.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "staging_spaces.#", "1"),
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

	t.Run("error path - invalid rule when creating security group", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_security_group_invalid_rule")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclSecurityGroup(&SecurityGroupModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "ds_rule",
						Name:          strtostrptr("tf-unit-test"),
						Rules:         &invalidRules,
					}),
					ExpectError: regexp.MustCompile(`API Error Creating Security Group`),
				},
			},
		})
	})
	t.Run("error path - invalid name when creating security group", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_security_group_invalid_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclSecurityGroup(&SecurityGroupModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "ds_invalid_name",
						Name:          &testSpace,
					}),
					ExpectError: regexp.MustCompile(`API Error Creating Security Group`),
				},
			},
		})
	})

}
