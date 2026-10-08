---
page_title: "oneuptime_kubernetes_cluster Data Source - oneuptime"
subcategory: "Other"
description: |-
  Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered.
---

# oneuptime_kubernetes_cluster (Data Source)

Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one kubernetes cluster may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_kubernetes_cluster" "example" {
  name = "Example kubernetes cluster"
}

# Or by id:
data "oneuptime_kubernetes_cluster" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the OneUptime Kubernetes agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_access_credential_id` (String) ID of the Kubernetes credential the AI access Runner uses for this cluster; the credential must be assigned to that Runner. Empty for the in-cluster Runner, which is never given a credential. Binding a credential needs Project Owner, Project Admin or Edit Auto Remediation Rule, and also permission to read credentials (Project Owner, Project Admin or Read Runbook Credential); anyone who may edit the cluster can clear it. The ID of a `oneuptime_runbook_credential`.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running kubectl on this cluster, kept until the next successful command. Set by the server.
- `ai_access_runner_id` (String) ID of the Runner OneUptime AI uses to run kubectl against this cluster. Another cluster's in-cluster Runner cannot be bound. Binding a Runner, or switching to a different one, needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the cluster can clear it. The ID of a `oneuptime_runner`.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this cluster. RequireApproval: AI composes a kubectl plan and a human approves it with one click before anything runs. Any follow-up plan asks again. Automatic: AI runs safe changes without a human — each on ONE named object: rollout restart/undo/pause/resume of one workload, scale one workload above zero, delete one named pod, cordon/uncordon one node, label/annotate one pod or workload with unreserved keys. A riskier change (patch, set image, drain, taint, scale to zero, deleting workloads or jobs, anything touching several objects) never runs without one: when the round could only find riskier fixes it ends by proposing exactly those for one-click approval; when it also ran safe fixes, a riskier fix is proposed only if verification shows the safe ones did not recover the signal (the follow-up round, which asks). Shapes on the cluster's kubectl allowlist run on their own. BypassApproval: AI does not ask. Every change the policy allows — safe AND riskier — runs on its own, follow-up rounds included, except for what always asks (below). In EVERY mode, Bypass approval included: destructive commands (Denied tier) never run; a write in a protected namespace (kube-system, kube-public, kube-node-lease), a node drain, a node taint and a patch of a Node always need a human; the in-cluster Runner never changes its own namespace or anything outside the namespaces its chart may write; and an unattended run becomes a proposal when the hourly per-cluster circuit breaker trips or another unattended round already holds the cluster. Anyone who may edit the cluster can lower the mode (to Disabled, RequireApproval, or from BypassApproval to Automatic); raising it to Automatic or BypassApproval needs Project Owner, Project Admin or Edit Auto Remediation Rule.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `cluster_identifier` (String) Unique identifier for this cluster, sourced from the k8s.cluster.name OTel resource attribute.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description for this Kubernetes cluster.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only kubectl commands (get, describe, logs, events, top, rollout status) on this cluster, through the cluster's Kubernetes AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the cluster can turn it on or off.
- `is_archived` (Boolean) Is this Kubernetes cluster archived? Archived Kubernetes clusters are hidden from lists but keep collecting telemetry.
- `name` (String) Friendly name for this Kubernetes cluster.
- `namespace_count` (Number) Cached count of namespaces in this cluster.
- `node_count` (Number) Cached count of nodes in this cluster.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `pod_count` (Number) Cached count of pods in this cluster.
- `provider_value` (String) Cloud provider or platform running this cluster (EKS, GKE, AKS, self-managed, unknown).
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Kubernetes cluster. Leave blank to use the project-wide default.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this cluster was first configured: by the in-cluster Runner's first registration, or by anyone saving an AI access setting. Set by the server; never cleared, so an in-cluster Runner that registers later never overwrites a setting an operator chose.
- `ai_access_last_verified_at` (String) When a kubectl command from OneUptime AI last succeeded on this cluster. Set by the server.
- `ai_access_runner_bound_at` (String) When a Runner was first bound to this cluster for OneUptime AI access. Set by the server; never cleared, so a cluster whose Runner was cleared or deleted is not silently re-bound when the in-cluster Runner registers again.
- `ai_kubectl_command_allowlist` (String) Optional JSON array of kubectl command patterns that Automatic mode may run without approval even though they are riskier changes, for example: ["kubectl set image deployment/web * -n web"]. A pattern is compared with the command word by word: * stands for exactly one word (an image, a name), never for extra objects, flags or a second -n, every flag the command uses must be written out in the pattern, and the leading "kubectl" is optional. At most 100 patterns of at most 500 characters each; a pattern that is not one kubectl command line is refused. Destructive commands (Denied tier) never run regardless, and a write in a protected namespace (kube-system, kube-public, kube-node-lease) or a node drain still needs a human. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the cluster can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `archived_at` (String) When was this Kubernetes cluster archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this cluster.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Kubernetes cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the cluster default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
