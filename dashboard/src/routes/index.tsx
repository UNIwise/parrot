import { FC } from 'react';
import { Route, Routes as RouterRoutes } from 'react-router-dom';
import ProjectsOverview from '../views/ProjectsOverview';
import VersionsOverview from '../views/ProjectsOverview/VersionsOverview';
import { Frame } from './Frame';

const Routes: FC = () => {
    return (
        <RouterRoutes>
            <Route path="he" element={<Frame />}>
                <Route path="projects" element={<ProjectsOverview />}>
                    <Route path="/:projectId/version/:versionId" element={<VersionsOverview />} />
                </Route>
            </Route>
        </RouterRoutes>
    );
};

export default Routes;
