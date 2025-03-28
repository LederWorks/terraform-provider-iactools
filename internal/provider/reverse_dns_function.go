// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ function.Function = ReverseDNSFunction{}
)

// NewReverseDNSFunction is a helper function to create a new instance of ReverseDNSFunction.
func NewReverseDNSFunction() function.Function {
	return ReverseDNSFunction{}
}

// ReverseDNSFunction is the struct for the reverse DNS function.
type ReverseDNSFunction struct{}

// Metadata sets the metadata for the function.
func (f ReverseDNSFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "reverse_dns"
}

// Definition sets the definition for the function.
func (f ReverseDNSFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Calculate the reverse DNS name of an IP address",
		MarkdownDescription: "Accepts both IPv4 and IPv6 addresses and outputs their reverse DNS entry.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "ip_address",
				MarkdownDescription: "The IPv4 or IPv6 address itself",
			},
		},
		Return: function.StringReturn{},
	}
}

// Run executes the reverse DNS function.
func (f ReverseDNSFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var ipAddress string

	// Parse the argument
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &ipAddress))
	if resp.Error != nil {
		return
	}

	// Validate input argument
	if ipAddress == "" {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("The ip_address argument must be provided and valid"))
		return
	}

	// Parse the IP address
	parsedIP := net.ParseIP(ipAddress)
	if parsedIP == nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("Cannot parse IP address '%s'", ipAddress)))
		return
	}

	// Calculate reverse DNS
	var result string
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		result = ReverseDNSIPv4(ipv4.String())
	} else {
		result = ReverseDNSIPv6(parsedIP)
	}

	// Convert the result to a Terraform-compatible type
	stringValue := types.StringValue(result)

	// Set the result
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, stringValue))
}
