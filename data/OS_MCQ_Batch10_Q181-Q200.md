# OS MCQ/MSQ Bank — Batch 10 (Q181–Q200)

**What this batch covers (M7, page replacement):**
- Optimal (OPT), and OPT vs FIFO vs LRU on one string.
- LRU: a trace, how it's implemented, stack algorithms and why they can't show Belady's anomaly.
- Second-chance/Clock (trace and properties) and enhanced second-chance (R, M classes).
- Dirty bits and the effective fault-service time.
- Reference bits and the additional-reference-bits algorithm.
- LFU and MFU.
- Locality: temporal vs spatial, and how loop order changes the fault count.
- What Linux uses, and scan pollution under LRU.

**Progress report #2 follows the answer key at the end of this file.**

> **Trace conventions:**
> - All frames start empty, and every first reference is a fault.
> - F = fault, H = hit. Frame contents are shown after each reference.
> - Ties in OPT (two pages never used again) go to the page in the lower-numbered frame.

---

## Q181. [M7 | Optimal replacement | Easy | MCQ | Interview]

Which page does **Optimal (OPT, also called MIN)** replacement evict?

- A. The page that was loaded earliest
- B. The page that won't be used for the longest time in the future. It gives the fewest possible faults, but needs knowledge of the future, so it's used as a benchmark
- C. The least recently used page
- D. The page with the fewest references so far

**Answer:** B

**Why:** Of the pages in memory, evict the one whose next use is farthest away, or that is never used again. Any other choice brings in a page that will be needed sooner, so it can't do better. The OS can't know future references, so OPT can't be implemented. It serves as the yardstick: "LRU is within X% of optimal".

**Why the other options fail:**
- A: That's FIFO.
- C: That's LRU.
- D: That's LFU.

**Concept recap:** OPT = farthest next use. It minimizes faults for any given reference string and number of frames, and it never shows Belady's anomaly.

**Trap alert:** "Optimal" in the name doesn't mean any OS actually runs it. It's a theoretical bound.

**Interview angle:** "Why can't we implement OPT?" It needs to know future references. Follow-up: *"Then why study it?"* To measure how close practical algorithms come to the best possible result.

---

## Q182. [M7 | OPT vs FIFO vs LRU on one string | Hard | Numerical | OA]

Reference string: **2 3 4 2 1 3 7 5 4 3**, with **3 frames**. How many faults do **OPT, FIFO and LRU** each produce, in that order?

- A. 6, 8, 9
- B. 6, 9, 8
- C. 5, 8, 9
- D. 6, 8, 8

**Answer:** A

**Why:** Each algorithm picks a different victim at the same fault, and the differences add up. Notice that on this string, **LRU does worse than FIFO**.

**Working: OPT** (next use is shown as a position in the string)

| # | Ref | Frames | Result | Reason |
|---|---|---|---|---|
| 1–3 | 2 3 4 | 2 3 4 | F F F | |
| 4 | 2 | 2 3 4 | H | |
| 5 | 1 | 1 3 4 | F | next uses: 2 never, 3 at #6, 4 at #9 → evict 2 |
| 6 | 3 | 1 3 4 | H | |
| 7 | 7 | 7 3 4 | F | 1 never, 3 at #10, 4 at #9 → evict 1 |
| 8 | 5 | 5 3 4 | F | 7 never → evict 7 |
| 9 | 4 | 5 3 4 | H | |
| 10 | 3 | 5 3 4 | H | |

OPT faults = **6**.

**Working: FIFO** (queue oldest → newest)

| Ref | 2 | 3 | 4 | 2 | 1 | 3 | 7 | 5 | 4 | 3 |
|---|---|---|---|---|---|---|---|---|---|---|
| Result | F | F | F | H | F (evict 2) | H | F (evict 3) | F (evict 4) | F (evict 1) | F (evict 7) |
| Queue after | 2 | 2 3 | 2 3 4 | 2 3 4 | 3 4 1 | 3 4 1 | 4 1 7 | 1 7 5 | 7 5 4 | 5 4 3 |

FIFO faults = **8**.

**Working: LRU** (recency list, least recent → most recent)

| Ref | 2 | 3 | 4 | 2 | 1 | 3 | 7 | 5 | 4 | 3 |
|---|---|---|---|---|---|---|---|---|---|---|
| Result | F | F | F | H | F (evict 3) | F (evict 4) | F (evict 2) | F (evict 1) | F (evict 3) | F (evict 7) |
| Recency after | 2 | 2 3 | 2 3 4 | 3 4 2 | 4 2 1 | 2 1 3 | 1 3 7 | 3 7 5 | 7 5 4 | 5 4 3 |

LRU faults = **9**.

**Why the other options fail:**
- B: Swaps the FIFO and LRU counts.
- C: OPT faults 6 times here, not 5.
- D: Forgets one of LRU's faults.

**Concept recap:** OPT ≤ any other algorithm on the same string. LRU is usually better than FIFO **on typical programs** (which have locality), but not on every string.

**Trap alert:** "LRU always beats FIFO" is false. This string has little reuse of recent pages, so LRU's assumption doesn't pay off.

**Interview angle:** "Is LRU always better than FIFO?" No. Give an example like this one. Follow-up: *"Why is LRU still preferred?"* Real programs have strong temporal locality, which LRU exploits.

---

## Q183. [M7 | LRU trace | Medium | Numerical | OA + Interview]

Reference string: **5 0 1 0 2 0 3 0 4 2 3**, with **3 frames**, using **LRU**. How many page faults occur?

