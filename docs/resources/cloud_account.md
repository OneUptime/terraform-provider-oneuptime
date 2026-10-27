---
page_title: "oneuptime_cloud_account Resource - oneuptime"
subcategory: "Other"
description: |-
  AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them.
---

# oneuptime_cloud_account (Resource)

AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them.

## Example Usage

```terraform
resource "oneuptime_cloud_account" "example" {
  name           = "Example cloud account"
  provider_value = "Example short text"
  account_id     = "Example short text"
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `account_id` (String) The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.
- `name` (String) What the account is called in OneUptime. Starts as the name its cloud gives it, and can be changed.
- `provider_value` (String) Which cloud the account is in: 'aws', 'azure' or 'gcp'.

### Optional

- `description` (String) What the account is for, in your team's words.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.

### Read-Only

- `cloud_connection_id` (String) ID of the connection that reads this account, if one does. The ID of a `oneuptime_cloud_connection`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this account. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `inventory_listed_at` (String) When the account's inventory last listed every resource it has.
- `last_seen_at` (String) When a cloud connection last reached the account.
- `project_id` (String) ID of the project this account belongs to. The ID of a `oneuptime_project`.
- `provider_account_name` (String) What the cloud calls the account: an AWS account alias, an Azure subscription's or a Google Cloud project's name.
- `region_count` (Number) How many regions the account's live resources are in.
- `resource_count` (Number) How many live resources the account has in OneUptime.
- `resource_summary` (String) The account's live resources counted by type and by region (CloudAccountSummary). A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing cloud account by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_account.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_account.example <id>
```
