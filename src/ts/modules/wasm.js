import * as _ from "../../api/wasm_exec.js";
import libwasm from "../../api/lib.wasm";

console.log(libwasm);

export const initializeWasm = async () => {
    const go = new Go();
    await libwasm({ ...go.importObject }).then(({ instance }) => {
        go.run(instance)
    });
}
