---
page_title: "oneuptime_status_page Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage status pages for your project.
---

# oneuptime_status_page (Resource)

Manage status pages for your project.

## Example Usage

```terraform
resource "oneuptime_status_page" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `page_title` (String) Title of your Status Page. This is used for SEO...
- `page_description` (String) Description of your Status Page. This is used for SEO...
- `enable_search_engine_indexing` (Bool) Should search engines like Google and Bing be allowed to index this status page? Turn this off to keep the page reachable by link but out of search results...
- `description` (String) Friendly description that will help you remember..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID..
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID..
- `cover_image_file_id` (String) A unique identifier for an object, represented as a UUID..
- `header_html` (String) Status Page Custom HTML Header. Served only from a verified custom domain...
- `footer_html` (String) Status Page Custom HTML Footer. Served only from a verified custom domain...
- `custom_css` (String) Status Page Custom CSS. Served only from a verified custom domain...
- `custom_java_script` (String) Status Page Custom JavaScript. This runs when the status page is loaded from a verified custom domain...
- `is_public_status_page` (Bool) Is this status page public?..
- `enable_mcp_server` (Bool) Can AI agents read this status page over the public OneUptime MCP server? This does not affect the status page website, its RSS feed, or its public JSON API...
- `enable_master_password` (Bool) Require visitors to enter a master password before viewing a private status page...
- `master_password` (String, Sensitive) Password required to unlock a private status page. This value is stored as a secure hash...
- `show_incident_labels_on_status_page` (Bool) Show Incident Labels on Status Page?..
- `show_scheduled_event_labels_on_status_page` (Bool) Show Scheduled Event Labels on Status Page?..
- `enable_email_subscribers` (Bool) Can email subscribers subscribe to this Status Page?..
- `allow_subscribers_to_choose_resources` (Bool) Can subscribers choose which resources to subscribe to?..
- `allow_subscribers_to_choose_event_types` (Bool) Can subscribers choose which event type like Announcements, Incidents, Scheduled Events to subscribe to?..
- `enable_sms_subscribers` (Bool) Can SMS subscribers subscribe to this Status Page?..
- `enable_slack_subscribers` (Bool) Can Slack subscribers subscribe to this Status Page?..
- `enable_microsoft_teams_subscribers` (Bool) Can Microsoft Teams subscribers subscribe to this Status Page?..
- `enable_webhook_subscribers` (Bool) Can Webhook subscribers subscribe to this Status Page?..
- `copyright_text` (String) Copyright Text..
- `logo_alt_text` (String) Alternative text for the logo image, read by screen readers for accessibility...
- `cover_image_alt_text` (String) Alternative text for the cover image, read by screen readers for accessibility. Leave blank if the cover image is purely decorative...
- `custom_fields` (String) Custom Fields on this resource...
- `require_sso_for_login` (Bool) Should SSO be required to login to Private Status Page..
- `smtp_config_id` (String) A unique identifier for an object, represented as a UUID..
- `call_sms_config_id` (String) A unique identifier for an object, represented as a UUID..
- `show_incident_history_in_days` (Number) How many days of incident history should be shown on the status page (in days)?..
- `show_announcement_history_in_days` (Number) How many days of announcement history should be shown on the status page (in days)?..
- `show_scheduled_event_history_in_days` (Number) How many days of scheduled event history should be shown on the status page (in days)?..
- `overview_page_description` (String) Overview Page description for your status page. This is a markdown field...
- `hide_powered_by_one_uptime_branding` (Bool) Hide Powered By OneUptime Branding?..
- `default_bar_color` (String) Color object.
- `subscriber_timezones` (String) Timezones of subscribers to this status page...
- `is_report_enabled` (Bool) Is Report Enabled for this Status Page?..
- `report_start_date_time` (String) A date time object..
- `report_recurring_interval` (String) How often would you like to send the report?..
- `send_next_report_by` (String) A date time object..
- `report_data_in_days` (Number) How many days of data should be included in the report?..
- `report_period_type` (String) Should the report cover a rolling number of days, or the previous whole calendar period?..
- `report_timezone` (String) The timezone report periods and send times are resolved in. A monthly report in this timezone runs from the 1st at 00:00 to the last day at 23:59...
- `show_overall_uptime_percent_on_status_page` (Bool) Show Overall Uptime Percent on Status Page?..
- `overall_uptime_percent_precision` (String) Overall Precision of uptime percent for this status page...
- `subscriber_email_notification_footer_text` (String) Text to send to subscribers in the footer of the email...
- `enable_custom_subscriber_email_notification_footer_text` (Bool) Enable custom footer text in subscriber email notifications...
- `show_incidents_on_status_page` (Bool) Show Incidents on Status Page?..
- `show_announcements_on_status_page` (Bool) Show Announcements on Status Page?..
- `show_episodes_on_status_page` (Bool) Show Incident Episodes on Status Page?..
- `show_episode_history_in_days` (Number) How many days of episode history to show on the status page..
- `show_episode_labels_on_status_page` (Bool) Show Episode Labels on Status Page?..
- `show_scheduled_maintenance_events_on_status_page` (Bool) Show Scheduled Maintenance Events on Status Page?..
- `show_subscriber_page_on_status_page` (Bool) Show Subscriber Page on Status Page?..
- `ip_whitelist` (String) IP Whitelist for this Status Page. One IP per line. Only used if the status page is private...
- `enable_embedded_overall_status` (Bool) Enable embedded overall status badge that can be displayed on external websites?..
- `show_uptime_history_in_days` (Number) How many days of uptime history should be shown on the status page? Maximum is 90 days...
- `embedded_overall_status_token` (String) Security token required to access the embedded overall status badge. This token must be provided in the URL...
- `default_language` (String) Default language that the status page is shown in when a visitor arrives for the first time...
- `enabled_languages` (String) Languages offered in the footer language switcher. Leave empty to offer all supported languages...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?..
- `downtime_monitor_statuses` (Set) List of monitors statuses that are considered as "down" for this status page...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page.example <id>
```
