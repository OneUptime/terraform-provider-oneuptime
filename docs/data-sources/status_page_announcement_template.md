---
page_title: "oneuptime_status_page_announcement_template Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcement templates for your status page
---

# oneuptime_status_page_announcement_template (Data Source)

Manage announcement templates for your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page announcement template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_announcement_template" "example" {
  template_name = "example-template-name"
}

# Or by id:
data "oneuptime_status_page_announcement_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Text of the announcement. This is in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about announcements created from this template?
- `template_description` (String) Description of the announcement template.
- `template_name` (String) Name of the announcement template.
- `title` (String) Title of the announcement.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `monitors` (Set of String) List of monitors affected by this announcement template. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_pages` (Set of String) Status Pages to show this announcement on. IDs of `oneuptime_status_page` resources.
- `updated_at` (String) Date and Time when the object was updated.
