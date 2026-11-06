# Aptitude MCQ/MSQ/NAT Bank — Batch 1 (Q1–Q20)

**What this batch covers (all M1):** percentage change and the asymmetry trap, successive percentage change, consumption–expenditure, turning a percentage equation into a ratio, layered-base percentages (elections), pass-mark problems, extra quantity after a price drop, ratio splitting with fractions, chained ratios, proportion properties (mean, third and fourth proportional, componendo–dividendo), ratio change after transfers, mixing two ratio mixtures, income–expenditure–savings, partnership (simple, time-varying capital, working partner, reverse-solving for time), and age problems.

**Patterns:** P01 ×7, P02 ×6, P03 ×5, X1 ×2

**Conventions for this batch:**
- All money is in ₹.
- "x% more/less than Y" always means x% of Y.
- In partnership questions, profit is shared in proportion to capital × time, unless the question says otherwise.
- All months are calendar months of equal length.

---

### Q1. [M1 | Percentage change | Easy | MCQ | OA-Speed | Service]
The price of a laptop is increased by 20%. Later, the new price is decreased by 20%. Compared with the original price, the final price is:

* A. Unchanged
* B. 4% higher
* C. 4% lower
* D. 2% lower

**Answer:** C

**Why:** The two 20% changes are taken on different bases. The increase is 20% of the original price. The decrease is 20% of the larger, increased price, so more comes off than went on. The changes multiply; they do not cancel.

**Shortcut / Working:**
- Given: +20%, then −20%. Find: net % change.
- Formula for two successive changes: net = a + b + ab/100.
- Net = 20 + (−20) + (20 × −20)/100 = 0 − 4 = **−4%**.
- Faster rule: an increase of x% followed by a decrease of x% (in either order) always gives a net loss of x²/100 %. Here 20²/100 = 4.
- Check with 100: 100 → 120 → 120 × 0.8 = 96, which is 4% lower.

**Why the other options fail:**
* A: Assumes +20 and −20 cancel. That is additive thinking, which is the trap here.
* B: Right size, wrong sign. The xy/100 term is negative when the signs differ.
* D: Halves the correct answer, perhaps from 20 × 20/200. There is no such formula.

**Concept recap:** Successive changes multiply: final = initial × (1 + a/100)(1 + b/100). Net % = a + b + ab/100. If the changes are +x and −x, net = −x²/100 %.

**Trap alert:** "Up 20, down 20, back to start" is the most common wrong answer in OAs. The base changes between steps.

**Time target:** 10 s with the x²/100 rule; 30 s otherwise.

---

### Q2. [M1 | Successive change (3 factors) | Medium | MCQ | OA-Speed | Service]
A rectangular water tank is redesigned. Its length increases by 20%, its breadth decreases by 20%, and its height increases by 10%. By what percentage does its volume change?

* A. +10%
* B. +5.6%
* C. −4%
* D. +6%

**Answer:** B

**Why:** Volume = l × b × h, so each % change acts as a multiplying factor. Apply the two-change formula twice: first to length and breadth, then to that result and height.

**Shortcut / Working:**
- Given: +20%, −20%, +10%. Find: % change in the product (volume).
- Step 1 (l and b): 20 − 20 + (20 × −20)/100 = −4%.
- Step 2 (combine −4% with +10%): −4 + 10 + (−4 × 10)/100 = 6 − 0.4 = **+5.6%**.
- Multiplier check: 1.2 × 0.8 × 1.1 = 0.96 × 1.1 = 1.056, so +5.6%.

**Why the other options fail:**
* A: Adds the three changes (20 − 20 + 10). Ignores compounding entirely.
* C: Leaves out the height change; −4% is only the l × b result.
* D: Gets −4% correctly, then adds +10% without the cross term (−4 × 10)/100 = −0.4.

**Concept recap:** For a product of quantities, net % change = combine the changes two at a time with a + b + ab/100, or multiply the factors (1 ± r/100). This works for area, volume, revenue (price × quantity) and expenditure.

**Trap alert:** Students apply the formula correctly once, then add the third change. Every step needs the cross term.

**Time target:** 20 s with the two-step formula.

---

### Q3. [M1 | Consumption–expenditure | Easy | MCQ | OA-Speed | Service]
The price of sugar rises by 25%. By what percentage must a household reduce its sugar consumption so that its spending on sugar stays the same?

* A. 20%
* B. 25%
* C. 18%
* D. 22.5%

**Answer:** A

