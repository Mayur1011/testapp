# OS MCQ/MSQ Bank — Batch 9 (Q161–Q180)

**What this batch covers:**
- **M6, memory management completed (Q161–169):** multilevel page-table sizing, memory accesses per page walk, inverted page tables, TLB-based EMAT and the hit ratio needed for a target, TLB properties, shared pages, buddy allocation.
- **M7, virtual memory begins (Q170–180):** pure demand paging, the steps of handling a page fault, valid/invalid bits, EMAT with page faults, the fault rate needed for a slowdown target, FIFO replacement, Belady's anomaly, restarting faulted instructions, minor vs major faults.

> **Conventions:**
> - KB = 2¹⁰ bytes, MB = 2²⁰ bytes. 1 ms = 10⁶ ns.
> - Unless stated, a TLB lookup happens **before** the page-table walk, so its time is paid on hits and misses.
> - Page-replacement traces start with all frames empty, and every first reference to a page is a fault.

---

## Q161. [M6 | Multilevel page tables: memory used | Hard | Numerical | OA]

- 32-bit virtual addresses, 4 KB pages, 4-byte PTEs.
- **Two-level** paging, with the address split **10 | 10 | 12** (outer index | inner index | offset). The outer table fits exactly in one page.
- The process uses exactly two regions:
  - the **lowest 4 MB** of its address space (code, data, heap): 0x00000000–0x003FFFFF
  - the **highest 4 MB** (stack): 0xFFC00000–0xFFFFFFFF

How much memory do its page tables occupy?

- A. 4 MB
- B. 12 KB
- C. 8 KB
- D. 8 MB

**Answer:** B

**Why:** With multiple levels, **inner tables exist only for regions that are actually used**. A single-level table needs an entry for every possible page, even the unused ones.

**Working:**
- **Single-level, for comparison:** 2²⁰ entries × 4 B = **4 MB**.
- **Two-level:**
  1. Outer table: 2¹⁰ entries × 4 B = **4 KB**, which is always present.
  2. Each inner table has 2¹⁰ entries mapping 4 KB pages, so it covers 2¹⁰ × 4 KB = **4 MB** of address space.
  3. The bottom 4 MB is exactly one aligned 4 MB chunk → 1 inner table. The top 4 MB is another → 1 inner table.
  4. Inner tables: 2 × 4 KB = **8 KB**.
  5. Total = 4 KB + 8 KB = **12 KB**.

**Why the other options fail:**
- A: The single-level size. It ignores that unused inner tables are never created.
- C: Forgets the outer table.
- D: Mixes up table size with the amount of memory mapped.

**Concept recap:** Multilevel table memory = outer table + (number of inner tables actually needed) × inner table size. Each inner table covers (entries per table) × (page size) of the address space.

**Trap alert:** Count inner tables by **aligned chunks touched**. A 4 MB region that straddles a chunk boundary would need 2 inner tables.

**Interview angle:** "Why do we use multilevel page tables?" Address spaces are **sparse**, so most inner tables are never allocated. Follow-up: *"What's the cost?"* More memory accesses on a TLB miss (Q162).

---

## Q162. [M6 | Cost of a page walk | Medium | MCQ | OA + Interview]

x86-64 uses **4-level** page tables. Suppose there is a TLB miss and none of the page-table entries are cached anywhere. How many memory accesses does a single load instruction need, counting the data access?

- A. 1
- B. 2
- C. 4
- D. 5

**Answer:** D

**Why:** Translation reads one entry at each level, from the top-level table (PML4) down to the last-level table that holds the frame number. That's **4 accesses**. Then the data itself is read, for **1 more**. Total = 4 + 1 = **5**.

**Why the other options fail:**
- A: Assumes a TLB hit.
- B: The single-level case (Q151 E).
- C: Forgets the data access.

**Concept recap:** On a TLB miss with an n-level page table, one load costs **n + 1** memory accesses. That's why TLB hit ratios (and page-walk caches) matter so much.

**Trap alert:** "4 levels ⇒ 4 accesses" is wrong. Always add the access for the data itself.

