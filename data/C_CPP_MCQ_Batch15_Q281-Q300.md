# C/C++ MCQ/MSQ Bank — Batch 15 (Q281–Q300)

**What this batch covers:**
- **M17 Overloading & operator overloading (Q281–Q292):**
  - how the compiler ranks overloads: exact match, promotion, conversion
  - ambiguous calls
  - a standard conversion beating a user-defined one (the `"hi"` → `bool` trap)
  - default arguments that cause ambiguity, and the rules for default arguments
  - overloading on `T&` / `const T&` / `T&&`
  - chained and parenthesised assignment
  - self-assignment going wrong
  - arithmetic and comparison operators
  - why `operator<<` must be a non-member
  - conversion operators that introduce ambiguity
  - the rules for implicit conversions
- **M18 Copy/move, part 1 (Q293–Q300):**
  - value categories: lvalue, xvalue and prvalue
  - choosing between the copy and move constructor, including `const` objects that can't really be moved
  - a named rvalue reference is an lvalue
  - when the compiler doesn't generate a move constructor
  - shallow copies of owning classes (double delete), and the rule-of-three fix
  - copy-and-swap

**Patterns:** P16, P59, P69, P70 (preview)

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - Object traces show every copy/move/assignment and why it was chosen.

> **Conventions for every question in this batch:**
> - **C++17**, exceptions enabled, with guaranteed copy elision from prvalues.
> - Platform: a typical 64-bit system.
> - Needed headers are included. In code, `cout` means `std::cout`.

---

## Q281. [M17 | Ranking overloads: exact, promotion, conversion | Easy | MCQ | [C++] · C++17 | OA] · Type T6 · P16

```cpp
#include <iostream>
using std::cout;
void f(int)    { cout << "int "; }
void f(double) { cout << "double "; }
void f(long)   { cout << "long "; }
int main() { f(5); f(5.0f); f('a'); f(5L); cout << '\n'; }
```

What does this program print?

- A. `int double int int`
- B. `int int int long`
- C. `int double int long`
- D. `int double long long`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** When several functions share a name, C++ scores each one by how much it would have to **change the argument** to call it, and picks the smallest change:
1. **Exact match** (no change): best.
2. **Promotion** (a small type moved up to its natural bigger type: `char` → `int`, `float` → `double`): next best.
3. **Conversion** (any other change: `int` → `double`, `int` → `long`): third.

**Why:**

| Call | `f(int)` | `f(double)` | `f(long)` | Winner |
|---|---|---|---|---|
| `f(5)` | **exact** | conversion | conversion | int |
| `f(5.0f)` | conversion | **promotion** | conversion | double |
| `f('a')` | **promotion** | conversion | conversion | int |
| `f(5L)` | conversion | conversion | **exact** | long |

**Why the other options fail:**
- A: Misses the exact match for `5L`.
- B: Treats `float` → `int` as better than `float` → `double`.
- D: Thinks `char` → `long` beats `char` → `int`.

**Concept recap:** The full ranking ([over.ics.rank]): exact match (including trivial adjustments like array → pointer) > promotion > standard conversion > user-defined conversion > ellipsis (`...`).

**Trap alert:** `char` → `int` is a promotion, but `char` → `long` is only a conversion. "Promotion" means specifically "to `int`" (or `unsigned int`), and `float` → `double`.

**Interview angle:** *"How does the compiler pick among overloads?"* "It ranks the conversion needed for each argument: exact > promotion > conversion > user-defined. The best function must be at least as good for every argument and strictly better for at least one; otherwise the call is ambiguous."

---

## Q282. [M17 | Ambiguous calls | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P16

```cpp
void g(long)   {}
void g(double) {}
```

Select **all** calls that are **ambiguous** (compile errors).

- A. `g(5)`
- B. `g(5.0f)`
- C. `g('x')`
- D. `g(5L)`
- E. `g(5.0)`

**Answer:** A, C

**In simple words:** An ambiguity happens when **two candidates need the same "amount" of change** and neither is better. An `int` can become a `long` or a `double`, and both are equally ranked "conversions", so C++ refuses to choose. A `char` would have a promotion to `int`, but there's no `g(int)`, so it's also stuck between two conversions.

**Why each option:**

| Call | → long | → double | Result |
|---|---|---|---|
| A `g(5)` | conversion | conversion | **ambiguous** |
| B `g(5.0f)` | conversion | **promotion** | g(double) |
| C `g('x')` | conversion | conversion | **ambiguous** |
| D `g(5L)` | **exact** | conversion | g(long) |
| E `g(5.0)` | conversion | **exact** | g(double) |

**Concept recap:** Fixes: add `g(int)`, cast at the call site (`g(5L)`, `g(5.0)`), or use a template.

**Trap alert:** A "bigger" target type doesn't make a conversion better. `int` → `long` and `int` → `double` have exactly the same rank.

**Interview angle:** *"Why is `g(5)` ambiguous with `g(long)` and `g(double)`?"* "Both need a standard conversion of the same rank, so neither is better. Only exact matches and promotions outrank conversions."

---

## Q283. [M17 | A standard conversion beats a user-defined one | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T6 · P16

```cpp
#include <iostream>
#include <string>
using std::cout;
struct S { S(int) {} };
void h(S)                  { cout << "S "; }
void h(double)             { cout << "d "; }
void k(const std::string &) { cout << "str "; }
void k(bool)               { cout << "bool "; }
int main() { h(1); k("hi"); cout << '\n'; }
```

What does this program print?

- A. `d bool`
- B. `S str`
- C. `d str`
- D. Compile error: `k("hi")` is ambiguous

**Answer:** A

**Behaviour class:** Well-defined (and a notorious trap).

**In simple words:** Turning a value into **your own class type** (through a constructor) is the **least preferred** kind of conversion. A built-in conversion, even an odd one like "pointer → bool" (true if non-null), always wins over it. So `"hi"` (a `const char*`) quietly becomes `true` instead of a `std::string`.

**Why:**

| Call | Candidate | Conversion | Rank |
|---|---|---|---|
| `h(1)` | `h(S)` | int → S via a constructor | user-defined |
| | `h(double)` | int → double | **standard** → wins → `d` |
| `k("hi")` | `k(const std::string&)` | `const char*` → `std::string` via a constructor | user-defined |
| | `k(bool)` | `const char*` → bool (pointer-to-bool conversion) | **standard** → wins → `bool` |

**Why the other options fail:**
- B, C: Expect the "intended" overload. Ranking ignores intent.
- D: The ranks differ, so there's no ambiguity.

**Concept recap:** Fixes: add a `k(const char*)` overload, pass `std::string("hi")` or `"hi"s` (C++14 literal), or make `k(bool)` harder to reach (e.g. delete a `k(const char*)` overload to force a decision).

