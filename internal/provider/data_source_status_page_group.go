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
var _ datasource.DataSource = &StatusPageGroupDataSource{}

func NewStatusPageGroupDataSource() datasource.DataSource {
    return &StatusPageGroupDataSource{}
}

// StatusPageGroupDataSource defines the data source implementation.
type StatusPageGroupDataSource struct {
    client *Client
}

// StatusPageGroupDataSourceModel describes the data source data model.
type StatusPageGroupDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    ParentStatusPageGroupId types.String `tfsdk:"parent_status_page_group_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Order types.Number `tfsdk:"order"`
    IsExpandedByDefault types.Bool `tfsdk:"is_expanded_by_default"`
    ShowCurrentStatus types.Bool `tfsdk:"show_current_status"`
    ShowUptimePercent types.Bool `tfsdk:"show_uptime_percent"`
    UptimePercentPrecision types.String `tfsdk:"uptime_percent_precision"`
    ViewMode types.String `tfsdk:"view_mode"`
    RowAxisLabel types.String `tfsdk:"row_axis_label"`
    ColumnAxisLabel types.String `tfsdk:"column_axis_label"`
    RowAxisValues types.String `tfsdk:"row_axis_values"`
    ColumnAxisValues types.String `tfsdk:"column_axis_values"`
}

func (d *StatusPageGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_group"
}

func (d *StatusPageGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage groups on your status page and categorize resources like monitors into these groups. Look up an existing status page group by `id`, or by any of its other arguments (`name`, `column_axis_label`, `column_axis_values`, ...): each one set must match, and exactly one status page group may match them all.",

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
            "parent_status_page_group_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Status Page Group this group is nested under. Empty for top level groups. The ID of a `oneuptime_status_page_group`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of the Group.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description for this group. This is visible on Status Page. This can be in markdown format.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order / Priority of this resource.",
                Optional: true,
                Computed: true,
            },
            "is_expanded_by_default": schema.BoolAttribute{
                MarkdownDescription: "Is this group expanded by default on Status Page?",
                Optional: true,
                Computed: true,
            },
            "show_current_status": schema.BoolAttribute{
                MarkdownDescription: "Show current status like offline, operational or degraded.",
                Optional: true,
                Computed: true,
            },
            "show_uptime_percent": schema.BoolAttribute{
                MarkdownDescription: "Show uptime percent of this group for the last 90 days.",
                Optional: true,
                Computed: true,
            },
            "uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Precision of uptime percent of this group for the last 90 days.",
                Optional: true,
                Computed: true,
            },
            "view_mode": schema.StringAttribute{
                MarkdownDescription: "Layout of this group on the status page. 'List' renders resources stacked vertically (default). 'Grid' renders resources as a matrix using row and column axes.",
                Optional: true,
                Computed: true,
            },
            "row_axis_label": schema.StringAttribute{
                MarkdownDescription: "Label shown above the row axis when the group is rendered as a grid (e.g. 'Service', 'Tenant'). Free-form so you can use any dimension you like.",
                Optional: true,
                Computed: true,
            },
            "column_axis_label": schema.StringAttribute{
                MarkdownDescription: "Label shown above the column axis when the group is rendered as a grid (e.g. 'Region', 'Environment'). Free-form so you can use any dimension you like.",
                Optional: true,
                Computed: true,
            },
            "row_axis_values": schema.StringAttribute{
                MarkdownDescription: "Comma-separated list of row labels for the grid (e.g. 'Auth, API, Database'). Determines row order in the grid layout.",
                Optional: true,
                Computed: true,
            },
            "column_axis_values": schema.StringAttribute{
                MarkdownDescription: "Comma-separated list of column labels for the grid (e.g. 'US-East, EU-West, Asia'). Determines column order in the grid layout.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *StatusPageGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageGroupDataSourceModel

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
    if !data.ParentStatusPageGroupId.IsNull() && !data.ParentStatusPageGroupId.IsUnknown() {
        filters["parentStatusPageGroupId"] = data.ParentStatusPageGroupId.ValueString()
        filterNames = append(filterNames, "parent_status_page_group_id = "+fmt.Sprintf("%q", data.ParentStatusPageGroupId.ValueString()))
    }
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
    if !data.Order.IsNull() && !data.Order.IsUnknown() {
        filters["order"] = lookupNumber(data.Order)
        filterNames = append(filterNames, "order = "+data.Order.ValueBigFloat().String())
    }
    if !data.IsExpandedByDefault.IsNull() && !data.IsExpandedByDefault.IsUnknown() {
        filters["isExpandedByDefault"] = data.IsExpandedByDefault.ValueBool()
        filterNames = append(filterNames, "is_expanded_by_default = "+fmt.Sprintf("%t", data.IsExpandedByDefault.ValueBool()))
    }
    if !data.ShowCurrentStatus.IsNull() && !data.ShowCurrentStatus.IsUnknown() {
        filters["showCurrentStatus"] = data.ShowCurrentStatus.ValueBool()
        filterNames = append(filterNames, "show_current_status = "+fmt.Sprintf("%t", data.ShowCurrentStatus.ValueBool()))
    }
    if !data.ShowUptimePercent.IsNull() && !data.ShowUptimePercent.IsUnknown() {
        filters["showUptimePercent"] = data.ShowUptimePercent.ValueBool()
        filterNames = append(filterNames, "show_uptime_percent = "+fmt.Sprintf("%t", data.ShowUptimePercent.ValueBool()))
    }
    if !data.UptimePercentPrecision.IsNull() && !data.UptimePercentPrecision.IsUnknown() {
        filters["uptimePercentPrecision"] = data.UptimePercentPrecision.ValueString()
        filterNames = append(filterNames, "uptime_percent_precision = "+fmt.Sprintf("%q", data.UptimePercentPrecision.ValueString()))
    }
    if !data.ViewMode.IsNull() && !data.ViewMode.IsUnknown() {
        filters["viewMode"] = data.ViewMode.ValueString()
        filterNames = append(filterNames, "view_mode = "+fmt.Sprintf("%q", data.ViewMode.ValueString()))
    }
    if !data.RowAxisLabel.IsNull() && !data.RowAxisLabel.IsUnknown() {
        filters["rowAxisLabel"] = data.RowAxisLabel.ValueString()
        filterNames = append(filterNames, "row_axis_label = "+fmt.Sprintf("%q", data.RowAxisLabel.ValueString()))
    }
    if !data.ColumnAxisLabel.IsNull() && !data.ColumnAxisLabel.IsUnknown() {
        filters["columnAxisLabel"] = data.ColumnAxisLabel.ValueString()
        filterNames = append(filterNames, "column_axis_label = "+fmt.Sprintf("%q", data.ColumnAxisLabel.ValueString()))
    }
    if !data.RowAxisValues.IsNull() && !data.RowAxisValues.IsUnknown() {
        filters["rowAxisValues"] = data.RowAxisValues.ValueString()
        filterNames = append(filterNames, "row_axis_values = "+fmt.Sprintf("%q", data.RowAxisValues.ValueString()))
    }
    if !data.ColumnAxisValues.IsNull() && !data.ColumnAxisValues.IsUnknown() {
        filters["columnAxisValues"] = data.ColumnAxisValues.ValueString()
        filterNames = append(filterNames, "column_axis_values = "+fmt.Sprintf("%q", data.ColumnAxisValues.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the status page group up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page group up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "statusPageId": true,
        "parentStatusPageGroupId": true,
        "name": true,
        "slug": true,
        "description": true,
        "createdByUserId": true,
        "order": true,
        "isExpandedByDefault": true,
        "showCurrentStatus": true,
        "showUptimePercent": true,
        "uptimePercentPrecision": true,
        "viewMode": true,
        "rowAxisLabel": true,
        "columnAxisLabel": true,
        "rowAxisValues": true,
        "columnAxisValues": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-group/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_group, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page group found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_group: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-group/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_group, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_group: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page group matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page group matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_group.")
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
    if obj, ok := item["parentStatusPageGroupId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ParentStatusPageGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ParentStatusPageGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ParentStatusPageGroupId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ParentStatusPageGroupId = types.StringValue(string(jsonBytes))
        } else {
            data.ParentStatusPageGroupId = types.StringNull()
        }
    } else if val, ok := item["parentStatusPageGroupId"].(string); ok {
        data.ParentStatusPageGroupId = types.StringValue(val)
    } else {
        data.ParentStatusPageGroupId = types.StringNull()
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
    if val, ok := item["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["order"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Order = types.NumberValue(big.NewFloat(val))
        } else {
            data.Order = types.NumberNull()
        }
    } else {
        data.Order = types.NumberNull()
    }
    if val, ok := item["isExpandedByDefault"].(bool); ok {
        data.IsExpandedByDefault = types.BoolValue(val)
    } else {
        data.IsExpandedByDefault = types.BoolNull()
    }
    if val, ok := item["showCurrentStatus"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    } else {
        data.ShowCurrentStatus = types.BoolNull()
    }
    if val, ok := item["showUptimePercent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    } else {
        data.ShowUptimePercent = types.BoolNull()
    }
    if obj, ok := item["uptimePercentPrecision"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UptimePercentPrecision = types.StringValue(string(jsonBytes))
        } else {
            data.UptimePercentPrecision = types.StringNull()
        }
    } else if val, ok := item["uptimePercentPrecision"].(string); ok {
        data.UptimePercentPrecision = types.StringValue(val)
    } else {
        data.UptimePercentPrecision = types.StringNull()
    }
    if obj, ok := item["viewMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ViewMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ViewMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ViewMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ViewMode = types.StringValue(string(jsonBytes))
        } else {
            data.ViewMode = types.StringNull()
        }
    } else if val, ok := item["viewMode"].(string); ok {
        data.ViewMode = types.StringValue(val)
    } else {
        data.ViewMode = types.StringNull()
    }
    if obj, ok := item["rowAxisLabel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RowAxisLabel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RowAxisLabel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RowAxisLabel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RowAxisLabel = types.StringValue(string(jsonBytes))
        } else {
            data.RowAxisLabel = types.StringNull()
        }
    } else if val, ok := item["rowAxisLabel"].(string); ok {
        data.RowAxisLabel = types.StringValue(val)
    } else {
        data.RowAxisLabel = types.StringNull()
    }
    if obj, ok := item["columnAxisLabel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ColumnAxisLabel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ColumnAxisLabel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ColumnAxisLabel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ColumnAxisLabel = types.StringValue(string(jsonBytes))
        } else {
            data.ColumnAxisLabel = types.StringNull()
        }
    } else if val, ok := item["columnAxisLabel"].(string); ok {
        data.ColumnAxisLabel = types.StringValue(val)
    } else {
        data.ColumnAxisLabel = types.StringNull()
    }
    if obj, ok := item["rowAxisValues"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RowAxisValues = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RowAxisValues = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RowAxisValues = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RowAxisValues = types.StringValue(string(jsonBytes))
        } else {
            data.RowAxisValues = types.StringNull()
        }
    } else if val, ok := item["rowAxisValues"].(string); ok {
        data.RowAxisValues = types.StringValue(val)
    } else {
        data.RowAxisValues = types.StringNull()
    }
    if obj, ok := item["columnAxisValues"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ColumnAxisValues = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ColumnAxisValues = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ColumnAxisValues = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ColumnAxisValues = types.StringValue(string(jsonBytes))
        } else {
            data.ColumnAxisValues = types.StringNull()
        }
    } else if val, ok := item["columnAxisValues"].(string); ok {
        data.ColumnAxisValues = types.StringValue(val)
    } else {
        data.ColumnAxisValues = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
