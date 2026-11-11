# OS MCQ/MSQ Bank — Batch 15 (Q281–Q300)

**What this batch covers:**
- **M11, Linux/Unix internals completed (Q281–288):** terminal signals and process groups, signal state across `fork`/`exec`, the CFS/EEVDF scheduler, `nice`, the vDSO (system calls that skip the kernel), `/proc`, overcommit and the OOM killer, containers vs VMs for running another OS.
- **M12, integrated questions (Q289–300):** priority inversion (scheduling × locks), EMAT with both TLB misses and page faults, `fork` × shared vs private memory, deadlock × memory allocation, inherited sockets × `exec`, dirty-page writeback × latency, file layout × disk access time, races × multicore, memory isolation × security, a system call's full path through scheduler and interrupt, a first-touch page fault on an `mmap`'d file, and a whole-system pipeline capstone.

**Progress report #3 follows the answer key at the end of this file.** The final deliverables and the three mock OAs are in separate files.

---

## Q281. [M11 | Ctrl+C and process groups | Easy | MCQ | Interview]

A terminal is running the pipeline `a | b` in the foreground. You press **Ctrl+C**. What happens?

- A. Only `b` receives SIGINT
- B. Only the shell receives SIGINT
- C. The terminal driver sends SIGINT to every process in the foreground process group, so both `a` and `b` receive it (and by default terminate); the shell is in a different process group, so it survives
- D. `a` receives SIGKILL

**Answer:** C

**Why:** The shell puts each pipeline in its own **process group** and makes that group the terminal's **foreground** group. Special keys are turned into signals by the terminal driver and delivered to the whole foreground group:
- Ctrl+C → **SIGINT** (interrupt)
- Ctrl+Z → **SIGTSTP** (stop, so the job can be resumed with `fg`/`bg`)
- Ctrl+\ → **SIGQUIT** (quit, with a core dump)

**Why the other options fail:**
- A: The signal goes to the whole group, not just the last process.
- B: The shell is deliberately *not* in the foreground group, so Ctrl+C doesn't kill your shell.
- D: Ctrl+C sends SIGINT, which a program can catch. Nothing sends SIGKILL here.

**Concept recap:** Process group = a pipeline (job). Session = everything started from one terminal. The terminal sends key-generated signals to the foreground process group.

**Trap alert:** A program that catches SIGINT (to clean up, or to ignore it) won't die on Ctrl+C. That's a choice made by the program, not a failure of the signal.

**Interview angle:** "What happens when you press Ctrl+C?" Terminal driver → SIGINT → foreground process group. Follow-up: *"How does `kill -INT -1234` differ from `kill -INT 1234`?"* A negative PID addresses the whole process group 1234.

---

## Q282. [M11 | Signal state across fork and exec | Medium | MSQ | Interview]

POSIX/Linux. **Select all that apply.**

- A. `kill(pid, sig)` sends a signal to one process, and `kill(-pgid, sig)` sends it to a whole process group.
- B. The signal mask (the set of blocked signals) is inherited across `fork` and preserved across `exec`.
- C. A child's set of pending signals starts empty after `fork`.
- D. Signals with custom handlers are reset to their default action on `exec`; signals set to "ignore" stay ignored.
- E. A process can block SIGKILL with `sigprocmask`.

**Answer:** A, B, C, D

**Why:**
- **A:** Negative PIDs address process groups (this is what job control uses).
- **B:** The mask describes which signals the thread is choosing to defer. Nothing about `fork` or `exec` invalidates that choice, so it carries over.
- **C:** Covered in Q29: a signal sent to the parent shouldn't be delivered twice.
- **D:** After `exec`, the handler code no longer exists, so "run this handler" can't survive. "Ignore" doesn't refer to any code, so it can (Q27).

**Why the other options fail:**
- E: SIGKILL and SIGSTOP can't be caught, ignored **or blocked** (Q45). Attempts to block them are silently ignored.

**Concept recap:**

| Across `fork` | Across `exec` |
|---|---|
| handlers inherited | handlers → default |
| ignored stays ignored | ignored stays ignored |
| mask inherited | mask preserved |
| pending cleared | pending preserved |

**Trap alert:** A process that starts with SIGINT ignored (for example via `nohup`, or because its parent blocked or ignored it) may not respond to Ctrl+C, because that state survived `exec`.

**Interview angle:** "What does a child inherit about signals?" Use the table. Follow-up: *"Why does `nohup` work?"* It sets SIGHUP to "ignore" before `exec`, and that setting survives into the program.

---

## Q283. [M11 | Linux's fair scheduler | Medium | MCQ | Interview]

Which description best matches how Linux's **CFS** (Completely Fair Scheduler) chooses the next task?

- A. Strict round robin with a fixed 100 ms quantum for every task
- B. Static priority queues, with no notion of fairness
- C. Shortest job first, using predicted CPU bursts
- D. Each task accumulates a **virtual runtime** that grows more slowly for higher-weight (lower-nice) tasks, and the scheduler runs the runnable task with the smallest virtual runtime, kept in a red-black tree. (Since Linux 6.6, the default fair scheduler is EEVDF, which builds on the same weighted-fairness idea and adds deadlines to improve latency)

**Answer:** D

**Why:** CFS aims to give every runnable task its **fair share** of the CPU, weighted by priority.
- Each task's `vruntime` increases as it runs, scaled by its weight: a low-nice (high-priority) task's `vruntime` grows more slowly, so it gets picked more often.
- The task with the **smallest** `vruntime` is the one that has received the least of its fair share, so it runs next.
- A red-black tree keeps tasks sorted by `vruntime`, making "find the minimum" fast.

**Why the other options fail:**
- A: CFS has no fixed quantum. Each task's slice depends on how many tasks are runnable and on their weights.
- B: Fairness is the central idea.
- C: CFS doesn't predict bursts.

