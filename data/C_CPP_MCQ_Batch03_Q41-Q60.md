# C/C++ MCQ/MSQ Bank — Batch 3 (Q41–Q60)

**What this batch covers:**
- **M3 Control flow (Q41–Q50):** `switch` fall-through and `break` inside a loop; `default` placed in the middle of a `switch`; `continue` and `goto` in nested loops; counting iterations of nested loops (doubling bounds, triangular triple loops); the stray `;` after `for`; the dangling `else`; `do-while` with a post-decrement; short-circuit evaluation in a loop condition; unsigned loop counters that never terminate.
- **M4 Functions & recursion (Q51–Q60):** swapping by value vs through pointers; a pointer is itself passed by value; activation records and stack depth; tracing recursive output when the print sits between two recursive calls; post-order printing while a value is built up; counting the calls made by naive Fibonacci; `static` locals that persist across calls and across separate recursive runs; indirect (mutual) recursion; tail recursion.

**Patterns:** P04, P11, P12, P13, P15

> **Conventions for every question in this batch (unless the question says otherwise):**
> - C11, `gcc -std=c11`, 64-bit Linux / GCC / x86-64
> - `sizeof(int)=4`, plain `char` signed, two's complement
> - `stdio.h` is included
> - The output does not depend on optimisation level.
> - Every program here was compiled and run, and its output matches the answer key.

---

## Q41. [M3 | switch fall-through + break inside a loop | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    for (int i = 0; i < 4; i++) {
        switch (i) {
            case 0: printf("A");
            case 1: printf("B"); break;
            case 2: printf("C");
            default: printf("D");
        }
    }
    printf("\n");
    return 0;
}
```

What does this program print?

- A. `ABCD`
- B. `AB`
- C. `ABBCDD`
- D. `ABBCD`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** A `case` label is only a jump target. Once execution enters the `switch`, it runs **straight down** through the following labels until it reaches a `break` or the closing `}`. A `break` inside a `switch` exits **only the switch**, not the enclosing loop.

| i | Jumps to | Statements executed | Output so far |
|---|---|---|---|
| 0 | `case 0` | `A`, falls into `case 1`: `B`, `break` | `AB` |
| 1 | `case 1` | `B`, `break` | `ABB` |
| 2 | `case 2` | `C`, falls into `default`: `D`, reaches `}` | `ABBCD` |
| 3 | no match → `default` | `D` | `ABBCDD` |
| 4 | loop ends | `\n` | `ABBCDD` |

**Why the other options fail:**
- A: Assumes every case stops automatically, so each `i` prints one letter. C has no automatic break.
- B: Assumes the `break` ends the `for` loop. It only leaves the `switch`.
- D: Forgets that i = 3 matches no case and so jumps to `default`.

**Concept recap:**
- A `switch` jumps to the matching `case`, or to `default` if no case matches. If there is no `default` either, it skips the whole body.
- From the label it jumps to, execution continues in source order.
- `break` exits the innermost `switch` or loop. `continue` affects only loops, never a `switch`.

**Trap alert:** The missing-`break` bug is so common that GCC's `-Wimplicit-fallthrough` warns about it. When fall-through is intentional, mark it with a `/* fall through */` comment or C23 `[[fallthrough]];`.

**Interview angle:** *"When is fall-through useful?"* "To group cases that share code, for example `case 'a': case 'e': case 'i': … vowel++; break;`, or for staged initialisation such as Duff's device." Follow-up: *"How do you exit a loop from inside a switch?"* Use a flag, `goto`, or move the loop into a function and `return`. A `break` won't reach the loop.

---

## Q42. [M3 | default in the middle of a switch | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int x = 5, sum = 0;
    switch (x) {
        case 1: sum += 1;
        case 2: sum += 2;
        case 3: sum += 3;
        default: sum += 10;
        case 4: sum += 4;
    }
    printf("%d\n", sum);
    return 0;
}
```

What does this program print?

- A. `10`
- B. `14`
- C. `20`
- D. `0`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:** The position of `default` doesn't affect *when* it is chosen: it is used only when no `case` matches, wherever it appears. But once control jumps there, execution falls through to whatever labels **follow it in the source**.

| Step | Action | sum |
|---|---|---|
| 1 | x = 5 matches none of 1, 2, 3, 4 → jump to `default` | 0 |
| 2 | `sum += 10` | 10 |
| 3 | no `break` → fall through into `case 4:` → `sum += 4` | 14 |
| 4 | reach `}` | 14 |

**Why the other options fail:**
- A: Assumes execution stops after `default`, or that `default` "must be last". Neither is true.
- C: Runs every statement from `case 1` (1 + 2 + 3 + 10 + 4 = 20), as if the switch started at the top.
- D: Assumes an unmatched value skips the body. That only happens when there is **no** `default`.

**Concept recap:** `default` can appear anywhere, at most once per `switch`. Case labels must be distinct integer constant expressions.

**Trap alert:** If x were 4, the switch would jump straight to `case 4` → sum = 4. The `default` block that appears earlier in the source would **not** run.

**Interview angle:** *"Can a switch use strings?"* "In C and C++, no. Case labels must be integral constants (enums and characters are fine). Use `if`/`strcmp`, a hash, or a lookup table." In C++ you can switch on a `constexpr` hash of the string.

---

## Q43. [M3 | continue + goto in nested loops | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int i, j, count = 0;
    for (i = 0; i < 3; i++) {
        for (j = 0; j < 3; j++) {
            if (j == 1) continue;
            if (i == 2) goto done;
            count++;
        }
    }
