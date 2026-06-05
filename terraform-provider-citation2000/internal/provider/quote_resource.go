// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/remieven/citation2000"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.ResourceWithConfigure = &QuoteResource{}
)

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

func (r *QuoteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quote"
}

func (r *QuoteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

	q := citation2000.Quote{
		Message: data.Message.ValueString(),
		Author:  data.Author.ValueString(),
	}

	id, err := citation2000.CreateQuoteFile(r.folderPath, q)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to create quote", "failed to create quote: "+err.Error()))
		return
	}

	data.ID = types.StringValue(id)

	tflog.Trace(ctx, "created quote "+id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuoteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data QuoteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	q, err := citation2000.ReadQuote(r.folderPath, id)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to read quote", "failed to read quote: "+err.Error()))
		return
	} else if q == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data = QuoteResourceModel{
		Message: types.StringValue(q.Message),
		Author:  types.StringValue(q.Author),
		ID:      types.StringValue(id),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuoteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data QuoteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		id = data.ID.ValueString()
		q  = citation2000.Quote{
			Message: data.Message.ValueString(),
			Author:  data.Author.ValueString(),
		}
	)
	if err := citation2000.WriteQuoteFile(r.folderPath, id, q); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to update quote", "failed to update quote: "+err.Error()))
		return
	}

	tflog.Trace(ctx, "updated quote "+id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuoteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data QuoteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	if err := citation2000.DeleteQuoteFile(r.folderPath, id); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to delete quote", "failed to delete quote: "+err.Error()))
		return
	}

	tflog.Trace(ctx, "deleted quote "+id)
}