**Why:** Expenditure = price × consumption. If price is multiplied by 5/4, consumption must be multiplied by 4/5 to keep the product constant. That is a cut of 1/5, or 20%. The percentage cut is measured on the old consumption, which is why it is smaller than 25%.

**Shortcut / Working:**
- Given: price +25%. Find: % cut in consumption for constant expenditure.
- Formula: reduction = r / (100 + r) × 100 = 25/125 × 100 = **20%**.
- Fraction shortcut: 25% = 1/4. Price becomes ×5/4, so consumption becomes ×4/5, a cut of 1/5 = 20%.
- Check: 100 units at ₹1 = ₹100. New price ₹1.25 × 80 units = ₹100 ✓.

**Why the other options fail:**
* B: Assumes the cut equals the rise, which is the asymmetry trap again.
* C: A guess near 20% from estimating instead of using the formula. The exact value is a clean 20%.
* D: Averages 20 and 25. Not a valid operation.

**Concept recap:** If price rises by r%, cut consumption by r/(100 + r) × 100 %. If price falls by r%, raise consumption by r/(100 − r) × 100 %. The same formula answers "A is r% more than B; how much less is B than A?"

**Trap alert:** If you mix up the rise case and the fall case, you will use the wrong denominator: (100 + r) for a rise, (100 − r) for a fall.

**Time target:** 10 s.

---

### Q4. [M1 | Percentage equation → ratio | Easy | MCQ | OA-Speed | Service]
If 40% of A equals 60% of B, then A is what percent more than B?

* A. 20%
* B. 33.33%
* C. 40%
* D. 50%

**Answer:** D

**Why:** Turn the percentage equation into a ratio first. 0.4A = 0.6B gives A/B = 0.6/0.4 = 3/2. "A is how much more than B" measures the gap on B's base: (3 − 2)/2 = 50%.

**Shortcut / Working:**
- Given: 40% of A = 60% of B. Find: (A − B)/B × 100.
- Cross rule: if x% of A = y% of B, then A : B = y : x. So A : B = 60 : 40 = 3 : 2.
- (A − B)/B = 1/2 = **50%**.

**Why the other options fail:**
* A: Subtracts the percentages (60 − 40). Percentages of different bases cannot be subtracted.
* B: That is how much *less* B is than A ((3 − 2)/3). The wrong base.
* C: Reads 40 directly from the question.

**Concept recap:** x% of A = y% of B ⇒ A : B = y : x (the larger percentage goes with the smaller number). "More than B" uses B as the base. "Less than A" uses A as the base.

**Trap alert:** The base word comes after "than". "More than B" means divide by B.

**Time target:** 10 s.

---

### Q5. [M1 | Layered-base percentages (election) | Hard | NAT | OA | Service]
In a town election between two candidates:
- 20% of the registered voters did not vote.
- Of the votes cast, 5% were declared invalid.
- The winner received 55% of the valid votes and won by 1,520 votes.

How many registered voters are there? (Bases: "did not vote" is a % of registered voters, "invalid" is a % of votes cast, and "55%" is a % of valid votes.)

*Answer is a number — no options.*

**Answer:** 20000

**Why:** There are three nested bases. Valid votes = 80% of 95% of registered = 76% of registered. With only two candidates, the winning margin is (55 − 45) = 10% of valid votes. Set that equal to 1,520 and solve backwards.

**Shortcut / Working:**
- Given: non-voters 20% (of V), invalid 5% (of cast), winner 55% (of valid), margin 1,520. Find: V.
- Valid votes = V × 0.80 × 0.95 = 0.76V.
- Two-candidate margin = (2 × 55 − 100)% of valid = 10% of valid.
- 0.10 × 0.76V = 0.076V = 1,520, so V = 1,520 / 0.076 = **20,000**.
- Verify: cast 16,000; invalid 800; valid 15,200; winner 8,360; loser 6,840; margin 1,520 ✓.

**Concept recap:** Chain the multipliers along the bases (registered → cast → valid). In a two-way contest, margin = (2w − 100)% of valid votes, where w is the winner's %.

**Trap alert:** Taking the 5% invalid as a % of *registered* voters (V × 0.75) gives 21,052.6, which is not even an integer. That is your cue that the base is wrong. Using margin = 55% of valid (forgetting to subtract the loser) gives a wrong answer too.

**Time target:** 60 s.

---

### Q6. [M1 | Pass-mark percentages | Easy | MCQ | OA-Speed | Service]
Ravi scores 30% of the maximum marks and fails by 30 marks. Meera scores 45% of the maximum marks and gets 15 marks more than the pass mark. What is the pass mark?

* A. 300
* B. 120
* C. 135
* D. 90

