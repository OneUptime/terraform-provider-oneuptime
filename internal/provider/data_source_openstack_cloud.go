package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OpenstackCloudDataSource{}

func NewOpenstackCloudDataSource() datasource.DataSource {
    return &OpenstackCloudDataSource{}
}

// NewOpenstackCloudLegacyDataSource registers this data source under the
// name it had before type names kept mixed-case words whole,
// oneuptime_open_stack_cloud. Deprecated.
func NewOpenstackCloudLegacyDataSource() datasource.DataSource {
    return &OpenstackCloudDataSource{isLegacyAlias: true}
}

// OpenstackCloudDataSource defines the data source implementation.
type OpenstackCloudDataSource struct {
    client *Client
    // Registered under the name this data source had before.
    isLegacyAlias bool
}

// OpenstackCloudDataSourceModel describes the data source data model.
type OpenstackCloudDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    RegionCount types.Number `tfsdk:"region_count"`
    AvailabilityZoneCount types.Number `tfsdk:"availability_zone_count"`
    HypervisorCount types.Number `tfsdk:"hypervisor_count"`
    InstanceCount types.Number `tfsdk:"instance_count"`
    ActiveInstanceCount types.Number `tfsdk:"active_instance_count"`
    VolumeCount types.Number `tfsdk:"volume_count"`
    NetworkCount types.Number `tfsdk:"network_count"`
    RouterCount types.Number `tfsdk:"router_count"`
    FloatingIpCount types.Number `tfsdk:"floating_ip_count"`
    ProjectCount types.Number `tfsdk:"project_count"`
    ServiceCount types.Number `tfsdk:"service_count"`
    DownServiceCount types.Number `tfsdk:"down_service_count"`
    VcpuCapacity types.Number `tfsdk:"vcpu_capacity"`
    VcpuAllocated types.Number `tfsdk:"vcpu_allocated"`
    MemoryCapacityBytes types.Number `tfsdk:"memory_capacity_bytes"`
    MemoryAllocatedBytes types.Number `tfsdk:"memory_allocated_bytes"`
    StorageCapacityBytes types.Number `tfsdk:"storage_capacity_bytes"`
    StorageUsedBytes types.Number `tfsdk:"storage_used_bytes"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    CollectionMethod types.String `tfsdk:"collection_method"`
    AuthUrl types.String `tfsdk:"auth_url"`
    EndpointInterface types.String `tfsdk:"endpoint_interface"`
    CollectionRegions types.String `tfsdk:"collection_regions"`
    CredentialKind types.String `tfsdk:"credential_kind"`
    ApplicationCredentialId types.String `tfsdk:"application_credential_id"`
    KeystoneUsername types.String `tfsdk:"keystone_username"`
    UserDomainName types.String `tfsdk:"user_domain_name"`
    ProjectName types.String `tfsdk:"project_name"`
    ProjectDomainName types.String `tfsdk:"project_domain_name"`
    IsCredentialSet types.Bool `tfsdk:"is_credential_set"`
    CredentialsUpdatedAt types.String `tfsdk:"credentials_updated_at"`
    CollectionProbeId types.String `tfsdk:"collection_probe_id"`
    TrustedCertificateFingerprints types.String `tfsdk:"trusted_certificate_fingerprints"`
    CollectionIntervalInMinutes types.Number `tfsdk:"collection_interval_in_minutes"`
    CollectionOptions types.String `tfsdk:"collection_options"`
    InstanceMetricsOptions types.String `tfsdk:"instance_metrics_options"`
    CapacitySnapshot types.String `tfsdk:"capacity_snapshot"`
    SpendSnapshot types.String `tfsdk:"spend_snapshot"`
    UsageReadAt types.String `tfsdk:"usage_read_at"`
    PriceList types.String `tfsdk:"price_list"`
    CollectionStatus types.String `tfsdk:"collection_status"`
    CollectionErrorCode types.String `tfsdk:"collection_error_code"`
    CollectionError types.String `tfsdk:"collection_error"`
    PresentedCertificate types.String `tfsdk:"presented_certificate"`
    CollectionSummary types.String `tfsdk:"collection_summary"`
    NextCollectionAt types.String `tfsdk:"next_collection_at"`
    LastCollectionAt types.String `tfsdk:"last_collection_at"`
    LastSuccessfulCollectionAt types.String `tfsdk:"last_successful_collection_at"`
    CollectionSettingsVersion types.Number `tfsdk:"collection_settings_version"`
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiCommandAllowlist types.String `tfsdk:"ai_command_allowlist"`
    AiAccessLastVerifiedAt types.String `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt types.String `tfsdk:"ai_access_configured_at"`
    AiAgentResetAt types.String `tfsdk:"ai_agent_reset_at"`
}

func (d *OpenstackCloudDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    if d.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_open_stack_cloud"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_openstack_cloud"
}

