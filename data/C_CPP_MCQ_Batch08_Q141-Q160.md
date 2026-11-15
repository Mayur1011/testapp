# C/C++ MCQ/MSQ Bank — Batch 8 (Q141–Q160)

**What this batch covers:**
- **M8 Structs & unions, finishing (Q141–Q149):**
  - how big a union is, and how a union differs from a struct
  - reading a union through a different member (type-punning) to see byte order or float bits
  - bit-field values and the rules for bit-fields
  - self-referential structs (linked-list nodes)
  - anonymous structs and unions inside a struct (C11)
  - flexible array members
- **M9 Dynamic memory, starting (Q150–Q160):**
  - `calloc` vs `malloc`
  - getting the `sizeof` wrong inside `malloc`
  - what happens when allocation fails, and `free(NULL)`
  - returning heap memory from a function
  - use-after-free through a second pointer, and double free
  - why `realloc` makes old pointers invalid
  - counting leaked bytes, and spotting which snippets leak
  - freeing a linked list correctly
  - forgetting the `+1` for `'\0'` when allocating a string

**Patterns:** P27, P28, P34, P35, P36, P45, P46

> **Rules in effect:** questions are compiler-independent (answers follow from the standard + the stated platform). Every question has an **"In simple words"** explanation, and UB questions describe typical behaviour only in general terms.

> **Conventions for every question in this batch:**
> - Standard: **C11** (C++ differences are noted where they matter).
> - Platform: a typical 64-bit system: `char` 1, `short` 2, `int` 4, `float` 4, `double` 8, pointers 8 bytes; plain `char` signed; IEEE-754; ASCII; **little-endian**.
> - Struct and union layout uses natural alignment (as in Batch 7).
> - Every `malloc`/`calloc` in a well-defined question is assumed to **succeed**, unless the question is about failure.
> - The needed standard headers (`stdio.h`, `stdlib.h`, `string.h`, `stddef.h`) are included.

---

## Q141. [M8 | Size of a union | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1

```c
#include <stdio.h>
union U { char c; int i; double d; };
union V { char s[5]; int i; };
int main(void) {
    printf("%zu %zu\n", sizeof(union U), sizeof(union V));
    return 0;
}
```

What does this program print?

- A. `13 9`
- B. `8 5`
- C. `8 8`
- D. `16 8`

**Answer:** C

**Behaviour class:** Well-defined. The sizes are implementation-defined, and fixed by the stated natural-alignment rule.

**In simple words:** A struct is a **row of separate lockers**: every member gets its own space, so the sizes add up. A union is **one locker shared by everyone**: all members use the same space, so it only needs to be as big as the largest member (then rounded up so that it stays properly aligned).

**Why:**

| Union | Largest member | Strictest alignment | Size |
|---|---|---|---|
| U | `double` 8 | 8 | **8** |
| V | `char s[5]` → 5 | 4 (from `int`) | 5 rounded up to a multiple of 4 → **8** |

**Why the other options fail:**
- A: Adds the member sizes, as for a struct (1 + 4 + 8, and 5 + 4).
- B: Forgets V's trailing padding (needed so an array of V keeps every `int` aligned).
- D: Thinks a union with a `double` needs 16.

**Concept recap:** sizeof(union) = the largest member size, rounded up to a multiple of the strictest member alignment. All members start at **offset 0**.

**Trap alert:** "Union size = largest member" is the GATE shortcut, but it misses the round-up. V would be 5 without it.

**Interview angle:** *"When would you use a union?"* "To save memory when only one of several representations is needed at a time, usually as a *tagged union*: a struct with a `kind` field plus a union, the C version of a variant. Also for low-level byte views (Q143)."

---

## Q142. [M8 | Struct vs union | Medium | MSQ | [C] · C11 | GATE+OA] · Type T3

Select **all** true statements (C11).

- A. All members of a union start at offset 0.
- B. The size of a union is at least the size of its largest member.
- C. The members of a struct may overlap in memory.
- D. Storing a value into a union member and then reading **the same** member is always well-defined.
- E. A union holds valid values for all of its members at the same time.

**Answer:** A, B, D

**In simple words:** A union is one storage slot with several **labels**. Whatever you last stored under one label is what's in the slot. Reading it back under the same label is always fine. Reading it under a *different* label means reinterpreting the same bytes (Q143–Q144).

**Why each option:**
- **A ✓:** N1570 §6.7.2.1p16: a pointer to a union, converted, points to each of its members, so all members begin at the start of the union.
- **B ✓:** It must hold any one member, plus possibly trailing padding (Q141).
- **C ✗:** Struct members are laid out in order at increasing addresses, and never overlap (§6.7.2.1p15).
- **D ✓:** That's the normal, intended use of a union.
- **E ✗:** Only the last-stored member has a meaningful value. Storing to one member overwrites the shared bytes.

**Concept recap:**

| | struct | union |
|---|---|---|
| memory | sum of members + padding | max member + padding |
| members valid at once | all | one (the last stored) |
| offsets | increasing | all 0 |
| typical use | records | variants, byte views |

**Trap alert:** Writing `u.i` and then reading `u.d` doesn't "convert" anything. The bytes are reinterpreted (see Q144).

**Interview angle:** *"Implement a variant type in C."* `struct Value { enum { INT, DBL } kind; union { int i; double d; }; };`. Always check `kind` before reading. That's exactly what C++17's `std::variant` automates (M22).

---

## Q143. [M8 | Seeing byte order through a union | Medium | MCQ | [C] · C11 | All] · Type T1 · P34

Platform: **little-endian**.

```c
#include <stdio.h>
int main(void) {
    union { unsigned int i; unsigned char b[4]; } u;
    u.i = 0x11223344;
    printf("%x %x\n", u.b[0], u.b[3]);
    return 0;
}
```

What does this program print?

