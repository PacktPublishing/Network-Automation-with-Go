package cvx

import (
	"net"
)

route: [net.IPv4 & string]: protocol: [string]: "entry-index": [string]: {
    distance: string
    metric: int
    via: [string]: type: string
}

