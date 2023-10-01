package pigeon

type SMTPRequest struct {
	to      []string
	from    string
	subject string
	body    string
}