- A. 6
- B. 9
- C. 8
- D. 7

**Answer:** C

**Why:** LRU evicts the page whose **last use** is oldest. Every hit moves that page to the "most recent" end of the list. That's the key difference from FIFO.

**Working:** (recency list, least recent → most recent)

| # | Ref | Frames | Result | Recency after |
|---|---|---|---|---|
| 1 | 5 | 5 | F | 5 |
| 2 | 0 | 5 0 | F | 5 0 |
| 3 | 1 | 5 0 1 | F | 5 0 1 |
| 4 | 0 | 5 0 1 | H | 5 1 0 |
| 5 | 2 | 2 0 1 | F (evict 5) | 1 0 2 |
| 6 | 0 | 2 0 1 | H | 1 2 0 |
| 7 | 3 | 2 0 3 | F (evict 1) | 2 0 3 |
| 8 | 0 | 2 0 3 | H | 2 3 0 |
| 9 | 4 | 4 0 3 | F (evict 2) | 3 0 4 |
| 10 | 2 | 4 0 2 | F (evict 3) | 0 4 2 |
| 11 | 3 | 4 3 2 | F (evict 0) | 4 2 3 |

Faults = **8**.

**Why the other options fail:**
- A: The **OPT** count on this string. OPT keeps 2 and 3 at #9 by evicting 0, which is never used again.
- B: The **FIFO** count. FIFO evicts 0 at #7 even though 0 is constantly used.
- D: Usually comes from forgetting to update the recency list on a hit.

**Concept recap:** LRU victim = the least recently used page. **Update the recency list on every reference**, hits included.

**Trap alert:** At #11, page 0 was used often earlier, but its last use (#8) is now the oldest, so LRU evicts it. Recency, not frequency, decides.

**Interview angle:** "Implement an LRU cache." Use a hash map plus a doubly linked list, giving O(1) get and put. Follow-up: *"How is that different from LRU page replacement in an OS?"* The OS can't afford to update a list on every memory access, so it approximates (Q184).

---

## Q184. [M7 | Implementing LRU | Medium | MCQ | Interview]

Why do real OSes **approximate** LRU instead of implementing it exactly?

- A. LRU suffers from Belady's anomaly
- B. Exact LRU needs knowledge of the future
- C. LRU can't be implemented at all
- D. Exact LRU must update a timestamp, or move a list entry, on **every memory reference**. Without special hardware that's far too costly, so the hardware sets a reference bit and the OS approximates LRU from it

**Answer:** D

**Why:** There are two textbook ways to implement exact LRU:
- **Counters:** store a "time of last use" in each page entry on every access, and on a fault search for the minimum.
- **Stack:** keep a doubly linked list and move the page to the top on every access. That's several pointer updates per access.

Either way, the work happens on **every load and store**, billions of times a second. No commodity CPU does that. Instead, hardware just sets an **accessed (reference) bit**, and the OS samples those bits.

**Why the other options fail:**
- A: LRU is a stack algorithm, so it's immune to Belady's anomaly (Q185).
- B: That's OPT. LRU looks at the past.
- C: It can be implemented; it's just too expensive.

**Concept recap:** Exact LRU costs work on every reference. Practical systems approximate it with reference bits: second-chance/Clock (Q186–187), additional-reference-bits (Q191), and Linux's active/inactive lists (Q195).

**Trap alert:** Software LRU caches (in Redis or your own code) *are* exact, because they run per *request*, not per memory access.

**Interview angle:** "How would you implement LRU page replacement in a kernel?" Explain why exact LRU is infeasible and describe Clock. Follow-up: *"What hardware support does Clock need?"* Just the reference bit in each PTE.

---

## Q185. [M7 | Stack algorithms and Belady's anomaly | Hard | MSQ | Interview]

**Select all that apply.**

- A. LRU never shows Belady's anomaly.
- B. OPT never shows Belady's anomaly.
- C. For a stack algorithm, the set of pages in memory with n frames is always a subset of the set with n + 1 frames.
- D. FIFO is a stack algorithm.
- E. Second-chance (Clock) can show Belady's anomaly.

**Answer:** A, B, C, E

**Why:** **Stack algorithms** have the **inclusion property** (C). With more frames, the memory contents are a superset of what they'd be with fewer frames, so anything that hits with n frames also hits with n + 1. Faults can therefore never increase. LRU (the n most recently used pages) and OPT (the n pages used soonest) both have this property.

**Why the other options fail:**
- D: FIFO's contents depend on load order, which changes with the number of frames. Q176 showed 9 → 10 faults as frames went from 3 to 4.

**Concept recap:** Stack algorithm ⇒ immune to Belady's anomaly. LRU and OPT are stack algorithms. FIFO and Clock are not.

**Trap alert:** "Belady's anomaly in LRU" is a classic false statement in exams. LRU can't show it.

**Interview angle:** "Which algorithms are immune to Belady's anomaly, and why?" Explain inclusion. Follow-up: *"Why isn't Clock immune, if it approximates LRU?"* When all reference bits are 1, Clock behaves like FIFO (Q197 A), so it inherits FIFO's anomaly.

---

## Q186. [M7 | Clock (second-chance) trace | Hard | Numerical | OA]

Clock replacement with **3 frames** (F0, F1, F2). The hand starts at F0, and every frame starts empty.
- **On a fault:** if the frame under the hand is empty, or its R bit is 0, load the page there, set R = 1, and advance the hand. If R = 1, set it to 0, advance the hand, and repeat.
- **On a hit:** set the page's R bit to 1. The hand doesn't move.

