# C/C++ MCQ/MSQ Bank — Batch 14 (Q261–Q280)

**What this batch covers (all M16 Polymorphism):**
- `final` and `override`: what they forbid and what they enable
- overloading vs overriding vs hiding:
  - virtual overloads where only one is overridden
  - overload resolution uses the *static* type
  - a full classification exercise
- virtual destructors:
  - the correct case
  - the UB when the base destructor isn't virtual
  - the rules, including smart-pointer behaviour
- abstract classes and pure virtual functions:
  - a pure virtual function that has a body
  - a pure virtual destructor without a definition (a link error)
- virtual calls inside constructors and destructors:
  - the order of calls
  - why dispatch stops at the class currently being built
  - a pure virtual call during construction (UB)
- covariant return types
- default arguments on virtual functions
- object slicing, including partial assignment through a base reference
- `dynamic_cast` cross-casts and the rules of `dynamic_cast`

**Patterns:** P42, P60, P61, P62, P63, P66, P72

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - Dispatch traces explain the vtable/vptr reasoning.

> **Conventions for every question in this batch:**
> - **C++17**, exceptions enabled.
> - Platform: a typical 64-bit system. Polymorphic objects hold one 8-byte vptr (the typical ABI).
> - Needed headers are included. In code, `cout` means `std::cout`.

---

## Q261. [M16 | final and override | Easy | MCQ | [C++] · C++17 | OA] · Type T2

```cpp
struct B {
    virtual void f() {}
    virtual void g() final {}
};
struct D final : B {
    void f() override {}     // L1
    void g() override {}     // L2
};
struct E : D {};             // L3
```

Which lines fail to compile?

- A. Only L2
- B. L2 and L3
- C. L1, L2 and L3
- D. Only L3

**Answer:** B

**Behaviour class:** Ill-formed (L2, L3).

**In simple words:**
- `override` on a function says **"I promise this replaces a base virtual function"**, and the compiler checks the promise.
- `final` on a function says **"nobody below me may replace this"**.
- `final` on a class says **"nobody may inherit from me"**.

**Why:**

| Line | Rule | Result |
|---|---|---|
| L1 | `f` is virtual in B, with the same signature → a valid override | ✓ |
| L2 | `B::g` is `final` → can't be overridden | ✗ "overriding final function" |
| L3 | `D` is `final` → can't be derived from | ✗ "cannot derive from 'final' base" |

**Why the other options fail:**
- A: Misses that `final` on a class blocks inheritance.
- C: L1 is a correct override.
- D: Misses the final function.

**Concept recap:** `final` documents design intent ("this behaviour is fixed"), and also lets the compiler **devirtualise** calls made through a `D*` / `D&`, because no further override can exist.

**Trap alert:** `override` and `final` are *contextual* keywords: special only in these positions. `int final = 0;` is still legal.

**Interview angle:** *"Why mark a class final?"* "To prevent inheritance from types not designed for it (e.g. ones with no virtual destructor), and to enable devirtualisation as an optimisation."

---

## Q262. [M16 | What final and override do | Medium | MSQ | [C++] · C++17 | Interview] · Type — concept

Select **all** true statements.

- A. `override` is a contextual keyword, so a variable may be named `override`.
- B. `final` applied to a class prevents any class from deriving from it.
- C. Marking a virtual function `final` lets the compiler turn calls through a pointer to that class into direct (non-virtual) calls.
- D. `override` may be used on a function that hides a **non-virtual** base function.
- E. A function marked `final` must also be written with the `virtual` or `override` keyword.

**Answer:** A, B, C

**In simple words:** These keywords are **safety labels**. They don't change *how* virtual calls work. They make the compiler check your intent (`override`), stop further customisation (`final`), and, as a bonus, let it skip the vtable lookup when it knows nothing can override further.

**Why each option:**
- **A ✓:** They're only special in declarator positions, so older code that uses them as names still compiles.
- **B ✓:** Q261.
- **C ✓:** With a `D*` where `D::f` is final, the only possible target is `D::f`, so a direct (even inlined) call is allowed.
- **D ✗:** `override` requires that the function actually overrides a **virtual** function. Otherwise it's a compile error.
- **E ✗:** `final` requires the function to be *virtual*, but it can get that by overriding: `void f() final;` alone is enough when the base declares `f` virtual.

**Concept recap:** Good style:
- write `virtual` only in the **base** declaration;
- write `override` (or `final`) in every overriding declaration;
- don't repeat `virtual` in derived classes.

**Trap alert:** `virtual void f() final;` in a class that has **no** base version of `f` is legal, but pointless: a virtual function that nobody can override.

**Interview angle:** *"Does `final` improve performance?"* "It can, by allowing devirtualisation and inlining of calls through the final type. But its main value is expressing design intent."

---

## Q263. [M16 | Virtual overloads with only one overridden | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P63

```cpp
#include <iostream>
using std::cout;
struct B {
    virtual void f(int)    { cout << "Bi "; }
    virtual void f(double) { cout << "Bd "; }
    virtual ~B() = default;
};
struct D : B { void f(int) override { cout << "Di "; } };
int main() {
    D d;
    B &rb = d;
    rb.f(1);
    rb.f(2.5);
    d.f(2.5);
    cout << '\n';
}
```

What does this program print?

- A. `Di Bd Di`
- B. `Di Bd Bd`
- C. `Di Di Di`
- D. `Bi Bd Di`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** Two separate steps decide every call:
1. **Which overload?** This is chosen at compile time, using the type you're calling *through*.
2. **Which class's version of that overload?** This is chosen at run time, through the vtable.

