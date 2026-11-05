# OS MCQ/MSQ Bank — Batch 6 (Q101–Q120)

**What this batch covers:**
- **M4, synchronization completed (Q101–114):** monitors, condition variables (atomic release, `while` vs `if`, lost signals, Hoare vs Mesa), producer-consumer with semaphores, readers-writers, dining philosophers, barriers, lock ordering, reentrant locks.
- **M5, deadlocks begins (Q115–120):** the four necessary conditions, resource-allocation graphs (RAG), wait-for graphs, deadlocks that don't involve locks.

> Conventions:
> - Semaphores use the blocking implementation from Q93.
> - Condition variables follow **pthreads/Java (Mesa) semantics** unless stated.
> - `lock`/`unlock` mean `pthread_mutex_lock`/`unlock`.

---

## Q101. [M4 | Monitors | Easy | MCQ | Interview]

What does a **monitor** guarantee?

- A. Every procedure in it runs in kernel mode
- B. Threads that use it never have to wait
- C. At most one thread is active inside the monitor's procedures at any time; the language or runtime provides this mutual exclusion automatically, and condition variables let threads wait inside it
- D. Deadlock is impossible

**Answer:** C

**Why:** A **monitor** packages shared data together with the procedures that operate on it. Calling any of those procedures acquires the monitor's hidden lock, and returning releases it, so the programmer can't forget to `wait` or `signal` (compare Q98). When a thread must wait for something *inside* the monitor, it uses a **condition variable**, which releases the monitor while the thread sleeps.

**Why the other options fail:**
- A: Monitors are a language-level construct. They have nothing to do with CPU privilege.
- B: Threads still wait, both to enter the monitor and on condition variables.
- D: Two monitors that call each other can still deadlock.

**Concept recap:** Monitor = shared data + procedures + an implicit mutex + condition variables. Examples: Java `synchronized` with `wait`/`notify`, and C# `lock` with `Monitor.Wait`/`Pulse`.

**Trap alert:** A monitor provides mutual exclusion **automatically**, but *conditions* such as "the buffer is not empty" still need condition variables and correct `while` loops (Q103).

**Interview angle:** "Monitor vs semaphore?" A monitor is higher-level and structured, so mismatched `wait`/`signal` pairs can't happen. A semaphore is a flexible but error-prone primitive. Follow-up: *"How does Java implement monitors?"* Every object has an intrinsic lock and a wait set.

---

## Q102. [M4 | What condition-variable wait does | Medium | MCQ | OA + Interview]

A thread calls `pthread_cond_wait(&c, &m)` while holding mutex `m`. Which is correct?

- A. It atomically unlocks `m` and puts the thread to sleep on `c`; before returning, it re-locks `m`
- B. It sleeps while still holding `m`
- C. It unlocks `m` and returns without re-locking it
- D. It spins until `c` is signalled

**Answer:** A

**Why:**
- The waiting thread must release `m`, or the thread that will change the condition (and signal) could never get in.
- The unlock and the sleep must happen as **one atomic step**. If they were separate, a signal sent between "unlock" and "sleep" would be lost, and the thread would sleep forever.
- On wake-up, the call re-acquires `m` before returning, so the caller can safely re-check the condition.

**Why the other options fail:**
- B: Sleeping while holding `m` deadlocks the signaller.
- C: The caller expects to hold `m` again when the call returns.
- D: Condition variables sleep; they don't spin.

**Concept recap:** `cond_wait(c, m)` = { atomically: unlock(m) + sleep on c } → woken → lock(m) → return.

**Trap alert:** Calling `cond_wait` *without* holding `m` is undefined behaviour, and it reopens the lost-wakeup race.

**Interview angle:** "Why does a condition variable need a mutex?" To make "check the condition, then sleep" atomic with respect to the thread that changes the condition and signals. Follow-up: *"Should you signal while holding the mutex?"* Either is correct. Signalling under the lock is simpler to reason about.

---

## Q103. [M4 | `while` vs `if` around wait | Medium | MCQ | OA + Interview]

There are two consumer threads and one producer thread. Each consumer runs:

```c
lock(&m);
if (count == 0)                       // ← should this be `while`?
    pthread_cond_wait(&not_empty, &m);
item = buf[--count];
unlock(&m);
```

Why must the `if` be a `while`?

- A. `while` executes faster than `if`
- B. Using `if` is a compile error
- C. Because `pthread_cond_wait` may not release the mutex
- D. Under Mesa semantics a woken thread only becomes *ready*; another consumer may take the item first, and spurious wakeups are also allowed, so the condition must be re-checked after waking

**Answer:** D

