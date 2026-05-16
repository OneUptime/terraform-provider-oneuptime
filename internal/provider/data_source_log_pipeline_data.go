package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LogPipelineDataDataSource{}

func NewLogPipelineDataDataSource() datasource.DataSource {
    return &LogPipelineDataDataSource{}
}

// LogPipelineDataDataSource defines the data source implementation.
type LogPipelineDataDataSource struct {
    client *Client
}

// LogPipelineDataDataSourceModel describes the data source data model.
type LogPipelineDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    FilterQuery types.String `tfsdk:"filter_query"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *LogPipelineDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log_pipeline_data"
}

func (d *LogPipelineDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "log_pipeline_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this log pipeline does.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Log Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline]",
                Computed: true,
            },
            "filter_query": schema.StringAttribute{
                MarkdownDescription: "Filter expression that determines which logs this pipeline applies to.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Log Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this log pipeline is active.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Log Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Determines the execution order of this pipeline relative to others.. Permissions - Create: [Project Owner, Project Admin, Create Log Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Log Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Log Pipeline]",
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

func (d *LogPipelineDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogPipelineDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LogPipelineDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "log-pipeline" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log_pipeline_data, got error: %s", err))
        return
    }

    var logPipelineDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &logPipelineDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log_pipeline_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := logPipelineDataResponse["data"].(map[string]interface{}); ok {
        logPipelineDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := logPipelineDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logPipelineDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["filter_query"].(string); ok {
        data.FilterQuery = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := logPipelineDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logPipelineDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := logPipelineDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
