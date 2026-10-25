package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AiRuleDataSource{}

func NewAiRuleDataSource() datasource.DataSource {
    return &AiRuleDataSource{}
}

// NewAiRuleLegacyDataSource registers this data source under the
// name it had before type names kept mixed-case words whole,
// oneuptime_auto_remediation_rule. Deprecated.
func NewAiRuleLegacyDataSource() datasource.DataSource {
    return &AiRuleDataSource{isLegacyAlias: true}
}

// AiRuleDataSource defines the data source implementation.
type AiRuleDataSource struct {
    client *Client
    // Registered under the name this data source had before.
    isLegacyAlias bool
}

// AiRuleDataSourceModel describes the data source data model.
type AiRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Criteria types.String `tfsdk:"criteria"`
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
    CommandAllowlist types.String `tfsdk:"command_allowlist"`
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AiRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    if d.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_auto_remediation_rule"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_ai_rule"
}

func (d *AiRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Which new incidents or alerts OneUptime AI investigates and fixes, and how. A rule investigates what it matches (right away, or only if still open after a wait), and fixes it or not, asking first or not. A fix may run the project's runbooks (only the rule's, when it names some), Bash commands on a Runner the rule picks, the tools of the project's external MCP servers (only the rule's, when it names some), and open fix and telemetry pull requests in the project's code repositories (only the rule's, when it names some). With no rule, OneUptime AI investigates every one, and fixes every one while automatic fixing is on. These were auto remediation rules, and investigation rules are AI rules now too. A rule whose Trigger Entity Type is Profile decides the same about the slow code OneUptime finds in a service's profiles, and one whose Trigger Entity Type is Trace about the slow and failing operations OneUptime finds in a service's traces, matched by their service; their fix is a pull request in one of the project's code repositories. Look up an existing ai rule by `id`, or by any of its other arguments (`name`, `ai_composes_commands`, `ai_context`, ...): each one set must match, and exactly one ai rule may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "criteria": schema.StringAttribute{
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this AI rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this AI rule.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Optional: true,
                Computed: true,
            },
            "trigger_entity_type": schema.StringAttribute{
                MarkdownDescription: "Which kind of signal this rule decides about: Incident or Alert - or Profile, the slow code OneUptime finds in a service's profiles, or Trace, the slow and failing operations OneUptime finds in a service's traces.",
                Optional: true,
                Computed: true,
            },
            "investigates": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI investigates the incidents or alerts this rule matches. With enabled rules that investigate, only their matches are investigated. Off, the rule only fixes, and does not change which incidents are investigated. The dashboard turns it on for every rule; left out of an API write, it is off.",
                Optional: true,
                Computed: true,
            },
            "only_if_open_after_minutes": schema.NumberAttribute{
                MarkdownDescription: "Investigate only incidents or alerts still open this many minutes after they were created: OneUptime AI waits that long, then investigates if it is still open. Leave empty to investigate right away. When several rules match, the shortest wait decides. At most 1440 minutes (a day).",
                Optional: true,
                Computed: true,
            },
            "fixes": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI fixes the incidents or alerts this rule matches, the way its Execution Mode says. Off, the rule only investigates them. With enabled rules, an incident is fixed only when a rule that fixes matches it. On by default.",
                Optional: true,
                Computed: true,
            },
            "execution_mode": schema.StringAttribute{
                MarkdownDescription: "How a rule that fixes asks. Suggest asks before fixing: every fix the rule starts waits for one-click human approval. FullAuto fixes without asking: OneUptime AI fixes on its own where the cluster's or resource's AI agent page allows, a runbook it picks from the ones the rule names starts right away (one it picks from all of the project's runbooks still waits for approval), a change through an MCP server the rule names is made right away (through any of the project's MCP servers, it still waits for approval), and a rule that runs its runbooks every time starts them immediately.",
                Optional: true,
                Computed: true,
            },
            "remediation_action": schema.StringAttribute{
                MarkdownDescription: "OneUptimeAI (the default): OneUptime AI fixes the matched incident or alert - on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, with the project's runbooks (only the rule's Runbooks when it names some) and with commands on the rule's Runner when it picks one. Runbooks: the rule starts every one of its Runbooks for every match, without OneUptime AI deciding. Whether a person approves first is the rule's Execution Mode.",
                Optional: true,
                Computed: true,
            },
            "ai_selects_runbook": schema.BoolAttribute{
                MarkdownDescription: "Kept for rules saved before rules were simplified: an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached), always proposing it for approval, and the rule does nothing else. Every OneUptime AI rule now lets AI pick a runbook on its own, so new rules leave this off.",
                Optional: true,
                Computed: true,
            },
            "ai_composes_commands": schema.BoolAttribute{
                MarkdownDescription: "Kept for rules saved before rules were simplified: the AI investigates the incident/alert and composes Bash commands for the rule's Runner and Command Runners, and the rule does nothing else. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. A OneUptime AI rule now composes commands on its Runner when it picks one, so new rules leave this off.",
                Optional: true,
                Computed: true,
            },
            "command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "command_runners": schema.SetAttribute{
                MarkdownDescription: "Kept for rules saved before rules were simplified: more Runners the AI may target with composed commands, besides the rule's Runner. A Runner runs a rule's commands only when the rule names it. IDs of `oneuptime_runner` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "title_pattern": schema.StringAttribute{
                MarkdownDescription: "Case-insensitive regex matched against the entity's title. Leave empty to match any title.",
                Optional: true,
                Computed: true,
            },
            "description_pattern": schema.StringAttribute{
                MarkdownDescription: "Case-insensitive regex matched against the entity's description. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "services": schema.SetAttribute{
                MarkdownDescription: "Profile and Trace AI rules only: match what OneUptime found in one of these services. Leave empty to match any service. IDs of `oneuptime_service` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_labels": schema.SetAttribute{
                MarkdownDescription: "Profile and Trace AI rules only: match what OneUptime found in a service that carries at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Profile and Trace AI rules only: case-insensitive regex matched against the name of the service OneUptime found the problem in. Leave empty to match any name.",
                Optional: true,
                Computed: true,
            },
            "runbooks": schema.SetAttribute{
                MarkdownDescription: "The runbooks OneUptime AI may pick and run for this rule. Leave empty to let it pick any of the project's runbooks. A rule whose Fix With is Runbooks starts every one of them for every match instead. IDs of `oneuptime_runbook` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "external_mcp_servers": schema.SetAttribute{
                MarkdownDescription: "The external MCP servers OneUptime AI may use while it investigates and fixes what this rule matches (an investigation uses only their tools that read). Leave empty to let it use all of the project's MCP servers; changes through them are then always proposed for approval first. IDs of `oneuptime_external_mcp_server` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "runner_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Runner OneUptime AI runs this rule's Bash commands on: one of the project's own Runners. Empty means none: OneUptime AI runs no commands of its own for the rule. Commands on the affected Kubernetes cluster or host always go through its AI agent. The ID of a `oneuptime_runner`.",
                Optional: true,
                Computed: true,
            },
            "code_repositories": schema.SetAttribute{
                MarkdownDescription: "The code repositories OneUptime AI may read while it investigates and fixes what this rule matches, and open the rule's pull requests in. Leave empty to let it use every repository the project has connected. IDs of `oneuptime_code_repository` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "llm_provider_id": schema.StringAttribute{
                MarkdownDescription: "ID of the LLM provider OneUptime AI uses for everything this rule starts. Empty: the project's default provider, else the global one. The ID of a `oneuptime_llm_provider`.",
                Optional: true,
                Computed: true,
            },
            "ai_context": schema.StringAttribute{
                MarkdownDescription: "Notes OneUptime AI reads after the project's AI context when it investigates or fixes what this rule matches. Markdown, at most 5,000 characters.",
                Optional: true,
                Computed: true,
            },
            "open_code_fix_pull_requests": schema.BoolAttribute{
                MarkdownDescription: "When an investigation of an incident or alert this rule matches finds the cause in your code, OneUptime AI opens a pull request with the fix, for review. Nothing is merged. Acts only while fixing new incidents (or alerts) automatically is on. On by default.",
                Optional: true,
                Computed: true,
            },
            "open_telemetry_pull_requests": schema.BoolAttribute{
                MarkdownDescription: "When missing logs, traces or metrics stop an investigation of an incident or alert this rule matches, OneUptime AI opens a pull request that adds them, for review. Nothing is merged. Acts only while fixing new incidents (or alerts) automatically is on. On by default.",
                Optional: true,
                Computed: true,
            },
            "verification_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes.",
                Optional: true,
                Computed: true,
            },
            "auto_resolve_on_verified_recovery": schema.BoolAttribute{
                MarkdownDescription: "When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
    if d.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_auto_remediation_rule has been renamed to oneuptime_ai_rule. Use the new name; the old one keeps working until then."
    }
}

