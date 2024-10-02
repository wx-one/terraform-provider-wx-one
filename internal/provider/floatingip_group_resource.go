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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &floatingIPGroupResource{}
	_ resource.ResourceWithConfigure   = &floatingIPGroupResource{}
	_ resource.ResourceWithImportState = &floatingIPGroupResource{}
)

// NewFloatingIPGroupResource is a helper function to simplify the provider implementation.
func NewFloatingIPGroupResource() resource.Resource {
	return &floatingIPGroupResource{}
}

// floatingIPGroupResource is the resource implementation.
type floatingIPGroupResource struct {
	wxOneClients *WxOneClients
}

// Metadata returns the resource type name.
func (r *floatingIPGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_floatingIPGroup"
}

// Schema defines the schema for the resource.
func (r *floatingIPGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a floatingIPGroup.",
		Attributes: map[string]schema.Attribute{

			"floatingip_id": schema.StringAttribute{
				Description: "ID of the floatingIP in UUID format.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "ID of the project in UUID format.",
				Required:    true,
			},
			"nat": schema.BoolAttribute{
				Description: "If the trafficn needs to be natted instead of routing towards the vm.",
				Optional:    true,
			},
			"vms": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the assignment in UUID format.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"priority": schema.Int32Attribute{
							Description: "Priority of the vm, yet unused.",
							Required:    true,
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.RequiresReplace(),
							},
						},
						"vm_id": schema.StringAttribute{
							Description: "ID of the VM in UUID format.",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

type floatingIPGroupResourceModel struct {
	FloatingIPID types.String `tfsdk:"floatingip_id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Nat          types.Bool   `tfsdk:"nat"`
	VMs          []vmModel    `tfsdk:"vms"`
}

type vmModel struct {
	ID       types.String `tfsdk:"id"`
	Priority types.Int32  `tfsdk:"priority"`
	VMID     types.String `tfsdk:"vm_id"`
}

// Configure adds the provider configured client to the resource.
func (r *floatingIPGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *floatingIPGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan floatingIPGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmInput := make([]FloatingGroupVmInput, len(plan.VMs))
	for i, item := range plan.VMs {
		vmInput[i] = FloatingGroupVmInput{
			Priority: int(item.Priority.ValueInt32()),
			Vm:       item.VMID.ValueString(),
		}
	}

	// Create new floatingIPGroup
	floatingIPGroup, err := createFloatingGroup(ctx, r.wxOneClients.graphqlClient, plan.FloatingIPID.ValueString(), plan.ProjectID.ValueString(), vmInput, plan.Nat.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating floatingIPGroup",
			"Could not create floatingIPGroup, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	for i, item := range floatingIPGroup.CreateFloatingGroup.Msg {
		plan.VMs[i].ID = types.StringValue(item.Id)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *floatingIPGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state floatingIPGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed floatingIPGroup value from WX-ONE
	floatingIPAttachment, err := getFloatingIPAttachment(ctx, r.wxOneClients.graphqlClient, state.FloatingIPID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error Reading WX-ONE",
				"Could not read WX-ONE floatingIPGroup ID "+state.FloatingIPID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	if floatingIPAttachment.GetFloatingIPAttachment.Msg.GetTypename() != "FloatingIPInstanceAttachment" {
		resp.Diagnostics.AddError(
			"Error IP manipulated outside of terraform",
			"Type is not FloatingIPInstanceAttachment, but "+floatingIPAttachment.GetFloatingIPAttachment.Msg.GetTypename()+": "+err.Error(),
		)
		return
	}

	// refresh the vms in the asset
	attachments := floatingIPAttachment.GetFloatingIPAttachment.Msg.(*getFloatingIPAttachmentGetFloatingIPAttachmentFloatingIPAttachmentResponseMsgFloatingIPInstanceAttachment)
	state.Nat = types.BoolValue(attachments.NatToVmsPrivateIp)
	state.VMs = make([]vmModel, len(attachments.Vms))
	for i, item := range attachments.Vms {
		state.VMs[i].ID = types.StringValue(item.Id)
		state.VMs[i].Priority = types.Int32Value(int32(item.Priority))
		state.VMs[i].VMID = types.StringValue(item.VmId)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *floatingIPGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"This resource can not yet be updated",
		"",
	)
	return

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *floatingIPGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state floatingIPGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing floatingIPGroup
	_, err := deleteFloatingGroupByFloatingIpIdAndInstanceId(ctx, r.wxOneClients.graphqlClient, state.ProjectID.ValueString(), state.FloatingIPID.ValueString(), state.VMs[0].VMID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WX-ONE FloatingIPGroup",
			"Could not delete floatingIPGroup, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *floatingIPGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
