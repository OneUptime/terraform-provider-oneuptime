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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentVideoCallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_video_call"
}

func (d *IncidentVideoCallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Video calls for an incident: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link. Look up an existing incident video call by `id`, or by any of its other arguments (`created_by_user_id`, `external_meeting_id`, `incident_id`, ...): each one set must match, and exactly one incident video call may match them all.",

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
