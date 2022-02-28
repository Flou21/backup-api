package metrics

import (
	"net/http"

	"github.com/Coflnet/db-backup/backup-api/db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	backupsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "backup_backups_created",
		Help: "The total number of backups created",
	})
)

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func IncBackupTarget(target *db.Target) {

}
