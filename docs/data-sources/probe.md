---
page_title: "oneuptime_probe Data Source - oneuptime"
subcategory: "Probes"
description: |-
  Manages custom probes. Deploy probes anywhere in the world and connect it to your project.
---

# oneuptime_probe (Data Source)

Manages custom probes. Deploy probes anywhere in the world and connect it to your project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one probe may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_probe" "example" {
  name = "Example probe"
}

# Or by id:
data "oneuptime_probe" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `connection_status` (String) Connection Status of the Probe.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `icon_file_id` (String) Probe Page Icon File ID. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `key` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]
- `name` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]
- `should_auto_enable_probe_on_new_monitors` (Boolean) Auto Enable Probe on New Monitors.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_alive` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Settings Admin, Settings Member, Settings Viewer, Read Probe], Update: [No access - you don't have permission for this operation]
- `packet_capture_capability` (String) What the probe last reported about packet capture: whether it is turned on (PROBE_PACKET_CAPTURE_ENABLED on the probe), whether tcpdump is installed, the network interfaces it can capture on, and the limits its operator set. Managed by the probe. A JSON value: write it with `jsonencode()`.
- `probe_version` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]
- `project_id` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Public], Update: [No access - you don't have permission for this operation]
- `updated_at` (String) Date and Time when the object was updated.
