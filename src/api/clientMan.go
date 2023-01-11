package main

import (
	"errors"
	"math/rand"
	"syscall/js"
	"time"
	"unsafe"
)

type clientManager struct {
	clients map[string]*client
	rndSrc  rand.Source
}

func newClientManager() *clientManager {
	src := rand.NewSource(time.Now().UnixNano())

	return &clientManager{
		clients: make(map[string]*client),
		rndSrc:  src,
	}
}

func (cm *clientManager) exposeToJS(rootModule js.Value) {
	rootModule.Set("client", make(map[string]interface{}))
	clientModule := rootModule.Get("client")

	clientModule.Set("newClient", js.FuncOf(cm.NewClient))
	clientModule.Set("initClient", js.FuncOf(cm.InitClient))
	clientModule.Set("isInitialized", js.FuncOf(cm.IsInitialized))
	clientModule.Set("registrationInit", js.FuncOf(cm.RegistrationInit))
	clientModule.Set("registrationFinalize", js.FuncOf(cm.RegistrationFinalize))
	clientModule.Set("loginInit", js.FuncOf(cm.LoginInit))
	clientModule.Set("loginFinish", js.FuncOf(cm.LoginFinish))
}

// NewClient creates new emptys client instance with identifier. It returns the identifier.
func (cm *clientManager) NewClient(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		clid := cm.GenerateRandomID()
		cm.clients[clid] = newClient()
		resolve.Invoke(clid)
	}
	return promiser(runner)
}

// InitClient initializes the already existing client instance with configuration.
func (cm *clientManager) InitClient(this js.Value, inputs []js.Value) any {

	runner := func(resolve js.Value, reject js.Value) {
		cl, err := cm.getClient(inputs, 3)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenSuiteName := inputs[1]
		chosenServerID := inputs[2]

		if err := checkIsString(chosenSuiteName, "suiteName"); err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenServerID, "serverID"); err != nil {
			rejectErr(reject, err)
			return
		}

		err = cl.InitializeClient(chosenSuiteName.String(), chosenServerID.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}
		resolve.Invoke()
	}

	return promiser(runner)
}

// IsInitialized returns whether the client has been initialized
func (cm *clientManager) IsInitialized(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		cl, err := cm.getClient(inputs, 1)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		resolve.Invoke(cl.IsInitialized())
	}
	return promiser(runner)
}

func (cm *clientManager) RegistrationInit(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		cl, err := cm.getClient(inputs, 2)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenPassword := inputs[1]
		if err := checkIsString(chosenPassword, "password"); err != nil {
			rejectErr(reject, err)
			return
		}

		regState, regReq, err := cl.RegistrationInit(chosenPassword.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		returnObj := make(map[string]interface{})
		returnObj["registrationState"] = copyBytesToJS(regState)
		returnObj["registrationRequest"] = copyBytesToJS(regReq)

		resolve.Invoke(returnObj)
	}
	return promiser(runner)
}

func (cm *clientManager) RegistrationFinalize(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		cl, err := cm.getClient(inputs, 4)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenRegState := inputs[1]
		chosenRegRes := inputs[2]
		chosenClientIdentity := inputs[3]

		regState, err := copyBytesToGo(chosenRegState, "registrationState")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		regRes, err := copyBytesToGo(chosenRegRes, "registrationResponse")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenClientIdentity, "clientIdentity"); err != nil {
			rejectErr(reject, err)
			return
		}

		regRecord, exportKey, err := cl.RegistrationFinalize(regState, regRes, chosenClientIdentity.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		returnObj := make(map[string]interface{})
		returnObj["registrationRecord"] = copyBytesToJS(regRecord)
		returnObj["exportKey"] = copyBytesToJS(exportKey)

		resolve.Invoke(returnObj)
	}
	return promiser(runner)
}

func (cm *clientManager) LoginInit(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		cl, err := cm.getClient(inputs, 2)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenPassword := inputs[1]
		if err := checkIsString(chosenPassword, "password"); err != nil {
			rejectErr(reject, err)
			return
		}

		loginState, ke1, err := cl.LoginInit(chosenPassword.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		returnObj := make(map[string]interface{})
		returnObj["loginState"] = copyBytesToJS(loginState)
		returnObj["ke1"] = copyBytesToJS(ke1)

		resolve.Invoke(returnObj)
	}
	return promiser(runner)
}

func (cm *clientManager) LoginFinish(this js.Value, inputs []js.Value) any {
	runner := func(resolve js.Value, reject js.Value) {
		cl, err := cm.getClient(inputs, 4)
		if err != nil {
			rejectErr(reject, err)
			return
		}

		chosenLoginState := inputs[1]
		chosenKE2 := inputs[2]
		chosenClientIdentity := inputs[3]

		loginState, err := copyBytesToGo(chosenLoginState, "loginState")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		ke2Message, err := copyBytesToGo(chosenKE2, "ke2Message")
		if err != nil {
			rejectErr(reject, err)
			return
		}

		if err := checkIsString(chosenClientIdentity, "clientIdentity"); err != nil {
			rejectErr(reject, err)
			return
		}

		ke3Message, sessionKey, exportKey, err := cl.LoginFinish(loginState, ke2Message, chosenClientIdentity.String())
		if err != nil {
			rejectErr(reject, err)
			return
		}

		returnObj := make(map[string]interface{})
		returnObj["ke3"] = copyBytesToJS(ke3Message)
		returnObj["sessionKey"] = copyBytesToJS(sessionKey)
		returnObj["exportKey"] = copyBytesToJS(exportKey)

		resolve.Invoke(returnObj)
	}
	return promiser(runner)
}

func (cm *clientManager) getClient(inputs []js.Value, inputLen int) (*client, error) {
	if err := checkInputLen(inputs, inputLen); err != nil {
		return nil, err
	}

	clIdentifier := inputs[0]

	if err := checkIsString(clIdentifier, "clientID"); err != nil {
		return nil, err
	}

	cl, ok := cm.clients[clIdentifier.String()]
	if !ok {
		return nil, errors.New("client not found")
	}

	return cl, nil
}

// GenerateRandomID generates random and unique client identifier
func (cm *clientManager) GenerateRandomID() string {
	n := 8

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, cm.rndSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = cm.rndSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	generatedString := *(*string)(unsafe.Pointer(&b))
	_, ok := cm.clients[generatedString]
	if ok {
		return cm.GenerateRandomID()
	}
	return generatedString
}
