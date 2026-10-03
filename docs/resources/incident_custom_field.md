---
page_title: "oneuptime_incident_custom_field Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage custom fields for your incident.
---

# oneuptime_incident_custom_field (Resource)

Manage custom fields for your incident.

## Example Usage

```terraform
resource "oneuptime_incident_custom_field" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description of this custom field that will help you remember..
- `custom_field_type` (String) Is this field Text, Number or Boolean?..
- `dropdown_options` (String) Options and optional colors for dropdown fields. Plain one-per-line values remain supported...
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand...
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from...
- `is_required_on_create` (Bool) When on, an incident declared from the dashboard cannot be created until this field is filled in (a Boolean field must be ticked). It applies only to fields shown on create. Incidents created by monitors, the API, Slack, Microsoft Teams or AI can leave it empty, and it stays optional when an incident is edited later...
- `sort_order` (Number) Where this field appears among the incident's custom fields, lowest number first. A new field is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them...
- `show_on_create` (Bool) When on, this field is asked for in a Details step when an incident is declared from the dashboard, and incident templates can fill it in...
- `include_in_subscriber_notifications` (Bool) When on, this field and its value appear in the messages status page subscribers get about an incident: the default email, Slack and Microsoft Teams messages, and webhooks (under customFields, by the field's template variable key). The default SMS is kept short and leaves it out. Subscribers are usually people outside your team, so turn this on only for fields that are safe to share with them...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `variable_key` (String) The key this field is reached by in templates, as {{incident.customFields.<key>}}. Made from the field's name when it is created - lowercase letters, digits and underscores, with _2, _3 and so on added when another field already has it - and never changed afterwards, so renaming the field does not break templates that use it...
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_custom_field.example <id>
```
