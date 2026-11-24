package shard02

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard02", []OrderedResource{
        {Order: 58, New: NewPipelineUserOwnerResource},
        {Order: 92, New: NewUserOverrideResource},
        {Order: 154, New: NewAlertEpisodeLabelRuleResource},
        {Order: 220, New: NewCephClusterOwnerRuleResource},
        {Order: 234, New: NewDomainResource},
        {Order: 258, New: NewRunbookResource},
        {Order: 317, New: NewAiAgentResource},
        {Order: 321, New: NewAiConversationResource},
        {Order: 357, New: NewVcenterFeedResource},
        {Order: 381, New: NewNetworkSiteResource},
        {Order: 423, New: NewHostUserOwnerResource},
        {Order: 471, New: NewVcenterFeedLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 357, Name: "vcenter_feed", LegacyName: "v_center_feed", New: NewVcenterFeedResource, NewLegacy: NewVcenterFeedLegacyResource},
    })
}
