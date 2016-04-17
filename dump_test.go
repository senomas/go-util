package util

import (
	"fmt"
	"testing"
)

type Struct1 struct {
	Name   string
	secret string `dump:"ignore"`
	Codes  []string
	Config map[string]*Address
	Age    uint
}

type Address struct {
	Street string
	City   string
}

func (address Address) String() string {
	return Dump(address)
}

// TestDump test
func TestDump(t *testing.T) {
	t.Logf("dump string(TEST) %s\n", Dump("TEST"))
	t.Logf("dump int(12) %s\n", Dump(12))
	config := make(map[string]*Address)
	config["kampung"] = &Address{"Jalan", "Solo"}
	config["rumah"] = &Address{"Kampret1", "Jakarta"}
	config["kantor"] = &Address{"Kampret2", "Jakarta"}
	t.Logf("dump Struct1{}\n%s\n", Dump(Struct1{"seno", "secret",
		[]string{"sono", "keling"},
		config,
		17}))
	t.Logf("fmt.Sprintf %s", fmt.Sprintf("HERE\n%v", &Address{"seno", "solo"}))
}
