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
var _ datasource.DataSource = &TracePipelineDataDataSource{}

func NewTracePipelineDataDataSource() datasource.DataSource {
    return &TracePipelineDataDataSource{}
}

// TracePipelineDataDataSource defines the data source implementation.
type TracePipelineDataDataSource struct {
    client *Client
}

// TracePipelineDataDataSourceModel describes the data source data model.
type TracePipelineDataDataSourceModel struct {
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

func (d *TracePipelineDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_trace_pipeline_data"
}

func (d *TracePipelineDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "trace_pipeline_data data source",

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
                MarkdownDescription: "Description of what this trace pipeline does.. Permissions - Create: [Project Owner, Project Admin, Create Trace Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Pipeline]",
                Computed: true,
            },
            "filter_query": schema.StringAttribute{
                MarkdownDescription: "Filter expression that determines which spans this pipeline applies to.. Permissions - Create: [Project Owner, Project Admin, Create Trace Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Pipeline]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this trace pipeline is active.. Permissions - Create: [Project Owner, Project Admin, Create Trace Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Pipeline]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Determines the execution order of this pipeline relative to others.. Permissions - Create: [Project Owner, Project Admin, Create Trace Pipeline], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Pipeline, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Pipeline]",
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

func (d *TracePipelineDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TracePipelineDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TracePipelineDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "trace-pipeline" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trace_pipeline_data, got error: %s", err))
        return
    }

    var tracePipelineDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &tracePipelineDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse trace_pipeline_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := tracePipelineDataResponse["data"].(map[string]interface{}); ok {
        tracePipelineDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := tracePipelineDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := tracePipelineDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["filter_query"].(string); ok {
        data.FilterQuery = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := tracePipelineDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := tracePipelineDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := tracePipelineDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
