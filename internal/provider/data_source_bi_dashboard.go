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
var _ datasource.DataSource = &BiDashboardDataSource{}

func NewBiDashboardDataSource() datasource.DataSource {
    return &BiDashboardDataSource{}
}

// BiDashboardDataSource defines the data source implementation.
type BiDashboardDataSource struct {
    client *Client
}

// BiDashboardDataSourceModel describes the data source data model.
type BiDashboardDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataBiConnectionId types.String `tfsdk:"data_bi_connection_id"`
    Tool types.String `tfsdk:"tool"`
    ExternalId types.String `tfsdk:"external_id"`
    Name types.String `tfsdk:"name"`
    Url types.String `tfsdk:"url"`
    Folder types.String `tfsdk:"folder"`
    Owner types.String `tfsdk:"owner"`
    ChangedAt types.String `tfsdk:"changed_at"`
    TableCount types.Number `tfsdk:"table_count"`
    DataAssets types.Set `tfsdk:"data_assets"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *BiDashboardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_bi_dashboard"
}

func (d *BiDashboardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A dashboard of a connected BI tool, with the watched tables its charts query. Look up an existing bi dashboard by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `data_bi_connection_id`, ...): each one set must match, and exactly one bi dashboard may match them all.",

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
            "data_bi_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the BI tool the dashboard was read from. The ID of a `oneuptime_bi_tool`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_bi_tool"),
            },
            "tool": schema.StringAttribute{
                MarkdownDescription: "Which BI tool the dashboard is in.",
                Optional: true,
                Computed: true,
            },
            "external_id": schema.StringAttribute{
                MarkdownDescription: "The dashboard's own ID in its tool.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The dashboard's name, as its tool shows it.",
                Optional: true,
                Computed: true,
            },
            "url": schema.StringAttribute{
                MarkdownDescription: "Where the dashboard opens in its tool.",
                Optional: true,
                Computed: true,
            },
            "folder": schema.StringAttribute{
                MarkdownDescription: "The collection, project, folder or workspace it is in.",
                Optional: true,
                Computed: true,
            },
            "owner": schema.StringAttribute{
                MarkdownDescription: "Who owns it in its tool.",
                Optional: true,
                Computed: true,
            },
            "changed_at": schema.StringAttribute{
                MarkdownDescription: "When the dashboard last changed in its tool.",
                Computed: true,
            },
            "table_count": schema.NumberAttribute{
                MarkdownDescription: "How many watched tables its charts query.",
                Optional: true,
                Computed: true,
            },
            "data_assets": schema.SetAttribute{
                MarkdownDescription: "The watched tables the dashboard's charts query. IDs of `oneuptime_data_asset` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *BiDashboardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BiDashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data BiDashboardDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DataBiConnectionId.IsNull() && !data.DataBiConnectionId.IsUnknown() {
        filters["dataBiConnectionId"] = data.DataBiConnectionId.ValueString()
        filterNames = append(filterNames, "data_bi_connection_id = "+fmt.Sprintf("%q", data.DataBiConnectionId.ValueString()))
    }
    if !data.Tool.IsNull() && !data.Tool.IsUnknown() {
        filters["tool"] = data.Tool.ValueString()
        filterNames = append(filterNames, "tool = "+fmt.Sprintf("%q", data.Tool.ValueString()))
    }
    if !data.ExternalId.IsNull() && !data.ExternalId.IsUnknown() {
        filters["externalId"] = data.ExternalId.ValueString()
        filterNames = append(filterNames, "external_id = "+fmt.Sprintf("%q", data.ExternalId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Url.IsNull() && !data.Url.IsUnknown() {
        filters["url"] = data.Url.ValueString()
        filterNames = append(filterNames, "url = "+fmt.Sprintf("%q", data.Url.ValueString()))
    }
    if !data.Folder.IsNull() && !data.Folder.IsUnknown() {
        filters["folder"] = data.Folder.ValueString()
        filterNames = append(filterNames, "folder = "+fmt.Sprintf("%q", data.Folder.ValueString()))
    }
    if !data.Owner.IsNull() && !data.Owner.IsUnknown() {
        filters["owner"] = data.Owner.ValueString()
        filterNames = append(filterNames, "owner = "+fmt.Sprintf("%q", data.Owner.ValueString()))
    }
    if !data.TableCount.IsNull() && !data.TableCount.IsUnknown() {
        filters["tableCount"] = lookupNumber(data.TableCount)
        filterNames = append(filterNames, "table_count = "+data.TableCount.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the bi dashboard up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the bi dashboard up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataBiConnectionId": true,
        "tool": true,
        "externalId": true,
        "name": true,
        "url": true,
        "folder": true,
        "owner": true,
        "changedAt": true,
        "tableCount": true,
        "dataAssets": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-dashboard/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read bi_dashboard, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No bi dashboard found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read bi_dashboard: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-dashboard/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list bi_dashboard, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list bi_dashboard: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No bi dashboard matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one bi dashboard matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for bi_dashboard.")
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
    if obj, ok := item["dataBiConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataBiConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataBiConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataBiConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataBiConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.DataBiConnectionId = types.StringNull()
        }
    } else if val, ok := item["dataBiConnectionId"].(string); ok {
        data.DataBiConnectionId = types.StringValue(val)
    } else {
        data.DataBiConnectionId = types.StringNull()
    }
    if obj, ok := item["tool"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Tool = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Tool = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Tool = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Tool = types.StringValue(string(jsonBytes))
        } else {
            data.Tool = types.StringNull()
        }
    } else if val, ok := item["tool"].(string); ok {
        data.Tool = types.StringValue(val)
    } else {
        data.Tool = types.StringNull()
    }
    if obj, ok := item["externalId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExternalId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExternalId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExternalId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExternalId = types.StringValue(string(jsonBytes))
        } else {
            data.ExternalId = types.StringNull()
        }
    } else if val, ok := item["externalId"].(string); ok {
        data.ExternalId = types.StringValue(val)
    } else {
        data.ExternalId = types.StringNull()
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
    } else if val, ok := item["url"].(string); ok {
        data.Url = types.StringValue(val)
    } else {
        data.Url = types.StringNull()
    }
    if obj, ok := item["folder"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Folder = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Folder = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Folder = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Folder = types.StringValue(string(jsonBytes))
        } else {
            data.Folder = types.StringNull()
        }
    } else if val, ok := item["folder"].(string); ok {
        data.Folder = types.StringValue(val)
    } else {
        data.Folder = types.StringNull()
    }
    if obj, ok := item["owner"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Owner = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Owner = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Owner = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Owner = types.StringValue(string(jsonBytes))
        } else {
            data.Owner = types.StringNull()
        }
    } else if val, ok := item["owner"].(string); ok {
        data.Owner = types.StringValue(val)
    } else {
        data.Owner = types.StringNull()
    }
    if obj, ok := item["changedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChangedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ChangedAt = types.StringNull()
        }
    } else if val, ok := item["changedAt"].(string); ok {
        data.ChangedAt = types.StringValue(val)
    } else {
        data.ChangedAt = types.StringNull()
    }
    data.TableCount = numberFromAPI(data.TableCount, item["tableCount"])
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
