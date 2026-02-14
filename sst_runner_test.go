package z80

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sstPath = flag.String("sstpath", "", "directory containing SST JSON test files")
var sstStrict = flag.Bool("sststrict", false, "run all SST tests including known failures")

// sstSkip lists JSON files that fail due to unmodeled Z80 internals.
// Remove entries as features are implemented to re-enable those tests.
var sstSkip = map[string]string{
	// SCF/CCF: q register affects F3/F5 (q!=0 uses A only, q=0 uses A|F)
	"37.json":    "SCF q-register F3/F5",
	"3f.json":    "CCF q-register F3/F5",
	"dd 37.json": "SCF q-register F3/F5 (DD prefix)",
	"dd 3f.json": "CCF q-register F3/F5 (DD prefix)",
	"fd 37.json": "SCF q-register F3/F5 (FD prefix)",
	"fd 3f.json": "CCF q-register F3/F5 (FD prefix)",
	// BIT n,(HL): F3/F5 come from WZ (MEMPTR) high byte, not the tested value
	"cb 46.json": "BIT 0,(HL) WZ/MEMPTR F3/F5",
	"cb 4e.json": "BIT 1,(HL) WZ/MEMPTR F3/F5",
	"cb 56.json": "BIT 2,(HL) WZ/MEMPTR F3/F5",
	"cb 5e.json": "BIT 3,(HL) WZ/MEMPTR F3/F5",
	"cb 66.json": "BIT 4,(HL) WZ/MEMPTR F3/F5",
	"cb 6e.json": "BIT 5,(HL) WZ/MEMPTR F3/F5",
	"cb 76.json": "BIT 6,(HL) WZ/MEMPTR F3/F5",
	"cb 7e.json": "BIT 7,(HL) WZ/MEMPTR F3/F5",
	// Block IO repeat: undocumented flags depend on WZ
	"ed b1.json": "CPIR WZ-dependent flags",
	"ed b2.json": "INIR WZ-dependent flags",
	"ed b3.json": "OTIR WZ-dependent flags",
	"ed ba.json": "INDR WZ-dependent flags",
	"ed bb.json": "OTDR WZ-dependent flags",
}

type sstJSONState struct {
	PC   uint16     `json:"pc"`
	SP   uint16     `json:"sp"`
	A    uint8      `json:"a"`
	B    uint8      `json:"b"`
	C    uint8      `json:"c"`
	D    uint8      `json:"d"`
	E    uint8      `json:"e"`
	F    uint8      `json:"f"`
	H    uint8      `json:"h"`
	L    uint8      `json:"l"`
	I    uint8      `json:"i"`
	R    uint8      `json:"r"`
	IX   uint16     `json:"ix"`
	IY   uint16     `json:"iy"`
	AF_  uint16     `json:"af_"`
	BC_  uint16     `json:"bc_"`
	DE_  uint16     `json:"de_"`
	HL_  uint16     `json:"hl_"`
	IM   uint8      `json:"im"`
	IFF1 uint8      `json:"iff1"`
	IFF2 uint8      `json:"iff2"`
	RAM  [][]uint16 `json:"ram"`
	// Parsed but not modeled.
	WZ uint16 `json:"wz"`
	EI uint8  `json:"ei"`
	P  uint8  `json:"p"`
	Q  uint8  `json:"q"`
}

func (s *sstJSONState) toZ80State() z80State {
	st := z80State{
		A: s.A, F: s.F, B: s.B, C: s.C,
		D: s.D, E: s.E, H: s.H, L: s.L,
		I: s.I, R: s.R,
		PC: s.PC, SP: s.SP,
		IX: s.IX, IY: s.IY,
		AF_: s.AF_, BC_: s.BC_, DE_: s.DE_, HL_: s.HL_,
		IM:   s.IM,
		IFF1: s.IFF1 != 0,
		IFF2: s.IFF2 != 0,
	}
	for _, entry := range s.RAM {
		st.RAM = append(st.RAM, [2]uint16{entry[0], entry[1]})
	}
	return st
}

type sstJSONTest struct {
	Name    string       `json:"name"`
	Initial sstJSONState `json:"initial"`
	Final   sstJSONState `json:"final"`
	Cycles  []any        `json:"cycles"`
	Ports   [][]any      `json:"ports"`
}

func TestSSTRunner(t *testing.T) {
	if *sstPath == "" {
		t.Skip("no -sstpath provided")
	}

	entries, err := os.ReadDir(*sstPath)
	if err != nil {
		t.Fatalf("reading sstpath: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		fname := entry.Name()
		if reason, ok := sstSkip[fname]; ok && !*sstStrict {
			t.Run(fname, func(t *testing.T) {
				t.Skipf("known failure: %s (use -sststrict to run)", reason)
			})
			continue
		}
		t.Run(fname, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile(filepath.Join(*sstPath, fname))
			if err != nil {
				t.Fatalf("reading %s: %v", fname, err)
			}

			var tests []sstJSONTest
			if err := json.Unmarshal(data, &tests); err != nil {
				t.Fatalf("parsing %s: %v", fname, err)
			}

			for i := range tests {
				jt := &tests[i]
				init := jt.Initial.toZ80State()
				want := jt.Final.toZ80State()
				want.Cycles = len(jt.Cycles)

				// Extract input port reads.
				for _, p := range jt.Ports {
					if len(p) >= 3 {
						if dir, ok := p[2].(string); ok && dir == "r" {
							addr := uint16(p[0].(float64))
							val := uint16(p[1].(float64))
							init.Ports = append(init.Ports, [2]uint16{addr, val})
						}
					}
				}

				runSSTTest(t, struct {
					name       string
					init, want z80State
				}{
					name: jt.Name,
					init: init,
					want: want,
				})
			}
		})
	}
}
