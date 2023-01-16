import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import typescript from "@rollup/plugin-typescript";
import clear from 'rollup-plugin-clear';
import terser from '@rollup/plugin-terser';

import { files, isDev, getTsConfigPath, pkg, getAbsPath } from './helpers.mjs';
import { gzipInlineWasm } from '../plugin/index.mjs'

export default {
    input: [
        ...files(getAbsPath('../src/ts'))
    ],
    output: [
        {
            exports: 'named',
            sourcemap: isDev(),
            name: "Opaque",
            format: 'iife',
            file: pkg.browser
        }
    ],
    plugins: [
        clear({
            targets: [getAbsPath('../dist/js')],
            watch: true
        }),
        gzipInlineWasm(),
        nodeResolve(),
        commonjs(),
        typescript({
            tsconfig: getTsConfigPath(),
            sourceMap: isDev(),
        }),
        !isDev() && terser(
            {
                format: {
                    comments: false
                }
            }
        ),
    ]
}