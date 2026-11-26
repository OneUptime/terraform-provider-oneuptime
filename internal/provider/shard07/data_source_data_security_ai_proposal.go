package shard07

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSecurityAiProposalDataSource{}

func NewDataSecurityAiProposalDataSource() datasource.DataSource {
    return &DataSecurityAiProposalDataSource{}
}

// DataSecurityAiProposalDataSource defines the data source implementation.
type DataSecurityAiProposalDataSource struct {
    client *Client
}

// DataSecurityAiProposalDataSourceModel describes the data source data model.
type DataSecurityAiProposalDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    DataSecurityIssueId types.String `tfsdk:"data_security_issue_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    DataSecurityIssueGroupId types.String `tfsdk:"data_security_issue_group_id"`
    Kind types.String `tfsdk:"kind"`
    Work types.String `tfsdk:"work"`
    Safety types.String `tfsdk:"safety"`
    Status types.String `tfsdk:"status"`
    Title types.String `tfsdk:"title"`
    ChangeMarkdown types.String `tfsdk:"change_markdown"`
    ReasonMarkdown types.String `tfsdk:"reason_markdown"`
    ProofMarkdown types.String `tfsdk:"proof_markdown"`
    UndoMarkdown types.String `tfsdk:"undo_markdown"`
    Plan types.String `tfsdk:"plan"`
    PlanKey types.String `tfsdk:"plan_key"`
    PlanFingerprint types.String `tfsdk:"plan_fingerprint"`
    Result types.String `tfsdk:"result"`
    RiskScore types.Number `tfsdk:"risk_score"`
    Autonomy types.String `tfsdk:"autonomy"`
    AskReason types.String `tfsdk:"ask_reason"`
    AutoRemediationRuleId types.String `tfsdk:"auto_remediation_rule_id"`
    AutoRemediationRuleName types.String `tfsdk:"auto_remediation_rule_name"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    FirstApprovedByUserId types.String `tfsdk:"first_approved_by_user_id"`
    FirstApprovedAt types.String `tfsdk:"first_approved_at"`
    FirstApprovedPlanFingerprint types.String `tfsdk:"first_approved_plan_fingerprint"`
    ApprovedByUserId types.String `tfsdk:"approved_by_user_id"`
    ApprovedAt types.String `tfsdk:"approved_at"`
    AppliedAt types.String `tfsdk:"applied_at"`
    DismissedByUserId types.String `tfsdk:"dismissed_by_user_id"`
    DismissedAt types.String `tfsdk:"dismissed_at"`
    DismissReason types.String `tfsdk:"dismiss_reason"`
    UndoneByUserId types.String `tfsdk:"undone_by_user_id"`
    UndoneAt types.String `tfsdk:"undone_at"`
    ExpiredAt types.String `tfsdk:"expired_at"`
    FailedAt types.String `tfsdk:"failed_at"`
    FailureMessage types.String `tfsdk:"failure_message"`
    VerifiedAt types.String `tfsdk:"verified_at"`
}

func (d *DataSecurityAiProposalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_security_ai_proposal"
}

