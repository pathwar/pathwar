package randstring

import (
	"fmt"
	"math/rand"
)

func ExampleRandString() {
	rand.Seed(42)
	fmt.Println(RandString(10))
	fmt.Println(RandString(42))
	// output:
	// CBzlgwF4Xt
	// dzFemEMgBqNznwB199sND0jQ6KJ402CKj1s8Oquw5O
}
