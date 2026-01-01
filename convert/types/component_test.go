package types

import (
	"fmt"
	"sort"
	"testing"

	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/util"
)

func TestManager_Range(t *testing.T) {
	names := make([]string, 0)
	GlobalManager().Range(func(ins named.Named) {
		names = append(names, ins.Name())
	})
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\n", name)
	}
}

func TestManager_Keys(t *testing.T) {
	fmt.Printf("%s\n", util.ToJson(GlobalManager().Keys()))
}