**Concept recap:** CFS = weighted fair sharing. Lowest `vruntime` runs next. I/O-bound tasks naturally have low `vruntime` (they sleep a lot), so they get the CPU quickly when they wake. This gives good interactivity without the explicit MLFQ rules of Q71.

**Trap alert:** Real-time tasks (`SCHED_FIFO`, `SCHED_RR`) use a separate, strict-priority scheduling class that always runs before CFS/EEVDF tasks.

**Interview angle:** "How does the Linux scheduler work, at a high level?" Weighted fair sharing via virtual runtime, plus separate real-time classes. Follow-up: *"Why use a red-black tree?"* O(log n) insertion and fast retrieval of the minimum.

---

## Q284. [M11 | nice values | Easy | MCQ | Interview]

A user runs `nice -n 10 ./job`. What is the effect?

- A. `job` runs at a lower priority, getting a smaller share of the CPU when other tasks compete. An unprivileged user can raise nice values like this, but can't set negative (higher-priority) values
- B. `job` runs at a higher priority than normal
- C. Nice values affect only disk I/O
- D. `job` is suspended for 10 seconds before it starts

**Answer:** A

**Why:** **nice** values range from **−20 (highest priority)** to **19 (lowest)**, with 0 the default. A higher nice value means "be nice to others": a lower weight, so a smaller CPU share under contention. Ordinary users may only *lower* their own priority. Raising it (a negative nice value) needs privilege (`CAP_SYS_NICE`).

**Why the other options fail:**
- B: Reversed. Higher nice = lower priority.
- C: CPU priority is what `nice` changes. I/O priority has its own tool, `ionice`.
- D: `nice` doesn't delay anything.

**Concept recap:** nice is a weight in CFS/EEVDF (Q283). Each nice step changes the weight by roughly 1.25×, so a few steps make a noticeable difference.

**Trap alert:** On an otherwise idle machine, a niced job still runs at full speed. Nice matters only when tasks **compete** for the CPU.

**Interview angle:** "How do you lower the impact of a background job?" Use `nice`, and `ionice` for its disk I/O. Follow-up: *"Why can't ordinary users set negative nice values?"* They could otherwise grab CPU at everyone else's expense.

---

## Q285. [M11 | System calls that skip the kernel (vDSO) | Medium | MCQ | OA + Interview]

Why can `clock_gettime()` on Linux often complete **without** a real system call?

- A. The compiler computes the time at build time
- B. Reading the time is a privileged instruction, so the call is handled by hardware
- C. The kernel stores the time in each process's heap through `malloc`
- D. The kernel maps a small shared library, the **vDSO**, into every process. Together with a read-only data page that the kernel keeps updated, it lets `clock_gettime()` read the time in user mode, avoiding the trap into the kernel

**Answer:** D

**Why:** A real system call costs a mode switch (Q10), which is significant for functions called millions of times. Reading the current time doesn't need privilege; it just needs data the kernel maintains. So the kernel exposes that data read-only, plus code (the **vDSO**, "virtual dynamic shared object") that computes the time from it, entirely in user space.

**Why the other options fail:**
- A: The time changes constantly, so it can't be computed in advance.
- B: Reading a clock isn't privileged (Q3 C), but the kernel's notion of time still lives in kernel data.
- C: Kernel data can't be placed in the process's `malloc` heap.

**Concept recap:** vDSO = kernel-provided code + data, mapped into user space, for read-only queries such as `clock_gettime` and `gettimeofday`. It appears in `/proc/<pid>/maps` as `[vdso]`.

**Trap alert:** "Every libc function that talks to the kernel makes a system call" isn't true. Q18 shows buffering avoiding calls, and the vDSO avoids them for some queries.

**Interview angle:** "What is the vDSO and why does it exist?" To make frequent, read-only kernel queries as cheap as a function call. Follow-up: *"Why can't `write()` use the vDSO?"* It changes kernel state and needs privileged work.

---

## Q286. [M11 | The /proc file system | Easy | MCQ | Interview]

Which statement about `/proc` on Linux is correct?

- A. `/proc` is a virtual file system: its files are generated by the kernel when they are read, exposing process and kernel information (e.g., `/proc/<pid>/maps`, `/proc/meminfo`), and they take no disk space
- B. `/proc` is a disk partition that stores system logs
- C. `/proc` contains the executables for every running program
- D. Any user can write to any file in `/proc` to change any process

**Answer:** A

**Why:** `/proc` (procfs) is an interface to kernel data that *looks* like a file system. Reading `/proc/meminfo` asks the kernel to describe memory right now. Examples:
- `/proc/<pid>/status`, `/proc/<pid>/maps` (memory mappings), `/proc/<pid>/fd/` (open fds)
- `/proc/cpuinfo`, `/proc/meminfo`, `/proc/buddyinfo` (Q169)
- `/proc/sys/...` for tunable kernel parameters (via `sysctl`)

**Why the other options fail:**
- B: Logs live in `/var/log` or the journal.
- C: `/proc/<pid>/exe` is a *link* to a program's executable, not a copy of every program.
- D: Writable entries exist, but most require root, and the kernel validates every write.

**Concept recap:** "Everything is a file" extends even to kernel state. `ps`, `top` and `free` are largely readers of `/proc`.

**Trap alert:** `/proc` files often report a size of 0 even though reading them returns data, because the content is generated on demand.

**Interview angle:** "How would you inspect a running process's open files and memory?" `/proc/<pid>/fd` and `/proc/<pid>/maps`. Follow-up: *"What is `/sys`?"* sysfs, a similar virtual file system for devices and drivers.

---

## Q287. [M11 | Overcommit, OOM and memory statistics | Medium | MSQ | Interview]

Linux, default settings. **Select all that apply.**

