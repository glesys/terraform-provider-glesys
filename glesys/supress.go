package glesys

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// IgnoreCase check if the strings match when both are in lowercase
func IgnoreCase(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
