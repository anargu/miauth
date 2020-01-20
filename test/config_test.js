const assert = require('assert')
const fs = require('fs')
const YAML = require('yaml')

describe('Reading miauth.config.yml', function() {
    const testConfigFile = 'test.config.yml'
    
    it('should return the name field of yml file', function() {
        try {
            const data = fs.readFileSync(`./test/${testConfigFile}`, { encoding: 'utf-8' })
            const configData = YAML.parse(data)
            assert.equal(configData.name, 'HelloMyMiauthImplementation') 
        } catch (error) {
            assert.fail(error)
        }
    })
})