**Answer:** B

**Why:** The pass mark is the same for both students. The gap between their percentages (15% of the maximum) must therefore equal the gap between their marks (30 + 15 = 45). That gives the maximum marks, and the pass mark follows.

**Shortcut / Working:**
- Given: 30% is 30 marks below pass; 45% is 15 marks above pass. Find: pass mark.
- 45% − 30% = 15% of max = 30 + 15 = 45, so max = 300.
- Pass = 30% of 300 + 30 = 90 + 30 = **120**.
- Check with Meera: 45% of 300 = 135 = 120 + 15 ✓.

**Why the other options fail:**
* A: That is the maximum marks, an intermediate value.
* C: That is Meera's score.
* D: That is Ravi's score.

**Concept recap:** When two candidates are measured against the same threshold, (difference in %) × max = (shortfall + surplus). Use subtraction if both are on the same side of the threshold.

**Trap alert:** Subtracting 30 − 15 = 15 marks instead of adding. One student is below the pass mark and one is above, so the gaps add.

**Time target:** 20 s.

---

### Q7. [M1 | Price drop → extra quantity | Medium | MCQ | OA-Speed | Service] — shortcut required
A 20% fall in the price of rice lets a shopper buy 5 kg more rice for ₹800. What was the original price per kg?

* A. ₹32
* B. ₹36
* C. ₹40
* D. ₹38.40

**Answer:** C

**Why:** After the price falls, the shopper saves 20% of ₹800 on the old quantity. That ₹160 is exactly what buys the extra 5 kg at the *new* price. That gives the new price; divide by 0.8 to get the original.

**Shortcut / Working:**
- Given: price −20%, spend ₹800, extra 5 kg. Find: original price per kg.
- Money freed = 20% × 800 = ₹160.
- New price = 160 / 5 = ₹32 per kg.
- Original price = 32 / 0.8 = **₹40 per kg**.
- Check: ₹800 buys 20 kg at ₹40 and 25 kg at ₹32. 25 − 20 = 5 kg ✓.
- Long method (why the shortcut matters): 800/(0.8p) − 800/p = 5 leads to 1000/p − 800/p = 5, so p = 40. That needs fraction algebra under time pressure, which is slow.

**Why the other options fail:**
* A: That is the *new* price. The question asks for the original.
* B: Divides ₹160 by 5 correctly, then adds 12.5% (32 × 1.125) instead of dividing by 0.8. Wrong reversal.
* D: Reverses the 20% fall with a 20% rise (32 × 1.2). A 20% fall is undone by a 25% rise, i.e. dividing by 0.8.

**Concept recap:** In a price-drop-to-extra-quantity problem, (r% × total spend) / extra quantity = new price, and original = new / (1 − r/100). For a price rise (less quantity), the freed amount gives the *original* price instead.

**Trap alert:** The ₹160 buys the extra rice at the *new* price, so it gives the new price, not the original. Stopping at ₹32 is the classic mistake.

**Time target:** 20 s with the shortcut; 90 s with algebra.

---

### Q8. [M1 | Ratio splitting (fractional ratio) | Easy | MCQ | OA-Speed | Service]
₹3,900 is divided among A, B and C in the ratio 1/2 : 1/3 : 1/4. What is C's share?

* A. ₹1,800
* B. ₹1,733.33
* C. ₹975
* D. ₹900

**Answer:** D

**Why:** A fractional ratio becomes a whole-number ratio when you multiply every term by the LCM of the denominators. The LCM of 2, 3 and 4 is 12, so the ratio is 6 : 4 : 3, which is 13 parts.

**Shortcut / Working:**
- Given: total 3,900; ratio 1/2 : 1/3 : 1/4. Find: C's share.
- Multiply by 12: 6 : 4 : 3, a total of 13 parts.
- 1 part = 3,900 / 13 = 300.
- C = 3 × 300 = **₹900** (A = 1,800, B = 1,200).

**Why the other options fail:**
* A: That is A's share.
* B: Uses the denominators as the ratio (2 : 3 : 4): 4/9 × 3,900. Fractional ratios run inversely to their denominators.
* C: Takes 1/4 of the total directly. The fractions are ratio terms, not shares of the total; they add up to 13/12, not 1.

**Concept recap:** To clear a fractional ratio, multiply every term by the LCM of the denominators. Ratio a : b : c of a total T gives shares aT/(a+b+c), and so on.

**Trap alert:** "1/4 means a quarter of the money" is wrong. The fractions here sum to more than 1, so they cannot be shares of the total.

**Time target:** 15 s.

---

