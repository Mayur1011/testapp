# C/C++ MCQ/MSQ Bank — Batch 12 (Q221–Q240)

**What this batch covers:**
- **M14 Constructors, destructors & lifetime, finishing (Q221–Q238):**
  - members are initialised in **declaration order**, not initializer-list order, including the UB that causes
  - the full construction and destruction sequence: base → members → body, then the reverse
  - default member initialisers
  - destruction order for locals, function-local statics and globals; for arrays and members; and the destructor that never runs for leaked heap objects
  - temporaries and when their lifetime is extended
  - `= delete` for copying, and for blocking conversions
  - delegating constructors
  - value, zero and aggregate initialisation
  - narrowing in aggregate initialisation
  - guaranteed elision vs NRVO vs copying from `?:`
  - the static-initialisation-order fiasco
  - implicit move when returning a parameter
- **M15 Inheritance, starting (Q239–Q240):** access to base members, and public/protected/private inheritance seen from outside the class.

**Patterns:** P56, P57, P58, P71, P99

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - Object traces show every constructor and destructor call, and why it happens.

> **Conventions for every question in this batch:**
> - **C++17**, exceptions enabled, with guaranteed copy elision from prvalues.
> - Platform: a typical 64-bit system.
> - Needed headers are included, and `std::` is written explicitly. In multi-statement traces, `cout` means `std::cout`.

---

## Q221. [M14 | Members are initialised in declaration order | Easy | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P58

```cpp
#include <iostream>
struct Tr { Tr(const char *s) { std::cout << s; } };
struct M {
    Tr b{"b"};
    Tr a;
    Tr c;
    M() : c("c"), a("a") {}
};
int main() { M m; std::cout << '\n'; }
```

What does this program print?

- A. `cab`
- B. `bac`
- C. `cba`
- D. `abc`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** Members are built in the order they're **declared in the class**, like seats filled from the front row back. The order you *write* in the constructor's initializer list doesn't matter. It only says *how* each member is built, not *when*. `b` isn't in the list, so it uses its default initializer `{"b"}`.

**Why:**

| Declaration order | How it's initialised | Prints |
|---|---|---|
| `b` | default member initializer `{"b"}` | b |
| `a` | initializer list `a("a")` | a |
| `c` | initializer list `c("c")` | c |

Prints `bac`.

**Why the other options fail:**
- A: Follows the initializer-list order (`c`, then `a`), with `b` at the end.
- C, D: Other orderings, unrelated to the declaration order.

**Concept recap:** The rule ([class.base.init]p13): (1) virtual bases, (2) direct bases in declaration order, (3) non-static data members in **declaration order**, (4) the constructor body.

**Trap alert:** Compilers warn about this mismatch (a "member will be initialized after" warning), because it's a common source of bugs (Q222).

**Interview angle:** *"In what order are members initialised?"* "In declaration order in the class, always, regardless of the order in the initializer list. Write the initializer list in the same order as the declarations to avoid confusion."

---

## Q222. [M14 | Initialisation order causing UB | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T3 · P58

> **This code has undefined behaviour.**

```cpp
#include <iostream>
struct R {
    int area;
    int w, h;
    R(int w_, int h_) : w(w_), h(h_), area(w * h) {}
};
int main() {
    R r(3, 4);
    std::cout << r.area << '\n';
}
```

Which statement is correct?

- A. It prints `12`.
- B. It prints `0`.
- C. The behaviour is undefined: `area` is declared first, so it is initialised first, reading `w` and `h` before they have been given values.
- D. Compile error: `area` uses members that come later in the initializer list.

**Answer:** C

**Behaviour class:** **Undefined behaviour.**

**In simple words:** The initializer list *looks* as if `w` and `h` are set before `area`. But initialisation follows the **class's declaration order**, and `area` is declared first. So `area` is computed from `w` and `h` while they still contain garbage, like adding up a bill before anyone has written the prices.

> **⚠ UB box**
> - **Category:** reading uninitialised (indeterminate) member values.
> - **Rule:** [class.base.init]p13 (members are initialised in declaration order) together with [basic.indet]p2 (using an indeterminate value is UB).
> - **What you may see in practice:** often 0 or a random-looking number, and occasionally 12 by luck. Compilers usually warn about the order mismatch. **Not guaranteed.**

**Why:**

| Order | Member | Initialiser | State |
|---|---|---|---|
| 1st | `area` | `w * h` → reads **uninitialised** w and h | **UB** |
| 2nd | `w` | `w_` = 3 | |
| 3rd | `h` | `h_` = 4 | |

**Why the other options fail:**
- A, B: Each is one possible observation, presented as a guarantee.
- D: The code compiles. Only a warning is typical.

**Concept recap:** Fixes:
- reorder the declarations (`int w, h; int area;`);
- or compute from the **parameters**: `area(w_ * h_)`;
- or compute `area` in the constructor body.

**Trap alert:** The same bug happens with objects, such as a `std::string` member initialised from another `std::string` member that is declared *later*. That's UB too, and it often crashes.

**Interview angle:** A classic code-review question: *"Find the bug."* Name the rule (declaration order), and show the fix of initialising from constructor parameters rather than from other members.

---

## Q223. [M14 | Full construction and destruction sequence | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P58

```cpp
#include <iostream>
struct Tr2 {
    const char *s;
    Tr2(const char *s) : s(s) { std::cout << s; }
    ~Tr2() { std::cout << "~" << s; }
};
struct Base { Base() { std::cout << "B"; } ~Base() { std::cout << "~B"; } };
struct D : Base {
    Tr2 m1{"1"};
    Tr2 m2;
    D() : m2("2"), Base() { std::cout << "D"; }
    ~D() { std::cout << "~D"; }
};
int main() { { D d; std::cout << "|"; } std::cout << '\n'; }
```

