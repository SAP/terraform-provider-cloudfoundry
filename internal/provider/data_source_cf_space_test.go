package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type SpaceDataSourceModelPtr struct {
	Name        *string
	Id          *string
	OrgName     *string
	Org         *string
	Quota       *string
	Labels      *map[string]string
	Annotations *map[string]string
}

func hclDataSourceSpace(sdsmp *SpaceDataSourceModelPtr) string {
	if sdsmp != nil {
		s := `
			data "cloudfoundry_space" "ds" {
			{{- if .Name}}
				name = "{{.Name}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .OrgName}}
				org_name = "{{.OrgName}}"
			{{- end -}}
			{{if .Org}}
				org = "{{.Org}}"
			{{- end -}}
			{{if .Quota}}
				quota = "{{.Quota}}"
			{{- end -}}
			{{if .Labels}}
				labels = "{{.Labels}}"
			{{- end -}}
			{{if .Annotations}}
				annotations = "{{.Annotations}}"
			{{- end }}
			}`
		tmpl, err := template.New("datasource_space").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, sdsmp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return `data "cloudfoundry_space" "ds" {}`
}

func TestSpaceDataSource_Configure(t *testing.T) {
	t.Parallel()
	t.Run("get available datasource space by orgID", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_space_orgid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceSpace(&SpaceDataSourceModelPtr{
						Name: strtostrptr("PerformanceTeamBLR"),
						Org:  strtostrptr("784b4cd0-4771-4e4d-9052-a07e178bae56"),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.cloudfoundry_space.ds", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.cloudfoundry_space.ds", "org_name", "PerformanceTeamBLR"),
						resource.TestCheckResourceAttr("data.cloudfoundry_space.ds", "quota", ""),
					),
				},
			},
		})
	})
	t.Run("get available datasource space by org_name", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_space_orgname")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceSpace(&SpaceDataSourceModelPtr{
						Name:    strtostrptr("PerformanceTeamBLR"),
						OrgName: strtostrptr("PerformanceTeamBLR"),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.cloudfoundry_space.ds", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.cloudfoundry_space.ds", "org", "784b4cd0-4771-4e4d-9052-a07e178bae56"),
						resource.TestCheckResourceAttr("data.cloudfoundry_space.ds", "quota", ""),
					),
				},
			},
		})
	})
	t.Run("error path - get unavailable datasource space", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_space_invalid_spacename")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceSpace(&SpaceDataSourceModelPtr{
						Name: strtostrptr("PerformanceTeamBLR12"),
						Org:  strtostrptr("784b4cd0-4771-4e4d-9052-a07e178bae56"),
					}),
					ExpectError: regexp.MustCompile(`Error: Unable to find space data in list`),
				},
			},
		})
	})
	t.Run("error path - org does not exist", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_space_invalid_orgname")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceSpace(&SpaceDataSourceModelPtr{
						Name:    strtostrptr("PerformanceTeamBLR"),
						OrgName: strtostrptr("PerformanceTeamBLR12"),
					}),
					ExpectError: regexp.MustCompile(`Error: Unable to find org data in list`),
				},
			},
		})
	})
	t.Run("error path - missing org attributes", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_space_invalid_attributes")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceSpace(&SpaceDataSourceModelPtr{
						Name: strtostrptr("PerformanceTeamBLR"),
					}),
					ExpectError: regexp.MustCompile(`Error: Neither Org GUID nor Org Name is present`),
				},
			},
		})
	})
	t.Run("error path - both org attributes provided", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_space_invalid_attributes")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceSpace(&SpaceDataSourceModelPtr{
						Name:    strtostrptr("PerformanceTeamBLR"),
						Org:     strtostrptr("784b4cd0-4771-4e4d-9052-a07e178bae56"),
						OrgName: strtostrptr("PerformanceTeamBLR"),
					}),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Combination`),
				},
			},
		})
	})
}
