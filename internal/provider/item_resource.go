package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/superorbital/inventory-service/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &itemResource{}
	_ resource.ResourceWithConfigure   = &itemResource{}
	_ resource.ResourceWithImportState = &itemResource{}
)

// NewItemResource is a helper function to simplify the provider implementation.
func NewItemResource() resource.Resource {
	return &itemResource{}
}

// itemResource is the resource implementation.
type itemResource struct {
	client *client.Client
}

// itemResourceModel maps the resource schema data.
type itemResourceModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Tag  types.String `tfsdk:"tag"`
}

// Configure adds the provider configured client to the resource.
func (r *itemResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		tflog.Error(ctx, "Unable to prepare client")
		return
	}
	r.client = client

}

// Metadata returns the resource type name.
func (r *itemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_item"
}

// Schema defines the schema for the resource.
func (r *itemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an item.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Identifier for this inventory item.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name for this inventory item.",
				Required:    true,
			},
			"tag": schema.StringAttribute{
				Description: "The tag for this inventory item.",
				Optional:    true,
			},
		},
	}
}

func (r *itemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// If our ID was a string then we could do this
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id, err := strconv.ParseInt(req.ID, 10, 64)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing item",
			"Could not import item, unexpected error (ID should be an integer): "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Create a new resource.
func (r *itemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Preparing to create item resource")
	// Retrieve values from plan
	var plan itemResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	tag := plan.Tag.ValueString()

	item := client.NewItem{
		Name: name,
		Tag:  &tag,
	}

	// Create new item

	itemResponse, err := r.client.AddItem(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Item",
			err.Error(),
		)
		return
	}

	var newItem client.Item
	if err := json.NewDecoder(itemResponse.Body).Decode(&newItem); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}

	// Map response body to model
	plan.ID = types.Int64Value(newItem.Id)
	plan.Name = types.StringValue(newItem.Name)
	plan.Tag = types.StringValue(*newItem.Tag)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created item resource", map[string]any{"success": true})
}

// Read resource information.
func (r *itemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read item resource")
	// Get current state
	var state itemResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	itemResponse, err := r.client.FindItemById(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Item",
			err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to remove/recreate resource
	if itemResponse.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if itemResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Item",
			itemResponse.Status,
		)
		return
	}

	var newItem client.Item
	if err := json.NewDecoder(itemResponse.Body).Decode(&newItem); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = itemResourceModel{
		ID:   types.Int64Value(newItem.Id),
		Name: types.StringValue(newItem.Name),
		Tag:  types.StringValue(*newItem.Tag),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading item resource", map[string]any{"success": true})
}

func (r *itemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update item resource")
	// Retrieve values from plan
	var plan itemResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	tag := plan.Tag.ValueString()

	item := client.NewItem{
		Name: name,
		Tag:  &tag,
	}

	// update item
	itemResponse, err := r.client.UpdateItem(ctx, plan.ID.ValueInt64(), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Item",
			err.Error(),
		)
		return
	}

	if itemResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Item",
			itemResponse.Status,
		)
		return
	}

	var newItem client.Item
	if err := json.NewDecoder(itemResponse.Body).Decode(&newItem); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	plan = itemResourceModel{
		ID:   types.Int64Value(newItem.Id),
		Name: types.StringValue(newItem.Name),
		Tag:  types.StringValue(*newItem.Tag),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updated item resource", map[string]any{"success": true})
}

func (r *itemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete item resource")
	// Retrieve values from state
	var state itemResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete item
	_, err := r.client.DeleteItem(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Item",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}
