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
        {Order: 35, New: NewSecurityCasePriorityDataSource},
        {Order: 70, New: NewDataReconciliationDataSource},
        {Order: 199, New: NewIncidentEpisodeOnCallRuleDataSource},
        {Order: 209, New: NewScheduledMaintenanceMeasurementDataSource},
        {Order: 304, New: NewStatusPageSsoDataSource},
        {Order: 346, New: NewIncomingCallPolicyTeamOwnerDataSource},
        {Order: 378, New: NewSourceMapDataSource},
        {Order: 387, New: NewDashboardSsoDataSource},
        {Order: 477, New: NewHostUserOwnerDataSource},
    }, nil)
}
