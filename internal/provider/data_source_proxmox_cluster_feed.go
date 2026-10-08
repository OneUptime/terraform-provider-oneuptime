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
var _ datasource.DataSource = &ProxmoxClusterFeedDataSource{}

func NewProxmoxClusterFeedDataSource() datasource.DataSource {
    return &ProxmoxClusterFeedDataSource{}
}

// ProxmoxClusterFeedDataSource defines the data source implementation.
type ProxmoxClusterFeedDataSource struct {
    client *Client
}

// ProxmoxClusterFeedDataSourceModel describes the data source data model.
type ProxmoxClusterFeedDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ProxmoxClusterId types.String `tfsdk:"proxmox_cluster_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    FeedInfoInMarkdown types.String `tfsdk:"feed_info_in_markdown"`
    MoreInformationInMarkdown types.String `tfsdk:"more_information_in_markdown"`
    ProxmoxClusterFeedEventType types.String `tfsdk:"proxmox_cluster_feed_event_type"`
    DisplayColor types.String `tfsdk:"display_color"`
    UserId types.String `tfsdk:"user_id"`
    PostedAt types.String `tfsdk:"posted_at"`
}

func (d *ProxmoxClusterFeedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_proxmox_cluster_feed"
}

func (d *ProxmoxClusterFeedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Log of everything that happened to this Proxmox cluster - creation, updates, owner changes and the rules that made them. Look up an existing proxmox cluster feed by `id`, or by any of its other arguments (`created_by_user_id`, `feed_info_in_markdown`, `more_information_in_markdown`, ...): each one set must match, and exactly one proxmox cluster feed may match them all.",

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
            "proxmox_cluster_id": schema.StringAttribute{
                MarkdownDescription: "Relation to Proxmox Cluster ID in which this resource belongs. The ID of a `oneuptime_proxmox_cluster`.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "feed_info_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Log of the Proxmox cluster change in Markdown.",
                Optional: true,
                Computed: true,
            },
            "more_information_in_markdown": schema.StringAttribute{
                MarkdownDescription: "More information in Markdown.",
                Optional: true,
                Computed: true,
            },
            "proxmox_cluster_feed_event_type": schema.StringAttribute{
                MarkdownDescription: "Proxmox Cluster Feed Event.",
                Optional: true,
                Computed: true,
            },
            "display_color": schema.StringAttribute{
                MarkdownDescription: "Display color for this feed item.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "Date and time when the feed was posted.",
                Computed: true,
            },
        },
    }
}

func (d *ProxmoxClusterFeedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProxmoxClusterFeedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProxmoxClusterFeedDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ProxmoxClusterId.IsNull() && !data.ProxmoxClusterId.IsUnknown() {
        filters["proxmoxClusterId"] = data.ProxmoxClusterId.ValueString()
        filterNames = append(filterNames, "proxmox_cluster_id = "+fmt.Sprintf("%q", data.ProxmoxClusterId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.FeedInfoInMarkdown.IsNull() && !data.FeedInfoInMarkdown.IsUnknown() {
        filters["feedInfoInMarkdown"] = data.FeedInfoInMarkdown.ValueString()
        filterNames = append(filterNames, "feed_info_in_markdown = "+fmt.Sprintf("%q", data.FeedInfoInMarkdown.ValueString()))
    }
    if !data.MoreInformationInMarkdown.IsNull() && !data.MoreInformationInMarkdown.IsUnknown() {
        filters["moreInformationInMarkdown"] = data.MoreInformationInMarkdown.ValueString()
        filterNames = append(filterNames, "more_information_in_markdown = "+fmt.Sprintf("%q", data.MoreInformationInMarkdown.ValueString()))
    }
    if !data.ProxmoxClusterFeedEventType.IsNull() && !data.ProxmoxClusterFeedEventType.IsUnknown() {
        filters["proxmoxClusterFeedEventType"] = data.ProxmoxClusterFeedEventType.ValueString()
        filterNames = append(filterNames, "proxmox_cluster_feed_event_type = "+fmt.Sprintf("%q", data.ProxmoxClusterFeedEventType.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the proxmox cluster feed up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the proxmox cluster feed up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "proxmoxClusterId": true,
        "createdByUserId": true,
        "feedInfoInMarkdown": true,
        "moreInformationInMarkdown": true,
        "proxmoxClusterFeedEventType": true,
        "displayColor": true,
        "userId": true,
        "postedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/proxmox-cluster-feed/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read proxmox_cluster_feed, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No proxmox cluster feed found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read proxmox_cluster_feed: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/proxmox-cluster-feed/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list proxmox_cluster_feed, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list proxmox_cluster_feed: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No proxmox cluster feed matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one proxmox cluster feed matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for proxmox_cluster_feed.")
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
    if obj, ok := item["proxmoxClusterId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProxmoxClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProxmoxClusterId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProxmoxClusterId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProxmoxClusterId = types.StringValue(string(jsonBytes))
        } else {
            data.ProxmoxClusterId = types.StringNull()
        }
    } else if val, ok := item["proxmoxClusterId"].(string); ok {
        data.ProxmoxClusterId = types.StringValue(val)
    } else {
        data.ProxmoxClusterId = types.StringNull()
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
    if obj, ok := item["feedInfoInMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FeedInfoInMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FeedInfoInMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FeedInfoInMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FeedInfoInMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.FeedInfoInMarkdown = types.StringNull()
        }
    } else if val, ok := item["feedInfoInMarkdown"].(string); ok {
        data.FeedInfoInMarkdown = types.StringValue(val)
    } else {
        data.FeedInfoInMarkdown = types.StringNull()
    }
    if obj, ok := item["moreInformationInMarkdown"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MoreInformationInMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MoreInformationInMarkdown = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MoreInformationInMarkdown = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MoreInformationInMarkdown = types.StringValue(string(jsonBytes))
        } else {
            data.MoreInformationInMarkdown = types.StringNull()
        }
    } else if val, ok := item["moreInformationInMarkdown"].(string); ok {
        data.MoreInformationInMarkdown = types.StringValue(val)
    } else {
        data.MoreInformationInMarkdown = types.StringNull()
    }
    if obj, ok := item["proxmoxClusterFeedEventType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProxmoxClusterFeedEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProxmoxClusterFeedEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProxmoxClusterFeedEventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProxmoxClusterFeedEventType = types.StringValue(string(jsonBytes))
        } else {
            data.ProxmoxClusterFeedEventType = types.StringNull()
        }
    } else if val, ok := item["proxmoxClusterFeedEventType"].(string); ok {
        data.ProxmoxClusterFeedEventType = types.StringValue(val)
    } else {
        data.ProxmoxClusterFeedEventType = types.StringNull()
    }
    if obj, ok := item["displayColor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisplayColor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DisplayColor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DisplayColor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DisplayColor = types.StringValue(string(jsonBytes))
        } else {
            data.DisplayColor = types.StringNull()
        }
    } else if val, ok := item["displayColor"].(string); ok {
        data.DisplayColor = types.StringValue(val)
    } else {
        data.DisplayColor = types.StringNull()
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
    if obj, ok := item["postedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostedAt = types.StringNull()
        }
    } else if val, ok := item["postedAt"].(string); ok {
        data.PostedAt = types.StringValue(val)
    } else {
        data.PostedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
