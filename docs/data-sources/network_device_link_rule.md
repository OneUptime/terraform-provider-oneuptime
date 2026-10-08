---
page_title: "oneuptime_network_device_link_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.
---

# oneuptime_network_device_link_rule (Data Source)

Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device link rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device_link_rule" "example" {
  name = "Example network device link rule"
}

# Or by id:
data "oneuptime_network_device_link_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule draws links. Disable to take its edges off the map without deleting the rule.
- `name` (String) Friendly name for this rule.
- `scope` (String) How wide the 'exactly one parent device' question is asked. Project (the default) looks for one parent across the whole project. Site asks once per site, so the same rule can draw an uplink in every building. Rules created before this existed are Project.

### Read-Only

- `child_device_labels` (Set of String) Devices carrying ALL of these labels each get one uplink drawn to the parent device. Empty matches nothing — a rule that linked every device in the project is never what anyone meant. IDs of `oneuptime_label` resources.
- `created_at` (String) Date and Time when the object was created.
- `parent_device_labels` (Set of String) The device carrying ALL of these labels is what the children uplink to. It has to identify exactly one device: match none and the rule draws nothing, match several and the rule is ambiguous and also draws nothing. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
