package timestamp

import (
	"encoding/json"
	"fmt"
)

type Version[T any, K any] struct {
	data      T
	timestamp DistributedClock[K]
	id        uint64
}

func NewVersion[T any, K any](
	data T, timestamp DistributedClock[K], id uint64,
) Version[T, K] {
	return Version[T, K]{
		data:      data,
		timestamp: timestamp,
		id:        id,
	}
}

func (v Version[T, K]) String() string {
	return fmt.Sprintf("Version{data: %v, Ts: %v, Source: %v}", v.data, v.timestamp, v.id)
}

func (v Version[T, K]) GetTs() DistributedClock[K] {
	return v.timestamp
}

func (v Version[T, K]) GetData() T {
	return v.data
}

func (v Version[T, K]) GetId() uint64 {
	return v.id
}

func (v Version[T, K]) Clone() Version[T, K] {
	return Version[T, K]{
		data:      v.GetData(),
		timestamp: v.GetTs().Clone(),
		id:        v.GetId(),
	}
}

func (v Version[T, K]) MarshalJSON() ([]byte, error) {
	// Define a temporary struct for JSON serialization
	type versionJSON struct {
		Data      T               `json:"data"`
		Timestamp json.RawMessage `json:"timestamp"`
		Id        uint64          `json:"id"`
	}

	// Marshal the timestamp field separately
	tsData, err := json.Marshal(v.timestamp)
	if err != nil {
		return nil, err
	}

	// Create and populate the temporary struct
	temp := versionJSON{
		Data:      v.data,
		Timestamp: tsData,
		Id:        v.id,
	}

	// Marshal the entire structure to JSON
	return json.Marshal(temp)
}

// UnmarshalJSON implements custom JSON deserialization for Version
func (v *Version[T, K]) UnmarshalJSON(data []byte) error {
	// Define a temporary struct for JSON deserialization
	type versionJSON struct {
		Data      T               `json:"data"`
		Timestamp json.RawMessage `json:"timestamp"`
		Id        uint64          `json:"id"`
	}

	// Unmarshal JSON into the temporary struct
	var temp versionJSON
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Unmarshal the timestamp field separately

	// target := reflect.New(v.timestamp)
	// Unmarshal the data to an interface to the concrete value (which will act as a pointer, don't ask why)
	// if err := json.Unmarshal(temp.Timestamp, v.timestamp); err != nil {
	// 	return err
	// }
	// Now we get the element value of the target and convert it to the interface type (this is to get rid of a pointer type instead of a plain struct value)

	var distClockJson DistClockJson[K]

	if err := json.Unmarshal(temp.Timestamp, &distClockJson); err != nil {
		return err
	}
	converted := GetDistType[K](distClockJson.ClockType)
	converted.Set(distClockJson.Data)
	v.timestamp = converted

	// Populate the Version fields
	v.data = temp.Data
	v.id = temp.Id
	return nil
}
