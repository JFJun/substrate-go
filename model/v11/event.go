package v11

/*
解析event结构
*/
type EventResponse struct {
	EventId      string       `json:"event_id"`
	EventIdx     int          `json:"event_idx"`
	ExtrinsicIdx int          `json:"extrinsic_idx"`
	ModuleId     string       `json:"module_id"`
	Phase        int          `json:"phase"`
	Type         string       `json:"type"`
	Params       []EventParam `json:"params"`
}
type EventParam struct {
	Type     string      `json:"type"`
	Value    interface{} `json:"value"` //string or struct
	ValueRaw string      `json:"value_raw"`
}

type EventResult struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
	ExtrinsicIdx int    `json:"extrinsic_idx"`
	EventIdx     int    `json:"event_idx"`
	Status       string `json:"status"`
}
