# COA MCQ/MSQ/NAT Bank — Batch 1 (Q1–Q20)

**What this batch covers (all M1):** base conversion (hex ↔ octal, unknown base, fractions), BCD addition, 2's complement (encoding, overflow, wrap-around, the asymmetric range), signed magnitude vs 1's vs 2's complement, Booth's algorithm (recoding, full trace, worst case), restoring and non-restoring division, IEEE 754 single-precision encoding (including rounding) and decoding (including ULP / machine epsilon).

> Unless a question says otherwise: integers are 2's complement; IEEE 754 single precision = 1 sign bit, 8 exponent bits (bias 127), 23 fraction bits; rounding is round-to-nearest-even; C code assumes a 2's-complement machine with IEEE 754 floats.

---

## Q1. [M1 | Base conversion | Easy | MCQ | OA]

Convert the hexadecimal number 0x2F7 to octal.

- A. 1367
- B. 5734
- C. 759
- D. 2177

**Answer:** A

**Why:** Hex and octal are both powers of 2, so you convert through binary. Each hex digit becomes 4 bits. Then you regroup the bits into 3s **starting from the right**, because each octal digit is 3 bits.

**Shortcut / Working:**
```
Hex:     2    F    7
Binary:  0010 1111 0111
Regroup from the RIGHT in 3s:  001 011 110 111
Octal:                          1   3   6   7   → 1367
Check: 0x2F7 = 2·256 + 15·16 + 7 = 759 ; 1367₈ = 512 + 192 + 48 + 7 = 759 ✓
```

**Why the other options fail:**
- B: The bits were grouped from the left, with zeros padded on the right. That multiplies the value by 4.
- C: 759 is the decimal value, written as if it were octal. The digit 9 doesn't exist in octal.
- D: Each hex digit was converted separately (F → 17). You can't convert digit by digit between base 16 and base 8.

**Concept recap:** Hex ↔ octal always goes through binary. Group bits outward from the binary point: to the left for the integer part, to the right for the fraction.

**Trap alert:** Padding on the wrong side. Integer parts get zeros added on the **left**; fractions get zeros added on the **right**.

**Interview angle:** Shows up in OAs as a quick warm-up. Follow-up: *"Why do programmers use hex rather than octal for bytes?"* A byte is 8 bits, which is exactly 2 hex digits. In octal it is an awkward 2⅔ digits.

**Time target:** 30 s

---

## Q2. [M1 | Base conversion (unknown base) | Medium | NAT | GATE-NAT]

(121)ᵣ = (144)₈, where r is a whole-number base greater than 2. Find r.

*Answer is a number — no options.*

**Answer:** 9

**Why:** Convert both sides to decimal. The right side is a known number. The left side becomes a polynomial in r, which you solve.

**Shortcut / Working:**
```
(144)₈ = 1·64 + 4·8 + 4 = 100
(121)ᵣ = r² + 2r + 1 = (r + 1)²
(r + 1)² = 100  →  r + 1 = 10  →  r = 9     (r = −11 rejected)
Digit check: largest digit 2 < 9 ✓
```
Shortcut: "121" in any base equals (r + 1)².

**Common wrong answers:** 11, from treating (144)₈ as decimal 144 and solving (r + 1)² = 144.

**Concept recap:** (dₙ…d₁d₀)ᵣ = Σ dᵢ·rⁱ. Every digit must be smaller than r.

**Trap alert:** Convert the known side to decimal first. Then check that the r you found is larger than every digit.

**Interview angle:** A GATE-style question; rare in interviews. The same skill (positional value) is behind `atoi`/`strtol` with a base argument, which does come up in coding rounds.

**Time target:** 45 s

---

## Q3. [M1 | Base conversion (fractions) | Medium | MSQ | GATE]

Which of the following equal decimal 45.625? (One or more may be correct.)

- A. (101101.101)₂
- B. (55.5)₈
- C. (2D.5)₁₆
- D. (231.22)₄

**Answer:** A, B, D

**Why:** In binary, 45 = 101101 and 0.625 = 0.101 (because 0.5 + 0.125 = 0.625). To get the other bases, regroup these bits. Fraction groups are counted rightward from the point, and padding zeros go on the right.

