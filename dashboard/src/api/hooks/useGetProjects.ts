import { useQuery } from '@tanstack/react-query';
import { GetProjectsResponse } from '../../interfaces/projects';
import client from '../client';
import '../mocks/useGetProjects.mock';

const getProjects = async () => {
  const response = await client.get<GetProjectsResponse>(
    `/api/v1/projects`,
  );

  return response.data;
};

export const useGetProjects = () => {
  return useQuery({
    queryKey: ['api', 'projects'],
    queryFn: () => getProjects(),
  });
};
