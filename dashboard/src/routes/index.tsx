import { FC } from "react";
import { Navigate, Route, Routes as RouterRoutes } from "react-router-dom";
import { ProjectsOverview } from "../views/Projects";
import { VersionsOverview } from "../views/Projects/Versions";
import { Frame } from "./Frame";

const Routes: FC = () => {
  return (
    <RouterRoutes>
      <Route element={<Frame />}>
        <Route path="projects">
          <Route index element={<ProjectsOverview />} />
          <Route path=":projectId/versions" element={<VersionsOverview />} />
        </Route>

        <Route path="*" element={<Navigate to="projects" />} />
      </Route>
    </RouterRoutes>
  );
};

export default Routes;
