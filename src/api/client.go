package main

import (
	"errors"

	"github.com/cymony/cryptomony/opaque"
)

type client struct {
	isInitialized bool
	cConf         *opaque.ClientConfiguration
	c             opaque.Client
}

func newClient() *client {
	return &client{isInitialized: false, cConf: nil, c: nil}
}

// RegistrationInit wasm wrapper for opaque.Client.CreateRegistrationRequest
// Takes one argument and it is string, returns []byte for registration request
// Prototype Go: RegistrationInit(password string) []byte
// Prototype JS: registrationInit(password: string) Uint8Array
func (c *client) RegistrationInit(password string) ([]byte, []byte, error) {
	if !c.IsInitialized() {
		return nil, nil, errors.New("client must be initialized first")
	}

	regState, regReq, err := c.c.CreateRegistrationRequest([]byte(password))
	if err != nil {
		return nil, nil, err
	}

	encodedRegState, err := regState.Encode()
	if err != nil {
		return nil, nil, err
	}

	encodedRegReq, err := regReq.Encode()
	if err != nil {
		return nil, nil, err
	}

	return encodedRegState, encodedRegReq, nil
}

// RegistrationFinalize wasm wrapper for opaque.Client.FinalizeRegistrationRequest
// Takes two argument, first one is string and second one is []byte
// returns []byte for registration record and []byte for exportKey
// Prototype Go: RegistrationFinalize(clientIdentity string, registrationRes []byte) ([]byte, []byte)
// Prototype JS: registrationFinalize(clientIdentity: string, registrationRes: Uint8Array) Object(Uint8Array, Uint8Array)
func (c *client) RegistrationFinalize(regState, regRes []byte, clientIdentity string) ([]byte, []byte, error) {
	if !c.IsInitialized() {
		return nil, nil, errors.New("client must be initialized first")
	}

	regisState := &opaque.ClientRegistrationState{}
	if err := regisState.Decode(c.cConf.OpaqueSuite.New(), regState); err != nil {
		return nil, nil, err
	}

	regRecord, exportKey, err := c.c.FinalizeRegistrationRequest(regisState, []byte(clientIdentity), regRes)
	if err != nil {
		return nil, nil, err
	}

	encodedRegRec, err := regRecord.Encode()
	if err != nil {
		return nil, nil, err
	}

	return encodedRegRec, exportKey, nil
}

// LoginInit wasm wrapper for opaque.Client.ClientInit
// Takes one argument and it is string, returns []byte for ke1 message
// Prototype Go: LoginInit(password string) []byte
// Prototype JS: loginInit(password: string) Uint8Array
func (c *client) LoginInit(password string) ([]byte, []byte, error) {
	if !c.IsInitialized() {
		return nil, nil, errors.New("client must be initialized first")
	}

	loginState, ke1Message, err := c.c.ClientInit([]byte(password))
	if err != nil {
		return nil, nil, err
	}

	encodedLoginState, err := loginState.Encode()
	if err != nil {
		return nil, nil, err
	}

	encodedKE1Message, err := ke1Message.Encode()
	if err != nil {
		return nil, nil, err
	}

	return encodedLoginState, encodedKE1Message, nil
}

// LoginFinish wasm wrapper for opaque.Client.ClientFinish
// Takes two argument, first one is string and second one is []byte
// returns []byte for ke3 message, []byte for sessionKey and []byte for exportKey
// Prototype Go: LoginFinish(clientIdentity string, ke2Message []byte) ([]byte, []byte, []byte)
// Prototype JS: loginFinish(clientIdentity: string, ke2Message Uint8Array) Object(Uint8Array, Uint8Array, Uint8Array)
func (c *client) LoginFinish(loginState, ke2 []byte, clientIdentity string) ([]byte, []byte, []byte, error) {
	if !c.IsInitialized() {
		return nil, nil, nil, errors.New("client must be initialized first")
	}

	logState := &opaque.ClientLoginState{}
	if err := logState.Decode(c.cConf.OpaqueSuite.New(), loginState); err != nil {
		return nil, nil, nil, err
	}

	ke3Message, sessionKey, exportKey, err := c.c.ClientFinish(logState, []byte(clientIdentity), ke2)
	if err != nil {
		return nil, nil, nil, err
	}

	encodedKE3Message, err := ke3Message.Encode()
	if err != nil {
		return nil, nil, nil, err
	}

	return encodedKE3Message, sessionKey, exportKey, nil
}

func (c *client) IsInitialized() bool {
	return c.isInitialized
}

// InitializeClient wasm wrapper for opaque.NewClient
// Takes two argument, both are string, returns nothing. But resolve promise if successful.
// Prototype Go: InitializeClient(suiteName string, serverID string)
// Prototype JS: initClient(suiteName: string, serverID: string)
func (c *client) InitializeClient(suiteName string, serverID string) error {
	cConf := &opaque.ClientConfiguration{}

	suiteID, err := strToSuite(suiteName)
	if err != nil {
		return err
	}

	cConf.OpaqueSuite = suiteID
	cConf.ServerID = []byte(serverID)

	c.c = opaque.NewClient(cConf)
	c.isInitialized = true
	c.cConf = cConf

	return nil
}
