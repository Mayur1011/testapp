# OS MCQ/MSQ Bank — Batch 3 (Q41–Q60)

**What this batch covers:**
- **M2, IPC (Q41–48):** shared memory vs message passing, pipes, named pipes (FIFOs), message queues, signals and signal handlers, sockets and file descriptors, pipe write atomicity.
- **M3, CPU scheduling begins (Q49–60):** long-, medium- and short-term schedulers, the dispatcher, preemptive vs non-preemptive scheduling, CPU-bound vs I/O-bound mix, scheduling criteria, FCFS, SJF, and predicting bursts by exponential averaging.

> In all code questions, assume the standard headers are included and Linux + glibc. In all scheduling numericals, times are in ms; there is no I/O and no context-switch cost unless the question says otherwise.

---

## Q41. [M2 | Shared memory vs message passing | Easy | MCQ | Interview]

Two processes on one machine exchange data. Which is correct?

- A. With message passing, processes read and write a common memory region directly, without the kernel
- B. With shared memory, once the region is set up, reads and writes are ordinary memory accesses with no system call, but the processes themselves must synchronize access
- C. Shared memory needs a system call for every byte exchanged
- D. Message passing cannot be used between processes on the same machine

**Answer:** B

**Why:**
- **Shared memory:** the kernel is involved only at setup (`shm_open` + `mmap`, or `shmget` + `shmat`). After that, the processes exchange data with plain loads and stores, which is the fastest IPC. Nothing synchronizes them, though, so they must use semaphores or mutexes to avoid race conditions.
- **Message passing:** every send and receive is a system call, and the kernel copies the data. In exchange, it gives synchronization for free, because a receive blocks until a message is there.

**Why the other options fail:**
- A: Describes shared memory, not message passing.
- C: Backwards. Shared memory avoids per-access system calls.
- D: Pipes, message queues and Unix domain sockets are all message passing on one machine.

**Concept recap:**
- Shared memory: fastest; you synchronize it yourself; works on one machine only.
- Message passing: simpler and safer; costs a system call and a kernel copy per message; can work across machines (sockets).

**Trap alert:** "Shared memory is automatically synchronized" is false. Two processes incrementing a shared counter will race exactly like two threads (covered in M4).

**Interview angle:** "Shared memory vs message passing?" Follow-up: *"Workers must share a 4 GB read-mostly matrix. Which do you pick?"* Shared memory. Copying 4 GB through the kernel would be wasteful.

---

## Q42. [M2 | Pipes | Medium | Code-trace | OA + Interview]

```c
int main(void) {
    int fd[2];
    pipe(fd);
    if (fork() == 0) {            // child
        close(fd[0]);
        write(fd[1], "OS", 2);
        exit(0);
    }
    close(fd[1]);                 // parent
    char buf[8];
    int n = read(fd[0], buf, sizeof buf);
    buf[n] = '\0';
    printf("%d %s\n", n, buf);
    return 0;
}
```

Assume all calls succeed, and either process may be scheduled first. What does the parent print?

- A. `0` followed by nothing, if the parent's `read` runs before the child's `write`
- B. `8 OS`
- C. Nothing: the program can deadlock if the parent reaches `read()` first
- D. `2 OS`

**Answer:** D

**Why:** `read` on an empty pipe **blocks** as long as some process still holds the pipe's write end open. It returns as soon as *any* data arrives, giving back however many bytes are available, up to the requested count. It returns 0 (end-of-file) only when the pipe is empty **and** every write end is closed.

**Working:**
1. `pipe(fd)` → `fd[0]` = read end, `fd[1]` = write end. These are the two lowest unused fd numbers, usually 3 and 4.
2. **Case 1: the parent runs first.** The pipe is empty and the child still holds `fd[1]`, so `read` blocks. The child writes 2 bytes, the parent wakes, and `read` returns n = 2.
3. **Case 2: the child runs first.** 2 bytes are already waiting in the pipe, so `read` returns n = 2 immediately.
4. Either way, `buf` = "OS" and the output is `2 OS`.

**Why the other options fail:**
- A: Confuses "empty" with "end-of-file". The child's write end is still open, so `read` waits instead of returning 0.
- B: `read` returns the number of bytes *actually read*, not the buffer size.
- C: A blocked `read` is woken by the child's `write`. Nothing is waiting in a cycle, so there's no deadlock.

**Concept recap:** A pipe is a **one-way, kernel-buffered byte stream**. Mnemonic: `fd[0]` reads (like stdin, fd 0) and `fd[1]` writes (like stdout, fd 1).
- `read` on an empty pipe with writers still open → blocks.
- `read` on an empty pipe with no writers left → returns 0.

