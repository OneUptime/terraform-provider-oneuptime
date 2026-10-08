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
var _ datasource.DataSource = &RunnerDataSource{}

func NewRunnerDataSource() datasource.DataSource {
    return &RunnerDataSource{}
}

// RunnerDataSource defines the data source implementation.
type RunnerDataSource struct {
    client *Client
}

// RunnerDataSourceModel describes the data source data model.
type RunnerDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    Key types.String `tfsdk:"key"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastAlive types.String `tfsdk:"last_alive"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
    CanRunRunbooks types.Bool `tfsdk:"can_run_runbooks"`
    CanRunCodeFixTasks types.Bool `tfsdk:"can_run_code_fix_tasks"`
    CanRunAiCommands types.Bool `tfsdk:"can_run_ai_commands"`
    HostInfo types.String `tfsdk:"host_info"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *RunnerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runner"
}

func (d *RunnerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them. Look up an existing runner by `id`, or by any of its other arguments (`name`, `can_run_ai_commands`, `can_run_code_fix_tasks`, ...): each one set must match, and exactly one runner may match them all.",

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
                MarkdownDescription: "Friendly name for this agent. Names starting with \"kubernetes-agent/\" are reserved for the in-cluster Runners the Kubernetes agent chart registers, and such a Runner cannot be renamed.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Optional description for this agent.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Secret key the agent presents on every request. Anyone who can read this key can claim work as this Runner and receive its secrets in plaintext. Never share it; reset it to revoke the agent.",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Self-reported version of the Runner binary. Updated on each heartbeat.",
                Computed: true,
            },
            "last_alive": schema.StringAttribute{
                MarkdownDescription: "Most recent heartbeat from this agent.",
                Computed: true,
            },
            "connection_status": schema.StringAttribute{
                MarkdownDescription: "Connected if the agent has heartbeated recently.",
                Optional: true,
                Computed: true,
            },
            "can_run_runbooks": schema.BoolAttribute{
                MarkdownDescription: "Whether this Runner executes runbook steps. On by default — this is why most Runners are installed. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only.",
                Optional: true,
                Computed: true,
            },
            "can_run_code_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "Whether this Runner works in your code repository to open AI fix pull requests. Off by default; it requires a connected code repository. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only.",
                Optional: true,
                Computed: true,
            },
            "can_run_ai_commands": schema.BoolAttribute{
                MarkdownDescription: "Whether OneUptime AI may run commands through this Runner: read-only kubectl while investigating a cluster it is bound to, and policy-checked remediation commands (Bash, SSH, kubectl) that either match an allowlist or wait for one-click human approval. Off by default; the in-cluster Runner installed by the Kubernetes agent chart turns it on.",
                Optional: true,
                Computed: true,
            },
            "host_info": schema.StringAttribute{
                MarkdownDescription: "Self-reported host info (hostname, OS, arch). Updated on each heartbeat. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *RunnerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunnerDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Key.IsNull() && !data.Key.IsUnknown() {
        filters["key"] = data.Key.ValueString()
        filterNames = append(filterNames, "key = "+fmt.Sprintf("%q", data.Key.ValueString()))
    }
    if !data.ConnectionStatus.IsNull() && !data.ConnectionStatus.IsUnknown() {
        filters["connectionStatus"] = data.ConnectionStatus.ValueString()
        filterNames = append(filterNames, "connection_status = "+fmt.Sprintf("%q", data.ConnectionStatus.ValueString()))
    }
    if !data.CanRunRunbooks.IsNull() && !data.CanRunRunbooks.IsUnknown() {
        filters["canRunRunbooks"] = data.CanRunRunbooks.ValueBool()
        filterNames = append(filterNames, "can_run_runbooks = "+fmt.Sprintf("%t", data.CanRunRunbooks.ValueBool()))
    }
    if !data.CanRunCodeFixTasks.IsNull() && !data.CanRunCodeFixTasks.IsUnknown() {
        filters["canRunCodeFixTasks"] = data.CanRunCodeFixTasks.ValueBool()
        filterNames = append(filterNames, "can_run_code_fix_tasks = "+fmt.Sprintf("%t", data.CanRunCodeFixTasks.ValueBool()))
    }
    if !data.CanRunAiCommands.IsNull() && !data.CanRunAiCommands.IsUnknown() {
        filters["canRunAiCommands"] = data.CanRunAiCommands.ValueBool()
        filterNames = append(filterNames, "can_run_ai_commands = "+fmt.Sprintf("%t", data.CanRunAiCommands.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the runner up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the runner up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "slug": true,
        "key": true,
        "agentVersion": true,
        "lastAlive": true,
        "connectionStatus": true,
        "canRunRunbooks": true,
        "canRunCodeFixTasks": true,
        "canRunAiCommands": true,
        "hostInfo": true,
        "createdByUserId": true,
        "labels": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/runner/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runner, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runner found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read runner: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/runner/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list runner, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list runner: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runner matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one runner matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for runner.")
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
    if obj, ok := item["key"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Key = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if val, ok := item["key"].(string); ok {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
    }
    if obj, ok := item["agentVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := item["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
    }
    if obj, ok := item["lastAlive"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAlive = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAlive = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAlive = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAlive = types.StringValue(string(jsonBytes))
        } else {
            data.LastAlive = types.StringNull()
        }
    } else if val, ok := item["lastAlive"].(string); ok {
        data.LastAlive = types.StringValue(val)
    } else {
        data.LastAlive = types.StringNull()
    }
    if obj, ok := item["connectionStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConnectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConnectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConnectionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConnectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ConnectionStatus = types.StringNull()
        }
    } else if val, ok := item["connectionStatus"].(string); ok {
        data.ConnectionStatus = types.StringValue(val)
    } else {
        data.ConnectionStatus = types.StringNull()
    }
    if val, ok := item["canRunRunbooks"].(bool); ok {
        data.CanRunRunbooks = types.BoolValue(val)
    } else {
        data.CanRunRunbooks = types.BoolNull()
    }
    if val, ok := item["canRunCodeFixTasks"].(bool); ok {
        data.CanRunCodeFixTasks = types.BoolValue(val)
    } else {
        data.CanRunCodeFixTasks = types.BoolNull()
    }
    if val, ok := item["canRunAiCommands"].(bool); ok {
        data.CanRunAiCommands = types.BoolValue(val)
    } else {
        data.CanRunAiCommands = types.BoolNull()
    }
    if obj, ok := item["hostInfo"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostInfo = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostInfo = types.StringValue(string(jsonBytes))
        } else {
            data.HostInfo = types.StringNull()
        }
    } else if val, ok := item["hostInfo"].(string); ok {
        data.HostInfo = types.StringValue(val)
    } else {
        data.HostInfo = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
