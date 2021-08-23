package objects

type MQConnectionOptions struct {
	connectionUrl     *string
	clientId          *string
	autoReconnect     *bool
	keepAliveInterval *int
	connectionTimeout *int
	authKey *string

	// END USER-------------
	userId   *string
	deviceId *string
	jwtToken *string
	// --------------------

	// SYSTEM INTERNAL-----
	msToken *string
	msCliId *string
	// --------------------
}

func (M *MQConnectionOptions) AuthKey() *string {
	return M.authKey
}

func (M *MQConnectionOptions) SetAuthKey(authKey *string) {
	M.authKey = authKey
}

func (M *MQConnectionOptions) GetHeaders() *map[string]string {
	var headers map[string]string
	if M.msCliId != nil && M.msToken != nil {
		headers = make(map[string]string)
		headers["master-tkn"] = *M.msToken
		headers["master-clid"] = *M.msCliId
	} else {
		if M.jwtToken != nil && M.userId != nil && M.deviceId != nil {
			headers = make(map[string]string)
			headers["X-Token"] = *M.jwtToken
			headers["X-UID"] = *M.userId
			headers["X-DeviceID"] = *M.deviceId
		}
	}
	return &headers
}

func (M *MQConnectionOptions) GetClientID() *string {
	if M.clientId == nil {
		if M.msCliId != nil && M.msToken != nil {
			*M.clientId = *M.msToken + "|" + *M.msCliId
		} else {
			if M.jwtToken != nil && M.userId != nil && M.deviceId != nil {
				*M.clientId = *M.userId + "|" + *M.deviceId
			}
		}
	}
	return M.clientId
}

func (M *MQConnectionOptions) ConnectionUrl() *string {
	return M.connectionUrl
}

func (M *MQConnectionOptions) SetConnectionUrl(connectionUrl *string) {
	M.connectionUrl = connectionUrl
}

//func (M *MQConnectionOptions) ClientId() *string {
//	return M.clientId
//}

func (M *MQConnectionOptions) SetClientId(clientId *string) {
	M.clientId = clientId
}

func (M *MQConnectionOptions) AutoReconnect() *bool {
	return M.autoReconnect
}

func (M *MQConnectionOptions) SetAutoReconnect(autoReconnect *bool) {
	M.autoReconnect = autoReconnect
}

func (M *MQConnectionOptions) KeepAliveInterval() *int {
	return M.keepAliveInterval
}

func (M *MQConnectionOptions) SetKeepAliveInterval(keepAliveInterval *int) {
	M.keepAliveInterval = keepAliveInterval
}

func (M *MQConnectionOptions) ConnectionTimeout() *int {
	return M.connectionTimeout
}

func (M *MQConnectionOptions) SetConnectionTimeout(connectionTimeout *int) {
	M.connectionTimeout = connectionTimeout
}

func (M *MQConnectionOptions) UserId() *string {
	return M.userId
}

func (M *MQConnectionOptions) SetUserId(userId *string) {
	M.userId = userId
}

func (M *MQConnectionOptions) DeviceId() *string {
	return M.deviceId
}

func (M *MQConnectionOptions) SetDeviceId(deviceId *string) {
	M.deviceId = deviceId
}

func (M *MQConnectionOptions) JwtToken() *string {
	return M.jwtToken
}

func (M *MQConnectionOptions) SetJwtToken(jwtToken *string) {
	M.jwtToken = jwtToken
}

func (M *MQConnectionOptions) MsToken() *string {
	return M.msToken
}

func (M *MQConnectionOptions) SetMsToken(msToken *string) {
	M.msToken = msToken
}

func (M *MQConnectionOptions) MsCliId() *string {
	return M.msCliId
}

func (M *MQConnectionOptions) SetMsCliId(msCliId *string) {
	M.msCliId = msCliId
}
