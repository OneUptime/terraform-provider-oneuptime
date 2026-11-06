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
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Tool types.String `tfsdk:"tool"`
    BaseUrl types.String `tfsdk:"base_url"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    LastSyncedAt types.String `tfsdk:"last_synced_at"`
    LastSyncError types.String `tfsdk:"last_sync_error"`
    DashboardCount types.Number `tfsdk:"dashboard_count"`
    MatchedTableCount types.Number `tfsdk:"matched_table_count"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *BiToolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_bi_tool"
}

func (d *BiToolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A BI tool whose dashboards, and the tables behind them, Data Observability reads, so a table's alerts list the dashboards built on it. Look up an existing bi tool by `id`, or by any of its other arguments (`name`, `base_url`, `created_by_user_id`, ...): each one set must match, and exactly one bi tool may match them all.",

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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
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
        filters["dashboardCount"] = lookupNumber(data.DashboardCount)
        filterNames = append(filterNames, "dashboard_count = "+data.DashboardCount.ValueBigFloat().String())
    }
    if !data.MatchedTableCount.IsNull() && !data.MatchedTableCount.IsUnknown() {
        filters["matchedTableCount"] = lookupNumber(data.MatchedTableCount)
        filterNames = append(filterNames, "matched_table_count = "+data.MatchedTableCount.ValueBigFloat().String())
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
        "projectId": true,
        "name": true,
        "tool": true,
        "baseUrl": true,
        "isEnabled": true,
        "lastSyncedAt": true,
        "lastSyncError": true,
        "dashboardCount": true,
        "matchedTableCount": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No bi tool matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one bi tool matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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
    } else if val, ok := item["lastSyncError"].(string); ok {
        data.LastSyncError = types.StringValue(val)
    } else {
        data.LastSyncError = types.StringNull()
    }
    data.DashboardCount = numberFromAPI(data.DashboardCount, item["dashboardCount"])
    data.MatchedTableCount = numberFromAPI(data.MatchedTableCount, item["matchedTableCount"])
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
