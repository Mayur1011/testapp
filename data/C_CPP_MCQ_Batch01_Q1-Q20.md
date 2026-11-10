# C/C++ MCQ/MSQ Bank — Batch 1 (Q1–Q20)

**What this batch covers:**
- **M1 (Q1–Q18):** type sizes and signedness; unsigned wrap-around; `sizeof` and the fact that its operand is not evaluated; `sizeof('a')` in C vs C++; octal, hex and escape literals; `static`/`extern` linkage across files; compile errors vs link errors; reading uninitialised variables; integer truncation; the signedness of plain `char`; what the standard guarantees about sizes; valid `main` signatures; tentative definitions; scope vs linkage with a block-scope `extern`; the `typedef` vs `#define` pointer trap; enum values; compilation stages.
- **M2 (Q19–Q20):** operator precedence and associativity.

**Patterns:** P01, P02 (secondary), P03, P05, P10, P14, P42, P48, P50, P97, P98

> **Conventions for every question in this batch (unless the question says otherwise):** C11, compiled with `gcc -std=c11`. 64-bit Linux / GCC / x86-64. `sizeof(short)=2`, `sizeof(int)=4`, `sizeof(long)=8`, `sizeof(void*)=8`. Plain `char` is signed. ASCII character set. Two's-complement integers. `NDEBUG` is not defined. The headers shown are the only ones included. Optimisation level does not matter unless stated.

---

## Q1. [M1 | Data-type sizes | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P01

Platform: the default above (LP64: `int` = 4 bytes, `long` = 8 bytes).

```c
#include <stdio.h>
int main(void) {
    printf("%zu %zu %zu %zu\n",
           sizeof(char), sizeof(short), sizeof(int), sizeof(long));
    return 0;
}
```

What does this program print?

- A. `1 2 4 8`
- B. `1 2 4 4`
- C. `1 2 2 4`
- D. `2 2 4 8`

**Answer:** A

**Behaviour class:** Well-defined. The program is well-defined, but the numbers it prints are **implementation-defined**, except `sizeof(char)`, which is always 1.

**Why:**

| Line | Statement executed | Values |
|---|---|---|
| 3–4 | evaluate the four `sizeof`s | `sizeof(char)=1` (always, by definition), `sizeof(short)=2`, `sizeof(int)=4`, `sizeof(long)=8` (LP64 data model) |
| 3 | `printf` | prints `1 2 4 8` |

`sizeof` gives a `size_t`, so the correct format specifier is `%zu`. Using `%d` for a `size_t` is a type mismatch, and that is UB.

**Why the other options fail:**
- B: `long` is 4 bytes under **LLP64**, the 64-bit Windows/MSVC model. It is 8 bytes on 64-bit Linux (**LP64**).
- C: `int` was 2 bytes on 16-bit DOS compilers such as Turbo C. Some old GATE books still assume this. Always check the stated platform.
- D: `sizeof(char)` is 1 by definition (N1570 §6.5.3.4p4), on every platform.

**Concept recap:** The common data models are:
- **ILP32** (32-bit): `int`, `long` and pointers are all 4 bytes.
- **LP64** (64-bit Linux/macOS): `long` and pointers are 8 bytes; `int` is 4.
- **LLP64** (64-bit Windows): only `long long` and pointers are 8 bytes; `long` stays at 4.

**Trap alert:** GATE questions sometimes assume a 2-byte `int`. Read the question's assumption before computing a size.

**Interview angle:** *"What is sizeof(long) on a 64-bit machine?"* Model answer: "It depends on the data model, not the CPU. It's 8 on Linux (LP64) and 4 on Windows (LLP64). The only size the standard fixes is `sizeof(char) == 1`. For exact widths, use `int32_t`, `int64_t` and friends from `<stdint.h>`." Follow-up: *"Why does `int` stay 4 bytes on 64-bit machines?"* Because widening it would have broken a huge amount of existing code and data formats, and 32-bit arithmetic is just as fast on x86-64.

---

## Q2. [M1 | Signedness & unsigned wrap-around | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P01, P02

```c
#include <stdio.h>
int main(void) {
    unsigned int u = -1;
    unsigned char c = 255;
    c = c + 1;
    printf("%u %d\n", u, c);
    return 0;
}
```

What does this program print?

- A. `-1 256`
- B. `4294967295 256`
- C. `4294967295 0`
- D. Nothing is guaranteed: unsigned overflow is undefined behaviour

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** Converting a value to an unsigned type is always defined: the result is the value **modulo 2ⁿ**, where n is the number of bits in the type (N1570 §6.3.1.3p2).

| Line | Statement executed | Variables after execution |
|---|---|---|
| 3 | `unsigned int u = -1;` | −1 + 2³² = **u = 4294967295** (`UINT_MAX`) |
| 4 | `unsigned char c = 255;` | c = 255 |
| 5 | `c + 1` | **integer promotion**: `c` becomes `int` 255, then 255 + 1 = `int` 256. Nothing overflows. |
| 5 | `c = 256` | stored into `unsigned char`: 256 mod 256 = **c = 0** |
| 6 | `printf` | `c` is promoted to `int` 0 for `%d`. Prints `4294967295 0` |

**Why the other options fail:**
- A: `u` holds no negative value. `%u` prints it as unsigned.
- B: An `unsigned char` cannot hold 256. The assignment wraps it to 0.
- D: Only **signed** overflow is UB. Unsigned arithmetic, and conversion to an unsigned type, wrap around by definition.

**Concept recap:**
- **Integer promotion:** before arithmetic, `char` and `short` (signed or unsigned) are converted to `int`, because `int` can hold all their values.
- **Conversion to unsigned:** take the value modulo 2ⁿ. This is always defined.
- **Conversion to signed** when the value doesn't fit: implementation-defined (see Q11).

**Trap alert:** `unsigned int u = -1;` is the standard idiom for "all bits set". It is well-defined, not a bug.

