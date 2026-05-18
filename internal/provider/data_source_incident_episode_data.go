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
var _ datasource.DataSource = &IncidentEpisodeDataDataSource{}

func NewIncidentEpisodeDataDataSource() datasource.DataSource {
    return &IncidentEpisodeDataDataSource{}
}

// IncidentEpisodeDataDataSource defines the data source implementation.
type IncidentEpisodeDataDataSource struct {
    client *Client
}

// IncidentEpisodeDataDataSourceModel describes the data source data model.
type IncidentEpisodeDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    EpisodeNumber types.Number `tfsdk:"episode_number"`
    EpisodeNumberWithPrefix types.String `tfsdk:"episode_number_with_prefix"`
    CurrentIncidentStateId types.String `tfsdk:"current_incident_state_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    RootCause types.String `tfsdk:"root_cause"`
    LastIncidentAddedAt types.String `tfsdk:"last_incident_added_at"`
    ResolvedAt types.String `tfsdk:"resolved_at"`
    AllIncidentsResolvedAt types.String `tfsdk:"all_incidents_resolved_at"`
    AssignedToUserId types.String `tfsdk:"assigned_to_user_id"`
    AssignedToTeamId types.String `tfsdk:"assigned_to_team_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    IsOnCallPolicyExecuted types.Bool `tfsdk:"is_on_call_policy_executed"`
    IncidentCount types.Number `tfsdk:"incident_count"`
    TitleTemplate types.String `tfsdk:"title_template"`
    DescriptionTemplate types.String `tfsdk:"description_template"`
    IsManuallyCreated types.Bool `tfsdk:"is_manually_created"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsOwnerNotifiedOfEpisodeCreation types.Bool `tfsdk:"is_owner_notified_of_episode_creation"`
    GroupingKey types.String `tfsdk:"grouping_key"`
    IncidentGroupingRuleId types.String `tfsdk:"incident_grouping_rule_id"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    PostmortemNote types.String `tfsdk:"postmortem_note"`
    PostUpdatesToWorkspaceChannels types.String `tfsdk:"post_updates_to_workspace_channels"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    DeclaredAt types.String `tfsdk:"declared_at"`
    ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_episode_created"`
    SubscriberNotificationStatusOnEpisodeCreated types.String `tfsdk:"subscriber_notification_status_on_episode_created"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    IsPrivate types.Bool `tfsdk:"is_private"`
}

func (d *IncidentEpisodeDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_episode_data"
}

func (d *IncidentEpisodeDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_episode_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this incident episode. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this incident episode. This is in markdown format.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "episode_number": schema.NumberAttribute{
                MarkdownDescription: "Auto-incrementing episode number per project. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "episode_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Episode number with prefix (e.g., 'IE-42' or '#42'). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "User-documented root cause of this episode. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "last_incident_added_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "resolved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "all_incidents_resolved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "assigned_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "assigned_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for this episode.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_on_call_policy_executed": schema.BoolAttribute{
                MarkdownDescription: "Whether the on-call policy has been executed for this episode. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incident_count": schema.NumberAttribute{
                MarkdownDescription: "Denormalized count of incidents in this episode. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "title_template": schema.StringAttribute{
                MarkdownDescription: "Template used to generate the episode title. Stored for dynamic variable updates.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description_template": schema.StringAttribute{
                MarkdownDescription: "Template used to generate the episode description. Stored for dynamic variable updates.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_manually_created": schema.BoolAttribute{
                MarkdownDescription: "Whether this episode was manually created vs auto-created by a rule. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_owner_notified_of_episode_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified when this episode is created?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "grouping_key": schema.StringAttribute{
                MarkdownDescription: "Key used for grouping incidents into this episode. Generated from groupByFields of the matching rule.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incident_grouping_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "User-documented remediation steps and notes for this episode. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "postmortem_note": schema.StringAttribute{
                MarkdownDescription: "User-documented postmortem summary for this episode. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "post_updates_to_workspace_channels": schema.StringAttribute{
                MarkdownDescription: "Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams). Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should this episode be visible on the status page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
            "declared_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_episode_created": schema.BoolAttribute{
                MarkdownDescription: "Should status page subscribers be notified when this episode is created?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "subscriber_notification_status_on_episode_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers when this episode was created. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_private": schema.BoolAttribute{
                MarkdownDescription: "If true, this incident episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Create Incident Episode], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Admin, Incident Member, Incident Viewer, Read Incident Episode], Update: [Project Owner, Project Admin, Project Member, Incident Admin, Incident Member, Edit Incident Episode]",
                Computed: true,
            },
        },
    }
}

func (d *IncidentEpisodeDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentEpisodeDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentEpisodeDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-episode" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_episode_data, got error: %s", err))
        return
    }

    var incidentEpisodeDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentEpisodeDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_episode_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentEpisodeDataResponse["data"].(map[string]interface{}); ok {
        incidentEpisodeDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentEpisodeDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentEpisodeDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["episode_number"].(float64); ok {
        data.EpisodeNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentEpisodeDataResponse["episode_number_with_prefix"].(string); ok {
        data.EpisodeNumberWithPrefix = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["current_incident_state_id"].(string); ok {
        data.CurrentIncidentStateId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["incident_severity_id"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["root_cause"].(string); ok {
        data.RootCause = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["last_incident_added_at"].(string); ok {
        data.LastIncidentAddedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["resolved_at"].(string); ok {
        data.ResolvedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["all_incidents_resolved_at"].(string); ok {
        data.AllIncidentsResolvedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["assigned_to_user_id"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["assigned_to_team_id"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["on_call_duty_policies"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.OnCallDutyPolicies = setValue
    }
    if val, ok := incidentEpisodeDataResponse["is_on_call_policy_executed"].(bool); ok {
        data.IsOnCallPolicyExecuted = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["incident_count"].(float64); ok {
        data.IncidentCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentEpisodeDataResponse["title_template"].(string); ok {
        data.TitleTemplate = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["description_template"].(string); ok {
        data.DescriptionTemplate = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["is_manually_created"].(bool); ok {
        data.IsManuallyCreated = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := incidentEpisodeDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["is_owner_notified_of_episode_creation"].(bool); ok {
        data.IsOwnerNotifiedOfEpisodeCreation = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["grouping_key"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["incident_grouping_rule_id"].(string); ok {
        data.IncidentGroupingRuleId = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["remediation_notes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["postmortem_note"].(string); ok {
        data.PostmortemNote = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["post_updates_to_workspace_channels"].(string); ok {
        data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["is_visible_on_status_page"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["declared_at"].(string); ok {
        data.DeclaredAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["should_status_page_subscribers_be_notified_on_episode_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated = types.BoolValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["subscriber_notification_status_on_episode_created"].(string); ok {
        data.SubscriberNotificationStatusOnEpisodeCreated = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["subscriber_notification_status_message"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    }
    if val, ok := incidentEpisodeDataResponse["is_private"].(bool); ok {
        data.IsPrivate = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