**Trap alert:** Swapping the two ends is the classic bug. Writing to `fd[0]` fails with `EBADF`. Why the parent must close its own `fd[1]` (or `read` may never see EOF) is covered in M11.

**Interview angle:** "How does a shell implement `ls | wc`?" With `pipe`, `fork`, `dup2` and `exec` (covered in M11). Follow-up: *"Is a pipe bidirectional?"* No. Use two pipes, or a `socketpair()`.

---

## Q43. [M2 | Named pipes (FIFOs) | Medium | MCQ | OA]

Two **unrelated** programs, started separately from different terminals on the same Linux machine, must stream data through a pipe-like channel. Which mechanism fits directly?

- A. A named pipe (FIFO) created with `mkfifo`; it has a name in the file system that both programs can `open`
- B. An anonymous pipe created with `pipe()` by the first program
- C. A pipe created by the second program and handed over by its PID
- D. Shared memory, because pipes of any kind require a common parent process

**Answer:** A

**Why:** An anonymous pipe exists only as two fds inside one process. Only processes that **inherit** those fds through `fork`, i.e. related processes, can use it. A **FIFO** is a special file with a path. Any process with permission can `open` it. The data still flows through a kernel buffer and is never stored on disk.

**Why the other options fail:**
- B: An unrelated process has no way to get those fds. (Advanced exception: fds can be passed over a Unix domain socket using `SCM_RIGHTS`.)
- C: There is no "hand over a pipe by PID" mechanism.
- D: The premise is false. FIFOs exist precisely so unrelated processes can use a pipe.

**Concept recap:**
- Anonymous pipe: related processes only, no name, gone when all its fds are closed.
- FIFO: any process with permission, has a path, and the path stays until `unlink`ed. Its data is never written to disk.

**Trap alert:** "FIFO data is written to the disk file" is false. Only the directory entry is on disk.

**Interview angle:** "Pipe vs named pipe?" Follow-up: *"What happens if you open a FIFO for reading and no writer exists?"* By default the `open` blocks until a writer opens the other end, unless you pass `O_NONBLOCK`.

---

## Q44. [M2 | Message queues | Medium | MSQ | Interview]

POSIX or System V message queues on Linux. **Select all that apply.**

- A. They preserve message boundaries: each receive returns exactly one whole message.
- B. Messages carry a type (System V) or a priority (POSIX), so a receiver can take a particular kind of message first.
- C. A queue lives in the kernel until it is explicitly removed, even after the process that created it exits.
- D. A pipe also preserves boundaries: two separate 10-byte writes are always returned by two separate reads.
- E. Only processes related by `fork` can use the same message queue.

**Answer:** A, B, C

**Why:** A message queue is **record-oriented**: the kernel stores whole messages as separate items, not a stream of bytes. Queues are named globally (a numeric key in System V, a name like `/myq` in POSIX), so any permitted process can open one. They have **kernel persistence**: a queue survives until someone removes it or the system reboots.

**Why the other options fail:**
- D: A pipe is a byte stream. The reader might get all 20 bytes in one `read`, or 5 bytes and then 15.
- E: Queues are found by key or name, so the processes don't need to be related.

**Concept recap:**
- Pipe = byte stream, FIFO order, disappears when closed.
- Message queue = discrete messages with a type or priority, persists in the kernel until removed.

**Trap alert:** Forgotten queues **leak** kernel resources. Clean them up with `ipcrm` or `mq_unlink`. You can inspect them with `ipcs`, or under `/dev/mqueue` for POSIX queues.

**Interview angle:** "When would you pick a message queue over a pipe?" When you need message boundaries, priorities, or unrelated processes. Follow-up: *"How do you send messages over a byte stream like TCP or a pipe?"* Add framing: a length prefix or a delimiter.

---

## Q45. [M2 | Signals | Medium | MCQ | OA + Interview]

Which two signals can never be caught, blocked, or ignored by a process?

- A. SIGINT and SIGTERM
- B. SIGSEGV and SIGFPE
- C. SIGKILL and SIGSTOP
- D. SIGCHLD and SIGALRM

**Answer:** C

**Why:** A **signal** is an asynchronous notification delivered to a process, a kind of software interrupt. For each signal a process can choose its **disposition** (how it reacts): the default action, ignore it, or run a handler. The OS must always keep one guaranteed way to **kill** a process and one to **stop** it, or a buggy or malicious program could make itself unkillable. So SIGKILL and SIGSTOP alone ignore whatever the process chose.

**Why the other options fail:**
- A: SIGINT (Ctrl+C) and SIGTERM (what plain `kill` sends) can be caught. This is how programs shut down gracefully.
- B: These can be caught, but returning from a SIGSEGV handler after a real fault re-runs the faulting instruction, which faults again. Handlers usually just log and exit.
- D: Both can be caught. SIGCHLD's default is to ignore, and SIGALRM's default is to terminate.

