package poker

import (
	"fmt"
	"sort"
	"strings"
)

// HandVal represents the value of a poker hand
type HandVal int

const (
	HighCard HandVal = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (hv HandVal) String() string {
	return [...]string{
		"High Card",
		"One Pair",
		"Two Pair",
		"Three of a Kind",
		"Straight",
		"Flush",
		"Full House",
		"Four of a Kind",
		"Straight Flush",
		"Royal Flush",
	}[hv]
}

// Hand represents a collection of cards (5 to 7)
type Hand struct {
	Cards []*Card
}

// NewHand creates a new Hand from a space-separated string of cards
func NewHand(s string) (*Hand, error) {
	parts := strings.Fields(s)
	// Relaxed constraint slightly to allow creating hands for testing, 
	// but keeping the core logic consistent. 
	// The assignment implies 5 or 7 cards usually.
	if len(parts) < 5 || len(parts) > 7 {
		return nil, fmt.Errorf("hand must have 5-7 cards, got: %d", len(parts))
	}

	hand := &Hand{Cards: make([]*Card, 0, len(parts))}
	seen := make(map[string]bool)

	for _, part := range parts {
		card, err := NewCard(part)
		if err != nil {
			return nil, err
		}

		key := card.String()
		if seen[key] {
			return nil, fmt.Errorf("duplicate card in hand: %s", key)
		}
		seen[key] = true

		hand.Cards = append(hand.Cards, card)
	}

	return hand, nil
}

// Evaluate calculates the value of the hand
func (h *Hand) Evaluate() (HandVal, []Rank) {
	if len(h.Cards) == 5 {
		return h.evaluateFive()
	}
	if len(h.Cards) == 7 {
		return h.findBestFive()
	}
	// Fallback/Error case, though NewHand restricts size
	return HighCard, nil
}

// evaluateFive evaluates a standard 5-card hand
func (h *Hand) evaluateFive() (HandVal, []Rank) {
	// Create a slice of ranks for easier processing
	ranks := make([]Rank, len(h.Cards))
	for i, card := range h.Cards {
		ranks[i] = card.Rank
	}
	// Sort ranks descending
	sort.Slice(ranks, func(i, j int) bool { return ranks[i] > ranks[j] })

	isFlush := h.isFlush()
	isStraight, straightHighCard := h.isStraight()

	if isFlush && isStraight {
		if straightHighCard == Ace {
			return RoyalFlush, []Rank{Ace}
		}
		return StraightFlush, []Rank{straightHighCard}
	}

	rankCounts := h.getRankCounts()
	// Sort counts to easily identify 4-of-a-kind, Full House, etc.
	counts := make([]int, 0, len(rankCounts))
	for _, count := range rankCounts {
		counts = append(counts, count)
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })

	if counts[0] == 4 {
		return FourOfAKind, h.getKickers(rankCounts, 4)
	}

	if counts[0] == 3 && len(counts) > 1 && counts[1] == 2 {
		return FullHouse, h.getKickers(rankCounts, 3, 2)
	}

	if isFlush {
		return Flush, ranks
	}

	if isStraight {
		return Straight, []Rank{straightHighCard}
	}

	if counts[0] == 3 {
		return ThreeOfAKind, h.getKickers(rankCounts, 3)
	}

	if counts[0] == 2 && len(counts) > 1 && counts[1] == 2 {
		return TwoPair, h.getKickers(rankCounts, 2, 2)
	}

	if counts[0] == 2 {
		return OnePair, h.getKickers(rankCounts, 2)
	}

	return HighCard, ranks
}

