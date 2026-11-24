package shard29

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard29", []OrderedDataSource{
        {Order: 85, New: NewEntityRiskDataSource},
        {Order: 93, New: NewLlmCostBudgetDataSource},
        {Order: 124, New: NewMonitorCustomFieldDataSource},
        {Order: 182, New: NewIncidentAlertDataSource},
        {Order: 269, New: NewStatusPageSubscriberDataSource},
        {Order: 414, New: NewOnCallTimeLogDataSource},
        {Order: 471, New: NewStorageArrayUserOwnerDataSource},
        {Order: 482, New: NewServerlessFunctionUserOwnerDataSource},
        {Order: 500, New: NewDatabaseLabelRuleDataSource},
    }, nil)
}
