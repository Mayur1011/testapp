# C/C++ MCQ/MSQ Bank — Batch 10 (Q181–Q200)

**What this batch covers:**
- **M11 stdio, finishing (Q181–Q184):**
  - what `scanf` returns when the input doesn't match
  - the `char c = getchar()` vs `EOF` trap
  - binary file modes with `fseek`/`ftell`/`rewind`
  - stdout vs stderr buffering
- **M12 C++ fundamentals (Q185–Q200):**
  - references vs pointers, and the rules for references
  - `const` member functions
  - the four C++ casts: choosing the right one, and which ones compile
  - `typeid` with polymorphic vs non-polymorphic types
  - brace initialisation and narrowing
  - `auto` deduction, including `auto x = {…}`
  - `constexpr` functions
  - range-for and structured bindings
  - `0` vs `nullptr` in overload resolution
  - `malloc`/`free` vs `new`/`delete`
  - `new[]` paired with `delete` (UB)
  - the dangling reference from `std::max` on temporaries
  - `malloc` paired with `delete` (UB)

**Patterns:** P26, P43, P44, P51, P52, P54, P55, P56, P78, P81, P82, P100, P101, P97 (secondary)

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - UB questions describe typical behaviour only in general terms.

> **Conventions for every question in this batch:**
> - **C questions (Q181–Q184):** C11. **C++ questions (Q185–Q200):** **C++17**.
> - Platform: a typical 64-bit system: `int` 4, `long` 8, pointers 8 bytes; plain `char` signed (unless the question says otherwise); ASCII.
> - Needed headers are included. C++ code writes `std::` explicitly.
> - Exceptions are enabled. Allocations in well-defined questions succeed.

---

## Q181. [M11 | What scanf returns on a mismatch | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

The program reads the input line `12 abc 7`.

```c
#include <stdio.h>
int main(void) {
    int a = 0, b = 0;
    int r = scanf("%d %d", &a, &b);
    printf("%d %d %d\n", r, a, b);
    return 0;
}
```

What does this program print?

- A. `2 12 7`
- B. `2 12 0`
- C. `1 12 0`
- D. `0 0 0`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** `scanf` reads the input from left to right, filling in variables **until something doesn't fit**. It returns **how many variables it actually filled**. "abc" isn't a number, so `scanf` stops there: only `a` got a value, and `b` keeps its old value.

**Why:**

| Step | Input consumed | Result |
|---|---|---|
| `%d` | `12` | a = 12 ✓ (1 item assigned) |
| ` ` (space in the format) | skips the whitespace | — |
| `%d` | sees `a`, not a digit → **matching failure** | stop; b unchanged (0) |
| return | | **1** |

`abc 7` is left unread in the input buffer. Prints `1 12 0`.

**Why the other options fail:**
- A, B: Think `scanf` skips bad input and carries on. It stops at the first mismatch.
- D: The first conversion succeeded.

**Concept recap:** `scanf` returns:
- the number of items **assigned**;
- `EOF` if the input ends (or an error happens) **before** the first conversion.

The standard check is `if (scanf("%d %d", &a, &b) != 2) { /* handle bad input */ }`.

**Trap alert:** The unread `abc` is still in the buffer. A loop like `while (scanf("%d", &x) != 1) { }` spins forever on the same bad characters. You must discard the rest of the line (e.g. read until `'\n'`) before retrying.

**Interview angle:** *"How do you read input safely in C?"* "Read a whole line with `fgets`, then parse it with `sscanf` or `strtol`, checking the return value and the end pointer. That gives full control over bad input and no leftover characters."

---

## Q182. [M11 | Storing getchar's result in a char | Hard | MSQ | [C] · C11 | All] · Type T3 · P97

```c
#include <stdio.h>
int main(void) {
    char c;
    int count = 0;
    while ((c = getchar()) != EOF)
        count++;
    printf("%d\n", count);
    return 0;
}
```

Select **all** true statements.

- A. Where plain `char` is **signed**, an input byte `0xFF` makes the loop stop early, as if the input had ended.
- B. Where plain `char` is **unsigned**, the loop never ends, even after the real end of input.
- C. Declaring `int c;` instead fixes both problems.
- D. `getchar` returns a `char`.
- E. The loop is correct on every platform as long as the input is plain ASCII text.

**Answer:** A, B, C

**In simple words:** `getchar` returns **one of 257 possible results**: the 256 byte values (0–255) or the special "no more input" value `EOF` (which is −1). A `char` can only hold 256 different values, so squeezing the result into a `char` loses information:
- If `char` is signed, the byte 255 turns into −1 and looks just like `EOF`.
- If `char` is unsigned, `EOF` turns into 255 and never matches −1, so the loop can't see the end.

**Why each option:**
- **A ✓:** On the stated platform (signed `char`), the byte 0xFF is stored as −1 (an implementation-defined conversion), which compares equal to `EOF`. Checked: the input `ab\xffcd` gives a count of **2**.
- **B ✓:** With unsigned `char`, `EOF` (−1) stored into `c` becomes 255. Then `c` is promoted to `int` 255, which is never equal to −1, so the loop is infinite.
- **C ✓:** An `int` holds all 257 possible results, so EOF and every byte stay distinct.
- **D ✗:** `getchar` returns `int` (N1570 §7.21.7.6).
- **E ✗:** It fails on unsigned-char platforms (option B) even with pure ASCII input.

**Concept recap:** The correct idiom: `int c; while ((c = getchar()) != EOF) { … }`. Convert to `char` only *after* the EOF check.

**Trap alert:** This exact bug appears in GATE and OA code-reading questions as "why does this loop hang or stop early?". The root cause is always the type of `c`.

**Interview angle:** *"Why does getchar return int?"* "So that it can return every possible byte value *and* a distinct EOF value. A `char` can't hold 257 distinct values."

---

## Q183. [M11 | Binary file update with fseek/ftell | Medium | MCQ | [C] · C11 | GATE+OA] · Type T5

```c
#include <stdio.h>
int main(void) {
    FILE *f = fopen("t.bin", "w+b");
    if (!f) return 1;
    fputs("HELLO", f);
    long a = ftell(f);
    fseek(f, 1, SEEK_SET);  int c = fgetc(f);
    fseek(f, -1, SEEK_END); int d = fgetc(f);
    rewind(f);              fputc('J', f);
    fseek(f, 0, SEEK_SET);
    char buf[8]; fgets(buf, sizeof buf, f);
    printf("%ld %c %c %s\n", a, c, d, buf);
    fclose(f); remove("t.bin");
    return 0;
}
```

What does this program print?

- A. `4 H O JELLO`
- B. `5 E L HELLO`
- C. `6 E O JELLO`
- D. `5 E O JELLO`

**Answer:** D

**Behaviour class:** Well-defined. Binary mode allows arbitrary `fseek` offsets, and every switch between reading and writing goes through `fseek` or `rewind`, as update mode requires.

**In simple words:** An open file has a **cursor**, like the text cursor in an editor. Writing and reading both happen at the cursor and move it forward. `ftell` asks "where is the cursor?", `fseek` moves it, and `rewind` puts it back at the start. Writing over existing text replaces characters, like the "overwrite" mode in an editor.

**Why:**

