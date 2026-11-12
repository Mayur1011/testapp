# C/C++ MCQ/MSQ Bank — Batch 5 (Q81–Q100)

**What this batch covers:**
- **M5 Arrays, finishing (Q81–Q84):**
  - out-of-bounds writes and reads (UB), and which pointer expressions near the array end are still legal
  - `sizeof` on a VLA *is* evaluated
  - arrays can't be assigned, but structs that contain arrays can
- **M6 Pointers, part 1 of 2 (Q85–Q100):**
  - null, wild and `void` pointers
  - the first line that causes UB in a pointer program
  - basic address-of and dereference traces
  - pointer arithmetic scaled by the pointed-to type, and byte-level access (endianness)
  - pointer differences and which pointer operations are legal
  - comparing pointers to different objects
  - the `*p++` / `(*p)++` / `++*p` / `*++p` family, including the `while (*d++ = *s++)` copy loop
  - `const int *` vs `int *const`
  - casting away `const` (sometimes legal, sometimes UB)

**Patterns:** P05, P17, P18, P28, P29

> **Conventions for every question in this batch (unless the question says otherwise):**
> - C11, `gcc -std=c11`, 64-bit Linux / GCC 13 / x86-64, **little-endian**
> - `sizeof(int)=4`, `sizeof(double)=8`, `sizeof(void*)=8`; plain `char` signed; two's complement; ASCII
> - Addresses in memory tables are illustrative; only the differences between them matter.
> - Every well-defined program was compiled and run, and its output matches the key.
> - For UB cases, observed GCC behaviour is reported but never claimed as guaranteed.

---

## Q81. [M5 | Writing one past the end | Easy | MCQ | [C] · C11 | All] · Type T3 · P28

> **This code has undefined behaviour.**

```c
#include <stdio.h>
int main(void) {
    int a[5] = {1, 2, 3, 4, 5};
    int b = 99;
    a[5] = 0;
    printf("%d %d\n", a[4], b);
    return 0;
}
```

Which statement is correct?

- A. It prints `5 0`, guaranteed, because `b` sits directly after `a` in memory.
- B. It prints `5 99`, guaranteed, because the write lands in unused padding.
- C. The behaviour is undefined: `a[5]` writes outside the array.
- D. Compile error: index 5 is out of range.

**Answer:** C

**Behaviour class:** **Undefined behaviour.**

> **⚠ UB box**
> - **Category:** out-of-bounds access (writing one past the end of an array).
> - **Rule:** N1570 §6.5.6p8. `a + 5` is a legal "one past the end" pointer, but **dereferencing** it, which `a[5] = …` does, is undefined.
> - **On GCC 13 x86-64** this printed `5 99` at both `-O0` and `-O2` (verified). GCC happened to place `b` *below* `a` on the stack, so the stray write hit some other stack slot. With a different variable order, compiler or flag, it could overwrite `b`, corrupt a saved register or return address, or trigger the stack protector (`*** stack smashing detected ***`). **No result is guaranteed, and none may be relied on.**

**Why:** C does no bounds checking. `a[5]` is just `*(a + 5)`, a write to whatever memory follows the array. The language gives that write no meaning at all.

**Why the other options fail:**
- A, B: Each describes one possible memory layout. The standard doesn't fix how locals are laid out, and the write itself is UB regardless of what it hits.
- D: The compiler isn't required to diagnose this. GCC may warn (`-Warray-bounds`, mainly with optimisation on), but must still accept the code.

**Concept recap:** Valid indices are `0 … n−1`. Forming the pointer `a + n` is allowed, but `a[n]` is not.

**Trap alert:** "It printed the right answer, so it works" is the most dangerous belief in C. Many security exploits (stack buffer overflows) start from exactly this kind of off-by-one write.

**Interview angle:** *"How do you detect out-of-bounds bugs?"* "AddressSanitizer (`-fsanitize=address`), which adds red zones around arrays; `-D_FORTIFY_SOURCE=2` for library calls; and in C++, `std::array::at()` / `vector::at()`, which throw `std::out_of_range`. Static analysers (`-fanalyzer`, clang-tidy) catch some cases too."

---

## Q82. [M5 | Which accesses near the end are UB? | Hard | MSQ | [C] · C11 | All] · Type T3 · P28

Given `int a[5];` with every element initialised, select **all** expressions whose **evaluation** is undefined behaviour.

- A. `a[4]`
- B. `a[5]` (reading its value)
- C. `&a[5]`
- D. `a + 6` (only computing the pointer, not dereferencing it)
- E. `a[-1]`

**Answer:** B, D, E

**Behaviour class per option:** A = well-defined · B = UB · C = **well-defined** · D = UB · E = UB

**Why each option:**
- **A ✗:** The last valid element.
- **B ✓ UB:** It reads one past the end.
- **C ✗ (well-defined):** N1570 §6.5.3.2p3 says `&a[5]` is evaluated **as if** it were `a + 5`. The implied `*` is never evaluated, so this produces the legal one-past-the-end pointer. That's why `for (p = a; p != &a[5]; p++)` is correct C.
- **D ✓ UB:** Just **computing** a pointer more than one past the end is already UB (§6.5.6p8), even if it is never dereferenced.
- **E ✓ UB:** A pointer before the first element (`a − 1`) is not allowed either.

**Concept recap:** The legal pointer values for an n-element array are `a + 0 … a + n`. Of those, only `a + 0 … a + n − 1` may be dereferenced.

**Trap alert:** Reverse loops written as `for (p = a + n - 1; p >= a; p--)` compute `a − 1` on their final decrement, which is technically UB. The safe form is `for (p = a + n; p != a; ) { --p; … }`.

**Interview angle:** *"Why does the standard allow one-past-the-end at all?"* "So that half-open ranges `[begin, end)` work. That's the basis of C++ iterators, `std::end(arr)`, and every `for (p = begin; p != end; ++p)` loop."

---

## Q83. [M5 | sizeof on a VLA evaluates its operand | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P05

```c
#include <stdio.h>
int main(void) {
    int n = 3;
    int vla[n];
    size_t s1 = sizeof(vla);
    size_t s2 = sizeof(int[n++]);
    size_t s3 = sizeof(n++);
    printf("%zu %zu %zu %d\n", s1, s2, s3, n);
    return 0;
}
```

What does this program print?

- A. `12 12 4 4`
- B. `12 16 4 5`
- C. `12 12 4 3`
- D. `12 16 4 4`

**Answer:** A

**Behaviour class:** Well-defined. The elements of `vla` are never read.

**Why:** Normally `sizeof` does not evaluate its operand. The **one exception** is an operand whose type is a **variable-length array** type: its size is only known at run time, so the expression that gives the size must be evaluated (N1570 §6.5.3.4p2).

| Line | Statement executed | Evaluated? | Variables after execution |
|---|---|---|---|
| 3 | `int n = 3;` | — | n = 3 |
| 4 | `int vla[n];` | the size is fixed at 3 when the declaration is reached | |
| 5 | `sizeof(vla)` | the VLA's size is computed at run time: 3 × 4 | s1 = 12 |
| 6 | `sizeof(int[n++])` | **VLA type → `n++` IS evaluated.** The size uses the old n (3): 3 × 4; then n becomes 4 | s2 = 12, n = 4 |
| 7 | `sizeof(n++)` | the type is `int`, not a VLA → **not evaluated** | s3 = 4, n = 4 |
| 8 | `printf` | — | `12 12 4 4` |

**Why the other options fail:**
- B: Uses the incremented n for the size (4 × 4 = 16) and also evaluates the second `n++`.
- C: Assumes `sizeof` never evaluates anything. True except for VLAs.
- D: Uses the post-increment value for the size. `n++` yields the **old** value.

