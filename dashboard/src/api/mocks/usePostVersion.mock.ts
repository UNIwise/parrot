import { Version } from "../../interfaces/versions";
import { mock } from "../client";
import { mockedVersionsResponse } from "./useGetVersions.mock";

mock.onPost(/^\/api\/v1\/projects\/\d+\/versions$/).reply((req) => {
  const request = JSON.parse(req.data);

  const newVersion: Version = {
    id: Math.floor(Math.random() * 1000),
    name: request.name,
    createdAt: new Date().toISOString(),
  };

  mockedVersionsResponse.versions.push(newVersion);

  return [201, newVersion];
});
