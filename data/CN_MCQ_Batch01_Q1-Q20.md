# Computer Networks Question Bank — Batch 01 (Q1–Q20)

**Modules:** M1 Fundamentals, layering, switching (Q1–Q14) · M2 Physical layer, delay & capacity — part 1 (Q15–Q20)
**Patterns targeted:** P01, P02, P03, P04 (+ layering / encapsulation / demux / topology / delivery-model groundwork, and bandwidth-vs-throughput-vs-latency traps)

**Global conventions for this batch (restated inside each numerical):** 1 Mbps = 10^6 bps, 1 Gbps = 10^9 bps; data sizes are stated explicitly; IPv4 header = 20 B (no options); TCP header = 20 B (no options); Ethernet header = 14 B + FCS 4 B (= 18 B overhead), preamble+SFD = 8 B and inter-frame gap = 12 B only when the question says "on the wire".

---

### Q1. [M1 | Layering rationale | Easy | MCQ | OA + Interview]

A laptop's network connection is switched from Wi-Fi (802.11) to wired Ethernet (802.3). The browser, the HTTP library, the TCP/IP stack and every open web page keep working with no code changes. Which property of layered network design best explains this?

A. Lower layers repeat the functions of the upper layers, so if one layer fails another can take over.
B. Layering keeps the total header overhead as small as possible.
C. Every layer has to be built in hardware, so swapping hardware is transparent.
D. Each layer offers a service to the layer above through a fixed interface, so a layer's implementation can be replaced without changing the layers above it.

**Answer:** D

**Why:** In a layered architecture each layer is a "black box" to the layer above it. The network layer (IP) only needs a "send this packet to that next-hop link address" service from the link layer. Whether Wi-Fi or Ethernet provides that service is hidden behind the interface. This property is called **modularity / separation of concerns**: you can change the implementation as long as the interface and the service stay the same.

**Working:** Here is what changes and what doesn't when Wi-Fi is replaced by Ethernet:

```
Layer        Wi-Fi case                 Ethernet case              Changed?
Application  HTTP                       HTTP                       No
Transport    TCP                        TCP                        No
Network      IP (same IP stack)         IP                         No (address may change via DHCP, but code doesn't)
Link         802.11 MAC, CSMA/CA, ACKs  802.3 MAC, full-duplex     YES
Physical     Radio (2.4/5 GHz)          Twisted pair               YES
```

Only the bottom two layers changed. IP still hands a packet down and gets packets back up through the same interface.

**Why the other options fail:**
- A — Layers don't duplicate each other as backup. Some functions do repeat (for example, error detection at both the link and transport layers), but they exist for end-to-end correctness, not as failover.
- B — Layering actually *adds* overhead, because every layer puts on its own header. Overhead is the price of modularity.
- C — Most layers (TCP, IP, HTTP) are written in software. Layering is about interfaces, not about hardware.

**Concept recap:** Layer N uses the service of layer N−1 and provides a service to layer N+1. Peer layers on different hosts talk through *protocols*; adjacent layers on the same host talk through *interfaces*. Benefits: modularity, easier debugging, independent evolution. Cost: header overhead, and sometimes duplicated functionality.

**Trap alert:** Students mix up *protocol* (horizontal: peer-to-peer rules, e.g., TCP on host A ↔ TCP on host B) with *interface* (vertical: service between adjacent layers on the same host). The swap in this question only works because the *interface* stays the same.

**Interview angle:** "Why do we layer network stacks?" Model answer: "To split the problem into independent pieces. Each layer solves one problem (moving bits, delivering on one link, routing between networks, reliable delivery between processes, application semantics) and exposes a stable service upward. That lets us swap Wi-Fi for Ethernet, or HTTP/1.1 for HTTP/2, without touching the other layers." Likely follow-up: "Name a case where layering is violated in practice." Good answers: NAT rewrites transport-layer ports even though it is a network-layer device; TCP checksum uses a pseudo-header that contains IP addresses; TCP MSS is chosen based on the link MTU.

---

### Q2. [M1 | OSI vs TCP/IP | Easy | MCQ | OA + Interview]

The OSI model has 7 layers; the TCP/IP (Internet) model is usually described with 4 or 5. Which OSI layers' responsibilities are absorbed into the **application layer** of the TCP/IP model?

A. Transport and Network
B. Data Link and Physical
C. Session and Presentation
D. Network and Session

**Answer:** C

**Why:** TCP/IP has no separate session layer (managing dialogs and checkpoints) and no separate presentation layer (data format, encoding, compression, encryption). Applications handle those jobs themselves, for example HTTP cookies for sessions and JSON/UTF-8/gzip for presentation, or they use libraries such as TLS. So OSI layers 5, 6 and 7 together correspond to the TCP/IP application layer.

**Working:** Mapping:

```
OSI (7)          TCP/IP (5-layer view)     Example
7 Application  \
6 Presentation  >  Application            HTTP, DNS, gzip, JSON, (TLS sits here-ish)
5 Session      /
4 Transport    ->  Transport              TCP, UDP, QUIC (over UDP)
3 Network      ->  Network (Internet)     IP, ICMP
2 Data Link    ->  Link                   Ethernet, 802.11, ARP*
1 Physical     ->  Physical               Copper, fibre, radio
(*ARP is usually placed between L2 and L3)
```

**Why the other options fail:**
- A — Transport and Network exist as their own layers in TCP/IP (TCP/UDP and IP).
- B — Data Link and Physical are merged into the "link / network access" layer in the 4-layer TCP/IP model, not into the application layer. This option catches students who remember "something gets merged" but not which layers.
- D — Network is a distinct layer in both models.

**Concept recap:** OSI is a *reference model* (7 layers, designed before protocols). TCP/IP is a *protocol suite* (4–5 layers, described after the protocols were built). OSI L5+L6+L7 → TCP/IP Application; OSI L1+L2 → TCP/IP Link (in the 4-layer view).

**Trap alert:** "TLS is the presentation layer" is not standard in either model. In practice TLS runs *on top of* TCP and *below* HTTP. Interviewers accept "between transport and application"; don't claim it is strictly L6.

**Interview angle:** "Compare OSI and TCP/IP." Model answer: "OSI is a 7-layer conceptual model; TCP/IP is the 4–5 layer model the Internet actually runs. TCP/IP folds session and presentation into the application layer. In practice we debug using TCP/IP terms but use OSI numbers as shorthand, like 'L4 load balancer' or 'L7 proxy'." Follow-up: "What is the difference between an L4 and an L7 load balancer?" (L4 routes on IP and port; L7 parses HTTP and routes on path, headers or cookies.)

---

### Q3. [M1 | PDUs per layer | Medium | MSQ | OA]

Select all that apply. Which of the following layer → PDU (protocol data unit) pairings are correct?

