package rcon

type ServerInfo struct {
	addr     string
	password string
	conn     *Conn
}

var serverList = make([]ServerInfo, 0)

func Get(addr string) *Conn {
	for _, s := range serverList {
		if s.addr == addr {
			return s.conn
		}
	}
	conn, _ := New(addr, "11223344")
	serverList = append(serverList, ServerInfo{addr, "11223344", conn})
	return conn
}
