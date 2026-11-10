# C/C++ MCQ/MSQ Bank — Batch 2 (Q21–Q40)

**What this batch covers (all M2):**
- integer vs floating-point division
- the usual arithmetic conversions across `int` / `unsigned` / `long`
- the types of expressions after promotion
- chained relational operators
- short-circuit evaluation combined with precedence
- pre/post increment on scalars and array elements
- the comma operator and the conditional operator `?:` (including C vs C++ lvalue rules)
- bitwise operations and right shifts on negative numbers
- signed overflow (UB), signed vs unsigned comparisons, the `sizeof` comparison trap
- unsequenced modifications (UB)
- floating → integer truncation, and out-of-range float → int conversion (UB)
- `/` and `%` with negative operands
- shift-count rules
- logical vs bitwise operators

**Patterns:** P02, P03, P04, P06, P07, P08, P09, P11, P49, P50

> **Conventions for every question in this batch (unless the question says otherwise):**
> - C11, `gcc -std=c11`, 64-bit Linux / GCC / x86-64
> - `sizeof(short)=2`, `sizeof(int)=4`, `sizeof(long)=8`, `size_t` = `unsigned long`
> - plain `char` is signed; two's-complement integers; IEEE-754 `double`; ASCII
> - `stdio.h` (and `limits.h` where `INT_MAX`/`INT_MIN` appear) are included
> - optimisation level does not matter unless stated

---

## Q21. [M2 | Integer vs floating division | Easy | MCQ | [C] · C11 | GATE] · Type T1 · P02

```c
#include <stdio.h>
int main(void) {
    int a = 5, b = 2;
    double d1 = a / b;
    double d2 = (double)a / b;
    double d3 = (double)(a / b);
    int x = 7 / 2 * 2.0;
    printf("%.1f %.1f %.1f %d\n", d1, d2, d3, x);
    return 0;
}
```

What does this program print?

- A. `2.5 2.5 2.5 7`
- B. `2.0 2.5 2.0 6`
- C. `2.0 2.5 2.5 6`
- D. `2.0 2.5 2.0 7`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:** The type of an arithmetic operation is decided by its **operands**, never by the variable the result is stored into.

| Line | Statement executed | Evaluation | Variables after execution |
|---|---|---|---|
| 3 | `int a = 5, b = 2;` | — | a=5, b=2 |
| 4 | `d1 = a / b` | `int / int` → integer division, truncated: **2**, then converted to 2.0 | d1 = 2.0 |
| 5 | `d2 = (double)a / b` | the cast applies to `a` only (unary operators bind tighter): `5.0 / b`. `b` is converted to `double` (usual arithmetic conversion) → 2.5 | d2 = 2.5 |
| 6 | `d3 = (double)(a / b)` | the division happens first, in `int`: 2, then cast → 2.0 | d3 = 2.0 |
| 7 | `x = 7 / 2 * 2.0` | left to right: `7 / 2` = 3 (int), then `3 * 2.0` = 6.0 (double), converted to `int` | x = 6 |
| 8 | `printf` | — | `2.0 2.5 2.0 6` |

**Why the other options fail:**
- A: Assumes the `double` destination makes the division floating-point.
- C: Thinks casting *after* the integer division can recover the .5. The fraction is already gone.
- D: Evaluates `7 * 2.0 / 2`. `/` and `*` are left-associative, so `7 / 2` happens first.

**Concept recap:** **Usual arithmetic conversions** — before a binary operator, both operands are converted to a common type:
1. If either operand is `long double` / `double` / `float`, convert the other to that type.
2. Otherwise, apply integer promotions, then convert to the higher-ranked type (with signed/unsigned rules; see Q22).

**Trap alert:** `average = sum / count;` with `int` operands truncates even when `average` is `double`. Write `(double)sum / count`.

**Interview angle:** *"Why does `double avg = total / n;` give a whole number?"* "Both operands are int, so it's integer division; the conversion to double happens only on assignment. Cast one operand first." Follow-up: *"Which way does C integer division round?"* Toward zero (since C99). See Q38.

---

## Q22. [M2 | Usual arithmetic conversions: int / unsigned / long | Hard | MCQ | [C] · C11 | All] · Type T1 · P02, P11

Platform: LP64, so `int` = 4 bytes and `long` = 8 bytes.

```c
#include <stdio.h>
int main(void) {
    int i = -1;
    unsigned int u = 1;
    long l = -1;
    printf("%d %d\n", i < u, l < u);
    return 0;
}
```

What does this program print?

- A. `1 1`
- B. `0 0`
- C. `1 0`
- D. `0 1`

**Answer:** D

**Behaviour class:** Well-defined. The result depends on the stated type sizes (implementation-defined).

**Why:** The comparisons go through the usual arithmetic conversions (N1570 §6.3.1.8).

| Comparison | Operand types | Rule applied | Converted values | Result |
|---|---|---|---|---|
| `i < u` | `int` vs `unsigned int` (same rank) | Same rank, one unsigned → **both become `unsigned int`** | `-1` → 4294967295; `4294967295 < 1` | **0** |
| `l < u` | `long` (rank higher) vs `unsigned int` | The signed type has higher rank **and can represent every `unsigned int` value** (64 bits > 32) → **both become `long`** | `-1 < 1` | **1** |

Prints `0 1`.

**Why the other options fail:**
- A: Assumes ordinary mathematical comparison everywhere. It is broken for `int` vs `unsigned`.
- B: This is the answer on **LLP64 (Windows)**. There `long` is 32 bits, so it can't hold every `unsigned int` value, and both operands become `unsigned long` → `-1` turns huge again.
- C: Has the two rules backwards.

**Concept recap:** For two integer operands, after promotion:
1. Same type → done.
2. Both signed or both unsigned → convert to the higher rank.
3. The unsigned one has rank ≥ the signed one → convert to the **unsigned** type.
4. The signed one can represent all values of the unsigned one → convert to the **signed** type.
5. Otherwise → convert to the unsigned version of the signed type.

**Trap alert:** The *same* expression `l < u` gives different answers on Linux and Windows. That is a real portability bug, not a curiosity.

