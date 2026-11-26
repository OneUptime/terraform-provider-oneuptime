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
        {Order: 33, New: NewSecurityCaseStatusResource},
        {Order: 38, New: NewDataStoreResource},
        {Order: 43, New: NewDataStoreLabelRuleResource},
        {Order: 99, New: NewIncomingCallLogItemResource},
        {Order: 131, New: NewAlertResource},
        {Order: 143, New: NewAlertEpisodeUserOwnerResource},
        {Order: 153, New: NewAlertLabelRuleResource},
        {Order: 170, New: NewIncidentOnCallRuleResource},
        {Order: 188, New: NewServiceLevelObjectiveResource},
        {Order: 208, New: NewServiceOwnerRuleResource},
        {Order: 287, New: NewScheduledMaintenanceTemplateTeamOwnerResource},
        {Order: 306, New: NewOnCallScheduleLabelRuleResource},
        {Order: 408, New: NewDockerSwarmClusterTeamOwnerResource},
    }, nil)
}
