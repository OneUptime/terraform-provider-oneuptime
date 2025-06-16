package oneuptime

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &oneuptimeProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func NewProvider(version string) func() provider.Provider {
    return func() provider.Provider {
        return &oneuptimeProvider{
            version: version,
        }
    }
}

// oneuptimeProvider is the provider implementation.
type oneuptimeProvider struct {
    version string
}

// Metadata returns the provider type name.
func (p *oneuptimeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "oneuptime"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *oneuptimeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Interact with OneUptime.",
        Attributes: map[string]schema.Attribute{
            "api_url": schema.StringAttribute{
                Description: "OneUptime API URL. May also be provided via ONEUPTIME_API_URL environment variable.",
                Optional:    true,
            },
            "api_key": schema.StringAttribute{
                Description: "OneUptime API Key. May also be provided via ONEUPTIME_API_KEY environment variable.",
                Optional:    true,
                Sensitive:   true,
            },
        },
    }
}

type oneuptimeProviderModel struct {
    ApiUrl types.String `tfsdk:"api_url"`
    ApiKey types.String `tfsdk:"api_key"`
}

// Configure prepares a OneUptime API client for data sources and resources.
func (p *oneuptimeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var config oneuptimeProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If configuration values are known, set them here
    // This is where you would initialize your API client
}

// DataSources defines the data sources implemented in the provider.
func (p *oneuptimeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        // Add your data sources here
    }
}

// Resources defines the resources implemented in the provider.
func (p *oneuptimeProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        // Add your resources here
    }
}