**Interview angle:** *"Why does `-1 < 1u` evaluate to false?"* "Same rank, one operand unsigned, so -1 is converted to UINT_MAX." Follow-up: *"How do you compare signed and unsigned values safely?"* Check the sign first (`i < 0 || (unsigned)i < u`), use C++20 `std::cmp_less`, and compile with `-Wsign-compare`.

---

## Q23. [M2 | Types of expressions after promotion | Medium | MSQ | [C] · C11 | GATE+OA] · Type T3 · P02

Given:

```c
char c = 1;
short s = 1;
unsigned char uc = 200;
```

Select **all** expressions whose type is **`int`**.

- A. `c + c`
- B. `uc * 2u`
- C. `s << 1L`
- D. `c + 1L`
- E. `-uc`

**Answer:** A, C, E

**Why each option:**
- **A ✓:** Both `char` operands are promoted to `int`. `int + int` is `int`.
- **B ✗:** `uc` becomes `int`. Then `int` vs `unsigned int` (same rank, one unsigned) → `unsigned int`.
- **C ✓:** For **shifts**, the result type is the **promoted left operand**. The usual arithmetic conversions are *not* applied between the two operands (N1570 §6.5.7p3). So `s` becomes `int`, and the `long` right operand does not affect the type.
- **D ✗:** `c` becomes `int`. Then `int` vs `long` → `long`.
- **E ✓:** Unary minus also promotes its operand: `uc` → `int` 200, so `-uc` = −200, an `int`. It is **not** a wrapped `unsigned char`.

(Checked with C11 `_Generic` on GCC: `int`, `unsigned`, `int`, `long`, `int`.)

**Concept recap:** **Integer promotion:** any operand narrower than `int` (`char`, `signed char`, `unsigned char`, `short`, `unsigned short`, bit-fields, `_Bool`) is converted to `int` if `int` can hold all its values, else to `unsigned int`. On this platform `int` holds all `unsigned short` values, so every small type becomes `int`. This is why `uc + (unsigned char)100` is the `int` 300 and doesn't wrap to 44.

**Trap alert:** `unsigned char a = 200, b = 100; if (a + b > 255)` is **true**, because the addition happens in `int`. Wrap-around happens only when the result is *stored* back into an `unsigned char`.

**Interview angle:** *"What's the type of `'a' + 1` in C?"* `int`. In C, `'a'` is already an `int`; in C++ it is a `char`, which gets promoted to `int`. Either way, character arithmetic is done in `int`.

---

## Q24. [M2 | Chained relational operators | Medium | MCQ | [C] · C11 | GATE] · Type T1 · P03

```c
#include <stdio.h>
int main(void) {
    int a = 1, b = 2, c = 3;
    int r = a < b < c;
    int s = c > b > a;
    printf("%d %d\n", r, s);
    return 0;
}
```

What does this program print?

- A. `1 0`
- B. `1 1`
- C. `0 0`
- D. `0 1`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:** Relational operators are **left-associative**, and each one yields the `int` 0 or 1. C has no mathematical "chaining".

| Line | Expression | Grouping | Step-by-step | Value |
|---|---|---|---|---|
| 4 | `a < b < c` | `(a < b) < c` | `(1 < 2)` = 1 → `1 < 3` = 1 | r = 1 |
| 5 | `c > b > a` | `(c > b) > a` | `(3 > 2)` = 1 → `1 > 1` = 0 | s = 0 |

Prints `1 0`.

**Why the other options fail:**
- B: Treats `c > b > a` as the mathematical "3 > 2 > 1", which is true. C compares the **boolean result** 1 against `a`.
- C, D: Misread the first comparison.

**Concept recap:** To express "a < b < c" in C, write `a < b && b < c`. The result of any relational, equality or logical operator has type `int` with value 0 or 1.

**Trap alert:** `a < b < c` happens to give the "right" answer for 1, 2, 3, which hides the bug. Try `a = 3, b = 2, c = 1`: `(3 < 2) < 1` → `0 < 1` → **1**, which is wrong.

**Interview angle:** *"Does `if (0 < x < 10)` check a range?"* "No. It's `(0 < x) < 10`, and both 0 and 1 are less than 10, so it's always true. Use `0 < x && x < 10`." Compilers warn with `-Wparentheses`.

---

## Q25. [M2 | Short-circuit evaluation + precedence | Hard | MCQ | [C] · C11 | GATE] · Type T1 · P04, P03

```c
#include <stdio.h>
int main(void) {
    int i = 0, j = 0, k = 0;
    int r = ++i || ++j && ++k;
    printf("%d %d %d %d\n", i, j, k, r);
    return 0;
}
```

What does this program print?

- A. `1 1 1 1`
- B. `1 1 0 1`
- C. `1 0 0 1`
- D. `1 0 1 1`

**Answer:** C

**Behaviour class:** Well-defined. There is a sequence point after the left operand of `||` and of `&&`.

**Why:** First parse, then evaluate. `&&` has **higher precedence** than `||`, so the expression is:

`r = (++i) || ((++j) && (++k))`

| Step | Action | Variables |
|---|---|---|
| 1 | Evaluate the left operand of `||`: `++i` → 1 (non-zero = true) | i=1 |
| 2 | `||` already has a true left operand → **the whole right operand `(++j && ++k)` is skipped** | j=0, k=0 |
| 3 | The result of `||` is 1 | r=1 |
| 4 | `printf` | `1 0 0 1` |

**Why the other options fail:**
- A: Ignores short-circuiting completely.
- B: Evaluates `++j` and skips only `++k`, as if `||` skipped only its *next* operand. It skips its entire right operand, which is the whole `&&` subexpression.
- D: Parses the expression as `(++i || ++j) && ++k` (wrong precedence). Then `++j` is skipped but `++k` runs, giving `1 0 1 1`.

**Concept recap:**
- `A && B`: evaluate A; if it is 0, the result is 0 and **B is not evaluated**.
- `A || B`: evaluate A; if it is non-zero, the result is 1 and **B is not evaluated**.
- Both guarantee left-to-right evaluation with a sequence point between the operands. That is why `p && p->next` is safe.