**Interview angle:** "Why is a TLB miss expensive on x86-64?" Up to 4 extra memory reads for the walk. Follow-up: *"How do huge pages help?"* A 2 MB page skips the last level, and one TLB entry covers 512× more memory.

---

## Q163. [M6 | Inverted page tables | Medium | MCQ | Interview]

Which statement about **inverted page tables** is correct?

- A. There is one entry per **physical frame**, holding ⟨process id, virtual page number⟩. The table is far smaller than per-process tables, but a lookup must search it (usually through a hash table), and sharing a frame between processes is awkward
- B. There is one entry per virtual page, per process
- C. It removes the need for a TLB
- D. It makes shared memory especially easy

**Answer:** A

**Why:** Normal page tables grow with the **virtual** address space, and there's one per process. An inverted table grows with **physical** memory, and there's **one for the whole system**. The catch: given a ⟨pid, page⟩ pair, you can't index the table directly, so you must search it, which in practice means hashing.

**Why the other options fail:**
- B: That describes ordinary per-process page tables.
- C: The lookup is slower, so a TLB matters *more*.
- D: Each frame has exactly one ⟨pid, page⟩ owner, so mapping one frame into several processes is hard.

**Concept recap:** Inverted table size = (number of frames) × entry size. It saves memory but makes lookups and sharing harder. PowerPC and IA-64 used hashed or inverted designs.

**Trap alert:** "Inverted = faster" is wrong. It's **smaller**, not faster.

**Interview angle:** "Hierarchical vs hashed vs inverted page tables?" Follow-up: *"Why do inverted tables make shared memory hard?"* A frame can record only one ⟨pid, page⟩ owner.

---

## Q164. [M6 | EMAT with a TLB | Medium | Numerical | OA + Interview]

- TLB lookup: **10 ns**. Memory access: **100 ns**. TLB hit ratio: **90%**.
- Single-level page table. The TLB is checked first. On a miss, the page table is read from memory, and then the data.

What is the effective memory access time (EMAT)?

- A. 110 ns
- B. 119 ns
- C. 120 ns
- D. 200 ns

**Answer:** C

**Why:**
- On a **hit**: TLB lookup + data access.
- On a **miss**: TLB lookup + page-table read + data access.

Weight each case by its probability.

**Working:**
- **Formula:** EMAT = h·(t + m) + (1 − h)·(t + 2m), where t = TLB time and m = memory time.
- Hit cost = 10 + 100 = 110 ns.
- Miss cost = 10 + 100 + 100 = 210 ns.
- EMAT = 0.9 × 110 + 0.1 × 210 = 99 + 21 = **120 ns**.

**Why the other options fail:**
- A: Assumes every access is a hit.
- B: Leaves the TLB time out of the miss case (0.9 × 110 + 0.1 × 200).
- D: Assumes there's no TLB at all (2 × 100).

**Concept recap:** EMAT = h(t + m) + (1 − h)(t + (n + 1)m) for an n-level table. When n = 1, the miss term is t + 2m.

**Trap alert:** Read whether the TLB lookup overlaps with the page-table access (parallel) or comes first (serial). Exams use both conventions, and they give different answers (119 vs 120 here).

**Interview angle:** "What is a TLB and why does it matter?" Follow-up: *"What hit ratio do real systems get?"* Typically above 99% for most workloads, thanks to locality.

---

## Q165. [M6 | TLB properties | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. A TLB is a small, fast, associative cache of recent page → frame translations.
- B. On a context switch between processes, the TLB must be flushed, unless its entries are tagged with an ASID (address-space ID).
- C. TLB hit ratios are high mainly because of locality of reference.
- D. The TLB stores the data values held in memory.
- E. On x86, a TLB miss is handled by a hardware page walker; on some RISC designs (e.g., MIPS), the OS handles it in software.

**Answer:** A, B, C, E

**Why:**
- **A:** "Associative" means all entries are compared at once, which is why a TLB can be small and fast.
- **B:** Without tags, the old process's mappings would translate the new process's addresses.
- **C:** Programs keep touching the same few pages, so recent translations get reused.
- **E:** A **software-managed TLB** traps to the OS, which walks the table and inserts the entry. That makes the page-table format flexible.

**Why the other options fail:**
- D: That's the data cache. The TLB holds **translations** only.

