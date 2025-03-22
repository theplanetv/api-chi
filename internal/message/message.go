package message

// Success message
const (
	AUTH_SUCCESS        = "Authorize success!"
	LOGIN_SUCCESS       = "Login success!"
	LOGOUT_SUCCESS      = "Logout success!"
	GET_DATA_SUCCESS    = "Get data success!"
	CREATE_DATA_SUCCESS = "Create data success!"
	UPDATE_DATA_SUCCESS = "Update data success!"
	REMOVE_DATA_SUCCESS = "Remove data success!"
)

// Failed message
const (
	INVALID_INPUT      = "Invalid input!"
	AUTH_FAILED        = "Authorize failed!"
	LOGIN_FAILED       = "Login failed!"
	GET_DATA_FAILED    = "Get data failed!"
	CREATE_DATA_FAILED = "Create data failed!"
	UPDATE_DATA_FAILED = "Update data failed!"
	REMOVE_DATA_FAILED = "Remove data failed!"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