**Concept recap:**
- `kill -9 pid` → SIGKILL (9): cannot be handled.
- `kill pid` → SIGTERM (15): can be caught to clean up.
- SIGSTOP pauses a process and SIGCONT resumes it. (Ctrl+Z sends SIGTSTP, which *can* be caught.)

**Trap alert:** Using `kill -9` as the first resort skips all cleanup: buffers aren't flushed and temp files and locks are left behind. Send SIGTERM first.

**Interview angle:** "SIGTERM vs SIGKILL?" Follow-up: *"Why might `kill -9` not remove a process immediately?"* The process may be in uninterruptible sleep (state `D`) waiting on I/O, or it may already be a zombie.

---

## Q46. [M2 | Signal handlers | Hard | MSQ | Interview]

Linux/POSIX. **Select all that apply.**

- A. `printf()` is not async-signal-safe but `write()` is, so a handler should use `write()` for output.
- B. Calling `malloc()` inside a signal handler is safe, because handlers run in kernel mode.
- C. A running handler can be interrupted by a different signal unless that signal is blocked (e.g., via `sa_mask`).
- D. SIGTERM's default action is to terminate the process, and it can be caught to perform a graceful cleanup.
- E. After a handler returns, an interrupted blocking system call always fails with `EINTR`; it can never be restarted automatically.

**Answer:** A, C, D

**Why:** A handler runs **asynchronously, in user mode**, at whatever point the program happened to be. That could be in the middle of `malloc` or `printf` while they hold an internal lock. If the handler then calls the same function, it can deadlock or corrupt that function's data. **Async-signal-safe** functions are the ones POSIX guarantees are safe to call from a handler, for example `write`, `_exit` and `sigaction`.

**Why the other options fail:**
- B: Handlers run in **user** mode, and `malloc` is not async-signal-safe.
- E: With the `SA_RESTART` flag (set via `sigaction`), many interrupted system calls are restarted automatically.

**Concept recap:** The safe handler pattern: the handler only sets a flag declared as `volatile sig_atomic_t`, and the main loop checks that flag and does the real work.

**Trap alert:** `printf` in a handler "works fine in testing" and then deadlocks once in a million runs. Bugs that depend on rare timing never show up in a quick test.

**Interview angle:** "What can you safely do in a signal handler?" Follow-up: *"How do you turn signals into normal events in an event loop?"* With the self-pipe trick (the handler writes one byte to a pipe) or Linux `signalfd`. Either way, the signal shows up as a readable fd.

---

## Q47. [M2 | Sockets and file descriptors | Easy | MCQ | OA]

Which statement about sockets on Linux is correct?

- A. Sockets only work between processes on the same machine
- B. A socket is represented in the process by a file descriptor, so `read()`/`write()` (as well as `send()`/`recv()`) work on it, and with TCP/IP it can connect processes on different machines
- C. Sockets require the two processes to set up shared memory first
- D. A TCP socket preserves message boundaries, so each `recv()` returns exactly one `send()`'s data

**Answer:** B

**Why:** A **file descriptor** (fd) is a small non-negative integer that indexes the process's table of open I/O objects. Unix gives files, pipes, terminals and sockets the same fd interface. A **socket** is a communication endpoint with an address: an IP address and port for TCP/UDP, or a file path for a Unix domain socket.

**Why the other options fail:**
- A: Unix domain sockets are local, but TCP and UDP sockets work across machines.
- C: Sockets don't involve shared memory.
- D: TCP is a **byte stream**, so boundaries are lost (like a pipe). UDP datagrams do preserve boundaries.

**Concept recap:** A new fd always gets the **lowest unused number**. Since 0, 1 and 2 are stdin, stdout and stderr, the first socket or file a program opens is usually fd 3.

**Trap alert:** Assuming `recv` on TCP returns one whole message is a classic network-programming bug. Loop until you have read the length you expect.

**Interview angle:** "Which IPC works across machines?" Sockets. Follow-up: *"Unix domain socket vs TCP over loopback (127.0.0.1)?"* Unix sockets skip the TCP/IP stack and can pass fds and credentials between processes.

---

## Q48. [M2 | Pipe write atomicity | Hard | MSQ | OA + Interview]

Several processes write log records into the same pipe, and one reader reads them. Each record is written with a single `write()` call. Linux. **Select all that apply.**

- A. A single `write()` of at most `PIPE_BUF` bytes is atomic: its bytes are never interleaved with other writers' data.
- B. On Linux, `PIPE_BUF` is 4096 bytes.
- C. A `write()` larger than `PIPE_BUF` may have its bytes interleaved with other writers' data.
- D. Each `read()` by the reader always returns exactly one writer's record.
- E. A pipe's total capacity equals `PIPE_BUF`, so a 5000-byte write always fails.

