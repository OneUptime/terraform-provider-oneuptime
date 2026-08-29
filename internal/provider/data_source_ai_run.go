package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AiRunDataSource{}

func NewAiRunDataSource() datasource.DataSource {
    return &AiRunDataSource{}
}

// AiRunDataSource defines the data source implementation.
type AiRunDataSource struct {
    client *Client
}

// AiRunDataSourceModel describes the data source data model.
type AiRunDataSourceModel struct {
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
    CodeFixRecommendation types.String `tfsdk:"code_fix_recommendation"`
    UserId types.String `tfsdk:"user_id"`
    ConversationId types.String `tfsdk:"conversation_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    TriggeredByTelemetryExceptionId types.String `tfsdk:"triggered_by_telemetry_exception_id"`
    TriggeredByAiInsightId types.String `tfsdk:"triggered_by_ai_insight_id"`
    TriggeredByAutoRemediationSuggestionId types.String `tfsdk:"triggered_by_auto_remediation_suggestion_id"`
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
    AnalysisTldr types.String `tfsdk:"analysis_tldr"`
    HumanVerdict types.String `tfsdk:"human_verdict"`
    HumanVerdictAt types.String `tfsdk:"human_verdict_at"`
    HumanVerdictByUserId types.String `tfsdk:"human_verdict_by_user_id"`
    AutoGrade types.String `tfsdk:"auto_grade"`
    AutoGradeAt types.String `tfsdk:"auto_grade_at"`
}

func (d *AiRunDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_run"
}

func (d *AiRunDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "One AI agent execution: LLM calls, tool calls, cost, and the egress manifest of what was sent to the LLM. Look up an existing ai_run by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
            "run_type": schema.StringAttribute{
                MarkdownDescription: "Type of AI run: Chat or Investigation..",
                Computed: true,
            },
            "code_fix_task_type": schema.StringAttribute{
                MarkdownDescription: "For CodeFix runs: which task recipe this run executes (fix the exception, write a regression test, ...). Null means FixException — rows created before task recipes existed..",
                Computed: true,
            },
            "task_number": schema.NumberAttribute{
                MarkdownDescription: "Per-project sequential number for this AI task (code-fix runs only)..",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this run..",
                Computed: true,
            },
            "code_fix_recommendation": schema.StringAttribute{
                MarkdownDescription: "For incident/alert investigations: whether the structured investigation outcome recommends opening a code-fix pull request..",
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
            "triggered_by_auto_remediation_suggestion_id": schema.StringAttribute{
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
                MarkdownDescription: "How many times a worker has claimed this run for execution. Incremented on each claim; the queue stops retrying after the maximum..",
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
                MarkdownDescription: "Number of LLM calls made during this run..",
                Computed: true,
            },
            "tool_call_count": schema.NumberAttribute{
                MarkdownDescription: "Number of tool calls executed during this run..",
                Computed: true,
            },
            "total_tokens": schema.NumberAttribute{
                MarkdownDescription: "Total LLM tokens used during this run..",
                Computed: true,
            },
            "total_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total billed cost of this run in USD cents..",
                Computed: true,
            },
            "egress_manifest": schema.StringAttribute{
                MarkdownDescription: "What data was sent to which LLM during this run: provider, model, and per-tool row/byte/redaction counts..",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Error message if the run failed..",
                Computed: true,
            },
            "analysis_tldr": schema.StringAttribute{
                MarkdownDescription: "For investigation runs: a one or two sentence plain-text summary of the posted analysis, generated by AI..",
                Computed: true,
            },
            "human_verdict": schema.StringAttribute{
                MarkdownDescription: "For investigation runs: the one-click human verdict on the posted analysis (Confirmed or Rejected). Null until a user weighs in..",
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
                MarkdownDescription: "For investigation runs: how the posted analysis compared to the incident's final recorded root cause (Match, Partial or Mismatch)..",
                Computed: true,
            },
            "auto_grade_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *AiRunDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiRunDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a ai_run.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "runType": true,
        "codeFixTaskType": true,
        "taskNumber": true,
        "status": true,
        "codeFixRecommendation": true,
        "userId": true,
        "conversationId": true,
        "triggeredByIncidentId": true,
        "triggeredByAlertId": true,
        "triggeredByTelemetryExceptionId": true,
        "triggeredByAiInsightId": true,
        "triggeredByAutoRemediationSuggestionId": true,
        "monitorId": true,
        "aiAgentId": true,
        "attemptCount": true,
        "startedAt": true,
        "completedAt": true,
        "lastHeartbeatAt": true,
        "llmCallCount": true,
        "toolCallCount": true,
        "totalTokens": true,
        "totalCostInUSDCents": true,
        "egressManifest": true,
        "errorMessage": true,
        "analysisTldr": true,
        "humanVerdict": true,
        "humanVerdictAt": true,
        "humanVerdictByUserId": true,
        "autoGrade": true,
        "autoGradeAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/ai-run/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_run, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_run found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_run: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/ai-run/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list ai_run, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list ai_run: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_run found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai_run matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for ai_run.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["runType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunType = types.StringValue(string(jsonBytes))
        } else {
            data.RunType = types.StringNull()
        }
    } else if val, ok := item["runType"].(string); ok {
        data.RunType = types.StringValue(val)
    } else {
        data.RunType = types.StringNull()
    }
    if obj, ok := item["codeFixTaskType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeFixTaskType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CodeFixTaskType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CodeFixTaskType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CodeFixTaskType = types.StringValue(string(jsonBytes))
        } else {
            data.CodeFixTaskType = types.StringNull()
        }
    } else if val, ok := item["codeFixTaskType"].(string); ok {
        data.CodeFixTaskType = types.StringValue(val)
    } else {
        data.CodeFixTaskType = types.StringNull()
    }
    if val, ok := item["taskNumber"].(float64); ok {
        data.TaskNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["taskNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TaskNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.TaskNumber = types.NumberNull()
        }
    } else {
        data.TaskNumber = types.NumberNull()
    }
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["codeFixRecommendation"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeFixRecommendation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CodeFixRecommendation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CodeFixRecommendation = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CodeFixRecommendation = types.StringValue(string(jsonBytes))
        } else {
            data.CodeFixRecommendation = types.StringNull()
        }
    } else if val, ok := item["codeFixRecommendation"].(string); ok {
        data.CodeFixRecommendation = types.StringValue(val)
    } else {
        data.CodeFixRecommendation = types.StringNull()
    }
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if obj, ok := item["conversationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConversationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConversationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConversationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConversationId = types.StringValue(string(jsonBytes))
        } else {
            data.ConversationId = types.StringNull()
        }
    } else if val, ok := item["conversationId"].(string); ok {
        data.ConversationId = types.StringValue(val)
    } else {
        data.ConversationId = types.StringNull()
    }
    if obj, ok := item["triggeredByIncidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByIncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByIncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByIncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByIncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByIncidentId = types.StringNull()
        }
    } else if val, ok := item["triggeredByIncidentId"].(string); ok {
        data.TriggeredByIncidentId = types.StringValue(val)
    } else {
        data.TriggeredByIncidentId = types.StringNull()
    }
    if obj, ok := item["triggeredByAlertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByAlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByAlertId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByAlertId = types.StringNull()
        }
    } else if val, ok := item["triggeredByAlertId"].(string); ok {
        data.TriggeredByAlertId = types.StringValue(val)
    } else {
        data.TriggeredByAlertId = types.StringNull()
    }
    if obj, ok := item["triggeredByTelemetryExceptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByTelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := item["triggeredByTelemetryExceptionId"].(string); ok {
        data.TriggeredByTelemetryExceptionId = types.StringValue(val)
    } else {
        data.TriggeredByTelemetryExceptionId = types.StringNull()
    }
    if obj, ok := item["triggeredByAiInsightId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByAiInsightId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByAiInsightId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByAiInsightId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByAiInsightId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByAiInsightId = types.StringNull()
        }
    } else if val, ok := item["triggeredByAiInsightId"].(string); ok {
        data.TriggeredByAiInsightId = types.StringValue(val)
    } else {
        data.TriggeredByAiInsightId = types.StringNull()
    }
    if obj, ok := item["triggeredByAutoRemediationSuggestionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByAutoRemediationSuggestionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByAutoRemediationSuggestionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByAutoRemediationSuggestionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByAutoRemediationSuggestionId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByAutoRemediationSuggestionId = types.StringNull()
        }
    } else if val, ok := item["triggeredByAutoRemediationSuggestionId"].(string); ok {
        data.TriggeredByAutoRemediationSuggestionId = types.StringValue(val)
    } else {
        data.TriggeredByAutoRemediationSuggestionId = types.StringNull()
    }
    if obj, ok := item["monitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := item["aiAgentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAgentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAgentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAgentId = types.StringValue(string(jsonBytes))
        } else {
            data.AiAgentId = types.StringNull()
        }
    } else if val, ok := item["aiAgentId"].(string); ok {
        data.AiAgentId = types.StringValue(val)
    } else {
        data.AiAgentId = types.StringNull()
    }
    if val, ok := item["attemptCount"].(float64); ok {
        data.AttemptCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["attemptCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AttemptCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.AttemptCount = types.NumberNull()
        }
    } else {
        data.AttemptCount = types.NumberNull()
    }
    if obj, ok := item["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartedAt = types.StringNull()
        }
    } else if val, ok := item["startedAt"].(string); ok {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if obj, ok := item["completedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CompletedAt = types.StringNull()
        }
    } else if val, ok := item["completedAt"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    } else {
        data.CompletedAt = types.StringNull()
    }
    if obj, ok := item["lastHeartbeatAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastHeartbeatAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastHeartbeatAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastHeartbeatAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastHeartbeatAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastHeartbeatAt = types.StringNull()
        }
    } else if val, ok := item["lastHeartbeatAt"].(string); ok {
        data.LastHeartbeatAt = types.StringValue(val)
    } else {
        data.LastHeartbeatAt = types.StringNull()
    }
    if val, ok := item["llmCallCount"].(float64); ok {
        data.LlmCallCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["llmCallCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LlmCallCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.LlmCallCount = types.NumberNull()
        }
    } else {
        data.LlmCallCount = types.NumberNull()
    }
    if val, ok := item["toolCallCount"].(float64); ok {
        data.ToolCallCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["toolCallCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ToolCallCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ToolCallCount = types.NumberNull()
        }
    } else {
        data.ToolCallCount = types.NumberNull()
    }
    if val, ok := item["totalTokens"].(float64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["totalTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TotalTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.TotalTokens = types.NumberNull()
        }
    } else {
        data.TotalTokens = types.NumberNull()
    }
    if val, ok := item["totalCostInUSDCents"].(float64); ok {
        data.TotalCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["totalCostInUSDCents"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TotalCostInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.TotalCostInUsdCents = types.NumberNull()
        }
    } else {
        data.TotalCostInUsdCents = types.NumberNull()
    }
    if obj, ok := item["egressManifest"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EgressManifest = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EgressManifest = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EgressManifest = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EgressManifest = types.StringValue(string(jsonBytes))
        } else {
            data.EgressManifest = types.StringNull()
        }
    } else if val, ok := item["egressManifest"].(string); ok {
        data.EgressManifest = types.StringValue(val)
    } else {
        data.EgressManifest = types.StringNull()
    }
    if obj, ok := item["errorMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ErrorMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ErrorMessage = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorMessage = types.StringNull()
        }
    } else if val, ok := item["errorMessage"].(string); ok {
        data.ErrorMessage = types.StringValue(val)
    } else {
        data.ErrorMessage = types.StringNull()
    }
    if obj, ok := item["analysisTldr"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AnalysisTldr = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AnalysisTldr = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AnalysisTldr = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AnalysisTldr = types.StringValue(string(jsonBytes))
        } else {
            data.AnalysisTldr = types.StringNull()
        }
    } else if val, ok := item["analysisTldr"].(string); ok {
        data.AnalysisTldr = types.StringValue(val)
    } else {
        data.AnalysisTldr = types.StringNull()
    }
    if obj, ok := item["humanVerdict"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HumanVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HumanVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HumanVerdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HumanVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.HumanVerdict = types.StringNull()
        }
    } else if val, ok := item["humanVerdict"].(string); ok {
        data.HumanVerdict = types.StringValue(val)
    } else {
        data.HumanVerdict = types.StringNull()
    }
    if obj, ok := item["humanVerdictAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HumanVerdictAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HumanVerdictAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HumanVerdictAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HumanVerdictAt = types.StringValue(string(jsonBytes))
        } else {
            data.HumanVerdictAt = types.StringNull()
        }
    } else if val, ok := item["humanVerdictAt"].(string); ok {
        data.HumanVerdictAt = types.StringValue(val)
    } else {
        data.HumanVerdictAt = types.StringNull()
    }
    if obj, ok := item["humanVerdictByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HumanVerdictByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HumanVerdictByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HumanVerdictByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HumanVerdictByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.HumanVerdictByUserId = types.StringNull()
        }
    } else if val, ok := item["humanVerdictByUserId"].(string); ok {
        data.HumanVerdictByUserId = types.StringValue(val)
    } else {
        data.HumanVerdictByUserId = types.StringNull()
    }
    if obj, ok := item["autoGrade"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoGrade = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoGrade = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoGrade = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoGrade = types.StringValue(string(jsonBytes))
        } else {
            data.AutoGrade = types.StringNull()
        }
    } else if val, ok := item["autoGrade"].(string); ok {
        data.AutoGrade = types.StringValue(val)
    } else {
        data.AutoGrade = types.StringNull()
    }
    if obj, ok := item["autoGradeAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoGradeAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoGradeAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoGradeAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoGradeAt = types.StringValue(string(jsonBytes))
        } else {
            data.AutoGradeAt = types.StringNull()
        }
    } else if val, ok := item["autoGradeAt"].(string); ok {
        data.AutoGradeAt = types.StringValue(val)
    } else {
        data.AutoGradeAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
