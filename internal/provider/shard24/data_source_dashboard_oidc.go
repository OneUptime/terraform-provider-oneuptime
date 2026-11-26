package shard24

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
var _ datasource.DataSource = &DashboardOidcDataSource{}

func NewDashboardOidcDataSource() datasource.DataSource {
    return &DashboardOidcDataSource{}
}

// DashboardOidcDataSource defines the data source implementation.
type DashboardOidcDataSource struct {
    client *Client
}

// DashboardOidcDataSourceModel describes the data source data model.
type DashboardOidcDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    DashboardId types.String `tfsdk:"dashboard_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    DiscoveryUrl types.String `tfsdk:"discovery_url"`
    IssuerUrl types.String `tfsdk:"issuer_url"`
    ClientId types.String `tfsdk:"client_id"`
    ClientSecret types.String `tfsdk:"client_secret"`
    Scopes types.String `tfsdk:"scopes"`
    EmailClaimName types.String `tfsdk:"email_claim_name"`
    NameClaimName types.String `tfsdk:"name_claim_name"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    IsTested types.Bool `tfsdk:"is_tested"`
}

func (d *DashboardOidcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_dashboard_oidc"
}

func (d *DashboardOidcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage OpenID Connect (OIDC) sign-in for a dashboard that asks its visitors to sign in Look up an existing dashboard oidc by `id`, or by any of its other arguments (`name`, `client_id`, `created_by_actor_id`, ...): each one set must match, and exactly one dashboard oidc may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "dashboard_id": schema.StringAttribute{
                MarkdownDescription: "ID of the dashboard this provider signs people in to. The ID of a `oneuptime_dashboard`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_dashboard"),
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What this sign-in is for, such as which identity provider it uses.",
                Optional: true,
                Computed: true,
            },
            "discovery_url": schema.StringAttribute{
                MarkdownDescription: "OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints.",
                Optional: true,
                Computed: true,
            },
            "issuer_url": schema.StringAttribute{
                MarkdownDescription: "Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider.",
                Optional: true,
                Computed: true,
            },
            "client_id": schema.StringAttribute{
                MarkdownDescription: "OIDC client ID issued by the identity provider.",
                Optional: true,
                Computed: true,
            },
            "client_secret": schema.StringAttribute{
                MarkdownDescription: "OIDC client secret issued by the identity provider. Stored encrypted at rest.",
                Computed: true,
                Sensitive: true,
            },
            "scopes": schema.StringAttribute{
                MarkdownDescription: "Space-separated list of OIDC scopes to request. Must include 'openid'.",
                Optional: true,
                Computed: true,
            },
            "email_claim_name": schema.StringAttribute{
                MarkdownDescription: "Claim name in the ID token (or userinfo response) that contains the user's email address.",
                Optional: true,
                Computed: true,
            },
            "name_claim_name": schema.StringAttribute{
                MarkdownDescription: "Claim name in the ID token (or userinfo response) that contains the user's display name.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, visitors of a dashboard that asks them to sign in can sign in with this provider.",
                Optional: true,
                Computed: true,
            },
            "is_tested": schema.BoolAttribute{
                MarkdownDescription: "Whether a sign-in with this provider has been tested. It does not turn the provider on.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DashboardOidcDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DashboardOidcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DashboardOidcDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.DashboardId.IsNull() && !data.DashboardId.IsUnknown() {
        filters["dashboardId"] = data.DashboardId.ValueString()
        filterNames = append(filterNames, "dashboard_id = "+fmt.Sprintf("%q", data.DashboardId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.DiscoveryUrl.IsNull() && !data.DiscoveryUrl.IsUnknown() {
        filters["discoveryURL"] = data.DiscoveryUrl.ValueString()
        filterNames = append(filterNames, "discovery_url = "+fmt.Sprintf("%q", data.DiscoveryUrl.ValueString()))
    }
    if !data.IssuerUrl.IsNull() && !data.IssuerUrl.IsUnknown() {
        filters["issuerURL"] = data.IssuerUrl.ValueString()
        filterNames = append(filterNames, "issuer_url = "+fmt.Sprintf("%q", data.IssuerUrl.ValueString()))
    }
    if !data.ClientId.IsNull() && !data.ClientId.IsUnknown() {
        filters["clientId"] = data.ClientId.ValueString()
        filterNames = append(filterNames, "client_id = "+fmt.Sprintf("%q", data.ClientId.ValueString()))
    }
    if !data.Scopes.IsNull() && !data.Scopes.IsUnknown() {
        filters["scopes"] = data.Scopes.ValueString()
        filterNames = append(filterNames, "scopes = "+fmt.Sprintf("%q", data.Scopes.ValueString()))
    }
    if !data.EmailClaimName.IsNull() && !data.EmailClaimName.IsUnknown() {
        filters["emailClaimName"] = data.EmailClaimName.ValueString()
        filterNames = append(filterNames, "email_claim_name = "+fmt.Sprintf("%q", data.EmailClaimName.ValueString()))
    }
    if !data.NameClaimName.IsNull() && !data.NameClaimName.IsUnknown() {
        filters["nameClaimName"] = data.NameClaimName.ValueString()
        filterNames = append(filterNames, "name_claim_name = "+fmt.Sprintf("%q", data.NameClaimName.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.IsTested.IsNull() && !data.IsTested.IsUnknown() {
        filters["isTested"] = data.IsTested.ValueBool()
        filterNames = append(filterNames, "is_tested = "+fmt.Sprintf("%t", data.IsTested.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the dashboard oidc up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the dashboard oidc up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "dashboardId": true,
        "name": true,
        "description": true,
        "discoveryURL": true,
        "issuerURL": true,
        "clientId": true,
        "clientSecret": true,
        "scopes": true,
        "emailClaimName": true,
        "nameClaimName": true,
        "createdByUserId": true,
        "isEnabled": true,
        "isTested": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/dashboard-oidc/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dashboard_oidc, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dashboard oidc found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read dashboard_oidc: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/dashboard-oidc/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list dashboard_oidc, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list dashboard_oidc: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dashboard oidc matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one dashboard oidc matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for dashboard_oidc.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["dashboardId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DashboardId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DashboardId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DashboardId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DashboardId = types.StringValue(string(jsonBytes))
        } else {
            data.DashboardId = types.StringNull()
        }
    } else if arr, ok := item["dashboardId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DashboardId = types.StringValue(string(jsonBytes))
        } else {
            data.DashboardId = types.StringNull()
        }
    } else if val, ok := item["dashboardId"].(string); ok {
        data.DashboardId = types.StringValue(val)
    } else {
        data.DashboardId = types.StringNull()
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
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["discoveryURL"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveryUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DiscoveryUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DiscoveryUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DiscoveryUrl = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveryUrl = types.StringNull()
        }
    } else if arr, ok := item["discoveryURL"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DiscoveryUrl = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveryUrl = types.StringNull()
        }
    } else if val, ok := item["discoveryURL"].(string); ok {
        data.DiscoveryUrl = types.StringValue(val)
    } else {
        data.DiscoveryUrl = types.StringNull()
    }
    if obj, ok := item["issuerURL"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IssuerUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IssuerUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IssuerUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IssuerUrl = types.StringValue(string(jsonBytes))
        } else {
            data.IssuerUrl = types.StringNull()
        }
    } else if arr, ok := item["issuerURL"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IssuerUrl = types.StringValue(string(jsonBytes))
        } else {
            data.IssuerUrl = types.StringNull()
        }
    } else if val, ok := item["issuerURL"].(string); ok {
        data.IssuerUrl = types.StringValue(val)
    } else {
        data.IssuerUrl = types.StringNull()
    }
    if obj, ok := item["clientId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClientId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if arr, ok := item["clientId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if val, ok := item["clientId"].(string); ok {
        data.ClientId = types.StringValue(val)
    } else {
        data.ClientId = types.StringNull()
    }
    if obj, ok := item["clientSecret"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientSecret = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClientSecret = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClientSecret = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClientSecret = types.StringValue(string(jsonBytes))
        } else {
            data.ClientSecret = types.StringNull()
        }
    } else if arr, ok := item["clientSecret"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ClientSecret = types.StringValue(string(jsonBytes))
        } else {
            data.ClientSecret = types.StringNull()
        }
    } else if val, ok := item["clientSecret"].(string); ok {
        data.ClientSecret = types.StringValue(val)
    } else {
        data.ClientSecret = types.StringNull()
    }
    if obj, ok := item["scopes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Scopes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Scopes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Scopes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Scopes = types.StringValue(string(jsonBytes))
        } else {
            data.Scopes = types.StringNull()
        }
    } else if arr, ok := item["scopes"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Scopes = types.StringValue(string(jsonBytes))
        } else {
            data.Scopes = types.StringNull()
        }
    } else if val, ok := item["scopes"].(string); ok {
        data.Scopes = types.StringValue(val)
    } else {
        data.Scopes = types.StringNull()
    }
    if obj, ok := item["emailClaimName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmailClaimName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EmailClaimName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EmailClaimName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EmailClaimName = types.StringValue(string(jsonBytes))
        } else {
            data.EmailClaimName = types.StringNull()
        }
    } else if arr, ok := item["emailClaimName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EmailClaimName = types.StringValue(string(jsonBytes))
        } else {
            data.EmailClaimName = types.StringNull()
        }
    } else if val, ok := item["emailClaimName"].(string); ok {
        data.EmailClaimName = types.StringValue(val)
    } else {
        data.EmailClaimName = types.StringNull()
    }
    if obj, ok := item["nameClaimName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NameClaimName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NameClaimName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NameClaimName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NameClaimName = types.StringValue(string(jsonBytes))
        } else {
            data.NameClaimName = types.StringNull()
        }
    } else if arr, ok := item["nameClaimName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.NameClaimName = types.StringValue(string(jsonBytes))
        } else {
            data.NameClaimName = types.StringNull()
        }
    } else if val, ok := item["nameClaimName"].(string); ok {
        data.NameClaimName = types.StringValue(val)
    } else {
        data.NameClaimName = types.StringNull()
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["isTested"].(bool); ok {
        data.IsTested = types.BoolValue(val)
    } else {
        data.IsTested = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
