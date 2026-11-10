# OS MCQ/MSQ Bank — Batch 12 (Q221–Q240)

**What this batch covers:**
- **M8, file systems completed (Q221–237):**
  - Inode pointers: maximum file size, and which pointer covers a given offset.
  - Open-file tables, independent offsets, fd leaks.
  - Hard vs symbolic links, and directory link counts.
  - Free-space management (bitmap size, lists, grouping, counting).
  - Unix permissions, including on directories.
  - Mounting, journaling and its modes, crash-consistency ordering.
  - The page cache and `fsync`, and the superblock.
- **M9, I/O begins (Q238–240):** memory-mapped I/O, polling vs interrupts, DMA.

> Conventions: KB = 2¹⁰ bytes, MB = 2²⁰, GB = 2³⁰, TB = 2⁴⁰. Platform-specific behaviour is Linux/ext4 unless stated.

---

## Q221. [M8 | Maximum file size with inode pointers | Hard | Numerical | OA + Interview]

A Unix-style inode has **12 direct pointers**, **1 single-indirect**, **1 double-indirect** and **1 triple-indirect** pointer. The block size is **4 KB** and each block pointer is **4 bytes**. Considering only the pointer structure, what is the maximum file size?

- A. ≈ 4 GB
- B. ≈ 4 MB
- C. ≈ 4 TB (exactly 4 TB + 4 GB + 4 MB + 48 KB)
- D. ≈ 16 TB

**Answer:** C

**Why:** Each kind of pointer reaches a different number of data blocks. Add up all the blocks the inode can reach, then multiply by the block size.

**Working:**
- **Pointers per block:** 4096 / 4 = **1024 = 2¹⁰**.
- **Direct:** 12 blocks → 12 × 4 KB = **48 KB**.
- **Single indirect:** 1024 blocks → 2¹⁰ × 2¹² B = **4 MB**.
- **Double indirect:** 1024² blocks → 2²⁰ × 2¹² B = **4 GB**.
- **Triple indirect:** 1024³ blocks → 2³⁰ × 2¹² B = **4 TB**.
- **Total** = 4 TB + 4 GB + 4 MB + 48 KB ≈ **4 TB**. The triple-indirect term dominates.

**Why the other options fail:**
- A: Stops at the double-indirect level.
- B: Counts only the direct and single-indirect pointers.
- D: Uses 2048 pointers per block, i.e. 2-byte pointers.

**Concept recap:**
- k = B / ptr (pointers per block).
- Max file size = (D + k + k² + k³) × B, where D is the number of direct pointers.

**Trap alert:** The question asks only about the **pointer structure**. Real file systems can hit other limits first. For example, ext3 with 4 KB blocks caps files at 2 TB, because its block-count field counts 512-byte sectors in 32 bits. Also, don't forget the direct blocks when computing an exact answer.

**Interview angle:** "Why do inodes mix direct and indirect pointers?" Small files, which are most files, need no indirect reads, while huge files are still possible. Follow-up: *"What do ext4 and XFS use instead?"* Extent trees. One extent maps a long contiguous run of blocks.

---

## Q222. [M8 | Which pointer covers a given offset | Medium | MCQ | OA]

Same inode as Q221 (12 direct, 4 KB blocks, 1024 pointers per block). The inode is in memory, and nothing else is cached. To read the byte at file offset **10 MB**, which kind of pointer is used, and how many disk reads are needed?

- A. Direct; 1 read
- B. Single indirect; 2 reads
- C. Double indirect; 3 reads
- D. Triple indirect; 4 reads

**Answer:** C

**Why:** Find the logical block number, then see which pointer range it falls in.
- Logical block = 10 × 2²⁰ / 2¹² = **2560**.
- Ranges: direct = blocks 0–11; single indirect = 12–1035; double indirect = 1036 to 1036 + 1024² − 1.
- 2560 > 1035, so the block is reached through the **double-indirect** pointer.
- Reads: the double-indirect block (1) → the single-indirect block it points to (1) → the data block (1) = **3 reads**.

