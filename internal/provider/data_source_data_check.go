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
var _ datasource.DataSource = &DataCheckDataSource{}

func NewDataCheckDataSource() datasource.DataSource {
    return &DataCheckDataSource{}
}

// DataCheckDataSource defines the data source implementation.
type DataCheckDataSource struct {
    client *Client
}

// DataCheckDataSourceModel describes the data source data model.
type DataCheckDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataAssetId types.String `tfsdk:"data_asset_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    CheckType types.String `tfsdk:"check_type"`
    ColumnName types.String `tfsdk:"column_name"`
    SegmentColumn types.String `tfsdk:"segment_column"`
    AcceptedValues types.String `tfsdk:"accepted_values"`
    Pattern types.String `tfsdk:"pattern"`
    FormatPreset types.String `tfsdk:"format_preset"`
    SamplePercent types.Number `tfsdk:"sample_percent"`
    ReferenceTableKey types.String `tfsdk:"reference_table_key"`
    ReferenceColumnName types.String `tfsdk:"reference_column_name"`
    CustomSql types.String `tfsdk:"custom_sql"`
    IntervalInMinutes types.Number `tfsdk:"interval_in_minutes"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    IsManagedByFile types.Bool `tfsdk:"is_managed_by_file"`
    AlertCondition types.String `tfsdk:"alert_condition"`
    AlertThreshold types.Number `tfsdk:"alert_threshold"`
    MonitorId types.String `tfsdk:"monitor_id"`
    DataReconciliationId types.String `tfsdk:"data_reconciliation_id"`
    LastRunAt types.String `tfsdk:"last_run_at"`
    LastRunStatus types.String `tfsdk:"last_run_status"`
    LastRunMessage types.String `tfsdk:"last_run_message"`
    LastValue types.Number `tfsdk:"last_value"`
    FailingRowsSample types.String `tfsdk:"failing_rows_sample"`
    LastSegmentCount types.Number `tfsdk:"last_segment_count"`
    ExternalId types.String `tfsdk:"external_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataCheckDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_check"
}

func (d *DataCheckDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A check on the values of a table - empty or repeated values, values outside a list or pattern, a query of your own - run on a schedule and written as a metric. Look up an existing data check by `id`, or by any of its other arguments (`name`, `alert_condition`, `alert_threshold`, ...): each one set must match, and exactly one data check may match them all.",

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
            "data_asset_id": schema.StringAttribute{
                MarkdownDescription: "ID of the table the check reads. The ID of a `oneuptime_data_asset`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_asset"),
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "ID of the data store the table is in, set from the table.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the check is called, as its metric and alerts name it.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Why the check matters, and what to do when it fails.",
                Optional: true,
                Computed: true,
            },
            "check_type": schema.StringAttribute{
                MarkdownDescription: "What the check measures.",
                Optional: true,
                Computed: true,
            },
            "column_name": schema.StringAttribute{
                MarkdownDescription: "The column the check reads, for a check on one column.",
                Optional: true,
                Computed: true,
            },
            "segment_column": schema.StringAttribute{
                MarkdownDescription: "Run the check once per value of this column - a tenant, a region - one series each, the largest 100.",
                Optional: true,
                Computed: true,
            },
            "accepted_values": schema.StringAttribute{
                MarkdownDescription: "For Invalid Percent: the values the column may hold. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "pattern": schema.StringAttribute{
                MarkdownDescription: "For Format Mismatch Percent: the regular expression every value must match.",
                Optional: true,
                Computed: true,
            },
            "format_preset": schema.StringAttribute{
                MarkdownDescription: "For Format Mismatch Percent: a format everyone wants, instead of writing the pattern.",
                Optional: true,
                Computed: true,
            },
            "sample_percent": schema.NumberAttribute{
                MarkdownDescription: "Read this percent of the table's rows instead of all of them, for a check that reads a rate.",
                Optional: true,
                Computed: true,
            },
            "reference_table_key": schema.StringAttribute{
                MarkdownDescription: "For Orphan Count: the table the column refers to, as schema.table.",
                Optional: true,
                Computed: true,
            },
            "reference_column_name": schema.StringAttribute{
                MarkdownDescription: "For Orphan Count: the column of that table it refers to.",
                Optional: true,
                Computed: true,
            },
            "custom_sql": schema.StringAttribute{
                MarkdownDescription: "For Custom SQL: one read-only SELECT returning a number - or, segmented, a segment and a number per row.",
                Optional: true,
                Computed: true,
            },
            "interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often the check runs: 15, 60, 360 or 1440 minutes.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the check runs.",
                Optional: true,
                Computed: true,
            },
            "is_managed_by_file": schema.BoolAttribute{
                MarkdownDescription: "Whether the check came from a file in a repository (oneuptime data-checks apply).",
                Optional: true,
                Computed: true,
            },
            "alert_condition": schema.StringAttribute{
                MarkdownDescription: "When the check's value should alert: unusually high or low for the hour of the week (learned), or above or below a threshold.",
                Optional: true,
                Computed: true,
            },
            "alert_threshold": schema.NumberAttribute{
                MarkdownDescription: "The threshold an Above or Below alert compares with.",
                Optional: true,
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Metrics monitor that alerts on the check. The ID of a `oneuptime_monitor`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_monitor"),
            },
            "data_reconciliation_id": schema.StringAttribute{
                MarkdownDescription: "For a check that reads one side of a reconciliation: the reconciliation's ID. Such a check is changed only from its reconciliation.",
                Optional: true,
                Computed: true,
            },
            "last_run_at": schema.StringAttribute{
                MarkdownDescription: "When the check last ran.",
                Computed: true,
            },
            "last_run_status": schema.StringAttribute{
                MarkdownDescription: "How the last run went: Ran, or Failed.",
                Optional: true,
                Computed: true,
            },
            "last_run_message": schema.StringAttribute{
                MarkdownDescription: "Why the last run failed, when it did.",
                Optional: true,
                Computed: true,
            },
            "last_value": schema.NumberAttribute{
                MarkdownDescription: "The value the last run read - the largest segment's, for a segmented check.",
                Optional: true,
                Computed: true,
            },
            "failing_rows_sample": schema.StringAttribute{
                MarkdownDescription: "A masked sample of the rows the last run found failing, when the project keeps them. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "last_segment_count": schema.NumberAttribute{
                MarkdownDescription: "How many segments the last run read.",
                Optional: true,
                Computed: true,
            },
            "external_id": schema.StringAttribute{
                MarkdownDescription: "For a dbt test, the test's unique id in the dbt project, which each upload of its artifacts is matched by.",
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

func (d *DataCheckDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataCheckDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DataAssetId.IsNull() && !data.DataAssetId.IsUnknown() {
        filters["dataAssetId"] = data.DataAssetId.ValueString()
        filterNames = append(filterNames, "data_asset_id = "+fmt.Sprintf("%q", data.DataAssetId.ValueString()))
    }
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CheckType.IsNull() && !data.CheckType.IsUnknown() {
        filters["checkType"] = data.CheckType.ValueString()
        filterNames = append(filterNames, "check_type = "+fmt.Sprintf("%q", data.CheckType.ValueString()))
    }
    if !data.ColumnName.IsNull() && !data.ColumnName.IsUnknown() {
        filters["columnName"] = data.ColumnName.ValueString()
        filterNames = append(filterNames, "column_name = "+fmt.Sprintf("%q", data.ColumnName.ValueString()))
    }
    if !data.SegmentColumn.IsNull() && !data.SegmentColumn.IsUnknown() {
        filters["segmentColumn"] = data.SegmentColumn.ValueString()
        filterNames = append(filterNames, "segment_column = "+fmt.Sprintf("%q", data.SegmentColumn.ValueString()))
    }
    if !data.Pattern.IsNull() && !data.Pattern.IsUnknown() {
        filters["pattern"] = data.Pattern.ValueString()
        filterNames = append(filterNames, "pattern = "+fmt.Sprintf("%q", data.Pattern.ValueString()))
    }
    if !data.FormatPreset.IsNull() && !data.FormatPreset.IsUnknown() {
        filters["formatPreset"] = data.FormatPreset.ValueString()
        filterNames = append(filterNames, "format_preset = "+fmt.Sprintf("%q", data.FormatPreset.ValueString()))
    }
    if !data.SamplePercent.IsNull() && !data.SamplePercent.IsUnknown() {
        filters["samplePercent"] = lookupNumber(data.SamplePercent)
        filterNames = append(filterNames, "sample_percent = "+data.SamplePercent.ValueBigFloat().String())
    }
    if !data.ReferenceTableKey.IsNull() && !data.ReferenceTableKey.IsUnknown() {
        filters["referenceTableKey"] = data.ReferenceTableKey.ValueString()
        filterNames = append(filterNames, "reference_table_key = "+fmt.Sprintf("%q", data.ReferenceTableKey.ValueString()))
    }
    if !data.ReferenceColumnName.IsNull() && !data.ReferenceColumnName.IsUnknown() {
        filters["referenceColumnName"] = data.ReferenceColumnName.ValueString()
        filterNames = append(filterNames, "reference_column_name = "+fmt.Sprintf("%q", data.ReferenceColumnName.ValueString()))
    }
    if !data.CustomSql.IsNull() && !data.CustomSql.IsUnknown() {
        filters["customSql"] = data.CustomSql.ValueString()
        filterNames = append(filterNames, "custom_sql = "+fmt.Sprintf("%q", data.CustomSql.ValueString()))
    }
    if !data.IntervalInMinutes.IsNull() && !data.IntervalInMinutes.IsUnknown() {
        filters["intervalInMinutes"] = lookupNumber(data.IntervalInMinutes)
        filterNames = append(filterNames, "interval_in_minutes = "+data.IntervalInMinutes.ValueBigFloat().String())
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.IsManagedByFile.IsNull() && !data.IsManagedByFile.IsUnknown() {
        filters["isManagedByFile"] = data.IsManagedByFile.ValueBool()
        filterNames = append(filterNames, "is_managed_by_file = "+fmt.Sprintf("%t", data.IsManagedByFile.ValueBool()))
    }
    if !data.AlertCondition.IsNull() && !data.AlertCondition.IsUnknown() {
        filters["alertCondition"] = data.AlertCondition.ValueString()
        filterNames = append(filterNames, "alert_condition = "+fmt.Sprintf("%q", data.AlertCondition.ValueString()))
    }
    if !data.AlertThreshold.IsNull() && !data.AlertThreshold.IsUnknown() {
        filters["alertThreshold"] = lookupNumber(data.AlertThreshold)
        filterNames = append(filterNames, "alert_threshold = "+data.AlertThreshold.ValueBigFloat().String())
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.DataReconciliationId.IsNull() && !data.DataReconciliationId.IsUnknown() {
        filters["dataReconciliationId"] = data.DataReconciliationId.ValueString()
        filterNames = append(filterNames, "data_reconciliation_id = "+fmt.Sprintf("%q", data.DataReconciliationId.ValueString()))
    }
    if !data.LastRunStatus.IsNull() && !data.LastRunStatus.IsUnknown() {
        filters["lastRunStatus"] = data.LastRunStatus.ValueString()
        filterNames = append(filterNames, "last_run_status = "+fmt.Sprintf("%q", data.LastRunStatus.ValueString()))
    }
    if !data.LastRunMessage.IsNull() && !data.LastRunMessage.IsUnknown() {
        filters["lastRunMessage"] = data.LastRunMessage.ValueString()
        filterNames = append(filterNames, "last_run_message = "+fmt.Sprintf("%q", data.LastRunMessage.ValueString()))
    }
    if !data.LastValue.IsNull() && !data.LastValue.IsUnknown() {
        filters["lastValue"] = lookupNumber(data.LastValue)
        filterNames = append(filterNames, "last_value = "+data.LastValue.ValueBigFloat().String())
    }
    if !data.LastSegmentCount.IsNull() && !data.LastSegmentCount.IsUnknown() {
        filters["lastSegmentCount"] = lookupNumber(data.LastSegmentCount)
        filterNames = append(filterNames, "last_segment_count = "+data.LastSegmentCount.ValueBigFloat().String())
    }
    if !data.ExternalId.IsNull() && !data.ExternalId.IsUnknown() {
        filters["externalId"] = data.ExternalId.ValueString()
        filterNames = append(filterNames, "external_id = "+fmt.Sprintf("%q", data.ExternalId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data check up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data check up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataAssetId": true,
        "dataStoreId": true,
        "name": true,
        "description": true,
        "checkType": true,
        "columnName": true,
        "segmentColumn": true,
        "acceptedValues": true,
        "pattern": true,
        "formatPreset": true,
        "samplePercent": true,
        "referenceTableKey": true,
        "referenceColumnName": true,
        "customSql": true,
        "intervalInMinutes": true,
        "isEnabled": true,
        "isManagedByFile": true,
        "alertCondition": true,
        "alertThreshold": true,
        "monitorId": true,
        "dataReconciliationId": true,
        "lastRunAt": true,
        "lastRunStatus": true,
        "lastRunMessage": true,
        "lastValue": true,
        "failingRowsSample": true,
        "lastSegmentCount": true,
        "externalId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-check/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_check, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data check found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_check: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-check/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_check, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_check: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data check matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data check matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_check.")
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
    if obj, ok := item["dataAssetId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataAssetId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.DataAssetId = types.StringNull()
        }
    } else if val, ok := item["dataAssetId"].(string); ok {
        data.DataAssetId = types.StringValue(val)
    } else {
        data.DataAssetId = types.StringNull()
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
    if obj, ok := item["checkType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CheckType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CheckType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CheckType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CheckType = types.StringValue(string(jsonBytes))
        } else {
            data.CheckType = types.StringNull()
        }
    } else if val, ok := item["checkType"].(string); ok {
        data.CheckType = types.StringValue(val)
    } else {
        data.CheckType = types.StringNull()
    }
    if obj, ok := item["columnName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ColumnName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ColumnName = types.StringValue(string(jsonBytes))
        } else {
            data.ColumnName = types.StringNull()
        }
    } else if val, ok := item["columnName"].(string); ok {
        data.ColumnName = types.StringValue(val)
    } else {
        data.ColumnName = types.StringNull()
    }
    if obj, ok := item["segmentColumn"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SegmentColumn = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SegmentColumn = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SegmentColumn = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SegmentColumn = types.StringValue(string(jsonBytes))
        } else {
            data.SegmentColumn = types.StringNull()
        }
    } else if val, ok := item["segmentColumn"].(string); ok {
        data.SegmentColumn = types.StringValue(val)
    } else {
        data.SegmentColumn = types.StringNull()
    }
    if obj, ok := item["acceptedValues"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcceptedValues = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcceptedValues = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcceptedValues = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcceptedValues = types.StringValue(string(jsonBytes))
        } else {
            data.AcceptedValues = types.StringNull()
        }
    } else if val, ok := item["acceptedValues"].(string); ok {
        data.AcceptedValues = types.StringValue(val)
    } else {
        data.AcceptedValues = types.StringNull()
    }
    if obj, ok := item["pattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Pattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Pattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Pattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Pattern = types.StringValue(string(jsonBytes))
        } else {
            data.Pattern = types.StringNull()
        }
    } else if val, ok := item["pattern"].(string); ok {
        data.Pattern = types.StringValue(val)
    } else {
        data.Pattern = types.StringNull()
    }
    if obj, ok := item["formatPreset"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FormatPreset = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FormatPreset = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FormatPreset = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FormatPreset = types.StringValue(string(jsonBytes))
        } else {
            data.FormatPreset = types.StringNull()
        }
    } else if val, ok := item["formatPreset"].(string); ok {
        data.FormatPreset = types.StringValue(val)
    } else {
        data.FormatPreset = types.StringNull()
    }
    data.SamplePercent = numberFromAPI(data.SamplePercent, item["samplePercent"])
    if obj, ok := item["referenceTableKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReferenceTableKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReferenceTableKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReferenceTableKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReferenceTableKey = types.StringValue(string(jsonBytes))
        } else {
            data.ReferenceTableKey = types.StringNull()
        }
    } else if val, ok := item["referenceTableKey"].(string); ok {
        data.ReferenceTableKey = types.StringValue(val)
    } else {
        data.ReferenceTableKey = types.StringNull()
    }
    if obj, ok := item["referenceColumnName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReferenceColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReferenceColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReferenceColumnName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReferenceColumnName = types.StringValue(string(jsonBytes))
        } else {
            data.ReferenceColumnName = types.StringNull()
        }
    } else if val, ok := item["referenceColumnName"].(string); ok {
        data.ReferenceColumnName = types.StringValue(val)
    } else {
        data.ReferenceColumnName = types.StringNull()
    }
    if obj, ok := item["customSql"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomSql = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomSql = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomSql = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomSql = types.StringValue(string(jsonBytes))
        } else {
            data.CustomSql = types.StringNull()
        }
    } else if val, ok := item["customSql"].(string); ok {
        data.CustomSql = types.StringValue(val)
    } else {
        data.CustomSql = types.StringNull()
    }
    data.IntervalInMinutes = numberFromAPI(data.IntervalInMinutes, item["intervalInMinutes"])
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["isManagedByFile"].(bool); ok {
        data.IsManagedByFile = types.BoolValue(val)
    } else {
        data.IsManagedByFile = types.BoolNull()
    }
    if obj, ok := item["alertCondition"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertCondition = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertCondition = types.StringValue(string(jsonBytes))
        } else {
            data.AlertCondition = types.StringNull()
        }
    } else if val, ok := item["alertCondition"].(string); ok {
        data.AlertCondition = types.StringValue(val)
    } else {
        data.AlertCondition = types.StringNull()
    }
    data.AlertThreshold = numberFromAPI(data.AlertThreshold, item["alertThreshold"])
    if obj, ok := item["monitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := item["dataReconciliationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataReconciliationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataReconciliationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataReconciliationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataReconciliationId = types.StringValue(string(jsonBytes))
        } else {
            data.DataReconciliationId = types.StringNull()
        }
    } else if val, ok := item["dataReconciliationId"].(string); ok {
        data.DataReconciliationId = types.StringValue(val)
    } else {
        data.DataReconciliationId = types.StringNull()
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
    if obj, ok := item["lastRunStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunStatus = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunStatus = types.StringNull()
        }
    } else if val, ok := item["lastRunStatus"].(string); ok {
        data.LastRunStatus = types.StringValue(val)
    } else {
        data.LastRunStatus = types.StringNull()
    }
    if obj, ok := item["lastRunMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunMessage = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunMessage = types.StringNull()
        }
    } else if val, ok := item["lastRunMessage"].(string); ok {
        data.LastRunMessage = types.StringValue(val)
    } else {
        data.LastRunMessage = types.StringNull()
    }
    data.LastValue = numberFromAPI(data.LastValue, item["lastValue"])
    if obj, ok := item["failingRowsSample"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FailingRowsSample = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FailingRowsSample = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FailingRowsSample = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FailingRowsSample = types.StringValue(string(jsonBytes))
        } else {
            data.FailingRowsSample = types.StringNull()
        }
    } else if val, ok := item["failingRowsSample"].(string); ok {
        data.FailingRowsSample = types.StringValue(val)
    } else {
        data.FailingRowsSample = types.StringNull()
    }
    data.LastSegmentCount = numberFromAPI(data.LastSegmentCount, item["lastSegmentCount"])
    if obj, ok := item["externalId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExternalId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExternalId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExternalId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExternalId = types.StringValue(string(jsonBytes))
        } else {
            data.ExternalId = types.StringNull()
        }
    } else if val, ok := item["externalId"].(string); ok {
        data.ExternalId = types.StringValue(val)
    } else {
        data.ExternalId = types.StringNull()
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
