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
var _ datasource.DataSource = &HuntRunDataSource{}

func NewHuntRunDataSource() datasource.DataSource {
    return &HuntRunDataSource{}
}

// HuntRunDataSource defines the data source implementation.
type HuntRunDataSource struct {
    client *Client
}

// HuntRunDataSourceModel describes the data source data model.
type HuntRunDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    SecurityHuntId types.String `tfsdk:"security_hunt_id"`
    WindowStartAt types.String `tfsdk:"window_start_at"`
    WindowEndAt types.String `tfsdk:"window_end_at"`
    EventCount types.Number `tfsdk:"event_count"`
    Error types.String `tfsdk:"error"`
}

func (d *HuntRunDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_hunt_run"
}

func (d *HuntRunDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Each scheduled run of a hunt: the events it searched, how many it found, or why it could not run. Kept for 90 days. Look up an existing hunt run by `id`, or by any of its other arguments (`error`, `event_count`, `security_hunt_id`): each one set must match, and exactly one hunt run may match them all.",

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
                MarkdownDescription: "ID of the project this run's hunt belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_hunt_id": schema.StringAttribute{
                MarkdownDescription: "ID of the hunt that ran. The ID of a `oneuptime_hunt`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_hunt"),
            },
            "window_start_at": schema.StringAttribute{
                MarkdownDescription: "The start of the events the run searched: the moment after the run before it stopped.",
                Computed: true,
            },
            "window_end_at": schema.StringAttribute{
                MarkdownDescription: "The end of the events the run searched, two minutes before it ran.",
                Computed: true,
            },
            "event_count": schema.NumberAttribute{
                MarkdownDescription: "How many events the run's search found in its window. Empty when the run could not be made.",
                Optional: true,
                Computed: true,
            },
            "error": schema.StringAttribute{
                MarkdownDescription: "Why the run could not be made. Empty when it ran; its window is searched again by the next run.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *HuntRunDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HuntRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HuntRunDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.SecurityHuntId.IsNull() && !data.SecurityHuntId.IsUnknown() {
        filters["securityHuntId"] = data.SecurityHuntId.ValueString()
        filterNames = append(filterNames, "security_hunt_id = "+fmt.Sprintf("%q", data.SecurityHuntId.ValueString()))
    }
    if !data.EventCount.IsNull() && !data.EventCount.IsUnknown() {
        filters["eventCount"] = lookupNumber(data.EventCount)
        filterNames = append(filterNames, "event_count = "+data.EventCount.ValueBigFloat().String())
    }
    if !data.Error.IsNull() && !data.Error.IsUnknown() {
        filters["error"] = data.Error.ValueString()
        filterNames = append(filterNames, "error = "+fmt.Sprintf("%q", data.Error.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the hunt run up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the hunt run up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "securityHuntId": true,
        "windowStartAt": true,
        "windowEndAt": true,
        "eventCount": true,
        "error": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-hunt-run/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hunt_run, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No hunt run found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read hunt_run: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-hunt-run/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list hunt_run, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list hunt_run: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No hunt run matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one hunt run matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for hunt_run.")
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
    if obj, ok := item["securityHuntId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityHuntId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityHuntId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityHuntId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityHuntId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityHuntId = types.StringNull()
        }
    } else if val, ok := item["securityHuntId"].(string); ok {
        data.SecurityHuntId = types.StringValue(val)
    } else {
        data.SecurityHuntId = types.StringNull()
    }
    if obj, ok := item["windowStartAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WindowStartAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WindowStartAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WindowStartAt = types.StringValue(string(jsonBytes))
        } else {
            data.WindowStartAt = types.StringNull()
        }
    } else if val, ok := item["windowStartAt"].(string); ok {
        data.WindowStartAt = types.StringValue(val)
    } else {
        data.WindowStartAt = types.StringNull()
    }
    if obj, ok := item["windowEndAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WindowEndAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WindowEndAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WindowEndAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WindowEndAt = types.StringValue(string(jsonBytes))
        } else {
            data.WindowEndAt = types.StringNull()
        }
    } else if val, ok := item["windowEndAt"].(string); ok {
        data.WindowEndAt = types.StringValue(val)
    } else {
        data.WindowEndAt = types.StringNull()
    }
    data.EventCount = numberFromAPI(data.EventCount, item["eventCount"])
    if obj, ok := item["error"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Error = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Error = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Error = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Error = types.StringValue(string(jsonBytes))
        } else {
            data.Error = types.StringNull()
        }
    } else if val, ok := item["error"].(string); ok {
        data.Error = types.StringValue(val)
    } else {
        data.Error = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