**Trap alert:** Precedence decides **grouping**; short-circuiting decides **what runs**. Always draw the parse tree first. A common mistake is to evaluate "left to right, skipping one operand".

**Interview angle:** *"Why is `if (ptr != NULL && ptr->val > 0)` safe?"* "`&&` short-circuits: if the pointer is null, the right side is never evaluated, so there is no null dereference." Follow-up: *"Does `&` short-circuit?"* No. Bitwise `&` evaluates both operands (in an unspecified order).

---

## Q26. [M2 | Pre/post increment (well-defined uses) | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P06

```c
#include <stdio.h>
int main(void) {
    int x = 5, y;
    y = x++ + 10;
    y = ++x + y;
    printf("%d %d\n", x, y);
    return 0;
}
```

What does this program print?

- A. `6 21`
- B. `7 21`
- C. `6 22`
- D. `7 22`

**Answer:** D

**Behaviour class:** Well-defined. Each statement modifies `x` at most once. In line 5, `y` is read only to compute the value being stored into `y`, which is allowed (N1570 §6.5p2).

**Why:**

| Line | Statement executed | Evaluation | Variables after execution |
|---|---|---|---|
| 3 | `int x = 5, y;` | — | x=5, y=? |
| 4 | `y = x++ + 10` | `x++` yields the **old** value 5, then x becomes 6 → 5 + 10 | x=6, y=15 |
| 5 | `y = ++x + y` | `++x` increments **first** and yields 7 → 7 + 15 | x=7, y=22 |
| 6 | `printf` | — | `7 22` |

**Why the other options fail:**
- A: Forgets the `++x` on line 5 (x stays 6), so y = 6 + 15 = 21.
- B: Uses `x++` semantics on line 5 too (yields 6, x becomes 7): 6 + 15 = 21.
- C: Increments x correctly but adds the wrong value.

**Concept recap:**
- `x++` (postfix): the expression's value is the old x; x is incremented as a side effect.
- `++x` (prefix): x is incremented, and the expression's value is the new x.
- The side effect is guaranteed to be complete by the next sequence point (here, the end of the full expression).

**Trap alert:** A single statement that modifies `x` **twice**, or reads it in another part of the same expression while modifying it, is UB (Q35, Q36). This question is safe because each statement modifies x only once.

**Interview angle:** *"Prefer `++i` or `i++` in a loop?"* "For `int` there's no difference; compilers generate identical code. For C++ iterators and class types, `++i` avoids making a copy of the old value, so it's the habit to have."

---

## Q27. [M2 | Increments on index vs element | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P06

```c
#include <stdio.h>
int main(void) {
    int a[] = {10, 20, 30, 40};
    int i = 0;
    int p = a[i++];
    int q = a[++i];
    int r = a[i]++;
    int s = ++a[i];
    printf("%d %d %d %d %d\n", p, q, r, s, i);
    return 0;
}
```

What does this program print?

- A. `10 30 31 32 2`
- B. `10 30 30 32 2`
- C. `10 20 30 32 2`
- D. `11 30 30 31 2`

**Answer:** B

**Behaviour class:** Well-defined. Each statement modifies only one object, once.

**Why:** Postfix `[]` and postfix `++` share the highest precedence level, and prefix `++` is applied to the whole postfix expression. So `a[i]++` and `++a[i]` both increment **the element**, while `a[i++]` and `a[++i]` increment **the index**.

| Line | Statement executed | What is incremented | a[] after | i | Result |
|---|---|---|---|---|---|
| 3–4 | init | — | {10, 20, 30, 40} | 0 | |
| 5 | `p = a[i++]` | index: uses i = 0, then i = 1 | {10, 20, 30, 40} | 1 | p = 10 |
| 6 | `q = a[++i]` | index: i = 2 first, then a[2] | {10, 20, 30, 40} | 2 | q = 30 |
| 7 | `r = a[i]++` | element a[2]: yields the old 30, then a[2] = 31 | {10, 20, **31**, 40} | 2 | r = 30 |
| 8 | `s = ++a[i]` | element a[2]: becomes 32 first, then yields 32 | {10, 20, **32**, 40} | 2 | s = 32 |

Prints `10 30 30 32 2`.

**Why the other options fail:**
- A: Uses the *new* value for the postfix `a[i]++` (r = 31).
- C: Treats `a[++i]` as if it used i before incrementing (giving 20).
- D: Treats `a[i++]` as incrementing the element.

**Concept recap:** Read the operand of `++` carefully. In `a[i++]` the operand is `i`. In `a[i]++` and `++a[i]` it is the element `a[i]`.

**Trap alert:** `a[i] = i++;` or `a[i++] = a[i];` *in one statement* is UB in C, because `i` is read in another subexpression while being modified, with no sequencing (Q36).

**Interview angle:** These are warm-up questions in OA coding rounds. The interview version is `*p++` vs `(*p)++` vs `++*p` with pointers (M6, pattern P29), which follows exactly the same reasoning.

---

## Q28. [M2 | Comma operator | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P07

```c
#include <stdio.h>
int main(void) {
    int a = (1, 2, 3);
    int b;
    b = 4, 5;
    int c = (b++, b + 10);
    printf("%d %d %d\n", a, b, c);
    return 0;
}
```

What does this program print?

- A. `3 5 15`
- B. `1 5 15`
- C. `3 6 16`
- D. `3 5 14`

**Answer:** A

**Behaviour class:** Well-defined. The comma operator has a sequence point after its left operand.

**Why:** The comma *operator* evaluates its left operand, discards the result, then evaluates its right operand, which is the value of the whole expression. It has the **lowest precedence** of all operators, below assignment.

| Line | Statement executed | Parsed as | Variables after execution |
|---|---|---|---|
| 3 | `int a = (1, 2, 3);` | the parentheses make it a comma expression → value 3 | a=3 |
| 5 | `b = 4, 5;` | `(b = 4), 5;` → assignment binds tighter; the `5` is evaluated and discarded | b=4 |
| 6 | `int c = (b++, b + 10);` | left: `b++` → b=5. **Sequence point.** Right: `b + 10` = 15 | b=5, c=15 |
| 7 | `printf` | — | `3 5 15` |