- A. A large `malloc` can succeed even if that much physical memory isn't free, because Linux overcommits by default; physical pages are allocated only when first touched.
- B. When memory is exhausted and reclaim can't free enough, the **OOM killer** chooses a process to kill.
- C. RSS (resident set size) is the physical memory a process is currently using.
- D. A process's virtual size (VSZ) is always smaller than or equal to its RSS.
- E. `/proc/<pid>/maps` lists the process's virtual memory areas: their address ranges, permissions, and backing (file or anonymous).

**Answer:** A, B, C, E

**Why:**
- **A:** This is demand paging (Q170) applied to `malloc`. Reserving address space is cheap, and frames are assigned on first use. Many programs never touch everything they allocate.
- **B:** If that bet loses, the kernel must free memory somehow. The OOM killer selects a victim using a score (mostly how much memory it uses, adjustable with `oom_score_adj`).
- **C, E:** Standard memory statistics and the VMA listing (Q172 B).

**Why the other options fail:**
- D: Reversed. VSZ counts all *mapped* virtual memory, including untouched and swapped-out pages, so **VSZ ≥ RSS**, often by a lot.

**Concept recap:** Virtual memory lets allocation (address space) and use (frames) differ. RSS ≤ VSZ. Overcommit trades the risk of an OOM kill for better memory utilization.

**Trap alert:** "`malloc` returned non-NULL, so the memory definitely exists" is false under overcommit. Failure can come later, as an OOM kill on first touch.

**Interview angle:** "What is the OOM killer?" Follow-up: *"How do you protect a critical service from it?"* Lower its `oom_score_adj`, limit other workloads with cgroups, or change the overcommit policy (`vm.overcommit_memory`).

---

## Q288. [M11 | Containers vs VMs: running another OS | Medium | MCQ | Interview]

You must run a **Windows-only** application, isolated, on a **Linux** host. Which fits, and why?

- A. A Linux container, because containers emulate any OS
- B. A `chroot` jail
- C. A virtual machine: containers share the host's kernel, so on a Linux host they can run only Linux programs, while a Windows application needs a Windows kernel, which a VM provides
- D. A process with a different nice value

**Answer:** C

**Why:** A container is just a set of isolated *processes* on the host kernel (Q267). Its programs make system calls to **the host's** kernel, so they must be built for that kernel. A VM virtualizes the *hardware*, so it can boot any OS, Windows included, with its own kernel.

**Why the other options fail:**
- A: Containers don't emulate operating systems. They share the host kernel.
- B: `chroot` only changes the apparent root directory, for processes on the same kernel.
- D: Nice values have nothing to do with isolation or OS compatibility.

**Concept recap:**

| | Container | VM |
|---|---|---|
| Kernel | shared with the host | its own guest kernel |
| Can run a different OS | no | yes |
| Startup | about a second or less | seconds to minutes |
| Overhead | low | higher |
| Isolation | namespaces + cgroups | hypervisor + virtual hardware |

**Trap alert:** "Docker on Windows or macOS runs Linux containers natively" isn't quite right. Docker Desktop runs them inside a lightweight **Linux VM**, precisely because containers need a Linux kernel.

**Interview angle:** "Containers vs VMs: when do you need a VM?" A different OS or kernel, a custom kernel, or stronger isolation for untrusted tenants. Follow-up: *"What are microVMs like Firecracker?"* Very lightweight VMs that aim for near-container startup times with VM-level isolation.

---

## Q289. [M12 | Scheduling × locks: priority inversion | Medium | MCQ | Interview]

A uniprocessor uses **preemptive priority scheduling**. Three tasks exist: **H** (high), **M** (medium) and **L** (low).
- L locks mutex R.
- H wakes up, preempts L, tries to lock R, and blocks.
- M (CPU-bound, never touches R) then becomes ready.

What happens, and what is the standard fix?

- A. H runs immediately, because it has the highest priority
- B. L always finishes first, because holding a lock prevents preemption
- C. M and H deadlock
- D. M preempts L, so L can't run to release R, and H ends up waiting behind M, a lower-priority task, for an unbounded time (**priority inversion**). Fix: **priority inheritance**, where L temporarily runs at H's priority until it releases R

**Answer:** D

**Why:** H is blocked on a lock that only L can release. M outranks L, so M runs instead of L. In effect, a **medium**-priority task is delaying a **high**-priority one indefinitely. With priority inheritance, L inherits H's priority while H waits for R. M can no longer preempt L, so L finishes its critical section and releases R, and H proceeds.

**Why the other options fail:**
- A: H is *blocked* on R, so its priority doesn't help it run.
- B: An ordinary mutex doesn't stop preemption. (A spinlock in the kernel does, Q91 B, but that's a different mechanism.)
- C: There's no cycle. H waits for L, and L waits for the CPU, which M holds. The wait is unbounded, but it isn't a deadlock.

**Concept recap:** Priority inversion = high waits on low, while medium runs. Fixes: priority inheritance (POSIX `PTHREAD_PRIO_INHERIT`), or the priority-ceiling protocol. This is the classic interaction between scheduling (M3) and synchronization (M4).

**Trap alert:** It takes **three** tasks to make this unbounded. With only H and L, H waits just for L's critical section, which is bounded. Q100 A showed a related but different danger: spinning instead of blocking.

**Interview angle:** "What is priority inversion?" The NASA Mars Pathfinder lander (1997) suffered repeated resets from it, fixed by enabling priority inheritance remotely. Follow-up: *"Inheritance vs ceiling?"* Inheritance raises the holder's priority only when a higher task actually blocks. Ceiling raises it as soon as the lock is taken, to the highest priority of any task that might use the lock.

---

## Q290. [M12 | EMAT with TLB misses and page faults | Hard | Numerical | OA + Interview]

- TLB lookup: **20 ns**. Memory access: **100 ns**. TLB hit ratio: **80%**. Single-level page table.
- A TLB hit means the page is resident.
- On a TLB miss, **0.1% of those accesses** also page-fault. A faulting access costs **10 ms** in total.
- A TLB miss without a fault costs TLB lookup + page-table read + data access.

What is the effective memory access time?

