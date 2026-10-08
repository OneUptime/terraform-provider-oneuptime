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
var _ datasource.DataSource = &IncidentDataSource{}

func NewIncidentDataSource() datasource.DataSource {
    return &IncidentDataSource{}
}

// IncidentDataSource defines the data source implementation.
type IncidentDataSource struct {
    client *Client
}

// IncidentDataSourceModel describes the data source data model.
type IncidentDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    DeclaredAt types.String `tfsdk:"declared_at"`
    ImpactStartedAt types.String `tfsdk:"impact_started_at"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Monitors types.Set `tfsdk:"monitors"`
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
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    Labels types.Set `tfsdk:"labels"`
    CurrentIncidentStateId types.String `tfsdk:"current_incident_state_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    SubscriberNotificationStatusOnIncidentCreated types.String `tfsdk:"subscriber_notification_status_on_incident_created"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    SubscriberNotificationStatusOnPostmortemPublished types.String `tfsdk:"subscriber_notification_status_on_postmortem_published"`
    SubscriberNotificationStatusMessageOnPostmortemPublished types.String `tfsdk:"subscriber_notification_status_message_on_postmortem_published"`
    ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_incident_created"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    RootCause types.String `tfsdk:"root_cause"`
    PostmortemNote types.String `tfsdk:"postmortem_note"`
    ShowPostmortemOnStatusPage types.Bool `tfsdk:"show_postmortem_on_status_page"`
    NotifySubscribersOnPostmortemPublished types.Bool `tfsdk:"notify_subscribers_on_postmortem_published"`
    PostmortemPostedAt types.String `tfsdk:"postmortem_posted_at"`
    PostmortemAttachments types.Set `tfsdk:"postmortem_attachments"`
    CreatedStateLog types.String `tfsdk:"created_state_log"`
    CreatedCriteriaId types.String `tfsdk:"created_criteria_id"`
    CreatedIncidentTemplateId types.String `tfsdk:"created_incident_template_id"`
    SeriesFingerprint types.String `tfsdk:"series_fingerprint"`
    SeriesLabels types.String `tfsdk:"series_labels"`
    MonitorSummary types.String `tfsdk:"monitor_summary"`
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
    HoldsMonitors types.Bool `tfsdk:"holds_monitors"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.String `tfsdk:"telemetry_query"`
    IncidentNumber types.Number `tfsdk:"incident_number"`
    IncidentNumberWithPrefix types.String `tfsdk:"incident_number_with_prefix"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    StatusPages types.Set `tfsdk:"status_pages"`
    IsScopedToStatusPages types.Bool `tfsdk:"is_scoped_to_status_pages"`
    StatusPagesNotifiedOnCreation types.String `tfsdk:"status_pages_notified_on_creation"`
    IsPrivate types.Bool `tfsdk:"is_private"`
    EnableReminders types.Bool `tfsdk:"enable_reminders"`
    NextReminderNotificationAt types.String `tfsdk:"next_reminder_notification_at"`
    ReminderNotificationSentCount types.Number `tfsdk:"reminder_notification_sent_count"`
    IncidentEpisodeId types.String `tfsdk:"incident_episode_id"`
}

func (d *IncidentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident"
}