**Answer:** A, B, C

**Why:** POSIX guarantees that a write of at most `PIPE_BUF` bytes lands in the pipe as one contiguous block. POSIX requires `PIPE_BUF` ≥ 512, and Linux uses 4096. So 200-byte records arrive intact, though their order across different writers is arbitrary. Larger writes can be split and interleaved with other writers' data.

**Why the other options fail:**
- D: The pipe is still a byte stream. A single `read` may return several records joined together, or part of one. The reader must split records itself, for example on newlines.
- E: Capacity is separate from the atomicity limit. On Linux the default pipe capacity is 65,536 bytes (16 pages), adjustable with `fcntl(F_SETPIPE_SZ)`. In blocking mode, a large write **waits** for space; it doesn't fail.

**Concept recap:** A record of at most `PIPE_BUF` bytes is never split. Pipe capacity ≠ `PIPE_BUF`.

**Trap alert:** Mixing up `PIPE_BUF` (the atomicity limit) with the pipe's capacity (how much it can hold) is the classic confusion.

**Interview angle:** "Is it safe for many processes to write into one pipe?" Yes, if each record is at most `PIPE_BUF` bytes and sent in one `write`. Follow-up: *"And to one log file?"* Open it with `O_APPEND`. On local file systems, moving to the end and writing then happen as one step, so appends don't overwrite each other. This isn't guaranteed on NFS.

---

## Q49. [M3 | Long-, medium-, short-term schedulers | Easy | MCQ | OA]

Which scheduler controls the **degree of multiprogramming**, i.e. how many processes are admitted into memory?

- A. The short-term (CPU) scheduler
- B. The dispatcher
- C. The I/O (disk) scheduler
- D. The long-term (job) scheduler

**Answer:** D

**Why:** The **long-term scheduler** chooses which jobs from the pool on disk are admitted into memory and the ready queue. That choice sets how many processes compete at once. It runs rarely (seconds to minutes). The **short-term scheduler** picks, among processes already admitted, which one runs next on the CPU, and runs every few milliseconds.

**Why the other options fail:**
- A: The short-term scheduler chooses only among processes already in memory.
- B: The dispatcher performs the CPU switch. It makes no admission decisions.
- C: The disk scheduler orders disk requests (covered in M9).

**Concept recap:**

| Scheduler | Decides | Runs |
|---|---|---|
| Long-term | admission into memory | rarely |
| Medium-term | swapping processes out and back in | occasionally |
| Short-term | which ready process gets the CPU | every few ms |

**Trap alert:** Time-sharing systems like Linux and Windows have almost no long-term scheduler. Every program you launch is admitted, limited only by available resources.

**Interview angle:** "Name the three schedulers and how often each runs." Follow-up: *"Which one barely exists on a Linux desktop?"* The long-term scheduler.

---

## Q50. [M3 | Medium-term scheduler | Medium | MCQ | OA + Interview]

Memory is overcommitted and the system is paging heavily. The OS picks a few processes, moves their whole state out to disk, and brings them back later when memory frees up. Which scheduler does this, and what is the action called?

- A. Medium-term scheduler; swapping (suspending a process out to disk and later back in)
- B. Long-term scheduler; admission
- C. Short-term scheduler; preemption
- D. Dispatcher; context switching

**Answer:** A

**Why:** The **medium-term scheduler** temporarily *reduces* the degree of multiprogramming by swapping whole processes out. A swapped-out process enters the Ready-Suspended or Blocked-Suspended state (the extra states mentioned in Q20). Freeing memory this way lets the remaining processes run without constant page faults.

**Why the other options fail:**
- B: Admission brings *new* jobs in. It can't remove processes that are already running.
- C: A preempted process stays in memory, in the ready queue.
- D: The dispatcher only hands over the CPU.

**Concept recap:** Swapping = moving a whole process between memory and disk. It is used to relieve memory pressure or to rebalance the CPU-bound/I/O-bound mix.

**Trap alert:** Don't confuse swapping a **whole process** with **paging** individual pages in and out (covered in M7).

**Interview angle:** "What does the medium-term scheduler do?" Follow-up: *"Swapping vs paging?"* Swapping moves whole processes. Paging moves individual pages on demand.

---

## Q51. [M3 | Dispatcher and dispatch latency | Medium | MSQ | OA]

Which of these are part of the **dispatcher's** job? **Select all that apply.**

- A. Switching context: loading the saved registers of the chosen process
- B. Switching the CPU to user mode
- C. Jumping to the right location in the user program to resume it
- D. Deciding which ready process should run next
- E. Admitting new jobs from disk into memory

**Answer:** A, B, C

