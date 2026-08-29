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
var _ datasource.DataSource = &IncidentMeasurementDataSource{}

func NewIncidentMeasurementDataSource() datasource.DataSource {
    return &IncidentMeasurementDataSource{}
}

// IncidentMeasurementDataSource defines the data source implementation.
type IncidentMeasurementDataSource struct {
    client *Client
}

// IncidentMeasurementDataSourceModel describes the data source data model.
type IncidentMeasurementDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Key types.String `tfsdk:"key"`
    Description types.String `tfsdk:"description"`
    MetricName types.String `tfsdk:"metric_name"`
    StartAnchorType types.String `tfsdk:"start_anchor_type"`
    EndAnchorType types.String `tfsdk:"end_anchor_type"`
    StartIncidentStateId types.String `tfsdk:"start_incident_state_id"`
    EndIncidentStateId types.String `tfsdk:"end_incident_state_id"`
    StartIncidentStateRole types.String `tfsdk:"start_incident_state_role"`
    EndIncidentStateRole types.String `tfsdk:"end_incident_state_role"`
    StartStateOccurrence types.String `tfsdk:"start_state_occurrence"`
    EndStateOccurrence types.String `tfsdk:"end_state_occurrence"`
    Unit types.String `tfsdk:"unit"`
    AggregationType types.String `tfsdk:"aggregation_type"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ShowOnIncidentView types.Bool `tfsdk:"show_on_incident_view"`
    Order types.Number `tfsdk:"order"`
    IsSystemDefined types.Bool `tfsdk:"is_system_defined"`
    BackfillRequestedAt types.String `tfsdk:"backfill_requested_at"`
    BackfillCursorCreatedAt types.String `tfsdk:"backfill_cursor_created_at"`
    BackfillCompletedAt types.String `tfsdk:"backfill_completed_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentMeasurementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_measurement"
}

func (d *IncidentMeasurementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A named duration between two points in an incident's life, computed automatically for every incident Look up an existing incident_measurement by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Stable, machine readable identifier for this measurement, unique within the project. It is immutable once created because it is used to build the metric name that every recorded point is written under - changing it would orphan all the history. Pick it carefully; to rename a measurement, change the Name instead..",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this measurement means to your team.",
                Computed: true,
            },
            "metric_name": schema.StringAttribute{
                MarkdownDescription: "The metric name every recorded point of this measurement is written under. Derived from the key as oneuptime.incident.measurement.<key> and maintained for you..",
                Computed: true,
            },
            "start_anchor_type": schema.StringAttribute{
                MarkdownDescription: "Where this measurement starts. One of: Impact Started At, Declared At, Created At, Timeline Start, State Entered, State Role Entered, Postmortem Posted At..",
                Computed: true,
            },
            "end_anchor_type": schema.StringAttribute{
                MarkdownDescription: "Where this measurement ends. One of: Impact Started At, Declared At, Created At, Timeline Start, State Entered, State Role Entered, Postmortem Posted At..",
                Computed: true,
            },
            "start_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "end_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "start_incident_state_role": schema.StringAttribute{
                MarkdownDescription: "The role of the state this measurement starts at - Created, Acknowledged or Resolved. Used when the Start Anchor Type is State Role Entered. Resolving by role keeps the measurement working when a project renames or replaces the state that plays that part..",
                Computed: true,
            },
            "end_incident_state_role": schema.StringAttribute{
                MarkdownDescription: "The role of the state this measurement ends at - Created, Acknowledged or Resolved. Used when the End Anchor Type is State Role Entered..",
                Computed: true,
            },
            "start_state_occurrence": schema.StringAttribute{
                MarkdownDescription: "Which entry to use when the start state is entered more than once - First or Last. First matches the built-in incident metrics; Last follows a reopened incident to its final pass through that state..",
                Computed: true,
            },
            "end_state_occurrence": schema.StringAttribute{
                MarkdownDescription: "Which entry to use when the end state is entered more than once - First or Last. First matches the built-in incident metrics; Last follows a reopened incident to its final pass through that state..",
                Computed: true,
            },
            "unit": schema.StringAttribute{
                MarkdownDescription: "The unit this measurement's values are displayed in. Values are always stored in seconds; this only changes how they are rendered..",
                Computed: true,
            },
            "aggregation_type": schema.StringAttribute{
                MarkdownDescription: "The aggregation this measurement's charts default to - Avg, Max, Min, P50, P90, P95 or P99. Sum is deliberately absent because summing durations across incidents produces a number with no meaning..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement is computed for new and updated incidents.",
                Computed: true,
            },
            "show_on_incident_view": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement is shown on the incident page alongside the incident's other timings.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order in which this measurement is displayed. Lowest first..",
                Computed: true,
            },
            "is_system_defined": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement was seeded by OneUptime rather than created by your team.",
                Computed: true,
            },
            "backfill_requested_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "backfill_cursor_created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "backfill_completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentMeasurementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentMeasurementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentMeasurementDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a incident_measurement.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "key": true,
        "description": true,
        "metricName": true,
        "startAnchorType": true,
        "endAnchorType": true,
        "startIncidentStateId": true,
        "endIncidentStateId": true,
        "startIncidentStateRole": true,
        "endIncidentStateRole": true,
        "startStateOccurrence": true,
        "endStateOccurrence": true,
        "unit": true,
        "aggregationType": true,
        "isEnabled": true,
        "showOnIncidentView": true,
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
        readPath := "/incident-measurement/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_measurement, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_measurement found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_measurement: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-measurement/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_measurement, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_measurement: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_measurement found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident_measurement matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_measurement.")
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["startIncidentStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.StartIncidentStateId = types.StringNull()
        }
    } else if val, ok := item["startIncidentStateId"].(string); ok {
        data.StartIncidentStateId = types.StringValue(val)
    } else {
        data.StartIncidentStateId = types.StringNull()
    }
    if obj, ok := item["endIncidentStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.EndIncidentStateId = types.StringNull()
        }
    } else if val, ok := item["endIncidentStateId"].(string); ok {
        data.EndIncidentStateId = types.StringValue(val)
    } else {
        data.EndIncidentStateId = types.StringNull()
    }
    if obj, ok := item["startIncidentStateRole"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartIncidentStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartIncidentStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartIncidentStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartIncidentStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.StartIncidentStateRole = types.StringNull()
        }
    } else if val, ok := item["startIncidentStateRole"].(string); ok {
        data.StartIncidentStateRole = types.StringValue(val)
    } else {
        data.StartIncidentStateRole = types.StringNull()
    }
    if obj, ok := item["endIncidentStateRole"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndIncidentStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndIncidentStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndIncidentStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndIncidentStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.EndIncidentStateRole = types.StringNull()
        }
    } else if val, ok := item["endIncidentStateRole"].(string); ok {
        data.EndIncidentStateRole = types.StringValue(val)
    } else {
        data.EndIncidentStateRole = types.StringNull()
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
    if val, ok := item["showOnIncidentView"].(bool); ok {
        data.ShowOnIncidentView = types.BoolValue(val)
    } else {
        data.ShowOnIncidentView = types.BoolNull()
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
