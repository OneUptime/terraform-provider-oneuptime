---
page_title: "oneuptime_status_page Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage status pages for your project.
---

# oneuptime_status_page (Data Source)

Manage status pages for your project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page" "by_name" {
  name = "example-status_page"
}

data "oneuptime_status_page" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `page_title` (String) Title of your Status Page. This is used for SEO... Computed.
- `page_description` (String) Description of your Status Page. This is used for SEO... Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `cover_image_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `header_html` (String) Status Page Custom HTML Header.. Computed.
- `footer_html` (String) Status Page Custom HTML Footer.. Computed.
- `custom_css` (String) Status Page Custom CSS Header.. Computed.
- `custom_java_script` (String) Status Page Custom JavaScript. This runs when the status page is loaded... Computed.
- `is_public_status_page` (Bool) Is this status page public?.. Computed.
- `enable_mcp_server` (Bool) Can AI agents read this status page over the public OneUptime MCP server? This does not affect the status page website, its RSS feed, or its public JSON API... Computed.
- `enable_master_password` (Bool) Require visitors to enter a master password before viewing a private status page... Computed.
- `master_password` (String, Sensitive) Password required to unlock a private status page. This value is stored as a secure hash... Computed.
- `show_incident_labels_on_status_page` (Bool) Show Incident Labels on Status Page?.. Computed.
- `show_scheduled_event_labels_on_status_page` (Bool) Show Scheduled Event Labels on Status Page?.. Computed.
- `enable_email_subscribers` (Bool) Can email subscribers subscribe to this Status Page?.. Computed.
- `allow_subscribers_to_choose_resources` (Bool) Can subscribers choose which resources to subscribe to?.. Computed.
- `allow_subscribers_to_choose_event_types` (Bool) Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?.. Computed.
- `enable_sms_subscribers` (Bool) Can SMS subscribers subscribe to this Status Page?.. Computed.
- `enable_slack_subscribers` (Bool) Can Slack subscribers subscribe to this Status Page?.. Computed.
- `enable_microsoft_teams_subscribers` (Bool) Can Microsoft Teams subscribers subscribe to this Status Page?.. Computed.
- `enable_webhook_subscribers` (Bool) Can Webhook subscribers subscribe to this Status Page?.. Computed.
- `copyright_text` (String) Copyright Text.. Computed.
- `logo_alt_text` (String) Alternative text for the logo image, read by screen readers for accessibility... Computed.
- `cover_image_alt_text` (String) Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative... Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `require_sso_for_login` (Bool) Should SSO be required to login to Private Status Page.. Computed.
- `smtp_config_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `call_sms_config_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?.. Computed.
- `show_incident_history_in_days` (Number) How many days of incident history should be shown on the status page (in days)?.. Computed.
- `show_announcement_history_in_days` (Number) How many days of announcement history should be shown on the status page (in days)?.. Computed.
- `show_scheduled_event_history_in_days` (Number) How many days of scheduled event history should be shown on the status page (in days)?.. Computed.
- `overview_page_description` (String) Overview Page description for your status page. This is a markdown field... Computed.
- `hide_powered_by_one_uptime_branding` (Bool) Hide Powered By OneUptime Branding?.. Computed.
- `default_bar_color` (String) Color object. Computed.
- `downtime_monitor_statuses` (Set) List of monitors statuses that are considered as "down" for this status page... Computed.
- `subscriber_timezones` (String) Timezones of subscribers to this status page... Computed.
- `is_report_enabled` (Bool) Is Report Enabled for this Status Page?.. Computed.
- `report_start_date_time` (String) A date time object.. Computed.
- `report_recurring_interval` (String) How often would you like to send the report?.. Computed.
- `send_next_report_by` (String) A date time object.. Computed.
- `report_data_in_days` (Number) How many days of data should be included in the report?.. Computed.
- `report_period_type` (String) Should the report cover a rolling number of days, or the previous whole calendar period?.. Computed.
- `report_timezone` (String) The timezone report periods and send times are resolved in. A monthly report in this timezone runs from the 1st at 00:00 to the last day at 23:59... Computed.
- `show_overall_uptime_percent_on_status_page` (Bool) Show Overall Uptime Percent on Status Page?.. Computed.
- `overall_uptime_percent_precision` (String) Overall Precision of uptime percent for this status page... Computed.
- `subscriber_email_notification_footer_text` (String) Text to send to subscribers in the footer of the email... Computed.
- `enable_custom_subscriber_email_notification_footer_text` (Bool) Enable custom footer text in subscriber email notifications... Computed.
- `show_incidents_on_status_page` (Bool) Show Incidents on Status Page?.. Computed.
- `show_announcements_on_status_page` (Bool) Show Announcements on Status Page?.. Computed.
- `show_episodes_on_status_page` (Bool) Show Incident Episodes on Status Page?.. Computed.
- `show_episode_history_in_days` (Number) How many days of episode history to show on the status page.. Computed.
- `show_episode_labels_on_status_page` (Bool) Show Episode Labels on Status Page?.. Computed.
- `show_scheduled_maintenance_events_on_status_page` (Bool) Show Scheduled Maintenance Events on Status Page?.. Computed.
- `show_subscriber_page_on_status_page` (Bool) Show Subscriber Page on Status Page?.. Computed.
- `ip_whitelist` (String) IP Whitelist for this Status Page. One IP per line. Only used if the status page is private... Computed.
- `enable_embedded_overall_status` (Bool) Enable embedded overall status badge that can be displayed on external websites?.. Computed.
- `show_uptime_history_in_days` (Number) How many days of uptime history should be shown on the status page? Maximum is 90 days... Computed.
- `embedded_overall_status_token` (String) Security token required to access the embedded overall status badge. This token must be provided in the URL... Computed.
- `default_language` (String) Default language that the status page is shown in when a visitor arrives for the first time... Computed.
- `enabled_languages` (String) Languages offered in the footer language switcher. Leave empty to offer all supported languages... Computed.
