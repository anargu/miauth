const bcrypt = require('bcryptjs')
const miauthConfig = require('../config')
const saltRounds = parseInt(miauthConfig.SALT) || 10

const hashPassword = async (password) => {
    let encoded = await bcrypt.hash(password, saltRounds)
    return encoded
}

const verifyPassword = async (inputPassword, hashedPassword) => {
    let result = await bcrypt.compare(inputPassword, hashedPassword)
    return result
}

module.exports = {
    hashPassword,
    verifyPassword
}