A. Transport layer using TCP → segment
B. Network layer → frame
C. Data link layer → frame
D. Transport layer using UDP → datagram
E. Physical layer → packet

**Answer:** A, C, D

**Why:** The PDU is the "unit" a layer works with: its header plus the payload it carries. TCP's PDU is a *segment*. UDP's is a *datagram* (RFC 768 is titled "User Datagram Protocol"). IP's is a *packet* (also called an IP datagram). The link layer's is a *frame*. The physical layer works on *bits/symbols*.

**Working:** Encapsulation view:

```
Application data: "GET / HTTP/1.1..."
+ TCP header  -> TCP segment              (L4 PDU)
+ IP header   -> IP packet / IP datagram  (L3 PDU)
+ Eth hdr+FCS -> Ethernet frame           (L2 PDU)
-> encoded as bits / symbols on the medium (L1)
```

**Why the other options fail:**
- B — "Frame" belongs to layer 2. Layer 3's unit is a packet; mixing these up suggests unclear encapsulation boundaries.
- E — The physical layer carries bits, not packets. A "packet" exists only at L3.

**Concept recap:** Application → message · TCP → segment · UDP → datagram · IP → packet (datagram) · Link → frame · Physical → bits. Each layer's PDU becomes the *payload (SDU)* of the layer below.

**Trap alert:** "IP datagram" and "UDP datagram" both use the word *datagram*. Questions that say "datagram" without saying which layer are ambiguous, so check the layer before answering.

**Interview angle:** "Walk me through encapsulation when you send an HTTP request." Mention each header being added and which device reads which header: switches read L2, routers read L3, and NAT/L4 load balancers read L4. Follow-up: "Which headers change hop by hop?" L2 headers are rebuilt at every router. L3 changes only TTL and checksum (and addresses if NAT). L4 is end-to-end unless NAT/PAT rewrites ports.

---

### Q4. [M1 | Encapsulation & overhead | Medium | MCQ | OA + Interview]

A TCP sender places **1460 bytes** of application data in one segment. Assume a TCP header of 20 B (no options), an IPv4 header of 20 B (no options), and Ethernet II framing with a 14 B header plus a 4 B FCS. Ignore the preamble and inter-frame gap. What are (i) the value of the IPv4 **Total Length** field and (ii) the Ethernet frame size (header + payload + FCS)?

A. (i) 1480 B, (ii) 1500 B
B. (i) 1500 B, (ii) 1518 B
C. (i) 1460 B, (ii) 1518 B
D. (i) 1500 B, (ii) 1514 B

**Answer:** B

**Why:** The IPv4 Total Length field counts the *whole IP packet*: the IP header plus everything after it (TCP header + data). The Ethernet frame then wraps that packet with a 14 B header in front and a 4 B FCS at the end. 1460 is exactly the standard **MSS** (maximum segment size: TCP payload only) for a 1500 B Ethernet **MTU** (maximum IP packet size).

**Working:**

```
TCP payload (MSS)          = 1460 B
+ TCP header               =   20 B  -> TCP segment = 1480 B
+ IPv4 header              =   20 B  -> IP packet   = 1500 B   <- IP Total Length = 1500 (= MTU)
+ Ethernet header          =   14 B
+ Ethernet FCS             =    4 B  -> Ethernet frame = 1518 B
(on the wire add 8 B preamble+SFD and 12 B inter-frame gap -> 1538 B "slot")

Payload efficiency (frame level) = 1460 / 1518 = 0.9618 -> 96.2 %
```

**Why the other options fail:**
- A — Counts TCP + data as the IP Total Length but leaves out the IP header itself. Total Length includes the IP header.
- C — Treats the IP Total Length as the application payload, which confuses MSS with IP length.
- D — 1514 leaves out the 4 B FCS. Many packet captures show 1514 because the NIC strips the FCS before handing the frame to the capture tool, but the frame on the medium is 1518.

**Concept recap:**
MTU = max IP packet size (IP header + transport header + data) = 1500 B on Ethernet.
MSS = MTU − IP header − TCP header = 1500 − 20 − 20 = 1460 B (no options).
Ethernet frame = 14 + IP packet + 4 → 64 B minimum, 1518 B maximum (untagged).

**Trap alert:** **MSS is payload only; MTU is the full IP packet.** Another trap: Wireshark shows "1514 bytes on wire" for a full-size frame because the FCS is usually not captured. Don't "correct" the standard to 1514.

**Interview angle:** "What's the difference between MTU and MSS, and why does it matter?" Model answer: "MTU is the largest IP packet a link can carry; MSS is the largest TCP payload. TCP announces its MSS in the SYN so the peer never sends segments that would need IP fragmentation." Follow-up: "What happens with a VPN tunnel that adds 50–60 B of headers?" The effective MTU drops, so you need MSS clamping or PMTUD; otherwise you get the "MTU black hole" (pattern P89, covered later).

---

### Q5. [M1 | Multiplexing / demultiplexing per layer | Medium | MCQ | Interview]

A host has one Ethernet NIC and runs a browser (two HTTPS connections to the same server) and an SSH session at the same time. An Ethernet frame arrives. Which sequence of fields does the receiving stack use, layer by layer, to hand the payload to the correct socket?

A. Ethernet **EtherType** → IPv4 **Protocol** field → TCP **4-tuple** (src IP, src port, dst IP, dst port)
B. Destination MAC address → IPv4 TTL → TCP sequence number
C. Ethernet EtherType → IPv4 Type-of-Service → TCP destination port only
D. 802.1Q VLAN tag → IPv4 Identification → TCP window size

**Answer:** A

**Why:** Each layer has a "which upper-layer protocol/instance is this for?" field:
- EtherType 0x0800 means the payload is IPv4 (0x86DD is IPv6, 0x0806 is ARP).
- IP Protocol 6 means TCP (17 is UDP, 1 is ICMP).
- For an *established* TCP connection, the stack matches the full 4-tuple. That is essential here: both HTTPS connections use destination port 443 on the server side and local port 443 is not involved, so only the client's ephemeral port tells them apart.

**Working:** Demux trace for an incoming packet from server 203.0.113.10:443 to client 192.168.1.5:51234:

```
| Layer    | Field used        | Value seen        | Decision                          |
|----------|-------------------|-------------------|-----------------------------------|
| Ethernet | EtherType         | 0x0800            | hand to IPv4                      |
| IPv4     | Protocol          | 6                 | hand to TCP                       |
| TCP      | 4-tuple lookup    | 203.0.113.10:443  | match socket of HTTPS conn #1     |
|          |                   | -> 192.168.1.5:51234 | (conn #2 uses local port 51235) |
```

**Why the other options fail:**
- B — The destination MAC only decides whether the NIC accepts the frame; TTL is about hop limits; the sequence number is about ordering within a connection, not about finding it.
- C — ToS/DSCP is a QoS marking, not a protocol selector. Destination port alone cannot separate the two HTTPS connections (on the client both are "some local ephemeral port", and on a server every client hits port 443).
- D — The VLAN tag selects a broadcast domain, IP Identification is for reassembling fragments, and the window field is for flow control. None of them demultiplex.