- A. `11 44`
- B. `44 44`
- C. `33 22`
- D. `44 11`

**Answer:** D

**Behaviour class:** Well-defined in C. The byte order is **implementation-defined**, and fixed here by the little-endian assumption. In C, reading a union member other than the one last written reinterprets the bytes (N1570 §6.5.2.3, footnote 95), and reading bytes through `unsigned char` is always allowed.

**In simple words:** The union lets you look at the same 4 bytes two ways: as one number, or as 4 separate bytes. On a little-endian machine the number is stored **"little end first"**: the lowest byte (0x44) goes in the first slot and the highest (0x11) in the last, like writing a number's digits backwards.

**Why:**

| b[0] | b[1] | b[2] | b[3] |
|---|---|---|---|
| **44** | 33 | 22 | **11** |

Prints `44 11`.

**Why the other options fail:**
- A: The **big-endian** layout (11 22 33 44).
- B, C: Mix up the index order.

**Concept recap:** A union of an integer and an `unsigned char` array is the classic C way to inspect bytes. A `memcpy` into a byte array gives the same result.

**Trap alert:** **In C++, reading an inactive union member is UB** ([class.union]). C++ code must use `memcpy`, or `std::bit_cast` (C++20). The same source can be fine in C and UB in C++.

**Interview angle:** *"Detect endianness with a union."* This code. Then mention the C++-safe alternatives, and `htonl`/`ntohl` for network byte order.

---

## Q144. [M8 | Float bits through a union | Hard | MCQ | [C] · C11 | OA+Interview] · Type T3 · P34

Assume IEEE-754 single precision, and `sizeof(float) == sizeof(unsigned int) == 4`.

```c
#include <stdio.h>
int main(void) {
    union { float f; unsigned int u; } x;
    x.f = -2.0f;
    printf("%08x\n", x.u);
    return 0;
}
```

Which statement is correct?

- A. It prints `c0000000`. In C, reading `u` after writing `f` reinterprets the bytes (type-punning through a union is allowed). The same code in **C++** would be undefined behaviour.
- B. It prints `40000000`.
- C. The behaviour is undefined in C, because a different member was read.
- D. It prints `-2`.

**Answer:** A

**Behaviour class:** C: well-defined, given the stated IEEE-754 and size assumptions. C++: **UB**.

**In simple words:** A `float` is stored as three packed fields: a **sign bit**, an **exponent** and a **fraction**. Reading the same 4 bytes as an `unsigned int` shows the raw bit pattern rather than the number. For −2.0: sign = 1 (negative), exponent = 128 (since 2 = 1.0 × 2¹ and the bias is 127), fraction = 0.

**Why:**

| Field | Bits | Value for −2.0f |
|---|---|---|
| sign | 1 | 1 |
| exponent (bias 127) | 8 | 1 + 127 = 128 = `1000 0000` |
| fraction | 23 | 0 |

Bits: `1 10000000 000…0` = `1100 0000 0000 …` = **0xC0000000**.

**Why the other options fail:**
- B: That's **+2.0f** (sign bit 0).
- C: C explicitly allows reading a different member: the bytes are reinterpreted (§6.5.2.3, fn 95). It's C++ that makes this UB.
- D: That would be a value *conversion*, `(int)x.f`, not a reinterpretation of the bytes.

