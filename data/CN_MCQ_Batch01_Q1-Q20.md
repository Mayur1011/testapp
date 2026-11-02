# CN MCQ/MSQ Bank — Batch 1 (Q1–Q20)

**What this batch covers (M1 + start of M2):** why networks are layered, OSI vs TCP/IP, PDU names, encapsulation overhead (MSS vs MTU), demultiplexing at each layer, connection-oriented vs connectionless service, circuit vs packet vs message vs virtual-circuit switching delay, optimal packet size, store-and-forward vs cut-through, mesh topology link count, client-server vs P2P distribution time, unicast/broadcast/multicast/anycast, bandwidth vs latency vs throughput vs goodput, bottleneck throughput, transmission delay, propagation delay, and the transmission = propagation crossover.

> In all numericals: 1 Mbps = 10^6 bps and 1 Gbps = 10^9 bps (link rates are always decimal); data sizes are stated in each question; IPv4 header = 20 B and TCP header = 20 B (no options); Ethernet header = 14 B + FCS 4 B; preamble+SFD (8 B) and inter-frame gap (12 B) are counted only when the question says "on the wire".

---

## Q1. [M1 | Layering rationale | Easy | MCQ | OA + Interview]

A laptop switches from Wi-Fi to wired Ethernet. The browser, the HTTP library and the TCP/IP stack keep working with no code changes. Which property of layered design explains this?

- A. Lower layers repeat upper-layer functions, so one can take over if another fails
- B. Layering keeps total header overhead as small as possible
- C. Every layer is built in hardware, so swapping hardware is invisible
- D. Each layer gives a service to the layer above through a fixed interface, so a layer's implementation can be swapped without touching the layers above

**Answer:** D

**Why:** IP only needs one thing from the link layer: "deliver this packet to that next-hop address on this link." Whether Wi-Fi or Ethernet does it is hidden behind the interface. Only L1 and L2 changed. L3 through L7 are untouched.

**Why the other options fail:**
- A: Layers aren't backups for each other. Some functions do repeat, such as error checks at L2 and L4, but that is for end-to-end correctness, not failover.
- B: Layering *adds* overhead, because every layer puts on its own header. That is the price you pay for modularity.
- C: TCP, IP and HTTP are software. Layering is about interfaces, not hardware.

**Concept recap:** Layer N uses the service of layer N−1 and offers a service to layer N+1.
- **Protocol** = rules between *peer* layers on different hosts (horizontal).
- **Interface** = the service boundary between *adjacent* layers on the same host (vertical).
- Benefits: modularity, independent evolution, easier debugging. Cost: header overhead.

**Trap alert:** Mixing up protocol and interface. The Wi-Fi→Ethernet swap works because the *interface* stays fixed, even though the L2 *protocol* changed.

**Interview angle:** "Why are network stacks layered?" Each layer solves one problem: bits on a medium → delivery on one link → routing across networks → process-to-process delivery → application meaning. Each exposes a stable service upward, so Wi-Fi↔Ethernet or HTTP/1.1↔HTTP/2 can change independently. Follow-up: *"Where is layering violated in practice?"* NAT rewrites L4 ports; the TCP checksum covers IP addresses through a pseudo-header; TCP MSS is derived from the link MTU.

---

## Q2. [M1 | OSI vs TCP/IP | Easy | MCQ | OA + Interview]

Which OSI layers are absorbed into the **application layer** of the TCP/IP model?

- A. Transport and Network
- B. Data Link and Physical
- C. Session and Presentation
- D. Network and Session

**Answer:** C

**Why:** TCP/IP has no separate session layer (dialog control) and no presentation layer (encoding, compression, encryption). Applications do those jobs themselves: cookies for sessions, JSON/UTF-8/gzip for presentation. So OSI L5+L6+L7 together map to TCP/IP's application layer.

```
OSI                         TCP/IP          Examples
7 Application  \
6 Presentation  >  ------>  Application     HTTP, DNS, gzip, JSON
5 Session      /
4 Transport        ------>  Transport       TCP, UDP
3 Network          ------>  Internet        IP, ICMP
2 Data Link    \
1 Physical     /   ------>  Link            Ethernet, 802.11
```

**Why the other options fail:**
- A: Transport and Network are separate layers in both models.
- B: L1+L2 are merged into TCP/IP's *link* layer, not its application layer. This option catches students who remember "something gets merged" but not what.
- D: Network is a distinct layer in both models.

**Concept recap:** OSI = a 7-layer *reference model*, designed before the protocols existed. TCP/IP = a 4–5 layer *protocol suite*, described after the protocols were built. Industry still uses OSI numbers as shorthand ("L4 LB", "L7 proxy").

**Trap alert:** "TLS is the presentation layer" isn't standard in either model. Say "TLS sits between TCP and the application."

