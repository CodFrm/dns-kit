package message

import "encoding/json"

type CreateCertMessage struct {
	ID int64 `json:"id"`
}

func (c *CreateCertMessage) Marshal() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *CreateCertMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}
