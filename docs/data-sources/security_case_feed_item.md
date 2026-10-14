---
page_title: "oneuptime_security_case_feed_item Data Source - oneuptime"
subcategory: "Other"
description: |-
  The timeline of a security case: detections joining, changes of status, priority and assignee, notes, and every step a playbook took.
---

# oneuptime_security_case_feed_item (Data Source)

The timeline of a security case: detections joining, changes of status, priority and assignee, notes, and every step a playbook took.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case feed item may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_feed_item" "example" {
  security_case_id = oneuptime_security_case.example.id
}

# Or by id:
data "oneuptime_security_case_feed_item" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person whose action wrote the entry. The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) What happened, in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) Detail behind the entry, shown when it is expanded.
- `security_case_feed_event_type` (String) What kind of thing happened.
- `security_case_id` (String) ID of the case this feed item is about. The ID of a `oneuptime_security_case`.
- `user_id` (String) ID of the person who did it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) The color the entry is drawn in.
- `posted_at` (String) When it happened.
- `project_id` (String) ID of the project this feed item belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
