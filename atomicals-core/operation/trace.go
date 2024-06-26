package atomicals

import (
	"encoding/json"
	"time"

	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"

	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) TraceBlock() {
	startTime := time.Now()
	location, err := m.CurrentLocation()
	if err != nil {
		log.Log.Panicf("CurrentLocation err:%v", err)
	}
	height := location.Height
	blockInfo, err := m.GetBlockByHeight(location.Height)
	if err != nil {
		log.Log.Panicf("GetBlockByHeight err:%v height:%v", err, location.Height+1)
	}
	if location.TxIndex+1 == int64(len(blockInfo.Tx)) {
		height++
		blockInfo, err = m.GetBlockByHeight(height)
		if err != nil {
			log.Log.Panicf("GetBlockByHeight err:%v height:%v", err, location.Height+1)
		}
	}
	for index, tx := range blockInfo.Tx {
		if location.Height == height && index <= int(location.TxIndex) {
			continue
		}
		if err := m.InsertBtcTx(&postsql.BtcTx{TxID: tx.Txid, BlockHeight: blockInfo.Height}); err != nil {
			log.Log.Panicf("InsertBtcTx err:%v", err)
		}
		m.TraceTx(tx, height)
		if err := m.UpdateCurrentHeightAndExecAllSql(height, int64(index)); err != nil {
			log.Log.Panicf("UpdateCurrentHeight err:%v", err)
		}
	}
	if height%1000 == 0 {
		log.Log.Infof("height:%v, TraceBlock take time:%v", blockInfo.Height, time.Since(startTime))
	}
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) error {
	operation := witness.ParseWitness(tx, height)
	// step 1: transfer nft, transfer ft
	// m.transferNft(operation, tx)
	// if height < utils.AtOMICALS_FT_PARTIAL_SPLITING_HEIGHT {
	// 	m.transferFt(operation, tx)
	// } else {
	// 	m.transferFtPartialSpliting(operation, tx)
	// }

	// step 2: process operation
	userPk := tx.Vout[utils.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
	if operation.Op == "dmt" {
		// m.mintDistributedFt(operation, tx.Vout, userPk)
	} else {
		switch operation.Op {
		case "dft":
			// m.deployDistributedFt(operation, userPk)
		case "ft":
			// m.mintDirectFt(operation, tx.Vout, userPk)
		case "nft":
			m.mintNft(operation, userPk, height)
		case "mod":
			r, _ := json.Marshal(operation.Payload.Dmint)
			m.InsertMod(&postsql.ModInfo{
				AtomicalsID: operation.AtomicalsID,
				LocationID:  operation.LocationID,
				Mod:         string(r),
			})
		case "evt":
		case "dat":
		case "sl":
		default:
		}
	}

	// step 3 check payment
	return nil
}