**Shortcut / Working:**
```
Binary: 101101.101
Octal  (3s): 101 101 . 101          → 55.5      ✓
Hex    (4s): 0010 1101 . 1010       → 2D.A      (not 2D.5)
Base-4 (2s): 10 11 01 . 10 10       → 231.22    ✓  (= 45 + 2/4 + 2/16)
```

**Why the other options fail:**
- C: The fraction bits .101 were padded on the left (0101 = 5) instead of the right (1010 = A). (2D.5)₁₆ = 45.3125.

**Concept recap:** Zeros added after the last fraction digit don't change the value (0.5 = 0.50). Zeros added before it do.

**Trap alert:** Integer part: pad on the left. Fraction part: pad on the right.

**Interview angle:** Rare in interviews. Useful background for *"why can't 0.1 be stored exactly in binary?"*: some fractions never terminate in base 2.

**Time target:** 60 s

---

## Q4. [M1 | BCD addition | Easy | MCQ | GATE]

The decimal numbers 92 and 38 are stored in packed BCD as 1001 0010 and 0011 1000. A plain binary adder gives 1100 1010 (0xCA). What correction produces the correct BCD result?

- A. Add 0110 to the low nibble only
- B. Add 0110 to the high nibble only
- C. Add 0110 to both nibbles, giving 1 0011 0000
- D. No correction; 1100 1010 is valid BCD

**Answer:** C

**Why:** In BCD (Binary-Coded Decimal: each decimal digit is stored in its own 4-bit nibble), any nibble greater than 9, or any nibble that produced a carry, must have 6 (0110) added. Adding 6 skips over the six unused codes 1010–1111. Here both nibbles are invalid: the low nibble is 1010 and the high nibble is 1100.

**Shortcut / Working:**
```
Low nibble : 0010 + 1000 = 1010 (>9) → +0110 → 0000, carry 1
High nibble: 1001 + 0011 + 1 = 1101 (>9) → +0110 → 0011, carry 1
Result: 1 0011 0000 = BCD 130 ✓   (92 + 38 = 130)
Shortcut: 0xCA + 0x66 = 0x130
```

**Why the other options fail:**
- A: Fixes only the low digit. The high nibble stays at 1101, which isn't valid BCD.
- B: Ignores both the invalid low nibble and the carry it produces.
- D: Codes 1010–1111 don't exist in BCD.

**Concept recap:** BCD addition: add in binary. For each nibble, if it is > 9 or it produced a carry, add 0110 and pass the carry to the next nibble.

**Trap alert:** The carry from correcting the low nibble can push the high nibble above 9.

**Interview angle:** Rare in interviews. Context: x86 had a dedicated instruction (`DAA`) for this correction, and financial and embedded systems still use BCD.

**Time target:** 40 s

---

## Q5. [M1 | 2's complement | Easy | MCQ | SDE-COA]

In C: `int8_t x = -75; printf("%02X", (uint8_t)x);`. Assume a 2's-complement machine. What is printed?

- A. CB
- B. B4
- C. 4B
- D. B5

**Answer:** D

**Why:** Casting a signed value to an unsigned type of the same width leaves the bits unchanged. So the program prints the 2's-complement bit pattern of −75: invert the bits of +75, then add 1.

**Shortcut / Working:**
```
+75      = 0100 1011
invert   = 1011 0100
+1       = 1011 0101 = 0xB5
Shortcut: 256 − 75 = 181 = 0xB5
```

**Why the other options fail:**
- A: This is signed magnitude: sign bit 1 followed by 100 1011.
- B: This is 1's complement; the +1 was forgotten.
- C: This is +75; the sign was dropped.

**Concept recap:** In n-bit 2's complement, −x is stored as 2ⁿ − x. A signed-to-unsigned cast of the same width doesn't change any bits.

**Trap alert:** Forgetting the "+1" gives the 1's-complement pattern.

**Interview angle:** *"What does this print?"* Model answer: "The cast changes how the bits are read, not the bits themselves. −75 is 256 − 75 = 181 = 0xB5." Follow-up: *"What if it were cast to `uint16_t`?"* The value is sign-extended first, giving 0xFFB5.

