---
page_title: "oneuptime_docker_swarm_cluster_label_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to DockerSwarm clusters when matching DockerSwarm clusters are created
---

# oneuptime_docker_swarm_cluster_label_rule (Data Source)

Configure rules for automatically attaching labels to DockerSwarm clusters when matching DockerSwarm clusters are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_docker_swarm_cluster_label_rule" "by_name" {
  name = "example-docker_swarm_cluster_label_rule"
}

data "oneuptime_docker_swarm_cluster_label_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this DockerSwarm cluster label rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `docker_swarm_cluster_labels` (Set) Only trigger for DockerSwarm clusters that already have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `docker_swarm_cluster_name_pattern` (String) Regex (case-insensitive) matched against the DockerSwarm cluster name. Leave empty to match any name... Computed.
- `docker_swarm_cluster_description_pattern` (String) Regex (case-insensitive) matched against the DockerSwarm cluster description. Leave empty to match any description... Computed.
- `labels_to_add` (Set) Labels to attach to the DockerSwarm cluster when this rule matches. Already-attached labels are not duplicated... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
