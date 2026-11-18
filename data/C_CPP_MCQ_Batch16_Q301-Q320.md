# C/C++ MCQ/MSQ Bank — Batch 16 (Q301–Q320)

**What this batch covers:**
- **M18 Copy/move, finishing (Q301–Q310):**
  - writing a move constructor and a move assignment operator
  - what `std::move` really does
  - `std::forward` vs `std::move` vs no cast at all
  - the rule of five: declaring a move operation deletes the copy operations
  - the rule of zero with a `unique_ptr` member
  - `return std::move(local)` disables NRVO
  - `noexcept` move constructors and `vector` reallocation
  - RAII guards during an exception
  - what you may do with a moved-from object
- **M19 Exceptions (Q311–Q318):**
  - catching by value vs by reference
  - the order of catch clauses
  - destructors running during stack unwinding
  - rethrowing from a nested handler
  - a constructor that throws
  - `noexcept` violations
  - `throw;` vs `throw e;`
  - the exception-safety guarantees
- **M20 Templates, starting (Q319–Q320):** template argument deduction, conflicting deductions, and deduction with arrays and `const`.

**Patterns:** P59, P70, P71, P76, P77, P83

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - Traces show every copy/move/destructor call and why it was chosen.

> **Conventions for every question in this batch:**
> - **C++17**, exceptions enabled, with guaranteed copy elision from prvalues.
> - Platform: a typical 64-bit system (`int` 4 bytes, pointers 8).
> - Needed headers are included. In code, `cout` means `std::cout`.

---

## Q301. [M18 | Writing a move constructor | Easy | MCQ | [C++] · C++17 | OA] · Type T5 · P76

```cpp
#include <cstddef>
#include <iostream>
#include <utility>
struct Buf {
    int *p; std::size_t n;
    Buf(std::size_t n) : p(new int[n]), n(n) {}
    Buf(Buf &&o) noexcept : p(o.p), n(o.n) { o.p = nullptr; o.n = 0; }
    ~Buf() { delete[] p; }
};
int main() {
    Buf a(5);
    Buf b = std::move(a);
    std::cout << b.n << ' ' << a.n << ' ' << (a.p == nullptr) << '\n';
}
```

What does this program print?

- A. `5 5 0`
- B. `5 0 1`
- C. `0 5 1`
- D. Undefined behaviour: the array is deleted twice

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** Moving is **handing over the keys** instead of building a copy of the house. The new object takes the old object's pointer, and the old object is left holding **nothing** (`nullptr`, size 0). When `a` is destroyed, `delete[] nullptr` does nothing, so the array is freed only once, by `b`.

**Why:**

| Step | a.p | a.n | b.p | b.n |
|---|---|---|---|---|
| `Buf a(5)` | → X | 5 | | |
| move constructor | **nullptr** | **0** | → X | 5 |
| destroy b | | | `delete[] X` | |
| destroy a | `delete[] nullptr` (a no-op) | | | |

Prints `5 0 1`.

**Why the other options fail:**
- A: Describes a copy.
- C: Swaps source and target.
- D: The source was nulled, so only one real `delete[]` happens.

**Concept recap:** A move constructor's three steps: (1) steal the resource, (2) copy the plain fields, (3) **leave the source empty but destructible**. Mark it `noexcept` (Q308).

**Trap alert:** Forgetting step 3 (`o.p = nullptr`) creates exactly the double delete from Q298.

**Interview angle:** *"Implement a move constructor for a buffer class."* Write this, and explain: "Steal, then null the source so its destructor is harmless. It's O(1) and `noexcept`, so containers will use it."

---

## Q302. [M18 | Move assignment | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T5 · P76

Add this move assignment operator to `Buf` from Q301:

```cpp
Buf &operator=(Buf &&o) noexcept {
    if (this != &o) { delete[] p; p = o.p; n = o.n; o.p = nullptr; o.n = 0; }
    return *this;
}
```

```cpp
int main() {
    Buf a(3), b(7);
    a = std::move(b);
    std::cout << a.n << b.n;
    b = Buf(2);
    std::cout << b.n << '\n';
}
```

What does this program print?

- A. `372`
- B. `770`
- C. `730`
- D. `702`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Move assignment is "**give up what I own, then take over yours**". `a` frees its own 3-element array, takes `b`'s 7-element array, and `b` is left empty. A moved-from object can still be **reused**: assigning a fresh value to `b` makes it whole again.

**Why:**

| Step | a.n | b.n | Note |
|---|---|---|---|
| start | 3 | 7 | |
| `a = std::move(b)` | 7 | 0 | a's old array freed; b emptied |
| print | `7` | `0` | |
| `b = Buf(2)` (the temporary is moved in) | | 2 | the temporary is destroyed afterwards, holding nullptr |
| print | | `2` | |

Output: `702`.

**Why the other options fail:**
- A: Treats the operation as a swap that keeps both values.
- B: Assumes `b` keeps its size after being moved from.
- C: Forgets the reassignment of `b`.

**Concept recap:** Moved-from objects must remain **valid**: they can be assigned to, destroyed, or (for standard types) used via operations without preconditions (Q310).

**Trap alert:** Without the `this != &o` check, a self-move (`a = std::move(a)`) would delete the array and then "steal" the dangling pointer. Standard library types tolerate self-move-assignment, leaving the object valid but unspecified; your own types should too.

**Interview angle:** *"Should move assignment check for self-assignment?"* "A self-move is rare but must not corrupt the object. A cheap `this != &o` check works, or implement move assignment as swap-based (`swap(*this, o)`), which is self-safe automatically."

---

## Q303. [M18 | What std::move really does | Medium | MSQ | [C++] · C++17 | Interview] · Type T3 · P77

