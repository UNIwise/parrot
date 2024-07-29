import { useQuery } from '@tanstack/react-query';
import { Project } from '../../interfaces/projects';
import client from '../client';
import '../mocks/useGetProject.mock';

const getProject = async (projectId?: number) => {
  const response = await client.get<Project>(
    `/api/v1/projects/${projectId}`,
  );

  return response.data;
};

export const useGetProject = (projectId?: number) => {
  return useQuery({
    queryKey: ['api', 'projects', projectId],
    queryFn: () => getProject(projectId),
  });
};
