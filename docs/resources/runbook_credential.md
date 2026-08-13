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
  name = "Example short text"
  credential_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `credential_type` (String) SSH, or Kubernetes...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `ssh_hostname` (String) Hostname or IP address the Runner connects to...
- `ssh_port` (Number) Defaults to 22 when unset...
- `ssh_username` (String) The user the Runner authenticates as...
- `ssh_private_key` (String) PEM private key used to authenticate. Encrypted at rest and never returned by the API...
- `ssh_passphrase` (String) Passphrase protecting the private key, when it has one. Encrypted at rest and never returned by the API...
- `ssh_password` (String) Password authentication, for hosts without key access. Encrypted at rest and never returned by the API...
- `kubernetes_api_server_url` (String) For example https://10.0.0.1:6443..
- `kubernetes_service_account_token` (String) Bearer token of a service account bound to a role that permits only the actions your runbooks need. Encrypted at rest and never returned by the API...
- `kubernetes_ca_certificate` (String) PEM certificate authority for the API server. Leave empty only if the API server presents a certificate your Runner already trusts...
- `runners` (Set) The Runners allowed to use this credential. A step referencing it must target one of them...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runbook_credential.example <id>
```
