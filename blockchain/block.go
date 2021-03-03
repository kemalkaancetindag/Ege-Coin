package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	//Block tipi olusturuyoruz (Obje gibi)
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func CreateBlock(txs []*Transaction, prevhash []byte) *Block {
	//Block pointeri olusturuyoruz ardindan block icin hash olusturuyoruz ve blocku derefer edip donuyoruz
	block := &Block{[]byte{}, txs, prevhash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

//Baslangic blockumuzu olusturuyoruz
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() []byte {
	//BadgerDB kullanacagimiz icin blocklarimizi byte slicelari haline getirmemiz gerekiyor
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b) //Blocku encodeliyoruz

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	/*
		Byte haline donmus blocku tekrar normal yapisina dondurur
	*/

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block

}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
