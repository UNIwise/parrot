import { Route, Routes as RouterRoutes } from 'react-router-dom';
import ProjectsOverview from '../views/ProjectsOverview';
import VersionOverview from '../views/VersionOverview';
import { Frame } from './Frame';

const Routes: React.FC = () => {
    return (
        <RouterRoutes>
            <Route path="he" element={<Frame />}>
                <Route path="projects" element={<ProjectsOverview />}>
                    <Route path="/:projectId/version/:versionId" element={<VersionOverview />} />
                </Route>
            </Route>
        </RouterRoutes>
    );
};

export default Routes;
