package dispatcher

import (
    "net"
)

type Dispatcher interface {
    Task(conn net.Conn)
}