**Interview angle:** *"Is unsigned overflow UB?"* "No. Unsigned arithmetic is defined modulo 2ⁿ. Only signed overflow is UB." Follow-up: *"Why the difference?"* Signed overflow being UB lets the compiler assume `i + 1 > i` for signed `i`, which enables loop optimisations. Unsigned wrap-around is needed for hashing, checksums and modular arithmetic.

---

## Q3. [M1 | sizeof operand is not evaluated | Easy | MCQ | [C] · C11 | All] · Type T1 · P05

```c
#include <stdio.h>
int main(void) {
    int i = 5;
    size_t s = sizeof(i++);
    printf("%d %zu\n", i, s);
    return 0;
}
```

What does this program print?

- A. `6 4`
- B. `5 4`
- C. `5 8`
- D. `6 8`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:** `sizeof` only needs the **type** of its operand, and the type is known at compile time. So the operand is **not evaluated**, unless its type is a variable-length array (N1570 §6.5.3.4p2). The `++` never runs.

| Line | Statement executed | Variables after execution |
|---|---|---|
| 3 | `int i = 5;` | i = 5 |
| 4 | `sizeof(i++)` | type of `i++` is `int`, so the result is 4. **`i++` is not executed** | i = 5, s = 4 |
| 5 | `printf` | prints `5 4` |

**Why the other options fail:**
- A: This assumes `i++` ran. It does not.
- C, D: `i++` has type `int` (4 bytes). `size_t` being 8 bytes is irrelevant; `s` stores the size of an `int`.

**Concept recap:** `sizeof expr` is a compile-time constant (except for VLAs). Function calls, `++` and assignments inside it never execute: `sizeof(f())` does not call `f`.

**Trap alert:** The VLA exception. With `int n = 3;`, `sizeof(int[n++])` **does** evaluate `n++`, because the size is only known at run time. This is covered in M5.

**Interview angle:** *"Does `sizeof(func())` call `func`?"* "No. It's an unevaluated operand; only the return type matters." Follow-up: *"Name other unevaluated contexts in C++."* `decltype`, `noexcept(expr)`, `typeid` on a non-polymorphic type, and `alignof`.

---

## Q4. [M1 | sizeof('a') in C vs C++ | Medium | MCQ | [C] · C11 | All] · Type T1 · P05, P02

```c
#include <stdio.h>
int main(void) {
    char c = 'a';
    printf("%zu %zu %zu\n", sizeof('a'), sizeof(c), sizeof(c + 1));
    return 0;
}
```

Compiled **as C11**, what does this print?

- A. `1 1 1`
- B. `1 1 4`
- C. `4 1 1`
- D. `4 1 4`

**Answer:** D

**Behaviour class:** Well-defined (sizes as per the stated platform).

**Why:**

| Expression | Type in C | sizeof |
|---|---|---|
| `'a'` | **`int`**: in C, a character constant has type `int` (N1570 §6.4.4.4p10) | 4 |
| `c` | `char` | 1 |
| `c + 1` | `c` is promoted to `int`, and `int + int` is `int` | 4 |

Prints `4 1 4`.

**Why the other options fail:**
- A: This assumes `c + 1` stays a `char`. Arithmetic never happens in `char`; promotion comes first.
- B: This is the **C++** answer. In C++, `'a'` has type `char`, so the output is `1 1 4`.
- C: This gets both of the "surprising" ones backwards.

**Concept recap:** The same source line gives a different result in C and C++. C keeps `int` character constants for historical reasons. C++ made them `char` so that overloads like `f(char)` vs `f(int)` pick the expected function for `f('a')`.

**Trap alert:** Any question showing `sizeof('x')` → check the language tag first.

**Interview angle:** *"Give a program that is valid in both C and C++ but behaves differently."* `printf("%zu", sizeof('a'));` prints 4 in C and 1 in C++. Other examples: the meaning of an empty parameter list `f()`, and implicit conversion from `void*`.

---

## Q5. [M1 | Literal forms | Easy | MCQ | [C] · C11 | GATE] · Type T1 · P10

Assume ASCII.

```c
#include <stdio.h>
int main(void) {
    int x = 010 + 0x10 + 10;
    int y = 'A' + '\n';
    printf("%d %d\n", x, y);
    return 0;
}
```

What does this program print?

- A. `30 75`
- B. `36 75`
- C. `34 75`
- D. `34 66`

**Answer:** C

**Behaviour class:** Well-defined.

**Why:** In C, a leading `0` means **octal** and a leading `0x` means **hex**.

| Line | Statement executed | Variables after execution |
|---|---|---|
| 3 | `010` = 8 (octal), `0x10` = 16 (hex), `10` = 10 → 8 + 16 + 10 | x = 34 |
| 4 | `'A'` = 65, `'\n'` = 10 (ASCII line feed) → 65 + 10 | y = 75 |
| 5 | `printf` | prints `34 75` |

**Why the other options fail:**
- A: Treats both `010` and `0x10` as decimal 10.
- B: Treats `010` as decimal 10: 10 + 16 + 10 = 36.
- D: Treats `'\n'` as 1. It is ASCII 10.

**Concept recap:**
- Integer literal prefixes: `0` = octal, `0x`/`0X` = hex, no prefix = decimal. C23 adds `0b` for binary (a GCC extension before that).
- Common escapes: `'\0'` = 0, `'\n'` = 10, `'\t'` = 9, `'0'` = 48, `'A'` = 65, `'a'` = 97.

**Trap alert:** Zero-padding numbers "for alignment", such as `int codes[] = {001, 010, 100};`, silently creates octal values. Also, `09` is a **compile error**, because 9 is not an octal digit.

**Interview angle:** *"Why is `0` itself technically an octal literal?"* The grammar defines an octal constant as `0` followed by zero or more octal digits, so `0` qualifies. It has the same value either way, so it doesn't matter, but it's a nice detail to know.

---

## Q6. [M1 | Escape sequences inside strings | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P10

Assume ASCII.

```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char s[] = "a\0b\x41\101";
    printf("%zu %zu\n", sizeof s, strlen(s));
    printf("%s\n", s + 2);
    return 0;
}
```

What does this program print?

- A. `6 1` then `bAA`
- B. `6 5` then `bAA`
- C. `13 1` then `b`
- D. `5 1` then `bA`

