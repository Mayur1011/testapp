# OS MCQ/MSQ Bank — Batch 11 (Q201–Q220)

**What this batch covers:**
- **M7, virtual memory completed (Q201–209):**
  - Working-set computation, and the working-set model as a cure for thrashing.
  - Thrashing, and page-fault-frequency (PFF) control.
  - Proportional frame allocation, global vs local replacement, minimum frames per process.
  - Copy-on-write, `mmap` (shared vs private, lazy loading), swapping.
  - The optimal page size.
- **M8, file systems begins (Q210–220):**
  - What an inode holds, why `open()` exists, access methods, relative paths.
  - Directory structures, and why hard links to directories are forbidden.
  - Contiguous, linked (and FAT) and indexed allocation, with disk-read counts.
  - Implementing large directories.

> Conventions: KB = 2¹⁰ bytes, MB = 2²⁰ bytes. Disk blocks are numbered from 0, and "logical block i" means the (i + 1)-th block of a file.

---

## Q201. [M7 | Working set | Medium | Numerical | OA]

Here is a reference string, with the reference times shown above each page:

| t | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| page | 1 | 2 | 3 | 2 | 1 | 4 | 4 | 4 | 5 | 4 | 6 | 6 |

The working-set window is **Δ = 4**. What are the working sets at **t = 8** and at **t = 12**?

- A. {1, 4} and {4, 5, 6}
- B. {1, 2, 3, 4} and {4, 5, 6}
- C. {4} and {4, 6}
- D. {4} and {6}

**Answer:** A

**Why:** The **working set** WS(t, Δ) is the set of **distinct pages** referenced in the last Δ references, i.e. times t − Δ + 1 through t. It estimates which pages a process needs right now.

**Working:**
- **t = 8:** the window covers t = 5…8, with references 1, 4, 4, 4 → WS = **{1, 4}**, so WSS (working-set size) = 2.
- **t = 12:** the window covers t = 9…12, with references 5, 4, 6, 6 → WS = **{4, 5, 6}**, so WSS = 3.

**Why the other options fail:**
- B: The first set uses Δ = 8 (the window t = 1…8).
- C: Uses Δ = 3 (t = 6…8 and t = 10…12).
- D: Uses Δ = 2.

**Concept recap:** WS(t, Δ) = distinct pages among the last Δ references, and WSSᵢ = |WSᵢ|. Small Δ → misses part of the current locality. Large Δ → spans several localities.

**Trap alert:** Count **distinct** pages. Four references to page 4 still contribute one page.

**Interview angle:** "What is a working set?" The pages a process is actively using in its current phase. Follow-up: *"How does an OS measure it without logging every reference?"* By sampling reference bits on a timer interrupt, like the reference-bit history in Q191.

---

## Q202. [M7 | Working-set model vs thrashing | Medium | MCQ | OA + Interview]

Five processes have working-set sizes of **30, 45, 20, 60 and 35** frames. The system has **180** frames. Under the working-set model, what should the OS do?

- A. Admit another process, because memory is nearly balanced
- B. Do nothing
- C. Total demand D = 190 > 180, so thrashing is likely: suspend (swap out) one process and give its frames to the others
- D. Shrink every process's allocation equally until they fit

**Answer:** C

**Why:** The working-set model keeps a process running only if its whole working set fits in memory. Total demand is D = ΣWSSᵢ = 30 + 45 + 20 + 60 + 35 = **190**. That exceeds the 180 frames available. Squeezing everyone means *every* process faults constantly. Suspending one process lets the rest each hold their full working set.

**Why the other options fail:**
- A: Admitting another process makes the shortage worse.
- B: Leaving it alone allows thrashing.
- D: Every process would be below its working set, which is exactly the recipe for thrashing.

**Concept recap:** If D > m (available frames), suspend a process. If D < m, there's room to admit another. This keeps the degree of multiprogramming as high as possible without thrashing.

**Trap alert:** "Fairness: shrink everyone a bit" sounds reasonable but makes things worse. It's better to run fewer processes at full speed.

