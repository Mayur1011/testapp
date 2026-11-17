# C/C++ MCQ/MSQ Bank — Batch 13 (Q241–Q260)

**What this batch covers:**
- **M15 Inheritance, finishing (Q241–Q254):**
  - access through two levels of inheritance (protected inheritance, then public)
  - constructor and destructor order: multilevel, multiple inheritance, and passing arguments up to a base constructor
  - name hiding, and `using` to bring base overloads back
  - ambiguity in multiple inheritance, which name lookup reports before overload resolution
  - the diamond without and with `virtual` inheritance
  - which class constructs a virtual base
  - upcasting and downcasting rules
  - how the pointer value changes when converting to the second base
- **M16 Polymorphism, starting (Q255–Q260):**
  - static vs dynamic binding, and qualified calls that skip dispatch
  - virtual dispatch through an array of base pointers
  - inheriting an override through an intermediate class
  - a signature mismatch that silently *doesn't* override
  - the template-method pattern, with a *private* override

**Patterns:** P57, P62, P63, P64, P65

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - Object traces show every constructor and destructor call; dispatch traces explain the vtable/vptr reasoning.

> **Conventions for every question in this batch:**
> - **C++17**, exceptions enabled.
> - Platform: a typical 64-bit system.
> - **Layout (typical ABIs):** base-class subobjects are laid out in declaration order, followed by the derived class's own members. Each polymorphic object holds one hidden vptr.
> - Needed headers are included. In code, `cout` means `std::cout`.

---

## Q241. [M15 | Access through two levels of inheritance | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P57

```cpp
class B { public: int a; protected: int b; private: int c; };
class D : protected B {};
class E : public D { void f(); };
int main() { D d; E e; }
```

Select **all** accesses that **compile**.

- A. Inside `E::f`, using `a`
- B. Inside `E::f`, using `b`
- C. Inside `E::f`, using `c`
- D. Inside `main`, using `d.a`
- E. Inside `main`, using `e.a`

**Answer:** A, B

**In simple words:** Each inheritance step can **lower the visibility** of what's passed down, never raise it. `D : protected B` turns B's public `a` into *protected* inside D. Protected members are still visible to D's own children (like E), but not to outsiders like `main`. B's *private* `c` is never visible to any derived class.

**Why each option:**

| Member | In B | In D (protected inheritance) | In E (public from D) | Inside `E::f`? | From `main`? |
|---|---|---|---|---|---|
| a | public | **protected** | protected | ✓ (A) | ✗ (D, E) |
| b | protected | protected | protected | ✓ (B) | ✗ |
| c | private | inaccessible | inaccessible | ✗ (C) | ✗ |

