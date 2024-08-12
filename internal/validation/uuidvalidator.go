package validation

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	UuidRegexp = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
	ShaRegexp  = regexp.MustCompile(`\b([a-f0-9]{64})\b`)
)

// ValidUUID checks that the String held in the attribute is a valid UUID.
func ValidUUID() validator.String {
	return stringvalidator.RegexMatches(UuidRegexp, "value must be a valid UUID")
}
