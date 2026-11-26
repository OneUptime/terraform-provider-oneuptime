package provider

// generatedShards is how many packages hold the generated resources and data
// sources, internal/provider/shard00 to shard31: each hands what it
// holds to this package (shards.go).
const generatedShards = 32

// The resources written by hand, in their places among the generated ones.
var handwrittenResources = []OrderedResource{
    {Order: 460, New: NewOnCallScheduleLayersResource},
}
