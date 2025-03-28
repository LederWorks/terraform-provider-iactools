// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package acceptance_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestInverseCidrFunction(t *testing.T) {
	testCases := map[string]struct {
		parentCIDR   string
		childCIDR    string
		inverseCIDRs []string
	}{
		"ipv4-example": {
			parentCIDR:   "192.168.0.0/16",
			childCIDR:    "192.168.1.0/24",
			inverseCIDRs: []string{"192.168.0.0/24", "192.168.2.0/23"},
		},
		"ipv6-example": {
			parentCIDR:   "2001:db8::/32",
			childCIDR:    "2001:db8:1::/48",
			inverseCIDRs: []string{"2001:db8:0::/48", "2001:db8:2::/47"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../functions/inverse_cidr",
				Vars: map[string]interface{}{
					"parent_cidr": testCase.parentCIDR,
					"child_cidr":  testCase.childCIDR,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			expectedOutput := strings.Join(testCase.inverseCIDRs, ",")
			actualOutput := terraform.Output(t, terraformOptions, "inverse_cidr")
			assert.Equal(t, expectedOutput, actualOutput, "inverse_cidr")
		})
	}
}
