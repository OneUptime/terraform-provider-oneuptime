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
var _ datasource.DataSource = &OnCallUserGroupDataSource{}

func NewOnCallUserGroupDataSource() datasource.DataSource {
    return &OnCallUserGroupDataSource{}
}

// OnCallUserGroupDataSource defines the data source implementation.
type OnCallUserGroupDataSource struct {
    client *Client
}

// OnCallUserGroupDataSourceModel describes the data source data model.
type OnCallUserGroupDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    WorkspaceType types.String `tfsdk:"workspace_type"`
    WorkspaceGroupId types.String `tfsdk:"workspace_group_id"`
    WorkspaceTeamId types.String `tfsdk:"workspace_team_id"`
    Name types.String `tfsdk:"name"`
    Handle types.String `tfsdk:"handle"`
    OnCallDutyPolicySchedules types.Set `tfsdk:"on_call_duty_policy_schedules"`
    SyncState types.String `tfsdk:"sync_state"`
    SyncedAt types.String `tfsdk:"synced_at"`
    NextSyncAt types.String `tfsdk:"next_sync_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *OnCallUserGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_user_group"
}

func (d *OnCallUserGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A Slack user group or a Microsoft Teams tag whose members OneUptime keeps equal to whoever is on call right now in the on-call schedules it follows, so anyone can mention it to reach them. Look up an existing on call user group by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `handle`, ...): each one set must match, and exactly one on call user group may match them all.",

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
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Where the group is: Slack (a user group) or MicrosoftTeams (a tag in a team).",
                Optional: true,
                Computed: true,
            },
            "workspace_group_id": schema.StringAttribute{
                MarkdownDescription: "The group's id in the workspace: a Slack user group id, such as S0604QSJC, or a Microsoft Teams tag id.",
                Optional: true,
                Computed: true,
            },
            "workspace_team_id": schema.StringAttribute{
                MarkdownDescription: "Microsoft Teams only: the team the tag is in (its Microsoft 365 group id).",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The group's name in the workspace, such as Payments on-call. Read from Slack or Microsoft Teams when the group is linked.",
                Optional: true,
                Computed: true,
            },
            "handle": schema.StringAttribute{
                MarkdownDescription: "Slack only: the handle people type to mention the group, without the @, such as payments-oncall.",
                Optional: true,
                Computed: true,
            },
            "on_call_duty_policy_schedules": schema.SetAttribute{
                MarkdownDescription: "The on-call schedules the group follows: it holds whoever is on call in any of them right now. IDs of `oneuptime_on_call_policy_schedule` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "sync_state": schema.StringAttribute{
                MarkdownDescription: "What the last sync came to: who OneUptime put in the group, who is on call but could not be added and why, whether the group is turned off while nobody is on call, and what stopped the sync, if anything. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "synced_at": schema.StringAttribute{
                MarkdownDescription: "When the group last matched who was on call.",
                Computed: true,
            },
            "next_sync_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime looks at the group again, besides every hand-off: to undo changes made to it by hand, and to try again after the workspace refused a change.",
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

func (d *OnCallUserGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallUserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallUserGroupDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.WorkspaceType.IsNull() && !data.WorkspaceType.IsUnknown() {
        filters["workspaceType"] = data.WorkspaceType.ValueString()
        filterNames = append(filterNames, "workspace_type = "+fmt.Sprintf("%q", data.WorkspaceType.ValueString()))
    }
    if !data.WorkspaceGroupId.IsNull() && !data.WorkspaceGroupId.IsUnknown() {
        filters["workspaceGroupId"] = data.WorkspaceGroupId.ValueString()
        filterNames = append(filterNames, "workspace_group_id = "+fmt.Sprintf("%q", data.WorkspaceGroupId.ValueString()))
    }
    if !data.WorkspaceTeamId.IsNull() && !data.WorkspaceTeamId.IsUnknown() {
        filters["workspaceTeamId"] = data.WorkspaceTeamId.ValueString()
        filterNames = append(filterNames, "workspace_team_id = "+fmt.Sprintf("%q", data.WorkspaceTeamId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Handle.IsNull() && !data.Handle.IsUnknown() {
        filters["handle"] = data.Handle.ValueString()
        filterNames = append(filterNames, "handle = "+fmt.Sprintf("%q", data.Handle.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the on call user group up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the on call user group up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "workspaceType": true,
        "workspaceGroupId": true,
        "workspaceTeamId": true,
        "name": true,
        "handle": true,
        "onCallDutyPolicySchedules": true,
        "syncState": true,
        "syncedAt": true,
        "nextSyncAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workspace-user-group/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_user_group, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call user group found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read on_call_user_group: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workspace-user-group/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list on_call_user_group, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list on_call_user_group: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call user group matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one on call user group matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for on_call_user_group.")
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
    if obj, ok := item["workspaceGroupId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkspaceGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkspaceGroupId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkspaceGroupId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkspaceGroupId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkspaceGroupId = types.StringNull()
        }
    } else if val, ok := item["workspaceGroupId"].(string); ok {
        data.WorkspaceGroupId = types.StringValue(val)
    } else {
        data.WorkspaceGroupId = types.StringNull()
    }
    if obj, ok := item["workspaceTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkspaceTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkspaceTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkspaceTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkspaceTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkspaceTeamId = types.StringNull()
        }
    } else if val, ok := item["workspaceTeamId"].(string); ok {
        data.WorkspaceTeamId = types.StringValue(val)
    } else {
        data.WorkspaceTeamId = types.StringNull()
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
    if obj, ok := item["handle"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Handle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Handle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Handle = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Handle = types.StringValue(string(jsonBytes))
        } else {
            data.Handle = types.StringNull()
        }
    } else if val, ok := item["handle"].(string); ok {
        data.Handle = types.StringValue(val)
    } else {
        data.Handle = types.StringNull()
    }
    if val, ok := item["onCallDutyPolicySchedules"].([]interface{}); ok {
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
        data.OnCallDutyPolicySchedules = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OnCallDutyPolicySchedules = types.SetNull(types.StringType)
    }
    if obj, ok := item["syncState"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SyncState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SyncState = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SyncState = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SyncState = types.StringValue(string(jsonBytes))
        } else {
            data.SyncState = types.StringNull()
        }
    } else if val, ok := item["syncState"].(string); ok {
        data.SyncState = types.StringValue(val)
    } else {
        data.SyncState = types.StringNull()
    }
    if obj, ok := item["syncedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SyncedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SyncedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SyncedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SyncedAt = types.StringValue(string(jsonBytes))
        } else {
            data.SyncedAt = types.StringNull()
        }
    } else if val, ok := item["syncedAt"].(string); ok {
        data.SyncedAt = types.StringValue(val)
    } else {
        data.SyncedAt = types.StringNull()
    }
    if obj, ok := item["nextSyncAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextSyncAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextSyncAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextSyncAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextSyncAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextSyncAt = types.StringNull()
        }
    } else if val, ok := item["nextSyncAt"].(string); ok {
        data.NextSyncAt = types.StringValue(val)
    } else {
        data.NextSyncAt = types.StringNull()
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