Select **all** true statements.

- A. `std::move(x)` by itself performs no move. It's only a cast to an rvalue reference.
- B. After `std::string t = std::move(s);`, the standard guarantees that `s` is empty.
- C. Passing `std::move(c)` where `c` is a `const std::string` results in a **copy**.
- D. After moving from a `std::vector v`, calling `v.clear()` and then reusing `v` is fine.
- E. After `int a = 5; int b = std::move(a);`, `a` is 0.

**Answer:** A, C, D

**In simple words:** `std::move` is a **label**, not an action. It marks something as "okay to take from". The actual taking happens only if a **move constructor or move assignment** receives it. Built-in types like `int` have nothing to take, so "moving" them just copies. A moved-from standard object is in a "valid but unspecified" state: you may reset it and reuse it, but you shouldn't assume what's inside.

**Why each option:**
- **A ✓:** It's `static_cast<T&&>(x)`.
- **B ✗:** For standard types it's "valid but unspecified". It's *usually* empty in practice, but that isn't guaranteed (short-string optimisation implementations may leave contents behind).
- **C ✓:** `const std::string&&` can't bind to the move constructor's `std::string&&`, so it falls back to the copy constructor (Q295).
- **D ✓:** `clear()` has no preconditions. Afterwards the vector is empty and fully usable.
- **E ✗:** Verified: `a` stays 5. Moving a scalar is a copy.

**Concept recap:** On a moved-from object, safe operations are those with **no preconditions**: assignment, `clear()`, `size()`, destruction. Unsafe ones assume contents: `front()`, `back()`, `operator[]`.

**Trap alert:** A static analyser's "use after move" warning flags reading a moved-from object. Reassign it before reading.

**Interview angle:** *"What does std::move do?"* "It casts its argument to an rvalue reference so that move operations can be selected. It moves nothing itself, and the source is left valid but unspecified."

---

## Q304. [M18 | std::forward vs std::move vs no cast | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T6 · P77

```cpp
#include <iostream>
#include <string>
#include <utility>
using std::cout;
void take(const std::string &) { cout << "L "; }
void take(std::string &&)      { cout << "R "; }
template <class T> void fwd(T &&x)   { take(std::forward<T>(x)); }
template <class T> void nofwd(T &&x) { take(x); }
template <class T> void mv(T &&x)    { take(std::move(x)); }
int main() {
    std::string s = "a";
    fwd(s);
    fwd(std::string("b"));
    nofwd(std::string("c"));
    mv(s);
    cout << '\n';
}
```

What does this program print?

- A. `L R L R`
- B. `L R R R`
- C. `L L L R`
- D. `R R L L`

**Answer:** A

**Behaviour class:** Well-defined. Note that `mv(s)` treats the caller's `s` as movable, which is dangerous in real code.

**In simple words:**
- `T&&` in a template is a **forwarding reference**: it accepts anything, and `T` remembers whether the caller passed a named variable or a temporary.
- `std::forward<T>(x)` **passes it on exactly as received**: named stays named, temporary stays temporary.
- Plain `x` always looks like a named variable (Q296).
- `std::move(x)` always looks like a temporary, even when the caller's variable shouldn't be emptied.

**Why:**

| Call | T deduced as | Expression passed to `take` | Category | Prints |
|---|---|---|---|---|
| `fwd(s)` | `std::string&` | `forward<std::string&>(x)` | lvalue | `L` |
| `fwd(std::string("b"))` | `std::string` | `forward<std::string>(x)` | rvalue | `R` |
| `nofwd(std::string("c"))` | `std::string` | `x` (a name) | **lvalue** | `L` |
| `mv(s)` | `std::string&` | `std::move(x)` | rvalue | `R` (s may now be emptied!) |

**Why the other options fail:**
- B: Thinks a named forwarding reference forwards automatically.
- C: Thinks `std::forward` always produces an lvalue.
- D: Reverses the rules.

**Concept recap:** **Perfect forwarding** = `T&&` + `std::forward<T>`. The reference-collapsing rules: `T& &&` → `T&`; `T&& &&` → `T&&`. This is how `emplace_back` and `make_unique` pass constructor arguments through without extra copies.

**Trap alert:** Using `std::move` on a forwarding reference (as `mv` does) can **steal from the caller's lvalue**. Use `std::forward` in templates and `std::move` on rvalue references you own.

**Interview angle:** *"std::move vs std::forward?"* "`move` unconditionally casts to an rvalue. `forward<T>` casts to an rvalue only if the original argument was an rvalue, preserving its value category. Use `forward` with forwarding references (`T&&` with deduced `T`)."

---

## Q305. [M18 | Declaring a move operation deletes the copies | Hard | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P70

```cpp
#include <utility>
#include <vector>
struct M {
    M() = default;
    M(M &&) = default;
};
M a;
```

Select **all** statements that **compile**.

- A. `M b = std::move(a);`
- B. `M c = a;`
- C. `M d; d = std::move(a);`
- D. `std::vector<M> v; v.push_back(M());`
- E. `M e = M();`

**Answer:** A, D, E

**In simple words:** Writing your own move constructor tells C++ "this type is special about ownership". The compiler then **deletes** the automatic copy constructor and copy assignment, and doesn't create a move *assignment* either. So `M` can be **built** by moving, but not copied, and can't be assigned at all unless you declare the assignment operators yourself.

**Why each option:**

| Option | Needs | Available? | Result |
|---|---|---|---|
| A | move constructor | ✓ (defaulted) | compiles |
| B | copy constructor | **deleted** (because a move constructor is declared) | error |
| C | move assignment → not declared → falls back to copy assignment, which is **deleted** | | error |
| D | `push_back(M&&)` → move constructor (and moves on reallocation) | ✓ | compiles (a move-only element type) |
| E | nothing (C++17 guaranteed elision) | — | compiles |