Reference string: **1 2 3 2 4 1 5 2**. How many faults occur, and what is in memory at the end?

- A. 7 faults; {2, 1, 5}
- B. 6 faults; {4, 2, 5}
- C. 7 faults; {4, 1, 5}
- D. 6 faults; {1, 2, 5}

**Answer:** A

**Why:** Clock is FIFO order plus a **second chance** for any page whose reference bit is set. Watch what happens to page 2. It's hit at step 4, but the sweep at step 5 clears its bit, so by step 6 it's an ordinary victim again.

**Working:** (frame contents with R bits; ▸ marks the hand position *after* the step)

| # | Ref | Action | F0 | F1 | F2 | Result |
|---|---|---|---|---|---|---|
| 1 | 1 | load into empty F0 | 1(1) | – | ▸– | F |
| 2 | 2 | load into empty F1 | 1(1) | 2(1) | ▸– | F |
| 3 | 3 | load into empty F2 | ▸1(1) | 2(1) | 3(1) | F |
| 4 | 2 | hit: R(2) = 1 | ▸1(1) | 2(1) | 3(1) | H |
| 5 | 4 | sweep clears F0, F1, F2; back at F0, R = 0 → replace 1 | 4(1) | ▸2(0) | 3(0) | F |
| 6 | 1 | F1 has R = 0 → replace 2 | 4(1) | 1(1) | ▸3(0) | F |
| 7 | 5 | F2 has R = 0 → replace 3 | ▸4(1) | 1(1) | 5(1) | F |
| 8 | 2 | sweep clears F0, F1, F2; back at F0 → replace 4 | 2(1) | ▸1(0) | 5(0) | F |

Faults = **7**; final contents **{2, 1, 5}**.

**Why the other options fail:**
- B, D: Count only 6 faults, which wrongly assumes page 2 survives because of its earlier hit.
- C: Stops before the final replacement at step 8.

**Concept recap:** Clock: on a fault, sweep; if R = 1, clear it and move on; if R = 0, evict. When **every** R bit is 1, the hand sweeps all the way round and evicts the page it started at, which is pure FIFO (steps 5 and 8).

**Trap alert:** A hit protects a page only until the next sweep reaches it. The second chance is a **single** chance, used up when the hand passes.

**Interview angle:** "Explain the Clock algorithm." Follow-up: *"What happens when memory pressure is high and every page has R = 1?"* Clock degenerates into FIFO, and the hand moves fast. A fast-moving hand is a sign of memory pressure.

---

## Q187. [M7 | Second-chance rule | Easy | MCQ | OA]

In the second-chance (Clock) algorithm, when the hand points at a page whose reference bit is **1**, the algorithm:

- A. Evicts that page immediately
- B. Sets the bit to 1 and stops
- C. Clears the bit to 0 and moves the hand to the next page, giving this page a second chance
- D. Writes the page to disk and keeps it in memory

**Answer:** C

**Why:** R = 1 means "used since the hand last passed". Clock treats that as evidence of recent use, so it spares the page *once* by clearing the bit. If the page isn't used again before the hand comes back, its bit will still be 0, and it gets evicted.

**Why the other options fail:**
- A: That would be plain FIFO.
- B: The algorithm never sets bits. Hardware sets them on access.
- D: Writing pages to disk is about the **dirty** bit, not the reference bit.

**Concept recap:** R = 1 → clear the bit and move on. R = 0 → evict. Hardware sets R on any access.

**Trap alert:** Don't mix up the R (reference) bit with the M (dirty) bit. R is about recency. M is about whether the page must be written back.

**Interview angle:** Follow-up: *"Why is it called 'Clock'?"* The frames are arranged in a circular list, and the pointer sweeps around it like a clock hand.

---

## Q188. [M7 | Enhanced second-chance | Medium | MCQ | OA + Interview]

The enhanced second-chance algorithm classifies pages by the pair **(R, M)**: reference bit and modify (dirty) bit. In which order does it prefer victims, from best to worst?

- A. (0,0) → (0,1) → (1,0) → (1,1)
- B. (0,0) → (1,0) → (0,1) → (1,1)
- C. (1,1) → (1,0) → (0,1) → (0,0)
- D. (0,1) → (0,0) → (1,1) → (1,0)

**Answer:** A

**Why:**
- **(0,0):** not recently used and clean. The best victim: it's probably not needed, and there's no write-back cost.
- **(0,1):** not recently used, but dirty. It must be written out first, but it isn't in active use.
- **(1,0):** recently used, but clean. It's likely to be needed again soon.
- **(1,1):** recently used and dirty. The worst victim: likely needed soon, *and* costly to evict.

Recency (R) outranks cleanliness (M). Evicting a page that will be needed again soon costs a whole future fault, which is worse than one extra write.

**Why the other options fail:**
- B: Puts "clean" above "not recently used". That keeps an unused dirty page while evicting an active one.
- C: Completely reversed.
- D: Puts a dirty page ahead of a clean, unused one.

**Concept recap:** Victim priority: (0,0) < (0,1) < (1,0) < (1,1). The hand may circle several times to find the lowest non-empty class.

**Trap alert:** Many students rank (1,0) ahead of (0,1) because "clean is cheaper". The rule is that **recent use matters more**.

**Interview angle:** "How does the dirty bit influence replacement?" Follow-up: *"Can you get rid of the extra write cost altogether?"* A background page-out daemon cleans dirty pages ahead of time (Q198 C).

