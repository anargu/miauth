const jwt = require('jsonwebtoken')
const miauthConfig = require('../config')

function expirationOffset () {
    return Math.floor(Date.now() / 1000) + (parseInt(process.env.JWT_EXP_OFFSET) || (5 * 60 * 60))
}

async function tokenize (payload = {}) {
    const expires_in = expirationOffset()
    const token = await jwt.sign({
        exp: expireAt,
        ...payload
    }, process.env.JWT_SECRET || 'm14uth')

    const refresh_token = await jwt.sign({
        exp: expireAt,
    }, process.env.REFRESH_SECRET || 'm14uth-refresh')

    return { token, refresh_token, expires_in }
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
