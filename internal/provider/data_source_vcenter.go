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
var _ datasource.DataSource = &VcenterDataSource{}

func NewVcenterDataSource() datasource.DataSource {
    return &VcenterDataSource{}
}

// NewVcenterLegacyDataSource registers this data source under the
// name it had before type names kept mixed-case words whole,
// oneuptime_v_center. Deprecated.
func NewVcenterLegacyDataSource() datasource.DataSource {
    return &VcenterDataSource{isLegacyAlias: true}
}

// VcenterDataSource defines the data source implementation.
type VcenterDataSource struct {
    client *Client
    // Registered under the name this data source had before.
    isLegacyAlias bool
}

// VcenterDataSourceModel describes the data source data model.
type VcenterDataSourceModel struct {
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
    DatacenterCount types.Number `tfsdk:"datacenter_count"`
    ClusterCount types.Number `tfsdk:"cluster_count"`
    HostCount types.Number `tfsdk:"host_count"`
    VmCount types.Number `tfsdk:"vm_count"`
    PoweredOnVmCount types.Number `tfsdk:"powered_on_vm_count"`
    DatastoreCount types.Number `tfsdk:"datastore_count"`
    ResourcePoolCount types.Number `tfsdk:"resource_pool_count"`
    DatastoreCapacityBytes types.Number `tfsdk:"datastore_capacity_bytes"`
    DatastoreUsedBytes types.Number `tfsdk:"datastore_used_bytes"`
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

func (d *VcenterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    if d.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_v_center"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_vcenter"
}

func (d *VcenterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "vSphere endpoints (a vCenter Server, or a standalone ESXi host) that are being monitored in this project. Each vCenter is auto-discovered when the OneUptime VMware Agent sends metrics, or can be manually registered. Look up an existing vcenter by `id`, or by any of its other arguments (`name`, `agent_version`, `ai_access_last_error`, ...): each one set must match, and exactly one vcenter may match them all.",

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
                MarkdownDescription: "Name of this vCenter (a vCenter Server, or a standalone ESXi host). This is the join key — it must match the vmware.vcenter.name OTel resource attribute stamped by the OneUptime VMware Agent (VMWARE_VCENTER_NAME).",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this vCenter.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime VMware Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When metrics were last received from this vCenter.",
                Computed: true,
            },
            "datacenter_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of vSphere datacenters reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries datacenter metrics.",
                Optional: true,
                Computed: true,
            },
            "cluster_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of vSphere clusters reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries cluster metrics.",
                Optional: true,
                Computed: true,
            },
            "host_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of ESXi hosts reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries host metrics.",
                Optional: true,
                Computed: true,
            },
            "vm_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of virtual machines (excluding VM templates) reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries virtual machine metrics.",
                Optional: true,
                Computed: true,
            },
            "powered_on_vm_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of virtual machines inferred to be powered on (the vcenter receiver emits CPU metrics only for powered-on VMs). Rendered as 'VMs X/Y powered on' next to vmCount.",
                Optional: true,
                Computed: true,
            },
            "datastore_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of datastores reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries datastore metrics.",
                Optional: true,
                Computed: true,
            },
            "resource_pool_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of resource pools reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries resource pool metrics.",
                Optional: true,
                Computed: true,
            },
            "datastore_capacity_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached total capacity in bytes summed over every datastore reported through this vCenter (used + available vcenter.datastore.disk.usage). The denominator for the storage-used bar on the overview page. Stored as bigint.",
                Optional: true,
                Computed: true,
            },
            "datastore_used_bytes": schema.NumberAttribute{
                MarkdownDescription: "Cached used space in bytes summed over every datastore reported through this vCenter (vcenter.datastore.disk.usage{disk_state=used}). Stored as bigint.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this vCenter archived? Archived vCenters are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this vCenter archived?",
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
                MarkdownDescription: "Number of days to retain telemetry data for this vCenter. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this vCenter (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the vCenter default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only commands (govc about, ls, vm.info, host.info, events) on this vCenter, through its VMware AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the vCenter can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this vCenter. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the vCenter's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the vCenter can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.",
                Optional: true,
                Computed: true,
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of command patterns that Automatic mode may run on this vCenter without approval even though they are riskier changes. Each pattern is one command line for this vCenter's agent (govc) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this vCenter is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the vCenter can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When a command from OneUptime AI last succeeded on this vCenter through its VMware AI agent. Set by the server.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running a command on this vCenter, kept until the next successful command. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this vCenter was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a VMware AI agent that registers later never overwrites a setting an operator chose.",
                Computed: true,
            },
        },
    }
    if d.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_v_center has been renamed to oneuptime_vcenter. Use the new name; the old one keeps working until then."
    }
}

