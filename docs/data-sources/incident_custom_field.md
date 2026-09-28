---
page_title: "oneuptime_incident_custom_field Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage custom fields for your incident.
---

# oneuptime_incident_custom_field (Data Source)

Manage custom fields for your incident. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_custom_field" "by_name" {
  name = "example-incident_custom_field"
}

data "oneuptime_incident_custom_field" "by_id" {
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
- `description` (String) Friendly description of this custom field that will help you remember.. Computed.
- `custom_field_type` (String) Is this field Text, Number or Boolean?.. Computed.
- `dropdown_options` (String) Options and optional colors for dropdown fields. Plain one-per-line values remain supported... Computed.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand... Computed.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from... Computed.
- `is_required_on_create` (Bool) When on, an incident declared from the dashboard cannot be created until this field is filled in (a Boolean field must be ticked). It applies only to fields shown on create. Incidents created by monitors, the API, Slack, Microsoft Teams or AI can leave it empty, and it stays optional when an incident is edited later... Computed.
- `sort_order` (Number) Where this field appears among the incident's custom fields, lowest first. Fields with no order come after the ones that have one... Computed.
- `show_on_create` (Bool) When on, this field is asked for in a Details step when an incident is declared from the dashboard, and incident templates can fill it in... Computed.
- `include_in_subscriber_notifications` (Bool) When on, this field and its value appear in the messages status page subscribers get about an incident: the default email, Slack and Microsoft Teams messages, and webhooks (under customFields, by the field's template variable key). The default SMS is kept short and leaves it out. Subscribers are usually people outside your team, so turn this on only for fields that are safe to share with them... Computed.
- `variable_key` (String) The key this field is reached by in templates, as {{customFields.<key>}}. Made from the field's name when it is created - lowercase letters, digits and underscores, with _2, _3 and so on added when another field already has it - and never changed afterwards, so renaming the field does not break templates that use it... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
