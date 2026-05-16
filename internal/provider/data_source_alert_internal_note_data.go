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
var _ datasource.DataSource = &AlertInternalNoteDataDataSource{}

func NewAlertInternalNoteDataDataSource() datasource.DataSource {
    return &AlertInternalNoteDataDataSource{}
}

// AlertInternalNoteDataDataSource defines the data source implementation.
type AlertInternalNoteDataDataSource struct {
    client *Client
}

// AlertInternalNoteDataDataSourceModel describes the data source data model.
type AlertInternalNoteDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AlertId types.String `tfsdk:"alert_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Note types.String `tfsdk:"note"`
    Attachments types.Set `tfsdk:"attachments"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    PostedFromSlackMessageId types.String `tfsdk:"posted_from_slack_message_id"`
}

func (d *AlertInternalNoteDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_internal_note_data"
}

func (d *AlertInternalNoteDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_internal_note_data data source",

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
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "note": schema.StringAttribute{
                MarkdownDescription: "Notes in markdown. Permissions - Create: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Create Alert Internal Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Internal Note, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Edit Alert Internal Note]",
                Computed: true,
            },
            "attachments": schema.SetAttribute{
                MarkdownDescription: "Files attached to this note. Permissions - Create: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Create Alert Internal Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Internal Note, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Edit Alert Internal Note]",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of this resource ownership?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Internal Note, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "posted_from_slack_message_id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.. Permissions - Create: [Project Owner, Project Admin, Project Member, Alert Admin, Alert Member, Create Alert Internal Note], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Internal Note, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *AlertInternalNoteDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertInternalNoteDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertInternalNoteDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-internal-note" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_internal_note_data, got error: %s", err))
        return
    }

    var alertInternalNoteDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertInternalNoteDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_internal_note_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertInternalNoteDataResponse["data"].(map[string]interface{}); ok {
        alertInternalNoteDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertInternalNoteDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertInternalNoteDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["note"].(string); ok {
        data.Note = types.StringValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["attachments"].([]interface{}); ok {
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
    if val, ok := alertInternalNoteDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := alertInternalNoteDataResponse["posted_from_slack_message_id"].(string); ok {
        data.PostedFromSlackMessageId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
