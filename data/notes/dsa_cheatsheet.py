# =============================================================================
#  DSA / LeetCode CHEAT SHEET — PYTHON  (3.10+; sections match the Go/C++/Java files)
#  Run:  python3 dsa_cheatsheet.py
# =============================================================================
#
#  SECTIONS
#   0. Template & fast I/O              10. Queue
#   1. Numbers, limits, math            11. Deque
#   2. Strings & chars                  12. Heap / priority queue
#   3. Arrays / lists                   13. Pairs, tuples, structs
#   4. Sorting                          14. Linked list & binary tree (LeetCode nodes)
#   5. Binary search                    15. Graphs: BFS, DFS, grid, Dijkstra, topo, DSU
#   6. Hash map                         16. Recursion, memoization, DP tables
#   7. Hash set                         17. Bit manipulation
#   8. Ordered map / set                18. Common patterns (+ itertools)
#   9. Stack                            19. Python gotchas
# =============================================================================

import bisect
import heapq
import itertools
import math
import sys
from collections import Counter, defaultdict, deque, OrderedDict
from dataclasses import dataclass
from functools import cache, cmp_to_key, lru_cache, reduce
from typing import List, Optional

MOD = 10**9 + 7
INF = float("inf")


# =============================================================================
# 0. TEMPLATE & FAST I/O
# =============================================================================
# On LeetCode you only write the method inside `class Solution:`.
# For competitive programming:
def fast_io_template():
    input = sys.stdin.readline  # much faster than built-in input()
    n = int(input())
    a = list(map(int, input().split()))
    sys.stdout.write(f"{n} {a}\n")
    # read everything at once: data = sys.stdin.buffer.read().split()


sys.setrecursionlimit(10**6)  # default is ~1000 — DFS on big inputs will crash without this