**Answer:** A

**Behaviour class:** Well-defined.

**Why:** First count the characters the literal actually contains. Each escape sequence is **one** character.

| Index | Source text | Character |
|---|---|---|
| 0 | `a` | `'a'` |
| 1 | `\0` | NUL (0) |
| 2 | `b` | `'b'` |
| 3 | `\x41` | hex 41 = 65 = `'A'` (a hex escape stops at the first non-hex character, here `\`) |
| 4 | `\101` | octal 101 = 65 = `'A'` (an octal escape takes at most 3 digits) |
| 5 | (implicit) | the terminating NUL |

| Line | Statement executed | Result |
|---|---|---|
| 4 | `char s[] = ...` | array of **6** `char`s |
| 5 | `sizeof s` | 6, the whole array including both NULs |
| 5 | `strlen(s)` | counts until the first NUL, at index 1 → **1** |
| 6 | `s + 2` | points at index 2 → prints `b`, `A`, `A`, then stops at index 5 → `bAA` |

**Why the other options fail:**
- B: `strlen` stops at the embedded `\0`. It never sees the last four characters.
- C: 13 counts the backslashes and digits as separate characters, i.e. 12 source characters plus 1. Escapes are single characters.
- D: Forgets the implicit terminating NUL.

**Concept recap:** `sizeof` on a char **array** = number of elements, including every NUL. `strlen` = number of characters before the first NUL. The two only agree (as `strlen + 1 == sizeof`) when there is no embedded NUL.

**Trap alert:** A hex escape is greedy: `"\x41BC"` is a single escape `\x41BC`, which is out of range for `char` (GCC reports "hex escape sequence out of range"). To write "A" followed by "BC", split the literal: `"\x41" "BC"`. Adjacent literals are concatenated.

**Interview angle:** *"How do you store binary data containing zero bytes in C?"* Use an array plus an explicit length (`memcpy`, `fwrite`), never the `str*` functions, which stop at the first zero byte.

---

## Q7. [M1 | static vs extern linkage across files | Medium | MCQ | [C] · C11 | GATE+OA] · Type T3 · P14, P42

The program is built with `gcc a.c b.c -o prog`.

```c
/* a.c */
static int count = 10;
int get(void) { return count; }
```
```c
/* b.c */
#include <stdio.h>
extern int count;
int get(void);
int main(void) {
    printf("%d\n", get() + count);
    return 0;
}
```

What is the result?

- A. Prints `20`
- B. Compile error in `b.c`: `count` is `static` in `a.c`
- C. Prints `10`, because `b.c` gets its own zero-initialised `count`
- D. Link error: undefined reference to `count`

**Answer:** D

**Behaviour class:** Linker error.

**Why:**
- `static` at file scope gives `count` **internal linkage**. The name is visible only inside `a.c`, and no global symbol `count` is exported from `a.o` (it appears as a local symbol, `d count`, in `nm`).
- `b.c` compiles cleanly. `extern int count;` only *declares* that some other file defines `count`. Each `.c` file is compiled separately, so the compiler of `b.c` cannot know what `a.c` contains.
- At link time, `b.o` needs a symbol `count` with external linkage. None exists → `undefined reference to 'count'`.
- `get` links fine. It has external linkage (the default for functions), and it can read its own file's `count`.

**Why the other options fail:**
- A: This would be the output if `count` were not `static`.
- B: Each translation unit is compiled independently, so `b.c` cannot see `a.c`'s `static`.
- C: `extern int count;` is a declaration only; it does not allocate storage. If `b.c` instead had `int count;` (a tentative definition), the program would link and print `10 + 0 = 10`.

**Concept recap:** Linkage controls whether the same name in different files refers to the same object.
- **External linkage:** shared across files. This is the default for file-scope objects and functions.
- **Internal linkage:** private to one file. Created by `static` at file scope.
- **No linkage:** local variables, including block-scope `static` locals.

**Trap alert:** `static` means two unrelated things. At **file scope** it means internal linkage. At **block scope** it means static storage duration (the variable lives for the whole program run).

**Interview angle:** *"How do you make a global variable private to one .c file?"* "Declare it `static` at file scope. In C++, an unnamed namespace is the preferred way." Follow-up: *"Is this a compile error or a link error, and how can you tell?"* It's a link error: compilation of each file succeeded, and the message comes from `ld` ("undefined reference"), not from `cc1`.

---

## Q8. [M1 | Storage classes & linkage rules | Medium | MSQ | [C] · C11 | GATE] · Type — concept + code reading · P14

Select **all** true statements about C11.

- A. A file-scope declaration `static int x;` gives `x` internal linkage.
- B. A block-scope declaration `static int y;` gives `y` no linkage but static storage duration.
- C. A file-scope `int z;` (no initializer) is only a declaration; it never defines storage.
- D. The file-scope declaration `extern int w = 5;` is a definition of `w`.
- E. A `static` function in one file can be called from another file if that file declares it with `extern`.

**Answer:** A, B, D

**Why each option:**
- **A ✓:** `static` at file scope means internal linkage (N1570 §6.2.2p3).
- **B ✓:** A block-scope `static` lives for the whole program (it is initialised once, before `main`) but has no linkage, so no other scope can refer to it by name.
- **C ✗:** `int z;` at file scope is a **tentative definition**. If no definition with an initializer appears later in the file, it becomes a real definition initialised to 0 (§6.9.2p2).
- **D ✓:** A file-scope declaration with an initializer is an external definition, even with `extern` (§6.9.2p1). GCC warns "'w' initialized and declared 'extern'", but it is valid. At **block** scope, `extern int w = 5;` is a constraint violation and a compile error.
- **E ✗:** The `static` function's name is not exported, so the other file gets a link error ("undefined reference").

**Concept recap:**

| Declaration | Scope | Linkage | Storage duration |
|---|---|---|---|
| `int g;` (file) | file | external | static |
| `static int g;` (file) | file | internal | static |
| `int a;` (block) | block | none | automatic |
| `static int a;` (block) | block | none | static |
| `extern int a;` (block) | block | external | static (refers to a definition elsewhere) |

**Trap alert:** **Scope**, **lifetime** and **linkage** are three separate properties. Scope = where the name is visible. Lifetime = when the storage exists. Linkage = whether the name connects across files. A block-scope `static` has small scope but a program-long lifetime.

**Interview angle:** *"What does `static` mean in C?"* Give all three uses: (1) a file-scope variable or function → internal linkage; (2) a block-scope variable → persists across calls; (3) in a parameter, `int a[static 10]` → a C99 promise to the compiler that at least 10 elements are passed.

---

## Q9. [M1 | Compile error vs link error | Medium | MSQ | [C] · C11 | All] · Type T3 · P42

Select **all** situations that produce a **linker** error (not a compile error) when building a multi-file program with `gcc -std=c11`.

- A. `main.c` calls `int square(int);`, which is declared with a prototype but defined in no file.
- B. `main.c` uses an identifier `total` that is not declared anywhere.
- C. `a.c` and `b.c` both contain, at file scope, `int g = 1;`
- D. `main.c` passes two arguments to a function whose visible prototype takes one.
- E. `main.c` contains `extern int n;` and reads `n`, but no file defines `n`.

**Answer:** A, C, E

**Why each option:**
- **A ✓:** The prototype satisfies the compiler, so `main.o` is produced with an unresolved reference to `square` → `undefined reference to 'square'`.
- **B ✗:** Using an undeclared identifier is a compile error ("'total' undeclared").
- **C ✓:** Both object files define the external symbol `g` → `multiple definition of 'g'`. Strictly speaking, the C standard calls having two external definitions undefined behaviour with no diagnostic required (§6.9p5). In practice every linker reports it as an error. (GCC 10+ uses `-fno-common` by default, so even two `int g;` without initializers fail to link. Older GCCs merged them as "common" symbols.)
- **D ✗:** A wrong number of arguments when a prototype is in scope is a constraint violation, so the compiler must issue a diagnostic → compile error.
- **E ✓:** Same idea as A, but for an object: the declaration compiles, and the definition is missing at link time.

**Concept recap:** The compiler checks **one translation unit** (one preprocessed `.c` file) at a time: syntax, types, declarations. The linker checks **across** units: that every referenced symbol is defined exactly once.

**Trap alert:** "It compiled" does not mean it will link. When asked "compile error or link error?", ask: *could a compiler that sees only this one file detect the problem?* If yes, it's a compile error; if not, it's a link error.

**Interview angle:** *"What does 'undefined reference to X' mean, and how do you fix it?"* "The linker couldn't find a definition of X. Usually the `.c`/`.o` file that defines it wasn't passed to the linker, a library is missing (`-lm`), the names don't match (C vs C++ name mangling, so `extern "C"` is needed), or the definition is `static`."

---

## Q10. [M1 | Default initialisation & uninitialised reads | Hard | MCQ | [C] · C11 | All] · Type T3 · P48

```c
#include <stdio.h>
int g;
int main(void) {
    static int s;
    int a;
    printf("%d %d %d\n", g, s, a);
    return 0;
}
```

Which statement is correct?

- A. The program prints `0 0 0`.
- B. The behaviour is undefined, because it reads `a`, an uninitialised automatic variable whose address is never taken, so no output is guaranteed. (`g` and `s` by themselves would be guaranteed to be 0.)
- C. The program prints `0 0` followed by some unspecified garbage value, and is otherwise well-defined.
- D. Compile error: `a` is used uninitialised.

**Answer:** B

**Behaviour class:** **Undefined behaviour.**

> **⚠ UB box**
> - **Category:** reading an uninitialised (indeterminate) automatic variable.
> - **Rule:** N1570 §6.3.2.1p2: if an automatic object could have been declared `register` (its address is never taken) and it is used before being initialised, the behaviour is undefined. In C++ it is UB in all such cases ([basic.indet]).
> - **On GCC x86-64,** `-O0` typically prints `0 0` and then whatever was left in that stack slot (often 0, or a value like 32767). At `-O2` it may print anything, because the optimiser is allowed to assume the read never happens. **Neither result is guaranteed, and you must not rely on either.**

**Why:**

| Line | Statement executed | Variables after execution |
|---|---|---|
| 2 | `int g;` | static storage, no initializer → **zero-initialised**: g = 0 (§6.7.9p10) |
| 4 | `static int s;` | static storage → s = 0 |
| 5 | `int a;` | automatic storage, no initializer → value **indeterminate** |
| 6 | `printf(..., a)` | reads `a` → **UB** |

Once a program contains UB, the standard places no requirements on any of its output, including the part that prints `g` and `s`.

**Why the other options fail:**
- A: Treats `a` as zero-initialised. Only static-storage objects get automatic zero-initialisation.
- C: "Garbage but well-defined" is the classic misconception. With the address never taken, the standard makes this undefined, not just an unknown value.
- D: GCC and Clang only **warn** (`-Wuninitialized`). The code is not ill-formed.

**Concept recap:** Objects with **static** storage duration (globals and `static` locals) are zero-initialised before the program starts. **Automatic** variables (ordinary locals) are not initialised at all.

**Trap alert:** GATE answer keys sometimes say "garbage value". The precise answer is *undefined behaviour*: the compiler may delete code, print different values on different runs, or appear to work.

**Interview angle:** *"What's the value of an uninitialised local int?"* "It's indeterminate, and reading it is undefined behaviour. Globals and statics are zero-initialised; locals aren't." Follow-up: *"How do you catch these?"* `-Wall -Wuninitialized`, MemorySanitizer (`-fsanitize=memory`, Clang), Valgrind's memcheck, or always initialising at declaration.

---

## Q11. [M1 | Integer truncation on conversion | Hard | MCQ | [C] · C11 | GATE+OA] · Type T3 · P50

```c
#include <stdio.h>
int main(void) {
    int big = 300;
    unsigned char uc = big;
    signed char sc = 200;
    short sh = 70000;
    printf("%d %d %d\n", uc, sc, sh);
    return 0;
}
```

Which statement is correct? (`SCHAR_MAX` = 127, `SHRT_MAX` = 32767.)

- A. It prints `44 -56 4464`, and the standard guarantees all three values.
- B. Undefined behaviour, because of signed overflow when storing into `sc` and `sh`.
- C. `uc == 44` is guaranteed by the standard. The values of `sc` and `sh` are **implementation-defined**; GCC documents modulo-2ⁿ reduction, so it prints `44 -56 4464`.
- D. Compile error: narrowing conversion in initialisation.

**Answer:** C

**Behaviour class:** **Implementation-defined** for `sc` and `sh`. The `uc` conversion is well-defined.

**Why:**

| Line | Statement executed | Rule | Value (GCC) |
|---|---|---|---|
| 4 | `uc = 300` | to **unsigned**: modulo 2⁸ (§6.3.1.3p2), guaranteed | 300 − 256 = **44** |
| 5 | `sc = 200` | to **signed**, value doesn't fit: implementation-defined result, or an implementation-defined signal is raised (§6.3.1.3p3) | GCC: 200 − 256 = **−56** |
| 6 | `sh = 70000` | same rule | GCC: 70000 − 65536 = **4464** |
| 7 | `printf` | all three are promoted to `int` for `%d` | `44 -56 4464` |

**Why the other options fail:**
- A: Right output on GCC, but the standard guarantees only the unsigned conversion.
- B: **Conversion is not arithmetic overflow.** Signed overflow UB applies to operations like `INT_MAX + 1`. An out-of-range conversion to a signed type is implementation-defined. This is exactly the UB vs implementation-defined confusion the taxonomy warns about.
- D: Narrowing is an error only in C++ **brace** initialisation (`short sh{70000};`). C has no such rule; GCC only warns with `-Woverflow`.

**Concept recap:** Narrowing an integer:
- to an unsigned type → keep the low n bits (well-defined);
- to a signed type, value fits → the value is unchanged;
- to a signed type, value doesn't fit → implementation-defined.

**Trap alert:** **Overflow** (an operation produces an out-of-range result) ≠ **conversion** (an existing value is stored into a narrower type). They have different behaviour classes.

**Interview angle:** *"Is `char c = 200;` UB?"* "No. If `char` is signed, the result is implementation-defined; GCC gives −56. It would be UB only if an *arithmetic operation* on signed values overflowed." Follow-up: *"What did C23 change?"* C23 requires two's complement representation, but the conversion rule is still worded as implementation-defined.

---

## Q12. [M1 | Signedness of plain char | Medium | MCQ | [C] · C11 | All] · Type T1 · P97

Platform: x86-64 GCC, where plain `char` is **signed** (the default stated above).

```c
#include <stdio.h>
int main(void) {
    char c = 0xFF;
    if (c == 0xFF)
        printf("equal\n");
    else
        printf("not equal\n");
    printf("%d\n", c);
    return 0;
}
```

What does this program print?

- A. `not equal` then `-1`
- B. `equal` then `-1`
- C. `equal` then `255`
- D. Undefined behaviour: `0xFF` does not fit in `char`

**Answer:** A

**Behaviour class:** **Implementation-defined.** Both the signedness of plain `char` (§6.2.5p15) and the result of storing 255 into a signed `char` (§6.3.1.3p3) are implementation-defined. Under the stated platform the output is determined.

**Why:**

| Line | Statement executed | Variables after execution |
|---|---|---|
| 3 | `char c = 0xFF;` | `0xFF` is the `int` 255. It doesn't fit in a signed `char`, so the value is implementation-defined; GCC stores bit pattern 1111 1111 = **c = −1** |
| 4 | `c == 0xFF` | `c` is promoted to `int` → **−1**. The right side is `int` **255**. −1 ≠ 255 → false |
| 7 | `printf("not equal\n")` | — |
| 8 | `printf("%d", c)` | promoted to `int` −1 → prints `-1` |

**Why the other options fail:**
- B: The bit patterns are "the same", but the comparison happens on promoted `int` values, which differ.
- C: This is the output on **ARM Linux / PowerPC**, where plain `char` is **unsigned**, or with `gcc -funsigned-char`.
- D: This is a conversion, not overflow, so it is implementation-defined, not UB (see Q11).

**Concept recap:** `char`, `signed char` and `unsigned char` are **three distinct types**. Whether plain `char` behaves like the signed or the unsigned one is up to the implementation. Byte-processing code should use `unsigned char` or `uint8_t`.

**Trap alert:** The same issue causes the classic `getchar()` bug: `char c; while ((c = getchar()) != EOF)`. EOF is −1 (an `int`), and storing the result in a `char` breaks the comparison. This is covered in M11.

**Interview angle:** *"Is char signed or unsigned?"* "It's implementation-defined. It's signed on x86 GCC/MSVC and unsigned on ARM Linux. Never rely on it; use `signed char`/`unsigned char` explicitly when the sign matters."

---

## Q13. [M1 | What the standard guarantees about integer types | Medium | MSQ | [C] · C11 | GATE] · Type — concept · P98

Select **all** properties that the **C11 standard itself** guarantees on every conforming implementation.

- A. `sizeof(char) == 1`
- B. Every value representable in `int` is also representable in `long`.
- C. `int` can represent at least the range −32767 to +32767.
- D. `sizeof(int) == 4`
- E. `CHAR_BIT == 8` (a byte has exactly 8 bits)

**Answer:** A, B, C

**Why each option:**
- **A ✓:** By definition (§6.5.3.4p4). `sizeof` counts in units of `char`.
- **B ✓:** For signed types of increasing **rank** (`signed char` < `short` < `int` < `long` < `long long`), each type's range is a subrange of the next one's (§6.2.5p8).
- **C ✓:** `INT_MIN` ≤ −32767 and `INT_MAX` ≥ +32767 (§5.2.4.2.1). So `int` is at least 16 bits.
- **D ✗:** This is **implementation-defined**. It is 4 on all mainstream desktop ABIs, but 2 on many 16-bit microcontrollers.
- **E ✗:** The standard only requires `CHAR_BIT ≥ 8`. Some DSPs use 16- or 32-bit bytes. (POSIX additionally requires 8.)

**Concept recap:** Guaranteed minimum widths: `char` ≥ 8 bits, `short` ≥ 16, `int` ≥ 16, `long` ≥ 32, `long long` ≥ 64. Everything else (exact sizes, padding, alignment) is implementation-defined.

**Trap alert:** "Implementation-defined" is not UB. A program that prints `sizeof(int)` is perfectly well-defined; the number it prints is documented by the implementation.

**Interview angle:** *"How would you write a struct for a network packet with a 32-bit field?"* "Use `uint32_t` from `<stdint.h>`, never `int`/`long`. Then handle byte order with `htonl`/`ntohl`, and handle padding (`static_assert` the struct size, or serialise field by field)."

---

## Q14. [M1 | Valid signatures of main | Easy | MSQ | [C] · C11 | GATE] · Type T2 · —

In a **hosted** C11 environment (one with an OS and the full standard library), which definitions of `main` must **every** conforming implementation accept?

- A. `int main(void) { return 0; }`
- B. `int main(int argc, char *argv[]) { return 0; }`
- C. `void main(void) { }`
- D. `int main(int argc, char **argv) { return 0; }`
- E. `main() { return 0; }`

**Answer:** A, B, D

**Why each option:**
- **A ✓, B ✓:** These are the two forms the standard lists (§5.1.2.2.1).
- **D ✓:** `char **argv` is equivalent to `char *argv[]`, because an array parameter is adjusted to a pointer.
- **C ✗:** The standard also allows main to be defined "in some other implementation-defined manner". So `void main` *may* be accepted by a particular compiler, but it is not portable, and the exit status becomes unspecified.
- **E ✗:** This relies on **implicit int**, which was removed in C99. A declaration must now have at least one type specifier (§6.7.2p2, a constraint), so this needs a diagnostic. GCC warns by default; `-pedantic-errors` makes it an error.

**Concept recap:** Since C99, if control reaches the closing `}` of `main` without a `return`, it is as if `return 0;` were executed (§5.1.2.2.3). This applies only to `main`. Falling off the end of any other non-`void` function and then *using* the return value is UB.

**Trap alert:** `int main()` (empty parentheses) is accepted in practice. In C before C23, `()` means "unspecified parameters"; C23 makes `()` mean `(void)`, as C++ always has.

**Interview angle:** *"Why not `void main`?"* "The return value of `main` is the process exit status that the shell, scripts and parent processes check. `void main` is non-standard, and the status becomes garbage or unspecified."

---

## Q15. [M1 | Declaration vs definition: tentative definitions | Hard | MCQ | [C] · C11 | GATE] · Type T1 · —

The whole program is the single file below.

```c
#include <stdio.h>
int x;
int x;
int x = 5;
int main(void) {
    printf("%d\n", x);
    return 0;
}
```

What happens?

- A. Compile error: redefinition of `x`
- B. It prints `0`
- C. Link error: multiple definition of `x`
- D. It prints `5`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:** At file scope, a declaration with no initializer and no storage class (or with `static`) is a **tentative definition** (§6.9.2p2). Within one file you may have any number of them. If one real definition (with an initializer) exists, the tentative ones simply refer to it.

| Line | Declaration | Status |
|---|---|---|
| 2 | `int x;` | tentative definition |
| 3 | `int x;` | tentative definition (allowed to repeat) |
| 4 | `int x = 5;` | **the** external definition |
| 6 | `printf` | `x` = 5 → prints `5` |

**Why the other options fail:**
- A: This is the answer **in C++**, which has no tentative definitions: `int x; int x;` is a redefinition there (one-definition rule).
- B: Line 4's initializer wins. Zero-initialisation is used only if no line provides an initializer.
- C: It's all one file, so there's only one definition for the linker. (Across *two* files, see Q9 option C.)

**Concept recap:**
- **Declaration:** introduces a name and its type.
- **Definition:** a declaration that also allocates storage (objects) or provides a body (functions).
- `extern int x;` = declaration only. `int x = 5;` = definition. `int x;` at file scope = tentative definition.

**Trap alert:** `int x = 5; int x = 6;` in the same file **is** a compile error, even in C. Two initializers mean two definitions.

**Interview angle:** *"Difference between declaration and definition?"* Give the three-line summary above, plus: "A program can declare something many times but must define it exactly once (the one-definition rule). Headers hold declarations; one `.c` file holds the definition."

---

## Q16. [M1 | Scope, shadowing & block-scope extern | Hard | MCQ | [C] · C11 | GATE+OA] · Type T1 · P14

```c
#include <stdio.h>
int x = 1;
int main(void) {
    int x = 2;
    {
        int x = 3;
        printf("%d ", x);
    }
    printf("%d ", x);
    {
        extern int x;
        printf("%d\n", x);
    }
    return 0;
}
```

What does this program print?

- A. `3 3 3`
- B. `3 2 1`
- C. `3 2 2`
- D. Compile error: conflicting declarations of `x`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:**

| Line | Statement executed | Which `x` is visible | Output so far |
|---|---|---|---|
| 2 | `int x = 1;` | global x (external linkage) = 1 | |
| 4 | `int x = 2;` | main's local x = 2 hides the global | |
| 6 | `int x = 3;` | inner block's x = 3 hides main's x | |
| 7 | `printf` | inner x → 3 | `3 ` |
| 8 | `}` | inner x's lifetime ends | |
| 9 | `printf` | main's x → 2 | `3 2 ` |
| 11 | `extern int x;` | This block-scope declaration has **external linkage**. The visible prior declaration (main's x) has *no* linkage, so by §6.2.2p4 the new one gets external linkage and refers to the **global** x | |
| 12 | `printf` | global x → 1 | `3 2 1` |

**Why the other options fail:**
- A, C: Assume that an inner declaration changes the outer variables, or that `extern` picks up the nearest x. `extern` always refers to an object with linkage, never to a local variable.
- D: Each `x` is declared in a **different block scope**, so nothing conflicts. Redeclaring a no-linkage name is an error only in the *same* scope.

**Concept recap:** Name lookup walks outward from the innermost block. An inner declaration *hides* outer ones for the rest of that block. `extern` inside a block is the one way to reach a global that a local has hidden. In C++ you would write `::x` instead.

**Trap alert:** Shadowing never changes or "copies" the outer variable. They are completely separate objects that happen to have the same name.

**Interview angle:** *"How do you access a global hidden by a local?"* "In C++: `::x`. In C: open a new block and declare `extern int x;`, or better, don't shadow. Compile with `-Wshadow`."

---

## Q17. [M1 | typedef vs #define & enum values | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P05

Platform: default (`sizeof(int*)=8`, `sizeof(int)=4`).

```c
#include <stdio.h>
typedef int *IntPtr;
#define INTPTR int *
enum level { LOW, MID = 5, HIGH, TOP = 3, MAX };
int main(void) {
    IntPtr p1, p2;
    INTPTR p3, p4;
    printf("%zu %zu %zu %zu\n", sizeof p1, sizeof p2, sizeof p3, sizeof p4);
    printf("%d %d %d\n", HIGH, TOP, MAX);
    return 0;
}
```

What does this program print?

- A. `8 8 8 4` then `6 3 4`
- B. `8 8 8 8` then `6 3 4`
- C. `8 4 8 4` then `6 3 4`
- D. `8 8 8 4` then `2 3 4`

**Answer:** A

**Behaviour class:** Well-defined. The pointers are never read: `sizeof` doesn't evaluate its operand, so their being uninitialised does not matter.

**Why:**

| Line | After processing | Types / values |
|---|---|---|
| 6 | `IntPtr p1, p2;` | `typedef` creates a real **type alias**, so it applies to every declarator → p1: `int*`, p2: `int*` |
| 7 | the preprocessor turns it into `int * p3, p4;` | `*` binds to the **declarator** p3 only → p3: `int*`, **p4: `int`** |
| 8 | sizes | 8, 8, 8, 4 |
| 4 | enum values | LOW = 0, MID = 5, HIGH = 6 (previous + 1), TOP = 3, MAX = 4 (previous + 1) |
| 9 | `printf` | `6 3 4` |

**Why the other options fail:**
- B: Treats the macro like a typedef. A macro is plain text substitution.
- C: Gets the typedef case wrong. A typedef'd pointer type applies to every name declared.
- D: Counts HIGH as 2, as if the enumerators were numbered by position. An enumerator without `=` is the **previous enumerator + 1**.

**Concept recap:**
- `typedef` is handled by the compiler and understands C types.
- `#define` is handled by the preprocessor and only replaces tokens.
- Enum constants have type `int` in C, and duplicate values are allowed.

