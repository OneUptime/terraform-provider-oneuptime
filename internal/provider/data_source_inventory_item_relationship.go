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
var _ datasource.DataSource = &InventoryItemRelationshipDataSource{}

func NewInventoryItemRelationshipDataSource() datasource.DataSource {
    return &InventoryItemRelationshipDataSource{}
}

// InventoryItemRelationshipDataSource defines the data source implementation.
type InventoryItemRelationshipDataSource struct {
    client *Client
}

// InventoryItemRelationshipDataSourceModel describes the data source data model.
type InventoryItemRelationshipDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    FromEntityKey types.String `tfsdk:"from_entity_key"`
    ToEntityKey types.String `tfsdk:"to_entity_key"`
    RelationshipType types.String `tfsdk:"relationship_type"`
    Source types.String `tfsdk:"source"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    CallCount types.Number `tfsdk:"call_count"`
    ErrorCount types.Number `tfsdk:"error_count"`
    AvgDurationMs types.Number `tfsdk:"avg_duration_ms"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *InventoryItemRelationshipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_inventory_item_relationship"
}

func (d *InventoryItemRelationshipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence. Look up an existing inventory item relationship by `id`, or by any of its other arguments (`avg_duration_ms`, `call_count`, `created_by_user_id`, ...): each one set must match, and exactly one inventory item relationship may match them all.",

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
            "from_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the source entity of this edge.",
                Optional: true,
                Computed: true,
            },
            "to_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the target entity of this edge.",
                Optional: true,
                Computed: true,
            },
            "relationship_type": schema.StringAttribute{
                MarkdownDescription: "The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of).",
                Optional: true,
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When this relationship was first observed in telemetry.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "Most recent time this relationship was observed (bumped, throttled). Drives staleness pruning.",
                Computed: true,
            },
            "call_count": schema.NumberAttribute{
                MarkdownDescription: "Calls observed over this edge in the most recent computation window (depends-on edges only).",
                Optional: true,
                Computed: true,
            },
            "error_count": schema.NumberAttribute{
                MarkdownDescription: "Errored calls observed over this edge in the most recent computation window (depends-on edges only).",
                Optional: true,
                Computed: true,
            },
            "avg_duration_ms": schema.NumberAttribute{
                MarkdownDescription: "Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only).",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *InventoryItemRelationshipDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InventoryItemRelationshipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data InventoryItemRelationshipDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.FromEntityKey.IsNull() && !data.FromEntityKey.IsUnknown() {
        filters["fromEntityKey"] = data.FromEntityKey.ValueString()
        filterNames = append(filterNames, "from_entity_key = "+fmt.Sprintf("%q", data.FromEntityKey.ValueString()))
    }
    if !data.ToEntityKey.IsNull() && !data.ToEntityKey.IsUnknown() {
        filters["toEntityKey"] = data.ToEntityKey.ValueString()
        filterNames = append(filterNames, "to_entity_key = "+fmt.Sprintf("%q", data.ToEntityKey.ValueString()))
    }
    if !data.RelationshipType.IsNull() && !data.RelationshipType.IsUnknown() {
        filters["relationshipType"] = data.RelationshipType.ValueString()
        filterNames = append(filterNames, "relationship_type = "+fmt.Sprintf("%q", data.RelationshipType.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.CallCount.IsNull() && !data.CallCount.IsUnknown() {
        filters["callCount"] = lookupNumber(data.CallCount)
        filterNames = append(filterNames, "call_count = "+data.CallCount.ValueBigFloat().String())
    }
    if !data.ErrorCount.IsNull() && !data.ErrorCount.IsUnknown() {
        filters["errorCount"] = lookupNumber(data.ErrorCount)
        filterNames = append(filterNames, "error_count = "+data.ErrorCount.ValueBigFloat().String())
    }
    if !data.AvgDurationMs.IsNull() && !data.AvgDurationMs.IsUnknown() {
        filters["avgDurationMs"] = lookupNumber(data.AvgDurationMs)
        filterNames = append(filterNames, "avg_duration_ms = "+data.AvgDurationMs.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the inventory item relationship up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the inventory item relationship up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "fromEntityKey": true,
        "toEntityKey": true,
        "relationshipType": true,
        "source": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "callCount": true,
        "errorCount": true,
        "avgDurationMs": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/inventory-item-relationship/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory_item_relationship, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No inventory item relationship found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read inventory_item_relationship: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/inventory-item-relationship/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list inventory_item_relationship, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list inventory_item_relationship: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No inventory item relationship matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one inventory item relationship matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for inventory_item_relationship.")
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
    if obj, ok := item["fromEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FromEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.FromEntityKey = types.StringNull()
        }
    } else if val, ok := item["fromEntityKey"].(string); ok {
        data.FromEntityKey = types.StringValue(val)
    } else {
        data.FromEntityKey = types.StringNull()
    }
    if obj, ok := item["toEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ToEntityKey = types.StringNull()
        }
    } else if val, ok := item["toEntityKey"].(string); ok {
        data.ToEntityKey = types.StringValue(val)
    } else {
        data.ToEntityKey = types.StringNull()
    }
    if obj, ok := item["relationshipType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RelationshipType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RelationshipType = types.StringValue(string(jsonBytes))
        } else {
            data.RelationshipType = types.StringNull()
        }
    } else if val, ok := item["relationshipType"].(string); ok {
        data.RelationshipType = types.StringValue(val)
    } else {
        data.RelationshipType = types.StringNull()
    }
    if obj, ok := item["source"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := item["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
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
    if val, ok := item["callCount"].(float64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["callCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CallCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.CallCount = types.NumberNull()
        }
    } else {
        data.CallCount = types.NumberNull()
    }
    if val, ok := item["errorCount"].(float64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["errorCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ErrorCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorCount = types.NumberNull()
        }
    } else {
        data.ErrorCount = types.NumberNull()
    }
    if val, ok := item["avgDurationMs"].(float64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["avgDurationMs"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.AvgDurationMs = types.NumberNull()
        }
    } else {
        data.AvgDurationMs = types.NumberNull()
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