**Concept recap:** Changing `n` after the declaration doesn't resize `vla`. Its length is fixed when its declaration is executed.

**Trap alert:** VLAs became an **optional** feature in C11 (implementations may define `__STDC_NO_VLA__`). They are not part of standard C++, though GCC accepts them in C++ as an extension. A large or attacker-controlled `n` can overflow the stack with no way to detect the failure.

**Interview angle:** *"Should I use VLAs?"* "Rarely. There's no failure reporting if `n` is too large, and they're banned in many codebases (the Linux kernel removed them; `-Wvla` warns). Use `malloc`, or a fixed maximum plus a check. VLA *types* in parameters (`int m[r][c]`) are still useful for multi-dimensional indexing."

---

## Q84. [M5 | Arrays are not assignable; structs are | Medium | MCQ | [C] · C11 | GATE+OA] · Type T2

```c
int main(void) {
    int a[3] = {1, 2, 3}, b[3];
    b = a;                                   /* line X */

    struct W { int v[3]; } s = {{1, 2, 3}}, t;
    t = s;                                   /* line Y */
    return 0;
}
```

Which lines compile?

- A. Both X and Y
- B. Only X
- C. Only Y
- D. Neither

**Answer:** C

**Behaviour class:** Line X is ill-formed. Line Y is well-defined.

**Why:**
- **Line X:** `b` is an array, and an array is **not a modifiable lvalue** (N1570 §6.3.2.1p1), so it can't appear on the left of `=`. (On the right, `a` would decay to a pointer anyway.) GCC: *"assignment to expression with array type"* (verified).
- **Line Y:** Struct assignment copies the whole struct, **including any arrays inside it**, member by member (effectively a `memcpy`). After it runs, `t.v` = {1, 2, 3}.

**Why the other options fail:**
- A: Assumes arrays can be copied with `=`, as `std::array` or Python lists can.
- B: Has the rules backwards.
- D: Assumes a struct containing an array can't be assigned either. It can.

**Concept recap:** Ways to copy an array: `memcpy(b, a, sizeof a);`, a loop, wrapping it in a struct, or in C++ using `std::array<int,3>`, which supports `=`.

**Trap alert:** The same "wrap it in a struct" trick lets a C function **return** or **take by value** a fixed-size array. A bare array parameter becomes a pointer (Q71), but a struct member array does not.

**Interview angle:** *"Why can't C arrays be assigned or returned?"* "For historical reasons. Arrays decay to pointers in almost every expression, so `a = b` would mean assigning an address to an array. C++ fixed this with `std::array`, which is a struct holding a C array, giving value semantics with zero overhead."

---

## Q85. [M6 | Null, wild and void pointers | Medium | MSQ | [C] · C11 | GATE+OA] · Type T3

Select **all** true statements (C11).

- A. After `int *p;` at block scope, `p` is a wild (uninitialised) pointer, and dereferencing it is undefined behaviour.
- B. `int *p = NULL; if (p) *p = 1;` is well-defined.
- C. After `static int *p;`, `p` is a null pointer.
- D. `int x; void *v = &x; *v = 5;` is valid C.
- E. A null pointer is guaranteed to be represented in memory as all-bits-zero.

**Answer:** A, B, C

**Why each option:**
- **A ✓:** An automatic pointer with no initializer holds an indeterminate value. It may point anywhere, so dereferencing it (and in fact even using its value) is UB.
- **B ✓:** The `if` reads `p` (null → false), so the dereference never runs. Checking a pointer against null is always well-defined.
- **C ✓:** Static-storage objects are zero-initialised, and for pointers that means a **null pointer** (N1570 §6.7.9p10). This holds even on unusual hardware where null is not all-bits-zero.
- **D ✗:** `void` is an incomplete type, so `*v` has type `void` and can't be assigned. GCC: *"invalid use of void expression"* (verified). Convert first: `*(int *)v = 5;`.
- **E ✗:** The standard says the **integer constant** `0`, used in a pointer context, produces a null pointer (§6.3.2.3p3). It says nothing about how a null pointer is stored. All-bits-zero is universal on mainstream hardware but not guaranteed, so `memset(&ptr, 0, sizeof ptr)` and `calloc` are not *portable* ways to get null pointers.

**Concept recap:**

| Kind | What it is | Safe to dereference? |
|---|---|---|
| Null | deliberately points to nothing | no (UB), but safe to test |
| Wild | never initialised | no; even reading its value is a bug |
| Dangling | pointed at an object whose lifetime has ended | no (UB) (Q64, M9) |
| `void *` | generic address with no element type | not until converted to an object pointer type |

**Trap alert:** `if (p = NULL)` (assignment instead of `==`) sets `p` to null and is always false. Writing it as `NULL == p` turns the typo into a compile error.

**Interview angle:** *"What happens when you dereference a null pointer?"* "In the language: UB. In practice, on Linux, page 0 is unmapped, so the program usually crashes with SIGSEGV. But the compiler may assume a dereferenced pointer is non-null and delete later null checks. That's a real class of kernel bugs, which is why the kernel builds with `-fno-delete-null-pointer-checks`."

---

## Q86. [M6 | Guarding a null pointer with short-circuiting | Easy | MCQ | [C] · C11 | GATE] · Type T1

```c
#include <stdio.h>
int main(void) {
    int *p = NULL;
    int x = 5;
    if (p != NULL && *p > 0) x = 1;
    if (p == NULL || *p > 0) x += 10;
    printf("%d\n", x);
    return 0;
}
```

What does this program print?

- A. `1`
- B. `11`
- C. Nothing guaranteed: `p` is dereferenced while null
- D. `15`

**Answer:** D

**Behaviour class:** Well-defined. `*p` is never evaluated.

**Why:**

| Line | Test | Short-circuit? | `*p` evaluated? | x |
|---|---|---|---|---|
| 5 | `p != NULL` → false | `&&` stops | **no** | 5 |
| 6 | `p == NULL` → true | `||` stops | **no** | 5 + 10 = 15 |

**Why the other options fail:**
- A, B: Run the first assignment, which requires the `&&` to succeed. It fails at the null check.
- C: Ignores short-circuiting. The dereference sits in the operand that is never evaluated.

**Concept recap:** Always put the null check **first** in an `&&` chain: `p && p->next && p->next->val`.

**Trap alert:** Writing the checks in the other order, `*p > 0 && p != NULL`, dereferences first, which is UB when `p` is null. Order matters.

**Interview angle:** *"Traverse a linked list safely in a condition."* `while (cur != NULL && cur->val != key) cur = cur->next;`. Explain that `&&` guarantees `cur->val` is read only when `cur` is non-null.

---

## Q87. [M6 | Which line first causes UB? | Hard | MCQ | [C] · C11 | All] · Type T3

```c
#include <stdio.h>
int main(void) {
    int *p;                    /* L1 */
    int x = 10;                /* L2 */
    int *q = &x;               /* L3 */
    *q = 20;                   /* L4 */
    p = q + 1;                 /* L5 */
    printf("%d\n", *(p - 1));  /* L6 */
    *p = 30;                   /* L7 */
    return 0;
}
```

Which line is the **first** to invoke undefined behaviour?

- A. L1: declaring a pointer without initialising it
- B. L5: making a pointer that points past `x`
- C. L6: dereferencing `p - 1`
- D. L7: writing through `p`

**Answer:** D

**Behaviour class:** UB at L7. Lines L1–L6 are well-defined, and L6 prints 20.

**Why:**

**Memory table** (say `x` is at 1000)

