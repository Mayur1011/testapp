# OS MCQ/MSQ Bank — Batch 2 (Q21–Q40)

**What this batch covers (all M2):** the PCB, `fork` return values, process queues, `fork` tracing (loops, short-circuit evaluation, variables, stdio buffers), `exec`, `wait`/`waitpid`, what a child inherits, zombies and orphans, process vs thread, user-level vs kernel-level threads, threading models, multithreading benefits, concurrency vs parallelism, and `fork` in a multithreaded process.

> In all code questions, assume the standard headers (`stdio.h`, `unistd.h`, `sys/wait.h`) are included and Linux + glibc.

---

## Q21. [M2 | PCB | Easy | MCQ | OA]

Which of the following is NOT stored in a process's PCB?

- A. The process state (ready, running, waiting, …)
- B. The saved program counter and registers, while the process is not running
- C. The full contents of the process's heap
- D. Pointers to its memory-management information (page-table base) and its open-file table

**Answer:** C

**Why:** The **PCB** (Process Control Block) is the kernel's small record *about* a process. It holds metadata and pointers, not the process's actual memory. The heap lives in the process's own pages. The PCB stores only where to find them, through the page-table pointer.

**Why the other options fail:**
- A: The state is needed by the scheduler.
- B: The saved CPU state is what lets a process resume after a context switch.
- D: Memory and file information are core PCB fields.

**Concept recap:** PCB fields, grouped:
- **Identity:** PID, PPID, UID
- **CPU state:** PC, registers
- **Scheduling:** priority, pointers into queues
- **Memory:** page-table base, limits
- **I/O:** open files
- **Accounting:** CPU time used

On Linux, the PCB is `struct task_struct`.

**Trap alert:** "PCB = everything about the process" leads students to put data or code in it. It holds only bookkeeping.

**Interview angle:** "What's in a PCB?" Answer with the six groups above. Follow-up: *"Where does the PCB live?"* In kernel memory, which the process itself cannot read.

---

## Q22. [M2 | fork() return values | Easy | MCQ | OA + Interview]

What does `fork()` return?

- A. 0 in the parent; the child's PID in the child
- B. The child's PID in the parent; 0 in the child; −1 in the parent (and no child is created) on failure
- C. The same new PID in both processes
- D. The parent's PID in the child; 0 in the parent

**Answer:** B

**Why:** After a successful `fork` there are two processes running the same code. The return value is the only way each one knows which it is.
- The parent needs the child's PID, so it can `wait` for or `kill` it later.
- The child gets 0. It can find its own PID with `getpid()` and its parent's with `getppid()`.
- 0 is never the PID of a user process, so this is unambiguous.

**Why the other options fail:**
- A: Parent and child are swapped. This is the most common mistake.
- C: Two processes cannot share a PID.
- D: The child does not get the parent's PID from `fork`.