done:
    printf("%d %d %d\n", i, j, count);
    return 0;
}
```

What does this program print?

- A. `3 3 6`
- B. `2 1 4`
- C. `3 3 4`
- D. `2 0 4`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:** `continue` jumps to the inner loop's increment step (`j++`). `goto done` leaves **both** loops immediately, without running any increment.

| i | j | `j == 1`? | `i == 2`? | Action | count |
|---|---|---|---|---|---|
| 0 | 0 | no | no | count++ | 1 |
| 0 | 1 | **yes** | — | continue | 1 |
| 0 | 2 | no | no | count++ | 2 |
| 1 | 0 | no | no | count++ | 3 |
| 1 | 1 | yes | — | continue | 3 |
| 1 | 2 | no | no | count++ | 4 |
| 2 | 0 | no | **yes** | `goto done` (i, j frozen at 2, 0) | 4 |

Prints `2 0 4`.

**Why the other options fail:**
- A: Ignores both `continue` and `goto`.
- B: Assumes `j` had already moved past 0. The `goto` fires on the very first inner iteration of i = 2.
- C: Counts correctly, but assumes the loops ran to completion (i = 3, j = 3).

**Concept recap:**
- `continue` in a `for` loop → jump to the increment expression, then re-test the condition.
- `continue` in a `while`/`do-while` → jump straight to the condition.
- `goto` may jump to any label in the same function. It is mainly used to break out of deeply nested loops and for cleanup paths.

**Trap alert:** `continue` inside a `while` loop can skip the update at the bottom of the body, creating an infinite loop: `while (i < n) { if (bad) continue; i++; }`.

**Interview angle:** *"Is `goto` ever acceptable?"* "Yes, for breaking out of nested loops and for the single-exit cleanup pattern in C (`goto cleanup;` that frees resources in reverse order). The Linux kernel uses this heavily. In C++, RAII replaces the cleanup use."

---

## Q44. [M3 | Iteration count with a doubling bound | Easy | MCQ | [C] · C11 | GATE] · Type T1

```c
#include <stdio.h>
int main(void) {
    int count = 0;
    for (int i = 1; i <= 10; i *= 2)
        for (int j = 0; j < i; j++)
            count++;
    printf("%d\n", count);
    return 0;
}
```

What does this program print?

- A. `15`
- B. `31`
- C. `10`
- D. `40`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:**

| Outer `i` | `i <= 10`? | Inner iterations (`j = 0 … i−1`) | count after |
|---|---|---|---|
| 1 | yes | 1 | 1 |
| 2 | yes | 2 | 3 |
| 4 | yes | 4 | 7 |
| 8 | yes | 8 | 15 |
| 16 | **no** | — | 15 |

Total = 1 + 2 + 4 + 8 = 2⁴ − 1 = **15**.

**Why the other options fail:**
- B: Includes i = 16 (1 + 2 + 4 + 8 + 16 = 31). 16 > 10, so the outer loop has already stopped.
- C: Counts each outer value as 10/i, or confuses the bound with the count.
- D: Multiplies the outer count (4) by the bound (10), as if the inner loop always ran 10 times.

**Concept recap:** For a doubling outer loop up to n, the total work is 1 + 2 + 4 + … + 2^⌊log₂n⌋ < 2n, which is **O(n)**, not O(n log n). This is the same argument that gives amortised O(1) for `vector::push_back`.

**Trap alert:** "The outer loop runs log n times and the inner one up to n times, so it's O(n log n)." That's a valid upper bound but not tight. Sum the actual inner counts.

**Interview angle:** *"What is the time complexity?"* "O(n). The geometric series sums to less than 2n." This is a favourite way to test whether a candidate multiplies loop bounds blindly.

---

## Q45. [M3 | Triangular triple loop | Hard | MCQ | [C] · C11 | GATE] · Type T1

```c
#include <stdio.h>
int main(void) {
    int count = 0;
    for (int i = 0; i < 5; i++)
        for (int j = i; j < 5; j++)
            for (int k = j; k < 5; k++)
                count++;
    printf("%d\n", count);
    return 0;
}
```

What does this program print?

- A. `125`
- B. `10`
- C. `35`
- D. `15`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** The loops visit every triple (i, j, k) with 0 ≤ i ≤ j ≤ k ≤ 4 exactly once.

Counting directly, for each i: the number of pairs j ≤ k in [i, 4] is T(5 − i), where T(m) = m(m+1)/2.

| i | m = 5 − i | pairs (j, k) = T(m) | running count |
|---|---|---|---|
| 0 | 5 | 15 | 15 |
| 1 | 4 | 10 | 25 |
| 2 | 3 | 6 | 31 |
| 3 | 2 | 3 | 34 |
| 4 | 1 | 1 | 35 |

Combinatorial check: choosing 3 values from 5 **with repetition, order ignored**, is C(5 + 3 − 1, 3) = C(7, 3) = **35**.

**Why the other options fail:**
- A: 5³. That would be the count if every loop started at 0.
- B: C(5, 3), strictly increasing i < j < k. The loops start at `j = i` and `k = j`, so equal values are included.
- D: C(6, 2), counting only pairs (one level of nesting).

**Concept recap:**
- Loops starting at the previous index (`j = i`, `k = j`) count **non-decreasing** tuples: C(n + r − 1, r).
- Loops starting one after (`j = i + 1`, `k = j + 1`) count **strictly increasing** tuples: C(n, r).

**Trap alert:** Check each loop's start (`= i` vs `= i + 1`) and its condition (`<` vs `<=`). Those two details decide between 10, 35 and 125.

**Interview angle:** *"Complexity of this triple loop for general n?"* "About n³/6, so it's still Θ(n³). The constant is smaller, but the growth rate is the same."

---

## Q46. [M3 | Stray semicolon after for | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int i, sum = 0;
    for (i = 0; i < 5; i++);
        sum += i;
    printf("%d %d\n", i, sum);
    return 0;
}
```

What does this program print?

- A. `5 10`
- B. `4 10`
- C. `5 15`
- D. `5 5`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:** The `;` immediately after `for (...)` is an **empty statement**, and that empty statement *is* the loop body. The indented `sum += i;` is not part of the loop; it runs **once**, after the loop ends.

