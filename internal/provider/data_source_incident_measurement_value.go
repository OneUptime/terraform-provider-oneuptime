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
var _ datasource.DataSource = &IncidentMeasurementValueDataSource{}

func NewIncidentMeasurementValueDataSource() datasource.DataSource {
    return &IncidentMeasurementValueDataSource{}
}

// IncidentMeasurementValueDataSource defines the data source implementation.
type IncidentMeasurementValueDataSource struct {
    client *Client
}

// IncidentMeasurementValueDataSourceModel describes the data source data model.
type IncidentMeasurementValueDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    IncidentMeasurementId types.String `tfsdk:"incident_measurement_id"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    ValueInSeconds types.Number `tfsdk:"value_in_seconds"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    StartIncidentStateTimelineId types.String `tfsdk:"start_incident_state_timeline_id"`
    EndIncidentStateTimelineId types.String `tfsdk:"end_incident_state_timeline_id"`
    ComputedAt types.String `tfsdk:"computed_at"`
}

func (d *IncidentMeasurementValueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_measurement_value"
}

func (d *IncidentMeasurementValueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The computed value of one incident measurement for one incident, recomputed from the incident's timeline rather than accumulated Look up an existing incident measurement value by `id`, or by any of its other arguments (`end_incident_state_timeline_id`, `incident_id`, `incident_measurement_id`, ...): each one set must match, and exactly one incident measurement value may match them all.",

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
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident this measurement value was computed for. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "incident_measurement_id": schema.StringAttribute{
                MarkdownDescription: "ID of the measurement definition this value was computed from. The ID of a `oneuptime_incident_measurement`.",
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
            "start_incident_state_timeline_id": schema.StringAttribute{
                MarkdownDescription: "The incident state timeline entry the start anchor resolved to. Recorded for provenance only - it carries no foreign key, so deleting a timeline entry never blocks or rewrites this row; the next recompute simply produces the right answer.",
                Optional: true,
                Computed: true,
            },
            "end_incident_state_timeline_id": schema.StringAttribute{
                MarkdownDescription: "The incident state timeline entry the end anchor resolved to. Recorded for provenance only - it carries no foreign key, so deleting a timeline entry never blocks or rewrites this row; the next recompute simply produces the right answer.",
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

func (d *IncidentMeasurementValueDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentMeasurementValueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentMeasurementValueDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.IncidentMeasurementId.IsNull() && !data.IncidentMeasurementId.IsUnknown() {
        filters["incidentMeasurementId"] = data.IncidentMeasurementId.ValueString()
        filterNames = append(filterNames, "incident_measurement_id = "+fmt.Sprintf("%q", data.IncidentMeasurementId.ValueString()))
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
    if !data.StartIncidentStateTimelineId.IsNull() && !data.StartIncidentStateTimelineId.IsUnknown() {
        filters["startIncidentStateTimelineId"] = data.StartIncidentStateTimelineId.ValueString()
        filterNames = append(filterNames, "start_incident_state_timeline_id = "+fmt.Sprintf("%q", data.StartIncidentStateTimelineId.ValueString()))
    }
    if !data.EndIncidentStateTimelineId.IsNull() && !data.EndIncidentStateTimelineId.IsUnknown() {
        filters["endIncidentStateTimelineId"] = data.EndIncidentStateTimelineId.ValueString()
        filterNames = append(filterNames, "end_incident_state_timeline_id = "+fmt.Sprintf("%q", data.EndIncidentStateTimelineId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident measurement value up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident measurement value up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incidentId": true,
        "incidentMeasurementId": true,
        "startedAt": true,
        "endedAt": true,
        "valueInSeconds": true,
        "status": true,
        "statusMessage": true,
        "startIncidentStateTimelineId": true,
        "endIncidentStateTimelineId": true,
        "computedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-measurement-value/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_measurement_value, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident measurement value found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_measurement_value: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-measurement-value/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_measurement_value, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_measurement_value: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident measurement value matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident measurement value matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_measurement_value.")
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
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := item["incidentMeasurementId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentMeasurementId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentMeasurementId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentMeasurementId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentMeasurementId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentMeasurementId = types.StringNull()
        }
    } else if val, ok := item["incidentMeasurementId"].(string); ok {
        data.IncidentMeasurementId = types.StringValue(val)
    } else {
        data.IncidentMeasurementId = types.StringNull()
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
    if obj, ok := item["startIncidentStateTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartIncidentStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartIncidentStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartIncidentStateTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartIncidentStateTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.StartIncidentStateTimelineId = types.StringNull()
        }
    } else if val, ok := item["startIncidentStateTimelineId"].(string); ok {
        data.StartIncidentStateTimelineId = types.StringValue(val)
    } else {
        data.StartIncidentStateTimelineId = types.StringNull()
    }
    if obj, ok := item["endIncidentStateTimelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndIncidentStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndIncidentStateTimelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndIncidentStateTimelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndIncidentStateTimelineId = types.StringValue(string(jsonBytes))
        } else {
            data.EndIncidentStateTimelineId = types.StringNull()
        }
    } else if val, ok := item["endIncidentStateTimelineId"].(string); ok {
        data.EndIncidentStateTimelineId = types.StringValue(val)
    } else {
        data.EndIncidentStateTimelineId = types.StringNull()
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
