const jwt = require('jsonwebtoken')
const miauthConfig = require('../config')

const defaultExpirationOffsetTime = (5 * 60 * 60)

function expirationOffset (exp) {
    return Math.floor(Date.now() / 1000) + (parseInt(exp) || defaultExpirationOffsetTime)
}

async function tokenize (payload = {}, secret, exp = undefined) {
    const expires_in = (exp === undefined) ? undefined : expirationOffset(exp)

    const JWTSignPayload = {...payload}
    if(expires_in)
        JWTSignPayload['exp'] = expires_in
    
    const token = await jwt.sign({...JWTSignPayload}, secret)

    return token
}

function introspect (token) {
    return jwt.decode(token)
}

async function verify (token, secret) {
    try {
        const payload = await jwt.verify(token, secret)
        return { isOk: true, payload: { ...payload } }
    } catch (error) {
        return { isOk: false, error: error }
    }
}

module.exports = {
    expirationOffset,
    tokenize,
    introspect,
    verify,
}
