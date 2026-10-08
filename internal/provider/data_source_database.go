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
var _ datasource.DataSource = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
    return &DatabaseDataSource{}
}

// DatabaseDataSource defines the data source implementation.
type DatabaseDataSource struct {
    client *Client
}

// DatabaseDataSourceModel describes the data source data model.
type DatabaseDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    DatabaseIdentifier types.String `tfsdk:"database_identifier"`
    WorkloadIdentifier types.String `tfsdk:"workload_identifier"`
    DbSystem types.String `tfsdk:"db_system"`
    ServerAddress types.String `tfsdk:"server_address"`
    ServerPort types.Number `tfsdk:"server_port"`
    DbVersion types.String `tfsdk:"db_version"`
    DiscoverySource types.String `tfsdk:"discovery_source"`
    KubernetesClusterId types.String `tfsdk:"kubernetes_cluster_id"`
    KubernetesNamespace types.String `tfsdk:"kubernetes_namespace"`
    WorkloadKind types.String `tfsdk:"workload_kind"`
    WorkloadName types.String `tfsdk:"workload_name"`
    DockerHostId types.String `tfsdk:"docker_host_id"`
    PodmanHostId types.String `tfsdk:"podman_host_id"`
    MemberEntityKeys types.String `tfsdk:"member_entity_keys"`
    InstanceCount types.Number `tfsdk:"instance_count"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    CollectorLastSeenAt types.String `tfsdk:"collector_last_seen_at"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    AutoArchivedAt types.String `tfsdk:"auto_archived_at"`
    ManuallyRestoredAt types.String `tfsdk:"manually_restored_at"`
    DbSystemSource types.String `tfsdk:"db_system_source"`
    WorkloadLastSeenAt types.String `tfsdk:"workload_last_seen_at"`
    AutomaticAssignments types.String `tfsdk:"automatic_assignments"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiCommandAllowlist types.String `tfsdk:"ai_command_allowlist"`
    AiAccessLastVerifiedAt types.String `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt types.String `tfsdk:"ai_access_configured_at"`
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Database servers this project runs or talks to. Each database is discovered from OpenTelemetry Collector database receivers, from database calls in application traces, from database workloads on monitored Kubernetes clusters and Docker / Podman hosts, or added manually. Look up an existing database by `id`, or by any of its other arguments (`name`, `agent_version`, `ai_access_last_error`, ...): each one set must match, and exactly one database may match them all.",

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
                MarkdownDescription: "Name of this database. Discovered databases are named after their engine and endpoint, e.g. PostgreSQL orders-db.internal:5432. Not unique - rename freely.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this database.",
                Optional: true,
                Computed: true,
            },
            "database_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identity of this database within the project: the engine family and its canonical endpoint joined with '|' (e.g. postgresql|orders-db.internal:5432), or the workload form for databases discovered on Kubernetes / Docker / Podman (e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db). A fork keys like the engine it forks (a MariaDB database is mysql|..., a Valkey one redis|...), so learning which of the two a database is never changes its identity. Computed by the server for manually added databases.",
                Optional: true,
                Computed: true,
            },
            "workload_identifier": schema.StringAttribute{
                MarkdownDescription: "Identity of the Kubernetes / Docker / Podman workload this database runs as, e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db. The prefix is the engine family, not the engine (a MariaDB workload is mysql|..., a Valkey one redis|...). Empty for databases that are only known by their endpoint.",
                Optional: true,
                Computed: true,
            },
            "db_system": schema.StringAttribute{
                MarkdownDescription: "Database engine family as an OpenTelemetry db.system.name value, e.g. postgresql, mysql, redis, mongodb, microsoft.sql_server.",
                Optional: true,
                Computed: true,
            },
            "server_address": schema.StringAttribute{
                MarkdownDescription: "Host name or IP address of the primary endpoint of this database (the OpenTelemetry server.address attribute).",
                Optional: true,
                Computed: true,
            },
            "server_port": schema.NumberAttribute{
                MarkdownDescription: "Port of the primary endpoint of this database (the OpenTelemetry server.port attribute). Defaults to the engine's standard port when not given.",
                Optional: true,
                Computed: true,
            },
            "db_version": schema.StringAttribute{
                MarkdownDescription: "Engine version last reported for this database, from the collector receiver or the container image tag.",
                Optional: true,
                Computed: true,
            },
            "discovery_source": schema.StringAttribute{
                MarkdownDescription: "How this database was first discovered: collector, client-spans, kubernetes, docker, podman or manual.",
                Optional: true,
                Computed: true,
            },
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Kubernetes cluster this database runs on, when it was discovered as a workload on a monitored cluster. The ID of a `oneuptime_kubernetes_cluster`.",
                Optional: true,
                Computed: true,
            },
            "kubernetes_namespace": schema.StringAttribute{
                MarkdownDescription: "Kubernetes namespace of the workload this database runs as.",
                Optional: true,
                Computed: true,
            },
            "workload_kind": schema.StringAttribute{
                MarkdownDescription: "Kind of workload this database runs as, e.g. StatefulSet, Deployment, Cluster (an operator-managed cluster) or Container.",
                Optional: true,
                Computed: true,
            },
            "workload_name": schema.StringAttribute{
                MarkdownDescription: "Name of the workload (StatefulSet, Deployment, operator cluster or container group) this database runs as.",
                Optional: true,
                Computed: true,
            },
            "docker_host_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Docker host this database runs on, when it was discovered as a container on a monitored host. The ID of a `oneuptime_docker_host`.",
                Optional: true,
                Computed: true,
            },
            "podman_host_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Podman host this database runs on, when it was discovered as a container on a monitored host. The ID of a `oneuptime_podman_host`.",
                Optional: true,
                Computed: true,
            },
            "member_entity_keys": schema.StringAttribute{
                MarkdownDescription: "Telemetry entity keys of the pods / containers this database runs as, each with the time it was last seen. Maintained by discovery. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "instance_count": schema.NumberAttribute{
                MarkdownDescription: "Number of running instances (pods / containers) last seen for this database's workload.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Whether engine metrics are currently being received from a collector / agent for this database (connected) or a collector reported and has stopped (disconnected). Empty when no collector has ever reported.",
                Optional: true,
                Computed: true,
            },
            "collector_last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When engine telemetry from a collector / agent was last received for this database.",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime agent reporting this database, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When this database was last seen by any discovery source - engine telemetry, application database calls or the workload inventory.",
                Computed: true,
            },
            "auto_archived_at": schema.StringAttribute{
                MarkdownDescription: "When this database was archived automatically because no discovery source had seen it for a while. Empty when it was archived by a person or is not archived.",
                Computed: true,
            },
            "manually_restored_at": schema.StringAttribute{
                MarkdownDescription: "When a person last restored this database from the archive. Automatic archiving leaves it alone until it is seen again or a grace period passes.",
                Computed: true,
            },
            "db_system_source": schema.StringAttribute{
                MarkdownDescription: "What the database engine was last determined from: manual, collector, container or client-spans. Stronger evidence may correct the engine; weaker evidence never changes it.",
                Optional: true,
                Computed: true,
            },
            "workload_last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When Kubernetes, Docker or Podman discovery last saw the workload this database runs as. Empty for a database that is not a discovered workload.",
                Computed: true,
            },
            "automatic_assignments": schema.StringAttribute{
                MarkdownDescription: "Label and owner ids that label rules, owner rules or telemetry attached automatically. Maintained by OneUptime. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this database archived? Archived databases are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this database archived?",
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
                MarkdownDescription: "Number of days to retain telemetry collected from this database by a collector / agent. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for telemetry collected from this database by a collector / agent (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the database default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (db diagnostics from a fixed catalog such as ping, version and sessions; never free-form SQL) on this database server, through its Database AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the database server can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this database server. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the database server's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the database server can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.",
                Optional: true,
                Computed: true,
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of command patterns that Automatic mode may run on this database server without approval even though they are riskier changes. Each pattern is one command line for this database server's agent (db, a fixed catalog of diagnostics, never free-form SQL) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this database server is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the database server can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this database server through its Database AI agent. Set by the server.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this database server, kept until the next successful command. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this database server was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Database AI agent that registers later never overwrites a setting an operator chose.",
                Computed: true,
            },
        },
    }
}

