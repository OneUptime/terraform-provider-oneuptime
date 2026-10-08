---
page_title: "oneuptime_iot_device_credential Data Source - oneuptime"
subcategory: "Other"
description: |-
  Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.
---

# oneuptime_iot_device_credential (Data Source)

Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one iot device credential may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_io_t_device_credential` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_iot_device_credential" "example" {
  name = "Example iot device credential"
}

# Or by id:
data "oneuptime_iot_device_credential" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `external_id` (String) The device id — must match the device.id label the device stamps on its datapoints. It is also the <device> segment of the device's MQTT topics, so a device that reports directly over MQTT cannot use an id containing '/', '+', or '#' (such devices can still report through a gateway).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `iot_fleet_id` (String) ID of the IoT Fleet this device belongs to. The ID of a `oneuptime_iot_fleet`.
- `is_enabled` (Boolean) Disabled credentials are rejected at MQTT CONNECT and stop the device's silent-death offline detection.
- `name` (String) Any friendly name of this device.
- `secret_key` (String) Secret this device presents as the MQTT password (with this credential's ID as the username).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_connected_at` (String) When this credential last authenticated an MQTT connection.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
