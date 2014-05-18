package db

import (
	"errors"

	"github.com/gitchain/gitchain/block"
	"github.com/gitchain/gitchain/types"
)

func (db *T) PutBlock(b *block.Block, last bool) error {
	dbtx, err := db.DB.Begin(true)
	success := false
	defer func() {
		if success {
			dbtx.Commit()
		} else {
			dbtx.Rollback()
		}
	}()
	if err != nil {
		return err
	}
	bucket, err := dbtx.CreateBucketIfNotExists([]byte("blocks"))
	if err != nil {
		return err
	}
	encoded, err := b.Encode()
	if err != nil {
		return err
	}
	err = bucket.Put(b.Hash(), encoded)
	if err != nil {
		return err
	}
	// store a reference to this block as a next block
	err = bucket.Put(append(b.PreviousBlockHash, []byte("+")...), b.Hash())
	if err != nil {
		return err
	}

	for i := range b.Transactions {
		err = bucket.Put(b.Transactions[i].Hash(), b.Hash())
		if err != nil {
			return err
		}
	}

	if last {
		bucket.Put([]byte("last"), b.Hash())
	}
	success = true
	return nil
}

func (db *T) GetBlock(hash []byte) (*block.Block, error) {
	dbtx, err := db.DB.Begin(false)
	defer dbtx.Rollback()
	if err != nil {
		return nil, err
	}
	bucket := dbtx.Bucket([]byte("blocks"))
	if bucket == nil {
		return nil, errors.New("blocks bucket does not exist")
	}
	b := bucket.Get(hash)
	if b == nil {
		return nil, errors.New("block not found")
	}
	decodedBlock, err := block.Decode(b)
	if err != nil {
		return decodedBlock, err
	}
	return decodedBlock, nil
}

func (db *T) GetLastBlock() (*block.Block, error) {
	dbtx, err := db.DB.Begin(false)
	defer dbtx.Rollback()
	if err != nil {
		return nil, err
	}
	bucket := dbtx.Bucket([]byte("blocks"))
	if bucket == nil {
		return nil, nil
	}
	last := bucket.Get([]byte("last"))
	if last == nil {
		return nil, nil
	}
	b := bucket.Get(last)
	if b == nil {
		return nil, errors.New("block not found")
	}
	decodedBlock, err := block.Decode(b)
	if err != nil {
		return decodedBlock, err
	}
	return decodedBlock, nil
}

func (db *T) GetNextBlock(hash types.Hash) (*block.Block, error) {
	dbtx, err := db.DB.Begin(false)
	defer dbtx.Rollback()
	if err != nil {
		return nil, err
	}
	bucket := dbtx.Bucket([]byte("blocks"))
	if bucket == nil {
		return nil, nil
	}
	next := bucket.Get(append(hash, []byte("+")...))
	if next == nil {
		return nil, nil
	}
	b := bucket.Get(next)
	if b == nil {
		return nil, errors.New("block not found")
	}
	decodedBlock, err := block.Decode(b)
	if err != nil {
		return decodedBlock, err
	}
	return decodedBlock, nil
}