**Trap alert:** `int* a, b;` declares one pointer and one `int`, however the spacing looks. Declare one pointer per line.

**Interview angle:** *"typedef vs #define for type aliases?"* "typedef is scoped, type-checked and handles compound declarators; #define is blind text substitution. In C++, prefer `using IntPtr = int*;`."

---

## Q18. [M1 | Compilation stages | Easy | MCQ | [C] · C11 | GATE+OA] · Type T2 · P42

```c
#include <stdio.h>
#define SIZE 10
int main(void) {
    int arr[SIZE] = {0}
    printf("%d\n", arr[0]);
    return 0;
}
```

The four commands below are run one at a time. Which is the **first** to fail?

- A. `gcc -E main.c` (preprocess only)
- B. None of them: the program builds, and the error appears only at run time
- C. `gcc -S main.c` (preprocess + compile to assembly)
- D. `gcc main.c` (the final link step only)

**Answer:** C

**Behaviour class:** Ill-formed (syntax error: missing `;` on line 4).

**Why:**

| Stage | gcc flag | Input → output | What happens here |
|---|---|---|---|
| Preprocess | `-E` | `.c` → `.i` | Pastes in `stdio.h`, replaces `SIZE` with `10`, removes comments. It does **not** parse C syntax, so it succeeds. |
| Compile | `-S` | `.i` → `.s` | Parses the C grammar → **"expected ',' or ';' before 'printf'"** — fails here |
| Assemble | `-c` | `.s` → `.o` | never reached |
| Link | (none) | `.o` + libs → executable | never reached |

