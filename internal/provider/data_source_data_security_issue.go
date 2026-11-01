package provider

import (
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
var _ datasource.DataSource = &DataSecurityIssueDataSource{}

func NewDataSecurityIssueDataSource() datasource.DataSource {
    return &DataSecurityIssueDataSource{}
}

// DataSecurityIssueDataSource defines the data source implementation.
type DataSecurityIssueDataSource struct {
    client *Client
}

// DataSecurityIssueDataSourceModel describes the data source data model.
type DataSecurityIssueDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    IssueKey types.String `tfsdk:"issue_key"`
    CheckId types.String `tfsdk:"check_id"`
    CheckGroup types.String `tfsdk:"check_group"`
    Title types.String `tfsdk:"title"`
    Evidence types.String `tfsdk:"evidence"`
    Severity types.String `tfsdk:"severity"`
    RiskScore types.Number `tfsdk:"risk_score"`
    RiskFactors types.String `tfsdk:"risk_factors"`
    Status types.String `tfsdk:"status"`
    DataTypeIds types.String `tfsdk:"data_type_ids"`
    Frameworks types.String `tfsdk:"frameworks"`
    AcceptedReason types.String `tfsdk:"accepted_reason"`
    AcceptedUntil types.String `tfsdk:"accepted_until"`
    AcceptedByUserId types.String `tfsdk:"accepted_by_user_id"`
    AssignedToUserId types.String `tfsdk:"assigned_to_user_id"`
    AssignedToTeamId types.String `tfsdk:"assigned_to_team_id"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    AlertId types.String `tfsdk:"alert_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    FirstDetectedAt types.String `tfsdk:"first_detected_at"`
    LastDetectedAt types.String `tfsdk:"last_detected_at"`
    ResolvedAt types.String `tfsdk:"resolved_at"`
    ReopenCount types.Number `tfsdk:"reopen_count"`
    RecheckRequestedAt types.String `tfsdk:"recheck_requested_at"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
}

func (d *DataSecurityIssueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_security_issue"
}

func (d *DataSecurityIssueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Everything wrong with a store that a check found, ranked by risk: what is wrong, why it matters, how to fix it, and who owns it. Look up an existing data security issue by `id`, or by any of its other arguments (`accepted_by_user_id`, `accepted_reason`, `alert_id`, ...): each one set must match, and exactly one data security issue may match them all.",

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
                MarkdownDescription: "ID of the project this issue belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store the issue is on. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
            },
            "issue_key": schema.StringAttribute{
                MarkdownDescription: "The store and check, which identify the issue within the project.",
                Optional: true,
                Computed: true,
            },
            "check_id": schema.StringAttribute{
                MarkdownDescription: "The check that failed.",
                Optional: true,
                Computed: true,
            },
            "check_group": schema.StringAttribute{
                MarkdownDescription: "The check's group: Exposure, Encryption, Recovery, Visibility, Access, Data hygiene, Telemetry, Residency or AI.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "What is wrong, in a few words.",
                Optional: true,
                Computed: true,
            },
            "evidence": schema.StringAttribute{
                MarkdownDescription: "What the check saw, in one sentence.",
                Optional: true,
                Computed: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "How severe the issue is for this store's data: Critical, High, Medium or Low.",
                Optional: true,
                Computed: true,
            },
            "risk_score": schema.NumberAttribute{
                MarkdownDescription: "Orders issues inside a severity, from 0 to 100: severity, exposure, volume, recent activity and production.",
                Optional: true,
                Computed: true,
            },
            "risk_factors": schema.StringAttribute{
                MarkdownDescription: "Each part of the risk score, with its points. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Open, Accepted, Resolved or Not an issue.",
                Optional: true,
                Computed: true,
            },
            "data_type_ids": schema.StringAttribute{
                MarkdownDescription: "The sensitive data types the store holds, which make the issue as severe as it is. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "frameworks": schema.StringAttribute{
                MarkdownDescription: "The framework requirements the check is evidence for. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "accepted_reason": schema.StringAttribute{
                MarkdownDescription: "Why the risk was accepted.",
                Optional: true,
                Computed: true,
            },
            "accepted_until": schema.StringAttribute{
                MarkdownDescription: "When an accepted issue opens again if its check still fails.",
                Computed: true,
            },
            "accepted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "The person who accepted the risk. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "assigned_to_user_id": schema.StringAttribute{
                MarkdownDescription: "The person working on the issue. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "assigned_to_team_id": schema.StringAttribute{
                MarkdownDescription: "The team working on the issue. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "The security case opened for the issue, if any. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "The alert raised for the issue, if any. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "The incident declared for the issue, if any. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "first_detected_at": schema.StringAttribute{
                MarkdownDescription: "When the check first failed.",
                Computed: true,
            },
            "last_detected_at": schema.StringAttribute{
                MarkdownDescription: "When the check last failed.",
                Computed: true,
            },
            "resolved_at": schema.StringAttribute{
                MarkdownDescription: "When the check last passed after failing.",
                Computed: true,
            },
            "reopen_count": schema.NumberAttribute{
                MarkdownDescription: "How many times the issue opened again after it was resolved.",
                Optional: true,
                Computed: true,
            },
            "recheck_requested_at": schema.StringAttribute{
                MarkdownDescription: "When a person last asked for the store to be checked again.",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Whether the store's owners were told about the issue.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DataSecurityIssueDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSecurityIssueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataSecurityIssueDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.IssueKey.IsNull() && !data.IssueKey.IsUnknown() {
        filters["issueKey"] = data.IssueKey.ValueString()
        filterNames = append(filterNames, "issue_key = "+fmt.Sprintf("%q", data.IssueKey.ValueString()))
    }
    if !data.CheckId.IsNull() && !data.CheckId.IsUnknown() {
        filters["checkId"] = data.CheckId.ValueString()
        filterNames = append(filterNames, "check_id = "+fmt.Sprintf("%q", data.CheckId.ValueString()))
    }
    if !data.CheckGroup.IsNull() && !data.CheckGroup.IsUnknown() {
        filters["checkGroup"] = data.CheckGroup.ValueString()
        filterNames = append(filterNames, "check_group = "+fmt.Sprintf("%q", data.CheckGroup.ValueString()))
    }
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Evidence.IsNull() && !data.Evidence.IsUnknown() {
        filters["evidence"] = data.Evidence.ValueString()
        filterNames = append(filterNames, "evidence = "+fmt.Sprintf("%q", data.Evidence.ValueString()))
    }
    if !data.Severity.IsNull() && !data.Severity.IsUnknown() {
        filters["severity"] = data.Severity.ValueString()
        filterNames = append(filterNames, "severity = "+fmt.Sprintf("%q", data.Severity.ValueString()))
    }
    if !data.RiskScore.IsNull() && !data.RiskScore.IsUnknown() {
        filters["riskScore"] = lookupNumber(data.RiskScore)
        filterNames = append(filterNames, "risk_score = "+data.RiskScore.ValueBigFloat().String())
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.AcceptedReason.IsNull() && !data.AcceptedReason.IsUnknown() {
        filters["acceptedReason"] = data.AcceptedReason.ValueString()
        filterNames = append(filterNames, "accepted_reason = "+fmt.Sprintf("%q", data.AcceptedReason.ValueString()))
    }
    if !data.AcceptedByUserId.IsNull() && !data.AcceptedByUserId.IsUnknown() {
        filters["acceptedByUserId"] = data.AcceptedByUserId.ValueString()
        filterNames = append(filterNames, "accepted_by_user_id = "+fmt.Sprintf("%q", data.AcceptedByUserId.ValueString()))
    }
    if !data.AssignedToUserId.IsNull() && !data.AssignedToUserId.IsUnknown() {
        filters["assignedToUserId"] = data.AssignedToUserId.ValueString()
        filterNames = append(filterNames, "assigned_to_user_id = "+fmt.Sprintf("%q", data.AssignedToUserId.ValueString()))
    }
    if !data.AssignedToTeamId.IsNull() && !data.AssignedToTeamId.IsUnknown() {
        filters["assignedToTeamId"] = data.AssignedToTeamId.ValueString()
        filterNames = append(filterNames, "assigned_to_team_id = "+fmt.Sprintf("%q", data.AssignedToTeamId.ValueString()))
    }
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.AlertId.IsNull() && !data.AlertId.IsUnknown() {
        filters["alertId"] = data.AlertId.ValueString()
        filterNames = append(filterNames, "alert_id = "+fmt.Sprintf("%q", data.AlertId.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.ReopenCount.IsNull() && !data.ReopenCount.IsUnknown() {
        filters["reopenCount"] = lookupNumber(data.ReopenCount)
        filterNames = append(filterNames, "reopen_count = "+data.ReopenCount.ValueBigFloat().String())
    }
    if !data.IsOwnerNotified.IsNull() && !data.IsOwnerNotified.IsUnknown() {
        filters["isOwnerNotified"] = data.IsOwnerNotified.ValueBool()
        filterNames = append(filterNames, "is_owner_notified = "+fmt.Sprintf("%t", data.IsOwnerNotified.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data security issue up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data security issue up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataStoreId": true,
        "issueKey": true,
        "checkId": true,
        "checkGroup": true,
        "title": true,
        "evidence": true,
        "severity": true,
        "riskScore": true,
        "riskFactors": true,
        "status": true,
        "dataTypeIds": true,
        "frameworks": true,
        "acceptedReason": true,
        "acceptedUntil": true,
        "acceptedByUserId": true,
        "assignedToUserId": true,
        "assignedToTeamId": true,
        "securityCaseId": true,
        "alertId": true,
        "incidentId": true,
        "firstDetectedAt": true,
        "lastDetectedAt": true,
        "resolvedAt": true,
        "reopenCount": true,
        "recheckRequestedAt": true,
        "isOwnerNotified": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-issue/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_issue, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security issue found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_security_issue: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-issue/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_security_issue, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_security_issue: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security issue matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data security issue matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_security_issue.")
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
    } else if val, ok := item["dataStoreId"].(string); ok {
        data.DataStoreId = types.StringValue(val)
    } else {
        data.DataStoreId = types.StringNull()
    }
    if obj, ok := item["issueKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IssueKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IssueKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IssueKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IssueKey = types.StringValue(string(jsonBytes))
        } else {
            data.IssueKey = types.StringNull()
        }
    } else if val, ok := item["issueKey"].(string); ok {
        data.IssueKey = types.StringValue(val)
    } else {
        data.IssueKey = types.StringNull()
    }
    if obj, ok := item["checkId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CheckId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CheckId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CheckId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CheckId = types.StringValue(string(jsonBytes))
        } else {
            data.CheckId = types.StringNull()
        }
    } else if val, ok := item["checkId"].(string); ok {
        data.CheckId = types.StringValue(val)
    } else {
        data.CheckId = types.StringNull()
    }
    if obj, ok := item["checkGroup"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CheckGroup = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CheckGroup = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CheckGroup = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CheckGroup = types.StringValue(string(jsonBytes))
        } else {
            data.CheckGroup = types.StringNull()
        }
    } else if val, ok := item["checkGroup"].(string); ok {
        data.CheckGroup = types.StringValue(val)
    } else {
        data.CheckGroup = types.StringNull()
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
    } else if val, ok := item["evidence"].(string); ok {
        data.Evidence = types.StringValue(val)
    } else {
        data.Evidence = types.StringNull()
    }
    if obj, ok := item["severity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Severity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Severity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Severity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Severity = types.StringValue(string(jsonBytes))
        } else {
            data.Severity = types.StringNull()
        }
    } else if val, ok := item["severity"].(string); ok {
        data.Severity = types.StringValue(val)
    } else {
        data.Severity = types.StringNull()
    }
    data.RiskScore = numberFromAPI(data.RiskScore, item["riskScore"])
    if obj, ok := item["riskFactors"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RiskFactors = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RiskFactors = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RiskFactors = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RiskFactors = types.StringValue(string(jsonBytes))
        } else {
            data.RiskFactors = types.StringNull()
        }
    } else if val, ok := item["riskFactors"].(string); ok {
        data.RiskFactors = types.StringValue(val)
    } else {
        data.RiskFactors = types.StringNull()
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
    if obj, ok := item["dataTypeIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataTypeIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataTypeIds = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeIds = types.StringNull()
        }
    } else if val, ok := item["dataTypeIds"].(string); ok {
        data.DataTypeIds = types.StringValue(val)
    } else {
        data.DataTypeIds = types.StringNull()
    }
    if obj, ok := item["frameworks"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Frameworks = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Frameworks = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Frameworks = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Frameworks = types.StringValue(string(jsonBytes))
        } else {
            data.Frameworks = types.StringNull()
        }
    } else if val, ok := item["frameworks"].(string); ok {
        data.Frameworks = types.StringValue(val)
    } else {
        data.Frameworks = types.StringNull()
    }
    if obj, ok := item["acceptedReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcceptedReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcceptedReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcceptedReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcceptedReason = types.StringValue(string(jsonBytes))
        } else {
            data.AcceptedReason = types.StringNull()
        }
    } else if val, ok := item["acceptedReason"].(string); ok {
        data.AcceptedReason = types.StringValue(val)
    } else {
        data.AcceptedReason = types.StringNull()
    }
    if obj, ok := item["acceptedUntil"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcceptedUntil = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcceptedUntil = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcceptedUntil = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcceptedUntil = types.StringValue(string(jsonBytes))
        } else {
            data.AcceptedUntil = types.StringNull()
        }
    } else if val, ok := item["acceptedUntil"].(string); ok {
        data.AcceptedUntil = types.StringValue(val)
    } else {
        data.AcceptedUntil = types.StringNull()
    }
    if obj, ok := item["acceptedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AcceptedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AcceptedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AcceptedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AcceptedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AcceptedByUserId = types.StringNull()
        }
    } else if val, ok := item["acceptedByUserId"].(string); ok {
        data.AcceptedByUserId = types.StringValue(val)
    } else {
        data.AcceptedByUserId = types.StringNull()
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
    } else if val, ok := item["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
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
    if obj, ok := item["resolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedAt = types.StringNull()
        }
    } else if val, ok := item["resolvedAt"].(string); ok {
        data.ResolvedAt = types.StringValue(val)
    } else {
        data.ResolvedAt = types.StringNull()
    }
    data.ReopenCount = numberFromAPI(data.ReopenCount, item["reopenCount"])
    if obj, ok := item["recheckRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecheckRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecheckRequestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecheckRequestedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecheckRequestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RecheckRequestedAt = types.StringNull()
        }
    } else if val, ok := item["recheckRequestedAt"].(string); ok {
        data.RecheckRequestedAt = types.StringValue(val)
    } else {
        data.RecheckRequestedAt = types.StringNull()
    }
    if val, ok := item["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else {
        data.IsOwnerNotified = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