func (d *OpenstackCloudDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "OpenStack clouds that are being monitored in this project. A cloud is discovered when the OneUptime OpenStack Agent sends its first metrics, or connected with a read-only credential that one of the project's probes collects it with. Look up an existing openstack cloud by `id`, or by any of its other arguments (`name`, `active_instance_count`, `agent_version`, ...): each one set must match, and exactly one openstack cloud may match them all.",

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
                MarkdownDescription: "Name of this OpenStack cloud. This is the join key — it must match the openstack.cloud.name OTel resource attribute stamped by the OneUptime OpenStack Agent.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this OpenStack cloud.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime OpenStack agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When metrics were last received from this cloud.",
                Computed: true,
            },
            "region_count": schema.NumberAttribute{
                MarkdownDescription: "How many regions this cloud reports.",
                Optional: true,
                Computed: true,
            },
            "availability_zone_count": schema.NumberAttribute{
                MarkdownDescription: "How many availability zones this cloud reports.",
                Optional: true,
                Computed: true,
            },
            "hypervisor_count": schema.NumberAttribute{
                MarkdownDescription: "How many hypervisors this cloud reports.",
                Optional: true,
                Computed: true,
            },
            "instance_count": schema.NumberAttribute{
                MarkdownDescription: "How many instances this cloud reports, in every project.",
                Optional: true,
                Computed: true,
            },
            "active_instance_count": schema.NumberAttribute{
                MarkdownDescription: "How many instances are ACTIVE.",
                Optional: true,
                Computed: true,
            },
            "volume_count": schema.NumberAttribute{
                MarkdownDescription: "How many volumes this cloud reports, in every project.",
                Optional: true,
                Computed: true,
            },
            "network_count": schema.NumberAttribute{
                MarkdownDescription: "How many networks this cloud reports, in every project.",
                Optional: true,
                Computed: true,
            },
            "router_count": schema.NumberAttribute{
                MarkdownDescription: "How many routers this cloud reports, in every project.",
                Optional: true,
                Computed: true,
            },
            "floating_ip_count": schema.NumberAttribute{
                MarkdownDescription: "How many floating IPs this cloud reports.",
                Optional: true,
                Computed: true,
            },
            "project_count": schema.NumberAttribute{
                MarkdownDescription: "How many OpenStack projects this cloud reports.",
                Optional: true,
                Computed: true,
            },
            "service_count": schema.NumberAttribute{
                MarkdownDescription: "How many services and agents (nova-compute, neutron agents, cinder-volume and so on) this cloud reports.",
                Optional: true,
                Computed: true,
            },
            "down_service_count": schema.NumberAttribute{
                MarkdownDescription: "How many enabled services and agents are down. A service someone disabled is not counted.",
                Optional: true,
                Computed: true,
            },
            "vcpu_capacity": schema.NumberAttribute{
                MarkdownDescription: "vCPUs the cloud can promise to instances, as Placement counts them: (total - reserved) * allocation ratio.",
                Optional: true,
                Computed: true,
            },
            "vcpu_allocated": schema.NumberAttribute{
                MarkdownDescription: "vCPUs promised to instances. This is what is promised, not what is busy.",
                Optional: true,
                Computed: true,
            },
            "memory_capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Memory the cloud can promise to instances, as Placement counts it.",
                Optional: true,
                Computed: true,
            },
            "memory_allocated_bytes": schema.NumberAttribute{
                MarkdownDescription: "Memory promised to instances. This is what is promised, not what is in use.",
                Optional: true,
                Computed: true,
            },
            "storage_capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Total size of the cloud's Cinder storage pools.",
                Optional: true,
                Computed: true,
            },
            "storage_used_bytes": schema.NumberAttribute{
                MarkdownDescription: "How much of the cloud's Cinder storage pools is used.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this OpenStack cloud archived? Archived OpenStack clouds are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this OpenStack cloud archived?",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this OpenStack cloud. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this OpenStack cloud (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the OpenStack cloud default, then the project's retention settings. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "collection_method": schema.StringAttribute{
                MarkdownDescription: "How this cloud's data reaches OneUptime. Probe: one of your OneUptime probes reads the cloud's APIs with the credential saved here, with no agent to install. Agent: the OneUptime OpenStack Agent you run sends it. Switching to Agent forgets the saved credential.",
                Optional: true,
                Computed: true,
            },
            "auth_url": schema.StringAttribute{
                MarkdownDescription: "The HTTPS address of the cloud's Keystone, which the probe logs in at and reads the service catalog from, such as https://keystone.example.com:5000/v3. Stored as https://host[:port][/path], without the version. Changing it asks for the credential again.",
                Optional: true,
                Computed: true,
            },
            "endpoint_interface": schema.StringAttribute{
                MarkdownDescription: "Which of the catalog's addresses the probe uses for each service: public, internal or admin. Most probes sit where public answers. One in the control plane's own network may want internal.",
                Optional: true,
                Computed: true,
            },
            "collection_regions": schema.StringAttribute{
                MarkdownDescription: "The regions the probe reads, as a JSON array of the names the catalog uses, such as [\"RegionOne\"]. The connection test lists the ones the catalog names. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "credential_kind": schema.StringAttribute{
                MarkdownDescription: "Which kind of credential the probe logs in with: applicationCredential or password. Worked out from what is written with a new secret, so it can be left out.",
                Optional: true,
                Computed: true,
            },
            "application_credential_id": schema.StringAttribute{
                MarkdownDescription: "The ID of the Keystone application credential the probe logs in with, the better of the two kinds because it can be limited to reading and can expire. Its secret goes in Credential Secret. Changing it asks for the secret again.",
                Optional: true,
                Computed: true,
            },
            "keystone_username": schema.StringAttribute{
                MarkdownDescription: "The Keystone user the probe logs in as, when it logs in with a password instead of an application credential. The password goes in Credential Secret. Changing it asks for the password again.",
                Optional: true,
                Computed: true,
            },
            "user_domain_name": schema.StringAttribute{
                MarkdownDescription: "The Keystone domain of that user. Default, when left out.",
                Optional: true,
                Computed: true,
            },
            "project_name": schema.StringAttribute{
                MarkdownDescription: "The project that user logs in to, when it logs in with a password. An application credential already names its own project.",
                Optional: true,
                Computed: true,
            },
            "project_domain_name": schema.StringAttribute{
                MarkdownDescription: "The Keystone domain of that project. Default, when left out.",
                Optional: true,
                Computed: true,
            },
            "is_credential_set": schema.BoolAttribute{
                MarkdownDescription: "Whether a credential is saved. Set by the server from the secret itself.",
                Optional: true,
                Computed: true,
            },
            "credentials_updated_at": schema.StringAttribute{
                MarkdownDescription: "When the credential was last changed. Set by the server.",
                Computed: true,
            },
            "collection_probe_id": schema.StringAttribute{
                MarkdownDescription: "ID of the OneUptime probe that collects this cloud. Changing it asks for the credential again. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_probe"),
            },
            "trusted_certificate_fingerprints": schema.StringAttribute{
                MarkdownDescription: "The SHA-256 fingerprints of the certificates the probe accepts from this cloud's APIs, as a JSON array, for APIs whose certificates are not from a public authority. More than one, because each service can have a host and a certificate of its own. Empty: every certificate must be signed by an authority the probe's machine trusts. Verification is never skipped. Trusting another certificate asks for the credential again. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "collection_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often the probe collects this cloud, in whole minutes from 1 to 60. Two minutes is the OpenStack agent's own default. Collect large clouds less often to go easier on their APIs.",
                Optional: true,
                Computed: true,
            },
            "collection_options": schema.StringAttribute{
                MarkdownDescription: "What the probe reads beyond the agent's defaults, as JSON. collectPorts reads ports, scoped by portNetworkIds and portProjectIds and paged. collectProjectLimits reads each project's quotas, four calls per project. Both are off unless asked for. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "instance_metrics_options": schema.StringAttribute{
                MarkdownDescription: "Whether real per-instance CPU, memory, disk and network series are kept for this cloud, as JSON: mode (automatic, all, selected or off), and for selected, the projectIds and aggregates. Per-instance series multiply fast, so automatic stops keeping them past a stated number of instances. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "capacity_snapshot": schema.StringAttribute{
                MarkdownDescription: "What this cloud's capacity added up to when a worker last read its inventory, as JSON: what is allocated and what can be, by zone, aggregate, region and cloud, how many more of each flavor fit, when each thing fills up, and what no instance explains. Written by the server. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "spend_snapshot": schema.StringAttribute{
                MarkdownDescription: "What this cloud's usage added up to when a worker last read its inventory, as JSON: this month against last, the largest domains and flavors, the top growers, and what could be reclaimed. Estimates, never an invoice. Written by the server. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "usage_read_at": schema.StringAttribute{
                MarkdownDescription: "The time this cloud's usage has been counted up to. Written by the server, and how a time is never counted twice.",
                Computed: true,
            },
            "price_list": schema.StringAttribute{
                MarkdownDescription: "What this cloud's operator says things cost, per hour, in one currency, as JSON. Optional: usage works without it, and a price of zero means no price has been set. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "collection_status": schema.StringAttribute{
                MarkdownDescription: "Where probe collection stands: Pending (the probe has not tried the current settings yet), Succeeded, PartiallySucceeded (some of the cloud could not be read) or Failed. Written by the server and the probe.",
                Optional: true,
                Computed: true,
            },
            "collection_error_code": schema.StringAttribute{
                MarkdownDescription: "Why the last collection failed or was partial, as a code the dashboard turns into a next step (InvalidCredential, UntrustedCertificate, ConnectionTimedOut, ...). Written by the probe.",
                Optional: true,
                Computed: true,
            },
            "collection_error": schema.StringAttribute{
                MarkdownDescription: "What went wrong in the last collection, in words. Written by the probe.",
                Optional: true,
                Computed: true,
            },
            "presented_certificate": schema.StringAttribute{
                MarkdownDescription: "The certificate an API presented when the last collection did not trust it: its SHA-256 fingerprint, subject, issuer, validity and the host that showed it, so a person can check it and trust exactly that certificate. Written by the probe. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "collection_summary": schema.StringAttribute{
                MarkdownDescription: "What the last collection found: the regions, how each service answered, how many hypervisors, instances, networks, volumes and projects it read, and what it left out. Written by the probe. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "next_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe collects this cloud next. Written by the server.",
                Computed: true,
            },
            "last_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe last tried to collect this cloud, whether it worked or not. Written by the server.",
                Computed: true,
            },
            "last_successful_collection_at": schema.StringAttribute{
                MarkdownDescription: "When a collection of this cloud last worked, in full or in part. Written by the server.",
                Computed: true,
            },
            "collection_settings_version": schema.NumberAttribute{
                MarkdownDescription: "Counts changes to the probe collection settings, so a report from a collection that started before a change is not shown as the current status. Written by the server.",
                Optional: true,
                Computed: true,
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (openstack server, hypervisor, service and agent reads, each one call to the cloud's API) on this OpenStack cloud, through its OpenStack AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the OpenStack cloud can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI only reads an OpenStack cloud today: it has no command that changes one, so this stays Disabled and any other mode is refused. Disabled: AI never proposes or runs a change on this OpenStack cloud. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the OpenStack cloud's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the OpenStack cloud can lower the mode; raising it needs Project Owner, Project Admin or Edit AI Rule.",
                Optional: true,
                Computed: true,
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI only reads an OpenStack cloud today, and a read needs no approval, so no pattern is valid here yet and this stays empty. Optional JSON array of command patterns that Automatic mode may run on this OpenStack cloud without approval even though they are riskier changes. Each pattern is one command line for this OpenStack cloud's agent (openstack) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this OpenStack cloud is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit AI Rule; anyone who may edit the OpenStack cloud can remove patterns or clear the list. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this OpenStack cloud through its OpenStack AI agent. Set by the server.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this OpenStack cloud, kept until the next successful command. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this OpenStack cloud was first configured by anyone saving an AI access setting. Set by the server; never cleared, so an OpenStack AI agent that registers later never overwrites a setting an operator chose.",
                Computed: true,
            },
            "ai_agent_reset_at": schema.StringAttribute{
                MarkdownDescription: "When this OpenStack cloud's AI agent was last reset: its key stopped working at once, and it registers again by itself. Set by the server.",
                Computed: true,
            },
        },
    }
    if d.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_open_stack_cloud has been renamed to oneuptime_openstack_cloud. Use the new name; the old one keeps working until then."
    }
}

