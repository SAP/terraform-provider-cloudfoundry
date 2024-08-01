package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMtaResource_Configure(t *testing.T) {
	var (
		resourceName         = "cloudfoundry_mta.rs"
		spaceGuid            = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
		namespace            = "test"
		mtarPath             = "../../assets/a.cf.app.mtar"
		mtarPath2            = "../../assets/my-mta_1.0.0.mtar"
		mtarUrl              = "https://github.com/Dray56/mtar-archive/releases/download/v1.0.0/a.cf.app.mtar"
		extensionDescriptors = `["../../assets/prod.mtaext","../../assets/prod-scale-vertically.mtaext"]`
	)
	t.Parallel()
	t.Run("happy path - create/update/delete mta", func(t *testing.T) {

		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_mta")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						MtarPath:      strtostrptr(mtarPath),
						Space:         strtostrptr(spaceGuid),
						Namespace:     strtostrptr(namespace),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "mtar_path", mtarPath),
						resource.TestCheckResourceAttr(resourceName, "space", spaceGuid),
						resource.TestCheckResourceAttr(resourceName, "mta.metadata.namespace", namespace),
					),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						MtarUrl:       strtostrptr(mtarUrl),
						Space:         strtostrptr(spaceGuid),
						Namespace:     strtostrptr(namespace),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "mtar_url", mtarUrl),
						resource.TestCheckNoResourceAttr(resourceName, "mtar_path"),
						resource.TestCheckResourceAttr(resourceName, "space", spaceGuid),
						resource.TestCheckResourceAttr(resourceName, "mta.metadata.namespace", namespace),
					),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:              hclObjectResource,
						HclObjectName:        "rs",
						MtarPath:             strtostrptr(mtarPath2),
						Space:                strtostrptr(spaceGuid),
						Namespace:            strtostrptr(namespace),
						ExtensionDescriptors: strtostrptr(extensionDescriptors),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "mtar_path", mtarPath2),
						resource.TestCheckResourceAttr(resourceName, "space", spaceGuid),
						resource.TestCheckResourceAttr(resourceName, "mta.metadata.namespace", namespace),
					),
				},
			},
		})
	})
	t.Run("error path - create mtar from invalid path/file", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_mta_invalid_mta_path")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Space:         strtostrptr(spaceGuid),
						MtarPath:      strtostrptr(invalidOrgGUID),
					}),
					ExpectError: regexp.MustCompile(`Unable to upload mtar file`),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Space:         strtostrptr(spaceGuid),
						MtarPath:      strtostrptr(""),
					}),
					ExpectError: regexp.MustCompile(`Unable to upload mtar file`),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Space:         strtostrptr(spaceGuid),
						MtarPath:      strtostrptr("../../assets/provider-config-local.txt"),
					}),
					ExpectError: regexp.MustCompile(`MTA ID missing`),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:              hclObjectResource,
						HclObjectName:        "rs",
						Space:                strtostrptr(spaceGuid),
						MtarPath:             strtostrptr(mtarPath),
						ExtensionDescriptors: strtostrptr(`["../../assets/pr"]`),
					}),
					ExpectError: regexp.MustCompile(`Unable to upload mta extension descriptor`),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:              hclObjectResource,
						HclObjectName:        "rs",
						Space:                strtostrptr(spaceGuid),
						MtarPath:             strtostrptr(mtarPath),
						ExtensionDescriptors: strtostrptr(`[""]`),
					}),
					ExpectError: regexp.MustCompile(`Unable to upload mta extension descriptor`),
				},
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:              hclObjectResource,
						HclObjectName:        "rs",
						Space:                strtostrptr(spaceGuid),
						MtarPath:             strtostrptr(mtarPath),
						ExtensionDescriptors: strtostrptr(`["../../assets/provider-config-local.txt"]`),
					}),
					ExpectError: regexp.MustCompile(`Failure in polling MTA operation`),
				},
			},
		})
	})
	t.Run("error path - create mtar for invalid namespace", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_mta_invalid_namespace")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						MtarPath:      strtostrptr(mtarPath),
						Space:         strtostrptr(spaceGuid),
						Namespace:     strtostrptr("Hello"),
					}),
					ExpectError: regexp.MustCompile(`Failure in polling MTA operation`),
				},
			},
		})
	})
	t.Run("error path - create mtar from empty URL", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_mta_invalid_empty_url")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclResourceMta(&MtaResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Space:         strtostrptr(spaceGuid),
						MtarUrl:       strtostrptr(""),
					}),
					ExpectError: regexp.MustCompile(`Unable to upload remote mtar file`),
				},
			},
		})
	})
}
