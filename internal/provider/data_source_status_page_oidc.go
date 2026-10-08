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
var _ datasource.DataSource = &StatusPageOidcDataSource{}

func NewStatusPageOidcDataSource() datasource.DataSource {
    return &StatusPageOidcDataSource{}
}

// StatusPageOidcDataSource defines the data source implementation.
type StatusPageOidcDataSource struct {
    client *Client
}

// StatusPageOidcDataSourceModel describes the data source data model.
type StatusPageOidcDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
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

func (d *StatusPageOidcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_oidc"
}

func (d *StatusPageOidcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage OpenID Connect (OIDC) authentication for your status page Look up an existing status page oidc by `id`, or by any of its other arguments (`name`, `client_id`, `client_secret`, ...): each one set must match, and exactly one status page oidc may match them all.",

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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page OIDC], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page OIDC], Update: [Project Owner, Project Admin, Edit Status Page OIDC]",
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
                Optional: true,
                Computed: true,
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
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page OIDC], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page OIDC], Update: [Project Owner, Project Admin, Edit Status Page OIDC]",
                Optional: true,
                Computed: true,
            },
            "is_tested": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page OIDC], Read: [Project Owner, Project Admin, Read Status Page OIDC], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *StatusPageOidcDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageOidcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageOidcDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.StatusPageId.IsNull() && !data.StatusPageId.IsUnknown() {
        filters["statusPageId"] = data.StatusPageId.ValueString()
        filterNames = append(filterNames, "status_page_id = "+fmt.Sprintf("%q", data.StatusPageId.ValueString()))
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
    if !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown() {
        filters["clientSecret"] = data.ClientSecret.ValueString()
        filterNames = append(filterNames, "client_secret = "+fmt.Sprintf("%q", data.ClientSecret.ValueString()))
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
            "Look the status page oidc up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page oidc up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "statusPageId": true,
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
        readPath := "/status-page-oidc/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_oidc, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page oidc found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_oidc: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-oidc/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_oidc, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_oidc: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page oidc matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page oidc matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_oidc.")
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
    if obj, ok := item["statusPageId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := item["statusPageId"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
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
