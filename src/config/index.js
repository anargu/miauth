const YAML = require('yaml')
const fs = require('fs')

let miauthConfig = null

function initConfig() {
    let miauthConfigString
    try {
        const { 
            MIAUTH_CONFIG_FILE
        } = process.env
        miauthConfigString = fs.readFileSync(MIAUTH_CONFIG_FILE, { encoding: 'utf-8' })                
    } catch (error) {
        throw new Error('MIAUTH_CONFIG_FILE variable given incorrectly or just not provided')
    }

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