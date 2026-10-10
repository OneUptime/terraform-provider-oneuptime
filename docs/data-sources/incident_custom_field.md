---
page_title: "oneuptime_incident_custom_field Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage custom fields for your incident.
---

# oneuptime_incident_custom_field (Data Source)

Manage custom fields for your incident.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident custom field may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_custom_field" "example" {
  name = "Example incident custom field"
}

# Or by id:
data "oneuptime_incident_custom_field" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description of this custom field that will help you remember.
- `dropdown_options` (String) Options and optional colors for dropdown fields, in the order they are listed. Plain one-per-line values remain supported. Records store an option as its text: changing an option here keeps the values records already hold. To rename an option and move those values with it, send the renames with the update, in miscDataProps: {"renamedDropdownOptions": [{"from": "Old text", "to": "New text"}]}. "to" must be one of the options.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `include_in_subscriber_notifications` (Boolean) When on, this field and its value appear in the messages status page subscribers get about an incident: the default email, Slack and Microsoft Teams messages, and webhooks (under customFields, by the field's template variable key). The default SMS is kept short and leaves it out. Subscribers are usually people outside your team, so turn this on only for fields that are safe to share with them.
- `is_required_on_create` (Boolean) When on, an incident declared from the dashboard cannot be created until this field is filled in (a Boolean field must be ticked). It applies only to fields shown on create. Incidents created by monitors, the API, Slack, Microsoft Teams or AI can leave it empty, and it stays optional when an incident is edited later.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand.
- `name` (String) Any friendly name of this object.
- `show_on_create` (Boolean) When on, this field is asked for in a Details step when an incident is declared from the dashboard, and incident templates can fill it in.
- `sort_order` (Number) Where this field appears among the incident's custom fields, lowest number first. A new field is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `variable_key` (String) The key this field is reached by in templates, as {{incident.customFields.<key>}}. Made from the field's name when it is created - lowercase letters, digits and underscores, with _2, _3 and so on added when another field already has it - and never changed afterwards, so renaming the field does not break templates that use it.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_field_type` (String) Is this field Text, Number or Boolean? A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
