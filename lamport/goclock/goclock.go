package goclock

import (
	"lamport/timestamp"
)

func Main() {
	// config for data type

	f := MachineFactory[int, *timestamp.LamportClock]{
		&timestamp.LamportClock{},
	}.UseDataType(StringType)

}
