---
page_title: "oneuptime_snmp_credential_profile Resource - oneuptime"
subcategory: "Other"
description: |-
  A reusable set of SNMP credentials. Attach a profile to a device or to a site and every device it covers is walked over SNMP with these credentials instead of carrying its own.
---

# oneuptime_snmp_credential_profile (Resource)

A reusable set of SNMP credentials. Attach a profile to a device or to a site and every device it covers is walked over SNMP with these credentials instead of carrying its own.

## Example Usage

```terraform
resource "oneuptime_snmp_credential_profile" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `snmp_version` (String) SNMP version devices using this profile are polled with (V1, V2c, V3)..
- `snmp_community_string` (String) Community string used for SNMP v1/v2c polling..
- `snmp_port` (Number) UDP port used for SNMP polling..
- `snmp_v3_security_level` (String) SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv..
- `snmp_v3_username` (String) Security name (username) used for SNMP v3 polling..
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512..
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase..
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256..
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_snmp_credential_profile.example <id>
```
