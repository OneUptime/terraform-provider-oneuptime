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
var _ datasource.DataSource = &LlmLogDataSource{}

func NewLlmLogDataSource() datasource.DataSource {
    return &LlmLogDataSource{}
}

// LlmLogDataSource defines the data source implementation.
type LlmLogDataSource struct {
    client *Client
}

// LlmLogDataSourceModel describes the data source data model.
type LlmLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    LlmProviderId types.String `tfsdk:"llm_provider_id"`
    LlmProviderName types.String `tfsdk:"llm_provider_name"`
    LlmType types.String `tfsdk:"llm_type"`
    ModelName types.String `tfsdk:"model_name"`
    IsGlobalProvider types.Bool `tfsdk:"is_global_provider"`
    TotalTokens types.Number `tfsdk:"total_tokens"`
    CompletionTokens types.Number `tfsdk:"completion_tokens"`
    CachedInputTokens types.Number `tfsdk:"cached_input_tokens"`
    CacheCreationTokens types.Number `tfsdk:"cache_creation_tokens"`
    CostInUsdCents types.Number `tfsdk:"cost_in_usd_cents"`
    WasBilled types.Bool `tfsdk:"was_billed"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    Feature types.String `tfsdk:"feature"`
    RequestPrompt types.String `tfsdk:"request_prompt"`
    ResponsePreview types.String `tfsdk:"response_preview"`
    IncidentId types.String `tfsdk:"incident_id"`
    AlertId types.String `tfsdk:"alert_id"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    UserId types.String `tfsdk:"user_id"`
    RequestStartedAt types.String `tfsdk:"request_started_at"`
    RequestCompletedAt types.String `tfsdk:"request_completed_at"`
    DurationMs types.Number `tfsdk:"duration_ms"`
}

func (d *LlmLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_llm_log"
}

