package shard23

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
var _ datasource.DataSource = &SecurityCaseAiActionDataSource{}

func NewSecurityCaseAiActionDataSource() datasource.DataSource {
    return &SecurityCaseAiActionDataSource{}
}

// SecurityCaseAiActionDataSource defines the data source implementation.
type SecurityCaseAiActionDataSource struct {
    client *Client
}

// SecurityCaseAiActionDataSourceModel describes the data source data model.
type SecurityCaseAiActionDataSourceModel struct {
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
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    Kind types.String `tfsdk:"kind"`
    Area types.String `tfsdk:"area"`
    Mode types.String `tfsdk:"mode"`
    Status types.String `tfsdk:"status"`
    Reason types.String `tfsdk:"reason"`
    HandOffReason types.String `tfsdk:"hand_off_reason"`
    StepKey types.String `tfsdk:"step_key"`
    Target types.String `tfsdk:"target"`
    SecurityPlaybookRunId types.String `tfsdk:"security_playbook_run_id"`
    DetectionRuleId types.String `tfsdk:"detection_rule_id"`
    Evidence types.String `tfsdk:"evidence"`
    Report types.String `tfsdk:"report"`
    AskedAt types.String `tfsdk:"asked_at"`
    DecidedAt types.String `tfsdk:"decided_at"`
    DecidedByUserId types.String `tfsdk:"decided_by_user_id"`
    DoneAt types.String `tfsdk:"done_at"`
    UndoneAt types.String `tfsdk:"undone_at"`
    UndoneByUserId types.String `tfsdk:"undone_by_user_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityCaseAiActionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_case_ai_action"
}

func (d *SecurityCaseAiActionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What OneUptime AI did, asked a person about or left to a person on a security case it works, and why. Look up an existing security case ai action by `id`, or by any of its other arguments (`ai_run_id`, `area`, `created_by_actor_id`, ...): each one set must match, and exactly one security case ai action may match them all.",

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
                MarkdownDescription: "ID of the project this action belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case OneUptime AI worked. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of OneUptime AI's run that did, asked or left this.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What OneUptime AI did, or would have done: TakeCase, SetPriority, MarkStepDone, MarkContained, BlockAddress, ContainAccount, ContainHost, DeclareIncident, RequestCodeFix, CloseCase, LiftBlock, OfferRuleException, StepForPeople or HandToPeople.",
                Optional: true,
                Computed: true,
            },
            "area": schema.StringAttribute{
                MarkdownDescription: "The kind of action the project's setting for it covers: Case, Edge, Accounts, Hosts, Code or Close. Empty for a step no setting covers.",
                Optional: true,
                Computed: true,
            },
            "mode": schema.StringAttribute{
                MarkdownDescription: "What the project's setting for this kind of action said when OneUptime AI acted: Off, Ask or OnItsOwn.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "What became of it: Doing, Done, Asked (waiting for a person), Declined, Left (OneUptime AI would not), Failed or Undone.",
                Optional: true,
                Computed: true,
            },
            "reason": schema.StringAttribute{
                MarkdownDescription: "Why it is as it is, as a code the case's AI card words.",
                Optional: true,
                Computed: true,
            },
            "hand_off_reason": schema.StringAttribute{
                MarkdownDescription: "For a hand-off to people: why OneUptime AI handed the case over.",
                Optional: true,
                Computed: true,
            },
            "step_key": schema.StringAttribute{
                MarkdownDescription: "The step of the case's Response Plan it is about: what the step is, and what it is about.",
                Optional: true,
                Computed: true,
            },
            "target": schema.StringAttribute{
                MarkdownDescription: "What it acts on, as the case names it: an address, a user, a host, a route, or the status it sets.",
                Optional: true,
                Computed: true,
            },
            "security_playbook_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the playbook engine's run that carries it out, or waits for its approval.",
                Optional: true,
                Computed: true,
            },
            "detection_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the detection rule that opened the case: the rule an exception is offered for, and what makes cases alike.",
                Optional: true,
                Computed: true,
            },
            "evidence": schema.StringAttribute{
                MarkdownDescription: "What OneUptime AI read that made it do, ask or leave this: facts, checks and links. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "report": schema.StringAttribute{
                MarkdownDescription: "For a close: OneUptime AI's closing report, in Markdown - what happened, the impact, what was done, what was fixed and what to watch.",
                Optional: true,
                Computed: true,
            },
            "asked_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI asked a person to approve it.",
                Computed: true,
            },
            "decided_at": schema.StringAttribute{
                MarkdownDescription: "When a person approved or declined it.",
                Computed: true,
            },
            "decided_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who approved or declined it.",
                Optional: true,
                Computed: true,
            },
            "done_at": schema.StringAttribute{
                MarkdownDescription: "When it was done.",
                Computed: true,
            },
            "undone_at": schema.StringAttribute{
                MarkdownDescription: "When a person undid it.",
                Computed: true,
            },
            "undone_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who undid it.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "Unused: OneUptime AI writes every action. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *SecurityCaseAiActionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityCaseAiActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityCaseAiActionDataSourceModel

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
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.AiRunId.IsNull() && !data.AiRunId.IsUnknown() {
        filters["aiRunId"] = data.AiRunId.ValueString()
        filterNames = append(filterNames, "ai_run_id = "+fmt.Sprintf("%q", data.AiRunId.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.Area.IsNull() && !data.Area.IsUnknown() {
        filters["area"] = data.Area.ValueString()
        filterNames = append(filterNames, "area = "+fmt.Sprintf("%q", data.Area.ValueString()))
    }
    if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
        filters["mode"] = data.Mode.ValueString()
        filterNames = append(filterNames, "mode = "+fmt.Sprintf("%q", data.Mode.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.Reason.IsNull() && !data.Reason.IsUnknown() {
        filters["reason"] = data.Reason.ValueString()
        filterNames = append(filterNames, "reason = "+fmt.Sprintf("%q", data.Reason.ValueString()))
    }
    if !data.HandOffReason.IsNull() && !data.HandOffReason.IsUnknown() {
        filters["handOffReason"] = data.HandOffReason.ValueString()
        filterNames = append(filterNames, "hand_off_reason = "+fmt.Sprintf("%q", data.HandOffReason.ValueString()))
    }
    if !data.StepKey.IsNull() && !data.StepKey.IsUnknown() {
        filters["stepKey"] = data.StepKey.ValueString()
        filterNames = append(filterNames, "step_key = "+fmt.Sprintf("%q", data.StepKey.ValueString()))
    }
    if !data.Target.IsNull() && !data.Target.IsUnknown() {
        filters["target"] = data.Target.ValueString()
        filterNames = append(filterNames, "target = "+fmt.Sprintf("%q", data.Target.ValueString()))
    }
    if !data.SecurityPlaybookRunId.IsNull() && !data.SecurityPlaybookRunId.IsUnknown() {
        filters["securityPlaybookRunId"] = data.SecurityPlaybookRunId.ValueString()
        filterNames = append(filterNames, "security_playbook_run_id = "+fmt.Sprintf("%q", data.SecurityPlaybookRunId.ValueString()))
    }
    if !data.DetectionRuleId.IsNull() && !data.DetectionRuleId.IsUnknown() {
        filters["detectionRuleId"] = data.DetectionRuleId.ValueString()
        filterNames = append(filterNames, "detection_rule_id = "+fmt.Sprintf("%q", data.DetectionRuleId.ValueString()))
    }
    if !data.Report.IsNull() && !data.Report.IsUnknown() {
        filters["report"] = data.Report.ValueString()
        filterNames = append(filterNames, "report = "+fmt.Sprintf("%q", data.Report.ValueString()))
    }
    if !data.DecidedByUserId.IsNull() && !data.DecidedByUserId.IsUnknown() {
        filters["decidedByUserId"] = data.DecidedByUserId.ValueString()
        filterNames = append(filterNames, "decided_by_user_id = "+fmt.Sprintf("%q", data.DecidedByUserId.ValueString()))
    }
    if !data.UndoneByUserId.IsNull() && !data.UndoneByUserId.IsUnknown() {
        filters["undoneByUserId"] = data.UndoneByUserId.ValueString()
        filterNames = append(filterNames, "undone_by_user_id = "+fmt.Sprintf("%q", data.UndoneByUserId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security case ai action up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security case ai action up by.",
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
        "securityCaseId": true,
        "aiRunId": true,
        "kind": true,
        "area": true,
        "mode": true,
        "status": true,
        "reason": true,
        "handOffReason": true,
        "stepKey": true,
        "target": true,
        "securityPlaybookRunId": true,
        "detectionRuleId": true,
        "evidence": true,
        "report": true,
        "askedAt": true,
        "decidedAt": true,
        "decidedByUserId": true,
        "doneAt": true,
        "undoneAt": true,
        "undoneByUserId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-case-ai-action/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case_ai_action, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case ai action found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_case_ai_action: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-case-ai-action/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_case_ai_action, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_case_ai_action: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case ai action matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security case ai action matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_case_ai_action.")
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
    if obj, ok := item["securityCaseId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := item["securityCaseId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := item["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
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
    if obj, ok := item["area"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Area = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Area = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Area = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Area = types.StringValue(string(jsonBytes))
        } else {
            data.Area = types.StringNull()
        }
    } else if arr, ok := item["area"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Area = types.StringValue(string(jsonBytes))
        } else {
            data.Area = types.StringNull()
        }
    } else if val, ok := item["area"].(string); ok {
        data.Area = types.StringValue(val)
    } else {
        data.Area = types.StringNull()
    }
    if obj, ok := item["mode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Mode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Mode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Mode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Mode = types.StringValue(string(jsonBytes))
        } else {
            data.Mode = types.StringNull()
        }
    } else if arr, ok := item["mode"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Mode = types.StringValue(string(jsonBytes))
        } else {
            data.Mode = types.StringNull()
        }
    } else if val, ok := item["mode"].(string); ok {
        data.Mode = types.StringValue(val)
    } else {
        data.Mode = types.StringNull()
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
    if obj, ok := item["reason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if arr, ok := item["reason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := item["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := item["handOffReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HandOffReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HandOffReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HandOffReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HandOffReason = types.StringValue(string(jsonBytes))
        } else {
            data.HandOffReason = types.StringNull()
        }
    } else if arr, ok := item["handOffReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.HandOffReason = types.StringValue(string(jsonBytes))
        } else {
            data.HandOffReason = types.StringNull()
        }
    } else if val, ok := item["handOffReason"].(string); ok {
        data.HandOffReason = types.StringValue(val)
    } else {
        data.HandOffReason = types.StringNull()
    }
    if obj, ok := item["stepKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StepKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StepKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StepKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StepKey = types.StringValue(string(jsonBytes))
        } else {
            data.StepKey = types.StringNull()
        }
    } else if arr, ok := item["stepKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StepKey = types.StringValue(string(jsonBytes))
        } else {
            data.StepKey = types.StringNull()
        }
    } else if val, ok := item["stepKey"].(string); ok {
        data.StepKey = types.StringValue(val)
    } else {
        data.StepKey = types.StringNull()
    }
    if obj, ok := item["target"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Target = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Target = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Target = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Target = types.StringValue(string(jsonBytes))
        } else {
            data.Target = types.StringNull()
        }
    } else if arr, ok := item["target"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Target = types.StringValue(string(jsonBytes))
        } else {
            data.Target = types.StringNull()
        }
    } else if val, ok := item["target"].(string); ok {
        data.Target = types.StringValue(val)
    } else {
        data.Target = types.StringNull()
    }
    if obj, ok := item["securityPlaybookRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if arr, ok := item["securityPlaybookRunId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if val, ok := item["securityPlaybookRunId"].(string); ok {
        data.SecurityPlaybookRunId = types.StringValue(val)
    } else {
        data.SecurityPlaybookRunId = types.StringNull()
    }
    if obj, ok := item["detectionRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DetectionRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionRuleId = types.StringNull()
        }
    } else if arr, ok := item["detectionRuleId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DetectionRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionRuleId = types.StringNull()
        }
    } else if val, ok := item["detectionRuleId"].(string); ok {
        data.DetectionRuleId = types.StringValue(val)
    } else {
        data.DetectionRuleId = types.StringNull()
    }
    if obj, ok := item["evidence"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Evidence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Evidence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Evidence = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Evidence = types.StringValue(string(jsonBytes))
        } else {
            data.Evidence = types.StringNull()
        }
    } else if arr, ok := item["evidence"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Evidence = types.StringValue(string(jsonBytes))
        } else {
            data.Evidence = types.StringNull()
        }
    } else if val, ok := item["evidence"].(string); ok {
        data.Evidence = types.StringValue(val)
    } else {
        data.Evidence = types.StringNull()
    }
    if obj, ok := item["report"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Report = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Report = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Report = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Report = types.StringValue(string(jsonBytes))
        } else {
            data.Report = types.StringNull()
        }
    } else if arr, ok := item["report"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Report = types.StringValue(string(jsonBytes))
        } else {
            data.Report = types.StringNull()
        }
    } else if val, ok := item["report"].(string); ok {
        data.Report = types.StringValue(val)
    } else {
        data.Report = types.StringNull()
    }
    if obj, ok := item["askedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AskedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AskedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AskedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AskedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AskedAt = types.StringNull()
        }
    } else if arr, ok := item["askedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AskedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AskedAt = types.StringNull()
        }
    } else if val, ok := item["askedAt"].(string); ok {
        data.AskedAt = types.StringValue(val)
    } else {
        data.AskedAt = types.StringNull()
    }
    if obj, ok := item["decidedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DecidedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DecidedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DecidedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DecidedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DecidedAt = types.StringNull()
        }
    } else if arr, ok := item["decidedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DecidedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DecidedAt = types.StringNull()
        }
    } else if val, ok := item["decidedAt"].(string); ok {
        data.DecidedAt = types.StringValue(val)
    } else {
        data.DecidedAt = types.StringNull()
    }
    if obj, ok := item["decidedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DecidedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DecidedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DecidedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DecidedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DecidedByUserId = types.StringNull()
        }
    } else if arr, ok := item["decidedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DecidedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DecidedByUserId = types.StringNull()
        }
    } else if val, ok := item["decidedByUserId"].(string); ok {
        data.DecidedByUserId = types.StringValue(val)
    } else {
        data.DecidedByUserId = types.StringNull()
    }
    if obj, ok := item["doneAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DoneAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DoneAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DoneAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DoneAt = types.StringValue(string(jsonBytes))
        } else {
            data.DoneAt = types.StringNull()
        }
    } else if arr, ok := item["doneAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DoneAt = types.StringValue(string(jsonBytes))
        } else {
            data.DoneAt = types.StringNull()
        }
    } else if val, ok := item["doneAt"].(string); ok {
        data.DoneAt = types.StringValue(val)
    } else {
        data.DoneAt = types.StringNull()
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
