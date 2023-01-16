import esConf from "./.build/es.config.mjs";
import iifeConf from './.build/iife.config.mjs';
import cjsConf from './.build/cjs.config.mjs';

export default [
    ...esConf,
    { ...iifeConf },
    { ...cjsConf }
]