import { ProjectsOverviewTable } from "./components/Table";
import { PaginationSection } from "./components/Table/PaginationSection";


export const ProjectsOverview = () => {
  return (
    <>
      <ProjectsOverviewTable />
      <PaginationSection />
    </>
  );
};

