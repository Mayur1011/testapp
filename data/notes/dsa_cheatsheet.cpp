// =============================================================================
//  DSA / LeetCode CHEAT SHEET — C++  (C++17; sections match the Go/Py/Java files)
//  Run:  g++ -std=c++17 -O2 dsa_cheatsheet.cpp -o cs && ./cs
// =============================================================================
//
//  SECTIONS
//   0. Template & fast I/O              10. Queue
//   1. Numbers, limits, math            11. Deque
//   2. Strings & chars                  12. Heap / priority queue
//   3. Arrays / vectors                 13. Pairs, tuples, structs
//   4. Sorting                          14. Linked list & binary tree (LeetCode nodes)
//   5. Binary search                    15. Graphs: BFS, DFS, grid, Dijkstra, topo, DSU
//   6. Hash map                         16. Recursion, memoization, DP tables
//   7. Hash set                         17. Bit manipulation
//   8. Ordered map / set                18. Common patterns (prefix sum, 2 ptr, window, mono stack)
//   9. Stack                            19. C++ gotchas
// =============================================================================

#include <bits/stdc++.h> // includes the whole STL (GCC). LeetCode already includes everything.
using namespace std;

using ll = long long;
using pii = pair<int, int>;
const int MOD = 1'000'000'007;

template <class T>
void print(const string& label, const vector<T>& v) {
    cout << label;
    for (auto& x : v) cout << x << ' ';
    cout << '\n';
}

// =============================================================================
// 0. TEMPLATE & FAST I/O
// =============================================================================
// On LeetCode you only write the method inside `class Solution { public: ... };`
// For competitive programming:
void fastIOTemplate() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    int n;
    cin >> n;
    vector<int> a(n);
    for (auto& x : a) cin >> x;
    cout << n << '\n'; // prefer '\n' over endl (endl flushes = slow)
}

// =============================================================================
// 1. NUMBERS, LIMITS, MATH
// =============================================================================
ll powMod(ll base, ll exp, ll mod) {
    ll res = 1;
    base %= mod;
    while (exp > 0) {
        if (exp & 1) res = res * base % mod;
        base = base * base % mod;
        exp >>= 1;
    }
    return res;
}

void sectionNumbers() {
    int a = 10;              // 32-bit: up to ~2.1e9
    ll b = 1LL * a * 300000000; // force 64-bit BEFORE multiplying
    double f = a / 4.0;      // 2.5  (a / 4 would be 2)
    cout << "types: " << a << ' ' << b << ' ' << f << ' ' << (int)f << '\n';

    cout << "limits: " << INT_MAX << ' ' << INT_MIN << ' ' << LLONG_MAX << '\n';
    cout << "numeric_limits: " << numeric_limits<int>::max() << '\n';
    const int INF = 1e9; // common "infinity" that won't overflow when you add two of them

    // division truncates toward zero; % keeps sign of dividend
    cout << "div/mod: " << 7 / 2 << ' ' << -7 / 2 << ' ' << 7 % 3 << ' ' << -7 % 3 << '\n'; // 3 -3 1 -1
    cout << "non-negative mod: " << ((-7 % 3) + 3) % 3 << '\n';

    cout << "min/max: " << min(3, 1) << ' ' << max({3, 1, 2}) << '\n'; // initializer list for >2
    cout << "abs: " << abs(-5) << ' ' << llabs(-5LL) << ' ' << fabs(-2.5) << '\n';
    cout << "sqrt/pow: " << sqrt(16) << ' ' << pow(2, 10) << ' ' << (int)sqrt(17) << '\n';
    cout << "floor/ceil: " << floor(2.7) << ' ' << ceil(2.1) << '\n';
    cout << "ceil div 7/2: " << (7 + 2 - 1) / 2 << '\n';
    cout << "gcd/lcm: " << gcd(12, 18) << ' ' << lcm(4, 6) << '\n'; // <numeric>, C++17
    cout << "powMod: " << powMod(2, 30, MOD) << '\n';
    ll x = 123456789, y = 987654321;
    cout << "mod mult: " << x % MOD * (y % MOD) % MOD << '\n';
    (void)INF;
}