**Why:** The work splits into **policy** and **mechanism**. The short-term scheduler is the policy: it decides *who* runs. The dispatcher is the mechanism: it actually *hands over* the CPU. **Dispatch latency** is the time the dispatcher takes to stop one process and start another. It is pure overhead, so it should be small.

**Why the other options fail:**
- D: That decision belongs to the short-term scheduler.
- E: That is the long-term scheduler's job.

**Concept recap:** Dispatcher = context switch + switch to user mode + jump to the saved PC. Dispatch latency is the "context-switch time δ" that some Gantt-chart questions include.

**Trap alert:** Many textbooks loosely call both "the scheduler", but exam questions separate policy (who runs) from mechanism (how the switch happens).

**Interview angle:** "Scheduler vs dispatcher?" Follow-up: *"Why does dispatch latency matter when choosing a Round-Robin quantum?"* If the quantum is close to the switch cost, most CPU time is wasted on switching (covered in the next batch).

---

## Q52. [M3 | Preemptive vs non-preemptive | Medium | MCQ | OA]

A scheduling decision may be needed when a process:
1. switches Running → Waiting,
2. switches Running → Ready (e.g., a timer interrupt),
3. switches Waiting → Ready,
4. terminates.

Under a purely **non-preemptive** scheduler, the scheduler is invoked only in cases:

- A. 2 and 3
- B. 1, 2, 3 and 4
- C. 1 and 4
- D. 3 and 4

**Answer:** C

**Why:** In **non-preemptive** scheduling (also called cooperative), a process keeps the CPU until it gives it up voluntarily, by blocking (case 1) or exiting (case 4). In those two cases the CPU is free, so the scheduler *must* choose someone. In cases 2 and 3, choosing would mean taking the CPU away from a running process, which only a **preemptive** scheduler does.

**Why the other options fail:**
- A: These are exactly the cases that only preemptive schedulers act on.
- B: Acting in all four cases means the scheduler is preemptive.
- D: Case 3 (a process waking up) triggers rescheduling only under preemption.

**Concept recap:**
- Non-preemptive: FCFS, SJF, non-preemptive priority.
- Preemptive: Round Robin, SRTF, preemptive priority.

**Trap alert:** Interrupts still happen under a non-preemptive scheduler. Handlers run, but they don't give the CPU to another process. Also note that preemption is what creates race conditions on shared data, including inside the kernel (covered in M4).

**Interview angle:** "Preemptive vs non-preemptive?" Follow-up: *"What problems does preemption bring?"* Race conditions on shared data, and the need for locks inside a preemptible kernel.

---

## Q53. [M3 | CPU-bound vs I/O-bound mix | Easy | MCQ | Interview]

The long-term scheduler admits **only I/O-bound** processes (short CPU bursts, frequent I/O). What happens?

- A. The CPU is fully busy and the I/O devices sit idle
- B. Throughput is maximized, because I/O-bound jobs finish fast
- C. Nothing; the mix of processes doesn't matter
- D. The ready queue is often empty, so the CPU sits idle while the device queues grow long

**Answer:** D

**Why:** I/O-bound processes spend most of their time waiting for devices, so few of them are ready to run at any moment. With only those processes, the CPU often has nothing to do. With only CPU-bound processes, the opposite happens: the ready queue is long and the devices sit idle. A **balanced mix** keeps both the CPU and the devices busy.

**Why the other options fail:**
- A: Describes the all-CPU-bound case.
- B: An idle CPU wastes capacity, so throughput drops.
- C: The mix directly determines how busy the CPU and devices are.

**Concept recap:**
- CPU-bound = a few long CPU bursts.
- I/O-bound = many short CPU bursts.

Measured burst distributions show most CPU bursts are short.

**Trap alert:** "I/O-bound" doesn't mean "slow program". It describes *where its time goes*.

**Interview angle:** "Why do schedulers favour I/O-bound processes?" Giving them the CPU quickly lets them issue their next I/O right away, which keeps devices busy and interactive programs snappy. MLFQ does this automatically (covered in the next batch).

---

## Q54. [M3 | Scheduling criteria | Medium | MSQ | OA + Interview]

For a process that does **no I/O**, **select all correct definitions.**

- A. Turnaround time = completion time − arrival time
- B. Waiting time = turnaround time − burst time
- C. Response time = time from arrival until the process first gets the CPU
- D. Throughput = time taken to complete one process
- E. Waiting time includes time spent waiting in device (I/O) queues

**Answer:** A, B, C

**Why:**
- **Turnaround** measures the whole stay in the system, from arrival to completion.
- **Waiting** is the part of that stay spent in the ready queue, so it is turnaround minus the time actually running.
- **Response** is how long until the process first runs, which is what an interactive user feels.

