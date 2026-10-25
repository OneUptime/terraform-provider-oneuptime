---
page_title: "oneuptime_security_case_priority Resource - oneuptime"
subcategory: "Other"
description: |-
  How urgent a security case is, most urgent first (Critical, High, Medium, Low, Informational and your own). Add, rename, reorder or archive them.
---

# oneuptime_security_case_priority (Resource)

How urgent a security case is, most urgent first (Critical, High, Medium, Low, Informational and your own). Add, rename, reorder or archive them.

## Example Usage

```terraform
resource "oneuptime_security_case_priority" "example" {
  name        = "Example security case priority"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `color` (String) The priority's colour, in hex (#32a852 for example).
- `name` (String) What the priority is called, like High or P1.

### Optional

- `description` (String) How fast a case at this priority needs a person, in a sentence.
- `is_archived` (Boolean) Archived priorities stay on the cases at them, but can't be picked again. The built-in priorities can't be archived. Defaults to `false`.
- `is_critical_priority` (Boolean) Is it the priority Critical detections open cases at? From it up, a case can't wait. A project has one; it can't be deleted or archived.
- `is_high_priority` (Boolean) Is it the priority High detections open cases at? From it up, a case counts as urgent. A project has one; it can't be deleted or archived.
- `is_informational_priority` (Boolean) Is it the priority Informational detections open cases at? A project has one, and it can't be deleted or archived.
- `is_low_priority` (Boolean) Is it the priority Low detections open cases at? A project has one, and it can't be deleted or archived.
- `is_medium_priority` (Boolean) Is it the priority Medium detections, and the cases people open, start at? A project has one; it can't be deleted or archived.
- `order` (Number) Where this priority ranks among the project's security case priorities: 1 is the most urgent. Detections open their cases at the critical, high, medium, low and informational priorities, so those have to stay in that order. A new priority without a number goes to the end of the list. Setting a number moves the priority to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `archived_at` (String) When this priority was archived. Empty while it is not archived.
- `archived_by_user_id` (String) ID of the person who archived this priority. The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added this priority. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this priority belongs to. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for the priority.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case priority by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_priority.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_priority.example <id>
```