// =============================================================================
// 2. STRINGS & CHARS   (std::string is MUTABLE)
// =============================================================================
void sectionStrings() {
    string s = "hello world";
    cout << "len: " << s.size() << " s[0]: " << s[0] << " sub: " << s.substr(0, 5) << '\n'; // substr(pos, LEN)

    for (char c : s) { (void)c; }          // iterate chars
    for (size_t i = 0; i < s.size(); i++) {} // by index

    char c = 'c';
    cout << "c-'a': " << c - 'a' << " 'a'+3: " << char('a' + 3) << " '7'-'0': " << '7' - '0' << '\n';
    cout << "isdigit/isalpha/isalnum: " << isdigit('5') << ' ' << (bool)isalpha(c) << ' ' << (bool)isalnum('_') << '\n';
    cout << "toupper: " << (char)toupper('q') << '\n';
    string up = "abc";
    transform(up.begin(), up.end(), up.begin(), ::toupper);
    cout << "upper string: " << up << '\n';

    s[0] = 'H'; // mutable
    s += "!";   // append (amortised O(1))
    s.push_back('?');
    s.pop_back();
    cout << "modified: " << s << '\n';

    // find returns string::npos if not found
    cout << "find: " << s.find("o") << ' ' << s.rfind("o") << ' ' << (s.find("z") == string::npos) << '\n';
    s.insert(0, ">> ");
    s.erase(0, 3); // erase(pos, len)
    s.replace(0, 5, "HOWDY");
    cout << "insert/erase/replace: " << s << '\n';

    // number <-> string
    cout << "stoi: " << stoi("123") + 1 << " stoll: " << stoll("9000000000") << " to_string: " << to_string(456) + "!" << '\n';
    cout << "binary: " << bitset<8>(10) << " parse base2: " << stoi("1010", nullptr, 2) << '\n';

    // split by delimiter / whitespace with stringstream
    stringstream ss("a,b,c");
    string tok;
    vector<string> parts;
    while (getline(ss, tok, ',')) parts.push_back(tok);
    print("split: ", parts);
    stringstream ws("  a  b   c ");
    vector<string> words;
    while (ws >> tok) words.push_back(tok);
    print("split ws: ", words);

    // join
    string joined;
    for (size_t i = 0; i < parts.size(); i++) joined += (i ? "/" : "") + parts[i];
    cout << "join: " << joined << '\n';

    string r = "abc";
    reverse(r.begin(), r.end());
    cout << "reversed: " << r << " repeat: " << string(3, 'x') << '\n';
    cout << "compare: " << (string("apple") < "banana") << ' ' << string("a").compare("b") << '\n';

    int freq[26] = {};
    for (char ch : string("hello")) if (islower(ch)) freq[ch - 'a']++;
    cout << "freq of 'l': " << freq['l' - 'a'] << '\n';
}

// =============================================================================
// 3. ARRAYS / VECTORS
// =============================================================================
void sectionArrays() {
    vector<int> a = {5, 3, 8, 1};
    vector<int> b(5);      // 5 zeros
    vector<int> c(5, -1);  // 5 x -1
    vector<int> d;
    d.reserve(10);         // pre-allocate
    print("a: ", a);

    int rows = 3, cols = 4;
    vector<vector<int>> grid(rows, vector<int>(cols, 0)); // 2D
    grid[1][2] = 7;
    cout << "grid[1][2]: " << grid[1][2] << '\n';
    int arr[5] = {};       // C-array, zero-initialised
    fill(begin(arr), end(arr), 3);
    memset(arr, 0, sizeof arr); // memset only for 0 or -1!

    a.push_back(10);            // push back
    int last = a.back();        // back
    a.pop_back();               // pop back
    int first = a.front();
    a.insert(a.begin() + 1, 99); // insert at index 1  O(n)
    a.erase(a.begin() + 1);      // erase index 1      O(n)
    a.erase(a.begin());          // pop front          O(n) -> use deque if frequent
    cout << "last/first: " << last << ' ' << first << '\n';
    print("after ops: ", a);

    vector<int> cp = a;          // COPY (vectors copy by value!)
    cp[0] = 1000;
    cout << "a[0] unchanged: " << a[0] << '\n';
    vector<int> sub(a.begin(), a.begin() + 2); // slice [0,2)
    print("slice: ", sub);

    vector<int> s = {4, 2, 9, 2};
    cout << "sum: " << accumulate(s.begin(), s.end(), 0LL) // 0LL => long long sum
         << " max: " << *max_element(s.begin(), s.end())
         << " argmax: " << max_element(s.begin(), s.end()) - s.begin()
         << " count(2): " << count(s.begin(), s.end(), 2)
         << " find 9 idx: " << find(s.begin(), s.end(), 9) - s.begin() << '\n';
    reverse(s.begin(), s.end());
    print("reversed: ", s);

    vector<int> idx(5);
    iota(idx.begin(), idx.end(), 0); // 0 1 2 3 4
    print("iota: ", idx);

    vector<int> t = {3, 1, 3, 2, 1};
    sort(t.begin(), t.end());
    t.erase(unique(t.begin(), t.end()), t.end()); // remove duplicates
    print("unique: ", t);

    a.assign(3, 7); // reset to 3 x 7
    a.clear();
    cout << "empty: " << a.empty() << '\n';
    (void)b; (void)c;
}

