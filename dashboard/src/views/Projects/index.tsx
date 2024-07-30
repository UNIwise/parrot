import SearchIcon from "@mui/icons-material/Search";
import {
  Box,
  FormControl,
  FormLabel,
  Input,
  Sheet,
  Table,
  Typography,
} from "@mui/joy";
import { useEffect, useState } from "react";
import { useGetProjects } from "../../api/hooks/useGetProjects";
import { TablePaginationSection } from "../../components/TablePaginationSection";

import { Placeholder } from '../../components/Placeholder';
import { GetProjectsResponse, Project } from "../../interfaces/projects";
import { ProjectTableRow } from "./components/TableRow";

const ITEMS_PER_PAGE = 20;

export const ProjectsOverview = () => {
  const [searchBar, setSearchBar] = useState("");
  const { data: projects, isLoading: isProjectsLoading } = useGetProjects();
  const [projectsList, setProjectsList] = useState<GetProjectsResponse>();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageCount, setPageCount] = useState(1);
  const [paginatedVersions, setPaginatedVersions] = useState<Project[]>();

  useEffect(() => {
    setProjectsList(projects);
  }, [projects]);

  const projectSearchHandle = (projectName: string) => {
    if (!projects) return;

    const filteredProjects = projects.projects.filter((project: Project) =>
      project.name.toLowerCase().includes(projectName.toLowerCase()),
    );
    setProjectsList({ projects: filteredProjects });
  };

  useEffect(() => {
    if (!projectsList || !projectsList.projects) return;

    const pageCount = Math.ceil(projectsList.projects.length / ITEMS_PER_PAGE);
    const paginatedVersions = projectsList.projects.slice(
      (currentPage - 1) * ITEMS_PER_PAGE,
      currentPage * ITEMS_PER_PAGE,
    );
    setPageCount(pageCount);
    setPaginatedVersions(paginatedVersions);
  }, [projectsList, currentPage]);

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  if (isProjectsLoading) {
    return <Placeholder />;
  }

  return (
    <>
      <Box
        sx={{
          borderRadius: "sm",
          py: 0.5,
          display: { xs: "none", sm: "flex" },
          flexWrap: "wrap",
          gap: 1.5,
          "& > *": {
            minWidth: { xs: "120px", md: "160px" },
          },
        }}
      >
        <Typography
          level="h2"
          component="h1"
          style={{
            alignSelf: "center",
            fontSize: "3rem",
            marginRight: "1.5rem",
          }}
        >
          Projects
        </Typography>

        <FormControl sx={{ flex: 1, pb: "1.1rem" }} size="sm">
          <FormLabel>Search for project</FormLabel>

          <Input
            size="sm"
            placeholder="Enter your favorite project name... Like FlowUI or WISEflow"
            startDecorator={<SearchIcon />}
            onChange={(e) => {
              setSearchBar(e.target.value);
              projectSearchHandle(e.target.value);
            }}
            value={searchBar}
          />
        </FormControl>
      </Box>

      <Sheet
        variant="outlined"
        sx={{
          display: { xs: "none", sm: "initial" },
          width: "100%",
          borderRadius: "sm",
          flexShrink: 1,
          overflow: "auto",
          minHeight: 0,
        }}
      >
        <Table
          aria-labelledby="tableTitle"
          stickyHeader
          hoverRow
          sx={{
            "--TableCell-headBackground":
              "var(--joy-palette-background-level1)",
            "--Table-headerUnderlineThickness": "1px",
            "--TableRow-hoverBackground":
              "var(--joy-palette-background-level1)",
            "--TableCell-paddingY": "4px",
            "--TableCell-paddingX": "8px",
          }}
        >
          <thead>
            <tr>
              <th style={{ width: 240, padding: "0.7rem 1.5rem" }}>Name</th>
              <th style={{ width: 100, padding: "0.7rem 0.5rem" }}>
                Number of versions
              </th>
              <th style={{ width: 340, padding: "0.7rem 0.5rem" }}>
                Created At
              </th>
              <th
                style={{
                  width: 140,
                  padding: "0.7rem 4.5rem",
                  textAlign: "end",
                }}
              >
              </th>
            </tr>
          </thead>

          {paginatedVersions && (
            <tbody style={{ fontSize: "1rem" }}>
              {paginatedVersions.map((project) => (
                <ProjectTableRow
                  projectId={project.id}
                  key={project.id}
                  projectName={project.name}
                  createdAt={project.createdAt}
                  numberOfVersions={project.numberOfVersions}
                />
              ))}
            </tbody>
          )}
        </Table>
      </Sheet>

      <TablePaginationSection
        currentPage={currentPage}
        pageCount={pageCount}
        onPageChange={handlePageChange}
      />
    </>
  );
};