**Concept recap:** Demux keys: L2 EtherType · L3 Protocol (IPv4) / Next Header (IPv6) · L4: UDP → (dst IP, dst port); TCP listening socket → (dst IP, dst port); TCP established socket → (src IP, src port, dst IP, dst port).

**Trap alert:** "TCP demultiplexes by destination port" is only true for *listening* sockets. A server accepts thousands of connections all on port 443; each accepted socket is identified by the 4-tuple. That is why one server port can serve many clients.

**Interview angle:** "How can a web server handle 10,000 connections on port 443?" Answer: each connection is a distinct 4-tuple; the listening socket only matches new SYNs, and `accept()` returns a new socket per 4-tuple. Follow-up: "Then what limits connections from one client to one server:port?" The client's ephemeral port range (about 28k by default on Linux), which leads to port exhaustion (P89).

---

### Q6. [M1 | Connection-oriented vs connectionless | Medium | MSQ | OA + Interview]

Select all that apply.

A. A connectionless service can never deliver packets in order.
B. TCP sets up state at the two endpoints before data flows; ordinary IP routers on the path keep no per-TCP-connection state.
C. Virtual-circuit networks (e.g., ATM, Frame Relay) keep per-connection state in the intermediate switches.
D. "Connection-oriented" means the network has reserved bandwidth for the connection.
E. UDP datagrams between the same pair of sockets may arrive out of order, duplicated, or not at all.

**Answer:** B, C, E

**Why:** "Connection" can mean two different things. TCP's connection is an *end-to-end* agreement: sequence numbers, windows and buffers exist only at the two hosts, while the IP network in between stays connectionless (each packet is routed independently). A virtual-circuit network, by contrast, installs per-connection state *inside* the switches. UDP adds no ordering, de-duplication or retransmission on top of IP, so it inherits all of IP's best-effort behaviour.

**Working:** Where connection state lives:

```
Service              State at end hosts?   State in network nodes?   Ordering guaranteed?
UDP over IP          minimal (socket)      none                      no
TCP over IP          yes (seq, windows)    none                      yes (by end hosts)
Virtual circuit      yes                   yes (VC tables)           yes (single path)
Circuit switching    yes                   yes (reserved capacity)   yes
```

**Why the other options fail:**
- A — Connectionless packets often *do* arrive in order (for example, when they all take the same path). Order is simply not *guaranteed*.
- D — TCP is connection-oriented but reserves no bandwidth; it competes through congestion control. Reservation is a property of circuit switching (or of explicit QoS reservation).

**Concept recap:** Connection-oriented = a setup phase creates shared state before data flows, followed by teardown. Connectionless = each unit is independent. Reliability, ordering and reservation are *separate* properties; don't merge them with "connection-oriented".

**Trap alert:** "TCP guarantees delivery" is false. TCP guarantees an *ordered, reliable byte stream or an error to the application* (connection reset or timeout). If the network is cut, data is not delivered; the application is told that the connection failed.

**Interview angle:** "TCP vs UDP: when would you choose UDP?" Answer: when latency matters more than perfect delivery (VoIP, games, live video), for one-shot request/response (DNS), or when you want your own reliability on top (QUIC). Follow-up: "Is QUIC connection-oriented?" Yes. It has a handshake and connection IDs, but runs over UDP, so the *network* still sees connectionless datagrams.

---

### Q7. [M1 | Circuit vs packet vs message switching delay | Medium | MCQ | OA]

A message of **10^6 bits** goes from host A to host B through **2 intermediate switches** (so 3 links). Every link runs at **1 Mbps (10^6 bps)** and has a propagation delay of **1 ms**. Ignore processing and queuing delays and header bits.

- **Circuit switching:** circuit setup takes 50 ms. After setup, the message is sent as one continuous bit stream at 1 Mbps (switches do not store and forward).
- **Packet switching:** the message is split into packets of 10,000 bits each, and every switch stores and forwards each packet.

What is the total time from the start (setup start, for circuit switching) until B receives the last bit, for circuit and packet switching respectively?

A. 1053 ms and 1023 ms
B. 1003 ms and 1023 ms
C. 1053 ms and 3003 ms
D. 1050 ms and 1020 ms

**Answer:** A

**Why:** In circuit switching the bits flow through the switches without being stored, so the message is transmitted only once at link rate, plus setup and three propagation delays. In packet switching, every packet is stored and forwarded at each switch, but packets overlap in time across links (pipelining). As a result only one extra *packet* transmission time is added per extra hop, not one extra *message* transmission time.

**Working:**

```
Given: M = 10^6 bits, R = 10^6 bps, k = 3 links, Tp = 1 ms/link, setup = 50 ms

Circuit switching:
  T = setup + M/R + k*Tp
    = 50 ms + 1000 ms + 3 ms = 1053 ms

Packet switching (store-and-forward, pipelined):
  packet size p = 10^4 bits -> n = M/p = 100 packets
  Tt(packet) = p/R = 10^4 / 10^6 = 10 ms
  Last packet leaves A at        n*Tt          = 100*10  = 1000 ms
  It then needs (k-1) more transmissions  (k-1)*Tt = 2*10 = 20 ms
  Plus propagation on k links    k*Tp          = 3*1    =    3 ms
  T = (n + k - 1)*Tt + k*Tp = 102*10 + 3 = 1023 ms

Message switching (for comparison - whole message stored at each switch):
  T = k*(M/R) + k*Tp = 3*1000 + 3 = 3003 ms
```

**Why the other options fail:**
- B — 1003 ms forgets the 50 ms circuit setup time.
- C — 3003 ms is *message* switching (whole message stored at each hop). It ignores the pipelining that packetization buys.
- D — Leaves out propagation delay in both cases (1053 − 3, 1023 − 3).

**Concept recap:**
Circuit: T = T_setup + M/R + k·Tp.
Packet (store-and-forward, n equal packets, k links, equal rates): T = (n + k − 1)·(p/R) + k·Tp.
Message switching: T = k·(M/R) + k·Tp.

**Trap alert:** Students multiply the per-packet time by the number of hops for *every* packet (100 × 3 × 10 ms), forgetting that packets are pipelined. Only the *last* packet's journey matters: it leaves the source at n·Tt and then needs (k−1) more transmissions.

**Interview angle:** "Why did the Internet choose packet switching over circuit switching?" Answer: statistical multiplexing (idle users consume no capacity), no setup latency, robustness to failures (rerouting), and the pipelining shown above. The cost: no guaranteed bandwidth, plus variable queuing delay (jitter). Follow-up: "When would circuit-like behaviour still be preferred?" For constant-bit-rate traffic with strict latency needs, such as old telephony and some industrial or real-time links.