(B and C verified: "use of deleted function".)

**Concept recap:** The **rule of five:** if you declare any of the destructor, copy constructor, copy assignment, move constructor or move assignment, think about all five. Declare the ones you want, and `= default` / `= delete` the rest explicitly.

**Trap alert:** Option C surprises people: a defaulted move **constructor** doesn't give you move **assignment**. Add `M &operator=(M &&) = default;`.

**Interview angle:** *"What happens to the copy operations if I declare a move constructor?"* "They're implicitly deleted, and the move assignment isn't declared at all. That's why you should declare all five intentionally, or none (the rule of zero)."

---

## Q306. [M18 | The rule of zero with a unique_ptr member | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T2 · P70

```cpp
#include <memory>
#include <string>
#include <utility>
#include <vector>
struct Widget {
    std::string name;
    std::vector<int> data;
    std::unique_ptr<int> extra;
};
int main() {
    Widget w;                         // L1
    Widget w2 = std::move(w);         // L2
    Widget w3 = w2;                   // L3
}
```

Which line fails to compile?

- A. L1, because `Widget` has no constructor
- B. L2, because `Widget` has no move constructor
- C. L3, because the implicit copy constructor is **deleted**: `std::unique_ptr` isn't copyable
- D. None of them

**Answer:** C

**Behaviour class:** Ill-formed (L3).

**In simple words:** With the **rule of zero**, you write **no** special member functions. Each member already knows how to copy, move and clean itself up, and the class simply **inherits the abilities of its members**. A `unique_ptr` can be moved but not copied, so `Widget` becomes move-only automatically. That's exactly right for a unique owner.

**Why:**

| Operation | `std::string` | `std::vector` | `std::unique_ptr` | Widget |
|---|---|---|---|---|
| default construct | ✓ | ✓ | ✓ | ✓ (L1) |
| move | ✓ | ✓ | ✓ | ✓ (L2) |
| copy | ✓ | ✓ | **✗ deleted** | **✗ deleted** (L3) |
| destroy | ✓ | ✓ | ✓ (frees) | ✓ automatic, no leaks |

**Why the other options fail:**
- A: The default constructor is implicitly generated.
- B: The implicit move constructor is generated (nothing suppresses it).
- D: Misses that the copy is deleted.

**Concept recap:** Rule of zero: classes should either manage **one** resource (and then follow the rule of five, as `unique_ptr` itself does) or manage **none** directly (and declare none of the five). Most classes should be the second kind.

**Trap alert:** If you need deep copies of `extra`, write a copy constructor **in a small wrapper type**, and keep `Widget` itself at the rule of zero.

**Interview angle:** *"What is the rule of zero?"* "Don't write destructors or copy/move operations. Compose classes from RAII members (`std::string`, `std::vector`, smart pointers) and let the compiler generate correct behaviour, including move-only behaviour when a member is move-only."

---

## Q307. [M18 | `return std::move(local)` disables NRVO | Hard | MCQ | [C++] · C++17 | Interview] · Type T7 · P71

```cpp
#include <iostream>
#include <utility>
using std::cout;
struct T {
    T() {}
    T(const T &)     { cout << "C"; }
    T(T &&) noexcept { cout << "M"; }
};
T good() { T t; return t; }
T bad()  { T t; return std::move(t); }
int main() {
    T a = good();
    cout << '|';
    T b = bad();
    cout << '\n';
}
```

Which statement is correct?

- A. `good()` prints nothing or `M` (NRVO may build `t` directly in `a`; otherwise `t` is moved). `bad()` always prints exactly `M`, because `std::move(t)` is no longer just the name of a local, so NRVO can't apply.
- B. Both always print `M`.
- C. `good()` prints `C`, and `bad()` prints `M`.
- D. Both always print nothing.

**Answer:** A

**Behaviour class:** Well-defined. `good()`'s output is unspecified between the two allowed forms (typical optimising compilers print nothing).

**In simple words:** For `return t;`, the compiler may build `t` **directly in the caller's variable**, so there's nothing to copy or move. If it doesn't, it moves `t` automatically. Writing `return std::move(t);` *looks* like a speed-up, but it **turns off** that first trick and forces a move every time. Just write `return t;`.

**Why:**

| Function | Returned expression | NRVO allowed? | Output |
|---|---|---|---|
| `good` | `t`, the name of a local | ✓ (optional) | `` or `M` |
| `bad` | `std::move(t)`, a function call result | ✗ | `M` |

(Typical result: `|M`.)

**Why the other options fail:**
- B: NRVO commonly removes the move in `good`.
- C: A local is never *copied* on return when it has a move constructor (an implicit move happens first).
- D: `bad` can't elide.

**Concept recap:** The only common reasons to write `std::move` in a return statement: returning a **member** or a **reference parameter** that you want to move out. Locals and by-value parameters are moved automatically (Q238).

**Trap alert:** Compilers often warn about this ("moving a local object in a return statement prevents copy elision"). Treat that warning as an error.

**Interview angle:** *"Should I write `return std::move(result);`?"* "No, for a local. It blocks NRVO and buys nothing, because the implicit move already happens."

---

## Q308. [M18 | noexcept move constructors and vector reallocation | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T8

```cpp
#include <iostream>
#include <vector>
using std::cout;
struct A { A() {} A(const A &) { cout << "C"; } A(A &&)          { cout << "M"; } };
struct B { B() {} B(const B &) { cout << "C"; } B(B &&) noexcept { cout << "M"; } };
int main() {
    std::vector<A> va; va.reserve(1); va.emplace_back(); va.emplace_back();
    cout << '|';
    std::vector<B> vb; vb.reserve(1); vb.emplace_back(); vb.emplace_back();
    cout << '\n';
}
```