**Interview angle:** "Compare OSI and TCP/IP." Follow-up: *"What's the difference between an L4 and an L7 load balancer?"* L4 routes on IP + port and never reads the payload. L7 parses HTTP and routes on path, headers or cookies.

---

## Q3. [M1 | PDUs per layer | Medium | MSQ | OA]

Select all that apply. Which layer → PDU pairings are correct?

- A. Transport (TCP) → segment
- B. Network → frame
- C. Data link → frame
- D. Transport (UDP) → datagram
- E. Physical → packet

**Answer:** A, C, D

**Why:** A **PDU** (protocol data unit) is a layer's header plus the payload it carries. The PDU of layer N becomes the payload of layer N−1:

```
App data
+ TCP header  -> segment          (L4)   [UDP: datagram]
+ IP header   -> packet/datagram  (L3)
+ Eth hdr+FCS -> frame            (L2)
-> bits on the medium             (L1)
```

**Why the other options fail:**
- B: A frame is L2. L3's unit is a packet.
- E: L1 carries bits/symbols. A "packet" exists only at L3.

**Concept recap:** App → message · TCP → segment · UDP → datagram · IP → packet · Link → frame · Physical → bits.

**Trap alert:** "Datagram" is used for both UDP and IP. Check which layer the question means.

**Interview angle:** "Walk me through encapsulation for an HTTP request." Follow-up: *"Which headers change hop by hop?"* L2 is rebuilt at every router. At L3 only TTL and checksum change (addresses too if NAT is involved). L4 is end-to-end unless NAT/PAT rewrites ports.

---

## Q4. [M1 | Encapsulation & overhead | Medium | MCQ | OA + Interview]

A TCP sender puts **1460 B** of application data in one segment (TCP header 20 B, IPv4 header 20 B, Ethernet header 14 B + FCS 4 B; ignore preamble and inter-frame gap). What are (i) the IPv4 **Total Length** field and (ii) the Ethernet frame size?

- A. (i) 1480 B, (ii) 1500 B
- B. (i) 1500 B, (ii) 1518 B
- C. (i) 1460 B, (ii) 1518 B
- D. (i) 1500 B, (ii) 1514 B

**Answer:** B

**Why:** IPv4 Total Length counts the *whole IP packet*: IP header + TCP header + data. The Ethernet frame then adds a 14 B header in front and a 4 B FCS at the end. 1460 is the standard **MSS** (TCP payload only) for a 1500 B **MTU** (full IP packet).

**Working:**
```
Payload (MSS)       1460 B
+ TCP header          20 B  -> segment  = 1480 B
+ IPv4 header         20 B  -> packet   = 1500 B  <- Total Length (= MTU)
+ Eth header + FCS    18 B  -> frame    = 1518 B
(on the wire: + 8 B preamble/SFD + 12 B IFG = 1538 B)
Frame-level payload efficiency = 1460/1518 = 96.2 %
```

**Why the other options fail:**
- A: Leaves the IP header itself out of Total Length.
- C: Treats Total Length as the application payload, confusing MSS with IP length.
- D: Leaves out the FCS. Wireshark shows 1514 only because NICs usually strip the FCS before capture.

**Concept recap:**
- MTU = largest IP packet a link can carry (1500 on Ethernet).
- MSS = MTU − IP header − TCP header = 1460.
- Ethernet frame = 64 B min, 1518 B max (untagged).

**Trap alert:** **MSS is payload only; MTU is the full IP packet.** Don't "correct" 1518 to 1514 because of a capture.

**Interview angle:** "MTU vs MSS?" MSS is advertised in the SYN so the peer never sends segments that need IP fragmentation. Follow-up: *"What happens inside a VPN that adds ~60 B of headers?"* The effective MTU drops. You need MSS clamping or PMTUD, otherwise you get an MTU black hole (P89).

---

## Q5. [M1 | Demultiplexing per layer | Medium | MCQ | Interview]

A host has one NIC and runs a browser (two HTTPS connections to the *same* server) plus an SSH session. A frame arrives. Which fields, layer by layer, deliver its payload to the correct socket?

- A. Ethernet **EtherType** → IPv4 **Protocol** → TCP **4-tuple** (src IP, src port, dst IP, dst port)
- B. Destination MAC → IPv4 TTL → TCP sequence number
- C. Ethernet EtherType → IPv4 ToS → TCP destination port only
- D. 802.1Q VLAN tag → IPv4 Identification → TCP window

**Answer:** A

**Why:** Each layer has a "who is this for?" field. For an **established** TCP connection, the stack matches the full 4-tuple. Both HTTPS connections come from server port 443, so only the client's two different ephemeral ports tell them apart.

