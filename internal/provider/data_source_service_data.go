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
var _ datasource.DataSource = &ServiceDataDataSource{}

func NewServiceDataDataSource() datasource.DataSource {
    return &ServiceDataDataSource{}
}

// ServiceDataDataSource defines the data source implementation.
type ServiceDataDataSource struct {
    client *Client
}

// ServiceDataDataSourceModel describes the data source data model.
type ServiceDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    ServiceColor types.String `tfsdk:"service_color"`
    ServiceLanguage types.String `tfsdk:"service_language"`
    TechStack types.String `tfsdk:"tech_stack"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    MetricCardinalityBudget types.Number `tfsdk:"metric_cardinality_budget"`
    MetricDownsamplingRetentionDays types.String `tfsdk:"metric_downsampling_retention_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
}

func (d *ServiceDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_data"
}

func (d *ServiceDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
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
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "service_language": schema.StringAttribute{
                MarkdownDescription: "Language in which this service is written",
                Computed: true,
            },
            "tech_stack": schema.StringAttribute{
                MarkdownDescription: "Tech stack used in the service. This will help other developers understand the service better.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
                Computed: true,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this service. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
                Computed: true,
            },
            "metric_cardinality_budget": schema.NumberAttribute{
                MarkdownDescription: "Max number of distinct metric series this service may emit per metric. When exceeded, the highest-cardinality attribute is auto-bucketed. Null inherits the project default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
                Computed: true,
            },
            "metric_downsampling_retention_days": schema.StringAttribute{
                MarkdownDescription: "Per-tier retention override (raw, 1m, 5m, 1h, 1d) in days. Null fields inherit the project default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this service (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the service default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Service], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Service], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Service]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *ServiceDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "service" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_data, got error: %s", err))
        return
    }

    var serviceDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &serviceDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := serviceDataResponse["data"].(map[string]interface{}); ok {
        serviceDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := serviceDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := serviceDataResponse["service_color"].(string); ok {
        data.ServiceColor = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["service_language"].(string); ok {
        data.ServiceLanguage = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["tech_stack"].(string); ok {
        data.TechStack = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceDataResponse["metric_cardinality_budget"].(float64); ok {
        data.MetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := serviceDataResponse["metric_downsampling_retention_days"].(string); ok {
        data.MetricDownsamplingRetentionDays = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }
    if val, ok := serviceDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
