package localize

import (
	"net"
)

/* LocalizeRequest describes a request object that
contains information about the requesters IP address
and the ID of the data they are requesting */
type LocalizeRequest interface {
	GetIPAddress() net.IP
	GetDataID() string
}
