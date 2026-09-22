---
page_title: "oneuptime_slo_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to SLOs when matching SLOs are created
---

# oneuptime_slo_label_rule (Resource)

Configure rules for automatically attaching labels to SLOs when matching SLOs are created

## Example Usage

```terraform
resource "oneuptime_slo_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this SLO label rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this SLO label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `service_level_objective_labels` (Set) Only trigger for SLOs that already have at least one of these labels. Leave empty to match regardless of labels...
- `service_level_objective_name_pattern` (String) Regex (case-insensitive) matched against the SLO name. Leave empty to match any name...
- `service_level_objective_description_pattern` (String) Regex (case-insensitive) matched against the SLO description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the SLO when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_slo_label_rule.example <id>
```
