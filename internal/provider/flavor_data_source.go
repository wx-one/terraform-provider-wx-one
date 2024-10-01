// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &flavorDataSource{}
	_ datasource.DataSourceWithConfigure = &flavorDataSource{}
)

// NewFlavorDataSource is a helper function to simplify the provider implementation.
func NewFlavorDataSource() datasource.DataSource {
	return &flavorDataSource{}
}

// flavorDataSource is the data source implementation.
type flavorDataSource struct {
	wxOneClients *WxOneClients
}

// Configure adds the provider configured client to the data source.
func (d *flavorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.wxOneClients = wxOneClients
}

// Metadata returns the data source type name.
func (d *flavorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flavor"
}

// Schema defines the schema for the data source.
func (d *flavorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve an flavor to be used for example to create an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the flavor in UUID format.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the flavor.",
				Required:    true,
			},
		},
	}
}

type flavorDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Read refreshes the Terraform state with the latest data.

func (d *flavorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data flavorDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flavor, err := getFlavorByName(ctx, d.wxOneClients.graphqlClient, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read WX-ONE Flavor",
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(flavor.GetFlavorByName.Msg.Id)

	if data.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"No Flavor Found", fmt.Sprintf("No flavor found with name: %s", data.Name.ValueString()),
		)
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
