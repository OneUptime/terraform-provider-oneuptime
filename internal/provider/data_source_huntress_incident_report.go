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
var _ datasource.DataSource = &HuntressIncidentReportDataSource{}

func NewHuntressIncidentReportDataSource() datasource.DataSource {
    return &HuntressIncidentReportDataSource{}
}

// HuntressIncidentReportDataSource defines the data source implementation.
type HuntressIncidentReportDataSource struct {
    client *Client
}

// HuntressIncidentReportDataSourceModel describes the data source data model.
type HuntressIncidentReportDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    HuntressConnectionId types.String `tfsdk:"huntress_connection_id"`
    HuntressAccountId types.String `tfsdk:"huntress_account_id"`
    HuntressIncidentReportId types.String `tfsdk:"huntress_incident_report_id"`
    OrganizationId types.String `tfsdk:"organization_id"`
    OrganizationName types.String `tfsdk:"organization_name"`
    AffectedName types.String `tfsdk:"affected_name"`
    Subject types.String `tfsdk:"subject"`
    Severity types.String `tfsdk:"severity"`
    Status types.String `tfsdk:"status"`
    IncidentId types.String `tfsdk:"incident_id"`
    Outcome types.String `tfsdk:"outcome"`
    PagedOnCall types.Bool `tfsdk:"paged_on_call"`
    LastEventType types.String `tfsdk:"last_event_type"`
    LastEventReceivedAt types.String `tfsdk:"last_event_received_at"`
}

func (d *HuntressIncidentReportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_huntress_incident_report"
}

func (d *HuntressIncidentReportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Huntress incident reports received through a Huntress connection, each with the incident it opened or the reason it opened none. Look up an existing huntress incident report by `id`, or by any of its other arguments (`affected_name`, `huntress_account_id`, `huntress_connection_id`, ...): each one set must match, and exactly one huntress incident report may match them all.",

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
                MarkdownDescription: "ID of the project this report was received in. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "huntress_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the connection that received this report, or empty once it is deleted. The ID of a `oneuptime_huntress_connection`.",
                Optional: true,
                Computed: true,
            },
            "huntress_account_id": schema.StringAttribute{
                MarkdownDescription: "The id of the Huntress account the report belongs to, or empty when Huntress did not say.",
                Optional: true,
                Computed: true,
            },
            "huntress_incident_report_id": schema.StringAttribute{
                MarkdownDescription: "The incident report's id in Huntress.",
                Optional: true,
                Computed: true,
            },
            "organization_id": schema.StringAttribute{
                MarkdownDescription: "The id of the Huntress organization the report is about.",
                Optional: true,
                Computed: true,
            },
            "organization_name": schema.StringAttribute{
                MarkdownDescription: "The name of the Huntress organization the report is about.",
                Optional: true,
                Computed: true,
            },
            "affected_name": schema.StringAttribute{
                MarkdownDescription: "The host or identity the report is about, as Huntress names it in the report's subject.",
                Optional: true,
                Computed: true,
            },
            "subject": schema.StringAttribute{
                MarkdownDescription: "The report's subject in Huntress.",
                Optional: true,
                Computed: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "The report's severity in Huntress: critical, high or low, as Huntress last sent it.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "The report's status in Huntress as it last sent it, such as sent, closed or dismissed.",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident this report opened. Empty when it opened none, or when the incident was deleted. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "outcome": schema.StringAttribute{
                MarkdownDescription: "What was done with the report: Opening, IncidentOpened, IncidentResolved, OrganizationNotWatched or ClosedBeforeReceived.",
                Optional: true,
                Computed: true,
            },
            "paged_on_call": schema.BoolAttribute{
                MarkdownDescription: "Whether the incident was opened with the connection's on-call policies, which pages them.",
                Optional: true,
                Computed: true,
            },
            "last_event_type": schema.StringAttribute{
                MarkdownDescription: "The last event Huntress sent about this report.",
                Optional: true,
                Computed: true,
            },
            "last_event_received_at": schema.StringAttribute{
                MarkdownDescription: "When Huntress last sent an event about this report.",
                Computed: true,
            },
        },
    }
}

