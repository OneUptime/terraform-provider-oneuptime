---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_status_page Resource - oneuptime"
subcategory: ""
description: |-
  Status page resource
---

# oneuptime_status_page (Resource)

Status page resource

## Example Usage

```terraform
resource "oneuptime_status_page" "example" {
  name = "example-resource"
  description = "Example resource"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `name` (String) Name. Required.
- `page_title` (String) Page Title. Optional.
- `page_description` (String) Page Description. Optional.
- `description` (String) Description. Optional.
- `labels` (List) Labels. Optional.
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `cover_image_file_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `header_h_t_m_l` (String) Header HTML. Optional.
- `footer_h_t_m_l` (String) Footer HTML. Optional.
- `custom_c_s_s` (String) CSS. Optional.
- `custom_java_script` (String) JavaScript. Optional.
- `is_public_status_page` (Bool) Public Status Page. Optional.
- `show_incident_labels_on_status_page` (Bool) Show Incident Labels. Optional.
- `show_scheduled_event_labels_on_status_page` (Bool) Show Scheduled Event Labels. Optional.
- `enable_subscribers` (Bool) Enable Subscribers. Optional.
- `enable_email_subscribers` (Bool) Enable Email Subscribers. Optional.
- `allow_subscribers_to_choose_resources` (Bool) Allow Subscribers to Choose Resources. Optional.
- `allow_subscribers_to_choose_event_types` (Bool) Allow Subscribers to subscribe to event types. Optional.
- `enable_sms_subscribers` (Bool) Enable SMS Subscribers. Optional.
- `enable_slack_subscribers` (Bool) Enable Slack Subscribers. Optional.
- `copyright_text` (String) Copyright Text. Optional.
- `custom_fields` (Map) Custom Fields. Optional.
- `require_sso_for_login` (Bool) Require SSO. Optional.
- `smtp_config_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `call_sms_config_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `is_owner_notified_of_resource_creation` (Bool) Are Owners Notified Of Resource Creation?. Optional.
- `show_incident_history_in_days` (Number) Show incident history in days. Optional.
- `show_announcement_history_in_days` (Number) Show announcement history in days. Optional.
- `show_scheduled_event_history_in_days` (Number) Show scheduled event history in days. Optional.
- `overview_page_description` (String) Overview Page Description. Optional.
- `hide_powered_by_one_uptime_branding` (Bool) Hide Powered By OneUptime Branding. Optional.
- `default_bar_color` (Map) Color object. Optional.
- `downtime_monitor_statuses` (List) Downtime Monitor Statuses. Optional.
- `subscriber_timezones` (Map) Subscriber Timezones. Optional.
- `is_report_enabled` (Bool) Is Report Enabled. Optional.
- `report_start_date_time` (Map) A date time object.. Optional.
- `report_recurring_interval` (Map) Report Recurring Interval. Optional.
- `send_next_report_by` (Map) A date time object.. Optional.
- `report_data_in_days` (Number) Report data for the last N days. Optional.
- `show_overall_uptime_percent_on_status_page` (Bool) Show Overall Uptime Percent on Status Page. Optional.
- `overall_uptime_percent_precision` (String) Overall Uptime Percent Precision. Optional.
- `subscriber_email_notification_footer_text` (String) Subscriber Email Notification Footer Text. Optional.
- `show_incidents_on_status_page` (Bool) Show Incidents on Status Page. Optional.
- `show_announcements_on_status_page` (Bool) Show Announcements on Status Page. Optional.
- `show_scheduled_maintenance_events_on_status_page` (Bool) Show Scheduled Maintenance Events on Status Page. Optional.
- `show_subscriber_page_on_status_page` (Bool) Show Subscriber Page on Status Page. Optional.
- `ip_whitelist` (String) IP Whitelist. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `slug` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page.example <id>
```
