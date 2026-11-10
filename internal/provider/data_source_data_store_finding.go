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
var _ datasource.DataSource = &DataStoreFindingDataSource{}

func NewDataStoreFindingDataSource() datasource.DataSource {
    return &DataStoreFindingDataSource{}
}

// DataStoreFindingDataSource defines the data source implementation.
type DataStoreFindingDataSource struct {
    client *Client
}

// DataStoreFindingDataSourceModel describes the data source data model.
type DataStoreFindingDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    FindingKey types.String `tfsdk:"finding_key"`
    DataTypeId types.String `tfsdk:"data_type_id"`
    DataTypeName types.String `tfsdk:"data_type_name"`
    DataCategory types.String `tfsdk:"data_category"`
    Sensitivity types.String `tfsdk:"sensitivity"`
    Location types.String `tfsdk:"location"`
    LocationKind types.String `tfsdk:"location_kind"`
    Confidence types.String `tfsdk:"confidence"`
    MatchCount types.Number `tfsdk:"match_count"`
    ValueCount types.Number `tfsdk:"value_count"`
    MatchPercent types.Number `tfsdk:"match_percent"`
    DistinctValues types.Number `tfsdk:"distinct_values"`
    TestValueCount types.Number `tfsdk:"test_value_count"`
    MaskedSamples types.String `tfsdk:"masked_samples"`
    IsFreeText types.Bool `tfsdk:"is_free_text"`
    IsNameHinted types.Bool `tfsdk:"is_name_hinted"`
    IsFoundByName types.Bool `tfsdk:"is_found_by_name"`
    IsScrubbed types.Bool `tfsdk:"is_scrubbed"`
    IsSuppressed types.Bool `tfsdk:"is_suppressed"`
    SuppressedReason types.String `tfsdk:"suppressed_reason"`
    IsConfirmed types.Bool `tfsdk:"is_confirmed"`
    IsAiSuggested types.Bool `tfsdk:"is_ai_suggested"`
    AiReason types.String `tfsdk:"ai_reason"`
    LastScanId types.String `tfsdk:"last_scan_id"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
}

func (d *DataStoreFindingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_store_finding"
}

func (d *DataStoreFindingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What Data Security found in a store: one data type at one location, with its confidence, counts and masked samples. Look up an existing data store finding by `id`, or by any of its other arguments (`ai_reason`, `confidence`, `data_category`, ...): each one set must match, and exactly one data store finding may match them all.",

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
                MarkdownDescription: "ID of the project this finding belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store the finding is in. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_store"),
            },
            "finding_key": schema.StringAttribute{
                MarkdownDescription: "The data type and location, which identify the finding within its store.",
                Optional: true,
                Computed: true,
            },
            "data_type_id": schema.StringAttribute{
                MarkdownDescription: "The data type found: a built-in id such as credit-card-number, or custom: and a custom type's id.",
                Optional: true,
                Computed: true,
            },
            "data_type_name": schema.StringAttribute{
                MarkdownDescription: "The data type's name when it was found.",
                Optional: true,
                Computed: true,
            },
            "data_category": schema.StringAttribute{
                MarkdownDescription: "The data type's category.",
                Optional: true,
                Computed: true,
            },
            "sensitivity": schema.StringAttribute{
                MarkdownDescription: "The data type's sensitivity.",
                Optional: true,
                Computed: true,
            },
            "location": schema.StringAttribute{
                MarkdownDescription: "Where in the store it was found: a table and column, a field path, an object path pattern, or a telemetry attribute.",
                Optional: true,
                Computed: true,
            },
            "location_kind": schema.StringAttribute{
                MarkdownDescription: "What the location is: Column, Field, Object, Attribute, Body or Text.",
                Optional: true,
                Computed: true,
            },
            "confidence": schema.StringAttribute{
                MarkdownDescription: "How sure the scan is: High, Medium or Low.",
                Optional: true,
                Computed: true,
            },
            "match_count": schema.NumberAttribute{
                MarkdownDescription: "How many sampled values matched.",
                Optional: true,
                Computed: true,
            },
            "value_count": schema.NumberAttribute{
                MarkdownDescription: "How many values were sampled at the location.",
                Optional: true,
                Computed: true,
            },
            "match_percent": schema.NumberAttribute{
                MarkdownDescription: "The share of sampled values that matched, from 0 to 100.",
                Optional: true,
                Computed: true,
            },
            "distinct_values": schema.NumberAttribute{
                MarkdownDescription: "How many different values matched, estimated from hashes that never leave the scan.",
                Optional: true,
                Computed: true,
            },
            "test_value_count": schema.NumberAttribute{
                MarkdownDescription: "How many matches were well-known test values, which are never findings.",
                Optional: true,
                Computed: true,
            },
            "masked_samples": schema.StringAttribute{
                MarkdownDescription: "At most three masked samples, such as j***@e******.com. Never a raw value. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "is_free_text": schema.BoolAttribute{
                MarkdownDescription: "The location holds free text, so the finding is about values inside it rather than the whole field.",
                Optional: true,
                Computed: true,
            },
            "is_name_hinted": schema.BoolAttribute{
                MarkdownDescription: "The location's name hints at the data type.",
                Optional: true,
                Computed: true,
            },
            "is_found_by_name": schema.BoolAttribute{
                MarkdownDescription: "Found from the location's name alone, for a data type no value pattern can confirm.",
                Optional: true,
                Computed: true,
            },
            "is_scrubbed": schema.BoolAttribute{
                MarkdownDescription: "For telemetry: a scrub rule already masks this data at ingest.",
                Optional: true,
                Computed: true,
            },
            "is_suppressed": schema.BoolAttribute{
                MarkdownDescription: "A person marked the finding as not sensitive. It no longer counts toward the store's sensitivity or checks.",
                Optional: true,
                Computed: true,
            },
            "suppressed_reason": schema.StringAttribute{
                MarkdownDescription: "Why the finding was marked not sensitive.",
                Optional: true,
                Computed: true,
            },
            "is_confirmed": schema.BoolAttribute{
                MarkdownDescription: "A person confirmed the finding holds this data. A confirmed finding counts as High confidence toward the store's sensitivity and checks, whatever the scan or OneUptime AI was sure of.",
                Optional: true,
                Computed: true,
            },
            "is_ai_suggested": schema.BoolAttribute{
                MarkdownDescription: "OneUptime AI suggested this data type for a field no detector placed, from the field's name, the names beside it and the shapes of its values - never a value. It counts only once a person confirms it.",
                Optional: true,
                Computed: true,
            },
            "ai_reason": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI's reason for its suggestion, in a sentence.",
                Optional: true,
                Computed: true,
            },
            "last_scan_id": schema.StringAttribute{
                MarkdownDescription: "The scan that last saw the finding. The ID of a `oneuptime_data_security_scan` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_data_security_scan"),
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a scan first found it.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a scan last found it.",
                Computed: true,
            },
        },
    }
}

func (d *DataStoreFindingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataStoreFindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataStoreFindingDataSourceModel

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
    if !data.FindingKey.IsNull() && !data.FindingKey.IsUnknown() {
        filters["findingKey"] = data.FindingKey.ValueString()
        filterNames = append(filterNames, "finding_key = "+fmt.Sprintf("%q", data.FindingKey.ValueString()))
    }
    if !data.DataTypeId.IsNull() && !data.DataTypeId.IsUnknown() {
        filters["dataTypeId"] = data.DataTypeId.ValueString()
        filterNames = append(filterNames, "data_type_id = "+fmt.Sprintf("%q", data.DataTypeId.ValueString()))
    }
    if !data.DataTypeName.IsNull() && !data.DataTypeName.IsUnknown() {
        filters["dataTypeName"] = data.DataTypeName.ValueString()
        filterNames = append(filterNames, "data_type_name = "+fmt.Sprintf("%q", data.DataTypeName.ValueString()))
    }
    if !data.DataCategory.IsNull() && !data.DataCategory.IsUnknown() {
        filters["dataCategory"] = data.DataCategory.ValueString()
        filterNames = append(filterNames, "data_category = "+fmt.Sprintf("%q", data.DataCategory.ValueString()))
    }
    if !data.Sensitivity.IsNull() && !data.Sensitivity.IsUnknown() {
        filters["sensitivity"] = data.Sensitivity.ValueString()
        filterNames = append(filterNames, "sensitivity = "+fmt.Sprintf("%q", data.Sensitivity.ValueString()))
    }
    if !data.Location.IsNull() && !data.Location.IsUnknown() {
        filters["location"] = data.Location.ValueString()
        filterNames = append(filterNames, "location = "+fmt.Sprintf("%q", data.Location.ValueString()))
    }
    if !data.LocationKind.IsNull() && !data.LocationKind.IsUnknown() {
        filters["locationKind"] = data.LocationKind.ValueString()
        filterNames = append(filterNames, "location_kind = "+fmt.Sprintf("%q", data.LocationKind.ValueString()))
    }
    if !data.Confidence.IsNull() && !data.Confidence.IsUnknown() {
        filters["confidence"] = data.Confidence.ValueString()
        filterNames = append(filterNames, "confidence = "+fmt.Sprintf("%q", data.Confidence.ValueString()))
    }
    if !data.MatchCount.IsNull() && !data.MatchCount.IsUnknown() {
        filters["matchCount"] = lookupNumber(data.MatchCount)
        filterNames = append(filterNames, "match_count = "+data.MatchCount.ValueBigFloat().String())
    }
    if !data.ValueCount.IsNull() && !data.ValueCount.IsUnknown() {
        filters["valueCount"] = lookupNumber(data.ValueCount)
        filterNames = append(filterNames, "value_count = "+data.ValueCount.ValueBigFloat().String())
    }
    if !data.MatchPercent.IsNull() && !data.MatchPercent.IsUnknown() {
        filters["matchPercent"] = lookupNumber(data.MatchPercent)
        filterNames = append(filterNames, "match_percent = "+data.MatchPercent.ValueBigFloat().String())
    }
    if !data.DistinctValues.IsNull() && !data.DistinctValues.IsUnknown() {
        filters["distinctValues"] = lookupNumber(data.DistinctValues)
        filterNames = append(filterNames, "distinct_values = "+data.DistinctValues.ValueBigFloat().String())
    }
    if !data.TestValueCount.IsNull() && !data.TestValueCount.IsUnknown() {
        filters["testValueCount"] = lookupNumber(data.TestValueCount)
        filterNames = append(filterNames, "test_value_count = "+data.TestValueCount.ValueBigFloat().String())
    }
    if !data.IsFreeText.IsNull() && !data.IsFreeText.IsUnknown() {
        filters["isFreeText"] = data.IsFreeText.ValueBool()
        filterNames = append(filterNames, "is_free_text = "+fmt.Sprintf("%t", data.IsFreeText.ValueBool()))
    }
    if !data.IsNameHinted.IsNull() && !data.IsNameHinted.IsUnknown() {
        filters["isNameHinted"] = data.IsNameHinted.ValueBool()
        filterNames = append(filterNames, "is_name_hinted = "+fmt.Sprintf("%t", data.IsNameHinted.ValueBool()))
    }
    if !data.IsFoundByName.IsNull() && !data.IsFoundByName.IsUnknown() {
        filters["isFoundByName"] = data.IsFoundByName.ValueBool()
        filterNames = append(filterNames, "is_found_by_name = "+fmt.Sprintf("%t", data.IsFoundByName.ValueBool()))
    }
    if !data.IsScrubbed.IsNull() && !data.IsScrubbed.IsUnknown() {
        filters["isScrubbed"] = data.IsScrubbed.ValueBool()
        filterNames = append(filterNames, "is_scrubbed = "+fmt.Sprintf("%t", data.IsScrubbed.ValueBool()))
    }
    if !data.IsSuppressed.IsNull() && !data.IsSuppressed.IsUnknown() {
        filters["isSuppressed"] = data.IsSuppressed.ValueBool()
        filterNames = append(filterNames, "is_suppressed = "+fmt.Sprintf("%t", data.IsSuppressed.ValueBool()))
    }
    if !data.SuppressedReason.IsNull() && !data.SuppressedReason.IsUnknown() {
        filters["suppressedReason"] = data.SuppressedReason.ValueString()
        filterNames = append(filterNames, "suppressed_reason = "+fmt.Sprintf("%q", data.SuppressedReason.ValueString()))
    }
    if !data.IsConfirmed.IsNull() && !data.IsConfirmed.IsUnknown() {
        filters["isConfirmed"] = data.IsConfirmed.ValueBool()
        filterNames = append(filterNames, "is_confirmed = "+fmt.Sprintf("%t", data.IsConfirmed.ValueBool()))
    }
    if !data.IsAiSuggested.IsNull() && !data.IsAiSuggested.IsUnknown() {
        filters["isAiSuggested"] = data.IsAiSuggested.ValueBool()
        filterNames = append(filterNames, "is_ai_suggested = "+fmt.Sprintf("%t", data.IsAiSuggested.ValueBool()))
    }
    if !data.AiReason.IsNull() && !data.AiReason.IsUnknown() {
        filters["aiReason"] = data.AiReason.ValueString()
        filterNames = append(filterNames, "ai_reason = "+fmt.Sprintf("%q", data.AiReason.ValueString()))
    }
    if !data.LastScanId.IsNull() && !data.LastScanId.IsUnknown() {
        filters["lastScanId"] = data.LastScanId.ValueString()
        filterNames = append(filterNames, "last_scan_id = "+fmt.Sprintf("%q", data.LastScanId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data store finding up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data store finding up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataStoreId": true,
        "findingKey": true,
        "dataTypeId": true,
        "dataTypeName": true,
        "dataCategory": true,
        "sensitivity": true,
        "location": true,
        "locationKind": true,
        "confidence": true,
        "matchCount": true,
        "valueCount": true,
        "matchPercent": true,
        "distinctValues": true,
        "testValueCount": true,
        "maskedSamples": true,
        "isFreeText": true,
        "isNameHinted": true,
        "isFoundByName": true,
        "isScrubbed": true,
        "isSuppressed": true,
        "suppressedReason": true,
        "isConfirmed": true,
        "isAiSuggested": true,
        "aiReason": true,
        "lastScanId": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-store-finding/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_store_finding, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store finding found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_store_finding: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-store-finding/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_store_finding, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_store_finding: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store finding matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data store finding matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_store_finding.")
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
    } else if arr, ok := item["dataStoreId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreId = types.StringNull()
        }
    } else if val, ok := item["dataStoreId"].(string); ok {
        data.DataStoreId = types.StringValue(val)
    } else {
        data.DataStoreId = types.StringNull()
    }
    if obj, ok := item["findingKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FindingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FindingKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FindingKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FindingKey = types.StringValue(string(jsonBytes))
        } else {
            data.FindingKey = types.StringNull()
        }
    } else if arr, ok := item["findingKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FindingKey = types.StringValue(string(jsonBytes))
        } else {
            data.FindingKey = types.StringNull()
        }
    } else if val, ok := item["findingKey"].(string); ok {
        data.FindingKey = types.StringValue(val)
    } else {
        data.FindingKey = types.StringNull()
    }
    if obj, ok := item["dataTypeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataTypeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataTypeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataTypeId = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeId = types.StringNull()
        }
    } else if arr, ok := item["dataTypeId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataTypeId = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeId = types.StringNull()
        }
    } else if val, ok := item["dataTypeId"].(string); ok {
        data.DataTypeId = types.StringValue(val)
    } else {
        data.DataTypeId = types.StringNull()
    }
    if obj, ok := item["dataTypeName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataTypeName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataTypeName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataTypeName = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeName = types.StringNull()
        }
    } else if arr, ok := item["dataTypeName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataTypeName = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeName = types.StringNull()
        }
    } else if val, ok := item["dataTypeName"].(string); ok {
        data.DataTypeName = types.StringValue(val)
    } else {
        data.DataTypeName = types.StringNull()
    }
    if obj, ok := item["dataCategory"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategory = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataCategory = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataCategory = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataCategory = types.StringValue(string(jsonBytes))
        } else {
            data.DataCategory = types.StringNull()
        }
    } else if arr, ok := item["dataCategory"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataCategory = types.StringValue(string(jsonBytes))
        } else {
            data.DataCategory = types.StringNull()
        }
    } else if val, ok := item["dataCategory"].(string); ok {
        data.DataCategory = types.StringValue(val)
    } else {
        data.DataCategory = types.StringNull()
    }
    if obj, ok := item["sensitivity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Sensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if arr, ok := item["sensitivity"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if val, ok := item["sensitivity"].(string); ok {
        data.Sensitivity = types.StringValue(val)
    } else {
        data.Sensitivity = types.StringNull()
    }
    if obj, ok := item["location"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Location = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Location = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Location = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Location = types.StringValue(string(jsonBytes))
        } else {
            data.Location = types.StringNull()
        }
    } else if arr, ok := item["location"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Location = types.StringValue(string(jsonBytes))
        } else {
            data.Location = types.StringNull()
        }
    } else if val, ok := item["location"].(string); ok {
        data.Location = types.StringValue(val)
    } else {
        data.Location = types.StringNull()
    }
    if obj, ok := item["locationKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LocationKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LocationKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LocationKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LocationKind = types.StringValue(string(jsonBytes))
        } else {
            data.LocationKind = types.StringNull()
        }
    } else if arr, ok := item["locationKind"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LocationKind = types.StringValue(string(jsonBytes))
        } else {
            data.LocationKind = types.StringNull()
        }
    } else if val, ok := item["locationKind"].(string); ok {
        data.LocationKind = types.StringValue(val)
    } else {
        data.LocationKind = types.StringNull()
    }
    if obj, ok := item["confidence"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Confidence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Confidence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Confidence = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Confidence = types.StringValue(string(jsonBytes))
        } else {
            data.Confidence = types.StringNull()
        }
    } else if arr, ok := item["confidence"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Confidence = types.StringValue(string(jsonBytes))
        } else {
            data.Confidence = types.StringNull()
        }
    } else if val, ok := item["confidence"].(string); ok {
        data.Confidence = types.StringValue(val)
    } else {
        data.Confidence = types.StringNull()
    }
    data.MatchCount = numberFromAPI(data.MatchCount, item["matchCount"])
    data.ValueCount = numberFromAPI(data.ValueCount, item["valueCount"])
    data.MatchPercent = numberFromAPI(data.MatchPercent, item["matchPercent"])
    data.DistinctValues = numberFromAPI(data.DistinctValues, item["distinctValues"])
    data.TestValueCount = numberFromAPI(data.TestValueCount, item["testValueCount"])
    if obj, ok := item["maskedSamples"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MaskedSamples = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MaskedSamples = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MaskedSamples = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MaskedSamples = types.StringValue(string(jsonBytes))
        } else {
            data.MaskedSamples = types.StringNull()
        }
    } else if arr, ok := item["maskedSamples"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MaskedSamples = types.StringValue(string(jsonBytes))
        } else {
            data.MaskedSamples = types.StringNull()
        }
    } else if val, ok := item["maskedSamples"].(string); ok {
        data.MaskedSamples = types.StringValue(val)
    } else {
        data.MaskedSamples = types.StringNull()
    }
    if val, ok := item["isFreeText"].(bool); ok {
        data.IsFreeText = types.BoolValue(val)
    } else {
        data.IsFreeText = types.BoolNull()
    }
    if val, ok := item["isNameHinted"].(bool); ok {
        data.IsNameHinted = types.BoolValue(val)
    } else {
        data.IsNameHinted = types.BoolNull()
    }
    if val, ok := item["isFoundByName"].(bool); ok {
        data.IsFoundByName = types.BoolValue(val)
    } else {
        data.IsFoundByName = types.BoolNull()
    }
    if val, ok := item["isScrubbed"].(bool); ok {
        data.IsScrubbed = types.BoolValue(val)
    } else {
        data.IsScrubbed = types.BoolNull()
    }
    if val, ok := item["isSuppressed"].(bool); ok {
        data.IsSuppressed = types.BoolValue(val)
    } else {
        data.IsSuppressed = types.BoolNull()
    }
    if obj, ok := item["suppressedReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SuppressedReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SuppressedReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SuppressedReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SuppressedReason = types.StringValue(string(jsonBytes))
        } else {
            data.SuppressedReason = types.StringNull()
        }
    } else if arr, ok := item["suppressedReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SuppressedReason = types.StringValue(string(jsonBytes))
        } else {
            data.SuppressedReason = types.StringNull()
        }
    } else if val, ok := item["suppressedReason"].(string); ok {
        data.SuppressedReason = types.StringValue(val)
    } else {
        data.SuppressedReason = types.StringNull()
    }
    if val, ok := item["isConfirmed"].(bool); ok {
        data.IsConfirmed = types.BoolValue(val)
    } else {
        data.IsConfirmed = types.BoolNull()
    }
    if val, ok := item["isAiSuggested"].(bool); ok {
        data.IsAiSuggested = types.BoolValue(val)
    } else {
        data.IsAiSuggested = types.BoolNull()
    }
    if obj, ok := item["aiReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiReason = types.StringValue(string(jsonBytes))
        } else {
            data.AiReason = types.StringNull()
        }
    } else if arr, ok := item["aiReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiReason = types.StringValue(string(jsonBytes))
        } else {
            data.AiReason = types.StringNull()
        }
    } else if val, ok := item["aiReason"].(string); ok {
        data.AiReason = types.StringValue(val)
    } else {
        data.AiReason = types.StringNull()
    }
    if obj, ok := item["lastScanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastScanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastScanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastScanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastScanId = types.StringValue(string(jsonBytes))
        } else {
            data.LastScanId = types.StringNull()
        }
    } else if arr, ok := item["lastScanId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastScanId = types.StringValue(string(jsonBytes))
        } else {
            data.LastScanId = types.StringNull()
        }
    } else if val, ok := item["lastScanId"].(string); ok {
        data.LastScanId = types.StringValue(val)
    } else {
        data.LastScanId = types.StringNull()
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
    } else if arr, ok := item["firstSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["lastSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
