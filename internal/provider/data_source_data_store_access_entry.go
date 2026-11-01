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
var _ datasource.DataSource = &DataStoreAccessEntryDataSource{}

func NewDataStoreAccessEntryDataSource() datasource.DataSource {
    return &DataStoreAccessEntryDataSource{}
}

// DataStoreAccessEntryDataSource defines the data source implementation.
type DataStoreAccessEntryDataSource struct {
    client *Client
}

// DataStoreAccessEntryDataSourceModel describes the data source data model.
type DataStoreAccessEntryDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    AccessKey types.String `tfsdk:"access_key"`
    Principal types.String `tfsdk:"principal"`
    PrincipalName types.String `tfsdk:"principal_name"`
    PrincipalType types.String `tfsdk:"principal_type"`
    AccessLevel types.String `tfsdk:"access_level"`
    GrantedVia types.String `tfsdk:"granted_via"`
    IsExternal types.Bool `tfsdk:"is_external"`
    IsPublic types.Bool `tfsdk:"is_public"`
    IsSeenInTraces types.Bool `tfsdk:"is_seen_in_traces"`
    LastUsedAt types.String `tfsdk:"last_used_at"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
}

func (d *DataStoreAccessEntryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_store_access_entry"
}

func (d *DataStoreAccessEntryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Who and what can reach each store: the internet, other accounts, identities and the services seen using it, with the policy that grants it. Look up an existing data store access entry by `id`, or by any of its other arguments (`access_key`, `access_level`, `data_store_id`, ...): each one set must match, and exactly one data store access entry may match them all.",

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
                MarkdownDescription: "ID of the project this access entry belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store that can be reached. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
            },
            "access_key": schema.StringAttribute{
                MarkdownDescription: "The principal, level and grant, which identify the entry within its store.",
                Optional: true,
                Computed: true,
            },
            "principal": schema.StringAttribute{
                MarkdownDescription: "Who or what has access, as its cloud names it: *, an account id, an ARN, an email, a service account.",
                Optional: true,
                Computed: true,
            },
            "principal_name": schema.StringAttribute{
                MarkdownDescription: "The principal's name, for people.",
                Optional: true,
                Computed: true,
            },
            "principal_type": schema.StringAttribute{
                MarkdownDescription: "What the principal is: Public, Account, Role, User, Group, Service Principal, Service Account or Service.",
                Optional: true,
                Computed: true,
            },
            "access_level": schema.StringAttribute{
                MarkdownDescription: "What it may do: Read, Write or Admin.",
                Optional: true,
                Computed: true,
            },
            "granted_via": schema.StringAttribute{
                MarkdownDescription: "What grants the access: a bucket policy, an ACL, a role assignment, an IAM binding, a snapshot share, or traffic seen in traces.",
                Optional: true,
                Computed: true,
            },
            "is_external": schema.BoolAttribute{
                MarkdownDescription: "The principal is outside the project's connected accounts.",
                Optional: true,
                Computed: true,
            },
            "is_public": schema.BoolAttribute{
                MarkdownDescription: "The principal is anyone on the internet.",
                Optional: true,
                Computed: true,
            },
            "is_seen_in_traces": schema.BoolAttribute{
                MarkdownDescription: "OneUptime saw this principal use the store in traces.",
                Optional: true,
                Computed: true,
            },
            "last_used_at": schema.StringAttribute{
                MarkdownDescription: "When the access was last used, where OneUptime can tell.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the access was first read.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the access was last read.",
                Computed: true,
            },
        },
    }
}

func (d *DataStoreAccessEntryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataStoreAccessEntryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataStoreAccessEntryDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.AccessKey.IsNull() && !data.AccessKey.IsUnknown() {
        filters["accessKey"] = data.AccessKey.ValueString()
        filterNames = append(filterNames, "access_key = "+fmt.Sprintf("%q", data.AccessKey.ValueString()))
    }
    if !data.Principal.IsNull() && !data.Principal.IsUnknown() {
        filters["principal"] = data.Principal.ValueString()
        filterNames = append(filterNames, "principal = "+fmt.Sprintf("%q", data.Principal.ValueString()))
    }
    if !data.PrincipalName.IsNull() && !data.PrincipalName.IsUnknown() {
        filters["principalName"] = data.PrincipalName.ValueString()
        filterNames = append(filterNames, "principal_name = "+fmt.Sprintf("%q", data.PrincipalName.ValueString()))
    }
    if !data.PrincipalType.IsNull() && !data.PrincipalType.IsUnknown() {
        filters["principalType"] = data.PrincipalType.ValueString()
        filterNames = append(filterNames, "principal_type = "+fmt.Sprintf("%q", data.PrincipalType.ValueString()))
    }
    if !data.AccessLevel.IsNull() && !data.AccessLevel.IsUnknown() {
        filters["accessLevel"] = data.AccessLevel.ValueString()
        filterNames = append(filterNames, "access_level = "+fmt.Sprintf("%q", data.AccessLevel.ValueString()))
    }
    if !data.GrantedVia.IsNull() && !data.GrantedVia.IsUnknown() {
        filters["grantedVia"] = data.GrantedVia.ValueString()
        filterNames = append(filterNames, "granted_via = "+fmt.Sprintf("%q", data.GrantedVia.ValueString()))
    }
    if !data.IsExternal.IsNull() && !data.IsExternal.IsUnknown() {
        filters["isExternal"] = data.IsExternal.ValueBool()
        filterNames = append(filterNames, "is_external = "+fmt.Sprintf("%t", data.IsExternal.ValueBool()))
    }
    if !data.IsPublic.IsNull() && !data.IsPublic.IsUnknown() {
        filters["isPublic"] = data.IsPublic.ValueBool()
        filterNames = append(filterNames, "is_public = "+fmt.Sprintf("%t", data.IsPublic.ValueBool()))
    }
    if !data.IsSeenInTraces.IsNull() && !data.IsSeenInTraces.IsUnknown() {
        filters["isSeenInTraces"] = data.IsSeenInTraces.ValueBool()
        filterNames = append(filterNames, "is_seen_in_traces = "+fmt.Sprintf("%t", data.IsSeenInTraces.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data store access entry up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data store access entry up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataStoreId": true,
        "accessKey": true,
        "principal": true,
        "principalName": true,
        "principalType": true,
        "accessLevel": true,
        "grantedVia": true,
        "isExternal": true,
        "isPublic": true,
        "isSeenInTraces": true,
        "lastUsedAt": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-store-access/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_store_access_entry, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store access entry found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_store_access_entry: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-store-access/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_store_access_entry, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_store_access_entry: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store access entry matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data store access entry matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_store_access_entry.")
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
    if obj, ok := item["dataStoreId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataStoreId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreId = types.StringNull()
        }
    } else if val, ok := item["dataStoreId"].(string); ok {
        data.DataStoreId = types.StringValue(val)
    } else {
        data.DataStoreId = types.StringNull()
    }
    if obj, ok := item["accessKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AccessKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AccessKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AccessKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AccessKey = types.StringValue(string(jsonBytes))
        } else {
            data.AccessKey = types.StringNull()
        }
    } else if val, ok := item["accessKey"].(string); ok {
        data.AccessKey = types.StringValue(val)
    } else {
        data.AccessKey = types.StringNull()
    }
    if obj, ok := item["principal"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Principal = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Principal = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Principal = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Principal = types.StringValue(string(jsonBytes))
        } else {
            data.Principal = types.StringNull()
        }
    } else if val, ok := item["principal"].(string); ok {
        data.Principal = types.StringValue(val)
    } else {
        data.Principal = types.StringNull()
    }
    if obj, ok := item["principalName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrincipalName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrincipalName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrincipalName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrincipalName = types.StringValue(string(jsonBytes))
        } else {
            data.PrincipalName = types.StringNull()
        }
    } else if val, ok := item["principalName"].(string); ok {
        data.PrincipalName = types.StringValue(val)
    } else {
        data.PrincipalName = types.StringNull()
    }
    if obj, ok := item["principalType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrincipalType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrincipalType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrincipalType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrincipalType = types.StringValue(string(jsonBytes))
        } else {
            data.PrincipalType = types.StringNull()
        }
    } else if val, ok := item["principalType"].(string); ok {
        data.PrincipalType = types.StringValue(val)
    } else {
        data.PrincipalType = types.StringNull()
    }
    if obj, ok := item["accessLevel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AccessLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AccessLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AccessLevel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AccessLevel = types.StringValue(string(jsonBytes))
        } else {
            data.AccessLevel = types.StringNull()
        }
    } else if val, ok := item["accessLevel"].(string); ok {
        data.AccessLevel = types.StringValue(val)
    } else {
        data.AccessLevel = types.StringNull()
    }
    if obj, ok := item["grantedVia"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GrantedVia = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GrantedVia = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GrantedVia = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GrantedVia = types.StringValue(string(jsonBytes))
        } else {
            data.GrantedVia = types.StringNull()
        }
    } else if val, ok := item["grantedVia"].(string); ok {
        data.GrantedVia = types.StringValue(val)
    } else {
        data.GrantedVia = types.StringNull()
    }
    if val, ok := item["isExternal"].(bool); ok {
        data.IsExternal = types.BoolValue(val)
    } else {
        data.IsExternal = types.BoolNull()
    }
    if val, ok := item["isPublic"].(bool); ok {
        data.IsPublic = types.BoolValue(val)
    } else {
        data.IsPublic = types.BoolNull()
    }
    if val, ok := item["isSeenInTraces"].(bool); ok {
        data.IsSeenInTraces = types.BoolValue(val)
    } else {
        data.IsSeenInTraces = types.BoolNull()
    }
    if obj, ok := item["lastUsedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastUsedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastUsedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastUsedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastUsedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastUsedAt = types.StringNull()
        }
    } else if val, ok := item["lastUsedAt"].(string); ok {
        data.LastUsedAt = types.StringValue(val)
    } else {
        data.LastUsedAt = types.StringNull()
    }
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
