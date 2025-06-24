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
var _ datasource.DataSource = &StatusPageDataDataSource{}

func NewStatusPageDataDataSource() datasource.DataSource {
    return &StatusPageDataDataSource{}
}

// StatusPageDataDataSource defines the data source implementation.
type StatusPageDataDataSource struct {
    client *Client
}

// StatusPageDataDataSourceModel describes the data source data model.
type StatusPageDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    PageTitle types.String `tfsdk:"page_title"`
    PageDescription types.String `tfsdk:"page_description"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Labels types.List `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    FaviconFileId types.String `tfsdk:"favicon_file_id"`
    LogoFileId types.String `tfsdk:"logo_file_id"`
    CoverImageFileId types.String `tfsdk:"cover_image_file_id"`
    HeaderHtml types.String `tfsdk:"header_html"`
    FooterHtml types.String `tfsdk:"footer_html"`
    CustomCss types.String `tfsdk:"custom_css"`
    CustomJavaScript types.String `tfsdk:"custom_java_script"`
    IsPublicStatusPage types.Bool `tfsdk:"is_public_status_page"`
    ShowIncidentLabelsOnStatusPage types.Bool `tfsdk:"show_incident_labels_on_status_page"`
    ShowScheduledEventLabelsOnStatusPage types.Bool `tfsdk:"show_scheduled_event_labels_on_status_page"`
    EnableSubscribers types.Bool `tfsdk:"enable_subscribers"`
    EnableEmailSubscribers types.Bool `tfsdk:"enable_email_subscribers"`
    AllowSubscribersToChooseResources types.Bool `tfsdk:"allow_subscribers_to_choose_resources"`
    AllowSubscribersToChooseEventTypes types.Bool `tfsdk:"allow_subscribers_to_choose_event_types"`
    EnableSmsSubscribers types.Bool `tfsdk:"enable_sms_subscribers"`
    EnableSlackSubscribers types.Bool `tfsdk:"enable_slack_subscribers"`
    CopyrightText types.String `tfsdk:"copyright_text"`
    CustomFields types.String `tfsdk:"custom_fields"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    SmtpConfigId types.String `tfsdk:"smtp_config_id"`
    CallSmsConfigId types.String `tfsdk:"call_sms_config_id"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    ShowIncidentHistoryInDays types.Number `tfsdk:"show_incident_history_in_days"`
    ShowAnnouncementHistoryInDays types.Number `tfsdk:"show_announcement_history_in_days"`
    ShowScheduledEventHistoryInDays types.Number `tfsdk:"show_scheduled_event_history_in_days"`
    OverviewPageDescription types.String `tfsdk:"overview_page_description"`
    HidePoweredByOneUptimeBranding types.Bool `tfsdk:"hide_powered_by_one_uptime_branding"`
    DefaultBarColor types.String `tfsdk:"default_bar_color"`
    DowntimeMonitorStatuses types.List `tfsdk:"downtime_monitor_statuses"`
    SubscriberTimezones types.String `tfsdk:"subscriber_timezones"`
    IsReportEnabled types.Bool `tfsdk:"is_report_enabled"`
    ReportStartDateTime types.String `tfsdk:"report_start_date_time"`
    ReportRecurringInterval types.String `tfsdk:"report_recurring_interval"`
    SendNextReportBy types.String `tfsdk:"send_next_report_by"`
    ReportDataInDays types.Number `tfsdk:"report_data_in_days"`
    ShowOverallUptimePercentOnStatusPage types.Bool `tfsdk:"show_overall_uptime_percent_on_status_page"`
    OverallUptimePercentPrecision types.String `tfsdk:"overall_uptime_percent_precision"`
    SubscriberEmailNotificationFooterText types.String `tfsdk:"subscriber_email_notification_footer_text"`
    ShowIncidentsOnStatusPage types.Bool `tfsdk:"show_incidents_on_status_page"`
    ShowAnnouncementsOnStatusPage types.Bool `tfsdk:"show_announcements_on_status_page"`
    ShowScheduledMaintenanceEventsOnStatusPage types.Bool `tfsdk:"show_scheduled_maintenance_events_on_status_page"`
    ShowSubscriberPageOnStatusPage types.Bool `tfsdk:"show_subscriber_page_on_status_page"`
    IpWhitelist types.String `tfsdk:"ip_whitelist"`
}

func (d *StatusPageDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_data"
}

func (d *StatusPageDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_data data source",

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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "favicon_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "logo_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "cover_image_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "header_html": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "footer_html": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "custom_css": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "enable_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page, Public], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "smtp_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "show_incident_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "default_bar_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "downtime_monitor_statuses": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
                ElementType: types.StringType,
            },
            "subscriber_timezones": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "report_start_date_time": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "report_recurring_interval": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "send_next_report_by": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_data, got error: %s", err))
        return
    }

    var statusPageDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageDataResponse["data"].(map[string]interface{}); ok {
        statusPageDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["page_title"].(string); ok {
        data.PageTitle = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["page_description"].(string); ok {
        data.PageDescription = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }
    if val, ok := statusPageDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["favicon_file_id"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["logo_file_id"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["cover_image_file_id"].(string); ok {
        data.CoverImageFileId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["header_html"].(string); ok {
        data.HeaderHtml = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["footer_html"].(string); ok {
        data.FooterHtml = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["custom_css"].(string); ok {
        data.CustomCss = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["custom_java_script"].(string); ok {
        data.CustomJavaScript = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["is_public_status_page"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_incident_labels_on_status_page"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_scheduled_event_labels_on_status_page"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["enable_subscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["enable_email_subscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["allow_subscribers_to_choose_resources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["allow_subscribers_to_choose_event_types"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["enable_sms_subscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["enable_slack_subscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["copyright_text"].(string); ok {
        data.CopyrightText = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["require_sso_for_login"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["smtp_config_id"].(string); ok {
        data.SmtpConfigId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["call_sms_config_id"].(string); ok {
        data.CallSmsConfigId = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["is_owner_notified_of_resource_creation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_incident_history_in_days"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["show_announcement_history_in_days"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["show_scheduled_event_history_in_days"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["overview_page_description"].(string); ok {
        data.OverviewPageDescription = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["hide_powered_by_one_uptime_branding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["default_bar_color"].(string); ok {
        data.DefaultBarColor = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["downtime_monitor_statuses"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.DowntimeMonitorStatuses = listValue
    }
    if val, ok := statusPageDataResponse["subscriber_timezones"].(string); ok {
        data.SubscriberTimezones = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["is_report_enabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["report_start_date_time"].(string); ok {
        data.ReportStartDateTime = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["report_recurring_interval"].(string); ok {
        data.ReportRecurringInterval = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["send_next_report_by"].(string); ok {
        data.SendNextReportBy = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["report_data_in_days"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["show_overall_uptime_percent_on_status_page"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["overall_uptime_percent_precision"].(string); ok {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["subscriber_email_notification_footer_text"].(string); ok {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["show_incidents_on_status_page"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_announcements_on_status_page"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_scheduled_maintenance_events_on_status_page"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_subscriber_page_on_status_page"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["ip_whitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