**Why the other options fail:**
- D: Throughput is a *rate*: processes completed per unit of time.
- E: Waiting time counts only time in the **ready queue**. Time doing I/O is not waiting time.

**Concept recap:** The five criteria: **maximize** CPU utilization and throughput; **minimize** turnaround, waiting and response time.
- TAT = CT − AT
- WT = TAT − BT (no I/O)
- RT = (first time on CPU) − AT

(CT = completion time, AT = arrival time, BT = burst time.)

**Trap alert:** For **non-preemptive** algorithms, RT = WT, because a process runs only once, straight through to completion. Under preemption, RT ≤ WT.

**Interview angle:** "Which criterion matters most for interactive systems?" Response time. For batch systems it's throughput and turnaround. Follow-up: *"Why might you minimize the variance of response time rather than its average?"* Users prefer predictable delays to delays that are usually short but sometimes very long.

---

## Q55. [M3 | FCFS and order dependence | Easy | MCQ | OA]

Three jobs arrive at time 0 with CPU bursts of 10, 4 and 2 ms, and are scheduled FCFS. Which arrival order gives the lowest average waiting time?

- A. 10, 4, 2
- B. 2, 4, 10
- C. 4, 10, 2
- D. Every order gives the same average waiting time

**Answer:** B

**Why:** In FCFS, each job waits for the sum of the bursts queued ahead of it. Putting short jobs first keeps those sums small:
- Order 2, 4, 10 → waits 0, 2, 6 → average = 8/3 ≈ **2.67 ms**
- Order 4, 10, 2 → waits 0, 4, 14 → average = 18/3 = 6 ms
- Order 10, 4, 2 → waits 0, 10, 14 → average = 24/3 = 8 ms

**Why the other options fail:**
- A: The longest job first gives the worst average.
- C: Places the 10 ms job before the 2 ms job.
- D: The *total work* is the same in every order (16 ms), but the waits are not.

**Concept recap:** Under FCFS, WTᵢ = sum of the bursts of jobs ahead of i. Shortest-first minimizes the total, and that is exactly the idea behind SJF.

**Trap alert:** "FCFS is fair, so the order doesn't matter" is wrong. FCFS performance depends heavily on arrival order. One long job at the front makes everyone wait (the convoy effect, covered in the next batch).

**Interview angle:** "Why is FCFS poor for interactive systems?" Follow-up: *"What is the convoy effect?"*

---

## Q56. [M3 | FCFS with arrival times and idle CPU | Medium | Numerical | OA]

FCFS, non-preemptive.

| Process | Arrival | Burst |
|---|---|---|
| P1 | 0 | 3 |
| P2 | 5 | 2 |
| P3 | 6 | 4 |

Find the **average turnaround time** and the **CPU utilization** over the interval from 0 to the completion of the last process.

- A. 7.00 ms and 81.8%
- B. 3.33 ms and 100%
- C. 3.33 ms and 81.8%
- D. 0.33 ms and 81.8%

**Answer:** C

**Why:** A process cannot start before it arrives. If the CPU finishes early and nothing has arrived yet, it sits **idle**, and that idle time counts against utilization.

**Working:**
- **Formulas:** TAT = CT − AT; WT = TAT − BT; Utilization = Σ bursts ÷ (time of last completion − 0).
- **Gantt chart:**

```
| P1 | idle | P2 | P3   |
0    3      5    7      11
```

| Process | AT | BT | CT | TAT = CT − AT | WT = TAT − BT |
|---|---|---|---|---|---|
| P1 | 0 | 3 | 3 | 3 | 0 |
| P2 | 5 | 2 | 7 | 2 | 0 |
| P3 | 6 | 4 | 11 | 5 | 1 |
| **Sum** | | 9 | | 10 | 1 |

- Average TAT = 10 / 3 = **3.33 ms**
- Utilization = busy time / total time = 9 / 11 = **81.8%**

**Why the other options fail:**
- A: Averages the *completion* times (21/3 = 7) and forgets to subtract arrival times.
- B: Assumes the CPU never idles, which would mean starting P2 before it arrives.
- D: 0.33 is the average *waiting* time, not turnaround.

**Concept recap:** TAT = CT − AT. Utilization = Σ BT / makespan, where makespan = time from start until the last process finishes. Idle gaps count in the denominator.

**Trap alert:** Never "slide" a process earlier to fill a gap. Always check whether it has arrived yet.

**Interview angle:** Rarely asked verbally. Follow-up: *"When is FCFS actually a good choice?"* For batch jobs of similar length, where its simplicity and lack of starvation pay off.

---

## Q57. [M3 | Non-preemptive SJF with arrivals and ties | Hard | Numerical | OA + Interview]

Non-preemptive SJF. The scheduler decides only when the CPU becomes free. Ties in burst length are broken by earlier arrival.

