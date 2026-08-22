---
page_title: "oneuptime_network_device_link_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.
---

# oneuptime_network_device_link_rule (Resource)

Draw uplinks on the network topology map from labels: every device carrying the child labels is linked to the single device carrying the parent labels.

## Example Usage

```terraform
resource "oneuptime_network_device_link_rule" "example" {
  name = "Example short text"
  child_device_labels = [{
    id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }]
  parent_device_labels = [{
    id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }]
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this rule..
- `child_device_labels` (Set) Devices carrying ALL of these labels each get one uplink drawn to the parent device. Empty matches nothing — a rule that linked every device in the project is never what anyone meant...
- `parent_device_labels` (Set) The device carrying ALL of these labels is what the children uplink to. It has to identify exactly one device: match none and the rule draws nothing, match several and the rule is ambiguous and also draws nothing...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this rule..
- `is_enabled` (Bool) Whether this rule draws links. Disable to take its edges off the map without deleting the rule...
- `scope` (String) How wide the 'exactly one parent device' question is asked. Project (the default) looks for one parent across the whole project. Site asks once per site, so the same rule can draw an uplink in every building. Rules created before this existed are Project...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_link_rule.example <id>
```
