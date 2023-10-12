package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v3client"
	"go.uber.org/zap"
	"social/pkg/log"
	"sync"
)

// EtcdServer is the singleton of embedded etcd server
var (
	initOnce   sync.Once
	closeOnce  sync.Once
	etcdServer *embed.Etcd
)

// GetEmbedEtcdClient returns client of embed etcd server
func GetEmbedEtcdClient() (*clientv3.Client, error) {
	client := v3client.New(etcdServer.Server)
	return client, nil
}

// InitEtcdServer initializes embedded etcd server singleton.
func InitEtcdServer(etcdCfg *EtcdConfig) error {
	if etcdCfg.UseEmbedEtcd {
		var initError error
		initOnce.Do(func() {
			path := etcdCfg.ConfigPath
			var cfg *embed.Config
			if len(path) > 0 {
				cfgFromFile, err := embed.ConfigFromFile(path)
				if err != nil {
					initError = err
				}
				cfg = cfgFromFile
			} else {
				cfg = embed.NewConfig()
			}
			cfg.Dir = etcdCfg.DataDir
			cfg.LogOutputs = []string{etcdCfg.EtcdLogPath}
			cfg.LogLevel = etcdCfg.EtcdLogLevel
			e, err := embed.StartEtcd(cfg)
			if err != nil {
				log.Error("failed to init embedded Etcd server", zap.Error(err))
				initError = err
			}
			etcdServer = e
			log.Info("finish init Etcd config", zap.String("path", path), zap.String("data", etcdCfg.DataDir))
		})
		return initError
	}
	return nil
}

// StopEtcdServer stops embedded etcd server singleton.
func StopEtcdServer() {
	if etcdServer != nil {
		closeOnce.Do(func() {
			etcdServer.Close()
		})
	}
}