| Line | Statement | Variable | Stored value / Points to | Legal? |
|---|---|---|---|---|
| L1 | `int *p;` | p | indeterminate | ✓ declaring is fine; only *using* p would be UB |
| L2 | `int x = 10;` | x @1000 | 10 | ✓ |
| L3 | `q = &x` | q | 1000 → x | ✓ |
| L4 | `*q = 20` | x | 20 | ✓ |
| L5 | `p = q + 1` | p | 1004 (one past x) | ✓ A single object counts as an array of length 1 for pointer arithmetic (N1570 §6.5.6p7), so one-past-the-end is legal |
| L6 | `*(p - 1)` | — | 1000 → x = 20 | ✓ prints `20` |
| L7 | `*p = 30` | ? | dereferences one-past-the-end | **✗ UB** (§6.5.6p8) |

**Why the other options fail:**
- A: An uninitialised pointer is harmless until its value is used. L5 assigns it before any use.
- B: One past a single object is explicitly allowed. **Two** past (`q + 2`) would be UB.
- C: `p - 1` brings the pointer back to `x`, a valid object.

**Concept recap:** The pointer-arithmetic rules treat **every object as an array**, even a lone `int`. That's why `&x + 1` is a valid end pointer for a "range" of one element.

**Trap alert:** Since L7 is UB, the **whole program** has undefined behaviour. The compiler is not obliged to keep L6's output intact either. Asking for "the *first* UB line" is about locating the bug, not a promise about what runs before it.

**Interview angle:** *"Can you treat a single variable as a one-element array?"* "Yes. `int x; process(&x, 1);` with `process(int *a, size_t n)` iterating `a[0 … n−1]` is valid, and so is `&x + 1` as the end pointer."

---

## Q88. [M6 | Address-of and dereference trace | Easy | MCQ | [C] · C11 | GATE] · Type T5

```c
#include <stdio.h>
int main(void) {
    int x = 10, y = 20;
    int *p = &x, *q = &y;
    *p = *q + 5;
    p = q;
    *p = 7;
    printf("%d %d %d %d\n", x, y, *p, *q);
    return 0;
}
```

What does this program print?

- A. `25 20 25 20`
- B. `25 7 7 7`
- C. `7 20 7 20`
- D. `25 7 25 7`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:**

**Memory table** (x @1000, y @1004)

| Line | Statement | x | y | p | q |
|---|---|---|---|---|---|
| 3 | init | 10 | 20 | — | — |
| 4 | `p = &x, q = &y` | 10 | 20 | 1000 → x | 1004 → y |
| 5 | `*p = *q + 5` → x = 20 + 5 | **25** | 20 | → x | → y |
| 6 | `p = q` (copies the **address**) | 25 | 20 | **1004 → y** | → y |
| 7 | `*p = 7` → y = 7 | 25 | **7** | → y | → y |
| 8 | print x, y, *p (y), *q (y) | | | | `25 7 7 7` |

**Why the other options fail:**
- A: Assumes `p = q` copies the *value* (`*p = *q`) rather than the pointer.
- C: Has line 5 write to the wrong place and miss the final write to y.
- D: Assumes `p` still points to x after `p = q`.

**Concept recap:**
- `p = q;` → both pointers now point to the **same** object (aliasing).
- `*p = *q;` → copies the **value** from one object into another.

Being clear about which of these a line does is the whole of pointer basics.

**Trap alert:** After `p = q`, any write through `p` is visible through `q`. A write through one name that changes what another name reads is called **aliasing**, and it's what makes pointer code hard to analyse.

**Interview angle:** *"What does `p = q` vs `*p = *q` do?"* Draw the boxes-and-arrows picture. Interviewers want to see you separate a pointer's own value (an address) from the value it points to.

---