---

## Q189. [M7 | Dirty bit and fault-service time | Medium | MCQ | OA]

Servicing a page fault takes **8 ms** if the victim page is **clean**, and **20 ms** if it is **dirty** (it must be written out first). **70%** of victims are dirty. What is the effective page-fault service time?

- A. 14 ms
- B. 28 ms
- C. 20 ms
- D. 16.4 ms

**Answer:** D

**Why:** The service time is a weighted average over the two cases:
- (1 − d) × clean + d × dirty
- = 0.3 × 8 + 0.7 × 20
- = 2.4 + 14 = **16.4 ms**

**Why the other options fail:**
- A: Counts only the dirty term (0.7 × 20).
- B: Adds both times, as if every fault involved both a clean and a dirty case.
- C: Assumes every victim is dirty.

**Concept recap:** Effective fault service = (1 − d)·t_clean + d·t_dirty. Plug this into EAT = (1 − p)·ma + p·(fault service) (Q173).

**Trap alert:** Check whether "20 ms" is the *total* time for a dirty victim (as here) or the *extra* write time to add on top of 8 ms. Exams phrase this both ways.

**Interview angle:** "Why does the dirty bit matter?" Without it, the OS would have to write *every* victim back, roughly doubling the I/O. Follow-up: *"Which pages are usually clean?"* Code pages and unmodified file pages.

---

## Q190. [M7 | Reference bits | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Hardware sets a page's reference (accessed) bit whenever the page is read or written.
- B. The OS periodically clears reference bits, so it can see which pages were used recently.
- C. The additional-reference-bits algorithm keeps an 8-bit history per page. At each timer interval it shifts the R bit into the top of the history; the page with the smallest value is the best LRU-like victim.
- D. The dirty bit is set on reads.
- E. On x86, the accessed (A) and dirty (D) bits are stored in the PTE.

**Answer:** A, B, C, E

**Why:**
- **A, B:** One bit can't say *when* a page was used, only *whether* it was used since the bit was last cleared. Clearing the bits periodically turns them into a "used in the last interval?" signal.
- **C:** Shifting the bit into a history byte keeps 8 intervals of usage, which approximates LRU order (Q191).
- **E:** The hardware page walker updates these bits directly in the PTE.

**Why the other options fail:**
- D: The dirty bit is set **only on writes**.

**Concept recap:**
- R (accessed) = any access → approximates recency.
- M (dirty) = writes only → decides whether a write-back is needed.

**Trap alert:** After the OS clears an R bit in the PTE, it may also need to flush that page's TLB entry. Otherwise the hardware may not set the bit again on the next access, because the cached translation already recorded the access.

**Interview angle:** "What hardware support does page replacement need?" The accessed and dirty bits. Follow-up: *"How does the OS 'sample' page use?"* Clear the bits, wait, then read them back.

---

## Q191. [M7 | Additional-reference-bits algorithm | Hard | MCQ | OA]

Each page keeps an 8-bit history. The **most significant bit** is the **most recent** interval. Treating each byte as an unsigned number, which page is evicted?

| Page | History |
|---|---|
| P1 | 11000100 |
| P2 | 01110111 |
| P3 | 00110011 |
| P4 | 10000000 |

- A. P4
- B. P3
- C. P2
- D. P1

**Answer:** B

**Why:** Because recent intervals sit in the high bits, a smaller number means "least recently used". Convert each history:
- P1 = 11000100₂ = **196**
- P2 = 01110111₂ = **119**
- P3 = 00110011₂ = **51** ← smallest
- P4 = 10000000₂ = **128**

P3 hasn't been used in the last two intervals, while every other page was used in at least one of them. **Evict P3.**

**Why the other options fail:**
- A: P4 has the **fewest 1-bits** (it was used only once), but that use was in the *most recent* interval. Counting 1-bits is LFU thinking, not LRU.
- C: P2 was used in the second-most-recent interval and three more after that.
- D: P1 was used in the two most recent intervals.

**Concept recap:** Interpret the history as an unsigned integer and evict the smallest. Recency outweighs frequency because the high bits dominate.

**Trap alert:** Don't count the ones. Compare the values, where the most recent bit carries the most weight.

**Interview angle:** Follow-up: *"What if two pages have the same history?"* Break the tie with FIFO, or evict all pages that share the smallest value.

---

## Q192. [M7 | LFU and MFU | Medium | MCQ | OA]

What is the main problem with plain **LFU** (Least Frequently Used) replacement, and what is a standard fix?

- A. It needs knowledge of the future
- B. It always suffers from Belady's anomaly
- C. It always evicts the newest page
- D. A page used heavily early on keeps a high count and stays in memory long after it stops being used. The fix is to decay the counts periodically (e.g., shift each count right by 1)

**Answer:** D

**Why:** LFU counts *all* past uses equally, so old popularity lingers. Decaying the counts makes old uses matter less, so the count reflects recent popularity. **MFU** (Most Frequently Used) argues the opposite: a page with a small count may have *just been loaded* and is about to be used. Neither is common in practice, because both approximate OPT poorly.

**Why the other options fail:**
- A: That's OPT.
- B: Belady's anomaly isn't LFU's defining weakness.
- C: Plain LFU actually tends to *penalize* new pages, since they start with low counts.

**Concept recap:** LFU = evict the smallest count. MFU = evict the largest count. Use decay or aging to forget old history.

