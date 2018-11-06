package storage

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/orvice/monitor-client/mod"
)

type InfluxStorage struct {
	client    client.Client
	dbName    string
	tableName string
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
		}

		pt, err := client.NewPoint(
			i.tableName,
			tags,
			fields,
			time.Now(),
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
