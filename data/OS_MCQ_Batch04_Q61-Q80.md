# OS MCQ/MSQ Bank — Batch 4 (Q61–Q80)

**What this batch covers:**
- **M3, CPU scheduling completed (Q61–76):** SRTF, preemptive and non-preemptive priority scheduling, starvation and aging, Round Robin (Gantt charts, choosing the quantum, context-switch overhead, starvation myth), the convoy effect, multilevel queues (MLQ), multilevel feedback queues (MLFQ), processor affinity, load balancing, multiprocessor run queues, and how the algorithms reduce to each other.
- **M4, synchronization begins (Q77–80):** race conditions (lost-update trace) and the three requirements for a critical-section solution.

> Scheduling conventions for every numerical in this batch, unless the question says otherwise:
> - Times are in ms; there is no I/O and no context-switch cost.
> - A process that arrives at time *t* can be scheduled at *t*.
> - Ties among ready processes go to the earlier arrival.

---

## Q61. [M3 | SRTF | Medium | Numerical | OA]

Preemptive SJF (SRTF, Shortest Remaining Time First). An arriving process preempts the running one only if its burst is **strictly less** than the running process's remaining time. On an equal value, the running process continues.

| Process | Arrival | Burst |
|---|---|---|
| P1 | 0 | 9 |
| P2 | 1 | 4 |
| P3 | 3 | 2 |
| P4 | 6 | 3 |

What is the **average waiting time**?

- A. 6.00 ms
- B. 6.75 ms
- C. 3.00 ms
- D. 7.50 ms

**Answer:** C

**Why:** In SRTF, the scheduler re-decides every time a process **arrives** or **completes**. At each decision it runs whichever ready process has the least *remaining* time.

**Working:**
- **Formulas:** TAT = CT − AT; WT = TAT − BT.
- **Event-by-event:**
  1. t = 0: only P1 (remaining 9) → run P1.
  2. t = 1: P2 arrives (4) and P1 has 8 left. 4 < 8 → **preempt**; run P2.
  3. t = 3: P3 arrives (2) and P2 has 2 left. 2 is not < 2 → P2 continues.
  4. t = 5: P2 finishes. Ready: P1 (8), P3 (2) → run P3.
  5. t = 6: P4 arrives (3) and P3 has 1 left → P3 continues.
  6. t = 7: P3 finishes. Ready: P1 (8), P4 (3) → run P4.
  7. t = 10: P4 finishes → run P1 until 18.

```
| P1 | P2      | P3  | P4     | P1               |
0    1         5     7        10                 18
```

| Process | AT | BT | CT | TAT | WT |
|---|---|---|---|---|---|
| P1 | 0 | 9 | 18 | 18 | 9 |
| P2 | 1 | 4 | 5 | 4 | 0 |
| P3 | 3 | 2 | 7 | 4 | 2 |
| P4 | 6 | 3 | 10 | 4 | 1 |
| **Sum** | | | | 30 | 12 |

- Average WT = 12 / 4 = **3.00 ms**. (Average TAT = 7.5 ms.)

**Why the other options fail:**
- A: The non-preemptive SJF answer: P1 0–9, P3 9–11, P4 11–14, P2 14–18, giving WT = 0, 13, 6, 5. Catches forgetting to preempt.
- B: The FCFS answer: WT = 0, 8, 10, 9.
- D: 7.5 is the average *turnaround* time.

**Concept recap:** SRTF re-decides at every arrival and every completion. WT = CT − AT − BT, and preemption pieces count together.

**Trap alert:** With this data, flipping the tie rule (letting P3 preempt P2 at t = 3) still gives an average of 3.00, but the individual waits change. Always state the tie convention before computing.

**Interview angle:** "How do you compute waiting time for a preempted process?" Use WT = CT − AT − BT, which automatically adds up all the gaps. Follow-up: *"Why is SRTF hard to implement exactly?"* The remaining burst must be predicted (Q58).

---

## Q62. [M3 | SRTF properties | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. SRTF is the preemptive version of SJF.
- B. If burst lengths are known exactly, SRTF gives the minimum average waiting time among all schedules, preemptive or not.
- C. SRTF generally causes more context switches than non-preemptive SJF.
- D. SRTF cannot starve a long job.
- E. A real system can implement SRTF only by predicting remaining burst times.

**Answer:** A, B, C, E

**Why:**
- SRTF applies SJF's rule at every arrival and completion, which is what makes it preemptive.
- The same exchange argument as Q60 works for any pair of time slices, so SRTF is optimal over all schedules.
- Every preemption is an extra context switch.
- Real systems don't know burst lengths, so they must predict them.

**Why the other options fail:**
- D: A long job is preempted by every shorter arrival. A steady stream of short jobs starves it, just as in SJF.

**Concept recap:** SRTF = optimal average WT + starvation risk + more context switches + needs predicted bursts.

