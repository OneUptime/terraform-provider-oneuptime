package provider

import (
    "context"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/diag"
)

// Config holds the provider configuration
type Config struct {
    OneuptimeUrl string
    ApiKey       string
    Client       *Client
}

// NewConfig creates a new configuration from the provider model
func NewConfig(ctx context.Context, model OneuptimeProviderModel) (*Config, diag.Diagnostics) {
    var diags diag.Diagnostics

    config := &Config{}

    // Set oneuptime_url
    if model.OneuptimeUrl.IsNull() {
        config.OneuptimeUrl = os.Getenv("ONEUPTIME_URL")
        if config.OneuptimeUrl == "" {
            config.OneuptimeUrl = "oneuptime.com"
        }
    } else {
        config.OneuptimeUrl = model.OneuptimeUrl.ValueString()
    }

    // Set API key
    if model.ApiKey.IsNull() {
        config.ApiKey = os.Getenv("ONEUPTIME_API_KEY")
        if config.ApiKey == "" {
            diags.AddError(
                "Missing API Key",
                "API key is required for authentication. "+
                    "Please provide it via the api_key attribute or the ONEUPTIME_API_KEY environment variable.",
            )
            return nil, diags
        }
    } else {
        config.ApiKey = model.ApiKey.ValueString()
    }

    // Create client
    client, err := NewClient(config.OneuptimeUrl, config.ApiKey)
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
