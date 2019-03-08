package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"

	u "gx/ipfs/QmNohiVssaPw3KVLZik59DBVGTSm2dGvYT9eoXt5DQ36Yz/go-ipfs-util"

	crypto "github.com/libp2p/go-libp2p-crypto"
	kb "github.com/libp2p/go-libp2p-kbucket"
	peer "github.com/libp2p/go-libp2p-peer"
)

var (
	alphabet      = regexp.MustCompile("^[123456789abcdefghijklmnopqrstuvwxyz]+$")
	numWorkers    = runtime.NumCPU()
	byteN         = 2
	precision     = 256
	examplePeerID = "QmaSCVHThE4syxb8hDnjMgCPvjsN9gedNBD2u2UeSs1hJk"
)

// Key stores PrettyID containing desired substring at Index
type Key struct {
	GenPrettyID  string
	DestPrettyID string
	DHTGen       []byte
	DHTDest      []byte
	Xor          []byte
	Dist         int
}

func main() {
	if len(os.Args) > 2 {
		fmt.Printf(`
This tool generates IPFS public and private keypair until it finds public key
which contains required substring. Keys are stored in local directory. If you
like one of them, you can move it to ~/.ipfs/keystore/ to use it with IPFS.

Future Goal: Given a peer ID, this tool generates IPFS keypairs close to dest
ination ID.

Usage:
	%s {part} {destination ID}
		For fast results suggested length of public key part is 4-5 characters
`, os.Args[0])
		os.Exit(1)
	}

	destinationID := examplePeerID

	if len(os.Args) > 1 {
		destinationID = os.Args[1]
	}

	runtime.GOMAXPROCS(numWorkers)
	keyChan := make(chan Key)
	for i := 0; i < numWorkers; i++ {
		go func() {
			err := generateKey(destinationID, keyChan)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	for key := range keyChan {
		fmt.Printf("%s\n", key.DestPrettyID)
		fmt.Printf("%s\n", key.GenPrettyID)
		fmt.Print("\n")
		printByte(key.DHTDest, byteN)
		printByte(key.DHTGen, byteN)
		printByte(key.Xor, byteN)
		fmt.Printf("%d", key.Dist)
		fmt.Print("\n")
	}
}

func printByte(byteSlice []byte, bytes int) {
	for i := 0; i < bytes; i++ {
		fmt.Printf("%08b", byteSlice[i])
	}
	fmt.Print("\n")
}

func power(a, n int) int {
	var i, result int
	result = 1
	for i = 0; i < n; i++ {
		result *= a
	}
	return result
}

func byteArrayToInt(byteSlice []byte, bytes int) int {
	sum := 0
	for i := 0; i < bytes; i++ {
		sum = sum + power(2, ((bytes-i-1)*8))*int(byteSlice[i])
	}

	return sum
}

func generateKey(destPrettyID string, keyChan chan Key) error {
	for {
		//privateKey, publicKey, err := crypto.GenerateEd25519Key(rand.Reader)
		privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048, rand.Reader)
		if err != nil {
			return err
		}

		genID, err := peer.IDFromPublicKey(publicKey)
		if err != nil {
			return err
		}

		dhtGen := kb.ConvertPeerID(genID)

		destID, err := peer.IDB58Decode(destPrettyID)
		if err != nil {
			return err
		}

		dhtDest := kb.ConvertPeerID(destID)

		xor := u.XOR(dhtGen, dhtDest)

		dist := byteArrayToInt(xor, byteN)

		if dist > precision {
			continue
		}

		genPrettyID := genID.Pretty()

		privateKeyBytes, err := privateKey.Bytes()
		if err != nil {
			return err
		}

		err = ioutil.WriteFile( /*prettyID*/ "NUL", privateKeyBytes, 0600)
		if err != nil {
			return err
		}

		keyChan <- Key{
			GenPrettyID:  genPrettyID,
			DestPrettyID: destPrettyID,
			DHTGen:       dhtGen,
			DHTDest:      dhtDest,
			Xor:          xor,
			Dist:         dist,
		}
	}
}
