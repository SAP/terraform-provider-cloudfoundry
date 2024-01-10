package provider

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
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
	}

	if err != nil {
		t.Fatal()
	}
	rec.SetMatcher(cfg.GetMatcher(t, rec))
	return rec
}
