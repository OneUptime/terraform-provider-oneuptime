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
var _ datasource.DataSource = &HostDataSource{}

func NewHostDataSource() datasource.DataSource {
    return &HostDataSource{}
}

// HostDataSource defines the data source implementation.
type HostDataSource struct {
    client *Client
}

// HostDataSourceModel describes the data source data model.
type HostDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    HostIdentifier types.String `tfsdk:"host_identifier"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    OsType types.String `tfsdk:"os_type"`
    OsVersion types.String `tfsdk:"os_version"`
    HostId types.String `tfsdk:"host_id"`
    HostArch types.String `tfsdk:"host_arch"`
    HostType types.String `tfsdk:"host_type"`
    HostIpAddresses types.String `tfsdk:"host_ip_addresses"`
    CpuCores types.Number `tfsdk:"cpu_cores"`
    TotalMemoryBytes types.Number `tfsdk:"total_memory_bytes"`
    ProcessCount types.Number `tfsdk:"process_count"`
    ContainerRuntime types.String `tfsdk:"container_runtime"`
    DockerHostId types.String `tfsdk:"docker_host_id"`
    KubernetesClusterId types.String `tfsdk:"kubernetes_cluster_id"`
    ProxmoxClusterId types.String `tfsdk:"proxmox_cluster_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    DeploymentEnvironment types.String `tfsdk:"deployment_environment"`
    RuntimeName types.String `tfsdk:"runtime_name"`
    RuntimeVersion types.String `tfsdk:"runtime_version"`
    CloudProvider types.String `tfsdk:"cloud_provider"`
    CloudPlatform types.String `tfsdk:"cloud_platform"`
    CloudRegion types.String `tfsdk:"cloud_region"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiCommandAllowlist types.String `tfsdk:"ai_command_allowlist"`
    AiAccessLastVerifiedAt types.String `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt types.String `tfsdk:"ai_access_configured_at"`
}

func (d *HostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *HostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered. Look up an existing host by `id`, or by any of its other arguments (`name`, `agent_version`, `ai_access_last_error`, ...): each one set must match, and exactly one host may match them all.",

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
                MarkdownDescription: "Friendly name for this host.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this host.",
                Optional: true,
                Computed: true,
            },
            "host_identifier": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for this host, sourced from the host.name OTel resource attribute.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector reporting on this host (connected or disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime agent reporting telemetry on this host, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When telemetry was last received from this host.",
                Computed: true,
            },
            "os_type": schema.StringAttribute{
                MarkdownDescription: "Operating system type of the host.",
                Optional: true,
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Operating system version of the host.",
                Optional: true,
                Computed: true,
            },
            "host_id": schema.StringAttribute{
                MarkdownDescription: "Stable host identifier reported by the OTel host.id resource attribute.",
                Optional: true,
                Computed: true,
            },
            "host_arch": schema.StringAttribute{
                MarkdownDescription: "CPU architecture from the OTel host.arch resource attribute.",
                Optional: true,
                Computed: true,
            },
            "host_type": schema.StringAttribute{
                MarkdownDescription: "Cloud-instance class reported by the OTel host.type resource attribute.",
                Optional: true,
                Computed: true,
            },
            "host_ip_addresses": schema.StringAttribute{
                MarkdownDescription: "Comma-separated list of every IP address reported by the OTel host.ip resource attribute, in the order the collector reported them, deduplicated. The Hosts list shows the most routable one (IPv4, non-loopback, non-link-local) first; the host detail page groups them all by category.",
                Optional: true,
                Computed: true,
            },
            "cpu_cores": schema.NumberAttribute{
                MarkdownDescription: "Logical CPU core count, sourced from system.cpu.logical.count metric.",
                Optional: true,
                Computed: true,
            },
            "total_memory_bytes": schema.NumberAttribute{
                MarkdownDescription: "Total physical memory in bytes, sourced from system.memory.usage metric (sum of all states).",
                Optional: true,
                Computed: true,
            },
            "process_count": schema.NumberAttribute{
                MarkdownDescription: "Most recent process count from system.processes.count metric.",
                Optional: true,
                Computed: true,
            },
            "container_runtime": schema.StringAttribute{
                MarkdownDescription: "Container runtime detected on this host, if any (e.g. docker, containerd).",
                Optional: true,
                Computed: true,
            },
            "docker_host_id": schema.StringAttribute{
                MarkdownDescription: "Optional FK to the DockerHost record for this same host when it is also running the Docker runtime. The ID of a `oneuptime_docker_host`.",
                Optional: true,
                Computed: true,
            },
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "Optional FK to the KubernetesCluster this host belongs to when k8s.cluster.name is reported. The ID of a `oneuptime_kubernetes_cluster`.",
                Optional: true,
                Computed: true,
            },
            "proxmox_cluster_id": schema.StringAttribute{
                MarkdownDescription: "Optional FK to the ProxmoxCluster this host runs inside (as a guest VM). The ID of a `oneuptime_proxmox_cluster`.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this host archived? Archived hosts are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this host archived?",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this host. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the host default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "deployment_environment": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging.",
                Optional: true,
                Computed: true,
            },
            "runtime_name": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the process.runtime.name OpenTelemetry resource attribute.",
                Optional: true,
                Computed: true,
            },
            "runtime_version": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the process.runtime.version OpenTelemetry resource attribute.",
                Optional: true,
                Computed: true,
            },
            "cloud_provider": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.",
                Optional: true,
                Computed: true,
            },
            "cloud_platform": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ec2, gcp_compute_engine.",
                Optional: true,
                Computed: true,
            },
            "cloud_region": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1.",
                Optional: true,
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.account.id OpenTelemetry resource attribute.",
                Optional: true,
                Computed: true,
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (systemctl status, journalctl, df, free, uptime, ps, ss) on this host, through its Host AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the host can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.",
                Optional: true,
                Computed: true,
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of command patterns that Automatic mode may run on this host without approval even though they are riskier changes. Each pattern is one command line for this host's agent (systemctl, journalctl and the other host programs it may run) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the host can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this host through its Host AI agent. Set by the server.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this host, kept until the next successful command. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this host was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Host AI agent that registers later never overwrites a setting an operator chose.",
                Computed: true,
            },
        },
    }
}

