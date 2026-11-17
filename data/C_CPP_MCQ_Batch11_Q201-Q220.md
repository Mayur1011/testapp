# C/C++ MCQ/MSQ Bank — Batch 11 (Q201–Q220)

**What this batch covers:**
- **M13 Classes (Q201–Q214):**
  - class vs struct
  - `this`: returning `*this` by reference vs by value, and the type of `this` inside a `const` member
  - static data members and object counters
  - the rules for static member functions
  - a static member hidden in a derived class
  - `const` + `mutable`
  - overloading a member on `const`
  - friend functions and friend classes (friendship isn't inherited), and the rules of friendship
  - object size: empty classes, the empty-base optimisation, the vptr, and what *doesn't* count toward size
  - nested classes
  - access-control errors
- **M14 Constructors, part 1 (Q215–Q220):**
  - which constructor or assignment runs
  - a missing default constructor
  - the most vexing parse
  - copies made when passing by value, and C++17 guaranteed elision
  - `explicit` constructors
  - `explicit operator bool`

**Patterns:** P33, P52, P53, P63 (secondary), P67, P68, P71 (secondary)

> **Rules in effect:**
> - Questions are compiler-independent: answers follow from the standard plus the stated platform.
> - Every question has an **"In simple words"** explanation.
> - Object traces show which constructor, destructor or operator runs, and why.

> **Conventions for every question in this batch:**
> - **C++17**, exceptions enabled. C++17 **guaranteed copy elision** applies (initialising an object from a temporary of the same type creates the object directly).
> - Platform: a typical 64-bit system: `int` 4, `double` 8, pointers 8 bytes, natural alignment.
> - **Object layout (typical ABIs):** a class with virtual functions holds one hidden **vptr** (8 bytes) at the start, and empty base classes take no space (the empty-base optimisation). The standard leaves layout to the implementation, but these are universal on mainstream 64-bit compilers.
> - Needed headers are included, and `std::` is written explicitly.

---

## Q201. [M13 | class vs struct | Easy | MCQ | [C++] · C++17 | OA] · Type — concept

What are the **only** language-level differences between `class` and `struct` in C++?

- A. The default member access (`private` for `class`, `public` for `struct`) and the default base-class access when inheriting (also `private` vs `public`).
- B. A `struct` can't have member functions.
- C. A `struct` can't have constructors, virtual functions or private members.
- D. `class` objects are created on the heap, and `struct` objects on the stack.

**Answer:** A

**In simple words:** In C++, `struct` and `class` are **the same tool with different default settings**. A `struct` starts out "open to everyone" (public), and a `class` starts out "private". Anything one can do, the other can do too.

**Why:**

```cpp
struct S { int x; };          // x is public
class  C { int x; };          // x is private
struct D1 : S {};             // public inheritance by default
class  D2 : S {};             // private inheritance by default
```

**Why the other options fail:**
- B, C: Describe C structs. In C++, a `struct` can have member functions, constructors, virtual functions, private members, templates, and so on.
- D: Where an object lives depends on **how it's created** (`T x;` vs `new T`), never on the keyword.

**Concept recap:** The common convention: `struct` for passive data aggregates with public fields (like a C struct), and `class` for types that maintain invariants behind private data and a public interface.

**Trap alert:** Forgetting `public` in `class Derived : Base` gives **private** inheritance. Then `Base *p = &derived;` fails to compile outside `Derived` (M15).

**Interview angle:** *"struct vs class in C++?"* "Only the default access, for members and for base classes. Use `struct` for plain data and `class` when you're enforcing invariants." Contrast with C#, where struct means a value type: that's not true in C++.

---

## Q202. [M13 | Returning `*this` by reference vs by value | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
struct Acc {
    int v = 0;
    Acc &add(int x)     { v += x; return *this; }   // returns a reference
    Acc  copy_add(int x){ v += x; return *this; }   // returns a copy
};
int main() {
    Acc a;
    a.add(1).add(2);
    a.copy_add(10).add(100);
    std::cout << a.v << '\n';
}
```

What does this program print?

- A. `113`
- B. `3`
- C. `13`
- D. `103`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** `this` is a hidden pointer to "the object I was called on", so `*this` *is* that object. Returning `Acc&` hands back **the same object**, so the next call keeps changing `a`. Returning `Acc` (by value) hands back a **copy**, so any later call in the chain changes only the copy.

**Why:**

| Call | Acts on | a.v | Copy's v |
|---|---|---|---|
| `a.add(1)` | a | 1 | |
| `.add(2)` | a (reference returned) | 3 | |
| `a.copy_add(10)` | a | **13** | returns a copy with v = 13 |
| `.add(100)` | **the copy** | 13 | 113, then the copy is destroyed |

Prints `13`.

**Why the other options fail:**
- A: Thinks `copy_add` returns a reference, so `add(100)` would hit `a`.
- B: Forgets that `copy_add` itself modifies `a` before copying.
- D: Mixes the two.

**Concept recap:** Return `*this` **by reference** (`T&`) for chaining ("fluent" interfaces such as `std::cout << a << b`) and from assignment operators (`T& operator=(...)`, M17).

**Trap alert:** Returning by value in an assignment operator compiles, but `(a = b) = c` then assigns to a temporary copy, so `c` never reaches `a`.

**Interview angle:** *"What is the `this` pointer?"* "An implicit parameter of every non-static member function, pointing to the object the function was called on. It's `X*` normally and `const X*` in `const` member functions (Q203). Common uses: chaining (`return *this`), self-assignment checks (`if (this == &other)`), and passing the object to other functions."

---

## Q203. [M13 | The type of `this` in a const member function | Easy | MCQ | [C++] · C++17 | OA+Interview] · Type — concept · P52

Inside `int get() const { … }` of class `X`, what is the type of `this`?

- A. `X *`
- B. `const X *`
- C. `X *const`
- D. `const X &`

**Answer:** B

**In simple words:** The `const` after a member function's parameter list is really a `const` on the hidden `this` parameter. It turns "pointer to X" into "pointer to a **read-only** X". That's why the function can't change any member.

**Why:**

| Member function | Type of `this` |
|---|---|
| `void f()` | `X *` |
| `void f() const` | **`const X *`** |
| `void f() volatile` | `volatile X *` |

(`this` itself is a prvalue, a plain value that can't be assigned to, so the question of `X *const` doesn't arise. You can never write `this = …`.)

**Why the other options fail:**
- A: The type in a non-const member function.
- C: Describes a pointer that can't be reseated, which is irrelevant here. What matters is the read-only object.
- D: `this` is a pointer, not a reference (a historical choice: references didn't exist yet when `this` was introduced).

**Concept recap:** Since `this` is `const X *`, `this->n = 5` inside a const member is a compile error, and so is calling a non-const member on it. `mutable` members are the exception (Q207).

**Trap alert:** A `const` member function can still modify objects it reaches **through pointer members**. `const` is shallow: `int *p` becomes `int *const p`, not `const int *p`. This is sometimes called the "const-pointer loophole".

**Interview angle:** *"Why can a const member function change `*ptr_member`?"* "Constness applies to the object's own members. For a pointer member, that makes the pointer constant, not the data it points to. Use `std::experimental::propagate_const`, or design accessors carefully."

---

## Q204. [M13 | Static data member as an object counter | Easy | MCQ | [C++] · C++17 | OA] · Type T7 · P68

```cpp
#include <iostream>
struct Obj {
    static int count;
    int id;
    Obj() : id(++count) {}
};
int Obj::count = 0;
int main() {
    Obj a, b;
    { Obj c; }
    Obj d;
    std::cout << Obj::count << ' ' << d.id << ' ' << sizeof(Obj) << '\n';
}
```

What does this program print?

- A. `3 3 4`
- B. `4 4 8`
- C. `3 4 4`
- D. `4 4 4`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** A `static` data member is **one shared variable for the whole class**, like a single tally counter on the classroom wall rather than one per student. Every object's constructor bumps the same counter. And since it isn't stored inside any object, it doesn't add to `sizeof`.

**Why:**

| Object created | count after | its id |
|---|---|---|
| a | 1 | 1 |
| b | 2 | 2 |
| c (destroyed at its `}`; nothing decrements) | 3 | 3 |
| d | 4 | 4 |

`sizeof(Obj)` = just `id` = 4. Prints `4 4 4`.

**Why the other options fail:**
- A: Assumes destroying `c` reduces the count. There's no destructor here that decrements it.
- B: Counts the static member as part of each object.
- C: Gets `d.id` right but the count wrong.

**Concept recap:**
- Declare the static member in the class, and define it once in a `.cpp` file (`int Obj::count = 0;`).
- Since C++17, `static inline int count = 0;` can be defined right inside the class, which is convenient in headers.

**Trap alert:** Forgetting the out-of-class definition (before C++17 `inline`) gives a **linker** error ("undefined reference to `Obj::count`"), not a compile error. This is the same idea as P42.

**Interview angle:** *"Count live objects of a class."* Increment in **every** constructor (including the copy and move constructors!) and decrement in the destructor. Forgetting the copy constructor is the classic bug: copies aren't counted, but they are decremented when destroyed.

---

## Q205. [M13 | Rules for static member functions | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P68

Select **all** true statements about static member functions.

- A. A static member function has no `this` pointer.
- B. A static member function can directly access the non-static data members of its class.
- C. A static member function can be called as `X::f()` without any object.
- D. A static member function can be declared `const`.
- E. A static member function can be `virtual`.

**Answer:** A, C

**In simple words:** A static member function belongs to **the class itself, not to any particular object**. It's like a function in the class's "namespace" that is also allowed to see private members. With no object, there's no `this`, so it can't read any particular object's fields unless you pass it an object.

**Why each option:**
- **A ✓:** There's no object, so there's no `this`.
- **B ✗:** It needs an object to say *whose* member: `obj.n` is fine if an object is passed in, but a bare `n` is an error.
- **C ✓:** `X::create()` is the classic factory-function form.
- **D ✗:** `const` on a member function qualifies `this`, which doesn't exist here, so it's ill-formed ("static member function cannot have cv-qualifier").
- **E ✗:** Virtual dispatch happens through an object's vptr. With no object, nothing can be dispatched.

**Concept recap:** Typical uses: factory functions, counters (`Obj::count()`), utility functions closely tied to a class, and callbacks for C APIs (a static member has no hidden `this` parameter, so it fits a plain C function pointer).

**Trap alert:** A static member function **can** access private members of its class, including those of an object passed to it as a parameter. Access control is per class, not per object.

**Interview angle:** *"Can a static member function be virtual or const?"* "Neither. Both depend on `this`, which a static function doesn't have."

---

## Q206. [M13 | Static members with inheritance and hiding | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P68, P63

```cpp
#include <iostream>
struct Base { static int n; Base() { ++n; } };
int Base::n = 0;
struct D1 : Base {};
struct D2 : Base { static int n; D2() { ++n; } };
int D2::n = 100;
int main() {
    D1 x;
    D2 y, z;
    std::cout << Base::n << ' ' << D1::n << ' ' << D2::n << '\n';
}
```

What does this program print?

- A. `3 3 102`
- B. `1 1 102`
- C. `3 1 102`
- D. `3 3 103`

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** A derived class **shares** its base's static variable; it doesn't get its own copy. So `D1::n` is simply another name for `Base::n`. `D2` declares a **new** variable also called `n`, which hides the base's. Inside `D2`'s constructor, a bare `n` means `D2`'s own `n`. And every object of a derived class runs the **base** constructor first.

**Why:**

| Object | Base() runs → Base::n | D2() runs → D2::n |
|---|---|---|
| x (D1) | 1 | — |
| y (D2) | 2 | 101 |
| z (D2) | 3 | 102 |

| Expression | Refers to | Value |
|---|---|---|
| `Base::n` | Base's static | 3 |
| `D1::n` | **the same** `Base::n` (inherited name) | 3 |
| `D2::n` | D2's own static (it hides Base::n) | 102 |

**Why the other options fail:**
- B: Thinks each derived class gets its own copy of `Base::n`.
- C: Thinks `D1::n` is separate from `Base::n`.
- D: Thinks D2's constructor runs for `x` as well.

**Concept recap:** Name lookup checks the derived class first. A member with the same name in a derived class **hides** the base member, for data and functions alike (pattern P63, M15–M16). `Base::n` remains reachable by qualifying it.

**Trap alert:** "Per-class counters" via a static in a base class count **all** derived objects together. To count each derived type separately, use CRTP: `template<class T> struct Counter { static inline int n = 0; };` gives one `n` per `T` (M20).

**Interview angle:** *"Is a static member inherited?"* "The **name** is accessible in derived classes, but there's only one object, shared with the base. A derived class can hide it with its own static member of the same name."

---

## Q207. [M13 | const member functions + mutable | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P52

```cpp
#include <iostream>
class Cache {
    int value = 10;
    mutable int hits = 0;
public:
    int get() const { ++hits; return value; }
    int hitCount() const { return hits; }
};
int main() {
    const Cache c{};
    c.get(); c.get();
    int g = c.get();
    std::cout << g << ' ' << c.hitCount() << '\n';
}
```

What does this program print?

- A. `10 3`
- B. `10 0`
- C. Compile error: `++hits` inside a `const` member function
- D. Compile error: calling `get()` on a `const` object

**Answer:** A

**Behaviour class:** Well-defined.

**In simple words:** `const` promises that the object's **meaningful state** won't change. Some members are just behind-the-scenes bookkeeping (a usage counter, a cache, a lock). `mutable` marks them as "allowed to change even inside a `const` function", because changing them doesn't change what the object *means*.

**Why:**

| Call | Allowed? | hits after |
|---|---|---|
| `c.get()` on a const object | ✓ `get` is const | 1 |
| `c.get()` | ✓ | 2 |
| `g = c.get()` | ✓ | 3 |
| `c.hitCount()` | ✓ | 3 |

`++hits` is legal because `hits` is `mutable`. `value` stays read-only. Prints `10 3`.

**Why the other options fail:**
- B: Thinks `mutable` changes are discarded.
- C: That would be correct **without** `mutable`.
- D: `get` is a `const` member, so const objects may call it.

**Concept recap:** Legitimate uses of `mutable`: caches and memoisation, access counters, and `std::mutex` members (locking changes the mutex but not the logical state). Lambdas also use the keyword: `[x]() mutable { ++x; }` lets a lambda modify its by-value copy of a captured variable.

**Trap alert:** `const Cache c;` (with no braces) is only allowed because every data member has a default initializer. For a `const` object, the class must be "const-default-constructible". Writing `{}` avoids the question.

**Interview angle:** *"Logical vs physical constness?"* "Physical (bitwise) constness: no bits change. Logical constness: the observable state doesn't change. C++ `const` plus `mutable` lets you express logical constness, such as a cached computation behind a const getter."

---

## Q208. [M13 | Overloading on const | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P52

```cpp
#include <iostream>
struct Arr {
    int a[3] = {1, 2, 3};
    int &at(int i)             { std::cout << "N"; return a[i]; }
    const int &at(int i) const { std::cout << "C"; return a[i]; }
};
int main() {
    Arr x;
    const Arr y{};
    x.at(0) = 10;
    int v = y.at(1);
    const Arr &rx = x;
    int w = rx.at(0);
    std::cout << ' ' << v + w << '\n';
}
```

What does this program print?

- A. `NCN 12`
- B. `NNC 3`
- C. `NCC 3`
- D. `NCC 12`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** A class can provide **two versions** of the same function: one for changeable objects and one for read-only objects. C++ picks one based on **how you're holding the object**: through a non-const name you get the "N" version, which lets you modify; through a `const` name or reference you get the "C" version, which is read-only. Holding the *same* object `x` through a `const` reference picks the const version.

**Why:**

| Call | Expression type | Chosen | Prints | Effect / value |
|---|---|---|---|---|
| `x.at(0) = 10` | `Arr` (non-const) | non-const `at` | `N` | x.a[0] = 10 |
| `y.at(1)` | `const Arr` | const `at` | `C` | v = 2 |
| `rx.at(0)` | `const Arr&` (refers to x) | const `at` | `C` | w = x.a[0] = **10** |

`v + w` = 2 + 10 = 12. Prints `NCC 12`.

**Why the other options fail:**
- A: Assumes the version is chosen by the underlying object (`x` is non-const), not by the expression used to reach it.
- B, C: Forget that `x.at(0) = 10` changed `x.a[0]` before `rx.at(0)` reads it.

**Concept recap:** That's exactly how `std::vector::operator[]` and `at()` work. The const version returns `const T&`, so you can read through a `const vector&` but not write.

**Trap alert:** If only the **non-const** overload existed, `y.at(1)` and `rx.at(0)` would fail to compile. Always provide a const overload for read access.

**Interview angle:** *"How do you avoid duplicating logic between const and non-const overloads?"* "Implement the const version, and have the non-const version call it via `const_cast<T&>(std::as_const(*this).at(i))`. That's safe, because the object is known to be non-const. C++23's 'deducing this' removes the duplication entirely."

---

## Q209. [M13 | Friend functions and friend classes | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T2 · P67

```cpp
class Box {
    int w;
public:
    Box(int w) : w(w) {}
    friend int area(const Box &b);
    friend class Inspector;
};
int area(const Box &b) { return b.w * b.w; }                          // (1)
class Inspector {
public:
    static int peek(const Box &b) { return b.w; }                     // (2)
};
class Child : public Inspector {
public:
    static int peek2(const Box &b) { return b.w; }                    // (3)
};
```

Which line fails to compile?

- A. (1): `area` isn't a member, so it can't access `w`.
- B. (3): friendship **isn't inherited**, so `Child` has no access to `Box::w`.
- C. (2): a friend class can only access public members.
- D. None: all three compile.

**Answer:** B

**Behaviour class:** Ill-formed (line 3: "'int Box::w' is private within this context").

**In simple words:** A `friend` declaration is a **personal invitation**: `Box` says "this particular function and this particular class may see my private parts". The invitation goes to `Inspector` by name. It doesn't extend to `Inspector`'s children, any more than a friend's kids get a key to your house.

**Why:**

| Line | Who accesses `Box::w` | Is it a friend? | Result |
|---|---|---|---|
| (1) | the free function `area` | ✓ declared a friend | compiles |
| (2) | `Inspector::peek` | ✓ (`Inspector` is a friend class) | compiles |
| (3) | `Child::peek2` | ✗ friendship isn't inherited | **error** |

**Why the other options fail:**
- A: The whole point of a friend function is access without being a member.
- C: Friends get full access, including private and protected members.
- D: Misses the inheritance rule.

**Concept recap:** Typical legitimate friends: `operator<<` / `operator>>` for streams (they must be non-members, because the stream is the left operand; M17), symmetric binary operators, and tightly coupled helper classes (an iterator and its container).

**Trap alert:** Granting friendship widens encapsulation. Prefer a public member function, or a narrow friend, over befriending whole classes.

**Interview angle:** *"Does friend break encapsulation?"* "It extends the class's own boundary to a specific function or class that the class itself chooses. It's not a back door anyone can use. Used narrowly (stream operators, iterators), it improves encapsulation compared with adding public getters just for those helpers."

---

## Q210. [M13 | Rules of friendship | Medium | MSQ | [C++] · C++17 | Interview] · Type — concept · P67

Select **all** true statements.

- A. Friendship is granted by a class; a function or class can't declare itself a friend of another class.
- B. If class A is a friend of class B, then B is automatically a friend of A.
- C. If A is a friend of B, and B is a friend of C, then A is a friend of C.
- D. A friend function declared inside a class isn't a member function, and has no `this`.
- E. Whether a `friend` declaration takes effect depends on whether it appears in a `public` or `private` section.

**Answer:** A, D

**In simple words:** Friendship in C++ is **one-way, non-transitive, not inherited, and granted only by the owner**. It's exactly like lending your house key: you decide who gets it, they don't automatically lend you theirs, and they can't pass it on to their friends.

**Why each option:**
- **A ✓:** Only the class whose privates are being opened can write the `friend` declaration.
- **B ✗:** Not symmetric.
- **C ✗:** Not transitive.
- **D ✓:** It's an ordinary non-member function with special access.
- **E ✗:** Access specifiers don't affect `friend` declarations. They can go in any section.

**Concept recap:** The four "nots" of friendship: not symmetric, not transitive, not inherited, not claimed by the friend. They also appear among the summary tables at the end of the bank.

**Trap alert:** A friend function **defined** inside the class body (the "hidden friend" idiom) is found only through argument-dependent lookup (ADL), i.e. when an argument has the class's type. That's how `operator<<` is often written in modern C++.

**Interview angle:** *"Is friendship inherited or transitive?"* "Neither. And it isn't symmetric. Each relationship must be granted explicitly by the class being accessed."

---

## Q211. [M13 | Empty classes, empty bases, the vptr | Hard | MCQ | [C++] · C++17 | Interview] · Type T1 · P33

Assume the typical-ABI conventions stated at the top (one 8-byte vptr per polymorphic object; empty bases take no space).

```cpp
#include <iostream>
struct Empty {};
struct E2 : Empty { int x; };
struct V  { virtual void f() {} };
struct VI { int x; virtual void f() {} };
struct S  { char c; Empty e; };
int main() {
    std::cout << sizeof(Empty) << ' ' << sizeof(E2) << ' ' << sizeof(V) << ' '
              << sizeof(VI) << ' ' << sizeof(S) << '\n';
}
```

What does this program print?

- A. `0 4 8 12 1`
- B. `1 8 8 16 2`
- C. `1 4 8 16 2`
- D. `1 4 0 4 2`

**Answer:** C

**Behaviour class:** Well-defined. The sizes are implementation-defined, except that `sizeof(Empty)` must be at least 1, and they're fixed by the stated ABI conventions.

**In simple words:**
- Every complete object needs its **own address**, so even an empty class takes at least 1 byte.
- When an empty class is used as a **base**, the compiler can squeeze it out entirely (the empty-base optimisation).
- A class with virtual functions secretly stores a **pointer to its function table** (the vptr), which costs 8 bytes and pushes other members to 8-byte alignment.

**Why:**

| Type | Layout | sizeof |
|---|---|---|
| `Empty` | no members, but must be ≥ 1 so that distinct objects have distinct addresses | **1** |
| `E2 : Empty` | the empty base takes 0 bytes (EBO), plus `int x` | **4** |
| `V` | the vptr | **8** |
| `VI` | vptr (8) + int (4) + 4 padding (alignment 8) | **16** |
| `S` | char (1) + an Empty **member** (1; EBO applies only to bases) | **2** |

**Why the other options fail:**
- A: `sizeof` is never 0 for a complete object, and VI needs trailing padding.
- B: Ignores the empty-base optimisation in `E2`.
- D: Forgets the vptr.

**Concept recap:** Member functions (virtual or not) and static data members **don't** take space in objects. Only non-static data members, the vptr(s), padding, and virtual-base bookkeeping do.

**Trap alert:** C++20's `[[no_unique_address]]` lets an empty **member** take no space too, which is what libraries use for stateless allocators and comparators.

**Interview angle:** *"What does adding one virtual function cost?"* "One vptr per object (8 bytes on 64-bit, plus any padding it causes), one vtable per class, and an indirect call that usually can't be inlined. That's why small value types such as `std::pair` are never polymorphic."

---

## Q212. [M13 | What counts toward an object's size | Medium | MCQ | [C++] · C++17 | OA] · Type T1 · P33

```cpp
#include <iostream>
class P {
    char tag;
    double val;
    static int total;
    void f() {}
    int g() const { return 0; }
public:
    int id;
};
int main() { std::cout << sizeof(P) << '\n'; }
```

Under natural alignment, what does this program print?

- A. `24`
- B. `13`
- C. `32`
- D. `17`

**Answer:** A

**Behaviour class:** Well-defined (layout as per the stated rule).

**In simple words:** An object only stores its **own data**. Functions live once in the program's code, and static members live once, outside every object, so neither adds to `sizeof`. The size comes from `tag`, `val` and `id`, plus padding. Access sections like `public:` don't change the order members are laid out in.

**Why:**

| Member | Offset | Size |
|---|---|---|
| `tag` (char) | 0 | 1 |
| padding | 1–7 | 7 |
| `val` (double) | 8 | 8 |
| `id` (int) | 16 | 4 |
| trailing padding (alignment 8) | 20–23 | 4 |
| **total** | | **24** |

**Why the other options fail:**
- B: Adds the member sizes with no padding.
- C: Also counts the static `int` and/or the functions.
- D: Adds the static member but no padding.

**Concept recap:** Within each access section, members are laid out in declaration order. (Since C++23 the order is also kept across sections; before that, implementations were *allowed* to reorder between sections, but mainstream ABIs never did.)

**Trap alert:** Putting `id` right after `tag` (char, int, double) would make it 16 bytes. Member ordering matters in C++ exactly as in C (Q132–Q133).

**Interview angle:** *"Do member functions increase object size?"* "No. Code is stored once per function. Only a vptr (if the class has virtual functions) adds per-object overhead related to functions."

---

## Q213. [M13 | Nested classes and access | Easy | MCQ | [C++] · C++17 | OA] · Type T1

```cpp
#include <iostream>
class Outer {
    int secret = 42;
    static int s;
public:
    class Inner {
    public:
        int read(const Outer &o) { return o.secret; }
        int readStatic() { return s; }
    };
};
int Outer::s = 7;
int main() {
    Outer o;
    Outer::Inner in;
    std::cout << in.read(o) + in.readStatic() << ' ' << sizeof(Outer) << '\n';
}
```

What does this program print?

- A. Compile error: `Inner` can't access `Outer`'s private members
- B. `42 4`
- C. `49 8`
- D. `49 4`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** A nested class is **a member of the outer class**, so (like any member) it can see the outer class's private parts. But it's still a **separate type**: an `Outer` object doesn't contain an `Inner` object, and an `Inner` object doesn't automatically know about any `Outer` object, so you must pass one in.

**Why:**

| Item | Value |
|---|---|
| `in.read(o)` → `o.secret` (allowed: a nested class is a member, C++11 [class.access.nest]) | 42 |
| `in.readStatic()` → `Outer::s` | 7 |
| sum | **49** |
| `sizeof(Outer)` = just `secret` (`s` is static; `Inner` is only a type) | **4** |

**Why the other options fail:**
- A: True in old C++98 rules, but not since C++11.
- B: Forgets `readStatic`.
- C: Thinks the nested type adds to `Outer`'s size.

**Concept recap:** Access works **one way**. The outer class has **no** special access to the inner class's private members, unless the inner class declares it a friend.

**Trap alert:** Unlike Java's non-static inner classes, a C++ nested class has **no hidden pointer** to an outer object. You must pass one explicitly, as `read(o)` does.

**Interview angle:** *"Where are nested classes used?"* "For implementation details tied to one class, like a linked list's `Node` or a container's `iterator` (`std::vector<int>::iterator`). Nesting keeps the name scoped and can make it private."

---

## Q214. [M13 | Access-control errors | Easy | MCQ | [C++] · C++17 | OA] · Type T2

```cpp
class Account {
    double balance = 0;
protected:
    int id = 1;
public:
    void deposit(double x) { balance += x; }
};
int main() {
    Account a;
    a.deposit(5);      // L1
    a.balance = 9;     // L2
    a.id = 3;          // L3
}
```

Which lines fail to compile?

- A. Only L2
- B. L2 and L3
- C. L1, L2 and L3
- D. Only L3

**Answer:** B

**Behaviour class:** Ill-formed (L2, L3).

**In simple words:** Access specifiers are **doors with different locks**:
- **public:** anyone can use it.
- **protected:** the class itself and its derived classes can use it.
- **private:** only the class itself (and its friends) can use it.

`main` is an outsider, so it can only go through the public door, `deposit`.

**Why:**

| Line | Member | Access | From `main` |
|---|---|---|---|
| L1 | `deposit` | public | ✓ |
| L2 | `balance` | private | ✗ |
| L3 | `id` | protected | ✗ (main isn't a derived class) |

**Why the other options fail:**
- A: Thinks protected members are accessible from outside.
- C: Public member functions are always callable.
- D: Thinks private members are accessible if they have default initializers.

**Concept recap:** `deposit` may modify `balance` because access is checked **per class**: member functions of `Account` can touch every `Account` member. Encapsulation means the outside world changes `balance` only through `deposit`, which can enforce rules (e.g. reject negative amounts).

**Trap alert:** Access control is checked at **compile time** and only on **names**. It doesn't affect layout or run-time behaviour, and it doesn't stop pointer tricks (which are UB anyway).

**Interview angle:** *"Why make data members private?"* "To protect invariants, such as a balance never going negative, behind a validated interface; to allow changing the representation later without breaking callers; and to keep the class's responsibilities clear." This is the textbook definition of **encapsulation**.

---

## Q215. [M14 | Which constructor or assignment runs? | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
struct T {
    T()                     { std::cout << "D"; }
    T(int)                  { std::cout << "I"; }
    T(const T &)            { std::cout << "C"; }
    T &operator=(const T &) { std::cout << "A"; return *this; }
};
int main() {
    T a;
    T b(5);
    T c = b;
    T d;
    d = a;
    T e = 7;
    std::cout << '\n';
}
```

