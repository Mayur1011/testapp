# OS MCQ/MSQ Bank — Batch 14 (Q261–Q280)

**What this batch covers:**
- **M10, protection & security completed (Q261–271):** authentication vs authorization, the access matrix, ACLs vs capability lists, protection domains and domain switching, the principle of least privilege, isolation mechanisms, memory protection, and privilege escalation at the conceptual level (what it means and how the OS defends against it).
- **M11, Linux/Unix internals begins (Q272–280):** the process model, file descriptors 0/1/2, `fork`/`exec`/`wait` in practice, fd inheritance across `fork`, shared file offsets, `dup2` redirection, and close-on-exec.

> Platform is Linux/POSIX unless stated. Code assumes the usual headers and that calls succeed unless the question says otherwise.

---

## Q261. [M10 | Authentication vs authorization | Easy | MCQ | Interview]

Which pair of definitions is correct?

- A. Authentication decides what you may do; authorization decides who you are
- B. Authentication verifies **who** you are (identity); authorization decides **what** that identity is allowed to do
- C. They are the same step
- D. Authentication applies to users; authorization applies only to files

**Answer:** B

**Why:**
- **Authentication** (authN) establishes identity: a password, a key, a fingerprint, a token.
- **Authorization** (authZ) decides, for an already-identified principal, which actions on which resources are permitted.

Authentication always comes first: you can't decide what someone may do until you know who they are.

**Why the other options fail:**
- A: The two definitions are swapped.
- C: They are distinct steps.
- D: Authorization governs every resource type: files, memory, devices, system calls.

**Concept recap:** authN = "who are you?" → authZ = "what may you do?". A third A, **accounting/auditing**, records what was done.

**Trap alert:** Logging in successfully (authN) doesn't mean you may read a given file (authZ). The two checks are separate.

**Interview angle:** "Authentication vs authorization?" Follow-up: *"Where does each live in a web app?"* authN at the login/session layer, authZ at each protected endpoint or resource check.

---

## Q262. [M10 | The access matrix | Medium | MCQ | OA]

In an **access matrix**, what do the rows and columns represent, and what does an entry contain?

- A. Rows are objects, columns are domains; an entry is a password
- B. Rows are domains (subjects), columns are objects (resources); an entry lists the access rights that domain has on that object
- C. Rows are users, columns are groups; an entry is a priority
- D. Rows and columns are both processes; an entry is a message

**Answer:** B

**Why:** The **access matrix** is the abstract model of protection. Each **domain** (a subject: a user, a process, or a role) is a row. Each **object** (a file, a device, memory) is a column. Entry (i, j) holds the rights domain i has on object j, such as {read, write, execute}.

**Why the other options fail:**
- A, C, D: None match the domain × object structure with rights in the cells.

