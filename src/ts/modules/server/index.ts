import { getWasmServer, Suite } from '../consts'

export interface ServerConfiguration {
    suiteName: Suite
    serverID: string
    privateKey: Uint8Array | null
}

export class Server {
    private _identifier: string = "";

    constructor() {
        const wasmSv = getWasmServer();
        let svID = wasmSv.newServer();
        this._identifier = svID;
    }

    private get identifier(): string {
        return this._identifier;
    }

    initServer(conf: ServerConfiguration): Promise<void> {
        const wasmSv = getWasmServer();
        return wasmSv.initServer(this.identifier, conf.suiteName, conf.serverID, conf.privateKey);
    }

    isInitialized(): Promise<boolean> {
        const wasmSv = getWasmServer();
        return wasmSv.isInitialized(this.identifier);
    }

    /**
    * generateOprfSeed generates
    * @returns Promise<Uint8array>
    */
    generateOprfSeed(): Promise<Uint8Array> {
        const wasmSv = getWasmServer();
        return wasmSv.generateOprfSeed(this.identifier);
    }

    registrationEval(registrationRequest: Uint8Array, oprfSeed: Uint8Array, credentialIdentifier: string): Promise<Uint8Array> {
        const wasmSv = getWasmServer();
        return wasmSv.registrationEval(this.identifier, registrationRequest, oprfSeed, credentialIdentifier);
    }

    loginInit(record: Uint8Array, ke1: Uint8Array, oprfSeed: Uint8Array, credID: string, clientIdentity: string): Promise<{
        loginState: Uint8Array
        ke2: Uint8Array
    }> {
        const wasmSv = getWasmServer();
        return wasmSv.loginInit(this.identifier, record, ke1, oprfSeed, credID, clientIdentity);
    }

    loginFinish(loginState: Uint8Array, ke3: Uint8Array): Promise<Uint8Array> {
        const wasmSv = getWasmServer();
        return wasmSv.loginFinish(this.identifier, loginState, ke3);
    }
}
