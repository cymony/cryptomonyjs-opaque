const wasmRootEl: string = "__cryptomonyjsopaque__";
const clientRootEl: string = "client";
const serverRootEl: string = "server";

export type Suite = 'Ristretto255Suite' | 'P256Suite'

export const isNode = typeof process !== "undefined" && process.versions != null &&
    process.versions.node != null;

export const getWasmClient = () => {
    return globalThis[wasmRootEl][clientRootEl]
}

export const getWasmServer = () => {
    return globalThis[wasmRootEl][serverRootEl]
}
