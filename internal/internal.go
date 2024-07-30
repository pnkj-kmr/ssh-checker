package internal

// SSHChecker interface help to freeze
// the input and outout function definations
type SSHChecker interface {
	GetInput() []Input
	ProduceOutput(<-chan Output, chan<- struct{})
}

// Input represent the input
type Input struct {
	Tag      string   `json:"tag,omitempty"`
	Host     string   `json:"host"`
	Port     int      `json:"port,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
	Commands []string `json:"commands"`
}

// Output represent the output
type Output struct {
	I   Input    `json:"input"`
	Err []string `json:"error,omitempty"`
	O   []string `json:"output,omitempty"`
}

// CmdPipe helps to define the Command line Agrs
type CmdPipe struct {
	Ifile         string
	Ofile         string
	Workers       int
	Port          int
	Timeout       int
	CircleTimeout int
	Usr           string
	Passwd        string
	KeyPath       string
	HostsPath     string
}
