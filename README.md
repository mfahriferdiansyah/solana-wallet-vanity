# Solana Wallet Generator

A high-performance CLI tool to generate Solana wallets with public keys that match a specified prefix. This tool uses Go's concurrency features for fast parallel processing and writes the matched results (public and private keys) to a `result.txt` file.

## Features
- Generates Solana wallet keypairs (public and private keys).
- Matches public keys against a specified prefix.
- Supports parallel processing with configurable worker threads.
- Saves matched results to a `result.txt` file for later use.

## Requirements
- [Go (Golang)](https://golang.org/dl/) version 1.19 or later.
- Internet connection to install the necessary dependencies.

## Installation

1. Clone the repository:
    ```
    git clone https://github.com/your-username/solana-wallet-generator.git
    cd solana-wallet-generator
    ```

2. Install dependencies:
    ```
    go mod tidy
    ```

3. Build the executable:
    ```
    go build -o solana-wallet-gen
    ```

## Usage

Run the program with the desired options:

```
./solana-wallet-gen --prefix <desired-prefix> [--workers <number-of-workers>] [--attempts <max-attempts>]
```

### Options
- `--prefix` (required): The prefix the generated public key should start with.
- `--workers` (optional): Number of concurrent workers (default: 4).
- `--attempts` (optional): Maximum number of attempts per worker (default: unlimited).

### Example
To generate a wallet with a public key starting with `sol`, using 8 workers and up to 1,000,000 attempts per worker:
```
./solana-wallet-gen --prefix "sol" --workers 8 --attempts 1000000
```

### Output
When a match is found:
1. The public and private keys are printed to the terminal.
2. The matched results are saved in the `result.txt` file in the following format:
    ```
    Public Key: solXYZ12345abc...
    Private Key: 5TxUHoN123ABCxyz...
    ```
