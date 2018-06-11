package bt2

/*网络传输对象*/

// dht方法
const(
	pingMethod = "ping"
)

// dht方法，y属性
const (
	// 查询 状态
	statusQuery = "q"
	// 回复 状态
	statusReceive = "r"
	// 异常状态
	statusError = "e"
)

// dht方法，属性名
const (
	// 状态
	statusField = "y"
	// 消息id
	messageIdField = "t"
	// 请求方法
	requestMethodField = "q"
	// 请求主体
	requestBodyField = "a"
	// nodeId
	nodeIdField = "id"
)

/**
	构建基础参数map
 */
func buildBaseParam(messageId, status string) *map[string]interface{} {
	return &map[string]interface{}{
		statusField:    status,
		messageIdField: messageId,
	}
}

/**
	构建请求参数
 */
func buildRequestParam(messageId, status, requestMethod string) *map[string]interface{} {
	data := *buildBaseParam(messageId, status)
	data[requestMethodField] = requestMethod
	return &data
}

/**
	构建ping请求
 */
func buildPingRequest(messageId, status, requestMethod,nodeId string) *map[string]interface{}{
	data := *buildRequestParam(messageId, status, requestMethod)
	data[requestMethodField] = pingMethod
	data[requestBodyField] = map[string]interface{}{
		nodeIdField : nodeId,
	}
	return &data
}