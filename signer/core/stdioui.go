
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450110486482944>

//

package core

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type StdIOUI struct {
	client rpc.Client
	mu     sync.Mutex
}

func NewStdIOUI() *StdIOUI {
	log.Info("NewStdIOUI")
client, err := rpc.DialContext(context.Background(), "stdio://“”
	if err != nil {
		log.Crit("Could not create stdio client", "err", err)
	}
	return &StdIOUI{client: *client}
}

//调度通过stdio发送请求
func (ui *StdIOUI) dispatch(serviceMethod string, args interface{}, reply interface{}) error {
	err := ui.client.Call(&reply, serviceMethod, args)
	if err != nil {
		log.Info("Error", "exc", err.Error())
	}
	return err
}

func (ui *StdIOUI) ApproveTx(request *SignTxRequest) (SignTxResponse, error) {
	var result SignTxResponse
	err := ui.dispatch("ApproveTx", request, &result)
	return result, err
}

func (ui *StdIOUI) ApproveSignData(request *SignDataRequest) (SignDataResponse, error) {
	var result SignDataResponse
	err := ui.dispatch("ApproveSignData", request, &result)
	return result, err
}

func (ui *StdIOUI) ApproveExport(request *ExportRequest) (ExportResponse, error) {
	var result ExportResponse
	err := ui.dispatch("ApproveExport", request, &result)
	return result, err
}

func (ui *StdIOUI) ApproveImport(request *ImportRequest) (ImportResponse, error) {
	var result ImportResponse
	err := ui.dispatch("ApproveImport", request, &result)
	return result, err
}

func (ui *StdIOUI) ApproveListing(request *ListRequest) (ListResponse, error) {
	var result ListResponse
	err := ui.dispatch("ApproveListing", request, &result)
	return result, err
}

func (ui *StdIOUI) ApproveNewAccount(request *NewAccountRequest) (NewAccountResponse, error) {
	var result NewAccountResponse
	err := ui.dispatch("ApproveNewAccount", request, &result)
	return result, err
}

func (ui *StdIOUI) ShowError(message string) {
	err := ui.dispatch("ShowError", &Message{message}, nil)
	if err != nil {
		log.Info("Error calling 'ShowError'", "exc", err.Error(), "msg", message)
	}
}

func (ui *StdIOUI) ShowInfo(message string) {
	err := ui.dispatch("ShowInfo", Message{message}, nil)
	if err != nil {
		log.Info("Error calling 'ShowInfo'", "exc", err.Error(), "msg", message)
	}
}
func (ui *StdIOUI) OnApprovedTx(tx ethapi.SignTransactionResult) {
	err := ui.dispatch("OnApprovedTx", tx, nil)
	if err != nil {
		log.Info("Error calling 'OnApprovedTx'", "exc", err.Error(), "tx", tx)
	}
}

func (ui *StdIOUI) OnSignerStartup(info StartupInfo) {
	err := ui.dispatch("OnSignerStartup", info, nil)
	if err != nil {
		log.Info("Error calling 'OnSignerStartup'", "exc", err.Error(), "info", info)
	}
}
func (ui *StdIOUI) OnInputRequired(info UserInputRequest) (UserInputResponse, error) {
	var result UserInputResponse
	err := ui.dispatch("OnInputRequired", info, &result)
	if err != nil {
		log.Info("Error calling 'OnInputRequired'", "exc", err.Error(), "info", info)
	}
	return result, err
}

