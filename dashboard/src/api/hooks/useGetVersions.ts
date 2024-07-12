import { useQuery } from '@tanstack/react-query';
import { GetVersionsResponse } from '../../interfaces/versions';
import client from '../client';
import '../mocks/useGetVersions.mock';

const getVersions = async (projectId?: string) => {
  const response = await client.get<GetVersionsResponse>(
    `/api/v1/projects/${projectId}/versions`,
  );

  return response.data;
};

export const useGetVersions = (projectId?: string) => {
  return useQuery({
    queryKey: ['api', 'projects', projectId, 'versions'],
    queryFn: () => getVersions(projectId),
  });
};
