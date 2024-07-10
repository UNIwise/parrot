import { useQuery } from '@tanstack/react-query';
import { ENV } from '../../constants/env';
import { GetProjectsResponse } from '../../interfaces/projects';
import './useGetApplicationCategories.mock';


// Get application categories that have been defined on the license
const getApplicationCategories = async () => {
  const response = await client.get<GetProjectsResponse[]>(
    `/device-monitor/v1/admin/settings/application-categories`,
    {
      // TODO: replave with the real one
      baseURL: ENV.MOCKED,
    },
  );

  return response.data;
};

export const useGetApplicationCategories = () => {
  return useQuery(
    ['device-monitor', 'v1', 'admin', 'settings', 'application-categories'],
    () => getApplicationCategories(),
  );
};