What does this program print?

- A. `DIADAI`
- B. `DICDAIC`
- C. `DICDAI`
- D. `DIDDAI`

**Answer:** C

**Behaviour class:** Well-defined.

**In simple words:** The `=` sign means two different things in C++:
- In a **declaration** (`T c = b;`) it means "**create** c as a copy" → a **constructor**.
- On an **existing** object (`d = a;`) it means "**overwrite** d" → the **assignment operator**.

The rule of thumb: if a new name is being born, a constructor runs.

**Why:**

| Line | Statement | What runs | Prints |
|---|---|---|---|
| 1 | `T a;` | default constructor | D |
| 2 | `T b(5);` | `T(int)` | I |
| 3 | `T c = b;` | copy **constructor** (a new object) | C |
| 4 | `T d;` | default constructor | D |
| 5 | `d = a;` | copy **assignment** (`d` already exists) | A |
| 6 | `T e = 7;` | `T(int)`, a converting constructor. In C++17 the temporary is elided: e is built directly from 7 | I |

Prints `DICDAI`.

**Why the other options fail:**
- A: Treats `T c = b;` as an assignment.
- B: Adds a copy for `T e = 7;`. In C++17, initialising from a prvalue of the same type constructs directly, with no copy. (Before C++17, the copy could be elided but was allowed to happen.)
- D: Treats `T c = b;` as default construction followed by assignment.

