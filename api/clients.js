const express = require('express')
const { checkSchema } = require('express-validator')

const { Session } = require('../models')

const miauthConfig = require('../config')


const adminApi = express.Router()

adminApi.post('/revokeAll', checkSchema({
    userId: {
        isString: true,
        nullable: false,
    }
}), async (req, res) => {
    
    const sessionsDeleted = await Session.revokeAll({ userId: req.body.userId })

    res.status(200).json({ sessions_deleted: sessionsDeleted })    
})