**Why:** With **Mesa semantics** (signal-and-continue), `signal` just moves a waiter to the ready queue. By the time that waiter actually re-acquires the mutex, the world may have changed.

**Working (the failure with `if`):**
1. count = 0. C1 locks `m`, sees 0, waits (releasing `m`).
2. The producer adds one item (count = 1), signals → C1 becomes ready but hasn't run yet.
3. C2 arrives, locks `m` *before* C1, sees count = 1, takes the item → count = 0, unlocks.
4. C1 finally re-locks `m`. Its `if` doesn't re-check → `buf[--count]` with count = 0 → count = −1 → **buffer underflow**.

With `while`, C1 re-checks, sees 0, and goes back to waiting.

**Why the other options fail:**
- A, B: Not true.
- C: It *does* release the mutex (Q102). The problem is the stale condition.

**Concept recap:** **Always** call `wait` inside a `while (!condition)` loop. POSIX explicitly allows **spurious wakeups**, where a thread returns from `wait` without any signal.

**Trap alert:** With a single consumer and no spurious wakeups, `if` happens to work, so tests pass. Adding a second consumer breaks it.

**Interview angle:** "Why `while` instead of `if` around `wait()`?" Give both reasons: stolen wakeups and spurious wakeups. Follow-up: *"Would `if` be correct under Hoare semantics?"* Yes, apart from spurious wakeups (Q105).

---

## Q104. [M4 | The lost signal | Medium | MCQ | OA]

```c
// shared: pthread_mutex_t m;  pthread_cond_t c;   (no other shared state)
Thread A:                     Thread B:
  lock(&m);                     lock(&m);
  cond_signal(&c);              cond_wait(&c, &m);
  unlock(&m);                   unlock(&m);
                                print("done");
```

Assume no spurious wakeups, and either thread may run first. Which is correct?

- A. "done" is always printed
- B. "done" is never printed
- C. If A signals before B starts waiting, the signal is lost and B blocks forever; "done" is printed only if B is already waiting when A signals
- D. Both threads deadlock on `m`

**Answer:** C

**Why:** A condition variable has **no memory**. `signal` wakes a thread that is waiting *right now*. If none is waiting, the signal simply vanishes. Here there is no shared state recording "the event already happened", so a late-arriving B waits for a signal that has already gone.

**Why the other options fail:**
- A: Fails when A runs first.
- B: Works when B is already waiting.
- D: B's `wait` releases `m`, so A can lock it.

**Concept recap:** Fix by pairing the condition variable with a **predicate** (shared state):
- A: `lock; ready = 1; signal; unlock;`
- B: `lock; while (!ready) wait; unlock;`

**Trap alert:** A semaphore `signal` **is remembered** (it increments the value). A condition-variable `signal` **is not** (contrast with Q94 E).

**Interview angle:** "What is a lost wakeup and how do you prevent it?" Use a predicate checked in a `while` loop, under the same mutex. Follow-up: *"`signal` vs `broadcast`?"* `broadcast` wakes all waiters. Use it when several waiters might be able to proceed, or when different waiters wait for different conditions on the same condition variable.

---

## Q105. [M4 | Hoare vs Mesa semantics | Hard | MSQ | Interview]

**Select all that apply.**

- A. Under Hoare semantics, the signaller hands the monitor to the waiter immediately, and the waiter runs knowing the condition is true.
- B. Under Mesa semantics, the signaller keeps running; the waiter is only made ready, so it must re-check the condition.
- C. pthreads and Java use Mesa semantics.
- D. With Hoare semantics and no spurious wakeups, an `if` is enough to check the condition before waiting.
- E. Under Mesa semantics, `signal()` blocks the signaller until the woken thread has finished.

**Answer:** A, B, C, D

**Why:**
- **Hoare (signal-and-wait):** the signaller is suspended and control passes straight to the waiter. Nothing can run in between, so the condition the signaller just established still holds.
- **Mesa (signal-and-continue):** the signal is only a *hint*. Other threads may run before the waiter does.

Mesa is easier to implement and allows `broadcast`, so it is what real systems use.

**Why the other options fail:**
- E: That's closer to Hoare's behaviour. Mesa's signaller never blocks on `signal`.

**Concept recap:**
- Hoare: the condition is guaranteed on wake-up, but there are extra context switches.
- Mesa: the condition is a hint, so use a `while` loop, but it is cheaper and more flexible.

**Trap alert:** Students learn monitors with Hoare semantics in theory and then write `if` in pthreads or Java code. In real code, always use `while`.

**Interview angle:** "Why do practical systems use Mesa semantics?" Fewer forced context switches, `broadcast` support, and robustness against spurious wakeups. Follow-up: *"What does Java's `notify` vs `notifyAll` correspond to?"* `signal` vs `broadcast`.

