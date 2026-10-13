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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VcenterResource{}
var _ resource.ResourceWithImportState = &VcenterResource{}
var _ resource.ResourceWithMoveState = &VcenterResource{}

func NewVcenterResource() resource.Resource {
    return &VcenterResource{}
}

// NewVcenterLegacyResource registers this resource under the name it
// had before resource type names kept mixed-case words whole,
// oneuptime_v_center, so configurations that use it keep
// working. Deprecated: plans that use it say so.
func NewVcenterLegacyResource() resource.Resource {
    return &VcenterResource{isLegacyAlias: true}
}

// VcenterResource defines the resource implementation.
type VcenterResource struct {
    client *Client
    // Registered under the name this resource had before (see NewVcenterLegacyResource).
    isLegacyAlias bool
}

// VcenterResourceModel describes the resource data model.
type VcenterResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig JSONSubsetValue `tfsdk:"telemetry_retention_config"`
    CollectionMethod types.String `tfsdk:"collection_method"`
    VcenterUrl types.String `tfsdk:"vcenter_url"`
    VcenterUsername types.String `tfsdk:"vcenter_username"`
    VcenterPassword types.String `tfsdk:"vcenter_password"`
    CollectionProbeId types.String `tfsdk:"collection_probe_id"`
    TrustedCertificateFingerprint types.String `tfsdk:"trusted_certificate_fingerprint"`
    CollectionIntervalInMinutes types.Number `tfsdk:"collection_interval_in_minutes"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
    DatacenterCount types.Number `tfsdk:"datacenter_count"`
    ClusterCount types.Number `tfsdk:"cluster_count"`
    HostCount types.Number `tfsdk:"host_count"`
    VmCount types.Number `tfsdk:"vm_count"`
    PoweredOnVmCount types.Number `tfsdk:"powered_on_vm_count"`
    DatastoreCount types.Number `tfsdk:"datastore_count"`
    ResourcePoolCount types.Number `tfsdk:"resource_pool_count"`
    DatastoreCapacityBytes types.Number `tfsdk:"datastore_capacity_bytes"`
    DatastoreUsedBytes types.Number `tfsdk:"datastore_used_bytes"`
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiCommandAllowlist JSONSubsetValue `tfsdk:"ai_command_allowlist"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ArchivedAt RFC3339Value `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    AiAccessLastVerifiedAt RFC3339Value `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt RFC3339Value `tfsdk:"ai_access_configured_at"`
    IsVCenterPasswordSet types.Bool `tfsdk:"is_v_center_password_set"`
    VcenterCredentialsUpdatedAt RFC3339Value `tfsdk:"vcenter_credentials_updated_at"`
    CollectionStatus types.String `tfsdk:"collection_status"`
    CollectionErrorCode types.String `tfsdk:"collection_error_code"`
    CollectionError types.String `tfsdk:"collection_error"`
    PresentedCertificate JSONSubsetValue `tfsdk:"presented_certificate"`
    CollectionSummary JSONSubsetValue `tfsdk:"collection_summary"`
    NextCollectionAt RFC3339Value `tfsdk:"next_collection_at"`
    LastCollectionAt RFC3339Value `tfsdk:"last_collection_at"`
    LastSuccessfulCollectionAt RFC3339Value `tfsdk:"last_successful_collection_at"`
    CollectionSettingsVersion types.Number `tfsdk:"collection_settings_version"`
}

func (r *VcenterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    if r.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_v_center"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_vcenter"
}

func (r *VcenterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
    if r.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_v_center has been renamed to oneuptime_vcenter. Rename the resource in your configuration and add a moved block (from = oneuptime_v_center.<name>, to = oneuptime_vcenter.<name>) to keep the existing vcenter; the old name keeps working until then."
    }
}

