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
var _ datasource.DataSource = &IncidentEpisodeOwnerRuleDataDataSource{}

func NewIncidentEpisodeOwnerRuleDataDataSource() datasource.DataSource {
    return &IncidentEpisodeOwnerRuleDataDataSource{}
}

// IncidentEpisodeOwnerRuleDataDataSource defines the data source implementation.
type IncidentEpisodeOwnerRuleDataDataSource struct {
    client *Client
}

// IncidentEpisodeOwnerRuleDataDataSourceModel describes the data source data model.
type IncidentEpisodeOwnerRuleDataDataSourceModel struct {
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
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    EpisodeLabels types.Set `tfsdk:"episode_labels"`
    EpisodeTitlePattern types.String `tfsdk:"episode_title_pattern"`
    EpisodeDescriptionPattern types.String `tfsdk:"episode_description_pattern"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentEpisodeOwnerRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_episode_owner_rule_data"
}

func (d *IncidentEpisodeOwnerRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_episode_owner_rule_data data source",

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
                MarkdownDescription: "Description of this incident episode owner rule. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Send notifications to owner users and teams when they are added by this rule. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for episodes with these severities. Leave empty to match any severity.. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for episodes that have at least one of these labels.. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "episode_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the episode title.. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
            },
            "episode_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the episode description.. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the incident episode when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the incident episode when this rule matches.. Permissions - Create: [Project Owner, Project Admin, Create Incident Episode Owner Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode Owner Rule], Update: [Project Owner, Project Admin, Edit Incident Episode Owner Rule]",
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

func (d *IncidentEpisodeOwnerRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentEpisodeOwnerRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentEpisodeOwnerRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-episode-owner-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_episode_owner_rule_data, got error: %s", err))
        return
    }

    var incidentEpisodeOwnerRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentEpisodeOwnerRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_episode_owner_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentEpisodeOwnerRuleDataResponse["data"].(map[string]interface{}); ok {
        incidentEpisodeOwnerRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentEpisodeOwnerRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["notify_owners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["incident_severities"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.IncidentSeverities = setValue
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["episode_labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EpisodeLabels = setValue
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["episode_title_pattern"].(string); ok {
        data.EpisodeTitlePattern = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["episode_description_pattern"].(string); ok {
        data.EpisodeDescriptionPattern = types.StringValue(val)
    }
    if val, ok := incidentEpisodeOwnerRuleDataResponse["owner_users"].([]interface{}); ok {
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
    if val, ok := incidentEpisodeOwnerRuleDataResponse["owner_teams"].([]interface{}); ok {
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
    if val, ok := incidentEpisodeOwnerRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
