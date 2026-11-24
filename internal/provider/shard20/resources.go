package shard20

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard20", []OrderedResource{
        {Order: 8, New: NewCustomRolePermissionResource},
        {Order: 9, New: NewTeamCustomRoleResource},
        {Order: 18, New: NewLogPipelineProcessorResource},
        {Order: 32, New: NewSecurityCaseStatusResource},
        {Order: 37, New: NewDataStoreResource},
        {Order: 42, New: NewDataStoreLabelRuleResource},
        {Order: 95, New: NewIncomingCallLogItemResource},
        {Order: 127, New: NewAlertResource},
        {Order: 139, New: NewAlertEpisodeUserOwnerResource},
        {Order: 149, New: NewAlertLabelRuleResource},
        {Order: 166, New: NewIncidentOnCallRuleResource},
        {Order: 184, New: NewServiceLevelObjectiveResource},
        {Order: 204, New: NewServiceOwnerRuleResource},
        {Order: 283, New: NewScheduledMaintenanceTemplateTeamOwnerResource},
        {Order: 302, New: NewOnCallScheduleLabelRuleResource},
        {Order: 404, New: NewDockerSwarmClusterTeamOwnerResource},
    }, nil)
}
