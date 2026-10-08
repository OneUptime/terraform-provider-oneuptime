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
var _ datasource.DataSource = &DashboardDataSource{}

func NewDashboardDataSource() datasource.DataSource {
    return &DashboardDataSource{}
}

// DashboardDataSource defines the data source implementation.
type DashboardDataSource struct {
    client *Client
}

// DashboardDataSourceModel describes the data source data model.
type DashboardDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    DashboardViewConfig types.String `tfsdk:"dashboard_view_config"`
    PageTitle types.String `tfsdk:"page_title"`
    PageDescription types.String `tfsdk:"page_description"`
    LogoFileId types.String `tfsdk:"logo_file_id"`
    FaviconFileId types.String `tfsdk:"favicon_file_id"`
    IsPublicDashboard types.Bool `tfsdk:"is_public_dashboard"`
    EnableMasterPassword types.Bool `tfsdk:"enable_master_password"`
    MasterPassword types.String `tfsdk:"master_password"`
    IpWhitelist types.String `tfsdk:"ip_whitelist"`
}

func (d *DashboardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (d *DashboardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Create and manage Dashboards to visualize your data in a single place Look up an existing dashboard by `id`, or by any of its other arguments (`name`, `archived_by_user_id`, `created_by_user_id`, ...): each one set must match, and exactly one dashboard may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Archived dashboards are hidden from the Dashboards list and their public link stops working. Unarchiving restores them as they were.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When this dashboard was archived. Empty while it is not archived.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "dashboard_view_config": schema.StringAttribute{
                MarkdownDescription: "Configuration of Dashboard View. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Title of the public dashboard page. This will be used for SEO and the browser tab.",
                Optional: true,
                Computed: true,
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Description of the public dashboard page. This will be used for SEO.",
                Optional: true,
                Computed: true,
            },
            "logo_file_id": schema.StringAttribute{
                MarkdownDescription: "Dashboard Logo File ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "favicon_file_id": schema.StringAttribute{
                MarkdownDescription: "Dashboard Favicon File ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "is_public_dashboard": schema.BoolAttribute{
                MarkdownDescription: "Is this dashboard public?",
                Optional: true,
                Computed: true,
            },
            "enable_master_password": schema.BoolAttribute{
                MarkdownDescription: "Require visitors to enter a master password before viewing a public dashboard.",
                Optional: true,
                Computed: true,
            },
            "master_password": schema.StringAttribute{
                MarkdownDescription: "Password required to unlock a public dashboard. This value is stored as a secure hash.",
                Computed: true,
                Sensitive: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist for this dashboard. One IP per line. Only used when the dashboard is public.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DashboardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DashboardDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.PageTitle.IsNull() && !data.PageTitle.IsUnknown() {
        filters["pageTitle"] = data.PageTitle.ValueString()
        filterNames = append(filterNames, "page_title = "+fmt.Sprintf("%q", data.PageTitle.ValueString()))
    }
    if !data.PageDescription.IsNull() && !data.PageDescription.IsUnknown() {
        filters["pageDescription"] = data.PageDescription.ValueString()
        filterNames = append(filterNames, "page_description = "+fmt.Sprintf("%q", data.PageDescription.ValueString()))
    }
    if !data.LogoFileId.IsNull() && !data.LogoFileId.IsUnknown() {
        filters["logoFileId"] = data.LogoFileId.ValueString()
        filterNames = append(filterNames, "logo_file_id = "+fmt.Sprintf("%q", data.LogoFileId.ValueString()))
    }
    if !data.FaviconFileId.IsNull() && !data.FaviconFileId.IsUnknown() {
        filters["faviconFileId"] = data.FaviconFileId.ValueString()
        filterNames = append(filterNames, "favicon_file_id = "+fmt.Sprintf("%q", data.FaviconFileId.ValueString()))
    }
    if !data.IsPublicDashboard.IsNull() && !data.IsPublicDashboard.IsUnknown() {
        filters["isPublicDashboard"] = data.IsPublicDashboard.ValueBool()
        filterNames = append(filterNames, "is_public_dashboard = "+fmt.Sprintf("%t", data.IsPublicDashboard.ValueBool()))
    }
    if !data.EnableMasterPassword.IsNull() && !data.EnableMasterPassword.IsUnknown() {
        filters["enableMasterPassword"] = data.EnableMasterPassword.ValueBool()
        filterNames = append(filterNames, "enable_master_password = "+fmt.Sprintf("%t", data.EnableMasterPassword.ValueBool()))
    }
    if !data.IpWhitelist.IsNull() && !data.IpWhitelist.IsUnknown() {
        filters["ipWhitelist"] = data.IpWhitelist.ValueString()
        filterNames = append(filterNames, "ip_whitelist = "+fmt.Sprintf("%q", data.IpWhitelist.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the dashboard up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the dashboard up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "labels": true,
        "dashboardViewConfig": true,
        "pageTitle": true,
        "pageDescription": true,
        "logoFileId": true,
        "faviconFileId": true,
        "isPublicDashboard": true,
        "enableMasterPassword": true,
        "masterPassword": true,
        "ipWhitelist": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/dashboard/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dashboard, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dashboard found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read dashboard: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/dashboard/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list dashboard, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list dashboard: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dashboard matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one dashboard matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for dashboard.")
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := item["labels"].([]interface{}); ok {
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
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if obj, ok := item["dashboardViewConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DashboardViewConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DashboardViewConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DashboardViewConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DashboardViewConfig = types.StringValue(string(jsonBytes))
        } else {
            data.DashboardViewConfig = types.StringNull()
        }
    } else if val, ok := item["dashboardViewConfig"].(string); ok {
        data.DashboardViewConfig = types.StringValue(val)
    } else {
        data.DashboardViewConfig = types.StringNull()
    }
    if obj, ok := item["pageTitle"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PageTitle = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PageTitle = types.StringValue(string(jsonBytes))
        } else {
            data.PageTitle = types.StringNull()
        }
    } else if val, ok := item["pageTitle"].(string); ok {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if obj, ok := item["pageDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.PageDescription = types.StringNull()
        }
    } else if val, ok := item["pageDescription"].(string); ok {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if obj, ok := item["logoFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LogoFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if val, ok := item["logoFileId"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if obj, ok := item["faviconFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if val, ok := item["faviconFileId"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if val, ok := item["isPublicDashboard"].(bool); ok {
        data.IsPublicDashboard = types.BoolValue(val)
    } else {
        data.IsPublicDashboard = types.BoolNull()
    }
    if val, ok := item["enableMasterPassword"].(bool); ok {
        data.EnableMasterPassword = types.BoolValue(val)
    } else {
        data.EnableMasterPassword = types.BoolNull()
    }
    if obj, ok := item["masterPassword"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MasterPassword = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MasterPassword = types.StringValue(string(jsonBytes))
        } else {
            data.MasterPassword = types.StringNull()
        }
    } else if val, ok := item["masterPassword"].(string); ok {
        data.MasterPassword = types.StringValue(val)
    } else {
        data.MasterPassword = types.StringNull()
    }
    if obj, ok := item["ipWhitelist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := item["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
