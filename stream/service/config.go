package service

import (
	"io/ioutil"
	"log"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"github.com/uber-go/zap"
)

type Config struct {
	CommonC *CommonConfig
	EtcdC   *EtcdConfig
	GrpcC   *GrpcConfig
}

func (c *Config) Show() {
	log.Println(c.CommonC, c.EtcdC, c.GrpcC)
}

type CommonConfig struct {
	Version  string
	IsDebug  bool
	LogLevel string
	LogPath  string
}

type EtcdConfig struct {
	Addrs      []string
	Dltimeout  int
	Rqtimeout  int
	ReportTime int64
	Reportdir  string
	TTL        int64
}

type GrpcConfig struct {
	Addr string
}

var Conf = &Config{}

func initConf() {
	Conf = &Config{
		CommonC: &CommonConfig{},
		EtcdC:   &EtcdConfig{},
		GrpcC:   &GrpcConfig{},
	}
}

func loadConfig(staticConf bool) {
	var contents []byte
	var err error

	// 初始化Conf
	initConf()
	if staticConf {
		//静态配置
		contents, err = ioutil.ReadFile("configs/stream.toml")
	} else {
		contents, err = ioutil.ReadFile("/etc/gomqtt/stream.toml")
	}

	if err != nil {
		log.Fatal("load config error", zap.Error(err))
	}

	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("parse config error", zap.Error(err))
	}
	// 解析CommonConfig
	parseCommon(tbl)

	// 初始化Logger
	InitLogger(Conf.CommonC.LogPath, Conf.CommonC.LogLevel, Conf.CommonC.IsDebug)

	// 解析EtcdConfig
	parseEtcd(tbl)

	// 解析grpc
	parseGrpc(tbl)

	Conf.Show()

}

func parseCommon(tbl *ast.Table) {
	if val, ok := tbl.Fields["common"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] parse common config: ", subTbl)
		}

		err := toml.UnmarshalTable(subTbl, Conf.CommonC)
		if err != nil {
			log.Fatalln("[FATAL] parseCommon: ", err, subTbl)
		}
	}
}

func parseEtcd(tbl *ast.Table) {
	if val, ok := tbl.Fields["etcd"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] parse etcd config: ", subTbl)
		}

		err := toml.UnmarshalTable(subTbl, Conf.EtcdC)
		if err != nil {
			log.Fatalln("[FATAL] parseEtcd: ", err, subTbl)
		}
	}
}

func parseGrpc(tbl *ast.Table) {
	if val, ok := tbl.Fields["grpc"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] parse grpc config: ", subTbl)
		}

		err := toml.UnmarshalTable(subTbl, Conf.GrpcC)
		if err != nil {
			log.Fatalln("[FATAL] parseGrpc: ", err, subTbl)
		}
	}
}
