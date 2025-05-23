const path = require('path');
const dotenv = require('dotenv');

const appEnv = process.env.APP_ENV || 'development';

const envPath = appEnv === 'production'
  ? path.resolve(__dirname, '.env')
  : path.resolve(__dirname, '../../.env');

dotenv.config({ path: envPath });

const config = {
  APP_ENV: appEnv,

  CORE_API_URL: process.env.CORE_API_URL,
  DASHBOARD_URL: process.env.DASHBOARD_URL,
};

module.exports = {
  config,
};
