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
		"ipv4-example-class-a": {
			parentCIDR:   "10.103.0.0/18",
			childCIDR:    "10.103.1.0/24",
			inverseCIDRs: []string{"10.103.32.0/19", "10.103.16.0/20", "10.103.8.0/21", "10.103.4.0/22", "10.103.2.0/23", "10.103.0.0/24"},
		},
		"ipv4-example-class-c": {
			parentCIDR:   "192.168.0.0/22",
			childCIDR:    "192.168.3.96/27",
			inverseCIDRs: []string{"192.168.0.0/23", "192.168.2.0/24", "192.168.3.128/25", "192.168.3.0/26", "192.168.3.64/27"},
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
                            output "result" {
                                value = provider::iactools::inverse_cidr("%s", "%s")
                            }
                        `, testCase.parentCIDR, testCase.childCIDR),
						Check: resource.ComposeAggregateTestCheckFunc(
							func(state *terraform.State) error {
								outputRaw, ok := state.RootModule().Outputs["result"]
								if !ok {
									return fmt.Errorf("output 'result' not found")
								}

								output, ok := outputRaw.Value.([]interface{})
								if !ok {
									return fmt.Errorf("expected list output, got %T", outputRaw.Value)
								}

								// Convert the list of interfaces to a list of strings
								var result []string
								for _, v := range output {
									str, ok := v.(string)
									if !ok {
										return fmt.Errorf("expected string in list, got %T", v)
									}
									result = append(result, str)
								}

								// Check that the output matches the expected result
								if len(result) != len(testCase.inverseCIDRs) {
									return fmt.Errorf("expected %d elements, got %d", len(testCase.inverseCIDRs), len(result))
								}

								for i, cidr := range testCase.inverseCIDRs {
									if result[i] != cidr {
										return fmt.Errorf("expected %s at position %d, got %s", cidr, i, result[i])
									}
								}

								return nil
							},
						),
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
			error:      `(?s)Call to function "provider::iactools::inverse_cidr" failed.*The parent_cidr.*argument must be provided and valid`,
		},
		"empty-child-cidr": {
			parentCIDR: "192.168.0.0/16",
			childCIDR:  "",
			error:      `(?s)Call to function "provider::iactools::inverse_cidr" failed.*The child_cidr.*argument must be provided and valid`,
		},
		"invalid-parent-cidr": {
			parentCIDR: "invalid-cidr",
			childCIDR:  "192.168.1.0/24",
			error:      `(?s)Call to function "provider::iactools::inverse_cidr" failed.*invalid parent CIDR: invalid CIDR address: invalid-cidr`,
		},
		"invalid-child-cidr": {
			parentCIDR: "192.168.0.0/16",
			childCIDR:  "invalid-cidr",
			error:      `(?s)Call to function "provider::iactools::inverse_cidr" failed.*invalid child CIDR: invalid CIDR address: invalid-cidr`,
		},
		"parentless-child-cidr": {
			parentCIDR: "192.168.0.0/16",
			childCIDR:  "172.16.0.0/24",
			error:      `(?s)Call to function "provider::iactools::inverse_cidr" failed.*Error calculating.*inverse CIDRs: child CIDR 172.16.0.0/24 is not within parent CIDR.*192.168.0.0/16`,
		},
		"mask-too-large": {
			parentCIDR: "192.168.84.42/32",
			childCIDR:  "192.168.84.42/32",
			error:      `(?s)Call to function "provider::iactools::inverse_cidr" failed.*Error calculating.*inverse CIDRs: child CIDR 172.16.0.0/24 is not within parent CIDR.*192.168.0.0/16`,
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
