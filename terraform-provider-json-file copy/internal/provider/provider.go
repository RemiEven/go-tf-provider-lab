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

// JsonFileProvider defines the provider implementation.
type JsonFileProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// JsonFileProviderModel describes the provider data model.
type JsonFileProviderModel struct {
	FolderPath types.String `tfsdk:"folder_path"`
}

func (p *JsonFileProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jsonfile"
	resp.Version = p.version
}

func (p *JsonFileProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"folder_path": schema.StringAttribute{
				MarkdownDescription: "Path of the folder containing the json files",
				Required:            true,
			},
		},
	}
}

func (p *JsonFileProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JsonFileProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = data.FolderPath.ValueString()
}

func (p *JsonFileProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQuoteResource,
	}
}

func (p *JsonFileProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JsonFileProvider{
			version: version,
		}
	}
}