---

## Q106. [M4 | Producer-consumer: the swap that deadlocks | Hard | Code-trace | OA + Interview]

```c
// semaphores: mutex = 1, empty = N, full = 0;   buffer size N = 2
Producer:                         Consumer:
while (true) {                    while (true) {
    produce item;                     wait(full);
    wait(mutex);   // ← swapped       wait(mutex);
    wait(empty);   // ← swapped       remove item;
    add item;                         signal(mutex);
    signal(mutex);                    signal(empty);
    signal(full);                     consume item;
}                                 }
```

The producer runs alone until it blocks, and then the consumer runs. What happens?

- A. It works correctly
- B. After producing 2 items, the producer blocks on `empty` while holding `mutex`; the consumer then blocks on `mutex`: deadlock
- C. The consumer removes an item without mutual exclusion
- D. The buffer overflows to 3 items

**Answer:** B

**Why:** The producer acquires the mutex **before** checking for a free slot. When the buffer is full, it sleeps on `empty` *while still holding `mutex`*. The only process that could free a slot, the consumer, needs that same mutex first.

**Working (semaphore values after each step):**

| Step | Action | mutex | empty | full |
|---|---|---|---|---|
| start | | 1 | 2 | 0 |
| item 1 | P: wait(mutex), wait(empty), add, signal(mutex), signal(full) | 1 | 1 | 1 |
| item 2 | same | 1 | 0 | 2 |
| item 3 | P: wait(mutex) → 0; wait(empty) → −1 → **P blocks holding mutex** | 0 | −1 | 2 |
| consumer | C: wait(full) → 1 (OK); wait(mutex) → −1 → **C blocks** | −1 | −1 | 1 |

P waits for C to `signal(empty)`. C waits for P to `signal(mutex)`. Neither can proceed: **deadlock**.

**Why the other options fail:**
- A: The trace above shows it fails once the buffer is full.
- C: The mutex is still respected. Nobody gets in at all.
- D: `empty` correctly stops the third add; it blocks rather than overflowing.

**Concept recap:** In the bounded-buffer solution, **wait on the counting semaphore (`empty`/`full`) before the mutex.** The order of the `signal`s doesn't matter for deadlock.