- A. 140 ns
- B. ≈ 2.0 ms
- C. ≈ 2.14 µs
- D. 10 ms

**Answer:** C

**Why:** Three outcomes are possible, each with its own cost. Weight each one by its probability. Even a tiny fault rate dominates, because a fault is about 100,000× more expensive than a memory access.

**Working:**
- **TLB hit:** 20 + 100 = **120 ns**, with probability 0.8.
- **TLB miss, no fault:** 20 + 100 + 100 = **220 ns**, with probability 0.2 × 0.999.
- **TLB miss + fault:** **10,000,000 ns**, with probability 0.2 × 0.001.
- EMAT = 0.8 × 120 + 0.2 × (0.999 × 220 + 0.001 × 10,000,000)
  = 96 + 0.2 × (219.78 + 10,000)
  = 96 + 0.2 × 10,219.78
  = 96 + 2043.96 ≈ **2140 ns ≈ 2.14 µs**.

**Why the other options fail:**
- A: Ignores page faults entirely: 0.8 × 120 + 0.2 × 220 = 140 ns.
- B: A units slip (µs read as ms).
- D: The cost of one fault, not the average.

**Concept recap:** Nested expectation: EMAT = h·hit cost + (1 − h)·[(1 − q)·miss cost + q·fault cost]. The fault term dominates, so the TLB hit ratio barely matters once faults happen at all.

**Trap alert:** Apply the fault probability **only to TLB misses** when the question says so. Here, a fault rate of 0.1% of misses = 0.02% of all accesses.

**Interview angle:** "What dominates memory performance: TLB misses or page faults?" Page faults, by orders of magnitude, whenever they occur. Follow-up: *"So how do you tune it?"* First keep the working set resident (no faults), then improve TLB reach (huge pages).

---

## Q291. [M12 | fork × private vs shared memory | Hard | Code-trace | OA + Interview]

```c
int x = 0;                                   // ordinary global: private after fork (COW)
int *p = mmap(NULL, 4096, PROT_READ | PROT_WRITE,
              MAP_SHARED | MAP_ANONYMOUS, -1, 0);   // shared page
*p = 0;
if (fork() == 0) {          // child
    x = 5;
    *p = 5;
    _exit(0);
}
wait(NULL);
printf("%d %d\n", x, *p);   // parent
```

Assume all calls succeed. What does the parent print?

- A. `0 0`
- B. `5 5`
- C. `0 5`
- D. `5 0`

**Answer:** C

**Why:** The two variables live in two different kinds of memory:
- `x` is ordinary process memory. After `fork`, it's copy-on-write (Q207). The child's write creates the **child's own copy**, so the parent's `x` stays **0**.
- `*p` is in a **`MAP_SHARED`** mapping. Both processes map the **same physical frame**, so the child's write is visible to the parent: **5**.

**Working:**
1. Before `fork`: x = 0 (private page), *p = 0 (shared frame F).
2. Child: `x = 5` → COW fault → child gets a private copy; the parent's page is untouched. `*p = 5` → writes into frame F, which both map.
3. Parent, after `wait`: reads its own x = **0** and frame F's value = **5**.

**Why the other options fail:**
- A: Treats the shared mapping as private.
- B: Treats ordinary memory as shared (the Q26 misconception).
- D: Reversed.

**Concept recap:** After `fork`, ordinary memory is private (COW), while `MAP_SHARED` mappings stay shared. Shared anonymous `mmap` is the simplest way for related processes to share memory (Q41).

**Trap alert:** `MAP_PRIVATE` would behave like `x` and give `0 0`. The mapping flag, not the use of `mmap`, decides whether memory is shared.

**Interview angle:** "How can a parent and child share a counter?" Shared memory (`MAP_SHARED` or POSIX shm), plus synchronization such as a process-shared mutex or atomics, because it's still a race (Q77). Follow-up: *"Can threads skip all this?"* Yes. Threads already share their whole address space (Q33).

---

## Q292. [M12 | Deadlock × memory allocation | Easy | MCQ | OA]

A memory pool has **8 frames**. Each process needs **3 frames** from the pool to finish. It acquires them one at a time and holds them until it finishes. What is the **largest** number of such processes that can run concurrently with deadlock guaranteed impossible?

- A. 3
- B. 2
- C. 4
- D. 8

**Answer:** A

**Why:** This is the resource-count condition from Q134, applied to memory frames.
- Deadlock-free ⇔ n(k − 1) + 1 ≤ R, with k = 3 and R = 8: 2n + 1 ≤ 8 → n ≤ 3.5 → **n = 3**.
- **Check with 4 processes:** each could hold 2 frames (4 × 2 = 8, the whole pool) while waiting for a third. No one can finish, so they deadlock.
- **Check with 3 processes:** at most 3 × 2 = 6 frames are held while waiting, leaving 2 free. Someone gets a third frame, finishes, and frees all 3.

**Why the other options fail:**
- B: Safe, but not the maximum.
- C: Can deadlock, as shown above.
- D: Ignores the hold-and-wait problem.

**Concept recap:** n_max = ⌊(R − 1) / (k − 1)⌋ = ⌊7/2⌋ = 3.

**Trap alert:** Memory allocation can deadlock just like locks. This is why kernels reserve emergency memory pools for code that must allocate in order to *free* memory, such as writeback.

**Interview angle:** Follow-up: *"How do real systems avoid this?"* Allocate everything up front (breaking hold-and-wait), use timeouts, or preempt by swapping a process out.

---

## Q293. [M12 | Inherited sockets × exec | Hard | MCQ | Interview]

A web server listens on port 8080. For each request it `fork`s a child, which `exec`s a script. The listening socket was **not** opened with close-on-exec. One script keeps running for hours. The administrator restarts the server. What is the likely result?

