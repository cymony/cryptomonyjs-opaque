import { getWasmServer, Suite } from '../consts'

export interface ServerConfiguration {
    suiteName: Suite
    serverID: string
    privateKey: Uint8Array | null
}

export class Server {
    constructor(conf?: ServerConfiguration) {
        if (conf) {
            this.initServer(conf);
        }
    }

    initServer(conf: ServerConfiguration): Promise<boolean> {
        const wasmSv = getWasmServer();
        return wasmSv.initServer(conf.suiteName, conf.serverID, conf.privateKey);
    }

    isInitialized(): Promise<boolean> {
        const wasmCl = getWasmServer();
        return wasmCl.isInitialized();
    }

    /**
     * generateOprfSeed generates
     * @returns Promise<Uint8array>
     */
    generateOprfSeed(): Promise<Uint8Array> {
        const wasmCl = getWasmServer();
        return wasmCl.generateOprfSeed();
    }

    registrationRes(regReq: Uint8Array, credID: string, oprfSeed: Uint8Array): Promise<Uint8Array> {
        const wasmCl = getWasmServer();
        return wasmCl.registrationRes(regReq, credID, oprfSeed);
    }

    loginInit(record: Uint8Array, ke1Message: Uint8Array, credID: string, clientIdentity: string, oprfSeed: Uint8Array): Promise<Uint8Array> {
        const wasmCl = getWasmServer();
        return wasmCl.loginInit(record, ke1Message, credID, clientIdentity, oprfSeed);
    }

    loginFinish(ke3Message: Uint8Array): Promise<Uint8Array> {
        const wasmCl = getWasmServer();
        return wasmCl.loginFinish(ke3Message);
    }
}
