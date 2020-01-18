const express = require('express')
const bodyParser = require('body-parser')

const miauthConfig = require('./config')

const forgotRoute = require('./routes/forgot.js')
const adminRoute = require('./routes/admin.js')
const userRoute = require('./routes/user.js')
const { handleError } = require('./utils/error')

const app = express()

app.set('view engine', 'html');
app.engine('html', require('hbs').__express);

app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }))

app.use('/forgot', forgotRoute)
app.use('/admin', adminRoute)
app.use('/auth', userRoute)

app.use(function(err, req, res, next) {
    if (err)
        handleError(err, res)
});

const HOST = '0.0.0.0'
const srv = app.listen(miauthConfig.PORT, HOST, () => {
    console.log(
        `▀▄▀▄▀▄ [ MiAuth started & listening on port ${miauthConfig.PORT} ] ▄▀▄▀▄▀`
    )
})

module.exports = srv
