package logium

type Fields map[string]interface{}

const (
	AccountIDField        = "account_id"
	AccountSessionIDField = "account_session_id"

	HTTPMethodField = "http_method"
	HTTPPathField   = "http_path"

	OperationField = "operation"
	ServiceField   = "service"
	ComponentField = "component"
)
