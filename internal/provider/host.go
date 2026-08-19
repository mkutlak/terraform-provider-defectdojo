package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefectDojo case-folds every host it stores. dojo/url/models.py clean_host()
// runs an IP through ipaddress.ip_address().compressed and a name through
// idna.encode(..., uts46=True); both answer in lower case, and the fallback for
// a name IDNA cannot encode is a bare host.lower(). Verified on 3.1.101:
// POST /api/v2/url/ {"host": "API.Example.COM"} returns 201 carrying
// "host": "api.example.com", and the GET readback agrees.
//
// url.host is Required, so the plan carries the configured literal and the
// create response contradicts it. Without preservation every apply fails with
//
//	.host: was cty.StringVal("API.Example.COM"),
//	       but now cty.StringVal("api.example.com")
//
// after the URL has already been created server-side. Terraform records that
// object as tainted, so the next apply destroys it, creates it again and fails
// identically: the resource never converges and no amount of retrying helps.
//
// Preserving the configured literal rather than rejecting upper case is the
// choice a236878 made for revenue, and hostnames make it an easy one: DNS names
// are case-insensitive (RFC 4343), so the two spellings denote the same host
// and there is nothing to warn the practitioner about.

// preserveHostCase picks the string to store in state for a server-provided
// host. When the attribute already holds a spelling that differs only by letter
// case - the practitioner's configured value on create/update, or the prior
// state value on refresh - that spelling is kept verbatim. Otherwise the
// server's own rendering is used.
//
// Only a case-insensitive match is preserved, so a genuinely different host is
// still reported as drift, and the other rewrites clean_host() can perform -
// IDNA punycoding ("Bücher.example" -> "xn--bcher-kva.example") and IPv6
// compression ("2001:db8:0:0:0:0:0:1" -> "2001:db8::1") - still reach state as
// the server spelled them. An import has nothing configured to preserve and so
// stores the server's spelling too, exactly like preserveDecimalLiteral and
// preserveDateTimeLiteral.
func preserveHostCase(current types.String, server string) types.String {
	if !current.IsNull() && !current.IsUnknown() && strings.EqualFold(current.ValueString(), server) {
		return current
	}

	return types.StringValue(server)
}
