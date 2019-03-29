package main

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mit-dci/go-bverify/client"
	"github.com/mit-dci/go-bverify/logging"
)

func main() {
	var err error
	hostName := flag.String("host", "localhost", "Host to connect to")
	hostPort := flag.Int("port", 9100, "Port to connect to")
	outputPath := flag.String("out", "out", "Directory to export the proofs to")
	flag.Parse()

	// This is a fixed hash that the server will commit to first before even
	// becoming available to clients. This to ensure we always get the entire
	// chain when requesting commitments.
	logging.Debugf("Starting new client and connecting to %s:%d...", *hostName, *hostPort)
	cli, err := client.NewClient([]byte{}, fmt.Sprintf("%s:%d", *hostName, *hostPort))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if cli.Ready {
				break
			}
			time.Sleep(time.Second * 1)
		}
		lastIdx := uint64(0)

		fresh := false
		logId := [32]byte{}
		if _, err := os.Stat("state.hex"); os.IsNotExist(err) {
			// Create new log
			fresh = true
		} else {
			logIdDisk, err := ioutil.ReadFile("state.hex")
			if err != nil {
				logging.Fatalf("Could not read logId from disk: %s", err.Error())
				return
			}
			copy(logId[:], logIdDisk[:32])
			lastIdx = binary.BigEndian.Uint64(logIdDisk[32:])
			logging.Debugf("LogID from disk: %x - Last index : %d", logId, lastIdx)
		}

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logging.Warn("Could not determine executable path")
		}

		for {
			time.Sleep(time.Second * 5)
			if !fresh {
				idx, _, _ := cli.GetLastCommittedLog(logId)
				if uint64(idx) != lastIdx {
					logging.Debugf("Still waiting for last statement to be committed")
					continue
				}
			}

			sense := exec.Command("/bin/bash", filepath.Join(dir, "sensor.sh"))
			stdout, err := sense.Output()

			statement := string(stdout)

			if fresh {
				// Add some randomness to prevent duplicate hashes
				r := [16]byte{}
				rand.Read(r[:])
				statement = fmt.Sprintf("[Initial log %x] %s", r, statement)
			}

			logging.Debugf("Witnessing statement [%s] to b_verify server", statement)

			if fresh {
				logId, err = cli.StartLogText(statement)
				if err != nil {
					logging.Fatalf("Could not start log: %s", err.Error())
					return
				}
				err = saveState(logId, 0)
				if err != nil {
					logging.Fatalf("Could not save state: %s", err.Error())
					return
				}
				fresh = false
			} else {

				// Right before sending a new log, export our current statement
				// (last committed one)
				fs, err := cli.ExportLog(logId)
				if err != nil {
					logging.Fatalf("Could not export last log: %s", err.Error())
					return
				}

				err = ioutil.WriteFile(filepath.Join(*outputPath, fmt.Sprintf("%d-%x.bin", time.Now().Unix(), fs.Proof.Commitment())), fs.Bytes(), 0644)
				if err != nil {
					logging.Fatalf("Could not write export to disk: %s", err.Error())
					return
				}

				lastIdx++
				err = cli.AppendLogText(lastIdx, logId, statement)
				if err != nil {
					logging.Fatalf("Could not append log: %s", err.Error())
					return
				}
				err = saveState(logId, lastIdx)
				if err != nil {
					logging.Fatalf("Could not save state: %s", err.Error())
					return
				}
			}
		}
	}()

	err = cli.Run(false)
	if err != nil {
		panic(err)
	}
}

func saveState(logId [32]byte, lastIdx uint64) error {
	state := [40]byte{}
	copy(state[:], logId[:])
	binary.BigEndian.PutUint64(state[32:], lastIdx)
	return ioutil.WriteFile("state.hex", state[:], 0600)
}
