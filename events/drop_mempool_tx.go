// POST /drop_mempool_tx
// =======================

// This payload includes txids removed from the node's mempool.

// Example request body: ./example_events/drop_mempool_tx.json

// We use Pipelining to delete the txids in a single operation from Redis

// =============================================================================