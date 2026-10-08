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
var _ datasource.DataSource = &WorkflowVariableDataSource{}

func NewWorkflowVariableDataSource() datasource.DataSource {
    return &WorkflowVariableDataSource{}
}

// WorkflowVariableDataSource defines the data source implementation.
type WorkflowVariableDataSource struct {
    client *Client
}

// WorkflowVariableDataSourceModel describes the data source data model.
type WorkflowVariableDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    WorkflowId types.String `tfsdk:"workflow_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsSecret types.Bool `tfsdk:"is_secret"`
    VariableType types.String `tfsdk:"variable_type"`
    OauthGrantType types.String `tfsdk:"oauth_grant_type"`
    OauthTokenUrl types.String `tfsdk:"oauth_token_url"`
    OauthClientId types.String `tfsdk:"oauth_client_id"`
    OauthScope types.String `tfsdk:"oauth_scope"`
    OauthAdditionalParameters types.String `tfsdk:"oauth_additional_parameters"`
    OauthClientAuthenticationMethod types.String `tfsdk:"oauth_client_authentication_method"`
    OauthAccessTokenExpiresAt types.String `tfsdk:"oauth_access_token_expires_at"`
    OauthLastRefreshedAt types.String `tfsdk:"oauth_last_refreshed_at"`
    OauthLastRefreshError types.String `tfsdk:"oauth_last_refresh_error"`
    OauthLastRefreshErrorAt types.String `tfsdk:"oauth_last_refresh_error_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *WorkflowVariableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflow_variable"
}

