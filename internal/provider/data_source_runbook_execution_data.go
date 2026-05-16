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
var _ datasource.DataSource = &RunbookExecutionDataDataSource{}

func NewRunbookExecutionDataDataSource() datasource.DataSource {
    return &RunbookExecutionDataDataSource{}
}

// RunbookExecutionDataDataSource defines the data source implementation.
type RunbookExecutionDataDataSource struct {
    client *Client
}

// RunbookExecutionDataDataSourceModel describes the data source data model.
type RunbookExecutionDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    RunbookId types.String `tfsdk:"runbook_id"`
    RunbookNameSnapshot types.String `tfsdk:"runbook_name_snapshot"`
    Status types.String `tfsdk:"status"`
    StepExecutions types.String `tfsdk:"step_executions"`
    IncidentId types.String `tfsdk:"incident_id"`
    AlertId types.String `tfsdk:"alert_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    FailureReason types.String `tfsdk:"failure_reason"`
}

func (d *RunbookExecutionDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook_execution_data"
}

func (d *RunbookExecutionDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "runbook_execution_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
            "runbook_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "runbook_name_snapshot": schema.StringAttribute{
                MarkdownDescription: "Name of the runbook at the time this execution was created (preserved even if the runbook is later renamed or deleted).. Permissions - Create: [Project Owner, Project Admin, Create Runbook Execution, Project Member, Runbook Admin, Runbook Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Execution, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this runbook execution.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Execution, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Runbook Execution, Runbook Admin]",
                Computed: true,
            },
            "step_executions": schema.StringAttribute{
                MarkdownDescription: "Per-step execution state. Each entry mirrors a step from the runbook with status, output, and timestamps.. Permissions - Create: [Project Owner, Project Admin, Create Runbook Execution, Project Member, Runbook Admin, Runbook Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Execution, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Runbook Execution, Runbook Admin]",
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "failure_reason": schema.StringAttribute{
                MarkdownDescription: "Reason this runbook execution failed (if it did).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Runbook Admin, Runbook Member, Runbook Viewer, Read Runbook Execution, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Runbook Execution, Runbook Admin]",
                Computed: true,
            },
        },
    }
}

func (d *RunbookExecutionDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookExecutionDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookExecutionDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "runbook-execution" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook_execution_data, got error: %s", err))
        return
    }

    var runbookExecutionDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &runbookExecutionDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse runbook_execution_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := runbookExecutionDataResponse["data"].(map[string]interface{}); ok {
        runbookExecutionDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := runbookExecutionDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := runbookExecutionDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["runbook_id"].(string); ok {
        data.RunbookId = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["runbook_name_snapshot"].(string); ok {
        data.RunbookNameSnapshot = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["step_executions"].(string); ok {
        data.StepExecutions = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["triggered_by_user_id"].(string); ok {
        data.TriggeredByUserId = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["started_at"].(string); ok {
        data.StartedAt = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["completed_at"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    }
    if val, ok := runbookExecutionDataResponse["failure_reason"].(string); ok {
        data.FailureReason = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
