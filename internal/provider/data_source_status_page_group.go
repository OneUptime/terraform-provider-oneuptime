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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    ParentStatusPageGroupId types.String `tfsdk:"parent_status_page_group_id"`
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
        MarkdownDescription: "Manage groups on your status page and categorize resources like monitors into these groups. Look up an existing status_page_group by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "parent_status_page_group_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description for this group. This is visible on Status Page. This can be in markdown format..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order / Priority of this resource.",
                Computed: true,
            },
            "is_expanded_by_default": schema.BoolAttribute{
                MarkdownDescription: "Is this group expanded by default on Status Page?.",
                Computed: true,
            },
            "show_current_status": schema.BoolAttribute{
                MarkdownDescription: "Show current status like offline, operational or degraded..",
                Computed: true,
            },
            "show_uptime_percent": schema.BoolAttribute{
                MarkdownDescription: "Show uptime percent of this group for the last 90 days.",
                Computed: true,
            },
            "uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Precision of uptime percent of this group for the last 90 days.",
                Computed: true,
            },
            "view_mode": schema.StringAttribute{
                MarkdownDescription: "Layout of this group on the status page. 'List' renders resources stacked vertically (default). 'Grid' renders resources as a matrix using row and column axes..",
                Computed: true,
            },
            "row_axis_label": schema.StringAttribute{
                MarkdownDescription: "Label shown above the row axis when the group is rendered as a grid (e.g. 'Service', 'Tenant'). Free-form so you can use any dimension you like..",
                Computed: true,
            },
            "column_axis_label": schema.StringAttribute{
                MarkdownDescription: "Label shown above the column axis when the group is rendered as a grid (e.g. 'Region', 'Environment'). Free-form so you can use any dimension you like..",
                Computed: true,
            },
            "row_axis_values": schema.StringAttribute{
                MarkdownDescription: "Comma-separated list of row labels for the grid (e.g. 'Auth, API, Database'). Determines row order in the grid layout..",
                Computed: true,
            },
            "column_axis_values": schema.StringAttribute{
                MarkdownDescription: "Comma-separated list of column labels for the grid (e.g. 'US-East, EU-West, Asia'). Determines column order in the grid layout..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a status_page_group.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "statusPageId": true,
        "parentStatusPageGroupId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status_page_group found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status_page_group found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status_page_group matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
