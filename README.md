# GoSQLite

A SQLite database engine implementation written in Go, built from scratch to understand the internals of SQLite's file format and B-tree storage.

## Overview

GoSQLite is an educational project that implements core SQLite functionality in Go. It provides low-level access to SQLite database files, allowing you to read and parse the binary format directly.

## Features

### Currently Implemented

- **File Operations**: Memory-mapped file access for efficient page reading
- **Database Header Parsing**: Reads SQLite database metadata including page size, encoding, and schema information
- **B-tree Implementation**: Basic B-tree leaf page parsing and navigation
- **Cell Processing**: Extraction and parsing of database cells containing row data
- **Varint Processing**: SQLite's variable-length integer format encoding/decoding
- **Record Parsing**: Structured data extraction supporting multiple data types:
  - Integers (8, 16, 24, 32, 48, 64-bit)
  - Floating-point numbers
  - Text and BLOB data
  - NULL values

### Architecture

```
gosqlite/
├── main_db.go          # Entry point and usage example
└── db/
    ├── pager.go        # File I/O and page management
    ├── btree.go        # B-tree data structures and operations
    ├── page_reader.go  # Database header and page parsing
    └── bit_processing.go # Varint and binary data processing
```

## Usage

```go
package main

import "github.com/rishabhBudhouliya/gosqlite/db"

func main() {
    // Initialize pager and open database file
    pager := db.Pager{}
    pager.Open("sample.db")
    
    // Read the first page (contains database header)
    result, _ := pager.GetPage(1, 4096)
    
    // Parse the page content
    db.Read(1, result)
}
```

## Current Limitations

This implementation is in early development and currently supports:
- Reading SQLite database files
- Parsing database headers and metadata
- B-tree leaf page processing
- Basic record extraction
What it doesn't support:
- SQL query parsing (not implemented)
- Writing/modifying databases (read-only)
- Index operations
- Transaction support

## Building and Running

```bash
# Clone the repository
git clone https://github.com/rishabhBudhouliya/gosqlite.git
cd gosqlite

# Run the main program
go run main_db.go

# Run tests
go test ./db/
```

## Dependencies

- Go 1.23.2+
- `golang.org/x/exp` for memory mapping support

## SQLite Format Understanding

This implementation follows the [SQLite file format specification](https://www.sqlite.org/fileformat.html) and includes:

- **Page Structure**: 4096-byte pages with headers and cell arrays
- **Varint Encoding**: Variable-length integer compression
- **B-tree Organization**: Hierarchical data storage for efficient access
- **Record Format**: Row data serialization with type information

## Contributing

This is an educational project focused on understanding SQLite internals. Contributions that improve the understanding of the SQLite format or add well-documented features are welcome.

## License

MIT License - see LICENSE file for details.

## References

- [SQLite File Format Documentation](https://www.sqlite.org/fileformat.html)
- [SQLite B-tree Module](https://www.sqlite.org/btreemodule.html)
- [Database Internals by Alex Petrov](https://databass.dev/)
