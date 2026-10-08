package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OnCallPolicyScheduleDataSource{}

func NewOnCallPolicyScheduleDataSource() datasource.DataSource {
    return &OnCallPolicyScheduleDataSource{}
}

// OnCallPolicyScheduleDataSource defines the data source implementation.
type OnCallPolicyScheduleDataSource struct {
    client *Client
}

// OnCallPolicyScheduleDataSourceModel describes the data source data model.
type OnCallPolicyScheduleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Labels types.Set `tfsdk:"labels"`
    Description types.String `tfsdk:"description"`
    Timezone types.String `tfsdk:"timezone"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    CurrentUserIdOnRoster types.String `tfsdk:"current_user_id_on_roster"`
    NextUserIdOnRoster types.String `tfsdk:"next_user_id_on_roster"`
    RosterHandoffAt types.String `tfsdk:"roster_handoff_at"`
    RosterNextHandoffAt types.String `tfsdk:"roster_next_handoff_at"`
    RosterNextStartAt types.String `tfsdk:"roster_next_start_at"`
    RosterStartAt types.String `tfsdk:"roster_start_at"`
    ShiftConfigVersion types.Number `tfsdk:"shift_config_version"`
}

func (d *OnCallPolicyScheduleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_policy_schedule"
}

func (d *OnCallPolicyScheduleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage schedules and rotations for your on-call duty policy. Look up an existing on call policy schedule by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `current_user_id_on_roster`, ...): each one set must match, and exactly one on call policy schedule may match them all.",

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
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "timezone": schema.StringAttribute{
                MarkdownDescription: "IANA timezone this schedule's restriction and hand-off wall-clock times are interpreted in. When empty, times are interpreted in the server's local timezone (legacy behavior).",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "current_user_id_on_roster": schema.StringAttribute{
                MarkdownDescription: "User ID who is currently on roster. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "next_user_id_on_roster": schema.StringAttribute{
                MarkdownDescription: "Next ID who is currently on roster. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "roster_handoff_at": schema.StringAttribute{
                MarkdownDescription: "When does the roster handoff occur for this schedule for the current user?",
                Computed: true,
            },
            "roster_next_handoff_at": schema.StringAttribute{
                MarkdownDescription: "When does the next roster handoff occur for this schedule for the next user?",
                Computed: true,
            },
            "roster_next_start_at": schema.StringAttribute{
                MarkdownDescription: "When does the next event start for this schedule for the next user?",
                Computed: true,
            },
            "roster_start_at": schema.StringAttribute{
                MarkdownDescription: "When does the current event start for this schedule for the current user?",
                Computed: true,
            },
            "shift_config_version": schema.NumberAttribute{
                MarkdownDescription: "Incremented whenever the schedule's layers, members, overrides or policy attachments change. Used as the calendar feed SEQUENCE.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *OnCallPolicyScheduleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallPolicyScheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallPolicyScheduleDataSourceModel

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
    if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
        filters["timezone"] = data.Timezone.ValueString()
        filterNames = append(filterNames, "timezone = "+fmt.Sprintf("%q", data.Timezone.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.CurrentUserIdOnRoster.IsNull() && !data.CurrentUserIdOnRoster.IsUnknown() {
        filters["currentUserIdOnRoster"] = data.CurrentUserIdOnRoster.ValueString()
        filterNames = append(filterNames, "current_user_id_on_roster = "+fmt.Sprintf("%q", data.CurrentUserIdOnRoster.ValueString()))
    }
    if !data.NextUserIdOnRoster.IsNull() && !data.NextUserIdOnRoster.IsUnknown() {
        filters["nextUserIdOnRoster"] = data.NextUserIdOnRoster.ValueString()
        filterNames = append(filterNames, "next_user_id_on_roster = "+fmt.Sprintf("%q", data.NextUserIdOnRoster.ValueString()))
    }
    if !data.ShiftConfigVersion.IsNull() && !data.ShiftConfigVersion.IsUnknown() {
        filters["shiftConfigVersion"] = lookupNumber(data.ShiftConfigVersion)
        filterNames = append(filterNames, "shift_config_version = "+data.ShiftConfigVersion.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the on call policy schedule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the on call policy schedule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "labels": true,
        "description": true,
        "timezone": true,
        "slug": true,
        "createdByUserId": true,
        "currentUserIdOnRoster": true,
        "nextUserIdOnRoster": true,
        "rosterHandoffAt": true,
        "rosterNextHandoffAt": true,
        "rosterNextStartAt": true,
        "rosterStartAt": true,
        "shiftConfigVersion": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/on-call-duty-policy-schedule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_policy_schedule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call policy schedule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read on_call_policy_schedule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/on-call-duty-policy-schedule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list on_call_policy_schedule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list on_call_policy_schedule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No on call policy schedule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one on call policy schedule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for on_call_policy_schedule.")
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if obj, ok := item["currentUserIdOnRoster"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentUserIdOnRoster = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentUserIdOnRoster = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentUserIdOnRoster = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentUserIdOnRoster = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentUserIdOnRoster = types.StringNull()
        }
    } else if val, ok := item["currentUserIdOnRoster"].(string); ok {
        data.CurrentUserIdOnRoster = types.StringValue(val)
    } else {
        data.CurrentUserIdOnRoster = types.StringNull()
    }
    if obj, ok := item["nextUserIdOnRoster"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextUserIdOnRoster = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextUserIdOnRoster = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextUserIdOnRoster = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextUserIdOnRoster = types.StringValue(string(jsonBytes))
        } else {
            data.NextUserIdOnRoster = types.StringNull()
        }
    } else if val, ok := item["nextUserIdOnRoster"].(string); ok {
        data.NextUserIdOnRoster = types.StringValue(val)
    } else {
        data.NextUserIdOnRoster = types.StringNull()
    }
    if obj, ok := item["rosterHandoffAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RosterHandoffAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RosterHandoffAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RosterHandoffAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RosterHandoffAt = types.StringValue(string(jsonBytes))
        } else {
            data.RosterHandoffAt = types.StringNull()
        }
    } else if val, ok := item["rosterHandoffAt"].(string); ok {
        data.RosterHandoffAt = types.StringValue(val)
    } else {
        data.RosterHandoffAt = types.StringNull()
    }
    if obj, ok := item["rosterNextHandoffAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RosterNextHandoffAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RosterNextHandoffAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RosterNextHandoffAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RosterNextHandoffAt = types.StringValue(string(jsonBytes))
        } else {
            data.RosterNextHandoffAt = types.StringNull()
        }
    } else if val, ok := item["rosterNextHandoffAt"].(string); ok {
        data.RosterNextHandoffAt = types.StringValue(val)
    } else {
        data.RosterNextHandoffAt = types.StringNull()
    }
    if obj, ok := item["rosterNextStartAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RosterNextStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RosterNextStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RosterNextStartAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RosterNextStartAt = types.StringValue(string(jsonBytes))
        } else {
            data.RosterNextStartAt = types.StringNull()
        }
    } else if val, ok := item["rosterNextStartAt"].(string); ok {
        data.RosterNextStartAt = types.StringValue(val)
    } else {
        data.RosterNextStartAt = types.StringNull()
    }
    if obj, ok := item["rosterStartAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RosterStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RosterStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RosterStartAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RosterStartAt = types.StringValue(string(jsonBytes))
        } else {
            data.RosterStartAt = types.StringNull()
        }
    } else if val, ok := item["rosterStartAt"].(string); ok {
        data.RosterStartAt = types.StringValue(val)
    } else {
        data.RosterStartAt = types.StringNull()
    }
    if val, ok := item["shiftConfigVersion"].(float64); ok {
        data.ShiftConfigVersion = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["shiftConfigVersion"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ShiftConfigVersion = types.NumberValue(big.NewFloat(val))
        } else {
            data.ShiftConfigVersion = types.NumberNull()
        }
    } else {
        data.ShiftConfigVersion = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
