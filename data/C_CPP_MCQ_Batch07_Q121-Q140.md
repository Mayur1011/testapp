# C/C++ MCQ/MSQ Bank — Batch 7 (Q121–Q140)

**What this batch covers:**
- **M7 Strings, finishing (Q121–Q131):**
  - whether identical string literals share storage
  - returning a string literal vs returning a local array
  - `strcpy` from the middle of a string
  - the sign of `strcmp`
  - `strchr` / `strrchr` / `strstr` positions
  - `sizeof` of a string parameter
  - `strlen` with no `'\0'` (UB)
  - the two `strncpy` traps
  - which string calls overflow a buffer
  - `snprintf`'s return value when it truncates
  - how `strtok` modifies its input
- **M8 Structs, starting (Q132–Q140):**
  - padding and member order
  - trailing padding and nested structs
  - member offsets inside arrays of structs
  - designated initialisers for structs
  - passing a struct by value vs by pointer
  - `->` vs `.`
  - pointer arithmetic over an array of structs
  - nested struct layout
  - `offsetof`

**Patterns:** P20, P25, P32, P33

> **New in every question from this batch: "In simple words".** Each question now has a short plain-language explanation of the idea, usually with an everyday analogy, before the detailed technical reasoning. Read it first if a concept is new, then the step-by-step "Why".

> **Conventions for every question in this batch:**
> - Standard: **C11**. The questions don't depend on any particular compiler.
> - Platform: a typical 64-bit system: `char` 1, `short` 2, `int` 4, `long` 8, `double` 8, pointers 8 bytes; plain `char` is signed; ASCII.
> - **Struct layout:** natural alignment. Each member is aligned to its own size (`char` 1, `short` 2, `int` 4, `double`/pointer 8), and a struct's size is rounded up to a multiple of its largest member's alignment. The standard leaves padding *implementation-defined*, but this is the rule every mainstream 64-bit ABI uses, and the one GATE assumes.
> - The needed standard headers are included.

---

## Q121. [M7 | Do identical literals share storage? | Medium | MCQ | [C] · C11 | All] · Type T3 · P25

```c
#include <stdio.h>
int main(void) {
    char *p = "abc", *q = "abc";
    char a[] = "abc", b[] = "abc";
    printf("%d %d\n", p == q, a == b);
    return 0;
}
```

Which statement is correct?

- A. It always prints `1 1`.
- B. `p == q` is **unspecified** (it may be 0 or 1), while `a == b` is always 0.
- C. It always prints `0 0`.
- D. Both comparisons are unspecified.

**Answer:** B

**Behaviour class:** Unspecified (`p == q`). Well-defined (`a == b` → 0).

**In simple words:** A string literal is like a printed notice the compiler stores once in read-only memory. If you write the same notice twice, the compiler is *allowed* to keep only one copy, but it doesn't *have* to. So two pointers to `"abc"` may or may not hold the same address. An array like `char a[] = "abc"` is different: it's your own photocopy. Two photocopies are always two separate sheets, at two different addresses.

**Why:**
- N1570 §6.4.5p7: *"It is unspecified whether these arrays are distinct provided their elements have the appropriate values."* So the compiler may merge identical literals, and `p == q` can be 1 or 0.
- `a` and `b` are two separate arrays, two distinct objects, so their addresses always differ. (`a == b` compares addresses, because both arrays decay to pointers.)

**Why the other options fail:**
- A: Treats merging as guaranteed. It's allowed, not required.
- C: Treats merging as forbidden.
- D: Two distinct array objects can never have the same address.

**Concept recap:** `==` on strings compares **addresses**, never contents. Compare contents with `strcmp(a, b) == 0`.

**Trap alert:** `if (s == "yes")` compiles, and may even "work" sometimes because of merging, but it's a bug. In C++, `std::string == "yes"` does compare contents.

**Interview angle:** *"Why doesn't `str == "hello"` work in C?"* "`==` compares two pointers, not the characters. Use `strcmp`. Whether identical literals share an address is unspecified, so you can't rely on pointer equality either."

---

## Q122. [M7 | Returning a string literal vs a local array | Easy | MCQ | [C] · C11 | GATE+OA] · Type T5 · P25

```c
const char *f(void) { return "ok"; }
char *g(void) { char s[] = "no"; return s; }
```

The caller prints the result of each function with `printf("%s", …)`. Which is safe?

- A. Only `f`
- B. Only `g`
- C. Both
- D. Neither

**Answer:** A

**Behaviour class:** `f` is well-defined. Using `g`'s result is **UB** (a dangling pointer).

**In simple words:** A string literal lives for the **whole program**, like a sign permanently painted on a wall, so handing out its address is always safe. The local array `s` is a sticky note that is thrown away the moment `g` returns. Handing out its address is like giving someone directions to a note that no longer exists.

**Why:**
- `"ok"` has **static storage duration** (N1570 §6.4.5p6), so it exists for the entire run, and the returned pointer stays valid.
- `s` is an automatic array. Its lifetime ends when `g` returns, so the returned pointer dangles (§6.2.4p2). Compilers warn "function returns address of local variable".

**Why the other options fail:**
- B: Gets the two functions backwards.
- C: Ignores that `s`'s lifetime ends.
- D: Literals are not destroyed when a function returns.

**Concept recap:** Safe ways to return a string:
- a literal (read-only, so return `const char *`);
- a `static` buffer (shared between calls, so not reentrant);
- `malloc`'d memory (the caller must `free` it);
- a buffer passed in by the caller.

**Trap alert:** If `f` returned `char *` and a caller wrote through it, that would modify a literal, which is UB (Q109). That's why `f` returns `const char *`.

**Interview angle:** *"How do you return a string from a C function?"* List the four options above with their trade-offs. Interviewers like hearing "caller-provided buffer plus length" (as `snprintf` does) as the most robust choice.

