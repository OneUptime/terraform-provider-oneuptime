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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
    Labels types.List `tfsdk:"labels"`
    FaviconFileId types.String `tfsdk:"favicon_file_id"`
    LogoFileId types.String `tfsdk:"logo_file_id"`
    CoverImageFileId types.String `tfsdk:"cover_image_file_id"`
    HeaderHTML types.String `tfsdk:"header_h_t_m_l"`
    FooterHTML types.String `tfsdk:"footer_h_t_m_l"`
    CustomCSS types.String `tfsdk:"custom_c_s_s"`
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
    CustomFields types.Map `tfsdk:"custom_fields"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    SmtpConfigId types.String `tfsdk:"smtp_config_id"`
    CallSmsConfigId types.String `tfsdk:"call_sms_config_id"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    ShowIncidentHistoryInDays types.Number `tfsdk:"show_incident_history_in_days"`
    ShowAnnouncementHistoryInDays types.Number `tfsdk:"show_announcement_history_in_days"`
    ShowScheduledEventHistoryInDays types.Number `tfsdk:"show_scheduled_event_history_in_days"`
    OverviewPageDescription types.String `tfsdk:"overview_page_description"`
    HidePoweredByOneUptimeBranding types.Bool `tfsdk:"hide_powered_by_one_uptime_branding"`
    DefaultBarColor types.Map `tfsdk:"default_bar_color"`
    DowntimeMonitorStatuses types.List `tfsdk:"downtime_monitor_statuses"`
    SubscriberTimezones types.Map `tfsdk:"subscriber_timezones"`
    IsReportEnabled types.Bool `tfsdk:"is_report_enabled"`
    ReportStartDateTime types.Map `tfsdk:"report_start_date_time"`
    ReportRecurringInterval types.Map `tfsdk:"report_recurring_interval"`
    SendNextReportBy types.Map `tfsdk:"send_next_report_by"`
    ReportDataInDays types.Number `tfsdk:"report_data_in_days"`
    ShowOverallUptimePercentOnStatusPage types.Bool `tfsdk:"show_overall_uptime_percent_on_status_page"`
    OverallUptimePercentPrecision types.String `tfsdk:"overall_uptime_percent_precision"`
    SubscriberEmailNotificationFooterText types.String `tfsdk:"subscriber_email_notification_footer_text"`
    ShowIncidentsOnStatusPage types.Bool `tfsdk:"show_incidents_on_status_page"`
    ShowAnnouncementsOnStatusPage types.Bool `tfsdk:"show_announcements_on_status_page"`
    ShowScheduledMaintenanceEventsOnStatusPage types.Bool `tfsdk:"show_scheduled_maintenance_events_on_status_page"`
    ShowSubscriberPageOnStatusPage types.Bool `tfsdk:"show_subscriber_page_on_status_page"`
    IpWhitelist types.String `tfsdk:"ip_whitelist"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
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
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name",
                Required: true,
            },
            "page_title": schema.StringAttribute{
                MarkdownDescription: "Page Title",
                Optional: true,
            },
            "page_description": schema.StringAttribute{
                MarkdownDescription: "Page Description",
                Optional: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description",
                Optional: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Labels",
                Optional: true,
                ElementType: types.StringType,
            },
            "favicon_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "logo_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "cover_image_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "header_h_t_m_l": schema.StringAttribute{
                MarkdownDescription: "Header HTML",
                Optional: true,
            },
            "footer_h_t_m_l": schema.StringAttribute{
                MarkdownDescription: "Footer HTML",
                Optional: true,
            },
            "custom_c_s_s": schema.StringAttribute{
                MarkdownDescription: "CSS",
                Optional: true,
            },
            "custom_java_script": schema.StringAttribute{
                MarkdownDescription: "JavaScript",
                Optional: true,
            },
            "is_public_status_page": schema.BoolAttribute{
                MarkdownDescription: "Public Status Page",
                Optional: true,
            },
            "show_incident_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incident Labels",
                Optional: true,
            },
            "show_scheduled_event_labels_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Event Labels",
                Optional: true,
            },
            "enable_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Enable Subscribers",
                Optional: true,
            },
            "enable_email_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Enable Email Subscribers",
                Optional: true,
            },
            "allow_subscribers_to_choose_resources": schema.BoolAttribute{
                MarkdownDescription: "Allow Subscribers to Choose Resources",
                Optional: true,
            },
            "allow_subscribers_to_choose_event_types": schema.BoolAttribute{
                MarkdownDescription: "Allow Subscribers to subscribe to event types",
                Optional: true,
            },
            "enable_sms_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS Subscribers",
                Optional: true,
            },
            "enable_slack_subscribers": schema.BoolAttribute{
                MarkdownDescription: "Enable Slack Subscribers",
                Optional: true,
            },
            "copyright_text": schema.StringAttribute{
                MarkdownDescription: "Copyright Text",
                Optional: true,
            },
            "custom_fields": schema.MapAttribute{
                MarkdownDescription: "Custom Fields",
                Optional: true,
                ElementType: types.StringType,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Require SSO",
                Optional: true,
            },
            "smtp_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are Owners Notified Of Resource Creation?",
                Optional: true,
            },
            "show_incident_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "Show incident history in days",
                Optional: true,
            },
            "show_announcement_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "Show announcement history in days",
                Optional: true,
            },
            "show_scheduled_event_history_in_days": schema.NumberAttribute{
                MarkdownDescription: "Show scheduled event history in days",
                Optional: true,
            },
            "overview_page_description": schema.StringAttribute{
                MarkdownDescription: "Overview Page Description",
                Optional: true,
            },
            "hide_powered_by_one_uptime_branding": schema.BoolAttribute{
                MarkdownDescription: "Hide Powered By OneUptime Branding",
                Optional: true,
            },
            "default_bar_color": schema.MapAttribute{
                MarkdownDescription: "Color object",
                Optional: true,
                ElementType: types.StringType,
            },
            "downtime_monitor_statuses": schema.ListAttribute{
                MarkdownDescription: "Downtime Monitor Statuses",
                Optional: true,
                ElementType: types.StringType,
            },
            "subscriber_timezones": schema.MapAttribute{
                MarkdownDescription: "Subscriber Timezones",
                Optional: true,
                ElementType: types.StringType,
            },
            "is_report_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Report Enabled",
                Optional: true,
            },
            "report_start_date_time": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "report_recurring_interval": schema.MapAttribute{
                MarkdownDescription: "Report Recurring Interval",
                Optional: true,
                ElementType: types.StringType,
            },
            "send_next_report_by": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "report_data_in_days": schema.NumberAttribute{
                MarkdownDescription: "Report data for the last N days",
                Optional: true,
            },
            "show_overall_uptime_percent_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Overall Uptime Percent on Status Page",
                Optional: true,
            },
            "overall_uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Overall Uptime Percent Precision",
                Optional: true,
            },
            "subscriber_email_notification_footer_text": schema.StringAttribute{
                MarkdownDescription: "Subscriber Email Notification Footer Text",
                Optional: true,
            },
            "show_incidents_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Incidents on Status Page",
                Optional: true,
            },
            "show_announcements_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Announcements on Status Page",
                Optional: true,
            },
            "show_scheduled_maintenance_events_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Scheduled Maintenance Events on Status Page",
                Optional: true,
            },
            "show_subscriber_page_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Show Subscriber Page on Status Page",
                Optional: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "IP Whitelist",
                Optional: true,
            },
            "created_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "updated_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "deleted_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
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
        "projectId": data.ProjectId.ValueString(),
        "name": data.Name.ValueString(),
        "pageTitle": data.PageTitle.ValueString(),
        "pageDescription": data.PageDescription.ValueString(),
        "description": data.Description.ValueString(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "faviconFileId": data.FaviconFileId.ValueString(),
        "logoFileId": data.LogoFileId.ValueString(),
        "coverImageFileId": data.CoverImageFileId.ValueString(),
        "headerHTML": data.HeaderHTML.ValueString(),
        "footerHTML": data.FooterHTML.ValueString(),
        "customCSS": data.CustomCSS.ValueString(),
        "customJavaScript": data.CustomJavaScript.ValueString(),
        "isPublicStatusPage": data.IsPublicStatusPage.ValueBool(),
        "showIncidentLabelsOnStatusPage": data.ShowIncidentLabelsOnStatusPage.ValueBool(),
        "showScheduledEventLabelsOnStatusPage": data.ShowScheduledEventLabelsOnStatusPage.ValueBool(),
        "enableSubscribers": data.EnableSubscribers.ValueBool(),
        "enableEmailSubscribers": data.EnableEmailSubscribers.ValueBool(),
        "allowSubscribersToChooseResources": data.AllowSubscribersToChooseResources.ValueBool(),
        "allowSubscribersToChooseEventTypes": data.AllowSubscribersToChooseEventTypes.ValueBool(),
        "enableSmsSubscribers": data.EnableSmsSubscribers.ValueBool(),
        "enableSlackSubscribers": data.EnableSlackSubscribers.ValueBool(),
        "copyrightText": data.CopyrightText.ValueString(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "requireSsoForLogin": data.RequireSsoForLogin.ValueBool(),
        "smtpConfigId": data.SmtpConfigId.ValueString(),
        "callSmsConfigId": data.CallSmsConfigId.ValueString(),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "showIncidentHistoryInDays": data.ShowIncidentHistoryInDays.ValueBigFloat(),
        "showAnnouncementHistoryInDays": data.ShowAnnouncementHistoryInDays.ValueBigFloat(),
        "showScheduledEventHistoryInDays": data.ShowScheduledEventHistoryInDays.ValueBigFloat(),
        "overviewPageDescription": data.OverviewPageDescription.ValueString(),
        "hidePoweredByOneUptimeBranding": data.HidePoweredByOneUptimeBranding.ValueBool(),
        "defaultBarColor": r.convertTerraformMapToInterface(data.DefaultBarColor),
        "downtimeMonitorStatuses": r.convertTerraformListToInterface(data.DowntimeMonitorStatuses),
        "subscriberTimezones": r.convertTerraformMapToInterface(data.SubscriberTimezones),
        "isReportEnabled": data.IsReportEnabled.ValueBool(),
        "reportStartDateTime": r.convertTerraformMapToInterface(data.ReportStartDateTime),
        "reportRecurringInterval": r.convertTerraformMapToInterface(data.ReportRecurringInterval),
        "sendNextReportBy": r.convertTerraformMapToInterface(data.SendNextReportBy),
        "reportDataInDays": data.ReportDataInDays.ValueBigFloat(),
        "showOverallUptimePercentOnStatusPage": data.ShowOverallUptimePercentOnStatusPage.ValueBool(),
        "overallUptimePercentPrecision": data.OverallUptimePercentPrecision.ValueString(),
        "subscriberEmailNotificationFooterText": data.SubscriberEmailNotificationFooterText.ValueString(),
        "showIncidentsOnStatusPage": data.ShowIncidentsOnStatusPage.ValueBool(),
        "showAnnouncementsOnStatusPage": data.ShowAnnouncementsOnStatusPage.ValueBool(),
        "showScheduledMaintenanceEventsOnStatusPage": data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool(),
        "showSubscriberPageOnStatusPage": data.ShowSubscriberPageOnStatusPage.ValueBool(),
        "ipWhitelist": data.IpWhitelist.ValueString(),
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["pageTitle"].(string); ok && val != "" {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if val, ok := dataMap["pageDescription"].(string); ok && val != "" {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["faviconFileId"].(string); ok && val != "" {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if val, ok := dataMap["logoFileId"].(string); ok && val != "" {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if val, ok := dataMap["coverImageFileId"].(string); ok && val != "" {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if val, ok := dataMap["headerHTML"].(string); ok && val != "" {
        data.HeaderHTML = types.StringValue(val)
    } else {
        data.HeaderHTML = types.StringNull()
    }
    if val, ok := dataMap["footerHTML"].(string); ok && val != "" {
        data.FooterHTML = types.StringValue(val)
    } else {
        data.FooterHTML = types.StringNull()
    }
    if val, ok := dataMap["customCSS"].(string); ok && val != "" {
        data.CustomCSS = types.StringValue(val)
    } else {
        data.CustomCSS = types.StringNull()
    }
    if val, ok := dataMap["customJavaScript"].(string); ok && val != "" {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    } else if dataMap["isPublicStatusPage"] == nil {
        data.IsPublicStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showIncidentLabelsOnStatusPage"] == nil {
        data.ShowIncidentLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showScheduledEventLabelsOnStatusPage"] == nil {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["enableSubscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
    } else if dataMap["enableSubscribers"] == nil {
        data.EnableSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    } else if dataMap["enableEmailSubscribers"] == nil {
        data.EnableEmailSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    } else if dataMap["allowSubscribersToChooseResources"] == nil {
        data.AllowSubscribersToChooseResources = types.BoolNull()
    }
    if val, ok := dataMap["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    } else if dataMap["allowSubscribersToChooseEventTypes"] == nil {
        data.AllowSubscribersToChooseEventTypes = types.BoolNull()
    }
    if val, ok := dataMap["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    } else if dataMap["enableSmsSubscribers"] == nil {
        data.EnableSmsSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    } else if dataMap["enableSlackSubscribers"] == nil {
        data.EnableSlackSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["copyrightText"].(string); ok && val != "" {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else if dataMap["requireSsoForLogin"] == nil {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if val, ok := dataMap["smtpConfigId"].(string); ok && val != "" {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if val, ok := dataMap["callSmsConfigId"].(string); ok && val != "" {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showIncidentHistoryInDays"] == nil {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showAnnouncementHistoryInDays"] == nil {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showScheduledEventHistoryInDays"] == nil {
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["overviewPageDescription"].(string); ok && val != "" {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    } else if dataMap["hidePoweredByOneUptimeBranding"] == nil {
        data.HidePoweredByOneUptimeBranding = types.BoolNull()
    }
    if val, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DefaultBarColor = mapValue
    } else if dataMap["defaultBarColor"] == nil {
        data.DefaultBarColor = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.DowntimeMonitorStatuses = listValue
    } else if dataMap["downtimeMonitorStatuses"] == nil {
        data.DowntimeMonitorStatuses = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberTimezones = mapValue
    } else if dataMap["subscriberTimezones"] == nil {
        data.SubscriberTimezones = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    } else if dataMap["isReportEnabled"] == nil {
        data.IsReportEnabled = types.BoolNull()
    }
    if val, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ReportStartDateTime = mapValue
    } else if dataMap["reportStartDateTime"] == nil {
        data.ReportStartDateTime = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ReportRecurringInterval = mapValue
    } else if dataMap["reportRecurringInterval"] == nil {
        data.ReportRecurringInterval = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SendNextReportBy = mapValue
    } else if dataMap["sendNextReportBy"] == nil {
        data.SendNextReportBy = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["reportDataInDays"] == nil {
        data.ReportDataInDays = types.NumberNull()
    }
    if val, ok := dataMap["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    } else if dataMap["showOverallUptimePercentOnStatusPage"] == nil {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok && val != "" {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok && val != "" {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := dataMap["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showIncidentsOnStatusPage"] == nil {
        data.ShowIncidentsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showAnnouncementsOnStatusPage"] == nil {
        data.ShowAnnouncementsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showScheduledMaintenanceEventsOnStatusPage"] == nil {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    } else if dataMap["showSubscriberPageOnStatusPage"] == nil {
        data.ShowSubscriberPageOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["ipWhitelist"].(string); ok && val != "" {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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
        "showIncidentLabelsOnStatusPage": true,
        "showScheduledEventLabelsOnStatusPage": true,
        "enableSubscribers": true,
        "enableEmailSubscribers": true,
        "allowSubscribersToChooseResources": true,
        "allowSubscribersToChooseEventTypes": true,
        "enableSmsSubscribers": true,
        "enableSlackSubscribers": true,
        "copyrightText": true,
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
        "showOverallUptimePercentOnStatusPage": true,
        "overallUptimePercentPrecision": true,
        "subscriberEmailNotificationFooterText": true,
        "showIncidentsOnStatusPage": true,
        "showAnnouncementsOnStatusPage": true,
        "showScheduledMaintenanceEventsOnStatusPage": true,
        "showSubscriberPageOnStatusPage": true,
        "ipWhitelist": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["pageTitle"].(string); ok && val != "" {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if val, ok := dataMap["pageDescription"].(string); ok && val != "" {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["faviconFileId"].(string); ok && val != "" {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if val, ok := dataMap["logoFileId"].(string); ok && val != "" {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if val, ok := dataMap["coverImageFileId"].(string); ok && val != "" {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if val, ok := dataMap["headerHTML"].(string); ok && val != "" {
        data.HeaderHTML = types.StringValue(val)
    } else {
        data.HeaderHTML = types.StringNull()
    }
    if val, ok := dataMap["footerHTML"].(string); ok && val != "" {
        data.FooterHTML = types.StringValue(val)
    } else {
        data.FooterHTML = types.StringNull()
    }
    if val, ok := dataMap["customCSS"].(string); ok && val != "" {
        data.CustomCSS = types.StringValue(val)
    } else {
        data.CustomCSS = types.StringNull()
    }
    if val, ok := dataMap["customJavaScript"].(string); ok && val != "" {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    } else if dataMap["isPublicStatusPage"] == nil {
        data.IsPublicStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showIncidentLabelsOnStatusPage"] == nil {
        data.ShowIncidentLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showScheduledEventLabelsOnStatusPage"] == nil {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["enableSubscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
    } else if dataMap["enableSubscribers"] == nil {
        data.EnableSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    } else if dataMap["enableEmailSubscribers"] == nil {
        data.EnableEmailSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    } else if dataMap["allowSubscribersToChooseResources"] == nil {
        data.AllowSubscribersToChooseResources = types.BoolNull()
    }
    if val, ok := dataMap["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    } else if dataMap["allowSubscribersToChooseEventTypes"] == nil {
        data.AllowSubscribersToChooseEventTypes = types.BoolNull()
    }
    if val, ok := dataMap["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    } else if dataMap["enableSmsSubscribers"] == nil {
        data.EnableSmsSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    } else if dataMap["enableSlackSubscribers"] == nil {
        data.EnableSlackSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["copyrightText"].(string); ok && val != "" {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else if dataMap["requireSsoForLogin"] == nil {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if val, ok := dataMap["smtpConfigId"].(string); ok && val != "" {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if val, ok := dataMap["callSmsConfigId"].(string); ok && val != "" {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showIncidentHistoryInDays"] == nil {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showAnnouncementHistoryInDays"] == nil {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showScheduledEventHistoryInDays"] == nil {
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["overviewPageDescription"].(string); ok && val != "" {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    } else if dataMap["hidePoweredByOneUptimeBranding"] == nil {
        data.HidePoweredByOneUptimeBranding = types.BoolNull()
    }
    if val, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DefaultBarColor = mapValue
    } else if dataMap["defaultBarColor"] == nil {
        data.DefaultBarColor = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.DowntimeMonitorStatuses = listValue
    } else if dataMap["downtimeMonitorStatuses"] == nil {
        data.DowntimeMonitorStatuses = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberTimezones = mapValue
    } else if dataMap["subscriberTimezones"] == nil {
        data.SubscriberTimezones = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    } else if dataMap["isReportEnabled"] == nil {
        data.IsReportEnabled = types.BoolNull()
    }
    if val, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ReportStartDateTime = mapValue
    } else if dataMap["reportStartDateTime"] == nil {
        data.ReportStartDateTime = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ReportRecurringInterval = mapValue
    } else if dataMap["reportRecurringInterval"] == nil {
        data.ReportRecurringInterval = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SendNextReportBy = mapValue
    } else if dataMap["sendNextReportBy"] == nil {
        data.SendNextReportBy = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["reportDataInDays"] == nil {
        data.ReportDataInDays = types.NumberNull()
    }
    if val, ok := dataMap["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    } else if dataMap["showOverallUptimePercentOnStatusPage"] == nil {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok && val != "" {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok && val != "" {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := dataMap["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showIncidentsOnStatusPage"] == nil {
        data.ShowIncidentsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showAnnouncementsOnStatusPage"] == nil {
        data.ShowAnnouncementsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showScheduledMaintenanceEventsOnStatusPage"] == nil {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    } else if dataMap["showSubscriberPageOnStatusPage"] == nil {
        data.ShowSubscriberPageOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["ipWhitelist"].(string); ok && val != "" {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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
        "data": map[string]interface{}{
        "name": data.Name.ValueString(),
        "pageTitle": data.PageTitle.ValueString(),
        "pageDescription": data.PageDescription.ValueString(),
        "description": data.Description.ValueString(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "faviconFileId": data.FaviconFileId.ValueString(),
        "logoFileId": data.LogoFileId.ValueString(),
        "coverImageFileId": data.CoverImageFileId.ValueString(),
        "headerHTML": data.HeaderHTML.ValueString(),
        "footerHTML": data.FooterHTML.ValueString(),
        "customCSS": data.CustomCSS.ValueString(),
        "customJavaScript": data.CustomJavaScript.ValueString(),
        "isPublicStatusPage": data.IsPublicStatusPage.ValueBool(),
        "showIncidentLabelsOnStatusPage": data.ShowIncidentLabelsOnStatusPage.ValueBool(),
        "showScheduledEventLabelsOnStatusPage": data.ShowScheduledEventLabelsOnStatusPage.ValueBool(),
        "enableSubscribers": data.EnableSubscribers.ValueBool(),
        "enableEmailSubscribers": data.EnableEmailSubscribers.ValueBool(),
        "allowSubscribersToChooseResources": data.AllowSubscribersToChooseResources.ValueBool(),
        "allowSubscribersToChooseEventTypes": data.AllowSubscribersToChooseEventTypes.ValueBool(),
        "enableSmsSubscribers": data.EnableSmsSubscribers.ValueBool(),
        "enableSlackSubscribers": data.EnableSlackSubscribers.ValueBool(),
        "copyrightText": data.CopyrightText.ValueString(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "requireSsoForLogin": data.RequireSsoForLogin.ValueBool(),
        "smtpConfigId": data.SmtpConfigId.ValueString(),
        "callSmsConfigId": data.CallSmsConfigId.ValueString(),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "showIncidentHistoryInDays": data.ShowIncidentHistoryInDays.ValueBigFloat(),
        "showAnnouncementHistoryInDays": data.ShowAnnouncementHistoryInDays.ValueBigFloat(),
        "showScheduledEventHistoryInDays": data.ShowScheduledEventHistoryInDays.ValueBigFloat(),
        "overviewPageDescription": data.OverviewPageDescription.ValueString(),
        "hidePoweredByOneUptimeBranding": data.HidePoweredByOneUptimeBranding.ValueBool(),
        "defaultBarColor": r.convertTerraformMapToInterface(data.DefaultBarColor),
        "downtimeMonitorStatuses": r.convertTerraformListToInterface(data.DowntimeMonitorStatuses),
        "subscriberTimezones": r.convertTerraformMapToInterface(data.SubscriberTimezones),
        "isReportEnabled": data.IsReportEnabled.ValueBool(),
        "reportStartDateTime": r.convertTerraformMapToInterface(data.ReportStartDateTime),
        "reportRecurringInterval": r.convertTerraformMapToInterface(data.ReportRecurringInterval),
        "sendNextReportBy": r.convertTerraformMapToInterface(data.SendNextReportBy),
        "reportDataInDays": data.ReportDataInDays.ValueBigFloat(),
        "showOverallUptimePercentOnStatusPage": data.ShowOverallUptimePercentOnStatusPage.ValueBool(),
        "overallUptimePercentPrecision": data.OverallUptimePercentPrecision.ValueString(),
        "subscriberEmailNotificationFooterText": data.SubscriberEmailNotificationFooterText.ValueString(),
        "showIncidentsOnStatusPage": data.ShowIncidentsOnStatusPage.ValueBool(),
        "showAnnouncementsOnStatusPage": data.ShowAnnouncementsOnStatusPage.ValueBool(),
        "showScheduledMaintenanceEventsOnStatusPage": data.ShowScheduledMaintenanceEventsOnStatusPage.ValueBool(),
        "showSubscriberPageOnStatusPage": data.ShowSubscriberPageOnStatusPage.ValueBool(),
        "ipWhitelist": data.IpWhitelist.ValueString(),
        },
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
        "showIncidentLabelsOnStatusPage": true,
        "showScheduledEventLabelsOnStatusPage": true,
        "enableSubscribers": true,
        "enableEmailSubscribers": true,
        "allowSubscribersToChooseResources": true,
        "allowSubscribersToChooseEventTypes": true,
        "enableSmsSubscribers": true,
        "enableSlackSubscribers": true,
        "copyrightText": true,
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
        "showOverallUptimePercentOnStatusPage": true,
        "overallUptimePercentPrecision": true,
        "subscriberEmailNotificationFooterText": true,
        "showIncidentsOnStatusPage": true,
        "showAnnouncementsOnStatusPage": true,
        "showScheduledMaintenanceEventsOnStatusPage": true,
        "showSubscriberPageOnStatusPage": true,
        "ipWhitelist": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["pageTitle"].(string); ok && val != "" {
        data.PageTitle = types.StringValue(val)
    } else {
        data.PageTitle = types.StringNull()
    }
    if val, ok := dataMap["pageDescription"].(string); ok && val != "" {
        data.PageDescription = types.StringValue(val)
    } else {
        data.PageDescription = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["faviconFileId"].(string); ok && val != "" {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
    }
    if val, ok := dataMap["logoFileId"].(string); ok && val != "" {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if val, ok := dataMap["coverImageFileId"].(string); ok && val != "" {
        data.CoverImageFileId = types.StringValue(val)
    } else {
        data.CoverImageFileId = types.StringNull()
    }
    if val, ok := dataMap["headerHTML"].(string); ok && val != "" {
        data.HeaderHTML = types.StringValue(val)
    } else {
        data.HeaderHTML = types.StringNull()
    }
    if val, ok := dataMap["footerHTML"].(string); ok && val != "" {
        data.FooterHTML = types.StringValue(val)
    } else {
        data.FooterHTML = types.StringNull()
    }
    if val, ok := dataMap["customCSS"].(string); ok && val != "" {
        data.CustomCSS = types.StringValue(val)
    } else {
        data.CustomCSS = types.StringNull()
    }
    if val, ok := dataMap["customJavaScript"].(string); ok && val != "" {
        data.CustomJavaScript = types.StringValue(val)
    } else {
        data.CustomJavaScript = types.StringNull()
    }
    if val, ok := dataMap["isPublicStatusPage"].(bool); ok {
        data.IsPublicStatusPage = types.BoolValue(val)
    } else if dataMap["isPublicStatusPage"] == nil {
        data.IsPublicStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showIncidentLabelsOnStatusPage"].(bool); ok {
        data.ShowIncidentLabelsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showIncidentLabelsOnStatusPage"] == nil {
        data.ShowIncidentLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showScheduledEventLabelsOnStatusPage"].(bool); ok {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showScheduledEventLabelsOnStatusPage"] == nil {
        data.ShowScheduledEventLabelsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["enableSubscribers"].(bool); ok {
        data.EnableSubscribers = types.BoolValue(val)
    } else if dataMap["enableSubscribers"] == nil {
        data.EnableSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["enableEmailSubscribers"].(bool); ok {
        data.EnableEmailSubscribers = types.BoolValue(val)
    } else if dataMap["enableEmailSubscribers"] == nil {
        data.EnableEmailSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["allowSubscribersToChooseResources"].(bool); ok {
        data.AllowSubscribersToChooseResources = types.BoolValue(val)
    } else if dataMap["allowSubscribersToChooseResources"] == nil {
        data.AllowSubscribersToChooseResources = types.BoolNull()
    }
    if val, ok := dataMap["allowSubscribersToChooseEventTypes"].(bool); ok {
        data.AllowSubscribersToChooseEventTypes = types.BoolValue(val)
    } else if dataMap["allowSubscribersToChooseEventTypes"] == nil {
        data.AllowSubscribersToChooseEventTypes = types.BoolNull()
    }
    if val, ok := dataMap["enableSmsSubscribers"].(bool); ok {
        data.EnableSmsSubscribers = types.BoolValue(val)
    } else if dataMap["enableSmsSubscribers"] == nil {
        data.EnableSmsSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["enableSlackSubscribers"].(bool); ok {
        data.EnableSlackSubscribers = types.BoolValue(val)
    } else if dataMap["enableSlackSubscribers"] == nil {
        data.EnableSlackSubscribers = types.BoolNull()
    }
    if val, ok := dataMap["copyrightText"].(string); ok && val != "" {
        data.CopyrightText = types.StringValue(val)
    } else {
        data.CopyrightText = types.StringNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else if dataMap["requireSsoForLogin"] == nil {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if val, ok := dataMap["smtpConfigId"].(string); ok && val != "" {
        data.SmtpConfigId = types.StringValue(val)
    } else {
        data.SmtpConfigId = types.StringNull()
    }
    if val, ok := dataMap["callSmsConfigId"].(string); ok && val != "" {
        data.CallSmsConfigId = types.StringValue(val)
    } else {
        data.CallSmsConfigId = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["showIncidentHistoryInDays"].(float64); ok {
        data.ShowIncidentHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showIncidentHistoryInDays"] == nil {
        data.ShowIncidentHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showAnnouncementHistoryInDays"].(float64); ok {
        data.ShowAnnouncementHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showAnnouncementHistoryInDays"] == nil {
        data.ShowAnnouncementHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["showScheduledEventHistoryInDays"].(float64); ok {
        data.ShowScheduledEventHistoryInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["showScheduledEventHistoryInDays"] == nil {
        data.ShowScheduledEventHistoryInDays = types.NumberNull()
    }
    if val, ok := dataMap["overviewPageDescription"].(string); ok && val != "" {
        data.OverviewPageDescription = types.StringValue(val)
    } else {
        data.OverviewPageDescription = types.StringNull()
    }
    if val, ok := dataMap["hidePoweredByOneUptimeBranding"].(bool); ok {
        data.HidePoweredByOneUptimeBranding = types.BoolValue(val)
    } else if dataMap["hidePoweredByOneUptimeBranding"] == nil {
        data.HidePoweredByOneUptimeBranding = types.BoolNull()
    }
    if val, ok := dataMap["defaultBarColor"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DefaultBarColor = mapValue
    } else if dataMap["defaultBarColor"] == nil {
        data.DefaultBarColor = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["downtimeMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.DowntimeMonitorStatuses = listValue
    } else if dataMap["downtimeMonitorStatuses"] == nil {
        data.DowntimeMonitorStatuses = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["subscriberTimezones"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberTimezones = mapValue
    } else if dataMap["subscriberTimezones"] == nil {
        data.SubscriberTimezones = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isReportEnabled"].(bool); ok {
        data.IsReportEnabled = types.BoolValue(val)
    } else if dataMap["isReportEnabled"] == nil {
        data.IsReportEnabled = types.BoolNull()
    }
    if val, ok := dataMap["reportStartDateTime"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ReportStartDateTime = mapValue
    } else if dataMap["reportStartDateTime"] == nil {
        data.ReportStartDateTime = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["reportRecurringInterval"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ReportRecurringInterval = mapValue
    } else if dataMap["reportRecurringInterval"] == nil {
        data.ReportRecurringInterval = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["sendNextReportBy"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SendNextReportBy = mapValue
    } else if dataMap["sendNextReportBy"] == nil {
        data.SendNextReportBy = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["reportDataInDays"].(float64); ok {
        data.ReportDataInDays = types.NumberValue(big.NewFloat(val))
    } else if dataMap["reportDataInDays"] == nil {
        data.ReportDataInDays = types.NumberNull()
    }
    if val, ok := dataMap["showOverallUptimePercentOnStatusPage"].(bool); ok {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolValue(val)
    } else if dataMap["showOverallUptimePercentOnStatusPage"] == nil {
        data.ShowOverallUptimePercentOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["overallUptimePercentPrecision"].(string); ok && val != "" {
        data.OverallUptimePercentPrecision = types.StringValue(val)
    } else {
        data.OverallUptimePercentPrecision = types.StringNull()
    }
    if val, ok := dataMap["subscriberEmailNotificationFooterText"].(string); ok && val != "" {
        data.SubscriberEmailNotificationFooterText = types.StringValue(val)
    } else {
        data.SubscriberEmailNotificationFooterText = types.StringNull()
    }
    if val, ok := dataMap["showIncidentsOnStatusPage"].(bool); ok {
        data.ShowIncidentsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showIncidentsOnStatusPage"] == nil {
        data.ShowIncidentsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showAnnouncementsOnStatusPage"].(bool); ok {
        data.ShowAnnouncementsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showAnnouncementsOnStatusPage"] == nil {
        data.ShowAnnouncementsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showScheduledMaintenanceEventsOnStatusPage"].(bool); ok {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolValue(val)
    } else if dataMap["showScheduledMaintenanceEventsOnStatusPage"] == nil {
        data.ShowScheduledMaintenanceEventsOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["showSubscriberPageOnStatusPage"].(bool); ok {
        data.ShowSubscriberPageOnStatusPage = types.BoolValue(val)
    } else if dataMap["showSubscriberPageOnStatusPage"] == nil {
        data.ShowSubscriberPageOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["ipWhitelist"].(string); ok && val != "" {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
