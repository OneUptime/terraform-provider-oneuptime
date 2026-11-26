package shard19

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
var _ datasource.DataSource = &SecurityCaseEvidenceDataSource{}

func NewSecurityCaseEvidenceDataSource() datasource.DataSource {
    return &SecurityCaseEvidenceDataSource{}
}

// SecurityCaseEvidenceDataSource defines the data source implementation.
type SecurityCaseEvidenceDataSource struct {
    client *Client
}

// SecurityCaseEvidenceDataSourceModel describes the data source data model.
type SecurityCaseEvidenceDataSourceModel struct {
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
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    Kind types.String `tfsdk:"kind"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    FileId types.String `tfsdk:"file_id"`
    Url types.String `tfsdk:"url"`
    FileName types.String `tfsdk:"file_name"`
    FileType types.String `tfsdk:"file_type"`
    FileSizeInBytes types.Number `tfsdk:"file_size_in_bytes"`
    Sha256 types.String `tfsdk:"sha256"`
    EventUid types.String `tfsdk:"event_uid"`
    EventTime types.String `tfsdk:"event_time"`
    EventSnapshot types.String `tfsdk:"event_snapshot"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityCaseEvidenceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_case_evidence"
}

func (d *SecurityCaseEvidenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Files, links and security events kept on a security case as evidence, each with who added it and when: a file with the SHA-256 of what was uploaded, an event as OneUptime read it when it was kept. Look up an existing security case evidence by `id`, or by any of its other arguments (`name`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one security case evidence may match them all.",

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
                MarkdownDescription: "ID of the project this evidence belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case this evidence is kept on. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What the evidence is: a File uploaded to OneUptime, a Link to evidence kept somewhere else, or an Event bookmarked from Security Events. It can't be changed.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the evidence is, in a few words. A file's own name when none is given.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the evidence shows, and where and how it was collected.",
                Optional: true,
                Computed: true,
            },
            "file_id": schema.StringAttribute{
                MarkdownDescription: "ID of the uploaded file, for File evidence: upload it first, then name it here. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_file"),
            },
            "url": schema.StringAttribute{
                MarkdownDescription: "Where the evidence is kept, for Link evidence: an http or https address.",
                Optional: true,
                Computed: true,
            },
            "file_name": schema.StringAttribute{
                MarkdownDescription: "The file's name as it was uploaded, for File evidence.",
                Optional: true,
                Computed: true,
            },
            "file_type": schema.StringAttribute{
                MarkdownDescription: "The file's type as it was uploaded, like text/plain.",
                Optional: true,
                Computed: true,
            },
            "file_size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "The file's size as it was uploaded, in bytes.",
                Optional: true,
                Computed: true,
            },
            "sha256": schema.StringAttribute{
                MarkdownDescription: "The SHA-256 of the file's bytes as they were uploaded, in hex: compare it with a download's to check it is what was kept.",
                Optional: true,
                Computed: true,
            },
            "event_uid": schema.StringAttribute{
                MarkdownDescription: "For Event evidence: the security event's ID, as its source gave it.",
                Optional: true,
                Computed: true,
            },
            "event_time": schema.StringAttribute{
                MarkdownDescription: "For Event evidence: when the security event happened.",
                Computed: true,
            },
            "event_snapshot": schema.StringAttribute{
                MarkdownDescription: "For Event evidence: the security event as OneUptime read it when it was kept - its severity, class, message, who and where - and the hunt that found it. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who added the evidence. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *SecurityCaseEvidenceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityCaseEvidenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityCaseEvidenceDataSourceModel

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
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.FileId.IsNull() && !data.FileId.IsUnknown() {
        filters["fileId"] = data.FileId.ValueString()
        filterNames = append(filterNames, "file_id = "+fmt.Sprintf("%q", data.FileId.ValueString()))
    }
    if !data.Url.IsNull() && !data.Url.IsUnknown() {
        filters["url"] = data.Url.ValueString()
        filterNames = append(filterNames, "url = "+fmt.Sprintf("%q", data.Url.ValueString()))
    }
    if !data.FileName.IsNull() && !data.FileName.IsUnknown() {
        filters["fileName"] = data.FileName.ValueString()
        filterNames = append(filterNames, "file_name = "+fmt.Sprintf("%q", data.FileName.ValueString()))
    }
    if !data.FileType.IsNull() && !data.FileType.IsUnknown() {
        filters["fileType"] = data.FileType.ValueString()
        filterNames = append(filterNames, "file_type = "+fmt.Sprintf("%q", data.FileType.ValueString()))
    }
    if !data.FileSizeInBytes.IsNull() && !data.FileSizeInBytes.IsUnknown() {
        filters["fileSizeInBytes"] = LookupNumber(data.FileSizeInBytes)
        filterNames = append(filterNames, "file_size_in_bytes = "+data.FileSizeInBytes.ValueBigFloat().String())
    }
    if !data.Sha256.IsNull() && !data.Sha256.IsUnknown() {
        filters["sha256"] = data.Sha256.ValueString()
        filterNames = append(filterNames, "sha256 = "+fmt.Sprintf("%q", data.Sha256.ValueString()))
    }
    if !data.EventUid.IsNull() && !data.EventUid.IsUnknown() {
        filters["eventUid"] = data.EventUid.ValueString()
        filterNames = append(filterNames, "event_uid = "+fmt.Sprintf("%q", data.EventUid.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security case evidence up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security case evidence up by.",
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
        "securityCaseId": true,
        "kind": true,
        "name": true,
        "description": true,
        "fileId": true,
        "url": true,
        "fileName": true,
        "fileType": true,
        "fileSizeInBytes": true,
        "sha256": true,
        "eventUid": true,
        "eventTime": true,
        "eventSnapshot": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-case-evidence/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case_evidence, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case evidence found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_case_evidence: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-case-evidence/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_case_evidence, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_case_evidence: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case evidence matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security case evidence matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_case_evidence.")
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
    if obj, ok := item["securityCaseId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := item["securityCaseId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := item["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
    }
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if arr, ok := item["kind"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
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
    if obj, ok := item["fileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FileId = types.StringValue(string(jsonBytes))
        } else {
            data.FileId = types.StringNull()
        }
    } else if arr, ok := item["fileId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FileId = types.StringValue(string(jsonBytes))
        } else {
            data.FileId = types.StringNull()
        }
    } else if val, ok := item["fileId"].(string); ok {
        data.FileId = types.StringValue(val)
    } else {
        data.FileId = types.StringNull()
    }
    if obj, ok := item["url"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Url = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Url = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Url = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Url = types.StringValue(string(jsonBytes))
        } else {
            data.Url = types.StringNull()
        }
    } else if arr, ok := item["url"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Url = types.StringValue(string(jsonBytes))
        } else {
            data.Url = types.StringNull()
        }
    } else if val, ok := item["url"].(string); ok {
        data.Url = types.StringValue(val)
    } else {
        data.Url = types.StringNull()
    }
    if obj, ok := item["fileName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FileName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FileName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FileName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FileName = types.StringValue(string(jsonBytes))
        } else {
            data.FileName = types.StringNull()
        }
    } else if arr, ok := item["fileName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FileName = types.StringValue(string(jsonBytes))
        } else {
            data.FileName = types.StringNull()
        }
    } else if val, ok := item["fileName"].(string); ok {
        data.FileName = types.StringValue(val)
    } else {
        data.FileName = types.StringNull()
    }
    if obj, ok := item["fileType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FileType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FileType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FileType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FileType = types.StringValue(string(jsonBytes))
        } else {
            data.FileType = types.StringNull()
        }
    } else if arr, ok := item["fileType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FileType = types.StringValue(string(jsonBytes))
        } else {
            data.FileType = types.StringNull()
        }
    } else if val, ok := item["fileType"].(string); ok {
        data.FileType = types.StringValue(val)
    } else {
        data.FileType = types.StringNull()
    }
    data.FileSizeInBytes = NumberFromAPI(data.FileSizeInBytes, item["fileSizeInBytes"])
    if obj, ok := item["sha256"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sha256 = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Sha256 = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Sha256 = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Sha256 = types.StringValue(string(jsonBytes))
        } else {
            data.Sha256 = types.StringNull()
        }
    } else if arr, ok := item["sha256"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Sha256 = types.StringValue(string(jsonBytes))
        } else {
            data.Sha256 = types.StringNull()
        }
    } else if val, ok := item["sha256"].(string); ok {
        data.Sha256 = types.StringValue(val)
    } else {
        data.Sha256 = types.StringNull()
    }
    if obj, ok := item["eventUid"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventUid = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventUid = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventUid = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventUid = types.StringValue(string(jsonBytes))
        } else {
            data.EventUid = types.StringNull()
        }
    } else if arr, ok := item["eventUid"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EventUid = types.StringValue(string(jsonBytes))
        } else {
            data.EventUid = types.StringNull()
        }
    } else if val, ok := item["eventUid"].(string); ok {
        data.EventUid = types.StringValue(val)
    } else {
        data.EventUid = types.StringNull()
    }
    if obj, ok := item["eventTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventTime = types.StringValue(string(jsonBytes))
        } else {
            data.EventTime = types.StringNull()
        }
    } else if arr, ok := item["eventTime"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EventTime = types.StringValue(string(jsonBytes))
        } else {
            data.EventTime = types.StringNull()
        }
    } else if val, ok := item["eventTime"].(string); ok {
        data.EventTime = types.StringValue(val)
    } else {
        data.EventTime = types.StringNull()
    }
    if obj, ok := item["eventSnapshot"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventSnapshot = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.EventSnapshot = types.StringNull()
        }
    } else if arr, ok := item["eventSnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EventSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.EventSnapshot = types.StringNull()
        }
    } else if val, ok := item["eventSnapshot"].(string); ok {
        data.EventSnapshot = types.StringValue(val)
    } else {
        data.EventSnapshot = types.StringNull()
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