---

## Q123. [M7 | strcpy from the middle, then truncate | Easy | MCQ | [C] · C11 | GATE] · Type T1 · P32

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char a[10] = "hello";
    char b[10];
    strcpy(b, a + 2);
    b[1] = '\0';
    printf("%s %zu %zu\n", b, strlen(a), strlen(b));
    return 0;
}
```

What does this program print?

- A. `llo 5 3`
- B. `he 5 2`
- C. `l 5 3`
- D. `l 5 1`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** `a + 2` means "start reading from the 3rd letter", so `strcpy` copies `"llo"`. Putting `'\0'` at position 1 is like putting a full stop after the first letter: everything after it is still in memory but is no longer part of the string.

**Why:**

| Line | Statement | b | strlen(b) |
|---|---|---|---|
| 4 | `a = "hello"` | — | — |
| 6 | `strcpy(b, a + 2)` copies from index 2: `l l o \0` | `llo` | 3 |
| 7 | `b[1] = '\0'` | `l \0 o \0` → reads as `"l"` | **1** |
| 8 | print | `l`, strlen(a) = 5, strlen(b) = 1 | |

**Why the other options fail:**
- A: Ignores line 7.
- B: Copies from the start of `a`.
- C: Thinks `strlen` is unaffected by the new `'\0'`. It stops at the first `'\0'`.

**Concept recap:** A C string ends at the **first** `'\0'`. Writing a `'\0'` into the middle truncates the string instantly, without moving any memory.

**Trap alert:** The `'o'` is still in `b[2]`. `sizeof b` is still 10. Only the *string view* got shorter.

**Interview angle:** *"Trim a string at the first newline."* `s[strcspn(s, "\n")] = '\0';`. This is the standard idiom after `fgets`, using exactly this idea.

---

## Q124. [M7 | Sign of strcmp | Medium | MSQ | [C] · C11 | GATE+OA] · Type T1 · P32

Assume ASCII. Select **all** calls that return a **positive** value.

- A. `strcmp("abc", "abd")`
- B. `strcmp("b", "abc")`
- C. `strcmp("abc", "ab")`
- D. `strcmp("Zebra", "apple")`
- E. `strcmp("apple", "Apple")`

**Answer:** B, C, E

**Behaviour class:** Well-defined. Only the **sign** of the result is guaranteed, not its size.

**In simple words:** `strcmp` works like looking words up in a dictionary, letter by letter, except that its "alphabet" is the ASCII table. In ASCII **every uppercase letter comes before every lowercase letter**, and "end of word" comes before any letter. It returns positive if the first string would come *after* the second.

**Why:** It compares characters as `unsigned char` until they differ or a `'\0'` is reached (N1570 §7.24.4.2).

| Option | First difference | Comparison | Sign |
|---|---|---|---|
| A | index 2: 'c' (99) vs 'd' (100) | 99 < 100 | negative |
| B | index 0: 'b' (98) vs 'a' (97) | 98 > 97 | **positive** ✓ |
| C | index 2: 'c' (99) vs '\0' (0) | 99 > 0 | **positive** ✓ |
| D | index 0: 'Z' (90) vs 'a' (97) | 90 < 97 | negative |
| E | index 0: 'a' (97) vs 'A' (65) | 97 > 65 | **positive** ✓ |

**Concept recap:** `strcmp` < 0, == 0, > 0 means "a sorts before b", "equal", "a sorts after b". Use `strncmp` to compare at most n characters. `strcasecmp` (POSIX, not ISO C) ignores case.

**Trap alert:** Don't test `strcmp(a, b) == 1` or `== -1`. Implementations may return any positive or negative value (often the difference between the two characters).

**Interview angle:** *"Implement strcmp."* `while (*a && *a == *b) { a++; b++; } return (unsigned char)*a - (unsigned char)*b;`. Mention the `unsigned char` cast: without it, characters ≥ 128 compare wrongly when `char` is signed.

---

## Q125. [M7 | strchr / strrchr / strstr positions | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P32

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    const char *s = "programming";
    const char *p = strchr(s, 'g');
    const char *q = strrchr(s, 'g');
    const char *r = strstr(s, "am");
    printf("%td %td %td %s\n", p - s, q - s, r - s, r);
    return 0;
}
```

What does this program print?

- A. `3 10 5 am`
- B. `4 11 6 amming`
- C. `3 10 5 amming`
- D. `3 3 5 amming`

**Answer:** C

**Behaviour class:** Well-defined. Every search succeeds, so no NULL pointer is used.

**In simple words:** These functions don't copy anything. They return a **pointer to where they found a match**, like a finger placed on the page. Subtracting the start of the string gives the 0-based index. Printing the pointer with `%s` shows everything from that finger to the end.

**Why:**

| Index | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| s | p | r | o | **g** | r | **a** | m | m | i | n | **g** |

| Call | Finds | Index |
|---|---|---|
| `strchr(s, 'g')` | the **first** 'g' | 3 |
| `strrchr(s, 'g')` | the **last** 'g' | 10 |
| `strstr(s, "am")` | the first occurrence of the substring "am" | 5 → `%s` prints `amming` |

**Why the other options fail:**
- A: Prints only the match. `%s` prints up to the `'\0'`.
- B: Uses 1-based positions.
- D: Treats `strrchr` like `strchr`.

**Concept recap:** All three return `NULL` when nothing matches. Always check before using the result. `strchr(s, '\0')` returns a pointer to the terminator (it counts as part of the string).

**Trap alert:** In C these functions take `const char *` but return a plain `char *`, a historical wart, so you could accidentally write into a literal through the result. Store the result as `const char *`, as here.

**Interview angle:** *"Find the file extension / split a key=value pair."* `strchr(s, '=')` gives the separator. The key is `[s, eq)` and the value is `eq + 1`. Interviewers like answers that don't allocate.

