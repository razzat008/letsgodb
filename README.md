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
├── main.go                   # Entry point (e.g., REPL)
├── internal/                 # Internal components (not exposed as APIs)
│   ├── repl/                 # Interactive shell
│   │   └── repl.go
│   ├── tokenizer/            # Tokenizing SQL-like inputs
│   │   └── tokenizer.go
│   ├── parser/               # SQL/command parser
│   │   └── parser.go
│   ├── ast/                  # Abstract Syntax Tree node definitions
│   │   └── node.go
│   ├── vm/                   # Bytecode execution engine (VM)
│   │   └── vm.go
│   ├── btree/                # Immutable B+Tree logic
│   │   ├── btree.go
│   │   └── node.go
│   ├── storage/              # Core storage layer (disk persistence)
│   ├── db/                   # High-level DB API (Insert, Delete, Query)
│   │   └── database.go
│   └── util/                 # Utility functions (assert, logging, etc.)
│       └── assert.go
├── testdata/                 # Sample files and testing fixtures
│   └── example.db
├── scripts/                  # Optional setup or debug tools
│   └── init.sh
├── go.mod
├── go.sum
└── README.md
``` 
