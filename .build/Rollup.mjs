import { readFileSync } from "fs";
import path from "path";
import { fileURLToPath } from "node:url";

import { nodeResolve } from '@rollup/plugin-node-resolve';
import typescript from "@rollup/plugin-typescript";
import terser from '@rollup/plugin-terser';
import clear from 'rollup-plugin-clear';
import commonjs from '@rollup/plugin-commonjs';
import dts from 'rollup-plugin-dts';

import { files, isDev } from './helpers.mjs';
import { gzipInlineWasm, emitModulePackageFile, polyfiller } from '../plugin/index.mjs'


export const getConfig = () => ({
    input: [
        ...files('src/ts')
    ],
    output: [
        {
            exports: 'named',
            sourcemap: isDev(),
            format: 'es',
            file: pkg.module
        },
        {
            exports: 'named',
            sourcemap: isDev(),
            name: "Opaque",
            format: 'iife',
            file: pkg.browser
        }
    ],
    plugins: [
        gzipInlineWasm(),
        nodeResolve(),
        commonjs(),
        typescript({
            tsconfig: "./tsconfig.json",
            sourceMap: isDev(),
        }),
        // !isDev() && terser(
        //     {
        //         format: {
        //             comments: false
        //         }
        //     }
        // )
        clear({
            targets: ['dist/es'],
            watch: true
        }),
        emitModulePackageFile()

    ]
})

const getOutput = (format) => {
    const output = {
        exports: 'named',
        sourcemap: isDev()
    }

    switch (format) {
        case 'cjs':
            return {
                ...output,
                format: 'cjs',
                file: pkg.main,
                footer: 'module.exports = Object.assign(exports.default, exports);',
            }
        case 'es':
            return {
                ...output,
                format: 'es',
                file: pkg.module,
            }
        case 'iife':
            return {
                ...output,
                name: "Opaque",
                format: 'iife',
                file: pkg.browser,
            }
        default:
            throw Error('unrecognized format')
    }
}

const getPlugins = (format) => {
    const plugs = [
        gzipInlineWasm(),
        nodeResolve(),
        commonjs(),
        typescript({
            tsconfig: "./tsconfig.json",
            sourceMap: isDev(),
        }),
        // !isDev() && terser(
        //     {
        //         format: {
        //             comments: false
        //         }
        //     }
        // )
    ]

    switch (format) {
        case 'cjs':
            return [
                ...plugs,
                clear({
                    targets: ['dist/cjs'],
                    watch: true
                }),
                polyfiller([fileURLToPath(new URL('../polyfils/wasm_exec_node.js', import.meta.url))])
            ]
        case 'es':
            return [
                ...plugs,
                clear({
                    targets: ['dist/es'],
                    watch: true
                }),
                emitModulePackageFile()
            ]
        case 'iife':
            return [
                ...plugs,
                clear({
                    targets: ['dist/js'],
                    watch: true
                }),
            ]
        default:
            throw Error('unrecognized format')
    }
}

export const configGen = (formats) => {

    let genConf = formats.map((format) => {
        let obj = {
            input: [
                ...files("src/ts")
            ],
            output: [
                { ...getOutput(format) }
            ],
            plugins: [
                ...getPlugins(format)
            ]
        }

        if (format === "cjs") {
            console.log(fileURLToPath(new URL('../polyfils/wasm_exec_node.js', import.meta.url)))
            return { ...obj, external: [fileURLToPath(new URL('../polyfils/wasm_exec_node.js', import.meta.url))] }
        }
        return { ...obj }
    })

    let returnObj = [
        ...genConf
    ]

    if (!isDev()) {
        returnObj = [
            ...returnObj,
            {
                input: 'dist/temptypes/ts/index.d.ts',
                output: [
                    { file: 'dist/types/cryptomonyjs-opaque.d.ts', format: 'es' }
                ],
                plugins: [
                    clear({
                        targets: ['dist/types'],
                    }),
                    dts()
                ]
            }
        ]
    }

    return returnObj
} 