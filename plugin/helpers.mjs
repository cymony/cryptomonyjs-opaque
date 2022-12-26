
export const HELPERS_ID = '\0gzipInlineWasmHelper.js';


export const getHelperModule = () => `
function _loadWasmModule (src, imports, ungzipper) {
    function _instantiateOrCompile(source, imports, stream) {
      var instantiateFunc = stream ? WebAssembly.instantiateStreaming : WebAssembly.instantiate;
      var compileFunc = stream ? WebAssembly.compileStreaming : WebAssembly.compile;
  
      if (imports) {
        return instantiateFunc(source, imports)
      } else {
        return compileFunc(source)
      }
    }
  
    
    var buf = null;

    var raw = globalThis.atob(src);
   
    
    var rawLength = raw.length;
    buf = new Uint8Array(new ArrayBuffer(rawLength));
    for(var i = 0; i < rawLength; i++) {
        buf[i] = raw.charCodeAt(i);
    }

    var ungzipedData = ungzipper(buf);

    return _instantiateOrCompile(ungzipedData, imports, false)
}
export { _loadWasmModule };
`