// =============================================================================
// 4. SORTING
// =============================================================================
void sectionSorting() {
    vector<int> a = {5, 2, 9, 1};
    sort(a.begin(), a.end());
    print("asc: ", a);
    sort(a.begin(), a.end(), greater<int>());
    print("desc: ", a);
    sort(a.rbegin(), a.rend()); // also ascending-reversed = descending trick
    sort(a.begin(), a.end(), [](int x, int y) { return x < y; }); // lambda: return true if x goes BEFORE y
    print("lambda: ", a);

    // intervals: start asc, then end desc
    vector<vector<int>> iv = {{1, 5}, {0, 2}, {1, 9}};
    sort(iv.begin(), iv.end(), [](const vector<int>& p, const vector<int>& q) {
        if (p[0] != q[0]) return p[0] < q[0];
        return p[1] > q[1];
    });
    cout << "intervals: ";
    for (auto& v : iv) cout << '[' << v[0] << ',' << v[1] << "] ";
    cout << '\n';

    // pairs sort lexicographically by default
    vector<pair<int, string>> ps = {{25, "bob"}, {25, "amy"}, {20, "cat"}};
    sort(ps.begin(), ps.end());
    cout << "pairs: ";
    for (auto& [age, name] : ps) cout << name << ':' << age << ' ';
    cout << '\n';
    stable_sort(ps.begin(), ps.end(), [](auto& p, auto& q) { return p.first > q.first; });

    string str = "dcba";
    sort(str.begin(), str.end());
    cout << "sorted chars: " << str << '\n';

    vector<int> vals = {30, 10, 20}, idx = {0, 1, 2};
    sort(idx.begin(), idx.end(), [&](int i, int j) { return vals[i] < vals[j]; }); // argsort
    print("argsort: ", idx);

    vector<int> k = {7, 2, 9, 4, 1};
    nth_element(k.begin(), k.begin() + 2, k.end()); // k[2] = 3rd smallest, O(n) avg
    cout << "3rd smallest: " << k[2] << '\n';
    partial_sort(k.begin(), k.begin() + 2, k.end()); // first 2 sorted
}

// =============================================================================
// 5. BINARY SEARCH
// =============================================================================
void sectionBinarySearch() {
    vector<int> a = {1, 3, 3, 5, 8};
    int lb = lower_bound(a.begin(), a.end(), 3) - a.begin(); // first >= 3 -> 1
    int ub = upper_bound(a.begin(), a.end(), 3) - a.begin(); // first >  3 -> 3
    cout << "lb/ub/count: " << lb << ' ' << ub << ' ' << ub - lb
         << " binary_search(5): " << binary_search(a.begin(), a.end(), 5) << '\n';
    // largest element <= x : prev(upper_bound(x)) if it's not begin()

    auto cond = [&](int i) { return a[i] >= 4; }; // F..F T..T
    int lo = 0, hi = a.size();
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (cond(mid)) hi = mid; else lo = mid + 1;
    }
    cout << "first >= 4 at: " << lo << '\n';

    int l = 0, r = 50;
    while (l < r) {
        int m = (l + r) / 2;
        if (m * m >= 50) r = m; else l = m + 1;
    }
    cout << "ceil sqrt(50): " << l << '\n';
}

// =============================================================================
// 6. HASH MAP  (unordered_map: O(1) avg)
// =============================================================================
void sectionHashMap() {
    unordered_map<string, int> m;
    m["apple"] = 3;
    m["banana"]++; // operator[] INSERTS default 0 if missing
    m["banana"]++;

    cout << "get: " << m["apple"] << " count(cherry): " << m.count("cherry") << " size: " << m.size() << '\n';
    if (auto it = m.find("apple"); it != m.end()) cout << "found " << it->first << '=' << it->second << '\n';
    // GOTCHA: `if (m["x"])` inserts "x". Use count()/find() (or contains() in C++20) to test.
    m.erase("apple");

    for (auto& [k, v] : m) cout << "kv: " << k << ' ' << v << '\n'; // order is arbitrary

    unordered_map<int, int> freq;
    for (int x : {1, 2, 2, 3, 3, 3}) freq[x]++;
    cout << "freq[3]: " << freq[3] << '\n';

    unordered_map<int, vector<int>> g;
    g[1].push_back(2);
    g[1].push_back(3);
    cout << "map of lists g[1].size(): " << g[1].size() << '\n';

    // pair keys: use map (ordered) or encode as a long long key, since unordered_map<pair> has no hash
    map<pii, int> pm;
    pm[{1, 2}] = 5;
    unordered_map<ll, int> enc;
    int r = 1, c = 2;
    enc[(ll)r * 100000 + c] = 5;
    cout << "pair keys: " << pm[{1, 2}] << ' ' << enc[100002] << '\n';
    m.clear();
}

