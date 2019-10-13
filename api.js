const express = require('express')
const User = require('./models').User
const utils = require('./utils/utils')
const { encodePassword, verifyPassword } = require('./utils/auth_utils')


function settingUpEndpoints (app) {
    const router = express.Router()

    // request inputs { email, password }
    router.post('/signup', async (req, res) => {
        const email = req.body.email
        const password = req.body.password
        try {
            const hash = await encodePassword(password)

            const user = User.init(email, hash)
            await user.save()

            res.status(200).json({
                id: user.id,
                email: user.email,
                hash: user.hash
            })
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
        const email = req.body.email
        // const username = req.body.username
        const password = req.body.password
        if (email === '' || email === undefined ||
        password === '' || password === undefined) {
            res.status(400).json(utils.errorMessage(
                'empty email or password',
                'empty input data',
                103
            ))
            return
        }

        const user = await User.find(email)
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
        try {
            const result = await verifyPassword(password, user.hash)
            if (result === true) {
                const token = await utils.tokenize({ user_email: email })
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

    router.post('/refreh-token', (req, res) => {
    })

    app.use('/auth', router)
}

module.exports = {
    settingUpEndpoints
}