- A. The restart works normally
- B. The script automatically receives SIGKILL
- C. The script's memory is freed
- D. The new server may fail to bind port 8080 (`EADDRINUSE`), and clients can still connect to the old listening socket, where nobody accepts them. The long-running script holds an **inherited fd** for that socket. Fix: create the listening socket with `SOCK_CLOEXEC` (or set `FD_CLOEXEC`)

**Answer:** D

**Why:** Each piece comes from a different module:
- **fork** copies the fd table, so the child holds the listening socket (Q275).
- **exec** keeps open fds unless they're marked close-on-exec (Q27, Q280).
- A socket stays open, still listening on 8080, as long as **any** process holds an fd for it (Q223). Killing the server doesn't close the copy the script holds.
- So the new server's `bind` conflicts with that still-listening socket.

**Why the other options fail:**
- A: The script's inherited fd keeps the old socket alive.
- B: Nothing signals the script automatically.
- C: Memory is irrelevant here. The problem is the fd.

**Concept recap:** Rule of thumb: **open every fd with close-on-exec**, unless it is explicitly meant to be passed to the child.

**Trap alert:** This bug shows up only when a child runs a **long time** after its parent exits. In short tests, children finish quickly, so it never appears.

**Interview angle:** "Why might a service fail to restart with 'Address already in use'?" A leftover process holds the socket, often an inherited fd, or the address is stuck in TIME_WAIT. Follow-up: *"How do you find the culprit?"* `ss -ltnp` or `lsof -i :8080`.

---

## Q294. [M12 | Dirty pages × latency spikes | Medium | MCQ | Interview]

An application quickly writes **20 GB** into the page cache. The system then comes under memory pressure, and **unrelated** processes suddenly see latency spikes. Why?

- A. The CPU scheduler stops every process
- B. Clean pages are lost
- C. Reclaim can't reuse dirty pages until they've been written back, so memory allocations stall behind a burst of writeback I/O. Linux limits this by throttling writers once dirty memory crosses thresholds such as `vm.dirty_background_ratio` and `vm.dirty_ratio`
- D. The page cache is never written back

**Answer:** C

**Why:**
- `write()` returns once data reaches the page cache (Q236), so 20 GB of dirty pages can pile up quickly.
- To reuse a dirty frame, the kernel must first write it out (Q189, Q198). Under memory pressure, other processes' allocations end up waiting on that disk I/O.
- The kernel limits the damage by starting background writeback early (`dirty_background_ratio`) and making heavy writers wait once dirty memory reaches `dirty_ratio`.

**Why the other options fail:**
- A: The scheduler keeps running. Processes are blocked on memory and I/O, not stopped.
- B: Clean pages can be dropped instantly and harmlessly. They aren't the problem.
- D: Writeback does happen. The issue is *when*, and how much at once.

**Concept recap:** This combines write-back caching (M8), page replacement with dirty bits (M7), and disk throughput (M9). The more dirty memory there is, the longer reclaim takes.

**Trap alert:** "Lots of free RAM in the page cache means memory is fine" is wrong when that cache is **dirty**, because dirty pages aren't immediately reusable.

**Interview angle:** "Why does a big file copy sometimes freeze a desktop?" Writeback and reclaim stalls. Follow-up: *"Mitigations?"* Lower the dirty thresholds, call `fsync` periodically in the writer, use `O_DIRECT` for bulk copies, or limit the writer with cgroups.

---

## Q295. [M12 | File layout × disk access time | Medium | Numerical | OA + Interview]

An HDD has an average seek of **4 ms**, an average rotational latency of **4 ms**, and a transfer time of **0.04 ms per block**. A **100-block** file is read from start to finish. How long does it take if the file is stored **contiguously** (one seek and one rotational delay, then every block in sequence), and if it is **scattered** (every block needs its own seek and rotational delay)?

- A. 12 ms and ≈ 804 ms
- B. 8.04 ms and ≈ 804 ms
- C. 12 ms and 12 ms
- D. ≈ 804 ms and 12 ms

**Answer:** A

**Why:** For contiguous blocks, the mechanical delay is paid **once**. For scattered blocks, it's paid **for every block**, and it's 200× larger than the transfer itself.

**Working:**
- **Contiguous:** 4 + 4 + 100 × 0.04 = 8 + 4 = **12 ms**.
- **Scattered:** 100 × (4 + 4 + 0.04) = 100 × 8.04 = **804 ms**.
- Ratio ≈ **67×** slower when scattered.

**Why the other options fail:**
- B: Charges the contiguous case for only one block of transfer.
- C: Ignores the per-block seek and rotation when scattered.
- D: Swapped.

**Concept recap:** Sequential HDD I/O ≈ one positioning cost + n transfers. Random HDD I/O ≈ n × (positioning + transfer). This is why file systems use extents and try to keep files contiguous (Q215, Q218), and why fragmentation hurts HDDs much more than SSDs (Q251).

**Trap alert:** On an SSD, the gap shrinks dramatically, because there's no seek or rotation. Layout matters far less there.

**Interview angle:** "Why do databases prefer sequential I/O?" Logs, LSM trees and large sequential scans avoid paying the positioning cost per record. Follow-up: *"How does read-ahead help?"* It turns predictable sequential reads into fewer, larger I/Os.

---

## Q296. [M12 | Races × multicore | Hard | MSQ | OA + Interview]

A race condition shows up regularly on a 4-core machine, but never in testing on a 1-core VM. **Select all that apply.**

- A. Races are still possible on a single core, because preemption can interleave threads at any point.
- B. True parallelism makes some bad interleavings far more likely, so race bugs surface more often on multicore machines.
- C. On multicore hardware, memory reordering (for example store buffers) can expose races that sequentially consistent execution would hide.
- D. Adding a `printf` that makes the bug disappear is not a fix: it only changes the timing (a "heisenbug").
- E. Running on a single core makes races impossible.

**Answer:** A, B, C, D

