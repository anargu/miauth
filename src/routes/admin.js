const express = require('express')
const { checkSchema } = require('express-validator')

const { Session } = require('../models')


const adminApi = express.Router()

adminApi.post('/revoke_all', checkSchema({
    userId: {
        isString: true,
        notEmpty: true,
    }
}), async (req, res) => {
    
    const sessionsDeleted = await Session.revokeAll({ userId: req.body.userId })

    res.status(200).json({ sessions_deleted: sessionsDeleted })    
})


module.exports = adminApi