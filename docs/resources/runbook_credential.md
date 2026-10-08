---
page_title: "oneuptime_runbook_credential Resource - oneuptime"
subcategory: "Other"
description: |-
  Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step.
---

# oneuptime_runbook_credential (Resource)

Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step.

## Example Usage

```terraform
resource "oneuptime_runbook_credential" "example" {
  name            = "Example runbook credential"
  credential_type = "Example short text"
  description     = "Managed by Terraform"
}
```

## Schema

### Required

- `credential_type` (String) SSH, or Kubernetes.
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `kubernetes_api_server_url` (String) For example https://10.0.0.1:6443.
- `kubernetes_ca_certificate` (String) PEM certificate authority for the API server. Leave empty only if the API server presents a certificate your Runner already trusts.
- `kubernetes_service_account_token` (String) Bearer token of a service account bound to a role that permits only the actions your runbooks need. Encrypted at rest and never returned by the API.
- `runners` (Set of String) The Runners allowed to use this credential. A step referencing it must target one of them. IDs of `oneuptime_runner` resources.
- `ssh_hostname` (String) Hostname or IP address the Runner connects to.
- `ssh_passphrase` (String) Passphrase protecting the private key, when it has one. Encrypted at rest and never returned by the API.
- `ssh_password` (String) Password authentication, for hosts without key access. Encrypted at rest and never returned by the API.
- `ssh_port` (Number) Defaults to 22 when unset.
- `ssh_private_key` (String) PEM private key used to authenticate. Encrypted at rest and never returned by the API.
- `ssh_username` (String) The user the Runner authenticates as.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook credential by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook_credential.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook_credential.example <id>
```