---

### Q8. [M1 | Packet size vs end-to-end delay (header overhead trade-off) | Hard | MCQ | GATE]

A message with **8000 bits of payload** travels over a path of **6 store-and-forward links** (5 intermediate routers). All links run at **1 Mbps (10^6 bps)**. The message is split into equal packets carrying **p** payload bits each; every packet also carries a **100-bit header**. Ignore propagation, processing and queuing delays (propagation adds the same constant whatever p is). Which payload size per packet gives the **minimum** time until the destination receives the last bit?

A. 200 bits
B. 800 bits
C. 400 bits
D. 1000 bits
E. 8000 bits (single packet)

**Answer:** C (400 bits → 12.5 ms)

**Why:** Smaller packets give more pipelining across the 6 links, which lowers delay. But every packet carries a fixed 100-bit header, so very small packets waste capacity, which raises delay. Total delay is (number of packets + extra hops) × (time to send one packet including its header), and minimizing it gives a square-root optimum.

**Working:**

```
X = 8000 payload bits, h = 100 header bits, k = 6 links, R = 10^6 bps
n = X/p packets, each of size (p + h) bits
T(p) = (n + k - 1) * (p + h) / R

p = 200  : n = 40 -> (40+5) * 300  = 13,500 bits -> 13.5 ms
p = 400  : n = 20 -> (20+5) * 500  = 12,500 bits -> 12.5 ms   <- minimum
p = 800  : n = 10 -> (10+5) * 900  = 13,500 bits -> 13.5 ms
p = 1000 : n = 8  -> ( 8+5) * 1100 = 14,300 bits -> 14.3 ms
p = 8000 : n = 1  -> ( 1+5) * 8100 = 48,600 bits -> 48.6 ms

Analytical optimum:
T(p)*R = (X/p + k-1)(p + h) = X + X*h/p + (k-1)*p + (k-1)*h
dT/dp = 0  ->  -X*h/p^2 + (k-1) = 0  ->  p* = sqrt(X*h/(k-1))
p* = sqrt(8000*100/5) = sqrt(160,000) = 400 bits  (check: 8000/400 = 20 packets, an integer)
```

**Why the other options fail:**
- A — 200 bits: too many packets (40), so header overhead (40 × 100 = 4000 extra bits) cancels out the pipelining gain.
- B — 800 bits: too few packets, so less pipelining. It ties with 200 at 13.5 ms, showing the curve is U-shaped.
- D — 1000 bits: even less pipelining, 14.3 ms.
- E — Single packet: no pipelining at all; each of the 6 links must carry the full 8100 bits one after another.

**Concept recap:** T(p) = (X/p + k − 1)(p + h)/R; optimum p* = √(X·h/(k − 1)). More hops → smaller optimal packets. Bigger headers → larger optimal packets.

**Trap alert:** Assuming "smaller packets are always faster". That holds only while header overhead is negligible. Always include the header in the per-packet transmission time.

**Interview angle:** Rarely asked verbally in that form. The SDE version is: "Why not send 1-byte packets for the lowest latency?" Answer: per-packet header and processing overhead dominate. This is why Nagle's algorithm batches small writes, and why TCP_NODELAY (P85) is a deliberate latency-vs-efficiency choice.

---

### Q9. [M1 | Virtual-circuit switching delay | Medium | MCQ | GATE]

A virtual-circuit (VC) network connects A to B over **3 links** (2 switches). Each link: **1 Mbps (10^6 bps)**, propagation delay **2 ms**. Before any data is sent, a VC must be set up: a setup request travels hop by hop from A to B, and an accept message travels back from B to A. Treat setup and accept messages as having negligible transmission and processing time. Once A receives the accept, it sends **5 data packets of 1000 bits each** (headers included), stored and forwarded at each switch. Ignore queuing. What is the total time from the start of setup until B receives the last bit of the 5th packet?

A. 13 ms
B. 19 ms
C. 33 ms
D. 25 ms

**Answer:** D

**Why:** A VC pays a round trip of setup delay before any data moves (request forward, accept back). After that, data packets are pipelined just like datagram packets, because each switch looks up a short VC identifier in its VC table instead of doing a full route lookup.

**Working:**

```
Setup:   request A->B = 3 links * 2 ms = 6 ms
         accept  B->A = 3 links * 2 ms = 6 ms
         T_setup = 12 ms

Data:    Tt = 1000 bits / 10^6 bps = 1 ms per packet per link
         n = 5 packets, k = 3 links
         T_data = (n + k - 1)*Tt + k*Tp = (5 + 2)*1 + 3*2 = 7 + 6 = 13 ms

Total  = 12 + 13 = 25 ms

VC table at switch S1 (example; VCIs are link-local and rewritten each hop):
| In-port | In-VCI | Out-port | Out-VCI |
|---------|--------|----------|---------|
|   1     |   12   |    3     |   22    |
```

**Why the other options fail:**
- A — 13 ms is the data phase only, i.e., the datagram answer. It ignores VC setup.
- B — 19 ms counts only the one-way setup (6 ms). Data cannot flow until the accept comes back.
- C — 33 ms = 12 + (5 × 3 × 1) + 6: it forgets pipelining and sends each packet across all hops one after another.

**Concept recap:** VC switching = setup (≈ 1 RTT along the path) + pipelined data + teardown. Switches keep per-VC state (in-port, in-VCI → out-port, out-VCI). All packets follow one path, so they arrive in order. Headers are smaller (VCI instead of full addresses).

**Trap alert:** A VCI is *not* globally unique. It is local to each link and is swapped hop by hop. That is exactly how MPLS labels work today.

**Interview angle:** "Datagram vs virtual-circuit networks?" Answer: datagram = no per-flow state in routers, each packet routed on its full destination address, robust to failures. VC = setup cost, per-flow state in switches, in-order delivery, easier QoS and traffic engineering. Follow-up: "Where does VC-style forwarding exist in modern networks?" In MPLS label switching inside ISP backbones.

---

### Q10. [M1 | Store-and-forward vs cut-through delay | Medium | MCQ | OA]

A **1500-byte Ethernet frame** (total frame size) goes from host A to host B through **3 Ethernet switches** (so 4 links). Every link runs at **1 Gbps (10^9 bps)**. Ignore propagation, processing and queuing delays. A **cut-through** switch starts forwarding a frame once it has received the first **14 bytes** (the Ethernet header, which contains the destination MAC). What is the time until B receives the last bit, with (i) all switches store-and-forward and (ii) all switches cut-through?

A. (i) 48 µs, (ii) 12.336 µs
B. (i) 12 µs, (ii) 12 µs
C. (i) 48 µs, (ii) 36.336 µs
D. (i) 36 µs, (ii) 12.336 µs

**Answer:** A