**Concept recap:** TLB = a cache for translations. Reach = number of entries × page size. ASIDs (PCID on x86) avoid flushes on context switches.

**Trap alert:** A TLB miss isn't a page fault. On a miss the page is usually in memory and just needs a page walk. A page fault means the page isn't in memory, or the access is invalid.

**Interview angle:** "TLB miss vs page fault?" Follow-up: *"What is TLB reach, and how do you increase it?"* Entries × page size. Increase it with huge pages.

---

## Q166. [M6 | TLB hit ratio needed for a target EMAT | Hard | MCQ | OA]

**Two-level** page table. TLB lookup = 10 ns, memory access = 100 ns, and the TLB is checked first. What TLB hit ratio is needed for an EMAT of at most **120 ns**?

- A. 90%
- B. 95%
- C. 98%
- D. 99%

**Answer:** B

**Why:**
- Hit = 10 + 100 = **110 ns**.
- Miss = 10 + 100 (outer table) + 100 (inner table) + 100 (data) = **310 ns**.
- EMAT = 110h + 310(1 − h) = 310 − 200h.
- Require 310 − 200h ≤ 120 → 200h ≥ 190 → **h ≥ 0.95**.

**Why the other options fail:**
- A: At h = 0.90, EMAT = 310 − 180 = 130 ns, which is too slow. (90% would be enough for a *single*-level table, Q164.)
- C, D: These meet the target, but the question asks for the *minimum* hit ratio needed.

**Concept recap:** Solve h(hit cost) + (1 − h)(miss cost) ≤ target for h. Each extra page-table level adds one m to the miss cost.

**Trap alert:** Use the miss cost for the **correct number of levels**: 2 levels → 3 memory accesses on a miss.

**Interview angle:** Follow-up: *"Why do deeper tables need better TLBs?"* Each miss gets costlier as levels are added.

---

## Q167. [M6 | Shared pages | Easy | MCQ | Interview]

Thirty users run the same text editor. Its code is **200 KB** and is **reentrant**, so it is shared. Each user also has **40 KB** of private data. How much memory does this take with sharing?

- A. 1400 KB
- B. 7200 KB
- C. 6200 KB
- D. 1240 KB

**Answer:** A

**Why:** With paging, the page tables of all 30 processes map the editor's code pages to the **same frames**, so only one copy of the code is in memory. Each user still needs private data.
- With sharing: 200 + 30 × 40 = 200 + 1200 = **1400 KB**.
- Without sharing: 30 × (200 + 40) = 7200 KB.

**Why the other options fail:**
- B: The total without sharing.
- C: Shares the data instead of the code (30 × 200 + 200).
- D: An arithmetic slip.

**Concept recap:** Only **reentrant** (pure) code can be shared. Reentrant code never modifies itself, so it can be mapped read-only. Shared libraries (`libc.so`) work the same way.

**Trap alert:** **Data** is never shared this way. Each process needs its own writable copy, or copy-on-write (covered in Batch 11).

**Interview angle:** "How do shared libraries save memory?" One physical copy of the code is mapped into every process. Follow-up: *"Why must the shared code be position-independent?"* Different processes may map it at different virtual addresses.

---

## Q168. [M6 | Buddy system allocation | Easy | MCQ | OA]

A buddy allocator manages one **1 MB** block. It receives a request for **100 KB**. What block size is allocated, and how much internal fragmentation is there?

- A. 100 KB; 0 KB
- B. 256 KB; 156 KB
- C. 128 KB; 28 KB
- D. 128 KB; 0 KB

**Answer:** C

**Why:** A buddy allocator hands out blocks whose size is a **power of two**, rounding each request up. It splits larger blocks in half repeatedly until it gets the right size.
- Round 100 KB up → **128 KB**.
- Splits: 1024 → 512 + 512 → 256 + 256 → 128 + 128. One 128 KB block is allocated.
- Waste = 128 − 100 = **28 KB**.
- Still free: 128 KB, 256 KB and 512 KB.

**Why the other options fail:**
- A: The buddy system never allocates an exact, non-power-of-two size.
- B: 256 KB is not the smallest power of two that fits 100 KB.
- D: The rounding-up waste is real.

