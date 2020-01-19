const path = require('path')

module.exports = {
    miauthSetup: () => {
        process.env.MIAUTH_CONFIG_FILE = path.join(__dirname, '../test.config.yml')
        process.env.PORT = '8000'
        process.env.ACCESS_TOKEN_SECRET = 'a baby secret'
        process.env.ACCESS_TOKEN_EXPIRATION = (2 * 60).toString()
        process.env.REFRESH_SECRET = 'a baby secret 2'

        require(path.join(__dirname, '../../src/config')).initConfig()
    }
}