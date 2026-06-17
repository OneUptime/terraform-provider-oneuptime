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
var _ datasource.DataSource = &DockerSwarmClusterDataDataSource{}

func NewDockerSwarmClusterDataDataSource() datasource.DataSource {
    return &DockerSwarmClusterDataDataSource{}
}

// DockerSwarmClusterDataDataSource defines the data source implementation.
type DockerSwarmClusterDataDataSource struct {
    client *Client
}

// DockerSwarmClusterDataDataSourceModel describes the data source data model.
type DockerSwarmClusterDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    DockerVersion types.String `tfsdk:"docker_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    NodeCount types.Number `tfsdk:"node_count"`
    ReadyNodeCount types.Number `tfsdk:"ready_node_count"`
    ManagerNodeCount types.Number `tfsdk:"manager_node_count"`
    ServiceCount types.Number `tfsdk:"service_count"`
    TaskCount types.Number `tfsdk:"task_count"`
    RunningTaskCount types.Number `tfsdk:"running_task_count"`
    StackCount types.Number `tfsdk:"stack_count"`
    NetworkCount types.Number `tfsdk:"network_count"`
    SwarmId types.String `tfsdk:"swarm_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *DockerSwarmClusterDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_docker_swarm_cluster_data"
}

func (d *DockerSwarmClusterDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "docker_swarm_cluster_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this DockerSwarm cluster. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Docker Swarm Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime DockerSwarm agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "docker_version": schema.StringAttribute{
                MarkdownDescription: "Docker Engine version reported by the swarm manager this agent talks to.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "ready_node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes whose status is 'ready' in this cluster. Rendered as 'Nodes X/Y ready' next to nodeCount.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "manager_node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes with the manager role in this cluster.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "service_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of swarm services in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "task_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of swarm tasks (service instances) in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "running_task_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of tasks in the running state. Rendered as 'Tasks X/Y running' next to taskCount.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "stack_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of deployed compose stacks in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "network_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of swarm-scoped (overlay) networks in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "swarm_id": schema.StringAttribute{
                MarkdownDescription: "The Docker Swarm cluster ID (docker info -> Swarm.Cluster.ID) reported by the manager. Stable for the lifetime of the swarm; informational only — the join key is the cluster name.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Docker Swarm Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Docker Swarm Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Docker Swarm Cluster]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this DockerSwarm cluster. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Docker Swarm Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Docker Swarm Cluster]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this DockerSwarm cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the DockerSwarm cluster default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Docker Swarm Cluster], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Docker Swarm Cluster], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Docker Swarm Cluster]",
                Computed: true,
            },
        },
    }
}

func (d *DockerSwarmClusterDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DockerSwarmClusterDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DockerSwarmClusterDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "docker-swarm-cluster" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read docker_swarm_cluster_data, got error: %s", err))
        return
    }

    var dockerSwarmClusterDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &dockerSwarmClusterDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse docker_swarm_cluster_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := dockerSwarmClusterDataResponse["data"].(map[string]interface{}); ok {
        dockerSwarmClusterDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := dockerSwarmClusterDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["docker_version"].(string); ok {
        data.DockerVersion = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["node_count"].(float64); ok {
        data.NodeCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["ready_node_count"].(float64); ok {
        data.ReadyNodeCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["manager_node_count"].(float64); ok {
        data.ManagerNodeCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["service_count"].(float64); ok {
        data.ServiceCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["task_count"].(float64); ok {
        data.TaskCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["running_task_count"].(float64); ok {
        data.RunningTaskCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["stack_count"].(float64); ok {
        data.StackCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["network_count"].(float64); ok {
        data.NetworkCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["swarm_id"].(string); ok {
        data.SwarmId = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := dockerSwarmClusterDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := dockerSwarmClusterDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := dockerSwarmClusterDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
