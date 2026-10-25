---
page_title: "oneuptime_entity_risk Data Source - oneuptime"
subcategory: "Other"
description: |-
  The risk of each user, host and IP address: what the detection rules, threat intel and connected products found about it in the last 7 days adds up to.
---

# oneuptime_entity_risk (Data Source)

The risk of each user, host and IP address: what the detection rules, threat intel and connected products found about it in the last 7 days adds up to.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one entity risk may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_entity_risk" "example" {
  entity_type = "example-entity-type"
}

# Or by id:
data "oneuptime_entity_risk" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `contributor_count` (Number) How many different detection rules, threat intel feeds and connected products' rules add to its risk.
- `criticality` (String) How much it matters, from the project's Entity Criticality list: Low (x0.5), High (x1.5) or Extreme (x2). Empty counts as it is.
- `entity_key` (String) The entity's value, trimmed and in lower case: one entity's risk, whatever case its findings write it in.
- `entity_type` (String) What the entity is: a User, a Host or an IP Address.
- `entity_value` (String) The user's sign-in name, the host's name or the address, as the first finding about it named it.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `risk_score` (Number) What its sources of the last 7 days add up to, weighted by how much it matters. A case opens at the project's threshold.

### Read-Only

- `contributions` (String) What adds to its risk: each source with its severity, its points, how often it fired and when. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `first_risk_at` (String) When the sources that count now began counting.
- `last_risk_at` (String) When a finding last added to its risk.
- `mitre_tactics` (String) ATT&CK tactic ids (TA0006) its sources are tagged with. A JSON value: write it with `jsonencode()`.
- `mitre_techniques` (String) ATT&CK technique ids (T1110) its sources are tagged with. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this entity's risk belongs to. The ID of a `oneuptime_project`.
- `risk_case_at` (String) When its risk last reached the threshold and opened a case, or joined the open case about it.
- `risk_expires_at` (String) When its score next falls, as the first of its sources stops counting.
- `risk_history` (String) Its risk score at the end of each of the last 30 days. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
