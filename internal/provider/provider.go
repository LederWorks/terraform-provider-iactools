// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure iactoolsProvider satisfies various provider interfaces.
var _ provider.Provider = &iactoolsProvider{}
var _ provider.ProviderWithFunctions = &iactoolsProvider{}
var _ provider.ProviderWithEphemeralResources = &iactoolsProvider{}

// iactoolsProvider defines the provider implementation.
type iactoolsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// iactoolsProviderModel describes the provider data model.
type iactoolsProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *iactoolsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "iactools"
	resp.Version = p.version
}

func (p *iactoolsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "LederWorks [iactools](https://github.com/lederworks/terraform-provider-iactools) provider. Requires Terraform 1.8 or later.",
				Optional:            true,
			},
		},
	}
}

func (p *iactoolsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data iactoolsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *iactoolsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *iactoolsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *iactoolsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *iactoolsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewInverseCIDRFunction,
		NewReverseDNSFunction,
	}
}

// New creates a new instance of the iactools provider with the specified version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &iactoolsProvider{
			version: version,
		}
	}
}
