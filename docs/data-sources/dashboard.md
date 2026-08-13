---
page_title: "oneuptime_dashboard Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Create and manage Dashboards to visualize your data in a single place
---

# oneuptime_dashboard (Data Source)

Create and manage Dashboards to visualize your data in a single place Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_dashboard" "by_name" {
  name = "example-dashboard"
}

data "oneuptime_dashboard" "by_id" {
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
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `dashboard_view_config` (String) Configuration of Dashboard View.. Computed.
- `page_title` (String) Title of the public dashboard page. This will be used for SEO and the browser tab... Computed.
- `page_description` (String) Description of the public dashboard page. This will be used for SEO... Computed.
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_public_dashboard` (Bool) Is this dashboard public?.. Computed.
- `enable_master_password` (Bool) Require visitors to enter a master password before viewing a public dashboard... Computed.
- `master_password` (String, Sensitive) Password required to unlock a public dashboard. This value is stored as a secure hash... Computed.
- `ip_whitelist` (String) IP Whitelist for this dashboard. One IP per line. Only used when the dashboard is public... Computed.