**Why:** A store-and-forward switch must receive the *entire* frame before it begins sending it out, so every link adds a full frame transmission time. A cut-through switch starts sending once it has read the header, so each switch adds only the header-reception time. The frame's body streams through.

**Working:**

```
Frame = 1500 B = 12,000 bits;  R = 10^9 bps
Tt(frame)  = 12,000 / 10^9 = 12 us
Tt(header) = 14*8 = 112 bits / 10^9 = 0.112 us

(i) Store-and-forward, 4 links:
    T = 4 * 12 us = 48 us

(ii) Cut-through, 3 switches:
    A transmits the whole frame once:    12 us
    each switch adds header latency:     3 * 0.112 = 0.336 us
    T = 12 + 0.336 = 12.336 us
```

**Why the other options fail:**
- B — Treats the switches as having zero delay. Even cut-through must read the header first.
- C — Adds a full frame time per hop in the cut-through case too (3 × 12 + 0.336). That is store-and-forward thinking applied to cut-through.
- D — 36 µs counts only 3 links (one per switch) for store-and-forward. With 3 switches there are 4 links.

**Concept recap:**
Store-and-forward, k links, equal rates: T = k·(L/R) + Σ Tp.
Cut-through: T ≈ L/R + (k − 1)·(H/R) + Σ Tp, where H = bytes needed to make the forwarding decision.
Cut-through works only when the output link is not faster than the input (otherwise the switch would run out of bits to send, called underrun).

**Trap alert:** Count *links* = switches + 1. And remember: cut-through cannot check the FCS before forwarding, so corrupted frames get forwarded and are dropped only at the next store-and-forward point or at the receiver.

**Interview angle:** "What do store-and-forward and cut-through mean for a switch?" Answer with the latency difference and the error-checking trade-off. Follow-up: "Is a router cut-through?" Usually no. Routers typically store-and-forward because they must check the header, decrement TTL, recompute the checksum and often queue the packet.

---

### Q11. [M1 | Store-and-forward vs cut-through trade-offs | Easy | MSQ | Interview]

Select all that apply.

A. Store-and-forward latency per hop does not depend on frame size.
B. A store-and-forward switch can detect a bad FCS and drop the frame before forwarding it.
C. Cut-through switching reduces per-hop latency to roughly the header-reception time, not the full-frame transmission time.
D. A cut-through switch can forward a frame whose FCS turns out to be bad.
E. IP routers must use cut-through because the IP header is at the start of the packet.

**Answer:** B, C, D

**Why:** The FCS (frame check sequence) sits at the **end** of an Ethernet frame. A store-and-forward switch has the whole frame and can verify the FCS before forwarding. A cut-through switch has already started sending by the time the FCS arrives, so it cannot stop a corrupted frame. That is the trade-off: low latency versus filtering errors.

**Working:** Per-hop latency comparison for a 1 Gbps link:

```
Frame size   Store-and-forward per hop   Cut-through per hop (14 B header)
64 B         0.512 us                    0.112 us
1518 B       12.144 us                   0.112 us   <- S&F grows with size; cut-through doesn't
```

**Why the other options fail:**
- A — Store-and-forward per-hop delay = L/R, which is directly proportional to frame size L.
- E — Where the header sits does not force cut-through. Routers usually store-and-forward (TTL decrement, checksum update, queueing).

**Concept recap:** S&F: per-hop delay L/R, error-checks every frame. Cut-through: per-hop delay H/R, forwards errors; fragment-free mode waits for the first 64 B so that collision fragments are not forwarded.

**Trap alert:** Treating "cut-through is faster" as "cut-through is always better". In data centres it is used for low latency, but it forwards corrupted frames and can't be used from a slower link to a faster one.

**Interview angle:** "How do HFT firms or data centres cut latency at L2?" Cut-through switches, kernel-bypass NICs, shorter fibre paths. Follow-up: "Where do corrupted frames get caught in a cut-through fabric?" At the receiving NIC's FCS check, or at the next store-and-forward hop.

---

### Q12. [M1 | Topologies — mesh link count | Easy | MCQ | OA]

A company connects **10 routers in a full mesh** (every router has a dedicated point-to-point link to every other router). How many links are needed, and how many **new** links are needed when an 11th router joins and must also be fully meshed?

A. 45 links; 11 new links
B. 90 links; 20 new links
C. 55 links; 10 new links
D. 45 links; 10 new links

**Answer:** D

**Why:** In a full mesh each pair of nodes needs one link, so the count is "n choose 2". A new node must connect to every existing node: n new links.

**Working:**

```
Links in full mesh = n(n-1)/2
n = 10 : 10*9/2 = 45 links;  ports per router = n - 1 = 9
n = 11 : 11*10/2 = 55 links
New links = 55 - 45 = 10   (= number of existing routers)

Compare (n = 10):
Topology       Links    Single link failure isolates a node?
Full mesh       45       no
Star            10       yes, if the hub fails all are isolated
Ring            10       no for one failure (path goes the other way); two failures can split it
Bus              1 trunk yes (trunk failure)
```

**Why the other options fail:**
- A — 11 new links would count a link to itself.
- B — n(n−1) = 90 counts every link twice (A→B and B→A).
- C — 55 is the count *after* the 11th node joins, not before.

**Concept recap:** Full mesh: n(n−1)/2 links, n−1 ports per node, O(n²) growth. That is why real networks use partial mesh or hierarchical designs, and why logical full meshes (e.g., iBGP full mesh) are replaced by route reflectors.

**Trap alert:** Double counting with n(n−1). Links are undirected pairs, so divide by 2.

**Interview angle:** Asked lightly, usually as "Why don't we fully mesh everything?" Answer: O(n²) cost and ports. Follow-up for SDE roles: "In a microservices system where every service talks to every other, what's the analogous problem?" n² connections/configs, which is why service meshes, connection pooling, or a message broker in a hub-like role are used.

---

### Q13. [M1 | Client-server vs P2P distribution time | Hard | MCQ | OA]

A **1 Gbit (10^9 bits)** file must be distributed from one server to **N = 100** peers.
- Server upload rate us = **100 Mbps**.
- Every peer has upload rate u = **10 Mbps** and download rate d = **20 Mbps**.
- Treat the network core as having unlimited capacity; ignore propagation and protocol overhead. Use the standard minimum distribution-time bounds.

What are the minimum distribution times for (i) client-server and (ii) peer-to-peer (BitTorrent-style)?

A. (i) 1000 s, (ii) 50 s
B. (i) 1000 s, (ii) 90.9 s
C. (i) 50 s, (ii) 90.9 s
D. (i) 1000 s, (ii) 10 s

**Answer:** B

**Why:** In client-server, the server alone must upload N full copies, so its uplink is the bottleneck. In P2P, every peer that downloads also uploads, so total upload capacity grows with N. Three bounds apply, and the distribution time is the largest of them: the server must send at least one copy, the slowest downloader must receive F bits, and the system must push N·F bits in total using the combined upload capacity.

