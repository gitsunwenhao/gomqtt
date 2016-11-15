package service

import (
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/aiyun/gomqtt/proto"
	"github.com/uber-go/zap"
)

type Grpc struct {
}

func NewGrpc() *Grpc {
	g := &Grpc{}
	return g
}

func (g *Grpc) Init() {

}

func (g *Grpc) Start() {
	l, err := net.Listen("tcp", Conf.GrpcC.Addr)
	if err != nil {
		Logger.Panic("Init", zap.Error(err))
	}
	s := grpc.NewServer()

	proto.RegisterRpcServer(s, &Grpc{})
	s.Serve(l)
}

func (g *Grpc) Close() error {
	return nil
}

// BPush 广播推送
func (g *Grpc) BPush(ctx context.Context, bm *proto.BPushMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// SPush 单播推送
func (g *Grpc) SPush(ctx context.Context, sm *proto.SPushMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// PChat 私聊
func (g *Grpc) PChat(ctx context.Context, pm *proto.PChatMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// GChat 群聊
func (g *Grpc) GChat(ctx context.Context, gm *proto.GChatMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}
