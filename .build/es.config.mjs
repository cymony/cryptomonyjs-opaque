import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import typescript from "@rollup/plugin-typescript";
import clear from 'rollup-plugin-clear';
import terser from '@rollup/plugin-terser';
import dts from 'rollup-plugin-dts';

import { files, isDev, getTsConfigPath, pkg, getAbsPath } from './helpers.mjs';
import { gzipInlineWasm, emitModulePackageFile } from '../plugin/index.mjs'

const typesConf = [
    {
        input: getAbsPath('../dist/temptypes/ts/index.d.ts'),
        output: [
            { file: getAbsPath('../dist/types/cryptomonyjs-opaque.d.ts'), format: 'es' }
        ],
        plugins: [
            clear({
                targets: [getAbsPath('../dist/types')],
            }),
            dts()
        ]
    }
]

const esConf = [
    {
        input: [
            ...files(getAbsPath('../src/ts'))
        ],
        output: [
            {
                exports: 'named',
                sourcemap: isDev(),
                format: 'es',
                file: pkg.module
            }
        ],
        plugins: [
            clear({
                targets: [getAbsPath('../dist/es')],
                watch: true
            }),
            gzipInlineWasm(),
            nodeResolve(),
            commonjs(),
            typescript({
                tsconfig: getTsConfigPath(),
                sourceMap: isDev()
            }),
            !isDev() && terser(
                {
                    format: {
                        comments: false
                    }
                }
            ),
            emitModulePackageFile()
        ]
    },
    ...(!isDev() ? typesConf : [])
]

export default esConf