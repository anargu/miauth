const Sequelize = require('sequelize')

const token = require('../utils/token')
const miauthConfig = require('../config')

module.exports = (sequelize) => {
    const Session = sequelize.define('session', {
        uuid: {
          primaryKey: true,
          type: Sequelize.UUID,
          defaultValue: Sequelize.UUIDV4,
          allowNull: false
        },
        issuedAt: {
            type: Sequelize.DATE,
            defaultValue: Sequelize.NOW
        },
        access_token: {
            type: Sequelize.STRING,
            allowNull: false
        },
        refresh_token: {
            type: Sequelize.STRING,
            allowNull: !miauthConfig.refresh_token.enabled
        },
        scope: {
            type: Sequelize.STRING,
            allowNull: true
        },
        expires_in: {
            type: Sequelize.STRING,
        },
    }, {
        // persist "deleted" sessions
        paranoid: true
    })

    Session.associate = function(models) {
        models.Session.belongsTo(models.User, { foreignKey: 'userId' })        
    }

    /**
    * Assign the project to an employee.
    * @param {Object} input - Input object which have required values to create a Session
    * @param {string} input.userId - The User ID which this Session are gonna be related with
    * @param {string} input.email - The User Email to inject in JWT payload as user_email.
    * @param {string} input.scope - [Optional] The scope value to inject in JWT payload as scope.
    */
    Session.createSession = async function(input) {
        if(!(input.userId && input.email)) {
            throw new Error('input does not contain needed parameters')
        }

        const access_token = await token.tokenize({
            userId: input.userId,
            user_email: input.email
        }, miauthConfig.access_token.secret, miauthConfig.access_token.expires_in)
        const expires_in = token.expirationOffset(miauthConfig.access_token.expires_in)
        // only if refresh was enabled
        let refresh_token
        if(miauthConfig.refresh_token.enabled) {
            refresh_token = await token.tokenize({
                userId: input.userId,
                user_email: input.email
            }, miauthConfig.refresh_token.secret)
        }
        
        const _session = await Session.create({
            userId: input.userId,
            scope: input.scope,
            access_token,
            expires_in,
            refresh_token
        })

        return _session.toJSON()
    }

    Session.revokeSession = async function(input) {
        if(!input.sessionId) {
            throw new Error('input does not contain sessionId field')
        }

        const _session = await Session.findOne({
            where: { uuid: sessionId }
        })

        await _session.destroy()

        return _session
    }

    /**
    * Assign the project to an employee.
    * @param {Object} input - Input object which have required values to revoke all Sessions of User
    * @param {string} input.userId - The User ID to look in DB for all his sessions
    */
    Session.revokeAllSessions = async function(input) {
        if(!input.userId) {
            throw new Error('input does not contain userId field')
        }

        const _deletedCount = await Session.destroy({
            where: { userId: input.userId }
        })

        return _deletedCount
    }

    Session.getAllSessions = async function() {
        return await Session.findAll({})
    }

    Session.getSessionById = async function(input) {
        if(!input.sessionId) {
            throw new Error('input does not contain sessionId field')
        }
        return await Session.findOne({ where: {
            uuid: sessionId
        }})
    }

    Session.refreshSession = async function(input) {
        if(!input.sessionId) {
            throw new Error('input does not contain sessionId field')
        }

        await Session.revokeSession(input)
        await Session.createSession(input)
    }

    return Session
}