

const express = require('express')
const { checkSchema, oneOf, validationResult } = require('express-validator')

const miauthConfig = require('../config')

const { verifyPassword } = require('../utils/auth')
const { errorMessage } = require('../utils/misc')
const {
    check_username,
    check_email,
    check_password,
    check_grant_type,
    check_refresh_token,
    check_scope
} = require('../middlewares/validations')

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

    authApi.post('/login', oneOf([
        // case 1 username
        ... (miauthConfig.user.username) ? [check_username(), check_password()] : [],
        // case 2 email
        ... (miauthConfig.user.email) ? [check_email(), check_password()] : []
    ]), authenticate)
    
    authApi.post('/signup', [
        ... (miauthConfig.user.username) ? [check_username()] : [],
        ... (miauthConfig.user.check_email) ? [check_email()] : [],
        check_password(),
    ], async (req, res) => {
    
        const _user = await User.createUser({
            username: req.body.username,
            email: req.body.email,
            password: req.body.password
        })
    
        res.status(200).json(_user)
    })
    
    if(miauthConfig.refresh) {
        authApi.post('/token/refresh', [
            check_grant_type(),
            check_refresh_token(),
            check_scope()
        ], async (req, res) => {
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