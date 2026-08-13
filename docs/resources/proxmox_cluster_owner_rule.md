---
page_title: "oneuptime_proxmox_cluster_owner_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching Proxmox clusters are created
---

# oneuptime_proxmox_cluster_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching Proxmox clusters are created

## Example Usage

```terraform
resource "oneuptime_proxmox_cluster_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Proxmox cluster owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this Proxmox cluster owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `proxmox_cluster_labels` (Set) Only trigger for Proxmox clusters that have at least one of these labels. Leave empty to match regardless of labels...
- `proxmox_cluster_name_pattern` (String) Regex (case-insensitive) matched against the Proxmox cluster name. Leave empty to match any name...
- `proxmox_cluster_description_pattern` (String) Regex (case-insensitive) matched against the Proxmox cluster description. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the Proxmox cluster when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the Proxmox cluster when this rule matches...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_proxmox_cluster_owner_rule.example <id>
```
