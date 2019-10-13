const saltRounds = parseInt(process.env.SALT) || 10
const bcrypt = require('bcryptjs')
const uuidv4 = require('uuid/v4')

const generateUUID = () => {
    return uuidv4()
}

const encodePassword = async (password) => {
    let encoded = await bcrypt.hash(password, saltRounds)
    return encoded
}

const verifyPassword = async (inputPassword, hashedPassword) => {
    let result = await bcrypt.compare(inputPassword, hashedPassword)
    return result
}

module.exports = {
    generateUUID,
    encodePassword,
    verifyPassword
}