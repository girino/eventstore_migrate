package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fiatjaf/eventstore/badger"
	"github.com/fiatjaf/eventstore/lmdb"
	"github.com/nbd-wtf/go-nostr"
)

const MAXLIMIT = 10_000_000

func main() {
	lmdbPath := flag.String("lmdb", "", "Path to the LMDB store directory")
	badgerPath := flag.String("badger", "", "Path to the Badger store directory")
	mapSize := flag.Int64("mapsize", 1<<34, "Map size for the LMDB store (optional)")
	flag.Parse()

	// Check that we got both paths
	if *lmdbPath == "" || *badgerPath == "" {
		fmt.Fprintln(os.Stderr, "Please specify both --lmdb and --badger paths.")
		os.Exit(1)
	}

	ctx := context.Background()

	// Initialize the LMDB store (source)
	lmdbStore := &lmdb.LMDBBackend{
		Path:     *lmdbPath,
		MaxLimit: MAXLIMIT, // adjust as needed
		MapSize:  *mapSize,
	}
	if err := lmdbStore.Init(); err != nil {
		log.Fatalf("failed to initialize LMDB store at '%s': %v", *lmdbPath, err)
	}

	// Initialize the Badger store (destination)
	badgerStore := &badger.BadgerBackend{
		Path:     *badgerPath,
		MaxLimit: MAXLIMIT, // adjust as needed
	}
	if err := badgerStore.Init(); err != nil {
		log.Fatalf("failed to initialize Badger store at '%s': %v", *badgerPath, err)
	}

	// Query all events from LMDB
	ch, err := lmdbStore.QueryEvents(ctx, nostr.Filter{})
	if err != nil {
		log.Fatalf("error querying events from LMDB: %v", err)
	}

	var count int
	for evt := range ch {
		if err := badgerStore.SaveEvent(ctx, evt); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save event '%s' to Badger: %v\n", evt.ID, err)
			os.Exit(123)
		}
		count++
		if count%10_000 == 0 {
			log.Printf("Migrated %d events so far...", count)
		}
	}

	// Close LMDB store
	lmdbStore.Close()

	// Close Badger store
	badgerStore.Close()

	fmt.Printf("Migration completed successfully. %d events moved.\n", count)
}
