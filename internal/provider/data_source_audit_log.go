package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuditLogDataSource{}

func NewAuditLogDataSource() datasource.DataSource {
    return &AuditLogDataSource{}
}

// AuditLogDataSource defines the data source implementation.
type AuditLogDataSource struct {
    client *Client
}

// AuditLogDataSourceModel describes the data source data model.
type AuditLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    ResourceName types.String `tfsdk:"resource_name"`
    RootResourceType types.String `tfsdk:"root_resource_type"`
    RootResourceId types.String `tfsdk:"root_resource_id"`
    Action types.String `tfsdk:"action"`
    UserId types.String `tfsdk:"user_id"`
    UserName types.String `tfsdk:"user_name"`
    UserEmail types.String `tfsdk:"user_email"`
    UserType types.String `tfsdk:"user_type"`
    ApiKeyId types.String `tfsdk:"api_key_id"`
    ApiKeyName types.String `tfsdk:"api_key_name"`
    McpOAuthGrantId types.String `tfsdk:"mcp_o_auth_grant_id"`
    McpClientName types.String `tfsdk:"mcp_client_name"`
    WorkflowId types.String `tfsdk:"workflow_id"`
    WorkflowName types.String `tfsdk:"workflow_name"`
    Changes types.Set `tfsdk:"changes"`
}

func (d *AuditLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_audit_log"
}