**Time target:** 20 s

---

## Q6. [M1 | 2's complement overflow | Medium | MSQ | GATE]

Each pair below is added with an 8-bit 2's-complement adder. Which additions overflow? (One or more may be correct.)

- A. 0x7F + 0x01
- B. 0x80 + 0xFF
- C. 0xC0 + 0xC0
- D. 0x50 + 0xB0

**Answer:** A, B

**Why:** Overflow means the true result falls outside −128 to +127. In hardware, it happens exactly when the carry **into** the most significant bit (MSB) differs from the carry **out** of it. Equivalently: two operands with the same sign produce a result with the opposite sign.

**Shortcut / Working:**
```
     Operands (decimal)   True sum   8-bit result   Overflow?
A    127 + 1               128        0x80 (−128)    YES (+ + → −)
B   −128 + (−1)           −129        0x7F (+127)    YES (− − → +)
C    −64 + (−64)          −128        0x80 (−128)    NO  (fits exactly)
D     80 + (−80)             0        0x00            NO  (carry-out = 1, but signs differ)
```

**Why the other options fail:**
- C: −128 is the most negative value that fits, so there's no overflow. This catches the belief that "two negatives always overflow."
- D: There is a carry out of the MSB, but carry-out is not overflow. Adding operands of opposite sign can never overflow.

**Concept recap:** V = C_in(MSB) XOR C_out(MSB). Equivalently, overflow happens if and only if the operands have the same sign and the result's sign is different.

**Trap alert:** Mixing up the carry flag (for unsigned numbers) with the overflow flag (for signed numbers).

**Interview angle:** *"How do you detect signed overflow in C?"* Model answer: "Signed overflow is undefined behaviour in C, so check before adding: `if (b > 0 && a > INT_MAX - b)`, or use `__builtin_add_overflow`." Follow-up: *"How does the CPU detect it?"* It XORs the carry into the MSB with the carry out of it.

**Time target:** 60 s

---

## Q7. [M1 | 2's complement subtraction | Medium | NAT | GATE-NAT]

X = 0xA7 and Y = 0x3E are 8-bit 2's-complement numbers. An 8-bit ALU computes X − Y and keeps only 8 bits. What is the decimal value of the 8-bit result, read as 2's complement?

*Answer is a number — no options.*

**Answer:** 105

**Why:** X = −89 and Y = +62, so the true result is −151. That is below −128, so it overflows. The 8-bit result is the true value modulo 256, which happens to be positive.

**Shortcut / Working:**
```
X = 0xA7 = 167 − 256 = −89 ;  Y = 0x3E = 62
True X − Y = −151 (outside −128…127 → overflow)
Hardware: 0xA7 − 0x3E = 0x69 = 105        (same as −151 + 256)
Sign check: negative − positive → positive result ⇒ V = 1
```