## Q89. [M6 | Precedence of `*` with `+` and `&*` | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int a[3] = {5, 10, 15};
    int *p = a;
    printf("%d %d %d %d\n", *p + 1, *(p + 1), *&a[2], (int)(&*p == p));
    return 0;
}
```

What does this program print?

- A. `6 10 15 1`
- B. `10 10 15 1`
- C. `6 10 15 0`
- D. `6 6 15 1`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:**

| Expression | Grouping | Value |
|---|---|---|
| `*p + 1` | unary `*` binds tighter than `+`: `(*p) + 1` = 5 + 1 | **6** |
| `*(p + 1)` | the pointer moves one int forward → a[1] | **10** |
| `*&a[2]` | `*` and `&` cancel → a[2] | **15** |
| `&*p == p` | `&*p` is exactly `p` (N1570 §6.5.3.2p3: neither operator is evaluated) | **1** |

**Why the other options fail:**
- B: Reads `*p + 1` as `*(p + 1)`.
- C: Thinks `&*p` produces a different pointer. It is defined to be identical. (`&*p` is even well-defined when `p` is null, because nothing is dereferenced.)
- D: Reads `*(p + 1)` as `*p + 1`.

**Concept recap:** Unary operators (`* & ++ -- ! ~ -` and casts) bind tighter than every binary operator. So `*p + 1` adds to the **value**, while `*(p + 1)` moves the **pointer**.

**Trap alert:** `*p++` is different again: postfix `++` binds *tighter* than prefix `*`. See Q96.

**Interview angle:** *"`*p + 1` vs `*(p + 1)`?"* "The first is the pointed-to value plus one. The second is the next element. The difference is precedence: unary `*` binds tighter than `+`."

---

## Q90. [M6 | Arithmetic scales by the pointed-to type | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P18

```c
#include <stdio.h>
int main(void) {
    double d[4];
    char *c = (char *)d;
    double *p = d;
    printf("%td %td\n", (char *)(p + 2) - c, (p + 3) - p);
    return 0;
}
```

What does this program print?

- A. `2 2`
- B. `16 16`
- C. `2 3`
- D. `16 3`

**Answer:** D

**Behaviour class:** Well-defined. Only addresses are computed (the elements are never read), and converting any object pointer to `char *` is allowed.

**Why:**

| Expression | Rule | Value |
|---|---|---|
| `p + 2` | moves by 2 × `sizeof(double)` = **16 bytes** | address d + 16 |
| `(char *)(p + 2) - c` | the difference of two `char *` pointers is counted in **bytes** | **16** |
| `(p + 3) - p` | the difference of two `double *` pointers is counted in **elements** | **3** |

**Why the other options fail:**
- A: Treats every pointer as moving by 1 byte, or measures both results in elements.
- B: Measures both results in bytes. Subtracting two `double *` pointers gives elements.
- C: Treats `p + 2` as moving 2 bytes.

**Concept recap:**
- `p + n` → address + n × `sizeof(*p)`.
- `q − p` → (address difference) / `sizeof(*p)`, as a `ptrdiff_t` (print with `%td`).
- Both pointers in a subtraction must point into (or one past) the **same** array.

**Trap alert:** To measure a byte distance between pointers of other types, cast **both** to `char *` (or `unsigned char *`) first, exactly as in the first expression.

**Interview angle:** *"Why can't you do arithmetic on `void *`?"* "The standard doesn't define `sizeof(void)`, so there's no step size. GCC allows it as an extension, treating the size as 1. Portable code casts to `char *`."

---

## Q91. [M6 | Byte-level view of ints (endianness) | Hard | MCQ | [C] · C11 | OA+Interview] · Type T1 · P18

Platform: x86-64, **little-endian** (the least significant byte is stored at the lowest address).

```c
#include <stdio.h>
int main(void) {
    int a[] = {0x01020304, 0x05060708};
    unsigned char *c = (unsigned char *)a;
    printf("%x %x\n", *(c + 1), *(c + 4));
    return 0;
}
```

What does this program print?

- A. `2 5`
- B. `4 8`
- C. `3 8`
- D. `3 5`

**Answer:** C

**Behaviour class:** Well-defined. Byte order is **implementation-defined**, and the stated little-endian platform fixes it. Reading any object through `unsigned char *` is explicitly allowed by the aliasing rules (N1570 §6.5p7).

**Why:**

**Memory table** (a @1000)

| Address | 1000 | 1001 | 1002 | 1003 | 1004 | 1005 | 1006 | 1007 |
|---|---|---|---|---|---|---|---|---|
| Byte (hex) | 04 | **03** | 02 | 01 | **08** | 07 | 06 | 05 |
| Belongs to | a[0], least significant byte → | | | most significant | a[1], least significant → | | | |

- `c + 1` → address 1001 → **0x03**
- `c + 4` → address 1004 → the first byte of a[1] → **0x08**

`%x` prints in hex without leading zeros → `3 8`.

**Why the other options fail:**
- A: The **big-endian** answer (bytes stored 01 02 03 04, so 1001 holds 02 and 1004 holds 05).
- B: Off by one byte in the first expression (it reads `c + 0`).
- D: Mixes the two byte orders.

**Concept recap:** Character pointers (`char`, `unsigned char`, `signed char`) may inspect the bytes of **any** object. Any other cross-type access, such as reading a `float` through an `int *`, breaks strict aliasing (P30, Batch 6).

**Trap alert:** Use `unsigned char` for raw bytes. With a signed `char`, a byte like 0xF0 would print as `fffffff0` under `%x` after sign-extension.

**Interview angle:** *"Write code to detect endianness at run time."* `int x = 1; return *(unsigned char *)&x == 1;` → 1 means little-endian. Follow-up: *"Why does it matter?"* Network protocols are big-endian ("network byte order"), so use `htonl`/`ntohl`, and binary file formats must specify their byte order.

---

## Q92. [M6 | Pointer arithmetic and differences | Medium | MCQ | [C] · C11 | GATE+OA] · Type T5 · P18

```c
#include <stdio.h>
int main(void) {
    int a[] = {10, 20, 30, 40, 50};
    int *p = a + 1;
    int *q = &a[4];
    p += 2;
    q -= p - a;
    printf("%d %d %td\n", *p, *q, p - q);
    return 0;
}
```

What does this program print?

- A. `40 30 1`
- B. `40 20 2`
- C. `30 20 1`
- D. `40 20 3`

**Answer:** B

**Behaviour class:** Well-defined. Every pointer stays inside the array.

**Why** (positions written as array indices):

| Line | Statement | p | q | Notes |
|---|---|---|---|---|
| 4 | `p = a + 1` | 1 | — | |
| 5 | `q = &a[4]` | 1 | 4 | |
| 6 | `p += 2` | **3** | 4 | |
| 7 | `p - a` | | | = 3 (elements) |
| 7 | `q -= 3` | 3 | **1** | |
| 8 | `*p` = a[3] = 40; `*q` = a[1] = 20; `p - q` = 3 − 1 = 2 | | | `40 20 2` |

**Why the other options fail:**
- A: Computes `p - a` as 2 (using p's starting position), giving q = 2 and `*q` = 30.
- C: Forgets `p += 2` happens before line 7's computation.
- D: Gets the difference as 3 by using q's old index.

**Concept recap:** Treating pointers as indices into an array (as in the table) is the fastest way to trace these questions. Convert back to values only at the end.

**Trap alert:** `p - q` is **negative** if `q` is after `p`. The result is a signed `ptrdiff_t`, so print it with `%td`, not `%u` or `%zu`.

**Interview angle:** *"Find the index of an element given a pointer to it."* `ptr - arr`. This works for any element type, because pointer subtraction divides by the element size automatically.

---

## Q93. [M6 | Which pointer operations are legal? | Hard | MSQ | [C] · C11 | GATE+OA] · Type T3 · P18

Given `int a[5]; int *p = a, *q = a + 3;`, select **all** expressions that **compile** and have **well-defined** behaviour.

- A. `q - p`
- B. `p + q`
- C. `p + 5`
- D. `q - 4`
- E. `p < q`

**Answer:** A, C, E

**Behaviour class per option:** A = well-defined (3) · B = **ill-formed** · C = well-defined · D = **UB** · E = well-defined (1)

**Why each option:**
- **A ✓:** The difference of two pointers into the same array = **3** elements (verified).
- **B ✗ ill-formed:** **Adding two pointers** is a constraint violation, a compile error ("invalid operands to binary +"). An address plus an address has no meaning.
- **C ✓:** `a + 5` is the one-past-the-end pointer. Legal, as long as it isn't dereferenced.
- **D ✗ UB:** `a + 3 - 4` = `a - 1`, before the array's start. Computing it is UB (§6.5.6p8).
- **E ✓:** A relational comparison between pointers into the **same** array is defined → **1**.

**Concept recap:** The allowed pointer operations:

| Operation | Result |
|---|---|
| pointer ± integer | pointer |
| pointer − pointer (same array) | integer (`ptrdiff_t`) |
| pointer `== != < > <= >=` pointer (same array for the relational ones) | `int` 0/1 |
| pointer + pointer, × or ÷ involving pointers | **not allowed** |

**Trap alert:** A binary-search midpoint written as `mid = (lo + hi) / 2` doesn't even compile for pointers. The pointer version `mid = lo + (hi - lo) / 2` also avoids integer overflow when used with indices.

**Interview angle:** *"Why is `lo + (hi - lo) / 2` preferred over `(lo + hi) / 2`?"* "With ints, `lo + hi` can overflow (that was the famous Java `binarySearch` bug). With pointers, adding two pointers isn't even legal. The difference-based form works for both."

---

## Q94. [M6 | Measuring a string with pointer difference | Easy | MCQ | [C] · C11 | GATE] · Type T6

```c
#include <stdio.h>
int my_strlen(const char *s) {
    const char *e = s;
    while (*e) e++;
    return e - s;
}
int main(void) {
    printf("%d %d\n", my_strlen("pointer"), my_strlen(""));
    return 0;
}
```

What does this program print?

- A. `7 0`
- B. `8 1`
- C. `7 1`
- D. `6 0`

**Answer:** A

**Behaviour class:** Well-defined. The two calls are independent, so their unspecified evaluation order doesn't matter here.

**Why:** `e` walks forward until it points at the `'\0'`. The distance `e - s` is the number of characters before it.

| Call | Characters stepped over | e stops at | `e - s` |
|---|---|---|---|
| `"pointer"` | p o i n t e r (7) | index 7 (`'\0'`) | **7** |
| `""` | none | index 0 (already `'\0'`) | **0** |

**Why the other options fail:**
- B, C: Count the terminator as part of the length.
- D: Off by one in the other direction.

**Concept recap:** `while (*e)` is short for `while (*e != '\0')`. `e - s` is a `ptrdiff_t`, converted here to the `int` return type (fine for short strings). The real `strlen` returns `size_t`.

**Trap alert:** Calling `my_strlen` on a char array with no `'\0'` (like `t` in Q80) runs off the end, which is UB. This is exactly why `strlen` on unterminated buffers is dangerous (P32, M7).

**Interview angle:** *"Implement `strlen` without indexing."* This is the answer. A follow-up may ask about the optimised version, which reads 8 bytes at a time and tests for a zero byte with bit tricks, as glibc does.

---

## Q95. [M6 | Comparing pointers to different objects | Hard | MCQ | [C] · C11 | All] · Type T3

```c
#include <stdio.h>
int main(void) {
    int x = 1, y = 2;
    int *p = &x, *q = &y;
    printf("%d\n", p == q);    /* line A */
    printf("%d\n", p < q);     /* line B */
    return 0;
}
```

How should lines A and B be classified in **C11**?

- A. Both are well-defined; line B prints 1, because `y` is declared after `x`.
- B. Line A is well-defined (prints `0`). Line B is **undefined behaviour** in C, because `<` is applied to pointers to different objects. (In C++ the result is unspecified, not UB.)
- C. Both are undefined behaviour.
- D. Line A is undefined behaviour; line B is well-defined.

**Answer:** B

**Behaviour class:** Line A is well-defined. Line B is **UB** in C (and unspecified in C++).

**Why:**
- **Equality** (`==`, `!=`) works on **any** two pointers of compatible type. Pointers to different objects compare unequal (N1570 §6.5.9p6). → `0`.
- **Relational** comparisons (`<`, `>`, `<=`, `>=`) are defined only when both pointers point into the **same** array (or the same struct, or one past its end). Otherwise the behaviour is undefined (§6.5.8p5).

> **⚠ UB box (line B)**
> - **Category:** relational comparison of pointers to unrelated objects.
> - **Rule:** N1570 §6.5.8p5. C++ ([expr.rel]) makes the result **unspecified** instead, and guarantees that `std::less<T*>` gives a consistent total order.
> - **On GCC x86-64,** this simply compares the two addresses. Whether it prints 0 or 1 depends on how locals are laid out on the stack. **Not guaranteed.**

**Why the other options fail:**
- A: Assumes the declaration order sets the memory order. It doesn't, and even if it did, the comparison would still be UB in C.
- C: Equality between any two valid pointers is always defined.
- D: Has the rules backwards.

**Concept recap:** Asking "are these two pointers the same?" is always OK. Asking "which comes first?" makes sense only within one array.

**Trap alert:** Code like `if (p >= buf && p < buf + n)` ("is p inside buf?") is technically UB when `p` points elsewhere. Portable code compares `(uintptr_t)` values, which is implementation-defined but not UB (Batch 6), or keeps track of an index instead.

**Interview angle:** *"How do you sort a container of raw pointers in C++?"* "Use `std::less<T*>` (the default for `std::set<T*>` and `std::map`). It is guaranteed to give a strict total order even where the built-in `<` is unspecified."

---

## Q96. [M6 | `*p++` vs `(*p)++` vs `++*p` vs `*++p` | Medium | MCQ | [C] · C11 | GATE+OA] · Type T5 · P29

```c
#include <stdio.h>
int main(void) {
    int a[] = {10, 20, 30};
    int *p = a;
    int x = *p++;
    int y = (*p)++;
    int z = ++*p;
    int w = *++p;
    printf("%d %d %d %d | %d %d %d\n", x, y, z, w, a[0], a[1], a[2]);
    return 0;
}
```

What does this program print?

- A. `10 20 21 30 | 10 21 30`
- B. `11 20 22 30 | 11 22 30`
- C. `10 21 22 20 | 10 22 30`
- D. `10 20 22 30 | 10 22 30`

**Answer:** D

**Behaviour class:** Well-defined. Each statement modifies only one object (either `p` or one element).

**Why:** Parse each expression first. Postfix `++` binds **tighter** than prefix `*`. Prefix operators apply right to left.

| Expression | Parsed as | Incremented | Value yielded |
|---|---|---|---|
| `*p++` | `*(p++)` | the **pointer** (after use) | the old `*p` |
| `(*p)++` | | the **element** (after use) | the old element |
| `++*p` | `++(*p)` | the **element** (before use) | the new element |
| `*++p` | `*(++p)` | the **pointer** (before use) | the element at the new position |

**Memory trace** (p shown as an index)

| Line | Statement | p | a[] | Result |
|---|---|---|---|---|
| 4 | `p = a` | 0 | {10, 20, 30} | |
| 5 | `x = *p++` | **1** | {10, 20, 30} | x = 10 |
| 6 | `y = (*p)++` | 1 | {10, **21**, 30} | y = 20 |
| 7 | `z = ++*p` | 1 | {10, **22**, 30} | z = 22 |
| 8 | `w = *++p` | **2** | {10, 22, 30} | w = 30 |

Prints `10 20 22 30 | 10 22 30`.

**Why the other options fail:**
- A: Forgets that `++*p` increments again, on top of the earlier `(*p)++`.
- B: Treats `*p++` as incrementing the element.
- C: Has `(*p)++` yield the new value, and `*++p` read the old position.

**Concept recap:** "Postfix binds tighter, prefix reads right to left." `*p++` is the standard idiom for "use the current element, then advance". It's the heart of every pointer loop.

**Trap alert:** `++(*p)` and `++*p` are identical. `*p++ = *q++` is fine (it changes two different objects), but `*p++ = *p` reads `p` and modifies it without sequencing, which is UB.

**Interview angle:** *"What does `*p++` do?"* "It returns the value `p` points to, and advances `p` to the next element. The `++` applies to the pointer, because postfix binds tighter than unary `*`." This exact question is asked constantly.

---

## Q97. [M6 | The `*d++ = *s++` copy loop | Hard | MCQ | [C] · C11 | All] · Type T5 · P29

```c
#include <stdio.h>
int main(void) {
    char src[] = "abc";
    char dst[8];
    char *s = src, *d = dst;
    while ((*d++ = *s++))
        ;
    printf("%s %td %td\n", dst, s - src, d - dst);
    return 0;
}
```

What does this program print?

- A. `abc 3 3`
- B. `abc 3 4`
- C. `abc 4 4`
- D. `ab 3 3`

**Answer:** C

**Behaviour class:** Well-defined. `dst` has room for 8 bytes, and only 4 are written.

**Why:** Each test of the loop condition copies one character, advances **both** pointers, and then tests the character that was **copied**. The loop stops right after copying `'\0'`, but that last iteration has still advanced both pointers.

| Iteration | Copies | s after | d after | Copied value | Continue? |
|---|---|---|---|---|---|
| 1 | `'a'` | src + 1 | dst + 1 | 'a' ≠ 0 | yes |
| 2 | `'b'` | src + 2 | dst + 2 | 'b' | yes |
| 3 | `'c'` | src + 3 | dst + 3 | 'c' | yes |
| 4 | `'\0'` | **src + 4** | **dst + 4** | 0 | **no** |

`dst` = "abc" (terminated). `s - src` = 4, and `d - dst` = 4. Prints `abc 4 4`.

**Why the other options fail:**
- A: Assumes the pointers stop *at* the `'\0'`. They advance past it on the final iteration.
- B: Advances only one of the two pointers in the final iteration.
- D: Assumes the loop exits before copying `'c'`.

**Concept recap:** This loop is the classic K&R `strcpy`. The extra parentheses around the assignment tell the compiler "yes, I mean `=`, not `==`" and silence `-Wparentheses`.

**Trap alert:** After the loop, `s` points **one past** the source's `'\0'`, and `d` one past the copied `'\0'`. If you need `d` at the terminator (to append more text, as in `strcat`), use `d - 1`.

**Interview angle:** *"What's wrong with `strcpy`?"* "It has no idea how big `dst` is. If `src` is longer, it overflows. Use `snprintf`, `strlcpy` (BSD/glibc 2.38+), or `memcpy` with a known length, and in C++ use `std::string`." (Buffer-overflow UB is covered in M7.)

---

## Q98. [M6 | Which expressions evaluate to 6? | Hard | MSQ | [C] · C11 | GATE+OA] · Type T1 · P29

Given `int a[] = {1, 5, 9}; int *p = a + 1;`, evaluate each expression **on a fresh copy** of `a` and `p` (the options don't affect each other). Select **all** expressions whose **value** is 6.

- A. `*p + 1`
- B. `*(p + 1)`
- C. `++*p`
- D. `*p++ + 1`
- E. `(*p)++`

**Answer:** A, C, D

**Why each option** (each verified separately):

| Option | Parsed as | Side effect | Value |
|---|---|---|---|
| A | `(*p) + 1` = 5 + 1 | none | **6** ✓ |
| B | `*(p + 1)` = a[2] | none | 9 |
| C | `++(*p)` | a[1] becomes 6 | **6** ✓ |
| D | `(*(p++)) + 1` = 5 + 1 | p moves to a[2] | **6** ✓ |
| E | `(*p)++` | a[1] becomes 6 | **5** (the old value) |

**Concept recap:** Ask two separate questions about each expression: (1) what **value** does it yield, and (2) what does it **change**? Option E changes `a[1]` to 6 but yields 5, which is exactly the trap.

**Trap alert:** In D, `++` applies to `p`, not to `*p`, and the `+ 1` is applied to the value read *before* `p` moved.

**Interview angle:** Interviewers often ask for "the value of the expression **and** the state afterwards". Always report both.

---

## Q99. [M6 | const pointer vs pointer to const | Medium | MCQ | [C] · C11 | All] · Type T2 · P17

```c
int main(void) {
    int x = 1, y = 2;
    const int *p1 = &x;          /* pointer to const int     */
    int *const p2 = &x;          /* const pointer to int     */
    const int *const p3 = &x;    /* const pointer to const   */
    *p1 = 5;   /* L1 */
    p1 = &y;   /* L2 */
    *p2 = 5;   /* L3 */
    p2 = &y;   /* L4 */
    p3 = &y;   /* L5 */
    return 0;
}
```

Which lines **fail to compile**?

- A. L1, L4, L5
- B. L2, L3
- C. L1, L5
- D. L3, L4, L5

**Answer:** A

**Behaviour class:** Ill-formed (L1, L4 and L5 are constraint violations: assigning to something that isn't a modifiable lvalue).

**Why:** Read the declaration **right to left**, starting at the variable name.

| Declaration | Read as | Can change `*p`? | Can change `p`? |
|---|---|---|---|
| `const int *p1` | p1 is a pointer to an int that is const | **no** | yes |
| `int *const p2` | p2 is a **const** pointer to int | yes | **no** |
| `const int *const p3` | p3 is a const pointer to a const int | no | no |

| Line | Statement | Result |
|---|---|---|
| L1 | `*p1 = 5` | ✗ "assignment of read-only location '*p1'" |
| L2 | `p1 = &y` | ✓ |
| L3 | `*p2 = 5` | ✓ (x becomes 5) |
| L4 | `p2 = &y` | ✗ "assignment of read-only variable 'p2'" |
| L5 | `p3 = &y` | ✗ "assignment of read-only variable 'p3'" |

(GCC 13 error messages verified.)

**Why the other options fail:**
- B: Lists exactly the lines that **do** compile.
- C: Misses L4 (a const pointer can't be re-pointed).
- D: Treats `int *const` as making the *pointee* const.

**Concept recap:** `const` applies to whatever is on its **left**. If nothing is on its left, it applies to what's on its right. So `const int *` and `int const *` mean the same thing, and the position of `const` relative to `*` is what matters.

**Trap alert:** `const int *p` means "you can't modify x **through p**". It doesn't make `x` itself constant: `x = 5;` directly is still fine.

**Interview angle:** *"Explain `const char *`, `char *const` and `const char *const`."* Use the right-to-left reading and the table. A common follow-up: *"Why is `char **` → `const char **` not an implicit conversion?"* Because it would let you store a `const char *` into a slot that's later used as a `char *`, which would allow modifying a const object without any cast.

---

## Q100. [M6 | Casting away const | Hard | MCQ | [C] · C11 | All] · Type T3 · P17

```c
#include <stdio.h>
int main(void) {
    int a = 10;
    const int b = 20;
    const int *pa = &a;
    int *q1 = (int *)pa;
    *q1 = 11;
    int *q2 = (int *)&b;
    *q2 = 21;
    printf("%d %d %d\n", a, b, *q2);
    return 0;
}
```

Which statement is correct?

- A. Both writes are well-defined, and it prints `11 21 21`.
- B. Both writes are undefined behaviour, because both cast away `const`.
- C. Compile error at the casts: `const` can't be cast away in C.
- D. The write through `q1` is well-defined, because `a` itself is not const. The write through `q2` is **undefined behaviour**, because `b` was **defined** `const`.

**Answer:** D

**Behaviour class:** The write through `q1` is well-defined. The write through `q2` is **UB**.

**Why:** What matters is whether the **object itself** was defined `const`, not the type of the pointer used to reach it.

| Write | Object | Was the object defined const? | Result |
|---|---|---|---|
| `*q1 = 11` | `a` | no; it was only *viewed* through a `const int *` | ✓ a = 11 |
| `*q2 = 21` | `b` | **yes** | **UB** (N1570 §6.7.3p6) |

> **⚠ UB box (the write through `q2`)**
> - **Category:** modifying an object defined with a `const`-qualified type.
> - **Rule:** N1570 §6.7.3p6: an attempt to modify an object defined with a const-qualified type through a non-const lvalue is UB. (C++: [dcl.type.cv]p4.)
> - **Observed on GCC 13 x86-64:** at **`-O0`** it prints `11 21 21`; at **`-O2`** it prints **`11 20 21`** (verified). The optimiser treated `b` as the constant 20 and never reloaded it from memory, so `b` and `*q2` disagree even though they're the same object. A `const` global might even live in read-only memory, and then the write would crash. **No result is guaranteed.**

**Why the other options fail:**
- A: Correct for `a`, wrong for `b`. The `-O2` result shows the problem directly.
- B: Casting away const is fine when the underlying object isn't const. This is how a function that takes `const T *` can legally return a non-const `T *` (as `strchr` does).
- C: C allows the explicit cast. (C++ requires `const_cast` for it, which is pattern P54.)

**Concept recap:** `const` on a **pointer** limits what *you* can do through that pointer. `const` on an **object definition** promises the object never changes, and the compiler is allowed to rely on that promise.

**Trap alert:** The cast compiles silently, with no warning by default. The bug shows up only with optimisation or when the object lives in read-only memory, which makes it hard to diagnose.

**Interview angle:** *"When is `const_cast` safe in C++?"* "Only when the object was not originally declared const. For example, calling a legacy API that takes `char *` but doesn't modify it, or implementing a non-const getter in terms of the const one. Writing to an object that was defined const is UB." (M12.)

---

## Batch 5 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 81 | M5 | E | MCQ | C | **UB** |
| 82 | M5 | H | MSQ | B, D, E | UB / well-defined mix |
| 83 | M5 | M | MCQ | A | Well-defined |
| 84 | M5 | M | MCQ | C | Ill-formed / well-defined |
| 85 | M6 | M | MSQ | A, B, C | — (concepts) |
| 86 | M6 | E | MCQ | D | Well-defined |
| 87 | M6 | H | MCQ | D | **UB** at L7 |
| 88 | M6 | E | MCQ | B | Well-defined |
| 89 | M6 | M | MCQ | A | Well-defined |
| 90 | M6 | M | MCQ | D | Well-defined |
| 91 | M6 | H | MCQ | C | Well-defined (byte order implementation-defined) |
| 92 | M6 | M | MCQ | B | Well-defined |
| 93 | M6 | H | MSQ | A, C, E | Ill-formed / UB / OK mix |
| 94 | M6 | E | MCQ | A | Well-defined |
| 95 | M6 | H | MCQ | B | **UB** (C) / unspecified (C++) |
| 96 | M6 | M | MCQ | D | Well-defined |
| 97 | M6 | H | MCQ | C | Well-defined |
| 98 | M6 | H | MSQ | A, C, D | Well-defined |
| 99 | M6 | M | MCQ | A | Ill-formed |
| 100 | M6 | H | MCQ | D | **UB** (`q2` write) |

- **Difficulty:** 4 E / 8 M / 8 H.
  - M5 total (Q69–Q84): 4/7/5, as planned.
  - M6 so far: 3/6/7.
- **Format:** 16 MCQ / 4 MSQ.
  - M5 total: 13/3, as planned.
- **MCQ answer keys:** A×4, B×3, C×4, D×5. Running total Q1–Q100: A 20, B 20, C 20, D 21.

---

**Completed: 100 / 400 Remaining: 300**

**Ledger, newly ✓ this batch:**
- P17 ×2 (Q99, Q100) ★ complete
- P18 ×4 (Q90, Q91, Q92, Q93) ★ complete
- P28 ×2 (Q81, Q82) ★ complete
- P29 ×3 (Q96, Q97, Q98) ★ complete
- P05 → 4 (Q83)

**M5 complete** (Q69–Q84, 16 questions). **M6 in progress** (16 of 35 done).

---
---

# SELF-AUDIT #1 (after Batch 5, Q1–Q100)

## A. Actual vs planned splits

**Difficulty**

| | Easy | Medium | Hard |
|---|---|---|---|
| Actual, Q1–Q100 | 27 | 41 | 32 |
| Module plan for exactly the material covered (M1–M5 + 16 of M6's 35) | 27 | 41 | 32 |
| Bank-wide target (25/40/35 over 400) | 25% | 40% | 35% |

✅ An exact match with the per-module plan. The bank-wide Hard share (35%) will be reached in the C++ modules M14, M16, M18 and M25, which are planned Hard-heavy.

**Format:** MCQ 81 / MSQ 19 (target 80/20). ✅

**Language:** [C] 100 / [C++] 0. ✅ As designed: all C modules come first, and the bank-wide target of 195 C / 205 C++ is unchanged.

**MCQ answer-key balance:** A 20 · B 20 · C 20 · D 21. ✅ "All/none of the above" used: **0** (limit 5).

## B. Relevance-tag split

| Tag | Actual Q1–Q100 | Expected pace for 100 C questions (out of the 195 C planned) | Status |
|---|---|---|---|
| [GATE] | 20 | ≈ 19.5 | ✅ |
| [GATE+OA] | 48 | ≈ 57 | ⚠ under |
| [All] | 29 | ≈ 18 | ⚠ over |
| [OA+Interview] | 3 | 5 for all C questions | ✅ (2 left) |
| [OA] / [Interview] | 0 | 2 / 3 for all C questions | ✅ |

**Correction:** The UB-heavy pointer and expression questions naturally earn [All], because UB awareness is tested in GATE, OAs *and* interviews. The remaining C modules (M6 part 2, M7, M9) are also UB-heavy. So I'm **re-targeting the C tags to [All] ≈ 45 and [GATE+OA] ≈ 102** (the C total stays 195). For the rest of the C modules, a plain output-tracing question defaults to [GATE+OA], and [All] is kept for questions whose main lesson is UB or memory safety.

## C. Pointer + OOP share

- M6 so far: 16 / 100 = 16%.
- Planned final: M6 (35) + M13–M18 (110) = 145 / 400 = 36.25%.
- ✅ On schedule. The OOP family starts at Q201.

## D. Full Pattern Ledger (✓ = covered, with the count; ☐ = not yet)

| ID | Pattern | ★ | Count | Questions | Status |
|---|---|---|---|---|---|
| P01 | Data-type sizes & signedness | ★ | 2 | Q1, Q2 | ✓ |
| P02 | Integer promotion & usual arithmetic conversions | ★ | 5 | Q2, Q4, Q21, Q22, Q23 | ✓ |
| P03 | Precedence & associativity | ★ | 4 | Q19, Q20, Q24, Q25 | ✓ |
| P04 | Short-circuit evaluation | | 2 | Q25, Q49 | ✓ |
| P05 | sizeof expression vs type | ★ | 4 | Q3, Q4, Q17, Q83 | ✓ |
| P06 | Pre/post increment | ★ | 2 | Q26, Q27 | ✓ |
| P07 | Comma / ternary | | 2 | Q28, Q29 | ✓ |
| P08 | Bitwise on signed integers | | 1 | Q30 | ✓ (1 more planned in M25) |
| P09 | Signed overflow — UB | ★ | 2 | Q31, Q32 | ✓ |
| P10 | Literal forms | | 2 | Q5, Q6 | ✓ |
| P11 | Signed vs unsigned comparison | ★ | 4 | Q22, Q33, Q34, Q50 | ✓ |
| P12 | Recursion trace | ★ | 6 | Q54, Q55, Q56, Q58, Q61, Q62 | ✓ |
| P13 | Static local persistence | ★ | 2 | Q57, Q58 | ✓ |
| P14 | extern vs static linkage | | 3 | Q7, Q8, Q16 | ✓ |
| P15 | Stack frame / parameter passing | | 3 | Q51, Q52, Q53 | ✓ |
| P16 | Function overloading resolution (C++) | | 0 | — | ☐ M17 |
| P17 | const ptr vs ptr to const | ★ | 2 | Q99, Q100 | ✓ (+1 planned in B6) |
| P18 | Pointer arithmetic | ★ | 4 | Q90, Q91, Q92, Q93 | ✓ |
| P19 | Array decay | ★ | 5 | Q71, Q73, Q74, Q75, Q78 | ✓ |
| P20 | sizeof array vs pointer | ★ | 3 | Q71, Q72, Q75 | ✓ |
| P21 | Multi-dim address calculation | | 3 | Q76, Q77, Q78 | ✓ |
| P22 | Pointer to pointer | | 0 | — | ☐ B6 |
| P23 | Function pointer | | 1 | Q63 | ✓ (more in B6) |
| P24 | void* cast semantics | | 0 | — | ☐ B6 |
| P25 | String literal vs char array | ★ | 0 | — | ☐ B6/B7 |
| P26 | NULL vs nullptr vs 0 | | 0 | — | ☐ M12 |
| P27 | Dangling pointer | | 1 | Q64 | ✓ |
| P28 | Out-of-bounds — UB | ★ | 2 | Q81, Q82 | ✓ |
| P29 | `*p++` family | ★ | 3 | Q96, Q97, Q98 | ✓ |
| P30 | Aliasing / strict aliasing | | 0 | — | ☐ B6 |
| P31 | Array init & partial init | | 2 | Q69, Q70 | ✓ |
| P32 | String functions | ★ | 0 | — | ☐ M7 |
| P33 | Struct padding | ★ | 0 | — | ☐ M8 |
| P34 | Union type-punning | ★ | 0 | — | ☐ M8 |
| P35 | Bit-fields | | 0 | — | ☐ M8 |
| P36 | Self-referential struct | | 0 | — | ☐ M8/M24 |
| P37 | Macro expansion trace | ★ | 0 | — | ☐ M10 |
| P38 | Macro side effects | ★ | 0 | — | ☐ M10 |
| P39 | Stringify / token paste | | 0 | — | ☐ M10 |
| P40 | Conditional compilation | | 0 | — | ☐ M10 |
| P41 | Header guards vs `#pragma once` | | 0 | — | ☐ M10 |
| P42 | Compile vs link error | ★ | 3 | Q7, Q9, Q18 | ✓ |
| P43 | malloc/free vs new/delete | ★ | 0 | — | ☐ M12 |
| P44 | new[]/delete[] mismatch | | 0 | — | ☐ M12 |
| P45 | Double free / use-after-free | ★ | 0 | — | ☐ M9 |
| P46 | Memory leaks | | 0 | — | ☐ M9 |
| P47 | realloc NULL / size 0 | | 0 | — | ☐ M9 |
| P48 | Uninitialised read — UB | ★ | 2 | Q10, Q65 | ✓ |
| P49 | Sequence-point violation — UB | ★ | 2 | Q35, Q36 | ✓ |
| P50 | Integer truncation | | 2 | Q11, Q37 | ✓ |
| P51–P96 | C++ fundamentals, OOP, modern C++, templates, STL, concurrency | | 0 | — | ☐ M12–M23 (as planned) |
| P97 | char signedness | | 1 | Q12 | ✓ |
| P98 | sizeof implementation-defined | | 1 | Q13 | ✓ |
| P99–P105 | Static init fiasco, dangling reference, mixing malloc/delete, integrated patterns | | 0 | — | ☐ M12/M14/M25 |

