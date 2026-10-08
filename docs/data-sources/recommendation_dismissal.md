---
page_title: "oneuptime_recommendation_dismissal Data Source - oneuptime"
subcategory: "Other"
description: |-
  Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.
---

# oneuptime_recommendation_dismissal (Data Source)

Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one recommendation dismissal may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_recommendation_dismissal" "example" {
  recommendation_type = "example-recommendation-type"
}

# Or by id:
data "oneuptime_recommendation_dismissal" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dismissal_reason` (String) Optional note explaining why this recommendation was dismissed, shown to whoever finds it in the dismissed list later.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `recommendation_id` (String) The catalog-wide id of the dismissed recommendation, for example Kubernetes:k8s-hpa-at-max-replicas.
- `recommendation_type` (String) Which family of recommendation this dismissal belongs to. See the RecommendationType enum.
- `resource_id` (String) ID of the resource this recommendation was shown on. Polymorphic — it points at whichever table Resource Type names — so it carries no foreign key.
- `resource_type` (String) The kind of resource this recommendation was shown on, for example Kubernetes or Docker. Empty for recommendations that are not scoped to a resource.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
