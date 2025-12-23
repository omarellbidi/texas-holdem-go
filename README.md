# Texas Hold'em Hand Evaluator

This project implements a Texas Hold'em hand evaluator in Go, adhering strictly to object-oriented design principles as specified in the assignment. It supports evaluating 5-card hands, finding the best 5-card hand from 7 cards, and comparing hands.

## Prerequisites

- [Go](https://go.dev/dl/) (version 1.16 or later)

## Usage

The application can be run in two modes: **Interactive Mode** and **Argument Mode**.

### 1. Interactive Mode

Run the program without any arguments to enter an interactive shell where you can type hands manually.

```bash
go run cmd/poker/main.go
```

**Example Session:**
```text
Enter poker hands (e.g., 'H2 SQ C2 D2 CQ') or 'quit' to exit:
> H2 SQ C2 D2 CQ
Hand: H2 SQ C2 D2 CQ
Value: Full House
Kickers: [2 12]
> quit
```

### 2. Argument Mode

You can evaluate hands directly by passing them as command-line arguments. Make sure to enclose each hand string in quotes.

```bash
go run cmd/poker/main.go "H2 SQ C2 D2 CQ" "H5 S6 C7 D8 H9"
```

## Running Tests

The project includes a comprehensive test suite based on the assignment's test cases.

To run all tests with verbose output:

```bash
go test -v ./...
```

### Test Coverage

To check the code coverage:

1.  Run the tests and generate a profile:
    ```bash
    go test -coverprofile=coverage.out ./...
    ```
2.  View the function-level coverage report:
    ```bash
    go tool cover -func=coverage.out
    ```

## Project Structure

- **`poker/`**: Contains the core logic.
    - `card.go`: Defines `Card`, `Suit`, and `Rank` types and methods.
    - `hand.go`: Defines `Hand` type, evaluation logic (`Evaluate`, `findBestFive`), and comparison logic (`Compare`).
    - `hand_test.go`: Unit tests using strict assignment data.
- **`cmd/poker/`**: Contains the main entry point (`main.go`).
- **`TEST_REPORT.md`**: Summary of test results and coverage.
