const Sequelize = require('sequelize')

const miauthConfig = require('../config')
const hashPassword = require('../utils/auth').hashPassword

module.exports = (sequelize) => {
    const User = sequelize.define('user', {
        uuid: {
          primaryKey: true,
          type: Sequelize.UUID,
          defaultValue: Sequelize.UUIDV4,
          allowNull: false
        },
        username: {
            type: Sequelize.STRING,
            unique: miauthConfig.user.username,
            allowNull: !miauthConfig.user.username,
            validate: miauthConfig.user.username
            ? ({
                is: new RegExp(miauthConfig.field_validations.username.pattern, 'g'),
                len: {
                    args: [...miauthConfig.field_validations.username.len],
                    msg: miauthConfig.field_validations.username.invalid_pattern_error_message
                }
            })
            : undefined
        },
        email: {
            type: Sequelize.STRING,
            unique: miauthConfig.user.username,
            allowNull: !miauthConfig.user.username,
            validate: miauthConfig.user.username
            ? ({
              is: new RegExp(miauthConfig.field_validations.email.pattern, 'g'),
              len: {
                  args: [...miauthConfig.field_validations.email.len],
                  msg: miauthConfig.field_validations.email.invalid_pattern_error_message
              }
            })
            : undefined
        },
        hash: {
            type: Sequelize.STRING,
            allowNull: false,
        }
    }, {
        indexes: (() => {
            let _indexes = []
            if(miauthConfig.user.username) {
                _indexes.push({ unique: true, fields: ['username'] })
            }
            if(miauthConfig.user.email) {
                _indexes.push({ unique: true, fields: ['email'] })
            }
            return _indexes
        })(),
        // persist "deleted" sessions
        paranoid: true
    })

    User.associate = function(models) {
        models.User.hasMany(models.Session, {
            foreignKey: 'userId'
        })
    }

    User.createUser = async (input) => {
        if(!((input.username || input.email) && input.password)) {
            throw new Error('input does not contain needed parameters')
        }
    
        const _user = await User.create({
            username: input.username,
            email: input.email,
            hash: await hashPassword(input.password)
        })
        return _user.toJSON()
    }

    /**
    * Assign the project to an employee.
    * @param {Object} userQueryField - Object which contain field & value props to defined the query params for User
    * @param {string} userQueryField.field - Field name from query that will be used to find a User
    * * @param {string} userQueryField.value - Value of the field that will be used to find a User
    */
    User.updatePassword = async (userQueryField, newPassword) => {
        if(!(newPassword)) {
            throw new Error('input does not contain needed parameters')
        } else if(!userQueryField) {
            throw new Error('user field filter null or undefined. Then, how can I query for the user?')
        } else if (!(userQueryField['field'] && userQueryField['value'])) {
            throw new Error('user query values incompleted. How can I query for the user?')
        } else if(['username', 'email', 'userId'].indexOf(userQueryField['field']) === -1) {
            throw new Error(`Unsupported user query field: ${userQueryField['field']}`)
        } else if (userQueryField['field'] === 'userId') {
            userQueryField['field'] = 'uuid'
        }
    
        const whereQuery = {}
        whereQuery[userQueryField['field']] = userQueryField['value']
    
        const _user = await User.findOne({
            where: { ...whereQuery }
        })
        if (_user === null) {
            throw new Error('user is not found')
        }
    
        _user.hash = await hashPassword(newPassword)
        await _user.save();
    
        return _user.toJSON()
    }
    
    User.findByUsername = async (username) => {
        const _user = await User.findOne({
            where: { username: username, }
        })

        if (_user === null)
            throw new Error('user is not found')            
        return _user.toJSON()
    }

    User.findByEmail = async (email) => {
        const _user = await User.findOne({
            where: { email: email, }
        })

        if (_user === null)
            throw new Error('user is not found')            
        return _user.toJSON()
    }
    
    return User
}