**Concept recap:** `fork()` → **> 0** (child's PID) in the parent · **0** in the child · **−1** on failure (only the parent exists).

**Trap alert:** On success, `fork` "returns twice", once in each process. On failure it returns once, because there is no child.

**Interview angle:** "How does the same code know whether it's the parent or the child?" Follow-up: *"Why not just give the child its own PID?"* It can call `getpid()` itself. 0 is a cheap, unambiguous marker.

---

## Q23. [M2 | Scheduling queues | Easy | MCQ | OA]

Process P calls `read()` on a disk file. The data is not cached, so P must wait. Where is P placed?

- A. The job queue
- B. The ready queue
- C. The ready queue of another CPU
- D. The device (wait) queue for that disk, i.e. the wait queue for that I/O event

**Answer:** D

**Why:** A blocked process cannot use the CPU, so it leaves the **ready queue** (processes ready to run) and waits on the **device queue** for its I/O. When the disk finishes, the interrupt handler (or code it triggers) wakes P and moves it back to the ready queue.

**Why the other options fail:**
- A: In the textbook model, the job queue holds every process in the system. It is not the place for waiting on an event.
- B: A blocked process must not be picked by the scheduler, so it can't stay in the ready queue.
- C: Moving between CPUs is load balancing. It doesn't apply to a blocked process.

**Concept recap:**
- Job queue = all processes.
- Ready queue = in memory and runnable.
- Device/wait queues = blocked on an event.

On Linux, wait queues are attached to the specific event being waited for.

**Trap alert:** A blocked process does **not** sit in the ready queue waiting its turn.

**Interview angle:** "What queues does an OS maintain?" Follow-up: *"Who moves the process back when I/O completes?"* The I/O-completion path, started by the device interrupt.

---

## Q24. [M2 | fork() in a loop | Medium | Code-trace | OA]

```c
int main(void) {
    for (int i = 0; i < 3; i++)
        fork();
    printf("x\n");
    return 0;
}
```

Assume every `fork` succeeds. How many **child processes** are created, not counting the original?

- A. 7
- B. 8
- C. 3
- D. 6

**Answer:** A

**Why:** Every process alive at a `fork()` call executes it, and so does every child born in an earlier iteration, because it continues the same loop. The number of processes therefore doubles at each iteration.

**Working:**
- **Given:** 3 iterations; each process alive at an iteration forks once.
- **Formula:** n sequential forks executed by everyone → total processes = 2ⁿ; children = 2ⁿ − 1.

| After iteration | Processes alive |
|---|---|
| start | 1 (P) |
| i = 0 | 2 |
| i = 1 | 4 |
| i = 2 | 8 |

Cross-check with the process tree:
- P forks at i = 0, 1, 2 → **3** children (C1, C2, C3).
- C1 (born at i = 0) forks at i = 1, 2 → **2** children (C4, C5).
- C2 (born at i = 1) forks at i = 2 → **1** child (C6).
- C4 (born at i = 1 inside C1) forks at i = 2 → **1** child (C7).

Total = 3 + 2 + 1 + 1 = **7** children. All 8 processes print "x".

**Why the other options fail:**
- B: 8 is the total number of processes, including the original. Catches 2ⁿ vs 2ⁿ − 1.
- C: Counts one child per `fork` written in the code, ignoring that the children fork too.
- D: 2n, which is a linear guess.

**Concept recap:** Total = 2ⁿ, children = 2ⁿ − 1, and the output line is printed 2ⁿ times.

**Trap alert:** Read the question carefully: "lines printed" (8) vs "children created" (7) vs "total processes" (8).

**Interview angle:** "How many processes does this create? Draw the tree." Follow-up: *"How do you make the parent alone create exactly n children?"* Have each child leave the loop: `if (fork() == 0) break;` (or call `exit`).

---

## Q25. [M2 | fork() with short-circuit evaluation | Hard | Code-trace | OA]

```c
int main(void) {
    fork() && fork() || fork();
    printf("Z\n");
    return 0;
}
```

Rules: in C, `&&` binds tighter than `||`, and both **short-circuit** (the right side is evaluated only if needed). Assume every `fork` succeeds. How many times is "Z" printed?

- A. 4
- B. 6
- C. 5
- D. 8

**Answer:** C

**Why:** The expression is `(F1 && F2) || F3`. In a parent, `fork` returns a PID > 0, which is *true*. In a child it returns 0, which is *false*. Each new child resumes evaluating the expression from the point where its `fork` returned.

**Working:**

| Process | F1 value | F2? | (F1 && F2) | F3? | Notes |
|---|---|---|---|---|---|
| P (original) | true (PID) | runs → true | true | skipped (`\|\|` short-circuits) | prints |
| C1 (child of F1) | false (0) | skipped (`&&` short-circuits) | false | runs | C1 and its child C4 both print |
| C2 (child of F2 in P) | true (inherited from P) | false (0) | false | runs | C2 and its child C3 both print |

Printing processes: P, C1, C4, C2, C3 = **5**.

**Why the other options fail:**
- A: Forgets that the child created by F3 also reaches `printf`.
- B: Assumes F3 runs in every process.
- D: Ignores short-circuiting completely (2³).

**Concept recap:**
- `a && b`: skip b if a is false.
- `a || b`: skip b if a is true.
- `fork()` is true in the parent and false (0) in the child.

**Trap alert:** A child inherits the values already computed. C2 inherits F1 = true from P, so it goes on to evaluate F2's result (0, false).

**Interview angle:** Follow-up: *"What about `fork() || fork() && fork()`?"* This is `F1 || (F2 && F3)`.
- P: F1 is true → skips the rest.
- C1: F2 → C1 gets true and runs F3, creating C3. C1's F2 child, C2, gets 0, so `&&` short-circuits.
- Processes: P, C1, C2, C3 → **4** prints.

---

## Q26. [M2 | fork() and separate address spaces | Medium | Code-trace | OA + Interview]

```c
int x = 10;
int main(void) {
    pid_t p = fork();
    if (p == 0) {
        x += 5;
        printf("child %d\n", x);
    } else {
        wait(NULL);
        x -= 5;
        printf("parent %d\n", x);
    }
    return 0;
}
```

Assume `fork` succeeds and output goes to a terminal. What is printed?

- A. `child 15` then `parent 10`
- B. `child 15` and `parent 5`, in either order (nondeterministic)
- C. `child 15` only; the parent blocks forever in `wait()`
- D. `child 15` then `parent 5`, in exactly this order

**Answer:** D

**Why:** `fork` gives the child its **own copy** of the address space. The kernel implements this with **copy-on-write** (COW): the pages are physically shared until one side writes, and then that page is copied. So changing `x` in one process never affects the other. `wait(NULL)` makes the parent sleep until the child terminates, which fixes the order.

**Working:**
1. After `fork`: the parent has x = 10, and the child has its own x = 10.
2. Child: x = 10 + 5 = 15 → prints `child 15` → returns from main → exits.
3. Parent: `wait` blocks until the child has exited → then x = 10 − 5 = **5** → prints `parent 5`.
4. Order: the child's line always comes first, because the parent prints only after `wait` returns.

**Why the other options fail:**
- A: Assumes the variable is shared (15 − 5 = 10). This is the separate-address-space misconception.
- B: The values are right, but it ignores `wait`, which forces the order.
- C: `wait` returns as soon as the child exits. It doesn't block forever.

**Concept recap:** After `fork`, memory is *logically* copied (COW). Writes are private to each process. `wait()` = block until some child terminates, then reap it.

**Trap alert:** Remove the `wait(NULL)` and the order becomes nondeterministic. The values stay 15 and 5.

**Interview angle:** "Does a child see the parent's later changes to a global?" No. Follow-up: *"Then how would they share data?"* Through IPC: shared memory (`mmap` with `MAP_SHARED`), pipes and similar (covered in Batch 3).

---

## Q27. [M2 | exec() | Medium | MCQ | OA + Interview]

Run in a terminal (stdout is line-buffered):

```c
printf("start\n");
execlp("ls", "ls", NULL);
printf("end\n");
```

Which is correct?

- A. "end" is printed only if `execlp` fails; on success the process's program is replaced by `ls`, and its PID stays the same
- B. "end" is printed after `ls` finishes
- C. `execlp` creates a new child process to run `ls`, and the caller continues to print "end"
- D. `execlp` succeeds and returns 0, then "end" is printed

**Answer:** A

**Why:** `exec` **replaces** the calling process's program image (code, data, heap, stack) with a new program. The original code, including the `printf("end")` line, no longer exists, so there is nothing to return to. That is why `exec` returns only on failure, with −1.

**Why the other options fail:**
- B: Assumes `exec` waits for the program and then comes back. That's `system()`, not `exec`.
- C: Confuses `exec` with `fork`. `exec` never creates a process.
- D: A successful `exec` never returns.

**Concept recap:** Kept across `exec`:
- PID and PPID
- open fds, unless marked close-on-exec
- current working directory
- signals set to *ignore*

Reset across `exec`: signals with custom *handlers* go back to their default, because the handler code is gone.

**Trap alert:** If stdout is redirected to a file, "start\n" may still be sitting in the stdio buffer when `exec` wipes user memory, and it is silently **lost**. Call `fflush(stdout)` before `exec`.

**Interview angle:** "Why does Unix split process creation into `fork` + `exec`?" Between the two calls, the child can rearrange its fds (redirection), change directory, or drop privileges. Follow-up: *"Alternatives?"* `posix_spawn` and `vfork`.

---

## Q28. [M2 | wait / waitpid | Medium | MSQ | Interview]

POSIX/Linux. **Select all that apply.**

- A. `wait(&status)` blocks until any one child terminates (returning immediately if one already has) and returns that child's PID.
- B. `waitpid(pid, &status, WNOHANG)` returns 0 immediately if that child exists but has not yet changed state.
- C. `wait()` returns −1 with `errno = ECHILD` if the caller has no children.
- D. A successful `wait()` reaps the child: its zombie entry is removed from the process table.
- E. `waitpid()` can wait only for one specific PID; it cannot wait for "any child".

**Answer:** A, B, C, D

**Why:**
- `wait` = "give me any child that has finished".
- `waitpid` adds control over *which* child, plus non-blocking mode (`WNOHANG`).
- Collecting a child's status is what **reaps** it, meaning the kernel frees its process-table slot.

**Why the other options fail:**
- E: `waitpid(-1, …)` waits for any child. `pid = 0` means any child in the caller's process group, and `pid < −1` means any child in group |pid|.

**Concept recap:** Decode `status` with macros, never directly:
- `WIFEXITED(status)` → did it exit normally?
- `WEXITSTATUS(status)` → the exit code (low 8 bits).
- `WIFSIGNALED` / `WTERMSIG` → was it killed by a signal, and which one?

**Trap alert:** Printing `status` raw after the child calls `exit(1)` typically shows **256**, not 1, because the code is stored in the second byte. Always use `WEXITSTATUS`.

**Interview angle:** "`wait` vs `waitpid`?" Follow-up: *"How does a server reap children without blocking?"* Install a SIGCHLD handler that loops `while (waitpid(-1, NULL, WNOHANG) > 0);`.

---

## Q29. [M2 | What a child inherits on fork | Medium | MSQ | OA]

POSIX/Linux. Which does the child receive from the parent after `fork()`? **Select all that apply.**

- A. A copy of the parent's environment variables and current working directory
- B. The parent's PID
- C. The parent's signal dispositions (installed handlers, ignored signals)
- D. The parent's set of pending signals
- E. A copy of the parent's memory contents (implemented copy-on-write)

**Answer:** A, C, E

**Why:** The child is a near-copy of the parent's *execution environment*. That includes:
- its memory (logically copied),
- its environment variables and current directory,
- its signal handlers (the handler code is in the copied memory, so it still makes sense).

Things that describe the *parent's own identity or pending events* are not copied.

**Why the other options fail:**
- B: The child gets a new, unique PID.
- D: POSIX says the child's pending-signal set starts **empty**, so a signal sent to the parent isn't delivered twice.

**Concept recap:** Not inherited: the PID, pending signals, `fcntl` record locks, pending alarms/timers, and resource-usage counters (reset to 0).

**Trap alert:** Signal *dispositions* are inherited by `fork`, but `exec` resets custom handlers to their default (see Q27). Keep the two calls separate in your head.

**Interview angle:** "What does a child inherit?" Follow-up: *"Do parent and child share file offsets?"* We cover this in M11 (fd tables).

---

## Q30. [M2 | Zombie vs orphan | Medium | MCQ | Interview]

Which pair of definitions is correct?

- A. Zombie: running but has no parent; Orphan: has exited but its parent hasn't called `wait`
- B. Zombie: has exited, but its parent hasn't yet collected its exit status with `wait`; Orphan: still running, but its parent has exited, so it is adopted by init (PID 1) or a subreaper
- C. Both are processes stuck in the Waiting state
- D. Zombie: a process using 100% CPU; Orphan: a process blocked forever on I/O

**Answer:** B

**Why:**
- When a process exits, the kernel frees its memory and files. It keeps a tiny process-table entry (PID, exit status, resource usage) so the parent can still read the result. That leftover entry is a **zombie**, shown as `Z` / `<defunct>` in `ps`.
- If the *parent* dies first, the child becomes an **orphan**. The kernel re-parents it to init, or to a "subreaper" on Linux (`prctl(PR_SET_CHILD_SUBREAPER)`). Init calls `wait` regularly, so orphans don't stay zombies for long.

**Why the other options fail:**
- A: The definitions are swapped. This is the classic mix-up.
- C: A zombie isn't waiting for anything. It's finished.
- D: Invented definitions.

**Concept recap:**
- Zombie = dead child, not yet reaped.
- Orphan = live child, dead parent → re-parented to PID 1 (or a subreaper).

**Trap alert:** A process can be *both* briefly. An orphan that exits is a zombie until its new parent (init) reaps it.

**Interview angle:** "Zombie vs orphan?" Follow-up: *"Are zombies harmful?"* See Q31.

---

## Q31. [M2 | Zombies in practice | Hard | MSQ | OA + Interview]

A Linux server forks one child per request and never calls `wait()`. It runs for days. **Select all that apply.**

- A. Each finished child remains a zombie until the server reaps it or the server itself exits.
- B. Each zombie keeps its full heap and stack memory allocated until it is reaped.
- C. Enough zombies can exhaust the PID/process limit, so later `fork()` calls fail.
- D. On Linux (POSIX.1-2001), explicitly setting SIGCHLD's disposition to `SIG_IGN` makes the kernel reap children automatically, so no zombies form.
- E. Sending SIGKILL to a zombie removes it.

**Answer:** A, C, D

**Why:** A zombie costs very little memory but holds a **PID and a process-table slot**. Thousands of them can hit the limit (`pid_max`, `ulimit -u`), and then `fork` fails with `EAGAIN`.

**Why the other options fail:**
- B: The memory, fds and other resources were freed at exit. Only the small table entry remains.
- E: A zombie is already dead, so there is nothing to kill. Remove it by having the parent reap it, or by killing the *parent* so init adopts and reaps it.

**Concept recap:** Ways to avoid zombies:
1. A SIGCHLD handler that loops `waitpid(-1, NULL, WNOHANG)`.
2. `signal(SIGCHLD, SIG_IGN)`.
3. The double-fork trick: the grandchild is orphaned, adopted by init, and reaped by init.

**Trap alert:** "`kill -9` the zombie" is the wrong instinct.

**Interview angle:** "`ps` shows hundreds of `<defunct>` processes under your server. Diagnose and fix." Follow-up: *"Why loop inside the SIGCHLD handler?"* Standard signals don't queue. Five children exiting together may deliver only one SIGCHLD.

---

## Q32. [M2 | fork() and stdio buffers | Hard | Code-trace | OA + Interview]

Linux/glibc, run in a terminal (stdout is line-buffered):

```c
printf("A");      // no newline
fork();
printf("B\n");
```

Assume `fork` succeeds. Ignoring which process's line appears first, what is printed?

- A. `AB` on one line and `B` on another
- B. `AB` once only
- C. `B` twice; the `A` is lost
- D. `AB` twice, on two lines

**Answer:** D

**Why:** `printf` writes into a **user-space buffer**, and a line-buffered stream is flushed only at a newline. At the moment of `fork`, "A" is still sitting in the buffer. That buffer is ordinary process memory, so the child gets a copy of it.

**Working:**
1. `printf("A")` → buffer = "A". There's no newline, so no `write()` yet.
2. `fork()` → both processes now have buffer = "A".
3. Each process runs `printf("B\n")` → its buffer = "AB\n" → the newline triggers a flush → `write(1, "AB\n", 3)`.
4. Result: two lines of `AB`, from 2 `write` system calls in total.

**Why the other options fail:**
- A: Assumes "A" was written to the terminal before the `fork`.
- B: Forgets that the child prints too.
- C: Assumes the child gets an empty buffer, or that the buffer is discarded.

**Concept recap:** `fork` copies **everything** in user memory, including unflushed stdio buffers. Fix: call `fflush(stdout)` before `fork`, or use `write()` directly.

**Trap alert:** If stdout were redirected to a file or pipe, even `printf("A\n")` before the `fork` would be duplicated, because the stream is then fully buffered (this links to Q18).

**Interview angle:** "Why does my log line appear twice only when output is piped?" Follow-up: *"How do you make the code robust?"* Flush before `fork`, or use `setvbuf` to change the buffering mode.

---

## Q33. [M2 | What threads share | Easy | MCQ | Interview]

What do threads of the same process share?

- A. Their stack
- B. Their register set, including the PC
- C. Code, global data, heap, and open files
- D. Their thread ID

**Answer:** C

**Why:** A thread is an independent **flow of execution** inside a process. Each needs its own PC, registers and stack to follow its own sequence of calls. Everything else belongs to the process and is shared.

**Why the other options fail:**
- A: Each thread has its own stack for its own local variables and return addresses.
- B: Registers and PC are per-thread execution state.
- D: Each thread has a unique ID.

**Concept recap:**
- **Shared:** code, globals, heap, fds, signal *handlers*, address space.
- **Private:** thread ID, PC, registers, stack, signal *mask*, `errno`.

**Trap alert:** `errno` is per-thread in modern libc (thread-local). Otherwise threads would overwrite each other's error codes.

**Interview angle:** "What's shared and what's private between threads?" Follow-up: *"Why must each thread have its own stack?"* Each one has a different chain of active function calls.

---

## Q34. [M2 | Process vs thread trade-offs | Medium | MSQ | OA + Interview]

Linux, pthreads. **Select all that apply.**

- A. If one thread dereferences NULL (SIGSEGV, not handled), the entire process is killed, not just that thread.
- B. A thread can read and modify another thread's local variables if it has a pointer to them, because all stacks are in the same address space.
- C. Creating a thread is generally cheaper than creating a process with `fork()`.
- D. Each thread gets its own private copy of every global variable.
- E. Threads can exchange data through shared variables without kernel IPC such as pipes (though they need synchronization).

**Answer:** A, B, C, E

**Why:** Threads share one address space. That makes communication cheap (E) and creation cheap, since no address space has to be set up (C). It also means there's **no isolation**: a crash or a stray pointer in one thread affects all of them (A, B).

**Why the other options fail:**
- D: Globals are shared. Per-thread copies need explicit thread-local storage (`_Thread_local` / `__thread`).

**Concept recap:**
- Processes: isolation, at the cost of expensive creation and communication.
- Threads: cheap sharing and creation, with no isolation.

**Trap alert:** "Threads are isolated like processes" is false. Thread stacks are just regions of shared memory, with no protection between them.

**Interview angle:** "When would you choose processes over threads?" For fault isolation and security, such as browser tab processes or one PostgreSQL backend process per connection. Choose threads for cheap sharing. Follow-up: *"How does Chrome benefit from multiple processes?"* A crashing tab doesn't take down the browser, and a compromised renderer is sandboxed.

---

## Q35. [M2 | User-level threads and blocking | Medium | MCQ | OA + Interview]

A process uses a pure user-level thread library with a many-to-one model: all ULTs (user-level threads) run on one kernel-schedulable entity. One ULT calls a blocking `read()` on a socket that has no data. What happens to the other ULTs?

- A. They all stop too, because the kernel sees one schedulable entity and blocks it
- B. The kernel schedules another ULT of the process
- C. The ULT library automatically converts the `read` into a new kernel thread
- D. Only the calling ULT blocks; the others continue on other cores

**Answer:** A

**Why:** User-level threads are created and scheduled by a library *inside* the process. The kernel doesn't know they exist. When the single kernel entity blocks in `read`, the library's scheduler can't run either, so every ULT stops.

**Why the other options fail:**
- B: The kernel can't schedule ULTs, because it doesn't know about them.
- C: A pure ULT library doesn't create kernel threads.
- D: Many-to-one has a single kernel entity, so there are no "other cores" to run on.

**Concept recap:** In a many-to-one model, any blocking system call *or page fault* in one ULT blocks the whole process, and there is no multicore parallelism.

**Trap alert:** Page faults count too. A page fault in one ULT blocks every ULT, even without an explicit system call.

**Interview angle:** "Drawbacks of user-level threads?" Follow-up: *"How do green-thread runtimes avoid this?"* They use non-blocking I/O behind the scenes and multiplex the ULTs over several kernel threads (M:N).

---

## Q36. [M2 | ULT vs KLT | Hard | MSQ | Interview]

**Select all that apply.**

- A. Switching between two ULTs of the same process needs no switch into kernel mode.
- B. The kernel gives each ULT its own time slice.
- C. A ULT library can run on a kernel that has no thread support at all.
- D. KLTs (kernel-level threads) of one process can run truly in parallel on different cores.
- E. Linux's NPTL (the standard pthreads library) uses a many-to-one model.

**Answer:** A, C, D

**Why:**
- A ULT switch is just the library saving and loading registers in user space, like a function call, so no trap is needed.
- Because ULTs need nothing from the kernel, the library runs on any OS.
- KLTs are known to the kernel and scheduled individually, so they can be placed on different cores at once.

**Why the other options fail:**
- B: The kernel is unaware of ULTs. It schedules the process (or its kernel thread), not individual ULTs.
- E: NPTL is **one-to-one**. Each pthread is a kernel task created with `clone()` and sharing the address space.

**Concept recap:**
- ULT: fast creation and switching, application-specific scheduling, portable; but blocking stalls all threads and there is no parallelism.
- KLT: true parallelism and independent blocking; but each operation costs a system call.

**Trap alert:** "Kernel threads are always slow" is outdated. Linux made kernel thread creation cheap enough that it chose 1:1.

**Interview angle:** "ULT vs KLT?" Follow-up: *"Why did Linux drop M:N?"* Two cooperating schedulers and messy signal handling cost more than the speed they saved.

---

## Q37. [M2 | Multithreading models | Medium | MCQ | OA]

A runtime maps 8 user threads onto a pool of 3 kernel threads. If one user thread blocks in the kernel, the other user threads keep running on the remaining kernel threads, and up to 3 can run in parallel. Which model is this?

- A. Many-to-one
- B. One-to-one
- C. Many-to-many
- D. One-to-many

**Answer:** C

**Why:** Many ULTs are multiplexed onto a smaller or equal number of KLTs (M:N). This gives parallelism up to the number of KLTs, and one blocked thread doesn't stall the rest.

**Why the other options fail:**
- A: A single KLT means no parallelism, and one blocking call stalls everything.
- B: Would need 8 kernel threads.
- D: Not a standard model. A single instruction stream can't run on several kernel threads at once.

**Concept recap:**
- M:1 = many ULTs on 1 KLT.
- 1:1 = one KLT per ULT.
- M:N = M ULTs on N ≤ M KLTs.
- Two-level = M:N plus the option to *bind* a ULT to its own dedicated KLT.

Example: Go's runtime schedules goroutines M:N over OS threads.

**Trap alert:** Don't decide the model from the number of cores. Decide it from how ULTs map to KLTs.

**Interview angle:** "Explain the threading models with one example each." Follow-up: *"What happens in M:N if all N kernel threads block?"* Runtimes typically create extra kernel threads to keep the pool running.

---

## Q38. [M2 | Benefits of multithreading | Easy | MCQ | Interview]

In a GUI app, the "Open" button's click handler reads a 2 GB file, and the window freezes until the read finishes. What is the best fix using threads?

- A. Raise the process's priority
- B. Call `fork()` for each screen repaint
- C. Use a larger read buffer so the file loads faster
- D. Do the file loading in a background (worker) thread, keeping the UI thread free to handle events, and post the result back to the UI thread

**Answer:** D

**Why:** The UI thread runs the **event loop** (the loop that handles clicks and repaints). While it is stuck in a long `read`, no events get handled, so the window freezes. Moving the slow work to a worker thread keeps the app **responsive**.

**Why the other options fail:**
- A: Priority doesn't help a thread that is blocked in I/O. The event loop still isn't running.
- B: A separate process doesn't share the UI state, and this is heavyweight.
- C: Faster I/O is still blocking. It shortens the freeze but doesn't remove it.

**Concept recap:** The four benefits of multithreading (Silberschatz): **responsiveness, resource sharing, economy** (cheaper than processes), and **scalability** (using multiple cores).

**Trap alert:** Most GUI toolkits allow UI updates only from the UI thread. The worker must post its result back to the UI thread rather than touching widgets directly.

**Interview angle:** "Why use threads at all?" Name the four benefits with one example each. Follow-up: *"Why can't the worker update the UI directly?"* The toolkit's data structures aren't thread-safe.

---

## Q39. [M2 | Concurrency vs parallelism | Medium | MCQ | OA + Interview]

Which of these is **parallelism**, not merely concurrency?

- A. Four threads on a 4-core CPU, each summing one quarter of an array at the same time
- B. Two threads time-sliced on one core
- C. A single-threaded event loop serving 1000 sockets
- D. Coroutines on one thread taking turns at `await` points

**Answer:** A

**Why:**
- **Concurrency** = several tasks *in progress* during overlapping time spans. It is a property of how the program is structured, and it works on one core by interleaving.
- **Parallelism** = several tasks *executing at the same instant*, which needs multiple hardware execution units.

Only option A has instructions from different tasks running simultaneously.

**Why the other options fail:**
- B, C, D: All interleave on a single execution unit. They are concurrent but not parallel.

**Concept recap:** Every parallel program is concurrent, but not every concurrent program is parallel. Rob Pike's summary: concurrency is about *dealing with* many things at once; parallelism is about *doing* many things at once.

**Trap alert:** "Concurrent" ≠ "simultaneous".

**Interview angle:** "Concurrency vs parallelism, with examples." Follow-up: *"Can adding cores make a concurrent program slower?"* Yes, through lock contention, cache-line bouncing, or false sharing (covered in M4).

---

## Q40. [M2 | fork() in a multithreaded process | Hard | MCQ | OA + Interview]

POSIX/Linux. A process has 4 threads. Thread T2 calls `fork()`. How many threads exist in the child?

- A. 4, copies of all threads
- B. 1, a copy of T2 only; the other threads do not exist in the child
- C. 0 until the child calls `exec`
- D. 4, but only T2 is runnable

**Answer:** B

**Why:** POSIX specifies that the child contains a single thread: a copy of the one that called `fork`. The **entire address space** is still copied, including the memory and the mutexes that belonged to the other threads.

**Why the other options fail:**
- A, D: Duplicating every thread could repeat work already in progress, such as writes that are half done. POSIX deliberately doesn't do this.
- C: A process always has at least one thread while it runs. The child keeps running T2's code.

**Concept recap:** Thread count in the child = 1. If another thread held a lock at the moment of `fork` (for example `malloc`'s internal lock), the child's copy of that lock is **locked forever**, because the owner doesn't exist in the child. Calling `malloc` in the child can then deadlock.

**Trap alert:** So in a multithreaded program, the child should call only **async-signal-safe** functions until it calls `exec`. `pthread_atfork()` handlers exist to make the locks safe across `fork`.

**Interview angle:** "What happens when a multithreaded program forks?" Follow-up: *"Why is fork-then-exec still safe?"* `exec` replaces the whole image, stale locks included, before any risky code runs.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 21 | C | 26 | D | 31 | A, C, D | 36 | A, C, D |
| 22 | B | 27 | A | 32 | D | 37 | C |
| 23 | D | 28 | A, B, C, D | 33 | C | 38 | D |
| 24 | A | 29 | A, C, E | 34 | A, B, C, E | 39 | A |
| 25 | C | 30 | B | 35 | A | 40 | B |

---

**Concepts tested so far (to avoid repeats):** Batch 1 list, plus: PCB fields · `fork` return values · ready vs device queues · loop fork counting (2ⁿ − 1) · fork with `&&`/`||` · COW and separate variables with `wait` ordering · `exec` replacing the image and returning only on failure · `wait`/`waitpid`/`WNOHANG`/`ECHILD`/status macros · what a child inherits (environment, current directory, dispositions; not pending signals) · zombie vs orphan definitions · zombie costs and fixes · stdio buffer duplicated by `fork` · what threads share vs keep private · thread isolation trade-offs · ULT blocking · ULT vs KLT properties and NPTL 1:1 · M:N model · responsiveness · concurrency vs parallelism · `fork` copying only the calling thread
*(File-descriptor sharing and file offsets after `fork` are left for M11.)*

**Completed: 40 / 300    Remaining: 260**
