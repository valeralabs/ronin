package mempool

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"log"
)

var Logger = log.Default()

func NewTX(TXs []string) ([]string, error) {
	var rawTXs [][]byte
	var TXIDs []string

	// remove `0x` prefix from each tx string and decode to hex with hex.DecodeString
	for _, TX := range TXs {
		if len(TX) < 2 {
			return []string{}, errors.New("TX is too short")
		}

		rawTX, err := hex.DecodeString(TX[2:])

		if err != nil {
			return []string{}, errors.New("failed to decode TX")
		}

		rawTXs = append(rawTXs, rawTX)
		sum := sha512.Sum512_256(rawTX)
		TXIDs = append(TXIDs, hex.EncodeToString(sum[:]))
	}

	// log hex-encoded txs using hex.EncodeToString(tx)
	// here we'd probably do something with the txs
	for _, TX := range rawTXs {
		sum := sha512.Sum512_256(TX)
		TXID := hex.EncodeToString(sum[:])
		Logger.Println("new mempool tx:", hex.EncodeToString([]byte(TXID)))
	}

	return TXIDs, nil
}