---

## Q126. [M7 | sizeof of a string inside a function | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P20

```c
#include <stdio.h>
#include <string.h>
void show(char s[]) { printf("%zu %zu\n", sizeof s, strlen(s)); }
int main(void) {
    char msg[50] = "hi";
    printf("%zu %zu %zu\n", sizeof msg, strlen(msg), sizeof "hi");
    show(msg);
    return 0;
}
```

What does this program print?

- A. `50 2 2` then `50 2`
- B. `50 2 3` then `8 2`
- C. `3 2 3` then `8 2`
- D. `50 2 3` then `50 2`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** `sizeof` measures the **box**; `strlen` counts the **letters in it**. `msg` is a 50-byte box holding 2 letters. A literal `"hi"` is a box made exactly to fit: 2 letters plus the end marker = 3. Inside `show`, `s` isn't a box at all, just a **note with the box's address**, and a note (pointer) is 8 bytes.

**Why:**

| Expression | Type | Value |
|---|---|---|
| `sizeof msg` | `char[50]` | 50 |
| `strlen(msg)` | — | 2 |
| `sizeof "hi"` | `char[3]` (the literal includes its `'\0'`) | 3 |
| `sizeof s` in `show` | `char *` (an array parameter is adjusted to a pointer) | 8 |
| `strlen(s)` | — | 2 |

**Why the other options fail:**
- A: Forgets the `'\0'` in `sizeof "hi"`, and treats the parameter as an array.
- C: Sizes `msg` from its initializer. The declared `[50]` wins.
- D: Treats the parameter as an array.

**Concept recap:** Inside a function, the only way to know a buffer's capacity is to pass it as a separate argument. That's why safe APIs look like `fgets(buf, sizeof buf, f)` and `snprintf(buf, sizeof buf, …)`, with `sizeof` used at the call site, where the array type is still known.

**Trap alert:** `snprintf(s, sizeof s, …)` **inside** `show` would pass 8, not 50, silently truncating the output.

**Interview angle:** *"Why do C string functions need a size parameter?"* "Arrays decay to pointers when passed, so the function can't learn the buffer's size any other way. Without it, overflow is unavoidable."

---

## Q127. [M7 | strlen on a buffer with no '\0' | Medium | MCQ | [C] · C11 | All] · Type T3 · P32

> **This code has undefined behaviour.**

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char buf[4];
    memcpy(buf, "ABCD", 4);
    printf("%zu\n", strlen(buf));
    return 0;
}
```

Which statement is correct?

- A. It prints `4`.
- B. It prints `5`.
- C. The behaviour is undefined: `buf` holds no `'\0'`, so `strlen` reads past the end of the array.
- D. It prints `0`.

**Answer:** C

**Behaviour class:** **Undefined behaviour.**

**In simple words:** `strlen` is like someone counting train carriages until they see the "END" sign. `memcpy` copied exactly 4 letters and **no END sign**, so the counter walks off the end of the train into whatever happens to be parked next to it.

> **⚠ UB box**
> - **Category:** out-of-bounds read through a string function on a non-terminated array.
> - **Rule:** N1570 §7.24.1p1 (string functions operate on null-terminated arrays) together with §6.5.6p8 (reading past the end of an array).
> - **What you may see in practice:** often a number a little larger than 4 (it keeps counting until it hits some zero byte in neighbouring memory), sometimes exactly 4 by luck, occasionally a crash. **Not guaranteed.**

**Why:** `memcpy(buf, "ABCD", 4)` copies 4 bytes: `A B C D`. The literal's `'\0'` is the 5th byte, and it was not copied. `buf` is a valid char array, but it isn't a *string*.

**Why the other options fail:**
- A, B, D: Each is one possible lucky outcome, presented as a guarantee.

**Concept recap:** Every function that takes a "string" assumes a `'\0'` within the array. Copy `n + 1` bytes, or terminate explicitly: `buf[n] = '\0'` (and size the buffer `n + 1`).

**Trap alert:** The same bug hides in `char t[4] = "ABCD";` (legal in C, no terminator; Q80) and in `strncpy` (Q128).

**Interview angle:** *"What's the most common source of C string bugs?"* "A missing terminator or a buffer one byte too small. Size buffers as length + 1, prefer `snprintf`, and test with AddressSanitizer."

---

## Q128. [M7 | The two strncpy traps | Hard | MCQ | [C] · C11 | All] · Type T1 · P32

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char d[6];
    strncpy(d, "planet", sizeof d);
    d[sizeof d - 1] = '\0';
    printf("%s ", d);

    char e[8];
    strncpy(e, "ab", sizeof e);
    printf("%s %d\n", e, e[7]);
    return 0;
}
```

What does this program print?

- A. `plane ab 0`
- B. `planet ab 0`
- C. `plane ab` followed by a garbage value
- D. Nothing is guaranteed: `d` is never terminated

**Answer:** A

**Behaviour class:** Well-defined. Line 6 adds the terminator that `strncpy` left out.

**In simple words:** `strncpy(dst, src, n)` always writes **exactly n characters**, with two odd rules. (1) If `src` is **too long**, it fills all n slots and adds **no end marker**, leaving an unterminated string. (2) If `src` is **short**, it fills every remaining slot with `'\0'`. It's a fixed-width form-filling function, not a "safe strcpy".

**Why:**

| Line | Statement | Buffer contents |
|---|---|---|
| 5 | `strncpy(d, "planet", 6)`: "planet" has 6 characters, so all 6 slots are filled and **no `'\0'`** is added | `p l a n e t` (unterminated) |
| 6 | `d[5] = '\0'` → manual terminator, overwriting `'t'` | `p l a n e \0` → `"plane"` |
| 10 | `strncpy(e, "ab", 8)`: the source is short, so the rest is padded with zeros | `a b \0 \0 \0 \0 \0 \0` |
| 11 | `e[7]` | 0 (guaranteed by the padding) |

