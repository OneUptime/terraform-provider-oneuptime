package shard09

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AiRuleResource{}
var _ resource.ResourceWithImportState = &AiRuleResource{}
var _ resource.ResourceWithModifyPlan = &AiRuleResource{}
var _ resource.ResourceWithMoveState = &AiRuleResource{}

func NewAiRuleResource() resource.Resource {
    return &AiRuleResource{}
}

// NewAiRuleLegacyResource registers this resource under the name it
// had before resource type names kept mixed-case words whole,
// oneuptime_auto_remediation_rule, so configurations that use it keep
// working. Deprecated: plans that use it say so.
func NewAiRuleLegacyResource() resource.Resource {
    return &AiRuleResource{isLegacyAlias: true}
}

// AiRuleResource defines the resource implementation.
type AiRuleResource struct {
    client *Client
    // Registered under the name this resource had before (see NewAiRuleLegacyResource).
    isLegacyAlias bool
}

// AiRuleResourceModel describes the resource data model.
type AiRuleResourceModel struct {
    Id types.String `tfsdk:"id"`
    Criteria JSONSubsetValue `tfsdk:"criteria"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    TriggerEntityType types.String `tfsdk:"trigger_entity_type"`
    Investigates types.Bool `tfsdk:"investigates"`
    OnlyIfOpenAfterMinutes types.Number `tfsdk:"only_if_open_after_minutes"`
    Fixes types.Bool `tfsdk:"fixes"`
    ExecutionMode types.String `tfsdk:"execution_mode"`
    RemediationAction types.String `tfsdk:"remediation_action"`
    AiSelectsRunbook types.Bool `tfsdk:"ai_selects_runbook"`
    AiComposesCommands types.Bool `tfsdk:"ai_composes_commands"`
    CommandAllowlist JSONSubsetValue `tfsdk:"command_allowlist"`
    CommandRunners types.Set `tfsdk:"command_runners"`
    Monitors types.Set `tfsdk:"monitors"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    AlertSeverities types.Set `tfsdk:"alert_severities"`
    Labels types.Set `tfsdk:"labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    TitlePattern types.String `tfsdk:"title_pattern"`
    DescriptionPattern types.String `tfsdk:"description_pattern"`
    Services types.Set `tfsdk:"services"`
    ServiceLabels types.Set `tfsdk:"service_labels"`
    ServiceNamePattern types.String `tfsdk:"service_name_pattern"`
    DataSecurityAnalyze types.String `tfsdk:"data_security_analyze"`
    DataSecurityContain types.String `tfsdk:"data_security_contain"`
    DataSecurityFix types.String `tfsdk:"data_security_fix"`
    DataSecurityJudge types.String `tfsdk:"data_security_judge"`
    DataSecurityChecks JSONSubsetValue `tfsdk:"data_security_checks"`
    DataSecurityCheckGroups JSONSubsetValue `tfsdk:"data_security_check_groups"`
    DataSecuritySeverities JSONSubsetValue `tfsdk:"data_security_severities"`
    DataStoreSensitivities JSONSubsetValue `tfsdk:"data_store_sensitivities"`
    DataStoreKinds JSONSubsetValue `tfsdk:"data_store_kinds"`
    DataStoreProviders JSONSubsetValue `tfsdk:"data_store_providers"`
    DataStoreEnvironments JSONSubsetValue `tfsdk:"data_store_environments"`
    DataStoreLabels types.Set `tfsdk:"data_store_labels"`
    DataCategories JSONSubsetValue `tfsdk:"data_categories"`
    Runbooks types.Set `tfsdk:"runbooks"`
    ExternalMcpServers types.Set `tfsdk:"external_mcp_servers"`
    RunnerId types.String `tfsdk:"runner_id"`
    CodeRepositories types.Set `tfsdk:"code_repositories"`
    LlmProviderId types.String `tfsdk:"llm_provider_id"`
    AiContext types.String `tfsdk:"ai_context"`
    OpenCodeFixPullRequests types.Bool `tfsdk:"open_code_fix_pull_requests"`
    OpenTelemetryPullRequests types.Bool `tfsdk:"open_telemetry_pull_requests"`
    VerificationWindowMinutes types.Number `tfsdk:"verification_window_minutes"`
    AutoResolveOnVerifiedRecovery types.Bool `tfsdk:"auto_resolve_on_verified_recovery"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt RFC3339Value `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    EmptiedLimits JSONSubsetValue `tfsdk:"emptied_limits"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *AiRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    if r.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_auto_remediation_rule"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_ai_rule"
}

func (r *AiRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
    if r.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_auto_remediation_rule has been renamed to oneuptime_ai_rule. Rename the resource in your configuration and add a moved block (from = oneuptime_auto_remediation_rule.<name>, to = oneuptime_ai_rule.<name>) to keep the existing ai rule; the old name keeps working until then."
    }
}

