// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
)

// NewInstanceResource is a helper function to simplify the provider implementation.
func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

// instanceResource is the resource implementation.
type instanceResource struct {
	wxOneClients *WxOneClients
}

// Metadata returns the resource type name.
func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema defines the schema for the resource.
func (r *instanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the instance in UUID format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the instance.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the instance.",
				Computed:    true,
			},
			"flavor_id": schema.StringAttribute{
				Description: "Flavor of the instance.",
				Required:    true,
			},
			"subnet_id": schema.StringAttribute{
				Description: "Network subnet of the instance.",
				Required:    true,
			},
			"image_id": schema.StringAttribute{
				Description: "Image of the instance.",
				Required:    true,
			},
			"availability_zone": schema.StringAttribute{
				Description: "Availability zone of the instance.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("wx_dus_1"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "Project id of the instance.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_keys": schema.ListAttribute{
				Required:    true,
				Description: "Id of the subnet.",
				ElementType: types.StringType,
			},
			// "floating_ips": schema.ListNestedAttribute{
			// 	Required: true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"id": schema.StringAttribute{
			// 				Description: "Id of the subnet.",
			// 				Computed:    true,
			// 			},
			// 			"ip": schema.StringAttribute{
			// 				Description: "IP Version of the subnet.",
			// 				Computed:    true,
			// 			},
			// 			"nat_to_vms_private_ip": schema.BoolAttribute{
			// 				Description: "Is NAT enabled or not for this public IP.",
			// 				Computed:    true,
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

type instanceResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	FlavorID         types.String `tfsdk:"flavor_id"`
	NetworkID        types.String `tfsdk:"subnet_id"`
	ImageID          types.String `tfsdk:"image_id"`
	Status           types.String `tfsdk:"status"`
	ProjectID        types.String `tfsdk:"project_id"`
	// FloatingIPs      []floatingIPModel `tfsdk:"floating_ips"`
	SSHKeys []types.String `tfsdk:"ssh_keys"`
}

// type floatingIPModel struct {
// 	ID  types.String `tfsdk:"id"`
// 	IP  types.String `tfsdk:"ip"`
// 	Nat types.String `tfsdk:"nat_to_vms_private_ip"`
// }

// Configure adds the provider configured client to the resource.
func (r *instanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan instanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKeysInput := make([]string, len(plan.SSHKeys))
	for i, item := range plan.SSHKeys {
		sshKeysInput[i] = item.ValueString()
	}

	// Create new instance
	instance, err := createInstance(ctx, r.wxOneClients.graphqlClient, plan.NetworkID.ValueString(), plan.FlavorID.ValueString(), plan.ImageID.ValueString(), plan.ProjectID.ValueString(), plan.Name.ValueString(), sshKeysInput, AvailabilityZone(plan.AvailabilityZone.ValueString()), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating instance",
			"Could not create instance, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(instance.CreateInstance.Msg.Id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed instance value from WX-ONE
	instance, err := getInstance(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error Reading WX-ONE",
				"Could not read WX-ONE instance ID "+state.ID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	state.Name = types.StringValue(instance.GetInstance.Msg.Name)
	state.AvailabilityZone = types.StringValue((string(instance.GetInstance.Msg.AvailabilityZone)))

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Error Updating WX-ONE Instance",
		"An instance can not be updated",
	)
	return

	// var plan instanceResourceModel
	// diags := req.Plan.Get(ctx, &plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Update existing instance
	// _, err := updateInstance(ctx, r.wxOneClients.graphqlClient, plan.ID.ValueString(), plan.ProjectID.ValueString(), plan.Name.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Updating WX-ONE Instance",
	// 		"Could not update instance, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// diags = resp.State.Set(ctx, plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing instance
	_, err := deleteInstance(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WX-ONE Instance",
			"Could not delete instance, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *instanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
