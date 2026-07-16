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
var _ datasource.DataSource = &AiRunDataDataSource{}

func NewAiRunDataDataSource() datasource.DataSource {
    return &AiRunDataDataSource{}
}

// AiRunDataDataSource defines the data source implementation.
type AiRunDataDataSource struct {
    client *Client
}

// AiRunDataDataSourceModel describes the data source data model.
type AiRunDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    RunType types.String `tfsdk:"run_type"`
    CodeFixTaskType types.String `tfsdk:"code_fix_task_type"`
    TaskNumber types.Number `tfsdk:"task_number"`
    Status types.String `tfsdk:"status"`
    UserId types.String `tfsdk:"user_id"`
    ConversationId types.String `tfsdk:"conversation_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    TriggeredByTelemetryExceptionId types.String `tfsdk:"triggered_by_telemetry_exception_id"`
    TriggeredByAiInsightId types.String `tfsdk:"triggered_by_ai_insight_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    AiAgentId types.String `tfsdk:"ai_agent_id"`
    AttemptCount types.Number `tfsdk:"attempt_count"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    LastHeartbeatAt types.String `tfsdk:"last_heartbeat_at"`
    LlmCallCount types.Number `tfsdk:"llm_call_count"`
    ToolCallCount types.Number `tfsdk:"tool_call_count"`
    TotalTokens types.Number `tfsdk:"total_tokens"`
    TotalCostInUsdCents types.Number `tfsdk:"total_cost_in_usd_cents"`
    EgressManifest types.String `tfsdk:"egress_manifest"`
    ErrorMessage types.String `tfsdk:"error_message"`
    HumanVerdict types.String `tfsdk:"human_verdict"`
    HumanVerdictAt types.String `tfsdk:"human_verdict_at"`
    HumanVerdictByUserId types.String `tfsdk:"human_verdict_by_user_id"`
    AutoGrade types.String `tfsdk:"auto_grade"`
    AutoGradeAt types.String `tfsdk:"auto_grade_at"`
}

func (d *AiRunDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_run_data"
}

func (d *AiRunDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_run_data data source",

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
            "run_type": schema.StringAttribute{
                MarkdownDescription: "Type of AI run: Chat or Investigation.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "code_fix_task_type": schema.StringAttribute{
                MarkdownDescription: "For CodeFix runs: which task recipe this run executes (fix the exception, write a regression test, ...). Null means FixException — rows created before task recipes existed.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "task_number": schema.NumberAttribute{
                MarkdownDescription: "Per-project sequential number for this AI task (code-fix runs only).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this run.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "conversation_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_telemetry_exception_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_ai_insight_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "ai_agent_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "attempt_count": schema.NumberAttribute{
                MarkdownDescription: "How many times a worker has claimed this run for execution. Incremented on each claim; the queue stops retrying after the maximum.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
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
            "last_heartbeat_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "llm_call_count": schema.NumberAttribute{
                MarkdownDescription: "Number of LLM calls made during this run.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "tool_call_count": schema.NumberAttribute{
                MarkdownDescription: "Number of tool calls executed during this run.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "total_tokens": schema.NumberAttribute{
                MarkdownDescription: "Total LLM tokens used during this run.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "total_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total billed cost of this run in USD cents.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "egress_manifest": schema.StringAttribute{
                MarkdownDescription: "What data was sent to which LLM during this run: provider, model, and per-tool row/byte/redaction counts.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Error message if the run failed.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "human_verdict": schema.StringAttribute{
                MarkdownDescription: "For investigation runs: the one-click human verdict on the posted analysis (Confirmed or Rejected). Null until a user weighs in.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "human_verdict_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "human_verdict_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "auto_grade": schema.StringAttribute{
                MarkdownDescription: "For investigation runs: how the posted analysis compared to the incident's final recorded root cause (Match, Partial or Mismatch).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "auto_grade_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *AiRunDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiRunDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiRunDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-run" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_run_data, got error: %s", err))
        return
    }

    var aiRunDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiRunDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_run_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiRunDataResponse["data"].(map[string]interface{}); ok {
        aiRunDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiRunDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["run_type"].(string); ok {
        data.RunType = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["code_fix_task_type"].(string); ok {
        data.CodeFixTaskType = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["task_number"].(float64); ok {
        data.TaskNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["conversation_id"].(string); ok {
        data.ConversationId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["triggered_by_incident_id"].(string); ok {
        data.TriggeredByIncidentId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["triggered_by_alert_id"].(string); ok {
        data.TriggeredByAlertId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["triggered_by_telemetry_exception_id"].(string); ok {
        data.TriggeredByTelemetryExceptionId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["triggered_by_ai_insight_id"].(string); ok {
        data.TriggeredByAiInsightId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["ai_agent_id"].(string); ok {
        data.AiAgentId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["attempt_count"].(float64); ok {
        data.AttemptCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["started_at"].(string); ok {
        data.StartedAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["completed_at"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["last_heartbeat_at"].(string); ok {
        data.LastHeartbeatAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["llm_call_count"].(float64); ok {
        data.LlmCallCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["tool_call_count"].(float64); ok {
        data.ToolCallCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["total_tokens"].(float64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["total_cost_in_usd_cents"].(float64); ok {
        data.TotalCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiRunDataResponse["egress_manifest"].(string); ok {
        data.EgressManifest = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["error_message"].(string); ok {
        data.ErrorMessage = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["human_verdict"].(string); ok {
        data.HumanVerdict = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["human_verdict_at"].(string); ok {
        data.HumanVerdictAt = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["human_verdict_by_user_id"].(string); ok {
        data.HumanVerdictByUserId = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["auto_grade"].(string); ok {
        data.AutoGrade = types.StringValue(val)
    }
    if val, ok := aiRunDataResponse["auto_grade_at"].(string); ok {
        data.AutoGradeAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
