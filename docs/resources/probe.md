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
- `probe_version` (String) A version label for the probe. Leave it out and it is 1.0.0.
- `should_auto_enable_probe_on_new_monitors` (Boolean) Auto Enable Probe on New Monitors.

### Read-Only

- `checks_last_late_at` (String) When OneUptime last found one of this probe's monitors waiting more than a minute past its check time. Set by OneUptime.
- `checks_late_since` (String) When this probe's monitors started waiting more than a minute past their check time. Set by OneUptime; empty while they are checked on time.
- `connection_status` (String) Connection Status of the Probe.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `device_polls_last_late_at` (String) When OneUptime last found one of the network devices this probe polls waiting more than a minute past its poll time. Set by OneUptime.
- `device_polls_late_since` (String) When the network devices this probe polls started waiting more than a minute past their poll time. Set by OneUptime; empty while they are polled on time.
- `egress_ip_address` (String) The IP address this probe last checked in from, as OneUptime saw it. Set by OneUptime from the probe's own requests; the blocklist never blocks it.
- `id` (String) Unique identifier for the resource.
- `interface_ip_addresses` (String) The IP addresses on this probe's own network interfaces, as the probe last listed them when it checked in, separated by commas. Set by OneUptime from the probe's own requests; Security Events knows the probe's own polls and checks by them, and the blocklist never blocks them.
- `last_alive` (String) When the probe last checked in. Set by OneUptime.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `major_version_behind_notified_for` (Number) The major release of OneUptime this probe's owners were last told the probe is behind. Set by OneUptime once a day; empty while the probe runs OneUptime's major release.
- `packet_capture_capability` (String) What the probe last reported about packet capture: whether it is turned on (PROBE_PACKET_CAPTURE_ENABLED on the probe), whether tcpdump is installed, the network interfaces it can capture on, and the limits its operator set. Managed by the probe. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project the probe belongs to. A global probe belongs to none.
- `reported_version` (String) The OneUptime release this probe last said it runs. Set by OneUptime from the probe's own requests. Empty until the probe connects, and for a probe older than 17.1.0, which does not say.
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
