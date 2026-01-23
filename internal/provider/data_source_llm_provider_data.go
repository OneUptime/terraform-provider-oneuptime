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
var _ datasource.DataSource = &LlmProviderDataDataSource{}

func NewLlmProviderDataDataSource() datasource.DataSource {
    return &LlmProviderDataDataSource{}
}

// LlmProviderDataDataSource defines the data source implementation.
type LlmProviderDataDataSource struct {
    client *Client
}

// LlmProviderDataDataSourceModel describes the data source data model.
type LlmProviderDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    LlmType types.String `tfsdk:"llm_type"`
    ApiKey types.String `tfsdk:"api_key"`
    ModelName types.String `tfsdk:"model_name"`
    BaseUrl types.String `tfsdk:"base_url"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsDefault types.Bool `tfsdk:"is_default"`
    CostPerMillionTokensInUsdCents types.Number `tfsdk:"cost_per_million_tokens_in_usd_cents"`
}

func (d *LlmProviderDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_llm_provider_data"
}

func (d *LlmProviderDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "llm_provider_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this LLM configuration.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create LLM], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit LLM]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Public], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "llm_type": schema.StringAttribute{
                MarkdownDescription: "The type of LLM provider (OpenAI, Anthropic, Ollama, etc.). Permissions - Create: [Project Owner, Project Admin, Project Member, Create LLM], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit LLM]",
                Computed: true,
            },
            "api_key": schema.StringAttribute{
                MarkdownDescription: "The API key for the LLM provider. Required for OpenAI and Anthropic.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create LLM], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Edit LLM]",
                Computed: true,
            },
            "model_name": schema.StringAttribute{
                MarkdownDescription: "The name of the model to use (e.g., gpt-4, claude-3-opus, llama2).. Permissions - Create: [Project Owner, Project Admin, Project Member, Create LLM], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit LLM]",
                Computed: true,
            },
            "base_url": schema.StringAttribute{
                MarkdownDescription: "The base URL for the LLM API. Required for Ollama, optional for others.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create LLM], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit LLM]",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_default": schema.BoolAttribute{
                MarkdownDescription: "Is this the default LLM provider for the project? When set, the global LLM provider will not be used.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create LLM], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Edit LLM]",
                Computed: true,
            },
            "cost_per_million_tokens_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost per million tokens in USD cents. Used for billing when using global LLM providers.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Public], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *LlmProviderDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LlmProviderDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LlmProviderDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "llm-provider" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read llm_provider_data, got error: %s", err))
        return
    }

    var llmProviderDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &llmProviderDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse llm_provider_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := llmProviderDataResponse["data"].(map[string]interface{}); ok {
        llmProviderDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := llmProviderDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := llmProviderDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["llm_type"].(string); ok {
        data.LlmType = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["api_key"].(string); ok {
        data.ApiKey = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["model_name"].(string); ok {
        data.ModelName = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["base_url"].(string); ok {
        data.BaseUrl = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := llmProviderDataResponse["is_default"].(bool); ok {
        data.IsDefault = types.BoolValue(val)
    }
    if val, ok := llmProviderDataResponse["cost_per_million_tokens_in_usd_cents"].(float64); ok {
        data.CostPerMillionTokensInUsdCents = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
