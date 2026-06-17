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
var _ datasource.DataSource = &AlertOwnerRuleDataDataSource{}

func NewAlertOwnerRuleDataDataSource() datasource.DataSource {
    return &AlertOwnerRuleDataDataSource{}
}

// AlertOwnerRuleDataDataSource defines the data source implementation.
type AlertOwnerRuleDataDataSource struct {
    client *Client
}

// AlertOwnerRuleDataDataSourceModel describes the data source data model.
type AlertOwnerRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    NotifyOwners types.Bool `tfsdk:"notify_owners"`
    Monitors types.Set `tfsdk:"monitors"`
    AlertSeverities types.Set `tfsdk:"alert_severities"`
    AlertLabels types.Set `tfsdk:"alert_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    AlertTitlePattern types.String `tfsdk:"alert_title_pattern"`
    AlertDescriptionPattern types.String `tfsdk:"alert_description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    InheritOwnersFromMonitors types.Bool `tfsdk:"inherit_owners_from_monitors"`
    InheritOwnersFromHosts types.Bool `tfsdk:"inherit_owners_from_hosts"`
    InheritOwnersFromKubernetesClusters types.Bool `tfsdk:"inherit_owners_from_kubernetes_clusters"`
    InheritOwnersFromDockerHosts types.Bool `tfsdk:"inherit_owners_from_docker_hosts"`
    InheritOwnersFromPodmanHosts types.Bool `tfsdk:"inherit_owners_from_podman_hosts"`
    InheritOwnersFromServices types.Bool `tfsdk:"inherit_owners_from_services"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AlertOwnerRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_owner_rule_data"
}

func (d *AlertOwnerRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_owner_rule_data data source",

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
                MarkdownDescription: "Description of this alert owner rule. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Send notifications to owner users and teams when they are added by this rule. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for alerts with these severities. Leave empty to match alerts of any severity.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for alerts that have at least one of these labels. Leave empty to match regardless of alert labels.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for alerts from monitors that have at least one of these labels.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the alert title.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "alert_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the alert description.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the alert's monitor name.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the alert's monitor description.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the alert when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the alert when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "inherit_owners_from_monitors": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the alert's monitor to the alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "inherit_owners_from_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the alert's affected hosts to the alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "inherit_owners_from_kubernetes_clusters": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the alert's affected Kubernetes clusters to the alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "inherit_owners_from_docker_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the alert's affected Docker hosts to the alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "inherit_owners_from_podman_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the alert's affected Podman hosts to the alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "inherit_owners_from_services": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the alert's affected services to the alert.. Permissions - Create: [Project Owner, Project Admin, Create Alert Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Owner Rule], Update: [Project Owner, Project Admin, Edit Alert Owner Rule]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *AlertOwnerRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertOwnerRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertOwnerRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-owner-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_owner_rule_data, got error: %s", err))
        return
    }

    var alertOwnerRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertOwnerRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_owner_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertOwnerRuleDataResponse["data"].(map[string]interface{}); ok {
        alertOwnerRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertOwnerRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertOwnerRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["notify_owners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["monitors"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Monitors = setValue
    }
    if val, ok := alertOwnerRuleDataResponse["alert_severities"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AlertSeverities = setValue
    }
    if val, ok := alertOwnerRuleDataResponse["alert_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AlertLabels = setValue
    }
    if val, ok := alertOwnerRuleDataResponse["monitor_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.MonitorLabels = setValue
    }
    if val, ok := alertOwnerRuleDataResponse["alert_title_pattern"].(string); ok {
        data.AlertTitlePattern = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["alert_description_pattern"].(string); ok {
        data.AlertDescriptionPattern = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["monitor_name_pattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["monitor_description_pattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["owner_users"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OwnerUsers = setValue
    }
    if val, ok := alertOwnerRuleDataResponse["owner_teams"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OwnerTeams = setValue
    }
    if val, ok := alertOwnerRuleDataResponse["inherit_owners_from_monitors"].(bool); ok {
        data.InheritOwnersFromMonitors = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["inherit_owners_from_hosts"].(bool); ok {
        data.InheritOwnersFromHosts = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["inherit_owners_from_kubernetes_clusters"].(bool); ok {
        data.InheritOwnersFromKubernetesClusters = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["inherit_owners_from_docker_hosts"].(bool); ok {
        data.InheritOwnersFromDockerHosts = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["inherit_owners_from_podman_hosts"].(bool); ok {
        data.InheritOwnersFromPodmanHosts = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["inherit_owners_from_services"].(bool); ok {
        data.InheritOwnersFromServices = types.BoolValue(val)
    }
    if val, ok := alertOwnerRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