**Working:**
```
| Layer    | Field         | Value                                   | Decision           |
|----------|---------------|-----------------------------------------|--------------------|
| Ethernet | EtherType     | 0x0800                                  | -> IPv4            |
| IPv4     | Protocol      | 6                                       | -> TCP             |
| TCP      | 4-tuple       | 203.0.113.10:443 -> 192.168.1.5:51234   | -> HTTPS socket #1 |
|          |               | (socket #2 uses local port 51235)       |                    |
```

**Why the other options fail:**
- B: The destination MAC only decides whether the NIC accepts the frame. TTL is a hop limit. The sequence number orders bytes *within* a connection; it doesn't identify the connection.
- C: ToS/DSCP is a QoS mark. A destination port alone can't separate two connections to the same service.
- D: The VLAN tag picks a broadcast domain; Identification is for reassembling fragments; window is flow control.

**Concept recap:** L2 EtherType (0x0800 IPv4, 0x86DD IPv6, 0x0806 ARP) → L3 Protocol (6 TCP, 17 UDP, 1 ICMP) → L4:
- UDP: (dst IP, dst port)
- TCP listening socket: (dst IP, dst port)
- TCP established socket: the full 4-tuple

**Trap alert:** "TCP demuxes by destination port" is true only for *listening* sockets.

**Interview angle:** "How does one server port 443 serve 10,000 clients?" Every connection is a distinct 4-tuple, and `accept()` returns a new socket for each one. Follow-up: *"Then what limits connections from one client to one server:port?"* The client's ephemeral port range (~28k by default on Linux), which leads to port exhaustion.

---

## Q6. [M1 | Connection-oriented vs connectionless | Medium | MSQ | OA + Interview]

Select all that apply.

- A. A connectionless service can never deliver packets in order.
- B. TCP keeps connection state only at the two endpoints; ordinary IP routers on the path hold no per-TCP-connection state.
- C. Virtual-circuit networks (ATM, Frame Relay) keep per-connection state inside intermediate switches.
- D. "Connection-oriented" means the network has reserved bandwidth for the connection.
- E. UDP datagrams between the same two sockets may arrive out of order, duplicated, or not at all.

**Answer:** B, C, E

**Why:** TCP's "connection" is an *end-to-end* agreement: sequence numbers, windows and buffers live only at the hosts. The IP core stays connectionless. A virtual circuit puts state *inside* the network. UDP adds nothing on top of IP's best effort.

```
Service            State at hosts   State in network       Ordered?
UDP over IP        minimal          none                   no
TCP over IP        yes              none                   yes (hosts reorder)
Virtual circuit    yes              VC tables              yes (single path)
Circuit switching  yes              reserved capacity      yes
```

**Why the other options fail:**
- A: Connectionless packets often *do* arrive in order. Order just isn't guaranteed.
- D: TCP reserves no bandwidth. Reservation belongs to circuit switching or explicit QoS.

**Concept recap:** Connection-oriented = a setup phase creates shared state, then data, then teardown. Reliability, ordering and reservation are *separate* properties.

**Trap alert:** "TCP guarantees delivery" is false. TCP guarantees an *ordered, reliable byte stream, or an error* (reset or timeout). It cannot deliver across a cut cable.

**Interview angle:** "When would you choose UDP over TCP?" Latency-sensitive traffic (VoIP, games), one-shot request/response (DNS), or when you build your own reliability (QUIC). Follow-up: *"Is QUIC connection-oriented?"* Yes: it has a handshake and connection IDs. But it runs over UDP, so routers still see connectionless datagrams.

---

## Q7. [M1 | Circuit vs packet vs message switching delay | Medium | MCQ | OA]

A **10^6-bit** message goes from A to B through **2 switches (3 links)**. Each link: **1 Mbps**, propagation **1 ms**. Ignore processing, queuing and headers.
- **Circuit switching:** 50 ms setup, then the message streams through at 1 Mbps; switches do not store it.
- **Packet switching:** packets of 10,000 bits, stored and forwarded at each switch.

What is the total time until B receives the last bit, for circuit and packet switching?

- A. 1053 ms and 1023 ms
- B. 1003 ms and 1023 ms
- C. 1053 ms and 3003 ms
- D. 1050 ms and 1020 ms

**Answer:** A

**Why:** In a circuit, bits pass straight through the switches, so the message is transmitted only *once*. With packets, each switch stores and forwards, but packets overlap across links (**pipelining**). Each extra hop therefore adds only one *packet* time, not one *message* time.

**Working:**
```
Circuit:  T = setup + M/R + k*Tp = 50 + 1000 + 3*1 = 1053 ms

Packet:   n = 10^6 / 10^4 = 100 packets;  Tt = 10^4 / 10^6 = 10 ms
          last packet leaves A at n*Tt          = 1000 ms
          then (k-1) more transmissions        =   20 ms
          plus k propagation delays            =    3 ms
          T = (n + k - 1)*Tt + k*Tp = 102*10 + 3 = 1023 ms

Message switching (whole message stored per hop):
          T = k*M/R + k*Tp = 3000 + 3 = 3003 ms
```

