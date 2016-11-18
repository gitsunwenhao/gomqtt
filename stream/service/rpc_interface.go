package service

import (
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/aiyun/gomqtt/proto"
	"github.com/uber-go/zap"
)

type Rpc struct {
	gs *grpc.Server
}

func NewRpc() *Rpc {
	rpc := &Rpc{}
	return rpc
}

func (rpc *Rpc) Init() {

}

func (rpc *Rpc) Start() {
	l, err := net.Listen("tcp", Conf.GrpcC.Addr)
	if err != nil {
		Logger.Panic("Init", zap.Error(err))
	}
	rpc.gs = grpc.NewServer()

	proto.RegisterRpcServer(rpc.gs, &Rpc{})
	go rpc.gs.Serve(l)
}

func (r *Rpc) Close() error {
	r.gs.Stop()
	return nil
}

// 推送流程

// 接收到消息查看在线
// 在线推送
// 是否要推送apns
// 存放消息

// Ack流程
// 消息Ack

// 用户相关设置流程

// 群流程

// ---------------- 消息推送相关接口  ----------------

// BPush 广播推送
func (rpc *Rpc) BPush(ctx context.Context, bm *proto.BPushMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// SPush 单播推送
func (rpc *Rpc) SPush(ctx context.Context, sm *proto.SPushMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// PChat 私聊
func (rpc *Rpc) PChat(ctx context.Context, pm *proto.PChatMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// GChat 群聊
func (rpc *Rpc) GChat(ctx context.Context, gm *proto.GChatMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// ---------------- 用户相关接口  ----------------

// LogIn 登陆
func (rpc *Rpc) LogIn(ctx context.Context, am *proto.AccMsg) (*proto.Reply, error) {
	var user *User
	user, ok := gStream.cache.As.GetUser(am.An, am.Un)
	if !ok {
		// 数据库中拉取
	}
	user.Update(am)

	Logger.Info("LogIn", zap.Object("user", user))

	return &proto.Reply{}, nil
}

// LogOut 登出
func (rpc *Rpc) LogOut(ctx context.Context, am *proto.AccMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// ---------------- 订阅相关接口  ----------------

// Subscribe 订阅
func (rpc *Rpc) Subscribe(ctx context.Context, tm *proto.TcMsg) (*proto.Reply, error) {
	return &proto.Reply{}, nil
}

// UnSubscribe 取消订阅
func (rpc *Rpc) UnSubscribe(ctx context.Context, tm *proto.TcMsg) (*proto.Reply, error) {
	return &proto.Reply{}, nil
}

// BPull 拉取广播推送
func (rpc *Rpc) BPull(ctx context.Context, bm *proto.BPushMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// SPull 拉取单播推送
func (rpc *Rpc) SPull(ctx context.Context, sm *proto.SPushMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// PPull 拉取私聊
func (rpc *Rpc) PPull(ctx context.Context, pm *proto.PChatMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}

// GPull 拉取群聊
func (rpc *Rpc) GPull(ctx context.Context, gm *proto.GChatMsg) (*proto.Reply, error) {

	return &proto.Reply{}, nil
}