**Trap alert:** "Preemptive ⇒ fair" is wrong. Preemption *increases* the starvation risk for long jobs under SRTF.

**Interview angle:** "SJF vs SRTF?" Follow-up: *"Which gives better response time for a newly arriving short job?"* SRTF, because the short job doesn't have to wait for the current job to finish.

---

## Q63. [M3 | Preemptive priority | Medium | Numerical | OA]

Preemptive priority scheduling. **A lower number means a higher priority.** A higher-priority arrival preempts the running process immediately. Equal priorities run FCFS. A process arriving at the same instant another one completes is included in that decision.

| Process | Arrival | Burst | Priority |
|---|---|---|---|
| P1 | 0 | 5 | 3 |
| P2 | 1 | 3 | 1 |
| P3 | 2 | 2 | 4 |
| P4 | 4 | 4 | 2 |

What is the **average waiting time**?

- A. 4.25 ms
- B. 4.50 ms
- C. 3.75 ms
- D. 7.75 ms

**Answer:** A

**Why:** At every arrival and completion, the CPU goes to the ready process with the smallest priority number.

**Working:**
1. t = 0: only P1 (priority 3) → run P1.
2. t = 1: P2 arrives (priority 1, higher) → **preempt** P1, which has 4 left. Run P2.
3. t = 2: P3 arrives (priority 4, lower) → no preemption.
4. t = 4: P2 finishes, and P4 arrives (priority 2). Ready: P1 (3), P3 (4), P4 (2) → run P4.
5. t = 8: P4 finishes. Ready: P1 (3), P3 (4) → run P1 for its remaining 4.
6. t = 12: P1 finishes → run P3 until 14.

```
| P1 | P2     | P4       | P1       | P3  |
0    1        4          8          12    14
```

| Process | AT | BT | CT | TAT | WT |
|---|---|---|---|---|---|
| P1 | 0 | 5 | 12 | 12 | 7 |
| P2 | 1 | 3 | 4 | 3 | 0 |
| P3 | 2 | 2 | 14 | 12 | 10 |
| P4 | 4 | 4 | 8 | 4 | 0 |
| **Sum** | | | | 31 | 17 |

- Average WT = 17 / 4 = **4.25 ms**

**Why the other options fail:**
- B: The **non-preemptive** answer: P1 0–5, P2 5–8, P4 8–12, P3 12–14, giving WT = 0, 4, 10, 4.
- C: Reads the priorities backwards (higher number = higher priority): P1 0–2, P3 2–4, P1 4–7, P4 7–11, P2 11–14.
- D: 7.75 is the average *turnaround* time.

**Concept recap:** Preemptive priority = always run the highest-priority ready process. Always check whether a *lower number* means higher priority. Textbooks differ.

**Trap alert:** P3 waits 10 ms out of 14, even though it needs only 2. Low priority can mean very long waits (see Q64).

**Interview angle:** "How do you prevent low-priority starvation?" Aging. Follow-up: *"What is priority inversion?"* A high-priority task waits on a lock held by a low-priority task (covered in M4 and M12).

---

## Q64. [M3 | Starvation and aging | Easy | MCQ | OA + Interview]

Priorities range from 0 (highest) to 127 (lowest). With **aging**, a waiting process's priority number decreases by 1 every 15 minutes. Roughly how long until a process that starts at 127 reaches 0?

- A. 127 minutes
- B. About 15 hours
- C. Never; aging cannot overcome a priority that low
- D. About 32 hours

**Answer:** D

**Why:** **Aging** gradually raises the priority of a process the longer it waits, so eventually it becomes the highest-priority process and runs. That guarantees no process waits forever. Here it needs 127 steps of 15 minutes: 127 × 15 = 1905 min ≈ **31.75 hours**.

**Why the other options fail:**
- A: Treats each step as 1 minute.
- B: Arithmetic slip.
- C: The whole point of aging is that any process eventually reaches the top priority.

**Concept recap:** Time to reach the top = (levels to climb) × (aging interval). **Starvation** (also called indefinite blocking) = a ready process that never gets the CPU.

**Trap alert:** Aging guarantees *eventual* service, but slow aging can still mean a very long wait. Choosing the aging rate is a trade-off.

**Interview angle:** "What is starvation and how does aging fix it?" Follow-up: *"Starvation vs deadlock?"* A starving process *could* run but keeps losing to others. Deadlocked processes can never proceed (covered in M5).

---

## Q65. [M3 | Round Robin trace | Hard | Numerical | OA]

Round Robin with quantum q = 2. **Queue convention:** if a new process arrives at the same instant the running process's quantum expires, the **new arrival joins the tail first**, and the preempted process joins after it.

| Process | Arrival | Burst |
|---|---|---|
| P1 | 0 | 5 |
| P2 | 1 | 3 |
| P3 | 2 | 1 |
| P4 | 3 | 4 |