**Why:**
- **A:** The lost update in Q77 needs only a context switch in the middle of `count++`, which a timer interrupt can cause on one core.
- **B:** On one core, interleavings happen only at preemption points, which are rare. On many cores, instructions truly overlap all the time.
- **C:** Peterson's algorithm failing on real hardware (Q82 C) is exactly this: a later read overtakes an earlier write.
- **D:** Print statements add delays and implicit synchronization, so the bad window closes, but only by accident.

**Why the other options fail:**
- E: Directly contradicted by A.

**Concept recap:** A race is a **correctness** bug that timing merely **reveals**. The fix is synchronization (locks, atomics, barriers), never a change of timing.

**Trap alert:** "It passed 10,000 test runs" proves nothing about a race (Q78). Use ThreadSanitizer, stress tests on many cores, and design reviews.

**Interview angle:** "Your bug only reproduces in production. How do you approach it?" Suspect races or memory ordering, run ThreadSanitizer, add assertions, and stress-test on a machine like production's. Follow-up: *"What does `volatile` do here?"* Nothing useful. It neither makes operations atomic nor orders them between threads (Q77).

---

## Q297. [M12 | Memory × isolation × security | Medium | MSQ | Interview]

**Select all that apply.**

- A. ASLR (address-space layout randomization) places the stack, heap and libraries at random addresses, making memory-corruption attacks much harder to aim.
- B. Guard pages (deliberately unmapped pages) below a thread's stack turn a stack overflow into an immediate fault, instead of silently corrupting adjacent memory.
- C. KPTI (kernel page-table isolation) removes most kernel mappings from user-mode page tables, to mitigate Meltdown-style leaks of kernel memory.
- D. The NX bit makes the stack executable.
- E. Because each process has its own page table, a stray pointer in one process cannot corrupt another process's memory.

**Answer:** A, B, C, E

**Why:** Each defence uses the paging hardware (M6/M7) for protection (M10):
- **A:** Randomizing the layout means an attacker doesn't know where useful code or data is.
- **B:** An unmapped page faults on any access (Q172), so an overflow is caught at its first step past the end.
- **C:** If kernel memory isn't mapped while user code runs, speculative-execution attacks have much less to read. The cost is extra page-table switching on system calls (the KPTI note in Q10).
- **E:** Isolation by construction (Q151 D).

**Why the other options fail:**
- D: Reversed. NX marks the stack and heap **non**-executable (Q152, Q266).

**Concept recap:** Paging is a **security** mechanism as much as a memory-management one: isolation, W^X, guard pages, ASLR and KPTI all rely on the page tables.

**Trap alert:** Per-process isolation does **not** stop a bug from corrupting the process's **own** memory. That's why ASLR, NX and stack canaries exist as additional layers.

**Interview angle:** "How does the OS make exploiting memory bugs harder?" Name ASLR, NX/W^X, guard pages and stack canaries, and add KPTI for kernel memory. Follow-up: *"What did Meltdown change about system-call cost?"* KPTI adds a page-table switch on entry and exit. PCID tagging (Q165 B) reduces the TLB cost of that switch.

---

## Q298. [M12 | A blocking read, end to end | Easy | MCQ | Interview]

A process calls `read()` on a socket that has no data yet. Put these events in order:

1. The process blocks on the socket's wait queue
2. The system call traps into the kernel (a mode switch)
3. A packet arrives; the network card raises an interrupt, and the kernel marks the process ready
4. The scheduler context-switches to another process
5. Later, the scheduler dispatches the process; `read()` copies the data and returns to user mode

- A. 1, 2, 4, 3, 5
- B. 2, 4, 1, 3, 5
- C. 2, 1, 3, 4, 5
- D. 2, 1, 4, 3, 5

**Answer:** D

**Why:** This one `read()` touches five modules:
- **Trap (M1):** the call enters the kernel (2).
- **Block (M2):** with no data, the process moves to a wait queue (1, Q23).
- **Schedule (M3):** the CPU is given to another process (4).
- **Interrupt (M9):** data arrives and the handler wakes the waiting process (3). It becomes **Ready**, not Running (Q20).
- **Dispatch:** eventually the scheduler runs it, `read` finishes, and control returns to user mode (5).

**Why the other options fail:**
- A: The process can't block before it has even entered the kernel.
- B: The scheduler can't switch away before the process has blocked, because nothing has given up the CPU.
- C: Has the packet arrive before the CPU is handed to anyone else. Usually other work runs while the process waits.

**Concept recap:** Running → (trap) kernel → Waiting → (context switch) … (interrupt) → Ready → (dispatch) Running → return to user mode.

**Trap alert:** A process woken by an interrupt goes to **Ready**. It runs only when dispatched (Q20).

**Interview angle:** "Walk me through what happens during a blocking network read." Use exactly this sequence. Follow-up: *"Where do the copies happen?"* NIC → kernel buffer (via DMA), then kernel → user buffer during `read`.

---

## Q299. [M12 | First touch of an mmap'd file page | Hard | MSQ | OA + Interview]

A process accesses a page of an `mmap`'d file for the first time. **Select all that can happen.**

- A. A TLB miss.
- B. The page walk finds an invalid PTE, so the MMU raises a page fault (a trap into the kernel).
- C. If the page is already in the page cache (another process read it), the kernel just maps that frame: a minor fault.
- D. If the page isn't cached, the kernel issues a disk read and blocks the process, and the scheduler runs something else meanwhile.
- E. After the page is loaded, the faulting instruction is skipped, and execution continues at the next instruction.

**Answer:** A, B, C, D

**Why:**
- **A, B:** There's no translation yet, so the TLB misses and the page walk finds an invalid PTE (Q172). The fault is legitimate, because the address lies inside a mapped region.
- **C:** Minor fault, with no I/O (Q179).
- **D:** Major fault. The disk read blocks this process, and the scheduler keeps the CPU busy with others (Q171, Q298).

**Why the other options fail:**
- E: After the page is in, the faulting instruction is **re-executed** (Q7 A), which is why it succeeds the second time.