**Totals:** 31 / 105 patterns ✓. ★ patterns complete: **18 / 51** (P01, P02, P03, P05, P06, P09, P11, P12, P13, P17, P18, P19, P20, P28, P29, P42, P48, P49). Every pattern that was scheduled for Batches 1–5 is ✓. There are no stragglers.

## E. Re-trace of every output; corrections

Every well-defined output question in Q1–Q100 was compiled with GCC 13 (and Clang 18 where evaluation order mattered) and run. **All 100 answer keys match.** The following corrections were made — two of them in this audit, the rest before the batch shipped:

| Where | Issue | Fix |
|---|---|---|
| Batch 3, Q55, option D | Its explanation didn't add up to the printed numbers | Option D changed to `3 8 10 \| 10`, with a consistent derivation |
| Batch 3, Q58, option C | Its derivation gave 15, not the 17 printed | Option C changed to `15 35`, with a consistent derivation |
| Batch 4, Q77, options A–C | Two distractors had no consistent derivation | Replaced with 1068 / 1316 / 1082, each from a named mistake |
| **Batch 4, Q67 (found in this audit)** | Calling an undeclared function was described as a "constraint violation" | **More precisely, it's a *syntax* violation:** N1570 §6.5.1p2, footnote 91: "an undeclared identifier is a violation of the syntax". The consequence (a diagnostic is required, i.e. a compile error) and the answer are unchanged. **Batch 4 has been updated.** (Q14's implicit-`int` case really is a constraint violation, §6.7.2p2, so it was already correct.) |
| **Batch 2, Q39, option B (found in this audit)** | The note "C++14–17 made `1 << 31` defined" was too loose | **More precisely:** in C++14–17, `1 << 31` is not UB. The result is the unsigned value 2³¹ converted to `int`, which is implementation-defined (−2³¹ on GCC). C++20 makes all left shifts well-defined modulo 2ⁿ. The C11 answer (UB) is unchanged. **Batch 2 has been updated.** |

