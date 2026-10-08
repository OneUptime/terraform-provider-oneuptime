package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
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
        MarkdownDescription: "A message in an AI conversation. Assistant messages carry citations, tool events and cost. Look up an existing ai conversation message by `id`, or by any of its other arguments (`ai_run_id`, `content_in_markdown`, `conversation_id`, ...): each one set must match, and exactly one ai conversation message may match them all.",

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
                MarkdownDescription: "ID of the project this message belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "conversation_id": schema.StringAttribute{
                MarkdownDescription: "ID of the conversation this message belongs to. The ID of a `oneuptime_ai_conversation`.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who owns the conversation. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "role": schema.StringAttribute{
                MarkdownDescription: "Who authored this message: User or Assistant.",
                Optional: true,
                Computed: true,
            },
            "content_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Message content in markdown.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this message.",
                Optional: true,
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the AI run that produced this assistant message.",
                Optional: true,
                Computed: true,
            },
            "citations": schema.StringAttribute{
                MarkdownDescription: "Server-minted citations for this assistant message. Each citation records the tool, the exact validated query arguments and the row count. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "widgets": schema.StringAttribute{
                MarkdownDescription: "Inline widgets (charts, tables, trace waterfalls, resource cards) built from this assistant message's tool results and rendered inline in the chat. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "tool_actions": schema.StringAttribute{
                MarkdownDescription: "Mutating actions the agent proposed or performed in this turn, with their approval status (pending, approved, denied, executed). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Error message if this message failed to generate.",
                Optional: true,
                Computed: true,
            },
            "user_feedback": schema.StringAttribute{
                MarkdownDescription: "Thumbs feedback the user left on this assistant message: Up or Down.",
                Optional: true,
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

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ConversationId.IsNull() && !data.ConversationId.IsUnknown() {
        filters["conversationId"] = data.ConversationId.ValueString()
        filterNames = append(filterNames, "conversation_id = "+fmt.Sprintf("%q", data.ConversationId.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.Role.IsNull() && !data.Role.IsUnknown() {
        filters["role"] = data.Role.ValueString()
        filterNames = append(filterNames, "role = "+fmt.Sprintf("%q", data.Role.ValueString()))
    }
    if !data.ContentInMarkdown.IsNull() && !data.ContentInMarkdown.IsUnknown() {
        filters["contentInMarkdown"] = data.ContentInMarkdown.ValueString()
        filterNames = append(filterNames, "content_in_markdown = "+fmt.Sprintf("%q", data.ContentInMarkdown.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.AiRunId.IsNull() && !data.AiRunId.IsUnknown() {
        filters["aiRunId"] = data.AiRunId.ValueString()
        filterNames = append(filterNames, "ai_run_id = "+fmt.Sprintf("%q", data.AiRunId.ValueString()))
    }
    if !data.ErrorMessage.IsNull() && !data.ErrorMessage.IsUnknown() {
        filters["errorMessage"] = data.ErrorMessage.ValueString()
        filterNames = append(filterNames, "error_message = "+fmt.Sprintf("%q", data.ErrorMessage.ValueString()))
    }
    if !data.UserFeedback.IsNull() && !data.UserFeedback.IsUnknown() {
        filters["userFeedback"] = data.UserFeedback.ValueString()
        filterNames = append(filterNames, "user_feedback = "+fmt.Sprintf("%q", data.UserFeedback.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the ai conversation message up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the ai conversation message up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai conversation message found with id %q.", data.Id.ValueString()))
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
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/ai-conversation-message/get-list", listBody)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai conversation message matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai conversation message matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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
