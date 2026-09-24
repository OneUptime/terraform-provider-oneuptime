---
page_title: "oneuptime_kubernetes_cluster Data Source - oneuptime"
subcategory: "Other"
description: |-
  Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered.
---

# oneuptime_kubernetes_cluster (Data Source)

Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_kubernetes_cluster" "by_name" {
  name = "example-kubernetes_cluster"
}

data "oneuptime_kubernetes_cluster" "by_id" {
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
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description for this Kubernetes cluster.. Computed.
- `cluster_identifier` (String) Unique identifier for this cluster, sourced from the k8s.cluster.name OTel resource attribute.. Computed.
- `provider` (String) Cloud provider or platform running this cluster (EKS, GKE, AKS, self-managed, unknown).. Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).. Computed.
- `agent_version` (String) Version of the OneUptime Kubernetes agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `node_count` (Number) Cached count of nodes in this cluster.. Computed.
- `pod_count` (Number) Cached count of pods in this cluster.. Computed.
- `namespace_count` (Number) Cached count of namespaces in this cluster.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this Kubernetes cluster archived? Archived Kubernetes clusters are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Kubernetes cluster. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Kubernetes cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the cluster default, then the project's retention settings... Computed.
- `ai_access_runner_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `ai_access_credential_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_ai_investigation_enabled` (Bool) When on, OneUptime AI runs read-only kubectl commands (get, describe, logs, events, top, rollout status) on this cluster while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. Anyone who may edit the cluster can turn it on or off... Computed.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this cluster. RequireApproval: AI composes a kubectl plan and a human approves it with one click before anything runs. Any follow-up plan asks again. Automatic: AI runs safe changes without a human — each on ONE named object: rollout restart/undo/pause/resume of one workload, scale one workload above zero, delete one named pod, cordon/uncordon one node, label/annotate one pod or workload with unreserved keys. A riskier change (patch, set image, drain, taint, scale to zero, deleting workloads or jobs, anything touching several objects) never runs without one: when the round could only find riskier fixes it ends by proposing exactly those for one-click approval; when it also ran safe fixes, a riskier fix is proposed only if verification shows the safe ones did not recover the signal (the follow-up round, which asks). Shapes on the cluster's kubectl allowlist run on their own. BypassApproval: AI does not ask. Every change the policy allows — safe AND riskier — runs on its own, follow-up rounds included, except for what always asks (below). In EVERY mode, Bypass approval included: destructive commands (Denied tier) never run; a write in a protected namespace (kube-system, kube-public, kube-node-lease), a node drain, a node taint and a patch of a Node always need a human; the in-cluster Runner never changes its own namespace or anything outside the namespaces its chart may write; and an unattended run becomes a proposal when the hourly per-cluster circuit breaker trips or another unattended round already holds the cluster. Anyone who may edit the cluster can lower the mode (to Disabled, RequireApproval, or from BypassApproval to Automatic); raising it to Automatic or BypassApproval needs Project Owner, Project Admin or Edit Auto Remediation Rule... Computed.
- `ai_kubectl_command_allowlist` (String) Optional JSON array of kubectl command patterns that Automatic mode may run without approval even though they are riskier changes, for example: ["kubectl set image deployment/web * -n web"]. A pattern is compared with the command word by word: * stands for exactly one word (an image, a name), never for extra objects, flags or a second -n, every flag the command uses must be written out in the pattern, and the leading "kubectl" is optional. At most 100 patterns of at most 500 characters each; a pattern that is not one kubectl command line is refused. Destructive commands (Denied tier) never run regardless, and a write in a protected namespace (kube-system, kube-public, kube-node-lease) or a node drain still needs a human. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the cluster can remove patterns or clear the list... Computed.
- `ai_access_last_verified_at` (String) A date time object.. Computed.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running kubectl on this cluster, kept until the next successful command. Set by the server... Computed.
- `ai_access_configured_at` (String) A date time object.. Computed.
- `ai_access_runner_bound_at` (String) A date time object.. Computed.