| Process | Arrival | Burst |
|---|---|---|
| P1 | 0 | 6 |
| P2 | 1 | 8 |
| P3 | 2 | 3 |
| P4 | 3 | 3 |

Find the **average waiting time** and the **average turnaround time**.

- A. 7.75 ms and 12.75 ms
- B. 5.25 ms and 10.25 ms
- C. 6.75 ms and 11.75 ms
- D. 4.75 ms and 9.75 ms

**Answer:** B

**Why:** At each decision point, SJF picks the shortest burst **among the processes that have already arrived**. Because it is non-preemptive, it never interrupts the running process.

**Working:**
1. t = 0: only P1 has arrived → run P1 from 0 to 6.
2. t = 6: ready are P2 (8), P3 (3), P4 (3). The shortest is 3, a tie between P3 and P4. P3 arrived earlier, so run P3 from 6 to 9.
3. t = 9: ready are P2 (8) and P4 (3) → run P4 from 9 to 12.
4. t = 12: only P2 is left → run P2 from 12 to 20.

```
| P1     | P3  | P4  | P2       |
0        6     9     12         20
```

| Process | AT | BT | CT | TAT | WT |
|---|---|---|---|---|---|
| P1 | 0 | 6 | 6 | 6 | 0 |
| P2 | 1 | 8 | 20 | 19 | 11 |
| P3 | 2 | 3 | 9 | 7 | 4 |
| P4 | 3 | 3 | 12 | 9 | 6 |
| **Sum** | | | | 41 | 21 |

- Average WT = 21 / 4 = **5.25 ms**
- Average TAT = 41 / 4 = **10.25 ms**

**Why the other options fail:**
- A: The FCFS result (P1, P2, P3, P4 in arrival order): WT = 0, 5, 12, 14.
- C: Uses the start time as the waiting time and the completion time as the turnaround time, forgetting to subtract arrival.
- D: The **preemptive** SRTF result. Its Gantt chart is P1 0–2 | P3 2–5 | P4 5–8 | P1 8–12 | P2 12–20. Catches mixing up SJF with SRTF.

**Concept recap:** Non-preemptive SJF: when the CPU becomes free, choose min(burst) among *arrived* processes. Then run that process to completion.

**Trap alert:** Don't pick P3 at t = 0. It hasn't arrived yet. And don't interrupt P1 when shorter jobs arrive; that would be SRTF.

**Interview angle:** "What information does SJF need?" The length of each process's next CPU burst. Follow-up: *"What changes if it can preempt?"* It becomes SRTF (covered in the next batch).

---

## Q58. [M3 | SJF in practice | Medium | MCQ | Interview]

SJF is optimal for average waiting time. Why can't a real OS's short-term scheduler implement it exactly?

- A. SJF requires special hardware support
- B. SJF works only on multiprocessors
- C. SJF suffers from Belady's anomaly
- D. The length of a process's next CPU burst isn't known in advance; it must be predicted, typically by exponential averaging of past bursts

**Answer:** D

**Why:** The scheduler can't see the future. It can only predict the next burst from the bursts it has measured so far. The standard predictor is the **exponential average**:

**τₙ₊₁ = α · tₙ + (1 − α) · τₙ**

where tₙ is the measured length of the latest burst, τₙ is the previous prediction, and 0 ≤ α ≤ 1.

**Why the other options fail:**
- A, B: Not true. SJF is purely a policy decision.
- C: Belady's anomaly is about page replacement (covered in M7), not CPU scheduling.

**Concept recap:**
- α = 0: ignore new measurements; the prediction never changes.
- α = 1: use only the last burst.
- 0 < α < 1: older bursts count for less and less, by a factor of (1 − α) each step.

**Trap alert:** In short-term scheduling, SJF uses the **next CPU burst**, not the job's total running time. Batch systems sometimes used a time limit the user supplied instead.

**Interview angle:** "How would you approximate SJF in a real kernel?" Follow-up: *"What does α control?"* The trade-off between reacting quickly to recent behaviour and staying stable.

---

## Q59. [M3 | Exponential averaging | Hard | Numerical | OA]

Predict the next CPU burst with **τₙ₊₁ = α · tₙ + (1 − α) · τₙ**, using α = 0.25 and τ₀ = 12 ms. The measured bursts are t₀ = 8, t₁ = 16 and t₂ = 4 ms. Find τ₃.

- A. 10.19 ms
- B. 9.33 ms
- C. 4.00 ms
- D. 6.56 ms

**Answer:** A

**Why:** Each new prediction blends the burst just measured (weight α) with the previous prediction (weight 1 − α). You must apply the formula step by step.