**Trap alert:** The same trap hits `std::variant<bool, std::string> v = "x";`. Before C++20's fix, it held `bool`.

**Interview angle:** *"Why does passing a string literal call my bool overload?"* "Pointer → bool is a standard conversion, which beats the user-defined conversion to `std::string`. Add a `const char*` or `std::string_view` overload."

---

## Q284. [M17 | Default arguments causing ambiguity | Medium | MCQ | [C++] · C++17 | OA] · Type T2

```cpp
void p(int a, int b = 10) {}
void p(int a) {}
int main() {
    p(1, 2);   // L1
    p(1);      // L2
}
```

What happens?

- A. Both compile; L2 calls `p(int)`, because it is an exact match without a default.
- B. L1 fails, because two-argument calls are ambiguous.
- C. Declaring both functions is itself an error.
- D. L1 compiles. L2 is **ambiguous**: both functions accept a single `int` equally well.

**Answer:** D

**Behaviour class:** Ill-formed (L2).

**In simple words:** A default argument effectively creates a **second way to call** the function with fewer arguments. `p(int a, int b = 10)` can be called as `p(1)`, and so can `p(int a)`. Both match exactly, so the call is ambiguous. Declaring both is fine; only the call that can't be decided is an error.

**Why:**

| Call | `p(int, int=10)` | `p(int)` | Result |
|---|---|---|---|
| `p(1, 2)` | ✓ exact | ✗ too many arguments | p(int, int) |
| `p(1)` | ✓ exact (uses the default) | ✓ exact | **ambiguous** |

**Why the other options fail:**
- A: A function doesn't get a better rank for "not needing a default". Defaults don't count in ranking.
- B: Only one function takes two arguments.
- C: The two declarations have different parameter lists, so they're legal overloads.

**Concept recap:** Don't combine overloading with default arguments that make the parameter counts overlap. Pick one technique.

**Trap alert:** Default arguments belong to the **declaration**. Specifying them in both a header declaration and the definition is an error, so put them in the header only.

**Interview angle:** *"Default arguments vs overloading?"* "Defaults are simpler when the behaviour is identical with a default value. Overloads are better when the behaviour differs. Never use both for the same arity."

---

## Q285. [M17 | Default argument rules | Easy | MCQ | [C++] · C++17 | OA] · Type T1

```cpp
#include <iostream>
int add(int a, int b = 2, int c = 3) { return a + b + c; }
int main() {
    std::cout << add(1) << ' ' << add(1, 5) << ' ' << add(1, 5, 7) << '\n';
}
```

What does this program print?

- A. `3 6 13`
- B. `6 9 13`
- C. `6 9 15`
- D. Compile error: default arguments must be given for every parameter

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** Default values fill in **missing arguments from the right**. You can leave off the last ones, but you can't skip a middle one.

**Why:**

| Call | a | b | c | Sum |
|---|---|---|---|---|
| `add(1)` | 1 | 2 (default) | 3 (default) | 6 |
| `add(1, 5)` | 1 | 5 | 3 (default) | 9 |
| `add(1, 5, 7)` | 1 | 5 | 7 | 13 |

**Why the other options fail:**
- A: Forgets the defaults.
- C: Miscounts the last call.
- D: Only the *trailing* parameters need defaults, and here they have them.

**Concept recap:** Once a parameter has a default, every parameter **after it** must have one too. `int bad(int a = 1, int b);` is an error ("default argument missing for parameter 2").

**Trap alert:** There's no way to give `c` but skip `b` (`add(1, , 7)` isn't valid C++). Use an options struct or overloads for that.

**Interview angle:** *"Where should default arguments be specified?"* "In the function's declaration (usually the header), only once, for trailing parameters only. And never redefine them in overrides of virtual functions (Q276)."

---

## Q286. [M17 | Overloading on `T&`, `const T&` and `T&&` | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T6

```cpp
#include <iostream>
#include <utility>
using std::cout;
void r(int &)       { cout << "lref "; }
void r(const int &) { cout << "cref "; }
void r(int &&)      { cout << "rref "; }
int main() {
    int x = 1;
    const int cx = 2;
    r(x); r(cx); r(3); r(std::move(x));
    cout << '\n';
}
```

What does this program print?

- A. `lref cref rref rref`
- B. `lref cref rref lref`
- C. `lref lref rref rref`
- D. `cref cref cref cref`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** The three reference kinds mean:
- `int&`: **"a variable I'm allowed to change"**.
- `const int&`: **"anything I may only read"**.
- `int&&`: **"a temporary, or something you've said I may take over"**.

C++ picks the most specific one that fits. `std::move(x)` doesn't move anything; it just **labels** `x` as "may be taken over".

**Why:**

| Argument | Category | `int&` | `const int&` | `int&&` | Chosen |
|---|---|---|---|---|---|
| `x` | modifiable lvalue | ✓ best | ✓ | ✗ | lref |
| `cx` | const lvalue | ✗ | ✓ | ✗ | cref |
| `3` | prvalue | ✗ | ✓ | ✓ **better** | rref |
| `std::move(x)` | xvalue | ✗ | ✓ | ✓ **better** | rref |

**Why the other options fail:**
- B: Treats `std::move(x)` as still an lvalue.
- C: Lets a const lvalue bind to `int&`. That's not allowed.
- D: Ignores the more specific overloads.

**Concept recap:** For rvalues, `T&&` beats `const T&` ([over.ics.rank]). This is how copy vs move constructors are chosen (Q295).

**Trap alert:** Without the `int&&` overload, `r(3)` would call the `const int&` version. That fallback is why a class with only a copy constructor still "works" with `std::move`: it just copies (Q297).

**Interview angle:** *"What does `std::move` do?"* "Nothing at run time. It's a cast to an rvalue reference (an xvalue), so overload resolution picks the `&&` overload, such as the move constructor."

---

## Q287. [M17 | Chained and parenthesised assignment | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P69

```cpp
#include <iostream>
using std::cout;
struct C {
    int v;
    C(int v = 0) : v(v) {}
    C &operator=(const C &o) { cout << "="; v = o.v; return *this; }
};
int main() {
    C a(1), b(2), c(3);
    a = b = c;
    c.v = 9;
    (a = b) = c;
    cout << " " << a.v << b.v << c.v << '\n';
}
```

What does this program print?

- A. `==== 333`
- B. `=== 939`
- C. `==== 939`
- D. `==== 999`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** `operator=` returns **a reference to the left-hand object**, so the result of an assignment *is* that object, ready to be used again. `a = b = c` is done right-to-left: `b = c`, then `a = (that b)`. `(a = b) = c` assigns to `a` twice: first from `b`, then from `c`.

**Why:**

