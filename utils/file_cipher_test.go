package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"testing"
)

// TestFileCipher_Encrypt 加密功能测试
func TestFileCipher_Encrypt(t *testing.T) {
	fc := NewFileCipher("a very very very very secret key")
	if err := fc.Encrypt("../a.c"); err != nil {
		t.Error(err)
	}
}

// TestFileCipher_Decrypt 解密功能测试
func TestFileCipher_Decrypt(t *testing.T) {
	fc := NewFileCipher("a very very very very secret key")
	if err := fc.Decrypt("../a.c.sea"); err != nil {
		t.Error(err)
	}
}

func TestPBKDF2Key(t *testing.T) {
	fmt.Println([]byte("a very very very very secret key"))
	var (
		password                  = []byte("goodlick")
		salt                      = []byte("adcakfndndfg")
		iter                      = 1024
		keyLen                    = 32
		h        func() hash.Hash = func() hash.Hash {
			return sha256.New()
		}
	)
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	fmt.Println(dk[:keyLen])
	for i := 0; i < keyLen; i++ {
		fmt.Printf("%c ", dk[i])
	}
}

// BenchmarkFileCipher_Encrypt 加密性能测试
func BenchmarkFileCipher_Encrypt(b *testing.B) {
	fc := NewFileCipher("a very very very very secret key")
	fc.Encrypt("./testdata/5G.iso")
}

// BenchmarkFileCipher_Decrypt 解密性能测试
func BenchmarkFileCipher_Decrypt(b *testing.B) {
	fc := NewFileCipher("a very very very very secret key")
	fc.Decrypt("./testdata/5G.iso.sea")
}