func (d *HuntressIncidentReportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HuntressIncidentReportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HuntressIncidentReportDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.HuntressConnectionId.IsNull() && !data.HuntressConnectionId.IsUnknown() {
        filters["huntressConnectionId"] = data.HuntressConnectionId.ValueString()
        filterNames = append(filterNames, "huntress_connection_id = "+fmt.Sprintf("%q", data.HuntressConnectionId.ValueString()))
    }
    if !data.HuntressAccountId.IsNull() && !data.HuntressAccountId.IsUnknown() {
        filters["huntressAccountId"] = data.HuntressAccountId.ValueString()
        filterNames = append(filterNames, "huntress_account_id = "+fmt.Sprintf("%q", data.HuntressAccountId.ValueString()))
    }
    if !data.HuntressIncidentReportId.IsNull() && !data.HuntressIncidentReportId.IsUnknown() {
        filters["huntressIncidentReportId"] = data.HuntressIncidentReportId.ValueString()
        filterNames = append(filterNames, "huntress_incident_report_id = "+fmt.Sprintf("%q", data.HuntressIncidentReportId.ValueString()))
    }
    if !data.OrganizationId.IsNull() && !data.OrganizationId.IsUnknown() {
        filters["organizationId"] = data.OrganizationId.ValueString()
        filterNames = append(filterNames, "organization_id = "+fmt.Sprintf("%q", data.OrganizationId.ValueString()))
    }
    if !data.OrganizationName.IsNull() && !data.OrganizationName.IsUnknown() {
        filters["organizationName"] = data.OrganizationName.ValueString()
        filterNames = append(filterNames, "organization_name = "+fmt.Sprintf("%q", data.OrganizationName.ValueString()))
    }
    if !data.AffectedName.IsNull() && !data.AffectedName.IsUnknown() {
        filters["affectedName"] = data.AffectedName.ValueString()
        filterNames = append(filterNames, "affected_name = "+fmt.Sprintf("%q", data.AffectedName.ValueString()))
    }
    if !data.Subject.IsNull() && !data.Subject.IsUnknown() {
        filters["subject"] = data.Subject.ValueString()
        filterNames = append(filterNames, "subject = "+fmt.Sprintf("%q", data.Subject.ValueString()))
    }
    if !data.Severity.IsNull() && !data.Severity.IsUnknown() {
        filters["severity"] = data.Severity.ValueString()
        filterNames = append(filterNames, "severity = "+fmt.Sprintf("%q", data.Severity.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.Outcome.IsNull() && !data.Outcome.IsUnknown() {
        filters["outcome"] = data.Outcome.ValueString()
        filterNames = append(filterNames, "outcome = "+fmt.Sprintf("%q", data.Outcome.ValueString()))
    }
    if !data.PagedOnCall.IsNull() && !data.PagedOnCall.IsUnknown() {
        filters["pagedOnCall"] = data.PagedOnCall.ValueBool()
        filterNames = append(filterNames, "paged_on_call = "+fmt.Sprintf("%t", data.PagedOnCall.ValueBool()))
    }
    if !data.LastEventType.IsNull() && !data.LastEventType.IsUnknown() {
        filters["lastEventType"] = data.LastEventType.ValueString()
        filterNames = append(filterNames, "last_event_type = "+fmt.Sprintf("%q", data.LastEventType.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the huntress incident report up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the huntress incident report up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "huntressConnectionId": true,
        "huntressAccountId": true,
        "huntressIncidentReportId": true,
        "organizationId": true,
        "organizationName": true,
        "affectedName": true,
        "subject": true,
        "severity": true,
        "status": true,
        "incidentId": true,
        "outcome": true,
        "pagedOnCall": true,
        "lastEventType": true,
        "lastEventReceivedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/huntress-incident-report/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read huntress_incident_report, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No huntress incident report found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read huntress_incident_report: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/huntress-incident-report/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list huntress_incident_report, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list huntress_incident_report: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No huntress incident report matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one huntress incident report matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for huntress_incident_report.")
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
    if obj, ok := item["huntressConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HuntressConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HuntressConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HuntressConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HuntressConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.HuntressConnectionId = types.StringNull()
        }
    } else if val, ok := item["huntressConnectionId"].(string); ok {
        data.HuntressConnectionId = types.StringValue(val)
    } else {
        data.HuntressConnectionId = types.StringNull()
    }
    if obj, ok := item["huntressAccountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HuntressAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HuntressAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HuntressAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HuntressAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.HuntressAccountId = types.StringNull()
        }
    } else if val, ok := item["huntressAccountId"].(string); ok {
        data.HuntressAccountId = types.StringValue(val)
    } else {
        data.HuntressAccountId = types.StringNull()
    }
    if obj, ok := item["huntressIncidentReportId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HuntressIncidentReportId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HuntressIncidentReportId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HuntressIncidentReportId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HuntressIncidentReportId = types.StringValue(string(jsonBytes))
        } else {
            data.HuntressIncidentReportId = types.StringNull()
        }
    } else if val, ok := item["huntressIncidentReportId"].(string); ok {
        data.HuntressIncidentReportId = types.StringValue(val)
    } else {
        data.HuntressIncidentReportId = types.StringNull()
    }
    if obj, ok := item["organizationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OrganizationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OrganizationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OrganizationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OrganizationId = types.StringValue(string(jsonBytes))
        } else {
            data.OrganizationId = types.StringNull()
        }
    } else if val, ok := item["organizationId"].(string); ok {
        data.OrganizationId = types.StringValue(val)
    } else {
        data.OrganizationId = types.StringNull()
    }
    if obj, ok := item["organizationName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OrganizationName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OrganizationName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OrganizationName = types.StringValue(string(jsonBytes))
        } else {
            data.OrganizationName = types.StringNull()
        }
    } else if val, ok := item["organizationName"].(string); ok {
        data.OrganizationName = types.StringValue(val)
    } else {
        data.OrganizationName = types.StringNull()
    }
    if obj, ok := item["affectedName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AffectedName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AffectedName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AffectedName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AffectedName = types.StringValue(string(jsonBytes))
        } else {
            data.AffectedName = types.StringNull()
        }
    } else if val, ok := item["affectedName"].(string); ok {
        data.AffectedName = types.StringValue(val)
    } else {
        data.AffectedName = types.StringNull()
    }
    if obj, ok := item["subject"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Subject = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Subject = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Subject = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Subject = types.StringValue(string(jsonBytes))
        } else {
            data.Subject = types.StringNull()
        }
    } else if val, ok := item["subject"].(string); ok {
        data.Subject = types.StringValue(val)
    } else {
        data.Subject = types.StringNull()
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
    if obj, ok := item["outcome"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Outcome = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Outcome = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Outcome = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Outcome = types.StringValue(string(jsonBytes))
        } else {
            data.Outcome = types.StringNull()
        }
    } else if val, ok := item["outcome"].(string); ok {
        data.Outcome = types.StringValue(val)
    } else {
        data.Outcome = types.StringNull()
    }
    if val, ok := item["pagedOnCall"].(bool); ok {
        data.PagedOnCall = types.BoolValue(val)
    } else {
        data.PagedOnCall = types.BoolNull()
    }
    if obj, ok := item["lastEventType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastEventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastEventType = types.StringValue(string(jsonBytes))
        } else {
            data.LastEventType = types.StringNull()
        }
    } else if val, ok := item["lastEventType"].(string); ok {
        data.LastEventType = types.StringValue(val)
    } else {
        data.LastEventType = types.StringNull()
    }
    if obj, ok := item["lastEventReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEventReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastEventReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastEventReceivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastEventReceivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastEventReceivedAt = types.StringNull()
        }
    } else if val, ok := item["lastEventReceivedAt"].(string); ok {
        data.LastEventReceivedAt = types.StringValue(val)
    } else {
        data.LastEventReceivedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