Prints `plane ab 0`.

**Why the other options fail:**
- B: Forgets that line 6 overwrites the last slot with `'\0'`.
- C: `strncpy` **guarantees** the zero padding (N1570 §7.24.2.4p3).
- D: Line 6 fixes the missing terminator before `d` is printed.

**Concept recap:** The safe idiom: `strncpy(d, src, sizeof d - 1); d[sizeof d - 1] = '\0';`. Better still: `snprintf(d, sizeof d, "%s", src)`, which always terminates and reports truncation.

**Trap alert:** Removing line 6 would make `printf("%s", d)` UB, because it reads past the array looking for a `'\0'` (as in Q127). The zero-padding also makes `strncpy` slow when copying a short string into a huge buffer.

**Interview angle:** *"Is strncpy a safe strcpy?"* "No. It doesn't guarantee termination and it pads wastefully. It was designed for fixed-width fields in old Unix directory entries. Use `snprintf`, or `strlcpy` where it's available, and terminate manually if you must use `strncpy`."

---

## Q129. [M7 | Which calls overflow the buffer? | Hard | MSQ | [C] · C11 | All] · Type T5

Each option starts from a fresh `char d[5];`. Select **all** calls that **write outside** `d` (undefined behaviour).

- A. `strcpy(d, "four");`
- B. `strcpy(d, "fives");`
- C. `strncpy(d, "fives", 5);`
- D. `snprintf(d, sizeof d, "%s", "fives");`
- E. `strcat(strcpy(d, "ab"), "cde");`

**Answer:** B, E

**Behaviour class per option:** A = OK · B = **UB** · C = OK, but leaves `d` unterminated · D = OK (truncates to `"five"`) · E = **UB**

**In simple words:** Count the bytes each call needs, **including the end marker**, and compare with the 5 seats available. If a call needs a 6th seat, it sits in someone else's memory. `snprintf` is the polite one: it stops at the last seat and always leaves room for the end marker.

**Why:**

| Option | Bytes written | Fits in 5? |
|---|---|---|
| A | "four" + `'\0'` = 5 | ✓ exactly |
| B | "fives" + `'\0'` = **6** | ✗ overflow → UB |
| C | exactly 5 (`'\0'` dropped) | ✓ no overflow, but not a valid string afterwards |
| D | at most 5: `"five"` + `'\0'` (truncated) | ✓ |
| E | "ab" + "cde" + `'\0'` = **6** | ✗ overflow → UB |

**Concept recap:** `strcpy` and `strcat` **never check** the destination's size. The caller must guarantee `strlen(src) + 1` (or `strlen(dst) + strlen(src) + 1` for `strcat`) bytes of space.

**Trap alert:** Option C is "safe" from overflow, but the next `printf("%s", d)` or `strlen(d)` would read out of bounds. That's a delayed bug.

**Interview angle:** *"What is a buffer overflow and why is it a security risk?"* "Writing past an array's end overwrites neighbouring memory: other variables, saved registers, or the return address on the stack. Attackers use that to take over control flow. Defences: bounds-checked APIs, stack canaries, ASLR, non-executable stacks, and memory-safe languages or containers."

---

## Q130. [M7 | snprintf's return value when it truncates | Hard | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char buf[8];
    int n = snprintf(buf, sizeof buf, "%d-%s", 2026, "GATE");
    printf("%d %s %zu\n", n, buf, strlen(buf));
    return 0;
}
```

What does this program print?

- A. `7 2026-GA 7`
- B. `8 2026-GAT 8`
- C. `-1 2026-GA 7`
- D. `9 2026-GA 7`

**Answer:** D

**Behaviour class:** Well-defined. Truncation is defined behaviour for `snprintf`.

**In simple words:** `snprintf` is like a printer with a small page. It prints as much as fits and always ends the page properly. Then it tells you how many characters the **full** message would have needed, so you can tell it got cut off, and how big a page you'd need.

**Why:**

| Step | Value |
|---|---|
| Full text `"%d-%s"` → `2026-GATE` | 9 characters |
| Capacity `sizeof buf` = 8 → at most 7 characters + `'\0'` | `2026-GA` |
| Return value: the length it **would** have written | **9** |
| `strlen(buf)` | 7 |

Prints `9 2026-GA 7`. Truncation test: `n >= (int)sizeof buf` → 9 ≥ 8 → truncated.

**Why the other options fail:**
- A: Thinks it returns the number of characters actually stored.
- B: Forgets to reserve a slot for `'\0'`.
- C: A negative return means an *encoding error*, not truncation.

**Concept recap:** The two-pass idiom: `int need = snprintf(NULL, 0, fmt, …);` gives the required length; then allocate `need + 1` bytes and call `snprintf` again. (A size of 0 with a NULL buffer is explicitly allowed, §7.21.6.5p2.)

**Trap alert:** `sprintf` (without the `n`) has no size limit, so it's the same overflow risk as `strcpy`. Prefer `snprintf` everywhere.

**Interview angle:** *"How do you safely build a formatted string of unknown length in C?"* The two-pass `snprintf` idiom above. In C++, `std::format` (C++20) or `std::ostringstream`.

---

## Q131. [M7 | strtok modifies its input | Hard | MCQ | [C] · C11 | GATE+OA] · Type T1 · P32

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char s[] = "a,,b;c";
    int count = 0;
    for (char *t = strtok(s, ",;"); t; t = strtok(NULL, ",;")) {
        printf("[%s]", t);
        count++;
    }
    printf(" %d %zu\n", count, strlen(s));
    return 0;
}
```

What does this program print?

- A. `[a][][b][c] 4 6`
- B. `[a][b][c] 3 1`
- C. `[a][b][c] 3 6`
- D. `[a][][b][c] 4 1`

