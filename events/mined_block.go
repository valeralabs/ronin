// POST /mined_block
// ===================

// This payload includes data related to block mined by this Stacks node. This
// will never be invoked if the node is configured only as a follower. This is 
// invoked when the miner assembles the block; this block may or may not win 
// the sortition.

// This endpoint will only broadcast events to observers that explicitly 
// register for `MinedBlocks` events, `AnyEvent` observers will not receive the 
// events by default.

// Example request body: ./example_events/mined_block.json

// We can ignore this event as the API doesn't need to store this information.

// =============================================================================