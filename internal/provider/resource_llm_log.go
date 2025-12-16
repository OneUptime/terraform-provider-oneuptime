package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "encoding/json"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LlmLogResource{}
var _ resource.ResourceWithImportState = &LlmLogResource{}

func NewLlmLogResource() resource.Resource {
    return &LlmLogResource{}
}

// LlmLogResource defines the resource implementation.
type LlmLogResource struct {
    client *Client
}

// LlmLogResourceModel describes the resource data model.
type LlmLogResourceModel struct {
    Id types.String `tfsdk:"id"`
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
    CostInUsdCents types.Number `tfsdk:"cost_in_usd_cents"`
    WasBilled types.Bool `tfsdk:"was_billed"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    Feature types.String `tfsdk:"feature"`
    RequestPrompt types.String `tfsdk:"request_prompt"`
    ResponsePreview types.String `tfsdk:"response_preview"`
    IncidentId types.String `tfsdk:"incident_id"`
    AlertId types.String `tfsdk:"alert_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    UserId types.String `tfsdk:"user_id"`
    RequestStartedAt types.String `tfsdk:"request_started_at"`
    RequestCompletedAt types.String `tfsdk:"request_completed_at"`
    DurationMs types.Number `tfsdk:"duration_ms"`
}

func (r *LlmLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_llm_log"
}