// =============================================================================
// 7. HASH SET
// =============================================================================
void sectionHashSet() {
    unordered_set<int> st = {5, 7};
    st.insert(9);
    cout << "contains 5: " << st.count(5) << " contains 4: " << st.count(4) << " size: " << st.size() << '\n';
    st.erase(5);

    vector<int> v = {3, 1, 3, 2};
    unordered_set<int> fromVec(v.begin(), v.end());
    cout << "unique count: " << fromVec.size() << '\n';

    unordered_set<int> A = {1, 2, 3}, B = {2, 3, 4};
    vector<int> inter;
    for (int x : A) if (B.count(x)) inter.push_back(x);
    sort(inter.begin(), inter.end());
    print("intersection: ", inter);
}

// =============================================================================
// 8. ORDERED MAP / SET  (red-black tree: O(log n), sorted)
// =============================================================================
void sectionOrdered() {
    set<int> s = {5, 1, 9, 3};
    s.insert(5); // duplicates ignored
    s.erase(9);
    cout << "set: ";
    for (int x : s) cout << x << ' ';
    cout << "| min: " << *s.begin() << " max: " << *s.rbegin() << '\n';

    auto it = s.lower_bound(4);                   // ceiling: smallest >= 4
    cout << "ceil(4): " << (it != s.end() ? *it : -1);
    auto it2 = s.upper_bound(4);                  // floor: largest <= 4 = prev(upper_bound)
    cout << " floor(4): " << (it2 != s.begin() ? *prev(it2) : -1) << '\n';

    multiset<int> ms = {2, 2, 3};
    ms.erase(ms.find(2)); // erase ONE copy (ms.erase(2) erases ALL copies)
    cout << "multiset count(2): " << ms.count(2) << '\n';

    map<string, int> om = {{"b", 2}, {"a", 1}, {"c", 3}};
    for (auto& [k, v] : om) cout << k << '=' << v << ' '; // sorted by key
    cout << "| first key: " << om.begin()->first << " last key: " << om.rbegin()->first << '\n';
    auto mit = om.lower_bound("bb"); // first key >= "bb"
    cout << "lower_bound(bb): " << mit->first << '\n';
}

// =============================================================================
// 9. STACK
// =============================================================================
void sectionStack() {
    stack<int> st;
    st.push(1); st.push(2); st.push(3);
    int top = st.top(); // peek
    st.pop();           // pop returns void!
    cout << "stack top: " << top << " size: " << st.size() << " empty: " << st.empty() << '\n';

    auto isValid = [](const string& s) {
        unordered_map<char, char> pair = {{')', '('}, {']', '['}, {'}', '{'}};
        stack<char> stk;
        for (char ch : s) {
            if (pair.count(ch)) {
                if (stk.empty() || stk.top() != pair[ch]) return false;
                stk.pop();
            } else stk.push(ch);
        }
        return stk.empty();
    };
    cout << "valid parens: " << isValid("({[]})") << ' ' << isValid("(]") << '\n';
    // a vector also works as a stack: push_back / back / pop_back
}

// =============================================================================
// 10. QUEUE
// =============================================================================
void sectionQueue() {
    queue<int> q;
    q.push(1); q.push(2); q.push(3);
    int front = q.front();
    q.pop();
    cout << "queue front: " << front << " back: " << q.back() << " size: " << q.size() << '\n';

    q = queue<int>();
    q.push(1);
    for (int level = 0; level < 3; level++) {
        for (int sz = q.size(); sz > 0; sz--) {
            int x = q.front(); q.pop();
            q.push(2 * x); q.push(2 * x + 1);
        }
    }
    cout << "nodes at level 3: " << q.size() << " starting at " << q.front() << '\n';
}

// =============================================================================
// 11. DEQUE
// =============================================================================
void sectionDeque() {
    deque<int> dq;
    dq.push_back(2); dq.push_back(3); dq.push_front(1);
    cout << "front/back: " << dq.front() << ' ' << dq.back() << " dq[1]: " << dq[1] << '\n'; // random access
    dq.pop_front(); dq.pop_back();
    cout << "remaining: " << dq.front() << '\n';
}

