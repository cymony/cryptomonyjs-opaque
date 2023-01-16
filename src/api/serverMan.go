package main

import (
	"errors"
	"math/rand"
	"syscall/js"
	"time"
	"unsafe"
)

type serverManager struct {
	servers map[string]*server
	rndSrc  rand.Source
}

func newServerManager() *serverManager {
	src := rand.NewSource(time.Now().UnixNano())

	return &serverManager{
		servers: make(map[string]*server),
		rndSrc:  src,
	}
}

func (sm *serverManager) exposeServer(rootModule js.Value) {
	rootModule.Set("server", make(map[string]interface{}))
	serverModule := rootModule.Get("server")

	serverModule.Set("newServer", js.FuncOf(sm.NewServer))
	serverModule.Set("initServer", js.FuncOf(sm.InitializeServer))
	serverModule.Set("isInitialized", js.FuncOf(sm.IsInitialized))
	serverModule.Set("generateOprfSeed", js.FuncOf(sm.GenerateOprfSeed))
	serverModule.Set("registrationEval", js.FuncOf(sm.RegistrationEval))
	serverModule.Set("loginInit", js.FuncOf(sm.LoginInit))
	serverModule.Set("loginFinish", js.FuncOf(sm.LoginFinish))
}

// NewServer creates new empty server instance with identifier. It returns the identifier.
func (sm *serverManager) NewServer(this js.Value, inputs []js.Value) any {
	clid := sm.GenerateRandomID()
	sm.servers[clid] = newServer()
	return clid
}

// initServer(identifier: string, suiteName: string, serverID: string, privKey: Uint8Array) Promise<void>
func (sm *serverManager) InitializeServer(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		sv, err := sm.getServer(inputs, 4)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenSuite := inputs[1]
		chosenServerID := inputs[2]
		chosenPrivKey := inputs[3]

		if err := checkIsString(chosenSuite, "suiteName"); err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenServerID, "serverID"); err != nil {
			rejectErr(reject, err)
			return
		}

		var privKey []byte

		if chosenPrivKey.IsNull() || chosenPrivKey.IsUndefined() || chosenPrivKey.IsNaN() {
			privKey = nil
		} else {
			privKey, err = copyBytesToGo(chosenPrivKey, "privKey")
			if err != nil {
				rejectErr(reject, err)
				return
			}
		}

		if err := sv.InitializeServer(chosenSuite.String(), chosenServerID.String(), privKey); err != nil {
			rejectErr(reject, err)
			return
		}

		resolve.Invoke()
	}

	return promiser(runner)
}

// isInitialized(identifier: string) Promise<boolean>
func (sm *serverManager) IsInitialized(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		sv, err := sm.getServer(inputs, 1)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		resolve.Invoke(sv.IsInitialized())
	}

	return promiser(runner)
}

// generateOprfSeed(identifier: string) Promise<Uint8Array>
func (sm *serverManager) GenerateOprfSeed(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		sv, err := sm.getServer(inputs, 1)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		oprfSeed, err := sv.GenerateOprfSeed()
		if err != nil {
			rejectErr(reject, err)
			return
		}

		dataJS := copyBytesToJS(oprfSeed)
		resolve.Invoke(dataJS)
	}

	return promiser(runner)
}

// registrationEval(identifier: string, registrationRequest: Uint8Array, oprfSeed: Uint8Array, credentialIdentifier: string) Promise<Uint8Array>
func (sm *serverManager) RegistrationEval(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		sv, err := sm.getServer(inputs, 4)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenRegReq := inputs[1]
		chosenOprfSeed := inputs[2]
		chosenCredID := inputs[3]

		regRes, err := copyBytesToGo(chosenRegReq, "registrationRequest")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		oprfSeed, err := copyBytesToGo(chosenOprfSeed, "oprfSeed")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenCredID, "credentialIdentifier"); err != nil {
			rejectErr(reject, err)
			return
		}

		regResponse, err := sv.RegistrationEval(regRes, oprfSeed, chosenCredID.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		dataJS := copyBytesToJS(regResponse)
		resolve.Invoke(dataJS)
	}

	return promiser(runner)
}

/*
* loginInit(identifier: string,
*   record: Uint8Array,
*   ke1: Uint8Array,
*   oprfSeed Uint8Array,
*   credentialID string,
*   clientIdentity string) Promise<{
*	loginState: Uint8Array,
*	ke2: Uint8Array}>
 */
func (sm *serverManager) LoginInit(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		sv, err := sm.getServer(inputs, 6)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenRecord := inputs[1]
		chosenKE1 := inputs[2]
		chosenOprfSeed := inputs[3]
		chosenCredID := inputs[4]
		chosenClientIdentity := inputs[5]

		record, err := copyBytesToGo(chosenRecord, "record")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		ke1, err := copyBytesToGo(chosenKE1, "ke1")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		oprfSeed, err := copyBytesToGo(chosenOprfSeed, "oprfSeed")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenCredID, "credentialID"); err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenClientIdentity, "clientIdentity"); err != nil {
			rejectErr(reject, err)
			return
		}

		loginState, ke2, err := sv.LoginInit(record, ke1, oprfSeed, chosenCredID.String(), chosenClientIdentity.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		returnObj := make(map[string]interface{})
		returnObj["loginState"] = copyBytesToJS(loginState)
		returnObj["ke2"] = copyBytesToJS(ke2)

		resolve.Invoke(returnObj)
	}

	return promiser(runner)
}

// loginFinish(identifier: string, loginState: Uint8Array, ke3: Uint8Array) Promise<Uint8Array>
func (sm *serverManager) LoginFinish(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		sv, err := sm.getServer(inputs, 3)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenLoginState := inputs[1]
		chosenKE3 := inputs[2]

		loginState, err := copyBytesToGo(chosenLoginState, "loginState")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		ke3, err := copyBytesToGo(chosenKE3, "ke3")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		sessionKey, err := sv.LoginFinish(loginState, ke3)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		dataJS := copyBytesToJS(sessionKey)
		resolve.Invoke(dataJS)
	}

	return promiser(runner)
}

func (sm *serverManager) getServer(inputs []js.Value, inputLen int) (*server, error) {
	if err := checkInputLen(inputs, inputLen); err != nil {
		return nil, err
	}

	svIdentifier := inputs[0]

	if err := checkIsString(svIdentifier, "identifier"); err != nil {
		return nil, err
	}

	sv, ok := sm.servers[svIdentifier.String()]
	if !ok {
		return nil, errors.New("server not found")
	}

	return sv, nil
}

// GenerateRandomID generates random and unique server identifier
func (sm *serverManager) GenerateRandomID() string {
	n := 8

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, sm.rndSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = sm.rndSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	generatedString := *(*string)(unsafe.Pointer(&b))
	_, ok := sm.servers[generatedString]
	if ok {
		return sm.GenerateRandomID()
	}
	return generatedString
}
