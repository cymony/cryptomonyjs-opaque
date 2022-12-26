import typescript from "@rollup/plugin-typescript";
import terser from '@rollup/plugin-terser';
import clear from 'rollup-plugin-clear';
import { isDev, files } from "./.build/Rollup.mjs";
import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import { gzipInlineWasm } from './plugin/index.mjs'

const targetDir = "dist"



export default [
    {
        input: [
            ...files("src/ts")
        ],
        output: [
            {
                dir: targetDir,
                format: "es",
                sourcemap: isDev()
            }
        ],
        plugins: [
            gzipInlineWasm(),
            nodeResolve(),
            commonjs(),
            clear({
                targets: [targetDir],
                watch: true
            }),
            typescript({
                tsconfig: "./tsconfig.json"
            }),
            !isDev() && terser(
                {
                    format: {
                        comments: false
                    }
                }
            )
        ]
    },
]

