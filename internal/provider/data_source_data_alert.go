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
var _ datasource.DataSource = &DataAlertDataSource{}

func NewDataAlertDataSource() datasource.DataSource {
    return &DataAlertDataSource{}
}

// DataAlertDataSource defines the data source implementation.
type DataAlertDataSource struct {
    client *Client
}

// DataAlertDataSourceModel describes the data source data model.
type DataAlertDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    AlertId types.String `tfsdk:"alert_id"`
    DataAssetId types.String `tfsdk:"data_asset_id"`
    DataPipelineId types.String `tfsdk:"data_pipeline_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    CauseAlertId types.String `tfsdk:"cause_alert_id"`
    Feedback types.String `tfsdk:"feedback"`
    FeedbackAt types.String `tfsdk:"feedback_at"`
    FeedbackByUserId types.String `tfsdk:"feedback_by_user_id"`
    IsCauseConfirmed types.Bool `tfsdk:"is_cause_confirmed"`
    DetectedAfterInSeconds types.Number `tfsdk:"detected_after_in_seconds"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataAlertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_alert"
}

func (d *DataAlertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What a data monitor's alert is about - a table, a pipeline or a store - and the alert it is a consequence of, found by lineage. Look up an existing data alert by `id`, or by any of its other arguments (`alert_id`, `cause_alert_id`, `created_by_user_id`, ...): each one set must match, and exactly one data alert may match them all.",

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
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the alert this is about. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_alert"),
            },
            "data_asset_id": schema.StringAttribute{
                MarkdownDescription: "ID of the table the alert is about, when it is about one. The ID of a `oneuptime_data_asset`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_asset"),
            },
            "data_pipeline_id": schema.StringAttribute{
                MarkdownDescription: "ID of the pipeline the alert is about, when it is about one. The ID of a `oneuptime_data_pipeline` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_data_pipeline"),
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "ID of the database the alert's table lives in, or of the database itself when the whole of it stopped reporting. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_store"),
            },
            "cause_alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the alert this one is a consequence of, found by lineage. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_alert"),
            },
            "feedback": schema.StringAttribute{
                MarkdownDescription: "What whoever handled the alert said of it: useful, expected, or not an issue.",
                Optional: true,
                Computed: true,
            },
            "feedback_at": schema.StringAttribute{
                MarkdownDescription: "When the feedback was given.",
                Computed: true,
            },
            "feedback_by_user_id": schema.StringAttribute{
                MarkdownDescription: "Who gave the feedback.",
                Optional: true,
                Computed: true,
            },
            "is_cause_confirmed": schema.BoolAttribute{
                MarkdownDescription: "Whether the likely cause OneUptime named was the real one, once someone says.",
                Optional: true,
                Computed: true,
            },
            "detected_after_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long after a table's update was due the alert opened; empty for an alert about anything else.",
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

