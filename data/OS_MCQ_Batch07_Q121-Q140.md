# OS MCQ/MSQ Bank — Batch 7 (Q121–Q140)

**What this batch covers:**
- **M5, deadlocks completed (Q121–134):**
  - Prevention (breaking each of the four conditions).
  - Safe vs unsafe states.
  - Banker's safety algorithm (single and multiple resource types) and resource-request algorithm.
  - Avoidance with claim edges.
  - Detection with the Request matrix, and recovery.
  - Deadlock vs starvation vs livelock.
  - What real OSes do (the ostrich approach).
  - Minimum instances that rule out deadlock.
- **M6, memory management begins (Q135–140):** memory hierarchy, logical vs physical addresses, address-binding times, the MMU with base/limit registers.

> **Banker's conventions:**
> - Need = Max − Allocation. Vectors are compared component-wise: X ≤ Y means every component of X is ≤ the matching component of Y.
> - When several processes could run next, pick the **lowest-indexed unfinished process whose Need ≤ Work**.

---

## Q121. [M5 | Prevention: hold and wait | Easy | MCQ | OA]

A system requires every process to request, and be granted, **all** its resources before it starts executing. Which deadlock condition does this eliminate?

- A. Hold and wait
- B. Mutual exclusion
- C. No preemption
- D. Circular wait

**Answer:** A

**Why:** **Hold and wait** means "holding something while waiting for more". A process that gets everything up front never waits while holding anything, so that condition can't arise. An alternative protocol has the same effect: a process may request resources only when it holds none, so it must release everything before asking again.

**Why the other options fail:**
- B: The resources are still used exclusively.
- C: Nothing is ever taken away from a process.
- D: Circular wait becomes impossible *as a consequence*, but the condition directly broken is hold-and-wait.

**Concept recap:** The price of breaking hold-and-wait:
- **Low resource utilization:** resources sit allocated but idle.
- **Possible starvation:** a process that needs many popular resources may never get them all at once.

**Trap alert:** Students often answer "circular wait" because a cycle can't form. Name the condition the rule **directly** violates.

**Interview angle:** "How would you prevent deadlock by attacking hold-and-wait?" Follow-up: *"Why is that rarely practical?"* Processes often don't know their full needs in advance, and holding everything for the whole run wastes resources.

---

## Q122. [M5 | Deadlock prevention | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Imposing a total order on resource types, and requiring every process to request them in increasing order, prevents circular wait.
- B. Making every resource sharable is always possible, so mutual exclusion is easy to break.
- C. Preempting the resources of a process that requests something unavailable works for resources whose state is easy to save and restore (CPU registers, memory pages), but not for mutex locks or a half-finished print job.
- D. Hold-and-wait prevention can lower resource utilization and cause starvation.
- E. Prevention requires knowing each process's maximum future resource needs.

**Answer:** A, C, D

**Why:** Prevention means designing the system so that one Coffman condition can never hold:
- **A:** If every process acquires resources in increasing order, a cycle would need some process to hold a higher-numbered resource while waiting for a lower one. That violates the rule, so no cycle can form.
- **C:** Preemption works only when you can save the resource's state and restore it later.
- **D:** See Q121.

