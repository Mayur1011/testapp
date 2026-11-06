# OS MCQ/MSQ Bank — Batch 8 (Q141–Q160)

**What this batch covers (M6, memory management):**
- Fixed vs variable partitioning.
- First / best / worst / next fit, compared on one trace.
- The 50-percent rule.
- Internal vs external fragmentation, and compaction.
- Paging: address split, translation, page-table sizing, PTE protection bits, frame counts, logical address-space size, internal fragmentation of the last page.
- Segmentation: translation and traps, segmentation vs paging, segmentation on x86-64.

> Units in this batch: **KB = 2¹⁰ bytes, MB = 2²⁰ bytes, GB = 2³⁰ bytes**, unless a question explicitly uses decimal (as a trap).

---

## Q141. [M6 | Fixed partitioning | Easy | MCQ | OA]

Memory is divided into fixed partitions of 100 KB each. A 70 KB process is loaded into one of them. The remaining 30 KB of that partition sits unused and can't be given to anyone else. This waste is called:

- A. External fragmentation
- B. Thrashing
- C. Internal fragmentation
- D. Compaction

**Answer:** C

**Why:** **Internal fragmentation** is wasted space *inside* an allocated region, because the allocation unit is bigger than the request. The 30 KB belongs to this process's partition, so nobody else can use it.

**Why the other options fail:**
- A: External fragmentation is free memory *between* allocations that's too scattered to use.
- B: Thrashing means excessive paging (covered in M7).
- D: Compaction is a fix for external fragmentation, not a kind of waste.

**Concept recap:**
- **Internal** fragmentation = allocated but unused, *inside* a block. Comes from fixed-size units.
- **External** fragmentation = free but unusable, *between* blocks. Comes from variable-size units.

**Trap alert:** Ask yourself "who owns the wasted space?" If it's owned by the process, the fragmentation is internal. If it's free but too scattered, it's external.

**Interview angle:** "Internal vs external fragmentation, with one example each." Follow-up: *"Which does paging have?"* Internal only, in the last page (Q146, Q159).

---

## Q142. [M6 | Variable partitioning | Easy | MCQ | Interview]

Compared with fixed partitions, **variable (dynamic) partitioning**, where each process gets exactly the size it asks for:

- A. Largely eliminates internal fragmentation, but suffers external fragmentation as holes form between processes
- B. Eliminates both kinds of fragmentation
- C. Eliminates external fragmentation but causes internal fragmentation
- D. Requires paging hardware

**Answer:** A

**Why:** Each process gets exactly what it needs, so nothing is wasted inside allocations. But as processes come and go, they leave **holes** of varying sizes. Eventually there may be enough free memory in total, yet no single hole is big enough for the next request.

**Why the other options fail:**
- B: External fragmentation is the defining problem of variable partitioning.
- C: Backwards.
- D: Base and limit registers are enough (Q139).

**Concept recap:** Variable partitioning = a list of holes plus a placement policy (Q143). Allocation is usually rounded up to some small block size, which leaves a little internal waste, hence "largely".

**Trap alert:** "No internal fragmentation ⇒ efficient" is wrong. External fragmentation can waste up to a third of memory (Q145).

**Interview angle:** "What problem does variable partitioning create?" External fragmentation. Follow-up: *"Name two ways to fight it."* Compaction (Q147), and non-contiguous allocation: paging or segmentation.

---

## Q143. [M6 | First / best / worst / next fit | Hard | Numerical | OA]

The free holes, in **address order**, are: **150, 450, 250, 350, 550** KB. Requests arrive in this order: **230, 420, 140, 400** KB. Each request is carved from the start of the chosen hole, and the leftover remains a hole in the same position. A request that fits nowhere waits. Which algorithm(s) can place **all four** requests?

- A. First fit and next fit
- B. Best fit only
- C. Worst fit only
- D. First fit, best fit and worst fit all place all four

**Answer:** B

**Why:**
- **First fit** takes the first hole big enough.
- **Best fit** takes the smallest hole big enough.
- **Worst fit** takes the largest hole.
- **Next fit** works like first fit, but resumes the search where the last one stopped.