**Why the other options fail:**
- A: 2560 is far beyond the 12 direct blocks.
- B: The single-indirect range ends at block 1035.
- D: The triple-indirect range starts beyond about 4 GB.

**Concept recap:** Reads needed = (number of indirection levels) + 1 for the data, when nothing is cached.

**Trap alert:** When finding the range, remember that the **12 direct blocks come first**. The single-indirect range starts at block 12, not 0.

**Interview angle:** Follow-up: *"Why is reading a large file sequentially still efficient?"* Each indirect block is read once and cached, then reused for its next 1024 data blocks.

---

## Q223. [M8 | Independent opens of the same file | Medium | MCQ | OA + Interview]

Processes A and B each call `open()` on the same file **independently** (no `fork`, no `dup`). A reads 100 bytes. Where does B's next `read()` start?

- A. At offset 0. Each `open()` creates its own entry in the system-wide open-file table, with its own offset; both entries point to the same in-memory inode
- B. At offset 100
- C. At the end of the file
- D. The result is undefined

**Answer:** A

**Why:** Unix keeps three levels of tables:
1. **The per-process fd table:** fd number → an entry in the open-file table.
2. **The system-wide open-file table:** **the current offset**, the access mode, and a reference count, with a pointer to the in-memory inode.
3. **In-memory inodes** (v-nodes): one per open file, shared by everyone who has that file open.

Each separate `open()` creates a separate entry at level 2, so each process has its own offset.

**Why the other options fail:**
- B: Offsets are shared only when two fds point to the **same** level-2 entry, as happens after `fork` or `dup`.
- C, D: Nothing in a fresh open sets the offset to the end of the file. (`O_APPEND` would make *writes* go to the end.)

**Concept recap:** Separate opens → separate offsets, but the same inode. Shared entries after `fork`/`dup` are covered in M11.

**Trap alert:** Writes from A *are* visible to B's later reads, because both go through the same inode and page cache. Only the **position** is independent.

**Interview angle:** "Describe the Unix file tables." Follow-up: *"Where is the file offset stored, and why does that matter?"* In the open-file-table entry. That's why it's shared after `fork` and `dup`.

---

## Q224. [M8 | File-descriptor leaks | Easy | MCQ | OA]

A long-running server accepts connections but never closes the sockets. After a while, `accept()` fails with `EMFILE`. Why?

- A. The disk is full
- B. The file system is mounted read-only
- C. There are too many processes in the system
- D. The process has reached its per-process limit on open file descriptors (e.g., `ulimit -n`), because it never closes its sockets: an fd leak

**Answer:** D

**Why:** Every socket uses an fd, and a process may hold only a limited number of fds at once. If the server never calls `close()`, it uses up the fd table. `EMFILE` means "this process has too many open files". (`ENFILE` would mean the *system-wide* table is full.)

**Why the other options fail:**
- A: A full disk gives `ENOSPC`.
- B: A read-only file system gives `EROFS`.
- C: Too many processes makes `fork` fail with `EAGAIN`.

**Concept recap:** fds are a limited resource, like memory. Close them on every path. `ls /proc/<pid>/fd | wc -l` counts the ones open.

**Trap alert:** Raising the limit only postpones the failure. Fix the leak itself: close on every path, use RAII, and use try-with-resources in Java.

**Interview angle:** "Your service starts failing with 'Too many open files'. What do you check?" `lsof -p <pid>`, the listing of `/proc/<pid>/fd`, and the `ulimit`. Follow-up: *"Does a fork-heavy server leak fds into its children?"* Yes, unless the fds are marked close-on-exec (covered in M11).

---

## Q225. [M8 | Hard vs symbolic links in action | Medium | MCQ | OA + Interview]

```
$ echo hi > a
$ ln a b          # hard link
$ ln -s a c       # symbolic link
$ rm a
$ cat b ; cat c
```

What happens?