**Why the other options fail:**
- B: Takes the *first* operand of a comma expression. The value is the *last*.
- C: Thinks `b = 4, 5` assigns 5 (b would be 5, then 6 after `b++`, and c = 16).
- D: Uses the old b in `b + 10`. The sequence point guarantees the increment has already happened.

**Concept recap:** Not every comma is the comma operator. In `int x = 1, y = 2;` and `f(a, b)`, the commas are **separators**. The separator commas between function arguments are *not* sequence points (see Q36). Without parentheses, `int a = 1, 2, 3;` is a **syntax error**, because after `,` the declaration expects another declarator.

**Trap alert:** `for (i = 0, j = n - 1; i < j; i++, j--)`: these commas are operators, and this is the main legitimate use of the comma operator.

**Interview angle:** *"Is `f(i++, i++)` the same as using the comma operator?"* "No. Those commas are argument separators with no sequencing, so it's UB in C (unsequenced) and unspecified in C++17. The comma operator `(i++, i++)` is fully sequenced."

---

## Q29. [M2 | Conditional operator: associativity & lvalue-ness | Hard | MCQ | [C] · C11 | GATE+OA] · Type T2 · P07

```c
#include <stdio.h>
int main(void) {
    int a = 5, b = 10;
    int m = a > b ? a : b > 7 ? 100 : 200;
    printf("%d\n", m);
    (a < b ? a : b) = 0;            /* line 6 */
    printf("%d %d\n", a, b);
    return 0;
}
```

Compiled as **C11**, what happens?

- A. It prints `100` then `0 10`
- B. It prints `10` then `0 10`
- C. Compile error at line 6: in C, the result of `?:` is not an lvalue. (Compiled as C++, it would print `100` then `0 10`.)
- D. It prints `100` then `5 0`

**Answer:** C

**Behaviour class:** Ill-formed (a constraint violation: an assignment needs a modifiable lvalue on its left).

**Why:**

| Line | Analysis |
|---|---|
| 4 | `?:` is **right-associative**: `a > b ? a : (b > 7 ? 100 : 200)`. `5 > 10` is false → evaluate `b > 7 ? 100 : 200` → 100. If the program compiled, m would be 100. |
| 6 | In **C**, `?:` always yields a value (an rvalue), not an object. `(…) = 0` → GCC: *"lvalue required as left operand of assignment"*. The whole file fails to compile, so nothing is printed. |

In **C++**, if both branches are lvalues of the same type, `?:` yields an lvalue. Line 6 then assigns 0 to `a` (since `a < b`), and the program prints `100` / `0 10` (checked with g++).

**Why the other options fail:**
- A: Correct output for C++, but the question compiles as C.
- B: Groups the expression left-to-right, `(a > b ? a : b) > 7 ? …`, which doesn't even yield 10. Either way, line 6 does not compile in C.
- D: Assigns to the wrong variable, and still ignores the C rule.

**Concept recap:** Portable C version of line 6: `*(a < b ? &a : &b) = 0;`. Choose a *pointer* with `?:`, then dereference it.

**Trap alert:** Nested ternaries group from the right, like `if / else if / else`: `x ? A : y ? B : C` means `x ? A : (y ? B : C)`.

**Interview angle:** *"Name expressions that are lvalues in C++ but not in C."* `?:` with two lvalue branches, prefix `++x`, and assignment expressions such as `(a = b) = c`. All are lvalues in C++ and rvalues in C.

---

## Q30. [M2 | Bitwise operations on signed integers | Hard | MCQ | [C] · C11 | All] · Type T3 · P08

```c
#include <stdio.h>
int main(void) {
    int a = -8;
    int b = a >> 1;
    int c = a & 7;
    int d = ~a;
    int e = a ^ -1;
    printf("%d %d %d %d\n", b, c, d, e);
    return 0;
}
```

Which statement is correct?

- A. It prints `-4 0 7 7`. The value of `b` is **implementation-defined** (GCC documents an arithmetic shift); the other three follow from the stated two's-complement representation.
- B. The behaviour is undefined, because it right-shifts a negative value.
- C. It prints `2147483644 0 7 7`, and the standard guarantees a logical shift.
- D. It prints `-4 0 -7 7`.

**Answer:** A

