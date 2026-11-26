package shard12

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard12", []OrderedDataSource{
        {Order: 15, New: NewMetricSavedViewDataSource},
        {Order: 38, New: NewSecurityCasePriorityDataSource},
        {Order: 75, New: NewDataReconciliationDataSource},
        {Order: 209, New: NewIncidentEpisodeOnCallRuleDataSource},
        {Order: 219, New: NewScheduledMaintenanceMeasurementDataSource},
        {Order: 314, New: NewStatusPageSsoDataSource},
        {Order: 356, New: NewIncomingCallPolicyTeamOwnerDataSource},
        {Order: 388, New: NewSourceMapDataSource},
        {Order: 397, New: NewDashboardSsoDataSource},
        {Order: 487, New: NewHostUserOwnerDataSource},
    }, nil)
}
