package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SloBurnRateRuleDataSource{}

func NewSloBurnRateRuleDataSource() datasource.DataSource {
    return &SloBurnRateRuleDataSource{}
}

// SloBurnRateRuleDataSource defines the data source implementation.
type SloBurnRateRuleDataSource struct {
    client *Client
}

// SloBurnRateRuleDataSourceModel describes the data source data model.
type SloBurnRateRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceLevelObjectiveId types.String `tfsdk:"service_level_objective_id"`
    Name types.String `tfsdk:"name"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    BurnRateThreshold types.Number `tfsdk:"burn_rate_threshold"`
    LongWindowInMinutes types.Number `tfsdk:"long_window_in_minutes"`
    ShortWindowInMinutes types.Number `tfsdk:"short_window_in_minutes"`
    MinimumSampleCount types.Number `tfsdk:"minimum_sample_count"`
    RefireSuppressionMinutes types.Number `tfsdk:"refire_suppression_minutes"`
    ShouldCreateAlert types.Bool `tfsdk:"should_create_alert"`
    ShouldCreateIncident types.Bool `tfsdk:"should_create_incident"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    AlertTitleTemplate types.String `tfsdk:"alert_title_template"`
    AlertDescriptionTemplate types.String `tfsdk:"alert_description_template"`
    AlertRemediationNotes types.String `tfsdk:"alert_remediation_notes"`
    IsAlertPrivate types.Bool `tfsdk:"is_alert_private"`
    AutoResolveAlert types.Bool `tfsdk:"auto_resolve_alert"`
    AlertLabels types.Set `tfsdk:"alert_labels"`
    AlertOwnerTeams types.Set `tfsdk:"alert_owner_teams"`
    AlertOwnerUsers types.Set `tfsdk:"alert_owner_users"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    IncidentOnCallDutyPolicies types.Set `tfsdk:"incident_on_call_duty_policies"`
    IncidentTitleTemplate types.String `tfsdk:"incident_title_template"`
    IncidentDescriptionTemplate types.String `tfsdk:"incident_description_template"`
    IncidentRemediationNotes types.String `tfsdk:"incident_remediation_notes"`
    IsIncidentPrivate types.Bool `tfsdk:"is_incident_private"`
    AutoResolveIncident types.Bool `tfsdk:"auto_resolve_incident"`
    IncidentLabels types.Set `tfsdk:"incident_labels"`
    IncidentOwnerTeams types.Set `tfsdk:"incident_owner_teams"`
    IncidentOwnerUsers types.Set `tfsdk:"incident_owner_users"`
    AddSloOwnersAsOwners types.Bool `tfsdk:"add_slo_owners_as_owners"`
    LastAlertCreatedAt types.String `tfsdk:"last_alert_created_at"`
    LastAlertResolvedAt types.String `tfsdk:"last_alert_resolved_at"`
    LastIncidentCreatedAt types.String `tfsdk:"last_incident_created_at"`
    LastIncidentResolvedAt types.String `tfsdk:"last_incident_resolved_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SloBurnRateRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_slo_burn_rate_rule"
}

