package shard02

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard02", []OrderedResource{
        {Order: 59, New: NewPipelineUserOwnerResource},
        {Order: 96, New: NewUserOverrideResource},
        {Order: 158, New: NewAlertEpisodeLabelRuleResource},
        {Order: 224, New: NewCephClusterOwnerRuleResource},
        {Order: 238, New: NewDomainResource},
        {Order: 262, New: NewRunbookResource},
        {Order: 321, New: NewAiAgentResource},
        {Order: 325, New: NewAiConversationResource},
        {Order: 361, New: NewVcenterFeedResource},
        {Order: 385, New: NewNetworkSiteResource},
        {Order: 427, New: NewHostUserOwnerResource},
        {Order: 475, New: NewVcenterFeedLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 361, Name: "vcenter_feed", LegacyName: "v_center_feed", New: NewVcenterFeedResource, NewLegacy: NewVcenterFeedLegacyResource},
    })
}
