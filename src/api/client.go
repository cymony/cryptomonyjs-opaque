package main

import (
	"errors"
	"syscall/js"

	"github.com/cymony/cryptomony/opaque"
)

type client struct {
	isInitialized     bool
	registrationState *opaque.ClientRegistrationState
	loginState        *opaque.ClientLoginState
	c                 opaque.Client
}

func newClient() *client {
	return &client{isInitialized: false, registrationState: nil, loginState: nil, c: nil}
}

func (c *client) exposeClient(rootModule js.Value) {
	rootModule.Set("client", make(map[string]interface{}))
	clientModule := rootModule.Get("client")

	clientModule.Set("initClient", js.FuncOf(c.InitializeClient))
	clientModule.Set("isInitialized", js.FuncOf(c.IsInitialized))
	clientModule.Set("registrationInit", js.FuncOf(c.RegistrationInit))
	clientModule.Set("registrationFinalize", js.FuncOf(c.RegistrationFinalize))
	clientModule.Set("loginInit", js.FuncOf(c.LoginInit))
	clientModule.Set("loginFinish", js.FuncOf(c.LoginFinish))
}

// RegistrationInit wasm wrapper for opaque.Client.CreateRegistrationRequest
// Takes one argument and it is string, returns []byte for registration request
// Prototype Go: RegistrationInit(password string) []byte
// Prototype JS: registrationInit(password: string) Uint8Array
func (c *client) RegistrationInit(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !c.isInitialized {
			rejectErr(reject, errors.New("client must be initialized first"))
			return
		}

		if err := checkInputLen(inputs, 1); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenPassword := inputs[0]

		if err := checkIsString(chosenPassword, "password"); err != nil {
			rejectErr(reject, err)
			return
		}

		regState, regReq, err := c.c.CreateRegistrationRequest([]byte(chosenPassword.String()))
		if err != nil {
			rejectErr(reject, err)
			return
		}

		encodedRegReq, err := regReq.Encode()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		c.registrationState = regState

		// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
		dataJS := copyBytesToJS(encodedRegReq)

		resolve.Invoke(dataJS)
	}

	return promiser(runner)
}

// RegistrationFinalize wasm wrapper for opaque.Client.FinalizeRegistrationRequest
// Takes two argument, first one is string and second one is []byte
// returns []byte for registration record and []byte for exportKey
// Prototype Go: RegistrationFinalize(clientIdentity string, registrationRes []byte) ([]byte, []byte)
// Prototype JS: registrationFinalize(clientIdentity: string, registrationRes: Uint8Array) Object(Uint8Array, Uint8Array)
func (c *client) RegistrationFinalize(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !c.isInitialized {
			rejectErr(reject, errors.New("client must be initialized first"))
			return
		}

		if c.registrationState == nil {
			rejectErr(reject, errors.New("registrationInit must be executed first"))
			return
		}

		if err := checkInputLen(inputs, 2); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenClientIdentity := inputs[0]
		chosenRegRes := inputs[1]

		if err := checkIsString(chosenClientIdentity, "clientIdentity"); err != nil {
			rejectErr(reject, err)
			return
		}

		regRes, err := copyBytesToGo(chosenRegRes, "registrationRes")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		regRecord, exportKey, err := c.c.FinalizeRegistrationRequest(c.registrationState, []byte(chosenClientIdentity.String()), regRes)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		encodedRegRec, err := regRecord.Encode()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
		regRecJS := copyBytesToJS(encodedRegRec)
		exportKeyJS := copyBytesToJS(exportKey)

		returnObj := make(map[string]interface{})
		returnObj["record"] = regRecJS
		returnObj["exportKey"] = exportKeyJS

		// zero out state
		c.registrationState = nil

		resolve.Invoke(returnObj)
	}

	return promiser(runner)
}