**Interview angle:** "How does the working-set model prevent thrashing?" Follow-up: *"Its practical cost?"* Tracking every process's working set continuously. Real systems approximate it.

---

## Q203. [M7 | Thrashing | Easy | MCQ | Interview]

Which best describes **thrashing**?

- A. The CPU is doing too much useful work
- B. The disk sits idle while processes compute
- C. Many processes are waiting for keyboard input
- D. Processes spend more time paging than executing. Each one faults and takes frames from others (under global replacement), CPU utilization drops, and a naive long-term scheduler then admits *more* processes, making it worse

**Answer:** D

**Why:** If a process has fewer frames than its working set, it faults almost constantly. With global replacement, it steals frames from other processes, which then fault too. Everyone queues at the paging disk, so the CPU goes idle. A scheduler that treats low CPU utilization as "add more work" then admits more processes, which deepens the collapse.

**Why the other options fail:**
- A: The CPU is mostly *idle* during thrashing.
- B: Backwards. The paging disk is saturated.
- C: Waiting for keyboard input has nothing to do with paging.

**Concept recap:** On a plot of CPU utilization against the degree of multiprogramming, utilization rises up to a peak and then **collapses**. Thrashing begins when ΣWSS > m.

**Trap alert:** "Low CPU utilization ⇒ add more processes" is the wrong conclusion during thrashing.

**Interview angle:** "How do you recognize and fix thrashing?" Signs: a high fault rate, the disk busy, low CPU utilization. Fix it by lowering the degree of multiprogramming, adding RAM, or using working-set/PFF control. Follow-up: *"Can local replacement prevent it?"* Local replacement *contains* it to one process, but that process still thrashes (Q206).

---

## Q204. [M7 | Page-fault frequency (PFF) | Medium | MCQ | OA]

A system uses **page-fault-frequency** control, with an upper and a lower bound on each process's fault rate. One process's fault rate rises **above the upper bound**, and **no free frames** remain. What should the OS do?

- A. Take frames away from that process
- B. Suspend (swap out) some process and redistribute its frames
- C. Lower the upper bound
- D. Ignore it until the rate drops

**Answer:** B

**Why:** PFF works like a thermostat on each process's fault rate.
- **Above the upper bound:** the process needs more frames, so give it some.
- **Below the lower bound:** it has more than it needs, so take some away.

If the process needs frames and none are free, the only safe move is to **lower the degree of multiprogramming**: suspend a process and redistribute its frames.

**Why the other options fail:**
- A: The process is already short of frames. Taking more away makes it worse.
- C: Changing the target doesn't fix the shortage.
- D: Letting it continue risks thrashing.

**Concept recap:** PFF = measure faults directly, and adjust frames between the two bounds. It's simpler than tracking working sets.

**Trap alert:** PFF and the working-set model reach the same remedy (suspend a process) from different measurements. PFF measures the fault rate; the working-set model estimates the pages needed.

**Interview angle:** "Working set vs PFF?" Follow-up: *"Why might PFF react badly to a phase change?"* At the start of a new locality, the fault rate spikes. PFF hands out frames that may not be needed once the new phase settles.

---

## Q205. [M7 | Proportional frame allocation | Easy | MCQ | OA]

There are **90** free frames. Three processes have sizes of **20, 60 and 100** pages. Using **proportional allocation**, how many frames does each get?

- A. 10, 30, 50
- B. 30, 30, 30
- C. 20, 60, 100
- D. 9, 30, 51

**Answer:** A

**Why:** Proportional allocation gives each process frames in proportion to its size: aᵢ = (sᵢ / S) × m, where S = Σsᵢ = 180 and m = 90.
- 20/180 × 90 = **10**
- 60/180 × 90 = **30**
- 100/180 × 90 = **50**

**Why the other options fail:**
- B: That's **equal** allocation (m/n).
- C: Gives each process its full size, which needs 180 frames.
- D: Off by one at each end, with no rule behind it.

**Concept recap:** aᵢ = sᵢ/S × m, rounded to integers, and each process must get at least the minimum its instruction set requires (Q206 E). *Priority* allocation uses priority, or a mix of priority and size, instead of size alone.

