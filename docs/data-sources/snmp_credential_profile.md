---
page_title: "oneuptime_snmp_credential_profile Data Source - oneuptime"
subcategory: "Other"
description: |-
  A reusable set of SNMP credentials. Attach a profile to a device or to a site and every device it covers is walked over SNMP with these credentials instead of carrying its own.
---

# oneuptime_snmp_credential_profile (Data Source)

A reusable set of SNMP credentials. Attach a profile to a device or to a site and every device it covers is walked over SNMP with these credentials instead of carrying its own. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_snmp_credential_profile" "by_name" {
  name = "example-snmp_credential_profile"
}

data "oneuptime_snmp_credential_profile" "by_id" {
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
- `description` (String) Friendly description that will help you remember.. Computed.
- `snmp_version` (String) SNMP version devices using this profile are polled with (V1, V2c, V3).. Computed.
- `snmp_community_string` (String) Community string used for SNMP v1/v2c polling.. Computed.
- `snmp_port` (Number) UDP port used for SNMP polling.. Computed.
- `snmp_v3_security_level` (String) SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv.. Computed.
- `snmp_v3_username` (String) Security name (username) used for SNMP v3 polling.. Computed.
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.. Computed.
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase.. Computed.
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.. Computed.
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
