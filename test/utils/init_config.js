
module.exports = {
    miauthSetup: () => {
        process.env.MIAUTH_CONFIG_FILE = path.join(__dirname, '../test.config.yml')
        process.env.ACCESS_TOKEN_SECRET = 'a baby secret'
        process.env.ACCESS_TOKEN_EXPIRATION = (2 * 60).toString()
        process.env.REFRESH_SECRET = 'a baby secret 2'
    }
}