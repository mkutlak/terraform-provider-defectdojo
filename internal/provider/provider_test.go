package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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

func testAccPreCheck(t *testing.T) {
	if os.Getenv("DEFECTDOJO_BASEURL") == "" {
		t.Fatal("DEFECTDOJO_BASEURL must be set for acceptance tests")
	}
	if os.Getenv("DEFECTDOJO_APIKEY") == "" &&
		(os.Getenv("DEFECTDOJO_USERNAME") == "" || os.Getenv("DEFECTDOJO_PASSWORD") == "") {
		t.Fatal("DEFECTDOJO_APIKEY or both DEFECTDOJO_USERNAME and DEFECTDOJO_PASSWORD must be set for acceptance tests")
	}
}
