package shard14

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard14", []OrderedResource{
        {Order: 57, New: NewTableTeamOwnerResource},
        {Order: 84, New: NewTraceRecordingRuleResource},
        {Order: 89, New: NewStatusPageResource},
        {Order: 97, New: NewIncomingCallPolicyResource},
        {Order: 165, New: NewIncidentEpisodeTeamOwnerResource},
        {Order: 186, New: NewScheduledMaintenanceMeasurementResource},
        {Order: 193, New: NewSloMonitorRuleResource},
        {Order: 231, New: NewRunbookLabelRuleResource},
        {Order: 253, New: NewScheduledMaintenanceEventResource},
        {Order: 337, New: NewDashboardResource},
        {Order: 351, New: NewOnCallUserGroupResource},
        {Order: 409, New: NewProxmoxClusterUserOwnerResource},
        {Order: 458, New: NewRumApplicationOwnerRuleResource},
    }, nil)
}