- A. Both `cat` commands print `hi`
- B. `cat b` prints `hi`; `cat c` fails, because `c` is now a dangling symlink ("No such file or directory")
- C. `cat b` fails; `cat c` prints `hi`
- D. Both fail

**Answer:** B

**Why:**
- A **hard link** is just another directory entry for the **same inode**. `rm a` removes one name and decrements the link count from 2 to 1. The data stays, and `b` still reaches it.
- A **symbolic link** is a separate small file that stores the **path** "a". With `a` gone, following that path fails.

**Why the other options fail:**
- A: The symlink stores a name, not an inode, so it can't find the data once the name is gone.
- C: Reversed.
- D: The hard link keeps the inode alive.

**Concept recap:**
- Hard link = name → inode (link count + 1).
- Symlink = name → path string (link count unchanged).
- Data is freed when the link count reaches 0 **and** no process has the file open.

**Trap alert:** If you later create a *new* file named `a`, `c` will work again, but it now shows the new file's contents.

**Interview angle:** "Hard link vs soft link?" Follow-up: *"Which survives moving the target?"* A hard link always works, because it doesn't use the path. A relative symlink breaks if the target moves.

---

## Q226. [M8 | Link properties | Medium | MSQ | Interview]

**Select all that apply.**

- A. A hard link has the same inode number as the original file.
- B. A hard link cannot cross file systems.
- C. A symbolic link can cross file systems and can point to a directory.
- D. Deleting a hard-linked name decrements the link count; the data is freed only when the count reaches 0 and no process still has the file open.
- E. Creating a symbolic link increments the target's link count.

**Answer:** A, B, C, D

**Why:**
- **A, B:** A hard link stores an **inode number**, and inode numbers are meaningful only within one file system.
- **C:** A symlink stores a **path**, which can point anywhere, including other mounts, directories, and paths that don't exist.
- **D:** This is reference counting on the inode, plus open handles (Q211).

**Why the other options fail:**
- E: Symlinks don't affect the target's link count at all.

**Concept recap:**

| | Hard link | Symbolic link |
|---|---|---|
| Points to | an inode | a path |
| Crosses file systems | no | yes |
| Can point to a directory | no (Q213) | yes |
| Survives deleting the original name | yes | no (dangling) |
| Extra inode | no | yes |

**Trap alert:** `ls -i` shows the same number for hard links, and a different one for the symlink itself.

**Interview angle:** Follow-up: *"How do package managers use links?"* Symlinks for version switching, such as `libfoo.so → libfoo.so.1.2`.

---

## Q227. [M8 | Directory link count | Medium | MCQ | OA]

Directory `d` contains **3 subdirectories** and **5 regular files**. With classic Unix semantics (ext2/3/4), what is `d`'s link count?

- A. 8
- B. 2
- C. 10
- D. 5

**Answer:** D

**Why:** A directory's link count counts the directory entries that name **this directory's inode**:
- Its entry `d` in the parent directory → 1
- Its own `.` entry → 1
- The `..` entry in each of its 3 subdirectories → 3

Total = 1 + 1 + 3 = **5**. Regular files inside `d` don't point back to `d`, so they don't count.

**Why the other options fail:**
- A: Counts the entries `d` contains (3 + 5) instead of the links pointing to it.
- B: The count for an empty directory.
- C: Counts the files as well, plus extra.

**Concept recap:** A directory's link count = 2 + (number of subdirectories). A new empty directory has a link count of 2.

**Trap alert:** This holds for classic Unix file systems. Some file systems (for example btrfs) don't follow it exactly.

**Interview angle:** Rarely asked verbally. Follow-up: *"How could you count the subdirectories without listing them?"* Subtract 2 from the directory's link count, on file systems that follow this rule.

---

## Q228. [M8 | Free-space bitmap size | Medium | Numerical | OA]

A **1 TB** disk uses **4 KB** blocks. Free space is tracked with a **bitmap** of one bit per block. How large is the bitmap?

- A. 256 MB
- B. 32 MB
- C. 4 MB
- D. 32 KB

