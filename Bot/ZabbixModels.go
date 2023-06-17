package main

type UserLoginRequest struct {
	Jsonrpc string	`json:"jsonrpc"`
	Method string	`json:"method"`
	Id int			`json:"id"`
	Params struct{
		Username string	`json:"username"`
		Password string	`json:"password"`
	}				`json:"params"`
}

func newUserLoginRequest(id int, username string, password string) UserLoginRequest{
	ans := UserLoginRequest{}
	ans.Jsonrpc = "2.0"
	ans.Method = "user.login"
	ans.Id = id
	ans.Params.Username = username
	ans.Params.Password = password
	return ans
}

type UserLoginResponse struct {
	Jsonrpc string		`json:"jsonrpc"`
	Result string		`json:"result"`
	Id int				`json:"id"`
	Error ErrorResponse	`json:"error"`
}

type EventGetRequest struct {
	Jsonrpc string	`json:"jsonrpc"`
	Method string	`json:"method"`
	Id int			`json:"id"`
	Params struct{
		Eventids	[]string	`json:"eventids"`
	}				`json:"params"`
}

func newEventGetRequest() EventGetRequest {
	ans := EventGetRequest{}
	ans.Jsonrpc = "2.0"
	ans.Method = "event.get"
	return ans
}

type EventGetResponse struct {
	Jsonrpc string				`json:"jsonrpc"`
	Result []struct{
		Eventid string			`json:"eventid"`
		Name string 			`json:"name"`
		Severity string			`json:"severity"`
		Clock string			`json:"clock"`
		Acknowledges []struct{
			Clock string 		`json:"clock"`
			Username string		`json:"username"`
			Message string		`json:"message"`
		}						`json:"acknowledges"`
	}							`json:"result"`
	Id int						`json:"id"`
}

type ProblemGetRequest struct {
	Jsonrpc string	`json:"jsonrpc"`
	Method string	`json:"method"`
	Id int			`json:"id"`
	Params struct{
	}			`json:"params"`
}

func newProblemGetRequest() ProblemGetRequest {
	ans := ProblemGetRequest{}
	ans.Jsonrpc = "2.0"
	ans.Method = "problem.get"
	return ans
}

type EventAcknowledgeRequest struct {
	Jsonrpc string	`json:"jsonrpc"`
	Method string	`json:"method"`
	Id int			`json:"id"`
	Params struct{
		Eventids	string	`json:"eventids"`
		Action      int		`json:"action"`
		Message		string	`json:"message"`
		Severity	int		`json:"severity"`
	}				`json:"params"`
}

func newEventAcknowledgeRequest() EventAcknowledgeRequest{
	ans := EventAcknowledgeRequest{}
	ans.Jsonrpc = "2.0"
	ans.Method = "event.acknowledge"
	return ans
}

type ScriptGetscriptsbyeventsRequest struct {
	Jsonrpc string	`json:"jsonrpc"`
	Method string	`json:"method"`
	Id int			`json:"id"`
	Params struct{
		Eventids	string	`json:"eventids"`
	}				`json:"params"`
}

func newScriptGetscriptsbyeventsRequest() ScriptGetscriptsbyeventsRequest{
	ans := ScriptGetscriptsbyeventsRequest{}
	ans.Jsonrpc = "2.0"
	ans.Method = "script.getscriptsbyevents"
	return ans
}

type ScriptGetscriptsbyeventsResponse struct {
	Jsonrpc string				`json:"jsonrpc"`
	Result map[string][]struct {
		Scriptid string 		`json:"scriptid"`
		Name     string 		`json:"name"`
	}							`json:"result"`
	Id int						`json:"id"`
}

type ScriptExecuteRequest struct {
	Jsonrpc string	`json:"jsonrpc"`
	Method string	`json:"method"`
	Id int			`json:"id"`
	Params struct{
		Scriptid string	`json:"scriptid"`
		Eventid	string		`json:"eventid"`
	}				`json:"params"`
}

func newScriptExecuteRequest() ScriptExecuteRequest{
	ans := ScriptExecuteRequest{}
	ans.Jsonrpc = "2.0"
	ans.Method = "script.execute"
	return ans
}

type ScriptExecuteResponse struct {
	Jsonrpc string				`json:"jsonrpc"`
	Result struct {
		Response string			`json:"response"`
		Value string			`json:"value"`
	}							`json:"result"`
	Id int						`json:"id"`
	Error ErrorResponse			`json:"error"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (r ProblemGetRequest) setZabbixId(id int){
	r.Id = id
}
func (r EventGetRequest) setZabbixId(id int){
	r.Id = id
}
func (r EventAcknowledgeRequest) setZabbixId(id int){
	r.Id = id
}
func (r ScriptExecuteRequest) setZabbixId(id int){
	r.Id = id
}
func (r ScriptGetscriptsbyeventsRequest) setZabbixId(id int){
	r.Id = id
}