| Step | Action | Cursor | File | Value |
|---|---|---|---|---|
| `fopen("w+b")` | create/truncate, read + write, binary | 0 | (empty) | |
| `fputs("HELLO")` | write 5 bytes | 5 | HELLO | |
| `ftell` | | | | a = 5 |
| `fseek(1, SEEK_SET)`, `fgetc` | read index 1 | 2 | | c = 'E' |
| `fseek(-1, SEEK_END)`, `fgetc` | read the last byte | 5 | | d = 'O' |
| `rewind`, `fputc('J')` | overwrite index 0 | 1 | **J**ELLO | |
| `fseek(0)`, `fgets` | read the whole line | | | buf = "JELLO" |

**Why the other options fail:**
- A: Counts 0-based positions as lengths, and reads the wrong index.
- B: Misses the overwrite, and reads the wrong final byte.
- C: Counts a terminator that `fputs` doesn't write.

**Concept recap:**

| Mode | Meaning |
|---|---|
| `"r"` | read; the file must exist |
| `"w"` | write; create or **truncate** |
| `"a"` | append; every write goes to the end |
| `"+"` | read and write ("update") |
| `"b"` | binary (no newline translation) |

In update mode you must call `fseek`, `fsetpos` or `rewind` (or `fflush` after writing) between a write and a following read, and vice versa (N1570 §7.21.5.3p7).

**Trap alert:** In **text** mode, the standard only guarantees `fseek` with offset 0, or with a value previously returned by `ftell` (§7.21.9.2p4). Byte arithmetic like `-1, SEEK_END` is portable only in binary mode, which is why the question uses `"w+b"`.

**Interview angle:** *"How do you get a file's size in C?"* "`fseek(f, 0, SEEK_END); long n = ftell(f);` in **binary** mode. For very large files, use `fstat`/`stat` (POSIX) or `ftello` instead."

---

## Q184. [M11 | stdout vs stderr buffering | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

Assume the usual behaviour when output goes to an interactive terminal: **stdout is line-buffered, and stderr is unbuffered.**

```c
#include <stdio.h>
int main(void) {
    printf("A");
    fprintf(stderr, "B");
    printf("C\n");
    return 0;
}
```

What appears on the terminal?

- A. `ABC`
- B. `BAC`
- C. `ACB`
- D. `CAB`

**Answer:** B

**Behaviour class:** Well-defined given the stated buffering. The standard says stderr is **not fully buffered**, and stdout is fully buffered **only if** it isn't connected to an interactive device; the exact mode otherwise is implementation-defined (N1570 §7.21.3p7).

**In simple words:** stdout is like a **mailbox that's emptied when a full line is ready**. stderr is like **shouting**: it goes out immediately. "A" waits in the mailbox, "B" is shouted straight away, then "C\n" completes the line, and the mailbox ("AC") is delivered.

**Why:**

| Call | Stream | Buffer | Screen |
|---|---|---|---|
| `printf("A")` | stdout | "A" (held) | |
| `fprintf(stderr, "B")` | stderr | — (written immediately) | `B` |
| `printf("C\n")` | stdout | "AC\n" → newline → flushed | `BAC` |

**Why the other options fail:**
- A: Assumes both streams are unbuffered.
- C: Assumes stderr is buffered and stdout isn't.
- D: There's no reordering *within* stdout.

**Concept recap:** Ways to force output out now:
- `fflush(stdout)`;
- end the line (when line-buffered);
- `setvbuf(stdout, NULL, _IONBF, 0)` to disable buffering.

**Trap alert:** If stdout is **redirected to a file or a pipe**, it usually becomes fully buffered, so "AC" may appear only when the program exits. Debug prints that "disappear" before a crash were often still sitting in the buffer. Print debug messages to stderr.

**Interview angle:** *"Why did my printf output vanish when the program crashed?"* "It was still in stdout's buffer. A crash (e.g. a segfault) doesn't flush buffers. Use stderr, call `fflush`, or end lines with `\n` when writing to a terminal."

---

## Q185. [M12 | References vs pointers | Easy | MCQ | [C++] · C++17 | OA] · Type T5 · P51

```cpp
#include <iostream>
int main() {
    int a = 10, b = 20;
    int &r = a;
    int *p = &a;
    r = b;
    p = &b;
    *p = 30;
    std::cout << a << ' ' << b << ' ' << r << '\n';
}
```

What does this program print?

- A. `20 30 20`
- B. `10 30 30`
- C. `20 20 20`
- D. `30 30 30`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** A **reference is a permanent nickname**. `r` *is* `a` forever, so `r = b` means "set a to b's value", not "make r refer to b". A **pointer is a sign with an arrow**, and you can turn the arrow to point somewhere else (`p = &b`).

**Why:**

| Line | Statement | a | b | r refers to | p → |
|---|---|---|---|---|---|
| 3–5 | init | 10 | 20 | a | a |
| 6 | `r = b` → assigns to **a** | **20** | 20 | a | a |
| 7 | `p = &b` → re-points p | 20 | 20 | a | **b** |
| 8 | `*p = 30` | 20 | **30** | a | b |
| 9 | print a, b, r (= a) | | | | `20 30 20` |

**Why the other options fail:**
- B: Thinks `r = b` rebinds `r` to `b`, so `a` would stay 10 and `r` would read 30.
- C: Misses `*p = 30`.
- D: Treats everything as aliases of everything.

**Concept recap:**

| | Reference | Pointer |
|---|---|---|
| must be initialised | yes | no |
| can be re-pointed | **never** | yes |
| can be null | no (no null references in valid code) | yes |
| syntax | used like the object | `*p`, `p->` |
| arithmetic | none | yes |

**Trap alert:** "Assigning to a reference" always writes to the object it refers to. There's no syntax for rebinding a reference.

**Interview angle:** *"When do you use a reference instead of a pointer?"* "References for parameters and return values that must refer to a valid object (`const T&` for read-only, `T&` for output). Pointers when 'no object' (null) is a valid state, when the target must change, or for raw ownership in low-level code (prefer smart pointers)."

---

## Q186. [M12 | Rules for references | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P51

Select **all** true statements (C++17).

- A. A reference must be initialised when it is declared.
- B. A reference can later be reseated to refer to a different object.
- C. A well-defined program has no "null references".
- D. `int &r = 5;` is valid.
- E. `const int &r = 5;` is valid, and it extends the lifetime of the temporary holding 5 to match `r`.

**Answer:** A, C, E

**In simple words:** A reference is a nickname, and a nickname needs **someone to belong to** from the start. A non-`const` nickname must name a real, changeable object, not a bare number like 5. A `const` nickname may name a temporary value, and C++ then **keeps that temporary alive** for as long as the nickname exists.

**Why each option:**
- **A ✓:** `int &r;` is ill-formed.
- **B ✗:** There's no rebinding syntax (Q185).
- **C ✓:** Creating a reference by dereferencing a null pointer (`int &r = *p;` with `p` null) is already UB, so a valid program never contains one.
- **D ✗:** A non-const lvalue reference can't bind to an rvalue (a temporary). It's ill-formed.
- **E ✓:** A `const T&` **can** bind to a temporary, and binding *directly* to it extends the temporary's lifetime to the reference's scope ([class.temporary]p6).

**Concept recap:**

| Reference type | Binds to lvalues (named objects) | Binds to rvalues (temporaries) |
|---|---|---|
| `T&` | yes | no |
| `const T&` | yes | yes (lifetime extended) |
| `T&&` | no | yes (used for move semantics, M18) |