**Trap alert:** When the shares aren't whole numbers, round them so the total still adds up to m.

**Interview angle:** "How should an OS divide frames among processes?" Follow-up: *"Why not divide them equally?"* A 10-page process doesn't need as many frames as a 1000-page one.

---

## Q206. [M7 | Global vs local replacement | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. With global replacement, a process can take frames from other processes.
- B. With local replacement, a process's fault rate depends only on its own behaviour.
- C. Global replacement generally gives better throughput and is the more common choice.
- D. With local replacement, thrashing in one process spreads to the others.
- E. The minimum number of frames a process needs is set by the instruction set: the largest number of pages a single instruction can touch.

**Answer:** A, B, C, E

**Why:**
- **A, C:** Global replacement treats all frames as one pool, so frames flow to whoever needs them. This tends to use memory better.
- **B:** Under local replacement, a process replaces only its own pages, so other processes can't affect its fault rate.
- **E:** If an instruction can touch k pages (the instruction itself, its operands, indirect addresses), it needs all k resident to complete. With fewer frames, it would fault forever.

**Why the other options fail:**
- D: That's the risk of **global** replacement. Local replacement keeps thrashing inside one process.

**Concept recap:**
- Global: better throughput, but processes interfere with each other.
- Local: predictable per process, but memory may be used poorly.
- Every process needs at least the ISA minimum number of frames.

**Trap alert:** Local replacement doesn't prevent thrashing. It just **contains** it to one process.

**Interview angle:** "Global or local page replacement?" Follow-up: *"How many frames does an instruction like `MOV [mem1] → [mem2]` need at minimum, if each operand can straddle two pages?"* Count the pages: the instruction itself (it may span 2) plus each operand (up to 2 each). That gives 6.

---

## Q207. [M7 | Copy-on-write | Medium | MCQ | OA + Interview]

After `fork()`, the parent and child share all pages **copy-on-write**. The child then writes to a variable on page P. What happens?

- A. The write goes to the shared frame and the parent sees it
- B. The child is killed for writing to a read-only page
- C. The write-protection fault traps to the kernel, which copies P into a new frame, maps that copy as writable in the child, and restarts the write; the parent keeps the original page
- D. Every page of the parent is copied at that moment

**Answer:** C

**Why:** With **COW**, `fork` doesn't copy any memory. It marks all writable pages **read-only** in *both* page tables and increments each frame's reference count. The first write to a shared page faults. The kernel sees it's a COW page, makes a private copy for the writer, and lets the write proceed. (If only one process still refers to the frame, the kernel can just make it writable without copying.)

**Why the other options fail:**
- A: That would break the separate address spaces that `fork` promises (Q26).
- B: The kernel recognizes a COW fault as legitimate, not a protection violation (Q152).
- D: Only the page being written is copied, and only when it's written.

**Concept recap:** COW = share until written, then copy that **one page**. That's why fork + exec is cheap: almost nothing gets copied before `exec` replaces the image.

**Trap alert:** COW is the reason "`fork` copies the address space" is still true **logically** (Q26), even though almost nothing is physically copied.

**Interview angle:** "Why is `fork()` fast on Linux?" COW plus copying only the page tables. Follow-up: *"What is `vfork()`?"* The child borrows the parent's memory outright, and the parent is suspended until the child calls `exec` or `_exit`. It's faster still, but dangerous if the child modifies memory.

---

## Q208. [M7 | mmap and swapping | Hard | MSQ | Interview]

Linux. **Select all that apply.**

- A. With `mmap(MAP_SHARED)` of a file, writes become visible to every process that maps the same file, and are eventually written back to the file.
- B. With `mmap(MAP_PRIVATE)`, writes go to private copy-on-write pages and never reach the file.
- C. `mmap` reads the entire file from disk at the moment of the call.
- D. Anonymous `mmap` memory, such as large `malloc` blocks, is backed by swap space.
- E. Modern systems swap individual pages, not whole processes.

**Answer:** A, B, D, E

