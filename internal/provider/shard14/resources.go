package shard14

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard14", []OrderedResource{
        {Order: 56, New: NewTableTeamOwnerResource},
        {Order: 80, New: NewTraceRecordingRuleResource},
        {Order: 85, New: NewStatusPageResource},
        {Order: 93, New: NewIncomingCallPolicyResource},
        {Order: 161, New: NewIncidentEpisodeTeamOwnerResource},
        {Order: 182, New: NewScheduledMaintenanceMeasurementResource},
        {Order: 189, New: NewSloMonitorRuleResource},
        {Order: 227, New: NewRunbookLabelRuleResource},
        {Order: 249, New: NewScheduledMaintenanceEventResource},
        {Order: 333, New: NewDashboardResource},
        {Order: 347, New: NewOnCallUserGroupResource},
        {Order: 405, New: NewProxmoxClusterUserOwnerResource},
        {Order: 454, New: NewRumApplicationOwnerRuleResource},
    }, nil)
}