What does this program print?

- A. `M|M`
- B. `C|M`
- C. `C|C`
- D. `M|C`

**Answer:** B

**Behaviour class:** Well-defined (this behaviour is required by the strong exception guarantee of `push_back`/`emplace_back`).

**In simple words:** When a vector outgrows its storage, it must transfer the old elements into new storage. If something goes wrong halfway, it wants to **undo cleanly**. Copying is safe to undo (the originals are untouched). Moving isn't (the originals have already been gutted). So the vector **moves only if the move constructor promises not to throw** (`noexcept`). Otherwise it plays safe and **copies**.

**Why:**

| Vector | Capacity 1 → second `emplace_back` → reallocation | Move constructor `noexcept`? | Existing element transferred by | Prints |
|---|---|---|---|---|
| `va` | yes | ✗ | **copy** (`std::move_if_noexcept`) | `C` |
| `vb` | yes | ✓ | **move** | `M` |

The new element is constructed in place by `emplace_back`, so it prints nothing.

**Why the other options fail:**
- A: Assumes vectors always move.
- C: Assumes vectors always copy.
- D: Reversed.

**Concept recap:** **Always mark move constructors and move assignments `noexcept`** when they can't throw (pointer swaps can't). Otherwise containers silently copy, which can be a huge performance cliff. A defaulted move constructor is `noexcept` automatically when all members' moves are.

**Trap alert:** If the type is **move-only** (no copy constructor), the vector moves even without `noexcept`. It has no choice, and the strong guarantee is lost.

**Interview angle:** *"Why should move constructors be noexcept?"* "`std::vector` (via `move_if_noexcept`) only moves elements during reallocation if moving can't throw. That keeps `push_back`'s strong exception guarantee. Without `noexcept`, you get copies."

---

## Q309. [M18 | RAII guards during an exception | Easy | MCQ | [C++] · C++17 | OA] · Type T7

```cpp
#include <iostream>
#include <stdexcept>
using std::cout;
struct Guard {
    const char *n;
    Guard(const char *n) : n(n) { cout << "+" << n; }
    ~Guard() { cout << "-" << n; }
};
void work(bool fail) {
    Guard g("A");
    if (fail) throw std::runtime_error("x");
    Guard h("B");
}
int main() {
    try { work(false); work(true); }
    catch (...) { cout << "!"; }
    cout << '\n';
}
```

What does this program print?

- A. `+A+A!`
- B. `+A+B-B-A+A!`
- C. `+A+B-B-A+A-A!`
- D. `+A+B-A-B+A-A!`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** **RAII** = "**R**esource **A**cquisition **I**s **I**nitialisation": tie a resource (a lock, a file, memory) to an object's lifetime. Its destructor releases the resource **however** the scope ends: a normal return *or* an exception. Here `-A` prints even when `work` is interrupted by an exception.

**Why:**

| Call | Events | Output |
|---|---|---|
| `work(false)` | A, B constructed; normal exit → destroy B, then A | `+A+B-B-A` |
| `work(true)` | A constructed; throw → **stack unwinding** destroys A; B was never constructed | `+A-A` |
| catch | | `!` |

**Why the other options fail:**
- A, B: Forget that destructors run during unwinding.
- D: Destroys in construction order instead of reverse order.

**Concept recap:** RAII is C++'s replacement for `finally`. Standard RAII types: `std::lock_guard`, `std::unique_ptr`, `std::fstream`, `std::vector`.

**Trap alert:** If an exception is **never caught**, whether the stack is unwound before `std::terminate` is **implementation-defined**, so destructors may not run. Catch exceptions somewhere (at least in `main`) if clean-up matters.

**Interview angle:** *"What is RAII and why does it matter for exceptions?"* "Each resource is owned by a stack object whose destructor releases it. Because destructors run during stack unwinding, resources are released on every exit path without try/finally, giving exception safety by construction."

---

## Q310. [M18 | Using a moved-from object | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T8 · P76

```cpp
#include <iostream>
#include <utility>
#include <vector>
int main() {
    std::vector<int> v{1, 2, 3};
    std::vector<int> w = std::move(v);
    v.push_back(9);
    std::cout << w.size() << ' ' << v.back() << '\n';
}
```

What does this program print?

- A. Undefined behaviour: `v` is used after being moved from
- B. `3 3`
- C. `3 9`
- D. `0 9`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** After a move, `v` is like an **emptied box that's still a perfectly good box**: you just aren't told exactly what's left in it. You may put new things in. `push_back` works on any valid vector, and afterwards the last element is certainly the 9 you just added.

**Why:**

| Step | w | v |
|---|---|---|
| construct v | | {1, 2, 3} |
| `w = std::move(v)` | {1, 2, 3} | valid but unspecified (in practice empty) |
| `v.push_back(9)` | | ends with 9 |
| print | `3` | `v.back()` = `9` |

**Why the other options fail:**
- A: Using a moved-from object is fine for operations without preconditions. `push_back` has none, and after it, `back()`'s precondition (non-empty) holds.
- B: Thinks `v` still holds its old elements, and confuses `back()`.
- D: Thinks `w` is the moved-from one.

**Concept recap:** Printing `v.size()` right after the move would be **unspecified**, not UB. Every implementation gives 0 for vectors, but the standard only guarantees "valid". This question was built to avoid relying on it.

**Trap alert:** `v.front()` or `v[0]` right after the move would be UB **if** `v` is empty (as it typically is). Those operations have a non-empty precondition.

