package base

type BaseSong struct {
	ID       string `json:"id"`
	Mid      string `json:"mid"`
	Name     string `json:"name"`
	Singer   string `json:"singer"`
	Interval int    `json:"interval"`
}
