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
var _ datasource.DataSource = &AiConversationMessageDataDataSource{}

func NewAiConversationMessageDataDataSource() datasource.DataSource {
    return &AiConversationMessageDataDataSource{}
}

// AiConversationMessageDataDataSource defines the data source implementation.
type AiConversationMessageDataDataSource struct {
    client *Client
}

// AiConversationMessageDataDataSourceModel describes the data source data model.
type AiConversationMessageDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ConversationId types.String `tfsdk:"conversation_id"`
    UserId types.String `tfsdk:"user_id"`
    Role types.String `tfsdk:"role"`
    ContentInMarkdown types.String `tfsdk:"content_in_markdown"`
    Status types.String `tfsdk:"status"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    Citations types.String `tfsdk:"citations"`
    Widgets types.String `tfsdk:"widgets"`
    ToolActions types.String `tfsdk:"tool_actions"`
    ErrorMessage types.String `tfsdk:"error_message"`
}

func (d *AiConversationMessageDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_conversation_message_data"
}

func (d *AiConversationMessageDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "ai_conversation_message_data data source",

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
            "conversation_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "role": schema.StringAttribute{
                MarkdownDescription: "Who authored this message: User or Assistant.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "content_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Message content in markdown.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this message.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "citations": schema.StringAttribute{
                MarkdownDescription: "Server-minted citations for this assistant message. Each citation records the tool, the exact validated query arguments and the row count.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "widgets": schema.StringAttribute{
                MarkdownDescription: "Inline widgets (charts, tables, trace waterfalls, resource cards) built from this assistant message's tool results and rendered inline in the chat.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "tool_actions": schema.StringAttribute{
                MarkdownDescription: "Mutating actions the agent proposed or performed in this turn, with their approval status (pending, approved, denied, executed).. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Error message if this message failed to generate.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *AiConversationMessageDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiConversationMessageDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiConversationMessageDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "ai-conversation-message" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_conversation_message_data, got error: %s", err))
        return
    }

    var aiConversationMessageDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &aiConversationMessageDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_conversation_message_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := aiConversationMessageDataResponse["data"].(map[string]interface{}); ok {
        aiConversationMessageDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := aiConversationMessageDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := aiConversationMessageDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["conversation_id"].(string); ok {
        data.ConversationId = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["role"].(string); ok {
        data.Role = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["content_in_markdown"].(string); ok {
        data.ContentInMarkdown = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["ai_run_id"].(string); ok {
        data.AiRunId = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["citations"].(string); ok {
        data.Citations = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["widgets"].(string); ok {
        data.Widgets = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["tool_actions"].(string); ok {
        data.ToolActions = types.StringValue(val)
    }
    if val, ok := aiConversationMessageDataResponse["error_message"].(string); ok {
        data.ErrorMessage = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
