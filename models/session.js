const Sequelize = require('sequelize')
const DataTypes = require('sequelize').DataTypes

const token = require('../utils/token')
const miauthConfig = require('../config')

module.exports = (sequelize) => {
    const Session = sequelize.define('session', {
        uuid: {
          primaryKey: true,
          type: Sequelize.UUID,
          defaultValue: DataTypes.UUIDV4,
          allowNull: false
        },
        issuedAt: {
            type: Sequelize.DATE,
            defaultValue: DataTypes.NOW
        },
        access_token: {
            type: Sequelize.STRING,
            allowNull: false
        },
        refresh_token: {
            type: Sequelize.STRING,
            allowNull: !miauthConfig.refresh
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
        models.Session.belongsTo(models.User)        
    }


    // TODO: validate what input params should be at creating a new session
    Session.createSession = async function(input) {
        if(!input.userId) {
            throw new Error('input does not contain needed parameters')
        }

        const tokenResult = token.tokenize(null, {...input})
        const access_token = tokenResult['access_token']
        const expires_in = tokenResult['expires_in']
        // only if refresh was enabled
        let refresh_token
        if(miauthConfig.refresh) {
            refresh_token = tokenResult['refresh_token']
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