**Why the other options fail:**
- A: The preprocessor works on tokens and `#` directives only. It can't tell that a `;` is missing.
- B: Syntax errors are always detected at compile time.
- D: Linking needs an object file, which was never produced.

**Concept recap:** The pipeline is preprocess (`cpp`) → compile (`cc1`) → assemble (`as`) → link (`ld`). Errors from each stage look different:
- Preprocessor: "No such file or directory" for a missing header, `#error`.
- Compiler: syntax and type errors.
- Linker: "undefined reference", "multiple definition".

**Trap alert:** A missing **header file** fails at `-E`. A missing **function definition** fails only at link time.

**Interview angle:** *"Walk me through what happens when you run `gcc main.c`."* List the four stages with their files: `.i`, `.s`, `.o`, `a.out`. Mention that `-save-temps` keeps all the intermediate files, so you can inspect each stage.

---

## Q19. [M2 | Operator precedence | Easy | MCQ | [C] · C11 | GATE+OA] · Type T1 · P03

```c
#include <stdio.h>
int main(void) {
    int a = 2, b = 3, c = 4;
    int r = a + b * c % 5 - 1;
    int s = a & b == 2;
    printf("%d %d\n", r, s);
    return 0;
}
```

What does this program print?

- A. `13 1`
- B. `3 1`
- C. `13 0`
- D. `3 0`

