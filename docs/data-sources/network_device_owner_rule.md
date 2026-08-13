---
page_title: "oneuptime_network_device_owner_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching network devices are created
---

# oneuptime_network_device_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching network devices are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device_owner_rule" "by_name" {
  name = "example-network_device_owner_rule"
}

data "oneuptime_network_device_owner_rule" "by_id" {
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
- `description` (String) Description of this network device owner rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule.. Computed.
- `network_device_labels` (Set) Only trigger for network devices that have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `network_device_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the network device name. Leave empty to match any name... Computed.
- `network_device_description_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the network device description. Leave empty to match any description... Computed.
- `owner_users` (Set) Users to add as owners on the network device when this rule matches... Computed.
- `owner_teams` (Set) Teams to add as owners on the network device when this rule matches... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
