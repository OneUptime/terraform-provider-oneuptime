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
var _ datasource.DataSource = &IncidentSlaDataSource{}

func NewIncidentSlaDataSource() datasource.DataSource {
    return &IncidentSlaDataSource{}
}

// IncidentSlaDataSource defines the data source implementation.
type IncidentSlaDataSource struct {
    client *Client
}

// IncidentSlaDataSourceModel describes the data source data model.
type IncidentSlaDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    IncidentSlaRuleId types.String `tfsdk:"incident_sla_rule_id"`
    ResponseDeadline types.String `tfsdk:"response_deadline"`
    ResolutionDeadline types.String `tfsdk:"resolution_deadline"`
    Status types.String `tfsdk:"status"`
    RespondedAt types.String `tfsdk:"responded_at"`
    ResolvedAt types.String `tfsdk:"resolved_at"`
    LastInternalNoteReminderSentAt types.String `tfsdk:"last_internal_note_reminder_sent_at"`
    LastPublicNoteReminderSentAt types.String `tfsdk:"last_public_note_reminder_sent_at"`
    BreachNotificationSentAt types.String `tfsdk:"breach_notification_sent_at"`
    SlaStartedAt types.String `tfsdk:"sla_started_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentSlaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_sla"
}

func (d *IncidentSlaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Track SLA status and deadlines for incidents Look up an existing incident sla by `id`, or by any of its other arguments (`created_by_user_id`, `incident_id`, `incident_sla_rule_id`, ...): each one set must match, and exactly one incident sla may match them all.",

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
                MarkdownDescription: "ID of the incident this SLA record is tracking. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "incident_sla_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the SLA rule that was applied to this incident. The ID of a `oneuptime_incident_sla_rule`.",
                Optional: true,
                Computed: true,
            },
            "response_deadline": schema.StringAttribute{
                MarkdownDescription: "The deadline by which the incident must be acknowledged to meet the SLA.",
                Computed: true,
            },
            "resolution_deadline": schema.StringAttribute{
                MarkdownDescription: "The deadline by which the incident must be resolved to meet the SLA.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current SLA status (On Track, At Risk, Breached, Met).",
                Optional: true,
                Computed: true,
            },
            "responded_at": schema.StringAttribute{
                MarkdownDescription: "The actual time when the incident was acknowledged.",
                Computed: true,
            },
            "resolved_at": schema.StringAttribute{
                MarkdownDescription: "The actual time when the incident was resolved.",
                Computed: true,
            },
            "last_internal_note_reminder_sent_at": schema.StringAttribute{
                MarkdownDescription: "The last time an internal note reminder was sent.",
                Computed: true,
            },
            "last_public_note_reminder_sent_at": schema.StringAttribute{
                MarkdownDescription: "The last time a public note reminder was sent.",
                Computed: true,
            },
            "breach_notification_sent_at": schema.StringAttribute{
                MarkdownDescription: "The time when breach notification was sent to incident owners.",
                Computed: true,
            },
            "sla_started_at": schema.StringAttribute{
                MarkdownDescription: "The time when SLA tracking started (usually the incident declaredAt time).",
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

func (d *IncidentSlaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentSlaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentSlaDataSourceModel

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
    if !data.IncidentSlaRuleId.IsNull() && !data.IncidentSlaRuleId.IsUnknown() {
        filters["incidentSlaRuleId"] = data.IncidentSlaRuleId.ValueString()
        filterNames = append(filterNames, "incident_sla_rule_id = "+fmt.Sprintf("%q", data.IncidentSlaRuleId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident sla up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident sla up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incidentId": true,
        "incidentSlaRuleId": true,
        "responseDeadline": true,
        "resolutionDeadline": true,
        "status": true,
        "respondedAt": true,
        "resolvedAt": true,
        "lastInternalNoteReminderSentAt": true,
        "lastPublicNoteReminderSentAt": true,
        "breachNotificationSentAt": true,
        "slaStartedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-sla/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_sla, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident sla found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_sla: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-sla/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_sla, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_sla: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident sla matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident sla matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_sla.")
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
    if obj, ok := item["incidentSlaRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSlaRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSlaRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSlaRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSlaRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSlaRuleId = types.StringNull()
        }
    } else if val, ok := item["incidentSlaRuleId"].(string); ok {
        data.IncidentSlaRuleId = types.StringValue(val)
    } else {
        data.IncidentSlaRuleId = types.StringNull()
    }
    if obj, ok := item["responseDeadline"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResponseDeadline = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResponseDeadline = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResponseDeadline = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResponseDeadline = types.StringValue(string(jsonBytes))
        } else {
            data.ResponseDeadline = types.StringNull()
        }
    } else if val, ok := item["responseDeadline"].(string); ok {
        data.ResponseDeadline = types.StringValue(val)
    } else {
        data.ResponseDeadline = types.StringNull()
    }
    if obj, ok := item["resolutionDeadline"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolutionDeadline = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResolutionDeadline = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResolutionDeadline = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResolutionDeadline = types.StringValue(string(jsonBytes))
        } else {
            data.ResolutionDeadline = types.StringNull()
        }
    } else if val, ok := item["resolutionDeadline"].(string); ok {
        data.ResolutionDeadline = types.StringValue(val)
    } else {
        data.ResolutionDeadline = types.StringNull()
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
    if obj, ok := item["respondedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RespondedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RespondedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RespondedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RespondedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RespondedAt = types.StringNull()
        }
    } else if val, ok := item["respondedAt"].(string); ok {
        data.RespondedAt = types.StringValue(val)
    } else {
        data.RespondedAt = types.StringNull()
    }
    if obj, ok := item["resolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedAt = types.StringNull()
        }
    } else if val, ok := item["resolvedAt"].(string); ok {
        data.ResolvedAt = types.StringValue(val)
    } else {
        data.ResolvedAt = types.StringNull()
    }
    if obj, ok := item["lastInternalNoteReminderSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastInternalNoteReminderSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastInternalNoteReminderSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastInternalNoteReminderSentAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastInternalNoteReminderSentAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastInternalNoteReminderSentAt = types.StringNull()
        }
    } else if val, ok := item["lastInternalNoteReminderSentAt"].(string); ok {
        data.LastInternalNoteReminderSentAt = types.StringValue(val)
    } else {
        data.LastInternalNoteReminderSentAt = types.StringNull()
    }
    if obj, ok := item["lastPublicNoteReminderSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastPublicNoteReminderSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastPublicNoteReminderSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastPublicNoteReminderSentAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastPublicNoteReminderSentAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastPublicNoteReminderSentAt = types.StringNull()
        }
    } else if val, ok := item["lastPublicNoteReminderSentAt"].(string); ok {
        data.LastPublicNoteReminderSentAt = types.StringValue(val)
    } else {
        data.LastPublicNoteReminderSentAt = types.StringNull()
    }
    if obj, ok := item["breachNotificationSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BreachNotificationSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BreachNotificationSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BreachNotificationSentAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BreachNotificationSentAt = types.StringValue(string(jsonBytes))
        } else {
            data.BreachNotificationSentAt = types.StringNull()
        }
    } else if val, ok := item["breachNotificationSentAt"].(string); ok {
        data.BreachNotificationSentAt = types.StringValue(val)
    } else {
        data.BreachNotificationSentAt = types.StringNull()
    }
    if obj, ok := item["slaStartedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlaStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SlaStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SlaStartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SlaStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.SlaStartedAt = types.StringNull()
        }
    } else if val, ok := item["slaStartedAt"].(string); ok {
        data.SlaStartedAt = types.StringValue(val)
    } else {
        data.SlaStartedAt = types.StringNull()
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
