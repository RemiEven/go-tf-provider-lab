// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Citation2000Provider defines the provider implementation.
type Citation2000Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Citation2000ProviderModel describes the provider data model.
type Citation2000ProviderModel struct {
	FolderPath types.String `tfsdk:"folder_path"`
}

func (p *Citation2000Provider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "citation2000"
	resp.Version = p.version
}

func (p *Citation2000Provider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"folder_path": schema.StringAttribute{
				MarkdownDescription: "Path of the folder containing the json files",
				Required:            true,
			},
		},
	}
}

func (p *Citation2000Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Citation2000ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = data.FolderPath.ValueString()
}

func (p *Citation2000Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQuoteResource,
	}
}

func (p *Citation2000Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Citation2000Provider{
			version: version,
		}
	}
}
