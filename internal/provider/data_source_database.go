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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Database servers this project runs or talks to. Each database is discovered from OpenTelemetry Collector database receivers, from database calls in application traces, from database workloads on monitored Kubernetes clusters and Docker / Podman hosts, or added manually. Look up an existing database by `id` or by `name`.",

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
                MarkdownDescription: "Friendly description for this database.",
                Computed: true,
            },
            "database_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identity of this database within the project: the engine family and its canonical endpoint joined with '|' (e.g. postgresql|orders-db.internal:5432), or the workload form for databases discovered on Kubernetes / Docker / Podman (e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db). A fork keys like the engine it forks (a MariaDB database is mysql|..., a Valkey one redis|...), so learning which of the two a database is never changes its identity. Computed by the server for manually added databases..",
                Computed: true,
            },
            "workload_identifier": schema.StringAttribute{
                MarkdownDescription: "Identity of the Kubernetes / Docker / Podman workload this database runs as, e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db. The prefix is the engine family, not the engine (a MariaDB workload is mysql|..., a Valkey one redis|...). Empty for databases that are only known by their endpoint..",
                Computed: true,
            },
            "db_system": schema.StringAttribute{
                MarkdownDescription: "Database engine family as an OpenTelemetry db.system.name value, e.g. postgresql, mysql, redis, mongodb, microsoft.sql_server..",
                Computed: true,
            },
            "server_address": schema.StringAttribute{
                MarkdownDescription: "Host name or IP address of the primary endpoint of this database (the OpenTelemetry server.address attribute)..",
                Computed: true,
            },
            "server_port": schema.NumberAttribute{
                MarkdownDescription: "Port of the primary endpoint of this database (the OpenTelemetry server.port attribute). Defaults to the engine's standard port when not given..",
                Computed: true,
            },
            "db_version": schema.StringAttribute{
                MarkdownDescription: "Engine version last reported for this database, from the collector receiver or the container image tag..",
                Computed: true,
            },
            "discovery_source": schema.StringAttribute{
                MarkdownDescription: "How this database was first discovered: collector, client-spans, kubernetes, docker, podman or manual..",
                Computed: true,
            },
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "kubernetes_namespace": schema.StringAttribute{
                MarkdownDescription: "Kubernetes namespace of the workload this database runs as..",
                Computed: true,
            },
            "workload_kind": schema.StringAttribute{
                MarkdownDescription: "Kind of workload this database runs as, e.g. StatefulSet, Deployment, Cluster (an operator-managed cluster) or Container..",
                Computed: true,
            },
            "workload_name": schema.StringAttribute{
                MarkdownDescription: "Name of the workload (StatefulSet, Deployment, operator cluster or container group) this database runs as..",
                Computed: true,
            },
            "docker_host_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "podman_host_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "member_entity_keys": schema.StringAttribute{
                MarkdownDescription: "Telemetry entity keys of the pods / containers this database runs as, each with the time it was last seen. Maintained by discovery..",
                Computed: true,
            },
            "instance_count": schema.NumberAttribute{
                MarkdownDescription: "Number of running instances (pods / containers) last seen for this database's workload..",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Whether engine metrics are currently being received from a collector / agent for this database (connected) or a collector reported and has stopped (disconnected). Empty when no collector has ever reported..",
                Computed: true,
            },
            "collector_last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime agent reporting this database, as self-reported via the oneuptime.agent.version resource attribute.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "auto_archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "manually_restored_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "db_system_source": schema.StringAttribute{
                MarkdownDescription: "What the database engine was last determined from: manual, collector, container or client-spans. Stronger evidence may correct the engine; weaker evidence never changes it..",
                Computed: true,
            },
            "workload_last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "automatic_assignments": schema.StringAttribute{
                MarkdownDescription: "Label and owner ids that label rules, owner rules or telemetry attached automatically. Maintained by OneUptime..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this database archived? Archived databases are hidden from lists but keep collecting telemetry..",
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
                MarkdownDescription: "Number of days to retain telemetry collected from this database by a collector / agent. Leave blank to use the project-wide default..",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for telemetry collected from this database by a collector / agent (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the database default, then the project's retention settings..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a database.",
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
        "deletedByUserId": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No database found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one database matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