func (d *OpenstackCloudDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OpenstackCloudDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OpenstackCloudDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.OtelCollectorStatus.IsNull() && !data.OtelCollectorStatus.IsUnknown() {
        filters["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
        filterNames = append(filterNames, "otel_collector_status = "+fmt.Sprintf("%q", data.OtelCollectorStatus.ValueString()))
    }
    if !data.AgentVersion.IsNull() && !data.AgentVersion.IsUnknown() {
        filters["agentVersion"] = data.AgentVersion.ValueString()
        filterNames = append(filterNames, "agent_version = "+fmt.Sprintf("%q", data.AgentVersion.ValueString()))
    }
    if !data.RegionCount.IsNull() && !data.RegionCount.IsUnknown() {
        filters["regionCount"] = lookupNumber(data.RegionCount)
        filterNames = append(filterNames, "region_count = "+data.RegionCount.ValueBigFloat().String())
    }
    if !data.AvailabilityZoneCount.IsNull() && !data.AvailabilityZoneCount.IsUnknown() {
        filters["availabilityZoneCount"] = lookupNumber(data.AvailabilityZoneCount)
        filterNames = append(filterNames, "availability_zone_count = "+data.AvailabilityZoneCount.ValueBigFloat().String())
    }
    if !data.HypervisorCount.IsNull() && !data.HypervisorCount.IsUnknown() {
        filters["hypervisorCount"] = lookupNumber(data.HypervisorCount)
        filterNames = append(filterNames, "hypervisor_count = "+data.HypervisorCount.ValueBigFloat().String())
    }
    if !data.InstanceCount.IsNull() && !data.InstanceCount.IsUnknown() {
        filters["instanceCount"] = lookupNumber(data.InstanceCount)
        filterNames = append(filterNames, "instance_count = "+data.InstanceCount.ValueBigFloat().String())
    }
    if !data.ActiveInstanceCount.IsNull() && !data.ActiveInstanceCount.IsUnknown() {
        filters["activeInstanceCount"] = lookupNumber(data.ActiveInstanceCount)
        filterNames = append(filterNames, "active_instance_count = "+data.ActiveInstanceCount.ValueBigFloat().String())
    }
    if !data.VolumeCount.IsNull() && !data.VolumeCount.IsUnknown() {
        filters["volumeCount"] = lookupNumber(data.VolumeCount)
        filterNames = append(filterNames, "volume_count = "+data.VolumeCount.ValueBigFloat().String())
    }
    if !data.NetworkCount.IsNull() && !data.NetworkCount.IsUnknown() {
        filters["networkCount"] = lookupNumber(data.NetworkCount)
        filterNames = append(filterNames, "network_count = "+data.NetworkCount.ValueBigFloat().String())
    }
    if !data.RouterCount.IsNull() && !data.RouterCount.IsUnknown() {
        filters["routerCount"] = lookupNumber(data.RouterCount)
        filterNames = append(filterNames, "router_count = "+data.RouterCount.ValueBigFloat().String())
    }
    if !data.FloatingIpCount.IsNull() && !data.FloatingIpCount.IsUnknown() {
        filters["floatingIpCount"] = lookupNumber(data.FloatingIpCount)
        filterNames = append(filterNames, "floating_ip_count = "+data.FloatingIpCount.ValueBigFloat().String())
    }
    if !data.ProjectCount.IsNull() && !data.ProjectCount.IsUnknown() {
        filters["projectCount"] = lookupNumber(data.ProjectCount)
        filterNames = append(filterNames, "project_count = "+data.ProjectCount.ValueBigFloat().String())
    }
    if !data.ServiceCount.IsNull() && !data.ServiceCount.IsUnknown() {
        filters["serviceCount"] = lookupNumber(data.ServiceCount)
        filterNames = append(filterNames, "service_count = "+data.ServiceCount.ValueBigFloat().String())
    }
    if !data.DownServiceCount.IsNull() && !data.DownServiceCount.IsUnknown() {
        filters["downServiceCount"] = lookupNumber(data.DownServiceCount)
        filterNames = append(filterNames, "down_service_count = "+data.DownServiceCount.ValueBigFloat().String())
    }
    if !data.VcpuCapacity.IsNull() && !data.VcpuCapacity.IsUnknown() {
        filters["vcpuCapacity"] = lookupNumber(data.VcpuCapacity)
        filterNames = append(filterNames, "vcpu_capacity = "+data.VcpuCapacity.ValueBigFloat().String())
    }
    if !data.VcpuAllocated.IsNull() && !data.VcpuAllocated.IsUnknown() {
        filters["vcpuAllocated"] = lookupNumber(data.VcpuAllocated)
        filterNames = append(filterNames, "vcpu_allocated = "+data.VcpuAllocated.ValueBigFloat().String())
    }
    if !data.MemoryCapacityBytes.IsNull() && !data.MemoryCapacityBytes.IsUnknown() {
        filters["memoryCapacityBytes"] = lookupNumber(data.MemoryCapacityBytes)
        filterNames = append(filterNames, "memory_capacity_bytes = "+data.MemoryCapacityBytes.ValueBigFloat().String())
    }
    if !data.MemoryAllocatedBytes.IsNull() && !data.MemoryAllocatedBytes.IsUnknown() {
        filters["memoryAllocatedBytes"] = lookupNumber(data.MemoryAllocatedBytes)
        filterNames = append(filterNames, "memory_allocated_bytes = "+data.MemoryAllocatedBytes.ValueBigFloat().String())
    }
    if !data.StorageCapacityBytes.IsNull() && !data.StorageCapacityBytes.IsUnknown() {
        filters["storageCapacityBytes"] = lookupNumber(data.StorageCapacityBytes)
        filterNames = append(filterNames, "storage_capacity_bytes = "+data.StorageCapacityBytes.ValueBigFloat().String())
    }
    if !data.StorageUsedBytes.IsNull() && !data.StorageUsedBytes.IsUnknown() {
        filters["storageUsedBytes"] = lookupNumber(data.StorageUsedBytes)
        filterNames = append(filterNames, "storage_used_bytes = "+data.StorageUsedBytes.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        filters["retainTelemetryDataForDays"] = lookupNumber(data.RetainTelemetryDataForDays)
        filterNames = append(filterNames, "retain_telemetry_data_for_days = "+data.RetainTelemetryDataForDays.ValueBigFloat().String())
    }
    if !data.CollectionMethod.IsNull() && !data.CollectionMethod.IsUnknown() {
        filters["collectionMethod"] = data.CollectionMethod.ValueString()
        filterNames = append(filterNames, "collection_method = "+fmt.Sprintf("%q", data.CollectionMethod.ValueString()))
    }
    if !data.AuthUrl.IsNull() && !data.AuthUrl.IsUnknown() {
        filters["authUrl"] = data.AuthUrl.ValueString()
        filterNames = append(filterNames, "auth_url = "+fmt.Sprintf("%q", data.AuthUrl.ValueString()))
    }
    if !data.EndpointInterface.IsNull() && !data.EndpointInterface.IsUnknown() {
        filters["endpointInterface"] = data.EndpointInterface.ValueString()
        filterNames = append(filterNames, "endpoint_interface = "+fmt.Sprintf("%q", data.EndpointInterface.ValueString()))
    }
    if !data.CredentialKind.IsNull() && !data.CredentialKind.IsUnknown() {
        filters["credentialKind"] = data.CredentialKind.ValueString()
        filterNames = append(filterNames, "credential_kind = "+fmt.Sprintf("%q", data.CredentialKind.ValueString()))
    }
    if !data.ApplicationCredentialId.IsNull() && !data.ApplicationCredentialId.IsUnknown() {
        filters["applicationCredentialId"] = data.ApplicationCredentialId.ValueString()
        filterNames = append(filterNames, "application_credential_id = "+fmt.Sprintf("%q", data.ApplicationCredentialId.ValueString()))
    }
    if !data.KeystoneUsername.IsNull() && !data.KeystoneUsername.IsUnknown() {
        filters["keystoneUsername"] = data.KeystoneUsername.ValueString()
        filterNames = append(filterNames, "keystone_username = "+fmt.Sprintf("%q", data.KeystoneUsername.ValueString()))
    }
    if !data.UserDomainName.IsNull() && !data.UserDomainName.IsUnknown() {
        filters["userDomainName"] = data.UserDomainName.ValueString()
        filterNames = append(filterNames, "user_domain_name = "+fmt.Sprintf("%q", data.UserDomainName.ValueString()))
    }
    if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
        filters["projectName"] = data.ProjectName.ValueString()
        filterNames = append(filterNames, "project_name = "+fmt.Sprintf("%q", data.ProjectName.ValueString()))
    }
    if !data.ProjectDomainName.IsNull() && !data.ProjectDomainName.IsUnknown() {
        filters["projectDomainName"] = data.ProjectDomainName.ValueString()
        filterNames = append(filterNames, "project_domain_name = "+fmt.Sprintf("%q", data.ProjectDomainName.ValueString()))
    }
    if !data.IsCredentialSet.IsNull() && !data.IsCredentialSet.IsUnknown() {
        filters["isCredentialSet"] = data.IsCredentialSet.ValueBool()
        filterNames = append(filterNames, "is_credential_set = "+fmt.Sprintf("%t", data.IsCredentialSet.ValueBool()))
    }
    if !data.CollectionProbeId.IsNull() && !data.CollectionProbeId.IsUnknown() {
        filters["collectionProbeId"] = data.CollectionProbeId.ValueString()
        filterNames = append(filterNames, "collection_probe_id = "+fmt.Sprintf("%q", data.CollectionProbeId.ValueString()))
    }
    if !data.CollectionIntervalInMinutes.IsNull() && !data.CollectionIntervalInMinutes.IsUnknown() {
        filters["collectionIntervalInMinutes"] = lookupNumber(data.CollectionIntervalInMinutes)
        filterNames = append(filterNames, "collection_interval_in_minutes = "+data.CollectionIntervalInMinutes.ValueBigFloat().String())
    }
    if !data.CollectionStatus.IsNull() && !data.CollectionStatus.IsUnknown() {
        filters["collectionStatus"] = data.CollectionStatus.ValueString()
        filterNames = append(filterNames, "collection_status = "+fmt.Sprintf("%q", data.CollectionStatus.ValueString()))
    }
    if !data.CollectionErrorCode.IsNull() && !data.CollectionErrorCode.IsUnknown() {
        filters["collectionErrorCode"] = data.CollectionErrorCode.ValueString()
        filterNames = append(filterNames, "collection_error_code = "+fmt.Sprintf("%q", data.CollectionErrorCode.ValueString()))
    }
    if !data.CollectionError.IsNull() && !data.CollectionError.IsUnknown() {
        filters["collectionError"] = data.CollectionError.ValueString()
        filterNames = append(filterNames, "collection_error = "+fmt.Sprintf("%q", data.CollectionError.ValueString()))
    }
    if !data.CollectionSettingsVersion.IsNull() && !data.CollectionSettingsVersion.IsUnknown() {
        filters["collectionSettingsVersion"] = lookupNumber(data.CollectionSettingsVersion)
        filterNames = append(filterNames, "collection_settings_version = "+data.CollectionSettingsVersion.ValueBigFloat().String())
    }
    if !data.IsAiInvestigationEnabled.IsNull() && !data.IsAiInvestigationEnabled.IsUnknown() {
        filters["isAiInvestigationEnabled"] = data.IsAiInvestigationEnabled.ValueBool()
        filterNames = append(filterNames, "is_ai_investigation_enabled = "+fmt.Sprintf("%t", data.IsAiInvestigationEnabled.ValueBool()))
    }
    if !data.AiRemediationMode.IsNull() && !data.AiRemediationMode.IsUnknown() {
        filters["aiRemediationMode"] = data.AiRemediationMode.ValueString()
        filterNames = append(filterNames, "ai_remediation_mode = "+fmt.Sprintf("%q", data.AiRemediationMode.ValueString()))
    }
    if !data.AiAccessLastError.IsNull() && !data.AiAccessLastError.IsUnknown() {
        filters["aiAccessLastError"] = data.AiAccessLastError.ValueString()
        filterNames = append(filterNames, "ai_access_last_error = "+fmt.Sprintf("%q", data.AiAccessLastError.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the openstack cloud up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the openstack cloud up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "regionCount": true,
        "availabilityZoneCount": true,
        "hypervisorCount": true,
        "instanceCount": true,
        "activeInstanceCount": true,
        "volumeCount": true,
        "networkCount": true,
        "routerCount": true,
        "floatingIpCount": true,
        "projectCount": true,
        "serviceCount": true,
        "downServiceCount": true,
        "vcpuCapacity": true,
        "vcpuAllocated": true,
        "memoryCapacityBytes": true,
        "memoryAllocatedBytes": true,
        "storageCapacityBytes": true,
        "storageUsedBytes": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "collectionMethod": true,
        "authUrl": true,
        "endpointInterface": true,
        "collectionRegions": true,
        "credentialKind": true,
        "applicationCredentialId": true,
        "keystoneUsername": true,
        "userDomainName": true,
        "projectName": true,
        "projectDomainName": true,
        "isCredentialSet": true,
        "credentialsUpdatedAt": true,
        "collectionProbeId": true,
        "trustedCertificateFingerprints": true,
        "collectionIntervalInMinutes": true,
        "collectionOptions": true,
        "instanceMetricsOptions": true,
        "capacitySnapshot": true,
        "spendSnapshot": true,
        "usageReadAt": true,
        "priceList": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "aiAgentResetAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/openstack-cloud/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read openstack_cloud, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No openstack cloud found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read openstack_cloud: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/openstack-cloud/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list openstack_cloud, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list openstack_cloud: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No openstack cloud matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one openstack cloud matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for openstack_cloud.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
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
    } else if arr, ok := item["slug"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["otelCollectorStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := item["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := item["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := item["agentVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := item["agentVersion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := item["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if arr, ok := item["lastSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    data.RegionCount = numberFromAPI(data.RegionCount, item["regionCount"])
    data.AvailabilityZoneCount = numberFromAPI(data.AvailabilityZoneCount, item["availabilityZoneCount"])
    data.HypervisorCount = numberFromAPI(data.HypervisorCount, item["hypervisorCount"])
    data.InstanceCount = numberFromAPI(data.InstanceCount, item["instanceCount"])
    data.ActiveInstanceCount = numberFromAPI(data.ActiveInstanceCount, item["activeInstanceCount"])
    data.VolumeCount = numberFromAPI(data.VolumeCount, item["volumeCount"])
    data.NetworkCount = numberFromAPI(data.NetworkCount, item["networkCount"])
    data.RouterCount = numberFromAPI(data.RouterCount, item["routerCount"])
    data.FloatingIpCount = numberFromAPI(data.FloatingIpCount, item["floatingIpCount"])
    data.ProjectCount = numberFromAPI(data.ProjectCount, item["projectCount"])
    data.ServiceCount = numberFromAPI(data.ServiceCount, item["serviceCount"])
    data.DownServiceCount = numberFromAPI(data.DownServiceCount, item["downServiceCount"])
    data.VcpuCapacity = numberFromAPI(data.VcpuCapacity, item["vcpuCapacity"])
    data.VcpuAllocated = numberFromAPI(data.VcpuAllocated, item["vcpuAllocated"])
    data.MemoryCapacityBytes = numberFromAPI(data.MemoryCapacityBytes, item["memoryCapacityBytes"])
    data.MemoryAllocatedBytes = numberFromAPI(data.MemoryAllocatedBytes, item["memoryAllocatedBytes"])
    data.StorageCapacityBytes = numberFromAPI(data.StorageCapacityBytes, item["storageCapacityBytes"])
    data.StorageUsedBytes = numberFromAPI(data.StorageUsedBytes, item["storageUsedBytes"])
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if arr, ok := item["archivedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := item["archivedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
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
    data.RetainTelemetryDataForDays = numberFromAPI(data.RetainTelemetryDataForDays, item["retainTelemetryDataForDays"])
    if obj, ok := item["telemetryRetentionConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryRetentionConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if arr, ok := item["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if val, ok := item["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    } else {
        data.TelemetryRetentionConfig = types.StringNull()
    }
    if obj, ok := item["collectionMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if arr, ok := item["collectionMethod"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if val, ok := item["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := item["authUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AuthUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if arr, ok := item["authUrl"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if val, ok := item["authUrl"].(string); ok {
        data.AuthUrl = types.StringValue(val)
    } else {
        data.AuthUrl = types.StringNull()
    }
    if obj, ok := item["endpointInterface"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if arr, ok := item["endpointInterface"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if val, ok := item["endpointInterface"].(string); ok {
        data.EndpointInterface = types.StringValue(val)
    } else {
        data.EndpointInterface = types.StringNull()
    }
    if obj, ok := item["collectionRegions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionRegions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionRegions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionRegions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionRegions = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionRegions = types.StringNull()
        }
    } else if arr, ok := item["collectionRegions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionRegions = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionRegions = types.StringNull()
        }
    } else if val, ok := item["collectionRegions"].(string); ok {
        data.CollectionRegions = types.StringValue(val)
    } else {
        data.CollectionRegions = types.StringNull()
    }
    if obj, ok := item["credentialKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CredentialKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if arr, ok := item["credentialKind"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if val, ok := item["credentialKind"].(string); ok {
        data.CredentialKind = types.StringValue(val)
    } else {
        data.CredentialKind = types.StringNull()
    }
    if obj, ok := item["applicationCredentialId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if arr, ok := item["applicationCredentialId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if val, ok := item["applicationCredentialId"].(string); ok {
        data.ApplicationCredentialId = types.StringValue(val)
    } else {
        data.ApplicationCredentialId = types.StringNull()
    }
    if obj, ok := item["keystoneUsername"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if arr, ok := item["keystoneUsername"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if val, ok := item["keystoneUsername"].(string); ok {
        data.KeystoneUsername = types.StringValue(val)
    } else {
        data.KeystoneUsername = types.StringNull()
    }
    if obj, ok := item["userDomainName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if arr, ok := item["userDomainName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if val, ok := item["userDomainName"].(string); ok {
        data.UserDomainName = types.StringValue(val)
    } else {
        data.UserDomainName = types.StringNull()
    }
    if obj, ok := item["projectName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if arr, ok := item["projectName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if val, ok := item["projectName"].(string); ok {
        data.ProjectName = types.StringValue(val)
    } else {
        data.ProjectName = types.StringNull()
    }
    if obj, ok := item["projectDomainName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if arr, ok := item["projectDomainName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if val, ok := item["projectDomainName"].(string); ok {
        data.ProjectDomainName = types.StringValue(val)
    } else {
        data.ProjectDomainName = types.StringNull()
    }
    if val, ok := item["isCredentialSet"].(bool); ok {
        data.IsCredentialSet = types.BoolValue(val)
    } else {
        data.IsCredentialSet = types.BoolNull()
    }
    if obj, ok := item["credentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CredentialsUpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CredentialsUpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CredentialsUpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CredentialsUpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialsUpdatedAt = types.StringNull()
        }
    } else if arr, ok := item["credentialsUpdatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CredentialsUpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialsUpdatedAt = types.StringNull()
        }
    } else if val, ok := item["credentialsUpdatedAt"].(string); ok {
        data.CredentialsUpdatedAt = types.StringValue(val)
    } else {
        data.CredentialsUpdatedAt = types.StringNull()
    }
    if obj, ok := item["collectionProbeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if arr, ok := item["collectionProbeId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if val, ok := item["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := item["trustedCertificateFingerprints"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprints = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrustedCertificateFingerprints = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrustedCertificateFingerprints = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrustedCertificateFingerprints = types.StringValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = types.StringNull()
        }
    } else if arr, ok := item["trustedCertificateFingerprints"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TrustedCertificateFingerprints = types.StringValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = types.StringNull()
        }
    } else if val, ok := item["trustedCertificateFingerprints"].(string); ok {
        data.TrustedCertificateFingerprints = types.StringValue(val)
    } else {
        data.TrustedCertificateFingerprints = types.StringNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, item["collectionIntervalInMinutes"])
    if obj, ok := item["collectionOptions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionOptions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionOptions = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionOptions = types.StringNull()
        }
    } else if arr, ok := item["collectionOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionOptions = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionOptions = types.StringNull()
        }
    } else if val, ok := item["collectionOptions"].(string); ok {
        data.CollectionOptions = types.StringValue(val)
    } else {
        data.CollectionOptions = types.StringNull()
    }
    if obj, ok := item["instanceMetricsOptions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InstanceMetricsOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InstanceMetricsOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InstanceMetricsOptions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InstanceMetricsOptions = types.StringValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = types.StringNull()
        }
    } else if arr, ok := item["instanceMetricsOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.InstanceMetricsOptions = types.StringValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = types.StringNull()
        }
    } else if val, ok := item["instanceMetricsOptions"].(string); ok {
        data.InstanceMetricsOptions = types.StringValue(val)
    } else {
        data.InstanceMetricsOptions = types.StringNull()
    }
    if obj, ok := item["capacitySnapshot"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CapacitySnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CapacitySnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CapacitySnapshot = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CapacitySnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = types.StringNull()
        }
    } else if arr, ok := item["capacitySnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CapacitySnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = types.StringNull()
        }
    } else if val, ok := item["capacitySnapshot"].(string); ok {
        data.CapacitySnapshot = types.StringValue(val)
    } else {
        data.CapacitySnapshot = types.StringNull()
    }
    if obj, ok := item["spendSnapshot"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpendSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SpendSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SpendSnapshot = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SpendSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = types.StringNull()
        }
    } else if arr, ok := item["spendSnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SpendSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = types.StringNull()
        }
    } else if val, ok := item["spendSnapshot"].(string); ok {
        data.SpendSnapshot = types.StringValue(val)
    } else {
        data.SpendSnapshot = types.StringNull()
    }
    if obj, ok := item["usageReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UsageReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UsageReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UsageReadAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UsageReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.UsageReadAt = types.StringNull()
        }
    } else if arr, ok := item["usageReadAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UsageReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.UsageReadAt = types.StringNull()
        }
    } else if val, ok := item["usageReadAt"].(string); ok {
        data.UsageReadAt = types.StringValue(val)
    } else {
        data.UsageReadAt = types.StringNull()
    }
    if obj, ok := item["priceList"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PriceList = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PriceList = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PriceList = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PriceList = types.StringValue(string(jsonBytes))
        } else {
            data.PriceList = types.StringNull()
        }
    } else if arr, ok := item["priceList"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PriceList = types.StringValue(string(jsonBytes))
        } else {
            data.PriceList = types.StringNull()
        }
    } else if val, ok := item["priceList"].(string); ok {
        data.PriceList = types.StringValue(val)
    } else {
        data.PriceList = types.StringNull()
    }
    if obj, ok := item["collectionStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if arr, ok := item["collectionStatus"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if val, ok := item["collectionStatus"].(string); ok {
        data.CollectionStatus = types.StringValue(val)
    } else {
        data.CollectionStatus = types.StringNull()
    }
    if obj, ok := item["collectionErrorCode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if arr, ok := item["collectionErrorCode"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if val, ok := item["collectionErrorCode"].(string); ok {
        data.CollectionErrorCode = types.StringValue(val)
    } else {
        data.CollectionErrorCode = types.StringNull()
    }
    if obj, ok := item["collectionError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if arr, ok := item["collectionError"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if val, ok := item["collectionError"].(string); ok {
        data.CollectionError = types.StringValue(val)
    } else {
        data.CollectionError = types.StringNull()
    }
    if obj, ok := item["presentedCertificate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PresentedCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PresentedCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PresentedCertificate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PresentedCertificate = types.StringValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = types.StringNull()
        }
    } else if arr, ok := item["presentedCertificate"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PresentedCertificate = types.StringValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = types.StringNull()
        }
    } else if val, ok := item["presentedCertificate"].(string); ok {
        data.PresentedCertificate = types.StringValue(val)
    } else {
        data.PresentedCertificate = types.StringNull()
    }
    if obj, ok := item["collectionSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionSummary = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionSummary = types.StringNull()
        }
    } else if arr, ok := item["collectionSummary"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CollectionSummary = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionSummary = types.StringNull()
        }
    } else if val, ok := item["collectionSummary"].(string); ok {
        data.CollectionSummary = types.StringValue(val)
    } else {
        data.CollectionSummary = types.StringNull()
    }
    if obj, ok := item["nextCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextCollectionAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextCollectionAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextCollectionAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextCollectionAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextCollectionAt = types.StringNull()
        }
    } else if arr, ok := item["nextCollectionAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.NextCollectionAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextCollectionAt = types.StringNull()
        }
    } else if val, ok := item["nextCollectionAt"].(string); ok {
        data.NextCollectionAt = types.StringValue(val)
    } else {
        data.NextCollectionAt = types.StringNull()
    }
    if obj, ok := item["lastCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastCollectionAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastCollectionAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastCollectionAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastCollectionAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastCollectionAt = types.StringNull()
        }
    } else if arr, ok := item["lastCollectionAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastCollectionAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastCollectionAt = types.StringNull()
        }
    } else if val, ok := item["lastCollectionAt"].(string); ok {
        data.LastCollectionAt = types.StringValue(val)
    } else {
        data.LastCollectionAt = types.StringNull()
    }
    if obj, ok := item["lastSuccessfulCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSuccessfulCollectionAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSuccessfulCollectionAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSuccessfulCollectionAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSuccessfulCollectionAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSuccessfulCollectionAt = types.StringNull()
        }
    } else if arr, ok := item["lastSuccessfulCollectionAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSuccessfulCollectionAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSuccessfulCollectionAt = types.StringNull()
        }
    } else if val, ok := item["lastSuccessfulCollectionAt"].(string); ok {
        data.LastSuccessfulCollectionAt = types.StringValue(val)
    } else {
        data.LastSuccessfulCollectionAt = types.StringNull()
    }
    data.CollectionSettingsVersion = numberFromAPI(data.CollectionSettingsVersion, item["collectionSettingsVersion"])
    if val, ok := item["isAiInvestigationEnabled"].(bool); ok {
        data.IsAiInvestigationEnabled = types.BoolValue(val)
    } else {
        data.IsAiInvestigationEnabled = types.BoolNull()
    }
    if obj, ok := item["aiRemediationMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if arr, ok := item["aiRemediationMode"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if val, ok := item["aiRemediationMode"].(string); ok {
        data.AiRemediationMode = types.StringValue(val)
    } else {
        data.AiRemediationMode = types.StringNull()
    }
    if obj, ok := item["aiCommandAllowlist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiCommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiCommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiCommandAllowlist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiCommandAllowlist = types.StringValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = types.StringNull()
        }
    } else if arr, ok := item["aiCommandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiCommandAllowlist = types.StringValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = types.StringNull()
        }
    } else if val, ok := item["aiCommandAllowlist"].(string); ok {
        data.AiCommandAllowlist = types.StringValue(val)
    } else {
        data.AiCommandAllowlist = types.StringNull()
    }
    if obj, ok := item["aiAccessLastVerifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastVerifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessLastVerifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessLastVerifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessLastVerifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastVerifiedAt = types.StringNull()
        }
    } else if arr, ok := item["aiAccessLastVerifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiAccessLastVerifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastVerifiedAt = types.StringNull()
        }
    } else if val, ok := item["aiAccessLastVerifiedAt"].(string); ok {
        data.AiAccessLastVerifiedAt = types.StringValue(val)
    } else {
        data.AiAccessLastVerifiedAt = types.StringNull()
    }
    if obj, ok := item["aiAccessLastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if arr, ok := item["aiAccessLastError"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if val, ok := item["aiAccessLastError"].(string); ok {
        data.AiAccessLastError = types.StringValue(val)
    } else {
        data.AiAccessLastError = types.StringNull()
    }
    if obj, ok := item["aiAccessConfiguredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessConfiguredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessConfiguredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessConfiguredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessConfiguredAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessConfiguredAt = types.StringNull()
        }
    } else if arr, ok := item["aiAccessConfiguredAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiAccessConfiguredAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessConfiguredAt = types.StringNull()
        }
    } else if val, ok := item["aiAccessConfiguredAt"].(string); ok {
        data.AiAccessConfiguredAt = types.StringValue(val)
    } else {
        data.AiAccessConfiguredAt = types.StringNull()
    }
    if obj, ok := item["aiAgentResetAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAgentResetAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAgentResetAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAgentResetAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAgentResetAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAgentResetAt = types.StringNull()
        }
    } else if arr, ok := item["aiAgentResetAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiAgentResetAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAgentResetAt = types.StringNull()
        }
    } else if val, ok := item["aiAgentResetAt"].(string); ok {
        data.AiAgentResetAt = types.StringValue(val)
    } else {
        data.AiAgentResetAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
