import * as JWT from '@uniwise/jwt';
import { AxiosInstance, AxiosHeaders } from 'axios';

function registerAuthorizationInterceptor(inst: AxiosInstance) {
  inst.interceptors.request.use((config) => {
    if (!config.headers) {
      config.headers = new AxiosHeaders();
    }
    config.headers['X-csrftoken'] = JWT.get()?.wiseflow.userInfo.csrfId || 0;
    config.headers['Authorization'] = `Bearer ${JWT.getString()}`;
    return config;
  });
}

export default registerAuthorizationInterceptor;