**Trap alert:** Lifetime extension works **only** when the reference binds *directly* to the temporary. If the temporary passes through a function that returns a reference, there's no extension, and you get a dangling reference (Q199).

**Interview angle:** *"Why can `const T&` bind to temporaries but `T&` can't?"* "Modifying a temporary that's about to disappear would almost always be a bug, so the language forbids it for non-const references. Read-only binding is safe, and it lets functions taking `const std::string&` accept `"literals"` without extra overloads."

---

## Q187. [M12 | Calling members on a const object | Easy | MCQ | [C++] · C++17 | OA] · Type T2 · P52

```cpp
#include <iostream>
struct Counter {
    int n = 0;
    int get() const { return n; }
    void inc() { ++n; }
};
void show(const Counter &c) {
    std::cout << c.get();   // line 1
    c.inc();                // line 2
}
```

Which line(s) fail to compile?

- A. Only line 1
- B. Only line 2
- C. Both
- D. Neither

**Answer:** B

**Behaviour class:** Ill-formed (line 2).

**In simple words:** Marking a member function `const` is a **promise: "I only look, I don't touch."** When you hold an object through a `const` reference, you may only call member functions that made that promise. `inc` changes `n`, so it isn't `const`, and calling it is refused.

**Why:**

| Call | Function | `this` inside it | Allowed on a `const Counter&`? |
|---|---|---|---|
| `c.get()` | `int get() const` | `const Counter*` | ✓ |
| `c.inc()` | `void inc()` | `Counter*` | ✗ "passing 'const Counter' as 'this' argument discards qualifiers" |

**Why the other options fail:** see the table.

**Concept recap:** A `const` member function can't modify ordinary data members, and can only call other `const` member functions. Mark **every** member function that doesn't change the object as `const`. Otherwise it can't be used through `const` references, which are the most common way objects are passed.

**Trap alert:** `mutable` members (like a cache or a mutex) may still change inside `const` functions. That's covered in M13.

**Interview angle:** *"What does `const` after a member function mean?"* "`this` becomes a pointer to const. The function can't modify the object's (non-mutable) state, and it's callable on const objects. It also enables overloading on constness, as `operator[]` does for `const` and non-`const` access."

---

## Q188. [M12 | Choosing the right cast | Easy | MCQ | [C++] · C++17 | OA+Interview] · Type — concept · P54

Match each task to the most appropriate cast:

1. Convert a `double` to an `int`.
2. Remove `const` from a pointer to an object that was **not** originally defined `const`, to call a legacy API.
3. Convert a `Base*` to a `Derived*`, **checking at run time** whether it really points to a `Derived`.
4. Reinterpret an object's address as an integer (e.g. a `std::uintptr_t`).

- A. 1 `const_cast`, 2 `static_cast`, 3 `reinterpret_cast`, 4 `dynamic_cast`
- B. 1 `static_cast`, 2 `reinterpret_cast`, 3 `dynamic_cast`, 4 `const_cast`
- C. 1 `static_cast`, 2 `const_cast`, 3 `dynamic_cast`, 4 `reinterpret_cast`
- D. 1 `dynamic_cast`, 2 `const_cast`, 3 `static_cast`, 4 `reinterpret_cast`

**Answer:** C

**In simple words:** C++ splits the one C cast into **four labelled tools**, so the code says *what* you're doing, and the compiler refuses the dangerous mixes:
- **static_cast:** an ordinary, sensible conversion.
- **const_cast:** add or remove `const`.
- **dynamic_cast:** a safe downcast with a run-time check.
- **reinterpret_cast:** "treat these bits as something else" (the most dangerous).

**Why:**

| Cast | What it does | Example |
|---|---|---|
| `static_cast` | well-defined conversions known at compile time | `static_cast<int>(3.7)` → 3 |
| `const_cast` | only adds or removes `const`/`volatile` | `const_cast<char*>(s)` |
| `dynamic_cast` | polymorphic down/cross casts with a run-time type check: returns `nullptr` (for pointers) or throws `std::bad_cast` (for references) on failure | `dynamic_cast<Derived*>(b)` |
| `reinterpret_cast` | low-level reinterpretation of bits or addresses | `reinterpret_cast<std::uintptr_t>(p)` |

**Why the other options fail:** A, B and D each put at least one cast on a job it can't do (`const_cast` can't convert `double` → `int`; `dynamic_cast` doesn't apply to `double`; and so on).

**Concept recap:** A C-style cast `(T)x` tries `const_cast`, then `static_cast`, then `reinterpret_cast` (and combinations) until one works. So it can silently do something dangerous. Named casts make that visible.

**Trap alert:** `dynamic_cast` requires the source type to be **polymorphic** (it must have at least one virtual function). Otherwise it doesn't compile for downcasts.

**Interview angle:** *"Explain the four C++ casts."* Give one line each, plus "prefer `static_cast`; `const_cast` only for legacy APIs; `dynamic_cast` signals a design smell if overused (prefer virtual functions); `reinterpret_cast` only for low-level code".

---

## Q189. [M12 | Which casts compile? | Hard | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P54

```cpp
struct B { virtual ~B() {} };
struct D : B {};
const int ci = 5;
int i = 3;
double d = 2.5;
B *pb = new D;
```

Select **all** lines that **compile**.

- A. `int x = static_cast<int>(d);`
- B. `int *p = static_cast<int *>(&ci);`
- C. `D *pd = dynamic_cast<D *>(pb);`
- D. `long *pl = static_cast<long *>(&i);`
- E. `int *q = const_cast<int *>(&ci);`

**Answer:** A, C, E

**In simple words:** Each named cast is allowed to do **only its own job**. `static_cast` refuses to remove `const`, and refuses to treat an `int*` as an unrelated `long*`, because those are jobs for `const_cast` and `reinterpret_cast`. That refusal is the point: you can't do something dangerous by accident.

**Why each option:**
- **A ✓:** A standard numeric conversion → 2.
- **B ✗:** `static_cast` can't cast away `const` ("invalid static_cast from 'const int*' to 'int*'").
- **C ✓:** `B` is polymorphic (it has a virtual destructor), so `dynamic_cast` is allowed. It succeeds here (`pb` really points to a `D`).
- **D ✗:** `int*` → `long*` are unrelated pointer types. Only `reinterpret_cast` would accept this, and then accessing through it would break strict aliasing (Q118-E).
- **E ✓:** It **compiles**. But `ci` was *defined* `const`, so **writing** through `q` would be UB (the same rule as Q100). Reading through it is fine.

**Concept recap:** Compiling isn't the same as being safe. `const_cast` compiling only means the types fit. Whether a later write is legal depends on how the object was **defined**.

**Trap alert:** A failed `dynamic_cast` on a **reference** throws `std::bad_cast`. On a **pointer** it returns `nullptr`. Always check the pointer result.

**Interview angle:** *"Can `const_cast` cause UB?"* "Yes, if you then modify an object that was originally declared const. The cast compiles, but the write is undefined. It's only safe when the underlying object isn't const."

---

## Q190. [M12 | typeid: polymorphic vs non-polymorphic | Hard | MCQ | [C++] · C++17 | Interview] · Type T1 · P55