| Line | Statement executed | Variables after execution |
|---|---|---|
| 3 | `int i, sum = 0;` | sum = 0 |
| 4 | `for (...) ;` runs 5 times doing nothing | i = 5 (the value that made `i < 5` false) |
| 5 | `sum += i;` (once) | sum = 5 |
| 6 | `printf` | `5 5` |

**Why the other options fail:**
- A: Treats `sum += i` as the loop body (0 + 1 + 2 + 3 + 4 = 10).
- B: Assumes i ends at 4. A `for` loop exits only when its condition fails, which happens at i = 5.
- C: Sums 1 to 5.

**Concept recap:** Without braces, a loop or `if` body is exactly **one statement**, and `;` alone counts as a statement. Indentation means nothing to the compiler. GCC's `-Wmisleading-indentation` (in `-Wall`) catches many of these cases.

**Trap alert:** The same bug appears as `while (x > 0);` (an infinite loop if `x` is positive) and `if (cond);` (the following block always runs).

**Interview angle:** *"Why do style guides require braces even for one-line bodies?"* "To prevent stray-semicolon and dangling-else bugs, and to stop a second statement added later from silently landing outside the body. Apple's 2014 `goto fail;` TLS bug is the famous real-world case."

---

## Q47. [M3 | Dangling else | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int a = 0, b = 1;
    if (a)
        if (b)
            printf("X");
    else
        printf("Y");
    printf("Z\n");
    return 0;
}
```

What does this program print?

- A. `YZ`
- B. `Z`
- C. `XZ`
- D. `XYZ`

**Answer:** B

**Behaviour class:** Well-defined. (GCC warns: "suggest explicit braces to avoid ambiguous 'else'".)

**Why:** An `else` attaches to the **nearest preceding `if` that doesn't already have an `else`**. Here that's `if (b)`, regardless of the indentation. The code actually means:

```c
if (a) {
    if (b) printf("X");
    else   printf("Y");
}
printf("Z\n");
```

| Step | Test | Result |
|---|---|---|
| 1 | `if (a)` → a = 0 → false | the whole inner `if/else` is skipped |
| 2 | `printf("Z\n")` | prints `Z` |

**Why the other options fail:**
- A: Reads the code by its indentation, attaching the `else` to `if (a)`.
- C, D: Assume the inner `if` runs even though `a` is 0.

**Concept recap:** The rule follows from the C grammar (the "dangling else" ambiguity is resolved in favour of the nearest `if`). To pair the `else` with the outer `if`, put braces around the inner one: `if (a) { if (b) … } else …`.

**Trap alert:** GATE and OA papers deliberately indent dangling-else code misleadingly. Ignore the whitespace and parse the code.

**Interview angle:** *"How do languages avoid the dangling-else problem?"* "Some require braces (Go, Rust, Swift). Python uses indentation as syntax. C and C++ resolve it with the nearest-if rule."

---

## Q48. [M3 | do-while with post-decrement | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int n = 0, count = 0;
    do {
        count++;
    } while (n-- > 0);
    printf("%d %d\n", n, count);
    return 0;
}
```

What does this program print?

- A. `-1 1`
- B. `0 1`
- C. `0 0`
- D. `-1 0`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:** A `do-while` runs its body **before** checking the condition, so the body runs at least once.

| Step | Statement executed | Variables after execution |
|---|---|---|
| 1 | body: `count++` | n = 0, count = 1 |
| 2 | test `n-- > 0`: compares the **old** n (0) with 0 → false; n is still decremented | n = −1 |
| 3 | loop exits → `printf` | `-1 1` |

**Why the other options fail:**
- B: Forgets that the post-decrement still happens even though the comparison failed.
- C, D: Treat the loop as a `while`, which would check before the first iteration.

**Concept recap:**
- `while (c) S;` runs S **0 or more** times.
- `do S while (c);` runs S **1 or more** times. Note the required `;` after the `while (...)`.
- The side effect inside the condition runs on **every** test, including the final failing one.

**Trap alert:** After `while (n-- > 0) {...}` finishes, n is **−1**, not 0, because the last (failing) test still decremented it.

**Interview angle:** *"Give a natural use of do-while."* "Input validation: read, then repeat while the input is invalid. Also, in macros, `#define M(x) do { … } while (0)` makes a multi-statement macro behave like a single statement (covered in M10)."

---

## Q49. [M3 | Short-circuit in a loop condition | Medium | MCQ | [C] · C11 | GATE+OA] · Type T6 · P04

```c
#include <stdio.h>
int calls = 0;
int ok(int x) { calls++; return x < 3; }
int main(void) {
    int i = 0;
    while (i < 2 || ok(i))
        i++;
    printf("%d %d\n", i, calls);
    return 0;
}
```

What does this program print?

- A. `3 4`
- B. `2 0`
- C. `3 2`
- D. `3 3`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** `||` evaluates `ok(i)` **only when** `i < 2` is false.

| Test # | i | `i < 2` | `ok(i)` called? | `ok` returns | Condition | calls | Action |
|---|---|---|---|---|---|---|---|
| 1 | 0 | true | **no** (short-circuit) | — | true | 0 | i → 1 |
| 2 | 1 | true | no | — | true | 0 | i → 2 |
| 3 | 2 | false | **yes** | 2 < 3 = 1 | true | 1 | i → 3 |
| 4 | 3 | false | **yes** | 3 < 3 = 0 | false | 2 | exit |

Prints `3 2`.

**Why the other options fail:**
- A: Calls `ok` on every test (4 tests → 4 calls).
- B: Stops when `i < 2` first fails, ignoring the right operand of `||`.
- D: Counts one call too many, e.g. by also calling `ok` on test 2.

**Concept recap:** Putting a cheap check before an expensive (or side-effecting) one in `&&`/`||` is both a correctness tool (guards) and an optimisation. But side effects inside conditions make code hard to trace, as this question shows.

