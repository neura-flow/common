package handlers

import (
	"strings"

	"github.com/neura-flow/common/convert/types"
)

func GetSuffixEngines() map[string][]string {
	kvm := make(map[string][]string)
	for _, key := range types.GlobalManager().Keys() {
		slice := strings.Split(key, ".")
		if val, ok := kvm[slice[0]]; ok {
			kvm[slice[0]] = append(val, slice[1])
		} else {
			kvm[slice[0]] = []string{slice[1]}
		}
	}
	return kvm
}