**Answer:** D

**Behaviour class:** Well-defined.

**Why:**

| Line | Statement executed | How it groups | Variables after execution |
|---|---|---|---|
| 3 | declarations | — | a=2, b=3, c=4 |
| 4 | `a + b * c % 5 - 1` | `*` and `%` share a precedence level, left to right: `(b*c) % 5` = 12 % 5 = 2; then `+`/`-` left to right: (2 + 2) − 1 | r = 3 |
| 5 | `a & b == 2` | `==` binds **tighter** than `&`: `a & (b == 2)` = 2 & 0 | s = 0 |
| 6 | `printf` | — | prints `3 0` |

**Why the other options fail:**
- A: Applies `%` before `*` (3 × (4 % 5) = 12, so r = 13) and `&` before `==` ((2 & 3) == 2 → 1).
- B: Gets r right but `&` vs `==` wrong.
- C: Gets `&` vs `==` right but `*`/`%` order wrong.

**Concept recap:** C precedence, from high to low (the ones that matter most):
1. postfix `() [] -> . ++ --`
2. unary `! ~ ++ -- + - * & sizeof` and casts
3. `* / %`
4. `+ -`
5. `<< >>`
6. `< <= > >=`
7. `== !=`
8. `&`
9. `^`
10. `|`
11. `&&`
12. `||`
13. `?:`
14. assignment
15. `,`