**Concept recap:** Forms of initialisation:
- `T x;` → default constructor
- `T x(args);` → direct-initialisation
- `T x = expr;` → copy-initialisation (uses non-`explicit` constructors)
- `T x{args};` → list-initialisation

**Trap alert:** If `T(int)` were `explicit`, `T e = 7;` would **not compile**, but `T e(7);` would (Q219).

**Interview angle:** *"Copy constructor vs copy assignment: when is each called?"* "Copy constructor: when a **new** object is initialised from an existing one (declaration with `=`, pass by value, return by value without elision). Copy assignment: when an **existing** object gets a new value."

---

## Q216. [M14 | A missing default constructor | Easy | MCQ | [C++] · C++17 | OA] · Type T2

```cpp
struct P {
    P(int) {}
};
int main() {
    P a(1);   // L1
    P b;      // L2
    P c{2};   // L3
    P d = 3;  // L4
}
```

Which line fails to compile?

- A. L2
- B. L4
- C. L1 and L2
- D. None

**Answer:** A

**Behaviour class:** Ill-formed (L2: "no matching function for call to 'P::P()'").

**In simple words:** The compiler gives a class a free "no-argument" constructor **only if you don't write any constructor yourself**. As soon as you write `P(int)`, you've taken charge, and the free default constructor disappears. So `P b;` has nothing to call.

