package provider

import (
	"context"
	"os"

	"github.com/superorbital/inventory-service/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &inventoryProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &inventoryProvider{
			version: version,
		}
	}
}

// inventoryProvider is the provider implementation.
type inventoryProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// inventoryProviderModel maps provider schema data to a Go type.
type inventoryProviderModel struct {
	Host types.String `tfsdk:"host"`
	Port types.String `tfsdk:"port"`
}

// Metadata returns the provider type name.
func (p *inventoryProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "inventory"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *inventoryProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The hostname or IP address for the inventory service endpoint. May also be provided via the INVENTORY_HOST environment variable.",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "The port to connect to. May also be provided via the INVENTORY_PORT environment variable.",
			},
		},
		Blocks:      map[string]schema.Block{},
		Description: "Interface with the Inventory service API.",
	}
}

// Configure prepares a Inventory API client for data sources and resources.
//
//gocyclo:ignore
func (p *inventoryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Inventory client")

	// Retrieve provider data from configuration
	var config inventoryProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Inventory service Host",
			"The provider cannot create the Inventory API client as there is an unknown configuration value for the Inventory API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INVENTORY_HOST environment variable.",
		)
	}

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Unknown Inventory service Port",
			"The provider cannot create the Inventory API client as there is an unknown configuration value for the Inventory API port. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INVENTORY_PORT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("INVENTORY_HOST")
	port := os.Getenv("INVENTORY_PORT")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = config.Port.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("host"),
			"Missing Inventory API Host (using default value: 127.0.0.1)",
			"The provider is using a default value as there is a missing or empty value for the Inventory API host. "+
				"Set the host value in the configuration or use the INVENTORY_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		host = "127.0.0.1"
	}

	if port == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("port"),
			"Missing Inventory API port (using default value: 8080)",
			"The provider is using a default value as there is a missing or empty value for the Inventory API host. "+
				"Set the host value in the configuration or use the INVENTORY_PORT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		port = "8080"
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Inventory client")

	// Instantiate the client that we will use to talk to the Inventory server
	serverURL := "http://" + host + ":" + port + "/"
	api, err := client.NewClient(serverURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Inventory API Client",
			"An unexpected error occurred when creating the Inventory API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Inventory Client Error: "+err.Error(),
		)
		return
	}
	// Test that we have some basic connectivity
	_, err = api.FindItemById(ctx, int64(1))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Inventory API Client",
			"An unexpected error occurred when creating the Inventory API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Inventory Client Error: "+err.Error(),
		)
		return
	}

	// Make the Inventory client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = api
	resp.ResourceData = api

	tflog.Info(ctx, "Configured Inventory client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *inventoryProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewItemDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *inventoryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewItemResource,
	}
}
