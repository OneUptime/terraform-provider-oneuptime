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
var _ datasource.DataSource = &AlertDataSource{}

func NewAlertDataSource() datasource.DataSource {
    return &AlertDataSource{}
}

// AlertDataSource defines the data source implementation.
type AlertDataSource struct {
    client *Client
}

// AlertDataSourceModel describes the data source data model.
type AlertDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    ImpactStartedAt types.String `tfsdk:"impact_started_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    Hosts types.Set `tfsdk:"hosts"`
    KubernetesClusters types.Set `tfsdk:"kubernetes_clusters"`
    KubernetesResources types.Set `tfsdk:"kubernetes_resources"`
    KubernetesContainers types.Set `tfsdk:"kubernetes_containers"`
    DockerHosts types.Set `tfsdk:"docker_hosts"`
    PodmanHosts types.Set `tfsdk:"podman_hosts"`
    ProxmoxClusters types.Set `tfsdk:"proxmox_clusters"`
    VmwareVCenters types.Set `tfsdk:"vmware_v_centers"`
    IotFleets types.Set `tfsdk:"iot_fleets"`
    DockerSwarmClusters types.Set `tfsdk:"docker_swarm_clusters"`
    CephClusters types.Set `tfsdk:"ceph_clusters"`
    StorageArrays types.Set `tfsdk:"storage_arrays"`
    DatabaseServers types.Set `tfsdk:"database_servers"`
    DockerResources types.Set `tfsdk:"docker_resources"`
    PodmanResources types.Set `tfsdk:"podman_resources"`
    Services types.Set `tfsdk:"services"`
    ServiceLevelObjectives types.Set `tfsdk:"service_level_objectives"`
    Labels types.Set `tfsdk:"labels"`
    CurrentAlertStateId types.String `tfsdk:"current_alert_state_id"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    MonitorStatusWhenThisAlertWasCreatedId types.String `tfsdk:"monitor_status_when_this_alert_was_created_id"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfAlertCreation types.Bool `tfsdk:"is_owner_notified_of_alert_creation"`
    RootCause types.String `tfsdk:"root_cause"`
    CreatedStateLog types.String `tfsdk:"created_state_log"`
    CreatedCriteriaId types.String `tfsdk:"created_criteria_id"`
    SeriesFingerprint types.String `tfsdk:"series_fingerprint"`
    SeriesLabels types.String `tfsdk:"series_labels"`
    MonitorSummary types.String `tfsdk:"monitor_summary"`
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.String `tfsdk:"telemetry_query"`
    AlertNumber types.Number `tfsdk:"alert_number"`
    AlertNumberWithPrefix types.String `tfsdk:"alert_number_with_prefix"`
    AlertEpisodeId types.String `tfsdk:"alert_episode_id"`
    IsPrivate types.Bool `tfsdk:"is_private"`
    EnableReminders types.Bool `tfsdk:"enable_reminders"`
    NextReminderNotificationAt types.String `tfsdk:"next_reminder_notification_at"`
    ReminderNotificationSentCount types.Number `tfsdk:"reminder_notification_sent_count"`
}

func (d *AlertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert"
}

func (d *AlertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage alerts for your project Look up an existing alert by `id`, or by any of its other arguments (`alert_episode_id`, `alert_number`, `alert_number_with_prefix`, ...): each one set must match, and exactly one alert may match them all.",

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
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this alert.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Short description of this alert. This will be visible on the status page. This is in markdown.",
                Optional: true,
                Computed: true,
            },
            "impact_started_at": schema.StringAttribute{
                MarkdownDescription: "When customer impact actually began. Left blank until someone records it - never inferred, because a guessed value is worse than no value.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the monitor this alert belongs to. The ID of a `oneuptime_monitor`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies affected by this alert. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "hosts": schema.SetAttribute{
                MarkdownDescription: "List of hosts affected by this alert. IDs of `oneuptime_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes clusters affected by this alert. IDs of `oneuptime_kubernetes_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_resources": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this alert.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_containers": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes containers affected by this alert.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Docker hosts affected by this alert. IDs of `oneuptime_docker_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Podman hosts affected by this alert. IDs of `oneuptime_podman_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "proxmox_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Proxmox clusters affected by this alert. IDs of `oneuptime_proxmox_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "vmware_v_centers": schema.SetAttribute{
                MarkdownDescription: "List of vCenters affected by this alert. IDs of `oneuptime_vcenter` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "iot_fleets": schema.SetAttribute{
                MarkdownDescription: "List of IoT fleets affected by this alert. IDs of `oneuptime_iot_fleet` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_swarm_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Docker Swarm clusters affected by this alert. IDs of `oneuptime_docker_swarm_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "ceph_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Ceph clusters affected by this alert. IDs of `oneuptime_ceph_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "storage_arrays": schema.SetAttribute{
                MarkdownDescription: "List of storage arrays affected by this alert. IDs of `oneuptime_storage_array` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "database_servers": schema.SetAttribute{
                MarkdownDescription: "List of databases affected by this alert. IDs of `oneuptime_database` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_resources": schema.SetAttribute{
                MarkdownDescription: "List of Docker resources (containers, images, networks, volumes) affected by this alert.",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_resources": schema.SetAttribute{
                MarkdownDescription: "List of Podman resources (containers, images, networks, volumes) affected by this alert.",
                Computed: true,
                ElementType: types.StringType,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "List of services affected by this alert. IDs of `oneuptime_service` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_level_objectives": schema.SetAttribute{
                MarkdownDescription: "List of Service Level Objectives (SLOs) affected by this alert. IDs of `oneuptime_service_level_objective` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_alert_state_id": schema.StringAttribute{
                MarkdownDescription: "Current Alert State ID. The ID of a `oneuptime_alert_state`.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "Alert Severity ID. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "monitor_status_when_this_alert_was_created_id": schema.StringAttribute{
                MarkdownDescription: "Monitor Status ID when this alert was created. The ID of a `oneuptime_monitor_status`.",
                Optional: true,
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_owner_notified_of_alert_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this alert is created?",
                Optional: true,
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "What is the root cause of this alert?",
                Optional: true,
                Computed: true,
            },
            "created_state_log": schema.StringAttribute{
                MarkdownDescription: "A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "If this alert was created by a Probe, this is the ID of the criteria that created it.",
                Optional: true,
                Computed: true,
            },
            "series_fingerprint": schema.StringAttribute{
                MarkdownDescription: "For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one alert is created per affected series.",
                Optional: true,
                Computed: true,
            },
            "series_labels": schema.StringAttribute{
                MarkdownDescription: "Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this alert was created from a per-series metric breach. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "monitor_summary": schema.StringAttribute{
                MarkdownDescription: "The monitor summary captured at the moment this alert was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "If this alert was created by a Probe, this is the ID of the probe that created it. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Is this alert created by OneUptime Probe or Workers automatically (and not created manually by a user)?",
                Optional: true,
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Notes on how to remediate this alert. This is in markdown.",
                Optional: true,
                Computed: true,
            },
            "telemetry_query": schema.StringAttribute{
                MarkdownDescription: "Telemetry query for this alert. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "alert_number": schema.NumberAttribute{
                MarkdownDescription: "Alert Number.",
                Optional: true,
                Computed: true,
            },
            "alert_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Alert number with prefix (e.g., 'ALT-42' or '#42').",
                Optional: true,
                Computed: true,
            },
            "alert_episode_id": schema.StringAttribute{
                MarkdownDescription: "The ID of the latest episode this alert is a member of, if any. Read-only: set by OneUptime when the alert is added to or removed from an episode's members (Alert Episode Member). The ID of a `oneuptime_alert_episode`.",
                Optional: true,
                Computed: true,
            },
            "is_private": schema.BoolAttribute{
                MarkdownDescription: "If true, this alert is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners.",
                Optional: true,
                Computed: true,
            },
            "enable_reminders": schema.BoolAttribute{
                MarkdownDescription: "Should reminder notifications be sent to owners while this alert is still open? Reminders are sent based on the reminder rules configured for this project.",
                Optional: true,
                Computed: true,
            },
            "next_reminder_notification_at": schema.StringAttribute{
                MarkdownDescription: "When will the next reminder notification be sent to owners of this alert? This is set automatically based on the reminder rules configured for this project.",
                Computed: true,
            },
            "reminder_notification_sent_count": schema.NumberAttribute{
                MarkdownDescription: "How many reminder notifications have been sent to owners of this alert so far.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *AlertDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.CurrentAlertStateId.IsNull() && !data.CurrentAlertStateId.IsUnknown() {
        filters["currentAlertStateId"] = data.CurrentAlertStateId.ValueString()
        filterNames = append(filterNames, "current_alert_state_id = "+fmt.Sprintf("%q", data.CurrentAlertStateId.ValueString()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.MonitorStatusWhenThisAlertWasCreatedId.IsNull() && !data.MonitorStatusWhenThisAlertWasCreatedId.IsUnknown() {
        filters["monitorStatusWhenThisAlertWasCreatedId"] = data.MonitorStatusWhenThisAlertWasCreatedId.ValueString()
        filterNames = append(filterNames, "monitor_status_when_this_alert_was_created_id = "+fmt.Sprintf("%q", data.MonitorStatusWhenThisAlertWasCreatedId.ValueString()))
    }
    if !data.IsOwnerNotifiedOfAlertCreation.IsNull() && !data.IsOwnerNotifiedOfAlertCreation.IsUnknown() {
        filters["isOwnerNotifiedOfAlertCreation"] = data.IsOwnerNotifiedOfAlertCreation.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_alert_creation = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfAlertCreation.ValueBool()))
    }
    if !data.RootCause.IsNull() && !data.RootCause.IsUnknown() {
        filters["rootCause"] = data.RootCause.ValueString()
        filterNames = append(filterNames, "root_cause = "+fmt.Sprintf("%q", data.RootCause.ValueString()))
    }
    if !data.CreatedCriteriaId.IsNull() && !data.CreatedCriteriaId.IsUnknown() {
        filters["createdCriteriaId"] = data.CreatedCriteriaId.ValueString()
        filterNames = append(filterNames, "created_criteria_id = "+fmt.Sprintf("%q", data.CreatedCriteriaId.ValueString()))
    }
    if !data.SeriesFingerprint.IsNull() && !data.SeriesFingerprint.IsUnknown() {
        filters["seriesFingerprint"] = data.SeriesFingerprint.ValueString()
        filterNames = append(filterNames, "series_fingerprint = "+fmt.Sprintf("%q", data.SeriesFingerprint.ValueString()))
    }
    if !data.CreatedByProbeId.IsNull() && !data.CreatedByProbeId.IsUnknown() {
        filters["createdByProbeId"] = data.CreatedByProbeId.ValueString()
        filterNames = append(filterNames, "created_by_probe_id = "+fmt.Sprintf("%q", data.CreatedByProbeId.ValueString()))
    }
    if !data.IsCreatedAutomatically.IsNull() && !data.IsCreatedAutomatically.IsUnknown() {
        filters["isCreatedAutomatically"] = data.IsCreatedAutomatically.ValueBool()
        filterNames = append(filterNames, "is_created_automatically = "+fmt.Sprintf("%t", data.IsCreatedAutomatically.ValueBool()))
    }
    if !data.RemediationNotes.IsNull() && !data.RemediationNotes.IsUnknown() {
        filters["remediationNotes"] = data.RemediationNotes.ValueString()
        filterNames = append(filterNames, "remediation_notes = "+fmt.Sprintf("%q", data.RemediationNotes.ValueString()))
    }
    if !data.AlertNumber.IsNull() && !data.AlertNumber.IsUnknown() {
        filters["alertNumber"] = lookupNumber(data.AlertNumber)
        filterNames = append(filterNames, "alert_number = "+data.AlertNumber.ValueBigFloat().String())
    }
    if !data.AlertNumberWithPrefix.IsNull() && !data.AlertNumberWithPrefix.IsUnknown() {
        filters["alertNumberWithPrefix"] = data.AlertNumberWithPrefix.ValueString()
        filterNames = append(filterNames, "alert_number_with_prefix = "+fmt.Sprintf("%q", data.AlertNumberWithPrefix.ValueString()))
    }
    if !data.AlertEpisodeId.IsNull() && !data.AlertEpisodeId.IsUnknown() {
        filters["alertEpisodeId"] = data.AlertEpisodeId.ValueString()
        filterNames = append(filterNames, "alert_episode_id = "+fmt.Sprintf("%q", data.AlertEpisodeId.ValueString()))
    }
    if !data.IsPrivate.IsNull() && !data.IsPrivate.IsUnknown() {
        filters["isPrivate"] = data.IsPrivate.ValueBool()
        filterNames = append(filterNames, "is_private = "+fmt.Sprintf("%t", data.IsPrivate.ValueBool()))
    }
    if !data.EnableReminders.IsNull() && !data.EnableReminders.IsUnknown() {
        filters["enableReminders"] = data.EnableReminders.ValueBool()
        filterNames = append(filterNames, "enable_reminders = "+fmt.Sprintf("%t", data.EnableReminders.ValueBool()))
    }
    if !data.ReminderNotificationSentCount.IsNull() && !data.ReminderNotificationSentCount.IsUnknown() {
        filters["reminderNotificationSentCount"] = lookupNumber(data.ReminderNotificationSentCount)
        filterNames = append(filterNames, "reminder_notification_sent_count = "+data.ReminderNotificationSentCount.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the alert up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the alert up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "title": true,
        "description": true,
        "impactStartedAt": true,
        "createdByUserId": true,
        "monitorId": true,
        "onCallDutyPolicies": true,
        "hosts": true,
        "kubernetesClusters": true,
        "kubernetesResources": true,
        "kubernetesContainers": true,
        "dockerHosts": true,
        "podmanHosts": true,
        "proxmoxClusters": true,
        "vmwareVCenters": true,
        "iotFleets": true,
        "dockerSwarmClusters": true,
        "cephClusters": true,
        "storageArrays": true,
        "databaseServers": true,
        "dockerResources": true,
        "podmanResources": true,
        "services": true,
        "serviceLevelObjectives": true,
        "labels": true,
        "currentAlertStateId": true,
        "alertSeverityId": true,
        "monitorStatusWhenThisAlertWasCreatedId": true,
        "customFields": true,
        "isOwnerNotifiedOfAlertCreation": true,
        "rootCause": true,
        "createdStateLog": true,
        "createdCriteriaId": true,
        "seriesFingerprint": true,
        "seriesLabels": true,
        "monitorSummary": true,
        "createdByProbeId": true,
        "isCreatedAutomatically": true,
        "remediationNotes": true,
        "telemetryQuery": true,
        "alertNumber": true,
        "alertNumberWithPrefix": true,
        "alertEpisodeId": true,
        "isPrivate": true,
        "enableReminders": true,
        "nextReminderNotificationAt": true,
        "reminderNotificationSentCount": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/alert/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read alert: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/alert/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list alert, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list alert: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No alert matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one alert matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for alert.")
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
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    if obj, ok := item["impactStartedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ImpactStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ImpactStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ImpactStartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ImpactStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ImpactStartedAt = types.StringNull()
        }
    } else if val, ok := item["impactStartedAt"].(string); ok {
        data.ImpactStartedAt = types.StringValue(val)
    } else {
        data.ImpactStartedAt = types.StringNull()
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
    if val, ok := item["hosts"].([]interface{}); ok {
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
        data.Hosts = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Hosts = types.SetNull(types.StringType)
    }
    if val, ok := item["kubernetesClusters"].([]interface{}); ok {
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
        data.KubernetesClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.KubernetesClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["kubernetesResources"].([]interface{}); ok {
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
        data.KubernetesResources = types.SetValueMust(types.StringType, setItems)
    } else {
        data.KubernetesResources = types.SetNull(types.StringType)
    }
    if val, ok := item["kubernetesContainers"].([]interface{}); ok {
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
        data.KubernetesContainers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.KubernetesContainers = types.SetNull(types.StringType)
    }
    if val, ok := item["dockerHosts"].([]interface{}); ok {
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
        data.DockerHosts = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DockerHosts = types.SetNull(types.StringType)
    }
    if val, ok := item["podmanHosts"].([]interface{}); ok {
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
        data.PodmanHosts = types.SetValueMust(types.StringType, setItems)
    } else {
        data.PodmanHosts = types.SetNull(types.StringType)
    }
    if val, ok := item["proxmoxClusters"].([]interface{}); ok {
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
        data.ProxmoxClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.ProxmoxClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["vmwareVCenters"].([]interface{}); ok {
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
        data.VmwareVCenters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.VmwareVCenters = types.SetNull(types.StringType)
    }
    if val, ok := item["iotFleets"].([]interface{}); ok {
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
        data.IotFleets = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IotFleets = types.SetNull(types.StringType)
    }
    if val, ok := item["dockerSwarmClusters"].([]interface{}); ok {
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
        data.DockerSwarmClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DockerSwarmClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["cephClusters"].([]interface{}); ok {
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
        data.CephClusters = types.SetValueMust(types.StringType, setItems)
    } else {
        data.CephClusters = types.SetNull(types.StringType)
    }
    if val, ok := item["storageArrays"].([]interface{}); ok {
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
        data.StorageArrays = types.SetValueMust(types.StringType, setItems)
    } else {
        data.StorageArrays = types.SetNull(types.StringType)
    }
    if val, ok := item["databaseServers"].([]interface{}); ok {
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
        data.DatabaseServers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DatabaseServers = types.SetNull(types.StringType)
    }
    if val, ok := item["dockerResources"].([]interface{}); ok {
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
        data.DockerResources = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DockerResources = types.SetNull(types.StringType)
    }
    if val, ok := item["podmanResources"].([]interface{}); ok {
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
        data.PodmanResources = types.SetValueMust(types.StringType, setItems)
    } else {
        data.PodmanResources = types.SetNull(types.StringType)
    }
    if val, ok := item["services"].([]interface{}); ok {
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
        data.Services = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Services = types.SetNull(types.StringType)
    }
    if val, ok := item["serviceLevelObjectives"].([]interface{}); ok {
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
        data.ServiceLevelObjectives = types.SetValueMust(types.StringType, setItems)
    } else {
        data.ServiceLevelObjectives = types.SetNull(types.StringType)
    }
    if val, ok := item["labels"].([]interface{}); ok {
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
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if obj, ok := item["currentAlertStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentAlertStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentAlertStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentAlertStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentAlertStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentAlertStateId = types.StringNull()
        }
    } else if val, ok := item["currentAlertStateId"].(string); ok {
        data.CurrentAlertStateId = types.StringValue(val)
    } else {
        data.CurrentAlertStateId = types.StringNull()
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
    if obj, ok := item["monitorStatusWhenThisAlertWasCreatedId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorStatusWhenThisAlertWasCreatedId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorStatusWhenThisAlertWasCreatedId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorStatusWhenThisAlertWasCreatedId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorStatusWhenThisAlertWasCreatedId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorStatusWhenThisAlertWasCreatedId = types.StringNull()
        }
    } else if val, ok := item["monitorStatusWhenThisAlertWasCreatedId"].(string); ok {
        data.MonitorStatusWhenThisAlertWasCreatedId = types.StringValue(val)
    } else {
        data.MonitorStatusWhenThisAlertWasCreatedId = types.StringNull()
    }
    if obj, ok := item["customFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := item["customFields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if val, ok := item["isOwnerNotifiedOfAlertCreation"].(bool); ok {
        data.IsOwnerNotifiedOfAlertCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfAlertCreation = types.BoolNull()
    }
    if obj, ok := item["rootCause"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootCause = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootCause = types.StringValue(string(jsonBytes))
        } else {
            data.RootCause = types.StringNull()
        }
    } else if val, ok := item["rootCause"].(string); ok {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if obj, ok := item["createdStateLog"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedStateLog = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedStateLog = types.StringNull()
        }
    } else if val, ok := item["createdStateLog"].(string); ok {
        data.CreatedStateLog = types.StringValue(val)
    } else {
        data.CreatedStateLog = types.StringNull()
    }
    if obj, ok := item["createdCriteriaId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedCriteriaId = types.StringNull()
        }
    } else if val, ok := item["createdCriteriaId"].(string); ok {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if obj, ok := item["seriesFingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeriesFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SeriesFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SeriesFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SeriesFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.SeriesFingerprint = types.StringNull()
        }
    } else if val, ok := item["seriesFingerprint"].(string); ok {
        data.SeriesFingerprint = types.StringValue(val)
    } else {
        data.SeriesFingerprint = types.StringNull()
    }
    if obj, ok := item["seriesLabels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeriesLabels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SeriesLabels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SeriesLabels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SeriesLabels = types.StringValue(string(jsonBytes))
        } else {
            data.SeriesLabels = types.StringNull()
        }
    } else if val, ok := item["seriesLabels"].(string); ok {
        data.SeriesLabels = types.StringValue(val)
    } else {
        data.SeriesLabels = types.StringNull()
    }
    if obj, ok := item["monitorSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorSummary = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorSummary = types.StringNull()
        }
    } else if val, ok := item["monitorSummary"].(string); ok {
        data.MonitorSummary = types.StringValue(val)
    } else {
        data.MonitorSummary = types.StringNull()
    }
    if obj, ok := item["createdByProbeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByProbeId = types.StringNull()
        }
    } else if val, ok := item["createdByProbeId"].(string); ok {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := item["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    } else {
        data.IsCreatedAutomatically = types.BoolNull()
    }
    if obj, ok := item["remediationNotes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationNotes = types.StringNull()
        }
    } else if val, ok := item["remediationNotes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if obj, ok := item["telemetryQuery"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryQuery = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryQuery = types.StringNull()
        }
    } else if val, ok := item["telemetryQuery"].(string); ok {
        data.TelemetryQuery = types.StringValue(val)
    } else {
        data.TelemetryQuery = types.StringNull()
    }
    if val, ok := item["alertNumber"].(float64); ok {
        data.AlertNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertNumber = types.NumberNull()
        }
    } else {
        data.AlertNumber = types.NumberNull()
    }
    if obj, ok := item["alertNumberWithPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := item["alertNumberWithPrefix"].(string); ok {
        data.AlertNumberWithPrefix = types.StringValue(val)
    } else {
        data.AlertNumberWithPrefix = types.StringNull()
    }
    if obj, ok := item["alertEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertEpisodeId = types.StringNull()
        }
    } else if val, ok := item["alertEpisodeId"].(string); ok {
        data.AlertEpisodeId = types.StringValue(val)
    } else {
        data.AlertEpisodeId = types.StringNull()
    }
    if val, ok := item["isPrivate"].(bool); ok {
        data.IsPrivate = types.BoolValue(val)
    } else {
        data.IsPrivate = types.BoolNull()
    }
    if val, ok := item["enableReminders"].(bool); ok {
        data.EnableReminders = types.BoolValue(val)
    } else {
        data.EnableReminders = types.BoolNull()
    }
    if obj, ok := item["nextReminderNotificationAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextReminderNotificationAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextReminderNotificationAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextReminderNotificationAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextReminderNotificationAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextReminderNotificationAt = types.StringNull()
        }
    } else if val, ok := item["nextReminderNotificationAt"].(string); ok {
        data.NextReminderNotificationAt = types.StringValue(val)
    } else {
        data.NextReminderNotificationAt = types.StringNull()
    }
    if val, ok := item["reminderNotificationSentCount"].(float64); ok {
        data.ReminderNotificationSentCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["reminderNotificationSentCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ReminderNotificationSentCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReminderNotificationSentCount = types.NumberNull()
        }
    } else {
        data.ReminderNotificationSentCount = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