func (d *HostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HostDataSourceModel

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
    if !data.HostIdentifier.IsNull() && !data.HostIdentifier.IsUnknown() {
        filters["hostIdentifier"] = data.HostIdentifier.ValueString()
        filterNames = append(filterNames, "host_identifier = "+fmt.Sprintf("%q", data.HostIdentifier.ValueString()))
    }
    if !data.OtelCollectorStatus.IsNull() && !data.OtelCollectorStatus.IsUnknown() {
        filters["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
        filterNames = append(filterNames, "otel_collector_status = "+fmt.Sprintf("%q", data.OtelCollectorStatus.ValueString()))
    }
    if !data.AgentVersion.IsNull() && !data.AgentVersion.IsUnknown() {
        filters["agentVersion"] = data.AgentVersion.ValueString()
        filterNames = append(filterNames, "agent_version = "+fmt.Sprintf("%q", data.AgentVersion.ValueString()))
    }
    if !data.OsType.IsNull() && !data.OsType.IsUnknown() {
        filters["osType"] = data.OsType.ValueString()
        filterNames = append(filterNames, "os_type = "+fmt.Sprintf("%q", data.OsType.ValueString()))
    }
    if !data.OsVersion.IsNull() && !data.OsVersion.IsUnknown() {
        filters["osVersion"] = data.OsVersion.ValueString()
        filterNames = append(filterNames, "os_version = "+fmt.Sprintf("%q", data.OsVersion.ValueString()))
    }
    if !data.HostId.IsNull() && !data.HostId.IsUnknown() {
        filters["hostId"] = data.HostId.ValueString()
        filterNames = append(filterNames, "host_id = "+fmt.Sprintf("%q", data.HostId.ValueString()))
    }
    if !data.HostArch.IsNull() && !data.HostArch.IsUnknown() {
        filters["hostArch"] = data.HostArch.ValueString()
        filterNames = append(filterNames, "host_arch = "+fmt.Sprintf("%q", data.HostArch.ValueString()))
    }
    if !data.HostType.IsNull() && !data.HostType.IsUnknown() {
        filters["hostType"] = data.HostType.ValueString()
        filterNames = append(filterNames, "host_type = "+fmt.Sprintf("%q", data.HostType.ValueString()))
    }
    if !data.HostIpAddresses.IsNull() && !data.HostIpAddresses.IsUnknown() {
        filters["hostIpAddresses"] = data.HostIpAddresses.ValueString()
        filterNames = append(filterNames, "host_ip_addresses = "+fmt.Sprintf("%q", data.HostIpAddresses.ValueString()))
    }
    if !data.CpuCores.IsNull() && !data.CpuCores.IsUnknown() {
        filters["cpuCores"] = lookupNumber(data.CpuCores)
        filterNames = append(filterNames, "cpu_cores = "+data.CpuCores.ValueBigFloat().String())
    }
    if !data.TotalMemoryBytes.IsNull() && !data.TotalMemoryBytes.IsUnknown() {
        filters["totalMemoryBytes"] = lookupNumber(data.TotalMemoryBytes)
        filterNames = append(filterNames, "total_memory_bytes = "+data.TotalMemoryBytes.ValueBigFloat().String())
    }
    if !data.ProcessCount.IsNull() && !data.ProcessCount.IsUnknown() {
        filters["processCount"] = lookupNumber(data.ProcessCount)
        filterNames = append(filterNames, "process_count = "+data.ProcessCount.ValueBigFloat().String())
    }
    if !data.ContainerRuntime.IsNull() && !data.ContainerRuntime.IsUnknown() {
        filters["containerRuntime"] = data.ContainerRuntime.ValueString()
        filterNames = append(filterNames, "container_runtime = "+fmt.Sprintf("%q", data.ContainerRuntime.ValueString()))
    }
    if !data.DockerHostId.IsNull() && !data.DockerHostId.IsUnknown() {
        filters["dockerHostId"] = data.DockerHostId.ValueString()
        filterNames = append(filterNames, "docker_host_id = "+fmt.Sprintf("%q", data.DockerHostId.ValueString()))
    }
    if !data.KubernetesClusterId.IsNull() && !data.KubernetesClusterId.IsUnknown() {
        filters["kubernetesClusterId"] = data.KubernetesClusterId.ValueString()
        filterNames = append(filterNames, "kubernetes_cluster_id = "+fmt.Sprintf("%q", data.KubernetesClusterId.ValueString()))
    }
    if !data.ProxmoxClusterId.IsNull() && !data.ProxmoxClusterId.IsUnknown() {
        filters["proxmoxClusterId"] = data.ProxmoxClusterId.ValueString()
        filterNames = append(filterNames, "proxmox_cluster_id = "+fmt.Sprintf("%q", data.ProxmoxClusterId.ValueString()))
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
    if !data.DeploymentEnvironment.IsNull() && !data.DeploymentEnvironment.IsUnknown() {
        filters["deploymentEnvironment"] = data.DeploymentEnvironment.ValueString()
        filterNames = append(filterNames, "deployment_environment = "+fmt.Sprintf("%q", data.DeploymentEnvironment.ValueString()))
    }
    if !data.RuntimeName.IsNull() && !data.RuntimeName.IsUnknown() {
        filters["runtimeName"] = data.RuntimeName.ValueString()
        filterNames = append(filterNames, "runtime_name = "+fmt.Sprintf("%q", data.RuntimeName.ValueString()))
    }
    if !data.RuntimeVersion.IsNull() && !data.RuntimeVersion.IsUnknown() {
        filters["runtimeVersion"] = data.RuntimeVersion.ValueString()
        filterNames = append(filterNames, "runtime_version = "+fmt.Sprintf("%q", data.RuntimeVersion.ValueString()))
    }
    if !data.CloudProvider.IsNull() && !data.CloudProvider.IsUnknown() {
        filters["cloudProvider"] = data.CloudProvider.ValueString()
        filterNames = append(filterNames, "cloud_provider = "+fmt.Sprintf("%q", data.CloudProvider.ValueString()))
    }
    if !data.CloudPlatform.IsNull() && !data.CloudPlatform.IsUnknown() {
        filters["cloudPlatform"] = data.CloudPlatform.ValueString()
        filterNames = append(filterNames, "cloud_platform = "+fmt.Sprintf("%q", data.CloudPlatform.ValueString()))
    }
    if !data.CloudRegion.IsNull() && !data.CloudRegion.IsUnknown() {
        filters["cloudRegion"] = data.CloudRegion.ValueString()
        filterNames = append(filterNames, "cloud_region = "+fmt.Sprintf("%q", data.CloudRegion.ValueString()))
    }
    if !data.CloudAccountId.IsNull() && !data.CloudAccountId.IsUnknown() {
        filters["cloudAccountId"] = data.CloudAccountId.ValueString()
        filterNames = append(filterNames, "cloud_account_id = "+fmt.Sprintf("%q", data.CloudAccountId.ValueString()))
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
            "Look the host up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the host up by.",
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
        "hostIdentifier": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "osType": true,
        "osVersion": true,
        "hostId": true,
        "hostArch": true,
        "hostType": true,
        "hostIpAddresses": true,
        "cpuCores": true,
        "totalMemoryBytes": true,
        "processCount": true,
        "containerRuntime": true,
        "dockerHostId": true,
        "kubernetesClusterId": true,
        "proxmoxClusterId": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "deploymentEnvironment": true,
        "runtimeName": true,
        "runtimeVersion": true,
        "cloudProvider": true,
        "cloudPlatform": true,
        "cloudRegion": true,
        "cloudAccountId": true,
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/host/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No host found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read host: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/host/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list host, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list host: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No host matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one host matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for host.")
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
    if obj, ok := item["hostIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.HostIdentifier = types.StringNull()
        }
    } else if val, ok := item["hostIdentifier"].(string); ok {
        data.HostIdentifier = types.StringValue(val)
    } else {
        data.HostIdentifier = types.StringNull()
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
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    if obj, ok := item["osType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OsType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OsType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OsType = types.StringValue(string(jsonBytes))
        } else {
            data.OsType = types.StringNull()
        }
    } else if val, ok := item["osType"].(string); ok {
        data.OsType = types.StringValue(val)
    } else {
        data.OsType = types.StringNull()
    }
    if obj, ok := item["osVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OsVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OsVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OsVersion = types.StringValue(string(jsonBytes))
        } else {
            data.OsVersion = types.StringNull()
        }
    } else if val, ok := item["osVersion"].(string); ok {
        data.OsVersion = types.StringValue(val)
    } else {
        data.OsVersion = types.StringNull()
    }
    if obj, ok := item["hostId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostId = types.StringValue(string(jsonBytes))
        } else {
            data.HostId = types.StringNull()
        }
    } else if val, ok := item["hostId"].(string); ok {
        data.HostId = types.StringValue(val)
    } else {
        data.HostId = types.StringNull()
    }
    if obj, ok := item["hostArch"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostArch = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostArch = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostArch = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostArch = types.StringValue(string(jsonBytes))
        } else {
            data.HostArch = types.StringNull()
        }
    } else if val, ok := item["hostArch"].(string); ok {
        data.HostArch = types.StringValue(val)
    } else {
        data.HostArch = types.StringNull()
    }
    if obj, ok := item["hostType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostType = types.StringValue(string(jsonBytes))
        } else {
            data.HostType = types.StringNull()
        }
    } else if val, ok := item["hostType"].(string); ok {
        data.HostType = types.StringValue(val)
    } else {
        data.HostType = types.StringNull()
    }
    if obj, ok := item["hostIpAddresses"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostIpAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostIpAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostIpAddresses = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostIpAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.HostIpAddresses = types.StringNull()
        }
    } else if val, ok := item["hostIpAddresses"].(string); ok {
        data.HostIpAddresses = types.StringValue(val)
    } else {
        data.HostIpAddresses = types.StringNull()
    }
    if val, ok := item["cpuCores"].(float64); ok {
        data.CpuCores = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["cpuCores"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CpuCores = types.NumberValue(big.NewFloat(val))
        } else {
            data.CpuCores = types.NumberNull()
        }
    } else {
        data.CpuCores = types.NumberNull()
    }
    if val, ok := item["totalMemoryBytes"].(float64); ok {
        data.TotalMemoryBytes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["totalMemoryBytes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TotalMemoryBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.TotalMemoryBytes = types.NumberNull()
        }
    } else {
        data.TotalMemoryBytes = types.NumberNull()
    }
    if val, ok := item["processCount"].(float64); ok {
        data.ProcessCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["processCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ProcessCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ProcessCount = types.NumberNull()
        }
    } else {
        data.ProcessCount = types.NumberNull()
    }
    if obj, ok := item["containerRuntime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerRuntime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContainerRuntime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContainerRuntime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContainerRuntime = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerRuntime = types.StringNull()
        }
    } else if val, ok := item["containerRuntime"].(string); ok {
        data.ContainerRuntime = types.StringValue(val)
    } else {
        data.ContainerRuntime = types.StringNull()
    }
    if obj, ok := item["dockerHostId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DockerHostId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DockerHostId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DockerHostId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DockerHostId = types.StringValue(string(jsonBytes))
        } else {
            data.DockerHostId = types.StringNull()
        }
    } else if val, ok := item["dockerHostId"].(string); ok {
        data.DockerHostId = types.StringValue(val)
    } else {
        data.DockerHostId = types.StringNull()
    }
    if obj, ok := item["kubernetesClusterId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KubernetesClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesClusterId = types.StringNull()
        }
    } else if val, ok := item["kubernetesClusterId"].(string); ok {
        data.KubernetesClusterId = types.StringValue(val)
    } else {
        data.KubernetesClusterId = types.StringNull()
    }
    if obj, ok := item["proxmoxClusterId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProxmoxClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProxmoxClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProxmoxClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProxmoxClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.ProxmoxClusterId = types.StringNull()
        }
    } else if val, ok := item["proxmoxClusterId"].(string); ok {
        data.ProxmoxClusterId = types.StringValue(val)
    } else {
        data.ProxmoxClusterId = types.StringNull()
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
    if val, ok := item["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
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
    } else if val, ok := item["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    } else {
        data.TelemetryRetentionConfig = types.StringNull()
    }
    if obj, ok := item["deploymentEnvironment"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeploymentEnvironment = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeploymentEnvironment = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeploymentEnvironment = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeploymentEnvironment = types.StringValue(string(jsonBytes))
        } else {
            data.DeploymentEnvironment = types.StringNull()
        }
    } else if val, ok := item["deploymentEnvironment"].(string); ok {
        data.DeploymentEnvironment = types.StringValue(val)
    } else {
        data.DeploymentEnvironment = types.StringNull()
    }
    if obj, ok := item["runtimeName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuntimeName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuntimeName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuntimeName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuntimeName = types.StringValue(string(jsonBytes))
        } else {
            data.RuntimeName = types.StringNull()
        }
    } else if val, ok := item["runtimeName"].(string); ok {
        data.RuntimeName = types.StringValue(val)
    } else {
        data.RuntimeName = types.StringNull()
    }
    if obj, ok := item["runtimeVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuntimeVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuntimeVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuntimeVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuntimeVersion = types.StringValue(string(jsonBytes))
        } else {
            data.RuntimeVersion = types.StringNull()
        }
    } else if val, ok := item["runtimeVersion"].(string); ok {
        data.RuntimeVersion = types.StringValue(val)
    } else {
        data.RuntimeVersion = types.StringNull()
    }
    if obj, ok := item["cloudProvider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudProvider = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudProvider = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudProvider = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudProvider = types.StringValue(string(jsonBytes))
        } else {
            data.CloudProvider = types.StringNull()
        }
    } else if val, ok := item["cloudProvider"].(string); ok {
        data.CloudProvider = types.StringValue(val)
    } else {
        data.CloudProvider = types.StringNull()
    }
    if obj, ok := item["cloudPlatform"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudPlatform = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudPlatform = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudPlatform = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudPlatform = types.StringValue(string(jsonBytes))
        } else {
            data.CloudPlatform = types.StringNull()
        }
    } else if val, ok := item["cloudPlatform"].(string); ok {
        data.CloudPlatform = types.StringValue(val)
    } else {
        data.CloudPlatform = types.StringNull()
    }
    if obj, ok := item["cloudRegion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudRegion = types.StringValue(string(jsonBytes))
        } else {
            data.CloudRegion = types.StringNull()
        }
    } else if val, ok := item["cloudRegion"].(string); ok {
        data.CloudRegion = types.StringValue(val)
    } else {
        data.CloudRegion = types.StringNull()
    }
    if obj, ok := item["cloudAccountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := item["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }
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
    } else if val, ok := item["aiAccessConfiguredAt"].(string); ok {
        data.AiAccessConfiguredAt = types.StringValue(val)
    } else {
        data.AiAccessConfiguredAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
