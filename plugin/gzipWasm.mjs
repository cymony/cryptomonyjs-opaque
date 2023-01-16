import * as fs from 'fs';
import { gzip } from 'pako';

const HELPERS_ID = '\0gzipInlineWasmHelper.js';

const getHelperModule = () => `
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
    
    var isNode = typeof process !== 'undefined' && process.versions != null && process.versions.node != null
    if (isNode) {
      buf = Buffer.from(src, 'base64')
    } else {
      var raw = globalThis.atob(src)
      var rawLength = raw.length
      buf = new Uint8Array(new ArrayBuffer(rawLength))
      for(var i = 0; i < rawLength; i++) {
        buf[i] = raw.charCodeAt(i)
      }
    }

    var ungzipedData = ungzipper(buf);

    return _instantiateOrCompile(ungzipedData, imports, false)
}
export { _loadWasmModule };
`

export function gzipInlineWasm() {
  return {
    name: 'rollup-plugin-gzip-inline-wasm',

    resolveId(id) {
      if (id === HELPERS_ID) {
        return id;
      }

      return null;
    },
    load(id) {
      if (id === HELPERS_ID) {
        return getHelperModule();
      }

      if (!/\.wasm$/.test(id)) {
        return null;
      }

      return Promise.all([fs.promises.stat(id), fs.promises.readFile(id)]).then(
        ([stats, buffer]) => {
          return buffer.toString('binary');
        }
      );
    },
    transform(code, id) {
      if (code && /\.wasm$/.test(id)) {
        let src;

        src = gzip(Buffer.from(code, 'binary'))
        src = Buffer.from(src, 'binary').toString('base64');
        src = `'${src}'`;
        return {
          map: {
            mappings: ''
          },
          code: `
                  import { _loadWasmModule } from ${JSON.stringify(HELPERS_ID)};
                  import { ungzip } from 'pako';
                  export default function(imports){return _loadWasmModule(${src}, imports, ungzip)}`
        };
      }
      return null;
    },

  }
};
