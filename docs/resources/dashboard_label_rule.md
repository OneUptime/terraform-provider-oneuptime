---
page_title: "oneuptime_dashboard_label_rule Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Configure rules for automatically attaching labels to dashboards when matching dashboards are created
---

# oneuptime_dashboard_label_rule (Resource)

Configure rules for automatically attaching labels to dashboards when matching dashboards are created

## Example Usage

```terraform
resource "oneuptime_dashboard_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this dashboard label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this dashboard label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `dashboard_labels` (Set) Only trigger for dashboards that already have at least one of these labels. Leave empty to match regardless of labels...
- `dashboard_name_pattern` (String) Regex (case-insensitive) matched against the dashboard name. Leave empty to match any name...
- `dashboard_description_pattern` (String) Regex (case-insensitive) matched against the dashboard description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the dashboard when this rule matches. Already-attached labels are not duplicated...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_dashboard_label_rule.example <id>
```
