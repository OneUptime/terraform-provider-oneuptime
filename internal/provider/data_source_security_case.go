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
var _ datasource.DataSource = &SecurityCaseDataSource{}

func NewSecurityCaseDataSource() datasource.DataSource {
    return &SecurityCaseDataSource{}
}

// SecurityCaseDataSource defines the data source implementation.
type SecurityCaseDataSource struct {
    client *Client
}

// SecurityCaseDataSourceModel describes the data source data model.
type SecurityCaseDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    CaseNumber types.Number `tfsdk:"case_number"`
    CaseNumberWithPrefix types.String `tfsdk:"case_number_with_prefix"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Status types.String `tfsdk:"status"`
    Priority types.String `tfsdk:"priority"`
    Verdict types.String `tfsdk:"verdict"`
    ClosingNote types.String `tfsdk:"closing_note"`
    AssignedToUserId types.String `tfsdk:"assigned_to_user_id"`
    AssignedToTeamId types.String `tfsdk:"assigned_to_team_id"`
    Source types.String `tfsdk:"source"`
    Entities types.String `tfsdk:"entities"`
    MitreTactics types.String `tfsdk:"mitre_tactics"`
    MitreTechniques types.String `tfsdk:"mitre_techniques"`
    DetectionCount types.Number `tfsdk:"detection_count"`
    FirstDetectedAt types.String `tfsdk:"first_detected_at"`
    LastDetectedAt types.String `tfsdk:"last_detected_at"`
    AcknowledgedAt types.String `tfsdk:"acknowledged_at"`
    ContainedAt types.String `tfsdk:"contained_at"`
    ClosedAt types.String `tfsdk:"closed_at"`
    ClosedByUserId types.String `tfsdk:"closed_by_user_id"`
    GroupingKey types.String `tfsdk:"grouping_key"`
    NotifiedPriority types.String `tfsdk:"notified_priority"`
    IncidentId types.String `tfsdk:"incident_id"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityCaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_case"
}

