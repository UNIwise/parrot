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
import { mockedProjectsResponse } from "../../api/mocks/projects.mock";
import { PaginationSection } from "../../components/TablePaginationSection";
import { TableRow } from "../../components/TableRow";
import { getProjectsResponse } from "../../interfaces/projects";

export const ProjectsOverview = () => {
  const [searchBar, setSearchBar] = useState("");
  //TODO: replace mocked data with the response from the API when react query hooks are implemented
  const projects = mockedProjectsResponse;
  const [projectsList, setProjectsList] = useState<getProjectsResponse>(
    mockedProjectsResponse,
  );

  useEffect(() => {
    if (searchBar === "") {
      return;
    }
    setProjectsList(projectsList);
  }, [projectsList, searchBar]);

  const projectSearchHandle = (projectName: string) => {
    const filteredProjects = projects.projects.filter((project) =>
      project.name.toLowerCase().includes(projectName.toLowerCase()),
    );
    setProjectsList({ projects: filteredProjects });
  };

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
        <Typography level="h2" component="h1">
          Projects
        </Typography>

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

          <tbody>
            {projectsList.projects.map((project) => (
              <TableRow
                key={project.id}
                name={project.name}
                createdAt={project.createdAt}
                numberOfVersions={project.numberOfVersions}
              />
            ))}
          </tbody>
        </Table>
      </Sheet>

      <PaginationSection />
    </>
  );
};
