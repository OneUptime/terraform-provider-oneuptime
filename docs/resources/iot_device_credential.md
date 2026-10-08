---
page_title: "oneuptime_iot_device_credential Resource - oneuptime"
subcategory: "Other"
description: |-
  Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.
---

# oneuptime_iot_device_credential (Resource)

Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection.

~> **Renamed:** this resource was called `oneuptime_io_t_device_credential` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing iot device credential:

```terraform
moved {
  from = oneuptime_io_t_device_credential.example
  to   = oneuptime_iot_device_credential.example
}
```

## Example Usage

```terraform
resource "oneuptime_iot_device_credential" "example" {
  iot_fleet_id = oneuptime_iot_fleet.example.id
  external_id  = "Example short text"
  name         = "Example iot device credential"
}
```

## Schema

### Required

- `external_id` (String) The device id — must match the device.id label the device stamps on its datapoints. It is also the <device> segment of the device's MQTT topics, so a device that reports directly over MQTT cannot use an id containing '/', '+', or '#' (such devices can still report through a gateway).
- `iot_fleet_id` (String) ID of the IoT Fleet this device belongs to. The ID of a `oneuptime_iot_fleet`.

### Optional

- `is_enabled` (Boolean) Disabled credentials are rejected at MQTT CONNECT and stop the device's silent-death offline detection. Defaults to `true`.
- `name` (String) Any friendly name of this device.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_connected_at` (String) When this credential last authenticated an MQTT connection.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `secret_key` (String) Secret this device presents as the MQTT password (with this credential's ID as the username).
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing iot device credential by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_iot_device_credential.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_iot_device_credential.example <id>
```
