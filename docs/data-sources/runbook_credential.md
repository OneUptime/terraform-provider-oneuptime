---
page_title: "oneuptime_runbook_credential Data Source - oneuptime"
subcategory: "Other"
description: |-
  Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step.
---

# oneuptime_runbook_credential (Data Source)

Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_runbook_credential" "by_name" {
  name = "example-runbook_credential"
}

data "oneuptime_runbook_credential" "by_id" {
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
- `credential_type` (String) SSH, or Kubernetes... Computed.
- `ssh_hostname` (String) Hostname or IP address the Runner connects to... Computed.
- `ssh_port` (Number) Defaults to 22 when unset... Computed.
- `ssh_username` (String) The user the Runner authenticates as... Computed.
- `kubernetes_api_server_url` (String) For example https://10.0.0.1:6443.. Computed.
- `kubernetes_ca_certificate` (String) PEM certificate authority for the API server. Leave empty only if the API server presents a certificate your Runner already trusts... Computed.
- `runners` (Set) The Runners allowed to use this credential. A step referencing it must target one of them... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
