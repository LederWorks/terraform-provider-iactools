// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestInverseCidrFunction_Valid(t *testing.T) {
	testCases := map[string]struct {
		parentCIDR   string
		childCIDR    string
		inverseCIDRs []string
	}{
		"ipv4-example": {
			parentCIDR:   "192.168.0.0/24",
			childCIDR:    "192.168.0.0/26",
			inverseCIDRs: []string{"192.168.64.0/26", "192.168.128.0/25"},
		},
		"ipv6-example": {
			parentCIDR:   "2001:db8::/48",
			childCIDR:    "2001:db8::/50",
			inverseCIDRs: []string{"2001:db8:0:0:4000::/50", "2001:db8:0:0:8000::/49"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_8_0),
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							output "result" {
								value = provider::iactools::inverse_cidr("%s", "%s")
							}
						`, testCase.parentCIDR, testCase.childCIDR),
						Check: resource.TestCheckOutput("result", strings.Join(testCase.inverseCIDRs, ",")),
					},
				},
			})
		})
	}
}

func TestInverseCidrFunction_Invalid(t *testing.T) {
	testCases := map[string]struct {
		parentCIDR string
		childCIDR  string
		error      string
	}{
		"empty-parent-cidr": {
			parentCIDR: "",
			childCIDR:  "192.168.1.0/24",
			error:      "The parent_cidr argument must be provided and valid.",
		},
		"empty-child-cidr": {
			parentCIDR: "192.168.0.0/16",
			childCIDR:  "",
			error:      "The child_cidr argument must be provided and valid.",
		},
		"invalid-parent-cidr": {
			parentCIDR: "invalid-cidr",
			childCIDR:  "192.168.1.0/24",
			error:      "Error calculating inverse CIDRs: invalid CIDR format",
		},
		"invalid-child-cidr": {
			parentCIDR: "192.168.0.0/16",
			childCIDR:  "invalid-cidr",
			error:      "Error calculating inverse CIDRs: invalid CIDR format",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_8_0),
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							output "result" {
								value = provider::iactools::inverse_cidr("%s", "%s")
							}
						`, testCase.parentCIDR, testCase.childCIDR),
						ExpectError: regexp.MustCompile(testCase.error),
					},
				},
			})
		})
	}
}
