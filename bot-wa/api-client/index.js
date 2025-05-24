const axios = require('axios')

// Base URL dari .env atau default fallback
const BASE_URL = process.env.CORE_API_URL;

const apiClient = axios.create({
  baseURL: BASE_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
});

const getDailyTransaction = async (phoneNumber) => {
  try {
    const response = await apiClient.get(`/users/${phoneNumber}/transactions/daily`);
    return response.data;
  } catch (error) {
    console.error(error.message);
    throw error;
  }
};

const getMonthlyTransaction = async (phoneNumber) => {
  try {
    const response = await apiClient.get(`/users/${phoneNumber}/transactions/monthly`);
    return response.data;
  } catch (error) {
    console.error(error.message);
    throw error;
  }
};

const registerUser = async (phoneNumber, fullName) => {
  try {
    const response = await apiClient.post(`/users/register`, {
      chat_id: phoneNumber,
      name: fullName
    });
    return response.data;
  } catch (error) {
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data;
      
      if (status >= 400 && status < 500) {
        return data;
      }
    }
    
    console.error(error.message);
    throw error;
  }
};

const checkUser = async (phoneNumber) => {
  try {
    const response = await apiClient.get(`/users/${phoneNumber}/exists`);
    return response.data;
  } catch (error) {
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data;
      
      if (status >= 400 && status < 500) {
        return data;
      }
    }
    
    console.error(error.message);
    throw error;
  }
}

const hitAiClassifyTransaction = async (phoneNumber, prompt) => {
  try {
    const response = await apiClient.post(`/users/${phoneNumber}/transactions/ai-classify`, {
      prompt
    });
    return response.data;
  } catch (error) {
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data;

      if (status >= 400 && status < 500) {
        return data; // kirim response error dari server sebagai hasil
      }
    }

    console.error(error.message);
    throw error;
  }
};


const saveTransaction = async (transaction) => {
  try {
    const response = await apiClient.post(`/transactions`, transaction);
    return response.data;
  } catch (error) {
    console.error(error.message);
    throw error;
  }
}

module.exports = {
  apiClient,
  getDailyTransaction,
  getMonthlyTransaction,
  registerUser,
  checkUser,
  hitAiClassifyTransaction,
  saveTransaction
};
