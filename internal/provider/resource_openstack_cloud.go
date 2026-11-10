package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OpenstackCloudResource{}
var _ resource.ResourceWithImportState = &OpenstackCloudResource{}
var _ resource.ResourceWithModifyPlan = &OpenstackCloudResource{}
var _ resource.ResourceWithMoveState = &OpenstackCloudResource{}

func NewOpenstackCloudResource() resource.Resource {
    return &OpenstackCloudResource{}
}

// NewOpenstackCloudLegacyResource registers this resource under the name it
// had before resource type names kept mixed-case words whole,
// oneuptime_open_stack_cloud, so configurations that use it keep
// working. Deprecated: plans that use it say so.
func NewOpenstackCloudLegacyResource() resource.Resource {
    return &OpenstackCloudResource{isLegacyAlias: true}
}

// OpenstackCloudResource defines the resource implementation.
type OpenstackCloudResource struct {
    client *Client
    // Registered under the name this resource had before (see NewOpenstackCloudLegacyResource).
    isLegacyAlias bool
}

// OpenstackCloudResourceModel describes the resource data model.
type OpenstackCloudResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig JSONSubsetValue `tfsdk:"telemetry_retention_config"`
    CollectionMethod types.String `tfsdk:"collection_method"`
    AuthUrl types.String `tfsdk:"auth_url"`
    EndpointInterface types.String `tfsdk:"endpoint_interface"`
    CollectionRegions JSONSubsetValue `tfsdk:"collection_regions"`
    CredentialKind types.String `tfsdk:"credential_kind"`
    ApplicationCredentialId types.String `tfsdk:"application_credential_id"`
    KeystoneUsername types.String `tfsdk:"keystone_username"`
    UserDomainName types.String `tfsdk:"user_domain_name"`
    ProjectName types.String `tfsdk:"project_name"`
    ProjectDomainName types.String `tfsdk:"project_domain_name"`
    CredentialSecret types.String `tfsdk:"credential_secret"`
    CollectionProbeId types.String `tfsdk:"collection_probe_id"`
    TrustedCertificateFingerprints JSONSubsetValue `tfsdk:"trusted_certificate_fingerprints"`
    CollectionIntervalInMinutes types.Number `tfsdk:"collection_interval_in_minutes"`
    CollectionOptions JSONSubsetValue `tfsdk:"collection_options"`
    InstanceMetricsOptions JSONSubsetValue `tfsdk:"instance_metrics_options"`
    PriceList JSONSubsetValue `tfsdk:"price_list"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
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
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiCommandAllowlist JSONSubsetValue `tfsdk:"ai_command_allowlist"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ArchivedAt RFC3339Value `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    IsCredentialSet types.Bool `tfsdk:"is_credential_set"`
    CredentialsUpdatedAt RFC3339Value `tfsdk:"credentials_updated_at"`
    CapacitySnapshot JSONSubsetValue `tfsdk:"capacity_snapshot"`
    SpendSnapshot JSONSubsetValue `tfsdk:"spend_snapshot"`
    UsageReadAt RFC3339Value `tfsdk:"usage_read_at"`
    CollectionStatus types.String `tfsdk:"collection_status"`
    CollectionErrorCode types.String `tfsdk:"collection_error_code"`
    CollectionError types.String `tfsdk:"collection_error"`
    PresentedCertificate JSONSubsetValue `tfsdk:"presented_certificate"`
    CollectionSummary JSONSubsetValue `tfsdk:"collection_summary"`
    NextCollectionAt RFC3339Value `tfsdk:"next_collection_at"`
    LastCollectionAt RFC3339Value `tfsdk:"last_collection_at"`
    LastSuccessfulCollectionAt RFC3339Value `tfsdk:"last_successful_collection_at"`
    CollectionSettingsVersion types.Number `tfsdk:"collection_settings_version"`
    AiAccessLastVerifiedAt RFC3339Value `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt RFC3339Value `tfsdk:"ai_access_configured_at"`
    AiAgentResetAt RFC3339Value `tfsdk:"ai_agent_reset_at"`
}

func (r *OpenstackCloudResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    if r.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_open_stack_cloud"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_openstack_cloud"
}

func (r *OpenstackCloudResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
    if r.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_open_stack_cloud has been renamed to oneuptime_openstack_cloud. Rename the resource in your configuration and add a moved block (from = oneuptime_open_stack_cloud.<name>, to = oneuptime_openstack_cloud.<name>) to keep the existing openstack cloud; the old name keeps working until then."
    }
}

