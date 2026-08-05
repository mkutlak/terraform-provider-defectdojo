package provider

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"defectdojo": func() (tfprotov6.ProviderServer, error) {
		provider := New("test")()
		server, err := providerserver.NewProtocol6WithError(provider)()
		return server, err
	},
}

// testAccConfig holds every environment-derived knob the acceptance suite
// needs, resolved once in TestMain instead of via scattered os.Getenv calls.
// It is only populated when TF_ACC is set; unit tests never read it.
type testAccConfig struct {
	baseURL   string
	apiKey    string
	username  string
	password  string
	ddVersion string
}

var testAccConf testAccConfig

// TestMain short-circuits when TF_ACC is unset so that the unit tests (the
// ddField audits, the schema meta-tests, the datetime tables) stay instant and
// need no backend. Without this, every `go test` invocation would pay for
// acceptance setup even though resource.Test would skip immediately.
func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC") == "" {
		os.Exit(m.Run())
	}

	testAccConf = testAccConfig{
		baseURL:   os.Getenv("DEFECTDOJO_BASEURL"),
		apiKey:    os.Getenv("DEFECTDOJO_APIKEY"),
		username:  os.Getenv("DEFECTDOJO_USERNAME"),
		password:  os.Getenv("DEFECTDOJO_PASSWORD"),
		ddVersion: os.Getenv("DD_VERSION"),
	}

	os.Exit(m.Run())
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("DEFECTDOJO_BASEURL") == "" {
		t.Fatal("DEFECTDOJO_BASEURL must be set for acceptance tests")
	}
	if os.Getenv("DEFECTDOJO_APIKEY") == "" &&
		(os.Getenv("DEFECTDOJO_USERNAME") == "" || os.Getenv("DEFECTDOJO_PASSWORD") == "") {
		t.Fatal("DEFECTDOJO_APIKEY or both DEFECTDOJO_USERNAME and DEFECTDOJO_PASSWORD must be set for acceptance tests")
	}
}

// testAccClient builds a DefectDojo client from the resolved test config. It is
// the single construction path for tests that need to talk to the API directly
// (seeding fixtures, verifying destroys, simulating out-of-band drift).
func testAccClient(ctx context.Context) (*dd.ClientWithResponses, error) {
	return newClient(ctx, testAccConf.baseURL, testAccConf.apiKey, testAccConf.username, testAccConf.password)
}

// skipIfDDBelow skips the test when the DefectDojo backend under test is older
// than min (e.g. "3.2"). tfversion gates the Terraform CLI version; there is no
// upstream equivalent for the backend, so capability gating lives here.
//
// When DD_VERSION is unset the test runs: an unknown backend is assumed current
// rather than silently skipping coverage.
func skipIfDDBelow(t *testing.T, min string) {
	t.Helper()

	if testAccConf.ddVersion == "" {
		return
	}
	if compareVersions(testAccConf.ddVersion, min) < 0 {
		t.Skipf("requires DefectDojo >= %s, running against %s", min, testAccConf.ddVersion)
	}
}

// compareVersions compares dotted numeric versions, returning -1, 0 or 1.
// Missing trailing components count as zero, so "3.1" == "3.1.0". Non-numeric
// components compare as zero rather than erroring: version gating must never be
// the reason a suite fails.
func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < max(len(aParts), len(bParts)); i++ {
		var av, bv int
		if i < len(aParts) {
			av, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bv, _ = strconv.Atoi(bParts[i])
		}
		if av != bv {
			if av < bv {
				return -1
			}
			return 1
		}
	}
	return 0
}