**Concept recap:** First touch of an `mmap`'d page = TLB miss → page walk → page fault → (cache hit: map it | cache miss: I/O + block) → update PTE → restart the instruction → TLB filled.

**Trap alert:** Under `MAP_PRIVATE`, the first *write* to a page causes another fault later, for copy-on-write (Q207), even if the page has already been read.

**Interview angle:** "What happens, step by step, when a program first reads an `mmap`'d file?" Follow-up: *"Why can the second process that reads the same file start faster?"* Its accesses are minor faults served from the page cache.

---

## Q300. [M12 | Capstone: `cat file.txt \| grep x` | Medium | MCQ | Interview]

A shell runs `cat file.txt | grep x`. Which statement is **FALSE**?

- A. The shell creates a pipe, forks two children, and uses `dup2` so that `cat`'s stdout is the pipe's write end and `grep`'s stdin is its read end
- B. While `cat` waits for a disk read, the scheduler can run `grep` (or anything else); if the pipe fills up, `cat` blocks until `grep` reads
- C. `grep` sees end-of-file only after `cat` exits and every write end of the pipe has been closed
- D. `cat` and `grep` share one address space, because they're connected by a pipe

**Answer:** D

**Why:** A, B and C are each correct, and each draws on a different module:
- **A:** Process creation and fd plumbing (Q274, Q278).
- **B:** Blocking I/O and scheduling (Q298), and pipe flow control. A full pipe blocks the writer, and an empty pipe blocks the reader (Q42, Q48).
- **C:** EOF semantics (Q278).

D is false. `cat` and `grep` are **separate processes** with **separate address spaces**. A pipe is **message passing** through a kernel buffer (Q41), not shared memory.

**Why the other options fail:** A, B and C are true statements, so they aren't the answer.

**Concept recap:** A shell pipeline uses, all at once: process creation (fork/exec), the fd table (dup2), IPC (pipes), scheduling (blocking and waking), and file systems and I/O (reading `file.txt` through the page cache).

**Trap alert:** "Connected processes share memory" is the misconception. Data crosses the pipe by being **copied** into and out of a kernel buffer.

**Interview angle:** "Explain everything that happens when you run `cat file.txt | grep x`." This is the classic whole-system OS interview question. Walk through the shell, `fork`, `pipe`, `dup2`, `exec`, `read` (disk or page cache), `write` into the pipe, scheduling, EOF and `wait`. Follow-up: *"Why might `grep` start printing before `cat` finishes?"* The pipe streams data. `grep` processes each chunk as it arrives.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 281 | C | 286 | A | 291 | C | 296 | A, B, C, D |
| 282 | A, B, C, D | 287 | A, B, C, E | 292 | A | 297 | A, B, C, E |
| 283 | D | 288 | C | 293 | D | 298 | D |
| 284 | A | 289 | D | 294 | C | 299 | A, B, C, D |
| 285 | D | 290 | C | 295 | A | 300 | D |

**Completed: 300 / 300    Remaining: 0** ✅

---
---

# 📊 Progress Report #3 (after Q300)

## 1. Topics completed

| Module | Questions | Status |
|---|---|---|
| M1 OS fundamentals & architecture | Q1–18 | ✅ |
| M2 Processes, threads, IPC | Q19–48 | ✅ |
| M3 CPU scheduling | Q49–76 | ✅ |
| M4 Synchronization & concurrency | Q77–114 | ✅ |
| M5 Deadlocks | Q115–134 | ✅ |
| M6 Memory management | Q135–169 | ✅ |
| M7 Virtual memory & page replacement | Q170–209 | ✅ |
| M8 File systems | Q210–237 | ✅ |
| M9 I/O, disks, storage | Q238–259 | ✅ |
| M10 Protection & security | Q260–271 | ✅ |
| M11 Linux/Unix internals | Q272–288 | ✅ |
| M12 Integrated & advanced | Q289–300 | ✅ |

## 2. Key concepts covered in Q201–Q300

