---
page_title: "oneuptime_data_source Data Source - oneuptime"
subcategory: "Other"
description: |-
  Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data.
---

# oneuptime_data_source (Data Source)

Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data source may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_source" "example" {
  name = "Example data source"
}

# Or by id:
data "oneuptime_data_source" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_source_type` (String) The kind of external system this data source connects to (Prometheus, PostgreSQL, MySQL, Microsoft SQL Server, ClickHouse, Loki, Elasticsearch, or REST API).
- `database_host` (String) Hostname or IP address for database sources (PostgreSQL, MySQL, SQL Server, ClickHouse).
- `database_name` (String) Database (or ClickHouse database) to connect to.
- `database_port` (Number) Port for database sources. Defaults per engine: PostgreSQL 5432, MySQL 3306, SQL Server 1433, ClickHouse 8123.
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.
- `url` (String) Base URL for HTTP-based sources (Prometheus, Loki, Elasticsearch, REST API) — e.g. https://prometheus.example.com.
- `username` (String) Username for database sources, or HTTP basic-auth username for HTTP sources. Use a READ-ONLY account — dashboards only ever read.

### Read-Only

- `additional_options` (String) Per-type options that are not secrets — e.g. { "sslEnabled": true, "elasticsearchIndex": "logs-*" }. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
