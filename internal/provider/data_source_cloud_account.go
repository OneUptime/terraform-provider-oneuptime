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
var _ datasource.DataSource = &CloudAccountDataSource{}

func NewCloudAccountDataSource() datasource.DataSource {
    return &CloudAccountDataSource{}
}

// CloudAccountDataSource defines the data source implementation.
type CloudAccountDataSource struct {
    client *Client
}

// CloudAccountDataSourceModel describes the data source data model.
type CloudAccountDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    ProviderValue types.String `tfsdk:"provider_value"`
    AccountId types.String `tfsdk:"account_id"`
    ProviderAccountName types.String `tfsdk:"provider_account_name"`
    CloudConnectionId types.String `tfsdk:"cloud_connection_id"`
    ResourceCount types.Number `tfsdk:"resource_count"`
    RegionCount types.Number `tfsdk:"region_count"`
    ResourceSummary types.String `tfsdk:"resource_summary"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    InventoryListedAt types.String `tfsdk:"inventory_listed_at"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsAiInvestigationEnabled types.Bool `tfsdk:"is_ai_investigation_enabled"`
    IsAiLogReadEnabled types.Bool `tfsdk:"is_ai_log_read_enabled"`
    AiRemediationMode types.String `tfsdk:"ai_remediation_mode"`
    AiCommandAllowlist types.String `tfsdk:"ai_command_allowlist"`
    AiAccessLastVerifiedAt types.String `tfsdk:"ai_access_last_verified_at"`
    AiAccessLastError types.String `tfsdk:"ai_access_last_error"`
    AiAccessConfiguredAt types.String `tfsdk:"ai_access_configured_at"`
}

func (d *CloudAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cloud_account"
}

func (d *CloudAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them. Look up an existing cloud account by `id`, or by any of its other arguments (`name`, `account_id`, `ai_access_last_error`, ...): each one set must match, and exactly one cloud account may match them all.",

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
                MarkdownDescription: "ID of the project this account belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the account is called in OneUptime. Starts as the name its cloud gives it, and can be changed.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the account is for, in your team's words.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Which cloud the account is in: 'aws', 'azure' or 'gcp'.",
                Optional: true,
                Computed: true,
            },
            "account_id": schema.StringAttribute{
                MarkdownDescription: "The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.",
                Optional: true,
                Computed: true,
            },
            "provider_account_name": schema.StringAttribute{
                MarkdownDescription: "What the cloud calls the account: an AWS account alias, an Azure subscription's or a Google Cloud project's name.",
                Optional: true,
                Computed: true,
            },
            "cloud_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the connection that reads this account, if one does. The ID of a `oneuptime_cloud_connection`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_cloud_connection"),
            },
            "resource_count": schema.NumberAttribute{
                MarkdownDescription: "How many live resources the account has in OneUptime.",
                Optional: true,
                Computed: true,
            },
            "region_count": schema.NumberAttribute{
                MarkdownDescription: "How many regions the account's live resources are in.",
                Optional: true,
                Computed: true,
            },
            "resource_summary": schema.StringAttribute{
                MarkdownDescription: "The account's live resources counted by type and by region (CloudAccountSummary). A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a cloud connection last reached the account.",
                Computed: true,
            },
            "inventory_listed_at": schema.StringAttribute{
                MarkdownDescription: "When the account's inventory last listed every resource it has.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this account. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_ai_investigation_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, OneUptime AI reads this account live while investigating incidents and alerts linked to its resources — a resource's settings and status, its recent events, load balancer target health and quota usage — through the account's read-only cloud connection, and uses what it read as evidence; it needs no Cloud AI agent for that. Nothing is ever changed by an investigation. On by default. Anyone who may edit the account can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "is_ai_log_read_enabled": schema.BoolAttribute{
                MarkdownDescription: "When on, an investigation may also read this account's recent log lines (CloudWatch Logs, Log Analytics, Cloud Logging) through the read-only connection. Log lines can hold customer data, so this is off until someone turns it on. Anyone who may edit the account can turn it on or off.",
                Optional: true,
                Computed: true,
            },
            "ai_remediation_mode": schema.StringAttribute{
                MarkdownDescription: "Disabled: AI never proposes or runs a change on this account. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the account's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). The Cloud AI agent's own configuration (ONEUPTIME_AI_FIXES, which its setup command writes) sets the mode once it connects. Anyone who may edit the account can lower the mode; raising it needs Project Owner, Project Admin or Edit AI Rule.",
                Optional: true,
                Computed: true,
            },
            "ai_command_allowlist": schema.StringAttribute{
                MarkdownDescription: "Optional JSON array of command patterns that Automatic mode may run in this account without approval even though they are riskier changes. Each pattern is one cloud catalog command (cloud <operation> RESOURCE [--flag value]) and is compared with the command word by word: * stands for exactly one word (a resource id, a count), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid change for this account's cloud is refused. Destructive operations (Denied tier) never run regardless, and an operation that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit AI Rule; anyone who may edit the account can remove patterns or clear the list. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "ai_access_last_verified_at": schema.StringAttribute{
                MarkdownDescription: "When an operation from OneUptime AI last succeeded in this account through its Cloud AI agent. Set by the server.",
                Computed: true,
            },
            "ai_access_last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent failure OneUptime AI hit in this account — an operation that failed, or a Cloud AI agent that could not connect to it — kept until the next success. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "ai_access_configured_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI access to this account was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Cloud AI agent that connects later never overwrites a setting an operator chose.",
                Computed: true,
            },
        },
    }
}

