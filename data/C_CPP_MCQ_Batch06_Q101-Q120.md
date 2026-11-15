# C/C++ MCQ/MSQ Bank — Batch 6 (Q101–Q120)

**What this batch covers:**
- **M6 Pointers, part 2 (Q101–Q119):**
  - `const` in function parameters
  - pointer-to-pointer: tracing it, letting a function change the caller's pointer, and using an array of pointers
  - `a[i] ≡ i[a]`, and how an array differs from a pointer
  - arrays of pointers vs pointers to arrays (including string tables)
  - modifying a string literal
  - walking a string with a pointer
  - dispatch tables of function pointers, reading complex declarations, `qsort` comparators
  - generic `void *` code and the rules for `void *`
  - rounding an address up to an alignment boundary
  - strict aliasing
  - pointers that dangle when a loop iteration ends
- **M7 Strings, starting (Q120):** `sizeof` vs `strlen`, and `strcat`.

**Patterns:** P17, P19, P22, P23, P24, P25, P27, P30, P32

> **How questions are written from this batch on:**
> - **No compiler-specific or version-specific questions.** Every answer follows from the C/C++ standard plus the platform assumptions stated below. These questions are about the language, not about one compiler.
> - **UB questions** still say what is *often observed in practice*, but only in general terms ("many compilers…"), never as something tied to a particular toolchain.
> - **Style matches recent GATE CSE and OA trends:**
>   - GATE: tracing well-defined code to an integer or string result (the kind of value a numerical-answer question asks for), "select all that apply" multiple-select questions, functions that take an array and its size, and pointer + array + string combinations. GATE never asks for the output of UB code.
>   - OAs: output prediction, find-the-bug, choosing the smallest fix, and callback/`qsort` usage.

> **Conventions for every question in this batch:**
> - Standard: **C11**.
> - Platform: a typical 64-bit system: `sizeof(int)=4`, `sizeof(long)=8`, pointers are 8 bytes, plain `char` is signed, two's complement, IEEE-754 `float`/`double`, ASCII, little-endian.
> - The needed standard headers are included.

---

## Q101. [M6 | const in a function parameter | Easy | MCQ | [C] · C11 | GATE+OA] · Type T4 · P17

```c
size_t count(____ s, char c) {
    size_t n = 0;
    for (; *s; s++)
        if (*s == c) n++;
    return n;
}
```

Which parameter type both **compiles** and **promises not to modify** the string?

- A. `char *const s`
- B. `const char *s`
- C. `const char *const s`
- D. `char *s`

**Answer:** B

**Behaviour class:** Option B is well-defined (`count("banana", 'a')` returns 3). Options A and C are ill-formed.

**Why:** The function changes **the pointer** (`s++`) but only **reads** the characters (`*s`).

| Option | Can `s` change (`s++`)? | Can `*s` change? | Result |
|---|---|---|---|
| A `char *const s` | **no** → `s++` is a compile error | yes | ✗ |
| B `const char *s` | yes | **no** | ✓ compiles and promises not to modify the text |
| C `const char *const s` | **no** → compile error | no | ✗ |
| D `char *s` | yes | yes | compiles, but makes no promise, and callers can't pass a `const char *` without a cast |

**Why the other options fail:** see the table.

**Concept recap:** `const char *` is the standard type for "read-only string" parameters (`strlen`, `strcmp`, `strchr`). It accepts both `char *` and `const char *` arguments: adding `const` is implicit, removing it is not.

**Trap alert:** A `const` placed after the `*` protects the **pointer variable**, which is a purely local copy. That's rarely what you want for a parameter.

**Interview angle:** *"Why is const correctness important in APIs?"* "It documents intent, lets callers pass string literals and other const data without casts, and the compiler enforces it, so an accidental write becomes a compile error instead of a bug."

---

## Q102. [M6 | Tracing a pointer-to-pointer | Easy | MCQ | [C] · C11 | GATE] · Type T5 · P22

```c
#include <stdio.h>
int main(void) {
    int x = 5, y = 9;
    int *p = &x;
    int **pp = &p;
    **pp = 7;
    *pp = &y;
    **pp += 1;
    printf("%d %d %d\n", x, y, *p);
    return 0;
}
```

What does this program print?

- A. `7 9 7`
- B. `8 9 8`
- C. `7 10 10`
- D. `5 10 10`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:**

**Memory table** (x @100, y @104, p @200, pp @300)

| Line | Statement | x | y | p (@200) | pp (@300) |
|---|---|---|---|---|---|
| 3–5 | init | 5 | 9 | 100 → x | 200 → p |
| 6 | `**pp = 7`: pp → p → x | **7** | 9 | → x | → p |
| 7 | `*pp = &y`: changes **p** itself | 7 | 9 | **104 → y** | → p |
| 8 | `**pp += 1`: pp → p → y | 7 | **10** | → y | → p |
| 9 | print x, y, *p | | | | `7 10 10` |

**Why the other options fail:**
- A: Ignores line 7, which re-points `p`.
- B: Applies the `+= 1` to x.
- D: Skips line 6.

**Concept recap:**
- `pp` = the address of `p`.
- `*pp` = `p` itself (assigning to it **re-points p**).
- `**pp` = whatever `p` currently points to.

**Trap alert:** After `*pp = &y`, every later `**pp` reaches `y`. A pointer-to-pointer follows the **current** value of the inner pointer, not the value it had when `pp` was set.

**Interview angle:** *"When do you need a double pointer?"* "When a function must change the **caller's pointer** (allocate a buffer for it, insert at a list head, advance a cursor), and for arrays of pointers like `char **argv`."

---

## Q103. [M6 | Changing the caller's pointer | Medium | MCQ | [C] · C11 | GATE+OA] · Type T6 · P22