- **Virtual memory (finished):** working sets, thrashing and its feedback loop, PFF, proportional allocation, global vs local replacement, copy-on-write, `mmap` shared vs private, the optimal page size.
- **File systems:** inodes (and what they don't store), `open()` and the three file tables, access methods, directory structures, contiguous/linked/FAT/indexed allocation, inode sizing, links and link counts, free-space bitmaps, permissions (including on directories), mounting, journaling and crash consistency, the page cache and `fsync`, the superblock.
- **I/O and storage:** MMIO, polling vs interrupts, DMA, buffering, caching and spooling, drivers, HDD access time, all six disk schedulers, RAID levels with parity and the small-write penalty, SSD internals, readiness vs completion I/O.
- **Protection and security:** authN vs authZ, the access matrix, ACLs vs capabilities, domains and setuid, least privilege, hardware memory protection, containers vs VMs, meta-rights, privilege escalation as a concept, mechanism vs policy.
- **Linux internals:** fds 0/1/2, the single task model, `fork`/`exec`/`wait` in the shell, shared offsets after `fork`, `dup2` and pipe plumbing, close-on-exec, signals and process groups, CFS/EEVDF and `nice`, the vDSO, `/proc`, overcommit and the OOM killer.
- **Integration:** priority inversion; EMAT with both TLB misses and faults; shared vs private memory after `fork`; memory-pool deadlock; inherited sockets; writeback stalls; layout vs access time; races on multicore; paging as a security mechanism; a blocking read end to end; a first-touch `mmap` fault; the shell pipeline capstone.

## 3. Misconceptions tested in Q201–Q300

| # | Misconception | Where |
|---|---|---|
| 67 | The working set counts repeated references | Q201 |
| 68 | Shrinking every process evenly cures thrashing | Q202 |
| 69 | Low CPU utilization ⇒ admit more processes | Q203 |
| 70 | Local replacement prevents thrashing | Q206 |
| 71 | `fork` physically copies all memory | Q207 |
| 72 | `mmap` reads the whole file at once | Q208 |
| 73 | Smaller (or larger) pages are always better | Q209 |
| 74 | The inode stores the file name | Q210 |
| 75 | Linked allocation is good for random access | Q218 |
| 76 | Reaching block i under linked allocation = i reads | Q219 |
| 77 | The indirect ranges start at block 0 (forgetting the direct blocks) | Q222 |
| 78 | Separate `open()`s share an offset; `fork` gives an independent one | Q223, Q275 |
| 79 | A symlink keeps the target's data alive | Q225 |
| 80 | A directory's link count includes its files | Q227 |
| 81 | Bits vs bytes in bitmap sizing | Q228 |
| 82 | Deleting a file needs write permission on the file itself | Q231 |
| 83 | Mounting deletes or merges the underlying files | Q232 |
| 84 | Journaling guarantees no recent data is lost | Q234 |
| 85 | Writing the inode first is the safest order | Q235 |
| 86 | `write()` returning ⇒ the data is on disk | Q236 |
| 87 | Polling is always wasteful | Q239 |
| 88 | DMA removes the need for interrupts | Q240 |
| 89 | A full rotation is the average rotational latency | Q245 |
| 90 | SSTF is optimal and never starves | Q247, Q250 |
| 91 | SCAN and LOOK turn at the same place; direction doesn't matter | Q248 |
| 92 | C-SCAN/C-LOOK serve requests on the return | Q249 |
| 93 | RAID 0 is redundant; RAID replaces backups | Q253, Q257 |
| 94 | A RAID 5 small write costs 1 or 2 I/Os | Q256 |
| 95 | Random small writes are fastest on SSDs | Q258 |
| 96 | epoll = asynchronous I/O | Q259 |
| 97 | Protection = security | Q260 |
| 98 | Unix permission bits are capabilities | Q263 |
| 99 | setuid gives permanent root | Q264 |
| 100 | Containers isolate as strongly as VMs, or can run any OS | Q267, Q288 |
| 101 | Authorized elevation (sudo, setuid passwd) = privilege escalation | Q269 |
| 102 | A PID is reused while the process is still a zombie | Q273 |
| 103 | Unused pipe ends don't need closing | Q278 |
| 104 | Close-on-exec affects `fork` | Q280 |
| 105 | VSZ ≤ RSS; `malloc` success ⇒ memory exists | Q287 |
| 106 | Priority inversion is a deadlock | Q289 |
| 107 | A single core makes races impossible; a `printf` "fixes" them | Q296 |
| 108 | The faulting instruction is skipped after a fault | Q299 |
| 109 | Processes connected by a pipe share memory | Q300 |

## 4. Formulas introduced in Q201–Q300

| Formula | Where |
|---|---|
| WS(t, Δ) = distinct pages in the last Δ references; thrashing if D = ΣWSS > m | Q201–202 |
| Proportional allocation aᵢ = (sᵢ / S) × m | Q205 |
| Optimal page size p = √(2·s·e) | Q209 |
| Contiguous: disk block = start + ⌊offset / B⌋; offset in block = offset mod B | Q215 |
| Reads to reach block i: contiguous 1; linked i + 1; indexed (levels + 1) | Q219 |
| Inode max size = (D + k + k² + k³)·B, where k = B / pointer size | Q221 |
| Directory link count = 2 + number of subdirectories | Q227 |
| Bitmap bytes = (disk size / block size) / 8 | Q228 |
| HDD access = seek + ½·(60/RPM) + block size / rate | Q245 |
| Disk capacity = surfaces × tracks × sectors × bytes per sector | Q246 |
| SCAN/LOOK/C-LOOK head-movement recipes | Q248–249 |
| RAID capacity: 0 → nS; 5 → (n − 1)S; 6 → (n − 2)S; 10 → (n/2)S | Q254 |
| Parity: P = ⊕Dᵢ; missing block = ⊕(survivors) | Q255 |
| Small-write I/Os: RAID 0: 1; 1/10: 2; 5: 4; 6: 6 | Q256 |
| Nested EMAT = h·hit + (1 − h)[(1 − q)·miss + q·fault] | Q290 |
| Max processes without deadlock = ⌊(R − 1)/(k − 1)⌋ | Q292 |
| Sequential vs random: 1 positioning + n transfers vs n × (positioning + transfer) | Q295 |

## 5. Areas to revise before the mock OAs

1. **Inode and file-allocation numericals (Q219, Q221–222, Q228):** be able to state the pointer ranges, read counts and bitmap size without notes.
2. **Disk scheduling (Q247–249):** redo all six algorithms on a fresh queue. Write down the direction and whether the return jump counts *before* you start.
3. **RAID (Q254–256):** capacity, XOR reconstruction and the small-write penalty for each level.
4. **fd plumbing (Q223, Q275–280):** separate opens vs shared offsets, `dup2`, closing pipe ends, close-on-exec. Practise writing `ls | wc -l` from memory.
5. **The integrated questions (Q289–300):** for each one, say aloud which modules it combines. That's exactly how interviewers expect you to reason.
6. **Re-attempt everything flagged in progress reports #1 and #2.**

## 6. Final bank statistics

| Measure | Planned | Delivered |
|---|---|---|
| Total questions | 300 | 300 |
| MCQ / MSQ / Numerical-or-code-trace | 179 / 75 / 46 | 177 / 75 / 48 |
| MSQs (cap 80) | ≤ 80 | 75 |
| "All / none of the above" (cap 5) | ≤ 5 | 0 |
| Correct-answer spread over 225 single-answer questions | ~56 each | A 54 · **B 61** · C 55 · D 55 |

The correct answers lean slightly toward **B** (61 of 225, about 27% instead of 25%). Most of that excess comes from Batch 14, where B was the answer far too often. Don't use "it's probably B" as a strategy: the lean is small, and the mock OAs contain no such pattern you can rely on.
