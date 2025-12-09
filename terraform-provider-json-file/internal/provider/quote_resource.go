// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-json-file/internal/quote"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = &QuoteResource{}
var _ resource.ResourceWithIdentity = &QuoteResource{}

func NewQuoteResource() resource.Resource {
	return &QuoteResource{}
}

// QuoteResource defines the resource implementation.
type QuoteResource struct {
	folderPath string
}

// QuoteResourceModel describes the resource data model.
type QuoteResourceModel struct {
	Message types.String `tfsdk:"message"`
	Author  types.String `tfsdk:"author"`
	ID      types.String `tfsdk:"id"`
}

type QuoteResourceIdentityModel struct {
	ID types.String `tfsdk:"id"`
}

func (r *QuoteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quote"
}

// IdentitySchema implements resource.ResourceWithIdentity.
func (r *QuoteResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				Description:       "ID of the quote",
				RequiredForImport: true,
			},
		},
	}
}

func (r *QuoteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Quote",

		Attributes: map[string]schema.Attribute{
			"message": schema.StringAttribute{
				MarkdownDescription: "Message of the quote",
				Required:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "Who said the quote",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the quote",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *QuoteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	folderPath, ok := req.ProviderData.(string)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.folderPath = folderPath
}

func (r *QuoteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data QuoteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	q := quote.Quote{
		Message: data.Message.ValueString(),
		Author:  data.Author.ValueString(),
	}

	id, err := quote.CreateQuoteFile(r.folderPath, q)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to create quote", "failed to create quote: "+err.Error()))
		return
	}

	identityData := QuoteResourceIdentityModel{
		ID: types.StringValue(id),
	}
	data.ID = types.StringValue(id)

	tflog.Trace(ctx, "created quote "+id)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identityData)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuoteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read identity data
	var identityData QuoteResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := identityData.ID.ValueString()
	q, err := quote.ReadQuote(r.folderPath, id)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to read quote", "failed to read quote: "+err.Error()))
		return
	}

	data := QuoteResourceModel{
		Message: types.StringValue(q.Message),
		Author:  types.StringValue(q.Author),
		ID:      types.StringValue(id),
	}

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identityData)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuoteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data QuoteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read identity data
	var identityData QuoteResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		id = identityData.ID.ValueString()
		q  = quote.Quote{
			Message: data.Message.ValueString(),
			Author:  data.Author.ValueString(),
		}
	)
	if err := quote.WriteQuoteFile(r.folderPath, id, q); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to update quote", "failed to update quote: "+err.Error()))
		return
	}

	tflog.Trace(ctx, "updated quote "+id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuoteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read identity data
	var identityData QuoteResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := identityData.ID.ValueString()
	if err := quote.DeleteQuoteFile(r.folderPath, id); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to delete quote", "failed to delete quote: "+err.Error()))
		return
	}

	tflog.Trace(ctx, "deleted quote "+id)
}
