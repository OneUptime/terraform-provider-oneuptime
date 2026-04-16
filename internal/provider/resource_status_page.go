package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
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
    Description types.String `tfsdk:"description"`
    Labels types.Set `tfsdk:"labels"`
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
    CopyrightText types.String `tfsdk:"copyright_text"`
    CustomFields types.String `tfsdk:"custom_fields"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    SmtpConfigId types.String `tfsdk:"smtp_config_id"`
    CallSmsConfigId types.String `tfsdk:"call_sms_config_id"`
    ShowIncidentHistoryInDays types.Number `tfsdk:"show_incident_history_in_days"`
    ShowAnnouncementHistoryInDays types.Number `tfsdk:"show_announcement_history_in_days"`
    ShowScheduledEventHistoryInDays types.Number `tfsdk:"show_scheduled_event_history_in_days"`
    OverviewPageDescription types.String `tfsdk:"overview_page_description"`
    HidePoweredByOneUptimeBranding types.Bool `tfsdk:"hide_powered_by_one_uptime_branding"`
    DefaultBarColor types.String `tfsdk:"default_bar_color"`
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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    DowntimeMonitorStatuses types.Set `tfsdk:"downtime_monitor_statuses"`
}

func (r *StatusPageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
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
                MarkdownDescription: "Any friendly name of this object. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Required: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Title of your Status Page. This is used for SEO.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Description of your Status Page. This is used for SEO.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
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
                MarkdownDescription: "Status Page Custom HTML Header. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "footer_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Footer. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_css": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom CSS Header. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom JavaScript. This runs when the status page is loaded.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Is this status page public?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_master_password": schema.BoolAttribute{
                MarkdownDescription: "Require visitors to enter a master password before viewing a private status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "master_password": schema.StringAttribute{
                MarkdownDescription: "Password required to unlock a private status page. This value is stored as a secure hash.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Labels on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Event Labels on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can email subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which resources to subscribe to?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can SMS subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Slack subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_microsoft_teams_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Microsoft Teams subscribers subscribe to this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Copyright Text. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Should SSO be required to login to Private Status Page. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Public, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
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
                MarkdownDescription: "How many days of incident history should be shown on the status page (in days)?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of announcement history should be shown on the status page (in days)?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of scheduled event history should be shown on the status page (in days)?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Overview Page description for your status page. This is a markdown field.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Hide Powered By OneUptime Branding?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "default_bar_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_timezones": schema.StringAttribute{
                MarkdownDescription: "Timezones of subscribers to this status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Report Enabled for this Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "report_start_date_time": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "report_recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often would you like to send the report?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "send_next_report_by": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of data should be included in the report?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(30)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Overall Uptime Percent on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Overall Precision of uptime percent for this status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("99.99% (Two Decimal)"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Text to send to subscribers in the footer of the email.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_custom_subscriber_email_notification_footer_text": schema.BoolAttribute{
                MarkdownDescription: "Enable custom footer text in subscriber email notifications.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incidents on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Announcements on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_episodes_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Episodes on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_episode_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of episode history to show on the status page. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(14)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "show_episode_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Episode Labels on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Maintenance Events on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Subscriber Page on Status Page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist for this Status Page. One IP per line. Only used if the status page is private.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_embedded_overall_status": schema.BoolAttribute{
                MarkdownDescription: "Enable embedded overall status badge that can be displayed on external websites?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_uptime_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of uptime history should be shown on the status page? Maximum is 90 days.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(90)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "embedded_overall_status_token": schema.StringAttribute{
                MarkdownDescription: "Security token required to access the embedded overall status badge. This token must be provided in the URL.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitors statuses that are considered as \"down\" for this status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page]",
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



    // Create API request body
    statusPageRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "name": data.Name.ValueString(),
        "pageTitle": data.PageTitle.ValueString(),
        "pageDescription": data.PageDescription.ValueString(),
        "description": data.Description.ValueString(),
        "labels": r.convertTerraformSetToInterface(data.Labels),
        "faviconFileId": data.FaviconFileId.ValueString(),
        "logoFileId": data.LogoFileId.ValueString(),
        "coverImageFileId": data.CoverImageFileId.ValueString(),
        "headerHTML": data.HeaderHtml.ValueString(),
        "footerHTML": data.FooterHtml.ValueString(),
        "customCSS": data.CustomCss.ValueString(),
        "customJavaScript": data.CustomJavaScript.ValueString(),
        "isPublicStatusPage": data.IsPublicStatusPage.ValueBool(),
        "enableMasterPassword": data.EnableMasterPassword.ValueBool(),
        "masterPassword": data.MasterPassword.ValueString(),
        "showIncidentLabelsOnStatusPage": data.ShowIncidentLabelsOnStatusPage.ValueBool(),
        "showScheduledEventLabelsOnStatusPage": data.ShowScheduledEventLabelsOnStatusPage.ValueBool(),
        "enableSubscribers": data.EnableSubscribers.ValueBool(),
        "enableEmailSubscribers": data.EnableEmailSubscribers.ValueBool(),
        "allowSubscribersToChooseResources": data.AllowSubscribersToChooseResources.ValueBool(),
        "allowSubscribersToChooseEventTypes": data.AllowSubscribersToChooseEventTypes.ValueBool(),
        "enableSmsSubscribers": data.EnableSmsSubscribers.ValueBool(),
        "enableSlackSubscribers": data.EnableSlackSubscribers.ValueBool(),
        "enableMicrosoftTeamsSubscribers": data.EnableMicrosoftTeamsSubscribers.ValueBool(),
        "copyrightText": data.CopyrightText.ValueString(),
        "customFields": r.parseJSONField(data.CustomFields),
        "requireSsoForLogin": data.RequireSsoForLogin.ValueBool(),
        "smtpConfigId": data.SmtpConfigId.ValueString(),
        "callSmsConfigId": data.CallSmsConfigId.ValueString(),
        "showIncidentHistoryInDays": r.bigFloatToFloat64(data.ShowIncidentHistoryInDays.ValueBigFloat()),
        "showAnnouncementHistoryInDays": r.bigFloatToFloat64(data.ShowAnnouncementHistoryInDays.ValueBigFloat()),
        "showScheduledEventHistoryInDays": r.bigFloatToFloat64(data.ShowScheduledEventHistoryInDays.ValueBigFloat()),
        "overviewPageDescription": data.OverviewPageDescription.ValueString(),
        "hidePoweredByOneUptimeBranding": data.HidePoweredByOneUptimeBranding.ValueBool(),
        "defaultBarColor": r.parseJSONField(data.DefaultBarColor),
        "subscriberTimezones": r.parseJSONField(data.SubscriberTimezones),
        "isReportEnabled": data.IsReportEnabled.ValueBool(),
        "reportStartDateTime": r.parseJSONField(data.ReportStartDateTime),
        "reportRecurringInterval": r.parseJSONField(data.ReportRecurringInterval),
        "sendNextReportBy": r.parseJSONField(data.SendNextReportBy),
        "reportDataInDays": r.bigFloatToFloat64(data.ReportDataInDays.ValueBigFloat()),
        "showOverallUptimePercentOnStatusPage": data.ShowOverallUptimePercentOnStatusPage.ValueBool(),
        "overallUptimePercentPrecision": data.OverallUptimePercentPrecision.ValueString(),
        "subscriberEmailNotificationFooterText": data.SubscriberEmailNotificationFooterText.ValueString(),
        "enableCustomSubscriberEmailNotificationFooterText": data.EnableCustomSubscriberEmailNotificationFooterText.ValueBool(),
        "showIncidentsOnStatusPage": data.ShowIncidentsOnStatusPage.ValueBool(),
        "showAnnouncementsOnStatusPage": data.ShowAnnouncementsOnStatusPage.ValueBool(),
        "showEpisodesOnStatusPage": data.ShowEpisodesOnStatusPage.ValueBool(),
        "showEpisodeHistoryInDays": r.bigFloatToFloat64(data.ShowEpisodeHistoryInDays.ValueBigFloat()),
        "showEpisodeLabelsOnStatusPage": data.ShowEpisodeLabelsOnStatusPage.ValueBool(),
        "showScheduledMaintenanceEventsOnStatusPage": data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool(),
        "showSubscriberPageOnStatusPage": data.ShowSubscriberPageOnStatusPage.ValueBool(),
        "ipWhitelist": data.IpWhitelist.ValueString(),
        "enableEmbeddedOverallStatus": data.EnableEmbeddedOverallStatus.ValueBool(),
        "showUptimeHistoryInDays": r.bigFloatToFloat64(data.ShowUptimeHistoryInDays.ValueBigFloat()),
        "embeddedOverallStatusToken": data.EmbeddedOverallStatusToken.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/status-page", statusPageRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page, got error: %s", err))
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["pageTitle"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["pageDescription"].(string); ok && val != "" {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
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
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["faviconFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["logoFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["coverImageFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["headerHTML"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["footerHTML"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["customCSS"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["customJavaScript"].(string); ok && val != "" {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
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
    } else if val, ok := dataMap["masterPassword"].(string); ok && val != "" {
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
    if val, ok := dataMap["enableSubscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
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
    } else if val, ok := dataMap["copyrightText"].(string); ok && val != "" {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok && val != "" {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
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
    } else if val, ok := dataMap["smtpConfigId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["callSmsConfigId"].(string); ok && val != "" {
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
    } else if dataMap["showIncidentHistoryInDays"] == nil {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["showAnnouncementHistoryInDays"] == nil {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["showScheduledEventHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["overviewPageDescription"].(string); ok && val != "" {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if obj, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultBarColor = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultBarColor = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultBarColor = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = types.StringNull()
        }
    } else if val, ok := dataMap["defaultBarColor"].(string); ok && val != "" {
        data.DefaultBarColor = types.StringValue(val)
    } else {
        data.DefaultBarColor = types.StringNull()
    }
    if obj, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberTimezones = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberTimezones = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberTimezones = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberTimezones"].(string); ok && val != "" {
        data.SubscriberTimezones = types.StringValue(val)
    } else {
        data.SubscriberTimezones = types.StringNull()
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportStartDateTime = types.StringValue(string(jsonBytes))
            } else {
                data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportStartDateTime = types.StringValue(string(jsonBytes))
            } else {
                data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportStartDateTime = types.StringValue(string(jsonBytes))
        } else {
            data.ReportStartDateTime = types.StringNull()
        }
    } else if val, ok := dataMap["reportStartDateTime"].(string); ok && val != "" {
        data.ReportStartDateTime = types.StringValue(val)
    } else {
        data.ReportStartDateTime = types.StringNull()
    }
    if obj, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = types.StringNull()
        }
    } else if val, ok := dataMap["reportRecurringInterval"].(string); ok && val != "" {
        data.ReportRecurringInterval = types.StringValue(val)
    } else {
        data.ReportRecurringInterval = types.StringNull()
    }
    if obj, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SendNextReportBy = types.StringValue(string(jsonBytes))
            } else {
                data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SendNextReportBy = types.StringValue(string(jsonBytes))
            } else {
                data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SendNextReportBy = types.StringValue(string(jsonBytes))
        } else {
            data.SendNextReportBy = types.StringNull()
        }
    } else if val, ok := dataMap["sendNextReportBy"].(string); ok && val != "" {
        data.SendNextReportBy = types.StringValue(val)
    } else {
        data.SendNextReportBy = types.StringNull()
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reportDataInDays"].(int); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reportDataInDays"].(int64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["reportDataInDays"] == nil {
        data.ReportDataInDays = types.NumberNull()
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
    } else if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok && val != "" {
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
    } else if dataMap["showEpisodeHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["ipWhitelist"].(string); ok && val != "" {
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
    } else if dataMap["showUptimeHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["embeddedOverallStatusToken"].(string); ok && val != "" {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
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
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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
        "description": true,
        "labels": true,
        "faviconFileId": true,
        "logoFileId": true,
        "coverImageFileId": true,
        "headerHTML": true,
        "footerHTML": true,
        "customCSS": true,
        "customJavaScript": true,
        "isPublicStatusPage": true,
        "enableMasterPassword": true,
        "masterPassword": true,
        "showIncidentLabelsOnStatusPage": true,
        "showScheduledEventLabelsOnStatusPage": true,
        "enableSubscribers": true,
        "enableEmailSubscribers": true,
        "allowSubscribersToChooseResources": true,
        "allowSubscribersToChooseEventTypes": true,
        "enableSmsSubscribers": true,
        "enableSlackSubscribers": true,
        "enableMicrosoftTeamsSubscribers": true,
        "copyrightText": true,
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
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "downtimeMonitorStatuses": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/status-page/" + data.Id.ValueString() + "/get-item", selectParam)
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["pageTitle"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["pageDescription"].(string); ok && val != "" {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
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
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["faviconFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["logoFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["coverImageFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["headerHTML"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["footerHTML"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["customCSS"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["customJavaScript"].(string); ok && val != "" {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
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
    } else if val, ok := dataMap["masterPassword"].(string); ok && val != "" {
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
    if val, ok := dataMap["enableSubscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
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
    } else if val, ok := dataMap["copyrightText"].(string); ok && val != "" {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok && val != "" {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
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
    } else if val, ok := dataMap["smtpConfigId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["callSmsConfigId"].(string); ok && val != "" {
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
    } else if dataMap["showIncidentHistoryInDays"] == nil {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["showAnnouncementHistoryInDays"] == nil {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["showScheduledEventHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["overviewPageDescription"].(string); ok && val != "" {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if obj, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultBarColor = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultBarColor = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultBarColor = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = types.StringNull()
        }
    } else if val, ok := dataMap["defaultBarColor"].(string); ok && val != "" {
        data.DefaultBarColor = types.StringValue(val)
    } else {
        data.DefaultBarColor = types.StringNull()
    }
    if obj, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberTimezones = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberTimezones = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberTimezones = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberTimezones"].(string); ok && val != "" {
        data.SubscriberTimezones = types.StringValue(val)
    } else {
        data.SubscriberTimezones = types.StringNull()
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportStartDateTime = types.StringValue(string(jsonBytes))
            } else {
                data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportStartDateTime = types.StringValue(string(jsonBytes))
            } else {
                data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportStartDateTime = types.StringValue(string(jsonBytes))
        } else {
            data.ReportStartDateTime = types.StringNull()
        }
    } else if val, ok := dataMap["reportStartDateTime"].(string); ok && val != "" {
        data.ReportStartDateTime = types.StringValue(val)
    } else {
        data.ReportStartDateTime = types.StringNull()
    }
    if obj, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = types.StringNull()
        }
    } else if val, ok := dataMap["reportRecurringInterval"].(string); ok && val != "" {
        data.ReportRecurringInterval = types.StringValue(val)
    } else {
        data.ReportRecurringInterval = types.StringNull()
    }
    if obj, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SendNextReportBy = types.StringValue(string(jsonBytes))
            } else {
                data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SendNextReportBy = types.StringValue(string(jsonBytes))
            } else {
                data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SendNextReportBy = types.StringValue(string(jsonBytes))
        } else {
            data.SendNextReportBy = types.StringNull()
        }
    } else if val, ok := dataMap["sendNextReportBy"].(string); ok && val != "" {
        data.SendNextReportBy = types.StringValue(val)
    } else {
        data.SendNextReportBy = types.StringNull()
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reportDataInDays"].(int); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reportDataInDays"].(int64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["reportDataInDays"] == nil {
        data.ReportDataInDays = types.NumberNull()
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
    } else if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok && val != "" {
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
    } else if dataMap["showEpisodeHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["ipWhitelist"].(string); ok && val != "" {
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
    } else if dataMap["showUptimeHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["embeddedOverallStatusToken"].(string); ok && val != "" {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
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
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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
    if !data.EnableSubscribers.IsUnknown() && !state.EnableSubscribers.IsUnknown() && !data.EnableSubscribers.Equal(state.EnableSubscribers) {
        requestDataMap["enableSubscribers"] = data.EnableSubscribers.ValueBool()
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
    if !data.CopyrightText.IsUnknown() && !state.CopyrightText.IsUnknown() && !data.CopyrightText.Equal(state.CopyrightText) {
        requestDataMap["copyrightText"] = data.CopyrightText.ValueString()
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
        var reportstartdatetimeData interface{}
        if err := json.Unmarshal([]byte(data.ReportStartDateTime.ValueString()), &reportstartdatetimeData); err == nil {
            requestDataMap["reportStartDateTime"] = reportstartdatetimeData
        } else {
            requestDataMap["reportStartDateTime"] = data.ReportStartDateTime.ValueString()
        }
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
        var sendnextreportbyData interface{}
        if err := json.Unmarshal([]byte(data.SendNextReportBy.ValueString()), &sendnextreportbyData); err == nil {
            requestDataMap["sendNextReportBy"] = sendnextreportbyData
        } else {
            requestDataMap["sendNextReportBy"] = data.SendNextReportBy.ValueString()
        }
    }
    if !data.ReportDataInDays.IsUnknown() && !state.ReportDataInDays.IsUnknown() && !data.ReportDataInDays.Equal(state.ReportDataInDays) {
        requestDataMap["reportDataInDays"] = r.bigFloatToFloat64(data.ReportDataInDays.ValueBigFloat())
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

    // Make API call
    httpResp, err := r.client.Put("/status-page/" + data.Id.ValueString() + "", statusPageRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status_page, got error: %s", err))
        return
    }

    // Parse the update response
    var statusPageResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "pageTitle": true,
        "pageDescription": true,
        "description": true,
        "labels": true,
        "faviconFileId": true,
        "logoFileId": true,
        "coverImageFileId": true,
        "headerHTML": true,
        "footerHTML": true,
        "customCSS": true,
        "customJavaScript": true,
        "isPublicStatusPage": true,
        "enableMasterPassword": true,
        "masterPassword": true,
        "showIncidentLabelsOnStatusPage": true,
        "showScheduledEventLabelsOnStatusPage": true,
        "enableSubscribers": true,
        "enableEmailSubscribers": true,
        "allowSubscribersToChooseResources": true,
        "allowSubscribersToChooseEventTypes": true,
        "enableSmsSubscribers": true,
        "enableSlackSubscribers": true,
        "enableMicrosoftTeamsSubscribers": true,
        "copyrightText": true,
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
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "downtimeMonitorStatuses": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/status-page/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page read response, got error: %s", err))
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["pageTitle"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["pageDescription"].(string); ok && val != "" {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
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
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["faviconFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["logoFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["coverImageFileId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["headerHTML"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["footerHTML"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["customCSS"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["customJavaScript"].(string); ok && val != "" {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
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
    } else if val, ok := dataMap["masterPassword"].(string); ok && val != "" {
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
    if val, ok := dataMap["enableSubscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
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
    } else if val, ok := dataMap["copyrightText"].(string); ok && val != "" {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok && val != "" {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
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
    } else if val, ok := dataMap["smtpConfigId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["callSmsConfigId"].(string); ok && val != "" {
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
    } else if dataMap["showIncidentHistoryInDays"] == nil {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showAnnouncementHistoryInDays"].(int64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["showAnnouncementHistoryInDays"] == nil {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["showScheduledEventHistoryInDays"].(int64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["showScheduledEventHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["overviewPageDescription"].(string); ok && val != "" {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    }
    if obj, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultBarColor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultBarColor = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultBarColor = types.StringValue(string(jsonBytes))
            } else {
                data.DefaultBarColor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultBarColor = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultBarColor = types.StringNull()
        }
    } else if val, ok := dataMap["defaultBarColor"].(string); ok && val != "" {
        data.DefaultBarColor = types.StringValue(val)
    } else {
        data.DefaultBarColor = types.StringNull()
    }
    if obj, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberTimezones = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberTimezones = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberTimezones = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberTimezones = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberTimezones = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberTimezones = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberTimezones"].(string); ok && val != "" {
        data.SubscriberTimezones = types.StringValue(val)
    } else {
        data.SubscriberTimezones = types.StringNull()
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportStartDateTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportStartDateTime = types.StringValue(string(jsonBytes))
            } else {
                data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportStartDateTime = types.StringValue(string(jsonBytes))
            } else {
                data.ReportStartDateTime = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportStartDateTime = types.StringValue(string(jsonBytes))
        } else {
            data.ReportStartDateTime = types.StringNull()
        }
    } else if val, ok := dataMap["reportStartDateTime"].(string); ok && val != "" {
        data.ReportStartDateTime = types.StringValue(val)
    } else {
        data.ReportStartDateTime = types.StringNull()
    }
    if obj, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ReportRecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.ReportRecurringInterval = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ReportRecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.ReportRecurringInterval = types.StringNull()
        }
    } else if val, ok := dataMap["reportRecurringInterval"].(string); ok && val != "" {
        data.ReportRecurringInterval = types.StringValue(val)
    } else {
        data.ReportRecurringInterval = types.StringNull()
    }
    if obj, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SendNextReportBy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SendNextReportBy = types.StringValue(string(jsonBytes))
            } else {
                data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SendNextReportBy = types.StringValue(string(jsonBytes))
            } else {
                data.SendNextReportBy = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SendNextReportBy = types.StringValue(string(jsonBytes))
        } else {
            data.SendNextReportBy = types.StringNull()
        }
    } else if val, ok := dataMap["sendNextReportBy"].(string); ok && val != "" {
        data.SendNextReportBy = types.StringValue(val)
    } else {
        data.SendNextReportBy = types.StringNull()
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["reportDataInDays"].(int); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["reportDataInDays"].(int64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["reportDataInDays"] == nil {
        data.ReportDataInDays = types.NumberNull()
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
    } else if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok && val != "" {
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
    } else if dataMap["showEpisodeHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["ipWhitelist"].(string); ok && val != "" {
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
    } else if dataMap["showUptimeHistoryInDays"] == nil {
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
    } else if val, ok := dataMap["embeddedOverallStatusToken"].(string); ok && val != "" {
        data.EmbeddedOverallStatusToken = types.StringValue(val)
    } else {
        data.EmbeddedOverallStatusToken = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
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
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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

func (r *StatusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/status-page/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status_page, got error: %s", err))
        return
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
func (r *StatusPageResource) parseJSONField(terraformString types.String) interface{} {
    if terraformString.IsNull() || terraformString.IsUnknown() || terraformString.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(terraformString.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return terraformString.ValueString()
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

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *StatusPageResource) isValidOneUptimeObjectType(typeStr string) bool {
    validTypes := map[string]bool{
        "ObjectID": true,
        "Decimal": true,
        "Name": true,
        "EqualTo": true,
        "EqualToOrNull": true,
        "MonitorSteps": true,
        "MonitorStep": true,
        "Recurring": true,
        "RestrictionTimes": true,
        "MonitorCriteria": true,
        "PositiveNumber": true,
        "MonitorCriteriaInstance": true,
        "NotEqual": true,
        "Email": true,
        "Phone": true,
        "Color": true,
        "Domain": true,
        "Version": true,
        "IP": true,
        "Route": true,
        "URL": true,
        "Permission": true,
        "Search": true,
        "GreaterThan": true,
        "GreaterThanOrEqual": true,
        "GreaterThanOrNull": true,
        "LessThanOrNull": true,
        "LessThan": true,
        "LessThanOrEqual": true,
        "Port": true,
        "Hostname": true,
        "HashedString": true,
        "DateTime": true,
        "Buffer": true,
        "InBetween": true,
        "NotNull": true,
        "IsNull": true,
        "Includes": true,
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
