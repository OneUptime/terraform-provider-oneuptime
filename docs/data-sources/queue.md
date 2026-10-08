---
page_title: "oneuptime_queue Data Source - oneuptime"
subcategory: "Other"
description: |-
  Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually.
---

# oneuptime_queue (Data Source)

Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one queue may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_queue" "example" {
  name = "Example queue"
}

# Or by id:
data "oneuptime_queue" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `broker_address` (String) The broker address (host[:port]) last reported for this queue by its telemetry, with any credentials removed. Display only - it is not part of the queue's identity.
- `broker_scope` (String) The Azure Service Bus or Event Hubs namespace this queue lives in, lowercased (the first part of <namespace>.servicebus.windows.net). Part of the queue's identity, because two namespaces can each hold a queue of the same name. Empty for every other messaging system.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description for this queue.
- `destination_name` (String) The queue, topic or subscription name as applications and the broker name it (the OpenTelemetry messaging.destination.name value), normalized the way discovery normalizes it: an SQS queue URL or an SNS topic ARN becomes its name, a Pub/Sub resource path its id, a Pulsar short name its persistent://public/default/ topic, and each UUID {uuid}.
- `discovery_source` (String) How this queue was first discovered: traces (messaging spans of instrumented applications), broker-metrics (the broker's own metrics) or manual.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Is this queue archived? Archived queues are hidden from lists but their telemetry is still collected.
- `messaging_system` (String) The broker this queue lives on, as an OpenTelemetry messaging.system value, e.g. kafka, rabbitmq, activemq, jms, aws_sqs, aws.sns, gcp_pubsub, servicebus, eventhubs, pulsar, rocketmq, nats or bullmq. A queue first seen through JMS is refined to its broker (activemq) once that broker's metrics are seen.
- `name` (String) Name of this queue. Discovered queues are named after their destination, e.g. orders.created. Not unique - two brokers may each have a queue of the same name. Rename freely.
- `queue_identifier` (String) Stable identity of this queue within the project: the messaging system's identity family, the broker scope (the Azure Service Bus / Event Hubs namespace, empty for every other system) and the destination, canonicalized and joined with '|' (e.g. kafka||orders.created, servicebus|orders-prod|orders). An ActiveMQ queue keys as jms, like the JMS clients that use it, so both are one queue. Computed by the server for manually added queues.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When was this queue archived?
- `auto_archived_at` (String) When this queue was archived automatically because no discovery source had seen it for a while. Empty when it was archived by a person or is not archived.
- `automatic_assignments` (String) Label and owner ids that label rules or owner rules attached automatically. Maintained by OneUptime. A JSON value: write it with `jsonencode()`.
- `broker_metrics_last_seen_at` (String) When the broker's own metrics for this queue (backlog, lag, dead letters, ...) were last received. Empty when they never were.
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When this queue was last seen by any discovery source - application messaging spans or broker metrics.
- `manually_restored_at` (String) When a person last restored this queue from the archive. Automatic archiving leaves it alone until it is seen again or a grace period passes.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
