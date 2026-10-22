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
var _ datasource.DataSource = &CollectorConfigVersionDataSource{}

func NewCollectorConfigVersionDataSource() datasource.DataSource {
    return &CollectorConfigVersionDataSource{}
}

// CollectorConfigVersionDataSource defines the data source implementation.
type CollectorConfigVersionDataSource struct {
    client *Client
}

// CollectorConfigVersionDataSourceModel describes the data source data model.
type CollectorConfigVersionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    HostCollectorConfigId types.String `tfsdk:"host_collector_config_id"`
    VersionNumber types.Number `tfsdk:"version_number"`
    ConfigYaml types.String `tfsdk:"config_yaml"`
    ConfigHash types.String `tfsdk:"config_hash"`
    ChangeNote types.String `tfsdk:"change_note"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *CollectorConfigVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_collector_config_version"
}

func (d *CollectorConfigVersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Every saved version of a Collector Config, so each host can say which version it runs. Look up an existing collector config version by `id`, or by any of its other arguments (`change_note`, `config_hash`, `config_yaml`, ...): each one set must match, and exactly one collector config version may match them all.",

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
            "host_collector_config_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Collector Config this is a version of. The ID of a `oneuptime_collector_config`.",
                Optional: true,
                Computed: true,
            },
            "version_number": schema.NumberAttribute{
                MarkdownDescription: "The version's number: 1 for the config as created, then one more for every saved change.",
                Optional: true,
                Computed: true,
            },
            "config_yaml": schema.StringAttribute{
                MarkdownDescription: "The collector's config.yaml as this version saved it.",
                Optional: true,
                Computed: true,
            },
            "config_hash": schema.StringAttribute{
                MarkdownDescription: "The SHA-256 of this version's config, which a host reports back once it runs it.",
                Optional: true,
                Computed: true,
            },
            "change_note": schema.StringAttribute{
                MarkdownDescription: "What this version did: created the config, changed it, or put an older version back.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who saved this version. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *CollectorConfigVersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CollectorConfigVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CollectorConfigVersionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.HostCollectorConfigId.IsNull() && !data.HostCollectorConfigId.IsUnknown() {
        filters["hostCollectorConfigId"] = data.HostCollectorConfigId.ValueString()
        filterNames = append(filterNames, "host_collector_config_id = "+fmt.Sprintf("%q", data.HostCollectorConfigId.ValueString()))
    }
    if !data.VersionNumber.IsNull() && !data.VersionNumber.IsUnknown() {
        filters["versionNumber"] = lookupNumber(data.VersionNumber)
        filterNames = append(filterNames, "version_number = "+data.VersionNumber.ValueBigFloat().String())
    }
    if !data.ConfigYaml.IsNull() && !data.ConfigYaml.IsUnknown() {
        filters["configYaml"] = data.ConfigYaml.ValueString()
        filterNames = append(filterNames, "config_yaml = "+fmt.Sprintf("%q", data.ConfigYaml.ValueString()))
    }
    if !data.ConfigHash.IsNull() && !data.ConfigHash.IsUnknown() {
        filters["configHash"] = data.ConfigHash.ValueString()
        filterNames = append(filterNames, "config_hash = "+fmt.Sprintf("%q", data.ConfigHash.ValueString()))
    }
    if !data.ChangeNote.IsNull() && !data.ChangeNote.IsUnknown() {
        filters["changeNote"] = data.ChangeNote.ValueString()
        filterNames = append(filterNames, "change_note = "+fmt.Sprintf("%q", data.ChangeNote.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the collector config version up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the collector config version up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "hostCollectorConfigId": true,
        "versionNumber": true,
        "configYaml": true,
        "configHash": true,
        "changeNote": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/host-collector-config-version/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read collector_config_version, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No collector config version found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read collector_config_version: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/host-collector-config-version/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list collector_config_version, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list collector_config_version: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No collector config version matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one collector config version matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for collector_config_version.")
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
    if obj, ok := item["hostCollectorConfigId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostCollectorConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostCollectorConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostCollectorConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostCollectorConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.HostCollectorConfigId = types.StringNull()
        }
    } else if val, ok := item["hostCollectorConfigId"].(string); ok {
        data.HostCollectorConfigId = types.StringValue(val)
    } else {
        data.HostCollectorConfigId = types.StringNull()
    }
    data.VersionNumber = numberFromAPI(data.VersionNumber, item["versionNumber"])
    if obj, ok := item["configYaml"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConfigYaml = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConfigYaml = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConfigYaml = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConfigYaml = types.StringValue(string(jsonBytes))
        } else {
            data.ConfigYaml = types.StringNull()
        }
    } else if val, ok := item["configYaml"].(string); ok {
        data.ConfigYaml = types.StringValue(val)
    } else {
        data.ConfigYaml = types.StringNull()
    }
    if obj, ok := item["configHash"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConfigHash = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConfigHash = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConfigHash = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConfigHash = types.StringValue(string(jsonBytes))
        } else {
            data.ConfigHash = types.StringNull()
        }
    } else if val, ok := item["configHash"].(string); ok {
        data.ConfigHash = types.StringValue(val)
    } else {
        data.ConfigHash = types.StringNull()
    }
    if obj, ok := item["changeNote"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChangeNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChangeNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChangeNote = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeNote = types.StringNull()
        }
    } else if val, ok := item["changeNote"].(string); ok {
        data.ChangeNote = types.StringValue(val)
    } else {
        data.ChangeNote = types.StringNull()
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
