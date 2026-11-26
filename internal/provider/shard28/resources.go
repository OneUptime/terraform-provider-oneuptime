package shard28

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard28", []OrderedResource{
        {Order: 173, New: NewIncidentOwnerRuleOwningTeamResource},
        {Order: 221, New: NewDockerSwarmClusterLabelRuleResource},
        {Order: 286, New: NewScheduledMaintenanceTemplateResource},
        {Order: 297, New: NewOnCallScheduleLayerUserResource},
        {Order: 300, New: NewScheduledMaintenanceExternalCalendarResource},
        {Order: 405, New: NewDockerSwarmClusterResource},
        {Order: 406, New: NewProxmoxClusterTeamOwnerResource},
        {Order: 433, New: NewServerlessFunctionOwnerRuleResource},
    }, nil)
}
