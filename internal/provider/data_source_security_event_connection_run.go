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
var _ datasource.DataSource = &SecurityEventConnectionRunDataSource{}

func NewSecurityEventConnectionRunDataSource() datasource.DataSource {
    return &SecurityEventConnectionRunDataSource{}
}

// SecurityEventConnectionRunDataSource defines the data source implementation.
type SecurityEventConnectionRunDataSource struct {
    client *Client
}

// SecurityEventConnectionRunDataSourceModel describes the data source data model.
type SecurityEventConnectionRunDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    SecurityEventConnectionId types.String `tfsdk:"security_event_connection_id"`
    RequestedByUserId types.String `tfsdk:"requested_by_user_id"`
    Type types.String `tfsdk:"type"`
    Status types.String `tfsdk:"status"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    Request types.String `tfsdk:"request"`
    Result types.String `tfsdk:"result"`
    Error types.String `tfsdk:"error"`
}

func (d *SecurityEventConnectionRunDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_event_connection_run"
}

func (d *SecurityEventConnectionRunDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "History of connection tests, previews, scheduled polls and historical imports for security event connections. Credentials are never included. Look up an existing security event connection run by `id`, or by any of its other arguments (`error`, `requested_by_user_id`, `security_event_connection_id`, ...): each one set must match, and exactly one security event connection run may match them all.",

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
                MarkdownDescription: "ID of the project for this run. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_event_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the connection for this run. The ID of a `oneuptime_security_event_connection`.",
                Optional: true,
                Computed: true,
            },
            "requested_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who requested this run. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "type": schema.StringAttribute{
                MarkdownDescription: "Operation of this connection run.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this connection run.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When this run started.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "When this run completed.",
                Computed: true,
            },
            "request": schema.StringAttribute{
                MarkdownDescription: "Validated operation and selected time range. Contains no credentials. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "result": schema.StringAttribute{
                MarkdownDescription: "Counts, requested time range, checks and a bounded preview of records. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "error": schema.StringAttribute{
                MarkdownDescription: "The run failure with credentials redacted.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *SecurityEventConnectionRunDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityEventConnectionRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityEventConnectionRunDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.SecurityEventConnectionId.IsNull() && !data.SecurityEventConnectionId.IsUnknown() {
        filters["securityEventConnectionId"] = data.SecurityEventConnectionId.ValueString()
        filterNames = append(filterNames, "security_event_connection_id = "+fmt.Sprintf("%q", data.SecurityEventConnectionId.ValueString()))
    }
    if !data.RequestedByUserId.IsNull() && !data.RequestedByUserId.IsUnknown() {
        filters["requestedByUserId"] = data.RequestedByUserId.ValueString()
        filterNames = append(filterNames, "requested_by_user_id = "+fmt.Sprintf("%q", data.RequestedByUserId.ValueString()))
    }
    if !data.Type.IsNull() && !data.Type.IsUnknown() {
        filters["type"] = data.Type.ValueString()
        filterNames = append(filterNames, "type = "+fmt.Sprintf("%q", data.Type.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.Error.IsNull() && !data.Error.IsUnknown() {
        filters["error"] = data.Error.ValueString()
        filterNames = append(filterNames, "error = "+fmt.Sprintf("%q", data.Error.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security event connection run up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security event connection run up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "securityEventConnectionId": true,
        "requestedByUserId": true,
        "type": true,
        "status": true,
        "startedAt": true,
        "completedAt": true,
        "request": true,
        "result": true,
        "error": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-event-connection-run/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_event_connection_run, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security event connection run found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_event_connection_run: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-event-connection-run/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_event_connection_run, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_event_connection_run: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security event connection run matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security event connection run matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_event_connection_run.")
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
    if obj, ok := item["securityEventConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityEventConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityEventConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityEventConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityEventConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityEventConnectionId = types.StringNull()
        }
    } else if val, ok := item["securityEventConnectionId"].(string); ok {
        data.SecurityEventConnectionId = types.StringValue(val)
    } else {
        data.SecurityEventConnectionId = types.StringNull()
    }
    if obj, ok := item["requestedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RequestedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RequestedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RequestedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.RequestedByUserId = types.StringNull()
        }
    } else if val, ok := item["requestedByUserId"].(string); ok {
        data.RequestedByUserId = types.StringValue(val)
    } else {
        data.RequestedByUserId = types.StringNull()
    }
    if obj, ok := item["type"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Type = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Type = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Type = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Type = types.StringValue(string(jsonBytes))
        } else {
            data.Type = types.StringNull()
        }
    } else if val, ok := item["type"].(string); ok {
        data.Type = types.StringValue(val)
    } else {
        data.Type = types.StringNull()
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
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartedAt = types.StringNull()
        }
    } else if val, ok := item["startedAt"].(string); ok {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if obj, ok := item["completedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CompletedAt = types.StringNull()
        }
    } else if val, ok := item["completedAt"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    } else {
        data.CompletedAt = types.StringNull()
    }
    if obj, ok := item["request"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Request = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Request = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Request = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Request = types.StringValue(string(jsonBytes))
        } else {
            data.Request = types.StringNull()
        }
    } else if val, ok := item["request"].(string); ok {
        data.Request = types.StringValue(val)
    } else {
        data.Request = types.StringNull()
    }
    if obj, ok := item["result"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Result = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Result = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Result = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Result = types.StringValue(string(jsonBytes))
        } else {
            data.Result = types.StringNull()
        }
    } else if val, ok := item["result"].(string); ok {
        data.Result = types.StringValue(val)
    } else {
        data.Result = types.StringNull()
    }
    if obj, ok := item["error"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Error = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Error = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Error = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Error = types.StringValue(string(jsonBytes))
        } else {
            data.Error = types.StringNull()
        }
    } else if val, ok := item["error"].(string); ok {
        data.Error = types.StringValue(val)
    } else {
        data.Error = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
