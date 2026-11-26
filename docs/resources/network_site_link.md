---
page_title: "oneuptime_network_site_link Resource - oneuptime"
subcategory: "Network"
description: |-
  Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.
---

# oneuptime_network_site_link (Resource)

Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.

## Example Usage

```terraform
resource "oneuptime_network_site_link" "example" {
  from_site_id = oneuptime_network_site.example.id
  to_site_id   = oneuptime_network_site.example.id
  name         = "Example network site link"
}
```

## Schema

### Required

- `from_site_id` (String) ID of the Network Site this link starts from. The ID of a `oneuptime_network_site`.
- `to_site_id` (String) ID of the Network Site this link ends at. The ID of a `oneuptime_network_site`.

### Optional

- `monitor_id` (String) ID of the Monitor whose status colors this link on map views. The ID of a `oneuptime_monitor`.
- `name` (String) Friendly name for this link.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network site link by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_site_link.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_site_link.example <id>
```
