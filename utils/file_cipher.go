package utils

// 流模式下的文件加/解密器

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"log"
	"os"
	"strings"
)

type FileCipher struct {
	Key   []byte       // Key 加密密钥
	Block cipher.Block // 块
	Iter  int          // 密钥轮转次数
	Ext   string       //后缀
}

// NewFileCipher 初始化
// Key password + salt 后通过N次轮转获得32位的密钥
func NewFileCipher(pkey string) *FileCipher {
	conf := GetConfig()
	fc := new(FileCipher)
	fc.Iter = conf.File.Iter
	fc.Key = PBKDF2Key([]byte(pkey), []byte(conf.File.Salt), fc.Iter, 32, func() hash.Hash {
		return sha256.New()
	})
	var err error
	if fc.Block, err = aes.NewCipher(fc.Key); err != nil {
		log.Println("error utils.NewCipher():", err)
	}
	return fc
}

// initIV 返回一个计数器模式(CTR)的、底层采用block生成key流的Stream接口，初始向量iv的长度必须等于block的块尺寸
// 主要参考来源: https://github.com/golang/go/blob/master/src/crypto/cipher/ctr.go?name=release#26
// 加密时使用
func (this *FileCipher) initIV() (stream cipher.Stream) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Println("utils.file_cipher.initIV err:", err)
		return
	}
	stream = cipher.NewCTR(this.Block, iv)
	return
}

// initWithIV 返回一个计数器模式的Stream接口
// 解密时使用
func (this *FileCipher) initWithIV(iv []byte) cipher.Stream {
	return cipher.NewCTR(this.Block, iv)
}

// Encrypt 给定加密文件路径 将该文件加密
func (this *FileCipher) Encrypt(path string) (err error) {
	if this.Block == nil {
		return errors.New("Need to Initialize Block first.")
	}
	//打开需要加密的文件
	var inFile, outFile *os.File
	if inFile, err = os.Open(path); err != nil {
		return
	}
	defer inFile.Close()
	//文件加密后的路径
	savePath := this.Extension(path, "enc")
	if outFile, err = os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777); err != nil {
		return
	}
	defer outFile.Close()
	//流设置
	stream := this.initIV()
	writer := &cipher.StreamWriter{S: stream, W: outFile}
	if _, err = io.Copy(writer, inFile); err != nil {
		return
	}
	return os.Remove(path)
}

// Decrypt 给定解密文件路径 将文件解密
func (this *FileCipher) Decrypt(path string) (err error) {
	if this.Block == nil {
		return errors.New("Need to Initialize Block first.")
	}
	//打开被加密了的文件
	var inFile, outFile *os.File
	if inFile, err = os.Open(path); err != nil {
		return
	}
	defer inFile.Close()
	// 写入到一个新文件中
	savePath := this.Extension(path, "dec")
	if outFile, err = os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777); err != nil {
		return
	}
	defer outFile.Close()
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(inFile, iv[:])
	stream := this.initIV()
	inFile.Seek(aes.BlockSize, 0) // Read after the IV
	reader := &cipher.StreamReader{S: stream, R: inFile}
	if _, err = io.Copy(outFile, reader); err != nil {
		return
	}
	return os.Remove(path)
}

// PBKDF2Key 算法来源:https://github.com/golang/crypto/blob/master/pbkdf2/pbkdf2.go
// dk := pbkdf2.Key([]byte("some password"), []byte(salt), 4096, 32, sha1.New)
func PBKDF2Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
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
	return dk[:keyLen]
}

// Extension 文件后缀处理
// 传入文件路径、处理方式[加密enc|解密dec]
func (this *FileCipher) Extension(path, method string) (filename string) {
	if method == "enc" {
		filename = path + this.Ext
	}
	if method == "dec" {
		filename = strings.TrimSuffix(path, this.Ext)
	}
	return filename
}
