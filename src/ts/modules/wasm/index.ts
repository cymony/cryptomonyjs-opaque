import "../../../api/wasm_exec";
//@ts-ignore
import libwasm from "../../../api/lib.wasm";

export const initializeWasm = async () => {
    //@ts-ignore
    const go = new Go();
    await libwasm({ ...go.importObject }).then(({ instance }) => {
        go.run(instance)
    });
}
