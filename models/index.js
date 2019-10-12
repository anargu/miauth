const User = require('./user')

async function initializeModels (sequelize) {
    await User.init(sequelize)
}

module.exports = {
    initializeModels,
    User
}
