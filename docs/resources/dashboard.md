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
  name        = "Example dashboard"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `dashboard_view_config` (String) Configuration of Dashboard View. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.
- `enable_master_password` (Boolean) Require visitors to enter a master password before viewing a public dashboard. Defaults to `false`.
- `favicon_file_id` (String) Dashboard Favicon File ID. The ID of a `oneuptime_file`.
- `ip_whitelist` (String) IP Whitelist for this dashboard. One IP per line. Only used when the dashboard is public.
- `is_archived` (Boolean) Archived dashboards are hidden from the Dashboards list and their public link stops working. Unarchiving restores them as they were. Defaults to `false`.
- `is_public_dashboard` (Boolean) Is this dashboard public? Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `logo_file_id` (String) Dashboard Logo File ID. The ID of a `oneuptime_file`.
- `master_password` (String, Sensitive) Password required to unlock a public dashboard. This value is stored as a secure hash.
- `page_description` (String) Description of the public dashboard page. This will be used for SEO.
- `page_title` (String) Title of the public dashboard page. This will be used for SEO and the browser tab.

### Read-Only

- `archived_at` (String) When this dashboard was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing dashboard by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_dashboard.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_dashboard.example <id>
```
