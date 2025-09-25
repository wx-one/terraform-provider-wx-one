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
	_ resource.Resource                = &securityGroupResource{}
	_ resource.ResourceWithConfigure   = &securityGroupResource{}
	_ resource.ResourceWithImportState = &securityGroupResource{}
)

// NewsecurityGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupResource() resource.Resource {
	return &securityGroupResource{}
}

// securityGroupResource is the resource implementation.
type securityGroupResource struct {
	wxOneClients *WxOneClients
}

// Metadata returns the resource type name.
func (r *securityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_securityGroup"
}

// Schema defines the schema for the resource.
func (r *securityGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a securityGroup.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the securityGroup in UUID format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the securityGroup.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the securityGroup.",
				Optional:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "Project id of the securityGroup.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the subnet.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Description: "Name of the rule.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the rule.",
							Required:    true,
						},
						"direction": schema.StringAttribute{
							Description: "Direciton of the rule.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("ingress", "egress"),
							},
						},
						"ether_type": schema.StringAttribute{
							Description: "IP Version of the rule.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("IPv4", "IPv6"),
							},
						},
						"priority": schema.Int32Attribute{
							Description: "Priority of the rule.",
							Required:    true,
						},
						"ports": schema.StringAttribute{
							Description: "Port List, either single, all, comma separated or a range 1:256 of the rule.",
							Required:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "Protocol of the rule.",
							Required:    true,
						},
						"protocol_int": schema.Int32Attribute{
							Description: "Protocol Number of the rule if protocol is set to int.",
							Optional:    true,
						},
						"cidr": schema.StringAttribute{
							Description: "CIDR of the rule.",
							Required:    true,
						},
						"action": schema.StringAttribute{
							Description: "Action of the rule.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("accept", "drop", "reject"),
							},
						},
					},
				},
			},
		},
	}
}

type securityGroupResourceModel struct {
	ID          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	Description types.String             `tfsdk:"description"`
	ProjectID   types.String             `tfsdk:"project_id"`
	Rules       []securityGroupRuleModel `tfsdk:"rules"`
}

type securityGroupRuleModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Direction   types.String `tfsdk:"direction"`
	Priority    types.Int32  `tfsdk:"priority"`
	Protocol    types.String `tfsdk:"protocol"`
	ProtocolInt types.Int32  `tfsdk:"protocol_int"`
	EtherType   types.String `tfsdk:"ether_type"`
	Ports       types.String `tfsdk:"ports"`
	Cidr        types.String `tfsdk:"cidr"`
	Action      types.String `tfsdk:"action"`
}

// Configure adds the provider configured client to the resource.
func (r *securityGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan securityGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleInput := make([]*W1SecurityGroupRuleInput, len(plan.Rules))
	for i, item := range plan.Rules {
		protocolInt := int(item.ProtocolInt.ValueInt32())
		action := SecurityGroupRuleAction(item.Action.ValueString())

		ruleInput[i] = &W1SecurityGroupRuleInput{
			Name:        item.Name.ValueString(),
			Description: item.Description.ValueString(),
			Ports:       item.Ports.ValueString(),
			Cidr:        item.Cidr.ValueString(),
			Direction:   SecurityGroupRuleDirection(item.Direction.ValueString()),
			Priority:    int(item.Priority.ValueInt32()),
			EtherType:   EtherType(item.EtherType.ValueString()),
			Protocol:    Protocol(item.Protocol.ValueString()),
			ProtocolInt: &protocolInt,
			Action:      &action,
		}
	}

	// Create new securityGroup
	securityGroup, err := createSecurityGroup(ctx, r.wxOneClients.graphqlClient, plan.Name.ValueString(), plan.ProjectID.ValueString(), plan.Description.ValueStringPointer(), ruleInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating securityGroup",
			"Could not create securityGroup, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(securityGroup.CreateSecurityGroup.Msg.Id)

	for i, item := range securityGroup.CreateSecurityGroup.Msg.Rules {
		plan.Rules[i].ID = types.StringValue(item.Id)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state securityGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed securityGroup value from WX-ONE
	// securityGroup, err := getsecurityGroup(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	// if err != nil {
	// 	if isNotFoundError(err) {
	// 		resp.State.RemoveResource(ctx)
	// 		return
	// 	} else {
	// 		resp.Diagnostics.AddError(
	// 			"Error Reading WX-ONE",
	// 			"Could not read WX-ONE securityGroup ID "+state.ID.ValueString()+": "+err.Error(),
	// 		)
	// 		return
	// 	}
	// }

	// state.Name = types.StringValue(securityGroup.GetsecurityGroup.Msg.Name)
	// subnets := make([]subnetModel, len(securityGroup.GetsecurityGroup.Msg.Subnets))

	// for i, subnet := range securityGroup.GetsecurityGroup.Msg.Subnets {
	// 	subnets[i] = subnetModel{
	// 		ID:        types.StringValue(subnet.Id),
	// 		Name:      types.StringValue(subnet.Name),
	// 		IPVersion: types.StringValue(subnet.IpVersion),
	// 		CIDR:      types.StringValue(subnet.Cidr),
	// 	}
	// }
	// state.Subnets = subnets
	// state.AvailabilityZone = types.StringValue((string(securityGroup.GetsecurityGroup.Msg.AvailabilityZone)))

	// // Set refreshed state
	// diags = resp.State.Set(ctx, state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan securityGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleInput := make([]*W1SecurityGroupRuleInput, len(plan.Rules))
	for i, item := range plan.Rules {
		protocolInt := int(item.ProtocolInt.ValueInt32())
		action := SecurityGroupRuleAction(item.Action.ValueString())

		ruleInput[i] = &W1SecurityGroupRuleInput{
			Name:        item.Name.ValueString(),
			Description: item.Description.ValueString(),
			Ports:       item.Ports.ValueString(),
			Cidr:        item.Cidr.ValueString(),
			Direction:   SecurityGroupRuleDirection(item.Direction.ValueString()),
			Priority:    int(item.Priority.ValueInt32()),
			EtherType:   EtherType(item.EtherType.ValueString()),
			Protocol:    Protocol(item.Protocol.ValueString()),
			ProtocolInt: &protocolInt,
			Action:      &action,
		}
	}

	// Update existing securityGroup
	_, err := setSecurityGroup(ctx, r.wxOneClients.graphqlClient, plan.ID.ValueString(), plan.Name.ValueString(), plan.ProjectID.ValueString(), plan.Description.ValueStringPointer(), ruleInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating WX-ONE securityGroup",
			"Could not update securityGroup, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state securityGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing securityGroup
	// _, err := deletesecurityGroup(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), state.ProjectID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Deleting WX-ONE securityGroup",
	// 		"Could not delete securityGroup, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}

func (r *securityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
