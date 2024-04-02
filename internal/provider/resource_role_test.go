package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type RoleResourceModelPtr struct {
	HclType       string
	HclObjectName string
	ObjectName    string
	UserName      *string
	Origin        *string
	Type          *string
	User          *string
	Space         *string
	Id            *string
	Organization  *string
	CreatedAt     *string
	UpdatedAt     *string
}

func hclRoleResource(rrmp *RoleResourceModelPtr) string {
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
			{{if .UserName}}
				username = "{{.UserName}}"
			{{- end -}}
			{{if .Origin}}
				origin = "{{.Origin}}"
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
		testSpaceGUID = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
		testUserGUID  = "efed91b4-d808-40fb-b4b8-7a9449ff1e15"
		testUserName  = "kesavan.s@sap.com"
		//testUserName2 = "debaditya.ray@sap.com"
		origin = "sap.ids"
	)
	t.Parallel()
	t.Run("happy path - create space role", func(t *testing.T) {
		resourceName := "cloudfoundry_role.rf"
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_space")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRoleResource(&RoleResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rf",
						Type:          strtostrptr("space_auditor"),
						UserName:      strtostrptr(testUserName),
						Origin:        strtostrptr(origin),
						Space:         strtostrptr(testSpaceGUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr(resourceName, "user", regexpValidUUID),
						resource.TestCheckResourceAttr(resourceName, "space", testSpaceGUID),
					),
				},
			},
		})
	})
	t.Run("happy path - create org role", func(t *testing.T) {
		resourceName := "cloudfoundry_role.rs"
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_org")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRoleResource(&RoleResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs",
						Type:          strtostrptr("organization_user"),
						User:          strtostrptr(testUser2GUID),
						Organization:  strtostrptr(testOrg2GUID),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", regexpValidUUID),
						resource.TestMatchResourceAttr(resourceName, "created_at", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr(resourceName, "user", testUser2GUID),
						resource.TestCheckResourceAttr(resourceName, "org", testOrg2GUID),
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
	t.Run("error path - create role with invalid type", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_role_invalid-type")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclRoleResource(&RoleResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rs_invalid",
						Type:          strtostrptr("organization_manager"),
						User:          strtostrptr(testUserGUID),
						Space:         strtostrptr(testSpaceGUID),
					}),
					ExpectError: regexp.MustCompile(`Invalid Role Type`),
				},
				{
					Config: hclProvider(nil) + hclRoleResource(&RoleResourceModelPtr{
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
					Config: hclProvider(nil) + hclRoleResource(&RoleResourceModelPtr{
						HclType:       hclObjectResource,
						HclObjectName: "rsi",
						Type:          strtostrptr("organization_manager"),
						User:          strtostrptr(testUser2GUID),
						Organization:  strtostrptr(testOrg2GUID),
					}),
					ExpectError: regexp.MustCompile(`API Error Registering Role`),
				},
			},
		})
	})

}
