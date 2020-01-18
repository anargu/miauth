
const path = require('path')
const assert = require('assert')

describe('testing utility methods for generating JWT', () => {

    before('setting init env variables', () => {
        process.env.MIAUTH_CONFIG_FILE = path.join(__dirname, '../miauth.config.yml')
        process.env.ACCESS_TOKEN_SECRET = 'a baby secret'
    })

    it('should generate a valid token without expiration', async () => {
        
        const miauthConfig = require(path.join(__dirname, '../src/config')).initConfig()
        const { tokenize, decodeToken, verify } = require(path.join(__dirname, '../src/utils/token.js'))
        
        const jwt = await tokenize({}, miauthConfig.ACCESS_TOKEN_SECRET)
        assert.deepEqual(typeof jwt, 'string')
        
        const decodedJwt = decodeToken(jwt)
        assert.deepEqual(typeof decodedJwt, 'object')

        const validToken = await verify(jwt, miauthConfig.ACCESS_TOKEN_SECRET)
        assert.deepEqual(validToken.isOk, true)
    })

    it('should generate a valid token with expiration', async () => {
        process.env.ACCESS_TOKEN_EXPIRATION = '600'
        const miauthConfig = require(path.join(__dirname, '../src/config')).initConfig()
        const { tokenize, decodeToken, verify } = require(path.join(__dirname, '../src/utils/token.js'))

        const jwt = await tokenize(
            {},
            miauthConfig.ACCESS_TOKEN_SECRET,
            miauthConfig.ACCESS_TOKEN_EXPIRATION)

        assert.deepEqual(typeof jwt, 'string')

        const validToken = await verify(jwt, miauthConfig.ACCESS_TOKEN_SECRET)
        assert.deepEqual(validToken.isOk, true)
    })

    const { promisify } = require('util')
    const sleep = promisify(setTimeout)
    it('should invalidate an expired token', async function () {
        this.timeout(1500)
        process.env.ACCESS_TOKEN_EXPIRATION = '1' // 5 seconds
        miauthConfig = require(path.join(__dirname, '../src/config')).initConfig()
        const { tokenize, verify } = require(path.join(__dirname, '../src/utils/token.js'))

        const jwt = await tokenize(
            {},
            miauthConfig.ACCESS_TOKEN_SECRET,
            miauthConfig.ACCESS_TOKEN_EXPIRATION)

        await sleep(1200)

        const invalidToken = await verify(jwt, miauthConfig.ACCESS_TOKEN_SECRET)
        console.log('invalidToken', invalidToken)
        assert.deepEqual(invalidToken.isOk, false)
    })
})