**Why:**
- **A:** A shared mapping maps the **page-cache frames** themselves into each process, so everyone sees one copy, and `msync` or periodic writeback saves it to the file.
- **B:** A private mapping uses COW (Q207). Your changes stay private.
- **D:** Anonymous memory has no file behind it. It starts as zero-filled pages and is paged out to swap.
- **E:** Classic whole-process swapping is largely replaced by paging individual pages. "Swapping" now usually means paging anonymous pages out to swap.

**Why the other options fail:**
- C: `mmap` only sets up the mapping (a new VMA). Pages are loaded **lazily** through page faults, i.e. demand paging.

**Concept recap:** `mmap` turns file I/O into memory access. Shared vs private decides whether writes reach the file. File-backed pages go back to the file; anonymous pages go to swap (Q198).

**Trap alert:** "`mmap` is always faster than `read`" isn't guaranteed. For small sequential reads, the page-fault and mapping costs can outweigh the saved copy.

**Interview angle:** "`read()` vs `mmap()` for a large file?" `mmap` avoids the user–kernel copy and suits random access. `read` is simpler and can be faster for streaming. Follow-up: *"How do shared libraries use mmap?"* Each library is mapped `MAP_PRIVATE`, so its code pages are shared between processes while any writes stay private.

---

## Q209. [M7 | Optimal page size | Hard | MCQ | OA + Interview]

The average process size is **s = 1 MB** and each PTE takes **e = 8 bytes**. Model the per-process overhead as

  **overhead(p) = s·e / p (page table) + p / 2 (internal fragmentation)**,

where p is the page size. Which page size minimizes this overhead?

- A. 1 KB
- B. 2 KB
- C. 8 KB
- D. 4 KB

**Answer:** D

**Why:** Small pages mean many PTEs (a bigger page table). Large pages mean more waste in the last page. The minimum is where the two terms balance.
- Differentiate: d/dp [s·e/p + p/2] = −s·e/p² + 1/2 = 0 → **p = √(2·s·e)**.
- p = √(2 × 2²⁰ × 8) = √(2²⁴) = 2¹² = **4096 bytes = 4 KB**.
- Check: overhead(4 KB) = 2²³/2¹² + 2¹¹ = 2048 + 2048 = 4096 bytes, whereas overhead(2 KB) = 4096 + 1024 = 5120 bytes and overhead(8 KB) = 1024 + 4096 = 5120 bytes.

**Why the other options fail:**
- A, B: Too small. The page-table term dominates.
- C: Too large. The fragmentation term dominates.

**Concept recap:** **p_opt = √(2·s·e)**. At the optimum, the two overhead terms are equal.

**Trap alert:** This model counts only PT size and fragmentation. Real choices also weigh **TLB reach** (bigger pages are better), **disk I/O efficiency** (bigger is better), and **locality and resolution** (smaller is better). That's why systems offer both 4 KB and huge pages.

**Interview angle:** "Arguments for bigger vs smaller pages?" Bigger: smaller page tables, more TLB reach, fewer faults, more efficient I/O. Smaller: less internal fragmentation, finer-grained memory use. Follow-up: *"When would you use 2 MB huge pages?"* For large-memory, TLB-heavy workloads such as databases and JVM heaps.

---

## Q210. [M8 | What an inode stores | Medium | MCQ | OA]

In a Unix file system, which of these is **NOT** stored in a file's **inode**?

- A. The file's size
- B. The file's name
- C. Its permission bits and owner
- D. Its timestamps (modification, access, change)

**Answer:** B

**Why:** An **inode** (index node) holds a file's metadata: type, size, owner, permissions, timestamps, link count, and pointers to its data blocks. The **name** is stored in the **directory entry**, which maps a name to an inode number. Keeping names out of the inode is what lets one file have several names (hard links).

**Why the other options fail:** A, C and D are all standard inode fields. You can see them with `stat`.

**Concept recap:** Directory entry = ⟨name, inode number⟩. Inode = everything else about the file. On most systems, the FCB (file control block) is the inode.

**Trap alert:** "The inode is the file's record, so it holds the name" is the classic mistake. Renaming a file only changes a directory entry; the inode stays the same.

