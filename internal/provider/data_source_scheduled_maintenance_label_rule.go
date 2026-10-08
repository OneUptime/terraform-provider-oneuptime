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
var _ datasource.DataSource = &ScheduledMaintenanceLabelRuleDataSource{}

func NewScheduledMaintenanceLabelRuleDataSource() datasource.DataSource {
    return &ScheduledMaintenanceLabelRuleDataSource{}
}

// ScheduledMaintenanceLabelRuleDataSource defines the data source implementation.
type ScheduledMaintenanceLabelRuleDataSource struct {
    client *Client
}

// ScheduledMaintenanceLabelRuleDataSourceModel describes the data source data model.
type ScheduledMaintenanceLabelRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Criteria types.String `tfsdk:"criteria"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Monitors types.Set `tfsdk:"monitors"`
    ScheduledMaintenanceLabels types.Set `tfsdk:"scheduled_maintenance_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    TitlePattern types.String `tfsdk:"title_pattern"`
    DescriptionPattern types.String `tfsdk:"description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    LabelsToAdd types.Set `tfsdk:"labels_to_add"`
    InheritLabelsFromMonitors types.Bool `tfsdk:"inherit_labels_from_monitors"`
    InheritLabelsFromHosts types.Bool `tfsdk:"inherit_labels_from_hosts"`
    InheritLabelsFromKubernetesClusters types.Bool `tfsdk:"inherit_labels_from_kubernetes_clusters"`
    InheritLabelsFromDockerHosts types.Bool `tfsdk:"inherit_labels_from_docker_hosts"`
    InheritLabelsFromPodmanHosts types.Bool `tfsdk:"inherit_labels_from_podman_hosts"`
    InheritLabelsFromServices types.Bool `tfsdk:"inherit_labels_from_services"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ScheduledMaintenanceLabelRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_label_rule"
}

func (d *ScheduledMaintenanceLabelRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules for automatically attaching labels to scheduled maintenance events — including labels inherited from the event's monitors — when matching events are created Look up an existing scheduled maintenance label rule by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one scheduled maintenance label rule may match them all.",

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
            "criteria": schema.StringAttribute{
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this scheduled maintenance label rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this scheduled maintenance label rule.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Optional: true,
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for events on these monitors. Leave empty to match events on any monitor. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "scheduled_maintenance_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for events that already have at least one of these labels. Leave empty to match regardless of event labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for events on monitors that have at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the event title. Leave empty to match any title.",
                Optional: true,
                Computed: true,
            },
            "description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the event description. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against any of the event's monitor names. Leave empty to match any monitor.",
                Optional: true,
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against any of the event's monitor descriptions. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "labels_to_add": schema.SetAttribute{
                MarkdownDescription: "Labels to attach to the event when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "inherit_labels_from_monitors": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also copy every label of the event's monitors onto the event.",
                Optional: true,
                Computed: true,
            },
            "inherit_labels_from_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also copy every label of the event's affected hosts onto the event.",
                Optional: true,
                Computed: true,
            },
            "inherit_labels_from_kubernetes_clusters": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also copy every label of the event's affected Kubernetes clusters onto the event.",
                Optional: true,
                Computed: true,
            },
            "inherit_labels_from_docker_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also copy every label of the event's affected Docker hosts onto the event.",
                Optional: true,
                Computed: true,
            },
            "inherit_labels_from_podman_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also copy every label of the event's affected Podman hosts onto the event.",
                Optional: true,
                Computed: true,
            },
            "inherit_labels_from_services": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also copy every label of the event's affected services onto the event.",
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

func (d *ScheduledMaintenanceLabelRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceLabelRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceLabelRuleDataSourceModel

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
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.TitlePattern.IsNull() && !data.TitlePattern.IsUnknown() {
        filters["titlePattern"] = data.TitlePattern.ValueString()
        filterNames = append(filterNames, "title_pattern = "+fmt.Sprintf("%q", data.TitlePattern.ValueString()))
    }
    if !data.DescriptionPattern.IsNull() && !data.DescriptionPattern.IsUnknown() {
        filters["descriptionPattern"] = data.DescriptionPattern.ValueString()
        filterNames = append(filterNames, "description_pattern = "+fmt.Sprintf("%q", data.DescriptionPattern.ValueString()))
    }
    if !data.MonitorNamePattern.IsNull() && !data.MonitorNamePattern.IsUnknown() {
        filters["monitorNamePattern"] = data.MonitorNamePattern.ValueString()
        filterNames = append(filterNames, "monitor_name_pattern = "+fmt.Sprintf("%q", data.MonitorNamePattern.ValueString()))
    }
    if !data.MonitorDescriptionPattern.IsNull() && !data.MonitorDescriptionPattern.IsUnknown() {
        filters["monitorDescriptionPattern"] = data.MonitorDescriptionPattern.ValueString()
        filterNames = append(filterNames, "monitor_description_pattern = "+fmt.Sprintf("%q", data.MonitorDescriptionPattern.ValueString()))
    }
    if !data.InheritLabelsFromMonitors.IsNull() && !data.InheritLabelsFromMonitors.IsUnknown() {
        filters["inheritLabelsFromMonitors"] = data.InheritLabelsFromMonitors.ValueBool()
        filterNames = append(filterNames, "inherit_labels_from_monitors = "+fmt.Sprintf("%t", data.InheritLabelsFromMonitors.ValueBool()))
    }
    if !data.InheritLabelsFromHosts.IsNull() && !data.InheritLabelsFromHosts.IsUnknown() {
        filters["inheritLabelsFromHosts"] = data.InheritLabelsFromHosts.ValueBool()
        filterNames = append(filterNames, "inherit_labels_from_hosts = "+fmt.Sprintf("%t", data.InheritLabelsFromHosts.ValueBool()))
    }
    if !data.InheritLabelsFromKubernetesClusters.IsNull() && !data.InheritLabelsFromKubernetesClusters.IsUnknown() {
        filters["inheritLabelsFromKubernetesClusters"] = data.InheritLabelsFromKubernetesClusters.ValueBool()
        filterNames = append(filterNames, "inherit_labels_from_kubernetes_clusters = "+fmt.Sprintf("%t", data.InheritLabelsFromKubernetesClusters.ValueBool()))
    }
    if !data.InheritLabelsFromDockerHosts.IsNull() && !data.InheritLabelsFromDockerHosts.IsUnknown() {
        filters["inheritLabelsFromDockerHosts"] = data.InheritLabelsFromDockerHosts.ValueBool()
        filterNames = append(filterNames, "inherit_labels_from_docker_hosts = "+fmt.Sprintf("%t", data.InheritLabelsFromDockerHosts.ValueBool()))
    }
    if !data.InheritLabelsFromPodmanHosts.IsNull() && !data.InheritLabelsFromPodmanHosts.IsUnknown() {
        filters["inheritLabelsFromPodmanHosts"] = data.InheritLabelsFromPodmanHosts.ValueBool()
        filterNames = append(filterNames, "inherit_labels_from_podman_hosts = "+fmt.Sprintf("%t", data.InheritLabelsFromPodmanHosts.ValueBool()))
    }
    if !data.InheritLabelsFromServices.IsNull() && !data.InheritLabelsFromServices.IsUnknown() {
        filters["inheritLabelsFromServices"] = data.InheritLabelsFromServices.ValueBool()
        filterNames = append(filterNames, "inherit_labels_from_services = "+fmt.Sprintf("%t", data.InheritLabelsFromServices.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the scheduled maintenance label rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the scheduled maintenance label rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "monitors": true,
        "scheduledMaintenanceLabels": true,
        "monitorLabels": true,
        "titlePattern": true,
        "descriptionPattern": true,
        "monitorNamePattern": true,
        "monitorDescriptionPattern": true,
        "labelsToAdd": true,
        "inheritLabelsFromMonitors": true,
        "inheritLabelsFromHosts": true,
        "inheritLabelsFromKubernetesClusters": true,
        "inheritLabelsFromDockerHosts": true,
        "inheritLabelsFromPodmanHosts": true,
        "inheritLabelsFromServices": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/scheduled-maintenance-label-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_label_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance label rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_label_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/scheduled-maintenance-label-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scheduled_maintenance_label_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list scheduled_maintenance_label_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance label rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one scheduled maintenance label rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for scheduled_maintenance_label_rule.")
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
    if obj, ok := item["criteria"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Criteria = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Criteria = types.StringValue(string(jsonBytes))
        } else {
            data.Criteria = types.StringNull()
        }
    } else if val, ok := item["criteria"].(string); ok {
        data.Criteria = types.StringValue(val)
    } else {
        data.Criteria = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["monitors"].([]interface{}); ok {
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
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if val, ok := item["scheduledMaintenanceLabels"].([]interface{}); ok {
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
        data.ScheduledMaintenanceLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.ScheduledMaintenanceLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["monitorLabels"].([]interface{}); ok {
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
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MonitorLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["titlePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.TitlePattern = types.StringNull()
        }
    } else if val, ok := item["titlePattern"].(string); ok {
        data.TitlePattern = types.StringValue(val)
    } else {
        data.TitlePattern = types.StringNull()
    }
    if obj, ok := item["descriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["descriptionPattern"].(string); ok {
        data.DescriptionPattern = types.StringValue(val)
    } else {
        data.DescriptionPattern = types.StringNull()
    }
    if obj, ok := item["monitorNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorNamePattern = types.StringNull()
        }
    } else if val, ok := item["monitorNamePattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    } else {
        data.MonitorNamePattern = types.StringNull()
    }
    if obj, ok := item["monitorDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["monitorDescriptionPattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    } else {
        data.MonitorDescriptionPattern = types.StringNull()
    }
    if val, ok := item["labelsToAdd"].([]interface{}); ok {
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
        data.LabelsToAdd = types.SetValueMust(types.StringType, setItems)
    } else {
        data.LabelsToAdd = types.SetNull(types.StringType)
    }
    if val, ok := item["inheritLabelsFromMonitors"].(bool); ok {
        data.InheritLabelsFromMonitors = types.BoolValue(val)
    } else {
        data.InheritLabelsFromMonitors = types.BoolNull()
    }
    if val, ok := item["inheritLabelsFromHosts"].(bool); ok {
        data.InheritLabelsFromHosts = types.BoolValue(val)
    } else {
        data.InheritLabelsFromHosts = types.BoolNull()
    }
    if val, ok := item["inheritLabelsFromKubernetesClusters"].(bool); ok {
        data.InheritLabelsFromKubernetesClusters = types.BoolValue(val)
    } else {
        data.InheritLabelsFromKubernetesClusters = types.BoolNull()
    }
    if val, ok := item["inheritLabelsFromDockerHosts"].(bool); ok {
        data.InheritLabelsFromDockerHosts = types.BoolValue(val)
    } else {
        data.InheritLabelsFromDockerHosts = types.BoolNull()
    }
    if val, ok := item["inheritLabelsFromPodmanHosts"].(bool); ok {
        data.InheritLabelsFromPodmanHosts = types.BoolValue(val)
    } else {
        data.InheritLabelsFromPodmanHosts = types.BoolNull()
    }
    if val, ok := item["inheritLabelsFromServices"].(bool); ok {
        data.InheritLabelsFromServices = types.BoolValue(val)
    } else {
        data.InheritLabelsFromServices = types.BoolNull()
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
