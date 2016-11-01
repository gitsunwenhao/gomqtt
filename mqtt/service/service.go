package service

/* service包提供了mqtt服务器和客户端的接口服务
   Author - Sunface */
import "errors"

var (
	errInvalidConnectionType  = errors.New("service: Invalid connection type")
	errInvalidSubscriber      = errors.New("service: Invalid subscriber")
	errBufferNotReady         = errors.New("service: buffer is not ready")
	errBufferInsufficientData = errors.New("service: buffer has insufficient data")
)

// ------------------------------------------------
//           Server Api
// ------------------------------------------------
