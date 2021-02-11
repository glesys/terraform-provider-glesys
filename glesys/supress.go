package glesys

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func IgnoreCase(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
