package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StatusPageResource{}
var _ resource.ResourceWithImportState = &StatusPageResource{}

func NewStatusPageResource() resource.Resource {
    return &StatusPageResource{}
}

// StatusPageResource defines the resource implementation.
type StatusPageResource struct {
    client *Client
}

// StatusPageResourceModel describes the resource data model.
type StatusPageResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    PageTitle types.String `tfsdk:"page_title"`
    PageDescription types.String `tfsdk:"page_description"`
    EnableSearchEngineIndexing types.Bool `tfsdk:"enable_search_engine_indexing"`
    Description types.String `tfsdk:"description"`
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
    CustomFields JSONSubsetValue `tfsdk:"custom_fields"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    SmtpConfigId types.String `tfsdk:"smtp_config_id"`
    CallSmsConfigId types.String `tfsdk:"call_sms_config_id"`
    ShowIncidentHistoryInDays types.Number `tfsdk:"show_incident_history_in_days"`
    ShowAnnouncementHistoryInDays types.Number `tfsdk:"show_announcement_history_in_days"`
    ShowScheduledEventHistoryInDays types.Number `tfsdk:"show_scheduled_event_history_in_days"`
    OverviewPageDescription types.String `tfsdk:"overview_page_description"`
    HidePoweredByOneUptimeBranding types.Bool `tfsdk:"hide_powered_by_one_uptime_branding"`
    DefaultBarColor JSONSubsetValue `tfsdk:"default_bar_color"`
    SubscriberTimezones JSONSubsetValue `tfsdk:"subscriber_timezones"`
    IsReportEnabled types.Bool `tfsdk:"is_report_enabled"`
    ReportStartDateTime RFC3339Value `tfsdk:"report_start_date_time"`
    ReportRecurringInterval JSONSubsetValue `tfsdk:"report_recurring_interval"`
    SendNextReportBy RFC3339Value `tfsdk:"send_next_report_by"`
    ReportDataInDays types.Number `tfsdk:"report_data_in_days"`
    ReportPeriodType types.String `tfsdk:"report_period_type"`
    ReportTimezone types.String `tfsdk:"report_timezone"`
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
    EnabledLanguages JSONSubsetValue `tfsdk:"enabled_languages"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    DowntimeMonitorStatuses types.Set `tfsdk:"downtime_monitor_statuses"`
}