```c
#include <stdio.h>
void advance(int **pp) { (*pp)++; }
void advance_copy(int *p) { p++; }
int main(void) {
    int a[] = {3, 6, 9, 12};
    int *p = a;
    advance(&p);
    advance_copy(p);
    advance(&p);
    printf("%d %td\n", *p, p - a);
    return 0;
}
```

What does this program print?

- A. `9 2`
- B. `12 3`
- C. `6 1`
- D. `3 0`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:** `advance` receives the **address of `p`**, so `(*pp)++` moves main's `p`. `advance_copy` receives a **copy** of `p` and moves only that copy.

| Line | Call | main's p (as an index) |
|---|---|---|
| 6 | `p = a` | 0 |
| 7 | `advance(&p)` | **1** |
| 8 | `advance_copy(p)`: the copy moves to 2 and is then discarded | 1 |
| 9 | `advance(&p)` | **2** |
| 10 | `*p` = a[2] = 9, `p - a` = 2 | `9 2` |

**Why the other options fail:**
- B: Thinks `advance_copy` also moves main's pointer.
- C: Counts only one successful advance.
- D: Thinks no call can move main's pointer.

**Concept recap:** `(*pp)++` needs the parentheses. `*pp++` would parse as `*(pp++)`, moving the *local* `pp` and leaving the caller's pointer unchanged.

**Trap alert:** This is exactly the linked-list "insert at head" bug from Batch 3 (Q52). To change a pointer in the caller, pass a pointer to it.

**Interview angle:** *"Write `push_front` for a singly linked list in C."* `void push(Node **head, int v) { Node *n = malloc(sizeof *n); n->v = v; n->next = *head; *head = n; }`. Explain why it takes `Node **`.

---

## Q104. [M6 | Array of pointers + pointer-to-pointer arithmetic | Hard | MCQ | [C] · C11 | GATE] · Type T5 · P22, P18

```c
#include <stdio.h>
int main(void) {
    int a[] = {10, 20, 30, 40};
    int *b[] = {a + 3, a + 2, a + 1, a};
    int **pp = b;
    pp++;
    printf("%d ", **pp);
    printf("%d ", *(*pp + 1));
    printf("%d ", **++pp);
    printf("%td\n", *pp - a);
    return 0;
}
```

What does this program print?

- A. `30 20 20 2`
- B. `20 30 30 1`
- C. `30 40 10 0`
- D. `30 40 20 1`

**Answer:** D

**Behaviour class:** Well-defined. Every pointer stays inside its array, and each `printf` is a separate statement, so evaluation order is never in question.

**Why:**

**Memory table**

| Variable | Holds |
|---|---|
| a | {10, 20, 30, 40} at indices 0–3 |
| b[0], b[1], b[2], b[3] | a+3, a+2, a+1, a+0 |
| pp | starts at b[0] |

| Line | Expression | pp → | Evaluation | Prints |
|---|---|---|---|---|
| 6 | `pp++` | b[1] | — | |
| 7 | `**pp` | b[1] | `*b[1]` = `*(a+2)` = a[2] | **30** |
| 8 | `*(*pp + 1)` | b[1] | `*(a + 2 + 1)` = a[3] (the `+ 1` moves along **a**, not b) | **40** |
| 9 | `**++pp` | **b[2]** | `++pp` first → b[2]; `*b[2]` = `*(a+1)` = a[1] | **20** |
| 10 | `*pp - a` | b[2] | `(a + 1) - a` | **1** |

**Why the other options fail:**
- A: Reads `*(*pp + 1)` as `*(*(pp + 1))`, which moves along `b` instead of `a` (giving a[1] = 20).
- B: Forgets the initial `pp++`.
- C: Treats `**++pp` as `**pp` followed by an increment, and uses the wrong index.

**Concept recap:** In `*pp + 1`, the `+ 1` acts on an `int *` (it moves through `a`). In `pp + 1`, it acts on an `int **` (it moves through `b`). The **type** of the operand decides which array is being walked.

**Trap alert:** GATE versions of this question put all four expressions in a single `printf`. In C that's fine only if no argument has a side effect. Since `**++pp` modifies `pp` while `**pp` reads it, putting them in one call would be **UB**. That's why they are separate statements here.

**Interview angle:** *"What is `int *b[4]`, and how does it differ from `int (*b)[4]`?"* The first is an array of 4 pointers (used here); the second is a single pointer to an array of 4 ints (Q107).

---

## Q105. [M6 | `a[i] ≡ *(a + i) ≡ i[a]` | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P19

```c
#include <stdio.h>
int main(void) {
    int a[] = {4, 8, 12, 16};
    printf("%d %c %c\n", 3[a] - 1[a], "GATE"[2], 1["CSE"]);
    return 0;
}
```

What does this program print?

- A. Compile error: an integer can't be indexed
- B. `8 T S`
- C. `8 A C`
- D. `12 T S`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:** The subscript operator is *defined* as `E1[E2]` ≡ `*((E1) + (E2))` (N1570 §6.5.2.1p2). Addition is commutative, so `i[a]` ≡ `a[i]`.

| Expression | Rewritten | Value |
|---|---|---|
| `3[a] - 1[a]` | a[3] − a[1] = 16 − 8 | **8** |
| `"GATE"[2]` | a string literal is a `char` array → its element 2 | **'T'** |
| `1["CSE"]` | `"CSE"[1]` | **'S'** |

**Why the other options fail:**
- A: `3[a]` is legal (but never write it in real code).
- C: Uses 1-based indexing for the strings.
- D: Computes a[3] − a[0].

**Concept recap:** Indexing a string literal directly is a neat trick, e.g. `"0123456789ABCDEF"[n % 16]` converts a number to a hex digit.

**Trap alert:** You can **read** `"GATE"[2]`, but writing to it, as in `"GATE"[2] = 'X'`, modifies a string literal and is UB (Q109).

