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
## Simply run
```bash
bash <(curl -fsSL https://raw.githubusercontent.com/razzat008/letsgodb/main/script/script.sh)
```
## Project Tree(Expected)
[_auto generated_]
```
letsgodb
├── data
│   └── test                    # Database Name in Create Database
│       ├── catalog.db          # Table schema catalog (JSON lines)
│       └── test_table.db       # Example table data file (rows, binary)
├── go.mod
├── internal
│   ├── catalog                 # Table schema catalog (metadata)
│   │   ├── catalog.go
│   │   └── catalog_test.go
│   ├── db                      # High-level DB API (Insert, Select, Row serialization, validation)
│   │   ├── database.go
│   │   ├── eval.go
│   │   └── validate.go
│   ├── Parser                  # SQL/command parser (AST construction)
│   │   └── parser.go
│   ├── REPl                    # Interactive shell (REPL)
│   │   └── repl.go             
│   ├── storage                 # Core storage layer (pager, disk persistence)
│   │   └── pager.go
│   ├── Tokenizer
│   │   └── tokenizer.go        # Tokenizing SQL-like inputs
│   └── util                    # Utility functions (assert, etc.)
│       └── assert.go
├── main.go                     # Entry point (REPL, statement execution)
├── README.md
└── script                      # Optional setup or debug tools
    └── script.sh
```

## Help
- All commands end with `;`, semicolon
- The `help;` command displays short help message.
- The `helpall;` command displays all supported queries.
- The `\e;` command exits the program.

## Wiki
See [wiki](https://github.com/razzat008/letsgodb/wiki) for more.
