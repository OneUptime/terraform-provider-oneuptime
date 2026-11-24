package shard19

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
var _ datasource.DataSource = &WorkflowDataSource{}

func NewWorkflowDataSource() datasource.DataSource {
    return &WorkflowDataSource{}
}

// WorkflowDataSource defines the data source implementation.
type WorkflowDataSource struct {
    client *Client
}

// WorkflowDataSourceModel describes the data source data model.
type WorkflowDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Graph types.String `tfsdk:"graph"`
    GraphRevision types.Number `tfsdk:"graph_revision"`
    Labels types.Set `tfsdk:"labels"`
    RunAs types.String `tfsdk:"run_as"`
    RunAsCustomRoleId types.String `tfsdk:"run_as_custom_role_id"`
    RunAsGivenByUserId types.String `tfsdk:"run_as_given_by_user_id"`
    RunAsGivenByApiKeyId types.String `tfsdk:"run_as_given_by_api_key_id"`
    RunAsGivenAt types.String `tfsdk:"run_as_given_at"`
    WebhookSecretKey types.String `tfsdk:"webhook_secret_key"`
    IncomingEmailSecretKey types.String `tfsdk:"incoming_email_secret_key"`
    IsWebhookSigningSecretSet types.Bool `tfsdk:"is_webhook_signing_secret_set"`
    LastWebhookRequestAt types.String `tfsdk:"last_webhook_request_at"`
    LastWebhookError types.String `tfsdk:"last_webhook_error"`
    LastWebhookErrorAt types.String `tfsdk:"last_webhook_error_at"`
}

func (d *WorkflowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *WorkflowDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Integrate your OneUptime project with rest of your software stack. Look up an existing workflow by `id`, or by any of its other arguments (`name`, `archived_by_user_id`, `created_by_user_id`, ...): each one set must match, and exactly one workflow may match them all.",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Archived workflows are hidden from the Workflows list and never run, from any trigger. Unarchiving restores them as they were.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When this workflow was archived. Empty while it is not archived.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is this workflow enabled?",
                Optional: true,
                Computed: true,
            },
            "graph": schema.StringAttribute{
                MarkdownDescription: "Workflow Graph in JSON. Ideally, create this via UI and not via API. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "graph_revision": schema.NumberAttribute{
                MarkdownDescription: "How many times the workflow's graph has been saved. To save a graph only if nobody has saved one since you read the workflow, send this number as miscDataProps.graphRevision with the graph: a save made to an older graph is refused with 409 Conflict. Set by the server.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "run_as": schema.StringAttribute{
                MarkdownDescription: "What the workflow's steps can do in the project. Default: what a Project Admin can, except change who has access. ProjectOwner: what a Project Owner can. CustomRole: exactly what the custom role in runAsCustomRoleId allows. Only a person or a project API key that has those permissions can choose anything but Default, and only such a person or key can then change the workflow or run it by hand.",
                Optional: true,
                Computed: true,
            },
            "run_as_custom_role_id": schema.StringAttribute{
                MarkdownDescription: "ID of the custom role the workflow's steps act as, when Run As is CustomRole. Empty otherwise. The ID of a `oneuptime_custom_role`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_custom_role"),
            },
            "run_as_given_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person whose permissions the workflow runs with when it runs as more than Default. Set by OneUptime. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "run_as_given_by_api_key_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project API key whose permissions the workflow runs with when it runs as more than Default and a key gave them. Set by OneUptime. The ID of a `oneuptime_api_key`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_api_key"),
            },
            "run_as_given_at": schema.StringAttribute{
                MarkdownDescription: "When the person in runAsGivenByUserId last gave the workflow its permissions. Set by OneUptime.",
                Computed: true,
            },
            "webhook_secret_key": schema.StringAttribute{
                MarkdownDescription: "The secret part of the Webhook trigger's URL (/workflow/trigger/<key>). Anyone who has the URL can start the workflow, so only people who can edit the workflow can read the key. Generated when the workflow is created; set a new value to reset the URL.",
                Computed: true,
                Sensitive: true,
            },
            "incoming_email_secret_key": schema.StringAttribute{
                MarkdownDescription: "The secret part of the Incoming Email trigger's address (workflow-<key>@<inbound email domain>). Anyone who has the address can start the workflow, so only people who can edit the workflow can read the key. Given to the workflow when its graph first has an Incoming Email trigger; set a new UUID to reset the address. Unique across all workflows.",
                Computed: true,
                Sensitive: true,
            },
            "is_webhook_signing_secret_set": schema.BoolAttribute{
                MarkdownDescription: "Whether a webhook signing secret is saved. A trigger that checks signatures refuses every request until it is. Set by the server from the signing secret itself.",
                Optional: true,
                Computed: true,
            },
            "last_webhook_request_at": schema.StringAttribute{
                MarkdownDescription: "When a request signed with the webhook signing secret last arrived. Empty until the sender's first one.",
                Computed: true,
            },
            "last_webhook_error": schema.StringAttribute{
                MarkdownDescription: "Why the last request to a trigger that checks signatures was refused, if it was. Cleared by the next request that is accepted, and by saving a new signing secret.",
                Optional: true,
                Computed: true,
            },
            "last_webhook_error_at": schema.StringAttribute{
                MarkdownDescription: "When the last webhook error happened.",
                Computed: true,
            },
        },
    }
}