**Interview angle:** *"Why does `2[a]` compile?"* "Because `a[i]` is just shorthand for `*(a + i)`, and `a + i == i + a`. It shows that C arrays are really pointer arithmetic under the hood."

---

## Q106. [M6 | Array vs pointer | Medium | MSQ | [C] · C11 | All] · Type T3 · P19

Given `int a[5]; int *p = a;`, select **all** true statements.

- A. `a[3]` and `p[3]` refer to the same object.
- B. `sizeof a == sizeof p`
- C. `p = p + 1;` is valid, but `a = a + 1;` is not.
- D. `&a` and `&p` have the same type.
- E. `*(a + 4)` and `p[4]` refer to the same object.

**Answer:** A, C, E

**Why each option:**
- **A ✓:** `p` points to `a[0]`, so `p[3]` = `*(p + 3)` = `a[3]`.
- **B ✗:** `sizeof a` = 5 × 4 = 20 (the whole array). `sizeof p` = 8 (a pointer).
- **C ✓:** `p` is an ordinary variable. `a` is an array and can't be assigned (Q84).
- **D ✗:** `&a` is `int (*)[5]` (a pointer to an array of 5 ints). `&p` is `int **`.
- **E ✓:** Both are a[4].

**Concept recap:**

| | `int a[5]` | `int *p = a` |
|---|---|---|
| storage | 5 ints | one pointer |
| `sizeof` | 20 | 8 |
| can be assigned? | no | yes |
| `&x` type | `int (*)[5]` | `int **` |
| indexing | same result | same result |

**Trap alert:** `extern int a[];` in one file and `int *a;` defining it in another is a classic mismatch. The two declarations disagree, so the program is invalid (UB), even though both "look like arrays".

**Interview angle:** *"Are arrays and pointers the same thing?"* "No. They're different types with different storage. An array *converts* to a pointer in most expressions, which is why indexing works on both."

---

## Q107. [M6 | Array of pointers vs pointer to array | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
int main(void) {
    int *ap[3];
    int (*pa)[3];
    int m[2][3] = {{1, 2, 3}, {4, 5, 6}};
    pa = m;
    printf("%zu %zu %d %d\n", sizeof ap, sizeof pa, (*pa)[2], pa[1][0]);
    return 0;
}
```

What does this program print?

- A. `24 8 3 4`
- B. `8 24 3 4`
- C. `24 24 3 4`
- D. `24 8 2 4`

**Answer:** A

**Behaviour class:** Well-defined. The elements of `ap` are never read: `sizeof` doesn't evaluate them.

**Why:** Read declarations from the name outward. The parentheses change the meaning.

| Declaration | Read as | sizeof |
|---|---|---|
| `int *ap[3]` | ap is an **array of 3** pointers to int (`[]` binds tighter than `*`) | 3 × 8 = **24** |
| `int (*pa)[3]` | pa is a **pointer** to an array of 3 ints | **8** |

| Expression | Meaning | Value |
|---|---|---|
| `pa = m` | `m` decays to `int (*)[3]`, so the types match exactly | pa → row 0 |
| `(*pa)[2]` | row 0, element 2 | **3** |
| `pa[1][0]` | `*(pa + 1)` → row 1, element 0 | **4** |

**Why the other options fail:**
- B: Swaps the two declarations.
- C: Treats `pa` as an array.
- D: Off by one in `(*pa)[2]`.

**Concept recap:** `int (*)[N]` is the type to use when passing a 2D array with N columns to a function (Q79).

**Trap alert:** `*pa[2]` without the parentheses means `*(pa[2])`, which is row 2's first element. Here row 2 doesn't exist, so that would be out of bounds.

**Interview angle:** *"What's the difference between `int *a[10]` and `int (*a)[10]`?"* "The first is ten pointers; the second is one pointer to a block of ten ints. Read inside-out, with `[]` binding tighter than `*`."

---

## Q108. [M6 | String table: `char *[]` vs `char [][N]` | Hard | MCQ | [C] · C11 | GATE+OA] · Type T1 · P20

```c
#include <stdio.h>
int main(void) {
    char *ptrs[] = {"Ada", "Linus", "Ken"};
    char arr[][6] = {"Ada", "Linus", "Ken"};
    printf("%zu %zu %c %c\n", sizeof ptrs, sizeof arr,
           *(ptrs[1] + 2), arr[2][1]);
    return 0;
}
```

What does this program print?

- A. `24 24 n e`
- B. `12 18 n e`
- C. `24 18 n e`
- D. `24 18 i K`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:**

**Memory table**

| Object | Layout | sizeof |
|---|---|---|
| `ptrs` | 3 pointers, each to a read-only string literal stored elsewhere | 3 × 8 = **24** |
| `arr` | 3 rows × 6 chars, **copied** into the array: `Ada\0\0\0`, `Linus\0`, `Ken\0\0\0` | 3 × 6 = **18** |

| Expression | Evaluation | Value |
|---|---|---|
| `*(ptrs[1] + 2)` | "Linus", index 2 | **'n'** |
| `arr[2][1]` | "Ken", index 1 | **'e'** |

**Why the other options fail:**
- A: Treats `arr` as storing pointers.
- B: Uses 4-byte pointers (a 32-bit system).
- D: Off by one in both indexes.

**Concept recap:**
- `char *[]` is a *jagged* table: each string has its own length, but the strings are read-only literals.
- `char [][N]` is a fixed-width table: each row is a **writable** copy, padded to N (the longest string + 1).

**Trap alert:** `arr[0][0] = 'a';` is fine. `ptrs[0][0] = 'a';` modifies a string literal, which is UB (Q109). The safe declaration for the pointer table is `const char *ptrs[]`.

**Interview angle:** *"`main(int argc, char *argv[])` — which kind of table is argv?"* "An array of pointers to separate, null-terminated, **writable** strings, plus a `NULL` sentinel at `argv[argc]`."

---

## Q109. [M6 | Modifying a string literal | Medium | MCQ | [C] · C11 | All] · Type T3 · P25

```c
#include <stdio.h>
int main(void) {
    char a[] = "gate";
    char *p = "gate";
    a[0] = 'G';      /* write 1 */
    p[0] = 'G';      /* write 2 */
    printf("%s %s\n", a, p);
    return 0;
}
```

Which statement is correct?

- A. Both writes are fine; it prints `Gate Gate`.
- B. Write 1 is well-defined. Write 2 is **undefined behaviour**, because it modifies a string literal.
- C. Both writes are undefined behaviour.
- D. Compile error at write 2: `p` points to a constant.

**Answer:** B

**Behaviour class:** **Undefined behaviour** (write 2).

> **⚠ UB box**
> - **Category:** modifying a string literal.
> - **Rule:** N1570 §6.4.5p7: if the program attempts to modify a string literal's array, the behaviour is undefined.
> - **What you may see in practice:** string literals usually live in read-only memory, so the program typically **crashes** (segmentation fault). Some older or embedded systems silently change the text. Implementations may also merge identical literals, so changing one "gate" could change every "gate" in the program. **Not guaranteed either way.**

**Why:**
- `char a[] = "gate";` creates a **writable array** initialised with a copy of the literal.
- `char *p = "gate";` points **into the literal itself**.

In C the literal's type is `char[5]` (not `const char[5]`, for historical reasons), so line 5 compiles without a diagnostic. It is still UB.

**Why the other options fail:**
- A: Treats `p` like `a`.
- C: Write 1 changes the program's own array.
- D: In C, a literal isn't `const`-typed, so no error is required. (In **C++** the literal is `const char[5]`, and `char *p = "gate";` itself is ill-formed since C++11.)

**Concept recap:** Always write `const char *p = "…";`. Then an accidental write becomes a compile error.

**Trap alert:** `char *s = "hello"; strcpy(s, "world");` and `strtok(s, " ")` on a literal are the same bug, because both functions write into the string.

**Interview angle:** *"Why does `char *s = "abc"; s[0] = 'x';` crash?"* "String literals are stored in read-only memory and the standard makes modifying them UB. Use `char s[] = "abc";` for a writable copy." Follow-up: *"In C++?"* `char *s = "abc";` doesn't even compile in C++11 and later.

---

## Q110. [M6 | Walking a string with a pointer | Easy | MCQ | [C] · C11 | GATE] · Type T1 · P25

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char s[] = "GATECSE";
    char *p = s;
    p += 3;
    printf("%s %c %c %zu\n", p, *(p - 2), p[2], strlen(p));
    return 0;
}
```