func (d *DataSecurityAiProposalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What OneUptime AI proposes for a Data Security issue, or for who owns the project's stores: what will change, why, how it will be proven and how to undo it, and whether it waits for a person, was made, dismissed or undone. Written by OneUptime; approved, dismissed and undone through Data Security's own routes, with the person's own permission for the change. Look up an existing data security ai proposal by `id`, or by any of its other arguments (`ai_run_id`, `approved_by_user_id`, `ask_reason`, ...): each one set must match, and exactly one data security ai proposal may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this proposal belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_security_issue_id": schema.StringAttribute{
                MarkdownDescription: "The issue the proposal is for. Empty for an owner rule, which is for no one issue, and for the owners of a store that has no issue (#10374). The ID of a `oneuptime_data_security_issue` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_data_security_issue"),
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store the proposal is for: the issue's store, or the store whose owners it proposes. Empty for an owner rule, which is for the project's stores (#10374). The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_data_store"),
            },
            "data_security_issue_group_id": schema.StringAttribute{
                MarkdownDescription: "The group of issues whose cause the proposal fixes once (#10370): one pull request, or one approval for the stores made in the console. Such a proposal names no one issue.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What OneUptime AI proposes: ScrubTelemetry, the scrub rules that hide a telemetry issue's data from now on.",
                Optional: true,
                Computed: true,
            },
            "work": schema.StringAttribute{
                MarkdownDescription: "The kind of work it is, whose autonomy the project's AI rules set: Analyze, Contain, Fix or Judge.",
                Optional: true,
                Computed: true,
            },
            "safety": schema.StringAttribute{
                MarkdownDescription: "What kind of change it is, which decides what may run without asking: MaskOnly (hides data in OneUptime, undone by deleting what it made), OpensPullRequest, RecordsJudgement, ChangesAccess, Deletes, Pages, ChangesWhoIsTold or ChangesCloudSettings. Deleting data, paging people, changing who is told about a store and changing a cloud store's settings always ask.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Where it stands: Proposed (waiting for a person), Approved (being made), Applied, AppliedWithoutAsking, Dismissed, Expired (no longer needed), Undone or Failed.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "What OneUptime AI proposes, in a line.",
                Optional: true,
                Computed: true,
            },
            "change_markdown": schema.StringAttribute{
                MarkdownDescription: "A preview of the change a person can read before approving it: each rule it makes, and where.",
                Optional: true,
                Computed: true,
            },
            "reason_markdown": schema.StringAttribute{
                MarkdownDescription: "Why OneUptime AI proposes it, with the facts it rests on: the check, the data types and where they were found - never a value.",
                Optional: true,
                Computed: true,
            },
            "proof_markdown": schema.StringAttribute{
                MarkdownDescription: "How the change is proven to have worked: the issue's own check, run again - never OneUptime AI's word.",
                Optional: true,
                Computed: true,
            },
            "undo_markdown": schema.StringAttribute{
                MarkdownDescription: "What Undo does, and what it cannot take back.",
                Optional: true,
                Computed: true,
            },
            "plan": schema.StringAttribute{
                MarkdownDescription: "The change, as its action reads it to make it: for ScrubTelemetry, each finding's data type and place. Names and ids only. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "plan_key": schema.StringAttribute{
                MarkdownDescription: "What makes two proposals the same change: OneUptime AI never proposes one for an issue again, unless the last one expired.",
                Optional: true,
                Computed: true,
            },
            "plan_fingerprint": schema.StringAttribute{
                MarkdownDescription: "A fingerprint of the plan: an approval names the one its person was shown, and a plan that changed since is not approved.",
                Optional: true,
                Computed: true,
            },
            "result": schema.StringAttribute{
                MarkdownDescription: "What making it made, for Undo to take back: for ScrubTelemetry, the scrub rules it created. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "risk_score": schema.NumberAttribute{
                MarkdownDescription: "The issue's risk score when OneUptime AI proposed it: what waits for a person is listed riskiest first.",
                Optional: true,
                Computed: true,
            },
            "autonomy": schema.StringAttribute{
                MarkdownDescription: "What the project's AI rules (or its defaults) allowed for this kind of work when it was proposed: Suggest or WithoutAsking.",
                Optional: true,
                Computed: true,
            },
            "ask_reason": schema.StringAttribute{
                MarkdownDescription: "Why a change OneUptime AI could have made on its own waits for a person: the work asks first, the change always asks, a rule must name it, it could not prove it is safe, two people approve it, the hour's circuit breaker, or the repository's cap on open AI pull requests.",
                Optional: true,
                Computed: true,
            },
            "auto_remediation_rule_id": schema.StringAttribute{
                MarkdownDescription: "The AI rule that allowed it; empty when the project's defaults did.",
                Optional: true,
                Computed: true,
            },
            "auto_remediation_rule_name": schema.StringAttribute{
                MarkdownDescription: "What the AI rule that allowed it was called when it was proposed; empty when the project's defaults did.",
                Optional: true,
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "The AI run that made the proposal, when one did.",
                Optional: true,
                Computed: true,
            },
            "first_approved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "The person who gave the first of two approvals. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "first_approved_at": schema.StringAttribute{
                MarkdownDescription: "When the first of two approvals was given.",
                Computed: true,
            },
            "first_approved_plan_fingerprint": schema.StringAttribute{
                MarkdownDescription: "The plan the first approval was given to: a second approval runs only that plan.",
                Optional: true,
                Computed: true,
            },
            "approved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "The person who approved it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "approved_at": schema.StringAttribute{
                MarkdownDescription: "When it was approved.",
                Computed: true,
            },
            "applied_at": schema.StringAttribute{
                MarkdownDescription: "When the change was made: after approval, or without asking.",
                Computed: true,
            },
            "dismissed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "The person who dismissed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "dismissed_at": schema.StringAttribute{
                MarkdownDescription: "When it was dismissed.",
                Computed: true,
            },
            "dismiss_reason": schema.StringAttribute{
                MarkdownDescription: "Why the person dismissed it.",
                Optional: true,
                Computed: true,
            },
            "undone_by_user_id": schema.StringAttribute{
                MarkdownDescription: "The person who undid the change. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "undone_at": schema.StringAttribute{
                MarkdownDescription: "When the change was undone.",
                Computed: true,
            },
            "expired_at": schema.StringAttribute{
                MarkdownDescription: "When it stopped being needed: its issue resolved, or was accepted or set aside, before anyone acted on it.",
                Computed: true,
            },
            "failed_at": schema.StringAttribute{
                MarkdownDescription: "When making the change failed.",
                Computed: true,
            },
            "failure_message": schema.StringAttribute{
                MarkdownDescription: "What stopped the change, in the words of what refused it. Nothing changed.",
                Optional: true,
                Computed: true,
            },
            "verified_at": schema.StringAttribute{
                MarkdownDescription: "When the issue's own check passed after the change was made: the proof that it worked.",
                Computed: true,
            },
        },
    }
}