**Why:**

| Line | Needs | Available? |
|---|---|---|
| L1 `P a(1)` | `P(int)` | ✓ |
| L2 `P b` | `P()` | ✗ not generated, because a user-declared constructor exists |
| L3 `P c{2}` | `P(int)` | ✓ |
| L4 `P d = 3` | `P(int)`, non-explicit → implicit conversion | ✓ |

**Why the other options fail:**
- B: `P(int)` isn't `explicit`, so copy-initialisation from 3 works.
- C: L1 calls the declared constructor directly.
- D: Misses the suppressed default constructor.

**Concept recap:** To get the default constructor back, write `P() = default;` (or `P() : …{}`). Containers like `std::vector<P> v(10);` and arrays `P arr[5];` need a default constructor; `std::vector<P> v(10, P(1));` doesn't.

**Trap alert:** A class member of type `P` inside another class makes **that** class's implicit default constructor unusable too, unless the member gets an initializer (in the constructor's initializer list, or a default member initializer).

**Interview angle:** *"When does the compiler generate a default constructor?"* "Only when the class has no user-declared constructors. And it's deleted if some member or base can't be default-constructed (e.g. a reference member, a `const` member without an initializer, or a member like `P`)."

---

## Q217. [M14 | The most vexing parse | Medium | MCQ | [C++] · C++17 | OA+Interview] · Type T7