// findBestFive finds the best 5-card combination from 7 cards
func (h *Hand) findBestFive() (HandVal, []Rank) {
	bestVal := HighCard
	var bestKickers []Rank

	// Generate all combinations of 5 indices from 7
	// We use a simple iterative approach for 7C5
	indices := []int{0, 1, 2, 3, 4}
	n := 7
	k := 5

	for {
		// Construct the 5-card sub-hand
		subHand := &Hand{Cards: make([]*Card, 5)}
		for i, idx := range indices {
			subHand.Cards[i] = h.Cards[idx]
		}

		val, kickers := subHand.evaluateFive()
		
		// Update best if current is better
		if val > bestVal {
			bestVal = val
			bestKickers = kickers
		} else if val == bestVal {
			if compareKickers(kickers, bestKickers) > 0 {
				bestKickers = kickers
			}
		}

		// Generate next combination
		// Find the rightmost index that can be incremented
		i := k - 1
		for i >= 0 && indices[i] == n-k+i {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < k; j++ {
			indices[j] = indices[j-1] + 1
		}
	}

	return bestVal, bestKickers
}

func (h *Hand) isFlush() bool {
	suit := h.Cards[0].Suit
	for _, card := range h.Cards[1:] {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

func (h *Hand) isStraight() (bool, Rank) {
	ranks := make([]Rank, len(h.Cards))
	for i, card := range h.Cards {
		ranks[i] = card.Rank
	}
	// Sort ascending for straight check
	sort.Slice(ranks, func(i, j int) bool { return ranks[i] < ranks[j] })

	// Check for standard straight
	isStandard := true
	for i := 0; i < len(ranks)-1; i++ {
		if ranks[i+1] != ranks[i]+1 {
			isStandard = false
			break
		}
	}
	if isStandard {
		return true, ranks[len(ranks)-1]
	}

	// Check for Wheel (A, 2, 3, 4, 5)
	// In our enum: 2=0 ... A=12
	if ranks[4] == Ace && ranks[0] == Two && ranks[1] == Three && ranks[2] == Four && ranks[3] == Five {
		return true, Five
	}

	return false, 0
}

func (h *Hand) getRankCounts() map[Rank]int {
	counts := make(map[Rank]int)
	for _, card := range h.Cards {
		counts[card.Rank]++
	}
	return counts
}

// getKickers returns ranks matching the target counts, followed by remaining high cards
func (h *Hand) getKickers(rankCounts map[Rank]int, targetCounts ...int) []Rank {
	var result []Rank

	// 1. Add ranks that match the target counts (e.g., the Pair rank)
	for _, target := range targetCounts {
		// Iterate high to low to ensure best kickers come first
		for r := Ace; r >= Two; r-- {
			if count, ok := rankCounts[r]; ok && count == target {
				// Add it as many times as it appears? 
				// For evaluation comparison, we usually just need the rank value once for the set 
				// (e.g. "Pair of Kings" -> K). But typically we might compare "Pair of Kings" vs "Pair of Kings".
				// Then we check kickers.
				// However, standard comparison usually flattens this: 
				// Full House 3xK, 2x5 -> [K, 5] (or K, K, K, 5, 5 depending on logic).
				// Let's stick to adding the rank *once* for the structural part if implied, 
				// but let's see how Compare works.
				// Compare uses lexicographical comparison of this slice.
				// So for Full House K,K,K,5,5 we want [K, 5].
				// For Two Pair K,K,8,8,2 we want [K, 8, 2].
				// For One Pair K,K,A,9,4 we want [K, A, 9, 4].
				
				// Wait, the friend's code did:
				// for i := 0; i < rankCounts[rank]; i++ { result = append(result, rank) } ? 
				// No, friend's code: "result = append(result, rank)" inside the match loop.
				// But for the "rest" (kickers), it appended all occurrences?
				// Friend's code:
				// targetCounts loop: appends rank once.
				// rest loop: appends rank * count times.
				
				// Let's analyze Two Pair: K K 8 8 2.
				// Target 2, 2.
				// First target 2: finds K -> append K. finds 8 -> append 8. Result: [K, 8].
				// Rest loop: finds 2 (count 1). -> append 2. Result: [K, 8, 2].
				// This works for comparison.
				
				// What about Full House? K K K 5 5.
				// Target 3, 2.
				// First target 3: finds K. -> append K.
				// Second target 2: finds 5. -> append 5.
				// Result: [K, 5]. Correct.
				
				result = append(result, r)
				// We consume this rank so we don't add it again in the "rest" loop
				// But we can't easily modify the map while iterating or we need a 'used' set.
				// Friend checked "!contains(result, rank)".
			}
		}
	}

	// 2. Add remaining cards (real kickers)
	for r := Ace; r >= Two; r-- {
		if count, ok := rankCounts[r]; ok && count > 0 {
			if !contains(result, r) {
				for i := 0; i < count; i++ {
					result = append(result, r)
				}
			}
		}
	}
	return result
}

func contains(slice []Rank, r Rank) bool {
	for _, v := range slice {
		if v == r {
			return true
		}
	}
	return false
}

// Compare compares this hand with another. Returns 1 if this > other, -1 if this < other, 0 if equal.
func (h *Hand) Compare(other *Hand) int {
	val1, kickers1 := h.Evaluate()
	val2, kickers2 := other.Evaluate()

	if val1 > val2 {
		return 1
	}
	if val1 < val2 {
		return -1
	}
	return compareKickers(kickers1, kickers2)
}

func compareKickers(k1, k2 []Rank) int {
	minLen := len(k1)
	if len(k2) < minLen {
		minLen = len(k2)
	}
	for i := 0; i < minLen; i++ {
		if k1[i] > k2[i] {
			return 1
		}
		if k1[i] < k2[i] {
			return -1
		}
	}
	if len(k1) > len(k2) {
		return 1
	}
	if len(k1) < len(k2) {
		return -1
	}
	return 0
}
