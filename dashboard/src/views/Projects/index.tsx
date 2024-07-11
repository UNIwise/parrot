import SearchIcon from '@mui/icons-material/Search';
import { Box, FormControl, FormLabel, Input, Sheet, Table } from "@mui/joy";
import { useEffect, useState } from "react";

import { useGetProjects } from '../../api/hooks/useGetProjects';
import { GetProjectsResponse } from "../../interfaces/projects";
import { PaginationSection } from "./components/PaginationSection/PaginationSection";
import { ProjectTableRow } from "./components/Row";


export const ProjectsOverview = () => {
  const [searchBar, setSearchBar] = useState("");
  //TODO: replace mocked data with the response from the API when react query hooks are implemented
  const { data: projects } = useGetProjects();
  const [projectsList, setProjectsList] = useState<GetProjectsResponse>();

  console.log(projects);

  // useEffect(() => {
  //   if (searchBar === "") {
  //     return;
  //   }

  //   if (!projects) {
  //     return;
  //   }

  //   setProjectsList((prevList) => {
  //     if (prevList !== projectsList) {
  //       return projectsList;
  //     }
  //     return prevList;
  //   });
  // }, [projectsList, searchBar, projects]);

  useEffect(() => {
    setProjectsList(projects);
  }, [projects]);

  const projectSearchHandle = (projectName: string) => {
    if (!projects) {
      return;
    }

    const filteredProjects = projects.projects.filter((project) =>
      project.name.toLowerCase().includes(projectName.toLowerCase()),
    );
    setProjectsList({ projects: filteredProjects });
  };

  return (
    <>
      <Box
        className="SearchAndFilters-tabletUp"
        sx={{
          borderRadius: "sm",
          py: 2,
          display: { xs: "none", sm: "flex" },
          flexWrap: "wrap",
          gap: 1.5,
          "& > *": {
            minWidth: { xs: "120px", md: "160px" },
          },
        }}
      >
        <FormControl sx={{ flex: 1 }} size="sm">
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
        className="OrderTableContainer"
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
                Delete
              </th>
            </tr>
          </thead>

          {projectsList && (
            <tbody>
              {projectsList.projects.map((project) => (
                <ProjectTableRow key={project.id} projectInfo={project} />
              ))}
            </tbody>
          )}
        </Table>
      </Sheet>

      <PaginationSection />
    </>
  );
};
