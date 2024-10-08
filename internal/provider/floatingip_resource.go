// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &floatingIPResource{}
	_ resource.ResourceWithConfigure   = &floatingIPResource{}
	_ resource.ResourceWithImportState = &floatingIPResource{}
)

// NewFloatingIPResource is a helper function to simplify the provider implementation.
func NewFloatingIPResource() resource.Resource {
	return &floatingIPResource{}
}

// floatingIPResource is the resource implementation.
type floatingIPResource struct {
	wxOneClients *WxOneClients
}

// Metadata returns the resource type name.
func (r *floatingIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_floatingIP"
}

// Schema defines the schema for the resource.
func (r *floatingIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a floatingIP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the floatingIP in UUID format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": schema.StringAttribute{
				Description: "The floatingIP.",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "Project id of the floatingIP.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type floatingIPResourceModel struct {
	ID        types.String `tfsdk:"id"`
	IP        types.String `tfsdk:"ip"`
	ProjectID types.String `tfsdk:"project_id"`
}

// Configure adds the provider configured client to the resource.
func (r *floatingIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	wxOneClients, ok := req.ProviderData.(*WxOneClients)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *WxOneClients, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.wxOneClients = wxOneClients
}

// Create creates the resource and sets the initial Terraform state.
func (r *floatingIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan floatingIPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new floatingIP
	floatingIP, err := createFloatingIP(ctx, r.wxOneClients.graphqlClient, plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating floatingIP",
			"Could not create floatingIP, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(floatingIP.CreateFloatingIP.Msg.Id)
	plan.IP = types.StringValue(floatingIP.CreateFloatingIP.Msg.Ip)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *floatingIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state floatingIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed floatingIP value from WX-ONE
	floatingIP, err := getFloatingIP(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error Reading WX-ONE",
				"Could not read WX-ONE floatingIP ID "+state.ID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	state.IP = types.StringValue(floatingIP.GetFloatingIP.Msg.Ip)

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *floatingIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"This resource can not be updated",
		"",
	)
	return

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *floatingIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state floatingIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing floatingIP
	_, err := deleteFloatingIP(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WX-ONE FloatingIP",
			"Could not delete floatingIP, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *floatingIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format 'id:projectId', got '%s'", req.ID),
		)
		return
	}

	// Set each ID to the respective attributes in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[1])...)
}