```cpp
#include <iostream>
#include <typeinfo>
struct B { virtual ~B() {} };
struct D : B {};
struct NB {};
struct ND : NB {};
int main() {
    D d; ND nd;
    B *pb = &d;
    NB *pn = &nd;
    std::cout << (typeid(*pb) == typeid(D))
              << (typeid(*pn) == typeid(ND))
              << (typeid(pb) == typeid(B *)) << '\n';
}
```

What does this program print?

- A. `111`
- B. `011`
- C. `100`
- D. `101`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** `typeid` can only look up an object's **real** (run-time) type if the class carries type information, which it does only if it has at least one **virtual function** (it's "polymorphic"). For a plain class, `typeid` just reports the **declared** type. And `typeid` of a *pointer* always reports the pointer's declared type, never the pointed-to object's.

**Why:**

| Expression | Is the class polymorphic? | typeid gives | Compared with | Result |
|---|---|---|---|---|
| `typeid(*pb)` | B: **yes** (virtual destructor) | the **dynamic** type → D | D | **1** |
| `typeid(*pn)` | NB: **no** | the **static** type → NB | ND | **0** |
| `typeid(pb)` | it's a pointer | `B*` | `B*` | **1** |

Prints `101`.

**Why the other options fail:**
- A: Assumes `typeid` always finds the dynamic type.
- B: Has the polymorphic and non-polymorphic cases backwards.
- C: Thinks `typeid(pb)` looks through the pointer.

**Concept recap:** **RTTI** (run-time type information) is stored alongside the vtable, so only polymorphic classes have it. `dynamic_cast` uses the same information.

**Trap alert:** `typeid(*p)` where `p` is a **null** pointer to a polymorphic type throws `std::bad_typeid`. `typeid(x).name()` returns an implementation-defined (often mangled) string, so never compare names: compare `typeid` objects with `==`, or use `std::type_index`.

**Interview angle:** *"How does `dynamic_cast` know the real type?"* "Through RTTI reached via the object's vptr. That's why it only works on polymorphic types, and why it has a run-time cost."

---

## Q191. [M12 | Brace initialisation and narrowing | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T2 · P56

```cpp
int a = 3.7;      // L1
int b{3.7};       // L2
int c(3.7);       // L3
double d{3};      // L4
char e{300};      // L5
```

Which lines are **ill-formed**? (`char` is 8 bits.)

- A. L2 and L5
- B. L1, L2 and L5
- C. Only L2
- D. L2, L4 and L5

**Answer:** A

**Behaviour class:** L2 and L5 are ill-formed (a diagnostic is required). L1, L3 and L4 are well-formed (a = 3, c = 3, d = 3.0).

**In simple words:** Curly-brace initialisation is the **strict** form: it refuses any conversion that could **lose information**, called "narrowing". The old `=` and `()` forms quietly chop 3.7 down to 3. Braces say "no, that loses data". But a constant that fits exactly (3 → 3.0) is fine.

**Why:**

| Line | Form | Conversion | Narrowing? | Verdict |
|---|---|---|---|---|
| L1 | `=` | double → int | yes, but not checked | ✓ a = 3 |
| L2 | `{}` | double → int | **yes** | ✗ ill-formed |
| L3 | `()` | double → int | not checked | ✓ c = 3 |
| L4 | `{}` | int **constant** 3 → double | exactly representable → not narrowing | ✓ |
| L5 | `{}` | int constant 300 → char | **doesn't fit** in −128…127 | ✗ ill-formed |

**Why the other options fail:**
- B: L1 uses `=`, which allows narrowing (as C does).
- C: Misses L5: narrowing also applies to constants that don't fit.
- D: L4's constant fits exactly, so it isn't narrowing.

