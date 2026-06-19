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
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    FaviconFileId types.String `tfsdk:"favicon_file_id"`
    LogoFileId types.String `tfsdk:"logo_file_id"`
    CoverImageFileId types.String `tfsdk:"cover_image_file_id"`
    HeaderHtml types.String `tfsdk:"header_html"`
    FooterHtml types.String `tfsdk:"footer_html"`
    CustomCss types.String `tfsdk:"custom_css"`
    CustomJavaScript types.String `tfsdk:"custom_java_script"`
    IsPublicStatusPage types.Bool `tfsdk:"is_public_status_page"`
    EnableMasterPassword types.Bool `tfsdk:"enable_master_password"`
    MasterPassword types.String `tfsdk:"master_password"`
    ShowIncidentLabelsOnStatusPage types.Bool `tfsdk:"show_incident_labels_on_status_page"`
    ShowScheduledEventLabelsOnStatusPage types.Bool `tfsdk:"show_scheduled_event_labels_on_status_page"`
    EnableSubscribers types.Bool `tfsdk:"enable_subscribers"`
    EnableEmailSubscribers types.Bool `tfsdk:"enable_email_subscribers"`
    AllowSubscribersToChooseResources types.Bool `tfsdk:"allow_subscribers_to_choose_resources"`
    AllowSubscribersToChooseEventTypes types.Bool `tfsdk:"allow_subscribers_to_choose_event_types"`
    EnableSmsSubscribers types.Bool `tfsdk:"enable_sms_subscribers"`
    EnableSlackSubscribers types.Bool `tfsdk:"enable_slack_subscribers"`
    EnableMicrosoftTeamsSubscribers types.Bool `tfsdk:"enable_microsoft_teams_subscribers"`
    EnableWebhookSubscribers types.Bool `tfsdk:"enable_webhook_subscribers"`
    CopyrightText types.String `tfsdk:"copyright_text"`
    LogoAltText types.String `tfsdk:"logo_alt_text"`
    CoverImageAltText types.String `tfsdk:"cover_image_alt_text"`
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
    DowntimeMonitorStatuses types.Set `tfsdk:"downtime_monitor_statuses"`
    SubscriberTimezones types.String `tfsdk:"subscriber_timezones"`
    IsReportEnabled types.Bool `tfsdk:"is_report_enabled"`
    ReportStartDateTime types.String `tfsdk:"report_start_date_time"`
    ReportRecurringInterval types.String `tfsdk:"report_recurring_interval"`
    SendNextReportBy types.String `tfsdk:"send_next_report_by"`
    ReportDataInDays types.Number `tfsdk:"report_data_in_days"`
    ShowOverallUptimePercentOnStatusPage types.Bool `tfsdk:"show_overall_uptime_percent_on_status_page"`
    OverallUptimePercentPrecision types.String `tfsdk:"overall_uptime_percent_precision"`
    SubscriberEmailNotificationFooterText types.String `tfsdk:"subscriber_email_notification_footer_text"`
    EnableCustomSubscriberEmailNotificationFooterText types.Bool `tfsdk:"enable_custom_subscriber_email_notification_footer_text"`
    ShowIncidentsOnStatusPage types.Bool `tfsdk:"show_incidents_on_status_page"`
    ShowAnnouncementsOnStatusPage types.Bool `tfsdk:"show_announcements_on_status_page"`
    ShowEpisodesOnStatusPage types.Bool `tfsdk:"show_episodes_on_status_page"`
    ShowEpisodeHistoryInDays types.Number `tfsdk:"show_episode_history_in_days"`
    ShowEpisodeLabelsOnStatusPage types.Bool `tfsdk:"show_episode_labels_on_status_page"`
    ShowScheduledMaintenanceEventsOnStatusPage types.Bool `tfsdk:"show_scheduled_maintenance_events_on_status_page"`
    ShowSubscriberPageOnStatusPage types.Bool `tfsdk:"show_subscriber_page_on_status_page"`
    IpWhitelist types.String `tfsdk:"ip_whitelist"`
    EnableEmbeddedOverallStatus types.Bool `tfsdk:"enable_embedded_overall_status"`
    ShowUptimeHistoryInDays types.Number `tfsdk:"show_uptime_history_in_days"`
    EmbeddedOverallStatusToken types.String `tfsdk:"embedded_overall_status_token"`
    DefaultLanguage types.String `tfsdk:"default_language"`
    EnabledLanguages types.String `tfsdk:"enabled_languages"`
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
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Title of your Status Page. This is used for SEO.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Description of your Status Page. This is used for SEO.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
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
                MarkdownDescription: "Status Page Custom HTML Header. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "footer_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Footer. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "custom_css": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom CSS Header. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom JavaScript. This runs when the status page is loaded.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Is this status page public?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_master_password": schema.BoolAttribute{
                MarkdownDescription: "Require visitors to enter a master password before viewing a private status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "master_password": schema.StringAttribute{
                MarkdownDescription: "Password required to unlock a private status page. This value is stored as a secure hash.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Labels on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Event Labels on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can email subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which resources to subscribe to?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can SMS subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Slack subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_microsoft_teams_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Microsoft Teams subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_webhook_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Webhook subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Copyright Text. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "logo_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the logo image, read by screen readers for accessibility.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "cover_image_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Should SSO be required to login to Private Status Page. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page, Public], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
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
                MarkdownDescription: "Are owners notified of when this resource is created?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "show_incident_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of incident history should be shown on the status page (in days)?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of announcement history should be shown on the status page (in days)?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of scheduled event history should be shown on the status page (in days)?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Overview Page description for your status page. This is a markdown field.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Hide Powered By OneUptime Branding?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "default_bar_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitors statuses that are considered as \"down\" for this status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
                ElementType: types.StringType,
            },
            "subscriber_timezones": schema.StringAttribute{
                MarkdownDescription: "Timezones of subscribers to this status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Report Enabled for this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "report_start_date_time": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "report_recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often would you like to send the report?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "send_next_report_by": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of data should be included in the report?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Overall Uptime Percent on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Overall Precision of uptime percent for this status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Text to send to subscribers in the footer of the email.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_custom_subscriber_email_notification_footer_text": schema.BoolAttribute{
                MarkdownDescription: "Enable custom footer text in subscriber email notifications.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incidents on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Announcements on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_episodes_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Episodes on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_episode_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of episode history to show on the status page. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_episode_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Episode Labels on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Maintenance Events on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Subscriber Page on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist for this Status Page. One IP per line. Only used if the status page is private.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enable_embedded_overall_status": schema.BoolAttribute{
                MarkdownDescription: "Enable embedded overall status badge that can be displayed on external websites?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "show_uptime_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of uptime history should be shown on the status page? Maximum is 90 days.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "embedded_overall_status_token": schema.StringAttribute{
                MarkdownDescription: "Security token required to access the embedded overall status badge. This token must be provided in the URL.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "default_language": schema.StringAttribute{
                MarkdownDescription: "Default language that the status page is shown in when a visitor arrives for the first time.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
                Computed: true,
            },
            "enabled_languages": schema.StringAttribute{
                MarkdownDescription: "Languages offered in the footer language switcher. Leave empty to offer all supported languages.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page]",
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
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
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
    if val, ok := statusPageDataResponse["enable_master_password"].(bool); ok {
        data.EnableMasterPassword = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["master_password"].(string); ok {
        data.MasterPassword = types.StringValue(val)
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
    if val, ok := statusPageDataResponse["enable_microsoft_teams_subscribers"].(bool); ok {
        data.EnableMicrosoftTeamsSubscribers = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["enable_webhook_subscribers"].(bool); ok {
        data.EnableWebhookSubscribers = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["copyright_text"].(string); ok {
        data.CopyrightText = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["logo_alt_text"].(string); ok {
        data.LogoAltText = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["cover_image_alt_text"].(string); ok {
        data.CoverImageAltText = types.StringValue(val)
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
        setValue, _ := types.SetValue(types.StringType, elements)
        data.DowntimeMonitorStatuses = setValue
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
    if val, ok := statusPageDataResponse["enable_custom_subscriber_email_notification_footer_text"].(bool); ok {
        data.EnableCustomSubscriberEmailNotificationFooterText = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_incidents_on_status_page"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_announcements_on_status_page"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_episodes_on_status_page"].(bool); ok {
        data.ShowEpisodesOnStatusPage = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_episode_history_in_days"].(float64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["show_episode_labels_on_status_page"].(bool); ok {
        data.ShowEpisodeLabelsOnStatusPage = types.BoolValue(val)
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
    if val, ok := statusPageDataResponse["enable_embedded_overall_status"].(bool); ok {
        data.EnableEmbeddedOverallStatus = types.BoolValue(val)
    }
    if val, ok := statusPageDataResponse["show_uptime_history_in_days"].(float64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageDataResponse["embedded_overall_status_token"].(string); ok {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["default_language"].(string); ok {
        data.DefaultLanguage = types.StringValue(val)
    }
    if val, ok := statusPageDataResponse["enabled_languages"].(string); ok {
        data.EnabledLanguages = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
