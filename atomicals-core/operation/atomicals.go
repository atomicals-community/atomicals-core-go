package atomicals

import (
	"github.com/atomicals-go/pkg/btcsync"
	"github.com/atomicals-go/pkg/conf"
	"github.com/atomicals-go/repo"
)

type Atomicals struct {
	*btcsync.BtcSync
	repo.DB
}

func NewAtomicalsWithSQL(conf *conf.Config) *Atomicals {
	b, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	return &Atomicals{
		DB:      repo.NewSqlDB(conf.SqlDNS),
		BtcSync: b,
	}
}
