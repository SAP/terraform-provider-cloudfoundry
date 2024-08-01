package provider

import (
	"bytes"
	"regexp"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type MtaDataSourceModelPtr struct {
	HclType       string
	HclObjectName string
	Space         *string
	Id            *string
	Namespace     *string
	Mtas          *string
	DeployUrl     *string
}

type MtaResourceModelPtr struct {
	HclType              string
	HclObjectName        string
	MtarPath             *string
	MtarUrl              *string
	ExtensionDescriptors *string
	DeployUrl            *string
	Space                *string
	Mta                  *string
	Namespace            *string
	Id                   *string
}

func hclDataSourceMta(mdsmp *MtaDataSourceModelPtr) string {
	if mdsmp != nil {
		s := `
		{{.HclType}} "cloudfoundry_mta" {{.HclObjectName}} {
			{{- if .Space}}
				space = "{{.Space}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .Namespace}}
				namespace = {{.Namespace}}
			{{- end -}}
			{{if .Mtas}}
				mtas = "{{.Mtas}}"
			{{- end -}}
			{{if .DeployUrl}}
				deploy_url = "{{.DeployUrl}}"
			{{- end }}
			}`
		tmpl, err := template.New("datasource_mtar").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, mdsmp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return mdsmp.HclType + ` "cloudfoundry_mta "` + mdsmp.HclObjectName + ` {}`
}

func hclResourceMta(mrmp *MtaResourceModelPtr) string {
	if mrmp != nil {
		s := `
		{{.HclType}} "cloudfoundry_mta" {{.HclObjectName}} {
			{{- if .Space}}
				space = "{{.Space}}"
			{{- end -}}
			{{if .Id}}
				id = "{{.Id}}"
			{{- end -}}
			{{if .Namespace}}
				namespace = "{{.Namespace}}"
			{{- end -}}
			{{if .Mta}}
				mta = "{{.Mta}}"
			{{- end -}}
			{{if .MtarPath}}
				mtar_path = "{{.MtarPath}}"
			{{- end -}}
			{{if .MtarUrl}}
				mtar_url = "{{.MtarUrl}}"
			{{- end -}}
			{{if .ExtensionDescriptors}}
				extension_descriptors = {{.ExtensionDescriptors}}
			{{- end -}}
			{{if .DeployUrl}}
				deploy_url = "{{.DeployUrl}}"
			{{- end }}
			}`
		tmpl, err := template.New("resource_mtar").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, mrmp)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return mrmp.HclType + ` "cloudfoundry_mta "` + mrmp.HclObjectName + ` {}`
}

func TestMtaDataSource_Configure(t *testing.T) {
	var (
		//canary->tf-space-1
		mtaId     = "a.cf.app"
		spaceGuid = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
	)
	t.Parallel()
	dataSourceName := "data.cloudfoundry_mta.ds"
	t.Run("happy path - read mtar", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_mta")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceMta(&MtaDataSourceModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "ds",
						Space:         strtostrptr(spaceGuid),
						Id:            strtostrptr(mtaId),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "mtas.0.metadata.id", mtaId),
						resource.TestCheckResourceAttr(dataSourceName, "space", spaceGuid),
						resource.TestCheckNoResourceAttr(dataSourceName, "deploy_url"),
					),
				},
			},
		})
	})
	t.Run("error path - mtar does not exist", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_mta_invalid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + hclDataSourceMta(&MtaDataSourceModelPtr{
						HclType:       hclObjectDataSource,
						HclObjectName: "ds",
						Space:         strtostrptr(spaceGuid),
						Id:            strtostrptr(invalidOrgGUID),
					}),
					ExpectError: regexp.MustCompile(`Unable to fetch MTA details`),
				},
			},
		})
	})
}