**Answer:** B

**Why:** The bitmap needs one bit for every block on the disk.

**Working:**
1. Number of blocks = 2⁴⁰ / 2¹² = **2²⁸**.
2. Bits = 2²⁸ → bytes = 2²⁸ / 8 = **2²⁵ bytes = 32 MB**.

**Why the other options fail:**
- A: Uses one **byte** per block (2²⁸ bytes).
- C: Off by a factor of 8 in the other direction.
- D: A units slip (KB for MB).

**Concept recap:** Bitmap size = (disk size / block size) / 8 bytes. It's a small fraction of the disk (here 1/32,768), and it can be cached.

**Trap alert:** Divide by **8** to turn bits into bytes, exactly once.

**Interview angle:** "How does a file system track free blocks?" Follow-up: *"Why do bitmaps make it easy to find contiguous free space?"* Search for a run of 0 bits. CPUs even have instructions to find the first set or clear bit.

---

## Q229. [M8 | Free-space management methods | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. A bitmap makes it easy to find runs of contiguous free blocks.
- B. A linked free list (each free block points to the next free block) needs no extra space, but traversing it is slow.
- C. **Grouping** stores the addresses of n free blocks inside the first free block; the last of those addresses points to the next such block.
- D. **Counting** stores (start, length) pairs, and works well when free space comes in contiguous runs.
- E. A bitmap must be kept only on disk and never cached in memory.

**Answer:** A, B, C, D

**Why:**
- **A:** A bitmap gives a compact, searchable view of where free runs are.
- **B:** The free list's links live inside the free blocks themselves, so it costs no extra space. But finding many free blocks means reading each one.
- **C:** Grouping returns n free addresses with a single block read.
- **D:** Counting (extent-style) uses one record per free run.

**Why the other options fail:**
- E: Bitmaps are cached for speed and written back carefully. The journal protects them (Q233).

**Concept recap:** The four methods are bitmap, linked list, grouping and counting. Modern file systems use bitmaps (ext4) or extent/B-tree free-space structures (XFS, btrfs).

**Trap alert:** If the cached bitmap and the disk copy disagree after a crash, blocks may be lost or double-allocated. That's why journaling covers bitmap updates.

**Interview angle:** Follow-up: *"Which free-space method suits an SSD with TRIM?"* Extent-based tracking, so large freed ranges can be reported to the drive efficiently.

---

## Q230. [M8 | Unix permission bits | Easy | MCQ | OA]

After `chmod 754 file`, what does `ls -l` show for the permissions?

- A. `rwxr-xr--`
- B. `rwxr--r-x`
- C. `rw-r-xr--`
- D. `r-xrwxr--`

**Answer:** A

**Why:** Each octal digit encodes r = 4, w = 2, x = 1, for owner, group and others in that order.
- 7 = 4 + 2 + 1 = `rwx` (owner)
- 5 = 4 + 1 = `r-x` (group)
- 4 = `r--` (others)

**Why the other options fail:**
- B: Swaps the group and others digits.
- C: Gives the owner 6 (`rw-`).
- D: Scrambles the digits.

**Concept recap:** Common modes: 755 (programs, directories), 644 (ordinary files), 600 (private files), 700 (private directories).

**Trap alert:** Permissions are checked in order **owner → group → others**, using the **first** class that matches you. An owner with `---` is denied even if "others" have `rwx`.

**Interview angle:** Follow-up: *"What does the setuid bit do?"* An executable with setuid runs with its **owner's** privileges. `passwd` uses it to update a protected file, and it's a common privilege-escalation target (covered in M10).

---

## Q231. [M8 | Permissions on directories | Hard | MCQ | Interview]

Directory `/shared` has mode **711** (`rwx--x--x`) and is owned by alice. Inside it, `report.txt` has mode **644**. Bob is not in alice's group. He runs `ls /shared` and then `cat /shared/report.txt`. What happens?