**Trap alert:** LFU's two weaknesses: new pages get evicted too soon, and old popular pages stick around too long.

**Interview angle:** "Design an LFU cache with O(1) operations." Use a map from key to node plus a map from frequency to a list of nodes. Follow-up: *"When does LFU beat LRU?"* When popularity is stable over time, such as a CDN's long-tail content.

---

## Q193. [M7 | Temporal vs spatial locality | Easy | MCQ | Interview]

In `for (i = 0; i < n; i++) sum += a[i];`, which kind of locality do the accesses to `a[i]` show, and which do the repeated uses of `sum` and of the loop's own instructions show?

- A. Temporal; spatial
- B. Spatial; temporal
- C. Both spatial
- D. Both temporal

**Answer:** B

**Why:**
- **Spatial locality:** if you access an address, you'll soon access *nearby* addresses. `a[0]`, `a[1]`, … are consecutive.
- **Temporal locality:** if you access an address, you'll soon access the *same* address again. `sum` and the loop's instructions are reused on every iteration.

**Why the other options fail:**
- A: Swapped.
- C, D: Each kind of access shows a different kind of locality.

**Concept recap:** Locality is why caches, TLBs and demand paging work: a small fast store catches most accesses.

**Trap alert:** One element `a[i]` is used only once. Its locality is spatial: the *next* element is nearby, on the same cache line or page.

**Interview angle:** "Why are linked lists slower than arrays to traverse, even at the same O(n)?" Array elements are contiguous (spatial locality). List nodes are scattered, so traversal causes cache and TLB misses. Follow-up: *"How does a prefetcher exploit spatial locality?"* It detects sequential or strided access and fetches ahead.

---

## Q194. [M7 | Loop order and page faults | Hard | MCQ | OA + Interview]

`int A[1024][1024];` uses 4-byte ints, stored **row-major**, page-aligned. Page size = **4 KB**, so each row occupies exactly one page. The process has **only one frame** for the array. Code and loop variables live elsewhere and never fault.

```c
// Loop 1 (row-wise)                // Loop 2 (column-wise)
for (i = 0; i < 1024; i++)          for (j = 0; j < 1024; j++)
    for (j = 0; j < 1024; j++)          for (i = 0; i < 1024; i++)
        A[i][j] = 0;                        A[i][j] = 0;
```

How many page faults does each loop cause?

- A. 1024 and 1024
- B. 1,048,576 and 1024
- C. 1024 and 1,048,576
- D. 1,048,576 and 1,048,576

**Answer:** C

**Why:** The fault count depends on **how often consecutive accesses switch pages**.
- **Loop 1** walks along a row: 1024 writes to the same page, then the next row. That's 1 fault per row, **1024 faults** in total.
- **Loop 2** walks down a column: each `A[i][j]` with a new `i` is in a *different row*, so a different page. With one frame, **every access faults**: 1024 × 1024 = **1,048,576 faults**.

**Why the other options fail:**
- A: Assumes the traversal order doesn't matter.
- B: Swapped.
- D: Assumes both loops fault on every access.

**Concept recap:** Traverse arrays in **storage order**: row-major in C/C++, column-major in Fortran and MATLAB. This exploits spatial locality.

**Trap alert:** Even with more frames, Loop 2 is much worse. Once the number of rows touched exceeds the number of frames, each column sweep evicts the pages it will need again.

**Interview angle:** "Why can swapping two loops make code 100× faster?" Locality, at the level of cache lines, TLB entries and pages. Follow-up: *"How would you speed up a matrix transpose?"* Loop tiling (blocking), so each block fits in the cache.

---

## Q195. [M7 | What Linux actually uses | Medium | MCQ | OA + Interview]

Which description best matches how **Linux** chooses pages to reclaim?

- A. Exact LRU, using a timestamp updated on every access
- B. Pure FIFO
- C. Optimal replacement, using profiling data
- D. Approximations of LRU based on the hardware accessed bits, organized into active and inactive lists (newer kernels can also use a multi-generational LRU)

**Answer:** D

**Why:**
- Linux keeps pages on **active** and **inactive** lists, separately for file-backed and anonymous memory.
- A page on the inactive list that gets referenced is promoted to the active list.
- Pages that age out of the active list are demoted, and reclaim takes victims from the tail of the inactive list.
- This is an LRU **approximation** driven by the accessed bits, with no per-access bookkeeping.

**Why the other options fail:**
- A: Too costly to run on every access (Q184).
- B: Too blind to how pages are used.
- C: Impossible, because it needs the future.

**Concept recap:** Real kernels = Clock-like scans of accessed bits + lists that separate pages used once from pages used repeatedly.

**Trap alert:** Linux's "LRU" lists aren't exact LRU. The name refers to the intent.

**Interview angle:** "How does the Linux kernel decide what to evict under memory pressure?" Describe active/inactive lists and `kswapd`. Follow-up: *"Why separate file pages from anonymous pages?"* Clean file pages can be dropped for free, while anonymous pages need swap I/O.

---

## Q196. [M7 | General facts about replacement | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. For a given reference string and number of frames, OPT produces the fewest possible page faults.
- B. OPT is used as a yardstick for evaluating practical algorithms.
- C. OPT produces the same number of faults on a reference string and on its reverse, and the same is true of LRU.
- D. Adding frames always reduces the number of faults, for any algorithm.
- E. FIFO can evict a heavily used page.

**Answer:** A, B, C, E

