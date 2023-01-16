import { readdirSync, readFileSync } from "fs";
import path from "path";
import { fileURLToPath } from "url";

export const getAbsPath = (relative) => {
    return fileURLToPath(new URL(relative, import.meta.url))
}

export const isDev = () => {
    return !!process.argv.find(el => el === '--config-dev')
}

export const files = dir => {
    return readdirSync(dir).filter(el => path.extname(el) === '.ts').map(el => dir + "/" + el);
}

export const pkg = JSON.parse(readFileSync(getAbsPath('../package.json'), 'utf8'));

export const getTsConfigPath = () => {
    return getAbsPath('../tsconfig.json');
}

