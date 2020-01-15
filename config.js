const YAML = require('yaml')
const fs = require('fs')

let miauthConfig = null

async function initConfig() {
    const { MIAUTH_CONFIG_FILE } = process.env
    const miauthConfigString = await fs.readFileSync(MIAUTH_CONFIG_FILE, { encoding: 'utf-8' })
    miauthConfig = YAML.parse(miauthConfigString)

    return miauthconfig
}

module.exports = (() => {
    if (!miauthconfig) 
        return initConfig()
    return miauthConfig
})()

module.exports.initConfig = initConfig