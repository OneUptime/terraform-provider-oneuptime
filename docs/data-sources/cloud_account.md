---
page_title: "oneuptime_cloud_account Data Source - oneuptime"
subcategory: "Other"
description: |-
  AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them.
---

# oneuptime_cloud_account (Data Source)

AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one cloud account may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_cloud_account" "example" {
  name = "Example cloud account"
}

# Or by id:
data "oneuptime_cloud_account" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `account_id` (String) The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.
- `cloud_connection_id` (String) ID of the connection that reads this account, if one does. The ID of a `oneuptime_cloud_connection`.
- `created_by_user_id` (String) ID of the user who created this account. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the account is for, in your team's words.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What the account is called in OneUptime. Starts as the name its cloud gives it, and can be changed.
- `provider_account_name` (String) What the cloud calls the account: an AWS account alias, an Azure subscription's or a Google Cloud project's name.
- `provider_value` (String) Which cloud the account is in: 'aws', 'azure' or 'gcp'.
- `region_count` (Number) How many regions the account's live resources are in.
- `resource_count` (Number) How many live resources the account has in OneUptime.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `inventory_listed_at` (String) When the account's inventory last listed every resource it has.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When a cloud connection last reached the account.
- `project_id` (String) ID of the project this account belongs to. The ID of a `oneuptime_project`.
- `resource_summary` (String) The account's live resources counted by type and by region (CloudAccountSummary). A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
