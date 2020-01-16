const Sequelize = require('sequelize')
const DataTypes = require('sequelize').DataTypes

const miauthConfig = require('../config')
const hashPassword = require('../utils/auth')

module.exports = (sequelize) => {
    const User = sequelize.define('user', {
        uuid: {
          primaryKey: true,
          type: Sequelize.UUID,
          defaultValue: DataTypes.UUIDV4,
          allowNull: false
        },
        username: {
            type: Sequelize.STRING,
            unique: miauthConfig.user.username,
            allowNull: !miauthConfig.user.username,
            validate: miauthConfig.user.username
            ? ({
                is: new RegExp(miauthConfig.field_validations.username, 'g')
            })
            : undefined
        },
        email: {
            type: Sequelize.STRING,
            unique: miauthConfig.user.username,
            allowNull: !miauthConfig.user.username,
            validate: miauthConfig.user.username
            ? ({
              isEmail: true,
            })
            : undefined
        },
        hash: {
            type: Sequelize.STRING,
            set(plainTextPassword) {
                this.setDataValue('hash', hashPassword(plainTextPassword))
            },
            allowNull: false,
            validate: {
                len: [...miauthConfig.field_validations.password.len]
            }
        }
    }, {
        // Using `unique: true` in an attribute above is exactly the same as creating the index in the model's options:
        indexes: [
            { unique: true, fields: ['username'] },
            { unique: true, fields: ['email'] }
        ]
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
            hash: input.password
        })
        return _user.toJSON()
    }

    User.updatePassword = async (userFieldFilter, newPassword) => {
        if(!(newPassword)) {
            throw new Error('input does not contain needed parameters')
        }
        if (!userFieldFilter.username && userFieldFilter.email) {
            throw new Error('user field filter null or undefined. How can I query for the user?')
        }
    
        const whereQuery = {}
        const fieldFilter = !userFieldFilter.username ? 'email' :'username'
        const valueFilter = userFieldFilter.username || userFieldFilter.password
        whereQuery[fieldFilter] = valueFilter
    
        const _user = await User.findOne({
            where: { ...whereQuery }
        })
        if (_user === null) {
            throw new Error('user is not found')
        }
    
        _user.hash = newPassword
        await _user.save();
    
        return _user.toJSON()
    }
    
    User.findByUsername = async (username) => {
        const _user = await User.findOne({
            where: {
                username: username,
            }
        })

        if (_user === null) {
            throw new Error('user is not found')            
        }
        return _user
    }
    
    return User
}