| Step | Call | Prints | a | b | c |
|---|---|---|---|---|---|
| start | | | 1 | 2 | 3 |
| `a = b = c` → `b = c` | b.operator=(c) | = | 1 | 3 | 3 |
| … then `a = b` | a.operator=(b) | = | 3 | 3 | 3 |
| `c.v = 9` | | | 3 | 3 | 9 |
| `(a = b)` | a.operator=(b) | = | 3 | 3 | 9 |
| `… = c` | the returned `a&` .operator=(c) | = | **9** | 3 | 9 |

Prints `==== 939`.

**Why the other options fail:**
- A: Ignores `c.v = 9`.
- B: Miscounts the four assignments.
- D: Assigns into `b` as well.

**Concept recap:** Canonical signature: `T& operator=(const T&)`, returning `*this` by reference (Q202). Returning by value would make the second assignment in `(a = b) = c` modify a temporary.

**Trap alert:** Built-in types allow `(a = b) = c` as well. Returning `T&` keeps class types consistent with built-in ones.

**Interview angle:** *"Why does operator= return `T&`?"* "To allow chaining (`a = b = c`) with the same semantics as built-in types, and to avoid an unnecessary copy."

---

## Q288. [M17 | A self-assignment bug | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T4 · P69

```cpp
#include <cstddef>
struct Buf {
    int *p; std::size_t n;
    Buf(std::size_t n) : p(new int[n]{}), n(n) {}
    ~Buf() { delete[] p; }
    Buf &operator=(const Buf &o) {
        delete[] p;
        n = o.n;
        p = new int[n];
        for (std::size_t i = 0; i < n; ++i) p[i] = o.p[i];
        return *this;
    }
};
int main() { Buf b(3); b = b; }
```

What is wrong, and what is the smallest correct fix?

