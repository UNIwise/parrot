import { mock } from "../client";
import { mockedVersionsResponse } from "./useGetVersions.mock";

mock.onPost(/^\/api\/v1\/projects\/\d+\/versions$/).reply((req) => {
  const request = JSON.parse(req.data);
  const projectId = parseInt(req.url!.split("/")[4]);

  const newVersion = {
    id: Math.floor(Math.random() * 1000),
    name: request.name,
    projectId: projectId,
    version: request.version,
    createdAt: new Date().toISOString(),
  };

  mockedVersionsResponse.versions.push(newVersion);

  return [201, newVersion];
});

