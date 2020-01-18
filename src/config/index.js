const YAML = require('yaml')
const fs = require('fs')

let miauthConfig = null

function initConfig() {
    const { 
        MIAUTH_CONFIG_FILE,
        PORT,
        SALT,
        JWT_EXP_OFFSET,
        JWT_SECRET,
        REFRESH_SECRET
    } = process.env
    const miauthConfigString = fs.readFileSync(MIAUTH_CONFIG_FILE, { encoding: 'utf-8' })    

    miauthConfig = YAML.parse(miauthConfigString)
    miauthConfig['PORT'] = PORT
    
    miauthConfig['bcrypt'] = {}
    miauthConfig.bcrypt.SALT = SALT
    
    miauthConfig['JWT_EXP_OFFSET'] = JWT_EXP_OFFSET
    miauthConfig['JWT_SECRET'] = JWT_SECRET
    miauthConfig['REFRESH_SECRET'] = REFRESH_SECRET

    return miauthConfig
}

module.exports = (() => {
    if (!miauthConfig) 
        return initConfig()
    return miauthConfig
})()

module.exports.initConfig = initConfig