What is the **average waiting time**?

- A. 5.50 ms
- B. 5.00 ms
- C. 8.25 ms
- D. 1.75 ms

**Answer:** B

**Why:** RR gives each process up to q units, then sends it to the tail of the ready queue. The order of the queue at each expiry decides everything, so you must track the queue explicitly.

**Working:** (`rem` = remaining burst; the queue is shown after each event, front on the left)

| Time | Runs | Event | Ready queue afterwards |
|---|---|---|---|
| 0–2 | P1 (rem 5 → 3) | P2 arrives at 1; at 2, P3 arrives *and* P1 expires → P3 first, then P1 | [P2, P3, P1] |
| 2–4 | P2 (rem 3 → 1) | P4 arrives at 3; P2 expires at 4 | [P3, P1, P4, P2] |
| 4–5 | P3 (rem 1 → 0) | P3 **done** at 5 | [P1, P4, P2] |
| 5–7 | P1 (rem 3 → 1) | expires | [P4, P2, P1] |
| 7–9 | P4 (rem 4 → 2) | expires | [P2, P1, P4] |
| 9–10 | P2 (rem 1 → 0) | P2 **done** at 10 | [P1, P4] |
| 10–11 | P1 (rem 1 → 0) | P1 **done** at 11 | [P4] |
| 11–13 | P4 (rem 2 → 0) | P4 **done** at 13 | [] |

```
| P1 | P2 | P3 | P1 | P4 | P2 | P1 | P4 |
0    2    4    5    7    9    10   11   13
```

| Process | AT | BT | CT | TAT | WT = TAT − BT |
|---|---|---|---|---|---|
| P1 | 0 | 5 | 11 | 11 | 6 |
| P2 | 1 | 3 | 10 | 9 | 6 |
| P3 | 2 | 1 | 5 | 3 | 2 |
| P4 | 3 | 4 | 13 | 10 | 6 |
| **Sum** | | | | 33 | 20 |

- Average WT = 20 / 4 = **5.00 ms**

**Why the other options fail:**
- A: The result under the *opposite* convention (preempted P1 rejoins before the new arrival P3). P3 then runs 6–7 instead of 4–5, its WT rises to 4, and the average becomes 22/4. Catches ignoring the stated convention.
- C: 8.25 is the average *turnaround* time.
- D: 1.75 is the average *response* time (first run − arrival: 0, 1, 2, 4).

**Concept recap:** RR: run min(q, remaining), then go to the tail. WT = CT − AT − BT. RT = first run − AT.

**Trap alert:** The arrival-vs-preempted ordering at the same instant changes the answer. If an exam doesn't say, use the common textbook convention (new arrival first) and say so.

**Interview angle:** "Why is RR good for time-sharing?" Every process gets the CPU within a bounded time, so response time is good. Follow-up: *"What does RR cost compared with SJF?"* Usually a higher average turnaround, plus context-switch overhead.

---

## Q66. [M3 | Choosing the RR quantum | Medium | MCQ | Interview]

Which statement about choosing the Round Robin time quantum q is correct?

- A. As q → 0, RR becomes FCFS
- B. q should be much smaller than the context-switch time
- C. If q is larger than every CPU burst, RR behaves like FCFS; if q is very small, context-switch overhead dominates. A common rule of thumb is to make q longer than about 80% of CPU bursts
- D. A smaller q always gives a lower average turnaround time

**Answer:** C

**Why:**
- **Too large:** no process ever uses its full quantum, so nothing is preempted, and RR degenerates into FCFS, convoy effect included.
- **Too small:** the CPU spends a large share of its time context-switching instead of working.
- **About right:** most CPU bursts (≈80%) finish within one quantum, and q is much larger than the switch time.

**Why the other options fail:**
- A: Reversed. q → ∞ gives FCFS. As q → 0, RR approaches idealized "processor sharing", where each of n processes runs as if on a CPU 1/n as fast (ignoring overhead).
- B: Reversed. With q < switch time, most CPU time would be overhead.
- D: Turnaround doesn't improve steadily as q shrinks. Smaller q adds switches and can make turnaround worse.

**Concept recap:** q ≫ switch time; q ≈ just above the typical CPU burst; q → ∞ ⇒ FCFS.

**Trap alert:** Students think "smaller quantum = more responsive = better in every way". Responsiveness improves only until switch overhead starts to dominate.

**Interview angle:** "How would you choose the time quantum?" Follow-up: *"Typical values?"* Usually somewhere around 10–100 ms, while a context switch takes microseconds.

---

## Q67. [M3 | RR overhead and CPU utilization | Hard | Numerical | OA]

A system runs **10 I/O-bound tasks and 1 CPU-bound task**, all long-running.
- Each I/O-bound task issues an I/O request after every **1 ms** of CPU time, and each I/O takes **10 ms**. I/Os can proceed in parallel.
- Every context switch costs **0.1 ms**.
- Scheduling is Round Robin with **q = 10 ms**.

