# Ronin

[![Go](https://github.com/syvita/ronin/actions/workflows/go.yml/badge.svg)](https://github.com/syvita/ronin/actions/workflows/go.yml)

Hello there!

Ronin is a ultra-speed Stacks API server. It's super lightweight, but scales easily.

## Why are we making this?

Because we don't like slow APIs. Also because it's super hard to set up one using current codebases.

![Live footage of us trying to deploy a Stacks API](./repo-img/mendit.gif)

<small align="center"><i>Live footage of us trying to deploy a Stacks API before Ronin</i></small>

## What's the plan?

Simple stuff: Ronin will use the event dispatcher from a Stacks node to hydrate it with that sweet, sweet data.

![Smug face](./repo-img/smug.gif)

- Built in Go
- Uses Redis as primary datastore (which will fill up all of that juicy RAM of yours)
- Uses the Rocket web framework for *speeeeed*
- Sometimes you won't be able to run it entirely in memory because some people don't have enough RAM to store an entire blockchain in it (what are you doing???).

Jokes aside, Ronin eats RAM up. By the very nature of its design, it stores the entire Stacks chain data (plus all the other parsed data that an API needs) in-memory. This means that the API can access the data way, way faster than using something like PostgreSQL.

Instead of taking 1-15ms to find a transaction from its TXID, then loading it from disk and passing it back to the API server, an in-memory store works a lot quicker. On my M1 Max MacBook, I get around 175746.92 `GET`s per *second* for a value size of 5KB (more than the parsed response for 95%+ of Stacks transactions), so a latency of 5.69 *microseconds*. For reference, if you're sitting a foot away from your laptop, the time it takes the sound waves to travel to your ears is around 880 microseconds, or 154x slower. So, it's *really* fast.

The disadvantages to this are here too - everything has to be stored in RAM in the first place. RAM is limited and expensive, so most people won't be able to run Ronin optimally, entirely in memory (the API data is probably around 30-50GB). Though, many people have SSDs or even NVME drives, which are pretty fast too.

After your physical memory fills up, Ronin runs exactly the same, but your operating system will start pushing data from the datastore into "virtual memory" which is actually on your disk in a swap file. When the API receives a request for a tx, it asks the datastore for it, who then asks the OS' memory for it, who then has to swap it from the (relatively) slow disk to RAM. It's less efficient, but it's still way faster than using PostgreSQL or something similiar. 

This is one of the ways Ronin can run on anything from a Raspberry Pi to a very powerful cloud server in an efficient way. Obviously, the cloud server will be orders of magnitude faster than a Raspberry Pi, but at that point, the biggest latency will originate from the fiber cables to your home. If you want the fastest experience, run Ronin on a decently high performance server at home.

<small align="center"><i>But please, at least run it on an SSD or NVME... HDDs aren't good for your health I beg you</i></small>

It scales, it descales (is that even a word...) and yeah it's like really fast.

It uses one of those lovely cores on your computer to run the event handler (it handles events from the Stacks node... duh) and uses the rest for the primary API. This means it will run on a Raspberry Pi with an SD card OR you can run it on a 60 core and 128GB RAM monster.

Anyways I'm still working on this so it doesn't work fully yet or at all (lol...).
