package shard26

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
var _ datasource.DataSource = &BlocklistFeedDataSource{}

func NewBlocklistFeedDataSource() datasource.DataSource {
    return &BlocklistFeedDataSource{}
}

// BlocklistFeedDataSource defines the data source implementation.
type BlocklistFeedDataSource struct {
    client *Client
}

// BlocklistFeedDataSourceModel describes the data source data model.
type BlocklistFeedDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Firewall types.String `tfsdk:"firewall"`
    Family types.String `tfsdk:"family"`
    AllowedSourceAddresses types.String `tfsdk:"allowed_source_addresses"`
    SecretHint types.String `tfsdk:"secret_hint"`
    SecretMadeAt types.String `tfsdk:"secret_made_at"`
    PreviousSecretExpiresAt types.String `tfsdk:"previous_secret_expires_at"`
    RevokedAt types.String `tfsdk:"revoked_at"`
    LastFetchedAt types.String `tfsdk:"last_fetched_at"`
    LastFetchSourceIp types.String `tfsdk:"last_fetch_source_ip"`
    LastFetchUserAgent types.String `tfsdk:"last_fetch_user_agent"`
    UsualFetchIntervalInSeconds types.Number `tfsdk:"usual_fetch_interval_in_seconds"`
    StoppedReadingToldAt types.String `tfsdk:"stopped_reading_told_at"`
    ServedCount types.Number `tfsdk:"served_count"`
    ServedConfirmedAt types.String `tfsdk:"served_confirmed_at"`
}

func (d *BlocklistFeedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_blocklist_feed"
}

func (d *BlocklistFeedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A secret URL a firewall or a proxy reads the project's blocklist from, on a schedule: Palo Alto, FortiGate, Check Point, Cisco, OPNsense, pfSense, nftables, nginx and others. Look up an existing blocklist feed by `id`, or by any of its other arguments (`name`, `allowed_source_addresses`, `created_by_actor_id`, ...): each one set must match, and exactly one blocklist feed may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this feed belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the feed is called where it is enforced: the firewall or proxy that reads it, such as Palo Alto edge firewall.",
                Optional: true,
                Computed: true,
            },
            "firewall": schema.StringAttribute{
                MarkdownDescription: "The firewall or proxy that reads the feed, whose setup steps its page opens on. Any of them can read any feed.",
                Optional: true,
                Computed: true,
            },
            "family": schema.StringAttribute{
                MarkdownDescription: "What the feed serves: IPv4 and IPv6 addresses (Both), or one of them, for a firewall that keeps them apart.",
                Optional: true,
                Computed: true,
            },
            "allowed_source_addresses": schema.StringAttribute{
                MarkdownDescription: "The only addresses and ranges the feed answers, one a line: the firewall's own. Empty to answer any address that has the secret.",
                Optional: true,
                Computed: true,
            },
            "secret_hint": schema.StringAttribute{
                MarkdownDescription: "The last four characters of the feed's secret, shown in its place so a person can tell which secret a firewall holds. Never enough to use it.",
                Optional: true,
                Computed: true,
            },
            "secret_made_at": schema.StringAttribute{
                MarkdownDescription: "When the feed's current secret was made.",
                Computed: true,
            },
            "previous_secret_expires_at": schema.StringAttribute{
                MarkdownDescription: "When the secret the last rotation replaced stops working: a day after the rotation, or at once when the feed's secret is revoked.",
                Computed: true,
            },
            "revoked_at": schema.StringAttribute{
                MarkdownDescription: "When the feed's secret was revoked. A revoked feed answers no firewall until a new secret is made.",
                Computed: true,
            },
            "last_fetched_at": schema.StringAttribute{
                MarkdownDescription: "When a firewall last fetched the feed.",
                Computed: true,
            },
            "last_fetch_source_ip": schema.StringAttribute{
                MarkdownDescription: "The address the feed was last fetched from.",
                Optional: true,
                Computed: true,
            },
            "last_fetch_user_agent": schema.StringAttribute{
                MarkdownDescription: "The user agent the feed was last fetched with: the firewall or tool that read it.",
                Optional: true,
                Computed: true,
            },
            "usual_fetch_interval_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How often the firewall usually fetches the feed: the median gap between its latest fetches. Three times it without a fetch, and the firewall has stopped reading the feed.",
                Optional: true,
                Computed: true,
            },
            "stopped_reading_told_at": schema.StringAttribute{
                MarkdownDescription: "When the Project Owners were told the firewall stopped reading the feed. Cleared when it reads it again, so they are told once each time.",
                Computed: true,
            },
            "served_count": schema.NumberAttribute{
                MarkdownDescription: "How many addresses and ranges the last list a firewall fetched held.",
                Optional: true,
                Computed: true,
            },
            "served_confirmed_at": schema.StringAttribute{
                MarkdownDescription: "When a firewall last fetched the feed and was found to hold its last list: the time its blocks were last confirmed there.",
                Computed: true,
            },
        },
    }
}

