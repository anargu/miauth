const axios = require('axios')

module.exports = {
    sendResetPasswordEmail: async (miauthConfig, userData, linkUrl) => {
        if(
            miauthConfig.reset_password.mail_service['dosmj'] === undefined ||
            miauthConfig.reset_password.mail_service['dosmj'] === null) {
            throw new Error('No service known provided. Miauth is still in beta and is dependent from dosmj service')
        }
        try {
            const mailService = miauthConfig.reset_password.mail_service.dosmj

            const httpMethod = mailService.method
            const mailSenderUrl = mailService.endpoint
            const payload = mailService.payload
            
            payload.template_data.reset_link = linkUrl
            payload.email_specs.to = [
                { name: userData.email, email: userData.email }
            ]
                
            let response = await axios({
                method: httpMethod.toLowerCase(),
                url: mailSenderUrl,
                data: { ...payload }
            })
            return response
        } catch (error) {
            throw error
        }
    }
}