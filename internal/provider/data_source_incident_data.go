package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IncidentDataDataSource{}

func NewIncidentDataDataSource() datasource.DataSource {
    return &IncidentDataDataSource{}
}

// IncidentDataDataSource defines the data source implementation.
type IncidentDataDataSource struct {
    client *Client
}

// IncidentDataDataSourceModel describes the data source data model.
type IncidentDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    DeclaredAt types.String `tfsdk:"declared_at"`
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
    DockerSwarmClusters types.Set `tfsdk:"docker_swarm_clusters"`
    CephClusters types.Set `tfsdk:"ceph_clusters"`
    DockerResources types.Set `tfsdk:"docker_resources"`
    PodmanResources types.Set `tfsdk:"podman_resources"`
    Services types.Set `tfsdk:"services"`
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
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.String `tfsdk:"telemetry_query"`
    IncidentNumber types.Number `tfsdk:"incident_number"`
    IncidentNumberWithPrefix types.String `tfsdk:"incident_number_with_prefix"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    IsPrivate types.Bool `tfsdk:"is_private"`
    IncidentEpisodeId types.String `tfsdk:"incident_episode_id"`
}

func (d *IncidentDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_data"
}

func (d *IncidentDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Short description of this incident. This is in markdown and will be visible on the status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "declared_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "List of monitors affected by this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "hosts": schema.SetAttribute{
                MarkdownDescription: "List of hosts affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes clusters affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_resources": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_containers": schema.SetAttribute{
                MarkdownDescription: "List of Kubernetes containers affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Docker hosts affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_hosts": schema.SetAttribute{
                MarkdownDescription: "List of Podman hosts affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "proxmox_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Proxmox clusters affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_swarm_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Docker Swarm clusters affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "ceph_clusters": schema.SetAttribute{
                MarkdownDescription: "List of Ceph clusters affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "docker_resources": schema.SetAttribute{
                MarkdownDescription: "List of Docker resources (containers, images, networks, volumes) affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "podman_resources": schema.SetAttribute{
                MarkdownDescription: "List of Podman resources (containers, images, networks, volumes) affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "List of services affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "subscriber_notification_status_on_incident_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "subscriber_notification_status_on_postmortem_published": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this incident postmortem. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "subscriber_notification_status_message_on_postmortem_published": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications on postmortem published - includes success messages, failure reasons, or skip reasons. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this incident?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "What is the root cause of this incident?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "postmortem_note": schema.StringAttribute{
                MarkdownDescription: "Document the postmortem summary for this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "show_postmortem_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should the postmortem note and attachments be visible on the status page once published?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "notify_subscribers_on_postmortem_published": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified when the postmortem is published?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "postmortem_posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "postmortem_attachments": schema.SetAttribute{
                MarkdownDescription: "Files that accompany the postmortem note and can be shared publicly when enabled.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_state_log": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the criteria that created it.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_incident_template_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the incident template that was used for creation.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "series_fingerprint": schema.StringAttribute{
                MarkdownDescription: "For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one incident is created per affected series.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "series_labels": schema.StringAttribute{
                MarkdownDescription: "Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this incident was created from a per-series metric breach.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Is this incident created by OneUptime Probe or Workers automatically (and not created manually by a user)?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Notes on how to remediate this incident. This is in markdown.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "telemetry_query": schema.StringAttribute{
                MarkdownDescription: "Telemetry query for this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "incident_number": schema.NumberAttribute{
                MarkdownDescription: "Incident Number. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incident_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Incident number with prefix (e.g., 'INC-42' or '#42'). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should this incident be visible on the status page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "is_private": schema.BoolAttribute{
                MarkdownDescription: "If true, this incident is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners. Private incidents are hidden from status pages.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident]",
                Computed: true,
            },
            "incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_data, got error: %s", err))
        return
    }

    var incidentDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentDataResponse["data"].(map[string]interface{}); ok {
        incidentDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["declared_at"].(string); ok {
        data.DeclaredAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Monitors = setValue
    }
    if val, ok := incidentDataResponse["hosts"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Hosts = setValue
    }
    if val, ok := incidentDataResponse["kubernetes_clusters"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.KubernetesClusters = setValue
    }
    if val, ok := incidentDataResponse["kubernetes_resources"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.KubernetesResources = setValue
    }
    if val, ok := incidentDataResponse["kubernetes_containers"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.KubernetesContainers = setValue
    }
    if val, ok := incidentDataResponse["docker_hosts"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.DockerHosts = setValue
    }
    if val, ok := incidentDataResponse["podman_hosts"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.PodmanHosts = setValue
    }
    if val, ok := incidentDataResponse["proxmox_clusters"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.ProxmoxClusters = setValue
    }
    if val, ok := incidentDataResponse["docker_swarm_clusters"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.DockerSwarmClusters = setValue
    }
    if val, ok := incidentDataResponse["ceph_clusters"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.CephClusters = setValue
    }
    if val, ok := incidentDataResponse["docker_resources"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.DockerResources = setValue
    }
    if val, ok := incidentDataResponse["podman_resources"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.PodmanResources = setValue
    }
    if val, ok := incidentDataResponse["services"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Services = setValue
    }
    if val, ok := incidentDataResponse["on_call_duty_policies"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OnCallDutyPolicies = setValue
    }
    if val, ok := incidentDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := incidentDataResponse["current_incident_state_id"].(string); ok {
        data.CurrentIncidentStateId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["incident_severity_id"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["change_monitor_status_to_id"].(string); ok {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["subscriber_notification_status_on_incident_created"].(string); ok {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["subscriber_notification_status_message"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["subscriber_notification_status_on_postmortem_published"].(string); ok {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["subscriber_notification_status_message_on_postmortem_published"].(string); ok {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["should_status_page_subscribers_be_notified_on_incident_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["is_owner_notified_of_resource_creation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["root_cause"].(string); ok {
        data.RootCause = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["postmortem_note"].(string); ok {
        data.PostmortemNote = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["show_postmortem_on_status_page"].(bool); ok {
        data.ShowPostmortemOnStatusPage = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["notify_subscribers_on_postmortem_published"].(bool); ok {
        data.NotifySubscribersOnPostmortemPublished = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["postmortem_posted_at"].(string); ok {
        data.PostmortemPostedAt = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["postmortem_attachments"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.PostmortemAttachments = setValue
    }
    if val, ok := incidentDataResponse["created_state_log"].(string); ok {
        data.CreatedStateLog = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_criteria_id"].(string); ok {
        data.CreatedCriteriaId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_incident_template_id"].(string); ok {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["series_fingerprint"].(string); ok {
        data.SeriesFingerprint = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["series_labels"].(string); ok {
        data.SeriesLabels = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["created_by_probe_id"].(string); ok {
        data.CreatedByProbeId = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["is_created_automatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["remediation_notes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["telemetry_query"].(string); ok {
        data.TelemetryQuery = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["incident_number"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentDataResponse["incident_number_with_prefix"].(string); ok {
        data.IncidentNumberWithPrefix = types.StringValue(val)
    }
    if val, ok := incidentDataResponse["is_visible_on_status_page"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["is_private"].(bool); ok {
        data.IsPrivate = types.BoolValue(val)
    }
    if val, ok := incidentDataResponse["incident_episode_id"].(string); ok {
        data.IncidentEpisodeId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
