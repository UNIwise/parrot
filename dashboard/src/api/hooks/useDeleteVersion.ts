import { useMutation, useQueryClient } from '@tanstack/react-query';
import { GetVersionsResponse } from '../../interfaces/versions';
import client from '../client';
import '../mocks/useDeleteVersion.mock';

const deleteVersion = async (projectId?: number, versionId?: number) => {
  const response = await client.delete<GetVersionsResponse>(
    `/api/v1/projects/${projectId}/versions/${versionId}`,
  );

  return response.data;
};

export const useDeleteVersion = (projectId?: number, versionId?: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ['api', 'projects', projectId, 'versions'],
        exact: true,
      },
        { throwOnError: true });
    },
    mutationFn: () => deleteVersion(projectId, versionId),
  });
};


