// POST /new_block
// =================

// This payload includes data related to a newly processed block,
// and any events emitted from Stacks transactions during the block.

// If the transaction originally comes from the parent microblock stream 
// preceding this block, the microblock related fields will be filled in.

// Example request body: ./example_events/new_block.json

// We use this to update the included transactions in the block and to store 
// the block data itself.

// =============================================================================