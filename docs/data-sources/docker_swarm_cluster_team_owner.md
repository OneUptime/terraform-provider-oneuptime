---
page_title: "oneuptime_docker_swarm_cluster_team_owner Data Source - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your DockerSwarm clusters.
---

# oneuptime_docker_swarm_cluster_team_owner (Data Source)

Add teams as owners to your DockerSwarm clusters.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one docker swarm cluster team owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_docker_swarm_cluster_team_owner" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_docker_swarm_cluster_team_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `docker_swarm_cluster_id` (String) ID of your OneUptime DockerSwarm Cluster in which this object belongs. The ID of a `oneuptime_docker_swarm_cluster`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