func (r *LlmLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "llm_log resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
                MarkdownDescription: "Name of the LLM Provider at time of call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "llm_type": schema.StringAttribute{
                MarkdownDescription: "Type of LLM (OpenAI, Anthropic, Ollama). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "model_name": schema.StringAttribute{
                MarkdownDescription: "Name of the model used (e.g., gpt-4, claude-3-opus). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_global_provider": schema.BoolAttribute{
                MarkdownDescription: "Was a global LLM provider used for this call?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "total_tokens": schema.NumberAttribute{
                MarkdownDescription: "Total tokens used (input + output). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "was_billed": schema.BoolAttribute{
                MarkdownDescription: "Was the project charged for this API call?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the LLM API call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message (error details if failed). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "feature": schema.StringAttribute{
                MarkdownDescription: "The feature that triggered this API call (e.g., IncidentPostmortem). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "request_prompt": schema.StringAttribute{
                MarkdownDescription: "The prompt sent to the LLM (truncated). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "response_preview": schema.StringAttribute{
                MarkdownDescription: "Preview of the LLM response (truncated). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "Request duration in milliseconds. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read LLM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *LlmLogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *LlmLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data LlmLogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    llmLogRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/llm-log/count", llmLogRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create llm_log, got error: %s", err))
        return
    }

    var llmLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &llmLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse llm_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := llmLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = llmLogResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["llmProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderId"].(string); ok && val != "" {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := dataMap["llmProviderName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.LlmProviderName = types.StringValue(val)
        } else {
            data.LlmProviderName = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderName"].(string); ok && val != "" {
        data.LlmProviderName = types.StringValue(val)
    } else {
        data.LlmProviderName = types.StringNull()
    }
    if obj, ok := dataMap["llmType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.LlmType = types.StringValue(val)
        } else {
            data.LlmType = types.StringNull()
        }
    } else if val, ok := dataMap["llmType"].(string); ok && val != "" {
        data.LlmType = types.StringValue(val)
    } else {
        data.LlmType = types.StringNull()
    }
    if obj, ok := dataMap["modelName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ModelName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ModelName = types.StringValue(val)
        } else {
            data.ModelName = types.StringNull()
        }
    } else if val, ok := dataMap["modelName"].(string); ok && val != "" {
        data.ModelName = types.StringValue(val)
    } else {
        data.ModelName = types.StringNull()
    }
    if val, ok := dataMap["isGlobalProvider"].(bool); ok {
        data.IsGlobalProvider = types.BoolValue(val)
    }
    if val, ok := dataMap["totalTokens"].(float64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["totalTokens"].(int); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["totalTokens"].(int64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["totalTokens"] == nil {
        data.TotalTokens = types.NumberNull()
    }
    if val, ok := dataMap["costInUSDCents"].(float64); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["costInUSDCents"].(int); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["costInUSDCents"].(int64); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["costInUSDCents"] == nil {
        data.CostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["wasBilled"].(bool); ok {
        data.WasBilled = types.BoolValue(val)
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["feature"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Feature = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Feature = types.StringValue(val)
        } else {
            data.Feature = types.StringNull()
        }
    } else if val, ok := dataMap["feature"].(string); ok && val != "" {
        data.Feature = types.StringValue(val)
    } else {
        data.Feature = types.StringNull()
    }
    if obj, ok := dataMap["requestPrompt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestPrompt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RequestPrompt = types.StringValue(val)
        } else {
            data.RequestPrompt = types.StringNull()
        }
    } else if val, ok := dataMap["requestPrompt"].(string); ok && val != "" {
        data.RequestPrompt = types.StringValue(val)
    } else {
        data.RequestPrompt = types.StringNull()
    }
    if obj, ok := dataMap["responsePreview"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResponsePreview = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ResponsePreview = types.StringValue(val)
        } else {
            data.ResponsePreview = types.StringNull()
        }
    } else if val, ok := dataMap["responsePreview"].(string); ok && val != "" {
        data.ResponsePreview = types.StringValue(val)
    } else {
        data.ResponsePreview = types.StringNull()
    }
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok && val != "" {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["alertId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := dataMap["alertId"].(string); ok && val != "" {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if obj, ok := dataMap["scheduledMaintenanceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceId"].(string); ok && val != "" {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
    }
    if obj, ok := dataMap["userId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["requestStartedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RequestStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestStartedAt = types.StringNull()
        }
    } else if val, ok := dataMap["requestStartedAt"].(string); ok && val != "" {
        data.RequestStartedAt = types.StringValue(val)
    } else {
        data.RequestStartedAt = types.StringNull()
    }
    if val, ok := dataMap["requestCompletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RequestCompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestCompletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["requestCompletedAt"].(string); ok && val != "" {
        data.RequestCompletedAt = types.StringValue(val)
    } else {
        data.RequestCompletedAt = types.StringNull()
    }
    if val, ok := dataMap["durationMs"].(float64); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["durationMs"].(int); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["durationMs"].(int64); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["durationMs"] == nil {
        data.DurationMs = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LlmLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data LlmLogResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
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
        "costInUSDCents": true,
        "wasBilled": true,
        "status": true,
        "statusMessage": true,
        "feature": true,
        "requestPrompt": true,
        "responsePreview": true,
        "incidentId": true,
        "alertId": true,
        "scheduledMaintenanceId": true,
        "userId": true,
        "requestStartedAt": true,
        "requestCompletedAt": true,
        "durationMs": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/llm-log/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read llm_log, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var llmLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &llmLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse llm_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := llmLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = llmLogResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["llmProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderId"].(string); ok && val != "" {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := dataMap["llmProviderName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.LlmProviderName = types.StringValue(val)
        } else {
            data.LlmProviderName = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderName"].(string); ok && val != "" {
        data.LlmProviderName = types.StringValue(val)
    } else {
        data.LlmProviderName = types.StringNull()
    }
    if obj, ok := dataMap["llmType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.LlmType = types.StringValue(val)
        } else {
            data.LlmType = types.StringNull()
        }
    } else if val, ok := dataMap["llmType"].(string); ok && val != "" {
        data.LlmType = types.StringValue(val)
    } else {
        data.LlmType = types.StringNull()
    }
    if obj, ok := dataMap["modelName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ModelName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ModelName = types.StringValue(val)
        } else {
            data.ModelName = types.StringNull()
        }
    } else if val, ok := dataMap["modelName"].(string); ok && val != "" {
        data.ModelName = types.StringValue(val)
    } else {
        data.ModelName = types.StringNull()
    }
    if val, ok := dataMap["isGlobalProvider"].(bool); ok {
        data.IsGlobalProvider = types.BoolValue(val)
    }
    if val, ok := dataMap["totalTokens"].(float64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["totalTokens"].(int); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["totalTokens"].(int64); ok {
        data.TotalTokens = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["totalTokens"] == nil {
        data.TotalTokens = types.NumberNull()
    }
    if val, ok := dataMap["costInUSDCents"].(float64); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["costInUSDCents"].(int); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["costInUSDCents"].(int64); ok {
        data.CostInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["costInUSDCents"] == nil {
        data.CostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["wasBilled"].(bool); ok {
        data.WasBilled = types.BoolValue(val)
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["feature"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Feature = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Feature = types.StringValue(val)
        } else {
            data.Feature = types.StringNull()
        }
    } else if val, ok := dataMap["feature"].(string); ok && val != "" {
        data.Feature = types.StringValue(val)
    } else {
        data.Feature = types.StringNull()
    }
    if obj, ok := dataMap["requestPrompt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestPrompt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.RequestPrompt = types.StringValue(val)
        } else {
            data.RequestPrompt = types.StringNull()
        }
    } else if val, ok := dataMap["requestPrompt"].(string); ok && val != "" {
        data.RequestPrompt = types.StringValue(val)
    } else {
        data.RequestPrompt = types.StringNull()
    }
    if obj, ok := dataMap["responsePreview"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResponsePreview = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ResponsePreview = types.StringValue(val)
        } else {
            data.ResponsePreview = types.StringNull()
        }
    } else if val, ok := dataMap["responsePreview"].(string); ok && val != "" {
        data.ResponsePreview = types.StringValue(val)
    } else {
        data.ResponsePreview = types.StringNull()
    }
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok && val != "" {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["alertId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := dataMap["alertId"].(string); ok && val != "" {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if obj, ok := dataMap["scheduledMaintenanceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceId"].(string); ok && val != "" {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
    }
    if obj, ok := dataMap["userId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["requestStartedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RequestStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestStartedAt = types.StringNull()
        }
    } else if val, ok := dataMap["requestStartedAt"].(string); ok && val != "" {
        data.RequestStartedAt = types.StringValue(val)
    } else {
        data.RequestStartedAt = types.StringNull()
    }
    if val, ok := dataMap["requestCompletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RequestCompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RequestCompletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["requestCompletedAt"].(string); ok && val != "" {
        data.RequestCompletedAt = types.StringValue(val)
    } else {
        data.RequestCompletedAt = types.StringNull()
    }
    if val, ok := dataMap["durationMs"].(float64); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["durationMs"].(int); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["durationMs"].(int64); ok {
        data.DurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["durationMs"] == nil {
        data.DurationMs = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LlmLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *LlmLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *LlmLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *LlmLogResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *LlmLogResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)
    
    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to parse JSON field for complex objects
func (r *LlmLogResource) parseJSONField(terraformString types.String) interface{} {
    if terraformString.IsNull() || terraformString.IsUnknown() || terraformString.ValueString() == "" {
        return nil
    }
    
    var result interface{}
    if err := json.Unmarshal([]byte(terraformString.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return terraformString.ValueString()
    }
    
    return result
}
