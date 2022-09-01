package glesys

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IgnoreCase check if the strings match when both are in lowercase
func IgnoreCase(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
