package shard28

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard28", []OrderedResource{
        {Order: 169, New: NewIncidentOwnerRuleOwningTeamResource},
        {Order: 217, New: NewDockerSwarmClusterLabelRuleResource},
        {Order: 282, New: NewScheduledMaintenanceTemplateResource},
        {Order: 293, New: NewOnCallScheduleLayerUserResource},
        {Order: 296, New: NewScheduledMaintenanceExternalCalendarResource},
        {Order: 401, New: NewDockerSwarmClusterResource},
        {Order: 402, New: NewProxmoxClusterTeamOwnerResource},
        {Order: 429, New: NewServerlessFunctionOwnerRuleResource},
    }, nil)
}