func (d *DataAlertDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataAlertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataAlertDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.AlertId.IsNull() && !data.AlertId.IsUnknown() {
        filters["alertId"] = data.AlertId.ValueString()
        filterNames = append(filterNames, "alert_id = "+fmt.Sprintf("%q", data.AlertId.ValueString()))
    }
    if !data.DataAssetId.IsNull() && !data.DataAssetId.IsUnknown() {
        filters["dataAssetId"] = data.DataAssetId.ValueString()
        filterNames = append(filterNames, "data_asset_id = "+fmt.Sprintf("%q", data.DataAssetId.ValueString()))
    }
    if !data.DataPipelineId.IsNull() && !data.DataPipelineId.IsUnknown() {
        filters["dataPipelineId"] = data.DataPipelineId.ValueString()
        filterNames = append(filterNames, "data_pipeline_id = "+fmt.Sprintf("%q", data.DataPipelineId.ValueString()))
    }
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.CauseAlertId.IsNull() && !data.CauseAlertId.IsUnknown() {
        filters["causeAlertId"] = data.CauseAlertId.ValueString()
        filterNames = append(filterNames, "cause_alert_id = "+fmt.Sprintf("%q", data.CauseAlertId.ValueString()))
    }
    if !data.Feedback.IsNull() && !data.Feedback.IsUnknown() {
        filters["feedback"] = data.Feedback.ValueString()
        filterNames = append(filterNames, "feedback = "+fmt.Sprintf("%q", data.Feedback.ValueString()))
    }
    if !data.FeedbackByUserId.IsNull() && !data.FeedbackByUserId.IsUnknown() {
        filters["feedbackByUserId"] = data.FeedbackByUserId.ValueString()
        filterNames = append(filterNames, "feedback_by_user_id = "+fmt.Sprintf("%q", data.FeedbackByUserId.ValueString()))
    }
    if !data.IsCauseConfirmed.IsNull() && !data.IsCauseConfirmed.IsUnknown() {
        filters["isCauseConfirmed"] = data.IsCauseConfirmed.ValueBool()
        filterNames = append(filterNames, "is_cause_confirmed = "+fmt.Sprintf("%t", data.IsCauseConfirmed.ValueBool()))
    }
    if !data.DetectedAfterInSeconds.IsNull() && !data.DetectedAfterInSeconds.IsUnknown() {
        filters["detectedAfterInSeconds"] = lookupNumber(data.DetectedAfterInSeconds)
        filterNames = append(filterNames, "detected_after_in_seconds = "+data.DetectedAfterInSeconds.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data alert up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data alert up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "alertId": true,
        "dataAssetId": true,
        "dataPipelineId": true,
        "dataStoreId": true,
        "causeAlertId": true,
        "feedback": true,
        "feedbackAt": true,
        "feedbackByUserId": true,
        "isCauseConfirmed": true,
        "detectedAfterInSeconds": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-alert/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_alert, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data alert found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_alert: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-alert/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_alert, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_alert: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data alert matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data alert matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_alert.")
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
    if obj, ok := item["alertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := item["alertId"].(string); ok {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
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
    if obj, ok := item["dataPipelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataPipelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataPipelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataPipelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataPipelineId = types.StringValue(string(jsonBytes))
        } else {
            data.DataPipelineId = types.StringNull()
        }
    } else if val, ok := item["dataPipelineId"].(string); ok {
        data.DataPipelineId = types.StringValue(val)
    } else {
        data.DataPipelineId = types.StringNull()
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
    if obj, ok := item["causeAlertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CauseAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CauseAlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CauseAlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CauseAlertId = types.StringValue(string(jsonBytes))
        } else {
            data.CauseAlertId = types.StringNull()
        }
    } else if val, ok := item["causeAlertId"].(string); ok {
        data.CauseAlertId = types.StringValue(val)
    } else {
        data.CauseAlertId = types.StringNull()
    }
    if obj, ok := item["feedback"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Feedback = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Feedback = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Feedback = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Feedback = types.StringValue(string(jsonBytes))
        } else {
            data.Feedback = types.StringNull()
        }
    } else if val, ok := item["feedback"].(string); ok {
        data.Feedback = types.StringValue(val)
    } else {
        data.Feedback = types.StringNull()
    }
    if obj, ok := item["feedbackAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FeedbackAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FeedbackAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FeedbackAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FeedbackAt = types.StringValue(string(jsonBytes))
        } else {
            data.FeedbackAt = types.StringNull()
        }
    } else if val, ok := item["feedbackAt"].(string); ok {
        data.FeedbackAt = types.StringValue(val)
    } else {
        data.FeedbackAt = types.StringNull()
    }
    if obj, ok := item["feedbackByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FeedbackByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FeedbackByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FeedbackByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FeedbackByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.FeedbackByUserId = types.StringNull()
        }
    } else if val, ok := item["feedbackByUserId"].(string); ok {
        data.FeedbackByUserId = types.StringValue(val)
    } else {
        data.FeedbackByUserId = types.StringNull()
    }
    if val, ok := item["isCauseConfirmed"].(bool); ok {
        data.IsCauseConfirmed = types.BoolValue(val)
    } else {
        data.IsCauseConfirmed = types.BoolNull()
    }
    data.DetectedAfterInSeconds = numberFromAPI(data.DetectedAfterInSeconds, item["detectedAfterInSeconds"])
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
