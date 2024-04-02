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
		tmpl, err := template.New("datasource_role").Parse(s)
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

func TestRoleDataSource_Configure(t *testing.T) {
	testRoleGUID := "34ffe894-6026-4774-93e4-bf9a8e827558"
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
						resource.TestCheckResourceAttr(dataSourceName, "org", testOrg2GUID),
						resource.TestCheckResourceAttr(dataSourceName, "user", testUser2GUID),
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
