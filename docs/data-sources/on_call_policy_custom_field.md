---
page_title: "oneuptime_on_call_policy_custom_field Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage custom fields for your on-call policy
---

# oneuptime_on_call_policy_custom_field (Data Source)

Manage custom fields for your on-call policy Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_policy_custom_field" "by_name" {
  name = "example-on_call_policy_custom_field"
}

data "oneuptime_on_call_policy_custom_field" "by_id" {
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
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