What is the CPU utilization (useful work ÷ total CPU time)?

- A. 94.8%
- B. 90.9%
- C. 99.0%
- D. 100%

**Answer:** A

**Why:** I/O-bound tasks never use a full quantum. Each one runs 1 ms, blocks, and costs a switch. Only the CPU-bound task uses the whole 10 ms. Useful work per round divided by total time per round gives the utilization.

**Working:**
- **Given:** 10 I/O-bound tasks, 1 CPU-bound task, q = 10, switch s = 0.1.
- **One full round of the ready queue:**
  - 10 I/O-bound tasks × (1 ms run + 0.1 ms switch) = 11 ms, of which 10 ms is useful.
  - 1 CPU-bound task × (10 ms run + 0.1 ms switch) = 10.1 ms, of which 10 ms is useful.
  - Round length = 11 + 10.1 = **21.1 ms**; useful work = 10 + 10 = **20 ms**.
- **Readiness check:** after blocking, an I/O-bound task's I/O finishes 10 ms later. Its next turn comes ≈ 21.1 ms later, so it is always ready when its turn arrives ✓.
- **Utilization** = 20 / 21.1 = **94.8%**.

**Why the other options fail:**
- B: The answer for q = 1 ms, where every task pays a switch after each 1 ms: 1/1.1 = 90.9%.
- C: Counts only the CPU-bound task (10/10.1).
- D: Ignores the context-switch cost.

**Concept recap:** Utilization = useful time ÷ (useful time + overhead) over one repeating cycle. A task that blocks early still pays a full switch.

**Trap alert:** Don't assume every process uses its whole quantum. I/O-bound processes give the CPU up early, and each early release costs a switch.

**Interview angle:** "How does the job mix affect RR overhead?" Follow-up: *"How could you reduce the overhead here?"* Use a larger quantum for CPU-bound work, as MLFQ's lower queues do (see Q71).

---

## Q68. [M3 | Round Robin properties | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. RR with a finite quantum cannot starve any process.
- B. Average turnaround time always improves as the quantum decreases.
- C. With n processes in the ready queue and quantum q, no process waits more than (n − 1)·q before its next quantum (ignoring switch time).
- D. RR typically has a higher average turnaround time than SJF, but a better response time.
- E. RR is a non-preemptive algorithm.

**Answer:** A, C, D

**Why:**
- In RR, every ready process is in one FIFO queue and gets the CPU once per round, so its wait is bounded by (n − 1)·q.
- New arrivals join the **tail**, so they can't jump ahead of a waiting process. Starvation is impossible.
- RR trades turnaround (jobs are interleaved, so each finishes later) for response time (each job starts quickly).

**Why the other options fail:**
- B: Shrinking q adds switches, and turnaround doesn't improve steadily (see Q66).
- E: RR is preemptive by definition: the timer interrupt takes the CPU away when the quantum ends.

**Concept recap:** Max wait between turns ≤ (n − 1)·q (plus switch costs). RR is starvation-free.

**Trap alert:** "RR can starve long jobs" is a common myth. Long jobs *finish later* under RR, but they always make progress.

**Interview angle:** "Can Round Robin cause starvation?" No, and say why: FIFO order and a bounded wait. Follow-up: *"Which algorithms can starve?"* SJF, SRTF, priority scheduling (without aging), and MLQ with fixed priorities.

---

## Q69. [M3 | Convoy effect | Easy | MCQ | OA + Interview]

Under FCFS, one long CPU-bound process holds the CPU. Meanwhile many I/O-bound processes finish their I/O and pile up in the ready queue while the devices sit idle. When the long process finally blocks, the I/O-bound ones run briefly and block again, and the CPU goes idle. What is this called?

- A. Belady's anomaly
- B. Convoy effect
- C. Priority inversion
- D. Thrashing

**Answer:** B

**Why:** Short jobs stuck behind one long job are like cars stuck behind a slow truck. FCFS is non-preemptive, so nothing can pass the long job. This lowers both CPU and device utilization.

**Why the other options fail:**
- A: More frames giving more page faults under FIFO page replacement (covered in M7).
- C: A high-priority task blocked by a low-priority task that holds a lock.
- D: The system spends most of its time paging instead of doing useful work (covered in M7).

**Concept recap:** Convoy effect = FCFS + one long CPU-bound job + many short I/O-bound jobs. Fixes: preemption (RR) or favouring short jobs (SJF/SRTF, MLFQ).

**Trap alert:** The convoy effect is a symptom of FCFS. It doesn't occur under RR, because the long job is preempted.

**Interview angle:** "What's wrong with FCFS?" Name the convoy effect and its average-waiting-time sensitivity to order (Q55). Follow-up: *"Which algorithm fixes it most directly?"* RR.

---

