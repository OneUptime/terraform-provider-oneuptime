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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    AutoRemediationRuleId types.String `tfsdk:"auto_remediation_rule_id"`
    RuleNameSnapshot types.String `tfsdk:"rule_name_snapshot"`
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
        MarkdownDescription: "A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule. Look up an existing auto_remediation_suggestion by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "auto_remediation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "rule_name_snapshot": schema.StringAttribute{
                MarkdownDescription: "Name of the rule when this suggestion was created — survives rule deletion..",
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "runbook_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "runbook_name_snapshot": schema.StringAttribute{
                MarkdownDescription: "Name of the proposed runbook when this suggestion was created — survives runbook deletion..",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Lifecycle status: Planning, Suggested, Approved, AutoExecuted, Dismissed or NoneApplicable..",
                Computed: true,
            },
            "execution_mode": schema.StringAttribute{
                MarkdownDescription: "The rule's execution mode when this suggestion was created (Suggest or FullAuto)..",
                Computed: true,
            },
            "suggestion_type": schema.StringAttribute{
                MarkdownDescription: "Runbook suggestions propose starting a pre-authored runbook; CommandPlan suggestions carry an AI-composed command plan..",
                Computed: true,
            },
            "command_plan": schema.StringAttribute{
                MarkdownDescription: "The AI-composed command plan for CommandPlan suggestions, including per-command execution results once run..",
                Computed: true,
            },
            "rationale_markdown": schema.StringAttribute{
                MarkdownDescription: "Why this runbook was proposed — the AI planning run's reasoning for AI rules, or a short note for deterministic rules..",
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "runbook_execution_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "approved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "approved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "dismissed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "dismissed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "verification_status": schema.StringAttribute{
                MarkdownDescription: "Outcome verification after execution: Pending, Verified, Failed or Skipped. Empty until a runbook is started..",
                Computed: true,
            },
            "verification_deadline_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "verification_completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "verification_note": schema.StringAttribute{
                MarkdownDescription: "Why verification ended the way it did..",
                Computed: true,
            },
            "verification_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Snapshot of the rule's verification window when this suggestion was created..",
                Computed: true,
            },
            "auto_resolve_on_recovery": schema.BoolAttribute{
                MarkdownDescription: "Snapshot of the rule's auto-resolve-on-verified-recovery setting when this suggestion was created..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a auto_remediation_suggestion.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "autoRemediationRuleId": true,
        "ruleNameSnapshot": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No auto_remediation_suggestion found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No auto_remediation_suggestion found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one auto_remediation_suggestion matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
