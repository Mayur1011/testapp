// =============================================================================
//  DSA / LeetCode CHEAT SHEET — GO  (Go 1.21+; sections match the C++/Py/Java files)
//  Run:  go run dsa_cheatsheet.go
// =============================================================================
//
//  SECTIONS
//   0. Template & fast I/O              10. Queue
//   1. Numbers, limits, math            11. Deque
//   2. Strings & chars                  12. Heap / priority queue
//   3. Arrays / slices (dynamic array)  13. Pairs, tuples, structs
//   4. Sorting                          14. Linked list & binary tree (LeetCode nodes)
//   5. Binary search                    15. Graphs: BFS, DFS, grid, Dijkstra, topo, DSU
//   6. Hash map                         16. Recursion, memoization, DP tables
//   7. Hash set                         17. Bit manipulation
//   8. Ordered map / set                18. Common patterns (prefix sum, 2 ptr, window, mono stack)
//   9. Stack                            19. Go gotchas
// =============================================================================

package main

import (
	"bufio"
	"cmp"
	"container/heap"
	"container/list"
	"fmt"
	"math"
	"math/bits"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// =============================================================================
// 0. TEMPLATE & FAST I/O
// =============================================================================
// On LeetCode you only write the function, e.g.  func twoSum(nums []int, target int) []int { ... }
// For competitive programming (Codeforces etc.) use buffered I/O:
func fastIOTemplate() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush() // MUST flush or nothing is printed

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}
	fmt.Fprintln(out, n, a)
}

// =============================================================================
// 1. NUMBERS, LIMITS, MATH
// =============================================================================
const MOD = 1_000_000_007

func abs(x int) int { // no integer abs in std lib
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int { return a / gcd(a, b) * b }

func powMod(base, exp, mod int) int { // fast exponentiation O(log exp)
	res := 1
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func sectionNumbers() {
	// int is 64-bit on LeetCode (amd64). int32/int64/uint64 also exist. No implicit conversion!
	var a int = 10
	var b int64 = int64(a) * 3
	var f float64 = float64(a) / 4 // 2.5
	fmt.Println("types:", a, b, f, int(f))

	fmt.Println("limits:", math.MaxInt, math.MinInt, math.MaxInt32, math.MinInt32)
	inf := math.MaxInt  // use as "infinity" for ints
	finf := math.Inf(1) // +Inf float
	fmt.Println("inf:", inf, finf)

	// Division truncates toward zero; % keeps the sign of the dividend
	fmt.Println("div/mod:", 7/2, -7/2, 7%3, -7%3)  // 3 -3 1 -1
	fmt.Println("non-negative mod:", ((-7%3)+3)%3) // 2

	fmt.Println("min/max builtin:", min(3, 1, 2), max(3, 1, 2)) // Go 1.21+
	fmt.Println("abs:", abs(-5), math.Abs(-2.5))
	fmt.Println("sqrt/pow:", math.Sqrt(16), math.Pow(2, 10), int(math.Sqrt(float64(17)))) // 4 1024 4
	fmt.Println("floor/ceil:", math.Floor(2.7), math.Ceil(2.1))
	fmt.Println("ceil div 7/2:", (7+2-1)/2) // integer ceil for positives
	fmt.Println("gcd/lcm:", gcd(12, 18), lcm(4, 6))
	fmt.Println("powMod:", powMod(2, 30, MOD))
	x, y := 123456789, 987654321
	fmt.Println("mod mult:", x%MOD*(y%MOD)%MOD) // take % after every multiply
}

// =============================================================================
// 2. STRINGS & CHARS
// =============================================================================
// Go strings are IMMUTABLE byte sequences. s[i] is a byte (uint8), not a char.
// For LeetCode (mostly ASCII) work with bytes; use runes only for Unicode.
func sectionStrings() {
	s := "hello world"
	fmt.Println("len:", len(s), "s[0]:", s[0], string(s[0]), "sub:", s[0:5]) // 104 "h" "hello"

	// iterate bytes
	for i := 0; i < 3; i++ {
		c := s[i]
		_ = c
	}
	// iterate runes (unicode-safe); i = byte index
	for i, r := range "héllo" {
		_ = i
		_ = r
	}

	// char <-> int
	c := byte('c')
	fmt.Println("c-'a':", c-'a', "'a'+3:", string(rune('a'+3)), "'7'-'0':", '7'-'0')
	fmt.Println("isDigit/isLetter:", c >= '0' && c <= '9', unicode.IsLetter(rune(c)), unicode.IsDigit('5'))
	fmt.Println("upper:", string(unicode.ToUpper('q')), strings.ToUpper("abc"), strings.ToLower("ABC"))

	// modify a string -> convert to []byte, change, convert back
	bs := []byte(s)
	bs[0] = 'H'
	fmt.Println("modified:", string(bs))

	// building strings: use strings.Builder (s += x in a loop is O(n^2))
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		sb.WriteString("ab")
		sb.WriteByte('-')
	}
	sb.WriteRune('✓')
	fmt.Println("builder:", sb.String())

	// common functions
	fmt.Println(strings.Split("a,b,c", ","))           // [a b c]
	fmt.Println(strings.Fields("  a  b   c "))         // [a b c]  (split on whitespace)
	fmt.Println(strings.Join([]string{"x", "y"}, "/")) // x/y
	fmt.Println(strings.TrimSpace("  hi  "), strings.Trim("xxhixx", "x"))
	fmt.Println(strings.Contains(s, "wor"), strings.Index(s, "o"), strings.LastIndex(s, "o"), strings.Index(s, "z")) // true 4 7 -1
	fmt.Println(strings.HasPrefix(s, "he"), strings.HasSuffix(s, "ld"), strings.Count("banana", "a"))
	fmt.Println(strings.ReplaceAll("aXbXc", "X", "-"), strings.Repeat("ab", 3))

	// number <-> string
	n, err := strconv.Atoi("123") // string -> int
	fmt.Println("Atoi:", n, err, "Itoa:", strconv.Itoa(456)+"!")
	fmt.Println("binary:", strconv.FormatInt(10, 2)) // "1010"
	v, _ := strconv.ParseInt("1010", 2, 64)          // 10
	fmt.Println("parse base2:", v, fmt.Sprintf("%d-%s-%.2f", 1, "x", 3.14159))
	// GOTCHA: string(65) == "A" (rune conversion), NOT "65". Use strconv.Itoa.

	// reverse a string
	rb := []byte("abc")
	slices.Reverse(rb)
	fmt.Println("reversed:", string(rb))

	// compare: ==, <, > work lexicographically
	fmt.Println("compare:", "apple" < "banana", strings.Compare("a", "b"))

	// frequency of lowercase letters (very common)
	var freq [26]int
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			freq[s[i]-'a']++
		}
	}
	fmt.Println("freq of 'l':", freq['l'-'a'])
}

