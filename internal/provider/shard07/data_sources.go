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
        {Order: 34, New: NewSecurityCaseStepCheckDataSource},
        {Order: 57, New: NewDataSecurityAiProposalDataSource},
        {Order: 180, New: NewAlertOnCallRuleOwningTeamDataSource},
        {Order: 215, New: NewIncidentMeasurementDataSource},
        {Order: 243, New: NewServiceLabelRuleDataSource},
        {Order: 245, New: NewDockerHostLabelRuleDataSource},
        {Order: 288, New: NewScheduledMaintenanceStateTimelineDataSource},
        {Order: 312, New: NewTeamUserOwnerDataSource},
        {Order: 331, New: NewEmailLogDataSource},
        {Order: 355, New: NewOnCallScheduleOwnerRuleDataSource},
        {Order: 364, New: NewServiceTeamOwnerDataSource},
        {Order: 386, New: NewExceptionDataSource},
        {Order: 423, New: NewMetricTypeDataSource},
        {Order: 443, New: NewPacketCaptureDataSource},
        {Order: 451, New: NewNetworkSiteLinkDataSource},
        {Order: 458, New: NewDockerHostTeamOwnerDataSource},
        {Order: 535, New: NewSecurityEventDataSource},
    }, nil)
}