**Trap alert:** Counting function calls in a short-circuit expression is a favourite OA trick. Build a table with one row per evaluation of the condition.

**Interview angle:** *"How do short-circuit operators help performance?"* "Order conditions from cheapest or most-likely-to-decide to most expensive. For example, check a cached flag before doing a hash lookup or a system call."

---

## Q50. [M3 | Which loops terminate? (unsigned counters) | Hard | MSQ | [C] · C11 | All] · Type T3 · P11

Select **all** loops that are **guaranteed to terminate**. Each body is empty or does not touch the loop variable.

- A. `for (unsigned i = 5; i >= 0; i--) { }`
- B. `for (unsigned i = 5; i > 0; i--) { }`
- C. `for (unsigned i = 5; i-- > 0; ) { }`
- D. `for (int i = 0; i != 10; i += 3) { }`
- E. `for (unsigned char c = 0; c < 256; c++) { }`

**Answer:** B, C

**Behaviour class per option:** A = well-defined, infinite · B = terminates · C = terminates · D = **UB** · E = well-defined, infinite

**Why each option:**
- **A ✗:** An unsigned value is always ≥ 0. When i goes from 0 to `i--`, it wraps (well-defined) to UINT_MAX, and the loop never ends. GCC warns "comparison of unsigned expression in '>= 0' is always true".
- **B ✓:** i = 5, 4, 3, 2, 1, then 0 fails `i > 0`. Five iterations; the wrap never happens.
- **C ✓:** This is the idiomatic way to count down an unsigned index. The test uses the old value and then decrements: the tests see 5, 4, 3, 2, 1, so the **body sees i = 4, 3, 2, 1, 0**. The final test sees 0, fails, and leaves i = UINT_MAX, but it is never used. (Verified: 5 iterations.) This is how you loop `n-1 … 0` over array indices with `size_t`.
- **D ✗ (UB):** i takes 0, 3, 6, 9, 12, … and never equals 10. Eventually `i += 3` overflows a signed `int` → **undefined behaviour**. It isn't merely a long loop; the compiler may assume the overflow never happens.
- **E ✗:** `c` is promoted to `int` for the comparison, and its maximum is 255 < 256, so the condition is always true. `c++` wraps from 255 to 0 (well-defined for unsigned types).

**Concept recap:** Before trusting a loop, check:
1. Can the condition become false at all, given the variable's type?
2. Can the update skip past the exit value? (Use `<`, not `!=`, for strided loops.)
3. Can the update overflow? Signed → UB; unsigned → wrap.

**Trap alert:** `for (size_t i = n - 1; i >= 0; i--)` is the same bug as A, and it appears constantly in real reverse loops over arrays.

**Interview angle:** *"Write a reverse loop over a `std::vector` using `size_t`."* `for (size_t i = v.size(); i-- > 0; ) use(v[i]);`. Alternatives: reverse iterators, or a signed index with `std::ssize(v)` in C++20.

---

## Q51. [M4 | Swapping by value vs through pointers | Easy | MCQ | [C] · C11 | GATE+OA] · Type T5 · P15

```c
#include <stdio.h>
void swap1(int a, int b) { int t = a; a = b; b = t; }
void swap2(int *a, int *b) { int t = *a; *a = *b; *b = t; }
int main(void) {
    int x = 1, y = 2;
    swap1(x, y);
    printf("%d %d ", x, y);
    swap2(&x, &y);
    printf("%d %d\n", x, y);
    return 0;
}
```

What does this program print?

- A. `2 1 2 1`
- B. `1 2 2 1`
- C. `2 1 1 2`
- D. `1 2 1 2`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:** C passes **every** argument **by value**: the parameter is a new local variable initialised with a *copy* of the argument. To change the caller's variable, you pass a **copy of its address**, and the callee writes through that address.

**Memory table**

