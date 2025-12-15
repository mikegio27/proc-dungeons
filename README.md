# Procedurally Generated Dungeon Maps

Constraint-based procedural generation of ASCII dungeon layouts in Go.

Rooms of varying shapes are placed within a bounded grid using size, spacing, and total-coverage limits. Each room is then connected into a single corridor network using deterministic pathfinding that obeys strict routing and clearance rules.

The result is a readable dungeon with:

- Clearly separated rooms
- Corridors that (kind of) respect spatial buffers
- Walls generated automatically from topology

## Core Concepts

### Rooms

- Generated first, independently of corridors
- Shapes include:
	- Rectangle
	- Square
	- Circle (approximated)
	- Triangle (isosceles)
- Placement constraints:
	- Per-room size limits
	- Minimum spacing between rooms
	- Maximum total room area relative to grid size

### Corridors

- Each room receives exactly one door
- All rooms connect into one continuous corridor network
- First corridor starts at a random grid edge
- Subsequent corridors branch from existing corridor cells
- Corridors:
	- Never run adjacent to rooms
	- Respect a configurable buffer distance
	- Only violate spacing rules at doors (controlled, local exception)
	- Use BFS pathfinding (grid-aligned, shortest path)

### Walls

- Walls are derived, not generated
- This keeps wall logic completely out of generation rules
- Rendering concerns are isolated from generation concerns

## Legend

| Symbol | Meaning                               |
| ------ | ------------------------------------- |
| `.`    | Room floor                            |
| `#`    | Corridor floor                        |
| `▒`    | Wall (derived from adjacency)         |
| `+`    | Door (room ↔ corridor connection)     |
| `*`    | Corridor start (always on grid edge)  |
| ` `    | unused space                          |

## Seeds & Determinism

The generator uses an explicit RNG instance.

Currently:

```go
seed := time.Now().UnixNano()
```

Planned:

- Accept seed via CLI flag or config file
- Fully deterministic runs when seed is provided

## Configuration

Currently configured directly in `main.go`:

```go
gridX := int32(20)
gridY := int32(20)
maxRooms := 10
```

Additional tunables (via `generator.Config`):

- `RoomMinW`, `RoomMaxW`
- `RoomMinH`, `RoomMaxH`
- `CorridorW` (corridor thickness)
- `CorridorBuff` (minimum clearance from rooms)

Planned:

- External config file (TOML / YAML / JSON)
- CLI overrides for grid size, room count, seed

## Example Output

```text
$ go run main.go
Procedurally generating dungeon...
Using seed: 1765785529671690599
                                      ▒ . . . . . . . . . . . ▒
                                      ▒ . . . . . . . . . . . ▒
                                      ▒ . . . . . . . . . . . ▒
                                      ▒ . . . . . . . . . . . ▒
                                      ▒ . . . . . . + . . . . ▒ # # # # # # # #
                                        ▒ ▒ ▒ ▒ ▒ ▒ # ▒ ▒ ▒ ▒   #     ▒ # ▒   #
                        # # # # # # # # # # # # # # # # # # # # #   ▒ . + . ▒ #
                        #               ▒ ▒ ▒ ▒                     ▒ . . . ▒ #
                        #             ▒ . . . . ▒                   ▒ . . . ▒ #
                        #           ▒ . . . . . . ▒                 ▒ . . . ▒ #
                        #         ▒ . . . . . . . . ▒                 ▒ . ▒   #
                        #         ▒ . . . . . . . . ▒                   ▒     #
  ▒ ▒ ▒                 #         ▒ . . . . . . . . ▒         ▒ ▒ ▒ ▒ ▒ ▒ ▒   #
▒ . . . ▒               #         ▒ . . . . . . . . ▒       ▒ . . . . . . . ▒ #
▒ . . . ▒               #           ▒ . . . . . . ▒         ▒ . . . . . . . ▒ #
▒ + . . ▒               #             ▒ . . . + ▒           ▒ . . . . . . . ▒ #
  # ▒ ▒                 #               ▒ ▒ ▒ #             ▒ . . . . . . . ▒ #
  #     ▒ ▒ ▒ ▒ ▒ ▒ ▒   #               # # # # # # # # # # ▒ . . . . . . . ▒ #
  #   ▒ . . . . . . . ▒ #               #   ▒ ▒ ▒ ▒ ▒ ▒   # ▒ . . . . . . . ▒ # *
  #   ▒ . . . . . . . ▒ #               # ▒ . . . . . . ▒ # ▒ . . . . . . . ▒ #
  #   ▒ . . . . . . . ▒ #               # ▒ . . . . . . ▒ # ▒ . . . . . . . ▒ #
  #   ▒ + . . . . . . ▒ #               # ▒ . . . . . . ▒ # ▒ . . . . . . . ▒ #
  # # # # . . . . . . ▒ #               #   ▒ . . . . ▒   # ▒ . . . . . . . ▒ #
  # # ▒ . . . . . . . ▒ #               # # # + . . . ▒   # ▒ . . . . . . . ▒ #
  # # ▒ . . . . . . . ▒ #                     ▒ . . ▒     # ▒ + . . . . . . ▒ #
  # #   ▒ ▒ ▒ ▒ ▒ ▒ ▒   #                     ▒ . . ▒     # # # . . . . . . ▒ #
  # # # # # # # # # # # # # # # # # # # # # #   ▒ ▒       # ▒ . . . . . . . ▒ #
    #                                       # # # # # # # # ▒ . . . . . . . ▒ #
    # # # # # # # # # # # # # # # # # # # # # #           # ▒ . . . . . . . ▒ #
      ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒ ▒   #           #   ▒ ▒ ▒ ▒ ▒ ▒ ▒   #
    ▒ . . . . . . . . . . . . . . . . . . . ▒ # # # # # # # # # # # # # # # # #
      ▒ ▒ ▒ . . . . . . . . . . . . . ▒ ▒ ▒   #           #   #     ▒ ▒
            ▒ ▒ ▒ . + . . . . . ▒ ▒ ▒   # # # # # # # # # #   #   ▒ . . ▒
                  ▒ # ▒ . ▒ ▒ ▒   # # # #                     # ▒ . . . . ▒
                    #   ▒   # # # #                           # ▒ . . . . ▒
                    # # # # #                                 #   ▒ + . ▒
                                                              # # # # ▒




Rooms: [{{-16 -4} {-10 2} Square} {{-17 -13} {1 -10} Triangle} {{11 -8} {17 7} Rectangle} {{-2 5} {5 12} Circle} {{2 -5} {7 1} Triangle} {{15 10} {17 14} Triangle} {{0 16} {10 20} Rectangle} {{-19 5} {-17 7} Circle} {{13 -15} {16 -12} Circle}]
```