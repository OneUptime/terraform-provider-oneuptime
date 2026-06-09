package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServerlessFunctionDataDataSource{}

func NewServerlessFunctionDataDataSource() datasource.DataSource {
    return &ServerlessFunctionDataDataSource{}
}

// ServerlessFunctionDataDataSource defines the data source implementation.
type ServerlessFunctionDataDataSource struct {
    client *Client
}

// ServerlessFunctionDataDataSourceModel describes the data source data model.
type ServerlessFunctionDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    FunctionIdentifier types.String `tfsdk:"function_identifier"`
    CloudPlatform types.String `tfsdk:"cloud_platform"`
    CloudProvider types.String `tfsdk:"cloud_provider"`
    CloudRegion types.String `tfsdk:"cloud_region"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    FunctionVersion types.String `tfsdk:"function_version"`
    RuntimeName types.String `tfsdk:"runtime_name"`
    RuntimeVersion types.String `tfsdk:"runtime_version"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *ServerlessFunctionDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_serverless_function_data"
}

func (d *ServerlessFunctionDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "serverless_function_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Serverless Function], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Serverless Function]",
                Computed: true,
            },
            "function_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identifier from the faas.name OpenTelemetry resource attribute. Identity key for this function.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Serverless Function], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cloud_platform": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_lambda, gcp_cloud_functions, azure_functions.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cloud_provider": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cloud_region": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.account.id OpenTelemetry resource attribute.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "function_version": schema.StringAttribute{
                MarkdownDescription: "Last-seen faas.version OpenTelemetry resource attribute.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "runtime_name": schema.StringAttribute{
                MarkdownDescription: "Last-seen process.runtime.name OpenTelemetry resource attribute.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "runtime_version": schema.StringAttribute{
                MarkdownDescription: "Last-seen process.runtime.version OpenTelemetry resource attribute.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Whether telemetry is currently being received (connected) or has gone stale (disconnected).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime agent reporting this function.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Serverless Function], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Serverless Function]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this function. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Serverless Function], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Serverless Function]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this function. Unset fields fall back to the function default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Serverless Function], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Serverless Function], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Serverless Function]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *ServerlessFunctionDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *ServerlessFunctionDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServerlessFunctionDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "serverless-function" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read serverless_function_data, got error: %s", err))
        return
    }

    var serverlessFunctionDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serverlessFunctionDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse serverless_function_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serverlessFunctionDataResponse["data"].(map[string]interface{}); ok {
        serverlessFunctionDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serverlessFunctionDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serverlessFunctionDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["function_identifier"].(string); ok {
        data.FunctionIdentifier = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["cloud_platform"].(string); ok {
        data.CloudPlatform = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["cloud_provider"].(string); ok {
        data.CloudProvider = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["cloud_region"].(string); ok {
        data.CloudRegion = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["cloud_account_id"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["function_version"].(string); ok {
        data.FunctionVersion = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["runtime_name"].(string); ok {
        data.RuntimeName = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["runtime_version"].(string); ok {
        data.RuntimeVersion = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := serverlessFunctionDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serverlessFunctionDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := serverlessFunctionDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