**Interview angle:** "What's in an inode?" List the fields, then point out that the name isn't there. Follow-up: *"How does `mv` within one file system work so fast?"* It just rewrites directory entries. No data moves.

---

## Q211. [M8 | Why open() exists | Easy | MCQ | Interview]

Why does a program call `open()` once, and then pass a handle to `read()` and `write()`, rather than passing the path on every call?

- A. Because `read()` can't accept a path argument for historical reasons only
- B. `open()` copies the whole file into memory
- C. `open()` searches the directories once, checks permissions, and loads the file's metadata into in-memory tables; later calls use the returned handle (fd) and skip that lookup
- D. `open()` locks the file so no other process can use it

**Answer:** C

**Why:**
- Resolving a path such as `/home/a/b/c.txt` means reading each directory along the way and checking permissions at each step. That's expensive to repeat for every 4 KB read.
- `open()` does it **once**. It puts an entry in the system-wide **open-file table** (holding the current offset and access mode, and pointing to the in-memory inode) and returns a small integer, the **fd**, that indexes the process's per-process table.

**Why the other options fail:**
- A: The handle exists for performance and to keep state such as the file offset, not merely for history.
- B: `open()` reads no file data.
- D: Unix `open()` doesn't lock by default. Locking is separate (`flock`, `fcntl`).

**Concept recap:** `open` = path lookup + permission check + table entries → fd. `close` releases them. The tables themselves are covered in Batch 12.

**Trap alert:** Permissions are checked at **open** time. If a file's permissions change afterwards, a process that already has it open usually keeps its access.

**Interview angle:** "What happens inside `open()`?" Follow-up: *"Why might a file still take up disk space after you delete it?"* A process still has it open. The data is freed only when the last link is removed **and** the last open handle is closed.

---

## Q212. [M8 | File access methods | Easy | MCQ | OA]

A database must fetch record **#5000** from a file of fixed-size records, without reading the first 4,999. Which access method does it need?

- A. Direct (relative) access: compute the record's byte position (record number × record size), work out which block holds it, and read that block
- B. Sequential access
- C. Indexed-sequential access only
- D. Append-only access

**Answer:** A

**Why:** **Direct access** treats a file as numbered blocks or records that can be read in any order. This works because block = ⌊(record number × record size) / block size⌋ is simple arithmetic. In Unix, `lseek` (or `pread`) to that offset, then read.

**Why the other options fail:**
- B: Sequential access would have to read through all 4,999 earlier records.
- C: An index helps with variable keys, but "only" is wrong. Plain direct access is enough here.
- D: Append-only access can't jump backwards.

**Concept recap:**
- Sequential: read/write next. Typical for logs and media.
- Direct: read block n. Typical for databases.
- Indexed: look up a key in an index, then access directly.

**Trap alert:** Direct access at the *file* level doesn't mean the blocks are contiguous on *disk*. How data is placed on disk is the allocation method (Q215–219).

**Interview angle:** Follow-up: *"Which allocation method makes direct access slow?"* Linked allocation without a FAT (Q219).

---

## Q213. [M8 | Directory structures | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. In a tree-structured directory, every file has exactly one path from the root.
- B. An acyclic-graph directory lets the same file or subdirectory appear in more than one directory.
- C. Unix forbids user-created hard links to directories, to keep the directory graph acyclic.
- D. A general graph with cycles makes directory traversal and deletion simple.
- E. Two-level directories (one directory per user) solve name collisions between users.

**Answer:** A, B, C, E

**Why:**
- **A:** A tree has no sharing, so each file has one unique path.
- **B:** An acyclic graph allows sharing (links) without loops.
- **C:** A hard link to a directory could create a cycle (a → b → a). Then `find` or `du` could loop forever, and reference counts could never reach zero. The only exceptions are `.` and `..`, which the system creates itself.
- **E:** With one directory per user, both alice and bob can have their own `notes.txt`.

**Why the other options fail:**
- D: Cycles make traversal loop forever unless you track visited nodes. Deletion then needs garbage collection, because counts alone fail.

**Concept recap:** single-level → two-level → tree → acyclic graph (links) → general graph (avoided).