What does this program print?

- A. `TECSE A C 5`
- B. `ECSE T S 4`
- C. `ECSE A C 4`
- D. `ECSE A S 4`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:**

| Index | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
|---|---|---|---|---|---|---|---|---|
| s | G | A | T | **E** | C | S | E | \0 |

`p = s + 3` → points at index 3.

| Expression | Index | Value |
|---|---|---|
| `p` as `%s` | 3 onward | **"ECSE"** |
| `*(p - 2)` | 1 | **'A'** |
| `p[2]` | 5 | **'S'** |
| `strlen(p)` | characters from index 3 to the `'\0'` | **4** |

**Why the other options fail:**
- A: `p += 3` skips three characters, landing on index 3. This option starts at index 2.
- B: `p - 2` is index 1, not 2.
- C: `p[2]` is index 5, not 4.

**Concept recap:** A `char *` into the middle of a string *is* a string: the suffix starting at that position. That's how `strchr` and `strstr` return results without copying anything.

**Trap alert:** `printf("%s", p - 4)` would read before the start of `s`, which is UB. Moving a pointer backwards is safe only while it stays within the array.

**Interview angle:** *"Return the extension of a filename without allocating."* `const char *dot = strrchr(name, '.'); return dot ? dot + 1 : "";`. This returns a pointer into the original string, using exactly the idea in this question.

---

## Q111. [M6 | Dispatch table of function pointers | Medium | MCQ | [C] · C11 | GATE+OA] · Type T6 · P23

```c
#include <stdio.h>
int add(int a, int b) { return a + b; }
int sub(int a, int b) { return a - b; }
int mul(int a, int b) { return a * b; }
int main(void) {
    int (*ops[])(int, int) = {add, sub, mul};
    int acc = 2;
    for (int i = 0; i < 3; i++)
        acc = ops[i](acc, i + 3);
    printf("%d\n", acc);
    return 0;
}
```

What value is printed?

- A. `7`
- B. `9`
- C. `5`
- D. `1`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:**

| i | Function | Call | acc after |
|---|---|---|---|
| — | — | — | 2 |
| 0 | add | add(2, 3) | 5 |
| 1 | sub | sub(5, 4) | 1 |
| 2 | mul | mul(1, 5) | **5** |

**Why the other options fail:**
- A: Runs the table backwards: mul(2, 3) = 6, sub(6, 4) = 2, add(2, 5) = 7.
- B: Adds every second argument to 2 (2 + 3 + 4 = 9), ignoring which function is called.
- D: Stops after the `sub` step.

**Concept recap:** `int (*ops[])(int, int)` is an **array** of pointers to functions that take `(int, int)` and return `int`. A dispatch table like this replaces a long `switch`, as in interpreters, state machines and command handlers.

**Trap alert:** `ops[i](acc, i + 3)` and `(*ops[i])(acc, i + 3)` are the same call. Writing `*ops[i](…)` without the parentheses would instead call `ops[i]` and then dereference the returned `int`, which is a compile error.

