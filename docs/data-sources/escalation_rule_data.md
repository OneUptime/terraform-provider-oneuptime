---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_escalation_rule_data Data Source - oneuptime"
subcategory: ""
description: |-
  Escalation rule data data source
---

# oneuptime_escalation_rule_data (Data Source)

Escalation rule data data source

## Example Usage

```terraform
data "oneuptime_escalation_rule_data" "example" {
  name = "example-escalation_rule_data"
}
```

## Schema

- `id` (String) Identifier to filter by. Optional.
- `name` (String) Name to filter by. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Escalation Rule], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Escalation Rule], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Duty Policy Escalation Rule]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `escalate_after_in_minutes` (Number) Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Escalation Rule], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Escalation Rule], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Duty Policy Escalation Rule]. Computed.
- `order` (Number) Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Escalation Rule], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Escalation Rule], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Duty Policy Escalation Rule]. Computed.