// =============================================================================
// 12. HEAP / PRIORITY QUEUE   (default = MAX heap!)
// =============================================================================
void sectionHeap() {
    priority_queue<int> maxh;                            // max heap
    priority_queue<int, vector<int>, greater<int>> minh; // min heap
    for (int x : {5, 2, 8, 1}) { maxh.push(x); minh.push(x); }
    cout << "max top: " << maxh.top() << " min top: " << minh.top() << '\n';
    vector<int> out;
    while (!minh.empty()) { out.push_back(minh.top()); minh.pop(); }
    print("popped in order: ", out);

    // min-heap of pairs {dist, node}: pairs compare lexicographically
    priority_queue<pii, vector<pii>, greater<pii>> pq;
    pq.push({5, 1}); pq.push({2, 7}); pq.push({2, 3});
    cout << "pair heap top: " << pq.top().first << ',' << pq.top().second << '\n';

    // custom comparator: comparator returns true if a has LOWER priority than b
    auto cmp = [](const pii& a, const pii& b) { return a.second > b.second; }; // min by .second
    priority_queue<pii, vector<pii>, decltype(cmp)> custom(cmp);
    custom.push({1, 9}); custom.push({2, 4});
    cout << "custom top: " << custom.top().first << '\n';

    // k largest: min heap of size k
    int k = 2;
    priority_queue<int, vector<int>, greater<int>> kh;
    for (int x : {4, 1, 7, 3, 9}) { kh.push(x); if ((int)kh.size() > k) kh.pop(); }
    cout << "kth largest (k=2): " << kh.top() << '\n';

    vector<int> h = {3, 1, 4};
    make_heap(h.begin(), h.end()); // heapify a vector in place (max heap)
}

// =============================================================================
// 13. PAIRS, TUPLES, STRUCTS
// =============================================================================
struct Item {
    int val, idx;
    string name;
    bool operator<(const Item& o) const { return val < o.val; } // enables sort / set / priority_queue
};

void sectionPairs() {
    pii p = {3, 4};
    auto p2 = make_pair(1, 2);
    cout << "pair: " << p.first << ' ' << p.second << ' ' << p2.first << '\n';

    tuple<int, string, double> t = {1, "x", 2.5};
    cout << "tuple: " << get<0>(t) << ' ' << get<1>(t) << '\n';
    auto [a, b, c] = t; // structured bindings (C++17)
    cout << "bindings: " << a << b << c << '\n';

    vector<Item> items = {{2, 0, "y"}, {1, 1, "x"}};
    sort(items.begin(), items.end());
    cout << "struct sorted first: " << items[0].name << '\n';

    int x = 1, y = 2;
    swap(x, y);
    cout << "swap: " << x << ' ' << y << '\n';
}

// =============================================================================
// 14. LINKED LIST & BINARY TREE (LeetCode definitions)
// =============================================================================
struct ListNode {
    int val;
    ListNode* next;
    ListNode(int x = 0, ListNode* n = nullptr) : val(x), next(n) {}
};

struct TreeNode {
    int val;
    TreeNode *left, *right;
    TreeNode(int x = 0) : val(x), left(nullptr), right(nullptr) {}
};

ListNode* buildList(const vector<int>& vals) {
    ListNode dummy; // dummy head trick
    ListNode* cur = &dummy;
    for (int v : vals) { cur->next = new ListNode(v); cur = cur->next; }
    return dummy.next;
}

ListNode* reverseList(ListNode* head) {
    ListNode* prev = nullptr;
    while (head) {
        ListNode* nxt = head->next;
        head->next = prev;
        prev = head;
        head = nxt;
    }
    return prev;
}

ListNode* middleNode(ListNode* head) {
    ListNode *slow = head, *fast = head;
    while (fast && fast->next) { slow = slow->next; fast = fast->next->next; }
    return slow;
}

const int NUL = INT_MIN; // marker for null in level-order input
TreeNode* buildTree(const vector<int>& v) {
    if (v.empty() || v[0] == NUL) return nullptr;
    TreeNode* root = new TreeNode(v[0]);
    queue<TreeNode*> q;
    q.push(root);
    size_t i = 1;
    while (!q.empty() && i < v.size()) {
        TreeNode* n = q.front(); q.pop();
        if (i < v.size() && v[i] != NUL) { n->left = new TreeNode(v[i]); q.push(n->left); }
        i++;
        if (i < v.size() && v[i] != NUL) { n->right = new TreeNode(v[i]); q.push(n->right); }
        i++;
    }
    return root;
}

int maxDepth(TreeNode* n) { return n ? 1 + max(maxDepth(n->left), maxDepth(n->right)) : 0; }

