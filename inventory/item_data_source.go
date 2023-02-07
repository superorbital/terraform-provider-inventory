package inventory

import (
	"context"
	"encoding/json"

	"github.com/superorbital/inventory-service/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &itemDataSource{}
	_ datasource.DataSourceWithConfigure = &itemDataSource{}
)

// NewItemDataSource is a helper function to simplify the provider implementation.
func NewItemDataSource() datasource.DataSource {
	return &itemDataSource{}
}

// itemDataSource is the data source implementation.
type itemDataSource struct {
	client *client.Client
}

// itemDataSourceModel maps the data source schema data.
type itemDataSourceModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Tag  types.String `tfsdk:"tag"`
}

// Configure adds the provider configured client to the data source.
func (d *itemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)

}

// Metadata returns the data source type name.
func (d *itemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_item"
}

// Schema defines the schema for the data source.
func (d *itemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch an item.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Identifier for this inventory item.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name for this inventory item.",
				Computed:    true,
			},
			"tag": schema.StringAttribute{
				Description: "The tag for this inventory item.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *itemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read item data source")
	var state itemDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	itemResponse, err := d.client.FindItemById(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Item",
			err.Error(),
		)
		return
	}

	var newItem client.Item
	if itemResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Item",
			itemResponse.Status,
		)
		return
	}

	if err := json.NewDecoder(itemResponse.Body).Decode(&newItem); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = itemDataSourceModel{
		ID:   types.Int64Value(newItem.Id),
		Name: types.StringValue(newItem.Name),
		Tag:  types.StringValue(*newItem.Tag),
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, "Finished reading item data source", map[string]any{"success": true})
}