**Trap alert:** **Symbolic** links *can* point to directories and can even form cycles. Tools cope by not following symlinks by default, or by limiting how deep they follow them (`ELOOP`).

**Interview angle:** "Why can't you hard-link a directory?" Follow-up: *"Then why can you symlink one?"* A symlink is just a path. Traversals can choose not to follow it, and it doesn't affect reference counts.

---

## Q214. [M8 | Resolving a relative path | Easy | MCQ | OA]

A process's current directory is `/home/alice`. It opens the relative path `../bob/notes.txt`. Which file does it open?

- A. `/home/alice/bob/notes.txt`
- B. `/bob/notes.txt`
- C. `/home/alice/../notes.txt`
- D. `/home/bob/notes.txt`

**Answer:** D

**Why:** A relative path starts at the **current working directory** (cwd). `..` means the parent directory.
- `/home/alice` + `..` → `/home`
- `/home` + `bob` → `/home/bob`
- `/home/bob` + `notes.txt` → `/home/bob/notes.txt`

**Why the other options fail:**
- A: Ignores the `..`.
- B: Resolves from the root instead of from the cwd.
- C: Ends up in the wrong directory.

**Concept recap:** An absolute path starts with `/`. A relative path starts from the cwd, and every process has its own cwd, inherited across `fork` (Q29).

**Trap alert:** If `/home/alice` were a *symbolic link* elsewhere, then `..` would refer to the parent of the real directory, not `/home`. Shells sometimes hide this.

**Interview angle:** Follow-up: *"Why is `cd` a shell built-in rather than a separate program?"* A child process changing *its own* cwd can't change the shell's.

---

## Q215. [M8 | Contiguous allocation | Easy | MCQ | OA + Interview]

A file is allocated **contiguously** starting at disk block **1200**, with a length of 50 blocks. The block size is **512 bytes**. Which disk block holds the byte at file offset **20,000**, and at what offset within that block?

- A. Block 1239, offset 0
- B. Block 1240, offset 32
- C. Block 1239, offset 32
- D. Block 1200, offset 20,000

**Answer:** C

**Why:** With contiguous allocation, logical block i is at disk block start + i. That's one addition, so direct access is trivial.
- Logical block = ⌊20,000 / 512⌋ = **39** (since 39 × 512 = 19,968).
- Offset within the block = 20,000 − 19,968 = **32**.
- Disk block = 1200 + 39 = **1239**.

**Why the other options fail:**
- A: Drops the offset.
- B: Off by one block, from rounding up instead of down.
- D: Treats the byte offset as if it were within one block.

**Concept recap:**
- Contiguous allocation: fast sequential *and* random access, and simple metadata (just a start and a length).
- The cost: external fragmentation, and files are hard to grow.

**Trap alert:** Also check that the block is within the file: 39 < 50 ✓. Reading past the end of the file is an error, not another block.

**Interview angle:** "Pros and cons of contiguous file allocation?" Follow-up: *"Where is it still used?"* On CD-ROMs and DVDs (write-once media, so nothing ever grows), and as the goal of **extents** in ext4, XFS and NTFS.

---

## Q216. [M8 | FAT vs plain linked allocation | Medium | MCQ | Interview]

Why does a **FAT** (File Allocation Table) file system give much faster random access than plain linked allocation?

