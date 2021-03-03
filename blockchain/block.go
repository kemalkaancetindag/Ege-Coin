package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//Block tipi olusturuyoruz (Obje gibi)
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

/*
// Block icin hash uretiyoruz
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // Byte slicelarini birlestirmek icin kullabiyoruz
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}
*/

func CreateBlock(data string, prevhash []byte) *Block {
	//Block pointeri olusturuyoruz ardindan block icin hash olusturuyoruz ve blocku derefer edip donuyoruz
	block := &Block{[]byte{}, []byte(data), prevhash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

//Baslangic blockumuzu olusturuyoruz
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
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