**Why the other options fail:**
- B: Forgets the 50 ms setup.
- C: 3003 ms is *message* switching. It ignores pipelining.
- D: Leaves out propagation in both cases.

**Concept recap:**
- Circuit: T = T_setup + M/R + k·Tp
- Packet: T = (n + k − 1)·(p/R) + k·Tp
- Message: T = k·(M/R) + k·Tp

**Trap alert:** Multiplying *every* packet by the hop count (100 × 3 × 10 ms). Only the *last* packet's journey matters.

**Interview angle:** "Why did the Internet choose packet switching?" Statistical multiplexing (idle users cost nothing), no setup delay, rerouting around failures, pipelining. The cost: no guaranteed bandwidth, plus jitter. Follow-up: *"When is circuit-like behaviour still preferred?"* Constant-bit-rate traffic with strict latency, such as legacy telephony.

---

## Q8. [M1 | Optimal packet size | Hard | MCQ | GATE]

An **8000-bit payload** crosses **6 store-and-forward links** at **1 Mbps** each. It is split into equal packets of **p** payload bits, and every packet carries a **100-bit header**. Ignore propagation, processing and queuing. Which **p** minimizes the time until the last bit arrives?

- A. 200 bits
- B. 800 bits
- C. 400 bits
- D. 1000 bits
- E. 8000 bits (one packet)

**Answer:** C (12.5 ms)

**Why:** Smaller packets → more pipelining → less delay. But each packet carries a fixed header, so tiny packets waste capacity → more delay. Delay is U-shaped in p, with a square-root optimum.

**Working:**
```
T(p) = (X/p + k - 1) * (p + h) / R      X = 8000, h = 100, k = 6, R = 10^6

p = 200 : (40+5) * 300  = 13,500 bits -> 13.5 ms
p = 400 : (20+5) * 500  = 12,500 bits -> 12.5 ms   <- minimum
p = 800 : (10+5) * 900  = 13,500 bits -> 13.5 ms
p = 1000: ( 8+5) * 1100 = 14,300 bits -> 14.3 ms
p = 8000: ( 1+5) * 8100 = 48,600 bits -> 48.6 ms

Optimum: T*R = X + Xh/p + (k-1)p + (k-1)h
dT/dp = -Xh/p^2 + (k-1) = 0  ->  p* = sqrt(Xh/(k-1)) = sqrt(8000*100/5) = 400
```

**Why the other options fail:**
- A: 40 headers (4000 extra bits) cancel out the pipelining gain.
- B: Too few packets for good pipelining. It ties with 200, which shows the U shape.
- D: Even less pipelining.
- E: No pipelining at all. Each link carries all 8100 bits in turn.

**Concept recap:** p* = √(X·h/(k−1)). More hops → smaller optimal packets. Bigger headers → larger optimal packets.

**Trap alert:** "Smaller packets are always faster." That's true only while header overhead is negligible.

**Interview angle:** SDE version: "Why not send 1-byte packets for the lowest latency?" Per-packet header and processing overhead dominate. That is why Nagle's algorithm batches small writes, and why `TCP_NODELAY` is a deliberate trade-off.

---

## Q9. [M1 | Virtual-circuit switching delay | Medium | MCQ | GATE]

A VC network connects A to B over **3 links** (1 Mbps, 2 ms propagation each). A setup request travels hop by hop A→B and an accept returns B→A. Both are negligible in size and processing. After the accept arrives, A sends **5 packets of 1000 bits** (headers included), stored and forwarded at each switch. Ignore queuing. What is the total time from the start of setup until B gets the last bit?

- A. 13 ms
- B. 19 ms
- C. 33 ms
- D. 25 ms

**Answer:** D

**Why:** A VC first pays about one path RTT for setup. After that, data is pipelined exactly like datagrams. Each switch only looks up a short **VCI** (virtual-circuit identifier) instead of doing a full route lookup.

**Working:**
```
Setup:  6 ms (request) + 6 ms (accept)            = 12 ms
Data:   Tt = 1000/10^6 = 1 ms
        (n + k - 1)*Tt + k*Tp = (5+2)*1 + 3*2      = 13 ms
Total                                              = 25 ms

VC table at switch S1 (VCIs are link-local, rewritten per hop):
| In-port | In-VCI | Out-port | Out-VCI |
|    1    |   12   |    3     |   22    |
```

**Why the other options fail:**
- A: Data phase only (the datagram answer). Setup is ignored.
- B: Counts only the one-way setup. Data can't flow until the accept returns.
- C: 12 + 15 + 6. It sends the packets one after another over every hop, with no pipelining.

**Concept recap:** VC = setup (≈ 1 path RTT) + pipelined data + teardown. Per-VC state in the switches, in-order delivery, smaller headers.

