---
page_title: "oneuptime_status_page_announcement_template Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcement templates for your status page
---

# oneuptime_status_page_announcement_template (Data Source)

Manage announcement templates for your status page Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_announcement_template" "by_name" {
  name = "example-status_page_announcement_template"
}

data "oneuptime_status_page_announcement_template" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `template_name` (String) Name of the announcement template.. Computed.
- `template_description` (String) Description of the announcement template.. Computed.
- `title` (String) Title of the announcement.. Computed.
- `description` (String) Text of the announcement. This is in Markdown... Computed.
- `status_pages` (Set) Status Pages to show this announcement on... Computed.
- `monitors` (Set) List of monitors affected by this announcement template. If none are selected, all subscribers will be notified... Computed.
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about announcements created from this template?.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
