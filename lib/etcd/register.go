package etcd

import (
	"context"
	"fmt"

	"triones-one/lib/zlog"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var ServicePrefix = "grpclb"

func RegisterService(ctx context.Context, service string, host string, port int, ttl int64, client *clientv3.Client) (func() error, error) {
	leaseResp, err := client.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	leaseID := &leaseResp.ID

	key, value := fmt.Sprintf("/%s/%s/%s:%d", ServicePrefix, service, host, port), fmt.Sprintf("%s:%d", host, port)
	_, err = client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}

	keepAlive, err := client.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return nil, err
	}

	go func() {
	keep:
		for {
			select {
			case <-client.Ctx().Done():
				return

			case k, ok := <-keepAlive:
				if !ok {
					break keep
				}
				_ = k

			case <-ctx.Done():
				return
			}
		}

	release:
		leaseResp, err := client.Grant(ctx, ttl)
		if err != nil {
			zlog.Error("Grant lease failure", zap.Error(err))
			goto release
		}

		key, value := fmt.Sprintf("/%s/%s:%d", service, host, port), fmt.Sprintf("%s:%d", host, port)
		_, err = client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
		if err != nil {
			zlog.Error("Put k/v failure", zap.Error(err), zap.String("key", key), zap.String("value", value))
			goto release
		}

		keepAlive, err = client.KeepAlive(ctx, leaseResp.ID)
		if err != nil {
			zlog.Error("Keepalive failure", zap.Error(err), zap.Any("leaseID", leaseResp.ID))
			goto release
		}
		leaseID = &leaseResp.ID

		goto keep
	}()
	revokeFunc := func() error {
		_, err = client.Revoke(ctx, *leaseID)
		return err
	}

	return revokeFunc, nil
}
