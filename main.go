package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	// This is your engine's package
	"github.com/rishabhBudhouliya/gosqlite/db"
)

const (
	TABLE_ROOT_ID = 2
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <database_file>", os.Args[0])
	}
	dbPath := os.Args[1]

	myFlagSet := flag.NewFlagSet("myApp", flag.ExitOnError)

	// 3. Define your flags on this new flag set, NOT on the global one.
	cpuprofile := myFlagSet.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := myFlagSet.String("memprofile", "", "write memory profile to `file`")

	// 4. Parse the flags from the subset of arguments.
	myFlagSet.Parse(os.Args[2:])

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // Make sure to close the file
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		// The defer call ensures that the profile is stopped correctly before the program exits.
		defer pprof.StopCPUProfile()
	}

	// Open your database engine
	myDb := db.NewDatabase(dbPath)
	// myDB, err := db.Open(dbPath)
	// defer myDB.Close()

	// Read rowids from standard input, one per line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		rowID, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Printf("Could not parse line '%s', skipping.", line)
			continue
		}

		// Perform the lookup using your engine
		_ = db.Search(TABLE_ROOT_ID, &myDb.Pager, rowID)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // Make sure to close the file
		// WriteHeapProfile writes a snapshot of the memory profile.
		// It's typically done at the end of the program.
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}
