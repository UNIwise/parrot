import { AxiosInstance } from 'axios';
import * as JWT from '@uniwise/jwt';

function registerAuthorizationInterceptor(inst: AxiosInstance) {
  inst.interceptors.request.use((config) => {
    if (!config.headers) {
      config.headers = {};
    }
    config.headers['X-csrftoken'] = JWT.get()?.wiseflow.userInfo.csrfId || 0;
    config.headers['Authorization'] = `Bearer ${JWT.getString()}`;
    return config;
  });
}

export default registerAuthorizationInterceptor;
