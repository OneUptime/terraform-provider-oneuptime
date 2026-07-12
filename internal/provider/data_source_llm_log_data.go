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
var _ datasource.DataSource = &LlmLogDataDataSource{}

func NewLlmLogDataDataSource() datasource.DataSource {
    return &LlmLogDataDataSource{}
}

// LlmLogDataDataSource defines the data source implementation.
type LlmLogDataDataSource struct {
    client *Client
}

// LlmLogDataDataSourceModel describes the data source data model.
type LlmLogDataDataSourceModel struct {
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

func (d *LlmLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_llm_log_data"
}

func (d *LlmLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "llm_log_data data source",

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
            "llm_provider_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "llm_provider_name": schema.StringAttribute{
                MarkdownDescription: "Name of the LLM Provider at time of call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "llm_type": schema.StringAttribute{
                MarkdownDescription: "Type of LLM (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "model_name": schema.StringAttribute{
                MarkdownDescription: "Name of the model used (e.g., gpt-4, claude-3-opus). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_global_provider": schema.BoolAttribute{
                MarkdownDescription: "Was a global LLM provider used for this call?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "total_tokens": schema.NumberAttribute{
                MarkdownDescription: "Total tokens used (input + output). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cached_input_tokens": schema.NumberAttribute{
                MarkdownDescription: "Input tokens served from the provider's prompt cache (billed at a discount). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cache_creation_tokens": schema.NumberAttribute{
                MarkdownDescription: "Input tokens written to the provider's prompt cache on this call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "was_billed": schema.BoolAttribute{
                MarkdownDescription: "Was the project charged for this API call?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the LLM API call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (error details if failed). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "feature": schema.StringAttribute{
                MarkdownDescription: "The feature that triggered this API call (e.g., IncidentPostmortem). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "request_prompt": schema.StringAttribute{
                MarkdownDescription: "The prompt sent to the LLM (truncated). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "response_preview": schema.StringAttribute{
                MarkdownDescription: "Preview of the LLM response (truncated). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "Request duration in milliseconds. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *LlmLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LlmLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LlmLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "llm-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read llm_log_data, got error: %s", err))
        return
    }

    var llmLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &llmLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse llm_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := llmLogDataResponse["data"].(map[string]interface{}); ok {
        llmLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := llmLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := llmLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["llm_provider_id"].(string); ok {
        data.LlmProviderId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["llm_provider_name"].(string); ok {
        data.LlmProviderName = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["llm_type"].(string); ok {
        data.LlmType = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["model_name"].(string); ok {
        data.ModelName = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["is_global_provider"].(bool); ok {
        data.IsGlobalProvider = types.BoolValue(val)
    }
    if val, ok := llmLogDataResponse["total_tokens"].(float64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := llmLogDataResponse["cached_input_tokens"].(float64); ok {
        data.CachedInputTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := llmLogDataResponse["cache_creation_tokens"].(float64); ok {
        data.CacheCreationTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := llmLogDataResponse["cost_in_usd_cents"].(float64); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := llmLogDataResponse["was_billed"].(bool); ok {
        data.WasBilled = types.BoolValue(val)
    }
    if val, ok := llmLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["feature"].(string); ok {
        data.Feature = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["request_prompt"].(string); ok {
        data.RequestPrompt = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["response_preview"].(string); ok {
        data.ResponsePreview = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["incident_id"].(string); ok {
        data.IncidentId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["ai_run_id"].(string); ok {
        data.AiRunId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["request_started_at"].(string); ok {
        data.RequestStartedAt = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["request_completed_at"].(string); ok {
        data.RequestCompletedAt = types.StringValue(val)
    }
    if val, ok := llmLogDataResponse["duration_ms"].(float64); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
