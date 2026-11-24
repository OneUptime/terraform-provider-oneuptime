---
page_title: "oneuptime_protected_entity Resource - oneuptime"
subcategory: "Security"
description: |-
  Users, hosts and IP addresses that security playbooks never contain and the blocklist never blocks, such as break-glass accounts, domain controllers, VPN gateways and office addresses.
---

# oneuptime_protected_entity (Resource)

Users, hosts and IP addresses that security playbooks never contain and the blocklist never blocks, such as break-glass accounts, domain controllers, VPN gateways and office addresses.

## Example Usage

```terraform
resource "oneuptime_protected_entity" "example" {
  entity_type = "User"
  value       = "This is an example of longer text content that might be stored in this field."
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `entity_type` (String) User, Host or IP Address: what the value names. Allowed values: `User`, `Host`, `IP Address`, `Domain`, `URL`, `File Hash`, `Process`, `Data Store`, `Other`.
- `value` (String) The user's sign-in name or the host's name, matched whatever the case. * stands for any run of characters: svc-* protects every service account. For an IP Address, an IPv4 or IPv6 address or a range such as 203.0.113.0/24, which protects every address inside it.

### Optional

- `description` (String) Why it is protected, for whoever reads the list later.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who protected it. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this protected entity belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing protected entity by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_protected_entity.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_protected_entity.example <id>
```