# =============================================================================
# 1. NUMBERS, LIMITS, MATH
# =============================================================================
def section_numbers():
    a, f = 10, 10 / 4  # `/` is ALWAYS float division
    print("types:", a, f, int(f), type(a).__name__, type(f).__name__)
    print("big ints never overflow:", 2**100)
    print("inf:", INF, -INF, "sys.maxsize:", sys.maxsize)

    # // and % FLOOR toward -infinity (different from C++/Java/Go!)
    print("div/mod:", 7 // 2, -7 // 2, 7 % 3, -7 % 3)  # 3 -4 1 2
    print("truncate toward zero:", int(-7 / 2), math.trunc(-3.5))  # -3 -3
    print("divmod:", divmod(17, 5))  # (3, 2)

    print("min/max:", min(3, 1, 2), max([3, 1, 2]), min([], default=0))
    print("abs:", abs(-5), "pow:", 2**10, pow(2, 10))
    print("sqrt:", math.sqrt(16), "isqrt:", math.isqrt(17))  # isqrt = exact integer sqrt
    print("floor/ceil:", math.floor(2.7), math.ceil(2.1))
    print("ceil div 7/2:", -(-7 // 2), (7 + 2 - 1) // 2)
    print("gcd/lcm:", math.gcd(12, 18), math.lcm(4, 6))  # lcm: 3.9+
    print("powMod:", pow(2, 30, MOD), "mod inverse:", pow(3, -1, 7))  # 3.8+
    print("mod mult:", 123456789 * 987654321 % MOD)
    print("round/float fmt:", round(2.675, 2), f"{3.14159:.2f}")  # 2.67: floats are inexact; round() uses banker's rounding
    print("factorial/comb/perm:", math.factorial(5), math.comb(5, 2), math.perm(5, 2))
    print("log2:", math.log2(8), "int bit_length:", (8).bit_length())


# =============================================================================
# 2. STRINGS & CHARS   (str is IMMUTABLE; no separate char type)
# =============================================================================
def section_strings():
    s = "hello world"
    print("len:", len(s), "s[0]:", s[0], "s[-1]:", s[-1], "sub:", s[0:5], "step:", s[::2])

    for ch in s[:2]:
        pass
    for i, ch in enumerate(s[:2]):
        pass

    c = "c"
    print("ord/chr:", ord(c) - ord("a"), chr(ord("a") + 3), int("7"))
    print("isdigit/isalpha/isalnum:", "5".isdigit(), c.isalpha(), "_".isalnum())
    print("upper/lower:", "q".upper(), "ABC".lower(), "hello world".title(), "aBc".swapcase())

    # modify: convert to list, change, join back
    lst = list(s)
    lst[0] = "H"
    print("modified:", "".join(lst))

    # building strings: collect parts in a list, then "".join (s += x in a loop can be O(n^2))
    parts = []
    for _ in range(3):
        parts.append("ab")
    print("join:", "-".join(parts))

    print(s.split(), "a,b,c".split(","), "a,b,c".split(",", 1))
    print(repr("  hi  ".strip()), "xxhixx".strip("x"), "  hi".lstrip())
    print("wor" in s, s.find("o"), s.rfind("o"), s.find("z"), s.index("w"))  # index raises if missing
    print(s.startswith("he"), s.endswith("ld"), "banana".count("a"))
    print("aXbXc".replace("X", "-"), "ab" * 3)

    print("int/str:", int("123") + 1, str(456) + "!")
    print("binary:", bin(10), bin(10)[2:], format(10, "08b"), "parse base2:", int("1010", 2), int("ff", 16))
    print("f-string:", f"{1}-{'x'}-{3.14159:.2f}-{42:>5}")

    print("reversed:", "abc"[::-1], "palindrome?:", "racecar" == "racecar"[::-1])
    print("compare:", "apple" < "banana")

    freq = [0] * 26
    for ch in "hello":
        freq[ord(ch) - ord("a")] += 1
    print("freq of 'l':", freq[ord("l") - ord("a")], "Counter:", Counter("hello"))


# =============================================================================
# 3. ARRAYS / LISTS (dynamic array)
# =============================================================================
def section_arrays():
    a = [5, 3, 8, 1]
    b = [0] * 5
    print(a, b, list(range(5)), list(range(10, 0, -3)))

    rows, cols = 3, 4
    grid = [[0] * cols for _ in range(rows)]  # CORRECT 2D
    bad = [[0] * cols] * rows  # WRONG: all rows are the SAME list
    grid[1][2] = 7
    bad[1][2] = 7
    print("grid:", grid)
    print("bad (aliasing bug):", bad)

    a.append(10)          # push back
    last = a[-1]          # back
    a.pop()               # pop back O(1)
    first = a.pop(0)      # pop front O(n)! use deque
    a.insert(1, 99)       # O(n)
    a.remove(99)          # remove first occurrence of VALUE
    del a[0]              # delete by index
    a.extend([7, 7])
    print("after ops:", a, last, first)

    cp = a[:]             # shallow copy (also a.copy(), list(a))
    cp[0] = 1000
    import copy
    deep = copy.deepcopy(grid)  # for nested lists
    print("copy:", a[0], cp[0], deep is grid)

    s = [4, 2, 9, 2]
    print("sum/max/min:", sum(s), max(s), min(s), "argmax:", s.index(max(s)))
    print("count/in:", s.count(2), 9 in s)
    print("reversed:", s[::-1], list(reversed(s)))
    s.reverse()  # in place

    print("slice:", s[1:3], s[-2:], s[:-1])
    print("unique (sorted):", sorted(set([3, 1, 3, 2, 1])))
    print("unique (keep order):", list(dict.fromkeys([3, 1, 3, 2, 1])))

    print("comprehension:", [x * x for x in range(5) if x % 2 == 0])
    print("2D comprehension:", [[i * j for j in range(3)] for i in range(2)])
    print("enumerate/zip:", list(enumerate("ab")), list(zip([1, 2], "xy")))
    print("any/all:", any(x > 8 for x in s), all(x > 0 for x in s))
    print("transpose:", [list(r) for r in zip(*[[1, 2, 3], [4, 5, 6]])])
    print("flatten:", [x for row in grid for x in row][:5])
    m, *rest = [1, 2, 3]  # unpacking
    print("unpack:", m, rest)


# =============================================================================
# 4. SORTING   (Timsort, STABLE)
# =============================================================================
def section_sorting():
    a = [5, 2, 9, 1]
    print("sorted (new list):", sorted(a), "desc:", sorted(a, reverse=True))
    a.sort()  # in place, returns None!
    print("in place:", a)

    iv = [[1, 5], [0, 2], [1, 9]]
    iv.sort(key=lambda x: (x[0], -x[1]))  # start asc, end desc
    print("intervals:", iv)

    people = [("bob", 25), ("amy", 25), ("cat", 20)]
    print("by age then name:", sorted(people, key=lambda p: (p[1], p[0])))
    print("by name length:", sorted(["pear", "fig", "apple"], key=len))
    print("sorted chars:", "".join(sorted("dcba")))

    vals = [30, 10, 20]
    print("argsort:", sorted(range(len(vals)), key=lambda i: vals[i]))

    # custom comparator (when a key is awkward): e.g. "largest number"
    def cmp(x, y):
        return -1 if x + y > y + x else (1 if x + y < y + x else 0)

    nums = ["3", "30", "34", "5", "9"]
    print("largest number:", "".join(sorted(nums, key=cmp_to_key(cmp))))

    print("k smallest/largest:", heapq.nsmallest(2, [7, 2, 9, 4]), heapq.nlargest(2, [7, 2, 9, 4]))


# =============================================================================
# 5. BINARY SEARCH
# =============================================================================
def section_binary_search():
    a = [1, 3, 3, 5, 8]
    lb = bisect.bisect_left(a, 3)   # first >= 3 -> 1  (lower_bound)
    ub = bisect.bisect_right(a, 3)  # first >  3 -> 3  (upper_bound)
    found = lb < len(a) and a[lb] == 3
    print("lb/ub/count:", lb, ub, ub - lb, "found:", found)
    bisect.insort(a, 4)  # insert keeping sorted (O(n))
    print("after insort:", a)
    # key= for bisect: Python 3.10+  e.g. bisect_left(intervals, x, key=lambda iv: iv[0])

    a = [1, 3, 3, 5, 8]
    cond = lambda i: a[i] >= 4  # F..F T..T
    lo, hi = 0, len(a)
    while lo < hi:
        mid = (lo + hi) // 2
        if cond(mid):
            hi = mid
        else:
            lo = mid + 1
    print("first >= 4 at:", lo)

    l, r = 0, 50
    while l < r:
        m = (l + r) // 2
        if m * m >= 50:
            r = m
        else:
            l = m + 1
    print("ceil sqrt(50):", l)


# =============================================================================
# 6. HASH MAP  (dict keeps INSERTION order since 3.7)
# =============================================================================
def section_hash_map():
    m = {}
    m["apple"] = 3
    m["banana"] = m.get("banana", 0) + 1  # get with default
    m["banana"] = m.get("banana", 0) + 1
    print("get:", m["apple"], m.get("cherry"), m.get("cherry", 0), "in:", "apple" in m, "len:", len(m))
    # m["cherry"] -> KeyError!
    del m["apple"]               # KeyError if missing
    m.pop("nothing", None)       # safe delete
    m.setdefault("kiwi", 5)      # insert if missing, return value

    for k, v in m.items():
        print("kv:", k, v)
    print("keys/values:", list(m.keys()), list(m.values()), "sorted keys:", sorted(m))

    freq = Counter([1, 2, 2, 3, 3, 3])
    print("Counter:", freq, freq[3], freq[99], freq.most_common(1))  # missing key -> 0
    print("Counter math:", Counter("aab") - Counter("ab"), Counter("ab") == Counter("ba"))

    g = defaultdict(list)  # missing key -> []
    g[1].append(2)
    g[1].append(3)
    cnt = defaultdict(int)  # missing key -> 0
    cnt["x"] += 1
    print("defaultdict:", dict(g), dict(cnt))

    # tuple keys (lists are NOT hashable)
    seen = {(1, 2): True}
    groups = defaultdict(list)
    for w in ["eat", "tea", "tan"]:
        groups[tuple(sorted(w))].append(w)
    print("tuple keys:", seen[(1, 2)], list(groups.values()))

    print("dict comprehension:", {x: x * x for x in range(3)})
    od = OrderedDict()  # move_to_end / popitem(last=False) -> LRU cache
    od["a"] = 1
    od["b"] = 2
    od.move_to_end("a")
    print("OrderedDict:", list(od), od.popitem(last=False))


# =============================================================================
# 7. HASH SET
# =============================================================================
def section_hash_set():
    st = {5, 7}
    st.add(9)
    print("contains 5:", 5 in st, "size:", len(st))
    st.remove(5)   # KeyError if missing
    st.discard(42)  # safe
    empty = set()   # NOT {} (that's a dict)

    A, B = {1, 2, 3}, {2, 3, 4}
    print("| & - ^:", A | B, A & B, A - B, A ^ B)
    print("subset:", {2, 3} <= A, "frozenset (hashable):", frozenset([1, 2]))
    print("set comprehension:", {x % 3 for x in range(10)}, empty)


# =============================================================================
# 8. ORDERED MAP / SET
# =============================================================================
# Python has no built-in TreeMap. LeetCode HAS `sortedcontainers` preinstalled:
#   from sortedcontainers import SortedList, SortedDict, SortedSet
#   sl = SortedList([5, 1, 9]); sl.add(3); sl.remove(9)        # O(log n)
#   sl[0], sl[-1]                 -> min, max
#   sl.bisect_left(4)             -> index of first >= 4 (ceiling = sl[i])
#   sl.bisect_right(4) - 1        -> index of last <= 4 (floor)
#   sd = SortedDict(); sd.peekitem(0) -> smallest (key, value)
try:
    from sortedcontainers import SortedList
except ImportError:  # minimal fallback so this file runs anywhere (O(n) insert)
    class SortedList:
        def __init__(self, it=()):
            self.a = sorted(it)

        def add(self, x):
            bisect.insort(self.a, x)

        def remove(self, x):
            self.a.pop(bisect.bisect_left(self.a, x))

        def bisect_left(self, x):
            return bisect.bisect_left(self.a, x)

        def bisect_right(self, x):
            return bisect.bisect_right(self.a, x)

        def __getitem__(self, i):
            return self.a[i]

        def __len__(self):
            return len(self.a)

        def __repr__(self):
            return f"SortedList({self.a})"


def section_ordered():
    sl = SortedList([5, 1, 9, 3])
    sl.add(5)  # SortedList keeps duplicates (it's a multiset)
    sl.remove(9)
    i = sl.bisect_left(4)
    ceil4 = sl[i] if i < len(sl) else None
    j = sl.bisect_right(4) - 1
    floor4 = sl[j] if j >= 0 else None
    print("sorted list:", sl, "floor(4):", floor4, "ceil(4):", ceil4, "min/max:", sl[0], sl[-1])


# =============================================================================
# 9. STACK  (list)
# =============================================================================
def section_stack():
    st = []
    st.append(1)
    st.append(2)
    st.append(3)
    top = st[-1]
    st.pop()
    print("stack top:", top, "now:", st, "empty:", not st)

    def is_valid(s):
        pair = {")": "(", "]": "[", "}": "{"}
        stk = []
        for ch in s:
            if ch in pair:
                if not stk or stk[-1] != pair[ch]:
                    return False
                stk.pop()
            else:
                stk.append(ch)
        return not stk

    print("valid parens:", is_valid("({[]})"), is_valid("(]"))


# =============================================================================
# 10. QUEUE  (collections.deque — list.pop(0) is O(n))
# =============================================================================
def section_queue():
    q = deque()
    q.append(1)
    q.append(2)
    q.append(3)
    front = q[0]
    q.popleft()
    print("queue front:", front, "now:", list(q))

    q = deque([1])
    for _ in range(3):
        for _ in range(len(q)):  # process one level
            x = q.popleft()
            q.extend([2 * x, 2 * x + 1])
    print("nodes at level 3:", list(q))


# =============================================================================
# 11. DEQUE
# =============================================================================
def section_deque():
    dq = deque([2, 3])
    dq.appendleft(1)
    print("front/back:", dq[0], dq[-1], "len:", len(dq))
    dq.popleft()
    dq.pop()
    print("remaining:", list(dq))
    window = deque(maxlen=3)  # keeps only last 3 items
    for x in range(5):
        window.append(x)
    d2 = deque([1, 2, 3, 4])
    d2.rotate(1)
    print("maxlen window:", list(window), "rotate:", list(d2))


# =============================================================================
# 12. HEAP / PRIORITY QUEUE   (heapq = MIN heap on a plain list)
# =============================================================================
def section_heap():
    h = [5, 2, 8]
    heapq.heapify(h)  # O(n), in place
    heapq.heappush(h, 1)
    print("min top:", h[0])
    print("popped in order:", [heapq.heappop(h) for _ in range(len(h))])

    # MAX heap: push negatives
    mh = []
    for x in [3, 9, 4]:
        heapq.heappush(mh, -x)
    print("max top:", -mh[0])

    # tuples compare element by element: (priority, tiebreak, item)
    pq = []
    heapq.heappush(pq, (5, 1))
    heapq.heappush(pq, (2, 7))
    heapq.heappush(pq, (2, 3))
    print("tuple heap pop:", heapq.heappop(pq), heapq.heappop(pq))
    # if items aren't comparable (e.g. ListNode), add a counter: (val, idx, node)

    k, kh = 2, []
    for x in [4, 1, 7, 3, 9]:
        heapq.heappush(kh, x)
        if len(kh) > k:
            heapq.heappop(kh)
    print("kth largest (k=2):", kh[0])
    print("heappushpop:", heapq.heappushpop([1, 5], 3), "merge:", list(heapq.merge([1, 4], [2, 3])))


# =============================================================================
# 13. PAIRS, TUPLES, STRUCTS
# =============================================================================
@dataclass(order=True)  # order=True -> comparable by fields in order (usable in heaps/sorting)
class Item:
    val: int
    idx: int
    name: str = ""


def section_pairs():
    p = (3, 4)  # tuple: immutable, hashable
    a, b = p
    print("pair:", p[0], p[1], a, b)

    items = [Item(2, 0, "y"), Item(1, 1, "x")]
    items.sort()
    print("dataclass sorted first:", items[0].name, items[0])

    from collections import namedtuple
    Point = namedtuple("Point", ["r", "c"])
    pt = Point(1, 2)
    print("namedtuple:", pt.r, pt.c, pt)

    x, y = 1, 2
    x, y = y, x
    print("swap:", x, y)


# =============================================================================
# 14. LINKED LIST & BINARY TREE (LeetCode definitions)
# =============================================================================
class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def build_list(vals):
    dummy = cur = ListNode()  # dummy head trick
    for v in vals:
        cur.next = ListNode(v)
        cur = cur.next
    return dummy.next


def reverse_list(head):
    prev = None
    while head:
        head.next, prev, head = prev, head, head.next
    return prev


def middle_node(head):
    slow = fast = head
    while fast and fast.next:
        slow, fast = slow.next, fast.next.next
    return slow


def build_tree(vals):  # LeetCode level-order with None
    if not vals or vals[0] is None:
        return None
    root = TreeNode(vals[0])
    q, i = deque([root]), 1
    while q and i < len(vals):
        node = q.popleft()
        if i < len(vals) and vals[i] is not None:
            node.left = TreeNode(vals[i])
            q.append(node.left)
        i += 1
        if i < len(vals) and vals[i] is not None:
            node.right = TreeNode(vals[i])
            q.append(node.right)
        i += 1
    return root


def section_nodes():
    head = build_list([1, 2, 3, 4, 5])
    print("middle:", middle_node(head).val)
    head = reverse_list(head)
    out = []
    while head:
        out.append(head.val)
        head = head.next
    print("reversed list:", out)

    root = build_tree([3, 9, 20, None, None, 15, 7])

    res = []

    def inorder(n):
        if not n:
            return
        inorder(n.left)
        res.append(n.val)  # appending to outer list is fine; REASSIGNING needs `nonlocal`
        inorder(n.right)

    inorder(root)
    print("inorder:", res)

    def depth(n):
        return 1 + max(depth(n.left), depth(n.right)) if n else 0

    print("max depth:", depth(root))

    levels, q = [], deque([root])
    while q:
        lvl = []
        for _ in range(len(q)):
            n = q.popleft()
            lvl.append(n.val)
            if n.left:
                q.append(n.left)
            if n.right:
                q.append(n.right)
        levels.append(lvl)
    print("level order:", levels)

    it, stk, cur = [], [], root
    while cur or stk:
        while cur:
            stk.append(cur)
            cur = cur.left
        cur = stk.pop()
        it.append(cur.val)
        cur = cur.right
    print("iterative inorder:", it)


# =============================================================================
# 15. GRAPHS
# =============================================================================
class DSU:
    def __init__(self, n):
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x):
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]  # path halving
            x = self.parent[x]
        return x

    def union(self, a, b):
        ra, rb = self.find(a), self.find(b)
        if ra == rb:
            return False
        if self.size[ra] < self.size[rb]:
            ra, rb = rb, ra
        self.parent[rb] = ra
        self.size[ra] += self.size[rb]
        return True


def section_graphs():
    n = 5
    edges = [[0, 1], [0, 2], [1, 3], [2, 3], [3, 4]]
    adj = defaultdict(list)  # or [[] for _ in range(n)]
    for u, v in edges:
        adj[u].append(v)
        adj[v].append(u)

    dist = [-1] * n
    dist[0] = 0
    q = deque([0])
    while q:
        u = q.popleft()
        for v in adj[u]:
            if dist[v] == -1:
                dist[v] = dist[u] + 1
                q.append(v)
    print("BFS dist:", dist)

    visited, order = set(), []

    def dfs(u):
        visited.add(u)
        order.append(u)
        for v in adj[u]:
            if v not in visited:
                dfs(v)

    dfs(0)
    print("DFS order:", order)

    grid = ["110", "010", "001"]
    R, C = len(grid), len(grid[0])
    seen, islands = set(), 0
    for r in range(R):
        for c in range(C):
            if grid[r][c] == "1" and (r, c) not in seen:
                islands += 1
                seen.add((r, c))
                gq = deque([(r, c)])
                while gq:
                    cr, cc = gq.popleft()
                    for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
                        nr, nc = cr + dr, cc + dc
                        if 0 <= nr < R and 0 <= nc < C and grid[nr][nc] == "1" and (nr, nc) not in seen:
                            seen.add((nr, nc))
                            gq.append((nr, nc))
    print("islands:", islands)

    wedges = [[0, 1, 4], [0, 2, 1], [2, 1, 2], [1, 3, 1], [2, 3, 5]]
    wadj = defaultdict(list)
    for u, v, w in wedges:
        wadj[u].append((v, w))
    d = [INF] * 4
    d[0] = 0
    pq = [(0, 0)]  # (dist, node)
    while pq:
        du, u = heapq.heappop(pq)
        if du > d[u]:
            continue  # stale
        for v, w in wadj[u]:
            if du + w < d[v]:
                d[v] = du + w
                heapq.heappush(pq, (d[v], v))
    print("Dijkstra:", d)

    dag = [[0, 1], [0, 2], [1, 3], [2, 3]]
    g, indeg = defaultdict(list), [0] * 4
    for u, v in dag:
        g[u].append(v)
        indeg[v] += 1
    tq = deque(i for i in range(4) if indeg[i] == 0)
    topo = []
    while tq:
        u = tq.popleft()
        topo.append(u)
        for v in g[u]:
            indeg[v] -= 1
            if indeg[v] == 0:
                tq.append(v)
    print("topo:", topo, "has cycle:", len(topo) != 4)

    dsu = DSU(5)
    dsu.union(0, 1)
    dsu.union(3, 4)
    print("DSU same(0,1):", dsu.find(0) == dsu.find(1), "same(1,3):", dsu.find(1) == dsu.find(3))


# =============================================================================
# 16. RECURSION, MEMOIZATION, DP TABLES
# =============================================================================
def section_dp():
    @cache  # = lru_cache(maxsize=None); args must be hashable (tuples, not lists)
    def fib(n):
        return n if n < 2 else fib(n - 1) + fib(n - 2)

    print("fib(50):", fib(50))
    fib.cache_clear()  # clear between LeetCode test cases if defined at module level

    @lru_cache(maxsize=None)
    def paths(i, j):
        if i == 0 or j == 0:
            return 1
        return paths(i - 1, j) + paths(i, j - 1)

    print("unique paths 3x7 (memo):", paths(2, 6))

    a, b = "abcde", "ace"
    dp = [[0] * (len(b) + 1) for _ in range(len(a) + 1)]
    for i in range(1, len(a) + 1):
        for j in range(1, len(b) + 1):
            dp[i][j] = dp[i - 1][j - 1] + 1 if a[i - 1] == b[j - 1] else max(dp[i - 1][j], dp[i][j - 1])
    print("LCS:", dp[-1][-1])

    coins, amount = [1, 2, 5], 11
    cc = [0] + [INF] * amount
    for i in range(1, amount + 1):
        for c in coins:
            if c <= i:
                cc[i] = min(cc[i], cc[i - c] + 1)
    print("coin change:", cc[amount] if cc[amount] != INF else -1)

    nums, subsets, cur = [1, 2, 3], [], []

    def bt(i):
        if i == len(nums):
            subsets.append(cur[:])  # COPY!
            return
        cur.append(nums[i])
        bt(i + 1)
        cur.pop()
        bt(i + 1)

    bt(0)
    print("subsets:", subsets)

    count = 0

    def inc():
        nonlocal count  # needed to REASSIGN an outer variable
        count += 1

    inc()
    print("nonlocal count:", count)


# =============================================================================
# 17. BIT MANIPULATION
# =============================================================================
def section_bits():
    x = 0b10110  # 22
    print("and/or/xor/not:", 6 & 3, 6 | 3, 6 ^ 3, ~0)
    print("shifts:", 1 << 4, 32 >> 2)
    print("popcount:", bin(x).count("1"), x.bit_count())  # bit_count: 3.10+
    print("bit length:", x.bit_length(), "trailing zeros:", (x & -x).bit_length() - 1)
    print("binary:", bin(x), format(x, "08b"))
    i = 2
    print("check bit i:", (x >> i) & 1, "set:", x | (1 << i), "clear:", x & ~(1 << i), "toggle:", x ^ (1 << i))
    print("lowest set bit:", x & -x, "drop lowest:", x & (x - 1), "pow2?:", 16 & 15 == 0)

    mask, subs, s = 0b101, [], 0b101
    while s:
        subs.append(s)
        s = (s - 1) & mask
    print("all masks:", len(range(1 << 3)), "submasks of 101:", subs)
    print("single number:", reduce(lambda p, q: p ^ q, [4, 1, 2, 1, 2]))
    # Python ints are unbounded: to emulate 32-bit, mask with & 0xFFFFFFFF


# =============================================================================
# 18. COMMON PATTERNS (+ itertools)
# =============================================================================
def section_patterns():
    nums = [2, 1, 3, 4, 1]
    pre = [0] + list(itertools.accumulate(nums))
    print("sum[1..3]:", pre[4] - pre[1])

    s = [1, 2, 4, 7, 11]
    l, r, target = 0, len(s) - 1, 15
    while l < r:
        tot = s[l] + s[r]
        if tot == target:
            break
        if tot < target:
            l += 1
        else:
            r -= 1
    print("two-sum pair:", s[l], s[r])

    st, last, best, left = "abcabcbb", {}, 0, 0
    for right, ch in enumerate(st):
        if ch in last and last[ch] >= left:
            left = last[ch] + 1
        last[ch] = right
        best = max(best, right - left + 1)
    print("longest unique window:", best)

    arr = [2, 1, 2, 4, 3]
    nge, stk = [-1] * len(arr), []
    for i in range(len(arr) - 1, -1, -1):
        while stk and arr[stk[-1]] <= arr[i]:
            stk.pop()
        if stk:
            nge[i] = arr[stk[-1]]
        stk.append(i)
    print("next greater:", nge)

    best_sum, cur = -INF, 0
    for v in [-2, 1, -3, 4, -1, 2, 1, -5, 4]:
        cur = max(v, cur + v)
        best_sum = max(best_sum, cur)
    print("max subarray:", best_sum)

    print("permutations:", list(itertools.permutations([1, 2, 3], 2))[:3])
    print("combinations:", list(itertools.combinations([1, 2, 3], 2)))
    print("product:", list(itertools.product("ab", repeat=2)))
    print("pairwise:", list(itertools.pairwise([1, 2, 3])))  # 3.10+
    print("groupby:", [(k, len(list(g))) for k, g in itertools.groupby("aaabcc")])


# =============================================================================
# 19. PYTHON GOTCHAS
# =============================================================================
#  * [[0]*m]*n creates n references to ONE row. Use [[0]*m for _ in range(n)].
#  * Default mutable args persist: def f(a=[]) -> use a=None.
#  * Recursion limit ~1000: sys.setrecursionlimit(10**6), or go iterative.
#  * -7 // 2 == -4 and -7 % 2 == 1 (floor). Use int(a / b) for truncation (or C-style math).
#  * list.pop(0) and `x in list` are O(n). Use deque / set.
#  * sort() returns None; sorted() returns a new list.
#  * Strings are immutable: build with list + "".join().
#  * Assigning to an outer variable inside a nested function needs `nonlocal`.
#  * Copy lists when storing snapshots in backtracking: res.append(cur[:]).
#  * @cache on a method inside `class Solution` keeps state across calls; define it inside the method.
#  * Python is slow: ~1e7 simple ops/sec. Prefer builtins (sum, sorted, Counter) over manual loops.


if __name__ == "__main__":
    for name, fn in [
        ("1. Numbers", section_numbers), ("2. Strings", section_strings), ("3. Arrays", section_arrays),
        ("4. Sorting", section_sorting), ("5. Binary search", section_binary_search),
        ("6. Hash map", section_hash_map), ("7. Hash set", section_hash_set),
        ("8. Ordered set", section_ordered), ("9. Stack", section_stack), ("10. Queue", section_queue),
        ("11. Deque", section_deque), ("12. Heap", section_heap), ("13. Pairs", section_pairs),
        ("14. Nodes", section_nodes), ("15. Graphs", section_graphs), ("16. DP", section_dp),
        ("17. Bits", section_bits), ("18. Patterns", section_patterns),
    ]:
        print(f"\n===== {name} =====")
        fn()
