# Eventstore Migrate

This project is a migration tool to transfer events from an LMDB store to a Badger store.

## Prerequisites

- Go 1.23.3 or later

## Installation

Clone the repository:

```sh
git clone <repository-url>
cd eventstore_migrate
```

Install dependencies:

```sh
go mod tidy
```

## Usage

To run the migration tool, use the following command:

```sh
go run main.go --lmdb <path-to-lmdb> --badger <path-to-badger> [--mapsize <map-size>]
```

### Arguments

- `--lmdb`: Path to the LMDB store directory (required)
- `--badger`: Path to the Badger store directory (required)
- `--mapsize`: Map size for the LMDB store (optional, default is 1<<34)

### Example

```sh
go run main.go --lmdb /path/to/lmdb --badger /path/to/badger --mapsize 1<<34
```

## Code Overview

The main functionality is implemented in `main.go`. The key steps are:

1. Parse command-line arguments.
2. Initialize the LMDB store.
3. Initialize the Badger store.
4. Query all events from the LMDB store.
5. Save each event to the Badger store.
6. Close both stores.

## Dependencies

- `github.com/fiatjaf/eventstore`
- `github.com/nbd-wtf/go-nostr`
- `github.com/PowerDNS/lmdb-go`
- `github.com/dgraph-io/badger/v4`

## License

This project is licensed under Girino's Anarchist License. See the [license](https://license.girino.org) for details.