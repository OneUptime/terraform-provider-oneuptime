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
var _ datasource.DataSource = &LlmProviderDataSource{}

func NewLlmProviderDataSource() datasource.DataSource {
    return &LlmProviderDataSource{}
}

// LlmProviderDataSource defines the data source implementation.
type LlmProviderDataSource struct {
    client *Client
}

// LlmProviderDataSourceModel describes the data source data model.
type LlmProviderDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    LlmType types.String `tfsdk:"llm_type"`
    ApiKey types.String `tfsdk:"api_key"`
    ModelName types.String `tfsdk:"model_name"`
    BaseUrl types.String `tfsdk:"base_url"`
    AdditionalParams types.String `tfsdk:"additional_params"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsDefault types.Bool `tfsdk:"is_default"`
    CostPerMillionTokensInUsdCents types.Number `tfsdk:"cost_per_million_tokens_in_usd_cents"`
}

func (d *LlmProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_llm_provider"
}

func (d *LlmProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features. Look up an existing llm provider by `id`, or by any of its other arguments (`name`, `api_key`, `base_url`, ...): each one set must match, and exactly one llm provider may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "A friendly name for this LLM configuration.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this LLM configuration.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "llm_type": schema.StringAttribute{
                MarkdownDescription: "The type of LLM provider (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAICompatible, etc.).",
                Optional: true,
                Computed: true,
            },
            "api_key": schema.StringAttribute{
                MarkdownDescription: "The API key for the LLM provider. Required for OpenAI, Azure OpenAI, Anthropic, Groq, and Mistral.",
                Optional: true,
                Computed: true,
            },
            "model_name": schema.StringAttribute{
                MarkdownDescription: "The name of the model to use (e.g., gpt-4, claude-3-opus, llama2).",
                Optional: true,
                Computed: true,
            },
            "base_url": schema.StringAttribute{
                MarkdownDescription: "The base URL for the LLM API. Required for Azure OpenAI and Ollama, optional for others.",
                Optional: true,
                Computed: true,
            },
            "additional_params": schema.StringAttribute{
                MarkdownDescription: "Optional JSON object with extra parameters sent directly to the provider API. These are merged last and override any defaults. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this LLM belongs to. If null, it is a global LLM.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User).",
                Optional: true,
                Computed: true,
            },
            "is_default": schema.BoolAttribute{
                MarkdownDescription: "Is this the default LLM provider for the project? When set, the global LLM provider will not be used.",
                Optional: true,
                Computed: true,
            },
            "cost_per_million_tokens_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost per million tokens in USD cents. Used for billing when using global LLM providers.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *LlmProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LlmProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LlmProviderDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.LlmType.IsNull() && !data.LlmType.IsUnknown() {
        filters["llmType"] = data.LlmType.ValueString()
        filterNames = append(filterNames, "llm_type = "+fmt.Sprintf("%q", data.LlmType.ValueString()))
    }
    if !data.ApiKey.IsNull() && !data.ApiKey.IsUnknown() {
        filters["apiKey"] = data.ApiKey.ValueString()
        filterNames = append(filterNames, "api_key = "+fmt.Sprintf("%q", data.ApiKey.ValueString()))
    }
    if !data.ModelName.IsNull() && !data.ModelName.IsUnknown() {
        filters["modelName"] = data.ModelName.ValueString()
        filterNames = append(filterNames, "model_name = "+fmt.Sprintf("%q", data.ModelName.ValueString()))
    }
    if !data.BaseUrl.IsNull() && !data.BaseUrl.IsUnknown() {
        filters["baseUrl"] = data.BaseUrl.ValueString()
        filterNames = append(filterNames, "base_url = "+fmt.Sprintf("%q", data.BaseUrl.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() {
        filters["isDefault"] = data.IsDefault.ValueBool()
        filterNames = append(filterNames, "is_default = "+fmt.Sprintf("%t", data.IsDefault.ValueBool()))
    }
    if !data.CostPerMillionTokensInUsdCents.IsNull() && !data.CostPerMillionTokensInUsdCents.IsUnknown() {
        filters["costPerMillionTokensInUSDCents"] = lookupNumber(data.CostPerMillionTokensInUsdCents)
        filterNames = append(filterNames, "cost_per_million_tokens_in_usd_cents = "+data.CostPerMillionTokensInUsdCents.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the llm provider up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the llm provider up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "name": true,
        "description": true,
        "slug": true,
        "llmType": true,
        "apiKey": true,
        "modelName": true,
        "baseUrl": true,
        "additionalParams": true,
        "projectId": true,
        "createdByUserId": true,
        "isDefault": true,
        "costPerMillionTokensInUSDCents": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/llm-provider/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read llm_provider, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm provider found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read llm_provider: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/llm-provider/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list llm_provider, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list llm_provider: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm provider matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one llm provider matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for llm_provider.")
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
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := item["apiKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApiKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.ApiKey = types.StringNull()
        }
    } else if val, ok := item["apiKey"].(string); ok {
        data.ApiKey = types.StringValue(val)
    } else {
        data.ApiKey = types.StringNull()
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
    if obj, ok := item["baseUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BaseUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BaseUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BaseUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BaseUrl = types.StringValue(string(jsonBytes))
        } else {
            data.BaseUrl = types.StringNull()
        }
    } else if val, ok := item["baseUrl"].(string); ok {
        data.BaseUrl = types.StringValue(val)
    } else {
        data.BaseUrl = types.StringNull()
    }
    if obj, ok := item["additionalParams"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AdditionalParams = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AdditionalParams = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AdditionalParams = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AdditionalParams = types.StringValue(string(jsonBytes))
        } else {
            data.AdditionalParams = types.StringNull()
        }
    } else if val, ok := item["additionalParams"].(string); ok {
        data.AdditionalParams = types.StringValue(val)
    } else {
        data.AdditionalParams = types.StringNull()
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
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["isDefault"].(bool); ok {
        data.IsDefault = types.BoolValue(val)
    } else {
        data.IsDefault = types.BoolNull()
    }
    if val, ok := item["costPerMillionTokensInUSDCents"].(float64); ok {
        data.CostPerMillionTokensInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["costPerMillionTokensInUSDCents"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CostPerMillionTokensInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.CostPerMillionTokensInUsdCents = types.NumberNull()
        }
    } else {
        data.CostPerMillionTokensInUsdCents = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
