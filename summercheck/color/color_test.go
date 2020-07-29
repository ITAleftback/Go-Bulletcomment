package color

import (
	"fmt"
	"testing"
)

func Test(t *testing.T)  {
	fmt.Printf("\x1b[%dmhello world  \x1b[0m\n", 30)
	fmt.Printf("\x1b[%dmhello world  \x1b[0m\n", 31)
}
