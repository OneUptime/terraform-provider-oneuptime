---
page_title: "oneuptime_security_case Resource - oneuptime"
subcategory: "Other"
description: |-
  A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.
---

# oneuptime_security_case (Resource)

A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.

## Example Usage

```terraform
resource "oneuptime_security_case" "example" {
  title       = "This is an example of longer text content that might be stored in this field."
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `title` (String) What the case is about, in a line.

### Optional

- `assigned_to_team_id` (String) ID of the team whose queue the case is in. The ID of a `oneuptime_team`.
- `assigned_to_user_id` (String) ID of the person working the case. The ID of a `oneuptime_user` (see the data source).
- `closing_note` (String) Why the case was closed the way it was, in Markdown: what was found and what was done.
- `description` (String) What is known about the threat, in Markdown. A case opened by a detection starts with the detection's description.
- `entities` (String) The users, hosts, IP addresses, domains and files the case is about, each as { type, value }. Detections add theirs as they join the case. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Labels to filter and report on cases by. IDs of `oneuptime_label` resources.
- `mitre_tactics` (String) ATT&CK tactic ids (TA0006) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.
- `mitre_techniques` (String) ATT&CK technique ids (T1110) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.
- `priority` (String) Critical, High, Medium, Low or Informational. A case opened by a detection starts at the detection's severity and rises when a more severe one joins it. Defaults to `Medium`.
- `status` (String) New, Investigating, Contained or Closed. Every case starts New; closing it records a verdict. Defaults to `New`.
- `verdict` (String) What the case turned out to be, set when it is closed: True Positive, Benign Positive, False Positive or Undetermined.

### Read-Only

- `acknowledged_at` (String) When somebody first took the case on: assigned it to a person, or moved it out of New. The time to acknowledge is measured to here.
- `case_number` (Number) The case's number in this project, counting up from 1 as cases are opened.
- `case_number_with_prefix` (String) The case number with the project's case number prefix in front of it, like CASE-42.
- `closed_at` (String) When the case was closed. Cleared if it is reopened. The time to close is measured to here.
- `closed_by_user_id` (String) ID of the person who closed the case. The ID of a `oneuptime_user` (see the data source).
- `contained_at` (String) When the case first reached Contained. The time to contain is measured to here.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who opened the case by hand. The ID of a `oneuptime_user` (see the data source).
- `detection_count` (Number) How many distinct detections the case holds.
- `first_detected_at` (String) When the case's earliest detection fired.
- `grouping_key` (String) What new detections are matched on to join this case while it is open: the entity they are about, or the detection's own fingerprint.
- `id` (String) Unique identifier for the resource.
- `incident_id` (String) ID of the OneUptime incident declared for this case. The ID of a `oneuptime_incident`.
- `last_detected_at` (String) When a detection of the case last fired.
- `project_id` (String) ID of the project this case belongs to. The ID of a `oneuptime_project`.
- `source` (String) What opened the case: a Detection Rule, Threat Intel, a Connected Product, or a person (Manual).
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case.example <id>
```