func (r *OpenstackCloudResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "OpenStack clouds that are being monitored in this project. A cloud is discovered when the OneUptime OpenStack Agent sends its first metrics, or connected with a read-only credential that one of the project's probes collects it with.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this OpenStack cloud. This is the join key — it must match the openstack.cloud.name OTel resource attribute stamped by the OneUptime OpenStack Agent.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this OpenStack cloud.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this OpenStack cloud archived? Archived OpenStack clouds are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDSetValidators("oneuptime_label"),
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this OpenStack cloud. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this OpenStack cloud (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the OpenStack cloud default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "collection_method": schema.StringAttribute{
                MarkdownDescription: "How this cloud's data reaches OneUptime. Probe: one of your OneUptime probes reads the cloud's APIs with the credential saved here, with no agent to install. Agent: the OneUptime OpenStack Agent you run sends it. Switching to Agent forgets the saved credential.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Agent"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Agent", "Probe"),
                },
            },
            "auth_url": schema.StringAttribute{
                MarkdownDescription: "The HTTPS address of the cloud's Keystone, which the probe logs in at and reads the service catalog from, such as https://keystone.example.com:5000/v3. Stored as https://host[:port][/path], without the version. Changing it asks for the credential again.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "endpoint_interface": schema.StringAttribute{
                MarkdownDescription: "Which of the catalog's addresses the probe uses for each service: public, internal or admin. Most probes sit where public answers. One in the control plane's own network may want internal.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("public"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("public", "internal", "admin"),
                },
            },
            "collection_regions": schema.StringAttribute{
                MarkdownDescription: "The regions the probe reads, as a JSON array of the names the catalog uses, such as [\"RegionOne\"]. The connection test lists the ones the catalog names. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "credential_kind": schema.StringAttribute{
                MarkdownDescription: "Which kind of credential the probe logs in with: applicationCredential or password. Worked out from what is written with a new secret, so it can be left out.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("applicationCredential", "password"),
                },
            },
            "application_credential_id": schema.StringAttribute{
                MarkdownDescription: "The ID of the Keystone application credential the probe logs in with, the better of the two kinds because it can be limited to reading and can expire. Its secret goes in Credential Secret. Changing it asks for the secret again.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "keystone_username": schema.StringAttribute{
                MarkdownDescription: "The Keystone user the probe logs in as, when it logs in with a password instead of an application credential. The password goes in Credential Secret. Changing it asks for the password again.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "user_domain_name": schema.StringAttribute{
                MarkdownDescription: "The Keystone domain of that user. Default, when left out.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_name": schema.StringAttribute{
                MarkdownDescription: "The project that user logs in to, when it logs in with a password. An application credential already names its own project.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_domain_name": schema.StringAttribute{
                MarkdownDescription: "The Keystone domain of that project. Default, when left out.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "credential_secret": schema.StringAttribute{
                MarkdownDescription: "The application credential's secret, or the password. Write-only: encrypted at rest, never returned by the API, and sent only to the probe that collects this cloud. Leave it out of an update to keep the saved one. Changing the address or the probe, or trusting another certificate, needs it again.",
                Optional: true,
                Sensitive: true,
            },
            "collection_probe_id": schema.StringAttribute{
                MarkdownDescription: "ID of the OneUptime probe that collects this cloud. Changing it asks for the credential again. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDValidators("oneuptime_probe"),
            },
            "trusted_certificate_fingerprints": schema.StringAttribute{
                MarkdownDescription: "The SHA-256 fingerprints of the certificates the probe accepts from this cloud's APIs, as a JSON array, for APIs whose certificates are not from a public authority. More than one, because each service can have a host and a certificate of its own. Empty: every certificate must be signed by an authority the probe's machine trusts. Verification is never skipped. Trusting another certificate asks for the credential again. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "collection_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often the probe collects this cloud, in whole minutes from 1 to 60. Two minutes is the OpenStack agent's own default. Collect large clouds less often to go easier on their APIs.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(2)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "collection_options": schema.StringAttribute{
                MarkdownDescription: "What the probe reads beyond the agent's defaults, as JSON. collectPorts reads ports, scoped by portNetworkIds and portProjectIds and paged. collectProjectLimits reads each project's quotas, four calls per project. Both are off unless asked for. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "instance_metrics_options": schema.StringAttribute{
                MarkdownDescription: "Whether real per-instance CPU, memory, disk and network series are kept for this cloud, as JSON: mode (automatic, all, selected or off), and for selected, the projectIds and aggregates. Per-instance series multiply fast, so automatic stops keeping them past a stated number of instances. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "price_list": schema.StringAttribute{
                MarkdownDescription: "What this cloud's operator says things cost, per hour, in one currency, as JSON. Optional: usage works without it, and a price of zero means no price has been set. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime OpenStack agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When metrics were last received from this cloud.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    RFC3339KeepSameInstant(),
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "region_count": schema.NumberAttribute{
                MarkdownDescription: "How many regions this cloud reports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "availability_zone_count": schema.NumberAttribute{
                MarkdownDescription: "How many availability zones this cloud reports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "hypervisor_count": schema.NumberAttribute{
                MarkdownDescription: "How many hypervisors this cloud reports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "instance_count": schema.NumberAttribute{
                MarkdownDescription: "How many instances this cloud reports, in every project.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "active_instance_count": schema.NumberAttribute{
                MarkdownDescription: "How many instances are ACTIVE.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "volume_count": schema.NumberAttribute{
                MarkdownDescription: "How many volumes this cloud reports, in every project.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "network_count": schema.NumberAttribute{
                MarkdownDescription: "How many networks this cloud reports, in every project.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "router_count": schema.NumberAttribute{
                MarkdownDescription: "How many routers this cloud reports, in every project.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "floating_ip_count": schema.NumberAttribute{
                MarkdownDescription: "How many floating IPs this cloud reports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "project_count": schema.NumberAttribute{
                MarkdownDescription: "How many OpenStack projects this cloud reports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "service_count": schema.NumberAttribute{
                MarkdownDescription: "How many services and agents (nova-compute, neutron agents, cinder-volume and so on) this cloud reports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "down_service_count": schema.NumberAttribute{
                MarkdownDescription: "How many enabled services and agents are down. A service someone disabled is not counted.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "vcpu_capacity": schema.NumberAttribute{
                MarkdownDescription: "vCPUs the cloud can promise to instances, as Placement counts them: (total - reserved) * allocation ratio.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "vcpu_allocated": schema.NumberAttribute{
                MarkdownDescription: "vCPUs promised to instances. This is what is promised, not what is busy.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "memory_capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Memory the cloud can promise to instances, as Placement counts it.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "memory_allocated_bytes": schema.NumberAttribute{
                MarkdownDescription: "Memory promised to instances. This is what is promised, not what is in use.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "storage_capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Total size of the cloud's Cinder storage pools.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "storage_used_bytes": schema.NumberAttribute{
                MarkdownDescription: "How much of the cloud's Cinder storage pools is used.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (openstack server, hypervisor, service and agent reads, each one call to the cloud's API) on this OpenStack cloud, through its OpenStack AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the OpenStack cloud can turn it on or off.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI only reads an OpenStack cloud today: it has no command that changes one, so this stays Disabled and any other mode is refused. Disabled: AI never proposes or runs a change on this OpenStack cloud. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the OpenStack cloud's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the OpenStack cloud can lower the mode; raising it needs Project Owner, Project Admin or Edit AI Rule.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Disabled"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Disabled", "RequireApproval", "Automatic", "BypassApproval"),
                },
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI only reads an OpenStack cloud today, and a read needs no approval, so no pattern is valid here yet and this stays empty. Optional JSON array of command patterns that Automatic mode may run on this OpenStack cloud without approval even though they are riskier changes. Each pattern is one command line for this OpenStack cloud's agent (openstack) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this OpenStack cloud is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit AI Rule; anyone who may edit the OpenStack cloud can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this OpenStack cloud archived?",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "is_credential_set": schema.BoolAttribute{
                MarkdownDescription: "Whether a credential is saved. Set by the server from the secret itself.",
                Computed: true,
            },
            "credentials_updated_at": schema.StringAttribute{
                MarkdownDescription: "When the credential was last changed. Set by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "capacity_snapshot": schema.StringAttribute{
                MarkdownDescription: "What this cloud's capacity added up to when a worker last read its inventory, as JSON: what is allocated and what can be, by zone, aggregate, region and cloud, how many more of each flavor fit, when each thing fills up, and what no instance explains. Written by the server. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "spend_snapshot": schema.StringAttribute{
                MarkdownDescription: "What this cloud's usage added up to when a worker last read its inventory, as JSON: this month against last, the largest domains and flavors, the top growers, and what could be reclaimed. Estimates, never an invoice. Written by the server. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "usage_read_at": schema.StringAttribute{
                MarkdownDescription: "The time this cloud's usage has been counted up to. Written by the server, and how a time is never counted twice.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "collection_status": schema.StringAttribute{
                MarkdownDescription: "Where probe collection stands: Pending (the probe has not tried the current settings yet), Succeeded, PartiallySucceeded (some of the cloud could not be read) or Failed. Written by the server and the probe.",
                Computed: true,
            },
            "collection_error_code": schema.StringAttribute{
                MarkdownDescription: "Why the last collection failed or was partial, as a code the dashboard turns into a next step (InvalidCredential, UntrustedCertificate, ConnectionTimedOut, ...). Written by the probe.",
                Computed: true,
            },
            "collection_error": schema.StringAttribute{
                MarkdownDescription: "What went wrong in the last collection, in words. Written by the probe.",
                Computed: true,
            },
            "presented_certificate": schema.StringAttribute{
                MarkdownDescription: "The certificate an API presented when the last collection did not trust it: its SHA-256 fingerprint, subject, issuer, validity and the host that showed it, so a person can check it and trust exactly that certificate. Written by the probe. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "collection_summary": schema.StringAttribute{
                MarkdownDescription: "What the last collection found: the regions, how each service answered, how many hypervisors, instances, networks, volumes and projects it read, and what it left out. Written by the probe. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "next_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe collects this cloud next. Written by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe last tried to collect this cloud, whether it worked or not. Written by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_successful_collection_at": schema.StringAttribute{
                MarkdownDescription: "When a collection of this cloud last worked, in full or in part. Written by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "collection_settings_version": schema.NumberAttribute{
                MarkdownDescription: "Counts changes to the probe collection settings, so a report from a collection that started before a change is not shown as the current status. Written by the server.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this OpenStack cloud through its OpenStack AI agent. Set by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this OpenStack cloud, kept until the next successful command. Set by the server.",
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this OpenStack cloud was first configured by anyone saving an AI access setting. Set by the server; never cleared, so an OpenStack AI agent that registers later never overwrites a setting an operator chose.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "ai_agent_reset_at": schema.StringAttribute{
                MarkdownDescription: "When this OpenStack cloud's AI agent was last reset: its key stopped working at once, and it registers again by itself. Set by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
        },
    }
}

// MoveState lets a moved block bring state over from the old name:
//
//     moved {
//       from = oneuptime_open_stack_cloud.example
//       to   = oneuptime_openstack_cloud.example
//     }
func (r *OpenstackCloudResource) MoveState(ctx context.Context) []resource.StateMover {
    if r.isLegacyAlias {
        return nil
    }

    return []resource.StateMover{
        legacyNameStateMover("oneuptime_open_stack_cloud", r.schemaDefinition()),
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *OpenstackCloudResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    keepStateWhenNothingChanges(ctx, req, resp)
}

func (r *OpenstackCloudResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *OpenstackCloudResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data OpenstackCloudResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config OpenstackCloudResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    openstackCloudRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := openstackCloudRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        requestDataMap["retainTelemetryDataForDays"] = r.bigFloatToFloat64(data.RetainTelemetryDataForDays.ValueBigFloat())
    }
    if parsedTelemetryRetentionConfig := r.parseJSONField(data.TelemetryRetentionConfig); parsedTelemetryRetentionConfig != nil {
        requestDataMap["telemetryRetentionConfig"] = parsedTelemetryRetentionConfig
    }
    if !data.CollectionMethod.IsNull() && !data.CollectionMethod.IsUnknown() {
        requestDataMap["collectionMethod"] = data.CollectionMethod.ValueString()
    }
    if !data.AuthUrl.IsNull() && !data.AuthUrl.IsUnknown() {
        requestDataMap["authUrl"] = data.AuthUrl.ValueString()
    }
    if !data.EndpointInterface.IsNull() && !data.EndpointInterface.IsUnknown() {
        requestDataMap["endpointInterface"] = data.EndpointInterface.ValueString()
    }
    if parsedCollectionRegions := r.parseJSONField(data.CollectionRegions); parsedCollectionRegions != nil {
        requestDataMap["collectionRegions"] = parsedCollectionRegions
    }
    if !data.CredentialKind.IsNull() && !data.CredentialKind.IsUnknown() {
        requestDataMap["credentialKind"] = data.CredentialKind.ValueString()
    }
    if !data.ApplicationCredentialId.IsNull() && !data.ApplicationCredentialId.IsUnknown() {
        requestDataMap["applicationCredentialId"] = data.ApplicationCredentialId.ValueString()
    }
    if !data.KeystoneUsername.IsNull() && !data.KeystoneUsername.IsUnknown() {
        requestDataMap["keystoneUsername"] = data.KeystoneUsername.ValueString()
    }
    if !data.UserDomainName.IsNull() && !data.UserDomainName.IsUnknown() {
        requestDataMap["userDomainName"] = data.UserDomainName.ValueString()
    }
    if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
        requestDataMap["projectName"] = data.ProjectName.ValueString()
    }
    if !data.ProjectDomainName.IsNull() && !data.ProjectDomainName.IsUnknown() {
        requestDataMap["projectDomainName"] = data.ProjectDomainName.ValueString()
    }
    if !data.CredentialSecret.IsNull() && !data.CredentialSecret.IsUnknown() {
        requestDataMap["credentialSecret"] = data.CredentialSecret.ValueString()
    }
    if !data.CollectionProbeId.IsNull() && !data.CollectionProbeId.IsUnknown() {
        requestDataMap["collectionProbeId"] = data.CollectionProbeId.ValueString()
    }
    if parsedTrustedCertificateFingerprints := r.parseJSONField(data.TrustedCertificateFingerprints); parsedTrustedCertificateFingerprints != nil {
        requestDataMap["trustedCertificateFingerprints"] = parsedTrustedCertificateFingerprints
    }
    if !data.CollectionIntervalInMinutes.IsNull() && !data.CollectionIntervalInMinutes.IsUnknown() {
        requestDataMap["collectionIntervalInMinutes"] = r.bigFloatToFloat64(data.CollectionIntervalInMinutes.ValueBigFloat())
    }
    if parsedCollectionOptions := r.parseJSONField(data.CollectionOptions); parsedCollectionOptions != nil {
        requestDataMap["collectionOptions"] = parsedCollectionOptions
    }
    if parsedInstanceMetricsOptions := r.parseJSONField(data.InstanceMetricsOptions); parsedInstanceMetricsOptions != nil {
        requestDataMap["instanceMetricsOptions"] = parsedInstanceMetricsOptions
    }
    if parsedPriceList := r.parseJSONField(data.PriceList); parsedPriceList != nil {
        requestDataMap["priceList"] = parsedPriceList
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/openstack-cloud", openstackCloudRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create openstack_cloud, got error: %s", err))
        return
    }

    var openstackCloudResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &openstackCloudResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create openstack_cloud: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := openstackCloudResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := openstackCloudResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for openstack_cloud did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * openstack_cloud is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "isArchived": true,
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
        "collectionProbeId": true,
        "trustedCertificateFingerprints": true,
        "collectionIntervalInMinutes": true,
        "collectionOptions": true,
        "instanceMetricsOptions": true,
        "priceList": true,
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
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "isCredentialSet": true,
        "credentialsUpdatedAt": true,
        "capacitySnapshot": true,
        "spendSnapshot": true,
        "usageReadAt": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "aiAgentResetAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/openstack-cloud/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created openstack_cloud but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created openstack_cloud but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    data.RetainTelemetryDataForDays = numberFromAPI(data.RetainTelemetryDataForDays, dataMap["retainTelemetryDataForDays"])
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["collectionMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionMethod = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionMethod = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionMethod"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if val, ok := dataMap["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := dataMap["authUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AuthUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AuthUrl = types.StringValue(string(jsonBytes))
            } else {
                data.AuthUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AuthUrl = types.StringValue(string(jsonBytes))
            } else {
                data.AuthUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if arr, ok := dataMap["authUrl"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if val, ok := dataMap["authUrl"].(string); ok {
        data.AuthUrl = types.StringValue(val)
    } else {
        data.AuthUrl = types.StringNull()
    }
    if obj, ok := dataMap["endpointInterface"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndpointInterface = types.StringValue(string(jsonBytes))
            } else {
                data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndpointInterface = types.StringValue(string(jsonBytes))
            } else {
                data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if arr, ok := dataMap["endpointInterface"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if val, ok := dataMap["endpointInterface"].(string); ok {
        data.EndpointInterface = types.StringValue(val)
    } else {
        data.EndpointInterface = types.StringNull()
    }
    if obj, ok := dataMap["collectionRegions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionRegions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionRegions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionRegions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionRegions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionRegions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionRegions"].(string); ok {
        data.CollectionRegions = NewJSONSubsetValue(val)
    } else {
        data.CollectionRegions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["credentialKind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CredentialKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CredentialKind = types.StringValue(string(jsonBytes))
            } else {
                data.CredentialKind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CredentialKind = types.StringValue(string(jsonBytes))
            } else {
                data.CredentialKind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if arr, ok := dataMap["credentialKind"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if val, ok := dataMap["credentialKind"].(string); ok {
        data.CredentialKind = types.StringValue(val)
    } else {
        data.CredentialKind = types.StringNull()
    }
    if obj, ok := dataMap["applicationCredentialId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
            } else {
                data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
            } else {
                data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if arr, ok := dataMap["applicationCredentialId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if val, ok := dataMap["applicationCredentialId"].(string); ok {
        data.ApplicationCredentialId = types.StringValue(val)
    } else {
        data.ApplicationCredentialId = types.StringNull()
    }
    if obj, ok := dataMap["keystoneUsername"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.KeystoneUsername = types.StringValue(string(jsonBytes))
            } else {
                data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.KeystoneUsername = types.StringValue(string(jsonBytes))
            } else {
                data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if arr, ok := dataMap["keystoneUsername"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if val, ok := dataMap["keystoneUsername"].(string); ok {
        data.KeystoneUsername = types.StringValue(val)
    } else {
        data.KeystoneUsername = types.StringNull()
    }
    if obj, ok := dataMap["userDomainName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UserDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UserDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.UserDomainName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UserDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.UserDomainName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if arr, ok := dataMap["userDomainName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if val, ok := dataMap["userDomainName"].(string); ok {
        data.UserDomainName = types.StringValue(val)
    } else {
        data.UserDomainName = types.StringNull()
    }
    if obj, ok := dataMap["projectName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProjectName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProjectName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProjectName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if arr, ok := dataMap["projectName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if val, ok := dataMap["projectName"].(string); ok {
        data.ProjectName = types.StringValue(val)
    } else {
        data.ProjectName = types.StringNull()
    }
    if obj, ok := dataMap["projectDomainName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProjectDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProjectDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if arr, ok := dataMap["projectDomainName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if val, ok := dataMap["projectDomainName"].(string); ok {
        data.ProjectDomainName = types.StringValue(val)
    } else {
        data.ProjectDomainName = types.StringNull()
    }
    if obj, ok := dataMap["collectionProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionProbeId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := dataMap["trustedCertificateFingerprints"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["trustedCertificateFingerprints"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["trustedCertificateFingerprints"].(string); ok {
        data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
    } else {
        data.TrustedCertificateFingerprints = NewJSONSubsetNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, dataMap["collectionIntervalInMinutes"])
    if obj, ok := dataMap["collectionOptions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionOptions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionOptions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionOptions"].(string); ok {
        data.CollectionOptions = NewJSONSubsetValue(val)
    } else {
        data.CollectionOptions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["instanceMetricsOptions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InstanceMetricsOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InstanceMetricsOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["instanceMetricsOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["instanceMetricsOptions"].(string); ok {
        data.InstanceMetricsOptions = NewJSONSubsetValue(val)
    } else {
        data.InstanceMetricsOptions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["priceList"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PriceList = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PriceList = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PriceList = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PriceList = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PriceList = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PriceList = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["priceList"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PriceList = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PriceList = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["priceList"].(string); ok {
        data.PriceList = NewJSONSubsetValue(val)
    } else {
        data.PriceList = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["agentVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    data.RegionCount = numberFromAPI(data.RegionCount, dataMap["regionCount"])
    data.AvailabilityZoneCount = numberFromAPI(data.AvailabilityZoneCount, dataMap["availabilityZoneCount"])
    data.HypervisorCount = numberFromAPI(data.HypervisorCount, dataMap["hypervisorCount"])
    data.InstanceCount = numberFromAPI(data.InstanceCount, dataMap["instanceCount"])
    data.ActiveInstanceCount = numberFromAPI(data.ActiveInstanceCount, dataMap["activeInstanceCount"])
    data.VolumeCount = numberFromAPI(data.VolumeCount, dataMap["volumeCount"])
    data.NetworkCount = numberFromAPI(data.NetworkCount, dataMap["networkCount"])
    data.RouterCount = numberFromAPI(data.RouterCount, dataMap["routerCount"])
    data.FloatingIpCount = numberFromAPI(data.FloatingIpCount, dataMap["floatingIpCount"])
    data.ProjectCount = numberFromAPI(data.ProjectCount, dataMap["projectCount"])
    data.ServiceCount = numberFromAPI(data.ServiceCount, dataMap["serviceCount"])
    data.DownServiceCount = numberFromAPI(data.DownServiceCount, dataMap["downServiceCount"])
    data.VcpuCapacity = numberFromAPI(data.VcpuCapacity, dataMap["vcpuCapacity"])
    data.VcpuAllocated = numberFromAPI(data.VcpuAllocated, dataMap["vcpuAllocated"])
    data.MemoryCapacityBytes = numberFromAPI(data.MemoryCapacityBytes, dataMap["memoryCapacityBytes"])
    data.MemoryAllocatedBytes = numberFromAPI(data.MemoryAllocatedBytes, dataMap["memoryAllocatedBytes"])
    data.StorageCapacityBytes = numberFromAPI(data.StorageCapacityBytes, dataMap["storageCapacityBytes"])
    data.StorageUsedBytes = numberFromAPI(data.StorageUsedBytes, dataMap["storageUsedBytes"])
    if val, ok := dataMap["isAiInvestigationEnabled"].(bool); ok {
        data.IsAiInvestigationEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["aiRemediationMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiRemediationMode = types.StringValue(string(jsonBytes))
            } else {
                data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiRemediationMode = types.StringValue(string(jsonBytes))
            } else {
                data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if arr, ok := dataMap["aiRemediationMode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if val, ok := dataMap["aiRemediationMode"].(string); ok {
        data.AiRemediationMode = types.StringValue(val)
    } else {
        data.AiRemediationMode = types.StringNull()
    }
    if obj, ok := dataMap["aiCommandAllowlist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiCommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiCommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["aiCommandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["aiCommandAllowlist"].(string); ok {
        data.AiCommandAllowlist = NewJSONSubsetValue(val)
    } else {
        data.AiCommandAllowlist = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["archivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isCredentialSet"].(bool); ok {
        data.IsCredentialSet = types.BoolValue(val)
    }
    if obj, ok := dataMap["credentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CredentialsUpdatedAt = NewRFC3339Value(val)
        } else {
            data.CredentialsUpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["credentialsUpdatedAt"].(string); ok && val != "" {
        data.CredentialsUpdatedAt = NewRFC3339Value(val)
    } else {
        data.CredentialsUpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["capacitySnapshot"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CapacitySnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CapacitySnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["capacitySnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["capacitySnapshot"].(string); ok {
        data.CapacitySnapshot = NewJSONSubsetValue(val)
    } else {
        data.CapacitySnapshot = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["spendSnapshot"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpendSnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SpendSnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["spendSnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["spendSnapshot"].(string); ok {
        data.SpendSnapshot = NewJSONSubsetValue(val)
    } else {
        data.SpendSnapshot = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["usageReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UsageReadAt = NewRFC3339Value(val)
        } else {
            data.UsageReadAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["usageReadAt"].(string); ok && val != "" {
        data.UsageReadAt = NewRFC3339Value(val)
    } else {
        data.UsageReadAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["collectionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["collectionStatus"].(string); ok {
        data.CollectionStatus = types.StringValue(val)
    } else {
        data.CollectionStatus = types.StringNull()
    }
    if obj, ok := dataMap["collectionErrorCode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionErrorCode = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionErrorCode = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionErrorCode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if val, ok := dataMap["collectionErrorCode"].(string); ok {
        data.CollectionErrorCode = types.StringValue(val)
    } else {
        data.CollectionErrorCode = types.StringNull()
    }
    if obj, ok := dataMap["collectionError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionError = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionError = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionError"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if val, ok := dataMap["collectionError"].(string); ok {
        data.CollectionError = types.StringValue(val)
    } else {
        data.CollectionError = types.StringNull()
    }
    if obj, ok := dataMap["presentedCertificate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PresentedCertificate = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PresentedCertificate = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["presentedCertificate"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["presentedCertificate"].(string); ok {
        data.PresentedCertificate = NewJSONSubsetValue(val)
    } else {
        data.PresentedCertificate = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["collectionSummary"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionSummary = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionSummary = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionSummary = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionSummary"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionSummary = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionSummary"].(string); ok {
        data.CollectionSummary = NewJSONSubsetValue(val)
    } else {
        data.CollectionSummary = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["nextCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextCollectionAt = NewRFC3339Value(val)
        } else {
            data.NextCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextCollectionAt"].(string); ok && val != "" {
        data.NextCollectionAt = NewRFC3339Value(val)
    } else {
        data.NextCollectionAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastCollectionAt = NewRFC3339Value(val)
        } else {
            data.LastCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastCollectionAt"].(string); ok && val != "" {
        data.LastCollectionAt = NewRFC3339Value(val)
    } else {
        data.LastCollectionAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastSuccessfulCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSuccessfulCollectionAt = NewRFC3339Value(val)
        } else {
            data.LastSuccessfulCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSuccessfulCollectionAt"].(string); ok && val != "" {
        data.LastSuccessfulCollectionAt = NewRFC3339Value(val)
    } else {
        data.LastSuccessfulCollectionAt = NewRFC3339Null()
    }
    data.CollectionSettingsVersion = numberFromAPI(data.CollectionSettingsVersion, dataMap["collectionSettingsVersion"])
    if obj, ok := dataMap["aiAccessLastVerifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAccessLastVerifiedAt = NewRFC3339Value(val)
        } else {
            data.AiAccessLastVerifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAccessLastVerifiedAt"].(string); ok && val != "" {
        data.AiAccessLastVerifiedAt = NewRFC3339Value(val)
    } else {
        data.AiAccessLastVerifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiAccessLastError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiAccessLastError = types.StringValue(string(jsonBytes))
            } else {
                data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiAccessLastError = types.StringValue(string(jsonBytes))
            } else {
                data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if arr, ok := dataMap["aiAccessLastError"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if val, ok := dataMap["aiAccessLastError"].(string); ok {
        data.AiAccessLastError = types.StringValue(val)
    } else {
        data.AiAccessLastError = types.StringNull()
    }
    if obj, ok := dataMap["aiAccessConfiguredAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAccessConfiguredAt = NewRFC3339Value(val)
        } else {
            data.AiAccessConfiguredAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAccessConfiguredAt"].(string); ok && val != "" {
        data.AiAccessConfiguredAt = NewRFC3339Value(val)
    } else {
        data.AiAccessConfiguredAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiAgentResetAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentResetAt = NewRFC3339Value(val)
        } else {
            data.AiAgentResetAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAgentResetAt"].(string); ok && val != "" {
        data.AiAgentResetAt = NewRFC3339Value(val)
    } else {
        data.AiAgentResetAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")


    // Nothing a create made was imported (see importsettle.go).
    resp.Diagnostics.Append(forgetImportedFromServer(ctx, resp.Private)...)
    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpenstackCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data OpenstackCloudResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "isArchived": true,
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
        "collectionProbeId": true,
        "trustedCertificateFingerprints": true,
        "collectionIntervalInMinutes": true,
        "collectionOptions": true,
        "instanceMetricsOptions": true,
        "priceList": true,
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
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "isCredentialSet": true,
        "credentialsUpdatedAt": true,
        "capacitySnapshot": true,
        "spendSnapshot": true,
        "usageReadAt": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "aiAgentResetAt": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/openstack-cloud/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read openstack_cloud, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var openstackCloudResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &openstackCloudResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse openstack_cloud response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := openstackCloudResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = openstackCloudResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    data.RetainTelemetryDataForDays = numberFromAPI(data.RetainTelemetryDataForDays, dataMap["retainTelemetryDataForDays"])
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["collectionMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionMethod = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionMethod = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionMethod"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if val, ok := dataMap["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := dataMap["authUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AuthUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AuthUrl = types.StringValue(string(jsonBytes))
            } else {
                data.AuthUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AuthUrl = types.StringValue(string(jsonBytes))
            } else {
                data.AuthUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if arr, ok := dataMap["authUrl"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if val, ok := dataMap["authUrl"].(string); ok {
        data.AuthUrl = types.StringValue(val)
    } else {
        data.AuthUrl = types.StringNull()
    }
    if obj, ok := dataMap["endpointInterface"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndpointInterface = types.StringValue(string(jsonBytes))
            } else {
                data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndpointInterface = types.StringValue(string(jsonBytes))
            } else {
                data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if arr, ok := dataMap["endpointInterface"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if val, ok := dataMap["endpointInterface"].(string); ok {
        data.EndpointInterface = types.StringValue(val)
    } else {
        data.EndpointInterface = types.StringNull()
    }
    if obj, ok := dataMap["collectionRegions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionRegions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionRegions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionRegions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionRegions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionRegions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionRegions"].(string); ok {
        data.CollectionRegions = NewJSONSubsetValue(val)
    } else {
        data.CollectionRegions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["credentialKind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CredentialKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CredentialKind = types.StringValue(string(jsonBytes))
            } else {
                data.CredentialKind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CredentialKind = types.StringValue(string(jsonBytes))
            } else {
                data.CredentialKind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if arr, ok := dataMap["credentialKind"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if val, ok := dataMap["credentialKind"].(string); ok {
        data.CredentialKind = types.StringValue(val)
    } else {
        data.CredentialKind = types.StringNull()
    }
    if obj, ok := dataMap["applicationCredentialId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
            } else {
                data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
            } else {
                data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if arr, ok := dataMap["applicationCredentialId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if val, ok := dataMap["applicationCredentialId"].(string); ok {
        data.ApplicationCredentialId = types.StringValue(val)
    } else {
        data.ApplicationCredentialId = types.StringNull()
    }
    if obj, ok := dataMap["keystoneUsername"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.KeystoneUsername = types.StringValue(string(jsonBytes))
            } else {
                data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.KeystoneUsername = types.StringValue(string(jsonBytes))
            } else {
                data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if arr, ok := dataMap["keystoneUsername"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if val, ok := dataMap["keystoneUsername"].(string); ok {
        data.KeystoneUsername = types.StringValue(val)
    } else {
        data.KeystoneUsername = types.StringNull()
    }
    if obj, ok := dataMap["userDomainName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UserDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UserDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.UserDomainName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UserDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.UserDomainName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if arr, ok := dataMap["userDomainName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if val, ok := dataMap["userDomainName"].(string); ok {
        data.UserDomainName = types.StringValue(val)
    } else {
        data.UserDomainName = types.StringNull()
    }
    if obj, ok := dataMap["projectName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProjectName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProjectName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProjectName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if arr, ok := dataMap["projectName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if val, ok := dataMap["projectName"].(string); ok {
        data.ProjectName = types.StringValue(val)
    } else {
        data.ProjectName = types.StringNull()
    }
    if obj, ok := dataMap["projectDomainName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProjectDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProjectDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if arr, ok := dataMap["projectDomainName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if val, ok := dataMap["projectDomainName"].(string); ok {
        data.ProjectDomainName = types.StringValue(val)
    } else {
        data.ProjectDomainName = types.StringNull()
    }
    if obj, ok := dataMap["collectionProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionProbeId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := dataMap["trustedCertificateFingerprints"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["trustedCertificateFingerprints"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["trustedCertificateFingerprints"].(string); ok {
        data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
    } else {
        data.TrustedCertificateFingerprints = NewJSONSubsetNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, dataMap["collectionIntervalInMinutes"])
    if obj, ok := dataMap["collectionOptions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionOptions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionOptions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionOptions"].(string); ok {
        data.CollectionOptions = NewJSONSubsetValue(val)
    } else {
        data.CollectionOptions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["instanceMetricsOptions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InstanceMetricsOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InstanceMetricsOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["instanceMetricsOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["instanceMetricsOptions"].(string); ok {
        data.InstanceMetricsOptions = NewJSONSubsetValue(val)
    } else {
        data.InstanceMetricsOptions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["priceList"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PriceList = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PriceList = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PriceList = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PriceList = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PriceList = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PriceList = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["priceList"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PriceList = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PriceList = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["priceList"].(string); ok {
        data.PriceList = NewJSONSubsetValue(val)
    } else {
        data.PriceList = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["agentVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    data.RegionCount = numberFromAPI(data.RegionCount, dataMap["regionCount"])
    data.AvailabilityZoneCount = numberFromAPI(data.AvailabilityZoneCount, dataMap["availabilityZoneCount"])
    data.HypervisorCount = numberFromAPI(data.HypervisorCount, dataMap["hypervisorCount"])
    data.InstanceCount = numberFromAPI(data.InstanceCount, dataMap["instanceCount"])
    data.ActiveInstanceCount = numberFromAPI(data.ActiveInstanceCount, dataMap["activeInstanceCount"])
    data.VolumeCount = numberFromAPI(data.VolumeCount, dataMap["volumeCount"])
    data.NetworkCount = numberFromAPI(data.NetworkCount, dataMap["networkCount"])
    data.RouterCount = numberFromAPI(data.RouterCount, dataMap["routerCount"])
    data.FloatingIpCount = numberFromAPI(data.FloatingIpCount, dataMap["floatingIpCount"])
    data.ProjectCount = numberFromAPI(data.ProjectCount, dataMap["projectCount"])
    data.ServiceCount = numberFromAPI(data.ServiceCount, dataMap["serviceCount"])
    data.DownServiceCount = numberFromAPI(data.DownServiceCount, dataMap["downServiceCount"])
    data.VcpuCapacity = numberFromAPI(data.VcpuCapacity, dataMap["vcpuCapacity"])
    data.VcpuAllocated = numberFromAPI(data.VcpuAllocated, dataMap["vcpuAllocated"])
    data.MemoryCapacityBytes = numberFromAPI(data.MemoryCapacityBytes, dataMap["memoryCapacityBytes"])
    data.MemoryAllocatedBytes = numberFromAPI(data.MemoryAllocatedBytes, dataMap["memoryAllocatedBytes"])
    data.StorageCapacityBytes = numberFromAPI(data.StorageCapacityBytes, dataMap["storageCapacityBytes"])
    data.StorageUsedBytes = numberFromAPI(data.StorageUsedBytes, dataMap["storageUsedBytes"])
    if val, ok := dataMap["isAiInvestigationEnabled"].(bool); ok {
        data.IsAiInvestigationEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["aiRemediationMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiRemediationMode = types.StringValue(string(jsonBytes))
            } else {
                data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiRemediationMode = types.StringValue(string(jsonBytes))
            } else {
                data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if arr, ok := dataMap["aiRemediationMode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if val, ok := dataMap["aiRemediationMode"].(string); ok {
        data.AiRemediationMode = types.StringValue(val)
    } else {
        data.AiRemediationMode = types.StringNull()
    }
    if obj, ok := dataMap["aiCommandAllowlist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiCommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiCommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["aiCommandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["aiCommandAllowlist"].(string); ok {
        data.AiCommandAllowlist = NewJSONSubsetValue(val)
    } else {
        data.AiCommandAllowlist = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["archivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isCredentialSet"].(bool); ok {
        data.IsCredentialSet = types.BoolValue(val)
    }
    if obj, ok := dataMap["credentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CredentialsUpdatedAt = NewRFC3339Value(val)
        } else {
            data.CredentialsUpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["credentialsUpdatedAt"].(string); ok && val != "" {
        data.CredentialsUpdatedAt = NewRFC3339Value(val)
    } else {
        data.CredentialsUpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["capacitySnapshot"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CapacitySnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CapacitySnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["capacitySnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["capacitySnapshot"].(string); ok {
        data.CapacitySnapshot = NewJSONSubsetValue(val)
    } else {
        data.CapacitySnapshot = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["spendSnapshot"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpendSnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SpendSnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["spendSnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["spendSnapshot"].(string); ok {
        data.SpendSnapshot = NewJSONSubsetValue(val)
    } else {
        data.SpendSnapshot = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["usageReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UsageReadAt = NewRFC3339Value(val)
        } else {
            data.UsageReadAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["usageReadAt"].(string); ok && val != "" {
        data.UsageReadAt = NewRFC3339Value(val)
    } else {
        data.UsageReadAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["collectionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["collectionStatus"].(string); ok {
        data.CollectionStatus = types.StringValue(val)
    } else {
        data.CollectionStatus = types.StringNull()
    }
    if obj, ok := dataMap["collectionErrorCode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionErrorCode = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionErrorCode = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionErrorCode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if val, ok := dataMap["collectionErrorCode"].(string); ok {
        data.CollectionErrorCode = types.StringValue(val)
    } else {
        data.CollectionErrorCode = types.StringNull()
    }
    if obj, ok := dataMap["collectionError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionError = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionError = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionError"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if val, ok := dataMap["collectionError"].(string); ok {
        data.CollectionError = types.StringValue(val)
    } else {
        data.CollectionError = types.StringNull()
    }
    if obj, ok := dataMap["presentedCertificate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PresentedCertificate = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PresentedCertificate = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["presentedCertificate"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["presentedCertificate"].(string); ok {
        data.PresentedCertificate = NewJSONSubsetValue(val)
    } else {
        data.PresentedCertificate = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["collectionSummary"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionSummary = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionSummary = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionSummary = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionSummary"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionSummary = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionSummary"].(string); ok {
        data.CollectionSummary = NewJSONSubsetValue(val)
    } else {
        data.CollectionSummary = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["nextCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextCollectionAt = NewRFC3339Value(val)
        } else {
            data.NextCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextCollectionAt"].(string); ok && val != "" {
        data.NextCollectionAt = NewRFC3339Value(val)
    } else {
        data.NextCollectionAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastCollectionAt = NewRFC3339Value(val)
        } else {
            data.LastCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastCollectionAt"].(string); ok && val != "" {
        data.LastCollectionAt = NewRFC3339Value(val)
    } else {
        data.LastCollectionAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastSuccessfulCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSuccessfulCollectionAt = NewRFC3339Value(val)
        } else {
            data.LastSuccessfulCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSuccessfulCollectionAt"].(string); ok && val != "" {
        data.LastSuccessfulCollectionAt = NewRFC3339Value(val)
    } else {
        data.LastSuccessfulCollectionAt = NewRFC3339Null()
    }
    data.CollectionSettingsVersion = numberFromAPI(data.CollectionSettingsVersion, dataMap["collectionSettingsVersion"])
    if obj, ok := dataMap["aiAccessLastVerifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAccessLastVerifiedAt = NewRFC3339Value(val)
        } else {
            data.AiAccessLastVerifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAccessLastVerifiedAt"].(string); ok && val != "" {
        data.AiAccessLastVerifiedAt = NewRFC3339Value(val)
    } else {
        data.AiAccessLastVerifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiAccessLastError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiAccessLastError = types.StringValue(string(jsonBytes))
            } else {
                data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiAccessLastError = types.StringValue(string(jsonBytes))
            } else {
                data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if arr, ok := dataMap["aiAccessLastError"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if val, ok := dataMap["aiAccessLastError"].(string); ok {
        data.AiAccessLastError = types.StringValue(val)
    } else {
        data.AiAccessLastError = types.StringNull()
    }
    if obj, ok := dataMap["aiAccessConfiguredAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAccessConfiguredAt = NewRFC3339Value(val)
        } else {
            data.AiAccessConfiguredAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAccessConfiguredAt"].(string); ok && val != "" {
        data.AiAccessConfiguredAt = NewRFC3339Value(val)
    } else {
        data.AiAccessConfiguredAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiAgentResetAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentResetAt = NewRFC3339Value(val)
        } else {
            data.AiAgentResetAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAgentResetAt"].(string); ok && val != "" {
        data.AiAgentResetAt = NewRFC3339Value(val)
    } else {
        data.AiAgentResetAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpenstackCloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data OpenstackCloudResourceModel
    var state OpenstackCloudResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // The attributes this update writes from the configuration (see importsettle.go).
    writtenFromConfig := map[string]bool{
        "ai_command_allowlist": !data.AiCommandAllowlist.Equal(state.AiCommandAllowlist),
        "collection_options": !data.CollectionOptions.Equal(state.CollectionOptions),
        "collection_regions": !data.CollectionRegions.Equal(state.CollectionRegions),
        "instance_metrics_options": !data.InstanceMetricsOptions.Equal(state.InstanceMetricsOptions),
        "price_list": !data.PriceList.Equal(state.PriceList),
        "telemetry_retention_config": !data.TelemetryRetentionConfig.Equal(state.TelemetryRetentionConfig),
        "trusted_certificate_fingerprints": !data.TrustedCertificateFingerprints.Equal(state.TrustedCertificateFingerprints),
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config OpenstackCloudResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    openstackCloudRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := openstackCloudRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.OtelCollectorStatus.IsUnknown() && !state.OtelCollectorStatus.IsUnknown() && !data.OtelCollectorStatus.Equal(state.OtelCollectorStatus) {
        requestDataMap["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
    }
    if !data.AgentVersion.IsUnknown() && !state.AgentVersion.IsUnknown() && !data.AgentVersion.Equal(state.AgentVersion) {
        requestDataMap["agentVersion"] = data.AgentVersion.ValueString()
    }
    if !data.LastSeenAt.IsUnknown() && !state.LastSeenAt.IsUnknown() && !data.LastSeenAt.Equal(state.LastSeenAt) {
        requestDataMap["lastSeenAt"] = data.LastSeenAt.ValueString()
    }
    if !data.RegionCount.IsUnknown() && !state.RegionCount.IsUnknown() && !data.RegionCount.Equal(state.RegionCount) {
        requestDataMap["regionCount"] = r.bigFloatToFloat64(data.RegionCount.ValueBigFloat())
    }
    if !data.AvailabilityZoneCount.IsUnknown() && !state.AvailabilityZoneCount.IsUnknown() && !data.AvailabilityZoneCount.Equal(state.AvailabilityZoneCount) {
        requestDataMap["availabilityZoneCount"] = r.bigFloatToFloat64(data.AvailabilityZoneCount.ValueBigFloat())
    }
    if !data.HypervisorCount.IsUnknown() && !state.HypervisorCount.IsUnknown() && !data.HypervisorCount.Equal(state.HypervisorCount) {
        requestDataMap["hypervisorCount"] = r.bigFloatToFloat64(data.HypervisorCount.ValueBigFloat())
    }
    if !data.InstanceCount.IsUnknown() && !state.InstanceCount.IsUnknown() && !data.InstanceCount.Equal(state.InstanceCount) {
        requestDataMap["instanceCount"] = r.bigFloatToFloat64(data.InstanceCount.ValueBigFloat())
    }
    if !data.ActiveInstanceCount.IsUnknown() && !state.ActiveInstanceCount.IsUnknown() && !data.ActiveInstanceCount.Equal(state.ActiveInstanceCount) {
        requestDataMap["activeInstanceCount"] = r.bigFloatToFloat64(data.ActiveInstanceCount.ValueBigFloat())
    }
    if !data.VolumeCount.IsUnknown() && !state.VolumeCount.IsUnknown() && !data.VolumeCount.Equal(state.VolumeCount) {
        requestDataMap["volumeCount"] = r.bigFloatToFloat64(data.VolumeCount.ValueBigFloat())
    }
    if !data.NetworkCount.IsUnknown() && !state.NetworkCount.IsUnknown() && !data.NetworkCount.Equal(state.NetworkCount) {
        requestDataMap["networkCount"] = r.bigFloatToFloat64(data.NetworkCount.ValueBigFloat())
    }
    if !data.RouterCount.IsUnknown() && !state.RouterCount.IsUnknown() && !data.RouterCount.Equal(state.RouterCount) {
        requestDataMap["routerCount"] = r.bigFloatToFloat64(data.RouterCount.ValueBigFloat())
    }
    if !data.FloatingIpCount.IsUnknown() && !state.FloatingIpCount.IsUnknown() && !data.FloatingIpCount.Equal(state.FloatingIpCount) {
        requestDataMap["floatingIpCount"] = r.bigFloatToFloat64(data.FloatingIpCount.ValueBigFloat())
    }
    if !data.ProjectCount.IsUnknown() && !state.ProjectCount.IsUnknown() && !data.ProjectCount.Equal(state.ProjectCount) {
        requestDataMap["projectCount"] = r.bigFloatToFloat64(data.ProjectCount.ValueBigFloat())
    }
    if !data.ServiceCount.IsUnknown() && !state.ServiceCount.IsUnknown() && !data.ServiceCount.Equal(state.ServiceCount) {
        requestDataMap["serviceCount"] = r.bigFloatToFloat64(data.ServiceCount.ValueBigFloat())
    }
    if !data.DownServiceCount.IsUnknown() && !state.DownServiceCount.IsUnknown() && !data.DownServiceCount.Equal(state.DownServiceCount) {
        requestDataMap["downServiceCount"] = r.bigFloatToFloat64(data.DownServiceCount.ValueBigFloat())
    }
    if !data.VcpuCapacity.IsUnknown() && !state.VcpuCapacity.IsUnknown() && !data.VcpuCapacity.Equal(state.VcpuCapacity) {
        requestDataMap["vcpuCapacity"] = r.bigFloatToFloat64(data.VcpuCapacity.ValueBigFloat())
    }
    if !data.VcpuAllocated.IsUnknown() && !state.VcpuAllocated.IsUnknown() && !data.VcpuAllocated.Equal(state.VcpuAllocated) {
        requestDataMap["vcpuAllocated"] = r.bigFloatToFloat64(data.VcpuAllocated.ValueBigFloat())
    }
    if !data.MemoryCapacityBytes.IsUnknown() && !state.MemoryCapacityBytes.IsUnknown() && !data.MemoryCapacityBytes.Equal(state.MemoryCapacityBytes) {
        requestDataMap["memoryCapacityBytes"] = r.bigFloatToFloat64(data.MemoryCapacityBytes.ValueBigFloat())
    }
    if !data.MemoryAllocatedBytes.IsUnknown() && !state.MemoryAllocatedBytes.IsUnknown() && !data.MemoryAllocatedBytes.Equal(state.MemoryAllocatedBytes) {
        requestDataMap["memoryAllocatedBytes"] = r.bigFloatToFloat64(data.MemoryAllocatedBytes.ValueBigFloat())
    }
    if !data.StorageCapacityBytes.IsUnknown() && !state.StorageCapacityBytes.IsUnknown() && !data.StorageCapacityBytes.Equal(state.StorageCapacityBytes) {
        requestDataMap["storageCapacityBytes"] = r.bigFloatToFloat64(data.StorageCapacityBytes.ValueBigFloat())
    }
    if !data.StorageUsedBytes.IsUnknown() && !state.StorageUsedBytes.IsUnknown() && !data.StorageUsedBytes.Equal(state.StorageUsedBytes) {
        requestDataMap["storageUsedBytes"] = r.bigFloatToFloat64(data.StorageUsedBytes.ValueBigFloat())
    }
    if !data.IsArchived.IsUnknown() && !state.IsArchived.IsUnknown() && !data.IsArchived.Equal(state.IsArchived) {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.RetainTelemetryDataForDays.IsUnknown() && !state.RetainTelemetryDataForDays.IsUnknown() && !data.RetainTelemetryDataForDays.Equal(state.RetainTelemetryDataForDays) {
        requestDataMap["retainTelemetryDataForDays"] = r.bigFloatToFloat64(data.RetainTelemetryDataForDays.ValueBigFloat())
    }
    if !data.TelemetryRetentionConfig.IsUnknown() && !state.TelemetryRetentionConfig.IsUnknown() && !data.TelemetryRetentionConfig.Equal(state.TelemetryRetentionConfig) {
        var telemetryretentionconfigData interface{}
        if err := json.Unmarshal([]byte(data.TelemetryRetentionConfig.ValueString()), &telemetryretentionconfigData); err == nil {
            requestDataMap["telemetryRetentionConfig"] = telemetryretentionconfigData
        } else {
            requestDataMap["telemetryRetentionConfig"] = data.TelemetryRetentionConfig.ValueString()
        }
    }
    if !data.CollectionMethod.IsUnknown() && !state.CollectionMethod.IsUnknown() && !data.CollectionMethod.Equal(state.CollectionMethod) {
        requestDataMap["collectionMethod"] = data.CollectionMethod.ValueString()
    }
    if !data.AuthUrl.IsUnknown() && !state.AuthUrl.IsUnknown() && !data.AuthUrl.Equal(state.AuthUrl) {
        requestDataMap["authUrl"] = data.AuthUrl.ValueString()
    }
    if !data.EndpointInterface.IsUnknown() && !state.EndpointInterface.IsUnknown() && !data.EndpointInterface.Equal(state.EndpointInterface) {
        requestDataMap["endpointInterface"] = data.EndpointInterface.ValueString()
    }
    if !data.CollectionRegions.IsUnknown() && !state.CollectionRegions.IsUnknown() && !data.CollectionRegions.Equal(state.CollectionRegions) {
        var collectionregionsData interface{}
        if err := json.Unmarshal([]byte(data.CollectionRegions.ValueString()), &collectionregionsData); err == nil {
            requestDataMap["collectionRegions"] = collectionregionsData
        } else {
            requestDataMap["collectionRegions"] = data.CollectionRegions.ValueString()
        }
    }
    if !data.CredentialKind.IsUnknown() && !state.CredentialKind.IsUnknown() && !data.CredentialKind.Equal(state.CredentialKind) {
        requestDataMap["credentialKind"] = data.CredentialKind.ValueString()
    }
    if !data.ApplicationCredentialId.IsUnknown() && !state.ApplicationCredentialId.IsUnknown() && !data.ApplicationCredentialId.Equal(state.ApplicationCredentialId) {
        requestDataMap["applicationCredentialId"] = data.ApplicationCredentialId.ValueString()
    }
    if !data.KeystoneUsername.IsUnknown() && !state.KeystoneUsername.IsUnknown() && !data.KeystoneUsername.Equal(state.KeystoneUsername) {
        requestDataMap["keystoneUsername"] = data.KeystoneUsername.ValueString()
    }
    if !data.UserDomainName.IsUnknown() && !state.UserDomainName.IsUnknown() && !data.UserDomainName.Equal(state.UserDomainName) {
        requestDataMap["userDomainName"] = data.UserDomainName.ValueString()
    }
    if !data.ProjectName.IsUnknown() && !state.ProjectName.IsUnknown() && !data.ProjectName.Equal(state.ProjectName) {
        requestDataMap["projectName"] = data.ProjectName.ValueString()
    }
    if !data.ProjectDomainName.IsUnknown() && !state.ProjectDomainName.IsUnknown() && !data.ProjectDomainName.Equal(state.ProjectDomainName) {
        requestDataMap["projectDomainName"] = data.ProjectDomainName.ValueString()
    }
    if !data.CredentialSecret.IsUnknown() && !state.CredentialSecret.IsUnknown() && !data.CredentialSecret.Equal(state.CredentialSecret) {
        requestDataMap["credentialSecret"] = data.CredentialSecret.ValueString()
    }
    if !data.CollectionProbeId.IsUnknown() && !state.CollectionProbeId.IsUnknown() && !data.CollectionProbeId.Equal(state.CollectionProbeId) {
        requestDataMap["collectionProbeId"] = data.CollectionProbeId.ValueString()
    }
    if !data.TrustedCertificateFingerprints.IsUnknown() && !state.TrustedCertificateFingerprints.IsUnknown() && !data.TrustedCertificateFingerprints.Equal(state.TrustedCertificateFingerprints) {
        var trustedcertificatefingerprintsData interface{}
        if err := json.Unmarshal([]byte(data.TrustedCertificateFingerprints.ValueString()), &trustedcertificatefingerprintsData); err == nil {
            requestDataMap["trustedCertificateFingerprints"] = trustedcertificatefingerprintsData
        } else {
            requestDataMap["trustedCertificateFingerprints"] = data.TrustedCertificateFingerprints.ValueString()
        }
    }
    if !data.CollectionIntervalInMinutes.IsUnknown() && !state.CollectionIntervalInMinutes.IsUnknown() && !data.CollectionIntervalInMinutes.Equal(state.CollectionIntervalInMinutes) {
        requestDataMap["collectionIntervalInMinutes"] = r.bigFloatToFloat64(data.CollectionIntervalInMinutes.ValueBigFloat())
    }
    if !data.CollectionOptions.IsUnknown() && !state.CollectionOptions.IsUnknown() && !data.CollectionOptions.Equal(state.CollectionOptions) {
        var collectionoptionsData interface{}
        if err := json.Unmarshal([]byte(data.CollectionOptions.ValueString()), &collectionoptionsData); err == nil {
            requestDataMap["collectionOptions"] = collectionoptionsData
        } else {
            requestDataMap["collectionOptions"] = data.CollectionOptions.ValueString()
        }
    }
    if !data.InstanceMetricsOptions.IsUnknown() && !state.InstanceMetricsOptions.IsUnknown() && !data.InstanceMetricsOptions.Equal(state.InstanceMetricsOptions) {
        var instancemetricsoptionsData interface{}
        if err := json.Unmarshal([]byte(data.InstanceMetricsOptions.ValueString()), &instancemetricsoptionsData); err == nil {
            requestDataMap["instanceMetricsOptions"] = instancemetricsoptionsData
        } else {
            requestDataMap["instanceMetricsOptions"] = data.InstanceMetricsOptions.ValueString()
        }
    }
    if !data.PriceList.IsUnknown() && !state.PriceList.IsUnknown() && !data.PriceList.Equal(state.PriceList) {
        var pricelistData interface{}
        if err := json.Unmarshal([]byte(data.PriceList.ValueString()), &pricelistData); err == nil {
            requestDataMap["priceList"] = pricelistData
        } else {
            requestDataMap["priceList"] = data.PriceList.ValueString()
        }
    }
    if !data.IsAiInvestigationEnabled.IsUnknown() && !state.IsAiInvestigationEnabled.IsUnknown() && !data.IsAiInvestigationEnabled.Equal(state.IsAiInvestigationEnabled) {
        requestDataMap["isAiInvestigationEnabled"] = data.IsAiInvestigationEnabled.ValueBool()
    }
    if !data.AiRemediationMode.IsUnknown() && !state.AiRemediationMode.IsUnknown() && !data.AiRemediationMode.Equal(state.AiRemediationMode) {
        requestDataMap["aiRemediationMode"] = data.AiRemediationMode.ValueString()
    }
    if !data.AiCommandAllowlist.IsUnknown() && !state.AiCommandAllowlist.IsUnknown() && !data.AiCommandAllowlist.Equal(state.AiCommandAllowlist) {
        var aicommandallowlistData interface{}
        if err := json.Unmarshal([]byte(data.AiCommandAllowlist.ValueString()), &aicommandallowlistData); err == nil {
            requestDataMap["aiCommandAllowlist"] = aicommandallowlistData
        } else {
            requestDataMap["aiCommandAllowlist"] = data.AiCommandAllowlist.ValueString()
        }
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(openstackCloudRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/openstack-cloud/" + data.Id.ValueString() + "", openstackCloudRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update openstack_cloud, got error: %s", err))
            return
        }

        // Parse the update response
        var openstackCloudResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &openstackCloudResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update openstack_cloud: %s", err))
            return
        }
        _ = openstackCloudResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "isArchived": true,
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
        "collectionProbeId": true,
        "trustedCertificateFingerprints": true,
        "collectionIntervalInMinutes": true,
        "collectionOptions": true,
        "instanceMetricsOptions": true,
        "priceList": true,
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
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "isCredentialSet": true,
        "credentialsUpdatedAt": true,
        "capacitySnapshot": true,
        "spendSnapshot": true,
        "usageReadAt": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "aiAgentResetAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/openstack-cloud/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read openstack_cloud after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read openstack_cloud after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    data.RetainTelemetryDataForDays = numberFromAPI(data.RetainTelemetryDataForDays, dataMap["retainTelemetryDataForDays"])
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["telemetryRetentionConfig"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["collectionMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionMethod = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionMethod = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionMethod = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionMethod"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionMethod = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionMethod = types.StringNull()
        }
    } else if val, ok := dataMap["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := dataMap["authUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AuthUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AuthUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AuthUrl = types.StringValue(string(jsonBytes))
            } else {
                data.AuthUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AuthUrl = types.StringValue(string(jsonBytes))
            } else {
                data.AuthUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if arr, ok := dataMap["authUrl"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AuthUrl = types.StringValue(string(jsonBytes))
        } else {
            data.AuthUrl = types.StringNull()
        }
    } else if val, ok := dataMap["authUrl"].(string); ok {
        data.AuthUrl = types.StringValue(val)
    } else {
        data.AuthUrl = types.StringNull()
    }
    if obj, ok := dataMap["endpointInterface"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndpointInterface = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndpointInterface = types.StringValue(string(jsonBytes))
            } else {
                data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndpointInterface = types.StringValue(string(jsonBytes))
            } else {
                data.EndpointInterface = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if arr, ok := dataMap["endpointInterface"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.EndpointInterface = types.StringValue(string(jsonBytes))
        } else {
            data.EndpointInterface = types.StringNull()
        }
    } else if val, ok := dataMap["endpointInterface"].(string); ok {
        data.EndpointInterface = types.StringValue(val)
    } else {
        data.EndpointInterface = types.StringNull()
    }
    if obj, ok := dataMap["collectionRegions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionRegions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionRegions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionRegions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionRegions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionRegions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionRegions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionRegions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionRegions"].(string); ok {
        data.CollectionRegions = NewJSONSubsetValue(val)
    } else {
        data.CollectionRegions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["credentialKind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CredentialKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CredentialKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CredentialKind = types.StringValue(string(jsonBytes))
            } else {
                data.CredentialKind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CredentialKind = types.StringValue(string(jsonBytes))
            } else {
                data.CredentialKind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if arr, ok := dataMap["credentialKind"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CredentialKind = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialKind = types.StringNull()
        }
    } else if val, ok := dataMap["credentialKind"].(string); ok {
        data.CredentialKind = types.StringValue(val)
    } else {
        data.CredentialKind = types.StringNull()
    }
    if obj, ok := dataMap["applicationCredentialId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ApplicationCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
            } else {
                data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
            } else {
                data.ApplicationCredentialId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if arr, ok := dataMap["applicationCredentialId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ApplicationCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.ApplicationCredentialId = types.StringNull()
        }
    } else if val, ok := dataMap["applicationCredentialId"].(string); ok {
        data.ApplicationCredentialId = types.StringValue(val)
    } else {
        data.ApplicationCredentialId = types.StringNull()
    }
    if obj, ok := dataMap["keystoneUsername"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.KeystoneUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.KeystoneUsername = types.StringValue(string(jsonBytes))
            } else {
                data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.KeystoneUsername = types.StringValue(string(jsonBytes))
            } else {
                data.KeystoneUsername = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if arr, ok := dataMap["keystoneUsername"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.KeystoneUsername = types.StringValue(string(jsonBytes))
        } else {
            data.KeystoneUsername = types.StringNull()
        }
    } else if val, ok := dataMap["keystoneUsername"].(string); ok {
        data.KeystoneUsername = types.StringValue(val)
    } else {
        data.KeystoneUsername = types.StringNull()
    }
    if obj, ok := dataMap["userDomainName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UserDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UserDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UserDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.UserDomainName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UserDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.UserDomainName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if arr, ok := dataMap["userDomainName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.UserDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.UserDomainName = types.StringNull()
        }
    } else if val, ok := dataMap["userDomainName"].(string); ok {
        data.UserDomainName = types.StringValue(val)
    } else {
        data.UserDomainName = types.StringNull()
    }
    if obj, ok := dataMap["projectName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProjectName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProjectName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProjectName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if arr, ok := dataMap["projectName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ProjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectName = types.StringNull()
        }
    } else if val, ok := dataMap["projectName"].(string); ok {
        data.ProjectName = types.StringValue(val)
    } else {
        data.ProjectName = types.StringNull()
    }
    if obj, ok := dataMap["projectDomainName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProjectDomainName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProjectDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProjectDomainName = types.StringValue(string(jsonBytes))
            } else {
                data.ProjectDomainName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if arr, ok := dataMap["projectDomainName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ProjectDomainName = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectDomainName = types.StringNull()
        }
    } else if val, ok := dataMap["projectDomainName"].(string); ok {
        data.ProjectDomainName = types.StringValue(val)
    } else {
        data.ProjectDomainName = types.StringNull()
    }
    if obj, ok := dataMap["collectionProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionProbeId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := dataMap["trustedCertificateFingerprints"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprints = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["trustedCertificateFingerprints"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TrustedCertificateFingerprints = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprints = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["trustedCertificateFingerprints"].(string); ok {
        data.TrustedCertificateFingerprints = NewJSONSubsetValue(val)
    } else {
        data.TrustedCertificateFingerprints = NewJSONSubsetNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, dataMap["collectionIntervalInMinutes"])
    if obj, ok := dataMap["collectionOptions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionOptions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionOptions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionOptions"].(string); ok {
        data.CollectionOptions = NewJSONSubsetValue(val)
    } else {
        data.CollectionOptions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["instanceMetricsOptions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InstanceMetricsOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InstanceMetricsOptions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.InstanceMetricsOptions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["instanceMetricsOptions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.InstanceMetricsOptions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.InstanceMetricsOptions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["instanceMetricsOptions"].(string); ok {
        data.InstanceMetricsOptions = NewJSONSubsetValue(val)
    } else {
        data.InstanceMetricsOptions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["priceList"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PriceList = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PriceList = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PriceList = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PriceList = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PriceList = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PriceList = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PriceList = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["priceList"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PriceList = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PriceList = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["priceList"].(string); ok {
        data.PriceList = NewJSONSubsetValue(val)
    } else {
        data.PriceList = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["otelCollectorStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
            } else {
                data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["otelCollectorStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := dataMap["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := dataMap["agentVersion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AgentVersion = types.StringValue(string(jsonBytes))
            } else {
                data.AgentVersion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if arr, ok := dataMap["agentVersion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := dataMap["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    data.RegionCount = numberFromAPI(data.RegionCount, dataMap["regionCount"])
    data.AvailabilityZoneCount = numberFromAPI(data.AvailabilityZoneCount, dataMap["availabilityZoneCount"])
    data.HypervisorCount = numberFromAPI(data.HypervisorCount, dataMap["hypervisorCount"])
    data.InstanceCount = numberFromAPI(data.InstanceCount, dataMap["instanceCount"])
    data.ActiveInstanceCount = numberFromAPI(data.ActiveInstanceCount, dataMap["activeInstanceCount"])
    data.VolumeCount = numberFromAPI(data.VolumeCount, dataMap["volumeCount"])
    data.NetworkCount = numberFromAPI(data.NetworkCount, dataMap["networkCount"])
    data.RouterCount = numberFromAPI(data.RouterCount, dataMap["routerCount"])
    data.FloatingIpCount = numberFromAPI(data.FloatingIpCount, dataMap["floatingIpCount"])
    data.ProjectCount = numberFromAPI(data.ProjectCount, dataMap["projectCount"])
    data.ServiceCount = numberFromAPI(data.ServiceCount, dataMap["serviceCount"])
    data.DownServiceCount = numberFromAPI(data.DownServiceCount, dataMap["downServiceCount"])
    data.VcpuCapacity = numberFromAPI(data.VcpuCapacity, dataMap["vcpuCapacity"])
    data.VcpuAllocated = numberFromAPI(data.VcpuAllocated, dataMap["vcpuAllocated"])
    data.MemoryCapacityBytes = numberFromAPI(data.MemoryCapacityBytes, dataMap["memoryCapacityBytes"])
    data.MemoryAllocatedBytes = numberFromAPI(data.MemoryAllocatedBytes, dataMap["memoryAllocatedBytes"])
    data.StorageCapacityBytes = numberFromAPI(data.StorageCapacityBytes, dataMap["storageCapacityBytes"])
    data.StorageUsedBytes = numberFromAPI(data.StorageUsedBytes, dataMap["storageUsedBytes"])
    if val, ok := dataMap["isAiInvestigationEnabled"].(bool); ok {
        data.IsAiInvestigationEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["aiRemediationMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiRemediationMode = types.StringValue(string(jsonBytes))
            } else {
                data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiRemediationMode = types.StringValue(string(jsonBytes))
            } else {
                data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if arr, ok := dataMap["aiRemediationMode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if val, ok := dataMap["aiRemediationMode"].(string); ok {
        data.AiRemediationMode = types.StringValue(val)
    } else {
        data.AiRemediationMode = types.StringNull()
    }
    if obj, ok := dataMap["aiCommandAllowlist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiCommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiCommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiCommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["aiCommandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiCommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["aiCommandAllowlist"].(string); ok {
        data.AiCommandAllowlist = NewJSONSubsetValue(val)
    } else {
        data.AiCommandAllowlist = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ArchivedAt = NewRFC3339Value(val)
        } else {
            data.ArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["archivedAt"].(string); ok && val != "" {
        data.ArchivedAt = NewRFC3339Value(val)
    } else {
        data.ArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["archivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["archivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isCredentialSet"].(bool); ok {
        data.IsCredentialSet = types.BoolValue(val)
    }
    if obj, ok := dataMap["credentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CredentialsUpdatedAt = NewRFC3339Value(val)
        } else {
            data.CredentialsUpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["credentialsUpdatedAt"].(string); ok && val != "" {
        data.CredentialsUpdatedAt = NewRFC3339Value(val)
    } else {
        data.CredentialsUpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["capacitySnapshot"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CapacitySnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CapacitySnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CapacitySnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["capacitySnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CapacitySnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CapacitySnapshot = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["capacitySnapshot"].(string); ok {
        data.CapacitySnapshot = NewJSONSubsetValue(val)
    } else {
        data.CapacitySnapshot = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["spendSnapshot"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpendSnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SpendSnapshot = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SpendSnapshot = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["spendSnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SpendSnapshot = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SpendSnapshot = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["spendSnapshot"].(string); ok {
        data.SpendSnapshot = NewJSONSubsetValue(val)
    } else {
        data.SpendSnapshot = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["usageReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UsageReadAt = NewRFC3339Value(val)
        } else {
            data.UsageReadAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["usageReadAt"].(string); ok && val != "" {
        data.UsageReadAt = NewRFC3339Value(val)
    } else {
        data.UsageReadAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["collectionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionStatus"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["collectionStatus"].(string); ok {
        data.CollectionStatus = types.StringValue(val)
    } else {
        data.CollectionStatus = types.StringNull()
    }
    if obj, ok := dataMap["collectionErrorCode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionErrorCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionErrorCode = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionErrorCode = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionErrorCode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionErrorCode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionErrorCode = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionErrorCode = types.StringNull()
        }
    } else if val, ok := dataMap["collectionErrorCode"].(string); ok {
        data.CollectionErrorCode = types.StringValue(val)
    } else {
        data.CollectionErrorCode = types.StringNull()
    }
    if obj, ok := dataMap["collectionError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionError = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionError = types.StringValue(string(jsonBytes))
            } else {
                data.CollectionError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if arr, ok := dataMap["collectionError"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionError = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionError = types.StringNull()
        }
    } else if val, ok := dataMap["collectionError"].(string); ok {
        data.CollectionError = types.StringValue(val)
    } else {
        data.CollectionError = types.StringNull()
    }
    if obj, ok := dataMap["presentedCertificate"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PresentedCertificate = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PresentedCertificate = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.PresentedCertificate = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["presentedCertificate"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PresentedCertificate = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.PresentedCertificate = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["presentedCertificate"].(string); ok {
        data.PresentedCertificate = NewJSONSubsetValue(val)
    } else {
        data.PresentedCertificate = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["collectionSummary"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionSummary = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CollectionSummary = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CollectionSummary = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionSummary = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["collectionSummary"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CollectionSummary = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CollectionSummary = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["collectionSummary"].(string); ok {
        data.CollectionSummary = NewJSONSubsetValue(val)
    } else {
        data.CollectionSummary = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["nextCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextCollectionAt = NewRFC3339Value(val)
        } else {
            data.NextCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextCollectionAt"].(string); ok && val != "" {
        data.NextCollectionAt = NewRFC3339Value(val)
    } else {
        data.NextCollectionAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastCollectionAt = NewRFC3339Value(val)
        } else {
            data.LastCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastCollectionAt"].(string); ok && val != "" {
        data.LastCollectionAt = NewRFC3339Value(val)
    } else {
        data.LastCollectionAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastSuccessfulCollectionAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSuccessfulCollectionAt = NewRFC3339Value(val)
        } else {
            data.LastSuccessfulCollectionAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSuccessfulCollectionAt"].(string); ok && val != "" {
        data.LastSuccessfulCollectionAt = NewRFC3339Value(val)
    } else {
        data.LastSuccessfulCollectionAt = NewRFC3339Null()
    }
    data.CollectionSettingsVersion = numberFromAPI(data.CollectionSettingsVersion, dataMap["collectionSettingsVersion"])
    if obj, ok := dataMap["aiAccessLastVerifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAccessLastVerifiedAt = NewRFC3339Value(val)
        } else {
            data.AiAccessLastVerifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAccessLastVerifiedAt"].(string); ok && val != "" {
        data.AiAccessLastVerifiedAt = NewRFC3339Value(val)
    } else {
        data.AiAccessLastVerifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiAccessLastError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiAccessLastError = types.StringValue(string(jsonBytes))
            } else {
                data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiAccessLastError = types.StringValue(string(jsonBytes))
            } else {
                data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if arr, ok := dataMap["aiAccessLastError"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if val, ok := dataMap["aiAccessLastError"].(string); ok {
        data.AiAccessLastError = types.StringValue(val)
    } else {
        data.AiAccessLastError = types.StringNull()
    }
    if obj, ok := dataMap["aiAccessConfiguredAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAccessConfiguredAt = NewRFC3339Value(val)
        } else {
            data.AiAccessConfiguredAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAccessConfiguredAt"].(string); ok && val != "" {
        data.AiAccessConfiguredAt = NewRFC3339Value(val)
    } else {
        data.AiAccessConfiguredAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiAgentResetAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiAgentResetAt = NewRFC3339Value(val)
        } else {
            data.AiAgentResetAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiAgentResetAt"].(string); ok && val != "" {
        data.AiAgentResetAt = NewRFC3339Value(val)
    } else {
        data.AiAgentResetAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // What this update wrote is what the configuration says now (see importsettle.go).
    resp.Diagnostics.Append(keepImportedFromServer(ctx, resp.Private, writtenFromConfig)...)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpenstackCloudResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data OpenstackCloudResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/openstack-cloud/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete openstack_cloud, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete openstack_cloud: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *OpenstackCloudResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
    // These hold the record as OneUptime keeps it, keys it fills in itself
    // included, until an update writes them (see importsettle.go).
    resp.Diagnostics.Append(markImportedFromServer(ctx, resp.Private, []string{"ai_command_allowlist", "collection_options", "collection_regions", "instance_metrics_options", "price_list", "telemetry_retention_config", "trusted_certificate_fingerprints"})...)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *OpenstackCloudResource) keepPlannedValues(data *OpenstackCloudResourceModel, plan *OpenstackCloudResourceModel, config *OpenstackCloudResourceModel) {
    if config.Description.IsNull() && !plan.Description.IsUnknown() {
        data.Description = plan.Description
    }
    if config.IsArchived.IsNull() && !plan.IsArchived.IsUnknown() {
        data.IsArchived = plan.IsArchived
    }
    if config.Labels.IsNull() && !plan.Labels.IsUnknown() {
        data.Labels = plan.Labels
    }
    if config.RetainTelemetryDataForDays.IsNull() && !plan.RetainTelemetryDataForDays.IsUnknown() {
        data.RetainTelemetryDataForDays = plan.RetainTelemetryDataForDays
    }
    if config.TelemetryRetentionConfig.IsNull() && !plan.TelemetryRetentionConfig.IsUnknown() {
        data.TelemetryRetentionConfig = plan.TelemetryRetentionConfig
    }
    if config.CollectionMethod.IsNull() && !plan.CollectionMethod.IsUnknown() {
        data.CollectionMethod = plan.CollectionMethod
    }
    if config.AuthUrl.IsNull() && !plan.AuthUrl.IsUnknown() {
        data.AuthUrl = plan.AuthUrl
    }
    if config.EndpointInterface.IsNull() && !plan.EndpointInterface.IsUnknown() {
        data.EndpointInterface = plan.EndpointInterface
    }
    if config.CollectionRegions.IsNull() && !plan.CollectionRegions.IsUnknown() {
        data.CollectionRegions = plan.CollectionRegions
    }
    if config.CredentialKind.IsNull() && !plan.CredentialKind.IsUnknown() {
        data.CredentialKind = plan.CredentialKind
    }
    if config.ApplicationCredentialId.IsNull() && !plan.ApplicationCredentialId.IsUnknown() {
        data.ApplicationCredentialId = plan.ApplicationCredentialId
    }
    if config.KeystoneUsername.IsNull() && !plan.KeystoneUsername.IsUnknown() {
        data.KeystoneUsername = plan.KeystoneUsername
    }
    if config.UserDomainName.IsNull() && !plan.UserDomainName.IsUnknown() {
        data.UserDomainName = plan.UserDomainName
    }
    if config.ProjectName.IsNull() && !plan.ProjectName.IsUnknown() {
        data.ProjectName = plan.ProjectName
    }
    if config.ProjectDomainName.IsNull() && !plan.ProjectDomainName.IsUnknown() {
        data.ProjectDomainName = plan.ProjectDomainName
    }
    if config.CollectionProbeId.IsNull() && !plan.CollectionProbeId.IsUnknown() {
        data.CollectionProbeId = plan.CollectionProbeId
    }
    if config.TrustedCertificateFingerprints.IsNull() && !plan.TrustedCertificateFingerprints.IsUnknown() {
        data.TrustedCertificateFingerprints = plan.TrustedCertificateFingerprints
    }
    if config.CollectionIntervalInMinutes.IsNull() && !plan.CollectionIntervalInMinutes.IsUnknown() {
        data.CollectionIntervalInMinutes = plan.CollectionIntervalInMinutes
    }
    if config.CollectionOptions.IsNull() && !plan.CollectionOptions.IsUnknown() {
        data.CollectionOptions = plan.CollectionOptions
    }
    if config.InstanceMetricsOptions.IsNull() && !plan.InstanceMetricsOptions.IsUnknown() {
        data.InstanceMetricsOptions = plan.InstanceMetricsOptions
    }
    if config.PriceList.IsNull() && !plan.PriceList.IsUnknown() {
        data.PriceList = plan.PriceList
    }
    if config.OtelCollectorStatus.IsNull() && !plan.OtelCollectorStatus.IsUnknown() {
        data.OtelCollectorStatus = plan.OtelCollectorStatus
    }
    if config.AgentVersion.IsNull() && !plan.AgentVersion.IsUnknown() {
        data.AgentVersion = plan.AgentVersion
    }
    if config.LastSeenAt.IsNull() && !plan.LastSeenAt.IsUnknown() {
        data.LastSeenAt = plan.LastSeenAt
    }
    if config.RegionCount.IsNull() && !plan.RegionCount.IsUnknown() {
        data.RegionCount = plan.RegionCount
    }
    if config.AvailabilityZoneCount.IsNull() && !plan.AvailabilityZoneCount.IsUnknown() {
        data.AvailabilityZoneCount = plan.AvailabilityZoneCount
    }
    if config.HypervisorCount.IsNull() && !plan.HypervisorCount.IsUnknown() {
        data.HypervisorCount = plan.HypervisorCount
    }
    if config.InstanceCount.IsNull() && !plan.InstanceCount.IsUnknown() {
        data.InstanceCount = plan.InstanceCount
    }
    if config.ActiveInstanceCount.IsNull() && !plan.ActiveInstanceCount.IsUnknown() {
        data.ActiveInstanceCount = plan.ActiveInstanceCount
    }
    if config.VolumeCount.IsNull() && !plan.VolumeCount.IsUnknown() {
        data.VolumeCount = plan.VolumeCount
    }
    if config.NetworkCount.IsNull() && !plan.NetworkCount.IsUnknown() {
        data.NetworkCount = plan.NetworkCount
    }
    if config.RouterCount.IsNull() && !plan.RouterCount.IsUnknown() {
        data.RouterCount = plan.RouterCount
    }
    if config.FloatingIpCount.IsNull() && !plan.FloatingIpCount.IsUnknown() {
        data.FloatingIpCount = plan.FloatingIpCount
    }
    if config.ProjectCount.IsNull() && !plan.ProjectCount.IsUnknown() {
        data.ProjectCount = plan.ProjectCount
    }
    if config.ServiceCount.IsNull() && !plan.ServiceCount.IsUnknown() {
        data.ServiceCount = plan.ServiceCount
    }
    if config.DownServiceCount.IsNull() && !plan.DownServiceCount.IsUnknown() {
        data.DownServiceCount = plan.DownServiceCount
    }
    if config.VcpuCapacity.IsNull() && !plan.VcpuCapacity.IsUnknown() {
        data.VcpuCapacity = plan.VcpuCapacity
    }
    if config.VcpuAllocated.IsNull() && !plan.VcpuAllocated.IsUnknown() {
        data.VcpuAllocated = plan.VcpuAllocated
    }
    if config.MemoryCapacityBytes.IsNull() && !plan.MemoryCapacityBytes.IsUnknown() {
        data.MemoryCapacityBytes = plan.MemoryCapacityBytes
    }
    if config.MemoryAllocatedBytes.IsNull() && !plan.MemoryAllocatedBytes.IsUnknown() {
        data.MemoryAllocatedBytes = plan.MemoryAllocatedBytes
    }
    if config.StorageCapacityBytes.IsNull() && !plan.StorageCapacityBytes.IsUnknown() {
        data.StorageCapacityBytes = plan.StorageCapacityBytes
    }
    if config.StorageUsedBytes.IsNull() && !plan.StorageUsedBytes.IsUnknown() {
        data.StorageUsedBytes = plan.StorageUsedBytes
    }
    if config.IsAiInvestigationEnabled.IsNull() && !plan.IsAiInvestigationEnabled.IsUnknown() {
        data.IsAiInvestigationEnabled = plan.IsAiInvestigationEnabled
    }
    if config.AiRemediationMode.IsNull() && !plan.AiRemediationMode.IsUnknown() {
        data.AiRemediationMode = plan.AiRemediationMode
    }
    if config.AiCommandAllowlist.IsNull() && !plan.AiCommandAllowlist.IsUnknown() {
        data.AiCommandAllowlist = plan.AiCommandAllowlist
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *OpenstackCloudResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *OpenstackCloudResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *OpenstackCloudResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *OpenstackCloudResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *OpenstackCloudResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *OpenstackCloudResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *OpenstackCloudResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *OpenstackCloudResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