void sectionNodes() {
    ListNode* head = buildList({1, 2, 3, 4, 5});
    cout << "middle: " << middleNode(head)->val << '\n';
    head = reverseList(head);
    cout << "reversed list: ";
    for (ListNode* c = head; c; c = c->next) cout << c->val << ' ';
    cout << '\n';

    TreeNode* root = buildTree({3, 9, 20, NUL, NUL, 15, 7});

    // recursive lambda: std::function (simple) — or pass the lambda to itself (faster)
    vector<int> res;
    function<void(TreeNode*)> inorder = [&](TreeNode* n) {
        if (!n) return;
        inorder(n->left);
        res.push_back(n->val);
        inorder(n->right);
    };
    inorder(root);
    print("inorder: ", res);
    cout << "max depth: " << maxDepth(root) << '\n';

    // level order BFS
    vector<vector<int>> levels;
    queue<TreeNode*> q;
    q.push(root);
    while (!q.empty()) {
        vector<int> lvl;
        for (int sz = q.size(); sz > 0; sz--) {
            TreeNode* n = q.front(); q.pop();
            lvl.push_back(n->val);
            if (n->left) q.push(n->left);
            if (n->right) q.push(n->right);
        }
        levels.push_back(lvl);
    }
    cout << "levels: " << levels.size() << " last level size: " << levels.back().size() << '\n';

    // iterative inorder
    vector<int> it;
    stack<TreeNode*> stk;
    TreeNode* cur = root;
    while (cur || !stk.empty()) {
        while (cur) { stk.push(cur); cur = cur->left; }
        cur = stk.top(); stk.pop();
        it.push_back(cur->val);
        cur = cur->right;
    }
    print("iterative inorder: ", it);
}

// =============================================================================
// 15. GRAPHS
// =============================================================================
struct DSU {
    vector<int> parent, sz;
    DSU(int n) : parent(n), sz(n, 1) { iota(parent.begin(), parent.end(), 0); }
    int find(int x) { return parent[x] == x ? x : parent[x] = find(parent[x]); } // path compression
    bool unite(int a, int b) {
        a = find(a); b = find(b);
        if (a == b) return false;
        if (sz[a] < sz[b]) swap(a, b);
        parent[b] = a;
        sz[a] += sz[b];
        return true;
    }
};

void sectionGraphs() {
    int n = 5;
    vector<vector<int>> edges = {{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}};
    vector<vector<int>> adj(n);
    for (auto& e : edges) { adj[e[0]].push_back(e[1]); adj[e[1]].push_back(e[0]); }

    // BFS
    vector<int> dist(n, -1);
    queue<int> q;
    dist[0] = 0; q.push(0);
    while (!q.empty()) {
        int u = q.front(); q.pop();
        for (int v : adj[u]) if (dist[v] == -1) { dist[v] = dist[u] + 1; q.push(v); }
    }
    print("BFS dist: ", dist);

    // DFS
    vector<bool> vis(n, false);
    vector<int> order;
    function<void(int)> dfs = [&](int u) {
        vis[u] = true;
        order.push_back(u);
        for (int v : adj[u]) if (!vis[v]) dfs(v);
    };
    dfs(0);
    print("DFS order: ", order);

    // grid BFS: islands
    vector<string> grid = {"110", "010", "001"};
    int R = grid.size(), C = grid[0].size();
    int dr[] = {1, -1, 0, 0}, dc[] = {0, 0, 1, -1};
    vector<vector<bool>> seen(R, vector<bool>(C, false));
    int islands = 0;
    for (int r = 0; r < R; r++)
        for (int c = 0; c < C; c++) {
            if (grid[r][c] != '1' || seen[r][c]) continue;
            islands++;
            queue<pii> gq;
            gq.push({r, c});
            seen[r][c] = true;
            while (!gq.empty()) {
                auto [cr, cc] = gq.front(); gq.pop();
                for (int d = 0; d < 4; d++) {
                    int nr = cr + dr[d], nc = cc + dc[d];
                    if (nr >= 0 && nr < R && nc >= 0 && nc < C && grid[nr][nc] == '1' && !seen[nr][nc]) {
                        seen[nr][nc] = true;
                        gq.push({nr, nc});
                    }
                }
            }
        }
    cout << "islands: " << islands << '\n';

    // Dijkstra
    vector<vector<int>> wedges = {{0, 1, 4}, {0, 2, 1}, {2, 1, 2}, {1, 3, 1}, {2, 3, 5}};
    vector<vector<pii>> wadj(4); // {v, w}
    for (auto& e : wedges) wadj[e[0]].push_back({e[1], e[2]});
    vector<ll> d(4, LLONG_MAX);
    priority_queue<pair<ll, int>, vector<pair<ll, int>>, greater<>> pq;
    d[0] = 0; pq.push({0, 0});
    while (!pq.empty()) {
        auto [du, u] = pq.top(); pq.pop();
        if (du > d[u]) continue; // stale
        for (auto [v, w] : wadj[u])
            if (du + w < d[v]) { d[v] = du + w; pq.push({d[v], v}); }
    }
    print("Dijkstra: ", d);

    // Topological sort (Kahn)
    vector<vector<int>> dag = {{0, 1}, {0, 2}, {1, 3}, {2, 3}};
    vector<vector<int>> g(4);
    vector<int> indeg(4, 0);
    for (auto& e : dag) { g[e[0]].push_back(e[1]); indeg[e[1]]++; }
    queue<int> tq;
    for (int i = 0; i < 4; i++) if (indeg[i] == 0) tq.push(i);
    vector<int> topo;
    while (!tq.empty()) {
        int u = tq.front(); tq.pop();
        topo.push_back(u);
        for (int v : g[u]) if (--indeg[v] == 0) tq.push(v);
    }
    print("topo: ", topo);
    cout << "has cycle: " << (topo.size() != 4) << '\n';

    DSU dsu(5);
    dsu.unite(0, 1); dsu.unite(3, 4);
    cout << "DSU same(0,1): " << (dsu.find(0) == dsu.find(1)) << " same(1,3): " << (dsu.find(1) == dsu.find(3)) << '\n';
}

