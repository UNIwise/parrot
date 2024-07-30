import { mock } from "../client";
import { mockedProjectsResponse } from "./useGetProjects.mock";

mock.onGet(/^\/api\/v1\/projects\/\d+$/).reply((req) => {
  const projectId = parseInt(req.url!.split("/")[4]);

  const project = mockedProjectsResponse.projects.find(
    (project) => project.id === projectId
  );

  return [200, project];
});
