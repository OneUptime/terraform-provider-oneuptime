package shard00

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard00", []OrderedResource{
        {Order: 56, New: NewTableUserOwnerResource},
        {Order: 68, New: NewAutomaticResponsePolicyResource},
        {Order: 125, New: NewFormResource},
        {Order: 148, New: NewAlertOnCallRuleResource},
        {Order: 164, New: NewIncidentEpisodeUserOwnerResource},
        {Order: 202, New: NewMonitorLabelRuleResource},
        {Order: 226, New: NewStorageArrayOwnerRuleResource},
        {Order: 280, New: NewMonitorTeamOwnerResource},
        {Order: 319, New: NewProbeTeamOwnerResource},
        {Order: 333, New: NewInventoryItemResource},
        {Order: 350, New: NewWorkspaceNotificationSummaryResource},
        {Order: 356, New: NewCephClusterFeedResource},
        {Order: 379, New: NewNetworkDeviceRoleResource},
        {Order: 380, New: NewNetworkAlertPolicyResource},
        {Order: 381, New: NewSnmpCredentialProfileResource},
        {Order: 392, New: NewNetworkDeviceLinkResource},
        {Order: 425, New: NewHostResource},
        {Order: 471, New: NewProbeTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 319, Name: "probe_team_owner", LegacyName: "probe_owner_team", New: NewProbeTeamOwnerResource, NewLegacy: NewProbeTeamOwnerLegacyResource},
    })
}