**Working:**

```
F = 10^9 bits, N = 100, us = 10^8 bps, u = 10^7 bps, d = 2*10^7 bps

(i) Client-server:  D_cs = max( N*F/us , F/d_min )
    N*F/us = 100*10^9 / 10^8 = 1000 s
    F/d    = 10^9 / (2*10^7) =   50 s
    D_cs   = 1000 s

(ii) P2P:  D_p2p = max( F/us , F/d_min , N*F/(us + N*u) )
    F/us             = 10^9 / 10^8                  =  10 s
    F/d              = 10^9 / (2*10^7)              =  50 s
    N*F/(us + N*u)   = 10^11 / (10^8 + 100*10^7)
                     = 10^11 / (1.1*10^9)           = 90.9 s
    D_p2p = 90.9 s
```

**Why the other options fail:**
- A — Uses only the download bound for P2P and ignores that the total upload capacity must move N·F bits.
- C — For client-server it picks the download bound (50 s) and ignores that the server must upload 100 copies.
- D — Uses only the server bound F/us for P2P (10 s), as if one copy from the server were enough for everyone.

**Concept recap:**
D_cs ≥ max(N·F/us, F/dmin): grows linearly in N.
D_p2p ≥ max(F/us, F/dmin, N·F/(us + Σui)): levels off as N grows, because capacity scales with demand.

**Trap alert:** Taking the *minimum* of the bounds, or only one bound. A lower bound on time is the **maximum** of all necessary constraints.

**Interview angle:** "Why does BitTorrent scale and a single download server doesn't?" Answer: self-scaling; every downloader adds upload capacity. Follow-up for SDE roles: "How do companies distribute large build artifacts or container images to thousands of machines?" P2P-style distribution (e.g., BitTorrent-like internal tools), CDNs, or hierarchical caches. All of them avoid the single-uplink bottleneck.

---

### Q14. [M1 | Unicast / broadcast / multicast / anycast | Easy | MCQ | OA + Interview]

A public DNS resolver uses the same IP address (e.g., 8.8.8.8) at dozens of data centres worldwide. A user in Mumbai and a user in Frankfurt both send queries to that address, and each query reaches a *different*, topologically nearest site, chosen by normal BGP routing. Which delivery model is this?

A. Multicast
B. Broadcast
C. Anycast
D. Unicast with DNS round-robin

**Answer:** C

**Why:** In anycast, one IP address is advertised from many locations, and routing delivers each packet to **one** member of the group, whichever is "closest" by routing metrics. Each user gets a nearby server with no client-side logic.

**Working:** Delivery-model comparison:

```
Model      Sender -> receivers              Example
Unicast    one -> one specific host         normal TCP connection
Broadcast  one -> all hosts on the subnet   ARP request, DHCPDISCOVER (IPv4 only)
Multicast  one -> all members of a group    IPTV, OSPF hellos (224.0.0.5), IPv6 neighbour discovery
Anycast    one -> nearest one of a group    root DNS servers, public resolvers, CDN edges
```

**Why the other options fail:**
- A — Multicast delivers a copy to *every* group member, not to one nearest member.
- B — Broadcast is limited to a single L2 domain/subnet and is never routed across the Internet.
- D — DNS round-robin returns *different IP addresses* to different clients. Here every client uses the *same* address, and routing does the selection.

**Concept recap:** Anycast = same address at many places + routing picks the nearest. IPv6 has no broadcast; it uses multicast (e.g., ff02::1) instead.

**Trap alert:** Assuming anycast is load balancing. It picks the *routing-nearest* site, which may not be the least loaded. Also, a route change in the middle of a long TCP connection could send packets to a different site, which is why anycast suits mostly short, UDP-based exchanges (DNS), though CDNs do use it for TCP too.

**Interview angle:** "How does a CDN route users to the nearest edge?" Answer: anycast, DNS-based geo-steering, or both. Follow-up: "What's the risk of anycast for long TCP sessions?" A BGP path change can move the flow to another site that has no connection state, which shows up as a connection reset.

---

### Q15. [M2 | Bandwidth vs latency trap | Easy | MCQ | OA + Interview]

Two offices are linked by **3000 km of fibre** (signal speed 2 × 10^8 m/s). The link is upgraded from **100 Mbps to 1 Gbps (10^9 bps)**. Afterwards, `ping` with **100-byte** packets shows almost the same RTT as before (about 30 ms). What best explains this?

A. The upgrade must be misconfigured; a 10× bandwidth increase should cut RTT about 10×.
B. RTT here is dominated by propagation delay (≈15 ms each way), which depends on distance and signal speed. The transmission time of a 100-byte packet was only 8 µs even before the upgrade.
C. A higher link rate raises the signal propagation speed only on copper, not on fibre.
D. ICMP echo packets are rate-limited by routers, so bandwidth upgrades never show up in ping.

**Answer:** B

**Why:** Delay has separate parts. Bandwidth affects only the **transmission** part (pushing bits onto the wire). Propagation delay is distance divided by signal speed and does not change with link rate. For small packets over long distances, propagation dominates.

**Working:**

```
Propagation, one way: d/v = 3*10^6 m / (2*10^8 m/s) = 0.015 s = 15 ms
RTT (propagation only) = 30 ms

Transmission of 100 B = 800 bits:
  at 100 Mbps: 800 / 10^8 = 8 us
  at 1 Gbps:   800 / 10^9 = 0.8 us
Saving per transmission = 7.2 us; per RTT (request + reply) ~ 14.4 us
=> RTT goes from ~30.016 ms to ~30.002 ms : not visible in ping
```

**Why the other options fail:**
- A — This is the core misconception: bandwidth (capacity) is not latency.
- C — Signal speed depends on the medium (≈ 2/3 c in fibre and copper), not on the bit rate.
- D — ICMP rate-limiting exists, but it limits how many replies per second a router sends, not RTT. The fundamental reason here is propagation.

**Concept recap:** Total delay = transmission (L/R) + propagation (d/v) + queuing + processing. Bandwidth only shrinks L/R. Throughput for big transfers can improve 10× even when latency doesn't.

**Trap alert:** "Faster internet = lower latency." A bigger pipe moves *more* data per second; it doesn't make the first bit arrive sooner.

**Interview angle:** "Explain bandwidth vs latency." Use the analogy of a highway's lane count (bandwidth) vs its length (latency). Follow-up: "So what does reduce latency for global users?" Moving servers closer (CDN, edge), fewer round trips (TLS 1.3, HTTP/2 connection reuse, 0-RTT), and avoiding queuing (bufferbloat control).

---

### Q16. [M2 | Throughput vs goodput | Medium | NAT | OA]

