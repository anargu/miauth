
const path = require('path')
const express = require('express')
const miauthConfig = require('../config')
const { MiauthError } = require('../utils/error.js')
const { tokenize, expirationOffset, verify } = require('../utils/token')
const { check, oneOf, validationResult, query } = require('express-validator')
const { check_retyped_password, check_token, check_username, check_email, check_password } = require('../middlewares/validations')

module.exports = (db) => {
    const { User, Session } = db

    const forgotRoute = express.Router()

    // step 1: user request for reset password
    forgotRoute.post('/request', oneOf([
        ... (miauthConfig.user.username) ? [check_username()] : [],
        ... (miauthConfig.user.email) ? [check_email()] : []
    ]), async (req, res, next) => {
        try {
            validationResult(req).throw()

            const fieldValue = req.body.username || req.body.email
            const findByUsername = req.body.username === undefined ? false : true
            
            let _userFound
            if (findByUsername) {
                _userFound = await User.findByUsername(fieldValue)
            } else {
                _userFound = await User.findByEmail(fieldValue)
            }
            if (!_userFound)
                throw new MiauthError(400, 'user_not_found', 'Usuario no encontrado/registrado.')
            // user found, generate token
            const resetEmailToken = await tokenize(
                { userId: _userFound.uuid, email: _userFound.email },
                miauthConfig.reset_password.secret,
                expirationOffset(miauthConfig.reset_password.token_expiration),
                
            )
            // send email to user
            // TODO: Create Email Service
            next(new MiauthError(500, 'service_not_implemented', 'service_not_implemented'))
            return

            // respond ok
            res.status(200).json({
                message: 'Email sent to the user with the instructure to restore password.'
            })
        } catch (err) {
            throw err
        }
    })

    forgotRoute.get('/reset', (req, res) => {
        res.render(
            path.join(__dirname, '../../public', 'reset_password.html'),
            { title: 'Reset Password' })
    })

    forgotRoute.post('/reset', [
        check_token(),
        check_password('new_password'),
        check_retyped_password(),
    ], async (req, res, next) => {
        const errors = validationResult(req)
        try {
            if (!errors.isEmpty()) {
                return res.render(
                    path.join(__dirname, '../../public', 'reset_password_result_error.html'),
                    { errors: errors.array() })
            } else {
                // variables to be used
                // req.query.token
                // req.body.new_password
                const tokenValid = await verify(req.query.token)
                if(!tokenValid.isOk) {
                    next(
                        new MiauthError(
                            400,
                            'invalid_token',
                            tokenValid.error.name,
                            'Reset Password session has expired. Please request again to recover password')
                    )
                }

                // REMOVING ALL SESSIONS LOGGED IN OF USER
                await Session.revokeAllSessions({ userId: tokenValid.payload.userId })


                // UPDATING PASSWORD
                // const userUpdated = 
                await User.updatePassword({
                    field: 'userId',
                    value: tokenValid.payload.userId
                }, req.body.new_password)

                return res.render(
                    path.join(__dirname, '../../public', 'reset_password_result_success.html'))
            }            
        } catch (error) {
            next(error)
        }
    })

    return forgotRoute
}