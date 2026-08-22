---
page_title: "oneuptime_network_device_link_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.
---

# oneuptime_network_device_link_rule (Data Source)

Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device_link_rule" "by_name" {
  name = "example-network_device_link_rule"
}

data "oneuptime_network_device_link_rule" "by_id" {
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
- `description` (String) Description of this rule.. Computed.
- `is_enabled` (Bool) Whether this rule draws links. Disable to take its edges off the map without deleting the rule... Computed.
- `child_device_labels` (Set) Devices carrying ALL of these labels each get one uplink drawn to the parent device. Empty matches nothing — a rule that linked every device in the project is never what anyone meant... Computed.
- `parent_device_labels` (Set) The device carrying ALL of these labels is what the children uplink to. It has to identify exactly one device: match none and the rule draws nothing, match several and the rule is ambiguous and also draws nothing... Computed.
- `scope` (String) How wide the 'exactly one parent device' question is asked. Project (the default) looks for one parent across the whole project. Site asks once per site, so the same rule can draw an uplink in every building. Rules created before this existed are Project... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
