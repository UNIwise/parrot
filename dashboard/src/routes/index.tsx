import { FC } from 'react';
import { Route, Routes as RouterRoutes } from 'react-router-dom';
import ProjectsOverview from '../views/ProjectsOverview';
import VersionsOverview from '../views/VersionsOverview';
import { Frame } from './Frame';

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