- A. `ls` works; `cat` fails
- B. Both fail
- C. Both work
- D. `ls` fails (no read permission on the directory), but `cat` works (execute permission on the directory allows traversal by name, and the file itself is world-readable)

**Answer:** D

**Why:** On a **directory**, the permission bits mean:
- **r:** list the names inside (`ls`).
- **w:** create, delete or rename entries (this also needs x).
- **x:** **traverse** it, i.e. use it as part of a path to reach something inside.

Bob has only x (`--x`), so he can't list the directory, but he *can* open a file whose name he already knows. The file's own `r--` for others then lets him read it.

**Why the other options fail:**
- A: Reversed. Listing needs r, which Bob lacks.
- B: x on the directory is enough to reach a known name.
- C: Listing still needs r.

**Concept recap:** Directory r = list names; w = change entries; x = pass through. To open `/a/b/c`, you need x on `/`, `/a` and `/a/b`, plus the right permission on `c` itself.

**Trap alert:** Deleting a file depends on **w (and x) on the directory**, not on the file's own permissions. That's why `/tmp` sets the **sticky bit**, so users can delete only their own files there.

**Interview angle:** "Can you delete a read-only file?" Yes, if you have w and x on its directory. Follow-up: *"What does the sticky bit on `/tmp` do?"* It restricts deleting and renaming to the file's owner, the directory's owner, or root.

---

## Q232. [M8 | Mounting | Easy | MCQ | OA + Interview]

`/mnt/usb` contains a file `old.txt`. You then mount a USB file system on `/mnt/usb`. What happens to `old.txt`?

- A. It is deleted
- B. It is moved onto the USB drive
- C. It is hidden (unreachable through that path) until the USB file system is unmounted; it still exists on the original disk
- D. It is merged into the USB drive's listing

**Answer:** C

**Why:** **Mounting** attaches another file system's root at a directory, the **mount point**. Path lookups that reach `/mnt/usb` now cross into the USB file system, so whatever was underneath is **covered**, not removed. Unmounting reveals it again.

**Why the other options fail:**
- A, B: Mounting never touches the underlying directory's contents.
- D: Plain mounts don't merge directories. Union or overlay file systems (OverlayFS, used by containers) are special cases that do.

**Concept recap:** The mount table maps mount points to file systems. Path resolution checks for a mount at each step.

**Trap alert:** A process that already had `old.txt` open (or had its cwd under the old directory) can still use it after the mount.

**Interview angle:** "What does mounting do?" Follow-up: *"How do Docker images get layered file systems?"* OverlayFS: read-only image layers plus a writable upper layer, merged into one view.

---

## Q233. [M8 | Journaling | Medium | MCQ | OA + Interview]

What does a **journaling** file system do?

- A. Logs every user keystroke for auditing
- B. Before changing metadata in place, writes the change as a transaction to an on-disk log. After a crash, it replays committed transactions and discards incomplete ones, so recovery is fast and needs no full-disk scan
- C. Keeps two complete copies of the disk
- D. Writes data to disk only once a day

**Answer:** B

**Why:** A single file operation, such as an append, updates several structures: the inode, the bitmap, a directory. A crash in the middle leaves them inconsistent (Q235). **Write-ahead logging** first writes "here is what I'm about to change" to the journal, then a commit record, and only then the real locations. Recovery replays committed transactions and ignores uncommitted ones. Each operation therefore becomes **atomic**: it either happened fully or not at all.

**Why the other options fail:**
- A: This journal is about the file system's own updates, not user activity.
- C: That's mirroring (RAID 1).
- D: That's an absurd caching policy.

**Concept recap:** Journaling = write-ahead log + commit record + checkpoint (writing the changes to their real homes). Recovery time depends on the journal's size, not the disk's.

**Trap alert:** Journaling protects the **file system's consistency**. It doesn't promise that your most recent writes survive (Q234 D).

**Interview angle:** "Why journal?" Follow-up: *"How does this relate to databases?"* It's the same idea as a database's write-ahead log (WAL).

