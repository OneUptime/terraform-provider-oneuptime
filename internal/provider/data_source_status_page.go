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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    PageTitle types.String `tfsdk:"page_title"`
    PageDescription types.String `tfsdk:"page_description"`
    EnableSearchEngineIndexing types.Bool `tfsdk:"enable_search_engine_indexing"`
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
        MarkdownDescription: "Manage status pages for your project. Look up an existing status_page by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
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
                MarkdownDescription: "Title of your Status Page. This is used for SEO..",
                Computed: true,
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Description of your Status Page. This is used for SEO..",
                Computed: true,
            },
            "enable_search_engine_indexing": schema.BoolAttribute{
                MarkdownDescription: "Should search engines like Google and Bing be allowed to index this status page? Turn this off to keep the page reachable by link but out of search results..",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
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
                MarkdownDescription: "Status Page Custom HTML Header. Served only from a verified custom domain..",
                Computed: true,
            },
            "footer_html": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom HTML Footer. Served only from a verified custom domain..",
                Computed: true,
            },
            "custom_css": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom CSS. Served only from a verified custom domain..",
                Computed: true,
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "Status Page Custom JavaScript. This runs when the status page is loaded from a verified custom domain..",
                Computed: true,
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Is this status page public?.",
                Computed: true,
            },
            "enable_mcp_server": schema.BoolAttribute{
                MarkdownDescription: "Can AI agents read this status page over the public OneUptime MCP server? This does not affect the status page website, its RSS feed, or its public JSON API..",
                Computed: true,
            },
            "enable_master_password": schema.BoolAttribute{
                MarkdownDescription: "Require visitors to enter a master password before viewing a private status page..",
                Computed: true,
            },
            "master_password": schema.StringAttribute{
                MarkdownDescription: "Password required to unlock a private status page. This value is stored as a secure hash..",
                Computed: true,
                Sensitive: true,
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Labels on Status Page?.",
                Computed: true,
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Event Labels on Status Page?.",
                Computed: true,
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can email subscribers subscribe to this Status Page?.",
                Computed: true,
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which resources to subscribe to?.",
                Computed: true,
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?.",
                Computed: true,
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can SMS subscribers subscribe to this Status Page?.",
                Computed: true,
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Slack subscribers subscribe to this Status Page?.",
                Computed: true,
            },
            "enable_microsoft_teams_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Microsoft Teams subscribers subscribe to this Status Page?.",
                Computed: true,
            },
            "enable_webhook_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Can Webhook subscribers subscribe to this Status Page?.",
                Computed: true,
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Copyright Text.",
                Computed: true,
            },
            "logo_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the logo image, read by screen readers for accessibility..",
                Computed: true,
            },
            "cover_image_alt_text": schema.StringAttribute{
                MarkdownDescription: "Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative..",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource..",
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Should SSO be required to login to Private Status Page.",
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
                MarkdownDescription: "Are owners notified of when this resource is created?.",
                Computed: true,
            },
            "show_incident_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of incident history should be shown on the status page (in days)?.",
                Computed: true,
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of announcement history should be shown on the status page (in days)?.",
                Computed: true,
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of scheduled event history should be shown on the status page (in days)?.",
                Computed: true,
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Overview Page description for your status page. This is a markdown field..",
                Computed: true,
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Hide Powered By OneUptime Branding?.",
                Computed: true,
            },
            "default_bar_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "downtime_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "List of monitors statuses that are considered as \"down\" for this status page..",
                Computed: true,
                ElementType: types.StringType,
            },
            "subscriber_timezones": schema.StringAttribute{
                MarkdownDescription: "Timezones of subscribers to this status page..",
                Computed: true,
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Report Enabled for this Status Page?.",
                Computed: true,
            },
            "report_start_date_time": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "report_recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often would you like to send the report?.",
                Computed: true,
            },
            "send_next_report_by": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of data should be included in the report?.",
                Computed: true,
            },
            "report_period_type": schema.StringAttribute{
                MarkdownDescription: "Should the report cover a rolling number of days, or the previous whole calendar period?.",
                Computed: true,
            },
            "report_timezone": schema.StringAttribute{
                MarkdownDescription: "The timezone report periods and send times are resolved in. A monthly report in this timezone runs from the 1st at 00:00 to the last day at 23:59..",
                Computed: true,
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Overall Uptime Percent on Status Page?.",
                Computed: true,
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Overall Precision of uptime percent for this status page..",
                Computed: true,
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Text to send to subscribers in the footer of the email..",
                Computed: true,
            },
            "enable_custom_subscriber_email_notification_footer_text": schema.BoolAttribute{
                MarkdownDescription: "Enable custom footer text in subscriber email notifications..",
                Computed: true,
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incidents on Status Page?.",
                Computed: true,
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Announcements on Status Page?.",
                Computed: true,
            },
            "show_episodes_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Episodes on Status Page?.",
                Computed: true,
            },
            "show_episode_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of episode history to show on the status page.",
                Computed: true,
            },
            "show_episode_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Episode Labels on Status Page?.",
                Computed: true,
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Maintenance Events on Status Page?.",
                Computed: true,
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Subscriber Page on Status Page?.",
                Computed: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist for this Status Page. One IP per line. Only used if the status page is private..",
                Computed: true,
            },
            "enable_embedded_overall_status": schema.BoolAttribute{
                MarkdownDescription: "Enable embedded overall status badge that can be displayed on external websites?.",
                Computed: true,
            },
            "show_uptime_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "How many days of uptime history should be shown on the status page? Maximum is 90 days..",
                Computed: true,
            },
            "embedded_overall_status_token": schema.StringAttribute{
                MarkdownDescription: "Security token required to access the embedded overall status badge. This token must be provided in the URL..",
                Computed: true,
            },
            "default_language": schema.StringAttribute{
                MarkdownDescription: "Default language that the status page is shown in when a visitor arrives for the first time..",
                Computed: true,
            },
            "enabled_languages": schema.StringAttribute{
                MarkdownDescription: "Languages offered in the footer language switcher. Leave empty to offer all supported languages..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a status_page.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "pageTitle": true,
        "pageDescription": true,
        "enableSearchEngineIndexing": true,
        "description": true,
        "slug": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status_page found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status_page found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status_page matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