**Answer:** B

**Behaviour class:** Well-defined. `s` is a writable array, which `strtok` requires.

**In simple words:** `strtok` cuts a string into pieces by **writing `'\0'` over the separators in your own string**, like cutting a ribbon with scissors. It also treats several separators in a row as one cut, so it never returns empty pieces. Afterwards, your original string looks like just its first piece, because the first cut ends it.

**Why:**

| Call | Skips leading separators | Token | Writes `'\0'` at | s in memory |
|---|---|---|---|---|
| `strtok(s, ",;")` | — | `"a"` | index 1 (the first ',') | `a\0,b;c` |
| `strtok(NULL, …)` | skips the ',' at index 2 | `"b"` | index 4 (';') | `a\0,b\0c` |
| `strtok(NULL, …)` | — | `"c"` | reaches the end | |
| `strtok(NULL, …)` | — | `NULL` → loop ends | | |

count = 3. `strlen(s)` now stops at index 1 → **1**.

**Why the other options fail:**
- A, D: Expect an empty token between the two commas. `strtok` collapses runs of separators.
- C: Forgets that `strtok` overwrote separators with `'\0'`.

**Concept recap:** `strtok`:
- needs a **writable** string (passing a literal is UB);
- keeps hidden static state between calls, so it is **not reentrant** and not thread-safe;
- never returns empty tokens.

For empty fields, use POSIX `strsep`. For reentrancy, use `strtok_r` (POSIX) or `strtok_s` (C11 Annex K, optional).

**Trap alert:** Calling `strtok` on a *different* string inside this loop resets its hidden state and breaks the outer loop.

**Interview angle:** *"Why is strtok considered dangerous?"* "It modifies its input, uses hidden global state (so it can't be nested or used across threads), and silently drops empty fields, which is wrong for CSV parsing. Prefer `strtok_r`, or a manual `strchr` loop. In C++, use `std::string_view` with `find`."

---

## Q132. [M8 | Padding depends on member order | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P33

```c
#include <stdio.h>
struct A { char c; int i; char d; };
struct B { char c; char d; int i; };
int main(void) {
    printf("%zu %zu\n", sizeof(struct A), sizeof(struct B));
    return 0;
}
```

Under natural alignment, what does this program print?

- A. `6 6`
- B. `12 8`
- C. `8 8`
- D. `12 12`

**Answer:** B

**Behaviour class:** Well-defined. The sizes are implementation-defined, and fixed by the stated natural-alignment rule.

**In simple words:** Imagine a car park where an `int` may only park in a bay whose number is a multiple of 4. A `char` can park anywhere. If a `char` takes bay 0, the `int` must skip bays 1–3 and park at 4, and those skipped bays are **padding**. Put the small members together, and fewer bays are wasted.

**Why:**

**struct A**

| Offset | 0 | 1–3 | 4–7 | 8 | 9–11 |
|---|---|---|---|---|---|
| Content | c | pad | i | d | **trailing pad** |

Members end at 9. The struct's alignment is 4 (from `int`), so the size is rounded up to 12. (Trailing padding ensures that in an array `struct A arr[2]`, `arr[1].i` is also 4-aligned.)

**struct B**

| Offset | 0 | 1 | 2–3 | 4–7 |
|---|---|---|---|---|
| Content | c | d | pad | i |

Size = 8.

**Why the other options fail:**
- A: Adds up the member sizes and ignores padding.
- C: Forgets the trailing padding in `A`.
- D: Assumes reordering the members can't help.

**Concept recap:** The layout rules:
1. Each member starts at an offset that is a multiple of its alignment.
2. The struct's alignment = the largest member alignment.
3. The size is rounded up to a multiple of the struct's alignment.

**C never reorders members**: they appear in declaration order (§6.7.2.1p15).

**Trap alert:** GATE questions that "assume no padding" are asking for the plain sum of member sizes. Read the assumption before calculating.

**Interview angle:** *"How do you reduce a struct's size?"* "Order members from largest to smallest alignment. Packing (`#pragma pack`, `__attribute__((packed))`) is non-standard and can make access slower or even fault on some CPUs, so use it only for wire formats."

---

## Q133. [M8 | double members, trailing padding, nesting | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P33

```c
#include <stdio.h>
struct S { char a; double b; short c; };
struct T { double b; short c; char a; };
struct U { struct T t; char x; };
int main(void) {
    printf("%zu %zu %zu\n", sizeof(struct S), sizeof(struct T), sizeof(struct U));
    return 0;
}
```

Under natural alignment, what does this program print?

- A. `11 11 12`
- B. `16 16 24`
- C. `24 16 24`
- D. `24 16 17`

**Answer:** C

**Behaviour class:** Well-defined. The sizes are implementation-defined, and fixed by the stated rule.

**In simple words:** A `double` must sit at a multiple of 8. If a lone `char` comes first, 7 bytes are wasted before the `double`. And because the whole struct must be a multiple of 8 (so that arrays of it keep every `double` aligned), the end gets padded too. A struct nested inside another struct keeps its own size and its own alignment.

**Why:**

| Struct | Layout (offset: member) | Members end at | Alignment | Size |
|---|---|---|---|---|
| S | 0: a, 1–7: pad, 8–15: b, 16–17: c, 18–23: pad | 18 | 8 | **24** |
| T | 0–7: b, 8–9: c, 10: a, 11–15: pad | 11 | 8 | **16** |
| U | 0–15: t (alignment 8), 16: x, 17–23: pad | 17 | 8 | **24** |

**Why the other options fail:**
- A: No padding at all.
- B: Gets T right but under-counts S.
- D: Forgets U's trailing padding. Its alignment is 8, because it contains a `double`.

**Concept recap:** Sorting members by size (largest first) cut S's 24 bytes to T's 16, with the same three members. Across a million-element array, that saves 8 MB.

