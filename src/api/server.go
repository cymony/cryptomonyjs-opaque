package main

import (
	"errors"
	"syscall/js"

	"github.com/cymony/cryptomony/opaque"
)

type server struct {
	isInitialized bool
	loginState    *opaque.ServerLoginState
	s             opaque.Server
}

func newServer() *server {
	return &server{isInitialized: false, loginState: nil, s: nil}
}

func (s *server) exposeServer(rootModule js.Value) {
	rootModule.Set("server", make(map[string]interface{}))
	serverModule := rootModule.Get("server")

	serverModule.Set("initServer", js.FuncOf(s.InitializeServer))
	serverModule.Set("isInitialized", js.FuncOf(s.IsInitialized))
	serverModule.Set("generateOprfSeed", js.FuncOf(s.GenerateOprfSeed))
	serverModule.Set("registrationRes", js.FuncOf(s.RegistrationRes))
	serverModule.Set("loginInit", js.FuncOf(s.LoginInit))
	serverModule.Set("loginFinish", js.FuncOf(s.LoginFinish))
}

// LoginFinish wasm wrapper for opaque.Server.ServerFinish
// Takes one argument, ke3Message []byte, returns sessionKey []byte
// Prototype Go: LoginFinish(ke3Message []byte) []byte
// Prototype JS: loginFinish(ke3Message: Uint8Array) Uint8Array
func (s *server) LoginFinish(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !s.isInitialized {
			rejectErr(reject, errors.New("server must be initialized first"))
			return
		}

		if s.loginState == nil {
			rejectErr(reject, errors.New("loginInit must be executed first"))
			return
		}

		if err := checkInputLen(inputs, 1); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenKE3 := inputs[0]
		ke3, err := copyBytesToGo(chosenKE3, "ke3Message")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		sessionKey, err := s.s.ServerFinish(s.loginState, ke3)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		dataJS := copyBytesToJS(sessionKey)
		resolve.Invoke(dataJS)
	}
	return promiser(runner)
}

// LoginInit wasm wrapper for opaque.Server.ServerInit
// Takes five argument, clRecord []byte, ke1Message []byte, credentialIdentifier string, clientIdentity string, oprfSeed []byte
// returns ke2Message []byte
// Prototype Go: LoginInit(record []byte, ke1 []byte, credID string, clientIdentity string, oprfSeed []byte) []byte
// Prototype JS: loginInit(record: Uint8Array, ke1Message: Uint8Array, credID: string, clientIdentity: string, oprfSeed: Uint8Array) Uint8Array
func (s *server) LoginInit(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !s.isInitialized {
			rejectErr(reject, errors.New("server must be initialized first"))
			return
		}

		if err := checkInputLen(inputs, 5); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenRecord := inputs[0]
		record, err := copyBytesToGo(chosenRecord, "record")
		if err != nil {
			rejectErr(reject, err)
			return
		}
		chosenKe1 := inputs[1]
		ke1, err := copyBytesToGo(chosenKe1, "ke1Message")
		if err != nil {
			rejectErr(reject, err)
			return
		}
		chosenCredID := inputs[2]
		if err := checkIsString(chosenCredID, "credID"); err != nil {
			rejectErr(reject, err)
			return
		}
		chosenClientID := inputs[3]
		if err := checkIsString(chosenClientID, "clientIdentity"); err != nil {
			rejectErr(reject, err)
			return
		}
		chosenOprfSeed := inputs[4]
		orpfSeed, err := copyBytesToGo(chosenOprfSeed, "oprfSeed")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		loginState, ke2, err := s.s.ServerInit(record, ke1, []byte(chosenCredID.String()), []byte(chosenClientID.String()), orpfSeed)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		encodedKE2, err := ke2.Encode()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		dataJS := copyBytesToJS(encodedKE2)

		s.loginState = loginState
		resolve.Invoke(dataJS)
	}

	return promiser(runner)
}

// GenerateOprfSeed wasm wrapper for opaque.Server.GenerateOprfSeed
// Takes no argument, returns []byte
// Prototype Go: GenerateOprfSeed() []byte
// Prototype JS: generateOprfSeed() Uint8Array
func (s *server) GenerateOprfSeed(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !s.isInitialized {
			rejectErr(reject, errors.New("server must be initialized first"))
			return
		}

		oprfSeed := s.s.GenerateOprfSeed()

		dataJS := copyBytesToJS(oprfSeed)
		resolve.Invoke(dataJS)
	}
	return promiser(runner)
}

// RegistrationRes wasm wrapper for opaque.Server.CreateRegistrationResponse
// Takes three argument, registrationRequest []byte, credentialID string and oprfSeed []byte
// returns registrationResponse []byte
// Prototype Go: RegistrationRes(regReq []byte, credID string, oprfSeed []byte) []byte
// Prototype JS: registrationRes(regReq: Uint8Array, credID: string, oprfSeed: Uint8Array) Uint8Array
func (s *server) RegistrationRes(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		if !s.isInitialized {
			rejectErr(reject, errors.New("server must be initialized first"))
			return
		}

		if err := checkInputLen(inputs, 3); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenRegReq := inputs[0]
		regReq, err := copyBytesToGo(chosenRegReq, "regReq")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenCredID := inputs[1]
		if err := checkIsString(chosenCredID, "credID"); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenOprfSeed := inputs[2]
		oprfSeed, err := copyBytesToGo(chosenOprfSeed, "oprfSeed")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		regRes, err := s.s.CreateRegistrationResponse(regReq, []byte(chosenCredID.String()), oprfSeed)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		encodedRegRes, err := regRes.Encode()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		dataJS := copyBytesToJS(encodedRegRes)

		resolve.Invoke(dataJS)
	}
	return promiser(runner)
}

func (s *server) IsInitialized(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		resolve.Invoke(s.isInitialized)
	}

	return promiser(runner)
}

// InitializeServer wasm wrapper for opaque.NewServer
// Takes three argument, first two are strign and last one []byte, returns nothing. But resolve promise if successful.
// Prototype Go: InitializeServer(suiteName string, serverID string, privKey []byte)
// Prototype JS: initClient(suiteName: string, serverID: string, privKey: Uint8Array)
func (s *server) InitializeServer(this js.Value, inputs []js.Value) any {
	sConf := &opaque.ServerConfiguration{}

	runner := func(resolve js.Value, reject js.Value) {
		if err := checkInputLen(inputs, 3); err != nil {
			rejectErr(reject, err)
			return
		}

		chosenSuite := inputs[0]

		if err := checkIsString(chosenSuite, "suiteName"); err != nil {
			rejectErr(reject, err)
			return
		}

		ss, err := strToSuite(chosenSuite.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}
		sConf.OpaqueSuite = ss

		chosenServerID := inputs[1]
		if err := checkIsString(chosenServerID, "serverID"); err != nil {
			rejectErr(reject, err)
			return
		}
		sConf.ServerID = []byte(chosenServerID.String())

		chosenPrivKey := inputs[2]
		if chosenPrivKey.IsNull() {
			sConf.ServerPrivateKey = nil
		} else {
			privKey, err := copyBytesToGo(chosenPrivKey, "privKey")
			if err != nil {
				rejectErr(reject, err)
				return
			}
			sConf.ServerPrivateKey = privKey
		}

		server, err := opaque.NewServer(sConf)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		s.s = server
		s.isInitialized = true

		resolve.Invoke(true)
	}
	return promiser(runner)
}
