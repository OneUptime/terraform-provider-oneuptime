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
var _ datasource.DataSource = &ExternalMcpServerDataSource{}

func NewExternalMcpServerDataSource() datasource.DataSource {
    return &ExternalMcpServerDataSource{}
}

// ExternalMcpServerDataSource defines the data source implementation.
type ExternalMcpServerDataSource struct {
    client *Client
}

// ExternalMcpServerDataSourceModel describes the data source data model.
type ExternalMcpServerDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Url types.String `tfsdk:"url"`
    AuthType types.String `tfsdk:"auth_type"`
    HeaderNames types.String `tfsdk:"header_names"`
    OauthTokenUrl types.String `tfsdk:"oauth_token_url"`
    OauthClientId types.String `tfsdk:"oauth_client_id"`
    OauthScope types.String `tfsdk:"oauth_scope"`
    OauthAudience types.String `tfsdk:"oauth_audience"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Transport types.String `tfsdk:"transport"`
    ServerName types.String `tfsdk:"server_name"`
    ServerVersion types.String `tfsdk:"server_version"`
    ProtocolVersion types.String `tfsdk:"protocol_version"`
    Tools types.String `tfsdk:"tools"`
    LastCheckedAt types.String `tfsdk:"last_checked_at"`
    LastConnectedAt types.String `tfsdk:"last_connected_at"`
    LastError types.String `tfsdk:"last_error"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ExternalMcpServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_external_mcp_server"
}

