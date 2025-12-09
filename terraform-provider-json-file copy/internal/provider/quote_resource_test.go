// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"jsonfile_quote.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"jsonfile_quote.test",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("example value when not configured"),
					),
					statecheck.ExpectKnownValue(
						"jsonfile_quote.test",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("one"),
					),
				},
			},
			// // Update and Read testing
			// {
			// 	Config: testAccExampleResourceConfig("two"),
			// 	ConfigStateChecks: []statecheck.StateCheck{
			// 		statecheck.ExpectKnownValue(
			// 			"jsonfile_quote.test",
			// 			tfjsonpath.New("id"),
			// 			knownvalue.StringExact("example-id"),
			// 		),
			// 		statecheck.ExpectKnownValue(
			// 			"jsonfile_quote.test",
			// 			tfjsonpath.New("defaulted"),
			// 			knownvalue.StringExact("example value when not configured"),
			// 		),
			// 		statecheck.ExpectKnownValue(
			// 			"jsonfile_quote.test",
			// 			tfjsonpath.New("configurable_attribute"),
			// 			knownvalue.StringExact("two"),
			// 		),
			// 	},
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "jsonfile" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}
