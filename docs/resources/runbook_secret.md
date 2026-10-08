---
page_title: "oneuptime_runbook_secret Resource - oneuptime"
subcategory: "Other"
description: |-
  Runbook Secret is a secret variable that can be used by runbook agents. For example you can store auth tokens, passwords, etc. in Runbook Secret and use them in your runbook steps. Runbook Secret is encrypted and only accessible by the assigned agent.
---

# oneuptime_runbook_secret (Resource)

Runbook Secret is a secret variable that can be used by runbook agents. For example you can store auth tokens, passwords, etc. in Runbook Secret and use them in your runbook steps. Runbook Secret is encrypted and only accessible by the assigned agent.

## Example Usage

```terraform
resource "oneuptime_runbook_secret" "example" {
  name        = "Example runbook secret"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `runners` (Set of String) List of runbook agents that can access this secret. IDs of `oneuptime_runner` resources.
- `secret_value` (String) Secret value that you want to store in this object. This value will be encrypted and only accessible by the assigned runbook agent.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook secret by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook_secret.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook_secret.example <id>
```
