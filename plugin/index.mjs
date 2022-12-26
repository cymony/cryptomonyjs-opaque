import { HELPERS_ID, getHelperModule } from './helpers.mjs';
import * as fs from 'fs';
import { gzip } from 'pako';

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