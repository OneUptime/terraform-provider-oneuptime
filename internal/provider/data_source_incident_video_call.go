package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IncidentVideoCallDataSource{}

func NewIncidentVideoCallDataSource() datasource.DataSource {
    return &IncidentVideoCallDataSource{}
}

// IncidentVideoCallDataSource defines the data source implementation.
type IncidentVideoCallDataSource struct {
    client *Client
}

// IncidentVideoCallDataSourceModel describes the data source data model.
type IncidentVideoCallDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    ProviderValue types.String `tfsdk:"provider_value"`
    VideoCallConnectionId types.String `tfsdk:"video_call_connection_id"`
    Title types.String `tfsdk:"title"`
    JoinUrl types.String `tfsdk:"join_url"`
    ExternalMeetingId types.String `tfsdk:"external_meeting_id"`
    WorkspaceNotificationRuleId types.String `tfsdk:"workspace_notification_rule_id"`
    Transcript types.String `tfsdk:"transcript"`
    TranscriptSource types.String `tfsdk:"transcript_source"`
    TranscriptInfo types.String `tfsdk:"transcript_info"`
    TranscriptAddedAt types.String `tfsdk:"transcript_added_at"`
    TranscriptAddedByUserId types.String `tfsdk:"transcript_added_by_user_id"`
    CallSummary types.String `tfsdk:"call_summary"`
    CallSummaryStatus types.String `tfsdk:"call_summary_status"`
    CallSummaryStatusMessage types.String `tfsdk:"call_summary_status_message"`
    CallSummaryUpdatedAt types.String `tfsdk:"call_summary_updated_at"`
    RootCauseSuggestionStatus types.String `tfsdk:"root_cause_suggestion_status"`
    RootCauseSuggestionReviewedAt types.String `tfsdk:"root_cause_suggestion_reviewed_at"`
    RootCauseSuggestionReviewedByUserId types.String `tfsdk:"root_cause_suggestion_reviewed_by_user_id"`
    TranscriptFetchStatus types.String `tfsdk:"transcript_fetch_status"`
    TranscriptFetchReason types.String `tfsdk:"transcript_fetch_reason"`
    TranscriptFetchCheckedAt types.String `tfsdk:"transcript_fetch_checked_at"`
    TranscriptFetchNextAttemptAt types.String `tfsdk:"transcript_fetch_next_attempt_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentVideoCallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_video_call"
}

func (d *IncidentVideoCallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Video calls for an incident: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link. Look up an existing incident video call by `id`, or by any of its other arguments (`call_summary_status`, `call_summary_status_message`, `created_by_user_id`, ...): each one set must match, and exactly one incident video call may match them all.",

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
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident this call is for. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Where the call is held: Zoom, GoogleMeet, MicrosoftTeams, SlackHuddle, or CustomLink for a link a person provided. Taken from the connection when one is given.",
                Optional: true,
                Computed: true,
            },
            "video_call_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Zoom, Google Meet, Microsoft Teams or meeting link connection to start the call with. Leave it out to add a link of your own, or to start the huddle of the incident's Slack channel. The ID of a `oneuptime_video_call_connection`.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "What the call is called. A meeting a provider creates is named for the incident.",
                Optional: true,
                Computed: true,
            },
            "join_url": schema.StringAttribute{
                MarkdownDescription: "The link responders open to join the call. Set by OneUptime for a call started with a connection or a Slack huddle; required, as an https link, when you add a link of your own.",
                Optional: true,
                Computed: true,
            },
            "external_meeting_id": schema.StringAttribute{
                MarkdownDescription: "The provider's own id for the meeting: a Zoom meeting id, a Google Meet space name, a Microsoft Teams online meeting id or a Slack channel id.",
                Optional: true,
                Computed: true,
            },
            "workspace_notification_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Slack or Microsoft Teams notification rule that started this call. The ID of a `oneuptime_workspace_notification_rule`.",
                Optional: true,
                Computed: true,
            },
            "transcript": schema.StringAttribute{
                MarkdownDescription: "What was said on the call, one remark a line: how far into the call (HH:MM:SS, empty when unknown), a tab, who spoke (empty when unknown), a tab, and the words. Read from Zoom or Google Meet, or added by a person, through the call's transcript routes.",
                Optional: true,
                Computed: true,
            },
            "transcript_source": schema.StringAttribute{
                MarkdownDescription: "Where the transcript came from: Provider when OneUptime read it from Zoom or Google Meet, Upload when a person added it.",
                Optional: true,
                Computed: true,
            },
            "transcript_info": schema.StringAttribute{
                MarkdownDescription: "How long the transcript is - remarks, speakers, characters and how many seconds it spans - and the name of the file it was added from, so the call can be described without reading the transcript. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "transcript_added_at": schema.StringAttribute{
                MarkdownDescription: "When the transcript was added to the call.",
                Computed: true,
            },
            "transcript_added_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who added the transcript. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "call_summary": schema.StringAttribute{
                MarkdownDescription: "What OneUptime AI read in the transcript: the root cause the call points at (a suggestion, or that the call did not discuss it enough), what was found, tried, decided and left to do, and open questions - each with the transcript lines it rests on. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "call_summary_status": schema.StringAttribute{
                MarkdownDescription: "Where OneUptime AI's reading of the transcript stands: Pending while it reads, Ready, Failed (Call Summary Status Message says why), or Skipped when AI is off for the project.",
                Optional: true,
                Computed: true,
            },
            "call_summary_status_message": schema.StringAttribute{
                MarkdownDescription: "Why OneUptime AI did not read the transcript, or could not finish reading it.",
                Optional: true,
                Computed: true,
            },
            "call_summary_updated_at": schema.StringAttribute{
                MarkdownDescription: "When the Call Summary Status last changed: when OneUptime AI started reading the transcript, or finished.",
                Computed: true,
            },
            "root_cause_suggestion_status": schema.StringAttribute{
                MarkdownDescription: "What a person did with the root cause OneUptime AI suggested from the call: Suggested while it waits, Accepted once used as the incident's root cause, Dismissed. Empty when nothing was suggested.",
                Optional: true,
                Computed: true,
            },
            "root_cause_suggestion_reviewed_at": schema.StringAttribute{
                MarkdownDescription: "When a person accepted or dismissed the root cause suggested from the call.",
                Computed: true,
            },
            "root_cause_suggestion_reviewed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who accepted or dismissed the root cause suggested from the call. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "transcript_fetch_status": schema.StringAttribute{
                MarkdownDescription: "Whether OneUptime gets the call's transcript from Zoom or Google Meet by itself, as it does for the calls it starts while Record and transcribe incident calls is on: Waiting while it asks the provider (again at Transcript Fetch Next Attempt At), Stopped once it no longer asks (Transcript Fetch Reason says why). Empty for a call it does not ask about, and once it has the transcript.",
                Optional: true,
                Computed: true,
            },
            "transcript_fetch_reason": schema.StringAttribute{
                MarkdownDescription: "Why the provider had no transcript to give when OneUptime last asked: NotReady (still preparing it), NotTranscribed (nothing recorded or transcribed yet), NoMeeting, Expired, PermissionMissing, ConnectionProblem or NotSupported. Empty when the provider did not answer.",
                Optional: true,
                Computed: true,
            },
            "transcript_fetch_checked_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last asked the provider for the call's transcript.",
                Computed: true,
            },
            "transcript_fetch_next_attempt_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime next asks the provider for the call's transcript, while the Transcript Fetch Status is Waiting.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentVideoCallDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentVideoCallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentVideoCallDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.VideoCallConnectionId.IsNull() && !data.VideoCallConnectionId.IsUnknown() {
        filters["videoCallConnectionId"] = data.VideoCallConnectionId.ValueString()
        filterNames = append(filterNames, "video_call_connection_id = "+fmt.Sprintf("%q", data.VideoCallConnectionId.ValueString()))
    }
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.JoinUrl.IsNull() && !data.JoinUrl.IsUnknown() {
        filters["joinUrl"] = data.JoinUrl.ValueString()
        filterNames = append(filterNames, "join_url = "+fmt.Sprintf("%q", data.JoinUrl.ValueString()))
    }
    if !data.ExternalMeetingId.IsNull() && !data.ExternalMeetingId.IsUnknown() {
        filters["externalMeetingId"] = data.ExternalMeetingId.ValueString()
        filterNames = append(filterNames, "external_meeting_id = "+fmt.Sprintf("%q", data.ExternalMeetingId.ValueString()))
    }
    if !data.WorkspaceNotificationRuleId.IsNull() && !data.WorkspaceNotificationRuleId.IsUnknown() {
        filters["workspaceNotificationRuleId"] = data.WorkspaceNotificationRuleId.ValueString()
        filterNames = append(filterNames, "workspace_notification_rule_id = "+fmt.Sprintf("%q", data.WorkspaceNotificationRuleId.ValueString()))
    }
    if !data.Transcript.IsNull() && !data.Transcript.IsUnknown() {
        filters["transcript"] = data.Transcript.ValueString()
        filterNames = append(filterNames, "transcript = "+fmt.Sprintf("%q", data.Transcript.ValueString()))
    }
    if !data.TranscriptSource.IsNull() && !data.TranscriptSource.IsUnknown() {
        filters["transcriptSource"] = data.TranscriptSource.ValueString()
        filterNames = append(filterNames, "transcript_source = "+fmt.Sprintf("%q", data.TranscriptSource.ValueString()))
    }
    if !data.TranscriptAddedByUserId.IsNull() && !data.TranscriptAddedByUserId.IsUnknown() {
        filters["transcriptAddedByUserId"] = data.TranscriptAddedByUserId.ValueString()
        filterNames = append(filterNames, "transcript_added_by_user_id = "+fmt.Sprintf("%q", data.TranscriptAddedByUserId.ValueString()))
    }
    if !data.CallSummaryStatus.IsNull() && !data.CallSummaryStatus.IsUnknown() {
        filters["callSummaryStatus"] = data.CallSummaryStatus.ValueString()
        filterNames = append(filterNames, "call_summary_status = "+fmt.Sprintf("%q", data.CallSummaryStatus.ValueString()))
    }
    if !data.CallSummaryStatusMessage.IsNull() && !data.CallSummaryStatusMessage.IsUnknown() {
        filters["callSummaryStatusMessage"] = data.CallSummaryStatusMessage.ValueString()
        filterNames = append(filterNames, "call_summary_status_message = "+fmt.Sprintf("%q", data.CallSummaryStatusMessage.ValueString()))
    }
    if !data.RootCauseSuggestionStatus.IsNull() && !data.RootCauseSuggestionStatus.IsUnknown() {
        filters["rootCauseSuggestionStatus"] = data.RootCauseSuggestionStatus.ValueString()
        filterNames = append(filterNames, "root_cause_suggestion_status = "+fmt.Sprintf("%q", data.RootCauseSuggestionStatus.ValueString()))
    }
    if !data.RootCauseSuggestionReviewedByUserId.IsNull() && !data.RootCauseSuggestionReviewedByUserId.IsUnknown() {
        filters["rootCauseSuggestionReviewedByUserId"] = data.RootCauseSuggestionReviewedByUserId.ValueString()
        filterNames = append(filterNames, "root_cause_suggestion_reviewed_by_user_id = "+fmt.Sprintf("%q", data.RootCauseSuggestionReviewedByUserId.ValueString()))
    }
    if !data.TranscriptFetchStatus.IsNull() && !data.TranscriptFetchStatus.IsUnknown() {
        filters["transcriptFetchStatus"] = data.TranscriptFetchStatus.ValueString()
        filterNames = append(filterNames, "transcript_fetch_status = "+fmt.Sprintf("%q", data.TranscriptFetchStatus.ValueString()))
    }
    if !data.TranscriptFetchReason.IsNull() && !data.TranscriptFetchReason.IsUnknown() {
        filters["transcriptFetchReason"] = data.TranscriptFetchReason.ValueString()
        filterNames = append(filterNames, "transcript_fetch_reason = "+fmt.Sprintf("%q", data.TranscriptFetchReason.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident video call up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident video call up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incidentId": true,
        "provider": true,
        "videoCallConnectionId": true,
        "title": true,
        "joinUrl": true,
        "externalMeetingId": true,
        "workspaceNotificationRuleId": true,
        "transcript": true,
        "transcriptSource": true,
        "transcriptInfo": true,
        "transcriptAddedAt": true,
        "transcriptAddedByUserId": true,
        "callSummary": true,
        "callSummaryStatus": true,
        "callSummaryStatusMessage": true,
        "callSummaryUpdatedAt": true,
        "rootCauseSuggestionStatus": true,
        "rootCauseSuggestionReviewedAt": true,
        "rootCauseSuggestionReviewedByUserId": true,
        "transcriptFetchStatus": true,
        "transcriptFetchReason": true,
        "transcriptFetchCheckedAt": true,
        "transcriptFetchNextAttemptAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-video-call/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_video_call, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident video call found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_video_call: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-video-call/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_video_call, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_video_call: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident video call matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident video call matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_video_call.")
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
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
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
    if obj, ok := item["videoCallConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VideoCallConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VideoCallConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VideoCallConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VideoCallConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.VideoCallConnectionId = types.StringNull()
        }
    } else if val, ok := item["videoCallConnectionId"].(string); ok {
        data.VideoCallConnectionId = types.StringValue(val)
    } else {
        data.VideoCallConnectionId = types.StringNull()
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
    if obj, ok := item["joinUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.JoinUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.JoinUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.JoinUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.JoinUrl = types.StringValue(string(jsonBytes))
        } else {
            data.JoinUrl = types.StringNull()
        }
    } else if val, ok := item["joinUrl"].(string); ok {
        data.JoinUrl = types.StringValue(val)
    } else {
        data.JoinUrl = types.StringNull()
    }
    if obj, ok := item["externalMeetingId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExternalMeetingId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExternalMeetingId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExternalMeetingId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExternalMeetingId = types.StringValue(string(jsonBytes))
        } else {
            data.ExternalMeetingId = types.StringNull()
        }
    } else if val, ok := item["externalMeetingId"].(string); ok {
        data.ExternalMeetingId = types.StringValue(val)
    } else {
        data.ExternalMeetingId = types.StringNull()
    }
    if obj, ok := item["workspaceNotificationRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkspaceNotificationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkspaceNotificationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkspaceNotificationRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkspaceNotificationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkspaceNotificationRuleId = types.StringNull()
        }
    } else if val, ok := item["workspaceNotificationRuleId"].(string); ok {
        data.WorkspaceNotificationRuleId = types.StringValue(val)
    } else {
        data.WorkspaceNotificationRuleId = types.StringNull()
    }
    if obj, ok := item["transcript"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Transcript = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Transcript = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Transcript = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Transcript = types.StringValue(string(jsonBytes))
        } else {
            data.Transcript = types.StringNull()
        }
    } else if val, ok := item["transcript"].(string); ok {
        data.Transcript = types.StringValue(val)
    } else {
        data.Transcript = types.StringNull()
    }
    if obj, ok := item["transcriptSource"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptSource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptSource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptSource = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptSource = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptSource = types.StringNull()
        }
    } else if val, ok := item["transcriptSource"].(string); ok {
        data.TranscriptSource = types.StringValue(val)
    } else {
        data.TranscriptSource = types.StringNull()
    }
    if obj, ok := item["transcriptInfo"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptInfo = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptInfo = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptInfo = types.StringNull()
        }
    } else if val, ok := item["transcriptInfo"].(string); ok {
        data.TranscriptInfo = types.StringValue(val)
    } else {
        data.TranscriptInfo = types.StringNull()
    }
    if obj, ok := item["transcriptAddedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptAddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptAddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptAddedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptAddedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptAddedAt = types.StringNull()
        }
    } else if val, ok := item["transcriptAddedAt"].(string); ok {
        data.TranscriptAddedAt = types.StringValue(val)
    } else {
        data.TranscriptAddedAt = types.StringNull()
    }
    if obj, ok := item["transcriptAddedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptAddedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptAddedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptAddedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptAddedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptAddedByUserId = types.StringNull()
        }
    } else if val, ok := item["transcriptAddedByUserId"].(string); ok {
        data.TranscriptAddedByUserId = types.StringValue(val)
    } else {
        data.TranscriptAddedByUserId = types.StringNull()
    }
    if obj, ok := item["callSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallSummary = types.StringValue(string(jsonBytes))
        } else {
            data.CallSummary = types.StringNull()
        }
    } else if val, ok := item["callSummary"].(string); ok {
        data.CallSummary = types.StringValue(val)
    } else {
        data.CallSummary = types.StringNull()
    }
    if obj, ok := item["callSummaryStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSummaryStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallSummaryStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallSummaryStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallSummaryStatus = types.StringValue(string(jsonBytes))
        } else {
            data.CallSummaryStatus = types.StringNull()
        }
    } else if val, ok := item["callSummaryStatus"].(string); ok {
        data.CallSummaryStatus = types.StringValue(val)
    } else {
        data.CallSummaryStatus = types.StringNull()
    }
    if obj, ok := item["callSummaryStatusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSummaryStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallSummaryStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallSummaryStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallSummaryStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.CallSummaryStatusMessage = types.StringNull()
        }
    } else if val, ok := item["callSummaryStatusMessage"].(string); ok {
        data.CallSummaryStatusMessage = types.StringValue(val)
    } else {
        data.CallSummaryStatusMessage = types.StringNull()
    }
    if obj, ok := item["callSummaryUpdatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSummaryUpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallSummaryUpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallSummaryUpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallSummaryUpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CallSummaryUpdatedAt = types.StringNull()
        }
    } else if val, ok := item["callSummaryUpdatedAt"].(string); ok {
        data.CallSummaryUpdatedAt = types.StringValue(val)
    } else {
        data.CallSummaryUpdatedAt = types.StringNull()
    }
    if obj, ok := item["rootCauseSuggestionStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCauseSuggestionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootCauseSuggestionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootCauseSuggestionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootCauseSuggestionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.RootCauseSuggestionStatus = types.StringNull()
        }
    } else if val, ok := item["rootCauseSuggestionStatus"].(string); ok {
        data.RootCauseSuggestionStatus = types.StringValue(val)
    } else {
        data.RootCauseSuggestionStatus = types.StringNull()
    }
    if obj, ok := item["rootCauseSuggestionReviewedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCauseSuggestionReviewedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootCauseSuggestionReviewedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootCauseSuggestionReviewedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootCauseSuggestionReviewedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RootCauseSuggestionReviewedAt = types.StringNull()
        }
    } else if val, ok := item["rootCauseSuggestionReviewedAt"].(string); ok {
        data.RootCauseSuggestionReviewedAt = types.StringValue(val)
    } else {
        data.RootCauseSuggestionReviewedAt = types.StringNull()
    }
    if obj, ok := item["rootCauseSuggestionReviewedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCauseSuggestionReviewedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RootCauseSuggestionReviewedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RootCauseSuggestionReviewedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RootCauseSuggestionReviewedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.RootCauseSuggestionReviewedByUserId = types.StringNull()
        }
    } else if val, ok := item["rootCauseSuggestionReviewedByUserId"].(string); ok {
        data.RootCauseSuggestionReviewedByUserId = types.StringValue(val)
    } else {
        data.RootCauseSuggestionReviewedByUserId = types.StringNull()
    }
    if obj, ok := item["transcriptFetchStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptFetchStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptFetchStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptFetchStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptFetchStatus = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptFetchStatus = types.StringNull()
        }
    } else if val, ok := item["transcriptFetchStatus"].(string); ok {
        data.TranscriptFetchStatus = types.StringValue(val)
    } else {
        data.TranscriptFetchStatus = types.StringNull()
    }
    if obj, ok := item["transcriptFetchReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptFetchReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptFetchReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptFetchReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptFetchReason = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptFetchReason = types.StringNull()
        }
    } else if val, ok := item["transcriptFetchReason"].(string); ok {
        data.TranscriptFetchReason = types.StringValue(val)
    } else {
        data.TranscriptFetchReason = types.StringNull()
    }
    if obj, ok := item["transcriptFetchCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptFetchCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptFetchCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptFetchCheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptFetchCheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptFetchCheckedAt = types.StringNull()
        }
    } else if val, ok := item["transcriptFetchCheckedAt"].(string); ok {
        data.TranscriptFetchCheckedAt = types.StringValue(val)
    } else {
        data.TranscriptFetchCheckedAt = types.StringNull()
    }
    if obj, ok := item["transcriptFetchNextAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TranscriptFetchNextAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TranscriptFetchNextAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TranscriptFetchNextAttemptAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TranscriptFetchNextAttemptAt = types.StringValue(string(jsonBytes))
        } else {
            data.TranscriptFetchNextAttemptAt = types.StringNull()
        }
    } else if val, ok := item["transcriptFetchNextAttemptAt"].(string); ok {
        data.TranscriptFetchNextAttemptAt = types.StringValue(val)
    } else {
        data.TranscriptFetchNextAttemptAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
