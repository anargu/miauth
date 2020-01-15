const saltRounds = parseInt(process.env.SALT) || 10
const bcrypt = require('bcryptjs')

const encodePassword = async (password) => {
    let encoded = await bcrypt.hash(password, saltRounds)
    return encoded
}

const verifyPassword = async (inputPassword, hashedPassword) => {
    let result = await bcrypt.compare(inputPassword, hashedPassword)
    return result
}

module.exports = {
    encodePassword,
    verifyPassword
}