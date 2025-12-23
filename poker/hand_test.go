package poker

import (
	"fmt"
	"testing"
)

func TestNewHand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid 5 cards", "H2 SQ C2 D2 CQ", false},
		{"Valid 7 cards", "H2 SQ C2 D2 CQ HK SA", false},
		{"Too few cards", "H2 SQ C2", true},
		{"Too many cards", "H2 SQ C2 D2 CQ HK SA H3", true},
		{"Invalid card", "H2 XX C2 D2 CQ", true},
		{"Duplicate card", "H2 H2 SQ C2 D2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandEvaluate_ManualCases(t *testing.T) {
	tests := []struct {
		input string
		want  HandVal
	}{
		// Cases from Hands to be tested.txt
		{"CT CJ CQ CK CA", RoyalFlush},
		{"D8 DQ DJ DT D9", StraightFlush}, // 8-Q Straight Flush
		{"H8 HT HJ H7 H9", StraightFlush}, // 7-J Straight Flush
		{"HT SQ ST DT CT", FourOfAKind},
		{"HT SK ST DT CT", FourOfAKind},
		{"H2 SQ C2 D2 CQ", FullHouse},
		{"H2 SJ C2 D2 CJ", FullHouse},
		{"HK HQ H2 H4 H5", Flush},
		{"D5 D4 D2 DQ DK", Flush},
		{"H3 S7 H5 D6 H4", Straight},
		{"C9 CT SJ D7 H8", Straight},
		{"H4 S5 HA D3 H2", Straight}, // Wheel
		{"H2 SQ S2 D2 CK", ThreeOfAKind},
		{"H5 SQ C5 DT CT", TwoPair},
		{"H9 SQ C9 DT CT", TwoPair},
		{"H3 S8 H5 D8 CA", OnePair},
		{"S4 DA H3 CA HT", OnePair},
		{"H3 S8 H5 DK CA", HighCard},
		{"H3 S8 H5 DK CT", HighCard},
		{"H3 S8 H5 DK C2", HighCard},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			hand, err := NewHand(tt.input)
			if err != nil {
				t.Fatalf("NewHand() failed: %v", err)
			}
			got, _ := hand.Evaluate()
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandCompare_ManualCases(t *testing.T) {
	// 1 = Hand1 wins, -1 = Hand2 wins, 0 = Tie
	tests := []struct {
		h1   string
		h2   string
		want int
	}{
		// Straight Flush Q high vs J high
		{"D8 DQ DJ DT D9", "H8 HT HJ H7 H9", 1},
		// Four of a kind T (K kicker) vs T (Q kicker) ? 
		// "HT SQ ST DT CT" -> T T T T Q
		// "HT SK ST DT CT" -> T T T T K
		{"HT SQ ST DT CT", "HT SK ST DT CT", -1},
		
		// "H8 SQ S8 D8 C8" -> 8 8 8 8 Q
		// "H7 SK S7 D7 C7" -> 7 7 7 7 K
		// 8888 beats 7777
		{"H8 SQ S8 D8 C8", "H7 SK S7 D7 C7", 1},

		// Three of a kind comparisons
		// "H2 S7 S2 D2 C9" -> 2 2 2 9 7
		// "H2 S8 S2 D2 C9" -> 2 2 2 9 8
		// Hand 2 has better kicker (8 vs 7)
		{"H2 S7 S2 D2 C9", "H2 S8 S2 D2 C9", -1},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			hand1, err := NewHand(tt.h1)
			if err != nil {
				t.Fatalf("NewHand(h1) failed: %v", err)
			}
			hand2, err := NewHand(tt.h2)
			if err != nil {
				t.Fatalf("NewHand(h2) failed: %v", err)
			}

			if got := hand1.Compare(hand2); got != tt.want {
				t.Errorf("Compare() = %v, want %v\n%s\nvs\n%s", got, tt.want, tt.h1, tt.h2)
			}
		})
	}
}

func TestFindBestFive(t *testing.T) {
	// 7 Cards
	// H2 H5 H7 H9 HK (Flush) + S3 C4
	hand, err := NewHand("H2 H5 H7 H9 HK S3 C4")
	if err != nil {
		t.Fatalf("NewHand() failed: %v", err)
	}
	val, _ := hand.Evaluate()
	if val != Flush {
		t.Errorf("Evaluate(7) = %v, want Flush", val)
	}

	// 7 Cards with Quads
	// 2 2 2 2 K + 3 4
	hand2, err := NewHand("H2 S2 C2 D2 HK S3 C4")
	if err != nil {
		t.Fatalf("NewHand() failed: %v", err)
	}
	val2, _ := hand2.Evaluate()
	if val2 != FourOfAKind {
		t.Errorf("Evaluate(7) = %v, want FourOfAKind", val2)
	}
}

