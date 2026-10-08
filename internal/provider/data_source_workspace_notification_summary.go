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
var _ datasource.DataSource = &WorkspaceNotificationSummaryDataSource{}

func NewWorkspaceNotificationSummaryDataSource() datasource.DataSource {
    return &WorkspaceNotificationSummaryDataSource{}
}

// WorkspaceNotificationSummaryDataSource defines the data source implementation.
type WorkspaceNotificationSummaryDataSource struct {
    client *Client
}

// WorkspaceNotificationSummaryDataSourceModel describes the data source data model.
type WorkspaceNotificationSummaryDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    WorkspaceType types.String `tfsdk:"workspace_type"`
    SummaryType types.String `tfsdk:"summary_type"`
    RecurringInterval types.String `tfsdk:"recurring_interval"`
    NumberOfDaysOfData types.Number `tfsdk:"number_of_days_of_data"`
    SendFirstReportAt types.String `tfsdk:"send_first_report_at"`
    Timezone types.String `tfsdk:"timezone"`
    ChannelNames types.String `tfsdk:"channel_names"`
    TeamName types.String `tfsdk:"team_name"`
    SummaryItems types.String `tfsdk:"summary_items"`
    Filters types.String `tfsdk:"filters"`
    FilterCondition types.String `tfsdk:"filter_condition"`
    NextSendAt types.String `tfsdk:"next_send_at"`
    LastSentAt types.String `tfsdk:"last_sent_at"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *WorkspaceNotificationSummaryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace_notification_summary"
}

func (d *WorkspaceNotificationSummaryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams Look up an existing workspace notification summary by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one workspace notification summary may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of the Summary Rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of the Summary Rule.",
                Optional: true,
                Computed: true,
            },
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Type of Workspace - Slack, Microsoft Teams, etc.",
                Optional: true,
                Computed: true,
            },
            "summary_type": schema.StringAttribute{
                MarkdownDescription: "Type of summary - Incident, Alert, Incident Episode, or Alert Episode.",
                Optional: true,
                Computed: true,
            },
            "recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often should the summary be sent? A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "number_of_days_of_data": schema.NumberAttribute{
                MarkdownDescription: "How many days of data to include in the summary.",
                Optional: true,
                Computed: true,
            },
            "send_first_report_at": schema.StringAttribute{
                MarkdownDescription: "When should the first summary report be sent? Subsequent reports will follow the recurring interval from this date.",
                Computed: true,
            },
            "timezone": schema.StringAttribute{
                MarkdownDescription: "The IANA time zone the summary's schedule is read in, such as Europe/Berlin or America/New_York. The summary goes out at the same time of day there all year, also after the clocks change for daylight saving time. Left out when the summary is created, it is the time zone in the creator's profile, or UTC when no person creates it (an API key or a workflow). A summary without one is read in UTC.",
                Optional: true,
                Computed: true,
            },
            "channel_names": schema.StringAttribute{
                MarkdownDescription: "List of channel names to post the summary to. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "team_name": schema.StringAttribute{
                MarkdownDescription: "Microsoft Teams team name (only for Microsoft Teams).",
                Optional: true,
                Computed: true,
            },
            "summary_items": schema.StringAttribute{
                MarkdownDescription: "Checklist of items to include in the summary. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "Filter conditions for which items to include in the summary. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "filter_condition": schema.StringAttribute{
                MarkdownDescription: "How to combine filters - Any or All.",
                Optional: true,
                Computed: true,
            },
            "next_send_at": schema.StringAttribute{
                MarkdownDescription: "When the next summary should be sent.",
                Computed: true,
            },
            "last_sent_at": schema.StringAttribute{
                MarkdownDescription: "When the last summary was sent.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is this summary rule enabled?",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *WorkspaceNotificationSummaryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceNotificationSummaryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkspaceNotificationSummaryDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.WorkspaceType.IsNull() && !data.WorkspaceType.IsUnknown() {
        filters["workspaceType"] = data.WorkspaceType.ValueString()
        filterNames = append(filterNames, "workspace_type = "+fmt.Sprintf("%q", data.WorkspaceType.ValueString()))
    }
    if !data.SummaryType.IsNull() && !data.SummaryType.IsUnknown() {
        filters["summaryType"] = data.SummaryType.ValueString()
        filterNames = append(filterNames, "summary_type = "+fmt.Sprintf("%q", data.SummaryType.ValueString()))
    }
    if !data.NumberOfDaysOfData.IsNull() && !data.NumberOfDaysOfData.IsUnknown() {
        filters["numberOfDaysOfData"] = lookupNumber(data.NumberOfDaysOfData)
        filterNames = append(filterNames, "number_of_days_of_data = "+data.NumberOfDaysOfData.ValueBigFloat().String())
    }
    if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
        filters["timezone"] = data.Timezone.ValueString()
        filterNames = append(filterNames, "timezone = "+fmt.Sprintf("%q", data.Timezone.ValueString()))
    }
    if !data.TeamName.IsNull() && !data.TeamName.IsUnknown() {
        filters["teamName"] = data.TeamName.ValueString()
        filterNames = append(filterNames, "team_name = "+fmt.Sprintf("%q", data.TeamName.ValueString()))
    }
    if !data.FilterCondition.IsNull() && !data.FilterCondition.IsUnknown() {
        filters["filterCondition"] = data.FilterCondition.ValueString()
        filterNames = append(filterNames, "filter_condition = "+fmt.Sprintf("%q", data.FilterCondition.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the workspace notification summary up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the workspace notification summary up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "workspaceType": true,
        "summaryType": true,
        "recurringInterval": true,
        "numberOfDaysOfData": true,
        "sendFirstReportAt": true,
        "timezone": true,
        "channelNames": true,
        "teamName": true,
        "summaryItems": true,
        "filters": true,
        "filterCondition": true,
        "nextSendAt": true,
        "lastSentAt": true,
        "isEnabled": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workspace-notification-summary/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace_notification_summary, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workspace notification summary found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workspace_notification_summary: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workspace-notification-summary/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workspace_notification_summary, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workspace_notification_summary: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workspace notification summary matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workspace notification summary matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workspace_notification_summary.")
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
    if obj, ok := item["workspaceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkspaceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkspaceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkspaceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkspaceType = types.StringValue(string(jsonBytes))
        } else {
            data.WorkspaceType = types.StringNull()
        }
    } else if val, ok := item["workspaceType"].(string); ok {
        data.WorkspaceType = types.StringValue(val)
    } else {
        data.WorkspaceType = types.StringNull()
    }
    if obj, ok := item["summaryType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SummaryType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SummaryType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SummaryType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SummaryType = types.StringValue(string(jsonBytes))
        } else {
            data.SummaryType = types.StringNull()
        }
    } else if val, ok := item["summaryType"].(string); ok {
        data.SummaryType = types.StringValue(val)
    } else {
        data.SummaryType = types.StringNull()
    }
    if obj, ok := item["recurringInterval"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.RecurringInterval = types.StringNull()
        }
    } else if val, ok := item["recurringInterval"].(string); ok {
        data.RecurringInterval = types.StringValue(val)
    } else {
        data.RecurringInterval = types.StringNull()
    }
    if val, ok := item["numberOfDaysOfData"].(float64); ok {
        data.NumberOfDaysOfData = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["numberOfDaysOfData"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.NumberOfDaysOfData = types.NumberValue(big.NewFloat(val))
        } else {
            data.NumberOfDaysOfData = types.NumberNull()
        }
    } else {
        data.NumberOfDaysOfData = types.NumberNull()
    }
    if obj, ok := item["sendFirstReportAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendFirstReportAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SendFirstReportAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SendFirstReportAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SendFirstReportAt = types.StringValue(string(jsonBytes))
        } else {
            data.SendFirstReportAt = types.StringNull()
        }
    } else if val, ok := item["sendFirstReportAt"].(string); ok {
        data.SendFirstReportAt = types.StringValue(val)
    } else {
        data.SendFirstReportAt = types.StringNull()
    }
    if obj, ok := item["timezone"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Timezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Timezone = types.StringValue(string(jsonBytes))
        } else {
            data.Timezone = types.StringNull()
        }
    } else if val, ok := item["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    } else {
        data.Timezone = types.StringNull()
    }
    if obj, ok := item["channelNames"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChannelNames = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChannelNames = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChannelNames = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChannelNames = types.StringValue(string(jsonBytes))
        } else {
            data.ChannelNames = types.StringNull()
        }
    } else if val, ok := item["channelNames"].(string); ok {
        data.ChannelNames = types.StringValue(val)
    } else {
        data.ChannelNames = types.StringNull()
    }
    if obj, ok := item["teamName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TeamName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TeamName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TeamName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TeamName = types.StringValue(string(jsonBytes))
        } else {
            data.TeamName = types.StringNull()
        }
    } else if val, ok := item["teamName"].(string); ok {
        data.TeamName = types.StringValue(val)
    } else {
        data.TeamName = types.StringNull()
    }
    if obj, ok := item["summaryItems"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SummaryItems = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SummaryItems = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SummaryItems = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SummaryItems = types.StringValue(string(jsonBytes))
        } else {
            data.SummaryItems = types.StringNull()
        }
    } else if val, ok := item["summaryItems"].(string); ok {
        data.SummaryItems = types.StringValue(val)
    } else {
        data.SummaryItems = types.StringNull()
    }
    if obj, ok := item["filters"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Filters = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Filters = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Filters = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Filters = types.StringValue(string(jsonBytes))
        } else {
            data.Filters = types.StringNull()
        }
    } else if val, ok := item["filters"].(string); ok {
        data.Filters = types.StringValue(val)
    } else {
        data.Filters = types.StringNull()
    }
    if obj, ok := item["filterCondition"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FilterCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FilterCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FilterCondition = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FilterCondition = types.StringValue(string(jsonBytes))
        } else {
            data.FilterCondition = types.StringNull()
        }
    } else if val, ok := item["filterCondition"].(string); ok {
        data.FilterCondition = types.StringValue(val)
    } else {
        data.FilterCondition = types.StringNull()
    }
    if obj, ok := item["nextSendAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextSendAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextSendAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextSendAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextSendAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextSendAt = types.StringNull()
        }
    } else if val, ok := item["nextSendAt"].(string); ok {
        data.NextSendAt = types.StringValue(val)
    } else {
        data.NextSendAt = types.StringNull()
    }
    if obj, ok := item["lastSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSentAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSentAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSentAt = types.StringNull()
        }
    } else if val, ok := item["lastSentAt"].(string); ok {
        data.LastSentAt = types.StringValue(val)
    } else {
        data.LastSentAt = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
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
