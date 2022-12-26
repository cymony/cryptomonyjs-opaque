import { getWasmClient, Suite } from '../consts'

export interface ClientConfiguration {
    suiteName: Suite
    serverID: string
}

export class Client {
    constructor(conf?: ClientConfiguration) {
        if (conf) {
            this.initClient(conf);
        }
    }

    initClient(conf: ClientConfiguration): Promise<boolean> {
        const wasmCl = getWasmClient();
        return wasmCl.initClient(conf.suiteName, conf.serverID);
    }

    isInitialized(): Promise<boolean> {
        const wasmCl = getWasmClient();
        return wasmCl.isInitialized();
    }

    registrationInit(password: string): Promise<Uint8Array> {
        const wasmCl = getWasmClient();
        return wasmCl.registrationInit(password);
    }

    registrationFinalize(clientIdentity: string, registrationRes: Uint8Array): Promise<{
        record: Uint8Array;
        exportKey: Uint8Array;
    }> {
        const wasmCl = getWasmClient();
        return wasmCl.registrationFinalize(clientIdentity, registrationRes);
    }

    loginInit(password: string): Promise<Uint8Array> {
        const wasmCl = getWasmClient();
        return wasmCl.loginInit(password);
    }

    loginFinish(clientIdentity: string, ke2Message: Uint8Array): Promise<{
        ke3: Uint8Array
        sessionKey: Uint8Array
        exportKey: Uint8Array
    }> {
        const wasmCl = getWasmClient();
        return wasmCl.loginFinish(clientIdentity, ke2Message);
    }
}