**Why:**
- **A, B:** OPT is optimal by construction (Q181).
- **C:** A known symmetry result, stated in Silberschatz. It's useful for checking traces.
- **E:** FIFO ignores use entirely, so a page loaded early and used constantly is still the next victim (Q180).

**Why the other options fail:**
- D: FIFO can show **Belady's anomaly** (Q176). "Always" fails.

**Concept recap:** OPT is a lower bound on faults. Stack algorithms (OPT, LRU) never get worse with more frames. FIFO and Clock can.

**Trap alert:** The reverse-string symmetry is about the **fault count**. The order in which pages are evicted differs.

**Interview angle:** Follow-up: *"How would you compare two replacement policies fairly?"* Simulate both on real reference traces, across a range of frame counts, and plot faults against frames.

---

## Q197. [M7 | Clock properties | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. If every reference bit is 1 when a fault occurs, Clock evicts exactly the page FIFO would.
- B. Clock needs no bookkeeping on each memory reference, beyond the hardware setting the R bit.
- C. Clock always produces the same number of faults as LRU.
- D. The enhanced (R, M) version may make several passes around the circle to find the best victim class.
- E. In basic Clock, the hand moves only when a replacement is needed.

**Answer:** A, B, D, E

**Why:**
- **A:** The hand clears every bit in one sweep, comes back to its starting frame (the oldest, in hand order), and evicts it (Q186, steps 5 and 8).
- **B:** That's Clock's whole advantage over exact LRU.
- **D:** The enhanced version first looks for (0,0). If none exists, it looks for (0,1) while clearing R bits, and so on.
- **E:** Hits only set bits. The hand advances only during a fault.

**Why the other options fail:**
- C: Clock *approximates* LRU. It loses the order among pages whose bits are both 1, so counts often differ.

**Concept recap:** Clock = cheap + approximately LRU, but it degenerates to FIFO under heavy use. Many variants exist (two-handed clock, Clock-Pro).

**Trap alert:** Don't assume Clock gives LRU's answer in a trace. Run Clock's own rules.

**Interview angle:** "Why is Clock the practical choice?" It gives LRU-like behaviour from a single bit of hardware support. Follow-up: *"What does a two-handed clock add?"* A front hand clears bits and a back hand evicts. The gap between them sets how long a page has to prove it's in use.

---

## Q198. [M7 | Write-back and clean victims | Easy | MSQ | OA]

**Select all that apply.**

- A. A clean victim page can be discarded without any I/O.
- B. Code pages are usually clean, because they are read-only.
- C. A page-out daemon may write dirty pages back ahead of time, so they are already clean when chosen as victims.
- D. Dirty pages of a memory-mapped file are written to swap space.
- E. Preferring clean victims reduces page-fault service time.

**Answer:** A, B, C, E

**Why:**
- **A, B:** A clean page matches its copy on disk, so the OS can just drop it and re-read it later if needed.
- **C:** Writing pages back in advance takes the write cost off the critical path of a fault.
- **E:** See Q189: clean victims save a whole disk write.

**Why the other options fail:**
- D: Dirty **file-backed** pages are written back to **their file**. Only **anonymous** pages (heap, stack) go to swap.

**Concept recap:** Each page has a backing store: the file for file-backed pages, swap for anonymous ones. A dirty page must be written to its backing store before its frame can be reused.

**Trap alert:** "Swap holds everything that isn't in RAM" is wrong. File pages go back to their files.

**Interview angle:** "What happens to a modified `mmap`'d file page under memory pressure?" It's written back to the file, then freed. Follow-up: *"And a `malloc`'d page?"* It's written to swap, or the allocation fails if there's no swap.

---

## Q199. [M7 | Scan pollution | Medium | MCQ | OA + Interview]

A database repeatedly uses a small, **hot index**. Then a report does **one sequential scan** of a huge table, touching each page exactly once. Under plain **LRU**, what happens, and what helps?