func (d *WorkflowVariableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Store environment variables or secrets for your workflows. Look up an existing workflow variable by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one workflow variable may match them all.",

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
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "workflow_id": schema.StringAttribute{
                MarkdownDescription: "ID of Workflow this variable belong to. If this is null then this variable will be a global variable. The ID of a `oneuptime_workflow`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Variable Name.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "is_secret": schema.BoolAttribute{
                MarkdownDescription: "Is this variable a secret. If true, then it'll not be in the logs.",
                Optional: true,
                Computed: true,
            },
            "variable_type": schema.StringAttribute{
                MarkdownDescription: "Static: the content you save is used as is. OAuth 2.0: OneUptime fetches an access token from your identity provider and refreshes it automatically when a workflow uses it after it has expired.",
                Optional: true,
                Computed: true,
            },
            "oauth_grant_type": schema.StringAttribute{
                MarkdownDescription: "OAuth 2.0 variables only. Client Credentials for machine-to-machine access, or Refresh Token to keep delegated access alive with a refresh token you obtained once.",
                Optional: true,
                Computed: true,
            },
            "oauth_token_url": schema.StringAttribute{
                MarkdownDescription: "OAuth 2.0 variables only. The token endpoint of your identity provider.",
                Optional: true,
                Computed: true,
            },
            "oauth_client_id": schema.StringAttribute{
                MarkdownDescription: "OAuth 2.0 variables only. The client ID of the application registered with your identity provider.",
                Optional: true,
                Computed: true,
            },
            "oauth_scope": schema.StringAttribute{
                MarkdownDescription: "OAuth 2.0 variables only. Space-separated scopes to request. Leave empty to use the scopes your identity provider grants by default.",
                Optional: true,
                Computed: true,
            },
            "oauth_additional_parameters": schema.StringAttribute{
                MarkdownDescription: "OAuth 2.0 variables only. Extra form parameters sent with every token request, such as audience for Auth0 or resource for Azure AD v1. Readable by anyone who can read the variable, so do not put secrets here. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "oauth_client_authentication_method": schema.StringAttribute{
                MarkdownDescription: "OAuth 2.0 variables only. How the client ID and secret are sent: in an HTTP Basic header (client_secret_basic, the default) or in the request body (client_secret_post).",
                Optional: true,
                Computed: true,
            },
            "oauth_access_token_expires_at": schema.StringAttribute{
                MarkdownDescription: "When the cached access token expires, as reported by the identity provider (expires_in) or by the token itself (the JWT exp claim). Empty when neither says.",
                Computed: true,
            },
            "oauth_last_refreshed_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last fetched an access token for this variable. Cleared when the OAuth settings change.",
                Computed: true,
            },
            "oauth_last_refresh_error": schema.StringAttribute{
                MarkdownDescription: "Why the last attempt to fetch an access token failed. Cleared by the next successful refresh.",
                Optional: true,
                Computed: true,
            },
            "oauth_last_refresh_error_at": schema.StringAttribute{
                MarkdownDescription: "When the last failed attempt to fetch an access token happened.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *WorkflowVariableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkflowVariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkflowVariableDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.WorkflowId.IsNull() && !data.WorkflowId.IsUnknown() {
        filters["workflowId"] = data.WorkflowId.ValueString()
        filterNames = append(filterNames, "workflow_id = "+fmt.Sprintf("%q", data.WorkflowId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.IsSecret.IsNull() && !data.IsSecret.IsUnknown() {
        filters["isSecret"] = data.IsSecret.ValueBool()
        filterNames = append(filterNames, "is_secret = "+fmt.Sprintf("%t", data.IsSecret.ValueBool()))
    }
    if !data.VariableType.IsNull() && !data.VariableType.IsUnknown() {
        filters["variableType"] = data.VariableType.ValueString()
        filterNames = append(filterNames, "variable_type = "+fmt.Sprintf("%q", data.VariableType.ValueString()))
    }
    if !data.OauthGrantType.IsNull() && !data.OauthGrantType.IsUnknown() {
        filters["oauthGrantType"] = data.OauthGrantType.ValueString()
        filterNames = append(filterNames, "oauth_grant_type = "+fmt.Sprintf("%q", data.OauthGrantType.ValueString()))
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
    if !data.OauthClientAuthenticationMethod.IsNull() && !data.OauthClientAuthenticationMethod.IsUnknown() {
        filters["oauthClientAuthenticationMethod"] = data.OauthClientAuthenticationMethod.ValueString()
        filterNames = append(filterNames, "oauth_client_authentication_method = "+fmt.Sprintf("%q", data.OauthClientAuthenticationMethod.ValueString()))
    }
    if !data.OauthLastRefreshError.IsNull() && !data.OauthLastRefreshError.IsUnknown() {
        filters["oauthLastRefreshError"] = data.OauthLastRefreshError.ValueString()
        filterNames = append(filterNames, "oauth_last_refresh_error = "+fmt.Sprintf("%q", data.OauthLastRefreshError.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the workflow variable up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the workflow variable up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "workflowId": true,
        "name": true,
        "description": true,
        "isSecret": true,
        "variableType": true,
        "oauthGrantType": true,
        "oauthTokenUrl": true,
        "oauthClientId": true,
        "oauthScope": true,
        "oauthAdditionalParameters": true,
        "oauthClientAuthenticationMethod": true,
        "oauthAccessTokenExpiresAt": true,
        "oauthLastRefreshedAt": true,
        "oauthLastRefreshError": true,
        "oauthLastRefreshErrorAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workflow-variable/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workflow_variable, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow variable found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workflow_variable: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workflow-variable/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workflow_variable, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workflow_variable: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow variable matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workflow variable matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workflow_variable.")
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
    if val, ok := item["isSecret"].(bool); ok {
        data.IsSecret = types.BoolValue(val)
    } else {
        data.IsSecret = types.BoolNull()
    }
    if obj, ok := item["variableType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VariableType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VariableType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VariableType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VariableType = types.StringValue(string(jsonBytes))
        } else {
            data.VariableType = types.StringNull()
        }
    } else if val, ok := item["variableType"].(string); ok {
        data.VariableType = types.StringValue(val)
    } else {
        data.VariableType = types.StringNull()
    }
    if obj, ok := item["oauthGrantType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthGrantType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthGrantType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthGrantType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthGrantType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthGrantType = types.StringNull()
        }
    } else if val, ok := item["oauthGrantType"].(string); ok {
        data.OauthGrantType = types.StringValue(val)
    } else {
        data.OauthGrantType = types.StringNull()
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
    if obj, ok := item["oauthAdditionalParameters"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthAdditionalParameters = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthAdditionalParameters = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthAdditionalParameters = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthAdditionalParameters = types.StringValue(string(jsonBytes))
        } else {
            data.OauthAdditionalParameters = types.StringNull()
        }
    } else if val, ok := item["oauthAdditionalParameters"].(string); ok {
        data.OauthAdditionalParameters = types.StringValue(val)
    } else {
        data.OauthAdditionalParameters = types.StringNull()
    }
    if obj, ok := item["oauthClientAuthenticationMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthClientAuthenticationMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthClientAuthenticationMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthClientAuthenticationMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthClientAuthenticationMethod = types.StringValue(string(jsonBytes))
        } else {
            data.OauthClientAuthenticationMethod = types.StringNull()
        }
    } else if val, ok := item["oauthClientAuthenticationMethod"].(string); ok {
        data.OauthClientAuthenticationMethod = types.StringValue(val)
    } else {
        data.OauthClientAuthenticationMethod = types.StringNull()
    }
    if obj, ok := item["oauthAccessTokenExpiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthAccessTokenExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthAccessTokenExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthAccessTokenExpiresAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthAccessTokenExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.OauthAccessTokenExpiresAt = types.StringNull()
        }
    } else if val, ok := item["oauthAccessTokenExpiresAt"].(string); ok {
        data.OauthAccessTokenExpiresAt = types.StringValue(val)
    } else {
        data.OauthAccessTokenExpiresAt = types.StringNull()
    }
    if obj, ok := item["oauthLastRefreshedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthLastRefreshedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthLastRefreshedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthLastRefreshedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthLastRefreshedAt = types.StringValue(string(jsonBytes))
        } else {
            data.OauthLastRefreshedAt = types.StringNull()
        }
    } else if val, ok := item["oauthLastRefreshedAt"].(string); ok {
        data.OauthLastRefreshedAt = types.StringValue(val)
    } else {
        data.OauthLastRefreshedAt = types.StringNull()
    }
    if obj, ok := item["oauthLastRefreshError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthLastRefreshError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthLastRefreshError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthLastRefreshError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthLastRefreshError = types.StringValue(string(jsonBytes))
        } else {
            data.OauthLastRefreshError = types.StringNull()
        }
    } else if val, ok := item["oauthLastRefreshError"].(string); ok {
        data.OauthLastRefreshError = types.StringValue(val)
    } else {
        data.OauthLastRefreshError = types.StringNull()
    }
    if obj, ok := item["oauthLastRefreshErrorAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthLastRefreshErrorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthLastRefreshErrorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthLastRefreshErrorAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthLastRefreshErrorAt = types.StringValue(string(jsonBytes))
        } else {
            data.OauthLastRefreshErrorAt = types.StringNull()
        }
    } else if val, ok := item["oauthLastRefreshErrorAt"].(string); ok {
        data.OauthLastRefreshErrorAt = types.StringValue(val)
    } else {
        data.OauthLastRefreshErrorAt = types.StringNull()
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
