import { useMutation, useQueryClient } from '@tanstack/react-query';
import { GetVersionsResponse } from '../../interfaces/versions';
import client from '../client';
import '../mocks/useGetVersions.mock';

const postVersion = async (projectId: number) => {
  const response = await client.post<GetVersionsResponse>(
    `/api/v1/projects/${projectId}/versions`,
  );

  return response.data;
};

export const usePostVersion = (projectId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ['api', 'projects', projectId, 'versions'],
        exact: true,
      },
        { throwOnError: true });
    },
    mutationFn: () => postVersion(projectId),
  });
};


