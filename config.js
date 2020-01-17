const YAML = require('yaml')
const fs = require('fs')

let miauthConfig = null

function initConfig() {
    const { MIAUTH_CONFIG_FILE } = process.env
    const miauthConfigString = fs.readFileSync(MIAUTH_CONFIG_FILE, { encoding: 'utf-8' })
    miauthConfig = YAML.parse(miauthConfigString)

    return miauthConfig
}

module.exports = (() => {
    if (!miauthConfig) 
        return initConfig()
    return miauthConfig
})()

module.exports.initConfig = initConfig