```cpp
#include <iostream>
struct W { W() { std::cout << "W"; } };
int main() {
    W a;
    W b();
    W c{};
    W *p = new W;
    delete p;
    std::cout << '\n';
}
```

What does this program print?

- A. `WWWW`
- B. `WW`
- C. Compile error: `b` is never used
- D. `WWW`

**Answer:** D

**Behaviour class:** Well-defined.

**In simple words:** `W b();` **looks** like "create b using the no-argument constructor", but C++ reads it as **declaring a function** named `b` that takes nothing and returns a `W`. No object is created. This surprising rule is called the "most vexing parse". Writing `W b{};` (braces) or `W b;` avoids it.

**Why:**

| Line | Parsed as | Constructor runs? |
|---|---|---|
| `W a;` | an object, default-initialised | ✓ W |
| `W b();` | **a function declaration** `W b(void)` | ✗ |
| `W c{};` | an object, value-initialised | ✓ W |
| `new W` | a heap object | ✓ W |

Three W's → `WWW`.

**Why the other options fail:**
- A: Counts `b` as an object.
- B: Misses either `c` or the heap object.
- C: Declaring a function that's never defined or used is perfectly legal (at most a warning).

**Concept recap:** The rule: "if it can be parsed as a declaration, it is one" ([dcl.ambig.res]). The worse version: `std::vector<int> v(std::istream_iterator<int>(cin), std::istream_iterator<int>());` declares a function too.