func (r *StatusPageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage status pages for your project.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Required: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Title of your Status Page. This is used for SEO..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Description of your Status Page. This is used for SEO..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_search_engine_indexing": schema.BoolAttribute{
                MarkdownDescription: "Should search engines like Google and Bing be allowed to index this status page? Turn this off to keep the page reachable by link but out of search results..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "favicon_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "logo_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "cover_image_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "header_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Header.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "footer_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Footer.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_css": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom CSS Header.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom JavaScript. This runs when the status page is loaded..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Is this status page public?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_mcp_server": schema.BoolAttribute{
                MarkdownDescription: "Can AI agents read this status page over the public OneUptime MCP server? This does not affect the status page website, its RSS feed, or its public JSON API..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_master_password": schema.BoolAttribute{
                MarkdownDescription: "Require visitors to enter a master password before viewing a private status page..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "master_password": schema.StringAttribute{
                MarkdownDescription: "Password required to unlock a private status page. This value is stored as a secure hash..",
                Optional: true,
                Computed: true,
                Sensitive: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Labels on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Event Labels on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can email subscribers subscribe to this Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which resources to subscribe to?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can SMS subscribers subscribe to this Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Slack subscribers subscribe to this Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_microsoft_teams_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Microsoft Teams subscribers subscribe to this Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_webhook_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Webhook subscribers subscribe to this Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Copyright Text.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "logo_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the logo image, read by screen readers for accessibility..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "cover_image_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Should SSO be required to login to Private Status Page.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "smtp_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "show_incident_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of incident history should be shown on the status page (in days)?.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of announcement history should be shown on the status page (in days)?.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of scheduled event history should be shown on the status page (in days)?.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Overview Page description for your status page. This is a markdown field..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Hide Powered By OneUptime Branding?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "default_bar_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "subscriber_timezones": schema.StringAttribute{
                MarkdownDescription: "Timezones of subscribers to this status page..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Report Enabled for this Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "report_start_date_time": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "report_recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often would you like to send the report?.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "send_next_report_by": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of data should be included in the report?.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(30)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "report_period_type": schema.StringAttribute{
                MarkdownDescription: "Should the report cover a rolling number of days, or the previous whole calendar period?.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Rolling"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "report_timezone": schema.StringAttribute{
                MarkdownDescription: "The timezone report periods and send times are resolved in. A monthly report in this timezone runs from the 1st at 00:00 to the last day at 23:59..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("UTC"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Overall Uptime Percent on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Overall Precision of uptime percent for this status page..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("99.99% (Two Decimal)"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Text to send to subscribers in the footer of the email..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_custom_subscriber_email_notification_footer_text": schema.BoolAttribute{
                MarkdownDescription: "Enable custom footer text in subscriber email notifications..",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incidents on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Announcements on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_episodes_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Episodes on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_episode_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of episode history to show on the status page.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_episode_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Episode Labels on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Maintenance Events on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Subscriber Page on Status Page?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist for this Status Page. One IP per line. Only used if the status page is private..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_embedded_overall_status": schema.BoolAttribute{
                MarkdownDescription: "Enable embedded overall status badge that can be displayed on external websites?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_uptime_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of uptime history should be shown on the status page? Maximum is 90 days..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(90)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "embedded_overall_status_token": schema.StringAttribute{
                MarkdownDescription: "Security token required to access the embedded overall status badge. This token must be provided in the URL..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "default_language": schema.StringAttribute{
                MarkdownDescription: "Default language that the status page is shown in when a visitor arrives for the first time..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("en"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enabled_languages": schema.StringAttribute{
                MarkdownDescription: "Languages offered in the footer language switcher. Leave empty to offer all supported languages..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?.",
                Computed: true,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitors statuses that are considered as \"down\" for this status page..",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (r *StatusPageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *StatusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    statusPageRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := statusPageRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.PageTitle.IsNull() && !data.PageTitle.IsUnknown() {
        requestDataMap["pageTitle"] = data.PageTitle.ValueString()
    }
    if !data.PageDescription.IsNull() && !data.PageDescription.IsUnknown() {
        requestDataMap["pageDescription"] = data.PageDescription.ValueString()
    }
    if !data.EnableSearchEngineIndexing.IsNull() && !data.EnableSearchEngineIndexing.IsUnknown() {
        requestDataMap["enableSearchEngineIndexing"] = data.EnableSearchEngineIndexing.ValueBool()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }
    if !data.FaviconFileId.IsNull() && !data.FaviconFileId.IsUnknown() {
        requestDataMap["faviconFileId"] = data.FaviconFileId.ValueString()
    }
    if !data.LogoFileId.IsNull() && !data.LogoFileId.IsUnknown() {
        requestDataMap["logoFileId"] = data.LogoFileId.ValueString()
    }
    if !data.CoverImageFileId.IsNull() && !data.CoverImageFileId.IsUnknown() {
        requestDataMap["coverImageFileId"] = data.CoverImageFileId.ValueString()
    }
    if !data.HeaderHtml.IsNull() && !data.HeaderHtml.IsUnknown() {
        requestDataMap["headerHTML"] = data.HeaderHtml.ValueString()
    }
    if !data.FooterHtml.IsNull() && !data.FooterHtml.IsUnknown() {
        requestDataMap["footerHTML"] = data.FooterHtml.ValueString()
    }
    if !data.CustomCss.IsNull() && !data.CustomCss.IsUnknown() {
        requestDataMap["customCSS"] = data.CustomCss.ValueString()
    }
    if !data.CustomJavaScript.IsNull() && !data.CustomJavaScript.IsUnknown() {
        requestDataMap["customJavaScript"] = data.CustomJavaScript.ValueString()
    }
    if !data.IsPublicStatusPage.IsNull() && !data.IsPublicStatusPage.IsUnknown() {
        requestDataMap["isPublicStatusPage"] = data.IsPublicStatusPage.ValueBool()
    }
    if !data.EnableMcpServer.IsNull() && !data.EnableMcpServer.IsUnknown() {
        requestDataMap["enableMcpServer"] = data.EnableMcpServer.ValueBool()
    }
    if !data.EnableMasterPassword.IsNull() && !data.EnableMasterPassword.IsUnknown() {
        requestDataMap["enableMasterPassword"] = data.EnableMasterPassword.ValueBool()
    }
    if !data.MasterPassword.IsNull() && !data.MasterPassword.IsUnknown() {
        requestDataMap["masterPassword"] = data.MasterPassword.ValueString()
    }
    if !data.ShowIncidentLabelsOnStatusPage.IsNull() && !data.ShowIncidentLabelsOnStatusPage.IsUnknown() {
        requestDataMap["showIncidentLabelsOnStatusPage"] = data.ShowIncidentLabelsOnStatusPage.ValueBool()
    }
    if !data.ShowScheduledEventLabelsOnStatusPage.IsNull() && !data.ShowScheduledEventLabelsOnStatusPage.IsUnknown() {
        requestDataMap["showScheduledEventLabelsOnStatusPage"] = data.ShowScheduledEventLabelsOnStatusPage.ValueBool()
    }
    if !data.EnableEmailSubscribers.IsNull() && !data.EnableEmailSubscribers.IsUnknown() {
        requestDataMap["enableEmailSubscribers"] = data.EnableEmailSubscribers.ValueBool()
    }
    if !data.AllowSubscribersToChooseResources.IsNull() && !data.AllowSubscribersToChooseResources.IsUnknown() {
        requestDataMap["allowSubscribersToChooseResources"] = data.AllowSubscribersToChooseResources.ValueBool()
    }
    if !data.AllowSubscribersToChooseEventTypes.IsNull() && !data.AllowSubscribersToChooseEventTypes.IsUnknown() {
        requestDataMap["allowSubscribersToChooseEventTypes"] = data.AllowSubscribersToChooseEventTypes.ValueBool()
    }
    if !data.EnableSmsSubscribers.IsNull() && !data.EnableSmsSubscribers.IsUnknown() {
        requestDataMap["enableSmsSubscribers"] = data.EnableSmsSubscribers.ValueBool()
    }
    if !data.EnableSlackSubscribers.IsNull() && !data.EnableSlackSubscribers.IsUnknown() {
        requestDataMap["enableSlackSubscribers"] = data.EnableSlackSubscribers.ValueBool()
    }
    if !data.EnableMicrosoftTeamsSubscribers.IsNull() && !data.EnableMicrosoftTeamsSubscribers.IsUnknown() {
        requestDataMap["enableMicrosoftTeamsSubscribers"] = data.EnableMicrosoftTeamsSubscribers.ValueBool()
    }
    if !data.EnableWebhookSubscribers.IsNull() && !data.EnableWebhookSubscribers.IsUnknown() {
        requestDataMap["enableWebhookSubscribers"] = data.EnableWebhookSubscribers.ValueBool()
    }
    if !data.CopyrightText.IsNull() && !data.CopyrightText.IsUnknown() {
        requestDataMap["copyrightText"] = data.CopyrightText.ValueString()
    }
    if !data.LogoAltText.IsNull() && !data.LogoAltText.IsUnknown() {
        requestDataMap["logoAltText"] = data.LogoAltText.ValueString()
    }
    if !data.CoverImageAltText.IsNull() && !data.CoverImageAltText.IsUnknown() {
        requestDataMap["coverImageAltText"] = data.CoverImageAltText.ValueString()
    }
    if parsedCustomFields := r.parseJSONField(data.CustomFields); parsedCustomFields != nil {
        requestDataMap["customFields"] = parsedCustomFields
    }
    if !data.RequireSsoForLogin.IsNull() && !data.RequireSsoForLogin.IsUnknown() {
        requestDataMap["requireSsoForLogin"] = data.RequireSsoForLogin.ValueBool()
    }
    if !data.SmtpConfigId.IsNull() && !data.SmtpConfigId.IsUnknown() {
        requestDataMap["smtpConfigId"] = data.SmtpConfigId.ValueString()
    }
    if !data.CallSmsConfigId.IsNull() && !data.CallSmsConfigId.IsUnknown() {
        requestDataMap["callSmsConfigId"] = data.CallSmsConfigId.ValueString()
    }
    if !data.ShowIncidentHistoryInDays.IsNull() && !data.ShowIncidentHistoryInDays.IsUnknown() {
        requestDataMap["showIncidentHistoryInDays"] = r.bigFloatToFloat64(data.ShowIncidentHistoryInDays.ValueBigFloat())
    }
    if !data.ShowAnnouncementHistoryInDays.IsNull() && !data.ShowAnnouncementHistoryInDays.IsUnknown() {
        requestDataMap["showAnnouncementHistoryInDays"] = r.bigFloatToFloat64(data.ShowAnnouncementHistoryInDays.ValueBigFloat())
    }
    if !data.ShowScheduledEventHistoryInDays.IsNull() && !data.ShowScheduledEventHistoryInDays.IsUnknown() {
        requestDataMap["showScheduledEventHistoryInDays"] = r.bigFloatToFloat64(data.ShowScheduledEventHistoryInDays.ValueBigFloat())
    }
    if !data.OverviewPageDescription.IsNull() && !data.OverviewPageDescription.IsUnknown() {
        requestDataMap["overviewPageDescription"] = data.OverviewPageDescription.ValueString()
    }
    if !data.HidePoweredByOneUptimeBranding.IsNull() && !data.HidePoweredByOneUptimeBranding.IsUnknown() {
        requestDataMap["hidePoweredByOneUptimeBranding"] = data.HidePoweredByOneUptimeBranding.ValueBool()
    }
    if parsedDefaultBarColor := r.parseJSONField(data.DefaultBarColor); parsedDefaultBarColor != nil {
        requestDataMap["defaultBarColor"] = parsedDefaultBarColor
    }
    if parsedSubscriberTimezones := r.parseJSONField(data.SubscriberTimezones); parsedSubscriberTimezones != nil {
        requestDataMap["subscriberTimezones"] = parsedSubscriberTimezones
    }
    if !data.IsReportEnabled.IsNull() && !data.IsReportEnabled.IsUnknown() {
        requestDataMap["isReportEnabled"] = data.IsReportEnabled.ValueBool()
    }
    if !data.ReportStartDateTime.IsNull() && !data.ReportStartDateTime.IsUnknown() {
        requestDataMap["reportStartDateTime"] = data.ReportStartDateTime.ValueString()
    }
    if parsedReportRecurringInterval := r.parseJSONField(data.ReportRecurringInterval); parsedReportRecurringInterval != nil {
        requestDataMap["reportRecurringInterval"] = parsedReportRecurringInterval
    }
    if !data.SendNextReportBy.IsNull() && !data.SendNextReportBy.IsUnknown() {
        requestDataMap["sendNextReportBy"] = data.SendNextReportBy.ValueString()
    }
    if !data.ReportDataInDays.IsNull() && !data.ReportDataInDays.IsUnknown() {
        requestDataMap["reportDataInDays"] = r.bigFloatToFloat64(data.ReportDataInDays.ValueBigFloat())
    }
    if !data.ReportPeriodType.IsNull() && !data.ReportPeriodType.IsUnknown() {
        requestDataMap["reportPeriodType"] = data.ReportPeriodType.ValueString()
    }
    if !data.ReportTimezone.IsNull() && !data.ReportTimezone.IsUnknown() {
        requestDataMap["reportTimezone"] = data.ReportTimezone.ValueString()
    }
    if !data.ShowOverallUptimePercentOnStatusPage.IsNull() && !data.ShowOverallUptimePercentOnStatusPage.IsUnknown() {
        requestDataMap["showOverallUptimePercentOnStatusPage"] = data.ShowOverallUptimePercentOnStatusPage.ValueBool()
    }
    if !data.OverallUptimePercentPrecision.IsNull() && !data.OverallUptimePercentPrecision.IsUnknown() {
        requestDataMap["overallUptimePercentPrecision"] = data.OverallUptimePercentPrecision.ValueString()
    }
    if !data.SubscriberEmailNotificationFooterText.IsNull() && !data.SubscriberEmailNotificationFooterText.IsUnknown() {
        requestDataMap["subscriberEmailNotificationFooterText"] = data.SubscriberEmailNotificationFooterText.ValueString()
    }
    if !data.EnableCustomSubscriberEmailNotificationFooterText.IsNull() && !data.EnableCustomSubscriberEmailNotificationFooterText.IsUnknown() {
        requestDataMap["enableCustomSubscriberEmailNotificationFooterText"] = data.EnableCustomSubscriberEmailNotificationFooterText.ValueBool()
    }
    if !data.ShowIncidentsOnStatusPage.IsNull() && !data.ShowIncidentsOnStatusPage.IsUnknown() {
        requestDataMap["showIncidentsOnStatusPage"] = data.ShowIncidentsOnStatusPage.ValueBool()
    }
    if !data.ShowAnnouncementsOnStatusPage.IsNull() && !data.ShowAnnouncementsOnStatusPage.IsUnknown() {
        requestDataMap["showAnnouncementsOnStatusPage"] = data.ShowAnnouncementsOnStatusPage.ValueBool()
    }
    if !data.ShowEpisodesOnStatusPage.IsNull() && !data.ShowEpisodesOnStatusPage.IsUnknown() {
        requestDataMap["showEpisodesOnStatusPage"] = data.ShowEpisodesOnStatusPage.ValueBool()
    }
    if !data.ShowEpisodeHistoryInDays.IsNull() && !data.ShowEpisodeHistoryInDays.IsUnknown() {
        requestDataMap["showEpisodeHistoryInDays"] = r.bigFloatToFloat64(data.ShowEpisodeHistoryInDays.ValueBigFloat())
    }
    if !data.ShowEpisodeLabelsOnStatusPage.IsNull() && !data.ShowEpisodeLabelsOnStatusPage.IsUnknown() {
        requestDataMap["showEpisodeLabelsOnStatusPage"] = data.ShowEpisodeLabelsOnStatusPage.ValueBool()
    }
    if !data.ShowScheduledMaintenanceEventsOnStatusPage.IsNull() && !data.ShowScheduledMaintenanceEventsOnStatusPage.IsUnknown() {
        requestDataMap["showScheduledMaintenanceEventsOnStatusPage"] = data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool()
    }
    if !data.ShowSubscriberPageOnStatusPage.IsNull() && !data.ShowSubscriberPageOnStatusPage.IsUnknown() {
        requestDataMap["showSubscriberPageOnStatusPage"] = data.ShowSubscriberPageOnStatusPage.ValueBool()
    }
    if !data.IpWhitelist.IsNull() && !data.IpWhitelist.IsUnknown() {
        requestDataMap["ipWhitelist"] = data.IpWhitelist.ValueString()
    }
    if !data.EnableEmbeddedOverallStatus.IsNull() && !data.EnableEmbeddedOverallStatus.IsUnknown() {
        requestDataMap["enableEmbeddedOverallStatus"] = data.EnableEmbeddedOverallStatus.ValueBool()
    }
    if !data.ShowUptimeHistoryInDays.IsNull() && !data.ShowUptimeHistoryInDays.IsUnknown() {
        requestDataMap["showUptimeHistoryInDays"] = r.bigFloatToFloat64(data.ShowUptimeHistoryInDays.ValueBigFloat())
    }
    if !data.EmbeddedOverallStatusToken.IsNull() && !data.EmbeddedOverallStatusToken.IsUnknown() {
        requestDataMap["embeddedOverallStatusToken"] = data.EmbeddedOverallStatusToken.ValueString()
    }
    if !data.DefaultLanguage.IsNull() && !data.DefaultLanguage.IsUnknown() {
        requestDataMap["defaultLanguage"] = data.DefaultLanguage.ValueString()
    }
    if parsedEnabledLanguages := r.parseJSONField(data.EnabledLanguages); parsedEnabledLanguages != nil {
        requestDataMap["enabledLanguages"] = parsedEnabledLanguages
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/status-page", statusPageRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page, got error: %s", err))
        return
    }

    var statusPageResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create status_page: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := statusPageResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := statusPageResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for status_page did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * status_page is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "pageTitle": true,
        "pageDescription": true,
        "enableSearchEngineIndexing": true,
        "description": true,
        "labels": true,
        "createdByUserId": true,
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
        "showIncidentHistoryInDays": true,
        "showAnnouncementHistoryInDays": true,
        "showScheduledEventHistoryInDays": true,
        "overviewPageDescription": true,
        "hidePoweredByOneUptimeBranding": true,
        "defaultBarColor": true,
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
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "downtimeMonitorStatuses": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/status-page/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created status_page but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created status_page but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["pageTitle"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageTitle = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageTitle = types.StringValue(string(jsonBytes))
            } else {
                data.PageTitle = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageTitle = types.StringValue(string(jsonBytes))
            } else {
                data.PageTitle = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageTitle = types.StringValue(string(jsonBytes))
        } else {
            data.PageTitle = types.StringNull()
        }
    } else if val, ok := dataMap["pageTitle"].(string); ok {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if obj, ok := dataMap["pageDescription"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.PageDescription = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.PageDescription = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.PageDescription = types.StringNull()
        }
    } else if val, ok := dataMap["pageDescription"].(string); ok {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := dataMap["enableSearchEngineIndexing"].(bool); ok {
        data.EnableSearchEngineIndexing = types.BoolValue(val)
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["faviconFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FaviconFileId = types.StringValue(string(jsonBytes))
            } else {
                data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FaviconFileId = types.StringValue(string(jsonBytes))
            } else {
                data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if val, ok := dataMap["faviconFileId"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if obj, ok := dataMap["logoFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogoFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LogoFileId = types.StringValue(string(jsonBytes))
            } else {
                data.LogoFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LogoFileId = types.StringValue(string(jsonBytes))
            } else {
                data.LogoFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if val, ok := dataMap["logoFileId"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if obj, ok := dataMap["coverImageFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CoverImageFileId = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CoverImageFileId = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CoverImageFileId = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageFileId = types.StringNull()
        }
    } else if val, ok := dataMap["coverImageFileId"].(string); ok {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if obj, ok := dataMap["headerHTML"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HeaderHtml = types.StringValue(string(jsonBytes))
            } else {
                data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HeaderHtml = types.StringValue(string(jsonBytes))
            } else {
                data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HeaderHtml = types.StringValue(string(jsonBytes))
        } else {
            data.HeaderHtml = types.StringNull()
        }
    } else if val, ok := dataMap["headerHTML"].(string); ok {
        data.HeaderHtml = types.StringValue(val)
    } else {
        data.HeaderHtml = types.StringNull()
    }
    if obj, ok := dataMap["footerHTML"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FooterHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FooterHtml = types.StringValue(string(jsonBytes))
            } else {
                data.FooterHtml = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FooterHtml = types.StringValue(string(jsonBytes))
            } else {
                data.FooterHtml = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FooterHtml = types.StringValue(string(jsonBytes))
        } else {
            data.FooterHtml = types.StringNull()
        }
    } else if val, ok := dataMap["footerHTML"].(string); ok {
        data.FooterHtml = types.StringValue(val)
    } else {
        data.FooterHtml = types.StringNull()
    }
    if obj, ok := dataMap["customCSS"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomCss = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomCss = types.StringValue(string(jsonBytes))
            } else {
                data.CustomCss = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomCss = types.StringValue(string(jsonBytes))
            } else {
                data.CustomCss = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomCss = types.StringValue(string(jsonBytes))
        } else {
            data.CustomCss = types.StringNull()
        }
    } else if val, ok := dataMap["customCSS"].(string); ok {
        data.CustomCss = types.StringValue(val)
    } else {
        data.CustomCss = types.StringNull()
    }
    if obj, ok := dataMap["customJavaScript"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomJavaScript = types.StringValue(string(jsonBytes))
            } else {
                data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomJavaScript = types.StringValue(string(jsonBytes))
            } else {
                data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomJavaScript = types.StringValue(string(jsonBytes))
        } else {
            data.CustomJavaScript = types.StringNull()
        }
    } else if val, ok := dataMap["customJavaScript"].(string); ok {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMcpServer"].(bool); ok {
        data.EnableMcpServer = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMasterPassword"].(bool); ok {
        data.EnableMasterPassword = types.BoolValue(val)
    }
    if obj, ok := dataMap["masterPassword"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MasterPassword = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MasterPassword = types.StringValue(string(jsonBytes))
            } else {
                data.MasterPassword = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MasterPassword = types.StringValue(string(jsonBytes))
            } else {
                data.MasterPassword = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MasterPassword = types.StringValue(string(jsonBytes))
        } else {
            data.MasterPassword = types.StringNull()
        }
    } else if val, ok := dataMap["masterPassword"].(string); ok {
        data.MasterPassword = types.StringValue(val)
    } else {
        data.MasterPassword = types.StringNull()
    }
    if val, ok := dataMap["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    }
    if val, ok := dataMap["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMicrosoftTeamsSubscribers"].(bool); ok {
        data.EnableMicrosoftTeamsSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableWebhookSubscribers"].(bool); ok {
        data.EnableWebhookSubscribers = types.BoolValue(val)
    }
    if obj, ok := dataMap["copyrightText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CopyrightText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CopyrightText = types.StringValue(string(jsonBytes))
            } else {
                data.CopyrightText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CopyrightText = types.StringValue(string(jsonBytes))
            } else {
                data.CopyrightText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CopyrightText = types.StringValue(string(jsonBytes))
        } else {
            data.CopyrightText = types.StringNull()
        }
    } else if val, ok := dataMap["copyrightText"].(string); ok {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := dataMap["logoAltText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogoAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LogoAltText = types.StringValue(string(jsonBytes))
            } else {
                data.LogoAltText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LogoAltText = types.StringValue(string(jsonBytes))
            } else {
                data.LogoAltText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogoAltText = types.StringValue(string(jsonBytes))
        } else {
            data.LogoAltText = types.StringNull()
        }
    } else if val, ok := dataMap["logoAltText"].(string); ok {
        data.LogoAltText = types.StringValue(val)
    } else {
        data.LogoAltText = types.StringNull()
    }
    if obj, ok := dataMap["coverImageAltText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CoverImageAltText = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CoverImageAltText = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CoverImageAltText = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageAltText = types.StringNull()
        }
    } else if val, ok := dataMap["coverImageAltText"].(string); ok {
        data.CoverImageAltText = types.StringValue(val)
    } else {
        data.CoverImageAltText = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CustomFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok {
        data.CustomFields = NewJSONSubsetValue(val)
    } else {
        data.CustomFields = NewJSONSubsetNull()
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if obj, ok := dataMap["smtpConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SmtpConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SmtpConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SmtpConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.SmtpConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["smtpConfigId"].(string); ok {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if obj, ok := dataMap["callSmsConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallSmsConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallSmsConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallSmsConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.CallSmsConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["callSmsConfigId"].(string); ok {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := dataMap["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showIncidentHistoryInDays"].(int); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showIncidentHistoryInDays"].(int64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showIncidentHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowIncidentHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showAnnouncementHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowAnnouncementHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showScheduledEventHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowScheduledEventHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if obj, ok := dataMap["overviewPageDescription"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OverviewPageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OverviewPageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OverviewPageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.OverviewPageDescription = types.StringNull()
        }
    } else if val, ok := dataMap["overviewPageDescription"].(string); ok {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if obj, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultBarColor = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["defaultBarColor"].(string); ok {
        data.DefaultBarColor = NewJSONSubsetValue(val)
    } else {
        data.DefaultBarColor = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberTimezones = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["subscriberTimezones"].(string); ok {
        data.SubscriberTimezones = NewJSONSubsetValue(val)
    } else {
        data.SubscriberTimezones = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ReportStartDateTime = NewRFC3339Value(val)
        } else {
            data.ReportStartDateTime = NewRFC3339Null()
        }
    } else if val, ok := dataMap["reportStartDateTime"].(string); ok && val != "" {
        data.ReportStartDateTime = NewRFC3339Value(val)
    } else {
        data.ReportStartDateTime = NewRFC3339Null()
    }
    if obj, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportRecurringInterval = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["reportRecurringInterval"].(string); ok {
        data.ReportRecurringInterval = NewJSONSubsetValue(val)
    } else {
        data.ReportRecurringInterval = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SendNextReportBy = NewRFC3339Value(val)
        } else {
            data.SendNextReportBy = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sendNextReportBy"].(string); ok && val != "" {
        data.SendNextReportBy = NewRFC3339Value(val)
    } else {
        data.SendNextReportBy = NewRFC3339Null()
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reportDataInDays"].(int); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reportDataInDays"].(int64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["reportDataInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReportDataInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ReportDataInDays = types.NumberNull()
    }
    if obj, ok := dataMap["reportPeriodType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportPeriodType = types.StringValue(string(jsonBytes))
            } else {
                data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportPeriodType = types.StringValue(string(jsonBytes))
            } else {
                data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportPeriodType = types.StringValue(string(jsonBytes))
        } else {
            data.ReportPeriodType = types.StringNull()
        }
    } else if val, ok := dataMap["reportPeriodType"].(string); ok {
        data.ReportPeriodType = types.StringValue(val)
    } else {
        data.ReportPeriodType = types.StringNull()
    }
    if obj, ok := dataMap["reportTimezone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportTimezone = types.StringValue(string(jsonBytes))
            } else {
                data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportTimezone = types.StringValue(string(jsonBytes))
            } else {
                data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportTimezone = types.StringValue(string(jsonBytes))
        } else {
            data.ReportTimezone = types.StringNull()
        }
    } else if val, ok := dataMap["reportTimezone"].(string); ok {
        data.ReportTimezone = types.StringValue(val)
    } else {
        data.ReportTimezone = types.StringNull()
    }
    if val, ok := dataMap["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    }
    if obj, ok := dataMap["overallUptimePercentPrecision"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
            } else {
                data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
            } else {
                data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
        } else {
            data.OverallUptimePercentPrecision = types.StringNull()
        }
    } else if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if obj, ok := dataMap["subscriberEmailNotificationFooterText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmailNotificationFooterText = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := dataMap["enableCustomSubscriberEmailNotificationFooterText"].(bool); ok {
        data.EnableCustomSubscriberEmailNotificationFooterText = types.BoolValue(val)
    }
    if val, ok := dataMap["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showEpisodesOnStatusPage"].(bool); ok {
        data.ShowEpisodesOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showEpisodeHistoryInDays"].(float64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showEpisodeHistoryInDays"].(int); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showEpisodeHistoryInDays"].(int64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showEpisodeHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowEpisodeHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowEpisodeHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showEpisodeLabelsOnStatusPage"].(bool); ok {
        data.ShowEpisodeLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    }
    if obj, ok := dataMap["ipWhitelist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IpWhitelist = types.StringValue(string(jsonBytes))
            } else {
                data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IpWhitelist = types.StringValue(string(jsonBytes))
            } else {
                data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := dataMap["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := dataMap["enableEmbeddedOverallStatus"].(bool); ok {
        data.EnableEmbeddedOverallStatus = types.BoolValue(val)
    }
    if val, ok := dataMap["showUptimeHistoryInDays"].(float64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showUptimeHistoryInDays"].(int); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showUptimeHistoryInDays"].(int64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showUptimeHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowUptimeHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowUptimeHistoryInDays = types.NumberNull()
    }
    if obj, ok := dataMap["embeddedOverallStatusToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
            } else {
                data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
            } else {
                data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
        } else {
            data.EmbeddedOverallStatusToken = types.StringNull()
        }
    } else if val, ok := dataMap["embeddedOverallStatusToken"].(string); ok {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := dataMap["defaultLanguage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultLanguage = types.StringNull()
        }
    } else if val, ok := dataMap["defaultLanguage"].(string); ok {
        data.DefaultLanguage = types.StringValue(val)
    } else {
        data.DefaultLanguage = types.StringNull()
    }
    if obj, ok := dataMap["enabledLanguages"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EnabledLanguages = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EnabledLanguages = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EnabledLanguages = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["enabledLanguages"].(string); ok {
        data.EnabledLanguages = NewJSONSubsetValue(val)
    } else {
        data.EnabledLanguages = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "pageTitle": true,
        "pageDescription": true,
        "enableSearchEngineIndexing": true,
        "description": true,
        "labels": true,
        "createdByUserId": true,
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
        "showIncidentHistoryInDays": true,
        "showAnnouncementHistoryInDays": true,
        "showScheduledEventHistoryInDays": true,
        "overviewPageDescription": true,
        "hidePoweredByOneUptimeBranding": true,
        "defaultBarColor": true,
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
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "downtimeMonitorStatuses": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/status-page/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var statusPageResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["pageTitle"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageTitle = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageTitle = types.StringValue(string(jsonBytes))
            } else {
                data.PageTitle = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageTitle = types.StringValue(string(jsonBytes))
            } else {
                data.PageTitle = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageTitle = types.StringValue(string(jsonBytes))
        } else {
            data.PageTitle = types.StringNull()
        }
    } else if val, ok := dataMap["pageTitle"].(string); ok {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if obj, ok := dataMap["pageDescription"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.PageDescription = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.PageDescription = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.PageDescription = types.StringNull()
        }
    } else if val, ok := dataMap["pageDescription"].(string); ok {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := dataMap["enableSearchEngineIndexing"].(bool); ok {
        data.EnableSearchEngineIndexing = types.BoolValue(val)
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["faviconFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FaviconFileId = types.StringValue(string(jsonBytes))
            } else {
                data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FaviconFileId = types.StringValue(string(jsonBytes))
            } else {
                data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if val, ok := dataMap["faviconFileId"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if obj, ok := dataMap["logoFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogoFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LogoFileId = types.StringValue(string(jsonBytes))
            } else {
                data.LogoFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LogoFileId = types.StringValue(string(jsonBytes))
            } else {
                data.LogoFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if val, ok := dataMap["logoFileId"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if obj, ok := dataMap["coverImageFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CoverImageFileId = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CoverImageFileId = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CoverImageFileId = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageFileId = types.StringNull()
        }
    } else if val, ok := dataMap["coverImageFileId"].(string); ok {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if obj, ok := dataMap["headerHTML"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HeaderHtml = types.StringValue(string(jsonBytes))
            } else {
                data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HeaderHtml = types.StringValue(string(jsonBytes))
            } else {
                data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HeaderHtml = types.StringValue(string(jsonBytes))
        } else {
            data.HeaderHtml = types.StringNull()
        }
    } else if val, ok := dataMap["headerHTML"].(string); ok {
        data.HeaderHtml = types.StringValue(val)
    } else {
        data.HeaderHtml = types.StringNull()
    }
    if obj, ok := dataMap["footerHTML"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FooterHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FooterHtml = types.StringValue(string(jsonBytes))
            } else {
                data.FooterHtml = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FooterHtml = types.StringValue(string(jsonBytes))
            } else {
                data.FooterHtml = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FooterHtml = types.StringValue(string(jsonBytes))
        } else {
            data.FooterHtml = types.StringNull()
        }
    } else if val, ok := dataMap["footerHTML"].(string); ok {
        data.FooterHtml = types.StringValue(val)
    } else {
        data.FooterHtml = types.StringNull()
    }
    if obj, ok := dataMap["customCSS"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomCss = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomCss = types.StringValue(string(jsonBytes))
            } else {
                data.CustomCss = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomCss = types.StringValue(string(jsonBytes))
            } else {
                data.CustomCss = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomCss = types.StringValue(string(jsonBytes))
        } else {
            data.CustomCss = types.StringNull()
        }
    } else if val, ok := dataMap["customCSS"].(string); ok {
        data.CustomCss = types.StringValue(val)
    } else {
        data.CustomCss = types.StringNull()
    }
    if obj, ok := dataMap["customJavaScript"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomJavaScript = types.StringValue(string(jsonBytes))
            } else {
                data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomJavaScript = types.StringValue(string(jsonBytes))
            } else {
                data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomJavaScript = types.StringValue(string(jsonBytes))
        } else {
            data.CustomJavaScript = types.StringNull()
        }
    } else if val, ok := dataMap["customJavaScript"].(string); ok {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMcpServer"].(bool); ok {
        data.EnableMcpServer = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMasterPassword"].(bool); ok {
        data.EnableMasterPassword = types.BoolValue(val)
    }
    if obj, ok := dataMap["masterPassword"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MasterPassword = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MasterPassword = types.StringValue(string(jsonBytes))
            } else {
                data.MasterPassword = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MasterPassword = types.StringValue(string(jsonBytes))
            } else {
                data.MasterPassword = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MasterPassword = types.StringValue(string(jsonBytes))
        } else {
            data.MasterPassword = types.StringNull()
        }
    } else if val, ok := dataMap["masterPassword"].(string); ok {
        data.MasterPassword = types.StringValue(val)
    } else {
        data.MasterPassword = types.StringNull()
    }
    if val, ok := dataMap["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    }
    if val, ok := dataMap["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMicrosoftTeamsSubscribers"].(bool); ok {
        data.EnableMicrosoftTeamsSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableWebhookSubscribers"].(bool); ok {
        data.EnableWebhookSubscribers = types.BoolValue(val)
    }
    if obj, ok := dataMap["copyrightText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CopyrightText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CopyrightText = types.StringValue(string(jsonBytes))
            } else {
                data.CopyrightText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CopyrightText = types.StringValue(string(jsonBytes))
            } else {
                data.CopyrightText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CopyrightText = types.StringValue(string(jsonBytes))
        } else {
            data.CopyrightText = types.StringNull()
        }
    } else if val, ok := dataMap["copyrightText"].(string); ok {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := dataMap["logoAltText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogoAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LogoAltText = types.StringValue(string(jsonBytes))
            } else {
                data.LogoAltText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LogoAltText = types.StringValue(string(jsonBytes))
            } else {
                data.LogoAltText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogoAltText = types.StringValue(string(jsonBytes))
        } else {
            data.LogoAltText = types.StringNull()
        }
    } else if val, ok := dataMap["logoAltText"].(string); ok {
        data.LogoAltText = types.StringValue(val)
    } else {
        data.LogoAltText = types.StringNull()
    }
    if obj, ok := dataMap["coverImageAltText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CoverImageAltText = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CoverImageAltText = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CoverImageAltText = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageAltText = types.StringNull()
        }
    } else if val, ok := dataMap["coverImageAltText"].(string); ok {
        data.CoverImageAltText = types.StringValue(val)
    } else {
        data.CoverImageAltText = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CustomFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok {
        data.CustomFields = NewJSONSubsetValue(val)
    } else {
        data.CustomFields = NewJSONSubsetNull()
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if obj, ok := dataMap["smtpConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SmtpConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SmtpConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SmtpConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.SmtpConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["smtpConfigId"].(string); ok {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if obj, ok := dataMap["callSmsConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallSmsConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallSmsConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallSmsConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.CallSmsConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["callSmsConfigId"].(string); ok {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := dataMap["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showIncidentHistoryInDays"].(int); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showIncidentHistoryInDays"].(int64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showIncidentHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowIncidentHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showAnnouncementHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowAnnouncementHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showScheduledEventHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowScheduledEventHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if obj, ok := dataMap["overviewPageDescription"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OverviewPageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OverviewPageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OverviewPageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.OverviewPageDescription = types.StringNull()
        }
    } else if val, ok := dataMap["overviewPageDescription"].(string); ok {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if obj, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultBarColor = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["defaultBarColor"].(string); ok {
        data.DefaultBarColor = NewJSONSubsetValue(val)
    } else {
        data.DefaultBarColor = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberTimezones = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["subscriberTimezones"].(string); ok {
        data.SubscriberTimezones = NewJSONSubsetValue(val)
    } else {
        data.SubscriberTimezones = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ReportStartDateTime = NewRFC3339Value(val)
        } else {
            data.ReportStartDateTime = NewRFC3339Null()
        }
    } else if val, ok := dataMap["reportStartDateTime"].(string); ok && val != "" {
        data.ReportStartDateTime = NewRFC3339Value(val)
    } else {
        data.ReportStartDateTime = NewRFC3339Null()
    }
    if obj, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportRecurringInterval = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["reportRecurringInterval"].(string); ok {
        data.ReportRecurringInterval = NewJSONSubsetValue(val)
    } else {
        data.ReportRecurringInterval = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SendNextReportBy = NewRFC3339Value(val)
        } else {
            data.SendNextReportBy = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sendNextReportBy"].(string); ok && val != "" {
        data.SendNextReportBy = NewRFC3339Value(val)
    } else {
        data.SendNextReportBy = NewRFC3339Null()
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reportDataInDays"].(int); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reportDataInDays"].(int64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["reportDataInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReportDataInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ReportDataInDays = types.NumberNull()
    }
    if obj, ok := dataMap["reportPeriodType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportPeriodType = types.StringValue(string(jsonBytes))
            } else {
                data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportPeriodType = types.StringValue(string(jsonBytes))
            } else {
                data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportPeriodType = types.StringValue(string(jsonBytes))
        } else {
            data.ReportPeriodType = types.StringNull()
        }
    } else if val, ok := dataMap["reportPeriodType"].(string); ok {
        data.ReportPeriodType = types.StringValue(val)
    } else {
        data.ReportPeriodType = types.StringNull()
    }
    if obj, ok := dataMap["reportTimezone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportTimezone = types.StringValue(string(jsonBytes))
            } else {
                data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportTimezone = types.StringValue(string(jsonBytes))
            } else {
                data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportTimezone = types.StringValue(string(jsonBytes))
        } else {
            data.ReportTimezone = types.StringNull()
        }
    } else if val, ok := dataMap["reportTimezone"].(string); ok {
        data.ReportTimezone = types.StringValue(val)
    } else {
        data.ReportTimezone = types.StringNull()
    }
    if val, ok := dataMap["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    }
    if obj, ok := dataMap["overallUptimePercentPrecision"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
            } else {
                data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
            } else {
                data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
        } else {
            data.OverallUptimePercentPrecision = types.StringNull()
        }
    } else if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if obj, ok := dataMap["subscriberEmailNotificationFooterText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmailNotificationFooterText = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := dataMap["enableCustomSubscriberEmailNotificationFooterText"].(bool); ok {
        data.EnableCustomSubscriberEmailNotificationFooterText = types.BoolValue(val)
    }
    if val, ok := dataMap["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showEpisodesOnStatusPage"].(bool); ok {
        data.ShowEpisodesOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showEpisodeHistoryInDays"].(float64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showEpisodeHistoryInDays"].(int); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showEpisodeHistoryInDays"].(int64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showEpisodeHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowEpisodeHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowEpisodeHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showEpisodeLabelsOnStatusPage"].(bool); ok {
        data.ShowEpisodeLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    }
    if obj, ok := dataMap["ipWhitelist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IpWhitelist = types.StringValue(string(jsonBytes))
            } else {
                data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IpWhitelist = types.StringValue(string(jsonBytes))
            } else {
                data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := dataMap["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := dataMap["enableEmbeddedOverallStatus"].(bool); ok {
        data.EnableEmbeddedOverallStatus = types.BoolValue(val)
    }
    if val, ok := dataMap["showUptimeHistoryInDays"].(float64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showUptimeHistoryInDays"].(int); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showUptimeHistoryInDays"].(int64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showUptimeHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowUptimeHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowUptimeHistoryInDays = types.NumberNull()
    }
    if obj, ok := dataMap["embeddedOverallStatusToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
            } else {
                data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
            } else {
                data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
        } else {
            data.EmbeddedOverallStatusToken = types.StringNull()
        }
    } else if val, ok := dataMap["embeddedOverallStatusToken"].(string); ok {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := dataMap["defaultLanguage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultLanguage = types.StringNull()
        }
    } else if val, ok := dataMap["defaultLanguage"].(string); ok {
        data.DefaultLanguage = types.StringValue(val)
    } else {
        data.DefaultLanguage = types.StringNull()
    }
    if obj, ok := dataMap["enabledLanguages"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EnabledLanguages = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EnabledLanguages = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EnabledLanguages = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["enabledLanguages"].(string); ok {
        data.EnabledLanguages = NewJSONSubsetValue(val)
    } else {
        data.EnabledLanguages = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StatusPageResourceModel
    var state StatusPageResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    statusPageRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := statusPageRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.PageTitle.IsUnknown() && !state.PageTitle.IsUnknown() && !data.PageTitle.Equal(state.PageTitle) {
        requestDataMap["pageTitle"] = data.PageTitle.ValueString()
    }
    if !data.PageDescription.IsUnknown() && !state.PageDescription.IsUnknown() && !data.PageDescription.Equal(state.PageDescription) {
        requestDataMap["pageDescription"] = data.PageDescription.ValueString()
    }
    if !data.EnableSearchEngineIndexing.IsUnknown() && !state.EnableSearchEngineIndexing.IsUnknown() && !data.EnableSearchEngineIndexing.Equal(state.EnableSearchEngineIndexing) {
        requestDataMap["enableSearchEngineIndexing"] = data.EnableSearchEngineIndexing.ValueBool()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.FaviconFileId.IsUnknown() && !state.FaviconFileId.IsUnknown() && !data.FaviconFileId.Equal(state.FaviconFileId) {
        requestDataMap["faviconFileId"] = data.FaviconFileId.ValueString()
    }
    if !data.LogoFileId.IsUnknown() && !state.LogoFileId.IsUnknown() && !data.LogoFileId.Equal(state.LogoFileId) {
        requestDataMap["logoFileId"] = data.LogoFileId.ValueString()
    }
    if !data.CoverImageFileId.IsUnknown() && !state.CoverImageFileId.IsUnknown() && !data.CoverImageFileId.Equal(state.CoverImageFileId) {
        requestDataMap["coverImageFileId"] = data.CoverImageFileId.ValueString()
    }
    if !data.HeaderHtml.IsUnknown() && !state.HeaderHtml.IsUnknown() && !data.HeaderHtml.Equal(state.HeaderHtml) {
        requestDataMap["headerHTML"] = data.HeaderHtml.ValueString()
    }
    if !data.FooterHtml.IsUnknown() && !state.FooterHtml.IsUnknown() && !data.FooterHtml.Equal(state.FooterHtml) {
        requestDataMap["footerHTML"] = data.FooterHtml.ValueString()
    }
    if !data.CustomCss.IsUnknown() && !state.CustomCss.IsUnknown() && !data.CustomCss.Equal(state.CustomCss) {
        requestDataMap["customCSS"] = data.CustomCss.ValueString()
    }
    if !data.CustomJavaScript.IsUnknown() && !state.CustomJavaScript.IsUnknown() && !data.CustomJavaScript.Equal(state.CustomJavaScript) {
        requestDataMap["customJavaScript"] = data.CustomJavaScript.ValueString()
    }
    if !data.IsPublicStatusPage.IsUnknown() && !state.IsPublicStatusPage.IsUnknown() && !data.IsPublicStatusPage.Equal(state.IsPublicStatusPage) {
        requestDataMap["isPublicStatusPage"] = data.IsPublicStatusPage.ValueBool()
    }
    if !data.EnableMcpServer.IsUnknown() && !state.EnableMcpServer.IsUnknown() && !data.EnableMcpServer.Equal(state.EnableMcpServer) {
        requestDataMap["enableMcpServer"] = data.EnableMcpServer.ValueBool()
    }
    if !data.EnableMasterPassword.IsUnknown() && !state.EnableMasterPassword.IsUnknown() && !data.EnableMasterPassword.Equal(state.EnableMasterPassword) {
        requestDataMap["enableMasterPassword"] = data.EnableMasterPassword.ValueBool()
    }
    if !data.MasterPassword.IsUnknown() && !state.MasterPassword.IsUnknown() && !data.MasterPassword.Equal(state.MasterPassword) {
        requestDataMap["masterPassword"] = data.MasterPassword.ValueString()
    }
    if !data.ShowIncidentLabelsOnStatusPage.IsUnknown() && !state.ShowIncidentLabelsOnStatusPage.IsUnknown() && !data.ShowIncidentLabelsOnStatusPage.Equal(state.ShowIncidentLabelsOnStatusPage) {
        requestDataMap["showIncidentLabelsOnStatusPage"] = data.ShowIncidentLabelsOnStatusPage.ValueBool()
    }
    if !data.ShowScheduledEventLabelsOnStatusPage.IsUnknown() && !state.ShowScheduledEventLabelsOnStatusPage.IsUnknown() && !data.ShowScheduledEventLabelsOnStatusPage.Equal(state.ShowScheduledEventLabelsOnStatusPage) {
        requestDataMap["showScheduledEventLabelsOnStatusPage"] = data.ShowScheduledEventLabelsOnStatusPage.ValueBool()
    }
    if !data.EnableEmailSubscribers.IsUnknown() && !state.EnableEmailSubscribers.IsUnknown() && !data.EnableEmailSubscribers.Equal(state.EnableEmailSubscribers) {
        requestDataMap["enableEmailSubscribers"] = data.EnableEmailSubscribers.ValueBool()
    }
    if !data.AllowSubscribersToChooseResources.IsUnknown() && !state.AllowSubscribersToChooseResources.IsUnknown() && !data.AllowSubscribersToChooseResources.Equal(state.AllowSubscribersToChooseResources) {
        requestDataMap["allowSubscribersToChooseResources"] = data.AllowSubscribersToChooseResources.ValueBool()
    }
    if !data.AllowSubscribersToChooseEventTypes.IsUnknown() && !state.AllowSubscribersToChooseEventTypes.IsUnknown() && !data.AllowSubscribersToChooseEventTypes.Equal(state.AllowSubscribersToChooseEventTypes) {
        requestDataMap["allowSubscribersToChooseEventTypes"] = data.AllowSubscribersToChooseEventTypes.ValueBool()
    }
    if !data.EnableSmsSubscribers.IsUnknown() && !state.EnableSmsSubscribers.IsUnknown() && !data.EnableSmsSubscribers.Equal(state.EnableSmsSubscribers) {
        requestDataMap["enableSmsSubscribers"] = data.EnableSmsSubscribers.ValueBool()
    }
    if !data.EnableSlackSubscribers.IsUnknown() && !state.EnableSlackSubscribers.IsUnknown() && !data.EnableSlackSubscribers.Equal(state.EnableSlackSubscribers) {
        requestDataMap["enableSlackSubscribers"] = data.EnableSlackSubscribers.ValueBool()
    }
    if !data.EnableMicrosoftTeamsSubscribers.IsUnknown() && !state.EnableMicrosoftTeamsSubscribers.IsUnknown() && !data.EnableMicrosoftTeamsSubscribers.Equal(state.EnableMicrosoftTeamsSubscribers) {
        requestDataMap["enableMicrosoftTeamsSubscribers"] = data.EnableMicrosoftTeamsSubscribers.ValueBool()
    }
    if !data.EnableWebhookSubscribers.IsUnknown() && !state.EnableWebhookSubscribers.IsUnknown() && !data.EnableWebhookSubscribers.Equal(state.EnableWebhookSubscribers) {
        requestDataMap["enableWebhookSubscribers"] = data.EnableWebhookSubscribers.ValueBool()
    }
    if !data.CopyrightText.IsUnknown() && !state.CopyrightText.IsUnknown() && !data.CopyrightText.Equal(state.CopyrightText) {
        requestDataMap["copyrightText"] = data.CopyrightText.ValueString()
    }
    if !data.LogoAltText.IsUnknown() && !state.LogoAltText.IsUnknown() && !data.LogoAltText.Equal(state.LogoAltText) {
        requestDataMap["logoAltText"] = data.LogoAltText.ValueString()
    }
    if !data.CoverImageAltText.IsUnknown() && !state.CoverImageAltText.IsUnknown() && !data.CoverImageAltText.Equal(state.CoverImageAltText) {
        requestDataMap["coverImageAltText"] = data.CoverImageAltText.ValueString()
    }
    if !data.CustomFields.IsUnknown() && !state.CustomFields.IsUnknown() && !data.CustomFields.Equal(state.CustomFields) {
        var customfieldsData interface{}
        if err := json.Unmarshal([]byte(data.CustomFields.ValueString()), &customfieldsData); err == nil {
            requestDataMap["customFields"] = customfieldsData
        } else {
            requestDataMap["customFields"] = data.CustomFields.ValueString()
        }
    }
    if !data.RequireSsoForLogin.IsUnknown() && !state.RequireSsoForLogin.IsUnknown() && !data.RequireSsoForLogin.Equal(state.RequireSsoForLogin) {
        requestDataMap["requireSsoForLogin"] = data.RequireSsoForLogin.ValueBool()
    }
    if !data.SmtpConfigId.IsUnknown() && !state.SmtpConfigId.IsUnknown() && !data.SmtpConfigId.Equal(state.SmtpConfigId) {
        requestDataMap["smtpConfigId"] = data.SmtpConfigId.ValueString()
    }
    if !data.CallSmsConfigId.IsUnknown() && !state.CallSmsConfigId.IsUnknown() && !data.CallSmsConfigId.Equal(state.CallSmsConfigId) {
        requestDataMap["callSmsConfigId"] = data.CallSmsConfigId.ValueString()
    }
    if !data.ShowIncidentHistoryInDays.IsUnknown() && !state.ShowIncidentHistoryInDays.IsUnknown() && !data.ShowIncidentHistoryInDays.Equal(state.ShowIncidentHistoryInDays) {
        requestDataMap["showIncidentHistoryInDays"] = r.bigFloatToFloat64(data.ShowIncidentHistoryInDays.ValueBigFloat())
    }
    if !data.ShowAnnouncementHistoryInDays.IsUnknown() && !state.ShowAnnouncementHistoryInDays.IsUnknown() && !data.ShowAnnouncementHistoryInDays.Equal(state.ShowAnnouncementHistoryInDays) {
        requestDataMap["showAnnouncementHistoryInDays"] = r.bigFloatToFloat64(data.ShowAnnouncementHistoryInDays.ValueBigFloat())
    }
    if !data.ShowScheduledEventHistoryInDays.IsUnknown() && !state.ShowScheduledEventHistoryInDays.IsUnknown() && !data.ShowScheduledEventHistoryInDays.Equal(state.ShowScheduledEventHistoryInDays) {
        requestDataMap["showScheduledEventHistoryInDays"] = r.bigFloatToFloat64(data.ShowScheduledEventHistoryInDays.ValueBigFloat())
    }
    if !data.OverviewPageDescription.IsUnknown() && !state.OverviewPageDescription.IsUnknown() && !data.OverviewPageDescription.Equal(state.OverviewPageDescription) {
        requestDataMap["overviewPageDescription"] = data.OverviewPageDescription.ValueString()
    }
    if !data.HidePoweredByOneUptimeBranding.IsUnknown() && !state.HidePoweredByOneUptimeBranding.IsUnknown() && !data.HidePoweredByOneUptimeBranding.Equal(state.HidePoweredByOneUptimeBranding) {
        requestDataMap["hidePoweredByOneUptimeBranding"] = data.HidePoweredByOneUptimeBranding.ValueBool()
    }
    if !data.DefaultBarColor.IsUnknown() && !state.DefaultBarColor.IsUnknown() && !data.DefaultBarColor.Equal(state.DefaultBarColor) {
        var defaultbarcolorData interface{}
        if err := json.Unmarshal([]byte(data.DefaultBarColor.ValueString()), &defaultbarcolorData); err == nil {
            requestDataMap["defaultBarColor"] = defaultbarcolorData
        } else {
            requestDataMap["defaultBarColor"] = data.DefaultBarColor.ValueString()
        }
    }
    if !data.SubscriberTimezones.IsUnknown() && !state.SubscriberTimezones.IsUnknown() && !data.SubscriberTimezones.Equal(state.SubscriberTimezones) {
        var subscribertimezonesData interface{}
        if err := json.Unmarshal([]byte(data.SubscriberTimezones.ValueString()), &subscribertimezonesData); err == nil {
            requestDataMap["subscriberTimezones"] = subscribertimezonesData
        } else {
            requestDataMap["subscriberTimezones"] = data.SubscriberTimezones.ValueString()
        }
    }
    if !data.IsReportEnabled.IsUnknown() && !state.IsReportEnabled.IsUnknown() && !data.IsReportEnabled.Equal(state.IsReportEnabled) {
        requestDataMap["isReportEnabled"] = data.IsReportEnabled.ValueBool()
    }
    if !data.ReportStartDateTime.IsUnknown() && !state.ReportStartDateTime.IsUnknown() && !data.ReportStartDateTime.Equal(state.ReportStartDateTime) {
        requestDataMap["reportStartDateTime"] = data.ReportStartDateTime.ValueString()
    }
    if !data.ReportRecurringInterval.IsUnknown() && !state.ReportRecurringInterval.IsUnknown() && !data.ReportRecurringInterval.Equal(state.ReportRecurringInterval) {
        var reportrecurringintervalData interface{}
        if err := json.Unmarshal([]byte(data.ReportRecurringInterval.ValueString()), &reportrecurringintervalData); err == nil {
            requestDataMap["reportRecurringInterval"] = reportrecurringintervalData
        } else {
            requestDataMap["reportRecurringInterval"] = data.ReportRecurringInterval.ValueString()
        }
    }
    if !data.SendNextReportBy.IsUnknown() && !state.SendNextReportBy.IsUnknown() && !data.SendNextReportBy.Equal(state.SendNextReportBy) {
        requestDataMap["sendNextReportBy"] = data.SendNextReportBy.ValueString()
    }
    if !data.ReportDataInDays.IsUnknown() && !state.ReportDataInDays.IsUnknown() && !data.ReportDataInDays.Equal(state.ReportDataInDays) {
        requestDataMap["reportDataInDays"] = r.bigFloatToFloat64(data.ReportDataInDays.ValueBigFloat())
    }
    if !data.ReportPeriodType.IsUnknown() && !state.ReportPeriodType.IsUnknown() && !data.ReportPeriodType.Equal(state.ReportPeriodType) {
        requestDataMap["reportPeriodType"] = data.ReportPeriodType.ValueString()
    }
    if !data.ReportTimezone.IsUnknown() && !state.ReportTimezone.IsUnknown() && !data.ReportTimezone.Equal(state.ReportTimezone) {
        requestDataMap["reportTimezone"] = data.ReportTimezone.ValueString()
    }
    if !data.ShowOverallUptimePercentOnStatusPage.IsUnknown() && !state.ShowOverallUptimePercentOnStatusPage.IsUnknown() && !data.ShowOverallUptimePercentOnStatusPage.Equal(state.ShowOverallUptimePercentOnStatusPage) {
        requestDataMap["showOverallUptimePercentOnStatusPage"] = data.ShowOverallUptimePercentOnStatusPage.ValueBool()
    }
    if !data.OverallUptimePercentPrecision.IsUnknown() && !state.OverallUptimePercentPrecision.IsUnknown() && !data.OverallUptimePercentPrecision.Equal(state.OverallUptimePercentPrecision) {
        requestDataMap["overallUptimePercentPrecision"] = data.OverallUptimePercentPrecision.ValueString()
    }
    if !data.SubscriberEmailNotificationFooterText.IsUnknown() && !state.SubscriberEmailNotificationFooterText.IsUnknown() && !data.SubscriberEmailNotificationFooterText.Equal(state.SubscriberEmailNotificationFooterText) {
        requestDataMap["subscriberEmailNotificationFooterText"] = data.SubscriberEmailNotificationFooterText.ValueString()
    }
    if !data.EnableCustomSubscriberEmailNotificationFooterText.IsUnknown() && !state.EnableCustomSubscriberEmailNotificationFooterText.IsUnknown() && !data.EnableCustomSubscriberEmailNotificationFooterText.Equal(state.EnableCustomSubscriberEmailNotificationFooterText) {
        requestDataMap["enableCustomSubscriberEmailNotificationFooterText"] = data.EnableCustomSubscriberEmailNotificationFooterText.ValueBool()
    }
    if !data.ShowIncidentsOnStatusPage.IsUnknown() && !state.ShowIncidentsOnStatusPage.IsUnknown() && !data.ShowIncidentsOnStatusPage.Equal(state.ShowIncidentsOnStatusPage) {
        requestDataMap["showIncidentsOnStatusPage"] = data.ShowIncidentsOnStatusPage.ValueBool()
    }
    if !data.ShowAnnouncementsOnStatusPage.IsUnknown() && !state.ShowAnnouncementsOnStatusPage.IsUnknown() && !data.ShowAnnouncementsOnStatusPage.Equal(state.ShowAnnouncementsOnStatusPage) {
        requestDataMap["showAnnouncementsOnStatusPage"] = data.ShowAnnouncementsOnStatusPage.ValueBool()
    }
    if !data.ShowEpisodesOnStatusPage.IsUnknown() && !state.ShowEpisodesOnStatusPage.IsUnknown() && !data.ShowEpisodesOnStatusPage.Equal(state.ShowEpisodesOnStatusPage) {
        requestDataMap["showEpisodesOnStatusPage"] = data.ShowEpisodesOnStatusPage.ValueBool()
    }
    if !data.ShowEpisodeHistoryInDays.IsUnknown() && !state.ShowEpisodeHistoryInDays.IsUnknown() && !data.ShowEpisodeHistoryInDays.Equal(state.ShowEpisodeHistoryInDays) {
        requestDataMap["showEpisodeHistoryInDays"] = r.bigFloatToFloat64(data.ShowEpisodeHistoryInDays.ValueBigFloat())
    }
    if !data.ShowEpisodeLabelsOnStatusPage.IsUnknown() && !state.ShowEpisodeLabelsOnStatusPage.IsUnknown() && !data.ShowEpisodeLabelsOnStatusPage.Equal(state.ShowEpisodeLabelsOnStatusPage) {
        requestDataMap["showEpisodeLabelsOnStatusPage"] = data.ShowEpisodeLabelsOnStatusPage.ValueBool()
    }
    if !data.ShowScheduledMaintenanceEventsOnStatusPage.IsUnknown() && !state.ShowScheduledMaintenanceEventsOnStatusPage.IsUnknown() && !data.ShowScheduledMaintenanceEventsOnStatusPage.Equal(state.ShowScheduledMaintenanceEventsOnStatusPage) {
        requestDataMap["showScheduledMaintenanceEventsOnStatusPage"] = data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool()
    }
    if !data.ShowSubscriberPageOnStatusPage.IsUnknown() && !state.ShowSubscriberPageOnStatusPage.IsUnknown() && !data.ShowSubscriberPageOnStatusPage.Equal(state.ShowSubscriberPageOnStatusPage) {
        requestDataMap["showSubscriberPageOnStatusPage"] = data.ShowSubscriberPageOnStatusPage.ValueBool()
    }
    if !data.IpWhitelist.IsUnknown() && !state.IpWhitelist.IsUnknown() && !data.IpWhitelist.Equal(state.IpWhitelist) {
        requestDataMap["ipWhitelist"] = data.IpWhitelist.ValueString()
    }
    if !data.EnableEmbeddedOverallStatus.IsUnknown() && !state.EnableEmbeddedOverallStatus.IsUnknown() && !data.EnableEmbeddedOverallStatus.Equal(state.EnableEmbeddedOverallStatus) {
        requestDataMap["enableEmbeddedOverallStatus"] = data.EnableEmbeddedOverallStatus.ValueBool()
    }
    if !data.ShowUptimeHistoryInDays.IsUnknown() && !state.ShowUptimeHistoryInDays.IsUnknown() && !data.ShowUptimeHistoryInDays.Equal(state.ShowUptimeHistoryInDays) {
        requestDataMap["showUptimeHistoryInDays"] = r.bigFloatToFloat64(data.ShowUptimeHistoryInDays.ValueBigFloat())
    }
    if !data.EmbeddedOverallStatusToken.IsUnknown() && !state.EmbeddedOverallStatusToken.IsUnknown() && !data.EmbeddedOverallStatusToken.Equal(state.EmbeddedOverallStatusToken) {
        requestDataMap["embeddedOverallStatusToken"] = data.EmbeddedOverallStatusToken.ValueString()
    }
    if !data.DefaultLanguage.IsUnknown() && !state.DefaultLanguage.IsUnknown() && !data.DefaultLanguage.Equal(state.DefaultLanguage) {
        requestDataMap["defaultLanguage"] = data.DefaultLanguage.ValueString()
    }
    if !data.EnabledLanguages.IsUnknown() && !state.EnabledLanguages.IsUnknown() && !data.EnabledLanguages.Equal(state.EnabledLanguages) {
        var enabledlanguagesData interface{}
        if err := json.Unmarshal([]byte(data.EnabledLanguages.ValueString()), &enabledlanguagesData); err == nil {
            requestDataMap["enabledLanguages"] = enabledlanguagesData
        } else {
            requestDataMap["enabledLanguages"] = data.EnabledLanguages.ValueString()
        }
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(statusPageRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/status-page/" + data.Id.ValueString() + "", statusPageRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status_page, got error: %s", err))
            return
        }

        // Parse the update response
        var statusPageResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &statusPageResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update status_page: %s", err))
            return
        }
        _ = statusPageResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "pageTitle": true,
        "pageDescription": true,
        "enableSearchEngineIndexing": true,
        "description": true,
        "labels": true,
        "createdByUserId": true,
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
        "showIncidentHistoryInDays": true,
        "showAnnouncementHistoryInDays": true,
        "showScheduledEventHistoryInDays": true,
        "overviewPageDescription": true,
        "hidePoweredByOneUptimeBranding": true,
        "defaultBarColor": true,
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
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "downtimeMonitorStatuses": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/status-page/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["pageTitle"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageTitle = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageTitle = types.StringValue(string(jsonBytes))
            } else {
                data.PageTitle = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageTitle = types.StringValue(string(jsonBytes))
            } else {
                data.PageTitle = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageTitle = types.StringValue(string(jsonBytes))
        } else {
            data.PageTitle = types.StringNull()
        }
    } else if val, ok := dataMap["pageTitle"].(string); ok {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if obj, ok := dataMap["pageDescription"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.PageDescription = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.PageDescription = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.PageDescription = types.StringNull()
        }
    } else if val, ok := dataMap["pageDescription"].(string); ok {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := dataMap["enableSearchEngineIndexing"].(bool); ok {
        data.EnableSearchEngineIndexing = types.BoolValue(val)
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["faviconFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FaviconFileId = types.StringValue(string(jsonBytes))
            } else {
                data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FaviconFileId = types.StringValue(string(jsonBytes))
            } else {
                data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if val, ok := dataMap["faviconFileId"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if obj, ok := dataMap["logoFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogoFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LogoFileId = types.StringValue(string(jsonBytes))
            } else {
                data.LogoFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LogoFileId = types.StringValue(string(jsonBytes))
            } else {
                data.LogoFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if val, ok := dataMap["logoFileId"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if obj, ok := dataMap["coverImageFileId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CoverImageFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CoverImageFileId = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CoverImageFileId = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageFileId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CoverImageFileId = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageFileId = types.StringNull()
        }
    } else if val, ok := dataMap["coverImageFileId"].(string); ok {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if obj, ok := dataMap["headerHTML"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HeaderHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HeaderHtml = types.StringValue(string(jsonBytes))
            } else {
                data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HeaderHtml = types.StringValue(string(jsonBytes))
            } else {
                data.HeaderHtml = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HeaderHtml = types.StringValue(string(jsonBytes))
        } else {
            data.HeaderHtml = types.StringNull()
        }
    } else if val, ok := dataMap["headerHTML"].(string); ok {
        data.HeaderHtml = types.StringValue(val)
    } else {
        data.HeaderHtml = types.StringNull()
    }
    if obj, ok := dataMap["footerHTML"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FooterHtml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FooterHtml = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FooterHtml = types.StringValue(string(jsonBytes))
            } else {
                data.FooterHtml = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FooterHtml = types.StringValue(string(jsonBytes))
            } else {
                data.FooterHtml = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FooterHtml = types.StringValue(string(jsonBytes))
        } else {
            data.FooterHtml = types.StringNull()
        }
    } else if val, ok := dataMap["footerHTML"].(string); ok {
        data.FooterHtml = types.StringValue(val)
    } else {
        data.FooterHtml = types.StringNull()
    }
    if obj, ok := dataMap["customCSS"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomCss = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomCss = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomCss = types.StringValue(string(jsonBytes))
            } else {
                data.CustomCss = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomCss = types.StringValue(string(jsonBytes))
            } else {
                data.CustomCss = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomCss = types.StringValue(string(jsonBytes))
        } else {
            data.CustomCss = types.StringNull()
        }
    } else if val, ok := dataMap["customCSS"].(string); ok {
        data.CustomCss = types.StringValue(val)
    } else {
        data.CustomCss = types.StringNull()
    }
    if obj, ok := dataMap["customJavaScript"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomJavaScript = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomJavaScript = types.StringValue(string(jsonBytes))
            } else {
                data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomJavaScript = types.StringValue(string(jsonBytes))
            } else {
                data.CustomJavaScript = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomJavaScript = types.StringValue(string(jsonBytes))
        } else {
            data.CustomJavaScript = types.StringNull()
        }
    } else if val, ok := dataMap["customJavaScript"].(string); ok {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMcpServer"].(bool); ok {
        data.EnableMcpServer = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMasterPassword"].(bool); ok {
        data.EnableMasterPassword = types.BoolValue(val)
    }
    if obj, ok := dataMap["masterPassword"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MasterPassword = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MasterPassword = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MasterPassword = types.StringValue(string(jsonBytes))
            } else {
                data.MasterPassword = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MasterPassword = types.StringValue(string(jsonBytes))
            } else {
                data.MasterPassword = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MasterPassword = types.StringValue(string(jsonBytes))
        } else {
            data.MasterPassword = types.StringNull()
        }
    } else if val, ok := dataMap["masterPassword"].(string); ok {
        data.MasterPassword = types.StringValue(val)
    } else {
        data.MasterPassword = types.StringNull()
    }
    if val, ok := dataMap["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    }
    if val, ok := dataMap["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableMicrosoftTeamsSubscribers"].(bool); ok {
        data.EnableMicrosoftTeamsSubscribers = types.BoolValue(val)
    }
    if val, ok := dataMap["enableWebhookSubscribers"].(bool); ok {
        data.EnableWebhookSubscribers = types.BoolValue(val)
    }
    if obj, ok := dataMap["copyrightText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CopyrightText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CopyrightText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CopyrightText = types.StringValue(string(jsonBytes))
            } else {
                data.CopyrightText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CopyrightText = types.StringValue(string(jsonBytes))
            } else {
                data.CopyrightText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CopyrightText = types.StringValue(string(jsonBytes))
        } else {
            data.CopyrightText = types.StringNull()
        }
    } else if val, ok := dataMap["copyrightText"].(string); ok {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := dataMap["logoAltText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogoAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LogoAltText = types.StringValue(string(jsonBytes))
            } else {
                data.LogoAltText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LogoAltText = types.StringValue(string(jsonBytes))
            } else {
                data.LogoAltText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogoAltText = types.StringValue(string(jsonBytes))
        } else {
            data.LogoAltText = types.StringNull()
        }
    } else if val, ok := dataMap["logoAltText"].(string); ok {
        data.LogoAltText = types.StringValue(val)
    } else {
        data.LogoAltText = types.StringNull()
    }
    if obj, ok := dataMap["coverImageAltText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CoverImageAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CoverImageAltText = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CoverImageAltText = types.StringValue(string(jsonBytes))
            } else {
                data.CoverImageAltText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CoverImageAltText = types.StringValue(string(jsonBytes))
        } else {
            data.CoverImageAltText = types.StringNull()
        }
    } else if val, ok := dataMap["coverImageAltText"].(string); ok {
        data.CoverImageAltText = types.StringValue(val)
    } else {
        data.CoverImageAltText = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CustomFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok {
        data.CustomFields = NewJSONSubsetValue(val)
    } else {
        data.CustomFields = NewJSONSubsetNull()
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if obj, ok := dataMap["smtpConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SmtpConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SmtpConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SmtpConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.SmtpConfigId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SmtpConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.SmtpConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["smtpConfigId"].(string); ok {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if obj, ok := dataMap["callSmsConfigId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CallSmsConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CallSmsConfigId = types.StringValue(string(jsonBytes))
            } else {
                data.CallSmsConfigId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CallSmsConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.CallSmsConfigId = types.StringNull()
        }
    } else if val, ok := dataMap["callSmsConfigId"].(string); ok {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := dataMap["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showIncidentHistoryInDays"].(int); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showIncidentHistoryInDays"].(int64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showIncidentHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowIncidentHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showAnnouncementHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowAnnouncementHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showScheduledEventHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowScheduledEventHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if obj, ok := dataMap["overviewPageDescription"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OverviewPageDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OverviewPageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OverviewPageDescription = types.StringValue(string(jsonBytes))
            } else {
                data.OverviewPageDescription = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OverviewPageDescription = types.StringValue(string(jsonBytes))
        } else {
            data.OverviewPageDescription = types.StringNull()
        }
    } else if val, ok := dataMap["overviewPageDescription"].(string); ok {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if obj, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultBarColor = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultBarColor = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["defaultBarColor"].(string); ok {
        data.DefaultBarColor = NewJSONSubsetValue(val)
    } else {
        data.DefaultBarColor = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberTimezones = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberTimezones = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["subscriberTimezones"].(string); ok {
        data.SubscriberTimezones = NewJSONSubsetValue(val)
    } else {
        data.SubscriberTimezones = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ReportStartDateTime = NewRFC3339Value(val)
        } else {
            data.ReportStartDateTime = NewRFC3339Null()
        }
    } else if val, ok := dataMap["reportStartDateTime"].(string); ok && val != "" {
        data.ReportStartDateTime = NewRFC3339Value(val)
    } else {
        data.ReportStartDateTime = NewRFC3339Null()
    }
    if obj, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportRecurringInterval = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportRecurringInterval = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["reportRecurringInterval"].(string); ok {
        data.ReportRecurringInterval = NewJSONSubsetValue(val)
    } else {
        data.ReportRecurringInterval = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SendNextReportBy = NewRFC3339Value(val)
        } else {
            data.SendNextReportBy = NewRFC3339Null()
        }
    } else if val, ok := dataMap["sendNextReportBy"].(string); ok && val != "" {
        data.SendNextReportBy = NewRFC3339Value(val)
    } else {
        data.SendNextReportBy = NewRFC3339Null()
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reportDataInDays"].(int); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reportDataInDays"].(int64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["reportDataInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ReportDataInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ReportDataInDays = types.NumberNull()
    }
    if obj, ok := dataMap["reportPeriodType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportPeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportPeriodType = types.StringValue(string(jsonBytes))
            } else {
                data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportPeriodType = types.StringValue(string(jsonBytes))
            } else {
                data.ReportPeriodType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportPeriodType = types.StringValue(string(jsonBytes))
        } else {
            data.ReportPeriodType = types.StringNull()
        }
    } else if val, ok := dataMap["reportPeriodType"].(string); ok {
        data.ReportPeriodType = types.StringValue(val)
    } else {
        data.ReportPeriodType = types.StringNull()
    }
    if obj, ok := dataMap["reportTimezone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportTimezone = types.StringValue(string(jsonBytes))
            } else {
                data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportTimezone = types.StringValue(string(jsonBytes))
            } else {
                data.ReportTimezone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportTimezone = types.StringValue(string(jsonBytes))
        } else {
            data.ReportTimezone = types.StringNull()
        }
    } else if val, ok := dataMap["reportTimezone"].(string); ok {
        data.ReportTimezone = types.StringValue(val)
    } else {
        data.ReportTimezone = types.StringNull()
    }
    if val, ok := dataMap["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    }
    if obj, ok := dataMap["overallUptimePercentPrecision"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OverallUptimePercentPrecision = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
            } else {
                data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
            } else {
                data.OverallUptimePercentPrecision = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OverallUptimePercentPrecision = types.StringValue(string(jsonBytes))
        } else {
            data.OverallUptimePercentPrecision = types.StringNull()
        }
    } else if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if obj, ok := dataMap["subscriberEmailNotificationFooterText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberEmailNotificationFooterText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmailNotificationFooterText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberEmailNotificationFooterText = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmailNotificationFooterText = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := dataMap["enableCustomSubscriberEmailNotificationFooterText"].(bool); ok {
        data.EnableCustomSubscriberEmailNotificationFooterText = types.BoolValue(val)
    }
    if val, ok := dataMap["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showEpisodesOnStatusPage"].(bool); ok {
        data.ShowEpisodesOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showEpisodeHistoryInDays"].(float64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showEpisodeHistoryInDays"].(int); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showEpisodeHistoryInDays"].(int64); ok {
        data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showEpisodeHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowEpisodeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowEpisodeHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowEpisodeHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showEpisodeLabelsOnStatusPage"].(bool); ok {
        data.ShowEpisodeLabelsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    }
    if obj, ok := dataMap["ipWhitelist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IpWhitelist = types.StringValue(string(jsonBytes))
            } else {
                data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IpWhitelist = types.StringValue(string(jsonBytes))
            } else {
                data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := dataMap["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := dataMap["enableEmbeddedOverallStatus"].(bool); ok {
        data.EnableEmbeddedOverallStatus = types.BoolValue(val)
    }
    if val, ok := dataMap["showUptimeHistoryInDays"].(float64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showUptimeHistoryInDays"].(int); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showUptimeHistoryInDays"].(int64); ok {
        data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["showUptimeHistoryInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ShowUptimeHistoryInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShowUptimeHistoryInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ShowUptimeHistoryInDays = types.NumberNull()
    }
    if obj, ok := dataMap["embeddedOverallStatusToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EmbeddedOverallStatusToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
            } else {
                data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
            } else {
                data.EmbeddedOverallStatusToken = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EmbeddedOverallStatusToken = types.StringValue(string(jsonBytes))
        } else {
            data.EmbeddedOverallStatusToken = types.StringNull()
        }
    } else if val, ok := dataMap["embeddedOverallStatusToken"].(string); ok {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := dataMap["defaultLanguage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultLanguage = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultLanguage = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultLanguage = types.StringNull()
        }
    } else if val, ok := dataMap["defaultLanguage"].(string); ok {
        data.DefaultLanguage = types.StringValue(val)
    } else {
        data.DefaultLanguage = types.StringNull()
    }
    if obj, ok := dataMap["enabledLanguages"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EnabledLanguages = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EnabledLanguages = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EnabledLanguages = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EnabledLanguages = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EnabledLanguages = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["enabledLanguages"].(string); ok {
        data.EnabledLanguages = NewJSONSubsetValue(val)
    } else {
        data.EnabledLanguages = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DowntimeMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/status-page/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status_page, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete status_page: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *StatusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *StatusPageResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *StatusPageResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *StatusPageResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *StatusPageResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *StatusPageResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *StatusPageResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *StatusPageResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *StatusPageResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
