---
page_title: "oneuptime_probe Resource - oneuptime"
subcategory: "Probes"
description: |-
  Manages custom probes. Deploy probes anywhere in the world and connect it to your project.
---

# oneuptime_probe (Resource)

Manages custom probes. Deploy probes anywhere in the world and connect it to your project.

## Example Usage

```terraform
resource "oneuptime_probe" "example" {
  name        = "Example probe"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `icon_file_id` (String) Probe Page Icon File ID. The ID of a `oneuptime_file`.
- `key` (String, Sensitive) The secret the probe connects with, its PROBE_KEY. Leave it out and OneUptime generates one.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_alive` (String) When the probe last checked in. Set by OneUptime.
- `probe_version` (String) A version label for the probe. Leave it out and it is 1.0.0.
- `should_auto_enable_probe_on_new_monitors` (Boolean) Auto Enable Probe on New Monitors.

### Read-Only

- `connection_status` (String) Connection Status of the Probe.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `id` (String) Unique identifier for the resource.
- `packet_capture_capability` (String) What the probe last reported about packet capture: whether it is turned on (PROBE_PACKET_CAPTURE_ENABLED on the probe), whether tcpdump is installed, the network interfaces it can capture on, and the limits its operator set. Managed by the probe. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project the probe belongs to. A global probe belongs to none.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing probe by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_probe.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_probe.example <id>
```
