import SearchIcon from "@mui/icons-material/Search";
import { FormControl, Input, Sheet, Table } from "@mui/joy";
import { useMemo, useState } from "react";
import { useGetProjects } from "../../api/hooks/useGetProjects";
import { Header } from "../../components/Header";
import { Placeholder } from "../../components/Placeholder";
import { TablePaginationSection } from "../../components/TablePaginationSection";
import { Project } from "../../interfaces/projects";
import { ProjectTableRow } from "./components/TableRow";

const ITEMS_PER_PAGE = 25;

export const ProjectsOverview = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const { data: projects, isLoading } = useGetProjects();
  const [currentPage, setCurrentPage] = useState(1);

  const filteredProjects: Project[] = useMemo(() => {
    if (!projects || !projects.projects || projects.projects.length === 0)
      return [];

    return projects.projects.filter((project: Project) =>
      project.name.toLowerCase().includes(searchTerm.toLowerCase()),
    );
  }, [projects, searchTerm]);

  const pageCount = filteredProjects
    ? Math.ceil(filteredProjects.length / ITEMS_PER_PAGE)
    : 0;

  const paginatedProjects = useMemo(() => {
    if (filteredProjects.length === 0) return [];

    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredProjects.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredProjects, currentPage]);

  if (isLoading) return <Placeholder />;

  return (
    <>
      <Header items={[{ name: "Projects" }]} />

      <FormControl sx={{ flex: 1, pb: "1.1rem" }}>
        <Input
          placeholder="search for project..."
          startDecorator={<SearchIcon />}
          onChange={(e) => setSearchTerm(e.target.value)}
          value={searchTerm}
        />
      </FormControl>

      <Sheet
        variant="outlined"
        sx={{
          display: { xs: "initial", sm: "initial" },
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
              ></th>
            </tr>
          </thead>
          <tbody>
            {paginatedProjects.map((project) => (
              <ProjectTableRow
                key={project.id}
                projectId={project.id}
                projectName={project.name}
                createdAt={project.createdAt}
                numberOfVersions={project.numberOfVersions}
              />
            ))}
          </tbody>
        </Table>
      </Sheet>

      {filteredProjects.length > ITEMS_PER_PAGE && (
        <TablePaginationSection
          currentPage={currentPage}
          pageCount={pageCount}
          onPageChange={setCurrentPage}
        />
      )}
    </>
  );
};
