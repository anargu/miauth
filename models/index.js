const User = require('./user')

function initializeModels () {
    User.init()
}

module.exports = {
    initializeModels
}