func (d *SecurityCaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A possible threat to investigate: the detections that point at it, who is on it, its priority and status, and the verdict it closed with. Look up an existing security case by `id`, or by any of its other arguments (`assigned_to_team_id`, `assigned_to_user_id`, `case_number`, ...): each one set must match, and exactly one security case may match them all.",

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
                MarkdownDescription: "ID of the project this case belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "case_number": schema.NumberAttribute{
                MarkdownDescription: "The case's number in this project, counting up from 1 as cases are opened.",
                Optional: true,
                Computed: true,
            },
            "case_number_with_prefix": schema.StringAttribute{
                MarkdownDescription: "The case number with the project's case number prefix in front of it, like CASE-42.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "What the case is about, in a line.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What is known about the threat, in Markdown. A case opened by a detection starts with the detection's description.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "New, Investigating, Contained or Closed. Every case starts New; closing it records a verdict.",
                Optional: true,
                Computed: true,
            },
            "priority": schema.StringAttribute{
                MarkdownDescription: "Critical, High, Medium, Low or Informational. A case opened by a detection starts at the detection's severity and rises when a more severe one joins it.",
                Optional: true,
                Computed: true,
            },
            "verdict": schema.StringAttribute{
                MarkdownDescription: "What the case turned out to be, set when it is closed: True Positive, Benign Positive, False Positive or Undetermined.",
                Optional: true,
                Computed: true,
            },
            "closing_note": schema.StringAttribute{
                MarkdownDescription: "Why the case was closed the way it was, in Markdown: what was found and what was done.",
                Optional: true,
                Computed: true,
            },
            "assigned_to_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person working the case. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "assigned_to_team_id": schema.StringAttribute{
                MarkdownDescription: "ID of the team whose queue the case is in. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "What opened the case: a Detection Rule, Threat Intel, a Connected Product, or a person (Manual).",
                Optional: true,
                Computed: true,
            },
            "entities": schema.StringAttribute{
                MarkdownDescription: "The users, hosts, IP addresses, domains and files the case is about, each as { type, value }. Detections add theirs as they join the case. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "mitre_tactics": schema.StringAttribute{
                MarkdownDescription: "ATT&CK tactic ids (TA0006) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "mitre_techniques": schema.StringAttribute{
                MarkdownDescription: "ATT&CK technique ids (T1110) the case's detections are tagged with. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "detection_count": schema.NumberAttribute{
                MarkdownDescription: "How many distinct detections the case holds.",
                Optional: true,
                Computed: true,
            },
            "first_detected_at": schema.StringAttribute{
                MarkdownDescription: "When the case's earliest detection fired.",
                Computed: true,
            },
            "last_detected_at": schema.StringAttribute{
                MarkdownDescription: "When a detection of the case last fired.",
                Computed: true,
            },
            "acknowledged_at": schema.StringAttribute{
                MarkdownDescription: "When somebody first took the case on: assigned it to a person, or moved it out of New. The time to acknowledge is measured to here.",
                Computed: true,
            },
            "contained_at": schema.StringAttribute{
                MarkdownDescription: "When the case first reached Contained. The time to contain is measured to here.",
                Computed: true,
            },
            "closed_at": schema.StringAttribute{
                MarkdownDescription: "When the case was closed. Cleared if it is reopened. The time to close is measured to here.",
                Computed: true,
            },
            "closed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who closed the case. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "grouping_key": schema.StringAttribute{
                MarkdownDescription: "What new detections are matched on to join this case while it is open: the entity they are about, or the detection's own fingerprint.",
                Optional: true,
                Computed: true,
            },
            "notified_priority": schema.StringAttribute{
                MarkdownDescription: "The highest priority people have been told about this case at. They are told again only when it rises above it.",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the OneUptime incident declared for this case. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to filter and report on cases by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who opened the case by hand. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *SecurityCaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityCaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityCaseDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CaseNumber.IsNull() && !data.CaseNumber.IsUnknown() {
        filters["caseNumber"] = lookupNumber(data.CaseNumber)
        filterNames = append(filterNames, "case_number = "+data.CaseNumber.ValueBigFloat().String())
    }
    if !data.CaseNumberWithPrefix.IsNull() && !data.CaseNumberWithPrefix.IsUnknown() {
        filters["caseNumberWithPrefix"] = data.CaseNumberWithPrefix.ValueString()
        filterNames = append(filterNames, "case_number_with_prefix = "+fmt.Sprintf("%q", data.CaseNumberWithPrefix.ValueString()))
    }
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
        filters["priority"] = data.Priority.ValueString()
        filterNames = append(filterNames, "priority = "+fmt.Sprintf("%q", data.Priority.ValueString()))
    }
    if !data.Verdict.IsNull() && !data.Verdict.IsUnknown() {
        filters["verdict"] = data.Verdict.ValueString()
        filterNames = append(filterNames, "verdict = "+fmt.Sprintf("%q", data.Verdict.ValueString()))
    }
    if !data.ClosingNote.IsNull() && !data.ClosingNote.IsUnknown() {
        filters["closingNote"] = data.ClosingNote.ValueString()
        filterNames = append(filterNames, "closing_note = "+fmt.Sprintf("%q", data.ClosingNote.ValueString()))
    }
    if !data.AssignedToUserId.IsNull() && !data.AssignedToUserId.IsUnknown() {
        filters["assignedToUserId"] = data.AssignedToUserId.ValueString()
        filterNames = append(filterNames, "assigned_to_user_id = "+fmt.Sprintf("%q", data.AssignedToUserId.ValueString()))
    }
    if !data.AssignedToTeamId.IsNull() && !data.AssignedToTeamId.IsUnknown() {
        filters["assignedToTeamId"] = data.AssignedToTeamId.ValueString()
        filterNames = append(filterNames, "assigned_to_team_id = "+fmt.Sprintf("%q", data.AssignedToTeamId.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.DetectionCount.IsNull() && !data.DetectionCount.IsUnknown() {
        filters["detectionCount"] = lookupNumber(data.DetectionCount)
        filterNames = append(filterNames, "detection_count = "+data.DetectionCount.ValueBigFloat().String())
    }
    if !data.ClosedByUserId.IsNull() && !data.ClosedByUserId.IsUnknown() {
        filters["closedByUserId"] = data.ClosedByUserId.ValueString()
        filterNames = append(filterNames, "closed_by_user_id = "+fmt.Sprintf("%q", data.ClosedByUserId.ValueString()))
    }
    if !data.GroupingKey.IsNull() && !data.GroupingKey.IsUnknown() {
        filters["groupingKey"] = data.GroupingKey.ValueString()
        filterNames = append(filterNames, "grouping_key = "+fmt.Sprintf("%q", data.GroupingKey.ValueString()))
    }
    if !data.NotifiedPriority.IsNull() && !data.NotifiedPriority.IsUnknown() {
        filters["notifiedPriority"] = data.NotifiedPriority.ValueString()
        filterNames = append(filterNames, "notified_priority = "+fmt.Sprintf("%q", data.NotifiedPriority.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security case up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security case up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "caseNumber": true,
        "caseNumberWithPrefix": true,
        "title": true,
        "description": true,
        "status": true,
        "priority": true,
        "verdict": true,
        "closingNote": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "source": true,
        "entities": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "detectionCount": true,
        "firstDetectedAt": true,
        "lastDetectedAt": true,
        "acknowledgedAt": true,
        "containedAt": true,
        "closedAt": true,
        "closedByUserId": true,
        "groupingKey": true,
        "notifiedPriority": true,
        "incidentId": true,
        "labels": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-case/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_case: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-case/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_case, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_case: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security case matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_case.")
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
    data.CaseNumber = numberFromAPI(data.CaseNumber, item["caseNumber"])
    if obj, ok := item["caseNumberWithPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CaseNumberWithPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CaseNumberWithPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CaseNumberWithPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.CaseNumberWithPrefix = types.StringNull()
        }
    } else if val, ok := item["caseNumberWithPrefix"].(string); ok {
        data.CaseNumberWithPrefix = types.StringValue(val)
    } else {
        data.CaseNumberWithPrefix = types.StringNull()
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
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    if obj, ok := item["priority"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Priority = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Priority = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Priority = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Priority = types.StringValue(string(jsonBytes))
        } else {
            data.Priority = types.StringNull()
        }
    } else if val, ok := item["priority"].(string); ok {
        data.Priority = types.StringValue(val)
    } else {
        data.Priority = types.StringNull()
    }
    if obj, ok := item["verdict"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Verdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Verdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Verdict = types.StringValue(string(jsonBytes))
        } else {
            data.Verdict = types.StringNull()
        }
    } else if val, ok := item["verdict"].(string); ok {
        data.Verdict = types.StringValue(val)
    } else {
        data.Verdict = types.StringNull()
    }
    if obj, ok := item["closingNote"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClosingNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClosingNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClosingNote = types.StringValue(string(jsonBytes))
        } else {
            data.ClosingNote = types.StringNull()
        }
    } else if val, ok := item["closingNote"].(string); ok {
        data.ClosingNote = types.StringValue(val)
    } else {
        data.ClosingNote = types.StringNull()
    }
    if obj, ok := item["assignedToUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignedToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignedToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignedToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToUserId = types.StringNull()
        }
    } else if val, ok := item["assignedToUserId"].(string); ok {
        data.AssignedToUserId = types.StringValue(val)
    } else {
        data.AssignedToUserId = types.StringNull()
    }
    if obj, ok := item["assignedToTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignedToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignedToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignedToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignedToTeamId = types.StringNull()
        }
    } else if val, ok := item["assignedToTeamId"].(string); ok {
        data.AssignedToTeamId = types.StringValue(val)
    } else {
        data.AssignedToTeamId = types.StringNull()
    }
    if obj, ok := item["source"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := item["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := item["entities"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Entities = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Entities = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Entities = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Entities = types.StringValue(string(jsonBytes))
        } else {
            data.Entities = types.StringNull()
        }
    } else if val, ok := item["entities"].(string); ok {
        data.Entities = types.StringValue(val)
    } else {
        data.Entities = types.StringNull()
    }
    if obj, ok := item["mitreTactics"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTactics = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MitreTactics = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MitreTactics = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MitreTactics = types.StringValue(string(jsonBytes))
        } else {
            data.MitreTactics = types.StringNull()
        }
    } else if val, ok := item["mitreTactics"].(string); ok {
        data.MitreTactics = types.StringValue(val)
    } else {
        data.MitreTactics = types.StringNull()
    }
    if obj, ok := item["mitreTechniques"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTechniques = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MitreTechniques = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MitreTechniques = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MitreTechniques = types.StringValue(string(jsonBytes))
        } else {
            data.MitreTechniques = types.StringNull()
        }
    } else if val, ok := item["mitreTechniques"].(string); ok {
        data.MitreTechniques = types.StringValue(val)
    } else {
        data.MitreTechniques = types.StringNull()
    }
    data.DetectionCount = numberFromAPI(data.DetectionCount, item["detectionCount"])
    if obj, ok := item["firstDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstDetectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstDetectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstDetectedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstDetectedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstDetectedAt = types.StringNull()
        }
    } else if val, ok := item["firstDetectedAt"].(string); ok {
        data.FirstDetectedAt = types.StringValue(val)
    } else {
        data.FirstDetectedAt = types.StringNull()
    }
    if obj, ok := item["lastDetectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastDetectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastDetectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastDetectedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastDetectedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastDetectedAt = types.StringNull()
        }
    } else if val, ok := item["lastDetectedAt"].(string); ok {
        data.LastDetectedAt = types.StringValue(val)
    } else {
        data.LastDetectedAt = types.StringNull()
    }
    if obj, ok := item["acknowledgedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcknowledgedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcknowledgedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcknowledgedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcknowledgedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AcknowledgedAt = types.StringNull()
        }
    } else if val, ok := item["acknowledgedAt"].(string); ok {
        data.AcknowledgedAt = types.StringValue(val)
    } else {
        data.AcknowledgedAt = types.StringNull()
    }
    if obj, ok := item["containedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContainedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContainedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContainedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ContainedAt = types.StringNull()
        }
    } else if val, ok := item["containedAt"].(string); ok {
        data.ContainedAt = types.StringValue(val)
    } else {
        data.ContainedAt = types.StringNull()
    }
    if obj, ok := item["closedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClosedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClosedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClosedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ClosedAt = types.StringNull()
        }
    } else if val, ok := item["closedAt"].(string); ok {
        data.ClosedAt = types.StringValue(val)
    } else {
        data.ClosedAt = types.StringNull()
    }
    if obj, ok := item["closedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClosedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClosedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClosedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ClosedByUserId = types.StringNull()
        }
    } else if val, ok := item["closedByUserId"].(string); ok {
        data.ClosedByUserId = types.StringValue(val)
    } else {
        data.ClosedByUserId = types.StringNull()
    }
    if obj, ok := item["groupingKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GroupingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GroupingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GroupingKey = types.StringValue(string(jsonBytes))
        } else {
            data.GroupingKey = types.StringNull()
        }
    } else if val, ok := item["groupingKey"].(string); ok {
        data.GroupingKey = types.StringValue(val)
    } else {
        data.GroupingKey = types.StringNull()
    }
    if obj, ok := item["notifiedPriority"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotifiedPriority = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NotifiedPriority = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NotifiedPriority = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NotifiedPriority = types.StringValue(string(jsonBytes))
        } else {
            data.NotifiedPriority = types.StringNull()
        }
    } else if val, ok := item["notifiedPriority"].(string); ok {
        data.NotifiedPriority = types.StringValue(val)
    } else {
        data.NotifiedPriority = types.StringNull()
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
