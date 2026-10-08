package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StatusPageScimLogDataSource{}

func NewStatusPageScimLogDataSource() datasource.DataSource {
    return &StatusPageScimLogDataSource{}
}

// StatusPageScimLogDataSource defines the data source implementation.
type StatusPageScimLogDataSource struct {
    client *Client
}

// StatusPageScimLogDataSourceModel describes the data source data model.
type StatusPageScimLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageScimId types.String `tfsdk:"status_page_scim_id"`
    OperationType types.String `tfsdk:"operation_type"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    LogBody types.String `tfsdk:"log_body"`
    HttpMethod types.String `tfsdk:"http_method"`
    RequestPath types.String `tfsdk:"request_path"`
    HttpStatusCode types.Number `tfsdk:"http_status_code"`
    AffectedUserEmail types.String `tfsdk:"affected_user_email"`
}

func (d *StatusPageScimLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_scim_log"
}

func (d *StatusPageScimLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Logs of all SCIM provisioning operations for status pages. Look up an existing status page scim log by `id`, or by any of its other arguments (`http_method`, `http_status_code`, `log_body`, ...): each one set must match, and exactly one status page scim log may match them all.",

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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Status Page. The ID of a `oneuptime_status_page`.",
                Optional: true,
                Computed: true,
            },
            "status_page_scim_id": schema.StringAttribute{
                MarkdownDescription: "ID of your Status Page SCIM configuration. The ID of a `oneuptime_status_page_scim`.",
                Optional: true,
                Computed: true,
            },
            "operation_type": schema.StringAttribute{
                MarkdownDescription: "Type of SCIM operation (e.g., CreateUser, UpdateUser, DeleteUser, ListUsers, GetUser, BulkOperation).",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the SCIM operation.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Short error or status description.",
                Optional: true,
                Computed: true,
            },
            "log_body": schema.StringAttribute{
                MarkdownDescription: "Detailed JSON with request/response data.",
                Optional: true,
                Computed: true,
            },
            "http_method": schema.StringAttribute{
                MarkdownDescription: "HTTP method used (GET, POST, PUT, PATCH, DELETE).",
                Optional: true,
                Computed: true,
            },
            "request_path": schema.StringAttribute{
                MarkdownDescription: "The SCIM endpoint path.",
                Optional: true,
                Computed: true,
            },
            "http_status_code": schema.NumberAttribute{
                MarkdownDescription: "Response HTTP status code.",
                Optional: true,
                Computed: true,
            },
            "affected_user_email": schema.StringAttribute{
                MarkdownDescription: "Email of the user affected by this operation.",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageScimLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageScimLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageScimLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.StatusPageId.IsNull() && !data.StatusPageId.IsUnknown() {
        filters["statusPageId"] = data.StatusPageId.ValueString()
        filterNames = append(filterNames, "status_page_id = "+fmt.Sprintf("%q", data.StatusPageId.ValueString()))
    }
    if !data.StatusPageScimId.IsNull() && !data.StatusPageScimId.IsUnknown() {
        filters["statusPageScimId"] = data.StatusPageScimId.ValueString()
        filterNames = append(filterNames, "status_page_scim_id = "+fmt.Sprintf("%q", data.StatusPageScimId.ValueString()))
    }
    if !data.OperationType.IsNull() && !data.OperationType.IsUnknown() {
        filters["operationType"] = data.OperationType.ValueString()
        filterNames = append(filterNames, "operation_type = "+fmt.Sprintf("%q", data.OperationType.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.LogBody.IsNull() && !data.LogBody.IsUnknown() {
        filters["logBody"] = data.LogBody.ValueString()
        filterNames = append(filterNames, "log_body = "+fmt.Sprintf("%q", data.LogBody.ValueString()))
    }
    if !data.HttpMethod.IsNull() && !data.HttpMethod.IsUnknown() {
        filters["httpMethod"] = data.HttpMethod.ValueString()
        filterNames = append(filterNames, "http_method = "+fmt.Sprintf("%q", data.HttpMethod.ValueString()))
    }
    if !data.RequestPath.IsNull() && !data.RequestPath.IsUnknown() {
        filters["requestPath"] = data.RequestPath.ValueString()
        filterNames = append(filterNames, "request_path = "+fmt.Sprintf("%q", data.RequestPath.ValueString()))
    }
    if !data.HttpStatusCode.IsNull() && !data.HttpStatusCode.IsUnknown() {
        filters["httpStatusCode"] = lookupNumber(data.HttpStatusCode)
        filterNames = append(filterNames, "http_status_code = "+data.HttpStatusCode.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the status page scim log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page scim log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "statusPageId": true,
        "statusPageScimId": true,
        "operationType": true,
        "status": true,
        "statusMessage": true,
        "logBody": true,
        "httpMethod": true,
        "requestPath": true,
        "httpStatusCode": true,
        "affectedUserEmail": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-scim-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_scim_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page scim log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_scim_log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-scim-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_scim_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_scim_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page scim log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page scim log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_scim_log.")
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
    if obj, ok := item["statusPageId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := item["statusPageId"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := item["statusPageScimId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageScimId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageScimId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageScimId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageScimId = types.StringNull()
        }
    } else if val, ok := item["statusPageScimId"].(string); ok {
        data.StatusPageScimId = types.StringValue(val)
    } else {
        data.StatusPageScimId = types.StringNull()
    }
    if obj, ok := item["operationType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OperationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OperationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OperationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OperationType = types.StringValue(string(jsonBytes))
        } else {
            data.OperationType = types.StringNull()
        }
    } else if val, ok := item["operationType"].(string); ok {
        data.OperationType = types.StringValue(val)
    } else {
        data.OperationType = types.StringNull()
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
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["logBody"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogBody = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LogBody = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LogBody = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LogBody = types.StringValue(string(jsonBytes))
        } else {
            data.LogBody = types.StringNull()
        }
    } else if val, ok := item["logBody"].(string); ok {
        data.LogBody = types.StringValue(val)
    } else {
        data.LogBody = types.StringNull()
    }
    if obj, ok := item["httpMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HttpMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HttpMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HttpMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HttpMethod = types.StringValue(string(jsonBytes))
        } else {
            data.HttpMethod = types.StringNull()
        }
    } else if val, ok := item["httpMethod"].(string); ok {
        data.HttpMethod = types.StringValue(val)
    } else {
        data.HttpMethod = types.StringNull()
    }
    if obj, ok := item["requestPath"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestPath = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RequestPath = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RequestPath = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RequestPath = types.StringValue(string(jsonBytes))
        } else {
            data.RequestPath = types.StringNull()
        }
    } else if val, ok := item["requestPath"].(string); ok {
        data.RequestPath = types.StringValue(val)
    } else {
        data.RequestPath = types.StringNull()
    }
    if val, ok := item["httpStatusCode"].(float64); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["httpStatusCode"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.HttpStatusCode = types.NumberValue(big.NewFloat(val))
        } else {
            data.HttpStatusCode = types.NumberNull()
        }
    } else {
        data.HttpStatusCode = types.NumberNull()
    }
    if obj, ok := item["affectedUserEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AffectedUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AffectedUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AffectedUserEmail = types.StringValue(string(jsonBytes))
        } else {
            data.AffectedUserEmail = types.StringNull()
        }
    } else if val, ok := item["affectedUserEmail"].(string); ok {
        data.AffectedUserEmail = types.StringValue(val)
    } else {
        data.AffectedUserEmail = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