func (d *WorkflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkflowDataSourceModel

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
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.GraphRevision.IsNull() && !data.GraphRevision.IsUnknown() {
        filters["graphRevision"] = LookupNumber(data.GraphRevision)
        filterNames = append(filterNames, "graph_revision = "+data.GraphRevision.ValueBigFloat().String())
    }
    if !data.RunAs.IsNull() && !data.RunAs.IsUnknown() {
        filters["runAs"] = data.RunAs.ValueString()
        filterNames = append(filterNames, "run_as = "+fmt.Sprintf("%q", data.RunAs.ValueString()))
    }
    if !data.RunAsCustomRoleId.IsNull() && !data.RunAsCustomRoleId.IsUnknown() {
        filters["runAsCustomRoleId"] = data.RunAsCustomRoleId.ValueString()
        filterNames = append(filterNames, "run_as_custom_role_id = "+fmt.Sprintf("%q", data.RunAsCustomRoleId.ValueString()))
    }
    if !data.RunAsGivenByUserId.IsNull() && !data.RunAsGivenByUserId.IsUnknown() {
        filters["runAsGivenByUserId"] = data.RunAsGivenByUserId.ValueString()
        filterNames = append(filterNames, "run_as_given_by_user_id = "+fmt.Sprintf("%q", data.RunAsGivenByUserId.ValueString()))
    }
    if !data.RunAsGivenByApiKeyId.IsNull() && !data.RunAsGivenByApiKeyId.IsUnknown() {
        filters["runAsGivenByApiKeyId"] = data.RunAsGivenByApiKeyId.ValueString()
        filterNames = append(filterNames, "run_as_given_by_api_key_id = "+fmt.Sprintf("%q", data.RunAsGivenByApiKeyId.ValueString()))
    }
    if !data.IsWebhookSigningSecretSet.IsNull() && !data.IsWebhookSigningSecretSet.IsUnknown() {
        filters["isWebhookSigningSecretSet"] = data.IsWebhookSigningSecretSet.ValueBool()
        filterNames = append(filterNames, "is_webhook_signing_secret_set = "+fmt.Sprintf("%t", data.IsWebhookSigningSecretSet.ValueBool()))
    }
    if !data.LastWebhookError.IsNull() && !data.LastWebhookError.IsUnknown() {
        filters["lastWebhookError"] = data.LastWebhookError.ValueString()
        filterNames = append(filterNames, "last_webhook_error = "+fmt.Sprintf("%q", data.LastWebhookError.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the workflow up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the workflow up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "isEnabled": true,
        "graph": true,
        "graphRevision": true,
        "labels": true,
        "runAs": true,
        "runAsCustomRoleId": true,
        "runAsGivenByUserId": true,
        "runAsGivenByApiKeyId": true,
        "runAsGivenAt": true,
        "webhookSecretKey": true,
        "incomingEmailSecretKey": true,
        "isWebhookSigningSecretSet": true,
        "lastWebhookRequestAt": true,
        "lastWebhookError": true,
        "lastWebhookErrorAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workflow/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workflow, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workflow: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workflow/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workflow, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workflow: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workflow matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workflow.")
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
    } else if arr, ok := item["slug"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
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
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if arr, ok := item["archivedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := item["archivedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["graph"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Graph = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Graph = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Graph = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Graph = types.StringValue(string(jsonBytes))
        } else {
            data.Graph = types.StringNull()
        }
    } else if arr, ok := item["graph"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Graph = types.StringValue(string(jsonBytes))
        } else {
            data.Graph = types.StringNull()
        }
    } else if val, ok := item["graph"].(string); ok {
        data.Graph = types.StringValue(val)
    } else {
        data.Graph = types.StringNull()
    }
    data.GraphRevision = NumberFromAPI(data.GraphRevision, item["graphRevision"])
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
    if obj, ok := item["runAs"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAs = types.StringValue(string(jsonBytes))
        } else {
            data.RunAs = types.StringNull()
        }
    } else if arr, ok := item["runAs"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAs = types.StringValue(string(jsonBytes))
        } else {
            data.RunAs = types.StringNull()
        }
    } else if val, ok := item["runAs"].(string); ok {
        data.RunAs = types.StringValue(val)
    } else {
        data.RunAs = types.StringNull()
    }
    if obj, ok := item["runAsCustomRoleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsCustomRoleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsCustomRoleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsCustomRoleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsCustomRoleId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsCustomRoleId = types.StringNull()
        }
    } else if arr, ok := item["runAsCustomRoleId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsCustomRoleId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsCustomRoleId = types.StringNull()
        }
    } else if val, ok := item["runAsCustomRoleId"].(string); ok {
        data.RunAsCustomRoleId = types.StringValue(val)
    } else {
        data.RunAsCustomRoleId = types.StringNull()
    }
    if obj, ok := item["runAsGivenByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsGivenByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsGivenByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsGivenByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsGivenByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByUserId = types.StringNull()
        }
    } else if arr, ok := item["runAsGivenByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsGivenByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByUserId = types.StringNull()
        }
    } else if val, ok := item["runAsGivenByUserId"].(string); ok {
        data.RunAsGivenByUserId = types.StringValue(val)
    } else {
        data.RunAsGivenByUserId = types.StringNull()
    }
    if obj, ok := item["runAsGivenByApiKeyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsGivenByApiKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsGivenByApiKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsGivenByApiKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsGivenByApiKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByApiKeyId = types.StringNull()
        }
    } else if arr, ok := item["runAsGivenByApiKeyId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsGivenByApiKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByApiKeyId = types.StringNull()
        }
    } else if val, ok := item["runAsGivenByApiKeyId"].(string); ok {
        data.RunAsGivenByApiKeyId = types.StringValue(val)
    } else {
        data.RunAsGivenByApiKeyId = types.StringNull()
    }
    if obj, ok := item["runAsGivenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsGivenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsGivenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsGivenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsGivenAt = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenAt = types.StringNull()
        }
    } else if arr, ok := item["runAsGivenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsGivenAt = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenAt = types.StringNull()
        }
    } else if val, ok := item["runAsGivenAt"].(string); ok {
        data.RunAsGivenAt = types.StringValue(val)
    } else {
        data.RunAsGivenAt = types.StringNull()
    }
    if obj, ok := item["webhookSecretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WebhookSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WebhookSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WebhookSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WebhookSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.WebhookSecretKey = types.StringNull()
        }
    } else if arr, ok := item["webhookSecretKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.WebhookSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.WebhookSecretKey = types.StringNull()
        }
    } else if val, ok := item["webhookSecretKey"].(string); ok {
        data.WebhookSecretKey = types.StringValue(val)
    } else {
        data.WebhookSecretKey = types.StringNull()
    }
    if obj, ok := item["incomingEmailSecretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailSecretKey = types.StringNull()
        }
    } else if arr, ok := item["incomingEmailSecretKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailSecretKey = types.StringNull()
        }
    } else if val, ok := item["incomingEmailSecretKey"].(string); ok {
        data.IncomingEmailSecretKey = types.StringValue(val)
    } else {
        data.IncomingEmailSecretKey = types.StringNull()
    }
    if val, ok := item["isWebhookSigningSecretSet"].(bool); ok {
        data.IsWebhookSigningSecretSet = types.BoolValue(val)
    } else {
        data.IsWebhookSigningSecretSet = types.BoolNull()
    }
    if obj, ok := item["lastWebhookRequestAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWebhookRequestAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastWebhookRequestAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastWebhookRequestAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastWebhookRequestAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastWebhookRequestAt = types.StringNull()
        }
    } else if arr, ok := item["lastWebhookRequestAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastWebhookRequestAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastWebhookRequestAt = types.StringNull()
        }
    } else if val, ok := item["lastWebhookRequestAt"].(string); ok {
        data.LastWebhookRequestAt = types.StringValue(val)
    } else {
        data.LastWebhookRequestAt = types.StringNull()
    }
    if obj, ok := item["lastWebhookError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWebhookError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastWebhookError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastWebhookError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastWebhookError = types.StringValue(string(jsonBytes))
        } else {
            data.LastWebhookError = types.StringNull()
        }
    } else if arr, ok := item["lastWebhookError"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastWebhookError = types.StringValue(string(jsonBytes))
        } else {
            data.LastWebhookError = types.StringNull()
        }
    } else if val, ok := item["lastWebhookError"].(string); ok {
        data.LastWebhookError = types.StringValue(val)
    } else {
        data.LastWebhookError = types.StringNull()
    }
    if obj, ok := item["lastWebhookErrorAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastWebhookErrorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastWebhookErrorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastWebhookErrorAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastWebhookErrorAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastWebhookErrorAt = types.StringNull()
        }
    } else if arr, ok := item["lastWebhookErrorAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastWebhookErrorAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastWebhookErrorAt = types.StringNull()
        }
    } else if val, ok := item["lastWebhookErrorAt"].(string); ok {
        data.LastWebhookErrorAt = types.StringValue(val)
    } else {
        data.LastWebhookErrorAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
