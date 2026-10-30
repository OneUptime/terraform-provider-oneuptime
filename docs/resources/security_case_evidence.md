---
page_title: "oneuptime_security_case_evidence Resource - oneuptime"
subcategory: "Other"
description: |-
  Files, links and security events kept on a security case as evidence, each with who added it and when: a file with the SHA-256 of what was uploaded, an event as OneUptime read it when it was kept.
---

# oneuptime_security_case_evidence (Resource)

Files, links and security events kept on a security case as evidence, each with who added it and when: a file with the SHA-256 of what was uploaded, an event as OneUptime read it when it was kept.

## Example Usage

```terraform
resource "oneuptime_security_case_evidence" "example" {
  security_case_id = oneuptime_security_case.example.id
  kind             = "Example short text"
  name             = "Example security case evidence"
  description      = "Managed by Terraform"
}
```

## Schema

### Required

- `kind` (String) What the evidence is: a File uploaded to OneUptime, a Link to evidence kept somewhere else, or an Event bookmarked from Security Events. It can't be changed.
- `name` (String) What the evidence is, in a few words. A file's own name when none is given.
- `security_case_id` (String) ID of the case this evidence is kept on. The ID of a `oneuptime_security_case`.

### Optional

- `description` (String) What the evidence shows, and where and how it was collected.
- `file_id` (String) ID of the uploaded file, for File evidence: upload it first, then name it here. The ID of a `oneuptime_file`.
- `url` (String) Where the evidence is kept, for Link evidence: an http or https address.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added the evidence. The ID of a `oneuptime_user` (see the data source).
- `event_snapshot` (String) For Event evidence: the security event as OneUptime read it when it was kept - its severity, class, message, who and where - and the hunt that found it. A JSON value: write it with `jsonencode()`.
- `event_time` (String) For Event evidence: when the security event happened.
- `event_uid` (String) For Event evidence: the security event's ID, as its source gave it.
- `file_name` (String) The file's name as it was uploaded, for File evidence.
- `file_size_in_bytes` (Number) The file's size as it was uploaded, in bytes.
- `file_type` (String) The file's type as it was uploaded, like text/plain.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this evidence belongs to. The ID of a `oneuptime_project`.
- `sha256` (String) The SHA-256 of the file's bytes as they were uploaded, in hex: compare it with a download's to check it is what was kept.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case evidence by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_evidence.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_evidence.example <id>
```
