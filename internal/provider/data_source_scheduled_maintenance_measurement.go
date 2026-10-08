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
var _ datasource.DataSource = &ScheduledMaintenanceMeasurementDataSource{}

func NewScheduledMaintenanceMeasurementDataSource() datasource.DataSource {
    return &ScheduledMaintenanceMeasurementDataSource{}
}

// ScheduledMaintenanceMeasurementDataSource defines the data source implementation.
type ScheduledMaintenanceMeasurementDataSource struct {
    client *Client
}

// ScheduledMaintenanceMeasurementDataSourceModel describes the data source data model.
type ScheduledMaintenanceMeasurementDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Key types.String `tfsdk:"key"`
    Description types.String `tfsdk:"description"`
    MetricName types.String `tfsdk:"metric_name"`
    StartAnchorType types.String `tfsdk:"start_anchor_type"`
    EndAnchorType types.String `tfsdk:"end_anchor_type"`
    StartScheduledMaintenanceStateId types.String `tfsdk:"start_scheduled_maintenance_state_id"`
    EndScheduledMaintenanceStateId types.String `tfsdk:"end_scheduled_maintenance_state_id"`
    StartScheduledMaintenanceStateRole types.String `tfsdk:"start_scheduled_maintenance_state_role"`
    EndScheduledMaintenanceStateRole types.String `tfsdk:"end_scheduled_maintenance_state_role"`
    StartStateOccurrence types.String `tfsdk:"start_state_occurrence"`
    EndStateOccurrence types.String `tfsdk:"end_state_occurrence"`
    Unit types.String `tfsdk:"unit"`
    AggregationType types.String `tfsdk:"aggregation_type"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ShowOnScheduledMaintenanceView types.Bool `tfsdk:"show_on_scheduled_maintenance_view"`
    Order types.Number `tfsdk:"order"`
    IsSystemDefined types.Bool `tfsdk:"is_system_defined"`
    BackfillRequestedAt types.String `tfsdk:"backfill_requested_at"`
    BackfillCursorCreatedAt types.String `tfsdk:"backfill_cursor_created_at"`
    BackfillCompletedAt types.String `tfsdk:"backfill_completed_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ScheduledMaintenanceMeasurementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_measurement"
}

func (d *ScheduledMaintenanceMeasurementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A named duration between two points in a scheduled maintenance event's life, computed automatically for every event Look up an existing scheduled maintenance measurement by `id`, or by any of its other arguments (`name`, `aggregation_type`, `created_by_user_id`, ...): each one set must match, and exactly one scheduled maintenance measurement may match them all.",

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
                MarkdownDescription: "Human readable name of this measurement, such as Start Delay. This is what charts call it.",
                Optional: true,
                Computed: true,
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Stable, machine readable identifier for this measurement, unique within the project: lowercase letters, numbers and hyphens. Leave it out and it is made from the name - Time to Start becomes time-to-start, with -2, -3 and so on added when another measurement already has it. It is part of the metric name, so it cannot be changed once the measurement is created; to rename a measurement, change the Name instead.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this measurement means to your team.",
                Optional: true,
                Computed: true,
            },
            "metric_name": schema.StringAttribute{
                MarkdownDescription: "Name of the metric this measurement writes to, derived from the key as oneuptime.scheduled-maintenance.measurement.<key>.",
                Optional: true,
                Computed: true,
            },
            "start_anchor_type": schema.StringAttribute{
                MarkdownDescription: "Where the measurement starts - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role.",
                Optional: true,
                Computed: true,
            },
            "end_anchor_type": schema.StringAttribute{
                MarkdownDescription: "Where the measurement ends - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role.",
                Optional: true,
                Computed: true,
            },
            "start_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "ID of the state whose entry starts this measurement, when the start anchor is a specific state. The ID of a `oneuptime_scheduled_maintenance_state`.",
                Optional: true,
                Computed: true,
            },
            "end_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "ID of the state whose entry ends this measurement, when the end anchor is a specific state. The ID of a `oneuptime_scheduled_maintenance_state`.",
                Optional: true,
                Computed: true,
            },
            "start_scheduled_maintenance_state_role": schema.StringAttribute{
                MarkdownDescription: "The role of the state that starts this measurement (Scheduled, Ongoing, Ended or Resolved), when the start anchor is a state role. Resolving by role keeps working when a project renames or replaces the state.",
                Optional: true,
                Computed: true,
            },
            "end_scheduled_maintenance_state_role": schema.StringAttribute{
                MarkdownDescription: "The role of the state that ends this measurement (Scheduled, Ongoing, Ended or Resolved), when the end anchor is a state role. Resolving by role keeps working when a project renames or replaces the state.",
                Optional: true,
                Computed: true,
            },
            "start_state_occurrence": schema.StringAttribute{
                MarkdownDescription: "Which entry to use when the start state is entered more than once - the first time it was entered, or the last.",
                Optional: true,
                Computed: true,
            },
            "end_state_occurrence": schema.StringAttribute{
                MarkdownDescription: "Which entry to use when the end state is entered more than once - the first time it was entered, or the last.",
                Optional: true,
                Computed: true,
            },
            "unit": schema.StringAttribute{
                MarkdownDescription: "The unit this measurement's charts are in: seconds (the default), minutes, hours or days. Every value is worked out in seconds and stored that way on the scheduled maintenance event; each chart point is written in this unit, so a chart in hours reads 1.5 for an hour and a half. With seconds, charts show seconds, minutes, hours or days as the numbers grow. Changing it rewrites the measurement's chart points in the new unit. A value that is not a time unit charts in seconds.",
                Optional: true,
                Computed: true,
            },
            "aggregation_type": schema.StringAttribute{
                MarkdownDescription: "How this measurement's chart sums up many scheduled maintenance events by default - Avg (the default), P50, P90, P95, P99, Max or Min. View Chart in the dashboard opens the chart this way. Sum is deliberately absent: adding durations up across scheduled maintenance events produces a number with no meaning.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement is computed for scheduled maintenance events.",
                Optional: true,
                Computed: true,
            },
            "show_on_scheduled_maintenance_view": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement is shown on the scheduled maintenance event page.",
                Optional: true,
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Where this measurement appears in the list of measurements, lowest number first. A new measurement is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "is_system_defined": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement was created by OneUptime rather than by your team.",
                Optional: true,
                Computed: true,
            },
            "backfill_requested_at": schema.StringAttribute{
                MarkdownDescription: "When a backfill of this measurement over existing scheduled maintenance events was requested.",
                Computed: true,
            },
            "backfill_cursor_created_at": schema.StringAttribute{
                MarkdownDescription: "How far the backfill has walked this project, so a restart resumes instead of starting over.",
                Computed: true,
            },
            "backfill_completed_at": schema.StringAttribute{
                MarkdownDescription: "When the backfill of this measurement over existing scheduled maintenance events finished.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceMeasurementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceMeasurementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceMeasurementDataSourceModel

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
    if !data.Key.IsNull() && !data.Key.IsUnknown() {
        filters["key"] = data.Key.ValueString()
        filterNames = append(filterNames, "key = "+fmt.Sprintf("%q", data.Key.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.MetricName.IsNull() && !data.MetricName.IsUnknown() {
        filters["metricName"] = data.MetricName.ValueString()
        filterNames = append(filterNames, "metric_name = "+fmt.Sprintf("%q", data.MetricName.ValueString()))
    }
    if !data.StartAnchorType.IsNull() && !data.StartAnchorType.IsUnknown() {
        filters["startAnchorType"] = data.StartAnchorType.ValueString()
        filterNames = append(filterNames, "start_anchor_type = "+fmt.Sprintf("%q", data.StartAnchorType.ValueString()))
    }
    if !data.EndAnchorType.IsNull() && !data.EndAnchorType.IsUnknown() {
        filters["endAnchorType"] = data.EndAnchorType.ValueString()
        filterNames = append(filterNames, "end_anchor_type = "+fmt.Sprintf("%q", data.EndAnchorType.ValueString()))
    }
    if !data.StartScheduledMaintenanceStateId.IsNull() && !data.StartScheduledMaintenanceStateId.IsUnknown() {
        filters["startScheduledMaintenanceStateId"] = data.StartScheduledMaintenanceStateId.ValueString()
        filterNames = append(filterNames, "start_scheduled_maintenance_state_id = "+fmt.Sprintf("%q", data.StartScheduledMaintenanceStateId.ValueString()))
    }
    if !data.EndScheduledMaintenanceStateId.IsNull() && !data.EndScheduledMaintenanceStateId.IsUnknown() {
        filters["endScheduledMaintenanceStateId"] = data.EndScheduledMaintenanceStateId.ValueString()
        filterNames = append(filterNames, "end_scheduled_maintenance_state_id = "+fmt.Sprintf("%q", data.EndScheduledMaintenanceStateId.ValueString()))
    }
    if !data.StartScheduledMaintenanceStateRole.IsNull() && !data.StartScheduledMaintenanceStateRole.IsUnknown() {
        filters["startScheduledMaintenanceStateRole"] = data.StartScheduledMaintenanceStateRole.ValueString()
        filterNames = append(filterNames, "start_scheduled_maintenance_state_role = "+fmt.Sprintf("%q", data.StartScheduledMaintenanceStateRole.ValueString()))
    }
    if !data.EndScheduledMaintenanceStateRole.IsNull() && !data.EndScheduledMaintenanceStateRole.IsUnknown() {
        filters["endScheduledMaintenanceStateRole"] = data.EndScheduledMaintenanceStateRole.ValueString()
        filterNames = append(filterNames, "end_scheduled_maintenance_state_role = "+fmt.Sprintf("%q", data.EndScheduledMaintenanceStateRole.ValueString()))
    }
    if !data.StartStateOccurrence.IsNull() && !data.StartStateOccurrence.IsUnknown() {
        filters["startStateOccurrence"] = data.StartStateOccurrence.ValueString()
        filterNames = append(filterNames, "start_state_occurrence = "+fmt.Sprintf("%q", data.StartStateOccurrence.ValueString()))
    }
    if !data.EndStateOccurrence.IsNull() && !data.EndStateOccurrence.IsUnknown() {
        filters["endStateOccurrence"] = data.EndStateOccurrence.ValueString()
        filterNames = append(filterNames, "end_state_occurrence = "+fmt.Sprintf("%q", data.EndStateOccurrence.ValueString()))
    }
    if !data.Unit.IsNull() && !data.Unit.IsUnknown() {
        filters["unit"] = data.Unit.ValueString()
        filterNames = append(filterNames, "unit = "+fmt.Sprintf("%q", data.Unit.ValueString()))
    }
    if !data.AggregationType.IsNull() && !data.AggregationType.IsUnknown() {
        filters["aggregationType"] = data.AggregationType.ValueString()
        filterNames = append(filterNames, "aggregation_type = "+fmt.Sprintf("%q", data.AggregationType.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.ShowOnScheduledMaintenanceView.IsNull() && !data.ShowOnScheduledMaintenanceView.IsUnknown() {
        filters["showOnScheduledMaintenanceView"] = data.ShowOnScheduledMaintenanceView.ValueBool()
        filterNames = append(filterNames, "show_on_scheduled_maintenance_view = "+fmt.Sprintf("%t", data.ShowOnScheduledMaintenanceView.ValueBool()))
    }
    if !data.Order.IsNull() && !data.Order.IsUnknown() {
        filters["order"] = lookupNumber(data.Order)
        filterNames = append(filterNames, "order = "+data.Order.ValueBigFloat().String())
    }
    if !data.IsSystemDefined.IsNull() && !data.IsSystemDefined.IsUnknown() {
        filters["isSystemDefined"] = data.IsSystemDefined.ValueBool()
        filterNames = append(filterNames, "is_system_defined = "+fmt.Sprintf("%t", data.IsSystemDefined.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the scheduled maintenance measurement up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the scheduled maintenance measurement up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "key": true,
        "description": true,
        "metricName": true,
        "startAnchorType": true,
        "endAnchorType": true,
        "startScheduledMaintenanceStateId": true,
        "endScheduledMaintenanceStateId": true,
        "startScheduledMaintenanceStateRole": true,
        "endScheduledMaintenanceStateRole": true,
        "startStateOccurrence": true,
        "endStateOccurrence": true,
        "unit": true,
        "aggregationType": true,
        "isEnabled": true,
        "showOnScheduledMaintenanceView": true,
        "order": true,
        "isSystemDefined": true,
        "backfillRequestedAt": true,
        "backfillCursorCreatedAt": true,
        "backfillCompletedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/scheduled-maintenance-measurement/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance measurement found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/scheduled-maintenance-measurement/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scheduled_maintenance_measurement, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list scheduled_maintenance_measurement: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance measurement matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one scheduled maintenance measurement matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for scheduled_maintenance_measurement.")
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
    if obj, ok := item["key"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Key = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if val, ok := item["key"].(string); ok {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
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
    if obj, ok := item["metricName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MetricName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MetricName = types.StringValue(string(jsonBytes))
        } else {
            data.MetricName = types.StringNull()
        }
    } else if val, ok := item["metricName"].(string); ok {
        data.MetricName = types.StringValue(val)
    } else {
        data.MetricName = types.StringNull()
    }
    if obj, ok := item["startAnchorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.StartAnchorType = types.StringNull()
        }
    } else if val, ok := item["startAnchorType"].(string); ok {
        data.StartAnchorType = types.StringValue(val)
    } else {
        data.StartAnchorType = types.StringNull()
    }
    if obj, ok := item["endAnchorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.EndAnchorType = types.StringNull()
        }
    } else if val, ok := item["endAnchorType"].(string); ok {
        data.EndAnchorType = types.StringValue(val)
    } else {
        data.EndAnchorType = types.StringNull()
    }
    if obj, ok := item["startScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := item["startScheduledMaintenanceStateId"].(string); ok {
        data.StartScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := item["endScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := item["endScheduledMaintenanceStateId"].(string); ok {
        data.EndScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := item["startScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := item["startScheduledMaintenanceStateRole"].(string); ok {
        data.StartScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := item["endScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := item["endScheduledMaintenanceStateRole"].(string); ok {
        data.EndScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := item["startStateOccurrence"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.StartStateOccurrence = types.StringNull()
        }
    } else if val, ok := item["startStateOccurrence"].(string); ok {
        data.StartStateOccurrence = types.StringValue(val)
    } else {
        data.StartStateOccurrence = types.StringNull()
    }
    if obj, ok := item["endStateOccurrence"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.EndStateOccurrence = types.StringNull()
        }
    } else if val, ok := item["endStateOccurrence"].(string); ok {
        data.EndStateOccurrence = types.StringValue(val)
    } else {
        data.EndStateOccurrence = types.StringNull()
    }
    if obj, ok := item["unit"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Unit = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Unit = types.StringValue(string(jsonBytes))
        } else {
            data.Unit = types.StringNull()
        }
    } else if val, ok := item["unit"].(string); ok {
        data.Unit = types.StringValue(val)
    } else {
        data.Unit = types.StringNull()
    }
    if obj, ok := item["aggregationType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AggregationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AggregationType = types.StringValue(string(jsonBytes))
        } else {
            data.AggregationType = types.StringNull()
        }
    } else if val, ok := item["aggregationType"].(string); ok {
        data.AggregationType = types.StringValue(val)
    } else {
        data.AggregationType = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["showOnScheduledMaintenanceView"].(bool); ok {
        data.ShowOnScheduledMaintenanceView = types.BoolValue(val)
    } else {
        data.ShowOnScheduledMaintenanceView = types.BoolNull()
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
    if val, ok := item["isSystemDefined"].(bool); ok {
        data.IsSystemDefined = types.BoolValue(val)
    } else {
        data.IsSystemDefined = types.BoolNull()
    }
    if obj, ok := item["backfillRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BackfillRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BackfillRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BackfillRequestedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BackfillRequestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.BackfillRequestedAt = types.StringNull()
        }
    } else if val, ok := item["backfillRequestedAt"].(string); ok {
        data.BackfillRequestedAt = types.StringValue(val)
    } else {
        data.BackfillRequestedAt = types.StringNull()
    }
    if obj, ok := item["backfillCursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BackfillCursorCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BackfillCursorCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BackfillCursorCreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BackfillCursorCreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.BackfillCursorCreatedAt = types.StringNull()
        }
    } else if val, ok := item["backfillCursorCreatedAt"].(string); ok {
        data.BackfillCursorCreatedAt = types.StringValue(val)
    } else {
        data.BackfillCursorCreatedAt = types.StringNull()
    }
    if obj, ok := item["backfillCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BackfillCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BackfillCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BackfillCompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BackfillCompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.BackfillCompletedAt = types.StringNull()
        }
    } else if val, ok := item["backfillCompletedAt"].(string); ok {
        data.BackfillCompletedAt = types.StringValue(val)
    } else {
        data.BackfillCompletedAt = types.StringNull()
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
