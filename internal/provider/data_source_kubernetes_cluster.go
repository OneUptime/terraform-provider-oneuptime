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
var _ datasource.DataSource = &KubernetesClusterDataSource{}

func NewKubernetesClusterDataSource() datasource.DataSource {
    return &KubernetesClusterDataSource{}
}

// KubernetesClusterDataSource defines the data source implementation.
type KubernetesClusterDataSource struct {
    client *Client
}

// KubernetesClusterDataSourceModel describes the data source data model.
type KubernetesClusterDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    ClusterIdentifier types.String `tfsdk:"cluster_identifier"`
    ProviderValue types.String `tfsdk:"provider_value"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    NodeCount types.Number `tfsdk:"node_count"`
    PodCount types.Number `tfsdk:"pod_count"`
    NamespaceCount types.Number `tfsdk:"namespace_count"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    AiAccessRunnerId types.String `tfsdk:"ai_access_runner_id"`
    AiAccessCredentialId types.String `tfsdk:"ai_access_credential_id"`
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiKubectlCommandAllowlist types.String `tfsdk:"ai_kubectl_command_allowlist"`
    AiAccessLastVerifiedAt types.String `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt types.String `tfsdk:"ai_access_configured_at"`
    AiAccessRunnerBoundAt types.String `tfsdk:"ai_access_runner_bound_at"`
}

func (d *KubernetesClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster"
}

func (d *KubernetesClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered. Look up an existing kubernetes_cluster by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this Kubernetes cluster.",
                Computed: true,
            },
            "cluster_identifier": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for this cluster, sourced from the k8s.cluster.name OTel resource attribute.",
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Cloud provider or platform running this cluster (EKS, GKE, AKS, self-managed, unknown).",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected).",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime Kubernetes agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes in this cluster.",
                Computed: true,
            },
            "pod_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of pods in this cluster.",
                Computed: true,
            },
            "namespace_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of namespaces in this cluster.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this Kubernetes cluster archived? Archived Kubernetes clusters are hidden from lists but keep collecting telemetry..",
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this Kubernetes cluster. Leave blank to use the project-wide default..",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this Kubernetes cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the cluster default, then the project's retention settings..",
                Computed: true,
            },
            "ai_access_runner_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "ai_access_credential_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI runs read-only kubectl commands (get, describe, logs, events, top, rollout status) on this cluster while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. Anyone who may edit the cluster can turn it on or off..",
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this cluster. RequireApproval: AI composes a kubectl plan and a human approves it with one click before anything runs. Any follow-up plan asks again. Automatic: AI runs safe changes without a human — each on ONE named object: rollout restart/undo/pause/resume of one workload, scale one workload above zero, delete one named pod, cordon/uncordon one node, label/annotate one pod or workload with unreserved keys. A riskier change (patch, set image, drain, taint, scale to zero, deleting workloads or jobs, anything touching several objects) never runs without one: when the round could only find riskier fixes it ends by proposing exactly those for one-click approval; when it also ran safe fixes, a riskier fix is proposed only if verification shows the safe ones did not recover the signal (the follow-up round, which asks). Shapes on the cluster's kubectl allowlist run on their own. BypassApproval: AI does not ask. Every change the policy allows — safe AND riskier — runs on its own, follow-up rounds included, except for what always asks (below). In EVERY mode, Bypass approval included: destructive commands (Denied tier) never run; a write in a protected namespace (kube-system, kube-public, kube-node-lease), a node drain, a node taint and a patch of a Node always need a human; the in-cluster Runner never changes its own namespace or anything outside the namespaces its chart may write; and an unattended run becomes a proposal when the hourly per-cluster circuit breaker trips or another unattended round already holds the cluster. Anyone who may edit the cluster can lower the mode (to Disabled, RequireApproval, or from BypassApproval to Automatic); raising it to Automatic or BypassApproval needs Project Owner, Project Admin or Edit Auto Remediation Rule..",
                Computed: true,
            },
            "ai_kubectl_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of kubectl command patterns that Automatic mode may run without approval even though they are riskier changes, for example: [\"kubectl set image deployment/web * -n web\"]. A pattern is compared with the command word by word: * stands for exactly one word (an image, a name), never for extra objects, flags or a second -n, every flag the command uses must be written out in the pattern, and the leading \"kubectl\" is optional. At most 100 patterns of at most 500 characters each; a pattern that is not one kubectl command line is refused. Destructive commands (Denied tier) never run regardless, and a write in a protected namespace (kube-system, kube-public, kube-node-lease) or a node drain still needs a human. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the cluster can remove patterns or clear the list..",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit while running kubectl on this cluster, kept until the next successful command. Set by the server..",
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "ai_access_runner_bound_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *KubernetesClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KubernetesClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data KubernetesClusterDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a kubernetes_cluster.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "slug": true,
        "description": true,
        "clusterIdentifier": true,
        "provider": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "nodeCount": true,
        "podCount": true,
        "namespaceCount": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "aiAccessRunnerId": true,
        "aiAccessCredentialId": true,
        "isAiInvestigationEnabled": true,
        "aiRemediationMode": true,
        "aiKubectlCommandAllowlist": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "aiAccessRunnerBoundAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/kubernetes-cluster/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kubernetes_cluster, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No kubernetes_cluster found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read kubernetes_cluster: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/kubernetes-cluster/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list kubernetes_cluster, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list kubernetes_cluster: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No kubernetes_cluster found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one kubernetes_cluster matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for kubernetes_cluster.")
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["clusterIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClusterIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClusterIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClusterIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClusterIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.ClusterIdentifier = types.StringNull()
        }
    } else if val, ok := item["clusterIdentifier"].(string); ok {
        data.ClusterIdentifier = types.StringValue(val)
    } else {
        data.ClusterIdentifier = types.StringNull()
    }
    if obj, ok := item["provider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := item["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
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
    if val, ok := item["nodeCount"].(float64); ok {
        data.NodeCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["nodeCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.NodeCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.NodeCount = types.NumberNull()
        }
    } else {
        data.NodeCount = types.NumberNull()
    }
    if val, ok := item["podCount"].(float64); ok {
        data.PodCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["podCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PodCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.PodCount = types.NumberNull()
        }
    } else {
        data.PodCount = types.NumberNull()
    }
    if val, ok := item["namespaceCount"].(float64); ok {
        data.NamespaceCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["namespaceCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.NamespaceCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.NamespaceCount = types.NumberNull()
        }
    } else {
        data.NamespaceCount = types.NumberNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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
    if obj, ok := item["aiAccessRunnerId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessRunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessRunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessRunnerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessRunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessRunnerId = types.StringNull()
        }
    } else if val, ok := item["aiAccessRunnerId"].(string); ok {
        data.AiAccessRunnerId = types.StringValue(val)
    } else {
        data.AiAccessRunnerId = types.StringNull()
    }
    if obj, ok := item["aiAccessCredentialId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessCredentialId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessCredentialId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessCredentialId = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessCredentialId = types.StringNull()
        }
    } else if val, ok := item["aiAccessCredentialId"].(string); ok {
        data.AiAccessCredentialId = types.StringValue(val)
    } else {
        data.AiAccessCredentialId = types.StringNull()
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
    if obj, ok := item["aiKubectlCommandAllowlist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiKubectlCommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiKubectlCommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiKubectlCommandAllowlist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiKubectlCommandAllowlist = types.StringValue(string(jsonBytes))
        } else {
            data.AiKubectlCommandAllowlist = types.StringNull()
        }
    } else if val, ok := item["aiKubectlCommandAllowlist"].(string); ok {
        data.AiKubectlCommandAllowlist = types.StringValue(val)
    } else {
        data.AiKubectlCommandAllowlist = types.StringNull()
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
    if obj, ok := item["aiAccessRunnerBoundAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessRunnerBoundAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessRunnerBoundAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessRunnerBoundAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessRunnerBoundAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessRunnerBoundAt = types.StringNull()
        }
    } else if val, ok := item["aiAccessRunnerBoundAt"].(string); ok {
        data.AiAccessRunnerBoundAt = types.StringValue(val)
    } else {
        data.AiAccessRunnerBoundAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