**Trap alert:** A VCI is not globally unique. It is swapped at each hop, which is exactly how MPLS labels work.

**Interview angle:** "Datagram vs virtual-circuit networks?" Datagram: no per-flow state, robust, every packet routed on its full address. VC: setup cost, per-flow state, ordering, easier traffic engineering. Follow-up: *"Where does VC-style forwarding exist today?"* MPLS in ISP backbones.

---

## Q10. [M1 | Store-and-forward vs cut-through delay | Medium | MCQ | OA]

A **1500 B** frame (total size) crosses **3 switches (4 links)** at **1 Gbps**. Ignore propagation, processing and queuing. A cut-through switch starts forwarding after receiving the first **14 B** (the Ethernet header). What is the delivery time with (i) all store-and-forward and (ii) all cut-through switches?

- A. (i) 48 µs, (ii) 12.336 µs
- B. (i) 12 µs, (ii) 12 µs
- C. (i) 48 µs, (ii) 36.336 µs
- D. (i) 36 µs, (ii) 12.336 µs

**Answer:** A

**Why:** A store-and-forward switch waits for the *whole* frame, so every link adds a full frame time. A cut-through switch starts sending once it has the header, so each switch adds only the header time.

**Working:**
```
Tt(frame)  = 12,000 bits / 10^9 = 12 us
Tt(header) = 112 bits / 10^9    = 0.112 us

(i)  S&F:         4 * 12             = 48 us
(ii) Cut-through: 12 + 3 * 0.112     = 12.336 us
```

**Why the other options fail:**
- B: Treats the switches as having zero delay. The header must still be read.
- C: Adds a full frame time per hop in the cut-through case too.
- D: Counts 3 links instead of 4 (links = switches + 1).

**Concept recap:**
- S&F: T = k·(L/R) + ΣTp
- Cut-through: T ≈ L/R + (k−1)·(H/R) + ΣTp
- Cut-through works only if the output link is not faster than the input (otherwise it underruns).

**Trap alert:** Links = switches + 1. Also, cut-through can't check the FCS before forwarding.

**Interview angle:** "Store-and-forward vs cut-through?" Explain the latency difference and the error-check trade-off. Follow-up: *"Are routers cut-through?"* Usually not. They decrement TTL, update the checksum and queue packets.

---

## Q11. [M1 | S&F vs cut-through trade-offs | Easy | MSQ | Interview]

Select all that apply.

- A. Store-and-forward per-hop latency does not depend on frame size.
- B. A store-and-forward switch can detect a bad FCS and drop the frame before forwarding.
- C. Cut-through cuts per-hop latency to roughly the header-reception time.
- D. A cut-through switch can forward a frame whose FCS turns out to be bad.
- E. IP routers must use cut-through because the IP header is at the start of the packet.

**Answer:** B, C, D

**Why:** The FCS is at the **end** of the frame. A store-and-forward switch has it before forwarding and can drop a corrupt frame. A cut-through switch is already transmitting when the FCS arrives, so it forwards corrupt frames.

```
Frame     S&F per hop (1 Gbps)   Cut-through per hop
64 B      0.512 us               0.112 us
1518 B    12.144 us              0.112 us
```

**Why the other options fail:**
- A: S&F per-hop delay = L/R, proportional to frame size.
- E: Header position doesn't force cut-through. Routers typically store and forward.

**Concept recap:** S&F = L/R per hop, checks every frame. Cut-through = H/R per hop, forwards errors. Fragment-free mode waits for the first 64 B to filter out collision fragments.

**Trap alert:** Treating "faster" as "always better". Cut-through forwards corruption and can't go from a slow link to a fast one.

**Interview angle:** "How do HFT firms or data centres cut L2 latency?" Cut-through switches, kernel-bypass NICs, shorter paths. Follow-up: *"Where are corrupt frames caught in a cut-through fabric?"* At the receiving NIC, or at the next store-and-forward hop.

---

## Q12. [M1 | Mesh topology link count | Easy | MCQ | OA]

**10 routers** are connected in a full mesh. How many links are there, and how many **new** links are needed when an 11th router joins the mesh?

- A. 45 links; 11 new
- B. 90 links; 20 new
- C. 55 links; 10 new
- D. 45 links; 10 new

**Answer:** D

**Why:** A full mesh has one link per *pair* of nodes. A newcomer must connect to every existing node.

**Working:**
```
Links = n(n-1)/2:  n = 10 -> 45;  n = 11 -> 55;  new = 10
Ports per router = n - 1 = 9

Topology (n = 10)   Links   Single point of failure?
Full mesh            45     no
Star                 10     yes (hub)
Ring                 10     no for 1 failure; 2 failures split it
```

**Why the other options fail:**
- A: 11 would include a link to itself.
- B: n(n−1) counts each link twice.
- C: 55 is the total *after* the join.