func (d *AiRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *AiRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiRuleDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.TriggerEntityType.IsNull() && !data.TriggerEntityType.IsUnknown() {
        filters["triggerEntityType"] = data.TriggerEntityType.ValueString()
        filterNames = append(filterNames, "trigger_entity_type = "+fmt.Sprintf("%q", data.TriggerEntityType.ValueString()))
    }
    if !data.Investigates.IsNull() && !data.Investigates.IsUnknown() {
        filters["investigates"] = data.Investigates.ValueBool()
        filterNames = append(filterNames, "investigates = "+fmt.Sprintf("%t", data.Investigates.ValueBool()))
    }
    if !data.OnlyIfOpenAfterMinutes.IsNull() && !data.OnlyIfOpenAfterMinutes.IsUnknown() {
        filters["onlyIfOpenAfterMinutes"] = lookupNumber(data.OnlyIfOpenAfterMinutes)
        filterNames = append(filterNames, "only_if_open_after_minutes = "+data.OnlyIfOpenAfterMinutes.ValueBigFloat().String())
    }
    if !data.Fixes.IsNull() && !data.Fixes.IsUnknown() {
        filters["fixes"] = data.Fixes.ValueBool()
        filterNames = append(filterNames, "fixes = "+fmt.Sprintf("%t", data.Fixes.ValueBool()))
    }
    if !data.ExecutionMode.IsNull() && !data.ExecutionMode.IsUnknown() {
        filters["executionMode"] = data.ExecutionMode.ValueString()
        filterNames = append(filterNames, "execution_mode = "+fmt.Sprintf("%q", data.ExecutionMode.ValueString()))
    }
    if !data.RemediationAction.IsNull() && !data.RemediationAction.IsUnknown() {
        filters["remediationAction"] = data.RemediationAction.ValueString()
        filterNames = append(filterNames, "remediation_action = "+fmt.Sprintf("%q", data.RemediationAction.ValueString()))
    }
    if !data.AiSelectsRunbook.IsNull() && !data.AiSelectsRunbook.IsUnknown() {
        filters["aiSelectsRunbook"] = data.AiSelectsRunbook.ValueBool()
        filterNames = append(filterNames, "ai_selects_runbook = "+fmt.Sprintf("%t", data.AiSelectsRunbook.ValueBool()))
    }
    if !data.AiComposesCommands.IsNull() && !data.AiComposesCommands.IsUnknown() {
        filters["aiComposesCommands"] = data.AiComposesCommands.ValueBool()
        filterNames = append(filterNames, "ai_composes_commands = "+fmt.Sprintf("%t", data.AiComposesCommands.ValueBool()))
    }
    if !data.TitlePattern.IsNull() && !data.TitlePattern.IsUnknown() {
        filters["titlePattern"] = data.TitlePattern.ValueString()
        filterNames = append(filterNames, "title_pattern = "+fmt.Sprintf("%q", data.TitlePattern.ValueString()))
    }
    if !data.DescriptionPattern.IsNull() && !data.DescriptionPattern.IsUnknown() {
        filters["descriptionPattern"] = data.DescriptionPattern.ValueString()
        filterNames = append(filterNames, "description_pattern = "+fmt.Sprintf("%q", data.DescriptionPattern.ValueString()))
    }
    if !data.ServiceNamePattern.IsNull() && !data.ServiceNamePattern.IsUnknown() {
        filters["serviceNamePattern"] = data.ServiceNamePattern.ValueString()
        filterNames = append(filterNames, "service_name_pattern = "+fmt.Sprintf("%q", data.ServiceNamePattern.ValueString()))
    }
    if !data.RunnerId.IsNull() && !data.RunnerId.IsUnknown() {
        filters["runnerId"] = data.RunnerId.ValueString()
        filterNames = append(filterNames, "runner_id = "+fmt.Sprintf("%q", data.RunnerId.ValueString()))
    }
    if !data.LlmProviderId.IsNull() && !data.LlmProviderId.IsUnknown() {
        filters["llmProviderId"] = data.LlmProviderId.ValueString()
        filterNames = append(filterNames, "llm_provider_id = "+fmt.Sprintf("%q", data.LlmProviderId.ValueString()))
    }
    if !data.AiContext.IsNull() && !data.AiContext.IsUnknown() {
        filters["aiContext"] = data.AiContext.ValueString()
        filterNames = append(filterNames, "ai_context = "+fmt.Sprintf("%q", data.AiContext.ValueString()))
    }
    if !data.OpenCodeFixPullRequests.IsNull() && !data.OpenCodeFixPullRequests.IsUnknown() {
        filters["openCodeFixPullRequests"] = data.OpenCodeFixPullRequests.ValueBool()
        filterNames = append(filterNames, "open_code_fix_pull_requests = "+fmt.Sprintf("%t", data.OpenCodeFixPullRequests.ValueBool()))
    }
    if !data.OpenTelemetryPullRequests.IsNull() && !data.OpenTelemetryPullRequests.IsUnknown() {
        filters["openTelemetryPullRequests"] = data.OpenTelemetryPullRequests.ValueBool()
        filterNames = append(filterNames, "open_telemetry_pull_requests = "+fmt.Sprintf("%t", data.OpenTelemetryPullRequests.ValueBool()))
    }
    if !data.VerificationWindowMinutes.IsNull() && !data.VerificationWindowMinutes.IsUnknown() {
        filters["verificationWindowMinutes"] = lookupNumber(data.VerificationWindowMinutes)
        filterNames = append(filterNames, "verification_window_minutes = "+data.VerificationWindowMinutes.ValueBigFloat().String())
    }
    if !data.AutoResolveOnVerifiedRecovery.IsNull() && !data.AutoResolveOnVerifiedRecovery.IsUnknown() {
        filters["autoResolveOnVerifiedRecovery"] = data.AutoResolveOnVerifiedRecovery.ValueBool()
        filterNames = append(filterNames, "auto_resolve_on_verified_recovery = "+fmt.Sprintf("%t", data.AutoResolveOnVerifiedRecovery.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the ai rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the ai rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
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
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/auto-remediation-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_rule: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/auto-remediation-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list ai_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list ai_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for ai_rule.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["criteria"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Criteria = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Criteria = types.StringValue(string(jsonBytes))
        } else {
            data.Criteria = types.StringNull()
        }
    } else if val, ok := item["criteria"].(string); ok {
        data.Criteria = types.StringValue(val)
    } else {
        data.Criteria = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["triggerEntityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggerEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggerEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggerEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerEntityType = types.StringNull()
        }
    } else if val, ok := item["triggerEntityType"].(string); ok {
        data.TriggerEntityType = types.StringValue(val)
    } else {
        data.TriggerEntityType = types.StringNull()
    }
    if val, ok := item["investigates"].(bool); ok {
        data.Investigates = types.BoolValue(val)
    } else {
        data.Investigates = types.BoolNull()
    }
    data.OnlyIfOpenAfterMinutes = numberFromAPI(data.OnlyIfOpenAfterMinutes, item["onlyIfOpenAfterMinutes"])
    if val, ok := item["fixes"].(bool); ok {
        data.Fixes = types.BoolValue(val)
    } else {
        data.Fixes = types.BoolNull()
    }
    if obj, ok := item["executionMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExecutionMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExecutionMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExecutionMode = types.StringValue(string(jsonBytes))
        } else {
            data.ExecutionMode = types.StringNull()
        }
    } else if val, ok := item["executionMode"].(string); ok {
        data.ExecutionMode = types.StringValue(val)
    } else {
        data.ExecutionMode = types.StringNull()
    }
    if obj, ok := item["remediationAction"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RemediationAction = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RemediationAction = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RemediationAction = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationAction = types.StringNull()
        }
    } else if val, ok := item["remediationAction"].(string); ok {
        data.RemediationAction = types.StringValue(val)
    } else {
        data.RemediationAction = types.StringNull()
    }
    if val, ok := item["aiSelectsRunbook"].(bool); ok {
        data.AiSelectsRunbook = types.BoolValue(val)
    } else {
        data.AiSelectsRunbook = types.BoolNull()
    }
    if val, ok := item["aiComposesCommands"].(bool); ok {
        data.AiComposesCommands = types.BoolValue(val)
    } else {
        data.AiComposesCommands = types.BoolNull()
    }
    if obj, ok := item["commandAllowlist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CommandAllowlist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CommandAllowlist = types.StringValue(string(jsonBytes))
        } else {
            data.CommandAllowlist = types.StringNull()
        }
    } else if val, ok := item["commandAllowlist"].(string); ok {
        data.CommandAllowlist = types.StringValue(val)
    } else {
        data.CommandAllowlist = types.StringNull()
    }
    if val, ok := item["commandRunners"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.CommandRunners = types.SetValueMust(types.StringType, setItems)
    } else {
        data.CommandRunners = types.SetNull(types.StringType)
    }
    if val, ok := item["monitors"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentSeverities"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["alertSeverities"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if val, ok := item["monitorLabels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MonitorLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["titlePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if val, ok := item["titlePattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    } else {
        data.TitlePattern = types.StringNull()
    }
    if obj, ok := item["descriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["descriptionPattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    } else {
        data.DescriptionPattern = types.StringNull()
    }
    if val, ok := item["services"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Services = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Services = types.SetNull(types.StringType)
    }
    if val, ok := item["serviceLabels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.ServiceLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.ServiceLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["serviceNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamePattern = types.StringNull()
        }
    } else if val, ok := item["serviceNamePattern"].(string); ok {
        data.ServiceNamePattern = types.StringValue(val)
    } else {
        data.ServiceNamePattern = types.StringNull()
    }
    if val, ok := item["runbooks"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Runbooks = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Runbooks = types.SetNull(types.StringType)
    }
    if val, ok := item["externalMcpServers"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.ExternalMcpServers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.ExternalMcpServers = types.SetNull(types.StringType)
    }
    if obj, ok := item["runnerId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunnerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunnerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunnerId = types.StringValue(string(jsonBytes))
        } else {
            data.RunnerId = types.StringNull()
        }
    } else if val, ok := item["runnerId"].(string); ok {
        data.RunnerId = types.StringValue(val)
    } else {
        data.RunnerId = types.StringNull()
    }
    if val, ok := item["codeRepositories"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.CodeRepositories = types.SetValueMust(types.StringType, setItems)
    } else {
        data.CodeRepositories = types.SetNull(types.StringType)
    }
    if obj, ok := item["llmProviderId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmProviderId = types.StringNull()
        }
    } else if val, ok := item["llmProviderId"].(string); ok {
        data.LlmProviderId = types.StringValue(val)
    } else {
        data.LlmProviderId = types.StringNull()
    }
    if obj, ok := item["aiContext"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiContext = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiContext = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiContext = types.StringValue(string(jsonBytes))
        } else {
            data.AiContext = types.StringNull()
        }
    } else if val, ok := item["aiContext"].(string); ok {
        data.AiContext = types.StringValue(val)
    } else {
        data.AiContext = types.StringNull()
    }
    if val, ok := item["openCodeFixPullRequests"].(bool); ok {
        data.OpenCodeFixPullRequests = types.BoolValue(val)
    } else {
        data.OpenCodeFixPullRequests = types.BoolNull()
    }
    if val, ok := item["openTelemetryPullRequests"].(bool); ok {
        data.OpenTelemetryPullRequests = types.BoolValue(val)
    } else {
        data.OpenTelemetryPullRequests = types.BoolNull()
    }
    data.VerificationWindowMinutes = numberFromAPI(data.VerificationWindowMinutes, item["verificationWindowMinutes"])
    if val, ok := item["autoResolveOnVerifiedRecovery"].(bool); ok {
        data.AutoResolveOnVerifiedRecovery = types.BoolValue(val)
    } else {
        data.AutoResolveOnVerifiedRecovery = types.BoolNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