**Working:** (holes listed in address order after each step)

**First fit**

| Request | Chosen hole | Holes after |
|---|---|---|
| 230 | 450 (first ≥ 230) → leaves 220 | 150, 220, 250, 350, 550 |
| 420 | 550 → leaves 130 | 150, 220, 250, 350, 130 |
| 140 | 150 → leaves 10 | 10, 220, 250, 350, 130 |
| 400 | none ≥ 400 (largest is 350) | **waits ✗** |

**Best fit**

| Request | Chosen hole | Holes after |
|---|---|---|
| 230 | 250 (smallest ≥ 230) → leaves 20 | 150, 450, 20, 350, 550 |
| 420 | 450 → leaves 30 | 150, 30, 20, 350, 550 |
| 140 | 150 → leaves 10 | 10, 30, 20, 350, 550 |
| 400 | 550 → leaves 150 | **all placed ✓** |

**Worst fit**

| Request | Chosen hole | Holes after |
|---|---|---|
| 230 | 550 (largest) → leaves 320 | 150, 450, 250, 350, 320 |
| 420 | 450 → leaves 30 | 150, 30, 250, 350, 320 |
| 140 | 350 → leaves 210 | 150, 30, 250, 210, 320 |
| 400 | none ≥ 400 (largest is 320) | **waits ✗** |

**Next fit:**
- 230 → 450 (leaves 220).
- 420 → continues past it: 220 ✗, 250 ✗, 350 ✗, 550 ✓ (leaves 130).
- 140 → 130 ✗, wraps around to 150 ✓ (leaves 10).
- 400 → nothing is big enough → **waits ✗**. Same outcome as first fit here.

**Why the other options fail:**
- A: Both leave the 550 hole broken up before the 400 KB request arrives.
- C: Worst fit wastes the largest hole on the first, small request.
- D: Only best fit keeps a large hole (550) intact for the big final request.

**Concept recap:** No placement algorithm is best on every input. This trace happens to favour best fit, because it saves big holes for big requests.

**Trap alert:** In first fit, scan **from the beginning every time**. Resuming from where the last search stopped is next fit.

**Interview angle:** "Compare the placement algorithms." Follow-up: *"Which is fastest?"* First fit (or next fit), because the search stops early. Best and worst fit scan the whole list, unless the holes are kept sorted by size.

---

## Q144. [M6 | Placement algorithm properties | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. First fit is usually fastest, because the search stops at the first hole that is big enough.
- B. Best fit tends to leave tiny leftover holes that are too small to use.
- C. Worst fit leaves larger leftovers, but tends to use up the large holes that big requests need.
- D. Simulations show first fit and best fit generally beat worst fit in both speed and storage utilization.
- E. Next fit always produces less fragmentation than first fit.

**Answer:** A, B, C, D

**Why:**
- **B:** Best fit minimizes the leftover *each time*, and over time that adds up to many slivers that are too small to use.
- **C:** Worst fit avoids slivers by keeping leftovers big, but it destroys large holes that big requests need (Q143).
- **D:** This is the classic simulation result reported in Silberschatz.

**Why the other options fail:**
- E: Next fit spreads allocations across all of memory, and it often fragments the large hole at the end. It's sometimes *worse* than first fit, and certainly not always better.

**Concept recap:** First fit is simple and fast. Best fit leaves the smallest leftover each time. Worst fit leaves the largest. None of them eliminates external fragmentation.

**Trap alert:** "Best fit" is best by its *local* goal (smallest leftover), not necessarily for overall utilization.

**Interview angle:** "Which placement algorithm would you choose for an allocator?" First fit is a strong default. Modern allocators use **segregated free lists**: separate lists for different size classes. Follow-up: *"How does malloc avoid scanning one huge list?"* Size-class bins and per-thread caches.

---

## Q145. [M6 | The 50-percent rule | Hard | MCQ | OA + Interview]

