import { mock } from "../client";
import { mockedVersionsResponse } from "./useGetVersions.mock";

mock.onDelete(/^\/api\/v1\/projects\/\d+\/versions\/\d+$/).reply((req) => {
  const versionId = parseInt(req.url!.split("/")[6]);

  const versionIndex = mockedVersionsResponse.versions.findIndex(
    (version) => version.id === versionId
  );

  mockedVersionsResponse.versions.splice(versionIndex, 1);

  return [204];
});
