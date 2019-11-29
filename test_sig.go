package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/btcec"
)

const (
	merchantId  = "iJ7YKZdnTzUsmyJj4AaiyOo2"
	accessToken = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXJjaGFudF9pZCI6ImlKN1lLWmRuVHpVc215Smo0QWFpeU9vMiIsIndoaXRlbGFiZWwiOiJkZWZhdWx0IiwiZ3JvdXAiOiJjdXN0b2RpYWwiLCJleHAiOjE1NzUwMzE4MzR9.JU3hCxRJ4TrHQFNH4TjWKRohZyxao90TdLO9NwUZ5Awo8j1AuBA1sI4EAmEP8zB-n1HxQESsK_doFbgicbpFaA"
)

const (
	host = "https://hub-testnet.lightningpeach.com"
)

var (
	getHubPublicKeyUrl        = fmt.Sprintf("%v/%v", host, "api/v2/hub/pubkey")
	getPSSPublicKeyUrl        = fmt.Sprintf("%v/%v", host, "pss-walleto/api/v1/pss/public_key")
	getPSSDefaultPublicKeyUrl = fmt.Sprintf("%v/%v", host, "pss/api/v1/pss/public_key")
	getOnChainAddressUrl      = fmt.Sprintf("%v/%v", host, "api/v2/merchant/on_chain_address")
)

type GetHubPublicKeyResponse struct {
	Pubkey string `json: pubkey"`
}

type GetBtcAddressResponse struct {
	BtcAddress   string `json:"btc_address"`
	PssSignature string `json:"pss_signature"`
	HubSignature string `json:"hub_signature"`
}
type GenerateBitcoinAddressResp struct {
	BtcAddr      *BitcoinAddress `json:"btc_addr"`
	PssSignature string          `json:"pss_signature"`
}

type BitcoinAddress struct {
	Content string `json:"content"`
}

type PublicKey struct {
	Content string `json:"content"`
}

func SerializeGenerateBitcoinAddressResponse(resp *GenerateBitcoinAddressResp, userId string) []byte {
	payload := fmt.Sprintf("%v%v", resp.BtcAddr.Content, userId)
	return []byte(payload)
}

func getHubPublicKey() (*btcec.PublicKey, error) {
	httpResp, err := http.Get(getHubPublicKeyUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp GetHubPublicKeyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	fmt.Printf("HUB's public key: %v\n", resp.Pubkey)

	raw, err := base64.StdEncoding.DecodeString(resp.Pubkey)
	if err != nil {
		return nil, err
	}

	pk, err := btcec.ParsePubKey(raw, btcec.S256())
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func getPSSPublicKey() (*btcec.PublicKey, error) {
	httpResp, err := http.Get(getPSSPublicKeyUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp PublicKey
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	fmt.Printf("WALLETO PSS's public key: %v\n", resp.Content)

	raw, err := base64.StdEncoding.DecodeString(resp.Content)
	if err != nil {
		return nil, err
	}

	pk, err := btcec.ParsePubKey(raw, btcec.S256())
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func getPSSDefaultPublicKey() (*btcec.PublicKey, error) {
	httpResp, err := http.Get(getPSSDefaultPublicKeyUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp PublicKey
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	fmt.Printf("DEFAULT PSS's public key: %v\n", resp.Content)

	raw, err := base64.StdEncoding.DecodeString(resp.Content)
	if err != nil {
		return nil, err
	}

	pk, err := btcec.ParsePubKey(raw, btcec.S256())
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func getOnChainAddress() (*GetBtcAddressResponse, error) {
	req, err := http.NewRequest("GET", getOnChainAddressUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	var resp GetBtcAddressResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func validateHubSignature(pk *btcec.PublicKey, resp *GetBtcAddressResponse) error {
	payload := SerializeGenerateBitcoinAddressResponse(&GenerateBitcoinAddressResp{
		BtcAddr: &BitcoinAddress{
			Content: resp.BtcAddress,
		},
		PssSignature: resp.PssSignature,
	}, merchantId)
	payloadHash := sha256.Sum256(payload)

	rawSignature, err := hex.DecodeString(resp.HubSignature)
	signature, err := btcec.ParseDERSignature(rawSignature, btcec.S256())
	if err != nil {
		return err
	}

	valid := signature.Verify(payloadHash[:], pk)
	if !valid {
		return errors.New("Invalid signature")
	}
	return nil
}

func validatePSSSignature(pk *btcec.PublicKey, resp *GetBtcAddressResponse) error {
	payload := SerializeGenerateBitcoinAddressResponse(&GenerateBitcoinAddressResp{
		BtcAddr: &BitcoinAddress{
			Content: resp.BtcAddress,
		},
		PssSignature: resp.PssSignature,
	}, merchantId)
	payloadHash := sha256.Sum256(payload)

	rawSignature, err := hex.DecodeString(resp.PssSignature)
	signature, err := btcec.ParseDERSignature(rawSignature, btcec.S256())
	if err != nil {
		return err
	}

	valid := signature.Verify(payloadHash[:], pk)
	if !valid {
		return errors.New("Invalid signature")
	}
	return nil
}

func hub() error {
	pk, err := getHubPublicKey()
	if err != nil {
		return err
	}

	resp, err := getOnChainAddress()
	if err != nil {
		return err
	}

	if err := validateHubSignature(pk, resp); err != nil {
		return err
	}
	fmt.Println("HUB: DONE")
	return nil
}

func pss() error {
	pk, err := getPSSPublicKey()
	if err != nil {
		return err
	}

	resp, err := getOnChainAddress()
	if err != nil {
		return err
	}

	if err := validatePSSSignature(pk, resp); err != nil {
		return err
	}
	fmt.Println("PSS: DONE")
	return nil
}

func main() {
	if err := hub(); err != nil {
		log.Fatal(err)
	}
	if _, err := getPSSDefaultPublicKey(); err != nil {
		log.Fatal(err)
	}
	if err := pss(); err != nil {
		log.Fatal(err)
	}
}
