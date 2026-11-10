# OS MCQ/MSQ Bank — Batch 13 (Q241–Q260)

**What this batch covers:**
- **M9, I/O and storage completed (Q241–259):**
  - Double buffering; buffering vs caching vs spooling; device drivers.
  - Disk access-time and capacity numericals.
  - FCFS, SSTF, SCAN, LOOK and C-LOOK traces on one request queue.
  - Disk-scheduling properties, why SSDs change the picture, which algorithm gives uniform waits.
  - RAID 0/1/5/6/10: capacity, parity reconstruction, small-write penalty.
  - SSD internals; non-blocking vs asynchronous I/O.
- **M10, protection begins (Q260):** protection vs security.

> **Disk-scheduling conventions for Q247–Q249:**
> - Cylinders are numbered 0–199, and the head starts at **50**.
> - Requests, in arrival order: **82, 170, 43, 140, 24, 16, 190**.
> - "Head movement" means the total cylinders traversed. SCAN goes all the way to the end cylinder (199) before reversing. LOOK and C-LOOK reverse (or jump) at the last request in that direction.
> - **Return jumps in C-SCAN and C-LOOK count as head movement** unless a question says otherwise.

---

## Q241. [M9 | Double buffering | Easy | MCQ | OA]

A process reads data from a slow device and processes it. What does **double buffering** achieve?

- A. Data is written to disk twice, for safety
- B. The device fills one buffer while the process works on the other, then they swap, so I/O overlaps with computation
- C. Two processes share one buffer
- D. Each buffer becomes twice as fast

**Answer:** B

**Why:** With a single buffer, the process must wait while the device fills it, and the device must wait while the process empties it. With **two** buffers, the producer (the device) and the consumer (the process) work at the same time on different buffers. Each step then takes max(I/O time, compute time) instead of their sum.

**Why the other options fail:**
- A: Writing twice for safety is mirroring or journaling.
- C: Sharing one buffer is the bounded-buffer problem, not double buffering.
- D: Buffers don't change speed.

**Concept recap:** Buffering handles three mismatches: **speed** (device vs CPU), **transfer size** (network packets vs file blocks), and **copy semantics** (Q242 E). Double buffering → overlap. A circular buffer generalizes it to n buffers.

**Trap alert:** Double buffering helps only if **both sides can run at the same time**. If the process blocks on the device anyway, there's no overlap.

**Interview angle:** "Why buffer I/O at all?" Name the three mismatches. Follow-up: *"Where else do you see double buffering?"* In graphics: the GPU draws the next frame into a back buffer while the current front buffer is on screen.

---

## Q242. [M9 | Buffering vs caching vs spooling | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. A buffer holds data while it moves between two places, and may hold the **only** copy of that data.
- B. A cache holds a **copy** of data that is stored elsewhere, for faster access.
- C. Spooling queues output for a device that can serve only one job at a time (e.g., a printer).
- D. A buffer and a cache can never be the same piece of memory.
- E. Copy semantics: at `write()`, the kernel copies data out of the user's buffer, so later changes to that buffer don't affect what gets written.

**Answer:** A, B, C, E

**Why:**
- **A vs B:** A buffer is for data **in transit**. A cache keeps a **duplicate** for faster access.
- **C:** **Spooling** (originally "simultaneous peripheral operations on-line") collects each job's output in a separate disk file. Jobs don't interleave, and the device works through the queue one job at a time.
- **E:** Copy semantics make `write()` behave as if the data were written at the moment of the call.

**Why the other options fail:**
- D: The Linux **page cache** serves both roles. A dirty page is buffered output, and it's also a cached copy for future reads.

**Concept recap:** Buffer = data in transit. Cache = a fast copy. Spool = a queue for a device that can't interleave jobs.

**Trap alert:** Copy semantics cost a memory copy. Zero-copy APIs (`sendfile`, `splice`, `io_uring` with registered buffers) avoid it, at the price of stricter rules about when you may reuse the buffer.

**Interview angle:** "Buffer vs cache?" Follow-up: *"Why might a kernel prefer to copy at `write()` rather than DMA straight from user memory?"* Safety and simplicity. The user can't change data that's partway through being written, and the kernel doesn't have to pin the user's pages (Q240 D).

