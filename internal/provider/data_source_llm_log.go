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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
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
        MarkdownDescription: "Logs of all the LLM API calls for AI features in this project. Look up an existing llm_log by `id` or by `name`.",

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
            "llm_provider_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "llm_provider_name": schema.StringAttribute{
                MarkdownDescription: "Name of the LLM Provider at time of call.",
                Computed: true,
            },
            "llm_type": schema.StringAttribute{
                MarkdownDescription: "Type of LLM (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama).",
                Computed: true,
            },
            "model_name": schema.StringAttribute{
                MarkdownDescription: "Name of the model used (e.g., gpt-4, claude-3-opus).",
                Computed: true,
            },
            "is_global_provider": schema.BoolAttribute{
                MarkdownDescription: "Was a global LLM provider used for this call?.",
                Computed: true,
            },
            "total_tokens": schema.NumberAttribute{
                MarkdownDescription: "Total tokens used (input + output).",
                Computed: true,
            },
            "completion_tokens": schema.NumberAttribute{
                MarkdownDescription: "Output (completion) tokens generated by this call.",
                Computed: true,
            },
            "cached_input_tokens": schema.NumberAttribute{
                MarkdownDescription: "Input tokens served from the provider's prompt cache (billed at a discount).",
                Computed: true,
            },
            "cache_creation_tokens": schema.NumberAttribute{
                MarkdownDescription: "Input tokens written to the provider's prompt cache on this call.",
                Computed: true,
            },
            "cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost in USD cents.",
                Computed: true,
            },
            "was_billed": schema.BoolAttribute{
                MarkdownDescription: "Was the project charged for this API call?.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the LLM API call.",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (error details if failed).",
                Computed: true,
            },
            "feature": schema.StringAttribute{
                MarkdownDescription: "The feature that triggered this API call (e.g., IncidentPostmortem).",
                Computed: true,
            },
            "request_prompt": schema.StringAttribute{
                MarkdownDescription: "The prompt sent to the LLM (truncated).",
                Computed: true,
            },
            "response_preview": schema.StringAttribute{
                MarkdownDescription: "Preview of the LLM response (truncated).",
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
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "request_started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "request_completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "duration_ms": schema.NumberAttribute{
                MarkdownDescription: "Request duration in milliseconds.",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a llm_log.",
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm_log found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm_log found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one llm_log matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