### Q9. [M1 | Chained ratios | Hard | MCQ | OA-Speed | Product] — shortcut required
A : B = 2 : 3, B : C = 4 : 5 and C : D = 6 : 7. ₹5,250 is divided among A, B, C and D in this combined ratio. What is D's share?

* A. ₹1,750
* B. ₹1,500
* C. ₹1,200
* D. ₹2,161.76

**Answer:** A

**Why:** The ratios have to be joined on their common terms (B, then C). The chain-product rule does this in one step without step-by-step LCM work.

**Shortcut / Working:**
- Given: A:B = 2:3, B:C = 4:5, C:D = 6:7, total 5,250. Find: D's share.
- Chain-product rule for A:B = a:b, B:C = c:d, C:D = e:f:
  - A = a·c·e = 2·4·6 = 48
  - B = b·c·e = 3·4·6 = 72
  - C = b·d·e = 3·5·6 = 90
  - D = b·d·f = 3·5·7 = 105
- 48 : 72 : 90 : 105, which divides by 3 to give **16 : 24 : 30 : 35** (sum 105).
- 1 part = 5,250 / 105 = 50, so D = 35 × 50 = **₹1,750**.
- Long method: making B equal (12), then C equal (30), takes 3–4 lines and is error-prone at speed.

**Why the other options fail:**
* B: That is C's share (30 × 50). The right method, but the wrong term was read.
* C: That is B's share (24 × 50).
* D: Joins the ratios naively as 2 : 3 : 5 : 7 (dropping the mismatched middle terms): 7/17 × 5,250. The ratios were never made consistent.

**Concept recap:** Chain rule: the first term takes all the "first" numbers; each later term swaps in one more "second" number. A:B:C:D = ace : bce : bde : bdf. Always reduce before dividing the total.

**Trap alert:** Joining ratios by concatenation (2:3 then 4:5 written as 2:3:5) treats B as both 3 and 4. A shared term must be equalised first.

**Time target:** 30 s with the chain rule; 90 s with step-wise LCM.

---

### Q10. [M1 | Proportion properties | Medium | MSQ | GATE | GATE-only]
Which of the following statements are correct? (Select all that apply.)

* A. The mean proportional between 9 and 16 is 12.
* B. The third proportional to 4 and 6 is 9.
* C. The fourth proportional to 2, 3 and 4 is 8.
* D. If x : y = 5 : 3, then (x + y) : (x − y) = 4 : 1.

**Answer:** A, B, D

**Why:** Each statement tests a standard proportion definition. A, B and D follow directly from their definitions. C confuses the fourth proportional with doubling the last term.

**Shortcut / Working:**
- A: Mean proportional of a and b = √(ab) = √(9 × 16) = √144 = 12 ✓.
- B: Third proportional to a and b is c with a : b = b : c, so c = b²/a = 36/4 = 9 ✓.
- C: Fourth proportional to a, b, c is d with a : b = c : d, so d = bc/a = (3 × 4)/2 = 6, not 8 ✗.
- D: Componendo–dividendo: if x/y = 5/3, then (x + y)/(x − y) = (5 + 3)/(5 − 3) = 8/2 = 4/1 ✓. This requires x ≠ y, which holds here.

**Why the other options fail:**
* C (incorrect statement): 8 comes from seeing 2 → 4 as "×2" and applying it to 4. But the proportion runs a : b = c : d (2 is to 3 as 4 is to ?). The multiplier 3/2 maps a to b, so d = 4 × 3/2 = 6.

**Concept recap:**
- Mean proportional = √(ab).
- Third proportional = b²/a.
- Fourth proportional = bc/a.
- Componendo–dividendo: a/b = c/d ⇒ (a + b)/(a − b) = (c + d)/(c − d), provided a ≠ b.

**Trap alert:** In MSQs, check every option on its own. Getting A right does not make the others right. Watch the direction of a proportion: "a is to b as c is to d".

**Time target:** 40 s.

---

### Q11. [M1 | Ratio change after transfer | Medium | MCQ | OA-Speed | Service]
In a class, the ratio of boys to girls is 7 : 5. If 12 boys leave and 12 girls join, the numbers of boys and girls become equal. How many students were in the class originally?

* A. 120
* B. 132
* C. 144
* D. 156

**Answer:** C

**Why:** Write boys = 7k and girls = 5k. The difference 2k must be closed by a 24-student swing (12 fewer boys and 12 more girls). So 2k = 24, k = 12 and the total is 12k = 144.

