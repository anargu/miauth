const miauthConfig = require('../config')

const { check, query } = require('express-validator')

const validations = {
    // check username
    check_username: () => (
        check('username', miauthConfig.field_validations.username.invalid_pattern_error_message)
            .exists({ checkNull: true })
            .matches(new RegExp(miauthConfig.field_validations.username.pattern), 'g')
            .isLength({
                min: miauthConfig.field_validations.username.len[0],
                max: miauthConfig.field_validations.username.len[1]
            }).withMessage(`Invalid username. Username should be between\ 
            ${miauthConfig.field_validations.username.len[0]} and\ 
            ${miauthConfig.field_validations.username.len[1]} characteres`)
    ),
    // check email
    check_email: () => (
        check('email', miauthConfig.field_validations.email.invalid_pattern_error_message)
            .exists({ checkNull: true })
            .matches(new RegExp(miauthConfig.field_validations.email.pattern), 'g')
            .isLength({
                min: miauthConfig.field_validations.email.len[0],
                max: miauthConfig.field_validations.email.len[1]
            }).withMessage(`Invalid email. Email should be between\ 
            ${miauthConfig.field_validations.email.len[0]} and\ 
            ${miauthConfig.field_validations.email.len[1]} characteres`)
    ),
    // check token
    check_token: () => (
        query('token', 'Token is not setted')
            .exists({ checkNull: true })
    ),
    // check password
    check_password: (inputName = 'password') => (
        check([inputName], miauthConfig.field_validations.password.invalid_pattern_error_message)
            .exists({ checkNull: true })
            .isLength({
                min: miauthConfig.field_validations.password.len[0],
                max: miauthConfig.field_validations.password.len[1]
            })
    ),
    // check retyped_password
    check_retyped_password: (retyped_password_fieldname = 'retyped_password') => (
        check(retyped_password_fieldname, 'password value mismatch!')
            .isString()
            .isLength({
                min: miauthConfig.field_validations.password.len[0],
                max: miauthConfig.field_validations.password.len[1]
            }).withMessage(`repeat password input must be between\ 
            ${miauthConfig.field_validations.password.len[0]} \ 
            and ${miauthConfig.field_validations.password.len[1]} characters`)
            .custom(
                (retypedPassword, { req }) => (retypedPassword === req.body.new_password))
    ),
    // grant_type,
    check_grant_type: () => (
        check('grant_type').notEmpty()
        .custom((value) => {
            if (value !== 'refresh_token') {
                throw new Error('invalid grant type')
            }
            return true
        })
        .withMessage('Invalid or empty grant_type.')
    ),
    // refresh_token
    check_refresh_token: () => (
        check('refresh_token').notEmpty().withMessage('Empty refresh token')
    ),
    // scope
    check_scope: () => (
        check('scope').not().notEmpty()
    ),
    check_userId: () => (
        check('userId')
        .exists()
        .isString()
    )
}

module.exports = validations