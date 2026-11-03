---
page_title: "oneuptime_data_security_issue_feed_item Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  The timeline of a data security issue: when it opened, changes of severity, status and assignee, and the cases, alerts and notifications it caused.
---

# oneuptime_data_security_issue_feed_item (Data Source)

The timeline of a data security issue: when it opened, changes of severity, status and assignee, and the cases, alerts and notifications it caused.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security issue feed item may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_issue_feed_item" "example" {
  data_security_issue_id = data.oneuptime_data_security_issue.example.id
}

# Or by id:
data "oneuptime_data_security_issue_feed_item" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person whose action wrote the entry. The ID of a `oneuptime_user` (see the data source).
- `data_security_issue_feed_event_type` (String) What kind of thing happened.
- `data_security_issue_id` (String) ID of the case this feed item is about. The ID of a `oneuptime_data_security_issue` (see the data source).
- `feed_info_in_markdown` (String) What happened, in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) Detail behind the entry, shown when it is expanded.
- `user_id` (String) ID of the person who did it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) The color the entry is drawn in.
- `posted_at` (String) When it happened.
- `project_id` (String) ID of the project this feed item belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
