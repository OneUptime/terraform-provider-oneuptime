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
  template_name = "Example short text"
  title = "Example short text"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `template_name` (String) Name of the announcement template..
- `title` (String) Title of the announcement..
- `description` (String) Text of the announcement. This is in Markdown...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `template_description` (String) Description of the announcement template..
- `status_pages` (Set) Status Pages to show this announcement on...
- `monitors` (Set) List of monitors affected by this announcement template. If none are selected, all subscribers will be notified...
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about announcements created from this template?..
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
terraform import oneuptime_status_page_announcement_template.example <id>
```
