package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte
	/*
		Data basemiz icine blockchaini olusturuyoruz
		once databaseyi nereye olutrucagimizi seciyoruz
		ardindan opsiyonalrimiz ile badger data basemizi aciyoru
	*/

	opts := badger.DefaultOptions("")
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		/*
			Oncelikle bir blockchain databasemiz varmi diye kontrol ediyoruz
			eger yoksa genesis blockumuzu olusturuyoruz
			ardindan genesis blockumuzu "lh" keywordune deger olarak atiyoruz
			eger varas "lh" keywordunun degerini aliyoruz ve hash degerini lastHashe atiyoruz
			en sonunda blockchain referansi donuyoruz
		*/
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found!")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))

			Handle(err)

			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	/*
		Block chainimize yeni block eklemek icin BlockChain structimizin ustune bu fonku yaziyoruz
	*/
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		/*
			Block chainimizin icine read-only sekilde bakiyoruz
			Ardinden "lh" keywordumuzun degerini aliyoruz
			Onuda lastHash degiskenimizin degeri yapiyoruz
		*/
		item, err := txn.Get([]byte("lh"))

		Handle(err)

		lastHash, err = item.ValueCopy(nil)

		return err
	})

	/*
		Burda datamiz ve yeni lastHashimizi alip yeni bir block olusturuyoruz
	*/

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		/*
			Ardindan yeni blockumuzu data baseye eklerken
			"lh" keywordumuzude yeni blockumuzun hashine esitliyoruz
		*/
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.ValueCopy(nil)

		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
