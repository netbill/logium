package logium

type Fields map[string]interface{}

const (
	AccountIDField        = "account_id"
	AccountSessionIDField = "account_session_id"

	HTTPMethodField = "http_method"
	HTTPPathField   = "http_path"

	UploadAccountIdField    = "upload_account_id"
	UploadSessionIdField    = "upload_session_id"
	UploadResourceTypeField = "upload_resource"
	UploadResourceIdField   = "upload_resource_id"
)
