package confluence

type UpdatePagePayload struct {
	Version VersionPayload `json:"version"`
	Type    string         `json:"type"`
	Title   string         `json:"title"`
	Body    BodyPayload    `json:"body"`
}

type VersionPayload struct {
	Number    int  `json:"number"`
	MinorEdit bool `json:"minorEdit"`
}

type BodyPayload struct {
	Storage ViewPayload `json:"storage"`
}

type ViewPayload struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}
