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
var _ datasource.DataSource = &IncidentEpisodePublicNoteDataDataSource{}

func NewIncidentEpisodePublicNoteDataDataSource() datasource.DataSource {
    return &IncidentEpisodePublicNoteDataDataSource{}
}

// IncidentEpisodePublicNoteDataDataSource defines the data source implementation.
type IncidentEpisodePublicNoteDataDataSource struct {
    client *Client
}

// IncidentEpisodePublicNoteDataDataSourceModel describes the data source data model.
type IncidentEpisodePublicNoteDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentEpisodeId types.String `tfsdk:"incident_episode_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Note types.String `tfsdk:"note"`
    Attachments types.Set `tfsdk:"attachments"`
    SubscriberNotificationStatusOnNoteCreated types.String `tfsdk:"subscriber_notification_status_on_note_created"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    ShouldStatusPageSubscribersBeNotifiedOnNoteCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_note_created"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    PostedAt types.String `tfsdk:"posted_at"`
    PostedFromSlackMessageId types.String `tfsdk:"posted_from_slack_message_id"`
}

func (d *IncidentEpisodePublicNoteDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_episode_public_note_data"
}

func (d *IncidentEpisodePublicNoteDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident_episode_public_note_data data source",

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
            "incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "note": schema.StringAttribute{
                MarkdownDescription: "Notes in markdown. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Manager, Create Incident Episode Public Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Manager, Edit Incident Episode Public Note]",
                Computed: true,
            },
            "attachments": schema.SetAttribute{
                MarkdownDescription: "Files attached to this note. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Manager, Create Incident Episode Public Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Manager, Edit Incident Episode Public Note]",
                Computed: true,
                ElementType: types.StringType,
            },
            "subscriber_notification_status_on_note_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this note. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Manager, Create Incident Episode Public Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Manager, Edit Incident Episode Public Note]",
                Computed: true,
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Manager, Create Incident Episode Public Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Incident Manager, Edit Incident Episode Public Note]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified_on_note_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this note?. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Manager, Create Incident Episode Public Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of this resource ownership?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "posted_from_slack_message_id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.. Permissions - Create: [Project Owner, Project Admin, Project Member, Incident Manager, Create Incident Episode Public Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Incident Manager, Read Incident Episode Public Note, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *IncidentEpisodePublicNoteDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentEpisodePublicNoteDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentEpisodePublicNoteDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incident-episode-public-note" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_episode_public_note_data, got error: %s", err))
        return
    }

    var incidentEpisodePublicNoteDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incidentEpisodePublicNoteDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident_episode_public_note_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incidentEpisodePublicNoteDataResponse["data"].(map[string]interface{}); ok {
        incidentEpisodePublicNoteDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incidentEpisodePublicNoteDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["incident_episode_id"].(string); ok {
        data.IncidentEpisodeId = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["note"].(string); ok {
        data.Note = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["attachments"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Attachments = setValue
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["subscriber_notification_status_on_note_created"].(string); ok {
        data.SubscriberNotificationStatusOnNoteCreated = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["subscriber_notification_status_message"].(string); ok {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["should_status_page_subscribers_be_notified_on_note_created"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnNoteCreated = types.BoolValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["posted_at"].(string); ok {
        data.PostedAt = types.StringValue(val)
    }
    if val, ok := incidentEpisodePublicNoteDataResponse["posted_from_slack_message_id"].(string); ok {
        data.PostedFromSlackMessageId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
