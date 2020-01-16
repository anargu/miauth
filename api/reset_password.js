
const path = require('path')
const express = require('express')
const { checkSchema } = require('express-validator')
const fs = require('fs')

const resetPassApi = express.Router()

// step 1: user request for reset password
resetPassApi.post('/forgot', () => {
    
})

resetPassApi.get('/reset', (req, res) => {
    const htmlForm = fs.readFileSync(path.join(__dirname, '../public', 'reset_password.html'), { encoding: 'utf-8' })
    res.render(htmlForm)
})

resetPassApi.post('/reset', (req, res) => {
    
    req.query.token
    req.body.new_password
})
