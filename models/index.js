const User = require('./user')

function initializeModels (sequelize) {
    User.init(sequelize)
}

module.exports = {
    initializeModels,
    User
}