func (d *SloBurnRateRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure multi-window burn rate rules that raise alerts and/or declare incidents when a Service Level Objective consumes its error budget too quickly Look up an existing slo burn rate rule by `id`, or by any of its other arguments (`name`, `add_slo_owners_as_owners`, `alert_description_template`, ...): each one set must match, and exactly one slo burn rate rule may match them all.",

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
            "service_level_objective_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Service Level Objective this burn rate rule belongs to. The ID of a `oneuptime_service_level_objective`.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this burn rate rule.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this burn rate rule is enabled.",
                Optional: true,
                Computed: true,
            },
            "burn_rate_threshold": schema.NumberAttribute{
                MarkdownDescription: "Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4).",
                Optional: true,
                Computed: true,
            },
            "long_window_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it.",
                Optional: true,
                Computed: true,
            },
            "short_window_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped.",
                Optional: true,
                Computed: true,
            },
            "minimum_sample_count": schema.NumberAttribute{
                MarkdownDescription: "For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic.",
                Optional: true,
                Computed: true,
            },
            "refire_suppression_minutes": schema.NumberAttribute{
                MarkdownDescription: "Minimum number of minutes after an alert or incident resolves before this rule can declare that same record again. Each output is suppressed independently, from its own resolve. Defaults to the long window length when not set.",
                Optional: true,
                Computed: true,
            },
            "should_create_alert": schema.BoolAttribute{
                MarkdownDescription: "Raise an Alert when this burn rate rule fires. Enabled by default.",
                Optional: true,
                Computed: true,
            },
            "should_create_incident": schema.BoolAttribute{
                MarkdownDescription: "Declare an Incident when this burn rate rule fires. Disabled by default.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Alert Severity of the alert created when this burn rate rule fires. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "On-call duty policies attached to alerts created by this burn rate rule. Incidents have their own list. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_title_template": schema.StringAttribute{
                MarkdownDescription: "Title of the alert raised when this burn rate rule fires. Supports template variables such as {{sloName}}. Leave empty to use the default title.",
                Optional: true,
                Computed: true,
            },
            "alert_description_template": schema.StringAttribute{
                MarkdownDescription: "Description (in Markdown) of the alert raised when this burn rate rule fires. Supports template variables. Leave empty to use the default description.",
                Optional: true,
                Computed: true,
            },
            "alert_remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Remediation notes (in Markdown) attached to the alert raised when this burn rate rule fires. Supports template variables.",
                Optional: true,
                Computed: true,
            },
            "is_alert_private": schema.BoolAttribute{
                MarkdownDescription: "Make the alert raised by this burn rate rule private, so only its owners, project admins and project owners can see it. Disabled by default.",
                Optional: true,
                Computed: true,
            },
            "auto_resolve_alert": schema.BoolAttribute{
                MarkdownDescription: "Resolve the alert automatically when the burn rate over the long window drops back below the threshold. Enabled by default. When disabled, the alert stays open until someone resolves it.",
                Optional: true,
                Computed: true,
            },
            "alert_labels": schema.SetAttribute{
                MarkdownDescription: "Labels added to alerts raised by this burn rate rule. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams added as owners of alerts raised by this burn rate rule. IDs of `oneuptime_team` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_owner_users": schema.SetAttribute{
                MarkdownDescription: "Users added as owners of alerts raised by this burn rate rule. IDs of `oneuptime_user` records.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Incident Severity of the incident declared when this burn rate rule fires. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "incident_on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "On-call duty policies attached to incidents declared by this burn rate rule. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_template": schema.StringAttribute{
                MarkdownDescription: "Title of the incident declared when this burn rate rule fires. Supports template variables such as {{sloName}}. Leave empty to use the default title.",
                Optional: true,
                Computed: true,
            },
            "incident_description_template": schema.StringAttribute{
                MarkdownDescription: "Description (in Markdown) of the incident declared when this burn rate rule fires. Supports template variables. Leave empty to use the default description.",
                Optional: true,
                Computed: true,
            },
            "incident_remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Remediation notes (in Markdown) attached to the incident declared when this burn rate rule fires. Supports template variables.",
                Optional: true,
                Computed: true,
            },
            "is_incident_private": schema.BoolAttribute{
                MarkdownDescription: "Make the incident declared by this burn rate rule private, so only its owners, project admins and project owners can see it. Disabled by default.",
                Optional: true,
                Computed: true,
            },
            "auto_resolve_incident": schema.BoolAttribute{
                MarkdownDescription: "Resolve the incident automatically when the burn rate over the long window drops back below the threshold. Enabled by default. When disabled, the incident stays open until someone resolves it.",
                Optional: true,
                Computed: true,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Labels added to incidents declared by this burn rate rule. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams added as owners of incidents declared by this burn rate rule. IDs of `oneuptime_team` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_owner_users": schema.SetAttribute{
                MarkdownDescription: "Users added as owners of incidents declared by this burn rate rule. IDs of `oneuptime_user` records.",
                Computed: true,
                ElementType: types.StringType,
            },
            "add_slo_owners_as_owners": schema.BoolAttribute{
                MarkdownDescription: "Also add the owner users and owner teams of the Service Level Objective as owners of the alerts and incidents this burn rate rule creates. Disabled by default.",
                Optional: true,
                Computed: true,
            },
            "last_alert_created_at": schema.StringAttribute{
                MarkdownDescription: "The last time an alert was created by this burn rate rule. Computed by the worker.",
                Computed: true,
            },
            "last_alert_resolved_at": schema.StringAttribute{
                MarkdownDescription: "The last time an alert created by this burn rate rule was resolved. Computed by the worker.",
                Computed: true,
            },
            "last_incident_created_at": schema.StringAttribute{
                MarkdownDescription: "The last time an incident was declared by this burn rate rule. Computed by the worker.",
                Computed: true,
            },
            "last_incident_resolved_at": schema.StringAttribute{
                MarkdownDescription: "The last time an incident declared by this burn rate rule was resolved. Computed by the worker.",
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

func (d *SloBurnRateRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SloBurnRateRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SloBurnRateRuleDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ServiceLevelObjectiveId.IsNull() && !data.ServiceLevelObjectiveId.IsUnknown() {
        filters["serviceLevelObjectiveId"] = data.ServiceLevelObjectiveId.ValueString()
        filterNames = append(filterNames, "service_level_objective_id = "+fmt.Sprintf("%q", data.ServiceLevelObjectiveId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.BurnRateThreshold.IsNull() && !data.BurnRateThreshold.IsUnknown() {
        filters["burnRateThreshold"] = lookupNumber(data.BurnRateThreshold)
        filterNames = append(filterNames, "burn_rate_threshold = "+data.BurnRateThreshold.ValueBigFloat().String())
    }
    if !data.LongWindowInMinutes.IsNull() && !data.LongWindowInMinutes.IsUnknown() {
        filters["longWindowInMinutes"] = lookupNumber(data.LongWindowInMinutes)
        filterNames = append(filterNames, "long_window_in_minutes = "+data.LongWindowInMinutes.ValueBigFloat().String())
    }
    if !data.ShortWindowInMinutes.IsNull() && !data.ShortWindowInMinutes.IsUnknown() {
        filters["shortWindowInMinutes"] = lookupNumber(data.ShortWindowInMinutes)
        filterNames = append(filterNames, "short_window_in_minutes = "+data.ShortWindowInMinutes.ValueBigFloat().String())
    }
    if !data.MinimumSampleCount.IsNull() && !data.MinimumSampleCount.IsUnknown() {
        filters["minimumSampleCount"] = lookupNumber(data.MinimumSampleCount)
        filterNames = append(filterNames, "minimum_sample_count = "+data.MinimumSampleCount.ValueBigFloat().String())
    }
    if !data.RefireSuppressionMinutes.IsNull() && !data.RefireSuppressionMinutes.IsUnknown() {
        filters["refireSuppressionMinutes"] = lookupNumber(data.RefireSuppressionMinutes)
        filterNames = append(filterNames, "refire_suppression_minutes = "+data.RefireSuppressionMinutes.ValueBigFloat().String())
    }
    if !data.ShouldCreateAlert.IsNull() && !data.ShouldCreateAlert.IsUnknown() {
        filters["shouldCreateAlert"] = data.ShouldCreateAlert.ValueBool()
        filterNames = append(filterNames, "should_create_alert = "+fmt.Sprintf("%t", data.ShouldCreateAlert.ValueBool()))
    }
    if !data.ShouldCreateIncident.IsNull() && !data.ShouldCreateIncident.IsUnknown() {
        filters["shouldCreateIncident"] = data.ShouldCreateIncident.ValueBool()
        filterNames = append(filterNames, "should_create_incident = "+fmt.Sprintf("%t", data.ShouldCreateIncident.ValueBool()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.AlertTitleTemplate.IsNull() && !data.AlertTitleTemplate.IsUnknown() {
        filters["alertTitleTemplate"] = data.AlertTitleTemplate.ValueString()
        filterNames = append(filterNames, "alert_title_template = "+fmt.Sprintf("%q", data.AlertTitleTemplate.ValueString()))
    }
    if !data.AlertDescriptionTemplate.IsNull() && !data.AlertDescriptionTemplate.IsUnknown() {
        filters["alertDescriptionTemplate"] = data.AlertDescriptionTemplate.ValueString()
        filterNames = append(filterNames, "alert_description_template = "+fmt.Sprintf("%q", data.AlertDescriptionTemplate.ValueString()))
    }
    if !data.AlertRemediationNotes.IsNull() && !data.AlertRemediationNotes.IsUnknown() {
        filters["alertRemediationNotes"] = data.AlertRemediationNotes.ValueString()
        filterNames = append(filterNames, "alert_remediation_notes = "+fmt.Sprintf("%q", data.AlertRemediationNotes.ValueString()))
    }
    if !data.IsAlertPrivate.IsNull() && !data.IsAlertPrivate.IsUnknown() {
        filters["isAlertPrivate"] = data.IsAlertPrivate.ValueBool()
        filterNames = append(filterNames, "is_alert_private = "+fmt.Sprintf("%t", data.IsAlertPrivate.ValueBool()))
    }
    if !data.AutoResolveAlert.IsNull() && !data.AutoResolveAlert.IsUnknown() {
        filters["autoResolveAlert"] = data.AutoResolveAlert.ValueBool()
        filterNames = append(filterNames, "auto_resolve_alert = "+fmt.Sprintf("%t", data.AutoResolveAlert.ValueBool()))
    }
    if !data.IncidentSeverityId.IsNull() && !data.IncidentSeverityId.IsUnknown() {
        filters["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
        filterNames = append(filterNames, "incident_severity_id = "+fmt.Sprintf("%q", data.IncidentSeverityId.ValueString()))
    }
    if !data.IncidentTitleTemplate.IsNull() && !data.IncidentTitleTemplate.IsUnknown() {
        filters["incidentTitleTemplate"] = data.IncidentTitleTemplate.ValueString()
        filterNames = append(filterNames, "incident_title_template = "+fmt.Sprintf("%q", data.IncidentTitleTemplate.ValueString()))
    }
    if !data.IncidentDescriptionTemplate.IsNull() && !data.IncidentDescriptionTemplate.IsUnknown() {
        filters["incidentDescriptionTemplate"] = data.IncidentDescriptionTemplate.ValueString()
        filterNames = append(filterNames, "incident_description_template = "+fmt.Sprintf("%q", data.IncidentDescriptionTemplate.ValueString()))
    }
    if !data.IncidentRemediationNotes.IsNull() && !data.IncidentRemediationNotes.IsUnknown() {
        filters["incidentRemediationNotes"] = data.IncidentRemediationNotes.ValueString()
        filterNames = append(filterNames, "incident_remediation_notes = "+fmt.Sprintf("%q", data.IncidentRemediationNotes.ValueString()))
    }
    if !data.IsIncidentPrivate.IsNull() && !data.IsIncidentPrivate.IsUnknown() {
        filters["isIncidentPrivate"] = data.IsIncidentPrivate.ValueBool()
        filterNames = append(filterNames, "is_incident_private = "+fmt.Sprintf("%t", data.IsIncidentPrivate.ValueBool()))
    }
    if !data.AutoResolveIncident.IsNull() && !data.AutoResolveIncident.IsUnknown() {
        filters["autoResolveIncident"] = data.AutoResolveIncident.ValueBool()
        filterNames = append(filterNames, "auto_resolve_incident = "+fmt.Sprintf("%t", data.AutoResolveIncident.ValueBool()))
    }
    if !data.AddSloOwnersAsOwners.IsNull() && !data.AddSloOwnersAsOwners.IsUnknown() {
        filters["addSloOwnersAsOwners"] = data.AddSloOwnersAsOwners.ValueBool()
        filterNames = append(filterNames, "add_slo_owners_as_owners = "+fmt.Sprintf("%t", data.AddSloOwnersAsOwners.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the slo burn rate rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the slo burn rate rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "serviceLevelObjectiveId": true,
        "name": true,
        "isEnabled": true,
        "burnRateThreshold": true,
        "longWindowInMinutes": true,
        "shortWindowInMinutes": true,
        "minimumSampleCount": true,
        "refireSuppressionMinutes": true,
        "shouldCreateAlert": true,
        "shouldCreateIncident": true,
        "alertSeverityId": true,
        "onCallDutyPolicies": true,
        "alertTitleTemplate": true,
        "alertDescriptionTemplate": true,
        "alertRemediationNotes": true,
        "isAlertPrivate": true,
        "autoResolveAlert": true,
        "alertLabels": true,
        "alertOwnerTeams": true,
        "alertOwnerUsers": true,
        "incidentSeverityId": true,
        "incidentOnCallDutyPolicies": true,
        "incidentTitleTemplate": true,
        "incidentDescriptionTemplate": true,
        "incidentRemediationNotes": true,
        "isIncidentPrivate": true,
        "autoResolveIncident": true,
        "incidentLabels": true,
        "incidentOwnerTeams": true,
        "incidentOwnerUsers": true,
        "addSloOwnersAsOwners": true,
        "lastAlertCreatedAt": true,
        "lastAlertResolvedAt": true,
        "lastIncidentCreatedAt": true,
        "lastIncidentResolvedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/service-level-objective-burn-rate-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read slo_burn_rate_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No slo burn rate rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read slo_burn_rate_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/service-level-objective-burn-rate-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list slo_burn_rate_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list slo_burn_rate_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No slo burn rate rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one slo burn rate rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for slo_burn_rate_rule.")
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
    if obj, ok := item["serviceLevelObjectiveId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceLevelObjectiveId = types.StringNull()
        }
    } else if val, ok := item["serviceLevelObjectiveId"].(string); ok {
        data.ServiceLevelObjectiveId = types.StringValue(val)
    } else {
        data.ServiceLevelObjectiveId = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["burnRateThreshold"].(float64); ok {
        data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["burnRateThreshold"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.BurnRateThreshold = types.NumberValue(big.NewFloat(val))
        } else {
            data.BurnRateThreshold = types.NumberNull()
        }
    } else {
        data.BurnRateThreshold = types.NumberNull()
    }
    if val, ok := item["longWindowInMinutes"].(float64); ok {
        data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["longWindowInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LongWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.LongWindowInMinutes = types.NumberNull()
        }
    } else {
        data.LongWindowInMinutes = types.NumberNull()
    }
    if val, ok := item["shortWindowInMinutes"].(float64); ok {
        data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["shortWindowInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShortWindowInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShortWindowInMinutes = types.NumberNull()
        }
    } else {
        data.ShortWindowInMinutes = types.NumberNull()
    }
    if val, ok := item["minimumSampleCount"].(float64); ok {
        data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["minimumSampleCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MinimumSampleCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumSampleCount = types.NumberNull()
        }
    } else {
        data.MinimumSampleCount = types.NumberNull()
    }
    if val, ok := item["refireSuppressionMinutes"].(float64); ok {
        data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["refireSuppressionMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RefireSuppressionMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.RefireSuppressionMinutes = types.NumberNull()
        }
    } else {
        data.RefireSuppressionMinutes = types.NumberNull()
    }
    if val, ok := item["shouldCreateAlert"].(bool); ok {
        data.ShouldCreateAlert = types.BoolValue(val)
    } else {
        data.ShouldCreateAlert = types.BoolNull()
    }
    if val, ok := item["shouldCreateIncident"].(bool); ok {
        data.ShouldCreateIncident = types.BoolValue(val)
    } else {
        data.ShouldCreateIncident = types.BoolNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := item["onCallDutyPolicies"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OnCallDutyPolicies = types.SetNull(types.StringType)
    }
    if obj, ok := item["alertTitleTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertTitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertTitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.AlertTitleTemplate = types.StringNull()
        }
    } else if val, ok := item["alertTitleTemplate"].(string); ok {
        data.AlertTitleTemplate = types.StringValue(val)
    } else {
        data.AlertTitleTemplate = types.StringNull()
    }
    if obj, ok := item["alertDescriptionTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertDescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.AlertDescriptionTemplate = types.StringNull()
        }
    } else if val, ok := item["alertDescriptionTemplate"].(string); ok {
        data.AlertDescriptionTemplate = types.StringValue(val)
    } else {
        data.AlertDescriptionTemplate = types.StringNull()
    }
    if obj, ok := item["alertRemediationNotes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertRemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertRemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertRemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertRemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.AlertRemediationNotes = types.StringNull()
        }
    } else if val, ok := item["alertRemediationNotes"].(string); ok {
        data.AlertRemediationNotes = types.StringValue(val)
    } else {
        data.AlertRemediationNotes = types.StringNull()
    }
    if val, ok := item["isAlertPrivate"].(bool); ok {
        data.IsAlertPrivate = types.BoolValue(val)
    } else {
        data.IsAlertPrivate = types.BoolNull()
    }
    if val, ok := item["autoResolveAlert"].(bool); ok {
        data.AutoResolveAlert = types.BoolValue(val)
    } else {
        data.AutoResolveAlert = types.BoolNull()
    }
    if val, ok := item["alertLabels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.AlertLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["alertOwnerTeams"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.AlertOwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertOwnerTeams = types.SetNull(types.StringType)
    }
    if val, ok := item["alertOwnerUsers"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.AlertOwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertOwnerUsers = types.SetNull(types.StringType)
    }
    if obj, ok := item["incidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := item["incidentOnCallDutyPolicies"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.IncidentOnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentOnCallDutyPolicies = types.SetNull(types.StringType)
    }
    if obj, ok := item["incidentTitleTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentTitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentTitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentTitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentTitleTemplate = types.StringNull()
        }
    } else if val, ok := item["incidentTitleTemplate"].(string); ok {
        data.IncidentTitleTemplate = types.StringValue(val)
    } else {
        data.IncidentTitleTemplate = types.StringNull()
    }
    if obj, ok := item["incidentDescriptionTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentDescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentDescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentDescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentDescriptionTemplate = types.StringNull()
        }
    } else if val, ok := item["incidentDescriptionTemplate"].(string); ok {
        data.IncidentDescriptionTemplate = types.StringValue(val)
    } else {
        data.IncidentDescriptionTemplate = types.StringNull()
    }
    if obj, ok := item["incidentRemediationNotes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentRemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentRemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentRemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentRemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentRemediationNotes = types.StringNull()
        }
    } else if val, ok := item["incidentRemediationNotes"].(string); ok {
        data.IncidentRemediationNotes = types.StringValue(val)
    } else {
        data.IncidentRemediationNotes = types.StringNull()
    }
    if val, ok := item["isIncidentPrivate"].(bool); ok {
        data.IsIncidentPrivate = types.BoolValue(val)
    } else {
        data.IsIncidentPrivate = types.BoolNull()
    }
    if val, ok := item["autoResolveIncident"].(bool); ok {
        data.AutoResolveIncident = types.BoolValue(val)
    } else {
        data.AutoResolveIncident = types.BoolNull()
    }
    if val, ok := item["incidentLabels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.IncidentLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentOwnerTeams"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.IncidentOwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentOwnerTeams = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentOwnerUsers"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.IncidentOwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentOwnerUsers = types.SetNull(types.StringType)
    }
    if val, ok := item["addSloOwnersAsOwners"].(bool); ok {
        data.AddSloOwnersAsOwners = types.BoolValue(val)
    } else {
        data.AddSloOwnersAsOwners = types.BoolNull()
    }
    if obj, ok := item["lastAlertCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAlertCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAlertCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAlertCreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAlertCreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastAlertCreatedAt = types.StringNull()
        }
    } else if val, ok := item["lastAlertCreatedAt"].(string); ok {
        data.LastAlertCreatedAt = types.StringValue(val)
    } else {
        data.LastAlertCreatedAt = types.StringNull()
    }
    if obj, ok := item["lastAlertResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAlertResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAlertResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAlertResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAlertResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastAlertResolvedAt = types.StringNull()
        }
    } else if val, ok := item["lastAlertResolvedAt"].(string); ok {
        data.LastAlertResolvedAt = types.StringValue(val)
    } else {
        data.LastAlertResolvedAt = types.StringNull()
    }
    if obj, ok := item["lastIncidentCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastIncidentCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastIncidentCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastIncidentCreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastIncidentCreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastIncidentCreatedAt = types.StringNull()
        }
    } else if val, ok := item["lastIncidentCreatedAt"].(string); ok {
        data.LastIncidentCreatedAt = types.StringValue(val)
    } else {
        data.LastIncidentCreatedAt = types.StringNull()
    }
    if obj, ok := item["lastIncidentResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastIncidentResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastIncidentResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastIncidentResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastIncidentResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastIncidentResolvedAt = types.StringNull()
        }
    } else if val, ok := item["lastIncidentResolvedAt"].(string); ok {
        data.LastIncidentResolvedAt = types.StringValue(val)
    } else {
        data.LastIncidentResolvedAt = types.StringNull()
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