- A. Nothing; self-assignment works.
- B. It leaks memory. Fix: delete `o.p` too.
- C. Compile error: `o` is `const`.
- D. On self-assignment, `delete[] p` also frees `o.p` (they're the same array), so the loop reads freed memory, which is UB. Fix: add `if (this == &o) return *this;`, or better, allocate and copy into a new array *before* deleting the old one.

**Answer:** D

**Behaviour class:** **UB** (`b = b` → use-after-free).

**In simple words:** When an object is assigned to **itself**, "the other object" *is* this object. Throwing away my old data then also throws away the data I was about to copy from, like shredding a document before photocopying it.

> **⚠ UB box**
> - **Category:** use-after-free during self-assignment.
> - **Rule:** after `delete[] p`, `o.p` (the same pointer) points to freed storage. Reading `o.p[i]` is UB ([basic.stc.dynamic.deallocation]).
> - **What you may see in practice:** garbage values copied in, or a crash under memory checkers. **Not guaranteed.**

**Why:**

| Step | `this->p` | `o.p` (the same object) |
|---|---|---|
| `delete[] p` | freed | **freed** |
| `p = new int[n]` | new array | **now the new array too** (it's the same member) |
| copy loop | reads `o.p[i]` = the new, uninitialised array | |

Even worse, in a version that stored `o.p` in a local variable before the `delete`, that local would dangle.

**Why the other options fail:**
- A: Self-assignment is exactly the broken case.
- B: Deleting `o.p` would make things worse.
- C: The code compiles.

**Concept recap:** An exception-safe and self-assignment-safe pattern:
```cpp
Buf &operator=(const Buf &o) {
    int *np = new int[o.n];               // 1. may throw: the object is unchanged
    std::copy(o.p, o.p + o.n, np);         // 2. copy
    delete[] p; p = np; n = o.n;           // 3. commit
    return *this;
}
```
Or copy-and-swap (Q300).

**Trap alert:** The `if (this == &o)` check fixes self-assignment but **not** exception safety. If `new` throws after the `delete[]`, the object is left with a dangling `p`.

**Interview angle:** *"What are the pitfalls in writing operator=?"* "Self-assignment, exception safety (allocate before releasing), deep copying, returning `*this`. The copy-and-swap idiom handles all of them."

---

## Q289. [M17 | Arithmetic and comparison operators | Easy | MCQ | [C++] · C++17 | OA] · Type T1

```cpp
#include <iostream>
struct V { int x, y; };
V operator+(V a, V b) { return {a.x + b.x, a.y + b.y}; }
bool operator==(V a, V b) { return a.x == b.x && a.y == b.y; }
bool operator<(V a, V b)  { return a.x < b.x || (a.x == b.x && a.y < b.y); }
int main() {
    V p{1, 2}, q{1, 3};
    V s = p + q;
    std::cout << s.x << s.y << ' ' << (p == q) << (p < q) << (q < p) << (s == V{2, 5}) << '\n';
}
```

What does this program print?

- A. `25 0110`
- B. `25 0101`
- C. `25 1101`
- D. `34 0101`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** Operator overloading lets your own types use `+`, `==` and `<` like numbers do. Each operator is just a function with a special name. `operator<` here compares "first by x, then by y", like sorting words letter by letter.

**Why:**

| Expression | Evaluation | Value |
|---|---|---|
| `p + q` | {1 + 1, 2 + 3} | {2, 5} → prints `25` |
| `p == q` | 2 ≠ 3 | 0 |
| `p < q` | x equal, 2 < 3 | 1 |
| `q < p` | x equal, 3 < 2 is false | 0 |
| `s == V{2, 5}` | equal | 1 |

**Why the other options fail:**
- A: Swaps the two `<` results.
- C: Thinks `p == q`.
- D: Miscomputes the sum.

**Concept recap:** Non-member operators (as here) treat both operands symmetrically. They also allow implicit conversions on the **left** operand, which member operators don't.

**Trap alert:** If you define `<`, users expect `>`, `<=` and `>=` too. In C++20, `operator<=>` (the spaceship operator) and a defaulted `==` generate all six comparisons.

**Interview angle:** *"Member or non-member operator+?"* "Non-member (often a friend), for symmetry: `2 + v` works only if `+` is a non-member and `int` converts to V. Implement `+` in terms of a member `+=`."

---

## Q290. [M17 | Why `operator<<` must be a non-member | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T2

```cpp
#include <iostream>
struct P {
    int v;
    std::ostream &operator<<(std::ostream &os) const { return os << v; }
};
int main() {
    P p{5};
    std::cout << p;      // L1
    p << std::cout;      // L2
}
```

Which statement is correct?

- A. Both lines compile and print `5`.
- B. Only L1 compiles.
- C. Only L2 compiles (and prints `5`). L1 fails, because a **member** `operator<<` makes `P` the **left** operand. The fix is a non-member `std::ostream &operator<<(std::ostream &, const P &)`, usually a friend.
- D. Neither compiles.

**Answer:** C

**Behaviour class:** Ill-formed (L1).

**In simple words:** For a member operator, **the object it belongs to is always on the left**: `p << x` means `p.operator<<(x)`. But in `std::cout << p`, `cout` is on the left, and you can't add a member to `std::ostream`. So the operator has to be a free function that takes the stream first and your object second.

**Why:**

| Line | Meaning | Is there a matching function? |
|---|---|---|
| L1 `cout << p` | `cout.operator<<(p)` or `operator<<(cout, p)` | ✗ neither exists → "no match for operator<<" |
| L2 `p << cout` | `p.operator<<(cout)` | ✓ prints 5, but reads backwards |

**Why the other options fail:**
- A, B: Assume the member works with the stream on the left.
- D: L2 does match the member.

**Concept recap:** The idiomatic form:
```cpp
struct P {
    int v;
    friend std::ostream &operator<<(std::ostream &os, const P &p) { return os << p.v; }
};
```
It returns the stream by reference for chaining (`cout << a << b`), and takes `const P&`.

**Trap alert:** Defining it as a *hidden friend* (inside the class, as above) means it's found only through argument-dependent lookup, which keeps overload sets small and error messages short.

**Interview angle:** *"Why can't operator<< be a member of my class?"* "The left operand of a member operator is always the class object, but in `cout << obj` the left operand is the stream. So it must be a non-member, often a friend for private access, returning `std::ostream&`."

---

## Q291. [M17 | A conversion operator that creates ambiguity | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T6

```cpp
#include <iostream>
struct Frac {
    int n, d;
    Frac(int n, int d = 1) : n(n), d(d) {}
    operator double() const { return double(n) / d; }
};
int main() {
    Frac f(1, 2);
    double x = f + 1;
    Frac g = 3;
    std::cout << x << ' ' << double(g) << ' ' << (f < 1) << '\n';
}
```

What does it print? And what happens to `f + 1` if a non-member `Frac operator+(Frac, Frac)` is **also** defined?

- A. It prints `1.5 3 1`. With the extra `operator+`, `f + 1` becomes **ambiguous** (a compile error).
- B. It prints `1.5 3 1`. With the extra `operator+`, `f + 1` calls it, giving the fraction 3/2.
- C. It prints `1 3 1`. The extra `operator+` changes nothing.
- D. It doesn't compile even without the extra `operator+`.

**Answer:** A

**Behaviour class:** Well-defined as written. Ill-formed (ambiguous) once `operator+(Frac, Frac)` is added.

**In simple words:** `Frac` can turn **into** a double (the conversion operator), and an `int` can turn **into** a `Frac` (the constructor). With only the conversion to `double` available, `f + 1` becomes ordinary `0.5 + 1`. Add a `Frac + Frac` operator, and there are **two equally good routes**: convert `f` to double, or convert `1` to `Frac`. Each needs one user-defined conversion on a different argument, so neither is better, and the call is ambiguous.

**Why:**

| Expression | Route | Result |
|---|---|---|
| `f + 1` | built-in `double + int` via `operator double` | 1.5 |
| `Frac g = 3` | converting constructor `Frac(3, 1)` | `double(g)` = 3 |
| `f < 1` | 0.5 < 1 via `operator double` | 1 |

With `operator+(Frac, Frac)` added:

| Candidate | Argument 1 (`f`) | Argument 2 (`1`) |
|---|---|---|
| built-in `+(double, int)` | user-defined (Frac → double) | exact |
| `operator+(Frac, Frac)` | exact | user-defined (int → Frac) |

Each is better on one argument and worse on the other → **ambiguous** (verified: "ambiguous overload for 'operator+'").

**Why the other options fail:**
- B: Assumes the user's operator has priority. It doesn't.
- C: The division is `double(n) / d`, so the result is 0.5, not integer division.
- D: The code compiles as written.

**Concept recap:** The fix is to make one direction `explicit`: `explicit operator double()` (then you must write `double(f)`), or an `explicit` constructor.

**Trap alert:** Types that convert both *to* and *from* another type in both directions almost always cause ambiguities like this one.

**Interview angle:** *"Why should conversion operators usually be explicit?"* "Implicit conversions in both directions cause ambiguities and surprising overload choices. Make one or both explicit, and provide named functions like `to_double()`."

---

## Q292. [M17 | Rules for implicit conversions | Medium | MSQ | [C++] · C++17 | Interview] · Type — concept

Select **all** true statements.

- A. A non-explicit constructor callable with one argument acts as an implicit conversion.
- B. At most **one** user-defined conversion is applied implicitly in a single conversion sequence.
- C. `explicit operator int() const;` still allows `int x = obj;`.
- D. `explicit operator bool() const;` still allows `if (obj)`.
- E. A conversion operator is declared with a return type before `operator`, as in `int operator int();`.

**Answer:** A, B, D

**In simple words:** C++ will make **one** silent "custom" conversion for you, but it won't chain two together (for example `const char*` → `std::string` → `MyType`). `explicit` turns off the silent conversions, except in yes/no conditions for `operator bool`.

**Why each option:**
- **A ✓:** E.g. `Frac g = 3;` (Q291). A constructor with defaults for the remaining parameters counts too.
- **B ✓:** [class.conv]p4. If `MyType` has a constructor from `std::string`, then `MyType m = "hi";` fails, because it would need two user-defined steps.
- **C ✗:** `explicit` blocks copy-initialisation. You must write `int x = static_cast<int>(obj);` or `int x(obj);`.
- **D ✓:** Contextual conversion to bool (Q220).
- **E ✗:** The return type is the operator's name itself: `operator int() const;`. Writing a return type before it is ill-formed.

**Concept recap:** User-defined conversions can come from **converting constructors** (into your type) or **conversion functions** (out of your type). Both can be `explicit` (conversion functions since C++11).

**Trap alert:** Option B is why `std::string s = 'x';` doesn't compile, but `std::string s(1, 'x');` does. Implicit-conversion chains have strict limits.

**Interview angle:** *"How many user-defined conversions can happen implicitly?"* "One per conversion sequence, surrounded by standard conversions if needed (e.g. `short` → `int` → MyType)."

---

## Q293. [M18 | Which expressions are lvalues? | Easy | MSQ | [C++] · C++17 | OA+Interview] · Type T3

Given `int x; int &f(); int g();`, select **all** expressions that are **lvalues**.

- A. `x`
- B. `x + 1`
- C. `f()`
- D. `g()`
- E. `++x`

**Answer:** A, C, E

**In simple words:** An **lvalue** is an expression that names **an object with a location you can point to** (you can take its address and, unless it's `const`, assign to it). A **prvalue** is a **pure value** like `x + 1` or a returned copy: a result with no lasting home of its own.

**Why each option:**

| Expression | Category | Why |
|---|---|---|
| `x` | lvalue | names a variable (`&x` is valid) |
| `x + 1` | prvalue | a computed value |
| `f()` | lvalue | it returns `int&`, i.e. it refers to an existing object |
| `g()` | prvalue | returns `int` by value |
| `++x` | lvalue | prefix `++` returns the variable itself in C++ (Q29 noted that C differs) |

**Concept recap:** The simple test: **can you take its address with `&`?** If yes, it's an lvalue. (Bit-fields and a few edge cases aside.)

**Trap alert:** `x++` (postfix) is a **prvalue**: it returns a copy of the old value.

**Interview angle:** *"What's an lvalue vs an rvalue?"* "An lvalue has identity (an address) and persists beyond the expression. An rvalue is a temporary value: prvalues are pure values, and xvalues are expiring objects. Rvalues can bind to `T&&`, which is what enables move semantics."

---

## Q294. [M18 | xvalues and prvalues | Hard | MSQ | [C++] · C++17 | Interview] · Type T3

Given `std::string s;`, select **all** true statements.

- A. `std::move(s)` is an **xvalue**.
- B. `std::string("a")` is a **prvalue**.
- C. `s + "x"` is an **lvalue**.
- D. After `std::string &&r = std::move(s);`, the expression `r` is an **rvalue**.
- E. `static_cast<std::string &&>(s)` is an **xvalue**.

**Answer:** A, B, E

**In simple words:** There are three kinds of expression:
- **lvalue:** a named object that's staying around.
- **prvalue:** a brand-new value, like a freshly built temporary.
- **xvalue:** an **existing object marked as "about to expire, feel free to take its stuff"**. That marking is exactly what `std::move` does.

A surprising rule: **a named rvalue reference is itself an lvalue**, because it has a name, so it can be used again.

**Why each option:**
- **A ✓:** `std::move` returns `std::string&&`. A function call returning an rvalue reference is an xvalue.
- **B ✓:** A temporary created by a functional cast is a prvalue.
- **C ✗:** `operator+` returns `std::string` by value → a prvalue.
- **D ✗:** `r` has a name → it's an **lvalue**, even though its type is `std::string&&`. To move from it again, write `std::move(r)` (Q296).
- **E ✓:** That's literally what `std::move` does inside.

**Concept recap:**

```
              expression
              /        \
        glvalue        rvalue
        /     \        /     \
   lvalue     xvalue        prvalue
```

- glvalue = has identity (lvalue or xvalue).
- rvalue = can be moved from (xvalue or prvalue).

**Trap alert:** "Type" and "value category" are different properties. `r` has type `std::string&&` but category lvalue.

**Interview angle:** *"Is a variable of type `T&&` an rvalue?"* "No. The expression naming it is an lvalue. That's why, inside a function taking `T&& x`, you must write `std::move(x)` (or `std::forward<T>(x)` in templates) to pass it on as an rvalue."

---

## Q295. [M18 | Copy or move constructor? | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P59

```cpp
#include <iostream>
#include <utility>
using std::cout;
struct T {
    T() {}
    T(const T &)     { cout << "C"; }
    T(T &&) noexcept { cout << "M"; }
};
int main() {
    T a;
    T b = a;
    T c = std::move(a);
    T d = T();
    const T e;
    T f = std::move(e);
    cout << '\n';
}
```

What does this program print?

- A. `CMM`
- B. `CCC`
- C. `CM`
- D. `CMC`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:**
- **Copy** when the source is a named object you'll keep using (`b = a`).
- **Move** when the source is marked "expiring" (`std::move(a)`).
- **Do nothing** when building straight from a fresh temporary (`T()`: C++17 elision).

A `const` object can't really be moved from (moving means changing the source), so `std::move` on a const object quietly falls back to a **copy**.

**Why:**

| Line | Source category | `T(const T&)` | `T(T&&)` | Chosen | Prints |
|---|---|---|---|---|---|
| `T b = a` | lvalue | ✓ | ✗ | copy | C |
| `T c = std::move(a)` | xvalue `T&&` | ✓ | ✓ better | move | M |
| `T d = T()` | prvalue | — | — | **elided** (C++17) | — |
| `T f = std::move(e)` | xvalue of type **`const T&&`** | ✓ (binds to `const T&`) | ✗ (`T&&` can't bind const) | copy | C |

Prints `CMC`.

**Why the other options fail:**
- A: Assumes moving from a const object moves.
- B: Ignores the move constructor.
- C: Misses the const case.

**Concept recap:** `std::move` on a `const` object is a silent pessimisation: it compiles and then copies. Don't mark objects `const` if you plan to move from them, e.g. `const std::string s; return std::move(s);` copies.

**Trap alert:** A `T(const T&&)` constructor would match, but is almost never written, because it can't actually steal the resources.

**Interview angle:** *"When is the move constructor called?"* "When initialising from an rvalue (a prvalue that isn't elided, or an xvalue from `std::move` or a function returning `T&&`) of non-const T, and the class has an accessible move constructor."

---

## Q296. [M18 | A named rvalue reference is an lvalue | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P59

```cpp
#include <iostream>
#include <utility>
using std::cout;
struct T {
    T() {}
    T(const T &)     { cout << "C"; }
    T(T &&) noexcept { cout << "M"; }
};
void sink(T &&x) {
    T y = x;
    T z = std::move(x);
}
int main() { sink(T()); cout << '\n'; }
```

What does this program print?

- A. `MM`
- B. `CM`
- C. `CC`
- D. `MC`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** `x` was **bound** to a temporary, but inside the function it has a **name**, and anything with a name is treated as "still in use", i.e. an lvalue. So `T y = x;` copies. You must explicitly say `std::move(x)` to hand it over. This protects you from accidentally gutting a variable you might still need two lines later.

**Why:**

| Line | Expression | Category | Constructor |
|---|---|---|---|
| `sink(T())` | the temporary binds to `x` | — | none (just a reference binding) |
| `T y = x` | `x` (named) | **lvalue** | copy → C |
| `T z = std::move(x)` | xvalue | rvalue | move → M |

**Why the other options fail:**
- A: Thinks `x` is an rvalue because of its type.
- C: Thinks `std::move` doesn't trigger a move.
- D: Swaps them.

**Concept recap:** The rule: **"if it has a name, it's an lvalue."** Move explicitly, and only at the **last** use.

**Trap alert:** After `std::move(x)`, `x` is in a "valid but unspecified" state (for standard types). Don't read its value, but you may assign to it or destroy it.

**Interview angle:** *"Why doesn't a parameter of type `T&&` automatically move?"* "Because it has a name, so it's an lvalue. Automatic moving would be dangerous if the parameter were used again. Use `std::move` (or `std::forward` for forwarding references)."

---

## Q297. [M18 | When the compiler doesn't generate a move constructor | Hard | MCQ | [C++] · C++17 | Interview] · Type T7 · P59

```cpp
#include <iostream>
#include <utility>
using std::cout;
struct A {
    A() {}
    A(const A &) { cout << "C"; }      // user-declared copy constructor
};
struct W { A a; };                      // no user-declared special members
int main() {
    A x;
    A y = std::move(x);
    W w1;
    W w2 = std::move(w1);
    cout << '\n';
}
```

What does this program print?

- A. `MM`
- B. `CM`
- C. `CC`
- D. Compile error: `A` has no move constructor

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** If you write your own **copy constructor** (or copy assignment, or destructor), C++ assumes your class manages something special, and **stops generating a move constructor** for you. Then a "move" falls back to the copy constructor, which still works, because `const A&` accepts an rvalue. `W`'s generated move constructor moves each member, but moving the member `a` also ends up copying.

**Why:**

| Line | Needs | What's available | Prints |
|---|---|---|---|
| `A y = std::move(x)` | A's move constructor | **not declared** (suppressed by the user copy constructor) → copy constructor binds the rvalue | C |
| `W w2 = std::move(w1)` | W's implicit move constructor → moves member `a` | A has no move constructor → copy | C |

**Why the other options fail:**
- A, B: Assume A gets an implicit move constructor.
- D: There's no error, because a "move" quietly becomes a copy.

**Concept recap:** When the compiler generates special members:

| You declare… | Implicit move constructor / move assignment? |
|---|---|
| nothing | ✓ generated |
| a copy constructor, copy assignment, or destructor | ✗ not declared (moves fall back to copies) |
| a move operation | copy operations become **deleted** |

That's the reasoning behind the **rule of five / zero** (P70, Batch 16).

**Trap alert:** Adding an empty destructor "for debugging" (`~X() {}`) silently turns every move of `X` into a copy. That's a real performance bug in large codebases.

**Interview angle:** *"When does the compiler generate a move constructor?"* "Only if the class has no user-declared copy constructor, copy assignment, move assignment or destructor. Otherwise, rvalue copies use the copy constructor. Write `= default` for the move operations if you need them back."

---

## Q298. [M18 | Shallow copy of an owning class | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T3 · P70

> **This code has undefined behaviour.**

```cpp
#include <iostream>
struct S {
    int *p;
    S(int v) : p(new int(v)) {}
    ~S() { delete p; }
};
int main() {
    S a(5);
    S b = a;
    *b.p = 7;
    std::cout << *a.p << '\n';
}
```

Which statement is correct?

- A. It prints `5`, because `b` has its own copy.
- B. It prints `7`, and the program is well-defined.
- C. Compile error: `S` has no copy constructor.
- D. The behaviour is undefined: the implicitly generated copy constructor copies the **pointer** (a shallow copy), so both `a` and `b` own the same `int`, and both destructors `delete` it at the end of `main` (a double delete).

**Answer:** D

**Behaviour class:** **Undefined behaviour.**

**In simple words:** The compiler's automatic copy just copies each member's value, and for a pointer that means **copying the address, not the thing it points to**. Now two objects think they own the same `int`. Changing it through one shows up in the other, and when both are destroyed, the same memory is freed twice.

> **⚠ UB box**
> - **Category:** double delete (caused by a shallow copy of an owning raw pointer).
> - **Rule:** [expr.delete]: deleting a pointer whose storage has already been deallocated is UB.
> - **What you may see in practice:** it often prints `7` (shared data), then may crash or report a "double free" error as `main` ends. **Not guaranteed.** Since the program as a whole has UB, even the printed `7` isn't guaranteed.

**Why:**

| Step | a.p | b.p | The `int` |
|---|---|---|---|
| `S a(5)` | → X | | X = 5 |
| `S b = a` (memberwise copy) | → X | → **X** | |
| `*b.p = 7` | | | X = 7 |
| end of `main`: `~b`, then `~a` | | | `delete` X **twice** |

**Why the other options fail:**
- A: Describes a deep copy, which this class doesn't do.
- B: Ignores the double delete at the end.
- C: The copy constructor is implicitly generated.

**Concept recap:** **Rule of three:** if a class needs a user-written destructor (because it owns a resource), it almost certainly needs a user-written copy constructor and copy assignment too (Q299).

**Trap alert:** It's the C++ version of Q162 (C's shallow struct copy). The difference is that C++ runs destructors automatically, so the double free is guaranteed to be attempted.

**Interview angle:** *"What happens if a class with a raw owning pointer has no copy constructor?"* "The default copy is shallow, which leads to shared ownership, a double delete, and use-after-free after one copy is destroyed. Fix: the rule of three/five, or better, the rule of zero with `std::unique_ptr` or `std::vector` members."

---

## Q299. [M18 | Fixing the shallow copy: the rule of three | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T4 · P70

The class from Q298 (`S` owns `int *p`, and its destructor does `delete p`) must support correct copying. What is the **minimal** set of members to add?

- A. A deep **copy constructor** (`S(const S &o) : p(new int(*o.p)) {}`) **and** a deep **copy assignment** operator. Together with the existing destructor, that's the rule of three.
- B. Only a deep copy constructor; assignment is fine as generated.
- C. `S(const S &) = default;` to make copying explicit.
- D. Only a change to the destructor: check `if (p)` before deleting.

**Answer:** A

**Behaviour class:** — (a design fix).

**In simple words:** There are **two** ways an object gets someone else's contents: being **born** as a copy (the copy constructor) and being **overwritten** by assignment (copy assignment). The generated versions of both copy the pointer. So you need to fix **both**, and keep the destructor, which is the "three".

**Why:**

| Option | Fixes `S b = a;`? | Fixes `b = a;`? | Result |
|---|---|---|---|
| A | ✓ | ✓ | correct |
| B | ✓ | ✗ (assignment still shares the pointer, and leaks b's old int) | broken |
| C | ✗ (`= default` is still shallow) | ✗ | broken |
| D | ✗ (`p` isn't null; it's the same non-null pointer twice) | ✗ | broken |

A correct copy assignment, for example:
```cpp
S &operator=(const S &o) { *p = *o.p; return *this; }  // copies the value; both ints already exist
```

**Concept recap:**
- **Rule of three (C++98):** destructor + copy constructor + copy assignment.
- **Rule of five (C++11):** plus move constructor + move assignment.
- **Rule of zero:** use members that manage themselves (`std::unique_ptr<int>`, `std::vector`) and declare **none** of the five.

**Trap alert:** Option D is a common misconception: checking for null doesn't prevent double deletes, because the pointer isn't null; both copies hold the same valid address until the first `delete`.

**Interview angle:** *"Explain the rule of three/five/zero."* "If you manage a resource manually, define (or delete) all of destructor, copy constructor, copy assignment, plus the move operations. Better: use the rule of zero, so RAII members handle everything."

---

## Q300. [M18 | Copy-and-swap | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P69

```cpp
#include <algorithm>
#include <cstddef>
#include <iostream>
#include <utility>
using std::cout;
struct Buf {
    int *p; std::size_t n;
    Buf(std::size_t n) : p(new int[n]{}), n(n) {}
    Buf(const Buf &o) : p(new int[o.n]), n(o.n) { cout << "C"; std::copy(o.p, o.p + n, p); }
    ~Buf() { delete[] p; }
    friend void swap(Buf &a, Buf &b) noexcept { std::swap(a.p, b.p); std::swap(a.n, b.n); }
    Buf &operator=(Buf o) { cout << "="; swap(*this, o); return *this; }
};
int main() {
    Buf a(2), b(3);
    a = b;
    a = a;
    cout << ' ' << a.n << '\n';
}
```

What does this program print?

- A. `C= 3` (self-assignment is skipped)
- B. `C=C= 3`
- C. `=C= 2`
- D. `C=C= 2`

**Answer:** B

**Behaviour class:** Well-defined (and self-assignment-safe).

**In simple words:** **Copy-and-swap** does assignment in two easy, safe steps:
1. **Make a full copy** of the source (the by-value parameter). If that fails, nothing has been touched yet.
2. **Swap** our contents with the copy. When the copy is destroyed at the end of the function, it takes our **old** data with it.

Self-assignment needs no special case: you just make a copy of yourself and swap with it.

**Why:**

| Statement | Step | Prints | Effect |
|---|---|---|---|
| `a = b` | parameter `o` is copy-constructed from b | C | o = copy of b (n = 3) |
| | body | = | swap: a gets b's copy; o gets a's old buffer |
| | `o` destroyed | | a's old buffer freed |
| `a = a` | `o` copy-constructed from a | C | a fresh copy of a |
| | body | = | swap (a now owns the copy) |
| | `o` destroyed | | a's previous buffer freed |
| print | | ` 3` | a.n = 3 |

**Why the other options fail:**
- A: Copy-and-swap has no self-check; it always copies.
- C: Puts the `=` before the copy. The parameter is created before the body runs.
- D: Wrong final size: a now holds b's size, 3.

**Concept recap:** Benefits of copy-and-swap:
- **strong exception safety** (all-or-nothing);
- automatic self-assignment safety;
- reuses the copy constructor and destructor (no duplicated logic);
- because the parameter is taken by value, an rvalue argument is **moved** into `o`, so this one `operator=` handles both copy and move assignment.

**Trap alert:** The cost: self-assignment and some assignments always pay for a full copy, even when the existing buffer could have been reused. Performance-critical containers like `std::vector` don't use this idiom.

**Interview angle:** *"Explain the copy-and-swap idiom."* "`operator=(T other)` takes the parameter by value (copy or move happens at the call), swaps members with it using a `noexcept` swap, and the old state is destroyed with the parameter. It gives the strong exception guarantee and handles self-assignment."

---

## Batch 15 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 281 | M17 | E | MCQ | C | Well-defined |
| 282 | M17 | M | MSQ | A, C | Ill-formed (ambiguous) / OK |
| 283 | M17 | H | MCQ | A | Well-defined |
| 284 | M17 | M | MCQ | D | Ill-formed (L2) |
| 285 | M17 | E | MCQ | B | Well-defined |
| 286 | M17 | M | MCQ | A | Well-defined |
| 287 | M17 | M | MCQ | C | Well-defined |
| 288 | M17 | H | MCQ | D | **UB** (self-assignment) |
| 289 | M17 | E | MCQ | B | Well-defined |
| 290 | M17 | H | MCQ | C | Ill-formed (L1) |
| 291 | M17 | H | MCQ | A | Well-defined / ambiguous variant |
| 292 | M17 | M | MSQ | A, B, D | — (rules) |
| 293 | M18 | E | MSQ | A, C, E | — (value categories) |
| 294 | M18 | H | MSQ | A, B, E | — (value categories) |
| 295 | M18 | H | MCQ | D | Well-defined |
| 296 | M18 | M | MCQ | B | Well-defined |
| 297 | M18 | H | MCQ | C | Well-defined |
| 298 | M18 | M | MCQ | D | **UB** (double delete) |
| 299 | M18 | M | MCQ | A | — (design fix) |
| 300 | M18 | H | MCQ | B | Well-defined |

- **Difficulty:** 4 E / 8 M / 8 H.
  - M17 total: 3/5/4, as planned.
  - M18 so far: 1/3/4.
- **Format:** 16 MCQ / 4 MSQ.
  - M17 total: 10/2, as planned.
- **MCQ answer keys:** A×4, B×4, C×4, D×4. Running total Q1–Q300: **A 60 · B 60 · C 60 · D 60**.

---

**Completed: 300 / 400 Remaining: 100**

**Ledger, newly ✓ this batch:**
- P16 ×3 (Q281, Q282, Q283)
- P69 ×3 (Q287, Q288, Q300) ★ complete
- P59 ×3 (Q295, Q296, Q297) ★ complete
- P70 ×2 (Q298, Q299) ★ complete (+1 planned)

**M17 complete** (Q281–Q292, 12 questions). **M18 in progress** (8 of 18 done).

---
---

# SELF-AUDIT #3 (after Batch 15, Q1–Q300)

## A. Actual vs planned splits

| | Easy | Medium | Hard |
|---|---|---|---|
| Actual, Q1–Q300 | 77 (25.7%) | 123 (41.0%) | 100 (33.3%) |
| Plan for exactly the modules covered (M1–M17 + 8 of M18) | 77 | 123 | 100 |
| Bank-wide target | 25% | 40% | 35% |

✅ An exact match with the per-module plan. The remaining 100 questions are planned at 23 E / 37 M / 40 H, which lands the bank on **100 / 160 / 140**.

- **Format:** MCQ 240 / MSQ 60 = exactly 80 / 20. ✅
- **Language:** [C] 184 / [C++] 116. 11 C questions remain (M24, M25), for a final 195 / 205. ✅
- **MCQ answer keys:** A 60 · B 60 · C 60 · D 60. ✅ "All/none of the above" used: **0**.
- **Code-reading share:** about 83%. ✅

## B. Relevance tags

**Correction found in this audit:** Q267 (a C++ question) was tagged **[All]**, which implies GATE relevance. That's wrong, because GATE doesn't test C++. It has been retagged **[OA+Interview]**, and the file was updated. No other C++ question carries a GATE tag (checked automatically).

| Tag (C++ only, Q185–Q300) | Actual | Pace needed for the revised C++ targets (≈ 62 / 101 / 42 over 205) |
|---|---|---|
| [OA] | 26 | 35 |
| [OA+Interview] | 71 | 57 |
| [Interview] | 19 | 24 |

⚠ The OOP modules naturally scored "OA+Interview". The remaining C++ modules are well suited to rebalance: STL and template output questions → [OA]; modern C++ and concurrency concepts → [Interview].

**Plan for the remaining 89 C++ questions:** about 36 [OA], 30 [OA+Interview], 23 [Interview]. That gives final C++ figures of about 62 / 101 / 42, on target.

The C split is unchanged since Audit #2. The remaining 11 C questions will be tagged [GATE] or [GATE+OA] only.

## C. Pointer + OOP share

M6 (35) + M13–M18 so far (14 + 24 + 16 + 26 + 12 + 8 = 100) = 135 / 300. The final total will be 145 / 400 = **36.25%**. ✅

## D. Full Pattern Ledger

| Range | Status |
|---|---|
| P01–P15 | all ✓ |
| P16 | ✓ 3 |
| P17–P50 | all ✓ |
| P51 ✓2 · P52 ✓4 · P53 ✓2 · P54 ✓2 · P55 ✓1 · P56 ✓2 | all ✓ |
| P57 ✓3 · P58 ✓3 · P59 ✓3 · P60 ✓3 · P61 ✓2 · P62 ✓6 · P63 ✓7 · P64 ✓2 · P65 ✓3 · P66 ✓3 · P67 ✓2 · P68 ✓3 · P69 ✓3 · P70 ✓2 · P71 ✓3 · P72 ✓3 | all ✓ (OOP family complete) |
| P73–P77 smart pointers, move, `std::forward` | ☐ (P76, P77 in Batch 16; P73–P75 in M22) |
| P78 auto ★ | ✓ 1 — needs 1 more (M22) |
| P79, P80 lambdas | ☐ (M21/M22) |
| P81 ✓1 · P82 ✓1 | ✓ |
| P83–P96 templates, STL, concurrency | ☐ (M20–M23, M25) |
| P97 ✓2 · P98 ✓1 · P99 ✓1 · P100 ✓1 · P101 ✓1 | ✓ |
| P102–P105 integrated | ☐ (M20, M25) |

**Totals:** **79 / 105 patterns ✓. ★ patterns: 41 / 51 complete** (every ★ up to P72 is done). The remaining 10 ★ patterns are P73, P74, P75, P76, P78 (1 of 2), P79, P83, P86, P87 and P102. Every one is scheduled in Batches 16–20. *(A small inconsistency was corrected: Batch 14's note "P63 → 8" was a miscount; the real count is 7. It's still complete.)*

## E. Re-trace of every output; corrections

- All output questions in Batches 11–15 were compiled as C++17 and run, and every key matches.
- Every compile-error, link-error and abort case was reproduced: Q209, Q214, Q216, Q219, Q220, Q230, Q231, Q235, Q239, Q240, Q245, Q247, Q253, Q261, Q269, Q271, Q274, Q282, Q284, Q290, Q291.
- **Corrections made:**
  - Batch 12, Q226 and Q228: answer letters rebalanced (option order changed, keys updated). The content was correct.
  - Batch 13, Q249: a malformed tag was fixed.
  - Batch 14, Q267: tag changed from [All] to [OA+Interview] (this audit).
- **No answer key needed changing.**

## F. ⚠ Highest-priority check: deterministic output claimed for UB code?

- **UB questions in Batches 11–15** (each with "undefined" as the correct option, plus a general "what you may see" note): Q222, Q237 (the bad order), Q267, Q274, Q288, Q298, and the UB options in Q229, Q253 and Q268.
- **Code that looks risky but was re-checked as well-defined:**
  - Q236: `make2`'s output is described as *unspecified* between two forms, not fixed.
  - Q258 and Q279: virtual destructors make the deletes through base pointers correct.
  - Q278: partial assignment is defined behaviour.
  - Q300: self-assignment through copy-and-swap is safe.

**Result: 0 violations.**

---

# PROGRESS REPORT @ 300

**Modules complete:** M1–M17. The entire C language part, C++ fundamentals, and the whole OOP family except the second half of M18. **Next:** M18 (move semantics, rule of five), M19 exceptions, M20 templates, M21 STL, M22 modern C++, M23 threads, M24 data structures in C/C++, M25 integrated.

**Key OOP traps covered in Q201–Q300:**
1. `const` member overloads chosen by **how** you access the object (Q208); `mutable` (Q207).
2. Members are initialised in **declaration order** (Q221–Q222); bases come before members, and destruction is reversed (Q223, Q243).
3. Destruction order for locals, statics and globals (Q225); temporaries die at the `;` unless bound **directly** to a reference (Q228–Q229).
4. `W b();` declares a function (Q217). C++17 guaranteed elision (Q218, Q230, Q236), NRVO is optional, and returning a parameter moves it (Q238).
5. The static-initialisation-order fiasco (Q237).
6. Name hiding and `using` (Q245–Q246); lookup ambiguity happens before overload resolution (Q248).
7. Diamond: two copies vs one virtual base; the most-derived class constructs the virtual base (Q249–Q251).
8. Virtual dispatch and vtables (Q257–Q258); a missing `const` means no override (Q259); private overrides still dispatch (Q260).
9. Overloads are chosen by the **static** type (Q263–Q264); default arguments come from the static type too (Q276).
10. A non-virtual destructor makes `delete` through a base pointer UB (Q267); `shared_ptr` vs `unique_ptr` (Q268).
11. Virtual calls in constructors and destructors don't reach derived classes, and a pure virtual call there is UB (Q272–Q274).
12. Slicing and partial assignment (Q277–Q278).
13. `"hi"` → `bool` beats `"hi"` → `std::string` (Q283); conversion operators can make operators ambiguous (Q291).
14. A const object "moved" is actually copied (Q295); a named `T&&` is an lvalue (Q296); a user-declared copy constructor suppresses the implicit move (Q297).

**UB categories now covered (25):** the 21 from Audit #2, plus:
22. reading members in the wrong initialisation order
23. the static-initialisation-order fiasco
24. `delete` through a base pointer with a non-virtual destructor
25. a pure virtual call during construction

(Double deletes from shallow copies and self-assignment use-after-free fall under the earlier categories.)

**Still to come:** iterator invalidation, lambda capture lifetimes, data races, `string_view` dangling, moved-from object misuse.

**Areas to revise:** Q225 (static and global destruction order), Q229 (lifetime extension), Q246/Q248 (lookup vs overloading), Q251 (virtual base construction), Q263 (virtual overload hiding), Q273/Q274 (constructor dispatch), Q283/Q291 (conversion ranking), Q294–Q297 (value categories and implicit moves).

*Say **"Continue."** for Batch 16 (Q301–Q320: the rest of M18 — move constructor and assignment, moved-from state, `std::move` vs `std::forward`, rule of three/five/zero, `noexcept` moves with `vector`, RAII wrappers — then M19 exceptions, then the start of M20 templates).*
