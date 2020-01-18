
const path = require('path')
const express = require('express')
const { check, oneOf, validationResult, query } = require('express-validator');
const miauthConfig = require('../config')
const { MiauthError } = require('../utils/error.js')

const forgotRoute = express.Router()

// step 1: user request for reset password
forgotRoute.post('/request', [
    check('username', miauthConfig.field_validations.username.invalid_pattern_error_message)
    .exists({ checkNull: true })
    .matches(new RegExp(miauthConfig.field_validations.username.pattern, 'g'))
    .isLength({
        min: miauthConfig.field_validations.username.len[0],
        max: miauthConfig.field_validations.username.len[1]
    }).withMessage(`Invalid username. Username should be between\ 
    ${miauthConfig.field_validations.username.len[0]} and\ 
    ${miauthConfig.field_validations.username.len[1]} characteres`)
], () => {

    // 
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
            // req.query.token
            // req.query.new_password



            return res.render(
                path.join(__dirname, '../../public', 'reset_password_result_success.html'))
        }            
    } catch (error) {
        next(error)
    }
})


module.exports = forgotRoute