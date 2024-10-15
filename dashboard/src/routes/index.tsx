import { FC } from "react";
import { Navigate, Route, Routes as RouterRoutes } from "react-router-dom";
import { ProjectsOverview } from "../views/Projects";
import { VersionsOverview } from "../views/Projects/Versions";

const Routes: FC = () => {
  return (
    <RouterRoutes>
      <Route>
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