**Shortcut / Working:**
- Given: boys:girls = 7:5; boys −12, girls +12, then equal. Find: original total.
- 7k − 12 = 5k + 12, so 2k = 24 and k = 12.
- Total = 7k + 5k = 12k = **144**.
- Check: 84 − 12 = 72 and 60 + 12 = 72 ✓.
- Shortcut: the gap in parts (7 − 5 = 2 parts) equals the total swing (12 + 12 = 24), so 1 part = 12.

**Why the other options fail:**
* A: k = 10 gives 70 − 12 = 58 and 50 + 12 = 62. Not equal.
* B: Not a multiple of 12, so it cannot be 12k at all. Rule it out instantly.
* D: k = 13 gives 79 vs 77. Not equal.

**Concept recap:** When the change is "move/close the gap", (difference in ratio parts) × k = total swing. Here the total does not change (−12 + 12), which also confirms that the answer must be a multiple of 7 + 5 = 12.

**Trap alert:** Using only 12 as the swing (forgetting that both sides move) gives k = 6 and total 72, which is not even among the options.

**Time target:** 15 s.

---

### Q12. [M1 | Combining two ratio mixtures | Medium | MCQ | OA-Speed | Service]
Vessel P contains milk and water in the ratio 3 : 1. Vessel Q contains milk and water in the ratio 5 : 3. 2 litres from P and 3 litres from Q are poured into an empty vessel. What is the ratio of milk to water in the mixture?

* A. 2 : 1
* B. 27 : 13
* C. 11 : 5
* D. 7 : 3

**Answer:** B

**Why:** You cannot add ratio terms directly when the volumes are different. Convert each vessel to the fraction of milk and water per litre, scale by the volume taken, then add the milk amounts and the water amounts separately.

**Shortcut / Working:**
- Given: P = 3:1 (4 parts), Q = 5:3 (8 parts); take 2 L of P and 3 L of Q. Find: milk : water.
- Make the part totals equal: P becomes 6 : 2 (per 8 parts) and Q stays 5 : 3 (per 8 parts).
- Now each "unit" of P or Q is the same size, so weight by 2 and 3:
  - Milk = 2 × 6 + 3 × 5 = 12 + 15 = 27
  - Water = 2 × 2 + 3 × 3 = 4 + 9 = 13
- Ratio = **27 : 13**.
- Litre check: milk = 2(3/4) + 3(5/8) = 1.5 + 1.875 = 3.375; water = 0.5 + 1.125 = 1.625. 3.375 : 1.625 = 27 : 13 ✓.

**Why the other options fail:**
* A: Adds the ratio terms (3 + 5 : 1 + 3 = 8 : 4). This ignores both the volumes and the different part sizes.
* C: Mixes equal volumes (1 : 1) instead of 2 : 3.
* D: Swaps the volumes (3 L of P and 2 L of Q).

**Concept recap:** To combine mixtures, first make the part totals equal across vessels (LCM of the sums). Then weight each vessel's parts by the quantity taken and add component by component.

**Trap alert:** "Add the ratios" is wrong unless the quantities taken are equal *and* the ratio sums are equal.

**Time target:** 30 s.

---

### Q13. [M1 | Income–expenditure–savings ratio | Hard | MCQ | OA | Product]
The monthly incomes of A and B are in the ratio 5 : 4, and their monthly expenditures are in the ratio 3 : 2. Each of them saves ₹1,600 per month. What is A's monthly income?

* A. ₹3,200
* B. ₹2,400
* C. ₹8,000
* D. ₹4,000

**Answer:** D

**Why:** Savings = income − expenditure. Use separate multipliers x (income) and y (expenditure). That gives two equations, and the equal savings let you eliminate y quickly.

**Shortcut / Working:**
- Given: incomes 5x and 4x; expenditures 3y and 2y; each saving 1,600. Find: 5x.
- A: 5x − 3y = 1,600 … (1)
- B: 4x − 2y = 1,600 … (2)
- Shortcut: multiply (2) by 3/2 to match the 3y: 6x − 3y = 2,400 … (2′).
- (2′) − (1): x = 800.
- A's income = 5 × 800 = **₹4,000**.
- Check: y from (2): 3,200 − 2y = 1,600, so y = 800. A spends 2,400 and saves 1,600 ✓. B earns 3,200, spends 1,600 and saves 1,600 ✓.
- Faster shortcut: equal savings means income − expenditure is equal. Cross-multiply the ratios: x = S × (e₁ − e₂)/(i₂e₁ − i₁e₂) = 1,600 × (3 − 2)/(4·3 − 5·2) = 1,600/2 = 800. Precondition: the savings are equal.

