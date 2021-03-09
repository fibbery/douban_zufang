package resp

type Page struct {
	Count int    `json:"count"`
	Items []Item `json:"items"`
	Start int    `json:"start"`
	Total int    `json:"total"`
}

type Item struct {
	List Doulist `json:"doulist"`
	Time string  `json:"time"`
}

//豆列
type Doulist struct {
	ID         string `json:"id"`
	Owner      Owner  `json:"owner"`
	SharingURL string `json:"sharing_url"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	URI        string `json:"uri"`
	URL        string `json:"url"`
}

//用户
type Owner struct {
	Avatar   string `json:"avatar"`
	Followed bool   `json:"followed"`
	Gender   string `json:"gender"`
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	RegTime  string `json:"reg_time"`
	UID      string `json:"uid"`
	URI      string `json:"uri"`
	URL      string `json:"url"`
}
