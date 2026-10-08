---
page_title: "oneuptime_host_label_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to hosts when matching hosts are created
---

# oneuptime_host_label_rule (Data Source)

Configure rules for automatically attaching labels to hosts when matching hosts are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one host label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_host_label_rule" "example" {
  name = "Example host label rule"
}

# Or by id:
data "oneuptime_host_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this host label rule.
- `host_description_pattern` (String) Regex (case-insensitive) matched against the host description. Leave empty to match any description.
- `host_name_pattern` (String) Regex (case-insensitive) matched against the host name. Leave empty to match any name.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this host label rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `host_labels` (Set of String) Only trigger for hosts that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `labels_to_add` (Set of String) Labels to attach to the host when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
