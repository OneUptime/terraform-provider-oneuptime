---
page_title: "oneuptime_dashboard Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Create and manage Dashboards to visualize your data in a single place
---

# oneuptime_dashboard (Data Source)

Create and manage Dashboards to visualize your data in a single place

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard" "example" {
  name = "Example dashboard"
}

# Or by id:
data "oneuptime_dashboard" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `enable_master_password` (Boolean) Require visitors to enter a master password before viewing a public dashboard.
- `favicon_file_id` (String) Dashboard Favicon File ID. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ip_whitelist` (String) IP Whitelist for this dashboard. One IP per line. Only used when the dashboard is public.
- `is_archived` (Boolean) Archived dashboards are hidden from the Dashboards list and their public link stops working. Unarchiving restores them as they were.
- `is_public_dashboard` (Boolean) Is this dashboard public?
- `logo_file_id` (String) Dashboard Logo File ID. The ID of a `oneuptime_file`.
- `name` (String) Any friendly name of this object.
- `page_description` (String) Description of the public dashboard page. This will be used for SEO.
- `page_title` (String) Title of the public dashboard page. This will be used for SEO and the browser tab.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When this dashboard was archived. Empty while it is not archived.
- `created_at` (String) Date and Time when the object was created.
- `dashboard_view_config` (String) Configuration of Dashboard View. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `master_password` (String, Sensitive) Password required to unlock a public dashboard. This value is stored as a secure hash.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
