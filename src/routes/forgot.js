
const path = require('path')
const express = require('express')
const { check, oneOf, validationResult, query } = require('express-validator');
const miauthConfig = require('../config')
const { MiauthError } = require('../utils/error.js')
const { User, Session } = require('../models')
const { tokenize, expirationOffset, verify } = require('../utils/token')

const forgotRoute = express.Router()

// step 1: user request for reset password
forgotRoute.post('/request', oneOf([
    check('username', miauthConfig.field_validations.username.invalid_pattern_error_message)
        .exists({ checkNull: true })
        .matches(new RegExp(miauthConfig.field_validations.username.pattern, 'g'))
        .isLength({
            min: miauthConfig.field_validations.username.len[0],
            max: miauthConfig.field_validations.username.len[1]
        }).withMessage(`Invalid username. Username should be between\ 
        ${miauthConfig.field_validations.username.len[0]} and\ 
        ${miauthConfig.field_validations.username.len[1]} characteres`),
    check('email', miauthConfig.field_validations.email.invalid_pattern_error_message)
        .exists({ checkNull: true })
        .matches(new RegExp(miauthConfig.field_validations.email.pattern, 'g'))
        .isLength({
            min: miauthConfig.field_validations.email.len[0],
            max: miauthConfig.field_validations.email.len[1]
        }).withMessage(`Invalid email. Email should be between\ 
        ${miauthConfig.field_validations.email.len[0]} and\ 
        ${miauthConfig.field_validations.email.len[1]} characteres`)
]), async (req, res) => {
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
    query('token', 'Token is not setted')
        .exists({ checkNull: true }),
    check(['new_password'], miauthConfig.field_validations.password.invalid_pattern_error_message)
        .exists({ checkNull: true })
        .isString()
        .isLength({
            min: miauthConfig.field_validations.password.len[0],
            max: miauthConfig.field_validations.password.len[1]
        }),
    check('retyped_password', 'password value mismatch!')
        .isString()
        .isLength({
            min: miauthConfig.field_validations.password.len[0],
            max: miauthConfig.field_validations.password.len[1]
        }).withMessage(`repeat password input must be between\ 
        ${miauthConfig.field_validations.password.len[0]} \ 
        and ${miauthConfig.field_validations.password.len[1]} characters`)
        .custom(
            (retypedPassword, { req }) => (retypedPassword === req.body.new_password))
], (req, res, next) => {
    const errors = validationResult(req)
    try {
        if (!errors.isEmpty()) {
            // next(
            //     new MiauthError(
            //         400,
            //         'invalid_input_data',
            //         errors.array().map(e => e.msg).join('\n'))
            // )
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


module.exports = forgotRoute