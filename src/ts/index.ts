import { initializeWasm } from "./modules/wasm";
import { Client } from './modules/client';
import { Server } from "./modules/server";

// always wait for wasm initialized
await initializeWasm();

const suiteName = "Ristretto255Suite"
const serverID = "example.com"
const clientPassword = "SuperSecurePass"
const credID = "UniqueCredentialID"
const clientIdentity = "UniqueClientID"


try {
    /////// Registration Steps

    // Client Initialization
    const cl = new Client();
    await cl.initClient({ suiteName: suiteName, serverID: serverID }).then((res) => {
        console.log("Client Initialized:", res);
    })
    await cl.isInitialized().then((res) => { console.log("Is Client initialized:", res) });

    // Server Initialization
    const sv = new Server();
    await sv.initServer({ suiteName: suiteName, serverID: serverID, privateKey: null }).then((res) => {
        console.log("Server Initialized:", res);
    })
    await sv.isInitialized().then((res) => { console.log("Is Server initialized:", res) })

    // Generate Random Oprf Seed
    const oprfSeed: Uint8Array = await sv.generateOprfSeed().then((seed) => {
        console.log("Generated Oprf Seed: ", seed)
        return seed
    })

    // Client Registration Init
    const regReq: Uint8Array = await cl.registrationInit(clientPassword).then((req) => {
        console.log("Generated Registration Request: ", req);
        return req
    });

    // Server Registration Response
    const regRes: Uint8Array = await sv.registrationRes(regReq, credID, oprfSeed).then((regRes) => {
        console.log("Generated Registration Response: ", regRes)
        return regRes
    })

    // Client Registration Finalize
    const regFinalize = await cl.registrationFinalize(clientIdentity, regRes).then(({ record, exportKey }) => {
        console.log("Generated Registration Record: ", record)
        console.log("Generated Registration Export Key: ", exportKey)
        return { ...{ record, exportKey } }
    })

    /////// Login Steps

    // Client Login Init
    const ke1: Uint8Array = await cl.loginInit(clientPassword).then((ke1Message) => {
        console.log("Generated KE1 Message:", ke1Message)
        return ke1Message
    })

    // Server Login Init
    const ke2: Uint8Array = await sv.loginInit(regFinalize.record, ke1, credID, clientIdentity, oprfSeed).then((ke2Message) => {
        console.log("Generated KE2 Message:", ke2Message)
        return ke2Message
    })

    // Client Login Finish
    const logClFinish = await cl.loginFinish(clientIdentity, ke2).then(({ ke3, sessionKey, exportKey }) => {
        console.log("Generated KE3 Message:", ke3)
        console.log("Generated Client Session Key:", sessionKey)
        console.log("Generated Login Export Key: ", exportKey)
        return { ...{ ke3, sessionKey, exportKey } }
    })

    const serverSessionKey: Uint8Array = await sv.loginFinish(logClFinish.ke3).then((sessionKey) => {
        console.log("Generated Server Session Key:", sessionKey)
        return sessionKey
    })

} catch (err) {
    console.log(err)
}










