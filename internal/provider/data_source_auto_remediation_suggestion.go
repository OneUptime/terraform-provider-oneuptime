package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoRemediationSuggestionDataSource{}

func NewAutoRemediationSuggestionDataSource() datasource.DataSource {
    return &AutoRemediationSuggestionDataSource{}
}

// AutoRemediationSuggestionDataSource defines the data source implementation.
type AutoRemediationSuggestionDataSource struct {
    client *Client
}

// AutoRemediationSuggestionDataSourceModel describes the data source data model.
type AutoRemediationSuggestionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    AutoRemediationRuleId types.String `tfsdk:"auto_remediation_rule_id"`
    RuleNameSnapshot types.String `tfsdk:"rule_name_snapshot"`
    KubernetesClusterId types.String `tfsdk:"kubernetes_cluster_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    AlertId types.String `tfsdk:"alert_id"`
    RunbookId types.String `tfsdk:"runbook_id"`
    RunbookNameSnapshot types.String `tfsdk:"runbook_name_snapshot"`
    Status types.String `tfsdk:"status"`
    ExecutionMode types.String `tfsdk:"execution_mode"`
    SuggestionType types.String `tfsdk:"suggestion_type"`
    CommandPlan types.String `tfsdk:"command_plan"`
    RationaleMarkdown types.String `tfsdk:"rationale_markdown"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    RunbookExecutionId types.String `tfsdk:"runbook_execution_id"`
    ApprovedByUserId types.String `tfsdk:"approved_by_user_id"`
    ApprovedAt types.String `tfsdk:"approved_at"`
    DismissedByUserId types.String `tfsdk:"dismissed_by_user_id"`
    DismissedAt types.String `tfsdk:"dismissed_at"`
    VerificationStatus types.String `tfsdk:"verification_status"`
    VerificationDeadlineAt types.String `tfsdk:"verification_deadline_at"`
    VerificationCompletedAt types.String `tfsdk:"verification_completed_at"`
    VerificationNote types.String `tfsdk:"verification_note"`
    VerificationWindowMinutes types.Number `tfsdk:"verification_window_minutes"`
    AutoResolveOnRecovery types.Bool `tfsdk:"auto_resolve_on_recovery"`
}

func (d *AutoRemediationSuggestionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_auto_remediation_suggestion"
}

