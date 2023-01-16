<p align="center">
    <img width="450" src="assets/images/logo.png">
</p>

## About

This library implements the OPAQUE key exchange protocol using Webassembly.

The implementation is based on [Cryptomony](https://github.com/cymony/cryptomony)

## Installation
```sh
npm install --save @cymony/cryptomonyjs-opaque
```

## Example Usage
```js
import { Client, Server, initializeWasm } from '@cymony/cryptomonyjs-opaque';

(async () => {
    // Configuration Constants
    const clientPassword = "SuperSecurePassword";
    const clientIdentity = "example@example.com";
    const serverID = "example.com";
    const suiteName = "Ristretto255Suite";
    const credentialID = "UniqueCredIdentifier";

    //// Initialization Part
    console.log("*".repeat(10) + " Initialization Part " + "*".repeat(10));

    // run wasm / initialize wasm
    await initializeWasm().then(() => {
        console.log("Wasm initialized !!");
    });

    // create new client
    const client = new Client();
    console.log("New client created !!");

    // create new server
    const server = new Server();
    console.log("New server created !!");

    // check client initialized
    await client.isInitialized().then((result) => {
        console.log("client.isInitialized Result: ", result);
    });

    // initialize client / load configuration
    await client
        .initClient({
            suiteName: suiteName,
            serverID: serverID,
        })
        .then(() => {
            console.log("Client initialized !!");
        });

    // check client initialized
    await client.isInitialized().then((result) => {
        console.log("client.isInitialized Result: ", result);
    });

    // check server initialized
    await server.isInitialized().then((result) => {
        console.log("server.isInitialized Result: ", result);
    });

    // initialize server / load configuration
    await server
        .initServer({
            suiteName: suiteName,
            serverID: serverID,
        })
        .then(() => {
            console.log("Server initialized !!");
        });

    // check server initialized
    await server.isInitialized().then((result) => {
        console.log("server.isInitialized Result: ", result);
    });

    //// Registration Part
    console.log("*".repeat(10) + " Registration Part " + "*".repeat(10));

    // client registration init
    const { registrationState: clRegistrationState, registrationRequest } =
        await client.registrationInit(clientPassword).then((obj) => {
            console.log("RegistrationState: ", obj.registrationState);
            console.log("RegistrationRequest: ", obj.registrationRequest);
            return obj;
        });

    // server generate oprf key
    const oprfSeed = await server.generateOprfSeed().then((seed) => {
        console.log("Generated oprf seed: ", seed);
        return seed;
    });

    // server registration evulation
    const registrationResponse = await server
        .registrationEval(registrationRequest, oprfSeed, credentialID)
        .then((registrationResponse) => {
            console.log("Registration Response: ", registrationResponse);
            return registrationResponse;
        });

    // client registration finalize
    const { registrationRecord, exportKey: registrationExportKey } =
        await client
            .registrationFinalize(
                clRegistrationState,
                registrationResponse,
                clientIdentity
            )
            .then((obj) => {
                console.log("Registration Record: ", obj.registrationRecord);
                console.log("Registration ExportKey: ", obj.exportKey);
                return obj;
            });

    //// Login Part
    console.log("*".repeat(10) + " Login Part " + "*".repeat(10));

    // client login init
    const { loginState: clLoginState, ke1 } = await client
        .loginInit(clientPassword)
        .then((obj) => {
            console.log("Client Login State: ", obj.loginState);
            console.log("KE1: ", obj.ke1);
            return obj;
        });

    // server login init
    const { loginState: svLoginState, ke2 } = await server
        .loginInit(
            registrationRecord,
            ke1,
            oprfSeed,
            credentialID,
            clientIdentity
        )
        .then((obj) => {
            console.log("Server Login State: ", obj.loginState);
            console.log("KE2: ", obj.ke2);
            return obj;
        });

    // client login finish
    const {
        exportKey: loginExportKey,
        sessionKey: clSessionKey,
        ke3,
    } = await client
        .loginFinish(clLoginState, ke2, clientIdentity)
        .then((obj) => {
            console.log("Client Session Key: ", obj.sessionKey);
            console.log("Login Export Key: ", obj.exportKey);
            console.log("KE3: ", obj.ke3);
            return obj;
        });

    // server login finish
    const svSessionKey = await server
        .loginFinish(svLoginState, ke3)
        .then((svSessionKey) => {
            console.log("Server Session Key:", svSessionKey);
            return svSessionKey;
        });

    console.log("*".repeat(10) + " Equality Part " + "*".repeat(10));
    console.log(
        "Is Session Keys Equal: ",
        JSON.stringify(svSessionKey) === JSON.stringify(clSessionKey)
    );
    console.log(
        "Is Registration and Login Export Keys Equal: ",
        JSON.stringify(registrationExportKey) ===
        JSON.stringify(loginExportKey)
    );
})();
```

## License
This project is licensed under the [BSD 3-Clause](./LICENSE)