**Interview angle:** *"Can you use an object after moving from it?"* "Yes, for operations that don't assume a particular state: assigning, clearing, destroying, `push_back`, and so on. Reading its value without resetting it first is a bug."

---

## Q311. [M19 | Catching by value vs by reference | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P61

```cpp
#include <exception>
#include <iostream>
using std::cout;
struct Base : std::exception { const char *what() const noexcept override { return "base"; } };
struct Derived : Base { const char *what() const noexcept override { return "derived"; } };
int main() {
    try { throw Derived(); } catch (Base b)        { cout << b.what() << ' '; }
    try { throw Derived(); } catch (const Base &b) { cout << b.what(); }
    cout << '\n';
}
```

What does this program print?

- A. `derived derived`
- B. `base base`
- C. `derived base`
- D. `base derived`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Catching **by value** makes a **copy** of the exception as the declared type, which **slices** a `Derived` down to a plain `Base` (Q277). Catching **by reference** refers to the original thrown object, so the virtual `what()` still reaches `Derived`'s version.

**Why:**

| Handler | What the handler holds | `what()` |
|---|---|---|
| `catch (Base b)` | a copy constructed as `Base` (sliced) | "base" |
| `catch (const Base &b)` | the original `Derived` object | "derived" |

**Why the other options fail:**
- A: Ignores slicing in the by-value handler.
- B: Ignores dispatch in the by-reference handler.
- C: Reversed.

**Concept recap:** The rule: **throw by value, catch by (const) reference.** It avoids slicing and an extra copy, and allows polymorphic handling.

**Trap alert:** Catching by **pointer** (`throw new E;` with `catch (E *e)`) requires someone to `delete` the exception, and is almost never what you want.

**Interview angle:** *"How should you catch exceptions?"* "By `const&` (or `&` if you modify and rethrow), ordered from most derived to most general. Catching by value slices and copies."

---

## Q312. [M19 | The order of catch clauses | Easy | MCQ | [C++] · C++17 | OA] · Type T7

```cpp
#include <iostream>
#include <stdexcept>
int main() {
    try { throw std::out_of_range("oops"); }
    catch (const std::exception &e)    { std::cout << "E"; }
    catch (const std::out_of_range &e) { std::cout << "O"; }
    std::cout << '\n';
}
```

What does this program print?

- A. `E`
- B. `O`
- C. `EO`
- D. Compile error: the second handler is unreachable

**Answer:** A

**Behaviour class:** Well-defined. Compilers usually warn that the second handler will never be used.

**In simple words:** Catch clauses are checked **top to bottom, and the first one that fits wins**, not the best one. `out_of_range` *is a* `std::exception`, so the first clause catches it, and the more specific clause below never gets a chance.

**Why:**