**Concept recap:** Allocated size = 2^⌈log₂ request⌉. Internal fragmentation = allocated size − request.

**Trap alert:** A request just above a power of two wastes almost half the block. For example, 65 KB → 128 KB wastes 63 KB.

**Interview angle:** "How does the buddy system work?" Split on allocation, and merge with the buddy on free (Q169). Follow-up: *"Its main weakness?"* Internal fragmentation from rounding to powers of two.

---

## Q169. [M6 | Buddy system properties | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. For a block of size 2ᵏ starting at address X, aligned to 2ᵏ (offsets measured from the start of the managed region), its buddy is at address X XOR 2ᵏ.
- B. When a block is freed, it is merged with its buddy only if the buddy is entirely free.
- C. The buddy system suffers internal fragmentation, from rounding requests up to powers of two.
- D. Any two adjacent free blocks of the same size can be merged.
- E. Linux uses a buddy allocator to manage physical page frames.

**Answer:** A, B, C, E

**Why:**
- **A:** Two buddies differ only in bit k of their address, so XOR flips between them. This makes finding a buddy, and merging, fast.
- **B:** Merging repeats upward as long as each buddy is free.
- **E:** `/proc/buddyinfo` on Linux shows the free blocks of each order, where order k means blocks of 2ᵏ pages.

**Why the other options fail:**
- D: Only **buddies**, the two halves produced by the same split, can merge. Two neighbouring same-size blocks from different parents can't, because the result wouldn't be an aligned power-of-two block.

**Concept recap:** Buddy = fast split and merge, little external fragmentation, but internal fragmentation from rounding. In Linux, the slab/SLUB allocator sits on top of it for small objects.

**Trap alert:** "Adjacent and same size ⇒ buddies" is false. Check with the XOR rule.

**Interview angle:** "Why does Linux use both a buddy allocator and a slab allocator?" The buddy allocator handles whole pages. Slab carves them into small, same-sized objects, avoiding rounding waste and saving initialization time. Follow-up: *"What do the orders mean?"* Order k = blocks of 2ᵏ contiguous pages.

---

## Q170. [M7 | Pure demand paging | Easy | MCQ | Interview]

Under **pure demand paging**, what happens when a process starts?

- A. All of its pages are loaded before the first instruction runs
- B. Only its code pages are loaded in advance
- C. Its pages are loaded in order of page number
- D. None of its pages are in memory at first; the very first instruction fetch causes a page fault, and each page is brought in only when it is first referenced

**Answer:** D

**Why:** Demand paging loads pages **lazily**, only when they're referenced. *Pure* demand paging takes this to the limit and loads nothing in advance. The program faults its way in, page by page, until its working set is resident.

**Why the other options fail:**
- A: That's swapping the whole process in.
- B, C: Those describe *prepaging* or loading in order, not pure demand paging.

**Concept recap:** Demand paging needs:
1. Valid/invalid bits in the PTEs (Q172).
2. Backing store (swap space or the executable file).
3. Restartable instructions (Q177).

**Trap alert:** Pure demand paging causes many faults at startup. Real systems often **prepage** (read ahead) to reduce them.

**Interview angle:** "What is demand paging?" Follow-up: *"Why does it work at all?"* Locality: a program uses only a small part of its pages at any one time.

---

## Q171. [M7 | Handling a page fault | Medium | MCQ | OA + Interview]

Put these page-fault handling steps in order:

1. Find a free frame (or choose a victim page to replace)
2. Trap to the OS; check its internal table to confirm the reference is valid
3. Restart the instruction that faulted
4. Schedule a disk read of the page into that frame
5. Update the page table: set the valid bit and the frame number

- A. 2, 4, 1, 5, 3
- B. 2, 1, 4, 5, 3
- C. 1, 2, 4, 5, 3
- D. 2, 1, 5, 4, 3

**Answer:** B

**Why:**
- The OS must first learn about the fault (the trap), and decide whether the access is legal or should kill the process (2).
- It needs somewhere to put the page (1) before it can read it in (4).
- It may mark the page valid only after the data has arrived (5).
- Finally, the faulting instruction runs again (3).

