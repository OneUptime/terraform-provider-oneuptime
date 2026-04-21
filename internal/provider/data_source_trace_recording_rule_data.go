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
var _ datasource.DataSource = &TraceRecordingRuleDataDataSource{}

func NewTraceRecordingRuleDataDataSource() datasource.DataSource {
    return &TraceRecordingRuleDataDataSource{}
}

// TraceRecordingRuleDataDataSource defines the data source implementation.
type TraceRecordingRuleDataDataSource struct {
    client *Client
}

// TraceRecordingRuleDataDataSourceModel describes the data source data model.
type TraceRecordingRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    OutputMetricName types.String `tfsdk:"output_metric_name"`
    Definition types.String `tfsdk:"definition"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *TraceRecordingRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_trace_recording_rule_data"
}

func (d *TraceRecordingRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "trace_recording_rule_data data source",

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
                MarkdownDescription: "What this recording rule computes and why.. Permissions - Create: [Project Owner, Project Admin, Create Trace Recording Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Recording Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Recording Rule]",
                Computed: true,
            },
            "output_metric_name": schema.StringAttribute{
                MarkdownDescription: "Name of the new metric this rule writes (e.g. http.error_rate). Must be unique per project.. Permissions - Create: [Project Owner, Project Admin, Create Trace Recording Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Recording Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Recording Rule]",
                Computed: true,
            },
            "definition": schema.StringAttribute{
                MarkdownDescription: "Sources (aliased span aggregations), arithmetic expression, and optional group-by attribute.. Permissions - Create: [Project Owner, Project Admin, Create Trace Recording Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Recording Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Recording Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is evaluated by the recording rule cron.. Permissions - Create: [Project Owner, Project Admin, Create Trace Recording Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Recording Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Recording Rule]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Evaluation order when multiple rules exist.. Permissions - Create: [Project Owner, Project Admin, Create Trace Recording Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Recording Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Recording Rule]",
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

func (d *TraceRecordingRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TraceRecordingRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TraceRecordingRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "trace-recording-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trace_recording_rule_data, got error: %s", err))
        return
    }

    var traceRecordingRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &traceRecordingRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse trace_recording_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := traceRecordingRuleDataResponse["data"].(map[string]interface{}); ok {
        traceRecordingRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := traceRecordingRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceRecordingRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["output_metric_name"].(string); ok {
        data.OutputMetricName = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["definition"].(string); ok {
        data.Definition = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceRecordingRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := traceRecordingRuleDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