What does this program print?

- A. `12BD|~D~B~2~1`
- B. `BD12|~D~2~1~B`
- C. `B21D|~D~1~2~B`
- D. `B12D|~D~2~1~B`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Building an object is like building a house: **foundation first (the base class), then the rooms (the members, in declaration order), then the finishing touches (the constructor body)**. Demolition is the exact reverse: finishing touches first (the destructor body), then the rooms in reverse, then the foundation last.

**Why:**

| Phase | Step | Prints |
|---|---|---|
| construction | 1. base `Base` | B |
| | 2. members in declaration order: m1, m2 | 1 2 |
| | 3. the body of `D()` | D |
| | | `|` |
| destruction | 1. the body of `~D()` | ~D |
| | 2. members in **reverse** order: m2, m1 | ~2 ~1 |
| | 3. the base | ~B |

Output: `B12D|~D~2~1~B`.

**Why the other options fail:**
- A: Builds members before the base.
- B: Runs the constructor body before the members.
- C: Uses the initializer-list order (`m2` first).

**Concept recap:** Destruction is always the **exact reverse** of construction. That guarantees each part is destroyed while everything it might depend on (earlier members, the base) still exists.

**Trap alert:** Writing `Base()` last in the initializer list doesn't delay it. Bases are always initialised before members.

**Interview angle:** *"Order of constructor and destructor calls with inheritance and members?"* State it in one breath: "Base → members (declaration order) → body; destruction is the reverse." Then mention virtual bases (constructed first, by the most-derived class) for a bonus (M15).

---

## Q224. [M14 | Default member initialisers vs the initializer list | Easy | MCQ | [C++] · C++17 | OA] · Type T1

```cpp
#include <iostream>
struct C {
    int a = 1;
    int b = 2;
    int c;
    C() : b(20) {}
    C(int x) : a(x), c(x) {}
};
int main() {
    C p(5);
    std::cout << p.a << ' ' << p.b << ' ' << p.c << '\n';
}
```

What does this program print?

- A. `5 2 5`
- B. `1 2 5`
- C. `5 20 5`
- D. `5 2 0`

**Answer:** A

**Behaviour class:** Well-defined. Every member that is read was initialised.

**In simple words:** A default member initializer (`int a = 1;`) is a **fallback value**. It's used only when the constructor that runs doesn't mention that member. If the constructor's initializer list sets the member, the list wins, and the fallback is ignored entirely (not "set, then overwritten").

**Why:** `C(5)` runs `C(int x)`.

| Member | In this constructor's list? | Value |
|---|---|---|
| a | ✓ `a(x)` | 5 |
| b | ✗ → default `= 2` | 2 |
| c | ✓ `c(x)` | 5 |

The other constructor's `b(20)` is irrelevant, because it isn't the one being called.

**Why the other options fail:**
- B: Uses the default for `a`, even though the list overrides it.
- C: Mixes in the other constructor's `b(20)`.
- D: Treats `c` as zero-initialised. It's set to 5 here. (With `C()`, `c` would be **uninitialised**, and reading it would be UB.)

**Concept recap:** Default member initialisers avoid repeating the same values in every constructor and prevent "forgot to initialise in one constructor" bugs. Give every scalar member one (`int c = 0;`).

**Trap alert:** A member with **no** default initializer that isn't in the constructor's list (like `c` in `C()`) is left **indeterminate** for built-in types. That's the source of Q222-style bugs.

**Interview angle:** *"Default member initialiser vs constructor initializer list?"* "Defaults say what a member usually starts as. The initializer list overrides them for a specific constructor. Both are better than assigning in the constructor body, which first default-constructs the member and then assigns it."

---

## Q225. [M14 | Destruction order: locals, function-local statics, globals | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
struct L {
    const char *n;
    L(const char *n) : n(n) { std::cout << "+" << n; }
    ~L() { std::cout << "-" << n; }
};
L g("G");
void f() { static L s("S"); L t("T"); }
int main() {
    std::cout << "[";
    L a("A");
    f(); f();
    { L b("B"); }
    std::cout << "]";
}
```

What is the **complete** output (including output after `main` returns)?

- A. `+G[+A+S+T-T+T-T+B-B]-A-S-G`
- B. `[+G+A+S+T-T+S+T-T+B-B]-A-S-G`
- C. `+G[+A+S+T-T+T-T+B-B]-A-G-S`
- D. `+G[+A+S+T-T+T-T+B-B-A]-S-G`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** Three kinds of lifetime:
- **Globals** are created before `main` starts.
- **Function-local `static`s** are created the first time their line runs, and never again.
- Ordinary **locals** are created and destroyed every time their block runs.

At the end, everything still alive is destroyed in **reverse order of creation**: `main`'s locals when `main` ends, then the statics and globals after `main` returns.

**Why:**

| Event | Output |
|---|---|
| before `main`: global g | `+G` |
| `main` starts | `[` |
| a | `+A` |
| 1st `f()`: s created (the first time only), then t, then t destroyed | `+S+T-T` |
| 2nd `f()`: s already exists; t created and destroyed | `+T-T` |
| inner block: b | `+B-B` |
| | `]` |
| `main` ends: a | `-A` |
| after `main`: statics and globals in reverse order of construction (S was built after G) | `-S-G` |

**Why the other options fail:**
- B: Constructs the global after `main` starts, and constructs `s` twice.
- C: Destroys G before S. Reverse order means the later-constructed S goes first.
- D: Destroys `a` before printing `]`. It lives until `main`'s closing brace.

**Concept recap:** Since C++11, a function-local static's initialisation is **thread-safe** (exactly once, even with concurrent first calls).

**Trap alert:** Objects that are still alive when `exit()` is called: statics *are* destroyed, but **locals on the stack are not**. With `abort()` or `quick_exit()`, nothing is destroyed.

**Interview angle:** *"When are static local objects constructed and destroyed?"* "Constructed on the first pass through their declaration (thread-safely since C++11), destroyed after `main` returns, in reverse order of construction, interleaved with globals."

---

## Q226. [M14 | Destruction order: arrays and members | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
struct E {
    int id;
    E(int i) : id(i) { std::cout << id; }
    ~E() { std::cout << -id; }
};
struct H { E x{1}; E y{2}; };
int main() {
    E arr[2] = {E(3), E(4)};
    H h;
    std::cout << "|";
}
```

