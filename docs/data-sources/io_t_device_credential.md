---
page_title: "oneuptime_io_t_device_credential Data Source - oneuptime"
subcategory: "Other"
description: |-
  Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.
---

# oneuptime_io_t_device_credential (Data Source)

Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_io_t_device_credential" "by_name" {
  name = "example-io_t_device_credential"
}

data "oneuptime_io_t_device_credential" "by_id" {
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
- `iot_fleet_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `external_id` (String) The device id — must match the device.id label the device stamps on its datapoints. It is also the <device> segment of the device's MQTT topics, so a device that reports directly over MQTT cannot use an id containing '/', '+', or '#' (such devices can still report through a gateway)... Computed.
- `is_enabled` (Bool) Disabled credentials are rejected at MQTT CONNECT and stop the device's silent-death offline detection... Computed.
- `last_connected_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `secret_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