// =============================================================================
// 3. ARRAYS / SLICES (dynamic array)
// =============================================================================
// [N]T is a fixed array (value type, comparable, can be a map key).
// []T is a slice (dynamic array: pointer + len + cap). This is what you use 99% of the time.
func sectionArrays() {
	a := []int{5, 3, 8, 1}
	b := make([]int, 5)     // [0 0 0 0 0]
	c := make([]int, 0, 10) // len 0, cap 10 (pre-allocate)
	fmt.Println(a, b, len(c), cap(c))

	for i := range b { // fill
		b[i] = -1
	}

	// 2D grid  rows x cols
	rows, cols := 3, 4
	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}
	grid[1][2] = 7
	fmt.Println("grid:", grid)

	a = append(a, 10, 11)       // push back
	last := a[len(a)-1]         // back
	a = a[:len(a)-1]            // pop back
	first := a[0]               // front
	a = a[1:]                   // pop front (O(1), but memory not freed)
	a = slices.Insert(a, 1, 99) // insert 99 at index 1 (O(n))
	a = slices.Delete(a, 1, 2)  // delete index 1 (O(n))
	fmt.Println("after ops:", a, last, first)

	// copy — slices SHARE memory when you slice them!
	view := a[0:2]
	view[0] = 1000 // also changes a[0]
	cp := slices.Clone(a)
	cp2 := append([]int(nil), a...)
	cp3 := make([]int, len(a))
	copy(cp3, a)
	fmt.Println("shared:", a[0], "clones:", cp, cp2, cp3)

	// iterate
	for i, v := range a {
		_, _ = i, v // v is a COPY
	}
	// for i := range n { }   // Go 1.22+: loops i = 0..n-1

	s := []int{4, 2, 9, 2}
	sum := 0
	for _, v := range s {
		sum += v
	}
	fmt.Println("sum/max/min:", sum, slices.Max(s), slices.Min(s))
	fmt.Println("contains/index:", slices.Contains(s, 9), slices.Index(s, 2), slices.Equal(s, []int{4, 2, 9, 2}))
	slices.Reverse(s)
	fmt.Println("reversed:", s)

	// remove duplicates (after sorting)
	t := []int{3, 1, 3, 2, 1}
	slices.Sort(t)
	t = slices.Compact(t)
	fmt.Println("unique:", t)

	// fixed array as map key (e.g. group anagrams)
	var key [26]int
	key[0] = 2
	groups := map[[26]int][]string{}
	groups[key] = append(groups[key], "aa")
	fmt.Println("array key map:", len(groups))
}