**Trap alert:** Calling `b.method()` afterwards produces a confusing error ("request for member in 'b', which is of non-class type"), and that's when people first discover this rule.

**Interview angle:** *"What's the most vexing parse, and how does C++11 help?"* "An ambiguous statement that looks like object creation is parsed as a function declaration. Uniform brace initialisation (`T x{};`) can't be read as a declaration, so it removes the ambiguity."

---

## Q218. [M14 | Copies from pass-by-value, plus guaranteed elision | Hard | MCQ | [C++] · C++17 | OA+Interview] · Type T7 · P71

```cpp
#include <iostream>
struct T {
    T()          { std::cout << "D"; }
    T(const T &) { std::cout << "C"; }
    ~T()         { std::cout << "~"; }
};
void byVal(T) {}
void byRef(const T &) {}
int main() {
    T arr[2];
    byVal(arr[0]);
    byRef(arr[1]);
    T x = T();
    std::cout << '|';
}
```

What is the **complete** output?

- A. `DDC~CDC~|~~~`
- B. `DDC~D|~~~`
- C. `DDCD|~~~~`
- D. `DDC~~D|~~`

**Answer:** B

**Behaviour class:** Well-defined.

**In simple words:**
- Building an array builds **each element**.
- Passing **by value** hands the function a **photocopy** (copy constructor), which is shredded when the call finishes. Passing **by reference** hands over the original, so nothing is copied.
- `T x = T();` in C++17 builds `x` **directly**, with no temporary and no copy.
- At the end of `main`, objects are destroyed in **reverse order of creation**.

