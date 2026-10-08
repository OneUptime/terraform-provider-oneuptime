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
var _ datasource.DataSource = &DockerHostDataSource{}

func NewDockerHostDataSource() datasource.DataSource {
    return &DockerHostDataSource{}
}

// DockerHostDataSource defines the data source implementation.
type DockerHostDataSource struct {
    client *Client
}

// DockerHostDataSourceModel describes the data source data model.
type DockerHostDataSourceModel struct {
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
    ContainersRunning types.Number `tfsdk:"containers_running"`
    ContainersStopped types.Number `tfsdk:"containers_stopped"`
    ContainersPaused types.Number `tfsdk:"containers_paused"`
    OsType types.String `tfsdk:"os_type"`
    OsVersion types.String `tfsdk:"os_version"`
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

func (d *DockerHostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_docker_host"
}

func (d *DockerHostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Docker Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Docker Agent sends metrics, or can be manually registered. Look up an existing docker host by `id`, or by any of its other arguments (`name`, `agent_version`, `ai_access_last_error`, ...): each one set must match, and exactly one docker host may match them all.",

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
                MarkdownDescription: "Friendly name for this Docker host.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this Docker host.",
                Optional: true,
                Computed: true,
            },
            "host_identifier": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for this Docker host, sourced from the host.name OTel resource attribute.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Docker agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When metrics were last received from this host.",
                Computed: true,
            },
            "containers_running": schema.NumberAttribute{
                MarkdownDescription: "Cached count of running containers on this host.",
                Optional: true,
                Computed: true,
            },
            "containers_stopped": schema.NumberAttribute{
                MarkdownDescription: "Cached count of stopped containers on this host.",
                Optional: true,
                Computed: true,
            },
            "containers_paused": schema.NumberAttribute{
                MarkdownDescription: "Cached count of paused containers on this host.",
                Optional: true,
                Computed: true,
            },
            "os_type": schema.StringAttribute{
                MarkdownDescription: "Operating system type of the Docker host.",
                Optional: true,
                Computed: true,
            },
            "os_version": schema.StringAttribute{
                MarkdownDescription: "Operating system version of the Docker host.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this Docker host archived? Archived Docker hosts are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this Docker host archived?",
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
                MarkdownDescription: "Number of days to retain telemetry data for this Docker host. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this Docker host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Docker host default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (docker ps, inspect, logs, stats, events) on this Docker host, through its Docker AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Docker host can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this Docker host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Docker host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Docker host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.",
                Optional: true,
                Computed: true,
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of command patterns that Automatic mode may run on this Docker host without approval even though they are riskier changes. Each pattern is one command line for this Docker host's agent (docker) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Docker host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Docker host can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this Docker host through its Docker AI agent. Set by the server.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this Docker host, kept until the next successful command. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this Docker host was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Docker AI agent that registers later never overwrites a setting an operator chose.",
                Computed: true,
            },
        },
    }
}

func (d *DockerHostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DockerHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DockerHostDataSourceModel

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
    if !data.ContainersRunning.IsNull() && !data.ContainersRunning.IsUnknown() {
        filters["containersRunning"] = lookupNumber(data.ContainersRunning)
        filterNames = append(filterNames, "containers_running = "+data.ContainersRunning.ValueBigFloat().String())
    }
    if !data.ContainersStopped.IsNull() && !data.ContainersStopped.IsUnknown() {
        filters["containersStopped"] = lookupNumber(data.ContainersStopped)
        filterNames = append(filterNames, "containers_stopped = "+data.ContainersStopped.ValueBigFloat().String())
    }
    if !data.ContainersPaused.IsNull() && !data.ContainersPaused.IsUnknown() {
        filters["containersPaused"] = lookupNumber(data.ContainersPaused)
        filterNames = append(filterNames, "containers_paused = "+data.ContainersPaused.ValueBigFloat().String())
    }
    if !data.OsType.IsNull() && !data.OsType.IsUnknown() {
        filters["osType"] = data.OsType.ValueString()
        filterNames = append(filterNames, "os_type = "+fmt.Sprintf("%q", data.OsType.ValueString()))
    }
    if !data.OsVersion.IsNull() && !data.OsVersion.IsUnknown() {
        filters["osVersion"] = data.OsVersion.ValueString()
        filterNames = append(filterNames, "os_version = "+fmt.Sprintf("%q", data.OsVersion.ValueString()))
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
            "Look the docker host up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the docker host up by.",
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
        "containersRunning": true,
        "containersStopped": true,
        "containersPaused": true,
        "osType": true,
        "osVersion": true,
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
        readPath := "/docker-host/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read docker_host, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No docker host found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read docker_host: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/docker-host/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list docker_host, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list docker_host: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No docker host matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one docker host matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for docker_host.")
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
    if val, ok := item["containersRunning"].(float64); ok {
        data.ContainersRunning = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["containersRunning"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ContainersRunning = types.NumberValue(big.NewFloat(val))
        } else {
            data.ContainersRunning = types.NumberNull()
        }
    } else {
        data.ContainersRunning = types.NumberNull()
    }
    if val, ok := item["containersStopped"].(float64); ok {
        data.ContainersStopped = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["containersStopped"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ContainersStopped = types.NumberValue(big.NewFloat(val))
        } else {
            data.ContainersStopped = types.NumberNull()
        }
    } else {
        data.ContainersStopped = types.NumberNull()
    }
    if val, ok := item["containersPaused"].(float64); ok {
        data.ContainersPaused = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["containersPaused"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ContainersPaused = types.NumberValue(big.NewFloat(val))
        } else {
            data.ContainersPaused = types.NumberNull()
        }
    } else {
        data.ContainersPaused = types.NumberNull()
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
