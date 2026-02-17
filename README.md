# go-chip-z80

A Zilog Z80 CPU emulator written in Go.

```go
import z80 "github.com/user-none/go-chip-z80"
```

## Overview

go-chip-z80 is a cycle-counted Z80 emulator that models the full
programmer-visible state of the processor: all main and shadow registers,
index registers IX/IY, interrupt modes 0/1/2, NMI, the HALT state, and
the R refresh counter. It implements every documented opcode plus the
common undocumented ones (SLL, IX/IY half-register ops, etc.) and handles
undocumented flag behavior (F3/F5) for most instructions.

The emulator is intended to be embedded in larger system emulators. You
provide a `Bus` implementation for memory and I/O, create a CPU, and call
`Step()` in a loop.

## Usage

### Implement the Bus interface

The `Bus` interface is the only coupling point between the CPU and your
system. It has five methods covering the Z80's separate memory and I/O
address spaces:

```go
type Bus interface {
    Fetch(addr uint16) uint8        // M1 opcode fetch
    Read(addr uint16) uint8         // Memory read
    Write(addr uint16, val uint8)   // Memory write
    In(port uint16) uint8           // I/O port read
    Out(port uint16, val uint8)     // I/O port write
}
```

`Fetch` is called during M1 (opcode fetch) machine cycles. On real
hardware the M1 signal is asserted during this access, which some systems
use for wait-state insertion, memory contention, or bank switching. If
your system doesn't distinguish M1 from data reads, `Fetch` can delegate
to `Read`.

I/O methods receive the full 16-bit address bus. The low byte is the port
number from the instruction; the high byte is context-dependent (register
A for single-byte IN/OUT, register B or C for block I/O).

### Create and run the CPU

```go
cpu := z80.New(bus)

for {
    cycles := cpu.Step()
    // ... advance other system components by 'cycles' T-states ...
}
```

`Step` executes a single instruction (or services an interrupt) and
returns the number of T-states consumed.

### Cycle-budgeted execution

For frame-based emulation where you need to run the CPU for a fixed
number of cycles per frame:

```go
budget := cyclesPerFrame
for budget > 0 {
    budget -= cpu.StepCycles(budget)
}
```

`StepCycles` tracks a deficit when an instruction's cost exceeds the
remaining budget, paying it down on subsequent calls. This prevents cycle
drift across frame boundaries.

### External bus-hold cycles

When external hardware (such as a DMA controller) seizes the bus, the CPU
cannot execute but time still passes. Use `AddCycles` to advance the
cycle counter without executing an instruction:

```go
cpu.AddCycles(dmaTransferCycles)
```

### Cycle-accurate bus access

If your peripherals need to know the exact T-state of each bus access,
implement `CycleBus` in addition to `Bus`:

```go
type CycleBus interface {
    Bus
    CycleFetch(cycle uint64, addr uint16) uint8
    CycleRead(cycle uint64, addr uint16) uint8
    CycleWrite(cycle uint64, addr uint16, val uint8)
    CycleIn(cycle uint64, port uint16) uint8
    CycleOut(cycle uint64, port uint16, val uint8)
}
```

When the bus passed to `New` implements `CycleBus`, the cycle-aware
methods are used automatically. The `cycle` parameter is the CPU's
cumulative T-state counter at the time of the access.

### Interrupts

```go
// Assert the maskable interrupt line with a data bus value.
// Level-triggered: stays asserted until deasserted.
cpu.INT(true, 0xFF)

// Deassert.
cpu.INT(false, 0)

// Trigger a non-maskable interrupt (edge-triggered, latched).
cpu.NMI()
```

The CPU checks for interrupts at the start of each `Step` call. NMI has
priority over INT. Maskable interrupts are only serviced when IFF1 is set
and the one-instruction delay after EI has passed.

All three interrupt modes are supported:

| Mode | Behavior | T-states |
|------|----------|----------|
| IM 0 | Execute instruction from data bus `(RST n)` | 11 |
| IM 1 | Jump to 0x0038 | 13 |
| IM 2 | Vector table lookup at `(I << 8 \| data)` | 19 |

### Inspecting and restoring state

```go
regs := cpu.Registers()    // Snapshot of all registers
cpu.SetState(regs)         // Restore (e.g. for save states)
cpu.Cycles()               // Total T-states since last Reset
cpu.AddCycles(n)           // Advance counter without executing (DMA, etc.)
cpu.Halted()               // True if executing HALT
cpu.Reset()                // Power-on state: PC=0, SP=0xFFFF, AF=0xFFFF
```

### Save states

For save-state support (e.g. in game console emulators), the CPU provides
binary serialization of its complete internal state:

```go
buf := make([]byte, z80.SerializeSize) // pre-allocate once
cpu.Serialize(buf)                     // save
err := cpu.Deserialize(buf)            // restore
```

`SerializeSize` is a package-level constant; the buffer can be allocated once
and reused. The format is a compact little-endian encoding of all registers
and internal state. Bus references are not included — the caller handles
memory and I/O state separately.

## Design

### Instruction decoding

Instructions are decoded through five function-pointer tables, one per
prefix class:

