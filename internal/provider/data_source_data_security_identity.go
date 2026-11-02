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
var _ datasource.DataSource = &DataSecurityIdentityDataSource{}

func NewDataSecurityIdentityDataSource() datasource.DataSource {
    return &DataSecurityIdentityDataSource{}
}

// DataSecurityIdentityDataSource defines the data source implementation.
type DataSecurityIdentityDataSource struct {
    client *Client
}

// DataSecurityIdentityDataSourceModel describes the data source data model.
type DataSecurityIdentityDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IdentityKey types.String `tfsdk:"identity_key"`
    Principal types.String `tfsdk:"principal"`
    PrincipalName types.String `tfsdk:"principal_name"`
    PrincipalType types.String `tfsdk:"principal_type"`
    ProviderValue types.String `tfsdk:"provider_value"`
    IsExternal types.Bool `tfsdk:"is_external"`
    IsPublic types.Bool `tfsdk:"is_public"`
    IsUnused types.Bool `tfsdk:"is_unused"`
    IsSeenInTraces types.Bool `tfsdk:"is_seen_in_traces"`
    StoreCount types.Number `tfsdk:"store_count"`
    SensitiveStoreCount types.Number `tfsdk:"sensitive_store_count"`
    RestrictedStoreCount types.Number `tfsdk:"restricted_store_count"`
    HighestSensitivity types.String `tfsdk:"highest_sensitivity"`
    HighestAccessLevel types.String `tfsdk:"highest_access_level"`
    CanChangeSensitiveData types.Bool `tfsdk:"can_change_sensitive_data"`
    DataTypeIds types.String `tfsdk:"data_type_ids"`
    EstimatedRecords types.Number `tfsdk:"estimated_records"`
    RiskScore types.Number `tfsdk:"risk_score"`
    LastUsedAt types.String `tfsdk:"last_used_at"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
}

func (d *DataSecurityIdentityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_security_identity"
}

func (d *DataSecurityIdentityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Every identity that can reach the project's data stores, with what it reaches, how sensitive that is, and how much harm it could do. Look up an existing data security identity by `id`, or by any of its other arguments (`can_change_sensitive_data`, `estimated_records`, `highest_access_level`, ...): each one set must match, and exactly one data security identity may match them all.",

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
                MarkdownDescription: "ID of the project this identity belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "identity_key": schema.StringAttribute{
                MarkdownDescription: "The principal as its cloud writes it, which identifies the identity within its project.",
                Optional: true,
                Computed: true,
            },
            "principal": schema.StringAttribute{
                MarkdownDescription: "Who or what it is, as its cloud names it: *, an account id, an ARN, an email, a service account.",
                Optional: true,
                Computed: true,
            },
            "principal_name": schema.StringAttribute{
                MarkdownDescription: "The identity's name, for people.",
                Optional: true,
                Computed: true,
            },
            "principal_type": schema.StringAttribute{
                MarkdownDescription: "What it is: Public, Account, Role, User, Group, Service Principal, Service Account or Service.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Where the identity lives: AWS, Azure, Google Cloud, or OneUptime for a service seen in traces.",
                Optional: true,
                Computed: true,
            },
            "is_external": schema.BoolAttribute{
                MarkdownDescription: "The identity is outside the project's connected accounts.",
                Optional: true,
                Computed: true,
            },
            "is_public": schema.BoolAttribute{
                MarkdownDescription: "The identity is anyone on the internet, or any AWS account.",
                Optional: true,
                Computed: true,
            },
            "is_unused": schema.BoolAttribute{
                MarkdownDescription: "The identity has not been used in 90 days, as far as its cloud can say.",
                Optional: true,
                Computed: true,
            },
            "is_seen_in_traces": schema.BoolAttribute{
                MarkdownDescription: "OneUptime saw the identity use a store in traces.",
                Optional: true,
                Computed: true,
            },
            "store_count": schema.NumberAttribute{
                MarkdownDescription: "How many data stores the identity can reach.",
                Optional: true,
                Computed: true,
            },
            "sensitive_store_count": schema.NumberAttribute{
                MarkdownDescription: "How many of them hold Confidential or Restricted data.",
                Optional: true,
                Computed: true,
            },
            "restricted_store_count": schema.NumberAttribute{
                MarkdownDescription: "How many of them hold Restricted data.",
                Optional: true,
                Computed: true,
            },
            "highest_sensitivity": schema.StringAttribute{
                MarkdownDescription: "How sensitive the most sensitive store it reaches is: Restricted, Confidential, Internal or Public.",
                Optional: true,
                Computed: true,
            },
            "highest_access_level": schema.StringAttribute{
                MarkdownDescription: "The most it may do to any store it reaches: Read, Write or Admin.",
                Optional: true,
                Computed: true,
            },
            "can_change_sensitive_data": schema.BoolAttribute{
                MarkdownDescription: "It may write to, or control, a store holding Confidential or Restricted data.",
                Optional: true,
                Computed: true,
            },
            "data_type_ids": schema.StringAttribute{
                MarkdownDescription: "The data types found in the sensitive stores it reaches, by their ids. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "estimated_records": schema.NumberAttribute{
                MarkdownDescription: "About how many records with sensitive data the stores it reaches hold, where scans could tell.",
                Optional: true,
                Computed: true,
            },
            "risk_score": schema.NumberAttribute{
                MarkdownDescription: "How much harm the identity could do with what it reaches, from 0 to 100.",
                Optional: true,
                Computed: true,
            },
            "last_used_at": schema.StringAttribute{
                MarkdownDescription: "When the identity was last used, where OneUptime can tell.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the identity was first found reaching a store.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the identity was last found reaching a store.",
                Computed: true,
            },
        },
    }
}

func (d *DataSecurityIdentityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSecurityIdentityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataSecurityIdentityDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IdentityKey.IsNull() && !data.IdentityKey.IsUnknown() {
        filters["identityKey"] = data.IdentityKey.ValueString()
        filterNames = append(filterNames, "identity_key = "+fmt.Sprintf("%q", data.IdentityKey.ValueString()))
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
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.IsExternal.IsNull() && !data.IsExternal.IsUnknown() {
        filters["isExternal"] = data.IsExternal.ValueBool()
        filterNames = append(filterNames, "is_external = "+fmt.Sprintf("%t", data.IsExternal.ValueBool()))
    }
    if !data.IsPublic.IsNull() && !data.IsPublic.IsUnknown() {
        filters["isPublic"] = data.IsPublic.ValueBool()
        filterNames = append(filterNames, "is_public = "+fmt.Sprintf("%t", data.IsPublic.ValueBool()))
    }
    if !data.IsUnused.IsNull() && !data.IsUnused.IsUnknown() {
        filters["isUnused"] = data.IsUnused.ValueBool()
        filterNames = append(filterNames, "is_unused = "+fmt.Sprintf("%t", data.IsUnused.ValueBool()))
    }
    if !data.IsSeenInTraces.IsNull() && !data.IsSeenInTraces.IsUnknown() {
        filters["isSeenInTraces"] = data.IsSeenInTraces.ValueBool()
        filterNames = append(filterNames, "is_seen_in_traces = "+fmt.Sprintf("%t", data.IsSeenInTraces.ValueBool()))
    }
    if !data.StoreCount.IsNull() && !data.StoreCount.IsUnknown() {
        filters["storeCount"] = lookupNumber(data.StoreCount)
        filterNames = append(filterNames, "store_count = "+data.StoreCount.ValueBigFloat().String())
    }
    if !data.SensitiveStoreCount.IsNull() && !data.SensitiveStoreCount.IsUnknown() {
        filters["sensitiveStoreCount"] = lookupNumber(data.SensitiveStoreCount)
        filterNames = append(filterNames, "sensitive_store_count = "+data.SensitiveStoreCount.ValueBigFloat().String())
    }
    if !data.RestrictedStoreCount.IsNull() && !data.RestrictedStoreCount.IsUnknown() {
        filters["restrictedStoreCount"] = lookupNumber(data.RestrictedStoreCount)
        filterNames = append(filterNames, "restricted_store_count = "+data.RestrictedStoreCount.ValueBigFloat().String())
    }
    if !data.HighestSensitivity.IsNull() && !data.HighestSensitivity.IsUnknown() {
        filters["highestSensitivity"] = data.HighestSensitivity.ValueString()
        filterNames = append(filterNames, "highest_sensitivity = "+fmt.Sprintf("%q", data.HighestSensitivity.ValueString()))
    }
    if !data.HighestAccessLevel.IsNull() && !data.HighestAccessLevel.IsUnknown() {
        filters["highestAccessLevel"] = data.HighestAccessLevel.ValueString()
        filterNames = append(filterNames, "highest_access_level = "+fmt.Sprintf("%q", data.HighestAccessLevel.ValueString()))
    }
    if !data.CanChangeSensitiveData.IsNull() && !data.CanChangeSensitiveData.IsUnknown() {
        filters["canChangeSensitiveData"] = data.CanChangeSensitiveData.ValueBool()
        filterNames = append(filterNames, "can_change_sensitive_data = "+fmt.Sprintf("%t", data.CanChangeSensitiveData.ValueBool()))
    }
    if !data.EstimatedRecords.IsNull() && !data.EstimatedRecords.IsUnknown() {
        filters["estimatedRecords"] = lookupNumber(data.EstimatedRecords)
        filterNames = append(filterNames, "estimated_records = "+data.EstimatedRecords.ValueBigFloat().String())
    }
    if !data.RiskScore.IsNull() && !data.RiskScore.IsUnknown() {
        filters["riskScore"] = lookupNumber(data.RiskScore)
        filterNames = append(filterNames, "risk_score = "+data.RiskScore.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data security identity up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data security identity up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "identityKey": true,
        "principal": true,
        "principalName": true,
        "principalType": true,
        "provider": true,
        "isExternal": true,
        "isPublic": true,
        "isUnused": true,
        "isSeenInTraces": true,
        "storeCount": true,
        "sensitiveStoreCount": true,
        "restrictedStoreCount": true,
        "highestSensitivity": true,
        "highestAccessLevel": true,
        "canChangeSensitiveData": true,
        "dataTypeIds": true,
        "estimatedRecords": true,
        "riskScore": true,
        "lastUsedAt": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-identity/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_identity, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security identity found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_security_identity: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-identity/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_security_identity, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_security_identity: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security identity matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data security identity matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_security_identity.")
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
    if obj, ok := item["identityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IdentityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IdentityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IdentityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IdentityKey = types.StringValue(string(jsonBytes))
        } else {
            data.IdentityKey = types.StringNull()
        }
    } else if val, ok := item["identityKey"].(string); ok {
        data.IdentityKey = types.StringValue(val)
    } else {
        data.IdentityKey = types.StringNull()
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
    if val, ok := item["isUnused"].(bool); ok {
        data.IsUnused = types.BoolValue(val)
    } else {
        data.IsUnused = types.BoolNull()
    }
    if val, ok := item["isSeenInTraces"].(bool); ok {
        data.IsSeenInTraces = types.BoolValue(val)
    } else {
        data.IsSeenInTraces = types.BoolNull()
    }
    data.StoreCount = numberFromAPI(data.StoreCount, item["storeCount"])
    data.SensitiveStoreCount = numberFromAPI(data.SensitiveStoreCount, item["sensitiveStoreCount"])
    data.RestrictedStoreCount = numberFromAPI(data.RestrictedStoreCount, item["restrictedStoreCount"])
    if obj, ok := item["highestSensitivity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighestSensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HighestSensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HighestSensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HighestSensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.HighestSensitivity = types.StringNull()
        }
    } else if val, ok := item["highestSensitivity"].(string); ok {
        data.HighestSensitivity = types.StringValue(val)
    } else {
        data.HighestSensitivity = types.StringNull()
    }
    if obj, ok := item["highestAccessLevel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighestAccessLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HighestAccessLevel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HighestAccessLevel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HighestAccessLevel = types.StringValue(string(jsonBytes))
        } else {
            data.HighestAccessLevel = types.StringNull()
        }
    } else if val, ok := item["highestAccessLevel"].(string); ok {
        data.HighestAccessLevel = types.StringValue(val)
    } else {
        data.HighestAccessLevel = types.StringNull()
    }
    if val, ok := item["canChangeSensitiveData"].(bool); ok {
        data.CanChangeSensitiveData = types.BoolValue(val)
    } else {
        data.CanChangeSensitiveData = types.BoolNull()
    }
    if obj, ok := item["dataTypeIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataTypeIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataTypeIds = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeIds = types.StringNull()
        }
    } else if val, ok := item["dataTypeIds"].(string); ok {
        data.DataTypeIds = types.StringValue(val)
    } else {
        data.DataTypeIds = types.StringNull()
    }
    data.EstimatedRecords = numberFromAPI(data.EstimatedRecords, item["estimatedRecords"])
    data.RiskScore = numberFromAPI(data.RiskScore, item["riskScore"])
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
