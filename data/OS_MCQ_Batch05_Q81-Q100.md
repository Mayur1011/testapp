# OS MCQ/MSQ Bank — Batch 5 (Q81–Q100)

**What this batch covers (M4, synchronization):** Peterson's algorithm (proof idea, ordering, what breaks it), test-and-set (TAS), bounded-waiting TAS and ticket locks, compare-and-swap (CAS) and lock-free updates, the ABA problem, mutexes (ownership, relocking from the same thread, futex fast path), spinlocks vs blocking, busy-waiting hazards, semaphores (blocking implementation, binary vs counting, value arithmetic, ordering traces, common misuse).

> Conventions for this batch:
> - Memory is sequentially consistent unless a question says otherwise.
> - `wait()` is also called P/down, and `signal()` is also called V/up.
> - Semaphores use the **blocking** implementation (shown in Q93) unless stated.

**Progress report #1 follows the answer key at the end of this file.**

---

## Q81. [M4 | Peterson's algorithm: why the order matters | Hard | MCQ | OA + Interview]

```c
// shared:  bool flag[2] = {false, false};  int turn;
// Process Pi   (j = 1 - i)
flag[i] = true;                      // line 1: "I want to enter"
turn = j;                            // line 2: "but you go first"
while (flag[j] && turn == j) ;       // busy-wait
/* critical section */
flag[i] = false;
```

Suppose lines 1 and 2 are **swapped**, so each process writes `turn = j` first and then `flag[i] = true`. Which property can now fail?

- A. Only bounded waiting
- B. Mutual exclusion
- C. Progress only
- D. Nothing; the algorithm is still correct

**Answer:** B

**Why:** Peterson's correctness depends on announcing interest (`flag[i] = true`) *before* giving way (`turn = j`). With the order swapped, a process can give way and then get preempted before raising its flag. The other process then sees "not interested" and walks in. Once the first process resumes, it finds `turn` pointing at itself and walks in too.

**Working (a counterexample interleaving):**
1. P0: `turn = 1` → preempted before setting its flag.
2. P1: `turn = 0`; `flag[1] = true`.
3. P1 tests `flag[0] && turn == 0` → `flag[0]` is still false → **P1 enters the CS**.
4. P0 resumes: `flag[0] = true`.
5. P0 tests `flag[1] && turn == 1` → `turn` is 0 → **P0 enters the CS**.
6. Both processes are in the critical section → mutual exclusion is violated.

**Why the other options fail:**
- A, C: The swap breaks the **safety** property itself, not just fairness or liveness.
- D: The interleaving above is a concrete failure.

**Concept recap:** Peterson's algorithm = flag (intent) + turn (politeness). The **process that writes `turn` last waits.** The intent must be visible before the turn write.

**Trap alert:** "The two writes are independent, so order doesn't matter" is the misconception. In concurrent code, the **order in which other threads see your writes** is part of the algorithm's correctness.

**Interview angle:** "Explain Peterson's algorithm and prove mutual exclusion." Proof sketch: if both were inside, then `flag[0]` and `flag[1]` are both true, so each one got past the loop only because `turn` pointed at the other. But `turn` holds a single value, so that's impossible. Follow-up: *"Does it work on a real x86 machine?"* See Q82.

---

## Q82. [M4 | Peterson's algorithm: properties | Medium | MSQ | Interview]

**Select all that apply.**

- A. For two processes, it satisfies mutual exclusion, progress and bounded waiting.
- B. It uses busy waiting.
- C. On modern multicore CPUs, plain Peterson can fail unless memory barriers (fences) are added, because the hardware may reorder a store with a later load.
- D. As written, it works for any number of processes.
- E. It requires special atomic hardware instructions such as test-and-set.

**Answer:** A, B, C

**Why:**
- **A:** Progress holds because a process that isn't interested has its flag false, so it never blocks anyone. Bounded waiting holds because, after you set `turn = j`, the other process can enter at most once before you do.
- **B:** The `while` loop spins on the CPU, which is busy waiting.
- **C:** x86 lets a later load complete before an earlier store becomes visible to other cores (the store-buffer effect). Each process can then read the other's flag as false, and both enter. A **fence** instruction (e.g., `mfence` on x86) prevents that reordering.

**Why the other options fail:**
- D: It's a 2-process algorithm. For n processes you need a generalization, such as the filter lock or Lamport's bakery algorithm.
- E: It needs only ordinary atomic loads and stores, which is exactly why it's taught.

**Concept recap:** Peterson = a software-only 2-process solution with all three properties, under sequential consistency.

**Trap alert:** The proof assumes **sequentially consistent** memory: every thread sees all writes in one single global order. Real hardware and compilers don't guarantee that by default.

**Interview angle:** "Why don't we use Peterson's algorithm in practice?" It handles only two processes, it busy-waits, and it needs fences on modern hardware. Hardware atomics such as TAS and CAS are simpler and scale. Follow-up: *"What does a memory fence do?"* It forces memory operations before the fence to become visible before those after it.

---

## Q83. [M4 | Peterson's algorithm: who enters? | Medium | MCQ | OA]

Use the correct Peterson code from Q81. Both processes want to enter, and their steps interleave in exactly this order:
1. P0: `flag[0] = true`
2. P1: `flag[1] = true`
3. P0: `turn = 1`
4. P1: `turn = 0`
5. Both now evaluate their `while` conditions.

What happens?

- A. P0 enters the critical section; P1 spins
- B. P1 enters the critical section; P0 spins
- C. Both enter
- D. Neither enters; they deadlock

**Answer:** A

