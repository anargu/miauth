const YAML = require('yaml')
const fs = require('fs')

let miauthConfig = null

function initConfig() {
    const { 
        MIAUTH_CONFIG_FILE
    } = process.env
    const miauthConfigString = fs.readFileSync(MIAUTH_CONFIG_FILE, { encoding: 'utf-8' })    

    miauthConfig = YAML.parse(miauthConfigString)
    // miauthConfig['port']
    // bcrypt:
    //     salt: 
    // access_token:
    //     secret: 
    //     expires_in:
    // refresh_token:
    //     enabled: 
    //     secret: 
    // reset_password:
    //     expires_in: 
    //     secret: '

    return miauthConfig
}

module.exports = (() => {
    if (!miauthConfig) 
        return initConfig()
    return miauthConfig
})()

// fn just for testing purposes
module.exports.initConfig = initConfig