package shard07

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard07", []OrderedDataSource{
        {Order: 3, New: NewTeamDataSource},
        {Order: 170, New: NewAlertOnCallRuleOwningTeamDataSource},
        {Order: 205, New: NewIncidentMeasurementDataSource},
        {Order: 233, New: NewServiceLabelRuleDataSource},
        {Order: 235, New: NewDockerHostLabelRuleDataSource},
        {Order: 278, New: NewScheduledMaintenanceStateTimelineDataSource},
        {Order: 302, New: NewTeamUserOwnerDataSource},
        {Order: 321, New: NewEmailLogDataSource},
        {Order: 345, New: NewOnCallScheduleOwnerRuleDataSource},
        {Order: 354, New: NewServiceTeamOwnerDataSource},
        {Order: 376, New: NewExceptionDataSource},
        {Order: 413, New: NewMetricTypeDataSource},
        {Order: 433, New: NewPacketCaptureDataSource},
        {Order: 441, New: NewNetworkSiteLinkDataSource},
        {Order: 448, New: NewDockerHostTeamOwnerDataSource},
        {Order: 525, New: NewSecurityEventDataSource},
    }, nil)
}
