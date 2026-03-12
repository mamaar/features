package reconcile

import (
	"fmt"

	"github.com/mamaar/jsonchamp"
)

func Reconcile(incoming *jsonchamp.Map, base *jsonchamp.Map, head *jsonchamp.Map) {

	incomingBase := incoming.Diff(base)
	headBase := head.Diff(base)

	fmt.Println(jsonchamp.InformationPaths(incomingBase))
	fmt.Println(jsonchamp.InformationPaths(headBase))
	fmt.Println("======")
}