func (d *ExternalMcpServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "MCP servers outside OneUptime - such as AWS, Azure or GitHub - whose tools OneUptime AI may use while it investigates and fixes incidents and alerts. Look up an existing external mcp server by `id`, or by any of its other arguments (`name`, `auth_type`, `created_by_user_id`, ...): each one set must match, and exactly one external mcp server may match them all.",

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
                MarkdownDescription: "ID of the project this MCP server belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "A name for the server, such as AWS or Azure. OneUptime AI and the people reading what it did see the server by this name.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the server is for, such as 'Our production AWS account, read-only'. OneUptime AI reads it to choose the right server.",
                Optional: true,
                Computed: true,
            },
            "url": schema.StringAttribute{
                MarkdownDescription: "The server's MCP address, such as https://mcp.example.com/mcp. Streamable HTTP and SSE servers both work: OneUptime finds out which when it connects. Never put a token in it: use Authentication.",
                Optional: true,
                Computed: true,
            },
            "auth_type": schema.StringAttribute{
                MarkdownDescription: "How OneUptime signs in to the server: None, BearerToken (an API token), Headers (headers of the server's own, such as X-API-Key) or OAuthClientCredentials (a token from an OAuth token URL, with a client ID and secret).",
                Optional: true,
                Computed: true,
            },
            "header_names": schema.StringAttribute{
                MarkdownDescription: "The names of the headers the server signs in with, without their values. Worked out from the headers on every read.",
                Optional: true,
                Computed: true,
            },
            "oauth_token_url": schema.StringAttribute{
                MarkdownDescription: "For OAuthClientCredentials: the OAuth token endpoint, such as https://login.microsoftonline.com/<tenant>/oauth2/v2.0/token.",
                Optional: true,
                Computed: true,
            },
            "oauth_client_id": schema.StringAttribute{
                MarkdownDescription: "For OAuthClientCredentials: the OAuth client ID.",
                Optional: true,
                Computed: true,
            },
            "oauth_scope": schema.StringAttribute{
                MarkdownDescription: "For OAuthClientCredentials: the scope to ask for, such as api://<application-id>/.default. Leave empty when the token URL needs none.",
                Optional: true,
                Computed: true,
            },
            "oauth_audience": schema.StringAttribute{
                MarkdownDescription: "For OAuthClientCredentials: the audience to ask for, when the token URL needs one (Auth0 does). Leave empty otherwise.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI may use this server. A server that is off is never used, whatever a rule names.",
                Optional: true,
                Computed: true,
            },
            "transport": schema.StringAttribute{
                MarkdownDescription: "How the server talks, as OneUptime found out when it last connected: StreamableHttp or Sse.",
                Optional: true,
                Computed: true,
            },
            "server_name": schema.StringAttribute{
                MarkdownDescription: "What the server called itself when OneUptime last connected.",
                Optional: true,
                Computed: true,
            },
            "server_version": schema.StringAttribute{
                MarkdownDescription: "The version the server gave when OneUptime last connected, if any.",
                Optional: true,
                Computed: true,
            },
            "protocol_version": schema.StringAttribute{
                MarkdownDescription: "The MCP protocol version the server agreed on when OneUptime last connected.",
                Optional: true,
                Computed: true,
            },
            "tools": schema.StringAttribute{
                MarkdownDescription: "The tools the server offered when OneUptime last connected: each one's name, description and whether the server marks it read-only. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "last_checked_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last tried to connect to the server. Empty means never.",
                Computed: true,
            },
            "last_connected_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last connected to the server and read its tools.",
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "Why the last connection failed, with secrets removed. Cleared by the next connection that works.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who added this MCP server. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *ExternalMcpServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExternalMcpServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ExternalMcpServerDataSourceModel

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
    if !data.Url.IsNull() && !data.Url.IsUnknown() {
        filters["url"] = data.Url.ValueString()
        filterNames = append(filterNames, "url = "+fmt.Sprintf("%q", data.Url.ValueString()))
    }
    if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
        filters["authType"] = data.AuthType.ValueString()
        filterNames = append(filterNames, "auth_type = "+fmt.Sprintf("%q", data.AuthType.ValueString()))
    }
    if !data.HeaderNames.IsNull() && !data.HeaderNames.IsUnknown() {
        filters["headerNames"] = data.HeaderNames.ValueString()
        filterNames = append(filterNames, "header_names = "+fmt.Sprintf("%q", data.HeaderNames.ValueString()))
    }
    if !data.OauthTokenUrl.IsNull() && !data.OauthTokenUrl.IsUnknown() {
        filters["oauthTokenUrl"] = data.OauthTokenUrl.ValueString()
        filterNames = append(filterNames, "oauth_token_url = "+fmt.Sprintf("%q", data.OauthTokenUrl.ValueString()))
    }
    if !data.OauthClientId.IsNull() && !data.OauthClientId.IsUnknown() {
        filters["oauthClientId"] = data.OauthClientId.ValueString()
        filterNames = append(filterNames, "oauth_client_id = "+fmt.Sprintf("%q", data.OauthClientId.ValueString()))
    }
    if !data.OauthScope.IsNull() && !data.OauthScope.IsUnknown() {
        filters["oauthScope"] = data.OauthScope.ValueString()
        filterNames = append(filterNames, "oauth_scope = "+fmt.Sprintf("%q", data.OauthScope.ValueString()))
    }
    if !data.OauthAudience.IsNull() && !data.OauthAudience.IsUnknown() {
        filters["oauthAudience"] = data.OauthAudience.ValueString()
        filterNames = append(filterNames, "oauth_audience = "+fmt.Sprintf("%q", data.OauthAudience.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.Transport.IsNull() && !data.Transport.IsUnknown() {
        filters["transport"] = data.Transport.ValueString()
        filterNames = append(filterNames, "transport = "+fmt.Sprintf("%q", data.Transport.ValueString()))
    }
    if !data.ServerName.IsNull() && !data.ServerName.IsUnknown() {
        filters["serverName"] = data.ServerName.ValueString()
        filterNames = append(filterNames, "server_name = "+fmt.Sprintf("%q", data.ServerName.ValueString()))
    }
    if !data.ServerVersion.IsNull() && !data.ServerVersion.IsUnknown() {
        filters["serverVersion"] = data.ServerVersion.ValueString()
        filterNames = append(filterNames, "server_version = "+fmt.Sprintf("%q", data.ServerVersion.ValueString()))
    }
    if !data.ProtocolVersion.IsNull() && !data.ProtocolVersion.IsUnknown() {
        filters["protocolVersion"] = data.ProtocolVersion.ValueString()
        filterNames = append(filterNames, "protocol_version = "+fmt.Sprintf("%q", data.ProtocolVersion.ValueString()))
    }
    if !data.LastError.IsNull() && !data.LastError.IsUnknown() {
        filters["lastError"] = data.LastError.ValueString()
        filterNames = append(filterNames, "last_error = "+fmt.Sprintf("%q", data.LastError.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the external mcp server up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the external mcp server up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "url": true,
        "authType": true,
        "headerNames": true,
        "oauthTokenUrl": true,
        "oauthClientId": true,
        "oauthScope": true,
        "oauthAudience": true,
        "isEnabled": true,
        "transport": true,
        "serverName": true,
        "serverVersion": true,
        "protocolVersion": true,
        "tools": true,
        "lastCheckedAt": true,
        "lastConnectedAt": true,
        "lastError": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/external-mcp-server/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read external_mcp_server, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No external mcp server found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read external_mcp_server: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/external-mcp-server/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list external_mcp_server, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list external_mcp_server: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No external mcp server matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one external mcp server matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for external_mcp_server.")
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
    if obj, ok := item["url"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Url = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Url = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Url = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Url = types.StringValue(string(jsonBytes))
        } else {
            data.Url = types.StringNull()
        }
    } else if val, ok := item["url"].(string); ok {
        data.Url = types.StringValue(val)
    } else {
        data.Url = types.StringNull()
    }
    if obj, ok := item["authType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AuthType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if val, ok := item["authType"].(string); ok {
        data.AuthType = types.StringValue(val)
    } else {
        data.AuthType = types.StringNull()
    }
    if obj, ok := item["headerNames"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeaderNames = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HeaderNames = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HeaderNames = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HeaderNames = types.StringValue(string(jsonBytes))
        } else {
            data.HeaderNames = types.StringNull()
        }
    } else if val, ok := item["headerNames"].(string); ok {
        data.HeaderNames = types.StringValue(val)
    } else {
        data.HeaderNames = types.StringNull()
    }
    if obj, ok := item["oauthTokenUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthTokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthTokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthTokenUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthTokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.OauthTokenUrl = types.StringNull()
        }
    } else if val, ok := item["oauthTokenUrl"].(string); ok {
        data.OauthTokenUrl = types.StringValue(val)
    } else {
        data.OauthTokenUrl = types.StringNull()
    }
    if obj, ok := item["oauthClientId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthClientId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthClientId = types.StringValue(string(jsonBytes))
        } else {
            data.OauthClientId = types.StringNull()
        }
    } else if val, ok := item["oauthClientId"].(string); ok {
        data.OauthClientId = types.StringValue(val)
    } else {
        data.OauthClientId = types.StringNull()
    }
    if obj, ok := item["oauthScope"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthScope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthScope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthScope = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthScope = types.StringValue(string(jsonBytes))
        } else {
            data.OauthScope = types.StringNull()
        }
    } else if val, ok := item["oauthScope"].(string); ok {
        data.OauthScope = types.StringValue(val)
    } else {
        data.OauthScope = types.StringNull()
    }
    if obj, ok := item["oauthAudience"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthAudience = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthAudience = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthAudience = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthAudience = types.StringValue(string(jsonBytes))
        } else {
            data.OauthAudience = types.StringNull()
        }
    } else if val, ok := item["oauthAudience"].(string); ok {
        data.OauthAudience = types.StringValue(val)
    } else {
        data.OauthAudience = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["transport"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Transport = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Transport = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Transport = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Transport = types.StringValue(string(jsonBytes))
        } else {
            data.Transport = types.StringNull()
        }
    } else if val, ok := item["transport"].(string); ok {
        data.Transport = types.StringValue(val)
    } else {
        data.Transport = types.StringNull()
    }
    if obj, ok := item["serverName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServerName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServerName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServerName = types.StringValue(string(jsonBytes))
        } else {
            data.ServerName = types.StringNull()
        }
    } else if val, ok := item["serverName"].(string); ok {
        data.ServerName = types.StringValue(val)
    } else {
        data.ServerName = types.StringNull()
    }
    if obj, ok := item["serverVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServerVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServerVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServerVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ServerVersion = types.StringNull()
        }
    } else if val, ok := item["serverVersion"].(string); ok {
        data.ServerVersion = types.StringValue(val)
    } else {
        data.ServerVersion = types.StringNull()
    }
    if obj, ok := item["protocolVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProtocolVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProtocolVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProtocolVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProtocolVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ProtocolVersion = types.StringNull()
        }
    } else if val, ok := item["protocolVersion"].(string); ok {
        data.ProtocolVersion = types.StringValue(val)
    } else {
        data.ProtocolVersion = types.StringNull()
    }
    if obj, ok := item["tools"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Tools = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Tools = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Tools = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Tools = types.StringValue(string(jsonBytes))
        } else {
            data.Tools = types.StringNull()
        }
    } else if val, ok := item["tools"].(string); ok {
        data.Tools = types.StringValue(val)
    } else {
        data.Tools = types.StringNull()
    }
    if obj, ok := item["lastCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastCheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastCheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastCheckedAt = types.StringNull()
        }
    } else if val, ok := item["lastCheckedAt"].(string); ok {
        data.LastCheckedAt = types.StringValue(val)
    } else {
        data.LastCheckedAt = types.StringNull()
    }
    if obj, ok := item["lastConnectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastConnectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastConnectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastConnectedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastConnectedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastConnectedAt = types.StringNull()
        }
    } else if val, ok := item["lastConnectedAt"].(string); ok {
        data.LastConnectedAt = types.StringValue(val)
    } else {
        data.LastConnectedAt = types.StringNull()
    }
    if obj, ok := item["lastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := item["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
