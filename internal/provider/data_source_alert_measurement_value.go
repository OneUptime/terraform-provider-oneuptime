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
var _ datasource.DataSource = &AlertMeasurementValueDataSource{}

func NewAlertMeasurementValueDataSource() datasource.DataSource {
    return &AlertMeasurementValueDataSource{}
}

// AlertMeasurementValueDataSource defines the data source implementation.
type AlertMeasurementValueDataSource struct {
    client *Client
}

// AlertMeasurementValueDataSourceModel describes the data source data model.
type AlertMeasurementValueDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    AlertId types.String `tfsdk:"alert_id"`
    AlertMeasurementId types.String `tfsdk:"alert_measurement_id"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    ValueInSeconds types.Number `tfsdk:"value_in_seconds"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    StartAlertStateTimelineId types.String `tfsdk:"start_alert_state_timeline_id"`
    EndAlertStateTimelineId types.String `tfsdk:"end_alert_state_timeline_id"`
    ComputedAt types.String `tfsdk:"computed_at"`
}

func (d *AlertMeasurementValueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_measurement_value"
}

func (d *AlertMeasurementValueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The computed value of one alert measurement for one alert, recomputed from the alert's timeline rather than accumulated Look up an existing alert measurement value by `id`, or by any of its other arguments (`alert_id`, `alert_measurement_id`, `end_alert_state_timeline_id`, ...): each one set must match, and exactly one alert measurement value may match them all.",

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
                MarkdownDescription: "ID of the alert this measurement value was computed for. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "alert_measurement_id": schema.StringAttribute{
                MarkdownDescription: "ID of the measurement definition this value was computed from. The ID of a `oneuptime_alert_measurement`.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When this measurement's start anchor resolved to. Blank while the start anchor has not resolved.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "When this measurement's end anchor resolved to. Blank while the end anchor has not resolved.",
                Computed: true,
            },
            "value_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "The measured duration in seconds. Only set when the status is Recorded - a measurement that could not be computed is left blank rather than written as zero.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "The outcome of evaluating this measurement: Recorded, Pending, Not Applicable or Invalid.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Why this measurement has the status it has, in plain words - for example which anchor has not been reached yet, or by how much the end precedes the start.",
                Optional: true,
                Computed: true,
            },
            "start_alert_state_timeline_id": schema.StringAttribute{
                MarkdownDescription: "The alert state timeline entry the start anchor resolved to. Recorded for provenance only - it carries no foreign key, so deleting a timeline entry never blocks or rewrites this row; the next recompute simply produces the right answer.",
                Optional: true,
                Computed: true,
            },
            "end_alert_state_timeline_id": schema.StringAttribute{
                MarkdownDescription: "The alert state timeline entry the end anchor resolved to. Recorded for provenance only - it carries no foreign key, so deleting a timeline entry never blocks or rewrites this row; the next recompute simply produces the right answer.",
                Optional: true,
                Computed: true,
            },
            "computed_at": schema.StringAttribute{
                MarkdownDescription: "When this value was last recomputed.",
                Computed: true,
            },
        },
    }
}

func (d *AlertMeasurementValueDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertMeasurementValueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertMeasurementValueDataSourceModel

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
    if !data.AlertMeasurementId.IsNull() && !data.AlertMeasurementId.IsUnknown() {
        filters["alertMeasurementId"] = data.AlertMeasurementId.ValueString()
        filterNames = append(filterNames, "alert_measurement_id = "+fmt.Sprintf("%q", data.AlertMeasurementId.ValueString()))
    }
    if !data.ValueInSeconds.IsNull() && !data.ValueInSeconds.IsUnknown() {
        filters["valueInSeconds"] = lookupNumber(data.ValueInSeconds)
        filterNames = append(filterNames, "value_in_seconds = "+data.ValueInSeconds.ValueBigFloat().String())
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.StartAlertStateTimelineId.IsNull() && !data.StartAlertStateTimelineId.IsUnknown() {
        filters["startAlertStateTimelineId"] = data.StartAlertStateTimelineId.ValueString()
        filterNames = append(filterNames, "start_alert_state_timeline_id = "+fmt.Sprintf("%q", data.StartAlertStateTimelineId.ValueString()))
    }
    if !data.EndAlertStateTimelineId.IsNull() && !data.EndAlertStateTimelineId.IsUnknown() {
        filters["endAlertStateTimelineId"] = data.EndAlertStateTimelineId.ValueString()
        filterNames = append(filterNames, "end_alert_state_timeline_id = "+fmt.Sprintf("%q", data.EndAlertStateTimelineId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the alert measurement value up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the alert measurement value up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "alertId": true,
        "alertMeasurementId": true,
        "startedAt": true,
        "endedAt": true,
        "valueInSeconds": true,
        "status": true,
        "statusMessage": true,
        "startAlertStateTimelineId": true,
        "endAlertStateTimelineId": true,
        "computedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/alert-measurement-value/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_measurement_value, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert measurement value found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read alert_measurement_value: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/alert-measurement-value/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list alert_measurement_value, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list alert_measurement_value: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert measurement value matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one alert measurement value matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for alert_measurement_value.")
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
    if obj, ok := item["alertMeasurementId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertMeasurementId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertMeasurementId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertMeasurementId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertMeasurementId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertMeasurementId = types.StringNull()
        }
    } else if val, ok := item["alertMeasurementId"].(string); ok {
        data.AlertMeasurementId = types.StringValue(val)
    } else {
        data.AlertMeasurementId = types.StringNull()
    }
    if obj, ok := item["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartedAt = types.StringNull()
        }
    } else if val, ok := item["startedAt"].(string); ok {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if obj, ok := item["endedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndedAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndedAt = types.StringNull()
        }
    } else if val, ok := item["endedAt"].(string); ok {
        data.EndedAt = types.StringValue(val)
    } else {
        data.EndedAt = types.StringNull()
    }
    if val, ok := item["valueInSeconds"].(float64); ok {
        data.ValueInSeconds = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["valueInSeconds"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ValueInSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.ValueInSeconds = types.NumberNull()
        }
    } else {
        data.ValueInSeconds = types.NumberNull()
    }
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
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["startAlertStateTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartAlertStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartAlertStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartAlertStateTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartAlertStateTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.StartAlertStateTimelineId = types.StringNull()
        }
    } else if val, ok := item["startAlertStateTimelineId"].(string); ok {
        data.StartAlertStateTimelineId = types.StringValue(val)
    } else {
        data.StartAlertStateTimelineId = types.StringNull()
    }
    if obj, ok := item["endAlertStateTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndAlertStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndAlertStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndAlertStateTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndAlertStateTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.EndAlertStateTimelineId = types.StringNull()
        }
    } else if val, ok := item["endAlertStateTimelineId"].(string); ok {
        data.EndAlertStateTimelineId = types.StringValue(val)
    } else {
        data.EndAlertStateTimelineId = types.StringNull()
    }
    if obj, ok := item["computedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ComputedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ComputedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ComputedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ComputedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ComputedAt = types.StringNull()
        }
    } else if val, ok := item["computedAt"].(string); ok {
        data.ComputedAt = types.StringValue(val)
    } else {
        data.ComputedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
