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
var _ datasource.DataSource = &HuntressConnectionDataSource{}

func NewHuntressConnectionDataSource() datasource.DataSource {
    return &HuntressConnectionDataSource{}
}

// HuntressConnectionDataSource defines the data source implementation.
type HuntressConnectionDataSource struct {
    client *Client
}

// HuntressConnectionDataSourceModel describes the data source data model.
type HuntressConnectionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    IsSigningSecretSet types.Bool `tfsdk:"is_signing_secret_set"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    PageOnCallFor types.String `tfsdk:"page_on_call_for"`
    CriticalIncidentSeverityId types.String `tfsdk:"critical_incident_severity_id"`
    HighIncidentSeverityId types.String `tfsdk:"high_incident_severity_id"`
    LowIncidentSeverityId types.String `tfsdk:"low_incident_severity_id"`
    WatchedOrganizations types.String `tfsdk:"watched_organizations"`
    Labels types.Set `tfsdk:"labels"`
    ResolveIncidentWhenReportCloses types.Bool `tfsdk:"resolve_incident_when_report_closes"`
    LastEventReceivedAt types.String `tfsdk:"last_event_received_at"`
    LastEventType types.String `tfsdk:"last_event_type"`
    LastError types.String `tfsdk:"last_error"`
    LastErrorAt types.String `tfsdk:"last_error_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *HuntressConnectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_huntress_connection"
}

