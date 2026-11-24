package shard27

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard27", []OrderedResource{
        {Order: 22, New: NewWatchlistResource},
        {Order: 77, New: NewTracePipelineProcessorResource},
        {Order: 83, New: NewApiKeyCustomRoleResource},
        {Order: 97, New: NewMonitorSecretResource},
        {Order: 114, New: NewIncidentRoleResource},
        {Order: 129, New: NewAlertCustomFieldResource},
        {Order: 151, New: NewAlertEpisodeOnCallRuleOwningTeamResource},
        {Order: 157, New: NewIncidentAlertResource},
        {Order: 226, New: NewRunbookOwnerRuleResource},
        {Order: 244, New: NewStatusPageFooterLinkResource},
        {Order: 262, New: NewRunbookRuleResource},
        {Order: 275, New: NewMonitorProbeResource},
        {Order: 355, New: NewProxmoxClusterFeedResource},
        {Order: 365, New: NewKubernetesClusterTeamOwnerResource},
        {Order: 420, New: NewVcenterUserOwnerResource},
        {Order: 481, New: NewVcenterUserOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 420, Name: "vcenter_user_owner", LegacyName: "v_center_user_owner", New: NewVcenterUserOwnerResource, NewLegacy: NewVcenterUserOwnerLegacyResource},
    })
}
