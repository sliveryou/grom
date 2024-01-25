package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	notFormattedStr = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}
type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
  Students []Student ` + "`" + `json:"students"` + "`" + `
}
service A-api {
@server(
handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)

@server(
handler: GreetHandler2
  )
  get /greet/from2/:name(Request) returns (Response)
}
`

	formattedStr = `type Request {
	Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}
type Response {
	Message  string    ` + "`" + `json:"message"` + "`" + `
	Students []Student ` + "`" + `json:"students"` + "`" + `
}
service A-api {
	@server(
		handler: GreetHandler
	)
	get /greet/from/:name(Request) returns (Response)

	@server(
		handler: GreetHandler2
	)
	get /greet/from2/:name(Request) returns (Response)
}`
)

func TestAPIFormat(t *testing.T) {
	r, err := APIFormat(notFormattedStr)
	assert.Nil(t, err)
	assert.Equal(t, formattedStr, r)
}
