const jwt = require('jsonwebtoken')

function errorMessage (userMessage = '', message = '', code = 100) {
    return {
        user_message: userMessage,
        message,
        code
    }
}

function expirationOffset () {
    return Math.floor(Date.now() / 1000) + (parseInt(process.env.JWT_EXP_OFFSET) || (5 * 60 * 60))
}

async function tokenize (payload = {}) {
    const token = await jwt.sign({
        exp: expirationOffset(),
        ...payload
    }, process.env.JWT_SECRET || 'm14uth')
    return token
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

function isEmpty(value) {
    if (value instanceof String) {
        const _value = value.trim()
        return _value === ''
    } else {
        return (value === undefined || value === null)
    }
}

module.exports = {
    errorMessage,
    tokenize,
    decodeToken,
    verify,
    isEmpty
}