func (r *VcenterResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "vSphere endpoints (a vCenter Server, or a standalone ESXi host) that are being monitored in this project. A vCenter is collected by one of your OneUptime probes with a read-only account saved on it (no agent to run), or auto-discovered when the OneUptime VMware Agent sends its metrics.",

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
                MarkdownDescription: "Name of this vCenter (a vCenter Server, or a standalone ESXi host). This is the join key — it must match the vmware.vcenter.name OTel resource attribute stamped by the OneUptime VMware Agent (VMWARE_VCENTER_NAME).",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this vCenter.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this vCenter archived? Archived vCenters are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this vCenter. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this vCenter (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the vCenter default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "collection_method": schema.StringAttribute{
                MarkdownDescription: "How this vCenter's data reaches OneUptime. Probe: one of your OneUptime probes logs in to vCenter with the read-only account saved here and collects it - no agent to install. Agent: the OneUptime VMware Agent you run sends it. Switching to Agent forgets the saved password.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Agent"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "vcenter_url": schema.StringAttribute{
                MarkdownDescription: "The HTTPS address of vCenter Server (or a standalone ESXi host) the probe connects to, such as https://vcsa.example.com. Stored as https://host[:port]; a path such as /ui or /sdk is dropped. Changing it asks for the password again.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "vcenter_username": schema.StringAttribute{
                MarkdownDescription: "The vSphere user the probe logs in as, with its domain - such as oneuptime@vsphere.local. Give it the built-in Read-Only role on the top-level vCenter object, propagated to children.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "vcenter_password": schema.StringAttribute{
                MarkdownDescription: "The password of the vCenter user. Write-only: encrypted at rest, never returned by the API, and sent only to the probe that collects this vCenter. Leave it out of an update to keep the saved one; changing the address, the probe or the trusted certificate needs it again.",
                Optional: true,
                Sensitive: true,
            },
            "collection_probe_id": schema.StringAttribute{
                MarkdownDescription: "ID of the OneUptime probe that collects this vCenter. Changing it asks for the password again. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "trusted_certificate_fingerprint": schema.StringAttribute{
                MarkdownDescription: "The SHA-256 fingerprint of the one certificate the probe accepts from this vCenter, for a vCenter whose certificate is not from a public authority (vCenter's own VMCA certificate, by default). Empty: the certificate must be signed by an authority the probe's machine trusts. Verification is never skipped. Trusting a different certificate asks for the password again.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "collection_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often the probe collects this vCenter, in whole minutes from 1 to 60. Two minutes is the VMware agent's own default; collect large vCenters less often to go easier on vCenter.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(2)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
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
                MarkdownDescription: "Version of the OneUptime VMware Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When metrics were last received from this vCenter.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "datacenter_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of vSphere datacenters reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries datacenter metrics.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "cluster_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of vSphere clusters reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries cluster metrics.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "host_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of ESXi hosts reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries host metrics.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "vm_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of virtual machines (excluding VM templates) reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries virtual machine metrics.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "powered_on_vm_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of virtual machines inferred to be powered on (the vcenter receiver emits CPU metrics only for powered-on VMs). Rendered as 'VMs X/Y powered on' next to vmCount.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "datastore_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of datastores reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries datastore metrics.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "resource_pool_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of resource pools reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries resource pool metrics.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "datastore_capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached total capacity in bytes summed over every datastore reported through this vCenter (used + available vcenter.datastore.disk.usage). The denominator for the storage-used bar on the overview page. Stored as bigint.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "datastore_used_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached used space in bytes summed over every datastore reported through this vCenter (vcenter.datastore.disk.usage{disk_state=used}). Stored as bigint.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (govc about, ls, vm.info, host.info, events) on this vCenter, through its VMware AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the vCenter can turn it on or off.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this vCenter. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the vCenter's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the vCenter can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Disabled"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of command patterns that Automatic mode may run on this vCenter without approval even though they are riskier changes. Each pattern is one command line for this vCenter's agent (govc) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this vCenter is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the vCenter can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
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
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this vCenter archived?",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this vCenter through its VMware AI agent. Set by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this vCenter, kept until the next successful command. Set by the server.",
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this vCenter was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a VMware AI agent that registers later never overwrites a setting an operator chose.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "is_v_center_password_set": schema.BoolAttribute{
                MarkdownDescription: "Whether a vCenter password is saved. Set by the server from the password itself.",
                Computed: true,
            },
            "vcenter_credentials_updated_at": schema.StringAttribute{
                MarkdownDescription: "When the vCenter user name or password was last changed. Set by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "collection_status": schema.StringAttribute{
                MarkdownDescription: "Where probe collection stands: Pending (the probe has not tried the current settings yet), Succeeded or Failed. Written by the server and the probe.",
                Computed: true,
            },
            "collection_error_code": schema.StringAttribute{
                MarkdownDescription: "Why the last collection failed, as a code the dashboard turns into a next step (InvalidLogin, UntrustedCertificate, ConnectionTimedOut, ...). Written by the probe.",
                Computed: true,
            },
            "collection_error": schema.StringAttribute{
                MarkdownDescription: "What went wrong in the last collection, in words. Written by the probe.",
                Computed: true,
            },
            "presented_certificate": schema.StringAttribute{
                MarkdownDescription: "The certificate vCenter presented when the last collection did not trust it: its SHA-256 fingerprint, subject, issuer and validity, so a person can check it and trust exactly that certificate. Written by the probe. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "collection_summary": schema.StringAttribute{
                MarkdownDescription: "What the last successful collection found: vCenter's product and version, and how many datacenters, clusters, hosts, virtual machines, datastores and resource pools it reported. Written by the probe. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "next_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe collects this vCenter next. Written by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe last tried to collect this vCenter, whether it worked or not. Written by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_successful_collection_at": schema.StringAttribute{
                MarkdownDescription: "When the probe last collected this vCenter successfully. Written by the server.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "collection_settings_version": schema.NumberAttribute{
                MarkdownDescription: "Counts changes to the probe collection settings, so a report from a collection that started before a change is not shown as the current status. Written by the server.",
                Computed: true,
            },
        },
    }
}

// MoveState lets a moved block bring state over from the old name:
//
//     moved {
//       from = oneuptime_v_center.example
//       to   = oneuptime_vcenter.example
//     }
func (r *VcenterResource) MoveState(ctx context.Context) []resource.StateMover {
    if r.isLegacyAlias {
        return nil
    }

    return []resource.StateMover{
        legacyNameStateMover("oneuptime_v_center", r.schemaDefinition()),
    }
}

func (r *VcenterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *VcenterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data VcenterResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config VcenterResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    vcenterRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := vcenterRequest["data"].(map[string]interface{})

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
    if !data.VcenterUrl.IsNull() && !data.VcenterUrl.IsUnknown() {
        requestDataMap["vcenterUrl"] = data.VcenterUrl.ValueString()
    }
    if !data.VcenterUsername.IsNull() && !data.VcenterUsername.IsUnknown() {
        requestDataMap["vcenterUsername"] = data.VcenterUsername.ValueString()
    }
    if !data.VcenterPassword.IsNull() && !data.VcenterPassword.IsUnknown() {
        requestDataMap["vcenterPassword"] = data.VcenterPassword.ValueString()
    }
    if !data.CollectionProbeId.IsNull() && !data.CollectionProbeId.IsUnknown() {
        requestDataMap["collectionProbeId"] = data.CollectionProbeId.ValueString()
    }
    if !data.TrustedCertificateFingerprint.IsNull() && !data.TrustedCertificateFingerprint.IsUnknown() {
        requestDataMap["trustedCertificateFingerprint"] = data.TrustedCertificateFingerprint.ValueString()
    }
    if !data.CollectionIntervalInMinutes.IsNull() && !data.CollectionIntervalInMinutes.IsUnknown() {
        requestDataMap["collectionIntervalInMinutes"] = r.bigFloatToFloat64(data.CollectionIntervalInMinutes.ValueBigFloat())
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/vmware-vcenter", vcenterRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create vcenter, got error: %s", err))
        return
    }

    var vcenterResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &vcenterResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create vcenter: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := vcenterResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := vcenterResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for vcenter did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * vcenter is orphaned server-side — never refreshed, never
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
        "vcenterUrl": true,
        "vcenterUsername": true,
        "collectionProbeId": true,
        "trustedCertificateFingerprint": true,
        "collectionIntervalInMinutes": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "datacenterCount": true,
        "clusterCount": true,
        "hostCount": true,
        "vmCount": true,
        "poweredOnVmCount": true,
        "datastoreCount": true,
        "resourcePoolCount": true,
        "datastoreCapacityBytes": true,
        "datastoreUsedBytes": true,
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "isVCenterPasswordSet": true,
        "vcenterCredentialsUpdatedAt": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/vmware-vcenter/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created vcenter but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created vcenter but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    } else if val, ok := dataMap["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := dataMap["vcenterUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VcenterUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.VcenterUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.VcenterUrl = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.VcenterUrl = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.VcenterUrl = types.StringValue(string(jsonBytes))
        } else {
            data.VcenterUrl = types.StringNull()
        }
    } else if val, ok := dataMap["vcenterUrl"].(string); ok {
        data.VcenterUrl = types.StringValue(val)
    } else {
        data.VcenterUrl = types.StringNull()
    }
    if obj, ok := dataMap["vcenterUsername"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VcenterUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.VcenterUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.VcenterUsername = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.VcenterUsername = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.VcenterUsername = types.StringValue(string(jsonBytes))
        } else {
            data.VcenterUsername = types.StringNull()
        }
    } else if val, ok := dataMap["vcenterUsername"].(string); ok {
        data.VcenterUsername = types.StringValue(val)
    } else {
        data.VcenterUsername = types.StringNull()
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
    } else if val, ok := dataMap["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := dataMap["trustedCertificateFingerprint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrustedCertificateFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprint = types.StringNull()
        }
    } else if val, ok := dataMap["trustedCertificateFingerprint"].(string); ok {
        data.TrustedCertificateFingerprint = types.StringValue(val)
    } else {
        data.TrustedCertificateFingerprint = types.StringNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, dataMap["collectionIntervalInMinutes"])
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
    data.DatacenterCount = numberFromAPI(data.DatacenterCount, dataMap["datacenterCount"])
    data.ClusterCount = numberFromAPI(data.ClusterCount, dataMap["clusterCount"])
    data.HostCount = numberFromAPI(data.HostCount, dataMap["hostCount"])
    data.VmCount = numberFromAPI(data.VmCount, dataMap["vmCount"])
    data.PoweredOnVmCount = numberFromAPI(data.PoweredOnVmCount, dataMap["poweredOnVmCount"])
    data.DatastoreCount = numberFromAPI(data.DatastoreCount, dataMap["datastoreCount"])
    data.ResourcePoolCount = numberFromAPI(data.ResourcePoolCount, dataMap["resourcePoolCount"])
    data.DatastoreCapacityBytes = numberFromAPI(data.DatastoreCapacityBytes, dataMap["datastoreCapacityBytes"])
    data.DatastoreUsedBytes = numberFromAPI(data.DatastoreUsedBytes, dataMap["datastoreUsedBytes"])
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
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
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
    if val, ok := dataMap["isVCenterPasswordSet"].(bool); ok {
        data.IsVCenterPasswordSet = types.BoolValue(val)
    }
    if obj, ok := dataMap["vcenterCredentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.VcenterCredentialsUpdatedAt = NewRFC3339Value(val)
        } else {
            data.VcenterCredentialsUpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["vcenterCredentialsUpdatedAt"].(string); ok && val != "" {
        data.VcenterCredentialsUpdatedAt = NewRFC3339Value(val)
    } else {
        data.VcenterCredentialsUpdatedAt = NewRFC3339Null()
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

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VcenterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data VcenterResourceModel

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
        "vcenterUrl": true,
        "vcenterUsername": true,
        "collectionProbeId": true,
        "trustedCertificateFingerprint": true,
        "collectionIntervalInMinutes": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "datacenterCount": true,
        "clusterCount": true,
        "hostCount": true,
        "vmCount": true,
        "poweredOnVmCount": true,
        "datastoreCount": true,
        "resourcePoolCount": true,
        "datastoreCapacityBytes": true,
        "datastoreUsedBytes": true,
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "isVCenterPasswordSet": true,
        "vcenterCredentialsUpdatedAt": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/vmware-vcenter/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vcenter, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var vcenterResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &vcenterResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse vcenter response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := vcenterResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = vcenterResponse
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
    } else if val, ok := dataMap["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := dataMap["vcenterUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VcenterUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.VcenterUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.VcenterUrl = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.VcenterUrl = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.VcenterUrl = types.StringValue(string(jsonBytes))
        } else {
            data.VcenterUrl = types.StringNull()
        }
    } else if val, ok := dataMap["vcenterUrl"].(string); ok {
        data.VcenterUrl = types.StringValue(val)
    } else {
        data.VcenterUrl = types.StringNull()
    }
    if obj, ok := dataMap["vcenterUsername"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VcenterUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.VcenterUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.VcenterUsername = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.VcenterUsername = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.VcenterUsername = types.StringValue(string(jsonBytes))
        } else {
            data.VcenterUsername = types.StringNull()
        }
    } else if val, ok := dataMap["vcenterUsername"].(string); ok {
        data.VcenterUsername = types.StringValue(val)
    } else {
        data.VcenterUsername = types.StringNull()
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
    } else if val, ok := dataMap["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := dataMap["trustedCertificateFingerprint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrustedCertificateFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprint = types.StringNull()
        }
    } else if val, ok := dataMap["trustedCertificateFingerprint"].(string); ok {
        data.TrustedCertificateFingerprint = types.StringValue(val)
    } else {
        data.TrustedCertificateFingerprint = types.StringNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, dataMap["collectionIntervalInMinutes"])
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
    data.DatacenterCount = numberFromAPI(data.DatacenterCount, dataMap["datacenterCount"])
    data.ClusterCount = numberFromAPI(data.ClusterCount, dataMap["clusterCount"])
    data.HostCount = numberFromAPI(data.HostCount, dataMap["hostCount"])
    data.VmCount = numberFromAPI(data.VmCount, dataMap["vmCount"])
    data.PoweredOnVmCount = numberFromAPI(data.PoweredOnVmCount, dataMap["poweredOnVmCount"])
    data.DatastoreCount = numberFromAPI(data.DatastoreCount, dataMap["datastoreCount"])
    data.ResourcePoolCount = numberFromAPI(data.ResourcePoolCount, dataMap["resourcePoolCount"])
    data.DatastoreCapacityBytes = numberFromAPI(data.DatastoreCapacityBytes, dataMap["datastoreCapacityBytes"])
    data.DatastoreUsedBytes = numberFromAPI(data.DatastoreUsedBytes, dataMap["datastoreUsedBytes"])
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
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
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
    if val, ok := dataMap["isVCenterPasswordSet"].(bool); ok {
        data.IsVCenterPasswordSet = types.BoolValue(val)
    }
    if obj, ok := dataMap["vcenterCredentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.VcenterCredentialsUpdatedAt = NewRFC3339Value(val)
        } else {
            data.VcenterCredentialsUpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["vcenterCredentialsUpdatedAt"].(string); ok && val != "" {
        data.VcenterCredentialsUpdatedAt = NewRFC3339Value(val)
    } else {
        data.VcenterCredentialsUpdatedAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VcenterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data VcenterResourceModel
    var state VcenterResourceModel

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

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config VcenterResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    vcenterRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := vcenterRequest["data"].(map[string]interface{})

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
    if !data.DatacenterCount.IsUnknown() && !state.DatacenterCount.IsUnknown() && !data.DatacenterCount.Equal(state.DatacenterCount) {
        requestDataMap["datacenterCount"] = r.bigFloatToFloat64(data.DatacenterCount.ValueBigFloat())
    }
    if !data.ClusterCount.IsUnknown() && !state.ClusterCount.IsUnknown() && !data.ClusterCount.Equal(state.ClusterCount) {
        requestDataMap["clusterCount"] = r.bigFloatToFloat64(data.ClusterCount.ValueBigFloat())
    }
    if !data.HostCount.IsUnknown() && !state.HostCount.IsUnknown() && !data.HostCount.Equal(state.HostCount) {
        requestDataMap["hostCount"] = r.bigFloatToFloat64(data.HostCount.ValueBigFloat())
    }
    if !data.VmCount.IsUnknown() && !state.VmCount.IsUnknown() && !data.VmCount.Equal(state.VmCount) {
        requestDataMap["vmCount"] = r.bigFloatToFloat64(data.VmCount.ValueBigFloat())
    }
    if !data.PoweredOnVmCount.IsUnknown() && !state.PoweredOnVmCount.IsUnknown() && !data.PoweredOnVmCount.Equal(state.PoweredOnVmCount) {
        requestDataMap["poweredOnVmCount"] = r.bigFloatToFloat64(data.PoweredOnVmCount.ValueBigFloat())
    }
    if !data.DatastoreCount.IsUnknown() && !state.DatastoreCount.IsUnknown() && !data.DatastoreCount.Equal(state.DatastoreCount) {
        requestDataMap["datastoreCount"] = r.bigFloatToFloat64(data.DatastoreCount.ValueBigFloat())
    }
    if !data.ResourcePoolCount.IsUnknown() && !state.ResourcePoolCount.IsUnknown() && !data.ResourcePoolCount.Equal(state.ResourcePoolCount) {
        requestDataMap["resourcePoolCount"] = r.bigFloatToFloat64(data.ResourcePoolCount.ValueBigFloat())
    }
    if !data.DatastoreCapacityBytes.IsUnknown() && !state.DatastoreCapacityBytes.IsUnknown() && !data.DatastoreCapacityBytes.Equal(state.DatastoreCapacityBytes) {
        requestDataMap["datastoreCapacityBytes"] = r.bigFloatToFloat64(data.DatastoreCapacityBytes.ValueBigFloat())
    }
    if !data.DatastoreUsedBytes.IsUnknown() && !state.DatastoreUsedBytes.IsUnknown() && !data.DatastoreUsedBytes.Equal(state.DatastoreUsedBytes) {
        requestDataMap["datastoreUsedBytes"] = r.bigFloatToFloat64(data.DatastoreUsedBytes.ValueBigFloat())
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
    if !data.CollectionMethod.IsUnknown() && !state.CollectionMethod.IsUnknown() && !data.CollectionMethod.Equal(state.CollectionMethod) {
        requestDataMap["collectionMethod"] = data.CollectionMethod.ValueString()
    }
    if !data.VcenterUrl.IsUnknown() && !state.VcenterUrl.IsUnknown() && !data.VcenterUrl.Equal(state.VcenterUrl) {
        requestDataMap["vcenterUrl"] = data.VcenterUrl.ValueString()
    }
    if !data.VcenterUsername.IsUnknown() && !state.VcenterUsername.IsUnknown() && !data.VcenterUsername.Equal(state.VcenterUsername) {
        requestDataMap["vcenterUsername"] = data.VcenterUsername.ValueString()
    }
    if !data.VcenterPassword.IsUnknown() && !state.VcenterPassword.IsUnknown() && !data.VcenterPassword.Equal(state.VcenterPassword) {
        requestDataMap["vcenterPassword"] = data.VcenterPassword.ValueString()
    }
    if !data.CollectionProbeId.IsUnknown() && !state.CollectionProbeId.IsUnknown() && !data.CollectionProbeId.Equal(state.CollectionProbeId) {
        requestDataMap["collectionProbeId"] = data.CollectionProbeId.ValueString()
    }
    if !data.TrustedCertificateFingerprint.IsUnknown() && !state.TrustedCertificateFingerprint.IsUnknown() && !data.TrustedCertificateFingerprint.Equal(state.TrustedCertificateFingerprint) {
        requestDataMap["trustedCertificateFingerprint"] = data.TrustedCertificateFingerprint.ValueString()
    }
    if !data.CollectionIntervalInMinutes.IsUnknown() && !state.CollectionIntervalInMinutes.IsUnknown() && !data.CollectionIntervalInMinutes.Equal(state.CollectionIntervalInMinutes) {
        requestDataMap["collectionIntervalInMinutes"] = r.bigFloatToFloat64(data.CollectionIntervalInMinutes.ValueBigFloat())
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(vcenterRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/vmware-vcenter/" + data.Id.ValueString() + "", vcenterRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update vcenter, got error: %s", err))
            return
        }

        // Parse the update response
        var vcenterResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &vcenterResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update vcenter: %s", err))
            return
        }
        _ = vcenterResponse
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
        "vcenterUrl": true,
        "vcenterUsername": true,
        "collectionProbeId": true,
        "trustedCertificateFingerprint": true,
        "collectionIntervalInMinutes": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "datacenterCount": true,
        "clusterCount": true,
        "hostCount": true,
        "vmCount": true,
        "poweredOnVmCount": true,
        "datastoreCount": true,
        "resourcePoolCount": true,
        "datastoreCapacityBytes": true,
        "datastoreUsedBytes": true,
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "isVCenterPasswordSet": true,
        "vcenterCredentialsUpdatedAt": true,
        "collectionStatus": true,
        "collectionErrorCode": true,
        "collectionError": true,
        "presentedCertificate": true,
        "collectionSummary": true,
        "nextCollectionAt": true,
        "lastCollectionAt": true,
        "lastSuccessfulCollectionAt": true,
        "collectionSettingsVersion": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/vmware-vcenter/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vcenter after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read vcenter after update: %s", err))
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
    } else if val, ok := dataMap["collectionMethod"].(string); ok {
        data.CollectionMethod = types.StringValue(val)
    } else {
        data.CollectionMethod = types.StringNull()
    }
    if obj, ok := dataMap["vcenterUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VcenterUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.VcenterUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.VcenterUrl = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.VcenterUrl = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.VcenterUrl = types.StringValue(string(jsonBytes))
        } else {
            data.VcenterUrl = types.StringNull()
        }
    } else if val, ok := dataMap["vcenterUrl"].(string); ok {
        data.VcenterUrl = types.StringValue(val)
    } else {
        data.VcenterUrl = types.StringNull()
    }
    if obj, ok := dataMap["vcenterUsername"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VcenterUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.VcenterUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.VcenterUsername = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.VcenterUsername = types.StringValue(string(jsonBytes))
            } else {
                data.VcenterUsername = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.VcenterUsername = types.StringValue(string(jsonBytes))
        } else {
            data.VcenterUsername = types.StringNull()
        }
    } else if val, ok := dataMap["vcenterUsername"].(string); ok {
        data.VcenterUsername = types.StringValue(val)
    } else {
        data.VcenterUsername = types.StringNull()
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
    } else if val, ok := dataMap["collectionProbeId"].(string); ok {
        data.CollectionProbeId = types.StringValue(val)
    } else {
        data.CollectionProbeId = types.StringNull()
    }
    if obj, ok := dataMap["trustedCertificateFingerprint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrustedCertificateFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrustedCertificateFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.TrustedCertificateFingerprint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrustedCertificateFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.TrustedCertificateFingerprint = types.StringNull()
        }
    } else if val, ok := dataMap["trustedCertificateFingerprint"].(string); ok {
        data.TrustedCertificateFingerprint = types.StringValue(val)
    } else {
        data.TrustedCertificateFingerprint = types.StringNull()
    }
    data.CollectionIntervalInMinutes = numberFromAPI(data.CollectionIntervalInMinutes, dataMap["collectionIntervalInMinutes"])
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
    data.DatacenterCount = numberFromAPI(data.DatacenterCount, dataMap["datacenterCount"])
    data.ClusterCount = numberFromAPI(data.ClusterCount, dataMap["clusterCount"])
    data.HostCount = numberFromAPI(data.HostCount, dataMap["hostCount"])
    data.VmCount = numberFromAPI(data.VmCount, dataMap["vmCount"])
    data.PoweredOnVmCount = numberFromAPI(data.PoweredOnVmCount, dataMap["poweredOnVmCount"])
    data.DatastoreCount = numberFromAPI(data.DatastoreCount, dataMap["datastoreCount"])
    data.ResourcePoolCount = numberFromAPI(data.ResourcePoolCount, dataMap["resourcePoolCount"])
    data.DatastoreCapacityBytes = numberFromAPI(data.DatastoreCapacityBytes, dataMap["datastoreCapacityBytes"])
    data.DatastoreUsedBytes = numberFromAPI(data.DatastoreUsedBytes, dataMap["datastoreUsedBytes"])
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
    } else if val, ok := dataMap["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
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
    if val, ok := dataMap["isVCenterPasswordSet"].(bool); ok {
        data.IsVCenterPasswordSet = types.BoolValue(val)
    }
    if obj, ok := dataMap["vcenterCredentialsUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.VcenterCredentialsUpdatedAt = NewRFC3339Value(val)
        } else {
            data.VcenterCredentialsUpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["vcenterCredentialsUpdatedAt"].(string); ok && val != "" {
        data.VcenterCredentialsUpdatedAt = NewRFC3339Value(val)
    } else {
        data.VcenterCredentialsUpdatedAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VcenterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data VcenterResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/vmware-vcenter/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete vcenter, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete vcenter: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *VcenterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *VcenterResource) keepPlannedValues(data *VcenterResourceModel, plan *VcenterResourceModel, config *VcenterResourceModel) {
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
    if config.VcenterUrl.IsNull() && !plan.VcenterUrl.IsUnknown() {
        data.VcenterUrl = plan.VcenterUrl
    }
    if config.VcenterUsername.IsNull() && !plan.VcenterUsername.IsUnknown() {
        data.VcenterUsername = plan.VcenterUsername
    }
    if config.CollectionProbeId.IsNull() && !plan.CollectionProbeId.IsUnknown() {
        data.CollectionProbeId = plan.CollectionProbeId
    }
    if config.TrustedCertificateFingerprint.IsNull() && !plan.TrustedCertificateFingerprint.IsUnknown() {
        data.TrustedCertificateFingerprint = plan.TrustedCertificateFingerprint
    }
    if config.CollectionIntervalInMinutes.IsNull() && !plan.CollectionIntervalInMinutes.IsUnknown() {
        data.CollectionIntervalInMinutes = plan.CollectionIntervalInMinutes
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
    if config.DatacenterCount.IsNull() && !plan.DatacenterCount.IsUnknown() {
        data.DatacenterCount = plan.DatacenterCount
    }
    if config.ClusterCount.IsNull() && !plan.ClusterCount.IsUnknown() {
        data.ClusterCount = plan.ClusterCount
    }
    if config.HostCount.IsNull() && !plan.HostCount.IsUnknown() {
        data.HostCount = plan.HostCount
    }
    if config.VmCount.IsNull() && !plan.VmCount.IsUnknown() {
        data.VmCount = plan.VmCount
    }
    if config.PoweredOnVmCount.IsNull() && !plan.PoweredOnVmCount.IsUnknown() {
        data.PoweredOnVmCount = plan.PoweredOnVmCount
    }
    if config.DatastoreCount.IsNull() && !plan.DatastoreCount.IsUnknown() {
        data.DatastoreCount = plan.DatastoreCount
    }
    if config.ResourcePoolCount.IsNull() && !plan.ResourcePoolCount.IsUnknown() {
        data.ResourcePoolCount = plan.ResourcePoolCount
    }
    if config.DatastoreCapacityBytes.IsNull() && !plan.DatastoreCapacityBytes.IsUnknown() {
        data.DatastoreCapacityBytes = plan.DatastoreCapacityBytes
    }
    if config.DatastoreUsedBytes.IsNull() && !plan.DatastoreUsedBytes.IsUnknown() {
        data.DatastoreUsedBytes = plan.DatastoreUsedBytes
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
func (r *VcenterResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *VcenterResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *VcenterResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *VcenterResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *VcenterResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *VcenterResource) normalizeURLString(value string) string {
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
func (r *VcenterResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *VcenterResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