// =============================================================================
// 16. RECURSION, MEMOIZATION, DP TABLES
// =============================================================================
void sectionDP() {
    unordered_map<int, ll> memo;
    function<ll(int)> fib = [&](int n) -> ll {
        if (n < 2) return n;
        if (memo.count(n)) return memo[n];
        return memo[n] = fib(n - 1) + fib(n - 2);
    };
    cout << "fib(50): " << fib(50) << '\n';

    int m = 3, n = 7;
    vector<vector<int>> mm(m, vector<int>(n, -1));
    function<int(int, int)> paths = [&](int i, int j) -> int {
        if (i == 0 || j == 0) return 1;
        int& res = mm[i][j]; // reference into memo table
        if (res != -1) return res;
        return res = paths(i - 1, j) + paths(i, j - 1);
    };
    cout << "unique paths 3x7 (memo): " << paths(m - 1, n - 1) << '\n';

    string a = "abcde", b = "ace";
    vector<vector<int>> dp(a.size() + 1, vector<int>(b.size() + 1, 0));
    for (size_t i = 1; i <= a.size(); i++)
        for (size_t j = 1; j <= b.size(); j++)
            dp[i][j] = a[i - 1] == b[j - 1] ? dp[i - 1][j - 1] + 1 : max(dp[i - 1][j], dp[i][j - 1]);
    cout << "LCS: " << dp[a.size()][b.size()] << '\n';

    vector<int> coins = {1, 2, 5};
    int amount = 11;
    vector<int> cc(amount + 1, INT_MAX);
    cc[0] = 0;
    for (int i = 1; i <= amount; i++)
        for (int c : coins)
            if (c <= i && cc[i - c] != INT_MAX) cc[i] = min(cc[i], cc[i - c] + 1);
    cout << "coin change: " << cc[amount] << '\n';

    vector<int> nums = {1, 2, 3}, cur;
    vector<vector<int>> subsets;
    function<void(int)> bt = [&](int i) {
        if (i == (int)nums.size()) { subsets.push_back(cur); return; } // push_back copies
        cur.push_back(nums[i]);
        bt(i + 1);
        cur.pop_back();
        bt(i + 1);
    };
    bt(0);
    cout << "subsets count: " << subsets.size() << '\n';
}

// =============================================================================
// 17. BIT MANIPULATION
// =============================================================================
void sectionBits() {
    int x = 0b10110; // 22
    cout << "and/or/xor/not: " << (6 & 3) << ' ' << (6 | 3) << ' ' << (6 ^ 3) << ' ' << (~0) << '\n';
    cout << "shifts: " << (1 << 4) << ' ' << (32 >> 2) << " 64-bit: " << (1LL << 40) << '\n'; // 1<<40 overflows int!
    cout << "popcount: " << __builtin_popcount(x) << " ll: " << __builtin_popcountll(1LL << 40) << '\n';
    cout << "trailing zeros: " << __builtin_ctz(x) << " leading zeros: " << __builtin_clz(x)
         << " bit length: " << 32 - __builtin_clz(x) << " __lg: " << __lg(x) << '\n';
    cout << "binary: " << bitset<8>(x) << '\n';
    int i = 2;
    cout << "check bit i: " << ((x >> i) & 1) << " set: " << (x | (1 << i)) << " clear: " << (x & ~(1 << i))
         << " toggle: " << (x ^ (1 << i)) << '\n';
    cout << "lowest set bit: " << (x & -x) << " drop lowest: " << (x & (x - 1)) << " pow2?: " << ((16 & 15) == 0) << '\n';

    int mask = 0b101;
    vector<int> subs;
    for (int s = mask; s > 0; s = (s - 1) & mask) subs.push_back(s);
    print("submasks of 101: ", subs);

    int single = 0;
    for (int v : {4, 1, 2, 1, 2}) single ^= v;
    cout << "single number: " << single << '\n';
}