**Concept recap:** A full mesh needs O(n²) links. That is why real designs are hierarchical or partial meshes (and why an iBGP full mesh is replaced by route reflectors).

**Trap alert:** Forgetting the ÷2. Links are unordered pairs.

**Interview angle:** "Why not fully mesh everything?" O(n²) cost. Follow-up: *"What's the microservices version of this problem?"* n² service-to-service connections and configs, which service meshes, connection pools or a message broker help tame.

---

## Q13. [M1 | Client-server vs P2P distribution time | Hard | MCQ | OA]

A **10^9-bit** file goes to **N = 100** peers. Server upload us = **100 Mbps**. Each peer: upload u = **10 Mbps**, download d = **20 Mbps**. The core has unlimited capacity; ignore propagation and overhead. What are the minimum distribution times for (i) client-server and (ii) P2P?

- A. (i) 1000 s, (ii) 50 s
- B. (i) 1000 s, (ii) 90.9 s
- C. (i) 50 s, (ii) 90.9 s
- D. (i) 1000 s, (ii) 10 s

**Answer:** B

**Why:** In client-server, the server alone uploads N copies. In P2P, every downloader also uploads, so capacity grows with N. The distribution time is the **maximum** of all the bounds that must hold.

**Working:**
```
Client-server: max(N*F/us, F/d) = max(10^11/10^8, 10^9/2*10^7)
             = max(1000, 50) = 1000 s

P2P: max(F/us, F/d, N*F/(us + N*u))
   = max(10, 50, 10^11 / (10^8 + 10^9))
   = max(10, 50, 90.9) = 90.9 s
```

**Why the other options fail:**
- A: Ignores the aggregate-upload bound in P2P.
- C: Ignores the fact that the client-server server must upload 100 copies.
- D: Uses only F/us, as if one copy from the server were enough for everyone.

**Concept recap:**
- D_cs ≥ max(N·F/us, F/dmin): linear in N.
- D_p2p ≥ max(F/us, F/dmin, N·F/(us + Σu)): levels off as N grows (self-scaling).

**Trap alert:** Taking the *minimum* of the bounds. A lower bound on time is the **max** of the necessary constraints.

**Interview angle:** "Why does BitTorrent scale?" Every downloader adds upload capacity. Follow-up: *"How do companies push large container images to thousands of hosts?"* P2P-style distributors, CDNs, hierarchical caches. All of them avoid the single-uplink bottleneck.

---

## Q14. [M1 | Unicast/broadcast/multicast/anycast | Easy | MCQ | OA + Interview]

A public DNS resolver uses one IP (e.g., 8.8.8.8) at dozens of sites worldwide. Users in Mumbai and Frankfurt both query that address, and each query reaches a *different*, routing-nearest site via normal BGP. Which delivery model is this?

- A. Multicast
- B. Broadcast
- C. Anycast
- D. Unicast with DNS round-robin

**Answer:** C

**Why:** **Anycast** advertises one address from many locations. Routing delivers each packet to **one** member, the nearest by routing metric.

```
Unicast    one -> one                normal TCP connection
Broadcast  one -> all on subnet      ARP request, DHCPDISCOVER (IPv4 only)
Multicast  one -> all group members  IPTV, OSPF hellos 224.0.0.5
Anycast    one -> nearest member     root DNS, public resolvers, CDN edges
```

**Why the other options fail:**
- A: Multicast delivers to *every* member of the group.
- B: Broadcast stays within one subnet and is never routed.
- D: Round-robin hands out *different* IPs. Here everyone uses the same IP.

**Concept recap:** Anycast = same address at many sites + routing picks one. IPv6 has no broadcast; it uses multicast (e.g., ff02::1) instead.

**Trap alert:** Anycast picks the *routing-nearest* site, not the least-loaded one. It is not a load balancer.

**Interview angle:** "How does a CDN send users to the nearest edge?" Anycast and/or geo-DNS. Follow-up: *"What's the risk of anycast for long TCP sessions?"* A BGP change can move the flow to a site with no connection state, which shows up as a reset.

---

## Q15. [M2 | Bandwidth vs latency | Easy | MCQ | OA + Interview]

Two offices are joined by **3000 km of fibre** (signal speed 2 × 10^8 m/s). The link is upgraded from **100 Mbps to 1 Gbps**. Ping RTT with **100 B** packets stays at about 30 ms. Why?

- A. The upgrade is misconfigured; 10× bandwidth should cut RTT about 10×
- B. RTT is dominated by propagation (≈15 ms each way), which depends on distance, not link rate; transmitting 100 B took only 8 µs even before the upgrade
- C. A higher bit rate raises signal speed only on copper, not fibre
- D. Routers rate-limit ICMP, so bandwidth changes never show in ping

**Answer:** B

**Why:** Bandwidth affects only **transmission** delay (L/R). Propagation (d/v) doesn't change. For small packets over long distances, propagation dominates.

