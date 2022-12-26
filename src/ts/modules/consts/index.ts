const wasmRootEl: string = "__cryptomonyjsopaque__";
const clientRootEl: string = "client";
const serverRootEl: string = "server";

export type Suite = 'Ristretto255Suite' | 'P256Suite'

export const getWasmClient = () => {
    return window[wasmRootEl][clientRootEl]
}

export const getWasmServer = () => {
    return window[wasmRootEl][serverRootEl]
}