**Trap alert:** The bitwise operators `& ^ |` sit **below** `==`. That's a historical accident from before C had `&&`. `if (flags & MASK == 0)` is always a bug; write `(flags & MASK) == 0`.

**Interview angle:** *"Why is `x & 1 == 0` wrong for an even-check?"* "It parses as `x & (1 == 0)`, which is `x & 0`, which is always 0. Parenthesise bitwise operations inside comparisons."

---

## Q20. [M2 | Associativity | Medium | MCQ | [C] · C11 | GATE+OA] · Type T1 · P03

```c
#include <stdio.h>
int main(void) {
    int a, b, c;
    a = b = c = 5;
    int x = 20 / 5 / 2;
    int y = 2 << 1 + 1;
    int z = 10 - 4 - 3;
    printf("%d %d %d %d\n", a, x, y, z);
    return 0;
}
```

What does this program print?

- A. `5 10 5 9`
- B. `5 2 8 3`
- C. `5 2 5 3`
- D. `5 10 8 9`

**Answer:** B

**Behaviour class:** Well-defined.

**Why:**

| Line | Expression | Grouping rule | Result |
|---|---|---|---|
| 4 | `a = b = c = 5` | assignment is **right-to-left**: `a = (b = (c = 5))` | c = 5, b = 5, a = 5 |
| 5 | `20 / 5 / 2` | `/` is **left-to-right**: `(20 / 5) / 2` = 4 / 2 | x = 2 |
| 6 | `2 << 1 + 1` | `+` has **higher precedence** than `<<`: `2 << (1 + 1)` = 2 << 2 | y = 8 |
| 7 | `10 - 4 - 3` | left-to-right: `(10 - 4) - 3` | z = 3 |
| 8 | `printf` | — | `5 2 8 3` |

