package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	tcpGetOk  = "TCP-GET-OK"
	tcpAddOk  = "TCP-ADD-OK"
	tcpUnknow = "TCP-UNKNOW-REQST"
)

// If have error return true
func ErrWrap(err error) bool {
	if err != nil {
		fmt.Errorf(err.Error())
		return true
	}

	return false
}

type TcpRequestJson struct {
	RqstType string `json:"rqst_type"`
	FileName string `json:"file_name"`
	FileData []byte `json:"file_data"`
	Status   string `json:"status"`
}

// Server
func main() {
	fmt.Println("Tcp File transfer Start ... ")

	// Listen port 8080
	tcpLst, err := net.Listen("tcp", ":8080")
	if ErrWrap(err) {
		return
	}

	// Accept clients
	for {
		client, err := tcpLst.Accept()
		if ErrWrap(err) {
			return
		}

		request := TcpRequestJson{}

		rqstDecoder := json.NewDecoder(client)
		if ErrWrap(rqstDecoder.Decode(&request)) {
			return
		}
		request.Status = tcpUnknow

		switch request.RqstType {
		case "GET":
			{
				getFile, err := os.Open(request.FileName)
				if ErrWrap(err) {
					request.Status = err.Error()
					break
				}

				fileData, err := io.ReadAll(getFile)
				if ErrWrap(err) {
					request.Status = err.Error()
					break
				}

				request.Status = tcpGetOk
				request.FileData = fileData
				break
			}
		case "ADD":
			{
				_, err := os.Create(request.FileName)
				if ErrWrap(err) {
					request.Status = err.Error()
					break
				}

				request.Status = tcpAddOk
				break
			}

			break
		}

		// Client send file
		json.NewEncoder(client).Encode(request)

		client.Close()
	}
}