func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DatabaseDataSourceModel

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
    if !data.DatabaseIdentifier.IsNull() && !data.DatabaseIdentifier.IsUnknown() {
        filters["databaseIdentifier"] = data.DatabaseIdentifier.ValueString()
        filterNames = append(filterNames, "database_identifier = "+fmt.Sprintf("%q", data.DatabaseIdentifier.ValueString()))
    }
    if !data.WorkloadIdentifier.IsNull() && !data.WorkloadIdentifier.IsUnknown() {
        filters["workloadIdentifier"] = data.WorkloadIdentifier.ValueString()
        filterNames = append(filterNames, "workload_identifier = "+fmt.Sprintf("%q", data.WorkloadIdentifier.ValueString()))
    }
    if !data.DbSystem.IsNull() && !data.DbSystem.IsUnknown() {
        filters["dbSystem"] = data.DbSystem.ValueString()
        filterNames = append(filterNames, "db_system = "+fmt.Sprintf("%q", data.DbSystem.ValueString()))
    }
    if !data.ServerAddress.IsNull() && !data.ServerAddress.IsUnknown() {
        filters["serverAddress"] = data.ServerAddress.ValueString()
        filterNames = append(filterNames, "server_address = "+fmt.Sprintf("%q", data.ServerAddress.ValueString()))
    }
    if !data.ServerPort.IsNull() && !data.ServerPort.IsUnknown() {
        filters["serverPort"] = lookupNumber(data.ServerPort)
        filterNames = append(filterNames, "server_port = "+data.ServerPort.ValueBigFloat().String())
    }
    if !data.DbVersion.IsNull() && !data.DbVersion.IsUnknown() {
        filters["dbVersion"] = data.DbVersion.ValueString()
        filterNames = append(filterNames, "db_version = "+fmt.Sprintf("%q", data.DbVersion.ValueString()))
    }
    if !data.DiscoverySource.IsNull() && !data.DiscoverySource.IsUnknown() {
        filters["discoverySource"] = data.DiscoverySource.ValueString()
        filterNames = append(filterNames, "discovery_source = "+fmt.Sprintf("%q", data.DiscoverySource.ValueString()))
    }
    if !data.KubernetesClusterId.IsNull() && !data.KubernetesClusterId.IsUnknown() {
        filters["kubernetesClusterId"] = data.KubernetesClusterId.ValueString()
        filterNames = append(filterNames, "kubernetes_cluster_id = "+fmt.Sprintf("%q", data.KubernetesClusterId.ValueString()))
    }
    if !data.KubernetesNamespace.IsNull() && !data.KubernetesNamespace.IsUnknown() {
        filters["kubernetesNamespace"] = data.KubernetesNamespace.ValueString()
        filterNames = append(filterNames, "kubernetes_namespace = "+fmt.Sprintf("%q", data.KubernetesNamespace.ValueString()))
    }
    if !data.WorkloadKind.IsNull() && !data.WorkloadKind.IsUnknown() {
        filters["workloadKind"] = data.WorkloadKind.ValueString()
        filterNames = append(filterNames, "workload_kind = "+fmt.Sprintf("%q", data.WorkloadKind.ValueString()))
    }
    if !data.WorkloadName.IsNull() && !data.WorkloadName.IsUnknown() {
        filters["workloadName"] = data.WorkloadName.ValueString()
        filterNames = append(filterNames, "workload_name = "+fmt.Sprintf("%q", data.WorkloadName.ValueString()))
    }
    if !data.DockerHostId.IsNull() && !data.DockerHostId.IsUnknown() {
        filters["dockerHostId"] = data.DockerHostId.ValueString()
        filterNames = append(filterNames, "docker_host_id = "+fmt.Sprintf("%q", data.DockerHostId.ValueString()))
    }
    if !data.PodmanHostId.IsNull() && !data.PodmanHostId.IsUnknown() {
        filters["podmanHostId"] = data.PodmanHostId.ValueString()
        filterNames = append(filterNames, "podman_host_id = "+fmt.Sprintf("%q", data.PodmanHostId.ValueString()))
    }
    if !data.InstanceCount.IsNull() && !data.InstanceCount.IsUnknown() {
        filters["instanceCount"] = lookupNumber(data.InstanceCount)
        filterNames = append(filterNames, "instance_count = "+data.InstanceCount.ValueBigFloat().String())
    }
    if !data.OtelCollectorStatus.IsNull() && !data.OtelCollectorStatus.IsUnknown() {
        filters["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
        filterNames = append(filterNames, "otel_collector_status = "+fmt.Sprintf("%q", data.OtelCollectorStatus.ValueString()))
    }
    if !data.AgentVersion.IsNull() && !data.AgentVersion.IsUnknown() {
        filters["agentVersion"] = data.AgentVersion.ValueString()
        filterNames = append(filterNames, "agent_version = "+fmt.Sprintf("%q", data.AgentVersion.ValueString()))
    }
    if !data.DbSystemSource.IsNull() && !data.DbSystemSource.IsUnknown() {
        filters["dbSystemSource"] = data.DbSystemSource.ValueString()
        filterNames = append(filterNames, "db_system_source = "+fmt.Sprintf("%q", data.DbSystemSource.ValueString()))
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
            "Look the database up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the database up by.",
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
        "databaseIdentifier": true,
        "workloadIdentifier": true,
        "dbSystem": true,
        "serverAddress": true,
        "serverPort": true,
        "dbVersion": true,
        "discoverySource": true,
        "kubernetesClusterId": true,
        "kubernetesNamespace": true,
        "workloadKind": true,
        "workloadName": true,
        "dockerHostId": true,
        "podmanHostId": true,
        "memberEntityKeys": true,
        "instanceCount": true,
        "otelCollectorStatus": true,
        "collectorLastSeenAt": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "autoArchivedAt": true,
        "manuallyRestoredAt": true,
        "dbSystemSource": true,
        "workloadLastSeenAt": true,
        "automaticAssignments": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
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
        readPath := "/database-server/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No database found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read database: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/database-server/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list database, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list database: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No database matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one database matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for database.")
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
    if obj, ok := item["databaseIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DatabaseIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DatabaseIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DatabaseIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseIdentifier = types.StringNull()
        }
    } else if val, ok := item["databaseIdentifier"].(string); ok {
        data.DatabaseIdentifier = types.StringValue(val)
    } else {
        data.DatabaseIdentifier = types.StringNull()
    }
    if obj, ok := item["workloadIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkloadIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkloadIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkloadIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkloadIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.WorkloadIdentifier = types.StringNull()
        }
    } else if val, ok := item["workloadIdentifier"].(string); ok {
        data.WorkloadIdentifier = types.StringValue(val)
    } else {
        data.WorkloadIdentifier = types.StringNull()
    }
    if obj, ok := item["dbSystem"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DbSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DbSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DbSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DbSystem = types.StringValue(string(jsonBytes))
        } else {
            data.DbSystem = types.StringNull()
        }
    } else if val, ok := item["dbSystem"].(string); ok {
        data.DbSystem = types.StringValue(val)
    } else {
        data.DbSystem = types.StringNull()
    }
    if obj, ok := item["serverAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServerAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServerAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServerAddress = types.StringValue(string(jsonBytes))
        } else {
            data.ServerAddress = types.StringNull()
        }
    } else if val, ok := item["serverAddress"].(string); ok {
        data.ServerAddress = types.StringValue(val)
    } else {
        data.ServerAddress = types.StringNull()
    }
    if val, ok := item["serverPort"].(float64); ok {
        data.ServerPort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["serverPort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ServerPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.ServerPort = types.NumberNull()
        }
    } else {
        data.ServerPort = types.NumberNull()
    }
    if obj, ok := item["dbVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DbVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DbVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DbVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DbVersion = types.StringValue(string(jsonBytes))
        } else {
            data.DbVersion = types.StringNull()
        }
    } else if val, ok := item["dbVersion"].(string); ok {
        data.DbVersion = types.StringValue(val)
    } else {
        data.DbVersion = types.StringNull()
    }
    if obj, ok := item["discoverySource"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoverySource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DiscoverySource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DiscoverySource = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DiscoverySource = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoverySource = types.StringNull()
        }
    } else if val, ok := item["discoverySource"].(string); ok {
        data.DiscoverySource = types.StringValue(val)
    } else {
        data.DiscoverySource = types.StringNull()
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
    if obj, ok := item["kubernetesNamespace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesNamespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KubernetesNamespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KubernetesNamespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KubernetesNamespace = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesNamespace = types.StringNull()
        }
    } else if val, ok := item["kubernetesNamespace"].(string); ok {
        data.KubernetesNamespace = types.StringValue(val)
    } else {
        data.KubernetesNamespace = types.StringNull()
    }
    if obj, ok := item["workloadKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkloadKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkloadKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkloadKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkloadKind = types.StringValue(string(jsonBytes))
        } else {
            data.WorkloadKind = types.StringNull()
        }
    } else if val, ok := item["workloadKind"].(string); ok {
        data.WorkloadKind = types.StringValue(val)
    } else {
        data.WorkloadKind = types.StringNull()
    }
    if obj, ok := item["workloadName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkloadName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkloadName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkloadName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkloadName = types.StringValue(string(jsonBytes))
        } else {
            data.WorkloadName = types.StringNull()
        }
    } else if val, ok := item["workloadName"].(string); ok {
        data.WorkloadName = types.StringValue(val)
    } else {
        data.WorkloadName = types.StringNull()
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
    if obj, ok := item["podmanHostId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PodmanHostId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PodmanHostId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PodmanHostId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PodmanHostId = types.StringValue(string(jsonBytes))
        } else {
            data.PodmanHostId = types.StringNull()
        }
    } else if val, ok := item["podmanHostId"].(string); ok {
        data.PodmanHostId = types.StringValue(val)
    } else {
        data.PodmanHostId = types.StringNull()
    }
    if obj, ok := item["memberEntityKeys"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MemberEntityKeys = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MemberEntityKeys = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MemberEntityKeys = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MemberEntityKeys = types.StringValue(string(jsonBytes))
        } else {
            data.MemberEntityKeys = types.StringNull()
        }
    } else if val, ok := item["memberEntityKeys"].(string); ok {
        data.MemberEntityKeys = types.StringValue(val)
    } else {
        data.MemberEntityKeys = types.StringNull()
    }
    if val, ok := item["instanceCount"].(float64); ok {
        data.InstanceCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["instanceCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InstanceCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.InstanceCount = types.NumberNull()
        }
    } else {
        data.InstanceCount = types.NumberNull()
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
    if obj, ok := item["collectorLastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectorLastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectorLastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectorLastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectorLastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.CollectorLastSeenAt = types.StringNull()
        }
    } else if val, ok := item["collectorLastSeenAt"].(string); ok {
        data.CollectorLastSeenAt = types.StringValue(val)
    } else {
        data.CollectorLastSeenAt = types.StringNull()
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
    if obj, ok := item["autoArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AutoArchivedAt = types.StringNull()
        }
    } else if val, ok := item["autoArchivedAt"].(string); ok {
        data.AutoArchivedAt = types.StringValue(val)
    } else {
        data.AutoArchivedAt = types.StringNull()
    }
    if obj, ok := item["manuallyRestoredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ManuallyRestoredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ManuallyRestoredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ManuallyRestoredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ManuallyRestoredAt = types.StringValue(string(jsonBytes))
        } else {
            data.ManuallyRestoredAt = types.StringNull()
        }
    } else if val, ok := item["manuallyRestoredAt"].(string); ok {
        data.ManuallyRestoredAt = types.StringValue(val)
    } else {
        data.ManuallyRestoredAt = types.StringNull()
    }
    if obj, ok := item["dbSystemSource"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DbSystemSource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DbSystemSource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DbSystemSource = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DbSystemSource = types.StringValue(string(jsonBytes))
        } else {
            data.DbSystemSource = types.StringNull()
        }
    } else if val, ok := item["dbSystemSource"].(string); ok {
        data.DbSystemSource = types.StringValue(val)
    } else {
        data.DbSystemSource = types.StringNull()
    }
    if obj, ok := item["workloadLastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkloadLastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkloadLastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkloadLastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkloadLastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.WorkloadLastSeenAt = types.StringNull()
        }
    } else if val, ok := item["workloadLastSeenAt"].(string); ok {
        data.WorkloadLastSeenAt = types.StringValue(val)
    } else {
        data.WorkloadLastSeenAt = types.StringNull()
    }
    if obj, ok := item["automaticAssignments"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutomaticAssignments = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutomaticAssignments = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutomaticAssignments = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutomaticAssignments = types.StringValue(string(jsonBytes))
        } else {
            data.AutomaticAssignments = types.StringNull()
        }
    } else if val, ok := item["automaticAssignments"].(string); ok {
        data.AutomaticAssignments = types.StringValue(val)
    } else {
        data.AutomaticAssignments = types.StringNull()
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