A bulk TCP transfer saturates a **100 Mbps (10^8 bps) full-duplex Ethernet** link in one direction. Every data frame carries a full **1460-byte** TCP payload. Per frame, count on the wire: preamble+SFD 8 B, Ethernet header 14 B, IPv4 header 20 B (no options), TCP header 20 B (no options), payload 1460 B, FCS 4 B, inter-frame gap 12 B. ACKs travel in the other direction and do not use this direction's capacity. There are no losses or retransmissions. What is the **goodput** (application bytes per second, expressed in Mbps), to 2 decimal places?

Answer is a number — no options.

**Answer:** ≈ **94.93 Mbps**

**Why:** Goodput counts only useful application data delivered per second. The link moves 100 Mbps of *raw* bits, but each 1460-byte payload occupies a 1538-byte slot on the wire because of headers, the FCS, the preamble and the gap. Goodput = link rate × (useful bytes ÷ bytes on wire).

**Working:**

```
Bytes on wire per frame = 8 + 14 + 20 + 20 + 1460 + 4 + 12 = 1538 B
Useful bytes per frame  = 1460 B

Goodput = 100 Mbps * 1460 / 1538
        = 146,000 / 1538 Mbps
        = 94.928... Mbps ~ 94.93 Mbps

Same link measured at different layers:
  Raw line rate                         100.00 Mbps
  Ethernet frame throughput (1518/1538)  98.70 Mbps
  IP throughput (1500/1538)              97.53 Mbps
  TCP payload = goodput (1460/1538)      94.93 Mbps
```

**Why the other options fail:** (NAT — the common wrong answers)
- 100 Mbps — confuses line rate (bandwidth) with goodput.
- 96.18 Mbps (1460/1518) — forgets the preamble and inter-frame gap, which use wire time even though they aren't part of the frame.
- 97.33 Mbps (1460/1500) — counts only IP+TCP header overhead.

**Concept recap:** Bandwidth = raw capacity. Throughput = actual bits per second achieved at a stated layer. Goodput = application payload per second, after removing headers and retransmissions. Goodput ≤ throughput ≤ bandwidth.

**Trap alert:** Always ask "throughput *measured at which layer*?" A 1 Gbps NIC will never show 1 Gbps in `iperf`. About 941 Mbps of TCP goodput is normal and correct.

**Interview angle:** "My 1 Gbps link only gives 940 Mbps in iperf, is something wrong?" Explain header overhead (≈ 5.1% with 1500 MTU). Follow-up: "How would jumbo frames (9000 MTU) change that?" Overhead drops to about 1%, at the cost of needing every device on the path to support it.

---

### Q17. [M2 | Bottleneck throughput | Medium | MSQ | OA + Interview]

A client downloads a large file from a server. Path: server NIC **10 Gbps** → data-centre uplink → ISP link **200 Mbps** → home router → client Wi-Fi **50 Mbps**. Unless stated otherwise in an option, assume no other traffic. Select all that apply.

A. The long-run throughput of a single transfer is bounded by the bottleneck link, 50 Mbps.
B. Upgrading the server NIC to 25 Gbps would raise this client's throughput.
C. If 10 long-lived flows (including this one) share the 200 Mbps ISP link equally, this flow's throughput is bounded by min(50, 200/10) = 20 Mbps.
D. The latency of a single small packet on this path is set only by the bottleneck bandwidth.
E. Even if every link were 10 Gbps, a single TCP flow could still be limited by window size ÷ RTT.

**Answer:** A, C, E

**Why:** End-to-end throughput is limited by the slowest *effective* share along the path, which is the link rate divided among competing flows. Separately, TCP can't send more than one window per RTT, which caps throughput even on fast links.

**Working:**

```
Effective capacity per hop for this flow:
  server NIC           10,000 Mbps
  ISP link (alone)        200 Mbps    (shared by 10 flows -> 20 Mbps)
  client Wi-Fi             50 Mbps
Bottleneck (alone)  = min(10000, 200, 50) = 50 Mbps
Bottleneck (shared) = min(10000, 20, 50)  = 20 Mbps

Window limit example: window = 64 KB = 524,288 bits, RTT = 50 ms
  max throughput = 524,288 / 0.05 = 10.49 Mbps   (even with 10 Gbps links)
```

**Why the other options fail:**
- B — The server NIC is not the bottleneck. Speeding up a non-bottleneck link doesn't help.
- D — Small-packet latency is mostly propagation + queuing + per-hop transmission. Bandwidth affects only the transmission part.

**Concept recap:** Throughput ≈ min over links of (link capacity ÷ number of flows sharing it), and also ≤ window/RTT for TCP. Find the bottleneck first, then optimize it.

**Trap alert:** Option E is the classic "long fat pipe" situation. A large bandwidth-delay product needs a big window (TCP window scaling); otherwise the extra bandwidth goes unused.

**Interview angle:** "A user on a 1 Gbps line downloads from a server in another continent at only 30 Mbps. Why?" Candidates: bottleneck elsewhere, small TCP window vs high RTT, packet loss triggering congestion control, a shared/congested peering link. Follow-up: "How would you diagnose it?" Use `iperf`/`mtr`, check RTT and loss, inspect the window and retransmissions in a capture, then try parallel streams.

---

### Q18. [M2 | Transmission delay (unit conventions) | Easy | MCQ | OA]

A **2 KB packet** is transmitted on an **8 Mbps** link. Use **1 KB = 1024 bytes** for data size and **1 Mbps = 10^6 bps** for link rate. What is the transmission delay?

A. 2.000 ms
B. 1.953 ms
C. 2.048 ms
D. 0.256 ms

**Answer:** C

**Why:** Transmission delay = number of bits ÷ link rate. The packet is 2048 bytes = 16,384 bits, and the link moves 8,000,000 bits per second. The unit conventions decide the third digit, so they must be applied exactly as stated.

**Working:**

```
L = 2 KB = 2 * 1024 = 2048 bytes = 2048 * 8 = 16,384 bits
R = 8 Mbps = 8 * 10^6 bps
Tt = L / R = 16,384 / 8,000,000 = 0.002048 s = 2.048 ms
```

**Why the other options fail:**
- A — 2.000 ms uses 1 KB = 1000 B (16,000 bits), against the stated convention.
- B — 1.953 ms uses 1 Mbps = 2^20 bps (8,388,608 bps). Link rates are always decimal.
- D — 0.256 ms forgets to convert bytes to bits (2048 / 8×10^6).

**Concept recap:** Tt = L/R. Link rates are decimal (k = 10^3, M = 10^6, G = 10^9). Storage/memory sizes are often binary (KiB = 1024) when stated. Always convert bytes to bits (×8).

**Trap alert:** The ×8 slip is the most common OA error in delay numericals. Write the units next to every number: "bytes" or "bits".

