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
var _ datasource.DataSource = &StatusPageDataSource{}

func NewStatusPageDataSource() datasource.DataSource {
    return &StatusPageDataSource{}
}

// StatusPageDataSource defines the data source implementation.
type StatusPageDataSource struct {
    client *Client
}

// StatusPageDataSourceModel describes the data source data model.
type StatusPageDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    PageTitle types.String `tfsdk:"page_title"`
    PageDescription types.String `tfsdk:"page_description"`
    EnableSearchEngineIndexing types.Bool `tfsdk:"enable_search_engine_indexing"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    FaviconFileId types.String `tfsdk:"favicon_file_id"`
    LogoFileId types.String `tfsdk:"logo_file_id"`
    CoverImageFileId types.String `tfsdk:"cover_image_file_id"`
    HeaderHtml types.String `tfsdk:"header_html"`
    FooterHtml types.String `tfsdk:"footer_html"`
    CustomCss types.String `tfsdk:"custom_css"`
    CustomJavaScript types.String `tfsdk:"custom_java_script"`
    IsPublicStatusPage types.Bool `tfsdk:"is_public_status_page"`
    EnableMcpServer types.Bool `tfsdk:"enable_mcp_server"`
    EnableMasterPassword types.Bool `tfsdk:"enable_master_password"`
    MasterPassword types.String `tfsdk:"master_password"`
    ShowIncidentLabelsOnStatusPage types.Bool `tfsdk:"show_incident_labels_on_status_page"`
    ShowScheduledEventLabelsOnStatusPage types.Bool `tfsdk:"show_scheduled_event_labels_on_status_page"`
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
    ReportPeriodType types.String `tfsdk:"report_period_type"`
    ReportTimezone types.String `tfsdk:"report_timezone"`
    ShowOverallUptimePercentOnStatusPage types.Bool `tfsdk:"show_overall_uptime_percent_on_status_page"`
    OverallUptimePercentPrecision types.String `tfsdk:"overall_uptime_percent_precision"`
    SubscriberEmailNotificationFooterText types.String `tfsdk:"subscriber_email_notification_footer_text"`
    EnableCustomSubscriberEmailNotificationFooterText types.Bool `tfsdk:"enable_custom_subscriber_email_notification_footer_text"`
    ShowIncidentsOnStatusPage types.Bool `tfsdk:"show_incidents_on_status_page"`
    OnlyShowScopedIncidents types.Bool `tfsdk:"only_show_scoped_incidents"`
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

func (d *StatusPageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (d *StatusPageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage status pages for your project. Look up an existing status page by `id`, or by any of its other arguments (`name`, `allow_subscribers_to_choose_event_types`, `allow_subscribers_to_choose_resources`, ...): each one set must match, and exactly one status page may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Title of your Status Page. This is used for SEO.",
                Optional: true,
                Computed: true,
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Description of your Status Page. This is used for SEO.",
                Optional: true,
                Computed: true,
            },
            "enable_search_engine_indexing": schema.BoolAttribute{
                MarkdownDescription: "Should search engines like Google and Bing be allowed to index this status page? Turn this off to keep the page reachable by link but out of search results.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
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
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Archived status pages are hidden from the Status Pages list, are not served to visitors, and send nothing to their subscribers. Unarchiving puts them back online.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When this status page was archived. Empty while it is not archived.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "favicon_file_id": schema.StringAttribute{
                MarkdownDescription: "Status Page Favicon File ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "logo_file_id": schema.StringAttribute{
                MarkdownDescription: "Status Page Logo File ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "cover_image_file_id": schema.StringAttribute{
                MarkdownDescription: "Status Page Cover Image ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "header_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Header. Served only from a verified custom domain.",
                Optional: true,
                Computed: true,
            },
            "footer_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Footer. Served only from a verified custom domain.",
                Optional: true,
                Computed: true,
            },
            "custom_css": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom CSS. Served only from a verified custom domain.",
                Optional: true,
                Computed: true,
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom JavaScript. This runs when the status page is loaded from a verified custom domain.",
                Optional: true,
                Computed: true,
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Is this status page public?",
                Optional: true,
                Computed: true,
            },
            "enable_mcp_server": schema.BoolAttribute{
                MarkdownDescription: "Can AI agents read this status page over the public OneUptime MCP server? This does not affect the status page website, its RSS feed, or its public JSON API.",
                Optional: true,
                Computed: true,
            },
            "enable_master_password": schema.BoolAttribute{
                MarkdownDescription: "Require visitors to enter a master password before viewing a private status page.",
                Optional: true,
                Computed: true,
            },
            "master_password": schema.StringAttribute{
                MarkdownDescription: "Password required to unlock a private status page. This value is stored as a secure hash.",
                Computed: true,
                Sensitive: true,
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Labels on Status Page?",
                Optional: true,
                Computed: true,
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Event Labels on Status Page?",
                Optional: true,
                Computed: true,
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can email subscribers subscribe to this Status Page?",
                Optional: true,
                Computed: true,
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which resources to subscribe to?",
                Optional: true,
                Computed: true,
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?",
                Optional: true,
                Computed: true,
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can SMS subscribers subscribe to this Status Page?",
                Optional: true,
                Computed: true,
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Slack subscribers subscribe to this Status Page?",
                Optional: true,
                Computed: true,
            },
            "enable_microsoft_teams_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Microsoft Teams subscribers subscribe to this Status Page?",
                Optional: true,
                Computed: true,
            },
            "enable_webhook_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Webhook subscribers subscribe to this Status Page?",
                Optional: true,
                Computed: true,
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Copyright Text.",
                Optional: true,
                Computed: true,
            },
            "logo_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the logo image, read by screen readers for accessibility.",
                Optional: true,
                Computed: true,
            },
            "cover_image_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative.",
                Optional: true,
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Should SSO be required to login to Private Status Page.",
                Optional: true,
                Computed: true,
            },
            "smtp_config_id": schema.StringAttribute{
                MarkdownDescription: "ID of your SMTP Config Resource which is used to send email to subscribers.",
                Optional: true,
                Computed: true,
            },
            "call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "ID of your Call/SMS Config Resource which is used to send SMS to subscribers.",
                Optional: true,
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?",
                Optional: true,
                Computed: true,
            },
            "show_incident_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of incident history should be shown on the status page (in days)?",
                Optional: true,
                Computed: true,
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of announcement history should be shown on the status page (in days)?",
                Optional: true,
                Computed: true,
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of scheduled event history should be shown on the status page (in days)?",
                Optional: true,
                Computed: true,
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Overview Page description for your status page. This is a markdown field.",
                Optional: true,
                Computed: true,
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Hide Powered By OneUptime Branding?",
                Optional: true,
                Computed: true,
            },
            "default_bar_color": schema.StringAttribute{
                MarkdownDescription: "Default color of the bar on the overview page.",
                Computed: true,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitors statuses that are considered as \"down\" for this status page. IDs of `oneuptime_monitor_status` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "subscriber_timezones": schema.StringAttribute{
                MarkdownDescription: "Timezones of subscribers to this status page. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this status page's email subscribers get reports. Turned on without a schedule, reports go out on the 1st of every month at 09:00 in the report timezone, each covering the calendar month before it.",
                Optional: true,
                Computed: true,
            },
            "report_start_date_time": schema.StringAttribute{
                MarkdownDescription: "When the first report goes out. Every later one follows it by the recurring interval, at the same time of day. Left out when reports are turned on, it is 09:00 in the report timezone at the start of the next period of the interval: the next 1st of the month for a monthly schedule (the default), the next Monday for a weekly one, the next day for a daily one and the next 1 January for a yearly one. An hourly schedule starts at the next full hour.",
                Computed: true,
            },
            "report_recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often a report goes out. Left out when reports are turned on, it is every month. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "send_next_report_by": schema.StringAttribute{
                MarkdownDescription: "When the next report goes out. The server works it out from the schedule.",
                Computed: true,
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of data should be included in the report?",
                Optional: true,
                Computed: true,
            },
            "report_period_type": schema.StringAttribute{
                MarkdownDescription: "Should the report cover a rolling number of days, or the previous whole calendar period?",
                Optional: true,
                Computed: true,
            },
            "report_timezone": schema.StringAttribute{
                MarkdownDescription: "The timezone report periods and send times are resolved in. A monthly report in this timezone runs from the 1st at 00:00 to the last day at 23:59.",
                Optional: true,
                Computed: true,
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Overall Uptime Percent on Status Page?",
                Optional: true,
                Computed: true,
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Overall Precision of uptime percent for this status page.",
                Optional: true,
                Computed: true,
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Text to send to subscribers in the footer of the email.",
                Optional: true,
                Computed: true,
            },
            "enable_custom_subscriber_email_notification_footer_text": schema.BoolAttribute{
                MarkdownDescription: "Enable custom footer text in subscriber email notifications.",
                Optional: true,
                Computed: true,
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incidents on Status Page?",
                Optional: true,
                Computed: true,
            },
            "only_show_scoped_incidents": schema.BoolAttribute{
                MarkdownDescription: "When on, this status page shows and notifies its subscribers about only the incidents limited to it. Incidents that are not limited to any status page never reach it.",
                Optional: true,
                Computed: true,
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Announcements on Status Page?",
                Optional: true,
                Computed: true,
            },
            "show_episodes_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Episodes on Status Page?",
                Optional: true,
                Computed: true,
            },
            "show_episode_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of episode history to show on the status page.",
                Optional: true,
                Computed: true,
            },
            "show_episode_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Episode Labels on Status Page?",
                Optional: true,
                Computed: true,
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Maintenance Events on Status Page?",
                Optional: true,
                Computed: true,
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Subscriber Page on Status Page?",
                Optional: true,
                Computed: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist for this Status Page. One IP per line. Only used if the status page is private.",
                Optional: true,
                Computed: true,
            },
            "enable_embedded_overall_status": schema.BoolAttribute{
                MarkdownDescription: "Enable embedded overall status badge that can be displayed on external websites?",
                Optional: true,
                Computed: true,
            },
            "show_uptime_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of uptime history should be shown on the status page? Maximum is 90 days.",
                Optional: true,
                Computed: true,
            },
            "embedded_overall_status_token": schema.StringAttribute{
                MarkdownDescription: "Security token required to access the embedded overall status badge. This token must be provided in the URL.",
                Optional: true,
                Computed: true,
            },
            "default_language": schema.StringAttribute{
                MarkdownDescription: "Default language that the status page is shown in when a visitor arrives for the first time.",
                Optional: true,
                Computed: true,
            },
            "enabled_languages": schema.StringAttribute{
                MarkdownDescription: "Languages offered in the footer language switcher. Leave empty to offer all supported languages. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.PageTitle.IsNull() && !data.PageTitle.IsUnknown() {
        filters["pageTitle"] = data.PageTitle.ValueString()
        filterNames = append(filterNames, "page_title = "+fmt.Sprintf("%q", data.PageTitle.ValueString()))
    }
    if !data.PageDescription.IsNull() && !data.PageDescription.IsUnknown() {
        filters["pageDescription"] = data.PageDescription.ValueString()
        filterNames = append(filterNames, "page_description = "+fmt.Sprintf("%q", data.PageDescription.ValueString()))
    }
    if !data.EnableSearchEngineIndexing.IsNull() && !data.EnableSearchEngineIndexing.IsUnknown() {
        filters["enableSearchEngineIndexing"] = data.EnableSearchEngineIndexing.ValueBool()
        filterNames = append(filterNames, "enable_search_engine_indexing = "+fmt.Sprintf("%t", data.EnableSearchEngineIndexing.ValueBool()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.FaviconFileId.IsNull() && !data.FaviconFileId.IsUnknown() {
        filters["faviconFileId"] = data.FaviconFileId.ValueString()
        filterNames = append(filterNames, "favicon_file_id = "+fmt.Sprintf("%q", data.FaviconFileId.ValueString()))
    }
    if !data.LogoFileId.IsNull() && !data.LogoFileId.IsUnknown() {
        filters["logoFileId"] = data.LogoFileId.ValueString()
        filterNames = append(filterNames, "logo_file_id = "+fmt.Sprintf("%q", data.LogoFileId.ValueString()))
    }
    if !data.CoverImageFileId.IsNull() && !data.CoverImageFileId.IsUnknown() {
        filters["coverImageFileId"] = data.CoverImageFileId.ValueString()
        filterNames = append(filterNames, "cover_image_file_id = "+fmt.Sprintf("%q", data.CoverImageFileId.ValueString()))
    }
    if !data.HeaderHtml.IsNull() && !data.HeaderHtml.IsUnknown() {
        filters["headerHTML"] = data.HeaderHtml.ValueString()
        filterNames = append(filterNames, "header_html = "+fmt.Sprintf("%q", data.HeaderHtml.ValueString()))
    }
    if !data.FooterHtml.IsNull() && !data.FooterHtml.IsUnknown() {
        filters["footerHTML"] = data.FooterHtml.ValueString()
        filterNames = append(filterNames, "footer_html = "+fmt.Sprintf("%q", data.FooterHtml.ValueString()))
    }
    if !data.CustomCss.IsNull() && !data.CustomCss.IsUnknown() {
        filters["customCSS"] = data.CustomCss.ValueString()
        filterNames = append(filterNames, "custom_css = "+fmt.Sprintf("%q", data.CustomCss.ValueString()))
    }
    if !data.CustomJavaScript.IsNull() && !data.CustomJavaScript.IsUnknown() {
        filters["customJavaScript"] = data.CustomJavaScript.ValueString()
        filterNames = append(filterNames, "custom_java_script = "+fmt.Sprintf("%q", data.CustomJavaScript.ValueString()))
    }
    if !data.IsPublicStatusPage.IsNull() && !data.IsPublicStatusPage.IsUnknown() {
        filters["isPublicStatusPage"] = data.IsPublicStatusPage.ValueBool()
        filterNames = append(filterNames, "is_public_status_page = "+fmt.Sprintf("%t", data.IsPublicStatusPage.ValueBool()))
    }
    if !data.EnableMcpServer.IsNull() && !data.EnableMcpServer.IsUnknown() {
        filters["enableMcpServer"] = data.EnableMcpServer.ValueBool()
        filterNames = append(filterNames, "enable_mcp_server = "+fmt.Sprintf("%t", data.EnableMcpServer.ValueBool()))
    }
    if !data.EnableMasterPassword.IsNull() && !data.EnableMasterPassword.IsUnknown() {
        filters["enableMasterPassword"] = data.EnableMasterPassword.ValueBool()
        filterNames = append(filterNames, "enable_master_password = "+fmt.Sprintf("%t", data.EnableMasterPassword.ValueBool()))
    }
    if !data.ShowIncidentLabelsOnStatusPage.IsNull() && !data.ShowIncidentLabelsOnStatusPage.IsUnknown() {
        filters["showIncidentLabelsOnStatusPage"] = data.ShowIncidentLabelsOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_incident_labels_on_status_page = "+fmt.Sprintf("%t", data.ShowIncidentLabelsOnStatusPage.ValueBool()))
    }
    if !data.ShowScheduledEventLabelsOnStatusPage.IsNull() && !data.ShowScheduledEventLabelsOnStatusPage.IsUnknown() {
        filters["showScheduledEventLabelsOnStatusPage"] = data.ShowScheduledEventLabelsOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_scheduled_event_labels_on_status_page = "+fmt.Sprintf("%t", data.ShowScheduledEventLabelsOnStatusPage.ValueBool()))
    }
    if !data.EnableEmailSubscribers.IsNull() && !data.EnableEmailSubscribers.IsUnknown() {
        filters["enableEmailSubscribers"] = data.EnableEmailSubscribers.ValueBool()
        filterNames = append(filterNames, "enable_email_subscribers = "+fmt.Sprintf("%t", data.EnableEmailSubscribers.ValueBool()))
    }
    if !data.AllowSubscribersToChooseResources.IsNull() && !data.AllowSubscribersToChooseResources.IsUnknown() {
        filters["allowSubscribersToChooseResources"] = data.AllowSubscribersToChooseResources.ValueBool()
        filterNames = append(filterNames, "allow_subscribers_to_choose_resources = "+fmt.Sprintf("%t", data.AllowSubscribersToChooseResources.ValueBool()))
    }
    if !data.AllowSubscribersToChooseEventTypes.IsNull() && !data.AllowSubscribersToChooseEventTypes.IsUnknown() {
        filters["allowSubscribersToChooseEventTypes"] = data.AllowSubscribersToChooseEventTypes.ValueBool()
        filterNames = append(filterNames, "allow_subscribers_to_choose_event_types = "+fmt.Sprintf("%t", data.AllowSubscribersToChooseEventTypes.ValueBool()))
    }
    if !data.EnableSmsSubscribers.IsNull() && !data.EnableSmsSubscribers.IsUnknown() {
        filters["enableSmsSubscribers"] = data.EnableSmsSubscribers.ValueBool()
        filterNames = append(filterNames, "enable_sms_subscribers = "+fmt.Sprintf("%t", data.EnableSmsSubscribers.ValueBool()))
    }
    if !data.EnableSlackSubscribers.IsNull() && !data.EnableSlackSubscribers.IsUnknown() {
        filters["enableSlackSubscribers"] = data.EnableSlackSubscribers.ValueBool()
        filterNames = append(filterNames, "enable_slack_subscribers = "+fmt.Sprintf("%t", data.EnableSlackSubscribers.ValueBool()))
    }
    if !data.EnableMicrosoftTeamsSubscribers.IsNull() && !data.EnableMicrosoftTeamsSubscribers.IsUnknown() {
        filters["enableMicrosoftTeamsSubscribers"] = data.EnableMicrosoftTeamsSubscribers.ValueBool()
        filterNames = append(filterNames, "enable_microsoft_teams_subscribers = "+fmt.Sprintf("%t", data.EnableMicrosoftTeamsSubscribers.ValueBool()))
    }
    if !data.EnableWebhookSubscribers.IsNull() && !data.EnableWebhookSubscribers.IsUnknown() {
        filters["enableWebhookSubscribers"] = data.EnableWebhookSubscribers.ValueBool()
        filterNames = append(filterNames, "enable_webhook_subscribers = "+fmt.Sprintf("%t", data.EnableWebhookSubscribers.ValueBool()))
    }
    if !data.CopyrightText.IsNull() && !data.CopyrightText.IsUnknown() {
        filters["copyrightText"] = data.CopyrightText.ValueString()
        filterNames = append(filterNames, "copyright_text = "+fmt.Sprintf("%q", data.CopyrightText.ValueString()))
    }
    if !data.LogoAltText.IsNull() && !data.LogoAltText.IsUnknown() {
        filters["logoAltText"] = data.LogoAltText.ValueString()
        filterNames = append(filterNames, "logo_alt_text = "+fmt.Sprintf("%q", data.LogoAltText.ValueString()))
    }
    if !data.CoverImageAltText.IsNull() && !data.CoverImageAltText.IsUnknown() {
        filters["coverImageAltText"] = data.CoverImageAltText.ValueString()
        filterNames = append(filterNames, "cover_image_alt_text = "+fmt.Sprintf("%q", data.CoverImageAltText.ValueString()))
    }
    if !data.RequireSsoForLogin.IsNull() && !data.RequireSsoForLogin.IsUnknown() {
        filters["requireSsoForLogin"] = data.RequireSsoForLogin.ValueBool()
        filterNames = append(filterNames, "require_sso_for_login = "+fmt.Sprintf("%t", data.RequireSsoForLogin.ValueBool()))
    }
    if !data.SmtpConfigId.IsNull() && !data.SmtpConfigId.IsUnknown() {
        filters["smtpConfigId"] = data.SmtpConfigId.ValueString()
        filterNames = append(filterNames, "smtp_config_id = "+fmt.Sprintf("%q", data.SmtpConfigId.ValueString()))
    }
    if !data.CallSmsConfigId.IsNull() && !data.CallSmsConfigId.IsUnknown() {
        filters["callSmsConfigId"] = data.CallSmsConfigId.ValueString()
        filterNames = append(filterNames, "call_sms_config_id = "+fmt.Sprintf("%q", data.CallSmsConfigId.ValueString()))
    }
    if !data.IsOwnerNotifiedOfResourceCreation.IsNull() && !data.IsOwnerNotifiedOfResourceCreation.IsUnknown() {
        filters["isOwnerNotifiedOfResourceCreation"] = data.IsOwnerNotifiedOfResourceCreation.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_resource_creation = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfResourceCreation.ValueBool()))
    }
    if !data.ShowIncidentHistoryInDays.IsNull() && !data.ShowIncidentHistoryInDays.IsUnknown() {
        filters["showIncidentHistoryInDays"] = lookupNumber(data.ShowIncidentHistoryInDays)
        filterNames = append(filterNames, "show_incident_history_in_days = "+data.ShowIncidentHistoryInDays.ValueBigFloat().String())
    }
    if !data.ShowAnnouncementHistoryInDays.IsNull() && !data.ShowAnnouncementHistoryInDays.IsUnknown() {
        filters["showAnnouncementHistoryInDays"] = lookupNumber(data.ShowAnnouncementHistoryInDays)
        filterNames = append(filterNames, "show_announcement_history_in_days = "+data.ShowAnnouncementHistoryInDays.ValueBigFloat().String())
    }
    if !data.ShowScheduledEventHistoryInDays.IsNull() && !data.ShowScheduledEventHistoryInDays.IsUnknown() {
        filters["showScheduledEventHistoryInDays"] = lookupNumber(data.ShowScheduledEventHistoryInDays)
        filterNames = append(filterNames, "show_scheduled_event_history_in_days = "+data.ShowScheduledEventHistoryInDays.ValueBigFloat().String())
    }
    if !data.OverviewPageDescription.IsNull() && !data.OverviewPageDescription.IsUnknown() {
        filters["overviewPageDescription"] = data.OverviewPageDescription.ValueString()
        filterNames = append(filterNames, "overview_page_description = "+fmt.Sprintf("%q", data.OverviewPageDescription.ValueString()))
    }
    if !data.HidePoweredByOneUptimeBranding.IsNull() && !data.HidePoweredByOneUptimeBranding.IsUnknown() {
        filters["hidePoweredByOneUptimeBranding"] = data.HidePoweredByOneUptimeBranding.ValueBool()
        filterNames = append(filterNames, "hide_powered_by_one_uptime_branding = "+fmt.Sprintf("%t", data.HidePoweredByOneUptimeBranding.ValueBool()))
    }
    if !data.IsReportEnabled.IsNull() && !data.IsReportEnabled.IsUnknown() {
        filters["isReportEnabled"] = data.IsReportEnabled.ValueBool()
        filterNames = append(filterNames, "is_report_enabled = "+fmt.Sprintf("%t", data.IsReportEnabled.ValueBool()))
    }
    if !data.ReportDataInDays.IsNull() && !data.ReportDataInDays.IsUnknown() {
        filters["reportDataInDays"] = lookupNumber(data.ReportDataInDays)
        filterNames = append(filterNames, "report_data_in_days = "+data.ReportDataInDays.ValueBigFloat().String())
    }
    if !data.ReportPeriodType.IsNull() && !data.ReportPeriodType.IsUnknown() {
        filters["reportPeriodType"] = data.ReportPeriodType.ValueString()
        filterNames = append(filterNames, "report_period_type = "+fmt.Sprintf("%q", data.ReportPeriodType.ValueString()))
    }
    if !data.ReportTimezone.IsNull() && !data.ReportTimezone.IsUnknown() {
        filters["reportTimezone"] = data.ReportTimezone.ValueString()
        filterNames = append(filterNames, "report_timezone = "+fmt.Sprintf("%q", data.ReportTimezone.ValueString()))
    }
    if !data.ShowOverallUptimePercentOnStatusPage.IsNull() && !data.ShowOverallUptimePercentOnStatusPage.IsUnknown() {
        filters["showOverallUptimePercentOnStatusPage"] = data.ShowOverallUptimePercentOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_overall_uptime_percent_on_status_page = "+fmt.Sprintf("%t", data.ShowOverallUptimePercentOnStatusPage.ValueBool()))
    }
    if !data.OverallUptimePercentPrecision.IsNull() && !data.OverallUptimePercentPrecision.IsUnknown() {
        filters["overallUptimePercentPrecision"] = data.OverallUptimePercentPrecision.ValueString()
        filterNames = append(filterNames, "overall_uptime_percent_precision = "+fmt.Sprintf("%q", data.OverallUptimePercentPrecision.ValueString()))
    }
    if !data.SubscriberEmailNotificationFooterText.IsNull() && !data.SubscriberEmailNotificationFooterText.IsUnknown() {
        filters["subscriberEmailNotificationFooterText"] = data.SubscriberEmailNotificationFooterText.ValueString()
        filterNames = append(filterNames, "subscriber_email_notification_footer_text = "+fmt.Sprintf("%q", data.SubscriberEmailNotificationFooterText.ValueString()))
    }
    if !data.EnableCustomSubscriberEmailNotificationFooterText.IsNull() && !data.EnableCustomSubscriberEmailNotificationFooterText.IsUnknown() {
        filters["enableCustomSubscriberEmailNotificationFooterText"] = data.EnableCustomSubscriberEmailNotificationFooterText.ValueBool()
        filterNames = append(filterNames, "enable_custom_subscriber_email_notification_footer_text = "+fmt.Sprintf("%t", data.EnableCustomSubscriberEmailNotificationFooterText.ValueBool()))
    }
    if !data.ShowIncidentsOnStatusPage.IsNull() && !data.ShowIncidentsOnStatusPage.IsUnknown() {
        filters["showIncidentsOnStatusPage"] = data.ShowIncidentsOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_incidents_on_status_page = "+fmt.Sprintf("%t", data.ShowIncidentsOnStatusPage.ValueBool()))
    }
    if !data.OnlyShowScopedIncidents.IsNull() && !data.OnlyShowScopedIncidents.IsUnknown() {
        filters["onlyShowScopedIncidents"] = data.OnlyShowScopedIncidents.ValueBool()
        filterNames = append(filterNames, "only_show_scoped_incidents = "+fmt.Sprintf("%t", data.OnlyShowScopedIncidents.ValueBool()))
    }
    if !data.ShowAnnouncementsOnStatusPage.IsNull() && !data.ShowAnnouncementsOnStatusPage.IsUnknown() {
        filters["showAnnouncementsOnStatusPage"] = data.ShowAnnouncementsOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_announcements_on_status_page = "+fmt.Sprintf("%t", data.ShowAnnouncementsOnStatusPage.ValueBool()))
    }
    if !data.ShowEpisodesOnStatusPage.IsNull() && !data.ShowEpisodesOnStatusPage.IsUnknown() {
        filters["showEpisodesOnStatusPage"] = data.ShowEpisodesOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_episodes_on_status_page = "+fmt.Sprintf("%t", data.ShowEpisodesOnStatusPage.ValueBool()))
    }
    if !data.ShowEpisodeHistoryInDays.IsNull() && !data.ShowEpisodeHistoryInDays.IsUnknown() {
        filters["showEpisodeHistoryInDays"] = lookupNumber(data.ShowEpisodeHistoryInDays)
        filterNames = append(filterNames, "show_episode_history_in_days = "+data.ShowEpisodeHistoryInDays.ValueBigFloat().String())
    }
    if !data.ShowEpisodeLabelsOnStatusPage.IsNull() && !data.ShowEpisodeLabelsOnStatusPage.IsUnknown() {
        filters["showEpisodeLabelsOnStatusPage"] = data.ShowEpisodeLabelsOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_episode_labels_on_status_page = "+fmt.Sprintf("%t", data.ShowEpisodeLabelsOnStatusPage.ValueBool()))
    }
    if !data.ShowScheduledMaintenanceEventsOnStatusPage.IsNull() && !data.ShowScheduledMaintenanceEventsOnStatusPage.IsUnknown() {
        filters["showScheduledMaintenanceEventsOnStatusPage"] = data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_scheduled_maintenance_events_on_status_page = "+fmt.Sprintf("%t", data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool()))
    }
    if !data.ShowSubscriberPageOnStatusPage.IsNull() && !data.ShowSubscriberPageOnStatusPage.IsUnknown() {
        filters["showSubscriberPageOnStatusPage"] = data.ShowSubscriberPageOnStatusPage.ValueBool()
        filterNames = append(filterNames, "show_subscriber_page_on_status_page = "+fmt.Sprintf("%t", data.ShowSubscriberPageOnStatusPage.ValueBool()))
    }
    if !data.IpWhitelist.IsNull() && !data.IpWhitelist.IsUnknown() {
        filters["ipWhitelist"] = data.IpWhitelist.ValueString()
        filterNames = append(filterNames, "ip_whitelist = "+fmt.Sprintf("%q", data.IpWhitelist.ValueString()))
    }
    if !data.EnableEmbeddedOverallStatus.IsNull() && !data.EnableEmbeddedOverallStatus.IsUnknown() {
        filters["enableEmbeddedOverallStatus"] = data.EnableEmbeddedOverallStatus.ValueBool()
        filterNames = append(filterNames, "enable_embedded_overall_status = "+fmt.Sprintf("%t", data.EnableEmbeddedOverallStatus.ValueBool()))
    }
    if !data.ShowUptimeHistoryInDays.IsNull() && !data.ShowUptimeHistoryInDays.IsUnknown() {
        filters["showUptimeHistoryInDays"] = lookupNumber(data.ShowUptimeHistoryInDays)
        filterNames = append(filterNames, "show_uptime_history_in_days = "+data.ShowUptimeHistoryInDays.ValueBigFloat().String())
    }
    if !data.EmbeddedOverallStatusToken.IsNull() && !data.EmbeddedOverallStatusToken.IsUnknown() {
        filters["embeddedOverallStatusToken"] = data.EmbeddedOverallStatusToken.ValueString()
        filterNames = append(filterNames, "embedded_overall_status_token = "+fmt.Sprintf("%q", data.EmbeddedOverallStatusToken.ValueString()))
    }
    if !data.DefaultLanguage.IsNull() && !data.DefaultLanguage.IsUnknown() {
        filters["defaultLanguage"] = data.DefaultLanguage.ValueString()
        filterNames = append(filterNames, "default_language = "+fmt.Sprintf("%q", data.DefaultLanguage.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the status page up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "pageTitle": true,
        "pageDescription": true,
        "enableSearchEngineIndexing": true,
        "description": true,
        "slug": true,
        "labels": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "faviconFileId": true,
        "logoFileId": true,
        "coverImageFileId": true,
        "headerHTML": true,
        "footerHTML": true,
        "customCSS": true,
        "customJavaScript": true,
        "isPublicStatusPage": true,
        "enableMcpServer": true,
        "enableMasterPassword": true,
        "masterPassword": true,
        "showIncidentLabelsOnStatusPage": true,
        "showScheduledEventLabelsOnStatusPage": true,
        "enableEmailSubscribers": true,
        "allowSubscribersToChooseResources": true,
        "allowSubscribersToChooseEventTypes": true,
        "enableSmsSubscribers": true,
        "enableSlackSubscribers": true,
        "enableMicrosoftTeamsSubscribers": true,
        "enableWebhookSubscribers": true,
        "copyrightText": true,
        "logoAltText": true,
        "coverImageAltText": true,
        "customFields": true,
        "requireSsoForLogin": true,
        "smtpConfigId": true,
        "callSmsConfigId": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "showIncidentHistoryInDays": true,
        "showAnnouncementHistoryInDays": true,
        "showScheduledEventHistoryInDays": true,
        "overviewPageDescription": true,
        "hidePoweredByOneUptimeBranding": true,
        "defaultBarColor": true,
        "downtimeMonitorStatuses": true,
        "subscriberTimezones": true,
        "isReportEnabled": true,
        "reportStartDateTime": true,
        "reportRecurringInterval": true,
        "sendNextReportBy": true,
        "reportDataInDays": true,
        "reportPeriodType": true,
        "reportTimezone": true,
        "showOverallUptimePercentOnStatusPage": true,
        "overallUptimePercentPrecision": true,
        "subscriberEmailNotificationFooterText": true,
        "enableCustomSubscriberEmailNotificationFooterText": true,
        "showIncidentsOnStatusPage": true,
        "onlyShowScopedIncidents": true,
        "showAnnouncementsOnStatusPage": true,
        "showEpisodesOnStatusPage": true,
        "showEpisodeHistoryInDays": true,
        "showEpisodeLabelsOnStatusPage": true,
        "showScheduledMaintenanceEventsOnStatusPage": true,
        "showSubscriberPageOnStatusPage": true,
        "ipWhitelist": true,
        "enableEmbeddedOverallStatus": true,
        "showUptimeHistoryInDays": true,
        "embeddedOverallStatusToken": true,
        "defaultLanguage": true,
        "enabledLanguages": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page.")
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
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["pageTitle"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PageTitle = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PageTitle = types.StringValue(string(jsonBytes))
        } else {
            data.PageTitle = types.StringNull()
        }
    } else if val, ok := item["pageTitle"].(string); ok {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if obj, ok := item["pageDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.PageDescription = types.StringNull()
        }
    } else if val, ok := item["pageDescription"].(string); ok {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := item["enableSearchEngineIndexing"].(bool); ok {
        data.EnableSearchEngineIndexing = types.BoolValue(val)
    } else {
        data.EnableSearchEngineIndexing = types.BoolNull()
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := item["faviconFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if val, ok := item["faviconFileId"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if obj, ok := item["logoFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LogoFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if val, ok := item["logoFileId"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if obj, ok := item["coverImageFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CoverImageFileId = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageFileId = types.StringNull()
        }
    } else if val, ok := item["coverImageFileId"].(string); ok {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if obj, ok := item["headerHTML"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HeaderHtml = types.StringValue(string(jsonBytes))
        } else {
            data.HeaderHtml = types.StringNull()
        }
    } else if val, ok := item["headerHTML"].(string); ok {
        data.HeaderHtml = types.StringValue(val)
    } else {
        data.HeaderHtml = types.StringNull()
    }
    if obj, ok := item["footerHTML"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FooterHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FooterHtml = types.StringValue(string(jsonBytes))
        } else {
            data.FooterHtml = types.StringNull()
        }
    } else if val, ok := item["footerHTML"].(string); ok {
        data.FooterHtml = types.StringValue(val)
    } else {
        data.FooterHtml = types.StringNull()
    }
    if obj, ok := item["customCSS"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomCss = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomCss = types.StringValue(string(jsonBytes))
        } else {
            data.CustomCss = types.StringNull()
        }
    } else if val, ok := item["customCSS"].(string); ok {
        data.CustomCss = types.StringValue(val)
    } else {
        data.CustomCss = types.StringNull()
    }
    if obj, ok := item["customJavaScript"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomJavaScript = types.StringValue(string(jsonBytes))
        } else {
            data.CustomJavaScript = types.StringNull()
        }
    } else if val, ok := item["customJavaScript"].(string); ok {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := item["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    } else {
        data.IsPublicStatusPage = types.BoolNull()
    }
    if val, ok := item["enableMcpServer"].(bool); ok {
        data.EnableMcpServer = types.BoolValue(val)
    } else {
        data.EnableMcpServer = types.BoolNull()
    }
    if val, ok := item["enableMasterPassword"].(bool); ok {
        data.EnableMasterPassword = types.BoolValue(val)
    } else {
        data.EnableMasterPassword = types.BoolNull()
    }
    if obj, ok := item["masterPassword"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MasterPassword = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MasterPassword = types.StringValue(string(jsonBytes))
        } else {
            data.MasterPassword = types.StringNull()
        }
    } else if val, ok := item["masterPassword"].(string); ok {
        data.MasterPassword = types.StringValue(val)
    } else {
        data.MasterPassword = types.StringNull()
    }
    if val, ok := item["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowIncidentLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := item["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := item["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    } else {
        data.EnableEmailSubscribers = types.BoolNull()
    }
    if val, ok := item["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    } else {
        data.AllowSubscribersToChooseResources = types.BoolNull()
    }
    if val, ok := item["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    } else {
        data.AllowSubscribersToChooseEventTypes = types.BoolNull()
    }
    if val, ok := item["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    } else {
        data.EnableSmsSubscribers = types.BoolNull()
    }
    if val, ok := item["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    } else {
        data.EnableSlackSubscribers = types.BoolNull()
    }
    if val, ok := item["enableMicrosoftTeamsSubscribers"].(bool); ok {
        data.EnableMicrosoftTeamsSubscribers = types.BoolValue(val)
    } else {
        data.EnableMicrosoftTeamsSubscribers = types.BoolNull()
    }
    if val, ok := item["enableWebhookSubscribers"].(bool); ok {
        data.EnableWebhookSubscribers = types.BoolValue(val)
    } else {
        data.EnableWebhookSubscribers = types.BoolNull()
    }
    if obj, ok := item["copyrightText"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CopyrightText = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CopyrightText = types.StringValue(string(jsonBytes))
        } else {
            data.CopyrightText = types.StringNull()
        }
    } else if val, ok := item["copyrightText"].(string); ok {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := item["logoAltText"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LogoAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LogoAltText = types.StringValue(string(jsonBytes))
        } else {
            data.LogoAltText = types.StringNull()
        }
    } else if val, ok := item["logoAltText"].(string); ok {
        data.LogoAltText = types.StringValue(val)
    } else {
        data.LogoAltText = types.StringNull()
    }
    if obj, ok := item["coverImageAltText"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CoverImageAltText = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageAltText = types.StringNull()
        }
    } else if val, ok := item["coverImageAltText"].(string); ok {
        data.CoverImageAltText = types.StringValue(val)
    } else {
        data.CoverImageAltText = types.StringNull()
    }
    if obj, ok := item["customFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := item["customFields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if val, ok := item["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if obj, ok := item["smtpConfigId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SmtpConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.SmtpConfigId = types.StringNull()
        }
    } else if val, ok := item["smtpConfigId"].(string); ok {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if obj, ok := item["callSmsConfigId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallSmsConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.CallSmsConfigId = types.StringNull()
        }
    } else if val, ok := item["callSmsConfigId"].(string); ok {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := item["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := item["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["showIncidentHistoryInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowIncidentHistoryInDays = types.NumberNull()
        }
    } else {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := item["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["showAnnouncementHistoryInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowAnnouncementHistoryInDays = types.NumberNull()
        }
    } else {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := item["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["showScheduledEventHistoryInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowScheduledEventHistoryInDays = types.NumberNull()
        }
    } else {
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if obj, ok := item["overviewPageDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OverviewPageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.OverviewPageDescription = types.StringNull()
        }
    } else if val, ok := item["overviewPageDescription"].(string); ok {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := item["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    } else {
        data.HidePoweredByOneUptimeBranding = types.BoolNull()
    }
    if obj, ok := item["defaultBarColor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DefaultBarColor = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = types.StringNull()
        }
    } else if val, ok := item["defaultBarColor"].(string); ok {
        data.DefaultBarColor = types.StringValue(val)
    } else {
        data.DefaultBarColor = types.StringNull()
    }
    if val, ok := item["downtimeMonitorStatuses"].([]interface{}); ok {
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
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DowntimeMonitorStatuses = types.SetNull(types.StringType)
    }
    if obj, ok := item["subscriberTimezones"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberTimezones = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = types.StringNull()
        }
    } else if val, ok := item["subscriberTimezones"].(string); ok {
        data.SubscriberTimezones = types.StringValue(val)
    } else {
        data.SubscriberTimezones = types.StringNull()
    }
    if val, ok := item["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    } else {
        data.IsReportEnabled = types.BoolNull()
    }
    if obj, ok := item["reportStartDateTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReportStartDateTime = types.StringValue(string(jsonBytes))
        } else {
            data.ReportStartDateTime = types.StringNull()
        }
    } else if val, ok := item["reportStartDateTime"].(string); ok {
        data.ReportStartDateTime = types.StringValue(val)
    } else {
        data.ReportStartDateTime = types.StringNull()
    }
    if obj, ok := item["reportRecurringInterval"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = types.StringNull()
        }
    } else if val, ok := item["reportRecurringInterval"].(string); ok {
        data.ReportRecurringInterval = types.StringValue(val)
    } else {
        data.ReportRecurringInterval = types.StringNull()
    }
    if obj, ok := item["sendNextReportBy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SendNextReportBy = types.StringValue(string(jsonBytes))
        } else {
            data.SendNextReportBy = types.StringNull()
        }
    } else if val, ok := item["sendNextReportBy"].(string); ok {
        data.SendNextReportBy = types.StringValue(val)
    } else {
        data.SendNextReportBy = types.StringNull()
    }
    if val, ok := item["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["reportDataInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReportDataInDays = types.NumberNull()
        }
    } else {
        data.ReportDataInDays = types.NumberNull()
    }
    if obj, ok := item["reportPeriodType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReportPeriodType = types.StringValue(string(jsonBytes))
        } else {
            data.ReportPeriodType = types.StringNull()
        }
    } else if val, ok := item["reportPeriodType"].(string); ok {
        data.ReportPeriodType = types.StringValue(val)
    } else {
        data.ReportPeriodType = types.StringNull()
    }
    if obj, ok := item["reportTimezone"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReportTimezone = types.StringValue(string(jsonBytes))
        } else {
            data.ReportTimezone = types.StringNull()
        }
    } else if val, ok := item["reportTimezone"].(string); ok {
        data.ReportTimezone = types.StringValue(val)
    } else {
        data.ReportTimezone = types.StringNull()
    }
    if val, ok := item["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolNull()
    }
    if obj, ok := item["overallUptimePercentPrecision"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
        } else {
            data.OverallUptimePercentPrecision = types.StringNull()
        }
    } else if val, ok := item["overallUptimePercentPrecision"].(string); ok {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if obj, ok := item["subscriberEmailNotificationFooterText"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmailNotificationFooterText = types.StringNull()
        }
    } else if val, ok := item["subscriberEmailNotificationFooterText"].(string); ok {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := item["enableCustomSubscriberEmailNotificationFooterText"].(bool); ok {
        data.EnableCustomSubscriberEmailNotificationFooterText = types.BoolValue(val)
    } else {
        data.EnableCustomSubscriberEmailNotificationFooterText = types.BoolNull()
    }
    if val, ok := item["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowIncidentsOnStatusPage = types.BoolNull()
    }
    if val, ok := item["onlyShowScopedIncidents"].(bool); ok {
        data.OnlyShowScopedIncidents = types.BoolValue(val)
    } else {
        data.OnlyShowScopedIncidents = types.BoolNull()
    }
    if val, ok := item["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowAnnouncementsOnStatusPage = types.BoolNull()
    }
    if val, ok := item["showEpisodesOnStatusPage"].(bool); ok {
        data.ShowEpisodesOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowEpisodesOnStatusPage = types.BoolNull()
    }
    if val, ok := item["showEpisodeHistoryInDays"].(float64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["showEpisodeHistoryInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowEpisodeHistoryInDays = types.NumberNull()
        }
    } else {
        data.ShowEpisodeHistoryInDays = types.NumberNull()
    }
    if val, ok := item["showEpisodeLabelsOnStatusPage"].(bool); ok {
        data.ShowEpisodeLabelsOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowEpisodeLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := item["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolNull()
    }
    if val, ok := item["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    } else {
        data.ShowSubscriberPageOnStatusPage = types.BoolNull()
    }
    if obj, ok := item["ipWhitelist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := item["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := item["enableEmbeddedOverallStatus"].(bool); ok {
        data.EnableEmbeddedOverallStatus = types.BoolValue(val)
    } else {
        data.EnableEmbeddedOverallStatus = types.BoolNull()
    }
    if val, ok := item["showUptimeHistoryInDays"].(float64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["showUptimeHistoryInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowUptimeHistoryInDays = types.NumberNull()
        }
    } else {
        data.ShowUptimeHistoryInDays = types.NumberNull()
    }
    if obj, ok := item["embeddedOverallStatusToken"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
        } else {
            data.EmbeddedOverallStatusToken = types.StringNull()
        }
    } else if val, ok := item["embeddedOverallStatusToken"].(string); ok {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := item["defaultLanguage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DefaultLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultLanguage = types.StringNull()
        }
    } else if val, ok := item["defaultLanguage"].(string); ok {
        data.DefaultLanguage = types.StringValue(val)
    } else {
        data.DefaultLanguage = types.StringNull()
    }
    if obj, ok := item["enabledLanguages"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EnabledLanguages = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EnabledLanguages = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EnabledLanguages = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EnabledLanguages = types.StringValue(string(jsonBytes))
        } else {
            data.EnabledLanguages = types.StringNull()
        }
    } else if val, ok := item["enabledLanguages"].(string); ok {
        data.EnabledLanguages = types.StringValue(val)
    } else {
        data.EnabledLanguages = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