**Behaviour class:** Implementation-defined (for `a >> 1`, and for the bitwise results on signed values, which depend on the representation; C11 allows three representations, and C23 mandates two's complement).

**Why:** −8 in 32-bit two's complement is `1111 … 1111 1000`.

| Line | Statement executed | Bits | Variables after execution |
|---|---|---|---|
| 3 | `a = -8` | `…1111 1000` | a = −8 |
| 4 | `a >> 1` | **Right-shifting a negative value is implementation-defined** (N1570 §6.5.7p5). GCC copies the sign bit in (arithmetic shift): `…1111 1100` | b = −4 |
| 5 | `a & 7` | `…1000 & …0111` = 0 | c = 0 |
| 6 | `~a` | flip all bits: `0000 … 0111`. In two's complement, `~x == -x - 1` → 8 − 1 | d = 7 |
| 7 | `a ^ -1` | XOR with all ones = `~a` | e = 7 |

Prints `-4 0 7 7`.

**Why the other options fail:**
- B: Right-shifting a **negative** value is implementation-defined, *not* UB. (**Left**-shifting a negative value *is* UB; see Q39.)
- C: A logical shift (filling with 0) would give `0x7FFFFFFC` = 2147483644. That is a *possible* implementation choice, but not GCC's, and not guaranteed.
- D: `~(-8)` is +7. People confuse `~x` with `-x`.

**Concept recap:** For portable bit manipulation, use **unsigned** types. All bitwise operations and shifts on unsigned values are fully defined (except a shift count ≥ the width).

**Trap alert:** `x >> 1` is **not** always `x / 2` for negative `x`. Division rounds toward zero (−7 / 2 = −3), while an arithmetic shift rounds toward −∞ (−7 >> 1 = −4 on GCC).

**Interview angle:** *"How do you check whether a number is a power of two?"* `x && !(x & (x - 1))`, with `x` unsigned. Follow-up: *"Why unsigned?"* To avoid implementation-defined signed bit behaviour and signed-overflow UB in edge cases.

---

## Q31. [M2 | Signed integer overflow | Easy | MCQ | [C] · C11 | All] · Type T3 · P09

> **This code has undefined behaviour.**

```c
#include <stdio.h>
#include <limits.h>
int main(void) {
    int x = INT_MAX;
    int y = x + 1;
    printf("%d\n", y);
    return 0;
}
```

Which statement is correct?

- A. It prints `-2147483648`, and this is guaranteed because signed integers wrap around.
- B. It prints `2147483648`.
- C. Compile error: integer overflow.
- D. The behaviour is undefined (signed integer overflow), so no output is guaranteed.

**Answer:** D

**Behaviour class:** **Undefined behaviour.**

> **⚠ UB box**
> - **Category:** signed integer overflow.
> - **Rule:** N1570 §6.5p5: if the result of an expression is not mathematically representable in its type, the behaviour is undefined. (C++: [expr]p4.)
> - **On GCC/Clang x86-64** this typically prints `-2147483648` (verified with `gcc -O0`), because the hardware `add` instruction wraps. **That is not guaranteed and must not be relied on.** At `-O2`, code that depends on the wrap can be optimised away entirely (see P96 in M25).

**Why:**

| Line | Statement executed | Variables after execution |
|---|---|---|
| 4 | `int x = INT_MAX;` | x = 2147483647 |
| 5 | `x + 1` → the mathematical result 2147483648 does not fit in `int` | **UB**: y has no defined value |

**Why the other options fail:**
- A: The observed output, presented as a guarantee. "Wraps around" is true for **unsigned** types only.
- B: `y` is an `int`, so it cannot hold 2³¹.
- C: `x` is a variable, so this is not a constant expression and there is no compile-time diagnostic. (For the *constant* expression `INT_MAX + 1`, GCC warns but still compiles.)

**Concept recap:** Ways to get defined behaviour:
- compute in a wider type: `(long long)x + 1`;
- use unsigned arithmetic;
- check before the operation: `if (x > INT_MAX - 1)`;
- use `__builtin_add_overflow` (GCC/Clang) or C23 `ckd_add`.

**Trap alert:** "But it printed −2147483648 when I ran it" — that is exactly what UB looks like on this hardware. It is still not a guarantee.

**Interview angle:** *"Why is signed overflow UB when every modern CPU wraps?"* "It lets the compiler assume `i + 1 > i`, simplify loop bounds, and widen induction variables to 64 bits. The cost is that overflow bugs become optimiser-visible UB. Use `-fsanitize=undefined` or `-ftrapv` to catch them; `-fwrapv` makes wrapping defined, as a GCC extension."

---

## Q32. [M2 | Spotting overflow UB | Hard | MSQ | [C] · C11 | All] · Type T3 · P09

Assume a 32-bit `int` and that `limits.h` is included. Select **all** snippets whose behaviour is **undefined**.

- A. `int x = INT_MAX; x++;`
- B. `unsigned u = UINT_MAX; u++;`
- C. `int x = INT_MIN; int y = -x;`
- D. `int x = INT_MIN, m = -1; int q = x / m;`
- E. `long long z = (long long)INT_MAX + 1;`

**Answer:** A, C, D

**Behaviour class per option:** A = UB · B = well-defined · C = UB · D = UB · E = well-defined

**Why each option:**
- **A ✓ UB:** `x++` computes INT_MAX + 1 in `int`. That is signed overflow (§6.5p5).
- **B ✗:** Unsigned arithmetic wraps modulo 2³². `u` becomes 0, guaranteed.
- **C ✓ UB:** In two's complement, |INT_MIN| = 2147483648 > INT_MAX, so **negating INT_MIN overflows**. (`abs(INT_MIN)` is UB for the same reason.)
- **D ✓ UB:** The true quotient 2147483648 is not representable. N1570 §6.5.5p6 makes both `INT_MIN / -1` and `INT_MIN % -1` undefined. On x86-64 the `idiv` instruction raises a hardware exception, so the process **typically dies with SIGFPE**. A real crash, not a wrap.
- **E ✗:** The cast happens first, so the addition is done in 64-bit `long long`, where 2147483648 fits.

**Concept recap:** Signed overflow can come from `+`, `-`, `*`, unary `-`, `/` (only INT_MIN / −1), `%` (only INT_MIN % −1), `++`, `--`, and compound assignments.

**Trap alert:** Option D looks harmless (dividing makes numbers *smaller*), but it is the one that actually crashes on x86.

**Interview angle:** *"Write `abs` safely for any int."* "You can't return `abs(INT_MIN)` as an `int`. Return `unsigned`: `x < 0 ? 0u - (unsigned)x : (unsigned)x`, or document a precondition `x != INT_MIN`."

---

## Q33. [M2 | Signed vs unsigned comparison & wrap | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P11

```c
#include <stdio.h>
int main(void) {
    int a = -1;
    unsigned int b = 1;
    if (a < b)
        printf("less\n");
    else
        printf("not less\n");
    printf("%d\n", a + b > 0);
    return 0;
}
```

What does this program print?

- A. `less` then `0`
- B. `not less` then `1`
- C. `not less` then `0`
- D. `less` then `1`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:**

| Line | Expression | Conversion | Value |
|---|---|---|---|
| 5 | `a < b` | `int` vs `unsigned int` → `a` becomes `unsigned`: 4294967295 | `4294967295 < 1` → **0** → prints `not less` |
| 9 | `a + b` | same conversion: 4294967295u + 1u → wraps modulo 2³² (well-defined) | **0u** |
| 9 | `0u > 0` | `0` becomes `0u` → `0u > 0u` | **0** → prints `0` |

**Why the other options fail:**
- A, D: Compare −1 and 1 mathematically.
- B: Gets the comparison right but forgets that `a + b` is computed in `unsigned` and wraps to exactly 0.

**Concept recap:** Once an `unsigned int` operand appears next to an `int` of the same rank, the whole operation becomes unsigned. Negative values turn into huge positive ones.

**Trap alert:** The well-known `for (unsigned i = n - 1; i >= 0; i--)` loop is infinite, because `i >= 0` is always true for unsigned `i`. (This appears again in M3.)

**Interview angle:** *"GCC warned 'comparison of integer expressions of different signedness'. Should I silence it with a cast?"* "Only after checking that the signed value can't be negative. Otherwise the cast just hides the bug. Fix the types instead, e.g. use `size_t` for both."

---

## Q34. [M2 | The sizeof comparison trap | Medium | MCQ | [C] · C11 | All] · Type T1 · P11

```c
#include <stdio.h>
int main(void) {
    int arr[] = {1, 2, 3};
    int i = -1;
    if (i < sizeof(arr) / sizeof(arr[0]))
        printf("inside\n");
    else
        printf("outside\n");
    return 0;
}
```

What does this program print?

- A. `inside`
- B. `outside`
- C. Compile error: comparison between different types
- D. The behaviour is undefined

**Answer:** B

**Behaviour class:** Well-defined.

**Why:**

| Step | Expression | Type | Value |
|---|---|---|---|
| 1 | `sizeof(arr) / sizeof(arr[0])` | `size_t` = `unsigned long` (64-bit) | 12 / 4 = 3 |
| 2 | `i < 3UL` | `int` (rank lower) vs `unsigned long` → `i` becomes `unsigned long` | −1 → 18446744073709551615 |
| 3 | `18446744073709551615 < 3` | — | **false** → prints `outside` |

**Why the other options fail:**
- A: Mathematical comparison. `sizeof` yields an unsigned value, which drags `i` into unsigned.
- C: Mixed-sign comparisons are allowed (GCC only warns with `-Wsign-compare`).
- D: Converting a negative value to unsigned is defined (modulo 2⁶⁴).

**Concept recap:** `sizeof`, `strlen`, and C++ `.size()` all return **unsigned** types (`size_t`). Any comparison or subtraction with them is unsigned arithmetic.

**Trap alert:** `for (int j = 0; j < strlen(s) - 1; j++)` with an empty `s`: `0 - 1` wraps to `SIZE_MAX`, so the loop runs far past the end of the buffer.

**Interview angle:** *"Why did C++20 add `std::ssize()`?"* "To get a **signed** size, which avoids exactly this class of bug when comparing with signed loop variables or subtracting."

---

## Q35. [M2 | Unsequenced modification | Medium | MCQ | [C] · C11 | All] · Type T3 · P49

> **This code has undefined behaviour.**

```c
#include <stdio.h>
int main(void) {
    int i = 5;
    i = i++ + ++i;
    printf("%d\n", i);
    return 0;
}
```

Which is correct?

- A. It prints `11`
- B. It prints `12`
- C. It prints `13`
- D. The behaviour is undefined (unsequenced modifications of `i`), so no output is guaranteed

**Answer:** D

**Behaviour class:** **Undefined behaviour.**

> **⚠ UB box**
> - **Category:** sequence-point violation (unsequenced side effects on the same object).
> - **Rule:** N1570 §6.5p2: if a side effect on a scalar object is unsequenced relative to another side effect on the same object, or to a value computation using its value, the behaviour is undefined. Here `i` is modified **three times** (by `i++`, by `++i`, and by `=`) with no sequencing between them. The same applies in C++ ([intro.execution]).
> - **On GCC and Clang x86-64, both print `12`** for this exact code (verified). Clang also warns *"multiple unsequenced modifications to 'i'"*. **That is not guaranteed and must not be relied on.** Old Turbo C answer keys say 13, which shows the same point.

**Why:** The operands of `+` can be evaluated in any order, and each `++` can write back its result at any time before the next sequence point (the `;`). Nothing in the language picks one ordering, so the standard gives no meaning to the expression at all.

**Why the other options fail:**
- A, B, C: Each is one plausible hardware ordering, presented as *the* answer. GATE and many OA papers have used this exact snippet, and the only correct answer is "undefined".

**Concept recap:** **Sequence points** in C, the points where all earlier side effects must be complete:
- the end of a full expression (`;`), and the controlling expressions of `if`/`while`/`for`;
- after the first operand of `&&`, `||`, `?:` and the comma operator;
- before a function call (after its arguments are evaluated).

**Trap alert:** "Undefined" doesn't mean "an unknown one of 11/12/13". It means the program as a whole has **no** meaning. The compiler may produce any output, or reorder unrelated code.

**Interview angle:** *"What is `i = i++ + ++i`?"* "Undefined behaviour: three unsequenced writes to `i`. The C++17 sequencing rules made some patterns well-defined, like `i = i++ + 1`, but this one is still UB." Follow-up: *"How do you detect it?"* Clang's `-Wunsequenced`, GCC's `-Wsequence-point` (in `-Wall`), and UBSan.

---

## Q36. [M2 | Which expressions are unsequenced? | Hard | MSQ | [C] · C11 | All] · Type T3 · P49

Assume `int i = 1, x, y = 2, a[10];` and a function `void f(int, int);`. Select **all** statements whose behaviour is **undefined in C11**.

- A. `i = i + 1;`
- B. `a[i] = i++;`
- C. `x = y++ + y;`
- D. `x = i++ && i++;`
- E. `f(i++, i++);`

**Answer:** B, C, E

**Why each option:**
- **A ✗:** One modification of `i`; the read is used only to compute the value being stored. Well-defined.
- **B ✓ UB:** `i` is read on the left (to compute the address `a[i]`) and modified on the right (`i++`), and C11 does not sequence the two operands of `=`.
  - **In C++17 this is well-defined**: the right operand of `=` is sequenced before the left, so it stores 1 into `a[2]`.
- **C ✓ UB:** `y` is modified by `y++` and read by the other `y`, unsequenced (operands of `+` are unordered).
- **D ✗:** `&&` has a sequence point after its left operand. Left `i++` → 1 (i = 2), true, so the right `i++` → 2 (i = 3). x = 1. Well-defined.
- **E ✓ UB (in C):** Function **arguments** are unsequenced relative to each other in C (N1570 §6.5.2.2p10 only places a sequence point *after all* arguments). Two modifications of `i` → UB.
  - **In C++17** the arguments are *indeterminately sequenced*: each one is fully evaluated before the other starts, in an unspecified order. So it is **unspecified** (f receives (1, 2) or (2, 1)), not UB.

**Concept recap:** The general rule: in one expression, never modify a variable more than once, and never read it elsewhere in the same expression if you also modify it — unless an operator with a sequence point (`&&`, `||`, `?:`, `,`) separates the two.

**Trap alert:** The same line can be **UB in C, well-defined in C++17** (option B) or **UB in C, unspecified in C++17** (option E). Always check the standard tag on the question.

**Interview angle:** *"What did C++17 change about evaluation order?"* "It added sequencing for `=` (right before left), `<<`/`>>` (left before right, so `cout << i++ << i++` is now defined), `.`/`->`, subscripting, and made function arguments indeterminately sequenced instead of unsequenced."

---

## Q37. [M2 | Floating → integer conversion | Hard | MCQ | [C] · C11 | GATE+OA] · Type T3 · P50

```c
#include <stdio.h>
int main(void) {
    int a = 3.99;
    int b = -3.99;
    double big = 2.5e10;
    int c = (int)big;
    printf("%d %d %d\n", a, b, c);
    return 0;
}
```

Which statement is correct?

- A. `a = 3` and `b = -4` (rounded down); `c` is implementation-defined.
- B. `a = 3` and `b = -3` (truncated toward zero). Converting `2.5e10` to `int` is **undefined behaviour**, so the program as a whole has no guaranteed output.
- C. `a = 4` and `b = -4` (rounded to nearest); `c` is implementation-defined.
- D. `a = 3` and `b = -3`; `c` is implementation-defined (it is `INT_MIN` on x86).

**Answer:** B

**Behaviour class:** **Undefined behaviour** (because of `c`). The conversions for `a` and `b` are well-defined.

> **⚠ UB box**
> - **Category:** out-of-range floating → integer conversion.
> - **Rule:** N1570 §6.3.1.4p1: a floating value is converted to an integer type by discarding its fractional part (truncation toward zero). If the integral part cannot be represented in the target type, the behaviour is undefined. (C++: [conv.fpint]p1, same rule.)
> - **On GCC x86-64** this typically prints `3 -3 -2147483648` (verified at `-O0` and `-O2`), because the `cvttsd2si` instruction returns the "integer indefinite" value 0x80000000. When the compiler folds a *constant* at compile time, it may produce a different value. **That is not guaranteed and must not be relied on.**

**Why:**

| Line | Statement executed | Rule | Variables after execution |
|---|---|---|---|
| 3 | `int a = 3.99;` | truncate toward zero | a = 3 |
| 4 | `int b = -3.99;` | truncate toward zero (**not** floor) | b = −3 |
| 5 | `double big = 2.5e10;` | — | big = 25000000000.0 |
| 6 | `(int)big` | 2.5×10¹⁰ > INT_MAX ≈ 2.1×10⁹ → **UB** | c = ? |

**Why the other options fail:**
- A: Floor gives −4. C truncates toward zero.
- C: The conversion never rounds. Use `lround`/`round` if you want rounding.
- D: This mixes up two rules. **Integer → narrower signed integer** out of range is implementation-defined (Q11). **Floating → integer** out of range is **UB**.

**Concept recap:**

| Conversion | Out-of-range result |
|---|---|
| any integer → unsigned integer | wraps modulo 2ⁿ (defined) |
| integer → signed integer | implementation-defined |
| floating → any integer | **UB** |
| `double` → `float` (value outside float's range) | UB; within range, rounded |

**Trap alert:** Converting a negative `double` to an **unsigned** integer is also UB, e.g. `(unsigned)-1.0`. The "modulo" rule applies only to integer sources.

**Interview angle:** *"How do you safely convert a double to int?"* "Check the range first: `if (d >= INT_MIN && d < INT_MAX + 1.0)`, and handle NaN (every comparison with NaN is false). Then cast, or use `lrint` and check `errno`/the floating-point exceptions."

---

## Q38. [M2 | Integer division & modulo with negatives | Medium | MCQ | [C] · C11 | GATE] · Type T1 · —

```c
#include <stdio.h>
int main(void) {
    printf("%d %d %d %d\n", 7 / -2, -7 / 2, 7 % -2, -7 % 2);
    return 0;
}
```

What does this program print?

- A. `-3 -3 1 -1`
- B. `-4 -4 1 -1`
- C. `-3 -3 -1 1`
- D. `-4 -3 1 1`

**Answer:** A

**Behaviour class:** Well-defined (since C99; in C89 the rounding direction was implementation-defined).

**Why:** The two rules (N1570 §6.5.5p6):
1. `a / b` truncates toward zero.
2. `(a / b) * b + a % b == a`. As a consequence, **`a % b` has the sign of `a` (the dividend)**.

| Expression | Quotient (toward zero) | Remainder = a − q·b | Value |
|---|---|---|---|
| `7 / -2` | −3.5 → **−3** | — | −3 |
| `-7 / 2` | −3.5 → **−3** | — | −3 |
| `7 % -2` | q = −3 | 7 − (−3)(−2) = 7 − 6 | **1** |
| `-7 % 2` | q = −3 | −7 − (−3)(2) = −7 + 6 | **−1** |

Prints `-3 -3 1 -1`.

**Why the other options fail:**
- B: Uses floor division (Python's `//`). C truncates toward zero.
- C: Gives the remainder the divisor's sign (Python's `%`).
- D: Mixes both rules.

**Concept recap:** `x % 2 == 1` is **false** for negative odd `x` (−7 % 2 == −1). Test oddness with `x % 2 != 0` or `x & 1`. Division or modulo by zero is UB.

**Trap alert:** A circular-buffer index computed as `(i - 1) % n` goes negative when i = 0. Use `(i - 1 + n) % n`.

**Interview angle:** *"Why do C and Python disagree on `-7 % 2`?"* "C truncates toward zero, matching CPU division instructions, so the remainder takes the dividend's sign. Python floors, so the remainder takes the divisor's sign. Neither is wrong; they're different definitions."

---

## Q39. [M2 | Shift-count & left-shift rules | Hard | MSQ | [C] · C11 | All] · Type T3 · —

Assume a 32-bit `int` and `unsigned int`, and `int n = -1;`. Select **all** expressions whose behaviour is **undefined** in C11.

- A. `1u << 31`
- B. `1 << 31`
- C. `1u << 32`
- D. `-1 >> 1`
- E. `8 << n`

**Answer:** B, C, E

**Behaviour class per option:** A = well-defined · B = UB · C = UB · D = implementation-defined · E = UB

**Why each option** (N1570 §6.5.7):
- **A ✗:** An unsigned left shift is defined as E1 × 2^E2 mod 2³² = 2147483648u. Well-defined.
- **B ✓ UB:** For a **signed** left operand, E1 × 2^E2 must be representable in the result type (p4). 2³¹ > INT_MAX → UB in C. (C++14–17 made this particular case defined, and C++20 defines all left shifts as modular.) Use `1u << 31` or `INT_MIN`.
- **C ✓ UB:** A shift count **≥ the width** of the (promoted) left operand is UB (p3), even for unsigned. On x86 the hardware masks the count to 5 bits, so `1u << 32` often yields **1**, not 0. Never rely on it.
- **D ✗:** Right-shifting a negative value is **implementation-defined** (p5), not UB. GCC gives −1.
- **E ✓ UB:** A **negative** shift count is UB (p3).

**Concept recap:** Shift rules:
- the count must be in [0, width − 1];
- left-shifting a signed value is defined only if the result fits;
- right-shifting a negative value is implementation-defined;
- unsigned operands are always the safe choice.

**Trap alert:** Option C bites people who build masks: `(1u << nbits) - 1` is UB when `nbits == 32`. Handle that case separately.

**Interview angle:** *"Write a function that sets bit k of a 32-bit word."* `w | (1u << k)` with a precondition `k < 32`. Mention that the literal must be `1u`, not `1`, so that k = 31 works.

---

## Q40. [M2 | Logical vs bitwise operators | Easy | MSQ | [C] · C11 | GATE] · Type T1 · —

Given `int a = 2, b = 4;`, select **all** expressions that evaluate to **1**.

- A. `a && b`
- B. `a & b`
- C. `!!a`
- D. `(a & b) == 0`
- E. `a | b == 6`

**Answer:** A, C, D

**Why each option** (verified by compiling):

| Option | Evaluation | Value |
|---|---|---|
| A | `&&`: both operands are non-zero → logical true | **1** ✓ |
| B | `0010 & 0100` = `0000` | 0 |
| C | `!a` = 0 → `!0` = 1 (the `!!` idiom converts any non-zero value to exactly 1) | **1** ✓ |
| D | `a & b` = 0 → `0 == 0` | **1** ✓ |
| E | `==` binds tighter: `a | (b == 6)` = `2 | 0` | 2 |

**Concept recap:**

| Logical (`&& || !`) | Bitwise (`& | ^ ~`) |
|---|---|
| treats each operand as true/false | works on every bit |
| result is always 0 or 1 | result can be any value |
| short-circuits | evaluates both operands |

**Trap alert:** `if (a & b)` and `if (a && b)` disagree here: 2 & 4 = 0 (false), but 2 && 4 = 1 (true). Using `&` instead of `&&` in a condition is a classic silent bug.

**Interview angle:** *"What does `!!x` do and why use it?"* "It normalises any non-zero value to exactly 1. It's useful when storing a boolean into a 1-bit field or comparing two truth values (`!!a == !!b`)." Also `!` vs `~`: `!2` = 0 but `~2` = −3.

---

## Batch 2 summary

| Q | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|
| 21 | E | MCQ | B | Well-defined |
| 22 | H | MCQ | D | Well-defined (platform-dependent) |
| 23 | M | MSQ | A, C, E | — (expression types) |
| 24 | M | MCQ | A | Well-defined |
| 25 | H | MCQ | C | Well-defined |
| 26 | E | MCQ | D | Well-defined |
| 27 | M | MCQ | B | Well-defined |
| 28 | M | MCQ | A | Well-defined |
| 29 | H | MCQ | C | Ill-formed |
| 30 | H | MCQ | A | Implementation-defined |
| 31 | E | MCQ | D | **UB** |
| 32 | H | MSQ | A, C, D | UB / well-defined mix |
| 33 | M | MCQ | C | Well-defined |
| 34 | M | MCQ | B | Well-defined |
| 35 | M | MCQ | D | **UB** |
| 36 | H | MSQ | B, C, E | UB / well-defined mix |
| 37 | H | MCQ | B | **UB** |
| 38 | M | MCQ | A | Well-defined |
| 39 | H | MSQ | B, C, E | UB / impl-defined mix |
| 40 | E | MSQ | A, C, D | Well-defined |

- **Difficulty:** 4 E / 8 M / 8 H. With Q19–Q20, M2 totals 5/9/8, as planned.
- **Format:** 15 MCQ / 5 MSQ. With Q19–Q20, M2 totals 17/5, as planned.
- **MCQ answer keys:** A×4, B×4, C×3, D×4. Running total Q1–Q40: A 8, B 8, C 7, D 8.
- **Every UB question** (Q31, Q32, Q35, Q36, Q37, Q39) has "undefined" as its correct option and never a specific value. The typical GCC/Clang output is given only in the UB box.

---

**Completed: 40 / 400 Remaining: 360**

**Ledger, newly ✓ this batch:**
- P02 ×3 (Q21, Q22, Q23) ★ complete (5 total, including secondary uses)
- P04 ×1 (Q25)
- P06 ×2 (Q26, Q27) ★ complete
- P07 ×2 (Q28, Q29)
- P08 ×1 (Q30)
- P09 ×2 (Q31, Q32) ★ complete
- P11 ×3 (Q22, Q33, Q34) ★ complete
- P49 ×2 (Q35, Q36) ★ complete
- P50 → 2 (Q37)
- P03 → 3 (Q24, plus secondary in Q25)

**M2 complete** (Q19–Q40, 22 questions).

*Say **"Continue."** for Batch 3 (Q41–Q60: M3 control flow, then the start of M4 functions and recursion).*
