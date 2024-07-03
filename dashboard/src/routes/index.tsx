import { FC } from "react";
import { Route, Routes as RouterRoutes } from "react-router-dom";
import { ProjectsOverview } from "../views/Projects";
import { VersionsOverview } from "../views/Projects/Versions";
import { Frame } from "./Frame";

const Routes: FC = () => {
  return (
    <RouterRoutes>
      <Route element={<Frame />}>
        <Route path="projects" element={<ProjectsOverview />}>
          <Route path=":projectId/versions" element={<VersionsOverview />} />
        </Route>
      </Route>
    </RouterRoutes>
  );
};

export default Routes;
