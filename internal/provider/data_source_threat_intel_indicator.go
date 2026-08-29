package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ThreatIntelIndicatorDataSource{}

func NewThreatIntelIndicatorDataSource() datasource.DataSource {
    return &ThreatIntelIndicatorDataSource{}
}

// ThreatIntelIndicatorDataSource defines the data source implementation.
type ThreatIntelIndicatorDataSource struct {
    client *Client
}

// ThreatIntelIndicatorDataSourceModel describes the data source data model.
type ThreatIntelIndicatorDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    FeedId types.String `tfsdk:"feed_id"`
    FeedName types.String `tfsdk:"feed_name"`
    StixId types.String `tfsdk:"stix_id"`
    IndicatorType types.String `tfsdk:"indicator_type"`
    IndicatorValue types.String `tfsdk:"indicator_value"`
    IndicatorName types.String `tfsdk:"indicator_name"`
    Confidence types.Number `tfsdk:"confidence"`
    StixLabels types.Set `tfsdk:"stix_labels"`
    ValidFrom types.String `tfsdk:"valid_from"`
    ValidUntil types.String `tfsdk:"valid_until"`
    Revoked types.Bool `tfsdk:"revoked"`
    Version types.String `tfsdk:"version"`
    RetentionDate types.String `tfsdk:"retention_date"`
}

func (d *ThreatIntelIndicatorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_threat_intel_indicator"
}

func (d *ThreatIntelIndicatorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Threat Intel Indicator Look up an existing threat_intel_indicator by `id` or by `name`.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "feed_id": schema.StringAttribute{
                MarkdownDescription: "Feed ID",
                Computed: true,
            },
            "feed_name": schema.StringAttribute{
                MarkdownDescription: "Feed",
                Computed: true,
            },
            "stix_id": schema.StringAttribute{
                MarkdownDescription: "STIX ID",
                Computed: true,
            },
            "indicator_type": schema.StringAttribute{
                MarkdownDescription: "Indicator Type",
                Computed: true,
            },
            "indicator_value": schema.StringAttribute{
                MarkdownDescription: "Indicator Value",
                Computed: true,
            },
            "indicator_name": schema.StringAttribute{
                MarkdownDescription: "Name",
                Computed: true,
            },
            "confidence": schema.NumberAttribute{
                MarkdownDescription: "Confidence",
                Computed: true,
            },
            "stix_labels": schema.SetAttribute{
                MarkdownDescription: "Labels",
                Computed: true,
                ElementType: types.StringType,
            },
            "valid_from": schema.StringAttribute{
                MarkdownDescription: "Valid From",
                Computed: true,
            },
            "valid_until": schema.StringAttribute{
                MarkdownDescription: "Valid Until",
                Computed: true,
            },
            "revoked": schema.BoolAttribute{
                MarkdownDescription: "Revoked",
                Computed: true,
            },
            "version": schema.StringAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
            "retention_date": schema.StringAttribute{
                MarkdownDescription: "Retention Date",
                Computed: true,
            },
        },
    }
}

func (d *ThreatIntelIndicatorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ThreatIntelIndicatorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ThreatIntelIndicatorDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a threat_intel_indicator.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "projectId": true,
        "feedId": true,
        "feedName": true,
        "stixId": true,
        "indicatorType": true,
        "indicatorValue": true,
        "indicatorName": true,
        "confidence": true,
        "stixLabels": true,
        "validFrom": true,
        "validUntil": true,
        "revoked": true,
        "version": true,
        "retentionDate": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/threat-intel-indicators/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read threat_intel_indicator, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No threat_intel_indicator found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read threat_intel_indicator: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/threat-intel-indicators/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list threat_intel_indicator, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list threat_intel_indicator: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No threat_intel_indicator found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one threat_intel_indicator matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for threat_intel_indicator.")
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
    if obj, ok := item["feedId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FeedId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FeedId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FeedId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FeedId = types.StringValue(string(jsonBytes))
        } else {
            data.FeedId = types.StringNull()
        }
    } else if val, ok := item["feedId"].(string); ok {
        data.FeedId = types.StringValue(val)
    } else {
        data.FeedId = types.StringNull()
    }
    if obj, ok := item["feedName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FeedName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FeedName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FeedName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FeedName = types.StringValue(string(jsonBytes))
        } else {
            data.FeedName = types.StringNull()
        }
    } else if val, ok := item["feedName"].(string); ok {
        data.FeedName = types.StringValue(val)
    } else {
        data.FeedName = types.StringNull()
    }
    if obj, ok := item["stixId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StixId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StixId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StixId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StixId = types.StringValue(string(jsonBytes))
        } else {
            data.StixId = types.StringNull()
        }
    } else if val, ok := item["stixId"].(string); ok {
        data.StixId = types.StringValue(val)
    } else {
        data.StixId = types.StringNull()
    }
    if obj, ok := item["indicatorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IndicatorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if val, ok := item["indicatorType"].(string); ok {
        data.IndicatorType = types.StringValue(val)
    } else {
        data.IndicatorType = types.StringNull()
    }
    if obj, ok := item["indicatorValue"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IndicatorValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IndicatorValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IndicatorValue = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorValue = types.StringNull()
        }
    } else if val, ok := item["indicatorValue"].(string); ok {
        data.IndicatorValue = types.StringValue(val)
    } else {
        data.IndicatorValue = types.StringNull()
    }
    if obj, ok := item["indicatorName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IndicatorName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IndicatorName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IndicatorName = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorName = types.StringNull()
        }
    } else if val, ok := item["indicatorName"].(string); ok {
        data.IndicatorName = types.StringValue(val)
    } else {
        data.IndicatorName = types.StringNull()
    }
    if val, ok := item["confidence"].(float64); ok {
        data.Confidence = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["confidence"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Confidence = types.NumberValue(big.NewFloat(val))
        } else {
            data.Confidence = types.NumberNull()
        }
    } else {
        data.Confidence = types.NumberNull()
    }
    if val, ok := item["stixLabels"].([]interface{}); ok {
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
        data.StixLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.StixLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["validFrom"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ValidFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ValidFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ValidFrom = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ValidFrom = types.StringValue(string(jsonBytes))
        } else {
            data.ValidFrom = types.StringNull()
        }
    } else if val, ok := item["validFrom"].(string); ok {
        data.ValidFrom = types.StringValue(val)
    } else {
        data.ValidFrom = types.StringNull()
    }
    if obj, ok := item["validUntil"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ValidUntil = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ValidUntil = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ValidUntil = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ValidUntil = types.StringValue(string(jsonBytes))
        } else {
            data.ValidUntil = types.StringNull()
        }
    } else if val, ok := item["validUntil"].(string); ok {
        data.ValidUntil = types.StringValue(val)
    } else {
        data.ValidUntil = types.StringNull()
    }
    if val, ok := item["revoked"].(bool); ok {
        data.Revoked = types.BoolValue(val)
    } else {
        data.Revoked = types.BoolNull()
    }
    if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Version = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Version = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Version = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Version = types.StringValue(string(jsonBytes))
        } else {
            data.Version = types.StringNull()
        }
    } else if val, ok := item["version"].(string); ok {
        data.Version = types.StringValue(val)
    } else {
        data.Version = types.StringNull()
    }
    if obj, ok := item["retentionDate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RetentionDate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RetentionDate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RetentionDate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RetentionDate = types.StringValue(string(jsonBytes))
        } else {
            data.RetentionDate = types.StringNull()
        }
    } else if val, ok := item["retentionDate"].(string); ok {
        data.RetentionDate = types.StringValue(val)
    } else {
        data.RetentionDate = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