func (d *IncidentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage incidents for your project Look up an existing incident by `id`, or by any of its other arguments (`change_monitor_status_to_id`, `created_by_probe_id`, `created_by_user_id`, ...): each one set must match, and exactly one incident may match them all.",

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
                MarkdownDescription: "Title of this incident.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Short description of this incident. This is in markdown and will be visible on the status page.",
                Optional: true,
                Computed: true,
            },
            "declared_at": schema.StringAttribute{
                MarkdownDescription: "Date and time when this incident was declared.",
                Computed: true,
            },
            "impact_started_at": schema.StringAttribute{
                MarkdownDescription: "When customer impact actually began. Left blank until someone records it - never inferred, because a guessed value is worse than no value.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "List of monitors affected by this incident. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "hosts": schema.SetAttribute{
                MarkdownDescription: "List of hosts affected by this incident. IDs of `oneuptime_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes clusters affected by this incident. IDs of `oneuptime_kubernetes_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_resources": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this incident.",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_containers": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes containers affected by this incident.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Docker hosts affected by this incident. IDs of `oneuptime_docker_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Podman hosts affected by this incident. IDs of `oneuptime_podman_host` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "proxmox_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Proxmox clusters affected by this incident. IDs of `oneuptime_proxmox_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "vmware_v_centers": schema.SetAttribute{
                MarkdownDescription: "List of vCenters affected by this incident. IDs of `oneuptime_vcenter` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "iot_fleets": schema.SetAttribute{
                MarkdownDescription: "List of IoT fleets affected by this incident. IDs of `oneuptime_iot_fleet` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_swarm_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Docker Swarm clusters affected by this incident. IDs of `oneuptime_docker_swarm_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "ceph_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Ceph clusters affected by this incident. IDs of `oneuptime_ceph_cluster` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "storage_arrays": schema.SetAttribute{
                MarkdownDescription: "List of storage arrays affected by this incident. IDs of `oneuptime_storage_array` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "database_servers": schema.SetAttribute{
                MarkdownDescription: "List of databases affected by this incident. IDs of `oneuptime_database` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_resources": schema.SetAttribute{
                MarkdownDescription: "List of Docker resources (containers, images, networks, volumes) affected by this incident.",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_resources": schema.SetAttribute{
                MarkdownDescription: "List of Podman resources (containers, images, networks, volumes) affected by this incident.",
                Computed: true,
                ElementType: types.StringType,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "List of services affected by this incident. IDs of `oneuptime_service` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_level_objectives": schema.SetAttribute{
                MarkdownDescription: "List of Service Level Objectives (SLOs) affected by this incident. IDs of `oneuptime_service_level_objective` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies affected by this incident. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "Current Incident State ID. The ID of a `oneuptime_incident_state`.",
                Optional: true,
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "Incident Severity ID. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Monitor Status Object ID. All monitors connected to this incident will be changed to this status when the incident is created. The ID of a `oneuptime_monitor_status`.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_on_incident_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this incident.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_on_postmortem_published": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this incident postmortem.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message_on_postmortem_published": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications on postmortem published - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this incident?",
                Optional: true,
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "The incident's custom field values, keyed by each incident custom field's name. When a user or an API key creates or updates an incident, each value it sets or changes must fit its field - a number for a Number field, true or false for a Boolean, one of the options for a Dropdown, and so on - or the request is refused. Values left as they were, keys that are not the name of a field and empty values are not checked. Required on Create is not enforced here: it applies to the dashboard's Declare Incident form only. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?",
                Optional: true,
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "What is the root cause of this incident?",
                Optional: true,
                Computed: true,
            },
            "postmortem_note": schema.StringAttribute{
                MarkdownDescription: "Document the postmortem summary for this incident.",
                Optional: true,
                Computed: true,
            },
            "show_postmortem_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should the postmortem note and attachments be visible on the status page once published?",
                Optional: true,
                Computed: true,
            },
            "notify_subscribers_on_postmortem_published": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified when the postmortem is published?",
                Optional: true,
                Computed: true,
            },
            "postmortem_posted_at": schema.StringAttribute{
                MarkdownDescription: "Timestamp that will be shown alongside the published postmortem on the status page.",
                Computed: true,
            },
            "postmortem_attachments": schema.SetAttribute{
                MarkdownDescription: "Files that accompany the postmortem note and can be shared publicly when enabled. IDs of `oneuptime_file` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_state_log": schema.StringAttribute{
                MarkdownDescription: "A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the criteria that created it.",
                Optional: true,
                Computed: true,
            },
            "created_incident_template_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the incident template that was used for creation.",
                Optional: true,
                Computed: true,
            },
            "series_fingerprint": schema.StringAttribute{
                MarkdownDescription: "For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one incident is created per affected series.",
                Optional: true,
                Computed: true,
            },
            "series_labels": schema.StringAttribute{
                MarkdownDescription: "Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this incident was created from a per-series metric breach. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "monitor_summary": schema.StringAttribute{
                MarkdownDescription: "The monitor summary captured at the moment this incident was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the probe that created it. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Is this incident created by OneUptime Probe or Workers automatically (and not created manually by a user)?",
                Optional: true,
                Computed: true,
            },
            "holds_monitors": schema.BoolAttribute{
                MarkdownDescription: "Whether this incident is holding its monitors - keeping them in its monitor status, or their monitoring paused - so that resolving it gives them back: their monitoring resumes and their status returns to operational. True from when the incident is declared open, or from when an edit while it is open puts its monitors in its monitor status. False for an incident declared already resolved, which never held them, and once a resolve has given them back. Empty for incidents from before it was recorded, which give their monitors back when they are resolved. Set by OneUptime; it cannot be written.",
                Optional: true,
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Notes on how to remediate this incident. This is in markdown.",
                Optional: true,
                Computed: true,
            },
            "telemetry_query": schema.StringAttribute{
                MarkdownDescription: "Telemetry query for this incident. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "incident_number": schema.NumberAttribute{
                MarkdownDescription: "Incident Number.",
                Optional: true,
                Computed: true,
            },
            "incident_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Incident number with prefix (e.g., 'INC-42' or '#42').",
                Optional: true,
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should this incident be visible on the status page?",
                Optional: true,
                Computed: true,
            },
            "status_pages": schema.SetAttribute{
                MarkdownDescription: "Limit this incident to these status pages. When set, the incident is shown on, and notifies the subscribers of, only these pages among the status pages that list its monitors. Leave empty to reach every status page that lists its monitors. IDs of `oneuptime_status_page` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_scoped_to_status_pages": schema.BoolAttribute{
                MarkdownDescription: "Whether this incident is limited to the status pages in Status Pages. Derived from Status Pages; any value sent for it is ignored.",
                Optional: true,
                Computed: true,
            },
            "status_pages_notified_on_creation": schema.StringAttribute{
                MarkdownDescription: "IDs of the status pages whose subscribers were sent the notification that this incident was created. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_private": schema.BoolAttribute{
                MarkdownDescription: "If true, this incident is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners. Private incidents are hidden from status pages.",
                Optional: true,
                Computed: true,
            },
            "enable_reminders": schema.BoolAttribute{
                MarkdownDescription: "Should reminder notifications be sent to owners while this incident is still open? Reminders are sent based on the reminder rules configured for this project.",
                Optional: true,
                Computed: true,
            },
            "next_reminder_notification_at": schema.StringAttribute{
                MarkdownDescription: "When will the next reminder notification be sent to owners of this incident? This is set automatically based on the reminder rules configured for this project.",
                Computed: true,
            },
            "reminder_notification_sent_count": schema.NumberAttribute{
                MarkdownDescription: "How many reminder notifications have been sent to owners of this incident so far.",
                Optional: true,
                Computed: true,
            },
            "incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of the latest Incident Episode this incident is a member of. Read-only: set by OneUptime when the incident is added to or removed from an episode's members (Incident Episode Member). The ID of a `oneuptime_incident_episode`.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.CurrentIncidentStateId.IsNull() && !data.CurrentIncidentStateId.IsUnknown() {
        filters["currentIncidentStateId"] = data.CurrentIncidentStateId.ValueString()
        filterNames = append(filterNames, "current_incident_state_id = "+fmt.Sprintf("%q", data.CurrentIncidentStateId.ValueString()))
    }
    if !data.IncidentSeverityId.IsNull() && !data.IncidentSeverityId.IsUnknown() {
        filters["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
        filterNames = append(filterNames, "incident_severity_id = "+fmt.Sprintf("%q", data.IncidentSeverityId.ValueString()))
    }
    if !data.ChangeMonitorStatusToId.IsNull() && !data.ChangeMonitorStatusToId.IsUnknown() {
        filters["changeMonitorStatusToId"] = data.ChangeMonitorStatusToId.ValueString()
        filterNames = append(filterNames, "change_monitor_status_to_id = "+fmt.Sprintf("%q", data.ChangeMonitorStatusToId.ValueString()))
    }
    if !data.SubscriberNotificationStatusOnIncidentCreated.IsNull() && !data.SubscriberNotificationStatusOnIncidentCreated.IsUnknown() {
        filters["subscriberNotificationStatusOnIncidentCreated"] = data.SubscriberNotificationStatusOnIncidentCreated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_incident_created = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnIncidentCreated.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessage.IsNull() && !data.SubscriberNotificationStatusMessage.IsUnknown() {
        filters["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessage.ValueString()))
    }
    if !data.SubscriberNotificationStatusOnPostmortemPublished.IsNull() && !data.SubscriberNotificationStatusOnPostmortemPublished.IsUnknown() {
        filters["subscriberNotificationStatusOnPostmortemPublished"] = data.SubscriberNotificationStatusOnPostmortemPublished.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_postmortem_published = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnPostmortemPublished.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessageOnPostmortemPublished.IsNull() && !data.SubscriberNotificationStatusMessageOnPostmortemPublished.IsUnknown() {
        filters["subscriberNotificationStatusMessageOnPostmortemPublished"] = data.SubscriberNotificationStatusMessageOnPostmortemPublished.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message_on_postmortem_published = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessageOnPostmortemPublished.ValueString()))
    }
    if !data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.IsNull() && !data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"] = data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified_on_incident_created = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.ValueBool()))
    }
    if !data.IsOwnerNotifiedOfResourceCreation.IsNull() && !data.IsOwnerNotifiedOfResourceCreation.IsUnknown() {
        filters["isOwnerNotifiedOfResourceCreation"] = data.IsOwnerNotifiedOfResourceCreation.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_resource_creation = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfResourceCreation.ValueBool()))
    }
    if !data.RootCause.IsNull() && !data.RootCause.IsUnknown() {
        filters["rootCause"] = data.RootCause.ValueString()
        filterNames = append(filterNames, "root_cause = "+fmt.Sprintf("%q", data.RootCause.ValueString()))
    }
    if !data.PostmortemNote.IsNull() && !data.PostmortemNote.IsUnknown() {
        filters["postmortemNote"] = data.PostmortemNote.ValueString()
        filterNames = append(filterNames, "postmortem_note = "+fmt.Sprintf("%q", data.PostmortemNote.ValueString()))
    }
    if !data.ShowPostmortemOnStatusPage.IsNull() && !data.ShowPostmortemOnStatusPage.IsUnknown() {
        filters["showPostmortemOnStatusPage"] = data.ShowPostmortemOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_postmortem_on_status_page = "+fmt.Sprintf("%t", data.ShowPostmortemOnStatusPage.ValueBool()))
    }
    if !data.NotifySubscribersOnPostmortemPublished.IsNull() && !data.NotifySubscribersOnPostmortemPublished.IsUnknown() {
        filters["notifySubscribersOnPostmortemPublished"] = data.NotifySubscribersOnPostmortemPublished.ValueBool()
        filterNames = append(filterNames, "notify_subscribers_on_postmortem_published = "+fmt.Sprintf("%t", data.NotifySubscribersOnPostmortemPublished.ValueBool()))
    }
    if !data.CreatedCriteriaId.IsNull() && !data.CreatedCriteriaId.IsUnknown() {
        filters["createdCriteriaId"] = data.CreatedCriteriaId.ValueString()
        filterNames = append(filterNames, "created_criteria_id = "+fmt.Sprintf("%q", data.CreatedCriteriaId.ValueString()))
    }
    if !data.CreatedIncidentTemplateId.IsNull() && !data.CreatedIncidentTemplateId.IsUnknown() {
        filters["createdIncidentTemplateId"] = data.CreatedIncidentTemplateId.ValueString()
        filterNames = append(filterNames, "created_incident_template_id = "+fmt.Sprintf("%q", data.CreatedIncidentTemplateId.ValueString()))
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
    if !data.HoldsMonitors.IsNull() && !data.HoldsMonitors.IsUnknown() {
        filters["holdsMonitors"] = data.HoldsMonitors.ValueBool()
        filterNames = append(filterNames, "holds_monitors = "+fmt.Sprintf("%t", data.HoldsMonitors.ValueBool()))
    }
    if !data.RemediationNotes.IsNull() && !data.RemediationNotes.IsUnknown() {
        filters["remediationNotes"] = data.RemediationNotes.ValueString()
        filterNames = append(filterNames, "remediation_notes = "+fmt.Sprintf("%q", data.RemediationNotes.ValueString()))
    }
    if !data.IncidentNumber.IsNull() && !data.IncidentNumber.IsUnknown() {
        filters["incidentNumber"] = lookupNumber(data.IncidentNumber)
        filterNames = append(filterNames, "incident_number = "+data.IncidentNumber.ValueBigFloat().String())
    }
    if !data.IncidentNumberWithPrefix.IsNull() && !data.IncidentNumberWithPrefix.IsUnknown() {
        filters["incidentNumberWithPrefix"] = data.IncidentNumberWithPrefix.ValueString()
        filterNames = append(filterNames, "incident_number_with_prefix = "+fmt.Sprintf("%q", data.IncidentNumberWithPrefix.ValueString()))
    }
    if !data.IsVisibleOnStatusPage.IsNull() && !data.IsVisibleOnStatusPage.IsUnknown() {
        filters["isVisibleOnStatusPage"] = data.IsVisibleOnStatusPage.ValueBool()
        filterNames = append(filterNames, "is_visible_on_status_page = "+fmt.Sprintf("%t", data.IsVisibleOnStatusPage.ValueBool()))
    }
    if !data.IsScopedToStatusPages.IsNull() && !data.IsScopedToStatusPages.IsUnknown() {
        filters["isScopedToStatusPages"] = data.IsScopedToStatusPages.ValueBool()
        filterNames = append(filterNames, "is_scoped_to_status_pages = "+fmt.Sprintf("%t", data.IsScopedToStatusPages.ValueBool()))
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
    if !data.IncidentEpisodeId.IsNull() && !data.IncidentEpisodeId.IsUnknown() {
        filters["incidentEpisodeId"] = data.IncidentEpisodeId.ValueString()
        filterNames = append(filterNames, "incident_episode_id = "+fmt.Sprintf("%q", data.IncidentEpisodeId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "title": true,
        "description": true,
        "declaredAt": true,
        "impactStartedAt": true,
        "slug": true,
        "createdByUserId": true,
        "monitors": true,
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
        "onCallDutyPolicies": true,
        "labels": true,
        "currentIncidentStateId": true,
        "incidentSeverityId": true,
        "changeMonitorStatusToId": true,
        "subscriberNotificationStatusOnIncidentCreated": true,
        "subscriberNotificationStatusMessage": true,
        "subscriberNotificationStatusOnPostmortemPublished": true,
        "subscriberNotificationStatusMessageOnPostmortemPublished": true,
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "rootCause": true,
        "postmortemNote": true,
        "showPostmortemOnStatusPage": true,
        "notifySubscribersOnPostmortemPublished": true,
        "postmortemPostedAt": true,
        "postmortemAttachments": true,
        "createdStateLog": true,
        "createdCriteriaId": true,
        "createdIncidentTemplateId": true,
        "seriesFingerprint": true,
        "seriesLabels": true,
        "monitorSummary": true,
        "createdByProbeId": true,
        "isCreatedAutomatically": true,
        "holdsMonitors": true,
        "remediationNotes": true,
        "telemetryQuery": true,
        "incidentNumber": true,
        "incidentNumberWithPrefix": true,
        "isVisibleOnStatusPage": true,
        "statusPages": true,
        "isScopedToStatusPages": true,
        "statusPagesNotifiedOnCreation": true,
        "isPrivate": true,
        "enableReminders": true,
        "nextReminderNotificationAt": true,
        "reminderNotificationSentCount": true,
        "incidentEpisodeId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident.")
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
    if obj, ok := item["declaredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeclaredAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeclaredAt = types.StringNull()
        }
    } else if val, ok := item["declaredAt"].(string); ok {
        data.DeclaredAt = types.StringValue(val)
    } else {
        data.DeclaredAt = types.StringNull()
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if val, ok := item["monitors"].([]interface{}); ok {
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
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
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
    if obj, ok := item["currentIncidentStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentIncidentStateId = types.StringNull()
        }
    } else if val, ok := item["currentIncidentStateId"].(string); ok {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
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
    if obj, ok := item["changeMonitorStatusToId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMonitorStatusToId = types.StringNull()
        }
    } else if val, ok := item["changeMonitorStatusToId"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusOnIncidentCreated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnIncidentCreated"].(string); ok {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusMessage"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusOnPostmortemPublished"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnPostmortemPublished"].(string); ok {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusMessageOnPostmortemPublished"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusMessageOnPostmortemPublished"].(string); ok {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
    }
    if val, ok := item["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolNull()
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
    if val, ok := item["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
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
    if obj, ok := item["postmortemNote"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostmortemNote = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemNote = types.StringNull()
        }
    } else if val, ok := item["postmortemNote"].(string); ok {
        data.PostmortemNote = types.StringValue(val)
    } else {
        data.PostmortemNote = types.StringNull()
    }
    if val, ok := item["showPostmortemOnStatusPage"].(bool); ok {
        data.ShowPostmortemOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowPostmortemOnStatusPage = types.BoolNull()
    }
    if val, ok := item["notifySubscribersOnPostmortemPublished"].(bool); ok {
        data.NotifySubscribersOnPostmortemPublished = types.BoolValue(val)
    } else {
        data.NotifySubscribersOnPostmortemPublished = types.BoolNull()
    }
    if obj, ok := item["postmortemPostedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemPostedAt = types.StringNull()
        }
    } else if val, ok := item["postmortemPostedAt"].(string); ok {
        data.PostmortemPostedAt = types.StringValue(val)
    } else {
        data.PostmortemPostedAt = types.StringNull()
    }
    if val, ok := item["postmortemAttachments"].([]interface{}); ok {
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
        data.PostmortemAttachments = types.SetValueMust(types.StringType, setItems)
    } else {
        data.PostmortemAttachments = types.SetNull(types.StringType)
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
    if obj, ok := item["createdIncidentTemplateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedIncidentTemplateId = types.StringNull()
        }
    } else if val, ok := item["createdIncidentTemplateId"].(string); ok {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
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
    if val, ok := item["holdsMonitors"].(bool); ok {
        data.HoldsMonitors = types.BoolValue(val)
    } else {
        data.HoldsMonitors = types.BoolNull()
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
    if val, ok := item["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentNumber = types.NumberNull()
        }
    } else {
        data.IncidentNumber = types.NumberNull()
    }
    if obj, ok := item["incidentNumberWithPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := item["incidentNumberWithPrefix"].(string); ok {
        data.IncidentNumberWithPrefix = types.StringValue(val)
    } else {
        data.IncidentNumberWithPrefix = types.StringNull()
    }
    if val, ok := item["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := item["statusPages"].([]interface{}); ok {
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
        data.StatusPages = types.SetValueMust(types.StringType, setItems)
    } else {
        data.StatusPages = types.SetNull(types.StringType)
    }
    if val, ok := item["isScopedToStatusPages"].(bool); ok {
        data.IsScopedToStatusPages = types.BoolValue(val)
    } else {
        data.IsScopedToStatusPages = types.BoolNull()
    }
    if obj, ok := item["statusPagesNotifiedOnCreation"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPagesNotifiedOnCreation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPagesNotifiedOnCreation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPagesNotifiedOnCreation = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPagesNotifiedOnCreation = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPagesNotifiedOnCreation = types.StringNull()
        }
    } else if val, ok := item["statusPagesNotifiedOnCreation"].(string); ok {
        data.StatusPagesNotifiedOnCreation = types.StringValue(val)
    } else {
        data.StatusPagesNotifiedOnCreation = types.StringNull()
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
    if obj, ok := item["incidentEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeId = types.StringNull()
        }
    } else if val, ok := item["incidentEpisodeId"].(string); ok {
        data.IncidentEpisodeId = types.StringValue(val)
    } else {
        data.IncidentEpisodeId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