**Why:** P1 wrote `turn` last, so the final value is `turn = 0`.
- P0 checks `flag[1] && turn == 1`: true && false = **false** → it leaves the loop and **enters**.
- P1 checks `flag[0] && turn == 0`: true && true = **true** → it **spins**.

When P0 exits and sets `flag[0] = false`, P1's condition becomes false and P1 enters.

**Why the other options fail:**
- B: Assumes whoever wrote `turn` last "wins". It's the opposite: the last writer gave way last.
- C: Impossible, because `turn` can favour only one process.
- D: `turn` always lets exactly one process through, so they can't both wait.

**Concept recap:** Whoever writes `turn` **last** waits. `turn` breaks the tie when both flags are up.

**Trap alert:** Misreading `turn = j` as "it's my turn". It means "**your** turn", a polite offer.

**Interview angle:** "In Peterson's algorithm, which process enters when both try at once?" Follow-up: *"Why can't both spin forever?"* `turn` is either 0 or 1, so one of the two loop conditions must be false.

---

## Q84. [M4 | Test-and-set: why atomicity matters | Medium | MCQ | OA + Interview]

```c
bool TestAndSet(bool *target) {   // the hardware executes this ATOMICALLY
    bool rv = *target;
    *target = true;
    return rv;
}
// shared: bool lock = false;
while (TestAndSet(&lock)) ;   // spin until we saw "false"
/* critical section */
lock = false;
```

Suppose that, instead of one atomic instruction, `TestAndSet` were just these two ordinary statements, which could be interleaved between threads. What would go wrong?

- A. Nothing; the lock would still work
- B. The lock would be slower but still correct
- C. Only bounded waiting would be lost
- D. Two threads could both read `false` before either writes `true`, so both would enter the critical section, losing mutual exclusion

**Answer:** D

**Why:** The whole point of TAS is that **"read the old value and set it to true" happens as one indivisible step**. No other CPU can sneak in between. Split into two steps, it has the same read-then-write race as `count++` in Q77:
1. T1 reads `false`.
2. T2 reads `false`.
3. Both set `true`, and both see "it was free" → both enter.

**Why the other options fail:**
- A, B: Correctness itself is gone, not just speed.
- C: The lost property is mutual exclusion, which is far worse than losing bounded waiting.

**Concept recap:** TAS = an atomic read-modify-write. Real instructions: `xchg` on x86 (with an implicit bus lock), or `lock bts`. On ARM, an LL/SC pair (`ldrex`/`strex`, or `ldxr`/`stxr`).

**Trap alert:** "One line of C" ≠ "one atomic step". Atomicity is a property of the hardware instruction (or a library function such as `atomic_exchange`), not of how the source looks.

**Interview angle:** "How do you build a lock from hardware primitives?" Write the TAS spin loop. Follow-up: *"What's the problem with spinning on TAS itself?"* Every TAS is a write, so the cache line bounces between cores. The fix is "test-and-test-and-set": spin on a plain read, and attempt TAS only when the lock looks free.

---

## Q85. [M4 | Bounded-waiting TAS and ticket locks | Hard | MSQ | OA + Interview]

This bounded-waiting lock for n processes uses TAS (Silberschatz):

```c
// shared: bool waiting[n] = {false};  bool lock = false;
// Process Pi:
waiting[i] = true;  bool key = true;
while (waiting[i] && key)
    key = TestAndSet(&lock);
waiting[i] = false;
/* critical section */
int j = (i + 1) % n;
while (j != i && !waiting[j])
    j = (j + 1) % n;
if (j == i) lock = false;        // nobody waiting: release
else        waiting[j] = false;  // hand the CS straight to Pj
```

**Select all that apply.**

- A. A process that wants to enter waits for at most n − 1 other critical-section entries.
- B. On exit, if another process is waiting, the CS is handed directly to the next waiting process in cyclic order, without setting `lock = false`.
- C. A plain TAS spinlock (`while (TestAndSet(&lock));`) can let one unlucky process lose the race every time.
- D. To hand over, the exiting process must first set `lock = false`, so the next process can win the TAS.
- E. A ticket lock built on atomic fetch-and-add serves waiters in FIFO order.

**Answer:** A, B, C, E

**Why:**
- **A:** The exit code scans i+1, i+2, … in cyclic order, so every waiting process is reached within n − 1 hand-offs.
- **B:** Setting `waiting[j] = false` releases Pj from its loop while `lock` stays true. Ownership passes directly, and no third process can race in.
- **C:** With plain TAS, *whoever happens to execute TAS first after the release* wins. Nothing prevents the same process from losing every time.
- **E:** In a ticket lock, `my = fetch_add(&next, 1)` takes a ticket, and `while (serving != my);` waits for it to be called. Tickets are served strictly in order.

**Why the other options fail:**
- D: Setting `lock = false` would let *any* process (including a newcomer) win the TAS, which defeats the guaranteed cyclic order.

**Concept recap:** Bounded waiting needs an **ordering mechanism**, such as cyclic hand-off or tickets. Plain TAS gives only ME + progress.

**Trap alert:** "Atomic instruction ⇒ fair lock" is false. Atomicity guarantees safety, not fairness.

**Interview angle:** "How would you make a spinlock fair?" Use a ticket lock. Follow-up: *"Downside of ticket locks under heavy contention?"* Every waiter spins on the same `serving` variable, so each release invalidates every waiter's cached copy. Queue locks (MCS) let each waiter spin on its own variable.

---

## Q86. [M4 | CAS-based lock-free increment | Hard | Code-trace | OA + Interview]

