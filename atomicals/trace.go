package atomicals

import (
	"log"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) TraceBlock(blockInfo *btcjson.GetBlockVerboseTxResult) {
	for index, tx := range blockInfo.Tx {
		log.Printf("height:%v,txIndex:%v,txHash:%v", blockInfo.Height, index, tx.Hash)
		m.TraceTx(tx, blockInfo.Height)
	}
}
func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) {
	operation := witness.ParseOperationAndPayLoad(tx, height)
	// step 1: mint or deploy
	for _, vin := range tx.Vin {
		userPk := tx.Vout[common.VOUT_EXPECT_OUTPUT_BYTES].ScriptPubKey.Address
		switch operation.Op {
		case "dft":
			if err := m.deployFt(operation, vin, tx.Vout, userPk); err != nil {
				log.Printf("deployNft err:%+v", err)
			}
		case "dmt":
			if err := m.mintDistributedFt(operation, vin, tx.Vout, userPk); err != nil {
				log.Printf("mintDistributedFt err:%+v", err)
			}
		case "ft":
			if err := m.mintDirectFt(operation, vin, tx.Vout, userPk); err != nil {
				log.Printf("mintFt err:%+v", err)
			}
		case "nft":
			if err := m.mintNft(operation, vin, tx.Vout, userPk); err != nil {
				log.Printf("mintNft err:%+v", err)
			}
		case "mod":
		case "evt":
		case "dat":
		case "sl":
		default:
		}
	}
	// step 2: transfer nft
	if err := m.transferNft(operation, tx); err != nil {
		log.Printf("transferNft err:%+v", err)
	}
	// step 3: transfer ft
	if err := m.transferFt(operation, tx); err != nil {
		log.Printf("transferFt err:%+v", err)
	}
}