// LoginInit wasm wrapper for opaque.Client.ClientInit
// Takes one argument and it is string, returns []byte for ke1 message
// Prototype Go: LoginInit(password string) []byte
// Prototype JS: loginInit(password: string) Uint8Array
func (c *client) LoginInit(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !c.isInitialized {
			rejectErr(reject, errors.New("client must be initialized first"))
			return
		}

		if err := checkInputLen(inputs, 1); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenPassword := inputs[0]

		if err := checkIsString(chosenPassword, "password"); err != nil {
			rejectErr(reject, err)
			return
		}

		logState, ke1Message, err := c.c.ClientInit([]byte(chosenPassword.String()))
		if err != nil {
			rejectErr(reject, err)
			return
		}

		encodedKe1Message, err := ke1Message.Encode()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		c.loginState = logState

		// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
		dataJS := copyBytesToJS(encodedKe1Message)

		resolve.Invoke(dataJS)
	}
	return promiser(runner)
}

// LoginFinish wasm wrapper for opaque.Client.ClientFinish
// Takes two argument, first one is string and second one is []byte
// returns []byte for ke3 message, []byte for sessionKey and []byte for exportKey
// Prototype Go: LoginFinish(clientIdentity string, ke2Message []byte) ([]byte, []byte, []byte)
// Prototype JS: loginFinish(clientIdentity: string, ke2Message Uint8Array) Object(Uint8Array, Uint8Array, Uint8Array)
func (c *client) LoginFinish(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !c.isInitialized {
			rejectErr(reject, errors.New("client must be initialized first"))
			return
		}

		if c.loginState == nil {
			rejectErr(reject, errors.New("loginInit must be executed first"))
			return
		}
		if err := checkInputLen(inputs, 2); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenClientIdentity := inputs[0]
		chosenKe2Message := inputs[1]

		if err := checkIsString(chosenClientIdentity, "clientIdentity"); err != nil {
			rejectErr(reject, err)
			return
		}

		ke2Message, err := copyBytesToGo(chosenKe2Message, "ke2Message")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		ke3Message, sessionKey, exportKey, err := c.c.ClientFinish(c.loginState, []byte(chosenClientIdentity.String()), ke2Message)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		encodedKe3Message, err := ke3Message.Encode()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
		ke3JS := copyBytesToJS(encodedKe3Message)
		sessionKeyJS := copyBytesToJS(sessionKey)
		exportKeyJS := copyBytesToJS(exportKey)

		returnObj := make(map[string]interface{})
		returnObj["ke3"] = ke3JS
		returnObj["sessionKey"] = sessionKeyJS
		returnObj["exportKey"] = exportKeyJS

		// zero out state
		c.loginState = nil

		resolve.Invoke(returnObj)
	}

	return promiser(runner)
}

func (c *client) IsInitialized(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		resolve.Invoke(c.isInitialized)
	}
	return promiser(runner)
}

// InitializeClient wasm wrapper for opaque.NewClient
// Takes two argument, both are string, returns nothing. But resolve promise if successful.
// Prototype Go: InitializeClient(suiteName string, serverID string)
// Prototype JS: initClient(suiteName: string, serverID: string)
func (c *client) InitializeClient(this js.Value, inputs []js.Value) any {
	cConf := &opaque.ClientConfiguration{}

	runner := func(resolve js.Value, reject js.Value) {
		if err := checkInputLen(inputs, 2); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenSuite := inputs[0]
		chosenServerID := inputs[1]

		if err := checkIsString(chosenSuite, "suiteName"); err != nil {
			rejectErr(reject, err)
			return
		}

		suiteID, err := strToSuite(chosenSuite.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenServerID, "serverID"); err != nil {
			rejectErr(reject, err)
			return
		}

		cConf.OpaqueSuite = suiteID.New()
		cConf.ServerID = []byte(chosenServerID.String())

		c.c = opaque.NewClient(cConf)
		c.isInitialized = true
		resolve.Invoke(true)
	}
	return promiser(runner)
}
