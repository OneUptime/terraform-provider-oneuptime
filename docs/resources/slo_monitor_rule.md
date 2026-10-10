---
page_title: "oneuptime_slo_monitor_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand
---

# oneuptime_slo_monitor_rule (Resource)

Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand

## Example Usage

```terraform
resource "oneuptime_slo_monitor_rule" "example" {
  service_level_objective_id = oneuptime_service_level_objective.example.id
  name                       = "Example slo monitor rule"
  description                = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this SLO monitor rule.
- `service_level_objective_id` (String) ID of the Service Level Objective this monitor rule belongs to. The ID of a `oneuptime_service_level_objective`.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this SLO monitor rule.
- `is_enabled` (Boolean) Whether this rule is enabled. A disabled rule matches nothing, so monitors that only this rule attached are detached from the SLO. Defaults to `true`.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter.
- `monitor_labels` (Set of String) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor.
- `monitor_type` (String) Only match monitors of this type. Leave empty to skip the type filter. Allowed values: `Manual`, `Website`, `API`, `Ping`, `Kubernetes`, `Docker`, `Host`, `Podman`, `Docker Swarm`, `Proxmox`, `VMware`, `Ceph`, `Storage Array`, `IoT Device`, `IP`, `Incoming Request`, `Incoming Email`, `Port`, `Server`, `SSL Certificate`, `SQL Query`, `Database`, `Synthetic Monitor`, `Custom JavaScript Code`, `Logs`, `Metrics`, `Traces`, `Exceptions`, `Profiles`, `Security Events`, `AI / LLM`, `Network Device`, `DNS`, `DNSSEC`, `NTP`, `Domain`, `External Status Page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing slo monitor rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_slo_monitor_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_slo_monitor_rule.example <id>
```
