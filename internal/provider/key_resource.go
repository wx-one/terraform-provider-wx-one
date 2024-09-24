package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &keyResource{}
	_ resource.ResourceWithConfigure = &keyResource{}
)

// NewKeyResource is a helper function to simplify the provider implementation.
func NewKeyResource() resource.Resource {
	return &keyResource{}
}

// keyResource is the resource implementation.
type keyResource struct {
	wxOneClients *WxOneClients
}

// Metadata returns the resource type name.
func (r *keyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

// Schema defines the schema for the resource.
func (r *keyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"private_key": schema.StringAttribute{
				Optional: true,
			},
			"public_key": schema.StringAttribute{
				Optional: true,
			},
			"project_wide": schema.BoolAttribute{
				Optional: true,
			},
			"project_id": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type keyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PrivateKey  types.String `tfsdk:"private_key"`
	PublicKey   types.String `tfsdk:"public_key"`
	ProjectWide types.Bool   `tfsdk:"project_wide"`
	ProjectId   types.String `tfsdk:"project_id"`
}

// Configure adds the provider configured client to the resource.
func (r *keyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan keyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new key
	key, err := createKey(ctx, r.wxOneClients.graphqlClient, plan.Name.ValueString(), plan.PublicKey.ValueString(), plan.ProjectId.ValueString(), plan.ProjectWide.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating key",
			"Could not create key, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(key.CreateKey.Msg.Id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *keyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state keyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed key value from WX-ONE
	key, err := getKey(ctx, r.wxOneClients.graphqlClient, state.ID.ValueString(), "", "", (*bool)(nil))
	if err != nil {

		if errList, ok := err.(gqlerror.List); ok {
			tflog.Info(ctx, "####### if before parsing")

			gqlerr := &gqlerror.Error{}

			if errList.As(&gqlerr) {
				tflog.Info(ctx, "####### gqlerror", map[string]interface{}{"gqlerror": gqlerr.Extensions})
				if errorCode, ok := gqlerr.Extensions["code"].(int); ok {
					tflog.Info(ctx, "####### if", map[string]interface{}{"error": errorCode})
				}

			}
			tflog.Info(ctx, "####### if after parsing")

			// return
		} else {
			tflog.Info(ctx, "####### else", map[string]interface{}{"error": err})
			// Handle cases where the error message is not JSON
			resp.Diagnostics.AddError(
				"Error Reading WX-ONE",
				"Could not read WX-ONE key ID "+state.ID.ValueString()+": "+err.Error(),
			)
		}
	}

	tflog.Info(ctx, "#######", map[string]interface{}{"key": key.GetKey.Code})

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *keyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *keyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
