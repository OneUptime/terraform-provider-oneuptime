---
page_title: "oneuptime_security_case Data Source - oneuptime"
subcategory: "Security"
description: |-
  A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.
---

# oneuptime_security_case (Data Source)

A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case" "example" {
  last_modified_by_user_id = data.oneuptime_user.example.id
  assigned_to_user_email   = "alice@example.com"
}

# Or by id:
data "oneuptime_security_case" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_hand_off_reason` (String) Why OneUptime AI last handed the case to people, such as SteeringText, HighestPriority or ProtectedTarget.
- `ai_hold_ended_by_user_id` (String) The person who ended OneUptime AI's hold of the case's alerts, by pressing Page now or taking the case. Empty when OneUptime ended it.
- `ai_hold_outcome` (String) Whether OneUptime AI is responding first to the case before anyone is alerted (Holding), and how that ended: Resolved with nobody alerted, handed to people (NeedsApproval, HandedOff, TimeRanOut and the rest), or stopped by a person (PageNow, TakenByPerson) or a priority rise (RoseInPriority). Empty for a case never held.
- `ai_response_state` (String) Where OneUptime AI stands on the case when it works it: Working, WaitingForPeople, HandedToPeople or Finished. Empty for a case OneUptime AI never worked.
- `ai_suggested_playbook_id` (String) ID of the project's playbook OneUptime AI suggested running when it triaged the case. It runs only when a person runs it.
- `ai_suggested_security_case_priority_id` (String) ID of the priority OneUptime AI suggested when it triaged the case. The case's own priority changes only when a person sets it.
- `ai_triage_prompt_injection_suspected` (Boolean) Whether OneUptime AI found text in the case's data written to steer an AI when it triaged the case. It then leaves the call to a person.
- `ai_triage_report` (String) The report OneUptime AI wrote when it last triaged the case: what happened, its read and the evidence for it, and the next steps.
- `ai_verdict` (String) What OneUptime AI made of the case when it triaged it: Likely Malicious, Likely Benign, Likely False Positive or Inconclusive. A suggestion; a person sets the case's verdict.
- `assigned_to_team_id` (String) ID of the team whose queue the case is in. The ID of a `oneuptime_team`.
- `assigned_to_user_email` (String) Look it up by the email address of the person in `assigned_to_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `assigned_to_user_id` (String) ID of the person working the case. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `assigned_to_user_email`.
- `case_number` (Number) The case's number in this project, counting up from 1 as cases are opened.
- `case_number_with_prefix` (String) The case number with the project's case number prefix in front of it, like CASE-42.
- `closed_by_user_id` (String) ID of the person who closed the case. The ID of a `oneuptime_user` (see the data source).
- `closing_note` (String) Why the case was closed the way it was, in Markdown: what was found and what was done.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who opened the case by hand. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What is known about the threat, in Markdown. A case opened by a detection starts with the detection's description.
- `detection_count` (Number) How many distinct detections the case holds.
- `grouping_key` (String) What new detections are matched on to join this case while it is open: the entity they are about, or the detection's own fingerprint.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the OneUptime incident declared for this case. The ID of a `oneuptime_incident`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `notified_security_case_priority_id` (String) ID of the highest priority people have been told about this case at. They are told again only when it rises above it.
- `security_case_priority_id` (String) ID of the case's priority. Optional on create: a case opened by hand starts at the project's medium priority. The ID of a `oneuptime_security_case_priority`.
- `security_case_status_id` (String) ID of the case's status. Optional on create: a new case starts in the project's new status. The ID of a `oneuptime_security_case_status`.
- `source` (String) What opened the case: a Detection Rule, Threat Intel, a Connected Product, Entity Risk, or a person (Manual).
- `time_to_block_in_seconds` (Number) Seconds from the address's first malicious request to an automatic response policy's block of it taking hold.
- `title` (String) What the case is about, in a line.
- `verdict` (String) What the case turned out to be, set when it is closed: True Positive, Benign Positive, False Positive or Undetermined.

### Read-Only

- `acknowledge_due_at` (String) When the case must be acknowledged by, from its priority's Acknowledge Within target. Empty when its priority has no target, or when it was acknowledged as it opened.
- `acknowledge_missed_at` (String) When the case passed its acknowledge due time without being acknowledged. Empty while it has not.
- `acknowledge_warned_at` (String) When OneUptime told the people on the case that it was about to miss its acknowledge target.
- `acknowledged_at` (String) When somebody first took the case on: assigned it to a person, or moved it out of New. The time to acknowledge is measured to here.
- `ai_closed_at` (String) When OneUptime AI closed the case, on its own or once a person approved it. Empty for a case it did not close.
- `ai_hold_ended_at` (String) When OneUptime AI's hold of the case's alerts ended. Empty while it holds them, and for a case never held.
- `ai_hold_until` (String) When people are alerted about the case at the latest while OneUptime AI responds first to it. Empty when it is not held.
- `ai_next_look_at` (String) When OneUptime AI next looks at the case while it works it. Empty when it is not working the case.
- `ai_took_case_at` (String) When OneUptime AI took the case to work it. Empty for a case it never took.
- `ai_triaged_at` (String) When OneUptime AI last finished triaging the case.
- `blocked_at` (String) When an automatic response policy's block of one of the case's addresses first took hold at an edge.
- `closed_at` (String) When the case was closed. Cleared if it is reopened. The time to close is measured to here.
- `contain_due_at` (String) When the case must be contained or closed by, from its priority's Contain Within target. Empty when its priority has no target.
- `contain_missed_at` (String) When the case passed its contain due time without being contained or closed. Empty while it has not.
- `contain_warned_at` (String) When OneUptime told the people on the case that it was about to miss its contain target.
- `contained_at` (String) When the case first reached Contained. The time to contain is measured to here.
- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) The values of the project's own security case fields (Security Events > Settings > Custom Fields), by each field's name. A JSON value: read it with `jsondecode()`.
- `entities` (String) The users, hosts, IP addresses, domains and files the case is about, each as { type, value }. Detections add theirs as they join the case. A JSON value: read it with `jsondecode()`.
- `first_detected_at` (String) When the case's earliest detection fired.
- `labels` (Set of String) Labels to filter and report on cases by. IDs of `oneuptime_label` resources.
- `last_detected_at` (String) When a detection of the case last fired.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `mitre_tactics` (String) ATT&CK tactic ids (TA0006) the case's detections are tagged with. A JSON value: read it with `jsondecode()`.
- `mitre_techniques` (String) ATT&CK technique ids (T1110) the case's detections are tagged with. A JSON value: read it with `jsondecode()`.
- `next_due_at` (String) The earliest due time of a response target the case has not met yet, a missed one included. Empty when every target is met or the case has none.
- `next_step_check_at` (String) When OneUptime next checks the steps of the case's Response Plan itself. Empty once the case is closed.
- `next_target_check_at` (String) When OneUptime next checks the case's response targets: when one is about to be missed, or is missed. Empty when there is nothing left to check.
- `project_id` (String) ID of the project this case belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
