package shard04

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
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
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
    SampleRecords types.String `tfsdk:"sample_records"`
    IsFreeText types.Bool `tfsdk:"is_free_text"`
    IsNameHinted types.Bool `tfsdk:"is_name_hinted"`
    IsFoundByName types.Bool `tfsdk:"is_found_by_name"`
    IsFoundByModel types.Bool `tfsdk:"is_found_by_model"`
    IsScrubbed types.Bool `tfsdk:"is_scrubbed"`
    IsSuppressed types.Bool `tfsdk:"is_suppressed"`
    SuppressedReason types.String `tfsdk:"suppressed_reason"`
    IsConfirmed types.Bool `tfsdk:"is_confirmed"`
    IsAiSuggested types.Bool `tfsdk:"is_ai_suggested"`
    AiReason types.String `tfsdk:"ai_reason"`
    AiVerdict types.String `tfsdk:"ai_verdict"`
    AiVerdictAnalysis types.String `tfsdk:"ai_verdict_analysis"`
    AiVerdictAt types.String `tfsdk:"ai_verdict_at"`
    IsConfirmedByAi types.Bool `tfsdk:"is_confirmed_by_ai"`
    DecidedAt types.String `tfsdk:"decided_at"`
    LastScanId types.String `tfsdk:"last_scan_id"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    FixedAt types.String `tfsdk:"fixed_at"`
}

func (d *DataStoreFindingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_store_finding"
}

func (d *DataStoreFindingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What Data Security found in a store: one data type at one location, with its confidence, counts and masked samples. Look up an existing data store finding by `id`, or by any of its other arguments (`ai_reason`, `ai_verdict`, `confidence`, ...): each one set must match, and exactly one data store finding may match them all.",

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
                MarkdownDescription: "ID of the project this finding belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store the finding is in. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_data_store"),
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
            "sample_records": schema.StringAttribute{
                MarkdownDescription: "Telemetry only: up to ten records the last scan found it in, by when each was written and its trace, span or session. Never a value. A JSON value: read it with `jsondecode()`.",
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
            "is_found_by_model": schema.BoolAttribute{
                MarkdownDescription: "Found in free text by the names-and-addresses model that runs inside the scanner, not by a pattern.",
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
            "ai_verdict": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI's verdict on a finding that waits for a person: Confirm, NotSensitive or CannotTell. A person decides.",
                Optional: true,
                Computed: true,
            },
            "ai_verdict_analysis": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI's reason for its verdict, a reason a person could give, and the earlier decisions on findings like it that it cites. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "ai_verdict_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI last gave its verdict on the finding.",
                Computed: true,
            },
            "is_confirmed_by_ai": schema.BoolAttribute{
                MarkdownDescription: "OneUptime AI confirmed the finding without asking, as one of the project's AI rules allows. Undo it from Data Security > AI.",
                Optional: true,
                Computed: true,
            },
            "decided_at": schema.StringAttribute{
                MarkdownDescription: "When the finding was last confirmed or marked not sensitive.",
                Computed: true,
            },
            "last_scan_id": schema.StringAttribute{
                MarkdownDescription: "The scan that last saw the finding. The ID of a `oneuptime_data_security_scan` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_data_security_scan"),
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a scan first found it.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a scan last found it.",
                Computed: true,
            },
            "fixed_at": schema.StringAttribute{
                MarkdownDescription: "For telemetry: when a code fix was verified to have stopped this finding in production. It no longer counts toward its issue unless a scan finds it again.",
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
        filters["matchCount"] = LookupNumber(data.MatchCount)
        filterNames = append(filterNames, "match_count = "+data.MatchCount.ValueBigFloat().String())
    }
    if !data.ValueCount.IsNull() && !data.ValueCount.IsUnknown() {
        filters["valueCount"] = LookupNumber(data.ValueCount)
        filterNames = append(filterNames, "value_count = "+data.ValueCount.ValueBigFloat().String())
    }
    if !data.MatchPercent.IsNull() && !data.MatchPercent.IsUnknown() {
        filters["matchPercent"] = LookupNumber(data.MatchPercent)
        filterNames = append(filterNames, "match_percent = "+data.MatchPercent.ValueBigFloat().String())
    }
    if !data.DistinctValues.IsNull() && !data.DistinctValues.IsUnknown() {
        filters["distinctValues"] = LookupNumber(data.DistinctValues)
        filterNames = append(filterNames, "distinct_values = "+data.DistinctValues.ValueBigFloat().String())
    }
    if !data.TestValueCount.IsNull() && !data.TestValueCount.IsUnknown() {
        filters["testValueCount"] = LookupNumber(data.TestValueCount)
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
    if !data.IsFoundByModel.IsNull() && !data.IsFoundByModel.IsUnknown() {
        filters["isFoundByModel"] = data.IsFoundByModel.ValueBool()
        filterNames = append(filterNames, "is_found_by_model = "+fmt.Sprintf("%t", data.IsFoundByModel.ValueBool()))
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
    if !data.AiVerdict.IsNull() && !data.AiVerdict.IsUnknown() {
        filters["aiVerdict"] = data.AiVerdict.ValueString()
        filterNames = append(filterNames, "ai_verdict = "+fmt.Sprintf("%q", data.AiVerdict.ValueString()))
    }
    if !data.IsConfirmedByAi.IsNull() && !data.IsConfirmedByAi.IsUnknown() {
        filters["isConfirmedByAi"] = data.IsConfirmedByAi.ValueBool()
        filterNames = append(filterNames, "is_confirmed_by_ai = "+fmt.Sprintf("%t", data.IsConfirmedByAi.ValueBool()))
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
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
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
        "sampleRecords": true,
        "isFreeText": true,
        "isNameHinted": true,
        "isFoundByName": true,
        "isFoundByModel": true,
        "isScrubbed": true,
        "isSuppressed": true,
        "suppressedReason": true,
        "isConfirmed": true,
        "isAiSuggested": true,
        "aiReason": true,
        "aiVerdict": true,
        "aiVerdictAnalysis": true,
        "aiVerdictAt": true,
        "isConfirmedByAi": true,
        "decidedAt": true,
        "lastScanId": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "fixedAt": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store finding matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data store finding matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
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
    data.MatchCount = NumberFromAPI(data.MatchCount, item["matchCount"])
    data.ValueCount = NumberFromAPI(data.ValueCount, item["valueCount"])
    data.MatchPercent = NumberFromAPI(data.MatchPercent, item["matchPercent"])
    data.DistinctValues = NumberFromAPI(data.DistinctValues, item["distinctValues"])
    data.TestValueCount = NumberFromAPI(data.TestValueCount, item["testValueCount"])
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
    if obj, ok := item["sampleRecords"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SampleRecords = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SampleRecords = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SampleRecords = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SampleRecords = types.StringValue(string(jsonBytes))
        } else {
            data.SampleRecords = types.StringNull()
        }
    } else if arr, ok := item["sampleRecords"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SampleRecords = types.StringValue(string(jsonBytes))
        } else {
            data.SampleRecords = types.StringNull()
        }
    } else if val, ok := item["sampleRecords"].(string); ok {
        data.SampleRecords = types.StringValue(val)
    } else {
        data.SampleRecords = types.StringNull()
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
    if val, ok := item["isFoundByModel"].(bool); ok {
        data.IsFoundByModel = types.BoolValue(val)
    } else {
        data.IsFoundByModel = types.BoolNull()
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
    if obj, ok := item["aiVerdict"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiVerdict = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiVerdict = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdict = types.StringNull()
        }
    } else if arr, ok := item["aiVerdict"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiVerdict = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdict = types.StringNull()
        }
    } else if val, ok := item["aiVerdict"].(string); ok {
        data.AiVerdict = types.StringValue(val)
    } else {
        data.AiVerdict = types.StringNull()
    }
    if obj, ok := item["aiVerdictAnalysis"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiVerdictAnalysis = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiVerdictAnalysis = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiVerdictAnalysis = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiVerdictAnalysis = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdictAnalysis = types.StringNull()
        }
    } else if arr, ok := item["aiVerdictAnalysis"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiVerdictAnalysis = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdictAnalysis = types.StringNull()
        }
    } else if val, ok := item["aiVerdictAnalysis"].(string); ok {
        data.AiVerdictAnalysis = types.StringValue(val)
    } else {
        data.AiVerdictAnalysis = types.StringNull()
    }
    if obj, ok := item["aiVerdictAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiVerdictAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiVerdictAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiVerdictAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiVerdictAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdictAt = types.StringNull()
        }
    } else if arr, ok := item["aiVerdictAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiVerdictAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiVerdictAt = types.StringNull()
        }
    } else if val, ok := item["aiVerdictAt"].(string); ok {
        data.AiVerdictAt = types.StringValue(val)
    } else {
        data.AiVerdictAt = types.StringNull()
    }
    if val, ok := item["isConfirmedByAi"].(bool); ok {
        data.IsConfirmedByAi = types.BoolValue(val)
    } else {
        data.IsConfirmedByAi = types.BoolNull()
    }
    if obj, ok := item["decidedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DecidedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DecidedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DecidedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DecidedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DecidedAt = types.StringNull()
        }
    } else if arr, ok := item["decidedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DecidedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DecidedAt = types.StringNull()
        }
    } else if val, ok := item["decidedAt"].(string); ok {
        data.DecidedAt = types.StringValue(val)
    } else {
        data.DecidedAt = types.StringNull()
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
    if obj, ok := item["fixedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FixedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FixedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FixedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FixedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FixedAt = types.StringNull()
        }
    } else if arr, ok := item["fixedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FixedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FixedAt = types.StringNull()
        }
    } else if val, ok := item["fixedAt"].(string); ok {
        data.FixedAt = types.StringValue(val)
    } else {
        data.FixedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