## F. ⚠ Highest-priority check: does any explanation claim a definite output for code that has UB?

I went through every question whose code contains UB or something close to it:
- **Labelled UB, with the correct option saying "undefined" (never a value):** Q10, Q31, Q35, Q37, Q64, Q65, Q81, Q87 (L7), Q95 (line B), Q100 (q2), plus the UB options in the multiple-answer questions Q32, Q36, Q39, Q50 (D), Q66 (C, D), Q82, Q93 (D). Each one has the UB box: category, rule, and "typically prints X, not guaranteed".
- **Code that looks like UB but is actually well-defined, re-checked:**
  - Q17, Q71, Q75, Q83: uninitialised objects that appear only in `sizeof` or address computations.
  - Q74: `&a + 1`, the one-past-the-end pointer.
  - Q26: `y = ++x + y`, where `y` is read only to compute the value stored.
  - Q91: reading bytes through `unsigned char *`.
  - Q56 and Q58: evaluation order is unspecified, but the code was built so the result doesn't depend on it.
  - Q98: each option is evaluated on a fresh copy.

**Result: 0 questions claim a deterministic output for UB code. Nothing to repair.**

---

# PROGRESS REPORT @ 100

**Modules complete:** M1 (types, storage, linkage), M2 (operators, expressions, UB in expressions), M3 (control flow), M4 (functions, recursion), M5 (arrays). **In progress:** M6 (pointers, 16 of 35).

