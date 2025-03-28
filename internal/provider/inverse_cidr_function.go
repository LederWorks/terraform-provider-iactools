// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ function.Function = InverseCIDRFunction{}
)

// NewInverseCIDRFunction is a helper function to create a new instance of InverseCIDRFunction.
func NewInverseCIDRFunction() function.Function {
	return InverseCIDRFunction{}
}

// InverseCIDRFunction is the struct for the inverse CIDR function.
type InverseCIDRFunction struct{}

// Metadata sets the metadata for the function.
func (r InverseCIDRFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "inverse_cidrs"
}

// Definition sets the definition for the function.
func (r InverseCIDRFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Calculate the inverse CIDR ranges of a parent and a child CIDR",
		MarkdownDescription: "Accepts both IPv4 and IPv6 addresses and outputs their inverse CIDR ranges.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "parent_cidr",
				MarkdownDescription: "The CIDR of the parent network",
			},
			function.StringParameter{
				Name:                "child_cidr",
				MarkdownDescription: "The CIDR of the child network",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

// Run executes the inverse CIDR function.
func (r InverseCIDRFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var parentCIDR, childCIDR string

	// Parse the arguments
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &parentCIDR, &childCIDR))
	if resp.Error != nil {
		return
	}

	// Validate input arguments
	if parentCIDR == "" {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("The parent_cidr argument must be provided and valid"))
		return
	}
	if childCIDR == "" {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("The child_cidr argument must be provided and valid"))
		return
	}

	// Calculate inverse CIDRs
	inverseCIDRs, err := InverseCIDR(parentCIDR, childCIDR)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("Error calculating inverse CIDRs: %s", err.Error())))
		return
	}

	// Convert the result to a Terraform-compatible type
	inverseCIDRList := make([]attr.Value, len(inverseCIDRs))
	for i, cidr := range inverseCIDRs {
		inverseCIDRList[i] = types.StringValue(cidr)
	}

	// Set the result
	listValue, diags := types.ListValue(types.StringType, inverseCIDRList)
	if diags.HasError() {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.FuncErrorFromDiags(ctx, diags))
		return
	}
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, listValue))
}