**Interview angle:** *"How do you implement a calculator without a switch?"* "A table of function pointers indexed by an operator code. In C++, a `std::unordered_map<char, std::function<int(int,int)>>`, or an array of lambdas."

---

## Q112. [M6 | Reading complex declarations | Hard | MSQ | [C] · C11 | OA+Interview] · Type T3 · P23

Select **all** correct descriptions.

- A. `int *f(int);` — f is a function that takes an `int` and returns a pointer to `int`.
- B. `int (*f)(int);` — f is a pointer to a function that takes an `int` and returns `int`.
- C. `int (*f[4])(int);` — f is a pointer to an array of 4 functions.
- D. `int (*(*f)(int))[3];` — f is a pointer to a function that takes an `int` and returns a pointer to an array of 3 `int`s.
- E. `void (*signal(int, void (*)(int)))(int);` — signal is a pointer to a function.

**Answer:** A, B, D

**Why each option:** Use the **right-left rule**: start at the name, read right while you can (`()` = "function taking…", `[]` = "array of…"), then read left (`*` = "pointer to"), and jump out of each pair of parentheses.

- **A ✓:** f → `(int)` function taking int → `*` returning a pointer → `int`.
- **B ✓:** `(*f)` → f is a pointer → `(int)` to a function taking int → returning `int`.
- **C ✗:** f → `[4]` an **array of 4** → `*` pointers → `(int)` to functions taking int → returning `int`. It is an array of 4 function pointers. (An "array of functions" isn't even a legal type.)
- **D ✓:** `(*f)` pointer → `(int)` to a function taking int → `*` returning a pointer → `[3]` to an array of 3 → `int`.
- **E ✗:** `signal(int, void (*)(int))` comes first: signal is a **function** taking an `int` and a handler pointer, and it **returns** a pointer to a function `(int)` returning `void`. This is the real declaration of `signal` in `<signal.h>`.

**Concept recap:** A `typedef` makes these readable:
```c
typedef void (*handler_t)(int);
handler_t signal(int sig, handler_t h);
```

**Trap alert:** Writing an "array of functions" or a "function returning a function" is ill-formed. C only has arrays of *pointers* to functions, and functions returning *pointers* to functions.

**Interview angle:** *"Explain `void (*signal(int, void (*)(int)))(int);`"* Walk through the right-left rule, then show the `typedef` version. That combination is what interviewers want to hear.

---

## Q113. [M6 | A qsort comparator callback | Medium | MCQ | [C] · C11 | GATE+OA] · Type T8 · P23, P24

```c
#include <stdio.h>
#include <stdlib.h>
int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    return (x % 10) - (y % 10);
}
int main(void) {
    int v[] = {23, 51, 12, 40, 35};
    qsort(v, 5, sizeof v[0], cmp);
    for (int i = 0; i < 5; i++) printf("%d ", v[i]);
    printf("\n");
    return 0;
}
```

What does this program print?

- A. `40 51 12 23 35`
- B. `12 23 35 40 51`
- C. `51 40 35 23 12`
- D. `35 23 12 51 40`

**Answer:** A

**Behaviour class:** Well-defined. All the last digits (3, 1, 2, 0, 5) are different, so the result doesn't depend on `qsort` not being stable.

**Why:** `qsort` passes **pointers to elements** as `const void *`. The comparator converts them back to `const int *`, then returns a negative, zero or positive value meaning "less", "equal" or "greater". Here it compares the **last digits** (x % 10).

| Element | Key (x % 10) |
|---|---|
| 23 | 3 |
| 51 | 1 |
| 12 | 2 |
| 40 | 0 |
| 35 | 5 |

Sorted by key: 0 → 40, 1 → 51, 2 → 12, 3 → 23, 5 → 35.

**Why the other options fail:**
- B: Sorts by the whole value.
- C: Sorts in descending order of value.
- D: Sorts by key in descending order.

**Concept recap:** The `qsort` comparator contract:
- It receives element **addresses**, not element values.
- It must give a consistent total order.
- `qsort` is **not stable**: elements with equal keys may end up in any order.

**Trap alert:** The comparator `return x - y;` can overflow for large values of opposite sign (e.g. `INT_MIN - 1`), which is signed-overflow UB. The safe idiom is `return (x > y) - (x < y);`.

**Interview angle:** *"Sort structs by a field using qsort."* Cast to `const struct S *`, then compare the field with the `(a > b) - (a < b)` idiom. Mention that C++ `std::sort` with a lambda is type-safe, can be inlined, and is usually faster.

---

## Q114. [M6 | Generic swap with void* | Medium | MCQ | [C] · C11 | OA+Interview] · Type T5 · P24

```c
#include <stdio.h>
#include <string.h>
void swap(void *a, void *b, size_t n) {
    unsigned char tmp[16];
    memcpy(tmp, a, n);
    memcpy(a, b, n);
    memcpy(b, tmp, n);
}
int main(void) {
    double d[2] = {1.5, 2.5};
    int i[3] = {1, 2, 3};
    swap(&d[0], &d[1], sizeof(double));
    swap(i, i + 2, sizeof(int));
    printf("%.1f %.1f %d %d %d\n", d[0], d[1], i[0], i[1], i[2]);
    return 0;
}
```

What does this program print?

- A. `1.5 2.5 1 2 3`
- B. `2.5 1.5 2 1 3`
- C. `2.5 1.5 3 2 1`
- D. `2.5 2.5 3 2 3`

**Answer:** C

**Behaviour class:** Well-defined, because `n` ≤ 16. (A caller passing `n > 16` would overflow `tmp`, which is UB.)

**Why:** Any object pointer converts to `void *` implicitly and back again without loss. `memcpy` copies raw bytes, so the same function swaps a `double` or an `int`.

| Call | Swaps | After |
|---|---|---|
| `swap(&d[0], &d[1], 8)` | the 8 bytes of d[0] with those of d[1] | d = {2.5, 1.5} |
| `swap(i, i + 2, 4)` | i[0] with i[2] | i = {3, 2, 1} |

**Why the other options fail:**
- A: Assumes a `void *` function can't change the caller's data.
- B: Swaps the wrong pair of ints (i[0] and i[1]).
- D: Copies in one direction only (it forgets the temporary).

**Concept recap:** `void *` + a size is C's form of generic programming. `qsort`, `bsearch`, `memcpy` and `memset` all use it. The price is no type checking: the size and the types must be right, or the behaviour is undefined.

**Trap alert:** `swap(&x, &y, sizeof(&x))` (the size of a pointer, by mistake) copies 8 bytes for a 4-byte `int`, which is out-of-bounds UB. Always use `sizeof *ptr` or `sizeof(type)` to match what is being copied.

**Interview angle:** *"How would you write generic swap in C++?"* "`template<class T> void swap(T& a, T& b) { T t = std::move(a); a = std::move(b); b = std::move(t); }`. It's type-checked and can use move semantics. Or just call `std::swap`."

---

## Q115. [M6 | Rules for void* | Hard | MSQ | [C] · C11 | All] · Type T2 · P24

Treat each snippet independently. Assume `int x;` is in scope and `<stdlib.h>` is included. Select **all** snippets that are **valid, well-defined C11**.

- A. `void *v = &x; int *p = v;`
- B. `void *v = &x; v++;`
- C. `int *p = malloc(4 * sizeof *p);`
- D. `void *v = &x; int y = *v;`
- E. `void *v = &x; void (*fp)(void) = (void (*)(void))v;`

**Answer:** A, C

**Why each option:**
- **A ✓:** In C, `void *` converts implicitly to and from any **object** pointer type. (C++ requires a cast in the `void *` → `int *` direction.)
- **B ✗:** Arithmetic on `void *` is not defined, because `void` has no size. Some compilers allow it as an extension (treating the size as 1), but it is not standard C.
- **C ✓:** `malloc` returns `void *`, which converts implicitly. In C, *not* casting `malloc`'s result is idiomatic.
- **D ✗:** You can't dereference a `void *` to get a value: `*v` has type `void`. Ill-formed.
- **E ✗:** Converting an **object** pointer to a **function** pointer is not defined by ISO C. POSIX requires it to work for `dlsym`, but under C11 alone it is not well-defined.

**Concept recap:**

| Conversion | C | C++ |
|---|---|---|
| `T *` → `void *` | implicit | implicit |
| `void *` → `T *` | implicit | needs `static_cast` |
| `void *` ↔ function pointer | not defined | conditionally supported |

**Trap alert:** Casting `malloc` (`(int *)malloc(...)`) in C is harmless today. In old C89 code, it could hide a missing `#include <stdlib.h>`: the compiler would assume `malloc` returns `int`, and on 64-bit systems the pointer could be truncated.

**Interview angle:** *"Why does C++ require a cast for `void *` → `T *`?"* "Type safety. An implicit conversion would let any pointer silently become any other through `void *`. In C++ you use `new` or `static_cast` anyway."

---

## Q116. [M6 | Rounding an address up to alignment | Hard | MSQ | [C] · C11 | OA+Interview] · Type T1

A memory allocator rounds an address up to the next multiple of 8 using:

```c
#include <stdint.h>
uintptr_t align8(uintptr_t addr) { return (addr + 7) & ~(uintptr_t)7; }
```

Select **all** inputs for which `align8(addr) == 0x1008`.

- A. `0x1001`
- B. `0x1008`
- C. `0x1009`
- D. `0x1000`
- E. `0x1007`

**Answer:** A, B, E

**Behaviour class:** Well-defined. This is pure unsigned arithmetic on `uintptr_t`.

**Why:** `~(uintptr_t)7` is a mask with the low 3 bits cleared, so `& ~7` rounds **down** to a multiple of 8. Adding 7 first turns that into rounding **up**, while leaving values that are already multiples of 8 unchanged.

| addr | addr + 7 | & ~7 | = 0x1008? |
|---|---|---|---|
| 0x1001 | 0x1008 | **0x1008** | ✓ |
| 0x1008 | 0x100F | **0x1008** | ✓ |
| 0x1009 | 0x1010 | 0x1010 | ✗ |
| 0x1000 | 0x1007 | 0x1000 | ✗ (already aligned) |
| 0x1007 | 0x100E | **0x1008** | ✓ |

**Concept recap:**
- `uintptr_t` (from `<stdint.h>`) is an unsigned integer type that can hold any `void *` and convert back to the same pointer. It is technically *optional*, but available on all mainstream platforms.
- How a pointer maps to an integer is implementation-defined. On ordinary flat-memory systems, it's simply the byte address.
- Align up to a power of two N: `(x + N − 1) & ~(N − 1)`. Check whether x is aligned: `(x & (N − 1)) == 0`.

**Trap alert:** Writing the mask as `~7` (a signed `int`) instead of `~(uintptr_t)7` works only because of sign-extension when it is converted. Making the type explicit avoids surprises with `unsigned int` masks on 64-bit addresses: `~7u` would clear the high 32 bits.

**Interview angle:** *"Why does alignment matter?"* "Misaligned access is slower on x86 and can fault on some architectures. SIMD and atomic operations require alignment, and `malloc` guarantees alignment suitable for any fundamental type (`max_align_t`). C11 adds `alignof`, `_Alignas` and `aligned_alloc`."

---

## Q117. [M6 | Strict aliasing: reading a float as an int | Hard | MCQ | [C] · C11 | All] · Type T3 · P30

> **This code has undefined behaviour.**

```c
#include <stdio.h>
int main(void) {
    float f = 1.0f;
    int *ip = (int *)&f;
    printf("%x\n", (unsigned)*ip);
    return 0;
}
```

Which statement is correct?

- A. It prints `3f800000`, guaranteed, because the cast reinterprets the bits.
- B. The behaviour is undefined: a `float` object is read through an `int` lvalue, which violates the effective-type (strict aliasing) rule. The portable way is `memcpy(&u, &f, sizeof u)`, which yields the IEEE-754 pattern `3f800000`.
- C. It prints `1`.
- D. Compile error: `float *` can't be cast to `int *`.

**Answer:** B

**Behaviour class:** **Undefined behaviour.**

> **⚠ UB box**
> - **Category:** strict-aliasing (effective-type) violation.
> - **Rule:** N1570 §6.5p7: an object may be accessed only through an lvalue of a compatible type (or its signed/unsigned variant, a qualified version, a containing struct/union, or a **character type**). `int` is not compatible with `float`.
> - **What you may see in practice:** often `3f800000` without optimisation. With optimisation, compilers assume an `int *` and a `float *` never refer to the same memory, so they may reorder or drop the read, and the program can print stale or unexpected values. **Not guaranteed.**

**Why:** The cast itself is allowed. The **access** through `*ip` is the problem. The well-defined ways to get the bits:
- `unsigned u; memcpy(&u, &f, sizeof u);` → `0x3f800000` (sign 0, exponent 127, mantissa 0)
- a union in C (pattern P34, M8)
- `std::bit_cast` in C++20

**Why the other options fail:**
- A: The bit pattern is right, but the guarantee is wrong.
- C: That's a *value* conversion (`(int)f`), not a reinterpretation of the bits.
- D: Pointer casts between object types compile. Only the later access is invalid.

**Concept recap:** Strict aliasing lets the compiler assume that pointers to different types don't overlap, which enables optimisations like keeping a value in a register across writes through another pointer type.

**Trap alert:** Network-packet code like `uint32_t *w = (uint32_t *)buf;` (where `buf` is a `char` array) has the same problem in reverse: it accesses a `char` array through a `uint32_t` lvalue, which is UB. Use `memcpy` into a `uint32_t`; compilers turn it into a single load.

**Interview angle:** *"What is strict aliasing?"* "The rule that objects must be accessed through their own type, a signed/unsigned variant, or `char`. Violating it is UB, and optimisers rely on it. Use `memcpy` or `bit_cast` for type punning. `-fno-strict-aliasing` disables the optimisation but doesn't make the code standard."

---

## Q118. [M6 | What strict aliasing allows | Hard | MSQ | [C] · C11 | OA+Interview] · Type T3 · P30

Select **all** accesses that the effective-type (strict aliasing) rule **permits**. Each access is through a pointer to the given type, aimed at an object of the stated type.

- A. Reading an `int` object through an `unsigned int *`.
- B. Reading a `double` object through an `unsigned char *`.
- C. Reading a `float` object through an `int *`.
- D. Reading member `s.a` of `struct S { int a; } s;` through an `int *` that points at `s.a`.
- E. Writing a `long` object through an `int *`.

**Answer:** A, B, D

**Why each option** (N1570 §6.5p7):
- **A ✓:** The signed and unsigned variants of the same type may alias each other.
- **B ✓:** Character types may read the bytes of **any** object (this is how `memcpy` and hex dumps work).
- **C ✗:** `int` and `float` are unrelated types → UB (Q117).
- **D ✓:** The object being accessed is an `int` (the member), and it is read through an `int` lvalue, so the types match exactly.
- **E ✗:** `long` and `int` are **different types**, even on platforms where they have the same size → UB.

**Concept recap:** The `restrict` qualifier (C99) works in the opposite direction. `void f(int *restrict a, int *restrict b)` **promises** the compiler that `a` and `b` don't overlap, so it can vectorise loops. Breaking that promise is UB. `memcpy` declares its parameters `restrict`; `memmove` doesn't, which is why `memmove` must be used for overlapping regions.

**Trap alert:** Option E catches people who think "same size means same type". Aliasing rules are about **types**, not sizes.

**Interview angle:** *"`memcpy` vs `memmove`?"* "`memcpy` requires non-overlapping buffers (its parameters are `restrict`-qualified, so overlap is UB). `memmove` handles overlap correctly, as if it copied through a temporary buffer. Use `memmove` whenever the regions might overlap."

---

## Q119. [M6 | Pointers that dangle when a loop iteration ends | Medium | MCQ | [C] · C11 | All] · Type T5 · P27

```c
#include <stdio.h>
int main(void) {
    int *keep[3];
    for (int i = 0; i < 3; i++) {
        int v = i * 10;
        keep[i] = &v;
    }
    printf("%d\n", *keep[0]);
    return 0;
}
```

Which statement is correct?

- A. It prints `0`, the value stored when `i` was 0.
- B. It prints `20`, guaranteed, because all three pointers refer to the same variable `v`.
- C. Compile error: the address of a block-scope variable can't be stored in an outer array.
- D. The behaviour is undefined: each `v`'s lifetime ends at the end of its loop iteration, so `keep[0]` is a dangling pointer.

**Answer:** D

**Behaviour class:** **Undefined behaviour.**

> **⚠ UB box**
> - **Category:** dangling pointer (use after the end of an automatic object's lifetime).
> - **Rule:** N1570 §6.2.4p2: when a block is exited, the lifetime of its automatic objects ends, and pointers to them become indeterminate. Each loop iteration enters and exits the body block, so every iteration creates a **new** `v`.
> - **What you may see in practice:** often `20`, because compilers usually reuse the same stack slot for `v` on every iteration, so all three pointers hold the same address, and the last value written there was 20. **Not guaranteed.**

**Why:**

| Iteration | `v` (a new object each time) | keep[i] | After the iteration |
|---|---|---|---|
| 0 | v₀ = 0 | &v₀ | v₀'s lifetime ends |
| 1 | v₁ = 10 | &v₁ | v₁'s lifetime ends |
| 2 | v₂ = 20 | &v₂ | v₂'s lifetime ends |
| print | `*keep[0]` → v₀ no longer exists → **UB** | | |

**Why the other options fail:**
- A: v₀ no longer exists, so there is nothing to read.
- B: "Same address" is an implementation detail, and "guaranteed" is false.
- C: Storing the address is legal. Only *using* it after the lifetime ends is UB.

**Concept recap:** To keep values from a loop, store the **values** (`int keep[3]; keep[i] = v;`), or give each one its own storage (`malloc`, or an array declared outside the loop).

**Trap alert:** The same bug appears in C++ when a lambda captures a loop variable **by reference** and runs after the loop ends (P80, M22).

**Interview angle:** *"Why do all my stored pointers show the last value?"* "They all point at the same stack slot, which was reused for each iteration's variable. And reading it after the loop is UB in any case. Store values, or allocate separately."

---

## Q120. [M7 | sizeof vs strlen, and strcat | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P32

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char s[20] = "GATE2026";
    printf("%zu %zu\n", sizeof s, strlen(s));
    strcat(s, "-CS");
    printf("%s %zu\n", s, strlen(s));
    return 0;
}
```

What does this program print?

- A. `20 8` then `GATE2026-CS 11`
- B. `9 8` then `GATE2026-CS 11`
- C. `20 9` then `GATE2026-CS 12`
- D. `20 8` then `GATE2026 8`

**Answer:** A

**Behaviour class:** Well-defined. The result needs 11 characters + `'\0'` = 12 bytes, which fits in 20.

**Why:**

| Line | Statement | s contents | sizeof s | strlen(s) |
|---|---|---|---|---|
| 4 | `char s[20] = "GATE2026"` | `GATE2026\0` followed by zero-filled bytes | **20** (the declared size) | **8** |
| 6 | `strcat(s, "-CS")` | finds the `'\0'` at index 8, copies `-CS\0` from there | 20 | **11** |
| 7 | print | `GATE2026-CS` | | |

**Why the other options fail:**
- B: Sizes the array from the initializer. The explicit `[20]` fixes it at 20.
- C: Counts the `'\0'` in `strlen`.
- D: Assumes `strcat` returns a new string and leaves `s` unchanged. It modifies `s` in place.

**Concept recap:**
- `sizeof` = the size of the array (fixed at compile time).
- `strlen` = the current length of the text (computed at run time, up to the first `'\0'`).
- `strcat(dst, src)` appends in place. It is your responsibility that `dst` has room for `strlen(dst) + strlen(src) + 1` bytes.

**Trap alert:** With `char s[] = "GATE2026";` (sized 9), the same `strcat` would write past the end, which is UB. This is the classic buffer overflow, covered in the rest of M7.

**Interview angle:** *"Why is `strcat` O(n) per call, and why does that matter?"* "It scans `dst` to find the end every time, so appending in a loop is O(n²). Keep a pointer to the end, or use `snprintf` / `memcpy` with tracked lengths. In C++, `std::string::append` is amortised O(1) per character."

---

## Batch 6 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 101 | M6 | E | MCQ | B | Well-defined / ill-formed options |
| 102 | M6 | E | MCQ | C | Well-defined |
| 103 | M6 | M | MCQ | A | Well-defined |
| 104 | M6 | H | MCQ | D | Well-defined |
| 105 | M6 | E | MCQ | B | Well-defined |
| 106 | M6 | M | MSQ | A, C, E | — (concepts) |
| 107 | M6 | M | MCQ | A | Well-defined |
| 108 | M6 | H | MCQ | C | Well-defined |
| 109 | M6 | M | MCQ | B | **UB** (write 2) |
| 110 | M6 | E | MCQ | D | Well-defined |
| 111 | M6 | M | MCQ | C | Well-defined |
| 112 | M6 | H | MSQ | A, B, D | — (declarators) |
| 113 | M6 | M | MCQ | A | Well-defined |
| 114 | M6 | M | MCQ | C | Well-defined |
| 115 | M6 | H | MSQ | A, C | Ill-formed / not standard / OK mix |
| 116 | M6 | H | MSQ | A, B, E | Well-defined |
| 117 | M6 | H | MCQ | B | **UB** |
| 118 | M6 | H | MSQ | A, B, D | — (aliasing rules) |
| 119 | M6 | M | MCQ | D | **UB** |
| 120 | M7 | M | MCQ | A | Well-defined |

- **Difficulty:** 4 E / 9 M / 7 H.
  - M6 total (Q85–Q119): 7/14/14, as planned.
- **Format:** 15 MCQ / 5 MSQ.
  - M6 total: 27/8, as planned.
- **MCQ answer keys:** A×4, B×4, C×4, D×3. Running total Q1–Q120: A 24, B 24, C 24, D 24.

---

**Completed: 120 / 400 Remaining: 280**

**Ledger, newly ✓ this batch:**
- P17 → 3 (Q101)
- P19 → 7 (Q105, Q106)
- P22 ×3 (Q102, Q103, Q104)
- P23 → 4 (Q111, Q112, Q113)
- P24 ×2 (Q114, Q115; plus Q113 secondary)
- P25 ×2 (Q109, Q110) ★ complete
- P27 → 2 (Q119)
- P30 ×2 (Q117, Q118)
- P32 ×1 (Q120)

**M6 complete** (Q85–Q119, 35 questions). **M7 started** (1 of 12 done).

*Say **"Continue."** for Batch 7 (Q121–Q140: the rest of M7 strings — `strcpy`/`strncpy`/`strcmp`/`strchr`/`strstr`/`strtok`/`snprintf`, `strlen` on unterminated buffers, buffer overflow — then M8 structs and padding begins).*
