// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true, // this particular test is fast and only relies on local resources (else, set TF_ACC=true)
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `provider "jsonfile" {
					folder_path = "/workspaces/go-tf-provider-lab/eeee"
				}

				resource "jsonfile_quote" "joke" {
					author = "adibou"
					message = "Coucou me revoilou"
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"jsonfile_quote.joke",
						tfjsonpath.New("author"),
						knownvalue.StringExact("adibou"),
					),
					statecheck.ExpectKnownValue(
						"jsonfile_quote.joke",
						tfjsonpath.New("message"),
						knownvalue.StringExact("Coucou me revoilou"),
					),
					statecheck.ExpectKnownValue(
						"jsonfile_quote.joke",
						tfjsonpath.New("id"),
						&uuidCheck{},
					),
				},
			},
			// Update and Read testing
			{
				Config: `provider "jsonfile" {
					folder_path = "/workspaces/go-tf-provider-lab/eeee"
				}

				resource "jsonfile_quote" "joke" {
					author = "adibou"
					message = "Oh non, bouzigouloum!"
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"jsonfile_quote.joke",
						tfjsonpath.New("author"),
						knownvalue.StringExact("adibou"),
					),
					statecheck.ExpectKnownValue(
						"jsonfile_quote.joke",
						tfjsonpath.New("message"),
						knownvalue.StringExact("Oh non, bouzigouloum!"),
					),
					statecheck.ExpectKnownValue(
						"jsonfile_quote.joke",
						tfjsonpath.New("id"),
						&uuidCheck{},
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

type uuidCheck struct{}

// CheckValue implements knownvalue.Check.
func (u *uuidCheck) CheckValue(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value should be a string, is %T", value)
	}
	if _, err := uuid.Parse(str); err != nil {
		return fmt.Errorf("%q is not a UUID", str)
	}
	return nil
}

// String implements knownvalue.Check.
func (u *uuidCheck) String() string {
	return "UUID format"
}
