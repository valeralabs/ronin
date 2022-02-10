# Performance

RONIN is designed to have very little real performance limits. It's built in
Go which is compiled to native machine code (no JIT compilation like JS/TS)
and uses an in-memory datastore (Redis) for the majority of its data requests.

For example, using 10KB blobs, my M1 Max MacBook gets the following results for
benchmarking Redis:

```other
SET: 125628.14 requests per second, p50=0.055 msec
GET: 123152.71 requests per second, p50=0.055 msec
```

Though it's worth noting that the M1 Max is an SoC, where the CPU and RAM sit on
the same chip. The RAM is LDDR5 @ 6500MHz, and has a massive 400GB/s bandwidth
to the CPU, so Redis will run very fast on this machine.

For reference, Stripe's entire global network processes around 13,000 requests
per second. So, there's not going to be much of a bottleneck from Redis.
Theoretically, a single Ronin instance could serve the entire Stacks ecosystem.

If anything, the biggest bottleneck will be your networking. Ronin can probably
bottleneck Redis with enough cores, but assuming a 5KB total request at Redis'
total throughput, that's 5Gbps of throughput.

To sum up, performance isn't much of an issue here. >90% of your traffic's
latency will just be due to network latency over the cables. It's more important
to have many, small and cheap nodes on the edge running Ronin than having a few,
64-core monsters if you're looking for a production deployment for cloud
infrastructure.

As for everyone on a personal level, you should just run it on a $30 Raspberry
Pi on your home network. Speed simply won't be a thought, and it's way better
for decentralisation too.
