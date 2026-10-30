package provider

import (
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
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SecurityCaseResource{}
var _ resource.ResourceWithImportState = &SecurityCaseResource{}

func NewSecurityCaseResource() resource.Resource {
    return &SecurityCaseResource{}
}

// SecurityCaseResource defines the resource implementation.
type SecurityCaseResource struct {
    client *Client
}

// SecurityCaseResourceModel describes the resource data model.
type SecurityCaseResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    SecurityCaseStatusId types.String `tfsdk:"security_case_status_id"`
    SecurityCasePriorityId types.String `tfsdk:"security_case_priority_id"`
    Verdict types.String `tfsdk:"verdict"`
    ClosingNote types.String `tfsdk:"closing_note"`
    AssignedToUserId types.String `tfsdk:"assigned_to_user_id"`
    AssignedToTeamId types.String `tfsdk:"assigned_to_team_id"`
    Entities JSONSubsetValue `tfsdk:"entities"`
    MitreTactics JSONSubsetValue `tfsdk:"mitre_tactics"`
    MitreTechniques JSONSubsetValue `tfsdk:"mitre_techniques"`
    Labels types.Set `tfsdk:"labels"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CaseNumber types.Number `tfsdk:"case_number"`
    CaseNumberWithPrefix types.String `tfsdk:"case_number_with_prefix"`
    Source types.String `tfsdk:"source"`
    DetectionCount types.Number `tfsdk:"detection_count"`
    FirstDetectedAt RFC3339Value `tfsdk:"first_detected_at"`
    LastDetectedAt RFC3339Value `tfsdk:"last_detected_at"`
    AcknowledgedAt RFC3339Value `tfsdk:"acknowledged_at"`
    ContainedAt RFC3339Value `tfsdk:"contained_at"`
    ClosedAt RFC3339Value `tfsdk:"closed_at"`
    AcknowledgeDueAt RFC3339Value `tfsdk:"acknowledge_due_at"`
    ContainDueAt RFC3339Value `tfsdk:"contain_due_at"`
    AcknowledgeWarnedAt RFC3339Value `tfsdk:"acknowledge_warned_at"`
    AcknowledgeMissedAt RFC3339Value `tfsdk:"acknowledge_missed_at"`
    ContainWarnedAt RFC3339Value `tfsdk:"contain_warned_at"`
    ContainMissedAt RFC3339Value `tfsdk:"contain_missed_at"`
    NextDueAt RFC3339Value `tfsdk:"next_due_at"`
    NextTargetCheckAt RFC3339Value `tfsdk:"next_target_check_at"`
    ClosedByUserId types.String `tfsdk:"closed_by_user_id"`
    GroupingKey types.String `tfsdk:"grouping_key"`
    NotifiedSecurityCasePriorityId types.String `tfsdk:"notified_security_case_priority_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    AiVerdict types.String `tfsdk:"ai_verdict"`
    AiSuggestedSecurityCasePriorityId types.String `tfsdk:"ai_suggested_security_case_priority_id"`
    AiSuggestedPlaybookId types.String `tfsdk:"ai_suggested_playbook_id"`
    AiTriageReport types.String `tfsdk:"ai_triage_report"`
    AiTriagedAt RFC3339Value `tfsdk:"ai_triaged_at"`
    AiTriagePromptInjectionSuspected types.Bool `tfsdk:"ai_triage_prompt_injection_suspected"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *SecurityCaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_case"
}

func (r *SecurityCaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *SecurityCaseResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this case belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "What the case is about, in a line.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What is known about the threat, in Markdown. A case opened by a detection starts with the detection's description.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "security_case_status_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case's status. Optional on create: a new case starts in the project's new status. The ID of a `oneuptime_security_case_status`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "security_case_priority_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case's priority. Optional on create: a case opened by hand starts at the project's medium priority. The ID of a `oneuptime_security_case_priority`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "verdict": schema.StringAttribute{
                MarkdownDescription: "What the case turned out to be, set when it is closed: True Positive, Benign Positive, False Positive or Undetermined.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "closing_note": schema.StringAttribute{
                MarkdownDescription: "Why the case was closed the way it was, in Markdown: what was found and what was done.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "assigned_to_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person working the case. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "assigned_to_team_id": schema.StringAttribute{
                MarkdownDescription: "ID of the team whose queue the case is in. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "entities": schema.StringAttribute{
                MarkdownDescription: "The users, hosts, IP addresses, domains and files the case is about, each as { type, value }. Detections add theirs as they join the case. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "mitre_tactics": schema.StringAttribute{
                MarkdownDescription: "ATT&CK tactic ids (TA0006) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "mitre_techniques": schema.StringAttribute{
                MarkdownDescription: "ATT&CK technique ids (T1110) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to filter and report on cases by. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
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
            "case_number": schema.NumberAttribute{
                MarkdownDescription: "The case's number in this project, counting up from 1 as cases are opened.",
                Computed: true,
            },
            "case_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "The case number with the project's case number prefix in front of it, like CASE-42.",
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "What opened the case: a Detection Rule, Threat Intel, a Connected Product, Entity Risk, or a person (Manual).",
                Computed: true,
            },
            "detection_count": schema.NumberAttribute{
                MarkdownDescription: "How many distinct detections the case holds.",
                Computed: true,
            },
            "first_detected_at": schema.StringAttribute{
                MarkdownDescription: "When the case's earliest detection fired.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_detected_at": schema.StringAttribute{
                MarkdownDescription: "When a detection of the case last fired.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "acknowledged_at": schema.StringAttribute{
                MarkdownDescription: "When somebody first took the case on: assigned it to a person, or moved it out of New. The time to acknowledge is measured to here.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "contained_at": schema.StringAttribute{
                MarkdownDescription: "When the case first reached Contained. The time to contain is measured to here.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "closed_at": schema.StringAttribute{
                MarkdownDescription: "When the case was closed. Cleared if it is reopened. The time to close is measured to here.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "acknowledge_due_at": schema.StringAttribute{
                MarkdownDescription: "When the case must be acknowledged by, from its priority's Acknowledge Within target. Empty when its priority has no target, or when it was acknowledged as it opened.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "contain_due_at": schema.StringAttribute{
                MarkdownDescription: "When the case must be contained or closed by, from its priority's Contain Within target. Empty when its priority has no target.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "acknowledge_warned_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime told the people on the case that it was about to miss its acknowledge target.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "acknowledge_missed_at": schema.StringAttribute{
                MarkdownDescription: "When the case passed its acknowledge due time without being acknowledged. Empty while it has not.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "contain_warned_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime told the people on the case that it was about to miss its contain target.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "contain_missed_at": schema.StringAttribute{
                MarkdownDescription: "When the case passed its contain due time without being contained or closed. Empty while it has not.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "next_due_at": schema.StringAttribute{
                MarkdownDescription: "The earliest due time of a response target the case has not met yet, a missed one included. Empty when every target is met or the case has none.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "next_target_check_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime next checks the case's response targets: when one is about to be missed, or is missed. Empty when there is nothing left to check.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "closed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who closed the case. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "grouping_key": schema.StringAttribute{
                MarkdownDescription: "What new detections are matched on to join this case while it is open: the entity they are about, or the detection's own fingerprint.",
                Computed: true,
            },
            "notified_security_case_priority_id": schema.StringAttribute{
                MarkdownDescription: "ID of the highest priority people have been told about this case at. They are told again only when it rises above it.",
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the OneUptime incident declared for this case. The ID of a `oneuptime_incident`.",
                Computed: true,
            },
            "ai_verdict": schema.StringAttribute{
                MarkdownDescription: "What OneUptime AI made of the case when it triaged it: Likely Malicious, Likely Benign, Likely False Positive or Inconclusive. A suggestion; a person sets the case's verdict.",
                Computed: true,
            },
            "ai_suggested_security_case_priority_id": schema.StringAttribute{
                MarkdownDescription: "ID of the priority OneUptime AI suggested when it triaged the case. The case's own priority changes only when a person sets it.",
                Computed: true,
            },
            "ai_suggested_playbook_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project's playbook OneUptime AI suggested running when it triaged the case. It runs only when a person runs it.",
                Computed: true,
            },
            "ai_triage_report": schema.StringAttribute{
                MarkdownDescription: "The report OneUptime AI wrote when it last triaged the case: what happened, its read and the evidence for it, and the next steps.",
                Computed: true,
            },
            "ai_triaged_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI last finished triaging the case.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "ai_triage_prompt_injection_suspected": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI found text in the case's data written to steer an AI when it triaged the case. It then leaves the call to a person.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who opened the case by hand. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

func (r *SecurityCaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *SecurityCaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data SecurityCaseResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config SecurityCaseResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    securityCaseRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := securityCaseRequest["data"].(map[string]interface{})

    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        requestDataMap["title"] = data.Title.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.SecurityCaseStatusId.IsNull() && !data.SecurityCaseStatusId.IsUnknown() {
        requestDataMap["securityCaseStatusId"] = data.SecurityCaseStatusId.ValueString()
    }
    if !data.SecurityCasePriorityId.IsNull() && !data.SecurityCasePriorityId.IsUnknown() {
        requestDataMap["securityCasePriorityId"] = data.SecurityCasePriorityId.ValueString()
    }
    if !data.Verdict.IsNull() && !data.Verdict.IsUnknown() {
        requestDataMap["verdict"] = data.Verdict.ValueString()
    }
    if !data.ClosingNote.IsNull() && !data.ClosingNote.IsUnknown() {
        requestDataMap["closingNote"] = data.ClosingNote.ValueString()
    }
    if !data.AssignedToUserId.IsNull() && !data.AssignedToUserId.IsUnknown() {
        requestDataMap["assignedToUserId"] = data.AssignedToUserId.ValueString()
    }
    if !data.AssignedToTeamId.IsNull() && !data.AssignedToTeamId.IsUnknown() {
        requestDataMap["assignedToTeamId"] = data.AssignedToTeamId.ValueString()
    }
    if parsedEntities := r.parseJSONField(data.Entities); parsedEntities != nil {
        requestDataMap["entities"] = parsedEntities
    }
    if parsedMitreTactics := r.parseJSONField(data.MitreTactics); parsedMitreTactics != nil {
        requestDataMap["mitreTactics"] = parsedMitreTactics
    }
    if parsedMitreTechniques := r.parseJSONField(data.MitreTechniques); parsedMitreTechniques != nil {
        requestDataMap["mitreTechniques"] = parsedMitreTechniques
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/security-case", securityCaseRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create security_case, got error: %s", err))
        return
    }

    var securityCaseResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &securityCaseResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create security_case: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := securityCaseResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := securityCaseResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for security_case did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * security_case is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "securityCaseStatusId": true,
        "securityCasePriorityId": true,
        "verdict": true,
        "closingNote": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "entities": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "caseNumber": true,
        "caseNumberWithPrefix": true,
        "source": true,
        "detectionCount": true,
        "firstDetectedAt": true,
        "lastDetectedAt": true,
        "acknowledgedAt": true,
        "containedAt": true,
        "closedAt": true,
        "acknowledgeDueAt": true,
        "containDueAt": true,
        "acknowledgeWarnedAt": true,
        "acknowledgeMissedAt": true,
        "containWarnedAt": true,
        "containMissedAt": true,
        "nextDueAt": true,
        "nextTargetCheckAt": true,
        "closedByUserId": true,
        "groupingKey": true,
        "notifiedSecurityCasePriorityId": true,
        "incidentId": true,
        "aiVerdict": true,
        "aiSuggestedSecurityCasePriorityId": true,
        "aiSuggestedPlaybookId": true,
        "aiTriageReport": true,
        "aiTriagedAt": true,
        "aiTriagePromptInjectionSuspected": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/security-case/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created security_case but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created security_case but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["securityCaseStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCaseStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCaseStatusId"].(string); ok {
        data.SecurityCaseStatusId = types.StringValue(val)
    } else {
        data.SecurityCaseStatusId = types.StringNull()
    }
    if obj, ok := dataMap["securityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCasePriorityId"].(string); ok {
        data.SecurityCasePriorityId = types.StringValue(val)
    } else {
        data.SecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["verdict"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Verdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Verdict = types.StringValue(string(jsonBytes))
            } else {
                data.Verdict = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Verdict = types.StringValue(string(jsonBytes))
            } else {
                data.Verdict = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Verdict = types.StringValue(string(jsonBytes))
        } else {
            data.Verdict = types.StringNull()
        }
    } else if val, ok := dataMap["verdict"].(string); ok {
        data.Verdict = types.StringValue(val)
    } else {
        data.Verdict = types.StringNull()
    }
    if obj, ok := dataMap["closingNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClosingNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClosingNote = types.StringValue(string(jsonBytes))
            } else {
                data.ClosingNote = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClosingNote = types.StringValue(string(jsonBytes))
            } else {
                data.ClosingNote = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClosingNote = types.StringValue(string(jsonBytes))
        } else {
            data.ClosingNote = types.StringNull()
        }
    } else if val, ok := dataMap["closingNote"].(string); ok {
        data.ClosingNote = types.StringValue(val)
    } else {
        data.ClosingNote = types.StringNull()
    }
    if obj, ok := dataMap["assignedToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignedToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignedToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignedToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["assignedToUserId"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    } else {
        data.AssignedToUserId = types.StringNull()
    }
    if obj, ok := dataMap["assignedToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignedToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignedToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignedToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["assignedToTeamId"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    } else {
        data.AssignedToTeamId = types.StringNull()
    }
    if obj, ok := dataMap["entities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Entities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Entities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Entities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Entities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Entities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Entities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["entities"].(string); ok {
        data.Entities = NewJSONSubsetValue(val)
    } else {
        data.Entities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["mitreTactics"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTactics = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MitreTactics = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MitreTactics = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["mitreTactics"].(string); ok {
        data.MitreTactics = NewJSONSubsetValue(val)
    } else {
        data.MitreTactics = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["mitreTechniques"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTechniques = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MitreTechniques = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MitreTechniques = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["mitreTechniques"].(string); ok {
        data.MitreTechniques = NewJSONSubsetValue(val)
    } else {
        data.MitreTechniques = NewJSONSubsetNull()
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
    data.CaseNumber = numberFromAPI(data.CaseNumber, dataMap["caseNumber"])
    if obj, ok := dataMap["caseNumberWithPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.CaseNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["caseNumberWithPrefix"].(string); ok {
        data.CaseNumberWithPrefix = types.StringValue(val)
    } else {
        data.CaseNumberWithPrefix = types.StringNull()
    }
    if obj, ok := dataMap["source"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    data.DetectionCount = numberFromAPI(data.DetectionCount, dataMap["detectionCount"])
    if obj, ok := dataMap["firstDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FirstDetectedAt = NewRFC3339Value(val)
        } else {
            data.FirstDetectedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["firstDetectedAt"].(string); ok && val != "" {
        data.FirstDetectedAt = NewRFC3339Value(val)
    } else {
        data.FirstDetectedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastDetectedAt = NewRFC3339Value(val)
        } else {
            data.LastDetectedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastDetectedAt"].(string); ok && val != "" {
        data.LastDetectedAt = NewRFC3339Value(val)
    } else {
        data.LastDetectedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgedAt"].(string); ok && val != "" {
        data.AcknowledgedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainedAt = NewRFC3339Value(val)
        } else {
            data.ContainedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containedAt"].(string); ok && val != "" {
        data.ContainedAt = NewRFC3339Value(val)
    } else {
        data.ContainedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["closedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ClosedAt = NewRFC3339Value(val)
        } else {
            data.ClosedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["closedAt"].(string); ok && val != "" {
        data.ClosedAt = NewRFC3339Value(val)
    } else {
        data.ClosedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeDueAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeDueAt"].(string); ok && val != "" {
        data.AcknowledgeDueAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainDueAt = NewRFC3339Value(val)
        } else {
            data.ContainDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containDueAt"].(string); ok && val != "" {
        data.ContainDueAt = NewRFC3339Value(val)
    } else {
        data.ContainDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeWarnedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeWarnedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeWarnedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeWarnedAt"].(string); ok && val != "" {
        data.AcknowledgeWarnedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeWarnedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeMissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeMissedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeMissedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeMissedAt"].(string); ok && val != "" {
        data.AcknowledgeMissedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeMissedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containWarnedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainWarnedAt = NewRFC3339Value(val)
        } else {
            data.ContainWarnedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containWarnedAt"].(string); ok && val != "" {
        data.ContainWarnedAt = NewRFC3339Value(val)
    } else {
        data.ContainWarnedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containMissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainMissedAt = NewRFC3339Value(val)
        } else {
            data.ContainMissedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containMissedAt"].(string); ok && val != "" {
        data.ContainMissedAt = NewRFC3339Value(val)
    } else {
        data.ContainMissedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextDueAt = NewRFC3339Value(val)
        } else {
            data.NextDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextDueAt"].(string); ok && val != "" {
        data.NextDueAt = NewRFC3339Value(val)
    } else {
        data.NextDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextTargetCheckAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextTargetCheckAt = NewRFC3339Value(val)
        } else {
            data.NextTargetCheckAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextTargetCheckAt"].(string); ok && val != "" {
        data.NextTargetCheckAt = NewRFC3339Value(val)
    } else {
        data.NextTargetCheckAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["closedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClosedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClosedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClosedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ClosedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["closedByUserId"].(string); ok {
        data.ClosedByUserId = types.StringValue(val)
    } else {
        data.ClosedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["groupingKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GroupingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GroupingKey = types.StringValue(string(jsonBytes))
            } else {
                data.GroupingKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GroupingKey = types.StringValue(string(jsonBytes))
            } else {
                data.GroupingKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GroupingKey = types.StringValue(string(jsonBytes))
        } else {
            data.GroupingKey = types.StringNull()
        }
    } else if val, ok := dataMap["groupingKey"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    } else {
        data.GroupingKey = types.StringNull()
    }
    if obj, ok := dataMap["notifiedSecurityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotifiedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NotifiedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.NotifiedSecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["notifiedSecurityCasePriorityId"].(string); ok {
        data.NotifiedSecurityCasePriorityId = types.StringValue(val)
    } else {
        data.NotifiedSecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["aiVerdict"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiVerdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiVerdict = types.StringValue(string(jsonBytes))
            } else {
                data.AiVerdict = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiVerdict = types.StringValue(string(jsonBytes))
            } else {
                data.AiVerdict = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdict = types.StringNull()
        }
    } else if val, ok := dataMap["aiVerdict"].(string); ok {
        data.AiVerdict = types.StringValue(val)
    } else {
        data.AiVerdict = types.StringNull()
    }
    if obj, ok := dataMap["aiSuggestedSecurityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.AiSuggestedSecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["aiSuggestedSecurityCasePriorityId"].(string); ok {
        data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
    } else {
        data.AiSuggestedSecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["aiSuggestedPlaybookId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiSuggestedPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiSuggestedPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
        } else {
            data.AiSuggestedPlaybookId = types.StringNull()
        }
    } else if val, ok := dataMap["aiSuggestedPlaybookId"].(string); ok {
        data.AiSuggestedPlaybookId = types.StringValue(val)
    } else {
        data.AiSuggestedPlaybookId = types.StringNull()
    }
    if obj, ok := dataMap["aiTriageReport"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiTriageReport = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiTriageReport = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiTriageReport = types.StringValue(string(jsonBytes))
            } else {
                data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiTriageReport = types.StringValue(string(jsonBytes))
            } else {
                data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiTriageReport = types.StringValue(string(jsonBytes))
        } else {
            data.AiTriageReport = types.StringNull()
        }
    } else if val, ok := dataMap["aiTriageReport"].(string); ok {
        data.AiTriageReport = types.StringValue(val)
    } else {
        data.AiTriageReport = types.StringNull()
    }
    if obj, ok := dataMap["aiTriagedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiTriagedAt = NewRFC3339Value(val)
        } else {
            data.AiTriagedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiTriagedAt"].(string); ok && val != "" {
        data.AiTriagedAt = NewRFC3339Value(val)
    } else {
        data.AiTriagedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["aiTriagePromptInjectionSuspected"].(bool); ok {
        data.AiTriagePromptInjectionSuspected = types.BoolValue(val)
    } else {
        data.AiTriagePromptInjectionSuspected = types.BoolNull()
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

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecurityCaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data SecurityCaseResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "securityCaseStatusId": true,
        "securityCasePriorityId": true,
        "verdict": true,
        "closingNote": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "entities": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "caseNumber": true,
        "caseNumberWithPrefix": true,
        "source": true,
        "detectionCount": true,
        "firstDetectedAt": true,
        "lastDetectedAt": true,
        "acknowledgedAt": true,
        "containedAt": true,
        "closedAt": true,
        "acknowledgeDueAt": true,
        "containDueAt": true,
        "acknowledgeWarnedAt": true,
        "acknowledgeMissedAt": true,
        "containWarnedAt": true,
        "containMissedAt": true,
        "nextDueAt": true,
        "nextTargetCheckAt": true,
        "closedByUserId": true,
        "groupingKey": true,
        "notifiedSecurityCasePriorityId": true,
        "incidentId": true,
        "aiVerdict": true,
        "aiSuggestedSecurityCasePriorityId": true,
        "aiSuggestedPlaybookId": true,
        "aiTriageReport": true,
        "aiTriagedAt": true,
        "aiTriagePromptInjectionSuspected": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/security-case/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var securityCaseResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &securityCaseResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse security_case response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := securityCaseResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = securityCaseResponse
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
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["securityCaseStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCaseStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCaseStatusId"].(string); ok {
        data.SecurityCaseStatusId = types.StringValue(val)
    } else {
        data.SecurityCaseStatusId = types.StringNull()
    }
    if obj, ok := dataMap["securityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCasePriorityId"].(string); ok {
        data.SecurityCasePriorityId = types.StringValue(val)
    } else {
        data.SecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["verdict"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Verdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Verdict = types.StringValue(string(jsonBytes))
            } else {
                data.Verdict = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Verdict = types.StringValue(string(jsonBytes))
            } else {
                data.Verdict = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Verdict = types.StringValue(string(jsonBytes))
        } else {
            data.Verdict = types.StringNull()
        }
    } else if val, ok := dataMap["verdict"].(string); ok {
        data.Verdict = types.StringValue(val)
    } else {
        data.Verdict = types.StringNull()
    }
    if obj, ok := dataMap["closingNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClosingNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClosingNote = types.StringValue(string(jsonBytes))
            } else {
                data.ClosingNote = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClosingNote = types.StringValue(string(jsonBytes))
            } else {
                data.ClosingNote = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClosingNote = types.StringValue(string(jsonBytes))
        } else {
            data.ClosingNote = types.StringNull()
        }
    } else if val, ok := dataMap["closingNote"].(string); ok {
        data.ClosingNote = types.StringValue(val)
    } else {
        data.ClosingNote = types.StringNull()
    }
    if obj, ok := dataMap["assignedToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignedToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignedToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignedToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["assignedToUserId"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    } else {
        data.AssignedToUserId = types.StringNull()
    }
    if obj, ok := dataMap["assignedToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignedToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignedToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignedToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["assignedToTeamId"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    } else {
        data.AssignedToTeamId = types.StringNull()
    }
    if obj, ok := dataMap["entities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Entities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Entities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Entities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Entities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Entities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Entities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["entities"].(string); ok {
        data.Entities = NewJSONSubsetValue(val)
    } else {
        data.Entities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["mitreTactics"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTactics = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MitreTactics = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MitreTactics = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["mitreTactics"].(string); ok {
        data.MitreTactics = NewJSONSubsetValue(val)
    } else {
        data.MitreTactics = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["mitreTechniques"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTechniques = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MitreTechniques = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MitreTechniques = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["mitreTechniques"].(string); ok {
        data.MitreTechniques = NewJSONSubsetValue(val)
    } else {
        data.MitreTechniques = NewJSONSubsetNull()
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
    data.CaseNumber = numberFromAPI(data.CaseNumber, dataMap["caseNumber"])
    if obj, ok := dataMap["caseNumberWithPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.CaseNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["caseNumberWithPrefix"].(string); ok {
        data.CaseNumberWithPrefix = types.StringValue(val)
    } else {
        data.CaseNumberWithPrefix = types.StringNull()
    }
    if obj, ok := dataMap["source"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    data.DetectionCount = numberFromAPI(data.DetectionCount, dataMap["detectionCount"])
    if obj, ok := dataMap["firstDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FirstDetectedAt = NewRFC3339Value(val)
        } else {
            data.FirstDetectedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["firstDetectedAt"].(string); ok && val != "" {
        data.FirstDetectedAt = NewRFC3339Value(val)
    } else {
        data.FirstDetectedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastDetectedAt = NewRFC3339Value(val)
        } else {
            data.LastDetectedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastDetectedAt"].(string); ok && val != "" {
        data.LastDetectedAt = NewRFC3339Value(val)
    } else {
        data.LastDetectedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgedAt"].(string); ok && val != "" {
        data.AcknowledgedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainedAt = NewRFC3339Value(val)
        } else {
            data.ContainedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containedAt"].(string); ok && val != "" {
        data.ContainedAt = NewRFC3339Value(val)
    } else {
        data.ContainedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["closedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ClosedAt = NewRFC3339Value(val)
        } else {
            data.ClosedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["closedAt"].(string); ok && val != "" {
        data.ClosedAt = NewRFC3339Value(val)
    } else {
        data.ClosedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeDueAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeDueAt"].(string); ok && val != "" {
        data.AcknowledgeDueAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainDueAt = NewRFC3339Value(val)
        } else {
            data.ContainDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containDueAt"].(string); ok && val != "" {
        data.ContainDueAt = NewRFC3339Value(val)
    } else {
        data.ContainDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeWarnedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeWarnedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeWarnedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeWarnedAt"].(string); ok && val != "" {
        data.AcknowledgeWarnedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeWarnedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeMissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeMissedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeMissedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeMissedAt"].(string); ok && val != "" {
        data.AcknowledgeMissedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeMissedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containWarnedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainWarnedAt = NewRFC3339Value(val)
        } else {
            data.ContainWarnedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containWarnedAt"].(string); ok && val != "" {
        data.ContainWarnedAt = NewRFC3339Value(val)
    } else {
        data.ContainWarnedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containMissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainMissedAt = NewRFC3339Value(val)
        } else {
            data.ContainMissedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containMissedAt"].(string); ok && val != "" {
        data.ContainMissedAt = NewRFC3339Value(val)
    } else {
        data.ContainMissedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextDueAt = NewRFC3339Value(val)
        } else {
            data.NextDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextDueAt"].(string); ok && val != "" {
        data.NextDueAt = NewRFC3339Value(val)
    } else {
        data.NextDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextTargetCheckAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextTargetCheckAt = NewRFC3339Value(val)
        } else {
            data.NextTargetCheckAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextTargetCheckAt"].(string); ok && val != "" {
        data.NextTargetCheckAt = NewRFC3339Value(val)
    } else {
        data.NextTargetCheckAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["closedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClosedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClosedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClosedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ClosedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["closedByUserId"].(string); ok {
        data.ClosedByUserId = types.StringValue(val)
    } else {
        data.ClosedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["groupingKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GroupingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GroupingKey = types.StringValue(string(jsonBytes))
            } else {
                data.GroupingKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GroupingKey = types.StringValue(string(jsonBytes))
            } else {
                data.GroupingKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GroupingKey = types.StringValue(string(jsonBytes))
        } else {
            data.GroupingKey = types.StringNull()
        }
    } else if val, ok := dataMap["groupingKey"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    } else {
        data.GroupingKey = types.StringNull()
    }
    if obj, ok := dataMap["notifiedSecurityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotifiedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NotifiedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.NotifiedSecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["notifiedSecurityCasePriorityId"].(string); ok {
        data.NotifiedSecurityCasePriorityId = types.StringValue(val)
    } else {
        data.NotifiedSecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["aiVerdict"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiVerdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiVerdict = types.StringValue(string(jsonBytes))
            } else {
                data.AiVerdict = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiVerdict = types.StringValue(string(jsonBytes))
            } else {
                data.AiVerdict = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdict = types.StringNull()
        }
    } else if val, ok := dataMap["aiVerdict"].(string); ok {
        data.AiVerdict = types.StringValue(val)
    } else {
        data.AiVerdict = types.StringNull()
    }
    if obj, ok := dataMap["aiSuggestedSecurityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.AiSuggestedSecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["aiSuggestedSecurityCasePriorityId"].(string); ok {
        data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
    } else {
        data.AiSuggestedSecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["aiSuggestedPlaybookId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiSuggestedPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiSuggestedPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
        } else {
            data.AiSuggestedPlaybookId = types.StringNull()
        }
    } else if val, ok := dataMap["aiSuggestedPlaybookId"].(string); ok {
        data.AiSuggestedPlaybookId = types.StringValue(val)
    } else {
        data.AiSuggestedPlaybookId = types.StringNull()
    }
    if obj, ok := dataMap["aiTriageReport"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiTriageReport = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiTriageReport = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiTriageReport = types.StringValue(string(jsonBytes))
            } else {
                data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiTriageReport = types.StringValue(string(jsonBytes))
            } else {
                data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiTriageReport = types.StringValue(string(jsonBytes))
        } else {
            data.AiTriageReport = types.StringNull()
        }
    } else if val, ok := dataMap["aiTriageReport"].(string); ok {
        data.AiTriageReport = types.StringValue(val)
    } else {
        data.AiTriageReport = types.StringNull()
    }
    if obj, ok := dataMap["aiTriagedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiTriagedAt = NewRFC3339Value(val)
        } else {
            data.AiTriagedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiTriagedAt"].(string); ok && val != "" {
        data.AiTriagedAt = NewRFC3339Value(val)
    } else {
        data.AiTriagedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["aiTriagePromptInjectionSuspected"].(bool); ok {
        data.AiTriagePromptInjectionSuspected = types.BoolValue(val)
    } else {
        data.AiTriagePromptInjectionSuspected = types.BoolNull()
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

func (r *SecurityCaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data SecurityCaseResourceModel
    var state SecurityCaseResourceModel

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

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config SecurityCaseResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    securityCaseRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := securityCaseRequest["data"].(map[string]interface{})

    if !data.Title.IsUnknown() && !state.Title.IsUnknown() && !data.Title.Equal(state.Title) {
        requestDataMap["title"] = data.Title.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.SecurityCaseStatusId.IsUnknown() && !state.SecurityCaseStatusId.IsUnknown() && !data.SecurityCaseStatusId.Equal(state.SecurityCaseStatusId) {
        requestDataMap["securityCaseStatusId"] = data.SecurityCaseStatusId.ValueString()
    }
    if !data.SecurityCasePriorityId.IsUnknown() && !state.SecurityCasePriorityId.IsUnknown() && !data.SecurityCasePriorityId.Equal(state.SecurityCasePriorityId) {
        requestDataMap["securityCasePriorityId"] = data.SecurityCasePriorityId.ValueString()
    }
    if !data.Verdict.IsUnknown() && !state.Verdict.IsUnknown() && !data.Verdict.Equal(state.Verdict) {
        requestDataMap["verdict"] = data.Verdict.ValueString()
    }
    if !data.ClosingNote.IsUnknown() && !state.ClosingNote.IsUnknown() && !data.ClosingNote.Equal(state.ClosingNote) {
        requestDataMap["closingNote"] = data.ClosingNote.ValueString()
    }
    if !data.AssignedToUserId.IsUnknown() && !state.AssignedToUserId.IsUnknown() && !data.AssignedToUserId.Equal(state.AssignedToUserId) {
        requestDataMap["assignedToUserId"] = data.AssignedToUserId.ValueString()
    }
    if !data.AssignedToTeamId.IsUnknown() && !state.AssignedToTeamId.IsUnknown() && !data.AssignedToTeamId.Equal(state.AssignedToTeamId) {
        requestDataMap["assignedToTeamId"] = data.AssignedToTeamId.ValueString()
    }
    if !data.Entities.IsUnknown() && !state.Entities.IsUnknown() && !data.Entities.Equal(state.Entities) {
        var entitiesData interface{}
        if err := json.Unmarshal([]byte(data.Entities.ValueString()), &entitiesData); err == nil {
            requestDataMap["entities"] = entitiesData
        } else {
            requestDataMap["entities"] = data.Entities.ValueString()
        }
    }
    if !data.MitreTactics.IsUnknown() && !state.MitreTactics.IsUnknown() && !data.MitreTactics.Equal(state.MitreTactics) {
        var mitretacticsData interface{}
        if err := json.Unmarshal([]byte(data.MitreTactics.ValueString()), &mitretacticsData); err == nil {
            requestDataMap["mitreTactics"] = mitretacticsData
        } else {
            requestDataMap["mitreTactics"] = data.MitreTactics.ValueString()
        }
    }
    if !data.MitreTechniques.IsUnknown() && !state.MitreTechniques.IsUnknown() && !data.MitreTechniques.Equal(state.MitreTechniques) {
        var mitretechniquesData interface{}
        if err := json.Unmarshal([]byte(data.MitreTechniques.ValueString()), &mitretechniquesData); err == nil {
            requestDataMap["mitreTechniques"] = mitretechniquesData
        } else {
            requestDataMap["mitreTechniques"] = data.MitreTechniques.ValueString()
        }
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(securityCaseRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/security-case/" + data.Id.ValueString() + "", securityCaseRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update security_case, got error: %s", err))
            return
        }

        // Parse the update response
        var securityCaseResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &securityCaseResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update security_case: %s", err))
            return
        }
        _ = securityCaseResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "securityCaseStatusId": true,
        "securityCasePriorityId": true,
        "verdict": true,
        "closingNote": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "entities": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "caseNumber": true,
        "caseNumberWithPrefix": true,
        "source": true,
        "detectionCount": true,
        "firstDetectedAt": true,
        "lastDetectedAt": true,
        "acknowledgedAt": true,
        "containedAt": true,
        "closedAt": true,
        "acknowledgeDueAt": true,
        "containDueAt": true,
        "acknowledgeWarnedAt": true,
        "acknowledgeMissedAt": true,
        "containWarnedAt": true,
        "containMissedAt": true,
        "nextDueAt": true,
        "nextTargetCheckAt": true,
        "closedByUserId": true,
        "groupingKey": true,
        "notifiedSecurityCasePriorityId": true,
        "incidentId": true,
        "aiVerdict": true,
        "aiSuggestedSecurityCasePriorityId": true,
        "aiSuggestedPlaybookId": true,
        "aiTriageReport": true,
        "aiTriagedAt": true,
        "aiTriagePromptInjectionSuspected": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/security-case/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_case after update: %s", err))
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
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["securityCaseStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCaseStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCaseStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCaseStatusId"].(string); ok {
        data.SecurityCaseStatusId = types.StringValue(val)
    } else {
        data.SecurityCaseStatusId = types.StringNull()
    }
    if obj, ok := dataMap["securityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCasePriorityId"].(string); ok {
        data.SecurityCasePriorityId = types.StringValue(val)
    } else {
        data.SecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["verdict"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Verdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Verdict = types.StringValue(string(jsonBytes))
            } else {
                data.Verdict = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Verdict = types.StringValue(string(jsonBytes))
            } else {
                data.Verdict = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Verdict = types.StringValue(string(jsonBytes))
        } else {
            data.Verdict = types.StringNull()
        }
    } else if val, ok := dataMap["verdict"].(string); ok {
        data.Verdict = types.StringValue(val)
    } else {
        data.Verdict = types.StringNull()
    }
    if obj, ok := dataMap["closingNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClosingNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClosingNote = types.StringValue(string(jsonBytes))
            } else {
                data.ClosingNote = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClosingNote = types.StringValue(string(jsonBytes))
            } else {
                data.ClosingNote = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClosingNote = types.StringValue(string(jsonBytes))
        } else {
            data.ClosingNote = types.StringNull()
        }
    } else if val, ok := dataMap["closingNote"].(string); ok {
        data.ClosingNote = types.StringValue(val)
    } else {
        data.ClosingNote = types.StringNull()
    }
    if obj, ok := dataMap["assignedToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignedToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignedToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignedToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["assignedToUserId"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    } else {
        data.AssignedToUserId = types.StringNull()
    }
    if obj, ok := dataMap["assignedToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignedToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignedToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignedToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["assignedToTeamId"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    } else {
        data.AssignedToTeamId = types.StringNull()
    }
    if obj, ok := dataMap["entities"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Entities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Entities = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Entities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Entities = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Entities = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Entities = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Entities = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["entities"].(string); ok {
        data.Entities = NewJSONSubsetValue(val)
    } else {
        data.Entities = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["mitreTactics"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTactics = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MitreTactics = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTactics = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MitreTactics = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MitreTactics = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["mitreTactics"].(string); ok {
        data.MitreTactics = NewJSONSubsetValue(val)
    } else {
        data.MitreTactics = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["mitreTechniques"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTechniques = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MitreTechniques = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.MitreTechniques = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MitreTechniques = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.MitreTechniques = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["mitreTechniques"].(string); ok {
        data.MitreTechniques = NewJSONSubsetValue(val)
    } else {
        data.MitreTechniques = NewJSONSubsetNull()
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
    data.CaseNumber = numberFromAPI(data.CaseNumber, dataMap["caseNumber"])
    if obj, ok := dataMap["caseNumberWithPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.CaseNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["caseNumberWithPrefix"].(string); ok {
        data.CaseNumberWithPrefix = types.StringValue(val)
    } else {
        data.CaseNumberWithPrefix = types.StringNull()
    }
    if obj, ok := dataMap["source"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    data.DetectionCount = numberFromAPI(data.DetectionCount, dataMap["detectionCount"])
    if obj, ok := dataMap["firstDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FirstDetectedAt = NewRFC3339Value(val)
        } else {
            data.FirstDetectedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["firstDetectedAt"].(string); ok && val != "" {
        data.FirstDetectedAt = NewRFC3339Value(val)
    } else {
        data.FirstDetectedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastDetectedAt = NewRFC3339Value(val)
        } else {
            data.LastDetectedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastDetectedAt"].(string); ok && val != "" {
        data.LastDetectedAt = NewRFC3339Value(val)
    } else {
        data.LastDetectedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgedAt"].(string); ok && val != "" {
        data.AcknowledgedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainedAt = NewRFC3339Value(val)
        } else {
            data.ContainedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containedAt"].(string); ok && val != "" {
        data.ContainedAt = NewRFC3339Value(val)
    } else {
        data.ContainedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["closedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ClosedAt = NewRFC3339Value(val)
        } else {
            data.ClosedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["closedAt"].(string); ok && val != "" {
        data.ClosedAt = NewRFC3339Value(val)
    } else {
        data.ClosedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeDueAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeDueAt"].(string); ok && val != "" {
        data.AcknowledgeDueAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainDueAt = NewRFC3339Value(val)
        } else {
            data.ContainDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containDueAt"].(string); ok && val != "" {
        data.ContainDueAt = NewRFC3339Value(val)
    } else {
        data.ContainDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeWarnedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeWarnedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeWarnedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeWarnedAt"].(string); ok && val != "" {
        data.AcknowledgeWarnedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeWarnedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["acknowledgeMissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AcknowledgeMissedAt = NewRFC3339Value(val)
        } else {
            data.AcknowledgeMissedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["acknowledgeMissedAt"].(string); ok && val != "" {
        data.AcknowledgeMissedAt = NewRFC3339Value(val)
    } else {
        data.AcknowledgeMissedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containWarnedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainWarnedAt = NewRFC3339Value(val)
        } else {
            data.ContainWarnedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containWarnedAt"].(string); ok && val != "" {
        data.ContainWarnedAt = NewRFC3339Value(val)
    } else {
        data.ContainWarnedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["containMissedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContainMissedAt = NewRFC3339Value(val)
        } else {
            data.ContainMissedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["containMissedAt"].(string); ok && val != "" {
        data.ContainMissedAt = NewRFC3339Value(val)
    } else {
        data.ContainMissedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextDueAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextDueAt = NewRFC3339Value(val)
        } else {
            data.NextDueAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextDueAt"].(string); ok && val != "" {
        data.NextDueAt = NewRFC3339Value(val)
    } else {
        data.NextDueAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["nextTargetCheckAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextTargetCheckAt = NewRFC3339Value(val)
        } else {
            data.NextTargetCheckAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextTargetCheckAt"].(string); ok && val != "" {
        data.NextTargetCheckAt = NewRFC3339Value(val)
    } else {
        data.NextTargetCheckAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["closedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClosedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClosedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClosedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ClosedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["closedByUserId"].(string); ok {
        data.ClosedByUserId = types.StringValue(val)
    } else {
        data.ClosedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["groupingKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GroupingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GroupingKey = types.StringValue(string(jsonBytes))
            } else {
                data.GroupingKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GroupingKey = types.StringValue(string(jsonBytes))
            } else {
                data.GroupingKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GroupingKey = types.StringValue(string(jsonBytes))
        } else {
            data.GroupingKey = types.StringNull()
        }
    } else if val, ok := dataMap["groupingKey"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    } else {
        data.GroupingKey = types.StringNull()
    }
    if obj, ok := dataMap["notifiedSecurityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotifiedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.NotifiedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.NotifiedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.NotifiedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.NotifiedSecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["notifiedSecurityCasePriorityId"].(string); ok {
        data.NotifiedSecurityCasePriorityId = types.StringValue(val)
    } else {
        data.NotifiedSecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["incidentId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := dataMap["aiVerdict"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiVerdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiVerdict = types.StringValue(string(jsonBytes))
            } else {
                data.AiVerdict = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiVerdict = types.StringValue(string(jsonBytes))
            } else {
                data.AiVerdict = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdict = types.StringNull()
        }
    } else if val, ok := dataMap["aiVerdict"].(string); ok {
        data.AiVerdict = types.StringValue(val)
    } else {
        data.AiVerdict = types.StringNull()
    }
    if obj, ok := dataMap["aiSuggestedSecurityCasePriorityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedSecurityCasePriorityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiSuggestedSecurityCasePriorityId = types.StringValue(string(jsonBytes))
        } else {
            data.AiSuggestedSecurityCasePriorityId = types.StringNull()
        }
    } else if val, ok := dataMap["aiSuggestedSecurityCasePriorityId"].(string); ok {
        data.AiSuggestedSecurityCasePriorityId = types.StringValue(val)
    } else {
        data.AiSuggestedSecurityCasePriorityId = types.StringNull()
    }
    if obj, ok := dataMap["aiSuggestedPlaybookId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiSuggestedPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiSuggestedPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
            } else {
                data.AiSuggestedPlaybookId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiSuggestedPlaybookId = types.StringValue(string(jsonBytes))
        } else {
            data.AiSuggestedPlaybookId = types.StringNull()
        }
    } else if val, ok := dataMap["aiSuggestedPlaybookId"].(string); ok {
        data.AiSuggestedPlaybookId = types.StringValue(val)
    } else {
        data.AiSuggestedPlaybookId = types.StringNull()
    }
    if obj, ok := dataMap["aiTriageReport"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiTriageReport = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiTriageReport = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiTriageReport = types.StringValue(string(jsonBytes))
            } else {
                data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiTriageReport = types.StringValue(string(jsonBytes))
            } else {
                data.AiTriageReport = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiTriageReport = types.StringValue(string(jsonBytes))
        } else {
            data.AiTriageReport = types.StringNull()
        }
    } else if val, ok := dataMap["aiTriageReport"].(string); ok {
        data.AiTriageReport = types.StringValue(val)
    } else {
        data.AiTriageReport = types.StringNull()
    }
    if obj, ok := dataMap["aiTriagedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiTriagedAt = NewRFC3339Value(val)
        } else {
            data.AiTriagedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiTriagedAt"].(string); ok && val != "" {
        data.AiTriagedAt = NewRFC3339Value(val)
    } else {
        data.AiTriagedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["aiTriagePromptInjectionSuspected"].(bool); ok {
        data.AiTriagePromptInjectionSuspected = types.BoolValue(val)
    } else {
        data.AiTriagePromptInjectionSuspected = types.BoolNull()
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

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecurityCaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data SecurityCaseResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/security-case/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete security_case, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete security_case: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *SecurityCaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *SecurityCaseResource) keepPlannedValues(data *SecurityCaseResourceModel, plan *SecurityCaseResourceModel, config *SecurityCaseResourceModel) {
    if config.Description.IsNull() && !plan.Description.IsUnknown() {
        data.Description = plan.Description
    }
    if config.SecurityCaseStatusId.IsNull() && !plan.SecurityCaseStatusId.IsUnknown() {
        data.SecurityCaseStatusId = plan.SecurityCaseStatusId
    }
    if config.SecurityCasePriorityId.IsNull() && !plan.SecurityCasePriorityId.IsUnknown() {
        data.SecurityCasePriorityId = plan.SecurityCasePriorityId
    }
    if config.Verdict.IsNull() && !plan.Verdict.IsUnknown() {
        data.Verdict = plan.Verdict
    }
    if config.ClosingNote.IsNull() && !plan.ClosingNote.IsUnknown() {
        data.ClosingNote = plan.ClosingNote
    }
    if config.AssignedToUserId.IsNull() && !plan.AssignedToUserId.IsUnknown() {
        data.AssignedToUserId = plan.AssignedToUserId
    }
    if config.AssignedToTeamId.IsNull() && !plan.AssignedToTeamId.IsUnknown() {
        data.AssignedToTeamId = plan.AssignedToTeamId
    }
    if config.Entities.IsNull() && !plan.Entities.IsUnknown() {
        data.Entities = plan.Entities
    }
    if config.MitreTactics.IsNull() && !plan.MitreTactics.IsUnknown() {
        data.MitreTactics = plan.MitreTactics
    }
    if config.MitreTechniques.IsNull() && !plan.MitreTechniques.IsUnknown() {
        data.MitreTechniques = plan.MitreTechniques
    }
    if config.Labels.IsNull() && !plan.Labels.IsUnknown() {
        data.Labels = plan.Labels
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *SecurityCaseResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *SecurityCaseResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *SecurityCaseResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *SecurityCaseResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *SecurityCaseResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *SecurityCaseResource) normalizeURLString(value string) string {
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
func (r *SecurityCaseResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *SecurityCaseResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