Statistical analysis of **first fit** shows that when N blocks are allocated, about **0.5N** blocks are lost to external fragmentation. Roughly what fraction of memory may then be unusable?

- A. 1/2
- B. 1/4
- C. 2/3
- D. 1/3

**Answer:** D

**Why:** Treat the allocated blocks and the fragment holes as roughly the same average size. The total is then N allocated + 0.5N fragments = 1.5N units, so the unusable fraction is 0.5N ÷ 1.5N = **1/3**.

**Why the other options fail:**
- A: 0.5N/N, which divides by the *allocated* part only instead of by all of memory.
- B, C: Arithmetic slips.

**Concept recap:** **50-percent rule:** fragments ≈ 0.5 × allocated blocks ⇒ about 1/3 of memory may be unusable.

**Trap alert:** The "50%" in the name refers to blocks lost per allocated block, not to the fraction of memory wasted. The memory fraction is 1/3.

**Interview angle:** "How bad can external fragmentation get?" Quote this rule. Follow-up: *"Why did paging become the standard?"* It removes external fragmentation completely, because any free frame can hold any page.

---

## Q146. [M6 | Internal vs external fragmentation across schemes | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Paging eliminates external fragmentation, but causes internal fragmentation in a process's last page.
- B. Segmentation suffers from external fragmentation.
- C. Fixed partitioning causes internal fragmentation.
- D. Compaction fixes internal fragmentation.
- E. With page size P, the average internal fragmentation is about P/2 per process (per contiguous region).

**Answer:** A, B, C, E

**Why:**
- **A:** All frames are the same size, so any free frame can be used and there's no external fragmentation. But a process rarely ends exactly on a page boundary, so part of its last page is wasted.
- **B:** Segments have variable sizes, which brings back the hole problem.
- **E:** The wasted part of the last page is anywhere from 0 to P − 1 bytes, so it averages about P/2.

**Why the other options fail:**
- D: Compaction merges scattered *free* holes, which is external fragmentation. It can't recover space wasted *inside* allocations.

**Concept recap:**

| Scheme | Internal fragmentation | External fragmentation |
|---|---|---|
| Fixed partitions | yes | no (but the number of partitions limits multiprogramming) |
| Variable partitions | ≈ no | yes |
| Paging | yes (last page) | no |
| Segmentation | ≈ no | yes |

**Trap alert:** Paging *trades* external fragmentation for a small amount of internal fragmentation. It doesn't eliminate both.

**Interview angle:** "What fragmentation does paging have?" Follow-up: *"How does page size affect it?"* Bigger pages waste more in the last page, but need fewer page-table entries (covered in M7).

---

## Q147. [M6 | Compaction | Medium | MCQ | OA]

Compaction (shuffling memory contents so that all free space forms one large block) is possible only if:

- A. Relocation is dynamic (execution-time binding), so a moved process keeps working once its base register or mapping is updated
- B. Processes are compiled with absolute addresses
- C. Memory uses fixed partitions
- D. The TLB is disabled

**Answer:** A

**Why:** Compaction physically moves processes. If their addresses were fixed into the code (compile- or load-time binding), every internal pointer would now point to the wrong place. With execution-time binding (Q137), the OS copies the memory and updates the base register (or mapping), and the process doesn't notice.

**Why the other options fail:**
- B: Absolute addresses are exactly what makes moving a process impossible.
- C: Compaction is for **variable** partitions, which have external holes.
- D: TLB entries would need flushing after a move, but disabling the TLB isn't required.

**Concept recap:** Compaction = move processes to one end and merge the holes. It needs dynamic relocation, and it costs time proportional to the amount of memory copied.

**Trap alert:** A process with I/O in progress, where a device is DMA-ing into its buffer, can't simply be moved. Either pin it in place, or do the I/O through kernel buffers.

**Interview angle:** "Why don't modern OSes compact all of RAM?" Paging made it unnecessary for external fragmentation. Follow-up: *"Does Linux ever compact memory?"* Yes. It migrates pages to assemble physically contiguous runs for huge pages and contiguous allocations.

---

