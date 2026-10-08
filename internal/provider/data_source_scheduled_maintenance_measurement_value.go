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
var _ datasource.DataSource = &ScheduledMaintenanceMeasurementValueDataSource{}

func NewScheduledMaintenanceMeasurementValueDataSource() datasource.DataSource {
    return &ScheduledMaintenanceMeasurementValueDataSource{}
}

// ScheduledMaintenanceMeasurementValueDataSource defines the data source implementation.
type ScheduledMaintenanceMeasurementValueDataSource struct {
    client *Client
}

// ScheduledMaintenanceMeasurementValueDataSourceModel describes the data source data model.
type ScheduledMaintenanceMeasurementValueDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    ScheduledMaintenanceMeasurementId types.String `tfsdk:"scheduled_maintenance_measurement_id"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    ValueInSeconds types.Number `tfsdk:"value_in_seconds"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    StartScheduledMaintenanceStateTimelineId types.String `tfsdk:"start_scheduled_maintenance_state_timeline_id"`
    EndScheduledMaintenanceStateTimelineId types.String `tfsdk:"end_scheduled_maintenance_state_timeline_id"`
    ComputedAt types.String `tfsdk:"computed_at"`
}

func (d *ScheduledMaintenanceMeasurementValueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_measurement_value"
}

func (d *ScheduledMaintenanceMeasurementValueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The computed value of one measurement for one scheduled maintenance event. Written by OneUptime, never edited by hand. Look up an existing scheduled maintenance measurement value by `id`, or by any of its other arguments (`end_scheduled_maintenance_state_timeline_id`, `scheduled_maintenance_id`, `scheduled_maintenance_measurement_id`, ...): each one set must match, and exactly one scheduled maintenance measurement value may match them all.",

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
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Scheduled Maintenance Event this measurement value was computed for. The ID of a `oneuptime_scheduled_maintenance_event`.",
                Optional: true,
                Computed: true,
            },
            "scheduled_maintenance_measurement_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Scheduled Maintenance Measurement definition this value was computed from. The ID of a `oneuptime_scheduled_maintenance_measurement`.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the start anchor of this measurement resolved.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "When the end anchor of this measurement resolved.",
                Computed: true,
            },
            "value_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "The measured duration in seconds. Empty unless the status is Recorded.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Outcome of evaluating this measurement - Recorded, Pending, Not Applicable or Invalid.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Why this measurement has the status it has - which anchor is still open, or why it can never resolve.",
                Optional: true,
                Computed: true,
            },
            "start_scheduled_maintenance_state_timeline_id": schema.StringAttribute{
                MarkdownDescription: "The state timeline entry the start anchor resolved to, when it resolved to one. Not a foreign key - the value survives the timeline entry being removed.",
                Optional: true,
                Computed: true,
            },
            "end_scheduled_maintenance_state_timeline_id": schema.StringAttribute{
                MarkdownDescription: "The state timeline entry the end anchor resolved to, when it resolved to one. Not a foreign key - the value survives the timeline entry being removed.",
                Optional: true,
                Computed: true,
            },
            "computed_at": schema.StringAttribute{
                MarkdownDescription: "When this measurement value was last computed.",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceMeasurementValueDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceMeasurementValueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceMeasurementValueDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ScheduledMaintenanceId.IsNull() && !data.ScheduledMaintenanceId.IsUnknown() {
        filters["scheduledMaintenanceId"] = data.ScheduledMaintenanceId.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_id = "+fmt.Sprintf("%q", data.ScheduledMaintenanceId.ValueString()))
    }
    if !data.ScheduledMaintenanceMeasurementId.IsNull() && !data.ScheduledMaintenanceMeasurementId.IsUnknown() {
        filters["scheduledMaintenanceMeasurementId"] = data.ScheduledMaintenanceMeasurementId.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_measurement_id = "+fmt.Sprintf("%q", data.ScheduledMaintenanceMeasurementId.ValueString()))
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
    if !data.StartScheduledMaintenanceStateTimelineId.IsNull() && !data.StartScheduledMaintenanceStateTimelineId.IsUnknown() {
        filters["startScheduledMaintenanceStateTimelineId"] = data.StartScheduledMaintenanceStateTimelineId.ValueString()
        filterNames = append(filterNames, "start_scheduled_maintenance_state_timeline_id = "+fmt.Sprintf("%q", data.StartScheduledMaintenanceStateTimelineId.ValueString()))
    }
    if !data.EndScheduledMaintenanceStateTimelineId.IsNull() && !data.EndScheduledMaintenanceStateTimelineId.IsUnknown() {
        filters["endScheduledMaintenanceStateTimelineId"] = data.EndScheduledMaintenanceStateTimelineId.ValueString()
        filterNames = append(filterNames, "end_scheduled_maintenance_state_timeline_id = "+fmt.Sprintf("%q", data.EndScheduledMaintenanceStateTimelineId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the scheduled maintenance measurement value up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the scheduled maintenance measurement value up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "scheduledMaintenanceId": true,
        "scheduledMaintenanceMeasurementId": true,
        "startedAt": true,
        "endedAt": true,
        "valueInSeconds": true,
        "status": true,
        "statusMessage": true,
        "startScheduledMaintenanceStateTimelineId": true,
        "endScheduledMaintenanceStateTimelineId": true,
        "computedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/scheduled-maintenance-measurement-value/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement_value, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance measurement value found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement_value: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/scheduled-maintenance-measurement-value/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scheduled_maintenance_measurement_value, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list scheduled_maintenance_measurement_value: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance measurement value matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one scheduled maintenance measurement value matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for scheduled_maintenance_measurement_value.")
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
    if obj, ok := item["scheduledMaintenanceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
    }
    if obj, ok := item["scheduledMaintenanceMeasurementId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceMeasurementId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceMeasurementId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceMeasurementId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceMeasurementId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceMeasurementId = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceMeasurementId"].(string); ok {
        data.ScheduledMaintenanceMeasurementId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceMeasurementId = types.StringNull()
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
    if obj, ok := item["startScheduledMaintenanceStateTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartScheduledMaintenanceStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartScheduledMaintenanceStateTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartScheduledMaintenanceStateTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateTimelineId = types.StringNull()
        }
    } else if val, ok := item["startScheduledMaintenanceStateTimelineId"].(string); ok {
        data.StartScheduledMaintenanceStateTimelineId = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateTimelineId = types.StringNull()
    }
    if obj, ok := item["endScheduledMaintenanceStateTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndScheduledMaintenanceStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndScheduledMaintenanceStateTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndScheduledMaintenanceStateTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateTimelineId = types.StringNull()
        }
    } else if val, ok := item["endScheduledMaintenanceStateTimelineId"].(string); ok {
        data.EndScheduledMaintenanceStateTimelineId = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateTimelineId = types.StringNull()
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
