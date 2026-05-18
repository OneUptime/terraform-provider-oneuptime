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
var _ datasource.DataSource = &KubernetesClusterLabelRuleDataDataSource{}

func NewKubernetesClusterLabelRuleDataDataSource() datasource.DataSource {
    return &KubernetesClusterLabelRuleDataDataSource{}
}

// KubernetesClusterLabelRuleDataDataSource defines the data source implementation.
type KubernetesClusterLabelRuleDataDataSource struct {
    client *Client
}

// KubernetesClusterLabelRuleDataDataSourceModel describes the data source data model.
type KubernetesClusterLabelRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    KubernetesClusterLabels types.Set `tfsdk:"kubernetes_cluster_labels"`
    KubernetesClusterNamePattern types.String `tfsdk:"kubernetes_cluster_name_pattern"`
    KubernetesClusterDescriptionPattern types.String `tfsdk:"kubernetes_cluster_description_pattern"`
    LabelsToAdd types.Set `tfsdk:"labels_to_add"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *KubernetesClusterLabelRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_label_rule_data"
}

func (d *KubernetesClusterLabelRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "kubernetes_cluster_label_rule_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this Kubernetes cluster label rule. Permissions - Create: [Project Owner, Project Admin, Create Kubernetes Cluster Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Kubernetes Cluster Label Rule], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster Label Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Kubernetes Cluster Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Kubernetes Cluster Label Rule], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster Label Rule]",
                Computed: true,
            },
            "kubernetes_cluster_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for Kubernetes clusters that already have at least one of these labels. Leave empty to match regardless of labels.. Permissions - Create: [Project Owner, Project Admin, Create Kubernetes Cluster Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Kubernetes Cluster Label Rule], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster Label Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "kubernetes_cluster_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the Kubernetes cluster name. Leave empty to match any name.. Permissions - Create: [Project Owner, Project Admin, Create Kubernetes Cluster Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Kubernetes Cluster Label Rule], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster Label Rule]",
                Computed: true,
            },
            "kubernetes_cluster_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the Kubernetes cluster description. Leave empty to match any description.. Permissions - Create: [Project Owner, Project Admin, Create Kubernetes Cluster Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Kubernetes Cluster Label Rule], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster Label Rule]",
                Computed: true,
            },
            "labels_to_add": schema.SetAttribute{
                MarkdownDescription: "Labels to attach to the Kubernetes cluster when this rule matches. Already-attached labels are not duplicated.. Permissions - Create: [Project Owner, Project Admin, Create Kubernetes Cluster Label Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Kubernetes Cluster Label Rule], Update: [Project Owner, Project Admin, Edit Kubernetes Cluster Label Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *KubernetesClusterLabelRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KubernetesClusterLabelRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data KubernetesClusterLabelRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "kubernetes-cluster-label-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kubernetes_cluster_label_rule_data, got error: %s", err))
        return
    }

    var kubernetesClusterLabelRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &kubernetesClusterLabelRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse kubernetes_cluster_label_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := kubernetesClusterLabelRuleDataResponse["data"].(map[string]interface{}); ok {
        kubernetesClusterLabelRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := kubernetesClusterLabelRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["kubernetes_cluster_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.KubernetesClusterLabels = setValue
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["kubernetes_cluster_name_pattern"].(string); ok {
        data.KubernetesClusterNamePattern = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["kubernetes_cluster_description_pattern"].(string); ok {
        data.KubernetesClusterDescriptionPattern = types.StringValue(val)
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["labels_to_add"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.LabelsToAdd = setValue
    }
    if val, ok := kubernetesClusterLabelRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