| Variable | Address (illustrative) | Stored value / Points to |
|---|---|---|
| main: x | 0x1000 | 1 |
| main: y | 0x1004 | 2 |
| swap1: a | 0x2000 (swap1's frame) | copy of x = 1 → after the swap, 2 |
| swap1: b | 0x2004 | copy of y = 2 → after the swap, 1 |
| *(swap1 returns: its frame is gone; x, y are unchanged)* | | |
| swap2: a | 0x2000 | 0x1000 → points to **x** |
| swap2: b | 0x2008 | 0x1004 → points to **y** |
| after `*a = *b; *b = t;` | | x = 2, y = 1 |

| Line | Statement executed | x | y |
|---|---|---|---|
| 5 | init | 1 | 2 |
| 6 | `swap1(x, y)` (swaps the copies only) | 1 | 2 |
| 7 | prints `1 2 ` | | |
| 8 | `swap2(&x, &y)` (swaps through the pointers) | 2 | 1 |
| 9 | prints `2 1` | | |

**Why the other options fail:**
- A: Assumes `swap1` swaps the caller's variables.
- C: Gets both calls backwards.
- D: Assumes neither call works.

**Concept recap:** "Pass by pointer" in C is *still* pass by value; the value passed happens to be an address. C has no references. C++ adds real reference parameters (`void swap(int& a, int& b)`).

**Trap alert:** Calling `swap2(x, y)` (forgetting the `&`) is a compile error: `int` is not compatible with `int *`. Pre-standard compilers accepted it and crashed at run time.

**Interview angle:** *"Is C pass-by-reference?"* "No. It's always pass-by-value. Pointers let you simulate reference semantics." Follow-up: *"Swap without a temporary?"* The XOR trick `a ^= b; b ^= a; a ^= b;` breaks when `a` and `b` are the same object (it zeroes it). Prefer a temporary.

---

## Q52. [M4 | A pointer parameter is itself a copy | Medium | MCQ | [C] · C11 | All] · Type T5 · P15

```c
#include <stdio.h>
int g = 100;
void f(int *p) {
    p = &g;
    *p = 5;
}
int main(void) {
    int x = 1;
    int *q = &x;
    f(q);
    printf("%d %d %d\n", x, *q, g);
    return 0;
}
```

What does this program print?

- A. `5 5 100`
- B. `1 5 5`
- C. `5 5 5`
- D. `1 1 5`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:** `p` is a **copy** of `q`. Making `p` point somewhere else changes only the copy. The write `*p = 5` then goes to whatever `p` points to **at that moment**, which is `g`.

**Memory table**

| Step | Variable | Address | Stored value / Points to |
|---|---|---|---|
| start | g | 0x4000 | 100 |
| | x | 0x1000 | 1 |
| | q | 0x1008 | 0x1000 (→ x) |
| enter f | p (f's frame) | 0x2000 | 0x1000 (→ x), a copy of q |
| `p = &g` | p | 0x2000 | **0x4000 (→ g)**; q is unchanged |
| `*p = 5` | g | 0x4000 | **5** |
| return | p is destroyed | | q still → x, x = 1 |

Prints `x = 1`, `*q = x = 1`, `g = 5` → `1 1 5`.

**Why the other options fail:**
- A: Assumes the write went to x.
- B: Assumes reassigning p also changed q, so that `*q` now reads g.
- C: Assumes everything aliases everything.

**Concept recap:**
- To let a function change **what the caller's variable holds**, pass a pointer to it (`int *`).
- To let a function change **where the caller's pointer points**, pass a pointer to the pointer (`int **`).

**Trap alert:** This bug is common in linked-list code: `void insert_head(Node *head, int v) { head = new_node; }` doesn't change the caller's head. The fix is `Node **head` (covered in M6 and M24).

**Interview angle:** *"Why does my `allocate(int *p) { p = malloc(...); }` leave p NULL in main?"* "The pointer is passed by value, so you change a copy, and you also leak the allocation. Pass `int **p` and write `*p = malloc(...)`, or return the pointer."

---

## Q53. [M4 | Activation records during recursion | Medium | MSQ | [C] · C11 | All] · Type — concept + code reading · P15

```c
int fact(int n) {
    if (n <= 1) return 1;
    return n * fact(n - 1);
}
/* main calls fact(4) */
```

Select **all** true statements about the call `fact(4)`.

- A. At the deepest point, four activation records of `fact` exist on the call stack at the same time.
- B. Each activation record holds its own separate copy of `n`.
- C. When `fact(1)` returns, all four activation records are released together.
- D. Because `fact` is recursive, `n` must be stored in static memory rather than on the stack.
- E. The stack space used grows linearly with `n`, so a large enough `n` can overflow the stack.

**Answer:** A, B, E

**Why each option:**

Call stack, frame by frame (the top of the stack is the last row):

| Phase | Stack contents (bottom → top) | Notes |
|---|---|---|
| main calls `fact(4)` | main · fact(n=4) | 4 waits on `4 * fact(3)` |
| | main · fact(4) · fact(n=3) | |
| | main · fact(4) · fact(3) · fact(n=2) | |
| deepest | main · fact(4) · fact(3) · fact(2) · **fact(n=1)** | 4 fact frames: **A ✓** |
| fact(1) returns 1 | main · fact(4) · fact(3) · fact(2) | fact(2) computes 2 × 1 = 2 |
| fact(2) returns 2 | main · fact(4) · fact(3) | 3 × 2 = 6 |
| fact(3) returns 6 | main · fact(4) | 4 × 6 = 24 |
| fact(4) returns 24 | main | |

- **A ✓:** See the "deepest" row.
- **B ✓:** Every call gets a fresh frame with its own parameters and locals. That's what makes recursion work.
- **C ✗:** Frames are popped **one at a time**, as each call returns (LIFO order).
- **D ✗:** The reverse is true. Recursion *requires* per-call storage. (A `static` local would be one shared copy; see Q57–Q58.)
- **E ✓:** The depth is n, so the stack use is O(n). The default Linux main-thread stack is typically 8 MB, so very deep recursion can crash with SIGSEGV (stack overflow).

**Concept recap:** An **activation record** (stack frame) usually holds the return address, the saved frame pointer, the parameters (or spilled copies of them), locals, and saved registers. Calling a function pushes a frame; returning pops it.

**Trap alert:** "Recursion uses the heap." It doesn't. Ordinary function calls use the stack. Only explicit `malloc`/`new` uses the heap.

**Interview angle:** *"Recursion vs iteration?"* "Recursion is clearer for tree-shaped problems but costs O(depth) stack space and call overhead. Convert to iteration (or tail recursion plus optimisation) when the depth can be large or attacker-controlled."

---

## Q54. [M4 | Recursion with output between two calls | Medium | MCQ | [C] · C11 | GATE] · Type T6 · P12

```c
#include <stdio.h>
void f(int n) {
    if (n <= 0) return;
    f(n - 2);
    printf("%d ", n);
    f(n - 1);
}
int main(void) {
    f(3);
    printf("\n");
    return 0;
}
```

What does this program print?

- A. `1 3 2 1`
- B. `3 1 2 1`
- C. `1 3 1 2`
- D. `1 2 3 1`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:** Each call does three things in order: **(left call `f(n-2)`) → print n → (right call `f(n-1)`)**. This is an *in-order* traversal of the call tree.

Call tree:
```
            f(3)
          /  |   \
      f(1)  "3"  f(2)
     / | \       / |  \
 f(-1)"1"f(0) f(0) "2" f(1)
                       / | \
                   f(-1)"1" f(0)
```

| Order | Call executing | Action | Output so far |
|---|---|---|---|
| 1 | f(3) | calls f(1) | |
| 2 | f(1) | calls f(−1) → returns | |
| 3 | f(1) | prints 1 | `1 ` |
| 4 | f(1) | calls f(0) → returns; f(1) returns | |
| 5 | f(3) | prints 3 | `1 3 ` |
| 6 | f(3) | calls f(2) | |
| 7 | f(2) | calls f(0) → returns; prints 2 | `1 3 2 ` |
| 8 | f(2) | calls f(1): f(−1), prints 1, f(0) | `1 3 2 1 ` |

**Why the other options fail:**
- B: Prints before the first recursive call (a pre-order reading).
- C: Swaps the order inside the right subtree.
- D: Treats the two recursive calls as if they ran before any printing.

**Concept recap:** Where the `printf` sits relative to the recursive calls decides the order:

| Position of the print | Traversal order |
|---|---|
| before both calls | pre-order |
| between the calls | in-order |
| after both calls | post-order |

**Trap alert:** Under time pressure, **draw the call tree**. Tracing in your head mixes up which call is "left".

**Interview angle:** *"How many times is `printf` called by f(n)?"* It satisfies P(n) = P(n−2) + 1 + P(n−1) with P(≤0) = 0. That's Fibonacci-like, so exponential: P(n) = Fib(n+2) − 1 (here P(3) = 4). Recognising this recurrence is the point of the question.

---

## Q55. [M4 | Post-order printing while building a value | Hard | MCQ | [C] · C11 | GATE+OA] · Type T6 · P12

```c
#include <stdio.h>
int fun(int n) {
    if (n < 10) { printf("%d ", n); return n; }
    int r = fun(n / 10) * 2 + n % 10;
    printf("%d ", r);
    return r;
}
int main(void) {
    int v = fun(352);
    printf("| %d\n", v);
    return 0;
}
```

What does this program print?

- A. `24 11 3 | 24`
- B. `3 11 24 | 3`
- C. `3 11 24 | 24`
- D. `3 8 10 | 10`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** Each call first recurses on `n / 10` (all digits except the last), **then** combines and prints. So the printing happens **on the way back up**: innermost first.

| Frame | n | Action | r / return | Output so far |
|---|---|---|---|---|
| fun(352) | 352 | needs fun(35) | (waiting) | |
| fun(35) | 35 | needs fun(3) | (waiting) | |
| fun(3) | 3 | base case: prints 3, returns 3 | **3** | `3 ` |
| fun(35) | 35 | r = 3 × 2 + 35 % 10 = 6 + 5 | **11**, printed | `3 11 ` |
| fun(352) | 352 | r = 11 × 2 + 352 % 10 = 22 + 2 | **24**, printed | `3 11 24 ` |
| main | | v = 24 | | `3 11 24 \| 24` |

**Why the other options fail:**
- A: Prints in call order (pre-order). The prints come *after* the recursive call.
- B: Takes the base-case value as the final return value.
- D: Adds digits instead of doubling (3, then 3 + 5 = 8, then 8 + 2 = 10).

**Concept recap:** The pattern `fun(n / 10) … n % 10` processes digits **most-significant first** on the way back. That is the standard way to print a number's digits in order with recursion.

**Trap alert:** Main stores the result in `v` before printing it. Writing `printf("| %d\n", fun(352))` would *still* run `fun` before `printf` prints anything, because arguments are evaluated before the call. But with *two* arguments that print, their relative order would be unspecified (see Q56).

**Interview angle:** *"Print an integer's digits using recursion and no arrays."* `void p(int n) { if (n >= 10) p(n / 10); putchar('0' + n % 10); }` (handle negative numbers separately). The recursion depth is the number of digits, O(log₁₀ n).

---

## Q56. [M4 | Counting calls in naive Fibonacci | Hard | MCQ | [C] · C11 | GATE+OA] · Type T6 · P12

```c
#include <stdio.h>
int calls = 0;
int fib(int n) {
    calls++;
    if (n < 2) return n;
    return fib(n - 1) + fib(n - 2);
}
int main(void) {
    int r = fib(6);
    printf("%d %d\n", r, calls);
    return 0;
}
```

What does this program print?

- A. `8 15`
- B. `8 25`
- C. `13 25`
- D. `8 41`

**Answer:** B

**Behaviour class:** Well-defined. The order in which `fib(n-1)` and `fib(n-2)` are evaluated is unspecified, but both the sum and the call count are the same either way.

**Why:** Let F(n) be the value and C(n) the number of calls (including the top-level one).
- F(n) = F(n−1) + F(n−2), with F(0) = 0 and F(1) = 1.
- C(n) = 1 + C(n−1) + C(n−2), with C(0) = C(1) = 1.

| n | F(n) | C(n) = 1 + C(n−1) + C(n−2) |
|---|---|---|
| 0 | 0 | 1 |
| 1 | 1 | 1 |
| 2 | 1 | 1 + 1 + 1 = 3 |
| 3 | 2 | 1 + 3 + 1 = 5 |
| 4 | 3 | 1 + 5 + 3 = 9 |
| 5 | 5 | 1 + 9 + 5 = 15 |
| 6 | **8** | 1 + 15 + 9 = **25** |

Shortcut: C(n) = 2·Fib(n+1) − 1 = 2 × 13 − 1 = 25.

**Why the other options fail:**
- A: That's C(5). Off by one level.
- C: 13 is Fib(7); the value is off by one index.
- D: That's C(7) = 1 + 25 + 15 = 41.

**Concept recap:** Naive Fibonacci makes Θ(φⁿ) ≈ Θ(1.618ⁿ) calls because it recomputes the same subproblems repeatedly. Memoisation or a simple loop brings it down to O(n).

**Trap alert:** Main deliberately reads `calls` in a **separate statement**. Writing `printf("%d %d\n", fib(6), calls);` would be a bug: the two arguments are evaluated in an **unspecified order**, so `calls` might be read before `fib` runs (printing 0) or after (25). In C11 this is unspecified, not UB, because the calls to `fib` are *indeterminately sequenced* with respect to the read of `calls`.

**Interview angle:** *"How many calls does naive fib(n) make, and how do you fix it?"* "2·Fib(n+1) − 1, which is exponential. Fix it with memoisation (top-down DP), a bottom-up loop (O(n) time, O(1) space), or matrix exponentiation (O(log n))."

---

## Q57. [M4 | static local vs automatic local | Easy | MCQ | [C] · C11 | GATE+OA] · Type T6 · P13

```c
#include <stdio.h>
int f(int x) {
    static int s = 10;
    int a = 0;
    s += x;
    a += x;
    return s + a;
}
int main(void) {
    int p = f(1);
    int q = f(2);
    int r = f(3);
    printf("%d %d %d\n", p, q, r);
    return 0;
}
```

What does this program print?

- A. `12 14 16`
- B. `11 12 13`
- C. `12 15 18`
- D. `12 15 19`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:** `static int s = 10;` is initialised **once**, before the program starts, and keeps its value between calls. `int a = 0;` is re-created and set to 0 on **every** call.

| Call | s before | `s += x` | a | return s + a |
|---|---|---|---|---|
| f(1) | 10 | 11 | 1 | 12 |
| f(2) | 11 | 13 | 2 | 15 |
| f(3) | 13 | 16 | 3 | 19 |

Prints `12 15 19`.

**Why the other options fail:**
- A: Resets s to 10 on each call (10 + 1 + 1, 10 + 2 + 2, 10 + 3 + 3). That's how an automatic variable would behave.
- B: Ignores `a`, and resets s.
- C: Keeps s but miscounts the increments (adds 3 each time).

**Concept recap:** A block-scope `static` has the **scope** of a local (visible only inside `f`) and the **lifetime** of a global (it exists from program start to end). Its initializer must be a constant expression in C.

**Trap alert:** The line `static int s = 10;` is **not** re-executed on each call. Treat it as a one-time setup, not an assignment.

**Interview angle:** *"Give a real use of a static local."* "A call counter, a lazily initialised cache or lookup table, or a function-local singleton. In C++ (since C++11), initialisation of a function-local static is thread-safe ('magic statics'), which makes it the standard Meyers-singleton pattern." Follow-up: *"Downside?"* The function is no longer reentrant or thread-safe in C. It has hidden state.

---

## Q58. [M4 | static local inside recursion, called twice | Hard | MCQ | [C] · C11 | GATE] · Type T6 · P13, P12

```c
#include <stdio.h>
int f(int n) {
    static int r = 0;
    if (n <= 0) return 1;
    if (n > 3) {
        r += n;
        return f(n - 2) + 1;
    }
    return f(n - 1) + r;
}
int main(void) {
    int a = f(6);
    int b = f(6);
    printf("%d %d\n", a, b);
    return 0;
}
```

What does this program print?

- A. `23 43`
- B. `23 23`
- C. `15 35`
- D. `23 39`

**Answer:** A

**Behaviour class:** Well-defined. In `f(n - 1) + r`, the order of the two operands is unspecified, but the call `f(n-1)` (with n − 1 ≤ 2) never changes `r`, so either order gives the same sum. The question was built so this holds.

**Why:** There is **one** `r`, shared by every activation and by both top-level calls.

**First call, `a = f(6)`** (r starts at 0):

| Frame | n | Branch | r after | Returns |
|---|---|---|---|---|
| f(6) | 6 | n > 3: r += 6 | 6 | f(4) + 1 |
| f(4) | 4 | n > 3: r += 4 | 10 | f(2) + 1 |
| f(2) | 2 | else | 10 | f(1) + r |
| f(1) | 1 | else | 10 | f(0) + r |
| f(0) | 0 | base case | 10 | **1** |
| f(1) | | 1 + 10 | | **11** |
| f(2) | | 11 + 10 | | **21** |
| f(4) | | 21 + 1 | | **22** |
| f(6) | | 22 + 1 | | **23** |

**Second call, `b = f(6)`** (r starts at **10**):

| Frame | r after | Returns |
|---|---|---|
| f(6): r += 6 | 16 | f(4) + 1 |
| f(4): r += 4 | 20 | f(2) + 1 |
| f(0) | 20 | 1 |
| f(1) | | 1 + 20 = 21 |
| f(2) | | 21 + 20 = 41 |
| f(4) | | 42 |
| f(6) | | **43** |

Prints `23 43`.

**Why the other options fail:**
- B: Assumes `r` resets between the two calls from `main`.
- C: Uses r = 6, ignoring the `r += 4` done by f(4): f(1) = 7, f(2) = 13, f(4) = 14, f(6) = 15. But f(4) updates `r` *before* calling f(2), so the reads in f(1) and f(2) see 10.
- D: Resets `r` partway, or miscounts the second accumulation.

**Concept recap:** Every activation of a recursive function shares its `static` locals, and so does every later top-level call. A static local is effectively a global with restricted visibility.

**Trap alert:** If a function modifies a `static` in the recursive call and *also* reads that static in the same expression (`return f(n-1) + r;` where `f(n-1)` changes `r`), the result depends on the unspecified evaluation order. Always check whether the recursive call can modify the shared variable.

**Interview angle:** *"Why are static locals dangerous in recursive or multi-threaded code?"* "Every caller shares them, so a nested call or another thread can change the value mid-computation. Pass state through parameters (an accumulator) instead."

---

## Q59. [M4 | Indirect (mutual) recursion | Medium | MCQ | [C] · C11 | GATE+OA] · Type T6

```c
#include <stdio.h>
void odd(int n);
void even(int n) {
    if (n <= 0) return;
    printf("E%d ", n);
    odd(n - 1);
}
void odd(int n) {
    if (n <= 0) return;
    printf("O%d ", n);
    even(n / 2);
}
int main(void) {
    even(9);
    printf("\n");
    return 0;
}
```

What does this program print?

- A. `E9 O8 E7 O6 E5 O4 E3 O2 E1`
- B. `E9 O4 E3 O1`
- C. `E9 O8 E4 O3 E1`
- D. `E9 O8 E4 O2 E1`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** `even` calls `odd(n - 1)`, and `odd` calls `even(n / 2)` (integer division).

| Call | n | Prints | Next call |
|---|---|---|---|
| even | 9 | `E9` | odd(8) |
| odd | 8 | `O8` | even(4) |
| even | 4 | `E4` | odd(3) |
| odd | 3 | `O3` | even(1) (3 / 2 = 1) |
| even | 1 | `E1` | odd(0) → returns immediately |

Output: `E9 O8 E4 O3 E1`. The stack depth reaches 6 frames (main plus 5).

**Why the other options fail:**
- A: Ignores the `/ 2` and decrements every time.
- B: Halves in `even` instead of in `odd`.
- D: Treats `odd(3)`'s n as 2. It gets 4 − 1 = 3.

**Concept recap:** **Indirect recursion**: A calls B and B calls A. The forward declaration `void odd(int n);` is **required** in C99 and later, because `even` calls `odd` before `odd` is defined, and implicit function declarations were removed in C99.

**Trap alert:** Without the prototype, GCC 14+ refuses the code ("implicit declaration of function 'odd'"), and older GCCs only warn. This is covered in Batch 4.

**Interview angle:** *"Real uses of mutual recursion?"* "Recursive-descent parsers (`expr` → `term` → `factor` → `expr`), state machines, and the classic `is_even`/`is_odd` pair."

---

## Q60. [M4 | Tail recursion | Easy | MSQ | [C] · C11 | OA+Interview] · Type T3

Select **all** functions whose recursive call is a **tail call** (the call is the very last action, and its result is returned unchanged).

- A. `int fact(int n) { return n <= 1 ? 1 : n * fact(n - 1); }`
- B. `int fact_acc(int n, int acc) { return n <= 1 ? acc : fact_acc(n - 1, n * acc); }`
- C. `int sum(int n) { if (n == 0) return 0; return sum(n - 1) + n; }`
- D. `int gcd(int a, int b) { return b == 0 ? a : gcd(b, a % b); }`
- E. `void show(int n) { if (n == 0) return; show(n - 1); printf("%d ", n); }`

**Answer:** B, D

**Why each option:**
- **A ✗:** After `fact(n - 1)` returns, the caller still has to **multiply by n**. That pending work keeps its frame alive.
- **B ✓:** The pending multiplication was moved into the **accumulator** parameter. The result of the recursive call is returned directly.
- **C ✗:** `+ n` happens after the call returns.
- **D ✓:** Euclid's algorithm. The call is the last action.
- **E ✗:** `printf` runs after the recursive call returns.

**Concept recap:**
- A compiler can turn a tail call into a **jump**, reusing the current frame (tail-call optimisation), which makes stack use O(1).
- **The C and C++ standards do not require this.** GCC and Clang usually do it at `-O2`, but not at `-O0`.
- The accumulator technique (option B) is the standard way to convert linear recursion into tail recursion.

**Trap alert:** A call that is written last but whose result is then **used** (`return 1 + f(n-1);`) is **not** a tail call.

**Interview angle:** *"Can I rely on tail-call optimisation in C++?"* "No, it isn't guaranteed. Destructors of local objects that still have to run after the call also prevent it. For guaranteed constant stack space, write a loop."

---

## Batch 3 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 41 | M3 | E | MCQ | C | Well-defined |
| 42 | M3 | M | MCQ | B | Well-defined |
| 43 | M3 | M | MCQ | D | Well-defined |
| 44 | M3 | E | MCQ | A | Well-defined |
| 45 | M3 | H | MCQ | C | Well-defined |
| 46 | M3 | E | MCQ | D | Well-defined |
| 47 | M3 | M | MCQ | B | Well-defined |
| 48 | M3 | E | MCQ | A | Well-defined |
| 49 | M3 | M | MCQ | C | Well-defined |
| 50 | M3 | H | MSQ | B, C | Terminates / infinite / **UB** mix |
| 51 | M4 | E | MCQ | B | Well-defined |
| 52 | M4 | M | MCQ | D | Well-defined |
| 53 | M4 | M | MSQ | A, B, E | — (concept) |
| 54 | M4 | M | MCQ | A | Well-defined |
| 55 | M4 | H | MCQ | C | Well-defined |
| 56 | M4 | H | MCQ | B | Well-defined |
| 57 | M4 | E | MCQ | D | Well-defined |
| 58 | M4 | H | MCQ | A | Well-defined |
| 59 | M4 | M | MCQ | C | Well-defined |
| 60 | M4 | E | MSQ | B, D | — (concept) |

- **Difficulty:** 7 E / 8 M / 5 H.
  - M3 is 4/4/2, as planned.
  - M4's part here is 3/4/3. The remaining M4 questions in Batch 4 will be 2/3/3, making M4's total 5/7/6.
- **Format:** 17 MCQ / 3 MSQ.
  - M3: 9/1, as planned.
  - M4 so far: 8/2. Batch 4 adds 7/1, making M4's total 15/3.
- **MCQ answer keys:** A×4, B×4, C×5, D×4. Running total Q1–Q60: A 12, B 12, C 12, D 12.
- **Avoiding unspecified evaluation order:** Q56 and Q58 were built so that it doesn't change the answer, and both explanations say why.

---

**Completed: 60 / 400 Remaining: 340**

**Ledger, newly ✓ this batch:**
- P12 ×4 (Q54, Q55, Q56, Q58) ★ complete
- P13 ×2 (Q57, Q58) ★ complete
- P15 ×3 (Q51, Q52, Q53)
- P04 → 2 (Q49)
- P11 → 4 (Q50)

**M3 complete** (Q41–Q50, 10 questions). **M4 in progress** (10 of 18 done).

*Say **"Continue."** for Batch 4 (Q61–Q80: the rest of M4 — function pointers and callbacks, returning a local's address, uninitialised reads, variadic functions, implicit declarations, unspecified argument order, more recursion — then the start of M5 arrays).*
