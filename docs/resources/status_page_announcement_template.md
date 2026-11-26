---
page_title: "oneuptime_status_page_announcement_template Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcement templates for your status page
---

# oneuptime_status_page_announcement_template (Resource)

Manage announcement templates for your status page

## Example Usage

```terraform
resource "oneuptime_status_page_announcement_template" "example" {
  template_name = "Planned Upgrade Announcement"
  title         = "Scheduled System Upgrade - January 15, 2025"
  description   = "Managed by Terraform"
}
```

## Schema

### Required

- `description` (String) Text of the announcement. This is in Markdown.
- `template_name` (String) Name of the announcement template.
- `title` (String) Title of the announcement.

### Optional

- `monitors` (Set of String) List of monitors affected by this announcement template. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about announcements created from this template? Defaults to `true`.
- `status_pages` (Set of String) Status Pages to show this announcement on. IDs of `oneuptime_status_page` resources.
- `template_description` (String) Description of the announcement template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page announcement template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_announcement_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_announcement_template.example <id>
```