---

## Q234. [M8 | Journaling modes and alternatives | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. ext4's default **ordered** mode journals only metadata, but writes file data to disk before committing the related metadata.
- B. **data=journal** mode logs both data and metadata. It's the safest mode, but writes file data twice.
- C. Without journaling, crash recovery needs a full `fsck` scan, whose time grows with the size of the disk.
- D. Journaling guarantees that no data written in the last few seconds is ever lost in a crash.
- E. Copy-on-write file systems (btrfs, ZFS) keep themselves consistent by never overwriting live blocks in place.

**Answer:** A, B, C, E

**Why:**
- **A:** Ordered mode stops metadata from pointing at blocks that hold garbage, which is the Q235 "only inode written" case.
- **B:** Full data journaling costs write bandwidth, since everything is written twice.
- **C:** `fsck` must cross-check every inode, bitmap and directory.
- **E:** COW file systems write new versions elsewhere, then atomically switch a root pointer, so the on-disk tree is always consistent.

**Why the other options fail:**
- D: Data still in the page cache, not yet committed, is lost in a crash (Q236). Journaling recovers **consistency**, not every write.

**Concept recap:** ext4 has three modes: writeback (metadata only, no ordering) < ordered (the default) < journal (data and metadata). COW is an alternative design to journaling.

**Trap alert:** The defaults differ between file systems and versions. Name the platform when you discuss modes.

**Interview angle:** "Journaling vs copy-on-write file systems?" Follow-up: *"What else does COW make cheap?"* Snapshots: just keep the old root pointer.

---

## Q235. [M8 | Crash consistency | Hard | MCQ | OA]

Appending one block to a file needs three writes: the **data block D**, the **inode I** (new size and block pointer), and the **bitmap B** (marking the block as used). There's no journaling. A crash happens after exactly **one** of the three writes has reached the disk. Which single completed write leaves the file system **consistent**, losing only the new data?

- A. Only D (the data block) was written
- B. Only I (the inode) was written
- C. Only B (the bitmap) was written
- D. Any one of them, since the file system is always consistent

**Answer:** A

**Why:** Check what each structure claims afterwards:
- **Only D:** no inode points to the block, and the bitmap still says it's free. It's as if the append never happened → **consistent**, and the new data is simply lost.
- **Only I:** the inode points to a block that the bitmap says is free, and that block holds **garbage** → inconsistent. The block could be handed out twice, and a reader would see old data.
- **Only B:** the block is marked used, but nothing points to it → inconsistent. That's a space leak.

**Why the other options fail:**
- B, C: Each leaves metadata that contradicts itself.
- D: Without journaling, crashes do leave file systems inconsistent.

**Concept recap:** The safe write order is **data first, then metadata**. Ordered journaling enforces exactly this (Q234 A). `fsck` repairs the other cases after the fact.

**Trap alert:** "The inode is most important, so write it first" is the dangerous order. A pointer that's written before its target can point at garbage.

**Interview angle:** "In what order should a file system write data and metadata?" Data before the metadata that points to it. Follow-up: *"How does `fsck` detect the 'only B' case?"* It rebuilds the bitmap from what the inodes actually reference, and finds the block marked used with no owner.

---

## Q236. [M8 | Page cache and durability | Medium | MSQ | Interview]

Linux. **Select all that apply.**

- A. `write()` usually returns once the data has been copied into the page cache, before it reaches the disk.
- B. `fsync()` forces a file's dirty data and metadata to stable storage before returning.
- C. A unified page cache serves both `read()`/`write()` and `mmap`, which avoids caching the same data twice and keeps the two views consistent.
- D. `O_DIRECT` bypasses the page cache, for applications (such as databases) that manage their own caching.
- E. Data in the page cache is lost in a crash only if it is clean.

**Answer:** A, B, C, D

