---
page_title: "oneuptime_snmp_credential_profile Resource - oneuptime"
subcategory: "Network"
description: |-
  A reusable set of SNMP credentials. Attach a profile to a device or to a site and every device it covers is walked over SNMP with these credentials instead of carrying its own.
---

# oneuptime_snmp_credential_profile (Resource)

A reusable set of SNMP credentials. Attach a profile to a device or to a site and every device it covers is walked over SNMP with these credentials instead of carrying its own.

## Example Usage

```terraform
resource "oneuptime_snmp_credential_profile" "example" {
  name        = "Example snmp credential profile"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `snmp_community_string` (String, Sensitive) Community string used for SNMP v1/v2c polling.
- `snmp_port` (Number) UDP port used for SNMP polling.
- `snmp_v3_auth_key` (String, Sensitive) SNMP v3 authentication passphrase.
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.
- `snmp_v3_priv_key` (String, Sensitive) SNMP v3 privacy (encryption) passphrase.
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.
- `snmp_v3_security_level` (String) SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv.
- `snmp_v3_username` (String) Security name (username) used for SNMP v3 polling.
- `snmp_version` (String) SNMP version devices using this profile are polled with (V1, V2c, V3).

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

Import an existing snmp credential profile by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_snmp_credential_profile.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_snmp_credential_profile.example <id>
```
