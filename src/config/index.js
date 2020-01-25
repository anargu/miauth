const YAML = require('yaml')
const fs = require('fs')

let miauthConfig = null

function initConfig() {
    let miauthConfigString
    const { 
        MIAUTH_CONFIG
    } = process.env
    miauthConfigString = MIAUTH_CONFIG

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