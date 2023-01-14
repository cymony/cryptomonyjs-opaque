package main

import (
	"errors"

	"github.com/cymony/cryptomony/opaque"
)

type server struct {
	isInitialized bool
	s             opaque.Server
	sConf         *opaque.ServerConfiguration
}

func newServer() *server {
	return &server{isInitialized: false, sConf: nil, s: nil}
}

// LoginFinish wasm wrapper for opaque.Server.ServerFinish
// Takes one argument, ke3Message []byte, returns sessionKey []byte
// Prototype Go: LoginFinish(ke3Message []byte) []byte
// Prototype JS: loginFinish(ke3Message: Uint8Array) Uint8Array
func (s *server) LoginFinish(loginState []byte, ke3 []byte) ([]byte, error) {
	if !s.IsInitialized() {
		return nil, errors.New("server must be initialized first")
	}

	svLoginState := &opaque.ServerLoginState{}
	if err := svLoginState.Decode(s.sConf.OpaqueSuite.New(), loginState); err != nil {
		return nil, err
	}

	sessionKey, err := s.s.ServerFinish(svLoginState, ke3)
	if err != nil {
		return nil, err
	}

	return sessionKey, nil
}

// LoginInit wasm wrapper for opaque.Server.ServerInit
func (s *server) LoginInit(record, ke1, oprfSeed []byte, credID, clientIdentity string) ([]byte, []byte, error) {
	if !s.IsInitialized() {
		return nil, nil, errors.New("server must be initialized first")
	}

	loginState, ke2, err := s.s.ServerInit(record, ke1, []byte(credID), []byte(clientIdentity), oprfSeed)
	if err != nil {
		return nil, nil, err
	}

	encodedLoginState, err := loginState.Encode()
	if err != nil {
		return nil, nil, err
	}

	encodedKE2, err := ke2.Encode()
	if err != nil {
		return nil, nil, err
	}

	return encodedLoginState, encodedKE2, nil
}

// RegistrationRes wasm wrapper for opaque.Server.CreateRegistrationResponse
func (s *server) RegistrationEval(regRequest, oprfSeed []byte, credID string) ([]byte, error) {
	if !s.IsInitialized() {
		return nil, errors.New("server must be initialized first")
	}

	regResponse, err := s.s.CreateRegistrationResponse(regRequest, []byte(credID), oprfSeed)
	if err != nil {
		return nil, err
	}

	encodedRegRes, err := regResponse.Encode()
	if err != nil {
		return nil, err
	}
	return encodedRegRes, nil
}

// GenerateOprfSeed wasm wrapper for opaque.Server.GenerateOprfSeed
func (s *server) GenerateOprfSeed() ([]byte, error) {
	if !s.IsInitialized() {
		return nil, errors.New("server must be initialized first")
	}

	oprfSeed := s.s.GenerateOprfSeed()
	return oprfSeed, nil
}

func (s *server) IsInitialized() bool {
	return s.isInitialized
}

// InitializeServer wasm wrapper for opaque.NewServer
func (s *server) InitializeServer(suiteName, serverID string, privKey []byte) error {
	sConf := &opaque.ServerConfiguration{}

	suiteID, err := strToSuite(suiteName)
	if err != nil {
		return err
	}

	sConf.OpaqueSuite = suiteID
	sConf.ServerID = []byte(serverID)
	sConf.ServerPrivateKey = privKey

	sv, err := opaque.NewServer(sConf)
	if err != nil {
		return err
	}

	s.s = sv
	s.isInitialized = true
	s.sConf = sConf

	return nil
}
