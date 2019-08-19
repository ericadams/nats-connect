package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
)

// printHelp will print out the flag options for the server.
func printHelp() {
	var usageStr = `
Usage: nats-connect [options]
Server Options:
	-p, --port <port>                Port to listen on (default: 5120)
	-c, --conn <connection_string>   Connection string to the database

Common Options:
	-h, --help                       Show this message
	-v, --version                    Show version
`
	fmt.Printf("%s\n", usageStr)
}

// printVersion will print out the version for the server.
func printVersion() {
	fmt.Printf("Version: %s\n", version)
}

func parseFlags() (*option, error) {
	var (
		showHelp    bool
		showVersion bool
	)

	opts := &option{
		Port: 5120,
	}

	// Create a FlagSet and set the usage.
	fs := flag.NewFlagSet("nats-connect", flag.ExitOnError)
	fs.Usage = printHelp

	fs.BoolVar(&showHelp, "h", false, "Show this message.")
	fs.BoolVar(&showHelp, "help", false, "Show this message.")
	fs.BoolVar(&showVersion, "v", false, "Show version.")
	fs.BoolVar(&showVersion, "version", false, "Show version.")
	fs.IntVar(&opts.Port, "p", 5120, "Port to listen on.")
	fs.IntVar(&opts.Port, "port", 5120, "Port to listen on.")
	fs.StringVar(&opts.Connection, "c", "", "Connection string to the database.")
	fs.StringVar(&opts.Connection, "conn", "", "Connection string to the database.")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if showVersion {
		printVersion()
		return nil, nil
	}

	if showHelp {
		printHelp()
		return nil, nil
	}

	return opts, nil
}

func main() {
	log.SetOutput(os.Stderr)
	opts, err := parseFlags()
	if err != nil {
		log.Fatalf("could not parse flags: %s", err.Error())
	}

	conn, err := sql.Open("mssql", opts.Connection)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}

	srv := newServer(opts, conn)
	log.Fatal(srv.Run())
}
