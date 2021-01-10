package conn

type Pumper interface {
	writePump()
	readPump()
}

type Conn struct {

}