**Concept recap:** The effective access is the **more restrictive** of (the member's access in the base, the inheritance mode).

| Member in base | `public` inheritance | `protected` inheritance | `private` inheritance |
|---|---|---|---|
| public | public | protected | private |
| protected | protected | protected | private |
| private | inaccessible | inaccessible | inaccessible |

**Trap alert:** "Inaccessible" isn't "absent": `c` still occupies space in every D and E object.

**Interview angle:** *"Draw the access table for inheritance modes."* This 3×3 table is asked very often. Add: "Public inheritance is the only one that models is-a."

---

## Q242. [M15 | Construction and destruction order in multilevel inheritance | Easy | MCQ | [C++] · C++17 | OA] · Type T7

```cpp
#include <iostream>
using std::cout;
struct A { A() { cout << "A"; } ~A() { cout << "a"; } };
struct B : A { B() { cout << "B"; } ~B() { cout << "b"; } };
struct C : B { C() { cout << "C"; } ~C() { cout << "c"; } };
int main() { { C obj; } cout << '\n'; }
```

What does this program print?

- A. `CBAabc`
- B. `ABCabc`
- C. `CBAcba`
- D. `ABCcba`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Build from the **ground up** (the grandparent first, then the parent, then the child). Tear down from the **top down** (the child first, then the parent, then the grandparent), exactly the reverse.

**Why:**

| Phase | Order | Prints |
|---|---|---|
| construct | A → B → C | ABC |
| destroy | ~C → ~B → ~A | cba |

**Why the other options fail:**
- A, C: Construct the most-derived class first.
- B: Destroys in the same order as construction.

**Concept recap:** Each constructor first runs its base's constructor (explicitly listed or by default), then its own members, then its body. Destructors do the reverse: body, then members, then the base.

**Trap alert:** If `C`'s constructor throws after `B` is fully built, then `B` and `A` *are* destroyed (in reverse), but `~C` is **not** called, because C was never completely constructed (M19).

**Interview angle:** *"Why base first?"* "A derived constructor may use base members, so the base must already be valid. In reverse, the derived destructor may still use base members, so the base must be destroyed last."

---

## Q243. [M15 | Order with multiple bases and a member | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
using std::cout;
struct X { X() { cout << "X"; } ~X() { cout << "x"; } };
struct Y { Y() { cout << "Y"; } ~Y() { cout << "y"; } };
struct Z { Z() { cout << "Z"; } ~Z() { cout << "z"; } };
struct M : Y, X {
    Z z;
    M() : X(), Y() { cout << "M"; }
    ~M() { cout << "m"; }
};
int main() { { M m; } cout << '\n'; }
```

What does this program print?

- A. `XYZMmzyx`
- B. `ZYXMmxyz`
- C. `YXZMmzxy`
- D. `YXZMmxyz`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Bases are built in the order they're **listed after the colon** in the class header (`: Y, X`), not the order in the constructor's initializer list. Then members, then the body. Destruction reverses everything.

**Why:**

| Step | What | Prints |
|---|---|---|
| 1 | bases in **base-specifier order**: Y, then X | Y X |
| 2 | members in declaration order: z | Z |
| 3 | body | M |
| destroy | body, then z, then bases in reverse (X, then Y) | m z x y |

Output: `YXZMmzxy`.

**Why the other options fail:**
- A: Uses the initializer-list order for the bases.
- B: Builds members before bases.
- D: Destroys the bases in construction order.

**Concept recap:** The same "declaration order wins" rule as for members (Q221), applied to bases. Compilers warn when the initializer list's order differs.

**Trap alert:** Changing `struct M : Y, X` to `struct M : X, Y` changes the construction order, and possibly behaviour, if one base's constructor depends on another (it shouldn't).

**Interview angle:** *"In what order are multiple base classes constructed?"* "In the order they appear in the class's base list. Virtual bases come before all of them (Q250–Q251)."

---

## Q244. [M15 | Passing arguments up to base constructors | Easy | MCQ | [C++] · C++17 | OA] · Type T7

```cpp
#include <iostream>
using std::cout;
struct A { A(int x) { cout << "A" << x; } };
struct B : A { B(int x) : A(x * 2) { cout << "B" << x; } };
struct C : B {
    C() : B(5) { cout << "C"; }
    C(int y) : B(y + 1) { cout << "C" << y; }
};
int main() {
    C c1;
    cout << "|";
    C c2(1);
    cout << '\n';
}
```

What does this program print?

- A. `A10B5C|A4B2C1`
- B. `A5B5C|A2B2C1`
- C. `C|C1`
- D. `A10B5C|A2B1C1`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** A derived constructor **chooses what to pass up** to its base, in its initializer list. Each level can transform the value on the way up. Then the output is printed bottom-up: the base finishes first.

**Why:**

| Object | Calls | A receives | Output |
|---|---|---|---|
| `c1` | `C()` → `B(5)` → `A(5 × 2)` | 10 | `A10` `B5` `C` |
| `c2(1)` | `C(1)` → `B(1 + 1)` → `A(2 × 2)` | 4 | `A4` `B2` `C1` |

**Why the other options fail:**
- B: Forgets B's `x * 2`.
- C: Ignores the base constructors.
- D: Forgets C's `y + 1`.

**Concept recap:** If a base has **no default constructor**, every derived constructor **must** call one of its constructors explicitly in the initializer list. Otherwise it's a compile error ("no matching function for call to `A::A()`").

**Trap alert:** You can't initialise a **base's member** directly from the derived class (`C() : x(5)` where `x` belongs to A is an error). Pass the value to the base's constructor instead.

**Interview angle:** *"How do you pass parameters to a base class constructor?"* "Through the derived constructor's initializer list, `Derived(int v) : Base(v)`. C++11 also lets you inherit the base's constructors with `using Base::Base;`."

---

## Q245. [M15 | Name hiding across inheritance | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T2 · P63

```cpp
#include <iostream>
using std::cout;
struct B {
    void f(int)    { cout << "B:int "; }
    void f(double) { cout << "B:double "; }
};
struct D : B {
    void f(const char *) { cout << "D:str "; }
};
int main() {
    D d;
    d.f("x");     // L1
    d.f(1);       // L2
    d.B::f(1);    // L3
}
```

Which line fails to compile?

- A. L1
- B. L3
- C. L1 and L2
- D. L2

**Answer:** D

**Behaviour class:** Ill-formed (L2: "invalid conversion from 'int' to 'const char*'").

**In simple words:** When C++ looks up `d.f`, it searches **D first, and stops at the first class that has any `f` at all**. D has `f(const char*)`, so the search stops there, and B's `f(int)` and `f(double)` are never considered. They're **hidden**, even though they'd be a better match. Passing `1` to `f(const char*)` doesn't work, so L2 is an error.

**Why:**

| Line | Name lookup finds | Overload resolution | Result |
|---|---|---|---|
| L1 | `D::f(const char*)` only | exact | ✓ |
| L2 | `D::f(const char*)` only (B's `f`s are hidden) | int → const char*: no conversion | ✗ |
| L3 | qualified `B::f` → both B overloads | f(int) | ✓ |

**Why the other options fail:**
- A, C: L1 matches exactly.
- B: Qualified names bypass hiding.

**Concept recap:** **Hiding** happens by **name**, not by signature: any member named `f` in the derived class hides **every** `f` in the base. This applies to non-virtual and virtual functions alike.

**Trap alert:** This often bites when adding one new overload in a derived class: suddenly all the base's overloads "disappear". The fix is Q246's `using B::f;`.

**Interview angle:** *"Overloading vs overriding vs hiding?"*
- **Overloading:** same scope, same name, different parameters.
- **Overriding:** a derived class's virtual function with the same signature replaces the base's in dispatch.
- **Hiding:** a derived class's name blocks lookup of the base's same-named members.

---

## Q246. [M15 | `using` to bring back base overloads | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T6 · P63

```cpp
#include <iostream>
using std::cout;
struct B {
    void f(int)    { cout << "Bi "; }
    void f(double) { cout << "Bd "; }
};
struct D : B {
    using B::f;
    void f(int) { cout << "Di "; }
};
int main() {
    D d;
    d.f(1);
    d.f(2.5);
    d.B::f(1);
    cout << '\n';
}
```

What does this program print?

- A. Compile error: `f(int)` is declared twice in `D`.
- B. `Di Bd Bi`
- C. `Bi Bd Bi`
- D. `Di Di Bi`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** `using B::f;` says **"also let B's versions of `f` be found from D"**, which undoes the hiding. If D *also* declares its own version with exactly the same parameters, D's own version takes priority over the one brought in from B. Overloads that D doesn't redefine (like `f(double)`) still come from B.

**Why:**

| Call | Candidates visible in D | Chosen | Prints |
|---|---|---|---|
| `d.f(1)` | D::f(int) (D's own beats the using'd B::f(int)), B::f(double) | D::f(int), exact | `Di` |
| `d.f(2.5)` | same | B::f(double), exact | `Bd` |
| `d.B::f(1)` | qualified: B's overloads only | B::f(int) | `Bi` |

**Why the other options fail:**
- A: A member declared in the class **hides or overrides** a using-declared member with the same signature ([namespace.udecl]p15). That's not a conflict.
- C: Ignores D's own `f(int)`.
- D: Thinks D's `f(int)` also catches `2.5`. An exact match with `f(double)` wins.

**Concept recap:** The idiom to add or replace **one** overload while keeping the others: `using Base::name;` plus your own declaration.

**Trap alert:** `using` can also change access: `using B::f;` in a `public:` section exposes protected base overloads publicly.

**Interview angle:** *"I overrode one virtual overload and now the others don't compile when called on a Derived object. Why?"* "Name hiding. Add `using Base::func;` in the derived class to bring the remaining overloads back into scope."

---

## Q247. [M15 | Ambiguity in multiple inheritance | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T2 · P64

```cpp
#include <iostream>
struct A { void hello() { std::cout << "A"; } int v = 1; };
struct B { void hello() { std::cout << "B"; } int v = 2; };
struct C : A, B {};
int main() {
    C c;
    c.hello();               // L1
    c.A::hello();            // L2
    std::cout << c.B::v;     // L3
}
```

Which line fails to compile?

- A. L2
- B. L2 and L3
- C. L1
- D. None; L1 calls `A::hello` because A is listed first.

**Answer:** C

**Behaviour class:** Ill-formed (L1: "request for member 'hello' is ambiguous").

**In simple words:** `C` inherited **two different `hello`s**, one from each parent, and neither parent comes first. `c.hello()` is like shouting a name that two people share: C++ refuses to guess. Saying which parent you mean (`c.A::hello()`) settles it.

**Why:**

| Line | Lookup | Result |
|---|---|---|
| L1 | finds `hello` in both A and B → ambiguous | ✗ |
| L2 | qualified → A's | ✓ prints A |
| L3 | qualified → B's `v` | ✓ prints 2 |

**Why the other options fail:**
- A, B: Qualified names are unambiguous.
- D: The order of the bases doesn't break ties.

**Concept recap:** Ways to resolve it:
- qualify the call: `c.A::hello()`;
- add `using A::hello;` in C to pick one;
- define `C::hello()` that calls one or both.

**Trap alert:** **Declaring** `C : A, B` is fine. The error only appears when the ambiguous name is **used**.

**Interview angle:** *"Problems with multiple inheritance?"* "Name ambiguities, the diamond problem (duplicate base subobjects; Q249), more complex construction order and layout, and pointer adjustment on casts (Q254). That's why many guidelines allow multiple inheritance mainly for interfaces (pure abstract bases)."

---

## Q248. [M15 | Lookup ambiguity happens before overload resolution | Hard | MSQ | [C++] · C++17 | Interview] · Type T2 · P64

```cpp
struct A { void f(int) {} };
struct B { void f(double) {} };
struct C : A, B {};
C c;
```

Select **all** true statements.

- A. `c.f(1);` is ambiguous, even though `A::f(int)` would be an exact match.
- B. Adding `using A::f; using B::f;` inside `C` makes `c.f(1)` call `A::f`.
- C. If `B::f` were `private`, `c.f(1)` would compile, because inaccessible members are ignored during lookup.
- D. `static_cast<A &>(c).f(1);` compiles.
- E. Declaring `struct C : A, B {}` is itself an error, because A and B have members with the same name.

**Answer:** A, B, D

**In simple words:** C++ answers two questions in order: first **"which class does the name `f` come from?"** (name lookup), and only then **"which overload fits best?"** (overload resolution). Here the first question has two answers (A and B), so C++ stops with "ambiguous" before ever comparing `int` with `double`. A `using` declaration moves both `f`s into C, so they become ordinary overloads in one scope, and the best match wins.

**Why each option:**
- **A ✓:** Lookup finds `f` in two different base classes → ambiguous ([class.member.lookup]). Overload resolution never runs.
- **B ✓:** Both are now declared in C's scope, so normal overload resolution picks `A::f(int)` (exact match).
- **C ✗:** Access is checked **after** lookup and overload resolution. Private names still take part in lookup, so it's still ambiguous.
- **D ✓:** Converting to `A&` makes A the static type, so lookup finds only `A::f`.
- **E ✗:** Only *using* an ambiguous name is an error.

**Concept recap:** The pipeline: **name lookup → overload resolution → access check → (virtual dispatch at run time)**. Many surprising errors come from the first step.

**Trap alert:** The same "lookup first" rule causes Q245's name hiding. In both cases the compiler stops at the first scope that has the name, before comparing types.

**Interview angle:** *"Why does `c.f(1)` fail when A has `f(int)` and B has `f(double)`?"* "Overload resolution only works within a single set of candidates found by name lookup. Names from different bases form an ambiguous lookup result. Merge them with `using` declarations."

---

## Q249. [M15 | The diamond without virtual inheritance | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P65

```cpp
#include <iostream>
using std::cout;
struct Base { int x = 0; Base() { cout << "B"; } };
struct L : Base {};
struct R : Base {};
struct J : L, R {};
int main() {
    J j;
    j.L::x = 1;
    j.R::x = 2;
    cout << " " << j.L::x << j.R::x << " " << sizeof(J) / sizeof(int) << '\n';
}
```

What does this program print?

- A. `B 22 1`
- B. `BB 22 2`
- C. `B 12 1`
- D. `BB 12 2`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Without `virtual`, **each** parent brings **its own separate copy** of the grandparent. So a `J` contains **two** `Base` objects (two `x`s), built twice. That's the "diamond problem": data is duplicated, and an unqualified `j.x` or `Base *p = &j` becomes ambiguous ("which Base do you mean?").

**Why:**

| Item | Value |
|---|---|
| `Base()` runs for L's copy and for R's copy | `BB` |
| `j.L::x = 1`, `j.R::x = 2` | two independent ints |
| print | `12` |
| `sizeof(J)` = two ints = 8 bytes → / 4 | `2` |

```
     Base   Base      ← two separate subobjects
      |      |
      L      R
       \    /
         J
```

**Why the other options fail:**
- A, C: Assume one shared Base. That needs `virtual` inheritance (Q250).
- B: Assumes the two writes hit the same `x`.

**Concept recap:** `j.x` and `Base *p = &j;` are both **ambiguous** here (compile errors). You must say `j.L::x`, or convert through `L` first: `Base *p = static_cast<L *>(&j);`.

**Trap alert:** Sometimes two copies are actually what you want (independent sub-parts). The diamond problem is about the case where you meant **one** shared grandparent.

**Interview angle:** *"What is the diamond problem?"* Describe the two copies and the ambiguity, then the fix: `virtual` inheritance (Q250), which makes the most-derived class own one shared base.

---

## Q250. [M15 | Virtual inheritance: one shared base | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P65

```cpp
#include <iostream>
using std::cout;
struct Base { int x = 0; Base() { cout << "B"; } };
struct L : virtual Base { L() { cout << "L"; } };
struct R : virtual Base { R() { cout << "R"; } };
struct J : L, R { J() { cout << "J"; } };
int main() {
    J j;
    j.x = 5;
    cout << " " << j.L::x << j.R::x << '\n';
}
```

What does this program print?

- A. `BLRJ 55`
- B. `BLBRJ 50`
- C. `LRBJ 55`
- D. Compile error: `j.x` is ambiguous

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** `virtual` inheritance means **"share the grandparent"**: however many paths lead to `Base`, a `J` contains exactly **one** `Base`. It's built **once, first of all**, before L and R, by the most-derived class. Since there's only one `x`, `j.x` is no longer ambiguous, and both paths see the same value.

**Why:**

| Step | Event | Prints |
|---|---|---|
| 1 | virtual base `Base` is constructed **first**, and only once | B |
| 2 | non-virtual bases in order: L, then R | L R |
| 3 | body of J | J |
| `j.x = 5` | the single shared `x` | |
| `j.L::x`, `j.R::x` | both name the same `x` | 5 5 |

```
        Base        ← one shared subobject
       /    \
      L      R
       \    /
         J
```

**Why the other options fail:**
- B: The non-virtual diamond (two copies).
- C: Virtual bases are constructed before any non-virtual base.
- D: With a single shared `Base`, `j.x` is unambiguous.

**Concept recap:** The cost: each class with a virtual base typically stores extra bookkeeping (via the vptr / a virtual-base offset) so the shared base can be found at run time. Objects get slightly bigger, and access to the base is slightly slower.

**Trap alert:** The real-world classic is `std::iostream`, which inherits from `istream` and `ostream`, each of which inherits **virtually** from `ios_base` (through `basic_ios`), so there's one shared stream state.

**Interview angle:** *"How does virtual inheritance solve the diamond problem, and what does it cost?"* "One shared base subobject, constructed by the most-derived class. The costs: extra indirection to reach it, larger objects, and more complex constructor rules (Q251)."

---

## Q251. [M15 | Who constructs the virtual base? | Hard | MCQ | [C++] · C++17 | Interview] · Type T7 · P65

```cpp
#include <iostream>
using std::cout;
struct Base { Base(int v) { cout << "B" << v; } };
struct L : virtual Base { L() : Base(1) { cout << "L"; } };
struct R : virtual Base { R() : Base(2) { cout << "R"; } };
struct J : L, R { J() : Base(3) { cout << "J"; } };
int main() { J j; cout << '\n'; }
```

What does this program print?

- A. `B1LB2RJ`
- B. `B1LRJ`
- C. `B2LRJ`
- D. `B3LRJ`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** With a shared (virtual) base, L and R might both try to build it, with different arguments. To avoid that conflict, C++ gives the job to **the most-derived class** (`J`, the object actually being created). L's `Base(1)` and R's `Base(2)` are simply **ignored** when they're building part of a `J`.

**Why:**

| Step | Who | Action | Prints |
|---|---|---|---|
| 1 | J (most-derived) | constructs the virtual base: `Base(3)` | B3 |
| 2 | L's constructor | its `Base(1)` initializer is **skipped** | L |
| 3 | R's constructor | its `Base(2)` initializer is **skipped** | R |
| 4 | J's body | | J |

**Why the other options fail:**
- A: Treats the base as non-virtual (built twice).
- B, C: Let an intermediate class choose the argument.

**Concept recap:** If `J`'s constructor did **not** mention `Base(...)`, the compiler would try `Base()`, which doesn't exist → **compile error**, even though L and R both mention `Base`. Every most-derived class must be able to construct the virtual base.

**Trap alert:** If you create an `L` on its own (`L l;`), then **L** is the most-derived class, and its `Base(1)` *is* used. The same constructor behaves differently depending on what's being built.

**Interview angle:** *"In virtual inheritance, which constructor initialises the virtual base?"* "Always the most-derived class's. Intermediate classes' initializers for the virtual base are ignored. That's why designs with virtual bases usually give them default constructors."

---

## Q252. [M15 | Upcasting and downcasting | Easy | MCQ | [C++] · C++17 | OA] · Type T5

```cpp
#include <iostream>
struct B { int b = 1; };
struct D : B { int d = 2; };
int main() {
    D obj;
    B *pb = &obj;
    D *pd = static_cast<D *>(pb);
    std::cout << pd->d << '\n';
}
```

What does this program print?

- A. `1`
- B. `2`
- C. Compile error: `static_cast` can't downcast
- D. Undefined behaviour: `pb` is a `B*`

**Answer:** B

**Behaviour class:** Well-defined: `pb` really points to a `D`.

**In simple words:** **Upcasting** (D* → B*) is always safe and automatic: every D *is a* B. **Downcasting** (B* → D*) is a promise you make: "I know this B is actually part of a D". `static_cast` takes your word for it without checking. Here the promise is true, so it works.

**Why:**

| Step | Cast | Safe? |
|---|---|---|
| `B *pb = &obj;` | upcast, implicit | always |
| `static_cast<D *>(pb)` | downcast, unchecked | only because `*pb` really is part of a `D` |
| `pd->d` | | 2 |

**Why the other options fail:**
- A: Reads the base member.
- C: `static_cast` *does* allow downcasts along a non-virtual inheritance chain.
- D: The cast is valid because the object's real type is D.

**Concept recap:** If `pb` had pointed to a plain `B` (`B base; pb = &base;`), the same `static_cast` would compile, but using the result as a D would be **UB** (Q253-C). For a checked downcast, use `dynamic_cast` (which needs a polymorphic base; M16).

**Trap alert:** `static_cast` can't downcast from a **virtual** base (the offset isn't known at compile time). Only `dynamic_cast` can.

**Interview angle:** *"static_cast vs dynamic_cast for downcasting?"* "`static_cast`: no run-time check, zero cost, UB if you're wrong. `dynamic_cast`: checked at run time (returns `nullptr` or throws `bad_cast`), needs a polymorphic type, and has a cost."

---

## Q253. [M15 | Rules for downcasting | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2

Given a **non-polymorphic** hierarchy `struct B { int b; }; struct D : B { int d; };` and objects `B b; D d;`, select **all** true statements.

- A. `B *p = &d;` compiles.
- B. `D *q = &b;` compiles.
- C. `D *r = static_cast<D *>(&b);` compiles, but accessing `r->d` is undefined behaviour.
- D. `D *s = dynamic_cast<D *>(&b);` compiles and yields `nullptr`.
- E. `B &rb = d; D &rd = static_cast<D &>(rb);` is well-defined.

**Answer:** A, C, E

**In simple words:** Going **up** is automatic. Going **down** needs a cast, and with `static_cast` it's on your honour: lie about the real type (as in C) and anything can happen. `dynamic_cast` can check the real type, but only for classes with virtual functions; B has none, so `dynamic_cast` won't even compile here.

**Why each option:**
- **A ✓:** Derived → base is an implicit conversion.
- **B ✗:** There's no implicit downcast ("invalid conversion from 'B*' to 'D*'").
- **C ✓:** `static_cast` compiles (it trusts you), but `b` isn't a D, so `r->d` accesses memory that doesn't belong to any D → UB.
- **D ✗:** `dynamic_cast` for downcasting requires a **polymorphic** source type (at least one virtual function). Here it's ill-formed ("source type is not polymorphic").
- **E ✓:** `rb` actually refers to a D, so casting it back is valid.

**Concept recap:**

| Cast | Direction | Checked? | Requirement |
|---|---|---|---|
| implicit | up | — | accessible, unambiguous base |
| `static_cast` | down | no | the object really is a Derived, or UB |
| `dynamic_cast` | down / across | yes, at run time | polymorphic base |

**Trap alert:** `dynamic_cast` **upcasts** don't need polymorphism. Only down- and cross-casts do.

**Interview angle:** *"How do you safely downcast in C++?"* "`dynamic_cast` on a polymorphic base, checking for `nullptr`. Better still, redesign so you don't need to: add a virtual function to the base."

---

## Q254. [M15 | The pointer value changes when converting to the second base | Hard | MCQ | [C++] · C++17 | Interview] · Type T5

Typical layout: base subobjects in declaration order, followed by the class's own members.

```cpp
#include <iostream>
struct A { int a = 1; };
struct B { int b = 2; };
struct C : A, B { int c = 3; };
int main() {
    C obj;
    A *pa = &obj;
    B *pb = &obj;
    std::cout << (char *)pb - (char *)&obj << ' '
              << (pa == &obj) << ' ' << (pb == &obj) << '\n';
}
```

What does this program print?

- A. `0 1 1`
- B. `4 1 0`
- C. `4 1 1`
- D. `0 1 0`

**Answer:** C

**Behaviour class:** Well-defined. The offset value is ABI-dependent, and fixed here by the stated layout.

**In simple words:** A `C` object is laid out as **[A part][B part][c]**. A pointer to the A part starts at the same address as the whole object, but a pointer to the B part has to point **4 bytes in**. C++ adjusts the address **automatically** when converting `C*` to `B*`. And when comparing `pb == &obj`, it converts `&obj` to `B*` first (with the same adjustment), so the two still compare equal.

**Why:**

| Subobject | Offset |
|---|---|
| A (a) | 0 |
| B (b) | 4 |
| c | 8 |

| Expression | Value |
|---|---|
| `(char*)pb - (char*)&obj` | **4** (the pointer was adjusted) |
| `pa == &obj` | `&obj` → `A*` (+0) → equal → **1** |
| `pb == &obj` | `&obj` → `B*` (+4) → equal → **1** |

**Why the other options fail:**
- A, D: Assume no adjustment for the second base.
- B: Forgets that the comparison converts both sides to the same type.

**Concept recap:** That's why converting through `void*` is dangerous with multiple inheritance: `(B*)(void*)&obj` skips the adjustment and points at the A part. Always use `static_cast` / `dynamic_cast` between related class types, never via `void*` or `reinterpret_cast`.

**Trap alert:** `reinterpret_cast<B *>(&obj)` also skips the adjustment, and using the result is wrong (UB). `static_cast<B *>(&obj)` adjusts correctly.

**Interview angle:** *"Why can `static_cast<Base2*>(derived_ptr)` change the numeric address?"* "With multiple inheritance, each base subobject lives at its own offset inside the derived object, and the cast adds that offset. A null pointer stays null."

---

## Q255. [M16 | Static vs dynamic binding | Easy | MCQ | [C++] · C++17 | OA] · Type T7 · P62

```cpp
#include <iostream>
using std::cout;
struct B { void f() { cout << "Bf "; } virtual void g() { cout << "Bg "; } };
struct D : B { void f() { cout << "Df "; } void g() override { cout << "Dg "; } };
int main() {
    D d;
    B *p = &d;
    p->f();
    p->g();
    d.f();
    cout << '\n';
}
```

What does this program print?

- A. `Df Dg Df`
- B. `Bf Bg Df`
- C. `Df Bg Bf`
- D. `Bf Dg Df`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:**
- For a **non-virtual** function, C++ decides **at compile time**, using the pointer's declared type (B*), so `p->f()` calls B's version.
- For a **virtual** function, C++ decides **at run time**, using the actual object (a D), so `p->g()` calls D's version.

It's like a name badge vs a fingerprint check: non-virtual calls read the badge (the pointer type), virtual calls check who's really there.

**Why:**

| Call | Virtual? | Decided by | Calls |
|---|---|---|---|
| `p->f()` | no | static type `B*` → **static binding** | B::f |
| `p->g()` | yes | dynamic type D → **dynamic binding** | D::g |
| `d.f()` | no | static type `D` | D::f |

**Why the other options fail:**
- A: Treats `f` as virtual.
- B: Treats `g` as non-virtual.
- C: Swaps the rules.

**Concept recap:** Dynamic dispatch needs **two things**: (1) the function is virtual, and (2) the call goes through a **pointer or reference**. Calling on an object by value uses its static type (which, for a plain object, *is* its real type).

**Trap alert:** A non-virtual `D::f` with the same name and signature as `B::f` doesn't override it; it **hides** it (Q245). Behaviour then depends on the pointer type, which is usually a bug.

**Interview angle:** *"Static vs dynamic binding?"* "Static (early) binding: resolved at compile time from the declared type, for non-virtual functions, overloads and templates. Dynamic (late) binding: resolved at run time via the vtable, for virtual functions called through a pointer or reference."

---

## Q256. [M16 | A qualified call skips virtual dispatch | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P62

```cpp
#include <iostream>
using std::cout;
struct B { virtual void g() { cout << "Bg "; } virtual ~B() = default; };
struct D : B { void g() override { cout << "Dg "; } };
int main() {
    D d;
    B &r = d;
    r.g();
    r.B::g();
    d.g();
    cout << '\n';
}
```

What does this program print?

- A. `Dg Bg Dg`
- B. `Dg Dg Dg`
- C. `Bg Bg Dg`
- D. Compile error: `r.B::g()` on a reference to a `D`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** A **reference** behaves like a pointer for virtual calls, so `r.g()` reaches D's version. Writing the class name explicitly, as in `r.B::g()`, says **"call exactly B's version, don't look up the real type"**. Overrides use this to call the base version as part of their own work.

**Why:**

| Call | Mechanism | Calls |
|---|---|---|
| `r.g()` | virtual call through a reference → dynamic dispatch | D::g |
| `r.B::g()` | **qualified name** → no dispatch (static call) | B::g |
| `d.g()` | object of static type D | D::g |

**Why the other options fail:**
- B: Assumes qualification can't bypass dispatch.
- C: Assumes references don't dispatch.
- D: Qualified calls to an accessible base version are legal.

**Concept recap:** The common pattern: `void D::g() override { B::g(); /* extra work */ }`, which extends rather than replaces the base behaviour.

**Trap alert:** Calling a **pure** virtual function with a qualified name (`B::g()`) is a direct call. It's legal **only if** the pure virtual function also has a definition. Otherwise it's a link error (M16 later).

**Interview angle:** *"How do you call the base-class version of an overridden function?"* "Qualify the call: `Base::func()`. That turns off virtual dispatch for that call."

---

## Q257. [M16 | Virtual dispatch through an array of base pointers | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P62

```cpp
#include <iostream>
using std::cout;
struct Shape {
    virtual double area() const { return 0; }
    virtual const char *name() const { return "shape"; }
    virtual ~Shape() = default;
};
struct Sq : Shape {
    double s; Sq(double s) : s(s) {}
    double area() const override { return s * s; }
    const char *name() const override { return "sq"; }
};
struct Ci : Shape {
    double r; Ci(double r) : r(r) {}
    double area() const override { return 3 * r * r; }
};
int main() {
    Sq a(2); Ci b(1); Shape c;
    const Shape *arr[] = {&a, &b, &c};
    double total = 0;
    for (auto p : arr) { total += p->area(); cout << p->name() << ' '; }
    cout << total << '\n';
}
```

What does this program print?

- A. `sq ci shape 7`
- B. `shape shape shape 0`
- C. `sq shape shape 7`
- D. `sq shape shape 4`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Each object carries a hidden link to **its own class's table of functions** (the vtable). Calling a virtual function through a base pointer follows that link. If a class didn't override a function (like `Ci::name`), its table just points to the **inherited** version.

**Why — the vtables:**

| Class | `area` slot | `name` slot |
|---|---|---|
| Shape | Shape::area | Shape::name |
| Sq | **Sq::area** | **Sq::name** |
| Ci | **Ci::area** | Shape::name (not overridden) |

| p → | area | name | total |
|---|---|---|---|
| a (Sq) | 4 | "sq" | 4 |
| b (Ci) | 3 | "shape" | 7 |
| c (Shape) | 0 | "shape" | 7 |

Output: `sq shape shape 7`.

**Why the other options fail:**
- A: Ci never defines `name`. There's no "ci" string anywhere.
- B: Ignores dispatch.
- D: Omits Ci's area.

**Concept recap:** **vptr / vtable:** each polymorphic class has one table of function pointers (the vtable), and each object holds a pointer to its class's table (the vptr), set by the constructor. A virtual call is: load the vptr, load the slot, make an indirect call.

**Trap alert:** Storing the shapes **by value** in an array (`Shape arr[] = {a, b, c};`) would **slice** them into plain Shapes, and every call would print "shape" (P61, later in M16). Polymorphism needs pointers or references.

**Interview angle:** *"Explain how virtual functions work internally."* Draw the vtable table above, then explain that the vptr is set during construction. Mention the cost (one indirect call, which is hard to inline) and that the standard doesn't mandate vtables, but every mainstream compiler uses them.

---

## Q258. [M16 | Inheriting an override through an intermediate class | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P62

```cpp
#include <iostream>
using std::cout;
struct A { virtual void f() { cout << "A"; } virtual ~A() = default; };
struct B : A {};
struct C : B { void f() override { cout << "C"; } };
struct D : C {};
int main() {
    A *p = new D; p->f();
    A *q = new B; q->f();
    delete p; delete q;
    cout << '\n';
}
```

What does this program print?

- A. `AA`
- B. `CA`
- C. `CC`
- D. Compile error: `D` doesn't override `f`

**Answer:** B

**Behaviour class:** Well-defined. The virtual destructor makes the `delete` through `A*` correct.

**In simple words:** Each class's vtable **starts as a copy of its parent's**, then replaces the slots it overrides. B overrides nothing (it keeps A's `f`). C replaces `f`. D overrides nothing, so it **inherits C's** replacement. A call always goes to the nearest override on the path from the object's real class back up to the base.

**Why:**

| Class | vtable `f` slot |
|---|---|
| A | A::f |
| B | A::f (inherited) |
| C | **C::f** |
| D | C::f (inherited) |

- `p` → a D object → D's slot → **C::f** → `C`
- `q` → a B object → B's slot → **A::f** → `A`

**Why the other options fail:**
- A: Ignores C's override for D objects.
- C: Thinks B also gets C's override. Overrides flow **down** the hierarchy, never up.
- D: Derived classes aren't required to override.

**Concept recap:** The **final overrider** of a virtual function, for an object of type T, is the most-derived override on T's inheritance path.

**Trap alert:** Without `virtual ~A()`, the `delete p` (A* → D object) would be **UB** (P60, later in M16). That's why the code includes it.

**Interview angle:** *"If an intermediate class doesn't override a virtual function, which version runs?"* "The nearest one above it in the hierarchy. The vtable slot is inherited."

---

## Q259. [M16 | A signature mismatch that doesn't override | Hard | MSQ | [C++] · C++17 | OA+Interview] · Type T7 · P62, P63

```cpp
#include <iostream>
using std::cout;
struct B {
    virtual void f(int)   { cout << "B "; }
    virtual void g() const { cout << "Bg "; }
    virtual ~B() = default;
};
struct D : B {
    void f(int) { cout << "D "; }
    void g()    { cout << "Dg "; }
};
```

With `B *p = new D;`, select **all** true statements.

- A. `p->f(1)` calls `D::f`.
- B. `p->g()` calls `D::g`.
- C. `D::f` is virtual, even though D doesn't write the `virtual` keyword.
- D. Writing `void g() override` in D would be a compile error.
- E. `D d; static_cast<B &>(d).g();` calls `D::g`.

**Answer:** A, C, D

**In simple words:** To override, the derived function must match the base's **exact signature, including `const`**. `D::f(int)` matches, so it overrides, and it's automatically virtual. `D::g()` is missing the `const`, so it's a **different function** that merely hides the base's name. Calls through a `B` still reach `B::g`. The `override` keyword exists to catch exactly this mistake at compile time.

**Why each option** (A and B verified: the output is `D Bg`):
- **A ✓:** The same signature → it overrides → dynamic dispatch picks `D::f`.
- **B ✗:** `D::g()` (non-const) doesn't override `B::g() const`, so the vtable slot still holds `B::g`.
- **C ✓:** "Once virtual, always virtual": a function with an overriding signature is virtual whether or not it says so.
- **D ✓:** `override` makes the compiler check: "marked 'override', but does not override". It's a compile error, which is exactly the point.
- **E ✗:** Through a `B&`, `g` resolves to the vtable slot → `B::g`.

**Concept recap:** Things that silently break overriding: differences in `const`, parameter types (`int` vs `long`), reference qualifiers (`&`/`&&`), or a misspelled name. Always write `override`.

**Trap alert:** Even the **return type** must match (or be covariant, P63/M16 later). A different non-covariant return type with the same parameters is a **compile error**, not a silent new function.

**Interview angle:** *"Why use the `override` keyword?"* "It turns silent non-overrides (a const mismatch, a parameter typo) into compile errors, and documents intent. Pair it with `final` to forbid further overriding."

---

## Q260. [M16 | The template-method pattern, with a private override | Hard | MCQ | [C++] · C++17 | Interview] · Type T7 · P62

```cpp
#include <iostream>
using std::cout;
struct B {
    void run() { cout << "["; step(); cout << "]"; }
    virtual void step() { cout << "b"; }
    virtual ~B() = default;
};
struct D : B {
private:
    void step() override { cout << "d"; }
};
int main() {
    D d; B b;
    d.run();
    b.run();
    B &r = d;
    r.run();
    cout << '\n';
}
```

What does this program print?

- A. Compile error: `D::step` is private
- B. `[b][b][b]`
- C. `[d][b][b]`
- D. `[d][b][d]`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** `run()` is an ordinary (non-virtual) function, but **inside** it, the call `step()` is really `this->step()`, which *is* virtual. So it dispatches to the real object's version. And **access is checked against the name you use**, not the function that ends up running: `B::step` is public, so the call is allowed, even though D's version is private.

**Why:**

| Call | `this` points to | `step()` dispatches to | Output |
|---|---|---|---|
| `d.run()` | D | D::step (private, but reached via B's public name) | `[d]` |
| `b.run()` | B | B::step | `[b]` |
| `r.run()` | D (through a B&) | D::step | `[d]` |

**Why the other options fail:**
- A: Access is checked at the call site against the **static** type (`B::step`, which is public). The override's own access doesn't matter for virtual calls through the base ([class.access.virt]).
- B: Ignores the virtual call inside `run`.
- C: Thinks calling through a `B&` stops dispatch inside `run`.

**Concept recap:** This is the **template method / non-virtual interface (NVI)** pattern: a public non-virtual function fixes the overall algorithm (`[`, step, `]`), and private virtual hooks let derived classes customise the steps.

**Trap alert:** Calling `d.step()` directly **would** fail to compile, because the static type is D, where `step` is private. Same function, different rule, because access depends on the name used.

**Interview angle:** *"What is the Non-Virtual Interface idiom?"* "Public non-virtual functions call private (or protected) virtual hooks. The base controls pre- and post-conditions, logging and locking, while derived classes override only the customisable steps. Herb Sutter's guideline: 'make virtual functions private'."

---

## Batch 13 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 241 | M15 | M | MSQ | A, B | Ill-formed / OK mix |
| 242 | M15 | E | MCQ | D | Well-defined |
| 243 | M15 | M | MCQ | C | Well-defined |
| 244 | M15 | E | MCQ | A | Well-defined |
| 245 | M15 | M | MCQ | D | Ill-formed (L2) |
| 246 | M15 | H | MCQ | B | Well-defined |
| 247 | M15 | M | MCQ | C | Ill-formed (L1) |
| 248 | M15 | H | MSQ | A, B, D | Lookup rules |
| 249 | M15 | M | MCQ | D | Well-defined |
| 250 | M15 | H | MCQ | A | Well-defined |
| 251 | M15 | H | MCQ | D | Well-defined |
| 252 | M15 | E | MCQ | B | Well-defined |
| 253 | M15 | M | MSQ | A, C, E | Ill-formed / UB / OK mix |
| 254 | M15 | H | MCQ | C | Well-defined (ABI offset) |
| 255 | M16 | E | MCQ | D | Well-defined |
| 256 | M16 | M | MCQ | A | Well-defined |
| 257 | M16 | M | MCQ | C | Well-defined |
| 258 | M16 | M | MCQ | B | Well-defined |
| 259 | M16 | H | MSQ | A, C, D | Well-defined / ill-formed (D) |
| 260 | M16 | H | MCQ | D | Well-defined |

- **Difficulty:** 4 E / 9 M / 7 H.
  - M15 total (Q239–Q254): 4/7/5, as planned.
  - M16 so far: 1/3/2.
- **Format:** 16 MCQ / 4 MSQ.
  - M15 total: 13/3, as planned.
- **MCQ answer keys:** A×3, B×3, C×4, D×6. Running total Q1–Q260: A 52, B 52, C 52, D 53.

---

**Completed: 260 / 400 Remaining: 140**

**Ledger, newly ✓ this batch:**
- P57 → 3 (Q241)
- P63 → 5 (Q245, Q246, Q259) ★ complete
- P64 ×2 (Q247, Q248)
- P65 ×3 (Q249, Q250, Q251) ★ complete
- P62 ×5 (Q255–Q260) ★ complete

**M15 complete** (Q239–Q254, 16 questions). **M16 in progress** (6 of 26 done).

*Say **"Continue."** for Batch 14 (Q261–Q280: the rest of M16 — `override`/`final`, more overriding vs overloading vs hiding, virtual destructors and the UB without one, pure virtual functions and abstract classes, a pure virtual destructor that needs a definition (a link error), virtual calls inside constructors and destructors, covariant return types, default arguments with virtual functions, object slicing, `dynamic_cast` across multiple inheritance).*
