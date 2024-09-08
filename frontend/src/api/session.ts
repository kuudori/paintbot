import axios from 'axios';
import telegram from "@/utils/telegram";

const session = axios.create({
  baseURL: '/api/',
  headers: {
    'X-Auth': telegram.webapp.initData,
  }
});

export default session;
