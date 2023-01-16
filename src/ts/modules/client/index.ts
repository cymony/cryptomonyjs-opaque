import { getWasmClient, Suite } from '../consts'

export interface ClientConfiguration {
    suiteName: Suite
    serverID: string
}

export class Client {
    private _identifier: string = '';

    constructor() {
        const wasmCl = getWasmClient();
        let clid = wasmCl.newClient();
        this._identifier = clid;
    }

    private get identifier(): string {
        return this._identifier;
    }

    initClient(conf: ClientConfiguration): Promise<void> {
        const wasmCl = getWasmClient();
        return wasmCl.initClient(this.identifier, conf.suiteName, conf.serverID);
    }

    isInitialized(): Promise<boolean> {
        const wasmCl = getWasmClient();
        return wasmCl.isInitialized(this.identifier);
    }

    registrationInit(password: string): Promise<{ registrationState: Uint8Array, registrationRequest: Uint8Array }> {
        const wasmCl = getWasmClient();
        return wasmCl.registrationInit(this.identifier, password);
    }

    registrationFinalize(registrationState: Uint8Array, registrationRes: Uint8Array, clientIdentity: string): Promise<{
        registrationRecord: Uint8Array;
        exportKey: Uint8Array;
    }> {
        const wasmCl = getWasmClient();
        return wasmCl.registrationFinalize(this.identifier, registrationState, registrationRes, clientIdentity);
    }

    loginInit(password: string): Promise<{ loginState: Uint8Array, ke1: Uint8Array }> {
        const wasmCl = getWasmClient();
        return wasmCl.loginInit(this.identifier, password);
    }

    loginFinish(loginState: Uint8Array, ke2: Uint8Array, clientIdentity: string): Promise<{
        ke3: Uint8Array
        sessionKey: Uint8Array
        exportKey: Uint8Array
    }> {
        const wasmCl = getWasmClient();
        return wasmCl.loginFinish(this.identifier, loginState, ke2, clientIdentity);
    }
}
