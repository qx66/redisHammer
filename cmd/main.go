package main

import (
	"flag"
	"fmt"
	"github.com/qx66/redisHammer/pkg"
	"go.uber.org/zap"
)

//
var addr string
var password string
var db int

//
var scan bool
var scanMatch string
var keyType bool

//
var keyScan bool
var keyScanMatch string

//
var get bool
var key string

//
var del bool

//
var info bool

//
var clientList bool

func init() {
	// 连接参数
	flag.StringVar(&addr, "addr", "127.0.0.1:6379", "-addr")
	flag.StringVar(&password, "password", "", "-password")
	flag.IntVar(&db, "db", 0, "-db")
	
	// Scan
	flag.BoolVar(&scan, "scan", false, "-scan")
	flag.StringVar(&scanMatch, "scanMatch", "", "-scanMatch")
	// keyType 为 true 时，获取 key 类型，不过会让速度慢几百/几千倍以上
	flag.BoolVar(&keyType, "keyType", false, "-keyType")
	
	// scanKey
	flag.BoolVar(&keyScan, "keyScan", false, "-keyScan")
	flag.StringVar(&keyScanMatch, "keyScanMatch", "", "-keyScanMatch")
	
	// Get
	flag.BoolVar(&get, "get", false, "-get")
	
	// Get
	flag.BoolVar(&del, "del", false, "-del")
	
	// Info
	flag.BoolVar(&info, "info", false, "-info")
	
	// ClientList
	flag.BoolVar(&clientList, "clientList", false, "-clientList")
	
	// common
	flag.StringVar(&key, "key", "", "-key")
}

func main() {
	flag.Parse()
	
	//
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	//
	cli, err := pkg.NewClient(addr, password, db, logger)
	if err != nil {
		logger.Error(
			"redis连接失败",
			zap.Error(err),
		)
		return
	}
	
	switch {
	case scan:
		cli.Scan(scanMatch, 1000, keyType)
	case get:
		cli.Get(key)
	case del:
		cli.Del(key)
	case keyScan:
		cli.KeyScan(key, keyScanMatch, 1000)
	case info:
		cli.Info()
	case clientList:
		cli.ClientList()
	
	default:
		fmt.Println("Usage: --help")
	}
	
}
