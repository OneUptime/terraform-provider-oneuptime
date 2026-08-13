---
page_title: "oneuptime_dashboard Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Create and manage Dashboards to visualize your data in a single place
---

# oneuptime_dashboard (Resource)

Create and manage Dashboards to visualize your data in a single place

## Example Usage

```terraform
resource "oneuptime_dashboard" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `dashboard_view_config` (String) Configuration of Dashboard View..
- `page_title` (String) Title of the public dashboard page. This will be used for SEO and the browser tab...
- `page_description` (String) Description of the public dashboard page. This will be used for SEO...
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID..
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID..
- `is_public_dashboard` (Bool) Is this dashboard public?..
- `enable_master_password` (Bool) Require visitors to enter a master password before viewing a public dashboard...
- `master_password` (String, Sensitive) Password required to unlock a public dashboard. This value is stored as a secure hash...
- `ip_whitelist` (String) IP Whitelist for this dashboard. One IP per line. Only used when the dashboard is public...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_dashboard.example <id>
```
