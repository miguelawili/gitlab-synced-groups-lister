package confluence

type GetPageContent struct {
	Id      string
	Type    string
	Status  string
	Title   string
	Space   Space
	Version Version
}

type Space struct {
	Id   int
	Key  string
	Name string
}

type Version struct {
	Number int
}
