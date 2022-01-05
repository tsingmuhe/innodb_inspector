package page

type HexEditorTag struct {
	From    uint32 `json:"from"`
	To      uint32 `json:"to"`
	Color   string `json:"color"`
	Caption string `json:"caption"`
}
