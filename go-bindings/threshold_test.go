package blschia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreshold(t *testing.T) {
	n, m := 5, 3
	var hashes []Hash
	for i := 0; i < n; i++ {
		hashes = append(hashes, genHash(byte(i+1)))
	}
	hash := genHash(100)
	pks := make([]*G1Element, m)
	sigs := make([]*G2Element, m)
	sks := make([]*PrivateKey, m)
	for i := 0; i < m; i++ {
		sks[i] = mustPrivateKeyFromBytes(genSeed(byte(i+11)), true)
		pks[i] = mustGetG1(sks[i])
		sigs[i] = ThresholdSign(sks[i], hash)
		assert.True(t, ThresholdVerify(pks[i], hash, sigs[i]))
	}
	sig := ThresholdSign(sks[0], hash)
	assert.True(t, ThresholdVerify(pks[0], hash, sig))
	skShares := make([]*PrivateKey, n)
	pkShares := make([]*G1Element, n)
	sigShares := make([]*G2Element, n)
	for i := 0; i < n; i++ {
		skShares[i], _ = ThresholdPrivateKeyShare(sks, hashes[i])
		pkShares[i], _ = ThresholdPublicKeyShare(pks, hashes[i])
		sigShares[i], _ = ThresholdSignatureShare(sigs, hashes[i])
		assert.True(t, mustGetG1(skShares[i]).EqualTo(pkShares[i]))
		sigShare2 := ThresholdSign(skShares[i], hash)
		assert.True(t, sigShares[i].EqualTo(sigShare2))
		assert.True(t, ThresholdVerify(pkShares[i], hash, sigShares[i]))
	}

	recSk, _ := ThresholdPrivateKeyRecover(skShares[:m-1], hashes[:m-1])
	recPk, _ := ThresholdPublicKeyRecover(pkShares[:m-1], hashes[:m-1])
	recSig, _ := ThresholdSignatureRecover(sigShares[:m-1], hashes[:m-1])
	assert.False(t, recSk.EqualTo(sks[0]))
	assert.False(t, recPk.EqualTo(pks[0]))
	assert.False(t, recSig.EqualTo(sig))

	recSk, _ = ThresholdPrivateKeyRecover(skShares[:m], hashes[:m])
	recPk, _ = ThresholdPublicKeyRecover(pkShares[:m], hashes[:m])
	recSig, _ = ThresholdSignatureRecover(sigShares[:m], hashes[:m])
	assert.True(t, recSk.EqualTo(sks[0]))
	assert.True(t, recPk.EqualTo(pks[0]))
	assert.True(t, recSig.EqualTo(sig))
}

func mustGetG1(sk *PrivateKey) *G1Element {
	pk, err := sk.G1Element()
	if err != nil {
		panic(err)
	}
	return pk
}

func mustPrivateKeyFromBytes(data []byte, modOrder bool) *PrivateKey {
	sk, err := PrivateKeyFromBytes(data, modOrder)
	if err != nil {
		panic(err)
	}
	return sk
}

func genHash(v byte) [HashSize]byte {
	var hash Hash
	copy(hash[:], genSeed(v))
	return hash
}