```c
// executed ATOMICALLY by the hardware
int compare_and_swap(int *value, int expected, int new_value) {
    int old = *value;
    if (old == expected) *value = new_value;
    return old;
}

void increment(int *v) {
    int old;
    do {
        old = *v;                                       // read
    } while (compare_and_swap(v, old, old + 1) != old); // retry if changed
}
```

`v` starts at 0. Two threads each call `increment(&v)` **3 times**, concurrently. What is true about the final value of `v`?

- A. It can be any value from 2 to 6
- B. It can be any value from 3 to 6
- C. It is exactly 6; a CAS fails only when another thread's CAS succeeded in between, so updates are never lost (a thread may just retry)
- D. It is exactly 6, but the program may loop forever because both threads can keep failing

**Answer:** C

**Why:** CAS writes only if the value **still equals** what you read. If another thread changed it in the meantime, the CAS fails, writes nothing, and you retry with the fresh value. So no update can overwrite a newer one. That is exactly the lost update from Q77, and it can't happen here.

**Working (one conflict):**

| Step | T1 | T2 | v |
|---|---|---|---|
| 1 | old = 0 | | 0 |
| 2 | | old = 0 | 0 |
| 3 | | CAS(v, 0, 1) → succeeds | 1 |
| 4 | CAS(v, 0, 1) → v is 1 ≠ 0 → **fails**, returns 1 | | 1 |
| 5 | retry: old = 1 | | 1 |
| 6 | CAS(v, 1, 2) → succeeds | | 2 |

- Each *successful* CAS adds exactly 1, and each of the 6 calls ends with exactly one successful CAS → final value = **6**.
- **Why it terminates:** a CAS fails only because some *other* CAS succeeded in between. There are only 6 successes in total, so the number of failures is finite too.

**Why the other options fail:**
- A: The plain `count++` result from Q77. CAS removes the lost update.
- B: Also assumes some updates can be lost.
- D: Unbounded failure would require unboundedly many successes, but there are only 6.

**Concept recap:** The CAS loop is **read → compute → CAS → retry on failure**. It is **lock-free**: some thread always makes progress, though a particular thread may retry many times. **Wait-free** is stronger: every thread finishes in a bounded number of its own steps.

**Trap alert:** "No locks ⇒ no waiting" is false. Under heavy contention, threads can retry many times, and throughput may be worse than with a good lock.

**Interview angle:** "Implement an atomic counter without a mutex." Write the CAS loop, or just use `std::atomic<int>::fetch_add`. Follow-up: *"What is the ABA problem?"* See Q87.

---

## Q87. [M4 | The ABA problem | Hard | MCQ | Interview]

A lock-free stack pops like this:

```c
do {
    old  = top;
    next = old->next;
} while (!CAS(&top, old, next));
```

