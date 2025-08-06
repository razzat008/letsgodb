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
letsgodb/
├── main.go                   # Entry point (REPL, statement execution)
├── internal/
│   ├── REPl/                 # Interactive shell (REPL)
│   │   └── repl.go
│   ├── Tokenizer/            # Tokenizing SQL-like inputs
│   │   └── tokenizer.go
│   ├── Parser/               # SQL/command parser (AST construction)
│   │   └── parser.go
│   ├── storage/              # Core storage layer (pager, disk persistence)
│   │   └── pager.go
│   ├── db/                   # High-level DB API (Insert, Select, Row serialization, validation)
│   │   ├── database.go
│   │   ├── eval.go
│   │   └── validate.go
│   ├── catalog/              # Table schema catalog (metadata)
│   │   ├── catalog.go
│   │   └── catalog_test.go
│   └── btree/                # B+Tree logic (indexing, fast lookups)
│       ├── btree.go
│       └── node.go
├── script/                   # Setup or debug scripts
│   └── script.sh
├── data/                     # Persistent database files
│   └── test/
│       ├── catalog.db
│       └── users.db
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── README.md                 # Project documentation
```

## Help
- All commands end with `;`, semicolon
- The `help;` command displays short help message.
- The `helpall;` command displays all supported queries.
- The `\e;` command exits the program.

## Wiki
See [wiki](https://github.com/razzat008/letsgodb/wiki) for more.
