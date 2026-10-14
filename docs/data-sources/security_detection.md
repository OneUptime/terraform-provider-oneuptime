---
page_title: "oneuptime_security_detection Data Source - oneuptime"
subcategory: "Other"
description: |-
  A signal inside a security case: a detection rule's match, a threat intelligence hit, or an alert from a connected security product, with how often it fired.
---

# oneuptime_security_detection (Data Source)

A signal inside a security case: a detection rule's match, a threat intelligence hit, or an alert from a connected security product, with how often it fired.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security detection may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_detection" "example" {
  security_case_id = oneuptime_security_case.example.id
}

# Or by id:
data "oneuptime_security_detection" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `description` (String) What the detection saw, in Markdown.
- `detection_rule_id` (String) ID of the detection rule that fired. The ID of a `oneuptime_detection_rule`.
- `event_count` (Number) How many security events are behind the detection, over every time it fired.
- `event_uid` (String) eventUid of the security event that records the detection: the Detection Finding a rule or feed wrote, or the imported alert.
- `fingerprint` (String) Identifies the same signal across repeats: the source and what it fired for.
- `fire_count` (Number) How many times the detection fired while its case was open.
- `group_value` (String) The value the rule grouped its matches by, such as the host or IP address they share.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `rule_name` (String) The name of what fired: the detection rule, the threat intel feed, or the connected product's own rule.
- `sample_message` (String) The message of one of the events behind the detection.
- `security_case_id` (String) ID of the case this detection belongs to. The ID of a `oneuptime_security_case`.
- `security_event_connection_id` (String) ID of the connection the alert was imported through. The ID of a `oneuptime_security_event_connection`.
- `severity` (String) OCSF severity of the detection: Informational, Low, Medium, High, Critical or Fatal.
- `source` (String) What raised it: a Detection Rule, Threat Intel, a Connected Product, or a person (Manual).
- `threat_intel_feed_id` (String) ID of the feed whose indicator matched. The ID of a `oneuptime_threat_intel_feed`.
- `title` (String) What fired, about what: the rule's name and the value it grouped on.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `entities` (String) The users, hosts, IP addresses, domains and files the detection names, each as { type, value }. A JSON value: write it with `jsonencode()`.
- `first_seen_at` (String) When the detection first fired.
- `last_seen_at` (String) When the detection last fired.
- `mitre_tactics` (String) ATT&CK tactic ids the detection is tagged with. A JSON value: write it with `jsonencode()`.
- `mitre_techniques` (String) ATT&CK technique ids the detection is tagged with. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this detection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
