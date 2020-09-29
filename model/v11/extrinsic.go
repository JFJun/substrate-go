package v11

type ExtrinsicDecodeResponse struct {
	AccountId          string                 `json:"account_id"`
	CallCode           string                 `json:"call_code"`
	CallModule         string                 `json:"call_module"`
	Era                string                 `json:"era"`
	Nonce              int64                  `json:"nonce"`
	VersionInfo        string                 `json:"version_info"`
	Signature          string                 `json:"signature"`
	Params             []ExtrinsicDecodeParam `json:"params"`
	CallModuleFunction string                 `json:"call_module_function"`
}

type ExtrinsicDecodeParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"value_raw"`
}

type UtilityParamsValue struct {
	CallModule   string                  `json:"call_module"`
	CallFunction string                  `json:"call_function"`
	CallIndex    string                  `json:"call_index"`
	CallArgs     []UtilityParamsValueArg `json:"call_args"`
}

type UtilityParamsValueArg struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	ValueRaw string `json:"value_raw"`
}
