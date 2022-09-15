package async

import (
	"strconv"
	"testing"
)

func Test(t *testing.T) {

	Execute(func() {
		_, err := strconv.ParseUint("7770009", 10, 16)
		if err != nil {
			panic(err)
		}
	})

}