**Why the other options fail:**
* A: That is B's income (4x).
* B: That is A's expenditure (3y).
* C: Assumes the saving is 1 part of the income ratio (x = 1,600). The income and expenditure ratios use different multipliers.

**Concept recap:** Use different multipliers for each ratio (x for income, y for expenditure). Savings = income − expenditure. With equal savings S: x = S(e₁ − e₂)/(i₂e₁ − i₁e₂).

**Trap alert:** Using one multiplier k for both ratios (5k − 3k = 2k = 1,600, so k = 800) happens to give ₹4,000 here, but only by coincidence: x and y both turned out to be 800. Change the expenditure ratio to 3 : 1 and the single-k method gives a wrong answer. Income and expenditure are independent quantities, so always give them separate multipliers.

**Time target:** 45 s with the elimination shortcut; 90 s otherwise.

---

### Q14. [M1 | Partnership (capital × time) | Easy | MCQ | OA-Speed | Service]
A invests ₹40,000 for 12 months and B invests ₹60,000 for 8 months in a business. Out of a profit of ₹18,000, what is B's share?

* A. ₹9,000
* B. ₹10,800
* C. ₹7,200
* D. Cannot be determined without knowing when B joined

**Answer:** A

**Why:** Profit is shared by capital × time (rupee-months). A: 40,000 × 12 = 4.8 lakh rupee-months. B: 60,000 × 8 = 4.8 lakh rupee-months. They are equal, so the profit splits equally.

**Shortcut / Working:**
- Given: A 40k × 12; B 60k × 8; profit 18,000. Find: B's share.
- Work in thousands: A = 40 × 12 = 480, B = 60 × 8 = 480. Ratio 1 : 1.
- B = 18,000 / 2 = **₹9,000**.

**Why the other options fail:**
* B: Uses the capital ratio only (2 : 3), so 3/5 × 18,000. This ignores time.
* C: Uses the time ratio only (12 : 8 = 3 : 2), so 2/5 × 18,000. This ignores capital.
* D: The duration (8 months) is all that matters. When B joined does not change the rupee-months.

**Concept recap:** Profit share ∝ capital × time. With equal durations, use the capital ratio. With equal capitals, use the time ratio.

**Trap alert:** Using only one factor is the standard OA trap. Multiply both factors every time.

**Time target:** 10 s.

---

### Q15. [M1 | Partnership with changing capital | Medium | NAT | OA | Service]
A starts a business with ₹30,000. After 4 months, B joins with ₹40,000. At the end of the 6th month from the start, A withdraws ₹10,000. At the end of the year the profit is ₹31,000. What is B's share of the profit (in ₹)?

*Answer is a number — no options.*

**Answer:** 16000

**Why:** Split A's year into periods with different capital and add the rupee-months. B is in for 12 − 4 = 8 months with a constant ₹40,000. Then share the profit in the ratio of the totals.

**Shortcut / Working:**
- Given: A 30k for months 1–6, then 20k for months 7–12; B 40k for months 5–12; profit 31,000. Find: B's share.
- In thousand-rupee-months:
  - A = 30 × 6 + 20 × 6 = 180 + 120 = 300
  - B = 40 × 8 = 320
- Ratio A : B = 300 : 320 = 15 : 16, so 31 parts.
- 1 part = 31,000 / 31 = 1,000, so B = 16 × 1,000 = **₹16,000**.

**Concept recap:** For changing capital, total contribution = Σ (capital in each period × length of period). Drop the thousands to keep the numbers small.

**Trap alert:**
- Giving B 12 months (he joined after 4) gives 480 : 300, which is wrong.
- Applying A's withdrawal from month 4 instead of month 6 gives A = 30 × 4 + 20 × 8 = 280, which is also wrong.
- Draw a 12-month timeline before multiplying.

**Time target:** 40 s.

---

### Q16. [M1 | Working-partner salary + capital share | Medium | MCQ | OA | Service]
A and B are partners with capitals in the ratio 3 : 2. A is the working partner and first takes 10% of the total profit as salary. The remaining profit is divided in the capital ratio. If A receives ₹16,000 in all, what is B's share?

* A. ₹6,400
* B. ₹10,000
* C. ₹9,000
* D. ₹10,666.67

**Answer:** C

**Why:** A's total has two parts: the salary (10% of P) and 3/5 of the remaining 90% of P. Write it as one fraction of P, solve for P, then take B's 2/5 of the remaining 90%.