| Handler (in order) | Matches `std::out_of_range`? |
|---|---|
| `const std::exception&` | ✓ (it's a base class) → **chosen** |
| `const std::out_of_range&` | never examined |

**Why the other options fail:**
- B: Assumes best-match selection, the way overloading works. Handlers don't work that way.
- C: Only one handler runs per exception.
- D: An unreachable handler isn't ill-formed; it typically gets a warning.

**Concept recap:** Order handlers from **most derived to most general**, ending with `catch (...)` if needed.

**Trap alert:** The standard hierarchy: `std::exception` → `std::logic_error` (→ `out_of_range`, `invalid_argument`, `length_error`, `domain_error`), `std::runtime_error` (→ `overflow_error`, `range_error`, `system_error`), `std::bad_alloc`, `std::bad_cast`, and so on.

**Interview angle:** *"How is a catch handler selected?"* "The first handler, in textual order, whose type matches (exactly, as a base class, or via certain pointer conversions). No conversions like int → long apply."

---

## Q313. [M19 | Destructors during stack unwinding | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
using std::cout;
struct G {
    char c;
    G(char c) : c(c) { cout << c; }
    ~G() { cout << (char)(c - 32); }   // prints the uppercase letter
};
void f() { G x('x'); throw 1; G y('y'); }
void g() { G a('a'); f(); G b('b'); }
int main() {
    try { g(); } catch (int) { cout << "!"; }
    cout << '\n';
}
```

What does this program print?

- A. `axXA!`
- B. `ax!XA`
- C. `axyYXA!`
- D. `axAX!`

**Answer:** A

**Behaviour class:** Well-defined (ASCII assumed: lowercase − 32 = uppercase).

**In simple words:** When an exception flies out of a function, C++ **walks back up the call stack** to find a `catch`. On the way, it **destroys every fully built local object**, innermost first. Objects that were never reached (`y`, `b`) don't exist, so there's nothing to destroy. The catch block runs only **after** all that clean-up.

**Why:**

| Event | Output |
|---|---|
| `g`: construct a | `a` |
| `f`: construct x | `x` |
| `throw 1` → unwinding starts in `f`: destroy x | `X` |
| leave `g`: destroy a (b was never constructed) | `A` |
| reach the handler in `main` | `!` |

**Why the other options fail:**
- B: The handler runs after unwinding, not before.
- C: `y` is never constructed.
- D: Destroys in the wrong order (outer before inner).

**Concept recap:** Stack unwinding destroys automatic objects in reverse order of construction, frame by frame, until a matching handler is found ([except.ctor]).

**Trap alert:** A destructor that **throws** while another exception is already unwinding calls `std::terminate`. That's one reason destructors are implicitly `noexcept`.

**Interview angle:** *"What happens to local objects when an exception is thrown?"* "Every fully constructed automatic object between the throw and the handler is destroyed, in reverse order. That's what makes RAII clean-up automatic."

---

## Q314. [M19 | Rethrowing from a nested handler | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7

Using the `G` type from Q313:

```cpp
int main() {
    try {
        G a('a');
        try {
            G b('b');
            throw 1;
        } catch (int) {
            cout << "1";
            G c('c');
            throw;
        }
    } catch (int) { cout << "2"; }
    cout << '\n';
}
```

What does this program print?

- A. `ab1cB2CA`
- B. `abB1c2CA`
- C. `abB1cCA2`
- D. `abB1cAC2`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Follow the exception step by step:
- It leaves the **inner** `try`, destroying `b`.
- The inner `catch` runs and builds `c`, then `throw;` **passes the same exception on**.
- Leaving that catch block destroys `c`, and leaving the outer `try` destroys `a`.
- Finally, the outer `catch` runs.

**Why:**

| Event | Output |
|---|---|
| construct a, b | `ab` |
| `throw 1` → leave the inner try: destroy b | `B` |
| inner handler | `1` |
| construct c | `c` |
| `throw;` → leave the handler: destroy c | `C` |
| leave the outer try: destroy a | `A` |
| outer handler | `2` |

Output: `abB1cCA2`.

**Why the other options fail:**
- A: Delays destroying `b` until after the handler.
- B: Runs the outer handler before unwinding the inner scopes.
- D: Destroys `a` before `c`. `c` is in the inner (handler) scope, so it's destroyed first.

**Concept recap:** `throw;` rethrows the **current exception object** (the same object, not a copy). The handler's locals are destroyed as the exception leaves the handler.

**Trap alert:** `throw;` with no exception currently being handled calls `std::terminate` (Q317-C).

**Interview angle:** *"Trace the output of nested try/catch with rethrow."* Use the event table: construction, unwinding per scope, then the handler. Interviewers mainly check the ordering of destructors relative to handlers.

---

## Q315. [M19 | A constructor that throws | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
using std::cout;
struct M {
    const char *n;
    M(const char *n) : n(n) { cout << "+" << n; }
    ~M() { cout << "-" << n; }
};
struct W {
    M m1{"1"};
    M m2{"2"};
    W()  { cout << "W"; throw 0; }
    ~W() { cout << "~W"; }
};
int main() {
    try { W w; } catch (int) { cout << "!"; }
    cout << '\n';
}
```

What does this program print?

- A. `+1+2W~W-2-1!`
- B. `+1+2W-2-1!`
- C. `+1+2W!`
- D. `+1+2W-1-2!`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** An object counts as **"born" only when its constructor finishes**. `W`'s constructor threw, so `W` was never fully born, and its destructor **doesn't run**. But its members `m1` and `m2` *were* fully built, so they **are** destroyed (in reverse order). This is why resources should live in members (RAII), not be grabbed by hand in the constructor body.

**Why:**

| Event | Output |
|---|---|
| members built | `+1+2` |
| constructor body | `W` |
| `throw 0` → `W` incomplete → **no `~W`** | |
| fully built members destroyed in reverse | `-2-1` |
| handler | `!` |

**Why the other options fail:**
- A: Calls `~W` for an object that was never completely constructed.
- C: Forgets to destroy the members.
- D: Wrong order.

**Concept recap:** Leak-safe constructors: acquire every resource through RAII members (`unique_ptr`, `vector`), never `new` raw pointers in the body. Otherwise a later throw leaks the earlier allocations, because the destructor won't run.

**Trap alert:** With **delegating constructors**, the object counts as constructed once the *target* constructor finishes. If the delegating constructor's body then throws, the destructor **does** run (Q232).

**Interview angle:** *"If a constructor throws, is the destructor called?"* "Not for that object, because it was never fully constructed. But its fully constructed members and bases are destroyed. That's why RAII members prevent leaks."

---

## Q316. [M19 | Throwing out of a noexcept function | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T3

```cpp
#include <iostream>
void f() noexcept { throw 1; }
int main() {
    try { f(); }
    catch (...) { std::cout << "caught"; }
}
```

What happens?

- A. It prints `caught`.
- B. Compile error: `throw` inside a `noexcept` function.
- C. `std::terminate` is called: an exception can't escape a `noexcept` function, so the `catch` in `main` never runs (the program typically aborts).
- D. Undefined behaviour.

**Answer:** C

**Behaviour class:** Well-defined: `std::terminate` is required. (Whether the stack is unwound first is implementation-defined.)

**In simple words:** `noexcept` is a **promise**: "nothing escapes this function". If something tries to, C++ doesn't let it travel further. It **ends the program immediately** (`std::terminate`), so the `catch` in `main` never gets a chance. Breaking the promise isn't "undefined"; the punishment is clearly specified: termination.

**Why:**

| Step | Result |
|---|---|
| `throw 1` inside `f` | the exception reaches `f`'s boundary |
| `f` is `noexcept` | **`std::terminate()`** ([except.spec]p5) |
| `main`'s handler | never reached |

(Verified: the program aborts with a "terminate called after throwing an instance of 'int'" message.)

**Why the other options fail:**
- A: The handler is never reached.
- B: The code compiles. Compilers may warn that the `throw` will always call terminate.
- D: The outcome (terminate) is specified, so it isn't UB.

**Concept recap:** Destructors, `swap`, and move operations should be `noexcept`. Destructors are `noexcept` by default since C++11.

**Trap alert:** A `noexcept` function **may** throw internally if it also catches the exception inside itself. Only *escaping* is forbidden.

**Interview angle:** *"What does noexcept do?"* "It promises the function won't propagate exceptions. That enables optimisations (e.g. `vector` moving on reallocation, Q308). If it's violated, `std::terminate` is called."

---

## Q317. [M19 | `throw;` vs `throw e;` | Hard | MSQ | [C++] · C++17 | Interview] · Type T7

Using `Base` and `Derived` from Q311:

```cpp
try {
    try { throw Derived(); }
    catch (Base &e) { throw e; }     // version 1; version 2 uses:  throw;
}
catch (Derived &) { std::cout << "D"; }
catch (Base &)    { std::cout << "B"; }
```

Select **all** true statements.

- A. Version 1 (`throw e;`) throws a **new copy of type `Base`** (sliced), so it prints `B`.
- B. Version 2 (`throw;`) rethrows the **original** `Derived` object, so it prints `D`.
- C. `throw;` executed when no exception is being handled calls `std::terminate`.
- D. Catching `Base e` by value would avoid slicing in version 1.
- E. `throw e;` and `throw;` behave identically whenever `e` is caught by reference.

**Answer:** A, B, C

**In simple words:** `throw e;` means **"throw a new exception, copied from `e` as the type written in the code"**. `e` is declared `Base&`, so the copy is a plain `Base`: sliced. `throw;` means **"keep the original exception flying"**, so nothing is copied and the type is preserved.

**Why each option** (A and B verified: `B` and `D`):
- **A ✓:** The thrown object's type is the **static** type of the expression `e`, which is `Base`. A copy is made.
- **B ✓:** It rethrows the current exception object itself.
- **C ✓:** [except.throw]p9.
- **D ✗:** By-value catching slices *earlier*, as soon as the handler is entered (Q311).
- **E ✗:** As A and B show, they differ.

**Concept recap:** To rethrow, use `throw;`. To wrap or translate an exception, throw a new type (and optionally use `std::throw_with_nested`).

**Trap alert:** `std::exception_ptr p = std::current_exception();` with a later `std::rethrow_exception(p)` is the way to rethrow on a different thread, or after the handler has ended.

**Interview angle:** *"What's the difference between `throw;` and `throw e;`?"* "`throw;` rethrows the original object with its dynamic type preserved and no copy. `throw e;` copy-constructs a new exception of e's static type, which can slice."

---

## Q318. [M19 | Exception-safety guarantees | Easy | MSQ | [C++] · C++17 | Interview] · Type — concept

Select **all** true statements.

- A. **Basic guarantee:** if an exception is thrown, no resources leak, and all objects remain in a valid (but possibly changed) state.
- B. **Strong guarantee:** the operation either completes fully or has no effect at all (commit-or-rollback).
- C. **Nothrow guarantee:** the operation never throws. Destructors, `swap` and move operations should provide it.
- D. The strong guarantee means the operation never throws.
- E. The basic guarantee means the program state is unchanged after a failure.

**Answer:** A, B, C

**In simple words:** Three promises a function can make about what happens if something goes wrong inside it:
- **Basic:** "Nothing leaks or breaks, but things may have changed."
- **Strong:** "Either it all happened, or nothing happened", like a bank transfer.
- **Nothrow:** "I never fail."

**Why each option:**
- **A ✓, B ✓, C ✓:** The standard definitions (Dave Abrahams' guarantees).
- **D ✗:** That's the nothrow guarantee. A strong-guarantee operation may throw, it just rolls back.
- **E ✗:** That describes the strong guarantee.

**Concept recap:** Examples:
- `vector::push_back`: strong (when the element's move is `noexcept` or the type is copyable).
- Copy-and-swap assignment: strong (Q300).
- Most algorithms: basic.
- Destructors: nothrow.

**Trap alert:** The basic guarantee doesn't mean "nothing happened". A container can be left with fewer or different elements, as long as it's still usable.

**Interview angle:** *"What are the exception safety guarantees?"* "Nothrow > strong > basic (> none). Aim for the basic guarantee everywhere, the strong guarantee for key operations (often via copy-and-swap), and nothrow for destructors, swaps and moves."

---

## Q319. [M20 | Function template argument deduction | Easy | MCQ | [C++] · C++17 | OA] · Type T1 · P83

```cpp
#include <iostream>
template <class T> T mx(T a, T b) { return a > b ? a : b; }
int main() {
    std::cout << mx(3, 7) << ' ' << mx(2.5, 1.5) << ' ' << mx<double>(3, 2.5) << '\n';
    // mx(3, 2.5);   // the line in question
}
```

What does the program print, and what happens if the commented line is enabled?

- A. `7 2.5 3`; the commented line compiles and returns 3.
- B. `7 2.5 3`; the commented line is a **compile error**: `T` is deduced as `int` from 3 and as `double` from 2.5, which conflict.
- C. `7 2 3`; the commented line returns 2.5.
- D. `7 2.5 2.5`; the commented line is a compile error.

**Answer:** B

**Behaviour class:** Well-defined (as written). The commented line is ill-formed.

**In simple words:** A function template is a **recipe with a blank** (`T`), and the compiler fills in the blank by looking at the arguments. Each argument must give the **same answer** for `T`. `3` says "int" and `2.5` says "double", so the compiler can't decide. Unlike ordinary functions, deduction **doesn't convert** arguments to make them agree. Writing `mx<double>(3, 2.5)` fills in the blank yourself, and then normal conversion (3 → 3.0) is allowed.

**Why:**

| Call | Deduced T | Result |
|---|---|---|
| `mx(3, 7)` | int | 7 |
| `mx(2.5, 1.5)` | double | 2.5 |
| `mx<double>(3, 2.5)` | explicitly double (3 → 3.0) | 3 (prints as `3`) |
| `mx(3, 2.5)` | int vs double → **conflict** | "no matching function" |

**Why the other options fail:**
- A: Deduction never applies conversions to resolve a conflict.
- C: Treats 2.5 as an int.
- D: Wrong: max(3.0, 2.5) is 3.

**Concept recap:** Fixes: explicit template arguments, two type parameters with a common return type (`template<class A, class B> auto mx(A a, B b)`), or `std::max<double>(3, 2.5)`.

**Trap alert:** `std::max(3, 2.5)` fails for exactly this reason. So does `std::max(x, 0)` when `x` is a `size_t` (`size_t` vs `int`).

**Interview angle:** *"Why does `std::max(1, 2.0)` not compile?"* "Template argument deduction deduces `T` from each argument independently, and the results conflict. Implicit conversions aren't considered during deduction."

---

## Q320. [M20 | Deduction with arrays and const | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T6 · P83

```cpp
#include <iostream>
template <class T> void byVal(T x)  { std::cout << sizeof(x) << ' '; }
template <class T> void byRef(T &x) { std::cout << sizeof(x) << ' '; }
int main() {
    int arr[10];
    byVal(arr);
    byRef(arr);
    std::cout << '\n';
}
```

What does this program print, and what is `T` for `byVal(ci)` and `byRef(ci)`, where `const int ci = 1;`?

- A. `8 40`; for `ci`: `byVal` deduces `T = int`, and `byRef` deduces `T = const int`.
- B. `40 40`; for `ci`: both deduce `T = const int`.
- C. `8 8`; for `ci`: both deduce `T = int`.
- D. `40 8`; for `ci`: `byVal` deduces `T = const int`, and `byRef` deduces `T = int`.

**Answer:** A

**Behaviour class:** Well-defined. Only `sizeof` is applied to `arr`, so its elements are never read. The deduced types were checked with `std::is_same`.

**In simple words:**
- **By value** means "I'll get my own copy". Arrays can't be copied as a whole, so the array **decays to a pointer**, and a copy doesn't need to stay `const`.
- **By reference** means "I'll look at the original". So `T` keeps **everything**: the full array type `int[10]`, and the `const`.

**Why:**

| Call | Parameter | Deduced T | `sizeof(x)` |
|---|---|---|---|
| `byVal(arr)` | `T x` | `int*` (decay) | 8 |
| `byRef(arr)` | `T& x` | `int[10]` (no decay) | 40 |
| `byVal(ci)` | `T x` | `int` (top-level const dropped) | — |
| `byRef(ci)` | `T& x` | `const int` (const kept) | — |

**Why the other options fail:**
- B: Thinks arrays survive by-value deduction.
- C: Thinks references also decay.
- D: Reverses every rule.

**Concept recap:** These are exactly the `auto` rules from Q192, since `auto` uses template deduction. Pass by reference to keep an array's size: `template<class T, std::size_t N> std::size_t len(T (&)[N]) { return N; }`, which is how `std::size` works.

**Trap alert:** A forwarding reference `T&&` applied to an lvalue deduces `T` as a **reference** type (`int(&)[10]`, `const int&`), as in Q304.

**Interview angle:** *"How does template deduction treat arrays and const?"* "By value: arrays and functions decay to pointers, and top-level const/volatile are dropped. By reference: the exact type is preserved, including array bounds and const."

---

## Batch 16 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 301 | M18 | E | MCQ | B | Well-defined |
| 302 | M18 | M | MCQ | D | Well-defined |
| 303 | M18 | M | MSQ | A, C, D | — (rules) |
| 304 | M18 | H | MCQ | A | Well-defined |
| 305 | M18 | H | MSQ | A, D, E | Ill-formed / OK mix |
| 306 | M18 | M | MCQ | C | Ill-formed (L3) |
| 307 | M18 | H | MCQ | A | Well-defined (`good` unspecified between two forms) |
| 308 | M18 | H | MCQ | B | Well-defined |
| 309 | M18 | E | MCQ | C | Well-defined |
| 310 | M18 | M | MCQ | C | Well-defined |
| 311 | M19 | M | MCQ | D | Well-defined |
| 312 | M19 | E | MCQ | A | Well-defined |
| 313 | M19 | M | MCQ | A | Well-defined |
| 314 | M19 | H | MCQ | C | Well-defined |
| 315 | M19 | H | MCQ | B | Well-defined |
| 316 | M19 | M | MCQ | C | Well-defined (`std::terminate`) |
| 317 | M19 | H | MSQ | A, B, C | Well-defined |
| 318 | M19 | E | MSQ | A, B, C | — (concepts) |
| 319 | M20 | E | MCQ | B | Well-defined / ill-formed line |
| 320 | M20 | M | MCQ | A | Well-defined |

- **Difficulty:** 5 E / 8 M / 7 H.
  - M18 total: 3/7/8, as planned.
  - M19 total: 2/3/3, as planned.
- **Format:** 16 MCQ / 4 MSQ.
  - M18 total: 14/4, as planned.
  - M19 total: 6/2, as planned.
- **MCQ answer keys:** A×5, B×4, C×5, D×2. Running total Q1–Q320: A 65, B 64, C 65, D 62. Batch 17 will favour D.

---

**Completed: 320 / 400 Remaining: 80**

**Ledger, newly ✓ this batch:**
- P76 ×3 (Q301, Q302, Q310) ★ complete
- P77 ×2 (Q303, Q304)
- P70 → 4 (Q305, Q306)
- P71 → 4 (Q307)
- P83 ×2 (Q319, Q320) ★ complete (+1 planned)

**M18 complete** (Q293–Q310, 18 questions). **M19 complete** (Q311–Q318, 8 questions). **M20 started.**

*Say **"Continue."** for Batch 17 (Q321–Q340: M20 templates — more deduction, class templates with non-type parameters, full and partial specialisation, template vs non-template overload preference, `typename` in dependent names, variadic templates, SFINAE, C++20 concepts as a supplement, CRTP — then M21 STL begins with `push_back` vs `emplace_back`, iterator invalidation and container complexity).*
