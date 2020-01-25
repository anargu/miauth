'use strict';

const miauthConfig = require('../config');

module.exports = {
  up: async (queryInterface, Sequelize) => {
    /*
    Add altering commands here.
    Return a promise to correctly handle asynchronicity.
    
    Example:
    return queryInterface.createTable('users', { id: Sequelize.INTEGER });
    */
    await queryInterface.createTable('users', {
      uuid: {
        primaryKey: true,
        type: Sequelize.UUID,
        defaultValue: Sequelize.UUIDV4,
        allowNull: false
      },
      createdAt: {
        allowNull: false,
        type: Sequelize.DATE
      },
      updatedAt: {
        allowNull: false,
        type: Sequelize.DATE
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
        unique: miauthConfig.user.email,
        allowNull: !miauthConfig.user.email,
        validate: miauthConfig.user.email
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
      },
      deletedAt: {
          type: Sequelize.DATE,
          allowNull: true,
          validate: {}
      }
    })
  },
  
  down: (queryInterface, Sequelize) => {
    /*
    Add reverting commands here.
    Return a promise to correctly handle asynchronicity.
    
    Example:
    return queryInterface.dropTable('users');
    */
    queryInterface.dropTable('users')
  }
};