---

## Q243. [M9 | Spooling | Easy | MCQ | OA]

Which device most classically requires **spooling**?

- A. A disk
- B. A keyboard
- C. A printer
- D. A network card

**Answer:** C

**Why:** A printer can't interleave output from two jobs: lines from report A mixed with lines from report B would be garbage. The OS (or a print daemon) spools each job to a disk file and feeds the printer one complete job at a time.

**Why the other options fail:**
- A: A disk handles interleaved block requests naturally.
- B: A keyboard is an input device with one stream of events.
- D: A network card interleaves packets from many flows by design.

**Concept recap:** Spooling turns a device that must be used exclusively into a shared queue. It's also a classic way around deadlock: processes wait on a queue rather than holding the device.

**Trap alert:** Spooling is a **queue for a device**. Don't confuse it with ordinary buffering for speed.

**Interview angle:** Follow-up: *"Which deadlock condition does spooling help avoid for printers?"* Mutual exclusion, from the processes' point of view: only the spooler daemon ever owns the printer.

---

## Q244. [M9 | Device drivers | Easy | MCQ | Interview]

Which description of a **device driver** is correct?

- A. A device driver translates the kernel's generic I/O requests into the specific commands, register accesses and interrupt handling for one kind of controller. On Linux, most drivers run in kernel mode
- B. A user-space application shipped by the vendor, which runs only at boot
- C. Firmware that runs inside the device itself
- D. The part of the kernel that schedules CPU time

**Answer:** A

**Why:** The kernel's I/O subsystem offers a **uniform interface**: block devices (read/write blocks), character devices (byte streams), and network interfaces. A driver adapts that interface to one particular piece of hardware, so the rest of the kernel doesn't need to know any device's details.

