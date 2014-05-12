package main

import (
	"crypto/rand"
	"crypto/rsa"
	trans "gitchain/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewReservation(t *testing.T) {
	privateKey := generateKey(t)
	txn, rand := trans.NewNameReservation("my-new-repository", &privateKey.PublicKey)
	txn1, rand1 := trans.NewNameReservation("my-new-repository", &privateKey.PublicKey)
	assert.NotEqual(t, txn.Hash, txn1.Hash, "hashes should not be equal")
	assert.NotEqual(t, rand, rand1, "random numbers should not be equal")
}

func TestReservationEncodingDecoding(t *testing.T) {
	privateKey := generateKey(t)
	txn, _ := trans.NewNameReservation("my-new-repository", &privateKey.PublicKey)

	testTransactionEncodingDecoding(t, txn)
}

func TestNewAllocation(t *testing.T) {
	privateKey := generateKey(t)
	_, rand := trans.NewNameReservation("my-new-repository", &privateKey.PublicKey)
	txn2, err := trans.NewNameAllocation("my-new-repository", rand, privateKey)

	if err != nil {
		t.Errorf("error while creating name allocation transaction: %v", err)
	}

	assert.True(t, txn2.Verify(&privateKey.PublicKey))

	txn2.Name = "my-old-repository"
	assert.False(t, txn2.Verify(&privateKey.PublicKey))

}

func TestAllocationEncodingDecoding(t *testing.T) {
	privateKey := generateKey(t)
	_, rand := trans.NewNameReservation("my-new-repository", &privateKey.PublicKey)
	txn, _ := trans.NewNameAllocation("my-new-repository", rand, privateKey)

	testTransactionEncodingDecoding(t, txn)
}

func TestNewDeallocation(t *testing.T) {
	privateKey := generateKey(t)
	txn, err := trans.NewNameDeallocation("my-new-repository", privateKey)

	if err != nil {
		t.Errorf("error while creating name allocation transaction: %v", err)
	}

	assert.True(t, txn.Verify(&privateKey.PublicKey))

	txn.Name = "my-old-repository"
	assert.False(t, txn.Verify(&privateKey.PublicKey))

}

func TestDeallocationEncodingDecoding(t *testing.T) {
	privateKey := generateKey(t)
	txn, _ := trans.NewNameDeallocation("my-new-repository", privateKey)
	testTransactionEncodingDecoding(t, txn)
}

////

func testTransactionEncodingDecoding(t *testing.T, txn trans.T) {
	e, err := txn.Encode()
	if err != nil {
		t.Errorf("error while encoding transaction: %v", err)
	}

	decoded, err1 := trans.Decode(e)
	if err1 != nil {
		t.Errorf("error while encoding transaction: %v", err1)
	}

	assert.Equal(t, decoded, txn, "encoded and decoded transaction should be identical to the original one")
}

func generateKey(t *testing.T) *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("failed to generate a key")
	}
	return privateKey
}