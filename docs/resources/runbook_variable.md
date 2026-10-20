---
page_title: "oneuptime_runbook_variable Resource - oneuptime"
subcategory: "Other"
description: |-
  A value runbook steps use: a host, a namespace or a URL, or a secret such as an SSH key, a password, a kubeconfig or an API token. It belongs to the project, which every runbook can use, or to one runbook. A secret value is encrypted and can never be read back through the API; it reaches only the Runners it lets use it, when they take a step that names it.
---

# oneuptime_runbook_variable (Resource)

A value runbook steps use: a host, a namespace or a URL, or a secret such as an SSH key, a password, a kubeconfig or an API token. It belongs to the project, which every runbook can use, or to one runbook. A secret value is encrypted and can never be read back through the API; it reaches only the Runners it lets use it, when they take a step that names it.

## Example Usage

```terraform
resource "oneuptime_runbook_variable" "example" {
  name        = "Example runbook variable"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Letters, numbers, hyphens and underscores. Steps use the variable by this name: {{global.variables.NAME}} for a global variable, {{local.variables.NAME}} for one of the runbook's own.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_available_to_global_runners` (Boolean) Whether steps of this project that run on a global Runner get this secret. Global Runners are the ones OneUptime runs for every project; a step of another project never gets it. Defaults to `false`.
- `is_secret` (Boolean) A secret is encrypted, can never be read back, and is hidden in step output. It cannot be changed later. Defaults to `false`.
- `runbook_id` (String) ID of the runbook this variable belongs to. Empty for a global variable, which every runbook of the project can use. The ID of a `oneuptime_runbook`.
- `runners` (Set of String) The Runners of your own that may use this secret. A step that names it must run on one of them. IDs of `oneuptime_runner` resources.
- `secret_value` (String, Sensitive) The value of a secret variable. Encrypted at rest and never returned by the API.
- `value` (String) The value of a variable that is not a secret. A secret's value is Secret Value.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook variable by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook_variable.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook_variable.example <id>
```
