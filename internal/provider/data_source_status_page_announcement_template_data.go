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
var _ datasource.DataSource = &StatusPageAnnouncementTemplateDataDataSource{}

func NewStatusPageAnnouncementTemplateDataDataSource() datasource.DataSource {
    return &StatusPageAnnouncementTemplateDataDataSource{}
}

// StatusPageAnnouncementTemplateDataDataSource defines the data source implementation.
type StatusPageAnnouncementTemplateDataDataSource struct {
    client *Client
}

// StatusPageAnnouncementTemplateDataDataSourceModel describes the data source data model.
type StatusPageAnnouncementTemplateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    StatusPages types.List `tfsdk:"status_pages"`
    ShouldStatusPageSubscribersBeNotified types.Bool `tfsdk:"should_status_page_subscribers_be_notified"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *StatusPageAnnouncementTemplateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_announcement_template_data"
}

func (d *StatusPageAnnouncementTemplateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_announcement_template_data data source",

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
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement Template], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement Template], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement Template]",
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement Template], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement Template], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement Template]",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement Template], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement Template], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement Template]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement Template], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement Template], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement Template]",
                Computed: true,
            },
            "status_pages": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement Template], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement Template], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "should_status_page_subscribers_be_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Announcement Template], Read: [Project Owner, Project Admin, Project Member, Read Status Page Announcement Template], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Announcement Template]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageAnnouncementTemplateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageAnnouncementTemplateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageAnnouncementTemplateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-announcement-template" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_announcement_template_data, got error: %s", err))
        return
    }

    var statusPageAnnouncementTemplateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageAnnouncementTemplateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_announcement_template_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageAnnouncementTemplateDataResponse["data"].(map[string]interface{}); ok {
        statusPageAnnouncementTemplateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageAnnouncementTemplateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["template_name"].(string); ok {
        data.TemplateName = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["template_description"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["status_pages"].([]interface{}); ok {
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
    if val, ok := statusPageAnnouncementTemplateDataResponse["should_status_page_subscribers_be_notified"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotified = types.BoolValue(val)
    }
    if val, ok := statusPageAnnouncementTemplateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
