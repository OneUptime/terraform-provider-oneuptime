package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AlertFeedDataDataSource{}

func NewAlertFeedDataDataSource() datasource.DataSource {
    return &AlertFeedDataDataSource{}
}

// AlertFeedDataDataSource defines the data source implementation.
type AlertFeedDataDataSource struct {
    client *Client
}

// AlertFeedDataDataSourceModel describes the data source data model.
type AlertFeedDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AlertId types.String `tfsdk:"alert_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    FeedInfoInMarkdown types.String `tfsdk:"feed_info_in_markdown"`
    MoreInformationInMarkdown types.String `tfsdk:"more_information_in_markdown"`
    AlertFeedEventType types.String `tfsdk:"alert_feed_event_type"`
    DisplayColor types.String `tfsdk:"display_color"`
    UserId types.String `tfsdk:"user_id"`
    PostedAt types.String `tfsdk:"posted_at"`
}

func (d *AlertFeedDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_feed_data"
}

func (d *AlertFeedDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_feed_data data source",

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
                MarkdownDescription: "Version",
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
            "feed_info_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert Feed], Read: [Project Owner, Project Admin, Project Member, Read Alert Feed], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "more_information_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert Feed], Read: [Project Owner, Project Admin, Project Member, Read Alert Feed], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "alert_feed_event_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert Feed], Read: [Project Owner, Project Admin, Project Member, Read Alert Feed], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "display_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *AlertFeedDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertFeedDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertFeedDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert-feed" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_feed_data, got error: %s", err))
        return
    }

    var alertFeedDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertFeedDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_feed_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertFeedDataResponse["data"].(map[string]interface{}); ok {
        alertFeedDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertFeedDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertFeedDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["alert_id"].(string); ok {
        data.AlertId = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["feed_info_in_markdown"].(string); ok {
        data.FeedInfoInMarkdown = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["more_information_in_markdown"].(string); ok {
        data.MoreInformationInMarkdown = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["alert_feed_event_type"].(string); ok {
        data.AlertFeedEventType = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["display_color"].(string); ok {
        data.DisplayColor = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := alertFeedDataResponse["posted_at"].(string); ok {
        data.PostedAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