**Trap alert:** Nesting doesn't "merge" padding. `x` can't use T's 5 bytes of trailing padding, because `t` is a complete object that must keep its own size.

**Interview angle:** *"Why does a struct have padding at the end?"* "So that every element of an array of that struct is correctly aligned. The size must be a multiple of the struct's alignment."

---

## Q134. [M8 | Member offset inside an array of structs | Hard | MCQ | [C] · C11 | GATE] · Type T1 · P33

```c
#include <stdio.h>
struct P { int id; char tag; };
int main(void) {
    struct P arr[10];
    printf("%zu %td\n", sizeof arr, (char *)&arr[3].tag - (char *)arr);
    return 0;
}
```

Under natural alignment, what does this program print?

- A. `50 15`
- B. `80 24`
- C. `50 19`
- D. `80 28`

**Answer:** D

**Behaviour class:** Well-defined. Only the size and an address are computed; nothing is read.

**In simple words:** Each `struct P` takes 5 bytes of data but occupies **8** bytes (padded so the next element's `int` starts at a multiple of 4). So the elements sit 8 bytes apart, like houses on plots of equal width. To find `tag` in house #3: skip 3 whole plots (24 bytes), then walk to `tag` inside the house (4 bytes in).

**Why:**

| Item | Computation | Value |
|---|---|---|
| `sizeof(struct P)` | id 0–3, tag 4, pad 5–7 | 8 |
| `sizeof arr` | 10 × 8 | **80** |
| Start of arr[3] | 3 × 8 | 24 |
| `tag` inside an element | offset 4 | |
| `&arr[3].tag` − start | 24 + 4 | **28** |

**Why the other options fail:**
- A, C: Use 5 bytes per element (no padding).
- B: Finds the start of arr[3] but forgets tag's offset within it.

**Concept recap:** Address(arr[i].m) = base + i × sizeof(struct) + offsetof(struct, m). This is the same formula as for a 2D array, with the struct as one "row".

**Trap alert:** Use `sizeof(struct P)`, never the sum of the member sizes, when computing strides or reading binary records.

**Interview angle:** *"Array of structs vs struct of arrays?"* "AoS keeps each record together, which is good for accessing whole records. SoA (`int id[N]; char tag[N];`) has no per-record padding and suits SIMD and cache-friendly scans of a single field. Game engines and databases often use SoA."

---

## Q135. [M8 | Designated initialisers for structs | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
struct Pt { int x, y, z; };
int main(void) {
    struct Pt p = { .z = 3, .x = 1 };
    struct Pt q = { 7 };
    printf("%d %d %d %d %d\n", p.x, p.y, p.z, q.x, q.y);
    return 0;
}
```

What does this program print?

- A. `1 0 3 7 0`
- B. `3 1 0 7 0`
- C. `1 3 0 7 7`
- D. `1 <garbage> 3 7 <garbage>`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** `.z = 3` is like filling in a form by field name, so the order doesn't matter. Any field you leave blank is set to **zero**, as long as you filled in *at least one* field. `{ 7 }` without names fills the fields in order, so only the first field gets 7.

**Why:**

| Object | x | y | z |
|---|---|---|---|
| p | 1 (named) | **0** (not mentioned) | 3 (named) |
| q | 7 (first in order) | **0** | **0** |

**Why the other options fail:**
- B: Assigns the values in the order they are written, ignoring the names.
- C: Gives `p.y` the value 3.
- D: A partial initializer zero-fills the remaining members, even for a local variable (N1570 §6.7.9p21).

**Concept recap:** `struct Pt zero = {0};` zero-initialises the entire struct. This is the standard idiom.

**Trap alert:** C++ supports `.field =` only from C++20, only in declaration order, and without mixing positional and named values. `{ .z = 3, .x = 1 }` is **ill-formed in C++**.

**Interview angle:** *"Why use designated initialisers?"* "They stay correct when fields are reordered or added, document intent, and guarantee zeroed defaults. Linux kernel driver tables use them heavily (`.open = my_open`)."

---

## Q136. [M8 | Struct by value vs by pointer | Medium | MCQ | [C] · C11 | GATE+OA] · Type T5

```c
#include <stdio.h>
struct Box { int v[3]; int n; };
void bump(struct Box b) { b.v[0] = 99; b.n++; }
void bump_p(struct Box *b) { b->v[1] = 77; b->n++; }
int main(void) {
    struct Box a = {{1, 2, 3}, 0};
    struct Box c = a;
    bump(a);
    bump_p(&c);
    printf("%d %d %d | %d %d %d\n", a.v[0], a.v[1], a.n, c.v[0], c.v[1], c.n);
    return 0;
}
```

What does this program print?

- A. `99 2 1 | 99 77 2`
- B. `1 2 0 | 99 77 1`
- C. `1 2 0 | 1 77 1`
- D. `99 2 1 | 1 77 1`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Passing a struct by value hands the function a **photocopy**, including any arrays inside it. Scribbling on the photocopy leaves the original clean. Passing `&c` hands over the original's address, so the changes stick. And `c = a` makes a full, independent copy: changing `c` never affects `a`.

**Why:**

| Line | Action | a | c |
|---|---|---|---|
| 6 | init | {1,2,3}, n=0 | — |
| 7 | `c = a` (deep copy of the array member too) | same | {1,2,3}, n=0 |
| 8 | `bump(a)`: changes a **copy** | unchanged | |
| 9 | `bump_p(&c)`: changes c | | {1,**77**,3}, n=**1** |
| 10 | print | `1 2 0` | `1 77 1` |

**Why the other options fail:**
- A: Assumes `bump` changes `a`, and that `c` shares storage with `a`.
- B: Assumes `c = a` shares the array.
- D: Assumes passing by value still changes the original.

**Concept recap:** Unlike a bare array parameter (which decays to a pointer), an array **inside a struct** is copied by value. That's the standard trick for passing or returning an array by value in C.

**Trap alert:** If the struct contained a **pointer** (`int *v`), the copy would duplicate the pointer, not the data it points to (a shallow copy), and both structs would share that data. That's covered in M9 (deep vs shallow copy).

**Interview angle:** *"Should you pass large structs by value?"* "Usually not. Copying is O(size). Pass `const struct S *` for read-only access and `struct S *` to modify. Small structs (2–3 words) are fine by value, and are often passed in registers."

---

## Q137. [M8 | `->` vs `.` | Medium | MSQ | [C] · C11 | GATE+OA] · Type T2

```c
struct N { int val; struct N *next; };
struct N n2 = {20, NULL};
struct N n1 = {10, &n2};
struct N *p = &n1;
```

Select **all** expressions that are valid and evaluate to **20**.

- A. `p->next->val`
- B. `(*p).next->val`
- C. `*p.next->val`
- D. `(*(*p).next).val`
- E. `++p->val`

**Answer:** A, B, D

**In simple words:** Use `.` when you're holding the **struct itself**, and `->` when you're holding a **pointer** to it. `p->x` is just shorthand for `(*p).x`. The trap is that `.` binds tighter than `*`, so `*p.next` tries to use `.` on a pointer, which is like asking an address for its "next" field.

**Why each option:**
- **A ✓:** p → n1 → `next` (&n2) → `val` = 20.
- **B ✓:** `(*p)` is n1 → `.next` → `->val` = 20.
- **C ✗ ill-formed:** It parses as `*(p.next->val)`. `p` is a pointer, so `p.next` is a compile error.
- **D ✓:** `(*p).next` = &n2 → `*` gives n2 → `.val` = 20.
- **E ✗:** `->` binds tighter than prefix `++`, so this is `++(p->val)` → n1.val becomes **11**.

**Concept recap:** Precedence: postfix `->`, `.`, `[]`, `()` (highest) > unary `*`, `&`, `++`, `--`, `!` > the rest. `p->a.b->c` chains left to right.

**Trap alert:** `*p->next` means `*(p->next)`, which is n2 itself (the struct), not `(*p).next`.

**Interview angle:** *"Why does C have both `.` and `->`?"* "`->` combines a dereference and a member access into one readable operator. Without it, every linked-list step would be written `(*(*p).next).val`."

---

## Q138. [M8 | Pointer arithmetic over an array of structs | Medium | MCQ | [C] · C11 | GATE] · Type T5

```c
#include <stdio.h>
struct S { int x, y; };
int main(void) {
    struct S a[3] = {{1, 2}, {3, 4}, {5, 6}};
    struct S *p = a;
    p++;
    printf("%d %d %d\n", p->x, (p + 1)->y, (*p).y + a[0].x);
    return 0;
}
```

What does this program print?

- A. `2 3 3`
- B. `3 6 5`
- C. `3 4 5`
- D. `1 4 3`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** `p++` on a struct pointer jumps over **one whole struct** (8 bytes here), not over one `int`. So `p` moves from the first record to the second, exactly like moving to the next row in a table.

**Why:**

| Element | x | y |
|---|---|---|
| a[0] | 1 | 2 |
| a[1] | 3 | 4 |
| a[2] | 5 | 6 |

| Expression | p → | Value |
|---|---|---|
| `p++` | a[1] | |
| `p->x` | a[1].x | **3** |
| `(p + 1)->y` | a[2].y | **6** |
| `(*p).y + a[0].x` | 4 + 1 | **5** |

**Why the other options fail:**
- A: Moves `p` by one `int` (4 bytes) instead of one struct.
- C: Reads `(p + 1)->y` as `p->y`.
- D: Forgets the `p++`.

**Concept recap:** Pointer arithmetic always steps by `sizeof(*p)`, including for structs (padding included). That's why `p++` is the natural way to iterate over records.

**Trap alert:** `p->x + 1` means `(p->x) + 1` = 4. `(p + 1)->x` means the next element's x = 5. Watch the parentheses.

**Interview angle:** *"Iterate over an array of structs with a pointer."* `for (struct S *q = a; q != a + n; ++q) use(q->x);`. This uses the one-past-the-end pointer as a sentinel (Q82).

---

## Q139. [M8 | Nested struct access and layout | Hard | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
struct Date { int d, m, y; };
struct Student { char name[8]; struct Date dob; int marks[2]; };
int main(void) {
    struct Student s = {"Ravi", {5, 8, 2003}, {90, 85}};
    struct Student *ps = &s;
    printf("%c %d %d %d\n", ps->name[1], ps->dob.m, s.marks[1] - s.dob.d, (int)sizeof s);
    return 0;
}
```

Under natural alignment, what does this program print?

- A. `R 8 80 28`
- B. `a 5 80 32`
- C. `a 8 85 28`
- D. `a 8 80 28`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** A struct inside a struct is like a folder inside a folder: `ps->dob.m` means "open the student, open the date-of-birth folder, read the month". For the size, just lay the pieces end to end: 8 letters, then 3 ints for the date, then 2 ints for marks. Every piece is already 4-aligned, so nothing extra is needed.

**Why:**

| Expression | Evaluation | Value |
|---|---|---|
| `ps->name[1]` | "Ravi"[1] | **'a'** |
| `ps->dob.m` | the m field of {5, 8, 2003} | **8** |
| `s.marks[1] - s.dob.d` | 85 − 5 | **80** |

| Member | Offset | Size |
|---|---|---|
| name[8] | 0 | 8 |
| dob (3 ints, alignment 4) | 8 | 12 |
| marks[2] | 20 | 8 |
| **total** (alignment 4, 28 is a multiple of 4) | | **28** |

**Why the other options fail:**
- A: `name[1]` is the 2nd character.
- B: Reads `dob.d` instead of `dob.m`, and adds imaginary padding.
- C: Uses `marks[0]` instead of `marks[1]`: 90 − 5 = 85.

**Concept recap:** The brace structure of the initializer mirrors the nesting: `{"Ravi", {5, 8, 2003}, {90, 85}}`. With designators: `{ .dob = { .m = 8 } }`.

**Trap alert:** `ps->dob->m` is a compile error, because `dob` is a struct, not a pointer. Use `->` only when the thing on its left is a pointer.

**Interview angle:** *"Model a student record in C, then in C++."* Show this struct. Then, in C++, a class with a `std::string name`, a `Date` value member, and a `std::array<int,2>` for marks. Mention that the C++ version is no longer trivially copyable once it holds a `std::string`.

---

## Q140. [M8 | offsetof | Hard | MCQ | [C] · C11 | All] · Type T1 · P33

```c
#include <stdio.h>
#include <stddef.h>
struct R { char a; short b; char c; int d; double e; };
int main(void) {
    printf("%zu %zu %zu %zu %zu\n",
           offsetof(struct R, b), offsetof(struct R, c),
           offsetof(struct R, d), offsetof(struct R, e), sizeof(struct R));
    return 0;
}
```

Under natural alignment, what does this program print?

- A. `2 4 8 16 24`
- B. `1 3 4 8 16`
- C. `2 4 6 10 18`
- D. `2 4 8 12 20`

**Answer:** A

**Behaviour class:** Well-defined. The offsets are implementation-defined, and fixed by the stated rule.

**In simple words:** `offsetof(type, member)` answers "how many bytes from the start of the struct is this member?". Walk through the members one at a time: each one jumps forward to the next position that is a multiple of its own size, then occupies its bytes.

**Why:**

| Member | Size / alignment | Next free byte before it | Aligned offset | Ends at |
|---|---|---|---|---|
| a | 1 / 1 | 0 | **0** | 1 |
| b | 2 / 2 | 1 | **2** (1 pad byte) | 4 |
| c | 1 / 1 | 4 | **4** | 5 |
| d | 4 / 4 | 5 | **8** (3 pad bytes) | 12 |
| e | 8 / 8 | 12 | **16** (4 pad bytes) | 24 |
| size | alignment 8 | 24 | 24 is a multiple of 8 | **24** |

**Why the other options fail:**
- B: No padding at all.
- C: Aligns everything to 2.
- D: Aligns the `double` to 4 (that was true on some 32-bit ABIs, but not on the stated platform).

**Concept recap:** `offsetof` (from `<stddef.h>`) is a compile-time constant of type `size_t`. Typical uses: serialisation, `container_of` (getting from a member pointer back to the enclosing struct, as the Linux kernel does), and checking layouts with `_Static_assert(offsetof(struct R, e) == 16, "layout")`.

**Trap alert:** 8 of these 24 bytes are padding. Reordering to `double e; int d; short b; char a, c;` gives 8 + 4 + 2 + 1 + 1 = 16 bytes, with no padding at all.

**Interview angle:** *"How would you verify a struct's layout for a network protocol?"* "Use fixed-width types (`uint32_t`), order members to avoid padding, `_Static_assert` the `sizeof` and `offsetof` values, and handle byte order explicitly. Or serialise field by field and never `fwrite` the raw struct."

---

## Batch 7 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 121 | M7 | M | MCQ | B | **Unspecified** / well-defined |
| 122 | M7 | E | MCQ | A | Well-defined / UB (using g) |
| 123 | M7 | E | MCQ | D | Well-defined |
| 124 | M7 | M | MSQ | B, C, E | Well-defined |
| 125 | M7 | M | MCQ | C | Well-defined |
| 126 | M7 | E | MCQ | B | Well-defined |
| 127 | M7 | M | MCQ | C | **UB** |
| 128 | M7 | H | MCQ | A | Well-defined |
| 129 | M7 | H | MSQ | B, E | UB / OK mix |
| 130 | M7 | H | MCQ | D | Well-defined |
| 131 | M7 | H | MCQ | B | Well-defined |
| 132 | M8 | E | MCQ | B | Well-defined (impl-defined sizes) |
| 133 | M8 | M | MCQ | C | Well-defined (impl-defined sizes) |
| 134 | M8 | H | MCQ | D | Well-defined |
| 135 | M8 | E | MCQ | A | Well-defined |
| 136 | M8 | M | MCQ | C | Well-defined |
| 137 | M8 | M | MSQ | A, B, D | Ill-formed / OK mix |
| 138 | M8 | M | MCQ | B | Well-defined |
| 139 | M8 | H | MCQ | D | Well-defined |
| 140 | M8 | H | MCQ | A | Well-defined |

- **Difficulty:** 5 E / 8 M / 7 H.
  - M7 total (Q120–Q131): 3/5/4, as planned.
  - M8 so far: 2/4/3.
- **Format:** 17 MCQ / 3 MSQ.
  - M7 total: 10/2, as planned.
- **MCQ answer keys:** A×4, B×5, C×4, D×4. Running total Q1–Q140: A 28, B 29, C 28, D 28.

---

**Completed: 140 / 400 Remaining: 260**

**Ledger, newly ✓ this batch:**
- P25 → 4 (Q121, Q122)
- P32 → 9 (Q123, Q124, Q125, Q127, Q128, Q131) ★ complete
- P20 → 4 (Q126)
- P33 ×4 (Q132, Q133, Q134, Q140) ★ complete

**M7 complete** (Q120–Q131, 12 questions). **M8 in progress** (9 of 18 done).

*Say **"Continue."** for Batch 8 (Q141–Q160: the rest of M8 — unions and type-punning, union size, struct vs union, bit-fields, self-referential structs, anonymous members, flexible array members — then M9 dynamic memory begins).*
