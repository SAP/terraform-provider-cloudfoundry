package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/validation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	hclObjectResource        = "resource"
	hclObjectDataSource      = "data"
	regexpValidUUID          = validation.UuidRegexp
	regexpValidRFC3999Format = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	testOrg                  = "tf-test-do-not-delete"
	testOrgGUID              = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
	testOrg2GUID             = "784b4cd0-4771-4e4d-9052-a07e178bae56"
	testSpace                = "tf-test-do-not-delete"
	testSpaceGUID            = "3bc20dc4-1870-4835-8308-dda2d766e61e"
	testSpace2GUID           = "dd457c79-f7c9-4828-862b-35843d3b646d"
	testOrgQuota             = "tf-test-do-not-delete"
	invalidOrgGUID           = "40b73419-5e01-4be0-baea-932d46cea45b"
	testIsolationSegmentGUID = "5215e4df-79a4-4ce8-a933-837d6aa7a77b"
	testSpaceQuota           = "tf-test-do-not-delete"
	testCreateLabel          = "{purpose: \"testing\", landscape: \"test\"}"
	testUpdateLabel          = "{purpose: \"production\", status: \"fine\"}"
	testUser                 = "debaditya.ray@sap.com"
	testUserGUID             = "2334cf47-fead-4e5f-bd2a-6e7153e7f144"
	testUser2ID              = "space"
	testRoleGUID             = "e17839d9-cd4f-4e4b-baf0-18786f12fede"
	testOrgQuotaGUID         = "e3cef997-9ba5-4cb4-b25b-c79faa81a33f"
	createRules              = `[{
									protocol = "tcp"
									destination = "192.168.1.100"
									ports = "1883,8883"
									log = true
								},{
									protocol = "udp"
									destination = "192.168.1.100"
									ports = "1883,8883"
									log = false
								},
								{
									protocol = "icmp"
									type = 0
									code = 0
									destination = "192.168.1.100"
									log = false
								}]`
	invalidRules = `[{
									protocol = "tcp"
									type = 0
									code = 0
									destination = "192.168.1.100"
									log = false
								}]`
	stagingSpaces           = "[\"3bc20dc4-1870-4835-8308-dda2d766e61e\", \"e6886bba-e263-4b52-aaf1-85d410f15fc8\"]"
	runningSpaces           = "[\"e6886bba-e263-4b52-aaf1-85d410f15fc8\"]"
	testServiceInstance     = "tf-test-do-not-delete" // in canary -> PerformanceTeamBLR -> tf-space-1
	testServiceInstanceGUID = "5e2976bb-332e-41e1-8be3-53baafea9296"
)

func (cfg *CloudFoundryProviderConfigPtr) GetHook() func(i *cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		type reg struct {
			regexpattern []*regexp.Regexp
			redactString string
		}
		var regList []reg
		if cfg.Endpoint != nil {
			u := *cfg.Endpoint
			u = u[strings.Index(u, "."):]
			reg := reg{
				regexpattern: []*regexp.Regexp{
					regexp.MustCompile(u),
					regexp.MustCompile(url.QueryEscape(u)),
				},
				redactString: ".x.x.x.x.com",
			}
			regList = append(regList, reg)
		}
		if cfg.User != nil {
			reg := reg{
				regexpattern: []*regexp.Regexp{
					regexp.MustCompile(*cfg.User),
					regexp.MustCompile(url.QueryEscape(*cfg.User)),
				},
				redactString: *redactedTestUser.User,
			}
			regList = append(regList, reg)
		}
		if cfg.Password != nil {
			reg := reg{
				regexpattern: []*regexp.Regexp{
					regexp.MustCompile(*cfg.Password),
					regexp.MustCompile(url.QueryEscape(*cfg.Password)),
				},
				redactString: *redactedTestUser.Password,
			}
			regList = append(regList, reg)
		}
		if cfg.CFClientID != nil {
			reg := reg{
				regexpattern: []*regexp.Regexp{
					regexp.MustCompile(*cfg.CFClientID),
					regexp.MustCompile(url.QueryEscape(*cfg.CFClientID)),
				},
				redactString: *redactedTestUser.CFClientID,
			}
			regList = append(regList, reg)
		}
		if cfg.CFClientSecret != nil {
			reg := reg{
				regexpattern: []*regexp.Regexp{
					regexp.MustCompile(*cfg.CFClientSecret),
					regexp.MustCompile(url.QueryEscape(*cfg.CFClientSecret)),
				},
				redactString: *redactedTestUser.CFClientSecret,
			}
			regList = append(regList, reg)
		}
		interactionJson, err := json.Marshal(i)
		if err != nil {
			panic(err)
		}
		var casInteraction cassette.Interaction
		interactionFinal := string(interactionJson)
		for _, reg := range regList {
			for _, regPattern := range reg.regexpattern {
				interactionFinal = regPattern.ReplaceAllString(interactionFinal, reg.redactString)
			}
		}
		// `[A-Z]+[A-Za-z0-9_]` part added to enforce some capital letter to prevent cases like 12.1.2
		jwtTokenPattern := regexp.MustCompile(`([A-Za-z0-9_]*[A-Z]+[A-Za-z0-9_]*\.){2}[A-Za-z0-9-_]*`)
		interactionFinal = jwtTokenPattern.ReplaceAllString(interactionFinal, "redacted")
		err = json.Unmarshal([]byte(interactionFinal), &casInteraction)
		if err != nil {
			panic(err)
		}
		*i = casInteraction
		return nil
	}
}
func (cfg *CloudFoundryProviderConfigPtr) GetHookKind() recorder.HookKind {
	return recorder.BeforeSaveHook
}

func (cfg *CloudFoundryProviderConfigPtr) GetMatcher(t *testing.T, rec *recorder.Recorder) cassette.MatcherFunc {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Method != i.Method {
			return false
		}
		if r.URL.String() != i.URL {
			if !rec.IsRecording() {
				url, err := url.Parse(*redactedTestUser.Endpoint)
				if err != nil {
					panic(err)
				}
				r.URL = url
			} else {
				return false
			}
		}
		return true
	}
}
func (cfg *CloudFoundryProviderConfigPtr) SetupVCR(t *testing.T, cassetteName string) *recorder.Recorder {
	t.Helper()

	mode := recorder.ModeRecordOnce
	if force, _ := strconv.ParseBool(os.Getenv("TEST_FORCE_REC")); force {
		mode = recorder.ModeRecordOnly
	}

	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       cassetteName,
		Mode:               mode,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})

	if rec.IsRecording() {
		t.Logf("\nATTENTION: Recording '%s'", cassetteName)
		rec.AddHook(cfg.GetHook(), cfg.GetHookKind())
	} else {
		t.Logf("\nReplaying '%s'", cassetteName)
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		err = os.Setenv("CF_HOME", pwd+"/../../asset")
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		t.Fatal()
	}
	rec.SetMatcher(cfg.GetMatcher(t, rec))
	return rec
}

func strtostrptr(s string) *string {
	return &s
}

// Returns a pointer to a map
func maptomapptr(s map[string]string) *map[string]string {
	return &s
}

func getIdForImport(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}
