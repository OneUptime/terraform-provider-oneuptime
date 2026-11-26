package shard28

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
var _ datasource.DataSource = &BiToolDataSource{}

func NewBiToolDataSource() datasource.DataSource {
    return &BiToolDataSource{}
}

// BiToolDataSource defines the data source implementation.
type BiToolDataSource struct {
    client *Client
}

// BiToolDataSourceModel describes the data source data model.
type BiToolDataSourceModel struct {
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
    Name types.String `tfsdk:"name"`
    Tool types.String `tfsdk:"tool"`
    BaseUrl types.String `tfsdk:"base_url"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    LastSyncedAt types.String `tfsdk:"last_synced_at"`
    LastSyncError types.String `tfsdk:"last_sync_error"`
    DashboardCount types.Number `tfsdk:"dashboard_count"`
    MatchedTableCount types.Number `tfsdk:"matched_table_count"`
    LineageSummary types.String `tfsdk:"lineage_summary"`
    LastSyncWarning types.String `tfsdk:"last_sync_warning"`
    UsageReadFrom types.String `tfsdk:"usage_read_from"`
    UsageReadUntil types.String `tfsdk:"usage_read_until"`
    AlertOnFailedRefreshes types.Bool `tfsdk:"alert_on_failed_refreshes"`
    RefreshMonitorId types.String `tfsdk:"refresh_monitor_id"`
    FailedRefreshes types.String `tfsdk:"failed_refreshes"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *BiToolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_bi_tool"
}