## Q148. [M6 | Splitting a paged address | Easy | MCQ | OA]

With a **32-bit** logical address and **4 KB** pages, how many bits are used for the page number and for the offset?

- A. 22 and 10
- B. 12 and 20
- C. 20 and 12
- D. 16 and 16

**Answer:** C

**Why:** The offset must be able to address every byte within a page: 4 KB = 2¹² bytes, so the offset is 12 bits. The remaining high-order bits form the page number: 32 − 12 = **20** bits, giving 2²⁰ pages.

**Why the other options fail:**
- A: Would mean 1 KB pages (10-bit offset).
- B: Swaps the two fields.
- D: Would mean 64 KB pages.

**Concept recap:** For page size 2ⁿ bytes and an m-bit address: **offset = n bits, page number = m − n bits.**

**Trap alert:** Because the page size is a power of two, splitting the address is just taking bits, with no division needed (Q160 A).

**Interview angle:** Usually the first step of a bigger numerical. Follow-up: *"How many page-table entries does that imply?"* 2²⁰, about 1 million per process (Q150).

---

## Q149. [M6 | Paged address translation | Medium | Numerical | OA + Interview]

Page size = **1024 bytes**. A process's page table is:

| Page | 0 | 1 | 2 | 3 | 4 | 5 |
|---|---|---|---|---|---|---|
| Frame | 5 | 2 | 7 | 1 | 6 | 3 |

What physical address does logical address **4500** (decimal) map to?

- A. 6500
- B. 6548
- C. 4500
- D. 3476

**Answer:** B

**Why:** Split the logical address into page number and offset, look up the frame for that page, then rebuild the address with the **same offset** inside the frame.

**Working:**
- **Formulas:** page = ⌊L / P⌋; offset = L mod P; physical = frame × P + offset.
- page = ⌊4500 / 1024⌋ = **4** (since 4 × 1024 = 4096).
- offset = 4500 − 4096 = **404**.
- Page 4 → frame **6**.
- Physical = 6 × 1024 + 404 = 6144 + 404 = **6548**.

**Why the other options fail:**
- A: Uses 1000 instead of 1024 (page 4, offset 500 → 6 × 1000 + 500). This is the decimal-vs-binary units trap.
- C: Skips translation entirely.
- D: Uses frame 3, which belongs to page 5: 3 × 1024 + 404 = 3476. Catches misreading the table.

**Concept recap:** physical = PT[⌊L/P⌋] × P + (L mod P). The offset never changes during translation.

**Trap alert:** A "1 KB" page is **1024** bytes. Also double-check the page you look up: an off-by-one row read is a common slip.

**Interview angle:** "Walk me through translating a virtual address." Split it, index the page table, combine frame and offset. Follow-up: *"What if page 4's entry were marked invalid?"* A page fault (covered in M7).

---

## Q150. [M6 | Page-table size | Hard | Numerical | OA]

- Virtual address: **40 bits**. Page size: **8 KB**. Physical memory: **32 GB**.
- Each PTE (page-table entry) holds the frame number plus **4 extra bits** (valid, dirty, referenced, protection), **rounded up to a whole number of bytes**.
- Page tables are single-level.

What is the size of one process's page table?

- A. 384 MB
- B. 416 MB
- C. 512 MB
- D. 16 MB

**Answer:** C

**Why:** A single-level page table has one PTE for every virtual page. Its size is (number of pages) × (PTE size).

**Working:**
1. Number of pages = 2⁴⁰ / 2¹³ = **2²⁷** entries.
2. Number of frames = 32 GB / 8 KB = 2³⁵ / 2¹³ = 2²² → the frame number needs **22 bits**.
3. PTE = 22 + 4 = 26 bits → rounded up to whole bytes = **4 bytes** (3 bytes would hold only 24 bits).
4. Page-table size = 2²⁷ × 4 B = 2²⁹ B = **512 MB**.

**Why the other options fail:**
- A: Uses a 3-byte PTE (2²⁷ × 3 = 384 MB), which is too small to hold 26 bits.
- B: Uses the exact 26 bits without rounding up (≈ 416 MB).
- D: Not obtained from any consistent combination of these values.