**Working:**
```
Tp = 3*10^6 / 2*10^8 = 15 ms one way -> 30 ms RTT
Tt(800 bits): 8 us at 100 Mbps, 0.8 us at 1 Gbps
RTT: ~30.016 ms -> ~30.002 ms  (invisible in ping)
```

**Why the other options fail:**
- A: This is the core misconception: bandwidth ≠ latency.
- C: Signal speed depends on the medium (~2/3 c), never on the bit rate.
- D: ICMP rate-limiting caps how many replies a router sends, not RTT.

**Concept recap:** Total delay = transmission + propagation + queuing + processing. Bandwidth shrinks only the transmission term. Big transfers get faster; the first-byte latency doesn't.

**Trap alert:** "Faster internet = lower latency."

**Interview angle:** "Explain bandwidth vs latency." Use the analogy of highway lanes (bandwidth) vs highway length (latency). Follow-up: *"So what actually lowers latency for global users?"* Servers closer to users (CDN/edge), fewer round trips (TLS 1.3, connection reuse, 0-RTT), less queuing (bufferbloat control).

---

## Q16. [M2 | Throughput vs goodput | Medium | NAT | OA]

A TCP bulk transfer saturates one direction of a **100 Mbps** full-duplex Ethernet link with full **1460 B** payloads. Per frame on the wire: preamble+SFD 8 B, Ethernet header 14 B, IPv4 20 B, TCP 20 B, payload 1460 B, FCS 4 B, inter-frame gap 12 B. ACKs use the other direction. There is no loss. What is the **goodput** in Mbps, to 2 decimal places?

*Answer is a number — no options.*

**Answer:** ≈ **94.93 Mbps**

**Why:** Goodput counts only application bytes. Each 1460 B payload occupies a 1538 B slot on the wire.

**Working:**
```
Wire bytes per frame = 8 + 14 + 20 + 20 + 1460 + 4 + 12 = 1538 B
Goodput = 100 * 1460/1538 = 94.928 ~ 94.93 Mbps

Same link measured at different layers:
  line rate             100.00 Mbps
  Ethernet (1518/1538)   98.70 Mbps
  IP       (1500/1538)   97.53 Mbps
  TCP payload = goodput  94.93 Mbps
```

**Common wrong answers:**
- 100: line rate, not goodput.
- 96.18 (1460/1518): forgets the preamble and gap.
- 97.33 (1460/1500): counts only the IP+TCP headers as overhead.

**Concept recap:** Goodput ≤ throughput ≤ bandwidth. Throughput must always be quoted *at a specific layer*.

**Trap alert:** A 1 Gbps NIC showing ~941 Mbps in `iperf` is correct, not a fault.

**Interview angle:** "My 1 Gbps link shows 940 Mbps in iperf. Is something broken?" No: ~5% header overhead at a 1500 MTU. Follow-up: *"What would jumbo frames change?"* Overhead drops to ~1%, but every device on the path must support 9000 MTU.

---

## Q17. [M2 | Bottleneck throughput | Medium | MSQ | OA + Interview]

Path: server NIC **10 Gbps** → ISP link **200 Mbps** → client Wi-Fi **50 Mbps**. There is no other traffic unless an option says so. Select all that apply.

- A. A single long transfer's throughput is bounded by the bottleneck, 50 Mbps.
- B. Upgrading the server NIC to 25 Gbps would raise this client's throughput.
- C. If 10 long flows (including this one) share the ISP link equally, this flow is bounded by min(50, 200/10) = 20 Mbps.
- D. A single small packet's latency is set only by the bottleneck bandwidth.
- E. Even with all links at 10 Gbps, one TCP flow can still be capped at window ÷ RTT.

**Answer:** A, C, E

**Why:** Throughput is the minimum, over all links, of each link's *effective share* for this flow. Separately, TCP sends at most one window per RTT.

**Working:**
```
Alone:   min(10000, 200, 50)    = 50 Mbps
Shared:  min(10000, 200/10, 50) = 20 Mbps
Window cap: 64 KB = 524,288 bits, RTT 50 ms -> 524,288/0.05 = 10.49 Mbps
```

**Why the other options fail:**
- B: The NIC isn't the bottleneck. Upgrading a non-bottleneck link changes nothing.
- D: Small-packet latency is mostly propagation + queuing. Bandwidth affects only the small transmission term.

**Concept recap:** Throughput ≈ min(link rate ÷ number of sharing flows), and also ≤ window/RTT. Find the bottleneck before optimizing anything.

**Trap alert:** Option E is the "long fat pipe" case: a large bandwidth-delay product needs TCP window scaling.

**Interview angle:** "A user on 1 Gbps downloads from another continent at 30 Mbps. Why?" Possible causes: a bottleneck elsewhere, a small window with high RTT, loss triggering congestion control, congested peering. Follow-up: *"How would you diagnose it?"* `iperf`, `mtr`, RTT and loss checks, then look at the window and retransmissions in a capture.