// =============================================================================
// 4. SORTING
// =============================================================================
func sectionSorting() {
	a := []int{5, 2, 9, 1}
	slices.Sort(a) // ascending (or sort.Ints(a))
	fmt.Println("asc:", a)

	slices.SortFunc(a, func(x, y int) int { return cmp.Compare(y, x) }) // descending
	fmt.Println("desc:", a)

	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] }) // classic way: less(i, j)
	fmt.Println("sort.Slice:", a)

	// sort intervals by start asc, then end desc
	iv := [][]int{{1, 5}, {0, 2}, {1, 9}}
	sort.Slice(iv, func(i, j int) bool {
		if iv[i][0] != iv[j][0] {
			return iv[i][0] < iv[j][0]
		}
		return iv[i][1] > iv[j][1]
	})
	fmt.Println("intervals:", iv)

	// sort structs with cmp.Or (Go 1.22+) — chain of keys
	type P struct {
		name string
		age  int
	}
	ps := []P{{"bob", 25}, {"amy", 25}, {"cat", 20}}
	slices.SortFunc(ps, func(x, y P) int {
		return cmp.Or(cmp.Compare(x.age, y.age), strings.Compare(x.name, y.name))
	})
	fmt.Println("structs:", ps)

	slices.SortStableFunc(ps, func(x, y P) int { return cmp.Compare(y.age, x.age) }) // stable
	fmt.Println("stable:", ps)

	strs := []string{"pear", "apple", "fig"}
	slices.Sort(strs)
	fmt.Println("strings:", strs)

	// sort characters of a string
	b := []byte("dcba")
	slices.Sort(b)
	fmt.Println("sorted chars:", string(b))

	// sort INDICES by value (argsort)
	vals := []int{30, 10, 20}
	idx := []int{0, 1, 2}
	sort.Slice(idx, func(i, j int) bool { return vals[idx[i]] < vals[idx[j]] })
	fmt.Println("argsort:", idx)
}

// =============================================================================
// 5. BINARY SEARCH
// =============================================================================
func sectionBinarySearch() {
	a := []int{1, 3, 3, 5, 8}

	lb := sort.SearchInts(a, 3)                                     // lower_bound: first i with a[i] >= 3 -> 1
	ub := sort.Search(len(a), func(i int) bool { return a[i] > 3 }) // upper_bound: first i with a[i] > 3 -> 3
	pos, found := slices.BinarySearch(a, 5)                         // 3 true
	fmt.Println("lb/ub/count:", lb, ub, ub-lb, "BinarySearch:", pos, found)

	// manual template: first index in [lo, hi) where cond(i) is true (cond is F..F T..T)
	cond := func(i int) bool { return a[i] >= 4 }
	lo, hi := 0, len(a)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if cond(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Println("first >= 4 at:", lo)

	// binary search on the ANSWER: smallest integer sqrt r with r*r >= 50
	l, r := 0, 50
	for l < r {
		m := (l + r) / 2
		if m*m >= 50 {
			r = m
		} else {
			l = m + 1
		}
	}
	fmt.Println("ceil sqrt(50):", l)
}

// =============================================================================
// 6. HASH MAP
// =============================================================================
func sectionHashMap() {
	m := map[string]int{} // or make(map[string]int)
	m["apple"] = 3
	m["banana"]++ // missing key -> zero value, so counting just works
	m["banana"]++

	v, ok := m["cherry"] // ok tells if key exists
	fmt.Println("get:", m["apple"], v, ok, "len:", len(m))
	if _, exists := m["apple"]; exists {
		fmt.Println("apple exists")
	}
	delete(m, "apple") // safe even if key is missing

	// iteration order is RANDOM. Sort keys if you need order:
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys) // Go 1.23+: keys := slices.Sorted(maps.Keys(m))
	for _, k := range keys {
		fmt.Println("kv:", k, m[k])
	}

	// frequency counter
	freq := map[int]int{}
	for _, x := range []int{1, 2, 2, 3, 3, 3} {
		freq[x]++
	}
	fmt.Println("freq:", freq)

	// map of lists (adjacency / grouping)
	g := map[int][]int{}
	g[1] = append(g[1], 2) // append to nil slice is fine
	g[1] = append(g[1], 3)
	fmt.Println("map of lists:", g)

	// composite keys: arrays or structs (slices CANNOT be keys)
	seen := map[[2]int]bool{}
	seen[[2]int{1, 2}] = true
	type pt struct{ r, c int }
	dist := map[pt]int{{0, 0}: 0}
	fmt.Println("composite keys:", seen[[2]int{1, 2}], dist[pt{0, 0}])

	clear(m) // Go 1.21+: remove all entries
}

