package server

const (
	CLASS_UTILS      = 0x0B
	CLASS_AUTH       = 0x74
	CLASS_SERVERS    = 0x03
	CLASS_SERVER_ACK = 0x04
	CLASS_LOGINCLASS = 66
)

const (
	TYPE_PING           = 0x00
	TYPE_AUTHREQ        = 0x00
	TYPE_AUTHACK        = 0xD6
	TYPE_SERVREQ        = 0xD6
	TYPE_SERVACK        = 0x00
	TYPE_AUTHCLASSICREQ = 0
)