**Why:**

| Step | Event | Prints |
|---|---|---|
| `T arr[2];` | arr[0], arr[1] default-constructed | `DD` |
| `byVal(arr[0])` | the parameter is copy-constructed | `C` |
| end of that call | the parameter is destroyed (before the next statement) | `~` |
| `byRef(arr[1])` | binds a reference, no object made | — |
| `T x = T();` | C++17: `x` built directly from the prvalue (no temporary, no copy) | `D` |
| `'|'` | | `|` |
| end of `main` | destroy in reverse order: x, arr[1], arr[0] | `~~~` |

Output: `DDC~D|~~~`.

**Why the other options fail:**
- A: Adds an extra copy and destruction for `T x = T();` (pre-C++17 thinking with no elision), plus a spurious extra `D`.
- C: Never destroys the by-value parameter (so one destruction too many at the end).
- D: Double-destroys the parameter and loses array destructors.

**Concept recap:**
- **Guaranteed copy elision (C++17):** when an object is initialised from a prvalue of the same type (`T x = T();`, `return T();`), no temporary exists, so no copy or move constructor is called, and none even needs to exist.
- **NRVO** (returning a *named* local) is still optional (Q238 area).

**Trap alert:** Passing large objects **by value** silently creates copies (and destroys them). That's why `const T&` is the default for read-only parameters.

**Interview angle:** *"What changed with copy elision in C++17?"* "Elision from prvalues became mandatory. Even non-copyable, non-movable types can be returned from factory functions by value. Named-return-value optimisation remains optional."

---

## Q219. [M14 | explicit constructors | Medium | MSQ | [C++] · C++17 | OA+Interview] · Type T2 · P53

```cpp
struct Meters { explicit Meters(double v) : v(v) {} double v; };
struct Feet   { Feet(double v) : v(v) {} double v; };
void walkM(Meters) {}
void walkF(Feet) {}
```

Select **all** lines that **compile**.

- A. `walkF(3.0);`
- B. `walkM(3.0);`
- C. `walkM(Meters(3.0));`
- D. `Meters m = 2.0;`
- E. `Meters n{2.0};`

**Answer:** A, C, E

**In simple words:** A normal one-argument constructor also acts as a **silent converter**: anywhere a `Feet` is expected, a plain `3.0` is quietly turned into one. Marking a constructor `explicit` says **"don't convert silently. Make the programmer write `Meters(...)`."** That prevents accidentally passing a bare number where a unit was meant.

**Why each option:**
- **A ✓:** `Feet(double)` isn't explicit, so 3.0 converts implicitly.
- **B ✗:** That would need an implicit conversion, which `explicit` forbids ("could not convert '3.0' to 'Meters'").
- **C ✓:** An explicit construction.
- **D ✗:** `=` with a different type is **copy-initialisation**, which only uses non-explicit constructors.
- **E ✓:** Direct-list-initialisation may use explicit constructors.

**Concept recap:**

| Form | Uses explicit constructors? |
|---|---|
| `T x(arg);` / `T x{arg};` | ✓ (direct-initialisation) |
| `T x = arg;` / passing an argument / `return arg;` | ✗ (copy-initialisation) |
| `static_cast<T>(arg)` | ✓ |

