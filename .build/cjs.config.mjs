import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import typescript from "@rollup/plugin-typescript";
import clear from 'rollup-plugin-clear';
import terser from '@rollup/plugin-terser';
import emitFiles from 'rollup-plugin-emit-files';

import { gzipInlineWasm } from '../plugin/index.mjs'
import { files, isDev, getAbsPath, pkg, getTsConfigPath } from "./helpers.mjs"

export default {
    input: [
        ...files(getAbsPath('../src/ts'))
    ],
    output: [
        {
            format: 'cjs',
            file: pkg.main,
            intro: "require('./wasm_exec_node.js');"
        }
    ],
    plugins: [
        clear({
            targets: [getAbsPath('../dist/cjs')],
            watch: true
        }),
        gzipInlineWasm(),
        nodeResolve(),
        commonjs(),
        typescript({
            tsconfig: getTsConfigPath(),
            sourceMap: isDev(),
        }),
        emitFiles({ src: getAbsPath('../polyfils') }),
        !isDev() && terser(
            {
                format: {
                    comments: false
                }
            }
        ),
    ]
}