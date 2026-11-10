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
var _ datasource.DataSource = &DataReconciliationDataSource{}

func NewDataReconciliationDataSource() datasource.DataSource {
    return &DataReconciliationDataSource{}
}

// DataReconciliationDataSource defines the data source implementation.
type DataReconciliationDataSource struct {
    client *Client
}

// DataReconciliationDataSourceModel describes the data source data model.
type DataReconciliationDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IntervalInMinutes types.Number `tfsdk:"interval_in_minutes"`
    ToleranceType types.String `tfsdk:"tolerance_type"`
    Tolerance types.Number `tfsdk:"tolerance"`
    SourceDataAssetId types.String `tfsdk:"source_data_asset_id"`
    SourceMeasure types.String `tfsdk:"source_measure"`
    SourceColumnName types.String `tfsdk:"source_column_name"`
    SourceCustomSql types.String `tfsdk:"source_custom_sql"`
    SourceDataCheckId types.String `tfsdk:"source_data_check_id"`
    SourceValue types.Number `tfsdk:"source_value"`
    TargetDataAssetId types.String `tfsdk:"target_data_asset_id"`
    TargetMeasure types.String `tfsdk:"target_measure"`
    TargetColumnName types.String `tfsdk:"target_column_name"`
    TargetCustomSql types.String `tfsdk:"target_custom_sql"`
    TargetDataCheckId types.String `tfsdk:"target_data_check_id"`
    TargetValue types.Number `tfsdk:"target_value"`
    Difference types.Number `tfsdk:"difference"`
    DifferencePercent types.Number `tfsdk:"difference_percent"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    ComparedAt types.String `tfsdk:"compared_at"`
    StatusChangedAt types.String `tfsdk:"status_changed_at"`
    MonitorId types.String `tfsdk:"monitor_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataReconciliationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_reconciliation"
}