Through `B&`, both overloads are visible, so `2.5` picks `f(double)`, which D didn't override. Through `D`, only D's `f(int)` is visible (it **hides** B's `f(double)`), so 2.5 gets converted to an `int`.

**Why:**

| Call | Step 1: overload chosen (static type) | Step 2: dispatch | Prints |
|---|---|---|---|
| `rb.f(1)` | B's `f(int)` | overridden by D | `Di` |
| `rb.f(2.5)` | B's `f(double)` (exact match) | not overridden → B | `Bd` |
| `d.f(2.5)` | lookup in D finds only `f(int)` → 2.5 → 2 | D::f(int) | `Di` |

**Why the other options fail:**
- B: Forgets that D's `f` hides B's `f(double)` when calling through a `D`.
- C: Thinks dispatch picks the overload as well. Dispatch only picks *which class's version* of the already-chosen overload.
- D: Ignores the override for `rb.f(1)`.

**Concept recap:** The fix is `using B::f;` in D (Q246), which makes both overloads visible through D as well.

**Trap alert:** Compilers can warn about this ("hides overloaded virtual function"). It's a frequent source of "my call goes to the wrong function" bugs.

**Interview angle:** *"Does virtual dispatch consider argument types?"* "No. Overload resolution is static. Virtual dispatch only selects among the overrides of the one signature chosen at compile time."

---

## Q264. [M16 | Overload resolution uses the static type | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T6 · P63

```cpp
#include <iostream>
using std::cout;
struct Animal { virtual ~Animal() = default; };
struct Dog : Animal {};
void speak(Animal &) { cout << "animal "; }
void speak(Dog &)    { cout << "dog "; }
int main() {
    Dog d;
    Animal &a = d;
    speak(d);
    speak(a);
    cout << '\n';
}
```

What does this program print?

- A. `dog dog`
- B. `animal animal`
- C. `dog animal`
- D. Compile error: `speak(a)` is ambiguous

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** Choosing between **overloaded free functions** happens at **compile time**, based only on the declared type of the argument. `a` is declared `Animal&`, so `speak(Animal&)` is picked, even though the object really is a Dog. Only **virtual member functions** look at the real object at run time.

**Why:**

| Call | Static type of the argument | Chosen |
|---|---|---|
| `speak(d)` | `Dog` | `speak(Dog&)` (exact) |
| `speak(a)` | `Animal&` | `speak(Animal&)` |

**Why the other options fail:**
- A: Expects run-time selection of overloads (multiple dispatch). C++ only does single dispatch, on the object a virtual member function is called on.
- B: `speak(d)` has an exact `Dog&` overload.
- D: Only one overload accepts an `Animal&`, so there's no ambiguity.

**Concept recap:** To get run-time behaviour, put it **in the class**: `virtual void speak() const`. For dispatch on **two** objects' types (collisions, say), use the Visitor pattern or `std::variant` + `std::visit`.

**Trap alert:** The same happens with templates: `template<class T> void f(T&)` deduces the **static** type.

**Interview angle:** *"Does C++ support multiple dispatch?"* "No, only single dispatch, via virtual functions on the object. Overloads are resolved statically. Use the double-dispatch / Visitor pattern, or `std::visit`, to emulate it."

---

## Q265. [M16 | Overriding vs hiding: classification | Hard | MSQ | [C++] · C++17 | OA+Interview] · Type T7 · P63

```cpp
#include <iostream>
using std::cout;
struct B {
    virtual void f(int)      { cout << "Bf "; }
    void g()                 { cout << "Bg "; }
    virtual void h() const   { cout << "Bh "; }
    virtual ~B() = default;
};
struct D : B {
    void f(long)    { cout << "Df "; }
    void g()        { cout << "Dg "; }
    void h() const  { cout << "Dh "; }
};
```

Select **all** true statements.

- A. `D::f(long)` overrides `B::f(int)`.
- B. `D::g` hides `B::g`.
- C. `D::h` overrides `B::h`.
- D. With `B *p = new D;`, the call `p->f(1)` calls `B::f`.
- E. With `D d;`, the call `d.f(1)` calls `B::f(int)`.

**Answer:** B, C, D

**In simple words:** For each derived function, ask: **"Same name, same parameters, same `const`, and is the base version virtual?"**
- All yes → it **overrides** (the vtable slot is replaced).
- Same name but anything else differs, or the base version isn't virtual → it **hides**: it's a new function that merely blocks the base name when you call through D.

**Why each option** (verified: `p->f(1)` prints `Bf`, `p->h()` prints `Dh`, `d.f(1)` prints `Df`):
- **A ✗:** `long` ≠ `int`, so it's a different signature. It hides; it doesn't override.
- **B ✓:** `B::g` isn't virtual, so the same-named `D::g` hides it.
- **C ✓:** An exact signature match with a virtual base function.
- **D ✓:** The vtable slot for `f(int)` was never replaced.
- **E ✗:** Through D, lookup finds only `D::f(long)`, and 1 → long works → prints `Df`.

**Concept recap:**

| | Same scope? | Base function virtual? | Same signature? | Effect |
|---|---|---|---|---|
| Overloading | yes | — | no | chosen at compile time |
| Overriding | no (base / derived) | **yes** | **yes** | chosen at run time |
| Hiding | no (base / derived) | either | either | blocks base names during lookup |

**Trap alert:** Writing `override` on `D::f(long)` would have turned this silent bug into a compile error.

**Interview angle:** Draw this table. It's one of the most-asked C++ OOP interview questions.

---

## Q266. [M16 | A virtual destructor working correctly | Easy | MCQ | [C++] · C++17 | OA] · Type T7 · P60

```cpp
#include <iostream>
using std::cout;
struct B { virtual ~B() { cout << "~B"; } };
struct D : B { ~D() { cout << "~D"; } };
int main() {
    B *p = new D;
    delete p;
    cout << '\n';
}
```

What does this program print?

- A. `~B`
- B. `~B~D`
- C. `~D`
- D. `~D~B`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** When you `delete` through a base pointer, C++ needs to know the **real** type to clean up properly. Making the base destructor `virtual` puts it in the vtable, so `delete p` finds `~D` first. Then, as always, the base part is destroyed after the derived part.

**Why:**

| Step | Action | Prints |
|---|---|---|
| `delete p` | virtual destructor → dispatch to the dynamic type D | |
| run `~D` body | | `~D` |
| then the base subobject's destructor | | `~B` |
| free the memory | | |

**Why the other options fail:**
- A: What tends to happen **without** `virtual`, and that's UB (Q267).
- B: Destroys the base first.
- C: Forgets that base destructors always run after the derived one.

**Concept recap:** `~D` is automatically virtual, because `~B` is.

**Trap alert:** Writing `virtual ~B() = default;` is enough. You don't need a body.

**Interview angle:** *"When should a destructor be virtual?"* "Whenever a class is meant to be used polymorphically, i.e. it has any virtual function or will be deleted through a base pointer."

---

## Q267. [M16 | Deleting through a base with a non-virtual destructor | Medium | MCQ | [C++] · C++17 | All] · Type T3 · P60

> **This code has undefined behaviour.**

```cpp
#include <iostream>
using std::cout;
struct B { ~B() { cout << "~B"; } };
struct D : B {
    int *buf = new int[100];
    ~D() { delete[] buf; cout << "~D"; }
};
int main() {
    B *p = new D;
    delete p;
}
```

Which statement is correct?

- A. It prints `~D~B`.
- B. It prints `~B`, and leaks `buf`, but is otherwise well-defined.
- C. The behaviour is undefined: deleting a D through a `B*` when `B`'s destructor isn't virtual.
- D. Compile error: `delete` through a base pointer requires a virtual destructor.

**Answer:** C

**Behaviour class:** **Undefined behaviour.**

**In simple words:** Without `virtual`, `delete p` only knows the pointer's type (B). It runs `~B` as if the object were just a B, so D's clean-up never runs, and even freeing the memory may go wrong (the sizes don't match). The standard doesn't say "it just leaks". It says **anything can happen**.

> **⚠ UB box**
> - **Category:** `delete` through a base pointer whose destructor isn't virtual.
> - **Rule:** [expr.delete]p3: if the static type of the object being deleted differs from its dynamic type, the static type must be a base class with a **virtual destructor**, or the behaviour is undefined.
> - **What you may see in practice:** often only `~B` is printed, and `buf` leaks. With multiple inheritance or a custom `operator delete`, it can corrupt the heap or crash. **Not guaranteed.**

**Why the other options fail:**
- A: The virtual-destructor result (Q266).
- B: Describes a common symptom as if it were the defined behaviour. It isn't defined.
- D: The compiler isn't required to diagnose this (some warn about deleting a polymorphic object through a pointer with a non-virtual destructor).

**Concept recap:** Fixes:
- `virtual ~B() = default;`
- or, if the class isn't meant for polymorphic deletion, make `~B` **protected and non-virtual**, so `delete p` through a `B*` doesn't compile (Q268-D).

**Trap alert:** Inheriting from standard library containers (`class MyVec : public std::vector<int>`) and deleting through `std::vector<int>*` hits exactly this: STL containers have non-virtual destructors.

**Interview angle:** *"What happens if a base class destructor isn't virtual?"* "Deleting a derived object through a base pointer is UB: typically only the base destructor runs, leaking the derived resources. The classic interview answer, with the [expr.delete] rule."

---

## Q268. [M16 | Virtual destructors: the rules | Hard | MSQ | [C++] · C++17 | Interview] · Type T3 · P60

Assume `struct B { ~B(); }` (non-virtual) and `struct D : B {};` unless an option says otherwise. Select **all** true statements.

- A. If objects will be deleted through a base pointer, the base needs a virtual destructor.
- B. `std::shared_ptr<B> p = std::make_shared<D>();` destroys the `D` correctly even though `~B` isn't virtual.
- C. `std::unique_ptr<B> p = std::make_unique<D>();` destroys the `D` correctly even though `~B` isn't virtual.
- D. Making `~B` **protected** and non-virtual means `delete` through a `B*` fails to compile outside the class hierarchy.
- E. Declaring a virtual destructor in a class with no other virtual functions adds a vptr to each object.

**Answer:** A, B, D, E

**In simple words:**
- `shared_ptr` remembers **how to delete the original type** when it's created from `make_shared<D>()`, so it calls `~D` even through a `shared_ptr<B>`.
- `unique_ptr<B>` just calls `delete` on a `B*`, which is the Q267 bug.
- A protected destructor means "outsiders can't delete through a base pointer", which turns the bug into a compile error.
- A virtual destructor is a virtual function like any other, so it needs a vptr.

**Why each option** (B verified: prints `~ND~NB` with a non-virtual base):
- **A ✓:** The core rule (Q267).
- **B ✓:** `shared_ptr`'s control block stores a deleter for the **D** it was created from (type erasure).
- **C ✗:** `unique_ptr<B>`'s default deleter is `delete (B*)`, which is UB here.
- **D ✓:** The guideline: "a base class destructor should be either public and virtual, or protected and non-virtual."
- **E ✓:** Polymorphic types carry a vptr. That's the cost (Q211).

**Concept recap:** The rule of thumb: **if a class has any virtual function, give it a virtual destructor.**

**Trap alert:** Don't rely on option B as a design. Construct a `shared_ptr` from a raw `B*` (`std::shared_ptr<B>(raw_b_ptr)`) and the deleter is for B again.

**Interview angle:** *"Why does `shared_ptr<Base>` work without a virtual destructor but `unique_ptr<Base>` doesn't?"* "`shared_ptr` type-erases the deleter at construction, capturing `delete (D*)`. `unique_ptr`'s deleter is part of its type and calls `delete (B*)`."

---

## Q269. [M16 | Abstract classes | Easy | MSQ | [C++] · C++17 | OA] · Type T2 · P66

```cpp
struct Shape { virtual double area() const = 0; virtual ~Shape() = default; };
struct Sq : Shape { double s = 2; double area() const override { return s * s; } };
struct Half : Shape {};
```

Select **all** declarations that **compile**.

- A. `Sq q;`
- B. `Shape s;`
- C. `Half h;`
- D. `Shape *p = &q;` (after A)
- E. `Shape &r = q;` (after A)

**Answer:** A, D, E

**In simple words:** `= 0` marks a function as **"must be provided by a derived class"**. A class with any such unfinished function is **abstract**: a blueprint you can't build an object from directly. You *can* hold pointers and references to it, pointing at complete derived objects. A derived class that doesn't fill in every `= 0` is still abstract.

**Why each option:**

| Option | Type | Abstract? | Result |
|---|---|---|---|
| A | Sq | no (implements `area`) | ✓ |
| B | Shape | yes | ✗ "cannot declare variable 's' to be of abstract type" |
| C | Half | yes (inherits the pure function without implementing it) | ✗ |
| D | `Shape*` | pointers to abstract types are fine | ✓ |
| E | `Shape&` | references too | ✓ |

**Concept recap:** Abstract classes define **interfaces**. C++ has no `interface` keyword; a class with only pure virtual functions (and a virtual destructor) plays that role.

**Trap alert:** You also can't pass or return an abstract class **by value** (`void f(Shape s)` is an error). That's also slicing protection: you can't accidentally slice a Shape.

**Interview angle:** *"Abstract class vs interface in C++?"* "C++ has only classes. An 'interface' is the convention of an abstract class with only pure virtual functions and no data. Abstract classes may also have data, constructors and implemented functions."

---

## Q270. [M16 | A pure virtual function with a body | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P66

```cpp
#include <iostream>
using std::cout;
struct B { virtual void f() = 0; virtual ~B() = default; };
void B::f() { cout << "B::f "; }
struct D : B { void f() override { B::f(); cout << "D::f"; } };
int main() {
    D d;
    B &r = d;
    r.f();
    cout << '\n';
}
```

What does this program print?

- A. `B::f D::f`
- B. `D::f`
- C. Compile error: a pure virtual function can't have a body
- D. Compile error: `B::f()` is called inside `D::f`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** `= 0` means **"derived classes must override this"**, but it **doesn't** forbid the base from also providing a default implementation. Derived classes can then call that implementation explicitly (`B::f()`) as a starting point. B stays abstract either way.

**Why:**

| Step | Action | Prints |
|---|---|---|
| `r.f()` | virtual dispatch → `D::f` | |
| inside `D::f`: `B::f()` | a qualified, non-virtual call to B's body | `B::f ` |
| | rest of `D::f` | `D::f` |

**Why the other options fail:**
- B: Ignores the explicit base call.
- C: A pure virtual function **may** be defined, but outside the class body (as here). The syntax `virtual void f() = 0 { … }` inside the class isn't allowed.
- D: A qualified call to a defined pure virtual function is legal.

**Concept recap:** Uses: a "default behaviour you must opt into", and **pure virtual destructors**, which **must** have a body (Q271).

**Trap alert:** Calling `B::f()` when it's pure and has **no** definition is a **link error** (undefined reference).

**Interview angle:** *"Can a pure virtual function have an implementation?"* "Yes. The class is still abstract and derived classes must still override, but they can call `Base::f()` explicitly. It's required for pure virtual destructors."

---

## Q271. [M16 | A pure virtual destructor without a definition | Hard | MCQ | [C++] · C++17 | Interview] · Type T3 · P66, P42

```cpp
struct B { virtual ~B() = 0; };
struct D : B { ~D() {} };
int main() { D d; }
```

What happens when this program is built?

- A. It compiles and runs, doing nothing.
- B. Compile error: a destructor can't be pure virtual.
- C. **Link error:** undefined reference to `B::~B()`. Every derived destructor calls the base destructor, so it needs a body. Fix: add `B::~B() {}`.
- D. It compiles, then has undefined behaviour at run time.

**Answer:** C

**Behaviour class:** Linker error.

**In simple words:** Making the destructor pure (`= 0`) is a trick to make a class abstract when it has no other natural pure function. But unlike other functions, a destructor is **always** called automatically: when a `D` is destroyed, `~D` runs and then **`~B` must run too**. So `~B` needs a body even though it's pure. The compiler accepts the declaration, but the linker can't find the body.

**Why:**

| Stage | Result |
|---|---|
| compile | ✓ `virtual ~B() = 0;` is legal; `D` is concrete, since it has an implicit, non-pure destructor |
| link | `~D` (and the D object's destruction) calls `B::~B()` → **undefined reference** |

**Why the other options fail:**
- A: Ignores the missing definition.
- B: Pure virtual destructors are legal.
- D: Undefined symbols stop the program at link time, before it can run.

**Concept recap:** The pattern:
```cpp
struct B { virtual ~B() = 0; };
B::~B() = default;   // or: B::~B() {}
```
This makes B abstract while keeping destruction working.

**Trap alert:** This is a "compiles but doesn't link" question: classify it as a **linker error** (P42), not ill-formed.

**Interview angle:** *"Why would you declare a pure virtual destructor?"* "To make a class abstract when no other member function is naturally pure. It must still be defined, because every derived destructor implicitly calls it."

---

## Q272. [M16 | Virtual calls inside constructors and destructors | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P72

```cpp
#include <iostream>
using std::cout;
struct B {
    B()  { hello(); }
    virtual void hello() { cout << "B "; }
    virtual ~B() { hello(); }
};
struct D : B {
    D()  { hello(); }
    void hello() override { cout << "D "; }
    ~D() { hello(); }
};
int main() { { D d; } cout << '\n'; }
```

What does this program print?

- A. `D D D D`
- B. `B D D B`
- C. `B B B B`
- D. `D D D B`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** While `B`'s constructor runs, the D part of the object **doesn't exist yet**. So C++ treats the object as **just a B** during that time, and virtual calls go to B's version. It's like a house still under construction: you can't use rooms that haven't been built. The same applies in reverse during destruction: by the time `~B` runs, the D part is already gone.

**Why:**

| Phase | Object currently treated as | `hello()` calls | Prints |
|---|---|---|---|
| `B()` body | B (the vptr points to B's vtable) | B::hello | `B ` |
| `D()` body | D | D::hello | `D ` |
| `~D()` body | D | D::hello | `D ` |
| `~B()` body | **B** (D already destroyed; vptr reset to B's) | B::hello | `B ` |

**Why the other options fail:**
- A: Assumes full dispatch even during construction (Java and C# behave this way; C++ doesn't).
- C: Assumes no dispatch at all.
- D: Gets construction wrong.

**Concept recap:** The rule ([class.cdtor]p4): during the constructor or destructor of class X, a virtual call resolves as if the object's dynamic type were X.

**Trap alert:** If the function were **pure** in B, calling it from B's constructor or destructor (directly, or indirectly as in Q274) would be **UB**.

**Interview angle:** *"Can you call virtual functions in a constructor?"* "You can, but they won't dispatch to derived overrides: during B's constructor, the object is a B. Java does dispatch, and can then call a method on uninitialised fields. C++ avoids that hazard. Use a factory or a two-phase `init()` if you need derived behaviour."

---

## Q273. [M16 | Why dispatch stops at the class being built | Hard | MCQ | [C++] · C++17 | Interview] · Type T7 · P72

```cpp
#include <iostream>
#include <string>
using std::cout;
struct B {
    B() { init(); }
    void init() { cout << name() << ' '; }
    virtual const char *name() const { return "base"; }
    virtual ~B() = default;
};
struct D : B {
    std::string s = "derived";
    const char *name() const override { return s.c_str(); }
};
int main() {
    D d;
    d.init();
    cout << '\n';
}
```

What does this program print?

- A. `derived derived`
- B. Undefined behaviour: `init()` reads `s` before it is constructed
- C. `base base`
- D. `base derived`

**Answer:** D

**Behaviour class:** Well-defined, precisely **because** of the constructor-dispatch rule.

**In simple words:** If C++ dispatched to `D::name()` during B's constructor, it would read the string `s`, which **hasn't been built yet** (members of D are built after B finishes). That's exactly the danger the rule prevents. During B's constructor, `name()` goes to B's safe version. Once the object is fully built, `d.init()` reaches D's version.

**Why:**

| Moment | D's member `s` | `name()` resolves to | Prints |
|---|---|---|---|
| inside `B()` (called from D's construction) | **not yet constructed** | B::name (dynamic type is B for now) | `base` |
| after construction, `d.init()` | constructed | D::name | `derived` |

**Why the other options fail:**
- A: What Java-style dispatch would *attempt*, reading an unbuilt string.
- B: This would be UB in a language that dispatched to D here. C++'s rule makes it safe.
- C: Ignores the normal dispatch after construction.

**Concept recap:** It holds even when the virtual call is **indirect** (through the non-virtual `init()`): what matters is the object's state at that moment, not the call syntax.

**Trap alert:** Passing `this` out of a constructor (to a callback or registry) and having someone call virtual functions on it *before construction finishes* gives the same "base version" behaviour. With a pure function, it's UB (Q274).

**Interview angle:** *"Why does C++ set the vptr to the base vtable during the base constructor?"* "Because derived members aren't initialised yet. Dispatching to derived overrides would use uninitialised state. C++ trades 'surprising' for 'safe'."

---

## Q274. [M16 | A pure virtual call during construction | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T3 · P72

> **This code has undefined behaviour.**

```cpp
#include <iostream>
struct B {
    B() { call(); }
    void call() { f(); }
    virtual void f() = 0;
    virtual ~B() = default;
};
struct D : B { void f() override { std::cout << "D"; } };
int main() { D d; }
```

Which statement is correct?

- A. It prints `D`.
- B. It prints nothing; the call is skipped.
- C. The behaviour is undefined: during B's constructor the object is treated as a B, so `f()` resolves to the **pure** `B::f`, and a virtual call to a pure virtual function during construction is UB.
- D. Compile error: a pure virtual function can't be called from a constructor.

**Answer:** C

**Behaviour class:** **Undefined behaviour.**

**In simple words:** By the Q272 rule, inside `B()` the object counts as "just a B", and B's `f` has **no implementation** (it's pure). There's nothing sensible to run, so the standard calls it undefined. Most systems detect this and **abort** the program with a "pure virtual function called" message.

> **⚠ UB box**
> - **Category:** a virtual call to a pure virtual function during construction or destruction.
> - **Rule:** [class.abstract]p6: making a virtual call to a pure virtual function from a constructor or destructor of an abstract class is undefined behaviour.
> - **What you may see in practice:** a run-time message such as "pure virtual method called", followed by an abort (the vtable slot points to an error handler). **Not guaranteed.**

**Why the other options fail:**
- A: Assumes dispatch to D, which doesn't happen during B's constructor.
- B: Nothing guarantees that the call is skipped.
- D: Compilers may warn when the pure call is **direct** inside the constructor, but here it's hidden behind `call()`. A diagnostic isn't required in either case.

**Concept recap:** Never call virtual functions (directly or indirectly) from constructors and destructors when you need derived behaviour. Use factory functions (`create()` that constructs, then calls `init()`), or pass the needed data to the base constructor.

**Trap alert:** The same UB happens in **destructors**, when a base destructor calls a pure function after the derived part has already been destroyed.

**Interview angle:** *"What does 'pure virtual function called' mean?"* "A virtual call reached a pure function during construction or destruction (or through a dangling pointer to a destroyed object). The C++ runtime aborts. Look for virtual calls in constructors and destructors."

---

## Q275. [M16 | Covariant return types | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
using std::cout;
struct B {
    virtual B *clone() const { cout << "B"; return new B(*this); }
    virtual ~B() = default;
};
struct D : B {
    D *clone() const override { cout << "D"; return new D(*this); }
};
int main() {
    D d;
    B &r = d;
    B *c1 = r.clone();
    D *c2 = d.clone();
    cout << (dynamic_cast<D *>(c1) != nullptr) << '\n';
    delete c1; delete c2;
}
```

What does this program print?

- A. `DD1`
- B. `BD0`
- C. Compile error: the override's return type differs from the base's
- D. `DD0`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** An override normally must return **exactly** the same type as the base. There's one allowed exception: if the base returns `B*` (or `B&`), the override may return a pointer (or reference) to a **more derived** class, `D*`. That's safe, because a `D*` *is* a `B*`. Callers through a `B` get a `B*`; callers through a `D` get a `D*` directly, with no cast needed.

**Why:**

| Call | Dispatch | Returned static type | Prints |
|---|---|---|---|
| `r.clone()` | virtual → D::clone | `B*` (pointing to a new D) | `D` |
| `d.clone()` | D::clone | **`D*`**, no cast needed | `D` |
| `dynamic_cast<D*>(c1)` | c1 really points to a D | non-null | `1` |

**Why the other options fail:**
- B: Ignores dispatch.
- C: **Covariant return types** are allowed ([class.virtual]p8), for pointers and references to classes only.
- D: c1 does point to a D.

**Concept recap:** `virtual clone()` is the standard **"virtual copy constructor"** idiom: it copies an object through a base pointer without knowing its real type. Modern versions return `std::unique_ptr<B>`, which **isn't** covariant (smart pointers don't qualify), so derived overrides must also return `unique_ptr<B>`.

**Trap alert:** `int f()` in the base and `long f()` in the derived class isn't covariant: it's a compile error ("conflicting return type"), not a silently hidden function.

**Interview angle:** *"How do you copy an object through a base pointer?"* "A virtual `clone()` that each class implements as `return new Derived(*this);` (or `std::make_unique`). Covariant returns let derived callers avoid a cast."

---

## Q276. [M16 | Default arguments on virtual functions | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
using std::cout;
struct B { virtual void show(int x = 1) { cout << "B" << x << ' '; } virtual ~B() = default; };
struct D : B { void show(int x = 2) override { cout << "D" << x << ' '; } };
int main() {
    D d;
    B &r = d;
    r.show();
    d.show();
    cout << '\n';
}
```

What does this program print?

- A. `D2 D2`
- B. `D1 D2`
- C. `B1 D2`
- D. `D1 D1`

**Answer:** B

**Behaviour class:** Well-defined (and a well-known trap).

**In simple words:** Two things are decided at **different times**:
- **Which function runs** is decided at **run time** (the real object → D).
- **Which default argument is filled in** is decided at **compile time**, by the **declared type you call through** (`B&` → B's default, 1).

So you get D's code running with B's default value: a mix-and-match.

**Why:**

| Call | Default argument taken from (static type) | Function run (dynamic type) | Prints |
|---|---|---|---|
| `r.show()` | `B::show` → `x = 1` | `D::show` | `D1` |
| `d.show()` | `D::show` → `x = 2` | `D::show` | `D2` |

**Why the other options fail:**
- A: Assumes the default follows the dynamic type.
- C: Assumes no dispatch.
- D: Assumes the default comes only from the base.

**Concept recap:** Rule ([dcl.fct.default]p10): default arguments are bound using the static type of the pointer or reference; the override's own defaults are ignored for calls through the base. **Guideline: never redefine default arguments in overrides.** Use the same defaults, or none (with a non-virtual wrapper that supplies the default).

**Trap alert:** This shows up in OA "predict the output" questions almost word for word. Remember: "virtual function, static default".

**Interview angle:** *"What's wrong with different default arguments in an override?"* "Defaults are resolved at compile time from the static type, so calls through a base reference use the base's defaults with the derived implementation. It's confusing and error-prone."

---

## Q277. [M16 | Object slicing | Easy | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P61

```cpp
#include <iostream>
using std::cout;
struct B {
    int b = 1;
    virtual const char *who() const { return "B"; }
    virtual ~B() = default;
};
struct D : B {
    int d = 2;
    const char *who() const override { return "D"; }
};
void byVal(B x)        { cout << x.who() << ' '; }
void byRef(const B &x) { cout << x.who() << ' '; }
int main() {
    D d;
    byVal(d);
    byRef(d);
    B copy = d;
    cout << copy.who() << ' ' << sizeof(copy) << '\n';
}
```

What does this program print?

- A. `D D D 24`
- B. `B D D 16`
- C. `B B B 16`
- D. `B D B 16`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** Copying a derived object **into a base-class variable** keeps only the base part and **slices off** the rest, like photocopying only the first page of a document. The copy really *is* a plain B (with B's vtable and no `d`), so virtual calls on it give B's behaviour. References and pointers don't copy anything, so they keep the full D.

**Why:**

| Use | Copy made? | Object's real type | `who()` |
|---|---|---|---|
| `byVal(d)` | yes → a sliced B | B | `B` |
| `byRef(d)` | no | D | `D` |
| `B copy = d` | yes → sliced | B | `B` |
| `sizeof(copy)` | vptr 8 + int 4 + padding 4 | | 16 |

**Why the other options fail:**
- A: Assumes copies keep the derived type (that would need a heap clone, Q275).
- B: Thinks the named `copy` keeps D's behaviour.
- C: Thinks references slice too.

**Concept recap:** Slicing loses the derived data (`d`) **and** the derived behaviour: the B copy constructor builds a genuine B, with B's vptr.

**Trap alert:** `std::vector<B>` storing D objects slices every element. Use `std::vector<std::unique_ptr<B>>` for polymorphic collections.

**Interview angle:** *"What is object slicing and how do you avoid it?"* "Copying a derived object into a base object by value, which loses the derived part. Avoid it by passing by reference or pointer, storing smart pointers in containers, making base classes abstract, or deleting the base copy operations."

---

## Q278. [M16 | Slicing through a base reference, and in containers | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P61

```cpp
#include <iostream>
#include <vector>
using std::cout;
struct B {
    int b = 1;
    virtual const char *who() const { return "B"; }
    virtual ~B() = default;
};
struct D : B {
    int d = 2;
    const char *who() const override { return "D"; }
};
int main() {
    D d1, d2;
    d1.d = 10; d2.d = 20; d2.b = 5;
    B &r = d1;
    r = d2;
    std::vector<B> v;
    v.push_back(d1);
    cout << d1.b << ' ' << d1.d << ' ' << v[0].who() << '\n';
}
```

What does this program print?

- A. `5 10 B`
- B. `5 20 D`
- C. `1 10 D`
- D. `5 20 B`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:**
1. `r = d2` uses **B's** assignment operator, because `r` is a `B&` and assignment isn't virtual. So only the B part of `d1` is overwritten (`b` becomes 5), while `d1.d` keeps 10: half of one object, half of the other. This is "partial assignment".
2. Pushing `d1` into a `vector<B>` stores a **sliced copy**: a plain B.

**Why:**

| Step | d1.b | d1.d | Note |
|---|---|---|---|
| start | 1 | 10 | |
| `r = d2` → `B::operator=` copies only the B part | **5** | 10 | the D part is untouched |
| `v.push_back(d1)` | | | the vector stores a B copy (sliced) |
| `v[0].who()` | | | `B` |

Prints `5 10 B`.

**Why the other options fail:**
- B, D: Assume the full D was assigned.
- C: Assumes the base-reference assignment did nothing.

**Concept recap:** **Partial assignment** is the nastier relative of slicing: it silently corrupts an existing derived object into a mix of two objects.

**Trap alert:** Making `operator=` virtual doesn't really help, because the parameter is still `const B&`. The standard fixes are abstract bases, `= delete` on the base assignment, or the `clone()` pattern.

**Interview angle:** *"Why shouldn't you assign derived objects through base references?"* "Assignment isn't polymorphic: only the base part is copied, which leaves a corrupt mixed object. Prefer non-copyable polymorphic bases (`=delete` copy operations) and clone via virtual functions."

---

## Q279. [M16 | dynamic_cast as a cross-cast | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
using std::cout;
struct Printable { virtual void print() const = 0; virtual ~Printable() = default; };
struct Saveable  { virtual void save()  const = 0; virtual ~Saveable()  = default; };
struct Doc : Printable, Saveable {
    void print() const override { cout << "P"; }
    void save()  const override { cout << "S"; }
};
struct Img : Printable { void print() const override { cout << "I"; } };
int main() {
    Printable *a = new Doc, *b = new Img;
    if (auto s = dynamic_cast<Saveable *>(a)) s->save();
    if (auto s = dynamic_cast<Saveable *>(b)) s->save(); else cout << "-";
    cout << '\n';
    delete a; delete b;
}
```

What does this program print?

- A. `SS`
- B. `S-`
- C. `--`
- D. Compile error: `Printable` and `Saveable` are unrelated

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:** `dynamic_cast` can ask an object, at run time, **"do you *also* support this other interface?"**, even when the two interfaces aren't related to each other. A `Doc` is both Printable and Saveable, so the cast succeeds. An `Img` isn't Saveable, so the cast returns `nullptr`.

**Why:**

| Pointer | Real type | Is it a Saveable? | Result |
|---|---|---|---|
| a | Doc | yes | non-null → `save()` → `S` |
| b | Img | no | `nullptr` → prints `-` |

The cast from `Printable*` to `Saveable*` is a **cross-cast**. It works by inspecting the complete object's type information (RTTI), then adjusting the pointer to the Saveable subobject (Q254).

**Why the other options fail:**
- A: Img doesn't derive from Saveable.
- C: Doc does.
- D: `dynamic_cast` allows casts between unrelated polymorphic types and checks them at run time. (`static_cast` would refuse, Q280-E.)

**Concept recap:** Capability queries like "if it can be saved, save it" are the main legitimate use of `dynamic_cast`.

**Trap alert:** Frequent `dynamic_cast` chains (`if (dynamic_cast<A*>) … else if (dynamic_cast<B*>) …`) usually signal a missing virtual function in the base.

**Interview angle:** *"What is a cross-cast?"* "Casting between two base classes of the same object, which only `dynamic_cast` can do. It needs RTTI and polymorphic types, and returns `nullptr` if the object doesn't have the target base."

---

## Q280. [M16 | The rules of dynamic_cast | Hard | MSQ | [C++] · C++17 | Interview] · Type T3

Select **all** true statements.

- A. `dynamic_cast` to a **pointer** type returns `nullptr` on failure.
- B. `dynamic_cast` to a **reference** type throws `std::bad_cast` on failure.
- C. `dynamic_cast<void *>(p)` (with `p` pointing to a polymorphic object) returns a pointer to the **most-derived** (complete) object.
- D. `dynamic_cast` can downcast from a **non-polymorphic** base, because the check is done at compile time.
- E. `static_cast<Saveable *>(a)` would also perform the cross-cast in Q279.

**Answer:** A, B, C

**In simple words:** A failed `dynamic_cast` has to signal failure somehow: a pointer can be null, but a reference can't, so a reference cast throws instead. Casting to `void*` means "give me the address of the whole object", whatever base you're starting from. And `dynamic_cast` can only check types for classes that carry run-time type information, i.e. that have virtual functions.

**Why each option:**
- **A ✓:** The standard pointer-cast failure signal.
- **B ✓:** There's no "null reference", so it throws `std::bad_cast`.
- **C ✓:** Verified with Q279's `Doc`: `dynamic_cast<void*>` applied to a `Printable*` gives `&doc`. That's useful for checking whether two base pointers refer to the same object.
- **D ✗:** Downcasts and cross-casts require a polymorphic source type (Q253-D). The check happens at **run time**.
- **E ✗:** Printable and Saveable are unrelated classes, so `static_cast` between them is ill-formed. Only `dynamic_cast` (a run-time check) or `reinterpret_cast` (wrong and dangerous here) would compile.

**Concept recap:**

| Cast | Run-time check | Allowed directions |
|---|---|---|
| `static_cast` | none | up, down (unchecked), not across |
| `dynamic_cast` | yes (RTTI) | up, down, **across**, and to `void*` |

**Trap alert:** Some codebases disable RTTI to save space. Then `dynamic_cast` and `typeid` on polymorphic types aren't available. That's a non-standard build mode, but common in game and embedded development.

**Interview angle:** *"How does dynamic_cast work, and what does it cost?"* "It reads the type information through the object's vptr, walks the class hierarchy to find the target subobject, and adjusts the pointer. It costs time (proportional to the hierarchy's complexity) and needs RTTI. Use it rarely; prefer virtual functions."

---

## Batch 14 summary

| Q | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|
| 261 | E | MCQ | B | Ill-formed (L2, L3) |
| 262 | M | MSQ | A, B, C | — (rules) |
| 263 | H | MCQ | A | Well-defined |
| 264 | M | MCQ | C | Well-defined |
| 265 | H | MSQ | B, C, D | Well-defined |
| 266 | E | MCQ | D | Well-defined |
| 267 | M | MCQ | C | **UB** |
| 268 | H | MSQ | A, B, D, E | UB / OK rules |
| 269 | E | MSQ | A, D, E | Ill-formed / OK mix |
| 270 | M | MCQ | A | Well-defined |
| 271 | H | MCQ | C | **Linker error** |
| 272 | M | MCQ | B | Well-defined |
| 273 | H | MCQ | D | Well-defined |
| 274 | H | MCQ | C | **UB** |
| 275 | M | MCQ | A | Well-defined |
| 276 | H | MCQ | B | Well-defined |
| 277 | E | MCQ | D | Well-defined |
| 278 | H | MCQ | A | Well-defined |
| 279 | M | MCQ | B | Well-defined |
| 280 | H | MSQ | A, B, C | Ill-formed / OK rules |

- **Difficulty:** 4 E / 7 M / 9 H.
  - M16 total (Q255–Q280): 5/10/11, as planned.
- **Format:** 15 MCQ / 5 MSQ.
  - M16 total: 20/6, as planned.
- **MCQ answer keys:** A×4, B×4, C×4, D×3. Running total Q1–Q280: A 56, B 56, C 56, D 56.

---

**Completed: 280 / 400 Remaining: 120**

**Ledger, newly ✓ this batch:**
- P60 ×3 (Q266, Q267, Q268) ★ complete
- P61 ×2 (Q277, Q278) ★ complete (+1 more planned in M25)
- P66 ×3 (Q269, Q270, Q271)
- P72 ×3 (Q272, Q273, Q274) ★ complete
- P63 → 8 (Q263, Q264, Q265)
- P42 → 6 (Q271)

**M16 complete** (Q255–Q280, 26 questions).

*Say **"Continue."** for Batch 15 (Q281–Q300: M17 overloading — overload-resolution ranking, ambiguity and default arguments, const overloads, assignment / arithmetic / comparison / stream operators, conversion operators — then M18 copy/move begins with value categories and choosing between the copy and move constructors). Batch 15 ends with **Self-audit #3**, the full **Pattern Ledger**, and **Progress report @300**.*