**Concept recap:**
- PT size = (2^VA / page size) × PTE size.
- Frame-number bits = log₂(physical memory / page size).

**Trap alert:** Don't confuse the **number of entries**, which comes from the *virtual* address space, with the **width of each entry**, which comes from the *physical* frame count.

**Interview angle:** "Why is a 512 MB page table per process a problem?" It's far too big to keep in memory for every process. Follow-up: *"Fixes?"* Multilevel page tables, inverted page tables, or larger pages (Batch 9).

---

## Q151. [M6 | Page tables in the system | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Paging lets a process's physical memory be non-contiguous.
- B. Each process typically has its own page table.
- C. The page-table base register (PTBR; CR3 on x86) is reloaded on a context switch between processes.
- D. With paging, a user process can access any physical frame by computing a large enough logical address.
- E. With a single-level page table and no TLB, every data access needs two memory accesses: one for the PTE and one for the data.

**Answer:** A, B, C, E

**Why:**
- **A:** Pages can be placed in *any* free frames.
- **B, C:** Each process has its own address space, so the kernel switches to that process's page table when it switches processes (this links to Q11).
- **E:** The PTE itself lives in memory, so finding the frame costs one memory access before the data access. That doubling is why the TLB exists (Batch 9).

**Why the other options fail:**
- D: A process can reach only frames **mapped in its own page table**. Addresses beyond that are invalid and cause faults. This is how paging provides isolation.

**Concept recap:** Paging = non-contiguous physical placement + per-process mappings + isolation. Without a TLB, effective access ≈ 2 × the memory access time.

**Trap alert:** Large virtual addresses don't let you escape your page table. The MMU translates everything through it.

**Interview angle:** "How does paging isolate processes?" Each process has its own page table, and only the kernel can change the PTBR. Follow-up: *"How do two processes share memory, then?"* Their page tables map pages to the same frame (shared pages, Batch 9).

---

## Q152. [M6 | Protection bits in a PTE | Medium | MCQ | Interview]

A process performs a store to an address whose page is mapped **read-only** in its page table. No copy-on-write is involved. What happens?

- A. The write succeeds but is silently discarded
- B. The page is silently made writable
- C. The CPU writes to a different frame instead
- D. The MMU raises a protection fault; the kernel finds it is a real violation and sends the process SIGSEGV (on Linux)

**Answer:** D

**Why:** PTEs carry **protection bits**. The MMU checks them on every access, alongside translation. A write to a read-only page traps to the kernel. The kernel checks whether it should handle the fault itself (for example copy-on-write, covered in M7). If not, the access is illegal, and on Linux the process gets SIGSEGV.

**Why the other options fail:**
- A, B, C: The hardware never silently changes or redirects a protected access.

**Concept recap:** Common PTE bits:
- **valid/present**
- **R/W** (read-only vs writable)
- **U/S** (user vs supervisor, i.e. whether user mode may access it)
- **NX** (no-execute)
- plus the accessed and dirty bits (covered in M7)

**Trap alert:** The *same* write-protection fault is also used deliberately, for copy-on-write after `fork`. A protection fault isn't always an error.

**Interview angle:** "How does the OS make code pages read-only and the stack non-executable?" With the R/W and NX bits in the PTEs. Follow-up: *"What does W^X mean?"* No page should be both writable and executable at the same time, which blocks classic code-injection attacks.

---

## Q153. [M6 | Frames and physical address bits | Medium | MCQ | OA]

Physical memory is **64 MB** and the frame size is **4 KB**. How many frames are there, and how many bits are in a physical address?

- A. 2¹⁴ frames; 26 bits
- B. 2¹⁶ frames; 26 bits
- C. 2¹⁴ frames; 24 bits
- D. 2¹² frames; 32 bits

**Answer:** A

**Why:**
- 64 MB = 2²⁶ bytes, so a physical address needs **26 bits**.
- Frames = 2²⁶ / 2¹² = **2¹⁴** = 16,384.
- Check: a physical address = 14-bit frame number + 12-bit offset = 26 bits ✓.

