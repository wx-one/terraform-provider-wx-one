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
	_ resource.Resource                = &securityGroupAssociationResource{}
	_ resource.ResourceWithConfigure   = &securityGroupAssociationResource{}
	_ resource.ResourceWithImportState = &securityGroupAssociationResource{}
)

// NewFloatingIPGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupAssociationResource() resource.Resource {
	return &securityGroupAssociationResource{}
}

// securityGroupAssociationResource is the resource implementation.
type securityGroupAssociationResource struct {
	wxOneClients *WxOneClients
}

// Metadata returns the resource type name.
func (r *securityGroupAssociationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_securityGroupAssociation"
}

// Schema defines the schema for the resource.
func (r *securityGroupAssociationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a securityGroupAssociation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the securityGroupAssociation in UUID format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Description: "ID of the security group in UUID format.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "ID of the project in UUID format.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.StringAttribute{
				Description: "ID of the vm in UUID format.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Description: "ID of the subnet in UUID format (not yet supported)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type securityGroupAssociationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`
	ProjectID       types.String `tfsdk:"project_id"`
	VmID            types.String `tfsdk:"vm_id"`
	SubnetID        types.String `tfsdk:"subnet_id"`
}

// Configure adds the provider configured client to the resource.
func (r *securityGroupAssociationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *securityGroupAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan securityGroupAssociationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new securityGroupAssociation
	securityGroupAssociation, err := assignSecurityGroup(ctx, r.wxOneClients.graphqlClient, plan.SecurityGroupID.ValueString(), plan.ProjectID.ValueString(), plan.VmID.ValueStringPointer(), plan.SubnetID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating securityGroupAssociation",
			"Could not create securityGroupAssociation, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(securityGroupAssociation.CreateSecurityGroupAssociation.Msg.Id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securityGroupAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state securityGroupAssociationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityGroupAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"This resource can not yet be updated",
		"",
	)
	return

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityGroupAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state securityGroupAssociationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing securityGroupAssociation
	_, err := detachSecurityGroup(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WX-ONE FloatingIPGroup",
			"Could not delete securityGroupAssociation, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *securityGroupAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
