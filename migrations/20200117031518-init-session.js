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
    await queryInterface.createTable('sessions', {
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
        allowNull: !miauthConfig.refresh
      },
      scope: {
        type: Sequelize.STRING,
        allowNull: true
      },
      expires_in: {
        type: Sequelize.STRING,
      },
      userId: {
        foreignKey: true,
        type: Sequelize.UUID,
        allowNull: false,
        references: {
          model: 'users',
          key: 'uuid',
        },
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
    queryInterface.dropTable('sessions')
  }
};
