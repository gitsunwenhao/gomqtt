package gate

import (
	"context"
	"log"

	"google.golang.org/grpc"

	rpc "github.com/aiyun/gomqtt/proto"
	"github.com/uber-go/zap"
)

type Rpc struct {
	conn   *grpc.ClientConn
	client rpc.RpcClient
}

func (r *Rpc) Init(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := rpc.NewRpcClient(conn)

	r.conn = conn
	r.client = c
}

func (r *Rpc) Close(addr string) {
	r.conn.Close()
}

// 用户登录接口
func (r *Rpc) LogIn(acm *rpc.AccMsg) error {
	req, err := r.client.LogIn(context.Background(), acm)
	if err != nil {
		Logger.Error("LogIn", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

func (r *Rpc) LogOut(acm *rpc.AccMsg) error {
	req, err := r.client.LogOut(context.Background(), acm)
	if err != nil {
		Logger.Error("LogOut", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

// 用户订阅相关
func (r *Rpc) Subscribe(tm *rpc.TcMsg) error {
	req, err := r.client.Subscribe(context.Background(), tm)
	if err != nil {
		Logger.Error("Subscribe", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

func (r *Rpc) UnSubscribe(tm *rpc.TcMsg) error {
	req, err := r.client.Subscribe(context.Background(), tm)
	if err != nil {
		Logger.Error("UnSubscribe", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

// 推送接口
func (r *Rpc) BPush(ctx context.Context, bm *rpc.BPushMsg) error {
	req, err := r.client.BPush(context.Background(), bm)
	if err != nil {
		Logger.Error("BPush", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

func (r *Rpc) SPush(ctx context.Context, sp *rpc.SPushMsg) error {
	req, err := r.client.SPush(context.Background(), sp)
	if err != nil {
		Logger.Error("SPush", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

func (r *Rpc) PChat(ctx context.Context, pm *rpc.PChatMsg) error {
	req, err := r.client.PChat(context.Background(), pm)
	if err != nil {
		Logger.Error("PChat", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}

func (r *Rpc) GChat(ctx context.Context, gm *rpc.GChatMsg) error {
	req, err := r.client.GChat(context.Background(), gm)
	if err != nil {
		Logger.Error("GChat", zap.Error(err))
		return err
	}
	log.Println(req)
	return nil
}