What does this program print?

- A. `3412|-3-4-1-2`
- B. `3412|-1-2-3-4`
- C. `3344|-4-3-2-1`
- D. `3412|-2-1-4-3`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Think of a **stack of plates**: the last plate put down is the first one picked up. `arr` was created first (element 0, then element 1), then `h` (member x, then member y). So destruction starts with `h` (y, then x), then `arr` (element 1, then element 0).

**Why:**

| Construction | Prints |
|---|---|
| arr[0] = E(3): initialised directly (C++17, no temporary) | 3 |
| arr[1] = E(4) | 4 |
| h.x, then h.y | 1 2 |
| | `|` |

| Destruction (reverse) | Prints |
|---|---|
| h: members in reverse (y, then x) | -2 -1 |
| arr: elements in reverse (1, then 0) | -4 -3 |

Output: `3412|-2-1-4-3`.

**Why the other options fail:**
- A: Destroys in creation order.
- B: Destroys `h`'s members forwards and the array forwards.
- C: Assumes each `E(3)` creates a temporary that is copied. In C++17 the elements are initialised directly from the prvalues.

**Concept recap:** Array elements are destroyed from the **highest index to the lowest**. Members are destroyed in reverse declaration order. Whole objects are destroyed in reverse order of their creation.

**Trap alert:** If a constructor throws while building an array, the elements that were **already built** are destroyed in reverse order, and the rest never existed (M19).

**Interview angle:** *"Why is destruction reverse order?"* "Later objects may depend on earlier ones. Destroying in reverse means every object's dependencies are still alive while its destructor runs."

---

## Q227. [M14 | The destructor that never runs | Easy | MCQ | [C++] · C++17 | OA] · Type T7

```cpp
#include <iostream>
struct L {
    const char *n;
    L(const char *n) : n(n) { std::cout << "+" << n; }
    ~L() { std::cout << "-" << n; }
};
int main() {
    L *p = new L("H");
    L a("A");
}
```

What does this program print?

- A. `+H+A-A-H`
- B. `+H+A-A`
- C. `+H+A-H-A`
- D. `+A-A`

**Answer:** B

**Behaviour class:** Well-defined. The heap object is leaked, but that isn't UB.

**In simple words:** A local object is **cleaned up automatically** when its block ends. An object made with `new` is cleaned up **only when you write `delete`**. The program ends without a `delete`, so `H`'s destructor never runs. The memory is reclaimed by the OS, but any clean-up work in the destructor (flushing a file, logging) is simply skipped.

**Why:**

| Object | Storage | Destroyed when | Prints |
|---|---|---|---|
| H | heap (`new`) | only by `delete`, which never happens | `+H` |
| A | automatic | the end of `main` | `+A` … `-A` |

**Why the other options fail:**
- A, C: Assume heap objects are destroyed automatically at program exit.
- D: `new` does run the constructor.

**Concept recap:** That's the motivation for **RAII** and smart pointers: `auto p = std::make_unique<L>("H");` would print `-H` automatically when `p` goes out of scope (M18, M22).

**Trap alert:** Skipped destructors cause *real* bugs beyond memory leaks: unflushed buffers, unreleased locks, temporary files left behind, database transactions that are never committed.

**Interview angle:** *"Is a leak at program exit a problem if the OS reclaims the memory?"* "The memory isn't, but the skipped destructor might be: lost output, held locks, unreleased external resources. And in long-running programs, leaks accumulate. Use RAII so that clean-up always runs."

---

