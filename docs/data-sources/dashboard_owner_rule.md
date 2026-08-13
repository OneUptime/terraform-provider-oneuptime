---
page_title: "oneuptime_dashboard_owner_rule Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Configure rules for automatically assigning owner users and teams when matching dashboards are created
---

# oneuptime_dashboard_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching dashboards are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_dashboard_owner_rule" "by_name" {
  name = "example-dashboard_owner_rule"
}

data "oneuptime_dashboard_owner_rule" "by_id" {
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
- `description` (String) Description of this dashboard owner rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule.. Computed.
- `dashboard_labels` (Set) Only trigger for dashboards that have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `dashboard_name_pattern` (String) Regex (case-insensitive) matched against the dashboard name. Leave empty to match any name... Computed.
- `dashboard_description_pattern` (String) Regex (case-insensitive) matched against the dashboard description. Leave empty to match any description... Computed.
- `owner_users` (Set) Users to add as owners on the dashboard when this rule matches... Computed.
- `owner_teams` (Set) Teams to add as owners on the dashboard when this rule matches... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
