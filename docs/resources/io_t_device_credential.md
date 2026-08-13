---
page_title: "oneuptime_io_t_device_credential Resource - oneuptime"
subcategory: "Other"
description: |-
  Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.
---

# oneuptime_io_t_device_credential (Resource)

Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.

## Example Usage

```terraform
resource "oneuptime_io_t_device_credential" "example" {
  iot_fleet_id = "123e4567-e89b-12d3-a456-426614174000"
  external_id = "Example short text"
  name = "Example short text"
}
```

## Schema

### Required

- `iot_fleet_id` (String) A unique identifier for an object, represented as a UUID..
- `external_id` (String) The device id — must match the device.id label the device stamps on its datapoints. It is also the <device> segment of the device's MQTT topics, so a device that reports directly over MQTT cannot use an id containing '/', '+', or '#' (such devices can still report through a gateway)...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Any friendly name of this device..
- `is_enabled` (Bool) Disabled credentials are rejected at MQTT CONNECT and stop the device's silent-death offline detection...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_connected_at` (String) A date time object..
- `secret_key` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_io_t_device_credential.example <id>
```
