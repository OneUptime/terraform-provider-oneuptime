---
page_title: "oneuptime_security_case_priority Resource - oneuptime"
subcategory: "Security"
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

<a id="order"></a>
## Order

`order` is this security case priority's place among the security case priorities of its project, lowest first. Setting it places this security case priority and moves no other, so a configuration that numbers every security case priority of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.security_case_priorities`:

```terraform
order = index(local.security_case_priorities, each.key) + 1
```

- Give each security case priority a number of its own: two security case priorities at one number are shown the older first.
- A security case priority created without `order` goes after the last one of its project.
- The critical priority, high priority, medium priority, low priority and informational priority stay in that order, from the top: a write that would put two of them at one number, or out of order, is refused. An apply that moves them past each other's numbers can meet that refusal on one write; apply again and it finishes.
- Deleting a security case priority leaves a gap in the numbers, which changes nothing: security case priorities go by their numbers, not by counting.
- Moving a security case priority in the OneUptime Dashboard numbers the project's security case priorities 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `color` (String) The priority's colour, in hex (#32a852 for example).
- `name` (String) What the priority is called, like High or P1.

### Optional

- `acknowledge_within_in_minutes` (Number) How long a case at this priority may wait for somebody to take it on, in minutes, from when it opens. Empty for no target.
- `contain_within_in_minutes` (Number) How long a case at this priority may take to be contained or closed, in minutes, from when it opens. Empty for no target.
- `description` (String) How fast a case at this priority needs a person, in a sentence.
- `is_archived` (Boolean) Archived priorities stay on the cases at them, but can't be picked again. The built-in priorities can't be archived. Defaults to `false`.
- `is_critical_priority` (Boolean) Is it the priority Critical detections open cases at? From it up, a case can't wait. A project has one; it can't be deleted or archived.
- `is_high_priority` (Boolean) Is it the priority High detections open cases at? From it up, a case counts as urgent. A project has one; it can't be deleted or archived.
- `is_informational_priority` (Boolean) Is it the priority Informational detections open cases at? A project has one, and it can't be deleted or archived.
- `is_low_priority` (Boolean) Is it the priority Low detections open cases at? A project has one, and it can't be deleted or archived.
- `is_medium_priority` (Boolean) Is it the priority Medium detections, and the cases people open, start at? A project has one; it can't be deleted or archived.
- `order` (Number) Where this priority ranks among the project's security case priorities: 1 is the most urgent. Detections open their cases at the critical, high, medium, low and informational priorities, so those have to stay in that order. A new priority without a number goes to the end of the list. Setting a number places the priority there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

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