**Why the other options fail:**
- A: Groups everything right-to-left (20/(5/2) = 20/2 = 10; 10 − (4 − 3) = 9) and applies `<<` before `+` ((2 << 1) + 1 = 5).
- C: Gets associativity right but shift precedence wrong.
- D: Gets shift precedence right but associativity wrong.

**Concept recap:**
- **Precedence** decides which *different* operators group first.
- **Associativity** decides how *equal-precedence* operators chain.
- Right-to-left: unary operators, `?:`, and all assignment operators. Almost everything else is left-to-right.

**Trap alert:** Associativity is **not** evaluation order. In `f() + g() * h()`, the grouping is fixed, but the order in which `f`, `g` and `h` are *called* is unspecified (covered in M2/M4).

**Interview angle:** *"What does `a = b = c` do, and why does it work?"* "Assignment is right-associative and an assignment expression yields the assigned value, so `c = 5` gives 5, which is assigned to `b`, then to `a`." Follow-up: *"Is `x << 1 + 1` readable code?"* No. Shifts mixed with arithmetic should always be parenthesised.

---

## Batch 1 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 1 | M1 | E | MCQ | A | Well-defined (impl-defined sizes) |
| 2 | M1 | E | MCQ | C | Well-defined |
| 3 | M1 | E | MCQ | B | Well-defined |
| 4 | M1 | M | MCQ | D | Well-defined |
| 5 | M1 | E | MCQ | C | Well-defined |
| 6 | M1 | M | MCQ | A | Well-defined |
| 7 | M1 | M | MCQ | D | Linker error |
| 8 | M1 | M | MSQ | A, B, D | — |
| 9 | M1 | M | MSQ | A, C, E | Linker vs ill-formed |
| 10 | M1 | H | MCQ | B | **UB** |
| 11 | M1 | H | MCQ | C | Implementation-defined |
| 12 | M1 | M | MCQ | A | Implementation-defined |
| 13 | M1 | M | MSQ | A, B, C | — |
| 14 | M1 | E | MSQ | A, B, D | Ill-formed / impl-defined |
| 15 | M1 | H | MCQ | D | Well-defined |
| 16 | M1 | H | MCQ | B | Well-defined |
| 17 | M1 | M | MCQ | A | Well-defined |
| 18 | M1 | E | MCQ | C | Ill-formed |
| 19 | M2 | E | MCQ | D | Well-defined |
| 20 | M2 | M | MCQ | B | Well-defined |

- **Difficulty:** 7 E / 9 M / 4 H. M1 alone matches its plan of 6/8/4; Q19 is the first of M2's Easy questions and Q20 the first of its Medium ones.
- **Format:** 16 MCQ / 4 MSQ.
- **Answer keys (MCQ):** A×4, B×4, C×4, D×4.

---

**Completed: 20 / 400 Remaining: 380**

**Ledger, newly ✓ this batch:**
- P01 ×2 (Q1, Q2) ★ complete
- P03 ×2 (Q19, Q20) ★ complete
- P05 ×3 (Q3, Q4, Q17) ★ complete
- P10 ×2 (Q5, Q6)
- P14 ×3 (Q7, Q8, Q16)
- P42 ×3 (Q7, Q9, Q18) ★ complete
- P48 ×1 (Q10) ★ — needs 1 more (planned in M4)
- P50 ×1 (Q11)
- P97 ×1 (Q12)
- P98 ×1 (Q13)
- P02 ×2 as a secondary pattern (Q2, Q4). Primary P02 questions are coming in Batch 2.

*Say **"Continue."** for Batch 2 (Q21–Q40, M2).*