**Common wrong answers:** −151 (the true value, which doesn't fit in 8 bits); −105 (a sign slip).

**Concept recap:** The n-bit result is (true result) mod 2ⁿ, read as 2's complement. Subtraction can overflow only when the two operands have opposite signs.

**Trap alert:** If the question says "keep only 8 bits," give the wrapped value.

**Interview angle:** Same idea as integer wrap-around bugs, e.g. `uint8_t` counters or `size_t` underflow in loops like `for (size_t i = n - 1; i >= 0; i--)`, which never ends.

**Time target:** 45 s

---

## Q8. [M1 | 2's complement range asymmetry | Easy | MCQ | SDE-COA]

A 32-bit ALU negates 0x80000000 using "invert all bits, then add 1". What is the result?

- A. 0x7FFFFFFF
- B. 0x80000000
- C. 0x00000000
- D. 0xFFFFFFFF

**Answer:** B

**Why:** The n-bit 2's-complement range is −2ⁿ⁻¹ to 2ⁿ⁻¹ − 1. There is one more negative value than positive values, so the most negative value has no positive partner. Negating it wraps around to itself and sets the overflow flag.

**Shortcut / Working:**
```
0x80000000 → invert → 0x7FFFFFFF → +1 → 0x80000000   (V = 1)
```

**Why the other options fail:**
- A: The +1 was skipped. 0x7FFFFFFF is INT_MAX.
- C: The result doesn't wrap to zero.
- D: 0xFFFFFFFF is −1.

**Concept recap:** Range = [−2ⁿ⁻¹, 2ⁿ⁻¹ − 1]. Negating −2ⁿ⁻¹ overflows back to −2ⁿ⁻¹.

**Trap alert:** Assuming every number has a representable negation.

**Interview angle:** *"Why is `abs(INT_MIN)` a bug?"* Model answer: "2's complement has one more negative number than positive ones, so INT_MIN has no positive counterpart. In C, negating it is undefined behaviour; the hardware returns INT_MIN unchanged." Follow-up: *"How do you safely compute the midpoint of two ints?"* Use `lo + (hi - lo) / 2`.

**Time target:** 20 s

---

## Q9. [M1 | Signed magnitude / 1's / 2's complement | Medium | MCQ | OA]

The 8-bit pattern 1111 1010 is read in three ways: signed magnitude, 1's complement, and 2's complement. What are the three decimal values, in that order?

- A. −122, −6, −5
- B. −122, −5, −6
- C. −6, −5, −122
- D. −5, −122, −6

**Answer:** B

**Why:** In all three systems, an MSB of 1 means negative. They differ in how they get the magnitude:
- Signed magnitude reads the low 7 bits directly.
- 1's complement inverts all the bits.
- 2's complement inverts all the bits and adds 1.

**Shortcut / Working:**
```
Signed magnitude : magnitude = 111 1010 = 122          → −122
1's complement   : invert → 0000 0101 = 5              → −5
2's complement   : invert + 1 → 0000 0110 = 6          → −6
```

**Why the other options fail:**
- A: The 1's- and 2's-complement values are swapped.
- C, D: The three readings are in the wrong order.

**Concept recap:** For a negative pattern: SM = −(low n−1 bits); 1's = −(NOT of the pattern); 2's = −(NOT of the pattern + 1). Signed magnitude and 1's complement both have two zeros; 2's complement has one.

**Trap alert:** Signed magnitude never inverts any bits.

**Interview angle:** *"Why does every modern CPU use 2's complement?"* Model answer: "There's only one zero, and the same adder works for signed and unsigned numbers. Subtraction is just adding the 2's complement, so no separate subtractor is needed."

**Time target:** 40 s

---

## Q10. [M1 | Booth's algorithm (recoding) | Medium | MCQ | GATE]

Booth's algorithm multiplies two 8-bit numbers. The multiplier is 0111 0110 and Q₋₁ starts at 0. How many add/subtract operations are performed, and how many arithmetic-shift steps?

- A. 4 add/subtract operations, 8 shifts
- B. 4 add/subtract operations, 9 shifts
- C. 5 add/subtract operations, 8 shifts
- D. 2 add/subtract operations, 8 shifts

**Answer:** A

**Why:** At each step Booth looks at the bit pair (Qᵢ, Qᵢ₋₁). A pair of 10 means subtract M; 01 means add M; 00 and 11 mean do nothing. Each of the n steps ends with an arithmetic right shift, so the shift count is n = 8. The extra Q₋₁ bit makes the register wider but does **not** add a step.

**Shortcut / Working:** Append a 0 on the right of the multiplier, then count the places where adjacent bits differ.
```
Bits b7..b0 = 0 1 1 1 0 1 1 0 | Q−1 = 0
Pair (bᵢ, bᵢ₋₁): i=0:(0,0)–  i=1:(1,0)SUB  i=2:(1,1)–  i=3:(0,1)ADD
                 i=4:(1,0)SUB i=5:(1,1)–   i=6:(1,1)–  i=7:(0,1)ADD
Recoded: +1 0 0 −1 +1 0 −1 0  → 128 − 16 + 8 − 2 = 118 = 0111 0110 ✓
Ops = 4 (2 ADD, 2 SUB), shifts = 8
```

**Why the other options fail:**
- B: Adds an extra step for Q₋₁, which doesn't exist.
- C: Counts the 1-bits (5) instead of the transitions between bits.
- D: Counts only the additions.

**Concept recap:** n-bit Booth = n iterations, each made of an optional add/subtract followed by an arithmetic right shift. The number of operations equals the number of bit transitions in the multiplier with a 0 appended on the right.

**Trap alert:** "Steps = n + 1" is a common wrong belief. Q₋₁ makes the register one bit wider; it doesn't add an iteration.

**Interview angle:** GATE-focused; rarely asked in SDE interviews. If it comes up, the one-liner is: "Booth recoding turns a run of 1s into one subtraction at the start of the run and one addition at the end, which cuts the number of additions for multipliers with long runs of 1s."

**Time target:** 45 s

---

## Q11. [M1 | Booth's algorithm (full trace) | Hard | NAT | GATE-NAT]

Booth's algorithm with 4-bit registers. Multiplicand M = −5 (1011), multiplier Q = +3 (0011). Initially A = 0000 and Q₋₁ = 0. What is the decimal value of A, read as 4-bit 2's complement, at the end of iteration 3 (after that iteration's arithmetic right shift)?

*Answer is a number — no options.*

**Answer:** −2

**Why:** Each iteration looks at (Q₀, Q₋₁), adds or subtracts M if needed, then arithmetic-right-shifts the combined A : Q : Q₋₁ register. "Arithmetic" means the sign bit of A is copied into the vacated top bit.

**Shortcut / Working:** M = 1011, −M = 0101.
```
Iter  Q0Q−1  Action        A      Q      Q−1   (after ASR)
 0     –     init          0000   0011   0
 1     10    A=A−M=0101    0010   1001   1
 2     11    none          0001   0100   1
 3     01    A=A+M=1100    1110   0010   0     ← A = 1110 = −2
 4     00    none          1111   0001   0
Final A:Q = 1111 0001 = −15 = (−5)(3) ✓
```

**Common wrong answers:** 6 (0110), from using a logical shift instead of an arithmetic one; −4 (1100), from reading A before the shift.

**Concept recap:** Booth step: 10 → A = A − M; 01 → A = A + M; 00 or 11 → no operation. Then arithmetic-shift A : Q : Q₋₁ right by 1, copying the sign of A.

**Trap alert:** Booth always uses an arithmetic shift, never a logical one.

**Interview angle:** GATE-only. Skip it for SDE preparation.

**Time target:** 2 min

---

## Q12. [M1 | Booth's algorithm (worst case) | Medium | MCQ | GATE]

Q₋₁ starts at 0. For which 8-bit multiplier does Booth's algorithm perform the most add/subtract operations?

- A. 1111 1111
- B. 0101 0101
- C. 1010 1010
- D. 0000 1111

**Answer:** B

**Why:** Booth performs one operation for every transition between adjacent bits, counting the appended Q₋₁ = 0 as the bit to the right of b0. An alternating pattern whose b0 is 1 has a transition at every one of the 8 positions.

**Shortcut / Working:**
```
Append 0 on the right and count transitions:
A  1111 1111 0  → 1      B  0101 0101 0  → 8
C  1010 1010 0  → 7      D  0000 1111 0  → 2
```

**Why the other options fail:**
- A: All 1s is Booth's best case: a single subtraction.
- C: Also alternating, but b0 = 0 matches the appended 0, so there are only 7 transitions. This is the near-miss option.
- D: A single run of 1s needs only 2 operations.

**Concept recap:** Booth does best with long runs of 1s and worst with alternating bits, where it needs n operations, no better than plain shift-and-add.

**Trap alert:** Forgetting to count the appended 0.

**Interview angle:** GATE-focused. The underlying idea, that runs of 1s can be replaced by one subtraction and one addition, is the same trick as computing `x * 15` as `(x << 4) - x`.

**Time target:** 40 s

---

## Q13. [M1 | Restoring division | Medium | MCQ | OA]

Restoring division of 1011 (11) by 0011 (3), with a 5-bit A register (initially 0) and a 4-bit Q register. How many times is A restored (M added back)?

- A. 1
- B. 3
- C. 2
- D. 4

**Answer:** C

**Why:** Each step shifts A : Q left and subtracts M. If A goes negative, that quotient bit is 0 and A is restored by adding M back. So the number of restores equals the number of 0 bits in the n-bit quotient.

**Shortcut / Working:** Quotient = 11 ÷ 3 = 3 = 0011, which has two 0 bits, so there are **2** restores.
```
Step  After shift (A | Q)   A − M          q   Restore?   A after
 1    00001 | 011_          11110 (neg)    0   yes        00001
 2    00010 | 110_          11111 (neg)    0   yes        00010
 3    00101 | 100_          00010 (≥0)     1   no         00010
 4    00101 | 001_          00010 (≥0)     1   no         00010
Q = 0011 (3), remainder A = 00010 (2) ✓
```

**Why the other options fail:**
- A, B: Miscounted from the trace.
- D: Assumes every step restores.

**Concept recap:** Restoring division = n subtractions + (number of 0 bits in the quotient) restores.

**Trap alert:** A must be n + 1 bits wide, so the sign of A − M can be seen.

**Interview angle:** Rare in interviews. The related practical point: hardware division is much slower than multiplication, which is why compilers replace division by a constant with a multiply and a shift.

**Time target:** 30 s with the "count the quotient's 0s" shortcut; 2 min with a full trace

---

## Q14. [M1 | Non-restoring division | Hard | NAT | GATE-NAT]

Non-restoring division of 1000 (8) by 0011 (3), with a 5-bit A register (initially 0) and a 4-bit Q register. At each step: if A ≥ 0, shift then subtract M; if A < 0, shift then add M. At the end, if A < 0, add M once to correct the remainder. How many add/subtract operations involving M happen in total, including any final correction?

*Answer is a number — no options.*

**Answer:** 5

**Why:** Non-restoring division does exactly one add or subtract per step, so 4 for 4 bits. It needs one extra correction only if the final remainder is negative, and here it is.

**Shortcut / Working:**
```
Step  A sign before  Shift → A    Op     A after     q
 1    ≥0             00001        −M     11110 (−2)   0
 2    <0             11100        +M     11111 (−1)   0
 3    <0             11110        +M     00001 (+1)   1
 4    ≥0             00010        −M     11111 (−1)   0
Final A < 0 → correction +M → 00010 (2)
Q = 0010 (2), R = 2 ✓   Ops = 4 + 1 = 5
(Restoring on the same inputs: 4 subtractions + 3 restores = 7.)
```

**Common wrong answers:** 4 (forgot the final correction); 7 (the restoring-division count).

**Concept recap:** Non-restoring = n operations + 1 if the final remainder is negative. Restoring = n operations + (number of 0s in the quotient).

**Trap alert:** The final correction is part of the algorithm, not an optional extra.

**Interview angle:** GATE-only.

**Time target:** 2 min

---

## Q15. [M1 | IEEE 754 single encode | Easy | MCQ | OA]

What is the IEEE 754 single-precision encoding of −6.25?

- A. 0xC0C80000
- B. 0xC1480000
- C. 0x40C80000
- D. 0xC0E40000

**Answer:** A

**Why:** Convert to binary, normalise to the form 1.xxx × 2ᴱ, and store E + 127 in the exponent field. Store only the fraction bits after the point; the leading 1 is implied and not stored (the "hidden bit").

**Shortcut / Working:**
```
6.25 = 110.01₂ = 1.1001 × 2²
S = 1 | Exp = 2 + 127 = 129 = 1000 0001 | Frac = 1001 000…0
1 10000001 10010000000000000000000
1100 0000 1100 1000 0000 … → 0xC0C80000
```

**Why the other options fail:**
- B: Used bias 128, giving an exponent field of 130.
- C: Sign bit is 0, so this is +6.25.
- D: Stored the hidden 1 in the fraction field.

**Concept recap:** value = (−1)ˢ × 1.F × 2^(Exp − 127). The bias is 127 for single precision and 1023 for double.

**Trap alert:** Using bias 128 and storing the hidden 1 are the two most common encoding mistakes.

**Interview angle:** *"Walk me through how a float is stored."* Model answer: "One sign bit, 8 exponent bits stored with a bias of 127, and 23 fraction bits with an implied leading 1. So we get 24 bits of precision for the price of 23."

**Time target:** 60 s

---

## Q16. [M1 | IEEE 754 single encode | Medium | NAT | GATE-NAT]

−0.75 is encoded in IEEE 754 single precision (bias 127). How many 1 bits are in the 32-bit pattern?

*Answer is a number — no options.*

**Answer:** 8

**Why:** A negative exponent is stored below the bias. 2⁻¹ is stored as 127 − 1 = 126 = 0111 1110, which has six 1s.

**Shortcut / Working:**
```
0.75 = 0.11₂ = 1.1 × 2⁻¹
S = 1 (one 1) | Exp = 126 = 0111 1110 (six 1s) | Frac = 1000…0 (one 1)
Total = 1 + 6 + 1 = 8      Pattern = 0xBF400000
```

**Common wrong answers:** 9 (exponent stored as 127 = 0111 1111); 7 (the sign bit was dropped).

**Concept recap:** Stored exponent = true exponent + bias. A true exponent below 0 gives a stored exponent below 127.

**Trap alert:** The exponent field is biased, not 2's complement.

**Interview angle:** *"Why use a bias instead of 2's complement for the exponent?"* Model answer: "With a bias, positive floats sort in the same order as their bit patterns read as integers, so the hardware can compare floats using integer comparison."

**Time target:** 45 s

---

## Q17. [M1 | IEEE 754 encode with rounding | Hard | NAT | GATE-NAT]

1/3 is encoded in IEEE 754 single precision using round-to-nearest-even. What is the decimal value of the 4 least significant bits of the 23-bit fraction field, read as a binary number?

*Answer is a number — no options.*

**Answer:** 11

**Why:** 1/3 = 0.010101…₂ = 1.0101…₂ × 2⁻². The 23 fraction bits that fit end in …010. The bits that don't fit start 1010…, which is more than half of the last kept bit's weight, so the result rounds up to …011.

**Shortcut / Working:**
```
Normalised: 1.0101 0101 0101 0101 0101 010 | 1 0101… × 2⁻²
                 ← 23 fraction bits →         ↑ guard = 1, rest ≠ 0 ⇒ > half ⇒ round up
Kept bits end …0 1010 → +1 → …0 1011
Last 4 bits = 1011 = 11       Full encoding 0x3EAAAAAB
```

**Common wrong answers:** 10 (1010, from truncating instead of rounding).

**Concept recap:** Round-to-nearest-even: if the discarded bits are more than half, round up; if less than half, round down; if exactly half, round so the last kept bit is 0.

**Trap alert:** IEEE 754's default rounding mode is round-to-nearest-even, not truncation.

**Interview angle:** *"Why does `0.1 + 0.2 != 0.3`?"* Model answer: "None of the three has a finite binary expansion. Each is rounded to the nearest double, and the rounding errors in 0.1 and 0.2 don't add up to the rounding error in 0.3. Compare floats with a tolerance, not with `==`."

**Time target:** 90 s

---

## Q18. [M1 | IEEE 754 single decode | Easy | MCQ | GATE]

Decode the IEEE 754 single-precision pattern 0x41200000.

- A. 5.0
- B. 2.0
- C. 10.0
- D. 20.0

**Answer:** C

**Why:** Split the pattern into sign, exponent and fraction. Subtract 127 from the exponent and put back the hidden 1.

**Shortcut / Working:**
```
0x41200000 = 0100 0001 0010 0000 …
S = 0 | Exp = 1000 0010 = 130 → E = 3 | Frac = 01000… → 1.01₂ = 1.25
Value = 1.25 × 2³ = 10.0
```

**Why the other options fail:**
- A: Bias 128 gives 1.25 × 4 = 5.
- B: Forgetting the hidden 1 gives 0.25 × 8 = 2.
- D: An exponent off by one gives 1.25 × 16 = 20.

**Concept recap:** Decode: (−1)ˢ × (1 + F) × 2^(Exp − 127), valid for exponent fields 1 to 254.

**Trap alert:** If a positive float's hex pattern starts with 4, its exponent field starts with 100…, so the value is at least 2.

**Interview angle:** Debugger skill: being able to recognise 0x3F800000 = 1.0f and 0x40000000 = 2.0f in a memory dump is handy when debugging.

**Time target:** 40 s

---

## Q19. [M1 | IEEE 754 single decode | Medium | NAT | GATE-NAT]

Decode the IEEE 754 single-precision pattern 0xC2F50000 to a decimal number.

*Answer is a number — no options.*

**Answer:** −122.5

**Why:** Sign 1 means negative. The exponent field is 133, so E = 133 − 127 = 6. Shift the binary point of 1.F six places to the right.

**Shortcut / Working:**
```
0xC2F50000 = 1100 0010 1111 0101 0000 …
S = 1 | Exp = 1000 0101 = 133 → E = 6 | Frac = 1110101 0…
1.1110101 × 2⁶ = 1111010.1₂ = 64+32+16+8+2 + 0.5 = 122.5
Value = −122.5
```

**Common wrong answers:** −61.25 (bias 128); −58.5 (dropped the hidden 1).

**Concept recap:** Once you know E, move the binary point E places to the right in 1.F.

**Trap alert:** The exponent field (bits 30–23) straddles a hex-digit boundary, so split the bits carefully.

**Interview angle:** Rare as a direct interview question, but OAs often ask you to decode a hex float pattern like this.

**Time target:** 60 s

---

## Q20. [M1 | IEEE 754 decode / ULP | Hard | MCQ | SDE-COA]

```c
uint32_t u = 0x3F800001; float f;
memcpy(&f, &u, 4);
```
Which statement is true (IEEE 754 single precision)?

- A. f == 1.0f, because the last bit is rounded away
- B. f = 1 + 2⁻²⁴
- C. f is a denormal number
- D. f = 1 + 2⁻²³, the next float after 1.0, and f − 1.0f equals 2⁻²³ exactly

**Answer:** D

**Why:** 0x3F800000 is exactly 1.0 (exponent 127 means E = 0, and the fraction is 0). Setting the lowest fraction bit adds one **ULP** (unit in the last place). At 1.0 a ULP is 2⁻²³, which is FLT_EPSILON. The subtraction f − 1.0f is exact: when two floats are within a factor of 2 of each other, their difference is always representable exactly (Sterbenz's lemma).

**Shortcut / Working:**
```
0x3F800001: S=0 | Exp=0111 1111=127 → E=0 | Frac=000…001 → 1 + 2⁻²³
2⁻²³ ≈ 1.1920929e−7 = FLT_EPSILON
```

**Why the other options fail:**
- A: `memcpy` copies the bits exactly; no conversion or rounding happens.
- B: 2⁻²⁴ is half a ULP at 1.0 (the unit roundoff), not the gap between neighbouring floats.
- C: A denormal needs an exponent field of 0; here it is 127.

**Concept recap:** In the range [2ᴱ, 2ᴱ⁺¹), consecutive floats are 2^(E−23) apart. At 1.0 the gap is FLT_EPSILON = 2⁻²³.

**Trap alert:** Mixing up machine epsilon (2⁻²³) with unit roundoff (2⁻²⁴).

**Interview angle:** *"Why `memcpy` instead of `*(float*)&u`?"* Model answer: "Casting the pointer breaks C's strict-aliasing rule, which is undefined behaviour. `memcpy` is the well-defined way to reinterpret bits, and compilers optimise it to a single register move." Follow-up: *"What's FLT_EPSILON?"* The gap between 1.0 and the next float, 2⁻²³.

**Time target:** 60 s

---

**Completed:** 20 / 340 · **Remaining:** 320

**Ledger — newly ✓ this batch:** P01 ★ (3) · P02 ★ (4) · P03 (1) · P04 ★ (3) · P05 (2) · P06 ★ (3) · P07 ★ (3) · BCD extra (1)

**Batch stats:** MCQ 11 · MSQ 2 · NAT 7 | Easy 6 · Medium 10 · Hard 4 | GATE 6 · GATE-NAT 7 · OA 4 · SDE-COA 3 | MCQ key: A 3 · B 3 · C 3 · D 2 | Numerical/trace-based: 18 of 20
