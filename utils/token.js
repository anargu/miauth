const jwt = require('jsonwebtoken')
const miauthConfig = require('../config')

function expirationOffset (exp) {
    return Math.floor(Date.now() / 1000) + (parseInt(exp || process.env.JWT_EXP_OFFSET) || (5 * 60 * 60))
}

async function tokenize (exp = null, payload = {}) {
    const expires_in = expirationOffset(exp)
    
    const access_token = await jwt.sign({
        exp: expires_in,
        ...payload
    }, process.env.JWT_SECRET || 'm14uth')

    const refresh_token = await jwt.sign({
    }, process.env.REFRESH_SECRET || 'm14uth-refresh')

    return { access_token, refresh_token, expires_in }
}

function decodeToken (token) {
    return jwt.decode(token)
}

async function verify (token) {
    try {
        const payload = await jwt.verify(token, process.env.JWT_SECRET || 'm14uth')
        return { isOk: true, payload: { ...payload } }
    } catch (error) {
        return { isOk: false, error: error }
    }
}

module.exports = {
    tokenize,
    decodeToken,
    verify,
}
