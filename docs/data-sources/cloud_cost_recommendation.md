---
page_title: "oneuptime_cloud_cost_recommendation Data Source - oneuptime"
subcategory: "Cloud"
description: |-
  Ways to spend less in your connected cloud accounts - rightsizing, idle resources, newer generations and commitments - with what each saves a month, from AWS Cost Optimization Hub, Azure Advisor, Google Cloud's Recommender and OneUptime's own CPU and memory readings. Written by the cloud connections; not user-editable.
---

# oneuptime_cloud_cost_recommendation (Data Source)

Ways to spend less in your connected cloud accounts - rightsizing, idle resources, newer generations and commitments - with what each saves a month, from AWS Cost Optimization Hub, Azure Advisor, Google Cloud's Recommender and OneUptime's own CPU and memory readings. Written by the cloud connections; not user-editable.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one cloud cost recommendation may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_cloud_cost_recommendation" "example" {
  cloud_account_id = oneuptime_cloud_account.example.id
}

# Or by id:
data "oneuptime_cloud_cost_recommendation" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action` (String) What to do, in the recommender's words: 'Rightsize', 'Delete the unattached disk'.
- `cloud_account_id` (String) ID of the account the recommendation is for. The ID of a `oneuptime_cloud_account`.
- `cloud_resource_id` (String) ID of the resource the recommendation is about, when OneUptime holds it. The ID of a `oneuptime_cloud_resource`.
- `created_by_user_id` (String) ID of the user who created this row (the connection writes as root).
- `currency` (String) The currency the figures are in.
- `current_configuration` (String) What the resource is now, in the recommender's words.
- `effort` (String) How much work it is, in the recommender's words: 'Low', 'Medium', 'High'.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What it asks: rightsize, stop or delete something idle, upgrade to a newer generation, commit to use, or other.
- `monthly_cost` (Number) What the resource costs a month now, when the recommender says.
- `monthly_savings` (Number) What doing it saves a month, as the recommender estimates it, in the recommendation's currency.
- `provider_value` (String) Which cloud the account is in: 'aws', 'azure' or 'gcp'.
- `recommendation_id` (String) The recommender's own id for the recommendation, kept from one read to the next.
- `recommended_configuration` (String) What the recommender would make it.
- `region` (String) The Region the resource is in; empty when none is given.
- `resource_id` (String) The resource the recommendation is about, as its cloud names it: an ARN, an ARM id, a full resource name.
- `resource_name` (String) What the resource is called, when the recommender says.
- `resource_type` (String) What kind of resource it is, in the recommender's words.
- `source` (String) Which recommender made it: AWS Cost Optimization Hub, Azure Advisor, Google Cloud Recommender or OneUptime.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `evidence` (String) For a right-sizing OneUptime worked out from its own readings: the CPU and memory readings and the days it was judged from. A JSON value: read it with `jsondecode()`.
- `last_seen_at` (String) When the account's connection last read the recommendation.
- `project_id` (String) ID of the project this recommendation belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