T1 reads `top = A` and `next = B`, then is preempted. T2 pops A, pops B (B's memory is freed), and pushes A again (the same node address is reused). Now `top = A` and `A->next = C`. T1 resumes: its `CAS(&top, A, B)` **succeeds**, so `top` now points to the freed node B. What is this bug called?

- A. The ABA problem
- B. Priority inversion
- C. The convoy effect
- D. False sharing

**Answer:** A

**Why:** CAS checks only that the *value* is unchanged ("is top still A?"). It can't tell that the value went A → B → A in between, while the structure behind it changed. T1's saved `next = B` is stale, and the CAS accepts it.

**Why the other options fail:**
- B: A high-priority thread blocked behind a low-priority lock holder.
- C: Short jobs stuck behind a long one under FCFS (Q69).
- D: Unrelated variables in one cache line causing coherence traffic.

**Concept recap:** Standard fixes for ABA:
1. A **tagged pointer**: a version counter next to the pointer, updated together with a double-width CAS.
2. **Safe memory reclamation** (hazard pointers, epoch-based reclamation), so a node isn't reused while another thread may still hold it.
3. LL/SC hardware, which fails on *any* intervening write, not just a changed value.

**Trap alert:** "If CAS succeeded, nothing changed" is false. CAS only proves the value is the same *right now*.

**Interview angle:** "What problems do lock-free data structures face?" ABA, memory reclamation, and the difficulty of proving them correct. Follow-up: *"How does a tagged pointer fix ABA?"* The counter increments on every change, so A-with-version-5 ≠ A-with-version-7.

---

## Q88. [M4 | Mutex basics | Easy | MCQ | Interview]

Which correctly describes a (blocking) mutex's `acquire()` and `release()`?

- A. `acquire()` spins forever, burning CPU, until the lock is free
- B. If the lock is held, `acquire()` puts the caller to sleep; `release()` wakes a waiter. Only the thread that acquired the mutex should release it
- C. Any thread may release a mutex, even one that never acquired it
- D. A mutex lets up to N threads into the critical section at once

**Answer:** B

**Why:** A **mutex** (MUTual EXclusion lock) protects a critical section so that only one thread is inside at a time. A *blocking* mutex sleeps instead of spinning, so waiting threads use no CPU. A mutex also has an **owner**: the thread that locked it is the one that must unlock it.

**Why the other options fail:**
- A: That's a spinlock.
- C: Breaks ownership. An error-checking pthread mutex returns `EPERM` if a non-owner unlocks it, and for other types it's undefined behaviour.
- D: That's a counting semaphore initialized to N.

**Concept recap:** Mutex = lock + owner + (usually) sleep while waiting. Binary semaphore = a value that is 0 or 1, with **no owner**, so any thread may signal it.

**Trap alert:** "Mutex = binary semaphore" is wrong because of **ownership**. Ownership is what enables relock detection (Q89) and priority inheritance.

**Interview angle:** "Mutex vs binary semaphore?" Explain ownership and its consequences. Follow-up: *"When would you use a semaphore where a mutex won't do?"* For signalling between threads, e.g., "thread A finished step 1, so thread B may start" (Q94 C).

---

## Q89. [M4 | Relocking a mutex from the same thread | Medium | MCQ | OA + Interview]

```c
pthread_mutexattr_t attr;
pthread_mutexattr_init(&attr);
pthread_mutexattr_settype(&attr, PTHREAD_MUTEX_NORMAL);
pthread_mutex_init(&m, &attr);

pthread_mutex_lock(&m);
pthread_mutex_lock(&m);   // same thread locks it again
```

What happens at the second `lock`?

- A. It succeeds, and the lock count becomes 2
- B. It returns `EDEADLK` immediately
- C. The mutex is automatically released and then re-acquired
- D. The thread blocks forever waiting for itself (self-deadlock)

**Answer:** D

**Why:** A `NORMAL` mutex doesn't track ownership for relocking. To it, the second call is just another lock attempt on a held mutex, so the caller sleeps until the holder unlocks. The holder is the caller itself, so it never wakes.

**Why the other options fail:**
- A: That's the **`PTHREAD_MUTEX_RECURSIVE`** behaviour. Each lock increments a count, and the thread needs as many unlocks.
- B: That's the **`PTHREAD_MUTEX_ERRORCHECK`** behaviour.
- C: No mutex type does this.

**Concept recap:**

| Type | Same-thread relock |
|---|---|
| NORMAL | deadlocks |
| ERRORCHECK | returns `EDEADLK` |
| RECURSIVE | succeeds; count + 1 |
| DEFAULT | undefined (glibc acts like NORMAL) |

**Trap alert:** A function that holds lock `m` calls a helper that also locks `m`. With a non-recursive mutex, that hangs silently. It is the most common real-world self-deadlock.

**Interview angle:** "What happens if a thread locks a mutex it already holds?" Answer by type. Follow-up: *"Why not make every mutex recursive?"* Recursive locks hide design problems, and they make it hard to know when the protected invariants actually hold.

---

## Q90. [M4 | Spin or block? | Medium | MCQ | OA]

Multicore machine. A lock is held for at most **2 µs**, and the holder is running on another core. Blocking (sleeping and later being woken) costs about two context switches, roughly **10 µs** in total. What is the best way for a thread to wait for this lock?

- A. Block immediately with a mutex, because spinning is always wasteful
- B. Disable interrupts
- C. Spin, because the expected wait (≤ 2 µs) is far below the cost of blocking (≈ 10 µs)
- D. Call `sleep(1)` inside a retry loop

**Answer:** C

**Why:** Spinning wastes CPU only for as long as you wait. Blocking wastes the cost of two context switches, however short the wait is. When the holder is *running on another core* and will release soon, spinning gets the lock sooner and wastes less.

**Why the other options fail:**
- A: "Always" is wrong. For very short waits, blocking costs more than it saves.
- B: Disabling interrupts on your own core doesn't stop other cores. It's also privileged, so user code can't do it.
- D: Sleeping for a whole second to wait 2 µs is absurdly slow.

**Concept recap:** **Spin if the expected wait < the cost of 2 context switches AND the lock holder is running on another CPU.** Otherwise, block.

**Trap alert:** On a **single CPU**, spinning is pointless: the holder can't run while you spin, so you waste your whole quantum (Q91 A).

**Interview angle:** "Spinlock vs mutex: when would you use each?" Follow-up: *"What's an adaptive mutex?"* It spins briefly while the holder is running on another CPU and then sleeps. Solaris has these, and glibc offers a similar mutex type.

---

## Q91. [M4 | Spinlock hazards | Hard | MSQ | Interview]

**Select all that apply.**

- A. On a uniprocessor, a thread spinning on a lock held by a preempted thread wastes its whole quantum.
- B. In the Linux kernel, holding a spinlock disables preemption on that CPU.
- C. Code that holds a spinlock must not sleep (block).
- D. Spinlocks avoid context-switch overhead while waiting.
- E. A basic test-and-set spinlock guarantees FIFO fairness among waiters.

**Answer:** A, B, C, D

**Why:**
- **A:** With one CPU, the holder can't make progress while you spin, so every cycle you spend spinning is wasted.
- **B:** If the holder could be preempted, the other threads would spin until it was scheduled again. That's why `spin_lock()` disables preemption.
- **C:** Sleeping while holding a spinlock leaves every other waiter spinning for a long, unbounded time, and can deadlock if the code that would wake you needs the same lock.
- **D:** A spinning waiter never gives up the CPU, so there's no switch.

**Why the other options fail:**
- E: Plain TAS has no ordering, so whoever wins the race gets the lock (Q85 C). FIFO needs a ticket or queued lock. Linux now uses fair queued spinlocks, but a *basic* TAS lock isn't fair.

**Concept recap:** Spinlock rules: keep critical sections very short, never sleep while holding one, and use them only on multiprocessors or with preemption disabled.

**Trap alert:** Interrupt handlers can't sleep, so they can't take a blocking mutex. In the kernel they use spinlocks, and if the lock is shared with process context, that code must take it with `spin_lock_irqsave` to avoid self-deadlock.

**Interview angle:** "Why can't you sleep while holding a spinlock?" Follow-up: *"How can a spinlock shared with an interrupt handler deadlock a single CPU?"* A thread holds the lock, an interrupt arrives on the same CPU, and its handler spins forever on that lock. The holder can't run again until the handler returns, so neither ever proceeds.

---

## Q92. [M4 | futex fast path | Medium | MCQ | Interview]

On Linux, why is an **uncontended** `pthread_mutex_lock()` + `pthread_mutex_unlock()` pair cheap?

- A. The fast path is a single atomic instruction (e.g., CAS) on a word in user memory; the kernel is entered (via the `futex` system call) only when there is contention
- B. The mutex lives in a CPU register
- C. The kernel caches the lock state in a special table
- D. Locking temporarily disables interrupts

**Answer:** A

**Why:** A **futex** ("fast user-space mutex") keeps the lock state in an ordinary integer in user memory.
- **Lock:** CAS the word from 0 (free) to 1 (locked). If that succeeds, the lock is taken with no system call.
- **Contended:** only when the CAS finds the lock taken does the thread call `futex(FUTEX_WAIT)`, and the kernel puts it to sleep.
- **Unlock:** the thread calls `futex(FUTEX_WAKE)` only if someone is waiting.

**Why the other options fail:**
- B: The lock must be in shared memory so every thread can see it.
- C: No such cache exists. The kernel is involved only when threads must sleep.
- D: User programs can't disable interrupts.

**Concept recap:** Uncontended lock = user space only (≈ the cost of one atomic). Contended lock = a system call to sleep or wake. "Make the common case fast."

**Trap alert:** "Every mutex operation is a system call" is outdated. That's why modern mutexes are cheap when uncontended.

**Interview angle:** "How are pthread mutexes implemented on Linux?" Explain the futex fast and slow paths. Follow-up: *"What makes a mutex slow, then?"* Contention: sleeping, waking, cache-line bouncing, and convoys of waiting threads.

---

## Q93. [M4 | What a negative semaphore value means | Easy | MCQ | OA]

The blocking semaphore implementation:

```c
wait(S):   S.value--;  if (S.value < 0)  { add caller to S.list; block(); }
signal(S): S.value++;  if (S.value <= 0) { remove a process P from S.list; wakeup(P); }
```

At some instant, `S.value = −3`. What does that mean?

- A. The semaphore is corrupted
- B. Three more `signal()` calls are allowed
- C. Three resources are free
- D. Three processes are blocked on S

**Answer:** D

**Why:** In this implementation, `wait` decrements *before* deciding whether to block. Each process that found nothing available pushed the value one step below zero. So a negative value counts the waiters: **|value| = number of blocked processes**.

**Why the other options fail:**
- A: A negative value is normal in this design.
- B: `signal()` can always be called. It isn't "allowed" a certain number of times.
- C: Free resources correspond to *positive* values.

**Concept recap:**
- value > 0 → that many units are available.
- value = 0 → none available and no one waiting.
- value < 0 → |value| processes are waiting.

**Trap alert:** In the *textbook busy-wait* definition (`while (S <= 0); S--;`), the value never goes negative. Check which semantics a question uses.

**Interview angle:** "How is a semaphore implemented without busy waiting?" Show the value plus a list of waiting processes. Follow-up: *"Who wakes a blocked process?"* The next `signal()`.

---

## Q94. [M4 | Uses of semaphores | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. A counting semaphore initialized to N can limit concurrent access to N instances of a resource.
- B. A binary semaphore initialized to 1 can provide mutual exclusion.
- C. A semaphore initialized to 0 can enforce an ordering such as "S2 runs only after S1": one thread runs `S1; signal(s);`, and the other runs `wait(s); S2;`.
- D. A binary semaphore has an owner, just as a mutex does.
- E. A `signal()` on a semaphore with no waiters is lost.

**Answer:** A, B, C

**Why:** A semaphore is a counter with blocking. Its three classic uses follow from the initial value:
- **N** → resource pool.
- **1** → mutual exclusion.
- **0** → "wait until someone signals", i.e. ordering and event notification.

**Why the other options fail:**
- D: Semaphores have **no owner**. Any thread may signal, which is exactly what makes the ordering pattern in C possible.
- E: A semaphore **remembers** signals by incrementing its value. A later `wait()` consumes that saved signal immediately, without blocking. (A condition variable's `signal` with no waiters **is** lost; covered in the next batch.)

**Concept recap:** The initial value sets the role: N = pool, 1 = mutex, 0 = signalling.

**Trap alert:** "Semaphore signal and condition-variable signal behave the same" is wrong. The semaphore remembers, the condition variable forgets. This is a common trap.

**Interview angle:** "Show how to make thread B run strictly after thread A using a semaphore." Follow-up: *"Can you do that with a mutex?"* Not correctly, because B can't unlock a mutex that A locked.

---

## Q95. [M4 | Semaphore value arithmetic | Medium | Numerical | OA]

A counting semaphore S is initialized to **7**. Then **20** `wait` (P) operations and **15** `signal` (V) operations are performed, in some order, and all of them eventually complete. What is the final value of S?

- A. 0
- B. 2
- C. −8
- D. 12

**Answer:** B

**Why:** With the blocking implementation, every `wait` subtracts 1 and every `signal` adds 1, whatever the order. The order only decides who blocks along the way, not the final value.

**Working:**
- **Formula:** final value = initial − (# of waits) + (# of signals)
- 7 − 20 + 15 = **2**
- The final value is ≥ 0, so no process is left blocked. (A negative result would mean |value| processes still waiting.)

**Why the other options fail:**
- A: Assumes the value "can't move" past some limit, or confuses 2 with "nothing left".
- C: 7 − 15, i.e. the signal count was subtracted instead of the wait count.
- D: 7 + 20 − 15, i.e. waits and signals swapped.

**Concept recap:** S_final = S₀ − #P + #V.

**Trap alert:** P = wait = decrement; V = signal = increment. The letters come from Dutch (*proberen*, "to test", and *verhogen*, "to increase"). Mixing them up flips the answer.

**Interview angle:** Rarely asked verbally. Follow-up: *"If the final value were −2, what would that tell you?"* Two processes are still blocked on S.

---

## Q96. [M4 | Ordering with semaphores | Hard | Code-trace | OA + Interview]

Semaphores `s = 0` and `t = 0`. Each `print` is atomic. Two threads run concurrently:

```
Thread P:            Thread Q:
  print("A");          wait(s);
  signal(s);           print("X");
  print("B");          signal(t);
  wait(t);             print("Y");
  print("C");
```

Which output is **impossible**?

- A. `A B X Y C`
- B. `A X B C Y`
- C. `A X Y B C`
- D. `A B Y X C`

**Answer:** D

**Why:** Write down the "must come before" rules, then check each option against them.

**Working:**
- **Program order:** A < B < C (thread P) and X < Y (thread Q).
- **Semaphore `s`** (starts at 0): Q's `wait(s)` can't pass until P's `signal(s)`, which comes after A. So **A < X**.
- **Semaphore `t`** (starts at 0): P's `wait(t)` can't pass until Q's `signal(t)`, which comes after X. So **X < C**.

Checking each option:

| Output | A<B<C | X<Y | A<X | X<C | Possible? |
|---|---|---|---|---|---|
| A B X Y C | ✓ | ✓ | ✓ | ✓ | yes |
| A X B C Y | ✓ | ✓ | ✓ | ✓ | yes |
| A X Y B C | ✓ | ✓ | ✓ | ✓ | yes |
| A B Y X C | ✓ | **✗ (Y before X)** | ✓ | ✓ | **no** |

In total, exactly 5 outputs are possible: ABXYC, ABXCY, AXBYC, AXBCY, AXYBC.

**Why the other options fail:** A, B and C satisfy every constraint, so they are possible outputs, not the impossible one.

**Concept recap:** Each `signal(x)` → `wait(x)` pair on a semaphore that starts at 0 creates the edge "everything before the signal < everything after the wait". The possible outputs are exactly the orderings consistent with program order plus these edges.

**Trap alert:** Students check only the semaphore edges and forget **program order** inside a thread. That's how D can look possible at first glance.

**Interview angle:** "Use semaphores to make two threads print alternately (ABAB…)." Use two semaphores, one starting at 1 and the other at 0. Follow-up: *"How would you extend that to three threads printing in turn?"* Use a ring of three semaphores, only the first starting at 1.

---

## Q97. [M4 | Making wait/signal themselves atomic | Medium | MCQ | Interview]

With the blocking semaphore implementation, waiting processes sleep instead of spinning. On a **multiprocessor**, how are the `wait()` and `signal()` operations themselves made atomic?

- A. They don't need to be atomic
- B. By disabling interrupts on every CPU
- C. By protecting the short `wait()`/`signal()` code with a spinlock or atomic instructions. Busy waiting still exists, but only for those few instructions, not for the whole critical section
- D. By running them in user space

**Answer:** C

**Why:** `wait` and `signal` read and update the semaphore's value and its waiting list. They are critical sections themselves. On a multiprocessor, other CPUs could run them at the same moment, so the kernel guards them with a spinlock. The trick is that this spinlock is held for only a handful of instructions, while the *user's* critical section, which may be long, is waited for by sleeping.

**Why the other options fail:**
- A: Two simultaneous `wait()`s could corrupt the value or the list, exactly like Q77.
- B: Disabling interrupts on one CPU doesn't stop the others, and coordinating every CPU's interrupts is slow and impractical.
- D: Blocking and waking are kernel operations.

**Concept recap:**
- On a uniprocessor, disabling interrupts briefly is enough.
- On a multiprocessor, spinlocks or atomics are used.

Either way, busy waiting is limited to the short critical sections *inside the semaphore implementation*.

**Trap alert:** "Blocking semaphores eliminate busy waiting completely" is wrong. They *move* it to a tiny, bounded section.

**Interview angle:** "Doesn't a semaphore implementation itself need a lock?" Yes, a tiny spinlock. Follow-up: *"Why is that acceptable?"* It's held for only a few instructions.

---

## Q98. [M4 | Semaphore misuse | Easy | MCQ | OA]

`mutex` is a semaphore initialized to 1. A programmer writes:

```c
wait(mutex);
/* critical section */
wait(mutex);      // bug: should be signal(mutex)
```

What happens?

- A. The process blocks forever on the second `wait()`, and every other process then blocks at its first `wait()`: deadlock
- B. Mutual exclusion is violated
- C. It works correctly
- D. The second `wait()` increments the semaphore, releasing it

**Answer:** A

**Why:**
1. The first `wait` takes the value 1 → 0.
2. The second `wait` takes 0 → −1, so the process blocks and nobody will ever signal it.
3. Every other process does its first `wait` and blocks too, because the value stays ≤ 0 forever.

**Why the other options fail:**
- B: No one else can get in, so mutual exclusion holds. The system just freezes.
- C: The semaphore is never released.
- D: `wait` always decrements.

**Concept recap:** Common semaphore bugs:
- `signal(); CS; wait();` (order swapped) → **mutual exclusion is violated**, several processes in the CS.
- `wait(); CS; wait();` → **deadlock**.
- Omitting the `wait()` → ME is violated. Omitting the `signal()` → deadlock.

**Trap alert:** Semaphore bugs are easy to write because the compiler can't pair `wait` with `signal`. Monitors and scoped locks (C++ `std::lock_guard`) exist to prevent exactly this (next batch).

**Interview angle:** "What can go wrong with semaphores?" Mention misordering, a forgotten signal, and lock-order deadlocks. Follow-up: *"How does RAII (`lock_guard`) help?"* The destructor always releases the lock, even on an early return or an exception.

---

## Q99. [M4 | Counting semaphore as a resource pool | Easy | MCQ | OA + Interview]

A pool of **5** database connections is guarded by a counting semaphore initialized to 5. Each thread does `wait(sem)`, uses a connection, then `signal(sem)`. Five threads currently hold connections, and a sixth calls `wait(sem)`. What happens?

- A. It gets a connection anyway
- B. It blocks until some thread calls `signal(sem)`
- C. It immediately returns an error
- D. The semaphore resets to 5

**Answer:** B

**Why:** The semaphore counts free connections, 5 → 0, as the first five threads take them. The sixth `wait` finds nothing available and blocks. The next `signal`, when a connection is returned, wakes it.

**Why the other options fail:**
- A: That would overcommit the pool.
- C: A blocking `wait` waits. (Non-blocking variants such as `sem_trywait` return `EAGAIN` instead.)
- D: Semaphores never reset themselves.

**Concept recap:** A counting semaphore initialized to N allows at most N holders at once. A blocked `wait` resumes on the next `signal`.

**Trap alert:** The semaphore only **counts** connections. It doesn't say *which* connection is free. You still need a mutex-protected free list to pick one.

**Interview angle:** "Design a connection pool." Use a semaphore for the count plus a mutex for the list, or a condition variable. Follow-up: *"How do you add a timeout?"* `sem_timedwait()`.

---

## Q100. [M4 | Busy waiting vs blocking | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. On a single-CPU system with strict priority scheduling, a high-priority thread busy-waiting for a lock held by a low-priority thread can spin forever, because the low-priority thread never gets to run.
- B. Blocking (sleeping) instead of spinning adds about two context switches' worth of overhead.
- C. Busy waiting is acceptable for very short critical sections in a multiprocessor kernel.
- D. A blocked thread keeps consuming CPU cycles while it waits.
- E. Disabling interrupts is a valid way for user programs to get mutual exclusion on a multicore machine.

**Answer:** A, B, C

**Why:**
- **A:** The spinning high-priority thread always outranks the lock holder, so the holder never runs, never releases the lock, and the system is stuck. Blocking would let the low-priority thread run.
- **B:** One switch to put the waiter to sleep, and one to bring it back.
- **C:** On a multiprocessor, short waits are cheaper to spin than to block (Q90).

**Why the other options fail:**
- D: A blocked thread is off the ready queue and uses no CPU. That's the point of blocking.
- E: Disabling interrupts is privileged, so user code can't do it, and it affects only the current core anyway.

**Concept recap:**
- Busy waiting: no switch cost, but it burns CPU and is dangerous with strict priorities on one CPU.
- Blocking: costs two switches, but uses no CPU while waiting.

**Trap alert:** Scenario A looks like priority inversion, and it is a related danger. But *classic* priority inversion involves a **third, medium-priority** task preempting the lock holder (covered in M12).

**Interview angle:** "When is busy waiting acceptable?" Short waits, multiprocessors, and kernel code that can't sleep. Follow-up: *"What does priority inheritance fix?"* The lock holder temporarily inherits the waiter's high priority so it can finish and release the lock.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 81 | B | 86 | C | 91 | A, B, C, D | 96 | D |
| 82 | A, B, C | 87 | A | 92 | A | 97 | C |
| 83 | A | 88 | B | 93 | D | 98 | A |
| 84 | D | 89 | D | 94 | A, B, C | 99 | B |
| 85 | A, B, C, E | 90 | C | 95 | B | 100 | A, B, C |

---

**Concepts tested so far (to avoid repeats):** Batches 1–4, plus: Peterson's write order and the proof idea · who enters in Peterson (last `turn` writer waits) · Peterson and memory fences, two-process limit · TAS atomicity · bounded-waiting TAS hand-off, ticket locks · CAS loop, lock-free vs wait-free · ABA and its fixes · mutex ownership · same-thread relock by mutex type · spin-vs-block rule · spinlock rules (no sleeping, preemption off, unfair basic TAS) · futex fast path · negative semaphore value = number of waiters · semaphore roles (N / 1 / 0) and that semaphores remember signals · S₀ − P + V · ordering trace with two semaphores · spinlock inside wait/signal · wait/wait misuse → deadlock · semaphore as a counter for a pool · busy-wait danger with strict priorities
*(Still to come in M4, Batch 6: monitors, condition variables (including lost signals and spurious wakeups), producer-consumer, readers-writers, dining philosophers, barriers, lock ordering, reentrant locks.)*

**Completed: 100 / 300    Remaining: 200**

---
---

# 📊 Progress Report #1 (after Q100)

## 1. Topics completed

| Module | Status | Questions |
|---|---|---|
| M1 OS fundamentals & architecture | ✅ complete | Q1–18 |
| M2 Processes, threads, IPC | ✅ complete | Q19–48 |
| M3 CPU scheduling | ✅ complete | Q49–76 |
| M4 Synchronization & concurrency | ⏳ 24 of 38 done | Q77–100 (continues Q101–114) |

## 2. Key concepts covered

- **Hardware and kernel boundary:** mode bit, privileged instructions, the trap as the only gate into the kernel, interrupts vs traps vs faults, fault vs trap return address, top/bottom halves, interrupt vector table, mode switch ≠ context switch.
- **Kernel design:** monolithic vs microkernel (isolation vs IPC cost), loadable modules, layering, SMP vs AMP, boot chain.
- **Processes:** PCB, the five states, queues, `fork`/`exec`/`wait` semantics, copy-on-write, zombies vs orphans, what a child inherits, stdio buffers duplicated by `fork`.
- **Threads:** what is shared vs private, ULT vs KLT, M:1 / 1:1 / M:N, NPTL 1:1, `fork` in multithreaded programs, concurrency vs parallelism.
- **IPC:** shared memory vs message passing, pipes and FIFOs, message queues, signals (dispositions, async-signal safety), sockets as fds, `PIPE_BUF` atomicity.
- **Scheduling:** the three schedulers, the dispatcher, preemptive decision points, the five criteria, FCFS / SJF / SRTF / priority / RR / MLQ / MLFQ, aging, the convoy effect, quantum choice, affinity, load balancing, NUMA/SMT.
- **Synchronization so far:** race conditions, the three critical-section requirements, Peterson, TAS, ticket locks, CAS and lock-freedom, ABA, mutex types and futexes, spin vs block, semaphores (implementation, roles, tracing, misuse).

## 3. Misconceptions tested so far (the traps)

| # | Misconception | Where |
|---|---|---|
| 1 | The trap instruction is privileged | Q3 |
| 2 | Entering the kernel = a context switch | Q10 |
| 3 | A page fault skips the faulting instruction | Q7 |
| 4 | Every `printf` is one system call; output is the same whether piped or not | Q18, Q32 |
| 5 | `fork` returns 0 to the parent | Q22 |
| 6 | Children created = 2ⁿ (forgetting the −1) | Q24 |
| 7 | A child's write to a global is seen by the parent | Q26 |
| 8 | `exec` returns after the new program runs | Q27 |
| 9 | `kill -9` removes a zombie | Q31 |
| 10 | Raw `status` from `wait` is the exit code (it's 256 for `exit(1)`) | Q28 |
| 11 | ULTs run in parallel on multicore (M:1) | Q35 |
| 12 | A blocked process stays in the ready queue | Q23 |
| 13 | An empty pipe returns 0 (EOF) | Q42 |
| 14 | TCP preserves message boundaries | Q47 |
| 15 | `PIPE_BUF` = pipe capacity | Q48 |
| 16 | Waiting time includes I/O time | Q54 |
| 17 | FCFS average wait doesn't depend on order | Q55 |
| 18 | You can schedule a process before it arrives | Q56 |
| 19 | SJF minimizes the *maximum* wait | Q60 |
| 20 | SRTF cannot starve | Q62 |
| 21 | Round Robin can starve | Q68 |
| 22 | A smaller quantum always helps | Q66 |
| 23 | MLFQ = MLQ, or MLFQ needs known bursts | Q70, Q72 |
| 24 | `count++` is atomic | Q77 |
| 25 | Strict alternation is a correct solution | Q79 |
| 26 | ME + progress ⇒ bounded waiting | Q80 |
| 27 | The order of the two writes in Peterson doesn't matter | Q81 |
| 28 | Atomic instruction ⇒ fair lock | Q85 |
| 29 | CAS success ⇒ nothing changed (ABA) | Q87 |
| 30 | Mutex = binary semaphore | Q88 |
| 31 | Relocking your own mutex is harmless | Q89 |
| 32 | Spinning is always wasteful, or always fine | Q90, Q100 |
| 33 | A semaphore signal with no waiters is lost | Q94 |

## 4. Formulas introduced

| Formula | Where |
|---|---|
| Amdahl: Speedup = 1 / ((1 − p) + p/N) | Q16 |
| Full-buffered writes = ⌈N / B⌉ | Q18 |
| Sequential forks: total = 2ⁿ, children = 2ⁿ − 1 | Q24 |
| TAT = CT − AT; WT = TAT − BT; RT = first run − AT | Q54 |
| CPU utilization = Σ BT / makespan | Q56 |
| Exponential average: τₙ₊₁ = α·tₙ + (1 − α)·τₙ | Q58–59 |
| Aging time = levels × interval | Q64 |
| RR: max wait between turns ≤ (n − 1)·q | Q68 |
| RR utilization = useful / (useful + switches) per cycle | Q67 |
| Lost-update range for 2 threads × n increments = [2, 2n] | Q77 |
| Spin if expected wait < 2 × context-switch cost (holder on another CPU) | Q90 |
| Semaphore: S_final = S₀ − #P + #V; if S < 0, |S| = number of waiters | Q93, Q95 |

## 5. Areas to revise before moving on

Concepts listed in order of how often students slip on them in OAs:

1. **Scheduling traces with conventions:** RR (Q65) and SRTF (Q61). Redo both from a blank page, writing the ready queue after every event.
2. **Fork counting with short-circuit operators** (Q25), and the stdio-buffer duplication pair (Q18 + Q32).
3. **The three critical-section requirements**, especially *progress* (Q79–80). Be able to say which one a given broken protocol violates.
4. **Peterson's algorithm:** recite it, explain the order of the writes, and give the one-paragraph proof (Q81–83).
5. **Semaphore ordering:** derive the "must come before" edges and check program order (Q96).
6. **Interview-ready pairs:** mode switch vs context switch, zombie vs orphan, ULT vs KLT, mutex vs binary semaphore, spinlock vs mutex.