**Concept recap:** Narrowing conversions ([dcl.init.list]):
- floating → integer;
- a wider floating type → a narrower one (unless it's a constant that fits);
- integer → floating (unless it's a constant that fits exactly);
- a wider integer → a narrower integer (unless it's a constant that fits).

**Trap alert:** `std::vector<int> v{3.7};` also fails to compile, for the same reason. Brace initialisation is a cheap way to catch accidental truncation.

**Interview angle:** *"Why prefer `{}` initialisation?"* "It works uniformly for every type, forbids narrowing, avoids the 'most vexing parse' (`T x();` declares a function), and value-initialises with `T x{};`. Caveat: for types with an `initializer_list` constructor, `{}` prefers that constructor, so `vector<int>{5}` has one element, while `vector<int>(5)` has five."

---

## Q192. [M12 | auto deduction | Hard | MSQ | [C++] · C++17 | OA+Interview] · Type T3 · P78

```cpp
const int ci = 5;
int x = 1;
int &rx = x;
auto a = ci;
auto b = rx;
auto &c = ci;
auto d = {1, 2, 3};
auto e{7};
```

Select **all** true statements (C++17).

- A. `a` is `int`.
- B. `b` is `int&`.
- C. `c` is `const int&`.
- D. `d` is `std::initializer_list<int>`.
- E. `e` is `std::initializer_list<int>`.

**Answer:** A, C, D

**In simple words:** Plain `auto` makes a **fresh copy**. A copy is just a value, so it drops the "reference" and any top-level "const" of what it copies. If you write `auto&`, you ask for a nickname instead, and a nickname to a const object must stay const. `= {…}` with braces is special-cased to make a list. Since C++17, `auto e{7}` with a single value is simply an `int`.

**Why** (each option was checked with `static_assert(std::is_same_v<…>)`):

| Variable | Rule | Deduced type |
|---|---|---|
| `a = ci` | by value: top-level `const` dropped | **int** ✓ |
| `b = rx` | by value: the reference is dropped (it copies `x`) | **int** (B ✗) |
| `c` (`auto&` bound to `ci`) | a reference keeps the referent's `const` | **const int&** ✓ |
| `d = {1, 2, 3}` | copy-list-initialisation with `auto` | **std::initializer_list<int>** ✓ |
| `e{7}` | C++17 rule (N3922): a single element in direct-list-init gives that element's type | **int** (E ✗) |

**Concept recap:** `auto` follows **template argument deduction** rules (as for `template<class T> void f(T)`), except for the braced-list special case. `decltype(expr)` instead keeps references and `const` exactly.

**Trap alert:** `for (auto x : container)` copies each element. Use `auto&` to modify elements, and `const auto&` to avoid copies of large elements.

**Interview angle:** *"What's the type of `auto x = {1, 2};` vs `auto x{1};`?"* "`std::initializer_list<int>` vs `int` (since C++17). And `auto x{1, 2};` is ill-formed." That's a favourite trick question.

---

## Q193. [M12 | constexpr functions | Easy | MCQ | [C++] · C++17 | OA] · Type T1 · P81

```cpp
#include <iostream>
constexpr int fact(int n) { return n <= 1 ? 1 : n * fact(n - 1); }
int main() {
    constexpr int f5 = fact(5);
    int arr[fact(4)];
    int n = 3;
    int r = fact(n);
    static_assert(fact(3) == 6);
    std::cout << f5 << ' ' << sizeof(arr) / sizeof(arr[0]) << ' ' << r << '\n';
}
```

What does this program print?

- A. Compile error: `fact(n)` with a non-constant `n`
- B. `120 24 6`
- C. `120 4 6`
- D. Compile error: an array size can't come from a function call

**Answer:** B

**Behaviour class:** Well-defined. `arr`'s elements are never read; it's only measured with `sizeof`.

**In simple words:** `constexpr` on a function means **"this function *may* be run by the compiler"**. If you call it with constant arguments where a constant is needed (array sizes, `static_assert`, `constexpr` variables), the compiler computes the answer at build time. If you call it with a run-time value, it simply runs normally at run time.

**Why:**

| Use | Context | When it's evaluated | Value |
|---|---|---|---|
| `constexpr int f5 = fact(5)` | requires a constant | compile time | 120 |
| `int arr[fact(4)]` | an array bound requires a constant | compile time | 24 elements |
| `fact(n)` with a variable `n` | an ordinary call | run time | 6 |
| `static_assert(fact(3) == 6)` | compile time | ✓ passes | |

**Why the other options fail:**
- A: A `constexpr` function may also be called at run time. Only a C++20 **`consteval`** function would reject `fact(n)`.
- C: Confuses `fact(4)` (24) with 4.
- D: A constant-expression function call is a valid array bound.

**Concept recap:**
- `constexpr` **variable** → it *must* be computed at compile time.
- `constexpr` **function** → it *can* be computed at compile time.
- `consteval` (C++20) → it *must* be computed at compile time.

Compile-time computation means zero run-time cost.

**Trap alert:** Since C++14, `constexpr` functions can contain loops and local variables, not just a single `return` (that was C++11's restriction).

**Interview angle:** *"`const` vs `constexpr`?"* "`const` means 'I won't modify this', and its value may only be known at run time. `constexpr` means 'this value is known at compile time', so it's usable in array sizes, template arguments and `static_assert`."

---

## Q194. [M12 | Range-for & structured bindings | Easy | MCQ | [C++] · C++17 | OA] · Type T8 · P82

```cpp
#include <iostream>
#include <map>
#include <string>
#include <vector>
int main() {
    std::map<std::string, int> m{{"b", 2}, {"a", 1}, {"c", 3}};
    int sum = 0;
    std::string keys;
    for (const auto &[k, v] : m) { keys += k; sum += v; }
    std::vector<int> v{1, 2, 3};
    for (auto x : v) x *= 10;
    for (auto &x : v) x += 1;
    std::cout << keys << ' ' << sum << ' ' << v[0] << v[1] << v[2] << '\n';
}
```

What does this program print?

- A. `bac 6 234`
- B. `abc 6 112131`
- C. `abc 6 234`
- D. `abc 6 11213141`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:**
- `[k, v]` **unpacks** each key-value pair into two names, like opening an envelope and taking out two cards.
- A `std::map` always keeps its keys **sorted**, so the loop sees a, b, c.
- `for (auto x : v)` works on **copies**: changing `x` changes nothing in `v`. `auto &x` works on the **real elements**.

**Why:**

| Loop | Effect |
|---|---|
| over `m` (sorted by key) | keys = "abc", sum = 1 + 2 + 3 = 6 |
| `auto x` (copy) | `v` unchanged: {1, 2, 3} |
| `auto &x` (reference) | v = {2, 3, 4} |

Prints `abc 6 234`.

**Why the other options fail:**
- A: Uses insertion order. `std::map` is ordered by key (`std::unordered_map` would have no particular order).
- B, D: Think the copy loop modified `v`.

**Concept recap:** Structured bindings (C++17) work with pairs, tuples, arrays and simple structs: `auto [x, y] = point;`. Range-for works with anything that has `begin()`/`end()`.

**Trap alert:** `for (auto [k, v] : m)` (no `&`) copies each pair, including the `std::string` key. Use `const auto&` for read-only loops over containers of large objects.

**Interview angle:** *"How do you iterate a map and modify the values?"* `for (auto &[k, v] : m) v *= 2;`. The key is still `const` (you can't change keys in place), but the value is modifiable.

---

## Q195. [M12 | 0 vs nullptr in overload resolution | Easy | MCQ | [C++] · C++17 | OA] · Type T1 · P26

```cpp
#include <iostream>
void f(int)    { std::cout << "int "; }
void f(char *) { std::cout << "ptr "; }
int main() {
    f(0);
    f(nullptr);
    std::cout << '\n';
}
```

What does this program print?

- A. `ptr ptr`
- B. `int int`
- C. Compile error: `f(0)` is ambiguous
- D. `int ptr`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** `0` is first and foremost **a number**. It can *also* mean "null pointer", but that's a second-best match, so `f(int)` wins. `nullptr` is a special value that **can only mean "null pointer"**, so it always picks the pointer version. That's exactly why C++11 introduced it.

**Why:**

| Call | Argument type | `f(int)` | `f(char*)` | Chosen |
|---|---|---|---|---|
| `f(0)` | `int` | exact match | null-pointer conversion | **f(int)** |
| `f(nullptr)` | `std::nullptr_t` | no conversion exists | pointer conversion | **f(char*)** |

**Why the other options fail:**
- A: Thinks 0 means a pointer.
- B: Thinks `nullptr` converts to `int`. It doesn't.
- C: An exact match beats a conversion, so there's no ambiguity.

**Concept recap:** `NULL` in C++ is an implementation-defined null pointer constant: it may be `0`, `0L`, or a compiler built-in. So `f(NULL)` may call `f(int)`, or may be **ambiguous**, depending on the implementation. Always write `nullptr` in C++.

**Trap alert:** Templates make this worse: `template<class T> void g(T x) { f(x); }` called as `g(0)` deduces `T = int`, and `f(int)` is called even if you "meant" null.

**Interview angle:** *"Why was nullptr added?"* "`0` and `NULL` are integers, so they pick integer overloads and deduce as `int` in templates. `nullptr` has its own type (`std::nullptr_t`) that converts only to pointer types."

---

## Q196. [M12 | malloc/free vs new/delete | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P43

```cpp
#include <cstdlib>
#include <iostream>
struct T {
    T()  { std::cout << "C"; }
    ~T() { std::cout << "D"; }
};
int main() {
    T *a = new T;
    T *b = static_cast<T *>(std::malloc(sizeof(T)));
    delete a;
    std::free(b);
    std::cout << '\n';
}
```

What does this program print?

- A. `CD`
- B. `CCDD`
- C. `CDD`
- D. `CCD`

**Answer:** A

**Behaviour class:** Well-defined. `b`'s memory is never used as a `T` object (no `T` was constructed there), and it's released with the matching `free`.

**In simple words:** `new` does two jobs: **get memory, and build the object** (it runs the constructor). `delete` **tears the object down** (destructor) **and returns the memory**. `malloc`/`free` only handle the raw memory, like renting an empty plot of land without ever building a house on it.

**Why:**

| Step | Constructor? | Destructor? | Output |
|---|---|---|---|
| `new T` | ✓ | | `C` |
| `malloc(sizeof(T))` | ✗ (raw bytes only) | | |
| `delete a` | | ✓ | `D` |
| `free(b)` | | ✗ | |

Prints `CD`.

**Why the other options fail:**
- B, D: Assume `malloc` runs the constructor.
- C: Assumes `free` runs the destructor.

**Concept recap:**

| | `new` / `delete` | `malloc` / `free` |
|---|---|---|
| constructor / destructor | ✓ / ✓ | ✗ / ✗ |
| result type | `T*` | `void*` (needs a cast in C++) |
| on failure | throws `std::bad_alloc` | returns `NULL` |
| size | computed for you | you compute it |
| can be overloaded per class | yes | no |

**Trap alert:** Using `b` as a `T` (e.g. `b->member = 1;`) without constructing one first would be UB for a type with a non-trivial constructor. You'd need *placement new* (`new (b) T`) to build the object in that memory.

**Interview angle:** *"Can you mix malloc and delete?"* "No. Each allocation must be released by its matching function: `malloc`→`free`, `new`→`delete`, `new[]`→`delete[]`. Mixing them is UB (Q198, Q200)."

---

## Q197. [M12 | malloc vs new: the rules | Medium | MSQ | [C++] · C++17 | Interview] · Type T3 · P43, P101

Select **all** true statements (C++17).

- A. `new` calls the constructor; `malloc` does not.
- B. By default, `new` reports failure by throwing `std::bad_alloc`, while `malloc` returns a null pointer.
- C. Calling `delete` on memory from `malloc` is fine as long as the type is trivial (like `int`).
- D. With `new`, you must compute the byte size yourself using `sizeof`.
- E. `realloc` can safely resize an array created with `new[]`.

**Answer:** A, B

**In simple words:** `new`/`delete` and `malloc`/`free` are **two separate shops with separate receipts**. You must return what you bought to the shop you bought it from, even for simple things like an `int`. `new` also builds objects and works out sizes for you. `malloc` only hands over raw bytes.

**Why each option:**
- **A ✓:** See Q196.
- **B ✓:** `new (std::nothrow) T` is the variant that returns `nullptr` instead.
- **C ✗:** It's UB **regardless of the type** ([expr.delete]: the pointer must have come from `new`). The two systems may use different heaps or bookkeeping (Q200).
- **D ✗:** `new T[n]` computes the size itself.
- **E ✗:** `realloc` knows only about `malloc`-family memory. It also does a raw byte copy, which is wrong for objects that manage resources. Use `std::vector`, which grows correctly.

**Concept recap:** In modern C++, write neither `malloc` nor raw `new`/`delete` in application code. Use containers (`std::vector`, `std::string`) and smart pointers (`std::make_unique`), which call the right functions automatically.

**Trap alert:** "It worked when I mixed them" happens because many C++ runtimes implement `new` on top of `malloc`. That's an implementation detail, not a guarantee, and it breaks as soon as a custom allocator or class-specific `operator new` is involved.

**Interview angle:** *"Name differences between new and malloc."* Five bullets: constructors and destructors, type safety, failure handling, size computation, overloadability. Then: "In modern C++ I'd use `make_unique` or a container instead of either."

---

## Q198. [M12 | new[] released with delete | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T3 · P44

> **This code has undefined behaviour.**

```cpp
#include <iostream>
struct T { ~T() { std::cout << "D"; } };
int main() {
    T *p = new T[3];
    delete p;
    std::cout << '\n';
}
```

Which statement is correct?

- A. It prints `DDD`.
- B. It prints `D` (only the first element is destroyed), and that's the defined behaviour.
- C. The behaviour is undefined: memory from `new[]` must be released with `delete[]`.
- D. Compile error: `delete` can't be applied to an array pointer.

**Answer:** C

**Behaviour class:** **Undefined behaviour.**

**In simple words:** `new T[3]` builds **three** objects and secretly notes "there are 3 here", usually just before the array. `delete[]` reads that note and destroys all three. Plain `delete` thinks there's **one** object, so it reads the wrong bookkeeping. The result is undefined: it may destroy only one object, corrupt the heap, or crash.

> **⚠ UB box**
> - **Category:** mismatched array new/delete.
> - **Rule:** [expr.delete]p2: for a non-array `delete`, the operand must be a pointer to a non-array object created by `new`; for `delete[]`, it must come from an array `new`. Otherwise, the behaviour is undefined.
> - **What you may see in practice:** often a single `D` followed by a crash or a heap-corruption error, because the allocator receives a pointer that is offset from what it handed out. **Not guaranteed.**

**Why the other options fail:**
- A: That's what `delete[] p;` would print.
- B: One possible symptom, presented as defined behaviour.
- D: The compiler can't tell from `T*` whether it points to one object or to an array, so no diagnostic is required. (Some compilers can warn in obvious cases.)

**Concept recap:** Matching pairs:

| Allocation | Release |
|---|---|
| `new T` | `delete p` |
| `new T[n]` | `delete[] p` |
| `malloc` / `calloc` / `realloc` | `free` |

**Trap alert:** For trivial types like `int`, the mismatch often *appears* to work. It's still UB. The bug becomes visible with types that have destructors.

**Interview angle:** *"How do you avoid new[]/delete mismatches?"* "Don't write `new[]` at all: use `std::vector<T>` or `std::unique_ptr<T[]>`, which always call `delete[]`."

---

## Q199. [M12 | Dangling reference from std::max on temporaries | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T3 · P100

> **This code has undefined behaviour.**

```cpp
#include <algorithm>
#include <iostream>
int main() {
    const int &r = std::max(3, 7);
    std::cout << r << '\n';
}
```

Which statement is correct?

- A. It prints `7`, guaranteed, because binding to `const int&` extends the temporary's lifetime.
- B. The behaviour is undefined. `std::max` returns a reference to one of its **temporary** arguments, and those temporaries are destroyed at the end of the full expression. So `r` dangles; lifetime extension doesn't pass through a function's return value.
- C. Compile error: `std::max` can't take temporaries.
- D. It prints `3`.

**Answer:** B

**Behaviour class:** **Undefined behaviour.**

**In simple words:** `std::max(3, 7)` receives two **temporary** values that exist only until the end of the line. It returns a **nickname for one of them** (the 7), not a copy. The "keep the temporary alive" rule only applies when you bind a nickname *directly* to a temporary. Here the nickname came back through a function, so the 7 disappears at the `;`, and `r` names something that no longer exists.

> **⚠ UB box**
> - **Category:** dangling reference (a reference to a temporary whose lifetime has ended).
> - **Rule:** [class.temporary]p6: lifetime extension applies only when a reference is bound directly to a temporary. A temporary bound to a reference **parameter** lasts only until the end of the full-expression containing the call. `std::max` is declared `const T& max(const T&, const T&)`.
> - **What you may see in practice:** often `7`, because the stack slot hasn't been reused yet. With optimisation, or under AddressSanitizer (stack-use-after-scope), it can print garbage or be reported as an error. **Not guaranteed.**

**Why the other options fail:**
- A: The lifetime-extension rule (Q186-E) is real, but it doesn't apply through a function's return value.
- C: `const T&` parameters bind to temporaries just fine.
- D: The returned reference refers to the larger value, 7.

**Concept recap:** Fixes:
- copy the result: `int r = std::max(3, 7);`;
- or use `auto r = std::max(...)`, which deduces a non-reference type.

The same trap appears with `std::min` and `std::clamp`, and with any function that returns a reference to one of its parameters.

**Trap alert:** Returning a reference to a **local** variable (`const std::string& f() { std::string s; return s; }`) is the same bug in its simplest form. Compilers warn about that one, but usually not about the `std::max` case.

**Interview angle:** *"Is `const auto& x = std::max(a, b + 1);` safe?"* "No. `b + 1` is a temporary, and if it's the larger value, `x` dangles. Store the result by value. When in doubt, return and store by value."

---

## Q200. [M12 | Releasing malloc'd memory with delete | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T3 · P101

```cpp
#include <cstdlib>
int main() {
    int *p = static_cast<int *>(std::malloc(sizeof(int)));
    *p = 5;
    delete p;
}
```

Which statement is correct?

- A. It's fine, because `int` is trivial and has no destructor.
- B. It's a memory leak but otherwise well-defined.
- C. Compile error: `delete` can't be applied to memory from `malloc`.
- D. The behaviour is undefined: memory from `malloc` must be released with `free`.

**Answer:** D

**Behaviour class:** **Undefined behaviour.**

**In simple words:** This is returning an item to the wrong shop (Q197). Even though an `int` needs no clean-up, `delete` hands the memory to C++'s deallocation function (`operator delete`), which only promises to work with memory that came from `operator new`.

> **⚠ UB box**
> - **Category:** mismatched allocation and deallocation functions (`malloc` + `delete`).
> - **Rule:** [expr.delete]p2: the operand of `delete` must be a pointer obtained from a non-array `new` expression (or null). Otherwise, the behaviour is undefined.
> - **What you may see in practice:** often nothing visible, when the runtime happens to implement `operator new` on top of `malloc`. With a replaced `operator new`, a custom allocator or a debug heap, it can crash or be reported as an allocation mismatch (AddressSanitizer reports "alloc-dealloc-mismatch"). **Not guaranteed.**

**Why the other options fail:**
- A: The rule is about **which allocator** owns the memory, not about destructors.
- B: The memory is released, through the wrong function, which is worse than a leak.
- C: The compiler can't know where a pointer came from, so no diagnostic is required.

**Concept recap:** The reverse mistake, `free(new int)`, is UB too, for the same reason.

**Trap alert:** This often sneaks in when C libraries return `malloc`'d buffers to C++ code (e.g. `strdup`, `getline`). Release them with `free`, or wrap them in `std::unique_ptr<char, decltype(&std::free)>`.

**Interview angle:** *"How do you hold a malloc'd buffer safely in C++?"* `std::unique_ptr<char, void(*)(void*)> p(strdup(s), std::free);`. The custom deleter guarantees the matching release on every path.

---

## Batch 10 summary

| Q | Module | Lang | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|---|
| 181 | M11 | C | E | MCQ | C | Well-defined |
| 182 | M11 | C | H | MSQ | A, B, C | Impl-defined consequences |
| 183 | M11 | C | M | MCQ | D | Well-defined |
| 184 | M11 | C | E | MCQ | B | Well-defined (stated buffering) |
| 185 | M12 | C++ | E | MCQ | A | Well-defined |
| 186 | M12 | C++ | M | MSQ | A, C, E | — (rules) |
| 187 | M12 | C++ | E | MCQ | B | Ill-formed (line 2) |
| 188 | M12 | C++ | E | MCQ | C | — (concept) |
| 189 | M12 | C++ | H | MSQ | A, C, E | Ill-formed / OK mix |
| 190 | M12 | C++ | H | MCQ | D | Well-defined |
| 191 | M12 | C++ | M | MCQ | A | Ill-formed (L2, L5) |
| 192 | M12 | C++ | H | MSQ | A, C, D | — (deduction) |
| 193 | M12 | C++ | E | MCQ | B | Well-defined |
| 194 | M12 | C++ | E | MCQ | C | Well-defined |
| 195 | M12 | C++ | E | MCQ | D | Well-defined |
| 196 | M12 | C++ | M | MCQ | A | Well-defined |
| 197 | M12 | C++ | M | MSQ | A, B | — (rules) |
| 198 | M12 | C++ | M | MCQ | C | **UB** |
| 199 | M12 | C++ | H | MCQ | B | **UB** |
| 200 | M12 | C++ | M | MCQ | D | **UB** |

- **Difficulty:** 8 E / 7 M / 5 H.
  - M11 total (Q178–Q184): 3/3/1, as planned.
  - M12 total: 6/6/4, as planned.
- **Format:** 15 MCQ / 5 MSQ.
  - M11 total: 6/1, as planned.
  - M12 total: 12/4, as planned.
- **MCQ answer keys:** A×3, B×4, C×4, D×4. Running total Q1–Q200: **A 40 · B 40 · C 40 · D 40**.

---

**Completed: 200 / 400 Remaining: 200**

**Ledger, newly ✓ this batch:**
- P26 ×1 (Q195)
- P43 ×2 (Q196, Q197) ★ complete
- P44 ×1 (Q198)
- P51 ×2 (Q185, Q186) ★ complete
- P52 ×1 (Q187) ★ — 1 more planned in M13
- P54 ×2 (Q188, Q189) ★ complete
- P55 ×1 (Q190)
- P56 ×1 (Q191)
- P78 ×1 (Q192) ★ — 2 more planned in M22
- P81 ×1 (Q193)
- P82 ×1 (Q194)
- P100 ×1 (Q199)
- P101 ×1 (Q200)
- P97 → 2 (Q182)

**M11 complete** (Q178–Q184, 7 questions). **M12 complete** (Q185–Q200, 16 questions). **All eleven C modules are done.** The C++ OOP family begins at Q201.

---
---

# SELF-AUDIT #2 (after Batch 10, Q1–Q200)

## A. Actual vs planned splits

**Difficulty**

| | Easy | Medium | Hard |
|---|---|---|---|
| Actual, Q1–Q200 | 54 (27%) | 82 (41%) | 64 (32%) |
| Plan for exactly the modules covered (M1–M12) | 54 | 82 | 64 |
| Bank-wide target | 25% | 40% | 35% |

✅ An exact match with the per-module plan. The remaining modules (M14, M16, M18, M25) are planned Hard-heavy, so the bank finishes at 100 / 160 / 140.

**Format:** MCQ 160 / MSQ 40 = exactly 80 / 20. ✅

**Language:** [C] 184 / [C++] 16. ✅ On plan. The remaining C questions are 11 (M24: 6, M25: 5), for a final total of 195 C / 205 C++.

**MCQ answer keys:** A 40 · B 40 · C 40 · D 40. ✅ "All/none of the above" used: **0** (limit 5).

**Code-reading share:** about 84% of questions require reading code (the target was ≥70%). ✅

## B. Relevance-tag split — ⚠ deviation found

| Tag | C so far (184) | C target (195; revised after Audit #1) | C++ so far (16) |
|---|---|---|---|
| [GATE] | 29 | 38 | — |
| [GATE+OA] | 85 | 102 | — |
| [All] | 56 | 45 | — |
| [OA+Interview] | 14 | 5 | 9 |
| [OA] | 0 | 2 | 5 |
| [Interview] | 0 | 3 | 2 |

**What happened:** systems topics in M6–M10 (strict aliasing, alignment, flexible array members, unions vs C++, linking, macro double evaluation) naturally earned interview tags, and memory-safety questions earned [All]. Only 11 C questions remain, so the original C targets for [GATE] and [GATE+OA] **can't be met without padding**, which the rules forbid.

**Decision and revised final targets (whole bank):**
- **GATE relevance remains very strong.** 170 of 184 C questions (92%) carry a GATE-inclusive tag ([GATE], [GATE+OA] or [All]).
- Revised bank-wide targets: [GATE] ≈ 33 · [GATE+OA] ≈ 92 · [All] ≈ 56 · [OA] ≈ 64 · [OA+Interview] ≈ 110 · [Interview] ≈ 45 (total 400).
- The remaining 11 C questions (M24 and the C part of M25) will be tagged **[GATE] or [GATE+OA] only**: linked lists, trees, stacks and queues with well-defined code tracing, which matches the recent GATE emphasis.

## C. Pointer + OOP share

- M6 = 35 of 200 so far (17.5%).
- The OOP family M13–M18 (110 questions) starts now. The planned final share is 145 / 400 = **36.25%**. ✅ On schedule.

## D. Full Pattern Ledger (✓ = covered, with the count)

| Range | Status |
|---|---|
| **P01–P15** | all ✓ — P01 2, P02 5, P03 4, P04 2, P05 4, P06 2, P07 2, P08 1, P09 2, P10 2, P11 4, P12 6, P13 2, P14 3, P15 3 |
| P16 overload resolution | ☐ → M17 |
| **P17–P50** | all ✓ — P17 3, P18 5, P19 7, P20 4, P21 3, P22 3, P23 4, P24 3, P25 4, P26 1, P27 3, P28 3, P29 3, P30 2, P31 2, P32 9, P33 4, P34 2, P35 2, P36 1, P37 3, P38 2, P39 1, P40 1, P41 1, P42 5, P43 2, P44 1, P45 4, P46 2, P47 1, P48 2, P49 2, P50 2 |
| P51 reference vs pointer ★ | ✓ 2 |
| P52 const correctness ★ | ✓ 1 — **needs 1 more** (M13) |
| P53 explicit ★ | ☐ → M14 |
| P54 four casts ★ | ✓ 2 |
| P55 typeid | ✓ 1 |
| P56 narrowing | ✓ 1 (+1 planned in M14) |
| P57–P77 OOP / move | ☐ → M13–M18 |
| P78 auto ★ | ✓ 1 — **needs 1 more** (M22) |
| P79, P80 lambdas | ☐ → M21/M22 |
| P81 constexpr | ✓ 1 |
| P82 structured bindings / range-for | ✓ 1 |
| P83–P96 templates, STL, concurrency | ☐ → M20–M23, M25 |
| P97 char signedness | ✓ 2 |
| P98 sizeof implementation-defined | ✓ 1 |
| P99 static init fiasco | ☐ → M14 |
| P100 dangling reference | ✓ 1 |
| P101 malloc/delete | ✓ 1 |
| P102–P105 integrated | ☐ → M20/M25 |

**Totals:** **61 / 105 patterns ✓.** ★ patterns: **28 / 51 complete**. Every C-side ★ pattern is complete. Two C++ ★ patterns are partly done (P52 1/2, P78 1/2); 21 C++ ★ patterns are scheduled in M13–M25. Nothing scheduled for Batches 1–10 is missing.

## E. Re-trace of every output; corrections

- Every well-defined output question in Batches 6–10 was compiled and run (C as C11, C++ as C++17), and stdin/file questions were run with the stated input. **All answer keys match.**
- Every compile-error question was confirmed by compiling: Q115, Q146, Q187, Q189, Q191.
- Q174 was confirmed to fail at link time.
- **Corrections since Audit #1** (all were made before delivery):
  - Q111, Q139: distractor explanations rewritten so that each wrong option has a consistent derivation. The keys were unchanged.
  - **The user's request to remove compiler- and version-specific content** was applied to Batches 1–5 (conventions, UB boxes, and Q67's answer wording). No answer key changed.
- **Found in this audit:** none needing a fix.

## F. ⚠ Highest-priority check: any deterministic output claimed for UB code?

- **UB questions in Batches 6–10**, all with "undefined" as the correct option and a general "what you may see" note:
  - single-answer: Q109, Q117, Q119, Q127, Q154, Q156, Q169, Q180, Q198, Q199, Q200;
  - bug-fix questions whose original code is UB: Q151, Q159, Q160;
  - UB options inside multiple-answer questions: Q115, Q129, Q155, Q161.
- **Code that looks risky but was re-checked as well-defined:**
  - Q122: only the safe function's result is used.
  - Q162: no double free in the code shown.
  - Q190: rewritten to use stack objects, so there's no `delete` through a non-virtual base.
  - Q193: the array is only measured with `sizeof`.
  - Q196: the raw `malloc` block is never used as a `T`.
  - Q143–Q144: union punning is valid in C (and UB in C++, as stated).

**Result: 0 violations.**

---

# PROGRESS REPORT @ 200

**Modules complete:** M1–M12 — every C module plus C++ fundamentals. **Next:** M13–M18, the OOP family (Q201–Q310).

**Key traps covered in Q101–Q200:**
1. `**pp` vs `*pp + 1` vs `**++pp` on an array of pointers (Q104).
2. Modifying a string literal is UB, but `char a[] = "…"` is writable (Q109). Identical literals may or may not share an address (Q121).
3. Strict aliasing: reading a float through `int *` is UB; use `memcpy` (Q117–Q118).
4. `strncpy` may leave no terminator (Q128). `snprintf` returns the length it *wanted* to write (Q130). `strtok` overwrites your string (Q131).
5. Struct padding and trailing padding, member order, `offsetof` (Q132–Q140).
6. Union punning: allowed in C, UB in C++ (Q143–Q144). Bit-fields wrap, and have no address or `sizeof` (Q145–Q146).
7. `realloc` always ends the old object's lifetime, even at the same address (Q156). Freeing a linked list requires reading `next` before `free` (Q159).
8. Macros: `SQ(1+2)`, `100/SQ(5)`, a `MAX` without parentheses, `SQ(i++)` (UB), and `MAX(next(), 0)` calling `next` twice (well-defined) (Q166–Q170).
9. Header guards don't prevent multiple-definition link errors across files (Q174).
10. `char c = getchar()` breaks EOF detection (Q182). stderr output can appear before earlier stdout output (Q184).
11. A reference can't be rebound (Q185). `auto x = {…}` vs `auto x{…}` (Q192). `std::max` on temporaries returns a dangling reference (Q199).

**Misconceptions tested:** "union = struct" · "`strncpy` is a safe `strcpy`" · "header guards fix linking" · "`printf` converts arguments to match the format" · "`realloc` in place keeps old pointers valid" · "`const&` always extends lifetime" · "`delete` vs `delete[]` doesn't matter for simple types" · "`static_cast` can remove const".

**UB categories covered so far (21):** the 12 from Audit #1, plus:
13. modifying a string literal
14. strict-aliasing violation
15. `strlen` / `%s` on an unterminated buffer
16. heap buffer overflow
17. use-after-free (including through an alias and after `realloc`)
18. double free / freeing a non-heap pointer
19. `printf` format mismatch
20. `new[]` + `delete`, and `malloc` + `delete`
21. a dangling reference from returning a reference to a temporary

**Still to come:** `delete` through a base with a non-virtual destructor, pure-virtual calls during construction, data races, iterator invalidation, lambda capture lifetimes, the static-initialisation-order fiasco.

**Areas to revise:** Q104 (pointer-to-pointer arithmetic), Q112 (reading declarations), Q133/Q140 (padding), Q156 (realloc lifetime), Q171 (two-level stringify), Q174 (one-definition rule across files), Q190 (typeid with non-polymorphic types), Q192 (auto deduction).

*Say **"Continue."** for Batch 11 (Q201–Q220: M13 classes — class vs struct, `this`, static members, `const`/`mutable`, friends, object size, nested classes — then M14 constructors begins).*
