

const express = require('express')
const { checkSchema, oneOf, validationResult } = require('express-validator')

const miauthConfig = require('../config')

const { verifyPassword } = require('../utils/auth')
const { verify } = require('../utils/token')
const { MiauthError } = require('../utils/error')
const {
    check_username,
    check_email,
    check_password,
    check_grant_type,
    check_refresh_token,
    check_access_token
} = require('../middlewares/validations')

module.exports = (db) => {
    const { User, Session } = db

    const authApi = express.Router()

    const authenticate = async (req, res, next) => {    
        try {
            const errors = validationResult(req)
            if(!errors.isEmpty()) {
                throw new MiauthError(
                    400, 'ValidationError',
                    errors.array().map((err) => err.msg).join('. '))
            }

            const { username, email, password } = req.body

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
                throw new MiauthError(
                    400, 'authentication_failed',
                    'Invalid username or password')
            }        
        } catch (error) {
            next(error)
        }
    }

    authApi.post('/login', [
        check_password(),
        oneOf([
            ... (miauthConfig.user.username) ? [check_username()] : [],
            ... (miauthConfig.user.email) ? [check_email()] : [],
        ])
    ], authenticate)
    
    authApi.post('/signup', [
        ... (miauthConfig.user.username) ? [check_username()] : [],
        ... (miauthConfig.user.email) ? [check_email()] : [],
        check_password(),
    ], async (req, res, next) => {
        try {
            const errors = validationResult(req)
            if(!errors.isEmpty()) {
                throw new MiauthError(
                    400, 'ValidationError',
                    errors.array().map((err) => err.msg).join('. '))
            }

            const _user = await User.createUser({
                username: req.body.username,
                email: req.body.email,
                password: req.body.password
            })
        
            res.status(200).json(_user)                
        } catch (error) {
            next(error)
        }
    })

    authApi.get('/verify', [
        check_access_token()
    ], async (req, res, next) => {
        try {
            const errors = validationResult(req)
            if(!errors.isEmpty()) {
                throw new MiauthError(
                    400, 'ValidationError',
                    errors.array().map((err) => err.msg).join('. '))
            }

            const tokenValid = await verify(req.query.access_token, miauthConfig.access_token.secret)
            if (tokenValid.isOk) {
                res.status(200).json({...tokenValid})
            } else {
                throw new MiauthError(400,
                    tokenValid.error.name,
                    tokenValid.error.message,
                    'Invalid token')
            }
        } catch (error) {
            next(error)
        }
    })
    
    if(miauthConfig.refresh_token.enabled) {
        authApi.post('/token/refresh', [
            check_grant_type(),
            check_refresh_token(),
            // disabled for a while
            // check_scope()
        ], async (req, res, next) => {
            try {
                const errors = validationResult(req)
                if(!errors.isEmpty()) {
                    throw new MiauthError(
                        400, 'ValidationError',
                        errors.array().map((err) => err.msg).join('. '))
                }

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
        
                res.status(200).json({
                    old_refresh_token: session.refresh_token,
                    userId: newSession.userId,
                    createdAt: newSession.createdAt,
                    access_token: newSession.access_token,
                    refresh_token: newSession.refresh_token,
                    expires_in:  newSession.expires_in
                })
            } catch (error) {
                next(error)
            }
        })    
    }
 
    return authApi
}