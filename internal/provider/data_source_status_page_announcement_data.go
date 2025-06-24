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
var _ datasource.DataSource = &StatusPageAnnouncementDataDataSource{}

func NewStatusPageAnnouncementDataDataSource() datasource.DataSource {
    return &StatusPageAnnouncementDataDataSource{}
}

// StatusPageAnnouncementDataDataSource defines the data source implementation.
type StatusPageAnnouncementDataDataSource struct {
    client *Client
}

// StatusPageAnnouncementDataDataSourceModel describes the data source data model.
type StatusPageAnnouncementDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPages types.List `tfsdk:"status_pages"`
    Title types.String `tfsdk:"title"`
    ShowAnnouncementAt types.String `tfsdk:"show_announcement_at"`
    EndAnnouncementAt types.String `tfsdk:"end_announcement_at"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsStatusPageSubscribersNotified types.Bool `tfsdk:"is_status_page_subscribers_notified"`
    ShouldStatusPageSubscribersBeNotified types.Bool `tfsdk:"should_status_page_subscribers_be_notified"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (d *StatusPageAnnouncementDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_announcement_data"
}

func (d *StatusPageAnnouncementDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_announcement_data data source",

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
            "status_pages": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement]",
                Computed: true,
                ElementType: types.StringType,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement]",
                Computed: true,
            },
            "show_announcement_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "end_announcement_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_status_page_subscribers_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "should_status_page_subscribers_be_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageAnnouncementDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageAnnouncementDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageAnnouncementDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-announcement" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_announcement_data, got error: %s", err))
        return
    }

    var statusPageAnnouncementDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageAnnouncementDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_announcement_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageAnnouncementDataResponse["data"].(map[string]interface{}); ok {
        statusPageAnnouncementDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageAnnouncementDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageAnnouncementDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["status_pages"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.StatusPages = listValue
    }
    if val, ok := statusPageAnnouncementDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["show_announcement_at"].(string); ok {
        data.ShowAnnouncementAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["end_announcement_at"].(string); ok {
        data.EndAnnouncementAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["is_status_page_subscribers_notified"].(bool); ok {
        data.IsStatusPageSubscribersNotified = types.BoolValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["should_status_page_subscribers_be_notified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    }
    if val, ok := statusPageAnnouncementDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
