---
page_title: "oneuptime_status_page_label_rule Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure rules for automatically attaching labels to status pages when matching status pages are created
---

# oneuptime_status_page_label_rule (Data Source)

Configure rules for automatically attaching labels to status pages when matching status pages are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_label_rule" "example" {
  name = "Example status page label rule"
}

# Or by id:
data "oneuptime_status_page_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this status page label rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this status page label rule.
- `status_page_description_pattern` (String) Regex (case-insensitive) matched against the status page description. Leave empty to match any description.
- `status_page_name_pattern` (String) Regex (case-insensitive) matched against the status page name. Leave empty to match any name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `labels_to_add` (Set of String) Labels to attach to the status page when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_page_labels` (Set of String) Only trigger for status pages that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `updated_at` (String) Date and Time when the object was updated.
