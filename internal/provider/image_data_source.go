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
	_ datasource.DataSource              = &imageDataSource{}
	_ datasource.DataSourceWithConfigure = &imageDataSource{}
)

// NewImageDataSource is a helper function to simplify the provider implementation.
func NewImageDataSource() datasource.DataSource {
	return &imageDataSource{}
}

// imageDataSource is the data source implementation.
type imageDataSource struct {
	wxOneClients *WxOneClients
}

// Configure adds the provider configured client to the data source.
func (d *imageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *imageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Schema defines the schema for the data source.
func (d *imageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve an image to be used for example to create an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the image in UUID format.",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "ID of the project in UUID format.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the image.",
				Required:    true,
			},
		},
	}
}

type imageDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
}

// Read refreshes the Terraform state with the latest data.

func (d *imageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data imageDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageList, err := getImageList(ctx, d.wxOneClients.graphqlClient, data.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read WX-ONE Images",
			err.Error(),
		)
		return
	}

	imageFound := false
	for _, image := range imageList.GetImageList.Msg {
		if image.Name == data.Name.ValueString() {
			imageFound = true
			data.ID = types.StringValue(image.Id)
		}
	}

	if imageFound == false {
		resp.Diagnostics.AddError(
			"No Image Found", fmt.Sprintf("No image found with name: %s", data.Name.ValueString()),
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
