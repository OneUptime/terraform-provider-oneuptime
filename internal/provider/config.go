package provider

import (
    "context"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/diag"
)

// Config holds the provider configuration
type Config struct {
    Host     string
    ApiKey   string
    Client   *Client
}

// NewConfig creates a new configuration from the provider model
func NewConfig(ctx context.Context, model OneuptimeProviderModel) (*Config, diag.Diagnostics) {
    var diags diag.Diagnostics

    config := &Config{}

    // Set host
    if model.Host.IsNull() {
        config.Host = os.Getenv("ONEUPTIME_HOST")
        if config.Host == "" {
            config.Host = "oneuptime.com"
        }
    } else {
        config.Host = model.Host.ValueString()
    }

    // Set API key
    if model.ApiKey.IsNull() {
        config.ApiKey = os.Getenv("ONEUPTIME_API_KEY")
    } else {
        config.ApiKey = model.ApiKey.ValueString()
    }

    // Create client
    client, err := NewClient(config.Host, config.ApiKey)
    if err != nil {
        diags.AddError(
            "Unable to Create API Client",
            "An unexpected error occurred when creating the API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Client Error: "+err.Error(),
        )
        return nil, diags
    }

    config.Client = client
    return config, diags
}
