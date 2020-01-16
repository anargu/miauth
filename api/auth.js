

const express = require('express')
const { checkSchema } = require('express-validator')

const miauthConfig = require('../config')

const { User, Session } = require('../models')
const { verifyPassword } = require('../utils/auth')
const { errorMessage } = require('../utils/utils')

const userSchemaValidation = (() => {
    const _userSchemaValidation = {}
    if (miauthConfig.user.username) {
        _userSchemaValidation['username'] = {
            matches: {
                options: [new RegExp(miauthConfig.field_validations.username, 'g')],
                errorMessage: 'Only use alphanumeric characters, \'-\' and \'_\'.'
            },
            nullable: false,
        }
    }
    
    if (miauthConfig.user.email) {
        _userSchemaValidation['email'] = {
            matches: {
                options: [new RegExp(miauthConfig.field_validations.email, 'g')],
                errorMessage: 'Please type a valid'
            },
            nullable: false,
        }
    }

    _userSchemaValidation['password'] = {
        nullable: false,
        isString: true,
        isLength: { min: miauthConfig.field_validations.password[0] }
    }

    return _userSchemaValidation
})()


const authApi = express.Router()

const authenticate = async (req, res) => {
    const { username, email, password } = req.body

    try {
        // TODO: I'm here
        if (username) {

        }
        const user = await User.findByUsername(username)
        if(await verifyPassword(password, user.hash)) {
            const session = await Session.createSession({ userId: user.uuid })
            res.status(200).json(session)
        } else {
            res.status(400).json(errorMessage('Invalid username or password', 'incorrect user credentials'))
        }        
    } catch (error) {
        res.status(400).json(errorMessage('Invalid username or password', error.message))        
    }
}

authApi.post('/login', checkSchema(userSchemaValidation), authenticate)

authApi.post('/signup', checkSchema(userSchemaValidation), async (req, res) => {

    const
    const _user = await User.createUser({
        username: req.body.username,
        email: req.body.email,
        password: req.body.password
    })

    res.status(200).json(_user)
})

authApi.post('/token/refresh', null)

module.exports = authApi