const express = require('express')
const bodyParser = require('body-parser')

const miauthConfig = require('./config')

const forgotRoute = require('./routes/forgot.js')
const adminRoute = require('./routes/admin.js')
const userRoute = require('./routes/user.js')
const { handleError } = require('./utils/error')

const server = express()

server.set('view engine', 'html');
server.engine('html', require('hbs').__express);

server.use(bodyParser.json())
server.use(bodyParser.urlencoded({ extended: true }))

server.use('/forgot', forgotRoute)
server.use('/admin', adminRoute)
server.use('/auth', userRoute)

server.use(function(err, req, res, next) {
    if (err)
        handleError(err, res)
});

server.listen(miauthConfig.PORT, () => {
    console.log(
        `▀▄▀▄▀▄ [ MiAuth started & listening on port ${miauthConfig.PORT} ] ▄▀▄▀▄▀`
    )
})