**Why:**
- **A:** This is **write-back caching**: fast writes, delayed persistence. Dirty pages are flushed periodically, or under memory pressure.
- **B:** `fsync` is the application's way to say "this must survive a crash".
- **C:** Older Unix systems had separate buffer and page caches, which risked inconsistency and double copies. Merging them fixed that.
- **D:** Databases often use `O_DIRECT` to avoid double caching with their own buffer pool.

**Why the other options fail:**
- E: Reversed. **Dirty** cached data that hasn't been written back is lost. Clean data is also safely on disk.

**Concept recap:** write → page cache (dirty) → writeback → disk. Durability requires `fsync`, `fdatasync`, or `O_SYNC`.

**Trap alert:** "My program called `write()` and the machine lost power, so the data must be on disk" is false. Without `fsync`, it may not be. Also call `fsync` on the **directory** after creating or renaming a file, so the new entry itself is durable.

**Interview angle:** "How do you make a write durable?" Follow-up: *"How do editors save files atomically?"* Write a temporary file, `fsync` it, `rename` it over the original, then `fsync` the directory.

---

## Q237. [M8 | The superblock | Easy | MCQ | OA]

Which on-disk structure holds **file-system-wide** information, such as the block size, the total and free block counts, and the inode counts?

- A. An inode
- B. A directory entry
- C. The superblock
- D. A FAT entry

**Answer:** C

**Why:** The **superblock** describes the whole file system: its geometry, its sizes, where other structures start, and its state ("cleanly unmounted?"). It's read at mount time. Because losing it is catastrophic, ext file systems keep **backup copies** spread across the disk.

**Why the other options fail:**
- A: An inode describes one file.
- B: A directory entry maps one name to one inode.
- D: A FAT entry describes the next block for one block.

**Concept recap:** Main on-disk structures: boot block, superblock, inode tables, bitmaps, data blocks, and the journal (if the file system journals).

**Trap alert:** The superblock is **file-system-wide**. Don't confuse it with the per-file inode.

**Interview angle:** Follow-up: *"What happens if the primary superblock is corrupted?"* Mount or repair using a backup copy (for example `e2fsck -b <backup>`).

---

## Q238. [M9 | Memory-mapped I/O | Easy | MCQ | Interview]

What is **memory-mapped I/O**?

- A. Device registers can be accessed only through special IN/OUT instructions
- B. Device registers appear at physical addresses, so the CPU accesses them with ordinary load and store instructions
- C. Devices copy their contents into RAM at boot
- D. A device's memory is swapped to disk

**Answer:** B

**Why:** Each device **controller** exposes registers: status, control, data-in, data-out. With memory-mapped I/O, those registers sit at physical addresses, and a driver reads or writes them like memory. The kernel maps those addresses **uncached**, so every access really reaches the device.

**Why the other options fail:**
- A: That's **port-mapped I/O**, which x86 also has (the `in`/`out` instructions) alongside MMIO.
- C, D: Unrelated.

**Concept recap:** CPU ↔ (MMIO or port I/O) ↔ controller registers ↔ device. The driver is the OS code that knows a particular controller's registers.

**Trap alert:** MMIO addresses aren't RAM. Caching them, or letting the compiler reorder the accesses, breaks drivers. Drivers use `volatile` and memory barriers.

**Interview angle:** "How does software talk to hardware?" Through controller registers, accessed via MMIO or ports. Follow-up: *"Why can't user programs write to those addresses?"* They aren't mapped into user address spaces, and access to them is privileged.

---

## Q239. [M9 | Polling vs interrupts | Medium | MCQ | OA]

A device completes each operation within a few hundred nanoseconds, or a NIC (network card) is receiving an extremely high packet rate. Which completion mechanism is usually more efficient?

- A. One interrupt per event, always
- B. DMA with no completion notification at all
- C. A sleep-and-wake loop on a timer
- D. Polling: checking the status register costs a few cycles, whereas each interrupt costs a context save and restore plus handler dispatch. At very high event rates, interrupts can even cause livelock, so Linux NAPI switches to polling under load

**Answer:** D

