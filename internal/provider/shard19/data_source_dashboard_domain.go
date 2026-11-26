package shard19

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
var _ datasource.DataSource = &DashboardDomainDataSource{}

func NewDashboardDomainDataSource() datasource.DataSource {
    return &DashboardDomainDataSource{}
}

// DashboardDomainDataSource defines the data source implementation.
type DashboardDomainDataSource struct {
    client *Client
}

// DashboardDomainDataSourceModel describes the data source data model.
type DashboardDomainDataSourceModel struct {
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
    DomainId types.String `tfsdk:"domain_id"`
    DashboardId types.String `tfsdk:"dashboard_id"`
    Subdomain types.String `tfsdk:"subdomain"`
    FullDomain types.String `tfsdk:"full_domain"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsCnameVerified types.Bool `tfsdk:"is_cname_verified"`
    IsSslOrdered types.Bool `tfsdk:"is_ssl_ordered"`
    IsSslProvisioned types.Bool `tfsdk:"is_ssl_provisioned"`
    CustomCertificate types.String `tfsdk:"custom_certificate"`
    IsCustomCertificate types.Bool `tfsdk:"is_custom_certificate"`
    CertificateReissueRequestedAt types.String `tfsdk:"certificate_reissue_requested_at"`
}

func (d *DashboardDomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_dashboard_domain"
}

func (d *DashboardDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage custom domains for your dashboard Look up an existing dashboard domain by `id`, or by any of its other arguments (`created_by_actor_id`, `created_by_actor_type`, `created_by_user_id`, ...): each one set must match, and exactly one dashboard domain may match them all.",

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
            "domain_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Dashboard Domain, Dashboard Admin, Dashboard Member], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Dashboard Domain, Dashboard Admin, Dashboard Member, Dashboard Viewer], Update: [No access - you don't have permission for this operation] The ID of a `oneuptime_domain`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_domain"),
            },
            "dashboard_id": schema.StringAttribute{
                MarkdownDescription: "ID of your Dashboard resource where this object belongs. The ID of a `oneuptime_dashboard`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_dashboard"),
            },
            "subdomain": schema.StringAttribute{
                MarkdownDescription: "Subdomain label for your dashboard such as 'dashboard'. Leave blank or enter @ to use the root domain.",
                Optional: true,
                Computed: true,
            },
            "full_domain": schema.StringAttribute{
                MarkdownDescription: "Full domain of your dashboard (like dashboard.acmeinc.com). This is autogenerated and is derived from subdomain and domain.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_cname_verified": schema.BoolAttribute{
                MarkdownDescription: "Is CNAME Verified?",
                Optional: true,
                Computed: true,
            },
            "is_ssl_ordered": schema.BoolAttribute{
                MarkdownDescription: "Is SSL ordered?",
                Optional: true,
                Computed: true,
            },
            "is_ssl_provisioned": schema.BoolAttribute{
                MarkdownDescription: "Is SSL provisioned?",
                Optional: true,
                Computed: true,
            },
            "custom_certificate": schema.StringAttribute{
                MarkdownDescription: "Your certificate for the domain, in PEM, with its chain. Used when Is Custom Certificate is on.",
                Optional: true,
                Computed: true,
            },
            "is_custom_certificate": schema.BoolAttribute{
                MarkdownDescription: "When on, the domain is served with your own certificate (custom_certificate and its key) instead of one OneUptime gets from Let's Encrypt.",
                Optional: true,
                Computed: true,
            },
            "certificate_reissue_requested_at": schema.StringAttribute{
                MarkdownDescription: "When a reissue of the SSL certificate for this domain was last requested. Used to rate limit how often a certificate can be reissued.",
                Computed: true,
            },
        },
    }
}

func (d *DashboardDomainDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DashboardDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DashboardDomainDataSourceModel

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
    if !data.DomainId.IsNull() && !data.DomainId.IsUnknown() {
        filters["domainId"] = data.DomainId.ValueString()
        filterNames = append(filterNames, "domain_id = "+fmt.Sprintf("%q", data.DomainId.ValueString()))
    }
    if !data.DashboardId.IsNull() && !data.DashboardId.IsUnknown() {
        filters["dashboardId"] = data.DashboardId.ValueString()
        filterNames = append(filterNames, "dashboard_id = "+fmt.Sprintf("%q", data.DashboardId.ValueString()))
    }
    if !data.Subdomain.IsNull() && !data.Subdomain.IsUnknown() {
        filters["subdomain"] = data.Subdomain.ValueString()
        filterNames = append(filterNames, "subdomain = "+fmt.Sprintf("%q", data.Subdomain.ValueString()))
    }
    if !data.FullDomain.IsNull() && !data.FullDomain.IsUnknown() {
        filters["fullDomain"] = data.FullDomain.ValueString()
        filterNames = append(filterNames, "full_domain = "+fmt.Sprintf("%q", data.FullDomain.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsCnameVerified.IsNull() && !data.IsCnameVerified.IsUnknown() {
        filters["isCnameVerified"] = data.IsCnameVerified.ValueBool()
        filterNames = append(filterNames, "is_cname_verified = "+fmt.Sprintf("%t", data.IsCnameVerified.ValueBool()))
    }
    if !data.IsSslOrdered.IsNull() && !data.IsSslOrdered.IsUnknown() {
        filters["isSslOrdered"] = data.IsSslOrdered.ValueBool()
        filterNames = append(filterNames, "is_ssl_ordered = "+fmt.Sprintf("%t", data.IsSslOrdered.ValueBool()))
    }
    if !data.IsSslProvisioned.IsNull() && !data.IsSslProvisioned.IsUnknown() {
        filters["isSslProvisioned"] = data.IsSslProvisioned.ValueBool()
        filterNames = append(filterNames, "is_ssl_provisioned = "+fmt.Sprintf("%t", data.IsSslProvisioned.ValueBool()))
    }
    if !data.CustomCertificate.IsNull() && !data.CustomCertificate.IsUnknown() {
        filters["customCertificate"] = data.CustomCertificate.ValueString()
        filterNames = append(filterNames, "custom_certificate = "+fmt.Sprintf("%q", data.CustomCertificate.ValueString()))
    }
    if !data.IsCustomCertificate.IsNull() && !data.IsCustomCertificate.IsUnknown() {
        filters["isCustomCertificate"] = data.IsCustomCertificate.ValueBool()
        filterNames = append(filterNames, "is_custom_certificate = "+fmt.Sprintf("%t", data.IsCustomCertificate.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the dashboard domain up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the dashboard domain up by.",
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
        "domainId": true,
        "dashboardId": true,
        "subdomain": true,
        "fullDomain": true,
        "createdByUserId": true,
        "isCnameVerified": true,
        "isSslOrdered": true,
        "isSslProvisioned": true,
        "customCertificate": true,
        "isCustomCertificate": true,
        "certificateReissueRequestedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/dashboard-domain/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dashboard_domain, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dashboard domain found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read dashboard_domain: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/dashboard-domain/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list dashboard_domain, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list dashboard_domain: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dashboard domain matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one dashboard domain matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for dashboard_domain.")
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
    if obj, ok := item["domainId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DomainId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DomainId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DomainId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DomainId = types.StringValue(string(jsonBytes))
        } else {
            data.DomainId = types.StringNull()
        }
    } else if arr, ok := item["domainId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DomainId = types.StringValue(string(jsonBytes))
        } else {
            data.DomainId = types.StringNull()
        }
    } else if val, ok := item["domainId"].(string); ok {
        data.DomainId = types.StringValue(val)
    } else {
        data.DomainId = types.StringNull()
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
    if obj, ok := item["subdomain"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Subdomain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Subdomain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Subdomain = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Subdomain = types.StringValue(string(jsonBytes))
        } else {
            data.Subdomain = types.StringNull()
        }
    } else if arr, ok := item["subdomain"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Subdomain = types.StringValue(string(jsonBytes))
        } else {
            data.Subdomain = types.StringNull()
        }
    } else if val, ok := item["subdomain"].(string); ok {
        data.Subdomain = types.StringValue(val)
    } else {
        data.Subdomain = types.StringNull()
    }
    if obj, ok := item["fullDomain"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FullDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FullDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FullDomain = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FullDomain = types.StringValue(string(jsonBytes))
        } else {
            data.FullDomain = types.StringNull()
        }
    } else if arr, ok := item["fullDomain"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FullDomain = types.StringValue(string(jsonBytes))
        } else {
            data.FullDomain = types.StringNull()
        }
    } else if val, ok := item["fullDomain"].(string); ok {
        data.FullDomain = types.StringValue(val)
    } else {
        data.FullDomain = types.StringNull()
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
    if val, ok := item["isCnameVerified"].(bool); ok {
        data.IsCnameVerified = types.BoolValue(val)
    } else {
        data.IsCnameVerified = types.BoolNull()
    }
    if val, ok := item["isSslOrdered"].(bool); ok {
        data.IsSslOrdered = types.BoolValue(val)
    } else {
        data.IsSslOrdered = types.BoolNull()
    }
    if val, ok := item["isSslProvisioned"].(bool); ok {
        data.IsSslProvisioned = types.BoolValue(val)
    } else {
        data.IsSslProvisioned = types.BoolNull()
    }
    if obj, ok := item["customCertificate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomCertificate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomCertificate = types.StringValue(string(jsonBytes))
        } else {
            data.CustomCertificate = types.StringNull()
        }
    } else if arr, ok := item["customCertificate"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CustomCertificate = types.StringValue(string(jsonBytes))
        } else {
            data.CustomCertificate = types.StringNull()
        }
    } else if val, ok := item["customCertificate"].(string); ok {
        data.CustomCertificate = types.StringValue(val)
    } else {
        data.CustomCertificate = types.StringNull()
    }
    if val, ok := item["isCustomCertificate"].(bool); ok {
        data.IsCustomCertificate = types.BoolValue(val)
    } else {
        data.IsCustomCertificate = types.BoolNull()
    }
    if obj, ok := item["certificateReissueRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CertificateReissueRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CertificateReissueRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CertificateReissueRequestedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CertificateReissueRequestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CertificateReissueRequestedAt = types.StringNull()
        }
    } else if arr, ok := item["certificateReissueRequestedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CertificateReissueRequestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CertificateReissueRequestedAt = types.StringNull()
        }
    } else if val, ok := item["certificateReissueRequestedAt"].(string); ok {
        data.CertificateReissueRequestedAt = types.StringValue(val)
    } else {
        data.CertificateReissueRequestedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