func (r *AiRuleResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "Which new incidents or alerts OneUptime AI investigates and fixes, and how. A rule investigates what it matches (right away, or only if still open after a wait), and fixes it or not, asking first or not. A fix may run the project's runbooks (only the rule's, when it names some), Bash commands on a Runner the rule picks, the tools of the project's external MCP servers (only the rule's, when it names some), and open fix and telemetry pull requests in the project's code repositories (only the rule's, when it names some). With no rule, OneUptime AI investigates every one, and fixes every one while automatic fixing is on. These were auto remediation rules, and investigation rules are AI rules now too. A rule whose Trigger Entity Type is Profile, Trace, Exception, Log or Metric decides the same about the problems OneUptime finds in a service's profiles, traces, exceptions, logs or metrics, matched by their service; their fix is a pull request in one of the project's code repositories.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "criteria": schema.StringAttribute{
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this AI rule.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this AI rule.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "trigger_entity_type": schema.StringAttribute{
                MarkdownDescription: "Which kind of signal this rule decides about: Incident or Alert - or Profile, the slow code OneUptime finds in a service's profiles, or Trace, the slow and failing operations OneUptime finds in a service's traces, or Exception, the new and rising exceptions OneUptime finds in a service, or Log, the error, missing and noisy logs OneUptime finds in a service's logs, or Metric, the problems OneUptime finds in a service's metrics, or DataSecurityIssue, a Data Security check that failed on a store.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Incident", "Alert", "Profile", "Trace", "Exception", "Log", "Metric", "DataSecurityIssue"),
                },
            },
            "investigates": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI investigates the incidents or alerts this rule matches. With enabled rules that investigate, only their matches are investigated. Off, the rule only fixes, and does not change which incidents are investigated. The dashboard turns it on for every rule; left out of an API write, it is off.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "only_if_open_after_minutes": schema.NumberAttribute{
                MarkdownDescription: "Investigate only incidents or alerts still open this many minutes after they were created: OneUptime AI waits that long, then investigates if it is still open. Leave empty to investigate right away. When several rules match, the shortest wait decides. At most 1440 minutes (a day).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "fixes": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI fixes the incidents or alerts this rule matches, the way its Execution Mode says. Off, the rule only investigates them. With enabled rules, an incident is fixed only when a rule that fixes matches it. On by default.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "execution_mode": schema.StringAttribute{
                MarkdownDescription: "How a rule that fixes asks. Suggest asks before fixing: every fix the rule starts waits for one-click human approval. FullAuto fixes without asking: OneUptime AI fixes on its own where the cluster's or resource's AI agent page allows, a runbook it picks from the ones the rule names starts right away (one it picks from all of the project's runbooks still waits for approval), a change through an MCP server the rule names is made right away (through any of the project's MCP servers, it still waits for approval), and a rule that runs its runbooks every time starts them immediately.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Suggest"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Suggest", "FullAuto"),
                },
            },
            "remediation_action": schema.StringAttribute{
                MarkdownDescription: "OneUptimeAI (the default): OneUptime AI fixes the matched incident or alert - on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, with the project's runbooks (only the rule's Runbooks when it names some) and with commands on the rule's Runner when it picks one. Runbooks: the rule starts every one of its Runbooks for every match, without OneUptime AI deciding. Whether a person approves first is the rule's Execution Mode.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("OneUptimeAI"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("OneUptimeAI", "Runbooks"),
                },
            },
            "ai_selects_runbook": schema.BoolAttribute{
                MarkdownDescription: "Kept for rules saved before rules were simplified: an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached), always proposing it for approval, and the rule does nothing else. Every OneUptime AI rule now lets AI pick a runbook on its own, so new rules leave this off.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_composes_commands": schema.BoolAttribute{
                MarkdownDescription: "Kept for rules saved before rules were simplified: the AI investigates the incident/alert and composes Bash commands for the rule's Runner and Command Runners, and the rule does nothing else. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. A OneUptime AI rule now composes commands on its Runner when it picks one, so new rules leave this off.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "command_runners": schema.SetAttribute{
                MarkdownDescription: "Kept for rules saved before rules were simplified: more Runners the AI may target with composed commands, besides the rule's Runner. A Runner runs a rule's commands only when the rule names it. IDs of `oneuptime_runner` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_runner"),
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_monitor"),
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_incident_severity"),
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_alert_severity"),
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_label"),
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_label"),
            },
            "title_pattern": schema.StringAttribute{
                MarkdownDescription: "Case-insensitive regex matched against the entity's title. Leave empty to match any title.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "description_pattern": schema.StringAttribute{
                MarkdownDescription: "Case-insensitive regex matched against the entity's description. Leave empty to match any description.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "Profile, Trace, Exception, Log, Metric and Data Security AI rules only: match what OneUptime found in one of these services - for a Data Security issue, on a store that is the service's telemetry. Leave empty to match any service. IDs of `oneuptime_service` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_service"),
            },
            "service_labels": schema.SetAttribute{
                MarkdownDescription: "Profile, Trace, Exception, Log, Metric and Data Security AI rules only: match what OneUptime found in a service that carries at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_label"),
            },
            "service_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Profile, Trace, Exception, Log, Metric and Data Security AI rules only: case-insensitive regex matched against the name of the service OneUptime found the problem in. Leave empty to match any name.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "data_security_analyze": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: whether OneUptime AI analyzes the issues this rule matches on its own - what is wrong, what it puts at risk and how to fix it, written on the issue's timeline. Off or WithoutAsking: an analysis changes nothing. Left out of a create, WithoutAsking.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Off", "WithoutAsking"),
                },
            },
            "data_security_contain": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: how OneUptime AI contains the issues this rule matches - a scrub rule that hides a telemetry issue's data from now on. Off; Suggest, proposed for one-click approval; or WithoutAsking, made on its own and undone with one press. A change of who can reach data runs without asking only where a rule names it, and deleting data never does. Left out of a create, Suggest.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Off", "Suggest", "WithoutAsking"),
                },
            },
            "data_security_fix": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: whether OneUptime AI fixes the cause of the issues this rule matches with a pull request in one of the project's code repositories. Off; Suggest, offered for a person to open; or WithoutAsking, opened on its own for review. Nothing is ever merged. Left out of a create, Suggest.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Off", "Suggest", "WithoutAsking"),
                },
            },
            "data_security_judge": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: whether OneUptime AI judges what the issues this rule matches rest on - which findings are real, who owns the store. Off; Suggest, proposed for a person to confirm; or WithoutAsking, recorded on its own and set back with one press. Left out of a create, Suggest.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Off", "Suggest", "WithoutAsking"),
                },
            },
            "data_security_checks": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues of these checks, by their ids (public-read, secrets-in-telemetry). Leave empty to match any check. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_security_check_groups": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues of checks in these groups: Exposure, Encryption, Recovery, Visibility, Access, Data hygiene, Telemetry, Residency or AI. Leave empty to match any group. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_security_severities": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues of these severities: Critical, High, Medium or Low. Leave empty to match any severity. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_store_sensitivities": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues on stores of these sensitivities: Restricted, Confidential, Internal or Public. Leave empty to match any sensitivity. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_store_kinds": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues on stores of these kinds (AWS S3 Bucket, Logs). Leave empty to match any kind. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_store_providers": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues on stores of these providers: aws, azure, gcp, self-managed or oneuptime. Leave empty to match any provider. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_store_environments": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues on stores that serve these environments: Production, Staging, Development, Test or Unknown. Leave empty to match any environment. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "data_store_labels": schema.SetAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues on stores that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_label"),
            },
            "data_categories": schema.StringAttribute{
                MarkdownDescription: "Data Security AI rules only: match issues whose data types are of at least one of these categories: Contact and identity, Government ID, Financial, Health, Secrets and credentials, Special category or Business confidential. Leave empty to match any data. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "runbooks": schema.SetAttribute{
                MarkdownDescription: "The runbooks OneUptime AI may pick and run for this rule. Leave empty to let it pick any of the project's runbooks. A rule whose Fix With is Runbooks starts every one of them for every match instead. IDs of `oneuptime_runbook` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_runbook"),
            },
            "external_mcp_servers": schema.SetAttribute{
                MarkdownDescription: "The external MCP servers OneUptime AI may use while it investigates and fixes what this rule matches (an investigation uses only their tools that read). Leave empty to let it use all of the project's MCP servers; changes through them are then always proposed for approval first. IDs of `oneuptime_external_mcp_server` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_external_mcp_server"),
            },
            "runner_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Runner OneUptime AI runs this rule's Bash commands on: one of the project's own Runners. Empty means none: OneUptime AI runs no commands of its own for the rule. Commands on the affected Kubernetes cluster or host always go through its AI agent. The ID of a `oneuptime_runner`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDValidators("oneuptime_runner"),
            },
            "code_repositories": schema.SetAttribute{
                MarkdownDescription: "The code repositories OneUptime AI may read while it investigates and fixes what this rule matches, and open the rule's pull requests in. Leave empty to let it use every repository the project has connected. IDs of `oneuptime_code_repository` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDSetValidators("oneuptime_code_repository"),
            },
            "llm_provider_id": schema.StringAttribute{
                MarkdownDescription: "ID of the LLM provider OneUptime AI uses for everything this rule starts. Empty: the project's default provider, else the global one. The ID of a `oneuptime_llm_provider`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: RecordIDValidators("oneuptime_llm_provider"),
            },
            "ai_context": schema.StringAttribute{
                MarkdownDescription: "Notes OneUptime AI reads after the project's AI context when it investigates or fixes what this rule matches. Markdown, at most 5,000 characters.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "open_code_fix_pull_requests": schema.BoolAttribute{
                MarkdownDescription: "When an investigation of an incident or alert this rule matches finds the cause in your code, OneUptime AI opens a pull request with the fix, for review. Nothing is merged. Acts only while fixing new incidents (or alerts) automatically is on. On by default.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "open_telemetry_pull_requests": schema.BoolAttribute{
                MarkdownDescription: "When missing logs, traces or metrics stop an investigation of an incident or alert this rule matches, OneUptime AI opens a pull request that adds them, for review. Nothing is merged. Acts only while fixing new incidents (or alerts) automatically is on. On by default.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "verification_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_resolve_on_verified_recovery": schema.BoolAttribute{
                MarkdownDescription: "When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "emptied_limits": schema.StringAttribute{
                MarkdownDescription: "The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

// MoveState lets a moved block bring state over from the old name:
//
//     moved {
//       from = oneuptime_auto_remediation_rule.example
//       to   = oneuptime_ai_rule.example
//     }
func (r *AiRuleResource) MoveState(ctx context.Context) []resource.StateMover {
    if r.isLegacyAlias {
        return nil
    }

    return []resource.StateMover{
        LegacyNameStateMover("oneuptime_auto_remediation_rule", r.schemaDefinition()),
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *AiRuleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    // Under Crossplane, what a manifest leaves out is the record's (see crossplane.go).
    KeepLeftOutUnderCrossplane(ctx, req, resp)

    KeepStateWhenNothingChanges(ctx, req, resp)
}

func (r *AiRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *AiRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data AiRuleResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config AiRuleResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    aiRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := aiRuleRequest["data"].(map[string]interface{})

    if parsedCriteria := r.parseJSONField(data.Criteria); parsedCriteria != nil {
        requestDataMap["criteria"] = parsedCriteria
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.TriggerEntityType.IsNull() && !data.TriggerEntityType.IsUnknown() {
        requestDataMap["triggerEntityType"] = data.TriggerEntityType.ValueString()
    }
    if !data.Investigates.IsNull() && !data.Investigates.IsUnknown() {
        requestDataMap["investigates"] = data.Investigates.ValueBool()
    }
    if !data.OnlyIfOpenAfterMinutes.IsNull() && !data.OnlyIfOpenAfterMinutes.IsUnknown() {
        requestDataMap["onlyIfOpenAfterMinutes"] = r.bigFloatToFloat64(data.OnlyIfOpenAfterMinutes.ValueBigFloat())
    }
    if !data.Fixes.IsNull() && !data.Fixes.IsUnknown() {
        requestDataMap["fixes"] = data.Fixes.ValueBool()
    }
    if !data.ExecutionMode.IsNull() && !data.ExecutionMode.IsUnknown() {
        requestDataMap["executionMode"] = data.ExecutionMode.ValueString()
    }
    if !data.RemediationAction.IsNull() && !data.RemediationAction.IsUnknown() {
        requestDataMap["remediationAction"] = data.RemediationAction.ValueString()
    }
    if !data.AiSelectsRunbook.IsNull() && !data.AiSelectsRunbook.IsUnknown() {
        requestDataMap["aiSelectsRunbook"] = data.AiSelectsRunbook.ValueBool()
    }
    if !data.AiComposesCommands.IsNull() && !data.AiComposesCommands.IsUnknown() {
        requestDataMap["aiComposesCommands"] = data.AiComposesCommands.ValueBool()
    }
    if parsedCommandAllowlist := r.parseJSONField(data.CommandAllowlist); parsedCommandAllowlist != nil {
        requestDataMap["commandAllowlist"] = parsedCommandAllowlist
    }
    if !data.CommandRunners.IsNull() && !data.CommandRunners.IsUnknown() {
        requestDataMap["commandRunners"] = r.convertTerraformSetToInterface(data.CommandRunners)
    }
    if !data.Monitors.IsNull() && !data.Monitors.IsUnknown() {
        requestDataMap["monitors"] = r.convertTerraformSetToInterface(data.Monitors)
    }
    if !data.IncidentSeverities.IsNull() && !data.IncidentSeverities.IsUnknown() {
        requestDataMap["incidentSeverities"] = r.convertTerraformSetToInterface(data.IncidentSeverities)
    }
    if !data.AlertSeverities.IsNull() && !data.AlertSeverities.IsUnknown() {
        requestDataMap["alertSeverities"] = r.convertTerraformSetToInterface(data.AlertSeverities)
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.MonitorLabels.IsNull() && !data.MonitorLabels.IsUnknown() {
        requestDataMap["monitorLabels"] = r.convertTerraformSetToInterface(data.MonitorLabels)
    }
    if !data.TitlePattern.IsNull() && !data.TitlePattern.IsUnknown() {
        requestDataMap["titlePattern"] = data.TitlePattern.ValueString()
    }
    if !data.DescriptionPattern.IsNull() && !data.DescriptionPattern.IsUnknown() {
        requestDataMap["descriptionPattern"] = data.DescriptionPattern.ValueString()
    }
    if !data.Services.IsNull() && !data.Services.IsUnknown() {
        requestDataMap["services"] = r.convertTerraformSetToInterface(data.Services)
    }
    if !data.ServiceLabels.IsNull() && !data.ServiceLabels.IsUnknown() {
        requestDataMap["serviceLabels"] = r.convertTerraformSetToInterface(data.ServiceLabels)
    }
    if !data.ServiceNamePattern.IsNull() && !data.ServiceNamePattern.IsUnknown() {
        requestDataMap["serviceNamePattern"] = data.ServiceNamePattern.ValueString()
    }
    if !data.DataSecurityAnalyze.IsNull() && !data.DataSecurityAnalyze.IsUnknown() {
        requestDataMap["dataSecurityAnalyze"] = data.DataSecurityAnalyze.ValueString()
    }
    if !data.DataSecurityContain.IsNull() && !data.DataSecurityContain.IsUnknown() {
        requestDataMap["dataSecurityContain"] = data.DataSecurityContain.ValueString()
    }
    if !data.DataSecurityFix.IsNull() && !data.DataSecurityFix.IsUnknown() {
        requestDataMap["dataSecurityFix"] = data.DataSecurityFix.ValueString()
    }
    if !data.DataSecurityJudge.IsNull() && !data.DataSecurityJudge.IsUnknown() {
        requestDataMap["dataSecurityJudge"] = data.DataSecurityJudge.ValueString()
    }
    if parsedDataSecurityChecks := r.parseJSONField(data.DataSecurityChecks); parsedDataSecurityChecks != nil {
        requestDataMap["dataSecurityChecks"] = parsedDataSecurityChecks
    }
    if parsedDataSecurityCheckGroups := r.parseJSONField(data.DataSecurityCheckGroups); parsedDataSecurityCheckGroups != nil {
        requestDataMap["dataSecurityCheckGroups"] = parsedDataSecurityCheckGroups
    }
    if parsedDataSecuritySeverities := r.parseJSONField(data.DataSecuritySeverities); parsedDataSecuritySeverities != nil {
        requestDataMap["dataSecuritySeverities"] = parsedDataSecuritySeverities
    }
    if parsedDataStoreSensitivities := r.parseJSONField(data.DataStoreSensitivities); parsedDataStoreSensitivities != nil {
        requestDataMap["dataStoreSensitivities"] = parsedDataStoreSensitivities
    }
    if parsedDataStoreKinds := r.parseJSONField(data.DataStoreKinds); parsedDataStoreKinds != nil {
        requestDataMap["dataStoreKinds"] = parsedDataStoreKinds
    }
    if parsedDataStoreProviders := r.parseJSONField(data.DataStoreProviders); parsedDataStoreProviders != nil {
        requestDataMap["dataStoreProviders"] = parsedDataStoreProviders
    }
    if parsedDataStoreEnvironments := r.parseJSONField(data.DataStoreEnvironments); parsedDataStoreEnvironments != nil {
        requestDataMap["dataStoreEnvironments"] = parsedDataStoreEnvironments
    }
    if !data.DataStoreLabels.IsNull() && !data.DataStoreLabels.IsUnknown() {
        requestDataMap["dataStoreLabels"] = r.convertTerraformSetToInterface(data.DataStoreLabels)
    }
    if parsedDataCategories := r.parseJSONField(data.DataCategories); parsedDataCategories != nil {
        requestDataMap["dataCategories"] = parsedDataCategories
    }
    if !data.Runbooks.IsNull() && !data.Runbooks.IsUnknown() {
        requestDataMap["runbooks"] = r.convertTerraformSetToInterface(data.Runbooks)
    }
    if !data.ExternalMcpServers.IsNull() && !data.ExternalMcpServers.IsUnknown() {
        requestDataMap["externalMcpServers"] = r.convertTerraformSetToInterface(data.ExternalMcpServers)
    }
    if !data.RunnerId.IsNull() && !data.RunnerId.IsUnknown() {
        requestDataMap["runnerId"] = data.RunnerId.ValueString()
    }
    if !data.CodeRepositories.IsNull() && !data.CodeRepositories.IsUnknown() {
        requestDataMap["codeRepositories"] = r.convertTerraformSetToInterface(data.CodeRepositories)
    }
    if !data.LlmProviderId.IsNull() && !data.LlmProviderId.IsUnknown() {
        requestDataMap["llmProviderId"] = data.LlmProviderId.ValueString()
    }
    if !data.AiContext.IsNull() && !data.AiContext.IsUnknown() {
        requestDataMap["aiContext"] = data.AiContext.ValueString()
    }
    if !data.OpenCodeFixPullRequests.IsNull() && !data.OpenCodeFixPullRequests.IsUnknown() {
        requestDataMap["openCodeFixPullRequests"] = data.OpenCodeFixPullRequests.ValueBool()
    }
    if !data.OpenTelemetryPullRequests.IsNull() && !data.OpenTelemetryPullRequests.IsUnknown() {
        requestDataMap["openTelemetryPullRequests"] = data.OpenTelemetryPullRequests.ValueBool()
    }
    if !data.VerificationWindowMinutes.IsNull() && !data.VerificationWindowMinutes.IsUnknown() {
        requestDataMap["verificationWindowMinutes"] = r.bigFloatToFloat64(data.VerificationWindowMinutes.ValueBigFloat())
    }
    if !data.AutoResolveOnVerifiedRecovery.IsNull() && !data.AutoResolveOnVerifiedRecovery.IsUnknown() {
        requestDataMap["autoResolveOnVerifiedRecovery"] = data.AutoResolveOnVerifiedRecovery.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/auto-remediation-rule", aiRuleRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ai_rule, got error: %s", err))
        return
    }

    var aiRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &aiRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create ai_rule: %s", err))
        return
    }

    // What OneUptime left out because only it writes it, if anything (see entries.go).
    resp.Diagnostics.Append(FieldsLeftOutWarning(httpResp, "ai rule")...)

    // What the write does that someone should know, if anything (see entries.go).
    resp.Diagnostics.Append(WriteWarnings(httpResp, "ai rule")...)

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := aiRuleResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := aiRuleResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for ai_rule did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)


    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * ai_rule is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "triggerEntityType": true,
        "investigates": true,
        "onlyIfOpenAfterMinutes": true,
        "fixes": true,
        "executionMode": true,
        "remediationAction": true,
        "aiSelectsRunbook": true,
        "aiComposesCommands": true,
        "commandAllowlist": true,
        "commandRunners": true,
        "monitors": true,
        "incidentSeverities": true,
        "alertSeverities": true,
        "labels": true,
        "monitorLabels": true,
        "titlePattern": true,
        "descriptionPattern": true,
        "services": true,
        "serviceLabels": true,
        "serviceNamePattern": true,
        "dataSecurityAnalyze": true,
        "dataSecurityContain": true,
        "dataSecurityFix": true,
        "dataSecurityJudge": true,
        "dataSecurityChecks": true,
        "dataSecurityCheckGroups": true,
        "dataSecuritySeverities": true,
        "dataStoreSensitivities": true,
        "dataStoreKinds": true,
        "dataStoreProviders": true,
        "dataStoreEnvironments": true,
        "dataStoreLabels": true,
        "dataCategories": true,
        "runbooks": true,
        "externalMcpServers": true,
        "runnerId": true,
        "codeRepositories": true,
        "llmProviderId": true,
        "aiContext": true,
        "openCodeFixPullRequests": true,
        "openTelemetryPullRequests": true,
        "verificationWindowMinutes": true,
        "autoResolveOnVerifiedRecovery": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "emptiedLimits": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/auto-remediation-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created ai_rule but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created ai_rule but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["criteria"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Criteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Criteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Criteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Criteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Criteria = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["criteria"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Criteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Criteria = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["criteria"].(string); ok {
        data.Criteria = NewJSONSubsetValue(val)
    } else {
        data.Criteria = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["triggerEntityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TriggerEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TriggerEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if arr, ok := dataMap["triggerEntityType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if val, ok := dataMap["triggerEntityType"].(string); ok {
        data.TriggerEntityType = types.StringValue(val)
    } else {
        data.TriggerEntityType = types.StringNull()
    }
    if val, ok := dataMap["investigates"].(bool); ok {
        data.Investigates = types.BoolValue(val)
    }
    data.OnlyIfOpenAfterMinutes = NumberFromAPI(data.OnlyIfOpenAfterMinutes, dataMap["onlyIfOpenAfterMinutes"])
    if val, ok := dataMap["fixes"].(bool); ok {
        data.Fixes = types.BoolValue(val)
    }
    if obj, ok := dataMap["executionMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ExecutionMode = types.StringValue(string(jsonBytes))
            } else {
                data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ExecutionMode = types.StringValue(string(jsonBytes))
            } else {
                data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if arr, ok := dataMap["executionMode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if val, ok := dataMap["executionMode"].(string); ok {
        data.ExecutionMode = types.StringValue(val)
    } else {
        data.ExecutionMode = types.StringNull()
    }
    if obj, ok := dataMap["remediationAction"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RemediationAction = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RemediationAction = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationAction = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RemediationAction = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationAction = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if arr, ok := dataMap["remediationAction"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if val, ok := dataMap["remediationAction"].(string); ok {
        data.RemediationAction = types.StringValue(val)
    } else {
        data.RemediationAction = types.StringNull()
    }
    if val, ok := dataMap["aiSelectsRunbook"].(bool); ok {
        data.AiSelectsRunbook = types.BoolValue(val)
    }
    if val, ok := dataMap["aiComposesCommands"].(bool); ok {
        data.AiComposesCommands = types.BoolValue(val)
    }
    if obj, ok := dataMap["commandAllowlist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["commandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["commandAllowlist"].(string); ok {
        data.CommandAllowlist = NewJSONSubsetValue(val)
    } else {
        data.CommandAllowlist = NewJSONSubsetNull()
    }
    if val, ok := dataMap["commandRunners"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.CommandRunners = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.CommandRunners = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["incidentSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.IncidentSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["titlePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.TitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.TitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if arr, ok := dataMap["titlePattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if val, ok := dataMap["titlePattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    } else {
        data.TitlePattern = types.StringNull()
    }
    if obj, ok := dataMap["descriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if arr, ok := dataMap["descriptionPattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["descriptionPattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    } else {
        data.DescriptionPattern = types.StringNull()
    }
    if val, ok := dataMap["services"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Services = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Services = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["serviceLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.ServiceLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.ServiceLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["serviceNamePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if arr, ok := dataMap["serviceNamePattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if val, ok := dataMap["serviceNamePattern"].(string); ok {
        data.ServiceNamePattern = types.StringValue(val)
    } else {
        data.ServiceNamePattern = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityAnalyze"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityAnalyze = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityAnalyze = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityAnalyze = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityAnalyze"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityAnalyze = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityAnalyze"].(string); ok {
        data.DataSecurityAnalyze = types.StringValue(val)
    } else {
        data.DataSecurityAnalyze = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityContain"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityContain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityContain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityContain = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityContain = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityContain = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityContain = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityContain"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityContain = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityContain = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityContain"].(string); ok {
        data.DataSecurityContain = types.StringValue(val)
    } else {
        data.DataSecurityContain = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityFix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityFix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityFix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityFix = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityFix = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityFix = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityFix = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityFix"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityFix = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityFix = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityFix"].(string); ok {
        data.DataSecurityFix = types.StringValue(val)
    } else {
        data.DataSecurityFix = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityJudge"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityJudge = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityJudge = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityJudge = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityJudge = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityJudge = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityJudge = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityJudge"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityJudge = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityJudge = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityJudge"].(string); ok {
        data.DataSecurityJudge = types.StringValue(val)
    } else {
        data.DataSecurityJudge = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityChecks"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityChecks = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityChecks = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityChecks = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecurityChecks"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityChecks = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecurityChecks"].(string); ok {
        data.DataSecurityChecks = NewJSONSubsetValue(val)
    } else {
        data.DataSecurityChecks = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataSecurityCheckGroups"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityCheckGroups = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecurityCheckGroups"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityCheckGroups = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecurityCheckGroups"].(string); ok {
        data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
    } else {
        data.DataSecurityCheckGroups = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataSecuritySeverities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecuritySeverities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecuritySeverities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecuritySeverities = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecuritySeverities"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecuritySeverities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecuritySeverities"].(string); ok {
        data.DataSecuritySeverities = NewJSONSubsetValue(val)
    } else {
        data.DataSecuritySeverities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreSensitivities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreSensitivities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreSensitivities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreSensitivities = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreSensitivities"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreSensitivities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreSensitivities"].(string); ok {
        data.DataStoreSensitivities = NewJSONSubsetValue(val)
    } else {
        data.DataStoreSensitivities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreKinds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreKinds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreKinds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreKinds = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreKinds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreKinds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreKinds"].(string); ok {
        data.DataStoreKinds = NewJSONSubsetValue(val)
    } else {
        data.DataStoreKinds = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreProviders"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreProviders = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreProviders = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreProviders = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreProviders"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreProviders = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreProviders"].(string); ok {
        data.DataStoreProviders = NewJSONSubsetValue(val)
    } else {
        data.DataStoreProviders = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreEnvironments"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreEnvironments = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreEnvironments = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreEnvironments = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreEnvironments"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreEnvironments = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreEnvironments"].(string); ok {
        data.DataStoreEnvironments = NewJSONSubsetValue(val)
    } else {
        data.DataStoreEnvironments = NewJSONSubsetNull()
    }
    if val, ok := dataMap["dataStoreLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DataStoreLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DataStoreLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["dataCategories"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataCategories"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataCategories"].(string); ok {
        data.DataCategories = NewJSONSubsetValue(val)
    } else {
        data.DataCategories = NewJSONSubsetNull()
    }
    if val, ok := dataMap["runbooks"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Runbooks = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Runbooks = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["externalMcpServers"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.ExternalMcpServers = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.ExternalMcpServers = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["runnerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RunnerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RunnerId = types.StringValue(string(jsonBytes))
            } else {
                data.RunnerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RunnerId = types.StringValue(string(jsonBytes))
            } else {
                data.RunnerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if arr, ok := dataMap["runnerId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if val, ok := dataMap["runnerId"].(string); ok {
        data.RunnerId = types.StringValue(val)
    } else {
        data.RunnerId = types.StringNull()
    }
    if val, ok := dataMap["codeRepositories"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.CodeRepositories = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.CodeRepositories = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["llmProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LlmProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LlmProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if arr, ok := dataMap["llmProviderId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderId"].(string); ok {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := dataMap["aiContext"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiContext = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiContext = types.StringValue(string(jsonBytes))
            } else {
                data.AiContext = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiContext = types.StringValue(string(jsonBytes))
            } else {
                data.AiContext = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if arr, ok := dataMap["aiContext"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if val, ok := dataMap["aiContext"].(string); ok {
        data.AiContext = types.StringValue(val)
    } else {
        data.AiContext = types.StringNull()
    }
    if val, ok := dataMap["openCodeFixPullRequests"].(bool); ok {
        data.OpenCodeFixPullRequests = types.BoolValue(val)
    }
    if val, ok := dataMap["openTelemetryPullRequests"].(bool); ok {
        data.OpenTelemetryPullRequests = types.BoolValue(val)
    }
    data.VerificationWindowMinutes = NumberFromAPI(data.VerificationWindowMinutes, dataMap["verificationWindowMinutes"])
    if val, ok := dataMap["autoResolveOnVerifiedRecovery"].(bool); ok {
        data.AutoResolveOnVerifiedRecovery = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["createdByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["createdByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastModifiedAt = NewRFC3339Value(val)
        } else {
            data.LastModifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastModifiedAt"].(string); ok && val != "" {
        data.LastModifiedAt = NewRFC3339Value(val)
    } else {
        data.LastModifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastModifiedByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["emptiedLimits"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmptiedLimits = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EmptiedLimits = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EmptiedLimits = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["emptiedLimits"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EmptiedLimits = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["emptiedLimits"].(string); ok {
        data.EmptiedLimits = NewJSONSubsetValue(val)
    } else {
        data.EmptiedLimits = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")


    // Nothing a create made was imported (see importsettle.go).
    resp.Diagnostics.Append(ForgetImportedFromServer(ctx, resp.Private)...)
    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AiRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data AiRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "triggerEntityType": true,
        "investigates": true,
        "onlyIfOpenAfterMinutes": true,
        "fixes": true,
        "executionMode": true,
        "remediationAction": true,
        "aiSelectsRunbook": true,
        "aiComposesCommands": true,
        "commandAllowlist": true,
        "commandRunners": true,
        "monitors": true,
        "incidentSeverities": true,
        "alertSeverities": true,
        "labels": true,
        "monitorLabels": true,
        "titlePattern": true,
        "descriptionPattern": true,
        "services": true,
        "serviceLabels": true,
        "serviceNamePattern": true,
        "dataSecurityAnalyze": true,
        "dataSecurityContain": true,
        "dataSecurityFix": true,
        "dataSecurityJudge": true,
        "dataSecurityChecks": true,
        "dataSecurityCheckGroups": true,
        "dataSecuritySeverities": true,
        "dataStoreSensitivities": true,
        "dataStoreKinds": true,
        "dataStoreProviders": true,
        "dataStoreEnvironments": true,
        "dataStoreLabels": true,
        "dataCategories": true,
        "runbooks": true,
        "externalMcpServers": true,
        "runnerId": true,
        "codeRepositories": true,
        "llmProviderId": true,
        "aiContext": true,
        "openCodeFixPullRequests": true,
        "openTelemetryPullRequests": true,
        "verificationWindowMinutes": true,
        "autoResolveOnVerifiedRecovery": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "emptiedLimits": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/auto-remediation-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_rule, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var aiRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &aiRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse ai_rule response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := aiRuleResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = aiRuleResponse
    }

    if obj, ok := dataMap["criteria"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Criteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Criteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Criteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Criteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Criteria = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["criteria"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Criteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Criteria = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["criteria"].(string); ok {
        data.Criteria = NewJSONSubsetValue(val)
    } else {
        data.Criteria = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["triggerEntityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TriggerEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TriggerEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if arr, ok := dataMap["triggerEntityType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if val, ok := dataMap["triggerEntityType"].(string); ok {
        data.TriggerEntityType = types.StringValue(val)
    } else {
        data.TriggerEntityType = types.StringNull()
    }
    if val, ok := dataMap["investigates"].(bool); ok {
        data.Investigates = types.BoolValue(val)
    }
    data.OnlyIfOpenAfterMinutes = NumberFromAPI(data.OnlyIfOpenAfterMinutes, dataMap["onlyIfOpenAfterMinutes"])
    if val, ok := dataMap["fixes"].(bool); ok {
        data.Fixes = types.BoolValue(val)
    }
    if obj, ok := dataMap["executionMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ExecutionMode = types.StringValue(string(jsonBytes))
            } else {
                data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ExecutionMode = types.StringValue(string(jsonBytes))
            } else {
                data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if arr, ok := dataMap["executionMode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if val, ok := dataMap["executionMode"].(string); ok {
        data.ExecutionMode = types.StringValue(val)
    } else {
        data.ExecutionMode = types.StringNull()
    }
    if obj, ok := dataMap["remediationAction"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RemediationAction = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RemediationAction = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationAction = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RemediationAction = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationAction = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if arr, ok := dataMap["remediationAction"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if val, ok := dataMap["remediationAction"].(string); ok {
        data.RemediationAction = types.StringValue(val)
    } else {
        data.RemediationAction = types.StringNull()
    }
    if val, ok := dataMap["aiSelectsRunbook"].(bool); ok {
        data.AiSelectsRunbook = types.BoolValue(val)
    }
    if val, ok := dataMap["aiComposesCommands"].(bool); ok {
        data.AiComposesCommands = types.BoolValue(val)
    }
    if obj, ok := dataMap["commandAllowlist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["commandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["commandAllowlist"].(string); ok {
        data.CommandAllowlist = NewJSONSubsetValue(val)
    } else {
        data.CommandAllowlist = NewJSONSubsetNull()
    }
    if val, ok := dataMap["commandRunners"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.CommandRunners = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.CommandRunners = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["incidentSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.IncidentSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["titlePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.TitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.TitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if arr, ok := dataMap["titlePattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if val, ok := dataMap["titlePattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    } else {
        data.TitlePattern = types.StringNull()
    }
    if obj, ok := dataMap["descriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if arr, ok := dataMap["descriptionPattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["descriptionPattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    } else {
        data.DescriptionPattern = types.StringNull()
    }
    if val, ok := dataMap["services"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Services = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Services = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["serviceLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.ServiceLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.ServiceLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["serviceNamePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if arr, ok := dataMap["serviceNamePattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if val, ok := dataMap["serviceNamePattern"].(string); ok {
        data.ServiceNamePattern = types.StringValue(val)
    } else {
        data.ServiceNamePattern = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityAnalyze"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityAnalyze = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityAnalyze = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityAnalyze = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityAnalyze"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityAnalyze = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityAnalyze"].(string); ok {
        data.DataSecurityAnalyze = types.StringValue(val)
    } else {
        data.DataSecurityAnalyze = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityContain"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityContain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityContain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityContain = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityContain = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityContain = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityContain = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityContain"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityContain = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityContain = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityContain"].(string); ok {
        data.DataSecurityContain = types.StringValue(val)
    } else {
        data.DataSecurityContain = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityFix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityFix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityFix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityFix = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityFix = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityFix = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityFix = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityFix"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityFix = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityFix = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityFix"].(string); ok {
        data.DataSecurityFix = types.StringValue(val)
    } else {
        data.DataSecurityFix = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityJudge"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityJudge = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityJudge = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityJudge = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityJudge = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityJudge = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityJudge = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityJudge"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityJudge = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityJudge = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityJudge"].(string); ok {
        data.DataSecurityJudge = types.StringValue(val)
    } else {
        data.DataSecurityJudge = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityChecks"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityChecks = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityChecks = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityChecks = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecurityChecks"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityChecks = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecurityChecks"].(string); ok {
        data.DataSecurityChecks = NewJSONSubsetValue(val)
    } else {
        data.DataSecurityChecks = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataSecurityCheckGroups"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityCheckGroups = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecurityCheckGroups"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityCheckGroups = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecurityCheckGroups"].(string); ok {
        data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
    } else {
        data.DataSecurityCheckGroups = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataSecuritySeverities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecuritySeverities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecuritySeverities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecuritySeverities = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecuritySeverities"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecuritySeverities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecuritySeverities"].(string); ok {
        data.DataSecuritySeverities = NewJSONSubsetValue(val)
    } else {
        data.DataSecuritySeverities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreSensitivities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreSensitivities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreSensitivities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreSensitivities = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreSensitivities"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreSensitivities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreSensitivities"].(string); ok {
        data.DataStoreSensitivities = NewJSONSubsetValue(val)
    } else {
        data.DataStoreSensitivities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreKinds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreKinds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreKinds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreKinds = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreKinds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreKinds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreKinds"].(string); ok {
        data.DataStoreKinds = NewJSONSubsetValue(val)
    } else {
        data.DataStoreKinds = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreProviders"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreProviders = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreProviders = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreProviders = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreProviders"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreProviders = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreProviders"].(string); ok {
        data.DataStoreProviders = NewJSONSubsetValue(val)
    } else {
        data.DataStoreProviders = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreEnvironments"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreEnvironments = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreEnvironments = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreEnvironments = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreEnvironments"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreEnvironments = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreEnvironments"].(string); ok {
        data.DataStoreEnvironments = NewJSONSubsetValue(val)
    } else {
        data.DataStoreEnvironments = NewJSONSubsetNull()
    }
    if val, ok := dataMap["dataStoreLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DataStoreLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DataStoreLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["dataCategories"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataCategories"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataCategories"].(string); ok {
        data.DataCategories = NewJSONSubsetValue(val)
    } else {
        data.DataCategories = NewJSONSubsetNull()
    }
    if val, ok := dataMap["runbooks"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Runbooks = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Runbooks = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["externalMcpServers"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.ExternalMcpServers = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.ExternalMcpServers = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["runnerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RunnerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RunnerId = types.StringValue(string(jsonBytes))
            } else {
                data.RunnerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RunnerId = types.StringValue(string(jsonBytes))
            } else {
                data.RunnerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if arr, ok := dataMap["runnerId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if val, ok := dataMap["runnerId"].(string); ok {
        data.RunnerId = types.StringValue(val)
    } else {
        data.RunnerId = types.StringNull()
    }
    if val, ok := dataMap["codeRepositories"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.CodeRepositories = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.CodeRepositories = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["llmProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LlmProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LlmProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if arr, ok := dataMap["llmProviderId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderId"].(string); ok {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := dataMap["aiContext"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiContext = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiContext = types.StringValue(string(jsonBytes))
            } else {
                data.AiContext = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiContext = types.StringValue(string(jsonBytes))
            } else {
                data.AiContext = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if arr, ok := dataMap["aiContext"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if val, ok := dataMap["aiContext"].(string); ok {
        data.AiContext = types.StringValue(val)
    } else {
        data.AiContext = types.StringNull()
    }
    if val, ok := dataMap["openCodeFixPullRequests"].(bool); ok {
        data.OpenCodeFixPullRequests = types.BoolValue(val)
    }
    if val, ok := dataMap["openTelemetryPullRequests"].(bool); ok {
        data.OpenTelemetryPullRequests = types.BoolValue(val)
    }
    data.VerificationWindowMinutes = NumberFromAPI(data.VerificationWindowMinutes, dataMap["verificationWindowMinutes"])
    if val, ok := dataMap["autoResolveOnVerifiedRecovery"].(bool); ok {
        data.AutoResolveOnVerifiedRecovery = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["createdByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["createdByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastModifiedAt = NewRFC3339Value(val)
        } else {
            data.LastModifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastModifiedAt"].(string); ok && val != "" {
        data.LastModifiedAt = NewRFC3339Value(val)
    } else {
        data.LastModifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastModifiedByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["emptiedLimits"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmptiedLimits = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EmptiedLimits = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EmptiedLimits = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["emptiedLimits"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EmptiedLimits = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["emptiedLimits"].(string); ok {
        data.EmptiedLimits = NewJSONSubsetValue(val)
    } else {
        data.EmptiedLimits = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AiRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data AiRuleResourceModel
    var state AiRuleResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // The attributes this update writes from the configuration (see importsettle.go).
    writtenFromConfig := map[string]bool{
        "command_allowlist": !data.CommandAllowlist.Equal(state.CommandAllowlist),
        "criteria": !data.Criteria.Equal(state.Criteria),
        "data_categories": !data.DataCategories.Equal(state.DataCategories),
        "data_security_check_groups": !data.DataSecurityCheckGroups.Equal(state.DataSecurityCheckGroups),
        "data_security_checks": !data.DataSecurityChecks.Equal(state.DataSecurityChecks),
        "data_security_severities": !data.DataSecuritySeverities.Equal(state.DataSecuritySeverities),
        "data_store_environments": !data.DataStoreEnvironments.Equal(state.DataStoreEnvironments),
        "data_store_kinds": !data.DataStoreKinds.Equal(state.DataStoreKinds),
        "data_store_providers": !data.DataStoreProviders.Equal(state.DataStoreProviders),
        "data_store_sensitivities": !data.DataStoreSensitivities.Equal(state.DataStoreSensitivities),
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config AiRuleResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    aiRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := aiRuleRequest["data"].(map[string]interface{})

    if !data.Criteria.IsUnknown() && !state.Criteria.IsUnknown() && !data.Criteria.Equal(state.Criteria) {
        var criteriaData interface{}
        if err := json.Unmarshal([]byte(data.Criteria.ValueString()), &criteriaData); err == nil {
            requestDataMap["criteria"] = criteriaData
        } else {
            requestDataMap["criteria"] = data.Criteria.ValueString()
        }
    }
    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.Investigates.IsUnknown() && !state.Investigates.IsUnknown() && !data.Investigates.Equal(state.Investigates) {
        requestDataMap["investigates"] = data.Investigates.ValueBool()
    }
    if !data.OnlyIfOpenAfterMinutes.IsUnknown() && !state.OnlyIfOpenAfterMinutes.IsUnknown() && !data.OnlyIfOpenAfterMinutes.Equal(state.OnlyIfOpenAfterMinutes) {
        requestDataMap["onlyIfOpenAfterMinutes"] = r.bigFloatToFloat64(data.OnlyIfOpenAfterMinutes.ValueBigFloat())
    }
    if !data.Fixes.IsUnknown() && !state.Fixes.IsUnknown() && !data.Fixes.Equal(state.Fixes) {
        requestDataMap["fixes"] = data.Fixes.ValueBool()
    }
    if !data.ExecutionMode.IsUnknown() && !state.ExecutionMode.IsUnknown() && !data.ExecutionMode.Equal(state.ExecutionMode) {
        requestDataMap["executionMode"] = data.ExecutionMode.ValueString()
    }
    if !data.RemediationAction.IsUnknown() && !state.RemediationAction.IsUnknown() && !data.RemediationAction.Equal(state.RemediationAction) {
        requestDataMap["remediationAction"] = data.RemediationAction.ValueString()
    }
    if !data.AiSelectsRunbook.IsUnknown() && !state.AiSelectsRunbook.IsUnknown() && !data.AiSelectsRunbook.Equal(state.AiSelectsRunbook) {
        requestDataMap["aiSelectsRunbook"] = data.AiSelectsRunbook.ValueBool()
    }
    if !data.AiComposesCommands.IsUnknown() && !state.AiComposesCommands.IsUnknown() && !data.AiComposesCommands.Equal(state.AiComposesCommands) {
        requestDataMap["aiComposesCommands"] = data.AiComposesCommands.ValueBool()
    }
    if !data.CommandAllowlist.IsUnknown() && !state.CommandAllowlist.IsUnknown() && !data.CommandAllowlist.Equal(state.CommandAllowlist) {
        var commandallowlistData interface{}
        if err := json.Unmarshal([]byte(data.CommandAllowlist.ValueString()), &commandallowlistData); err == nil {
            requestDataMap["commandAllowlist"] = commandallowlistData
        } else {
            requestDataMap["commandAllowlist"] = data.CommandAllowlist.ValueString()
        }
    }
    if !data.CommandRunners.IsUnknown() && !state.CommandRunners.IsUnknown() && !data.CommandRunners.Equal(state.CommandRunners) {
        requestDataMap["commandRunners"] = r.convertTerraformSetToInterface(data.CommandRunners)
    }
    if !data.Monitors.IsUnknown() && !state.Monitors.IsUnknown() && !data.Monitors.Equal(state.Monitors) {
        requestDataMap["monitors"] = r.convertTerraformSetToInterface(data.Monitors)
    }
    if !data.IncidentSeverities.IsUnknown() && !state.IncidentSeverities.IsUnknown() && !data.IncidentSeverities.Equal(state.IncidentSeverities) {
        requestDataMap["incidentSeverities"] = r.convertTerraformSetToInterface(data.IncidentSeverities)
    }
    if !data.AlertSeverities.IsUnknown() && !state.AlertSeverities.IsUnknown() && !data.AlertSeverities.Equal(state.AlertSeverities) {
        requestDataMap["alertSeverities"] = r.convertTerraformSetToInterface(data.AlertSeverities)
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.MonitorLabels.IsUnknown() && !state.MonitorLabels.IsUnknown() && !data.MonitorLabels.Equal(state.MonitorLabels) {
        requestDataMap["monitorLabels"] = r.convertTerraformSetToInterface(data.MonitorLabels)
    }
    if !data.TitlePattern.IsUnknown() && !state.TitlePattern.IsUnknown() && !data.TitlePattern.Equal(state.TitlePattern) {
        requestDataMap["titlePattern"] = data.TitlePattern.ValueString()
    }
    if !data.DescriptionPattern.IsUnknown() && !state.DescriptionPattern.IsUnknown() && !data.DescriptionPattern.Equal(state.DescriptionPattern) {
        requestDataMap["descriptionPattern"] = data.DescriptionPattern.ValueString()
    }
    if !data.Services.IsUnknown() && !state.Services.IsUnknown() && !data.Services.Equal(state.Services) {
        requestDataMap["services"] = r.convertTerraformSetToInterface(data.Services)
    }
    if !data.ServiceLabels.IsUnknown() && !state.ServiceLabels.IsUnknown() && !data.ServiceLabels.Equal(state.ServiceLabels) {
        requestDataMap["serviceLabels"] = r.convertTerraformSetToInterface(data.ServiceLabels)
    }
    if !data.ServiceNamePattern.IsUnknown() && !state.ServiceNamePattern.IsUnknown() && !data.ServiceNamePattern.Equal(state.ServiceNamePattern) {
        requestDataMap["serviceNamePattern"] = data.ServiceNamePattern.ValueString()
    }
    if !data.DataSecurityAnalyze.IsUnknown() && !state.DataSecurityAnalyze.IsUnknown() && !data.DataSecurityAnalyze.Equal(state.DataSecurityAnalyze) {
        requestDataMap["dataSecurityAnalyze"] = data.DataSecurityAnalyze.ValueString()
    }
    if !data.DataSecurityContain.IsUnknown() && !state.DataSecurityContain.IsUnknown() && !data.DataSecurityContain.Equal(state.DataSecurityContain) {
        requestDataMap["dataSecurityContain"] = data.DataSecurityContain.ValueString()
    }
    if !data.DataSecurityFix.IsUnknown() && !state.DataSecurityFix.IsUnknown() && !data.DataSecurityFix.Equal(state.DataSecurityFix) {
        requestDataMap["dataSecurityFix"] = data.DataSecurityFix.ValueString()
    }
    if !data.DataSecurityJudge.IsUnknown() && !state.DataSecurityJudge.IsUnknown() && !data.DataSecurityJudge.Equal(state.DataSecurityJudge) {
        requestDataMap["dataSecurityJudge"] = data.DataSecurityJudge.ValueString()
    }
    if !data.DataSecurityChecks.IsUnknown() && !state.DataSecurityChecks.IsUnknown() && !data.DataSecurityChecks.Equal(state.DataSecurityChecks) {
        var datasecuritychecksData interface{}
        if err := json.Unmarshal([]byte(data.DataSecurityChecks.ValueString()), &datasecuritychecksData); err == nil {
            requestDataMap["dataSecurityChecks"] = datasecuritychecksData
        } else {
            requestDataMap["dataSecurityChecks"] = data.DataSecurityChecks.ValueString()
        }
    }
    if !data.DataSecurityCheckGroups.IsUnknown() && !state.DataSecurityCheckGroups.IsUnknown() && !data.DataSecurityCheckGroups.Equal(state.DataSecurityCheckGroups) {
        var datasecuritycheckgroupsData interface{}
        if err := json.Unmarshal([]byte(data.DataSecurityCheckGroups.ValueString()), &datasecuritycheckgroupsData); err == nil {
            requestDataMap["dataSecurityCheckGroups"] = datasecuritycheckgroupsData
        } else {
            requestDataMap["dataSecurityCheckGroups"] = data.DataSecurityCheckGroups.ValueString()
        }
    }
    if !data.DataSecuritySeverities.IsUnknown() && !state.DataSecuritySeverities.IsUnknown() && !data.DataSecuritySeverities.Equal(state.DataSecuritySeverities) {
        var datasecurityseveritiesData interface{}
        if err := json.Unmarshal([]byte(data.DataSecuritySeverities.ValueString()), &datasecurityseveritiesData); err == nil {
            requestDataMap["dataSecuritySeverities"] = datasecurityseveritiesData
        } else {
            requestDataMap["dataSecuritySeverities"] = data.DataSecuritySeverities.ValueString()
        }
    }
    if !data.DataStoreSensitivities.IsUnknown() && !state.DataStoreSensitivities.IsUnknown() && !data.DataStoreSensitivities.Equal(state.DataStoreSensitivities) {
        var datastoresensitivitiesData interface{}
        if err := json.Unmarshal([]byte(data.DataStoreSensitivities.ValueString()), &datastoresensitivitiesData); err == nil {
            requestDataMap["dataStoreSensitivities"] = datastoresensitivitiesData
        } else {
            requestDataMap["dataStoreSensitivities"] = data.DataStoreSensitivities.ValueString()
        }
    }
    if !data.DataStoreKinds.IsUnknown() && !state.DataStoreKinds.IsUnknown() && !data.DataStoreKinds.Equal(state.DataStoreKinds) {
        var datastorekindsData interface{}
        if err := json.Unmarshal([]byte(data.DataStoreKinds.ValueString()), &datastorekindsData); err == nil {
            requestDataMap["dataStoreKinds"] = datastorekindsData
        } else {
            requestDataMap["dataStoreKinds"] = data.DataStoreKinds.ValueString()
        }
    }
    if !data.DataStoreProviders.IsUnknown() && !state.DataStoreProviders.IsUnknown() && !data.DataStoreProviders.Equal(state.DataStoreProviders) {
        var datastoreprovidersData interface{}
        if err := json.Unmarshal([]byte(data.DataStoreProviders.ValueString()), &datastoreprovidersData); err == nil {
            requestDataMap["dataStoreProviders"] = datastoreprovidersData
        } else {
            requestDataMap["dataStoreProviders"] = data.DataStoreProviders.ValueString()
        }
    }
    if !data.DataStoreEnvironments.IsUnknown() && !state.DataStoreEnvironments.IsUnknown() && !data.DataStoreEnvironments.Equal(state.DataStoreEnvironments) {
        var datastoreenvironmentsData interface{}
        if err := json.Unmarshal([]byte(data.DataStoreEnvironments.ValueString()), &datastoreenvironmentsData); err == nil {
            requestDataMap["dataStoreEnvironments"] = datastoreenvironmentsData
        } else {
            requestDataMap["dataStoreEnvironments"] = data.DataStoreEnvironments.ValueString()
        }
    }
    if !data.DataStoreLabels.IsUnknown() && !state.DataStoreLabels.IsUnknown() && !data.DataStoreLabels.Equal(state.DataStoreLabels) {
        requestDataMap["dataStoreLabels"] = r.convertTerraformSetToInterface(data.DataStoreLabels)
    }
    if !data.DataCategories.IsUnknown() && !state.DataCategories.IsUnknown() && !data.DataCategories.Equal(state.DataCategories) {
        var datacategoriesData interface{}
        if err := json.Unmarshal([]byte(data.DataCategories.ValueString()), &datacategoriesData); err == nil {
            requestDataMap["dataCategories"] = datacategoriesData
        } else {
            requestDataMap["dataCategories"] = data.DataCategories.ValueString()
        }
    }
    if !data.Runbooks.IsUnknown() && !state.Runbooks.IsUnknown() && !data.Runbooks.Equal(state.Runbooks) {
        requestDataMap["runbooks"] = r.convertTerraformSetToInterface(data.Runbooks)
    }
    if !data.ExternalMcpServers.IsUnknown() && !state.ExternalMcpServers.IsUnknown() && !data.ExternalMcpServers.Equal(state.ExternalMcpServers) {
        requestDataMap["externalMcpServers"] = r.convertTerraformSetToInterface(data.ExternalMcpServers)
    }
    if !data.RunnerId.IsUnknown() && !state.RunnerId.IsUnknown() && !data.RunnerId.Equal(state.RunnerId) {
        requestDataMap["runnerId"] = data.RunnerId.ValueString()
    }
    if !data.CodeRepositories.IsUnknown() && !state.CodeRepositories.IsUnknown() && !data.CodeRepositories.Equal(state.CodeRepositories) {
        requestDataMap["codeRepositories"] = r.convertTerraformSetToInterface(data.CodeRepositories)
    }
    if !data.LlmProviderId.IsUnknown() && !state.LlmProviderId.IsUnknown() && !data.LlmProviderId.Equal(state.LlmProviderId) {
        requestDataMap["llmProviderId"] = data.LlmProviderId.ValueString()
    }
    if !data.AiContext.IsUnknown() && !state.AiContext.IsUnknown() && !data.AiContext.Equal(state.AiContext) {
        requestDataMap["aiContext"] = data.AiContext.ValueString()
    }
    if !data.OpenCodeFixPullRequests.IsUnknown() && !state.OpenCodeFixPullRequests.IsUnknown() && !data.OpenCodeFixPullRequests.Equal(state.OpenCodeFixPullRequests) {
        requestDataMap["openCodeFixPullRequests"] = data.OpenCodeFixPullRequests.ValueBool()
    }
    if !data.OpenTelemetryPullRequests.IsUnknown() && !state.OpenTelemetryPullRequests.IsUnknown() && !data.OpenTelemetryPullRequests.Equal(state.OpenTelemetryPullRequests) {
        requestDataMap["openTelemetryPullRequests"] = data.OpenTelemetryPullRequests.ValueBool()
    }
    if !data.VerificationWindowMinutes.IsUnknown() && !state.VerificationWindowMinutes.IsUnknown() && !data.VerificationWindowMinutes.Equal(state.VerificationWindowMinutes) {
        requestDataMap["verificationWindowMinutes"] = r.bigFloatToFloat64(data.VerificationWindowMinutes.ValueBigFloat())
    }
    if !data.AutoResolveOnVerifiedRecovery.IsUnknown() && !state.AutoResolveOnVerifiedRecovery.IsUnknown() && !data.AutoResolveOnVerifiedRecovery.Equal(state.AutoResolveOnVerifiedRecovery) {
        requestDataMap["autoResolveOnVerifiedRecovery"] = data.AutoResolveOnVerifiedRecovery.ValueBool()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(aiRuleRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/auto-remediation-rule/" + data.Id.ValueString() + "", aiRuleRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ai_rule, got error: %s", err))
            return
        }

        // Parse the update response
        var aiRuleResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &aiRuleResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update ai_rule: %s", err))
            return
        }
        _ = aiRuleResponse

        // What the write does that someone should know, if anything (see entries.go).
        resp.Diagnostics.Append(WriteWarnings(httpResp, "ai rule")...)
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "triggerEntityType": true,
        "investigates": true,
        "onlyIfOpenAfterMinutes": true,
        "fixes": true,
        "executionMode": true,
        "remediationAction": true,
        "aiSelectsRunbook": true,
        "aiComposesCommands": true,
        "commandAllowlist": true,
        "commandRunners": true,
        "monitors": true,
        "incidentSeverities": true,
        "alertSeverities": true,
        "labels": true,
        "monitorLabels": true,
        "titlePattern": true,
        "descriptionPattern": true,
        "services": true,
        "serviceLabels": true,
        "serviceNamePattern": true,
        "dataSecurityAnalyze": true,
        "dataSecurityContain": true,
        "dataSecurityFix": true,
        "dataSecurityJudge": true,
        "dataSecurityChecks": true,
        "dataSecurityCheckGroups": true,
        "dataSecuritySeverities": true,
        "dataStoreSensitivities": true,
        "dataStoreKinds": true,
        "dataStoreProviders": true,
        "dataStoreEnvironments": true,
        "dataStoreLabels": true,
        "dataCategories": true,
        "runbooks": true,
        "externalMcpServers": true,
        "runnerId": true,
        "codeRepositories": true,
        "llmProviderId": true,
        "aiContext": true,
        "openCodeFixPullRequests": true,
        "openTelemetryPullRequests": true,
        "verificationWindowMinutes": true,
        "autoResolveOnVerifiedRecovery": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "emptiedLimits": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/auto-remediation-rule/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_rule after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_rule after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["criteria"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Criteria = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Criteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Criteria = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Criteria = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Criteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Criteria = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["criteria"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Criteria = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Criteria = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["criteria"].(string); ok {
        data.Criteria = NewJSONSubsetValue(val)
    } else {
        data.Criteria = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["triggerEntityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TriggerEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TriggerEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if arr, ok := dataMap["triggerEntityType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if val, ok := dataMap["triggerEntityType"].(string); ok {
        data.TriggerEntityType = types.StringValue(val)
    } else {
        data.TriggerEntityType = types.StringNull()
    }
    if val, ok := dataMap["investigates"].(bool); ok {
        data.Investigates = types.BoolValue(val)
    }
    data.OnlyIfOpenAfterMinutes = NumberFromAPI(data.OnlyIfOpenAfterMinutes, dataMap["onlyIfOpenAfterMinutes"])
    if val, ok := dataMap["fixes"].(bool); ok {
        data.Fixes = types.BoolValue(val)
    }
    if obj, ok := dataMap["executionMode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ExecutionMode = types.StringValue(string(jsonBytes))
            } else {
                data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ExecutionMode = types.StringValue(string(jsonBytes))
            } else {
                data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if arr, ok := dataMap["executionMode"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if val, ok := dataMap["executionMode"].(string); ok {
        data.ExecutionMode = types.StringValue(val)
    } else {
        data.ExecutionMode = types.StringNull()
    }
    if obj, ok := dataMap["remediationAction"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RemediationAction = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RemediationAction = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationAction = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RemediationAction = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationAction = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if arr, ok := dataMap["remediationAction"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if val, ok := dataMap["remediationAction"].(string); ok {
        data.RemediationAction = types.StringValue(val)
    } else {
        data.RemediationAction = types.StringNull()
    }
    if val, ok := dataMap["aiSelectsRunbook"].(bool); ok {
        data.AiSelectsRunbook = types.BoolValue(val)
    }
    if val, ok := dataMap["aiComposesCommands"].(bool); ok {
        data.AiComposesCommands = types.BoolValue(val)
    }
    if obj, ok := dataMap["commandAllowlist"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CommandAllowlist = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CommandAllowlist = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["commandAllowlist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CommandAllowlist = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["commandAllowlist"].(string); ok {
        data.CommandAllowlist = NewJSONSubsetValue(val)
    } else {
        data.CommandAllowlist = NewJSONSubsetNull()
    }
    if val, ok := dataMap["commandRunners"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.CommandRunners = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.CommandRunners = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Monitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["incidentSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.IncidentSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["alertSeverities"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AlertSeverities = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["monitorLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.MonitorLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["titlePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.TitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TitlePattern = types.StringValue(string(jsonBytes))
            } else {
                data.TitlePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if arr, ok := dataMap["titlePattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if val, ok := dataMap["titlePattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    } else {
        data.TitlePattern = types.StringNull()
    }
    if obj, ok := dataMap["descriptionPattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DescriptionPattern = types.StringValue(string(jsonBytes))
            } else {
                data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if arr, ok := dataMap["descriptionPattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if val, ok := dataMap["descriptionPattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    } else {
        data.DescriptionPattern = types.StringNull()
    }
    if val, ok := dataMap["services"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Services = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Services = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["serviceLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.ServiceLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.ServiceLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["serviceNamePattern"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceNamePattern = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if arr, ok := dataMap["serviceNamePattern"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if val, ok := dataMap["serviceNamePattern"].(string); ok {
        data.ServiceNamePattern = types.StringValue(val)
    } else {
        data.ServiceNamePattern = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityAnalyze"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityAnalyze = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityAnalyze = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityAnalyze = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityAnalyze = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityAnalyze"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityAnalyze = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityAnalyze = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityAnalyze"].(string); ok {
        data.DataSecurityAnalyze = types.StringValue(val)
    } else {
        data.DataSecurityAnalyze = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityContain"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityContain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityContain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityContain = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityContain = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityContain = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityContain = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityContain = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityContain"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityContain = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityContain = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityContain"].(string); ok {
        data.DataSecurityContain = types.StringValue(val)
    } else {
        data.DataSecurityContain = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityFix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityFix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityFix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityFix = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityFix = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityFix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityFix = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityFix = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityFix"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityFix = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityFix = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityFix"].(string); ok {
        data.DataSecurityFix = types.StringValue(val)
    } else {
        data.DataSecurityFix = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityJudge"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityJudge = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityJudge = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityJudge = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityJudge = types.StringValue(string(jsonBytes))
            } else {
                data.DataSecurityJudge = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityJudge = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityJudge = types.StringNull()
        }
    } else if arr, ok := dataMap["dataSecurityJudge"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityJudge = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityJudge = types.StringNull()
        }
    } else if val, ok := dataMap["dataSecurityJudge"].(string); ok {
        data.DataSecurityJudge = types.StringValue(val)
    } else {
        data.DataSecurityJudge = types.StringNull()
    }
    if obj, ok := dataMap["dataSecurityChecks"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityChecks = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityChecks = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityChecks = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityChecks = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecurityChecks"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityChecks = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityChecks = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecurityChecks"].(string); ok {
        data.DataSecurityChecks = NewJSONSubsetValue(val)
    } else {
        data.DataSecurityChecks = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataSecurityCheckGroups"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecurityCheckGroups = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityCheckGroups = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecurityCheckGroups"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecurityCheckGroups = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecurityCheckGroups = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecurityCheckGroups"].(string); ok {
        data.DataSecurityCheckGroups = NewJSONSubsetValue(val)
    } else {
        data.DataSecurityCheckGroups = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataSecuritySeverities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecuritySeverities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataSecuritySeverities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataSecuritySeverities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecuritySeverities = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataSecuritySeverities"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataSecuritySeverities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataSecuritySeverities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataSecuritySeverities"].(string); ok {
        data.DataSecuritySeverities = NewJSONSubsetValue(val)
    } else {
        data.DataSecuritySeverities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreSensitivities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreSensitivities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreSensitivities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreSensitivities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreSensitivities = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreSensitivities"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreSensitivities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreSensitivities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreSensitivities"].(string); ok {
        data.DataStoreSensitivities = NewJSONSubsetValue(val)
    } else {
        data.DataStoreSensitivities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreKinds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreKinds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreKinds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreKinds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreKinds = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreKinds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreKinds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreKinds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreKinds"].(string); ok {
        data.DataStoreKinds = NewJSONSubsetValue(val)
    } else {
        data.DataStoreKinds = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreProviders"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreProviders = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreProviders = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreProviders = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreProviders = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreProviders"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreProviders = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreProviders = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreProviders"].(string); ok {
        data.DataStoreProviders = NewJSONSubsetValue(val)
    } else {
        data.DataStoreProviders = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataStoreEnvironments"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreEnvironments = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataStoreEnvironments = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataStoreEnvironments = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreEnvironments = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataStoreEnvironments"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataStoreEnvironments = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataStoreEnvironments = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataStoreEnvironments"].(string); ok {
        data.DataStoreEnvironments = NewJSONSubsetValue(val)
    } else {
        data.DataStoreEnvironments = NewJSONSubsetNull()
    }
    if val, ok := dataMap["dataStoreLabels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DataStoreLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DataStoreLabels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["dataCategories"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["dataCategories"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataCategories"].(string); ok {
        data.DataCategories = NewJSONSubsetValue(val)
    } else {
        data.DataCategories = NewJSONSubsetNull()
    }
    if val, ok := dataMap["runbooks"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Runbooks = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Runbooks = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["externalMcpServers"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.ExternalMcpServers = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.ExternalMcpServers = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["runnerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RunnerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RunnerId = types.StringValue(string(jsonBytes))
            } else {
                data.RunnerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RunnerId = types.StringValue(string(jsonBytes))
            } else {
                data.RunnerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if arr, ok := dataMap["runnerId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if val, ok := dataMap["runnerId"].(string); ok {
        data.RunnerId = types.StringValue(val)
    } else {
        data.RunnerId = types.StringNull()
    }
    if val, ok := dataMap["codeRepositories"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.CodeRepositories = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.CodeRepositories = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["llmProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LlmProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LlmProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if arr, ok := dataMap["llmProviderId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["llmProviderId"].(string); ok {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := dataMap["aiContext"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiContext = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiContext = types.StringValue(string(jsonBytes))
            } else {
                data.AiContext = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiContext = types.StringValue(string(jsonBytes))
            } else {
                data.AiContext = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if arr, ok := dataMap["aiContext"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if val, ok := dataMap["aiContext"].(string); ok {
        data.AiContext = types.StringValue(val)
    } else {
        data.AiContext = types.StringNull()
    }
    if val, ok := dataMap["openCodeFixPullRequests"].(bool); ok {
        data.OpenCodeFixPullRequests = types.BoolValue(val)
    }
    if val, ok := dataMap["openTelemetryPullRequests"].(bool); ok {
        data.OpenTelemetryPullRequests = types.BoolValue(val)
    }
    data.VerificationWindowMinutes = NumberFromAPI(data.VerificationWindowMinutes, dataMap["verificationWindowMinutes"])
    if val, ok := dataMap["autoResolveOnVerifiedRecovery"].(bool); ok {
        data.AutoResolveOnVerifiedRecovery = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["createdByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["createdByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastModifiedAt = NewRFC3339Value(val)
        } else {
            data.LastModifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastModifiedAt"].(string); ok && val != "" {
        data.LastModifiedAt = NewRFC3339Value(val)
    } else {
        data.LastModifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastModifiedByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["emptiedLimits"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmptiedLimits = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EmptiedLimits = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.EmptiedLimits = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EmptiedLimits = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["emptiedLimits"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.EmptiedLimits = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.EmptiedLimits = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["emptiedLimits"].(string); ok {
        data.EmptiedLimits = NewJSONSubsetValue(val)
    } else {
        data.EmptiedLimits = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // What this update wrote is what the configuration says now (see importsettle.go).
    resp.Diagnostics.Append(KeepImportedFromServer(ctx, resp.Private, writtenFromConfig)...)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AiRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data AiRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/auto-remediation-rule/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ai_rule, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete ai_rule: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *AiRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
    // These hold the record as OneUptime keeps it, keys it fills in itself
    // included, until an update writes them (see importsettle.go).
    resp.Diagnostics.Append(MarkImportedFromServer(ctx, resp.Private, []string{"command_allowlist", "criteria", "data_categories", "data_security_check_groups", "data_security_checks", "data_security_severities", "data_store_environments", "data_store_kinds", "data_store_providers", "data_store_sensitivities"})...)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *AiRuleResource) keepPlannedValues(data *AiRuleResourceModel, plan *AiRuleResourceModel, config *AiRuleResourceModel) {
    if config.Criteria.IsNull() && !plan.Criteria.IsUnknown() {
        data.Criteria = plan.Criteria
    }
    if config.Description.IsNull() && !plan.Description.IsUnknown() {
        data.Description = plan.Description
    }
    if config.IsEnabled.IsNull() && !plan.IsEnabled.IsUnknown() {
        data.IsEnabled = plan.IsEnabled
    }
    if config.Investigates.IsNull() && !plan.Investigates.IsUnknown() {
        data.Investigates = plan.Investigates
    }
    if config.OnlyIfOpenAfterMinutes.IsNull() && !plan.OnlyIfOpenAfterMinutes.IsUnknown() {
        data.OnlyIfOpenAfterMinutes = plan.OnlyIfOpenAfterMinutes
    }
    if config.Fixes.IsNull() && !plan.Fixes.IsUnknown() {
        data.Fixes = plan.Fixes
    }
    if config.ExecutionMode.IsNull() && !plan.ExecutionMode.IsUnknown() {
        data.ExecutionMode = plan.ExecutionMode
    }
    if config.RemediationAction.IsNull() && !plan.RemediationAction.IsUnknown() {
        data.RemediationAction = plan.RemediationAction
    }
    if config.AiSelectsRunbook.IsNull() && !plan.AiSelectsRunbook.IsUnknown() {
        data.AiSelectsRunbook = plan.AiSelectsRunbook
    }
    if config.AiComposesCommands.IsNull() && !plan.AiComposesCommands.IsUnknown() {
        data.AiComposesCommands = plan.AiComposesCommands
    }
    if config.CommandAllowlist.IsNull() && !plan.CommandAllowlist.IsUnknown() {
        data.CommandAllowlist = plan.CommandAllowlist
    }
    if config.CommandRunners.IsNull() && !plan.CommandRunners.IsUnknown() {
        data.CommandRunners = plan.CommandRunners
    }
    if config.Monitors.IsNull() && !plan.Monitors.IsUnknown() {
        data.Monitors = plan.Monitors
    }
    if config.IncidentSeverities.IsNull() && !plan.IncidentSeverities.IsUnknown() {
        data.IncidentSeverities = plan.IncidentSeverities
    }
    if config.AlertSeverities.IsNull() && !plan.AlertSeverities.IsUnknown() {
        data.AlertSeverities = plan.AlertSeverities
    }
    if config.Labels.IsNull() && !plan.Labels.IsUnknown() {
        data.Labels = plan.Labels
    }
    if config.MonitorLabels.IsNull() && !plan.MonitorLabels.IsUnknown() {
        data.MonitorLabels = plan.MonitorLabels
    }
    if config.TitlePattern.IsNull() && !plan.TitlePattern.IsUnknown() {
        data.TitlePattern = plan.TitlePattern
    }
    if config.DescriptionPattern.IsNull() && !plan.DescriptionPattern.IsUnknown() {
        data.DescriptionPattern = plan.DescriptionPattern
    }
    if config.Services.IsNull() && !plan.Services.IsUnknown() {
        data.Services = plan.Services
    }
    if config.ServiceLabels.IsNull() && !plan.ServiceLabels.IsUnknown() {
        data.ServiceLabels = plan.ServiceLabels
    }
    if config.ServiceNamePattern.IsNull() && !plan.ServiceNamePattern.IsUnknown() {
        data.ServiceNamePattern = plan.ServiceNamePattern
    }
    if config.DataSecurityAnalyze.IsNull() && !plan.DataSecurityAnalyze.IsUnknown() {
        data.DataSecurityAnalyze = plan.DataSecurityAnalyze
    }
    if config.DataSecurityContain.IsNull() && !plan.DataSecurityContain.IsUnknown() {
        data.DataSecurityContain = plan.DataSecurityContain
    }
    if config.DataSecurityFix.IsNull() && !plan.DataSecurityFix.IsUnknown() {
        data.DataSecurityFix = plan.DataSecurityFix
    }
    if config.DataSecurityJudge.IsNull() && !plan.DataSecurityJudge.IsUnknown() {
        data.DataSecurityJudge = plan.DataSecurityJudge
    }
    if config.DataSecurityChecks.IsNull() && !plan.DataSecurityChecks.IsUnknown() {
        data.DataSecurityChecks = plan.DataSecurityChecks
    }
    if config.DataSecurityCheckGroups.IsNull() && !plan.DataSecurityCheckGroups.IsUnknown() {
        data.DataSecurityCheckGroups = plan.DataSecurityCheckGroups
    }
    if config.DataSecuritySeverities.IsNull() && !plan.DataSecuritySeverities.IsUnknown() {
        data.DataSecuritySeverities = plan.DataSecuritySeverities
    }
    if config.DataStoreSensitivities.IsNull() && !plan.DataStoreSensitivities.IsUnknown() {
        data.DataStoreSensitivities = plan.DataStoreSensitivities
    }
    if config.DataStoreKinds.IsNull() && !plan.DataStoreKinds.IsUnknown() {
        data.DataStoreKinds = plan.DataStoreKinds
    }
    if config.DataStoreProviders.IsNull() && !plan.DataStoreProviders.IsUnknown() {
        data.DataStoreProviders = plan.DataStoreProviders
    }
    if config.DataStoreEnvironments.IsNull() && !plan.DataStoreEnvironments.IsUnknown() {
        data.DataStoreEnvironments = plan.DataStoreEnvironments
    }
    if config.DataStoreLabels.IsNull() && !plan.DataStoreLabels.IsUnknown() {
        data.DataStoreLabels = plan.DataStoreLabels
    }
    if config.DataCategories.IsNull() && !plan.DataCategories.IsUnknown() {
        data.DataCategories = plan.DataCategories
    }
    if config.Runbooks.IsNull() && !plan.Runbooks.IsUnknown() {
        data.Runbooks = plan.Runbooks
    }
    if config.ExternalMcpServers.IsNull() && !plan.ExternalMcpServers.IsUnknown() {
        data.ExternalMcpServers = plan.ExternalMcpServers
    }
    if config.RunnerId.IsNull() && !plan.RunnerId.IsUnknown() {
        data.RunnerId = plan.RunnerId
    }
    if config.CodeRepositories.IsNull() && !plan.CodeRepositories.IsUnknown() {
        data.CodeRepositories = plan.CodeRepositories
    }
    if config.LlmProviderId.IsNull() && !plan.LlmProviderId.IsUnknown() {
        data.LlmProviderId = plan.LlmProviderId
    }
    if config.AiContext.IsNull() && !plan.AiContext.IsUnknown() {
        data.AiContext = plan.AiContext
    }
    if config.OpenCodeFixPullRequests.IsNull() && !plan.OpenCodeFixPullRequests.IsUnknown() {
        data.OpenCodeFixPullRequests = plan.OpenCodeFixPullRequests
    }
    if config.OpenTelemetryPullRequests.IsNull() && !plan.OpenTelemetryPullRequests.IsUnknown() {
        data.OpenTelemetryPullRequests = plan.OpenTelemetryPullRequests
    }
    if config.VerificationWindowMinutes.IsNull() && !plan.VerificationWindowMinutes.IsUnknown() {
        data.VerificationWindowMinutes = plan.VerificationWindowMinutes
    }
    if config.AutoResolveOnVerifiedRecovery.IsNull() && !plan.AutoResolveOnVerifiedRecovery.IsUnknown() {
        data.AutoResolveOnVerifiedRecovery = plan.AutoResolveOnVerifiedRecovery
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *AiRuleResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *AiRuleResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *AiRuleResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *AiRuleResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *AiRuleResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *AiRuleResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *AiRuleResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *AiRuleResource) isValidOneUptimeObjectType(typeStr string) bool {
    return ValidOneUptimeObjectTypes[typeStr]
}
