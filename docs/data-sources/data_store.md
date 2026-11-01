---
page_title: "oneuptime_data_store Data Source - oneuptime"
subcategory: "Other"
description: |-
  Every bucket, database, warehouse and telemetry stream Data Security knows: what sensitive data it holds, who can reach it and how it is configured.
---

# oneuptime_data_store (Data Source)

Every bucket, database, warehouse and telemetry stream Data Security knows: what sensitive data it holds, who can reach it and how it is configured.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data store may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_store" "example" {
  name = "Example data store"
}

# Or by id:
data "oneuptime_data_store" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `cloud_account_id` (String) The cloud account the store belongs to, for a cloud store. The ID of a `oneuptime_cloud_account`.
- `cloud_account_number` (String) The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.
- `cloud_connection_id` (String) The cloud connection that reads the store's configuration, if any. The ID of a `oneuptime_cloud_connection`.
- `cloud_resource_id` (String) The cloud resource the store was discovered from, if any. The ID of a `oneuptime_cloud_resource`.
- `content_scan_status` (String) How the last content scan ended, or why it could not run.
- `created_by_user_id` (String) ID of the user who created this data store. The ID of a `oneuptime_user` (see the data source).
- `database_server_id` (String) The Database product server the store is, for a self-managed database. The ID of a `oneuptime_database`.
- `description` (String) What the store is for, in your team's words.
- `discovered_from` (String) How OneUptime found the store: cloud inventory, the Database product, telemetry, the service map, or added by hand.
- `environment` (String) Which environment the store serves: Production, Staging, Development, Test or Unknown. Read from its tags and name; people may set it.
- `estimated_records` (Number) An estimate of how many records with sensitive data the store holds, from its samples.
- `exposure` (String) Where the store sits relative to the internet and other organizations: Public, Other accounts, Internet reachable, Private or Unknown.
- `family` (String) The store's family: object storage, database, warehouse, search, cache, file system, snapshot or telemetry.
- `finding_count` (Number) How many data types were found at how many locations.
- `highest_issue_severity` (String) The severity of the store's worst open issue.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_data` (Boolean) The store holds data used to train or ground an AI model.
- `is_archived` (Boolean) The store is gone from its cloud or source; it is kept with its history.
- `is_content_scan_enabled` (Boolean) Whether Data Security samples the store's content. Configuration checks run either way.
- `is_public_by_design` (Boolean) The store is meant to be public, such as a website's assets. Public access checks pass for it.
- `kind` (String) What kind of store it is: an S3 bucket, an RDS instance, a PostgreSQL database, a service's logs.
- `name` (String) The store's name: a bucket, database, dataset or service name.
- `object_count` (Number) How many objects, rows or documents the store holds, when known.
- `open_issue_count` (Number) How many of the store's issues are open.
- `probe_connection_variable` (String) For a database scanned from a probe: the environment variable on that probe that holds its connection string. The connection string itself never reaches OneUptime.
- `provider_resource_id` (String) The store's id in its cloud: an ARN, an Azure resource id or a Google Cloud resource name.
- `provider_value` (String) Where the store lives: aws, azure, gcp, self-managed or oneuptime.
- `purpose` (String) Why the store holds personal data, for the GDPR record of processing.
- `region` (String) The region the store lives in.
- `retention` (String) How long data is kept in the store, for the GDPR record of processing.
- `risk_score` (Number) The risk score of the store's worst open issue, from 0 to 100.
- `scan_probe_id` (String) A probe that scans the store's content from inside your network instead of OneUptime's workers. The ID of a `oneuptime_probe`.
- `sensitivity` (String) How sensitive the store is: Restricted, Confidential, Internal or Public, from the most sensitive data type found in it.
- `sensitivity_override` (String) A sensitivity people set by hand, which the store is treated as instead of what the scans found.
- `service_id` (String) The telemetry service whose logs, traces, exceptions, replays or AI conversations the store is. The ID of a `oneuptime_service`.
- `size_in_bytes` (Number) How much data the store holds, from the metrics already collected.
- `store_key` (String) What identifies the store within the project: its provider, kind and native id. One row per store, however often it is discovered.

### Read-Only

- `cloud_tags` (String) The store's tags or labels in its cloud. A JSON value: write it with `jsonencode()`.
- `configuration_read_at` (String) When the store's configuration and access were last read.
- `content_scanned_at` (String) When the store's content was last sampled.
- `created_at` (String) Date and Time when the object was created.
- `data_categories` (String) The categories of those data types. A JSON value: write it with `jsonencode()`.
- `data_type_ids` (String) The data types found in the store with at least Medium confidence. A JSON value: write it with `jsonencode()`.
- `first_seen_at` (String) When OneUptime first found the store.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_activity_at` (String) When the store was last read or written, when its metrics or traces say.
- `last_seen_at` (String) When OneUptime last found the store.
- `next_content_scan_at` (String) When the store's content is next sampled.
- `posture` (String) What the store's configuration says about its safety: public access, encryption, backups, logging and the rest. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this data store belongs to. The ID of a `oneuptime_project`.
- `scan_exclusions` (String) Path prefixes, tables or attribute keys the content scan skips. A JSON value: write it with `jsonencode()`.
- `subject_countries` (String) Countries whose people's data the store holds, from country-bound data types (ISO 3166-1 alpha-2). A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
