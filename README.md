# Pokédex CLI

A **command-line Pokémon explorer and catcher** written in Go. Explore locations, encounter Pokémon, catch them, and maintain your personal Pokédex!

---

## Features

- Browse Pokémon **locations** with `map` and `bmap`.  
- **Explore a location** to see which Pokémon are present.  
- **Catch Pokémon** with a probability based on base experience.  
- Inspect Pokémon details (`ID`, height, weight, types, base experience).  
- Maintain a **personal Pokédex** of caught Pokémon.  
- Fully interactive **command-line interface**.

---

## Installation

1. Make sure you have **Go installed** (>= 1.20 recommended).  
2. Clone the repository:

```bash
git clone https://github.com/amiroxx12/pokedex.git
cd pokedex

```bash
go build -o pokedex
```

---

## Usage

Start the CLI by running:

```bash
./pokedex
```

You'll enter an interactive prompt. Type commands to explore locations, catch Pokémon, and view your Pokédex.

---

## Commands

| Command         | Description                                                |
|-----------------|-----------------------------------------------------------|
| `help`          | Show available commands                                   |
| `exit`          | Exit the Pokédex CLI                                      |
| `map`           | List Pokémon locations (next 20)                          |
| `bmap`          | List previous 20 Pokémon locations                        |
| `explore <loc>` | List wild Pokémon at the given location                   |
| `catch <name>`  | Attempt to catch a wild Pokémon by name                   |
| `inspect <name>`| Show details about a caught Pokémon                       |
| `pokedex`       | List all Pokémon you've caught                            |

---

## Examples

List locations:

```bash
map
```

Explore a location:

```bash
explore viridian-forest
```

Catch a Pokémon:

```bash
catch caterpie
```

View your Pokédex:

```bash
pokedex
```

Inspect a Pokémon:

```bash
inspect pikachu
```

---

## Catch Probability

The probability of catching a Pokémon is determined by its **base experience**:

```
Catch Probability = max(5%, 100% - (Base Experience / factor))
```

For example, a Pokémon with base experience 64:

```
100% - (64 / factor) 
```
And then a random number is generated between 0-100, if a the random number happens to be less than the catch probability, the Pokemon is successfully catched
The factor makes it either harder or easier to catch the pokemon, the higher the factor, the harder it is to catch the pokemon.
You always have at least a 5% chance to catch any Pokémon.

---

## Project Structure

```
.
├── main.go
├── internal/
│   ├── pokecache/
│   ├── pokecmd/
│   └── pokeapi/
├── go.mod
├── go.sum
└── README.md
```

---

## Dependencies

This project uses:

- [Go standard library](https://golang.org/pkg/)
- [pokeapi.co](https://pokeapi.co/) — for Pokémon data

No external dependencies are required.

---

## Thanks to Boot.dev for the guidance!


