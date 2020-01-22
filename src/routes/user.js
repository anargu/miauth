

const express = require('express')
const { checkSchema } = require('express-validator')

const miauthConfig = require('../config')

const { verifyPassword } = require('../utils/auth')
const { errorMessage } = require('../utils/misc')

const userSchemaValidation = (() => {
    const _userSchemaValidation = {}
    if (miauthConfig.user.username) {
        _userSchemaValidation['username'] = {
            matches: {
                options: [new RegExp(miauthConfig.field_validations.username.pattern, 'g')],
            },
            isLength: {
                errorMessage: `Invalid username. Username should be between\ 
                ${miauthConfig.field_validations.username.len[0]} and\ 
                ${miauthConfig.field_validations.username.len[1]} characteres`,
                options: {
                    min: miauthConfig.field_validations.username.len[0],
                    max: miauthConfig.field_validations.username.len[1],
                }
            },
            notEmpty: true,
            errorMessage: miauthConfig.field_validations.username.invalid_pattern_error_message,
        }
    }
    
    if (miauthConfig.user.email) {
        _userSchemaValidation['email'] = {
            matches: {
                options: [new RegExp(miauthConfig.field_validations.email.pattern, 'g')],
            },
            isLength: {
                errorMessage: `Invalid email. Email should be between\ 
                ${miauthConfig.field_validations.email.len[0]} and\ 
                ${miauthConfig.field_validations.email.len[1]} characteres`,
                options: {
                    min: miauthConfig.field_validations.email.len[0],
                    max: miauthConfig.field_validations.email.len[1],
                }
            },
            errorMessage: miauthConfig.field_validations.email.invalid_pattern_error_message,
            notEmpty: true,
        }
    }

    _userSchemaValidation['password'] = {
        notEmpty: true,
        isString: true,
        isLength: {
            errorMessage: `Invalid password. Password should be between\ 
            ${miauthConfig.field_validations.password.len[0]} and\ 
            ${miauthConfig.field_validations.password.len[1]} characteres`,
            options: {
                min: miauthConfig.field_validations.password.len[0],
                max: miauthConfig.field_validations.password.len[1]    
            }
        },
        errorMessage: miauthConfig.field_validations.password.invalid_pattern_error_message,
    }

    return _userSchemaValidation
})()

module.exports = (db) => {
    const { User, Session } = db

    const authApi = express.Router()

    const authenticate = async (req, res) => {
        const { username, email, password } = req.body
    
        try {
            let user
            if (username) {
                user = await User.findByUsername(username)
            } else {
                user = await User.findByEmail(email)
            }
            if(await verifyPassword(password, user.hash)) {
                const session = await Session.createSession({ userId: user.uuid, email: user.email })
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
    
        const _user = await User.createUser({
            username: req.body.username,
            email: req.body.email,
            password: req.body.password
        })
    
        res.status(200).json(_user)
    })
    
    if(miauthConfig.refresh) {
        authApi.post('/token/refresh', checkSchema({
            grant_type: {
                notEmpty: true,
                custom: {
                    shouldBeRefreshToken: (value) => {
                        if (value !== 'refresh_token') {
                            throw new Error('invalid grant type')
                        }
                        return true
                    }
                },
                errorMessage: 'Invalid or empty grant_type.'
            },
            refresh_token: {
                isString: true,
                notEmpty: true,
                errorMessage: 'Empty refresh token',
            },
            scope: {
                notEmpty: false
            }    
        }), async (req, res) => {
            try {
                const session = await Session.findOne({
                    where: {
                        refresh_token: req.body.refresh_token
                    }
                })
                if (session === null) {
                    throw new Error('session not found')
                }
                const _user = await User.findOne({
                    where: {
                        uuid: session.userId
                    }
                })
                if (_user === null) {
                    throw new Error('User not found. Critical error.')
                }
        
                const newSession = await Session.createSession({
                    userId: session.userId,
                    email: _user.email
                })
                await session.destroy()
        
                res.status(200).json(newSession)
            } catch (error) {
                next(err)
            }
        })    
    }
 
    return authApi
}