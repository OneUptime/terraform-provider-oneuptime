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
var _ datasource.DataSource = &TraceDropFilterDataDataSource{}

func NewTraceDropFilterDataDataSource() datasource.DataSource {
    return &TraceDropFilterDataDataSource{}
}

// TraceDropFilterDataDataSource defines the data source implementation.
type TraceDropFilterDataDataSource struct {
    client *Client
}

// TraceDropFilterDataDataSourceModel describes the data source data model.
type TraceDropFilterDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    FilterQuery types.String `tfsdk:"filter_query"`
    Action types.String `tfsdk:"action"`
    SamplePercentage types.Number `tfsdk:"sample_percentage"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *TraceDropFilterDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_trace_drop_filter_data"
}

func (d *TraceDropFilterDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "trace_drop_filter_data data source",

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
                MarkdownDescription: "Description of what this drop filter does.. Permissions - Create: [Project Owner, Project Admin, Create Trace Drop Filter], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Drop Filter, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Drop Filter]",
                Computed: true,
            },
            "filter_query": schema.StringAttribute{
                MarkdownDescription: "Filter expression that identifies which spans to drop or sample.. Permissions - Create: [Project Owner, Project Admin, Create Trace Drop Filter], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Drop Filter, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Drop Filter]",
                Computed: true,
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "What to do with matching spans: 'drop' to discard entirely, 'sample' to keep a percentage.. Permissions - Create: [Project Owner, Project Admin, Create Trace Drop Filter], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Drop Filter, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Drop Filter]",
                Computed: true,
            },
            "sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "When action is 'sample', the percentage of matching spans to keep (1-99).. Permissions - Create: [Project Owner, Project Admin, Create Trace Drop Filter], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Drop Filter, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Drop Filter]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this drop filter is active.. Permissions - Create: [Project Owner, Project Admin, Create Trace Drop Filter], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Drop Filter, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Drop Filter]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Determines the evaluation order of this filter relative to others.. Permissions - Create: [Project Owner, Project Admin, Create Trace Drop Filter], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Trace Drop Filter, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Drop Filter]",
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

func (d *TraceDropFilterDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TraceDropFilterDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TraceDropFilterDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "trace-drop-filter" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trace_drop_filter_data, got error: %s", err))
        return
    }

    var traceDropFilterDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &traceDropFilterDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse trace_drop_filter_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := traceDropFilterDataResponse["data"].(map[string]interface{}); ok {
        traceDropFilterDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := traceDropFilterDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceDropFilterDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["filter_query"].(string); ok {
        data.FilterQuery = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["action"].(string); ok {
        data.Action = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["sample_percentage"].(float64); ok {
        data.SamplePercentage = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceDropFilterDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := traceDropFilterDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceDropFilterDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := traceDropFilterDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