**Shortcut / Working:**
- Given: salary 10% of P to A; rest split 3 : 2; A's total 16,000. Find: B's share.
- A's fraction = 0.10 + (3/5)(0.90) = 0.10 + 0.54 = 0.64.
- 0.64P = 16,000, so P = 25,000.
- B = (2/5)(0.90)(25,000) = 0.36 × 25,000 = **₹9,000**.
- Check: salary 2,500 + 3/5 × 22,500 = 2,500 + 13,500 = 16,000 ✓. 16,000 + 9,000 = 25,000 ✓.

**Why the other options fail:**
* A: Takes 40% of A's ₹16,000. B's share is a fraction of the profit, not of A's receipt.
* B: Finds P = 25,000 correctly but gives B 2/5 of the full P, forgetting to subtract the salary first.
* D: Ignores the salary altogether (A = 3/5 P gives P = 26,666.67, and B = 2/5 of that).

**Concept recap:** Working-partner problems: deduct the salary or commission first, then share the remainder in the capital(-time) ratio. Express each partner's total as a single fraction of P.

**Trap alert:** The salary is taken out *before* the split. Splitting the full profit by capital and then adding the salary double-counts it.

**Time target:** 40 s.

---

### Q17. [M1 | Partnership — find joining time | Medium | MCQ | OA-Speed | Service]
A starts a business with ₹8,000. B joins some months later with ₹12,000. At the end of the year, the profit is shared A : B = 2 : 1. How many months after the start did B join?

* A. 4
* B. 6
* C. 3
* D. 8

**Answer:** D

**Why:** Work backwards from the profit ratio to rupee-months. A has 8,000 × 12. B's rupee-months must be half of A's, which fixes how long B was in. The question asks *when he joined*, which is 12 minus that duration.

**Shortcut / Working:**
- Given: A 8,000 × 12 months; B 12,000 × t months; profit ratio 2 : 1. Find: 12 − t.
- 8,000 × 12 : 12,000 × t = 2 : 1.
- 96,000 = 2 × 12,000t = 24,000t, so t = 4 months in the business.
- B joined after 12 − 4 = **8 months**.

**Why the other options fail:**
* A: 4 is how long B stayed, not when he joined. The question asks for the joining month.
* B: Assumes B also put in ₹8,000 (8,000 × 12 : 8,000t = 2 : 1 gives t = 6).
* C: Swaps the capitals (12,000 × 12 : 8,000t = 2 : 1 gives t = 9, joining after 3).

**Concept recap:** Reverse partnership: profit ratio = (C₁T₁) : (C₂T₂). Solve for the unknown T, then check whether the question wants the *duration* or the *joining point* (12 − T).

**Trap alert:** Duration vs joining time. Underline "after how many months did B join" before answering.

**Time target:** 20 s.

---

### Q18. [M1 | Partnership with multiple capital changes | Hard | MCQ | OA | Service]
A, B and C start a business with capitals in the ratio 3 : 4 : 5.
- After 4 months, A doubles his capital.
- At the end of 6 months from the start, B withdraws half of his capital.
- C's capital stays the same all year.

The year-end profit is ₹1,30,000. How much more does C receive than B?

* A. ₹10,833
* B. ₹20,000
* C. ₹30,000
* D. ₹0

**Answer:** B

**Why:** Treat the ratio terms as unit capitals (3, 4, 5) and build each partner's capital-months over the 12 months. The final ratio decides the split. The initial 3 : 4 : 5 is only a starting point.

**Shortcut / Working:**
- Given: start 3 : 4 : 5; A ×2 after month 4; B ×½ after month 6; C fixed; profit 1,30,000. Find: C − B.
- Capital-month units:
  - A = 3 × 4 + 6 × 8 = 12 + 48 = 60
  - B = 4 × 6 + 2 × 6 = 24 + 12 = 36
  - C = 5 × 12 = 60
- Ratio 60 : 36 : 60 = 5 : 3 : 5, a total of 13 parts.
- 1 part = 1,30,000 / 13 = 10,000. C = 50,000, B = 30,000.
- C − B = 2 parts = **₹20,000**.
- Shortcut: you only need C − B = (5 − 3) parts, so skip computing individual shares.

**Why the other options fail:**
* A: Uses the initial ratio 3 : 4 : 5 and ignores the changes (1 part of 12 = 10,833).
* C: That is B's share, not the difference.
* D: Notices that A and C end up equal (60 = 60) and confuses A with B.

**Concept recap:** For multi-change partnerships, the ratio terms can act as capitals when everyone starts together. Build per-period products, reduce, then answer only what is asked (a difference needs only the difference in parts).