- A. Nothing goes wrong; LRU handles scans well
- B. The scan floods memory with pages used only once, which evicts the hot index (cache pollution). Scan-resistant policies help: put newly loaded pages on a probationary or inactive list first (2Q, ARC, Linux's inactive list), or treat scan pages as most-recently-used victims
- C. LRU evicts the scan pages first, so the index survives
- D. Switching to FIFO fixes it

**Answer:** B

**Why:** LRU assumes "used recently ⇒ will be used again". A one-time scan breaks that assumption. Each scanned page is briefly the "most recent", pushing the genuinely hot pages out. Scan-resistant designs make a page **prove** it's reused, by being referenced a second time, before it can displace established pages.

**Why the other options fail:**
- A: This is LRU's best-known weakness.
- C: Under LRU the scan pages are *recent*, so they're evicted last.
- D: FIFO also lets the scan push out the index.

**Concept recap:** Sequential flooding is LRU's pathological case. Fixes: 2Q, LRU-K, ARC, active/inactive lists, or MRU for known scans.

**Trap alert:** A scan that's *larger than memory* and repeated (scanning the same big file over and over) is even worse under LRU. Every access misses, because each page is evicted just before it's needed again.

**Interview angle:** "Where does LRU fail?" Scans and loops slightly larger than the cache. Follow-up: *"How does Linux protect hot pages from a big `cp` or `grep`?"* New file pages start on the inactive list. Only a second access promotes them to the active list.

---

## Q200. [M7 | LFU trace | Medium | Numerical | OA]

**LFU** replacement with **3 frames**.
- A page's count starts at 1 when it is loaded and increases on each hit. The count is forgotten when the page is evicted.
- Ties on count are broken by evicting the page **loaded earliest**.

Reference string: **1 2 3 1 4 1 2 5 2 3**. How many faults occur?

- A. 7
- B. 6
- C. 8
- D. 10

**Answer:** A

**Why:** At each fault, evict the page with the smallest count, breaking ties by the earliest load.

**Working:** (each entry is page:count)

| # | Ref | Frames after | Result | Reason |
|---|---|---|---|---|
| 1 | 1 | 1:1 | F | |
| 2 | 2 | 1:1, 2:1 | F | |
| 3 | 3 | 1:1, 2:1, 3:1 | F | |
| 4 | 1 | **1:2**, 2:1, 3:1 | H | |
| 5 | 4 | 1:2, 3:1, 4:1 | F | min count 1: tie between 2 and 3 → 2 loaded earlier → evict 2 |
| 6 | 1 | **1:3**, 3:1, 4:1 | H | |
| 7 | 2 | 1:3, 4:1, 2:1 | F | tie between 3 and 4 → evict 3 (loaded earlier) |
| 8 | 5 | 1:3, 2:1, 5:1 | F | tie between 4 and 2 → evict 4 (loaded earlier) |
| 9 | 2 | 1:3, **2:2**, 5:1 | H | |
| 10 | 3 | 1:3, 2:2, 3:1 | F | min count is 5's → evict 5 |

Faults = **7** (hits = 3).

**Why the other options fail:**
- B: Typically from keeping a page's count after it's evicted, which wrongly protects page 2 at #7.
- C: Usually from breaking ties by the most recent load instead of the earliest.
- D: Counts every reference as a fault.

**Concept recap:** LFU victim = min(count), with ties broken by the stated rule. Whether counts survive eviction must always be stated, because it changes the answer.

**Trap alert:** Page 1 builds up count 3 early and is then "locked in" for the rest of the string. That's LFU's stale-popularity problem (Q192) in miniature.

**Interview angle:** Follow-up: *"How would decaying the counts change this trace?"* Halving every count at regular intervals would eventually let page 1 be evicted if it stopped being used.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 181 | B | 186 | A | 191 | B | 196 | A, B, C, E |
| 182 | A | 187 | C | 192 | D | 197 | A, B, D, E |
| 183 | C | 188 | A | 193 | B | 198 | A, B, C, E |
| 184 | D | 189 | D | 194 | C | 199 | B |
| 185 | A, B, C, E | 190 | A, B, C, E | 195 | D | 200 | A |

---

**Concepts tested so far (to avoid repeats):** Batches 1–9, plus: OPT definition and why it's a benchmark · one string traced under OPT, FIFO and LRU, with LRU losing to FIFO · LRU trace updating on hits · cost of exact LRU (counters vs stack) · stack algorithms, inclusion, Belady immunity · Clock trace where a hit's bit gets cleared · the second-chance rule · enhanced (R, M) victim order · dirty-weighted fault-service time · reference-bit sampling and the A/D bits · additional-reference-bits (value, not count of ones) · LFU decay and the MFU rationale · temporal vs spatial locality · row vs column traversal faults · Linux active/inactive lists and MGLRU · OPT/LRU reverse-string symmetry · Clock degenerating to FIFO, the two-handed clock · clean victims, page-out daemon, file vs swap backing · scan pollution and 2Q/ARC · LFU trace with tie-breaking
*(M7 finishes in Batch 11 with working sets, thrashing, PFF, frame allocation, global vs local replacement, copy-on-write, swapping, mmap, and page-size trade-offs.)*

**Completed: 200 / 300    Remaining: 100**

---
---

# 📊 Progress Report #2 (after Q200)

## 1. Topics completed

| Module | Status | Questions |
|---|---|---|
| M1 Fundamentals & architecture | ✅ | Q1–18 |
| M2 Processes, threads, IPC | ✅ | Q19–48 |
| M3 CPU scheduling | ✅ | Q49–76 |
| M4 Synchronization & concurrency | ✅ | Q77–114 |
| M5 Deadlocks | ✅ | Q115–134 |
| M6 Memory management | ✅ | Q135–169 |
| M7 Virtual memory & page replacement | ⏳ 31 of 40 done | Q170–200 (finishes Q201–209) |
| M8–M12 | ⬜ | Q210–300 |

## 2. Key concepts covered in Q101–Q200

- **Synchronization (finished):** monitors; condition variables (atomic release-and-sleep, `while` loops, lost signals, Hoare vs Mesa); bounded buffer and the deadlock caused by swapping the waits; readers-writers (reader preference, upgrade deadlock); dining philosophers; barriers (and why the one-time version isn't reusable); lock ordering; reentrant locks.
- **Deadlocks:** the four Coffman conditions; RAGs, reduction, and wait-for graphs; deadlocks that aren't about locks (communication, pipes); prevention for each condition; safe vs unsafe; the Banker's safety and request algorithms; detection with the Request matrix; recovery; livelock; the ostrich approach; n(k − 1) + 1.
- **Memory management:** the hierarchy; logical vs physical addresses; binding times; base and limit; partitioning; placement algorithms; the 50-percent rule; fragmentation; compaction; paging arithmetic; page-table sizing; PTE protection bits; segmentation (and on x86-64); multilevel and inverted page tables; TLB and EMAT; shared code; the buddy allocator.
- **Virtual memory so far:** demand paging; page-fault handling; the meanings of the valid bit; EAT with page faults; restarting instructions; minor vs major faults; FIFO, OPT, LRU, Clock, enhanced Clock, LFU/MFU; dirty and reference bits; locality; Linux's reclaim lists; scan pollution.

## 3. Misconceptions tested in Q101–Q200

| # | Misconception | Where |
|---|---|---|
| 34 | `if` is fine around `cond_wait` | Q103 |
| 35 | A condition-variable signal is remembered like a semaphore's | Q104 |
| 36 | The order of the two waits in producer-consumer doesn't matter | Q106 |
| 37 | `full + empty = N` at every instant | Q107 |
| 38 | RW locks are always faster than a mutex | Q109 |
| 39 | Deadlock-free ⇒ starvation-free (dining philosophers) | Q111 |
| 40 | Consistent-looking `lock(from); lock(to)` can't deadlock | Q113 |
| 41 | Starvation is a Coffman condition | Q115 |
| 42 | A cycle in a RAG ⇒ deadlock, even with multiple instances | Q116–117 |
| 43 | Unsafe = deadlocked | Q123 |
| 44 | The safe sequence is unique | Q124 |
| 45 | Memorized textbook safe sequences | Q125 |
| 46 | Request ≤ Available ⇒ grant | Q126 |
| 47 | "Someone can finish ⇒ no deadlock" | Q128 |
| 48 | The minimum for no deadlock is n(k − 1), without the + 1 | Q134 |
| 49 | Compaction fixes internal fragmentation | Q146 |
| 50 | A 1 KB page has 1000 bytes | Q149 |
| 51 | PTE width comes from the virtual address space | Q150 |
| 52 | offset = limit is a valid segment offset | Q156 |
| 53 | Linux relies on segmentation | Q158 |
| 54 | n levels ⇒ n memory accesses (forgetting the data access) | Q162 |
| 55 | Inverted page tables are faster | Q163 |
| 56 | TLB miss = page fault | Q165 |
| 57 | Adjacent same-size free blocks are buddies | Q169 |
| 58 | Every page fault means disk I/O, and every fault means two transfers | Q178–179 |
| 59 | Hits reorder the FIFO queue | Q175 |
| 60 | More frames ⇒ fewer faults, always | Q176, Q196 |
| 61 | LRU always beats FIFO | Q182 |
| 62 | LRU can show Belady's anomaly | Q185 |
| 63 | A hit permanently protects a page under Clock | Q186 |
| 64 | Clean pages outrank recently used pages as victims | Q188 |
| 65 | Count the 1-bits in reference histories | Q191 |
| 66 | Dirty file pages go to swap | Q198 |

## 4. Formulas introduced in Q101–Q200

| Formula | Where |
|---|---|
| Banker's: Need = Max − Allocation; Work = Available; loop: Need ≤ Work → Work += Allocation | Q124–125 |
| Request algorithm: ≤ Need → ≤ Available → pretend → safe? | Q126 |
| Deadlock-free minimum: R ≥ n(k − 1) + 1 | Q134 |
| Base/limit: valid ⇔ L < limit; PA = base + L | Q138 |
| Segments: valid ⇔ d < limit[s]; PA = base[s] + d | Q156 |
| Paging split: offset = log₂ P bits; page = VA − offset bits | Q148 |
| PA = PT[⌊L/P⌋] × P + (L mod P) | Q149 |
| PT size = (2^VA / P) × PTE; frame bits = log₂(PM / P) | Q150 |
| Logical size = #PTE × P; max physical = 2^frame-bits × P | Q154 |
| Pages = ⌈S/P⌉; internal fragmentation = ⌈S/P⌉P − S; average ≈ P/2 | Q146, Q159 |
| 50-percent rule → ⅓ of memory unusable | Q145 |
| Multilevel table memory = outer + (#inner needed) × size | Q161 |
| n-level walk on a TLB miss = n + 1 memory accesses | Q162 |
| EMAT = h(t + m) + (1 − h)(t + (n + 1)m) | Q164, Q166 |
| Shared code: code + users × private data | Q167 |
| Buddy: block = 2^⌈log₂ r⌉; buddy = X XOR 2ᵏ | Q168–169 |
| EAT = (1 − p)·ma + p·service; p_max ≈ allowed extra time / service | Q173–174 |
| Fault service = (1 − d)·t_clean + d·t_dirty | Q189 |
| Barrier turnstile: S₀ − #P + #V | Q112 |

## 5. Areas to revise before moving on

1. **Banker's (Q125–126):** redo both from scratch. Compute Need, then run the lowest-index loop, then the three-step request check. Re-verify that Work ends equal to Total.
2. **Page-table sizing (Q150, Q154, Q161):** practise keeping the *number of entries* (from VA) separate from the *entry width* (from PA). Then do one multilevel sizing problem.
3. **EMAT conventions (Q164, Q166, Q173):** serial vs parallel TLB lookup, and the number of levels. Always convert ms to ns first.
4. **Replacement traces (Q175, Q182–183, Q186, Q200):** one table per algorithm. Remember: FIFO doesn't change on hits, LRU updates on hits, Clock only *sets* bits on hits, and LFU needs its count-retention rule stated.
5. **Synchronization patterns (Q103–106):** write the bounded buffer twice from memory, once with semaphores and once with a mutex and condition variables, including the `while` loops.
6. **Interview-ready pairs:** safe vs unsafe vs deadlocked; prevention vs avoidance vs detection; deadlock vs livelock vs starvation; paging vs segmentation; TLB miss vs page fault; minor vs major fault; Clock vs LRU.
