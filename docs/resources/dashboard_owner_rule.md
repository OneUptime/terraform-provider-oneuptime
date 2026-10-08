---
page_title: "oneuptime_dashboard_owner_rule Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Configure rules for automatically assigning owner users and teams when matching dashboards are created
---

# oneuptime_dashboard_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching dashboards are created

## Example Usage

```terraform
resource "oneuptime_dashboard_owner_rule" "example" {
  name        = "Example dashboard owner rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this dashboard owner rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `dashboard_description_pattern` (String) Regex (case-insensitive) matched against the dashboard description. Leave empty to match any description.
- `dashboard_labels` (Set of String) Only trigger for dashboards that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `dashboard_name_pattern` (String) Regex (case-insensitive) matched against the dashboard name. Leave empty to match any name.
- `description` (String) Description of this dashboard owner rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `owner_teams` (Set of String) Teams to add as owners on the dashboard when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the dashboard when this rule matches. IDs of `oneuptime_user` records.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing dashboard owner rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_dashboard_owner_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_dashboard_owner_rule.example <id>
```