**Why the other options fail:**
- B: 2¹⁶ frames would require 2²⁸ bytes (256 MB).
- C: The frame count is right, but 24 bits addresses only 16 MB.
- D: Assumes a 32-bit physical address without computing it.

**Concept recap:** PA bits = log₂(physical memory). Frames = memory / frame size. Frame-number bits = PA bits − offset bits.

**Trap alert:** The virtual-address width and the physical-address width are **independent**. A 32-bit virtual address doesn't imply 32-bit physical addresses.

**Interview angle:** Usually part of a bigger sizing problem. Follow-up: *"How wide must the PTE's frame field be here?"* 14 bits.

---

## Q154. [M6 | Logical address space from page-table parameters | Hard | MCQ | OA]

A process's page table has **2¹⁸ entries**, one per page of its logical address space. The page size is **2 KB**, and each PTE's frame-number field is **16 bits**. What is the size of the logical address space, and what is the maximum physical memory the system can address?

- A. 256 MB; 64 MB
- B. 512 MB; 128 MB
- C. 1 GB; 128 MB
- D. 512 MB; 2¹⁶ bytes

**Answer:** B

**Why:** The logical size follows from the number of pages; the physical size follows from the number of frames the PTE can name.
- **Logical:** 2¹⁸ pages × 2¹¹ bytes = 2²⁹ bytes = **512 MB** (29-bit logical addresses).
- **Physical:** a 16-bit frame number names 2¹⁶ frames × 2¹¹ bytes = 2²⁷ bytes = **128 MB** (27-bit physical addresses).

**Why the other options fail:**
- A: Uses 1 KB pages for both calculations.
- C: Logical size is off by a factor of 2.
- D: Treats 2¹⁶ as a byte count, forgetting to multiply by the frame size.

