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
var _ datasource.DataSource = &KubernetesClusterDataDataSource{}

func NewKubernetesClusterDataDataSource() datasource.DataSource {
    return &KubernetesClusterDataDataSource{}
}

// KubernetesClusterDataDataSource defines the data source implementation.
type KubernetesClusterDataDataSource struct {
    client *Client
}

// KubernetesClusterDataDataSourceModel describes the data source data model.
type KubernetesClusterDataDataSourceModel struct {
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
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    NodeCount types.Number `tfsdk:"node_count"`
    PodCount types.Number `tfsdk:"pod_count"`
    NamespaceCount types.Number `tfsdk:"namespace_count"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *KubernetesClusterDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_data"
}

func (d *KubernetesClusterDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "kubernetes_cluster_data data source",

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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this Kubernetes cluster. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Kubernetes Cluster], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Kubernetes Cluster]",
                Computed: true,
            },
            "cluster_identifier": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for this cluster, sourced from the k8s.cluster.name OTel resource attribute. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Kubernetes Cluster], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Kubernetes Cluster]",
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Cloud provider or platform running this cluster (EKS, GKE, AKS, self-managed, unknown). Permissions - Create: [Project Owner, Project Admin, Project Member, Create Kubernetes Cluster], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Kubernetes Cluster]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "node_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of nodes in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster]",
                Computed: true,
            },
            "pod_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of pods in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster]",
                Computed: true,
            },
            "namespace_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of namespaces in this cluster. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster]",
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
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Kubernetes Cluster], Read: [Project Owner, Project Admin, Project Member, Read Kubernetes Cluster, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Edit Kubernetes Cluster]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *KubernetesClusterDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KubernetesClusterDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data KubernetesClusterDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "kubernetes-cluster" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kubernetes_cluster_data, got error: %s", err))
        return
    }

    var kubernetesClusterDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &kubernetesClusterDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse kubernetes_cluster_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := kubernetesClusterDataResponse["data"].(map[string]interface{}); ok {
        kubernetesClusterDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := kubernetesClusterDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesClusterDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["cluster_identifier"].(string); ok {
        data.ClusterIdentifier = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["node_count"].(float64); ok {
        data.NodeCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesClusterDataResponse["pod_count"].(float64); ok {
        data.PodCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesClusterDataResponse["namespace_count"].(float64); ok {
        data.NamespaceCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesClusterDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := kubernetesClusterDataResponse["labels"].([]interface{}); ok {
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
