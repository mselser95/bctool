# 🧠 bctool — Blockchain CLI Utility

`bctool` is a powerful and extensible command-line tool built in Go that fetches and displays blockchain-related data. It provides quick access to token prices, address inspection, and network metadata by leveraging real-time APIs.

---

## 📦 Features

- ✅ Fetch **live token prices** from Gate.io via the `prices` command
- ⚡ Fast: concurrent price fetching for low-latency output
- 🖥️ Pretty tabular CLI output
- 🔌 Mockable, testable HTTP client interface
- 🧪 Unit tested with error cases and retry scenarios
- 🧱 Built to be easily extended with new blockchain commands

---

## 📦 Installation

To install the bctool CLI via Homebrew:

```
brew tap mselser95/bctool
brew install bctool
```

## 💡 Example Usage

```
$ bctool prices BTC ETH LINK TRC UNI
```

### Output

```
+--------+------------+
| TICKER | PRICE USD  |
+--------+------------+
| ETH    | $1,589.34  |
| LINK   | $12.57     |
| BTC    | $84,483.00 |
| TRC    | $0.001580  |
| UNI    | $5.18      |
+--------+------------+
```

---

## 🛠️ How It Works

### Internals

- `cmd/prices.go`: Cobra subcommand that accepts tickers, calls internal client, renders result.
- `internal/prices/client.go`: Uses Gate.io API to fetch prices concurrently with an injectable HTTP client.
- `internal/render/render.go`: Pretty prints the result using `go-pretty/table` and formatted money.

### Concurrency

`prices.GetPrices()` uses a `sync.WaitGroup` to perform all API calls concurrently. Errors are captured through a channel and the first error encountered is returned.

---

## 🚀 Running the CLI

Make sure Go 1.20+ is installed.

```
git clone https://github.com/yourusername/bctool.git
cd bctool
go build -o bctool ./cmd
./bctool prices BTC ETH
```

---

## ✅ Running Tests

All logic is designed for testability with mocked HTTP clients and table-driven test cases.

```
go test ./internal/prices -v
```

---

## 📁 Directory Structure

```
bctool/
├── cmd/                    # CLI commands (Cobra)
│   └── prices.go           # 'prices' subcommand
├── internal/
│   ├── prices/             # Gate.io HTTP client
│   │   └── client.go
│   └── render/             # Pretty-print table output
│       └── render.go
├── main.go                 # Entry point
└── go.mod / go.sum         # Go modules
```

---

## 📬 Roadmap

- [ ] Add `inspect-address` command for on-chain address data
- [ ] Support alternate APIs (CoinGecko, Binance)
- [ ] Add global config / caching
- [ ] JSON and CSV output options

---

