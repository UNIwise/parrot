import { useMutation, useQueryClient } from "@tanstack/react-query";
import { postProjectRequest } from "../../interfaces/projects";
import { GetVersionsResponse } from "../../interfaces/versions";
import client from "../client";
import "../mocks/usePostVersion.mock";

const postVersion = async (projectId: number, request: postProjectRequest) => {
  const response = await client.post<GetVersionsResponse>(
    `/api/v1/projects/${projectId}/versions`,
    request,
  );

  return response.data;
};

export const usePostVersion = (projectId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    onSuccess: async () => {
      await queryClient.invalidateQueries(
        {
          queryKey: ["api", "projects", projectId, "versions"],
          exact: true,
        },
        { throwOnError: true },
      );
    },
    mutationFn: (request: postProjectRequest) =>
      postVersion(projectId, request),
  });
};
