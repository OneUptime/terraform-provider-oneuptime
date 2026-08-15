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
var _ datasource.DataSource = &AiConversationMessageDataSource{}

func NewAiConversationMessageDataSource() datasource.DataSource {
    return &AiConversationMessageDataSource{}
}

// AiConversationMessageDataSource defines the data source implementation.
type AiConversationMessageDataSource struct {
    client *Client
}

// AiConversationMessageDataSourceModel describes the data source data model.
type AiConversationMessageDataSourceModel struct {
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
    UserFeedback types.String `tfsdk:"user_feedback"`
}

func (d *AiConversationMessageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_conversation_message"
}

func (d *AiConversationMessageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A message in an AI conversation. Assistant messages carry citations, tool events and cost. Look up an existing ai_conversation_message by `id` or by `name`.",

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
            "conversation_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "role": schema.StringAttribute{
                MarkdownDescription: "Who authored this message: User or Assistant..",
                Computed: true,
            },
            "content_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Message content in markdown..",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this message..",
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "citations": schema.StringAttribute{
                MarkdownDescription: "Server-minted citations for this assistant message. Each citation records the tool, the exact validated query arguments and the row count..",
                Computed: true,
            },
            "widgets": schema.StringAttribute{
                MarkdownDescription: "Inline widgets (charts, tables, trace waterfalls, resource cards) built from this assistant message's tool results and rendered inline in the chat..",
                Computed: true,
            },
            "tool_actions": schema.StringAttribute{
                MarkdownDescription: "Mutating actions the agent proposed or performed in this turn, with their approval status (pending, approved, denied, executed)..",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Error message if this message failed to generate..",
                Computed: true,
            },
            "user_feedback": schema.StringAttribute{
                MarkdownDescription: "Thumbs feedback the user left on this assistant message: Up or Down..",
                Computed: true,
            },
        },
    }
}

func (d *AiConversationMessageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiConversationMessageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiConversationMessageDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a ai_conversation_message.",
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
        "conversationId": true,
        "userId": true,
        "role": true,
        "contentInMarkdown": true,
        "status": true,
        "aiRunId": true,
        "citations": true,
        "widgets": true,
        "toolActions": true,
        "errorMessage": true,
        "userFeedback": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/ai-conversation-message/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_conversation_message, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_conversation_message found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_conversation_message: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/ai-conversation-message/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list ai_conversation_message, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list ai_conversation_message: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai_conversation_message found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai_conversation_message matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for ai_conversation_message.")
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
    if obj, ok := item["role"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Role = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Role = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Role = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Role = types.StringValue(string(jsonBytes))
        } else {
            data.Role = types.StringNull()
        }
    } else if val, ok := item["role"].(string); ok {
        data.Role = types.StringValue(val)
    } else {
        data.Role = types.StringNull()
    }
    if obj, ok := item["contentInMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContentInMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContentInMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContentInMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContentInMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ContentInMarkdown = types.StringNull()
        }
    } else if val, ok := item["contentInMarkdown"].(string); ok {
        data.ContentInMarkdown = types.StringValue(val)
    } else {
        data.ContentInMarkdown = types.StringNull()
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
    if obj, ok := item["citations"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Citations = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Citations = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Citations = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Citations = types.StringValue(string(jsonBytes))
        } else {
            data.Citations = types.StringNull()
        }
    } else if val, ok := item["citations"].(string); ok {
        data.Citations = types.StringValue(val)
    } else {
        data.Citations = types.StringNull()
    }
    if obj, ok := item["widgets"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Widgets = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Widgets = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Widgets = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Widgets = types.StringValue(string(jsonBytes))
        } else {
            data.Widgets = types.StringNull()
        }
    } else if val, ok := item["widgets"].(string); ok {
        data.Widgets = types.StringValue(val)
    } else {
        data.Widgets = types.StringNull()
    }
    if obj, ok := item["toolActions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToolActions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToolActions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToolActions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToolActions = types.StringValue(string(jsonBytes))
        } else {
            data.ToolActions = types.StringNull()
        }
    } else if val, ok := item["toolActions"].(string); ok {
        data.ToolActions = types.StringValue(val)
    } else {
        data.ToolActions = types.StringNull()
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
    if obj, ok := item["userFeedback"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserFeedback = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserFeedback = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserFeedback = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserFeedback = types.StringValue(string(jsonBytes))
        } else {
            data.UserFeedback = types.StringNull()
        }
    } else if val, ok := item["userFeedback"].(string); ok {
        data.UserFeedback = types.StringValue(val)
    } else {
        data.UserFeedback = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