**Concept recap:** Ways to type-pun (reinterpret an object's bytes as another type):

| Method | C | C++ |
|---|---|---|
| union | ✓ | UB |
| pointer cast | UB (strict aliasing, Q117) | UB |
| `memcpy` | ✓ | ✓ |
| `std::bit_cast` | — | ✓ (C++20) |

**Trap alert:** Interviewers often show the union version, then ask "is this valid C++?". The answer is no.

**Interview angle:** *"How does the fast inverse square root trick read a float's bits?"* "The original used a pointer cast, which is UB by today's rules. Correct versions use `memcpy` or `std::bit_cast`, and compilers turn those into a single register move."

---

## Q145. [M8 | Bit-field values wrap | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P35

```c
#include <stdio.h>
struct F { unsigned int a : 3; unsigned int b : 5; };
int main(void) {
    struct F f = {9, 33};
    f.a += 7;
    printf("%u %u\n", f.a, f.b);
    return 0;
}
```

What does this program print?

- A. `16 33`
- B. `0 1`
- C. `1 1`
- D. `7 31`

**Answer:** B

**Behaviour class:** Well-defined. Storing a value into an **unsigned** bit-field reduces it modulo 2^width, just like any unsigned conversion. (Compilers may warn about the overflow in the initializer.)

**In simple words:** `a : 3` means "give `a` only 3 bits", so it can only count from 0 to 7, like a 3-digit binary odometer. Anything bigger wraps around: 9 becomes 1, and 33 in a 5-bit field (0–31) also wraps to 1.

**Why:**

| Field | Width | Range | Stored value |
|---|---|---|---|
| a = 9 | 3 | 0–7 | 9 mod 8 = **1** |
| b = 33 | 5 | 0–31 | 33 mod 32 = **1** |
| `f.a += 7` | | | (1 + 7) = 8 mod 8 = **0** |

Prints `0 1`.

**Why the other options fail:**
- A: Ignores the bit widths.
- C: Forgets the `+= 7` line.
- D: Clamps to the maximum value instead of wrapping.

**Concept recap:** A bit-field `type name : width;` packs several small fields into one machine word. Typical uses: hardware registers, protocol flags, compact records.

**Trap alert:** For a **signed** bit-field, an out-of-range value gives an implementation-defined result (it isn't wrapped by definition). Whether a plain `int` bit-field is signed at all is also implementation-defined (Q146). Use `unsigned` for predictable bit-fields.

**Interview angle:** *"Bit-fields vs manual masks?"* "Bit-fields are convenient, but the layout — bit order, how fields are packed into units, padding — is implementation-defined, so they're not portable across compilers or architectures. For wire formats and registers that must match exactly, use explicit shifts and masks."

---

## Q146. [M8 | Rules for bit-fields | Hard | MSQ | [C] · C11 | All] · Type T2 · P35

Given `struct G { int x : 3; unsigned y : 4; } g;`, select **all** true statements.

- A. `&g.y` is ill-formed: you can't take the address of a bit-field.
- B. Whether the plain `int x : 3` field is signed or unsigned is implementation-defined.
- C. On typical 64-bit ABIs, `sizeof(struct G)` is 4.
- D. After `g.x = 5;`, the standard guarantees that `g.x == 5`.
- E. `sizeof(g.y)` evaluates to 4.

**Answer:** A, B, C

**In simple words:** A bit-field may be only a few bits inside a byte, and an address can only point at a whole byte. So a bit-field has **no address** and **no `sizeof`**. Several small fields share one storage unit (here a single 4-byte `int`).

**Why each option:**
- **A ✓:** The operand of unary `&` can't be a bit-field (a constraint, N1570 §6.5.3.2p1). Compilers report "cannot take address of bit-field".
- **B ✓:** For bit-fields specifically, whether plain `int` means `signed int` or `unsigned int` is implementation-defined (§6.7.2p5). (For ordinary variables, `int` is always signed.)
- **C ✓:** Both fields (3 + 4 = 7 bits) fit into one `int`-sized unit → 4 bytes. (Layout is implementation-defined, but this is universal on mainstream 64-bit ABIs.)
- **D ✗:** If `x` is signed, 3 bits hold −4 … 3, so 5 doesn't fit and the stored value is implementation-defined (typically −3). If `x` happens to be unsigned, 5 fits. Either way, the standard doesn't guarantee 5 in general.
- **E ✗:** `sizeof` can't be applied to a bit-field (a constraint, §6.5.3.4p1), so this is a compile error.

**Concept recap:** Bit-fields can't be used with `&`, `sizeof`, arrays of bit-fields, or pointers to bit-fields. They can be used in expressions; they are promoted like small integers.

**Trap alert:** `scanf("%d", &g.y)` is ill-formed for the same reason as option A. Read into an ordinary variable, then assign it to the field.

**Interview angle:** *"Name some things you can't do with a bit-field."* Take its address, apply `sizeof`, make an array of them, or rely on a portable layout.

---

## Q147. [M8 | Self-referential struct (a linked list) | Easy | MCQ | [C] · C11 | GATE] · Type T1 · P36

```c
#include <stdio.h>
struct Node { int v; struct Node *next; };
int main(void) {
    struct Node c = {3, NULL}, b = {2, &c}, a = {1, &b};
    int sum = 0;
    for (struct Node *p = &a; p; p = p->next)
        sum = sum * 10 + p->v;
    printf("%d\n", sum);
    return 0;
}
```

What value is printed?

- A. `6`
- B. `321`
- C. `123`
- D. `12`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Each node is a box holding a number and an **arrow to the next box**. The loop follows the arrows from `a` until an arrow points nowhere (NULL). `sum * 10 + v` builds a number digit by digit, in the order the boxes are visited.

**Why:**

| Visit | p → | v | sum = sum × 10 + v |
|---|---|---|---|
| 1 | a | 1 | 1 |
| 2 | b | 2 | 12 |
| 3 | c | 3 | 123 |
| — | NULL | | stop |

**Why the other options fail:**
- A: Adds the digits.
- B: Visits the nodes in reverse.
- D: Stops one node early.

**Concept recap:** A struct may contain a **pointer** to its own type, because a pointer has a known size even while the struct is still being defined. That's what makes linked lists, trees and graphs possible.

**Trap alert:** `struct Node { int v; struct Node next; };` (no `*`) is **ill-formed**. The struct would have to contain a complete copy of itself, which would make it infinitely large.

**Interview angle:** *"Why must `next` be a pointer?"* "While the struct is being defined, its type is incomplete, and an incomplete type can't be used as a member. A pointer to an incomplete type is fine, because every pointer has the same known size."

---

## Q148. [M8 | Anonymous struct/union members (C11) | Medium | MCQ | [C] · C11 | OA+Interview] · Type T1

```c
#include <stdio.h>
struct Shape {
    int kind;
    union {
        struct { int w, h; };
        int r;
    };
};
int main(void) {
    struct Shape s = { .kind = 1, .w = 3, .h = 4 };
    printf("%d %d %zu\n", s.w * s.h, s.r, sizeof s);
    return 0;
}
```

What does this program print?

- A. `12 3 12`
- B. `12 0 16`
- C. Compile error: `w` is not a member of `struct Shape`
- D. `12 4 12`

**Answer:** A

**Behaviour class:** Well-defined. `r` and `w` are both `int` at offset 0 of the union, so reading `r` reads the same `int` that `w` stored.

**In simple words:** An *anonymous* union or struct is a group with **no name of its own**. Its members can be used as if they belonged directly to the outer struct, so you write `s.w` instead of `s.u.dims.w`. Here `w`/`h` and `r` share the same space: a shape is either a rectangle (w, h) or a circle (r).

**Why:**

| Member | Offset | Value |
|---|---|---|
| kind | 0 | 1 |
| union: { w, h } / r | 4 | w = 3, h = 4; r shares w's bytes → 3 |

- `s.w * s.h` = 12.
- `s.r` = 3 (the same bytes as `w`).
- sizeof: 4 (kind) + 8 (the union: max(8, 4)) = **12**, alignment 4.

**Why the other options fail:**
- B: Assumes `r` is a separate zero-initialised member, and miscounts the size.
- C: Anonymous members have been standard since **C11** (§6.7.2.1p13).
- D: `r` overlaps `w`, not `h`.

**Concept recap:** Anonymous unions are the idiomatic way to write tagged variants: check `kind` first, then read the matching members.

**Trap alert:** Anonymous unions also exist in C++, but anonymous **structs** are only a common compiler extension there, not standard C++.

**Interview angle:** *"How do you represent 'a shape is a circle or a rectangle' in C vs C++?"* "C: a tag plus an (anonymous) union, checked by hand. C++: `std::variant<Circle, Rect>` with `std::visit`, or a class hierarchy with virtual functions (M16)."

---

## Q149. [M8 | Flexible array member | Hard | MCQ | [C] · C11 | OA+Interview] · Type T1

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
struct Msg { int len; char data[]; };
int main(void) {
    struct Msg *m = malloc(sizeof *m + 6);
    if (!m) return 1;
    m->len = 5;
    memcpy(m->data, "hello", 6);
    printf("%zu %d %s\n", sizeof(struct Msg), m->len, m->data);
    free(m);
    return 0;
}
```

What does this program print?

- A. `12 5 hello`
- B. `8 5 hello`
- C. Compile error: an array member needs a size
- D. `4 5 hello`

**Answer:** D

**Behaviour class:** Well-defined. 6 extra bytes were allocated for `"hello"` plus its `'\0'`.

**In simple words:** `char data[];` at the **end** of a struct is a "this space extends as far as you allocate" marker. The struct header itself is only the `int` (4 bytes), and `data` takes no space in `sizeof`. You decide how much room `data` gets when you `malloc`, so one allocation holds both the header and the payload.

**Why:**

| Item | Value |
|---|---|
| `sizeof(struct Msg)` | 4: only `len`; the flexible member contributes 0 (and needs no padding here) |
| allocation | 4 + 6 = 10 bytes |
| `m->data` | starts right after `len`, with 6 usable bytes |
| output | `4 5 hello` |

**Why the other options fail:**
- A, B: Guess a size for `data`. It contributes nothing to `sizeof`.
- C: A **flexible array member** has been standard since C99 (N1570 §6.7.2.1p18). It must be the **last** member, and the struct must have at least one other named member.

**Concept recap:** Allocation idiom: `malloc(sizeof(struct Msg) + n * sizeof(m->data[0]))`, or the slightly tighter `offsetof(struct Msg, data) + n`. Uses: network packets, variable-length records, and string headers like `struct { size_t len; char s[]; }`.

**Trap alert:**
- `struct Msg x;` (a non-heap instance) gets **no** usable `data` storage, so writing `x.data[0]` is out-of-bounds UB.
- Copying `*m` by assignment copies only the header.
- FAMs aren't standard C++. Use `std::vector` or `std::string` there.

**Interview angle:** *"FAM vs a `char *data` member?"* "A FAM needs one allocation instead of two, keeps the payload next to its header (better locality), and is freed with a single `free`. The downsides: no struct copying, no arrays of such structs, and it must be the last member."

---

## Q150. [M9 | calloc vs malloc | Easy | MCQ | [C] · C11 | GATE+OA] · Type T5

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int main(void) {
    int *a = calloc(4, sizeof *a);
    int *b = malloc(4 * sizeof *b);
    for (int i = 0; i < 4; i++) b[i] = i;
    memcpy(a + 1, b, 2 * sizeof *b);
    printf("%d %d %d %d\n", a[0], a[1], a[2], a[3]);
    free(a); free(b);
    return 0;
}
```

What does this program print?

- A. `0 1 2 3`
- B. `0 0 1 0`
- C. `0 0 1 2`
- D. Some values are unpredictable, because `a` was never initialised

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** `calloc` gives you memory that has been **wiped clean** (all zeros). `malloc` gives you memory **as it was left**: it may contain anything until you write to it. Here `b` is fully written before use, and `a` starts as all zeros, so every value is known.

**Why:**

| Line | Action | a | b |
|---|---|---|---|
| 5 | `calloc(4, 4)` → zeroed | {0, 0, 0, 0} | |
| 6–7 | `malloc`, then b[i] = i | | {0, 1, 2, 3} |
| 8 | copy 2 ints from b into a starting at index 1 | {0, **0**, **1**, 0} | |
| 9 | print | `0 0 1 0` | |

**Why the other options fail:**
- A: Copies all of `b` into `a` from the start.
- C: Copies 3 ints.
- D: `calloc` guarantees all-bits-zero, which for `int` means 0.

**Concept recap:**

| | `malloc(n)` | `calloc(count, size)` |
|---|---|---|
| contents | indeterminate | all bits zero |
| arguments | total bytes | count × size, checked for overflow by the library |
| speed | slightly faster | may be just as fast for big blocks (fresh OS pages are already zero) |

**Trap alert:** All-bits-zero is guaranteed to mean 0 for integers, but **not** guaranteed to be a null pointer or 0.0 on every platform (Q85-E). It is in practice on mainstream systems.

**Interview angle:** *"Why does calloc take two arguments?"* "So the library can check `count × size` for overflow. `malloc(count * size)` can silently wrap to a small allocation, which is a known security bug class."

---

## Q151. [M9 | sizeof mistake inside malloc | Medium | MCQ | [C] · C11 | All] · Type T4

```c
#include <stdlib.h>
struct Point { double x, y; };
struct Point *make(int n) {
    struct Point *p = malloc(n * sizeof(p));
    for (int i = 0; i < n; i++) { p[i].x = i; p[i].y = -i; }
    return p;
}
```

What is the **smallest correct fix**?

- A. Change `sizeof(p)` to `sizeof(struct Point *)`.
- B. Cast the result: `(struct Point *)malloc(n * sizeof(p))`.
- C. Change `sizeof(p)` to `sizeof *p`.
- D. Use `calloc(n, sizeof(p))` instead.

**Answer:** C

**Behaviour class:** The original is **UB** for any n ≥ 1: it writes past the end of the allocation (heap overflow).

**In simple words:** `p` is a **pointer**, and `sizeof(p)` is the size of the pointer (8 bytes), not the size of the thing it points to (a `struct Point`, 16 bytes). So the code asks for room for half the data. `sizeof *p` means "the size of what p points to", which is always right, even if you later change p's type.

**Why:**

| Expression | Value | Bytes for n = 10 |
|---|---|---|
| `sizeof(p)` | 8 (a pointer) | 80, **too small** |
| `sizeof *p` | 16 (`struct Point`) | 160 ✓ |

The loop writes 16 × n bytes into an 8 × n block, so half of the writes go out of bounds.

**Why the other options fail:**
- A: Still the size of a pointer.
- B: Casting doesn't change how many bytes were allocated.
- D: `calloc` zeroes the memory but is still given the wrong size.

**Concept recap:** The idiom `T *p = malloc(n * sizeof *p);` repeats no type names, so it can't go out of sync with the declaration.

**Trap alert:** With `int *p` on a 64-bit system, `sizeof(p)` = 8 is *bigger* than `sizeof(int)` = 4. The same bug then over-allocates and "works", hiding the mistake until the type changes.

**Interview angle:** *"What's wrong with `malloc(n * sizeof(p))`?"* "It allocates n pointers' worth of space instead of n elements' worth. Use `sizeof *p`. A static analyser or ASan catches the resulting overflow."

---

## Q152. [M9 | Allocation failure, free(NULL), realloc | Easy | MSQ | [C] · C11 | All] · Type T3

Select **all** true statements (C11).

- A. `malloc` returns a null pointer when it can't allocate.
- B. `malloc` throws an exception when it can't allocate.
- C. Writing through `malloc`'s result without checking it is undefined behaviour when the allocation failed.
- D. `p = realloc(p, n);` leaks the original block if `realloc` fails.
- E. `free(NULL)` is undefined behaviour.

**Answer:** A, C, D

**In simple words:** When the system is out of memory, `malloc` just hands back **"nothing"** (NULL). It doesn't raise an alarm. Using that "nothing" as if it were memory is like moving into an address that doesn't exist. `realloc` returns NULL on failure but keeps the old block, so if you overwrite your only pointer with that NULL, you've lost the key to memory you still own.

**Why each option:**
- **A ✓:** N1570 §7.22.3.4p3.
- **B ✗:** C has no exceptions. (C++'s `new` throws `std::bad_alloc`.)
- **C ✓:** Dereferencing a null pointer is UB.
- **D ✓:** On failure, `realloc` returns NULL and the **old block is left untouched** (§7.22.3.5p4). Assigning the result straight back to `p` loses the only pointer to it, which is a leak. Correct pattern: `tmp = realloc(p, n); if (!tmp) { /* handle */ } else p = tmp;`
- **E ✗:** `free(NULL)` does nothing, guaranteed (§7.22.3.3p2). So `free(p); p = NULL;` is safe to repeat.

**Concept recap:** Check every allocation, and decide on a policy: return an error code, print and exit, or use a wrapper (`xmalloc`) that aborts.

**Trap alert:** On Linux, `malloc` rarely returns NULL even when memory is short (because of overcommit). The failure can show up later as the process being killed. That's no excuse to skip the check: portable code must handle NULL.

**Interview angle:** *"What does `new` do on failure, and can you make it return null?"* "It throws `std::bad_alloc`. `new (std::nothrow) T` returns `nullptr` instead."

---

## Q153. [M9 | Returning heap memory | Easy | MCQ | [C] · C11 | GATE+OA] · Type T6

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
char *dup(const char *s) {
    char *d = malloc(strlen(s) + 1);
    if (!d) return NULL;
    strcpy(d, s);
    return d;
}
int main(void) {
    char *a = dup("heap");
    a[0] = 'H';
    printf("%s %zu\n", a, strlen(a));
    free(a);
    return 0;
}
```

What does this program print?

- A. `Heap 4`
- B. `heap 4`
- C. `Heap 5`
- D. Undefined behaviour: `a` points into `dup`'s stack frame

**Answer:** A

**Behaviour class:** Well-defined. `malloc`'d memory lives until it is `free`d, regardless of which function allocated it.

**In simple words:** Local variables are **rented for the length of the function call**; heap memory is **rented until you explicitly give it back** with `free`. So `dup` can safely hand its heap copy to `main`, and `main` becomes responsible for returning it.

**Why:**

| Step | Action | Result |
|---|---|---|
| `malloc(4 + 1)` | 5 bytes: room for "heap" + `'\0'` | |
| `strcpy` | copies "heap\0" | |
| return | the pointer stays valid after `dup` returns | |
| `a[0] = 'H'` | writable: it's our own heap copy, not the literal | `"Heap"` |
| print | | `Heap 4` |

**Why the other options fail:**
- B: Ignores the write.
- C: Counts the `'\0'`.
- D: Confuses heap memory with a local array (Q122).

**Concept recap:** **Ownership:** whoever receives a pointer from an allocating function must `free` it, exactly once. Document this in the function's contract ("the caller must free the result").

**Trap alert:** `malloc(strlen(s))` without the `+ 1` is the most common off-by-one in C (Q160). Also note that `strdup` became standard only in **C23**; before that it was POSIX.

**Interview angle:** *"Who frees memory returned by a library function?"* "Whatever the API documents. In C that's usually the caller (`strdup`, `getline`). C++ avoids the question by returning a `std::string` or a `std::unique_ptr`, which own their memory and free it automatically."

---

## Q154. [M9 | Use-after-free through a second pointer | Medium | MCQ | [C] · C11 | All] · Type T3 · P45

> **This code has undefined behaviour.**

```c
#include <stdio.h>
#include <stdlib.h>
int main(void) {
    int *p = malloc(2 * sizeof *p);
    int *q = p;
    p[0] = 1; p[1] = 2;
    free(p);
    p = NULL;
    printf("%d\n", q[1]);
    return 0;
}
```

Which statement is correct?

- A. It prints `2`: `q` still points to the data.
- B. It prints `0`: `free` clears the memory.
- C. It crashes reliably, because setting `p = NULL` also nulls `q`.
- D. The behaviour is undefined: `q[1]` reads freed memory (use-after-free). Setting `p = NULL` doesn't protect the second pointer `q`.

**Answer:** D

**Behaviour class:** **Undefined behaviour.**

**In simple words:** `free` is like handing back the keys to a hotel room. `p` and `q` were two copies of the room key. Throwing away one copy (`p = NULL`) doesn't change the fact that the room has been returned. Using the other copy (`q`) means walking into a room that may already belong to someone else.

> **⚠ UB box**
> - **Category:** use-after-free (accessing an object after its storage was deallocated).
> - **Rule:** N1570 §7.22.3.3p2 (free) and §6.2.4p2: after `free`, the object's lifetime has ended, and every pointer to it — including copies — has an indeterminate value.
> - **What you may see in practice:** often `2` (the bytes haven't been reused yet) or some junk value (the allocator wrote its bookkeeping data there). It rarely crashes immediately. That's what makes use-after-free bugs so dangerous and so widely exploited. **Not guaranteed.**

**Why:**

| Step | p | q | Block |
|---|---|---|---|
| malloc | → block | → block | live |
| `free(p)` | indeterminate | **indeterminate** | freed |
| `p = NULL` | NULL | still the old address | freed |
| `q[1]` | | | **UB** |

**Why the other options fail:**
- A, B: Each is one possible observation, presented as a guarantee.
- C: `p` and `q` are separate variables. Changing `p` doesn't touch `q`.

**Concept recap:** `p = NULL` after `free` is good hygiene, but it protects only **that one pointer**. Aliases need a clear **ownership** rule (exactly one owner frees, and the others don't outlive it).

**Trap alert:** Use-after-free is behind a large share of real-world security vulnerabilities in browsers and OS kernels. The fix is ownership discipline, which C++ encodes in `unique_ptr` (M22).

**Interview angle:** *"How do you find use-after-free bugs?"* "AddressSanitizer (it quarantines freed memory and reports the access), Valgrind memcheck, and hardened allocators. Design-level fix: single ownership, and RAII in C++."

---

## Q155. [M9 | Which free sequence is safe? | Medium | MCQ | [C] · C11 | All] · Type T3 · P45

Assume `int *p = malloc(4 * sizeof *p), *q;` succeeded before each option. Which sequence is **well-defined**?

- A. `free(p); free(p);`
- B. `q = p; free(p); free(q);`
- C. `free(p); p = NULL; free(p);`
- D. `free(p + 1);`

**Answer:** C

**Behaviour class:** C is well-defined. A, B and D are **UB**.

**In simple words:** Each block from `malloc` must be returned **exactly once**, using **exactly the address `malloc` gave you**. Returning it twice (even through a different variable), or returning "the middle of it", confuses the allocator's records. `free(NULL)` is a safe "return nothing".

**Why:**

| Option | What happens | Class |
|---|---|---|
| A | the same block freed twice → **double free** | UB |
| B | `q` is the same address, so freeing both is a double free in disguise | UB |
| C | the second call is `free(NULL)`, which does nothing | ✓ well-defined |
| D | the address wasn't returned by `malloc`/`calloc`/`realloc` | UB |

(N1570 §7.22.3.3p2: UB if the argument doesn't match a pointer returned by an allocation function, or if the space was already freed.)

**Concept recap:** A double free can corrupt the allocator's internal lists, which can crash later at an unrelated `malloc` and is exploitable.

**Trap alert:** Option B is how double frees usually happen in real code: two structs both think they own the same buffer (a shallow copy, M9 part 2).

**Interview angle:** *"Why set pointers to NULL after free?"* "It turns an accidental second `free` into a harmless no-op, and a later dereference into an obvious null-pointer crash instead of silent corruption. But it doesn't fix aliases (Q154)."

---

## Q156. [M9 | realloc invalidates old pointers | Hard | MCQ | [C] · C11 | All] · Type T5 · P45

```c
#include <stdio.h>
#include <stdlib.h>
int main(void) {
    int *a = malloc(2 * sizeof *a);
    if (!a) return 1;
    a[0] = 1; a[1] = 2;
    int *first = &a[0];
    int *tmp = realloc(a, 1000 * sizeof *a);
    if (!tmp) { free(a); return 1; }
    a = tmp;
    printf("%d\n", *first);
    free(a);
    return 0;
}
```

Which statement is correct?

- A. It prints `1` if `realloc` grew the block in place, and is UB only if the block moved.
- B. The behaviour is undefined whenever `realloc` succeeds: the old object is **always** deallocated, even if the new block has the same address, so `first` is dangling.
- C. It prints `1`, because `realloc` copies the old contents.
- D. It prints `0`, because `realloc` zero-fills the block.

**Answer:** B

**Behaviour class:** **Undefined behaviour.**

**In simple words:** `realloc` doesn't "stretch" your block. Officially it **gives you a new block and throws away the old one**, copying the contents across. Sometimes the new block happens to be at the same address, but the standard still treats it as a *different* object. So any old pointer into the old block (`first`) is out of date. Always re-derive pointers from the new one, e.g. `&a[0]`.

> **⚠ UB box**
> - **Category:** use of a pointer to a deallocated object (realloc flavour).
> - **Rule:** N1570 §7.22.3.5p2: realloc **deallocates the old object** and returns a pointer to a new object, "which may have the same value as a pointer to the old object". Using the old pointer is therefore use-after-free.
> - **What you may see in practice:** if the block grew in place, often `1`; if it moved (very likely when growing to 4000 bytes), often a junk value. **Not guaranteed in either case.**

**Why the other options fail:**
- A: Even when the address is the same, the old object is formally gone. Some compilers even use this rule to optimise.
- C: The contents *are* copied into the new block (`a[0]` is 1), but `first` isn't the new block.
- D: `realloc` doesn't zero the extra space; it's indeterminate.

**Concept recap:** After `realloc`:
- recompute every pointer into the block from the new base pointer;
- or store **indices** (offsets) instead of pointers, which survive a move.

**Trap alert:** This is exactly C++'s **iterator invalidation** for `std::vector` (P87, M21). `push_back` past the capacity reallocates, and every pointer, reference and iterator into the vector becomes invalid.

**Interview angle:** *"Why does my pointer into a growing dynamic array go bad?"* "Growth reallocates and moves the elements. Keep indices, re-fetch pointers after each growth, or reserve enough capacity up front."

---

## Q157. [M9 | Counting leaked bytes | Medium | MCQ | [C] · C11 | GATE+OA] · Type T5 · P46

```c
#include <stdlib.h>
void f(int n) {
    int *buf = malloc(n * sizeof *buf);
    if (!buf) return;
    if (n > 5) return;          /* early exit */
    for (int i = 0; i < n; i++) buf[i] = i;
    free(buf);
}
int main(void) {
    f(3); f(8); f(10);
    return 0;
}
```

How many bytes are leaked (`sizeof(int)` = 4)?

- A. 72
- B. 0
- C. 84
- D. 40

**Answer:** A

**Behaviour class:** Well-defined. A leak isn't UB; it's a resource bug.

**In simple words:** A leak is memory you borrowed and never gave back, because you **lost track of it**. Here the function borrows a buffer, then for large `n` it leaves early through a side door (`return`) without reaching the `free`. Once the function returns, nothing points to that buffer any more, so it can never be freed.

**Why:**

| Call | Bytes allocated | Takes the early exit? | Freed? | Leaked |
|---|---|---|---|---|
| f(3) | 12 | no | yes | 0 |
| f(8) | 32 | yes | **no** | 32 |
| f(10) | 40 | yes | **no** | 40 |
| **total** | 84 | | | **72** |

**Why the other options fail:**
- B: Misses the early `return`.
- C: Counts every allocation as leaked.
- D: Counts only the last call.

**Concept recap:** Every exit path from a function must release what was acquired. The C idiom is a single cleanup label: `if (…) goto out; … out: free(buf); return;`. C++ makes this automatic with RAII: `std::vector<int> buf(n);` is freed on every path.

**Trap alert:** When the process exits, the OS reclaims all its memory, so a leak in a short program "doesn't matter". In a long-running server, a few bytes per request eventually exhaust memory.

**Interview angle:** *"How do you find memory leaks?"* "Valgrind (`--leak-check=full`), LeakSanitizer (part of ASan), heap profilers. And prevent them by design: RAII and smart pointers in C++, and single-exit cleanup in C."

---

## Q158. [M9 | Which snippets leak? | Hard | MSQ | [C] · C11 | All] · Type T5 · P46

Treat each snippet independently, as complete code inside some function. Assume every allocation succeeds and every pointer variable is declared. Select **all** snippets that **leak** memory.

- A. `p = malloc(10); p = malloc(20); free(p);`
- B. `p = malloc(10); q = p; free(q);`
- C. `p = malloc(10); p = realloc(p, 20); free(p);`
- D. `s = malloc(8); s++; free(s - 1);`
- E. `for (int i = 0; i < 3; i++) { int *t = malloc(4); }`

**Answer:** A, E

**In simple words:** For each snippet, ask: **"Is there still a pointer to every block when it's time to free it?"** Overwriting your only pointer (A), or letting it disappear at the end of a loop body (E), leaves the block unreachable. That's a leak.

**Why each option:**
- **A ✓ leak:** The first block's only pointer is overwritten → 10 bytes lost.
- **B ✗:** `q` is a copy of `p`, so `free(q)` frees the block. One block, one free.
- **C ✗:** `realloc` succeeds (as assumed), so the old block is released by `realloc` itself and the new one is freed. (If `realloc` could fail, this exact line would leak; Q152-D.)
- **D ✗:** `s - 1` is the original address, so the free is correct.
- **E ✓ leak:** Each iteration's `t` goes out of scope with its block never freed → 3 × 4 = 12 bytes.

**Concept recap:** A leak = an allocated block with **no remaining pointer to it**. Tools find leaks by scanning memory for pointers to every live block, just like a garbage collector's reachability check.

**Trap alert:** D looks scary, but it's fine. What matters is that `free` receives exactly the address `malloc` returned. `free(s)` (after `s++`) would be UB.

**Interview angle:** *"Can a garbage-collected language leak memory?"* "Yes, logically: you can keep references you'll never use again (growing caches, listeners that are never removed). GC only frees *unreachable* objects."

---

## Q159. [M9 | Freeing a linked list | Hard | MCQ | [C] · C11 | All] · Type T4 · P27, P45

```c
struct Node { int v; struct Node *next; };
void free_list(struct Node *head) {
    for (struct Node *p = head; p; p = p->next)
        free(p);
}
```

What is wrong, and what is the smallest correct fix?

- A. Nothing is wrong.
- B. It leaks the last node. Fix: add `free(head);` after the loop.
- C. It double-frees the head. Fix: start the loop from `head->next`.
- D. The update `p = p->next` reads `p->next` **after** `free(p)`, which is a use-after-free. Fix: save the next pointer first: `for (struct Node *p = head, *nx; p; p = nx) { nx = p->next; free(p); }`

**Answer:** D

**Behaviour class:** The original is **UB** for any non-empty list.

**In simple words:** It's like demolishing a house and *then* going inside to read the address of the next house written on its wall. Once a node is freed, you can't read its `next`. **Read the next address first, then demolish.**

**Why:**

| Step in the loop | Action | Problem |
|---|---|---|
| body | `free(p)` | node p is deallocated |
| update | `p = p->next` | reads a field of the **freed** node → UB |

Fixed version: save `nx = p->next` while `p` is still valid, then `free(p)`, then move to `nx`.

**Why the other options fail:**
- A: It often *appears* to work, because freed memory usually isn't overwritten immediately. It's still UB.
- B: Every node is freed. The problem is reading, not missing nodes.
- C: Nothing is freed twice.

**Concept recap:** Never touch an object after `free`, including reading its fields as part of a loop update. The same pattern applies to destroying trees (free the children before the parent) and to erasing from containers while iterating (C++ `erase` returns the next iterator).

**Trap alert:** A recursive free, `void f(struct Node *n) { if (n) { f(n->next); free(n); } }`, is correct, but uses O(n) stack space and can overflow on long lists.

**Interview angle:** *"Delete a linked list."* Write the save-then-free loop and mention the recursive alternative's stack cost. In C++, `std::unique_ptr<Node> next` frees the list automatically, but very long lists can still overflow the stack through the recursive destructors.

---

## Q160. [M9 | Forgetting +1 for the '\0' | Hard | MCQ | [C] · C11 | All] · Type T4 · P28

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int main(void) {
    const char *src = "hello";
    char *s = malloc(strlen(src));
    strcpy(s, src);
    printf("%s\n", s);
    free(s);
    return 0;
}
```

What is the smallest correct fix?

- A. Replace `strcpy` with `strncpy(s, src, strlen(src));`
- B. Cast the result of `malloc` to `char *`.
- C. Allocate `strlen(src) + 1` bytes.
- D. Add `s[5] = '\0';` after the `strcpy`.

**Answer:** C

**Behaviour class:** The original is **UB**: `strcpy` writes 6 bytes into a 5-byte heap block (heap buffer overflow).

**In simple words:** `strlen` counts the **letters only**, but storing a string also needs room for the invisible end marker `'\0'`. Allocating `strlen` bytes is like booking 5 seats for 5 people plus a driver.

**Why:**

| Item | Value |
|---|---|
| `strlen("hello")` | 5 |
| bytes needed | 5 + 1 = **6** |
| allocated | 5 → `strcpy` writes the 6th byte out of bounds |

**Why the other options fail:**
- A: Copies 5 bytes with no `'\0'`, so the following `printf("%s")` reads out of bounds. That's still UB.
- B: Casting doesn't change the size.
- D: `s[5]` is itself one past the 5-byte block, so this is another out-of-bounds write.

**Concept recap:** The rule for string buffers: **bytes = length + 1**. In C23, use `strdup(src)`, which does this correctly.

**Trap alert:** This bug often goes unnoticed because the allocator rounds small requests up (often to 16 or 24 bytes), so the extra byte lands in slack space. AddressSanitizer catches it immediately.

**Interview angle:** *"Find the bug in this string-copy code."* A very common OA question, and the answer is almost always the missing `+ 1`. Then mention `strdup`/`strndup`, or `std::string` in C++.

---

## Batch 8 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 141 | M8 | E | MCQ | C | Well-defined (impl-defined sizes) |
| 142 | M8 | M | MSQ | A, B, D | — (concepts) |
| 143 | M8 | M | MCQ | D | Well-defined in C (byte order impl-defined) |
| 144 | M8 | H | MCQ | A | Well-defined in C; UB in C++ |
| 145 | M8 | M | MCQ | B | Well-defined |
| 146 | M8 | H | MSQ | A, B, C | Ill-formed / impl-defined mix |
| 147 | M8 | E | MCQ | C | Well-defined |
| 148 | M8 | M | MCQ | A | Well-defined |
| 149 | M8 | H | MCQ | D | Well-defined |
| 150 | M9 | E | MCQ | B | Well-defined |
| 151 | M9 | M | MCQ | C | **UB** (original) |
| 152 | M9 | E | MSQ | A, C, D | — (concepts) |
| 153 | M9 | E | MCQ | A | Well-defined |
| 154 | M9 | M | MCQ | D | **UB** |
| 155 | M9 | M | MCQ | C | UB / OK mix |
| 156 | M9 | H | MCQ | B | **UB** |
| 157 | M9 | M | MCQ | A | Well-defined (leak) |
| 158 | M9 | H | MSQ | A, E | Leak analysis |
| 159 | M9 | H | MCQ | D | **UB** (original) |
| 160 | M9 | H | MCQ | C | **UB** (original) |

- **Difficulty:** 5 E / 8 M / 7 H.
  - M8 total (Q132–Q149): 4/8/6, as planned.
  - M9 so far: 3/4/4.
- **Format:** 16 MCQ / 4 MSQ.
  - M8 total: 15/3, as planned.
- **MCQ answer keys:** A×4, B×3, C×5, D×4. Running total Q1–Q160: A 32, B 32, C 33, D 32.

---

**Completed: 160 / 400 Remaining: 240**

**Ledger, newly ✓ this batch:**
- P34 ×2 (Q143, Q144) ★ complete
- P35 ×2 (Q145, Q146)
- P36 ×1 (Q147)
- P45 ×3 (Q154, Q155, Q156; plus Q159) ★ complete
- P46 ×2 (Q157, Q158)
- P27 → 3 (Q159)
- P28 → 3 (Q160)

**M8 complete** (Q132–Q149, 18 questions). **M9 in progress** (11 of 16 done).

*Say **"Continue."** for Batch 9 (Q161–Q180: the end of M9 — `realloc(NULL, n)` and size 0, shallow vs deep copy, dynamic 2D arrays, `malloc(0)` — then all of M10 preprocessor, compilation and linking, then the start of M11 stdio).*