**Trap alert:** Applying B's withdrawal from month 4 (the other event's month) gives B = 4 × 4 + 2 × 8 = 32. Map each event to its own month on a timeline.

**Time target:** 60 s.

---

### Q19. [M1 | Age problem — two time points | Easy | MCQ | OA-Speed | Service]
Ten years ago, a father was four times as old as his son. Ten years from now, he will be twice as old as his son. What is the present ratio of the father's age to the son's age?

* A. 5 : 2
* B. 4 : 1
* C. 3 : 1
* D. 2 : 1

**Answer:** A

**Why:** Set up one equation for each time point and solve. Present ages: son 20, father 50, so the ratio is 5 : 2.

**Shortcut / Working:**
- Given: f − 10 = 4(s − 10); f + 10 = 2(s + 10). Find: f : s.
- From the first: f = 4s − 30. From the second: f = 2s + 10.
- 4s − 30 = 2s + 10, so s = 20 and f = 50.
- Ratio = 50 : 20 = **5 : 2**.
- Check: 10 years ago 40 and 10 → 4 : 1 ✓. In 10 years 60 and 30 → 2 : 1 ✓.
- Age-difference shortcut: the difference never changes. Then: 4 : 1 means a difference of 3 parts (son = 1 part = a). Later: 2 : 1 means a difference of 1 part (son = b). So 3a = b, and b − a = 20 (20 years pass). That gives a = 10, so the son is 20 now.

**Why the other options fail:**
* B: The ratio 10 years ago, not now.
* C: Averages the two given ratios (4 and 2). Ratios at different times cannot be averaged.
* D: The ratio 10 years from now.

**Concept recap:** In age problems, the age difference is constant. Express each time point's ratio so the difference matches, and the time gap then fixes the scale.

**Trap alert:** Reading "now" vs "then" carelessly. The question asks for the *present* ratio; both given ratios are at other times.

**Time target:** 30 s.

---

### Q20. [M1 | Age ratio — scaling the difference | Medium | NAT | GATE | Product] — shortcut required
The present ages of P and Q are in the ratio 7 : 9. Six years ago, their ages were in the ratio 2 : 3. What is the sum of their present ages (in years)?

*Answer is a number — no options.*

**Answer:** 32

**Why:** The age difference is the same at both times. The present ratio has a difference of 2 parts and the old ratio a difference of 1 part. Scale the old ratio to 4 : 6 so both differences are 2 parts. Now each person's age dropped by 3 parts (7 → 4 and 9 → 6) in 6 years, so 1 part = 2 years.

**Shortcut / Working:**
- Given: now 7 : 9; 6 years ago 2 : 3. Find: sum of present ages.
- Differences: 9 − 7 = 2 parts now; 3 − 2 = 1 part then. Scale "then" by 2: 4 : 6.
- Change in each age: 7 − 4 = 3 parts = 6 years, so 1 part = 2 years.
- Present ages: 14 and 18. Sum = **32**.
- Check: 6 years ago, 8 and 12 → 2 : 3 ✓.
- Algebra: (7k − 6)/(9k − 6) = 2/3 gives 21k − 18 = 18k − 12, so k = 2. Ages 14 and 18.

**Concept recap:** When two age ratios have different part-differences, scale one so the differences match (the difference is invariant). Then the change in parts equals the time gap.

**Trap alert:** Comparing 7 : 9 with 2 : 3 directly (7 − 2 = 5 parts = 6 years) mixes parts of different sizes. Equalise the differences first.

**Time target:** 30 s with difference-scaling; 60 s with algebra.

---

## Batch 1 summary

| Split | This batch |
|---|---|
| MCQ / MSQ / NAT | 16 / 1 / 3 |
| Easy / Medium / Hard | 7 / 9 / 4 |
| OA-Speed / OA / GATE / SDE-Interview | 13 / 5 / 2 / 0 |
| Service / Product / GATE-only | 16 / 3 / 1 |
| MCQ answer key (A / B / C / D) | 4 / 4 / 4 / 4 |

**Traps planted:** increase–decrease asymmetry (Q1, Q3, Q4); additive vs multiplicative successive change (Q2); wrong percentage base (Q4, Q5); new vs original price (Q7); fractional-ratio inversion (Q8); joining ratios without equalising the shared term (Q9); adding ratio terms of mixtures (Q12); one-factor partnership (Q14); duration vs joining point (Q17).

**Completed:** 20 / 360  **Remaining:** 340

**Ledger:** P01 ✓ (7) · P02 ✓ (6) · P03 ✓ (5) · X1 ✓ (2 of 4; 2 more in Batch 2)
