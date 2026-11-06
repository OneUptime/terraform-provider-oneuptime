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
var _ datasource.DataSource = &WarehouseJobDataSource{}

func NewWarehouseJobDataSource() datasource.DataSource {
    return &WarehouseJobDataSource{}
}

// WarehouseJobDataSource defines the data source implementation.
type WarehouseJobDataSource struct {
    client *Client
}

// WarehouseJobDataSourceModel describes the data source data model.
type WarehouseJobDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    JobKey types.String `tfsdk:"job_key"`
    Name types.String `tfsdk:"name"`
    UserName types.String `tfsdk:"user_name"`
    WarehouseName types.String `tfsdk:"warehouse_name"`
    Team types.String `tfsdk:"team"`
    RunCount types.Number `tfsdk:"run_count"`
    FailedRunCount types.Number `tfsdk:"failed_run_count"`
    Credits types.Number `tfsdk:"credits"`
    BytesBilled types.Number `tfsdk:"bytes_billed"`
    WrittenTables types.String `tfsdk:"written_tables"`
    WritesOnlyUnusedTables types.Bool `tfsdk:"writes_only_unused_tables"`
    LastRunAt types.String `tfsdk:"last_run_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *WarehouseJobDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_warehouse_job"
}

func (d *WarehouseJobDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A query a warehouse runs again and again, with what it spent over the last day, its failures and the watched tables it writes. Look up an existing warehouse job by `id`, or by any of its other arguments (`name`, `bytes_billed`, `created_by_user_id`, ...): each one set must match, and exactly one warehouse job may match them all.",

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
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "ID of the warehouse the job runs in. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_store"),
            },
            "job_key": schema.StringAttribute{
                MarkdownDescription: "The warehouse's own key for the query, its values left out.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The job's tag, or what it does and the table it writes: never its query.",
                Optional: true,
                Computed: true,
            },
            "user_name": schema.StringAttribute{
                MarkdownDescription: "The login or service account that runs it.",
                Optional: true,
                Computed: true,
            },
            "warehouse_name": schema.StringAttribute{
                MarkdownDescription: "The warehouse or reservation it runs on.",
                Optional: true,
                Computed: true,
            },
            "team": schema.StringAttribute{
                MarkdownDescription: "The team its query tag or label names.",
                Optional: true,
                Computed: true,
            },
            "run_count": schema.NumberAttribute{
                MarkdownDescription: "How many times it ran in the last day.",
                Optional: true,
                Computed: true,
            },
            "failed_run_count": schema.NumberAttribute{
                MarkdownDescription: "How many of its runs in the last day failed.",
                Optional: true,
                Computed: true,
            },
            "credits": schema.NumberAttribute{
                MarkdownDescription: "Snowflake compute credits it spent in the last day.",
                Optional: true,
                Computed: true,
            },
            "bytes_billed": schema.NumberAttribute{
                MarkdownDescription: "BigQuery bytes billed for it in the last day.",
                Optional: true,
                Computed: true,
            },
            "written_tables": schema.StringAttribute{
                MarkdownDescription: "The watched tables it writes, each as { dataAssetId, table }. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "writes_only_unused_tables": schema.BoolAttribute{
                MarkdownDescription: "Whether every table it writes has gone unread for 90 days: spend nobody uses.",
                Optional: true,
                Computed: true,
            },
            "last_run_at": schema.StringAttribute{
                MarkdownDescription: "When it last ran.",
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

func (d *WarehouseJobDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WarehouseJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WarehouseJobDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.JobKey.IsNull() && !data.JobKey.IsUnknown() {
        filters["jobKey"] = data.JobKey.ValueString()
        filterNames = append(filterNames, "job_key = "+fmt.Sprintf("%q", data.JobKey.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.UserName.IsNull() && !data.UserName.IsUnknown() {
        filters["userName"] = data.UserName.ValueString()
        filterNames = append(filterNames, "user_name = "+fmt.Sprintf("%q", data.UserName.ValueString()))
    }
    if !data.WarehouseName.IsNull() && !data.WarehouseName.IsUnknown() {
        filters["warehouseName"] = data.WarehouseName.ValueString()
        filterNames = append(filterNames, "warehouse_name = "+fmt.Sprintf("%q", data.WarehouseName.ValueString()))
    }
    if !data.Team.IsNull() && !data.Team.IsUnknown() {
        filters["team"] = data.Team.ValueString()
        filterNames = append(filterNames, "team = "+fmt.Sprintf("%q", data.Team.ValueString()))
    }
    if !data.RunCount.IsNull() && !data.RunCount.IsUnknown() {
        filters["runCount"] = lookupNumber(data.RunCount)
        filterNames = append(filterNames, "run_count = "+data.RunCount.ValueBigFloat().String())
    }
    if !data.FailedRunCount.IsNull() && !data.FailedRunCount.IsUnknown() {
        filters["failedRunCount"] = lookupNumber(data.FailedRunCount)
        filterNames = append(filterNames, "failed_run_count = "+data.FailedRunCount.ValueBigFloat().String())
    }
    if !data.Credits.IsNull() && !data.Credits.IsUnknown() {
        filters["credits"] = lookupNumber(data.Credits)
        filterNames = append(filterNames, "credits = "+data.Credits.ValueBigFloat().String())
    }
    if !data.BytesBilled.IsNull() && !data.BytesBilled.IsUnknown() {
        filters["bytesBilled"] = lookupNumber(data.BytesBilled)
        filterNames = append(filterNames, "bytes_billed = "+data.BytesBilled.ValueBigFloat().String())
    }
    if !data.WritesOnlyUnusedTables.IsNull() && !data.WritesOnlyUnusedTables.IsUnknown() {
        filters["writesOnlyUnusedTables"] = data.WritesOnlyUnusedTables.ValueBool()
        filterNames = append(filterNames, "writes_only_unused_tables = "+fmt.Sprintf("%t", data.WritesOnlyUnusedTables.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the warehouse job up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the warehouse job up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataStoreId": true,
        "jobKey": true,
        "name": true,
        "userName": true,
        "warehouseName": true,
        "team": true,
        "runCount": true,
        "failedRunCount": true,
        "credits": true,
        "bytesBilled": true,
        "writtenTables": true,
        "writesOnlyUnusedTables": true,
        "lastRunAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-warehouse-job/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read warehouse_job, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No warehouse job found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read warehouse_job: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-warehouse-job/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list warehouse_job, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list warehouse_job: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No warehouse job matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one warehouse job matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for warehouse_job.")
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
    if obj, ok := item["dataStoreId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataStoreId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreId = types.StringNull()
        }
    } else if val, ok := item["dataStoreId"].(string); ok {
        data.DataStoreId = types.StringValue(val)
    } else {
        data.DataStoreId = types.StringNull()
    }
    if obj, ok := item["jobKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.JobKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.JobKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.JobKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.JobKey = types.StringValue(string(jsonBytes))
        } else {
            data.JobKey = types.StringNull()
        }
    } else if val, ok := item["jobKey"].(string); ok {
        data.JobKey = types.StringValue(val)
    } else {
        data.JobKey = types.StringNull()
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
    if obj, ok := item["userName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserName = types.StringValue(string(jsonBytes))
        } else {
            data.UserName = types.StringNull()
        }
    } else if val, ok := item["userName"].(string); ok {
        data.UserName = types.StringValue(val)
    } else {
        data.UserName = types.StringNull()
    }
    if obj, ok := item["warehouseName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WarehouseName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WarehouseName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WarehouseName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WarehouseName = types.StringValue(string(jsonBytes))
        } else {
            data.WarehouseName = types.StringNull()
        }
    } else if val, ok := item["warehouseName"].(string); ok {
        data.WarehouseName = types.StringValue(val)
    } else {
        data.WarehouseName = types.StringNull()
    }
    if obj, ok := item["team"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Team = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Team = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Team = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Team = types.StringValue(string(jsonBytes))
        } else {
            data.Team = types.StringNull()
        }
    } else if val, ok := item["team"].(string); ok {
        data.Team = types.StringValue(val)
    } else {
        data.Team = types.StringNull()
    }
    data.RunCount = numberFromAPI(data.RunCount, item["runCount"])
    data.FailedRunCount = numberFromAPI(data.FailedRunCount, item["failedRunCount"])
    data.Credits = numberFromAPI(data.Credits, item["credits"])
    data.BytesBilled = numberFromAPI(data.BytesBilled, item["bytesBilled"])
    if obj, ok := item["writtenTables"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WrittenTables = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WrittenTables = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WrittenTables = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WrittenTables = types.StringValue(string(jsonBytes))
        } else {
            data.WrittenTables = types.StringNull()
        }
    } else if val, ok := item["writtenTables"].(string); ok {
        data.WrittenTables = types.StringValue(val)
    } else {
        data.WrittenTables = types.StringNull()
    }
    if val, ok := item["writesOnlyUnusedTables"].(bool); ok {
        data.WritesOnlyUnusedTables = types.BoolValue(val)
    } else {
        data.WritesOnlyUnusedTables = types.BoolNull()
    }
    if obj, ok := item["lastRunAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunAt = types.StringNull()
        }
    } else if val, ok := item["lastRunAt"].(string); ok {
        data.LastRunAt = types.StringValue(val)
    } else {
        data.LastRunAt = types.StringNull()
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