**Working:**
- **Given:** α = 0.25, 1 − α = 0.75, τ₀ = 12.
- τ₁ = 0.25 · 8 + 0.75 · 12 = 2 + 9 = **11**
- τ₂ = 0.25 · 16 + 0.75 · 11 = 4 + 8.25 = **12.25**
- τ₃ = 0.25 · 4 + 0.75 · 12.25 = 1 + 9.1875 = **10.1875 ≈ 10.19 ms**
- **Check with the expanded form:** τ₃ = α·t₂ + (1−α)·α·t₁ + (1−α)²·α·t₀ + (1−α)³·τ₀ = 1 + 3 + 1.125 + 5.0625 = 10.1875 ✓

**Why the other options fail:**
- B: The simple average of the bursts (28/3). That ignores the formula.
- C: Equivalent to α = 1, i.e. using only the last burst.
- D: Swaps α and 1 − α (puts 0.75 on the new burst): τ₁ = 9, τ₂ = 14.25, τ₃ = 6.5625.

**Concept recap:** τₙ₊₁ = α·tₙ + (1 − α)·τₙ, which expands to Σ α(1 − α)ʲ·tₙ₋ⱼ + (1 − α)ⁿ⁺¹·τ₀.

**Trap alert:** Two common slips: putting α on the wrong term, and off-by-one indexing. τ₃ is the prediction for the burst *after* t₂.

**Interview angle:** Usually asked in words, without numbers: "What does α = 0.5 mean?" Equal weight on the latest burst and on all the history before it.

---

## Q60. [M3 | Properties of SJF | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. When all jobs are available at the same time, non-preemptive SJF gives the minimum average waiting time of any non-preemptive schedule.
- B. SJF can starve a long job if shorter jobs keep arriving.
- C. For the same set of jobs, a schedule that minimizes average waiting time also minimizes average turnaround time.
- D. SJF minimizes the maximum waiting time experienced by any single job.
- E. SJF is a special case of priority scheduling in which priority is set by the (predicted) next CPU burst: shorter burst, higher priority.

**Answer:** A, B, C, E

**Why:**
- **A (the exchange argument):** Suppose a longer job a runs immediately before a shorter job b. Swapping them changes the total waiting time by b − a < 0, so the total goes down. Repeating such swaps reaches shortest-first order, so SJF is optimal.
- **C:** TATᵢ = WTᵢ + BTᵢ, and Σ BT is fixed, so minimizing Σ WT also minimizes Σ TAT.
- **B:** A long job keeps being passed over as long as shorter ones keep arriving.
- **E:** SJF is priority scheduling with priority = 1/burst.

**Why the other options fail:**
- D: The long job waits behind *every* shorter job, so the worst-case wait can be very large, and unbounded if short jobs keep arriving. "Optimal on average" says nothing about the worst case.

**Concept recap:** SJF is optimal for the **average** wait, gives no **worst-case** guarantee, and starvation is possible. The usual fix is **aging**: a process's priority rises the longer it waits (covered in the next batch).

**Trap alert:** "Optimal" here means the *average* waiting time. It is neither fair nor optimal for the worst case.

**Interview angle:** "Prove SJF minimizes average waiting time." Give the exchange argument in two lines. Follow-up: *"How do you prevent starvation?"* Aging.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 41 | B | 46 | A, C, D | 51 | A, B, C | 56 | C |
| 42 | D | 47 | B | 52 | C | 57 | B |
| 43 | A | 48 | A, B, C | 53 | D | 58 | D |
| 44 | A, B, C | 49 | D | 54 | A, B, C | 59 | A |
| 45 | C | 50 | A | 55 | B | 60 | A, B, C, E |

---

**Concepts tested so far (to avoid repeats):** Batches 1–2, plus: shared memory vs message passing · pipe read semantics (block vs EOF, bytes returned) · FIFO for unrelated processes · message queues (boundaries, type/priority, kernel persistence) · uncatchable SIGKILL/SIGSTOP · async-signal-safety, `sa_mask`, `SA_RESTART` · sockets as fds, TCP as a byte stream, lowest-fd rule · `PIPE_BUF` atomicity vs pipe capacity · long-term scheduler and degree of multiprogramming · medium-term swapping · dispatcher's duties and dispatch latency · non-preemptive decision points · CPU-bound/I/O-bound mix · TAT/WT/RT/throughput definitions · FCFS order dependence · FCFS with idle CPU and utilization · non-preemptive SJF with ties · burst prediction by exponential averaging · SJF optimality, starvation, and SJF as priority scheduling
*(Pipe EOF and closing unused ends, `dup2` redirection, and sharing fd offsets after `fork` are reserved for M11. Convoy effect, SRTF, priority scheduling with aging, RR, MLQ/MLFQ and multiprocessor scheduling are coming in Batch 4.)*

**Completed: 60 / 300    Remaining: 240**