func (d *AuditLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Audit Log Look up an existing audit log by `id`, or by any of its other arguments (`action`, `api_key_id`, `api_key_name`, ...): each one set must match, and exactly one audit log may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID.",
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "Resource Type.",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "Resource ID.",
                Optional: true,
                Computed: true,
            },
            "resource_name": schema.StringAttribute{
                MarkdownDescription: "Resource Name.",
                Optional: true,
                Computed: true,
            },
            "root_resource_type": schema.StringAttribute{
                MarkdownDescription: "Root Resource Type.",
                Optional: true,
                Computed: true,
            },
            "root_resource_id": schema.StringAttribute{
                MarkdownDescription: "Root Resource ID.",
                Optional: true,
                Computed: true,
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "Action.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "User ID.",
                Optional: true,
                Computed: true,
            },
            "user_name": schema.StringAttribute{
                MarkdownDescription: "User Name.",
                Optional: true,
                Computed: true,
            },
            "user_email": schema.StringAttribute{
                MarkdownDescription: "User Email.",
                Optional: true,
                Computed: true,
            },
            "user_type": schema.StringAttribute{
                MarkdownDescription: "User Type.",
                Optional: true,
                Computed: true,
            },
            "api_key_id": schema.StringAttribute{
                MarkdownDescription: "API Key ID.",
                Optional: true,
                Computed: true,
            },
            "api_key_name": schema.StringAttribute{
                MarkdownDescription: "API Key Name.",
                Optional: true,
                Computed: true,
            },
            "mcp_o_auth_grant_id": schema.StringAttribute{
                MarkdownDescription: "MCP Client Authorization ID.",
                Optional: true,
                Computed: true,
            },
            "mcp_client_name": schema.StringAttribute{
                MarkdownDescription: "MCP Client Name.",
                Optional: true,
                Computed: true,
            },
            "workflow_id": schema.StringAttribute{
                MarkdownDescription: "Workflow ID.",
                Optional: true,
                Computed: true,
            },
            "workflow_name": schema.StringAttribute{
                MarkdownDescription: "Workflow Name.",
                Optional: true,
                Computed: true,
            },
            "changes": schema.SetAttribute{
                MarkdownDescription: "Changes.",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *AuditLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AuditLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AuditLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
        filters["resourceType"] = data.ResourceType.ValueString()
        filterNames = append(filterNames, "resource_type = "+fmt.Sprintf("%q", data.ResourceType.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.ResourceName.IsNull() && !data.ResourceName.IsUnknown() {
        filters["resourceName"] = data.ResourceName.ValueString()
        filterNames = append(filterNames, "resource_name = "+fmt.Sprintf("%q", data.ResourceName.ValueString()))
    }
    if !data.RootResourceType.IsNull() && !data.RootResourceType.IsUnknown() {
        filters["rootResourceType"] = data.RootResourceType.ValueString()
        filterNames = append(filterNames, "root_resource_type = "+fmt.Sprintf("%q", data.RootResourceType.ValueString()))
    }
    if !data.RootResourceId.IsNull() && !data.RootResourceId.IsUnknown() {
        filters["rootResourceId"] = data.RootResourceId.ValueString()
        filterNames = append(filterNames, "root_resource_id = "+fmt.Sprintf("%q", data.RootResourceId.ValueString()))
    }
    if !data.Action.IsNull() && !data.Action.IsUnknown() {
        filters["action"] = data.Action.ValueString()
        filterNames = append(filterNames, "action = "+fmt.Sprintf("%q", data.Action.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.UserName.IsNull() && !data.UserName.IsUnknown() {
        filters["userName"] = data.UserName.ValueString()
        filterNames = append(filterNames, "user_name = "+fmt.Sprintf("%q", data.UserName.ValueString()))
    }
    if !data.UserEmail.IsNull() && !data.UserEmail.IsUnknown() {
        filters["userEmail"] = data.UserEmail.ValueString()
        filterNames = append(filterNames, "user_email = "+fmt.Sprintf("%q", data.UserEmail.ValueString()))
    }
    if !data.UserType.IsNull() && !data.UserType.IsUnknown() {
        filters["userType"] = data.UserType.ValueString()
        filterNames = append(filterNames, "user_type = "+fmt.Sprintf("%q", data.UserType.ValueString()))
    }
    if !data.ApiKeyId.IsNull() && !data.ApiKeyId.IsUnknown() {
        filters["apiKeyId"] = data.ApiKeyId.ValueString()
        filterNames = append(filterNames, "api_key_id = "+fmt.Sprintf("%q", data.ApiKeyId.ValueString()))
    }
    if !data.ApiKeyName.IsNull() && !data.ApiKeyName.IsUnknown() {
        filters["apiKeyName"] = data.ApiKeyName.ValueString()
        filterNames = append(filterNames, "api_key_name = "+fmt.Sprintf("%q", data.ApiKeyName.ValueString()))
    }
    if !data.McpOAuthGrantId.IsNull() && !data.McpOAuthGrantId.IsUnknown() {
        filters["mcpOAuthGrantId"] = data.McpOAuthGrantId.ValueString()
        filterNames = append(filterNames, "mcp_o_auth_grant_id = "+fmt.Sprintf("%q", data.McpOAuthGrantId.ValueString()))
    }
    if !data.McpClientName.IsNull() && !data.McpClientName.IsUnknown() {
        filters["mcpClientName"] = data.McpClientName.ValueString()
        filterNames = append(filterNames, "mcp_client_name = "+fmt.Sprintf("%q", data.McpClientName.ValueString()))
    }
    if !data.WorkflowId.IsNull() && !data.WorkflowId.IsUnknown() {
        filters["workflowId"] = data.WorkflowId.ValueString()
        filterNames = append(filterNames, "workflow_id = "+fmt.Sprintf("%q", data.WorkflowId.ValueString()))
    }
    if !data.WorkflowName.IsNull() && !data.WorkflowName.IsUnknown() {
        filters["workflowName"] = data.WorkflowName.ValueString()
        filterNames = append(filterNames, "workflow_name = "+fmt.Sprintf("%q", data.WorkflowName.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the audit log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the audit log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "resourceType": true,
        "resourceId": true,
        "resourceName": true,
        "rootResourceType": true,
        "rootResourceId": true,
        "action": true,
        "userId": true,
        "userName": true,
        "userEmail": true,
        "userType": true,
        "apiKeyId": true,
        "apiKeyName": true,
        "mcpOAuthGrantId": true,
        "mcpClientName": true,
        "workflowId": true,
        "workflowName": true,
        "changes": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/audit-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read audit_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No audit log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read audit_log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/audit-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list audit_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list audit_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No audit log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one audit log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for audit_log.")
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
    if obj, ok := item["resourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceType = types.StringNull()
        }
    } else if val, ok := item["resourceType"].(string); ok {
        data.ResourceType = types.StringValue(val)
    } else {
        data.ResourceType = types.StringNull()
    }
    if obj, ok := item["resourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if val, ok := item["resourceId"].(string); ok {
        data.ResourceId = types.StringValue(val)
    } else {
        data.ResourceId = types.StringNull()
    }
    if obj, ok := item["resourceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceName = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceName = types.StringNull()
        }
    } else if val, ok := item["resourceName"].(string); ok {
        data.ResourceName = types.StringValue(val)
    } else {
        data.ResourceName = types.StringNull()
    }
    if obj, ok := item["rootResourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.RootResourceType = types.StringNull()
        }
    } else if val, ok := item["rootResourceType"].(string); ok {
        data.RootResourceType = types.StringValue(val)
    } else {
        data.RootResourceType = types.StringNull()
    }
    if obj, ok := item["rootResourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.RootResourceId = types.StringNull()
        }
    } else if val, ok := item["rootResourceId"].(string); ok {
        data.RootResourceId = types.StringValue(val)
    } else {
        data.RootResourceId = types.StringNull()
    }
    if obj, ok := item["action"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := item["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
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
    if obj, ok := item["userName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserName = types.StringValue(string(jsonBytes))
        } else {
            data.UserName = types.StringNull()
        }
    } else if val, ok := item["userName"].(string); ok {
        data.UserName = types.StringValue(val)
    } else {
        data.UserName = types.StringNull()
    }
    if obj, ok := item["userEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserEmail = types.StringValue(string(jsonBytes))
        } else {
            data.UserEmail = types.StringNull()
        }
    } else if val, ok := item["userEmail"].(string); ok {
        data.UserEmail = types.StringValue(val)
    } else {
        data.UserEmail = types.StringNull()
    }
    if obj, ok := item["userType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserType = types.StringValue(string(jsonBytes))
        } else {
            data.UserType = types.StringNull()
        }
    } else if val, ok := item["userType"].(string); ok {
        data.UserType = types.StringValue(val)
    } else {
        data.UserType = types.StringNull()
    }
    if obj, ok := item["apiKeyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApiKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApiKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApiKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApiKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.ApiKeyId = types.StringNull()
        }
    } else if val, ok := item["apiKeyId"].(string); ok {
        data.ApiKeyId = types.StringValue(val)
    } else {
        data.ApiKeyId = types.StringNull()
    }
    if obj, ok := item["apiKeyName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApiKeyName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApiKeyName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApiKeyName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApiKeyName = types.StringValue(string(jsonBytes))
        } else {
            data.ApiKeyName = types.StringNull()
        }
    } else if val, ok := item["apiKeyName"].(string); ok {
        data.ApiKeyName = types.StringValue(val)
    } else {
        data.ApiKeyName = types.StringNull()
    }
    if obj, ok := item["mcpOAuthGrantId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.McpOAuthGrantId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.McpOAuthGrantId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.McpOAuthGrantId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.McpOAuthGrantId = types.StringValue(string(jsonBytes))
        } else {
            data.McpOAuthGrantId = types.StringNull()
        }
    } else if val, ok := item["mcpOAuthGrantId"].(string); ok {
        data.McpOAuthGrantId = types.StringValue(val)
    } else {
        data.McpOAuthGrantId = types.StringNull()
    }
    if obj, ok := item["mcpClientName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.McpClientName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.McpClientName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.McpClientName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.McpClientName = types.StringValue(string(jsonBytes))
        } else {
            data.McpClientName = types.StringNull()
        }
    } else if val, ok := item["mcpClientName"].(string); ok {
        data.McpClientName = types.StringValue(val)
    } else {
        data.McpClientName = types.StringNull()
    }
    if obj, ok := item["workflowId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkflowId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkflowId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkflowId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkflowId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowId = types.StringNull()
        }
    } else if val, ok := item["workflowId"].(string); ok {
        data.WorkflowId = types.StringValue(val)
    } else {
        data.WorkflowId = types.StringNull()
    }
    if obj, ok := item["workflowName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkflowName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkflowName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkflowName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkflowName = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowName = types.StringNull()
        }
    } else if val, ok := item["workflowName"].(string); ok {
        data.WorkflowName = types.StringValue(val)
    } else {
        data.WorkflowName = types.StringNull()
    }
    if val, ok := item["changes"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Changes = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Changes = types.SetNull(types.StringType)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
