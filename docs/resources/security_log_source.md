---
page_title: "oneuptime_security_log_source Resource - oneuptime"
subcategory: "Security"
description: |-
  Whether OneUptime looks for security events in the logs the project already sends it, and what the last look found.
---

# oneuptime_security_log_source (Resource)

Whether OneUptime looks for security events in the logs the project already sends it, and what the last look found.

## Example Usage

```terraform
resource "oneuptime_security_log_source" "example" {

}
```

## Schema

### Optional

- `case_minimum_severity` (String) Threats security products write in the logs (CEF, LEEF, Palo Alto Networks, FortiGate, Cisco ASA) at or above this severity (Low, Medium, High or Critical) open security cases. Empty stores them as events only. Defaults to `High`.
- `is_enabled` (Boolean) Whether OneUptime looks for security events in the project's logs. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who turned this source on. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which logs have been read, as an ISO string: the end of the last window the scan finished.
- `id` (String) Unique identifier for the resource.
- `last_error` (String) Why the last scan failed, if it did. Cleared by the next scan that succeeds.
- `last_event_found_at` (String) When a security event was last found in the logs.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_scan_result` (String) What the last scan read and found: the log lines read, the security events found, and how many of each kind. A JSON value: read it with `jsondecode()`.
- `last_scanned_at` (String) When the logs were last read.
- `project_id` (String) ID of the project whose logs are read. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security log source by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_log_source.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_log_source.example <id>
```