// =============================================================================
// 7. HASH SET  (Go has no set type -> use map[T]bool or map[T]struct{})
// =============================================================================
func sectionHashSet() {
	set := map[int]bool{}
	set[5] = true
	set[7] = true
	fmt.Println("contains 5:", set[5], "contains 9:", set[9], "size:", len(set))
	delete(set, 5)

	// struct{} uses zero memory
	set2 := map[string]struct{}{}
	set2["a"] = struct{}{}
	_, has := set2["a"]
	fmt.Println("struct{} set has a:", has)

	// intersection
	A := map[int]bool{1: true, 2: true, 3: true}
	B := map[int]bool{2: true, 3: true, 4: true}
	inter := []int{}
	for k := range A {
		if B[k] {
			inter = append(inter, k)
		}
	}
	slices.Sort(inter)
	fmt.Println("intersection:", inter)
}

// =============================================================================
// 8. ORDERED MAP / SET  (TreeMap / std::map / SortedList)
// =============================================================================
// Go's standard library has NO balanced BST. Options:
//
//	(a) Sorted slice + binary search: O(log n) query, O(n) insert/delete. Fine for n <= ~1e5 inserts.
//	(b) LeetCode's Go environment has (historically) shipped github.com/emirpasic/gods —
//	    check the (i) language info on LeetCode before relying on it:
//	      import "github.com/emirpasic/gods/trees/redblacktree"
//	      t := redblacktree.NewWithIntComparator()
//	      t.Put(5, "x"); node, ok := t.Floor(4); node, ok = t.Ceiling(4)
//	      t.Left() (min), t.Right() (max), t.Remove(5), t.Size()
//	(c) Often a heap, a monotonic deque or sorting + two pointers avoids needing one.
type SortedSet struct{ a []int }

func (s *SortedSet) Add(x int) {
	i, found := slices.BinarySearch(s.a, x)
	if !found {
		s.a = slices.Insert(s.a, i, x)
	}
}
func (s *SortedSet) Remove(x int) {
	if i, found := slices.BinarySearch(s.a, x); found {
		s.a = slices.Delete(s.a, i, i+1)
	}
}

// Floor: largest element <= x
func (s *SortedSet) Floor(x int) (int, bool) {
	i := sort.SearchInts(s.a, x+1) // first > x
	if i == 0 {
		return 0, false
	}
	return s.a[i-1], true
}

// Ceiling: smallest element >= x
func (s *SortedSet) Ceiling(x int) (int, bool) {
	i := sort.SearchInts(s.a, x)
	if i == len(s.a) {
		return 0, false
	}
	return s.a[i], true
}

func sectionOrdered() {
	s := &SortedSet{}
	for _, x := range []int{5, 1, 9, 3, 5} {
		s.Add(x)
	}
	s.Remove(9)
	fl, _ := s.Floor(4)
	ce, _ := s.Ceiling(4)
	fmt.Println("sorted set:", s.a, "floor(4):", fl, "ceil(4):", ce, "min/max:", s.a[0], s.a[len(s.a)-1])
}

// =============================================================================
// 9. STACK  (slice)
// =============================================================================
func sectionStack() {
	st := []int{}
	st = append(st, 1, 2, 3) // push
	top := st[len(st)-1]     // peek
	st = st[:len(st)-1]      // pop
	fmt.Println("stack top:", top, "now:", st, "empty:", len(st) == 0)

	// valid parentheses
	isValid := func(s string) bool {
		pair := map[byte]byte{')': '(', ']': '[', '}': '{'}
		stk := []byte{}
		for i := 0; i < len(s); i++ {
			ch := s[i]
			if open, ok := pair[ch]; ok {
				if len(stk) == 0 || stk[len(stk)-1] != open {
					return false
				}
				stk = stk[:len(stk)-1]
			} else {
				stk = append(stk, ch)
			}
		}
		return len(stk) == 0
	}
	fmt.Println("valid parens:", isValid("({[]})"), isValid("(]"))
}

// =============================================================================
// 10. QUEUE  (slice; q = q[1:] is O(1))
// =============================================================================
func sectionQueue() {
	q := []int{}
	q = append(q, 1, 2, 3) // enqueue
	front := q[0]          // peek
	q = q[1:]              // dequeue
	fmt.Println("queue front:", front, "now:", q)

	// level-by-level processing (BFS levels)
	q = []int{1}
	level := 0
	for len(q) > 0 && level < 3 {
		size := len(q)
		for i := 0; i < size; i++ {
			x := q[0]
			q = q[1:]
			q = append(q, x*2, x*2+1)
		}
		level++
	}
	fmt.Println("nodes at level 3:", q)
}

