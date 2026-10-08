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
var _ datasource.DataSource = &IncidentEpisodeMemberDataSource{}

func NewIncidentEpisodeMemberDataSource() datasource.DataSource {
    return &IncidentEpisodeMemberDataSource{}
}

// IncidentEpisodeMemberDataSource defines the data source implementation.
type IncidentEpisodeMemberDataSource struct {
    client *Client
}

// IncidentEpisodeMemberDataSourceModel describes the data source data model.
type IncidentEpisodeMemberDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IncidentEpisodeId types.String `tfsdk:"incident_episode_id"`
    IncidentId types.String `tfsdk:"incident_id"`
    AddedAt types.String `tfsdk:"added_at"`
    AddedBy types.String `tfsdk:"added_by"`
    AddedByUserId types.String `tfsdk:"added_by_user_id"`
    MatchedRuleId types.String `tfsdk:"matched_rule_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsOwnerNotifiedOfIncidentAdded types.Bool `tfsdk:"is_owner_notified_of_incident_added"`
}

func (d *IncidentEpisodeMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_episode_member"
}

func (d *IncidentEpisodeMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Link between incidents and episodes Look up an existing incident episode member by `id`, or by any of its other arguments (`added_by`, `added_by_user_id`, `created_by_user_id`, ...): each one set must match, and exactly one incident episode member may match them all.",

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
            "incident_episode_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Incident Episode that this incident belongs to. The ID of a `oneuptime_incident_episode`.",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Incident that is a member of this episode. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "added_at": schema.StringAttribute{
                MarkdownDescription: "When this incident was added to the episode.",
                Computed: true,
            },
            "added_by": schema.StringAttribute{
                MarkdownDescription: "How this incident was added to the episode (rule, manual, or api).",
                Optional: true,
                Computed: true,
            },
            "added_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who manually added this incident to the episode. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "matched_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the grouping rule that matched this incident. The ID of a `oneuptime_incident_grouping_rule`.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_owner_notified_of_incident_added": schema.BoolAttribute{
                MarkdownDescription: "Has the owner been notified that this incident was added to the episode?",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentEpisodeMemberDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentEpisodeMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentEpisodeMemberDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IncidentEpisodeId.IsNull() && !data.IncidentEpisodeId.IsUnknown() {
        filters["incidentEpisodeId"] = data.IncidentEpisodeId.ValueString()
        filterNames = append(filterNames, "incident_episode_id = "+fmt.Sprintf("%q", data.IncidentEpisodeId.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.AddedBy.IsNull() && !data.AddedBy.IsUnknown() {
        filters["addedBy"] = data.AddedBy.ValueString()
        filterNames = append(filterNames, "added_by = "+fmt.Sprintf("%q", data.AddedBy.ValueString()))
    }
    if !data.AddedByUserId.IsNull() && !data.AddedByUserId.IsUnknown() {
        filters["addedByUserId"] = data.AddedByUserId.ValueString()
        filterNames = append(filterNames, "added_by_user_id = "+fmt.Sprintf("%q", data.AddedByUserId.ValueString()))
    }
    if !data.MatchedRuleId.IsNull() && !data.MatchedRuleId.IsUnknown() {
        filters["matchedRuleId"] = data.MatchedRuleId.ValueString()
        filterNames = append(filterNames, "matched_rule_id = "+fmt.Sprintf("%q", data.MatchedRuleId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsOwnerNotifiedOfIncidentAdded.IsNull() && !data.IsOwnerNotifiedOfIncidentAdded.IsUnknown() {
        filters["isOwnerNotifiedOfIncidentAdded"] = data.IsOwnerNotifiedOfIncidentAdded.ValueBool()
        filterNames = append(filterNames, "is_owner_notified_of_incident_added = "+fmt.Sprintf("%t", data.IsOwnerNotifiedOfIncidentAdded.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident episode member up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident episode member up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incidentEpisodeId": true,
        "incidentId": true,
        "addedAt": true,
        "addedBy": true,
        "addedByUserId": true,
        "matchedRuleId": true,
        "createdByUserId": true,
        "isOwnerNotifiedOfIncidentAdded": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-episode-member/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_episode_member, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident episode member found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_episode_member: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-episode-member/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_episode_member, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_episode_member: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident episode member matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident episode member matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_episode_member.")
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
    if obj, ok := item["incidentEpisodeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentEpisodeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentEpisodeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentEpisodeId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeId = types.StringNull()
        }
    } else if val, ok := item["incidentEpisodeId"].(string); ok {
        data.IncidentEpisodeId = types.StringValue(val)
    } else {
        data.IncidentEpisodeId = types.StringNull()
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
    if obj, ok := item["addedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AddedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AddedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AddedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AddedAt = types.StringNull()
        }
    } else if val, ok := item["addedAt"].(string); ok {
        data.AddedAt = types.StringValue(val)
    } else {
        data.AddedAt = types.StringNull()
    }
    if obj, ok := item["addedBy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AddedBy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AddedBy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AddedBy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AddedBy = types.StringValue(string(jsonBytes))
        } else {
            data.AddedBy = types.StringNull()
        }
    } else if val, ok := item["addedBy"].(string); ok {
        data.AddedBy = types.StringValue(val)
    } else {
        data.AddedBy = types.StringNull()
    }
    if obj, ok := item["addedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AddedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AddedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AddedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AddedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AddedByUserId = types.StringNull()
        }
    } else if val, ok := item["addedByUserId"].(string); ok {
        data.AddedByUserId = types.StringValue(val)
    } else {
        data.AddedByUserId = types.StringNull()
    }
    if obj, ok := item["matchedRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MatchedRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MatchedRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MatchedRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MatchedRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.MatchedRuleId = types.StringNull()
        }
    } else if val, ok := item["matchedRuleId"].(string); ok {
        data.MatchedRuleId = types.StringValue(val)
    } else {
        data.MatchedRuleId = types.StringNull()
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
    if val, ok := item["isOwnerNotifiedOfIncidentAdded"].(bool); ok {
        data.IsOwnerNotifiedOfIncidentAdded = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfIncidentAdded = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