## Q70. [M3 | Multilevel queue (MLQ) | Medium | MCQ | OA]

The ready queue is split **permanently** into three queues: System (highest), Interactive (RR), and Batch (FCFS, lowest). Priority between the queues is **strictly fixed**, and processes never move between queues. Which is a real risk, and what is the standard fix?

- A. Interactive jobs get worse response than batch jobs; fix it by using FCFS in every queue
- B. Batch processes migrate upward; fix it by forbidding movement
- C. Context switches disappear entirely; fix it by adding a timer
- D. Batch jobs may starve while the higher queues stay busy; fix it by time-slicing between queues (e.g., 80% to interactive, 20% to batch) or by using MLFQ with aging

**Answer:** D

**Why:** With strict priority, a lower queue runs **only when every higher queue is empty**. A steady interactive load can therefore starve the batch queue indefinitely. Two fixes:
- Give each queue a guaranteed share of CPU time.
- Let processes move between queues (MLFQ), and use aging so a long-waiting process moves up.

**Why the other options fail:**
- A: Interactive jobs are served first, so their response is good.
- B: In MLQ, processes are assigned to a queue permanently; nothing migrates.
- C: The RR queue still has timer preemption, and higher queues still preempt lower ones.

**Concept recap:**
- MLQ: classify processes once (by type), each queue has its own algorithm, and the queues are scheduled by fixed priority or by time-slicing.
- MLFQ: processes *move* between queues based on how they behave.

**Trap alert:** MLQ and MLFQ are easy to confuse. The difference is movement: MLQ has none, MLFQ has it.

**Interview angle:** "MLQ vs MLFQ?" Follow-up: *"How does a process end up in a given MLQ queue?"* By a fixed property, such as its type, memory size, or user-assigned class.

---

## Q71. [M3 | MLFQ behaviour | Medium | MCQ | OA]

An MLFQ has three queues:
- Q0: RR with q = 8 ms
- Q1: RR with q = 16 ms
- Q2: FCFS

Rules: a new process enters Q0. A process that uses its **full** quantum is demoted one level. A higher queue always preempts a lower one. A CPU-bound process with a **30 ms** burst runs **alone** in the system. In which queue does it finish, and how much CPU time does it get in that queue?

- A. Q1, 22 ms
- B. Q0, 30 ms
- C. Q2, 6 ms
- D. Q1, 16 ms

**Answer:** C

**Why:** The process uses its whole quantum at each level, so it keeps getting demoted until it reaches FCFS:
1. Q0: runs 8 ms (8 used, 22 left) → used the full quantum → demoted.
2. Q1: runs 16 ms (24 used, 6 left) → used the full quantum → demoted.
3. Q2: runs the remaining **6 ms** under FCFS → finishes.

**Why the other options fail:**
- A: Forgets that 22 ms remaining exceeds Q1's 16 ms quantum.
- B: Ignores demotion.
- D: Stops after Q1, with 6 ms still remaining.

**Concept recap:** In MLFQ, a CPU-bound process sinks to the long-quantum or FCFS queues, while an I/O-bound process (short bursts) stays in the high-priority queues. MLFQ thereby *learns* each process's burst behaviour, approximating SJF without knowing bursts in advance.

**Trap alert:** Demotion happens only when the process uses its **full** quantum. Blocking early does not demote it (see Q72).

**Interview angle:** "Explain MLFQ." Follow-up: *"What parameters define an MLFQ?"* The number of queues, each queue's algorithm and quantum, the promotion and demotion rules, and which queue a new process enters.

---

## Q72. [M3 | MLFQ properties | Medium | MSQ | Interview]

**Select all that apply.**

- A. MLFQ favours interactive and I/O-bound processes, because their short bursts keep them in the high-priority queues.
- B. A process that blocks for I/O is always moved to the lowest queue.
- C. Without a periodic priority boost (or aging), long CPU-bound jobs in the lowest queue can starve if interactive jobs keep arriving.
- D. A program can game a naive MLFQ by doing a tiny I/O just before its quantum expires, so it stays at the top. The fix is to charge the total CPU time used at a level, rather than resetting on every block.
- E. MLFQ requires burst lengths to be known in advance.

**Answer:** A, C, D

**Why:**
- MLFQ judges processes by what they *do*: using the full quantum looks CPU-bound, and blocking early looks interactive.
- Its weaknesses are **starvation** of the bottom queue and **gaming** by programs that block on purpose.
- Two standard fixes: periodically boost every process back to the top queue, and track how much of its quantum a process has used at a level, even across blocks.

**Why the other options fail:**
- B: Blocking early *keeps* (or raises) a process's level. Only using the full quantum demotes it.
- E: MLFQ learns behaviour as processes run, so no advance knowledge is needed.

**Concept recap:** Rules of a well-designed MLFQ: run higher queues first; RR within a level; new processes start at the top; demote when the time allotment for a level is used up; boost everyone periodically.