// =============================================================================
// 11. DEQUE  (container/list = doubly linked list, values are `any`)
// =============================================================================
func sectionDeque() {
	dq := list.New()
	dq.PushBack(2)
	dq.PushBack(3)
	dq.PushFront(1)
	fmt.Println("front/back:", dq.Front().Value.(int), dq.Back().Value.(int), "len:", dq.Len())
	dq.Remove(dq.Front()) // pop front
	dq.Remove(dq.Back())  // pop back
	for e := dq.Front(); e != nil; e = e.Next() {
		fmt.Println("deque item:", e.Value.(int))
	}
	// For an index-only deque (sliding window max) a slice also works:
	//   push back: d = append(d, i); pop back: d = d[:len(d)-1]; pop front: d = d[1:]
}

// =============================================================================
// 12. HEAP / PRIORITY QUEUE
// =============================================================================
// (a) Standard way: implement heap.Interface (Len, Less, Swap, Push, Pop)
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] } // use > for a MAX heap
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// (b) Copy-paste generic priority queue (no interface boilerplate per type)
type PQ[T any] struct {
	data []T
	less func(a, b T) bool
}

func NewPQ[T any](less func(a, b T) bool) *PQ[T] { return &PQ[T]{less: less} }
func (pq *PQ[T]) Len() int                       { return len(pq.data) }
func (pq *PQ[T]) Top() T                         { return pq.data[0] }
func (pq *PQ[T]) Push(x T) {
	pq.data = append(pq.data, x)
	i := len(pq.data) - 1
	for i > 0 {
		p := (i - 1) / 2
		if !pq.less(pq.data[i], pq.data[p]) {
			break
		}
		pq.data[i], pq.data[p] = pq.data[p], pq.data[i]
		i = p
	}
}
func (pq *PQ[T]) Pop() T {
	top := pq.data[0]
	n := len(pq.data) - 1
	pq.data[0] = pq.data[n]
	pq.data = pq.data[:n]
	i := 0
	for {
		l, r, best := 2*i+1, 2*i+2, i
		if l < n && pq.less(pq.data[l], pq.data[best]) {
			best = l
		}
		if r < n && pq.less(pq.data[r], pq.data[best]) {
			best = r
		}
		if best == i {
			break
		}
		pq.data[i], pq.data[best] = pq.data[best], pq.data[i]
		i = best
	}
	return top
}

func sectionHeap() {
	h := &MinHeap{5, 2, 8}
	heap.Init(h) // heapify O(n)
	heap.Push(h, 1)
	fmt.Println("min-heap top:", (*h)[0]) // peek = index 0
	out := []int{}
	for h.Len() > 0 {
		out = append(out, heap.Pop(h).(int))
	}
	fmt.Println("popped in order:", out)

	// max heap with generic PQ
	maxPQ := NewPQ(func(a, b int) bool { return a > b })
	for _, x := range []int{3, 9, 4} {
		maxPQ.Push(x)
	}
	fmt.Println("max-heap top:", maxPQ.Top())

	// heap of pairs [dist, node]: smallest dist first, tie -> smaller node
	pq := NewPQ(func(a, b [2]int) bool {
		if a[0] != b[0] {
			return a[0] < b[0]
		}
		return a[1] < b[1]
	})
	pq.Push([2]int{5, 1})
	pq.Push([2]int{2, 7})
	pq.Push([2]int{2, 3})
	fmt.Println("pair heap pop:", pq.Pop(), pq.Pop())

	// k largest elements: keep a MIN heap of size k
	k, kh := 2, NewPQ(func(a, b int) bool { return a < b })
	for _, x := range []int{4, 1, 7, 3, 9} {
		kh.Push(x)
		if kh.Len() > k {
			kh.Pop()
		}
	}
	fmt.Println("kth largest (k=2):", kh.Top())
}

// =============================================================================
// 13. PAIRS, TUPLES, STRUCTS
// =============================================================================
func divmod(a, b int) (int, int) { return a / b, a % b } // multiple return values

func sectionPairs() {
	p := [2]int{3, 4} // fixed array as a pair (comparable, usable as map key)
	fmt.Println("pair:", p[0], p[1])

	type Item struct {
		val, idx int
		name     string
	}
	it := Item{val: 1, idx: 0, name: "x"}
	items := []Item{it, {2, 1, "y"}}
	fmt.Println("struct:", it.val, items[1].name)

	anon := struct{ x, y int }{1, 2} // anonymous struct
	q, r := divmod(17, 5)
	fmt.Println("anon:", anon.x, "divmod:", q, r)

	// swap
	a, b := 1, 2
	a, b = b, a
	fmt.Println("swap:", a, b)
}

// =============================================================================
// 14. LINKED LIST & BINARY TREE (LeetCode definitions)
// =============================================================================
type ListNode struct {
	Val  int
	Next *ListNode
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func buildList(vals []int) *ListNode {
	dummy := &ListNode{} // dummy head trick
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	for head != nil {
		head.Next, prev, head = prev, head, head.Next
	}
	return prev
}

func middleNode(head *ListNode) *ListNode { // slow/fast pointers
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
	}
	return slow
}