---

## Q18. [M2 | Transmission delay (units) | Easy | MCQ | OA]

A **2 KB** packet (1 KB = 1024 B) is sent on an **8 Mbps** link (1 Mbps = 10^6 bps). What is the transmission delay?

- A. 2.000 ms
- B. 1.953 ms
- C. 2.048 ms
- D. 0.256 ms

**Answer:** C

**Why:** Tt = bits ÷ rate. The stated conventions decide the third digit.

**Working:**
```
L = 2*1024*8 = 16,384 bits;  R = 8*10^6 bps
Tt = 16,384 / 8,000,000 = 2.048 ms
```

**Why the other options fail:**
- A: Uses 1 KB = 1000 B.
- B: Uses 1 Mbps = 2^20 bps. Link rates are always decimal.
- D: Forgot to multiply by 8 to convert bytes to bits.

**Concept recap:** Tt = L/R. Link rates are decimal; data sizes follow whatever the question states. Always convert bytes to bits (×8).

**Trap alert:** Forgetting ×8 is the most common OA slip. Write the unit next to every number.

**Interview angle:** SDE version: "How long does a 1 MB API response take on a 10 Mbps link?" ≈ 0.8 s of transmission plus RTTs, which is why you compress, paginate and use a CDN.

---

## Q19. [M2 | Propagation delay | Easy | MCQ | OA]

Stations A and B talk via a geostationary satellite at **36,000 km** altitude. Both stations are directly beneath it, and the signal speed is **3 × 10^8 m/s**. What is the one-way propagation delay A→B?

- A. 240 ms
- B. 120 ms
- C. 480 ms
- D. 0.24 ms

**Answer:** A

**Why:** The signal goes up 36,000 km *and* down 36,000 km. Propagation depends only on distance and medium.

**Working:**
```
Path = 72,000 km = 7.2*10^7 m
Tp = 7.2*10^7 / 3*10^8 = 240 ms;  RTT (A->B->A) = 480 ms
```

**Why the other options fail:**
- B: Counts only one leg.
- C: That's the RTT.
- D: km/m units slip.

**Concept recap:** Tp = d/v. GEO one-way ≈ 240–280 ms. LEO (~550 km altitude) ≈ a few ms per leg, which is why LEO broadband feels interactive.

**Trap alert:** "Via satellite" means up *and* down.

**Interview angle:** "Why are interactive apps painful over GEO links?" ~500 ms RTT is paid on every TCP handshake, TLS handshake and request, and slow start ramps up slowly. Follow-up: *"How do you mitigate that?"* Fewer round trips (TLS 1.3, QUIC 0-RTT, connection reuse), performance-enhancing proxies, a larger initial window.

---

## Q20. [M2 | Transmission = propagation crossover | Medium | NAT | GATE]

A **2 km** link runs at **100 Mbps** with signal speed **2 × 10^8 m/s**. For what packet length **L (in bits)** does transmission delay equal propagation delay?

*Answer is a number — no options.*

**Answer:** **1000 bits** (= 125 B)

**Why:** Set L/R = d/v, so L = R·Tp. That is the number of bits "in flight" on the wire at any instant, i.e. the link's bandwidth-delay product.

**Working:**
```
Tp = 2000 / 2*10^8 = 10 us
L  = R * Tp = 10^8 * 10^-5 = 1000 bits

a = Tp/Tt:  L = 1000  -> a = 1
            L = 12000 -> a = 0.083 (transmission dominates)
            L = 100   -> a = 10    (propagation dominates)
```

**Common wrong answers:**
- 125: gives bytes when bits were asked.
- 10: gives Tp in µs.
- 100,000: km/m units slip.

**Concept recap:** **a = Tp/Tt** drives the data-link numericals that come next:
- Stop-and-wait efficiency = 1/(1 + 2a)
- CSMA/CD needs Tt ≥ 2Tp
- BDP = R·Tp

**Trap alert:** a is propagation ÷ transmission. Inverting it flips every efficiency formula.

**Interview angle:** The intuition interviewers want: on a short, slow link the packet is longer than the wire; on a long, fast link the wire holds many packets, so you need pipelining (sliding windows). Follow-up: *"How big must the TCP window be to fill such a pipe?"* At least the BDP (next batch, P06).

---

**Completed:** 20 / 340   **Remaining:** 320
**Ledger (new ✓ this batch):** P01 ×3 (Q7–Q9) · P02 ×2 (Q10, Q11) · P03 ×2 (Q18, Q20) · P04 ×2 (Q19, Q20)
**Batch split:** Easy/Medium/Hard 8/10/2 · MCQ/MSQ/NAT 14/4/2 · MCQ answer key A/B/C/D = 4/3/4/3