**Concept recap:**
- Logical size = (#PTEs) × page size.
- Max physical size = 2^(frame-field bits) × frame size.

**Trap alert:** The virtual address space here (512 MB) is *bigger* than physical memory (128 MB). That's perfectly normal, and it's exactly the situation virtual memory handles (covered in M7).

**Interview angle:** Follow-up: *"Can a process use all 512 MB at once?"* Not all in RAM at the same time. Demand paging keeps only the needed pages in frames.

---

## Q155. [M6 | Segmentation basics | Easy | MCQ | Interview]

Which statement describes **segmentation**?

- A. Memory is divided into fixed-size blocks that are invisible to the programmer
- B. A logical address is a page number plus an offset
- C. The logical address space is a set of variable-sized segments (code, data, stack, …), and an address is the pair ⟨segment number, offset⟩
- D. Segmentation eliminates external fragmentation

**Answer:** C

**Why:** Segmentation matches the **programmer's view** of a program as separate logical units: the main code, globals, the heap, the stack, libraries. Each unit is a segment with its own base and length, recorded in a **segment table**.

**Why the other options fail:**
- A, B: That's paging.
- D: Variable-sized segments bring back external fragmentation (Q146).

**Concept recap:** Segment table entry = ⟨base, limit⟩ (plus protection bits). An address ⟨s, d⟩ is valid if d < limit[s], and then physical = base[s] + d.

**Trap alert:** A segment number isn't a page number. Segments vary in size; pages are all the same size.

**Interview angle:** "Paging vs segmentation?" Follow-up: *"Which matches how programmers think?"* Segmentation, which is why it makes sharing and protection natural (Q157).

---

## Q156. [M6 | Segmented address translation | Medium | MCQ | OA]

The segment table is:

| Segment | Base | Limit |
|---|---|---|
| 0 | 2000 | 600 |
| 1 | 5200 | 300 |
| 2 | 900 | 1100 |

Which option gives correct translation result(s)?

- A. ⟨1, 300⟩ → 5500
- B. ⟨0, 580⟩ → 2580, and ⟨2, 1050⟩ → 1950
- C. ⟨2, 1050⟩ causes a trap
- D. ⟨0, 580⟩ → 580

**Answer:** B

**Why:** Check offset < limit first. If it holds, add the base. Otherwise trap.
- ⟨0, 580⟩: 580 < 600 ✓ → 2000 + 580 = **2580**.
- ⟨2, 1050⟩: 1050 < 1100 ✓ → 900 + 1050 = **1950**.
- ⟨1, 300⟩: 300 < 300 ✗ → **trap**. Valid offsets are 0 … 299.

**Why the other options fail:**
- A: An offset equal to the limit is out of range.
- C: 1050 is within segment 2's limit.
- D: Forgets to add the base.

**Concept recap:** valid ⇔ 0 ≤ d < limit; physical = base + d. Segment bases don't have to be in any order (segment 2 sits *below* segment 0 here).

**Trap alert:** **offset == limit is invalid.** This off-by-one shows up constantly in exams.

**Interview angle:** Usually an OA item. Follow-up: *"What happens to the process on that trap?"* The OS treats it as an addressing error, and on Unix the process typically gets SIGSEGV.

---

## Q157. [M6 | Segmentation vs paging | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. Segmentation matches the programmer's view of a program; paging is invisible to the programmer.
- B. Segmentation suffers from external fragmentation; paging suffers from internal fragmentation.
- C. Sharing and protection fit naturally at segment granularity (e.g., share the whole code segment read-only).
- D. Segments are all the same fixed size.
- E. Paged segmentation, as in MULTICS and legacy x86, combines the two: each segment is itself divided into pages.

**Answer:** A, B, C, E

**Why:**
- **A, C:** Segments correspond to meaningful units, so "share this library" or "make this read-only" is a single segment-table setting.
- **B:** Variable sizes cause holes; fixed sizes cause waste in the last page.
- **E:** Paged segmentation keeps logical segments *and* avoids external fragmentation, because each segment is stored in pages.

**Why the other options fail:**
- D: Variable size is the whole point of segments.

**Concept recap:**

| | Segmentation | Paging |
|---|---|---|
| Unit size | variable | fixed |
| Visible to programmer | yes | no |
| Fragmentation | external | internal |
| Sharing/protection unit | logical segment | page |

**Trap alert:** "Paging can't share code" is false. It can, by mapping the same frames in several page tables. It's just less natural than sharing one segment.

**Interview angle:** "Compare paging and segmentation." Use the table above. Follow-up: *"Why did paging win?"* Fixed-size frames make allocation trivial and remove external fragmentation.

---

## Q158. [M6 | Segmentation on x86-64 today | Medium | MCQ | OA + Interview]

On x86-64 Linux, how is segmentation used for ordinary user code?

- A. Every function gets its own segment
- B. Segmentation provides all the protection, and paging is disabled
- C. Segments are variable-sized and swapped individually
- D. Almost not at all: code and data segment bases are 0 with no effective limits (a "flat" model), paging handles translation and protection, and only the FS/GS bases are used, e.g. to point at thread-local storage

**Answer:** D

**Why:** In 64-bit mode, the CPU treats the bases of CS, DS, ES and SS as 0 and doesn't check segment limits. Every address is effectively flat, so segmentation does almost nothing. **Paging** does all the real work. FS and GS still have usable base registers, which OSes use to locate per-thread data (TLS, thread-local storage) or per-CPU data.

**Why the other options fail:**
- A: No such scheme exists.
- B: Paging is mandatory in x86-64 long mode.
- C: That's classic textbook segmentation, not what x86-64 does.

**Concept recap:** Modern x86-64 = "flat" segmentation + 4-level paging (5-level on newer CPUs).

**Trap alert:** Don't conclude "x86 uses segmentation, so Linux does". The segment registers still exist, but they're mostly vestigial.

**Interview angle:** "Does Linux use segmentation?" Essentially no, beyond FS/GS for TLS. Follow-up: *"How is `thread_local` accessed so cheaply?"* Through an FS-relative address on x86-64.

---

## Q159. [M6 | Internal fragmentation in the last page | Medium | MCQ | OA]

Page size = **4096 bytes**. A process needs **72,766 bytes**. How many pages does it need, and how much internal fragmentation is there?

- A. 18 pages; 962 bytes
- B. 17 pages; 3134 bytes
- C. 18 pages; 3134 bytes
- D. 17 pages; 962 bytes

**Answer:** A

**Why:** The process needs a whole number of pages, so round up. The waste is whatever is left unused in the last page.
- 72,766 / 4096 ≈ 17.77 → **18 pages**.
- 17 full pages = 69,632 bytes, so the last page holds 72,766 − 69,632 = 3134 bytes.
- Internal fragmentation = 4096 − 3134 = **962 bytes**.

**Why the other options fail:**
- B: 17 pages can't hold the process, and 3134 is the *used* part of the last page.
- C: Reports the used part instead of the waste.
- D: Doesn't round up.

**Concept recap:** pages = ⌈S / P⌉; internal fragmentation = ⌈S / P⌉ × P − S.

**Trap alert:** Don't report the **used** part of the last page as the fragmentation. The fragmentation is the **unused** part.

**Interview angle:** Follow-up: *"Worst case and average case?"* The worst case is P − 1 bytes, when the process needs one byte more than a page boundary. The average is about P/2.

---

## Q160. [M6 | Paging address arithmetic | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. The page size is a power of two, so splitting an address into page number and offset needs no division, only taking bits.
- B. The frame size equals the page size.
- C. Physical address = frame number × page size + offset.
- D. The number of pages in a logical address space must equal the number of frames in physical memory.
- E. The offset is identical in the logical and the physical address.

**Answer:** A, B, C, E

**Why:**
- **A:** With P = 2ⁿ, the offset is the low n bits and the page number is the rest. Hardware just routes the wires.
- **B, C, E:** A page maps onto an equal-sized frame and keeps its position within it, so translation replaces only the page number with the frame number.

**Why the other options fail:**
- D: Logical and physical sizes are independent (Q154). Usually there are many more pages than frames, and virtual memory relies on that.

**Concept recap:** Translation = **replace the page number with the frame number; keep the offset.** In bits: physical = (frame << n) | offset.

**Trap alert:** Translation never "adds" the page number to anything. The page number is replaced entirely, and the offset passes through unchanged.

**Interview angle:** "Why are page sizes powers of two?" So translation is just taking bits (fast hardware), with no division. Follow-up: *"Is it the same for segments?"* No. Segments need a full addition (base + offset) and a comparison against the limit.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 141 | C | 146 | A, B, C, E | 151 | A, B, C, E | 156 | B |
| 142 | A | 147 | A | 152 | D | 157 | A, B, C, E |
| 143 | B | 148 | C | 153 | A | 158 | D |
| 144 | A, B, C, D | 149 | B | 154 | B | 159 | A |
| 145 | D | 150 | C | 155 | C | 160 | A, B, C, E |

---

**Concepts tested so far (to avoid repeats):** Batches 1–7, plus: internal fragmentation from fixed partitions · variable partitions and external fragmentation · a four-algorithm placement trace where only best fit places everything · placement properties (speed, slivers, next fit) · 50-percent rule → 1/3 · fragmentation type by scheme, P/2 average · compaction needs dynamic relocation, DMA pinning · page-number/offset split · translation with the 1000-vs-1024 trap · page-table sizing with PTE width from the frame count · per-process page tables, PTBR, the 2× access cost without a TLB · read-only PTE → protection fault → SIGSEGV, W^X · frames and PA bits · logical size and max physical memory from PT parameters · segmentation definition · segment translation with the offset = limit trap · segmentation vs paging, paged segmentation · x86-64 flat model and FS/GS · last-page internal fragmentation · power-of-two address arithmetic
*(M6 continues in Batch 9 with multilevel paging, inverted page tables, TLB and EMAT, shared pages, and the buddy system. M7 (virtual memory) starts after that.)*

**Completed: 160 / 300    Remaining: 140**