**Trap alert:** The bug appears **only when the buffer is full** (or, for the mirror-image bug in the consumer, when it's empty). Light-load tests never trigger it.

**Interview angle:** "Write the bounded-buffer solution with semaphores." Then: *"What if I swap the two waits?"* Explain this trace. Follow-up: *"How would you write it with a mutex and two condition variables?"* Use `while (count == N) wait(not_full)` and `while (count == 0) wait(not_empty)`.

---

## Q107. [M4 | Producer-consumer semaphore roles | Medium | MSQ | OA]

Correct bounded-buffer solution with semaphores `empty`, `full` and `mutex`. **Select all that apply.**

- A. `empty` counts free slots and starts at N.
- B. `full` counts filled slots and starts at 0.
- C. `mutex` protects the shared buffer and its indices.
- D. With exactly one producer and one consumer, where the producer alone updates the `in` index and the consumer alone updates `out`, the mutex can be removed. `empty` and `full` already prevent both from touching the same slot at once and order their accesses.
- E. `full + empty` equals N at every instant.

**Answer:** A, B, C, D

**Why:**
- `empty` blocks the producer when there's no space; `full` blocks the consumer when there's nothing to take.
- With several producers (or several consumers), two of them could update the same index at once, so the mutex is needed.
- With exactly one of each, each index has a single writer, and the counting semaphores keep them on different slots. That's the classic **SPSC (single-producer single-consumer) ring buffer**.

**Why the other options fail:**
- E: Between the producer's `wait(empty)` and its `signal(full)`, one slot is counted in neither, so the sum is briefly N − 1. It equals N only when no operation is in progress.

**Concept recap:**
- **Counting semaphores** express *conditions*: is there space? is there data?
- The **mutex** protects *shared state* that several threads update.

**Trap alert:** Dropping the mutex is safe **only** for one producer and one consumer. With more of either, the shared index races.

**Interview angle:** "Implement a thread-safe bounded queue." Follow-up: *"How do lock-free SPSC queues work?"* Separate head and tail indices, each written by only one thread, with acquire/release memory ordering.

---

## Q108. [M4 | First readers-writers solution | Medium | MCQ | OA]

```c
// semaphores: rw_mutex = 1, mutex = 1;   int read_count = 0;
Writer:  wait(rw_mutex);  /* write */  signal(rw_mutex);

Reader:  wait(mutex);
         read_count++;
         if (read_count == 1) wait(rw_mutex);   // first reader locks out writers
         signal(mutex);
         /* read */
         wait(mutex);
         read_count--;
         if (read_count == 0) signal(rw_mutex); // last reader lets writers in
         signal(mutex);
```

What can happen under a continuous stream of **overlapping** readers?

- A. Readers starve
- B. Mutual exclusion between writers is broken
- C. Deadlock
- D. Writers may starve: while at least one reader is reading, new readers keep entering, so `read_count` never drops to 0 and `rw_mutex` is never released

**Answer:** D

**Why:** This is the **reader-preference** solution. Only the *first* reader takes `rw_mutex`, and only the *last* reader releases it. If a new reader always arrives before the current ones finish, the last-reader moment never comes, and a waiting writer never gets in.

**Why the other options fail:**
- A: Readers are favoured, not starved.
- B: Writers are serialized by `rw_mutex`.
- C: There's no circular wait here.

**Concept recap:**
- First readers-writers problem (**reader preference**) → writers can starve.
- Second readers-writers problem (**writer preference**) → once a writer waits, new readers are blocked, so readers can starve.
- A fair version uses a FIFO "turnstile" so readers and writers take turns in arrival order.

**Trap alert:** "Readers never conflict, so letting them in is always fine." Fairness to writers is the cost.

**Interview angle:** "Explain the readers-writers problem and its variants." Follow-up: *"Why does a reader take `mutex` around updating `read_count`?"* `read_count` is itself shared data, and `++`/`--` aren't atomic (Q77).

---

## Q109. [M4 | Readers-writer locks in practice | Hard | MSQ | Interview]

**Select all that apply.**

- A. A readers-writer lock lets several readers hold it at the same time.
- B. A writer needs exclusive access: no other reader or writer may hold the lock.
- C. A writer-preferring lock blocks new readers once a writer is waiting, which can starve readers.
- D. A readers-writer lock is always faster than a plain mutex.
- E. If two threads each hold a read lock and both try to upgrade to a write lock, they can deadlock.

**Answer:** A, B, C, E

**Why:**
- **A, B:** This is the basic read-shared / write-exclusive rule.
- **C:** Preferring writers simply shifts the starvation risk to readers.
- **E:** Each upgrader waits for all *other* readers to leave, but the other upgrader never leaves, because it's waiting too. That's a circular wait. Most RW locks therefore don't allow upgrades, or allow only one designated "upgradable" reader.

**Why the other options fail:**
- D: RW locks have extra bookkeeping (tracking the reader count). They help only when reads dominate and the critical sections are long enough. For short, frequent sections, a plain mutex is often faster.

**Concept recap:** Examples: `pthread_rwlock_t`, C++ `std::shared_mutex`, Java `ReentrantReadWriteLock`, and the Linux kernel's `rwlock_t` (with RCU as an alternative for read-mostly data).

**Trap alert:** "Read-heavy ⇒ always use an RW lock" is wrong. Measure first, because the reader-count updates themselves cause cache-line contention.

**Interview angle:** "When would you use a readers-writer lock?" For read-mostly data with long reads, such as a configuration table or cache. Follow-up: *"What is RCU?"* Read-Copy-Update: readers take no lock at all, and writers publish a new copy and free the old one only after all readers are done with it.

---

## Q110. [M4 | Barriers | Easy | MCQ | OA]

What does a **barrier** do in concurrent programming?

- A. It blocks each arriving thread until all N participating threads have arrived, then releases them all together
- B. It lets threads pass one at a time
- C. It prevents writes while reads are in progress
- D. It limits the critical section to at most N threads

**Answer:** A

**Why:** A barrier is a meeting point. Iterative parallel algorithms, such as simulations and matrix sweeps, use it so that **no thread starts phase k + 1 until every thread has finished phase k**.

**Why the other options fail:**
- B: That describes a mutex or turnstile.
- C: That describes a readers-writers lock.
- D: That describes a counting semaphore initialized to N.

**Concept recap:** APIs: `pthread_barrier_wait`, C++20 `std::barrier`, Java `CyclicBarrier`, and OpenMP `#pragma omp barrier`.

**Trap alert:** The barrier's count must equal the **number of participating threads**. If it's one too many, everyone waits forever.

**Interview angle:** "Where would you use a barrier?" In phase-based parallel computation. Follow-up: *"How do you implement a reusable barrier?"* See Q112.

---

## Q111. [M4 | Dining philosophers | Hard | MSQ | OA + Interview]

Five philosophers sit at a round table with five chopsticks, one between each pair of neighbours. A philosopher needs both adjacent chopsticks to eat. **Select all that apply.**

- A. In the naive solution (pick up the left chopstick, then the right), deadlock occurs if all five pick up their left chopstick at the same time.
- B. Allowing at most four philosophers to sit at the table at once prevents deadlock.
- C. A monitor solution in which a philosopher picks up chopsticks only when *neither* neighbour is eating is deadlock-free.
- D. That monitor solution also guarantees that no philosopher starves.
- E. Numbering the chopsticks and always picking up the lower-numbered one first prevents deadlock.

**Answer:** A, B, C, E

**Why:**
- **A:** Each philosopher holds one chopstick and waits for the next one, all the way around the table: a circular wait.
- **B:** With at most four seated, the pigeonhole principle means at least one philosopher can get both chopsticks.
- **C:** Taking both chopsticks at once removes *hold and wait*.
- **E:** A global order removes *circular wait*. The philosopher between chopsticks 4 and 0 picks up 0 first, which breaks the ring.

**Why the other options fail:**
- D: Two neighbours who take turns eating can keep the philosopher between them waiting forever. Deadlock-freedom ≠ starvation-freedom.

**Concept recap:** Each fix breaks one Coffman condition (Q115): seat limit → the cycle can't close; both-at-once → no hold-and-wait; resource ordering → no circular wait.

**Trap alert:** "No deadlock ⇒ everyone eats" is false. Starvation needs a separate fairness mechanism, such as queues or aging.

**Interview angle:** "Solve dining philosophers." Give at least two fixes and name the deadlock condition each one breaks. Follow-up: *"Is the asymmetric solution starvation-free?"* Not by itself. Fairness depends on how the underlying locks queue their waiters.

---

## Q112. [M4 | Barrier with semaphores | Medium | Code-trace | OA + Interview]

N threads use this one-time barrier:

```c
// shared: int count = 0;   semaphores: mutex = 1, turnstile = 0
wait(mutex);
count++;
if (count == N) signal(turnstile);   // last arrival opens the gate
signal(mutex);

wait(turnstile);
signal(turnstile);                   // each thread lets the next one through
/* ... code after the barrier ... */
```

After all N threads have passed, what is the value of `turnstile`, and what does that imply?

- A. 0, so the barrier can be reused as is
- B. 1, so if the barrier were reused, the next phase's first arrival would pass straight through without waiting; it is not reusable without a second turnstile
- C. N
- D. −N

**Answer:** B

**Why:** The last arrival opens the gate once. After that, each thread does a `wait` (go through) and a `signal` (hold the door for the next thread). The very last `signal` leaves the gate **open**.

**Working:**
- **Use the counting formula from Q95:** S_final = S₀ − #waits + #signals.
- S₀ = 0; waits on `turnstile` = N (one per thread); signals = 1 (from the last arrival) + N (one per thread after passing) = N + 1.
- Final value = 0 − N + (N + 1) = **1**.
- A second phase would find the gate open, so the first thread through wouldn't wait for the others.

**Why the other options fail:**
- A: Ignores the extra signal from the last arrival.
- C, D: Forget that each pass both decrements and increments.

**Concept recap:** A **reusable** barrier needs two turnstiles: gate 1 opens when count reaches N, gate 2 opens when count drops back to 0. Alternatively, use `pthread_barrier_t`, which handles reuse internally.

**Trap alert:** A barrier that works once is easy to write. Correctly *reusing* it is the hard part, and it's a classic interview follow-up.

**Interview angle:** "Implement a barrier with semaphores." Follow-up: *"Make it reusable."* Two-phase turnstile, or a generation counter checked with a condition variable.

---

## Q113. [M4 | Lock ordering | Easy | MCQ | Interview]

```c
void transfer(Account *from, Account *to, int amt) {
    lock(&from->m);
    lock(&to->m);
    from->balance -= amt;  to->balance += amt;
    unlock(&to->m);
    unlock(&from->m);
}
```

Thread T1 calls `transfer(A, B, 10)` while T2 calls `transfer(B, A, 20)` at the same time. What's the problem, and what is the standard fix?

- A. There is no problem
- B. There is a race on the balances; fix it with `volatile`
- C. Deadlock is possible: T1 holds A and waits for B while T2 holds B and waits for A. Fix it by always locking the two accounts in a global order (e.g., lower account ID first)
- D. Only T1 can starve; fix it with priorities

**Answer:** C

**Why:** The lock order depends on the *arguments*, so two calls in opposite directions acquire the same locks in opposite orders. That creates a **circular wait**. A single global ordering ("always lock the lower ID first") means every thread acquires the locks in the same sequence, so no cycle can form.

**Why the other options fail:**
- A: The deadlock is real under the right timing.
- B: The balances *are* protected. The problem is the order of acquiring locks, and `volatile` doesn't help anyway (Q77).
- D: This is a deadlock, not starvation, and priorities don't fix it.

**Concept recap:** **Lock ordering** breaks circular wait. Order by a stable key (an ID or address). C++17's `std::scoped_lock(a, b)` acquires multiple locks without deadlock.

**Trap alert:** Code can look consistent ("always lock `from`, then `to`") and still deadlock, because `from` and `to` swap between calls.

**Interview angle:** "How do you prevent deadlock with multiple locks?" Use a global lock order, `try_lock` with back-off, or a single coarse lock. Follow-up: *"What if the two accounts are the same?"* Lock only once, otherwise it's a same-thread relock (Q89).

---

## Q114. [M4 | Reentrant locks | Medium | MCQ | Interview]

```java
class Widget {
    public synchronized void doSomething() { /* ... */ }
}
class LoggingWidget extends Widget {
    @Override
    public synchronized void doSomething() {
        System.out.println("calling doSomething");
        super.doSomething();      // locks the same object again
    }
}
```

Why doesn't calling `doSomething()` on a `LoggingWidget` deadlock?

- A. `super` calls never use locks
- B. Java releases the lock before calling `super`
- C. `synchronized` instance methods lock the class, not the object
- D. Java's intrinsic locks are reentrant: the JVM records the owning thread and a hold count, so the same thread re-acquiring the lock just increments the count, and the lock is freed when the count returns to 0

**Answer:** D

**Why:** Both methods synchronize on the **same object** (`this`). If the lock weren't reentrant, the call to `super.doSomething()` would wait for a lock held by its own thread, the self-deadlock from Q89. Java's intrinsic locks track an owner and a hold count, so the owning thread can re-enter freely.

**Why the other options fail:**
- A: `super.doSomething()` is itself `synchronized`, so it does take the lock.
- B: Java holds the lock for the whole method body.
- C: Only `static synchronized` methods lock the class object.

**Concept recap:** Reentrant lock = owner + hold count. Each `lock()` increments the count and each `unlock()` decrements it; the lock is free at 0. Examples: Java `synchronized`, `ReentrantLock`, and `PTHREAD_MUTEX_RECURSIVE`.

**Trap alert:** With `ReentrantLock`, every `lock()` needs a matching `unlock()`. Lock twice and unlock once, and the lock is still held, so other threads block forever. Always `unlock()` in a `finally` block.

**Interview angle:** "What is a reentrant lock and why does Java use one?" So that synchronized methods can call each other and their overrides safely. Follow-up: *"Downside?"* Re-entering partway through a method can expose an object whose invariants are temporarily broken.

---

## Q115. [M5 | The four deadlock conditions | Easy | MCQ | OA]

Which of these is **NOT** one of the four necessary conditions for deadlock (the Coffman conditions)?

- A. Mutual exclusion
- B. Hold and wait
- C. Starvation
- D. Circular wait

**Answer:** C

**Why:** Deadlock can occur only if **all four** of these hold at the same time:
1. **Mutual exclusion:** at least one resource can be held by only one process at a time.
2. **Hold and wait:** a process holds at least one resource while waiting for another.
3. **No preemption:** a resource can't be taken away from its holder; only the holder can release it.
4. **Circular wait:** there is a cycle P₀ → P₁ → … → Pₙ → P₀, where each process waits for a resource held by the next.

Starvation isn't on the list. It's a *different* problem.

**Why the other options fail:** A, B and D are three of the four conditions. The fourth, "no preemption", isn't listed among the options.

**Concept recap:** ME + hold-and-wait + no preemption + circular wait = **necessary** for deadlock. Breaking any one of them prevents deadlock (covered in Batch 7).

**Trap alert:** The condition is "**no** preemption". An option worded "preemption" is a trick. Also, starvation can occur without deadlock, and vice versa.

**Interview angle:** "What are the necessary conditions for deadlock?" Name all four, with a one-line meaning each. Follow-up: *"Which is easiest to break in practice?"* Usually circular wait, by imposing a lock order (Q113).

---

## Q116. [M5 | Cycles and deadlock in a RAG | Medium | MSQ | OA + Interview]

**Select all that apply.**

- A. All four Coffman conditions must hold simultaneously for deadlock to occur.
- B. In a resource-allocation graph where every resource type has exactly one instance, a cycle implies deadlock.
- C. When resource types have multiple instances, a cycle does not necessarily imply deadlock.
- D. A resource-allocation graph with no cycle can still be deadlocked.
- E. Starvation is one of the four necessary conditions.

**Answer:** A, B, C

**Why:**
- **A:** The four conditions are jointly necessary.
- **B:** With single instances, every process on a cycle waits for a resource held by the next process on the cycle, and none can get it. So a cycle is both **necessary and sufficient**.
- **C:** With multiple instances, a process outside the cycle may hold another instance and eventually release it, which breaks the wait (Q117).

**Why the other options fail:**
- D: Deadlock requires a circular wait, so a graph with no cycle can't be deadlocked.
- E: See Q115.

**Concept recap:**
- No cycle ⇒ no deadlock.
- Cycle + single-instance resources ⇒ deadlock.
- Cycle + multiple instances ⇒ *possible* deadlock; you must run a reduction or detection algorithm.

**Trap alert:** "Cycle ⇒ deadlock" is true only when every resource type has a **single instance**.

**Interview angle:** "How do you detect deadlock from a resource graph?" Follow-up: *"What changes with multiple instances?"* Use a detection algorithm similar to Banker's (covered in Batch 7).

---

## Q117. [M5 | RAG with multiple instances | Medium | Numerical | OA]

**Resources:** R1 has 2 instances, and R2 has 2 instances.
**Assignment edges** (resource → holder): R1 → P2, R1 → P3, R2 → P1, R2 → P4.
**Request edges** (waiter → resource): P1 → R1, P2 → R2.

```
          ┌──── R1 (2) ────┐
   P1 ───►│ ●──► P2        │      P3 holds one R1, requests nothing
          │ ●──► P3        │
          └────────────────┘
   P2 ───► R2 (2): ●──► P1,  ●──► P4      P4 holds one R2, requests nothing
```

There is a cycle **P1 → R1 → P2 → R2 → P1**. Is the system deadlocked?

- A. Yes, because any cycle means deadlock
- B. No: P3 and P4 aren't waiting for anything. When either one finishes and releases its instance, a waiting process can proceed, and the cycle breaks
- C. Yes, because P1 and P2 both hold resources and wait
- D. It cannot be determined without process priorities

**Answer:** B

**Why:** Use **graph reduction**: repeatedly find a process whose requests can all be satisfied, let it finish, and give back its resources. If every process can be reduced, there is no deadlock.

**Working:**
1. **Available** instances: R1 = 2 − 2 = 0; R2 = 2 − 2 = 0.
2. P3 requests nothing → it finishes and releases one R1 → Available R1 = 1.
3. P1 requests R1 → granted → P1 finishes and releases R1 and R2 → Available: R1 = 2, R2 = 1.
4. P2 requests R2 → granted → P2 finishes.
5. P4 requests nothing → it finishes.
6. Every process has been reduced → **no deadlock**.

**Why the other options fail:**
- A: "Cycle ⇒ deadlock" holds only for single-instance resources (Q116).
- C: Hold-and-wait is necessary for deadlock, but not sufficient.
- D: Priorities are irrelevant to whether a deadlock exists.

**Concept recap:** Reduction means: find a process whose requests ≤ Available, let it finish, and add its allocation back to Available. Repeat. Any process left unreduced at the end is deadlocked.

**Trap alert:** Don't stop at "I found a cycle". With multiple instances, always run the reduction.

**Interview angle:** Rarely asked verbally. Follow-up: *"How is this reduction related to the Banker's safety algorithm?"* It's the same idea, with current requests in place of maximum future needs (Batch 7).

---

## Q118. [M5 | Wait-for graphs | Medium | MCQ | OA + Interview]

Which statement about **wait-for graphs** is correct?

- A. A wait-for graph is formed by removing the resource nodes from a RAG and collapsing the edges: Pi → Pj means Pi waits for a resource held by Pj. With single-instance resources, a cycle means deadlock
- B. A wait-for graph works for any number of instances per resource type
- C. A wait-for graph contains resource nodes and claim edges
- D. A cycle in a wait-for graph never implies deadlock

**Answer:** A

**Why:** If Pi → R → Pj in the RAG (Pi requests R, which is held by Pj), the wait-for graph has the direct edge **Pi → Pj**. Detecting deadlock then reduces to **cycle detection** among processes. The OS (or a database's lock manager) keeps this graph and checks it periodically.

**Why the other options fail:**
- B: With multiple instances, an edge can't say *which* holder you're waiting for. Use the matrix-based detection algorithm instead.
- C: Claim edges belong to the RAG used for **avoidance** (covered in Batch 7), not to the wait-for graph.
- D: With single-instance resources, a cycle *does* mean deadlock.

**Concept recap:** Wait-for graph = processes only. Cycle ⇔ deadlock (single instances). Cycle detection costs **O(n²)** for n processes, using DFS on a graph with up to n² edges.

**Trap alert:** Use the wait-for graph only when every resource has **one instance**, such as mutexes and row locks.

**Interview angle:** "How do databases detect deadlocks?" They build a wait-for graph among transactions, find a cycle, and abort a victim. Follow-up: *"How often should detection run?"* It's a trade-off: frequent checks find deadlocks quickly but cost CPU. Some systems run detection only when a lock request waits longer than a timeout.

---

## Q119. [M5 | Deadlocks beyond lock pairs | Hard | MSQ | Interview]

Which of these situations can deadlock? **Select all that apply.**

- A. Two threads each lock the same two mutexes, in opposite orders.
- B. A single thread calls `lock()` twice on the same non-recursive mutex.
- C. Two processes each block reading from a pipe that only the other one writes to, and neither writes first.
- D. Two threads acquire the same two locks, always in the same global order.
- E. A parent calls `waitpid()` on its child, while the child is blocked writing to a full pipe that only the parent reads.

**Answer:** A, B, C, E

**Why:** Deadlock is any **circular wait** in which no one can proceed. The resources don't have to be locks:
- **A:** The classic two-lock cycle (Q113).
- **B:** A cycle of length 1: the thread waits for itself (Q89).
- **C:** A **communication deadlock**: each process waits for a message the other will only send after receiving one.
- **E:** The parent waits for the child to exit; the child waits for the parent to drain the pipe. That's a cycle through two different kinds of resource.

**Why the other options fail:**
- D: A consistent global order makes a cycle impossible.

**Concept recap:** Model "waiting for X held by (or only producible by) Y" as an edge X → Y. Any cycle among non-preemptible waits is a deadlock.

**Trap alert:** E shows up in real code. For example, Python's `subprocess` documentation warns against `wait()` with large piped output and recommends `communicate()` for exactly this reason. Read the pipes while waiting.

**Interview angle:** "Give a deadlock example that doesn't involve two mutexes." C and E are good answers. Follow-up: *"How would you fix E?"* Drain the pipe with a separate thread or with `select`/`poll`, or read until EOF before calling `waitpid`.

---

## Q120. [M5 | RAG edge meanings | Easy | MCQ | OA]

In a resource-allocation graph, what does an edge **from resource R to process P** mean?

- A. P is requesting R
- B. P has released R
- C. P may need R in the future (a claim edge)
- D. An instance of R is assigned to (held by) P

**Answer:** D

**Why:** RAG edges point in the direction of "who is waiting on whom":
- **Request edge P → R:** P is waiting for R.
- **Assignment edge R → P:** an instance of R is held by P. (In drawings, it starts from a dot inside R's box, one dot per instance.)

**Why the other options fail:**
- A: A request goes P → R, the opposite direction.
- B: A release simply removes the edge.
- C: A claim edge is dashed and points P ⇢ R. It's used only in avoidance algorithms (covered in Batch 7).

**Concept recap:** P → R = request; R → P = assignment; P ⇢ R (dashed) = claim. Granting a request turns P → R into R → P.

**Trap alert:** Reversing the arrow direction makes a cycle look like a non-cycle, or the other way round. Fix the convention first.

**Interview angle:** Usually part of a larger question: "Draw the RAG for this scenario and decide whether it's deadlocked." Follow-up: *"How is a claim edge used in avoidance?"* A request is granted only if converting its claim edge into an assignment edge creates no cycle.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 101 | C | 106 | B | 111 | A, B, C, E | 116 | A, B, C |
| 102 | A | 107 | A, B, C, D | 112 | B | 117 | B |
| 103 | D | 108 | D | 113 | C | 118 | A |
| 104 | C | 109 | A, B, C, E | 114 | D | 119 | A, B, C, E |
| 105 | A, B, C, D | 110 | A | 115 | C | 120 | D |

---

**Concepts tested so far (to avoid repeats):** Batches 1–5, plus: what a monitor guarantees · `cond_wait` releasing the mutex atomically and re-acquiring it · `while` vs `if` (stolen and spurious wakeups) · condition-variable signals aren't remembered, and the predicate fix · Hoare vs Mesa · bounded buffer with the waits swapped → deadlock · roles of empty/full/mutex, the SPSC exception, the empty+full ≠ N moment · reader-preference starvation of writers · RW lock properties and the upgrade deadlock · barrier purpose · dining philosophers fixes and starvation · one-time barrier leaves the turnstile at 1 · lock-ordering deadlock and its fix · Java reentrancy (owner + hold count) · the four Coffman conditions · cycle vs deadlock with single vs multiple instances · RAG reduction with multiple instances · wait-for graph and O(n²) detection · communication and pipe deadlocks · RAG edge directions
*(M4 is complete. In M5, still to come in Batch 7: prevention, avoidance, safe/unsafe states, Banker's safety and request algorithms, detection, recovery, and deadlock vs starvation vs livelock.)*

**Completed: 120 / 300    Remaining: 180**