func (d *VcenterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VcenterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data VcenterDataSourceModel

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
    if !data.DatacenterCount.IsNull() && !data.DatacenterCount.IsUnknown() {
        filters["datacenterCount"] = lookupNumber(data.DatacenterCount)
        filterNames = append(filterNames, "datacenter_count = "+data.DatacenterCount.ValueBigFloat().String())
    }
    if !data.ClusterCount.IsNull() && !data.ClusterCount.IsUnknown() {
        filters["clusterCount"] = lookupNumber(data.ClusterCount)
        filterNames = append(filterNames, "cluster_count = "+data.ClusterCount.ValueBigFloat().String())
    }
    if !data.HostCount.IsNull() && !data.HostCount.IsUnknown() {
        filters["hostCount"] = lookupNumber(data.HostCount)
        filterNames = append(filterNames, "host_count = "+data.HostCount.ValueBigFloat().String())
    }
    if !data.VmCount.IsNull() && !data.VmCount.IsUnknown() {
        filters["vmCount"] = lookupNumber(data.VmCount)
        filterNames = append(filterNames, "vm_count = "+data.VmCount.ValueBigFloat().String())
    }
    if !data.PoweredOnVmCount.IsNull() && !data.PoweredOnVmCount.IsUnknown() {
        filters["poweredOnVmCount"] = lookupNumber(data.PoweredOnVmCount)
        filterNames = append(filterNames, "powered_on_vm_count = "+data.PoweredOnVmCount.ValueBigFloat().String())
    }
    if !data.DatastoreCount.IsNull() && !data.DatastoreCount.IsUnknown() {
        filters["datastoreCount"] = lookupNumber(data.DatastoreCount)
        filterNames = append(filterNames, "datastore_count = "+data.DatastoreCount.ValueBigFloat().String())
    }
    if !data.ResourcePoolCount.IsNull() && !data.ResourcePoolCount.IsUnknown() {
        filters["resourcePoolCount"] = lookupNumber(data.ResourcePoolCount)
        filterNames = append(filterNames, "resource_pool_count = "+data.ResourcePoolCount.ValueBigFloat().String())
    }
    if !data.DatastoreCapacityBytes.IsNull() && !data.DatastoreCapacityBytes.IsUnknown() {
        filters["datastoreCapacityBytes"] = lookupNumber(data.DatastoreCapacityBytes)
        filterNames = append(filterNames, "datastore_capacity_bytes = "+data.DatastoreCapacityBytes.ValueBigFloat().String())
    }
    if !data.DatastoreUsedBytes.IsNull() && !data.DatastoreUsedBytes.IsUnknown() {
        filters["datastoreUsedBytes"] = lookupNumber(data.DatastoreUsedBytes)
        filterNames = append(filterNames, "datastore_used_bytes = "+data.DatastoreUsedBytes.ValueBigFloat().String())
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
            "Look the vcenter up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the vcenter up by.",
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
        "datacenterCount": true,
        "clusterCount": true,
        "hostCount": true,
        "vmCount": true,
        "poweredOnVmCount": true,
        "datastoreCount": true,
        "resourcePoolCount": true,
        "datastoreCapacityBytes": true,
        "datastoreUsedBytes": true,
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
        readPath := "/vmware-vcenter/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vcenter, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No vcenter found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read vcenter: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/vmware-vcenter/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list vcenter, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list vcenter: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No vcenter matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one vcenter matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for vcenter.")
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
    if val, ok := item["datacenterCount"].(float64); ok {
        data.DatacenterCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["datacenterCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DatacenterCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.DatacenterCount = types.NumberNull()
        }
    } else {
        data.DatacenterCount = types.NumberNull()
    }
    if val, ok := item["clusterCount"].(float64); ok {
        data.ClusterCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["clusterCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ClusterCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ClusterCount = types.NumberNull()
        }
    } else {
        data.ClusterCount = types.NumberNull()
    }
    if val, ok := item["hostCount"].(float64); ok {
        data.HostCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["hostCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.HostCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.HostCount = types.NumberNull()
        }
    } else {
        data.HostCount = types.NumberNull()
    }
    if val, ok := item["vmCount"].(float64); ok {
        data.VmCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["vmCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.VmCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.VmCount = types.NumberNull()
        }
    } else {
        data.VmCount = types.NumberNull()
    }
    if val, ok := item["poweredOnVmCount"].(float64); ok {
        data.PoweredOnVmCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["poweredOnVmCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PoweredOnVmCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.PoweredOnVmCount = types.NumberNull()
        }
    } else {
        data.PoweredOnVmCount = types.NumberNull()
    }
    if val, ok := item["datastoreCount"].(float64); ok {
        data.DatastoreCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["datastoreCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DatastoreCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.DatastoreCount = types.NumberNull()
        }
    } else {
        data.DatastoreCount = types.NumberNull()
    }
    if val, ok := item["resourcePoolCount"].(float64); ok {
        data.ResourcePoolCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["resourcePoolCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ResourcePoolCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResourcePoolCount = types.NumberNull()
        }
    } else {
        data.ResourcePoolCount = types.NumberNull()
    }
    if val, ok := item["datastoreCapacityBytes"].(float64); ok {
        data.DatastoreCapacityBytes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["datastoreCapacityBytes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DatastoreCapacityBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.DatastoreCapacityBytes = types.NumberNull()
        }
    } else {
        data.DatastoreCapacityBytes = types.NumberNull()
    }
    if val, ok := item["datastoreUsedBytes"].(float64); ok {
        data.DatastoreUsedBytes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["datastoreUsedBytes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DatastoreUsedBytes = types.NumberValue(big.NewFloat(val))
        } else {
            data.DatastoreUsedBytes = types.NumberNull()
        }
    } else {
        data.DatastoreUsedBytes = types.NumberNull()
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