// build tree from LeetCode level-order array; use NULL for null
const NULL = math.MinInt

func buildTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == NULL {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	q := []*TreeNode{root}
	i := 1
	for len(q) > 0 && i < len(vals) {
		node := q[0]
		q = q[1:]
		if i < len(vals) && vals[i] != NULL {
			node.Left = &TreeNode{Val: vals[i]}
			q = append(q, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != NULL {
			node.Right = &TreeNode{Val: vals[i]}
			q = append(q, node.Right)
		}
		i++
	}
	return root
}

func sectionNodes() {
	head := buildList([]int{1, 2, 3, 4, 5})
	fmt.Println("middle:", middleNode(head).Val)
	head = reverseList(head)
	for cur := head; cur != nil; cur = cur.Next {
		fmt.Print(cur.Val, " ")
	}
	fmt.Println("<- reversed list")

	root := buildTree([]int{3, 9, 20, NULL, NULL, 15, 7})

	// recursive DFS with a closure (declare first so it can call itself)
	var inorder func(*TreeNode)
	res := []int{}
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)
		res = append(res, n.Val)
		inorder(n.Right)
	}
	inorder(root)
	fmt.Println("inorder:", res)

	var depth func(*TreeNode) int
	depth = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		return 1 + max(depth(n.Left), depth(n.Right))
	}
	fmt.Println("max depth:", depth(root))

	// level order BFS
	levels := [][]int{}
	q := []*TreeNode{root}
	for len(q) > 0 {
		lvl := []int{}
		for size := len(q); size > 0; size-- {
			n := q[0]
			q = q[1:]
			lvl = append(lvl, n.Val)
			if n.Left != nil {
				q = append(q, n.Left)
			}
			if n.Right != nil {
				q = append(q, n.Right)
			}
		}
		levels = append(levels, lvl)
	}
	fmt.Println("level order:", levels)

	// iterative inorder with a stack
	it := []int{}
	stk := []*TreeNode{}
	cur := root
	for cur != nil || len(stk) > 0 {
		for cur != nil {
			stk = append(stk, cur)
			cur = cur.Left
		}
		cur = stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		it = append(it, cur.Val)
		cur = cur.Right
	}
	fmt.Println("iterative inorder:", it)
}

// =============================================================================
// 15. GRAPHS
// =============================================================================
type DSU struct{ parent, size []int }

func NewDSU(n int) *DSU {
	d := &DSU{make([]int, n), make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}
func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]] // path halving
		x = d.parent[x]
	}
	return x
}
func (d *DSU) Union(a, b int) bool {
	ra, rb := d.Find(a), d.Find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

func sectionGraphs() {
	n := 5
	edges := [][]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}}

	// adjacency list (undirected)
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// BFS shortest path (unweighted)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	q := []int{0}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	fmt.Println("BFS dist:", dist)

	// DFS recursive (closure)
	visited := make([]bool, n)
	order := []int{}
	var dfs func(int)
	dfs = func(u int) {
		visited[u] = true
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v)
			}
		}
	}
	dfs(0)
	fmt.Println("DFS order:", order)

	// grid BFS: number of islands
	grid := []string{"110", "010", "001"}
	R, C := len(grid), len(grid[0])
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	seen := make([][]bool, R)
	for i := range seen {
		seen[i] = make([]bool, C)
	}
	islands := 0
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			if grid[r][c] == '1' && !seen[r][c] {
				islands++
				seen[r][c] = true
				gq := [][2]int{{r, c}}
				for len(gq) > 0 {
					cell := gq[0]
					gq = gq[1:]
					for _, d := range dirs {
						nr, nc := cell[0]+d[0], cell[1]+d[1]
						if nr >= 0 && nr < R && nc >= 0 && nc < C && grid[nr][nc] == '1' && !seen[nr][nc] {
							seen[nr][nc] = true
							gq = append(gq, [2]int{nr, nc})
						}
					}
				}
			}
		}
	}
	fmt.Println("islands:", islands)

	// Dijkstra (weighted, non-negative) with the generic PQ
	wedges := [][]int{{0, 1, 4}, {0, 2, 1}, {2, 1, 2}, {1, 3, 1}, {2, 3, 5}}
	wadj := make([][][2]int, 4) // wadj[u] = list of {v, w}
	for _, e := range wedges {
		wadj[e[0]] = append(wadj[e[0]], [2]int{e[1], e[2]})
	}
	d := make([]int, 4)
	for i := range d {
		d[i] = math.MaxInt
	}
	d[0] = 0
	pq := NewPQ(func(a, b [2]int) bool { return a[0] < b[0] }) // {dist, node}
	pq.Push([2]int{0, 0})
	for pq.Len() > 0 {
		top := pq.Pop()
		du, u := top[0], top[1]
		if du > d[u] {
			continue // stale entry
		}
		for _, e := range wadj[u] {
			v, w := e[0], e[1]
			if du+w < d[v] {
				d[v] = du + w
				pq.Push([2]int{d[v], v})
			}
		}
	}
	fmt.Println("Dijkstra:", d)

	// Topological sort (Kahn's algorithm) on a DAG; also detects cycles
	dag := [][]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}} // u -> v
	indeg := make([]int, 4)
	g := make([][]int, 4)
	for _, e := range dag {
		g[e[0]] = append(g[e[0]], e[1])
		indeg[e[1]]++
	}
	tq := []int{}
	for i, dgr := range indeg {
		if dgr == 0 {
			tq = append(tq, i)
		}
	}
	topo := []int{}
	for len(tq) > 0 {
		u := tq[0]
		tq = tq[1:]
		topo = append(topo, u)
		for _, v := range g[u] {
			indeg[v]--
			if indeg[v] == 0 {
				tq = append(tq, v)
			}
		}
	}
	fmt.Println("topo:", topo, "has cycle:", len(topo) != 4)

	// Union-Find
	dsu := NewDSU(5)
	dsu.Union(0, 1)
	dsu.Union(3, 4)
	fmt.Println("DSU same(0,1):", dsu.Find(0) == dsu.Find(1), "same(1,3):", dsu.Find(1) == dsu.Find(3))
}