func (d *BiToolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A BI tool whose dashboards, and the tables behind them, Data Observability reads, so a table's alerts list the dashboards built on it. Look up an existing bi tool by `id`, or by any of its other arguments (`name`, `alert_on_failed_refreshes`, `base_url`, ...): each one set must match, and exactly one bi tool may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "What the BI tool is called here.",
                Optional: true,
                Computed: true,
            },
            "tool": schema.StringAttribute{
                MarkdownDescription: "Which BI tool it is.",
                Optional: true,
                Computed: true,
            },
            "base_url": schema.StringAttribute{
                MarkdownDescription: "The tool's own address, such as https://metabase.acme.com or https://acme.cloud.looker.com. Power BI needs none.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether its dashboards are read every hour.",
                Optional: true,
                Computed: true,
            },
            "last_synced_at": schema.StringAttribute{
                MarkdownDescription: "When its dashboards were last read.",
                Computed: true,
            },
            "last_sync_error": schema.StringAttribute{
                MarkdownDescription: "Why the last read failed; empty when it did not.",
                Optional: true,
                Computed: true,
            },
            "dashboard_count": schema.NumberAttribute{
                MarkdownDescription: "How many dashboards the last read found.",
                Optional: true,
                Computed: true,
            },
            "matched_table_count": schema.NumberAttribute{
                MarkdownDescription: "How many watched tables the last read found behind its dashboards.",
                Optional: true,
                Computed: true,
            },
            "lineage_summary": schema.StringAttribute{
                MarkdownDescription: "Why the last read's dashboards are, or are not, tied to watched tables: the tables they read and how many matched, the data sources behind them and whether OneUptime watches each, and what keeps the rest from matching, such as a Power BI tenant setting that is off. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "last_sync_warning": schema.StringAttribute{
                MarkdownDescription: "What the last read could not read while it read the dashboards - views or refreshes - and why; empty when it read everything.",
                Optional: true,
                Computed: true,
            },
            "usage_read_from": schema.StringAttribute{
                MarkdownDescription: "The start of the hours of the tool's activity log read so far (Power BI): its dashboards' views are counted from them.",
                Computed: true,
            },
            "usage_read_until": schema.StringAttribute{
                MarkdownDescription: "The end of the hours of the tool's activity log read so far (Power BI): the next read goes on from here.",
                Computed: true,
            },
            "alert_on_failed_refreshes": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime keeps a monitor that opens an alert for each dataset whose refresh failed (Power BI), resolved once a later refresh completes.",
                Optional: true,
                Computed: true,
            },
            "refresh_monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Metrics monitor that alerts on the tool's failed refreshes.",
                Optional: true,
                Computed: true,
            },
            "failed_refreshes": schema.StringAttribute{
                MarkdownDescription: "The datasets whose last refresh failed or was cancelled (Power BI), each with its error, until a later refresh completes - and for an hour after that, marked recovered. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *BiToolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BiToolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data BiToolDataSourceModel

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
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Tool.IsNull() && !data.Tool.IsUnknown() {
        filters["tool"] = data.Tool.ValueString()
        filterNames = append(filterNames, "tool = "+fmt.Sprintf("%q", data.Tool.ValueString()))
    }
    if !data.BaseUrl.IsNull() && !data.BaseUrl.IsUnknown() {
        filters["baseUrl"] = data.BaseUrl.ValueString()
        filterNames = append(filterNames, "base_url = "+fmt.Sprintf("%q", data.BaseUrl.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.LastSyncError.IsNull() && !data.LastSyncError.IsUnknown() {
        filters["lastSyncError"] = data.LastSyncError.ValueString()
        filterNames = append(filterNames, "last_sync_error = "+fmt.Sprintf("%q", data.LastSyncError.ValueString()))
    }
    if !data.DashboardCount.IsNull() && !data.DashboardCount.IsUnknown() {
        filters["dashboardCount"] = LookupNumber(data.DashboardCount)
        filterNames = append(filterNames, "dashboard_count = "+data.DashboardCount.ValueBigFloat().String())
    }
    if !data.MatchedTableCount.IsNull() && !data.MatchedTableCount.IsUnknown() {
        filters["matchedTableCount"] = LookupNumber(data.MatchedTableCount)
        filterNames = append(filterNames, "matched_table_count = "+data.MatchedTableCount.ValueBigFloat().String())
    }
    if !data.LastSyncWarning.IsNull() && !data.LastSyncWarning.IsUnknown() {
        filters["lastSyncWarning"] = data.LastSyncWarning.ValueString()
        filterNames = append(filterNames, "last_sync_warning = "+fmt.Sprintf("%q", data.LastSyncWarning.ValueString()))
    }
    if !data.AlertOnFailedRefreshes.IsNull() && !data.AlertOnFailedRefreshes.IsUnknown() {
        filters["alertOnFailedRefreshes"] = data.AlertOnFailedRefreshes.ValueBool()
        filterNames = append(filterNames, "alert_on_failed_refreshes = "+fmt.Sprintf("%t", data.AlertOnFailedRefreshes.ValueBool()))
    }
    if !data.RefreshMonitorId.IsNull() && !data.RefreshMonitorId.IsUnknown() {
        filters["refreshMonitorId"] = data.RefreshMonitorId.ValueString()
        filterNames = append(filterNames, "refresh_monitor_id = "+fmt.Sprintf("%q", data.RefreshMonitorId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the bi tool up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the bi tool up by.",
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
        "name": true,
        "tool": true,
        "baseUrl": true,
        "isEnabled": true,
        "lastSyncedAt": true,
        "lastSyncError": true,
        "dashboardCount": true,
        "matchedTableCount": true,
        "lineageSummary": true,
        "lastSyncWarning": true,
        "usageReadFrom": true,
        "usageReadUntil": true,
        "alertOnFailedRefreshes": true,
        "refreshMonitorId": true,
        "failedRefreshes": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-bi-connection/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read bi_tool, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No bi tool found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read bi_tool: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-bi-connection/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list bi_tool, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list bi_tool: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No bi tool matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one bi tool matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for bi_tool.")
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
    if obj, ok := item["tool"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Tool = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Tool = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Tool = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Tool = types.StringValue(string(jsonBytes))
        } else {
            data.Tool = types.StringNull()
        }
    } else if arr, ok := item["tool"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Tool = types.StringValue(string(jsonBytes))
        } else {
            data.Tool = types.StringNull()
        }
    } else if val, ok := item["tool"].(string); ok {
        data.Tool = types.StringValue(val)
    } else {
        data.Tool = types.StringNull()
    }
    if obj, ok := item["baseUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BaseUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BaseUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BaseUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BaseUrl = types.StringValue(string(jsonBytes))
        } else {
            data.BaseUrl = types.StringNull()
        }
    } else if arr, ok := item["baseUrl"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.BaseUrl = types.StringValue(string(jsonBytes))
        } else {
            data.BaseUrl = types.StringNull()
        }
    } else if val, ok := item["baseUrl"].(string); ok {
        data.BaseUrl = types.StringValue(val)
    } else {
        data.BaseUrl = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["lastSyncedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSyncedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSyncedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSyncedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSyncedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncedAt = types.StringNull()
        }
    } else if arr, ok := item["lastSyncedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSyncedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncedAt = types.StringNull()
        }
    } else if val, ok := item["lastSyncedAt"].(string); ok {
        data.LastSyncedAt = types.StringValue(val)
    } else {
        data.LastSyncedAt = types.StringNull()
    }
    if obj, ok := item["lastSyncError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSyncError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSyncError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSyncError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSyncError = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncError = types.StringNull()
        }
    } else if arr, ok := item["lastSyncError"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSyncError = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncError = types.StringNull()
        }
    } else if val, ok := item["lastSyncError"].(string); ok {
        data.LastSyncError = types.StringValue(val)
    } else {
        data.LastSyncError = types.StringNull()
    }
    data.DashboardCount = NumberFromAPI(data.DashboardCount, item["dashboardCount"])
    data.MatchedTableCount = NumberFromAPI(data.MatchedTableCount, item["matchedTableCount"])
    if obj, ok := item["lineageSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LineageSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LineageSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LineageSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LineageSummary = types.StringValue(string(jsonBytes))
        } else {
            data.LineageSummary = types.StringNull()
        }
    } else if arr, ok := item["lineageSummary"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LineageSummary = types.StringValue(string(jsonBytes))
        } else {
            data.LineageSummary = types.StringNull()
        }
    } else if val, ok := item["lineageSummary"].(string); ok {
        data.LineageSummary = types.StringValue(val)
    } else {
        data.LineageSummary = types.StringNull()
    }
    if obj, ok := item["lastSyncWarning"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSyncWarning = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSyncWarning = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSyncWarning = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSyncWarning = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncWarning = types.StringNull()
        }
    } else if arr, ok := item["lastSyncWarning"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSyncWarning = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncWarning = types.StringNull()
        }
    } else if val, ok := item["lastSyncWarning"].(string); ok {
        data.LastSyncWarning = types.StringValue(val)
    } else {
        data.LastSyncWarning = types.StringNull()
    }
    if obj, ok := item["usageReadFrom"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UsageReadFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UsageReadFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UsageReadFrom = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UsageReadFrom = types.StringValue(string(jsonBytes))
        } else {
            data.UsageReadFrom = types.StringNull()
        }
    } else if arr, ok := item["usageReadFrom"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UsageReadFrom = types.StringValue(string(jsonBytes))
        } else {
            data.UsageReadFrom = types.StringNull()
        }
    } else if val, ok := item["usageReadFrom"].(string); ok {
        data.UsageReadFrom = types.StringValue(val)
    } else {
        data.UsageReadFrom = types.StringNull()
    }
    if obj, ok := item["usageReadUntil"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UsageReadUntil = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UsageReadUntil = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UsageReadUntil = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UsageReadUntil = types.StringValue(string(jsonBytes))
        } else {
            data.UsageReadUntil = types.StringNull()
        }
    } else if arr, ok := item["usageReadUntil"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UsageReadUntil = types.StringValue(string(jsonBytes))
        } else {
            data.UsageReadUntil = types.StringNull()
        }
    } else if val, ok := item["usageReadUntil"].(string); ok {
        data.UsageReadUntil = types.StringValue(val)
    } else {
        data.UsageReadUntil = types.StringNull()
    }
    if val, ok := item["alertOnFailedRefreshes"].(bool); ok {
        data.AlertOnFailedRefreshes = types.BoolValue(val)
    } else {
        data.AlertOnFailedRefreshes = types.BoolNull()
    }
    if obj, ok := item["refreshMonitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RefreshMonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RefreshMonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RefreshMonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RefreshMonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.RefreshMonitorId = types.StringNull()
        }
    } else if arr, ok := item["refreshMonitorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RefreshMonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.RefreshMonitorId = types.StringNull()
        }
    } else if val, ok := item["refreshMonitorId"].(string); ok {
        data.RefreshMonitorId = types.StringValue(val)
    } else {
        data.RefreshMonitorId = types.StringNull()
    }
    if obj, ok := item["failedRefreshes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FailedRefreshes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FailedRefreshes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FailedRefreshes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FailedRefreshes = types.StringValue(string(jsonBytes))
        } else {
            data.FailedRefreshes = types.StringNull()
        }
    } else if arr, ok := item["failedRefreshes"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FailedRefreshes = types.StringValue(string(jsonBytes))
        } else {
            data.FailedRefreshes = types.StringNull()
        }
    } else if val, ok := item["failedRefreshes"].(string); ok {
        data.FailedRefreshes = types.StringValue(val)
    } else {
        data.FailedRefreshes = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
