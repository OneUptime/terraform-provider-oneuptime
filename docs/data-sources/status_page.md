---
page_title: "oneuptime_status_page Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage status pages for your project.
---

# oneuptime_status_page (Data Source)

Manage status pages for your project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page" "example" {
  name = "Example status page"
}

# Or by id:
data "oneuptime_status_page" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `allow_subscribers_to_choose_event_types` (Boolean) Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?
- `allow_subscribers_to_choose_resources` (Boolean) Can subscribers choose which resources to subscribe to?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `call_sms_config_id` (String) ID of your Call/SMS Config Resource which is used to send SMS to subscribers.
- `copyright_text` (String) Copyright Text.
- `cover_image_alt_text` (String) Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative.
- `cover_image_file_id` (String) Status Page Cover Image ID. The ID of a `oneuptime_file`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_css` (String) Status Page Custom CSS. Served only from a verified custom domain.
- `custom_java_script` (String) Status Page Custom JavaScript. This runs when the status page is loaded from a verified custom domain.
- `default_language` (String) Default language that the status page is shown in when a visitor arrives for the first time.
- `description` (String) Friendly description that will help you remember.
- `embedded_overall_status_token` (String) Security token required to access the embedded overall status badge. This token must be provided in the URL.
- `enable_custom_subscriber_email_notification_footer_text` (Boolean) Enable custom footer text in subscriber email notifications.
- `enable_email_subscribers` (Boolean) Can email subscribers subscribe to this Status Page?
- `enable_embedded_overall_status` (Boolean) Enable embedded overall status badge that can be displayed on external websites?
- `enable_master_password` (Boolean) Require visitors to enter a master password before viewing a private status page.
- `enable_mcp_server` (Boolean) Can AI agents read this status page over the public OneUptime MCP server? This does not affect the status page website, its RSS feed, or its public JSON API.
- `enable_microsoft_teams_subscribers` (Boolean) Can Microsoft Teams subscribers subscribe to this Status Page?
- `enable_search_engine_indexing` (Boolean) Should search engines like Google and Bing be allowed to index this status page? Turn this off to keep the page reachable by link but out of search results.
- `enable_slack_subscribers` (Boolean) Can Slack subscribers subscribe to this Status Page?
- `enable_sms_subscribers` (Boolean) Can SMS subscribers subscribe to this Status Page?
- `enable_webhook_subscribers` (Boolean) Can Webhook subscribers subscribe to this Status Page?
- `favicon_file_id` (String) Status Page Favicon File ID. The ID of a `oneuptime_file`.
- `footer_html` (String) Status Page Custom HTML Footer. Served only from a verified custom domain.
- `header_html` (String) Status Page Custom HTML Header. Served only from a verified custom domain.
- `hide_powered_by_one_uptime_branding` (Boolean) Hide Powered By OneUptime Branding?
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ip_whitelist` (String) IP Whitelist for this Status Page. One IP per line. Only used if the status page is private.
- `is_archived` (Boolean) Archived status pages are hidden from the Status Pages list, are not served to visitors, and send nothing to their subscribers. Unarchiving puts them back online.
- `is_owner_notified_of_resource_creation` (Boolean) Are owners notified of when this resource is created?
- `is_public_status_page` (Boolean) Is this status page public?
- `is_report_enabled` (Boolean) Whether this status page's email subscribers get reports. Turned on without a schedule, reports go out on the 1st of every month at 09:00 in the report timezone, each covering the calendar month before it.
- `logo_alt_text` (String) Alternative text for the logo image, read by screen readers for accessibility.
- `logo_file_id` (String) Status Page Logo File ID. The ID of a `oneuptime_file`.
- `name` (String) Any friendly name of this object.
- `only_show_scoped_incidents` (Boolean) When on, this status page shows and notifies its subscribers about only the incidents limited to it. Incidents that are not limited to any status page never reach it.
- `overall_uptime_percent_precision` (String) Overall Precision of uptime percent for this status page.
- `overview_page_description` (String) Overview Page description for your status page. This is a markdown field.
- `page_description` (String) Description of your Status Page. This is used for SEO.
- `page_title` (String) Title of your Status Page. This is used for SEO.
- `report_data_in_days` (Number) How many days of data should be included in the report?
- `report_period_type` (String) Should the report cover a rolling number of days, or the previous whole calendar period?
- `report_timezone` (String) The timezone report periods and send times are resolved in. A monthly report in this timezone runs from the 1st at 00:00 to the last day at 23:59.
- `require_sso_for_login` (Boolean) Should SSO be required to login to Private Status Page.
- `show_announcement_history_in_days` (Number) How many days of announcement history should be shown on the status page (in days)?
- `show_announcements_on_status_page` (Boolean) Show Announcements on Status Page?
- `show_episode_history_in_days` (Number) How many days of episode history to show on the status page.
- `show_episode_labels_on_status_page` (Boolean) Show Episode Labels on Status Page?
- `show_episodes_on_status_page` (Boolean) Show Incident Episodes on Status Page?
- `show_incident_history_in_days` (Number) How many days of incident history should be shown on the status page (in days)?
- `show_incident_labels_on_status_page` (Boolean) Show Incident Labels on Status Page?
- `show_incidents_on_status_page` (Boolean) Show Incidents on Status Page?
- `show_overall_uptime_percent_on_status_page` (Boolean) Show Overall Uptime Percent on Status Page?
- `show_scheduled_event_history_in_days` (Number) How many days of scheduled event history should be shown on the status page (in days)?
- `show_scheduled_event_labels_on_status_page` (Boolean) Show Scheduled Event Labels on Status Page?
- `show_scheduled_maintenance_events_on_status_page` (Boolean) Show Scheduled Maintenance Events on Status Page?
- `show_subscriber_page_on_status_page` (Boolean) Show Subscriber Page on Status Page?
- `show_uptime_history_in_days` (Number) How many days of uptime history should be shown on the status page? Maximum is 90 days.
- `slug` (String) Friendly globally unique name for your object.
- `smtp_config_id` (String) ID of your SMTP Config Resource which is used to send email to subscribers.
- `subscriber_email_notification_footer_text` (String) Text to send to subscribers in the footer of the email.

### Read-Only

- `archived_at` (String) When this status page was archived. Empty while it is not archived.
- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `default_bar_color` (String) Default color of the bar on the overview page.
- `downtime_monitor_statuses` (Set of String) List of monitors statuses that are considered as "down" for this status page. IDs of `oneuptime_monitor_status` resources.
- `enabled_languages` (String) Languages offered in the footer language switcher. Leave empty to offer all supported languages. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `master_password` (String, Sensitive) Password required to unlock a private status page. This value is stored as a secure hash.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `report_recurring_interval` (String) How often a report goes out. Left out when reports are turned on, it is every month. A JSON value: write it with `jsonencode()`.
- `report_start_date_time` (String) When the first report goes out. Every later one follows it by the recurring interval, at the same time of day. Left out when reports are turned on, it is 09:00 in the report timezone at the start of the next period of the interval: the next 1st of the month for a monthly schedule (the default), the next Monday for a weekly one, the next day for a daily one and the next 1 January for a yearly one. An hourly schedule starts at the next full hour.
- `send_next_report_by` (String) When the next report goes out. The server works it out from the schedule.
- `subscriber_timezones` (String) Timezones of subscribers to this status page. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
