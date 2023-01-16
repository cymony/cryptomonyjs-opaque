const gulp = require('gulp');
const rimraf = require('rimraf');
const { series, watch, task, parallel } = gulp;

const browsersync = require('browser-sync');
const server = browsersync.create();

const shell = require('gulp-shell');

// Development server tasks
task("dev:serve", () => {
    server.init({
        files: [
            "dist",
            "assets"
        ],
        watchEvents: ["add", "change", "addDir"],
        server: ["assets", "dist"],
        port: 8000,
        browser: "google chrome"
    })
});

// Typescript tasks
task('ts:watch', shell.task('npx rollup --config rollup.config.mjs --config-dev --watch'));
task('ts:type:compile', shell.task("npx tsc -d --allowJS --emitDeclarationOnly --declarationDir dist/temptypes"))
task('ts:type:clean', (cb) => {
    rimraf('dist/temptypes', cb)
})
task('ts:compile', shell.task('npx rollup --config rollup.config.mjs'))

// Go tasks
task('go:clean', (cb) => {
    rimraf('./src/api/lib.wasm', cb);
})
task('go:compile', shell.task('cd src/api/ && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o lib.wasm $(find . -name "*.go" -type f -print0 | xargs -0)'))
task('go:watch', () => [
    watch([
        'src/api/**/*.go'
    ], series('go:compile'))
])

// Main tasks
task('dev', parallel(['go:watch', 'ts:watch', 'dev:serve']))
task('build', series(['go:clean', 'go:compile', 'ts:type:compile', 'ts:compile', 'ts:type:clean']))