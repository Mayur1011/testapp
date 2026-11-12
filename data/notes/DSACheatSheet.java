// =============================================================================
//  DSA / LeetCode CHEAT SHEET — JAVA  (Java 17+; sections match the Go/C++/Py files)
//  Run:  javac DSACheatSheet.java && java DSACheatSheet
//  (Java requires the file name to match the public class name.)
// =============================================================================
//
//  SECTIONS
//   0. Template & fast I/O              10. Queue
//   1. Numbers, limits, math            11. Deque
//   2. Strings & chars                  12. Heap / priority queue
//   3. Arrays & ArrayList               13. Pairs, tuples, records
//   4. Sorting                          14. Linked list & binary tree (LeetCode nodes)
//   5. Binary search                    15. Graphs: BFS, DFS, grid, Dijkstra, topo, DSU
//   6. Hash map                         16. Recursion, memoization, DP tables
//   7. Hash set                         17. Bit manipulation
//   8. Ordered map / set (TreeMap)      18. Common patterns (prefix sum, 2 ptr, window, mono stack)
//   9. Stack                            19. Java gotchas
// =============================================================================

import java.io.*;
import java.util.*;
import java.util.stream.*;

public class DSACheatSheet {

    static final int MOD = 1_000_000_007;

    // =========================================================================
    // 0. TEMPLATE & FAST I/O
    // =========================================================================
    // On LeetCode you only write the method inside `class Solution { ... }`.
    // For competitive programming (Scanner is slow):
    static void fastIOTemplate() throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        PrintWriter out = new PrintWriter(new BufferedWriter(new OutputStreamWriter(System.out)));
        int n = Integer.parseInt(br.readLine().trim());
        StringTokenizer st = new StringTokenizer(br.readLine());
        int[] a = new int[n];
        for (int i = 0; i < n; i++) a[i] = Integer.parseInt(st.nextToken());
        out.println(n + " " + Arrays.toString(a));
        out.flush(); // MUST flush
    }

    // =========================================================================
    // 1. NUMBERS, LIMITS, MATH
    // =========================================================================
    static long powMod(long base, long exp, long mod) {
        long res = 1;
        base %= mod;
        while (exp > 0) {
            if ((exp & 1) == 1) res = res * base % mod;
            base = base * base % mod;
            exp >>= 1;
        }
        return res;
    }

    static int gcd(int a, int b) { return b == 0 ? a : gcd(b, a % b); }

    static void sectionNumbers() {
        int a = 10;                     // 32-bit
        long b = (long) a * 300_000_000; // cast BEFORE multiplying
        double f = a / 4.0;             // 2.5  (a / 4 == 2)
        System.out.println("types: " + a + " " + b + " " + f + " " + (int) f);
        System.out.println("limits: " + Integer.MAX_VALUE + " " + Integer.MIN_VALUE + " " + Long.MAX_VALUE);
        System.out.println("overflow: " + (Integer.MAX_VALUE + 1)); // wraps silently!

        System.out.println("div/mod: " + 7 / 2 + " " + -7 / 2 + " " + 7 % 3 + " " + -7 % 3); // 3 -3 1 -1
        System.out.println("floorMod/floorDiv: " + Math.floorMod(-7, 3) + " " + Math.floorDiv(-7, 2)); // 2 -4

        System.out.println("min/max/abs: " + Math.min(3, 1) + " " + Math.max(3, 1) + " " + Math.abs(-5));
        System.out.println("sqrt/pow: " + Math.sqrt(16) + " " + (long) Math.pow(2, 10) + " " + (int) Math.sqrt(17));
        System.out.println("floor/ceil: " + Math.floor(2.7) + " " + Math.ceil(2.1));
        System.out.println("ceil div 7/2: " + (7 + 2 - 1) / 2);
        System.out.println("gcd/lcm: " + gcd(12, 18) + " " + (4 / gcd(4, 6) * 6));
        System.out.println("powMod: " + powMod(2, 30, MOD));
        long x = 123456789, y = 987654321;
        System.out.println("mod mult: " + x % MOD * (y % MOD) % MOD);
        System.out.println("addExact throws on overflow: " + Math.addExact(1, 2));
    }

    // =========================================================================
    // 2. STRINGS & CHARS  (String is IMMUTABLE; use StringBuilder to build)
    // =========================================================================
    static void sectionStrings() {
        String s = "hello world";
        System.out.println("len: " + s.length() + " charAt(0): " + s.charAt(0) + " sub: " + s.substring(0, 5)); // [begin, end)

        for (char c : s.toCharArray()) { }         // iterate chars
        for (int i = 0; i < s.length(); i++) { char c = s.charAt(i); }

        char c = 'c';
        System.out.println("c-'a': " + (c - 'a') + " 'a'+3: " + (char) ('a' + 3) + " '7'-'0': " + ('7' - '0'));
        System.out.println("isDigit/isLetter: " + Character.isDigit('5') + " " + Character.isLetter(c)
                + " " + Character.isLetterOrDigit('_'));
        System.out.println("upper/lower: " + Character.toUpperCase('q') + " " + "abc".toUpperCase() + " " + "ABC".toLowerCase());

        char[] arr = s.toCharArray();               // modify: char[] then new String
        arr[0] = 'H';
        System.out.println("modified: " + new String(arr) + " / " + String.valueOf(arr));

        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < 3; i++) sb.append("ab").append('-');
        sb.setCharAt(0, 'A');
        sb.deleteCharAt(sb.length() - 1);   // remove last
        sb.insert(0, ">");
        System.out.println("builder: " + sb + " reversed: " + new StringBuilder("abc").reverse());

        System.out.println(Arrays.toString("a,b,c".split(",")) + " " + Arrays.toString("  a  b c ".trim().split("\\s+"))); // split takes a REGEX
        System.out.println(String.join("/", List.of("x", "y")) + " " + "  hi  ".trim() + " " + "  hi ".strip());
        System.out.println(s.contains("wor") + " " + s.indexOf('o') + " " + s.lastIndexOf("o") + " " + s.indexOf("z"));
        System.out.println(s.startsWith("he") + " " + s.endsWith("ld") + " " + "aXbXc".replace("X", "-") + " " + "ab".repeat(3));

        System.out.println("parse: " + (Integer.parseInt("123") + 1) + " " + Long.parseLong("9000000000")
                + " toString: " + String.valueOf(456) + "!" + " " + Integer.toString(7));
        System.out.println("binary: " + Integer.toBinaryString(10) + " parse base2: " + Integer.parseInt("1010", 2)
                + " format: " + String.format("%d-%s-%.2f", 1, "x", 3.14159));

        // == compares REFERENCES. Always use equals() for content.
        String t1 = new String("hi"), t2 = "hi";
        System.out.println("== vs equals: " + (t1 == t2) + " " + t1.equals(t2) + " compareTo: " + "apple".compareTo("banana"));
        System.out.println("equalsIgnoreCase: " + "Hi".equalsIgnoreCase("hI") + " isEmpty: " + "".isEmpty());

        int[] freq = new int[26];
        for (char ch : "hello".toCharArray()) freq[ch - 'a']++;
        System.out.println("freq of 'l': " + freq['l' - 'a']);

        char[] sortedChars = "dcba".toCharArray();
        Arrays.sort(sortedChars);
        System.out.println("sorted chars: " + new String(sortedChars));
    }

    // =========================================================================
    // 3. ARRAYS & ARRAYLIST
    // =========================================================================
    static void sectionArrays() {
        int[] a = {5, 3, 8, 1};
        int[] b = new int[5];             // zeros
        Arrays.fill(b, -1);
        int[][] grid = new int[3][4];     // 2D, zeros
        grid[1][2] = 7;
        for (int[] row : grid) Arrays.fill(row, 0);
        System.out.println("arrays: " + Arrays.toString(a) + " " + Arrays.toString(b) + " " + Arrays.deepToString(grid));

        int[] cp = a.clone();                            // copy
        int[] part = Arrays.copyOfRange(a, 1, 3);        // [1, 3)
        int[] bigger = Arrays.copyOf(a, 6);              // pad with 0
        System.out.println("copies: " + Arrays.toString(cp) + Arrays.toString(part) + Arrays.toString(bigger));
        System.out.println("equals: " + Arrays.equals(a, cp) + " sum: " + Arrays.stream(a).sum()
                + " max: " + Arrays.stream(a).max().getAsInt());

        // ArrayList = dynamic array (holds objects -> Integer, not int)
        List<Integer> list = new ArrayList<>();
        list.add(5); list.add(3); list.add(8);
        list.add(0, 99);                 // insert at index O(n)
        list.set(0, 100);                // update
        int first = list.get(0);
        list.remove(0);                  // remove by INDEX
        list.remove(Integer.valueOf(8)); // remove by VALUE
        int last = list.get(list.size() - 1);
        list.remove(list.size() - 1);    // pop back
        System.out.println("list: " + list + " first: " + first + " last: " + last);

        List<Integer> l2 = new ArrayList<>(List.of(4, 2, 9, 2)); // List.of is IMMUTABLE -> wrap it
        System.out.println("contains/indexOf: " + l2.contains(9) + " " + l2.indexOf(2) + " max: " + Collections.max(l2)
                + " freq(2): " + Collections.frequency(l2, 2));
        Collections.reverse(l2);
        System.out.println("reversed: " + l2 + " subList: " + l2.subList(1, 3));

        // conversions
        int[] prim = l2.stream().mapToInt(Integer::intValue).toArray();
        List<Integer> boxed = Arrays.stream(prim).boxed().collect(Collectors.toList());
        System.out.println("convert: " + Arrays.toString(prim) + " " + boxed);

        List<List<Integer>> nested = new ArrayList<>();
        for (int i = 0; i < 3; i++) nested.add(new ArrayList<>());
        nested.get(1).add(42);
        System.out.println("list of lists: " + nested);

        List<Integer> uniq = new ArrayList<>(new TreeSet<>(List.of(3, 1, 3, 2, 1)));
        System.out.println("unique sorted: " + uniq);
    }

    // =========================================================================
    // 4. SORTING
    // =========================================================================
    static void sectionSorting() {
        int[] a = {5, 2, 9, 1};
        Arrays.sort(a); // primitives: ascending only (dual-pivot quicksort)
        System.out.println("asc: " + Arrays.toString(a));

        Integer[] boxed = {5, 2, 9, 1};                 // descending needs boxed type
        Arrays.sort(boxed, Collections.reverseOrder());
        System.out.println("desc: " + Arrays.toString(boxed));

        List<Integer> list = new ArrayList<>(List.of(5, 2, 9, 1));
        Collections.sort(list);
        list.sort((x, y) -> y - x);  // careful: y - x can overflow for big values
        list.sort(Comparator.reverseOrder());
        System.out.println("list desc: " + list);

        int[][] iv = {{1, 5}, {0, 2}, {1, 9}};
        Arrays.sort(iv, (p, q) -> p[0] != q[0] ? Integer.compare(p[0], q[0]) : Integer.compare(q[1], p[1]));
        System.out.println("intervals: " + Arrays.deepToString(iv));

        record Person(String name, int age) {}
        List<Person> ps = new ArrayList<>(List.of(new Person("bob", 25), new Person("amy", 25), new Person("cat", 20)));
        ps.sort(Comparator.comparingInt(Person::age).thenComparing(Person::name));
        System.out.println("by age then name: " + ps);
        ps.sort(Comparator.comparingInt(Person::age).reversed()); // List.sort is STABLE
        System.out.println("by age desc (stable): " + ps);

        String[] strs = {"pear", "apple", "fig"};
        Arrays.sort(strs);
        Arrays.sort(strs, (s1, s2) -> s1.length() - s2.length());
        System.out.println("by length: " + Arrays.toString(strs));

        int[] vals = {30, 10, 20};
        Integer[] idx = {0, 1, 2};
        Arrays.sort(idx, (i, j) -> Integer.compare(vals[i], vals[j]));
        System.out.println("argsort: " + Arrays.toString(idx));
    }

    // =========================================================================
    // 5. BINARY SEARCH
    // =========================================================================
    static int lowerBound(int[] a, int x) { // first index with a[i] >= x
        int lo = 0, hi = a.length;
        while (lo < hi) {
            int mid = lo + (hi - lo) / 2;
            if (a[mid] >= x) hi = mid; else lo = mid + 1;
        }
        return lo;
    }

    static int upperBound(int[] a, int x) { // first index with a[i] > x
        int lo = 0, hi = a.length;
        while (lo < hi) {
            int mid = lo + (hi - lo) / 2;
            if (a[mid] > x) hi = mid; else lo = mid + 1;
        }
        return lo;
    }

    static void sectionBinarySearch() {
        int[] a = {1, 3, 3, 5, 8};
        int lb = lowerBound(a, 3), ub = upperBound(a, 3);
        System.out.println("lb/ub/count: " + lb + " " + ub + " " + (ub - lb));
        // Arrays.binarySearch: index if found (ANY copy for duplicates), else -(insertionPoint) - 1
        System.out.println("Arrays.binarySearch(5): " + Arrays.binarySearch(a, 5) + " (4): " + Arrays.binarySearch(a, 4));
        int r = Arrays.binarySearch(a, 4);
        System.out.println("insertion point of 4: " + (r < 0 ? -r - 1 : r));

        int l = 0, h = 50;
        while (l < h) {
            int m = (l + h) / 2;
            if (m * m >= 50) h = m; else l = m + 1;
        }
        System.out.println("ceil sqrt(50): " + l);
    }

    // =========================================================================
    // 6. HASH MAP
    // =========================================================================
    static void sectionHashMap() {
        Map<String, Integer> m = new HashMap<>();
        m.put("apple", 3);
        m.put("banana", m.getOrDefault("banana", 0) + 1);
        m.merge("banana", 1, Integer::sum);   // counting shortcut
        System.out.println("get: " + m.get("apple") + " " + m.get("cherry") + " " + m.getOrDefault("cherry", 0)
                + " containsKey: " + m.containsKey("apple") + " size: " + m.size());
        m.remove("apple");
        m.putIfAbsent("kiwi", 5);

        for (Map.Entry<String, Integer> e : m.entrySet()) System.out.println("kv: " + e.getKey() + " " + e.getValue());
        for (String k : m.keySet()) { }
        for (int v : m.values()) { }

        Map<Integer, Integer> freq = new HashMap<>();
        for (int x : new int[]{1, 2, 2, 3, 3, 3}) freq.merge(x, 1, Integer::sum);
        System.out.println("freq: " + freq);

        Map<Integer, List<Integer>> g = new HashMap<>();
        g.computeIfAbsent(1, k -> new ArrayList<>()).add(2);  // map of lists
        g.computeIfAbsent(1, k -> new ArrayList<>()).add(3);
        System.out.println("map of lists: " + g);

        // composite keys: String, List<Integer>, a record, or encode as long (int[] does NOT work as a key!)
        Map<String, Integer> sk = new HashMap<>();
        sk.put(1 + "," + 2, 5);
        Map<List<Integer>, Integer> lk = new HashMap<>();
        lk.put(List.of(1, 2), 5);
        Map<Long, Integer> enc = new HashMap<>();
        enc.put(1L * 100000 + 2, 5);
        System.out.println("composite keys: " + sk.get("1,2") + " " + lk.get(List.of(1, 2)) + " " + enc.get(100002L));

        // group anagrams: key = sorted chars
        Map<String, List<String>> groups = new HashMap<>();
        for (String w : new String[]{"eat", "tea", "tan"}) {
            char[] ch = w.toCharArray();
            Arrays.sort(ch);
            groups.computeIfAbsent(new String(ch), k -> new ArrayList<>()).add(w);
        }
        System.out.println("groups: " + groups.values());

        // LinkedHashMap keeps insertion order (LRU cache with accessOrder=true + removeEldestEntry)
        Map<String, Integer> lhm = new LinkedHashMap<>();
        lhm.put("b", 2); lhm.put("a", 1);
        System.out.println("LinkedHashMap: " + lhm);
    }

    // =========================================================================
    // 7. HASH SET
    // =========================================================================
    static void sectionHashSet() {
        Set<Integer> st = new HashSet<>(List.of(5, 7));
        st.add(9);
        boolean added = st.add(9); // false if already present (handy for "seen" checks)
        System.out.println("contains 5: " + st.contains(5) + " size: " + st.size() + " added again: " + added);
        st.remove(5);

        Set<Integer> A = new HashSet<>(List.of(1, 2, 3)), B = new HashSet<>(List.of(2, 3, 4));
        Set<Integer> inter = new HashSet<>(A); inter.retainAll(B);
        Set<Integer> union = new HashSet<>(A); union.addAll(B);
        Set<Integer> diff = new HashSet<>(A); diff.removeAll(B);
        System.out.println("inter/union/diff: " + inter + " " + union + " " + diff);
    }

    // =========================================================================
    // 8. ORDERED MAP / SET  (TreeMap / TreeSet: red-black tree, O(log n))
    // =========================================================================
    static void sectionOrdered() {
        TreeSet<Integer> ts = new TreeSet<>(List.of(5, 1, 9, 3));
        ts.remove(9);
        System.out.println("TreeSet: " + ts + " first/last: " + ts.first() + " " + ts.last());
        System.out.println("floor(4)<=: " + ts.floor(4) + " ceiling(4)>=: " + ts.ceiling(4)
                + " lower(3)<: " + ts.lower(3) + " higher(3)>: " + ts.higher(3) + " floor(0): " + ts.floor(0)); // null if none
        System.out.println("headSet(<5): " + ts.headSet(5) + " tailSet(>=3): " + ts.tailSet(3) + " pollFirst: " + ts.pollFirst());

        TreeMap<String, Integer> tm = new TreeMap<>(Map.of("b", 2, "a", 1, "c", 3));
        System.out.println("TreeMap: " + tm + " firstKey: " + tm.firstKey() + " lastEntry: " + tm.lastEntry()
                + " ceilingKey(bb): " + tm.ceilingKey("bb") + " floorKey(bb): " + tm.floorKey("bb"));
        System.out.println("descending: " + tm.descendingMap());

        // multiset: TreeMap<value, count>
        TreeMap<Integer, Integer> ms = new TreeMap<>();
        for (int x : new int[]{2, 2, 3}) ms.merge(x, 1, Integer::sum);
        if (ms.merge(2, -1, Integer::sum) == 0) ms.remove(2); // remove one copy
        System.out.println("multiset: " + ms);
    }

    // =========================================================================
    // 9. STACK  (use ArrayDeque, not the legacy Stack class)
    // =========================================================================
    static boolean isValid(String s) {
        Map<Character, Character> pair = Map.of(')', '(', ']', '[', '}', '{');
        Deque<Character> stk = new ArrayDeque<>();
        for (char ch : s.toCharArray()) {
            if (pair.containsKey(ch)) {
                if (stk.isEmpty() || !stk.peek().equals(pair.get(ch))) return false; // Character objects: equals, not !=
                stk.pop();
            } else stk.push(ch);
        }
        return stk.isEmpty();
    }

    static void sectionStack() {
        Deque<Integer> st = new ArrayDeque<>();
        st.push(1); st.push(2); st.push(3);
        int top = st.peek();
        st.pop();
        System.out.println("stack top: " + top + " now: " + st + " empty: " + st.isEmpty());
        System.out.println("valid parens: " + isValid("({[]})") + " " + isValid("(]"));
    }

    // =========================================================================
    // 10. QUEUE
    // =========================================================================
    static void sectionQueue() {
        Queue<Integer> q = new ArrayDeque<>(); // or new LinkedList<>()
        q.offer(1); q.offer(2); q.offer(3);    // add
        int front = q.peek();                  // null if empty (element() throws)
        q.poll();                              // remove, null if empty (remove() throws)
        System.out.println("queue front: " + front + " now: " + q);

        q = new ArrayDeque<>(List.of(1));
        for (int level = 0; level < 3; level++) {
            for (int sz = q.size(); sz > 0; sz--) {
                int x = q.poll();
                q.offer(2 * x); q.offer(2 * x + 1);
            }
        }
        System.out.println("nodes at level 3: " + q);
    }

    // =========================================================================
    // 11. DEQUE
    // =========================================================================
    static void sectionDeque() {
        Deque<Integer> dq = new ArrayDeque<>(); // ArrayDeque does NOT allow null
        dq.offerLast(2); dq.offerLast(3); dq.offerFirst(1);
        System.out.println("front/back: " + dq.peekFirst() + " " + dq.peekLast() + " size: " + dq.size());
        dq.pollFirst(); dq.pollLast();
        System.out.println("remaining: " + dq);
    }

    // =========================================================================
    // 12. HEAP / PRIORITY QUEUE  (default = MIN heap)
    // =========================================================================
    static void sectionHeap() {
        PriorityQueue<Integer> minh = new PriorityQueue<>();
        PriorityQueue<Integer> maxh = new PriorityQueue<>(Collections.reverseOrder());
        for (int x : new int[]{5, 2, 8, 1}) { minh.offer(x); maxh.offer(x); }
        System.out.println("min top: " + minh.peek() + " max top: " + maxh.peek());
        List<Integer> out = new ArrayList<>();
        while (!minh.isEmpty()) out.add(minh.poll());
        System.out.println("popped in order: " + out); // NOTE: printing a PQ directly is NOT sorted

        // int[] {dist, node}: compare with Integer.compare (a[0] - b[0] can overflow)
        PriorityQueue<int[]> pq = new PriorityQueue<>((x, y) -> x[0] != y[0] ? Integer.compare(x[0], y[0]) : Integer.compare(x[1], y[1]));
        pq.offer(new int[]{5, 1}); pq.offer(new int[]{2, 7}); pq.offer(new int[]{2, 3});
        System.out.println("pair heap poll: " + Arrays.toString(pq.poll()) + " " + Arrays.toString(pq.poll()));

        PriorityQueue<String> byLen = new PriorityQueue<>(Comparator.comparingInt(String::length));
        byLen.offer("ccc"); byLen.offer("a");
        System.out.println("comparator heap: " + byLen.peek());

        int k = 2;
        PriorityQueue<Integer> kh = new PriorityQueue<>();
        for (int x : new int[]{4, 1, 7, 3, 9}) { kh.offer(x); if (kh.size() > k) kh.poll(); }
        System.out.println("kth largest (k=2): " + kh.peek());
    }

    // =========================================================================
    // 13. PAIRS, TUPLES, RECORDS
    // =========================================================================
    record Pair(int first, int second) {} // Java 16+: immutable, equals/hashCode auto -> usable as map key

    static class Item implements Comparable<Item> { // classic class with ordering
        int val, idx;
        String name;
        Item(int v, int i, String n) { val = v; idx = i; name = n; }
        public int compareTo(Item o) { return Integer.compare(val, o.val); }
    }

    static void sectionPairs() {
        int[] p = {3, 4}; // quickest "pair" (NOT usable as a HashMap key)
        Pair rp = new Pair(1, 2);
        Map<Pair, Integer> pm = new HashMap<>();
        pm.put(new Pair(1, 2), 5);
        System.out.println("pair: " + p[0] + " " + p[1] + " record: " + rp + " key lookup: " + pm.get(new Pair(1, 2)));

        List<Item> items = new ArrayList<>(List.of(new Item(2, 0, "y"), new Item(1, 1, "x")));
        Collections.sort(items);
        System.out.println("Comparable sorted first: " + items.get(0).name);

        Map.Entry<String, Integer> e = Map.entry("k", 1); // another built-in pair
        System.out.println("entry: " + e.getKey() + "=" + e.getValue());
    }

    // =========================================================================
    // 14. LINKED LIST & BINARY TREE (LeetCode definitions)
    // =========================================================================
    static class ListNode {
        int val;
        ListNode next;
        ListNode() {}
        ListNode(int val) { this.val = val; }
        ListNode(int val, ListNode next) { this.val = val; this.next = next; }
    }

    static class TreeNode {
        int val;
        TreeNode left, right;
        TreeNode() {}
        TreeNode(int val) { this.val = val; }
    }

    static ListNode buildList(int... vals) {
        ListNode dummy = new ListNode(), cur = dummy; // dummy head trick
        for (int v : vals) { cur.next = new ListNode(v); cur = cur.next; }
        return dummy.next;
    }

    static ListNode reverseList(ListNode head) {
        ListNode prev = null;
        while (head != null) {
            ListNode nxt = head.next;
            head.next = prev;
            prev = head;
            head = nxt;
        }
        return prev;
    }

    static ListNode middleNode(ListNode head) {
        ListNode slow = head, fast = head;
        while (fast != null && fast.next != null) { slow = slow.next; fast = fast.next.next; }
        return slow;
    }

    static TreeNode buildTree(Integer... v) { // LeetCode level order with null
        if (v.length == 0 || v[0] == null) return null;
        TreeNode root = new TreeNode(v[0]);
        Queue<TreeNode> q = new ArrayDeque<>();
        q.offer(root);
        int i = 1;
        while (!q.isEmpty() && i < v.length) {
            TreeNode n = q.poll();
            if (i < v.length && v[i] != null) { n.left = new TreeNode(v[i]); q.offer(n.left); }
            i++;
            if (i < v.length && v[i] != null) { n.right = new TreeNode(v[i]); q.offer(n.right); }
            i++;
        }
        return root;
    }

    static void inorder(TreeNode n, List<Integer> res) {
        if (n == null) return;
        inorder(n.left, res);
        res.add(n.val);
        inorder(n.right, res);
    }

    static int maxDepth(TreeNode n) { return n == null ? 0 : 1 + Math.max(maxDepth(n.left), maxDepth(n.right)); }

    static void sectionNodes() {
        ListNode head = buildList(1, 2, 3, 4, 5);
        System.out.println("middle: " + middleNode(head).val);
        head = reverseList(head);
        StringBuilder sb = new StringBuilder();
        for (ListNode c = head; c != null; c = c.next) sb.append(c.val).append(' ');
        System.out.println("reversed list: " + sb);

        TreeNode root = buildTree(3, 9, 20, null, null, 15, 7);
        List<Integer> res = new ArrayList<>();
        inorder(root, res);
        System.out.println("inorder: " + res + " max depth: " + maxDepth(root));

        List<List<Integer>> levels = new ArrayList<>();
        Queue<TreeNode> q = new ArrayDeque<>();
        q.offer(root);
        while (!q.isEmpty()) {
            List<Integer> lvl = new ArrayList<>();
            for (int sz = q.size(); sz > 0; sz--) {
                TreeNode n = q.poll();
                lvl.add(n.val);
                if (n.left != null) q.offer(n.left);
                if (n.right != null) q.offer(n.right);
            }
            levels.add(lvl);
        }
        System.out.println("level order: " + levels);

        List<Integer> it = new ArrayList<>();
        Deque<TreeNode> stk = new ArrayDeque<>();
        TreeNode cur = root;
        while (cur != null || !stk.isEmpty()) {
            while (cur != null) { stk.push(cur); cur = cur.left; }
            cur = stk.pop();
            it.add(cur.val);
            cur = cur.right;
        }
        System.out.println("iterative inorder: " + it);
    }

    // =========================================================================
    // 15. GRAPHS
    // =========================================================================
    static class DSU {
        int[] parent, size;
        DSU(int n) {
            parent = new int[n];
            size = new int[n];
            for (int i = 0; i < n; i++) { parent[i] = i; size[i] = 1; }
        }
        int find(int x) { return parent[x] == x ? x : (parent[x] = find(parent[x])); }
        boolean union(int a, int b) {
            a = find(a); b = find(b);
            if (a == b) return false;
            if (size[a] < size[b]) { int t = a; a = b; b = t; }
            parent[b] = a;
            size[a] += size[b];
            return true;
        }
    }

    static void dfs(int u, List<List<Integer>> adj, boolean[] vis, List<Integer> order) {
        vis[u] = true;
        order.add(u);
        for (int v : adj.get(u)) if (!vis[v]) dfs(v, adj, vis, order);
    }

    static void sectionGraphs() {
        int n = 5;
        int[][] edges = {{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}};
        List<List<Integer>> adj = new ArrayList<>();
        for (int i = 0; i < n; i++) adj.add(new ArrayList<>());
        for (int[] e : edges) { adj.get(e[0]).add(e[1]); adj.get(e[1]).add(e[0]); }

        int[] dist = new int[n];
        Arrays.fill(dist, -1);
        dist[0] = 0;
        Queue<Integer> q = new ArrayDeque<>(List.of(0));
        while (!q.isEmpty()) {
            int u = q.poll();
            for (int v : adj.get(u)) if (dist[v] == -1) { dist[v] = dist[u] + 1; q.offer(v); }
        }
        System.out.println("BFS dist: " + Arrays.toString(dist));

        List<Integer> order = new ArrayList<>();
        dfs(0, adj, new boolean[n], order);
        System.out.println("DFS order: " + order);

        String[] grid = {"110", "010", "001"};
        int R = grid.length, C = grid[0].length();
        int[][] dirs = {{1, 0}, {-1, 0}, {0, 1}, {0, -1}};
        boolean[][] seen = new boolean[R][C];
        int islands = 0;
        for (int r = 0; r < R; r++)
            for (int c = 0; c < C; c++) {
                if (grid[r].charAt(c) != '1' || seen[r][c]) continue;
                islands++;
                seen[r][c] = true;
                Queue<int[]> gq = new ArrayDeque<>();
                gq.offer(new int[]{r, c});
                while (!gq.isEmpty()) {
                    int[] cell = gq.poll();
                    for (int[] d : dirs) {
                        int nr = cell[0] + d[0], nc = cell[1] + d[1];
                        if (nr >= 0 && nr < R && nc >= 0 && nc < C && grid[nr].charAt(nc) == '1' && !seen[nr][nc]) {
                            seen[nr][nc] = true;
                            gq.offer(new int[]{nr, nc});
                        }
                    }
                }
            }
        System.out.println("islands: " + islands);

        int[][] wedges = {{0, 1, 4}, {0, 2, 1}, {2, 1, 2}, {1, 3, 1}, {2, 3, 5}};
        List<List<int[]>> wadj = new ArrayList<>();
        for (int i = 0; i < 4; i++) wadj.add(new ArrayList<>());
        for (int[] e : wedges) wadj.get(e[0]).add(new int[]{e[1], e[2]});
        long[] d = new long[4];
        Arrays.fill(d, Long.MAX_VALUE);
        d[0] = 0;
        PriorityQueue<long[]> pq = new PriorityQueue<>((x, y) -> Long.compare(x[0], y[0])); // {dist, node}
        pq.offer(new long[]{0, 0});
        while (!pq.isEmpty()) {
            long[] top = pq.poll();
            long du = top[0];
            int u = (int) top[1];
            if (du > d[u]) continue; // stale
            for (int[] e : wadj.get(u)) {
                int v = e[0], w = e[1];
                if (du + w < d[v]) { d[v] = du + w; pq.offer(new long[]{d[v], v}); }
            }
        }
        System.out.println("Dijkstra: " + Arrays.toString(d));

        int[][] dag = {{0, 1}, {0, 2}, {1, 3}, {2, 3}};
        List<List<Integer>> g = new ArrayList<>();
        for (int i = 0; i < 4; i++) g.add(new ArrayList<>());
        int[] indeg = new int[4];
        for (int[] e : dag) { g.get(e[0]).add(e[1]); indeg[e[1]]++; }
        Queue<Integer> tq = new ArrayDeque<>();
        for (int i = 0; i < 4; i++) if (indeg[i] == 0) tq.offer(i);
        List<Integer> topo = new ArrayList<>();
        while (!tq.isEmpty()) {
            int u = tq.poll();
            topo.add(u);
            for (int v : g.get(u)) if (--indeg[v] == 0) tq.offer(v);
        }
        System.out.println("topo: " + topo + " has cycle: " + (topo.size() != 4));

        DSU dsu = new DSU(5);
        dsu.union(0, 1); dsu.union(3, 4);
        System.out.println("DSU same(0,1): " + (dsu.find(0) == dsu.find(1)) + " same(1,3): " + (dsu.find(1) == dsu.find(3)));
    }

    // =========================================================================
    // 16. RECURSION, MEMOIZATION, DP TABLES
    // =========================================================================
    static Map<Integer, Long> fibMemo = new HashMap<>();

    static long fib(int n) {
        if (n < 2) return n;
        Long cached = fibMemo.get(n);
        if (cached != null) return cached;
        long r = fib(n - 1) + fib(n - 2);
        fibMemo.put(n, r); // computeIfAbsent with recursion throws ConcurrentModificationException!
        return r;
    }

    static int[][] memo;

    static int paths(int i, int j) {
        if (i == 0 || j == 0) return 1;
        if (memo[i][j] != -1) return memo[i][j];
        return memo[i][j] = paths(i - 1, j) + paths(i, j - 1);
    }

    static void backtrack(int i, int[] nums, List<Integer> cur, List<List<Integer>> res) {
        if (i == nums.length) { res.add(new ArrayList<>(cur)); return; } // COPY!
        cur.add(nums[i]);
        backtrack(i + 1, nums, cur, res);
        cur.remove(cur.size() - 1);
        backtrack(i + 1, nums, cur, res);
    }

    static void sectionDP() {
        System.out.println("fib(50): " + fib(50));

        memo = new int[3][7];
        for (int[] row : memo) Arrays.fill(row, -1);
        System.out.println("unique paths 3x7 (memo): " + paths(2, 6));

        String a = "abcde", b = "ace";
        int[][] dp = new int[a.length() + 1][b.length() + 1];
        for (int i = 1; i <= a.length(); i++)
            for (int j = 1; j <= b.length(); j++)
                dp[i][j] = a.charAt(i - 1) == b.charAt(j - 1) ? dp[i - 1][j - 1] + 1 : Math.max(dp[i - 1][j], dp[i][j - 1]);
        System.out.println("LCS: " + dp[a.length()][b.length()]);

        int[] coins = {1, 2, 5};
        int amount = 11;
        int[] cc = new int[amount + 1];
        Arrays.fill(cc, Integer.MAX_VALUE);
        cc[0] = 0;
        for (int i = 1; i <= amount; i++)
            for (int c : coins)
                if (c <= i && cc[i - c] != Integer.MAX_VALUE) cc[i] = Math.min(cc[i], cc[i - c] + 1);
        System.out.println("coin change: " + cc[amount]);

        List<List<Integer>> subsets = new ArrayList<>();
        backtrack(0, new int[]{1, 2, 3}, new ArrayList<>(), subsets);
        System.out.println("subsets: " + subsets);
    }

    // =========================================================================
    // 17. BIT MANIPULATION
    // =========================================================================
    static void sectionBits() {
        int x = 0b10110; // 22
        System.out.println("and/or/xor/not: " + (6 & 3) + " " + (6 | 3) + " " + (6 ^ 3) + " " + (~0));
        System.out.println("shifts: " + (1 << 4) + " " + (32 >> 2) + " 64-bit: " + (1L << 40) + " unsigned >>>: " + (-1 >>> 28));
        System.out.println("popcount: " + Integer.bitCount(x) + " long: " + Long.bitCount(1L << 40));
        System.out.println("trailing zeros: " + Integer.numberOfTrailingZeros(x) + " leading zeros: " + Integer.numberOfLeadingZeros(x)
                + " bit length: " + (32 - Integer.numberOfLeadingZeros(x)) + " highestOneBit: " + Integer.highestOneBit(x));
        System.out.println("binary: " + Integer.toBinaryString(x) + " padded: " + String.format("%8s", Integer.toBinaryString(x)).replace(' ', '0'));
        int i = 2;
        System.out.println("check bit i: " + ((x >> i) & 1) + " set: " + (x | (1 << i)) + " clear: " + (x & ~(1 << i)) + " toggle: " + (x ^ (1 << i)));
        System.out.println("lowest set bit: " + (x & -x) + " drop lowest: " + (x & (x - 1)) + " pow2?: " + ((16 & 15) == 0));

        int mask = 0b101;
        List<Integer> subs = new ArrayList<>();
        for (int s = mask; s > 0; s = (s - 1) & mask) subs.add(s);
        System.out.println("submasks of 101: " + subs);

        int single = 0;
        for (int v : new int[]{4, 1, 2, 1, 2}) single ^= v;
        System.out.println("single number: " + single);
    }

    // =========================================================================
    // 18. COMMON PATTERNS
    // =========================================================================
    static void sectionPatterns() {
        int[] nums = {2, 1, 3, 4, 1};
        long[] pre = new long[nums.length + 1];
        for (int i = 0; i < nums.length; i++) pre[i + 1] = pre[i] + nums[i];
        System.out.println("sum[1..3]: " + (pre[4] - pre[1]));

        int[] s = {1, 2, 4, 7, 11};
        int l = 0, r = s.length - 1, target = 15;
        while (l < r) {
            int sum = s[l] + s[r];
            if (sum == target) break;
            if (sum < target) l++; else r--;
        }
        System.out.println("two-sum pair: " + s[l] + " " + s[r]);

        String str = "abcabcbb";
        int[] last = new int[128];
        Arrays.fill(last, -1);
        int best = 0, left = 0;
        for (int right = 0; right < str.length(); right++) {
            left = Math.max(left, last[str.charAt(right)] + 1);
            last[str.charAt(right)] = right;
            best = Math.max(best, right - left + 1);
        }
        System.out.println("longest unique window: " + best);

        int[] arr = {2, 1, 2, 4, 3}, nge = new int[arr.length];
        Arrays.fill(nge, -1);
        Deque<Integer> stk = new ArrayDeque<>();
        for (int i = arr.length - 1; i >= 0; i--) {
            while (!stk.isEmpty() && arr[stk.peek()] <= arr[i]) stk.pop();
            if (!stk.isEmpty()) nge[i] = arr[stk.peek()];
            stk.push(i);
        }
        System.out.println("next greater: " + Arrays.toString(nge));

        int bestSum = Integer.MIN_VALUE, cur = 0;
        for (int v : new int[]{-2, 1, -3, 4, -1, 2, 1, -5, 4}) { cur = Math.max(v, cur + v); bestSum = Math.max(bestSum, cur); }
        System.out.println("max subarray: " + bestSum);

        // streams (concise, but slower than loops — fine for small things)
        List<Integer> evensSquared = IntStream.range(0, 6).filter(v -> v % 2 == 0).map(v -> v * v).boxed().collect(Collectors.toList());
        System.out.println("stream: " + evensSquared);
    }

    // =========================================================================
    // 19. JAVA GOTCHAS
    // =========================================================================
    //  * Integer caching: Integer a = 1000, b = 1000; a == b is FALSE. Use .equals() or unbox (intValue()).
    //    Classic bug: map.get(x) == map.get(y) for counts > 127.
    //  * String == compares references; use equals().
    //  * list.remove(int) removes by INDEX; list.remove(Integer.valueOf(x)) removes by VALUE.
    //  * int overflow is silent: use long, cast before multiplying ((long) a * b).
    //  * Comparator (a, b) -> a - b can overflow; use Integer.compare(a, b).
    //  * Arrays.sort(int[]) can't take a comparator; use Integer[] or sort int[][] by rows.
    //  * List.of / Map.of are immutable (and reject nulls). Wrap in new ArrayList<>(...).
    //  * Printing a PriorityQueue does not show sorted order.
    //  * int[] cannot be a HashMap key (uses identity). Use a String, List<Integer>, record or long.
    //  * Use ArrayDeque for stacks/queues (Stack and LinkedList are slower).
    //  * Static fields persist across LeetCode test cases — reset them in the method.
    //  * Deep recursion (~1e4+ frames) may StackOverflow; go iterative for big inputs.
    //  * % of a negative number is negative; use Math.floorMod for a non-negative mod.

    public static void main(String[] args) {
        Object[][] sections = {
            {"1. Numbers", (Runnable) DSACheatSheet::sectionNumbers},
            {"2. Strings", (Runnable) DSACheatSheet::sectionStrings},
            {"3. Arrays", (Runnable) DSACheatSheet::sectionArrays},
            {"4. Sorting", (Runnable) DSACheatSheet::sectionSorting},
            {"5. Binary search", (Runnable) DSACheatSheet::sectionBinarySearch},
            {"6. Hash map", (Runnable) DSACheatSheet::sectionHashMap},
            {"7. Hash set", (Runnable) DSACheatSheet::sectionHashSet},
            {"8. Ordered map/set", (Runnable) DSACheatSheet::sectionOrdered},
            {"9. Stack", (Runnable) DSACheatSheet::sectionStack},
            {"10. Queue", (Runnable) DSACheatSheet::sectionQueue},
            {"11. Deque", (Runnable) DSACheatSheet::sectionDeque},
            {"12. Heap", (Runnable) DSACheatSheet::sectionHeap},
            {"13. Pairs", (Runnable) DSACheatSheet::sectionPairs},
            {"14. Nodes", (Runnable) DSACheatSheet::sectionNodes},
            {"15. Graphs", (Runnable) DSACheatSheet::sectionGraphs},
            {"16. DP", (Runnable) DSACheatSheet::sectionDP},
            {"17. Bits", (Runnable) DSACheatSheet::sectionBits},
            {"18. Patterns", (Runnable) DSACheatSheet::sectionPatterns},
        };
        for (Object[] s : sections) {
            System.out.println("\n===== " + s[0] + " =====");
            ((Runnable) s[1]).run();
        }
    }
}