**Trap alert:** Without `explicit`, `std::vector<int> v = 5;` would create a vector of 5 elements. The standard library marks that constructor `explicit` precisely to prevent surprises like this.

**Interview angle:** *"When should a constructor be explicit?"* "By default, for **every** single-argument constructor (and, since C++11, multi-argument ones too, to block `= {a, b}` conversions), unless implicit conversion is clearly natural, as with `std::string` from `const char*`. It prevents silent, often expensive conversions and unit mix-ups."

---

## Q220. [M14 | explicit operator bool | Hard | MCQ | [C++] · C++17 | Interview] · Type T2 · P53

```cpp
struct Handle {
    int fd;
    explicit operator bool() const { return fd >= 0; }
};
int main() {
    Handle h{3};
    if (h) {}                              // L1
    bool b = h;                            // L2
    bool c = static_cast<bool>(h);         // L3
    int n = h + 1;                         // L4
    bool d = !h;                           // L5
}
```

Which lines **fail** to compile?

- A. L1, L2 and L4
- B. Only L4
- C. L2 and L4
- D. L2, L4 and L5

**Answer:** C

**Behaviour class:** Ill-formed (L2, L4).

**In simple words:** `explicit operator bool` lets an object answer **"are you valid?"**, but only in places that are clearly asking a yes/no question: `if`, `while`, `!`, `&&`, `||`, `?:`. It refuses to quietly turn into a `bool` value that could then sneak into arithmetic (`h + 1`) or ordinary assignments.

**Why:**

| Line | Context | Allowed? |
|---|---|---|
| L1 `if (h)` | **contextual conversion** to bool | ✓ |
| L2 `bool b = h;` | copy-initialisation needs an implicit conversion | ✗ |
| L3 `static_cast<bool>(h)` | explicit | ✓ |
| L4 `h + 1` | would need an implicit conversion to bool and then to int | ✗ "no match for operator+" |
| L5 `!h` | the operand of `!` is contextually converted | ✓ |

**Why the other options fail:**
- A: `if (h)` is the very context explicit bool is designed for.
- B: Misses L2.
- D: `!` counts as a contextual conversion, so L5 is fine.

**Concept recap:** This is the C++11 replacement for the old "safe bool idiom". `std::unique_ptr`, `std::shared_ptr`, `std::optional` and `std::ifstream` all use `explicit operator bool`, which is why `if (ptr)` and `while (file)` work but `int x = ptr;` doesn't.

**Trap alert:** A **non-explicit** `operator bool` would allow nonsense such as `h1 + h2` or `int n = h;` to compile silently.

**Interview angle:** *"Why is `std::unique_ptr`'s bool conversion explicit?"* "So you can write `if (p)`, but can't accidentally compare two unrelated pointers through their bool values, or add them, or pass them where an int is expected."

---

## Batch 11 summary

| Q | Module | Diff | Format | Answer | Behaviour class |
|---|---|---|---|---|---|
| 201 | M13 | E | MCQ | A | — (concept) |
| 202 | M13 | M | MCQ | C | Well-defined |
| 203 | M13 | E | MCQ | B | — (concept) |
| 204 | M13 | E | MCQ | D | Well-defined |
| 205 | M13 | M | MSQ | A, C | Ill-formed / OK rules |
| 206 | M13 | H | MCQ | A | Well-defined |
| 207 | M13 | M | MCQ | A | Well-defined |
| 208 | M13 | H | MCQ | D | Well-defined |
| 209 | M13 | M | MCQ | B | Ill-formed (line 3) |
| 210 | M13 | M | MSQ | A, D | — (rules) |
| 211 | M13 | H | MCQ | C | Well-defined (ABI sizes) |
| 212 | M13 | M | MCQ | A | Well-defined |
| 213 | M13 | E | MCQ | D | Well-defined |
| 214 | M13 | E | MCQ | B | Ill-formed |
| 215 | M14 | M | MCQ | C | Well-defined |
| 216 | M14 | E | MCQ | A | Ill-formed |
| 217 | M14 | M | MCQ | D | Well-defined |
| 218 | M14 | H | MCQ | B | Well-defined |
| 219 | M14 | M | MSQ | A, C, E | Ill-formed / OK mix |
| 220 | M14 | H | MCQ | C | Ill-formed |

- **Difficulty:** 6 E / 9 M / 5 H.
  - M13 total: 5/6/3, as planned.
  - M14 so far: 1/3/2.
- **Format:** 17 MCQ / 3 MSQ.
  - M13 total: 12/2, as planned.
- **MCQ answer keys:** A×5, B×4, C×4, D×4. Running total Q1–Q220: A 45, B 44, C 44, D 44.

---

**Completed: 220 / 400 Remaining: 180**

**Ledger, newly ✓ this batch:**
- P68 ×3 (Q204, Q205, Q206)
- P52 → 4 (Q203, Q207, Q208) ★ complete
- P67 ×2 (Q209, Q210)
- P33 → 6 (Q211, Q212)
- P53 ×2 (Q219, Q220) ★ complete
- P71 ×1 (Q218)
- P63 ×1, secondary (Q206)

**M13 complete** (Q201–Q214, 14 questions). **M14 in progress** (6 of 24 done).

*Say **"Continue."** for Batch 12 (Q221–Q240: M14 continued — member initialisation order vs initializer-list order, default member initialisers, destruction order for locals, globals, members and arrays, temporaries and reference lifetime extension, `= default` / `= delete`, delegating constructors, value/zero/aggregate initialisation, copy elision vs NRVO, the static-initialisation-order fiasco — then M15 inheritance begins).*
