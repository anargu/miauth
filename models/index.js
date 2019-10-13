const User = require('./user')

function initializeModels(re) {
    User.setup(re)
}

module.exports = {
    initializeModels,
    User
}