**Key traps covered:**
1. `sizeof` does **not** evaluate its operand, except for VLAs (Q3, Q83).
2. `sizeof('a')` is 4 in C and 1 in C++ (Q4).
3. Every array parameter is really a pointer: `sizeof` inside the function gives 8 (Q71).
4. In a mixed `int`/`unsigned` comparison, −1 becomes huge (Q22, Q33, Q34). On Windows (LLP64), `long` vs `unsigned` behaves differently (Q22).
5. `&&` binds tighter than `||`, and short-circuiting skips the *entire* right operand (Q25).
6. A `default:` placed mid-switch still falls through into the cases below it (Q42). A stray `;` becomes a loop body (Q46). The dangling `else` binds to the nearest `if` (Q47).
7. A pointer parameter is a copy: re-pointing it inside a function doesn't affect the caller (Q52).
8. A `static` local persists across calls *and* across separate recursive runs (Q57, Q58).
9. `a` vs `&a`: same address, different step sizes (Q74). A 2D array decays to `int (*)[C]`, never to `int **` (Q75, Q79).
10. `*p++` advances the pointer; `(*p)++` increments the element (Q96–Q98).

**Misconceptions tested:** "signed integers wrap around" · "garbage value, but well-defined" · "conversion overflow = UB" (it's implementation-defined for integer → integer, but UB for float → int) · "right shift of a negative number is UB" (it's implementation-defined) · "arguments are evaluated left to right" (GCC and Clang disagree, Q68) · "`static` fixes uninitialised variables" (it changes the semantics, Q65) · "a `const` cast is always safe" (Q100).

