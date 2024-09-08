package graph

import (
	"encoding/json"
	"fmt"
)

type Edge [2]string

func (e *Edge) String() string {
	return fmt.Sprintf("%s-to-%s", e[0], e[1])
}

func (e *Edge) From() string {
	return e[0]
}

func (e *Edge) To() string {
	return e[1]
}

func (e *Edge) Hash() string {
	return e.String() + ".json"
}

func (e *Edge) RawJson() json.RawMessage {
	j, _ := json.Marshal(e)
	return json.RawMessage(j)
}

// func (e *Edge) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(e)
// }

func (e *Edge) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// func (e *Edge) UnmarshalJSON(b []byte) error {
// 	return json.Unmarshal(b, e)
// }

func (e *Edge) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, e)
}