**Trap alert:** "MLFQ = MLQ with RR everywhere" is wrong. Feedback, meaning movement between levels, is the defining feature.

**Interview angle:** "How would you stop a process from gaming MLFQ?" Follow-up: *"Why boost every process periodically?"* It prevents starvation, and it lets a process whose behaviour changes (such as a CPU-bound job turning interactive) move back up.

---

## Q73. [M3 | Processor affinity | Medium | MCQ | OA + Interview]

Why does a multiprocessor scheduler try to keep a thread on the core it last ran on?

- A. Its data is probably still in that core's caches ("warm cache"); moving it means refilling caches on the new core
- B. Different cores have different instruction sets
- C. To avoid context switches entirely
- D. It is required for mutual exclusion

**Answer:** A

**Why:** A running thread fills its core's L1/L2 caches (and TLB) with its working data. If it moves to another core, those caches start "cold", so it suffers many cache misses until they fill again. **Processor affinity** is the scheduler's preference for keeping a thread where it last ran.

**Why the other options fail:**
- B: Typical SMP cores are identical. (Heterogeneous designs exist, but that isn't the reason for affinity.)
- C: Context switches still happen on the same core.
- D: Mutual exclusion is provided by locks, not by where a thread runs.

**Concept recap:**
- **Soft affinity:** the scheduler tries to keep the thread on the same core, but may move it.
- **Hard affinity:** the thread is pinned to a set of cores (on Linux, `sched_setaffinity()` or `taskset`).

**Trap alert:** Affinity conflicts with load balancing (Q75). A busy core holding many affine threads may need to give some away.

**Interview angle:** "What is processor affinity?" Follow-up: *"When would you pin a thread with hard affinity?"* For latency-critical work, such as trading systems or real-time loops, to avoid migration and cache disruption.

---

## Q74. [M3 | Load balancing | Easy | MCQ | OA]

Each core has its own run queue. An idle core takes a waiting task from a busy core's run queue. This is called:

- A. Push migration
- B. Pull migration
- C. Gang scheduling
- D. Hard affinity

**Answer:** B

**Why:** **Load balancing** spreads work evenly across cores. There are two directions:
- **Pull migration:** the idle core takes work from a busy core.
- **Push migration:** a periodic task moves work from overloaded cores to less loaded ones.

Linux uses both.

**Why the other options fail:**
- A: In push migration, the *busy* side, or a balancer, initiates the move.
- C: Gang scheduling runs related threads at the same time on different cores.
- D: Hard affinity *prevents* migration.

**Concept recap:** Load balancing matters only with per-core run queues. With one shared queue, any idle core simply takes the next task.

**Trap alert:** Who initiates the move decides the name: the idle core → pull; the busy side or a balancer → push.

**Interview angle:** "How does a multicore scheduler keep cores evenly loaded?" Follow-up: *"What does balancing cost?"* Lost cache affinity, plus the locking needed to touch another core's queue.

---

## Q75. [M3 | Multiprocessor scheduling trade-offs | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. Per-core run queues reduce lock contention compared with a single shared run queue.
- B. Load balancing can conflict with processor affinity.
- C. On a NUMA machine, migrating a thread to a core on another memory node can slow down its memory accesses.
- D. Two hyperthreads (SMT) on one physical core deliver the same throughput as two separate physical cores.
- E. With a single shared ready queue, two cores could pick the same process unless the queue is protected by a lock.

**Answer:** A, B, C, E

**Why:**
- A single shared queue must be locked on every scheduling decision, and that lock becomes a bottleneck as cores are added. Per-core queues avoid it.
- Balancing moves threads, which throws away their warm caches.
- On **NUMA** (Non-Uniform Memory Access) machines, memory attached to another node is slower to reach.
- Hyperthreads share one core's execution units and caches.

**Why the other options fail:**
- D: SMT (simultaneous multithreading) hides stalls by running two instruction streams on one core's units. It usually gives a modest throughput gain, far from 2×, because the threads compete for the same execution units and caches.

**Concept recap:** Multicore scheduling balances lock contention (shared vs per-core queues), affinity (cache warmth), load balance, and NUMA locality.

**Trap alert:** The OS sees each hyperthread as a CPU, but treating hyperthreads as full cores overestimates capacity.

**Interview angle:** "Single shared run queue vs per-core queues?" Follow-up: *"How does Linux's scheduler use NUMA topology?"* It groups CPUs into "scheduling domains" and prefers to balance within a node before moving work across nodes.

---

## Q76. [M3 | Relationships between algorithms | Hard | MCQ | OA]

Which of these equivalences is **FALSE**?

- A. Priority scheduling where priority = arrival time (earlier = higher) behaves as FCFS
- B. Preemptive priority scheduling where priority = remaining burst (shorter = higher) behaves as SRTF
- C. Non-preemptive priority scheduling where priority = next CPU burst (shorter = higher) behaves as SJF
- D. SJF with all bursts equal behaves as Round Robin

**Answer:** D

**Why:** When every burst is equal, SJF has no length difference to sort by, so it falls back to its tie-break rule (arrival order) and becomes **FCFS**. Each job then runs **to completion** with no time slicing, which is not RR behaviour.

**Why the other options fail:** These three are true, which is why they're not the answer:
- A: Earliest arrival first is exactly FCFS.
- B: Always running the shortest remaining time, with preemption, is the definition of SRTF.
- C: SJF is priority scheduling where priority = next burst (Q60 E).

**Concept recap:** Priority scheduling is the general framework. FCFS, SJF and SRTF are special cases defined by what the priority measures. RR differs because of **time slicing**, not because of its ordering rule.

**Trap alert:** "Equal bursts → everyone gets turns → RR" is wrong. Without preemption there are no turns.

**Interview angle:** "Which algorithms are special cases of priority scheduling?" Follow-up: *"Is RR a special case of anything?"* RR with q → ∞ is FCFS, and RR is also one level of an MLFQ.

---

## Q77. [M4 | Race condition: lost updates | Hard | Code-trace | OA + Interview]

Two threads share `int count = 0;`. Each runs:

```c
for (int k = 0; k < 3; k++)
    count++;          // compiled as: r = count;  r = r + 1;  count = r;
```

Each thread has its own register `r`. The load, add and store of each increment can interleave arbitrarily with the other thread's. No locks are used. What is the full range of possible final values of `count`?

- A. 3 to 6
- B. 2 to 6
- C. Exactly 6
- D. 1 to 6

**Answer:** B

**Why:** `count++` is **not atomic**. It's a separate load, add and store, so one thread's store can overwrite the other's update (a **lost update**). Many interleavings are possible:
- **6** is the maximum: the increments run one after another with no overlap.
- Many updates can be lost, but never enough to reach 1.

**Working:**
- **A schedule that ends at 2.** Below, L = load, S = add + store, and `count` is shown after each step.
  1. T1: L (r₁ = 0)
  2. T2: two complete increments → count = 2
  3. T1: S → count = r₁ + 1 = **1** (overwrites T2's work)
  4. T2: L for its 3rd increment (r₂ = 1)
  5. T1: its 2nd and 3rd increments, complete → count = 3
  6. T2: S → count = r₂ + 1 = **2**. Both threads are now done; final = 2 ✓
- **Why 1 is impossible:**
  - Every store writes (some value loaded) + 1, and every loaded value is ≥ 0, so every store writes ≥ 1. After any thread's first store, `count` is ≥ 1 forever.
  - Consider the thread that does the last store. Its last iteration's load happens after its own first store, so that load reads ≥ 1, and its final store writes ≥ 2.
  - So the final value is ≥ 2.

**Why the other options fail:**
- A: Assumes each thread's own 3 increments always "count", ignoring that the other thread can overwrite them.
- C: Assumes `++` is atomic. This is the core misconception behind races.
- D: Ignores the argument above showing the final value is at least 2.

**Concept recap:** Two threads each doing n non-atomic increments (n ≥ 2) give a final value in **[2, 2n]**. The fix is to make the read-modify-write atomic: use a mutex, or `atomic_fetch_add` / `std::atomic<int>`.

**Trap alert:** `volatile` does **not** fix this. It only forces real memory accesses; it doesn't make load-add-store atomic.

**Interview angle:** "Is `i++` thread-safe?" No, and explain the load/add/store steps. Follow-up: *"Smallest fix?"* `std::atomic<int>` with `fetch_add`, or a mutex around the increment.

---

## Q78. [M4 | What a race condition is | Easy | MCQ | Interview]

What is a **race condition**?

- A. When the result depends on the timing or interleaving of accesses by several threads or processes to shared data, where at least one access is a write
- B. When two processes compete for the CPU
- C. When a process runs faster than expected
- D. When two threads only read the same variable at the same time

**Answer:** A

**Why:** If concurrent tasks touch shared data without coordination and at least one of them writes, different interleavings produce different results, sometimes wrong ones. The code segment where shared data is accessed is the **critical section**, and preventing races means controlling entry to it.

**Why the other options fail:**
- B: Competing for the CPU is just scheduling. It's harmless by itself.
- C: Speed is irrelevant. What matters is the *ordering* of shared accesses.
- D: Concurrent reads alone are safe, because nothing changes.

**Concept recap:** Race = shared data + concurrent access + at least one write + no synchronization. The cure is to make the critical section execute as if it were indivisible.

**Trap alert:** Races are **timing-dependent**. Code can pass thousands of test runs and then fail in production. "It worked when I tested it" proves nothing.

**Interview angle:** "What is a race condition? Give an example." Use a bank balance updated by two withdrawals. Follow-up: *"How would you detect one?"* Thread sanitizers such as `-fsanitize=thread` (TSan), code review of shared state, and stress testing.

---

## Q79. [M4 | Critical-section requirements: strict alternation | Medium | MCQ | OA]

Two processes use this protocol:

```c
// shared:  int turn = 0;
// Process Pi  (i is 0 or 1;  j = 1 - i)
while (true) {
    while (turn != i)  ;     // busy-wait
    /* critical section */
    turn = j;
    /* remainder section */
}
```

Which requirement of a correct critical-section solution does this **violate**?

- A. Mutual exclusion
- B. Bounded waiting
- C. Progress
- D. Both mutual exclusion and bounded waiting

**Answer:** C

**Why:** The processes are forced to take turns strictly. Suppose P0 finishes its critical section (setting `turn = 1`) and P1 is busy in its long remainder section, not wanting to enter. If P0 wants to enter again, it can't, because it's P1's "turn". A process that doesn't even want the critical section is blocking one that does. That breaks **progress**.

**Why the other options fail:**
- A: Mutual exclusion **holds**. `turn` has only one value at a time, so only that process can be past the busy-wait.
- B: Bounded waiting **holds**. After a request, the other process can enter at most once before you get a turn.
- D: Both properties actually hold.

**Concept recap:** A correct critical-section solution needs **mutual exclusion**, **progress**, and **bounded waiting** (defined in Q80).

**Trap alert:** Students see "it alternates fairly" and call it correct. Fairness here costs progress whenever the two processes run at different rates.

**Interview angle:** "What's wrong with a simple turn variable?" Follow-up: *"How does Peterson's algorithm fix it?"* It adds a `flag[i]` ("I want in"), so you defer to the other process only if it actually wants to enter (next batch).

---

## Q80. [M4 | Critical-section requirements: definitions | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Mutual exclusion: at most one process is in its critical section at any time.
- B. Progress: if no process is in its critical section and some wish to enter, the choice of which enters next cannot be postponed indefinitely, and only processes not in their remainder sections take part in that choice.
- C. Bounded waiting: after a process requests entry, there is a limit on how many times other processes may enter their critical sections before the request is granted.
- D. Bounded waiting means each critical section must finish within a fixed amount of time.
- E. Any solution that satisfies mutual exclusion and progress automatically satisfies bounded waiting.

**Answer:** A, B, C

**Why:** The three properties guard against three different failures:
- Mutual exclusion → **correctness**: no two processes in the critical section at once.
- Progress → **no needless blocking**: the system doesn't stall, and processes that don't want in don't get a say.
- Bounded waiting → **fairness**: no process is bypassed forever.

**Why the other options fail:**
- D: Bounded waiting limits the **number of times others enter ahead of you**, not the critical section's running time.
- E: False. A simple test-and-set spinlock gives mutual exclusion and progress, yet an unlucky thread can lose the race every time and starve. Bounded waiting is a separate property (see the TAS question in the next batch).

**Concept recap:** Mutual exclusion = safety. Progress = liveness (the system keeps moving). Bounded waiting = fairness. A correct solution must satisfy **all three**.

**Trap alert:** Progress is the one students mis-state. Its key clause is that **processes in their remainder section must not influence** who enters next. Strict alternation (Q79) violates exactly this.

**Interview angle:** "What three properties must a critical-section solution satisfy?" Give a one-line definition of each. Follow-up: *"Which property does a plain spinlock lack?"* Bounded waiting.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 61 | C | 66 | C | 71 | C | 76 | D |
| 62 | A, B, C, E | 67 | A | 72 | A, C, D | 77 | B |
| 63 | A | 68 | A, C, D | 73 | A | 78 | A |
| 64 | D | 69 | B | 74 | B | 79 | C |
| 65 | B | 70 | D | 75 | A, B, C, E | 80 | A, B, C |

---

**Concepts tested so far (to avoid repeats):** Batches 1–3, plus: SRTF trace with strict-less preemption · SRTF optimality, switches, starvation · preemptive priority trace (lower number = higher priority) · aging arithmetic · RR trace with the arrival-vs-preempted convention · choosing the quantum (q → ∞ ⇒ FCFS, the 80% rule) · RR utilization with switch cost and an I/O-bound mix · RR starvation-freedom and the (n − 1)q bound · convoy effect · MLQ starvation and time-slicing between queues · MLFQ demotion trace · MLFQ gaming and priority boost · processor affinity (soft/hard) · pull vs push migration · per-core queues, NUMA, SMT · algorithms as special cases of priority scheduling · lost-update range [2, 2n] · race condition definition · strict alternation violating progress · ME/progress/bounded-waiting definitions
*(The M3 module is complete. Peterson, TAS, CAS, mutexes, spinlocks and semaphores follow in Batch 5.)*

**Completed: 80 / 300    Remaining: 220**
