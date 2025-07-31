package reconnector

import "context"

type Connector[Conn ErrorCloserClosedChecker] interface {
	ConnectContext(ctx context.Context) (Conn, error)
}