**Why:**
- **Interrupts** win when events are **rare or slow**: the CPU does other work and pays the interrupt cost only occasionally.
- **Polling** wins when events are **imminent or very frequent**: a quick check is cheaper than an interrupt, and batching avoids a storm of interrupts.

**Why the other options fail:**
- A: "Always" ignores the overhead and livelock risk at high rates (Q8 follow-up).
- B: The OS still needs to learn when transfers finish.
- C: Sleeping is far too slow for sub-microsecond completions.

**Concept recap:** Choose between polling and interrupts by comparing the expected wait with the interrupt cost. Hybrids: interrupt coalescing, NAPI, and polling-mode NVMe and DPDK.

**Trap alert:** "Polling is always wasteful" is wrong. It's the right choice for very fast devices and very high event rates.

**Interview angle:** "When would you poll instead of using interrupts?" Follow-up: *"What is interrupt coalescing?"* The device raises one interrupt for several completions, trading a little latency for less overhead.

---

## Q240. [M9 | DMA | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. A DMA controller transfers blocks of data between a device and memory without the CPU copying each byte.
- B. The CPU is typically interrupted once per completed block transfer, not once per byte.
- C. A DMA transfer can briefly compete with the CPU for the memory bus ("cycle stealing").
- D. The memory being transferred must stay at fixed physical locations (pinned) while the DMA is in progress.
- E. DMA removes the need for interrupts.

**Answer:** A, B, C, D

**Why:**
- **A, B:** The CPU programs the DMA engine with a source, a destination and a count. The engine moves the data, then raises **one** interrupt when it's done. The CPU is free in the meantime.
- **C:** The memory bus is shared, so DMA traffic occasionally delays CPU memory accesses.
- **D:** The device works with **physical** addresses. If the OS paged out or moved that memory mid-transfer, the data would go to the wrong place (compare Q147).

**Why the other options fail:**
- E: DMA still uses an interrupt (or polling) to report completion.

**Concept recap:**
- Programmed I/O: the CPU moves every byte.
- Interrupt-driven I/O: the CPU moves data on each interrupt.
- DMA: the device engine moves the data, and the CPU is interrupted once at the end.

**Trap alert:** DMA bypasses the CPU caches. The OS must keep the caches coherent with DMA buffers, by flushing or invalidating them on non-coherent systems.

**Interview angle:** "Explain DMA and why it matters." Follow-up: *"What is an IOMMU for?"* It translates device addresses and restricts which memory each device may access. That gives isolation, and it's how devices are safely passed through to VMs.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 221 | C | 226 | A, B, C, D | 231 | D | 236 | A, B, C, D |
| 222 | C | 227 | D | 232 | C | 237 | C |
| 223 | A | 228 | B | 233 | B | 238 | B |
| 224 | D | 229 | A, B, C, D | 234 | A, B, C, E | 239 | D |
| 225 | B | 230 | A | 235 | A | 240 | A, B, C, D |

---

**Concepts tested so far (to avoid repeats):** Batches 1–11, plus: inode maximum file size (D + k + k² + k³)·B, and the ext3 2 TB note · offset → pointer range → reads · independent opens have independent offsets, the three-level tables · EMFILE fd leaks · hard vs symlink after `rm` · link properties table · directory link count = 2 + subdirectories · bitmap size · bitmap, list, grouping, counting · octal permissions · directory r/w/x and the sticky bit · mount hides the underlying contents, OverlayFS · journaling (WAL, commit, replay) · ordered/journal/writeback modes, COW file systems · crash-consistency write order · page cache, `fsync`, `O_DIRECT`, atomic save by rename · superblock and its backups · memory-mapped vs port-mapped I/O · polling vs interrupts, NAPI, interrupt coalescing · DMA, pinning, IOMMU
*(M8 is complete. M9 continues in Batch 13 with buffering, spooling, drivers, disk access time, the six disk-scheduling algorithms, RAID, and SSD vs HDD, followed by the start of M10.)*

**Completed: 240 / 300    Remaining: 60**