// =============================================================================
// 16. RECURSION, MEMOIZATION, DP TABLES
// =============================================================================
func sectionDP() {
	// memo with a map (any key type that is comparable)
	memo := map[int]int{}
	var fib func(int) int
	fib = func(n int) int {
		if n < 2 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v
		}
		memo[n] = fib(n-1) + fib(n-2)
		return memo[n]
	}
	fmt.Println("fib(50):", fib(50))

	// memo with a 2D slice initialised to -1 (faster than a map)
	m, n := 3, 7
	mm := make([][]int, m)
	for i := range mm {
		mm[i] = make([]int, n)
		for j := range mm[i] {
			mm[i][j] = -1
		}
	}
	var paths func(int, int) int
	paths = func(i, j int) int {
		if i == 0 || j == 0 {
			return 1
		}
		if mm[i][j] != -1 {
			return mm[i][j]
		}
		mm[i][j] = paths(i-1, j) + paths(i, j-1)
		return mm[i][j]
	}
	fmt.Println("unique paths 3x7 (memo):", paths(m-1, n-1))

	// bottom-up 2D DP: LCS
	a, b := "abcde", "ace"
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	fmt.Println("LCS:", dp[len(a)][len(b)])

	// 1D DP: coin change (min coins)
	coins, amount := []int{1, 2, 5}, 11
	cc := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		cc[i] = math.MaxInt32
		for _, c := range coins {
			if c <= i && cc[i-c]+1 < cc[i] {
				cc[i] = cc[i-c] + 1
			}
		}
	}
	fmt.Println("coin change:", cc[amount])

	// backtracking: subsets
	nums := []int{1, 2, 3}
	subsets := [][]int{}
	cur := []int{}
	var bt func(int)
	bt = func(i int) {
		if i == len(nums) {
			subsets = append(subsets, slices.Clone(cur)) // COPY! cur keeps changing
			return
		}
		cur = append(cur, nums[i])
		bt(i + 1)
		cur = cur[:len(cur)-1]
		bt(i + 1)
	}
	bt(0)
	fmt.Println("subsets:", subsets)
}

// =============================================================================
// 17. BIT MANIPULATION
// =============================================================================
func sectionBits() {
	x := 0b10110                                           // 22
	fmt.Println("and/or/xor/andnot:", 6&3, 6|3, 6^3, 6&^3) // 2 7 5 4  (&^ = AND NOT)
	fmt.Println("shifts:", 1<<4, 32>>2)
	fmt.Println("popcount:", bits.OnesCount(uint(x)))           // 3
	fmt.Println("bit length:", bits.Len(uint(x)))               // 5
	fmt.Println("trailing zeros:", bits.TrailingZeros(uint(x))) // 1
	fmt.Println("binary:", strconv.FormatInt(int64(x), 2), fmt.Sprintf("%08b", x))
	i := 2
	fmt.Println("check bit i:", (x>>i)&1, "set:", x|(1<<i), "clear:", x&^(1<<i), "toggle:", x^(1<<i))
	fmt.Println("lowest set bit:", x&-x, "drop lowest:", x&(x-1), "power of 2?:", 16&(16-1) == 0)

	// iterate over all subsets of n items
	n := 3
	cnt := 0
	for mask := 0; mask < 1<<n; mask++ {
		cnt++
	}
	// iterate over submasks of a mask
	mask, subs := 0b101, []int{}
	for s := mask; s > 0; s = (s - 1) & mask {
		subs = append(subs, s)
	}
	fmt.Println("all masks:", cnt, "submasks of 101:", subs)

	// xor trick: single number
	single := 0
	for _, v := range []int{4, 1, 2, 1, 2} {
		single ^= v
	}
	fmt.Println("single number:", single)
}

