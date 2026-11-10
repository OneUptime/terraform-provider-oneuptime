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
var _ datasource.DataSource = &WorkflowVersionDataSource{}

func NewWorkflowVersionDataSource() datasource.DataSource {
    return &WorkflowVersionDataSource{}
}

// WorkflowVersionDataSource defines the data source implementation.
type WorkflowVersionDataSource struct {
    client *Client
}

// WorkflowVersionDataSourceModel describes the data source data model.
type WorkflowVersionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    WorkflowId types.String `tfsdk:"workflow_id"`
    VersionNumber types.Number `tfsdk:"version_number"`
    VersionType types.String `tfsdk:"version_type"`
    Graph types.String `tfsdk:"graph"`
    GraphFingerprint types.String `tfsdk:"graph_fingerprint"`
    Changes types.String `tfsdk:"changes"`
    ChangeCount types.Number `tfsdk:"change_count"`
    RestoredFromVersionNumber types.Number `tfsdk:"restored_from_version_number"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SavedThrough types.String `tfsdk:"saved_through"`
}

func (d *WorkflowVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflow_version"
}

func (d *WorkflowVersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Every saved version of a workflow's graph: who saved it, when, and what it changed, so an earlier one can be looked at and restored. Look up an existing workflow version by `id`, or by any of its other arguments (`change_count`, `created_by_user_id`, `graph_fingerprint`, ...): each one set must match, and exactly one workflow version may match them all.",

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
            "workflow_id": schema.StringAttribute{
                MarkdownDescription: "ID of the workflow this is a version of. The ID of a `oneuptime_workflow`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_workflow"),
            },
            "version_number": schema.NumberAttribute{
                MarkdownDescription: "The version's number: one more than the version before it. Only the newest 100 versions of a workflow are kept.",
                Optional: true,
                Computed: true,
            },
            "version_type": schema.StringAttribute{
                MarkdownDescription: "What made this version: Created (the workflow as it was created), Changed (changes saved in the builder, by Ask AI or through the API), Restored (an earlier version put back) or Earlier (the workflow as it was when version history began).",
                Optional: true,
                Computed: true,
            },
            "graph": schema.StringAttribute{
                MarkdownDescription: "The workflow's graph as this version saved it: its steps, where each one sits, and their connections. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "graph_fingerprint": schema.StringAttribute{
                MarkdownDescription: "SHA-256 of the graph without what only the builder's canvas keeps (which card is selected, how big each one is). Two versions with the same graph have the same fingerprint.",
                Optional: true,
                Computed: true,
            },
            "changes": schema.StringAttribute{
                MarkdownDescription: "What this version changed from the version before it, one plain line per change: steps added, removed or reconfigured, connections made or removed, steps moved. A setting is named, never quoted. At most 25 lines. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "change_count": schema.NumberAttribute{
                MarkdownDescription: "How many changes this version made, counted before its lines were cut to 25.",
                Optional: true,
                Computed: true,
            },
            "restored_from_version_number": schema.NumberAttribute{
                MarkdownDescription: "For a Restored version: the number of the earlier version it put back.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who saved this version. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "saved_through": schema.StringAttribute{
                MarkdownDescription: "AskAi when Ask AI saved this version for the person in Saved by User: each change Ask AI makes is a version of its own. Empty when the person saved it themselves, or the API did.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *WorkflowVersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkflowVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkflowVersionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.WorkflowId.IsNull() && !data.WorkflowId.IsUnknown() {
        filters["workflowId"] = data.WorkflowId.ValueString()
        filterNames = append(filterNames, "workflow_id = "+fmt.Sprintf("%q", data.WorkflowId.ValueString()))
    }
    if !data.VersionNumber.IsNull() && !data.VersionNumber.IsUnknown() {
        filters["versionNumber"] = lookupNumber(data.VersionNumber)
        filterNames = append(filterNames, "version_number = "+data.VersionNumber.ValueBigFloat().String())
    }
    if !data.VersionType.IsNull() && !data.VersionType.IsUnknown() {
        filters["versionType"] = data.VersionType.ValueString()
        filterNames = append(filterNames, "version_type = "+fmt.Sprintf("%q", data.VersionType.ValueString()))
    }
    if !data.GraphFingerprint.IsNull() && !data.GraphFingerprint.IsUnknown() {
        filters["graphFingerprint"] = data.GraphFingerprint.ValueString()
        filterNames = append(filterNames, "graph_fingerprint = "+fmt.Sprintf("%q", data.GraphFingerprint.ValueString()))
    }
    if !data.ChangeCount.IsNull() && !data.ChangeCount.IsUnknown() {
        filters["changeCount"] = lookupNumber(data.ChangeCount)
        filterNames = append(filterNames, "change_count = "+data.ChangeCount.ValueBigFloat().String())
    }
    if !data.RestoredFromVersionNumber.IsNull() && !data.RestoredFromVersionNumber.IsUnknown() {
        filters["restoredFromVersionNumber"] = lookupNumber(data.RestoredFromVersionNumber)
        filterNames = append(filterNames, "restored_from_version_number = "+data.RestoredFromVersionNumber.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.SavedThrough.IsNull() && !data.SavedThrough.IsUnknown() {
        filters["savedThrough"] = data.SavedThrough.ValueString()
        filterNames = append(filterNames, "saved_through = "+fmt.Sprintf("%q", data.SavedThrough.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the workflow version up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the workflow version up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "workflowId": true,
        "versionNumber": true,
        "versionType": true,
        "graph": true,
        "graphFingerprint": true,
        "changes": true,
        "changeCount": true,
        "restoredFromVersionNumber": true,
        "createdByUserId": true,
        "savedThrough": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workflow-version/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workflow_version, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow version found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workflow_version: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workflow-version/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workflow_version, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workflow_version: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow version matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workflow version matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workflow_version.")
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
    if obj, ok := item["workflowId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkflowId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkflowId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkflowId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkflowId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowId = types.StringNull()
        }
    } else if arr, ok := item["workflowId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.WorkflowId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowId = types.StringNull()
        }
    } else if val, ok := item["workflowId"].(string); ok {
        data.WorkflowId = types.StringValue(val)
    } else {
        data.WorkflowId = types.StringNull()
    }
    data.VersionNumber = numberFromAPI(data.VersionNumber, item["versionNumber"])
    if obj, ok := item["versionType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VersionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VersionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VersionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VersionType = types.StringValue(string(jsonBytes))
        } else {
            data.VersionType = types.StringNull()
        }
    } else if arr, ok := item["versionType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.VersionType = types.StringValue(string(jsonBytes))
        } else {
            data.VersionType = types.StringNull()
        }
    } else if val, ok := item["versionType"].(string); ok {
        data.VersionType = types.StringValue(val)
    } else {
        data.VersionType = types.StringNull()
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
    if obj, ok := item["graphFingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GraphFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GraphFingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GraphFingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GraphFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.GraphFingerprint = types.StringNull()
        }
    } else if arr, ok := item["graphFingerprint"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.GraphFingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.GraphFingerprint = types.StringNull()
        }
    } else if val, ok := item["graphFingerprint"].(string); ok {
        data.GraphFingerprint = types.StringValue(val)
    } else {
        data.GraphFingerprint = types.StringNull()
    }
    if obj, ok := item["changes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Changes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Changes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Changes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Changes = types.StringValue(string(jsonBytes))
        } else {
            data.Changes = types.StringNull()
        }
    } else if arr, ok := item["changes"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Changes = types.StringValue(string(jsonBytes))
        } else {
            data.Changes = types.StringNull()
        }
    } else if val, ok := item["changes"].(string); ok {
        data.Changes = types.StringValue(val)
    } else {
        data.Changes = types.StringNull()
    }
    data.ChangeCount = numberFromAPI(data.ChangeCount, item["changeCount"])
    data.RestoredFromVersionNumber = numberFromAPI(data.RestoredFromVersionNumber, item["restoredFromVersionNumber"])
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
    if obj, ok := item["savedThrough"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SavedThrough = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SavedThrough = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SavedThrough = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SavedThrough = types.StringValue(string(jsonBytes))
        } else {
            data.SavedThrough = types.StringNull()
        }
    } else if arr, ok := item["savedThrough"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SavedThrough = types.StringValue(string(jsonBytes))
        } else {
            data.SavedThrough = types.StringNull()
        }
    } else if val, ok := item["savedThrough"].(string); ok {
        data.SavedThrough = types.StringValue(val)
    } else {
        data.SavedThrough = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
