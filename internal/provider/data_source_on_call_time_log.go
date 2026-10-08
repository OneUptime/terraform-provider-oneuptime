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
var _ datasource.DataSource = &OnCallTimeLogDataSource{}

func NewOnCallTimeLogDataSource() datasource.DataSource {
    return &OnCallTimeLogDataSource{}
}

// OnCallTimeLogDataSource defines the data source implementation.
type OnCallTimeLogDataSource struct {
    client *Client
}

// OnCallTimeLogDataSourceModel describes the data source data model.
type OnCallTimeLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    TeamId types.String `tfsdk:"team_id"`
    MoreInfo types.String `tfsdk:"more_info"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    UserId types.String `tfsdk:"user_id"`
    StartsAt types.String `tfsdk:"starts_at"`
    EndsAt types.String `tfsdk:"ends_at"`
}

func (d *OnCallTimeLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_time_log"
}

func (d *OnCallTimeLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user. Look up an existing on call time log by `id`, or by any of its other arguments (`created_by_user_id`, `more_info`, `on_call_duty_policy_escalation_rule_id`, ...): each one set must match, and exactly one on call time log may match them all.",

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
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy where this escalation rule belongs.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_schedule_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy Schedule where this escalation rule belongs.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy Escalation Rule ID where this escalation rule belongs.",
                Optional: true,
                Computed: true,
            },
            "team_id": schema.StringAttribute{
                MarkdownDescription: "ID of your On-Call Policy Team ID where this escalation rule belongs.",
                Optional: true,
                Computed: true,
            },
            "more_info": schema.StringAttribute{
                MarkdownDescription: "More information about this log record.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "User ID for which this log belongs. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "When does this start?",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "When does this end?",
                Computed: true,
            },
        },
    }
}

func (d *OnCallTimeLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallTimeLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallTimeLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.OnCallDutyPolicyId.IsNull() && !data.OnCallDutyPolicyId.IsUnknown() {
        filters["onCallDutyPolicyId"] = data.OnCallDutyPolicyId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyId.ValueString()))
    }
    if !data.OnCallDutyPolicyScheduleId.IsNull() && !data.OnCallDutyPolicyScheduleId.IsUnknown() {
        filters["onCallDutyPolicyScheduleId"] = data.OnCallDutyPolicyScheduleId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_schedule_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyScheduleId.ValueString()))
    }
    if !data.OnCallDutyPolicyEscalationRuleId.IsNull() && !data.OnCallDutyPolicyEscalationRuleId.IsUnknown() {
        filters["onCallDutyPolicyEscalationRuleId"] = data.OnCallDutyPolicyEscalationRuleId.ValueString()
        filterNames = append(filterNames, "on_call_duty_policy_escalation_rule_id = "+fmt.Sprintf("%q", data.OnCallDutyPolicyEscalationRuleId.ValueString()))
    }
    if !data.TeamId.IsNull() && !data.TeamId.IsUnknown() {
        filters["teamId"] = data.TeamId.ValueString()
        filterNames = append(filterNames, "team_id = "+fmt.Sprintf("%q", data.TeamId.ValueString()))
    }
    if !data.MoreInfo.IsNull() && !data.MoreInfo.IsUnknown() {
        filters["moreInfo"] = data.MoreInfo.ValueString()
        filterNames = append(filterNames, "more_info = "+fmt.Sprintf("%q", data.MoreInfo.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the on call time log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the on call time log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "onCallDutyPolicyId": true,
        "onCallDutyPolicyScheduleId": true,
        "onCallDutyPolicyEscalationRuleId": true,
        "teamId": true,
        "moreInfo": true,
        "createdByUserId": true,
        "userId": true,
        "startsAt": true,
        "endsAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/on-call-duty-policy-time-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_time_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call time log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read on_call_time_log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/on-call-duty-policy-time-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list on_call_time_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list on_call_time_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call time log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one on call time log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for on_call_time_log.")
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
    if obj, ok := item["onCallDutyPolicyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyId"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyId = types.StringNull()
    }
    if obj, ok := item["onCallDutyPolicyScheduleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyScheduleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyScheduleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyScheduleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyScheduleId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyScheduleId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyScheduleId"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyScheduleId = types.StringNull()
    }
    if obj, ok := item["onCallDutyPolicyEscalationRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OnCallDutyPolicyEscalationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
        }
    } else if val, ok := item["onCallDutyPolicyEscalationRuleId"].(string); ok {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
    }
    if obj, ok := item["teamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TeamId = types.StringValue(string(jsonBytes))
        } else {
            data.TeamId = types.StringNull()
        }
    } else if val, ok := item["teamId"].(string); ok {
        data.TeamId = types.StringValue(val)
    } else {
        data.TeamId = types.StringNull()
    }
    if obj, ok := item["moreInfo"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MoreInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MoreInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MoreInfo = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MoreInfo = types.StringValue(string(jsonBytes))
        } else {
            data.MoreInfo = types.StringNull()
        }
    } else if val, ok := item["moreInfo"].(string); ok {
        data.MoreInfo = types.StringValue(val)
    } else {
        data.MoreInfo = types.StringNull()
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
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if obj, ok := item["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartsAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartsAt = types.StringNull()
        }
    } else if val, ok := item["startsAt"].(string); ok {
        data.StartsAt = types.StringValue(val)
    } else {
        data.StartsAt = types.StringNull()
    }
    if obj, ok := item["endsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndsAt = types.StringNull()
        }
    } else if val, ok := item["endsAt"].(string); ok {
        data.EndsAt = types.StringValue(val)
    } else {
        data.EndsAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
