package shard08

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
var _ datasource.DataSource = &CloudWafDataSource{}

func NewCloudWafDataSource() datasource.DataSource {
    return &CloudWafDataSource{}
}

// CloudWafDataSource defines the data source implementation.
type CloudWafDataSource struct {
    client *Client
}

// CloudWafDataSourceModel describes the data source data model.
type CloudWafDataSourceModel struct {
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
    CloudConnectionId types.String `tfsdk:"cloud_connection_id"`
    CloudResourceId types.String `tfsdk:"cloud_resource_id"`
    CloudProvider types.String `tfsdk:"cloud_provider"`
    WafId types.String `tfsdk:"waf_id"`
    Name types.String `tfsdk:"name"`
    Region types.String `tfsdk:"region"`
    ClientAddress types.String `tfsdk:"client_address"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    StatusReason types.String `tfsdk:"status_reason"`
    RuleDetails types.String `tfsdk:"rule_details"`
    CheckedAt types.String `tfsdk:"checked_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *CloudWafDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cloud_waf"
}

func (d *CloudWafDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The WAFs of connected cloud accounts that enforce the project's blocklist: AWS WAF web ACLs, Azure WAF policies and Google Cloud Armor policies, each holding OneUptime's own rules. Look up an existing cloud waf by `id`, or by any of its other arguments (`name`, `client_address`, `cloud_connection_id`, ...): each one set must match, and exactly one cloud waf may match them all.",

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
                MarkdownDescription: "ID of the project this WAF enforces the blocklist of. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "cloud_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the cloud connection OneUptime writes the WAF through. The ID of a `oneuptime_cloud_connection`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_cloud_connection"),
            },
            "cloud_resource_id": schema.StringAttribute{
                MarkdownDescription: "The WAF's row in the cloud inventory: a web ACL, Azure WAF policy or Cloud Armor policy of a connection whose blocking setting is on. The ID of a `oneuptime_cloud_resource`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_cloud_resource"),
            },
            "cloud_provider": schema.StringAttribute{
                MarkdownDescription: "The cloud the WAF is in: aws, azure or gcp.",
                Optional: true,
                Computed: true,
            },
            "waf_id": schema.StringAttribute{
                MarkdownDescription: "The WAF as its cloud names it: a web ACL's ARN, an Azure WAF policy's resource ID, or a Cloud Armor policy's name with its project.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The WAF's name: prod-alb.",
                Optional: true,
                Computed: true,
            },
            "region": schema.StringAttribute{
                MarkdownDescription: "The Region the WAF is in: us-east-1 for a CloudFront distribution's.",
                Optional: true,
                Computed: true,
            },
            "client_address": schema.StringAttribute{
                MarkdownDescription: "Where the WAF reads the client's address: the request's source address, or a header a proxy or CDN in front of it writes.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Where the WAF stands, as OneUptime last read it: Checking, Enforcing, Waiting for your rule, No free priority, Applying, Changed by hand, or Failed.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Why the WAF is not enforcing the blocklist, and what to do: empty when it is.",
                Optional: true,
                Computed: true,
            },
            "status_reason": schema.StringAttribute{
                MarkdownDescription: "Why the WAF is not enforcing the blocklist, as a sentence and its values that the Dashboard words in the reader's language: empty when it is. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "rule_details": schema.StringAttribute{
                MarkdownDescription: "What OneUptime keeps in the WAF: its rules' names and priorities, what they hold or reference, and the metric whose sampled requests are evidence. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "checked_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last read the WAF.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who picked the WAF. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *CloudWafDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudWafDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CloudWafDataSourceModel

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
    if !data.CloudConnectionId.IsNull() && !data.CloudConnectionId.IsUnknown() {
        filters["cloudConnectionId"] = data.CloudConnectionId.ValueString()
        filterNames = append(filterNames, "cloud_connection_id = "+fmt.Sprintf("%q", data.CloudConnectionId.ValueString()))
    }
    if !data.CloudResourceId.IsNull() && !data.CloudResourceId.IsUnknown() {
        filters["cloudResourceId"] = data.CloudResourceId.ValueString()
        filterNames = append(filterNames, "cloud_resource_id = "+fmt.Sprintf("%q", data.CloudResourceId.ValueString()))
    }
    if !data.CloudProvider.IsNull() && !data.CloudProvider.IsUnknown() {
        filters["cloudProvider"] = data.CloudProvider.ValueString()
        filterNames = append(filterNames, "cloud_provider = "+fmt.Sprintf("%q", data.CloudProvider.ValueString()))
    }
    if !data.WafId.IsNull() && !data.WafId.IsUnknown() {
        filters["wafId"] = data.WafId.ValueString()
        filterNames = append(filterNames, "waf_id = "+fmt.Sprintf("%q", data.WafId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Region.IsNull() && !data.Region.IsUnknown() {
        filters["region"] = data.Region.ValueString()
        filterNames = append(filterNames, "region = "+fmt.Sprintf("%q", data.Region.ValueString()))
    }
    if !data.ClientAddress.IsNull() && !data.ClientAddress.IsUnknown() {
        filters["clientAddress"] = data.ClientAddress.ValueString()
        filterNames = append(filterNames, "client_address = "+fmt.Sprintf("%q", data.ClientAddress.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the cloud waf up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the cloud waf up by.",
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
        "cloudConnectionId": true,
        "cloudResourceId": true,
        "cloudProvider": true,
        "wafId": true,
        "name": true,
        "region": true,
        "clientAddress": true,
        "status": true,
        "statusMessage": true,
        "statusReason": true,
        "ruleDetails": true,
        "checkedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-block-cloud-waf/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud_waf, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud waf found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read cloud_waf: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-block-cloud-waf/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud_waf, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list cloud_waf: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud waf matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one cloud waf matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for cloud_waf.")
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
    if obj, ok := item["cloudConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if arr, ok := item["cloudConnectionId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := item["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := item["cloudResourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if arr, ok := item["cloudResourceId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if val, ok := item["cloudResourceId"].(string); ok {
        data.CloudResourceId = types.StringValue(val)
    } else {
        data.CloudResourceId = types.StringNull()
    }
    if obj, ok := item["cloudProvider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudProvider = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudProvider = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudProvider = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudProvider = types.StringValue(string(jsonBytes))
        } else {
            data.CloudProvider = types.StringNull()
        }
    } else if arr, ok := item["cloudProvider"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CloudProvider = types.StringValue(string(jsonBytes))
        } else {
            data.CloudProvider = types.StringNull()
        }
    } else if val, ok := item["cloudProvider"].(string); ok {
        data.CloudProvider = types.StringValue(val)
    } else {
        data.CloudProvider = types.StringNull()
    }
    if obj, ok := item["wafId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WafId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WafId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WafId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WafId = types.StringValue(string(jsonBytes))
        } else {
            data.WafId = types.StringNull()
        }
    } else if arr, ok := item["wafId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.WafId = types.StringValue(string(jsonBytes))
        } else {
            data.WafId = types.StringNull()
        }
    } else if val, ok := item["wafId"].(string); ok {
        data.WafId = types.StringValue(val)
    } else {
        data.WafId = types.StringNull()
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
    if obj, ok := item["region"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if arr, ok := item["region"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := item["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
    }
    if obj, ok := item["clientAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClientAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClientAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClientAddress = types.StringValue(string(jsonBytes))
        } else {
            data.ClientAddress = types.StringNull()
        }
    } else if arr, ok := item["clientAddress"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ClientAddress = types.StringValue(string(jsonBytes))
        } else {
            data.ClientAddress = types.StringNull()
        }
    } else if val, ok := item["clientAddress"].(string); ok {
        data.ClientAddress = types.StringValue(val)
    } else {
        data.ClientAddress = types.StringNull()
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
    } else if arr, ok := item["status"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if arr, ok := item["statusMessage"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["statusReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusReason = types.StringValue(string(jsonBytes))
        } else {
            data.StatusReason = types.StringNull()
        }
    } else if arr, ok := item["statusReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StatusReason = types.StringValue(string(jsonBytes))
        } else {
            data.StatusReason = types.StringNull()
        }
    } else if val, ok := item["statusReason"].(string); ok {
        data.StatusReason = types.StringValue(val)
    } else {
        data.StatusReason = types.StringNull()
    }
    if obj, ok := item["ruleDetails"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleDetails = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleDetails = types.StringValue(string(jsonBytes))
        } else {
            data.RuleDetails = types.StringNull()
        }
    } else if arr, ok := item["ruleDetails"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RuleDetails = types.StringValue(string(jsonBytes))
        } else {
            data.RuleDetails = types.StringNull()
        }
    } else if val, ok := item["ruleDetails"].(string); ok {
        data.RuleDetails = types.StringValue(val)
    } else {
        data.RuleDetails = types.StringNull()
    }
    if obj, ok := item["checkedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CheckedAt = types.StringNull()
        }
    } else if arr, ok := item["checkedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CheckedAt = types.StringNull()
        }
    } else if val, ok := item["checkedAt"].(string); ok {
        data.CheckedAt = types.StringValue(val)
    } else {
        data.CheckedAt = types.StringNull()
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
