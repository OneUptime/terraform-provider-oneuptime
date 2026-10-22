---
page_title: "oneuptime_security_case Data Source - oneuptime"
subcategory: "Other"
description: |-
  A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.
---

# oneuptime_security_case (Data Source)

A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case" "example" {
  case_number_with_prefix = "example-case-number-with-prefix"
}

# Or by id:
data "oneuptime_security_case" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `assigned_to_team_id` (String) ID of the team whose queue the case is in. The ID of a `oneuptime_team`.
- `assigned_to_user_id` (String) ID of the person working the case. The ID of a `oneuptime_user` (see the data source).
- `case_number` (Number) The case's number in this project, counting up from 1 as cases are opened.
- `case_number_with_prefix` (String) The case number with the project's case number prefix in front of it, like CASE-42.
- `closed_by_user_id` (String) ID of the person who closed the case. The ID of a `oneuptime_user` (see the data source).
- `closing_note` (String) Why the case was closed the way it was, in Markdown: what was found and what was done.
- `created_by_user_id` (String) ID of the person who opened the case by hand. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What is known about the threat, in Markdown. A case opened by a detection starts with the detection's description.
- `detection_count` (Number) How many distinct detections the case holds.
- `grouping_key` (String) What new detections are matched on to join this case while it is open: the entity they are about, or the detection's own fingerprint.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the OneUptime incident declared for this case. The ID of a `oneuptime_incident`.
- `notified_priority` (String) The highest priority people have been told about this case at. They are told again only when it rises above it.
- `priority` (String) Critical, High, Medium, Low or Informational. A case opened by a detection starts at the detection's severity and rises when a more severe one joins it.
- `source` (String) What opened the case: a Detection Rule, Threat Intel, a Connected Product, or a person (Manual).
- `status` (String) New, Investigating, Contained or Closed. Every case starts New; closing it records a verdict.
- `title` (String) What the case is about, in a line.
- `verdict` (String) What the case turned out to be, set when it is closed: True Positive, Benign Positive, False Positive or Undetermined.

### Read-Only

- `acknowledged_at` (String) When somebody first took the case on: assigned it to a person, or moved it out of New. The time to acknowledge is measured to here.
- `closed_at` (String) When the case was closed. Cleared if it is reopened. The time to close is measured to here.
- `contained_at` (String) When the case first reached Contained. The time to contain is measured to here.
- `created_at` (String) Date and Time when the object was created.
- `entities` (String) The users, hosts, IP addresses, domains and files the case is about, each as { type, value }. Detections add theirs as they join the case. A JSON value: write it with `jsonencode()`.
- `first_detected_at` (String) When the case's earliest detection fired.
- `labels` (Set of String) Labels to filter and report on cases by. IDs of `oneuptime_label` resources.
- `last_detected_at` (String) When a detection of the case last fired.
- `mitre_tactics` (String) ATT&CK tactic ids (TA0006) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.
- `mitre_techniques` (String) ATT&CK technique ids (T1110) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this case belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
