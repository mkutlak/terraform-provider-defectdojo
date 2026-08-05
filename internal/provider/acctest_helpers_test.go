package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
)

// uniqueId returns a unique identifier for use in acceptance tests.
func uniqueId() string {
	return id.UniqueId()
}