**Concept recap:** The matrix is usually **sparse** (most domains can't touch most objects), so it's stored by rows or by columns rather than in full (Q263).

**Trap alert:** The matrix is a **model**, not a literal in-memory table. Real systems implement it as ACLs or capability lists.

**Interview angle:** "What is the access-matrix model?" Follow-up: *"How is it actually stored?"* By column (ACLs) or by row (capability lists), because storing the full matrix wastes space.

---

## Q263. [M10 | ACLs vs capability lists | Hard | MSQ | OA + Interview]

**Select all that apply.**

- A. An access control list (ACL) stores, **with each object**, the list of domains and their rights. This is the access matrix stored **by column**.
- B. A capability list stores, **with each domain**, the objects it may access and its rights. This is the matrix stored **by row**.
- C. "Who can access this file?" is easy to answer with ACLs; "what can this process access?" is easy with capability lists.
- D. Revoking one domain's access to one object is generally easier with ACLs than with capability lists.
- E. Unix file permissions (owner/group/other) are a form of capability list.

**Answer:** A, B, C, D

**Why:**
- **A, B:** The two schemes are the two ways to slice the matrix.
- **C:** Each scheme answers questions along the axis it's indexed by.
- **D:** With an ACL, revocation just edits that object's list. With capabilities, the right is held by the domains, possibly copied around, so finding and revoking every copy is harder.

**Why the other options fail:**
- E: Unix permissions are a compact **ACL** (stored with the file). A **file descriptor** is closer to a capability: once you hold it, you can use the file without re-checking the path's permissions (Q211 trap).

**Concept recap:**
- ACL: per object. Easy to see and revoke per object. Used by file systems.
- Capability: per domain. Easy to carry and delegate. Used by fds, Linux capabilities, and seL4.

**Trap alert:** A **capability** is an unforgeable token that *is* the right. Possessing it grants access. That makes delegation easy and revocation hard.

**Interview angle:** "ACL vs capabilities?" Follow-up: *"Which is a Unix file descriptor?"* A capability: permission is checked at `open`, and the fd then carries the right (even across `fork`).

---

## Q264. [M10 | Protection domains | Medium | MCQ | Interview]

A process running as a normal user executes `/usr/bin/passwd` to change its password. `passwd` is **setuid root**, and it updates `/etc/shadow`, which only root may write. What does the setuid bit do here, in protection-domain terms?

- A. It gives the user permanent root access
- B. It disables file permissions for that user
- C. During execution of `passwd`, the process's effective user ID becomes root, so it temporarily switches into a more privileged **protection domain** for exactly that task; it returns to the user's domain when `passwd` exits
- D. It copies `/etc/shadow` into the user's home directory

**Answer:** C

**Why:** A **protection domain** is the set of rights a process currently has. The **setuid** mechanism is a controlled way to *switch* domains: while a setuid-root program runs, the process's effective UID is root, so it can do the one privileged job the program was written for. When the program exits, the process is back to the user's ordinary rights.

**Why the other options fail:**
- A: The elevation lasts only while `passwd` runs, and only within what `passwd` chooses to do.
- B: Permissions still apply; the process just runs under a different effective identity.
- D: `passwd` edits the real file in place, under root's authority.

**Concept recap:** A domain switch grants more rights for a limited task. setuid is the Unix tool for it, and it's why `passwd`, `sudo` and `ping` (historically) are setuid root.

**Trap alert:** setuid programs are a prime security target, exactly because they bridge from a low-privilege domain to a high-privilege one. A bug in one can hand its elevated rights to an attacker (Q269).

**Interview angle:** "What does the setuid bit do?" Explain the temporary domain switch. Follow-up: *"Why is a setuid-root shell dangerous?"* It would give its caller a root shell outright — unrestricted, not scoped to one task.

---

## Q265. [M10 | Least privilege | Easy | MCQ | OA + Interview]

What does the **principle of least privilege** state?

- A. Every process should run as root for simplicity
- B. Each user, program or process should be granted only the privileges it needs to do its job, and no more
- C. Privileges should be assigned once and never changed
- D. All privileges should be shared equally among all users

**Answer:** B

**Why:** Least privilege **limits the damage** of any bug or compromise. If a web server runs as an unprivileged user confined to its own directory, a flaw in it can't read `/etc/shadow` or touch the kernel. The smaller the granted rights, the smaller the "blast radius".

**Why the other options fail:**
- A: Running everything as root is the opposite, and maximally dangerous.
- C: Privileges should be granted per task and dropped when no longer needed.
- D: Equal sharing ignores what each component actually needs.

**Concept recap:** Least privilege → smaller attack surface, better fault isolation, easier auditing. Applications of it: dropping privileges after startup, per-service user accounts, sandboxing, and fine-grained Linux capabilities instead of full root.

**Trap alert:** "Grant broad rights now to avoid permission errors later" is the anti-pattern this principle warns against.

**Interview angle:** "What is least privilege, and how do you apply it?" Drop privileges after binding a low port, use a dedicated service user, and sandbox. Follow-up: *"How do Linux capabilities help?"* They split root's power into units (e.g., `CAP_NET_BIND_SERVICE`), so a program can bind port 80 without all of root's other powers.

---

## Q266. [M10 | Memory protection | Medium | MSQ | OA]

**Select all that apply.**

- A. The user/supervisor bit in a page-table entry stops user-mode code from reading kernel pages.
- B. The no-execute (NX/XD) bit marks data pages non-executable, which blocks running injected code from the stack or heap.
- C. Each process having its own page table isolates one process's memory from another's.
- D. Base and limit registers can confine a process to a contiguous memory region.
- E. Memory protection is enforced by the compiler at build time.

**Answer:** A, B, C, D

**Why:**
- **A:** The MMU checks the U/S bit on every access, so user code faults if it touches a supervisor page (Q152).
- **B:** NX enforces **W^X** — a page is writable or executable, not both — which stops classic code injection (Q152 follow-up).
- **C:** Separate page tables mean process A simply has no mapping to process B's frames (Q151 D).
- **D:** Base/limit is the simplest hardware confinement (Q138).

**Why the other options fail:**
- E: Memory protection is enforced by **hardware** (the MMU) at **run time**, on every access. The compiler can add checks, but it can't enforce isolation between processes.

**Concept recap:** Memory protection is a **hardware** mechanism: the MMU checks translation and permission on every access, and only the kernel can change the mappings.

**Trap alert:** Language-level safety (bounds checks, type safety) is a separate layer. It doesn't replace hardware memory protection, and hardware protection doesn't prevent a buffer overflow *within* a process's own writable memory.

**Interview angle:** "How does the OS stop one process from reading another's memory?" Per-process page tables plus the MMU. Follow-up: *"What is W^X / DEP?"* No page is both writable and executable, which blocks executing injected data.

---

## Q267. [M10 | Isolation mechanisms | Medium | MCQ | Interview]

Which statement about **isolation** mechanisms is correct?

- A. A virtual machine and a container isolate workloads the same way, at the same level
- B. A container isolates processes using kernel features (namespaces for what a process can see, and cgroups for how much it can use), while sharing the host kernel; a virtual machine runs a full guest OS on virtual hardware under a hypervisor, giving stronger isolation at a higher cost
- C. Containers provide stronger isolation than virtual machines
- D. A hypervisor runs inside each container

**Answer:** B

**Why:**
- A **container** is a group of processes that the host kernel makes look like they have their own system: **namespaces** give each container its own view of PIDs, mounts, network and users, and **cgroups** limit its CPU, memory and I/O. All containers share **one** kernel.
- A **virtual machine** runs a full guest OS on virtual hardware, managed by a **hypervisor**. The guest has its own kernel, so a guest-kernel compromise is contained within that VM.

**Why the other options fail:**
- A: They isolate at different levels (OS features vs virtual hardware).
- C: VMs isolate more strongly, because they don't share a kernel.
- D: The hypervisor sits *below* the VMs, not inside a container.

**Concept recap:**
- Container: lightweight, fast to start, shares the kernel, weaker isolation.
- VM: heavier, slower to start, its own kernel, stronger isolation.

Full process isolation (one process ↔ one address space) is the baseline both build on.

**Trap alert:** A kernel vulnerability can break **container** isolation, because the kernel is shared. VMs are more resistant, since they expose a much smaller interface (the virtual hardware).

**Interview angle:** "Containers vs VMs?" Covered again in M11 (Q288). Follow-up: *"How do you harden a container?"* seccomp to restrict system calls, drop Linux capabilities, read-only root filesystem, user namespaces, and no privileged mode.

---

## Q268. [M10 | Protecting the access matrix itself | Hard | MSQ | Interview]

**Select all that apply.**

- A. The ability to change entries in the access matrix is itself a right (e.g., "owner", "control"), controlled by the matrix.
- B. A "copy" right lets a domain pass one of its rights on another object to a different domain.
- C. The "owner" right on an object lets a domain add or remove rights for that object.
- D. If any domain could freely edit the matrix, the whole protection system would be meaningless.
- E. The access matrix must be stored in user space so applications can edit it quickly.

**Answer:** A, B, C, D

**Why:** The access matrix must **protect itself**, or it's pointless (D). So the right to modify it is just another right recorded *in* the matrix:
- **copy (\*)**: propagate a right to another domain.
- **owner**: change the rights for an object you own.
- **control**: change the rights of another *domain* (its row).

**Why the other options fail:**
- E: The matrix is protection-critical, so it must live in **kernel-protected** storage (file metadata, kernel tables), not in editable user memory.

**Concept recap:** The matrix governs its own modification through meta-rights (copy, owner, control). This is "mechanism protecting policy".

**Trap alert:** A protection system is only as strong as the control over **who may change the policy**. If a user can edit their own permissions freely, enforcement collapses.

**Interview angle:** "Who can change file permissions in Unix?" The file's owner and root (`chmod`). Follow-up: *"Why can't an ordinary user grant themselves write access to a root-owned file?"* They lack the owner right on it, and only the kernel edits the metadata.

---

## Q269. [M10 | Privilege escalation (concept) | Medium | MCQ | OA + Interview]

In security terms, what is **privilege escalation**, and which OS principle most directly limits the damage it can cause?

- A. A scheduler raising a process's priority; it is limited by the time quantum
- B. Gaining rights beyond those originally granted — **vertical** (a normal user ending up with admin/root rights) or **horizontal** (acting as a different same-level user). The principle of least privilege limits the damage, by keeping each component's starting rights small
- C. A process using more memory than allowed; it is limited by paging
- D. A user legitimately logging in as root with the correct password

**Answer:** B

**Why:**
- **Vertical** escalation: a low-privilege principal obtains higher-privilege rights (user → root).
- **Horizontal** escalation: a principal gains the rights of another principal at the same level (reading another user's data).

Escalation usually exploits a flaw in a component that already bridges privilege levels, such as a setuid program or a kernel interface. **Least privilege** (Q265) is the main structural defence: if every component starts with minimal rights and a compromised one has little to offer, there is less to escalate *to*.

**Why the other options fail:**
- A, C: These misuse the term for scheduling or memory limits.
- D: Logging in with the correct root password is *authorized* access, not escalation.

**Concept recap:** Privilege escalation = ending up with more authority than granted. Defences the OS provides: least privilege, dropping privileges promptly, sandboxing (seccomp, namespaces), reducing the number and attack surface of setuid programs, and keeping the kernel and its interfaces patched.

**Trap alert:** The elevated-but-authorized case (setuid `passwd`, `sudo` with the right rules) is **not** escalation. Escalation means obtaining rights the policy never intended you to have.

**Interview angle:** "What is privilege escalation, at a high level?" Define vertical vs horizontal, and name least privilege and a small setuid surface as defences. Follow-up: *"Why are setuid-root programs audited so heavily?"* They are exactly the bridges between privilege levels, so a flaw in one is a natural escalation path.

---

## Q270. [M10 | Defending privileged programs | Medium | MSQ | OA + Interview]

A service must bind to a low port (80), then spend the rest of its life serving files from one directory. Which practices follow good OS-security principles? **Select all that apply.**

- A. Start with just enough privilege to bind the port, then **drop privileges** to an unprivileged user before handling any requests.
- B. Run the whole service as root the entire time, for convenience.
- C. Restrict the set of system calls the process may make (e.g., with seccomp), so a compromised process can do less.
- D. Confine the process's file view (e.g., a chroot, a mount namespace, or a container) so it can reach only its own directory.
- E. Grant only the specific capability needed (`CAP_NET_BIND_SERVICE`) rather than full root.

**Answer:** A, C, D, E

**Why:** Each practice shrinks what a compromise can achieve:
- **A:** Privilege is needed only momentarily, to bind the port. Dropping it afterward means a later flaw runs with few rights.
- **C:** A reduced system-call set limits the actions available to injected code.
- **D:** Confining the file view limits which files a compromise can touch.
- **E:** A single fine-grained capability avoids handing over all of root's powers.

**Why the other options fail:**
- B: Running as root throughout maximizes the blast radius. This is the anti-pattern least privilege exists to prevent.

**Concept recap:** The pattern: **acquire the minimum privilege, use it briefly, drop it, then sandbox**. This combines least privilege (Q265) with isolation (Q267).

**Trap alert:** Binding the port is the *only* step that needs elevation. Keeping root beyond that point is pure risk with no benefit.

**Interview angle:** "How would you run a public-facing service securely?" Describe the acquire-use-drop-sandbox sequence. Follow-up: *"What does seccomp add on top of dropping privileges?"* Even with low privileges, it narrows the kernel attack surface by blocking system calls the service never needs.

---

## Q271. [M10 | Separating mechanism from policy | Medium | MCQ | Interview]

Why do OS designers insist on separating **mechanism** from **policy** in protection?

- A. It makes the kernel larger
- B. The kernel provides general enforcement mechanisms (permission checks, the MMU, capabilities), while the specific rules — who may access what — are set as policy by administrators and users; this lets the policy change without rewriting the mechanism
- C. Policy and mechanism must always be identical
- D. It removes the need for authentication

**Answer:** B

**Why:** A **mechanism** is a general tool, like "check this permission bit before allowing the access". A **policy** is a specific decision, like "bob may read this file". Keeping them separate means the same kernel mechanism supports many different policies, and you can change the rules (via `chmod`, group membership, config) without touching kernel code.

**Why the other options fail:**
- A: Separation tends to keep the kernel *smaller* and more general.
- C: The whole point is that they're distinct.
- D: Authentication is still required; it's a separate concern.

**Concept recap:** Mechanism = how enforcement works (fixed, general, in the kernel). Policy = what is allowed (flexible, set by admins and users). The access matrix is the mechanism; the entries are the policy.

**Trap alert:** Baking a specific policy into the mechanism makes the system rigid. The flexibility comes precisely from the separation.

**Interview angle:** "Why separate mechanism from policy?" Flexibility and reuse. Follow-up: *"Give a scheduling example."* The kernel provides preemption and priorities (mechanism); which policy to run (CFS, real-time classes, nice values) is configurable.

---

## Q272. [M11 | File descriptors 0, 1, 2 | Easy | MCQ | OA]

On Linux, what are file descriptors **0, 1 and 2** by convention?

- A. 0 = stdout, 1 = stdin, 2 = stderr
- B. 0 = stdin, 1 = stdout, 2 = stderr
- C. 0 = stderr, 1 = stdout, 2 = stdin
- D. They are reserved and unusable

**Answer:** B

**Why:** By convention, every process starts with three open descriptors:
- **0 = standard input (stdin)** — where it reads input.
- **1 = standard output (stdout)** — where it writes normal output.
- **2 = standard error (stderr)** — where it writes error messages.

The shell sets these up (usually to the terminal) before `exec`ing a program. Keeping normal output (1) and errors (2) separate is what lets `./prog > out.txt 2> err.txt` split them.

**Why the other options fail:**
- A, C: Wrong ordering.
- D: They are ordinary descriptors and can be redirected or closed.

**Concept recap:** A new fd always gets the **lowest free number** (Q47). Since 0, 1 and 2 are taken, the first file a program opens is usually fd 3.

**Trap alert:** If a program closes fd 1 and then opens a file, that file becomes fd 1, so its normal output silently goes to the file. This is the basis of redirection (Q278).

**Interview angle:** "What are fds 0, 1, 2?" Follow-up: *"Why send errors to stderr instead of stdout?"* So that redirecting or piping normal output doesn't swallow the error messages.

---

## Q273. [M11 | The Linux process model | Medium | MSQ | Interview]

**Select all that apply.**

- A. Every process except the first has a parent; the process tree is rooted at PID 1 (init/systemd).
- B. On Linux, threads and processes are both created by `clone()`; a thread is a task that shares the address space (and more) with its creator.
- C. `getpid()` returns the process ID; `getppid()` returns the parent's PID.
- D. A process's PID is reused immediately after it exits, before it is reaped.
- E. Each process has its own page table, file-descriptor table, and current working directory.

**Answer:** A, B, C, E

**Why:**
- **A:** The kernel starts PID 1, which is the ancestor of every user process.
- **B:** Linux represents both threads and processes as `task_struct`s. `clone()` with sharing flags (address space, fds, signals) makes a thread; without them, it makes a process. `fork()` is a `clone()` with no sharing.
- **C, E:** Standard per-process state.

**Why the other options fail:**
- D: A PID stays reserved while the process is a **zombie** (Q30). It can be reused only after the parent reaps it.

**Concept recap:** Linux has one task abstraction. Threads and processes differ only in *what they share*. This is why Linux threading is "1:1" with cheap kernel tasks (Q36).

**Trap alert:** "Linux threads are totally different from processes" overstates it. They're tasks with different sharing flags.

**Interview angle:** "How does Linux represent threads and processes?" One `task_struct` each, created by `clone()` with different flags. Follow-up: *"What does `fork()` share with the child?"* Nothing writable directly; it copies (via COW) rather than shares.

---

## Q274. [M11 | fork/exec/wait in a shell | Medium | MCQ | OA + Interview]

When a shell runs `ls -l`, which sequence of system calls does it use?

- A. `exec("ls")` only, in the shell process
- B. `fork()` to create a child; the child calls `exec("ls", …)` to replace itself with `ls`; the parent shell calls `wait()` for the child to finish
- C. `fork()` only; the child runs `ls` code the shell already contains
- D. `clone()` to share memory, then run `ls` in a thread

**Answer:** B

**Why:** The shell must keep running after `ls` finishes, so it can't just `exec` `ls` in itself (that would replace the shell). The classic pattern:
1. `fork()` → a child process.
2. In the child, `exec("ls", …)` → replace its image with the `ls` program.
3. In the parent, `wait()` → block until the child exits, then show the next prompt.

**Why the other options fail:**
- A: `exec` in the shell itself would replace the shell with `ls`.
- C: The shell doesn't contain `ls`'s code; `exec` loads it from disk.
- D: `ls` runs as a separate process, not a thread sharing the shell's memory.

**Concept recap:** The Unix idiom is **fork + exec + wait**. The gap between fork and exec is where the child sets up redirection, changes directory, or drops privileges (Q278).

**Trap alert:** For a **background** command (`ls &`), the shell skips the `wait` (or reaps later via SIGCHLD) so it can accept the next command immediately.

**Interview angle:** "What happens when you type a command in a shell?" fork, exec, wait. Follow-up: *"Why fork first instead of exec directly?"* So the shell survives, and so the child can rearrange its environment before exec.

---

## Q275. [M11 | fd inheritance across fork | Medium | MCQ | OA]

A process opens a file (getting fd 3), then calls `fork()`. Which is true of the child?

- A. The child does not have fd 3
- B. The child has its own copy of fd 3, referring to the **same** open-file-table entry, so parent and child **share the file offset**
- C. The child has fd 3, but it points to a separate open-file entry with an independent offset
- D. fd 3 is closed in the parent after fork

**Answer:** B

**Why:** `fork` copies the **fd table** (the per-process level), but the copied entries point to the **same system-wide open-file-table entries** as the parent's. Since the **offset lives in that shared entry** (Q223), parent and child share one offset for fd 3. If one reads 100 bytes, the other's next read continues from there.

**Why the other options fail:**
- A: The child inherits all of the parent's open fds.
- C: A **separate** `open()` would give a separate entry (Q223). `fork` **shares** the existing one.
- D: The parent keeps its fds.

**Concept recap:** Separate `open()` → separate offset. `fork` or `dup` → **shared** offset. This shared offset is what makes `ls > f; echo done >> f`-style cooperation work, and it's how a shell pipeline's two ends stay coordinated.

**Trap alert:** This is the subtle complement to Q223. The rule is about *which open-file entry* the fds point to, not about having the same fd number.

**Interview angle:** "Do parent and child share the file offset after fork?" Yes, because their fds point to the same open-file entry. Follow-up: *"Why does appending from two forked children usually not interleave badly?"* They share the offset, and `O_APPEND` makes each write seek-to-end atomically (Q48).

---

## Q276. [M11 | Shared offset after fork | Hard | Code-trace | OA + Interview]

```c
int fd = open("out.txt", O_WRONLY | O_CREAT | O_TRUNC, 0644);
if (fork() == 0) {
    write(fd, "AAA", 3);   // child
    _exit(0);
}
wait(NULL);
write(fd, "BBB", 3);        // parent, after the child has finished
```

Assume all calls succeed. What does `out.txt` contain?

- A. `BBB` (the child's write is overwritten)
- B. `AAABBB`
- C. `AAA` only
- D. `BBBAAA`

**Answer:** B

**Why:** Parent and child share the open-file entry, so they share the **offset** (Q275). The `wait(NULL)` makes the parent write only after the child has finished.
1. Child writes "AAA" at offset 0 → offset advances to 3.
2. That advance is in the **shared** entry, so the parent sees offset 3.
3. Parent writes "BBB" at offset 3 → file = **"AAABBB"**.

**Why the other options fail:**
- A: Would happen only if each had its **own** offset, both starting at 0 (i.e., two separate `open()`s).
- C: Ignores the parent's write.
- D: Wrong order; `wait` guarantees the child goes first.

**Concept recap:** Shared offset + `wait` ordering → the writes append cleanly, "AAABBB".

**Trap alert:** Replace `fork`-sharing with two independent `open()`s and you'd get "BBB" (the second write overwrites, because each has its own offset at 0). The sharing is exactly what makes the difference.

**Interview angle:** "Why doesn't the parent's write overwrite the child's here?" The shared offset already advanced past "AAA". Follow-up: *"What if you removed the `wait`?"* The two writes could still both land (offset is shared), but their order would be nondeterministic.

---

## Q277. [M11 | dup2 and redirection | Hard | Code-trace | OA + Interview]

```c
int fd = open("log.txt", O_WRONLY | O_CREAT | O_TRUNC, 0644);
dup2(fd, 1);    // make fd 1 (stdout) refer to log.txt
close(fd);
printf("hello\n");
fflush(stdout);
```

Where does "hello" go, and what does `dup2(fd, 1)` do?

- A. To the terminal; `dup2` has no effect on stdout
- B. To `log.txt`; `dup2(fd, 1)` closes the old fd 1 and makes fd 1 a copy of `fd`, so both refer to `log.txt`; writing to stdout now writes to the file
- C. To both the terminal and `log.txt`
- D. Nowhere; closing `fd` breaks stdout

**Answer:** B

**Why:** `dup2(oldfd, newfd)` makes `newfd` refer to the **same open-file entry** as `oldfd`, first closing `newfd` if it was open. After `dup2(fd, 1)`, **fd 1 points at `log.txt`**. `printf` writes to stdout (fd 1), so the text goes to the file. Closing the original `fd` is harmless, because fd 1 still refers to the file.

**Why the other options fail:**
- A: `dup2` is exactly what redirects stdout.
- C: stdout now points **only** at the file.
- D: fd 1 still refers to the file after `close(fd)`; they're independent copies pointing at the same entry.

**Concept recap:** `dup2(fd, 1)` = redirect stdout to whatever `fd` is open on. This is how a shell implements `> file`: in the child, after `fork` and before `exec`, it opens the file and `dup2`s it onto fd 1.

**Trap alert:** Keep stdio buffering in mind. Without the `fflush` (or a newline to a line-buffered terminal), buffered output could be lost or reordered relative to unbuffered writes (Q18, Q32).

**Interview angle:** "How does `command > file` work?" `fork`, then in the child `open(file)`, `dup2(fd, 1)`, `close(fd)`, then `exec`. Follow-up: *"How do you redirect stderr to stdout (`2>&1`)?"* `dup2(1, 2)`, after stdout is already set up.

---

## Q278. [M11 | Building a pipe with dup2 | Hard | MCQ | OA + Interview]

A shell runs `ls | wc -l`. After `pipe(fd)` gives a read end `fd[0]` and a write end `fd[1]`, which setup is correct for the two children?

- A. `ls` child: `dup2(fd[1], 1)`; `wc` child: `dup2(fd[0], 0)`; both children close both original pipe fds; the parent closes both pipe fds too
- B. `ls` child: `dup2(fd[0], 0)`; `wc` child: `dup2(fd[1], 1)`
- C. Both children share one fd and no `dup2` is needed
- D. The parent keeps both pipe ends open forever

**Answer:** A

**Why:** Data flows `ls → pipe → wc`:
- `ls` writes its output into the pipe, so its **stdout (1)** must become the pipe's **write end**: `dup2(fd[1], 1)`.
- `wc` reads its input from the pipe, so its **stdin (0)** must become the pipe's **read end**: `dup2(fd[0], 0)`.

Every process must **close the pipe fds it doesn't use**, including the originals after `dup2`. Otherwise `wc` never sees end-of-file: a `read` on a pipe returns 0 (EOF) only when **all** write ends are closed (Q42). If the parent or `ls` leaves a write end open, `wc` blocks forever.

**Why the other options fail:**
- B: Reversed. It would make `ls` read from the pipe and `wc` write to it.
- C: Without `dup2`, `ls` and `wc` would still read/write their terminal, not the pipe.
- D: Leaving ends open prevents EOF, so `wc` hangs.

**Concept recap:** Pipe setup = `pipe` → `fork` twice → in each child `dup2` the right end onto 0 or 1 → **close all unused pipe fds** → `exec`.

**Trap alert:** The single most common bug is **forgetting to close unused pipe ends**, which hangs the reader because it never sees EOF.

**Interview angle:** "Implement `ls | wc -l` with system calls." Walk through `pipe`, `fork`, `dup2`, and the closes. Follow-up: *"Why must every write end be closed before the reader sees EOF?"* As long as any process could still write, the kernel can't declare the pipe finished.

---

## Q279. [M11 | File syscalls and the fd table | Medium | MSQ | OA]

**Select all that apply.**

- A. `open()` returns the lowest-numbered unused file descriptor.
- B. `read(fd, buf, n)` returns the number of bytes actually read, which may be fewer than `n`; it returns 0 at end of file.
- C. `close(fd)` frees that fd number for reuse and decrements the open-file entry's reference count.
- D. After `close(fd)`, using `fd` again is valid and refers to the same file.
- E. `lseek(fd, 0, SEEK_END)` moves the file offset to the end of the file.

**Answer:** A, B, C, E

**Why:**
- **A:** The lowest-free-fd rule (Q272) underpins redirection tricks.
- **B:** A **short read** is normal, especially from pipes, sockets and terminals. 0 means EOF.
- **C:** Closing drops this fd's reference. The open-file entry is freed only when its count hits 0 (e.g., after all `dup`s and the child's copy are closed).
- **E:** `lseek` repositions the offset; `SEEK_END` with offset 0 goes to the end.

**Why the other options fail:**
- D: After `close`, the fd number is invalid until reused by a new `open`/`dup`. Using it gives `EBADF`.

**Concept recap:** The fd lifecycle: `open`/`dup` allocates, `read`/`write`/`lseek` use it, `close` frees it. The offset lives in the shared open-file entry.

**Trap alert:** Always loop on `read`/`write` until you've transferred everything, because short counts are legal and common.

**Interview angle:** "What does `read` returning 3 when you asked for 100 mean?" You got 3 bytes; call again for more. Follow-up: *"When does `write` do a short write?"* On a full pipe, a signal interrupting it, or a nearly full disk.

---

## Q280. [M11 | Close-on-exec | Hard | MCQ | OA + Interview]

A server opens a sensitive config file (fd 5), then `fork`s and `exec`s an untrusted helper program. Without special care, the helper inherits fd 5. How is this prevented?

- A. It can't be prevented; `exec` always keeps all fds
- B. By setting the **close-on-exec** flag on fd 5 (via `O_CLOEXEC` at open, or `fcntl(fd, F_SETFD, FD_CLOEXEC)`), so the fd is automatically closed when the process calls `exec`
- C. By calling `fork()` twice
- D. By marking the file read-only

**Answer:** B

**Why:** Open fds normally survive `exec` (Q27), which is what makes redirection work. But that also means a child program inherits fds it shouldn't see. The **close-on-exec** flag marks an fd to be closed automatically at `exec`, so it's gone before the untrusted program's code runs.

**Why the other options fail:**
- A: Close-on-exec is precisely the mechanism to close specific fds at `exec`.
- C: Double-fork changes parentage (Q31), not fd inheritance.
- D: Read-only still lets the helper read a sensitive file.

**Concept recap:** `O_CLOEXEC` (set at `open` time) avoids a race where another thread could `fork`+`exec` between the `open` and a later `fcntl`. Modern code sets `O_CLOEXEC` by default and clears it only for fds meant to be inherited.

**Trap alert:** This is distinct from inheritance across **`fork`** (Q275): `fork` always copies fds, and close-on-exec affects only **`exec`**. An fd can be inherited by the forked child and still be closed when that child `exec`s.

**Interview angle:** "How do you stop a child program from inheriting your open files?" Open them with `O_CLOEXEC`. Follow-up: *"Why prefer `O_CLOEXEC` over a later `fcntl`?"* It avoids a window where another thread's `fork`+`exec` could leak the fd before you set the flag.

---

## Answer key (quick check)

| Q | Ans | Q | Ans | Q | Ans | Q | Ans |
|---|---|---|---|---|---|---|---|
| 261 | B | 266 | A, B, C, D | 271 | B | 276 | B |
| 262 | B | 267 | B | 272 | B | 277 | B |
| 263 | A, B, C, D | 268 | A, B, C, D | 273 | A, B, C, E | 278 | A |
| 264 | C | 269 | B | 274 | B | 279 | A, B, C, E |
| 265 | B | 270 | A, C, D, E | 275 | B | 280 | B |

---

**Concepts tested so far (to avoid repeats):** Batches 1–13, plus: authN vs authZ · access matrix (domains × objects) · ACLs (by column) vs capabilities (by row), fds as capabilities · setuid as a domain switch · least privilege · memory protection as hardware-enforced (U/S, NX, per-process tables) · containers (namespaces + cgroups) vs VMs (hypervisor) · meta-rights (copy/owner/control) protecting the matrix · privilege escalation (vertical vs horizontal) and least-privilege defence · acquire-use-drop-sandbox pattern, seccomp · mechanism vs policy · fds 0/1/2 · Linux one-task model and `clone` · fork+exec+wait in a shell · fd inheritance sharing the offset · shared-offset write trace (AAABBB) · `dup2` redirection · pipe construction and closing unused ends · fd lifecycle and short reads · close-on-exec
*(M10 is complete. M11 finishes in Batch 15 with signals, pipes recap, the CFS/EEVDF scheduler, user vs kernel space, the syscall transition, `/proc`, and containers vs VMs, followed by M12 integrated questions, Q300, and the final deliverables.)*

**Completed: 280 / 300    Remaining: 20**
