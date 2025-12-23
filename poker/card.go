package poker

import "fmt"

// Suit represents the suit of a playing card
type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	return [...]string{"H", "D", "C", "S"}[s]
}

// Rank represents the rank of a playing card
type Rank int

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	return [...]string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}[r]
}

// Card represents a playing card with a Suit and a Rank
type Card struct {
	Suit Suit
	Rank Rank
}

// NewCard creates a new Card from a string representation (e.g., "H2", "SQ")
func NewCard(s string) (*Card, error) {
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid card format: %s", s)
	}

	var suit Suit
	switch s[0] {
	case 'H':
		suit = Hearts
	case 'D':
		suit = Diamonds
	case 'C':
		suit = Clubs
	case 'S':
		suit = Spades
	default:
		return nil, fmt.Errorf("invalid suit: %c", s[0])
	}

	var rank Rank
	switch s[1] {
	case '2':
		rank = Two
	case '3':
		rank = Three
	case '4':
		rank = Four
	case '5':
		rank = Five
	case '6':
		rank = Six
	case '7':
		rank = Seven
	case '8':
		rank = Eight
	case '9':
		rank = Nine
	case 'T':
		rank = Ten
	case 'J':
		rank = Jack
	case 'Q':
		rank = Queen
	case 'K':
		rank = King
	case 'A':
		rank = Ace
	default:
		return nil, fmt.Errorf("invalid rank: %c", s[1])
	}

	return &Card{Suit: suit, Rank: rank}, nil
}

// String returns the string representation of the card
func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.Suit, c.Rank)
}