| Table | Prefix | Contents |
|-------|--------|----------|
| `baseOps[256]` | none | Unprefixed opcodes |
| `cbOps[256]` | CB | Bit test/set/reset, rotates and shifts |
| `edOps[256]` | ED | Extended instructions (block ops, I/O, 16-bit ALU) |
| `ixOps[256]` | DD / FD | IX/IY indexed operations |
| `ixcbOps[256]` | DD CB / FD CB | Indexed bit operations |

Each table entry is a function with signature `func(c *CPU, op uint8)`.
Tables are populated in `init()` functions across the `ops_*.go` files.
A nil entry is treated as an undocumented NOP with an appropriate cycle
cost.

DD and FD prefixes share the same `ixOps` table. Before dispatching, the
prefix handler sets an internal pointer (`ixiyReg`) to IX or IY. Register
accessors for H, L, and (HL) read through this pointer, so a single set
of handlers covers both IX and IY variants. If no IX/IY-specific handler
exists, the instruction falls through to `baseOps` with a 4 T-state
prefix penalty.

### Cycle counting

Each instruction handler increments the CPU's cycle counter directly.
There is no separate timing table. The counter is a `uint64` that
persists across the CPU's lifetime and is only reset by `Reset()`.

### Instruction organization

Instruction handlers are split across files by functional category:

| File | Category |
|------|----------|
| `ops_load.go` | LD, EX, PUSH, POP |
| `ops_arith.go` | ADD, ADC, SUB, SBC, AND, XOR, OR, CP, INC, DEC, DAA, NEG |
| `ops_rotate.go` | RLCA, RRCA, RLA, RRA, CB-prefix rotates and shifts |
| `ops_bit.go` | BIT, SET, RES (CB prefix) |
| `ops_branch.go` | JP, JR, DJNZ, CALL, RET, RETI, RETN |
| `ops_block.go` | LDI, LDD, LDIR, LDDR, CPI, CPD, CPIR, CPDR |
| `ops_io.go` | IN, OUT, block I/O (INI, OUTI, etc.) |
| `ops_ed.go` | ED-prefix extended instructions |
| `ops_ix.go` | DD/FD indexed operations and DD CB/FD CB |

## Limitations

The emulator intentionally does not model two internal Z80 registers that
are invisible to normal programs but affect undocumented flag bits (F3 and
F5) in specific edge cases:

### WZ (MEMPTR) register

The Z80 has an internal 16-bit temporary register called WZ (sometimes
referred to as MEMPTR in community documentation). It is updated by many
instructions and its value leaks into the F3 and F5 flags in two cases:

- **BIT n,(HL)**: F3 and F5 come from the high byte of WZ rather than
  from the tested value. The emulator sources these flags from the tested
  value instead.
- **Block I/O repeat instructions** (INIR, INDR, OTIR, OTDR): The
  undocumented flag computation on the repeat path involves WZ. The
  emulator uses an approximation.

### q register

The Z80 has an internal flag called q that tracks whether the last
instruction modified the F register. It affects SCF and CCF:

- When q = 0 (previous instruction did not write F): F3 and F5 are set
  from `A | F`.
- When q != 0 (previous instruction wrote F): F3 and F5 are set from
  `A` alone.

The emulator always uses the q = 0 formula (`A | F`).

### Practical impact

These differences only affect undocumented flag bits (F3 and F5) in
narrow situations. No production software depends on this behavior. The
emulator passes all documented flag behavior tests and the vast majority
of undocumented flag tests. See the Testing section below for specifics.

## Testing

### Unit tests

Each `ops_*_test.go` file contains hand-written tests for its instruction
family. Run them with:

```
go test ./...
```

### SingleStepTests

The emulator is verified against hardware-captured test vectors from the
[SingleStepTests/z80](https://github.com/SingleStepTests/z80) project.
These tests were captured from a real Z80 CPU and record the exact
register, flag, memory, and cycle state for every instruction.

#### Embedded tests

A curated subset of 162 test cases covering critical instructions is
embedded directly in the `ops_*_test.go` files as `TestSST_*` functions.
These run as part of the normal test suite:

```
go test ./...
```

#### Full test runner

The full SingleStepTests suite contains 1604 JSON files with 1000 test
cases each (~1.6 million total). These are not bundled with the module.
To run them, download the test vectors and point the runner at them:

```
git clone https://github.com/SingleStepTests/z80.git
go test -run TestSSTRunner -sstpath ./z80/v1/
```

Run a single opcode file:

```
go test -v -run 'TestSSTRunner/^00\.json$' -sstpath ./z80/v1/
```

The runner skips 19 of the 1604 files by default. These correspond to the
WZ and q register limitations described above. The remaining 1585 files
(~1.585 million test cases) pass. To include the known failures:

```
go test -run TestSSTRunner -sstpath ./z80/v1/ -sststrict
```

The skip list is in `sst_runner_test.go`. As limitations are addressed,
entries can be removed to re-enable those tests.

| Skip reason | Files | Opcodes |
|---|---|---|
| SCF/CCF q-register F3/F5 | 6 | 37, 3F, DD 37, DD 3F, FD 37, FD 3F |
| BIT n,(HL) WZ/MEMPTR F3/F5 | 8 | CB 46/4E/56/5E/66/6E/76/7E |
| Block I/O repeat WZ-dependent flags | 5 | ED B1/B2/B3/BA/BB |