## Q228. [M14 | When temporaries die | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
struct T {
    T()  { std::cout << "c"; }
    ~T() { std::cout << "d"; }
    int v() const { return 1; }
};
int main() {
    int x = T().v();
    std::cout << "|";
    const T &r = T();
    std::cout << "|";
}
```

What does this program print?

- A. `cd|c|d`
- B. `c|c|dd`
- C. `cd|cd|`
- D. `c|cd|d`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** A temporary object (an unnamed `T()`) normally lives **until the end of the statement it appears in** (the `;`), and is then destroyed. There's one exception: when you **bind it directly to a `const` reference** (or an rvalue reference), it lives as long as that reference does.

**Why:**

| Line | Event | Output |
|---|---|---|
| `int x = T().v();` | temporary created, used, destroyed at the `;` | `cd` |
| | | `|` |
| `const T &r = T();` | temporary created, **lifetime extended** to r's scope | `c` |
| | | `|` |
| end of `main` | r's scope ends → the temporary is destroyed | `d` |

**Why the other options fail:**
- C: Ignores lifetime extension.
- B: Extends the first temporary too. It wasn't bound to a reference.
- D: Swaps the two cases.

**Concept recap:** The full-expression rule: temporaries are destroyed at the end of the complete expression that created them, in reverse order of creation ([class.temporary]p4). Lifetime extension only applies when binding **directly** (Q199, Q229).

**Trap alert:** `const char *s = std::string("hi").c_str();` → the string dies at the `;`, so `s` dangles. This is a very common real-world bug.

**Interview angle:** *"How long does a temporary live?"* "Until the end of the full-expression, unless it's bound directly to a `const T&` or `T&&` local reference, which extends it to the reference's lifetime. Extension doesn't propagate through function returns or member access."

---

## Q229. [M14 | Which bindings extend a temporary's lifetime? | Hard | MSQ | [C++] · C++17 | OA+Interview] · Type T5

Select **all** declarations after which the variable can be **safely used on later lines**.

- A. `const std::string &r = std::string("a") + "b";`
- B. `std::string &&r = std::string("x");`
- C. `const int &r = std::max(1, 2);`
- D. `const char *p = std::string("hi").c_str();`
- E. `auto &&r = std::vector<int>{1, 2};`

**Answer:** A, B, E

**In simple words:** A temporary is kept alive only when you bind a reference **directly to the temporary itself**. If you bind to something **returned by a function** that refers to the temporary (`std::max`'s result, or the `char*` from `c_str()`), the temporary still dies at the end of the line, and your reference or pointer is left dangling.

**Why each option:**
- **A ✓:** `operator+` returns a **new `std::string` prvalue**, and `r` binds directly to it → extended.
- **B ✓:** An rvalue reference binding directly to a prvalue → extended.
- **C ✗:** `std::max` returns a `const int&` to one of its temporary arguments. Not a direct binding, so it dangles (Q199).
- **D ✗:** `p` is a plain pointer into the string's buffer. Pointers never extend lifetimes, so the string dies at the `;`.
- **E ✓:** `auto&&` (a forwarding reference) binding to a prvalue deduces `std::vector<int>&&` → extended.

**Concept recap:** "Directly" means the reference is initialised by the temporary object itself (possibly through a cast or a member subobject), not by the result of a function call that returns a reference into it.

**Trap alert:** Range-for over a temporary's member function result, like `for (auto x : make_obj().items())`, dangles before C++23 if `items()` returns a reference. The temporary from `make_obj()` dies before the loop body runs. (C++23 extends temporaries in a range-for initializer.)

**Interview angle:** *"Explain this crash: `auto &v = getConfig().values(); use(v);`."* "`getConfig()` returns a temporary, and `values()` returns a reference into it. The temporary dies at the semicolon, so `v` dangles. Store `getConfig()` in a named variable first."

---

## Q230. [M14 | `= delete` for copy operations | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2

```cpp
#include <utility>
struct NoCopy {
    NoCopy() = default;
    NoCopy(const NoCopy &) = delete;
    NoCopy &operator=(const NoCopy &) = delete;
};
```

Select **all** lines that **compile** (C++17). Assume `NoCopy a;` already exists.

- A. `NoCopy x;`
- B. `NoCopy b = a;`
- C. `NoCopy c; c = a;`
- D. `NoCopy d = NoCopy();`
- E. `NoCopy e = std::move(a);`

**Answer:** A, D

**In simple words:** `= delete` means **"this operation is forbidden"**: any attempt to use it is a compile error. Deleting the copy operations makes the type non-copyable. Option D still works because, in C++17, `NoCopy d = NoCopy();` doesn't copy or move anything: `d` is built directly in place. Option E fails because declaring a copy constructor (even a deleted one) means **no move constructor is generated**, so the "move" falls back to the deleted copy.

**Why each option:**

| Option | Operation needed | Result |
|---|---|---|
| A | default constructor (`= default`) | ✓ |
| B | copy constructor | ✗ "use of deleted function" |
| C | copy assignment | ✗ |
| D | **none**: guaranteed elision from a prvalue | ✓ |
| E | move constructor → not declared → overload resolution picks the deleted copy constructor | ✗ |

**Concept recap:** Examples of non-copyable types: `std::unique_ptr`, `std::mutex`, `std::thread`, file handles. Add `NoCopy(NoCopy&&) = default;` to make the type movable but not copyable (the `unique_ptr` style).

**Trap alert:** Before C++17, option D also needed an accessible copy or move constructor (even though the copy was usually elided), so it failed to compile in C++14.

**Interview angle:** *"How do you make a class non-copyable?"* "Delete the copy constructor and copy assignment (C++11), and optionally default the move operations to keep it movable. The old C++03 trick was to declare them private and never define them."

---

## Q231. [M14 | Deleting an overload to block conversions | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T2

```cpp
#include <iostream>
void g(int x) { std::cout << x << ' '; }
void g(double) = delete;
```

Which calls **compile**?

1. `g(3)`
2. `g('a')`
3. `g(3.5)`
4. `g(3.5f)`

- A. Only 1
- B. 1 and 2
- C. 1, 2 and 4
- D. All four

**Answer:** B

**Behaviour class:** Calls 3 and 4 are ill-formed. Calls 1 and 2 are well-defined (they print `3 97`).

**In simple words:** A **deleted function still takes part in choosing the best overload**. If it turns out to be the best match, the call is rejected. So `g(double) = delete` means "if you pass something that's closest to a double, refuse it", instead of quietly converting it to `int`.

**Why:**

| Call | Argument | `g(int)` | `g(double)` | Best | Result |
|---|---|---|---|---|---|
| 1 | int | exact | conversion | g(int) | ✓ prints 3 |
| 2 | char | **promotion** (char → int) | conversion | g(int) | ✓ prints 97 |
| 3 | double | conversion | exact | g(double), **deleted** | ✗ |
| 4 | float | conversion | **promotion** (float → double) | g(double), **deleted** | ✗ |

**Why the other options fail:**
- A: Misses that char → int is a promotion, which beats a conversion.
- C: Misses that float → double is a promotion, which beats float → int.
- D: Deleted functions can't be called.

**Concept recap:** Overload ranking: exact match > **promotion** (char/short → int, float → double) > standard conversion > user-defined conversion (M17).

**Trap alert:** To accept *only* `int`, delete a template catch-all: `template<class T> void g(T) = delete;`. Then *every* non-`int` argument hits the deleted template (`char` too).

**Interview angle:** *"Besides copy operations, what else is `= delete` good for?"* "Blocking unwanted implicit conversions in overload sets, deleting `operator new` to forbid heap allocation, and preventing particular template instantiations."

---

## Q232. [M14 | Delegating constructors | Easy | MCQ | [C++] · C++17 | OA] · Type T7

```cpp
#include <iostream>
struct Pt {
    int x, y;
    Pt(int a, int b) : x(a), y(b) { std::cout << "2"; }
    Pt(int a) : Pt(a, a * 10)     { std::cout << "1"; }
    Pt() : Pt(7)                  { std::cout << "0"; }
};
int main() {
    Pt p;
    std::cout << " " << p.x << " " << p.y << '\n';
}
```

What does this program print?

- A. `012 7 70`
- B. `210 7 70`
- C. `0 7 70`
- D. `210 0 0`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** One constructor can **hand the work to another** constructor first (C++11 "delegation"). It's like a manager passing a job down the chain. The innermost constructor finishes first; then each delegating constructor runs its own body, on the way back up.

**Why:**

| Step | Constructor | Action | Prints |
|---|---|---|---|
| 1 | `Pt()` | delegates to `Pt(7)` | |
| 2 | `Pt(7)` | delegates to `Pt(7, 70)` | |
| 3 | `Pt(7, 70)` | x = 7, y = 70, body | 2 |
| 4 | back in `Pt(7)` | body | 1 |
| 5 | back in `Pt()` | body | 0 |

Output: `210 7 70`.

**Why the other options fail:**
- A: Runs the bodies outer-first. They run after the delegated-to constructor completes.
- C: Ignores the inner bodies.
- D: Thinks delegation doesn't set the members.

**Concept recap:** A delegating constructor's initializer list may contain **only** the delegation, with no other member initialisers. The object counts as fully constructed once the **target** constructor finishes, so if a delegating body later throws, the destructor **does** run.

**Trap alert:** Constructors that delegate to each other in a cycle (`A() : A(1) {}` with `A(int) : A() {}`) make the program ill-formed, no diagnostic required. In practice it means infinite recursion and a stack overflow.

**Interview angle:** *"How do you avoid duplicated code across constructors?"* "Delegating constructors (C++11), default member initialisers, or default arguments. Delegation keeps a single 'real' constructor."

---

## Q233. [M14 | Which values are guaranteed zero? | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T5

```cpp
struct Agg      { int a; double b; };
struct UserCtor { int a; UserCtor() {} };
Agg x{};
Agg y{5};
int *p = new int();
UserCtor u{};
int z;                 // a local variable inside a function
```

Select **all** values that are **guaranteed** to be 0.

- A. `x.a`
- B. `y.b`
- C. `*p`
- D. `u.a`
- E. `z`

**Answer:** A, B, C

**In simple words:** Empty braces or parentheses usually mean **"start from zero"**. But if a class has its **own** constructor that forgets to set a member, that member is left as garbage, because the class's constructor is in charge. A plain local `int z;` with no initializer also starts as garbage.

**Why each option:**

| Option | Kind of initialisation | Result |
|---|---|---|
| A `x{}` (aggregate, empty braces) | value-initialisation → zero | 0 ✓ |
| B `y{5}` | aggregate: listed members get values; the rest are value-initialised | b = 0 ✓ |
| C `new int()` | the `()` means value-initialise → 0 | 0 ✓ |
| D `u{}` | `UserCtor` has a **user-provided** constructor, so it's called, and it doesn't set `a` | indeterminate |
| E `int z;` | default-initialisation of a local `int` | indeterminate |

(Verified for A–C: prints `0 0 0`.)

**Concept recap:**
- **default-init** (`T x;`, `new T`) → built-in types are left indeterminate; classes call their default constructor.
- **value-init** (`T x{};`, `T()`, `new T()`) → zero first, unless there's a user-provided default constructor.
- **zero-init** → also applied to statics and globals before anything else.

**Trap alert:** `new int` (no parentheses) vs `new int()`: the first gives garbage, the second gives 0. One pair of parentheses changes the guarantee.

**Interview angle:** *"Difference between `int *p = new int;` and `new int();`?"* "Default-initialisation (an indeterminate value) vs value-initialisation (zero). In general, prefer `{}`: `int x{};` is always 0."

---

## Q234. [M14 | Aggregate initialisation with nesting and default initialisers | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T1

```cpp
#include <iostream>
struct In  { int a, b; };
struct Out { In i; int arr[3]; int c = 9; };
int main() {
    Out o{{1}, {4, 5}};
    Out p{};
    std::cout << o.i.b << ' ' << o.arr[2] << ' ' << o.c << ' '
              << p.c << ' ' << p.arr[0] << '\n';
}
```

What does this program print?

- A. `1 5 9 0 0`
- B. `0 5 9 9 0`
- C. `0 0 9 9 0`
- D. Compile error: a struct with a default member initializer isn't an aggregate

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Brace initialisation **fills members in order, and fills nested parts with nested braces**. Anything you leave out is set to its default initializer if it has one (like `c = 9`), or to zero otherwise. So a short initializer list never leaves garbage in an aggregate.

**Why:**

| Member | `o{{1}, {4, 5}}` | `p{}` |
|---|---|---|
| `i.a` | 1 | 0 |
| `i.b` | **0** (not given) | 0 |
| `arr` | {4, 5, **0**} | {**0**, 0, 0} |
| `c` | **9** (default member initializer) | **9** |

Prints `0 0 9 9 0`.

**Why the other options fail:**
- A: Reads the wrong members, and ignores `c = 9` in `p`.
- B: Treats `arr[2]` as 5, as if the list `{4, 5}` started at index 1.
- D: Since **C++14**, classes with default member initialisers are still aggregates.

**Concept recap:** An **aggregate** (C++17) is an array, or a class with no user-provided/explicit/inherited constructors, no private or protected non-static data members, no virtual functions, and no virtual/private/protected bases. Aggregates support brace initialisation member by member (and designated initialisers in C++20).

**Trap alert:** Adding **any** user-provided constructor makes the struct a non-aggregate, and `Out o{{1}, {4, 5}};` would then need a matching constructor.

**Interview angle:** *"What is an aggregate, and why does it matter?"* "A simple class with public data and no user-written constructors. It can be brace-initialised member by member, is often trivially copyable, and suits configuration structs and plain data."

---

## Q235. [M14 | Narrowing in aggregate initialisation | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P56

```cpp
struct S { int i; char c; };
```

Select **all** declarations that are valid **C++17**. (`char` is 8 bits and signed.)

- A. `S a{1, 'x'};`
- B. `S b{2.5, 'y'};`
- C. `S c(1, 'x');`
- D. `S d = {3, 65};`
- E. `S e{3, 300};`

**Answer:** A, D

**In simple words:** Brace initialisation of an aggregate applies the same **no-narrowing** rule as Q191: a double can't be squeezed into an `int`, and 300 can't be squeezed into a `char`. A constant that **fits** (65 fits in a `char`) is fine. And in C++17, aggregates can only be initialised with **braces**, not parentheses.

**Why each option:**

| Option | Issue | Valid? |
|---|---|---|
| A | exact types | ✓ |
| B | 2.5 → int: **narrowing** | ✗ |
| C | parenthesised aggregate initialisation is **C++20** only (P0960) | ✗ in C++17 |
| D | 65 → char: a constant that fits → not narrowing | ✓ (c = 'A') |
| E | 300 → char: a constant that **doesn't fit** → narrowing | ✗ |

**Concept recap:** C++20 allows `S c(1, 'x');`, but parentheses **don't** check for narrowing, so `S c(2.5, 'y');` would compile there too, silently truncating.

**Trap alert:** `std::vector<S> v; v.emplace_back(1, 'x');` fails in C++17 for an aggregate (it uses parentheses internally) and works in C++20. Use `v.push_back({1, 'x'})` in C++17.

**Interview angle:** *"Why doesn't `emplace_back` work with my simple struct?"* "Before C++20, `emplace_back` constructs with parentheses, and aggregates had no constructor to call. Use `push_back({...})`, add a constructor, or compile as C++20."

---

## Q236. [M14 | Guaranteed elision vs NRVO vs copying from `?:` | Hard | MCQ | [C++] · C++17 | Interview] · Type T7 · P71

```cpp
#include <iostream>
struct T {
    T()                    { std::cout << "D"; }
    T(const T &)           { std::cout << "C"; }
    T(T &&) noexcept       { std::cout << "M"; }
};
T make1()       { return T(); }
T make2()       { T t; return t; }
T make3(bool b) { T x, y; return b ? x : y; }
```

Which statement correctly describes what is printed by `T a = make1();`, then `T b = make2();`, then `T c = make3(true);`?

- A. `make1` prints exactly `D`. `make2` prints either `D` (if the compiler applies NRVO) or `DM`. `make3` prints `DDC`.
- B. All three print exactly `D`.
- C. `make1` prints `DM`, `make2` prints `DM`, and `make3` prints `DDM`.
- D. `make1` prints `D`, `make2` prints exactly `D`, and `make3` prints `DDM`.

**Answer:** A

**Behaviour class:** Well-defined. `make2`'s exact output is **unspecified** between the two allowed forms.

**In simple words:**
1. Returning a **brand-new temporary** (`return T();`): C++17 *guarantees* the object is built straight into the caller's variable, so no copy or move ever happens.
2. Returning a **named local** (`return t;`): the compiler is *allowed* to build `t` directly in the caller's variable (NRVO), but not *required* to. If it doesn't, it **moves** `t`, because it knows `t` is about to die.
3. Returning `b ? x : y` isn't simply returning a local's name. It's an expression that refers to an existing object, so the automatic move doesn't apply, and it **copies**.

**Why:**

| Function | Rule | Output |
|---|---|---|
| `make1` | guaranteed elision (C++17, prvalue) | `D` |
| `make2` | NRVO **optional**; otherwise implicit move ([class.copy.elision]p3) | `D` or `DM` (typical optimising compilers give `D`) |
| `make3` | x and y default-constructed. `b ? x : y` is an lvalue expression, not a plain name, so there's no NRVO and no implicit move → **copy** | `DDC` |

**Why the other options fail:**
- B: `make3` must build two locals and copy one of them.
- C: Ignores guaranteed elision, and says `make3` moves.
- D: Treats NRVO as guaranteed, and says `make3` moves.

**Concept recap:** Implicit move on return applies when the returned expression is the **name** of a local variable or a by-value parameter (Q238). Writing `return std::move(t);` *disables* NRVO, which is usually a pessimisation.

**Trap alert:** To move in the `?:` case, write `return b ? std::move(x) : std::move(y);`, or use `if` / `else` with separate returns, which allows implicit move or NRVO.

**Interview angle:** *"RVO vs NRVO?"* "RVO from a prvalue (`return T(...)`) is **guaranteed** in C++17. NRVO for a named local is an **optional** optimisation; if it isn't applied, the local is moved. Don't write `return std::move(local)`."

---

## Q237. [M14 | The static-initialisation-order fiasco | Hard | MCQ | [C++] · C++17 | Interview] · Type T3 · P99

```cpp
// b.cpp
#include <string>
std::string greeting = "hi";
```
```cpp
// a.cpp
#include <iostream>
#include <string>
extern std::string greeting;
std::string message = greeting + "!";
int main() { std::cout << message << '\n'; }
```

Which statement is correct?

- A. It always prints `hi!`, because `b.cpp`'s globals are initialised first.
- B. It always prints `!`, because `greeting` is empty when `message` is built.
- C. The relative order of dynamic initialisation of globals in **different** files is unspecified. If `a.cpp`'s initializer runs first, it uses a `std::string` object that hasn't been constructed yet, which is UB. The usual fix is a function-local static: `std::string &greeting() { static std::string g = "hi"; return g; }`.
- D. Compile error: a global can't be initialised from another file's global.

**Answer:** C

**Behaviour class:** The order is **unspecified**, and in the bad order the behaviour is **undefined** (using an object before its lifetime begins).

**In simple words:** Within one file, globals are built top to bottom. **Between files, C++ doesn't say which file's globals are built first.** If `message` happens to be built before `greeting`, it reads a string that doesn't exist yet. The fix is to make the dependency **build itself on first use**: a function with a `static` local is constructed the first time it's called, whichever file calls it.

> **⚠ UB box (the bad order)**
> - **Category:** using an object before its constructor has run (the static-initialisation-order fiasco).
> - **Rule:** [basic.start.dynamic]: dynamic initialisation of non-local variables in different translation units is indeterminately ordered. Using `greeting` before its initialisation accesses an object outside its lifetime ([basic.life]).
> - **What you may see in practice:** it often depends on the **order the object files are linked**: one order prints `hi!`, the other prints `!` or crashes. **Not guaranteed.**

**Why the other options fail:**
- A, B: Each picks one order as if it were guaranteed.
- D: `extern` declarations across files are ordinary and valid.

**Concept recap:** Fixes:
1. a function-local static (the "construct on first use" / Meyers idiom), which is thread-safe since C++11;
2. `constexpr` / `constinit` (C++20) data, which is initialised at compile time and so has no order problem;
3. avoiding globals that depend on other files' globals.

**Trap alert:** Constant-initialised globals like `const int N = 5;` or `constexpr` values are set up before **any** dynamic initialisation, so they're immune. The fiasco affects objects with **run-time** constructors (`std::string`, `std::map`, …).

**Interview angle:** *"What is the static initialisation order fiasco, and how do you fix it?"* Give the two-file example, then the function-local static fix, plus `constinit` in C++20. Mention the mirror-image problem at program end (destruction order, the "static deinitialisation fiasco").

---

## Q238. [M14 | Implicit move when returning a parameter | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P71

```cpp
#include <iostream>
struct T {
    T()              { std::cout << "D"; }
    T(const T &)     { std::cout << "C"; }
    T(T &&) noexcept { std::cout << "M"; }
};
T pass(T t) { return t; }
int main() {
    T a;
    T b = pass(a);
    std::cout << '\n';
}
```

What does this program print?

- A. `DCC`
- B. `DC`
- C. `DMM`
- D. `DCM`

**Answer:** D

**Behaviour class:** Well-defined. NRVO isn't permitted for function **parameters**, so the output is fully determined.

**In simple words:** Passing `a` by value makes a **copy** for the parameter `t`, because `a` is still needed afterwards. Returning `t` can't be optimised away (a parameter lives in a slot the caller set up, not in the return slot). But C++ knows `t` is about to die, so it **moves** it into `b` instead of copying.

**Why:**

| Step | Operation | Prints |
|---|---|---|
| `T a;` | default construction | D |
| call `pass(a)` | `a` is an lvalue → the parameter `t` is **copy**-constructed | C |
| `return t;` | `t` is a parameter → **implicit move** (no elision allowed) | M |
| into `b` | `pass(a)` is a prvalue → `b` is initialised directly (C++17 elision) | — |

Output: `DCM`.

**Why the other options fail:**
- A: Ignores the implicit move on return.
- B: Assumes the return is elided. That's not permitted for parameters.
- C: Assumes `a` is moved into the parameter. It's an lvalue, so it's copied (unless you write `pass(std::move(a))`).

**Concept recap:** The "sink" pattern: take a parameter by value, then move from it. `void set(std::string s) { name = std::move(s); }` copies for lvalue arguments and moves for rvalue arguments, all with one overload.

**Trap alert:** If `T`'s move constructor were deleted, `return t;` would fail to compile in C++17: overload resolution picks the deleted move first and doesn't fall back to the copy. That was relaxed in later standards.

**Interview angle:** *"Should I write `return std::move(param);`?"* "No. Returning a by-value parameter or a local by name already moves it automatically. An explicit `std::move` adds nothing, and for locals it even prevents NRVO."

---

## Q239. [M15 | Accessing base members from a derived class | Easy | MCQ | [C++] · C++17 | OA] · Type T2 · P57

```cpp
class B {
public:    int pub = 1;
protected: int pro = 2;
private:   int pri = 3;
};
class D : public B {
public:
    int f() { return pub + pro; }   // L1
    int g() { return pri; }         // L2
};
```

Which line fails to compile?

- A. L1
- B. L2
- C. Both L1 and L2
- D. Neither

**Answer:** B

**Behaviour class:** Ill-formed (L2: "'int B::pri' is private within this context").

**In simple words:** A derived class is like a **child in the family business**: it can use what's public (open to everyone) and what's protected (family only). But the parent's **private** things stay private, **even from the child**. The child still *has* them (they're part of every D object), but can't touch them by name.

**Why:**

| Member of B | Accessible inside D's member functions? |
|---|---|
| `pub` (public) | ✓ |
| `pro` (protected) | ✓ |
| `pri` (private) | ✗, but it still exists inside every D object (`sizeof(D)` includes it) |

**Why the other options fail:**
- A: Protected members exist precisely so derived classes can use them.
- C, D: Misjudge the private rule.

**Concept recap:** Derived classes reach private base data through the base's public or protected member functions, which keeps the base in control of its own invariants.

**Trap alert:** Protected access has a subtle limit: inside `D`, you may access `pro` through a `D` object (or `this`), but **not** through some other `B` object: `void h(B &other) { other.pro; }` is an error. Protected means "for my own inherited part", not "for every B".

**Interview angle:** *"private vs protected?"* "Private: only the class itself and its friends. Protected: the class plus derived classes, accessed through the derived object. Prefer private data with protected or public functions; protected data couples derived classes to the base's representation."

---

## Q240. [M15 | public / protected / private inheritance, seen from outside | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T2 · P57

```cpp
class B { public: void hello() {} };
class Pu : public B {};
class Pr : private B {};
class Pt : protected B {};
int main() {
    Pu x; x.hello();      // L1
    Pr y; y.hello();      // L2
    Pt z; z.hello();      // L3
    B *p = &y;            // L4
}
```

Which lines **compile**?

- A. L1, L2 and L3
- B. L1 and L4
- C. Only L1
- D. L1 and L3

**Answer:** C

**Behaviour class:** L2, L3 and L4 are ill-formed.

**In simple words:** The keyword in front of the base class sets the **"visibility ceiling"** for what's inherited:
- **public** inheritance keeps public things public: "a `Pu` is-a `B`, and everyone knows it".
- **protected** inheritance turns public into protected.
- **private** inheritance turns everything into private: "`Pr` uses a `B` internally, but that's nobody else's business".

So from outside, only `Pu` exposes `hello()`, and only `Pu` can be treated as a `B*`.

**Why:**

| Inheritance | Inherited `public` members become | Outside code can call `hello()`? | Outside code can convert to `B*`? |
|---|---|---|---|
| `public` | public | ✓ (L1) | ✓ |
| `protected` | protected | ✗ (L3) | ✗ |
| `private` | private | ✗ (L2) | ✗ (L4: "'B' is an inaccessible base of 'Pr'") |

**Why the other options fail:**
- A: Ignores how the inheritance mode restricts access.
- B: Treats the conversion to a private base as allowed.
- D: Treats protected inheritance as exposing public members to outsiders.

**Concept recap:** **Public** inheritance models "is-a" (substitutability). **Private** inheritance models "is-implemented-in-terms-of" (prefer composition, i.e. a `B` member, unless you need to override virtual functions or use the empty-base optimisation). Protected inheritance is rare.

**Trap alert:** `class D : B {}` with the `class` keyword means **private** inheritance by default (Q201). Forgetting `public` is a frequent cause of "inaccessible base" errors.

**Interview angle:** *"When would you use private inheritance?"* "When you need a base's implementation, or to override its virtual functions, but don't want users to treat your class as that base, e.g. inheriting from a policy or helper class. Otherwise prefer composition."

---

## Batch 12 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 221 | M14 | E | MCQ | B | Well-defined |
| 222 | M14 | H | MCQ | C | **UB** |
| 223 | M14 | M | MCQ | D | Well-defined |
| 224 | M14 | E | MCQ | A | Well-defined |
| 225 | M14 | H | MCQ | A | Well-defined |
| 226 | M14 | M | MCQ | D | Well-defined |
| 227 | M14 | E | MCQ | B | Well-defined (leak) |
| 228 | M14 | M | MCQ | A | Well-defined |
| 229 | M14 | H | MSQ | A, B, E | Lifetime analysis |
| 230 | M14 | M | MSQ | A, D | Ill-formed / OK mix |
| 231 | M14 | H | MCQ | B | Ill-formed / OK mix |
| 232 | M14 | E | MCQ | B | Well-defined |
| 233 | M14 | M | MSQ | A, B, C | Indeterminate / zero mix |
| 234 | M14 | H | MCQ | C | Well-defined |
| 235 | M14 | M | MSQ | A, D | Ill-formed / OK mix |
| 236 | M14 | H | MCQ | A | Well-defined (`make2` unspecified between two forms) |
| 237 | M14 | H | MCQ | C | **Unspecified → UB** |
| 238 | M14 | M | MCQ | D | Well-defined |
| 239 | M15 | E | MCQ | B | Ill-formed |
| 240 | M15 | M | MCQ | C | Ill-formed / OK mix |

- **Difficulty:** 5 E / 8 M / 7 H.
  - M14 total (Q215–Q238): 5/10/9, as planned.
  - M15 so far: 1/1/0.
- **Format:** 16 MCQ / 4 MSQ.
  - M14 total: 19/5, as planned.
- **MCQ answer keys:** A×4, B×5, C×4, D×3. Running total Q1–Q240: A 49, B 49, C 48, D 47.

---

**Completed: 240 / 400 Remaining: 160**

**Ledger, newly ✓ this batch:**
- P58 ×3 (Q221, Q222, Q223) ★ complete
- P56 → 2 (Q235)
- P71 → 3 (Q236, Q238)
- P99 ×1 (Q237)
- P57 ×2 (Q239, Q240) ★ complete (+1 more planned)

**M14 complete** (Q215–Q238, 24 questions). **M15 in progress** (2 of 16 done).

*Say **"Continue."** for Batch 13 (Q241–Q260: M15 continued — inherited access rules, constructor and destructor order with multiple inheritance, name hiding and `using`, ambiguity in multiple inheritance, the diamond problem and virtual inheritance, upcasting and downcasting, object layout — then M16 polymorphism begins).*
