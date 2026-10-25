---
page_title: "oneuptime_entity_criticality Resource - oneuptime"
subcategory: "Other"
description: |-
  How much each user, host and IP address matters: the risk of an entity an entry names counts less (Low) or more (High, Extreme).
---

# oneuptime_entity_criticality (Resource)

How much each user, host and IP address matters: the risk of an entity an entry names counts less (Low) or more (High, Extreme).

## Example Usage

```terraform
resource "oneuptime_entity_criticality" "example" {
  value       = "This is an example of longer text content that might be stored in this field."
  criticality = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `criticality` (String) How much it matters: Low (its risk counts half), High (one and a half times) or Extreme (twice).
- `value` (String) The user's sign-in name, the host's name or the address, matched whatever the case. * stands for any run of characters: dc-* names every domain controller.

### Optional

- `description` (String) Why it matters this much, for whoever reads the list later.
- `entity_type` (String) User, Host or IP Address: the kind of entity it is for. Empty is any kind.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added the entry. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this entry belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing entity criticality by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_entity_criticality.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_entity_criticality.example <id>
```