**UB categories covered so far (12):**
1. signed overflow
2. `INT_MIN / -1`
3. unsequenced modifications
4. uninitialised reads
5. float → int out of range
6. shift-count and left-shift UB
7. `va_arg` type mismatch / reading too many arguments
8. returning a local's address (dangling pointer)
9. out-of-bounds access and pointer arithmetic beyond one-past-the-end
10. relational comparison of unrelated pointers
11. modifying a const-defined object
12. (in Q50) a loop that can only end through signed overflow

**Still to come:**
- C side: modifying string literals, `strlen` on unterminated buffers, use-after-free and double free, strict aliasing, `printf` format mismatches (mentioned so far, not yet tested)
- C++ side: `delete` through a base pointer with a non-virtual destructor, `new[]`/`delete` mismatch, pure-virtual calls, data races, union punning in C++

**Areas to revise (the questions most often missed on first try in this material):** Q22 (rank/representability rule), Q36 (C vs C++17 sequencing), Q45 (non-decreasing triples), Q58 (static + recursion), Q68 (unspecified vs UB), Q74/Q75 (pointer-to-array types), Q100 (const object vs const view).

*Say **"Continue."** for Batch 6 (Q101–Q120: M6 part 2 — the last const-pointer question, pointer-to-pointer, the array/pointer relationship, arrays of pointers vs pointers to arrays, string pointers, function pointers and declarators, `void*` and casts, `uintptr_t` and alignment, strict aliasing and `restrict`, dangling pointers — then M7 strings begins).*
