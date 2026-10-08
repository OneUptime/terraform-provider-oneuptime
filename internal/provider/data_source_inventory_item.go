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
var _ datasource.DataSource = &InventoryItemDataSource{}

func NewInventoryItemDataSource() datasource.DataSource {
    return &InventoryItemDataSource{}
}

// InventoryItemDataSource defines the data source implementation.
type InventoryItemDataSource struct {
    client *Client
}

// InventoryItemDataSourceModel describes the data source data model.
type InventoryItemDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    EntityType types.String `tfsdk:"entity_type"`
    EntityKey types.String `tfsdk:"entity_key"`
    DisplayName types.String `tfsdk:"display_name"`
    Source types.String `tfsdk:"source"`
    Description types.String `tfsdk:"description"`
    IdentifyingAttributes types.String `tfsdk:"identifying_attributes"`
    DescriptiveAttributes types.String `tfsdk:"descriptive_attributes"`
    Labels types.String `tfsdk:"labels"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    InventoryStatus types.String `tfsdk:"inventory_status"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    CustomFields types.String `tfsdk:"custom_fields"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *InventoryItemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_inventory_item"
}

func (d *InventoryItemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Catalog of everything OneUptime knows about your estate (service, host, k8s.pod, container, network device, ...), discovered from telemetry resource attributes, mirrored from inventory tables, or registered by hand. Look up an existing inventory item by `id`, or by any of its other arguments (`archived_by_user_id`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one inventory item may match them all.",

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
            "entity_type": schema.StringAttribute{
                MarkdownDescription: "The OpenTelemetry entity type (service, host, k8s.pod, container, ...).",
                Optional: true,
                Computed: true,
            },
            "entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns).",
                Optional: true,
                Computed: true,
            },
            "display_name": schema.StringAttribute{
                MarkdownDescription: "Human-readable name shown in the Inventory list.",
                Optional: true,
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "How this row came to exist: discovered from telemetry, mirrored from a OneUptime inventory table, or created manually by a user. Determines whether stale-entity pruning applies.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Free-text description. Primarily for manually created entities, where there are no telemetry attributes to explain what the thing is.",
                Optional: true,
                Computed: true,
            },
            "identifying_attributes": schema.StringAttribute{
                MarkdownDescription: "The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "descriptive_attributes": schema.StringAttribute{
                MarkdownDescription: "Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "labels": schema.StringAttribute{
                MarkdownDescription: "Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster).",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "Polymorphic pointer id to the rich typed row named by resourceType, if any.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When this entity was first observed in telemetry.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "Most recent time this entity was observed in telemetry (bumped, throttled). Drives staleness pruning.",
                Computed: true,
            },
            "inventory_status": schema.StringAttribute{
                MarkdownDescription: "Current heartbeat status: live, recent, stale, never seen, or not tracked.",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this item archived? Archived items are hidden from the default list but keep their identity and keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When this item was archived.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom fields on this item. A JSON value: write it with `jsonencode()`.",
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

func (d *InventoryItemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InventoryItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data InventoryItemDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.EntityType.IsNull() && !data.EntityType.IsUnknown() {
        filters["entityType"] = data.EntityType.ValueString()
        filterNames = append(filterNames, "entity_type = "+fmt.Sprintf("%q", data.EntityType.ValueString()))
    }
    if !data.EntityKey.IsNull() && !data.EntityKey.IsUnknown() {
        filters["entityKey"] = data.EntityKey.ValueString()
        filterNames = append(filterNames, "entity_key = "+fmt.Sprintf("%q", data.EntityKey.ValueString()))
    }
    if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
        filters["displayName"] = data.DisplayName.ValueString()
        filterNames = append(filterNames, "display_name = "+fmt.Sprintf("%q", data.DisplayName.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
        filters["resourceType"] = data.ResourceType.ValueString()
        filterNames = append(filterNames, "resource_type = "+fmt.Sprintf("%q", data.ResourceType.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.InventoryStatus.IsNull() && !data.InventoryStatus.IsUnknown() {
        filters["inventoryStatus"] = data.InventoryStatus.ValueString()
        filterNames = append(filterNames, "inventory_status = "+fmt.Sprintf("%q", data.InventoryStatus.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the inventory item up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the inventory item up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "entityType": true,
        "entityKey": true,
        "displayName": true,
        "source": true,
        "description": true,
        "identifyingAttributes": true,
        "descriptiveAttributes": true,
        "labels": true,
        "resourceType": true,
        "resourceId": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "inventoryStatus": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "customFields": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/inventory-item/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory_item, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No inventory item found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read inventory_item: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/inventory-item/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list inventory_item, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list inventory_item: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No inventory item matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one inventory item matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for inventory_item.")
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
    if obj, ok := item["entityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EntityType = types.StringValue(string(jsonBytes))
        } else {
            data.EntityType = types.StringNull()
        }
    } else if val, ok := item["entityType"].(string); ok {
        data.EntityType = types.StringValue(val)
    } else {
        data.EntityType = types.StringNull()
    }
    if obj, ok := item["entityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.EntityKey = types.StringNull()
        }
    } else if val, ok := item["entityKey"].(string); ok {
        data.EntityKey = types.StringValue(val)
    } else {
        data.EntityKey = types.StringNull()
    }
    if obj, ok := item["displayName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisplayName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DisplayName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DisplayName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DisplayName = types.StringValue(string(jsonBytes))
        } else {
            data.DisplayName = types.StringNull()
        }
    } else if val, ok := item["displayName"].(string); ok {
        data.DisplayName = types.StringValue(val)
    } else {
        data.DisplayName = types.StringNull()
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
    if obj, ok := item["identifyingAttributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IdentifyingAttributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IdentifyingAttributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IdentifyingAttributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IdentifyingAttributes = types.StringValue(string(jsonBytes))
        } else {
            data.IdentifyingAttributes = types.StringNull()
        }
    } else if val, ok := item["identifyingAttributes"].(string); ok {
        data.IdentifyingAttributes = types.StringValue(val)
    } else {
        data.IdentifyingAttributes = types.StringNull()
    }
    if obj, ok := item["descriptiveAttributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptiveAttributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DescriptiveAttributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DescriptiveAttributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DescriptiveAttributes = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptiveAttributes = types.StringNull()
        }
    } else if val, ok := item["descriptiveAttributes"].(string); ok {
        data.DescriptiveAttributes = types.StringValue(val)
    } else {
        data.DescriptiveAttributes = types.StringNull()
    }
    if obj, ok := item["labels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Labels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Labels = types.StringValue(string(jsonBytes))
        } else {
            data.Labels = types.StringNull()
        }
    } else if val, ok := item["labels"].(string); ok {
        data.Labels = types.StringValue(val)
    } else {
        data.Labels = types.StringNull()
    }
    if obj, ok := item["resourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceType = types.StringNull()
        }
    } else if val, ok := item["resourceType"].(string); ok {
        data.ResourceType = types.StringValue(val)
    } else {
        data.ResourceType = types.StringNull()
    }
    if obj, ok := item["resourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if val, ok := item["resourceId"].(string); ok {
        data.ResourceId = types.StringValue(val)
    } else {
        data.ResourceId = types.StringNull()
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
    if obj, ok := item["inventoryStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InventoryStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InventoryStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InventoryStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InventoryStatus = types.StringValue(string(jsonBytes))
        } else {
            data.InventoryStatus = types.StringNull()
        }
    } else if val, ok := item["inventoryStatus"].(string); ok {
        data.InventoryStatus = types.StringValue(val)
    } else {
        data.InventoryStatus = types.StringNull()
    }
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := item["customFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := item["customFields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
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
