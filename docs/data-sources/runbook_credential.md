---
page_title: "oneuptime_runbook_credential Data Source - oneuptime"
subcategory: "Other"
description: |-
  Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step.
---

# oneuptime_runbook_credential (Data Source)

Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one runbook credential may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_runbook_credential" "example" {
  name = "Example runbook credential"
}

# Or by id:
data "oneuptime_runbook_credential" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `credential_type` (String) SSH, or Kubernetes.
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kubernetes_api_server_url` (String) For example https://10.0.0.1:6443.
- `kubernetes_ca_certificate` (String) PEM certificate authority for the API server. Leave empty only if the API server presents a certificate your Runner already trusts.
- `name` (String) Any friendly name of this object.
- `ssh_hostname` (String) Hostname or IP address the Runner connects to.
- `ssh_port` (Number) Defaults to 22 when unset.
- `ssh_username` (String) The user the Runner authenticates as.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runners` (Set of String) The Runners allowed to use this credential. A step referencing it must target one of them. IDs of `oneuptime_runner` resources.
- `updated_at` (String) Date and Time when the object was updated.
