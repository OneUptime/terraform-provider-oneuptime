---
page_title: "oneuptime_service_ai_note Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.
---

# oneuptime_service_ai_note (Resource)

What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.

## Example Usage

```terraform
resource "oneuptime_service_ai_note" "example" {
  service_id = oneuptime_service.example.id
  note       = "Checkout's database failover takes five minutes; don't page for it."
}
```

## Schema

### Required

- `note` (String) What the team knows about the service that OneUptime AI should read when it investigates it, in plain text, at most 2000 characters.
- `service_id` (String) ID of the service the note is about. The ID of a `oneuptime_service`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who wrote the note. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this note belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service ai note by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_ai_note.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_ai_note.example <id>
```