func (d *DataReconciliationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "One measure of two tables - rows, a sum, distinct values, a query's number - compared, in one database or two, with a tolerance, and an alert when they drift. Look up an existing data reconciliation by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one data reconciliation may match them all.",

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
                MarkdownDescription: "What the reconciliation is called, as its alert names it.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Why the two must agree, and what to do when they do not.",
                Optional: true,
                Computed: true,
            },
            "interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often both sides are read and compared: every 15 minutes, hour, 6 hours or day.",
                Optional: true,
                Computed: true,
            },
            "tolerance_type": schema.StringAttribute{
                MarkdownDescription: "How the tolerance is read: Absolute, an amount, or Percent, of the source's value.",
                Optional: true,
                Computed: true,
            },
            "tolerance": schema.NumberAttribute{
                MarkdownDescription: "How far apart the two may be before they have drifted: 0 means they must agree exactly.",
                Optional: true,
                Computed: true,
            },
            "source_data_asset_id": schema.StringAttribute{
                MarkdownDescription: "ID of the source table the reconciliation reads. The ID of a `oneuptime_data_asset`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_asset"),
            },
            "source_measure": schema.StringAttribute{
                MarkdownDescription: "What is measured of the source table: Row Count, Sum, Distinct Count, or Custom SQL.",
                Optional: true,
                Computed: true,
            },
            "source_column_name": schema.StringAttribute{
                MarkdownDescription: "For Sum or Distinct Count: the column of the source table.",
                Optional: true,
                Computed: true,
            },
            "source_custom_sql": schema.StringAttribute{
                MarkdownDescription: "For Custom SQL: a read-only query of the source's database returning one number, as value.",
                Optional: true,
                Computed: true,
            },
            "source_data_check_id": schema.StringAttribute{
                MarkdownDescription: "ID of the data check that reads the source, made and kept by the reconciliation. The ID of a `oneuptime_data_check`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_check"),
            },
            "source_value": schema.NumberAttribute{
                MarkdownDescription: "The source's value at the latest comparison.",
                Optional: true,
                Computed: true,
            },
            "target_data_asset_id": schema.StringAttribute{
                MarkdownDescription: "ID of the target table the reconciliation reads. The ID of a `oneuptime_data_asset`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_asset"),
            },
            "target_measure": schema.StringAttribute{
                MarkdownDescription: "What is measured of the target table: Row Count, Sum, Distinct Count, or Custom SQL.",
                Optional: true,
                Computed: true,
            },
            "target_column_name": schema.StringAttribute{
                MarkdownDescription: "For Sum or Distinct Count: the column of the target table.",
                Optional: true,
                Computed: true,
            },
            "target_custom_sql": schema.StringAttribute{
                MarkdownDescription: "For Custom SQL: a read-only query of the target's database returning one number, as value.",
                Optional: true,
                Computed: true,
            },
            "target_data_check_id": schema.StringAttribute{
                MarkdownDescription: "ID of the data check that reads the target, made and kept by the reconciliation. The ID of a `oneuptime_data_check`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_check"),
            },
            "target_value": schema.NumberAttribute{
                MarkdownDescription: "The target's value at the latest comparison.",
                Optional: true,
                Computed: true,
            },
            "difference": schema.NumberAttribute{
                MarkdownDescription: "The target's value less the source's, at the latest comparison.",
                Optional: true,
                Computed: true,
            },
            "difference_percent": schema.NumberAttribute{
                MarkdownDescription: "The difference as a percentage of the source's value, at the latest comparison.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "In Step, Drifted, Waiting for both sides to be read together, or Failed when a side's check could not run.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Why the reconciliation is waiting or failed.",
                Optional: true,
                Computed: true,
            },
            "compared_at": schema.StringAttribute{
                MarkdownDescription: "When the later of the two compared values was read.",
                Computed: true,
            },
            "status_changed_at": schema.StringAttribute{
                MarkdownDescription: "When the status last changed.",
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Metrics monitor that alerts when the two drift. The ID of a `oneuptime_monitor`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_monitor"),
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

func (d *DataReconciliationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataReconciliationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataReconciliationDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.IntervalInMinutes.IsNull() && !data.IntervalInMinutes.IsUnknown() {
        filters["intervalInMinutes"] = lookupNumber(data.IntervalInMinutes)
        filterNames = append(filterNames, "interval_in_minutes = "+data.IntervalInMinutes.ValueBigFloat().String())
    }
    if !data.ToleranceType.IsNull() && !data.ToleranceType.IsUnknown() {
        filters["toleranceType"] = data.ToleranceType.ValueString()
        filterNames = append(filterNames, "tolerance_type = "+fmt.Sprintf("%q", data.ToleranceType.ValueString()))
    }
    if !data.Tolerance.IsNull() && !data.Tolerance.IsUnknown() {
        filters["tolerance"] = lookupNumber(data.Tolerance)
        filterNames = append(filterNames, "tolerance = "+data.Tolerance.ValueBigFloat().String())
    }
    if !data.SourceDataAssetId.IsNull() && !data.SourceDataAssetId.IsUnknown() {
        filters["sourceDataAssetId"] = data.SourceDataAssetId.ValueString()
        filterNames = append(filterNames, "source_data_asset_id = "+fmt.Sprintf("%q", data.SourceDataAssetId.ValueString()))
    }
    if !data.SourceMeasure.IsNull() && !data.SourceMeasure.IsUnknown() {
        filters["sourceMeasure"] = data.SourceMeasure.ValueString()
        filterNames = append(filterNames, "source_measure = "+fmt.Sprintf("%q", data.SourceMeasure.ValueString()))
    }
    if !data.SourceColumnName.IsNull() && !data.SourceColumnName.IsUnknown() {
        filters["sourceColumnName"] = data.SourceColumnName.ValueString()
        filterNames = append(filterNames, "source_column_name = "+fmt.Sprintf("%q", data.SourceColumnName.ValueString()))
    }
    if !data.SourceCustomSql.IsNull() && !data.SourceCustomSql.IsUnknown() {
        filters["sourceCustomSql"] = data.SourceCustomSql.ValueString()
        filterNames = append(filterNames, "source_custom_sql = "+fmt.Sprintf("%q", data.SourceCustomSql.ValueString()))
    }
    if !data.SourceDataCheckId.IsNull() && !data.SourceDataCheckId.IsUnknown() {
        filters["sourceDataCheckId"] = data.SourceDataCheckId.ValueString()
        filterNames = append(filterNames, "source_data_check_id = "+fmt.Sprintf("%q", data.SourceDataCheckId.ValueString()))
    }
    if !data.SourceValue.IsNull() && !data.SourceValue.IsUnknown() {
        filters["sourceValue"] = lookupNumber(data.SourceValue)
        filterNames = append(filterNames, "source_value = "+data.SourceValue.ValueBigFloat().String())
    }
    if !data.TargetDataAssetId.IsNull() && !data.TargetDataAssetId.IsUnknown() {
        filters["targetDataAssetId"] = data.TargetDataAssetId.ValueString()
        filterNames = append(filterNames, "target_data_asset_id = "+fmt.Sprintf("%q", data.TargetDataAssetId.ValueString()))
    }
    if !data.TargetMeasure.IsNull() && !data.TargetMeasure.IsUnknown() {
        filters["targetMeasure"] = data.TargetMeasure.ValueString()
        filterNames = append(filterNames, "target_measure = "+fmt.Sprintf("%q", data.TargetMeasure.ValueString()))
    }
    if !data.TargetColumnName.IsNull() && !data.TargetColumnName.IsUnknown() {
        filters["targetColumnName"] = data.TargetColumnName.ValueString()
        filterNames = append(filterNames, "target_column_name = "+fmt.Sprintf("%q", data.TargetColumnName.ValueString()))
    }
    if !data.TargetCustomSql.IsNull() && !data.TargetCustomSql.IsUnknown() {
        filters["targetCustomSql"] = data.TargetCustomSql.ValueString()
        filterNames = append(filterNames, "target_custom_sql = "+fmt.Sprintf("%q", data.TargetCustomSql.ValueString()))
    }
    if !data.TargetDataCheckId.IsNull() && !data.TargetDataCheckId.IsUnknown() {
        filters["targetDataCheckId"] = data.TargetDataCheckId.ValueString()
        filterNames = append(filterNames, "target_data_check_id = "+fmt.Sprintf("%q", data.TargetDataCheckId.ValueString()))
    }
    if !data.TargetValue.IsNull() && !data.TargetValue.IsUnknown() {
        filters["targetValue"] = lookupNumber(data.TargetValue)
        filterNames = append(filterNames, "target_value = "+data.TargetValue.ValueBigFloat().String())
    }
    if !data.Difference.IsNull() && !data.Difference.IsUnknown() {
        filters["difference"] = lookupNumber(data.Difference)
        filterNames = append(filterNames, "difference = "+data.Difference.ValueBigFloat().String())
    }
    if !data.DifferencePercent.IsNull() && !data.DifferencePercent.IsUnknown() {
        filters["differencePercent"] = lookupNumber(data.DifferencePercent)
        filterNames = append(filterNames, "difference_percent = "+data.DifferencePercent.ValueBigFloat().String())
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data reconciliation up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data reconciliation up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "intervalInMinutes": true,
        "toleranceType": true,
        "tolerance": true,
        "sourceDataAssetId": true,
        "sourceMeasure": true,
        "sourceColumnName": true,
        "sourceCustomSql": true,
        "sourceDataCheckId": true,
        "sourceValue": true,
        "targetDataAssetId": true,
        "targetMeasure": true,
        "targetColumnName": true,
        "targetCustomSql": true,
        "targetDataCheckId": true,
        "targetValue": true,
        "difference": true,
        "differencePercent": true,
        "status": true,
        "statusMessage": true,
        "comparedAt": true,
        "statusChangedAt": true,
        "monitorId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-reconciliation/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_reconciliation, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data reconciliation found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_reconciliation: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-reconciliation/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_reconciliation, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_reconciliation: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data reconciliation matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data reconciliation matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_reconciliation.")
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
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    data.IntervalInMinutes = numberFromAPI(data.IntervalInMinutes, item["intervalInMinutes"])
    if obj, ok := item["toleranceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToleranceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToleranceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToleranceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToleranceType = types.StringValue(string(jsonBytes))
        } else {
            data.ToleranceType = types.StringNull()
        }
    } else if arr, ok := item["toleranceType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ToleranceType = types.StringValue(string(jsonBytes))
        } else {
            data.ToleranceType = types.StringNull()
        }
    } else if val, ok := item["toleranceType"].(string); ok {
        data.ToleranceType = types.StringValue(val)
    } else {
        data.ToleranceType = types.StringNull()
    }
    data.Tolerance = numberFromAPI(data.Tolerance, item["tolerance"])
    if obj, ok := item["sourceDataAssetId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SourceDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SourceDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SourceDataAssetId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SourceDataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.SourceDataAssetId = types.StringNull()
        }
    } else if arr, ok := item["sourceDataAssetId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SourceDataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.SourceDataAssetId = types.StringNull()
        }
    } else if val, ok := item["sourceDataAssetId"].(string); ok {
        data.SourceDataAssetId = types.StringValue(val)
    } else {
        data.SourceDataAssetId = types.StringNull()
    }
    if obj, ok := item["sourceMeasure"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SourceMeasure = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SourceMeasure = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SourceMeasure = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SourceMeasure = types.StringValue(string(jsonBytes))
        } else {
            data.SourceMeasure = types.StringNull()
        }
    } else if arr, ok := item["sourceMeasure"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SourceMeasure = types.StringValue(string(jsonBytes))
        } else {
            data.SourceMeasure = types.StringNull()
        }
    } else if val, ok := item["sourceMeasure"].(string); ok {
        data.SourceMeasure = types.StringValue(val)
    } else {
        data.SourceMeasure = types.StringNull()
    }
    if obj, ok := item["sourceColumnName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SourceColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SourceColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SourceColumnName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SourceColumnName = types.StringValue(string(jsonBytes))
        } else {
            data.SourceColumnName = types.StringNull()
        }
    } else if arr, ok := item["sourceColumnName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SourceColumnName = types.StringValue(string(jsonBytes))
        } else {
            data.SourceColumnName = types.StringNull()
        }
    } else if val, ok := item["sourceColumnName"].(string); ok {
        data.SourceColumnName = types.StringValue(val)
    } else {
        data.SourceColumnName = types.StringNull()
    }
    if obj, ok := item["sourceCustomSql"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SourceCustomSql = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SourceCustomSql = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SourceCustomSql = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SourceCustomSql = types.StringValue(string(jsonBytes))
        } else {
            data.SourceCustomSql = types.StringNull()
        }
    } else if arr, ok := item["sourceCustomSql"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SourceCustomSql = types.StringValue(string(jsonBytes))
        } else {
            data.SourceCustomSql = types.StringNull()
        }
    } else if val, ok := item["sourceCustomSql"].(string); ok {
        data.SourceCustomSql = types.StringValue(val)
    } else {
        data.SourceCustomSql = types.StringNull()
    }
    if obj, ok := item["sourceDataCheckId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SourceDataCheckId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SourceDataCheckId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SourceDataCheckId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SourceDataCheckId = types.StringValue(string(jsonBytes))
        } else {
            data.SourceDataCheckId = types.StringNull()
        }
    } else if arr, ok := item["sourceDataCheckId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SourceDataCheckId = types.StringValue(string(jsonBytes))
        } else {
            data.SourceDataCheckId = types.StringNull()
        }
    } else if val, ok := item["sourceDataCheckId"].(string); ok {
        data.SourceDataCheckId = types.StringValue(val)
    } else {
        data.SourceDataCheckId = types.StringNull()
    }
    data.SourceValue = numberFromAPI(data.SourceValue, item["sourceValue"])
    if obj, ok := item["targetDataAssetId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetDataAssetId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetDataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.TargetDataAssetId = types.StringNull()
        }
    } else if arr, ok := item["targetDataAssetId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetDataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.TargetDataAssetId = types.StringNull()
        }
    } else if val, ok := item["targetDataAssetId"].(string); ok {
        data.TargetDataAssetId = types.StringValue(val)
    } else {
        data.TargetDataAssetId = types.StringNull()
    }
    if obj, ok := item["targetMeasure"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetMeasure = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetMeasure = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetMeasure = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetMeasure = types.StringValue(string(jsonBytes))
        } else {
            data.TargetMeasure = types.StringNull()
        }
    } else if arr, ok := item["targetMeasure"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetMeasure = types.StringValue(string(jsonBytes))
        } else {
            data.TargetMeasure = types.StringNull()
        }
    } else if val, ok := item["targetMeasure"].(string); ok {
        data.TargetMeasure = types.StringValue(val)
    } else {
        data.TargetMeasure = types.StringNull()
    }
    if obj, ok := item["targetColumnName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetColumnName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetColumnName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetColumnName = types.StringValue(string(jsonBytes))
        } else {
            data.TargetColumnName = types.StringNull()
        }
    } else if arr, ok := item["targetColumnName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetColumnName = types.StringValue(string(jsonBytes))
        } else {
            data.TargetColumnName = types.StringNull()
        }
    } else if val, ok := item["targetColumnName"].(string); ok {
        data.TargetColumnName = types.StringValue(val)
    } else {
        data.TargetColumnName = types.StringNull()
    }
    if obj, ok := item["targetCustomSql"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetCustomSql = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetCustomSql = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetCustomSql = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetCustomSql = types.StringValue(string(jsonBytes))
        } else {
            data.TargetCustomSql = types.StringNull()
        }
    } else if arr, ok := item["targetCustomSql"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetCustomSql = types.StringValue(string(jsonBytes))
        } else {
            data.TargetCustomSql = types.StringNull()
        }
    } else if val, ok := item["targetCustomSql"].(string); ok {
        data.TargetCustomSql = types.StringValue(val)
    } else {
        data.TargetCustomSql = types.StringNull()
    }
    if obj, ok := item["targetDataCheckId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetDataCheckId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetDataCheckId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetDataCheckId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetDataCheckId = types.StringValue(string(jsonBytes))
        } else {
            data.TargetDataCheckId = types.StringNull()
        }
    } else if arr, ok := item["targetDataCheckId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetDataCheckId = types.StringValue(string(jsonBytes))
        } else {
            data.TargetDataCheckId = types.StringNull()
        }
    } else if val, ok := item["targetDataCheckId"].(string); ok {
        data.TargetDataCheckId = types.StringValue(val)
    } else {
        data.TargetDataCheckId = types.StringNull()
    }
    data.TargetValue = numberFromAPI(data.TargetValue, item["targetValue"])
    data.Difference = numberFromAPI(data.Difference, item["difference"])
    data.DifferencePercent = numberFromAPI(data.DifferencePercent, item["differencePercent"])
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if arr, ok := item["status"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if arr, ok := item["statusMessage"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["comparedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ComparedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ComparedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ComparedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ComparedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ComparedAt = types.StringNull()
        }
    } else if arr, ok := item["comparedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ComparedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ComparedAt = types.StringNull()
        }
    } else if val, ok := item["comparedAt"].(string); ok {
        data.ComparedAt = types.StringValue(val)
    } else {
        data.ComparedAt = types.StringNull()
    }
    if obj, ok := item["statusChangedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusChangedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StatusChangedAt = types.StringNull()
        }
    } else if arr, ok := item["statusChangedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StatusChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StatusChangedAt = types.StringNull()
        }
    } else if val, ok := item["statusChangedAt"].(string); ok {
        data.StatusChangedAt = types.StringValue(val)
    } else {
        data.StatusChangedAt = types.StringNull()
    }
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
    } else if arr, ok := item["monitorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
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