**Interview angle:** Not asked verbally as such. The SDE version is: "How long to send a 1 MB response on a 10 Mbps link?" ≈ 0.8 s of transmission plus RTTs. This drives decisions such as compressing payloads, paginating APIs and using CDNs.

---

### Q19. [M2 | Propagation delay | Easy | MCQ | OA]

Earth stations A and B communicate through a geostationary satellite. Take the satellite altitude as **36,000 km**, assume both stations are directly beneath it (a straight-up and straight-down path), and use signal speed **3 × 10^8 m/s**. What is the one-way propagation delay from A to B?

A. 240 ms
B. 120 ms
C. 480 ms
D. 0.24 ms

**Answer:** A

**Why:** The signal travels up 36,000 km to the satellite and then down 36,000 km to B, a total of 72,000 km at the speed of light. Propagation depends only on distance and medium speed, not on link rate or packet size.

**Working:**

```
Path = up + down = 36,000 + 36,000 = 72,000 km = 7.2 * 10^7 m
Tp = 7.2*10^7 / 3*10^8 = 0.24 s = 240 ms
RTT (A -> B -> A) = 480 ms (ignoring transmission and processing)
```

**Why the other options fail:**
- B — 120 ms counts only one leg (ground to satellite).
- C — 480 ms is the RTT (A→B→A), not one-way.
- D — 0.24 ms is a units slip (km treated as m, or a 1000× error).

**Concept recap:** Tp = d/v. GEO satellite one-way ≈ 240–280 ms (depends on station geometry). LEO constellations (≈ 550 km altitude) cut this to a few ms per leg, which is why LEO broadband has much lower latency.

**Trap alert:** Count every leg of the path. "Via satellite" means up *and* down.

**Interview angle:** "Why are interactive apps painful over GEO satellite links?" Answer: ~500+ ms RTT. Every TCP handshake, TLS handshake and request/response pays it, and TCP slow start ramps up slowly at that RTT. Follow-up: "How do you mitigate it?" Fewer round trips (TLS 1.3, QUIC 0-RTT, connection reuse), performance-enhancing proxies, larger initial windows.

---

### Q20. [M2 | Transmission = propagation crossover | Medium | NAT | GATE]

A point-to-point link is **2 km** long, runs at **100 Mbps (10^8 bps)**, and the signal speed is **2 × 10^8 m/s**. For what packet length **L, in bits**, is the transmission delay exactly equal to the propagation delay?

Answer is a number — no options.

**Answer:** **1000 bits** (= 125 bytes)

**Why:** Set L/R = d/v and solve for L. The result, L = R × Tp, is the number of bits "in flight" on the wire at any instant, i.e., the bandwidth-delay product of the link. Packets shorter than this finish transmitting before their first bit reaches the far end.

**Working:**

```
Tp = d / v = 2000 m / (2*10^8 m/s) = 10^-5 s = 10 us
Set Tt = Tp:  L / R = Tp  ->  L = R * Tp
L = 10^8 bps * 10^-5 s = 1000 bits = 125 bytes

Ratio a = Tp / Tt:
  L = 1000 bits  -> a = 1     (crossover)
  L = 12000 bits -> Tt = 120 us -> a = 0.083 (transmission dominates)
  L = 100 bits   -> Tt = 1 us   -> a = 10    (propagation dominates)
```

**Why the other options fail:** (NAT — the common wrong answers)
- 125 — gives bytes when bits were asked.
- 10 — gives the propagation time in µs instead of the length.
- 100,000 bits — uses d = 2 km as 2 × 10^5 m or a slower medium; a units slip.

**Concept recap:** a = Tp/Tt is the single most important ratio in data-link numericals. Stop-and-wait efficiency = 1/(1 + 2a). CSMA/CD needs Tt ≥ 2Tp. Bandwidth-delay product = R × Tp bits.

**Trap alert:** a is "propagation ÷ transmission", not the reverse. Mixing it up flips every efficiency formula.

**Interview angle:** Mostly a GATE/OA item. The intuition interviewers want: "On a short, slow link the packet is longer than the wire; on a long, fast link the wire holds many packets, so you need pipelining (sliding windows) to use it." The follow-up leads into BDP and TCP window sizing (P06).

---

## Batch 01 summary

**Completed:** 20 / 340   **Remaining:** 320

**Ledger — newly ✓ this batch:**
| Pattern | Count | Questions |
|---|---|---|
| P01 Circuit vs packet vs VC switching delay | ✓ 3 | Q7, Q8, Q9 |
| P02 Store-and-forward vs cut-through | ✓ 2 | Q10, Q11 |
| P03 Transmission delay | ✓ 2 | Q18, Q20 |
| P04 Propagation delay | ✓ 2 | Q19, Q20 |

**Groundwork for later patterns (not counted toward the ledger):** P85 (Nagle/TCP_NODELAY intuition, Q8), P89 (port exhaustion Q5, MTU black hole Q4), P06 (BDP intuition Q17, Q20).

**Batch stats:**
| Metric | This batch | Running total |
|---|---|---|
| Easy / Medium / Hard | 8 / 10 / 2 | 8 / 10 / 2 |
| MCQ / MSQ / NAT | 14 / 4 / 2 | 14 / 4 / 2 |
| OA+Int / OA / Int / GATE | 7 / 8 / 2 / 3 | 7 / 8 / 2 / 3 |
| GATE-only | 0 | 0 |
| MCQ answer key A / B / C / D / E | 4 / 3 / 4 / 3 / 0 | 4 / 3 / 4 / 3 / 0 |
| "All/none of the above" used | 0 | 0 / 5 |

Hard questions are light here on purpose: M1 is foundational. They pick up in M2 (delay decomposition), M3, M4 and M12 to reach 25% overall.

**Traps covered so far:** bandwidth vs throughput vs latency vs goodput (Q15, Q16, Q17); MSS vs MTU (Q4); "TCP guarantees delivery" is false (Q6); TCP demux by 4-tuple vs port (Q5).

**Dedupe index (concepts used):** layering-interface swap · OSI↔TCP/IP mapping · PDU names · 1460/1500/1518 encapsulation · EtherType/Protocol/4-tuple demux · conn-oriented vs connectionless state location · circuit/packet/message timing (3 links, 10^6 bits) · optimal packet size √(Xh/(k−1)) · VC setup + pipelined data · S&F vs cut-through timing (1 Gbps, 3 switches) · cut-through FCS trade-off · mesh n(n−1)/2 · CS vs P2P distribution bounds · anycast · bandwidth≠latency (3000 km ping) · goodput 1460/1538 · bottleneck & window/RTT limit · Tt with KB=1024 · GEO propagation · L = R·Tp crossover.

*Next: Batch 02 (Q21–Q40), covering the rest of M2: total-delay decomposition (P05★), BDP & pipelining, RTT throughput, Nyquist, Shannon, combined, baud vs bit rate, FDM/TDM, duplex modes.*