func (d *DataSecurityAiProposalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSecurityAiProposalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataSecurityAiProposalDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.DataSecurityIssueId.IsNull() && !data.DataSecurityIssueId.IsUnknown() {
        filters["dataSecurityIssueId"] = data.DataSecurityIssueId.ValueString()
        filterNames = append(filterNames, "data_security_issue_id = "+fmt.Sprintf("%q", data.DataSecurityIssueId.ValueString()))
    }
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.DataSecurityIssueGroupId.IsNull() && !data.DataSecurityIssueGroupId.IsUnknown() {
        filters["dataSecurityIssueGroupId"] = data.DataSecurityIssueGroupId.ValueString()
        filterNames = append(filterNames, "data_security_issue_group_id = "+fmt.Sprintf("%q", data.DataSecurityIssueGroupId.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.Work.IsNull() && !data.Work.IsUnknown() {
        filters["work"] = data.Work.ValueString()
        filterNames = append(filterNames, "work = "+fmt.Sprintf("%q", data.Work.ValueString()))
    }
    if !data.Safety.IsNull() && !data.Safety.IsUnknown() {
        filters["safety"] = data.Safety.ValueString()
        filterNames = append(filterNames, "safety = "+fmt.Sprintf("%q", data.Safety.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.ChangeMarkdown.IsNull() && !data.ChangeMarkdown.IsUnknown() {
        filters["changeMarkdown"] = data.ChangeMarkdown.ValueString()
        filterNames = append(filterNames, "change_markdown = "+fmt.Sprintf("%q", data.ChangeMarkdown.ValueString()))
    }
    if !data.ReasonMarkdown.IsNull() && !data.ReasonMarkdown.IsUnknown() {
        filters["reasonMarkdown"] = data.ReasonMarkdown.ValueString()
        filterNames = append(filterNames, "reason_markdown = "+fmt.Sprintf("%q", data.ReasonMarkdown.ValueString()))
    }
    if !data.ProofMarkdown.IsNull() && !data.ProofMarkdown.IsUnknown() {
        filters["proofMarkdown"] = data.ProofMarkdown.ValueString()
        filterNames = append(filterNames, "proof_markdown = "+fmt.Sprintf("%q", data.ProofMarkdown.ValueString()))
    }
    if !data.UndoMarkdown.IsNull() && !data.UndoMarkdown.IsUnknown() {
        filters["undoMarkdown"] = data.UndoMarkdown.ValueString()
        filterNames = append(filterNames, "undo_markdown = "+fmt.Sprintf("%q", data.UndoMarkdown.ValueString()))
    }
    if !data.PlanKey.IsNull() && !data.PlanKey.IsUnknown() {
        filters["planKey"] = data.PlanKey.ValueString()
        filterNames = append(filterNames, "plan_key = "+fmt.Sprintf("%q", data.PlanKey.ValueString()))
    }
    if !data.PlanFingerprint.IsNull() && !data.PlanFingerprint.IsUnknown() {
        filters["planFingerprint"] = data.PlanFingerprint.ValueString()
        filterNames = append(filterNames, "plan_fingerprint = "+fmt.Sprintf("%q", data.PlanFingerprint.ValueString()))
    }
    if !data.RiskScore.IsNull() && !data.RiskScore.IsUnknown() {
        filters["riskScore"] = LookupNumber(data.RiskScore)
        filterNames = append(filterNames, "risk_score = "+data.RiskScore.ValueBigFloat().String())
    }
    if !data.Autonomy.IsNull() && !data.Autonomy.IsUnknown() {
        filters["autonomy"] = data.Autonomy.ValueString()
        filterNames = append(filterNames, "autonomy = "+fmt.Sprintf("%q", data.Autonomy.ValueString()))
    }
    if !data.AskReason.IsNull() && !data.AskReason.IsUnknown() {
        filters["askReason"] = data.AskReason.ValueString()
        filterNames = append(filterNames, "ask_reason = "+fmt.Sprintf("%q", data.AskReason.ValueString()))
    }
    if !data.AutoRemediationRuleId.IsNull() && !data.AutoRemediationRuleId.IsUnknown() {
        filters["autoRemediationRuleId"] = data.AutoRemediationRuleId.ValueString()
        filterNames = append(filterNames, "auto_remediation_rule_id = "+fmt.Sprintf("%q", data.AutoRemediationRuleId.ValueString()))
    }
    if !data.AutoRemediationRuleName.IsNull() && !data.AutoRemediationRuleName.IsUnknown() {
        filters["autoRemediationRuleName"] = data.AutoRemediationRuleName.ValueString()
        filterNames = append(filterNames, "auto_remediation_rule_name = "+fmt.Sprintf("%q", data.AutoRemediationRuleName.ValueString()))
    }
    if !data.AiRunId.IsNull() && !data.AiRunId.IsUnknown() {
        filters["aiRunId"] = data.AiRunId.ValueString()
        filterNames = append(filterNames, "ai_run_id = "+fmt.Sprintf("%q", data.AiRunId.ValueString()))
    }
    if !data.FirstApprovedByUserId.IsNull() && !data.FirstApprovedByUserId.IsUnknown() {
        filters["firstApprovedByUserId"] = data.FirstApprovedByUserId.ValueString()
        filterNames = append(filterNames, "first_approved_by_user_id = "+fmt.Sprintf("%q", data.FirstApprovedByUserId.ValueString()))
    }
    if !data.FirstApprovedPlanFingerprint.IsNull() && !data.FirstApprovedPlanFingerprint.IsUnknown() {
        filters["firstApprovedPlanFingerprint"] = data.FirstApprovedPlanFingerprint.ValueString()
        filterNames = append(filterNames, "first_approved_plan_fingerprint = "+fmt.Sprintf("%q", data.FirstApprovedPlanFingerprint.ValueString()))
    }
    if !data.ApprovedByUserId.IsNull() && !data.ApprovedByUserId.IsUnknown() {
        filters["approvedByUserId"] = data.ApprovedByUserId.ValueString()
        filterNames = append(filterNames, "approved_by_user_id = "+fmt.Sprintf("%q", data.ApprovedByUserId.ValueString()))
    }
    if !data.DismissedByUserId.IsNull() && !data.DismissedByUserId.IsUnknown() {
        filters["dismissedByUserId"] = data.DismissedByUserId.ValueString()
        filterNames = append(filterNames, "dismissed_by_user_id = "+fmt.Sprintf("%q", data.DismissedByUserId.ValueString()))
    }
    if !data.DismissReason.IsNull() && !data.DismissReason.IsUnknown() {
        filters["dismissReason"] = data.DismissReason.ValueString()
        filterNames = append(filterNames, "dismiss_reason = "+fmt.Sprintf("%q", data.DismissReason.ValueString()))
    }
    if !data.UndoneByUserId.IsNull() && !data.UndoneByUserId.IsUnknown() {
        filters["undoneByUserId"] = data.UndoneByUserId.ValueString()
        filterNames = append(filterNames, "undone_by_user_id = "+fmt.Sprintf("%q", data.UndoneByUserId.ValueString()))
    }
    if !data.FailureMessage.IsNull() && !data.FailureMessage.IsUnknown() {
        filters["failureMessage"] = data.FailureMessage.ValueString()
        filterNames = append(filterNames, "failure_message = "+fmt.Sprintf("%q", data.FailureMessage.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data security ai proposal up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data security ai proposal up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "dataSecurityIssueId": true,
        "dataStoreId": true,
        "dataSecurityIssueGroupId": true,
        "kind": true,
        "work": true,
        "safety": true,
        "status": true,
        "title": true,
        "changeMarkdown": true,
        "reasonMarkdown": true,
        "proofMarkdown": true,
        "undoMarkdown": true,
        "plan": true,
        "planKey": true,
        "planFingerprint": true,
        "result": true,
        "riskScore": true,
        "autonomy": true,
        "askReason": true,
        "autoRemediationRuleId": true,
        "autoRemediationRuleName": true,
        "aiRunId": true,
        "firstApprovedByUserId": true,
        "firstApprovedAt": true,
        "firstApprovedPlanFingerprint": true,
        "approvedByUserId": true,
        "approvedAt": true,
        "appliedAt": true,
        "dismissedByUserId": true,
        "dismissedAt": true,
        "dismissReason": true,
        "undoneByUserId": true,
        "undoneAt": true,
        "expiredAt": true,
        "failedAt": true,
        "failureMessage": true,
        "verifiedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-ai-proposal/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_ai_proposal, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security ai proposal found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_security_ai_proposal: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-ai-proposal/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_security_ai_proposal, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_security_ai_proposal: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security ai proposal matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data security ai proposal matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_security_ai_proposal.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["dataSecurityIssueId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityIssueId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataSecurityIssueId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataSecurityIssueId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataSecurityIssueId = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityIssueId = types.StringNull()
        }
    } else if arr, ok := item["dataSecurityIssueId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataSecurityIssueId = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityIssueId = types.StringNull()
        }
    } else if val, ok := item["dataSecurityIssueId"].(string); ok {
        data.DataSecurityIssueId = types.StringValue(val)
    } else {
        data.DataSecurityIssueId = types.StringNull()
    }
    if obj, ok := item["dataStoreId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataStoreId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreId = types.StringNull()
        }
    } else if arr, ok := item["dataStoreId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreId = types.StringNull()
        }
    } else if val, ok := item["dataStoreId"].(string); ok {
        data.DataStoreId = types.StringValue(val)
    } else {
        data.DataStoreId = types.StringNull()
    }
    if obj, ok := item["dataSecurityIssueGroupId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSecurityIssueGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataSecurityIssueGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataSecurityIssueGroupId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataSecurityIssueGroupId = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityIssueGroupId = types.StringNull()
        }
    } else if arr, ok := item["dataSecurityIssueGroupId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataSecurityIssueGroupId = types.StringValue(string(jsonBytes))
        } else {
            data.DataSecurityIssueGroupId = types.StringNull()
        }
    } else if val, ok := item["dataSecurityIssueGroupId"].(string); ok {
        data.DataSecurityIssueGroupId = types.StringValue(val)
    } else {
        data.DataSecurityIssueGroupId = types.StringNull()
    }
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if arr, ok := item["kind"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := item["work"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Work = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Work = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Work = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Work = types.StringValue(string(jsonBytes))
        } else {
            data.Work = types.StringNull()
        }
    } else if arr, ok := item["work"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Work = types.StringValue(string(jsonBytes))
        } else {
            data.Work = types.StringNull()
        }
    } else if val, ok := item["work"].(string); ok {
        data.Work = types.StringValue(val)
    } else {
        data.Work = types.StringNull()
    }
    if obj, ok := item["safety"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Safety = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Safety = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Safety = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Safety = types.StringValue(string(jsonBytes))
        } else {
            data.Safety = types.StringNull()
        }
    } else if arr, ok := item["safety"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Safety = types.StringValue(string(jsonBytes))
        } else {
            data.Safety = types.StringNull()
        }
    } else if val, ok := item["safety"].(string); ok {
        data.Safety = types.StringValue(val)
    } else {
        data.Safety = types.StringNull()
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
    } else if arr, ok := item["status"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if arr, ok := item["title"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := item["changeMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChangeMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChangeMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChangeMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMarkdown = types.StringNull()
        }
    } else if arr, ok := item["changeMarkdown"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ChangeMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMarkdown = types.StringNull()
        }
    } else if val, ok := item["changeMarkdown"].(string); ok {
        data.ChangeMarkdown = types.StringValue(val)
    } else {
        data.ChangeMarkdown = types.StringNull()
    }
    if obj, ok := item["reasonMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReasonMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReasonMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReasonMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReasonMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ReasonMarkdown = types.StringNull()
        }
    } else if arr, ok := item["reasonMarkdown"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ReasonMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ReasonMarkdown = types.StringNull()
        }
    } else if val, ok := item["reasonMarkdown"].(string); ok {
        data.ReasonMarkdown = types.StringValue(val)
    } else {
        data.ReasonMarkdown = types.StringNull()
    }
    if obj, ok := item["proofMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProofMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProofMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProofMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProofMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ProofMarkdown = types.StringNull()
        }
    } else if arr, ok := item["proofMarkdown"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProofMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.ProofMarkdown = types.StringNull()
        }
    } else if val, ok := item["proofMarkdown"].(string); ok {
        data.ProofMarkdown = types.StringValue(val)
    } else {
        data.ProofMarkdown = types.StringNull()
    }
    if obj, ok := item["undoMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UndoMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UndoMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UndoMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UndoMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.UndoMarkdown = types.StringNull()
        }
    } else if arr, ok := item["undoMarkdown"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UndoMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.UndoMarkdown = types.StringNull()
        }
    } else if val, ok := item["undoMarkdown"].(string); ok {
        data.UndoMarkdown = types.StringValue(val)
    } else {
        data.UndoMarkdown = types.StringNull()
    }
    if obj, ok := item["plan"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Plan = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Plan = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Plan = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Plan = types.StringValue(string(jsonBytes))
        } else {
            data.Plan = types.StringNull()
        }
    } else if arr, ok := item["plan"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Plan = types.StringValue(string(jsonBytes))
        } else {
            data.Plan = types.StringNull()
        }
    } else if val, ok := item["plan"].(string); ok {
        data.Plan = types.StringValue(val)
    } else {
        data.Plan = types.StringNull()
    }
    if obj, ok := item["planKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlanKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PlanKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PlanKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PlanKey = types.StringValue(string(jsonBytes))
        } else {
            data.PlanKey = types.StringNull()
        }
    } else if arr, ok := item["planKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PlanKey = types.StringValue(string(jsonBytes))
        } else {
            data.PlanKey = types.StringNull()
        }
    } else if val, ok := item["planKey"].(string); ok {
        data.PlanKey = types.StringValue(val)
    } else {
        data.PlanKey = types.StringNull()
    }
    if obj, ok := item["planFingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlanFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PlanFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PlanFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PlanFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.PlanFingerprint = types.StringNull()
        }
    } else if arr, ok := item["planFingerprint"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PlanFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.PlanFingerprint = types.StringNull()
        }
    } else if val, ok := item["planFingerprint"].(string); ok {
        data.PlanFingerprint = types.StringValue(val)
    } else {
        data.PlanFingerprint = types.StringNull()
    }
    if obj, ok := item["result"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Result = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Result = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Result = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Result = types.StringValue(string(jsonBytes))
        } else {
            data.Result = types.StringNull()
        }
    } else if arr, ok := item["result"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Result = types.StringValue(string(jsonBytes))
        } else {
            data.Result = types.StringNull()
        }
    } else if val, ok := item["result"].(string); ok {
        data.Result = types.StringValue(val)
    } else {
        data.Result = types.StringNull()
    }
    data.RiskScore = NumberFromAPI(data.RiskScore, item["riskScore"])
    if obj, ok := item["autonomy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Autonomy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Autonomy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Autonomy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Autonomy = types.StringValue(string(jsonBytes))
        } else {
            data.Autonomy = types.StringNull()
        }
    } else if arr, ok := item["autonomy"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Autonomy = types.StringValue(string(jsonBytes))
        } else {
            data.Autonomy = types.StringNull()
        }
    } else if val, ok := item["autonomy"].(string); ok {
        data.Autonomy = types.StringValue(val)
    } else {
        data.Autonomy = types.StringNull()
    }
    if obj, ok := item["askReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AskReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AskReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AskReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AskReason = types.StringValue(string(jsonBytes))
        } else {
            data.AskReason = types.StringNull()
        }
    } else if arr, ok := item["askReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AskReason = types.StringValue(string(jsonBytes))
        } else {
            data.AskReason = types.StringNull()
        }
    } else if val, ok := item["askReason"].(string); ok {
        data.AskReason = types.StringValue(val)
    } else {
        data.AskReason = types.StringNull()
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
    } else if arr, ok := item["autoRemediationRuleId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AutoRemediationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.AutoRemediationRuleId = types.StringNull()
        }
    } else if val, ok := item["autoRemediationRuleId"].(string); ok {
        data.AutoRemediationRuleId = types.StringValue(val)
    } else {
        data.AutoRemediationRuleId = types.StringNull()
    }
    if obj, ok := item["autoRemediationRuleName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoRemediationRuleName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoRemediationRuleName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoRemediationRuleName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoRemediationRuleName = types.StringValue(string(jsonBytes))
        } else {
            data.AutoRemediationRuleName = types.StringNull()
        }
    } else if arr, ok := item["autoRemediationRuleName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AutoRemediationRuleName = types.StringValue(string(jsonBytes))
        } else {
            data.AutoRemediationRuleName = types.StringNull()
        }
    } else if val, ok := item["autoRemediationRuleName"].(string); ok {
        data.AutoRemediationRuleName = types.StringValue(val)
    } else {
        data.AutoRemediationRuleName = types.StringNull()
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
    } else if arr, ok := item["aiRunId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.AiRunId = types.StringNull()
        }
    } else if val, ok := item["aiRunId"].(string); ok {
        data.AiRunId = types.StringValue(val)
    } else {
        data.AiRunId = types.StringNull()
    }
    if obj, ok := item["firstApprovedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstApprovedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstApprovedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstApprovedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstApprovedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.FirstApprovedByUserId = types.StringNull()
        }
    } else if arr, ok := item["firstApprovedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstApprovedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.FirstApprovedByUserId = types.StringNull()
        }
    } else if val, ok := item["firstApprovedByUserId"].(string); ok {
        data.FirstApprovedByUserId = types.StringValue(val)
    } else {
        data.FirstApprovedByUserId = types.StringNull()
    }
    if obj, ok := item["firstApprovedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstApprovedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstApprovedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstApprovedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstApprovedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstApprovedAt = types.StringNull()
        }
    } else if arr, ok := item["firstApprovedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstApprovedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstApprovedAt = types.StringNull()
        }
    } else if val, ok := item["firstApprovedAt"].(string); ok {
        data.FirstApprovedAt = types.StringValue(val)
    } else {
        data.FirstApprovedAt = types.StringNull()
    }
    if obj, ok := item["firstApprovedPlanFingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstApprovedPlanFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstApprovedPlanFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstApprovedPlanFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstApprovedPlanFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.FirstApprovedPlanFingerprint = types.StringNull()
        }
    } else if arr, ok := item["firstApprovedPlanFingerprint"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstApprovedPlanFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.FirstApprovedPlanFingerprint = types.StringNull()
        }
    } else if val, ok := item["firstApprovedPlanFingerprint"].(string); ok {
        data.FirstApprovedPlanFingerprint = types.StringValue(val)
    } else {
        data.FirstApprovedPlanFingerprint = types.StringNull()
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
    } else if arr, ok := item["approvedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["approvedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ApprovedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ApprovedAt = types.StringNull()
        }
    } else if val, ok := item["approvedAt"].(string); ok {
        data.ApprovedAt = types.StringValue(val)
    } else {
        data.ApprovedAt = types.StringNull()
    }
    if obj, ok := item["appliedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppliedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AppliedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AppliedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AppliedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AppliedAt = types.StringNull()
        }
    } else if arr, ok := item["appliedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AppliedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AppliedAt = types.StringNull()
        }
    } else if val, ok := item["appliedAt"].(string); ok {
        data.AppliedAt = types.StringValue(val)
    } else {
        data.AppliedAt = types.StringNull()
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
    } else if arr, ok := item["dismissedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["dismissedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DismissedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DismissedAt = types.StringNull()
        }
    } else if val, ok := item["dismissedAt"].(string); ok {
        data.DismissedAt = types.StringValue(val)
    } else {
        data.DismissedAt = types.StringNull()
    }
    if obj, ok := item["dismissReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DismissReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DismissReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DismissReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DismissReason = types.StringValue(string(jsonBytes))
        } else {
            data.DismissReason = types.StringNull()
        }
    } else if arr, ok := item["dismissReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DismissReason = types.StringValue(string(jsonBytes))
        } else {
            data.DismissReason = types.StringNull()
        }
    } else if val, ok := item["dismissReason"].(string); ok {
        data.DismissReason = types.StringValue(val)
    } else {
        data.DismissReason = types.StringNull()
    }
    if obj, ok := item["undoneByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UndoneByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UndoneByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UndoneByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UndoneByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.UndoneByUserId = types.StringNull()
        }
    } else if arr, ok := item["undoneByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UndoneByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.UndoneByUserId = types.StringNull()
        }
    } else if val, ok := item["undoneByUserId"].(string); ok {
        data.UndoneByUserId = types.StringValue(val)
    } else {
        data.UndoneByUserId = types.StringNull()
    }
    if obj, ok := item["undoneAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UndoneAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UndoneAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UndoneAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UndoneAt = types.StringValue(string(jsonBytes))
        } else {
            data.UndoneAt = types.StringNull()
        }
    } else if arr, ok := item["undoneAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UndoneAt = types.StringValue(string(jsonBytes))
        } else {
            data.UndoneAt = types.StringNull()
        }
    } else if val, ok := item["undoneAt"].(string); ok {
        data.UndoneAt = types.StringValue(val)
    } else {
        data.UndoneAt = types.StringNull()
    }
    if obj, ok := item["expiredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExpiredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExpiredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExpiredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExpiredAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiredAt = types.StringNull()
        }
    } else if arr, ok := item["expiredAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ExpiredAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiredAt = types.StringNull()
        }
    } else if val, ok := item["expiredAt"].(string); ok {
        data.ExpiredAt = types.StringValue(val)
    } else {
        data.ExpiredAt = types.StringNull()
    }
    if obj, ok := item["failedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FailedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FailedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FailedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FailedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FailedAt = types.StringNull()
        }
    } else if arr, ok := item["failedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FailedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FailedAt = types.StringNull()
        }
    } else if val, ok := item["failedAt"].(string); ok {
        data.FailedAt = types.StringValue(val)
    } else {
        data.FailedAt = types.StringNull()
    }
    if obj, ok := item["failureMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FailureMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FailureMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FailureMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FailureMessage = types.StringValue(string(jsonBytes))
        } else {
            data.FailureMessage = types.StringNull()
        }
    } else if arr, ok := item["failureMessage"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FailureMessage = types.StringValue(string(jsonBytes))
        } else {
            data.FailureMessage = types.StringNull()
        }
    } else if val, ok := item["failureMessage"].(string); ok {
        data.FailureMessage = types.StringValue(val)
    } else {
        data.FailureMessage = types.StringNull()
    }
    if obj, ok := item["verifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VerifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VerifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VerifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VerifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.VerifiedAt = types.StringNull()
        }
    } else if arr, ok := item["verifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.VerifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.VerifiedAt = types.StringNull()
        }
    } else if val, ok := item["verifiedAt"].(string); ok {
        data.VerifiedAt = types.StringValue(val)
    } else {
        data.VerifiedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