**Why the other options fail:**
- B: Drivers handle I/O the whole time the system runs, not just at boot. (User-space driver frameworks exist, but they're the exception.)
- C: Firmware runs *on* the device. The driver runs on the CPU and talks to that firmware.
- D: That's the CPU scheduler.

**Concept recap:** Layers: application → system call → VFS / block layer → driver → controller → device.

**Trap alert:** Because drivers run in kernel mode, a buggy driver can crash the whole system. Driver bugs are a major cause of OS crashes, which is part of the microkernel argument (Q12).

**Interview angle:** "What happens between `write()` and the disk?" Walk through the layers. Follow-up: *"What is a loadable kernel module?"* A driver loaded at runtime into kernel space (Q13).

---

## Q245. [M9 | Disk access time | Medium | Numerical | OA]

An HDD has an average seek time of **4 ms**, spins at **7200 RPM**, and transfers **100 MB/s** (MB = 10⁶ bytes). Ignore controller overhead. What is the average time to read one random **4096-byte** block?

- A. ≈ 12.37 ms
- B. ≈ 4.04 ms
- C. ≈ 4.21 ms
- D. ≈ 8.21 ms

**Answer:** D

**Why:** Access time = **seek + rotational latency + transfer**. On average, the head waits **half** a rotation for the right sector to come round.

**Working:**
1. Rotation period = 60 s / 7200 = **8.33 ms**. Average rotational latency = 8.33 / 2 = **4.17 ms**.
2. Transfer = 4096 B / (100 × 10⁶ B/s) = **0.041 ms**.
3. Total = 4 + 4.17 + 0.04 ≈ **8.21 ms**.

**Why the other options fail:**
- A: Uses a **full** rotation (8.33 ms) as the latency.
- B: Forgets rotational latency entirely.
- C: Forgets the seek.

**Concept recap:**
- T_access = T_seek + ½ × (60 / RPM) + (block size / transfer rate).
- For small random reads, the **mechanical** parts dominate. Transfer is almost negligible.

**Trap alert:** Use **half** a rotation for the *average* latency. Watch the units: decimal MB for transfer rates, binary units for block sizes.

**Interview angle:** "Why are random reads on HDDs so slow?" About 8 ms per read means only about 120 random IOPS. Follow-up: *"How do databases cope?"* Caching, sequential access patterns (logs, LSM trees), and SSDs.

---

## Q246. [M9 | Disk capacity | Easy | MCQ | OA]

A disk has **8** recording surfaces, **10,000** tracks per surface, **500** sectors per track, and **512** bytes per sector. What is its capacity?

- A. 20.48 × 10⁹ bytes (≈ 19.1 GiB)
- B. 2.56 × 10⁹ bytes
- C. 40.96 × 10⁹ bytes
- D. 20.48 × 10⁶ bytes

**Answer:** A

**Why:** Capacity = surfaces × tracks per surface × sectors per track × bytes per sector
= 8 × 10,000 × 500 × 512
= 80,000 tracks × 500 = 40,000,000 sectors × 512 = **20,480,000,000 bytes** ≈ 20.48 GB (decimal) ≈ 19.07 GiB.

**Why the other options fail:**
- B: Forgets the factor of 8 surfaces.
- C: Doubles the count.
- D: Off by a factor of 1000 (MB for GB).

**Concept recap:** A **cylinder** is the set of tracks at the same arm position, one on each surface. The number of cylinders equals the number of tracks per surface.

**Trap alert:** Drive makers quote decimal GB, while OSes often show binary GiB. That's why a "20 GB" disk shows up as about 19 GiB.

**Interview angle:** Rarely asked verbally. Follow-up: *"Why do modern disks have more sectors per track on outer tracks?"* Zoned bit recording: outer tracks are longer, so they can hold more sectors at the same bit density.

---

## Q247. [M9 | FCFS and SSTF disk scheduling | Medium | Numerical | OA]

Use the shared setup (head at 50; requests 82, 170, 43, 140, 24, 16, 190). What is the total head movement under **FCFS** and under **SSTF**?

- A. 600 and 208
- B. 642 and 208
- C. 642 and 226
- D. 208 and 642

**Answer:** B

**Why:**
- **FCFS** serves requests in arrival order.
- **SSTF** (Shortest Seek Time First) always moves to the closest pending request.

**Working: FCFS**

| Move | 50→82 | 82→170 | 170→43 | 43→140 | 140→24 | 24→16 | 16→190 |
|---|---|---|---|---|---|---|---|
| Cylinders | 32 | 88 | 127 | 97 | 116 | 8 | 174 |

Total = 32 + 88 + 127 + 97 + 116 + 8 + 174 = **642**.

**Working: SSTF**

| From | Closest pending request | Distance |
|---|---|---|
| 50 | 43 | 7 |
| 43 | 24 | 19 |
| 24 | 16 | 8 |
| 16 | 82 | 66 |
| 82 | 140 | 58 |
| 140 | 170 | 30 |
| 170 | 190 | 20 |

Total = 7 + 19 + 8 + 66 + 58 + 30 + 20 = **208**.

**Why the other options fail:**
- A: The FCFS total is miscounted.
- C: The SSTF total is miscounted (for example by going to 82 first).
- D: Swapped.

**Concept recap:** FCFS is fair but can zig-zag wildly. SSTF is greedy: low total movement, but it can **starve** far-away requests (Q250).

**Trap alert:** SSTF is **not optimal**. Being greedy at each step doesn't guarantee the minimum total movement.

**Interview angle:** "What's the downside of SSTF?" Starvation at the edges, while a steady stream of requests arrives near the head. Follow-up: *"How do elevator algorithms fix that?"* They sweep in a direction, so every request is eventually reached (Q248).

---

## Q248. [M9 | SCAN and LOOK | Hard | Numerical | OA + Interview]

Shared setup, with the head at 50 and **moving toward higher cylinder numbers**. What is the total head movement under **SCAN** and under **LOOK**?

- A. 332 and 314
- B. 240 and 314
- C. 314 and 332
- D. 332 and 340

**Answer:** A

**Why:** Both are "elevator" algorithms: serve requests in the current direction, then reverse.
- **SCAN** travels all the way to the end of the disk before reversing.
- **LOOK** reverses as soon as there are no more requests in the current direction.

**Working: SCAN** (moving up)
- 50 → 82 → 140 → 170 → 190 → **199** (the end) = 199 − 50 = **149**.
- Reverse: 199 → 43 → 24 → 16 = 199 − 16 = **183**.
- Total = 149 + 183 = **332**.

**Working: LOOK** (moving up)
- 50 → 82 → 140 → 170 → **190** (the last request upward) = 190 − 50 = **140**.
- Reverse: 190 → 43 → 24 → 16 = 190 − 16 = **174**.
- Total = 140 + 174 = **314**.

**Why the other options fail:**
- B: SCAN run in the **wrong direction**: 50 → 0 (50), then 0 → 190 (190), giving 240.
- C: Swapped.
- D: The LOOK total is miscounted.

**Concept recap:**
- SCAN = (end − head) + (end − lowest request), when moving up first.
- LOOK = (highest request − head) + (highest request − lowest request).

**Trap alert:** Read the **initial direction** carefully, and whether the algorithm is SCAN (goes to the end) or LOOK (turns at the last request). These two details decide the answer.

**Interview angle:** "SCAN vs LOOK?" LOOK avoids pointless travel to the edge of the disk. Follow-up: *"Why is SCAN called the elevator algorithm?"* It serves requests along one direction like an elevator serving floors, then reverses.

---

## Q249. [M9 | C-LOOK | Hard | Numerical | OA]

Shared setup, with the head at 50 moving toward higher cylinders. Under **C-LOOK**, after serving the highest request the head **jumps** back to the lowest pending request and continues upward. Counting that jump as head movement, what is the total?

- A. 167
- B. 391
- C. 341
- D. 314

**Answer:** C

**Why:** **C-LOOK** serves requests in **one direction only**. After the last request upward, it returns to the lowest pending request and sweeps up again. This treats the cylinders as a circular list.

**Working:**
1. Upward: 50 → 82 → 140 → 170 → 190 = **140**.
2. Jump: 190 → 16 = **174** (counted, per the convention).
3. Upward again: 16 → 24 → 43 = **27**.
4. Total = 140 + 174 + 27 = **341**.

**Why the other options fail:**
- A: Doesn't count the jump (140 + 27).
- B: **C-SCAN** counting the full trip: up to 199 (149), back to 0 (199), then up to 43 (43).
- D: That's LOOK (Q248), which serves requests on the way back down.

**Concept recap:**
- C-SCAN = (end − head) + (end − 0) + (highest request below the starting position − 0).
- C-LOOK = (highest − head) + (highest − lowest) + (last request served upward − lowest).
- Whether to count the jump is a **convention**, so always state it.

**Trap alert:** C-algorithms **don't serve requests on the return trip**. LOOK and SCAN do.

**Interview angle:** "Why use C-SCAN or C-LOOK instead of SCAN or LOOK?" They give more **uniform** waiting times (Q252). Follow-up: *"Is the return trip wasted?"* Partly. But it's a single fast seek with no stops, so it's cheap compared with a sweep that stops at every request.

---

## Q250. [M9 | Disk-scheduling properties | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. SSTF can starve requests far from the head when new requests keep arriving near it.
- B. SCAN gives more uniform waiting times than SSTF.
- C. C-SCAN treats the cylinders as circular, which gives more uniform waiting times than SCAN.
- D. LOOK reverses direction at the last request, rather than at the end of the disk.
- E. SSTF always minimizes the total head movement.

**Answer:** A, B, C, D

**Why:**
- **A:** SSTF is greedy, so close requests keep winning.
- **B:** SCAN guarantees every cylinder is visited on each sweep, so no request waits forever.
- **C:** Under SCAN, a request at the edge just behind the head waits about two full sweeps, while one in the middle waits less. C-SCAN evens this out.
- **D:** That's exactly the difference between LOOK and SCAN (Q248).

**Why the other options fail:**
- E: A greedy choice at each step can lead to a longer total path. The true minimum needs global planning, and the queue keeps changing anyway.

**Concept recap:**

| Algorithm | Movement | Starvation risk | Waiting-time uniformity |
|---|---|---|---|
| FCFS | high | none | fair but slow |
| SSTF | low | **yes** | poor |
| SCAN / LOOK | low–medium | none | medium (middle favoured) |
| C-SCAN / C-LOOK | medium | none | **best** |

**Trap alert:** "Lowest total movement" ≠ "best response time for every request".

**Interview angle:** "Which disk scheduler would you choose for a busy server?" C-LOOK or a deadline-based scheduler. Follow-up: *"What does Linux's mq-deadline add?"* Expiry times on requests, so none waits too long, plus a preference for reads.

---

## Q251. [M9 | Disk scheduling for SSDs | Medium | MCQ | Interview]

Why are seek-optimizing algorithms like SCAN largely irrelevant for NVMe SSDs?

- A. SSDs allow only sequential access
- B. SSDs have higher seek times than HDDs
- C. SSDs never receive more than one request at a time
- D. An SSD has no moving head, so access time barely depends on the "distance" between addresses. Linux typically uses a minimal I/O scheduler (e.g., `none` or `mq-deadline`), and the device's many parallel internal queues handle the ordering

**Answer:** D

**Why:** Disk-scheduling algorithms exist to minimize **mechanical** cost: arm movement and rotation. Flash storage has neither, so reordering requests by address saves almost nothing. What matters for SSDs is **parallelism** (many queues spread across many flash chips) and fairness and latency between workloads.

**Why the other options fail:**
- A: SSDs are excellent at random access.
- B: SSDs have essentially no seek time.
- C: NVMe supports many queues, each with deep queues of outstanding requests.

**Concept recap:** On an HDD, the scheduler reduces seeks. On an SSD, the scheduler mainly passes requests through, perhaps adding fairness or latency control.

**Trap alert:** Sequential I/O is still somewhat faster than random I/O on SSDs, because of internal parallelism, read-ahead and the flash translation layer. But the gap is small compared with HDDs.

**Interview angle:** "Does disk scheduling matter on SSDs?" Much less. Follow-up: *"What matters instead?"* Queue depth, parallelism, write amplification (Q258), and the latency of each request.

---

## Q252. [M9 | Uniform waiting time under load | Medium | MCQ | OA]

Under heavy load, which algorithm gives the most **uniform** waiting time across all cylinders?

- A. FCFS
- B. SSTF
- C. C-SCAN
- D. SCAN

**Answer:** C

**Why:** Under SCAN, a cylinder in the middle is visited twice per round trip, at roughly evenly spaced times. A cylinder at an edge is visited twice in quick succession and then waits about two whole sweeps. C-SCAN serves requests in only one direction, so **every cylinder waits about one full cycle**, wherever it is.

**Why the other options fail:**
- A: FCFS ordering is fair, but the long zig-zag movements make waits long and highly variable.
- B: SSTF can starve the edges.
- D: SCAN favours the middle cylinders.

**Concept recap:** C-SCAN/C-LOOK → uniform waits. SCAN → middle-biased. SSTF → biased toward where the head currently is.

**Trap alert:** "SCAN is fair because it visits everyone" is only partly true. *How often* each cylinder is visited still depends on where it is.

**Interview angle:** Follow-up: *"Why not just always use C-LOOK?"* Under light load, SSTF or plain LOOK may give a lower average latency, and modern schedulers are deadline- or budget-based anyway.

---

## Q253. [M9 | RAID 0 | Easy | MCQ | OA]

Data is striped with **RAID 0** across **4** disks, and one disk fails. What is the result?

- A. The array keeps working at reduced speed
- B. The array's data is lost: RAID 0 has no redundancy, and files larger than one stripe have pieces on every disk
- C. Exactly 25% of the files are affected, and the rest are fine
- D. The array rebuilds itself from parity

**Answer:** B

**Why:** **RAID 0** (striping) spreads consecutive chunks across disks for speed and full capacity, with **zero redundancy**. Losing any disk breaks almost every file larger than one stripe, so in practice the whole array is lost. The chance that the array fails also *grows* with the number of disks.

**Why the other options fail:**
- A: There's nothing to fall back on.
- C: Files are split across disks, not assigned one disk each.
- D: RAID 0 has no parity.

**Concept recap:** RAID 0: capacity = n × disk size, best performance, no fault tolerance. P(array fails) ≈ n × P(one disk fails), for small probabilities.

**Trap alert:** "RAID = redundancy" isn't true for RAID 0. The "R" in the name doesn't apply.

**Interview angle:** "When would you use RAID 0?" For scratch or temporary data where speed matters and losing it is acceptable. Follow-up: *"RAID 0 vs RAID 1?"* Speed vs safety.

---

## Q254. [M9 | RAID capacity | Medium | MCQ | OA + Interview]

You have **6 disks of 4 TB** each. What usable capacity do **RAID 5** and **RAID 6** give, respectively?

- A. 24 TB and 20 TB
- B. 20 TB and 12 TB
- C. 16 TB and 20 TB
- D. 20 TB and 16 TB

**Answer:** D

**Why:** RAID 5 uses one disk's worth of space for **parity**, and RAID 6 uses two.
- RAID 5: (n − 1) × size = 5 × 4 = **20 TB**. It survives any 1 disk failure.
- RAID 6: (n − 2) × size = 4 × 4 = **16 TB**. It survives any 2 disk failures.

**Why the other options fail:**
- A: RAID 5 is given RAID 0's capacity (no parity).
- B: The second number (12 TB) is the RAID 10 capacity (n/2 × size), not RAID 6.
- C: The two values are swapped.

**Concept recap:**

| Level | Usable capacity | Failures survived |
|---|---|---|
| RAID 0 | n·S | 0 |
| RAID 1 (2-way) | S per mirrored pair | 1 per pair |
| RAID 5 | (n − 1)·S | 1 |
| RAID 6 | (n − 2)·S | 2 |
| RAID 10 | (n/2)·S | 1 per mirrored pair |

**Trap alert:** Parity in RAID 5 is **distributed** across all disks. It still costs exactly one disk's worth of capacity.

**Interview angle:** "Why prefer RAID 6 over RAID 5 for large disks?" Rebuilding a big disk takes hours, and a second failure or an unreadable sector during the rebuild would lose the RAID 5 array. Follow-up: *"Why is RAID 10 popular for databases?"* Small writes have no parity penalty (Q256), and rebuilds are fast.

---

## Q255. [M9 | RAID 5 parity reconstruction | Hard | MCQ | OA]

In one RAID 5 stripe the data blocks are **D1 = 1011**, **D2 = 0110** and **D3 = 1100**, and the parity block is **P = D1 ⊕ D2 ⊕ D3**. The disk holding D2 fails. What is the reconstructed D2?

- A. 0110
- B. 0001
- C. 1101
- D. 0111

**Answer:** A

**Why:** XOR parity lets you rebuild **any single missing block** by XOR-ing everything that survives, because x ⊕ x = 0.
- First compute P: 1011 ⊕ 0110 = 1101; then 1101 ⊕ 1100 = **0001**.
- Rebuild D2 = D1 ⊕ D3 ⊕ P: 1011 ⊕ 1100 = 0111; then 0111 ⊕ 0001 = **0110** ✓. This matches the original D2.

**Why the other options fail:**
- B: That's the parity block P itself.
- C: An intermediate value (D1 ⊕ D2).
- D: D1 ⊕ D3, forgetting to include P.

**Concept recap:** P = ⊕(all data blocks); a missing block = ⊕(all surviving blocks, including P). Single parity ⇒ any **one** missing block can be recovered.

**Trap alert:** Two missing blocks can't be recovered from single parity. That's why RAID 6 adds a second, independent parity (a Reed–Solomon "Q" block).

**Interview angle:** "How does RAID 5 recover from a disk failure?" XOR the survivors. Follow-up: *"What happens to reads while the array is degraded?"* Each read of the failed disk's data needs reads from **every** surviving disk, so performance drops sharply.

---

## Q256. [M9 | RAID 5 small-write penalty | Hard | MCQ | OA + Interview]

A single data block in a RAID 5 stripe is updated using **read-modify-write**. How many disk I/Os does this take?

- A. 1
- B. 2
- C. 4: read the old data, read the old parity, write the new data, write the new parity (new P = old P ⊕ old D ⊕ new D)
- D. 6

**Answer:** C

**Why:** The parity must stay correct. Instead of reading the whole stripe, the controller uses the XOR identity:
- **P_new = P_old ⊕ D_old ⊕ D_new**
- That needs 2 reads (D_old, P_old) + 2 writes (D_new, P_new) = **4 I/Os** for one logical write.

**Why the other options fail:**
- A: That's a plain write with no redundancy (RAID 0).
- B: That's RAID 1, which writes both mirrors.
- D: That's **RAID 6**: read the data and both parities, then write the data and both parities.

**Concept recap:** Small-write cost per logical write: RAID 0 → 1; RAID 1/10 → 2; RAID 5 → 4; RAID 6 → 6. A full-stripe write avoids the reads altogether.

**Trap alert:** RAID 5's write penalty matters most for **small random writes**, which is exactly a database's OLTP pattern.

**Interview angle:** "Why do write-heavy databases avoid RAID 5?" The 4× I/O penalty, and the "write hole" (Q257 trap). Follow-up: *"How do controllers hide the penalty?"* A battery-backed write cache, and gathering writes into full stripes.

---

## Q257. [M9 | RAID levels | Medium | MSQ | OA]

**Select all that apply.**

- A. RAID 1 can serve a read from either mirror, which improves read throughput.
- B. RAID 5 spreads parity across all disks, avoiding the dedicated parity-disk bottleneck of RAID 4.
- C. RAID 6 survives any two disk failures.
- D. RAID replaces the need for backups.
- E. RAID 10 stripes data across mirrored pairs.

**Answer:** A, B, C, E

**Why:**
- **A:** Two copies give two places to read from.
- **B:** In RAID 4, every write updates the single parity disk, which becomes a bottleneck. RAID 5 rotates the parity across all disks.
- **C:** Two independent parities let RAID 6 recover any two missing blocks.
- **E:** RAID 10 = RAID 1 + 0: data is striped across mirrored pairs.

**Why the other options fail:**
- D: RAID protects against **disk failure** only. It doesn't protect against deletion, corruption, ransomware, or fire, because every one of those is faithfully copied to the redundancy too. You still need backups.

**Concept recap:** RAID gives **availability**. Backups give **recoverability**.

**Trap alert:** The **RAID 5 "write hole"**: if power fails between writing the data and writing the parity, they disagree, and silent corruption can follow. Battery-backed caches, journals, or ZFS-style copy-on-write close this hole.

**Interview angle:** "Is RAID a backup?" No, and say why. Follow-up: *"RAID 4 vs RAID 5?"* A dedicated parity disk vs parity spread across all disks.

---

## Q258. [M9 | SSD vs HDD internals | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. SSDs have no seek time and no rotational latency.
- B. Flash pages can't be overwritten in place, and erasing happens in larger blocks. So SSDs need a flash translation layer (FTL) and garbage collection.
- C. TRIM lets the OS tell the SSD which blocks are no longer in use.
- D. Flash cells survive only a limited number of program/erase cycles, so wear leveling spreads writes evenly.
- E. On any SSD, small random writes are faster than sequential writes.

**Answer:** A, B, C, D

**Why:**
- **A:** There are no moving parts.
- **B:** Writes go to fresh pages, and the FTL remaps logical addresses to the new physical pages. Garbage collection later erases blocks that contain only stale pages. This causes **write amplification**: the drive writes more physical data than the host asked it to.
- **C:** Without TRIM, the SSD thinks deleted data is still live, and keeps copying it during garbage collection.
- **D:** Wear leveling keeps a few heavily used blocks from wearing out early.

**Why the other options fail:**
- E: Small random writes increase garbage collection and write amplification. They're usually **slower** than sequential writes.

**Concept recap:**
- HDD: mechanical latency dominates.
- SSD: latency in microseconds, with the FTL, garbage collection, wear leveling and TRIM.
- Both still benefit from caching and batching I/O.

**Trap alert:** An SSD's performance can **drop** when it's nearly full or after heavy random writes, because garbage collection has less free space to work with.

**Interview angle:** "How is an SSD different from a disk, from the OS's point of view?" Follow-up: *"Why do log-structured designs suit flash?"* They write sequentially and never overwrite in place, which matches how flash works.

---

## Q259. [M9 | Non-blocking vs asynchronous I/O | Medium | MCQ | OA + Interview]

Which statement correctly distinguishes **non-blocking** I/O from **asynchronous** I/O?

- A. They are the same thing
- B. A non-blocking read returns immediately with whatever data is available (possibly none, with `EAGAIN`), and the caller must try again later. Asynchronous I/O starts the whole operation, and the kernel notifies the caller when it has **completed** (e.g., `io_uring`, POSIX AIO)
- C. Asynchronous I/O blocks until the operation is done
- D. Non-blocking I/O requires DMA

**Answer:** B

**Why:**
- **Blocking:** the call returns when the operation is done; the thread sleeps meanwhile.
- **Non-blocking:** the call returns at once with partial data or "try again". You usually pair it with **readiness notification** (`select`, `poll`, `epoll`) that says *when* to retry.
- **Asynchronous:** you submit a request and get notified when it's **finished**. The data is already in your buffer when the notification arrives.

**Why the other options fail:**
- A: Readiness and completion are different models.
- C: That describes blocking I/O.
- D: DMA is about how data moves inside the hardware, not about this API.

**Concept recap:** epoll = a **readiness** model ("you can read now"). io_uring and IOCP = a **completion** model ("your read is done").

**Trap alert:** `epoll` with non-blocking sockets is **not** asynchronous I/O in this strict sense. You still perform the `read` yourself after being told data is ready.

**Interview angle:** "How do high-performance servers handle 10,000 connections?" With event loops using non-blocking I/O and epoll, or with completion-based I/O (io_uring). Follow-up: *"Why was io_uring added to Linux?"* Batched submissions, fewer system calls, and true completion-based I/O for files as well as sockets.

---

## Q260. [M10 | Protection vs security | Easy | MCQ | Interview]

Which distinction is correct?

- A. Protection and security are synonyms
- B. Security controls which processes can access which resources inside the system; protection is about outside attackers
- C. Protection applies only to files
- D. Protection is the internal mechanism that controls which processes and users can access which system resources. Security is the broader defence of the whole system, including protection, against internal and external threats such as malware, intrusion and data theft

**Answer:** D

**Why:**
- **Protection** = *mechanisms*: memory protection (the MMU), file permissions, privileged instructions, access matrices.
- **Security** = *policies and defences* built on those mechanisms: authentication, encryption, auditing, and handling human factors and attacks.

Protection is necessary but not sufficient. Perfect permission checks don't help if an attacker steals a password.

**Why the other options fail:**
- A: The two are related but have different scope.
- B: The definitions are swapped.
- C: Protection also covers memory, the CPU and devices.

**Concept recap:** Protection answers "*who may do what* to which resource?". Security answers "can the system *resist* abuse, even by someone with valid credentials?".

**Trap alert:** "Separate mechanism from policy" is the guiding principle. The OS provides the mechanism (for example permission bits); administrators and users set the policy (which bits to use).

**Interview angle:** "Protection vs security?" Follow-up: *"Give an example where protection works but security fails."* Phishing a user's password: every access check passes, but the person using the account is the attacker.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 241 | B | 246 | A | 251 | D | 256 | C |
| 242 | A, B, C, E | 247 | B | 252 | C | 257 | A, B, C, E |
| 243 | C | 248 | A | 253 | B | 258 | A, B, C, D |
| 244 | A | 249 | C | 254 | D | 259 | B |
| 245 | D | 250 | A, B, C, D | 255 | A | 260 | D |

---

**Concepts tested so far (to avoid repeats):** Batches 1–12, plus: double buffering and the three buffering motives · buffer vs cache vs spool, copy semantics, zero-copy · spooling for printers · driver layering · access time with average latency = half a rotation · capacity from geometry, GB vs GiB · FCFS 642 vs SSTF 208 · SCAN 332 vs LOOK 314 with direction · C-LOOK 341 counting the jump, C-SCAN 391 · scheduling properties table · SSDs and the `none`/`mq-deadline` schedulers · C-SCAN uniformity · RAID 0 has no redundancy · RAID 5/6/10 capacity · XOR reconstruction · small-write penalty 4/6 · RAID 4 vs 5, RAID isn't a backup, write hole · FTL, garbage collection, TRIM, wear leveling, write amplification · readiness vs completion I/O · protection vs security
*(M9 is complete. Batch 14 continues M10 with authentication vs authorization, the access matrix, ACLs vs capabilities, protection domains, least privilege, isolation and privilege escalation, then starts M11 (Linux internals).)*

**Completed: 260 / 300    Remaining: 40**