func (d *AutoRemediationSuggestionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule. Look up an existing auto remediation suggestion by `id`, or by any of its other arguments (`ai_run_id`, `alert_id`, `approved_by_user_id`, ...): each one set must match, and exactly one auto remediation suggestion may match them all.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "auto_remediation_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the rule that produced this suggestion. The ID of a `oneuptime_auto_remediation_rule`.",
                Optional: true,
                Computed: true,
            },
            "rule_name_snapshot": schema.StringAttribute{
                MarkdownDescription: "Name of the rule when this suggestion was created — survives rule deletion.",
                Optional: true,
                Computed: true,
            },
            "kubernetes_cluster_id": schema.StringAttribute{
                MarkdownDescription: "ID of the cluster whose AI remediation mode produced this suggestion. The ID of a `oneuptime_kubernetes_cluster`.",
                Optional: true,
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "The kind of resource whose AI remediation mode produced this suggestion (DockerHost, PodmanHost, DockerSwarmCluster, ProxmoxCluster, VMwareVCenter, CephCluster, DatabaseServer or Host; resource-level remediation, no rule).",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "ID of the resource whose AI remediation mode produced this suggestion, in the table its resource type names.",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident this suggestion remediates. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the alert this suggestion remediates. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "runbook_id": schema.StringAttribute{
                MarkdownDescription: "ID of the proposed runbook. The ID of a `oneuptime_runbook`.",
                Optional: true,
                Computed: true,
            },
            "runbook_name_snapshot": schema.StringAttribute{
                MarkdownDescription: "Name of the proposed runbook when this suggestion was created — survives runbook deletion.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Lifecycle status: Planning, Suggested, Approved, AutoExecuted, Dismissed or NoneApplicable.",
                Optional: true,
                Computed: true,
            },
            "execution_mode": schema.StringAttribute{
                MarkdownDescription: "The rule's execution mode when this suggestion was created (Suggest or FullAuto).",
                Optional: true,
                Computed: true,
            },
            "suggestion_type": schema.StringAttribute{
                MarkdownDescription: "Runbook suggestions propose starting a pre-authored runbook; CommandPlan suggestions carry an AI-composed command plan.",
                Optional: true,
                Computed: true,
            },
            "command_plan": schema.StringAttribute{
                MarkdownDescription: "The AI-composed command plan for CommandPlan suggestions, including per-command execution results once run. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "rationale_markdown": schema.StringAttribute{
                MarkdownDescription: "Why this runbook was proposed — the AI planning run's reasoning for AI rules, or a short note for deterministic rules.",
                Optional: true,
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "The AI planning run that picked the runbook (AI rules only).",
                Optional: true,
                Computed: true,
            },
            "runbook_execution_id": schema.StringAttribute{
                MarkdownDescription: "The runbook execution started when this suggestion was approved or auto-executed.",
                Optional: true,
                Computed: true,
            },
            "approved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who approved this suggestion. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "approved_at": schema.StringAttribute{
                MarkdownDescription: "When this suggestion was approved.",
                Computed: true,
            },
            "dismissed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who dismissed this suggestion. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "dismissed_at": schema.StringAttribute{
                MarkdownDescription: "When this suggestion was dismissed.",
                Computed: true,
            },
            "verification_status": schema.StringAttribute{
                MarkdownDescription: "Outcome verification after execution: Pending, Verified, Failed or Skipped. Empty until a runbook is started.",
                Optional: true,
                Computed: true,
            },
            "verification_deadline_at": schema.StringAttribute{
                MarkdownDescription: "When the verification window closes — the monitors must be operational by this time.",
                Computed: true,
            },
            "verification_completed_at": schema.StringAttribute{
                MarkdownDescription: "When verification reached a terminal outcome.",
                Computed: true,
            },
            "verification_note": schema.StringAttribute{
                MarkdownDescription: "Why verification ended the way it did.",
                Optional: true,
                Computed: true,
            },
            "verification_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Snapshot of the rule's verification window when this suggestion was created.",
                Optional: true,
                Computed: true,
            },
            "auto_resolve_on_recovery": schema.BoolAttribute{
                MarkdownDescription: "Snapshot of the rule's auto-resolve-on-verified-recovery setting when this suggestion was created.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *AutoRemediationSuggestionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AutoRemediationSuggestionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AutoRemediationSuggestionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.AutoRemediationRuleId.IsNull() && !data.AutoRemediationRuleId.IsUnknown() {
        filters["autoRemediationRuleId"] = data.AutoRemediationRuleId.ValueString()
        filterNames = append(filterNames, "auto_remediation_rule_id = "+fmt.Sprintf("%q", data.AutoRemediationRuleId.ValueString()))
    }
    if !data.RuleNameSnapshot.IsNull() && !data.RuleNameSnapshot.IsUnknown() {
        filters["ruleNameSnapshot"] = data.RuleNameSnapshot.ValueString()
        filterNames = append(filterNames, "rule_name_snapshot = "+fmt.Sprintf("%q", data.RuleNameSnapshot.ValueString()))
    }
    if !data.KubernetesClusterId.IsNull() && !data.KubernetesClusterId.IsUnknown() {
        filters["kubernetesClusterId"] = data.KubernetesClusterId.ValueString()
        filterNames = append(filterNames, "kubernetes_cluster_id = "+fmt.Sprintf("%q", data.KubernetesClusterId.ValueString()))
    }
    if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
        filters["resourceType"] = data.ResourceType.ValueString()
        filterNames = append(filterNames, "resource_type = "+fmt.Sprintf("%q", data.ResourceType.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.AlertId.IsNull() && !data.AlertId.IsUnknown() {
        filters["alertId"] = data.AlertId.ValueString()
        filterNames = append(filterNames, "alert_id = "+fmt.Sprintf("%q", data.AlertId.ValueString()))
    }
    if !data.RunbookId.IsNull() && !data.RunbookId.IsUnknown() {
        filters["runbookId"] = data.RunbookId.ValueString()
        filterNames = append(filterNames, "runbook_id = "+fmt.Sprintf("%q", data.RunbookId.ValueString()))
    }
    if !data.RunbookNameSnapshot.IsNull() && !data.RunbookNameSnapshot.IsUnknown() {
        filters["runbookNameSnapshot"] = data.RunbookNameSnapshot.ValueString()
        filterNames = append(filterNames, "runbook_name_snapshot = "+fmt.Sprintf("%q", data.RunbookNameSnapshot.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.ExecutionMode.IsNull() && !data.ExecutionMode.IsUnknown() {
        filters["executionMode"] = data.ExecutionMode.ValueString()
        filterNames = append(filterNames, "execution_mode = "+fmt.Sprintf("%q", data.ExecutionMode.ValueString()))
    }
    if !data.SuggestionType.IsNull() && !data.SuggestionType.IsUnknown() {
        filters["suggestionType"] = data.SuggestionType.ValueString()
        filterNames = append(filterNames, "suggestion_type = "+fmt.Sprintf("%q", data.SuggestionType.ValueString()))
    }
    if !data.RationaleMarkdown.IsNull() && !data.RationaleMarkdown.IsUnknown() {
        filters["rationaleMarkdown"] = data.RationaleMarkdown.ValueString()
        filterNames = append(filterNames, "rationale_markdown = "+fmt.Sprintf("%q", data.RationaleMarkdown.ValueString()))
    }
    if !data.AiRunId.IsNull() && !data.AiRunId.IsUnknown() {
        filters["aiRunId"] = data.AiRunId.ValueString()
        filterNames = append(filterNames, "ai_run_id = "+fmt.Sprintf("%q", data.AiRunId.ValueString()))
    }
    if !data.RunbookExecutionId.IsNull() && !data.RunbookExecutionId.IsUnknown() {
        filters["runbookExecutionId"] = data.RunbookExecutionId.ValueString()
        filterNames = append(filterNames, "runbook_execution_id = "+fmt.Sprintf("%q", data.RunbookExecutionId.ValueString()))
    }
    if !data.ApprovedByUserId.IsNull() && !data.ApprovedByUserId.IsUnknown() {
        filters["approvedByUserId"] = data.ApprovedByUserId.ValueString()
        filterNames = append(filterNames, "approved_by_user_id = "+fmt.Sprintf("%q", data.ApprovedByUserId.ValueString()))
    }
    if !data.DismissedByUserId.IsNull() && !data.DismissedByUserId.IsUnknown() {
        filters["dismissedByUserId"] = data.DismissedByUserId.ValueString()
        filterNames = append(filterNames, "dismissed_by_user_id = "+fmt.Sprintf("%q", data.DismissedByUserId.ValueString()))
    }
    if !data.VerificationStatus.IsNull() && !data.VerificationStatus.IsUnknown() {
        filters["verificationStatus"] = data.VerificationStatus.ValueString()
        filterNames = append(filterNames, "verification_status = "+fmt.Sprintf("%q", data.VerificationStatus.ValueString()))
    }
    if !data.VerificationNote.IsNull() && !data.VerificationNote.IsUnknown() {
        filters["verificationNote"] = data.VerificationNote.ValueString()
        filterNames = append(filterNames, "verification_note = "+fmt.Sprintf("%q", data.VerificationNote.ValueString()))
    }
    if !data.VerificationWindowMinutes.IsNull() && !data.VerificationWindowMinutes.IsUnknown() {
        filters["verificationWindowMinutes"] = lookupNumber(data.VerificationWindowMinutes)
        filterNames = append(filterNames, "verification_window_minutes = "+data.VerificationWindowMinutes.ValueBigFloat().String())
    }
    if !data.AutoResolveOnRecovery.IsNull() && !data.AutoResolveOnRecovery.IsUnknown() {
        filters["autoResolveOnRecovery"] = data.AutoResolveOnRecovery.ValueBool()
        filterNames = append(filterNames, "auto_resolve_on_recovery = "+fmt.Sprintf("%t", data.AutoResolveOnRecovery.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the auto remediation suggestion up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the auto remediation suggestion up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "autoRemediationRuleId": true,
        "ruleNameSnapshot": true,
        "kubernetesClusterId": true,
        "resourceType": true,
        "resourceId": true,
        "incidentId": true,
        "alertId": true,
        "runbookId": true,
        "runbookNameSnapshot": true,
        "status": true,
        "executionMode": true,
        "suggestionType": true,
        "commandPlan": true,
        "rationaleMarkdown": true,
        "aiRunId": true,
        "runbookExecutionId": true,
        "approvedByUserId": true,
        "approvedAt": true,
        "dismissedByUserId": true,
        "dismissedAt": true,
        "verificationStatus": true,
        "verificationDeadlineAt": true,
        "verificationCompletedAt": true,
        "verificationNote": true,
        "verificationWindowMinutes": true,
        "autoResolveOnRecovery": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/auto-remediation-suggestion/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read auto_remediation_suggestion, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No auto remediation suggestion found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read auto_remediation_suggestion: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/auto-remediation-suggestion/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list auto_remediation_suggestion, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list auto_remediation_suggestion: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No auto remediation suggestion matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one auto remediation suggestion matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for auto_remediation_suggestion.")
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
    if obj, ok := item["autoRemediationRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoRemediationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoRemediationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoRemediationRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoRemediationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.AutoRemediationRuleId = types.StringNull()
        }
    } else if val, ok := item["autoRemediationRuleId"].(string); ok {
        data.AutoRemediationRuleId = types.StringValue(val)
    } else {
        data.AutoRemediationRuleId = types.StringNull()
    }
    if obj, ok := item["ruleNameSnapshot"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleNameSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleNameSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleNameSnapshot = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleNameSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.RuleNameSnapshot = types.StringNull()
        }
    } else if val, ok := item["ruleNameSnapshot"].(string); ok {
        data.RuleNameSnapshot = types.StringValue(val)
    } else {
        data.RuleNameSnapshot = types.StringNull()
    }
    if obj, ok := item["kubernetesClusterId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KubernetesClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KubernetesClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KubernetesClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesClusterId = types.StringNull()
        }
    } else if val, ok := item["kubernetesClusterId"].(string); ok {
        data.KubernetesClusterId = types.StringValue(val)
    } else {
        data.KubernetesClusterId = types.StringNull()
    }
    if obj, ok := item["resourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceType = types.StringNull()
        }
    } else if val, ok := item["resourceType"].(string); ok {
        data.ResourceType = types.StringValue(val)
    } else {
        data.ResourceType = types.StringNull()
    }
    if obj, ok := item["resourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if val, ok := item["resourceId"].(string); ok {
        data.ResourceId = types.StringValue(val)
    } else {
        data.ResourceId = types.StringNull()
    }
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := item["alertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := item["alertId"].(string); ok {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if obj, ok := item["runbookId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunbookId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunbookId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunbookId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunbookId = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookId = types.StringNull()
        }
    } else if val, ok := item["runbookId"].(string); ok {
        data.RunbookId = types.StringValue(val)
    } else {
        data.RunbookId = types.StringNull()
    }
    if obj, ok := item["runbookNameSnapshot"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunbookNameSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunbookNameSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunbookNameSnapshot = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunbookNameSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookNameSnapshot = types.StringNull()
        }
    } else if val, ok := item["runbookNameSnapshot"].(string); ok {
        data.RunbookNameSnapshot = types.StringValue(val)
    } else {
        data.RunbookNameSnapshot = types.StringNull()
    }
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
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
    if obj, ok := item["suggestionType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SuggestionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SuggestionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SuggestionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SuggestionType = types.StringValue(string(jsonBytes))
        } else {
            data.SuggestionType = types.StringNull()
        }
    } else if val, ok := item["suggestionType"].(string); ok {
        data.SuggestionType = types.StringValue(val)
    } else {
        data.SuggestionType = types.StringNull()
    }
    if obj, ok := item["commandPlan"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CommandPlan = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CommandPlan = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CommandPlan = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CommandPlan = types.StringValue(string(jsonBytes))
        } else {
            data.CommandPlan = types.StringNull()
        }
    } else if val, ok := item["commandPlan"].(string); ok {
        data.CommandPlan = types.StringValue(val)
    } else {
        data.CommandPlan = types.StringNull()
    }
    if obj, ok := item["rationaleMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RationaleMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RationaleMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RationaleMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RationaleMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.RationaleMarkdown = types.StringNull()
        }
    } else if val, ok := item["rationaleMarkdown"].(string); ok {
        data.RationaleMarkdown = types.StringValue(val)
    } else {
        data.RationaleMarkdown = types.StringNull()
    }
    if obj, ok := item["aiRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.AiRunId = types.StringNull()
        }
    } else if val, ok := item["aiRunId"].(string); ok {
        data.AiRunId = types.StringValue(val)
    } else {
        data.AiRunId = types.StringNull()
    }
    if obj, ok := item["runbookExecutionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunbookExecutionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunbookExecutionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunbookExecutionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunbookExecutionId = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookExecutionId = types.StringNull()
        }
    } else if val, ok := item["runbookExecutionId"].(string); ok {
        data.RunbookExecutionId = types.StringValue(val)
    } else {
        data.RunbookExecutionId = types.StringNull()
    }
    if obj, ok := item["approvedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApprovedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApprovedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApprovedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApprovedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ApprovedByUserId = types.StringNull()
        }
    } else if val, ok := item["approvedByUserId"].(string); ok {
        data.ApprovedByUserId = types.StringValue(val)
    } else {
        data.ApprovedByUserId = types.StringNull()
    }
    if obj, ok := item["approvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApprovedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApprovedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApprovedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApprovedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ApprovedAt = types.StringNull()
        }
    } else if val, ok := item["approvedAt"].(string); ok {
        data.ApprovedAt = types.StringValue(val)
    } else {
        data.ApprovedAt = types.StringNull()
    }
    if obj, ok := item["dismissedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DismissedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DismissedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DismissedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DismissedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DismissedByUserId = types.StringNull()
        }
    } else if val, ok := item["dismissedByUserId"].(string); ok {
        data.DismissedByUserId = types.StringValue(val)
    } else {
        data.DismissedByUserId = types.StringNull()
    }
    if obj, ok := item["dismissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DismissedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DismissedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DismissedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DismissedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DismissedAt = types.StringNull()
        }
    } else if val, ok := item["dismissedAt"].(string); ok {
        data.DismissedAt = types.StringValue(val)
    } else {
        data.DismissedAt = types.StringNull()
    }
    if obj, ok := item["verificationStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VerificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VerificationStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VerificationStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VerificationStatus = types.StringValue(string(jsonBytes))
        } else {
            data.VerificationStatus = types.StringNull()
        }
    } else if val, ok := item["verificationStatus"].(string); ok {
        data.VerificationStatus = types.StringValue(val)
    } else {
        data.VerificationStatus = types.StringNull()
    }
    if obj, ok := item["verificationDeadlineAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VerificationDeadlineAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VerificationDeadlineAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VerificationDeadlineAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VerificationDeadlineAt = types.StringValue(string(jsonBytes))
        } else {
            data.VerificationDeadlineAt = types.StringNull()
        }
    } else if val, ok := item["verificationDeadlineAt"].(string); ok {
        data.VerificationDeadlineAt = types.StringValue(val)
    } else {
        data.VerificationDeadlineAt = types.StringNull()
    }
    if obj, ok := item["verificationCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VerificationCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VerificationCompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VerificationCompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VerificationCompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.VerificationCompletedAt = types.StringNull()
        }
    } else if val, ok := item["verificationCompletedAt"].(string); ok {
        data.VerificationCompletedAt = types.StringValue(val)
    } else {
        data.VerificationCompletedAt = types.StringNull()
    }
    if obj, ok := item["verificationNote"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VerificationNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VerificationNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VerificationNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VerificationNote = types.StringValue(string(jsonBytes))
        } else {
            data.VerificationNote = types.StringNull()
        }
    } else if val, ok := item["verificationNote"].(string); ok {
        data.VerificationNote = types.StringValue(val)
    } else {
        data.VerificationNote = types.StringNull()
    }
    if val, ok := item["verificationWindowMinutes"].(float64); ok {
        data.VerificationWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["verificationWindowMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.VerificationWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.VerificationWindowMinutes = types.NumberNull()
        }
    } else {
        data.VerificationWindowMinutes = types.NumberNull()
    }
    if val, ok := item["autoResolveOnRecovery"].(bool); ok {
        data.AutoResolveOnRecovery = types.BoolValue(val)
    } else {
        data.AutoResolveOnRecovery = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
