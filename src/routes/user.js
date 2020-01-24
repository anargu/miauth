

const express = require('express')
const { checkSchema, oneOf, validationResult } = require('express-validator')

const miauthConfig = require('../config')

const { verifyPassword } = require('../utils/auth')
const { MiauthError } = require('../utils/error')
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
    
    if(miauthConfig.refresh_token.enabled) {
        authApi.post('/token/refresh', [
            check_grant_type(),
            check_refresh_token(),
            check_scope()
        ], async (req, res) => {
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
        
                res.status(200).json(newSession)
            } catch (error) {
                next(err)
            }
        })    
    }
 
    return authApi
}