func (d *LlmLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Logs of all the LLM API calls for AI features in this project. Look up an existing llm log by `id`, or by any of its other arguments (`ai_run_id`, `alert_id`, `cache_creation_tokens`, ...): each one set must match, and exactly one llm log may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "llm_provider_id": schema.StringAttribute{
                MarkdownDescription: "ID of LLM Provider used for this API call. The ID of a `oneuptime_llm_provider`.",
                Optional: true,
                Computed: true,
            },
            "llm_provider_name": schema.StringAttribute{
                MarkdownDescription: "Name of the LLM Provider at time of call.",
                Optional: true,
                Computed: true,
            },
            "llm_type": schema.StringAttribute{
                MarkdownDescription: "Type of LLM (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama).",
                Optional: true,
                Computed: true,
            },
            "model_name": schema.StringAttribute{
                MarkdownDescription: "Name of the model used (e.g., gpt-4, claude-3-opus).",
                Optional: true,
                Computed: true,
            },
            "is_global_provider": schema.BoolAttribute{
                MarkdownDescription: "Was a global LLM provider used for this call?",
                Optional: true,
                Computed: true,
            },
            "total_tokens": schema.NumberAttribute{
                MarkdownDescription: "Total tokens used (input + output).",
                Optional: true,
                Computed: true,
            },
            "completion_tokens": schema.NumberAttribute{
                MarkdownDescription: "Output (completion) tokens generated by this call.",
                Optional: true,
                Computed: true,
            },
            "cached_input_tokens": schema.NumberAttribute{
                MarkdownDescription: "Input tokens served from the provider's prompt cache (billed at a discount).",
                Optional: true,
                Computed: true,
            },
            "cache_creation_tokens": schema.NumberAttribute{
                MarkdownDescription: "Input tokens written to the provider's prompt cache on this call.",
                Optional: true,
                Computed: true,
            },
            "cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost in USD cents.",
                Optional: true,
                Computed: true,
            },
            "was_billed": schema.BoolAttribute{
                MarkdownDescription: "Was the project charged for this API call?",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the LLM API call.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (error details if failed).",
                Optional: true,
                Computed: true,
            },
            "feature": schema.StringAttribute{
                MarkdownDescription: "The feature that triggered this API call (e.g., IncidentPostmortem).",
                Optional: true,
                Computed: true,
            },
            "request_prompt": schema.StringAttribute{
                MarkdownDescription: "The prompt sent to the LLM (truncated).",
                Optional: true,
                Computed: true,
            },
            "response_preview": schema.StringAttribute{
                MarkdownDescription: "Preview of the LLM response (truncated).",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of Incident associated with this LLM call (if any). The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of Alert associated with this LLM call (if any). The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the AI run this LLM call was part of (if any).",
                Optional: true,
                Computed: true,
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "ID of Scheduled Maintenance associated with this LLM call (if any). The ID of a `oneuptime_scheduled_maintenance_event`.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "ID of User who triggered this LLM call (if any). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "request_started_at": schema.StringAttribute{
                MarkdownDescription: "When the LLM request started.",
                Computed: true,
            },
            "request_completed_at": schema.StringAttribute{
                MarkdownDescription: "When the LLM request completed.",
                Computed: true,
            },
            "duration_ms": schema.NumberAttribute{
                MarkdownDescription: "Request duration in milliseconds.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *LlmLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LlmLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LlmLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.LlmProviderId.IsNull() && !data.LlmProviderId.IsUnknown() {
        filters["llmProviderId"] = data.LlmProviderId.ValueString()
        filterNames = append(filterNames, "llm_provider_id = "+fmt.Sprintf("%q", data.LlmProviderId.ValueString()))
    }
    if !data.LlmProviderName.IsNull() && !data.LlmProviderName.IsUnknown() {
        filters["llmProviderName"] = data.LlmProviderName.ValueString()
        filterNames = append(filterNames, "llm_provider_name = "+fmt.Sprintf("%q", data.LlmProviderName.ValueString()))
    }
    if !data.LlmType.IsNull() && !data.LlmType.IsUnknown() {
        filters["llmType"] = data.LlmType.ValueString()
        filterNames = append(filterNames, "llm_type = "+fmt.Sprintf("%q", data.LlmType.ValueString()))
    }
    if !data.ModelName.IsNull() && !data.ModelName.IsUnknown() {
        filters["modelName"] = data.ModelName.ValueString()
        filterNames = append(filterNames, "model_name = "+fmt.Sprintf("%q", data.ModelName.ValueString()))
    }
    if !data.IsGlobalProvider.IsNull() && !data.IsGlobalProvider.IsUnknown() {
        filters["isGlobalProvider"] = data.IsGlobalProvider.ValueBool()
        filterNames = append(filterNames, "is_global_provider = "+fmt.Sprintf("%t", data.IsGlobalProvider.ValueBool()))
    }
    if !data.TotalTokens.IsNull() && !data.TotalTokens.IsUnknown() {
        filters["totalTokens"] = lookupNumber(data.TotalTokens)
        filterNames = append(filterNames, "total_tokens = "+data.TotalTokens.ValueBigFloat().String())
    }
    if !data.CompletionTokens.IsNull() && !data.CompletionTokens.IsUnknown() {
        filters["completionTokens"] = lookupNumber(data.CompletionTokens)
        filterNames = append(filterNames, "completion_tokens = "+data.CompletionTokens.ValueBigFloat().String())
    }
    if !data.CachedInputTokens.IsNull() && !data.CachedInputTokens.IsUnknown() {
        filters["cachedInputTokens"] = lookupNumber(data.CachedInputTokens)
        filterNames = append(filterNames, "cached_input_tokens = "+data.CachedInputTokens.ValueBigFloat().String())
    }
    if !data.CacheCreationTokens.IsNull() && !data.CacheCreationTokens.IsUnknown() {
        filters["cacheCreationTokens"] = lookupNumber(data.CacheCreationTokens)
        filterNames = append(filterNames, "cache_creation_tokens = "+data.CacheCreationTokens.ValueBigFloat().String())
    }
    if !data.CostInUsdCents.IsNull() && !data.CostInUsdCents.IsUnknown() {
        filters["costInUSDCents"] = lookupNumber(data.CostInUsdCents)
        filterNames = append(filterNames, "cost_in_usd_cents = "+data.CostInUsdCents.ValueBigFloat().String())
    }
    if !data.WasBilled.IsNull() && !data.WasBilled.IsUnknown() {
        filters["wasBilled"] = data.WasBilled.ValueBool()
        filterNames = append(filterNames, "was_billed = "+fmt.Sprintf("%t", data.WasBilled.ValueBool()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.Feature.IsNull() && !data.Feature.IsUnknown() {
        filters["feature"] = data.Feature.ValueString()
        filterNames = append(filterNames, "feature = "+fmt.Sprintf("%q", data.Feature.ValueString()))
    }
    if !data.RequestPrompt.IsNull() && !data.RequestPrompt.IsUnknown() {
        filters["requestPrompt"] = data.RequestPrompt.ValueString()
        filterNames = append(filterNames, "request_prompt = "+fmt.Sprintf("%q", data.RequestPrompt.ValueString()))
    }
    if !data.ResponsePreview.IsNull() && !data.ResponsePreview.IsUnknown() {
        filters["responsePreview"] = data.ResponsePreview.ValueString()
        filterNames = append(filterNames, "response_preview = "+fmt.Sprintf("%q", data.ResponsePreview.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.AlertId.IsNull() && !data.AlertId.IsUnknown() {
        filters["alertId"] = data.AlertId.ValueString()
        filterNames = append(filterNames, "alert_id = "+fmt.Sprintf("%q", data.AlertId.ValueString()))
    }
    if !data.AiRunId.IsNull() && !data.AiRunId.IsUnknown() {
        filters["aiRunId"] = data.AiRunId.ValueString()
        filterNames = append(filterNames, "ai_run_id = "+fmt.Sprintf("%q", data.AiRunId.ValueString()))
    }
    if !data.ScheduledMaintenanceId.IsNull() && !data.ScheduledMaintenanceId.IsUnknown() {
        filters["scheduledMaintenanceId"] = data.ScheduledMaintenanceId.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_id = "+fmt.Sprintf("%q", data.ScheduledMaintenanceId.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.DurationMs.IsNull() && !data.DurationMs.IsUnknown() {
        filters["durationMs"] = lookupNumber(data.DurationMs)
        filterNames = append(filterNames, "duration_ms = "+data.DurationMs.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the llm log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the llm log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "llmProviderId": true,
        "llmProviderName": true,
        "llmType": true,
        "modelName": true,
        "isGlobalProvider": true,
        "totalTokens": true,
        "completionTokens": true,
        "cachedInputTokens": true,
        "cacheCreationTokens": true,
        "costInUSDCents": true,
        "wasBilled": true,
        "status": true,
        "statusMessage": true,
        "feature": true,
        "requestPrompt": true,
        "responsePreview": true,
        "incidentId": true,
        "alertId": true,
        "aiRunId": true,
        "scheduledMaintenanceId": true,
        "userId": true,
        "requestStartedAt": true,
        "requestCompletedAt": true,
        "durationMs": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/llm-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read llm_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read llm_log: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/llm-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list llm_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list llm_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one llm log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for llm_log.")
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
    if obj, ok := item["llmProviderId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := item["llmProviderId"].(string); ok {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := item["llmProviderName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmProviderName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmProviderName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmProviderName = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderName = types.StringNull()
        }
    } else if val, ok := item["llmProviderName"].(string); ok {
        data.LlmProviderName = types.StringValue(val)
    } else {
        data.LlmProviderName = types.StringNull()
    }
    if obj, ok := item["llmType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmType = types.StringValue(string(jsonBytes))
        } else {
            data.LlmType = types.StringNull()
        }
    } else if val, ok := item["llmType"].(string); ok {
        data.LlmType = types.StringValue(val)
    } else {
        data.LlmType = types.StringNull()
    }
    if obj, ok := item["modelName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ModelName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ModelName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ModelName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ModelName = types.StringValue(string(jsonBytes))
        } else {
            data.ModelName = types.StringNull()
        }
    } else if val, ok := item["modelName"].(string); ok {
        data.ModelName = types.StringValue(val)
    } else {
        data.ModelName = types.StringNull()
    }
    if val, ok := item["isGlobalProvider"].(bool); ok {
        data.IsGlobalProvider = types.BoolValue(val)
    } else {
        data.IsGlobalProvider = types.BoolNull()
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
    if val, ok := item["completionTokens"].(float64); ok {
        data.CompletionTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["completionTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CompletionTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.CompletionTokens = types.NumberNull()
        }
    } else {
        data.CompletionTokens = types.NumberNull()
    }
    if val, ok := item["cachedInputTokens"].(float64); ok {
        data.CachedInputTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cachedInputTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CachedInputTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.CachedInputTokens = types.NumberNull()
        }
    } else {
        data.CachedInputTokens = types.NumberNull()
    }
    if val, ok := item["cacheCreationTokens"].(float64); ok {
        data.CacheCreationTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cacheCreationTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CacheCreationTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.CacheCreationTokens = types.NumberNull()
        }
    } else {
        data.CacheCreationTokens = types.NumberNull()
    }
    if val, ok := item["costInUSDCents"].(float64); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["costInUSDCents"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CostInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.CostInUsdCents = types.NumberNull()
        }
    } else {
        data.CostInUsdCents = types.NumberNull()
    }
    if val, ok := item["wasBilled"].(bool); ok {
        data.WasBilled = types.BoolValue(val)
    } else {
        data.WasBilled = types.BoolNull()
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
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["feature"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Feature = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Feature = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Feature = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Feature = types.StringValue(string(jsonBytes))
        } else {
            data.Feature = types.StringNull()
        }
    } else if val, ok := item["feature"].(string); ok {
        data.Feature = types.StringValue(val)
    } else {
        data.Feature = types.StringNull()
    }
    if obj, ok := item["requestPrompt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestPrompt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RequestPrompt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RequestPrompt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RequestPrompt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestPrompt = types.StringNull()
        }
    } else if val, ok := item["requestPrompt"].(string); ok {
        data.RequestPrompt = types.StringValue(val)
    } else {
        data.RequestPrompt = types.StringNull()
    }
    if obj, ok := item["responsePreview"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResponsePreview = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResponsePreview = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResponsePreview = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResponsePreview = types.StringValue(string(jsonBytes))
        } else {
            data.ResponsePreview = types.StringNull()
        }
    } else if val, ok := item["responsePreview"].(string); ok {
        data.ResponsePreview = types.StringValue(val)
    } else {
        data.ResponsePreview = types.StringNull()
    }
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := item["alertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := item["alertId"].(string); ok {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if obj, ok := item["aiRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.AiRunId = types.StringNull()
        }
    } else if val, ok := item["aiRunId"].(string); ok {
        data.AiRunId = types.StringValue(val)
    } else {
        data.AiRunId = types.StringNull()
    }
    if obj, ok := item["scheduledMaintenanceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
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
    if obj, ok := item["requestStartedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RequestStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RequestStartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RequestStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestStartedAt = types.StringNull()
        }
    } else if val, ok := item["requestStartedAt"].(string); ok {
        data.RequestStartedAt = types.StringValue(val)
    } else {
        data.RequestStartedAt = types.StringNull()
    }
    if obj, ok := item["requestCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RequestCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RequestCompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RequestCompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestCompletedAt = types.StringNull()
        }
    } else if val, ok := item["requestCompletedAt"].(string); ok {
        data.RequestCompletedAt = types.StringValue(val)
    } else {
        data.RequestCompletedAt = types.StringNull()
    }
    if val, ok := item["durationMs"].(float64); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["durationMs"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DurationMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.DurationMs = types.NumberNull()
        }
    } else {
        data.DurationMs = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