func (d *BlocklistFeedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BlocklistFeedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data BlocklistFeedDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Firewall.IsNull() && !data.Firewall.IsUnknown() {
        filters["firewall"] = data.Firewall.ValueString()
        filterNames = append(filterNames, "firewall = "+fmt.Sprintf("%q", data.Firewall.ValueString()))
    }
    if !data.Family.IsNull() && !data.Family.IsUnknown() {
        filters["family"] = data.Family.ValueString()
        filterNames = append(filterNames, "family = "+fmt.Sprintf("%q", data.Family.ValueString()))
    }
    if !data.AllowedSourceAddresses.IsNull() && !data.AllowedSourceAddresses.IsUnknown() {
        filters["allowedSourceAddresses"] = data.AllowedSourceAddresses.ValueString()
        filterNames = append(filterNames, "allowed_source_addresses = "+fmt.Sprintf("%q", data.AllowedSourceAddresses.ValueString()))
    }
    if !data.SecretHint.IsNull() && !data.SecretHint.IsUnknown() {
        filters["secretHint"] = data.SecretHint.ValueString()
        filterNames = append(filterNames, "secret_hint = "+fmt.Sprintf("%q", data.SecretHint.ValueString()))
    }
    if !data.LastFetchSourceIp.IsNull() && !data.LastFetchSourceIp.IsUnknown() {
        filters["lastFetchSourceIp"] = data.LastFetchSourceIp.ValueString()
        filterNames = append(filterNames, "last_fetch_source_ip = "+fmt.Sprintf("%q", data.LastFetchSourceIp.ValueString()))
    }
    if !data.LastFetchUserAgent.IsNull() && !data.LastFetchUserAgent.IsUnknown() {
        filters["lastFetchUserAgent"] = data.LastFetchUserAgent.ValueString()
        filterNames = append(filterNames, "last_fetch_user_agent = "+fmt.Sprintf("%q", data.LastFetchUserAgent.ValueString()))
    }
    if !data.UsualFetchIntervalInSeconds.IsNull() && !data.UsualFetchIntervalInSeconds.IsUnknown() {
        filters["usualFetchIntervalInSeconds"] = LookupNumber(data.UsualFetchIntervalInSeconds)
        filterNames = append(filterNames, "usual_fetch_interval_in_seconds = "+data.UsualFetchIntervalInSeconds.ValueBigFloat().String())
    }
    if !data.ServedCount.IsNull() && !data.ServedCount.IsUnknown() {
        filters["servedCount"] = LookupNumber(data.ServedCount)
        filterNames = append(filterNames, "served_count = "+data.ServedCount.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the blocklist feed up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the blocklist feed up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "name": true,
        "firewall": true,
        "family": true,
        "allowedSourceAddresses": true,
        "secretHint": true,
        "secretMadeAt": true,
        "previousSecretExpiresAt": true,
        "revokedAt": true,
        "lastFetchedAt": true,
        "lastFetchSourceIp": true,
        "lastFetchUserAgent": true,
        "usualFetchIntervalInSeconds": true,
        "stoppedReadingToldAt": true,
        "servedCount": true,
        "servedConfirmedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-blocklist-feed/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blocklist_feed, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No blocklist feed found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read blocklist_feed: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-blocklist-feed/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list blocklist_feed, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list blocklist_feed: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No blocklist feed matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one blocklist feed matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for blocklist_feed.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["firewall"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Firewall = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if arr, ok := item["firewall"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if val, ok := item["firewall"].(string); ok {
        data.Firewall = types.StringValue(val)
    } else {
        data.Firewall = types.StringNull()
    }
    if obj, ok := item["family"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Family = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if arr, ok := item["family"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if val, ok := item["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := item["allowedSourceAddresses"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if arr, ok := item["allowedSourceAddresses"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if val, ok := item["allowedSourceAddresses"].(string); ok {
        data.AllowedSourceAddresses = types.StringValue(val)
    } else {
        data.AllowedSourceAddresses = types.StringNull()
    }
    if obj, ok := item["secretHint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecretHint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if arr, ok := item["secretHint"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if val, ok := item["secretHint"].(string); ok {
        data.SecretHint = types.StringValue(val)
    } else {
        data.SecretHint = types.StringNull()
    }
    if obj, ok := item["secretMadeAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretMadeAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecretMadeAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecretMadeAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecretMadeAt = types.StringValue(string(jsonBytes))
        } else {
            data.SecretMadeAt = types.StringNull()
        }
    } else if arr, ok := item["secretMadeAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecretMadeAt = types.StringValue(string(jsonBytes))
        } else {
            data.SecretMadeAt = types.StringNull()
        }
    } else if val, ok := item["secretMadeAt"].(string); ok {
        data.SecretMadeAt = types.StringValue(val)
    } else {
        data.SecretMadeAt = types.StringNull()
    }
    if obj, ok := item["previousSecretExpiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PreviousSecretExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PreviousSecretExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PreviousSecretExpiresAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PreviousSecretExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.PreviousSecretExpiresAt = types.StringNull()
        }
    } else if arr, ok := item["previousSecretExpiresAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PreviousSecretExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.PreviousSecretExpiresAt = types.StringNull()
        }
    } else if val, ok := item["previousSecretExpiresAt"].(string); ok {
        data.PreviousSecretExpiresAt = types.StringValue(val)
    } else {
        data.PreviousSecretExpiresAt = types.StringNull()
    }
    if obj, ok := item["revokedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RevokedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RevokedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RevokedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RevokedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RevokedAt = types.StringNull()
        }
    } else if arr, ok := item["revokedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RevokedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RevokedAt = types.StringNull()
        }
    } else if val, ok := item["revokedAt"].(string); ok {
        data.RevokedAt = types.StringValue(val)
    } else {
        data.RevokedAt = types.StringNull()
    }
    if obj, ok := item["lastFetchedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastFetchedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastFetchedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastFetchedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchedAt = types.StringNull()
        }
    } else if arr, ok := item["lastFetchedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastFetchedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchedAt = types.StringNull()
        }
    } else if val, ok := item["lastFetchedAt"].(string); ok {
        data.LastFetchedAt = types.StringValue(val)
    } else {
        data.LastFetchedAt = types.StringNull()
    }
    if obj, ok := item["lastFetchSourceIp"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if arr, ok := item["lastFetchSourceIp"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if val, ok := item["lastFetchSourceIp"].(string); ok {
        data.LastFetchSourceIp = types.StringValue(val)
    } else {
        data.LastFetchSourceIp = types.StringNull()
    }
    if obj, ok := item["lastFetchUserAgent"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if arr, ok := item["lastFetchUserAgent"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if val, ok := item["lastFetchUserAgent"].(string); ok {
        data.LastFetchUserAgent = types.StringValue(val)
    } else {
        data.LastFetchUserAgent = types.StringNull()
    }
    data.UsualFetchIntervalInSeconds = NumberFromAPI(data.UsualFetchIntervalInSeconds, item["usualFetchIntervalInSeconds"])
    if obj, ok := item["stoppedReadingToldAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StoppedReadingToldAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StoppedReadingToldAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StoppedReadingToldAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StoppedReadingToldAt = types.StringValue(string(jsonBytes))
        } else {
            data.StoppedReadingToldAt = types.StringNull()
        }
    } else if arr, ok := item["stoppedReadingToldAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StoppedReadingToldAt = types.StringValue(string(jsonBytes))
        } else {
            data.StoppedReadingToldAt = types.StringNull()
        }
    } else if val, ok := item["stoppedReadingToldAt"].(string); ok {
        data.StoppedReadingToldAt = types.StringValue(val)
    } else {
        data.StoppedReadingToldAt = types.StringNull()
    }
    data.ServedCount = NumberFromAPI(data.ServedCount, item["servedCount"])
    if obj, ok := item["servedConfirmedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServedConfirmedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServedConfirmedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServedConfirmedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServedConfirmedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ServedConfirmedAt = types.StringNull()
        }
    } else if arr, ok := item["servedConfirmedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ServedConfirmedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ServedConfirmedAt = types.StringNull()
        }
    } else if val, ok := item["servedConfirmedAt"].(string); ok {
        data.ServedConfirmedAt = types.StringValue(val)
    } else {
        data.ServedConfirmedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
