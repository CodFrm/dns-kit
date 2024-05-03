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

type CreateCertAfterMessage struct {
	ID      int64 `json:"id"`
	Success bool  `json:"success"`
}

func (c *CreateCertAfterMessage) Marshal() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *CreateCertAfterMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}
