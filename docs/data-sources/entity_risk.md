---
page_title: "oneuptime_entity_risk Data Source - oneuptime"
subcategory: "Security"
description: |-
  The risk of each user, host, IP address and data store: what the detection rules, threat intel, connected products and Data Security found about it in the last 7 days adds up to.
---

# oneuptime_entity_risk (Data Source)

The risk of each user, host, IP address and data store: what the detection rules, threat intel, connected products and Data Security found about it in the last 7 days adds up to.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one entity risk may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_entity_risk" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_entity_risk" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `contributor_count` (Number) How many different detection rules, threat intel feeds and connected products' rules add to its risk.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `criticality` (String) How much it matters, from the project's Entity Criticality list: Low (x0.5), High (x1.5) or Extreme (x2). Empty counts as it is.
- `entity_key` (String) The entity's value, trimmed and in lower case: one entity's risk, whatever case its findings write it in.
- `entity_type` (String) What the entity is: a User, a Host, an IP Address or a Data Store.
- `entity_value` (String) The user's sign-in name, the host's name or the address, as the first finding about it named it.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `risk_score` (Number) What its sources of the last 7 days add up to, weighted by how much it matters. A case opens at the project's threshold.

### Read-Only

- `contributions` (String) What adds to its risk: each source with its severity, its points, how often it fired and when. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `first_risk_at` (String) When the sources that count now began counting.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_risk_at` (String) When a finding last added to its risk.
- `mitre_tactics` (String) ATT&CK tactic ids (TA0006) its sources are tagged with. A JSON value: read it with `jsondecode()`.
- `mitre_techniques` (String) ATT&CK technique ids (T1110) its sources are tagged with. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this entity's risk belongs to. The ID of a `oneuptime_project`.
- `risk_case_at` (String) When its risk last reached the threshold and opened a case, or joined the open case about it.
- `risk_expires_at` (String) When its score next falls, as the first of its sources stops counting.
- `risk_history` (String) When its risk score changed over the last 30 days, and to what: each day's score is read from it in the reader's time zone. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
