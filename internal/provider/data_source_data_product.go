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
var _ datasource.DataSource = &DataProductDataSource{}

func NewDataProductDataSource() datasource.DataSource {
    return &DataProductDataSource{}
}

// DataProductDataSource defines the data source implementation.
type DataProductDataSource struct {
    client *Client
}

// DataProductDataSourceModel describes the data source data model.
type DataProductDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Consumers types.String `tfsdk:"consumers"`
    DataAssets types.Set `tfsdk:"data_assets"`
    FreshByTimeOfDay types.String `tfsdk:"fresh_by_time_of_day"`
    MaxDataAgeInMinutes types.Number `tfsdk:"max_data_age_in_minutes"`
    SlaTargetPercentage types.Number `tfsdk:"sla_target_percentage"`
    OwnerTeamId types.String `tfsdk:"owner_team_id"`
    OwnerUserId types.String `tfsdk:"owner_user_id"`
    Status types.String `tfsdk:"status"`
    LateTableCount types.Number `tfsdk:"late_table_count"`
    LateDataAssetIds types.String `tfsdk:"late_data_asset_ids"`
    DataAsOf types.String `tfsdk:"data_as_of"`
    StatusChangedAt types.String `tfsdk:"status_changed_at"`
    StatusCheckedAt types.String `tfsdk:"status_checked_at"`
    MonitorId types.String `tfsdk:"monitor_id"`
    ServiceLevelObjectiveId types.String `tfsdk:"service_level_objective_id"`
    BadgeToken types.String `tfsdk:"badge_token"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    OwnerUserEmail types.String `tfsdk:"owner_user_email"`
}

func (d *DataProductDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_product"
}

func (d *DataProductDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A named group of watched tables with an owner, its consumers and a freshness SLA, on a status page and tracked as an SLO. Look up an existing data product by `id`, or by any of its other arguments (`name`, `consumers`, `created_by_user_id`, ...): each one set must match, and exactly one data product may match them all.",

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
                MarkdownDescription: "What the data product is called, as its consumers know it.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the data is, and what it is for.",
                Optional: true,
                Computed: true,
            },
            "consumers": schema.StringAttribute{
                MarkdownDescription: "Who uses the data: teams, dashboards, reports, AI agents.",
                Optional: true,
                Computed: true,
            },
            "data_assets": schema.SetAttribute{
                MarkdownDescription: "The watched tables the data product is made of. IDs of `oneuptime_data_asset` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "fresh_by_time_of_day": schema.StringAttribute{
                MarkdownDescription: "The time of day, in UTC (HH:MM), by which each day's data lands. Empty for no daily deadline.",
                Optional: true,
                Computed: true,
            },
            "max_data_age_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How old the data may be, in minutes, before the product is late. Empty for no limit.",
                Optional: true,
                Computed: true,
            },
            "sla_target_percentage": schema.NumberAttribute{
                MarkdownDescription: "The share of time the data must be on time, which its SLO tracks: 99 means late for at most 1% of the time.",
                Optional: true,
                Computed: true,
            },
            "owner_team_id": schema.StringAttribute{
                MarkdownDescription: "ID of the team that owns the data product. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_team"),
            },
            "owner_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who owns the data product. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `owner_user_email`.",
                Optional: true,
                Computed: true,
                Validators: personIDValidators("owner_user_email", false),
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Whether the data is fresh by the SLA: On Time, Late, or Unknown while no table has been seen changing.",
                Optional: true,
                Computed: true,
            },
            "late_table_count": schema.NumberAttribute{
                MarkdownDescription: "How many of the tables are late.",
                Optional: true,
                Computed: true,
            },
            "late_data_asset_ids": schema.StringAttribute{
                MarkdownDescription: "The IDs of the tables that are late. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "data_as_of": schema.StringAttribute{
                MarkdownDescription: "The oldest last change among the tables: what the data is as of.",
                Computed: true,
            },
            "status_changed_at": schema.StringAttribute{
                MarkdownDescription: "When the status last changed.",
                Computed: true,
            },
            "status_checked_at": schema.StringAttribute{
                MarkdownDescription: "When the tables were last checked against the SLA.",
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Metrics monitor that is Degraded, with an incident, while a table is late. The ID of a `oneuptime_monitor`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_monitor"),
            },
            "service_level_objective_id": schema.StringAttribute{
                MarkdownDescription: "ID of the SLO that tracks the data product's SLA, with burn-rate alerts. The ID of a `oneuptime_service_level_objective`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_service_level_objective"),
            },
            "badge_token": schema.StringAttribute{
                MarkdownDescription: "The secret in the badge's URL, so a badge can be embedded where people cannot sign in to OneUptime.",
                Computed: true,
                Sensitive: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "owner_user_email": schema.StringAttribute{
                MarkdownDescription: "Look it up by the email address of the person in `owner_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.",
                Optional: true,
                Computed: true,
                Validators: personEmailValidators(),
            },
        },
    }
}