// =============================================================================
// 18. COMMON PATTERNS
// =============================================================================
void sectionPatterns() {
    vector<int> nums = {2, 1, 3, 4, 1};
    vector<ll> pre(nums.size() + 1, 0);
    for (size_t i = 0; i < nums.size(); i++) pre[i + 1] = pre[i] + nums[i];
    // or: partial_sum(nums.begin(), nums.end(), pre.begin() + 1);
    cout << "sum[1..3]: " << pre[4] - pre[1] << '\n';

    vector<int> s = {1, 2, 4, 7, 11};
    int l = 0, r = s.size() - 1, target = 15;
    while (l < r) {
        int sum = s[l] + s[r];
        if (sum == target) break;
        sum < target ? l++ : r--;
    }
    cout << "two-sum pair: " << s[l] << ' ' << s[r] << '\n';

    string str = "abcabcbb";
    vector<int> last(128, -1);
    int best = 0, left = 0;
    for (int right = 0; right < (int)str.size(); right++) {
        left = max(left, last[str[right]] + 1);
        last[str[right]] = right;
        best = max(best, right - left + 1);
    }
    cout << "longest unique window: " << best << '\n';

    vector<int> arr = {2, 1, 2, 4, 3}, nge(arr.size(), -1);
    stack<int> stk;
    for (int i = arr.size() - 1; i >= 0; i--) {
        while (!stk.empty() && arr[stk.top()] <= arr[i]) stk.pop();
        if (!stk.empty()) nge[i] = arr[stk.top()];
        stk.push(i);
    }
    print("next greater: ", nge);

    int bestSum = INT_MIN, cur = 0;
    for (int v : {-2, 1, -3, 4, -1, 2, 1, -5, 4}) { cur = max(v, cur + v); bestSum = max(bestSum, cur); }
    cout << "max subarray: " << bestSum << '\n';
}

// =============================================================================
// 19. C++ GOTCHAS
// =============================================================================
//  * int overflows at ~2.1e9: use long long and write 1LL * a * b, 1LL << 40.
//  * v.size() is UNSIGNED: v.size() - 1 when empty = huge number. Cast: (int)v.size() - 1.
//  * unordered_map operator[] inserts missing keys; use count()/find() to test.
//  * priority_queue is a MAX heap by default; comparator logic is "reversed".
//  * Pass big containers by reference (const vector<int>&) — by value copies them.
//  * Sort comparator must be a strict weak ordering: use <, never <=  (else UB / crash).
//  * multiset::erase(value) removes ALL copies; erase(find(value)) removes one.
//  * memset works only for 0 and -1 on int arrays.
//  * Use '\n' instead of endl in loops (endl flushes).
//  * unordered_map can be hacked to O(n) per op on Codeforces; use a custom hash there.
//  * Recursion depth ~1e5–1e6 is OK by default; deeper may stack overflow.

int main() {
    vector<pair<string, function<void()>>> sections = {
        {"1. Numbers", sectionNumbers}, {"2. Strings", sectionStrings}, {"3. Arrays", sectionArrays},
        {"4. Sorting", sectionSorting}, {"5. Binary search", sectionBinarySearch}, {"6. Hash map", sectionHashMap},
        {"7. Hash set", sectionHashSet}, {"8. Ordered map/set", sectionOrdered}, {"9. Stack", sectionStack},
        {"10. Queue", sectionQueue}, {"11. Deque", sectionDeque}, {"12. Heap", sectionHeap},
        {"13. Pairs", sectionPairs}, {"14. Nodes", sectionNodes}, {"15. Graphs", sectionGraphs},
        {"16. DP", sectionDP}, {"17. Bits", sectionBits}, {"18. Patterns", sectionPatterns},
    };
    for (auto& [name, fn] : sections) {
        cout << "\n===== " << name << " =====\n";
        fn();
    }
    (void)fastIOTemplate; // not called: it reads stdin
    return 0;
}