// =============================================================================
// 18. COMMON PATTERNS
// =============================================================================
func sectionPatterns() {
	nums := []int{2, 1, 3, 4, 1}

	// prefix sum: sum(l..r) = pre[r+1] - pre[l]
	pre := make([]int, len(nums)+1)
	for i, v := range nums {
		pre[i+1] = pre[i] + v
	}
	fmt.Println("sum[1..3]:", pre[4]-pre[1])

	// two pointers on a sorted array: pair with target sum
	s := []int{1, 2, 4, 7, 11}
	l, r, target := 0, len(s)-1, 15
	for l < r {
		sum := s[l] + s[r]
		if sum == target {
			break
		} else if sum < target {
			l++
		} else {
			r--
		}
	}
	fmt.Println("two-sum pair:", s[l], s[r])

	// sliding window: longest substring without repeating chars
	str := "abcabcbb"
	last := map[byte]int{}
	best, left := 0, 0
	for right := 0; right < len(str); right++ {
		if p, ok := last[str[right]]; ok && p >= left {
			left = p + 1
		}
		last[str[right]] = right
		best = max(best, right-left+1)
	}
	fmt.Println("longest unique window:", best)

	// monotonic stack: next greater element
	arr := []int{2, 1, 2, 4, 3}
	nge := make([]int, len(arr))
	stk := []int{} // indices, values decreasing
	for i := len(arr) - 1; i >= 0; i-- {
		for len(stk) > 0 && arr[stk[len(stk)-1]] <= arr[i] {
			stk = stk[:len(stk)-1]
		}
		if len(stk) == 0 {
			nge[i] = -1
		} else {
			nge[i] = arr[stk[len(stk)-1]]
		}
		stk = append(stk, i)
	}
	fmt.Println("next greater:", nge)

	// Kadane: max subarray sum
	bestSum, curSum := math.MinInt, 0
	for _, v := range []int{-2, 1, -3, 4, -1, 2, 1, -5, 4} {
		curSum = max(v, curSum+v)
		bestSum = max(bestSum, curSum)
	}
	fmt.Println("max subarray:", bestSum)
}

// =============================================================================
// 19. GO GOTCHAS
// =============================================================================
//  * No ternary operator, no while (use `for cond {}`), no do-while. `for {}` = infinite loop.
//  * Unused variables / imports are COMPILE ERRORS. Use `_ = x` while debugging.
//  * Map iteration order is random. Sort keys when order matters.
//  * Slicing shares the underlying array; append may or may not reallocate. Clone when storing
//    a slice you'll keep modifying (backtracking: append(slices.Clone(cur))).
//  * `for _, v := range s` gives a COPY of v; modify via s[i].
//  * string(65) is "A"; use strconv.Itoa for numbers.  s[i] is a byte, not a string.
//  * Recursive closures must be declared first: `var dfs func(int) int; dfs = func(...) ...`
//  * Integer division truncates toward zero; % can be negative.
//  * int is 64-bit on LeetCode, but intermediate products can still overflow — mod early.
//  * min/max/clear builtins need Go 1.21+; `for i := range n` and cmp.Or need Go 1.22+.
//  * Global variables persist across LeetCode test cases — reset them inside the function.

func main() {
	sections := []struct {
		name string
		fn   func()
	}{
		{"1. Numbers", sectionNumbers},
		{"2. Strings", sectionStrings},
		{"3. Arrays", sectionArrays},
		{"4. Sorting", sectionSorting},
		{"5. Binary search", sectionBinarySearch},
		{"6. Hash map", sectionHashMap},
		{"7. Hash set", sectionHashSet},
		{"8. Ordered set", sectionOrdered},
		{"9. Stack", sectionStack},
		{"10. Queue", sectionQueue},
		{"11. Deque", sectionDeque},
		{"12. Heap", sectionHeap},
		{"13. Pairs", sectionPairs},
		{"14. Nodes", sectionNodes},
		{"15. Graphs", sectionGraphs},
		{"16. DP", sectionDP},
		{"17. Bits", sectionBits},
		{"18. Patterns", sectionPatterns},
	}
	for _, s := range sections {
		fmt.Println("\n=====", s.name, "=====")
		s.fn()
	}
	_ = fastIOTemplate // not called: it reads stdin
}
