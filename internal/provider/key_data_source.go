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
	_ datasource.DataSource              = &keyDataSource{}
	_ datasource.DataSourceWithConfigure = &keyDataSource{}
)

// NewKeyDataSource is a helper function to simplify the provider implementation.
func NewKeyDataSource() datasource.DataSource {
	return &keyDataSource{}
}

// keyDataSource is the data source implementation.
type keyDataSource struct {
	wxOneClients *WxOneClients
}

// Configure adds the provider configured client to the data source.
func (d *keyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *keyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

// Schema defines the schema for the data source.
func (d *keyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Computed: true,
			},
			"project_wide": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

type keyDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PrivateKey  types.String `tfsdk:"private_key"`
	PublicKey   types.String `tfsdk:"public_key"`
	projectWide types.Bool   `tfsdk:"project_wide"`
}

// Read refreshes the Terraform state with the latest data.

func (d *keyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var state keyDataSourceModel

	// coffees, err := d.client.GetCoffees()
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to Read HashiCups Coffees",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// // Map response body to model
	// for _, coffee := range coffees {
	// 	coffeeState := coffeesModel{
	// 		ID:          types.Int64Value(int64(coffee.ID)),
	// 		Name:        types.StringValue(coffee.Name),
	// 		Teaser:      types.StringValue(coffee.Teaser),
	// 		Description: types.StringValue(coffee.Description),
	// 		Price:       types.Float64Value(coffee.Price),
	// 		Image:       types.StringValue(coffee.Image),
	// 	}

	// 	for _, ingredient := range coffee.Ingredient {
	// 		coffeeState.Ingredients = append(coffeeState.Ingredients, coffeesIngredientsModel{
	// 			ID: types.Int64Value(int64(ingredient.ID)),
	// 		})
	// 	}

	// 	state.Coffees = append(state.Coffees, coffeeState)
	// }

	// // Set state
	// diags := resp.State.Set(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}
