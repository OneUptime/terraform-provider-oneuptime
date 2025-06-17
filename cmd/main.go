package main

import (
    "context"
    "flag"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/OneUptime/terraform-provider-oneuptime"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name oneuptime

var (
    // these will be set by the goreleaser configuration
    // to appropriate values for the compiled binary.
    version string = "dev"

    // goreleaser can pass other information to the main package, such as the specific commit
    // https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
    var debug bool

    flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
    flag.Parse()

    opts := providerserver.ServeOpts{
        Address: "registry.terraform.io/oneuptime/oneuptime",
        Debug:   debug,
    }

    err := providerserver.Serve(context.Background(), oneuptime.NewProvider(version), opts)
    if err != nil {
        log.Fatal(err.Error())
    }
}
