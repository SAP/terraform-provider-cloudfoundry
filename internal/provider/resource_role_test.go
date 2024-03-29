package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type RoleModelPtr struct {
	HclType       string
	HclObjectName string
	ObjectName    string
	Type          *string
	User          *string
	Space         *string
	Id            *string
	Organization  *string
	CreatedAt     *string
	UpdatedAt     *string
}

func hclRole(rrmp *RoleModelPtr) string {
	if rrmp != nil {
		s := `
		{{.HclType}} "cloudfoundry_role" {{.HclObjectName}} {
			{{- if .Type}}
				type = "{{.Type}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .User}}
				user = "{{.User}}"
			{{- end -}}
			{{if .Organization}}
				org = "{{.Organization}}"
			{{- end -}}
			{{if .Space}}
				space = "{{.Space}}"
			{{- end -}}
			{{if .CreatedAt}}
				created_at = "{{.CreatedAt}}"
			{{- end -}}
			{{if .UpdatedAt}}
				updated_at = "{{.UpdatedAt}}"
			{{- end }}
			}`
		tmpl, err := template.New("resource_role").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, rrmp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return rrmp.HclType + ` "cloudfoundry_role "` + rrmp.HclObjectName + ` {}`
}

func TestRoleResource_Configure(t *testing.T) {
	var (
		// in canary -> PerformanceTeamBLR -> tf-space-1
		testUserGUID  = "4467eb10-a5dd-4c46-904f-d5a1c86f05a2"
		testSpaceGUID = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
		testOrgGUID   = "784b4cd0-4771-4e4d-9052-a07e178bae56"
		testUser2GUID = "4595acb0-8a59-4461-83fc-989f0e0c42d7"
	)
	t.Parallel()
	t.Run("happy path - create/import/delete role", func(t *testing.T) {
		resourceName := "cloudfoundry_role.rs"
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_crud")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Type:          strtostrptr("organization_user"),
						User:          strtostrptr(testUserGUID),
						Organization:  strtostrptr(testOrgGUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr(resourceName, "user", testUserGUID),
						resource.TestCheckResourceAttr(resourceName, "org", testOrgGUID),
					),
				},
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Type:          strtostrptr("space_manager"),
						User:          strtostrptr(testUser2GUID),
						Space:         strtostrptr(testSpaceGUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr(resourceName, "user", testUser2GUID),
						resource.TestCheckResourceAttr(resourceName, "space", testSpaceGUID),
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
	t.Run("error path - create space role for user without org role", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_invalid_missing_org_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs_invalid",
						Type:          strtostrptr("space_manager"),
						User:          strtostrptr(testUserGUID),
						Space:         strtostrptr(testSpaceGUID),
					}),
					ExpectError: regexp.MustCompile(`Assigned User missing organization_user role`),
				},
			},
		})
	})
	t.Run("error path - create role with invalid type", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_invalid-type")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs_invalid",
						Type:          strtostrptr("organization_manager"),
						User:          strtostrptr(testUserGUID),
						Space:         strtostrptr(testSpaceGUID),
					}),
					ExpectError: regexp.MustCompile(`Invalid Role Type`),
				},
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs_invalid",
						Type:          strtostrptr("space_manager"),
						User:          strtostrptr(testUserGUID),
						Organization:  strtostrptr(testOrgGUID),
					}),
					ExpectError: regexp.MustCompile(`Invalid Role Type`),
				},
			},
		})
	})
	t.Run("error path - create role with existing id", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_invalid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRole(&RoleModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Type:          strtostrptr("organization_manager"),
						User:          strtostrptr(testUserGUID),
						Organization:  strtostrptr(testOrgGUID),
					}),
					ExpectError: regexp.MustCompile(`API Error Registering Role`),
				},
			},
		})
	})
}
