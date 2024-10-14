import SearchIcon from "@mui/icons-material/Search";
import { FormControl, FormLabel, Input, Sheet, Table } from "@mui/joy";
import { useMemo, useState } from "react";
import { useGetPageParams } from "../../../api/hooks/useGetPageParams";
import { useGetProject } from "../../../api/hooks/useGetProject";
import { useGetVersions } from "../../../api/hooks/useGetVersions";
import { Header } from "../../../components/Header";
import { ManageVersionModal } from "../../../components/ManageVersionModal";
import { Placeholder } from "../../../components/Placeholder";
import { TablePaginationSection } from "../../../components/TablePaginationSection";
import { Version } from "../../../interfaces/versions";
import { VersionTableRow } from "./components/VersionTableRow";

const ITEMS_PER_PAGE = 20;

export const VersionsOverview = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const { projectId } = useGetPageParams();
  const { data: project, isLoading: isProjectLoading } =
    useGetProject(projectId);
  const { data: versionsData, isLoading: isVersionsDataLoading } =
    useGetVersions(projectId);

  const filteredVersions: Version[] = useMemo(() => {
    if (
      !versionsData ||
      !versionsData.versions ||
      versionsData.versions.length === 0
    )
      return [];

    return versionsData.versions.filter((version: Version) =>
      version.name.toLowerCase().includes(searchTerm.toLowerCase()),
    );
  }, [versionsData, searchTerm]);

  const pageCount =
    filteredVersions.length > 0
      ? Math.ceil(filteredVersions.length / ITEMS_PER_PAGE)
      : 0;

  const paginatedVersions = useMemo(() => {
    if (filteredVersions.length === 0) return [];

    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredVersions.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredVersions, currentPage]);

  if (isProjectLoading || isVersionsDataLoading) {
    return <Placeholder />;
  }

  return (
    <>
      {/* <Typography
          level="h2"
          component="h1"
          sx={{
            alignSelf: "center",
            fontSize: "2rem",
            color: (t) => t.palette.primary[400],
            m: "0 1.5rem 2rem 0",
            border: "1px solid",
            borderColor: (t) => t.palette.primary[400],
            p: "1rem 2.5rem",
            borderRadius: "sm",
          }}
        >
          {project?.name ? `${project.name} versions` : "Versions"}
        </Typography> */}

      <Header
        items={[
          {
            name: "Projects",
            to: "/projects",
          },
          {
            name: project?.name || "",
            to: `/projects/${projectId}/versions`,
          },
        ]}
      />

      <FormControl sx={{ flex: 1, pb: "1rem" }} size="sm">
        <FormLabel>Search for version</FormLabel>

        <Input
          size="sm"
          placeholder="Look up for version name"
          startDecorator={<SearchIcon />}
          onChange={(e) => setSearchTerm(e.target.value)}
          value={searchTerm}
        />
      </FormControl>

      {projectId && <ManageVersionModal projectId={projectId} />}

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

          {projectId && (
            <tbody>
              {paginatedVersions.map((version: Version) => (
                <VersionTableRow
                  key={version.id}
                  projectId={projectId}
                  versionId={version.id}
                  versionName={version.name}
                  createdAt={version.createdAt}
                />
              ))}
            </tbody>
          )}
        </Table>
      </Sheet>

      <TablePaginationSection
        currentPage={currentPage}
        pageCount={pageCount}
        onPageChange={setCurrentPage}
      />
    </>
  );
};
