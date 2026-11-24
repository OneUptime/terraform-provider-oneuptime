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

    // What a Crossplane manifest last applied, for the one managed resource
    // a provider made for its reconcile plans (NewUnderCrossplane,
    // shared.go). Nil under Terraform.
    lastApplied *LastApplied
}

// OneuptimeProviderModel describes the provider data model.
type OneuptimeProviderModel struct {
    OneuptimeUrl   types.String `tfsdk:"oneuptime_url"`
    ApiKey         types.String `tfsdk:"api_key"`
    ExclusiveLists types.Bool   `tfsdk:"exclusive_lists"`
}

func (p *OneuptimeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "oneuptime"
    resp.Version = p.version
}

func (p *OneuptimeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage a OneUptime project as code, with Terraform or OpenTofu: monitors, status pages, incidents, on-call policies, teams and everything else the OneUptime API manages.",

        Attributes: map[string]schema.Attribute{
            "oneuptime_url": schema.StringAttribute{
                MarkdownDescription: "The address of your OneUptime: https://oneuptime.com, or your own installation's. Leave out /api: the provider adds it. Falls back to the `ONEUPTIME_URL` environment variable, then to https://oneuptime.com.",
                Optional:            true,
            },
            "api_key": schema.StringAttribute{
                MarkdownDescription: "An API key of the project to manage, from its **Settings** > **API Keys**. Falls back to the `ONEUPTIME_API_KEY` environment variable; the provider stops before any plan when neither is set.",
                Optional:            true,
                Sensitive:           true,
            },
            "exclusive_lists": schema.BoolAttribute{
                MarkdownDescription: "How a configuration holds a list OneUptime adds to as well, such as a record's `labels`, which label rules add to. Left out or `false`, Terraform adds and takes off the entries the configuration lists and leaves every other entry: one OneUptime or a person adds stays. `true` makes the configuration's list the record's whole list, as earlier versions of the provider did: Terraform takes off every entry the configuration does not list, one a rule added included.",
                Optional:            true,
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
    var oneuptimeUrl string
    var apiKey string

    if data.OneuptimeUrl.IsUnknown() {
        resp.Diagnostics.AddError(
            "Unknown Provider Configuration",
            "oneuptime_url is not known at configure time. Set it to a static value or resolve the reference before applying.",
        )
        return
    }

    if data.OneuptimeUrl.IsNull() {
        oneuptimeUrl = os.Getenv("ONEUPTIME_URL")
        if oneuptimeUrl == "" {
            oneuptimeUrl = "oneuptime.com"
        }
    } else {
        oneuptimeUrl = data.OneuptimeUrl.ValueString()
    }

    if data.ApiKey.IsUnknown() {
        resp.Diagnostics.AddError(
            "Unknown Provider Configuration",
            "api_key is not known at configure time. Set it to a static value or resolve the reference before applying.",
        )
        return
    }

    if data.ApiKey.IsNull() {
        apiKey = os.Getenv("ONEUPTIME_API_KEY")
        if apiKey == "" {
            resp.Diagnostics.AddError(
                "Missing API Key",
                "API key is required for authentication. "+
                    "Please provide it via the api_key attribute or the ONEUPTIME_API_KEY environment variable.",
            )
            return
        }
    } else {
        apiKey = data.ApiKey.ValueString()
    }

    client, err := NewClient(oneuptimeUrl, apiKey, p.version)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create Oneuptime API Client",
            "An unexpected error occurred when creating the Oneuptime API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Oneuptime Client Error: "+err.Error(),
        )
        return
    }

    client.lastApplied = p.lastApplied

    // A list OneUptime adds to as well is the configuration's whole list
    // only when the configuration says so (entries.go). Never under
    // Crossplane, which holds a manifest's own entries.
    client.exclusiveLists = !data.ExclusiveLists.IsNull() && !data.ExclusiveLists.IsUnknown() && data.ExclusiveLists.ValueBool()

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

// NewUnderCrossplane returns the provider for one Crossplane reconcile of one
// managed resource: what its manifest last applied reaches the plan of the
// attributes people and OneUptime both change (shared.go). xpprovider calls
// it.
func NewUnderCrossplane(version string, lastApplied LastApplied) func() provider.Provider {
    return func() provider.Provider {
        return &OneuptimeProvider{
            version:     version,
            lastApplied: &lastApplied,
        }
    }
}
