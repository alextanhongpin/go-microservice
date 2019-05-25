package mysqluuid

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

// New returns a new ordered uuid to be stored in the MySQL database. This
// function is useful when you need to generate the UUID from the client side,
// even though MySQL 8 and above provides the function uuid_to_bin and
// bin_to_uuid respectively.
func New() string {
	//  MySQL 5.7 uses uuid.v1.
	id := uuid.Must(uuid.NewV1())
	output := strings.Split(id.String(), "-")
	part3 := output[0]
	part2 := output[1]
	part1 := output[2]
	part4 := output[3]
	part5 := output[4]
	return strings.Join([]string{part1, part2, part3, part4, part5}, "")
}