func (d *HuntressConnectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Huntress webhook endpoints. Every Huntress incident report opens an incident that pages on-call, and closing the report in Huntress resolves it. Look up an existing huntress connection by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `critical_incident_severity_id`, ...): each one set must match, and exactly one huntress connection may match them all.",

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
                MarkdownDescription: "ID of the project this connection belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "A name for this connection, such as the Huntress account it receives from.",
                Optional: true,
                Computed: true,
            },
            "is_signing_secret_set": schema.BoolAttribute{
                MarkdownDescription: "Whether a signing secret is saved. Requests are only accepted once it is. Set from the signing secret itself; a value sent for it is ignored.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "The on-call policies paged for an incident report at or above the Page On-Call For severity. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "page_on_call_for": schema.StringAttribute{
                MarkdownDescription: "The lowest Huntress severity that pages the on-call policies: critical (critical reports only), high (high and critical reports) or low (every report). Every report opens an incident either way.",
                Optional: true,
                Computed: true,
            },
            "critical_incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity critical reports open at. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "high_incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity high reports open at. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "low_incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity low reports open at. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "watched_organizations": schema.StringAttribute{
                MarkdownDescription: "The Huntress organizations this connection opens incidents for, one per line, by name or by id. Empty watches every organization.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels added to every incident this connection opens, next to the label named after the report's organization. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "resolve_incident_when_report_closes": schema.BoolAttribute{
                MarkdownDescription: "Resolve the incident when the report is closed or dismissed in Huntress. When off, a note on the incident says so instead.",
                Optional: true,
                Computed: true,
            },
            "last_event_received_at": schema.StringAttribute{
                MarkdownDescription: "When a request signed with the signing secret last arrived. Empty until Huntress sends the first one.",
                Computed: true,
            },
            "last_event_type": schema.StringAttribute{
                MarkdownDescription: "The type of the last event received, such as incident_report.created.",
                Optional: true,
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "Why the last request was refused or could not be handled, if it was. Cleared by the next request that is handled.",
                Optional: true,
                Computed: true,
            },
            "last_error_at": schema.StringAttribute{
                MarkdownDescription: "When the last error happened.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *HuntressConnectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HuntressConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HuntressConnectionDataSourceModel

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
    if !data.IsSigningSecretSet.IsNull() && !data.IsSigningSecretSet.IsUnknown() {
        filters["isSigningSecretSet"] = data.IsSigningSecretSet.ValueBool()
        filterNames = append(filterNames, "is_signing_secret_set = "+fmt.Sprintf("%t", data.IsSigningSecretSet.ValueBool()))
    }
    if !data.PageOnCallFor.IsNull() && !data.PageOnCallFor.IsUnknown() {
        filters["pageOnCallFor"] = data.PageOnCallFor.ValueString()
        filterNames = append(filterNames, "page_on_call_for = "+fmt.Sprintf("%q", data.PageOnCallFor.ValueString()))
    }
    if !data.CriticalIncidentSeverityId.IsNull() && !data.CriticalIncidentSeverityId.IsUnknown() {
        filters["criticalIncidentSeverityId"] = data.CriticalIncidentSeverityId.ValueString()
        filterNames = append(filterNames, "critical_incident_severity_id = "+fmt.Sprintf("%q", data.CriticalIncidentSeverityId.ValueString()))
    }
    if !data.HighIncidentSeverityId.IsNull() && !data.HighIncidentSeverityId.IsUnknown() {
        filters["highIncidentSeverityId"] = data.HighIncidentSeverityId.ValueString()
        filterNames = append(filterNames, "high_incident_severity_id = "+fmt.Sprintf("%q", data.HighIncidentSeverityId.ValueString()))
    }
    if !data.LowIncidentSeverityId.IsNull() && !data.LowIncidentSeverityId.IsUnknown() {
        filters["lowIncidentSeverityId"] = data.LowIncidentSeverityId.ValueString()
        filterNames = append(filterNames, "low_incident_severity_id = "+fmt.Sprintf("%q", data.LowIncidentSeverityId.ValueString()))
    }
    if !data.WatchedOrganizations.IsNull() && !data.WatchedOrganizations.IsUnknown() {
        filters["watchedOrganizations"] = data.WatchedOrganizations.ValueString()
        filterNames = append(filterNames, "watched_organizations = "+fmt.Sprintf("%q", data.WatchedOrganizations.ValueString()))
    }
    if !data.ResolveIncidentWhenReportCloses.IsNull() && !data.ResolveIncidentWhenReportCloses.IsUnknown() {
        filters["resolveIncidentWhenReportCloses"] = data.ResolveIncidentWhenReportCloses.ValueBool()
        filterNames = append(filterNames, "resolve_incident_when_report_closes = "+fmt.Sprintf("%t", data.ResolveIncidentWhenReportCloses.ValueBool()))
    }
    if !data.LastEventType.IsNull() && !data.LastEventType.IsUnknown() {
        filters["lastEventType"] = data.LastEventType.ValueString()
        filterNames = append(filterNames, "last_event_type = "+fmt.Sprintf("%q", data.LastEventType.ValueString()))
    }
    if !data.LastError.IsNull() && !data.LastError.IsUnknown() {
        filters["lastError"] = data.LastError.ValueString()
        filterNames = append(filterNames, "last_error = "+fmt.Sprintf("%q", data.LastError.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the huntress connection up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the huntress connection up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "isSigningSecretSet": true,
        "onCallDutyPolicies": true,
        "pageOnCallFor": true,
        "criticalIncidentSeverityId": true,
        "highIncidentSeverityId": true,
        "lowIncidentSeverityId": true,
        "watchedOrganizations": true,
        "labels": true,
        "resolveIncidentWhenReportCloses": true,
        "lastEventReceivedAt": true,
        "lastEventType": true,
        "lastError": true,
        "lastErrorAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/huntress-connection/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read huntress_connection, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No huntress connection found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read huntress_connection: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/huntress-connection/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list huntress_connection, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list huntress_connection: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No huntress connection matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one huntress connection matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for huntress_connection.")
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
    if val, ok := item["isSigningSecretSet"].(bool); ok {
        data.IsSigningSecretSet = types.BoolValue(val)
    } else {
        data.IsSigningSecretSet = types.BoolNull()
    }
    if val, ok := item["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OnCallDutyPolicies = types.SetNull(types.StringType)
    }
    if obj, ok := item["pageOnCallFor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PageOnCallFor = types.StringValue(string(jsonBytes))
        } else {
            data.PageOnCallFor = types.StringNull()
        }
    } else if val, ok := item["pageOnCallFor"].(string); ok {
        data.PageOnCallFor = types.StringValue(val)
    } else {
        data.PageOnCallFor = types.StringNull()
    }
    if obj, ok := item["criticalIncidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.CriticalIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["criticalIncidentSeverityId"].(string); ok {
        data.CriticalIncidentSeverityId = types.StringValue(val)
    } else {
        data.CriticalIncidentSeverityId = types.StringNull()
    }
    if obj, ok := item["highIncidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.HighIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["highIncidentSeverityId"].(string); ok {
        data.HighIncidentSeverityId = types.StringValue(val)
    } else {
        data.HighIncidentSeverityId = types.StringNull()
    }
    if obj, ok := item["lowIncidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.LowIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["lowIncidentSeverityId"].(string); ok {
        data.LowIncidentSeverityId = types.StringValue(val)
    } else {
        data.LowIncidentSeverityId = types.StringNull()
    }
    if obj, ok := item["watchedOrganizations"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WatchedOrganizations = types.StringValue(string(jsonBytes))
        } else {
            data.WatchedOrganizations = types.StringNull()
        }
    } else if val, ok := item["watchedOrganizations"].(string); ok {
        data.WatchedOrganizations = types.StringValue(val)
    } else {
        data.WatchedOrganizations = types.StringNull()
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
    if val, ok := item["resolveIncidentWhenReportCloses"].(bool); ok {
        data.ResolveIncidentWhenReportCloses = types.BoolValue(val)
    } else {
        data.ResolveIncidentWhenReportCloses = types.BoolNull()
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
    if obj, ok := item["lastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := item["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
    }
    if obj, ok := item["lastErrorAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastErrorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastErrorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastErrorAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastErrorAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastErrorAt = types.StringNull()
        }
    } else if val, ok := item["lastErrorAt"].(string); ok {
        data.LastErrorAt = types.StringValue(val)
    } else {
        data.LastErrorAt = types.StringNull()
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
