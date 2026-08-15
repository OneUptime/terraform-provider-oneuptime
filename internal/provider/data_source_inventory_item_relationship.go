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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *InventoryItemRelationshipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_inventory_item_relationship"
}

func (d *InventoryItemRelationshipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence. Look up an existing inventory_item_relationship by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "from_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the source entity of this edge..",
                Computed: true,
            },
            "to_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the target entity of this edge..",
                Computed: true,
            },
            "relationship_type": schema.StringAttribute{
                MarkdownDescription: "The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of)..",
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies..",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "call_count": schema.NumberAttribute{
                MarkdownDescription: "Calls observed over this edge in the most recent computation window (depends-on edges only)..",
                Computed: true,
            },
            "error_count": schema.NumberAttribute{
                MarkdownDescription: "Errored calls observed over this edge in the most recent computation window (depends-on edges only)..",
                Computed: true,
            },
            "avg_duration_ms": schema.NumberAttribute{
                MarkdownDescription: "Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only)..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a inventory_item_relationship.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
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
        "deletedByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No inventory_item_relationship found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.Post(ctx, "/inventory-item-relationship/get-list", listBody)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No inventory_item_relationship found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one inventory_item_relationship matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
