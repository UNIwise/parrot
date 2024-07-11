import { useQuery } from '@tanstack/react-query';
import { GetProjectsResponse } from '../../interfaces/projects';
import client from '../client';
import '../mocks/projects.mock';

const getProjects = async () => {
  const response = await client.get<GetProjectsResponse>(
    `/parrot/v1/projects`,
    {
      // TODO: replace with the real one
      baseURL: '',
    },
  );

  return response.data;
};

export const useGetProjects = () => {
  return useQuery({
    queryKey: ['parrot', 'projects'],
    queryFn: () => getProjects(),
  });
};
