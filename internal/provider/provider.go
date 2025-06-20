package provider

import (
    "context"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &OneuptimeProvider{}

// OneuptimeProvider defines the provider implementation.
type OneuptimeProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// OneuptimeProviderModel describes the provider data model.
type OneuptimeProviderModel struct {
    Host     types.String `tfsdk:"host"`
    ApiKey   types.String `tfsdk:"api_key"`
}

func (p *OneuptimeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "oneuptime"
    resp.Version = p.version
}

func (p *OneuptimeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "OpenAPI specification for OneUptime. This document describes the API endpoints, request and response formats, and other details necessary for developers to interact with the OneUptime API.",

        Attributes: map[string]schema.Attribute{
            "host": schema.StringAttribute{
                MarkdownDescription: "The oneuptime host (without /api path). Defaults to 'oneuptime.com' if not specified. The provider automatically appends '/api' to the host.",
                Optional:            true,
            },
            "api_key": schema.StringAttribute{
                MarkdownDescription: "API key for authentication",
                Optional:            true,
                Sensitive:           true,
            },
        },
    }
}

func (p *OneuptimeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var data OneuptimeProviderModel

    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Configuration values are now available.
    var host string
    var apiKey string

    if data.Host.IsUnknown() {
        // Cannot connect to client with an unknown value
        resp.Diagnostics.AddWarning(
            "Unable to create client",
            "Cannot use unknown value as host",
        )
        return
    }

    if data.Host.IsNull() {
        host = os.Getenv("ONEUPTIME_HOST")
        if host == "" {
            host = "oneuptime.com"
        }
    } else {
        host = data.Host.ValueString()
    }

    if data.ApiKey.IsNull() {
        apiKey = os.Getenv("ONEUPTIME_API_KEY")
    } else {
        apiKey = data.ApiKey.ValueString()
    }

    client, err := NewClient(host, apiKey)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create Oneuptime API Client",
            "An unexpected error occurred when creating the Oneuptime API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Oneuptime Client Error: "+err.Error(),
        )
        return
    }

    resp.DataSourceData = client
    resp.ResourceData = client

    tflog.Info(ctx, "Configured Oneuptime client", map[string]any{"success": true})
}

func (p *OneuptimeProvider) Resources(ctx context.Context) []func() resource.Resource {
    return GetResources()
}

func (p *OneuptimeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return GetDataSources()
}

func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &OneuptimeProvider{
            version: version,
        }
    }
}
