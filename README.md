# letsgodb
A database built from scratch in Go.

# Running the project
Install go
```bash
sudo pacman -S go
```
Clone the repository
```bash
git clone https://github.com/razzat008/letsgodb.git
cd letsgodb
```
Run the project
```bash
go run .
```
## Project Tree(Expected)
[_auto generated_]
```
letsgodb/
├── main.go                   # Entry point (REPL, statement execution)
├── internal/
│   ├── repl/                 # Interactive shell
│   │   └── repl.go
│   ├── tokenizer/            # Tokenizing SQL-like inputs
│   │   └── tokenizer.go
│   ├── parser/               # SQL/command parser (AST construction)
│   │   └── parser.go
│   ├── storage/              # Core storage layer (pager, disk persistence)
│   │   └── pager.go
│   ├── db/                   # High-level DB API (Insert, Select, Row serialization)
│   │   └── database.go
│   ├── catalog/              # Table schema catalog (metadata)
│   │   ├── catalog.go
│   │   └── catalog_test.go
│   ├── btree/                # B+Tree logic (indexing, fast lookups)
│   │   ├── btree.go
│   │   └── node.go
│   └── util/                 # Utility functions (assert, etc.)
│       └── assert.go
├── scripts/                  # Optional setup or debug tools
│   └── init.sh
├── go.mod
├── go.sum
├── catalog.db                # Table schema catalog (JSON lines)
├── users.db                  # Example table data file (rows, binary)
└── README.md
```

## Help
- All commands end with `;`, semicolon
- The `help` command displays this message.
- The `\e` command exits the program.

## Wiki
See [wiki](https://github.com/razzat008/letsgodb/wiki) for more.
