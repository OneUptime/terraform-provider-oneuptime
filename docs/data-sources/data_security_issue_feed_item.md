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
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_security_issue_feed_item" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person whose action wrote the entry. The ID of a `oneuptime_user` (see the data source).
- `data_security_ai_proposal_id` (String) The OneUptime AI proposal this entry is about, when it is about one: a change made without asking offers its Undo here.
- `data_security_issue_feed_event_type` (String) What kind of thing happened.
- `data_security_issue_id` (String) ID of the case this feed item is about. The ID of a `oneuptime_data_security_issue` (see the data source).
- `feed_info_in_markdown` (String) What happened, in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `more_information_in_markdown` (String) Detail behind the entry, shown when it is expanded.
- `user_id` (String) ID of the person who did it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) The color the entry is drawn in.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `posted_at` (String) When it happened.
- `project_id` (String) ID of the project this feed item belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