func (d *DataProductDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataProductDataSourceModel

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
    if !data.Consumers.IsNull() && !data.Consumers.IsUnknown() {
        filters["consumers"] = data.Consumers.ValueString()
        filterNames = append(filterNames, "consumers = "+fmt.Sprintf("%q", data.Consumers.ValueString()))
    }
    if !data.FreshByTimeOfDay.IsNull() && !data.FreshByTimeOfDay.IsUnknown() {
        filters["freshByTimeOfDay"] = data.FreshByTimeOfDay.ValueString()
        filterNames = append(filterNames, "fresh_by_time_of_day = "+fmt.Sprintf("%q", data.FreshByTimeOfDay.ValueString()))
    }
    if !data.MaxDataAgeInMinutes.IsNull() && !data.MaxDataAgeInMinutes.IsUnknown() {
        filters["maxDataAgeInMinutes"] = lookupNumber(data.MaxDataAgeInMinutes)
        filterNames = append(filterNames, "max_data_age_in_minutes = "+data.MaxDataAgeInMinutes.ValueBigFloat().String())
    }
    if !data.SlaTargetPercentage.IsNull() && !data.SlaTargetPercentage.IsUnknown() {
        filters["slaTargetPercentage"] = lookupNumber(data.SlaTargetPercentage)
        filterNames = append(filterNames, "sla_target_percentage = "+data.SlaTargetPercentage.ValueBigFloat().String())
    }
    if !data.OwnerTeamId.IsNull() && !data.OwnerTeamId.IsUnknown() {
        filters["ownerTeamId"] = data.OwnerTeamId.ValueString()
        filterNames = append(filterNames, "owner_team_id = "+fmt.Sprintf("%q", data.OwnerTeamId.ValueString()))
    }
    if !data.OwnerUserId.IsNull() && !data.OwnerUserId.IsUnknown() {
        filters["ownerUserId"] = data.OwnerUserId.ValueString()
        filterNames = append(filterNames, "owner_user_id = "+fmt.Sprintf("%q", data.OwnerUserId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.LateTableCount.IsNull() && !data.LateTableCount.IsUnknown() {
        filters["lateTableCount"] = lookupNumber(data.LateTableCount)
        filterNames = append(filterNames, "late_table_count = "+data.LateTableCount.ValueBigFloat().String())
    }
    if !data.MonitorId.IsNull() && !data.MonitorId.IsUnknown() {
        filters["monitorId"] = data.MonitorId.ValueString()
        filterNames = append(filterNames, "monitor_id = "+fmt.Sprintf("%q", data.MonitorId.ValueString()))
    }
    if !data.ServiceLevelObjectiveId.IsNull() && !data.ServiceLevelObjectiveId.IsUnknown() {
        filters["serviceLevelObjectiveId"] = data.ServiceLevelObjectiveId.ValueString()
        filterNames = append(filterNames, "service_level_objective_id = "+fmt.Sprintf("%q", data.ServiceLevelObjectiveId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.OwnerUserEmail.IsNull() && !data.OwnerUserEmail.IsUnknown() {
        personID, ok := personIDForLookup(ctx, d.client, data.OwnerUserEmail.ValueString(), "owner_user_email", &resp.Diagnostics)
        if !ok {
            return
        }
        filters["ownerUserId"] = personID
        filterNames = append(filterNames, "owner_user_email = "+fmt.Sprintf("%q", data.OwnerUserEmail.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data product up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data product up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "consumers": true,
        "dataAssets": true,
        "freshByTimeOfDay": true,
        "maxDataAgeInMinutes": true,
        "slaTargetPercentage": true,
        "ownerTeamId": true,
        "ownerUserId": true,
        "status": true,
        "lateTableCount": true,
        "lateDataAssetIds": true,
        "dataAsOf": true,
        "statusChangedAt": true,
        "statusCheckedAt": true,
        "monitorId": true,
        "serviceLevelObjectiveId": true,
        "badgeToken": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-product/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_product, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data product found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_product: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-product/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_product, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_product: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data product matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data product matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_product.")
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
    if obj, ok := item["consumers"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Consumers = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Consumers = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Consumers = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Consumers = types.StringValue(string(jsonBytes))
        } else {
            data.Consumers = types.StringNull()
        }
    } else if val, ok := item["consumers"].(string); ok {
        data.Consumers = types.StringValue(val)
    } else {
        data.Consumers = types.StringNull()
    }
    if val, ok := item["dataAssets"].([]interface{}); ok {
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
        data.DataAssets = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DataAssets = types.SetNull(types.StringType)
    }
    if obj, ok := item["freshByTimeOfDay"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FreshByTimeOfDay = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FreshByTimeOfDay = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FreshByTimeOfDay = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FreshByTimeOfDay = types.StringValue(string(jsonBytes))
        } else {
            data.FreshByTimeOfDay = types.StringNull()
        }
    } else if val, ok := item["freshByTimeOfDay"].(string); ok {
        data.FreshByTimeOfDay = types.StringValue(val)
    } else {
        data.FreshByTimeOfDay = types.StringNull()
    }
    data.MaxDataAgeInMinutes = numberFromAPI(data.MaxDataAgeInMinutes, item["maxDataAgeInMinutes"])
    data.SlaTargetPercentage = numberFromAPI(data.SlaTargetPercentage, item["slaTargetPercentage"])
    if obj, ok := item["ownerTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OwnerTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OwnerTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OwnerTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OwnerTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.OwnerTeamId = types.StringNull()
        }
    } else if val, ok := item["ownerTeamId"].(string); ok {
        data.OwnerTeamId = types.StringValue(val)
    } else {
        data.OwnerTeamId = types.StringNull()
    }
    if obj, ok := item["ownerUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OwnerUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OwnerUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OwnerUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OwnerUserId = types.StringValue(string(jsonBytes))
        } else {
            data.OwnerUserId = types.StringNull()
        }
    } else if val, ok := item["ownerUserId"].(string); ok {
        data.OwnerUserId = types.StringValue(val)
    } else {
        data.OwnerUserId = types.StringNull()
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
    data.LateTableCount = numberFromAPI(data.LateTableCount, item["lateTableCount"])
    if obj, ok := item["lateDataAssetIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LateDataAssetIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LateDataAssetIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LateDataAssetIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LateDataAssetIds = types.StringValue(string(jsonBytes))
        } else {
            data.LateDataAssetIds = types.StringNull()
        }
    } else if val, ok := item["lateDataAssetIds"].(string); ok {
        data.LateDataAssetIds = types.StringValue(val)
    } else {
        data.LateDataAssetIds = types.StringNull()
    }
    if obj, ok := item["dataAsOf"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataAsOf = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataAsOf = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataAsOf = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataAsOf = types.StringValue(string(jsonBytes))
        } else {
            data.DataAsOf = types.StringNull()
        }
    } else if val, ok := item["dataAsOf"].(string); ok {
        data.DataAsOf = types.StringValue(val)
    } else {
        data.DataAsOf = types.StringNull()
    }
    if obj, ok := item["statusChangedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusChangedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StatusChangedAt = types.StringNull()
        }
    } else if val, ok := item["statusChangedAt"].(string); ok {
        data.StatusChangedAt = types.StringValue(val)
    } else {
        data.StatusChangedAt = types.StringNull()
    }
    if obj, ok := item["statusCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusCheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusCheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StatusCheckedAt = types.StringNull()
        }
    } else if val, ok := item["statusCheckedAt"].(string); ok {
        data.StatusCheckedAt = types.StringValue(val)
    } else {
        data.StatusCheckedAt = types.StringNull()
    }
    if obj, ok := item["monitorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorId = types.StringNull()
        }
    } else if val, ok := item["monitorId"].(string); ok {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if obj, ok := item["serviceLevelObjectiveId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceLevelObjectiveId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceLevelObjectiveId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceLevelObjectiveId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceLevelObjectiveId = types.StringNull()
        }
    } else if val, ok := item["serviceLevelObjectiveId"].(string); ok {
        data.ServiceLevelObjectiveId = types.StringValue(val)
    } else {
        data.ServiceLevelObjectiveId = types.StringNull()
    }
    if obj, ok := item["badgeToken"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BadgeToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BadgeToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BadgeToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BadgeToken = types.StringValue(string(jsonBytes))
        } else {
            data.BadgeToken = types.StringNull()
        }
    } else if val, ok := item["badgeToken"].(string); ok {
        data.BadgeToken = types.StringValue(val)
    } else {
        data.BadgeToken = types.StringNull()
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
    data.OwnerUserEmail = personEmailAfterLookup(ctx, d.client, data.OwnerUserId, data.OwnerUserEmail)

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
