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
  name                 = "Example network device link rule"
  child_device_labels  = [oneuptime_label.example.id]
  parent_device_labels = [oneuptime_label.example.id]
  description          = "Managed by Terraform"
}
```

## Schema

### Required

- `child_device_labels` (Set of String) Devices carrying ALL of these labels each get one uplink drawn to the parent device. Empty matches nothing — a rule that linked every device in the project is never what anyone meant. IDs of `oneuptime_label` resources.
- `name` (String) Friendly name for this rule.
- `parent_device_labels` (Set of String) The device carrying ALL of these labels is what the children uplink to. It has to identify exactly one device: match none and the rule draws nothing, match several and the rule is ambiguous and also draws nothing. IDs of `oneuptime_label` resources.

### Optional

- `description` (String) Description of this rule.
- `is_enabled` (Boolean) Whether this rule draws links. Disable to take its edges off the map without deleting the rule. Defaults to `true`.
- `scope` (String) How wide the 'exactly one parent device' question is asked. Project (the default) looks for one parent across the whole project. Site asks once per site, so the same rule can draw an uplink in every building. Rules created before this existed are Project.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device link rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_link_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_link_rule.example <id>
```
