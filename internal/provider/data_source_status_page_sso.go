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
var _ datasource.DataSource = &StatusPageSsoDataSource{}

func NewStatusPageSsoDataSource() datasource.DataSource {
    return &StatusPageSsoDataSource{}
}

// StatusPageSsoDataSource defines the data source implementation.
type StatusPageSsoDataSource struct {
    client *Client
}

// StatusPageSsoDataSourceModel describes the data source data model.
type StatusPageSsoDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    SignatureMethod types.String `tfsdk:"signature_method"`
    DigestMethod types.String `tfsdk:"digest_method"`
    SignOnUrl types.String `tfsdk:"sign_on_url"`
    IssuerUrl types.String `tfsdk:"issuer_url"`
    PublicCertificate types.String `tfsdk:"public_certificate"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    IsTested types.Bool `tfsdk:"is_tested"`
}

func (d *StatusPageSsoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_sso"
}

func (d *StatusPageSsoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure Status Page SSO Look up an existing status page sso by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one status page sso may match them all.",

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
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "signature_method": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "digest_method": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "sign_on_url": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "issuer_url": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "public_certificate": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "is_tested": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *StatusPageSsoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageSsoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageSsoDataSourceModel

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
    if !data.SignatureMethod.IsNull() && !data.SignatureMethod.IsUnknown() {
        filters["signatureMethod"] = data.SignatureMethod.ValueString()
        filterNames = append(filterNames, "signature_method = "+fmt.Sprintf("%q", data.SignatureMethod.ValueString()))
    }
    if !data.DigestMethod.IsNull() && !data.DigestMethod.IsUnknown() {
        filters["digestMethod"] = data.DigestMethod.ValueString()
        filterNames = append(filterNames, "digest_method = "+fmt.Sprintf("%q", data.DigestMethod.ValueString()))
    }
    if !data.SignOnUrl.IsNull() && !data.SignOnUrl.IsUnknown() {
        filters["signOnURL"] = data.SignOnUrl.ValueString()
        filterNames = append(filterNames, "sign_on_url = "+fmt.Sprintf("%q", data.SignOnUrl.ValueString()))
    }
    if !data.IssuerUrl.IsNull() && !data.IssuerUrl.IsUnknown() {
        filters["issuerURL"] = data.IssuerUrl.ValueString()
        filterNames = append(filterNames, "issuer_url = "+fmt.Sprintf("%q", data.IssuerUrl.ValueString()))
    }
    if !data.PublicCertificate.IsNull() && !data.PublicCertificate.IsUnknown() {
        filters["publicCertificate"] = data.PublicCertificate.ValueString()
        filterNames = append(filterNames, "public_certificate = "+fmt.Sprintf("%q", data.PublicCertificate.ValueString()))
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
            "Look the status page sso up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page sso up by.",
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
        "signatureMethod": true,
        "digestMethod": true,
        "signOnURL": true,
        "issuerURL": true,
        "publicCertificate": true,
        "createdByUserId": true,
        "isEnabled": true,
        "isTested": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-sso/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_sso, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page sso found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_sso: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-sso/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_sso, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_sso: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page sso matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page sso matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_sso.")
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
    if obj, ok := item["signatureMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SignatureMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SignatureMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SignatureMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SignatureMethod = types.StringValue(string(jsonBytes))
        } else {
            data.SignatureMethod = types.StringNull()
        }
    } else if val, ok := item["signatureMethod"].(string); ok {
        data.SignatureMethod = types.StringValue(val)
    } else {
        data.SignatureMethod = types.StringNull()
    }
    if obj, ok := item["digestMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DigestMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DigestMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DigestMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DigestMethod = types.StringValue(string(jsonBytes))
        } else {
            data.DigestMethod = types.StringNull()
        }
    } else if val, ok := item["digestMethod"].(string); ok {
        data.DigestMethod = types.StringValue(val)
    } else {
        data.DigestMethod = types.StringNull()
    }
    if obj, ok := item["signOnURL"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SignOnUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SignOnUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SignOnUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SignOnUrl = types.StringValue(string(jsonBytes))
        } else {
            data.SignOnUrl = types.StringNull()
        }
    } else if val, ok := item["signOnURL"].(string); ok {
        data.SignOnUrl = types.StringValue(val)
    } else {
        data.SignOnUrl = types.StringNull()
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
    if obj, ok := item["publicCertificate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PublicCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PublicCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PublicCertificate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PublicCertificate = types.StringValue(string(jsonBytes))
        } else {
            data.PublicCertificate = types.StringNull()
        }
    } else if val, ok := item["publicCertificate"].(string); ok {
        data.PublicCertificate = types.StringValue(val)
    } else {
        data.PublicCertificate = types.StringNull()
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
