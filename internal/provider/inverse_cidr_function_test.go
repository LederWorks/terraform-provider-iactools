// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestInverseCidrFunction_Valid(t *testing.T) {
	testCases := map[string]struct {
		parentCIDR   string
		childCIDR    string
		inverseCIDRs []string
	}{
		"ipv4-example": {
			parentCIDR:   "10.103.0.0/18",
			childCIDR:    "10.103.1.0/24",
			inverseCIDRs: []string{"10.103.32.0/19", "10.103.16.0/20", "10.103.8.0/21", "10.103.4.0/22", "10.103.2.0/23", "10.103.0.0/24"},
		},
		"ipv6-example": {
			parentCIDR:   "2001:db8::/44",
			childCIDR:    "2001:db8::/50",
			inverseCIDRs: []string{"2001:db8:8::/45", "2001:db8:4::/46", "2001:db8:2::/47", "2001:db8:1::/48", "2001:db8:0:8000::/49", "2001:db8:0:4000::/50"},
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
							data "iactools_inverse_cidr" "test" {  
								parent_cidr = "%s"  
								child_cidr  = "%s"  
							}  
  
							output "result" {  
								value = data.iactools_inverse_cidr.test.inverse_cidrs  
							}  
						`, testCase.parentCIDR, testCase.childCIDR),
						Check: func(s *terraform.State) error {
							rs := s.RootModule().Outputs["result"]
							if rs == nil {
								return fmt.Errorf("output 'result' not found")
							}
							got := rs.Value.([]interface{})
							if len(got) != len(testCase.inverseCIDRs) {
								return fmt.Errorf("expected %d inverse CIDRs, got %d", len(testCase.inverseCIDRs), len(got))
							}
							for i, v := range testCase.inverseCIDRs {
								if got[i].(string) != v {
									return fmt.Errorf("expected %s at index %d, got %s", v, i, got[i].(string))
								}
							}
							return nil
						},
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
