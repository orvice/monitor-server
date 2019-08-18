package storage

import (
	"fmt"
	_ "github.com/influxdata/influxdb-client-go"
	"time"
	"github.com/orvice/monitor-client/mod"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxStorage struct {
	client    client.Client
	dbName    string
	tableName string
}

func NewInfluxStorage(cli client.Client, dbName, tableName string) *InfluxStorage {
	return &InfluxStorage{
		cli,
		dbName,
		tableName,
	}
}

func (i *InfluxStorage) InsertNodeInfo(ns []mod.NodeStat) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: i.dbName,
	})
	if err != nil {
		return err
	}

	for _, n := range ns {
		tags := map[string]string{

			"id":   fmt.Sprintf("%d", n.NodeID),
			"name": n.NodeName,
		}

		fields := map[string]interface{}{
			"bytes_recv":  int64(n.Stat.NetSpeed.BytesRecv),
			"bytes_sent":  int64(n.Stat.NetSpeed.BytesSent),
			"bytes_total": int64(n.Stat.NetSpeed.BytesRecv + n.Stat.NetSpeed.BytesSent),
			"memory_used": n.Stat.MemoryStatus.UsedPercent,
			"load1":       n.Stat.AvgLoad.Load1,
			"load5":       n.Stat.AvgLoad.Load5,
			"load15":      n.Stat.AvgLoad.Load15,
			"net_status":  n.Stat.NetInfo.Status,
			"disk_usage":  n.Stat.DiskUsage.UsedPercent,
		}

		pt, err := client.NewPoint(
			i.tableName,
			tags,
			fields,
			time.Unix(n.Time, 0),
		)
		if err != nil {
			return err
		}

		bp.AddPoint(pt)
	}
	err = i.client.Write(bp)
	if err != nil {
		return err
	}
	return nil
}
