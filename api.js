const express = require('express')
const User = require('./models').User
const utils = require('./utils/utils')
const { encodePassword, verifyPassword } = require('./utils/auth')
const { check, validationResult } = require('express-validator');

function settingUpEndpoints (app) {
    const router = express.Router()

    // request inputs { email, password }
    router.post('/signup', [
        
    ], async (req, res) => {
        const email = req.body.email
        const username = req.body.username
        if (process.env.LOGIN_BY === LOGIN_BY_USERNAME && utils.isEmpty(username)) {
            res.status(400).json(utils.errorMessage(
                'Please provide an username.',
                'username field is needed.',
                702
            ))
            return                
        }
        const password = req.body.password
        if(utils.isEmpty(email) || utils.isEmpty(password)) {
            res.status(400).json(utils.errorMessage(
                'empty email or password',
                'empty input data',
                103
            ))
            return
        }
        try {
            const hash = await encodePassword(password)

            const user = User.init(email, username, hash)
            await user.save()

            res.status(200).json({...user})
        } catch (error) {
            res.status(500).json(utils.errorMessage(
                'Something went wrong. Please try again in a moment.',
                error.toString(),
                701
            ))
        }
    })

    // request inputs { email or username, password }
    router.post('/login', async (req, res) => {
        const userIdentifier = req.body.username || req.body.email

        const password = req.body.password
        if(utils.isEmpty(userIdentifier) || utils.isEmpty(password)) {
            res.status(400).json(utils.errorMessage(
                'empty email or password',
                'empty input data',
                103
            ))
            return
        }
        try {
            const user = await User.find(userIdentifier)
            // if user not found
            if (!user) {
                res.status(400).json(utils.errorMessage(
                    'email/username not found, are you typing your email/username correctly?',
                    'email not found',
                    101
                ))
                return
            }
            // if user exists then validate password
            const result = await verifyPassword(password, user.hash)
            if (result === true) {
                let payload = {}
                const payloadKey = `user_${process.env.LOGIN_BY}`
                payload[payloadKey] = userIdentifier
                payload['user_email'] = user.email
                const token = await utils.tokenize(payload)
                res.status(200).json({
                    status: 'ok',
                    token: token
                })
            } else {
                res.status(400).json(utils.errorMessage(
                    'Invalid password. Please remember it and try again',
                    'invalid password',
                    102
                ))
            }
        } catch (error) {
            res.status(500).json(utils.errorMessage(
                'Something went wrong. Please try again in a moment.',
                'fail on verifying password ' + error.toString(),
                700
            ))
        }
    })

    router.post('/verify', async (req, res) => {
        const token = req.body.token
        const result = await utils.verify(token)
        if (!result.isOk) {
            res.status(400).json(utils.errorMessage(
                'Invalid token',
                'invalid token',
                103
            ))
        } else {
            res.status(200).json({
                status: 'ok',
                payload: result.payload
            })
        }
    })

    const tokenRouter = express.Router()
    tokenRouter.post('/refresh', (req, res) => {})
    tokenRouter.post('/revoke', (req, res) => {})
    router.use('/token', tokenRouter)

    app.use('/auth', router)
}

module.exports = {
    settingUpEndpoints
}