While the disk read is in progress, the CPU runs other processes.

**Why the other options fail:**
- A: Reads from disk before a frame has been chosen.
- C: Picks a frame before the OS even knows a fault happened.
- D: Marks the page valid before its data has been read in, so another access could see garbage.

**Concept recap:** trap → validate → get frame (evict if needed, writing the victim back first if it's dirty) → read page → update PTE → restart.

**Trap alert:** "Invalid PTE ⇒ always load from disk" is wrong. If the address isn't part of the address space at all, the OS kills the process (SIGSEGV) instead.

**Interview angle:** "Walk me through a page fault." Follow-up: *"What makes servicing a fault so slow?"* The disk I/O. It takes milliseconds, versus about 100 ns for a memory access (Q173).

---

## Q172. [M7 | Valid/invalid bits | Medium | MSQ | OA]

**Select all that apply.**

- A. An invalid PTE can mean either "this page is not part of the address space (illegal)" or "this page is legal but not currently in memory".
- B. The OS tells these two cases apart using its own data structures (on Linux, the process's list of virtual memory areas).
- C. When a process accesses a page whose PTE is invalid, the MMU traps (a page fault).
- D. A set valid bit means the page has been modified (is dirty).
- E. Under demand paging, a page's PTE is marked invalid until the page is loaded.

**Answer:** A, B, C, E

**Why:** The hardware knows only one thing: valid or not. It traps on invalid. The OS then consults its own records of what the process *should* have mapped:
- If the address is legal → load the page.
- If not → it's a segmentation fault.

**Why the other options fail:**
- D: "Modified" is the separate **dirty bit** (covered in Batch 10). Valid means the translation can be used.

**Concept recap:** Invalid PTE + legal address → page-fault service. Invalid PTE + illegal address → SIGSEGV.

**Trap alert:** One hardware bit, two meanings. The OS decides which one applies.

**Interview angle:** "How does the kernel know whether a fault is legitimate?" On Linux, it searches the process's `vm_area_struct` list (its VMAs). Follow-up: *"What's in a VMA?"* An address range, its permissions, and its backing (a file, or anonymous memory).

---

## Q173. [M7 | EMAT with page faults | Medium | Numerical | OA + Interview]

- Memory access time = **100 ns**.
- Average page-fault service time = **10 ms**.
- **1 access in every 2000** causes a page fault.

What is the effective access time?

- A. ≈ 105 ns
- B. ≈ 10.1 µs
- C. ≈ 5.1 µs
- D. 10 ms

**Answer:** C

**Why:** Effective access time is a weighted average. Most accesses take 100 ns, but the rare fault costs **100,000 times more**, so faults dominate the average.

**Working:**
- **Formula:** EAT = (1 − p)·ma + p·service.
- p = 1/2000 = 0.0005; ma = 100 ns; service = 10 ms = 10,000,000 ns.
- EAT = 0.9995 × 100 + 0.0005 × 10,000,000 = 99.95 + 5000 = **5099.95 ns ≈ 5.1 µs**.
- Slowdown = 5100 / 100 ≈ **51×**.

**Why the other options fail:**
- A: Treats 10 ms as 10 µs. This is the unit-conversion trap.
- B: Uses p = 1/1000 instead of 1/2000.
- D: The cost of one fault, not the average.

**Concept recap:** EAT = (1 − p)·ma + p·(fault service time). Even a tiny p dominates the result.

**Trap alert:** Convert everything to **one unit** before computing. 1 ms = 10⁶ ns.

**Interview angle:** "Why are page faults so damaging to performance?" Disk access is around 10⁵ × slower than RAM. Follow-up: *"What fault rate keeps the slowdown under 10%?"* See Q174.

---

## Q174. [M7 | Fault rate needed for a slowdown target | Hard | MCQ | OA]

Same system as Q173: memory access = 100 ns, fault service = 10 ms. To keep the slowdown caused by paging **under 10%**, i.e. EAT < 110 ns, the page-fault rate must be about:

- A. Below about 1 fault per 1,000,000 accesses (p < 10⁻⁶)
- B. Below about 1 fault per 1000 accesses
- C. Below about 1 fault per 100,000 accesses
- D. Below about 1 fault per 10⁹ accesses

**Answer:** A

**Why:** Set up the inequality and solve for p:
- 110 > (1 − p) × 100 + p × 10,000,000
- 110 > 100 + p × (10,000,000 − 100)
- 10 > p × 9,999,900
- p < 10 / 9,999,900 ≈ **1.0 × 10⁻⁶**

**Why the other options fail:**
- B: At p = 10⁻³, EAT ≈ 10.1 µs, about 100× slower.
- C: At p = 10⁻⁵, EAT ≈ 200 ns, about 2× slower.
- D: Much stricter than necessary.

**Concept recap:** p_max ≈ (allowed extra time) / (fault service time). Here 10 ns / 10 ms = 10⁻⁶.

**Trap alert:** The condition is the slowdown being **under 10%**, i.e. EAT < 1.1 × ma. It's not "EAT under 10 ns".

**Interview angle:** Follow-up: *"How do OSes keep the fault rate that low?"* Enough frames per process (working sets, covered in Batch 11), good replacement policies, and read-ahead.

---

## Q175. [M7 | FIFO replacement trace | Medium | Numerical | OA]

Reference string: **4 7 6 1 7 6 1 2 7 2**, with **3 frames**, using **FIFO** replacement. How many page faults occur?

- A. 5
- B. 7
- C. 10
- D. 6

**Answer:** D

**Why:** FIFO evicts the page that was **loaded earliest**. Keep a queue ordered by load time. A hit doesn't change the queue.

**Working:** (F = fault, H = hit; the queue is listed oldest first)

| # | Ref | Frames after | Result | FIFO queue (oldest → newest) |
|---|---|---|---|---|
| 1 | 4 | 4 – – | F | 4 |
| 2 | 7 | 4 7 – | F | 4 7 |
| 3 | 6 | 4 7 6 | F | 4 7 6 |
| 4 | 1 | 1 7 6 | F (evict 4) | 7 6 1 |
| 5 | 7 | 1 7 6 | H | 7 6 1 |
| 6 | 6 | 1 7 6 | H | 7 6 1 |
| 7 | 1 | 1 7 6 | H | 7 6 1 |
| 8 | 2 | 1 2 6 | F (evict 7) | 6 1 2 |
| 9 | 7 | 1 2 7 | F (evict 6) | 1 2 7 |
| 10 | 2 | 1 2 7 | H | 1 2 7 |

Faults = **6** (hits = 4).

**Why the other options fail:**
- A: The result for **Optimal** replacement on this string. At reference 8, Optimal evicts 1 or 6 (neither is used again), keeps 7, and so avoids the fault at reference 9.
- B: Typically from updating the FIFO queue on hits.
- C: Counts every reference as a fault.

**Concept recap:** In FIFO, the victim is the head of the load-order queue. **Hits don't reorder the queue.** That's the difference from LRU.

**Trap alert:** Moving a page to the back of the queue on a hit turns FIFO into LRU by accident.

**Interview angle:** "FIFO page replacement: pros and cons?" It's simple, but it can evict a heavily used page just because it's old. It also suffers Belady's anomaly (Q176). Follow-up: *"What fixes FIFO's blindness to use?"* Second-chance (clock), covered in Batch 10.

---

## Q176. [M7 | Belady's anomaly with FIFO | Hard | Numerical | OA + Interview]

Reference string: **1 2 3 4 1 2 5 1 2 3 4 5**, using **FIFO**. How many page faults occur with **3 frames**, and with **4 frames**?

- A. 9 and 8
- B. 10 and 9
- C. 9 and 10
- D. 12 and 12

**Answer:** C

**Why:** With FIFO, **adding a frame can increase the number of faults**. That's **Belady's anomaly**.

**Working: 3 frames**

| Ref | 1 | 2 | 3 | 4 | 1 | 2 | 5 | 1 | 2 | 3 | 4 | 5 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Frames | 1 | 1 2 | 1 2 3 | 4 2 3 | 4 1 3 | 4 1 2 | 5 1 2 | 5 1 2 | 5 1 2 | 5 3 2 | 5 3 4 | 5 3 4 |
| Result | F | F | F | F | F | F | F | H | H | F | F | H |

Faults = **9**.

**Working: 4 frames**

| Ref | 1 | 2 | 3 | 4 | 1 | 2 | 5 | 1 | 2 | 3 | 4 | 5 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Frames | 1 | 1 2 | 1 2 3 | 1 2 3 4 | 1 2 3 4 | 1 2 3 4 | 5 2 3 4 | 5 1 3 4 | 5 1 2 4 | 5 1 2 3 | 4 1 2 3 | 4 5 2 3 |
| Result | F | F | F | F | H | H | F | F | F | F | F | F |

Faults = **10**.

**Why the other options fail:**
- A: Assumes more frames always means fewer faults.
- B: Arithmetic slips in both traces.
- D: Counts every reference as a fault.

**Concept recap:** Belady's anomaly = faults(n + 1 frames) > faults(n frames). FIFO can show it. LRU and Optimal cannot, because they're **stack algorithms** (covered in Batch 10).

**Trap alert:** With 4 frames, FIFO keeps an "older" set of pages that happens to be wrong for the rest of this string.

**Interview angle:** "What is Belady's anomaly?" Define it and name an algorithm that suffers from it. Follow-up: *"Why can't LRU show it?"* The set of pages LRU keeps with n frames is always a subset of the set it keeps with n + 1 frames.

---

## Q177. [M7 | Restarting a faulted instruction | Hard | MCQ | Interview]

A single instruction copies 256 bytes from a source region to a destination region. The regions may **overlap**, and either may cross a page boundary. Partway through the copy, a destination page faults. Why is simply restarting the instruction from the beginning dangerous, and what is a standard fix?

- A. It isn't dangerous; restarting is always correct
- B. If the regions overlap, part of the source may already have been overwritten, so re-running from the start copies wrong data. Fixes: touch (probe) every page of both regions before modifying anything, or save the overwritten values in temporary registers so the instruction can be undone
- C. The CPU loses the program counter when it faults
- D. The disk can't load a page while an instruction is in progress

**Answer:** B

**Why:** Demand paging assumes that any instruction can be **restarted** after a fault, as if it had never begun. An instruction that changes memory in several places breaks that assumption, if some changes already happened and the source data depended on them. Real hardware (IBM System/370's MVC, x86 string instructions with saved progress) handles this, either by checking everything up front or by recording enough state to resume or undo.

**Why the other options fail:**
- A: The overlap case shows that restarting can be wrong.
- C: The CPU saves the faulting PC precisely so the instruction can be restarted (Q7).
- D: The page is loaded while the instruction is paused, which is exactly what demand paging does.

**Concept recap:** Demand paging requires instructions to be **restartable**: either nothing visible happened before the fault, or the state can be resumed or undone.

**Trap alert:** "Re-executing an instruction is always harmless" is false for instructions with **side effects** spread across several pages.

**Interview angle:** "What hardware support does demand paging need?" Valid bits, a precise fault PC, and restartable instructions. Follow-up: *"How does x86 `rep movsb` handle it?"* Its progress lives in registers (RSI, RDI, RCX), so after a fault it resumes from where it stopped.

---

## Q178. [M7 | Benefits of demand paging | Easy | MSQ | OA + Interview]

**Select all that apply.**

- A. Demand paging lets a program larger than physical memory run.
- B. Programs can start faster, because only the pages they need are loaded.
- C. Pages that are never referenced are never loaded.
- D. Every page fault requires writing a page out to disk.
- E. Swap space (the backing store) holds pages that aren't currently in memory.

**Answer:** A, B, C, E

**Why:**
- **A:** Only the pages in active use need frames, so a program bigger than RAM can still run.
- **B, C:** Unused code, such as error-handling paths, is never read in, which saves both I/O and memory.
- **E:** Pages that aren't resident live on disk: in swap for anonymous memory, or in the original file for code and file data.

**Why the other options fail:**
- D: A write-out is needed only if the **victim page is dirty**. Clean pages, such as code, can simply be dropped and re-read from the file later (covered in Batch 10).

**Concept recap:** Demand paging gives a larger effective memory, faster starts and higher multiprogramming, at the cost of page faults.

**Trap alert:** A fault costs *one* disk transfer (read the page in) when the victim is clean, and *two* when it's dirty (write the victim out, then read the page in).

**Interview angle:** "Benefits of virtual memory?" Follow-up: *"What's backing a page of `malloc`'d memory versus a page of code?"* Swap (or the zero page before first use) for `malloc`; the executable file for code.

---

## Q179. [M7 | Minor vs major page faults | Medium | MCQ | OA]

On Linux, which of these is a **minor** (soft) page fault?

- A. A process touches a file-backed page that is already in the page cache (because another process read it); the kernel just maps the existing frame into this process
- B. A page that was swapped out must be read back from swap
- C. A file page that isn't in memory must be read from disk
- D. A process references an address that isn't mapped at all

**Answer:** A

**Why:** A **minor fault** needs no disk I/O. The data is already in RAM, and only this process's page table is missing the mapping. A **major fault** needs disk I/O. Minor faults cost microseconds; major ones cost milliseconds.

**Why the other options fail:**
- B, C: Both need disk I/O, so they're major faults.
- D: That isn't a fault the kernel can service. The process gets SIGSEGV.

**Concept recap:**
- Minor = map an existing frame (shared page-cache pages, the zero page, COW setup).
- Major = read from disk.

You can see the counts with `/usr/bin/time -v` or `ps -o min_flt,maj_flt`.

**Trap alert:** "Page fault ⇒ disk I/O" isn't always true. Many faults are minor.

**Interview angle:** "Your service shows lots of page faults. Is that bad?" First check whether they're major or minor. Follow-up: *"What causes a burst of major faults?"* Memory pressure that evicted the working set, or a cold start.

---

## Q180. [M7 | What FIFO evicts | Easy | MCQ | OA]

Which page does **FIFO** replacement evict?

- A. The least recently used page
- B. The least frequently used page
- C. The page that won't be needed for the longest time in the future
- D. The page that has been in memory the longest (loaded earliest), however recently it was used

**Answer:** D

**Why:** FIFO tracks only **when each page arrived**. It never looks at how recently or how often a page is used. That makes it simple, and also blind: it may evict a page in constant use, just because that page is old.

**Why the other options fail:**
- A: That's LRU.
- B: That's LFU.
- C: That's Optimal (OPT/MIN).

**Concept recap:** FIFO = oldest load. LRU = oldest *use*. OPT = farthest *future* use. LFU = fewest uses.

**Trap alert:** "Oldest page" is ambiguous. In FIFO it means the oldest by **load time**, not by last use.

**Interview angle:** "Name the classic page-replacement algorithms with a one-line rule each." Follow-up: *"Which one is optimal, and why can't we use it?"* OPT, because it needs to know future references (covered in Batch 10).

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 161 | B | 166 | B | 171 | B | 176 | C |
| 162 | D | 167 | A | 172 | A, B, C, E | 177 | B |
| 163 | A | 168 | C | 173 | C | 178 | A, B, C, E |
| 164 | C | 169 | A, B, C, E | 174 | A | 179 | A |
| 165 | A, B, C, E | 170 | D | 175 | D | 180 | D |

---

**Concepts tested so far (to avoid repeats):** Batches 1–8, plus: sizing a two-level table for a sparse address space · n + 1 accesses for an n-level walk · inverted page tables (one per system, hashed, hard to share) · EMAT with a serial TLB · hit ratio needed for a target with 2 levels · TLB facts (associative, ASID, locality, hardware vs software walk) · shared reentrant code arithmetic · buddy rounding and waste · buddy XOR, merging only with buddies, Linux buddyinfo · pure demand paging · page-fault handling order · two meanings of the invalid bit and VMAs · EAT with page faults · fault rate for a 10% slowdown · FIFO trace (hits don't reorder) · Belady's anomaly 9 → 10 · instruction-restart problem · minor vs major faults · what FIFO evicts
*(M6 is complete. Batch 10 continues M7 with Optimal, LRU (implementations and stack property), Second-chance/Clock and Enhanced clock, LFU/MFU, dirty and reference bits, and locality.)*

**Completed: 180 / 300    Remaining: 120**