func (d *CloudAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CloudAccountDataSourceModel

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
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.AccountId.IsNull() && !data.AccountId.IsUnknown() {
        filters["accountId"] = data.AccountId.ValueString()
        filterNames = append(filterNames, "account_id = "+fmt.Sprintf("%q", data.AccountId.ValueString()))
    }
    if !data.ProviderAccountName.IsNull() && !data.ProviderAccountName.IsUnknown() {
        filters["providerAccountName"] = data.ProviderAccountName.ValueString()
        filterNames = append(filterNames, "provider_account_name = "+fmt.Sprintf("%q", data.ProviderAccountName.ValueString()))
    }
    if !data.CloudConnectionId.IsNull() && !data.CloudConnectionId.IsUnknown() {
        filters["cloudConnectionId"] = data.CloudConnectionId.ValueString()
        filterNames = append(filterNames, "cloud_connection_id = "+fmt.Sprintf("%q", data.CloudConnectionId.ValueString()))
    }
    if !data.ResourceCount.IsNull() && !data.ResourceCount.IsUnknown() {
        filters["resourceCount"] = lookupNumber(data.ResourceCount)
        filterNames = append(filterNames, "resource_count = "+data.ResourceCount.ValueBigFloat().String())
    }
    if !data.RegionCount.IsNull() && !data.RegionCount.IsUnknown() {
        filters["regionCount"] = lookupNumber(data.RegionCount)
        filterNames = append(filterNames, "region_count = "+data.RegionCount.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsAiInvestigationEnabled.IsNull() && !data.IsAiInvestigationEnabled.IsUnknown() {
        filters["isAiInvestigationEnabled"] = data.IsAiInvestigationEnabled.ValueBool()
        filterNames = append(filterNames, "is_ai_investigation_enabled = "+fmt.Sprintf("%t", data.IsAiInvestigationEnabled.ValueBool()))
    }
    if !data.IsAiLogReadEnabled.IsNull() && !data.IsAiLogReadEnabled.IsUnknown() {
        filters["isAiLogReadEnabled"] = data.IsAiLogReadEnabled.ValueBool()
        filterNames = append(filterNames, "is_ai_log_read_enabled = "+fmt.Sprintf("%t", data.IsAiLogReadEnabled.ValueBool()))
    }
    if !data.AiRemediationMode.IsNull() && !data.AiRemediationMode.IsUnknown() {
        filters["aiRemediationMode"] = data.AiRemediationMode.ValueString()
        filterNames = append(filterNames, "ai_remediation_mode = "+fmt.Sprintf("%q", data.AiRemediationMode.ValueString()))
    }
    if !data.AiAccessLastError.IsNull() && !data.AiAccessLastError.IsUnknown() {
        filters["aiAccessLastError"] = data.AiAccessLastError.ValueString()
        filterNames = append(filterNames, "ai_access_last_error = "+fmt.Sprintf("%q", data.AiAccessLastError.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the cloud account up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the cloud account up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "provider": true,
        "accountId": true,
        "providerAccountName": true,
        "cloudConnectionId": true,
        "resourceCount": true,
        "regionCount": true,
        "resourceSummary": true,
        "lastSeenAt": true,
        "inventoryListedAt": true,
        "labels": true,
        "createdByUserId": true,
        "isAiInvestigationEnabled": true,
        "isAiLogReadEnabled": true,
        "aiRemediationMode": true,
        "aiCommandAllowlist": true,
        "aiAccessLastVerifiedAt": true,
        "aiAccessLastError": true,
        "aiAccessConfiguredAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/cloud-account/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud_account, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud account found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read cloud_account: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/cloud-account/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud_account, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list cloud_account: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud account matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one cloud account matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for cloud_account.")
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
    if obj, ok := item["provider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := item["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
    }
    if obj, ok := item["accountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AccountId = types.StringValue(string(jsonBytes))
        } else {
            data.AccountId = types.StringNull()
        }
    } else if val, ok := item["accountId"].(string); ok {
        data.AccountId = types.StringValue(val)
    } else {
        data.AccountId = types.StringNull()
    }
    if obj, ok := item["providerAccountName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderAccountName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderAccountName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderAccountName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderAccountName = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderAccountName = types.StringNull()
        }
    } else if val, ok := item["providerAccountName"].(string); ok {
        data.ProviderAccountName = types.StringValue(val)
    } else {
        data.ProviderAccountName = types.StringNull()
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
    } else if val, ok := item["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    data.ResourceCount = numberFromAPI(data.ResourceCount, item["resourceCount"])
    data.RegionCount = numberFromAPI(data.RegionCount, item["regionCount"])
    if obj, ok := item["resourceSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceSummary = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceSummary = types.StringNull()
        }
    } else if val, ok := item["resourceSummary"].(string); ok {
        data.ResourceSummary = types.StringValue(val)
    } else {
        data.ResourceSummary = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    if obj, ok := item["inventoryListedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InventoryListedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InventoryListedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InventoryListedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InventoryListedAt = types.StringValue(string(jsonBytes))
        } else {
            data.InventoryListedAt = types.StringNull()
        }
    } else if val, ok := item["inventoryListedAt"].(string); ok {
        data.InventoryListedAt = types.StringValue(val)
    } else {
        data.InventoryListedAt = types.StringNull()
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
    if val, ok := item["isAiInvestigationEnabled"].(bool); ok {
        data.IsAiInvestigationEnabled = types.BoolValue(val)
    } else {
        data.IsAiInvestigationEnabled = types.BoolNull()
    }
    if val, ok := item["isAiLogReadEnabled"].(bool); ok {
        data.IsAiLogReadEnabled = types.BoolValue(val)
    } else {
        data.IsAiLogReadEnabled = types.BoolNull()
    }
    if obj, ok := item["aiRemediationMode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiRemediationMode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiRemediationMode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiRemediationMode = types.StringValue(string(jsonBytes))
        } else {
            data.AiRemediationMode = types.StringNull()
        }
    } else if val, ok := item["aiRemediationMode"].(string); ok {
        data.AiRemediationMode = types.StringValue(val)
    } else {
        data.AiRemediationMode = types.StringNull()
    }
    if obj, ok := item["aiCommandAllowlist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiCommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiCommandAllowlist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiCommandAllowlist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiCommandAllowlist = types.StringValue(string(jsonBytes))
        } else {
            data.AiCommandAllowlist = types.StringNull()
        }
    } else if val, ok := item["aiCommandAllowlist"].(string); ok {
        data.AiCommandAllowlist = types.StringValue(val)
    } else {
        data.AiCommandAllowlist = types.StringNull()
    }
    if obj, ok := item["aiAccessLastVerifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastVerifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessLastVerifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessLastVerifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessLastVerifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastVerifiedAt = types.StringNull()
        }
    } else if val, ok := item["aiAccessLastVerifiedAt"].(string); ok {
        data.AiAccessLastVerifiedAt = types.StringValue(val)
    } else {
        data.AiAccessLastVerifiedAt = types.StringNull()
    }
    if obj, ok := item["aiAccessLastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessLastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessLastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessLastError = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessLastError = types.StringNull()
        }
    } else if val, ok := item["aiAccessLastError"].(string); ok {
        data.AiAccessLastError = types.StringValue(val)
    } else {
        data.AiAccessLastError = types.StringNull()
    }
    if obj, ok := item["aiAccessConfiguredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiAccessConfiguredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiAccessConfiguredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiAccessConfiguredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiAccessConfiguredAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiAccessConfiguredAt = types.StringNull()
        }
    } else if val, ok := item["aiAccessConfiguredAt"].(string); ok {
        data.AiAccessConfiguredAt = types.StringValue(val)
    } else {
        data.AiAccessConfiguredAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
