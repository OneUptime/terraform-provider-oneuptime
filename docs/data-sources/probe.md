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
- `description` (String) Friendly description that will help you remember.
- `egress_ip_address` (String) The IP address this probe last checked in from, as OneUptime saw it. Set by OneUptime from the probe's own requests; the blocklist never blocks it.
- `icon_file_id` (String) Probe Page Icon File ID. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `reported_version` (String) The OneUptime release this probe last said it runs. Set by OneUptime from the probe's own requests. Empty until the probe connects, and for a probe older than 17.0.1, which does not say.
- `should_auto_enable_probe_on_new_monitors` (Boolean) Auto Enable Probe on New Monitors.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `checks_last_late_at` (String) When OneUptime last found one of this probe's monitors waiting more than a minute past its check time. Set by OneUptime.
- `checks_late_since` (String) When this probe's monitors started waiting more than a minute past their check time. Set by OneUptime; empty while they are checked on time.
- `created_at` (String) Date and Time when the object was created.
- `device_polls_last_late_at` (String) When OneUptime last found one of the network devices this probe polls waiting more than a minute past its poll time. Set by OneUptime.
- `device_polls_late_since` (String) When the network devices this probe polls started waiting more than a minute past their poll time. Set by OneUptime; empty while they are polled on time.
- `key` (String, Sensitive) The secret the probe connects with, its PROBE_KEY. Leave it out and OneUptime generates one.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_alive` (String) When the probe last checked in. Set by OneUptime.
- `packet_capture_capability` (String) What the probe last reported about packet capture: whether it is turned on (PROBE_PACKET_CAPTURE_ENABLED on the probe), whether tcpdump is installed, the network interfaces it can capture on, and the limits its operator set. Managed by the probe. A JSON value: read it with `jsondecode()`.
- `probe_version` (String) A version label for the probe. Leave it out and it is 1.0.0.
- `project_id` (String) ID of the project the probe belongs to. A global probe belongs to none.
- `updated_at` (String) Date and Time when the object was updated.