- A. The links (each block's next-block number) are kept together in one table at the start of the volume, which can be cached in memory; reaching block i means following i table entries in memory instead of doing i disk reads
- B. FAT stores every file contiguously
- C. FAT uses a separate index block for each file
- D. FAT doesn't use links at all

**Answer:** A

**Why:** In plain linked allocation, each "next" pointer is stored **inside the data block**. Following the chain therefore means reading every block along the way. FAT moves all the pointers into one table, with one entry per disk block. That table is small enough to cache, so walking a chain is pure memory access, and only the target block needs a disk read.

**Why the other options fail:**
- B: FAT files can be scattered anywhere on the disk.
- C: That describes indexed allocation.
- D: FAT *is* linked allocation; it just keeps the links in a separate table.

**Concept recap:** FAT = linked allocation with the links moved into a cached table. The same table also tracks free blocks (an entry of 0 means free).

**Trap alert:** If the FAT isn't cached, random access is slow again. And on a huge volume the FAT becomes very large, which is one reason FAT32 limits individual files to 4 GB and newer file systems use extents.

**Interview angle:** "How does FAT work?" Follow-up: *"Why is FAT still everywhere?"* It's simple and universally supported, which is why USB drives and SD cards use FAT32 or exFAT.

---

## Q217. [M8 | Indexed allocation | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Each file has an index block holding pointers to all of its data blocks.
- B. It supports direct access without external fragmentation.
- C. It wastes space for tiny files, because even a one-block file needs a whole index block.
- D. A single index block limits the maximum file size, unless the scheme is extended (linked index blocks, multilevel indexes, or the combined inode scheme).
- E. Indexed allocation requires the data blocks to be contiguous.

**Answer:** A, B, C, D

**Why:**
- **A, B:** To find block i, read the index block and take entry i. Data blocks can sit anywhere, so there's no external fragmentation.
- **C:** Overhead: a 4 KB index block for a 100-byte file is mostly wasted.
- **D:** One index block holds (block size / pointer size) pointers, so file size is capped at that many blocks. Unix inodes fix this with direct plus indirect pointers (Batch 12).

**Why the other options fail:**
- E: Non-contiguity is the whole point of indexed allocation.

**Concept recap:** Max file size with one index block = (B / ptr) × B. For example, 4 KB blocks with 4-byte pointers give 1024 × 4 KB = 4 MB.

**Trap alert:** Indexed allocation gives fast *random* access, but sequential access may still seek all over the disk if the blocks are scattered.

**Interview angle:** "Contiguous vs linked vs indexed allocation?" Follow-up: *"How do Unix inodes avoid wasting an index block on tiny files?"* The first few pointers sit directly inside the inode.

---

## Q218. [M8 | Comparing allocation methods | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. Contiguous allocation gives the best performance for both sequential and random access.
- B. Linked allocation has no external fragmentation.
- C. With linked allocation, one damaged pointer can lose the rest of the file.
- D. FAT is a variant of linked allocation.
- E. Indexed allocation has to read the whole file to find block i.

**Answer:** A, B, C, D

**Why:**
- **A:** Contiguous blocks mean few seeks and simple arithmetic to find any block.
- **B:** Any free block can be linked in, so there's no external fragmentation.
- **C:** The chain is the only route to the later blocks, so a broken pointer cuts them off. (FAT keeps two copies of its table partly for this reason.)
- **D:** FAT is linked allocation with the pointers moved into a table (Q216).

**Why the other options fail:**
- E: Indexed allocation reads only the index block(s) and then the target block.

**Concept recap:**

| | Sequential | Random | External frag. | Growth |
|---|---|---|---|---|
| Contiguous | best | best | yes | hard |
| Linked | OK | poor (i reads) | no | easy |
| FAT | OK | good (if cached) | no | easy |
| Indexed | OK | good | no | easy (up to the limit) |

**Trap alert:** "Linked allocation is good for random access" is wrong. Without a cached FAT, it's the worst of the four.

**Interview angle:** "Which allocation method do modern file systems use?" **Extents** (contiguous runs described by a start and a length), organized in trees: ext4, XFS, NTFS, btrfs. Follow-up: *"Why extents?"* They give nearly contiguous performance, with indexed-style flexibility.

---

## Q219. [M8 | Disk reads to reach block i | Hard | Numerical | OA + Interview]

A process wants to read **logical block 150** of a file (block 0 is the first). The file's FCB/inode is already in memory, and **nothing else is cached**. Under indexed allocation, the file has a single index block, which is on disk. How many disk reads does reading block 150 take under **contiguous, linked and indexed** allocation, respectively?

- A. 1, 150, 2
- B. 1, 151, 2
- C. 150, 1, 2
- D. 1, 151, 1

**Answer:** B

**Why:** Count the reads needed to *find* the block, plus the read of the block itself.

**Working:**
- **Contiguous:** disk block = start + 150, computed from the FCB. Read it → **1** read.
- **Linked:** block 150's location is stored in block 149, block 149's in block 148, and so on. You must read blocks 0 … 149 (**150 reads**) to follow the chain, then read block 150 → **151** reads.
- **Indexed:** read the index block (1), take entry 150, then read the data block (1) → **2** reads.

**Why the other options fail:**
- A: Forgets the final read of block 150 itself under linked allocation.
- C: Swaps contiguous and linked.
- D: Assumes the index block is already cached, which the question rules out.

**Concept recap:**
- Contiguous: 1 read.
- Linked: i + 1 reads (or 1 read if the FAT is cached).
- Indexed: (number of index levels) + 1 reads.

**Trap alert:** Read the caching assumptions carefully. If the index block *were* cached, indexed allocation would need 1 read. Exams change these assumptions to change the answer.

**Interview angle:** Follow-up: *"How many reads with a two-level index, nothing cached?"* 3: the outer index, the inner index, then the data.

---

## Q220. [M8 | Implementing large directories | Medium | MCQ | Interview]

A directory holding **100,000** files is implemented as a **linear list** of entries. What is the problem, and what's a common fix?

- A. There is no problem; linear lists scale well
- B. Limit every directory to 100 files
- C. Store each directory contiguously on disk
- D. Each lookup or create must scan the list, which is O(n), and creating a file must also check for a duplicate name. Index the directory with a hash table or a B-tree (e.g., ext4's hashed B-tree, "htree")

**Answer:** D

**Why:**
- A linear list is simple, but every `open` of a name, and every `create`, scans up to all 100,000 entries.
- A **hash table** gives about O(1) lookups but must handle collisions and resizing.
- A **B-tree** keeps names sorted, scales well on disk, and allows ordered listings.

ext4 uses a hashed B-tree (`dir_index`), and XFS and btrfs use B-trees.

**Why the other options fail:**
- A: At this size, the O(n) scans dominate.
- B: An arbitrary restriction that doesn't fix the data structure.
- C: Contiguity speeds up reading the list, but searching it is still linear.

**Concept recap:** Directory = a map from name to inode number. Choose the data structure the way you would for any map: list, hash table, or tree.

**Trap alert:** The cost isn't only lookups. Create and delete must also find their spot and check for duplicates.

**Interview angle:** "Why is `ls` slow in a directory with millions of files?" Reading and sorting all the entries, plus a `stat` per file if `ls` shows details. Follow-up: *"How does the OS avoid repeated directory lookups?"* With the **dentry cache**, which caches name → inode lookups in memory.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 201 | A | 206 | A, B, C, E | 211 | C | 216 | A |
| 202 | C | 207 | C | 212 | A | 217 | A, B, C, D |
| 203 | D | 208 | A, B, D, E | 213 | A, B, C, E | 218 | A, B, C, D |
| 204 | B | 209 | D | 214 | D | 219 | B |
| 205 | A | 210 | B | 215 | C | 220 | D |

---

**Concepts tested so far (to avoid repeats):** Batches 1–10, plus: working-set computation (distinct pages in the window) · D = ΣWSS > m → suspend a process · thrashing and its feedback loop · PFF bounds and suspending a process · proportional allocation · global vs local replacement, ISA minimum frames · COW fault on write, `vfork` · `mmap` shared vs private, lazy loading, anonymous memory → swap, page-level swapping · p_opt = √(2se) and the TLB/I/O trade-offs · the name isn't in the inode · why `open()` exists, permissions checked at open · direct access by record number · directory structures, no hard links to directories · relative paths with `..` · contiguous block and offset arithmetic · FAT vs plain linked allocation · indexed allocation and its size limit · the allocation comparison table · disk reads to reach block i under each scheme · directories with hash tables or B-trees, the dentry cache
*(M7 is complete. M8 continues in Batch 12 with inode pointers and maximum file size, open-file tables and fds, hard vs soft links, free-space management, permissions, mounting, journaling and consistency, and caching.)*

**Completed: 220 / 300    Remaining: 80**
