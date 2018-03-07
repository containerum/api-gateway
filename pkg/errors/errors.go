package errors

type Error struct {
	Message string
	ID      ID
	Details []string `json:",omitempty"`
}

type ID struct {
	Service string
	Kind    int
}

func New(msg string, details ...string) Error {
	return Error{
		Message: msg,
		ID: ID{
			Service: "gateway",
			Kind:    1,
		},
		Details: details,
	}
}
