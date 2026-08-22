---
page_title: "oneuptime_monitor_template Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Reusable monitor template. Use it to create new monitors with the same configuration.
---

# oneuptime_monitor_template (Resource)

Reusable monitor template. Use it to create new monitors with the same configuration.

## Example Usage

```terraform
resource "oneuptime_monitor_template" "example" {
  template_name = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
  monitor_name = "Example short text"
  monitor_type = "Manual"
  monitor_steps = [
    {
      monitor_destination      = "https://your-service.example.com"
      monitor_destination_type = "URL"
      request_type             = "GET"
      criteria = [
        {
          name             = "Check if online"
          filter_condition = "All"
          filters = [
            {
              check_on = "Is Online"
            }
          ]
        }
      ]
    }
  ]
}
```

## Schema

### Required

- `template_name` (String) Name of the Monitor Template..
- `template_description` (String) Description of the Monitor Template..
- `monitor_name` (String) Default name applied to monitors created from this template. Users can override on creation...
- `monitor_type` (String) What is the type of monitor created from this template?.. Allowed values: `Manual`, `Website`, `API`, `Ping`, `Kubernetes`, `Docker`, `Host`, `Podman`, `Docker Swarm`, `Proxmox`, `Ceph`, `IoT Device`, `IP`, `Incoming Request`, `Incoming Email`, `Port`, `Server`, `SSL Certificate`, `SQL Query`, `Synthetic Monitor`, `Custom JavaScript Code`, `Logs`, `Metrics`, `Traces`, `Exceptions`, `Profiles`, `Security Events`, `Network Device`, `DNS`, `DNSSEC`, `Domain`, `External Status Page`.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_description` (String) Default description applied to monitors created from this template...
- `monitor_steps` (Block List) MonitorSteps object.
- `monitoring_interval` (String) Default monitoring interval for monitors created from this template..
- `labels` (Set) Default labels applied to monitors created from this template...
- `custom_fields` (String) Custom Fields on this resource...
- `minimum_probe_agreement` (Number) Default minimum number of probes that must agree on a status before the monitor status changes...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor_template.example <id>
```
