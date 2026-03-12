package reconcile

import (
	"fmt"

	"github.com/mamaar/features/feature"
	"github.com/mamaar/jsonchamp"
)

func Reconcile(incoming *feature.Feature, base *feature.Feature, head *feature.Feature) {

	incomingBase := incoming.Map().Diff(base.Map())
	headBase := head.Map().Diff(base.Map())

	fmt.Println(jsonchamp.InformationPaths(incomingBase))
	fmt.Println(jsonchamp.InformationPaths(headBase))
	fmt.Println("======")
}
