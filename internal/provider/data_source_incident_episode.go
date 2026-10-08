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
var _ datasource.DataSource = &IncidentEpisodeDataSource{}

func NewIncidentEpisodeDataSource() datasource.DataSource {
    return &IncidentEpisodeDataSource{}
}

// IncidentEpisodeDataSource defines the data source implementation.
type IncidentEpisodeDataSource struct {
    client *Client
}

// IncidentEpisodeDataSourceModel describes the data source data model.
type IncidentEpisodeDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
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

func (d *IncidentEpisodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_episode"
}

func (d *IncidentEpisodeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage incident episodes (groups of related incidents) for your project Look up an existing incident episode by `id`, or by any of its other arguments (`assigned_to_team_id`, `assigned_to_user_id`, `created_by_user_id`, ...): each one set must match, and exactly one incident episode may match them all.",

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
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this incident episode.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this incident episode. This is in markdown format.",
                Optional: true,
                Computed: true,
            },
            "episode_number": schema.NumberAttribute{
                MarkdownDescription: "Auto-incrementing episode number per project.",
                Optional: true,
                Computed: true,
            },
            "episode_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "Episode number with prefix (e.g., 'IE-42' or '#42').",
                Optional: true,
                Computed: true,
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "Current Incident State ID. The ID of a `oneuptime_incident_state`.",
                Optional: true,
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "Incident Severity ID. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "User-documented root cause of this episode.",
                Optional: true,
                Computed: true,
            },
            "last_incident_added_at": schema.StringAttribute{
                MarkdownDescription: "When the last incident was added to this episode.",
                Computed: true,
            },
            "resolved_at": schema.StringAttribute{
                MarkdownDescription: "When this episode was resolved.",
                Computed: true,
            },
            "all_incidents_resolved_at": schema.StringAttribute{
                MarkdownDescription: "When all incidents in this episode were first detected as resolved. Used for resolve delay calculation.",
                Computed: true,
            },
            "assigned_to_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who is assigned to this episode. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "assigned_to_team_id": schema.StringAttribute{
                MarkdownDescription: "Team ID that is assigned to this episode. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "List of on-call duty policies to execute for this episode. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_on_call_policy_executed": schema.BoolAttribute{
                MarkdownDescription: "Whether the on-call policy has been executed for this episode.",
                Optional: true,
                Computed: true,
            },
            "incident_count": schema.NumberAttribute{
                MarkdownDescription: "Denormalized count of incidents in this episode.",
                Optional: true,
                Computed: true,
            },
            "title_template": schema.StringAttribute{
                MarkdownDescription: "Template used to generate the episode title. Stored for dynamic variable updates.",
                Optional: true,
                Computed: true,
            },
            "description_template": schema.StringAttribute{
                MarkdownDescription: "Template used to generate the episode description. Stored for dynamic variable updates.",
                Optional: true,
                Computed: true,
            },
            "is_manually_created": schema.BoolAttribute{
                MarkdownDescription: "Whether this episode was manually created vs auto-created by a rule.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_owner_notified_of_episode_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified when this episode is created?",
                Optional: true,
                Computed: true,
            },
            "grouping_key": schema.StringAttribute{
                MarkdownDescription: "Key used for grouping incidents into this episode. Generated from groupByFields of the matching rule. When a private incident opened the episode, its title is in the key only as a keyed hash.",
                Optional: true,
                Computed: true,
            },
            "incident_grouping_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Incident Grouping Rule that created this episode (if applicable). The ID of a `oneuptime_incident_grouping_rule`.",
                Optional: true,
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "User-documented remediation steps and notes for this episode.",
                Optional: true,
                Computed: true,
            },
            "postmortem_note": schema.StringAttribute{
                MarkdownDescription: "User-documented postmortem summary for this episode.",
                Optional: true,
                Computed: true,
            },
            "post_updates_to_workspace_channels": schema.StringAttribute{
                MarkdownDescription: "Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should this episode be visible on the status page?",
                Optional: true,
                Computed: true,
            },
            "declared_at": schema.StringAttribute{
                MarkdownDescription: "When this episode was declared.",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_episode_created": schema.BoolAttribute{
                MarkdownDescription: "Should status page subscribers be notified when this episode is created?",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_on_episode_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers when this episode was created.",
                Optional: true,
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.",
                Optional: true,
                Computed: true,
            },
            "is_private": schema.BoolAttribute{
                MarkdownDescription: "If true, this incident episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentEpisodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentEpisodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentEpisodeDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.EpisodeNumber.IsNull() && !data.EpisodeNumber.IsUnknown() {
        filters["episodeNumber"] = lookupNumber(data.EpisodeNumber)
        filterNames = append(filterNames, "episode_number = "+data.EpisodeNumber.ValueBigFloat().String())
    }
    if !data.EpisodeNumberWithPrefix.IsNull() && !data.EpisodeNumberWithPrefix.IsUnknown() {
        filters["episodeNumberWithPrefix"] = data.EpisodeNumberWithPrefix.ValueString()
        filterNames = append(filterNames, "episode_number_with_prefix = "+fmt.Sprintf("%q", data.EpisodeNumberWithPrefix.ValueString()))
    }
    if !data.CurrentIncidentStateId.IsNull() && !data.CurrentIncidentStateId.IsUnknown() {
        filters["currentIncidentStateId"] = data.CurrentIncidentStateId.ValueString()
        filterNames = append(filterNames, "current_incident_state_id = "+fmt.Sprintf("%q", data.CurrentIncidentStateId.ValueString()))
    }
    if !data.IncidentSeverityId.IsNull() && !data.IncidentSeverityId.IsUnknown() {
        filters["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
        filterNames = append(filterNames, "incident_severity_id = "+fmt.Sprintf("%q", data.IncidentSeverityId.ValueString()))
    }
    if !data.RootCause.IsNull() && !data.RootCause.IsUnknown() {
        filters["rootCause"] = data.RootCause.ValueString()
        filterNames = append(filterNames, "root_cause = "+fmt.Sprintf("%q", data.RootCause.ValueString()))
    }
    if !data.AssignedToUserId.IsNull() && !data.AssignedToUserId.IsUnknown() {
        filters["assignedToUserId"] = data.AssignedToUserId.ValueString()
        filterNames = append(filterNames, "assigned_to_user_id = "+fmt.Sprintf("%q", data.AssignedToUserId.ValueString()))
    }
    if !data.AssignedToTeamId.IsNull() && !data.AssignedToTeamId.IsUnknown() {
        filters["assignedToTeamId"] = data.AssignedToTeamId.ValueString()
        filterNames = append(filterNames, "assigned_to_team_id = "+fmt.Sprintf("%q", data.AssignedToTeamId.ValueString()))
    }
    if !data.IsOnCallPolicyExecuted.IsNull() && !data.IsOnCallPolicyExecuted.IsUnknown() {
        filters["isOnCallPolicyExecuted"] = data.IsOnCallPolicyExecuted.ValueBool()
        filterNames = append(filterNames, "is_on_call_policy_executed = "+fmt.Sprintf("%t", data.IsOnCallPolicyExecuted.ValueBool()))
    }
    if !data.IncidentCount.IsNull() && !data.IncidentCount.IsUnknown() {
        filters["incidentCount"] = lookupNumber(data.IncidentCount)
        filterNames = append(filterNames, "incident_count = "+data.IncidentCount.ValueBigFloat().String())
    }
    if !data.TitleTemplate.IsNull() && !data.TitleTemplate.IsUnknown() {
        filters["titleTemplate"] = data.TitleTemplate.ValueString()
        filterNames = append(filterNames, "title_template = "+fmt.Sprintf("%q", data.TitleTemplate.ValueString()))
    }
    if !data.DescriptionTemplate.IsNull() && !data.DescriptionTemplate.IsUnknown() {
        filters["descriptionTemplate"] = data.DescriptionTemplate.ValueString()
        filterNames = append(filterNames, "description_template = "+fmt.Sprintf("%q", data.DescriptionTemplate.ValueString()))
    }
    if !data.IsManuallyCreated.IsNull() && !data.IsManuallyCreated.IsUnknown() {
        filters["isManuallyCreated"] = data.IsManuallyCreated.ValueBool()
        filterNames = append(filterNames, "is_manually_created = "+fmt.Sprintf("%t", data.IsManuallyCreated.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsOwnerNotifiedOfEpisodeCreation.IsNull() && !data.IsOwnerNotifiedOfEpisodeCreation.IsUnknown() {
        filters["isOwnerNotifiedOfEpisodeCreation"] = data.IsOwnerNotifiedOfEpisodeCreation.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_episode_creation = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfEpisodeCreation.ValueBool()))
    }
    if !data.GroupingKey.IsNull() && !data.GroupingKey.IsUnknown() {
        filters["groupingKey"] = data.GroupingKey.ValueString()
        filterNames = append(filterNames, "grouping_key = "+fmt.Sprintf("%q", data.GroupingKey.ValueString()))
    }
    if !data.IncidentGroupingRuleId.IsNull() && !data.IncidentGroupingRuleId.IsUnknown() {
        filters["incidentGroupingRuleId"] = data.IncidentGroupingRuleId.ValueString()
        filterNames = append(filterNames, "incident_grouping_rule_id = "+fmt.Sprintf("%q", data.IncidentGroupingRuleId.ValueString()))
    }
    if !data.RemediationNotes.IsNull() && !data.RemediationNotes.IsUnknown() {
        filters["remediationNotes"] = data.RemediationNotes.ValueString()
        filterNames = append(filterNames, "remediation_notes = "+fmt.Sprintf("%q", data.RemediationNotes.ValueString()))
    }
    if !data.PostmortemNote.IsNull() && !data.PostmortemNote.IsUnknown() {
        filters["postmortemNote"] = data.PostmortemNote.ValueString()
        filterNames = append(filterNames, "postmortem_note = "+fmt.Sprintf("%q", data.PostmortemNote.ValueString()))
    }
    if !data.IsVisibleOnStatusPage.IsNull() && !data.IsVisibleOnStatusPage.IsUnknown() {
        filters["isVisibleOnStatusPage"] = data.IsVisibleOnStatusPage.ValueBool()
        filterNames = append(filterNames, "is_visible_on_status_page = "+fmt.Sprintf("%t", data.IsVisibleOnStatusPage.ValueBool()))
    }
    if !data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated.IsNull() && !data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated.IsUnknown() {
        filters["shouldStatusPageSubscribersBeNotifiedOnEpisodeCreated"] = data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated.ValueBool()
        filterNames = append(filterNames, "should_status_page_subscribers_be_notified_on_episode_created = "+fmt.Sprintf("%t", data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated.ValueBool()))
    }
    if !data.SubscriberNotificationStatusOnEpisodeCreated.IsNull() && !data.SubscriberNotificationStatusOnEpisodeCreated.IsUnknown() {
        filters["subscriberNotificationStatusOnEpisodeCreated"] = data.SubscriberNotificationStatusOnEpisodeCreated.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_on_episode_created = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusOnEpisodeCreated.ValueString()))
    }
    if !data.SubscriberNotificationStatusMessage.IsNull() && !data.SubscriberNotificationStatusMessage.IsUnknown() {
        filters["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
        filterNames = append(filterNames, "subscriber_notification_status_message = "+fmt.Sprintf("%q", data.SubscriberNotificationStatusMessage.ValueString()))
    }
    if !data.IsPrivate.IsNull() && !data.IsPrivate.IsUnknown() {
        filters["isPrivate"] = data.IsPrivate.ValueBool()
        filterNames = append(filterNames, "is_private = "+fmt.Sprintf("%t", data.IsPrivate.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident episode up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident episode up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "title": true,
        "description": true,
        "episodeNumber": true,
        "episodeNumberWithPrefix": true,
        "currentIncidentStateId": true,
        "incidentSeverityId": true,
        "rootCause": true,
        "lastIncidentAddedAt": true,
        "resolvedAt": true,
        "allIncidentsResolvedAt": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "onCallDutyPolicies": true,
        "isOnCallPolicyExecuted": true,
        "incidentCount": true,
        "titleTemplate": true,
        "descriptionTemplate": true,
        "isManuallyCreated": true,
        "labels": true,
        "createdByUserId": true,
        "isOwnerNotifiedOfEpisodeCreation": true,
        "groupingKey": true,
        "incidentGroupingRuleId": true,
        "remediationNotes": true,
        "postmortemNote": true,
        "postUpdatesToWorkspaceChannels": true,
        "isVisibleOnStatusPage": true,
        "declaredAt": true,
        "shouldStatusPageSubscribersBeNotifiedOnEpisodeCreated": true,
        "subscriberNotificationStatusOnEpisodeCreated": true,
        "subscriberNotificationStatusMessage": true,
        "isPrivate": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-episode/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_episode, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident episode found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_episode: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-episode/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_episode, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_episode: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident episode matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident episode matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_episode.")
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
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    if val, ok := item["episodeNumber"].(float64); ok {
        data.EpisodeNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["episodeNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.EpisodeNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.EpisodeNumber = types.NumberNull()
        }
    } else {
        data.EpisodeNumber = types.NumberNull()
    }
    if obj, ok := item["episodeNumberWithPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EpisodeNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EpisodeNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EpisodeNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EpisodeNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.EpisodeNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := item["episodeNumberWithPrefix"].(string); ok {
        data.EpisodeNumberWithPrefix = types.StringValue(val)
    } else {
        data.EpisodeNumberWithPrefix = types.StringNull()
    }
    if obj, ok := item["currentIncidentStateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentIncidentStateId = types.StringNull()
        }
    } else if val, ok := item["currentIncidentStateId"].(string); ok {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if obj, ok := item["incidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if obj, ok := item["rootCause"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootCause = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootCause = types.StringValue(string(jsonBytes))
        } else {
            data.RootCause = types.StringNull()
        }
    } else if val, ok := item["rootCause"].(string); ok {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if obj, ok := item["lastIncidentAddedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastIncidentAddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastIncidentAddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastIncidentAddedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastIncidentAddedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastIncidentAddedAt = types.StringNull()
        }
    } else if val, ok := item["lastIncidentAddedAt"].(string); ok {
        data.LastIncidentAddedAt = types.StringValue(val)
    } else {
        data.LastIncidentAddedAt = types.StringNull()
    }
    if obj, ok := item["resolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedAt = types.StringNull()
        }
    } else if val, ok := item["resolvedAt"].(string); ok {
        data.ResolvedAt = types.StringValue(val)
    } else {
        data.ResolvedAt = types.StringNull()
    }
    if obj, ok := item["allIncidentsResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllIncidentsResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AllIncidentsResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AllIncidentsResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AllIncidentsResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AllIncidentsResolvedAt = types.StringNull()
        }
    } else if val, ok := item["allIncidentsResolvedAt"].(string); ok {
        data.AllIncidentsResolvedAt = types.StringValue(val)
    } else {
        data.AllIncidentsResolvedAt = types.StringNull()
    }
    if obj, ok := item["assignedToUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignedToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToUserId = types.StringNull()
        }
    } else if val, ok := item["assignedToUserId"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    } else {
        data.AssignedToUserId = types.StringNull()
    }
    if obj, ok := item["assignedToTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignedToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToTeamId = types.StringNull()
        }
    } else if val, ok := item["assignedToTeamId"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    } else {
        data.AssignedToTeamId = types.StringNull()
    }
    if val, ok := item["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OnCallDutyPolicies = types.SetNull(types.StringType)
    }
    if val, ok := item["isOnCallPolicyExecuted"].(bool); ok {
        data.IsOnCallPolicyExecuted = types.BoolValue(val)
    } else {
        data.IsOnCallPolicyExecuted = types.BoolNull()
    }
    if val, ok := item["incidentCount"].(float64); ok {
        data.IncidentCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentCount = types.NumberNull()
        }
    } else {
        data.IncidentCount = types.NumberNull()
    }
    if obj, ok := item["titleTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TitleTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TitleTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TitleTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.TitleTemplate = types.StringNull()
        }
    } else if val, ok := item["titleTemplate"].(string); ok {
        data.TitleTemplate = types.StringValue(val)
    } else {
        data.TitleTemplate = types.StringNull()
    }
    if obj, ok := item["descriptionTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DescriptionTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DescriptionTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DescriptionTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionTemplate = types.StringNull()
        }
    } else if val, ok := item["descriptionTemplate"].(string); ok {
        data.DescriptionTemplate = types.StringValue(val)
    } else {
        data.DescriptionTemplate = types.StringNull()
    }
    if val, ok := item["isManuallyCreated"].(bool); ok {
        data.IsManuallyCreated = types.BoolValue(val)
    } else {
        data.IsManuallyCreated = types.BoolNull()
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
    if val, ok := item["isOwnerNotifiedOfEpisodeCreation"].(bool); ok {
        data.IsOwnerNotifiedOfEpisodeCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfEpisodeCreation = types.BoolNull()
    }
    if obj, ok := item["groupingKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GroupingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GroupingKey = types.StringValue(string(jsonBytes))
        } else {
            data.GroupingKey = types.StringNull()
        }
    } else if val, ok := item["groupingKey"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    } else {
        data.GroupingKey = types.StringNull()
    }
    if obj, ok := item["incidentGroupingRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentGroupingRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentGroupingRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentGroupingRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentGroupingRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentGroupingRuleId = types.StringNull()
        }
    } else if val, ok := item["incidentGroupingRuleId"].(string); ok {
        data.IncidentGroupingRuleId = types.StringValue(val)
    } else {
        data.IncidentGroupingRuleId = types.StringNull()
    }
    if obj, ok := item["remediationNotes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationNotes = types.StringNull()
        }
    } else if val, ok := item["remediationNotes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if obj, ok := item["postmortemNote"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostmortemNote = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemNote = types.StringNull()
        }
    } else if val, ok := item["postmortemNote"].(string); ok {
        data.PostmortemNote = types.StringValue(val)
    } else {
        data.PostmortemNote = types.StringNull()
    }
    if obj, ok := item["postUpdatesToWorkspaceChannels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostUpdatesToWorkspaceChannels = types.StringValue(string(jsonBytes))
        } else {
            data.PostUpdatesToWorkspaceChannels = types.StringNull()
        }
    } else if val, ok := item["postUpdatesToWorkspaceChannels"].(string); ok {
        data.PostUpdatesToWorkspaceChannels = types.StringValue(val)
    } else {
        data.PostUpdatesToWorkspaceChannels = types.StringNull()
    }
    if val, ok := item["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if obj, ok := item["declaredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeclaredAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeclaredAt = types.StringNull()
        }
    } else if val, ok := item["declaredAt"].(string); ok {
        data.DeclaredAt = types.StringValue(val)
    } else {
        data.DeclaredAt = types.StringNull()
    }
    if val, ok := item["shouldStatusPageSubscribersBeNotifiedOnEpisodeCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated = types.BoolValue(val)
    } else {
        data.ShouldStatusPageSubscribersBeNotifiedOnEpisodeCreated = types.BoolNull()
    }
    if obj, ok := item["subscriberNotificationStatusOnEpisodeCreated"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnEpisodeCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusOnEpisodeCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusOnEpisodeCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusOnEpisodeCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnEpisodeCreated = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusOnEpisodeCreated"].(string); ok {
        data.SubscriberNotificationStatusOnEpisodeCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnEpisodeCreated = types.StringNull()
    }
    if obj, ok := item["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := item["subscriberNotificationStatusMessage"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if val, ok := item["isPrivate"].(bool); ok {
        data.IsPrivate = types.BoolValue(val)
    } else {
        data.IsPrivate = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
