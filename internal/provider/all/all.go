// Package all links every package of generated resources and data sources
// into the program that imports it: they hand what they hold to
// internal/provider from their init functions (its shards.go). main.go,
// xpprovider and the provider's tests import it.
package all

import (
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard00"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard01"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard02"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard03"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard04"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard05"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard06"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard07"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard08"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard09"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard10"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard11"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard12"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard13"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard14"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard15"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard16"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard17"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard18"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard19"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard20"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard21"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard22"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard23"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard24"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard25"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard26"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard27"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard28"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard29"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard30"
    _ "github.com/oneuptime/terraform-provider-oneuptime/internal/provider/shard31"
)