**Why the other options fail:**
- B: Some resources are inherently non-sharable, such as a mutex or a write lock. Read-only files are sharable, but you can't make everything sharable.
- E: Knowing maximum needs is what **avoidance** (Banker's) requires. Prevention needs no such knowledge.

**Concept recap:** Ways to break each condition:
- **ME:** make resources sharable (rarely possible).
- **Hold-and-wait:** request everything at once.
- **No preemption:** preempt resources whose state can be saved.
- **Circular wait:** acquire in a global order. This is the most practical one.

**Trap alert:** **Prevention vs avoidance.** Prevention restricts *how* processes may request resources. Avoidance decides *whether* to grant a request by looking ahead.

**Interview angle:** "Which deadlock condition do engineers usually break in practice?" Circular wait, through lock ordering (Q113). Follow-up: *"How does the Linux kernel enforce lock order?"* Its `lockdep` validator, enabled in debug builds, records the order in which locks are taken and warns about orders that could form a cycle.

---

## Q123. [M5 | Safe vs unsafe states | Medium | MCQ | OA + Interview]

Which statement is correct?

- A. Every unsafe state is a deadlocked state
- B. From a safe state, the OS can avoid deadlock by granting only requests that keep the state safe; an unsafe state *may* lead to deadlock but is not necessarily deadlocked
- C. A deadlocked state can be safe
- D. "Safe" means no process is currently waiting

**Answer:** B

**Why:** A state is **safe** if there is at least one **safe sequence**: an order in which every process can obtain its maximum remaining need, finish, and release what it holds. In an **unsafe** state, no such guarantee exists. Deadlock *could* follow, but processes may never actually ask for their full maximum, so it might not.

**Why the other options fail:**
- A: This is the misconception the question targets. Unsafe ≠ deadlocked.
- C: In a deadlocked state, some processes can never finish, so no safe sequence can exist.
- D: Processes may well be waiting in a safe state. Safety is about whether everyone *can* finish.

**Concept recap:** Deadlocked states ⊂ unsafe states. Safe states and deadlocked states never overlap. Avoidance means **never leaving the safe region**.

**Trap alert:** "Unsafe ⇒ deadlock is certain" is wrong. Unsafe only means the OS can no longer *guarantee* that deadlock won't happen.

**Interview angle:** "What is a safe state?" Define it using the safe sequence. Follow-up: *"Why not just let the system become unsafe and see?"* Because, once unsafe, a specific pattern of requests can force a deadlock, and the OS can't prevent it any more.

---

## Q124. [M5 | Safe sequence, one resource type | Medium | Numerical | OA]

A system has **10** instances of one resource type.

| Process | Max | Allocation |
|---|---|---|
| P0 | 7 | 3 |
| P1 | 4 | 2 |
| P2 | 5 | 2 |

Is the state safe, and which option gives a valid safe sequence?

- A. Unsafe
- B. Safe; <P0, P1, P2>
- C. Safe; <P1, P2, P0>
- D. Safe; and the **only** safe sequence is <P2, P0, P1>

**Answer:** C

**Why:** A sequence is safe if each process in turn can get its remaining Need from what's available plus everything released by the processes before it.

**Working:**
- **Formulas:** Need = Max − Allocation; Available = Total − ΣAllocation.
- Need: P0 = 4, P1 = 2, P2 = 3.
- Available = 10 − (3 + 2 + 2) = **3**.
- Check <P1, P2, P0>:

| Step | Process | Need ≤ Work? | Work after it finishes |
|---|---|---|---|
| start | | | 3 |
| 1 | P1 | 2 ≤ 3 ✓ | 3 + 2 = 5 |
| 2 | P2 | 3 ≤ 5 ✓ | 5 + 2 = 7 |
| 3 | P0 | 4 ≤ 7 ✓ | 7 + 3 = 10 |

All three finish → the state is safe.

**Why the other options fail:**
- A: The table above shows a safe sequence exists.
- B: P0 needs 4 but only 3 are available, so P0 can't go first.
- D: <P2, P0, P1> is safe (3 ≤ 3 → 5; 4 ≤ 5 → 8; 2 ≤ 8), but it is **not the only one**. <P1, P2, P0> and <P1, P0, P2> also work.

**Concept recap:** A state is safe ⇔ **at least one** safe sequence exists. There are often several.

**Trap alert:** Don't assume the safe sequence is unique. Read words like "only" carefully.

**Interview angle:** "Walk me through checking whether a state is safe." Follow-up: *"What is the safety algorithm's complexity?"* O(m · n²) for m resource types and n processes.

---

## Q125. [M5 | Banker's safety algorithm | Hard | Numerical | OA + Interview]

Resource types A, B, C, with totals **A = 9, B = 6, C = 6**.

| Process | Allocation (A B C) | Max (A B C) |
|---|---|---|
| P0 | 1 1 0 | 4 3 2 |
| P1 | 2 0 1 | 3 2 2 |
| P2 | 1 2 1 | 2 3 3 |
| P3 | 0 1 2 | 1 1 3 |
| P4 | 2 1 0 | 5 2 1 |

Using the lowest-index convention, which sequence does the safety algorithm produce?

- A. The state is unsafe
- B. Safe; <P2, P0, P1, P3, P4>
- C. Safe; <P1, P3, P4, P2, P0>
- D. Safe; <P0, P2, P1, P3, P4>

**Answer:** B

**Why:** The safety algorithm starts with Work = Available and repeatedly finds an unfinished process whose Need ≤ Work. It lets that process finish and adds its Allocation back to Work. If every process finishes, the state is safe.

**Working:**
- **Step 1, compute Available** = Total − ΣAllocation.
  - ΣAllocation: A = 1+2+1+0+2 = 6, B = 1+0+2+1+1 = 5, C = 0+1+1+2+0 = 4.
  - Available = (9−6, 6−5, 6−4) = **(3, 1, 2)**.
- **Step 2, compute Need** = Max − Allocation:

| Process | Need (A B C) |
|---|---|
| P0 | 3 2 2 |
| P1 | 1 2 1 |
| P2 | 1 1 2 |
| P3 | 1 0 1 |
| P4 | 3 1 1 |

- **Step 3, run the algorithm:**

| Step | Work before | Lowest-index process with Need ≤ Work | + Allocation | Work after |
|---|---|---|---|---|
| 1 | (3,1,2) | P0 ✗ (B: 2>1), P1 ✗ (B: 2>1), **P2 (1,1,2) ✓** | +(1,2,1) | (4,3,3) |
| 2 | (4,3,3) | **P0 (3,2,2) ✓** | +(1,1,0) | (5,4,3) |
| 3 | (5,4,3) | **P1 (1,2,1) ✓** | +(2,0,1) | (7,4,4) |
| 4 | (7,4,4) | **P3 (1,0,1) ✓** | +(0,1,2) | (7,5,6) |
| 5 | (7,5,6) | **P4 (3,1,1) ✓** | +(2,1,0) | (9,6,6) ✓ = Total |

All processes finish → **safe**, with sequence <P2, P0, P1, P3, P4>. (Work ends equal to Total, a quick self-check.)

**Why the other options fail:**
- A: The table shows every process can finish.
- C: The famous textbook answer, but it doesn't fit *this* data: P1 needs B = 2 when only 1 is available. Catches memorizing instead of computing.
- D: P0 can't go first, because it needs B = 2 > 1.

**Concept recap:**
- Need = Max − Alloc; Work = Available.
- Loop: find i with Need_i ≤ Work → Work += Alloc_i → mark i finished.
- Safe ⇔ every process gets marked finished.

**Trap alert:** Compare **every component**. P2's Need (1,1,2) fits (3,1,2) exactly on B and C. Students often reject it by misreading 1 ≤ 1 or 2 ≤ 2.

**Interview angle:** "Explain the Banker's algorithm." Name the data structures (Available, Max, Allocation, Need) and the safety loop. Follow-up: *"Why don't real OSes use it?"* Processes don't declare their maximum needs, the number of processes changes over time, and running it on every request costs O(m·n²).

---

## Q126. [M5 | Banker's resource-request algorithm | Hard | Numerical | OA]

Same state as Q125: Available = (3, 1, 2). **P0 requests (1, 1, 2).** What should the OS do?

- A. Grant it immediately, because the request ≤ Available
- B. Raise an error, because the request exceeds P0's maximum claim
- C. Grant it, because the resulting state is safe
- D. Deny it for now (P0 waits): pretending to grant leaves Available = (2, 0, 0), no process's Need then fits, so the state would be unsafe

**Answer:** D

**Why:** The resource-request algorithm has three checks:
1. Request ≤ Need? If not, it's an **error**: the process has exceeded its declared maximum.
2. Request ≤ Available? If not, the process **waits**.
3. **Pretend** to grant the request and run the safety algorithm. If the result is safe, grant it. Otherwise roll back and make the process **wait**.

**Working:**
1. Request (1,1,2) ≤ Need_P0 (3,2,2)? ✓
2. Request (1,1,2) ≤ Available (3,1,2)? ✓
3. Pretend to grant:
   - Available = (3,1,2) − (1,1,2) = **(2, 0, 0)**
   - Alloc_P0 = (1,1,0) + (1,1,2) = (2,2,2)
   - Need_P0 = (3,2,2) − (1,1,2) = (2,1,0)
4. Safety check with Work = (2, 0, 0):

| Process | Need | ≤ (2,0,0)? |
|---|---|---|
| P0 | 2 1 0 | ✗ (B) |
| P1 | 1 2 1 | ✗ (B, C) |
| P2 | 1 1 2 | ✗ (B, C) |
| P3 | 1 0 1 | ✗ (C) |
| P4 | 3 1 1 | ✗ (A, B, C) |

5. No process can proceed → **unsafe** → roll back, and P0 waits.

**Why the other options fail:**
- A: "Request ≤ Available ⇒ grant" skips the safety check, which is the entire point of avoidance.
- B: (1,1,2) ≤ (3,2,2), so it is within P0's claim.
- C: The safety check fails, as shown above.

**Concept recap:** Request ≤ Need → Request ≤ Available → pretend → safe? → grant : roll back and wait.

**Trap alert:** A request can be **available** and still be refused. Avoidance refuses requests that would lead into an unsafe state.

**Interview angle:** "How does Banker's decide on a request?" Give the three steps. Follow-up: *"What's the cost of being conservative?"* Lower utilization, because resources sit idle even when they're free.

---

## Q127. [M5 | Detection vs Banker's safety | Medium | MCQ | Interview]

In the deadlock-**detection** algorithm (multiple instances per type), what replaces the Need matrix used in the Banker's safety algorithm, and why?

- A. The Max matrix
- B. The Request matrix (each process's current outstanding requests), because detection asks whether processes can finish given what they are waiting for right now, assuming optimistically that a process that gets its request will finish and release everything
- C. The Allocation matrix
- D. The Available vector

**Answer:** B

**Why:**
- **Avoidance** looks at the *worst case in the future* ("what if everyone asks for their Max?"), so it uses Need.
- **Detection** looks at the *present* ("is anyone stuck right now?"), so it uses Request.

A process whose current request can be satisfied is assumed to finish. If it later asks for more, the next detection run will catch any resulting deadlock.

**Why the other options fail:**
- A: Detection doesn't need Max. That's part of why it's practical when processes can't declare their needs.
- C, D: Allocation and Available are still used, but they don't play Need's role.

**Concept recap:** Detection is the same loop as the safety algorithm, with Need replaced by Request. Two details:
- A process with **zero allocation** can be marked finished at the start.
- Any process **left unfinished** at the end is deadlocked.

**Trap alert:** Avoidance answers "could this state *become* deadlocked?" Detection answers "*is* it deadlocked now?" Different questions, different matrices.

**Interview angle:** "Avoidance vs detection: when would you choose detection?" When deadlocks are rare and declaring maximum needs is impractical; databases are a typical case. Follow-up: *"How often should detection run?"* It's a trade-off between overhead and how long deadlocked processes stay stuck. Some systems run it only when CPU utilization drops sharply.

---

## Q128. [M5 | Running the detection algorithm | Medium | MCQ | OA]

Resource types A and B. **Available = (0, 0).**

| Process | Allocation (A B) | Request (A B) |
|---|---|---|
| P1 | 1 0 | 0 1 |
| P2 | 0 1 | 1 0 |
| P3 | 0 0 | 0 0 |
| P4 | 1 0 | 1 0 |

Which processes are deadlocked?

- A. P1 and P2 only
- B. All four
- C. P1, P2 and P4. P3 can finish because it requests nothing, but it holds nothing, so finishing frees nothing
- D. None

**Answer:** C

**Why:** Mark as finished every process whose Request ≤ Work, add back its Allocation, and repeat. Whatever is left unmarked at the end is deadlocked.

**Working:**
- Work = (0, 0).
- **Pass 1:** P1 (0,1) ✗; P2 (1,0) ✗; **P3 (0,0) ✓** → Work += (0,0) = (0,0); P4 (1,0) ✗.
- **Pass 2:** Work is still (0,0), so P1, P2 and P4 all fail again. Nothing changes → stop.
- Unfinished: **P1, P2, P4** → deadlocked.
- In plain words: P1 waits for B, which only P2 holds. P2 and P4 both wait for A, which only P1 and P4 hold. Nobody among them can release anything.

**Why the other options fail:**
- A: Misses P4. P4 holds the other instance of A and waits for an A, so it's stuck in the same knot.
- B: P3 isn't waiting for anything, so it isn't deadlocked.
- D: Nothing ever adds to Work, so no stuck process can be marked finished.

**Concept recap:** A process with Request = 0 always finishes, but it frees only what it *holds*. A deadlocked set can include processes outside the obvious two-process cycle.

**Trap alert:** "Someone can finish ⇒ no deadlock" is false. Track **what that process releases**.

**Interview angle:** Follow-up: *"Having found {P1, P2, P4}, which would you abort?"* Pick the victim whose abort frees the most and costs the least. Aborting P1 frees one A, which P2 can then use (Q129).

---

## Q129. [M5 | Recovery from deadlock | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Aborting all deadlocked processes is simple, but loses all the work they have done.
- B. Aborting deadlocked processes one at a time requires re-running detection after each abort.
- C. Recovering by resource preemption requires rolling the victim back to a safe point, such as a checkpoint or a restart.
- D. If the same victim is chosen every time, it can starve; a common fix is to include the number of times it has been rolled back in the cost of choosing it.
- E. Recovery is unnecessary, because deadlocks eventually resolve themselves.

**Answer:** A, B, C, D

**Why:** There are two families of recovery:
1. **Terminate processes:** abort all of them (cheap to decide, expensive in lost work), or abort one at a time (cheaper in work, but detection must run again after each abort).
2. **Preempt resources:** take resources from a victim. The victim must be **rolled back**, because its partial work assumed it had those resources.

Victim selection weighs priority, work done, resources held, and resources still needed. It must also avoid **starvation** of a victim that keeps being picked.

**Why the other options fail:**
- E: By definition, nothing inside a deadlock can make progress. Without outside intervention it lasts forever.

**Concept recap:** Recovery = choose a victim → terminate it or preempt its resources (with rollback) → avoid repeatedly picking the same victim.

**Trap alert:** Aborting a process in the middle of updating a file or database can leave the data **inconsistent**. Databases avoid this by rolling the transaction back.

**Interview angle:** "A database finds a deadlock between two transactions. What does it do?" Abort one (the victim), roll it back, and usually retry it. Follow-up: *"How does it pick the victim?"* Typically the youngest transaction, or the one that has done the least work.

---

## Q130. [M5 | Avoidance with claim edges | Medium | MCQ | OA]

In the resource-allocation-graph **avoidance** algorithm (every resource has a single instance), each process declares its possible future requests as **claim edges** P ⇢ R in advance. P now requests R, which is free. When is the request granted?

- A. Only if converting the claim edge P ⇢ R into an assignment edge R → P does not create a cycle in the graph (claim edges included)
- B. Whenever R is free
- C. Only if P currently holds no other resources
- D. Only if the graph has no cycle right now, before the conversion

**Answer:** A

**Why:** If the conversion creates a cycle, the system would be in an unsafe state: a future request along a claim edge could complete a circular wait. So the request must wait even though R is free. Cycle detection costs O(n²).

**Why the other options fail:**
- B: R being free is necessary, but not sufficient. That's the difference between avoidance and simply granting requests.
- C: That's a prevention rule (hold-and-wait), not avoidance.
- D: The check must be on the graph **after** the conversion.

**Concept recap:**
- Single instances → RAG + claim edges, checked for cycles.
- Multiple instances → the Banker's algorithm.

**Trap alert:** Claim edges are part of the cycle check. A cycle that runs through claim edges signals an unsafe state, not a deadlock that has already happened.

**Interview angle:** "How do avoidance algorithms differ for single vs multiple instances?" Follow-up: *"What happens to a claim edge when the resource is released?"* The assignment edge R → P turns back into the claim edge P ⇢ R.

---

## Q131. [M5 | Deadlock, starvation, or livelock? | Easy | MCQ | Interview]

Two people meet in a narrow corridor. Each politely steps to the same side to let the other pass, then both step to the other side, and so on forever. Both keep moving, but neither gets by. In concurrency terms, this is:

- A. Deadlock
- B. Starvation
- C. Priority inversion
- D. Livelock

**Answer:** D

**Why:** In a **livelock**, the processes are *active*, changing state in response to each other, but no one makes progress. In a deadlock, everyone is *blocked* and stops completely.

**Why the other options fail:**
- A: Nobody is blocked; both keep moving.
- B: Starvation means *some* party is repeatedly passed over while others make progress. Here, nobody progresses.
- C: Priority inversion involves priorities and a lock holder.

**Concept recap:**
- **Deadlock:** blocked, waiting in a cycle; uses no CPU.
- **Livelock:** running, uses CPU, but no progress.
- **Starvation:** others progress, but one process never gets its turn.

**Trap alert:** Livelock often comes *from* deadlock-avoidance code, such as "if I can't get the second lock, release the first and retry".

**Interview angle:** "Deadlock vs livelock vs starvation, with examples." Follow-up: *"How do you break a livelock?"* Randomized or exponential back-off, so the two sides stop acting in lockstep.

---

## Q132. [M5 | Comparing deadlock, starvation and livelock | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. Deadlocked processes are blocked and use no CPU; livelocked processes keep running but make no progress.
- B. Starvation can occur without any deadlock.
- C. The processes in a deadlock wait forever, so they are also starved.
- D. Livelock can arise from deadlock-avoidance schemes such as repeated "try_lock, fail, release, retry" with identical timing.
- E. Randomized back-off cannot help with livelock.

**Answer:** A, B, C, D

**Why:**
- **A:** Blocked vs busy is the key observable difference.
- **B:** For example, SJF starving a long job (Q60), or reader-preference starving writers (Q108).
- **C:** Waiting forever is the extreme case of being starved.
- **D:** Two threads that keep backing off in lockstep collide again and again.

**Why the other options fail:**
- E: Random back-off breaks the symmetry. This is exactly how Ethernet's exponential back-off resolves repeated collisions.

**Concept recap:**

| | Blocked? | Uses CPU? | Progress? |
|---|---|---|---|
| Deadlock | yes | no | none (for those involved) |
| Livelock | no | yes | none |
| Starvation | maybe | maybe | others progress; the victim doesn't |

**Trap alert:** "No deadlock ⇒ no problem". Livelock and starvation are separate liveness failures and need separate fixes.

**Interview angle:** "Your service's CPU is at 100% but nothing completes. Deadlock?" More likely livelock, or a spin. A deadlock would usually show idle CPU. Follow-up: *"How would you confirm it?"* Take thread dumps a few seconds apart. Deadlocked threads sit blocked in the same place every time, while livelocked or spinning threads are runnable and keep cycling through the same code.

---

## Q133. [M5 | What real operating systems do | Easy | MCQ | OA + Interview]

How do general-purpose OSes such as Linux and Windows typically handle deadlocks among **user applications' locks**?

- A. They run the Banker's algorithm on every lock request
- B. They prevent all four Coffman conditions
- C. They largely ignore the problem (the "ostrich" approach) and leave it to application developers; specialized systems such as databases do detect and recover
- D. They reboot periodically to clear deadlocks

**Answer:** C

**Why:** Deadlocks among user locks are rare. Detection would cost overhead on every lock operation, and avoidance would need knowledge the OS doesn't have. So the OS **ignores** application-level deadlock, and preventing it, mainly through lock ordering, is the programmer's job. Kernels do guard *their own* locks; Linux's `lockdep` in debug builds is one example.

**Why the other options fail:**
- A: Too costly, and processes don't declare their maximum needs.
- B: Impossible in general; mutual exclusion is inherent to many resources.
- D: Not a strategy any OS uses.

**Concept recap:** The four strategies are prevention, avoidance, detection + recovery, and ignoring the problem. General-purpose OSes mostly ignore application deadlocks. Databases use detection + recovery.

**Trap alert:** "The OS will detect my deadlock" is false. A deadlocked user program just hangs.

**Interview angle:** "Does Linux detect deadlocks in user programs?" No, apart from tools: `gdb` thread dumps, `jstack` for Java, and thread sanitizers. Follow-up: *"Why is ignoring them reasonable?"* Deadlocks are rare, any check costs every lock operation something, and the programmer knows the lock structure best.

---

## Q134. [M5 | Minimum instances that rule out deadlock | Hard | MCQ | OA]

Five processes share a single resource type. Each may need **up to 3** instances at the same time, and holds what it gets until it finishes. What is the **minimum** total number of instances that guarantees deadlock can never occur?

- A. 15
- B. 11
- C. 10
- D. 6

**Answer:** B

**Why:** Consider the worst case. Each process holds one less than its maximum (k − 1 = 2) and waits for one more. That uses 5 × 2 = 10 instances, and if there are only 10, everyone waits forever. **One extra instance** lets some process reach its maximum, finish, and release everything, and then the rest follow.

Formula: minimum = **n(k − 1) + 1** = 5 × (3 − 1) + 1 = **11**.

**Why the other options fail:**
- A: n × k = 15 also guarantees no deadlock, but it isn't the *minimum*.
- C: With exactly 10, the worst case above is a deadlock. This option catches forgetting the + 1.
- D: Far too few. With 6, deadlock is possible as soon as three processes each hold 2.

**Concept recap:** For n processes each needing at most k instances: deadlock-free ⇔ total R ≥ n(k − 1) + 1. Solving for n gives the maximum number of processes: n ≤ (R − 1)/(k − 1).

**Trap alert:** The question asks for the **minimum** that guarantees safety, not just *a* safe number. And it's n(k − 1) **+ 1**, not n(k − 1).

**Interview angle:** Usually an OA item. Variant: *"With 6 instances and each process needing up to 2, how many processes can run without deadlock?"* n(2 − 1) + 1 ≤ 6 → **n ≤ 5**.

---

## Q135. [M6 | Memory hierarchy | Easy | MCQ | OA]

Which order is from **fastest to slowest** typical access time?

- A. L1 cache, registers, DRAM, SSD
- B. DRAM, L1 cache, registers, HDD
- C. SSD, DRAM, L2 cache, registers
- D. Registers, L1 cache, L2/L3 cache, DRAM, SSD, HDD

**Answer:** D

**Why:** Storage closer to the CPU is smaller, faster, and more expensive per byte. Rough orders of magnitude:
- Registers: below 1 ns
- L1 cache: about 1 ns
- L2/L3 caches: a few to tens of ns
- DRAM: about 100 ns
- SSD: tens of µs
- HDD: milliseconds (dominated by seek and rotation)

**Why the other options fail:**
- A: Registers are faster than L1.
- B: DRAM is not faster than the caches.
- C: Completely reversed.

**Concept recap:** The hierarchy works because of **locality**: programs reuse recent data (temporal locality) and access nearby data (spatial locality). So a small fast level catches most accesses.

**Trap alert:** The gap between DRAM and disk is huge, roughly 1,000× for SSD and 100,000× for HDD. That's why even a small page-fault rate dominates memory access time (covered in M7).

**Interview angle:** "Why do we need a memory hierarchy?" Speed, size and cost pull in opposite directions, and locality makes a hierarchy effective. Follow-up: *"Who manages each level?"* Hardware manages the caches, the OS manages RAM ↔ disk (paging), and the compiler manages registers.

---

## Q136. [M6 | Logical vs physical addresses | Easy | MCQ | Interview]

Which statement is correct?

- A. The CPU generates logical (virtual) addresses; the MMU translates each one to a physical address, and a user program never sees physical addresses
- B. The compiler generates physical addresses, and the CPU uses them directly
- C. Logical addresses are always larger than physical addresses
- D. Only the kernel uses logical addresses

**Answer:** A

**Why:**
- A **logical address** is what the running program uses: what's in its pointers and in the program counter.
- A **physical address** is the actual location in RAM.
- The **MMU** (Memory Management Unit) is hardware that translates every logical address to a physical one on every access.

This lets each process have its own private address space, and lets the OS place it anywhere in RAM.

**Why the other options fail:**
- B: True only for compile-time binding (Q137), which modern general-purpose systems don't use.
- C: Address-space sizes vary. A virtual space can be larger or smaller than physical memory.
- D: User programs use logical addresses all the time.

**Concept recap:** Logical address space = the set of addresses a program can generate. Physical address space = the actual RAM locations. The MMU maps one to the other.

**Trap alert:** Printing a pointer in C shows a **virtual** address. Two processes can print the same pointer value for completely different physical memory.

**Interview angle:** "What happens to an address between `*p = 5` in your code and the RAM chip?" Virtual address → TLB/MMU → physical address → cache/RAM. Follow-up: *"Why not let programs use physical addresses?"* There would be no isolation between processes, no relocation, and no virtual memory.

---

## Q137. [M6 | Address-binding times | Medium | MCQ | OA + Interview]

A process must be movable to a different area of memory **while it is running**, for example after compaction or when it is swapped back in at a different location. Which kind of address binding makes this possible?

- A. Compile-time binding
- B. Load-time binding
- C. Execution-time (run-time) binding, which needs hardware support such as an MMU
- D. Link-time binding

**Answer:** C

**Why:** **Binding** fixes the mapping from a program's addresses to physical memory. If the mapping is fixed *before* the program runs, moving the program would break every address stored inside it. With execution-time binding, the translation happens **on every access**, so the OS can move the process and simply update the mapping (for example, a new base-register value).

**Why the other options fail:**
- A: Compile time produces **absolute code**: the program must run at exactly that address.
- B: Load time produces **relocatable code** that is adjusted when loaded, but it's fixed after that.
- D: Linking resolves references between modules. It doesn't let a running process move.

**Concept recap:**

| Binding time | Code produced | Can move while running? |
|---|---|---|
| Compile | absolute | no |
| Load | relocatable, fixed at load | no |
| Execution | translated per access by the MMU | **yes** |

**Trap alert:** Only execution-time binding needs **hardware support**. The other two are purely software steps.

**Interview angle:** "When are addresses bound?" Follow-up: *"Which binding does Linux use for user programs?"* Execution-time, through paging and the MMU.

---

## Q138. [M6 | Base and limit registers | Easy | Numerical | OA]

A process runs with **base (relocation) register = 14000** and **limit register = 3000**. It issues logical addresses **2500** and **3200**. What happens?

- A. They map to physical addresses 16500 and 17200
- B. They are used as-is: physical 2500 and 3200
- C. 2500 traps, because 2500 < 14000
- D. 2500 maps to physical 16500; 3200 causes a trap (addressing error), because 3200 ≥ limit (3000)

**Answer:** D

**Why:** On every access, the MMU first checks **logical < limit**. If that holds, it adds the base. Otherwise it traps to the OS, which normally terminates the process.

**Working:**
- **Rule:** if (logical < limit) then physical = base + logical; else trap.
- 2500 < 3000 ✓ → physical = 14000 + 2500 = **16500**.
- 3200 < 3000? ✗ → **trap** (addressing error / segmentation fault).
- The legal physical range is [14000, 16999].

**Why the other options fail:**
- A: Ignores the limit check. 17200 would land in another process's memory.
- B: Ignores relocation altogether.
- C: The logical address is compared with the **limit**, not with the base.

**Concept recap:** Valid logical addresses are 0 ≤ L < limit, and physical = base + L. One comparison and one addition happen per access, in hardware.

**Trap alert:** Some textbooks show the check on the *physical* address (base ≤ P < base + limit). The result is the same. Just don't compare the logical address with the base.

**Interview angle:** "How does the simplest MMU protect processes from each other?" With base + limit registers. Follow-up: *"Who loads those registers?"* The kernel, during a context switch, using privileged instructions (Q139).

---

## Q139. [M6 | MMU with base and limit | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. Loading the base and limit registers is a privileged operation.
- B. The check is "logical address < limit", done before (or alongside) adding the base.
- C. The base-limit scheme lets the OS relocate a process by moving its memory and changing the base register.
- D. The base-limit scheme lets one process's memory be split into several non-contiguous pieces.
- E. Every memory access made by the user process is checked in hardware.

**Answer:** A, B, C, E

**Why:**
- **A:** If user code could change base or limit, it could reach any memory, so loading them is privileged (Q3).
- **B, E:** The check happens in hardware, on every access, at full speed.
- **C:** Relocation is the main benefit: copy the memory and update the base.

**Why the other options fail:**
- D: One base and one limit describe **one contiguous** block. Non-contiguous allocation needs segmentation or paging (Batch 8).

**Concept recap:** Base + limit gives **contiguous allocation** with relocation and protection. Its weakness is external fragmentation (Batch 8).

**Trap alert:** The kernel runs without these checks (or with its own mapping). They constrain *user* mode.

**Interview angle:** "What does an MMU do?" Translation plus protection, in hardware, on every access. Follow-up: *"What's the limitation of base and limit?"* Contiguity, which leads to fragmentation, and no sharing of parts of a process.

---

## Q140. [M6 | When logical = physical | Hard | MCQ | OA]

Which statement about logical and physical addresses is correct?

- A. Under compile-time and load-time binding, logical and physical addresses are the same; under execution-time binding they differ
- B. They are always the same
- C. They are always different
- D. They differ only under compile-time binding

**Answer:** A

**Why:** With compile-time or load-time binding, the final physical addresses are written **into the program itself** before it runs. The CPU then emits addresses that are already physical, and nothing translates them. With execution-time binding, the CPU emits logical (virtual) addresses, and the MMU translates them at run time, so the two differ.

**Why the other options fail:**
- B: False under execution-time binding, which is what paging systems use.
- C: False under compile-time and load-time binding.
- D: Backwards. Compile-time binding is the case where they're *identical*.

**Concept recap:** logical = physical ⇔ binding happens before execution. logical ≠ physical ⇔ execution-time binding, the virtual-address case.

**Trap alert:** "Logical address" and "virtual address" are used interchangeably **only** under execution-time binding.

**Interview angle:** Rarely asked verbally. Follow-up: *"What does an embedded system with no MMU do?"* It typically uses physical addresses directly (compile- or link-time binding), so there is no isolation between tasks.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 121 | A | 126 | D | 131 | D | 136 | A |
| 122 | A, C, D | 127 | B | 132 | A, B, C, D | 137 | C |
| 123 | B | 128 | C | 133 | C | 138 | D |
| 124 | C | 129 | A, B, C, D | 134 | B | 139 | A, B, C, E |
| 125 | B | 130 | A | 135 | D | 140 | A |

---

**Concepts tested so far (to avoid repeats):** Batches 1–6, plus: breaking hold-and-wait · prevention for each condition, prevention vs avoidance · unsafe ≠ deadlocked · single-type safe sequences (not unique) · Banker's safety with three resource types · request algorithm (available but refused) · detection uses Request · detection run with a partially deadlocked set · recovery (termination, preemption + rollback, victim starvation) · claim-edge avoidance · livelock vs deadlock vs starvation · the ostrich approach and `lockdep` · n(k − 1) + 1 · hierarchy latencies and locality · logical vs physical and the MMU · binding times · base/limit translation and trap · base/limit privilege and contiguity · logical = physical only for compile/load binding
*(M5 is complete. M6 continues in Batch 8 with fixed/variable partitioning, first/best/worst/next fit, internal vs external fragmentation, compaction, paging and address translation, page-table sizing, and segmentation.)*

**Completed: 140 / 300    Remaining: 160**
